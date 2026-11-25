// Package fuzz holds an engine parser to the engine contract on any input:
// Parse never panics, fails only with the engine's ParseError at a byte
// offset inside the input, and returns a tree whose every byte range lies
// inside the input. Each engine's parser package has one FuzzParse built
// from it; `go test` replays its compat.tsv as the seed corpus, and
// `go test -fuzz=FuzzParse` searches beyond it.
package fuzz

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/bytebase/omni/enginetest/compat"
)

// Seed adds the empty input and every distinct statement of the compat.tsv
// at path to f's seed corpus.
func Seed(f *testing.F, path string) {
	f.Helper()
	file, err := compat.ReadFile(path)
	if err != nil {
		f.Fatal(err)
	}
	seen := map[string]bool{"": true}
	f.Add("")
	for _, r := range file.Rows {
		if !seen[r.Statement] {
			seen[r.Statement] = true
			f.Add(r.Statement)
		}
	}
}

// Error fails t unless err is nil or holds an E, and every E it holds
// (each one of a list that unwraps to several) reports a byte offset in
// [0, len(input)] in its Position field.
func Error[E error](t testing.TB, err error, input string) {
	t.Helper()
	if err == nil {
		return
	}
	var e E
	if !errors.As(err, &e) {
		t.Fatalf("Parse(%s) failed with %T, not %T: %v", quote(input), err, e, err)
		return
	}
	if list, ok := err.(interface{ Unwrap() []error }); ok {
		for _, each := range list.Unwrap() {
			Error[E](t, each, input)
		}
		return
	}
	if pos := position(reflect.ValueOf(e)); pos < 0 || pos > len(input) {
		t.Fatalf("Parse(%s): %T at byte %d, outside the %d-byte input: %v", quote(input), e, pos, len(input), err)
	}
}

func position(v reflect.Value) int {
	for v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if f := v.FieldByName("Position"); f.IsValid() && f.CanInt() {
		return int(f.Int())
	}
	panic("fuzz: " + v.Type().String() + " has no int Position field")
}

// rangeFields names the pairs of int fields that hold a node's byte range:
// Loc{Start, End} in most engines, ClickHouse's span{Pos, End}, and the
// Elasticsearch console's ParsedRequest{StartOffset, EndOffset}.
var rangeFields = [][2]string{{"Start", "End"}, {"Pos", "End"}, {"StartOffset", "EndOffset"}}

// Tree fails t when a value reachable from tree carries a byte range
// outside input. A range is a struct with one of rangeFields'
// pairs of int fields; each bound is -1 (unknown) or in [0, len(input)], and a
// known Start is not past a known End. The walk keeps its own stack and
// visits a shared pointer once.
func Tree(t testing.TB, tree any, input string) {
	t.Helper()
	seen := map[uintptr]bool{}
	stack := []reflect.Value{reflect.ValueOf(tree)}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		switch v.Kind() {
		case reflect.Pointer:
			if v.IsNil() || seen[v.Pointer()] {
				continue
			}
			seen[v.Pointer()] = true
			stack = append(stack, v.Elem())
		case reflect.Interface:
			if !v.IsNil() {
				stack = append(stack, v.Elem())
			}
		case reflect.Slice, reflect.Array:
			for i := 0; i < v.Len(); i++ {
				stack = append(stack, v.Index(i))
			}
		case reflect.Map:
			for it := v.MapRange(); it.Next(); {
				stack = append(stack, it.Key(), it.Value())
			}
		case reflect.Struct:
			checkRange(t, v, input)
			for i := 0; i < v.NumField(); i++ {
				stack = append(stack, v.Field(i))
			}
		}
	}
}

func checkRange(t testing.TB, v reflect.Value, input string) {
	t.Helper()
	n := len(input)
	for _, pair := range rangeFields {
		start, end := intField(v, pair[0]), intField(v, pair[1])
		if !start.IsValid() || !end.IsValid() {
			continue
		}
		s, e := int(start.Int()), int(end.Int())
		if !bound(s, n) || !bound(e, n) || (s >= 0 && e >= 0 && s > e) {
			t.Fatalf("Parse(%s): %s{%s: %d, %s: %d} outside the %d-byte input", quote(input), v.Type(), pair[0], s, pair[1], e, n)
		}
		return
	}
}

// intField returns v's own int field name, not one promoted from an
// embedded struct, which the walk checks on its own.
func intField(v reflect.Value, name string) reflect.Value {
	f, ok := v.Type().FieldByName(name)
	if !ok || len(f.Index) != 1 {
		return reflect.Value{}
	}
	if fv := v.FieldByIndex(f.Index); fv.CanInt() {
		return fv
	}
	return reflect.Value{}
}

func bound(x, n int) bool { return x == -1 || (x >= 0 && x <= n) }

// quote shows up to 200 bytes of an input in a failure message.
func quote(input string) string {
	if len(input) > 200 {
		return strconv.Quote(input[:200]) + "..."
	}
	return strconv.Quote(input)
}
