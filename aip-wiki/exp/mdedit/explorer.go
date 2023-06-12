package mdedit

import (
	"image/color"
	"io/fs"
	"strings"

	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Explorer struct {
	currentDir  string
	rootDir     string
	homeBtn     widget.Clickable
	mkdirBtn    widget.Clickable
	bcrumbs     []breadcrumb
	entries     []explEntry
	bcrumbList  layout.List
	entryList   widget.List
	lastClicked int
	events      []ExplorerEvent
}

func NewExplorer(rootDir, openToDir string) Explorer {
	ex := Explorer{rootDir: rootDir}
	ex.entryList.Axis = layout.Vertical
	if openToDir == "" {
		openToDir = rootDir
	}
	ex.events = append(ex.events, DirChosenEvent{openToDir})
	return ex
}

func (ex *Explorer) Layout(gtx C, th *material.Theme) D {
	ex.update(gtx)
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return ex.layTopbar(gtx, th)
		}),
		layout.Rigid(rule{width: 1, color: th.Fg}.Layout),
		layout.Flexed(1, func(gtx C) D {
			return ex.layEntryList(gtx, th)
		}),
	)
}

func (ex *Explorer) update(gtx C) {
	if ex.homeBtn.Clicked() {
		ex.chooseDir(gtx, ex.rootDir)
	}
	for i := range ex.bcrumbs {
		if ex.bcrumbs[i].btn.Clicked() {
			ex.chooseDir(gtx, ex.bcrumbs[i].path)
		}
	}
}

func (ex *Explorer) chooseDir(gtx C, fpath string) {
	ex.events = append(ex.events, DirChosenEvent{fpath})
	op.InvalidateOp{}.Add(gtx.Ops)
}

func (ex *Explorer) layTopbar(gtx C, th *material.Theme) D {
	bcrumbBtns := []groupButton{
		{click: &ex.homeBtn, icon: iconHome},
	}
	for i, bc := range ex.bcrumbs {
		bcrumbBtns = append(bcrumbBtns, groupButton{
			click:    &ex.bcrumbs[i].btn,
			text:     bc.name,
			disabled: ex.currentDir == bc.path,
		})
	}
	return layout.UniformInset(5).Layout(gtx, func(gtx C) D {
		return ex.bcrumbList.Layout(gtx, 1, func(gtx C, i int) D {
			return buttonGroup{
				bg:       th.Bg,
				fg:       th.Fg,
				shaper:   th.Shaper,
				textSize: th.TextSize * 0.9,
			}.layout(gtx, bcrumbBtns)
		})
	})
}

func (ex *Explorer) layEntryList(gtx C, th *material.Theme) D {
	if len(ex.entries) == 0 {
		return layout.Center.Layout(gtx, material.Body1(th, "No files here!").Layout)
	}
	headers := func(gtx C) D {
		sp := layout.Spacer{Width: 5}.Layout
		return layout.UniformInset(8).Layout(gtx, func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Flexed(1, headerLbl(th, "Name")),
				layout.Rigid(sp),
				layout.Flexed(0.3, headerLbl(th, "Last Modified")),
				layout.Rigid(sp),
			)
		})
	}
	entryList := func(gtx C) D {
		return material.List(th, &ex.entryList).Layout(gtx, len(ex.entries), func(gtx C, i int) D {
			if i >= len(ex.entries) {
				return D{}
			}
			en := &ex.entries[i]
			if sel := en.selectionMade(); sel != nil {
				ex.makeSelection(sel, i)
				op.InvalidateOp{}.Add(gtx.Ops)
			}
			if en.wasActivated() {
				fpath := ex.currentDir + "/" + en.info.Name()
				if en.info.IsDir() {
					ex.chooseDir(gtx, fpath)
				} else {
					ex.events = append(ex.events, FilesChosenEvent{[]string{fpath}})
				}
				op.InvalidateOp{}.Add(gtx.Ops)
			}
			return en.layout(gtx, th)
		})
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(headers),
		layout.Rigid(rule{width: 1, color: th.Fg}.Layout),
		layout.Flexed(1, entryList),
	)
}

func (ex *Explorer) makeSelection(sel *entrySelection, i int) {
	en := &ex.entries[i]
	switch {
	case sel.isCtrl:
		en.selected = !en.selected
	case sel.isShift:
		if i > ex.lastClicked {
			for z := ex.lastClicked + 1; z <= i; z++ {
				ex.entries[z].selected = !ex.entries[z].selected
			}
		}
		if i < ex.lastClicked {
			for z := ex.lastClicked - 1; z >= i; z-- {
				ex.entries[z].selected = !ex.entries[z].selected
			}
		}
	default:
		ex.deselectAll()
		en.selected = true
	}
	ex.lastClicked = i
}

