// chris 2016-03-06

// TODO Extend index template injection to the command itself?

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
// This is a comparatively complex prun command.
//
// As soon as there is a non-successful termination of one of the
// commands, prunparallel will cease launching any new commands, wait
// for the currently-running commands to terminate, and return the exit
// code of that first non-successful termination.
//
// Sample Usage
//
// Here is a trivial example.
//
//	$ prunparallel 6 4 {} echo {}
//	3
//	1
//	2
//	0
//	4
//	5
//
// prunparallel will execute 6 runs in total, with 4 concurrent.  The
// indextemplate is the string "{}", and all it will do is echo the
// command index.
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
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

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

// Move NewInjectedProc into cmd package?

// NewInjectedProc returns a new cmd.Proc.  It first, however, runs
// through args, replacing occurrences of indextemplate with the decimal
// string representation of the given index.  If indextemplate is the
// empty string, no replacement occurs.
func NewInjectedProc(command string, args []string, indextemplate string, index uint64) *cmd.Proc {
	if len(indextemplate) != 0 {
		args2 := make([]string, len(args))
		copy(args2, args)
		args = args2
		new := fmt.Sprintf("%d", index)
		for i, arg := range args {
			args[i] = strings.Replace(arg, indextemplate, new, -1)
		}
	}
	return cmd.NewProc(command, args)
}

func worker(work chan *cmd.Proc, returncodes chan int, done chan struct{}) {
	log.Print("worker starting")
	workloop:
	for proc := range work {
		fns := []func() *cmd.ProcError{proc.StartError, proc.WaitError}
		for _, fn := range fns {
			if pe := fn(); pe != nil {
				pe.Print()
				returncodes <- pe.Code
				continue workloop
			}
		}
	}
	done <- struct{}{}
	log.Print("worker exiting")
}

func main() {
	// XXX Special-case when state.total == 0?  Why doesn't the
	// existing machinery just work and exit without doing anything?

	work := make(chan *cmd.Proc)
	returncodes := make(chan int)
	done := make(chan struct{})

	// Determine how many workers we'll need and start 'em all up.
	var workers uint64
	if state.concur > state.total {
		workers = state.total
	} else {
		workers = state.concur
	}
	log.Print("starting workers")
	for i := uint64(0); i < workers; i++ {
		go worker(work, returncodes, done)
	}

	// Machinery to signal work scheduler to abort scheduling more
	// work if we get a failure (immediately below).
	abort := false
	mu := &sync.Mutex{}
	aborting := func() bool {
		mu.Lock()
		defer mu.Unlock()
		return abort
	}
	setAbort := func() {
		mu.Lock()
		defer mu.Unlock()
		abort = true
	}

	// Simple work scheduler: create cmd.Proc objects based off of
	// indices and feed them into the workers.  Bug out on abort.
	go func() {
		for i := uint64(0); i < state.total; i++ {
			if aborting() {
				break
			}
			proc := NewInjectedProc(state.cmd.Cmd.Name, state.cmd.Cmd.Args, state.indextemplate, i)
			proc.Cmd.Stdout = os.Stdout
			proc.Cmd.Stderr = os.Stderr
			work <- proc
		}
		close(work)
		log.Print("closed work channel")
	}()

	workersdone := uint64(0)
	// The whole program will exit with the first non-zero
	// return code, if there is one.
	returncode := 0

	mainloop:
	for {
		log.Print("mainloop")
		select {
		case r := <- returncodes:
			if returncode == 0 && r != 0 {
				returncode = r
				setAbort()
			}
		case <-done:
			workersdone += 1
			if workersdone == workers {
				// Just in case something goes wrong and
				// someone tries to write to either of
				// these, we'll panic.
				close(returncodes)
				close(done)
				break mainloop
			}
		}
	}

	os.Exit(returncode)
}
