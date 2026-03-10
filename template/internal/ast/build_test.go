package ast

import (
	"testing"
)

func TestStatementBuilderBuildCommand(t *testing.T) {
	b := &statementBuilder{
		action: "create",
		entity: "vpc",
		newparams: map[string]interface{}{
			"name": InterfaceNode{i: "myvpc"},
		},
	}
	stmt := b.build()
	if stmt == nil {
		t.Fatal("expected non-nil statement")
	}
	cmd, ok := stmt.Node.(*CommandNode)
	if !ok {
		t.Fatalf("expected *CommandNode, got %T", stmt.Node)
	}
	if cmd.Action != "create" || cmd.Entity != "vpc" {
		t.Fatalf("unexpected action/entity: %s %s", cmd.Action, cmd.Entity)
	}
}

func TestStatementBuilderBuildWithDeclaration(t *testing.T) {
	b := &statementBuilder{
		action:                "create",
		entity:                "vpc",
		declarationIdentifier: "myvar",
	}
	stmt := b.build()
	if stmt == nil {
		t.Fatal("expected non-nil statement")
	}
	decl, ok := stmt.Node.(*DeclarationNode)
	if !ok {
		t.Fatalf("expected *DeclarationNode, got %T", stmt.Node)
	}
	if decl.Ident != "myvar" {
		t.Fatalf("expected ident 'myvar', got %q", decl.Ident)
	}
}

func TestStatementBuilderBuildValue(t *testing.T) {
	b := &statementBuilder{
		isValue:               true,
		declarationIdentifier: "x",
		currentNode:           InterfaceNode{i: "hello"},
	}
	stmt := b.build()
	if stmt == nil {
		t.Fatal("expected non-nil statement")
	}
	decl, ok := stmt.Node.(*DeclarationNode)
	if !ok {
		t.Fatalf("expected *DeclarationNode, got %T", stmt.Node)
	}
	expr, ok := decl.Expr.(*RightExpressionNode)
	if !ok {
		t.Fatalf("expected *RightExpressionNode, got %T", decl.Expr)
	}
	if expr.i != (InterfaceNode{i: "hello"}) {
		t.Fatalf("unexpected expression value")
	}
}

func TestStatementBuilderBuildEmpty(t *testing.T) {
	b := &statementBuilder{}
	stmt := b.build()
	if stmt != nil {
		t.Fatal("expected nil statement for empty builder")
	}
}

func TestStatementBuilderAddParamKeyAndValue(t *testing.T) {
	b := &statementBuilder{}
	b.addParamKey("mykey")
	if b.currentKey != "mykey" {
		t.Fatalf("expected currentKey 'mykey', got %q", b.currentKey)
	}
	b.addParamValue(InterfaceNode{i: "myval"})
	if b.newparams["mykey"] != (InterfaceNode{i: "myval"}) {
		t.Fatalf("expected param to be set")
	}
	if b.currentKey != "" {
		t.Fatal("expected currentKey to be reset")
	}
}

func TestStatementBuilderAddParamValueNoKey(t *testing.T) {
	b := &statementBuilder{}
	b.addParamValue(InterfaceNode{i: "val"})
	// When there's no key, the value should be stored as currentNode
	// but not added to newparams
	if len(b.newparams) != 0 {
		t.Fatal("expected no params without key")
	}
}

func TestStatementBuilderListBuilder(t *testing.T) {
	b := &statementBuilder{}
	b.addParamKey("ids")
	b.newList()
	b.addParamValue(InterfaceNode{i: "a"})
	b.addParamValue(InterfaceNode{i: "b"})
	b.buildList()

	list, ok := b.newparams["ids"].(ListNode)
	if !ok {
		t.Fatalf("expected ListNode, got %T", b.newparams["ids"])
	}
	if len(list.arr) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(list.arr))
	}
}

