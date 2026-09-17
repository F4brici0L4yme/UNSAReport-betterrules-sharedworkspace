---
name: report-content
description: Organize UNSAReport lab report body content — decide between tables, lists, callouts, and code blocks — and apply a custom table style (light-gray header cells, white content cells) that is visually distinct from the template's red styling. Use when writing sections of a report.typ.
---

# Report Content Organization

When writing the body of a UNSA `lab`/`multi-lab` report, choose the structure that fits the information, and use the project's custom table style (NOT the template's default red `#C8310E` styling) for any table you generate.

## Choosing how to organize information

| Need | Use |
|------|-----|
| Structured, comparable values (columns/rows, config, results matrix, parameters) | **Table** (custom style below) |
| Ordered steps, procedures, ranked items | `#enum` (numbered list) |
| Unordered items, features, bullet points | `#list` (bullets) |
| Term → definition (glossary, symbol, parameter name → value) | Typst `#terms()` / definition list |
| Short code or command output | `#raw` block (inline `#raw` for identifiers) |
| Larger code samples | `code-block` component (`#import "components/code-block.typ"`) |
| Highlighted note/objective/result callout | `#lab-section(title: ...)` block (the template's own callout) |

Prefer a table only when there are two or more columns of comparable data. Do not force a table for prose.

## Custom table style

The template's tables (`basic-info-table`, `page-header`) use a red header fill (`#C8310E`). For any table YOU generate inside the report body, use this distinct style instead: **light-gray header cells and white content cells**, with thin gray strokes.

```typst
#let gray-header-table(header, rows, ..args) = {
  set table.cell(inset: 0.5em)
  table(
    columns: (1fr,) * header.len(),
    stroke: 0.5pt + gray,
    table.header(
      ..header.map(h => table.cell(fill: rgb("#E0E0E0"))[
        #text(weight: "bold")[#h]
      ]),
    ),
    ..rows.map(r => r.map(c => table.cell(fill: white)[#c])).flatten(),
  )
}

#gray-header-table(
  (["Parámetro"], ["Valor"]),
  (
    (["Caché"], ["L1"]),
    (["Reloj"], ["2.4 GHz"]),
  ),
)
```

Rules:
- Header fill `#E0E0E0` (light gray), header text `weight: "bold"`.
- Content cells `fill: white`.
- Stroke `0.5pt + gray` (matches the template's `#808080` border gray, not the red).
- Left-align text; use `align: left + horizon` when cells hold varied-length text.

If the project has a `functions.typ` file, prefer defining the helper there and `#import`-ing it, instead of re-declaring it in every section. Keep the red `#C8310E` color reserved for the template's own header/info blocks; do not reuse it for body tables.

## Section structure

Body sections live under `sections/`. The top-level `report.typ` already `#include`s `sections/1-resultados.typ`, `2-cuestionario.typ`, `3-conclusiones.typ`, `4-referencias.typ`. Each of those wraps its content in `#lab-section(title: ...)[ ... ]` and `#include`s per-item files (e.g. `sections/1-resultados/ejercicio-1.typ`). Follow this pattern: add new content as an `#include` of a per-item file inside the matching `#lab-section` block.
