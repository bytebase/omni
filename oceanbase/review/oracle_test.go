package review

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/tape"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/review"

	"github.com/bytebase/omni/oceanbase/internal/testcontainer"
)

// The oracle tests of this package hold the review to OceanBase. Each replays
// its tape (testdata/tapes) and starts no server; a recording
// (OMNI_RECORD=1) asks the pinned one, which TestMain then removes.

func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.OceanBase.Terminate()
	os.Exit(code)
}

// oracleBase is the schema the rule oracles' cases start from.
const oracleBase = `
CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
  UNIQUE KEY uk_code (code), KEY idx_n (n),
  CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
  CONSTRAINT ck_n CHECK (n > 0));
CREATE TABLE nokey (a INT);
CREATE TABLE plain (id INT PRIMARY KEY, v INT);
INSERT INTO p VALUES (1, 'a'), (2, 'b');
INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
INSERT INTO nokey VALUES (1), (2);
INSERT INTO plain VALUES (1, 1), (2, 2);
CREATE VIEW pv AS SELECT id, note FROM p;
`

// openShop returns t's connection to an OceanBase that holds only the shop
// database, made from oracleBase and selected.
func openShop(t *testing.T) *sql.DB {
	t.Helper()
	emptyOceanBase(t)
	db := testcontainer.Open(t)
	mustExec(t, db, "CREATE DATABASE shop", "USE shop")
	mustExec(t, db, splitWalk(oracleBase)...)
	return db
}

