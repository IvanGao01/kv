package root

import (
	"github.com/ivangao01/kv/cmd/cli/get"
	"github.com/ivangao01/kv/cmd/cli/set"
	"github.com/spf13/cobra"
)

var Run = func(cmd *cobra.Command, args []string) {}

func init() {
	rootCmd.AddCommand(set.SetCmd)
	rootCmd.AddCommand(get.GetCmd)
}

var rootCmd = &cobra.Command{
	Use:  "tool",
	Long: "tool set <key value> | get <key>",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableDescriptions: true,
		DisableNoDescFlag:   true,
	},
}

func AddCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}

func Execute() error {
	return rootCmd.Execute()
}
