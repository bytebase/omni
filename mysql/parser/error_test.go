package parser

import (
	"strings"
	"testing"
)

// parseExpectError parses sql and asserts that an error is returned whose
// message contains expectedSubstring.
func parseExpectError(t *testing.T, sql, expectedSubstring string) {
	t.Helper()
	_, err := Parse(sql)
	if err == nil {
		t.Fatalf("expected error for %q, got nil", sql)
	}
	if !strings.Contains(err.Error(), expectedSubstring) {
		t.Errorf("error %q does not contain %q (sql: %q)", err.Error(), expectedSubstring, sql)
	}
}

func TestParseError_Section_1_1_ErrorInfrastructure(t *testing.T) {
	// Verify that parsing invalid SQL returns an error (basic harness test).
	// The exact error message content will be improved in Phase 2+.
	cases := []struct {
		name     string
		sql      string
		contains string
	}{
		// Truncated CREATE returns error.
		{"truncated_create", "CREATE", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatalf("expected error for %q, got nil", tc.sql)
			}
			if tc.contains != "" && !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("error %q does not contain %q (sql: %q)", err.Error(), tc.contains, tc.sql)
			}
		})
	}
}

func TestSyntaxErrorAtCur_EOF(t *testing.T) {
	p := &Parser{
		lexer: NewLexer(""),
	}
	p.advance()

	err := p.syntaxErrorAtCur()
	if !strings.Contains(err.Error(), "at end of input") {
		t.Errorf("expected 'at end of input', got %q", err.Error())
	}
	if err.Line == 0 {
		t.Error("expected Line to be set")
	}
}

func TestSyntaxErrorAtTok_NonEOF(t *testing.T) {
	p := &Parser{
		lexer: NewLexer("SELECT"),
	}
	p.advance()

	err := p.syntaxErrorAtTok(p.cur)
	if !strings.Contains(err.Error(), "at or near") {
		t.Errorf("expected 'at or near', got %q", err.Error())
	}
	if !strings.Contains(err.Error(), "SELECT") {
		t.Errorf("expected token text 'SELECT' in error, got %q", err.Error())
	}
	if err.Line != 1 || err.Column != 1 {
		t.Errorf("expected line 1, column 1, got line %d, column %d", err.Line, err.Column)
	}
}

func TestParseError_PositionFormat(t *testing.T) {
	pe := &ParseError{
		Message:  "syntax error at end of input",
		Position: 10,
		Line:     1,
		Column:   11,
	}
	got := pe.Error()
	if got != "syntax error at end of input (line 1, column 11)" {
		t.Errorf("unexpected Error() output: %q", got)
	}
}

func TestParseError_NoPosition(t *testing.T) {
	pe := &ParseError{
		Message:  "unexpected token",
		Position: 5,
	}
	got := pe.Error()
	if got != "unexpected token" {
		t.Errorf("expected plain message, got %q", got)
	}
}

