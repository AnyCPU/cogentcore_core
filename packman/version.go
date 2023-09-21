// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packman

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"goki.dev/goki/config"
	"goki.dev/grease"
)

// SetVersion updates the config and version file of the config project based
// on the config version and commits and pushes the changes.
//
//gti:add
func SetVersion(c *config.Config) error {
	// we need to update the config file with the new version
	// TODO: determine correct config file instead of just first one
	err := grease.Save(c, grease.ConfigFiles[0])
	if err != nil {
		return fmt.Errorf("error saving new version to config file: %w", err)
	}
	str, err := VersionFileString(c)
	if err != nil {
		return fmt.Errorf("error generating version file string: %w", err)
	}
	err = os.WriteFile(c.Release.VersionFile, []byte(str), 0666)
	if err != nil {
		return fmt.Errorf("error writing version string to version file: %w", err)
	}
	err = PushVersionFileGit(c)
	if err != nil {
		return fmt.Errorf("error pushing version file to Git: %w", err)
	}
	return nil
}

// VersionFileString returns the version file string
// for a project with the given config info.
func VersionFileString(c *config.Config) (string, error) {
	var b strings.Builder
	b.WriteString("// Code generated by \"goki " + ArgsString(os.Args[1:]) + "\"; DO NOT EDIT.\n\n")
	b.WriteString("package " + c.Release.Package + "\n\n")
	b.WriteString("const (\n")
	b.WriteString("\t// Version is the version of this package being used\n")
	b.WriteString("\tVersion = \"" + c.Version + "\"\n")

	gc := exec.Command("git", "rev-parse", "--short", "HEAD")
	res, err := gc.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error getting previous git commit: %w (%s)", err, res)
	}
	b.WriteString("\t// GitCommit is the commit just before the latest version commit\n")
	b.WriteString("\tGitCommit = \"" + strings.TrimSuffix(string(res), "\n") + "\"\n")

	date := time.Now().UTC().Format("2006-01-02 15:04")
	b.WriteString("\t// VersionDate is the date-time of the latest version commit in UTC (in the format 'YYYY-MM-DD HH:MM', which is the Go format '2006-01-02 15:04')\n")
	b.WriteString("\tVersionDate = \"" + date + "\"\n")
	b.WriteString(")\n")
	return b.String(), nil
}

// PushVersionFileGit makes and pushes a Git commit
// updating the version file based on the given
// config info. It does not actually update the
// version file; it only commits and pushes the
// changes that should have already been made by
// [UpdateVersion].
func PushVersionFileGit(c *config.Config) error {
	ac := exec.Command("git", "add", c.Release.VersionFile)
	_, err := RunCmd(ac)
	if err != nil {
		return fmt.Errorf("error adding version file: %w", err)
	}

	cc := exec.Command("git", "commit", "-am", "updated version to "+c.Version)
	_, err = RunCmd(cc)
	if err != nil {
		return fmt.Errorf("error commiting release commit: %w", err)
	}

	pc := exec.Command("git", "push")
	_, err = RunCmd(pc)
	if err != nil {
		return fmt.Errorf("error pushing commit: %w", err)
	}

	return nil
}