func (ex *Explorer) selectAll() {
	for i := range ex.entries {
		ex.entries[i].selected = true
	}
}

func (ex *Explorer) deselectAll() {
	for i := range ex.entries {
		en := &ex.entries[i]
		en.selected = false
	}
}

func (ex *Explorer) selectedEntries() (sel []fs.FileInfo) {
	for i := range ex.entries {
		en := &ex.entries[i]
		if en.selected {
			sel = append(sel, en.info)
		}
	}
	return
}

func (ex *Explorer) Events() []ExplorerEvent {
	e := ex.events
	ex.events = nil
	return e
}

func (ex *Explorer) Populate(dir string, files []fs.FileInfo) {
	parentPath := strings.TrimPrefix(strings.TrimPrefix(dir, ex.rootDir), "/")
	var bcrumbs []breadcrumb
	if parentPath != "" {
		parentDirs := strings.Split(parentPath, "/")
		for i, dname := range parentDirs {
			bcrumbs = append(bcrumbs, breadcrumb{
				name: dname,
				path: strings.Join(append([]string{ex.rootDir}, parentDirs[:i+1]...), "/"),
			})
		}
	}

	entries := make([]explEntry, len(files))
	for i, info := range files {
		entries[i] = newExplorerEntry(info)
	}

	ex.currentDir = dir
	ex.bcrumbs = bcrumbs
	ex.entries = entries
}

type breadcrumb struct {
	btn  widget.Clickable
	name string
	path string
}

type explEntry struct {
	info      fs.FileInfo
	lastmod   string
	icon      *widget.Icon
	click     widget.Clickable
	selected  bool
	activated bool
	selection *entrySelection
}

func newExplorerEntry(info fs.FileInfo) explEntry {
	icon := iconRegFile
	if info.IsDir() {
		icon = iconDirectory
	}
	return explEntry{
		info:    info,
		lastmod: info.ModTime().Format("2 Jan 2006 15:04"),
		icon:    icon,
	}
}

func (en *explEntry) layout(gtx C, th *material.Theme) D {
	en.update(gtx)

	macro := op.Record(gtx.Ops)
	dims := en.draw(gtx, th)
	call := macro.Stop()

	var bg color.NRGBA
	switch {
	case en.selected:
		bg = lighten(th.ContrastFg, 0.05)
	case en.click.Hovered():
		bg = darken(th.ContrastFg, 0.1)
	}

	paint.FillShape(gtx.Ops, bg, clip.Rect{Max: dims.Size}.Op())
	return en.click.Layout(gtx, func(gtx C) D {
		call.Add(gtx.Ops)
		return dims
	})
}

func (en *explEntry) update(gtx C) {
	for _, c := range en.click.Clicks() {
		switch c.NumClicks {
		case 2:
			en.activated = true
			en.selection = &entrySelection{}
		case 1:
			en.selection = &entrySelection{
				isCtrl:  c.Modifiers.Contain(key.ModCtrl),
				isShift: c.Modifiers.Contain(key.ModShift),
			}
		}
		op.InvalidateOp{}.Add(gtx.Ops)
	}
}

func (en *explEntry) draw(gtx C, th *material.Theme) D {
	icon := func(gtx C) D {
		return en.icon.Layout(gtx, th.Fg)
	}
	name := material.Body2(th, en.info.Name())
	lastmod := material.Body2(th, en.lastmod)
	if !en.selected && !en.click.Hovered() {
		lastmod.Color.A /= 4
	}
	return layout.Inset{Top: 2, Bottom: 2}.Layout(gtx, func(gtx C) D {
		sp := layout.Spacer{Width: 10}.Layout
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(sp),
			layout.Rigid(icon),
			layout.Rigid(sp),
			layout.Flexed(1, name.Layout),
			layout.Rigid(sp),
			layout.Flexed(0.3, lastmod.Layout),
			layout.Rigid(sp),
		)
	})
}

func (en *explEntry) selectionMade() *entrySelection {
	s := en.selection
	en.selection = nil
	return s
}

func (en *explEntry) wasActivated() bool {
	a := en.activated
	en.activated = false
	return a
}

type entrySelection struct {
	isCtrl  bool
	isShift bool
}

func headerLbl(th *material.Theme, txt string) layout.Widget {
	l := material.Body2(th, txt)
	l.Color.A /= 3
	l.Font.Weight = font.Bold
	return l.Layout
}

type ExplorerEvent interface{}

type FilesChosenEvent struct {
	Paths []string
}

type DirChosenEvent struct {
	Path string
}
