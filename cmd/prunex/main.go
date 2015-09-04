// chris 082815

// prunfor runs a command exclusively.
//
//	usage: prunex command [argument ...]
//
// Diagnostics
//
// prunfor will return with the following exit codes.
//
//	0 The command executed successfully.
//	1 The command coudn't be found, exited unsuccessfully, or some
//	  other error occurred when trying to run the command.
//	2 Invalid arguments.
//
// TODO Further documentation.
// TODO Better return codes.
package main

import (
	"fmt"
	"log"
	"os"

	"path/filepath"

	"chrispennello.com/go/prun/cmd"
	"chrispennello.com/go/util/lockfile"
)

var myargs struct {
	// Name of this program as it's invoked.
	myname string

	// Name of the command to run.
	command string

	// Optional arguments to pass to the program.
	args []string

	// Lock file names.
	globalname string
	localname  string
}

func usage() {
	cmd.BadArgs("usage: %s command [argument ...]\n", myargs.myname)
}

func init() {
	log.SetFlags(0)
	myargs.myname = filepath.Base(os.Args[0])

	if len(os.Args) < 2 {
		usage()
	}

	myargs.command = os.Args[1]
	myargs.args = os.Args[2:]

	tmp := os.TempDir()
	key := cmd.MakeKey(myargs.command, myargs.args)
	myargs.globalname = filepath.Join(tmp, fmt.Sprintf("%s_global", myargs.myname))
	myargs.localname  = filepath.Join(tmp, fmt.Sprintf("%s_local_%s", myargs.myname, key))
}

func main() {
	lc, err := lockfile.LockRm(myargs.globalname, myargs.localname)
	if err != nil {
		log.Fatal(err)
	}
	defer lc.Unlock()
	proc, err2 := cmd.NewProc(myargs.command, myargs.args)
	if err2 != nil {
		log.Fatal(err2)
	}
	success, err3 := proc.Wait()
	if err3 != nil {
		log.Fatal(err3)
	}
	if !success {
		os.Exit(1)
	}
}
