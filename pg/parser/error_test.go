package parser

import (
	"strings"
	"testing"
)

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name         string
		sql          string
		wantContains string
		wantPos      int // 0-based byte offset, -1 = end of input
	}{
		// Basic syntax errors
		{
			name:         "missing keyword",
			sql:          "SELECT * FRO t;",
			wantContains: `syntax error at or near "FRO"`,
			wantPos:      9,
		},
		{
			name:         "unexpected token",
			sql:          "CREATE TABLE (id int);",
			wantContains: `syntax error at or near "("`,
			wantPos:      13,
		},
		// Numeric table name
		{
			name:         "numeric table name",
			sql:          "CREATE TABLE 123 (id int);",
			wantContains: `syntax error at or near "123"`,
			wantPos:      13,
		},
		// Multiline — position is byte offset, not line/col
		{
			name:         "multiline error",
			sql:          "SELECT *\nFROM t\nWHER a = 1;",
			wantContains: `syntax error at or near "a"`,
			wantPos:      21,
		},
		// Error() format: pgx-style single line
		{
			name:         "pgx format",
			sql:          "SELECT * FRO t;",
			wantContains: "ERROR:",
		},
		{
			name:         "sqlstate in output",
			sql:          "SELECT * FRO t;",
			wantContains: "SQLSTATE 42601",
		},
		// Section 1.1: Logical & Comparison Operators — soft-fail nil checks
		{name: "OR no right operand", sql: "SELECT 1 OR", wantContains: "syntax error", wantPos: -1},
		{name: "AND no right operand", sql: "SELECT 1 AND", wantContains: "syntax error", wantPos: -1},
		{name: "less-than no right operand", sql: "SELECT 1 <", wantContains: "syntax error", wantPos: -1},
		{name: "greater-than no right operand", sql: "SELECT 1 >", wantContains: "syntax error", wantPos: -1},
		{name: "equals no right operand", sql: "SELECT 1 =", wantContains: "syntax error", wantPos: -1},
		{name: "less-equals no right operand", sql: "SELECT 1 <=", wantContains: "syntax error", wantPos: -1},
		{name: "greater-equals no right operand", sql: "SELECT 1 >=", wantContains: "syntax error", wantPos: -1},
		{name: "not-equals no right operand", sql: "SELECT 1 <>", wantContains: "syntax error", wantPos: -1},
		{name: "concat-op no right operand", sql: "SELECT 1 ||", wantContains: "syntax error", wantPos: -1},
		// Section 1.2: Arithmetic Operators — soft-fail nil checks
		{name: "plus no right operand", sql: "SELECT 1 +", wantContains: "syntax error", wantPos: -1},
		{name: "minus no right operand", sql: "SELECT 1 -", wantContains: "syntax error", wantPos: -1},
		{name: "multiply no right operand", sql: "SELECT 1 *", wantContains: "syntax error", wantPos: -1},
		{name: "divide no right operand", sql: "SELECT 1 /", wantContains: "syntax error", wantPos: -1},
		{name: "modulo no right operand", sql: "SELECT 1 %", wantContains: "syntax error", wantPos: -1},
		{name: "exponent no right operand", sql: "SELECT 1 ^", wantContains: "syntax error", wantPos: -1},
		// Section 1.3: IS DISTINCT FROM — soft-fail nil checks
		{name: "IS DISTINCT FROM no right expr", sql: "SELECT 1 IS DISTINCT FROM", wantContains: "syntax error", wantPos: -1},
		{name: "IS NOT DISTINCT FROM no right expr", sql: "SELECT 1 IS NOT DISTINCT FROM", wantContains: "syntax error", wantPos: -1},
		// Section 2.1: BETWEEN / LIKE / ILIKE / SIMILAR TO — soft-fail nil checks
		{name: "BETWEEN no lower bound", sql: "SELECT 1 BETWEEN", wantContains: "syntax error", wantPos: -1},
		{name: "BETWEEN AND no upper bound", sql: "SELECT 1 BETWEEN 0 AND", wantContains: "syntax error", wantPos: -1},
		{name: "NOT BETWEEN no lower bound", sql: "SELECT 1 NOT BETWEEN", wantContains: "syntax error", wantPos: -1},
		{name: "LIKE no pattern", sql: "SELECT 'a' LIKE", wantContains: "syntax error", wantPos: -1},
		{name: "LIKE ESCAPE no escape char", sql: "SELECT 'a' LIKE 'b' ESCAPE", wantContains: "syntax error", wantPos: -1},
		{name: "NOT LIKE no pattern", sql: "SELECT 'a' NOT LIKE", wantContains: "syntax error", wantPos: -1},
		{name: "ILIKE no pattern", sql: "SELECT 'a' ILIKE", wantContains: "syntax error", wantPos: -1},
		{name: "ILIKE ESCAPE no escape char", sql: "SELECT 'a' ILIKE 'b' ESCAPE", wantContains: "syntax error", wantPos: -1},
		{name: "SIMILAR TO no pattern", sql: "SELECT 'a' SIMILAR TO", wantContains: "syntax error", wantPos: -1},
		{name: "SIMILAR TO ESCAPE no escape char", sql: "SELECT 'a' SIMILAR TO 'b' ESCAPE", wantContains: "syntax error", wantPos: -1},
		// Section 2.2: COLLATE, TYPECAST & String Functions — soft-fail nil checks
		{name: "COLLATE no collation name", sql: "SELECT 'a' COLLATE", wantContains: "syntax error", wantPos: -1},
		{name: "TYPECAST no type name", sql: "SELECT 1::", wantContains: "syntax error", wantPos: -1},
		{name: "AT TIME ZONE no timezone expr", sql: "SELECT now() AT TIME ZONE", wantContains: "syntax error", wantPos: -1},
		{name: "OVERLAY PLACING no replacement", sql: "SELECT OVERLAY('abc' PLACING", wantContains: "syntax error", wantPos: -1},
		{name: "OVERLAY FOR no length", sql: "SELECT OVERLAY('abc' PLACING 'x' FROM 1 FOR", wantContains: "syntax error", wantPos: -1},
		{name: "POSITION IN no string", sql: "SELECT POSITION('a' IN", wantContains: "syntax error", wantPos: -1},
		{name: "SUBSTRING FROM no start", sql: "SELECT SUBSTRING('abc' FROM", wantContains: "syntax error", wantPos: -1},
		{name: "SUBSTRING FROM FOR no length", sql: "SELECT SUBSTRING('abc' FROM 1 FOR", wantContains: "syntax error", wantPos: -1},
		{name: "SUBSTRING SIMILAR no pattern", sql: "SELECT SUBSTRING('abc' SIMILAR", wantContains: "syntax error", wantPos: -1},
		{name: "TRIM FROM no source string", sql: "SELECT TRIM(LEADING FROM", wantContains: "syntax error", wantPos: -1},
		// Section 3.1: b_expr Binary Operators — soft-fail nil checks
		{name: "b_expr plus no right", sql: "SELECT CAST(1 + AS int)", wantContains: `syntax error at or near "AS"`, wantPos: 16},
		{name: "b_expr minus no right", sql: "SELECT CAST(1 - AS int)", wantContains: `syntax error at or near "AS"`, wantPos: 16},
		{name: "b_expr multiply no right", sql: "SELECT CAST(1 * AS int)", wantContains: `syntax error at or near "AS"`, wantPos: 16},
		{name: "b_expr less-than no right", sql: "SELECT CAST(1 < AS int)", wantContains: `syntax error at or near "AS"`, wantPos: 16},
		{name: "b_expr greater-than no right", sql: "SELECT CAST(1 > AS int)", wantContains: `syntax error at or near "AS"`, wantPos: 16},
		{name: "b_expr equals no right", sql: "SELECT CAST(1 = AS int)", wantContains: `syntax error at or near "AS"`, wantPos: 16},
		{name: "b_expr less-equals no right", sql: "SELECT CAST(1 <= AS int)", wantContains: `syntax error at or near "AS"`, wantPos: 17},
		{name: "b_expr greater-equals no right", sql: "SELECT CAST(1 >= AS int)", wantContains: `syntax error at or near "AS"`, wantPos: 17},
		{name: "b_expr not-equals no right", sql: "SELECT CAST(1 <> AS int)", wantContains: `syntax error at or near "AS"`, wantPos: 17},
		{name: "b_expr concat-op no right", sql: "SELECT CAST(1 || AS int)", wantContains: `syntax error at or near "AS"`, wantPos: 17},
		{name: "b_expr IS DISTINCT FROM no right", sql: "SELECT CAST(1 IS DISTINCT FROM AS int)", wantContains: `syntax error at or near "AS"`, wantPos: 31},
		{name: "b_expr TYPECAST no type", sql: "SELECT CAST(1:: AS int)", wantContains: `syntax error at or near "AS"`, wantPos: 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.sql)
			if err == nil {
				t.Fatalf("Parse(%q) expected error, got nil", tt.sql)
			}

			if !strings.Contains(err.Error(), tt.wantContains) {
				t.Errorf("Parse(%q) error = %q, want substring %q", tt.sql, err.Error(), tt.wantContains)
			}

			if tt.wantPos != 0 { // 0 means "don't check position"
				pe, ok := err.(*ParseError)
				if !ok {
					t.Fatalf("error is not *ParseError: %T", err)
				}
				if tt.wantPos == -1 {
					// For end-of-input, the EOF token has Loc=len(sql) or similar;
					// just verify Position >= len(strings.TrimRight(tt.sql, " \t\n;"))
					trimmed := strings.TrimRight(tt.sql, " \t\n;")
					if pe.Position < len(trimmed) {
						t.Errorf("Position = %d, want >= %d (end of input)", pe.Position, len(trimmed))
					}
				} else if pe.Position != tt.wantPos {
					t.Errorf("Position = %d, want %d", pe.Position, tt.wantPos)
				}
			}
		})
	}
}

