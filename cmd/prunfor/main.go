// chris 082815

// prunfor runs a command for an optionally limited amount of time.
//
//	usage: prunfor timelimit command [argument ...]
//
// timelimit is a non-negative time.Duration.  If timelimit is zero, no
// time timelimit will be applied to command's execution.
//
// Diagnostics
//
// prunfor will return with the following exit codes.
//
//	0 The command executed successfully in the time allotted.
//	1 The command coudn't be found, exited unsuccessfully, or some
//	  other error occurred when trying to run the command.
//	2 Invalid arguments.
//	3 The command timed out.
//
package main

import (
	"log"
	"os"
	"time"

	"path/filepath"

	"chrispennello.com/go/prun/cmd"
)

var myargs struct {
	// Name of this program as it's invoked.
	myname string

	// Time timelimit that the specified program can run for.
	timelimit time.Duration

	// Name of the command to run.
	command string

	// Optional arguments to pass to the program.
	args []string
}

func badargs(format string, a ...interface{}) {
	log.Printf(format, a...)
	os.Exit(2)
}

func usage() {
	badargs("usage: %s timelimit command [argument ...]\n", myargs.myname)
}

func argerr(err error) {
	badargs("%v\n", err)
}

func init() {
	log.SetFlags(0)
	myargs.myname = filepath.Base(os.Args[0])

	if len(os.Args) < 3 {
		usage()
	}

	var err error
	myargs.timelimit, err = time.ParseDuration(os.Args[1])
	if err != nil {
		argerr(err)
	}
	if myargs.timelimit < 0 {
		badargs("timelimit must be non-negative\n")
	}

	myargs.command = os.Args[2]
	myargs.args = os.Args[3:]
}

func main() {
	proc, err := cmd.NewProc(myargs.command, myargs.args)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		success, err := proc.Wait()
		if err != nil {
			log.Fatal(err)
		}
		if !success {
			os.Exit(1)
		}
		close(done)
	}()

	if myargs.timelimit > 0 {
		timeout := make(chan struct{})
		go func() {
			time.Sleep(myargs.timelimit)
			close(timeout)
		}()
		select {
		case <-done:
			// Process exited before timeout.  Thus, there's
			// no need to wait on that anymore in
			// combination with the timeout.
			break
		case <-timeout:
			log.Printf("timed out: %s\n", proc)
			if err := proc.Kill(); err != nil {
				log.Print(err)
			}
			// Don't care if this errors.
			proc.Wait()
			os.Exit(3)
		}
	}

	<-done
}
