package fuzz

import (
	"errors"
	"fmt"
	"testing"
)

// recorder stands in for the fuzz target's T and counts its failures.
type recorder struct {
	testing.TB
	failures []string
}

func (r *recorder) Helper() {}

func (r *recorder) Fatalf(format string, args ...any) {
	r.failures = append(r.failures, fmt.Sprintf(format, args...))
}

type parseError struct{ Position int }

func (e *parseError) Error() string { return fmt.Sprintf("at %d", e.Position) }

type parseErrors []parseError

func (e parseErrors) Error() string { return "errors" }

func (e parseErrors) Unwrap() []error {
	out := make([]error, len(e))
	for i := range e {
		out[i] = &e[i]
	}
	return out
}

type Loc struct{ Start, End int }

type span struct{ Pos, End int }

// vcpu holds numbers, not a byte range, in fields called Start and End.
type vcpu struct{ Start, End int64 }

type node struct {
	Loc  Loc
	Kids []any
	Next *node
	CPU  *vcpu
}

type spanNode struct {
	span
	Name string
}

func TestError(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
		fail bool
	}{
		{"nil", nil, false},
		{"start", &parseError{0}, false},
		{"end", &parseError{3}, false},
		{"unknown", &parseError{-1}, true},
		{"past end", &parseError{4}, true},
		{"wrapped", fmt.Errorf("stmt: %w", &parseError{1}), false},
		{"other type", errors.New("boom"), true},
		{"list", parseErrors{{0}, {2}}, false},
		{"list past end", parseErrors{{0}, {9}}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := &recorder{TB: t}
			Error[*parseError](r, tc.err, "abc")
			if got := len(r.failures) > 0; got != tc.fail {
				t.Errorf("failed = %v, want %v: %q", got, tc.fail, r.failures)
			}
		})
	}
}

func TestTree(t *testing.T) {
	cyclic := &node{Loc: Loc{0, 1}}
	cyclic.Next = cyclic
	for _, tc := range []struct {
		name string
		tree any
		fail bool
	}{
		{"nil", nil, false},
		{"inside", &node{Loc: Loc{0, 3}}, false},
		{"unknown", &node{Loc: Loc{-1, -1}}, false},
		{"empty", &node{Loc: Loc{2, 2}}, false},
		{"end before start", &node{Loc: Loc{2, 1}}, true},
		{"past end", &node{Loc: Loc{0, 4}}, true},
		{"below unknown", &node{Loc: Loc{-2, 1}}, true},
		{"nested in a slice", []any{&node{Kids: []any{node{Loc: Loc{5, 6}}}}}, true},
		{"in a map", map[string]*node{"a": {Loc: Loc{0, 9}}}, true},
		{"embedded span", &spanNode{span: span{3, 1}}, true},
		{"cycle", cyclic, false},
		{"other Start and End", &node{CPU: &vcpu{100, 101}}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := &recorder{TB: t}
			Tree(r, tc.tree, "abc")
			if got := len(r.failures) > 0; got != tc.fail {
				t.Errorf("failed = %v, want %v: %q", got, tc.fail, r.failures)
			}
		})
	}
}
