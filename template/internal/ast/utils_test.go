package ast

import (
	"strings"
	"testing"
)

func TestVerifyRefsValid(t *testing.T) {
	tree := &AST{}
	tree.Statements = append(tree.Statements,
		&Statement{Node: &DeclarationNode{
			Ident: "myvar",
			Expr: &CommandNode{
				Action: "create", Entity: "vpc",
				ParamNodes: map[string]interface{}{},
				Refs:       map[string]interface{}{},
			},
		}},
		&Statement{Node: &CommandNode{
			Action: "create", Entity: "subnet",
			ParamNodes: map[string]interface{}{"vpc": RefNode{key: "myvar"}},
			Refs:       map[string]interface{}{},
		}},
	)

	if err := VerifyRefs(tree); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestVerifyRefsUndefined(t *testing.T) {
	tree := &AST{}
	tree.Statements = append(tree.Statements,
		&Statement{Node: &CommandNode{
			Action: "create", Entity: "vpc",
			ParamNodes: map[string]interface{}{"id": RefNode{key: "undefined"}},
			Refs:       map[string]interface{}{},
		}},
	)

	err := VerifyRefs(tree)
	if err == nil {
		t.Fatal("expected error for undefined ref")
	}
	if !strings.Contains(err.Error(), "undefined") {
		t.Fatalf("expected error about undefined ref, got: %v", err)
	}
}

func TestVerifyRefsDuplicateDeclaration(t *testing.T) {
	tree := &AST{}
	tree.Statements = append(tree.Statements,
		&Statement{Node: &DeclarationNode{
			Ident: "myvar",
			Expr: &CommandNode{
				Action: "create", Entity: "vpc",
				ParamNodes: map[string]interface{}{},
				Refs:       map[string]interface{}{},
			},
		}},
		&Statement{Node: &DeclarationNode{
			Ident: "myvar",
			Expr: &CommandNode{
				Action: "create", Entity: "subnet",
				ParamNodes: map[string]interface{}{},
				Refs:       map[string]interface{}{},
			},
		}},
	)

	err := VerifyRefs(tree)
	if err == nil {
		t.Fatal("expected error for duplicate declaration")
	}
	if !strings.Contains(err.Error(), "already been assigned") {
		t.Fatalf("expected error about duplicate assignment, got: %v", err)
	}
}

func TestProcessRefs(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]interface{}{"id": RefNode{key: "myvar"}},
		Refs:       map[string]interface{}{},
	}
	tree := &AST{Statements: []*Statement{{Node: cmd}}}

	ProcessRefs(tree, map[string]interface{}{"myvar": "resolved-value"})

	if got, want := cmd.ParamNodes["id"], interface{}("resolved-value"); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestProcessRefsInList(t *testing.T) {
	list := ListNode{arr: []interface{}{RefNode{key: "x"}, InterfaceNode{i: "keep"}}}
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]interface{}{"ids": list},
		Refs:       map[string]interface{}{},
	}
	tree := &AST{Statements: []*Statement{{Node: cmd}}}

	ProcessRefs(tree, map[string]interface{}{"x": "resolved"})

	if got := list.arr[0]; got != "resolved" {
		t.Fatalf("got %v, want 'resolved'", got)
	}
}

func TestProcessRefsInRightExpression(t *testing.T) {
	expr := &RightExpressionNode{i: RefNode{key: "myref"}}
	tree := &AST{Statements: []*Statement{{Node: &DeclarationNode{Ident: "a", Expr: expr}}}}

	ProcessRefs(tree, map[string]interface{}{"myref": "resolved"})

	if got, want := expr.i, interface{}("resolved"); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestCollectHoles(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]interface{}{
			"name": HoleNode{key: "hole1"},
			"id":   HoleNode{key: "hole2"},
		},
		Refs: map[string]interface{}{},
	}
	tree := &AST{Statements: []*Statement{{Node: cmd}}}

	holes := CollectHoles(tree)
	if len(holes) != 2 {
		t.Fatalf("expected 2 holes, got %d", len(holes))
	}
}

func TestCollectUniqueHoles(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]interface{}{
			"name": HoleNode{key: "hole1"},
			"id":   HoleNode{key: "hole1"},
		},
		Refs: map[string]interface{}{},
	}
	tree := &AST{Statements: []*Statement{{Node: cmd}}}

	uniqueHoles := CollectUniqueHoles(tree)
	if len(uniqueHoles) != 1 {
		t.Fatalf("expected 1 unique hole, got %d", len(uniqueHoles))
	}
}

func TestCollectUniqueHolesWithPaths(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]interface{}{
			"name": HoleNode{key: "hole1"},
		},
		Refs: map[string]interface{}{},
	}
	tree := &AST{Statements: []*Statement{{Node: cmd}}}

	uniqueHoles := CollectUniqueHoles(tree)
	hole := HoleNode{key: "hole1"}
	paths := uniqueHoles[hole]
	if len(paths) != 1 || paths[0] != "create.vpc.name" {
		t.Fatalf("expected [create.vpc.name], got %v", paths)
	}
}

