// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Adapted in part from: https://github.com/magefile/mage
// Copyright presumably by Nate Finch, primary contributor
// Apache License, Version 2.0, January 2004

package xe

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Exec executes the command, piping its stdout and stderr to the given
// writers. If the command fails, it will return an error with the command output.
// Env is a list of environment variables to set when running the command,
// which override the current environment variables set
// (which are also passed to the command). cmd and args may include references
// to environment variables in $FOO format, in which case these will be
// expanded before the command is run.
//
// Ran reports if the command ran (rather than was not found or not executable).
// Code reports the exit code the command returned if it ran. If err == nil, ran
// is always true and code is always 0.
func Exec(env map[string]string, stdout, stderr io.Writer, cmd string, args ...string) (ran bool, err error) {
	expand := func(s string) string {
		s2, ok := env[s]
		if ok {
			return s2
		}
		return os.Getenv(s)
	}
	cmd = os.Expand(cmd, expand)
	for i := range args {
		args[i] = os.Expand(args[i], expand)
	}
	ran, code, err := run(env, stdout, stderr, cmd, args...)
	_ = code
	if err == nil {
		return true, nil
	}
	return ran, fmt.Errorf(`failed to run "%s %s: %v"`, cmd, strings.Join(args, " "), err)
}

func run(env map[string]string, stdout, stderr io.Writer, cmd string, args ...string) (ran bool, code int, err error) {
	c := exec.Command(cmd, args...)
	c.Env = os.Environ()
	for k, v := range env {
		c.Env = append(c.Env, k+"="+v)
	}
	c.Stderr = stderr
	c.Stdout = stdout
	c.Stdin = os.Stdin

	var quoted []string
	for i := range args {
		quoted = append(quoted, fmt.Sprintf("%q", args[i]))
	}
	// To protect against logging from doing exec in global variables
	// if mg.Verbose() {
	log.Println("exec:", cmd, strings.Join(quoted, " "))
	// }
	err = c.Run()
	return CmdRan(err), ExitStatus(err), err
}

// CmdRan examines the error to determine if it was generated as a result of a
// command running via os/exec.Command.  If the error is nil, or the command ran
// (even if it exited with a non-zero exit code), CmdRan reports true.  If the
// error is an unrecognized type, or it is an error from exec.Command that says
// the command failed to run (usually due to the command not existing or not
// being executable), it reports false.
func CmdRan(err error) bool {
	if err == nil {
		return true
	}
	ee, ok := err.(*exec.ExitError)
	if ok {
		return ee.Exited()
	}
	return false
}

type exitStatus interface {
	ExitStatus() int
}

// ExitStatus returns the exit status of the error if it is an exec.ExitError
// or if it implements ExitStatus() int.
// 0 if it is nil or 1 if it is a different error.
func ExitStatus(err error) int {
	if err == nil {
		return 0
	}
	if e, ok := err.(exitStatus); ok {
		return e.ExitStatus()
	}
	if e, ok := err.(*exec.ExitError); ok {
		if ex, ok := e.Sys().(exitStatus); ok {
			return ex.ExitStatus()
		}
	}
	return 1
}
