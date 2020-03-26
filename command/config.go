package command

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(cfgCmd)
}

var cfgCmd = &cobra.Command{
	Use:   "config",
	Short: "Set and get gh settings",
	RunE:  cfg,
	Long: `
	TODO
`,
}

func cfg(cmd *cobra.Command, args []string) error {
	// TODO decide on top level behavior. I'm thinking either just dump the config or go into an
	// interactive mode where config keys can be searched through.
	return nil
}
