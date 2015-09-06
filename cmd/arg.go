// chris 082915

package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"path/filepath"
)

// Args stores the name and arguments of a command-line command.
type Args struct {
	// The name of the command.
	Name string
	// Any of its arguments.
	Args []string
}

// State stores the basic command-line argument state for a prun
// command-line implementation.
type State struct {
	// The prun command-line program.
	Me Args

	// The command it's meant to invoke.
	Cmd Args
}

// BadArgs logs the message exits the process with exit status 2.
func BadArgs(message string) {
	log.Printf(message)
	os.Exit(2)
}

// usage calls BadArgs with a standard usage message, including any
// additional arguments the prun utility might take preceding the
// command.
func usage(name string, args []string) {
	m := fmt.Sprintf("usage: %s", name)
	if len(args) > 0 {
		m += " " + strings.Join(args, " ")
	}
	m += " command [argument ...]\n"
	BadArgs(m)
}

// Parse constructs a State given any additional arguments the utility
// might take preceding the command.  If the command-line invocation is
// incorrect, Parse calls displays a standard usage message and exits
// with exit status 2.
func Parse(args ...string) State {
	name := filepath.Base(os.Args[0])
	if len(os.Args) < 2+len(args) {
		usage(name, args)
	}
	return State{
		Me: Args{
			Name: name,
			Args: os.Args[1 : 1+len(args)],
		},
		Cmd: Args{
			Name: os.Args[1+len(args)],
			Args: os.Args[1+len(args)+1:],
		},
	}
}

// ArgError calls BadArgs with the given error.
func ArgError(err error) {
	BadArgs(err.Error() + "\n")
}
