package archiver

import (
	"log"
)

var err error

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s, %s", msg, err)
	}
}
