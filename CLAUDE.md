# random-marc

Go CLI. Generates random bibliographic MARC21 records via
[gomarc](https://github.com/beyto1974/gomarc), using
[gofakeit](https://github.com/brianvoe/gofakeit) v7 for realistic
random names/titles/publishers/cities/ISBNs.

## Layout

- `main.go` — flag parsing, output sink (stdout or file), format dispatch to the right gomarc writer (`run`/`newWriter`).
- `internal/genmarc/genmarc.go` — `Record(faker *gofakeit.Faker, seq int) (*marc.Record, error)`, the single place record shape lives.
- `*_test.go` alongside each package.

## Commands

- Build: `go build ./...`
- Test: `go test ./...`
- Vet: `go vet ./...`
- Run: `go run . -count N -format {mrc|json|xml|text} [-out file] [-seed N]`

## Conventions

- **Single random source.** Exactly one `*gofakeit.Faker` per run, created in `main.go` from `-seed` (0 → time-based) and threaded through `genmarc.Record`. Don't reintroduce `math/rand` into `internal/genmarc` — reproducibility (`-seed`) depends on this.
- **Curated pools stay curated.** `subjects` and `notes` in `genmarc.go` are hand-written, not faker-generated: LCSH-style subject headings and real catalog note phrasing ("Includes bibliographical references and index.") have no faker equivalent, and generic word/sentence generators produce unconvincing output for them. Everything else (names, titles, publishers, cities, ISBNs, years, page counts) comes from gofakeit.
- **`gofakeit.ISBNOptions.Separator: ""` does not mean "no separator"** — it's treated as unset and falls back to the default `"-"`. `genmarc.go` generates with `Separator: "-"` and strips dashes itself to get a bare digit string.
- **Writer formats.** All four (`mrc`/`json`/`xml`/`text`) implement `Write(*marc.Record) error`. `json` and `xml` wrap output (array / `<collection>`) and need `Close()` called after the last record; `mrc` and `text` don't.
- **Commits.** One commit per logical step for multi-step work (e.g. a library swap: failing tests → implementation → follow-up tests). A single self-contained change gets one commit.
