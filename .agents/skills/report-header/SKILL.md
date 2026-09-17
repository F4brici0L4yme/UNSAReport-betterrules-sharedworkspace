---
name: report-header
description: Collect the UNSA lab report header metadata (course, lab title/number, instructor, members, semester, dates) BEFORE writing or editing a report.typ, then fill the `unsa-report.with(...)` block. Use whenever building, scaffolding, or updating a UNSAReport Typst report.
---

# Report Header Metadata

The UNSA `lab` and `multi-lab` templates render a fixed header (page header + "INFORMACIÓN BÁSICA" table) from metadata passed to `#show: unsa-report.with(...)`. Before writing or editing any `report.typ`, collect this metadata from the user so the header is always complete.

## Workflow

1. **Ask first.** Before touching `report.typ`, ask the user for the values below. Use the `question` tool for a single compact batch (do not ask for fields that are already filled in the file).
2. **Fill the block.** Write the values into the `#show: unsa-report.with(...)` call at the top of `report.typ`.
3. **Never leave placeholders.** Do not ship a report that still has `"Título de la Práctica"`, `"Nombre del Docente"`, or the three `Apellidos... Nombres...` stub members.

## Metadata fields

These are the exact argument names of `unsa-report` (see `lib.typ` in the template):

| Field | Required? | Meaning | Notes |
|-------|-----------|---------|-------|
| `course_name` | yes | Asignatura | e.g. `"Ingeniería de Software"` |
| `lab_title` | yes | Título de la práctica | |
| `lab_number` | yes | Número de la práctica | string, zero-padded, e.g. `"01"` |
| `instructor_name` | yes | Docente (profesor) | |
| `members` | yes | Integrantes | an array of full-name strings, e.g. `("Apellido Nombre", "Apellido2 Nombre2")` |
| `year` | no | Año lectivo | defaults to current year if omitted |
| `sem_code` | no | Nro. de semestre | `"A"` or `"B"`; defaults by month (Jan–Jul → A, Aug–Dec → B) if omitted |
| `presentation_date` | no | Fecha de presentación | defaults to today (`d/m/y`) if omitted |
| `presentation_hour` | no | Hora de presentación | defaults to `"11:59:00"` |

Example block:

```typst
#show: unsa-report.with(
  course_name: "Ingeniería de Software",
  lab_title: "Título de la Práctica",
  lab_number: "01",
  instructor_name: "Nombre del Docente",
  members: (
    "Apellido1 Nombres1",
    "Apellido2 Nombres2",
  ),
  year: none,
  sem_code: "A",
  presentation_date: none,
)
```

## Important behaviors

- The metadata is exported to the CLI via `typst query <var_export>` and drives the output filename in `unsarep prepare` (variables like `{lab_number}`, `{members}`). Keep `lab_number` stable — it becomes part of the submission filename.
- `unsa-report` accepts `..custom_vars`; extra named arguments are exported as additional CLI variables but are NOT part of the header. Do not add custom fields expecting them to appear in the header.
- `members` is an array, not a string. Separate members as distinct entries.
- Multi-lab projects: each lab directory (`l1/`, `l2/`, ...) has its own `report.typ` with its own `lab_number`; collect per-session values.
