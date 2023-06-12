package mdedit

import (
	"errors"
	"image"
	"image/color"
	"io/fs"
	"log"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type FS interface {
	HomeDir() string
	WorkingDir() string
	ReadDir(fpath string) ([]fs.FileInfo, error)
	ReadFile(fpath string) ([]byte, error)
	WriteFile(fpath string, data []byte) error
}

type Session struct {
	fsys      FS
	win       *app.Window
	tabs      []tab
	tabList   layout.List
	activeTab int
}

type tab struct {
	btn     widget.Clickable
	content tabContent
}

func NewSession(fsys FS, win *app.Window) Session {
	return Session{
		fsys:    fsys,
		win:     win,
		tabList: layout.List{Axis: layout.Vertical},
	}
}

func (s *Session) Layout(gtx C, th *material.Theme) D {
	paint.Fill(gtx.Ops, th.Bg)
	if len(s.tabs) == 0 {
		return layout.Center.Layout(gtx, func(gtx C) D {
			return material.Body1(th, "Use Ctrl-O to open a file!").Layout(gtx)
		})
	}
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = 220
			gtx.Constraints.Max.X = 220
			return s.tabList.Layout(gtx, len(s.tabs), func(gtx C, i int) D {
				t := &s.tabs[i]
				if t.btn.Clicked() {
					s.SelectTab(i)
					op.InvalidateOp{}.Add(gtx.Ops)
				}
				// If this is the active tab, emphasize the text and invert the bg & fg.
				lbl := material.Body1(th, t.content.title())
				lbl.Font.Variant = "Mono"
				bg := th.Bg
				if i == s.activeTab {
					lbl.Font.Weight = font.Bold
					lbl.Color = bg
					bg = th.ContrastBg
				}
				// Record the layout in order to get the size for filling the background.
				m := op.Record(gtx.Ops)
				dims := t.btn.Layout(gtx, func(gtx C) D {
					return layout.UniformInset(5).Layout(gtx, lbl.Layout)
				})
				call := m.Stop()
				// Fill the background and draw the tab button.
				rect := clip.Rect{Max: dims.Size}
				paint.FillShape(gtx.Ops, bg, rect.Op())
				call.Add(gtx.Ops)
				return dims
			})
		}),
		layout.Rigid(func(gtx C) D {
			size := image.Point{1, gtx.Constraints.Max.Y}
			rect := clip.Rect{Max: size}.Op()
			paint.FillShape(gtx.Ops, th.Fg, rect)
			return D{Size: size}
		}),
		layout.Flexed(1, func(gtx C) D {
			return s.layTabContent(gtx, th, s.tabs[s.activeTab].content)
		}),
	)
}

func (s *Session) HandleKeyEvent(e key.Event) {
	if e.State != key.Press {
		return
	}
	switch e.Modifiers {
	case key.ModCtrl:
		switch e.Name {
		case "O":
			s.OpenFileExplorerTab()
		case "W":
			s.CloseActiveTab()
		case key.NameTab:
			s.NextTab()
		}
	case key.ModCtrl | key.ModShift:
		switch e.Name {
		case key.NamePageUp:
			s.SwapTabUp()
		case key.NamePageDown:
			s.SwapTabDown()
		case key.NameTab:
			s.PrevTab()
		}
	case key.ModAlt:
		if strings.Contains("123456789", e.Name) {
			if n, err := strconv.Atoi(e.Name); err == nil {
				s.SelectTab(n - 1)
			}
		}
	}
}

func (s *Session) layTabContent(gtx C, th *material.Theme, t tabContent) D {
	switch t := t.(type) {
	case *markdownTab:
		return s.layMarkdownTab(gtx, th, t)
	case *explorerTab:
		return s.layExplorerTab(gtx, th, t)
	default:
		return D{}
	}
}

func (s *Session) layMarkdownTab(gtx C, th *material.Theme, t *markdownTab) D {
	if t.view.Editor.SaveRequested() {
		go s.writeFile(t.name, t.view.Editor.Text())
	}
	return ViewStyle{
		Theme:      th,
		EditorFont: font.Font{Variant: "Mono"},
		Palette: Palette{
			Fg:         th.Fg,
			Bg:         th.Bg,
			LineNumber: color.NRGBA{200, 180, 4, 125},
			Heading:    color.NRGBA{200, 193, 255, 255},
			ListMarker: color.NRGBA{10, 190, 240, 255},
			BlockQuote: color.NRGBA{165, 165, 165, 230},
			CodeBlock:  color.NRGBA{162, 120, 70, 255},
		},
		View: &t.view,
	}.Layout(gtx)
}

