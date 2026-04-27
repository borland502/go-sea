package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/borland502/go-sea/internal/version"

	"github.com/spf13/cobra"
)

func TestVersionCommand(t *testing.T) {
	originalArgs := os.Args
	t.Cleanup(func() {
		os.Args = originalArgs
	})
	os.Args = []string{"/tmp/bin/go-sea", "version"}

	output := &bytes.Buffer{}
	command := &cobra.Command{Run: versionCmd.Run}
	command.SetOut(output)

	command.Run(command, nil)

	result := output.String()
	if result == "" {
		t.Fatal("version command produced no output")
	}

	if !bytes.Contains(output.Bytes(), []byte("go-sea version")) {
		t.Errorf("output missing binary version prefix: %s", result)
	}

	if !bytes.Contains(output.Bytes(), []byte(version.Version)) {
		t.Errorf("output missing version %q: %s", version.Version, result)
	}

	if version.Commit != "unknown" && !bytes.Contains(output.Bytes(), []byte(version.Commit)) {
		t.Errorf("output missing commit %q: %s", version.Commit, result)
	}

	if version.Date != "unknown" && !bytes.Contains(output.Bytes(), []byte(version.Date)) {
		t.Errorf("output missing build date %q: %s", version.Date, result)
	}
}

func TestFormatVersionOutputOmitsUnknownMetadata(t *testing.T) {
	result := formatVersionOutput("go-sea", "v0.1.2", "unknown", "unknown")

	if result != "go-sea version v0.1.2" {
		t.Fatalf("formatVersionOutput() = %q, want %q", result, "go-sea version v0.1.2")
	}
}

func TestFormatVersionOutputIncludesKnownMetadata(t *testing.T) {
	result := formatVersionOutput("go-sea", "v0.1.2", "abc123", "2026-04-27T17:21:40Z")
	expected := "go-sea version v0.1.2 (commit abc123, built 2026-04-27T17:21:40Z)"

	if result != expected {
		t.Fatalf("formatVersionOutput() = %q, want %q", result, expected)
	}
}

func TestDisplayBinaryNameUsesArgv0Base(t *testing.T) {
	originalArgs := os.Args
	t.Cleanup(func() {
		os.Args = originalArgs
	})

	os.Args = []string{"/tmp/bin/go-sea", "version"}

	if got := displayBinaryName(); got != "go-sea" {
		t.Fatalf("displayBinaryName() = %q, want %q", got, "go-sea")
	}
}
