# Guía rápida — UNSAReport

Herramienta para crear informes de laboratorio (UNSA) en Typst: monta el proyecto, gestiona plantillas/componentes, captura terminal, compila y empaqueta la entrega.

## Requisitos

El CLI valida las herramientas en tiempo de ejecución (no al compilar):

| Herramienta | Para qué |
|-------------|----------|
| [Typst](https://typst.app/) | Compilar el informe y leer metadatos |
| [Freeze](https://github.com/charmbracelet/freeze) | Capturar salida de terminal |
| [ImageMagick](https://imagemagick.org/) (`magick`) | Convertir SVG → PNG |
| [`gh`](https://cli.github.com) | `unsarep share` (compartir con el grupo) |
| [`puml`](https://www.npmjs.com/package/node-plantuml) | Diagramas PlantUML (`npm install -g node-plantuml`) |

Con Nix: `nix develop` (o `direnv allow`) deja todo listo.

## Flujo recomendado

```bash
unsarep install lab              # 1. crear el proyecto (plantilla)
# 2. editar report.typ (cabecera) y sections/ con el agente
unsarep prepare                  # 3. compilar + empaquetar en submission/
```

## Comandos

### `unsarep install`
Baja la plantilla y crea `unsareport.json`.

```bash
unsarep install                # selector interactivo
unsarep install lab            # plantilla específica
unsarep install lab@^1.0.0     # rango de versión
unsarep install multi-lab      # multi-laboratorio
unsarep install lab --dest ./mis-informes
```

**Salida esperada**: lista de `Created: ...`, `Components installed: N`, y al final `Installation complete!`.

### `unsarep prepare`
Compila el PDF y comprime el código fuente en `submission/`.

```bash
unsarep prepare                # single-lab
unsarep prepare l1             # multi-lab desde la raíz
cd l1 && unsarep prepare       # multi-lab desde dentro del lab
unsarep prepare --configure    # re-configurar nombres de archivo
```

**Salida esperada**:

```
Compiling typst report...
Archiving src to submission/Código Fuente_01.zip...

Report: submission/Informe_01.pdf
Code:   submission/Código Fuente_01.zip
```

El nombre se genera desde la plantilla `fileTemplate` usando las variables exportadas por `report.typ` (p. ej. `{output_type}_{lab_number}`).

### `unsarep capture`
Captura salida de terminal a PNG.

```bash
unsarep capture images/salida.png "ls -la" "cat README.md"
unsarep capture --cwd ./src out.png "python" "print('hi')" "k:enter" "c:d"
```

Instrucciones: texto = escribir + Enter; `w:2s` espera; `r:` texto crudo; `c:c` Ctrl+C; `k:enter/tab/esc` teclas; `x:node` comando interactivo. El log ANSI se guarda en `capture_logs/`.

### `unsarep component`
```bash
unsarep component list            # listar componentes disponibles
unsarep component add code-block  # instalar
unsarep component remove code-block
unsarep component update          # actualizar todos
```

### `unsarep update`
Sincroniza los archivos de la plantilla con el repositorio.

```bash
unsarep update            # interactivo (diff por archivo)
unsarep update --force    # aplicar todo sin preguntar
unsarep update --rollback # restaurar el último backup
```

### `unsarep watch`
Recompila el informe al guardar cambios (preview rápido).

```bash
unsarep watch              # single-lab
unsarep watch l1           # multi-lab
```

### `unsarep share`
Crea un **repo desechable en GitHub** desde la carpeta actual y lo sube, para que un compañero lo clone y siga editando con su agente.

```bash
unsarep share                         # privado, nombre desde la carpeta
unsarep share mi-informe-lab-01       # nombre explícito
unsarep share --public --description "Informe grupo 5"
```

Si la carpeta no es un repo git, `share` hace `git init` y un commit inicial automáticamente. Requiere `gh` autenticado (`gh auth login`).

**Salida esperada**: `Repository ready: https://github.com/<user>/<name>`.

## Trabajo en equipo

1. Uno crea el repo: `unsarep share`.
2. Añade colaboradores en GitHub (`Settings → Collaborators`).
3. Cada integrante clona y trabaja con su agente; al terminar `git commit` + `git push`.
4. Los demás hacen `git pull` y ven los cambios. Para controlar el flujo, usa `main` protegida + Pull Requests.

## Agentes (skills)

Al editar con un agente, estas skills guían el trabajo en `.agents/skills/`:

- `report-header` — pedir cabecera (curso, práctica, docente, integrantes) antes de editar `report.typ`.
- `report-content` — decidir tabla vs lista vs callout, y estilo de tabla propio (cabecera gris claro, celdas blancas).
- `report-code` — incrustar código fuente real con el componente `code-block` + snippets nombrados.
- `report-terminal` — capturas de terminal (`unsarep capture`).
- `report-references` — `bibliography.bib` + citas IEEE.
- `report-review` — checklist previo a la entrega.
- `browser-capture` — capturas de navegador paso a paso (Playwright).
- `puml-diagram` — diagramas PlantUML con `puml`.

## Configuración

`unsareport.json` controla el CLI. Campos útiles: `capture.columns`, `prepare.output.fileTemplate`, `prepare.output.submissionDir`. `unsareport.lock` se mantiene solo (no editar).

## Problemas comunes

| Síntoma | Causa / solución |
|---------|------------------|
| `unsareport.json not found` | Ejecuta dentro del proyecto o en una carpeta hija. |
| `missing required external tool on PATH` | Instala typst/freeze/magick o usa `nix develop`. |
| `source directory not found, skipping zip` | No existe `src/`; crea la carpeta con tu código. |
| `gh repo create` falla | `gh auth login` y verifica permisos. |
