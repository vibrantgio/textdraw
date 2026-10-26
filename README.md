# textdraw

Text measurement and drawing for [Vibrant Gio](https://github.com/vibrantgio),
a design system for native desktop applications on macOS, Windows and Linux,
written in pure Go on [Gio](https://gioui.org) — one level below the component
layer, straight onto a `*text.Shaper`.

Gio's `widget/material.Label` is a component: it takes the constraints it is
given, lays itself out, and returns its dimensions. That is the right shape
most of the time and the wrong one exactly when you are painting rather than
composing — a caption centred in a grid cell you computed, a row whose height
must not change when its text does, a rounded fill sized to the
text inside it. Doing any of those through the component layer means laying the
text out to find out how big it is, throwing that away, and laying it out again
inside the box you then built.

textdraw gives you the two halves separately. `MeasureText` shapes a string and
returns its size as an `image.Point`, drawing nothing. `FillText` paints a
string into an `image.Rectangle` you already decided on, at a fractional
alignment within it — `(0, 0.5)` for left-and-vertically-centred, `(0.5, 0.5)`
for dead centre — and returns nothing at all, because you already know where it
went. `Text` and `TextOnRRect` wrap the same code as `layout.Widget`s for the cases
where you do want a component.

The styling parameter is `TextStyle`: font, alignment, size in `unit.Sp`, line
limit, truncator and wrap policy. An application on the theme derives its
values from it: a small conversion turns a typography role (`theme/tokens.TextStyle`) into a
`textdraw.TextStyle`, so the typeface and sizes arrive through the theme and
this module only ever sees the struct it defines.

## Where it sits

Tier 0 of the stack — `mvu → theme → components → effects → patterns → markdown` —
a leaf that imports only Gio and `golang.org/x/image`. The
[organization page](https://github.com/vibrantgio) has the full tier table.

**This module is not deprecated.** The freeze on
[style](https://github.com/vibrantgio/style) does not apply to it.

```sh
go get github.com/vibrantgio/textdraw
```

The Gio and Go versions are the ones `go.mod` names.

## Packages

One package, at the module root.

| Symbol | |
| --- | --- |
| `TextStyle` | Font, alignment, size (`unit.Sp`), `MaxLines`, `Truncator` and `WrapPolicy`. Everything else here takes one. |
| `MeasureText` | Shapes a string and returns its size as an `image.Point`. Draws nothing, allocates no ops. |
| `FillText` | Paints a string into an `image.Rectangle` at fractional alignment `(ax, ay)` within it. Returns nothing. |
| `Text` | `FillText` as a `layout.Widget` over the incoming constraints, returning `Dimensions` with a computed `Baseline`. |
| `FillTextOnRRect` | A rounded fill in `fill` sized to the text, with the string drawn on it in `foreground`. The fill tracks `ax` horizontally and is clamped inside the rectangle. |
| `TextOnRRect` | `FillTextOnRRect` as a `layout.Widget` over the incoming constraints. |
| `Start`, `End`, `Middle` | Aliases of the `gioui.org/text` alignment constants, so a caller needs one text import instead of two. |
| `FontFace` | A type *alias* for `gioui.org/font.FontFace` — not a wrapper. A `[]textdraw.FontFace` drops straight into `text.WithCollection`. |
| `EN_US`, `NL`, `ZH_CN`, `Default` | Prebuilt `system.Locale` values. All three are left-to-right; `Default` is `EN_US`. |

## Usage

Measure, then place, then draw. This sizes a clickable row to one line of
its Title style — a theme type role such as Title3 converted to a
`textdraw.TextStyle`, with the theme's cached shaper — and vertically centres
the item's text in it:

```go
h := textdraw.MeasureText(gtx, typ.Shaper, typ.Title, "W").Y
size := image.Pt(gtx.Constraints.Max.X, h+gtx.Dp(Padding))
textdraw.FillText(gtx, typ.Shaper, typ.Title, image.Rectangle{Max: size}, 0.0, 0.5, textColor, item.Text)
return layout.Dimensions{Size: size}
```

Measuring the literal `"W"` rather than `item.Text` is the trick worth
stealing: it gives every row the same height, so the list does not reflow when
an item's text changes. `MeasureText` returns the height of one line for any
non-empty single-line string.

`ax` and `ay` are fractions of the rectangle, not pixels, which is what makes
centring in a computed cell a single call. `(0.5, 0.0)` centres a caption
under an icon:

```go
captionRect := image.Rect(cell.Min.X, gtx.Dp(8)+iconPx+gtx.Dp(4), cell.Max.X, cellH)
textdraw.FillText(gtx, t.typ.Shaper, t.typ.Caption, captionRect, 0.5, 0.0, p.Text, IconTable[icon].Name)
```

When you do want a component, `Text` is the same drawing wrapped for `layout` —
here an overlay title and an FPS readout pinned to opposite corners by their
alignment fractions alone:

```go
layout.UniformInset(12).Layout(gtx, textdraw.Text(shaper, style.H3, 0.0, 0.0, Grey900, "Gravity Well"))
layout.UniformInset(12).Layout(gtx, textdraw.Text(shaper, style.H4, 1.0, 1.0, Grey900, fmt.Sprint(fps, "fps")))
```

The `*text.Shaper` is always the caller's. In a Vibrant Gio application that
is the theme's — `Typography.Shaper()`, built once from the theme's faces and
cached in the value — which is what the first two snippets above pass as
`typ.Shaper`. A program without a theme builds one per window at
layer-building scope and passes it down; a self-built shaper owns glyph
caches and is not safe for concurrent use.

## For coding assistants

Read this repository's [AGENTS.md](./AGENTS.md) before you write code against this module.

Build and test the root module with `go build ./...` and `go test ./...`.

## Status

Honest about what does not work yet.

- **`TextStyle` is a drawing parameter, not a typography source.** Type
  decisions are in the `Typography` theme token — a full
  `theme/tokens.TextStyle` per typography role — so applications convert a
  typography role into one of these at the call site rather than declaring
  tables of them.
- **`FillTextOnRRect` can overflow its rectangle.** It and `TextOnRRect` are the only two
  functions here that paint a fill as well as glyphs. The fill's width
  is `textWidth + ⅔ × lineHeight`, and that padding is added *after* the text
  was measured against the rectangle, so the fill can be up to two-thirds of a
  line height wider than the rectangle it was told to stay inside. The two
  clamps then fight: the first pushes the fill left to fit the right edge, the second
  pushes it back to the left edge, and it overflows on the right.
- **`FillTextOnRRect` returns the wrong dimensions.** It returns
  `layout.Dimensions{Size: gtx.Constraints.Max}` regardless of the rectangle it
  actually drew into, so a caller cannot learn how much space the fill took.
  `FillText` returns nothing at all, by design — you already know the
  rectangle. `MeasureText` is how you find out.
- **The locales are not used by the drawing functions, and none is
  right-to-left.** `MeasureText`, `FillText` and `Text` shape with
  `gtx.Locale`, which comes from the Gio window, not with `EN_US`, `NL`,
  `ZH_CN` or `Default`. The three locales are all `system.LTR`, so they offer
  no RTL case to exercise.
- **Measuring and drawing shape the string twice.** `MeasureText` and
  `FillText` each call `shaper.LayoutString` and walk the glyph iterator, so
  the measure-then-draw pattern this module is built around costs two full
  shaping passes per frame per string. There is no cached-layout API. In
  practice callers measure a constant like `"W"` and draw the real text, which
  hides the cost; measuring the real string in a long list would not.
- **Height counts only lines that end in a line break.** Both functions
  accumulate `dy` from glyphs flagged `FlagLineBreak`, so the returned height is
  a whole number of shaped lines and never a partial one. That is correct for
  Gio's shaper — it flags every line, including the last — but it means the
  measurement is line-quantised: a height is a whole number of lines, with
  nothing in between.
- **`Text` reports a baseline; `TextOnRRect` does not.** `Text` computes
  `Baseline: size.Y - firstLineBaseline`, so it composes correctly in a
  baseline-aligned flex. `FillTextOnRRect` and `TextOnRRect` return a zero `Baseline`.
- **There are no tests and no golden images.** `go test ./...` reports "no
  test files", so none of its output is pinned.

## License

MIT — see [LICENSE](./LICENSE).
