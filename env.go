package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func stripQuotes(s string) string {
	if len(s) < 2 {
		return s
	}
	if (s[0] == '"' && s[len(s)-1] == '"') ||
		(s[0] == '\'' && s[len(s)-1] == '\'') {
		return s[1 : len(s)-1]
	}
	return s
}

func parseEnvLines(text string) ([][2]string, error) {
	var pairs [][2]string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			slog.Warn("skipping malformed line", "line", line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = stripQuotes(value)

		pairs = append(pairs, [2]string{key, value})
	}
	return pairs, nil
}

func loadEnvFile(path string) ([][2]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read env file %s: %w", path, err)
	}
	return parseEnvLines(string(data))
}

func applyEnvPairs(pairs [][2]string) error {
	for _, pair := range pairs {
		key := pair[0]
		value := os.ExpandEnv(pair[1])
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("failed to set %s: %w", key, err)
		}
	}
	return nil
}
