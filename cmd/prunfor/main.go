// chris 082815

// prunfor runs a command for an optionally limited amount of time.
//
//	usage: prunfor timelimit command [argument ...]
//
// timelimit is a non-negative time.Duration.  If timelimit is zero, no
// time limit will be applied to command's execution.
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

	"chrispennello.com/go/prun/cmd"
)

var state struct {
	cmd cmd.State

	// Time timelimit that the specified program can run for.
	timelimit time.Duration
}

func init() {
	log.SetFlags(0)
	state.cmd = cmd.Parse("timelimit")

	var err error
	state.timelimit, err = time.ParseDuration(state.cmd.Me.Args[0])
	if err != nil {
		cmd.ArgError(err)
	}
	if state.timelimit < 0 {
		cmd.BadArgs("timelimit must be non-negative\n")
	}
}

func main() {
	proc := cmd.NewProc(state.cmd.Cmd.Name, state.cmd.Cmd.Args)
	proc.Cmd.Stdout = os.Stdout
	proc.Cmd.Stderr = os.Stderr
	proc.StartExit()

	done := make(chan struct{})
	go func() {
		proc.WaitExit()
		close(done)
	}()

	if state.timelimit > 0 {
		timeout := make(chan struct{})
		go func() {
			time.Sleep(state.timelimit)
			close(timeout)
		}()
		select {
		case <-done:
			// Process exited successfully before timeout.
			// Thus, there's no need to wait on that anymore
			// in combination with the timeout.
			break
		case <-timeout:
			log.Printf("timed out: %s\n", proc)
			if err := proc.Cmd.Process.Kill(); err != nil {
				log.Print(err)
			}
			proc.Wait() // Don't care if this errors.
			os.Exit(3)
		}
	}

	<-done
}
