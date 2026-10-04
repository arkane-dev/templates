# site

Website: SvelteKit (Svelte 5), fully prerendered with adapter-static, styled with NEONDECK.

## Setup
```bash
source .venv/bin/activate          # node/npm live in the venv (uv + nodeenv)
npm run dev                        # http://localhost:5173
npm run build                      # static site in build/ (upload anywhere)
npm run check                      # svelte-check
```

## Layout
| Path | What |
|---|---|
| `src/lib/config.ts` | site name, brand, nav, hero text. Start here |
| `src/routes/+layout.svelte` | AppShell (top bar, rail, status bar) |
| `src/routes/+page.svelte` | home: hero + MoonScroll, numbered bilingual sections |
| `src/routes/about/` | long-read page on the paper surface |
| `src/routes/+error.svelte` | "signal lost" page; also written as `build/404.html` |

NEONDECK is copied in (`install-links=true`). After changing it in `sharable_assets`, rebuild it there,
then run `npm run update:neondeck` here.
