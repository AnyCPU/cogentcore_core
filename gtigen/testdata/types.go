// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testdata

import (
	"fmt"
	"image/color"

	"goki.dev/gti"
)

// Person represents a person and their attributes.
// The zero value of a Person is not valid.
//
//ki:flagtype NodeFlags -field Flag
type Person struct { //goki:embedder
	color.RGBA
	// Name is the name of the person
	//gi:toolbar -hide
	Name string //goki:setter
	// Age is the age of the person
	//gi:view inline
	Age int `json:"-"`
	// Type is the type of the person
	Type *gti.Type
}

var _ = fmt.Stringer(&Person{})

//gti:skip
func (p Person) String() string { return p.Name }

// Introduction returns an introduction for the person.
// It contains the name of the person and their age.
//
//gi:toolbar -name ShowIntroduction -icon play -show-result -confirm
func (p *Person) Introduction() string { //gti:add
	return fmt.Sprintf("%s is %d years old", p.Name, p.Age)
}

// Alert prints an alert with the given message
func Alert(msg string) {
	fmt.Println("Alert:", msg)
}

// we test various type omitted arg combinations

func TypeOmittedArgs0(x, y float32)                {}
func TypeOmittedArgs1(x int, y, z float32)         {}
func TypeOmittedArgs2(x, y, z int)                 {}
func TypeOmittedArgs3(x int, y, z bool, w float32) {}
func TypeOmittedArgs4(x, y, z string, w bool)      {}
