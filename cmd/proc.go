// chris 082815

package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"os/exec"
)

// ErrNoEnt is returned by NewProc when the specified command cannot be
// found.
var ErrNoEnt = errors.New("not found")

// ProcError represents a non-successful termination of the process.
// Print Msg to standard error and exit with status code Code.
type ProcError struct {
	Msg  string
	Code int
}

// Exit prints the message, if one is present, to standard error, and
// exits the parent process with the exit code.
func (pe *ProcError) Exit() {
	if pe.Msg != "" {
		log.Print(pe.Msg)
	}
	os.Exit(pe.Code)
}

// Proc is a wrapper around an exec.Cmd that provides common logic for
// prun command-line implementations.
type Proc struct {
	*exec.Cmd

	command string
	args    []string
}

// filterErrNoEnt replaces the given error with ErrNoEnt if appropriate
// and returns the original error otherwise.
func filterErrNoEnt(err error) error {
	execErr, ok := err.(*exec.Error)
	if ok {
		if execErr.Err == exec.ErrNotFound {
			return ErrNoEnt
		}
	}
	if os.IsNotExist(err) {
		return ErrNoEnt
	}
	return err
}

// NewProc returns the Proc struct to execute the named command with its
// optional arguments.
//
// Remember to set the Cmd Stdout and Stderr or capture with the
// corresponding pipes.  Otherwise, all output will be discarded.
func NewProc(command string, args []string) *Proc {
	return &Proc{
		Cmd: exec.Command(command, args...),
		command: command,
		args: args,
	}
}

// Start wraps the underlying exec.Cmd Start, filtering any returned
// errors and transforming them into an ErrNoEnt if appropriate.
func (p *Proc) Start() error {
	return filterErrNoEnt(p.Cmd.Start())
}

// StartError wraps Start.  It consolidates the various errors that can
// be returned into a single *ProcError.
//
// If the command could not be found, the exit status is 127.  For all
// other errors, the exit status is 1.
func (p *Proc) StartError() *ProcError {
	err := p.Start()
	if err != nil {
		if err == ErrNoEnt {
			return &ProcError{
				Msg:  fmt.Sprintf("%s: not found\n", p.command),
				Code: 127,
			}
		}
		return &ProcError{
			Msg:  err.Error(),
			Code: 1,
		}
		log.Fatal(err) // Implicitly exits with status 1.
	}
	return nil
}

// StartExit wraps StartError, but instead of returning an error when
// something goes wrong, it exits the parent process with the specified
// useful error message and exit status.
func (p *Proc) StartExit() {
	perr := p.StartError()
	if perr != nil {
		perr.Exit()
	}
}

// String returns a string representation of the command and its
// arguments.
func (p *Proc) String() string {
	if len(p.args) == 0 {
		return p.command
	}
	return fmt.Sprintf("%s %s", p.command, strings.Join(p.args, " "))
}

// Wait calls Wait on the underlying exec.Cmd's Process and, if the
// operating system supports it, returns the exit status.
//
// If an error occurs when waiting for the underlying process, the exit
// status will be -2, and the error will be returned.  If the operating
// system does not support determining the exit status, but the program
// exited successfully, the exit status will be 0.  If the operating
// system does not support determining the exit status and the program
// exited unsuccessfully, the exit status will be -1.
func (p *Proc) Wait() (exitStatus int, err error) {
	var ps *os.ProcessState
	ps, err = p.Cmd.Process.Wait()
	if err != nil {
		return -2, err
	}
	ws, ok := ps.Sys().(syscall.WaitStatus)
	if ok {
		return ws.ExitStatus(), nil
	}
	if ps.Success() {
		return 0, nil
	}
	return -1, nil
}

// WaitExit wraps Wait.  It consolidates the various errors that can be
// returned into a single *ProcError.
func (p *Proc) WaitError() *ProcError {
	exitStatus, err := p.Wait()
	if err != nil {
		return &ProcError{
			Msg:  err.Error(),
			Code: 1,
		}
	}
	if exitStatus != 0 {
		return &ProcError{
			Msg: "",
			Code: exitStatus,
		}
	}
	return nil
}

// WaitExit wraps WaitError and, given a *ProcError, exits the parent
// process with a useful message and exit status when something goes
// wrong.  If the underlying process exited successfully, it does
// nothing--that is, it does not exit the parent process).
func (p *Proc) WaitExit() {
	perr := p.WaitError()
	if perr != nil {
		perr.Exit()
	}
}
