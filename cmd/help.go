package cmd

import "github.com/spf13/cobra"

// runHelp will run the cobra Help() for given cobra command.
// ...
// pls help
func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
