package main

import (
	"log"
	"os"

	"github.com/roderickmonk/bittrex-go/internal/archiver"
)

func main() {

	err := archiver.NewArchiver()
	if err != nil {
		log.Fatal("Error returned from NewArchiver: ", err)
	}

	os.Exit(0)

}
