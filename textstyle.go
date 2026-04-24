package textdraw

import (
	"gioui.org/font"
	"gioui.org/text"
	"gioui.org/unit"
)

// TextStyle contains the styling parameters for a block of text.
type TextStyle struct {
	// Font describes the preferred typeface.
	Font font.Font
	// Alignment characterizes the positioning of text within the line. It does not directly
	// impact shaping, but is provided in order to allow efficient offset computation.
	Alignment text.Alignment
	// Size is the Sp size to shape the text with.
	Size unit.Sp
	// MaxLines limits the quantity of shaped lines. Zero means no limit.
	MaxLines int
	// Truncator is a string of text to insert where the shaped text was truncated, which
	// can currently ohly happen if MaxLines is nonzero and the text on the final line is
	// truncated.
	Truncator string
	// WrapPolicy configures how line breaks will be chosen when wrapping text across lines.
	WrapPolicy text.WrapPolicy
}
