// Copyright (c) 2023, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tree_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "cogentcore.org/core/tree"
	"cogentcore.org/core/tree/testdata"
	"cogentcore.org/core/types"
)

func TestNodeAddChild(t *testing.T) {
	parent := NewNodeBase()
	child := &NodeBase{}
	parent.AddChild(child)
	child.SetName("child1")
	assert.Equal(t, 1, len(parent.Kids))
	assert.Equal(t, parent, child.Parent())
	assert.Equal(t, "/node-base/child1", child.Path())
}

func TestNodeEmbedAddChild(t *testing.T) {
	parent := testdata.NewNodeEmbed()
	child := &testdata.NodeEmbed{}
	parent.AddChild(child)
	child.SetName("child1")
	assert.Len(t, parent.Kids, 1)
	assert.Equal(t, parent, child.Parent())
	assert.Equal(t, "/node-embed/child1", child.Path())
}

func TestNodeEmbedNewChild(t *testing.T) {
	parent := testdata.NewNodeEmbed()
	child := parent.NewChild(parent.NodeType())
	child.SetName("child1")
	assert.Len(t, parent.Kids, 1)
	assert.Equal(t, "/node-embed/child1", child.Path())
	assert.Equal(t, parent.NodeType(), child.NodeType())
}

func TestNodePath(t *testing.T) {
	parent := testdata.NewNodeEmbed()
	child1 := parent.NewChild(parent.NodeType())
	child2 := parent.NewChild(parent.NodeType())
	child3 := parent.NewChild(parent.NodeType())
	assert.Len(t, parent.Kids, 3)
	assert.Equal(t, "/node-embed/node-embed-0", child1.Path())
	assert.Equal(t, "/node-embed/node-embed-1", child2.Path())
	assert.Equal(t, "/node-embed/node-embed-2", child3.Path())
}

func TestNodeEscapePaths(t *testing.T) {
	parent := NewNodeBase()
	child1 := NewNodeBase(parent)
	child1.SetName("child1.go")
	child2 := NewNodeBase(parent)
	child2.SetName("child1/child1")
	child3 := NewNodeBase(parent)
	child3.SetName("child1/child1.go")
	schild2 := child2.NewChild(NodeBaseType)
	schild2.SetName("subchild1")
	assert.Len(t, parent.Kids, 3)
	assert.Equal(t, `/node-base/child1\,go`, child1.Path())
	assert.Equal(t, `/node-base/child1\\child1`, child2.Path())
	assert.Equal(t, `/node-base/child1\\child1\,go`, child3.Path())
	assert.Equal(t, `/node-base/child1\\child1/subchild1`, schild2.Path())
	assert.Equal(t, child1, parent.FindPath(child1.Path()))
	assert.Equal(t, child3, parent.FindPath(child3.Path()))
	assert.Equal(t, child3, parent.FindPath(child3.Path()))
	assert.Equal(t, schild2, parent.FindPath(schild2.Path()))
	assert.Equal(t, schild2, child2.FindPath(schild2.Path()))
}

func TestNodePathFrom(t *testing.T) {
	a := NewNodeBase()
	a.SetName("a")
	b := NewNodeBase(a)
	b.SetName("b")
	c := NewNodeBase(b)
	c.SetName("c")
	d := NewNodeBase(c)
	d.SetName("d")
	e := NewNodeBase(d)
	e.SetName("e")

	assert.Equal(t, "c/d", d.PathFrom(b))
}

func TestNodeDeleteChild(t *testing.T) {
	parent := NewNodeBase()
	child := NewNodeBase(parent)
	assert.Len(t, parent.Kids, 1)
	assert.True(t, parent.DeleteChild(child))
	assert.Len(t, parent.Kids, 0)
}

