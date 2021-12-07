package set

import (
	"fmt"
	"github.com/spf13/cobra"
)

var SetCmd = &cobra.Command{
	Short: "set key value",
	Use:   "set <key value>",
	Run:   setRun,
}

func setRun(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		fmt.Println("Please enter the correct parameters. ")
		return
	}

	//TODO write key value to _tool.tmp

}
