package main

import (
	"fmt"
	"os"
)

func main() {
	err := readConfig()
	if err != nil {
		errorProcess(err)
	}
	err = cmdMain()
	if err != nil {
		errorProcess(err)
	}
}

func errorProcess(err error) {
	fmt.Println(err)
	os.Exit(-1)
}