func TestNodeDeleteChildByName(t *testing.T) {
	parent := NewNodeBase()
	child := NewNodeBase(parent)
	child.SetName("child1")
	assert.Len(t, parent.Kids, 1)
	assert.True(t, parent.DeleteChildByName("child1"))
	assert.Len(t, parent.Kids, 0)
}

func TestNodeFindName(t *testing.T) {
	names := []string{"name0", "name1", "name2", "name3", "name4", "name5"}
	parent := NewNodeBase()
	for _, name := range names {
		child := NewNodeBase(parent)
		child.SetName(name)
	}
	assert.Len(t, parent.Kids, len(names))
	for i, nm := range names {
		for st := range names { // test all starting indexes
			idx, ok := parent.Children().IndexByName(nm, st)
			assert.True(t, ok)
			assert.Equal(t, i, idx)
		}
	}
}

func TestNodeFindType(t *testing.T) {
	parent := NewNodeBase()
	ne := testdata.NewNodeEmbed(parent)
	nb := NewNodeBase(parent)

	assert.True(t, ne.NodeType().HasEmbed(NodeBaseType))
	assert.True(t, nb.NodeType().HasEmbed(NodeBaseType))

	idx, ok := parent.Children().IndexByType(testdata.NodeEmbedType, NoEmbeds, 0)
	if assert.True(t, ok) {
		assert.Equal(t, 0, idx)
	}
	idx, ok = parent.Children().IndexByType(NodeBaseType, NoEmbeds, 0)
	if assert.True(t, ok) {
		assert.Equal(t, 1, idx)
	}
	_, err := parent.Children().ElemByTypeTry(NodeBaseType, NoEmbeds, 0)
	assert.NoError(t, err)
	idx, ok = parent.Children().IndexByType(NodeBaseType, Embeds, 0)
	if assert.True(t, ok) {
		assert.Equal(t, 0, idx)
	}
}

func TestNodeMove(t *testing.T) {
	parent := testdata.NewNodeEmbed()
	parent.Mbr1 = "bloop"
	parent.Mbr2 = 32
	child0 := testdata.NewNodeEmbed(parent)
	child0.SetName("child0")
	child1 := NewNodeBase(parent)
	child1.SetName("child1")
	child2 := testdata.NewNodeEmbed(parent)
	child2.SetName("child2")
	schild1 := child1.NewChild(testdata.NodeEmbedType)
	schild1.SetName("subchild1")
	child3 := NewNodeBase(parent)
	child3.SetName("child3")

	bf := fmt.Sprintf("mv before:\n%v\n", parent.Kids)
	parent.Children().Move(3, 1)
	a31 := fmt.Sprintf("mv 3 -> 1:\n%v\n", parent.Kids)
	parent.Children().Move(0, 3)
	a03 := fmt.Sprintf("mv 0 -> 3:\n%v\n", parent.Kids)
	parent.Children().Move(1, 2)
	a12 := fmt.Sprintf("mv 1 -> 2:\n%v\n", parent.Kids)

	bft := `mv before:
[/node-embed/child0 /node-embed/child1 /node-embed/child2 /node-embed/child3]
`
	assert.Equal(t, bft, bf)
	a31t := `mv 3 -> 1:
[/node-embed/child0 /node-embed/child3 /node-embed/child1 /node-embed/child2]
`
	assert.Equal(t, a31t, a31)
	a03t := `mv 0 -> 3:
[/node-embed/child3 /node-embed/child1 /node-embed/child2 /node-embed/child0]
`
	assert.Equal(t, a03t, a03)
	a12t := `mv 1 -> 2:
[/node-embed/child3 /node-embed/child2 /node-embed/child1 /node-embed/child0]
`
	assert.Equal(t, a12t, a12)
}

