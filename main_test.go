package main

import (
	"os"
	"path/filepath"
	"testing"
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

func TestRun_WritesToFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.mrc")
	if err := run(3, "mrc", path, 42); err != nil {
		t.Fatalf("run: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if len(data) == 0 {
		t.Error("output file is empty")
	}
}
