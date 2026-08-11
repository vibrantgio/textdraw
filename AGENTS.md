# AGENTS.md — textdraw

Text measurement and drawing straight onto a `*text.Shaper`, one level
below Gio's widget layer: the `TextStyle` struct, `MeasureText`, `FillText`
and `Text`, `FillLabel` and `Label`, and the alignment constants and
locales they take.

**Layer.** Tier 0 of ADR-001's table — a leaf one level below Gio's widget
layer, needing only Gio and `golang.org/x/image`. Its root module imports
nothing else in the organization. Imported by `style`. Outside the tier
table, also by the demo module `mvu/example`, the adapter modules
`ivg/raster/gio`, `svg/driver/gio` and `traer/gio` and the workbench
applications `iconbrowser`, `mindchat` and `todos`. Both directions are
measured rather than typed — `scripts/check-layers.sh --edges` reports the
graph and `scripts/sync-agents.sh` renders these sentences from it — so
correcting them here changes nothing.

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
moved the typeface into the theme: `theme/tokens.Typography` carries a full
`TextStyle` per MD3 role. Since it landed, this module's `TextStyle` is no
longer the type anything in the design system should be styling text with —
`style`, its only in-org library consumer, is frozen by the same ADR.

Read that as a warning against new dependencies, not as a removal: ADR-003
freezes `style` and says nothing about freezing `textdraw`. Phase C came and
went without touching this module; `MeasureText`, `FillText` and `FillLabel`
still have no replacement in the design system — they draw straight onto a
`*text.Shaper`, below the widget layer — and workbench applications still call
them directly. No phase has scheduled a replacement surface, so treat this
module as staying rather than as waiting to be removed.
