package review

import (
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/review"
)

// codeCase is a change that gives one finding of rule, which carries
// exactly want.
type codeCase struct {
	rule   review.Rule
	sql    string
	change review.Change
	target review.Target
	want   []review.FindingCode
}

func runCodeCases(t *testing.T, cases []codeCase) {
	t.Helper()
	for _, tt := range cases {
		t.Run(string(tt.want[0])+" "+tt.sql, func(t *testing.T) {
			result, err := Review(context.Background(), tt.sql, review.Options{Rules: []review.Rule{tt.rule}, Change: tt.change}, []review.Target{tt.target})
			if err != nil {
				t.Fatalf("Review: %v", err)
			}
			if len(result.Findings) != 1 {
				t.Fatalf("findings = %+v, want one", result.Findings)
			}
			f := result.Findings[0]
			if f.Rule != tt.rule || !slices.Equal(f.Codes, tt.want) {
				t.Errorf("finding %s %v %q, want %s %v", f.Rule, f.Codes, f.Message, tt.rule, tt.want)
			}
			if err := review.CheckCodes(f); err != nil {
				t.Error(err)
			}
		})
	}
}

// TestFindingCodes pins the code of each condition the package reports.
func TestFindingCodes(t *testing.T) {
	shopTarget := review.Target{Schema: shop(), BackupDatabaseExists: true}
	walkTarget := review.Target{Schema: walkSchema()}
	backup := review.Change{PriorBackup: true, MaxBackupSize: 1 << 20}
	codes := func(c ...review.FindingCode) []review.FindingCode { return c }
	runCodeCases(t, []codeCase{
		{review.Syntax, "SELEC 1;", review.Change{}, shopTarget, codes(review.SyntaxError)},

		{review.WalkThrough, "ALTER TABLE nosuch ADD COLUMN b int;", review.Change{}, walkTarget, codes(review.ObjectNotFound)},
		{review.WalkThrough, "ALTER TABLE t DROP COLUMN zz;", review.Change{}, walkTarget, codes(review.ObjectNotFound)},
		{review.WalkThrough, "CREATE TABLE t (a int);", review.Change{}, walkTarget, codes(review.NameTaken)},
		{review.WalkThrough, "ALTER TABLE t ADD PRIMARY KEY (a);", review.Change{}, walkTarget, codes(review.NameTaken)},
		{review.WalkThrough, "CREATE TABLE x (a int PRIMARY KEY, b int, PRIMARY KEY (b));", review.Change{}, walkTarget, codes(review.NameTaken)},
		{review.WalkThrough, "ALTER TABLE v ADD COLUMN b int;", review.Change{}, walkTarget, codes(review.WrongObjectKind)},
		// The message says the view is no table that exists.
		{review.WalkThrough, "DROP TABLE v;", review.Change{}, walkTarget, codes(review.WrongObjectKind)},
		{review.WalkThrough, "DROP SEQUENCE v;", review.Change{}, walkTarget, codes(review.WrongObjectKind)},
		{review.WalkThrough, "ALTER TABLE t ADD COLUMN b int, ADD COLUMN b int;", review.Change{}, walkTarget, codes(review.DuplicateName)},
		{review.WalkThrough, "CREATE TABLE x (a int, a int);", review.Change{}, walkTarget, codes(review.DuplicateName)},
		{review.WalkThrough, "DROP DATABASE shop;\nDROP SEQUENCE IF EXISTS sq;", review.Change{}, walkTarget, codes(review.NoCurrentSchema)},
		{review.WalkThrough, "CREATE TABLE " + strings.Repeat("x", 65) + " (a int);", review.Change{}, walkTarget, codes(review.NameTooLong)},
		{review.WalkThrough, "CREATE INDEX `PRIMARY` ON nokey (a);", review.Change{}, walkTarget, codes(review.ReservedName)},
		{review.WalkThrough, "RENAME TABLE v TO other.v;", review.Change{}, walkTarget, codes(review.UnsupportedForm)},
		{review.WalkThrough, "ALTER TABLE t COMMENT 'x', RENAME TO t2;", review.Change{}, walkTarget, codes(review.UnsupportedForm)},

		{review.OnlineMigration, "ALTER TABLE plain ADD COLUMN m INT;", review.Change{OnlineMigration: true}, shopTarget, codes(review.OnlineMigrationUnsupported)},

		{review.PriorBackup, "UPDATE t SET n = 1 WHERE id = 1;", review.Change{PriorBackup: true, MaxBackupSize: 10}, shopTarget, codes(review.ChangeTooLarge)},
		{review.PriorBackup, "UPDATE t SET n = 1 WHERE id = 1;", backup, review.Target{Schema: shop()}, codes(review.BackupLocationMissing)},
		{review.PriorBackup, "UPDATE t SET n = 1 WHERE id = 1;\nUPDATE t SET n = 2 WHERE id = 2;\nDELETE FROM shop.t WHERE id = 3;\nUPDATE t SET n = 4 WHERE id = 4;\nUPDATE t SET n = 5 WHERE id = 5;\nUPDATE p SET note = '' WHERE id = 1;", backup, shopTarget, codes(review.StatementsShareBackup)},
		{review.PriorBackup, "UPDATE t SET n = 1 WHERE id = 1;\nUPDATE t SET n = 2 WHERE id = 2;\nWITH x AS (SELECT 3 AS id) UPDATE t SET n = 3 WHERE id IN (SELECT id FROM x);\nUPDATE t SET n = 4 WHERE id = 4;\nUPDATE t SET n = 5 WHERE id = 5;\nUPDATE t SET n = 6 WHERE id = 6;", backup, shopTarget, codes(review.StatementsShareBackup)},
		// Over the limit, with a statement whose tables are not certain, both
		// reasons may stop the backup: one code still covers them.
		{review.PriorBackup, "UPDATE t SET n = 1 WHERE id = 1;\nUPDATE p SET note = '' WHERE id = 1;\nUPDATE plain SET v = 1 WHERE id = 1;\nDELETE FROM nokey WHERE a = 1;\nUPDATE trig SET id = 2 WHERE id = 1;\nWITH x AS (SELECT 1 AS id) UPDATE t JOIN p ON p.id = t.p_id SET t.n = 1 WHERE t.id IN (SELECT id FROM x);", backup, shopTarget, codes(review.StatementsShareBackup)},
		{review.PriorBackup, "DELETE y FROM t AS x WHERE x.id = 2;", backup, shopTarget, codes(review.BackupReadsWrongTable)},
		{review.PriorBackup, "UPDATE pv SET note = 'z' WHERE id = 1;", backup, shopTarget, codes(review.TableNotBackupable)},
		{review.PriorBackup, "UPDATE other.t SET n = 1 WHERE id = 1;", backup, shopTarget, codes(review.TableNotBackupable)},
		{review.PriorBackup, "BATCH ON id LIMIT 2 DELETE FROM t WHERE n > 0;", backup, shopTarget, codes(review.StatementNotBackupable)},
		{review.PriorBackup, "REPLACE INTO p VALUES (3, 'c');", backup, shopTarget, codes(review.StatementNotBackupable)},
		{review.PriorBackup, "EXPLAIN ANALYZE UPDATE t SET n = 1 WHERE id = 1;", backup, shopTarget, codes(review.StatementNotBackupable)},
		{review.PriorBackup, "CREATE TABLE n2 (id INT PRIMARY KEY);\nUPDATE n2 SET id = 2 WHERE id = 1;", backup, shopTarget, codes(review.BackupReadsNewObject)},

		{review.RequireIsNull, "SELECT * FROM t WHERE n = NULL;", review.Change{}, shopTarget, codes(review.NullComparison)},
		{review.RequireWhere, "DELETE FROM t;", review.Change{}, shopTarget, codes(review.MissingWhere)},

		{review.DisallowDropObject, "DROP TABLE t;", review.Change{}, shopTarget, codes(review.DropsData)},
		{review.DisallowDropObject, "DROP DATABASE shop;", review.Change{}, shopTarget, codes(review.DropsData)},
		{review.DisallowDropObject, "DROP SEQUENCE sq;", review.Change{}, shopTarget, codes(review.DropsData)},
		{review.DisallowDropObject, "ALTER TABLE t DROP COLUMN n;", review.Change{}, shopTarget, codes(review.DropsData)},
		{review.DisallowDropObject, "ALTER TABLE t DROP PARTITION p0;", review.Change{}, shopTarget, codes(review.DropsData)},
		{review.DisallowDropObject, "DROP VIEW pv;", review.Change{}, shopTarget, codes(review.DropsDefinition)},
		{review.DisallowDropObject, "DROP INDEX idx_n ON t;", review.Change{}, shopTarget, codes(review.DropsDefinition)},
		{review.DisallowDropObject, "ALTER TABLE t DROP INDEX idx_n;", review.Change{}, shopTarget, codes(review.DropsDefinition)},
		{review.DisallowDropObject, "DROP PLACEMENT POLICY pp;", review.Change{}, shopTarget, codes(review.DropsDefinition)},
		{review.DisallowDropObject, "DROP USER u;", review.Change{}, shopTarget, codes(review.DropsAccess)},
		{review.DisallowDropObject, "DROP ROLE r;", review.Change{}, shopTarget, codes(review.DropsAccess)},

		{review.DisallowTruncate, "TRUNCATE TABLE t;", review.Change{}, shopTarget, codes(review.TruncatesData)},
		{review.DisallowTruncate, "ALTER TABLE t TRUNCATE PARTITION p0;", review.Change{}, shopTarget, codes(review.TruncatesData)},
		{review.DisallowDropConstraint, "ALTER TABLE t DROP FOREIGN KEY fk_p;", review.Change{}, shopTarget, codes(review.DropsConstraint)},
		{review.DisallowDropConstraint, "DROP INDEX uk_code ON t;", review.Change{}, shopTarget, codes(review.DropsConstraint)},
		{review.DisallowRename, "RENAME TABLE t TO t2;", review.Change{}, shopTarget, codes(review.RenamesObject)},
		{review.DisallowRename, "ALTER TABLE t RENAME COLUMN n TO m;", review.Change{}, shopTarget, codes(review.RenamesObject)},

		{review.RequirePrimaryKey, "CREATE TABLE x (a INT);", review.Change{}, shopTarget, codes(review.TableWithoutPrimaryKey)},
		{review.RequirePrimaryKey, "ALTER TABLE plain DROP PRIMARY KEY;", review.Change{}, shopTarget, codes(review.PrimaryKeyRemoved)},
	})
}

