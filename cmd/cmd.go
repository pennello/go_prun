// chris 082815

// Package cmd provides common code for command-line prun
// implementations.
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"os/exec"
	"path/filepath"
)

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
func NewProc(command string, args []string) (*Proc, error) {
	origCommand := command
	if filepath.Base(origCommand) == origCommand {
		if lp, err := exec.LookPath(origCommand); err != nil {
			return nil, err
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
		return nil, err
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

// Kill simply calls Kill on the underlying os.Process.
func (p *Proc) Kill() error {
	return p.Process.Kill()
}

func BadArgs(format string, a ...interface{}) {
	log.Printf(format, a...)
	os.Exit(2)
}

func ArgError(err error) {
	BadArgs("%v\n", err)
}
