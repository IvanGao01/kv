package get

import (
	"fmt"
	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Short: "get key",
	Use:   "get key",
	Run:   getRun,
}

func getRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Errorf("Error: Please enter the correct format. ")
	}
	//TODO get the value by key with _tool.tmp
}
