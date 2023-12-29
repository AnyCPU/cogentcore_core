// Copyright (c) 2023, The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goki

import (
	"fmt"
	"runtime"
	"strings"

	"goki.dev/goki/config"
	"goki.dev/goki/mobile"
	"goki.dev/xe"
)

// Install installs the config package by looking for it in the list
// of supported packages. If the config ID is a filepath, it installs
// the package at that filepath on the local system. Install uses the
// same config info as build.
func Install(c *config.Config) error { //gti:add
	if c.Build.Package == "." || c.Build.Package == ".." || strings.Contains(c.Build.Package, "/") {
		return InstallLocal(c)
	}
	packages, err := LoadPackages()
	if err != nil {
		return fmt.Errorf("error loading packages: %w", err)
	}
	for _, pkg := range packages {
		if pkg.ID == c.Build.Package {
			return InstallPackage(pkg)
		}
	}
	return fmt.Errorf("error: could not find package %s", c.Build.Package)
}

// InstallPackage installs the given package object.
func InstallPackage(pkg Package) error {
	fmt.Println("Installing", pkg.Name)
	commands, ok := pkg.InstallCommands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("error: the requested package (%s) does not support your operating system (%s)", pkg.Name, runtime.GOOS)
	}
	for _, command := range commands {
		err := xe.Run(command.Name, command.Args...)
		if err != nil {
			return fmt.Errorf("error installing %s: %w", pkg.Name, err)
		}
	}
	return nil
}

// InstallLocal installs a local package from the filesystem
// on the user's device for the config target operating systems.
func InstallLocal(c *config.Config) error {
	for i, p := range c.Build.Target {
		err := config.OSSupported(p.OS)
		if err != nil {
			return fmt.Errorf("install: %w", err)
		}
		// if no arch is specified, we can assume it is the current arch,
		// as the user is running it (it could be a different arch when testing
		// on an external mobile device, but it is up to the user to specify
		// that arch in that case)
		if p.Arch == "*" {
			p.Arch = runtime.GOARCH
			c.Build.Target[i] = p
		}
		if p.OS == "android" || p.OS == "ios" {
			err := Build(c)
			if err != nil {
				return fmt.Errorf("error building: %w", err)
			}
			// we only want this target for install
			ot := c.Build.Target
			c.Build.Target = []config.Platform{p}
			err = mobile.Install(c)
			c.Build.Target = ot
			if err != nil {
				return fmt.Errorf("install: %w", err)
			}
			continue
		}
		if p.OS == "js" {
			// TODO: implement js
			continue
		}
		err = InstallLocalDesktop(c.Build.Package, p.OS)
		if err != nil {
			return fmt.Errorf("install: %w", err)
		}
	}
	return nil
}

// InstallLocalDesktop builds and installs an executable for the package at the given path for the given desktop platform.
// InstallLocalDesktop does not check whether operating systems are valid, so it should be called through Install in almost all cases.
func InstallLocalDesktop(pkgPath string, osName string) error {
	xc := xe.Major()
	xc.Env["GOOS"] = osName
	xc.Env["GOARCH"] = runtime.GOARCH
	err := xc.Run("go", "install", pkgPath)
	if err != nil {
		return fmt.Errorf("error installing on platform %s/%s: %w", osName, runtime.GOARCH, err)
	}
	return nil
}