func (s *Session) layExplorerTab(gtx C, th *material.Theme, t *explorerTab) D {
	for _, e := range t.expl.Events() {
		switch e := e.(type) {
		case DirChosenEvent:
			go s.openExplorerDir(t, e.Path)
		case FilesChosenEvent:
			go func() {
				for i, fpath := range e.Paths {
					s.OpenFile(fpath)
					if i == 0 {
						s.CloseActiveTab()
						s.win.Invalidate()
					}
				}
			}()
		}
	}
	return t.expl.Layout(gtx, th)
}

func (s *Session) OpenFileExplorerTab() {
	s.tabs = append(s.tabs, tab{content: &explorerTab{
		expl: NewExplorer(s.fsys.HomeDir(), s.fsys.WorkingDir()),
	}})
	s.NextTab()
}

func (s *Session) OpenFile(fpath string) {
	if fpath == "" {
		log.Println("open file: empty file path")
		return
	}
	data, err := s.fsys.ReadFile(fpath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Printf("reading '%s': %v\n", fpath, err)
		return
	}
	name := fpath
	if fpath[0] == '/' {
		rel, err := filepath.Rel(s.fsys.WorkingDir(), fpath)
		if err != nil {
			log.Printf("getting relative path '%s': %v\n", fpath, err)
			return
		}
		name = rel
	}
	md := &markdownTab{name: name}
	md.view.Editor.SetText(data)
	md.view.SplitRatio = 0.5
	s.tabs = append(s.tabs, tab{content: md})
	s.win.Invalidate()
}

func (s *Session) openExplorerDir(t *explorerTab, fpath string) {
	if fpath == "" {
		log.Println("open dir: empty file path")
		return
	}
	fpath = path.Clean(fpath)
	files, err := s.fsys.ReadDir(fpath)
	if err != nil {
		log.Println(err)
		return
	}
	sort.SliceStable(files, func(i, j int) bool {
		a, b := files[i], files[j]
		if a.IsDir() == b.IsDir() {
			return a.Name() < b.Name()
		}
		return a.IsDir()
	})
	t.expl.Populate(fpath, files)
	s.win.Invalidate()
}

func (s *Session) writeFile(fpath string, data []byte) {
	if err := s.fsys.WriteFile(fpath, data); err != nil {
		log.Println(err)
	}
}

func (s *Session) CloseActiveTab() {
	s.tabs = append(s.tabs[:s.activeTab], s.tabs[s.activeTab+1:]...)
	n := len(s.tabs)
	if s.activeTab > 0 && s.activeTab >= n {
		s.activeTab = n - 1
	}
	if n > 0 {
		s.tabs[s.activeTab].content.focus()
	}
}

func (s *Session) FocusActiveTab() {
	if i := s.activeTab; i >= 0 && i < len(s.tabs) {
		s.tabs[i].content.focus()
	}
}

func (s *Session) NextTab() {
	if s.activeTab < len(s.tabs)-1 {
		s.activeTab++
	}
}

func (s *Session) PrevTab() {
	if s.activeTab > 0 {
		s.activeTab--
	}
}

func (s *Session) SwapTabUp() {
	if s.activeTab == 0 {
		return
	}
	i := s.activeTab
	j := i - 1
	s.tabs[i], s.tabs[j] = s.tabs[j], s.tabs[i]
	s.activeTab--
}

func (s *Session) SwapTabDown() {
	if s.activeTab == len(s.tabs)-1 {
		return
	}
	i := s.activeTab
	j := i + 1
	s.tabs[i], s.tabs[j] = s.tabs[j], s.tabs[i]
	s.activeTab++
}

func (s *Session) SelectTab(n int) {
	if len(s.tabs) == 0 || n < 0 {
		return
	}
	if n >= len(s.tabs) {
		n = len(s.tabs) - 1
	}
	s.activeTab = n
	s.tabs[s.activeTab].content.focus()
}

type tabContent interface {
	title() string
	focus()
}

type markdownTab struct {
	name string
	view View
}

func (t *markdownTab) title() string {
	return t.name
}

func (t *markdownTab) focus() {
	t.view.Editor.Focus()
}

type explorerTab struct {
	expl Explorer
}

func (t *explorerTab) title() string {
	return "[Choose Files]"
}

func (t *explorerTab) focus() {}