func TestNodeConfig(t *testing.T) {
	parent := testdata.NewNodeEmbed()
	parent.Mbr1 = "bloop"
	parent.Mbr2 = 32
	child0 := testdata.NewNodeEmbed(parent)
	child0.SetName("child0")
	child1 := testdata.NewNodeEmbed(parent)
	child1.SetName("child1")
	child2 := testdata.NewNodeEmbed(parent)
	child2.SetName("child2")
	schild1 := child1.NewChild(testdata.NodeEmbedType)
	schild1.SetName("subchild1")
	child3 := NewNodeBase(parent)
	child3.SetName("child3")

	config1 := Config{
		{Type: testdata.NodeEmbedType, Name: "child2"},
		{Type: testdata.NodeEmbedType, Name: "child3"},
		{Type: testdata.NodeEmbedType, Name: "child1"},
	}

	parent.ConfigChildren(config1)

	cf1 := fmt.Sprintf("config1:\n%v\n", parent.Kids)

	config2 := Config{
		{testdata.NodeEmbedType, "child4"},
		{NodeBaseType, "child1"}, // note: changing this to Node type removes child1.subchild1
		{testdata.NodeEmbedType, "child5"},
		{testdata.NodeEmbedType, "child3"},
		{testdata.NodeEmbedType, "child6"},
	}

	parent.ConfigChildren(config2)

	cf2 := fmt.Sprintf("config2:\n%v\n", parent.Kids)

	cf1t := `config1:
[/node-embed/child2 /node-embed/child3 /node-embed/child1]
`
	assert.Equal(t, cf1t, cf1)

	cf2t := `config2:
[/node-embed/child4 /node-embed/child1 /node-embed/child5 /node-embed/child3 /node-embed/child6]
`
	assert.Equal(t, cf2t, cf2)
}

func TestNodeWalk(t *testing.T) {
	parent := testdata.NodeEmbed{}
	parent.InitName(&parent, "par1")
	typ := parent.NodeType()
	parent.Mbr1 = "bloop"
	parent.Mbr2 = 32
	// child1 :=
	parent.NewChild(typ, "child1")
	child2 := parent.NewChild(typ, "child2")
	// child3 :=
	parent.NewChild(typ, "child3")
	schild2 := child2.NewChild(typ, "subchild1")

	res := []string{}

	schild2.WalkUp(func(k Node) bool {
		res = append(res, fmt.Sprintf("%v", k.Name()))
		return Continue
	})
	//	fmt.Printf("result: %v\n", res)

	trg := []string{"subchild1", "child2", "par1"}
	assert.Equal(t, trg, res)
	res = res[:0]

	parent.WalkDownPost(func(k Node) bool {
		return Continue
	},
		func(k Node) bool {
			res = append(res, fmt.Sprintf("[%v]", k.Name()))
			return Continue
		})
	// fmt.Printf("node field fun result: %v\n", res)
	trg = []string{"[child1]", "[subchild1]", "[child2]", "[child3]", "[par1]"}
	assert.Equal(t, trg, res)
	res = res[:0]

	// test for return = false working
	parent.WalkDownPost(func(k Node) bool {
		if k.Name() == "child2" {
			return Break
		}
		return Continue
	},
		func(k Node) bool {
			if k.Name() == "child2" {
				return Break
			}
			res = append(res, fmt.Sprintf("[%v]", k.Name()))
			return Continue
		})
	// fmt.Printf("node field fun result: %v\n", res)
	trg = []string{"[child1]", "[child3]", "[par1]"}
	assert.Equal(t, trg, res)
	res = res[:0]

	parent.WalkDownBreadth(func(k Node) bool {
		res = append(res, fmt.Sprintf("[%v]", k.Name()))
		return Continue
	})
	// fmt.Printf("node field fun result: %v\n", res)
	trg = []string{"[par1]", "[child1]", "[child2]", "[child3]", "[subchild1]"}
	assert.Equal(t, trg, res)
	res = res[:0]

	// test for return false
	parent.WalkDownBreadth(func(k Node) bool {
		if k.Name() == "child2" {
			return Break
		}
		res = append(res, fmt.Sprintf("[%v]", k.Name()))
		return Continue
	})
	// fmt.Printf("node field fun result: %v\n", res)
	trg = []string{"[par1]", "[child1]", "[child3]"}
	assert.Equal(t, trg, res)
	res = res[:0]
}

