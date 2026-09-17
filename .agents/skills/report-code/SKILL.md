---
name: report-code
description: Embed source code in a UNSAReport report using the code-block component — reference real files under src/, slice named snippets with START/END-SNIPPET markers, and pick the right import path. Use when a report section needs to show code.
---

# Embedding Code in Reports

UNSA reports render source code through the `code-block` component, installed by `unsarep install` into `components/code-block.typ`. Always reference the real source file under `src/` instead of pasting code inline — the archive in `unsarep prepare` keeps them in sync.

## Component API

```typst
#let code-block(
  file: none,        // path to the source file, relative to the PROJECT ROOT (e.g. "src/main.go")
  snippet: none,     // optional named snippet to extract (see below)
  prefix: none,      // marker prefix; defaults to "//"
  lang: none,        // syntax-highlight language (e.g. "go", "python", "java", "rust", "text")
  fill: none,        // block background; defaults to a light gray
  breakable: none,   // allow page breaks inside the block (default true)
  width: none,       // default 100%
  inset: none,
  radius: none,
  spacing: none,
  clip: none,        // clip long lines (default false, wraps instead)
  text-size: none,   // default 7pt
)
```

`file:` is resolved against the project root (the CLI compiles with `--root .`), NOT the location of the importing file. A snippet of `src/main.go` in a Java project becomes `file: "src/Main.java"`.

## Import path

The import path is relative to the file doing the `#import`:

- Top-level `report.typ`: `#import "components/code-block.typ": code-block`
- Section file (`sections/1-resultados.typ`): `#import "../components/code-block.typ": code-block`
- Per-item file (`sections/1-resultados/ejercicio-1.typ`): `#import "../../components/code-block.typ": code-block`

## Named snippets

Mark a region of a source file with comment markers (the default prefix is `//`; use `prefix:` for languages like Python/`#`):

```go
func main() {
    // START-SNIPPET,setup
    db := connect()
    // END-SNIPPET
    db.Run()
}
```

Then render only that region:

```typst
#code-block(file: "src/main.go", snippet: "setup", lang: "go")
```

For Python, set `prefix: "#"` and use `# START-SNIPPET,setup` / `# END-SNIPPET`.

## Guidance

- Use `lang:` so syntax highlighting matches the source language.
- For short terminal-like output (a few lines), prefer a plain `#raw` block instead of `code-block` (see the `report-content` skill).
- Do not duplicate code by pasting it into the `.typ`; reference `file:` and let the component read the real source.
- Wrap each code block in a short introductory sentence and, when useful, a caption using `#figure`.
