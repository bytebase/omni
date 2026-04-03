package parser

import (
	"strings"
	"testing"
)

// TestSplitMysqldumpOutput tests against a realistic mysqldump-style output
// that bytebase users would import.
func TestSplitMysqldumpOutput(t *testing.T) {
	sql := `SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;

-- Table structure for employees
CREATE TABLE employees (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Procedure structure for citycount
DELIMITER ;;
CREATE DEFINER=` + "`root`@`%`" + ` PROCEDURE citycount(IN country CHAR(3), OUT cities INT)
BEGIN
    SELECT COUNT(*) INTO cities FROM city WHERE CountryCode = country;
END ;;
DELIMITER ;

-- Function structure for CalcIncome
DELIMITER ;;
CREATE DEFINER=` + "`root`@`%`" + ` FUNCTION CalcIncome(starting_value INT) RETURNS int
    DETERMINISTIC
BEGIN
   DECLARE income INT;
   SET income = 0;
   label1: WHILE income <= 3000 DO
     SET income = income + starting_value;
   END WHILE label1;
   RETURN income;
END ;;
DELIMITER ;

-- Trigger structure for trg_audit
DELIMITER ;;
CREATE TRIGGER trg_audit BEFORE INSERT ON employees
FOR EACH ROW
BEGIN
    SET NEW.created_at = NOW();
END ;;
DELIMITER ;

SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;`

	segs := Split(sql)

	// We expect: 2 SETs, 1 CREATE TABLE, 1 PROCEDURE, 1 FUNCTION, 1 TRIGGER, 2 SETs = 8 statements.
	wantCount := 8
	if len(segs) != wantCount {
		var texts []string
		for _, s := range segs {
			texts = append(texts, s.Text[:min(60, len(s.Text))]+"...")
		}
		t.Fatalf("got %d segments, want %d:\n%s", len(segs), wantCount, strings.Join(texts, "\n"))
	}

	// Verify key segments.
	assertSegContains(t, segs[0], "SET @OLD_UNIQUE_CHECKS")
	assertSegContains(t, segs[2], "CREATE TABLE employees")
	assertSegContains(t, segs[3], "PROCEDURE citycount")
	assertSegContains(t, segs[3], "SELECT COUNT(*) INTO cities")
	assertSegContains(t, segs[4], "FUNCTION CalcIncome")
	assertSegContains(t, segs[4], "END WHILE label1")
	assertSegContains(t, segs[5], "TRIGGER trg_audit")
	assertSegContains(t, segs[6], "SET FOREIGN_KEY_CHECKS")

	// Verify compound blocks are intact — no inner semicolons caused extra splits.
	assertSegNotContains(t, segs[3], "DELIMITER")
	assertSegNotContains(t, segs[4], "DELIMITER")
	assertSegNotContains(t, segs[5], "DELIMITER")
}

// TestSplitKeywordsTrappedInStrings verifies that keywords inside string
// literals and identifiers don't affect block depth or delimiter detection.
func TestSplitKeywordsTrappedInStrings(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want int // expected segment count
	}{
		{
			name: "BEGIN END in single-quoted string",
			sql:  "INSERT INTO t VALUES ('BEGIN; END;'); SELECT 1;",
			want: 2,
		},
		{
			name: "BEGIN END in double-quoted string",
			sql:  `INSERT INTO t VALUES ("BEGIN; END;"); SELECT 1;`,
			want: 2,
		},
		{
			name: "BEGIN in backtick identifier",
			sql:  "SELECT * FROM `BEGIN`; SELECT 1;",
			want: 2,
		},
		{
			name: "END in backtick identifier",
			sql:  "SELECT `END` FROM t; SELECT 1;",
			want: 2,
		},
		{
			name: "DELIMITER in single-quoted string",
			sql:  "INSERT INTO t VALUES ('DELIMITER ;;'); SELECT 1;",
			want: 2,
		},
		{
			name: "DELIMITER in block comment",
			sql:  "/* DELIMITER ;; */\nSELECT 1; SELECT 2;",
			want: 2,
		},
		{
			name: "DELIMITER in line comment",
			sql:  "-- DELIMITER ;;\nSELECT 1; SELECT 2;",
			want: 2,
		},
		{
			name: "DELIMITER in hash comment",
			sql:  "# DELIMITER ;;\nSELECT 1; SELECT 2;",
			want: 2,
		},
		{
			name: "CASE WHEN in string",
			sql:  "INSERT INTO t VALUES ('CASE WHEN x THEN END'); SELECT 1;",
			want: 2,
		},
		{
			name: "IF in string",
			sql:  "SELECT 'IF x THEN y; END IF;'; SELECT 1;",
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segs := Split(tt.sql)
			if len(segs) != tt.want {
				t.Errorf("got %d segments, want %d\nsegs: %q", len(segs), tt.want, segTexts(segs))
			}
		})
	}
}

