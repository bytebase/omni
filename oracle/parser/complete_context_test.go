package parser

import (
	"strings"
	"testing"
)

func TestOracleCompletionSelectTableReferenceSignals(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		qualifier MultipartName
	}{
		{name: "from", input: "SELECT * FROM |"},
		{name: "join", input: "SELECT * FROM t JOIN |"},
		{name: "schema qualified", input: "SELECT * FROM schema1.|", qualifier: MultipartName{Schema: "schema1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := collectCompletionMarked(t, tt.input)
			requireCompletionRule(t, ctx.Candidates, "table_ref")
			requireObjectKind(t, ctx.Intent, ObjectKindTable)
			if tt.qualifier.Schema != "" && !strings.EqualFold(ctx.Intent.Qualifier.Schema, tt.qualifier.Schema) {
				t.Fatalf("schema qualifier = %q, want %q", ctx.Intent.Qualifier.Schema, tt.qualifier.Schema)
			}
		})
	}
}

func TestOracleCompletionSelectColumnReferenceSignals(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		qualifier MultipartName
		refs      []string
	}{
		{name: "target list", input: "SELECT | FROM employees", refs: []string{"employees"}},
		{name: "qualified alias", input: "SELECT e.| FROM employees e", qualifier: MultipartName{Object: "e"}, refs: []string{"e"}},
		{name: "join on", input: "SELECT * FROM employees e JOIN departments d ON |", refs: []string{"e", "d"}},
		{name: "where", input: "SELECT * FROM employees e WHERE |", refs: []string{"e"}},
		{name: "group by", input: "SELECT id FROM employees e GROUP BY |", refs: []string{"e"}},
		{name: "order by", input: "SELECT id AS eid FROM employees e ORDER BY |", refs: []string{"e"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := collectCompletionMarked(t, tt.input)
			requireCompletionRule(t, ctx.Candidates, "columnref")
			requireObjectKind(t, ctx.Intent, ObjectKindColumn)
			if tt.qualifier.Object != "" && !strings.EqualFold(ctx.Intent.Qualifier.Object, tt.qualifier.Object) {
				t.Fatalf("object qualifier = %q, want %q", ctx.Intent.Qualifier.Object, tt.qualifier.Object)
			}
			for _, ref := range tt.refs {
				requireRangeReference(t, ctx.Scope, ref)
			}
		})
	}
}

func TestOracleCompletionCTEAndSubqueryScope(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		ref     string
		columns []string
	}{
		{
			name:  "cte table reference",
			input: "WITH x AS (SELECT * FROM employees) SELECT * FROM |",
			ref:   "x",
		},
		{
			name:    "explicit cte columns",
			input:   "WITH x(x1, x2) AS (SELECT * FROM employees) SELECT x.| FROM x",
			ref:     "x",
			columns: []string{"x1", "x2"},
		},
		{
			name:    "subquery alias columns",
			input:   "SELECT src.| FROM (SELECT c1, c2 FROM employees) src",
			ref:     "src",
			columns: []string{"c1", "c2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := collectCompletionMarked(t, tt.input)
			requireRangeReference(t, ctx.Scope, tt.ref)
			if len(tt.columns) > 0 {
				ref := findRangeReference(t, ctx.Scope, tt.ref)
				for _, col := range tt.columns {
					requireColumnName(t, ref.Columns, col)
				}
			}
		})
	}
}

