// Command random-marc generates random bibliographic MARC21 records.
package main

import (
	"flag"
	"fmt"
	"os"

	"random-marc/internal/generate"
)

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
	sink := os.Stdout
	if out != "" {
		f, err := os.Create(out)
		if err != nil {
			return fmt.Errorf("create output file: %w", err)
		}
		defer f.Close()
		sink = f
	}

	return generate.Records(count, format, seed, sink)
}
