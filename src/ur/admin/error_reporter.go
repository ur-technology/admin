package main

import (
	"fmt"
	"os"
)

type errorReporter interface {
	fatal(desc string, args ...interface{})
}

type osErrorReporter struct{}

func newOsErrorReporter() errorReporter {
	return &osErrorReporter{}
}

func (r *osErrorReporter) fatal(desc string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[EXITING] "+desc+"\n", args...)
	os.Exit(1)
}
