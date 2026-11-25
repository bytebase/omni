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

// rangeTypes names the struct types that hold a node's byte range and
// their bound fields: Loc{Start, End} in every engine but ClickHouse, whose
// nodes embed span{Pos, End}. Other fields called Start and End can hold
// anything (MySQL's VCPUSpec holds CPU numbers), so only these count.
var rangeTypes = map[string][2]string{"Loc": {"Start", "End"}, "span": {"Pos", "End"}}

// Tree fails t when a Loc (or span) reachable from tree lies outside input
// (see Range). The walk keeps its own stack and visits a shared pointer
// once.
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
			if pair, ok := rangeTypes[v.Type().Name()]; ok {
				start, end := v.FieldByName(pair[0]), v.FieldByName(pair[1])
				if start.Kind() == reflect.Int && end.Kind() == reflect.Int {
					Range(t, input, v.Type().String(), int(start.Int()), int(end.Int()))
				}
			}
			for i := 0; i < v.NumField(); i++ {
				stack = append(stack, v.Field(i))
			}
		}
	}
}

// Range fails t unless [start, end) is a byte range of input: each bound
// is -1 (unknown) or in [0, len(input)], and a known start is not past a
// known end. what names the range in the failure.
func Range(t testing.TB, input, what string, start, end int) {
	t.Helper()
	n := len(input)
	if !bound(start, n) || !bound(end, n) || (start >= 0 && end >= 0 && start > end) {
		t.Fatalf("Parse(%s): %s [%d, %d) outside the %d-byte input", quote(input), what, start, end, n)
	}
}

func bound(x, n int) bool { return x == -1 || (x >= 0 && x <= n) }

// quote shows up to 200 bytes of an input in a failure message.
func quote(input string) string {
	if len(input) > 200 {
		return strconv.Quote(input[:200]) + "..."
	}
	return strconv.Quote(input)
}
