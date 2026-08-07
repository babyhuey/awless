package awsat

import (
	"reflect"
	"testing"
)

// basicMock records calls and verifies inputs. Embedded in Mock.
type basicMock struct {
	t       *testing.T
	calls   map[string]int
	inputs  map[string]any
	ignored map[string]struct{}
}

func (m *basicMock) addCall(call string) {
	if m.calls == nil {
		m.calls = make(map[string]int)
	}
	m.calls[call]++
}

func (m *basicMock) Calls() map[string]int {
	return m.calls
}

func (m *basicMock) SetTesting(t *testing.T) {
	m.t = t
}

func (m *basicMock) SetInputs(inputs map[string]any) {
	m.inputs = inputs
}

func (m *basicMock) SetIgnored(ignored map[string]struct{}) {
	m.ignored = ignored
}

func (m *basicMock) verifyInput(call string, got any) {
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
