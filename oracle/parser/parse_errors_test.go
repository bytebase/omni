package parser

import (
	"errors"
	"testing"
)

// TestParse_ReturnsParseErrors pins the error contract every omni engine
// shares: Parse returns nil on valid SQL, and otherwise returns an error that
// errors.As resolves to a *ParseError carrying a non-empty Message and a byte
// Position pointing at the offending token.
func TestParse_ReturnsParseErrors(t *testing.T) {
	if _, err := Parse("SELECT 1 FROM dual"); err != nil {
		t.Fatalf("Parse(valid) = %v, want nil", err)
	}

	_, err := Parse("SELECT 1 FROM FROM")
	if err == nil {
		t.Fatal("Parse(invalid) returned nil error")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("errors.As(*ParseError) failed on %T", err)
	}
	if pe.Position != len("SELECT 1 FROM ") {
		t.Errorf("ParseError.Position = %d, want %d (the second FROM)", pe.Position, len("SELECT 1 FROM "))
	}
	if pe.Message == "" {
		t.Error("ParseError.Message is empty")
	}
}

// TestErrorPositions pins omni's exact ParseError.Position byte offsets across:
//   - representative compat.tsv rejections where omni pinpoints the offending
//     token more sharply than Oracle (stored PL/SQL ORA-24344 with no byte,
//     server-crash statements with no byte, statement-level errors Oracle
//     reports at byte 0, coarse clause-start positions, and FLASHBACK ARCHIVE
//     token-end positions), so switching compat.tsv to coarse accept/reject
//     verdicts never degrades diagnostic byte offsets; and
//   - core SQL/DDL/DML/PLSQL syntax, multi-line, multi-byte UTF-8, and lexer
//     errors.
func TestErrorPositions(t *testing.T) {
	cases := []struct {
		name string
		sql  string
		pos  int
	}{
		// 1. Stored PL/SQL units and crashes where Oracle reports no byte offset.
		{"create_function_missing_name", "CREATE FUNCTION RETURN NUMBER IS BEGIN RETURN 1; END;", 23},
		{"create_function_missing_return", "CREATE FUNCTION a NUMBER IS BEGIN RETURN 1; END;", 18},
		{"create_function_missing_return_type", "CREATE FUNCTION a RETURN IS BEGIN RETURN 1; END;", 28},
		{"create_function_missing_is", "CREATE FUNCTION a RETURN NUMBER BEGIN RETURN 1; END;", 32},
		{"create_function_truncated_after_is", "CREATE FUNCTION a RETURN NUMBER IS", 34},
		{"create_function_trailing_token", "CREATE FUNCTION a RETURN NUMBER IS BEGIN RETURN 1; END; zzz", 56},
		{"create_package_body_missing_is", "CREATE PACKAGE BODY a END;", 22},
		{"create_procedure_missing_body", "CREATE PROCEDURE a IS", 21},
		{"create_type_unclosed_object", "CREATE TYPE a AS OBJECT (a INT", 30},
		{"create_java_source_using_eof", "CREATE JAVA SOURCE USING", 24},

		// 2. Statement-level errors Oracle reports at byte 0.
		{"alter_dimension_missing_name", "ALTER DIMENSION COMPILE", 23},
		{"alter_dimension_missing_action", "ALTER DIMENSION a", 17},
		{"create_audit_policy_bad_option", "CREATE AUDIT POLICY a zzz", 22},
		{"alter_audit_policy_trailing_token", "ALTER AUDIT POLICY a CONDITION 's' EVALUATE PER STATEMENT zzz", 58},
		{"alter_json_duality_view_missing_name", "ALTER JSON RELATIONAL DUALITY VIEW ENABLE LOGICAL REPLICATION", 42},
		{"alter_json_duality_view_missing_enable", "ALTER JSON RELATIONAL DUALITY VIEW a LOGICAL REPLICATION", 37},
		{"alter_type_modify_instantiable", "ALTER TYPE a MODIFY INSTANTIABLE", 20},
		{"alter_type_cascade_force_missing_exceptions", "ALTER TYPE a ADD MAP MEMBER FUNCTION a RETURN NUMBER CASCADE FORCE", 66},
		{"alter_pluggable_database_recover_standby", "ALTER PLUGGABLE DATABASE RECOVER STANDBY", 33},

		// 3. Coarse clause-start positions where omni pinpoints the offending token.
		{"drop_tablespace_set_trailing_token", "DROP TABLESPACE SET a zzz", 22},
		{"alter_database_recover_until_time_trailing", "ALTER DATABASE RECOVER STANDBY DATABASE UNTIL TIME 'YYYY-MM-DD:HH24:MI:SS' zzz", 75},
		{"alter_dimension_add_hierarchy_trailing", "ALTER DIMENSION a ADD HIERARCHY a(a CHILD OF a) zzz", 48},
		{"alter_dimension_drop_attribute_determines", "ALTER DIMENSION a DROP ATTRIBUTE a DETERMINES(a)", 35},
		{"alter_dimension_drop_attribute_parens", "ALTER DIMENSION a DROP ATTRIBUTE a(a)", 34},

		// 4. Token-end and skipped-token positions (FLASHBACK ARCHIVE, MERGE).
		{"alter_flashback_archive_default", "ALTER FLASHBACK ARCHIVE a DEFAULT", 26},
		{"create_flashback_archive_missing_name", "CREATE FLASHBACK ARCHIVE TABLESPACE a RETENTION 1 DAY", 36},
		{"create_flashback_archive_extra_ident", "CREATE FLASHBACK ARCHIVE a a RETENTION 1 DAY", 27},
		{"create_flashback_archive_missing_retention", "CREATE FLASHBACK ARCHIVE a TABLESPACE a 1 DAY", 40},
		{"create_flashback_archive_trailing_token", "CREATE FLASHBACK ARCHIVE a TABLESPACE a RETENTION 1 DAY zzz", 56},
		{"merge_into_missing_using", "MERGE INTO a", 12},
		{"merge_into_unparenthesized_on", "MERGE INTO a USING b ON a.id = b.id WHEN MATCHED THEN UPDATE SET a.x = b.x", 24},

		// 5. Core SQL, DDL, DML, and PL/SQL syntax errors.
		{"select_missing_expr", "SELECT FROM dual", 7},
		{"select_trailing_from", "SELECT 1 FROM", 13},
		{"select_duplicate_from", "SELECT 1 FROM FROM", 14},
		{"select_unclosed_paren", "SELECT (1 FROM dual", 10},
		{"select_extra_alias", "SELECT 1 FROM dual a b", 21},
		{"create_table_missing_lparen", "CREATE TABLE a NO INMEMORY", 15},
		{"alter_table_missing_action", "ALTER TABLE a", 13},
		{"alter_table_invalid_option", "ALTER TABLE a zzz", 14},
		{"insert_missing_into", "INSERT a VALUES (1)", 7},
		{"update_missing_set", "UPDATE a WHERE 1 = 1", 9},
		{"delete_trailing_token", "DELETE FROM a zzz yyy", 18},
		{"plsql_block_missing_semicolon", "BEGIN NULL END;", 11},
		{"bad_statement_start", "ZZZ", 0},

		// 6. Multi-line and multi-byte UTF-8 byte offsets.
		{"multiline_error", "SELECT 1\n  FROM FROM", 16},
		{"multibyte_2byte_utf8", "SELECT 'é' FROM FROM", 17},
		{"multibyte_4byte_utf8", "SELECT '😀' FROM FROM", 19},

		// 7. Lexer errors.
		{"lexer_unterminated_string", "SELECT 'abc FROM dual", 7},
		{"lexer_unterminated_q_quote", "SELECT q'[abc' FROM dual", 7},
		{"lexer_unterminated_quoted_ident", `SELECT "abc FROM dual`, 7},
		{"lexer_empty_quoted_ident", `SELECT "" FROM dual`, 7},
		{"lexer_unterminated_comment", "SELECT 1 /* open FROM dual", 9},
		{"lexer_unexpected_char", "SELECT 1 FROM dual ~", 19},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatalf("Parse(%q) = nil, want ParseError at byte %d", tc.sql, tc.pos)
			}
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("Parse(%q) error type %T does not unwrap to *ParseError", tc.sql, err)
			}
			if pe.Position != tc.pos {
				t.Errorf("Parse(%q) Position = %d, want %d (message=%q)", tc.sql, pe.Position, tc.pos, pe.Message)
			}
			if pe.Message == "" {
				t.Errorf("Parse(%q) returned empty ParseError.Message", tc.sql)
			}
		})
	}
}
