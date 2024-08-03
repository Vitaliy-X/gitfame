package exceptions

import (
	"log"
)

func Exception(err error, errorString string) {
	if err != nil {
		log.Fatalf("[ERROR] %s:\n%v", errorString, err)
	}
}
