// Package genmarc builds random plausible bibliographic MARC21 records.
package genmarc

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	marc "github.com/beyto1974/gomarc"
)

var firstNames = []string{
	"James", "Mary", "Robert", "Patricia", "John", "Jennifer", "Michael", "Linda",
	"William", "Elizabeth", "David", "Barbara", "Richard", "Susan", "Joseph", "Jessica",
}

var lastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis",
	"Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas",
}

var titleWords = []string{
	"Shadow", "Garden", "River", "Silent", "Distant", "Broken", "Golden", "Hidden",
	"Winter", "Whisper", "Ember", "Voyage", "Echo", "Labyrinth", "Harbor", "Threshold",
}

var subtitleWords = []string{
	"a novel", "an introduction", "essays and reflections", "a memoir",
	"stories from the edge", "a practical guide", "collected writings", "a history",
}

var publishers = []string{
	"Northbridge Press", "Willow & Stone", "Alderbrook Publishing", "Ironwood Books",
	"Cascade House", "Meridian Editions", "Fernwood Press", "Cobalt & Co.",
}

var places = []string{
	"New York", "London", "Chicago", "Boston", "San Francisco", "Toronto", "Edinburgh", "Portland",
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

func pick(rng *rand.Rand, pool []string) string {
	return pool[rng.Intn(len(pool))]
}

// pickN returns n distinct random entries from pool (n must be <= len(pool)).
func pickN(rng *rand.Rand, pool []string, n int) []string {
	idx := rng.Perm(len(pool))[:n]
	out := make([]string, n)
	for i, j := range idx {
		out[i] = pool[j]
	}
	return out
}

// isbn13 returns a random 13-digit ISBN (978 prefix) with a valid EAN-13 check digit.
func isbn13(rng *rand.Rand) string {
	digits := make([]int, 12)
	digits[0], digits[1], digits[2] = 9, 7, 8
	for i := 3; i < 12; i++ {
		digits[i] = rng.Intn(10)
	}
	sum := 0
	for i, d := range digits {
		if i%2 == 0 {
			sum += d
		} else {
			sum += d * 3
		}
	}
	check := (10 - sum%10) % 10
	var sb strings.Builder
	for _, d := range digits {
		fmt.Fprintf(&sb, "%d", d)
	}
	fmt.Fprintf(&sb, "%d", check)
	return sb.String()
}

// Record builds one random bibliographic record. seq is used for the 001 control number.
func Record(rng *rand.Rand, seq int) (*marc.Record, error) {
	record, err := marc.NewRecord()
	if err != nil {
		return nil, fmt.Errorf("new record: %w", err)
	}

	now := time.Now()
	pubYear := 1980 + rng.Intn(45)
	first, last := pick(rng, firstNames), pick(rng, lastNames)
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

	record.AddField(marc.NewDataField("020", " ", " ",
		marc.Subfield{Code: "a", Value: isbn13(rng)},
	))

	record.AddField(marc.NewDataField("100", "1", " ",
		marc.Subfield{Code: "a", Value: author},
	))

	title := fmt.Sprintf("The %s %s", pick(rng, titleWords), pick(rng, titleWords))
	subfields245 := []marc.Subfield{
		{Code: "a", Value: title + " :"},
		{Code: "b", Value: pick(rng, subtitleWords) + " /"},
		{Code: "c", Value: fmt.Sprintf("%s %s.", first, last)},
	}
	record.AddField(marc.NewDataField("245", "1", "0", subfields245...))

	record.AddField(marc.NewDataField("264", " ", "1",
		marc.Subfield{Code: "a", Value: pick(rng, places) + " :"},
		marc.Subfield{Code: "b", Value: pick(rng, publishers) + ","},
		marc.Subfield{Code: "c", Value: fmt.Sprintf("%d.", pubYear)},
	))

	pages := 100 + rng.Intn(400)
	record.AddField(marc.NewDataField("300", " ", " ",
		marc.Subfield{Code: "a", Value: fmt.Sprintf("%d pages :", pages)},
		marc.Subfield{Code: "b", Value: "illustrations ;"},
		marc.Subfield{Code: "c", Value: "24 cm"},
	))

	subjectCount := 1 + rng.Intn(3)
	for _, s := range pickN(rng, subjects, subjectCount) {
		record.AddField(marc.NewDataField("650", " ", "0",
			marc.Subfield{Code: "a", Value: s + "."},
		))
	}

	if rng.Intn(2) == 0 {
		record.AddField(marc.NewDataField("500", " ", " ",
			marc.Subfield{Code: "a", Value: pick(rng, notes)},
		))
	}

	return record, nil
}
