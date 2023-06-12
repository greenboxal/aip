package mdedit

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type ViewStyle struct {
	Theme      *material.Theme
	EditorFont font.Font
	Palette    Palette
	View       *View
}

type Palette struct {
	Fg         color.NRGBA
	Bg         color.NRGBA
	LineNumber color.NRGBA
	Heading    color.NRGBA
	ListMarker color.NRGBA
	BlockQuote color.NRGBA
	CodeBlock  color.NRGBA
}

func (vs ViewStyle) Layout(gtx C) D {
	return vs.View.Layout(gtx, vs.Theme, vs.EditorFont, vs.Palette)
}

type ViewMode uint8

const (
	ViewModeSplit ViewMode = iota
	ViewModeSingle
)

type SingleViewWidget uint8

const (
	SingleViewEditor SingleViewWidget = iota
	SingleViewDocument
)

type View struct {
	Editor   Editor
	document Document

	Mode         ViewMode
	doSplitView  widget.Clickable
	doSingleView widget.Clickable

	SplitRatio   float32 // portion of total space to allow the first widget
	dividerPos   int     // position (in pixels) of the divider
	dividerDrag  gesture.Drag
	dividerClick gesture.Click

	SingleWidget SingleViewWidget
	showEditor   widget.Clickable
	showDocument widget.Clickable
}

func (vw *View) Layout(gtx C, th *material.Theme, edFnt font.Font, pal Palette) D {
	vw.update(gtx)
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			if vw.Mode == ViewModeSingle {
				return vw.laySingleView(gtx, th, edFnt, pal)
			}
			return vw.laySplitView(gtx, th, edFnt, pal)
		}),
		layout.Rigid(func(gtx C) D {
			return vw.layToolbar(gtx, th)
		}),
	)
}

func (vw *View) update(gtx C) {
	// Update view mode if view mode buttons are clicked.
	if vw.doSplitView.Clicked() {
		vw.Mode = ViewModeSplit
		op.InvalidateOp{}.Add(gtx.Ops)
	}
	if vw.doSingleView.Clicked() {
		vw.Mode = ViewModeSingle
		op.InvalidateOp{}.Add(gtx.Ops)
	}
	// Update which single view widget to show.
	if vw.showEditor.Clicked() {
		vw.SingleWidget = SingleViewEditor
		vw.Editor.Focus()
		op.InvalidateOp{}.Add(gtx.Ops)
	}
	if vw.showDocument.Clicked() {
		vw.SingleWidget = SingleViewDocument
		op.InvalidateOp{}.Add(gtx.Ops)
	}
}

func (vw *View) laySplitView(gtx C, th *material.Theme, edFnt font.Font, pal Palette) D {
	if vw.Editor.HasChanged() {
		vw.Editor.highlight()
		_ = vw.document.Render(vw.Editor.Text(), th)
	}

	maxWidth := float32(gtx.Constraints.Max.X)
	vw.dividerPos = int(vw.SplitRatio * maxWidth)

	m := op.Record(gtx.Ops)
	dims := vw.layDivider(gtx, func(gtx C) D {
		return layout.Inset{Left: 3, Right: 3}.Layout(gtx, rule{
			width: 2,
			color: color.NRGBA{90, 90, 90, 255},
			axis:  layout.Vertical,
		}.Layout)
	})
	c := m.Stop()

	vw.SplitRatio = float32(vw.dividerPos) / maxWidth
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(vw.SplitRatio, func(gtx C) D {
			return vw.Editor.Layout(gtx, th.Shaper, edFnt, th.TextSize, pal)
		}),
		layout.Rigid(func(gtx C) D {
			c.Add(gtx.Ops)
			return dims
		}),
		layout.Flexed(1-vw.SplitRatio, func(gtx C) D {
			return vw.document.Layout(gtx, th)
		}),
	)
}

func (vw *View) layDivider(gtx C, w layout.Widget) D {
	var de *pointer.Event
	for _, e := range vw.dividerDrag.Events(gtx.Metric, gtx, gesture.Horizontal) {
		if e.Type == pointer.Drag {
			de = &e
		}
	}
	if de != nil {
		vw.dividerPos += int(de.Position.X)
	}
	for _, e := range vw.dividerClick.Events(gtx) {
		if e.Type == gesture.TypeClick && e.NumClicks == 2 {
			vw.SplitRatio = 0.5
			vw.dividerPos = int(vw.SplitRatio * float32(gtx.Constraints.Max.X))
		}
	}

	if vw.dividerPos < 0 {
		vw.dividerPos = 0
	} else if vw.dividerPos > gtx.Constraints.Max.X {
		vw.dividerPos = gtx.Constraints.Max.X
	}

	dims := w(gtx)
	rect := image.Rectangle{Max: dims.Size}
	defer clip.Rect(rect).Push(gtx.Ops).Pop()

	vw.dividerDrag.Add(gtx.Ops)
	vw.dividerClick.Add(gtx.Ops)
	pointer.CursorColResize.Add(gtx.Ops)
	return dims
}

func (vw *View) laySingleView(gtx C, th *material.Theme, edFnt font.Font, pal Palette) D {
	if vw.SingleWidget == SingleViewDocument {
		return vw.document.Layout(gtx, th)
	}
	if vw.Editor.HasChanged() {
		vw.Editor.highlight()
		_ = vw.document.Render(vw.Editor.Text(), th)
	}
	return vw.Editor.Layout(gtx, th.Shaper, edFnt, th.TextSize, pal)
}

func (vw *View) layToolbar(gtx C, th *material.Theme) D {
	modeButtons := []groupButton{
		{
			click:    &vw.doSingleView,
			icon:     iconWebAsset,
			disabled: vw.Mode == ViewModeSingle,
		},
		{
			click:    &vw.doSplitView,
			icon:     iconReader,
			disabled: vw.Mode == ViewModeSplit,
		},
	}

	m := op.Record(gtx.Ops)
	dims := layout.UniformInset(8).Layout(gtx, func(gtx C) D {
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return D{Size: image.Point{gtx.Constraints.Max.X, 1}}
			}),
			layout.Rigid(func(gtx C) D {
				if vw.Mode != ViewModeSingle {
					return D{}
				}
				return buttonGroup{
					bg:       merge(th.Bg, th.Fg, 0.08),
					fg:       th.Fg,
					shaper:   th.Shaper,
					textSize: th.TextSize,
				}.layout(gtx, []groupButton{
					{
						click:    &vw.showEditor,
						icon:     iconEdit,
						disabled: vw.SingleWidget == SingleViewEditor,
					},
					{
						click:    &vw.showDocument,
						icon:     iconVisibility,
						disabled: vw.SingleWidget == SingleViewDocument,
					},
				})
			}),
			layout.Rigid(layout.Spacer{Width: 8}.Layout),
			layout.Rigid(func(gtx C) D {
				return buttonGroup{
					bg:       merge(th.Bg, th.Fg, 0.08),
					fg:       th.Fg,
					shaper:   th.Shaper,
					textSize: th.TextSize,
				}.layout(gtx, modeButtons)
			}),
		)
	})
	call := m.Stop()

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(rule{
			width: 2,
			color: merge(th.Fg, th.Bg, 0.7),
		}.Layout),
		layout.Rigid(func(gtx C) D {
			rect := clip.Rect{Max: dims.Size}
			paint.FillShape(gtx.Ops, merge(th.Bg, th.Fg, 0.04), rect.Op())
			call.Add(gtx.Ops)
			return dims
		}),
	)
}