func TestParseError_Section_1_2_StmtDispatch(t *testing.T) {
	// Section 1.2: Verify that truncated dispatch calls in stmt.go return errors.
	// All 10 sub-parsers already return proper errors (never nil,nil),
	// so these tests serve as regression guards.
	cases := []struct {
		name     string
		sql      string
		contains string
	}{
		// parseChangeMasterStmtInner — truncated after MASTER (needs TO + options)
		{"change_master_trunc", "CHANGE MASTER TO", "at end of input"},
		// parseStartReplicaStmt — truncated mid-clause (FOR without CHANNEL name)
		{"start_replica_for_trunc", "START REPLICA FOR", "at end of input"},
		// parseStartGroupReplicationStmt — valid as standalone; test CHANGE REPLICATION with no target
		{"change_replication_trunc", "CHANGE REPLICATION", "expected SOURCE or FILTER"},
		// parseStopReplicaStmt — truncated mid-clause (FOR without CHANNEL name)
		{"stop_replica_for_trunc", "STOP REPLICA FOR", "at end of input"},
		// parseStopGroupReplicationStmt is always valid; test STOP without valid keyword
		{"stop_no_keyword", "STOP", "expected REPLICA"},
		// parseLoadIndexIntoCacheStmt — truncated after CACHE (needs table ref)
		{"load_index_trunc", "LOAD INDEX INTO CACHE", "at end of input"},
		// parseCreateResourceGroupStmt — truncated after GROUP (needs name)
		{"create_resource_group_trunc", "CREATE RESOURCE GROUP", "at end of input"},
		// parseAlterResourceGroupStmt — truncated after GROUP (needs name)
		{"alter_resource_group_trunc", "ALTER RESOURCE GROUP", "at end of input"},
		// parseDropSpatialRefSysStmt — truncated after SPATIAL (needs REFERENCE SYSTEM)
		{"drop_spatial_trunc", "DROP SPATIAL", "expected REFERENCE"},
		// parseDropResourceGroupStmt — truncated after GROUP (needs name)
		{"drop_resource_group_trunc", "DROP RESOURCE GROUP", "at end of input"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatalf("expected error for %q, got nil", tc.sql)
			}
			if !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("error %q does not contain %q (sql: %q)", err.Error(), tc.contains, tc.sql)
			}
		})
	}
}

func TestParseError_Section_1_3_DDLIgnoredErrors(t *testing.T) {
	// Section 1.3: Verify that truncated DDL inputs where parseIdentifier()
	// errors were previously discarded now produce proper errors.
	cases := []struct {
		name     string
		sql      string
		contains string
	}{
		// create_function.go: LANGUAGE identifier
		{"func_language_trunc", "CREATE FUNCTION f() RETURNS INT LANGUAGE", "at end of input"},
		// create_function.go: SQL SECURITY identifier
		{"func_sql_security_trunc", "CREATE FUNCTION f() RETURNS INT SQL SECURITY", "at end of input"},
		// create_table.go: column COLLATE identifier
		{"col_collate_trunc", "CREATE TABLE t (a INT COLLATE)", "expected identifier"},
		// create_database.go: CHARACTER SET identifier
		{"db_charset_trunc", "CREATE DATABASE db CHARACTER SET", "at end of input"},
		// create_database.go: CHARSET identifier
		{"db_charset_short_trunc", "CREATE DATABASE db CHARSET", "at end of input"},
		// create_database.go: COLLATE identifier
		{"db_collate_trunc", "CREATE DATABASE db COLLATE", "at end of input"},
		// alter_table.go: AFTER column identifier
		{"alter_after_trunc", "ALTER TABLE t ADD COLUMN c INT AFTER", "at end of input"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatalf("expected error for %q, got nil", tc.sql)
			}
			if !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("error %q does not contain %q (sql: %q)", err.Error(), tc.contains, tc.sql)
			}
		})
	}
}