func TestNodeUpdate(t *testing.T) {
	parent := testdata.NodeEmbed{}
	parent.InitName(&parent, "par1")
	typ := parent.NodeType()
	parent.Mbr1 = "bloop"
	parent.Mbr2 = 32

	res := make([]string, 0, 10)
	parent.NewChild(typ, "child1")
	child2 := parent.NewChild(typ, "child2")
	parent.NewChild(typ, "child3")
	child2.NewChild(typ, "subchild1")

	parent.WalkDown(func(n Node) bool {
		res = append(res, n.Path())
		return Continue
	})
	assert.Equal(t, []string{"/par1", "/par1/child1", "/par1/child2", "/par1/child2/subchild1", "/par1/child3"}, res)
}

func TestProperties(t *testing.T) {
	n := testdata.NodeEmbed{}
	n.InitName(&n, "node")

	n.SetProperty("intprop", 42)
	assert.Equal(t, 42, n.Property("intprop"))

	n.SetProperty("floatprop", 42.0)
	assert.Equal(t, 42.0, n.Property("floatprop"))

	n.SetProperty("stringprop", "test string")
	assert.Equal(t, "test string", n.Property("stringprop"))

	n.DeleteProperty("floatprop")
	assert.Equal(t, nil, n.Property("floatprop"))

	assert.Equal(t, nil, n.Property("randomprop"))

	assert.Equal(t, map[string]any{"intprop": 42, "stringprop": "test string"}, n.Properties())
}

func TestPropertiesJSON(t *testing.T) {
	testProperties := map[string]any{
		"floatprop":  3.1415,
		"stringprop": "type string",
		"#subproperties": map[string]any{
			"sp1": "#FFE",
			"sp2": 42.2,
		},
	}

	b, err := json.MarshalIndent(testProperties, "", "  ")
	require.NoError(t, err)

	res := map[string]any{}
	err = json.Unmarshal(b, &res)
	require.NoError(t, err)

	assert.Equal(t, testProperties, res)
}

func TestDirectives(t *testing.T) {
	// test type directives: replacement for type properties
	n := testdata.NodeEmbed{}
	n.InitName(&n, "node")
	typ := n.NodeType()

	dir := typ.Directives[0]
	if dir.Tool != "direct" || dir.Directive != "value" {
		t.Errorf("Type directives error: directive should be `direct:value`, got %s", dir)
	}
}

func TestTreeMod(t *testing.T) {
	tree1 := NodeBase{}
	typ := tree1.NodeType()
	tree1.InitName(&tree1, "tree1")
	// child11 :=
	tree1.NewChild(typ, "child11")
	child12 := tree1.NewChild(typ, "child12")
	// child13 :=
	tree1.NewChild(typ, "child13")
	// schild12 :=
	child12.NewChild(typ, "subchild12")

	tree2 := NodeBase{}
	tree2.InitName(&tree2, "tree2")
	// child21 :=
	tree2.NewChild(typ, "child21")
	child22 := tree2.NewChild(typ, "child22")
	// child23 :=
	tree2.NewChild(typ, "child23")
	// schild22 :=
	child22.NewChild(typ, "subchild22")

	// fmt.Printf("#################################\n")

	// fmt.Printf("Trees before:\n%v%v", tree1, tree2)
	MoveToParent(child12.This(), tree2.This())

	// fmt.Printf("#################################\n")
	// fmt.Printf("Trees after add child12 move:\n%v%v", tree1, tree2)

	tree2.DeleteChild(child12)

	// fmt.Printf("#################################\n")

	// todo need actual tests in here!
}

