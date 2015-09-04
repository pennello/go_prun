// chris 082915

package cmd

import (
	"log"
	"os"
)

// BadArgs logs with the given format and additional arguments and exits
// the process with exit status 2.
func BadArgs(format string, a ...interface{}) {
	log.Printf(format, a...)
	os.Exit(2)
}

// ArgError calls BadArgs with the given error.
func ArgError(err error) {
	BadArgs("%v\n", err)
}
