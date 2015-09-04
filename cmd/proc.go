// chris 082815

package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"os/exec"
	"path/filepath"
)

// ErrNoEnt is returned by NewProc when the specified command cannot be
// found.
var ErrNoEnt = errors.New("not found")

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
// If command contains no path separators, Command uses LookPath to
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

// Kill simply calls Kill on the underlying os.Process.
func (p *Proc) Kill() error {
	return p.Process.Kill()
}
