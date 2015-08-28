package cmd

import (
	"os"

	"os/exec"
	"path/filepath"
)

// Proc is a wrapper around an os.Process that provides some of the
// high-level conveniences of an exec.Cmd, but with some more of the
// low-level utility of an os.Process.
type Proc struct {
	*os.Process
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
func NewProc(command string, args []string) (*Proc, error) {
	if filepath.Base(command) == command {
		if lp, err := exec.LookPath(command); err != nil {
			return nil, err
		} else {
			command = lp
		}
	}
	argv := append([]string{command}, args...)
	attr := new(os.ProcAttr)
	attr.Files = []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
	}
	process, err := os.StartProcess(command, argv, attr)
	if err != nil {
		return nil, err
	}
	return &Proc{Process: process}, nil
}

// Wait calls Wait on the underlying os.Process and returns whether or
// not the command exited successfully.
func (p *Proc) Wait() (exitedSuccessfully bool, err error) {
	var ps *os.ProcessState
	ps, err = p.Process.Wait()
	if err != nil {
		return false, err
	}
	return ps.Success(), nil
}
