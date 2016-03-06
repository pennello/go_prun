// chris 2016-01-07

// prunsleep runs a command after sleeping a random amount of time.
//
//	usage: prunsleep bound command [argument ...]
//
// bound is a non-negative time.Duration.  If bound is zero, the command
// will be executed immediately.  Generally, it represents a limit on
// how long prunsleep will wait before executing the specified command.
// Precisely, prunsleep will sleep for a pseudo-random number of
// nanoseconds in [0,bound).
//
// Sample Usage
//
// Suppose that you have a nightly cron to pull in some configuration
// data.  This cron runs on many hosts, so it would be good if they
// didn't all fire at the same time--otherwise you have a "thundering
// herd".  Thus, you might wrap your cron job with prunsleep.
//
//	@daily prunsleep 1h sh pull.sh
//
// As a point of comparison, sleeping for a random amount of time when
// run as a cron job is a feature in FreeBSD's portsnap cron command.
//
// Diagnostics
//
// prunsleep may return with the following exit codes.
//
//	  1 An unidentified error occurred when trying to run or wait on
//	    the command.
//	  2 Invalid arguments.
//	127 The command could not be found.
//
// And it will print an appropriate message to standard error.
//
// In addition, prunsleep may return with the following exit code.
//
//	255 The command exited unsuccessfully, but the underlying
//	    operating system does not support examining the exit status.
//
// Otherwise, prunsleep will return with the exit code of the command.
package main

import (
	"log"
	"os"
	"time"

	"math/rand"

	"chrispennello.com/go/prun/cmd"
)

var state struct {
	cmd cmd.State

	// Maximum time that we will sleep before running the specified
	// program.
	bound time.Duration
}

func init() {
	log.SetFlags(0)
	state.cmd = cmd.Parse("bound")

	var err error
	state.bound, err = time.ParseDuration(state.cmd.Me.Args[0])
	if err != nil {
		cmd.ArgError(err)
	}
	if state.bound < 0 {
		cmd.BadArgs("bound must be non-negative")
	}
}

func main() {
	if state.bound != 0 {
		time.Sleep(time.Duration(rand.Int63n(state.bound.Nanoseconds())))
	}
	proc := cmd.NewProc(state.cmd.Cmd.Name, state.cmd.Cmd.Args)
	proc.Cmd.Stdout = os.Stdout
	proc.Cmd.Stderr = os.Stderr
	proc.StartExit()
	proc.WaitExit()
}
