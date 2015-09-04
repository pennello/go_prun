// chris 082815

// prunfor runs a command for an optionally limited amount of time.
//
//	usage: prunfor timelimit command [argument ...]
//
// timelimit is a non-negative time.Duration.  If timelimit is zero, no
// timelimit will be applied to command's execution.
//
// Diagnostics
//
// prunfor may return with the following exit codes.
//
//	  1 An unidentified error occurred when trying to run or wait on
//	    the command.
//	  2 Invalid arguments.
//	  3 Timed out.
//	127 The command could not be found.
//
// And it will print an appropriate message to standard error.
//
// In addition, prunfor may return with the following exit code.
//
//	255 The command exited unsuccessfully, but the underlying
//	    operating system does not support examining the exit status.
//
// Otherwise, prunfor will return with the exit code of the command.
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

func usage() {
	cmd.BadArgs("usage: %s timelimit command [argument ...]\n", myargs.myname)
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
		cmd.ArgError(err)
	}
	if myargs.timelimit < 0 {
		cmd.BadArgs("timelimit must be non-negative\n")
	}

	myargs.command = os.Args[2]
	myargs.args = os.Args[3:]
}

func main() {
	proc := cmd.NewProcExit(myargs.command, myargs.args)

	done := make(chan struct{})
	go func() {
		proc.WaitExit()
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
