package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// These tests pin what the walk does with the ALTER TABLE clauses this parser
// newly reads. A clause the walk does not know falls to the default arm of
// alterCmd, which abandons the statement and leaves every later one working
// from an Unknown table, so the rules downstream stop seeing the facts they
// need. The cases below assert the opposite: a clause with no name of its own
// is understood and silent, and the clauses that do carry a name are followed.
//
// Every expectation was measured against pingcap/tidb:v8.5.5 on 2026-11-04.

// TestWalkAlterTableClausesAreSilent: the TiDB actions that name nothing the
// walk keeps must raise no finding, and must not stop the walk following the
// actions beside them.
func TestWalkAlterTableClausesAreSilent(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	for _, sql := range []string{
		"ALTER TABLE t CACHE;",
		"ALTER TABLE t NOCACHE;",
		"ALTER TABLE t CACHE, NOCACHE;",
		"ALTER TABLE t READ ONLY;",
		"ALTER TABLE t READ WRITE;",
		"ALTER TABLE t COMPACT;",
		"ALTER TABLE t COMPACT TIFLASH REPLICA;",
		"ALTER TABLE t COMPACT PARTITION p1, p2;",
		"ALTER TABLE t ATTRIBUTES = 'str1,str2';",
		"ALTER TABLE t ATTRIBUTES = DEFAULT;",
		"ALTER TABLE t STATS_OPTIONS = 'str';",
		"ALTER TABLE t ANALYZE PARTITION p0;",
		"ALTER TABLE t ANALYZE PARTITION p0 INDEX i WITH 4 BUCKETS;",
		"ALTER TABLE t PARTITION p0 ATTRIBUTES = 'str';",
		"ALTER TABLE t PARTITION p0 ENGINE = InnoDB MAX_ROWS 5;",
		"ALTER TABLE t LAST PARTITION LESS THAN (1000);",
		"ALTER TABLE t FIRST PARTITION LESS THAN (1);",
		"ALTER TABLE t SPLIT MAXVALUE PARTITION LESS THAN (1000);",
		"ALTER TABLE t MERGE FIRST PARTITION LESS THAN (1000);",
		"ALTER TABLE t FORCE AUTO_INCREMENT = 3;",
		"ALTER TABLE t ATTRIBUTES = 'x', ADD COLUMN c int;",
		"ALTER TABLE t CACHE, ADD COLUMN c int, READ ONLY;",
	} {
		t.Run(sql, func(t *testing.T) { noFinding(t, sql, shop) })
	}
}

// TestWalkAlterTableClausesStillFollowNames: a clause beside one of the silent
// ones must still be followed, so the statement after it sees what it changed.
// A reader that abandoned the statement would let the second statement pass.
func TestWalkAlterTableClausesStillFollowNames(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	oneFinding(t, "ALTER TABLE t CACHE, ADD COLUMN c int;\nALTER TABLE t ADD COLUMN c int;",
		1, "Column `c` already exists in `t`.", shop)
	oneFinding(t, "ALTER TABLE t ATTRIBUTES = 'x', DROP COLUMN zz;",
		0, "Column `zz` doesn't exist in `t`.", shop)
	oneFinding(t, "ALTER TABLE t READ ONLY, ADD INDEX i (a);",
		0, "Index `i` already exists in `t`.", shop)
}

// TestWalkAlterTableAddTableElementList: the parenthesized ADD carries
// constraints beside the columns, and the walk must follow both. An index of a
// name the table already has is 1061 whichever form adds it.
func TestWalkAlterTableAddTableElementList(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	oneFinding(t, "ALTER TABLE t ADD COLUMN (c int, INDEX i (c));",
		0, "Index `i` already exists in `t`.", shop)
	oneFinding(t, "ALTER TABLE t ADD (c int, PRIMARY KEY (c));",
		0, "`t` already has a primary key.", shop)
	// The columns of the list are followed too.
	oneFinding(t, "ALTER TABLE t ADD COLUMN (c int, d int);\nALTER TABLE t ADD COLUMN d int;",
		1, "Column `d` already exists in `t`.", shop)
	noFinding(t, "ALTER TABLE t ADD COLUMN (c int, UNIQUE KEY u (c));", shop)
}

