package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	configPath := flag.String("config", "envrun.yaml", "Path to YAML configuration file")
	profileName := flag.String("profile", "", "Profile name to use")
	showVersion := flag.Bool("version", false, "Print version information and exit")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  envrun -profile <name> [-config <path>] -- <command> [args...]\n")
		fmt.Fprintf(os.Stderr, "  ENVRUN_FILE=<path> envrun -- <command> [args...]\n")
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	cmdArgs := flag.Args()
	envrunFile := os.Getenv("ENVRUN_FILE")

	if len(cmdArgs) == 0 {
		fmt.Fprintln(os.Stderr, "Error: no command specified")
		os.Exit(1)
	}

	// Mutual exclusion check
	if *profileName != "" && envrunFile != "" {
		fmt.Fprintln(os.Stderr, "Error: cannot specify both -profile and ENVRUN_FILE")
		os.Exit(1)
	}

	if *profileName == "" && envrunFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Resolve environment variable pairs
	pairs, err := resolvePairs(envrunFile, *profileName, *configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Apply environment variables with expansion
	if err := applyEnvPairs(pairs); err != nil {
		slog.Warn("failed to apply environment variable", "error", err)
	}

	// Remove ENVRUN_FILE from child process environment
	os.Unsetenv("ENVRUN_FILE")

	// Execute command
	exitCode := executeCommand(cmdArgs)
	os.Exit(exitCode)
}
