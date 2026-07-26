// Package genmarc builds random plausible bibliographic MARC21 records.
package genmarc

import (
	"fmt"
	"strings"
	"time"

	marc "github.com/beyto1974/gomarc"
	"github.com/brianvoe/gofakeit/v7"
)

// subtitleWords, subjects and notes have no faker equivalent: LCSH-style
// subject headings and real catalog note phrasing need curated pools, not
// generic word/sentence generation.
var subtitleWords = []string{
	"a novel", "an introduction", "essays and reflections", "a memoir",
	"stories from the edge", "a practical guide", "collected writings", "a history",
}

var subjects = []string{
	"Fiction", "History", "Philosophy", "Science", "Biography", "Poetry",
	"Politics", "Economics", "Psychology", "Art", "Travel", "Technology",
}

var notes = []string{
	"Includes bibliographical references and index.",
	"Includes index.",
	"\"A selection of previously unpublished material\"--Preface.",
	"Translated from the French.",
	"Originally published in 2005.",
}

// pickN returns n distinct random entries from pool (n must be <= len(pool)).
func pickN(faker *gofakeit.Faker, pool []string, n int) []string {
	cp := append([]string(nil), pool...)
	faker.ShuffleStrings(cp)
	return cp[:n]
}

// Record builds one random bibliographic record. seq is used for the 001 control number.
func Record(faker *gofakeit.Faker, seq int) (*marc.Record, error) {
	record, err := marc.NewRecord()
	if err != nil {
		return nil, fmt.Errorf("new record: %w", err)
	}

	now := time.Now()
	pubYear := faker.IntRange(1980, 2024)
	first, last := faker.FirstName(), faker.LastName()
	author := fmt.Sprintf("%s, %s.", last, first)

	record.AddField(marc.NewControlField("001", fmt.Sprintf("%09d", seq)))
	record.AddField(marc.NewControlField("003", "RandomMARC"))
	record.AddField(marc.NewControlField("005", now.Format("20060102150405.0")))

	// 008 layout (books, 40 chars): 00-05 date entered, 06 date type,
	// 07-10 date1, 11-14 date2 (blank), 15-17 place, 18-34 blank,
	// 35-37 language, 38 blank, 39 cataloging source.
	field008 := fmt.Sprintf("%s%s%4d%4s%3s%17s%3s%1s%1s",
		now.Format("060102"), "s", pubYear, "", "xxu", "", "eng", "", "d")
	record.AddField(marc.NewControlField("008", field008))

	// ISBNOptions.Separator "" is treated as unset (defaults to "-"), so
	// strip the dashes ourselves to get a bare 13-digit ISBN.
	isbn := strings.ReplaceAll(faker.ProductISBN(&gofakeit.ISBNOptions{Version: "13", Separator: "-"}), "-", "")
	record.AddField(marc.NewDataField("020", " ", " ",
		marc.Subfield{Code: "a", Value: isbn},
	))

	record.AddField(marc.NewDataField("100", "1", " ",
		marc.Subfield{Code: "a", Value: author},
	))

	subfields245 := []marc.Subfield{
		{Code: "a", Value: faker.BookTitle() + " :"},
		{Code: "b", Value: faker.RandomString(subtitleWords) + " /"},
		{Code: "c", Value: fmt.Sprintf("%s %s.", first, last)},
	}
	record.AddField(marc.NewDataField("245", "1", "0", subfields245...))

	record.AddField(marc.NewDataField("264", " ", "1",
		marc.Subfield{Code: "a", Value: faker.City() + " :"},
		marc.Subfield{Code: "b", Value: faker.Company() + ","},
		marc.Subfield{Code: "c", Value: fmt.Sprintf("%d.", pubYear)},
	))

	pages := faker.IntRange(100, 499)
	record.AddField(marc.NewDataField("300", " ", " ",
		marc.Subfield{Code: "a", Value: fmt.Sprintf("%d pages :", pages)},
		marc.Subfield{Code: "b", Value: "illustrations ;"},
		marc.Subfield{Code: "c", Value: "24 cm"},
	))

	subjectCount := faker.IntRange(1, 3)
	for _, s := range pickN(faker, subjects, subjectCount) {
		record.AddField(marc.NewDataField("650", " ", "0",
			marc.Subfield{Code: "a", Value: s + "."},
		))
	}

	if faker.IntRange(0, 1) == 0 {
		record.AddField(marc.NewDataField("500", " ", " ",
			marc.Subfield{Code: "a", Value: faker.RandomString(notes)},
		))
	}

	return record, nil
}
