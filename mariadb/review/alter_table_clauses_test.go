package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// TestWalkAlterTableAddTableElementList: the parenthesized ADD carries
// constraints beside the columns, and the walk must follow both. A reader that
// took only the columns would let a duplicate index name through.
func TestWalkAlterTableAddTableElementList(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	oneFinding(t, "ALTER TABLE t ADD COLUMN (c int, INDEX i (c));",
		0, "Index `i` already exists in `t`.", shop)
	oneFinding(t, "ALTER TABLE t ADD (c int, PRIMARY KEY (c));",
		0, "`t` already has a primary key.", shop)
	oneFinding(t, "ALTER TABLE t ADD COLUMN (c int, d int);\nALTER TABLE t ADD COLUMN d int;",
		1, "Column `d` already exists in `t`.", shop)
	noFinding(t, "ALTER TABLE t ADD COLUMN (c int, UNIQUE KEY u (c));", shop)
}

// TestAlterTableClauseRules pins what the rules read from the keys the
// parenthesized ADD now carries: a primary key written beside the columns
// settles the primary-key rule, and a key written there is one the change
// made rather than one the schema had.
func TestAlterTableClauseRules(t *testing.T) {
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
