// chris 082915

package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// BadArgs logs the message exits the process with exit status 2.
func BadArgs(message string) {
	log.Printf(message)
	os.Exit(2)
}

// Usage calls BadArgs with a standard usage message, including any
// additional arguments the prun utility might take preceding the
// command.
func Usage(name string, args ...string) {
	m := fmt.Sprintf("usage: %s", name)
	if len(args) > 0 {
		m += " " + strings.Join(args, " ")
	}
	m += " command [argument ...]\n"
	BadArgs(m)
}

// ArgError calls BadArgs with the given error.
func ArgError(err error) {
	BadArgs(err.Error() + "\n")
}
