// Copyright (c) 2020, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ki

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testTree *NodeBase

func init() {
	testTree = &NodeBase{}
	typ := testTree.KiType()
	testTree.InitName(testTree, "root")
	// child1 :=
	testTree.NewChild(typ, "child0")
	var child2 = testTree.NewChild(typ, "child1")
	// child3 :=
	testTree.NewChild(typ, "child2")
	schild2 := child2.NewChild(typ, "subchild1")
	// sschild2 :=
	schild2.NewChild(typ, "subsubchild1")
	// child4 :=
	testTree.NewChild(typ, "child3")
}

func TestDown(t *testing.T) {
	cur := testTree
	res := []string{}
	for {
		res = append(res, cur.Path())
		curi := Next(cur)
		if curi == nil {
			break
		}
		cur = curi.(*NodeBase)
	}
	assert.Equal(t, []string{"/root", "/root/child0", "/root/child1", "/root/child1/subchild1", "/root/child1/subchild1/subsubchild1", "/root/child2", "/root/child3"}, res)
}

func TestUp(t *testing.T) {
	cur := Last(testTree)
	res := []string{}
	for {
		res = append(res, cur.Path())
		curi := Prev(cur)
		if curi == nil {
			break
		}
		cur = curi.(*NodeBase)
	}
	assert.Equal(t, []string{"/root/child3", "/root/child2", "/root/child1/subchild1/subsubchild1", "/root/child1/subchild1", "/root/child1", "/root/child0", "/root"}, res)
}
