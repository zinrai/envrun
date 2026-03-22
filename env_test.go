package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStripQuotes(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`"hello"`, "hello"},
		{`'hello'`, "hello"},
		{`hello`, "hello"},
		{`""`, ""},
		{`''`, ""},
		{`"`, `"`},
		{``, ""},
		{`"mismatched'`, `"mismatched'`},
		{`'mismatched"`, `'mismatched"`},
	}

	for _, tt := range tests {
		got := stripQuotes(tt.input)
		if got != tt.want {
			t.Errorf("stripQuotes(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestParseEnvLines(t *testing.T) {
	input := `
# This is a comment
KEY=value
SECRET="quoted value"
PASSWORD='also quoted'

EMPTY=
COMPOUND=with=equals
`
	pairs, err := parseEnvLines(input)
	if err != nil {
		t.Fatalf("parseEnvLines failed: %v", err)
	}

	expected := [][2]string{
		{"KEY", "value"},
		{"SECRET", "quoted value"},
		{"PASSWORD", "also quoted"},
		{"EMPTY", ""},
		{"COMPOUND", "with=equals"},
	}

	if len(pairs) != len(expected) {
		t.Fatalf("expected %d pairs, got %d", len(expected), len(pairs))
	}

	for i, exp := range expected {
		if pairs[i][0] != exp[0] || pairs[i][1] != exp[1] {
			t.Errorf("pair[%d] = %v, want %v", i, pairs[i], exp)
		}
	}
}

func TestParseEnvLinesMalformed(t *testing.T) {
	input := "GOOD=value\nmalformed_line\nALSO_GOOD=123"

	pairs, err := parseEnvLines(input)
	if err != nil {
		t.Fatalf("parseEnvLines failed: %v", err)
	}

	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs (skipping malformed), got %d", len(pairs))
	}
}

func TestApplyEnvPairs(t *testing.T) {
	// Save and restore environment
	origValue := os.Getenv("ENVRUN_TEST_VAR")
	defer os.Setenv("ENVRUN_TEST_VAR", origValue)

	pairs := [][2]string{
		{"ENVRUN_TEST_VAR", "hello"},
	}

	if err := applyEnvPairs(pairs); err != nil {
		t.Fatalf("applyEnvPairs failed: %v", err)
	}

	got := os.Getenv("ENVRUN_TEST_VAR")
	if got != "hello" {
		t.Errorf("ENVRUN_TEST_VAR = %q, want %q", got, "hello")
	}
}

func TestApplyEnvPairsExpansion(t *testing.T) {
	// Save and restore environment
	origA := os.Getenv("ENVRUN_TEST_A")
	origB := os.Getenv("ENVRUN_TEST_B")
	defer func() {
		os.Setenv("ENVRUN_TEST_A", origA)
		os.Setenv("ENVRUN_TEST_B", origB)
	}()

	pairs := [][2]string{
		{"ENVRUN_TEST_A", "/opt/tool"},
		{"ENVRUN_TEST_B", "$ENVRUN_TEST_A/bin"},
	}

	if err := applyEnvPairs(pairs); err != nil {
		t.Fatalf("applyEnvPairs failed: %v", err)
	}

	gotA := os.Getenv("ENVRUN_TEST_A")
	if gotA != "/opt/tool" {
		t.Errorf("ENVRUN_TEST_A = %q, want %q", gotA, "/opt/tool")
	}

	gotB := os.Getenv("ENVRUN_TEST_B")
	if gotB != "/opt/tool/bin" {
		t.Errorf("ENVRUN_TEST_B = %q, want %q", gotB, "/opt/tool/bin")
	}
}

func TestLoadEnvFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := "HOST=localhost\nPORT=8080\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	pairs, err := loadEnvFile(path)
	if err != nil {
		t.Fatalf("loadEnvFile failed: %v", err)
	}

	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(pairs))
	}
}

func TestLoadEnvFileNotFound(t *testing.T) {
	_, err := loadEnvFile("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for nonexistent env file")
	}
}
