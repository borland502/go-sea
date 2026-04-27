package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/borland502/go-sea/internal/version"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:         "version",
	Short:       "Show version information",
	Annotations: map[string]string{skipConfigLoadAnnotation: "true"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), formatVersionOutput(displayBinaryName(), version.Version, version.Commit, version.Date))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func displayBinaryName() string {
	if len(os.Args) == 0 || os.Args[0] == "" {
		return binaryName
	}

	name := filepath.Base(os.Args[0])
	if name == "." || name == string(filepath.Separator) || name == "" {
		return binaryName
	}

	return name
}

func formatVersionOutput(name, versionValue, commitValue, dateValue string) string {
	parts := []string{fmt.Sprintf("%s version %s", name, versionValue)}

	if commitValue != "" && commitValue != "unknown" {
		parts = append(parts, fmt.Sprintf("commit %s", commitValue))
	}

	if dateValue != "" && dateValue != "unknown" {
		parts = append(parts, fmt.Sprintf("built %s", dateValue))
	}

	if len(parts) == 1 {
		return parts[0]
	}

	return parts[0] + " (" + strings.Join(parts[1:], ", ") + ")"
}
