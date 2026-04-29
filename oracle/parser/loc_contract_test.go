package parser

import (
	"reflect"
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

func TestOracleLocSentinelContract(t *testing.T) {
	if got := ast.NoLoc(); got != (ast.Loc{Start: -1, End: -1}) {
		t.Fatalf("NoLoc() = %+v, want {-1, -1}", got)
	}
	if !ast.NoLoc().IsUnknown() {
		t.Fatal("NoLoc().IsUnknown() = false, want true")
	}
	if (ast.Loc{Start: 0, End: 0}).IsUnknown() {
		t.Fatal("Loc{0,0}.IsUnknown() = true, want false")
	}
	if (ast.Loc{Start: -1, End: 0}).IsUnknown() {
		t.Fatal("mixed sentinel IsUnknown() = true, want false")
	}
}

func TestOracleLocWalkerRejectsInvalidSentinels(t *testing.T) {
	tests := []struct {
		name string
		loc  ast.Loc
		want string
	}{
		{name: "zero span", loc: ast.Loc{Start: 0, End: 0}, want: "Start >= 0 but End <= Start"},
		{name: "mixed start", loc: ast.Loc{Start: -1, End: 1}, want: "mixed unknown Loc sentinel"},
		{name: "mixed end", loc: ast.Loc{Start: 1, End: -1}, want: "mixed unknown Loc sentinel"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var violations []LocViolation
			walkNodeLocs(reflect.ValueOf(&ast.NumberLiteral{Val: "1", Loc: tt.loc}), "node", &violations)
			if len(violations) != 1 {
				t.Fatalf("violations = %d, want 1", len(violations))
			}
			if violations[0].Reason != tt.want {
				t.Fatalf("violation reason = %q, want %q", violations[0].Reason, tt.want)
			}
		})
	}
}

func TestOracleLocContractScenarios(t *testing.T) {
	cases := []string{
		`SELECT "MixedCase" AS alias FROM "TableName" WHERE a(+) = b`,
		"INSERT INTO t (a, b) VALUES (1, JSON_VALUE(payload, '$.x'))",
		"UPDATE t SET a = CASE WHEN b > 0 THEN b ELSE 0 END WHERE id = :id",
		"DELETE FROM t WHERE EXISTS (SELECT 1 FROM u WHERE u.id = t.id)",
		"CREATE TABLE t (id NUMBER GENERATED ALWAYS AS IDENTITY, name VARCHAR2(30))",
		"ALTER TABLE t ADD (created_at DATE DEFAULT SYSDATE)",
		"BEGIN x := func(1, 2); IF x > 0 THEN NULL; END IF; END;",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			if violations := CheckLocations(t, sql); len(violations) > 0 {
				t.Fatalf("Loc violations for %q: %v", sql, violations)
			}
		})
	}
}