func TestOracleCompletionDMLSignals(t *testing.T) {
	tests := []struct {
		name  string
		input string
		rule  string
		kind  ObjectKind
		refs  []string
	}{
		{name: "insert table", input: "INSERT INTO |", rule: "table_ref", kind: ObjectKindTable},
		{name: "insert column list", input: "INSERT INTO employees (|)", rule: "columnref", kind: ObjectKindColumn, refs: []string{"employees"}},
		{name: "insert values expression", input: "INSERT INTO employees VALUES (|)", rule: "columnref", kind: ObjectKindColumn, refs: []string{"employees"}},
		{name: "insert select source", input: "INSERT INTO employees SELECT | FROM departments", rule: "columnref", kind: ObjectKindColumn, refs: []string{"departments"}},
		{name: "update table", input: "UPDATE | SET name = 'x'", rule: "table_ref", kind: ObjectKindTable},
		{name: "update set column", input: "UPDATE employees SET |", rule: "columnref", kind: ObjectKindColumn, refs: []string{"employees"}},
		{name: "update set expression", input: "UPDATE employees SET name = |", rule: "columnref", kind: ObjectKindColumn, refs: []string{"employees"}},
		{name: "update where", input: "UPDATE employees SET name = 'x' WHERE |", rule: "columnref", kind: ObjectKindColumn, refs: []string{"employees"}},
		{name: "delete table", input: "DELETE FROM |", rule: "table_ref", kind: ObjectKindTable},
		{name: "delete where", input: "DELETE FROM employees WHERE |", rule: "columnref", kind: ObjectKindColumn, refs: []string{"employees"}},
		{name: "merge target", input: "MERGE INTO |", rule: "table_ref", kind: ObjectKindTable},
		{name: "merge source", input: "MERGE INTO employees e USING |", rule: "table_ref", kind: ObjectKindTable, refs: []string{"e"}},
		{name: "merge on", input: "MERGE INTO employees e USING departments d ON |", rule: "columnref", kind: ObjectKindColumn, refs: []string{"e", "d"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := collectCompletionMarked(t, tt.input)
			requireCompletionRule(t, ctx.Candidates, tt.rule)
			requireObjectKind(t, ctx.Intent, tt.kind)
			for _, ref := range tt.refs {
				requireRangeReference(t, ctx.Scope, ref)
			}
		})
	}
}

func TestOracleCompletionDDLAndUtilitySignals(t *testing.T) {
	t.Run("create object types", func(t *testing.T) {
		ctx := collectCompletionMarked(t, "CREATE |")
		for _, tok := range []int{TABLE, VIEW, INDEX, SEQUENCE, SYNONYM, PROCEDURE, FUNCTION, PACKAGE, TRIGGER, USER, ROLE} {
			if !ctx.Candidates.HasToken(tok) {
				t.Fatalf("missing CREATE candidate %q; got %v", TokenName(tok), tokenNamesForTest(ctx.Candidates.Tokens))
			}
		}
	})

	tests := []struct {
		name  string
		input string
		rule  string
		kind  ObjectKind
		refs  []string
	}{
		{name: "create table datatype", input: "CREATE TABLE employees (name |)", rule: "type_name", kind: ObjectKindType},
		{name: "create table references", input: "CREATE TABLE employees (dept_id NUMBER REFERENCES |)", rule: "table_ref", kind: ObjectKindTable},
		{name: "alter table", input: "ALTER TABLE |", rule: "table_ref", kind: ObjectKindTable},
		{name: "alter drop column", input: "ALTER TABLE employees DROP COLUMN |", rule: "columnref", kind: ObjectKindColumn, refs: []string{"employees"}},
		{name: "drop table", input: "DROP TABLE |", rule: "table_ref", kind: ObjectKindTable},
		{name: "drop view", input: "DROP VIEW |", rule: "table_ref", kind: ObjectKindView},
		{name: "drop sequence", input: "DROP SEQUENCE |", rule: "sequence_ref", kind: ObjectKindSequence},
		{name: "truncate table", input: "TRUNCATE TABLE |", rule: "table_ref", kind: ObjectKindTable},
		{name: "comment table", input: "COMMENT ON TABLE |", rule: "table_ref", kind: ObjectKindTable},
		{name: "comment column", input: "COMMENT ON COLUMN employees.|", rule: "columnref", kind: ObjectKindColumn, refs: []string{"employees"}},
		{name: "grant on object", input: "GRANT SELECT ON |", rule: "table_ref", kind: ObjectKindTable},
		{name: "revoke on object", input: "REVOKE SELECT ON |", rule: "table_ref", kind: ObjectKindTable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := collectCompletionMarked(t, tt.input)
			requireCompletionRule(t, ctx.Candidates, tt.rule)
			requireObjectKind(t, ctx.Intent, tt.kind)
			for _, ref := range tt.refs {
				requireRangeReference(t, ctx.Scope, ref)
			}
		})
	}
}

func TestOracleCompletionAlterTableAddCandidates(t *testing.T) {
	ctx := collectCompletionMarked(t, "ALTER TABLE employees ADD |")
	for _, tok := range []int{kwCOLUMN, kwCONSTRAINT, kwPRIMARY, kwUNIQUE, kwFOREIGN, kwCHECK} {
		if !ctx.Candidates.HasToken(tok) {
			t.Fatalf("missing ALTER TABLE ADD candidate %q; got %v", TokenName(tok), tokenNamesForTest(ctx.Candidates.Tokens))
		}
	}
}