func TestParseError_Section_1_4_DMLIgnoredErrors(t *testing.T) {
	// Section 1.4: Verify that truncated DML & other file inputs where
	// parseIdentifier() errors were previously discarded now produce proper errors.
	cases := []struct {
		name     string
		sql      string
		contains string
	}{
		// select.go: table alias parseIdentifier (2 sites — subquery AS at EOF)
		{"subquery_as_trunc", "SELECT * FROM (SELECT 1) AS", "at end of input"},
		{"lateral_subquery_as_trunc", "SELECT * FROM t, LATERAL (SELECT 1) AS", "at end of input"},
		// select.go: JSON_TABLE alias (required) at EOF
		{"json_table_alias_trunc", "SELECT * FROM JSON_TABLE('[1]', '$[*]' COLUMNS (a INT PATH '$.a')) AS", "at end of input"},
		// select.go: INTO OUTFILE charset at EOF
		{"outfile_charset_trunc", "SELECT 1 INTO OUTFILE 'f' CHARACTER SET", "at end of input"},
		{"outfile_charset_short_trunc", "SELECT 1 INTO OUTFILE 'f' CHARSET", "at end of input"},
		// set_show.go -> grant.go: SET DEFAULT ROLE at EOF (already errors)
		{"set_default_role_trunc", "SET DEFAULT ROLE", "expected"},
		// set_show.go -> grant.go: SET ROLE at EOF (already errors)
		{"set_role_trunc", "SET ROLE", "expected"},
		// grant.go: IDENTIFIED WITH at EOF (4 sites)
		{"create_user_with_trunc", "CREATE USER u IDENTIFIED WITH", "at end of input"},
		{"alter_user_with_trunc", "ALTER USER u IDENTIFIED WITH", "at end of input"},
		{"grant_host_at_trunc", "SET DEFAULT ROLE role1 TO user@", "expected"},
		{"rename_user_trunc", "RENAME USER", "at end of input"},
		// load_data.go: charset/file identifier (2 sites)
		{"load_data_charset_trunc", "LOAD DATA INFILE 'f' INTO TABLE t CHARACTER SET", "at end of input"},
		{"load_data_charset_short_trunc", "LOAD DATA INFILE 'f' INTO TABLE t CHARSET", "at end of input"},
		// replication.go: channel identifier (UNTIL pos name)
		{"repl_until_pos_trunc", "START REPLICA UNTIL SOURCE_LOG_FILE = 'f',", "at end of input"},
		// utility.go: identifier (2 sites — RESET bare, HELP bare)
		{"reset_bare", "RESET", "at end of input"},
		{"help_bare", "HELP", "at end of input"},
		// expr.go: window spec identifier — truncated WINDOW clause definition
		{"window_clause_trunc", "SELECT 1 WINDOW w AS (", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatalf("expected error for %q, got nil", tc.sql)
			}
			if tc.contains != "" && !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("error %q does not contain %q (sql: %q)", err.Error(), tc.contains, tc.sql)
			}
		})
	}
}

func TestLineCol(t *testing.T) {
	p := &Parser{lexer: NewLexer("SELECT\n  1 + 2")}
	// offset 0 -> line 1, col 1
	l, c := p.lineCol(0)
	if l != 1 || c != 1 {
		t.Errorf("offset 0: got line %d col %d, want 1 1", l, c)
	}
	// offset 7 (the space after newline) -> line 2, col 1
	l, c = p.lineCol(7)
	if l != 2 || c != 1 {
		t.Errorf("offset 7: got line %d col %d, want 2 1", l, c)
	}
	// offset 9 (the '1' on line 2) -> line 2, col 3
	l, c = p.lineCol(9)
	if l != 2 || c != 3 {
		t.Errorf("offset 9: got line %d col %d, want 2 3", l, c)
	}
}

func TestParseError_Section_2_1_ArithmeticOperators(t *testing.T) {
	cases := []struct {
		name     string
		sql      string
		contains string
	}{
		{"plus_eof", "SELECT 1 +", "at end of input"},
		{"minus_eof", "SELECT 1 -", "at end of input"},
		{"star_eof", "SELECT 1 *", "at end of input"},
		{"slash_eof", "SELECT 1 /", "at end of input"},
		{"mod_eof", "SELECT 1 %", "at end of input"},
		{"div_eof", "SELECT 1 DIV", "at end of input"},
		{"mod_kw_eof", "SELECT 1 MOD", "at end of input"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.contains)
			}
		})
	}
}

func TestParseError_Section_2_2_ComparisonOperators(t *testing.T) {
	cases := []struct {
		name     string
		sql      string
		contains string
	}{
		{"eq_eof", "SELECT 1 =", "at end of input"},
		{"lt_eof", "SELECT 1 <", "at end of input"},
		{"gt_eof", "SELECT 1 >", "at end of input"},
		{"lte_eof", "SELECT 1 <=", "at end of input"},
		{"gte_eof", "SELECT 1 >=", "at end of input"},
		{"neq_eof", "SELECT 1 <>", "at end of input"},
		{"neq_bang_eof", "SELECT 1 !=", "at end of input"},
		{"spaceship_eof", "SELECT 1 <=>", "at end of input"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.contains)
			}
		})
	}
}