/*

func TestNodeFieldJSONSave(t *testing.T) {
	parent := testdata.NodeField2{}
	parent.InitName(&parent, "par1")
	typ := parent.NodeType()
	parent.Mbr1 = "bloop"
	parent.Mbr2 = 32
	// child1 :=
	parent.NewChild(typ, "child1")
	child2 := parent.NewChild(typ, "child1").(*testdata.NodeField2)
	// child3 :=
	parent.NewChild(typ, "child1")
	child2.NewChild(typ, "subchild1")

	var buf bytes.Buffer
	err := parent.WriteJSON(&buf, true)
	if err != nil {
		t.Error(err)
		// } else {
		// 	fmt.Printf("json output:\n%v\n", string(buf.Bytes()))
	}
	b := buf.Bytes()

	tstload := testdata.NodeField2{}
	tstload.InitName(&tstload, "")
	err = tstload.ReadJSON(bytes.NewReader(b))
	if err != nil {
		t.Error(err)
	} else {
		var buf2 bytes.Buffer
		err = tstload.WriteJSON(&buf2, true)
		if err != nil {
			t.Error(err)
		}
		tstb := buf2.Bytes()
		// fmt.Printf("test loaded json output: %v\n", string(tstb))
		if !bytes.Equal(tstb, b) {
			t.Error("original and unmarshal'd json rep are not equivalent")
			os.WriteFile("/tmp/jsonout1", b, 0644)
			os.WriteFile("/tmp/jsonout2", tstb, 0644)
		}
	}

	nwnd, err := ReadNewJSON(bytes.NewReader(b))
	if err != nil {
		t.Error(err)
	} else {
		var buf2 bytes.Buffer
		err = nwnd.WriteJSON(&buf2, true)
		if err != nil {
			t.Error(err)
		}
		tstb := buf2.Bytes()
		// fmt.Printf("test loaded json output: %v\n", string(tstb))
		if !bytes.Equal(tstb, b) {
			t.Error("original and unmarshal'd json rep are not equivalent")
		}
	}
}

*/

func TestClone(t *testing.T) {
	parent := testdata.NodeField2{}
	parent.InitName(&parent, "par1")
	typ := parent.NodeType()
	parent.Mbr1 = "bloop"
	parent.Mbr2 = 32
	// child1 :=
	parent.NewChild(typ, "child1")
	child2 := parent.NewChild(typ, "child1").(*testdata.NodeField2)
	// child3 :=
	parent.NewChild(typ, "child1")
	child2.NewChild(typ, "subchild1")

	/*
		var buf bytes.Buffer
		err := parent.WriteJSON(&buf, true)
		if err != nil {
			t.Error(err)
			// } else {
			// 	fmt.Printf("json output:\n%v\n", string(buf.Bytes()))
		}
		b := buf.Bytes()

		tstload := parent.Clone()
		var buf2 bytes.Buffer
		err = tstload.WriteJSON(&buf2, true)
		if err != nil {
			t.Error(err)
		}
		tstb := buf2.Bytes()
		// fmt.Printf("test loaded json output: %v\n", string(tstb))
		if !bytes.Equal(tstb, b) {
			t.Error("original and unmarshal'd json rep are not equivalent")
			os.WriteFile("/tmp/jsonout1", b, 0644)
			os.WriteFile("/tmp/jsonout2", tstb, 0644)
		}
	*/
}

func TestAutoTypeName(t *testing.T) {
	root := &NodeBase{}
	root.InitName(root, "root")

	child := root.NewChild(NodeBaseType)
	assert.Equal(t, "node-base-0", child.Name())
}

// BuildGuiTreeSlow builds a tree that is typical of GUI structures where there are
// many widgets in a container and each widget has some number of parts.
// Uses slow AddChild method instead of fast one.
func BuildGuiTreeSlow(widgets, parts int, typ *types.Type) Node {
	win := NewOfType(typ)
	win.InitName(win, "window")

	vp := win.NewChild(typ, "vp")
	frame := vp.NewChild(typ, "frame")
	for wi := 0; wi < widgets; wi++ {
		widg := frame.NewChild(typ, fmt.Sprintf("widg_%d", wi))

		for pi := 0; pi < parts; pi++ {
			widg.NewChild(typ, fmt.Sprintf("part_%d", pi))
		}
	}
	return win
}

