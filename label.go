package textdraw

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
)

// FillTextOnRRect paints str in foreground on a rounded fill sized to the text,
// placed within rect at fractional alignment (ax, ay).
func FillTextOnRRect(gtx layout.Context, shaper *text.Shaper, style TextStyle, rect image.Rectangle, ax, ay float32, radius int, fill, foreground color.Color, str string) layout.Dimensions {
	gtx.Constraints.Max.X = min(gtx.Constraints.Max.X, rect.Dx())

	textSize := MeasureText(gtx, shaper, style, str)

	textWidth := (textSize.X*3 + textSize.Y*2) / 3
	leading := rect.Min.X + int(ax*float32(rect.Dx())) - textWidth/2
	trailing := leading + textWidth
	if trailing > rect.Max.X {
		leading, trailing = rect.Max.X-textWidth, rect.Max.X
	}
	if leading < rect.Min.X {
		leading, trailing = rect.Min.X, rect.Min.X+textWidth
	}
	top := rect.Min.Y + int(ay*float32(rect.Dy()-textSize.Y))
	bottom := top + textSize.Y
	textrect := image.Rect(leading, top, trailing, bottom)

	c := color.NRGBAModel.Convert(fill).(color.NRGBA)
	shape := clip.UniformRRect(textrect, radius).Op(gtx.Ops)
	paint.FillShape(gtx.Ops, c, shape)

	FillText(gtx, shaper, style, textrect, 0.5, ay, foreground, str)

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

// TextOnRRect is [FillTextOnRRect] as a layout.Widget over the incoming
// constraints.
func TextOnRRect(shaper *text.Shaper, style TextStyle, ax, ay float32, radius int, fill, foreground color.Color, str string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return FillTextOnRRect(gtx, shaper, style, image.Rectangle{Max: gtx.Constraints.Max}, ax, ay, radius, fill, foreground, str)
	}
}
