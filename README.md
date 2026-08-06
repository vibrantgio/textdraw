# textdraw

Text measurement and drawing for [Vibrant Gio](https://github.com/vibrantgio),
a design system for native desktop applications on macOS, Windows and Linux,
written in pure Go on [Gio](https://gioui.org) — one level below the widget
layer, straight onto a `*text.Shaper`.

Gio's `widget/material.Label` is a widget: it takes the constraints it is
given, lays itself out, and returns its dimensions. That is the right shape
most of the time and the wrong one exactly when you are painting rather than
composing — a caption centred in a grid cell you computed, a row whose height
must not change when its text does, a pill of background colour sized to the
words inside it. Doing any of those through the widget layer means laying the
text out to find out how big it is, throwing that away, and laying it out again
inside the box you then built.

textdraw gives you the two halves separately. `MeasureText` shapes a string and
returns its size as an `image.Point`, drawing nothing. `FillText` paints a
string into an `image.Rectangle` you already decided on, at a fractional
alignment within it — `(0, 0.5)` for left-and-vertically-centred, `(0.5, 0.5)`
for dead centre — and returns nothing at all, because you already know where it
went. `Text` and `Label` wrap the same code as `layout.Widget`s for the cases
where you do want a widget.

The styling parameter is `TextStyle`: font, alignment, size in `unit.Sp`, line
limit, truncator and wrap policy. Where the values come from has two eras. The
frozen [style](https://github.com/vibrantgio/style) module is a table of them,
and the support repositories' example programs still draw with
`textdraw.FillText(gtx, shaper, style.H6, …)`. The workbench applications
instead derive theirs from the theme: a small per-app conversion turns a
`spectrum/tokens.TextStyle` role into a `textdraw.TextStyle` (see
`todos/theme.go`), so the typeface and sizes arrive through the theme and this
module still only ever sees the struct it defines.

## Where it sits

Tier 0 of the stack — `mvu → spectrum → prism → pulse → cadence → markdown` —
a leaf that imports only Gio and `golang.org/x/image`. The
[organization page](https://github.com/vibrantgio) has the full tier table.

Nothing inside the design system imports it: prism, pulse, cadence and markdown
draw their own text through Gio's widget layer.
[style](https://github.com/vibrantgio/style) imports it for its `TextStyle`
table, and the callers are applications — a dozen example mains under
`mvu/example`, `ivg/raster/gio/example`, `svg/driver/gio/example` and
`traer/gio`, plus three of the seven
[workbench](https://github.com/vibrantgio/workbench) applications: `todos`,
`iconbrowser` and `mindchat`.

**This module is not deprecated, and the fate of `style` does not touch it.**
ADR-003 froze `style` — F3.4 of the
[org plan](https://github.com/vibrantgio/.github) (planned) archives that
repository at v0.0.6 — but says nothing about textdraw, which stays a live
tier-0 module: no phase deprecates it, no phase deletes it, and `MeasureText`,
`FillText` and `FillLabel` have no replacement anywhere in the design system —
there is nothing else in the organization that measures a string or paints one
into a rectangle you chose. The workbench applications kept drawing through it
when F1 moved them off `style`; only the source of their `TextStyle` values
changed, from `style`'s table to the theme's Typography roles.

```sh
go get github.com/vibrantgio/textdraw
```

Every module in the organization is on gioui.org v0.10.1 and Go 1.25.1.

## Packages

One package, at the module root.

| Symbol | |
| --- | --- |
| `TextStyle` | Font, alignment, size (`unit.Sp`), `MaxLines`, `Truncator` and `WrapPolicy`. Everything else here takes one. |
| `MeasureText` | Shapes a string and returns its size as an `image.Point`. Draws nothing, allocates no ops. |
| `FillText` | Paints a string into an `image.Rectangle` at fractional alignment `(ax, ay)` within it. Returns nothing. |
| `Text` | `FillText` as a `layout.Widget` over the incoming constraints, returning `Dimensions` with a computed `Baseline`. |
| `FillLabel` | A rounded pill of `fill` sized to the text, with the string drawn on it in `onFill`. The pill tracks `ax` horizontally and is clamped inside the rectangle. |
| `Label` | `FillLabel` as a `layout.Widget` over the incoming constraints. |
| `Start`, `End`, `Middle` | Aliases of the `gioui.org/text` alignment constants, so a caller needs one text import instead of two. |
| `FontFace` | A type *alias* for `gioui.org/font.FontFace` — not a wrapper. A `[]textdraw.FontFace` drops straight into `text.WithCollection`. |
| `EN_US`, `NL`, `ZH_CN`, `Default` | Prebuilt `system.Locale` values. All three are left-to-right; `Default` is `EN_US`. |

## Usage

Measure, then place, then draw. This is `list.go` from
[workbench/todos](https://github.com/vibrantgio/workbench/tree/master/todos),
sizing a clickable row to one line of its Title style — the theme's TitleLarge
role converted to a `textdraw.TextStyle`, with the theme's cached shaper —
and vertically centring the item's text in it:

```go
h := textdraw.MeasureText(gtx, typ.Shaper, typ.Title, "W").Y
size := image.Pt(gtx.Constraints.Max.X, h+gtx.Dp(Padding))
textdraw.FillText(gtx, typ.Shaper, typ.Title, image.Rectangle{Max: size}, 0.0, 0.5, textColor, item.Text)
return layout.Dimensions{Size: size}
```

Measuring the literal `"W"` rather than `item.Text` is the trick worth
stealing: it gives every row the same height, so the list does not reflow when
a todo is renamed. `MeasureText` returns the height of one line for any
non-empty single-line string — at 16 sp with `PxPerSp: 1`, `"W"` measures
`(16, 20)`, `"Hello"` measures `(39, 20)`, and the empty string still measures
`(0, 20)`.

`ax` and `ay` are fractions of the rectangle, not pixels, which is what makes
centring in a computed cell a single call. From `iconbrowser/view.go` —
`(0.5, 0.5)` for an empty-state notice in the middle of the pane, `(0.5, 0.0)`
for a caption centred under an icon:

```go
textdraw.FillText(gtx, t.typ.Shaper, t.typ.Notice, image.Rectangle{Max: size}, 0.5, 0.5, p.Muted, notice)

captionRect := image.Rect(cell.Min.X, gtx.Dp(8)+iconPx+gtx.Dp(4), cell.Max.X, cellH)
textdraw.FillText(gtx, t.typ.Shaper, t.typ.Caption, captionRect, 0.5, 0.0, p.Text, IconTable[icon].Name)
```

When you do want a widget, `Text` is the same drawing wrapped for `layout` —
here from `traer/gio/gravity`, an overlay title and an FPS readout pinned to
opposite corners by their alignment fractions alone:

```go
layout.UniformInset(12).Layout(gtx, textdraw.Text(shaper, style.H3, 0.0, 0.0, Grey900, "Gravity Well"))
layout.UniformInset(12).Layout(gtx, textdraw.Text(shaper, style.H4, 1.0, 1.0, Grey900, fmt.Sprint(fps, "fps")))
```

The `*text.Shaper` is always the caller's. In a Vibrant Gio application that
is the theme's — `Typography.Shaper()`, built once from the theme's faces and
cached in the value — which is what the workbench snippets above pass as
`typ.Shaper`. A program without a theme builds one per window at
layer-building scope (the example programs use
[`style.FontFaces()`](https://github.com/vibrantgio/style)) and passes it
down; a self-built shaper owns glyph caches and is not safe for concurrent
use.

## For coding assistants

Read the canonical guide before writing code against this module — the module
inventory with current tags, the application skeleton, MVU and rx semantics,
typography, and the pitfalls that are not guessable:

<https://raw.githubusercontent.com/vibrantgio/.github/master/llms.txt>

[`AGENTS.md`](./AGENTS.md) in this repository has the build and test commands.

## Status

Honest about what does not work yet. Every number below was measured against
the built module.

- **`TextStyle` is superseded as a typography source, though this module is
  not.** ADR-003's `Typography` theme token — a full
  `spectrum/tokens.TextStyle` per MD3 role — is where type decisions live
  now, so *this* `TextStyle` is a drawing parameter, not a place to define a
  type system: applications convert a theme role into one at the call site
  (`todos/theme.go` is the recipe) rather than declaring tables of them.
  [style](https://github.com/vibrantgio/style), the one in-org library that
  does declare such a table, is frozen and F3.4 (planned) archives it. The
  drawing functions keep their job; nothing replaces them.
- **`FillLabel` and `Label` have no consumer anywhere in the organization.**
  Not a module, not a demo, not a workbench application. They are the only two
  functions here that paint a background as well as glyphs, and nothing has
  ever called them, so their behaviour is unexercised. The pill's width is
  `textWidth + ⅔ × lineHeight` — for `"Save"` at 16 sp, measured 37 px wide and
  drawn 50 px wide — and that padding is added *after* the text was measured
  against the rectangle, so the pill can be up to two-thirds of a line height
  wider than the rectangle it was told to stay inside. The two clamps then
  fight: the first pushes the pill left to fit the right edge, the second
  pushes it back to the left edge, and it overflows on the right.
- **`FillLabel` returns the wrong dimensions.** It returns
  `layout.Dimensions{Size: gtx.Constraints.Max}` regardless of the rectangle it
  actually drew into, so a caller cannot learn how much space the label took.
  `FillText` returns nothing at all, by design — you already know the
  rectangle. `MeasureText` is how you find out.
- **`End` and `Middle` are exported and unused; so are the locales.** Every
  `TextStyle` declared or derived in the organization —
  [style](https://github.com/vibrantgio/style)'s fourteen and the conversions
  the workbench applications build from theme roles — sets
  `Alignment: textdraw.Start`; alignment inside the rectangle is done with the
  `ax` fraction instead. Nothing anywhere references `EN_US`, `NL`, `ZH_CN` or
  `Default` — `gtx.Locale` comes from the Gio window instead. The three locales
  are also all `system.LTR`, so the type has no RTL case to exercise, and no
  RTL text has ever been drawn through this module.
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
  measurement is line-quantised: at 16 sp, one line is 20 px, two are 40, and
  there is nothing in between.
- **`Text` reports a baseline; `Label` does not.** `Text` computes
  `Baseline: size.Y - firstLineBaseline`, so it composes correctly in a
  baseline-aligned flex. `FillLabel` and `Label` return a zero `Baseline`.
- **There are no tests and no golden images.** `go test ./...` reports "no
  test files". This is the module the applications draw all their own text
  through, and none of its output is pinned.

## License

MIT — see [LICENSE](./LICENSE).