// TestSplitCaseExpression verifies that CASE expressions in SELECT don't
// interfere with statement splitting (they balance depth correctly).
func TestSplitCaseExpression(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want int
	}{
		{
			name: "simple CASE expression",
			sql:  "SELECT CASE WHEN x > 0 THEN 1 ELSE 0 END FROM t; SELECT 2;",
			want: 2,
		},
		{
			name: "nested CASE expressions",
			sql:  "SELECT CASE WHEN a THEN CASE WHEN b THEN 1 ELSE 2 END ELSE 3 END FROM t; SELECT 1;",
			want: 2,
		},
		{
			name: "CASE with semicolons in compound block",
			sql:  "CREATE PROCEDURE p() BEGIN CASE x WHEN 1 THEN SET y = 1; WHEN 2 THEN SET y = 2; END CASE; END; SELECT 1;",
			want: 2,
		},
		{
			name: "CASE expression outside compound block",
			sql:  "SELECT CASE status WHEN 'A' THEN 'Active' WHEN 'I' THEN 'Inactive' END AS label FROM t; INSERT INTO t2 SELECT 1;",
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segs := Split(tt.sql)
			if len(segs) != tt.want {
				t.Errorf("got %d segments, want %d\nsegs: %q", len(segs), tt.want, segTexts(segs))
			}
		})
	}
}

// TestSplitBackslashEscapeEdgeCases tests tricky backslash escape scenarios
// that could cause the quote skipper to lose track of string boundaries.
func TestSplitBackslashEscapeEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want int
	}{
		{
			name: "backslash-escaped quote in string",
			sql:  `SELECT '\''; SELECT 1;`,
			want: 2,
		},
		{
			name: "double backslash then quote",
			sql:  `SELECT '\\'; SELECT 1;`,
			want: 2,
		},
		{
			name: "backslash-escaped quote then semicolon in string",
			sql:  `SELECT '\';'; SELECT 1;`,
			want: 2,
		},
		{
			name: "backslash at end of string",
			sql:  `SELECT 'path\\'; SELECT 1;`,
			want: 2,
		},
		{
			name: "mixed escapes",
			sql:  `INSERT INTO t VALUES ('it''s a \"test\"'); SELECT 1;`,
			want: 2,
		},
		{
			name: "escaped double quote in double-quoted string",
			sql:  `SELECT "say \"hello\""; SELECT 1;`,
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segs := Split(tt.sql)
			if len(segs) != tt.want {
				t.Errorf("got %d segments, want %d\nsegs: %q", len(segs), tt.want, segTexts(segs))
			}
		})
	}
}