func TestStatementBuilderConcatenationBuilder(t *testing.T) {
	a := &AST{}
	a.NewStatement()
	a.addAction("create")
	a.addEntity("vpc")
	a.addParamKey("name")
	a.addFirstValueInConcatenation()
	a.addStringValue("prefix-")
	a.addParamHoleValue("suffix")
	a.lastValueInConcatenation()
	a.StatementDone()

	if len(a.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(a.Statements))
	}
	cmd := a.Statements[0].Node.(*CommandNode)
	concat, ok := cmd.ParamNodes["name"].(ConcatenationNode)
	if !ok {
		t.Fatalf("expected ConcatenationNode, got %T", cmd.ParamNodes["name"])
	}
	if len(concat.arr) != 2 {
		t.Fatalf("expected 2 elements in concatenation, got %d", len(concat.arr))
	}
}

func TestASTNewStatementAndStatementDone(t *testing.T) {
	a := &AST{}
	a.NewStatement()
	a.addAction("create")
	a.addEntity("vpc")
	a.StatementDone()

	if len(a.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(a.Statements))
	}
}

func TestASTAddParamValue(t *testing.T) {
	a := &AST{}
	a.NewStatement()
	a.addAction("create")
	a.addEntity("vpc")
	a.addParamKey("count")
	a.addParamValue("42")
	a.StatementDone()

	cmd := a.Statements[0].Node.(*CommandNode)
	val, ok := cmd.ParamNodes["count"].(InterfaceNode)
	if !ok {
		t.Fatalf("expected InterfaceNode, got %T", cmd.ParamNodes["count"])
	}
	if val.i != 42 {
		t.Fatalf("expected 42, got %v", val.i)
	}
}

func TestASTAddParamValueFloat(t *testing.T) {
	a := &AST{}
	a.NewStatement()
	a.addAction("create")
	a.addEntity("vpc")
	a.addParamKey("ratio")
	a.addParamValue("3.14")
	a.StatementDone()

	cmd := a.Statements[0].Node.(*CommandNode)
	val := cmd.ParamNodes["ratio"].(InterfaceNode)
	if val.i != 3.14 {
		t.Fatalf("expected 3.14, got %v", val.i)
	}
}

func TestASTAddParamValueString(t *testing.T) {
	a := &AST{}
	a.NewStatement()
	a.addAction("create")
	a.addEntity("vpc")
	a.addParamKey("name")
	a.addParamValue("hello")
	a.StatementDone()

	cmd := a.Statements[0].Node.(*CommandNode)
	val := cmd.ParamNodes["name"].(InterfaceNode)
	if val.i != "hello" {
		t.Fatalf("expected 'hello', got %v", val.i)
	}
}

func TestASTAddParamValueLeadingZero(t *testing.T) {
	a := &AST{}
	a.NewStatement()
	a.addAction("create")
	a.addEntity("vpc")
	a.addParamKey("code")
	a.addParamValue("042")
	a.StatementDone()

	cmd := a.Statements[0].Node.(*CommandNode)
	val := cmd.ParamNodes["code"].(InterfaceNode)
	// Leading zero integers should be kept as strings
	if val.i != "042" {
		t.Fatalf("expected '042', got %v", val.i)
	}
}

func TestASTAddRefValue(t *testing.T) {
	a := &AST{}
	a.NewStatement()
	a.addAction("create")
	a.addEntity("vpc")
	a.addParamKey("id")
	a.addParamRefValue("myref")
	a.StatementDone()

	cmd := a.Statements[0].Node.(*CommandNode)
	ref, ok := cmd.ParamNodes["id"].(RefNode)
	if !ok {
		t.Fatalf("expected RefNode, got %T", cmd.ParamNodes["id"])
	}
	if ref.key != "myref" {
		t.Fatalf("expected 'myref', got %q", ref.key)
	}
}