func TestOracleCompletionObjectKindSpecificDDL(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      ObjectKind
		notWanted ObjectKind
	}{
		{name: "drop table", input: "DROP TABLE |", want: ObjectKindTable, notWanted: ObjectKindView},
		{name: "drop view", input: "DROP VIEW |", want: ObjectKindView, notWanted: ObjectKindTable},
		{name: "comment table", input: "COMMENT ON TABLE |", want: ObjectKindTable, notWanted: ObjectKindSequence},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := collectCompletionMarked(t, tt.input)
			requireObjectKind(t, ctx.Intent, tt.want)
			if hasObjectKind(ctx.Intent, tt.notWanted) {
				t.Fatalf("intent kinds = %v, should not include %v", ctx.Intent.ObjectKinds, tt.notWanted)
			}
		})
	}
}

func TestOracleCompletionKeywordDrivenCandidates(t *testing.T) {
	t.Run("datatype context", func(t *testing.T) {
		ctx := collectCompletionMarked(t, "CREATE TABLE employees (name |)")
		for _, tok := range []int{kwNUMBER, kwVARCHAR2, kwDATE, kwTIMESTAMP} {
			if !ctx.Candidates.HasToken(tok) {
				t.Fatalf("missing datatype candidate %q; got %v", TokenName(tok), tokenNamesForTest(ctx.Candidates.Tokens))
			}
		}
	})

	t.Run("expression context", func(t *testing.T) {
		ctx := collectCompletionMarked(t, "SELECT | FROM dual")
		for _, tok := range []int{kwCAST, kwDECODE, kwJSON_VALUE, kwROWNUM, kwSYSDATE} {
			if !ctx.Candidates.HasToken(tok) {
				t.Fatalf("missing expression candidate %q; got %v", TokenName(tok), tokenNamesForTest(ctx.Candidates.Tokens))
			}
		}
	})
}

func collectCompletionMarked(t *testing.T, input string) *CompletionContext {
	t.Helper()
	cursor := strings.Index(input, "|")
	if cursor < 0 {
		t.Fatalf("input %q has no cursor marker", input)
	}
	sql := strings.Replace(input, "|", "", 1)
	ctx := CollectCompletion(sql, cursor)
	if ctx == nil {
		t.Fatal("CollectCompletion returned nil")
	}
	if ctx.Candidates == nil {
		t.Fatal("CollectCompletion returned nil Candidates")
	}
	if ctx.Scope == nil {
		t.Fatal("CollectCompletion returned nil Scope")
	}
	return ctx
}

func requireCompletionRule(t *testing.T, cs *CandidateSet, rule string) {
	t.Helper()
	if !cs.HasRule(rule) {
		t.Fatalf("missing rule %q; got rules=%v tokens=%v", rule, cs.Rules, tokenNamesForTest(cs.Tokens))
	}
}

func requireObjectKind(t *testing.T, intent *CompletionIntent, kind ObjectKind) {
	t.Helper()
	if intent == nil {
		t.Fatalf("nil completion intent; want kind %v", kind)
	}
	for _, got := range intent.ObjectKinds {
		if got == kind {
			return
		}
	}
	t.Fatalf("intent kinds = %v, want %v", intent.ObjectKinds, kind)
}

func hasObjectKind(intent *CompletionIntent, kind ObjectKind) bool {
	if intent == nil {
		return false
	}
	for _, got := range intent.ObjectKinds {
		if got == kind {
			return true
		}
	}
	return false
}

func requireRangeReference(t *testing.T, scope *ScopeSnapshot, aliasOrName string) {
	t.Helper()
	_ = findRangeReference(t, scope, aliasOrName)
}

func findRangeReference(t *testing.T, scope *ScopeSnapshot, aliasOrName string) RangeReference {
	t.Helper()
	for _, ref := range scope.References {
		if strings.EqualFold(ref.Alias, aliasOrName) || strings.EqualFold(ref.Name, aliasOrName) {
			return ref
		}
	}
	t.Fatalf("missing range reference %q in %#v", aliasOrName, scope.References)
	return RangeReference{}
}

func requireColumnName(t *testing.T, columns []string, want string) {
	t.Helper()
	for _, got := range columns {
		if strings.EqualFold(got, want) {
			return
		}
	}
	t.Fatalf("missing column %q in %v", want, columns)
}