func TestParseError_Section_2_3_LogicalBitwiseOperators(t *testing.T) {
	cases := []struct {
		name     string
		sql      string
		contains string
	}{
		{"and_eof", "SELECT 1 AND", "at end of input"},
		{"or_eof", "SELECT 1 OR", "at end of input"},
		{"xor_eof", "SELECT 1 XOR", "at end of input"},
		{"bitor_eof", "SELECT 1 |", "at end of input"},
		{"bitand_eof", "SELECT 1 &", "at end of input"},
		{"bitxor_eof", "SELECT 1 ^", "at end of input"},
		{"lshift_eof", "SELECT 1 <<", "at end of input"},
		{"rshift_eof", "SELECT 1 >>", "at end of input"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.contains)
			}
		})
	}
}

func TestParseError_Section_2_4_UnaryPrefixOperators(t *testing.T) {
	cases := []struct {
		name     string
		sql      string
		contains string
	}{
		{"not_eof", "SELECT NOT", "at end of input"},
		{"neg_eof", "SELECT -", "at end of input"},
		{"bitnot_eof", "SELECT ~", "at end of input"},
		{"bang_eof", "SELECT !", "at end of input"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.contains)
			}
		})
	}
}

func TestParseError_Section_3_1_BetweenLikeInRegexp(t *testing.T) {
	cases := []struct {
		name     string
		sql      string
		contains string
	}{
		{"between_eof", "SELECT 1 BETWEEN", "at end of input"},
		{"between_and_eof", "SELECT 1 BETWEEN 0 AND", "at end of input"},
		{"not_between_eof", "SELECT 1 NOT BETWEEN", "at end of input"},
		{"like_eof", "SELECT 'a' LIKE", "at end of input"},
		{"like_escape_eof", "SELECT 'a' LIKE 'b' ESCAPE", "at end of input"},
		{"not_like_eof", "SELECT 'a' NOT LIKE", "at end of input"},
		{"in_paren_eof", "SELECT 1 IN (", "at end of input"},
		{"in_comma_eof", "SELECT 1 IN (1,", "at end of input"},
		{"regexp_eof", "SELECT 1 REGEXP", "at end of input"},
		{"not_regexp_eof", "SELECT 1 NOT REGEXP", "at end of input"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.contains)
			}
		})
	}
}

func TestParseError_Section_3_2_CaseCastIsCollate(t *testing.T) {
	cases := []struct {
		name     string
		sql      string
		contains string
	}{
		{"case_when_eof", "SELECT CASE WHEN", "at end of input"},
		{"case_when_then_eof", "SELECT CASE WHEN 1 THEN", "at end of input"},
		{"cast_as_eof", "SELECT CAST(1 AS", "at end of input"},
		{"convert_using_eof", "SELECT CONVERT(1 USING", "at end of input"},
		{"is_eof", "SELECT 1 IS", "at end of input"},
		{"is_not_eof", "SELECT 1 IS NOT", "at end of input"},
		{"collate_eof", "SELECT a COLLATE", "at end of input"},
		{"sounds_eof", "SELECT a SOUNDS", "at end of input"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.contains)
			}
		})
	}
}

func TestParseError_Section_3_3_MidTokenErrors(t *testing.T) {
	cases := []struct {
		name     string
		sql      string
		contains string
	}{
		{"cast_plus_close_paren", "SELECT CAST(1 + )", "at or near"},
		{"plus_close_paren", "SELECT 1 + )", "at or near"},
		{"plus_semicolon", "SELECT 1 + ;", "at end of input"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.contains)
			}
		})
	}
}

