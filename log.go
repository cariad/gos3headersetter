package gos3headersetter

import (
	"fmt"
	"log"
)

func info(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	log.Printf("gos3headersetter: %s\n", msg)
}
