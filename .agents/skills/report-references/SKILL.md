---
name: report-references
description: Manage citations and the bibliography.bib file in a UNSAReport report using IEEE style. Use when a report section needs references, in-text citations, or a completed 4-referencias section.
---

# References and Bibliography

UNSA reports use BibTeX citations rendered in IEEE style. The template ships a `bibliography.bib` and a `sections/4-referencias.typ` that renders it with `#bibliography("../bibliography.bib", style: "ieee")`.

## Citing in text

Reference an entry by its BibTeX key:

```typst
#cite(<key>)            // e.g. #cite(<knuth1997>)
#cite(<k1, k2>)         // multiple
```

Example: `Según #cite(<gamma2011>), la arquitectura ...`

## Adding an entry

Append a correctly-typed BibTeX entry to `bibliography.bib`. Use the entry type that fits the source:

- `@article` — journal/paper
- `@book` — book
- `@inproceedings` — conference paper
- `@misc` / `@online` — web page, with `url` and `urldate`

```bibtex
@online{plantuml,
  title   = {PlantUML: Open-source UML diagramming tool},
  author  = {PlantUML Team},
  url     = {https://plantuml.com/},
  urldate = {2026-01-15},
}
```

## Guidelines

- Only add entries that are actually cited. The IEEE style lists only cited references.
- Keys must be unique and stable; prefer `authorYear` (e.g. `gamma2011`).
- Verify the citation renders before finishing — the reference list appears in `4-referencias`.
- For URLs, include `urldate` (date accessed) as IEEE requires it for online sources.
- If a reference is no longer cited, remove it from `bibliography.bib` so the list stays clean.
