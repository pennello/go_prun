// chris 090415

// TODO Write me.
// TODO Do something better than log.Fatal here.
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"path/filepath"

	"chrispennello.com/go/prun/cmd"
	"chrispennello.com/go/util/ringbuffer"
)

var state struct {
	cmd cmd.State

	// The maximum number of consecutive failures, after which the
	// output of the most recent failure will be emitted, and the
	// failure count will be reset.
	maxfail int

	// Log file name.
	logname string
}

func init() {
	log.SetFlags(0)
	state.cmd = cmd.Parse("maxfail")

	maxfail, err := strconv.ParseInt(state.cmd.Me.Args[0], 0, 0)
	if err != nil {
		cmd.ArgError(err)
	}
	if maxfail < 1 {
		cmd.BadArgs("maxfail must be positive")
	}
	state.maxfail = int(maxfail)

	tmp := os.TempDir()
	key := cmd.MakeKey(state.cmd.Cmd.Name, state.cmd.Cmd.Args)
	state.logname = filepath.Join(tmp, fmt.Sprintf("%s_%s.log", state.cmd.Me.Name, key))
}

func combinedOutput(proc *cmd.Proc) io.Reader {
	stdout, outerr := proc.Cmd.StdoutPipe()
	if outerr != nil {
		log.Fatal(outerr)
	}
	stderr, errerr := proc.Cmd.StderrPipe()
	if errerr != nil {
		log.Fatal(errerr)
	}
	return io.MultiReader(stdout, stderr)
}

func write(lf *logFile, rbuf *ringbuffer.B, tostderr bool) {
	data := rbuf.Bytes()
	if err := lf.write(data); err != nil {
		log.Fatal(err)
	}

	if tostderr {
		_, err := os.Stderr.Write(data)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func exit(lf *logFile, perr *cmd.ProcError, rbuf *ringbuffer.B) {
	lf.failures += 1
	if perr.Msg != "" {
		rbuf.Write([]byte(perr.Msg + "\n"))
	}
	write(lf, rbuf, lf.failures > state.maxfail)
	os.Exit(perr.Code)
}

func main() {
	proc := cmd.NewProc(state.cmd.Cmd.Name, state.cmd.Cmd.Args)

	// Copy combined output from process into ring buffer.  So in
	// the end, if a lot is written, we'll only be left with the
	// last bits.
	cout := combinedOutput(proc)
	rbuf := ringbuffer.New(maxLogSize)
	go func() {
		_, err := io.Copy(rbuf, cout)
		if err != nil {
			log.Fatal(err)
		}
	}()

	lf, lferr := newLogFile(state.logname)
	if lferr != nil {
		log.Fatal(lferr)
	}

	if perr := proc.StartError(); perr != nil {
		// There aren't any errors that could be caused by just
		// trying to start the process that we should elide.
		perr.Exit()
	}

	if perr := proc.WaitError(); perr != nil {
		// These are the errors that we'll want to potentially elide.
		exit(lf, perr, rbuf)
	}

	// Success: reset the failure count and not only do we log, but
	// the consolidated output also goes unconditionally to standard
	// error.
	lf.failures = 0
	write(lf, rbuf, true)
}
