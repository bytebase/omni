package parser

import "testing"

// What MariaDB 11.8.8 reads after ALTER [COLUMN] col_name, measured with PREPARE
// on mariadb:11.8.8 on 2026-11-06: a 1064 is the parser's refusal, and any other
// answer means the text parsed. ALTER COLUMN sets or drops a default; it has no
// SET NOT NULL or DROP NOT NULL, and a column's nullability changes through
// MODIFY or CHANGE. It has no SET VISIBLE or SET INVISIBLE either: a column is
// hidden with the INVISIBLE attribute, and MariaDB has no VISIBLE attribute.
//
// Ref: `ALTER opt_column opt_if_exists_table_element field_ident SET DEFAULT
// column_default_expr` and `... field_ident DROP DEFAULT` in MariaDB
// mariadb-11.8.8 sql/sql_yacc.yy, lines 8094 and 8111; `vcol_attribute: ... |
// INVISIBLE_SYM`, line 6384, with no VISIBLE_SYM arm.

type alterColumnFormCase struct {
	sql    string
	accept bool
}

var alterColumnFormCases = []alterColumnFormCase{
	// nullability, which ALTER COLUMN does not change.
	{"ALTER TABLE t ALTER COLUMN b SET NOT NULL", false},
	{"ALTER TABLE t ALTER b SET NOT NULL", false},
	{"ALTER TABLE t ALTER COLUMN b DROP NOT NULL", false},
	{"ALTER TABLE t ALTER b DROP NOT NULL", false},
	{"alter table t alter column b set not null", false},
	{"ALTER TABLE t ALTER COLUMN b SET NULL", false},
	{"ALTER TABLE t ALTER COLUMN b DROP NULL", false},
	{"ALTER TABLE t ALTER COLUMN b SET NOT", false},
	{"ALTER TABLE t ALTER COLUMN b DROP NOT", false},
	{"ALTER TABLE t ALTER COLUMN t.b SET NOT NULL", false},
	{"ALTER TABLE t ALTER COLUMN t.b DROP NOT NULL", false},
	{"ALTER TABLE t ALTER COLUMN b SET DEFAULT 1, ALTER COLUMN b SET NOT NULL", false},
	{"ALTER TABLE t ALTER COLUMN b DROP NOT NULL, ALTER COLUMN b DROP DEFAULT", false},

	// a column name with no action, or an action with nothing after it.
	{"ALTER TABLE t ALTER COLUMN b SET", false},
	{"ALTER TABLE t ALTER COLUMN b DROP", false},
	{"ALTER TABLE t ALTER COLUMN b", false},
	{"ALTER TABLE t ALTER b", false},

	// the actions ALTER COLUMN has.
	{"ALTER TABLE t ALTER COLUMN b SET DEFAULT 1", true},
	{"ALTER TABLE t ALTER COLUMN b SET DEFAULT NULL", true},
	{"ALTER TABLE t ALTER b DROP DEFAULT", true},
	{"ALTER TABLE t ALTER b DROP DEFAULT, ALTER b SET DEFAULT 2", true},

	// how a column's nullability does change.
	{"ALTER TABLE t MODIFY b INT NOT NULL", true},
	{"ALTER TABLE t MODIFY COLUMN b INT NULL", true},
	{"ALTER TABLE t CHANGE b b INT NOT NULL", true},

	// visibility, which ALTER COLUMN does not change either.
	{"ALTER TABLE t ALTER COLUMN b SET VISIBLE", false},
	{"ALTER TABLE t ALTER b SET VISIBLE", false},
	{"ALTER TABLE t ALTER COLUMN b SET INVISIBLE", false},
	{"ALTER TABLE t ALTER b SET INVISIBLE", false},
	{"alter table t alter column b set invisible", false},
	{"ALTER TABLE t ALTER COLUMN b SET DEFAULT 1, ALTER COLUMN b SET INVISIBLE", false},

	// how a column's visibility does change: INVISIBLE, never VISIBLE.
	{"CREATE TABLE t2 (a INT, b INT INVISIBLE)", true},
	{"CREATE TABLE t2 (a INT, b INT VISIBLE)", false},
	{"ALTER TABLE t ADD COLUMN c INT INVISIBLE", true},
	{"ALTER TABLE t ADD COLUMN c INT VISIBLE", false},
	{"ALTER TABLE t MODIFY b INT INVISIBLE", true},
	{"ALTER TABLE t MODIFY b INT VISIBLE", false},
	{"ALTER TABLE t CHANGE b b INT INVISIBLE", true},
	{"ALTER TABLE t CHANGE b b INT VISIBLE", false},
}

// TestAlterColumnForms holds the parser to the recorded verdicts without the
// server, so a regression fails however the suite is run.
func TestAlterColumnForms(t *testing.T) {
	for _, c := range alterColumnFormCases {
		t.Run(c.sql, func(t *testing.T) {
			omniAccepts, perr := parsesAsOneStatement(c.sql)
			if omniAccepts != c.accept {
				t.Errorf("omni accepts=%v, the record says %v for %q: %v", omniAccepts, c.accept, c.sql, perr)
			}
		})
	}
}
