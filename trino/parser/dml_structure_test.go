package parser

import (
	"testing"
)

// This file is the parser-dml node's correctness gate for the structural layer:
// it asserts the AST shape produced for INSERT / DELETE / UPDATE / MERGE /
// TRUNCATE. The authoritative accept/reject differential against the live
// Trino 481 oracle lives in oracle_dml_test.go.

// dmlParseOne parses exactly one statement via the public Parse entry point and
// returns it, failing the test if parsing errored or did not yield exactly one
// statement.
func dmlParseOne(t *testing.T, sql string) interface{} {
	t.Helper()
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q): unexpected errors: %v", sql, errs)
	}
	if file == nil || len(file.Stmts) != 1 {
		got := 0
		if file != nil {
			got = len(file.Stmts)
		}
		t.Fatalf("Parse(%q): got %d statements, want 1", sql, got)
	}
	return file.Stmts[0]
}

// dmlParseErr asserts that Parse reports at least one error for sql (a
// rejected-by-the-grammar input). The oracle differential confirms these
// rejections match Trino 481.
func dmlParseErr(t *testing.T, sql string) {
	t.Helper()
	_, errs := Parse(sql)
	if len(errs) == 0 {
		t.Errorf("Parse(%q): want at least one error, got none", sql)
	}
}

func TestInsert_Structure(t *testing.T) {
	t.Run("values", func(t *testing.T) {
		stmt, ok := dmlParseOne(t, "INSERT INTO cities VALUES (1, 'San Francisco')").(*InsertStmt)
		if !ok {
			t.Fatalf("got %T, want *InsertStmt", stmt)
		}
		if got := stmt.Target.Normalize(); got != "cities" {
			t.Errorf("target = %q, want cities", got)
		}
		if stmt.Columns != nil {
			t.Errorf("columns = %v, want nil", stmt.Columns)
		}
		if stmt.Source == nil {
			t.Fatalf("source = nil, want a query")
		}
		if _, ok := stmt.Source.Body.(*ValuesQuery); !ok {
			t.Errorf("source body = %T, want *ValuesQuery", stmt.Source.Body)
		}
	})

	t.Run("explicit_columns_select", func(t *testing.T) {
		stmt := dmlParseOne(t, "INSERT INTO nation (nationkey, name, regionkey) SELECT 1, 'a', 2").(*InsertStmt)
		if len(stmt.Columns) != 3 {
			t.Fatalf("got %d columns, want 3", len(stmt.Columns))
		}
		want := []string{"nationkey", "name", "regionkey"}
		for i, c := range stmt.Columns {
			if c.Normalize() != want[i] {
				t.Errorf("column[%d] = %q, want %q", i, c.Normalize(), want[i])
			}
		}
		if _, ok := stmt.Source.Body.(*QuerySpec); !ok {
			t.Errorf("source body = %T, want *QuerySpec", stmt.Source.Body)
		}
	})

	t.Run("select_star", func(t *testing.T) {
		stmt := dmlParseOne(t, "INSERT INTO orders SELECT * FROM new_orders").(*InsertStmt)
		if stmt.Columns != nil {
			t.Errorf("columns = %v, want nil", stmt.Columns)
		}
	})

	t.Run("qualified_target_with_quoted_part", func(t *testing.T) {
		stmt := dmlParseOne(t, `INSERT INTO a."b/c".d SELECT * FROM t`).(*InsertStmt)
		if len(stmt.Target.Parts) != 3 {
			t.Fatalf("got %d target parts, want 3", len(stmt.Target.Parts))
		}
		if stmt.Target.Parts[1].Value != "b/c" {
			t.Errorf("middle part = %q, want b/c", stmt.Target.Parts[1].Value)
		}
	})

	t.Run("parenthesized_source_query", func(t *testing.T) {
		// `( SELECT … )` is a parenthesized source query, NOT a column list.
		stmt := dmlParseOne(t, "INSERT INTO t (SELECT 1)").(*InsertStmt)
		if stmt.Columns != nil {
			t.Errorf("columns = %v, want nil (the parens are a subquery)", stmt.Columns)
		}
		if _, ok := stmt.Source.Body.(*ParenQuery); !ok {
			t.Errorf("source body = %T, want *ParenQuery", stmt.Source.Body)
		}
	})

	t.Run("table_query_source", func(t *testing.T) {
		stmt := dmlParseOne(t, "INSERT INTO t TABLE other").(*InsertStmt)
		if _, ok := stmt.Source.Body.(*TableQuery); !ok {
			t.Errorf("source body = %T, want *TableQuery", stmt.Source.Body)
		}
	})

	t.Run("with_cte_source", func(t *testing.T) {
		stmt := dmlParseOne(t, "INSERT INTO t WITH cte AS (SELECT 1) SELECT * FROM cte").(*InsertStmt)
		if stmt.Source.With == nil {
			t.Errorf("source.With = nil, want a CTE clause")
		}
	})

	// negatives (oracle confirms these are SYNTAX_ERROR in Trino 481)
	t.Run("missing_source", func(t *testing.T) { dmlParseErr(t, "INSERT INTO t") })
	t.Run("columns_no_source", func(t *testing.T) { dmlParseErr(t, "INSERT INTO t (a, b)") })
	t.Run("four_part_target", func(t *testing.T) { dmlParseErr(t, "INSERT INTO a.b.c.d VALUES (1)") })
	t.Run("missing_into", func(t *testing.T) { dmlParseErr(t, "INSERT t VALUES (1)") })
}

