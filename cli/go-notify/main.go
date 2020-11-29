package main

import (
	"github.com/Harry-027/go-notify/cli/go-notify/cmd"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load("./../../.env")
	err  = cmd.RootCmd.Execute() // Execute the root command
	if err != nil {
		log.Println("An error occurred: ", err)
		panic(err)
	}
}
