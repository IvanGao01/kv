package get

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"strings"
)

var GetCmd = &cobra.Command{
	Short: "get key",
	Use:   "get key",
	Run:   getRun,
}

func getRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		panic("Error: Please enter the correct format. ")
	}

	if strings.Contains(args[0], "=") {
		panic("Error: Please enter the correct format. ")
	}

	request, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3700/get", bytes.NewReader([]byte(args[0])))
	if err != nil {
		panic(err)
	}
	client := http.Client{}
	if err != nil {
		panic(err)
	}
	response, err := client.Do(request)
	var result []byte = make([]byte, 10)
	response.Body.Read(result)
	fmt.Println(string(result))

}
