package analysis

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/oracle/parser"
)

func TestParse(t *testing.T) {
	list, err := Parse("SELECT 1 FROM DUAL;\nSELECT A FROM T")
	require.NoError(t, err)
	require.Len(t, list.Items, 2)
	second, ok := list.Items[1].(*ast.RawStmt)
	require.True(t, ok)
	require.Equal(t, 20, second.Stmt.(*ast.SelectStmt).Loc.Start)

	list, err = Parse("-- only a comment")
	require.NoError(t, err)
	require.Empty(t, list.Items)

	_, err = Parse("SELECT 1 FROM DUAL;\nSELEC 2")
	var parseErr *parser.ParseError
	require.ErrorAs(t, err, &parseErr)
	require.GreaterOrEqual(t, parseErr.Position, 20)
}

func TestParseLocatesEveryNodeInTheScript(t *testing.T) {
	// Each statement's own nodes, reached by ast.Walk or not, lie inside it:
	// a DELETE's target, an UNPIVOT's value column, TABLE()'s argument, a
	// multi-table INSERT's WHEN condition.
	for _, statement := range []string{
		"DELETE a WHERE b = 1",
		"SELECT * FROM a UNPIVOT (v FOR k IN (c))",
		"SELECT * FROM TABLE(1)",
		"INSERT WHEN 1 = 1 THEN INTO a INTO b SELECT 1 FROM dual",
	} {
		script := "SELECT 1 FROM DUAL;\n" + statement
		list, err := Parse(script)
		require.NoError(t, err, statement)
		require.Len(t, list.Items, 2, statement)
		start := len("SELECT 1 FROM DUAL;\n")
		for _, loc := range locsOf(reflect.ValueOf(list.Items[1]), map[uintptr]bool{}) {
			if loc.Start < 0 {
				continue
			}
			require.GreaterOrEqual(t, loc.Start, start, statement)
			require.LessOrEqual(t, loc.End, len(script), statement)
		}
	}
}

// locsOf returns every Loc in a tree, following each field, not only those
// ast.Walk visits.
func locsOf(v reflect.Value, seen map[uintptr]bool) []ast.Loc {
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		if v.IsNil() {
			return nil
		}
		if v.Kind() == reflect.Pointer {
			if seen[v.Pointer()] {
				return nil
			}
			seen[v.Pointer()] = true
		}
		return locsOf(v.Elem(), seen)
	case reflect.Struct:
		if loc, ok := v.Interface().(ast.Loc); ok {
			return []ast.Loc{loc}
		}
		var locs []ast.Loc
		for i := 0; i < v.NumField(); i++ {
			if v.Type().Field(i).IsExported() {
				locs = append(locs, locsOf(v.Field(i), seen)...)
			}
		}
		return locs
	case reflect.Slice:
		var locs []ast.Loc
		for i := 0; i < v.Len(); i++ {
			locs = append(locs, locsOf(v.Index(i), seen)...)
		}
		return locs
	default:
		return nil
	}
}
