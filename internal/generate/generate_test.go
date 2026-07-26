package generate

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	marc "github.com/beyto1974/gomarc"
)

func TestRecords_InvalidCount(t *testing.T) {
	var buf bytes.Buffer
	if err := Records(0, "mrc", 1, &buf); err == nil {
		t.Fatal("expected error for count=0, got nil")
	}
}

func TestRecords_InvalidFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := Records(1, "bogus", 1, &buf); err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestRecords_WritesEachFormat(t *testing.T) {
	for _, format := range []string{"mrc", "json", "xml", "text"} {
		t.Run(format, func(t *testing.T) {
			var buf bytes.Buffer
			if err := Records(3, format, 42, &buf); err != nil {
				t.Fatalf("Records: %v", err)
			}

			n, err := countRecords(format, buf.Bytes())
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
	var buf bytes.Buffer
	if _, _, err := newWriter("bogus", &buf); err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func countRecords(format string, data []byte) (int, error) {
	switch format {
	case "mrc":
		r := marc.NewReader(bytes.NewReader(data))
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
		recs, err := marc.ParseJSON(data)
		if err != nil {
			return 0, err
		}
		return len(recs), nil
	case "xml":
		recs, err := marc.ParseXML(bytes.NewReader(data))
		if err != nil {
			return 0, err
		}
		return len(recs), nil
	case "text":
		if len(data) == 0 {
			return 0, nil
		}
		return strings.Count(string(data), "=001"), nil
	default:
		return 0, errors.New("unknown format")
	}
}
