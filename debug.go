package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	_debugOut     io.Writer = os.Stderr
	_debugEnabled           = os.Getenv("DEBUG_EXTRACT") != ""
)

func debug(msg string, args ...interface{}) {
	if !_debugEnabled {
		return
	}

	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Fprintf(_debugOut, msg, args...)
}
