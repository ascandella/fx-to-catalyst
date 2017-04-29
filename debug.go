package main

import (
	"fmt"
	"os"
	"strings"
)

var debugEnabled = os.Getenv("DEBUG_EXTRACT") != ""

func debug(msg string, args ...interface{}) {
	if !debugEnabled {
		return
	}

	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Fprintf(os.Stderr, msg, args...)
}
