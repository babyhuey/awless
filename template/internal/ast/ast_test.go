/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ast

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/bootswithdefer/awless/template/env"
	"github.com/bootswithdefer/awless/template/params"
)

func TestCloneAST(t *testing.T) {
	tree := &AST{}

	cmd := new(fakeCmd)

	tree.Statements = append(tree.Statements, &Statement{Node: &DeclarationNode{
		Ident: "myvar",
		Expr: &CommandNode{
			Action: "create", Entity: "vpc",
			ParamNodes: map[string]any{"count": InterfaceNode{i: 1}, "myname": RefNode{key: "name"}},
			Refs:       map[string]any{"here": "there"},
		}}}, &Statement{Node: &DeclarationNode{
		Ident: "myothervar",
		Expr: &CommandNode{
			Command: cmd,
			Action:  "create", Entity: "subnet",
			ParamNodes: map[string]any{"vpc": HoleNode{key: "myvar"}},
			Refs:       map[string]any{},
		}}}, &Statement{Node: &CommandNode{
		Action: "create", Entity: "instance",
		ParamNodes: map[string]any{"subnet": HoleNode{key: "myothervar"}},
		Refs:       map[string]any{"donald": "duck"},
	}},
	)

	clone := tree.Clone()

	if got, want := clone, tree; !reflect.DeepEqual(got, want) {
		t.Fatalf("\ngot %#v\n\nwant %#v", got, want)
	}

	clone.Statements[0].Node.(*DeclarationNode).Expr.(*CommandNode).ParamNodes["new"] = InterfaceNode{i: "mynode"}

	if got, want := clone.Statements, tree.Statements; reflect.DeepEqual(got, want) {
		t.Fatalf("\ngot %s\n\nwant %s", got, want)
	}
}

func TestIsQuoted(t *testing.T) {
	tcases := []struct {
		in  string
		out bool
	}{
		{"", false},
		{"'", false},
		{"\"", false},
		{"''", true},
		{"\"\"", true},
		{"\"'", false},
		{"'\"", false},
		{"'test\"", false},
		{"\"test'", false},
		{"\"test\"", true},
		{"'test'", true},
	}
	for i, tcase := range tcases {
		if got, want := isQuoted(tcase.in), tcase.out; got != want {
			t.Fatalf("%d: got %t, want %t", i+1, got, want)
		}
	}
}

func TestCommandNodeString(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]any{
			"name": InterfaceNode{i: "myvpc"},
			"cidr": InterfaceNode{i: "10.0.0.0/16"},
		},
		Refs: map[string]any{},
	}
	got := cmd.String()
	if !strings.HasPrefix(got, "create vpc") {
		t.Fatalf("expected to start with 'create vpc', got %q", got)
	}
	if !strings.Contains(got, "name=myvpc") {
		t.Fatalf("expected 'name=myvpc' in output, got %q", got)
	}
	if !strings.Contains(got, "cidr=10.0.0.0/16") {
		t.Fatalf("expected 'cidr=10.0.0.0/16' in output, got %q", got)
	}
}

func TestCommandNodeStringWithRefs(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "subnet",
		ParamNodes: map[string]any{},
		Refs:       map[string]any{"vpc": "$myvpc"},
	}
	got := cmd.String()
	if !strings.Contains(got, "vpc=$myvpc") {
		t.Fatalf("expected ref in output, got %q", got)
	}
}

func TestCommandNodeStringWithList(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "securitygroup",
		ParamNodes: map[string]any{
			"ports": []any{"80", "443"},
		},
		Refs: map[string]any{},
	}
	got := cmd.String()
	// String values that look like integers get quoted
	if !strings.Contains(got, "ports=[") {
		t.Fatalf("expected list in output, got %q", got)
	}
}

func TestCommandNodeStringWithInterfaceList(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "securitygroup",
		ParamNodes: map[string]any{
			"ports": []any{80, 443},
		},
		Refs: map[string]any{},
	}
	got := cmd.String()
	if !strings.Contains(got, "ports=[80,443]") {
		t.Fatalf("expected list in output, got %q", got)
	}
}

