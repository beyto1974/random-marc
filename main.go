// Command random-marc generates random bibliographic MARC21 records.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	marc "github.com/beyto1974/gomarc"
	"random-marc/internal/genmarc"
)

var validFormats = map[string]bool{
	"mrc":  true,
	"json": true,
	"xml":  true,
	"text": true,
}

func main() {
	count := flag.Int("count", 1, "number of records to generate")
	format := flag.String("format", "mrc", "output format: mrc, json, xml, text")
	out := flag.String("out", "", "output file path (default: stdout)")
	seed := flag.Int64("seed", 0, "random seed (default: time-based)")
	flag.Parse()

	if err := run(*count, *format, *out, *seed); err != nil {
		fmt.Fprintln(os.Stderr, "random-marc:", err)
		os.Exit(1)
	}
}

func run(count int, format, out string, seed int64) error {
	if count <= 0 {
		return fmt.Errorf("count must be > 0, got %d", count)
	}
	if !validFormats[format] {
		return fmt.Errorf("unknown format %q (valid: mrc, json, xml, text)", format)
	}

	sink := os.Stdout
	if out != "" {
		f, err := os.Create(out)
		if err != nil {
			return fmt.Errorf("create output file: %w", err)
		}
		defer f.Close()
		sink = f
	}

	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(seed))

	write, closeWriter, err := newWriter(format, sink)
	if err != nil {
		return err
	}

	for i := 1; i <= count; i++ {
		record, err := genmarc.Record(rng, i)
		if err != nil {
			return fmt.Errorf("generate record %d: %w", i, err)
		}
		if err := write(record); err != nil {
			return fmt.Errorf("write record %d: %w", i, err)
		}
	}

	if closeWriter != nil {
		if err := closeWriter(); err != nil {
			return fmt.Errorf("close writer: %w", err)
		}
	}

	return nil
}

// newWriter returns a write func for a single record plus an optional close
// func for formats that need a trailing wrapper (JSON array, XML collection).
func newWriter(format string, sink *os.File) (write func(*marc.Record) error, closeWriter func() error, err error) {
	switch format {
	case "mrc":
		w := marc.NewWriter(sink)
		return w.Write, nil, nil
	case "text":
		w := marc.NewTextWriter(sink)
		return w.Write, nil, nil
	case "json":
		w, err := marc.NewJSONWriter(sink)
		if err != nil {
			return nil, nil, fmt.Errorf("new json writer: %w", err)
		}
		return w.Write, w.Close, nil
	case "xml":
		w, err := marc.NewXMLWriter(sink)
		if err != nil {
			return nil, nil, fmt.Errorf("new xml writer: %w", err)
		}
		return w.Write, w.Close, nil
	default:
		return nil, nil, fmt.Errorf("unknown format %q", format)
	}
}
