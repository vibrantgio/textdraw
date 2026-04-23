package text

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"

	"golang.org/x/image/math/fixed"
)

// min returns the minimum of two integers.
func min(l, r int) int {
	if l < r {
		return l
	}
	return r
}

func MeasureText(gtx layout.Context, shaper *text.Shaper, style TextStyle, str string) image.Point {
	parameters := text.Parameters{
		Font:       style.Font,
		Alignment:  style.Alignment,
		PxPerEm:    fixed.I(gtx.Metric.Sp(style.Size)),
		MaxLines:   style.MaxLines,
		Truncator:  style.Truncator,
		WrapPolicy: style.WrapPolicy,
		MinWidth:   gtx.Constraints.Min.X,
		MaxWidth:   gtx.Constraints.Max.X,
		Locale:     gtx.Locale,
	}
	shaper.LayoutString(parameters, str)
	dx, dy := 0, 0
	for glyph, ok := shaper.NextGlyph(); ok; glyph, ok = shaper.NextGlyph() {
		if glyph.Flags&text.FlagLineBreak != 0 {
			if dx < glyph.X.Ceil()+glyph.Advance.Ceil() {
				dx = glyph.X.Ceil() + glyph.Advance.Ceil()
			}
			dy += glyph.Ascent.Ceil() + glyph.Descent.Ceil()
		}
	}
	return image.Pt(dx, dy)
}

func FillText(gtx layout.Context, shaper *text.Shaper, style TextStyle, rect image.Rectangle, ax, ay float32, fill color.Color, str string) {
	parameters := text.Parameters{
		Font:       style.Font,
		Alignment:  style.Alignment,
		PxPerEm:    fixed.I(gtx.Metric.Sp(style.Size)),
		MaxLines:   style.MaxLines,
		Truncator:  style.Truncator,
		WrapPolicy: style.WrapPolicy,
		MinWidth:   gtx.Constraints.Min.X,
		MaxWidth:   min(gtx.Constraints.Max.X, rect.Dx()),
		Locale:     gtx.Locale,
	}
	shaper.LayoutString(parameters, str)
	lines := [][]text.Glyph(nil)
	line := []text.Glyph(nil)
	dx, dy := 0, 0
	for glyph, ok := shaper.NextGlyph(); ok; glyph, ok = shaper.NextGlyph() {
		line = append(line, glyph)
		if glyph.Flags&text.FlagLineBreak != 0 {
			lines = append(lines, line)
			line = nil
			if dx < glyph.X.Ceil()+glyph.Advance.Ceil() {
				dx = glyph.X.Ceil() + glyph.Advance.Ceil()
			}
			dy += glyph.Ascent.Ceil() + glyph.Descent.Ceil()
		}
	}

	c := color.NRGBAModel.Convert(fill).(color.NRGBA)

	offset := rect.Min.Add(image.Pt(int(ax*float32(rect.Dx()-dx)), int(ay*float32(rect.Dy()-dy))))
	for _, line := range lines {
		shape := clip.Outline{Path: shaper.Shape(line)}.Op()
		glyph := line[len(line)-1]
		offset.Y += glyph.Ascent.Ceil()
		tstack := op.Offset(offset).Push(gtx.Ops)
		paint.FillShape(gtx.Ops, c, shape)
		offset.Y += glyph.Descent.Ceil()
		tstack.Pop()
	}
}

func Text(shaper *text.Shaper, style TextStyle, ax, ay float32, fill color.Color, str string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		size := gtx.Constraints.Max
		parameters := text.Parameters{
			Font:       style.Font,
			Alignment:  style.Alignment,
			PxPerEm:    fixed.I(gtx.Metric.Sp(style.Size)),
			MaxLines:   style.MaxLines,
			Truncator:  style.Truncator,
			WrapPolicy: style.WrapPolicy,
			MinWidth:   gtx.Constraints.Min.X,
			MaxWidth:   size.X,
			Locale:     gtx.Locale,
		}
		shaper.LayoutString(parameters, str)
		lines := [][]text.Glyph(nil)
		line := []text.Glyph(nil)
		dx, dy := 0, 0
		for glyph, ok := shaper.NextGlyph(); ok; glyph, ok = shaper.NextGlyph() {
			line = append(line, glyph)
			if glyph.Flags&text.FlagLineBreak != 0 {
				lines = append(lines, line)
				line = nil
				if dx < glyph.X.Ceil()+glyph.Advance.Ceil() {
					dx = glyph.X.Ceil() + glyph.Advance.Ceil()
				}
				dy += glyph.Ascent.Ceil() + glyph.Descent.Ceil()
			}
		}

		c := color.NRGBAModel.Convert(fill).(color.NRGBA)

		offset := image.Pt(int(ax*float32(size.X-dx)), int(ay*float32(size.Y-dy)))
		baseline := 0
		for _, line := range lines {
			shape := clip.Outline{Path: shaper.Shape(line)}.Op()
			glyph := line[len(line)-1]
			offset.Y += glyph.Ascent.Ceil()
			tstack := op.Offset(offset).Push(gtx.Ops)
			if baseline == 0 {
				baseline = offset.Y
			}
			paint.FillShape(gtx.Ops, c, shape)
			offset.Y += glyph.Descent.Ceil()
			tstack.Pop()
		}
		return layout.Dimensions{Size: size, Baseline: size.Y - baseline}
	}
}
