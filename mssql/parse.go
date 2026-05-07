package mssql

import (
	"reflect"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/mssql/parser"
)

// Statement is the result of parsing a single SQL statement.
type Statement struct {
	// Text is the SQL text including trailing semicolon if present.
	Text string
	// AST is the inner statement node. Nil for empty statements.
	AST ast.Node

	// ByteStart is the inclusive start byte offset in the original SQL.
	ByteStart int
	// ByteEnd is the exclusive end byte offset in the original SQL.
	ByteEnd int

	// Start is the start position (line:column) in the original SQL.
	Start Position
	// End is the exclusive end position (line:column) in the original SQL.
	End Position
}

// Position represents a location in source text.
type Position struct {
	// Line is 1-based line number.
	Line int
	// Column is 1-based column in bytes.
	Column int
}

// Empty returns true if this statement has no meaningful content.
func (s *Statement) Empty() bool {
	return s.AST == nil
}

// Parse splits and parses a SQL string into statements.
// Each statement includes the text, AST, and byte/line positions.
func Parse(sql string) ([]Statement, error) {
	list, err := parser.Parse(sql)
	if err != nil {
		return nil, err
	}
	if list == nil || len(list.Items) == 0 {
		return nil, nil
	}

	idx := buildLineIndex(sql)

	var stmts []Statement
	prevEnd := 0

	for _, item := range list.Items {
		if item == nil {
			continue
		}

		loc := nodeLoc(item)

		// Text starts where the previous statement ended.
		start := prevEnd

		// Text ends after the semicolon following the statement, or at the Loc.End.
		end := loc.End
		if end <= start {
			end = len(sql)
		}
		// Scan past trailing whitespace to find the semicolon.
		j := end
		for j < len(sql) && isSpace(sql[j]) {
			j++
		}
		if j < len(sql) && sql[j] == ';' {
			end = j + 1
		}

		// Start position points to the AST node start, while Text preserves any
		// preceding separators that belong to this statement segment.
		contentStart := loc.Start
		if contentStart < start || contentStart > end {
			contentStart = start
			for contentStart < end && isSpace(sql[contentStart]) {
				contentStart++
			}
		}

		stmts = append(stmts, Statement{
			Text:      sql[start:end],
			AST:       item,
			ByteStart: start,
			ByteEnd:   end,
			Start:     offsetToPosition(idx, contentStart),
			End:       offsetToPosition(idx, end),
		})

		prevEnd = end
	}
	return stmts, nil
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// nodeLoc extracts the Loc field from an AST node using reflection.
func nodeLoc(n ast.Node) ast.Loc {
	v := reflect.ValueOf(n)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		f := v.FieldByName("Loc")
		if f.IsValid() && f.Type() == reflect.TypeOf(ast.Loc{}) {
			return f.Interface().(ast.Loc)
		}
	}
	return ast.NoLoc()
}

type lineIndex []int

func buildLineIndex(s string) lineIndex {
	idx := lineIndex{0}
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			idx = append(idx, i+1)
		}
	}
	return idx
}

func offsetToPosition(idx lineIndex, offset int) Position {
	// Binary search for the line containing offset.
	lo, hi := 0, len(idx)-1
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if idx[mid] <= offset {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return Position{
		Line:   lo + 1,               // 1-based
		Column: offset - idx[lo] + 1, // 1-based
	}
}
