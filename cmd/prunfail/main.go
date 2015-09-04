// chris 090415

// TODO Write me.
package main

import (
	"log"
	"os"

	"path/filepath"

	"chrispennello.com/go/prun/cmd"
)

var myargs struct {
	// Name of this program as it's invoked.
	myname string

	// Name of the command to run.
	command string

	// Optional arguments to pass to the program.
	args []string

	// Log file name.
	logname string
}

func init() {
	log.SetFlags(0)
	myargs.myname = filepath.Base(os.Args[0])

	if len(os.Args) < 2 {
		cmd.Usage(myargs.myname)
	}
}

func main() {
}