func TestParseError_Section_4_1_SELECTClauses(t *testing.T) {
	cases := []struct {
		name     string
		sql      string
		contains string
	}{
		{"join_eof", "SELECT * FROM t JOIN", "at end of input"},
		{"left_join_eof", "SELECT * FROM t LEFT JOIN", "at end of input"},
		{"join_on_eof", "SELECT * FROM t JOIN t2 ON", "at end of input"},
		{"where_eof", "SELECT * FROM t WHERE", "at end of input"},
		{"group_by_eof", "SELECT * FROM t GROUP BY", "at end of input"},
		{"group_by_comma_eof", "SELECT * FROM t GROUP BY a,", "at end of input"},
		{"having_eof", "SELECT * FROM t HAVING", "at end of input"},
		{"order_by_eof", "SELECT * FROM t ORDER BY", "at end of input"},
		{"order_by_comma_eof", "SELECT * FROM t ORDER BY a,", "at end of input"},
		{"limit_eof", "SELECT * FROM t LIMIT", "at end of input"},
		{"union_eof", "SELECT 1 UNION", "at end of input"},
		{"union_all_eof", "SELECT 1 UNION ALL", "at end of input"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.contains)
			}
		})
	}
}

func TestParseError_Section_4_2_CTESubqueries(t *testing.T) {
	cases := []struct {
		name     string
		sql      string
		contains string
	}{
		{"with_eof", "WITH", "at end of input"},
		{"with_cte_as_paren_eof", "WITH cte AS (", "at end of input"},
		{"select_exists_paren_eof", "SELECT EXISTS (", "at end of input"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.contains)
			}
		})
	}
}

func TestParseError_Section_4_3_DDLTruncation(t *testing.T) {
	cases := []struct {
		name     string
		sql      string
		contains string
	}{
		{"create_table_default_eof", "CREATE TABLE t (a INT DEFAULT", "at end of input"},
		{"create_table_references_eof", "CREATE TABLE t (a INT REFERENCES", "at end of input"},
		{"create_table_check_paren_eof", "CREATE TABLE t (a INT CHECK (", "at end of input"},
		{"alter_table_add_column_eof", "ALTER TABLE t ADD COLUMN", "at end of input"},
		{"alter_table_drop_column_eof", "ALTER TABLE t DROP COLUMN", "at end of input"},
		{"alter_table_rename_to_eof", "ALTER TABLE t RENAME TO", "at end of input"},
		{"create_index_on_eof", "CREATE INDEX idx ON", "at end of input"},
		{"create_view_as_eof", "CREATE VIEW v AS", "at end of input"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.contains)
			}
		})
	}
}

func TestParseError_Section_4_4_DMLTruncation(t *testing.T) {
	cases := []struct {
		name     string
		sql      string
		contains string
	}{
		{"insert_values_paren_eof", "INSERT INTO t VALUES (", "at end of input"},
		{"insert_cols_paren_eof", "INSERT INTO t (", "at end of input"},
		{"insert_on_dup_key_eof", "INSERT INTO t VALUES (1) ON DUPLICATE KEY UPDATE a =", "at end of input"},
		{"update_set_eof", "UPDATE t SET", "at end of input"},
		{"update_set_eq_eof", "UPDATE t SET a =", "at end of input"},
		{"delete_where_eof", "DELETE FROM t WHERE", "at end of input"},
		{"set_var_eq_eof", "SET @x =", "at end of input"},
		{"use_eof", "USE", "at end of input"},
		{"drop_table_eof", "DROP TABLE", "at end of input"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.contains)
			}
		})
	}
}

func TestParseError_InvalidIntervalUnit(t *testing.T) {
	cases := []struct {
		name     string
		sql      string
		contains string
	}{
		{"INVALID_UNIT", "SELECT DATE_ADD(d, INTERVAL 1 INVALID_UNIT) FROM t", "invalid INTERVAL unit"},
		{"FOOBAR", "SELECT DATE_ADD(d, INTERVAL 1 FOOBAR) FROM t", "invalid INTERVAL unit"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parseExpectError(t, tc.sql, tc.contains)
		})
	}
}