// TestWalkAlterTableIfExistsClauses: IF NOT EXISTS and IF EXISTS turn the
// server's error into a note, so the walk must not report one, and must leave
// the table as the server leaves it.
func TestWalkAlterTableIfExistsClauses(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}

	// A duplicate column with IF NOT EXISTS is a note, not 1060.
	noFinding(t, "ALTER TABLE t ADD COLUMN IF NOT EXISTS a int;", shop)
	// And the column stays the one that was there, so a later DROP finds it.
	noFinding(t, "ALTER TABLE t ADD COLUMN IF NOT EXISTS a int;\nALTER TABLE t DROP COLUMN a;", shop)
	// A column that is not there is still added.
	oneFinding(t, "ALTER TABLE t ADD COLUMN IF NOT EXISTS c int;\nALTER TABLE t ADD COLUMN c int;",
		1, "Column `c` already exists in `t`.", shop)

	// A missing column with IF EXISTS is a note, not 1054.
	noFinding(t, "ALTER TABLE t MODIFY COLUMN IF EXISTS zz varchar(5);", shop)
	noFinding(t, "ALTER TABLE t CHANGE COLUMN IF EXISTS zz yy varchar(5);", shop)
	// CHANGE COLUMN IF EXISTS of a missing column renames nothing, so the new
	// name is not there afterwards.
	oneFinding(t, "ALTER TABLE t CHANGE COLUMN IF EXISTS zz yy varchar(5);\nALTER TABLE t DROP COLUMN yy;",
		1, "Column `yy` doesn't exist in `t`.", shop)
	// A column that is there is still changed.
	oneFinding(t, "ALTER TABLE t CHANGE COLUMN IF EXISTS a c int;\nALTER TABLE t DROP COLUMN a;",
		1, "Column `a` doesn't exist in `t`.", shop)

	// A duplicate index name with IF NOT EXISTS is a note, not 1061, and the
	// index stays the one that was there.
	noFinding(t, "ALTER TABLE t ADD INDEX IF NOT EXISTS i (a);", shop)
	noFinding(t, "ALTER TABLE t ADD KEY IF NOT EXISTS i (a);\nALTER TABLE t DROP INDEX i;", shop)
	// Without it the same statement is reported.
	oneFinding(t, "ALTER TABLE t ADD INDEX i (a);", 0, "Index `i` already exists in `t`.", shop)
	// A name that is not there is still added.
	oneFinding(t, "ALTER TABLE t ADD INDEX IF NOT EXISTS j (a);\nALTER TABLE t ADD INDEX j (a);",
		1, "Index `j` already exists in `t`.", shop)

	// A missing index with IF EXISTS is a note, not 1091.
	noFinding(t, "ALTER TABLE t DROP INDEX IF EXISTS nosuch;", shop)
	noFinding(t, "ALTER TABLE t DROP KEY IF EXISTS nosuch;", shop)
	// Without it the same statement is reported.
	oneFinding(t, "ALTER TABLE t DROP INDEX nosuch;", 0, "Index `nosuch` doesn't exist in `t`.", shop)
	// An index that is there still goes.
	oneFinding(t, "ALTER TABLE t DROP INDEX IF EXISTS i;\nALTER TABLE t DROP INDEX i;",
		1, "Index `i` doesn't exist in `t`.", shop)
}

// TestWalkAlterTableRenameQualified: the walk must follow a RENAME whose
// destination carries a database, and must not treat "db.t" as one name.
func TestWalkAlterTableRenameQualified(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	// Renaming onto a table that is already there is 1050, whichever
	// separator was written.
	for _, sql := range []string{
		"ALTER TABLE t RENAME TO shop.r;",
		"ALTER TABLE t RENAME shop.r;",
		"ALTER TABLE t RENAME = shop.r;",
		"ALTER TABLE t RENAME AS shop.r;",
	} {
		oneFinding(t, sql, 0, "Table `shop.r` already exists.", shop)
	}
	// A free name in the same database is fine, and the table is there after.
	noFinding(t, "ALTER TABLE t RENAME TO shop.t2;\nALTER TABLE shop.t2 ADD COLUMN c int;", shop)
	// The old name is gone.
	oneFinding(t, "ALTER TABLE t RENAME = shop.t2;\nALTER TABLE t ADD COLUMN c int;",
		1, "Table `t` doesn't exist.", shop)
}

// TestAlterTableClauseRules pins what the rules read from the clauses this
// change taught the parser. Three of them read a renamed table's new name and
// two read the keys an action adds, and both of those are shapes this family
// widened: the new name may carry a database, and `ADD [COLUMN] ( ... )`
// carries constraints beside the columns.
func TestAlterTableClauseRules(t *testing.T) {
	// A qualified RENAME is written as a table reference, not as one
	// identifier: reading NewName alone printed `db1.t1` as a single
	// backquoted name.
	runCases(t, review.DisallowRename, []ruleCase{
		{
			name: "a qualified new name",
			sql:  "ALTER TABLE t RENAME TO shop.t2;",
			want: []finding{{0, "RENAME TO shop.t2", "Renames table `t` to `shop.t2`.", []int{0}}},
		},
		{
			name: "an unqualified new name",
			sql:  "ALTER TABLE t RENAME = t2;",
			want: []finding{{0, "RENAME = t2", "Renames table `t` to `t2`.", []int{0}}},
		},
	})

	// A primary key added in the paren form settles the rule, as one added
	// with ADD PRIMARY KEY does.
	runCases(t, review.RequirePrimaryKey, []ruleCase{
		{
			name: "a primary key in the table element list",
			sql:  "CREATE TABLE a (id INT);\nALTER TABLE a ADD COLUMN (note INT, PRIMARY KEY (id));",
		},
		{
			name: "no primary key in the table element list",
			sql:  "CREATE TABLE a (id INT);\nALTER TABLE a ADD COLUMN (note INT, KEY k (note));",
			want: []finding{{0, "CREATE TABLE a (id INT)", "Creates table `a` without a primary key.", []int{0}}},
		},
	})

	// A key added in the paren form is a key the change made, and reads the
	// same as one added with ADD KEY: a plain key dropped again is nothing,
	// and a unique one is reported, which is what the rule already did for
	// `ADD UNIQUE KEY u (b)`.
	runCases(t, review.DisallowDropConstraint, []ruleCase{
		{
			name: "a plain key the table element list added",
			sql:  "ALTER TABLE nokey ADD COLUMN (b INT, KEY u (b));\nALTER TABLE nokey DROP INDEX u;",
		},
		{
			name: "a unique key the table element list added reads as ADD UNIQUE KEY does",
			sql:  "ALTER TABLE nokey ADD COLUMN (b INT, UNIQUE KEY u (b));\nALTER TABLE nokey DROP INDEX u;",
			want: []finding{{1, "DROP INDEX u", "Drops unique constraint `u` of `nokey`.", []int{0}}},
		},
	})
}