// read runs a query and returns its rows as text, NULL as "".
func read(t testing.TB, db *sql.DB, query string, args ...any) [][]string {
	t.Helper()
	rows, err := db.QueryContext(context.Background(), query, args...)
	if err != nil {
		t.Fatalf("%s: %v", query, err)
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		t.Fatal(err)
	}
	var out [][]string
	for rows.Next() {
		vals := make([]sql.NullString, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			t.Fatalf("%s: %v", query, err)
		}
		row := make([]string, len(cols))
		for i, v := range vals {
			row[i] = v.String
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("%s: %v", query, err)
	}
	return out
}

// oceanBaseSync reads shop into the metadata Bytebase's sync produces from
// OceanBase, the parts the rules read (walk-through spec §2): base tables
// with their columns, indexes and key parts, foreign keys, and triggers;
// views; external tables; materialized views from oceanbase.DBA_MVIEWS,
// none when unread, as for a user who may not read it; routines; and events
// by name, each once. OceanBase reports version 5.7.25, for which the sync
// reads no CHECK constraint, index expression, or partition.
//
// bytebase: backend/plugin/db/mysql/sync.go — SyncDBSchema, since
// bytebase/bytebase#21685 and bytebase/bytebase#21687
func oceanBaseSync(t testing.TB, db *sql.DB, unread bool) *metadata.DatabaseSchemaMetadata {
	t.Helper()
	s := &metadata.SchemaMetadata{}
	tables := make(map[string]*metadata.TableMetadata)
	externals := make(map[string]*metadata.ExternalTableMetadata)
	for _, r := range read(t, db, "SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME") {
		switch name := r[0]; r[1] {
		case "VIEW":
			s.Views = append(s.Views, &metadata.ViewMetadata{Name: name})
		case "EXTERNAL TABLE":
			externals[name] = &metadata.ExternalTableMetadata{Name: name}
			s.ExternalTables = append(s.ExternalTables, externals[name])
		default:
			tables[name] = &metadata.TableMetadata{Name: name, Engine: r[2]}
			s.Tables = append(s.Tables, tables[name])
		}
	}
	for _, r := range read(t, db, "SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION") {
		c := &metadata.ColumnMetadata{Name: r[1]}
		if x, ok := tables[r[0]]; ok {
			x.Columns = append(x.Columns, c)
		} else if x, ok := externals[r[0]]; ok {
			x.Columns = append(x.Columns, c)
		}
	}
	// A key part is its column and its length, -1 when it is not a prefix.
	indexes := make(map[[2]string]*metadata.IndexMetadata)
	for _, r := range read(t, db, "SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, ''), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX") {
		x, ok := tables[r[0]]
		if !ok {
			continue
		}
		i, ok := indexes[[2]string{r[0], r[1]}]
		if !ok {
			i = &metadata.IndexMetadata{Name: r[1], Primary: r[1] == "PRIMARY", Unique: r[2] == "0"}
			indexes[[2]string{r[0], r[1]}] = i
			x.Indexes = append(x.Indexes, i)
		}
		var length int64
		if _, err := fmt.Sscan(r[4], &length); err != nil {
			t.Fatalf("the length of a key part of %s.%s is %q", r[0], r[1], r[4])
		}
		i.Expressions = append(i.Expressions, r[3])
		i.KeyLength = append(i.KeyLength, length)
	}
	// A foreign key names the table it references without its database.
	foreignKeys := make(map[[2]string]*metadata.ForeignKeyMetadata)
	for _, r := range read(t, db, "SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME") {
		if x, ok := tables[r[0]]; ok {
			fk := &metadata.ForeignKeyMetadata{Name: r[1], ReferencedTable: r[2]}
			foreignKeys[[2]string{r[0], r[1]}] = fk
			x.ForeignKeys = append(x.ForeignKeys, fk)
		}
	}
	for _, r := range read(t, db, "SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION") {
		if fk, ok := foreignKeys[[2]string{r[0], r[1]}]; ok {
			fk.Columns = append(fk.Columns, r[2])
			fk.ReferencedColumns = append(fk.ReferencedColumns, r[3])
		}
	}
	for _, r := range read(t, db, "SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME") {
		if x, ok := tables[r[0]]; ok {
			x.Triggers = append(x.Triggers, &metadata.TriggerMetadata{Name: r[1]})
		}
	}
	for _, r := range read(t, db, "SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'shop' ORDER BY ROUTINE_TYPE, ROUTINE_NAME") {
		if r[0] == "PROCEDURE" {
			s.Procedures = append(s.Procedures, &metadata.ProcedureMetadata{Name: r[1]})
		} else {
			s.Functions = append(s.Functions, &metadata.FunctionMetadata{Name: r[1]})
		}
	}
	for _, r := range read(t, db, "SELECT DISTINCT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = 'shop' ORDER BY EVENT_NAME") {
		s.Events = append(s.Events, &metadata.EventMetadata{Name: r[0]})
	}
	if !unread {
		for _, r := range read(t, db, "SELECT MVIEW_NAME FROM oceanbase.DBA_MVIEWS WHERE OWNER = 'shop' ORDER BY MVIEW_NAME") {
			s.MaterializedViews = append(s.MaterializedViews, &metadata.MaterializedViewMetadata{Name: r[0], SkipDump: true})
		}
	}
	return &metadata.DatabaseSchemaMetadata{Name: "shop", Schemas: []*metadata.SchemaMetadata{s}}
}

// lowerCaseTableNames reads the server's lower_case_table_names.
func lowerCaseTableNames(t *testing.T, db *sql.DB) int {
	t.Helper()
	var v int
	if err := db.QueryRowContext(context.Background(), "SELECT @@lower_case_table_names").Scan(&v); err != nil {
		t.Fatal(err)
	}
	return v
}

// parseChange parses a change the way Review does.
func parseChange(t *testing.T, change string) []statement {
	t.Helper()
	stmts, syntax := parse(change, statementRanges(change))
	if syntax != nil {
		t.Fatalf("syntax: %s", syntax.Message)
	}
	return stmts
}

// TestScanRulesAgainstOceanBase runs each change on OceanBase, statement by
// statement, and derives from information_schema what the constraint rules
// should report: DisallowDropConstraint, a primary key, foreign key, unique
// key, or check constraint a statement removes; RequirePrimaryKey, a table a
// CREATE TABLE made, or one that had a primary key, that ends the change
// without one. Every finding must match; for a complete case, every expected
// finding must be reported.
//
// A statement OceanBase refuses ends the change, and the case names its
// error (stops): the findings on that statement and after it are not judged,
// and the case passes only while OceanBase stops the change so.
func TestScanRulesAgainstOceanBase(t *testing.T) {
	tests := []struct {
		name     string
		change   string
		complete bool
		stops    uint16
	}{
		// OceanBase takes one change of the primary key to a statement (1235),
		// and DROP CONSTRAINT for a foreign key or a check only (3940).
		{name: "constraint drops the SQL names", change: "ALTER TABLE t DROP PRIMARY KEY, DROP FOREIGN KEY fk_p, DROP CHECK ck_n;\nALTER TABLE t DROP CONSTRAINT uk_code;", complete: true, stops: 1235},
		{name: "constraint drops the SQL names, one to a statement", change: "ALTER TABLE t DROP PRIMARY KEY;\nALTER TABLE t DROP FOREIGN KEY fk_p;\nALTER TABLE t DROP CHECK ck_n;", complete: true},
		{name: "DROP CONSTRAINT of a foreign key and a check", change: "ALTER TABLE t DROP CONSTRAINT fk_p;\nALTER TABLE t DROP CONSTRAINT ck_n;", complete: true},
		{name: "DROP CONSTRAINT of a unique key", change: "ALTER TABLE t DROP CONSTRAINT uk_code;", complete: true, stops: 3940},
		{name: "a unique index dropped by name", change: "ALTER TABLE t DROP INDEX uk_code;\nDROP INDEX idx_n ON t;", complete: true},
		{name: "a unique key the change makes, then drops", change: "ALTER TABLE nokey ADD UNIQUE KEY uk_a (a);\nALTER TABLE nokey DROP INDEX uk_a;", complete: true},
		// DROP INDEX does not reach the primary key on OceanBase (1091).
		{name: "DROP INDEX of the primary key", change: "DROP INDEX `PRIMARY` ON plain;", complete: true, stops: 1091},
		{name: "a table created without a key", change: "CREATE TABLE a (id INT);\nCREATE TABLE b (id INT PRIMARY KEY);", complete: true},
		// OceanBase makes no table from a key alone and a query's columns (4028).
		{name: "CREATE TABLE AS SELECT, which copies no key", change: "CREATE TABLE c AS SELECT id, v FROM plain;\nCREATE TABLE d (PRIMARY KEY (id)) AS SELECT id FROM plain;", complete: true, stops: 4028},
		{name: "a key added later", change: "CREATE TABLE a (id INT);\nALTER TABLE a ADD PRIMARY KEY (id);", complete: true},
		{name: "a primary key replaced", change: "ALTER TABLE plain DROP PRIMARY KEY, ADD PRIMARY KEY (id, v);", complete: true},
		{name: "a primary key dropped, then added back", change: "ALTER TABLE plain DROP PRIMARY KEY;\nALTER TABLE plain ADD PRIMARY KEY (v);", complete: true},
		{name: "a table renamed after losing its key", change: "ALTER TABLE plain DROP PRIMARY KEY;\nRENAME TABLE plain TO plain2;"},
		{name: "CREATE TABLE IF NOT EXISTS of a table with a key", change: "CREATE TABLE IF NOT EXISTS plain (id INT);", complete: true},
		{name: "CREATE TABLE IF NOT EXISTS, not judged", change: "CREATE TABLE IF NOT EXISTS a (id INT);"},
		{name: "a table whose database the change drops", change: "CREATE DATABASE d;\nCREATE TABLE d.e (id INT);\nDROP DATABASE d;", complete: true},
		{name: "an index of a table CREATE TABLE IF NOT EXISTS leaves", change: "CREATE TABLE IF NOT EXISTS t (id INT);\nDROP INDEX uk_code ON t;", complete: true},
		{name: "an index of a table the change recreates", change: "DROP TABLE t;\nCREATE TABLE t (code INT, KEY uk_code (code));\nALTER TABLE t DROP INDEX uk_code;"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := openShop(t)
			schema := oceanBaseSync(t, db, false)
			result, err := ReviewOceanBase(context.Background(), tt.change, review.Options{
				Rules: []review.Rule{review.DisallowDropConstraint, review.RequirePrimaryKey},
			}, []review.Target{{Schema: schema, LowerCaseTableNames: lowerCaseTableNames(t, db)}})
			if err != nil {
				t.Fatalf("Review: %v", err)
			}
			want, ran := scanTruth(t, db, tt.change, tt.stops)
			got := make(map[string]bool)
			for _, f := range result.Findings {
				if err := review.CheckCodes(f); err != nil {
					t.Error(err)
				}
				if f.Rule == review.Syntax {
					t.Fatalf("syntax: %s", f.Message)
				}
				if f.Statement < ran {
					got[fmt.Sprintf("%s %d", f.Rule, f.Statement)] = true
				}
			}
			for k := range got {
				if !want[k] {
					t.Errorf("reported %s, which the server contradicts", k)
				}
			}
			if tt.complete {
				for k := range want {
					if !got[k] {
						t.Errorf("missed %s", k)
					}
				}
			}
		})
	}
}

