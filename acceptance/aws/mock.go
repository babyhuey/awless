package awsat

import (
	"reflect"
	"testing"
)

type mock interface {
	Calls() map[string]int
	SetInputs(map[string]interface{})
	SetIgnored(map[string]struct{})
	SetTesting(*testing.T)
}

type basicMock struct { //nolint:unused // embedded in generated mock types
	t       *testing.T
	calls   map[string]int
	inputs  map[string]interface{}
	ignored map[string]struct{}
}

func (m *basicMock) addCall(call string) { //nolint:unused
	if m.calls == nil {
		m.calls = make(map[string]int)
	}
	m.calls[call]++
}

func (m *basicMock) Calls() map[string]int { //nolint:unused
	return m.calls
}

func (m *basicMock) SetTesting(t *testing.T) { //nolint:unused
	m.t = t
}

func (m *basicMock) SetInputs(inputs map[string]interface{}) { //nolint:unused
	m.inputs = inputs
}

func (m *basicMock) SetIgnored(ignored map[string]struct{}) { //nolint:unused
	m.ignored = ignored
}

func (m *basicMock) verifyInput(call string, got interface{}) { //nolint:unused
	if m.t == nil {
		return
	}
	if _, ok := m.ignored[call]; ok {
		return
	}
	if want, ok := m.inputs[call]; ok {
		if !reflect.DeepEqual(got, want) {
			m.t.Fatalf("got \n%#v\n\nwant \n%#v\n", got, want)
		}
	}
}