func TestASTAddHoleValue(t *testing.T) {
	a := &AST{}
	a.NewStatement()
	a.addAction("create")
	a.addEntity("vpc")
	a.addParamKey("id")
	a.addParamHoleValue("myhole")
	a.StatementDone()

	cmd := a.Statements[0].Node.(*CommandNode)
	hole, ok := cmd.ParamNodes["id"].(HoleNode)
	if !ok {
		t.Fatalf("expected HoleNode, got %T", cmd.ParamNodes["id"])
	}
	if hole.key != "myhole" {
		t.Fatalf("expected 'myhole', got %q", hole.key)
	}
}

func TestASTAddAliasValue(t *testing.T) {
	a := &AST{}
	a.NewStatement()
	a.addAction("create")
	a.addEntity("vpc")
	a.addParamKey("id")
	a.addAliasParam("myalias")
	a.StatementDone()

	cmd := a.Statements[0].Node.(*CommandNode)
	alias, ok := cmd.ParamNodes["id"].(AliasNode)
	if !ok {
		t.Fatalf("expected AliasNode, got %T", cmd.ParamNodes["id"])
	}
	if alias.key != "myalias" {
		t.Fatalf("expected 'myalias', got %q", alias.key)
	}
}

func TestASTAddStringValue(t *testing.T) {
	a := &AST{}
	a.NewStatement()
	a.addAction("create")
	a.addEntity("vpc")
	a.addParamKey("name")
	a.addStringValue("my string value")
	a.StatementDone()

	cmd := a.Statements[0].Node.(*CommandNode)
	val := cmd.ParamNodes["name"].(InterfaceNode)
	if val.i != "my string value" {
		t.Fatalf("expected 'my string value', got %v", val.i)
	}
}

func TestASTListValues(t *testing.T) {
	a := &AST{}
	a.NewStatement()
	a.addAction("create")
	a.addEntity("vpc")
	a.addParamKey("ids")
	a.addFirstValueInList()
	a.addStringValue("id-1")
	a.addStringValue("id-2")
	a.lastValueInList()
	a.StatementDone()

	cmd := a.Statements[0].Node.(*CommandNode)
	list, ok := cmd.ParamNodes["ids"].(ListNode)
	if !ok {
		t.Fatalf("expected ListNode, got %T", cmd.ParamNodes["ids"])
	}
	if len(list.arr) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(list.arr))
	}
}

func TestASTDeclarationWithValue(t *testing.T) {
	a := &AST{}
	a.NewStatement()
	a.addDeclarationIdentifier("myvar")
	a.addValue()
	a.addStringValue("hello")
	a.StatementDone()

	if len(a.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(a.Statements))
	}
	decl, ok := a.Statements[0].Node.(*DeclarationNode)
	if !ok {
		t.Fatalf("expected *DeclarationNode, got %T", a.Statements[0].Node)
	}
	if decl.Ident != "myvar" {
		t.Fatalf("expected ident 'myvar', got %q", decl.Ident)
	}
}

func TestASTEmptyStatementDone(t *testing.T) {
	a := &AST{}
	a.NewStatement()
	a.StatementDone()

	if len(a.Statements) != 0 {
		t.Fatalf("expected 0 statements for empty builder, got %d", len(a.Statements))
	}
}

func TestASTAddActionPanicsOnInvalid(t *testing.T) {
	a := &AST{}
	a.NewStatement()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for invalid action")
		}
	}()
	a.addAction("invalidaction")
}

func TestASTAddEntityPanicsOnInvalid(t *testing.T) {
	a := &AST{}
	a.NewStatement()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for invalid entity")
		}
	}()
	a.addEntity("invalidentity")
}

func TestListValueBuilder(t *testing.T) {
	b := &listValueBuilder{}
	b.add(InterfaceNode{i: "a"})
	b.add(InterfaceNode{i: "b"})
	list := b.build()
	if len(list.arr) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(list.arr))
	}
}

func TestConcatenationValueBuilder(t *testing.T) {
	b := &concatenationValueBuilder{}
	b.add(InterfaceNode{i: "hello"})
	b.add(InterfaceNode{i: "world"})
	concat := b.build()
	if len(concat.arr) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(concat.arr))
	}
}