// TestSplitLabeledBlocks tests labeled compound blocks which are common in
// real stored procedures.
func TestSplitLabeledBlocks(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want []string
	}{
		{
			name: "labeled BEGIN END",
			sql:  "CREATE PROCEDURE p() myblock: BEGIN SELECT 1; END myblock; SELECT 2;",
			want: []string{
				"CREATE PROCEDURE p() myblock: BEGIN SELECT 1; END myblock",
				" SELECT 2",
			},
		},
		{
			name: "labeled WHILE",
			sql:  "CREATE PROCEDURE p() BEGIN lbl: WHILE TRUE DO LEAVE lbl; END WHILE lbl; END; SELECT 1;",
			want: []string{
				"CREATE PROCEDURE p() BEGIN lbl: WHILE TRUE DO LEAVE lbl; END WHILE lbl; END",
				" SELECT 1",
			},
		},
		{
			name: "labeled LOOP",
			sql:  "CREATE PROCEDURE p() BEGIN myloop: LOOP SELECT 1; LEAVE myloop; END LOOP myloop; END; SELECT 2;",
			want: []string{
				"CREATE PROCEDURE p() BEGIN myloop: LOOP SELECT 1; LEAVE myloop; END LOOP myloop; END",
				" SELECT 2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segs := Split(tt.sql)
			got := segTexts(segs)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d segments, want %d\n  got:  %q\n  want: %q", len(got), len(tt.want), got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("seg[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// TestSplitEmptyAndMinimalBodies tests edge cases with empty or minimal
// procedure/function bodies.
func TestSplitEmptyAndMinimalBodies(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want int
	}{
		{
			name: "empty BEGIN END",
			sql:  "CREATE PROCEDURE p() BEGIN END; SELECT 1;",
			want: 2,
		},
		{
			name: "BEGIN with only SET",
			sql:  "CREATE PROCEDURE p() BEGIN SET @x = 1; END; SELECT 1;",
			want: 2,
		},
		{
			name: "multiple empty procedures with DELIMITER",
			sql:  "DELIMITER ;;\nCREATE PROCEDURE p1() BEGIN END;;\nCREATE PROCEDURE p2() BEGIN END;;\nDELIMITER ;\nSELECT 1;",
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segs := Split(tt.sql)
			if len(segs) != tt.want {
				t.Errorf("got %d segments, want %d\nsegs: %q", len(segs), tt.want, segTexts(segs))
			}
		})
	}
}

// TestSplitDelimiterEdgeCases tests unusual but valid DELIMITER usage.
func TestSplitDelimiterEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want int
	}{
		{
			name: "DELIMITER at EOF without newline",
			sql:  "SELECT 1;\nDELIMITER ;;",
			want: 1,
		},
		{
			name: "DELIMITER changed twice",
			sql:  "DELIMITER ;;\nSELECT 1;;\nDELIMITER $$\nSELECT 2$$\nDELIMITER ;\nSELECT 3;",
			want: 3,
		},
		{
			name: "DELIMITER with CRLF",
			sql:  "DELIMITER ;;\r\nSELECT 1;;\r\nDELIMITER ;\r\nSELECT 2;",
			want: 2,
		},
		{
			name: "DELIMITER with tabs around delimiter string",
			sql:  "DELIMITER\t;;\nSELECT 1;;\nDELIMITER\t;\nSELECT 2;",
			want: 2,
		},
		{
			name: "multiple procedures sequential DELIMITER blocks",
			sql: `DELIMITER ;;
CREATE PROCEDURE p1() BEGIN SELECT 1; END;;
CREATE PROCEDURE p2() BEGIN SELECT 2; END;;
DELIMITER ;
SELECT 3;`,
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segs := Split(tt.sql)
			if len(segs) != tt.want {
				t.Errorf("got %d segments, want %d\nsegs: %q", len(segs), tt.want, segTexts(segs))
			}
		})
	}
}

// TestSplitConditionalComments tests MySQL version-conditional comments
// which are common in mysqldump output.
// Conditional comments /*!...*/ contain executable SQL and must NOT be
// treated as empty. Split preserves them as segments; the Lexer/Parser
// expands their content downstream.
func TestSplitConditionalComments(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want int
	}{
		{
			name: "conditional comment standalone is not empty",
			sql:  "/*!50000 SET @a = 1; SET @b = 2; */; SELECT 1;",
			want: 2,
		},
		{
			name: "conditional comment as part of CREATE TABLE",
			sql:  "CREATE TABLE t (\n  id INT\n) /*!50100 ENGINE=InnoDB */; SELECT 1;",
			want: 2,
		},
		{
			name: "mysqldump SET with conditional is not empty",
			sql:  "/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */; SELECT 1;",
			want: 2,
		},
		{
			name: "semicolons inside conditional comment dont split",
			sql:  "SELECT /*!50000 1, 2; */ 3; SELECT 4;",
			want: 2,
		},
		{
			name: "regular block comment is still empty",
			sql:  "/* just a comment */; SELECT 1;",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segs := Split(tt.sql)
			if len(segs) != tt.want {
				t.Errorf("got %d segments, want %d\nsegs: %q", len(segs), tt.want, segTexts(segs))
			}
		})
	}
}

