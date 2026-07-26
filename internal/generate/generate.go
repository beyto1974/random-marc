// Package generate drives genmarc.Record into one of gomarc's four writer
// formats, writing to an arbitrary io.Writer so callers (CLI, wasm) don't
// need to go through a file.
package generate

import (
	"fmt"
	"io"
	"time"

	marc "github.com/beyto1974/gomarc"
	"github.com/brianvoe/gofakeit/v7"
	"random-marc/internal/genmarc"
)

var validFormats = map[string]bool{
	"mrc":  true,
	"json": true,
	"xml":  true,
	"text": true,
}

// Records generates count random bibliographic records and writes them to w
// in the given format. seed == 0 seeds from the current time.
func Records(count int, format string, seed int64, w io.Writer) error {
	if count <= 0 {
		return fmt.Errorf("count must be > 0, got %d", count)
	}
	if !validFormats[format] {
		return fmt.Errorf("unknown format %q (valid: mrc, json, xml, text)", format)
	}

	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	faker := gofakeit.New(uint64(seed))

	write, closeWriter, err := newWriter(format, w)
	if err != nil {
		return err
	}

	for i := 1; i <= count; i++ {
		record, err := genmarc.Record(faker, i)
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
func newWriter(format string, sink io.Writer) (write func(*marc.Record) error, closeWriter func() error, err error) {
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