// BuildGuiTree builds a tree that is typical of GUI structures where there are
// many widgets in a container and each widget has some number of parts.
func BuildGuiTree(widgets, parts int, typ *types.Type) Node {
	win := NewOfType(typ)
	win.InitName(win, "window")

	vp := win.NewChild(typ, "vp")
	frame := vp.NewChild(typ, "frame")
	for wi := 0; wi < widgets; wi++ {
		widg := frame.NewChild(typ, fmt.Sprintf("widg_%d", wi))

		for pi := 0; pi < parts; pi++ {
			widg.NewChild(typ, fmt.Sprintf("part_%d", pi))
		}
	}
	return win
}

var TotNodes int
var TestGUITree_NodeEmbed Node
var TestGUITree_NodeField Node
var TestGUITree_NodeField2 Node

var NWidgets = 10000
var NParts = 5

func BenchmarkBuildGuiTree_NodeEmbed(b *testing.B) {
	for n := 0; n < b.N; n++ {
		wt := BuildGuiTree(NWidgets, NParts, testdata.NodeEmbedType)
		TestGUITree_NodeEmbed = wt
	}
}

func BenchmarkBuildGuiTree_NodeField(b *testing.B) {
	for n := 0; n < b.N; n++ {
		wt := BuildGuiTree(NWidgets, NParts, testdata.NodeFieldType)
		TestGUITree_NodeField = wt
	}
}

func BenchmarkBuildGuiTree_NodeField2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		wt := BuildGuiTree(NWidgets, NParts, testdata.NodeField2Type)
		TestGUITree_NodeField2 = wt
	}
}

func BenchmarkBuildGuiTreeSlow_NodeEmbed(b *testing.B) {
	// profile.Reset()
	// profile.Profiling = true
	for n := 0; n < b.N; n++ {
		wt := BuildGuiTreeSlow(NWidgets, NParts, testdata.NodeEmbedType)
		TestGUITree_NodeEmbed = wt
	}
	// profile.Report(time.Millisecond)
	// profile.Profiling = false
}

func BenchmarkWalkPre_NodeEmbed(b *testing.B) {
	wt := TestGUITree_NodeEmbed
	nnodes := 0
	for n := 0; n < b.N; n++ {
		wt.WalkDown(func(k Node) bool {
			nnodes++
			return Continue
		})
	}
	TotNodes = nnodes
	// fmt.Printf("tot nodes: %d\n", TotNodes)
}

func BenchmarkWalkPre_NodeField(b *testing.B) {
	wt := TestGUITree_NodeField
	nnodes := 0
	for n := 0; n < b.N; n++ {
		wt.WalkDown(func(k Node) bool {
			nnodes++
			return Continue
		})
	}
	TotNodes = nnodes
	// fmt.Printf("tot nodes: %d\n", TotNodes)
}

func BenchmarkWalkPre_NodeField2(b *testing.B) {
	wt := TestGUITree_NodeField2
	nnodes := 0
	for n := 0; n < b.N; n++ {
		wt.WalkDown(func(k Node) bool {
			nnodes++
			return Continue
		})
	}
	TotNodes = nnodes
	// fmt.Printf("tot nodes: %d\n", TotNodes)
}

func BenchmarkNewOfType(b *testing.B) {
	for n := 0; n < b.N; n++ {
		n := NewOfType(NodeBaseType)
		n.InitName(n)
	}
}

func BenchmarkStdNew(b *testing.B) {
	for n := 0; n < b.N; n++ {
		n := new(NodeBase)
		n.InitName(n)
	}
}
