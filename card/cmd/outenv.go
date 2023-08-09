package cmd

import "github.com/spf13/cobra"

var outEnvCmd = &cobra.Command{
	Use:   "outenv",
	Short: "Output all enviroment variables to std",
	Run: func(cmd *cobra.Command, args []string) {
		newService().OutEnv()
	},
}
