A set of handy process running utilities, written in Go.

Documentation
-------------
 - [Package GoDoc Documentation](https://godoc.org/chrispennello.com/go/prun)
 - [prunevery](https://godoc.org/chrispennello.com/go/prun/cmd/prunevery):
   Enforce a minimum period between executions of a command.
 - [prunex](https://godoc.org/chrispennello.com/go/prun/cmd/prunex):
   Run a command exclusively (Unix only).
 - [prunfail](https://godoc.org/chrispennello.com/go/prun/cmd/prunfail):
   Guard the output of a potentially or intermittently failing command.
 - [prunfor](https://godoc.org/chrispennello.com/go/prun/cmd/prunfor):
   Run a command for an optionally limited amount of time.
 - [prunsleep](https://godoc.org/chrispennello.com/go/prun/cmd/prunsleep):
   Run a command after sleeping a random amount of time.

Installation
------------
Individually:

    go get chrispennello.com/go/prun/cmd/prunevery
    go get chrispennello.com/go/prun/cmd/prunex
    go get chrispennello.com/go/prun/cmd/prunfail
    go get chrispennello.com/go/prun/cmd/prunfor
    go get chrispennello.com/go/prun/cmd/prunsleep

Everything:

    go get chrispennello.com/go/prun/cmd/...

Future Work
-----------
 - Explicitly handle interrupts and other signals.
