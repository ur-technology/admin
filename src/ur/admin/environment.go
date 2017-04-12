package main

import (
	"os"
	"strings"
)

type environmentMap map[string]string

func parseEnvironment() environmentMap {

	m := make(environmentMap)

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		m[pair[0]] = pair[1]
	}

	return m
}
