package parser

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/oracle/ast"
)

// TestInfixLocStartsAtFirstOperand pins the text every node built around
// an operand the parser read first covers: a binary operator, AND and OR,
// a postfix test, a quantified comparison, MULTISET, a member, a cursor
// attribute, AT TIME ZONE, a subscript, a choice range, a join, PIVOT,
// UNPIVOT, and MATCH_RECOGNIZE. Each starts where its first operand does,
// a '(' when the operand is parenthesized, and ends at its last token.
// The operands are compound, so a start at the operator shows.
func TestInfixLocStartsAtFirstOperand(t *testing.T) {
	tests := []struct {
		sql  string
		node string   // the node type
		want []string // the text of every node of that type, in tree order
	}{
		{"SELECT * FROM t WHERE a*b = c", "BinaryExpr", []string{"a*b = c", "a*b"}},
		{"SELECT * FROM t WHERE (a) = b", "BinaryExpr", []string{"(a) = b"}},
		{"SELECT 1 + 2 * 3 - 4 FROM dual", "BinaryExpr", []string{"1 + 2 * 3 - 4", "1 + 2 * 3", "2 * 3"}},
		{"SELECT 2 ** 3 ** 2 FROM dual", "BinaryExpr", []string{"2 ** 3 ** 2", "3 ** 2"}},
		{"SELECT * FROM t WHERE x || y LIKE z", "LikeExpr", []string{"x || y LIKE z"}},
		{"SELECT * FROM t WHERE x || y LIKE z", "BinaryExpr", []string{"x || y"}},
		{"SELECT * FROM t WHERE x || y NOT LIKE z ESCAPE '!'", "LikeExpr", []string{"x || y NOT LIKE z ESCAPE '!'"}},
		{"SELECT * FROM t WHERE a IS NULL AND b IN (1)", "BoolExpr", []string{"a IS NULL AND b IN (1)"}},
		{"SELECT * FROM t WHERE a IS NULL AND b IN (1)", "IsExpr", []string{"a IS NULL"}},
		{"SELECT * FROM t WHERE a IS NULL AND b IN (1)", "InExpr", []string{"b IN (1)"}},
		{"SELECT * FROM t WHERE a + 1 IS NOT NULL", "IsExpr", []string{"a + 1 IS NOT NULL"}},
		{"SELECT * FROM t WHERE a = 1 AND b = 2 AND c = 3 OR d = 4", "BoolExpr", []string{"a = 1 AND b = 2 AND c = 3 OR d = 4", "a = 1 AND b = 2 AND c = 3"}},
		{"SELECT * FROM t WHERE NOT (a = 1 OR b = 2)", "BoolExpr", []string{"NOT (a = 1 OR b = 2)", "a = 1 OR b = 2"}},
		{"SELECT * FROM t WHERE a + b BETWEEN c - 1 AND d", "BetweenExpr", []string{"a + b BETWEEN c - 1 AND d"}},
		{"SELECT * FROM t WHERE a - b NOT BETWEEN 1 AND 2", "BetweenExpr", []string{"a - b NOT BETWEEN 1 AND 2"}},
		{"SELECT * FROM t WHERE a || b NOT IN (SELECT c FROM u)", "InExpr", []string{"a || b NOT IN (SELECT c FROM u)"}},
		{"SELECT * FROM t WHERE a || b NOT IN (SELECT c FROM u)", "SubqueryExpr", []string{"SELECT c FROM u"}},
		{"SELECT * FROM t WHERE a * 2 > ANY (SELECT c FROM u)", "BinaryExpr", []string{"a * 2 > ANY (SELECT c FROM u)", "a * 2"}},
		{"SELECT * FROM t WHERE a * 2 > ANY (SELECT c FROM u)", "SubqueryExpr", []string{"SELECT c FROM u"}},
		{"SELECT * FROM t WHERE a - 1 = ALL (1, 2)", "BinaryExpr", []string{"a - 1 = ALL (1, 2)", "a - 1"}},
		{"SELECT n1 MULTISET UNION ALL n2 FROM t", "MultisetExpr", []string{"n1 MULTISET UNION ALL n2"}},
		{"SELECT CAST(d AS TIMESTAMP) AT TIME ZONE 'UTC' FROM t", "FuncCallExpr", []string{"CAST(d AS TIMESTAMP) AT TIME ZONE 'UTC'"}},
		{"SELECT * FROM t MODEL DIMENSION BY (d) MEASURES (s) RULES (s[1] = s[0] * 2)", "BinaryExpr", []string{"s[0] * 2"}},
		{"BEGIN v := rows(i).status || 'x'; END;", "BinaryExpr", []string{"rows(i).status || 'x'"}},
		{"BEGIN v := rows(i).status || 'x'; END;", "FieldAccessExpr", []string{"rows(i).status"}},
		{"BEGIN IF c%ROWCOUNT + 1 > 2 THEN NULL; END IF; END;", "BinaryExpr", []string{"c%ROWCOUNT + 1 > 2", "c%ROWCOUNT + 1"}},
		{"DECLARE x t := t(1 + 1..5 => 2); BEGIN NULL; END;", "ChoiceRange", []string{"1 + 1..5"}},
		{"DECLARE x t := t(1 + 1..5 => 2); BEGIN NULL; END;", "ChoiceAssocExpr", []string{"1 + 1..5 => 2"}},
		{"SELECT * FROM t a JOIN u b ON a.id = b.id LEFT JOIN v ON v.id = a.id", "JoinClause", []string{"t a JOIN u b ON a.id = b.id LEFT JOIN v ON v.id = a.id", "t a JOIN u b ON a.id = b.id"}},
		{"SELECT * FROM t a JOIN u b ON a.id = b.id LEFT JOIN v ON v.id = a.id", "BinaryExpr", []string{"a.id = b.id", "v.id = a.id"}},
		{"SELECT * FROM t PIVOT (SUM(a) FOR b IN (1)) p", "PivotClause", []string{"t PIVOT (SUM(a) FOR b IN (1)) p"}},
		{"SELECT * FROM (SELECT 1 a FROM dual) UNPIVOT (v FOR k IN (a))", "UnpivotClause", []string{"(SELECT 1 a FROM dual) UNPIVOT (v FOR k IN (a))"}},
		{"SELECT * FROM t MATCH_RECOGNIZE (ORDER BY id PATTERN (a) DEFINE a AS a.n > 0)", "MatchRecognizeClause", []string{"t MATCH_RECOGNIZE (ORDER BY id PATTERN (a) DEFINE a AS a.n > 0)"}},
		{"SELECT * FROM t MATCH_RECOGNIZE (ORDER BY id PATTERN (a) DEFINE a AS a.n > 0)", "TableRef", []string{"t"}},
	}
	for _, tt := range tests {
		t.Run(tt.node+"/"+tt.sql, func(t *testing.T) {
			tree, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			var got []string
			for _, loc := range locsOfType(tree, tt.node) {
				if loc.Start < 0 || loc.End > len(tt.sql) || loc.End < loc.Start {
					t.Fatalf("%s Loc = %+v, outside the input", tt.node, loc)
				}
				got = append(got, tt.sql[loc.Start:loc.End])
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("%s texts mismatch (-want +got):\n%s", tt.node, diff)
			}
			if violations := CheckLocations(t, tt.sql); len(violations) > 0 {
				t.Errorf("Loc violations: %v", violations)
			}
		})
	}
}

