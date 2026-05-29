package parser

import (
	"reflect"
	"strings"
	"testing"
)

func TestOracleCompletionSelectExpansionSignals(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		rule   string
		kind   ObjectKind
		tokens []int
		refs   []string
	}{
		{name: "target comma", input: "SELECT a, | FROM employees", rule: "columnref", kind: ObjectKindColumn, refs: []string{"employees"}},
		{name: "distinct target", input: "SELECT DISTINCT | FROM employees", rule: "columnref", kind: ObjectKindColumn, refs: []string{"employees"}},
		{name: "from comma", input: "SELECT * FROM employees, |", rule: "table_ref", kind: ObjectKindTable},
		{name: "left join", input: "SELECT * FROM employees e LEFT JOIN |", rule: "table_ref", kind: ObjectKindTable, refs: []string{"e"}},
		{name: "right join", input: "SELECT * FROM employees e RIGHT JOIN |", rule: "table_ref", kind: ObjectKindTable, refs: []string{"e"}},
		{name: "full outer join", input: "SELECT * FROM employees e FULL OUTER JOIN |", rule: "table_ref", kind: ObjectKindTable, refs: []string{"e"}},
		{name: "cross join", input: "SELECT * FROM employees e CROSS JOIN |", rule: "table_ref", kind: ObjectKindTable, refs: []string{"e"}},
		{name: "natural join", input: "SELECT * FROM employees e NATURAL JOIN |", rule: "table_ref", kind: ObjectKindTable, refs: []string{"e"}},
		{name: "join using", input: "SELECT * FROM employees e JOIN departments d USING (|", rule: "columnref", kind: ObjectKindColumn, refs: []string{"e", "d"}},
		{name: "where and", input: "SELECT * FROM employees e WHERE id = 1 AND |", rule: "columnref", kind: ObjectKindColumn, refs: []string{"e"}},
		{name: "where or", input: "SELECT * FROM employees e WHERE id = 1 OR |", rule: "columnref", kind: ObjectKindColumn, refs: []string{"e"}},
		{name: "where operator", input: "SELECT * FROM employees e WHERE salary + | > 0", rule: "columnref", kind: ObjectKindColumn, refs: []string{"e"}},
		{name: "in list", input: "SELECT * FROM employees e WHERE id IN (|", rule: "columnref", kind: ObjectKindColumn, refs: []string{"e"}},
		{name: "between start", input: "SELECT * FROM employees e WHERE salary BETWEEN |", rule: "columnref", kind: ObjectKindColumn, refs: []string{"e"}},
		{name: "between end", input: "SELECT * FROM employees e WHERE salary BETWEEN 1 AND |", rule: "columnref", kind: ObjectKindColumn, refs: []string{"e"}},
		{name: "exists subquery", input: "SELECT * FROM employees e WHERE EXISTS (|", tokens: []int{SELECT}, refs: []string{"e"}},
		{name: "group comma", input: "SELECT id FROM employees e GROUP BY id, |", rule: "columnref", kind: ObjectKindColumn, refs: []string{"e"}},
		{name: "having", input: "SELECT id FROM employees e GROUP BY id HAVING |", rule: "columnref", kind: ObjectKindColumn, refs: []string{"e"}},
		{name: "order comma", input: "SELECT id FROM employees e ORDER BY id, |", rule: "columnref", kind: ObjectKindColumn, refs: []string{"e"}},
		{name: "from table clause keywords", input: "SELECT * FROM employees |", tokens: []int{WHERE, JOIN, GROUP, ORDER}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := collectCompletionMarked(t, tt.input)
			if tt.rule != "" {
				requireCompletionRule(t, ctx.Candidates, tt.rule)
			}
			if tt.kind != ObjectKindUnknown {
				requireObjectKind(t, ctx.Intent, tt.kind)
			}
			for _, tok := range tt.tokens {
				if !ctx.Candidates.HasToken(tok) {
					t.Fatalf("missing token %q; got %v", TokenName(tok), tokenNamesForTest(ctx.Candidates.Tokens))
				}
			}
			for _, ref := range tt.refs {
				requireRangeReference(t, ctx.Scope, ref)
			}
		})
	}
}

