// chris 082815

// prunfor runs a command for an optionally limited amount of time.
//
//	usage: prunfor limit command [argument ...]
//
// limit is a non-negative time.Duration.  If limit is zero, no time
// limit will be applied to command's execution.
package main

import (
	"os"
	"log"
	"path"
	"time"
)

var myargs struct {
	// Name of this program as it's invoked.
	myname string

	// Time limit that the specified program can run for.
	limit time.Duration

	// Name of the program to run.
	name string

	// Optional arguments to pass to the program.
	args []string
}

func usage() {
	log.Printf("usage: %s limit command [argument ...]\n", myargs.myname)
	os.Exit(2)
}

func init() {
	log.SetFlags(0)
	myargs.myname = path.Base(os.Args[0])

	if len(os.Args) < 3 {
		usage()
	}

	var err error
	myargs.limit, err = time.ParseDuration(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	if myargs.limit < 0 {
		log.Fatal("limit must be non-negative")
	}

	myargs.name = os.Args[2]
	myargs.args = os.Args[3:]
}

func main() {
	log.Printf("%#v\n", myargs)
}
