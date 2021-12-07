package main

import (
	"fmt"
	"github.com/ivangao01/kv/cmd/cli/root"
)

func main() {
	if err := root.Execute(); err != nil {
		fmt.Printf("Error : %+v", err)
	}

}