func TestOracleCompletionScopeExpansion(t *testing.T) {
	t.Run("incomplete join keeps left scope", func(t *testing.T) {
		ctx := collectCompletionMarked(t, "SELECT * FROM employees e JOIN |")
		if got, want := oracleRefNames(ctx.Scope.LocalReferences), []string{"e"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("local refs = %v, want %v", got, want)
		}
	})

	t.Run("subquery inner table does not leak into outer table refs", func(t *testing.T) {
		ctx := collectCompletionMarked(t, "SELECT | FROM (SELECT salary FROM payroll) p")
		if got, want := oracleRefNames(ctx.Scope.LocalReferences), []string{"p"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("local refs = %v, want %v", got, want)
		}
	})

	t.Run("nested subquery records outer references", func(t *testing.T) {
		ctx := collectCompletionMarked(t, "SELECT * FROM employees e WHERE EXISTS (SELECT | FROM departments d WHERE d.id = e.id)")
		if got, want := oracleRefNames(ctx.Scope.LocalReferences), []string{"d"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("local refs = %v, want %v", got, want)
		}
		if len(ctx.Scope.OuterReferences) != 1 {
			t.Fatalf("outer reference levels = %d, want 1", len(ctx.Scope.OuterReferences))
		}
		if got, want := oracleRefNames(ctx.Scope.OuterReferences[0]), []string{"e"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("outer refs = %v, want %v", got, want)
		}
	})

	t.Run("union uses current arm scope", func(t *testing.T) {
		ctx := collectCompletionMarked(t, "SELECT id FROM employees UNION SELECT | FROM departments")
		if got, want := oracleRefNames(ctx.Scope.LocalReferences), []string{"departments"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("local refs = %v, want %v", got, want)
		}
	})
}

func TestOracleCompletionDMLExpansionSignals(t *testing.T) {
	tests := []struct {
		name  string
		input string
		refs  []string
	}{
		{name: "insert column comma", input: "INSERT INTO employees (id, |)", refs: []string{"employees"}},
		{name: "insert values comma", input: "INSERT INTO employees VALUES (1, |)", refs: []string{"employees"}},
		{name: "insert select with cte", input: "WITH x AS (SELECT id FROM departments) INSERT INTO employees SELECT | FROM x", refs: []string{"x"}},
		{name: "update assignment comma", input: "UPDATE employees SET id = 1, |", refs: []string{"employees"}},
		{name: "merge update set", input: "MERGE INTO employees e USING departments d ON e.id = d.id WHEN MATCHED THEN UPDATE SET name = |", refs: []string{"e", "d"}},
		{name: "merge insert column list", input: "MERGE INTO employees e USING departments d ON e.id = d.id WHEN NOT MATCHED THEN INSERT (|)", refs: []string{"e", "d"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := collectCompletionMarked(t, tt.input)
			requireCompletionRule(t, ctx.Candidates, "columnref")
			requireObjectKind(t, ctx.Intent, ObjectKindColumn)
			for _, ref := range tt.refs {
				requireRangeReference(t, ctx.Scope, ref)
			}
		})
	}
}

func TestOracleCompletionDDLExpansionSignals(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		rule   string
		kind   ObjectKind
		tokens []int
		refs   []string
	}{
		{name: "create table column comma", input: "CREATE TABLE t (id NUMBER, |)", tokens: []int{kwCONSTRAINT, kwPRIMARY, kwUNIQUE, kwFOREIGN, kwCHECK}},
		{name: "create table column options", input: "CREATE TABLE t (id NUMBER |)", tokens: []int{kwNOT, kwNULL, kwDEFAULT, kwPRIMARY, kwUNIQUE, kwREFERENCES, kwCHECK}},
		{name: "primary key columns", input: "CREATE TABLE t (id NUMBER, CONSTRAINT pk PRIMARY KEY (|))", rule: "columnref", kind: ObjectKindColumn, refs: []string{"t"}},
		{name: "foreign key columns", input: "CREATE TABLE t (dept_id NUMBER, FOREIGN KEY (|) REFERENCES departments(id))", rule: "columnref", kind: ObjectKindColumn, refs: []string{"t"}},
		{name: "referenced columns", input: "CREATE TABLE t (dept_id NUMBER REFERENCES departments(|))", rule: "columnref", kind: ObjectKindColumn, refs: []string{"departments"}},
		{name: "alter modify column", input: "ALTER TABLE employees MODIFY |", rule: "columnref", kind: ObjectKindColumn, refs: []string{"employees"}},
		{name: "alter drop constraint", input: "ALTER TABLE employees DROP CONSTRAINT |", rule: "constraint_ref", kind: ObjectKindConstraint, refs: []string{"employees"}},
		{name: "create index table", input: "CREATE INDEX idx ON |", rule: "table_ref", kind: ObjectKindTable},
		{name: "create index columns", input: "CREATE INDEX idx ON employees (|)", rule: "columnref", kind: ObjectKindColumn, refs: []string{"employees"}},
		{name: "drop index", input: "DROP INDEX |", rule: "index_ref", kind: ObjectKindIndex},
		{name: "drop synonym", input: "DROP SYNONYM |", rule: "synonym_ref", kind: ObjectKindSynonym},
		{name: "drop trigger", input: "DROP TRIGGER |", rule: "trigger_ref", kind: ObjectKindTrigger},
		{name: "alter sequence", input: "ALTER SEQUENCE |", rule: "sequence_ref", kind: ObjectKindSequence},
		{name: "alter view", input: "ALTER VIEW |", rule: "table_ref", kind: ObjectKindView},
		{name: "alter procedure", input: "ALTER PROCEDURE |", rule: "proc_ref", kind: ObjectKindProcedure},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := collectCompletionMarked(t, tt.input)
			if tt.rule != "" {
				requireCompletionRule(t, ctx.Candidates, tt.rule)
			}
			if tt.kind != ObjectKindUnknown {
				requireObjectKind(t, ctx.Intent, tt.kind)
			}
			for _, tok := range tt.tokens {
				if !ctx.Candidates.HasToken(tok) {
					t.Fatalf("missing token %q; got %v", TokenName(tok), tokenNamesForTest(ctx.Candidates.Tokens))
				}
			}
			for _, ref := range tt.refs {
				requireRangeReference(t, ctx.Scope, ref)
			}
		})
	}
}