func TestParseErrors_NameType(t *testing.T) {
	tests := []struct {
		name         string
		sql          string
		wantContains string
		wantPos      int
	}{
		{
			name:         "invalid table name",
			sql:          "CREATE TABLE 123 (id int);",
			wantContains: `syntax error at or near "123"`,
			wantPos:      13,
		},
		{
			name:         "invalid array bound",
			sql:          "CREATE TABLE t (id int[x]);",
			wantContains: `syntax error at or near "x"`,
			wantPos:      23,
		},
		{
			name:         "invalid column type",
			sql:          "CREATE TABLE t (id INT XYZ);",
			wantContains: `syntax error at or near "XYZ"`,
			wantPos:      23,
		},
		{
			name:         "invalid national type",
			sql:          "CREATE TABLE t (id NATIONAL 123);",
			wantContains: `syntax error at or near "123"`,
			wantPos:      28,
		},
		{
			name:         "leading garbage",
			sql:          ") SELECT 1;",
			wantContains: `syntax error at or near ")"`,
			wantPos:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.sql)
			if err == nil {
				t.Fatalf("Parse(%q) expected error, got nil", tt.sql)
			}

			if !strings.Contains(err.Error(), tt.wantContains) {
				t.Errorf("Parse(%q) error = %q, want substring %q", tt.sql, err.Error(), tt.wantContains)
			}

			if tt.wantPos > 0 {
				pe, ok := err.(*ParseError)
				if !ok {
					t.Fatalf("error is not *ParseError: %T", err)
				}
				if pe.Position != tt.wantPos {
					t.Errorf("Position = %d, want %d", pe.Position, tt.wantPos)
				}
			}
		})
	}
}
