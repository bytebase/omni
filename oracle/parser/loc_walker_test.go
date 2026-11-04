package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// LocViolation records a single location validation failure.
type LocViolation struct {
	Path    string // e.g. "Items[0](SelectStmt).TargetList[0](ResTarget)"
	NodeTag string // node type name
	Start   int
	End     int
	Reason  string
}

func (v LocViolation) String() string {
	return fmt.Sprintf("%s [%s]: Start=%d End=%d — %s", v.Path, v.NodeTag, v.Start, v.End, v.Reason)
}

// CheckLocations parses sql via Parse(), recursively walks the AST using
// reflection, and returns all Loc violations walkNodeLocs finds.
// This does NOT call t.Fatal — callers decide how to handle violations.
func CheckLocations(t *testing.T, sql string) []LocViolation {
	t.Helper()

	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("CheckLocations Parse(%q): %v", sql, err)
	}

	var violations []LocViolation
	if result != nil {
		for i, item := range result.Items {
			path := fmt.Sprintf("Items[%d]", i)
			walkNodeLocs(reflect.ValueOf(item), path, &violations)
		}
	}
	return violations
}

// locParent is the nearest node above a walked value whose Loc is known.
type locParent struct {
	tag string
	loc ast.Loc
}

// walkNodeLocs recursively walks a reflected AST value, checking every Loc
// field: an unknown Loc is the single {-1, -1} sentinel, a known one is not
// empty, and a known one lies within the known Loc of the nearest node
// above it. The last catches a node built around an operand the parser
// read first, a binary operator, a postfix test, a join, that starts at
// its own keyword instead of where that operand starts.
func walkNodeLocs(v reflect.Value, path string, violations *[]LocViolation) {
	walkNodeLocsIn(v, path, locParent{loc: ast.NoLoc()}, violations)
}

func walkNodeLocsIn(v reflect.Value, path string, parent locParent, violations *[]LocViolation) {
	// Dereference pointers and interfaces.
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		typeName := v.Type().Name()

		// Check if this struct has a Loc field of type ast.Loc.
		locField := v.FieldByName("Loc")
		if locField.IsValid() && locField.Type() == reflect.TypeOf(ast.Loc{}) {
			loc := locField.Interface().(ast.Loc)
			switch {
			case loc.IsUnknown():
				// Unknown locations must use the single {-1, -1} sentinel.
			case loc.Start < 0 || loc.End < 0:
				*violations = append(*violations, LocViolation{
					Path:    path,
					NodeTag: typeName,
					Start:   loc.Start,
					End:     loc.End,
					Reason:  "mixed unknown Loc sentinel",
				})
			case loc.End <= loc.Start:
				*violations = append(*violations, LocViolation{
					Path:    path,
					NodeTag: typeName,
					Start:   loc.Start,
					End:     loc.End,
					Reason:  "Start >= 0 but End <= Start",
				})
			default:
				if !parent.loc.IsUnknown() && (loc.Start < parent.loc.Start || loc.End > parent.loc.End) {
					*violations = append(*violations, LocViolation{
						Path:    path,
						NodeTag: typeName,
						Start:   loc.Start,
						End:     loc.End,
						Reason:  fmt.Sprintf("outside %s [%d,%d)", parent.tag, parent.loc.Start, parent.loc.End),
					})
				}
				parent = locParent{tag: typeName, loc: loc}
			}
		}

		// Recurse into all fields.
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			if !field.IsExported() {
				continue
			}
			if field.Name == "Loc" {
				continue // already checked
			}
			childPath := path + "." + field.Name
			within := parent
			if typeName == "InsertIntoClause" && field.Name == "When" {
				// A multitable INSERT's WHEN condition comes before the
				// INTO clauses of its branch, and each of them holds it.
				within = locParent{loc: ast.NoLoc()}
			}
			walkNodeLocsIn(v.Field(i), childPath, within, violations)
		}

	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			elemPath := fmt.Sprintf("%s[%d]", path, i)
			// Add type name for interface elements.
			actual := elem
			for actual.Kind() == reflect.Ptr || actual.Kind() == reflect.Interface {
				if actual.IsNil() {
					break
				}
				actual = actual.Elem()
			}
			if actual.IsValid() && actual.Kind() == reflect.Struct {
				elemPath = fmt.Sprintf("%s[%d](%s)", path, i, actual.Type().Name())
			}
			walkNodeLocsIn(elem, elemPath, parent, violations)
		}
	}
}
