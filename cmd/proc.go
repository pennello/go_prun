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
	"path/filepath"
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

// Proc is a wrapper around an os.Process that provides some of the
// high-level conveniences of an exec.Cmd, but with some more of the
// low-level utility of an os.Process.
type Proc struct {
	*os.Process

	command string
	args    []string
}

// NewProc returns the Proc struct to execute the named command with its
// optional arguments.
//
// If command contains no path separators, Command uses exec.LookPath to
// resolve the path to a complete command if possible. Otherwise it uses
// command directly.
//
// It sets the current process's standard in, output, and error to be
// those used by the new process.
//
// If the specified command cannot be found, ErrNoEnt is returned.  If
// there is any other error trying to find or start the given command,
// that error is returned.
func NewProc(command string, args []string) (*Proc, error) {
	filterError := func(err error) error {
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

	origCommand := command
	if filepath.Base(origCommand) == origCommand {
		if lp, err := exec.LookPath(origCommand); err != nil {
			return nil, filterError(err)
		} else {
			command = lp
		}
	}
	argv := append([]string{origCommand}, args...)
	attr := new(os.ProcAttr)
	attr.Files = []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
	}
	process, err := os.StartProcess(command, argv, attr)
	if err != nil {
		return nil, filterError(err)
	}
	p := &Proc{
		Process: process,
		command: origCommand,
		args: args,
	}
	return p, nil
}

// NewProcExit wraps NewProc.  It consolidates the various errors that
// can be returned into a single *ProcError.
//
// If the command could not be found, the exit status is 127.  For all
// other errors, the exit status is 1.
func NewProcError(command string, args []string) (*Proc, *ProcError) {
	proc, err := NewProc(command, args)
	if err != nil {
		if err == ErrNoEnt {
			return nil, &ProcError{
				Msg:  fmt.Sprintf("%s: not found\n", command),
				Code: 127,
			}
		}
		return nil, &ProcError{
			Msg:  err.Error(),
			Code: 1,
		}
		log.Fatal(err) // Implicitly exits with status 1.
	}
	return proc, nil
}

// NewProcExit wraps NewProcErr, but instead of returning an error when
// something goes wrong, it exits the parent process with the specified
// useful error message and exit status.
func NewProcExit(command string, args []string) *Proc {
	proc, perr := NewProcError(command, args)
	if perr != nil {
		perr.Exit()
	}
	return proc
}

// String returns a string representation of the command and its
// arguments.
func (p *Proc) String() string {
	if len(p.args) == 0 {
		return p.command
	}
	return fmt.Sprintf("%s %s", p.command, strings.Join(p.args, " "))
}

// Wait calls Wait on the underlying os.Process and, if the operating
// system supports it, returns the exit status.
//
// If an error occurs when waiting for the underlying process, the exit
// status will be -2, and the error will be returned.  If the operating
// system does not support determining the exit status, but the program
// exited successfully, the exit status will be 0.  If the operating
// system does not support determining the exit status and the program
// exited unsuccessfully, the exit status will be -1.
func (p *Proc) Wait() (exitStatus int, err error) {
	var ps *os.ProcessState
	ps, err = p.Process.Wait()
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
// wrong.  If the underlying os.Process exited successfully, it does
// nothing--that is, it does not exit the parent process).
func (p *Proc) WaitExit() {
	perr := p.WaitError()
	if perr != nil {
		perr.Exit()
	}
}

// Kill simply calls Kill on the underlying os.Process.
func (p *Proc) Kill() error {
	return p.Process.Kill()
}
