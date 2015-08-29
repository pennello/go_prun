// chris 082915

package cmd

import (
	"log"
	"os"
)

func BadArgs(format string, a ...interface{}) {
	log.Printf(format, a...)
	os.Exit(2)
}

func ArgError(err error) {
	BadArgs("%v\n", err)
}
