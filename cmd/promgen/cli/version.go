package cli

import (
	"fmt"
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/gustavooferreira/prometheus-metrics-generator/internal/build"
)

func newVersionCmd(parentCmd *cobra.Command) *cobra.Command {
	selfCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Long:  "Show version.",
		Args: func(cmd *cobra.Command, args []string) error {
			outputErr := os.Stderr

			if len(args) != 0 {
				msg := pterm.Error.Sprintfln("Accepts 0 args, received %d", len(args))
				_, _ = fmt.Fprint(outputErr, msg)
				return ErrValidation
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			outputInfo := os.Stdout

			info := build.Version()

			_, _ = fmt.Fprintf(outputInfo, "Version: %s\nSHA: %s\n", info.Version, info.SHA)

			return nil
		},
	}

	parentCmd.AddCommand(selfCmd)
	return selfCmd
}
