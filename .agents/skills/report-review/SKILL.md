---
name: report-review
description: Final pre-submission checklist for a UNSAReport report. Run before declaring a report done — verify the header is filled, required sections exist, code/images/diagrams resolve, citations are present, and `unsarep prepare` compiles cleanly. Use after editing a report.typ.
---

# Report Review Checklist

Before a report is considered done, run through this checklist and fix anything that fails. End by actually compiling.

## 1. Header (see report-header)

- [ ] `#show: unsa-report.with(...)` has real values — no `"Título de la Práctica"`, `"Nombre del Docente"`, or stub `Apellidos... Nombres...` members.
- [ ] `course_name`, `lab_title`, `lab_number`, `instructor_name`, `members` are all set.
- [ ] `lab_number` matches the session (per-directory in multi-lab).

## 2. Structure

- [ ] All required sections present and included from `report.typ`: `1-resultados`, `2-cuestionario`, `3-conclusiones`, `4-referencias` (or the template's equivalent).
- [ ] Each section body is wrapped in its `#lab-section(title: ...)[ ... ]` block.

## 3. Content integrity

- [ ] Every `code-block` `file:` points to a real file under `src/`.
- [ ] Every `#image(...)` path exists (browser captures in `images/`, terminal captures, `puml` diagrams).
- [ ] Every `#cite(<key>)` has a matching entry in `bibliography.bib`; no unused entries remain.
- [ ] Tables you generated use the gray-header/white-cell style (see report-content), not the template's red.

## 4. Compile

Run the actual compile and fix all errors:

```bash
unsarep prepare            # single-lab, from project root
unsarep prepare l1         # multi-lab, from project root
```

- [ ] `unsarep prepare` exits successfully and writes both the PDF and the ZIP under `submission/`.
- [ ] The generated filenames look correct (e.g. `Informe_01.pdf`, `Código Fuente_01.zip`).
- [ ] No Typst errors/warnings about missing imports, images, or undefined references.

## 5. Cleanliness

- [ ] No temp scripts or browser artifacts committed (Playwright scripts, `node_modules`, `capture_logs/`).
- [ ] `unsareport.lock` is present but not hand-edited.
