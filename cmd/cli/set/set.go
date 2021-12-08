package set

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"strings"
)

var Cmd = &cobra.Command{
	Short: "set key=value",
	Use:   "set key=value",
	Run:   setRun,
}

func setRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("Please enter the correct parameters. ")
		return
	}
	if len(strings.Split(args[0], "=")) != 2 {
		fmt.Println("Please enter the correct parameters. ")
	}
	request, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3700/set", bytes.NewReader([]byte(args[0])))
	if err != nil {
		fmt.Println(err)
	}
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	if response.StatusCode == http.StatusBadRequest {
		fmt.Println("Write failed! ")
	}
	if response.StatusCode == http.StatusOK {
		fmt.Println("Write successfully! ")
	}

}