// TestOracleLocWalkerRejectsChildOutsideParent checks that the Loc walker
// reports a node that starts after its first operand, as a binary
// expression starting at its operator would.
func TestOracleLocWalkerRejectsChildOutsideParent(t *testing.T) {
	// a = 1, its Loc starting at the '='.
	expr := &ast.BinaryExpr{
		Op:    "=",
		Left:  &ast.ColumnRef{Column: "A", Loc: ast.Loc{Start: 0, End: 1}},
		Right: &ast.NumberLiteral{Val: "1", Loc: ast.Loc{Start: 4, End: 5}},
		Loc:   ast.Loc{Start: 2, End: 5},
	}
	var violations []LocViolation
	walkNodeLocs(reflect.ValueOf(expr), "node", &violations)
	if len(violations) != 1 {
		t.Fatalf("violations = %v, want one", violations)
	}
	if v := violations[0]; v.Path != "node.Left" || v.Reason != "outside BinaryExpr [2,5)" {
		t.Fatalf("violation = %s, want node.Left outside BinaryExpr [2,5)", v)
	}
}

// locsOfType returns the Loc of every node of the named type in n, in
// tree order, each node once.
func locsOfType(n ast.Node, typeName string) []ast.Loc {
	type key struct {
		typ reflect.Type
		ptr uintptr
	}
	var out []ast.Loc
	seen := map[key]bool{}
	var walk func(v reflect.Value)
	walk = func(v reflect.Value) {
		for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
			if v.IsNil() {
				return
			}
			if v.Kind() == reflect.Ptr {
				k := key{v.Type(), v.Pointer()}
				if seen[k] {
					return
				}
				seen[k] = true
			}
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Struct:
			if v.Type().Name() == typeName {
				if f := v.FieldByName("Loc"); f.IsValid() && f.Type() == reflect.TypeOf(ast.Loc{}) {
					out = append(out, f.Interface().(ast.Loc))
				}
			}
			for i := 0; i < v.NumField(); i++ {
				if v.Type().Field(i).IsExported() {
					walk(v.Field(i))
				}
			}
		case reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				walk(v.Index(i))
			}
		}
	}
	walk(reflect.ValueOf(n))
	return out
}