func TestTruncate_Structure(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		stmt, ok := dmlParseOne(t, "TRUNCATE TABLE orders").(*TruncateStmt)
		if !ok {
			t.Fatalf("got %T, want *TruncateStmt", stmt)
		}
		if got := stmt.Target.Normalize(); got != "orders" {
			t.Errorf("target = %q, want orders", got)
		}
	})

	t.Run("three_part", func(t *testing.T) {
		stmt := dmlParseOne(t, "TRUNCATE TABLE a.b.c").(*TruncateStmt)
		if len(stmt.Target.Parts) != 3 {
			t.Errorf("got %d parts, want 3", len(stmt.Target.Parts))
		}
	})

	// negatives
	t.Run("missing_table_keyword", func(t *testing.T) { dmlParseErr(t, "TRUNCATE orders") })
	t.Run("four_part_target", func(t *testing.T) { dmlParseErr(t, "TRUNCATE TABLE a.b.c.d") })
	t.Run("trailing_token", func(t *testing.T) { dmlParseErr(t, "TRUNCATE TABLE orders extra") })
}

func TestDelete_Structure(t *testing.T) {
	t.Run("no_where", func(t *testing.T) {
		stmt, ok := dmlParseOne(t, "DELETE FROM orders").(*DeleteStmt)
		if !ok {
			t.Fatalf("got %T, want *DeleteStmt", stmt)
		}
		if got := stmt.Target.Normalize(); got != "orders" {
			t.Errorf("target = %q, want orders", got)
		}
		if stmt.Where != nil {
			t.Errorf("where = %v, want nil", stmt.Where)
		}
	})

	t.Run("with_where", func(t *testing.T) {
		stmt := dmlParseOne(t, "DELETE FROM lineitem WHERE shipmode = 'AIR'").(*DeleteStmt)
		if stmt.Where == nil {
			t.Errorf("where = nil, want a predicate")
		}
	})

	t.Run("where_subquery", func(t *testing.T) {
		stmt := dmlParseOne(t, "DELETE FROM lineitem WHERE orderkey IN (SELECT orderkey FROM orders WHERE x = 1)").(*DeleteStmt)
		if stmt.Where == nil {
			t.Errorf("where = nil, want an IN-subquery predicate")
		}
	})

	t.Run("quoted_target", func(t *testing.T) {
		stmt := dmlParseOne(t, `DELETE FROM "awesome table"`).(*DeleteStmt)
		if stmt.Target.Parts[0].Value != "awesome table" {
			t.Errorf("target = %q, want \"awesome table\"", stmt.Target.Parts[0].Value)
		}
	})

	// negatives
	t.Run("missing_from", func(t *testing.T) { dmlParseErr(t, "DELETE orders") })
	t.Run("dangling_where", func(t *testing.T) { dmlParseErr(t, "DELETE FROM orders WHERE") })
	t.Run("four_part_target", func(t *testing.T) { dmlParseErr(t, "DELETE FROM a.b.c.d") })
	t.Run("trailing_token", func(t *testing.T) { dmlParseErr(t, "DELETE FROM orders garbage") })
}