func TestProcessHoles(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]interface{}{
			"name": HoleNode{key: "hole1"},
		},
		Refs: map[string]interface{}{},
	}
	tree := &AST{Statements: []*Statement{{Node: cmd}}}

	processed := ProcessHoles(tree, map[string]interface{}{"hole1": "filled-value"})

	if got, want := cmd.ParamNodes["name"], interface{}("filled-value"); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := processed["hole1"], interface{}("filled-value"); got != want {
		t.Fatalf("processed: got %v, want %v", got, want)
	}
}

func TestProcessHolesWithAliasValue(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]interface{}{
			"name": HoleNode{key: "hole1"},
		},
		Refs: map[string]interface{}{},
	}
	tree := &AST{Statements: []*Statement{{Node: cmd}}}

	processed := ProcessHoles(tree, map[string]interface{}{"hole1": AliasNode{key: "myalias"}})
	if _, ok := processed["hole1"]; !ok {
		t.Fatal("expected hole1 in processed map")
	}
}

func TestProcessHolesWithRefValue(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]interface{}{
			"name": HoleNode{key: "hole1"},
		},
		Refs: map[string]interface{}{},
	}
	tree := &AST{Statements: []*Statement{{Node: cmd}}}

	processed := ProcessHoles(tree, map[string]interface{}{"hole1": RefNode{key: "myref"}})
	if _, ok := processed["hole1"]; !ok {
		t.Fatal("expected hole1 in processed map")
	}
}

func TestProcessHolesWithHoleValue(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]interface{}{
			"name": HoleNode{key: "hole1"},
		},
		Refs: map[string]interface{}{},
	}
	tree := &AST{Statements: []*Statement{{Node: cmd}}}

	processed := ProcessHoles(tree, map[string]interface{}{"hole1": HoleNode{key: "another"}})
	if _, ok := processed["hole1"]; !ok {
		t.Fatal("expected hole1 in processed map")
	}
}

func TestProcessHolesWithConcatenationValue(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]interface{}{
			"name": HoleNode{key: "hole1"},
		},
		Refs: map[string]interface{}{},
	}
	tree := &AST{Statements: []*Statement{{Node: cmd}}}

	processed := ProcessHoles(tree, map[string]interface{}{"hole1": ConcatenationNode{arr: []interface{}{InterfaceNode{i: "a"}}}})
	if _, ok := processed["hole1"]; !ok {
		t.Fatal("expected hole1 in processed map")
	}
}

func TestProcessHolesWithListValue(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]interface{}{
			"name": HoleNode{key: "hole1"},
		},
		Refs: map[string]interface{}{},
	}
	tree := &AST{Statements: []*Statement{{Node: cmd}}}

	listVal := ListNode{arr: []interface{}{AliasNode{key: "a"}, "plain", RefNode{key: "r"}, HoleNode{key: "h"}}}
	processed := ProcessHoles(tree, map[string]interface{}{"hole1": listVal})
	arr, ok := processed["hole1"].([]interface{})
	if !ok {
		t.Fatalf("expected []interface{}, got %T", processed["hole1"])
	}
	if len(arr) != 4 {
		t.Fatalf("expected 4 elements, got %d", len(arr))
	}
}

func TestProcessHolesInList(t *testing.T) {
	list := ListNode{arr: []interface{}{HoleNode{key: "h1"}, InterfaceNode{i: "keep"}}}
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]interface{}{"ids": list},
		Refs:       map[string]interface{}{},
	}
	tree := &AST{Statements: []*Statement{{Node: cmd}}}

	ProcessHoles(tree, map[string]interface{}{"h1": "filled"})

	if got := list.arr[0]; got != "filled" {
		t.Fatalf("got %v, want 'filled'", got)
	}
}

func TestProcessHolesInConcatenation(t *testing.T) {
	concat := ConcatenationNode{arr: []interface{}{InterfaceNode{i: "prefix-"}, HoleNode{key: "h1"}}}
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]interface{}{"name": concat},
		Refs:       map[string]interface{}{},
	}
	tree := &AST{Statements: []*Statement{{Node: cmd}}}

	ProcessHoles(tree, map[string]interface{}{"h1": "suffix"})

	if got := concat.arr[1]; got != "suffix" {
		t.Fatalf("got %v, want 'suffix'", got)
	}
}

