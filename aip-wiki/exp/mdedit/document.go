package mdedit

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/richtext"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type Document struct {
	renderer  *docRenderer
	textState richtext.InteractiveText
	elements  []spanGroup
	elemList  widget.List
}

func (d *Document) Render(data []byte, th *material.Theme) error {
	if d.renderer == nil {
		d.renderer = newDocRenderer()
	}
	elements, err := d.renderer.Render(th, data)
	if err != nil {
		return err
	}
	d.elements = elements
	return nil
}

func (d *Document) Layout(gtx C, th *material.Theme) D {
	if d.elemList.Axis != layout.Vertical {
		d.elemList.Axis = layout.Vertical
	}
	return layout.Inset{Left: 15, Right: 10}.Layout(gtx, func(gtx C) D {
		return material.List(th, &d.elemList).Layout(gtx, len(d.elements), func(gtx C, i int) D {
			return d.layBlock(gtx, th, &d.elements[i])
		})
	})
}

func (d *Document) layBlock(gtx C, th *material.Theme, blk *spanGroup) D {
	return layout.Inset{Bottom: 24}.Layout(gtx, func(gtx C) D {
		switch blk.mdata.(type) {
		case isHr:
			size := image.Point{gtx.Constraints.Max.X, 1}
			rect := clip.Rect{Max: size}.Op()
			paint.FillShape(gtx.Ops, color.NRGBA{120, 120, 120, 255}, rect)
			return D{Size: size}
		case isFencedCodeBlock:
			m := op.Record(gtx.Ops)
			dims := layout.UniformInset(15).Layout(gtx, func(gtx C) D {
				return richtext.Text(&d.textState, th.Shaper, blk.items...).Layout(gtx)
			})
			call := m.Stop()
			rect := clip.Rect{Max: dims.Size}.Op()
			paint.FillShape(gtx.Ops, color.NRGBA{190, 190, 190, 10}, rect)
			call.Add(gtx.Ops)
			return dims
		case isBlockquote:
			m := op.Record(gtx.Ops)
			gtx.Constraints.Max.X -= 28
			dims := richtext.Text(&d.textState, th.Shaper, blk.items...).Layout(gtx)
			call := m.Stop()
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: 5, Right: 15}.Layout(gtx, func(gtx C) D {
						size := image.Point{4, dims.Size.Y}
						rect := clip.Rect{Max: size}.Op()
						paint.FillShape(gtx.Ops, color.NRGBA{120, 120, 120, 255}, rect)
						return D{Size: size}
					})
				}),
				layout.Rigid(func(gtx C) D {
					call.Add(gtx.Ops)
					return dims
				}),
			)
		default:
			return richtext.Text(&d.textState, th.Shaper, blk.items...).Layout(gtx)
		}
	})
}

type spanGroup struct {
	mdata interface{}
	items []richtext.SpanStyle
}

type spanBuilder struct {
	theme   *material.Theme
	result  []spanGroup
	current spanGroup

	list listState
}

func (sb *spanBuilder) newSpan(l material.LabelStyle) {
	sb.current.items = append(sb.current.items, richtext.SpanStyle{})
	sb.useStyle(l)
}

func (sb *spanBuilder) currentSpan() *richtext.SpanStyle {
	if len(sb.current.items) == 0 {
		sb.newSpan(material.Body1(sb.theme, ""))
	}
	return &sb.current.items[len(sb.current.items)-1]
}

func (sb *spanBuilder) useStyle(l material.LabelStyle) {
	ss := sb.currentSpan()
	ss.Font = l.Font
	ss.Color = l.Color
	ss.Size = l.TextSize
}

func (sb *spanBuilder) commitGroup() {
	sb.result = append(sb.result, sb.current)
	sb.current = spanGroup{}
}

func (sb *spanBuilder) writeLines(source []byte, n ast.Node) {
	l := n.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		v := line.Value(source)
		v = bytes.ReplaceAll(v, []byte{'\t'}, []byte("    ")) // TODO this should only replace leading tab characters not in strings
		if i == l-1 && len(v) > 0 && v[len(v)-1] == '\n' {
			v = v[:len(v)-1]
		}
		sb.currentSpan().Content += string(v)
	}
}

func (sb *spanBuilder) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// blocks
	reg.Register(ast.KindDocument, sb.renderDocument)
	reg.Register(ast.KindHeading, sb.renderHeading)
	reg.Register(ast.KindBlockquote, sb.renderBlockquote)
	reg.Register(ast.KindCodeBlock, sb.renderCodeBlock)
	reg.Register(ast.KindFencedCodeBlock, sb.renderFencedCodeBlock)
	reg.Register(ast.KindHTMLBlock, sb.renderHTMLBlock)
	reg.Register(ast.KindList, sb.renderList)
	reg.Register(ast.KindListItem, sb.renderListItem)
	reg.Register(ast.KindParagraph, sb.renderParagraph)
	reg.Register(ast.KindTextBlock, sb.renderTextBlock)
	reg.Register(ast.KindThematicBreak, sb.renderThematicBreak)
	// inlines
	reg.Register(ast.KindAutoLink, sb.renderAutoLink)
	reg.Register(ast.KindCodeSpan, sb.renderCodeSpan)
	reg.Register(ast.KindEmphasis, sb.renderEmphasis)
	reg.Register(ast.KindImage, sb.renderImage)
	reg.Register(ast.KindLink, sb.renderLink)
	reg.Register(ast.KindRawHTML, sb.renderRawHTML)
	reg.Register(ast.KindText, sb.renderText)
	reg.Register(ast.KindString, sb.renderString)
}