func TestUpdate_Structure(t *testing.T) {
	t.Run("single_assignment_with_where", func(t *testing.T) {
		stmt, ok := dmlParseOne(t, "UPDATE purchases SET status = 'OVERDUE' WHERE ship_date IS NULL").(*UpdateStmt)
		if !ok {
			t.Fatalf("got %T, want *UpdateStmt", stmt)
		}
		if got := stmt.Target.Normalize(); got != "purchases" {
			t.Errorf("target = %q, want purchases", got)
		}
		if len(stmt.Assignments) != 1 {
			t.Fatalf("got %d assignments, want 1", len(stmt.Assignments))
		}
		if stmt.Assignments[0].Column.Normalize() != "status" {
			t.Errorf("assignment column = %q, want status", stmt.Assignments[0].Column.Normalize())
		}
		if stmt.Where == nil {
			t.Errorf("where = nil, want a predicate")
		}
	})

	t.Run("multi_assignment_no_where", func(t *testing.T) {
		stmt := dmlParseOne(t, "UPDATE customers SET account_manager = 'John Henry', assign_date = now()").(*UpdateStmt)
		if len(stmt.Assignments) != 2 {
			t.Fatalf("got %d assignments, want 2", len(stmt.Assignments))
		}
		if stmt.Where != nil {
			t.Errorf("where = %v, want nil", stmt.Where)
		}
	})

	t.Run("subquery_value", func(t *testing.T) {
		stmt := dmlParseOne(t, "UPDATE new_hires SET manager = (SELECT e.name FROM employees e WHERE e.employee_id = new_hires.manager_id)").(*UpdateStmt)
		if len(stmt.Assignments) != 1 {
			t.Fatalf("got %d assignments, want 1", len(stmt.Assignments))
		}
		if stmt.Assignments[0].Value == nil {
			t.Errorf("assignment value = nil, want a scalar subquery")
		}
	})

	t.Run("quoted_target_column", func(t *testing.T) {
		stmt := dmlParseOne(t, `UPDATE t SET "col one" = a + b * 2, x = CASE WHEN y THEN 1 ELSE 2 END`).(*UpdateStmt)
		if stmt.Assignments[0].Column.Value != "col one" {
			t.Errorf("first column = %q, want \"col one\"", stmt.Assignments[0].Column.Value)
		}
	})

	// negatives
	t.Run("missing_set", func(t *testing.T) { dmlParseErr(t, "UPDATE t a = 1") })
	t.Run("dangling_set", func(t *testing.T) { dmlParseErr(t, "UPDATE t SET") })
	t.Run("row_assignment", func(t *testing.T) { dmlParseErr(t, "UPDATE t SET (a, b) = (1, 2)") })
	t.Run("four_part_target", func(t *testing.T) { dmlParseErr(t, "UPDATE a.b.c.d SET x = 1") })
	t.Run("trailing_token", func(t *testing.T) { dmlParseErr(t, "UPDATE t SET a = 1 garbage") })
}

