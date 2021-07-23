package main

import (
	"fmt"
	"github.com/tianxingpan/go-ctl/cmd"
	"os"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Println("Execute err: ", err.Error())
		os.Exit(-1)
	}
}