// constraintSet reads every primary key, unique key, foreign key, and
// check constraint of the shop database, as table/kind/name/definition: a
// key's columns, a check's expression. A primary key is always named
// PRIMARY, so the definition tells a replaced one apart.
func constraintSet(t *testing.T, db *sql.DB) map[string]bool {
	t.Helper()
	out := make(map[string]bool)
	for _, r := range read(t, db, `SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY 1, 2, 3`) {
		out[strings.ToLower(strings.Join(r, "/"))] = true
	}
	return out
}

// primaryKeyed reads the tables of the shop database that have a primary
// key, and all its base tables.
func primaryKeyed(t *testing.T, db *sql.DB) (keyed, all map[string]bool) {
	t.Helper()
	keyed, all = make(map[string]bool), make(map[string]bool)
	for _, r := range read(t, db, `SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY 1`) {
		name := strings.ToLower(r[0])
		all[name] = true
		if r[1] == "1" {
			keyed[name] = true
		}
	}
	return keyed, all
}

// scanTruth runs the change statement by statement, as far as OceanBase
// takes it, and returns the findings the server's catalog calls for, as
// "RULE statement", and how many statements ran. A statement OceanBase
// refuses must fail with stops, and a change it runs whole must have none.
func scanTruth(t *testing.T, db *sql.DB, change string, stops uint16) (want map[string]bool, ran int) {
	t.Helper()
	want = make(map[string]bool)
	ranges := statementRanges(change)
	before := constraintSet(t, db)
	keyedBefore, _ := primaryKeyed(t, db)
	keyed := keyedBefore
	// lostAt is the statement that removed a table's key, created the
	// statement that made a table.
	lostAt, createdAt := make(map[string]int), make(map[string]int)
	stopped := false
	for _, s := range parseChange(t, change) {
		r := ranges[s.index]
		if _, err := db.ExecContext(context.Background(), change[r.Start:r.End]); err != nil {
			if stops == 0 || refusal(err) != stops {
				t.Fatalf("statement %d: %v", s.index, err)
			}
			t.Logf("OceanBase refuses statement %d (%v), which ends the change", s.index, err)
			stopped = true
			break
		}
		ran++
		after := constraintSet(t, db)
		for k := range before {
			if !after[k] {
				want[fmt.Sprintf("%s %d", review.DisallowDropConstraint, s.index)] = true
			}
		}
		keyedNow, _ := primaryKeyed(t, db)
		for name := range keyed {
			if !keyedNow[name] {
				lostAt[name] = s.index
			}
		}
		if v, ok := s.node.(*ast.CreateTableStmt); ok && !v.Temporary && v.Like == nil {
			createdAt[strings.ToLower(v.Table.Name)] = s.index
		}
		before, keyed = after, keyedNow
	}
	if stops != 0 && !stopped {
		t.Fatalf("OceanBase runs the whole change, which it stopped with %d: take stops off the case", stops)
	}
	keyedEnd, all := primaryKeyed(t, db)
	for name, at := range createdAt {
		if all[name] && !keyedEnd[name] {
			want[fmt.Sprintf("%s %d", review.RequirePrimaryKey, at)] = true
		}
	}
	for name, at := range lostAt {
		if _, made := createdAt[name]; !made && keyedBefore[name] && all[name] && !keyedEnd[name] {
			want[fmt.Sprintf("%s %d", review.RequirePrimaryKey, at)] = true
		}
	}
	return want, ran
}

