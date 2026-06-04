package parser

import (
	"testing"
)

// This file is the parser-dcl-tcl node's structural correctness gate for the
// prepared-statement family (PREPARE / DEALLOCATE / EXECUTE / EXECUTE IMMEDIATE
// / DESCRIBE INPUT / DESCRIBE OUTPUT). The authoritative accept/reject
// differential against the live Trino 481 oracle lives in
// oracle_dcl_tcl_test.go.

func TestPrepared_Structure(t *testing.T) {
	t.Run("prepare_select", func(t *testing.T) {
		stmt, ok := dclParseOne(t, "PREPARE my_select FROM SELECT * FROM nation").(*PrepareStmt)
		if !ok {
			t.Fatalf("got %T, want *PrepareStmt", stmt)
		}
		if stmt.Name.Normalize() != "my_select" {
			t.Errorf("Name = %q, want my_select", stmt.Name.Normalize())
		}
		if stmt.Body != "SELECT * FROM nation" {
			t.Errorf("Body = %q, want %q", stmt.Body, "SELECT * FROM nation")
		}
	})

	t.Run("prepare_insert", func(t *testing.T) {
		stmt := dclParseOne(t, "PREPARE my_insert FROM INSERT INTO cities VALUES (1, 'San Francisco')").(*PrepareStmt)
		if stmt.Body != "INSERT INTO cities VALUES (1, 'San Francisco')" {
			t.Errorf("Body = %q", stmt.Body)
		}
	})

	t.Run("prepare_with_params", func(t *testing.T) {
		// The placeholder '?' and the rest of the inner statement are captured
		// verbatim; the ';' inside a string would have been handled by Split,
		// but here the whole tail is one statement.
		stmt := dclParseOne(t, "PREPARE p FROM SELECT name FROM nation WHERE regionkey = ? AND nationkey < ?").(*PrepareStmt)
		want := "SELECT name FROM nation WHERE regionkey = ? AND nationkey < ?"
		if stmt.Body != want {
			t.Errorf("Body = %q, want %q", stmt.Body, want)
		}
	})

	t.Run("deallocate", func(t *testing.T) {
		stmt, ok := dclParseOne(t, "DEALLOCATE PREPARE my_query").(*DeallocateStmt)
		if !ok {
			t.Fatalf("got %T, want *DeallocateStmt", stmt)
		}
		if stmt.Name.Normalize() != "my_query" {
			t.Errorf("Name = %q, want my_query", stmt.Name.Normalize())
		}
	})

	t.Run("execute_no_using", func(t *testing.T) {
		stmt, ok := dclParseOne(t, "EXECUTE my_select1").(*ExecuteStmt)
		if !ok {
			t.Fatalf("got %T, want *ExecuteStmt", stmt)
		}
		if stmt.Name.Normalize() != "my_select1" {
			t.Errorf("Name = %q", stmt.Name.Normalize())
		}
		if len(stmt.Using) != 0 {
			t.Errorf("got %d USING args, want 0", len(stmt.Using))
		}
	})

	t.Run("execute_using", func(t *testing.T) {
		stmt := dclParseOne(t, "EXECUTE my_select2 USING 1, 3").(*ExecuteStmt)
		if len(stmt.Using) != 2 {
			t.Fatalf("got %d USING args, want 2", len(stmt.Using))
		}
	})

	t.Run("execute_immediate", func(t *testing.T) {
		stmt, ok := dclParseOne(t, "EXECUTE IMMEDIATE 'SELECT name FROM nation'").(*ExecuteImmediateStmt)
		if !ok {
			t.Fatalf("got %T, want *ExecuteImmediateStmt", stmt)
		}
		if stmt.SQL != "SELECT name FROM nation" {
			t.Errorf("SQL = %q", stmt.SQL)
		}
		if len(stmt.Using) != 0 {
			t.Errorf("got %d USING args, want 0", len(stmt.Using))
		}
	})

	t.Run("execute_immediate_using", func(t *testing.T) {
		stmt := dclParseOne(t, "EXECUTE IMMEDIATE 'SELECT name FROM nation WHERE regionkey = ? AND nationkey < ?' USING 1, 3").(*ExecuteImmediateStmt)
		if len(stmt.Using) != 2 {
			t.Fatalf("got %d USING args, want 2", len(stmt.Using))
		}
	})

	t.Run("execute_named_immediate", func(t *testing.T) {
		// IMMEDIATE not followed by a string is a prepared-statement NAME
		// (IMMEDIATE is non-reserved), so this is the named-execute form.
		stmt, ok := dclParseOne(t, "EXECUTE IMMEDIATE").(*ExecuteStmt)
		if !ok {
			t.Fatalf("got %T, want *ExecuteStmt (named execute of 'immediate')", stmt)
		}
		if stmt.Name.Normalize() != "immediate" {
			t.Errorf("Name = %q, want immediate", stmt.Name.Normalize())
		}
	})

	t.Run("execute_named_immediate_using", func(t *testing.T) {
		stmt := dclParseOne(t, "EXECUTE IMMEDIATE USING 1").(*ExecuteStmt)
		if stmt.Name.Normalize() != "immediate" || len(stmt.Using) != 1 {
			t.Errorf("got Name=%q Using=%d, want immediate / 1", stmt.Name.Normalize(), len(stmt.Using))
		}
	})

	t.Run("describe_input", func(t *testing.T) {
		stmt, ok := dclParseOne(t, "DESCRIBE INPUT my_select1").(*DescribeInputStmt)
		if !ok {
			t.Fatalf("got %T, want *DescribeInputStmt", stmt)
		}
		if stmt.Name.Normalize() != "my_select1" {
			t.Errorf("Name = %q", stmt.Name.Normalize())
		}
	})

	t.Run("describe_output", func(t *testing.T) {
		stmt, ok := dclParseOne(t, "DESCRIBE OUTPUT my_select1").(*DescribeOutputStmt)
		if !ok {
			t.Fatalf("got %T, want *DescribeOutputStmt", stmt)
		}
		if stmt.Name.Normalize() != "my_select1" {
			t.Errorf("Name = %q", stmt.Name.Normalize())
		}
	})
}

func TestPrepared_Negative(t *testing.T) {
	// PREPARE with no/invalid inner statement; DEALLOCATE missing PREPARE
	// keyword; EXECUTE IMMEDIATE <ident> (parsed as named-execute "immediate"
	// with a trailing token); EXECUTE with a dangling USING. The oracle
	// differential confirms each matches Trino 481.
	dclParseErr(t, "PREPARE p FROM")
	dclParseErr(t, "PREPARE p FROM 1")             // inner is not a statement
	dclParseErr(t, "PREPARE p FROM NOTASTATEMENT") // inner is not a statement
	dclParseErr(t, "DEALLOCATE my_query")
	dclParseErr(t, "EXECUTE IMMEDIATE my_select") // "immediate" name + trailing token
	dclParseErr(t, "EXECUTE my_select USING")
}
