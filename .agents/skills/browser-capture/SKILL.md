---
name: browser-capture
description: Capture step-by-step browser screenshots with Playwright and insert them into a UNSAReport report.typ as a guided "paso a paso" (step-by-step) section with figures and captions. Use when a report asks for screenshots of a web UI, a browser walkthrough, or an objective-driven procedure.
---

# Browser Capture for Reports

When a lab report asks for browser screenshots (a guided walkthrough, a web UI demo, "captura de pantalla paso a paso"), drive a headless browser with Playwright, save one screenshot per step, and wire them into `report.typ` as figures.

Do NOT use `unsarep capture` for this — that command captures **terminal** output via Freeze. This skill is for **browser** screenshots.

## Prerequisites

The Nix dev shell already provides `nodejs` and `pnpm`. Install the Chromium browser once (Playwright's own download):

```bash
pnpm dlx playwright@latest install chromium
```

If `pnpm dlx` is unavailable, use `npx playwright@latest install chromium`.

## Workflow

1. **Clarify the objective and the steps.** Ask the user for the URL, the sequence of actions (navigate, fill, click, wait), and what each screenshot must show. For a "paso a paso", define one screenshot per step.
2. **Write a throwaway script** in a temp location (e.g. `/tmp/opencode/<report>-capture.js`), never inside the report project, so it is not committed. Use the template below.
3. **Save screenshots into the project's `images/` directory** with stable, ordered names (`paso-1.png`, `paso-2.png`, …). Create the directory if needed.
4. **Reference the images in `report.typ`** as figures inside the relevant `#lab-section(...)` block (see below).
5. Delete the temp script after generating the screenshots.

## Script template

```js
const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch();
  const page = await browser.newPage({ viewport: { width: 1280, height: 720 } });

  const steps = [
    {
      url: 'https://example.com/login',
      actions: [
        { type: 'fill',  selector: '#user', value: 'admin' },
        { type: 'fill',  selector: '#pass', value: 'secret' },
        { type: 'click', selector: 'button[type=submit]' },
        { type: 'wait',  ms: 1500 },
      ],
      shot: 'images/paso-1.png',
      caption: 'Paso 1: Inicio de sesión',
    },
    {
      actions: [
        { type: 'click', selector: 'a[href="/dashboard"]' },
        { type: 'wait',  ms: 1000 },
      ],
      shot: 'images/paso-2.png',
      caption: 'Paso 2: Panel principal',
    },
  ];

  for (const step of steps) {
    if (step.url) await page.goto(step.url, { waitUntil: 'networkidle' });
    for (const a of step.actions || []) {
      if (a.type === 'fill')      await page.fill(a.selector, a.value);
      else if (a.type === 'click') await page.click(a.selector);
      else if (a.type === 'type')  await page.type(a.selector, a.value);
      else if (a.type === 'press') await page.press(a.selector, a.key);
      else if (a.type === 'wait')  await page.waitForTimeout(a.ms);
    }
    await page.screenshot({ path: step.shot, fullPage: false });
  }

  await browser.close();
})();
```

Run it from the project root (so `images/` resolves there):

```bash
node /tmp/opencode/<report>-capture.js
```

If Playwright is not resolvable from `/tmp`, install it locally in the temp dir first (`pnpm --dir /tmp/opencode add playwright`), or run `pnpm dlx` with the script's directory.

## Inserting screenshots into the report

Inside the matching `#lab-section(title: "RESULTADOS Y PRUEBAS")[...]` (or the section that fits), add one figure per step. The project compiles with `--root .`, so image paths are relative to the project root:

```typst
#figure(
  image("images/paso-1.png", width: 80%),
  caption: [Paso 1: Inicio de sesión],
)
#v(0.5em)
```

Guidelines:
- Keep one figure per step with a short `caption:` (the template sets `figure(supplement: [Figura])`).
- Order figures to match the written steps; the text should narrate each step before its figure.
- Use full-page screenshots (`fullPage: true`) only for long pages; clip the viewport otherwise to keep images legible.
- Do not commit the temp script or browser artifacts to the repo.