// TestFindingCodesAcrossTargets checks that findings merge across targets
// only with their codes: DROP TABLE of a missing table and of a view read
// the same, but one names no object and the other one of another kind.
func TestFindingCodesAcrossTargets(t *testing.T) {
	noView := walkSchema()
	noView.Schemas[0].Views = nil
	result, err := Review(context.Background(), "DROP TABLE v;", review.Options{Rules: []review.Rule{review.WalkThrough}},
		[]review.Target{{Schema: noView}, {Schema: walkSchema()}, {Schema: noView}})
	if err != nil {
		t.Fatalf("Review: %v", err)
	}
	type coded struct {
		codes   []review.FindingCode
		message string
		targets []int
	}
	var got []coded
	for _, f := range result.Findings {
		got = append(got, coded{f.Codes, f.Message, f.Targets})
	}
	want := []coded{
		{[]review.FindingCode{review.ObjectNotFound}, "Table v does not exist", []int{0, 2}},
		{[]review.FindingCode{review.WrongObjectKind}, "Table v does not exist", []int{1}},
	}
	if !slices.EqualFunc(got, want, func(a, b coded) bool {
		return slices.Equal(a.codes, b.codes) && a.message == b.message && slices.Equal(a.targets, b.targets)
	}) {
		t.Errorf("findings = %+v, want %+v", got, want)
	}
}
