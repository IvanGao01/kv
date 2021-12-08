package root

import (
	"github.com/ivangao01/kv/cmd/cli/get"
	"github.com/ivangao01/kv/cmd/cli/set"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "tool",
	Long: "tool set <key value> | get <key>",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableDescriptions: true,
		DisableNoDescFlag:   true,
	},
}

func Execute() error {
	rootCmd.AddCommand(set.Cmd)
	rootCmd.AddCommand(get.Cmd)
	return rootCmd.Execute()
}
