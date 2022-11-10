package main

import (
	"log"
	"os"
)

func main() {
	envs, err := ReadDir(os.Args[1])
	if err != nil {
		log.Print(err)
		os.Exit(-1)
	}

	retCode := RunCmd(os.Args[2:], envs)

	os.Exit(retCode)
}