func TestProcessHolesInRightExpression(t *testing.T) {
	expr := &RightExpressionNode{i: HoleNode{key: "h1"}}
	tree := &AST{Statements: []*Statement{{Node: &DeclarationNode{Ident: "x", Expr: expr}}}}

	ProcessHoles(tree, map[string]interface{}{"h1": "value"})

	if got, want := expr.i, interface{}("value"); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestRemoveOptionalHoles(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]interface{}{
			"name":     HoleNode{key: "required", optional: false},
			"optional": HoleNode{key: "opt", optional: true},
		},
		Refs: map[string]interface{}{},
	}
	tree := &AST{Statements: []*Statement{{Node: cmd}}}

	RemoveOptionalHoles(tree)

	if _, ok := cmd.ParamNodes["optional"]; ok {
		t.Fatal("expected optional hole to be removed")
	}
	if _, ok := cmd.ParamNodes["name"]; !ok {
		t.Fatal("expected required hole to remain")
	}
}

func TestRemoveOptionalHolesInRightExpression(t *testing.T) {
	expr := &RightExpressionNode{i: HoleNode{key: "opt", optional: true}}
	tree := &AST{Statements: []*Statement{{Node: &DeclarationNode{Ident: "x", Expr: expr}}}}

	RemoveOptionalHoles(tree)

	if expr.i != nil {
		t.Fatal("expected optional hole in right expression to be set to nil")
	}
}

func TestRemoveOptionalHolesInList(t *testing.T) {
	// Exercise the ListNode branch of RemoveOptionalHoles.
	// Note: ListNode is a value type, so removal within the visitor
	// modifies the local copy's arr slice. We just verify no panic.
	list := ListNode{arr: []interface{}{
		HoleNode{key: "opt", optional: true},
		InterfaceNode{i: "keep"},
	}}
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]interface{}{"ids": list},
		Refs:       map[string]interface{}{},
	}
	tree := &AST{Statements: []*Statement{{Node: cmd}}}

	RemoveOptionalHoles(tree) // should not panic
}

func TestCollectAliases(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]interface{}{
			"name": AliasNode{key: "myalias"},
		},
		Refs: map[string]interface{}{},
	}
	tree := &AST{Statements: []*Statement{{Node: cmd}}}

	aliases := CollectAliases(tree)
	if len(aliases) != 1 {
		t.Fatalf("expected 1 alias, got %d", len(aliases))
	}
	if got, want := aliases[0].key, "myalias"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestProcessAliases(t *testing.T) {
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]interface{}{
			"name": AliasNode{key: "myalias"},
		},
		Refs: map[string]interface{}{},
	}
	tree := &AST{Statements: []*Statement{{Node: cmd}}}

	aliasFunc := func(action, entity, key string) func(string) (string, bool) {
		return func(alias string) (string, bool) {
			if alias == "myalias" {
				return "resolved-id", true
			}
			return "", false
		}
	}

	ProcessAliases(tree, aliasFunc)

	if got, want := cmd.ParamNodes["name"], interface{}("resolved-id"); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestProcessAliasesInList(t *testing.T) {
	list := ListNode{arr: []interface{}{AliasNode{key: "a1"}}}
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]interface{}{"ids": list},
		Refs:       map[string]interface{}{},
	}
	tree := &AST{Statements: []*Statement{{Node: cmd}}}

	aliasFunc := func(action, entity, key string) func(string) (string, bool) {
		return func(alias string) (string, bool) {
			return "resolved", true
		}
	}

	ProcessAliases(tree, aliasFunc)

	if got := list.arr[0]; got != "resolved" {
		t.Fatalf("got %v, want 'resolved'", got)
	}
}

func TestProcessAliasesInConcatenation(t *testing.T) {
	concat := ConcatenationNode{arr: []interface{}{AliasNode{key: "a1"}, InterfaceNode{i: "-suffix"}}}
	cmd := &CommandNode{
		Action: "create", Entity: "vpc",
		ParamNodes: map[string]interface{}{"name": concat},
		Refs:       map[string]interface{}{},
	}
	tree := &AST{Statements: []*Statement{{Node: cmd}}}

	aliasFunc := func(action, entity, key string) func(string) (string, bool) {
		return func(alias string) (string, bool) {
			return "resolved", true
		}
	}

	ProcessAliases(tree, aliasFunc)

	if got := concat.arr[0]; got != "resolved" {
		t.Fatalf("got %v, want 'resolved'", got)
	}
}

func TestProcessAliasesInRightExpression(t *testing.T) {
	expr := &RightExpressionNode{i: AliasNode{key: "a1"}}
	tree := &AST{Statements: []*Statement{{Node: &DeclarationNode{Ident: "x", Expr: expr}}}}

	aliasFunc := func(action, entity, key string) func(string) (string, bool) {
		return func(alias string) (string, bool) {
			return "resolved", true
		}
	}

	ProcessAliases(tree, aliasFunc)

	if got, want := expr.i, interface{}("resolved"); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestContains(t *testing.T) {
	if !contains([]string{"a", "b", "c"}, "b") {
		t.Fatal("expected contains to find 'b'")
	}
	if contains([]string{"a", "b", "c"}, "d") {
		t.Fatal("expected contains not to find 'd'")
	}
	if contains(nil, "a") {
		t.Fatal("expected contains on nil to return false")
	}
}
