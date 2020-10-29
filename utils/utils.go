package utils

import (
	"log"
)

// FailOnError is an error handler that fails the program if an error is passed in
func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// LogOnError is an error handler that logs the error if one is passed in
func LogOnError(err error) {
	if err != nil {
		log.Printf("Error: %s", err)
	}
}
