package main

import (
	"log"
	"os"

	yadsh "github.com/ep4sh/yadsharez/pkg/yadsh"
)

func main() {
	args := os.Args[:]
	if len(args) < 2 || len(args) > 3 {
		yadsh.Help()
	}

	filePath := os.Args[1]
	log.Printf("Filepath: %s\n", filePath)

	token, err := yadsh.GetOAuthToken()
	if err != nil {
		log.Fatalf("Export env OAUTH_TOKEN\n")
	}

	yadsh.Upload(filePath, token)
	yadsh.Download(filePath, token)
}
