package ast

import (
	"reflect"
	"testing"
)

func TestRefNode(t *testing.T) {
	n := NewRefNode("myvar")
	if got, want := n.Ref(), "myvar"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if got, want := n.String(), "$myvar"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	clone := n.clone()
	if got, want := clone, Node(n); !reflect.DeepEqual(got, want) {
		t.Fatalf("clone mismatch: got %v, want %v", got, want)
	}
}

func TestAliasNode(t *testing.T) {
	n := NewAliasNode("my-alias")
	if got, want := n.Alias(), "my-alias"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if got, want := n.String(), "@my-alias"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	clone := n.clone()
	if got, want := clone, Node(n); !reflect.DeepEqual(got, want) {
		t.Fatalf("clone mismatch: got %v, want %v", got, want)
	}
}

func TestHoleNode(t *testing.T) {
	n := NewHoleNode("myhole")
	if got, want := n.Hole(), "myhole"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if got, want := n.String(), "{myhole}"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if n.IsOptional() {
		t.Fatal("expected non-optional hole")
	}
	clone := n.clone()
	if got, want := clone, Node(n); !reflect.DeepEqual(got, want) {
		t.Fatalf("clone mismatch: got %v, want %v", got, want)
	}
}

func TestOptionalHoleNode(t *testing.T) {
	n := NewOptionalHoleNode("opthole")
	if !n.IsOptional() {
		t.Fatal("expected optional hole")
	}
	if got, want := n.Hole(), "opthole"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestListNode(t *testing.T) {
	arr := []interface{}{"a", "b", "c"}
	n := NewListNode(arr)
	if got, want := n.String(), "[a,b,c]"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	elems := n.Elems()
	if got, want := len(elems), 3; got != want {
		t.Fatalf("got %d elems, want %d", got, want)
	}
	clone := n.clone()
	if got, want := clone, Node(n); !reflect.DeepEqual(got, want) {
		t.Fatalf("clone mismatch: got %v, want %v", got, want)
	}
}

func TestListNodeEmpty(t *testing.T) {
	n := NewListNode(nil)
	if got, want := n.String(), "[]"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestInterfaceNodeString(t *testing.T) {
	tcases := []struct {
		val    interface{}
		expect string
	}{
		{"simple", "simple"},
		{"hello world", "'hello world'"},
		{42, "42"},
		{3.14, "3.14"},
		{[]string{"a", "b"}, "[a,b]"},
		{true, "true"},
	}
	for _, tc := range tcases {
		n := InterfaceNode{i: tc.val}
		if got := n.String(); got != tc.expect {
			t.Errorf("InterfaceNode{%v}.String() = %q, want %q", tc.val, got, tc.expect)
		}
	}
}

func TestInterfaceNodeValue(t *testing.T) {
	n := InterfaceNode{i: "hello"}
	if got, want := n.Value(), interface{}("hello"); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestInterfaceNodeClone(t *testing.T) {
	n := InterfaceNode{i: 42}
	clone := n.clone()
	if got, want := clone, Node(n); !reflect.DeepEqual(got, want) {
		t.Fatalf("clone mismatch")
	}
}

func TestConcatenationNodeConcat(t *testing.T) {
	arr := []interface{}{
		InterfaceNode{i: "hello"},
		InterfaceNode{i: "-"},
		InterfaceNode{i: "world"},
	}
	n := NewConcatenationNode(arr)
	if got, want := n.Concat(), "hello-world"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestConcatenationNodeConcatWithNonInterface(t *testing.T) {
	arr := []interface{}{
		InterfaceNode{i: "prefix"},
		"suffix",
	}
	n := NewConcatenationNode(arr)
	if got, want := n.Concat(), "prefixsuffix"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestConcatenationNodeStringNoHoles(t *testing.T) {
	arr := []interface{}{
		InterfaceNode{i: "hello"},
		InterfaceNode{i: "world"},
	}
	n := NewConcatenationNode(arr)
	got := n.String()
	if got != "helloworld" {
		t.Fatalf("got %q, want %q", got, "helloworld")
	}
}

func TestConcatenationNodeStringWithHoles(t *testing.T) {
	arr := []interface{}{
		InterfaceNode{i: "prefix"},
		HoleNode{key: "myhole"},
	}
	n := NewConcatenationNode(arr)
	got := n.String()
	// With holes, strings are quoted and joined by +
	if got != "'prefix'+{myhole}" {
		t.Fatalf("got %q, want %q", got, "'prefix'+{myhole}")
	}
}

func TestConcatenationNodeClone(t *testing.T) {
	arr := []interface{}{InterfaceNode{i: "a"}}
	n := NewConcatenationNode(arr)
	clone := n.clone()
	if got, want := clone, Node(n); !reflect.DeepEqual(got, want) {
		t.Fatalf("clone mismatch")
	}
}

func TestRightExpressionNodeResultInterfaceNode(t *testing.T) {
	n := &RightExpressionNode{i: InterfaceNode{i: "hello"}}
	if got, want := n.Result(), interface{}("hello"); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if n.Err() != nil {
		t.Fatalf("expected nil error for InterfaceNode")
	}
}

func TestRightExpressionNodeResultRefNode(t *testing.T) {
	n := &RightExpressionNode{i: RefNode{key: "x"}}
	if n.Result() != nil {
		t.Fatal("expected nil result for RefNode")
	}
}

func TestRightExpressionNodeResultAliasNode(t *testing.T) {
	n := &RightExpressionNode{i: AliasNode{key: "x"}}
	if n.Result() != nil {
		t.Fatal("expected nil result for AliasNode")
	}
}

func TestRightExpressionNodeResultHoleNode(t *testing.T) {
	n := &RightExpressionNode{i: HoleNode{key: "x"}}
	if n.Result() != nil {
		t.Fatal("expected nil result for HoleNode")
	}
}

func TestRightExpressionNodeResultListNode(t *testing.T) {
	n := &RightExpressionNode{i: ListNode{arr: []interface{}{InterfaceNode{i: 1}, InterfaceNode{i: 2}}}}
	result := n.Result()
	arr, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected []interface{}, got %T", result)
	}
	if len(arr) != 2 || arr[0] != 1 || arr[1] != 2 {
		t.Fatalf("unexpected result: %v", arr)
	}
}

func TestRightExpressionNodeResultListWithRef(t *testing.T) {
	n := &RightExpressionNode{i: ListNode{arr: []interface{}{RefNode{key: "x"}}}}
	if n.Result() != nil {
		t.Fatal("expected nil result for list containing ref")
	}
}

func TestRightExpressionNodeResultListWithAlias(t *testing.T) {
	n := &RightExpressionNode{i: ListNode{arr: []interface{}{AliasNode{key: "x"}}}}
	if n.Result() != nil {
		t.Fatal("expected nil result for list containing alias")
	}
}

func TestRightExpressionNodeResultListWithHole(t *testing.T) {
	n := &RightExpressionNode{i: ListNode{arr: []interface{}{HoleNode{key: "x"}}}}
	if n.Result() != nil {
		t.Fatal("expected nil result for list containing hole")
	}
}

func TestRightExpressionNodeResultListWithPlainValue(t *testing.T) {
	n := &RightExpressionNode{i: ListNode{arr: []interface{}{"plain"}}}
	result := n.Result()
	arr, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected []interface{}, got %T", result)
	}
	if len(arr) != 1 || arr[0] != "plain" {
		t.Fatalf("unexpected result: %v", arr)
	}
}

func TestRightExpressionNodeResultConcatenation(t *testing.T) {
	n := &RightExpressionNode{i: ConcatenationNode{arr: []interface{}{InterfaceNode{i: "a"}, InterfaceNode{i: "b"}}}}
	if got, want := n.Result(), interface{}("ab"); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestRightExpressionNodeResultDefault(t *testing.T) {
	n := &RightExpressionNode{i: "plain-string"}
	if got, want := n.Result(), interface{}("plain-string"); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestRightExpressionNodeErr(t *testing.T) {
	n := &RightExpressionNode{i: "not-interface-node"}
	if n.Err() == nil {
		t.Fatal("expected error for non-InterfaceNode")
	}
}

func TestRightExpressionNodeString(t *testing.T) {
	n := &RightExpressionNode{i: InterfaceNode{i: "test"}}
	got := n.String()
	if got == "" {
		t.Fatal("expected non-empty string")
	}
}

func TestRightExpressionNodeNode(t *testing.T) {
	inner := InterfaceNode{i: 42}
	n := &RightExpressionNode{i: inner}
	if got, want := n.Node(), interface{}(inner); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestRightExpressionNodeClone(t *testing.T) {
	n := &RightExpressionNode{i: InterfaceNode{i: "x"}}
	clone := n.clone()
	if got, want := clone.(*RightExpressionNode).i, n.i; !reflect.DeepEqual(got, want) {
		t.Fatalf("clone mismatch")
	}
}