func TestOracleCompletionOracleSpecificSignals(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		rule      string
		kind      ObjectKind
		qualifier MultipartName
		tokens    []int
	}{
		{name: "sequence nextval member", input: "SELECT seq.| FROM dual", rule: "sequence_member_ref", kind: ObjectKindSequenceMember, qualifier: MultipartName{Object: "seq"}},
		{name: "schema package member", input: "SELECT pkg.| FROM dual", rule: "package_member_ref", kind: ObjectKindPackageMember, qualifier: MultipartName{Object: "pkg"}},
		{name: "package procedure call", input: "BEGIN pkg.|; END;", rule: "package_member_ref", kind: ObjectKindPackageMember, qualifier: MultipartName{Object: "pkg"}},
		{name: "table db link", input: "SELECT * FROM employees@|", rule: "database_link_ref", kind: ObjectKindDatabaseLink},
		{name: "select into variable", input: "SELECT id INTO | FROM employees", rule: "variable_ref", kind: ObjectKindVariable},
		{name: "plsql declaration type", input: "DECLARE v |; BEGIN NULL; END;", rule: "type_name", kind: ObjectKindType},
		{name: "plsql block statement", input: "BEGIN | END;", tokens: []int{SELECT, INSERT, UPDATE, DELETE}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := collectCompletionMarked(t, tt.input)
			if tt.rule != "" {
				requireCompletionRule(t, ctx.Candidates, tt.rule)
			}
			if tt.kind != ObjectKindUnknown {
				requireObjectKind(t, ctx.Intent, tt.kind)
			}
			if tt.qualifier != (MultipartName{}) && !equalMultipartNameFold(ctx.Intent.Qualifier, tt.qualifier) {
				t.Fatalf("qualifier = %+v, want %+v", ctx.Intent.Qualifier, tt.qualifier)
			}
			for _, tok := range tt.tokens {
				if !ctx.Candidates.HasToken(tok) {
					t.Fatalf("missing token %q; got %v", TokenName(tok), tokenNamesForTest(ctx.Candidates.Tokens))
				}
			}
		})
	}
}

func TestOracleCompletionPrefixExpansion(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		rule      string
		prefix    string
		qualifier MultipartName
	}{
		{name: "table prefix", input: "SELECT * FROM emp|", rule: "table_ref", prefix: "emp"},
		{name: "schema table prefix", input: "SELECT * FROM hr.emp|", rule: "table_ref", prefix: "emp", qualifier: MultipartName{Schema: "hr"}},
		{name: "column prefix", input: "SELECT e.na| FROM employees e", rule: "columnref", prefix: "na", qualifier: MultipartName{Object: "e"}},
		{name: "quoted table prefix", input: "SELECT * FROM \"Camel|", rule: "table_ref", prefix: "Camel"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := collectCompletionMarked(t, tt.input)
			requireCompletionRule(t, ctx.Candidates, tt.rule)
			if ctx.Prefix != tt.prefix {
				t.Fatalf("prefix = %q, want %q", ctx.Prefix, tt.prefix)
			}
			if tt.qualifier != (MultipartName{}) && !equalMultipartNameFold(ctx.Intent.Qualifier, tt.qualifier) {
				t.Fatalf("qualifier = %+v, want %+v", ctx.Intent.Qualifier, tt.qualifier)
			}
		})
	}
}

func oracleRefNames(refs []RangeReference) []string {
	var names []string
	for _, ref := range refs {
		if ref.Alias != "" {
			names = append(names, strings.ToLower(ref.Alias))
			continue
		}
		names = append(names, strings.ToLower(ref.Name))
	}
	return names
}

func equalMultipartNameFold(got MultipartName, want MultipartName) bool {
	return strings.EqualFold(got.Schema, want.Schema) && strings.EqualFold(got.Object, want.Object)
}
