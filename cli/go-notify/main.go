package main

import (
	"github.com/Harry-027/go-notify/cli/go-notify/cmd"
	"log"
)

func main() {
	err := cmd.RootCmd.Execute() // Execute the root command
	if err != nil {
		log.Println("An error occurred: ", err)
		panic(err)
	}
}