func (sb *spanBuilder) renderDocument(_ util.BufWriter, _ []byte, _ ast.Node, _ bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (sb *spanBuilder) renderHeading(_ util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		h := material.H6
		n := node.(*ast.Heading)
		switch n.Level {
		case 1:
			h = material.H1
		case 2:
			h = material.H2
		case 3:
			h = material.H3
		case 4:
			h = material.H4
		case 5:
			h = material.H5
		}
		l := h(sb.theme, "")
		sb.useStyle(l)
	} else {
		sb.commitGroup()
	}
	return ast.WalkContinue, nil
}

func (sb *spanBuilder) renderBlockquote(_ util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		sb.current.mdata = isBlockquote{}
		sb.currentSpan().Color.A = 120
		sb.currentSpan().Font.Style = font.Italic
	}
	return ast.WalkContinue, nil
}

func (sb *spanBuilder) renderCodeBlock(_ util.BufWriter, _ []byte, _ ast.Node, _ bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (sb *spanBuilder) renderFencedCodeBlock(_ util.BufWriter, src []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		sb.current.mdata = isFencedCodeBlock{}
		sb.currentSpan().Font.Variant = "Mono"
		sb.writeLines(src, n) // TODO store the text content in the spanGroup metadata
	} else {
		sb.commitGroup()
	}
	return ast.WalkContinue, nil
}

func (sb *spanBuilder) renderHTMLBlock(_ util.BufWriter, _ []byte, _ ast.Node, _ bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (sb *spanBuilder) renderList(_ util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.List)
		sb.list.isOrdered = n.IsOrdered()
		sb.list.index = 0
	} else {
		sb.currentSpan().Content = strings.TrimSuffix(sb.currentSpan().Content, "\n")
		sb.commitGroup()
	}
	return ast.WalkContinue, nil
}

func (sb *spanBuilder) renderListItem(_ util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		sb.newSpan(material.Body1(sb.theme, ""))
		prefix := "  •  "
		if sb.list.isOrdered {
			sb.list.index++
			prefix = fmt.Sprintf("  %d.  ", sb.list.index)
		}
		sb.currentSpan().Content = prefix
	} else {
		sb.currentSpan().Content += "\n"
	}
	return ast.WalkContinue, nil
}

func (sb *spanBuilder) renderParagraph(_ util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		sb.commitGroup()
	}
	return ast.WalkContinue, nil
}

func (sb *spanBuilder) renderTextBlock(_ util.BufWriter, src []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (sb *spanBuilder) renderThematicBreak(_ util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		sb.current.mdata = isHr{}
	} else {
		sb.commitGroup()
	}
	return ast.WalkContinue, nil
}

func (sb *spanBuilder) renderAutoLink(_ util.BufWriter, _ []byte, _ ast.Node, _ bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (sb *spanBuilder) renderCodeSpan(_ util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		l := material.Body1(sb.theme, "")
		l.Color = color.NRGBA{162, 120, 70, 255}
		l.Font.Variant = "Mono"
		sb.newSpan(l)
	} else {
		sb.newSpan(material.Body1(sb.theme, ""))
	}
	return ast.WalkContinue, nil
}

func (sb *spanBuilder) renderEmphasis(_ util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.Emphasis)
		l := material.Body1(sb.theme, "")
		if n.Level == 2 {
			l.Font.Weight = font.Bold
			l.Color.A = 255
		} else {
			l.Font.Style = font.Italic
		}
		sb.newSpan(l)
	} else {
		sb.newSpan(material.Body1(sb.theme, ""))
	}
	return ast.WalkContinue, nil
}

func (sb *spanBuilder) renderImage(_ util.BufWriter, _ []byte, _ ast.Node, _ bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (sb *spanBuilder) renderLink(_ util.BufWriter, _ []byte, _ ast.Node, _ bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (sb *spanBuilder) renderRawHTML(_ util.BufWriter, _ []byte, _ ast.Node, _ bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (sb *spanBuilder) renderText(_ util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.Text)
		segment := n.Segment
		sb.currentSpan().Content += string(segment.Value(src))
		if n.SoftLineBreak() {
			sb.currentSpan().Content += " "
		} else if n.HardLineBreak() {
			sb.currentSpan().Content += "\n"
		}
	}
	return ast.WalkContinue, nil
}

func (sb *spanBuilder) renderString(_ util.BufWriter, _ []byte, node ast.Node, _ bool) (ast.WalkStatus, error) {
	n := node.(*ast.String)
	sb.currentSpan().Content += string(n.Value)
	return ast.WalkContinue, nil
}

func (sb *spanBuilder) Result() []spanGroup {
	res := sb.result
	sb.result = nil
	return res
}

type docRenderer struct {
	sb *spanBuilder
	md goldmark.Markdown
}

func newDocRenderer() *docRenderer {
	sb := &spanBuilder{}
	md := goldmark.New(
		goldmark.WithRenderer(
			renderer.NewRenderer(
				renderer.WithNodeRenderers(util.Prioritized(sb, 0)),
			),
		),
	)
	return &docRenderer{sb, md}
}

func (r *docRenderer) Render(th *material.Theme, src []byte) ([]spanGroup, error) {
	if r.sb.theme != th {
		r.sb.theme = th
	}
	l := material.Body1(th, "")
	r.sb.useStyle(l)
	if err := r.md.Convert(src, io.Discard); err != nil {
		return nil, err
	}
	return r.sb.Result(), nil
}

type listState struct {
	isOrdered bool
	index     int
}

type isHr struct{}

type isFencedCodeBlock struct{}

type isBlockquote struct{}
