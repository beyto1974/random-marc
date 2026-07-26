# random-marc

Go CLI. Generates random bibliographic MARC21 records via
[gomarc](https://github.com/beyto1974/gomarc).

## Build

```sh
go build .
```

## Usage

```sh
random-marc -count N -format {mrc|json|xml|text} [-out file] [-seed N]
```

Flags:

- `-count` — number of records (default 1).
- `-format` — `mrc` (binary MARC21), `json` (MARC-in-JSON), `xml` (MARCXML), `text` (MARCMaker text). Default `mrc`.
- `-out` — output file path. Default: stdout.
- `-seed` — random seed. Default: time-based. Same seed → same generated content (the `005` control field still carries the real wall-clock timestamp).

## Examples

Print 5 records as MARCMaker text:

```sh
random-marc -count 5 -format text
```

Write 100 records as binary MARC21 to a file:

```sh
random-marc -count 100 -format mrc -out records.mrc
```

Write MARC-in-JSON to a file:

```sh
random-marc -count 20 -format json -out records.json
```

Write MARCXML to a file:

```sh
random-marc -count 20 -format xml -out records.xml
```

Reproducible batch (same seed → same content):

```sh
random-marc -count 10 -format mrc -seed 42 -out a.mrc
random-marc -count 10 -format mrc -seed 42 -out b.mrc
diff <(xxd a.mrc) <(xxd b.mrc)   # only 005 timestamps differ
```

## Record shape

Each record: `001`/`003`/`005`/`008` control fields, `020` (ISBN-13),
`100` (author), `245` (title), `264` (publication), `300` (physical
description), `650` (1–3 subjects), optional `500` (note).

## Web UI (WebAssembly)

Same generator, compiled to WASM, running entirely in the browser — no
server needed once loaded.

```sh
make serve
```

Opens on `http://localhost:8080`. `make wasm` alone just builds
`web/main.wasm` + copies `web/wasm_exec.js` without serving.
