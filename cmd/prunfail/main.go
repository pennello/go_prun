// chris 090415

// TODO Write me.
package main

import (
	"fmt"
	"log"
	"os"

	"path/filepath"

	"chrispennello.com/go/prun/cmd"
)

var state struct {
	cmd cmd.State

	// Log file name.
	logname string
}

func init() {
	log.SetFlags(0)
	state.cmd = cmd.Parse()

	tmp := os.TempDir()
	key := cmd.MakeKey(state.cmd.Cmd.Name, state.cmd.Cmd.Args)
	state.logname = filepath.Join(tmp, fmt.Sprintf("%s_%s.log", state.cmd.Me.Name, key))
}

func main() {
}
