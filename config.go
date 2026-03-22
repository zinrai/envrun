package main

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Profile struct {
	Vars string `yaml:"vars"`
	File string `yaml:"file"`
}

type Config struct {
	Profiles map[string]Profile `yaml:"profiles"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config %s: %w", path, err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config %s: %w", path, err)
	}

	return &config, nil
}

func resolveProfile(profile Profile) ([][2]string, error) {
	if profile.Vars != "" && profile.File != "" {
		return nil, fmt.Errorf("profile must specify either 'vars' or 'file', not both")
	}

	if profile.Vars == "" && profile.File == "" {
		return nil, fmt.Errorf("profile must specify 'vars' or 'file'")
	}

	if profile.File != "" {
		return loadEnvFile(profile.File)
	}

	return parseEnvLines(profile.Vars)
}

func resolvePairs(envrunFile, profileName, cfgPath string) ([][2]string, error) {
	if envrunFile != "" {
		return loadEnvFile(envrunFile)
	}

	config, err := loadConfig(cfgPath)
	if err != nil {
		return nil, err
	}

	profile, ok := config.Profiles[profileName]
	if !ok {
		return nil, fmt.Errorf("profile '%s' not found in %s", profileName, cfgPath)
	}

	return resolveProfile(profile)
}
