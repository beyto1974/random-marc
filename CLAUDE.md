# random-marc

Go CLI. Generates random bibliographic MARC21 records via
[gomarc](https://github.com/beyto1974/gomarc), using
[gofakeit](https://github.com/brianvoe/gofakeit) v7 for realistic
random names/titles/publishers/cities/ISBNs.

## Layout

- `main.go` — thin CLI wrapper: flag parsing, output sink (stdout or file), calls `generate.Records`.
- `internal/generate/generate.go` — `Records(count int, format string, seed int64, w io.Writer) error`, the count/format validation + gomarc writer dispatch (`newWriter`). The one place both the CLI and the wasm build call into.
- `internal/genmarc/genmarc.go` — `Record(faker *gofakeit.Faker, seq int) (*marc.Record, error)`, the single place record shape lives.
- `cmd/wasm/main.go` (`//go:build js && wasm`) — exposes `generate.Records` to the browser as `window.generateMARC(count, format, seed)`. Thin glue only; not built or tested by host `go build`/`go test` (build-tag-excluded, silently skipped by wildcard patterns).
- `web/` — static UI (`index.html`, `app.js`) served against `cmd/wasm`'s output. `web/main.wasm` and `web/wasm_exec.js` are build outputs (gitignored, produced by `make wasm`), not source.
- `*_test.go` alongside each package.

## Commands

- Build: `go build ./...`
- Test: `go test ./...`
- Vet: `go vet ./...`
- Run: `go run . -count N -format {mrc|json|xml|text} [-out file] [-seed N]`
- Build wasm UI: `make wasm` (outputs to `web/`), or `make serve` to also start a local server on :8080.
- Wasm code can't be unit-tested with plain `go test` (needs a JS runtime for `syscall/js`); it's smoke-tested by running the compiled `.wasm` under node with `wasm_exec.js` and calling `generateMARC` directly — see the commit that added `cmd/wasm` for the exact approach if this needs re-verifying.

## Conventions

- **Single random source.** Exactly one `*gofakeit.Faker` per run, created in `main.go` from `-seed` (0 → time-based) and threaded through `genmarc.Record`. Don't reintroduce `math/rand` into `internal/genmarc` — reproducibility (`-seed`) depends on this.
- **Curated pools stay curated.** `subjects` and `notes` in `genmarc.go` are hand-written, not faker-generated: LCSH-style subject headings and real catalog note phrasing ("Includes bibliographical references and index.") have no faker equivalent, and generic word/sentence generators produce unconvincing output for them. Everything else (names, titles, publishers, cities, ISBNs, years, page counts) comes from gofakeit.
- **`gofakeit.ISBNOptions.Separator: ""` does not mean "no separator"** — it's treated as unset and falls back to the default `"-"`. `genmarc.go` generates with `Separator: "-"` and strips dashes itself to get a bare digit string.
- **Writer formats.** All four (`mrc`/`json`/`xml`/`text`) implement `Write(*marc.Record) error`. `json` and `xml` wrap output (array / `<collection>`) and need `Close()` called after the last record; `mrc` and `text` don't.
- **Commits.** One commit per logical step for multi-step work (e.g. a library swap: failing tests → implementation → follow-up tests). A single self-contained change gets one commit.
