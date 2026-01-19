package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rawnly/super-bumper/internal"
	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "bump [version] [bump-type]",
	Short: "Increment semantic version numbers",
	Long: `A simple CLI to increment semantic version numbers.

Usage:
  super-bumper [version] [bump-type]    Bump the given version
  super-bumper [bump-type]              Bump version from manifest or stdin
  echo "1.0.0" | super-bumper [type]    Bump version from stdin

super-bumper types:
  major    Increment major version (1.0.0 -> 2.0.0)
  minor    Increment minor version (1.0.0 -> 1.1.0)
  patch    Increment patch version (1.0.0 -> 1.0.1) [default]

Examples:
  super-bumper 1.0.0 patch              # Output: 1.0.1
  super-bumper 1.0.0 minor              # Output: 1.1.0
  super-bumper 1.0.0 major              # Output: 2.0.0
  super-bumper 1.0.0                    # Output: 1.0.1 (default: patch)
  echo "2.3.4" | super-bumper minor     # Output: 2.4.0
  super-bumper patch                    # Reads version from package.json/composer.json/Cargo.toml`,
	Version:      version,
	SilenceUsage: true,
	RunE:         run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	var versionStr string
	var bumpType internal.BumpType = internal.Patch

	switch len(args) {
	case 0:
		// No args: try stdin, then manifest
		v, err := readFromStdin()
		if err == nil && v != "" {
			versionStr = v
		} else {
			v, err := readFromManifest()
			if err != nil {
				return fmt.Errorf("no version provided and %w", err)
			}
			versionStr = v
		}

	case 1:
		// One arg: either version or bump type
		if internal.IsBumpType(args[0]) {
			bt, _ := internal.ParseBumpType(args[0])
			bumpType = bt
			// Try stdin, then manifest
			v, err := readFromStdin()
			if err == nil && v != "" {
				versionStr = v
			} else {
				v, err := readFromManifest()
				if err != nil {
					return fmt.Errorf("no version provided and %w", err)
				}
				versionStr = v
			}
		} else {
			versionStr = args[0]
		}

	case 2:
		// Two args: version and bump type
		if internal.IsBumpType(args[0]) {
			bumpType, _ = internal.ParseBumpType(args[0])
			versionStr = args[1]
		} else {
			versionStr = args[0]
			bt, err := internal.ParseBumpType(args[1])
			if err != nil {
				return err
			}
			bumpType = bt
		}

	default:
		return fmt.Errorf("too many arguments")
	}

	v, err := internal.Parse(versionStr)
	if err != nil {
		return err
	}

	bumped := v.Bump(bumpType)
	fmt.Println(bumped.String())
	return nil
}

func readFromStdin() (string, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", fmt.Errorf("no stdin")
	}

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text()), nil
	}
	return "", fmt.Errorf("empty stdin")
}

func readFromManifest() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return internal.DetectVersion(cwd)
}
