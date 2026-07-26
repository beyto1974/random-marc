package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	marc "github.com/beyto1974/gomarc"
)

func TestRun_InvalidCount(t *testing.T) {
	if err := run(0, "mrc", "", 1); err == nil {
		t.Fatal("expected error for count=0, got nil")
	}
}

func TestRun_InvalidFormat(t *testing.T) {
	if err := run(1, "bogus", "", 1); err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestRun_WritesEachFormat(t *testing.T) {
	tests := []struct {
		format string
	}{
		{"mrc"}, {"json"}, {"xml"}, {"text"},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "out."+tt.format)
			if err := run(3, tt.format, path, 42); err != nil {
				t.Fatalf("run: %v", err)
			}

			n, err := countRecords(tt.format, path)
			if err != nil {
				t.Fatalf("countRecords: %v", err)
			}
			if n != 3 {
				t.Errorf("got %d records, want 3", n)
			}
		})
	}
}

func TestNewWriter_UnknownFormat(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "out")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if _, _, err := newWriter("bogus", f); err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func countRecords(format, path string) (int, error) {
	switch format {
	case "mrc":
		f, err := os.Open(path)
		if err != nil {
			return 0, err
		}
		defer f.Close()
		r := marc.NewReader(f)
		n := 0
		for {
			if _, err := r.Next(); errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				return 0, err
			}
			n++
		}
		return n, nil
	case "json":
		data, err := os.ReadFile(path)
		if err != nil {
			return 0, err
		}
		recs, err := marc.ParseJSON(data)
		if err != nil {
			return 0, err
		}
		return len(recs), nil
	case "xml":
		f, err := os.Open(path)
		if err != nil {
			return 0, err
		}
		defer f.Close()
		recs, err := marc.ParseXML(f)
		if err != nil {
			return 0, err
		}
		return len(recs), nil
	case "text":
		data, err := os.ReadFile(path)
		if err != nil {
			return 0, err
		}
		if len(data) == 0 {
			return 0, nil
		}
		return strings.Count(string(data), "=001"), nil
	default:
		return 0, errors.New("unknown format")
	}
}
