package root

import (
	"github.com/ivangao01/kv/cmd/cli/get"
	"github.com/ivangao01/kv/cmd/cli/set"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:  "set key value | get key",
	Long: "set key value | get key",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableDescriptions: true,
		DisableNoDescFlag:   true,
	},
}

func Execute() error {
	Cmd.AddCommand(set.Cmd)
	Cmd.AddCommand(get.Cmd)
	return Cmd.Execute()
}
