package genmarc

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
)

var testSeeds = []uint64{1, 2, 42}

func TestRecord_NoError(t *testing.T) {
	for _, seed := range testSeeds {
		faker := gofakeit.New(seed)
		if _, err := Record(faker, 1); err != nil {
			t.Errorf("seed %d: Record returned error: %v", seed, err)
		}
	}
}

func TestRecord_RequiredFieldsPresent(t *testing.T) {
	singleOccurrence := []string{"001", "003", "005", "008", "020", "100", "245", "264", "300"}

	for _, seed := range testSeeds {
		faker := gofakeit.New(seed)
		record, err := Record(faker, 1)
		if err != nil {
			t.Fatalf("seed %d: Record returned error: %v", seed, err)
		}

		for _, tag := range singleOccurrence {
			fields := record.GetFields(tag)
			if len(fields) != 1 {
				t.Errorf("seed %d: tag %s: got %d fields, want 1", seed, tag, len(fields))
			}
		}

		subjects := record.GetFields("650")
		if n := len(subjects); n < 1 || n > 3 {
			t.Errorf("seed %d: got %d 650 fields, want 1-3", seed, n)
		}
	}
}

func TestRecord_008Length(t *testing.T) {
	faker := gofakeit.New(1)
	record, err := Record(faker, 1)
	if err != nil {
		t.Fatalf("Record returned error: %v", err)
	}

	field := record.Get("008")
	if field == nil {
		t.Fatal("008 field missing")
	}
	if got := len([]rune(field.Value())); got != 40 {
		t.Errorf("008 length = %d, want 40", got)
	}
}

func TestRecord_ISBN13Format(t *testing.T) {
	faker := gofakeit.New(1)
	record, err := Record(faker, 1)
	if err != nil {
		t.Fatalf("Record returned error: %v", err)
	}

	isbn, ok := record.Get("020").Subfield("a")
	if !ok {
		t.Fatal("020$a missing")
	}
	if len(isbn) != 13 {
		t.Errorf("ISBN length = %d, want 13: %q", len(isbn), isbn)
	}
	for _, r := range isbn {
		if r < '0' || r > '9' {
			t.Errorf("ISBN %q contains non-digit %q", isbn, r)
			break
		}
	}
}

func TestRecord_SeqInControlNumber(t *testing.T) {
	faker := gofakeit.New(1)
	for _, seq := range []int{1, 42, 999} {
		record, err := Record(faker, seq)
		if err != nil {
			t.Fatalf("seq %d: Record returned error: %v", seq, err)
		}
		want := fmt.Sprintf("%09d", seq)
		if got := record.Get("001").Value(); got != want {
			t.Errorf("seq %d: 001 = %q, want %q", seq, got, want)
		}
	}
}
