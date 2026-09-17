---
name: report-terminal
description: Capture terminal output as images with `unsarep capture` (Freeze) and embed them into a UNSAReport report.typ. Use when a report needs to show command-line sessions, program output, or a terminal walkthrough (distinct from browser screenshots).
---

# Terminal Output in Reports

Use the `unsarep capture` command to render terminal sessions to PNG and embed them in the report. This is the terminal counterpart of the `browser-capture` skill (browser screenshots) and the `puml-diagram` skill (diagrams).

## Command

```bash
unsarep capture <output.png> <instruction...>
```

The first argument is the output PNG path; the rest are instructions executed in a virtual terminal.

## Instruction syntax

| Prefix | Behavior | Example |
|--------|----------|---------|
| *(none)* | Type text and press Enter | `"ls -la"` |
| `w:` | Wait/sleep | `"w:2s"`, `"w:500ms"` |
| `r:` | Write raw text without Enter | `"r:some text"` |
| `c:` | Send Ctrl+key | `"c:c"` (Ctrl+C) |
| `k:` | Send a control key | `"k:enter"`, `"k:tab"`, `"k:esc"` |
| `x:` | Type + Enter but don't wait for the command to finish (interactive programs) | `"x:node"` |

## Example

```bash
unsarep capture images/salida.png "x:python3" "r:print('hello')" "k:enter" "c:d"
```

With a custom working directory and delays:

```bash
unsarep capture --cwd ./src images/prueba.png "make" "w:1s" "./prog"
```

## Embedding in the report

Save captures into `images/` and reference them as figures (image paths are project-root relative):

```typst
#figure(
  image("images/salida.png", width: 80%),
  caption: [Salida de la ejecución],
)
```

## Notes

- A raw ANSI log of each capture is saved automatically to `capture_logs/` (gitignored).
- Terminal width, prompt, and colors come from the `capture` block in `unsareport.json`; pass `--freeze-flags "..."` for Freeze themes and `--save-freeze-flags` to persist them.
- Keep captures to a single logical command sequence each; use several captures rather than one huge scroll.
