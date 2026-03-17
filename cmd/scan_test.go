package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestScanCmd(t *testing.T) {
	// Redirect stdout to a buffer
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Execute the scan command with --dry-run
	rootCmd.SetArgs([]string{"scan", "--dry-run"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	t.Log(output)

	// Check if the output contains the table headers
	expectedHeaders := []string{"NAMESPACE", "NAME", "RESOURCE TYPE", "AGE", "REASON"}
	for _, header := range expectedHeaders {
		if !strings.Contains(output, header) {
			t.Errorf("expected output to contain header %q, but it did not", header)
		}
	}
}
