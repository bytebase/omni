package parser

import "testing"

// alterDropFormsOracleCases run against the tables assertParseCodes seeds:
// t (id, n, pid) with check ck_n, index i, and foreign key f, and pt with
// partitions p0, p1, p2.
var alterDropFormsOracleCases = []parseCodeCase{
	// opt_restrict follows DROP [COLUMN] [IF EXISTS] field_ident and nothing
	// else.
	{"ALTER TABLE t DROP COLUMN n RESTRICT", 0},
	{"ALTER TABLE t DROP COLUMN n CASCADE", 0},
	{"ALTER TABLE t DROP n RESTRICT", 0},
	{"ALTER TABLE t DROP pid CASCADE", 0},
	{"ALTER TABLE t DROP COLUMN IF EXISTS n RESTRICT", 0},
	{"ALTER TABLE t DROP IF EXISTS no_such CASCADE", 0},
	{"ALTER TABLE t DROP COLUMN n RESTRICT, DROP COLUMN pid CASCADE", 0},
	{"ALTER TABLE t DROP COLUMN n CASCADE, ADD COLUMN z INT", 0},
	{"ALTER TABLE t DROP COLUMN n /* c */ RESTRICT", 0},
	{"ALTER TABLE t DROP COLUMN `restrict`", 0},
	{"ALTER TABLE t DROP `cascade`", 0},
	{"ALTER TABLE t DROP COLUMN n RESTRICT CASCADE", 1064},
	{"ALTER TABLE t DROP COLUMN n RESTRICT RESTRICT", 1064},
	{"ALTER TABLE t DROP COLUMN n RESTRICT FIRST", 1064},
	{"ALTER TABLE t DROP COLUMN n, RESTRICT", 1064},
	{"ALTER TABLE t DROP COLUMN restrict", 1064},
	{"ALTER TABLE t DROP cascade", 1064},
	{"ALTER TABLE t DROP INDEX i RESTRICT", 1064},
	{"ALTER TABLE t DROP KEY i CASCADE", 1064},
	{"ALTER TABLE t DROP KEY IF EXISTS i CASCADE", 1064},
	{"ALTER TABLE t DROP FOREIGN KEY f RESTRICT", 1064},
	{"ALTER TABLE t DROP FOREIGN KEY IF EXISTS f RESTRICT", 1064},
	{"ALTER TABLE t DROP CONSTRAINT ck_n CASCADE", 1064},
	{"ALTER TABLE t DROP CONSTRAINT IF EXISTS ck_n RESTRICT", 1064},
	{"ALTER TABLE t DROP CONSTRAINT PRIMARY KEY RESTRICT", 1064},
	{"ALTER TABLE t DROP PRIMARY KEY CASCADE", 1064},
	{"ALTER TABLE pt DROP PARTITION p0 RESTRICT", 1064},
	// KEY is required after PRIMARY and FOREIGN.
	{"ALTER TABLE t DROP PRIMARY KEY", 0},
	{"ALTER TABLE t DROP FOREIGN KEY f", 0},
	{"ALTER TABLE t DROP PRIMARY KEY, DROP FOREIGN KEY f", 0},
	{"ALTER TABLE t DROP PRIMARY", 1064},
	{"ALTER TABLE t DROP PRIMARY, ADD PRIMARY KEY (id)", 1064},
	{"ALTER TABLE t DROP PRIMARY INDEX", 1064},
	{"ALTER TABLE t DROP FOREIGN f", 1064},
	{"ALTER TABLE t DROP FOREIGN INDEX f", 1064},
	{"ALTER TABLE t DROP FOREIGN IF EXISTS f", 1064},
}

// TestMariaDBAlterDropFormsOracle asserts omni agrees with a live
// mariadb:11.8.8 on the ALTER TABLE ... DROP forms in
// alterDropFormsOracleCases (container test).
func TestMariaDBAlterDropFormsOracle(t *testing.T) {
	assertParseCodes(t, alterDropFormsOracleCases)
}
