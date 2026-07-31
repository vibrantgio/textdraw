# AGENTS.md — textdraw

Text measurement and drawing straight onto a `*text.Shaper`, one level
below Gio's widget layer: the `TextStyle` struct, `MeasureText`, `FillText`
and `Text`, `FillLabel` and `Label`, and the alignment constants and
locales they take.

**Layer.** Tier 0 of ADR-001's table — a leaf that imports only Gio and
`golang.org/x/image`. `style` imports it, as do example mains under mvu,
ivg, svg and traer and three of the workbench applications; nothing in the
design system does.

**Read the canonical guide before you write code against this module.** It is
the organization's one agent guide — the module inventory with current tags,
the application skeleton, the MVU loop and rx semantics, typography, and the
pitfalls that are not guessable. It lives exactly once, in `vibrantgio/.github`,
and this file links it rather than copying it:

    https://raw.githubusercontent.com/vibrantgio/.github/master/llms.txt

**Module.** `github.com/vibrantgio/textdraw`, one module at the repository
root.

**Build and test.** From the repository root:

    go build ./... && go test ./...

**Superseded for typography, still live as a drawing primitive.** ADR-003
moves the typeface into the theme: `Typography` becomes a theme token carrying
a full `spectrum/tokens.TextStyle` per MD3 role. Once that lands, this
module's `TextStyle` is no longer the type anything in the design system
should be styling text with — `style`, its only in-org library consumer, is
frozen by the same ADR.

Read that as a warning against new dependencies, not as a removal: ADR-003
freezes `style`, and says nothing about freezing `textdraw`. Nothing in
Phase C touches this module, `MeasureText`, `FillText` and `FillLabel` have no
replacement in the design system — they draw straight onto a `*text.Shaper`,
below the widget layer — and three workbench applications call them directly
today. F2.3 revisits this note once Phase C has shipped and the real
replacement surface is known.
