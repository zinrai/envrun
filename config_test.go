package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	content := `
profiles:
  tfenv:
    vars: |
      TFENV_ROOT="/opt/tfenv"
      PATH="$TFENV_ROOT/bin:$PATH"
  myapp:
    file: ".env.prod"
`
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	config, err := loadConfig(path)
	if err != nil {
		t.Fatalf("loadConfig failed: %v", err)
	}

	if len(config.Profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(config.Profiles))
	}

	tfenv, ok := config.Profiles["tfenv"]
	if !ok {
		t.Fatal("profile 'tfenv' not found")
	}
	if tfenv.Vars == "" {
		t.Fatal("tfenv.Vars should not be empty")
	}
	if tfenv.File != "" {
		t.Fatal("tfenv.File should be empty")
	}

	myapp, ok := config.Profiles["myapp"]
	if !ok {
		t.Fatal("profile 'myapp' not found")
	}
	if myapp.File != ".env.prod" {
		t.Fatalf("expected myapp.File '.env.prod', got '%s'", myapp.File)
	}
	if myapp.Vars != "" {
		t.Fatal("myapp.Vars should be empty")
	}
}

func TestLoadConfigFileNotFound(t *testing.T) {
	_, err := loadConfig("/nonexistent/config.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent config file")
	}
}

func TestResolveProfileVars(t *testing.T) {
	profile := Profile{
		Vars: "KEY=value\nOTHER=123",
	}
	pairs, err := resolveProfile(profile)
	if err != nil {
		t.Fatalf("resolveProfile failed: %v", err)
	}
	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(pairs))
	}
}

func TestResolveProfileFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte("DB_HOST=localhost\n"), 0644); err != nil {
		t.Fatal(err)
	}

	profile := Profile{
		File: envPath,
	}
	pairs, err := resolveProfile(profile)
	if err != nil {
		t.Fatalf("resolveProfile failed: %v", err)
	}
	if len(pairs) != 1 {
		t.Fatalf("expected 1 pair, got %d", len(pairs))
	}
	if pairs[0][0] != "DB_HOST" || pairs[0][1] != "localhost" {
		t.Fatalf("unexpected pair: %v", pairs[0])
	}
}

func TestResolveProfileBothVarsAndFile(t *testing.T) {
	profile := Profile{
		Vars: "KEY=value",
		File: ".env",
	}
	_, err := resolveProfile(profile)
	if err == nil {
		t.Fatal("expected error when both vars and file are specified")
	}
}

func TestResolveProfileEmpty(t *testing.T) {
	profile := Profile{}
	_, err := resolveProfile(profile)
	if err == nil {
		t.Fatal("expected error when neither vars nor file is specified")
	}
}