// TestColumnDropAgainstOceanBase runs ALTER TABLE ... DROP COLUMN on a
// column of a table x's primary key on OceanBase, alone and beside what
// MySQL's rule follows around such a drop, and pins that x keeps its key.
// MySQL drops the column and, with its last column, the key, which
// REQUIRE_PRIMARY_KEY reports there. OceanBase refuses to drop a column of a
// primary key (1235), so ReviewOceanBase judges no DROP COLUMN, and must
// report nothing here. The day an OceanBase runs one of these drops, its
// case fails, and the rule has a drop to judge.
//
// Each case names the error that stops its change (stops), or the one that
// refuses its setup (refused) where the setup is of what only MySQL has.
//
// oceanbase: https://github.com/oceanbase/oceanbase-doc/blob/V4.3.5/en-US/700.reference/500.sql-reference/100.sql-syntax/200.common-tenant-of-mysql-mode/600.sql-statement-of-mysql-mode/1600.alter-table-of-mysql-mode.md — DROP [COLUMN]
func TestColumnDropAgainstOceanBase(t *testing.T) {
	tests := []struct {
		name    string
		setup   string
		change  string
		stops   uint16
		refused uint16
	}{
		// A key of one column.
		{name: "a key of one column", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); INSERT INTO x VALUES (1, 1), (2, 2)", change: "ALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a key of one AUTO_INCREMENT column", setup: "CREATE TABLE x (id INT AUTO_INCREMENT PRIMARY KEY, v INT); INSERT INTO x (v) VALUES (1), (2)", change: "ALTER TABLE x DROP id", stops: 1235},
		{name: "a key of a column prefix", setup: "CREATE TABLE x (s VARCHAR(50), v INT, PRIMARY KEY (s(10)))", change: "ALTER TABLE x DROP COLUMN s", refused: 1064},
		{name: "a key of one column, with a unique key on a NOT NULL column left", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT NOT NULL, UNIQUE KEY uv (v))", change: "ALTER TABLE x DROP COLUMN id", stops: 1235},
		// A key of two columns.
		{name: "one column of a key of two", setup: "CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b)); INSERT INTO x VALUES (1, 1, 1), (1, 2, 2)", change: "ALTER TABLE x DROP COLUMN a", stops: 1235},
		{name: "one column of a key of two, the other repeating", setup: "CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b)); INSERT INTO x VALUES (1, 1, 1), (2, 1, 2)", change: "ALTER TABLE x DROP COLUMN a", stops: 1235},
		{name: "both columns of a key of two, in two statements", setup: "CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b))", change: "ALTER TABLE x DROP COLUMN a;\nALTER TABLE x DROP COLUMN b", stops: 1235},
		{name: "both columns of a key of two, in two statements, a unique key of the other keeping its values apart", setup: "CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b), UNIQUE KEY ub (b)); INSERT INTO x VALUES (1, 1, 1), (1, 2, 2)", change: "ALTER TABLE x DROP COLUMN a;\nALTER TABLE x DROP COLUMN b", stops: 1235},
		{name: "both columns of a key of two, in two statements, a unique key of a prefix of the other keeping its values apart", setup: "CREATE TABLE x (a INT, b VARCHAR(20), v INT, PRIMARY KEY (a, b), UNIQUE KEY ub (b(2))); INSERT INTO x VALUES (1, 'aa1', 1), (1, 'bb2', 2)", change: "ALTER TABLE x DROP COLUMN a;\nALTER TABLE x DROP COLUMN b", stops: 1235},
		{name: "one column of a key of two whose other is a prefix, beside a unique key of the whole other", setup: "CREATE TABLE x (a INT, b VARCHAR(20), v INT, PRIMARY KEY (a, b(2)), UNIQUE KEY ub (b)); INSERT INTO x VALUES (1, 'xx1', 1), (2, 'xx2', 2)", change: "ALTER TABLE x DROP COLUMN a", refused: 1064},
		{name: "both columns of a key of two", setup: "CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b)); INSERT INTO x VALUES (1, 1, 1), (2, 1, 2)", change: "ALTER TABLE x DROP COLUMN b, DROP COLUMN a", stops: 1235},
		// Beside other columns, and what reads them.
		{name: "another column dropped before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, w INT)", change: "ALTER TABLE x DROP COLUMN v;\nALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a generated column reading it, dropped before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1))", change: "ALTER TABLE x DROP COLUMN g;\nALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a generated column reading it, dropped too", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1) STORED)", change: "ALTER TABLE x DROP COLUMN id, DROP COLUMN g", stops: 3108},
		{name: "a check of it alone", setup: "CREATE TABLE x (id INT PRIMARY KEY CHECK (id > 0), v INT)", change: "ALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a check of both columns of the key, both dropped", setup: "CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b), CONSTRAINT ck CHECK (a < b))", change: "ALTER TABLE x DROP COLUMN a, DROP COLUMN b", stops: 3959},
		{name: "another index holds it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, KEY iv (v, id))", change: "ALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a view or a trigger reads it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE VIEW xv AS SELECT id FROM x; CREATE TRIGGER xt BEFORE INSERT ON x FOR EACH ROW SET NEW.v = NEW.id", change: "ALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a foreign key of another table references another column", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, UNIQUE KEY uv (v)); CREATE TABLE ch (x_v INT, CONSTRAINT fk_ch FOREIGN KEY (x_v) REFERENCES x (v))", change: "ALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a foreign key of the table holds another column", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, CONSTRAINT fk_x FOREIGN KEY (v) REFERENCES p (id))", change: "ALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "statements on another table, on it that give it no key, or on none after", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)", change: "ALTER TABLE x DROP COLUMN id;\nINSERT INTO y VALUES (1);\nALTER TABLE y ADD COLUMN b INT;\nCREATE INDEX iv ON x (v);\nALTER TABLE x ADD COLUMN w INT;\nINSERT INTO x (v) VALUES (1);\nCREATE VIEW xv AS SELECT v FROM x", stops: 1235},
		{name: "a generated column naming it only in a literal", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, g VARCHAR(20) AS (CONCAT('id', v)))", change: "ALTER TABLE x DROP COLUMN id", stops: 1235},
		// What depends on the column refuses the drop before the key does.
		{name: "a foreign key of another table references it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE ch (x_id INT, CONSTRAINT fk_ch FOREIGN KEY (x_id) REFERENCES x (id))", change: "ALTER TABLE x DROP COLUMN id", stops: 5320},
		{name: "a foreign key of another table references it, with foreign_key_checks off", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE ch (x_id INT, CONSTRAINT fk_ch FOREIGN KEY (x_id) REFERENCES x (id))", change: "SET foreign_key_checks = 0;\nALTER TABLE x DROP COLUMN id", stops: 5320},
		{name: "a foreign key of another database references it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE DATABASE other; CREATE TABLE other.ch (x_id INT, CONSTRAINT fk_och FOREIGN KEY (x_id) REFERENCES shop.x (id))", change: "ALTER TABLE x DROP COLUMN id", stops: 5320},
		{name: "a foreign key of the table holds it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, CONSTRAINT fk_x FOREIGN KEY (id) REFERENCES p (id))", change: "ALTER TABLE x DROP COLUMN id", stops: 5320},
		{name: "a generated column reads it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1))", change: "ALTER TABLE x DROP COLUMN id", stops: 3108},
		{name: "a check reads it and another column", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, CONSTRAINT ck CHECK (id < v))", change: "ALTER TABLE x DROP COLUMN id", stops: 3959},
		{name: "a functional index reads it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, KEY fx ((id + v)))", change: "ALTER TABLE x DROP COLUMN id", stops: 3837},
		{name: "the table is partitioned by it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT) PARTITION BY HASH (id) PARTITIONS 2", change: "ALTER TABLE x DROP COLUMN id", stops: 3855},
		{name: "the table is partitioned by KEY (), its primary key", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT) PARTITION BY KEY () PARTITIONS 2", change: "ALTER TABLE x DROP COLUMN id", stops: 3855},
		{name: "it is the only column", setup: "CREATE TABLE x (id INT PRIMARY KEY)", change: "ALTER TABLE x DROP COLUMN id", stops: 1090},
		{name: "with a column the table lacks", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id, DROP COLUMN missing", stops: 1091},
		{name: "twice", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id, DROP COLUMN id", stops: 1091},
		// With another change to the table.
		{name: "it is added back under its name", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id, ADD COLUMN id INT", stops: 1060},
		{name: "a key is added back", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id, ADD PRIMARY KEY (v)", stops: 1068},
		{name: "a later statement adds a key back", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id;\nALTER TABLE x ADD PRIMARY KEY (v)", stops: 1235},
		{name: "a later statement that adds a key back, but fails", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id;\nALTER TABLE x DROP COLUMN missing, ADD PRIMARY KEY (v)", stops: 1235},
		{name: "a check of another column added before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); INSERT INTO x VALUES (1, 1)", change: "ALTER TABLE x ADD CHECK (v > 0);\nALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a foreign key to it added and dropped before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE ch (k INT PRIMARY KEY, x_id INT)", change: "ALTER TABLE ch ADD CONSTRAINT f FOREIGN KEY (x_id) REFERENCES x (id);\nALTER TABLE ch DROP FOREIGN KEY f;\nALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a table CREATE TABLE IF NOT EXISTS leaves, whose foreign key would reference it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE c (k INT PRIMARY KEY)", change: "CREATE TABLE IF NOT EXISTS c (k INT PRIMARY KEY, FOREIGN KEY (k) REFERENCES x (id));\nALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a temporary table of its name dropped twice", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "DROP TEMPORARY TABLE IF EXISTS x;\nDROP TEMPORARY TABLE x;\nALTER TABLE x DROP COLUMN id", stops: 1051},
		{name: "with a check the table lacks altered", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x ALTER CHECK missing ENFORCED, DROP COLUMN id", stops: 3821},
		{name: "a table made before whose column's REFERENCES, which MySQL ignores, names it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "CREATE TABLE c (k INT REFERENCES x (id), PRIMARY KEY (k));\nALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "with a partition dropped", setup: "CREATE TABLE x (id INT, v INT, PRIMARY KEY (id)) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN MAXVALUE)", change: "ALTER TABLE x DROP COLUMN id, DROP PARTITION p0", stops: 1064},
		// OceanBase has no sql_require_primary_key (1193).
		{name: "sql_require_primary_key on after", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id;\nSET SESSION sql_require_primary_key = ON", stops: 1235},
		{name: "sql_require_primary_key off before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "SET SESSION sql_require_primary_key = OFF;\nALTER TABLE x DROP COLUMN id", stops: 1193},
		{name: "sql_require_primary_key on before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "SET SESSION sql_require_primary_key = ON;\nALTER TABLE x DROP COLUMN id", stops: 1193},
		{name: "a user variable named sql_require_primary_key set before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "SET @sql_require_primary_key = 1;\nALTER TABLE x DROP COLUMN id", stops: 1235},
		// Invisible columns and defaults.
		{name: "only an invisible column left", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT INVISIBLE)", change: "ALTER TABLE x DROP COLUMN id", stops: 4028},
		{name: "an invisible column and a visible one left", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT INVISIBLE, w INT)", change: "ALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "an expression default reads it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT DEFAULT (id + 1))", change: "ALTER TABLE x DROP COLUMN id", refused: 1235},
		{name: "a literal default names it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v VARCHAR(10) DEFAULT 'id')", change: "ALTER TABLE x DROP COLUMN id", stops: 1235},
		// A table of another database renamed without a database lands in the
		// current one.
		{name: "a table of another database renamed to its name before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE DATABASE other; CREATE TABLE other.t (a INT PRIMARY KEY)", change: "ALTER TABLE other.t RENAME TO x;\nALTER TABLE x DROP COLUMN id", stops: 1050},
		// A statement before the drop that fails stops the change there.
		{name: "a statement on another table before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)", change: "ALTER TABLE y ADD COLUMN b INT;\nALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a table the database lacks dropped before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "DROP TABLE missing;\nALTER TABLE x DROP COLUMN id", stops: 1051},
		{name: "a database the server lacks dropped before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "DROP DATABASE missing_db;\nALTER TABLE x DROP COLUMN id", stops: 1008},
		{name: "sql_generate_invisible_primary_key on before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "SET SESSION sql_generate_invisible_primary_key = ON;\nALTER TABLE x DROP COLUMN id", stops: 1193},
		{name: "sql_generate_invisible_primary_key on, then off, before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "SET SESSION sql_generate_invisible_primary_key = ON;\nSET SESSION sql_generate_invisible_primary_key = OFF;\nALTER TABLE x DROP COLUMN id", stops: 1193},
		// TRUNCATE, LOCK TABLES, and DML around the drop.
		{name: "a truncate of a table a foreign key references before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE tp (a INT PRIMARY KEY); CREATE TABLE tc (a INT, FOREIGN KEY (a) REFERENCES tp (a))", change: "TRUNCATE TABLE tp;\nALTER TABLE x DROP COLUMN id", stops: 1701},
		{name: "a truncate after", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id;\nTRUNCATE TABLE x", stops: 1235},
		{name: "tables locked for reading before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "LOCK TABLES x READ;\nALTER TABLE x DROP COLUMN id;\nUNLOCK TABLES", stops: 1235},
		{name: "another table locked before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)", change: "LOCK TABLES y WRITE;\nALTER TABLE x DROP COLUMN id;\nUNLOCK TABLES", stops: 1235},
		{name: "tables locked after", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id;\nLOCK TABLES x READ;\nUNLOCK TABLES", stops: 1235},
		{name: "a row inserted before into a table the database lacks", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "INSERT INTO missing VALUES (1);\nALTER TABLE x DROP COLUMN id", stops: 1146},
		{name: "rows of the table updated before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); INSERT INTO x VALUES (1, 1)", change: "UPDATE x SET v = id + 1 WHERE id IN (SELECT id FROM (SELECT id FROM x) d);\nALTER TABLE x DROP COLUMN id", stops: 1235},
		// Statements before the drop that fail, and ones that cannot.
		{name: "a column the table lacks inserted into before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "INSERT INTO x (missing) VALUES (1);\nALTER TABLE x DROP COLUMN id", stops: 1054},
		{name: "the database's tables shown before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "SHOW TABLES FROM shop;\nALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a rollback to a savepoint the transaction lacks before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "BEGIN;\nROLLBACK TO SAVEPOINT nosuch;\nALTER TABLE x DROP COLUMN id", stops: 1305},
		{name: "statements that cannot fail before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "BEGIN;\nSAVEPOINT s;\nROLLBACK;\nCOMMIT;\nUNLOCK TABLES;\nANALYZE TABLE missing;\nANALYZE TABLE x UPDATE HISTOGRAM ON missing;\nSET foreign_key_checks = OFF, @@session.unique_checks = 0, @x = 'a', @h = X'41', @b = b'101', @n = -1.5e3, @h2 = X'aB', @h3 = 0xaB, @b2 = B'01', foreign_key_checks = DEFAULT, unique_checks = DEFAULT;\nALTER TABLE x DROP COLUMN id", stops: 1146},
		{name: "sql_require_primary_key set with a variable the server reads only before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "SET SESSION sql_require_primary_key = OFF, version = 'x';\nALTER TABLE x DROP COLUMN id", stops: 1193},
		// Unique keys that hold the column.
		{name: "a unique key of it and another column", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, UNIQUE KEY u (id, v)); INSERT INTO x VALUES (1, 1), (2, 2)", change: "ALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a unique key of it and another column, the other repeating", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, UNIQUE KEY u (id, v)); INSERT INTO x VALUES (1, 1), (2, 1)", change: "ALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a unique key of it and another column, both dropped", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, w INT, UNIQUE KEY u (id, v)); INSERT INTO x VALUES (1, 1, 1), (2, 1, 2)", change: "ALTER TABLE x DROP COLUMN id, DROP COLUMN v", stops: 1235},
		{name: "a unique key of it alone", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, UNIQUE KEY u (id)); INSERT INTO x VALUES (1, 1), (2, 1)", change: "ALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a column out of the key dropped before from a table partitioned by KEY ()", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE kp (id INT PRIMARY KEY, v INT) PARTITION BY KEY () PARTITIONS 2", change: "ALTER TABLE kp DROP COLUMN v;\nALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a column out of any key dropped before from a table without a primary key partitioned by KEY ()", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE kp (u INT NOT NULL, v INT, w INT, UNIQUE KEY uu (u)) PARTITION BY KEY () PARTITIONS 2", change: "ALTER TABLE kp DROP COLUMN v;\nALTER TABLE x DROP COLUMN id", refused: 1488},
		{name: "the unique key's column dropped before from a table without a primary key partitioned by KEY ()", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE kp (u INT NOT NULL, v INT, w INT, UNIQUE KEY uu (u)) PARTITION BY KEY () PARTITIONS 2", change: "ALTER TABLE kp DROP COLUMN u;\nALTER TABLE x DROP COLUMN id", refused: 1488},
		{name: "a user variable set to a double out of range before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "SET @x = 1e400;\nALTER TABLE x DROP COLUMN id", stops: 1367},
		{name: "a user variable set to a malformed hexadecimal value before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "SET @x = X'4';\nALTER TABLE x DROP COLUMN id", stops: 1064},
		{name: "a scope for a system variable before a user variable", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "SET SESSION foreign_key_checks = 1, @x = 1;\nALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a unique key of it and another column, beside a unique key of the other", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT NOT NULL, w INT, UNIQUE KEY u (v), UNIQUE KEY uv (id, v)); INSERT INTO x VALUES (1, 1, 1), (2, 2, 2)", change: "ALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a unique key of it and a prefix of another column, beside a unique key of the other", setup: "CREATE TABLE x (id INT PRIMARY KEY, v VARCHAR(20) NOT NULL, UNIQUE KEY u (v), UNIQUE KEY uv (id, v(2))); INSERT INTO x VALUES (1, 'xx1'), (2, 'xx2')", change: "ALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a table made of its name after, which fails", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id;\nCREATE TABLE x (id INT PRIMARY KEY)", stops: 1235},
		{name: "a table renamed onto it after, which fails", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)", change: "ALTER TABLE x DROP COLUMN id;\nRENAME TABLE y TO x", stops: 1235},
		{name: "a table altered to be renamed onto it after, which fails", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)", change: "ALTER TABLE x DROP COLUMN id;\nALTER TABLE y RENAME TO x", stops: 1235},
		{name: "a user variable set to a hexadecimal value of an uppercase prefix before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "SET @x = 0X41;\nALTER TABLE x DROP COLUMN id", stops: 1064},
		{name: "foreign_key_checks set with := before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "SET foreign_key_checks := OFF;\nALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a column of the key held by a prefix dropped before from a table partitioned by KEY (), which leaves the prefix out", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE kp (a VARCHAR(20), b INT, v INT, PRIMARY KEY (a(10), b), UNIQUE KEY ub (b)) PARTITION BY KEY () PARTITIONS 2", change: "ALTER TABLE kp DROP COLUMN a;\nALTER TABLE x DROP COLUMN id", refused: 1064},
		{name: "tables flushed after", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id;\nFLUSH TABLES", stops: 1235},
		{name: "tables flushed with a read lock before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "FLUSH TABLES WITH READ LOCK;\nALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a column of only some unique keys dropped before from a table without a primary key partitioned by KEY ()", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE kp (u INT NOT NULL, v INT NOT NULL, w INT, UNIQUE KEY uu (u), UNIQUE KEY uv (u, v)) PARTITION BY KEY () PARTITIONS 2", change: "ALTER TABLE kp DROP COLUMN v;\nALTER TABLE x DROP COLUMN id", refused: 1488},
		{name: "a rename onto a table the database holds after, which fails", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)", change: "ALTER TABLE x DROP COLUMN id;\nRENAME TABLE x TO y", stops: 1235},
		{name: "a key added back and a rename onto a table the database holds after, which fail", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)", change: "ALTER TABLE x DROP COLUMN id;\nALTER TABLE x ADD PRIMARY KEY (v), RENAME TO y", stops: 1235},
		{name: "a rename after with a pair onto a name an earlier pair moved a table back onto, which fails", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY); CREATE TABLE z (b INT PRIMARY KEY)", change: "ALTER TABLE x DROP COLUMN id;\nRENAME TABLE y TO spare, x TO y, z TO y", stops: 1235},
		{name: "a rename after with a pair onto a name an earlier pair made, which fails", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY); CREATE TABLE z (b INT PRIMARY KEY)", change: "ALTER TABLE x DROP COLUMN id;\nRENAME TABLE x TO spare, y TO spare", stops: 1235},
		{name: "a rename after with a pair onto a name an earlier pair made and a later one moved away", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY); CREATE TABLE z (b INT PRIMARY KEY)", change: "ALTER TABLE x DROP COLUMN id;\nRENAME TABLE x TO spare, spare TO s2, y TO spare", stops: 1235},
		{name: "a rename after onto a table an earlier pair moves away", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)", change: "ALTER TABLE x DROP COLUMN id;\nRENAME TABLE y TO y2, x TO y", stops: 1235},
		{name: "a temporary table of its name made after, given a key, and dropped", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id;\nCREATE TEMPORARY TABLE x (a INT);\nALTER TABLE x ADD PRIMARY KEY (a);\nDROP TABLE x", stops: 1235},
		{name: "a temporary table of its name dropped after, and then it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id;\nCREATE TEMPORARY TABLE x (a INT);\nDROP TABLE x;\nDROP TABLE x", stops: 1235},
		{name: "a drop after of it, which a foreign key of the database refuses", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, UNIQUE KEY uv (v)); CREATE TABLE c (xv INT, FOREIGN KEY (xv) REFERENCES x (v))", change: "ALTER TABLE x DROP COLUMN id;\nDROP TABLE x", stops: 1235},
		{name: "a drop after of it, which a foreign key naming a table of its name in another database does not refuse", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE DATABASE other; CREATE TABLE other.x (id INT PRIMARY KEY); CREATE TABLE c (xid INT, FOREIGN KEY (xid) REFERENCES other.x (id))", change: "ALTER TABLE x DROP COLUMN id;\nDROP TABLE x", stops: 1235},
		{name: "a rename of its name onto a table the database holds after, which MySQL runs on it though a temporary table of its name exists, and fails", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)", change: "ALTER TABLE x DROP COLUMN id;\nCREATE TEMPORARY TABLE x (a INT);\nRENAME TABLE x TO y;\nDROP TABLE x", stops: 1235},
		{name: "a temporary table renamed onto its name after, and dropped", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id;\nCREATE TEMPORARY TABLE tmp (a INT);\nALTER TABLE tmp RENAME TO x;\nDROP TABLE x", stops: 1235},
		{name: "a RENAME TABLE of a temporary table onto its name after, which MySQL does not rename", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id;\nCREATE TEMPORARY TABLE tmp (a INT);\nRENAME TABLE tmp TO x;\nDROP TABLE x", stops: 1235},
		{name: "a statement prepared and deallocated after, which runs nothing", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id;\nPREPARE s FROM 'ALTER TABLE x ADD PRIMARY KEY (v)';\nDEALLOCATE PREPARE s", stops: 1235},
		{name: "a statement prepared before of a table the database lacks", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "PREPARE s FROM 'SELECT * FROM missing';\nALTER TABLE x DROP COLUMN id", stops: 1146},
		{name: "a read-only transaction started before, which the drop commits", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "START TRANSACTION READ ONLY;\nALTER TABLE x DROP COLUMN id", stops: 1792},
		// ALGORITHM and LOCK.
		{name: "with ALGORITHM and LOCK left at DEFAULT", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id, ALGORITHM=DEFAULT, LOCK DEFAULT", stops: 1235},
		{name: "with ALGORITHM set to INSTANT, then DEFAULT", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id, ALGORITHM=INSTANT, ALGORITHM=DEFAULT", stops: 1235},
		{name: "with ALGORITHM set to INSTANT before LOCK left at DEFAULT", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id, ALGORITHM=INSTANT, LOCK=DEFAULT", stops: 1235},
		{name: "with ALGORITHM set to DEFAULT, then INSTANT", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id, ALGORITHM=DEFAULT, ALGORITHM=INSTANT", stops: 1235},
		{name: "with ALGORITHM=INPLACE", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id, ALGORITHM=INPLACE", stops: 1235},
		{name: "with LOCK=NONE", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id, LOCK=NONE", stops: 1235},
		{name: "with a quoted DEFAULT", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id, ALGORITHM=`DEFAULT`", stops: 1064},
		// USE before and after the drop.
		{name: "a USE of its database before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "USE shop;\nALTER TABLE x DROP COLUMN id", stops: 1235},
		{name: "a USE of a database the server lacks before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "USE nosuch_db;\nALTER TABLE x DROP COLUMN id", stops: 1049},
		{name: "a key added after a USE of another database, to a table of its name there", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE DATABASE other; CREATE TABLE other.x (id INT, v INT)", change: "ALTER TABLE x DROP COLUMN id;\nUSE other;\nALTER TABLE x ADD PRIMARY KEY (v)", stops: 1235},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := openShop(t)
			ctx := context.Background()
			if !setUp(t, db, splitWalk(tt.setup), tt.refused) {
				return
			}
			// On OceanBase the rule reads the statements alone, so the target
			// carries no synced schema.
			result, err := ReviewOceanBase(ctx, tt.change, review.Options{Rules: []review.Rule{review.RequirePrimaryKey}}, []review.Target{{LowerCaseTableNames: 1}})
			if err != nil {
				t.Fatalf("Review: %v", err)
			}
			for _, f := range result.Findings {
				if err := review.CheckCodes(f); err != nil {
					t.Error(err)
				}
				if f.Rule == review.Syntax {
					t.Fatalf("syntax: %s", f.Message)
				}
				t.Errorf("reported %q on statement %d, but OceanBase leaves x its key", f.Message, f.Statement)
			}
			before := primaryKeyOf(t, db, "x")
			if before == "" {
				t.Fatalf("x has no primary key before the change")
			}
			var stopped uint16
			for _, stmt := range splitWalk(tt.change) {
				if _, err := db.ExecContext(ctx, stmt); err != nil {
					if stopped = refusal(err); stopped == 0 {
						t.Fatalf("%s: %v", stmt, err)
					}
					t.Logf("OceanBase stops the change at %s: %v", stmt, err)
					break
				}
			}
			if stopped != tt.stops {
				t.Errorf("OceanBase stops the change with %d, want %d (0 for a change it runs whole)", stopped, tt.stops)
			}
			if after := primaryKeyOf(t, db, "x"); after != before {
				t.Errorf("the key of x is %q after the change, and was %q: OceanBase changed the key, so REQUIRE_PRIMARY_KEY has a DROP COLUMN to judge", after, before)
			}
		})
	}
}

// primaryKeyOf reads the columns of a shop table's primary key, comma
// separated, "" when it has none.
func primaryKeyOf(t *testing.T, db *sql.DB, table string) string {
	t.Helper()
	var key sql.NullString
	if err := db.QueryRowContext(context.Background(), `SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'`, table).Scan(&key); err != nil {
		t.Fatal(err)
	}
	return key.String
}
