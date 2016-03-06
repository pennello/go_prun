// chris 2016-03-06

// TODO
// prunparallel runs commands in parallel.
//
//	usage: prunparallel total concur indextemplate command [argument ...]
//
// total is the total number of commands to run.  concur is the positive
// number of maximum concurrent executions.  intdextemplate is a string
// that, if it appears in any of the given arguments, will be
// substituted with the 0-based index of the particular command being
// executed.  If indextemplate is the empty string, no substitution will
// occur.
//
// As soon as there is a non-successful termination of one of the
// commands, prunparallel will cease launching any new commands, wait
// for the currently-running commands to terminate, and return the exit
// code of the first non-successful termination.
//
// Sample Usage
//
// TODO
//
// Diagnostics
//
// prunparallel may return with the following exit codes.
//
//	  1 An unidentified error occurred when trying to run or wait on
//	    one of the commands.
//	  2 Invalid arguments.
//	127 The command could not be found.
//
// And it will print an appropriate message to standard error.
//
// In addition, prunparallel may return with the following exit code.
//
//	255 One of the commands exited unsuccessfully, but the underlying
//	    operating system does not support examining the exit status.
//
package main

import (
	"log"
	//"os"
	"strconv"

	"chrispennello.com/go/prun/cmd"
)

var state struct {
	cmd cmd.State

	total  uint64
	concur uint64

	indextemplate string
}

func init() {
	log.SetFlags(0)
	state.cmd = cmd.Parse("total", "concur", "indextemplate")

	var err error
	state.total, err = strconv.ParseUint(state.cmd.Me.Args[0], 0, 64)
	if err != nil {
		cmd.ArgError(err)
	}
	state.concur, err = strconv.ParseUint(state.cmd.Me.Args[1], 0, 64)
	if err != nil {
		cmd.ArgError(err)
	}
	if state.concur == 0 {
		cmd.BadArgs("concur must be positive")
	}
	state.indextemplate = state.cmd.Me.Args[2]
}

func worker(work chan uint64, returncodes chan int) {
	for index := range work {
		log.Println("index", index)
		returncodes <- 0
	}
}

func main() {
	work := make(chan uint64, state.concur)
	returncodes := make(chan int)

	// Determine how many workers we'll need and start 'em all up.
	var workers uint64
	if state.concur > state.total {
		workers = state.total
	} else {
		workers = state.concur
	}
	for i := uint64(0); i < workers; i++ {
		go worker(work, returncodes)
	}

	go func() {
		for i := uint64(0); i < state.total; i++ {
			work <- i
		}
		close(work)
	}()

	for r := range returncodes {
		log.Println("return", r)
	}
}