func TestMerge_Structure(t *testing.T) {
	t.Run("matched_delete", func(t *testing.T) {
		stmt, ok := dmlParseOne(t, "MERGE INTO accounts t USING monthly s ON t.customer = s.customer WHEN MATCHED THEN DELETE").(*MergeStmt)
		if !ok {
			t.Fatalf("got %T, want *MergeStmt", stmt)
		}
		if got := stmt.Target.Normalize(); got != "accounts" {
			t.Errorf("target = %q, want accounts", got)
		}
		if stmt.Alias == nil || stmt.Alias.Normalize() != "t" {
			t.Errorf("alias = %v, want t", stmt.Alias)
		}
		if stmt.Source == nil {
			t.Fatalf("source = nil")
		}
		if stmt.On == nil {
			t.Fatalf("on = nil")
		}
		if len(stmt.Whens) != 1 {
			t.Fatalf("got %d WHEN clauses, want 1", len(stmt.Whens))
		}
		if stmt.Whens[0].Kind != MergeDelete {
			t.Errorf("clause kind = %v, want MergeDelete", stmt.Whens[0].Kind)
		}
	})

	t.Run("as_alias", func(t *testing.T) {
		stmt := dmlParseOne(t, "MERGE INTO accounts AS t USING src s ON t.c = s.c WHEN MATCHED THEN DELETE").(*MergeStmt)
		if stmt.Alias == nil || stmt.Alias.Normalize() != "t" {
			t.Errorf("alias = %v, want t", stmt.Alias)
		}
	})

	t.Run("no_alias", func(t *testing.T) {
		stmt := dmlParseOne(t, "MERGE INTO accounts USING src s ON accounts.c = s.c WHEN MATCHED THEN DELETE").(*MergeStmt)
		if stmt.Alias != nil {
			t.Errorf("alias = %v, want nil", stmt.Alias)
		}
	})

	t.Run("update_insert", func(t *testing.T) {
		stmt := dmlParseOne(t, "MERGE INTO accounts t USING src s ON (t.customer = s.customer) "+
			"WHEN MATCHED THEN UPDATE SET purchases = s.purchases + t.purchases "+
			"WHEN NOT MATCHED THEN INSERT (customer, purchases, address) VALUES (s.customer, s.purchases, s.address)").(*MergeStmt)
		if len(stmt.Whens) != 2 {
			t.Fatalf("got %d WHEN clauses, want 2", len(stmt.Whens))
		}
		upd := stmt.Whens[0]
		if upd.Kind != MergeUpdate {
			t.Errorf("clause[0] kind = %v, want MergeUpdate", upd.Kind)
		}
		if len(upd.Assignments) != 1 {
			t.Errorf("update assignments = %d, want 1", len(upd.Assignments))
		}
		ins := stmt.Whens[1]
		if ins.Kind != MergeInsert {
			t.Errorf("clause[1] kind = %v, want MergeInsert", ins.Kind)
		}
		if len(ins.Columns) != 3 {
			t.Errorf("insert columns = %d, want 3", len(ins.Columns))
		}
		if len(ins.Values) != 3 {
			t.Errorf("insert values = %d, want 3", len(ins.Values))
		}
	})

	t.Run("conditions_and_multi_assignment", func(t *testing.T) {
		stmt := dmlParseOne(t, "MERGE INTO accounts t USING src s ON (t.customer = s.customer) "+
			"WHEN MATCHED AND s.address = 'Centreville' THEN DELETE "+
			"WHEN MATCHED THEN UPDATE SET purchases = s.purchases + t.purchases, address = s.address "+
			"WHEN NOT MATCHED THEN INSERT (customer, purchases, address) VALUES (s.customer, s.purchases, s.address)").(*MergeStmt)
		if len(stmt.Whens) != 3 {
			t.Fatalf("got %d WHEN clauses, want 3", len(stmt.Whens))
		}
		if stmt.Whens[0].Condition == nil {
			t.Errorf("clause[0] (matched delete) condition = nil, want AND guard")
		}
		if stmt.Whens[0].Kind != MergeDelete {
			t.Errorf("clause[0] kind = %v, want MergeDelete", stmt.Whens[0].Kind)
		}
		if len(stmt.Whens[1].Assignments) != 2 {
			t.Errorf("clause[1] assignments = %d, want 2", len(stmt.Whens[1].Assignments))
		}
	})

	t.Run("insert_no_column_list", func(t *testing.T) {
		stmt := dmlParseOne(t, "MERGE INTO acc t USING src s ON t.c = s.c WHEN NOT MATCHED THEN INSERT VALUES (s.c, s.p)").(*MergeStmt)
		ins := stmt.Whens[0]
		if ins.Columns != nil {
			t.Errorf("insert columns = %v, want nil", ins.Columns)
		}
		if len(ins.Values) != 2 {
			t.Errorf("insert values = %d, want 2", len(ins.Values))
		}
	})

	t.Run("subquery_source", func(t *testing.T) {
		stmt := dmlParseOne(t, "MERGE INTO acc t USING (SELECT * FROM src) s ON t.c = s.c WHEN MATCHED THEN DELETE").(*MergeStmt)
		if stmt.Source == nil {
			t.Errorf("source = nil, want a subquery relation")
		}
	})

	t.Run("legacy_corpus_full", func(t *testing.T) {
		// The full legacy examples/merge.sql case.
		stmt := dmlParseOne(t, "MERGE INTO inventory AS i USING changes AS c ON i.part = c.part "+
			"WHEN MATCHED AND c.action = 'mod' THEN UPDATE SET qty = qty + c.qty, ts = CURRENT_TIMESTAMP "+
			"WHEN MATCHED AND c.action = 'del' THEN DELETE "+
			"WHEN NOT MATCHED AND c.action = 'new' THEN INSERT (part, qty) VALUES (c.part, c.qty)").(*MergeStmt)
		if len(stmt.Whens) != 3 {
			t.Fatalf("got %d WHEN clauses, want 3", len(stmt.Whens))
		}
	})

	// negatives
	t.Run("no_when", func(t *testing.T) { dmlParseErr(t, "MERGE INTO accounts t USING src s ON t.c = s.c") })
	t.Run("update_set_parenthesized", func(t *testing.T) {
		dmlParseErr(t, "MERGE INTO accounts t USING src s ON t.c = s.c WHEN MATCHED THEN UPDATE SET (p = s.p)")
	})
	t.Run("insert_empty_values", func(t *testing.T) {
		dmlParseErr(t, "MERGE INTO acc t USING src s ON t.c = s.c WHEN NOT MATCHED THEN INSERT VALUES ()")
	})
	t.Run("four_part_target", func(t *testing.T) {
		dmlParseErr(t, "MERGE INTO a.b.c.d t USING s ON t.c = s.c WHEN MATCHED THEN DELETE")
	})
	t.Run("not_matched_action_must_be_insert", func(t *testing.T) {
		// A WHEN NOT MATCHED clause's only legal action is INSERT; DELETE/UPDATE/
		// any other token before VALUES must NOT be silently accepted as INSERT.
		dmlParseErr(t, "MERGE INTO acc t USING src s ON t.c = s.c WHEN NOT MATCHED THEN DELETE VALUES (1)")
		dmlParseErr(t, "MERGE INTO acc t USING src s ON t.c = s.c WHEN NOT MATCHED THEN UPDATE VALUES (1)")
		dmlParseErr(t, "MERGE INTO acc t USING src s ON t.c = s.c WHEN NOT MATCHED THEN foo VALUES (1)")
	})
}