// TestSplitByteOffsetIdentity verifies that for every segment,
// sql[seg.ByteStart:seg.ByteEnd] == seg.Text. This is critical for
// downstream consumers that need to map positions back to the original SQL.
func TestSplitByteOffsetIdentity(t *testing.T) {
	sqls := []string{
		"SELECT 1; SELECT 2; SELECT 3;",
		"DELIMITER ;;\nSELECT 1;;\nDELIMITER ;\nSELECT 2;",
		"CREATE PROCEDURE p() BEGIN SELECT 1; END; SELECT 2;",
		"SELECT 'hello'; -- comment\nSELECT 2;",
		"SELECT 1;\r\nSELECT 2;\r\nSELECT 3;",
		"SELECT * FROM 表名; INSERT INTO 表 VALUES (1);",
	}
	for _, sql := range sqls {
		segs := Split(sql)
		for i, s := range segs {
			got := sql[s.ByteStart:s.ByteEnd]
			if got != s.Text {
				t.Errorf("sql=%q seg[%d]: sql[%d:%d] = %q, but Text = %q",
					sql, i, s.ByteStart, s.ByteEnd, got, s.Text)
			}
		}
	}
}

// TestSplitColumnNamedBegin verifies that keywords used as identifiers
// (without backticks) in non-compound contexts don't confuse the splitter.
func TestSplitKeywordAsIdentifier(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want int
	}{
		{
			name: "column named begin in CREATE TABLE",
			// In practice MySQL doesn't allow unquoted BEGIN as column name,
			// but backtick-quoted is fine and common.
			sql:  "CREATE TABLE t (`begin` DATE, `end` DATE); SELECT 1;",
			want: 2,
		},
		{
			name: "table named loop",
			sql:  "SELECT * FROM `loop`; SELECT 1;",
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segs := Split(tt.sql)
			if len(segs) != tt.want {
				t.Errorf("got %d segments, want %d\nsegs: %q", len(segs), tt.want, segTexts(segs))
			}
		})
	}
}

// TestSplitMigrationScript tests a realistic migration script that a user
// might paste into bytebase's SQL review.
func TestSplitMigrationScript(t *testing.T) {
	sql := `-- Migration: add audit logging
ALTER TABLE users ADD COLUMN updated_at DATETIME;

CREATE INDEX idx_users_updated ON users(updated_at);

DELIMITER ;;
CREATE TRIGGER trg_users_update
BEFORE UPDATE ON users
FOR EACH ROW
BEGIN
    SET NEW.updated_at = NOW();
    INSERT INTO audit_log (table_name, row_id, action)
    VALUES ('users', NEW.id, 'UPDATE');
END;;
DELIMITER ;

-- Backfill existing rows
UPDATE users SET updated_at = created_at WHERE updated_at IS NULL;`

	segs := Split(sql)

	wantCount := 4
	if len(segs) != wantCount {
		t.Fatalf("got %d segments, want %d\nsegs: %q", len(segs), wantCount, segTexts(segs))
	}

	assertSegContains(t, segs[0], "ALTER TABLE users")
	assertSegContains(t, segs[1], "CREATE INDEX idx_users_updated")
	assertSegContains(t, segs[2], "CREATE TRIGGER trg_users_update")
	assertSegContains(t, segs[2], "INSERT INTO audit_log")
	assertSegContains(t, segs[3], "UPDATE users SET updated_at")

	// Trigger body must be intact.
	assertSegNotContains(t, segs[2], "DELIMITER")
}

// --- helpers ---

func assertSegContains(t *testing.T, seg Segment, substr string) {
	t.Helper()
	if !strings.Contains(seg.Text, substr) {
		t.Errorf("segment does not contain %q:\n  %s", substr, truncate(seg.Text, 120))
	}
}

func assertSegNotContains(t *testing.T, seg Segment, substr string) {
	t.Helper()
	if strings.Contains(seg.Text, substr) {
		t.Errorf("segment unexpectedly contains %q:\n  %s", substr, truncate(seg.Text, 120))
	}
}

func segTexts(segs []Segment) []string {
	var out []string
	for _, s := range segs {
		out = append(out, s.Text)
	}
	return out
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