func TestCommandNodeStringWithNonStringNonList(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]any{
			"count": InterfaceNode{i: 5},
		},
		Refs: map[string]any{},
	}
	got := cmd.String()
	if !strings.Contains(got, "count=") {
		t.Fatalf("expected 'count=' in output, got %q", got)
	}
}

func TestCommandNodeStringNoParams(t *testing.T) {
	cmd := &CommandNode{
		Action: "delete", Entity: "vpc",
		ParamNodes: map[string]any{},
		Refs:       map[string]any{},
	}
	if got, want := cmd.String(), "delete vpc"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestCommandNodeKeys(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]any{"name": InterfaceNode{i: "x"}, "cidr": InterfaceNode{i: "y"}},
		Refs:       map[string]any{"ref1": "val"},
	}
	keys := cmd.Keys()
	sort.Strings(keys)
	expected := []string{"cidr", "name", "ref1"}
	if !reflect.DeepEqual(keys, expected) {
		t.Fatalf("got %v, want %v", keys, expected)
	}
}

func TestCommandNodeResultAndErr(t *testing.T) {
	cmd := &CommandNode{
		CmdResult:  "result-value",
		CmdErr:     nil,
		ParamNodes: map[string]any{},
		Refs:       map[string]any{},
	}
	if got, want := cmd.Result(), any("result-value"); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if cmd.Err() != nil {
		t.Fatal("expected nil error")
	}
}

func TestCommandNodeClone(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]any{"name": InterfaceNode{i: "x"}},
		Refs:       map[string]any{"ref": "val"},
	}
	clone := cmd.clone().(*CommandNode)
	if clone.Action != cmd.Action || clone.Entity != cmd.Entity {
		t.Fatal("clone action/entity mismatch")
	}
	// Mutate clone and ensure original is not affected
	clone.ParamNodes["new"] = InterfaceNode{i: "y"}
	if _, ok := cmd.ParamNodes["new"]; ok {
		t.Fatal("original should not be mutated")
	}
}

func TestCommandNodeProcessRefs(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "subnet",
		ParamNodes: map[string]any{},
		Refs:       map[string]any{"vpc": RefNode{key: "myvar"}},
	}
	cmd.ProcessRefs(map[string]any{"myvar": "vpc-12345"})

	if got, want := cmd.ParamNodes["vpc"], any("vpc-12345"); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestCommandNodeProcessRefsWithList(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "subnet",
		ParamNodes: map[string]any{},
		Refs: map[string]any{
			"ids": ListNode{arr: []any{RefNode{key: "id1"}, InterfaceNode{i: "static"}}},
		},
	}
	cmd.ProcessRefs(map[string]any{"id1": "resolved-id"})

	arr, ok := cmd.ParamNodes["ids"].([]any)
	if !ok {
		t.Fatalf("expected []any, got %T", cmd.ParamNodes["ids"])
	}
	if arr[0] != "resolved-id" {
		t.Fatalf("expected 'resolved-id', got %v", arr[0])
	}
}

func TestCommandNodeToDriverParams(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]any{
			"name":  InterfaceNode{i: "myvpc"},
			"ref":   RefNode{key: "x"},
			"hole":  HoleNode{key: "y"},
			"alias": AliasNode{key: "z"},
			"plain": "plain-value",
		},
		Refs: map[string]any{},
	}
	params := cmd.ToDriverParams()

	if params["name"] != "myvpc" {
		t.Fatalf("expected 'myvpc', got %v", params["name"])
	}
	if params["plain"] != "plain-value" {
		t.Fatalf("expected 'plain-value', got %v", params["plain"])
	}
	// RefNode, HoleNode, AliasNode should be excluded
	if _, ok := params["ref"]; ok {
		t.Fatal("RefNode should be excluded from driver params")
	}
	if _, ok := params["hole"]; ok {
		t.Fatal("HoleNode should be excluded from driver params")
	}
	if _, ok := params["alias"]; ok {
		t.Fatal("AliasNode should be excluded from driver params")
	}
}

func TestCommandNodeToFillerParams(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]any{
			"name":  InterfaceNode{i: "myvpc"},
			"alias": AliasNode{key: "z"},
			"ref":   RefNode{key: "x"},
		},
		Refs: map[string]any{},
	}
	params := cmd.ToFillerParams()

	if params["name"] != "myvpc" {
		t.Fatalf("expected 'myvpc', got %v", params["name"])
	}
	if _, ok := params["alias"]; !ok {
		t.Fatal("AliasNode should be included in filler params")
	}
	// RefNode should not appear
	if _, ok := params["ref"]; ok {
		t.Fatal("RefNode should not be in filler params")
	}
}

func TestCommandNodeToFillerParamsWithList(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]any{
			"ids": ListNode{arr: []any{InterfaceNode{i: "id1"}, AliasNode{key: "a1"}}},
		},
		Refs: map[string]any{},
	}
	params := cmd.ToFillerParams()
	_, ok := params["ids"]
	if !ok {
		t.Fatal("expected 'ids' in filler params")
	}
}

func TestDeclarationNodeString(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]any{},
		Refs:       map[string]any{},
	}
	decl := &DeclarationNode{Ident: "myvar", Expr: cmd}
	got := decl.String()
	if !strings.HasPrefix(got, "myvar = ") {
		t.Fatalf("expected to start with 'myvar = ', got %q", got)
	}
}

func TestDeclarationNodeClone(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]any{},
		Refs:       map[string]any{},
	}
	decl := &DeclarationNode{Ident: "myvar", Expr: cmd}
	clone := decl.clone().(*DeclarationNode)
	if clone.Ident != decl.Ident {
		t.Fatal("ident mismatch in clone")
	}
}

func TestDeclarationNodeCloneNilExpr(t *testing.T) {
	decl := &DeclarationNode{Ident: "myvar", Expr: nil}
	clone := decl.clone().(*DeclarationNode)
	if clone.Expr != nil {
		t.Fatal("expected nil Expr in clone")
	}
}

func TestStatementClone(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]any{"k": InterfaceNode{i: "v"}},
		Refs:       map[string]any{},
	}
	stmt := &Statement{Node: cmd}
	clone := stmt.Clone()
	clonedCmd := clone.Node.(*CommandNode)
	clonedCmd.ParamNodes["new"] = InterfaceNode{i: "x"}
	if _, ok := cmd.ParamNodes["new"]; ok {
		t.Fatal("original should not be mutated")
	}
}

func TestASTString(t *testing.T) {
	tree := &AST{}
	tree.Statements = append(tree.Statements,
		&Statement{Node: &CommandNode{
			Action: "create", Entity: "vpc",
			ParamNodes: map[string]any{"name": InterfaceNode{i: "myvpc"}},
			Refs:       map[string]any{},
		}},
		&Statement{Node: &CommandNode{
			Action: "delete", Entity: "subnet",
			ParamNodes: map[string]any{},
			Refs:       map[string]any{},
		}},
	)
	got := tree.String()
	lines := strings.Split(got, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %q", len(lines), got)
	}
}

func TestASTCloneNode(t *testing.T) {
	tree := &AST{}
	tree.Statements = append(tree.Statements,
		&Statement{Node: &CommandNode{
			Action: "create", Entity: "vpc",
			ParamNodes: map[string]any{},
			Refs:       map[string]any{},
		}},
	)
	cloneNode := tree.clone()
	clone, ok := cloneNode.(*AST)
	if !ok {
		t.Fatalf("expected *AST, got %T", cloneNode)
	}
	if len(clone.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(clone.Statements))
	}
}

type fakeCmd struct{}

func (*fakeCmd) ParamsSpec() params.Spec                      { return nil }
func (*fakeCmd) Run(env.Running, map[string]any) (any, error) { return nil, nil }
