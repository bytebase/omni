package catalog

import (
	"strings"
	"testing"
)

// loadSDLForDB is a test helper: it loads SDL and returns the named database
// from the resulting catalog, failing the test on error or a missing database.
func loadSDLForDB(t *testing.T, sql, dbName string) *Database {
	t.Helper()
	c, err := LoadSDL(sql)
	if err != nil {
		t.Fatalf("LoadSDL error: %v", err)
	}
	if c == nil {
		t.Fatal("LoadSDL returned nil catalog")
	}
	db := c.GetDatabase(dbName)
	if db == nil {
		t.Fatalf("database %q not found in catalog", dbName)
	}
	return db
}

func TestSDLValidation(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr string // if non-empty, expect error containing this
	}{
		// ---- Accepted statements ----
		{
			name: "empty string returns empty catalog",
			sql:  "",
		},
		{
			name: "valid CREATE DATABASE",
			sql:  "CREATE DATABASE app;",
		},
		{
			name: "valid CREATE TABLE",
			sql:  "CREATE DATABASE app; USE app; CREATE TABLE t (id int);",
		},
		{
			name: "valid CREATE VIEW",
			sql:  "CREATE DATABASE app; USE app; CREATE TABLE t (id int); CREATE VIEW v AS SELECT id FROM t;",
		},
		{
			name: "valid CREATE INDEX",
			sql:  "CREATE DATABASE app; USE app; CREATE TABLE t (id int); CREATE INDEX idx ON t (id);",
		},
		{
			name: "valid CREATE FUNCTION",
			sql:  "CREATE DATABASE app; USE app; CREATE FUNCTION f() RETURNS int DETERMINISTIC RETURN 1;",
		},
		{
			name: "valid CREATE PROCEDURE",
			sql:  "CREATE DATABASE app; USE app; CREATE PROCEDURE p() BEGIN SELECT 1; END;",
		},
		{
			name: "valid CREATE TRIGGER",
			sql:  "CREATE DATABASE app; USE app; CREATE TABLE t (id int); CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW SET NEW.id = 1;",
		},
		{
			name: "valid CREATE EVENT",
			sql:  "CREATE DATABASE app; USE app; CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1;",
		},
		{
			name: "valid SET",
			sql:  "SET foreign_key_checks = 0;",
		},

		// ---- Rejected DML ----
		{
			name:    "INSERT rejected",
			sql:     "CREATE DATABASE app; USE app; CREATE TABLE t (id int); INSERT INTO t VALUES (1);",
			wantErr: "SDL does not allow INSERT statements",
		},
		{
			name:    "UPDATE rejected",
			sql:     "UPDATE t SET x = 1;",
			wantErr: "SDL does not allow UPDATE statements",
		},
		{
			name:    "DELETE rejected",
			sql:     "DELETE FROM t;",
			wantErr: "SDL does not allow DELETE statements",
		},
		{
			name:    "SELECT rejected",
			sql:     "SELECT 1;",
			wantErr: "SDL does not allow SELECT statements",
		},
		// ---- Rejected destructive DDL ----
		{
			name:    "DROP TABLE rejected",
			sql:     "DROP TABLE t;",
			wantErr: "SDL does not allow DROP statements",
		},
		{
			name:    "DROP DATABASE rejected",
			sql:     "DROP DATABASE app;",
			wantErr: "SDL does not allow DROP statements",
		},
		{
			name:    "DROP INDEX rejected",
			sql:     "DROP INDEX idx ON t;",
			wantErr: "SDL does not allow DROP statements",
		},
		{
			name:    "DROP VIEW rejected",
			sql:     "DROP VIEW v;",
			wantErr: "SDL does not allow DROP statements",
		},
		{
			name:    "TRUNCATE rejected",
			sql:     "TRUNCATE TABLE t;",
			wantErr: "SDL does not allow TRUNCATE statements",
		},
		{
			name:    "RENAME TABLE rejected",
			sql:     "RENAME TABLE a TO b;",
			wantErr: "SDL does not allow RENAME statements",
		},
		{
			name:    "ALTER TABLE ADD COLUMN rejected",
			sql:     "ALTER TABLE t ADD COLUMN x int;",
			wantErr: "SDL does not allow ALTER TABLE",
		},
		{
			name:    "ALTER TABLE DROP COLUMN rejected",
			sql:     "ALTER TABLE t DROP COLUMN x;",
			wantErr: "SDL does not allow ALTER TABLE",
		},
		// ---- Parse error ----
		{
			name:    "parse error returns error",
			sql:     "CREATE TABL t (id int);",
			wantErr: "", // non-empty error from parser, asserted specially below
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "parse error returns error" {
				_, err := LoadSDL(tt.sql)
				if err == nil {
					t.Fatal("expected parse error, got nil")
				}
				return
			}

			c, err := LoadSDL(tt.sql)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if c == nil {
				t.Fatal("expected non-nil catalog")
			}
		})
	}
}

// TestSDLLoad_UnterminatedComment_NoPanic guards the declarative entry point
// against malformed input. An unterminated /*! ... */ executable comment (no
// closing */) once panicked the lexer with "slice bounds out of range"; a
// malformed schema submitted to LoadSDL must return a clean error, never crash.
func TestSDLLoad_UnterminatedComment_NoPanic(t *testing.T) {
	cases := []struct {
		name string
		sql  string
	}{
		// The confirmed real-world repro (no closing */ on the version comment).
		{"versioned_partition_comment", "CREATE TABLE `t` (`id` int NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB /*!50100 PARTITION BY HASH (`id`)"},
		{"bare_exec_comment", "CREATE TABLE `t` (`id` int) /*! FOO"},
		{"plain_block_comment", "CREATE TABLE `t` (`id` int) /* nope"},
		{"unterminated_string", "CREATE TABLE `t` (`id` int) COMMENT='oops"},
		{"unterminated_ident", "CREATE TABLE `t"},
		// Trailing unterminated comment after a valid statement — the trailing
		// segment must not be silently dropped by the splitter.
		{"trailing_unterminated_comment", "CREATE TABLE `t` (`id` int); /*!50100"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("LoadSDL panicked on %q: %v", tc.sql, r)
				}
			}()
			c, err := LoadSDL(tc.sql)
			if err == nil {
				t.Fatalf("expected a clean error for malformed input, got nil (catalog=%v)", c)
			}
		})
	}
}

// TestSDLLoad_ValidVersionComment_Partition confirms the fix did not break the
// valid, terminated version-comment form: a /*!50100 PARTITION BY ... */ clause
// still loads and the table carries its partitioning.
func TestSDLLoad_ValidVersionComment_Partition(t *testing.T) {
	sql := "CREATE DATABASE app;\nUSE app;\nCREATE TABLE `t` (`id` int NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB /*!50100 PARTITION BY HASH (`id`) PARTITIONS 4 */;"
	db := loadSDLForDB(t, sql, "app")
	tbl := db.GetTable("t")
	if tbl == nil {
		t.Fatal("expected table t")
	}
	if tbl.Partitioning == nil {
		t.Error("expected partitioning to survive the /*!50100 ... */ version comment, got none")
	}
}

// TestSDLForwardForeignKey verifies a table with a forward FK reference (the
// referenced table is declared *after* the referencing table) loads cleanly.
// LoadSQL would fail here because foreign_key_checks defaults on and the parent
// table does not yet exist.
func TestSDLForwardForeignKey(t *testing.T) {
	sql := `
CREATE DATABASE app;
USE app;
CREATE TABLE child (
  id int PRIMARY KEY,
  parent_id int,
  CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent (id)
);
CREATE TABLE parent (
  id int PRIMARY KEY
);
`
	db := loadSDLForDB(t, sql, "app")
	child := db.GetTable("child")
	if child == nil {
		t.Fatal("expected table child")
	}
	if db.GetTable("parent") == nil {
		t.Fatal("expected table parent")
	}
	// The FK constraint must be present on child.
	var fkFound bool
	for _, con := range child.Constraints {
		if con.Type == ConForeignKey && strings.EqualFold(con.RefTable, "parent") {
			fkFound = true
		}
	}
	if !fkFound {
		t.Fatalf("expected FK constraint on child referencing parent, constraints: %+v", child.Constraints)
	}
}

// TestSDLForwardViewReference verifies a view referencing a table declared after
// it loads and resolves correctly (AnalyzedQuery populated regardless of order).
func TestSDLForwardViewReference(t *testing.T) {
	sql := `
CREATE DATABASE app;
USE app;
CREATE VIEW v AS SELECT id, name FROM t;
CREATE TABLE t (id int, name varchar(50));
`
	db := loadSDLForDB(t, sql, "app")
	v := db.Views["v"]
	if v == nil {
		t.Fatal("expected view v")
	}
	if v.AnalyzedQuery == nil {
		t.Fatal("expected view v to have a populated AnalyzedQuery (referenced table ordered before view)")
	}
}

// TestSDLOrderIndependence is the spine of the loader: the same set of objects
// declared in different statement orders must yield an equivalent catalog.
func TestSDLOrderIndependence(t *testing.T) {
	// A representative multi-object schema: tables with a forward FK, a view
	// referencing a later table, a secondary/composite index, a generated
	// column, a trigger, and a routine.
	objects := []string{
		`CREATE TABLE parent (id int PRIMARY KEY, code varchar(20))`,
		`CREATE TABLE child (
			id int PRIMARY KEY,
			parent_id int,
			amount int,
			doubled int AS (amount * 2) STORED,
			CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent (id)
		)`,
		`CREATE INDEX idx_child_parent ON child (parent_id, amount)`,
		`CREATE VIEW v_join AS SELECT c.id, p.code FROM child c JOIN parent p ON c.parent_id = p.id`,
		`CREATE TRIGGER trg_child BEFORE INSERT ON child FOR EACH ROW SET NEW.amount = NEW.amount`,
		`CREATE FUNCTION fn_double(x int) RETURNS int DETERMINISTIC RETURN x * 2`,
	}

	build := func(order []int) *Database {
		var b strings.Builder
		b.WriteString("CREATE DATABASE app; USE app;\n")
		for _, i := range order {
			b.WriteString(objects[i])
			b.WriteString(";\n")
		}
		return loadSDLForDB(t, b.String(), "app")
	}

	// Canonical (declaration) order.
	canonical := build([]int{0, 1, 2, 3, 4, 5})

	// Several shuffled orders that all place dependents before dependencies.
	shuffles := [][]int{
		{5, 4, 3, 2, 1, 0}, // fully reversed
		{3, 2, 1, 0, 5, 4}, // view & index first
		{4, 3, 5, 2, 0, 1}, // trigger/view/fn before tables
		{2, 5, 3, 4, 0, 1}, // index & view before tables, parent last
	}

	for idx, order := range shuffles {
		got := build(order)
		assertDatabasesEquivalent(t, canonical, got, idx)
	}
}

// assertDatabasesEquivalent checks that two databases loaded from the same set
// of objects (in different orders) contain the same objects with the same
// salient shape. It compares object names per kind plus key per-table details.
func assertDatabasesEquivalent(t *testing.T, want, got *Database, shuffleIdx int) {
	t.Helper()

	assertKeySetsEqual(t, "tables", keys(want.Tables), keys(got.Tables), shuffleIdx)
	assertKeySetsEqual(t, "views", keys(want.Views), keys(got.Views), shuffleIdx)
	assertKeySetsEqual(t, "functions", keys(want.Functions), keys(got.Functions), shuffleIdx)
	assertKeySetsEqual(t, "procedures", keys(want.Procedures), keys(got.Procedures), shuffleIdx)
	assertKeySetsEqual(t, "triggers", keys(want.Triggers), keys(got.Triggers), shuffleIdx)
	assertKeySetsEqual(t, "events", keys(want.Events), keys(got.Events), shuffleIdx)

	for name, wantTbl := range want.Tables {
		gotTbl := got.Tables[name]
		if gotTbl == nil {
			t.Fatalf("shuffle %d: table %q missing", shuffleIdx, name)
		}
		// Columns: same names in the same order.
		if len(wantTbl.Columns) != len(gotTbl.Columns) {
			t.Fatalf("shuffle %d: table %q column count: want %d got %d", shuffleIdx, name, len(wantTbl.Columns), len(gotTbl.Columns))
		}
		for i := range wantTbl.Columns {
			if !strings.EqualFold(wantTbl.Columns[i].Name, gotTbl.Columns[i].Name) {
				t.Fatalf("shuffle %d: table %q column %d: want %q got %q", shuffleIdx, name, i, wantTbl.Columns[i].Name, gotTbl.Columns[i].Name)
			}
			if (wantTbl.Columns[i].Generated == nil) != (gotTbl.Columns[i].Generated == nil) {
				t.Fatalf("shuffle %d: table %q column %q generated mismatch", shuffleIdx, name, wantTbl.Columns[i].Name)
			}
		}
		// Indexes: same names.
		assertStringSetsEqual(t, "indexes of "+name, sdlIndexNames(wantTbl), sdlIndexNames(gotTbl), shuffleIdx)
		// FK constraints: same names.
		assertStringSetsEqual(t, "fks of "+name, fkNames(wantTbl), fkNames(gotTbl), shuffleIdx)
	}

	// Views must resolve identically: AnalyzedQuery populated in both.
	for name, wantView := range want.Views {
		gotView := got.Views[name]
		if gotView == nil {
			t.Fatalf("shuffle %d: view %q missing", shuffleIdx, name)
		}
		if (wantView.AnalyzedQuery == nil) != (gotView.AnalyzedQuery == nil) {
			t.Fatalf("shuffle %d: view %q AnalyzedQuery populated mismatch (want nil=%v got nil=%v)",
				shuffleIdx, name, wantView.AnalyzedQuery == nil, gotView.AnalyzedQuery == nil)
		}
	}
}

func keys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func sdlIndexNames(t *Table) []string {
	out := make([]string, 0, len(t.Indexes))
	for _, idx := range t.Indexes {
		out = append(out, strings.ToLower(idx.Name))
	}
	return out
}

func fkNames(t *Table) []string {
	var out []string
	for _, con := range t.Constraints {
		if con.Type == ConForeignKey {
			out = append(out, strings.ToLower(con.Name))
		}
	}
	return out
}

func assertKeySetsEqual(t *testing.T, label string, want, got []string, shuffleIdx int) {
	t.Helper()
	assertStringSetsEqual(t, label, want, got, shuffleIdx)
}

func assertStringSetsEqual(t *testing.T, label string, want, got []string, shuffleIdx int) {
	t.Helper()
	wantSet := make(map[string]bool, len(want))
	for _, w := range want {
		wantSet[w] = true
	}
	gotSet := make(map[string]bool, len(got))
	for _, g := range got {
		gotSet[g] = true
	}
	for w := range wantSet {
		if !gotSet[w] {
			t.Fatalf("shuffle %d: %s: missing %q (want %v, got %v)", shuffleIdx, label, w, want, got)
		}
	}
	for g := range gotSet {
		if !wantSet[g] {
			t.Fatalf("shuffle %d: %s: unexpected %q (want %v, got %v)", shuffleIdx, label, g, want, got)
		}
	}
}

// TestSDLCoverage enumerates the object kinds the loader must handle and asserts
// each loads into the catalog.
func TestSDLCoverage(t *testing.T) {
	sql := `
CREATE DATABASE app CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE app;

CREATE TABLE parent (
  id int PRIMARY KEY,
  code varchar(20) NOT NULL,
  UNIQUE KEY uq_code (code)
);

CREATE TABLE t (
  id int AUTO_INCREMENT PRIMARY KEY,
  parent_id int,
  name varchar(100) NOT NULL,
  price decimal(10,2) DEFAULT 0.00,
  total decimal(12,2) AS (price * id) VIRTUAL,
  tags varchar(255),
  body text,
  CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent (id) ON DELETE CASCADE,
  CONSTRAINT chk_price CHECK (price >= 0),
  KEY idx_name_prefix (name(10)),
  KEY idx_expr ((price + 1)),
  FULLTEXT KEY ft_body (body)
);

CREATE INDEX idx_parent ON t (parent_id);

CREATE VIEW v AS SELECT id, name FROM t;

CREATE FUNCTION fn_inc(x int) RETURNS int DETERMINISTIC RETURN x + 1;

CREATE PROCEDURE pr_noop() BEGIN SELECT 1; END;

CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW SET NEW.name = NEW.name;

CREATE EVENT ev ON SCHEDULE EVERY 1 DAY DO SET @y = 1;
`
	db := loadSDLForDB(t, sql, "app")

	// database
	if db.Name == "" {
		t.Fatal("expected database name set")
	}
	// table + columns
	tbl := db.GetTable("t")
	if tbl == nil {
		t.Fatal("expected table t")
	}
	if tbl.GetColumn("id") == nil || tbl.GetColumn("name") == nil {
		t.Fatal("expected columns id and name on t")
	}
	// generated column
	total := tbl.GetColumn("total")
	if total == nil || total.Generated == nil {
		t.Fatal("expected generated column total")
	}
	// indexes: prefix, functional, plus standalone CREATE INDEX
	if len(tbl.Indexes) == 0 {
		t.Fatal("expected indexes on t")
	}
	var haveStandalone bool
	for _, idx := range tbl.Indexes {
		if strings.EqualFold(idx.Name, "idx_parent") {
			haveStandalone = true
		}
	}
	if !haveStandalone {
		t.Fatal("expected standalone CREATE INDEX idx_parent on t")
	}
	// FK + CHECK constraints
	var haveFK, haveCheck bool
	for _, con := range tbl.Constraints {
		switch con.Type {
		case ConForeignKey:
			haveFK = true
		case ConCheck:
			haveCheck = true
		}
	}
	if !haveFK {
		t.Fatal("expected FK constraint on t")
	}
	if !haveCheck {
		t.Fatal("expected CHECK constraint on t")
	}
	// view
	if db.Views["v"] == nil {
		t.Fatal("expected view v")
	}
	// function + procedure
	if db.Functions["fn_inc"] == nil {
		t.Fatal("expected function fn_inc")
	}
	if db.Procedures["pr_noop"] == nil {
		t.Fatal("expected procedure pr_noop")
	}
	// trigger
	if db.Triggers["trg"] == nil {
		t.Fatal("expected trigger trg")
	}
	// event
	if db.Events["ev"] == nil {
		t.Fatal("expected event ev")
	}
}

// TestSDLEquivalentToLoadSQLWhenOrdered verifies that, for a schema whose
// statements are already in dependency order, LoadSDL yields the same catalog
// shape as LoadSQL (LoadSDL must not lose or alter objects).
func TestSDLEquivalentToLoadSQLWhenOrdered(t *testing.T) {
	sql := `
CREATE DATABASE app;
USE app;
CREATE TABLE parent (id int PRIMARY KEY);
CREATE TABLE child (id int PRIMARY KEY, parent_id int,
  CONSTRAINT fk FOREIGN KEY (parent_id) REFERENCES parent (id));
CREATE INDEX idx ON child (parent_id);
CREATE VIEW v AS SELECT id FROM child;
`
	sdlCat, err := LoadSDL(sql)
	if err != nil {
		t.Fatalf("LoadSDL: %v", err)
	}
	sqlCat, err := LoadSQL(sql)
	if err != nil {
		t.Fatalf("LoadSQL: %v", err)
	}
	sdlDB := sdlCat.GetDatabase("app")
	sqlDB := sqlCat.GetDatabase("app")
	if sdlDB == nil || sqlDB == nil {
		t.Fatal("expected app database from both loaders")
	}
	assertKeySetsEqual(t, "tables", keys(sqlDB.Tables), keys(sdlDB.Tables), 0)
	assertKeySetsEqual(t, "views", keys(sqlDB.Views), keys(sdlDB.Views), 0)
}

// TestSDLMultiDatabaseScoping verifies that unqualified objects in a multi-database
// SDL are created in the database that was current at their declaration position —
// not in whichever database the last USE happened to select. The topological sort
// must not float USE statements past the objects they scope.
func TestSDLMultiDatabaseScoping(t *testing.T) {
	sql := `
CREATE DATABASE db1;
CREATE DATABASE db2;
USE db1;
CREATE TABLE t1 (id int);
USE db2;
CREATE TABLE t2 (id int);
`
	c, err := LoadSDL(sql)
	if err != nil {
		t.Fatalf("LoadSDL: %v", err)
	}
	db1 := c.GetDatabase("db1")
	db2 := c.GetDatabase("db2")
	if db1 == nil || db2 == nil {
		t.Fatal("expected both db1 and db2")
	}
	if db1.GetTable("t1") == nil {
		t.Errorf("t1 must be in db1, db1 tables: %v", keys(db1.Tables))
	}
	if db2.GetTable("t2") == nil {
		t.Errorf("t2 must be in db2, db2 tables: %v", keys(db2.Tables))
	}
	// And not crossed.
	if db2.GetTable("t1") != nil {
		t.Errorf("t1 wrongly created in db2")
	}
	if db1.GetTable("t2") != nil {
		t.Errorf("t2 wrongly created in db1")
	}
}

// TestSDLCircularForeignKey verifies that two tables with mutually-referencing
// foreign keys load successfully. With foreign_key_checks disabled (the loader's
// strategy), a FK does not require its parent to exist, so a FK cycle is a valid
// MySQL schema and must not be reported as a dependency cycle.
func TestSDLCircularForeignKey(t *testing.T) {
	sql := `
CREATE DATABASE app;
USE app;
CREATE TABLE a (
  id int PRIMARY KEY,
  b_id int,
  CONSTRAINT fk_a FOREIGN KEY (b_id) REFERENCES b (id)
);
CREATE TABLE b (
  id int PRIMARY KEY,
  a_id int,
  CONSTRAINT fk_b FOREIGN KEY (a_id) REFERENCES a (id)
);
`
	db := loadSDLForDB(t, sql, "app")
	if db.GetTable("a") == nil || db.GetTable("b") == nil {
		t.Fatalf("expected both tables a and b, got %v", keys(db.Tables))
	}
}

// TestSDLSelfReferentialForeignKey verifies a table whose FK references itself
// loads cleanly (no spurious self-cycle).
func TestSDLSelfReferentialForeignKey(t *testing.T) {
	sql := `
CREATE DATABASE app;
USE app;
CREATE TABLE node (
  id int PRIMARY KEY,
  parent_id int,
  CONSTRAINT fk_self FOREIGN KEY (parent_id) REFERENCES node (id)
);
`
	db := loadSDLForDB(t, sql, "app")
	if db.GetTable("node") == nil {
		t.Fatal("expected table node")
	}
}

// TestSDLSetForeignKeyChecksDoesNotBreakForwardFK verifies that a SET
// foreign_key_checks statement inside the SDL (as emitted by mysqldump) does not
// re-enable FK validation and break a forward FK reference. The loader owns the
// FK-check state for the duration of the load.
func TestSDLSetForeignKeyChecksDoesNotBreakForwardFK(t *testing.T) {
	sql := `
SET FOREIGN_KEY_CHECKS = 1;
CREATE DATABASE app;
USE app;
CREATE TABLE child (
  id int PRIMARY KEY,
  parent_id int,
  CONSTRAINT fk FOREIGN KEY (parent_id) REFERENCES parent (id)
);
CREATE TABLE parent (id int PRIMARY KEY);
SET FOREIGN_KEY_CHECKS = 1;
`
	db := loadSDLForDB(t, sql, "app")
	if db.GetTable("child") == nil || db.GetTable("parent") == nil {
		t.Fatalf("expected child and parent, got %v", keys(db.Tables))
	}
}

// TestSDLViewCTEShadowsTable verifies that a top-level CTE whose name collides
// with a declared table does not create a spurious dependency edge to that table
// (the CTE is query-local scope and shadows the real object). The view must still
// load; the point is that no false ordering constraint is introduced.
func TestSDLViewCTEShadowsTable(t *testing.T) {
	sql := `
CREATE DATABASE app;
USE app;
CREATE TABLE parent (id int PRIMARY KEY);
CREATE VIEW v AS WITH parent AS (SELECT 1 AS id) SELECT id FROM parent;
`
	// extractSDLDeps must not emit an edge v -> parent (the FROM parent resolves to
	// the CTE, not the table).
	c, err := LoadSDL(sql)
	if err != nil {
		t.Fatalf("LoadSDL: %v", err)
	}
	db := c.GetDatabase("app")
	if db == nil || db.Views["v"] == nil || db.GetTable("parent") == nil {
		t.Fatal("expected app.v and app.parent")
	}
}

// TestSDLViewNestedCTEDoesNotShadowOuterView verifies that a CTE defined only in
// a nested subquery does NOT suppress an unqualified reference to a real VIEW in
// the outer query. View-on-view dependencies share a priority layer, so the edge
// (not priority) is what orders them: if the nested CTE wrongly shadowed the outer
// reference, the dependent view could be analyzed before the view it builds on.
func TestSDLViewNestedCTEDoesNotShadowOuterView(t *testing.T) {
	// base_v is a real view. dependent_v references base_v unqualified in its outer
	// query and also nests a CTE named base_v in a subquery. The outer reference
	// must bind to the real view, forcing dependent_v to be ordered after base_v.
	// dependent_v is declared FIRST, so only a correct edge yields a populated
	// AnalyzedQuery.
	sql := `
CREATE DATABASE app;
USE app;
CREATE VIEW dependent_v AS
  SELECT id FROM base_v
  WHERE id IN (WITH base_v AS (SELECT 7 AS id) SELECT id FROM base_v);
CREATE TABLE t (id int PRIMARY KEY);
CREATE VIEW base_v AS SELECT id FROM t;
`
	db := loadSDLForDB(t, sql, "app")
	dep := db.Views["dependent_v"]
	base := db.Views["base_v"]
	if dep == nil || base == nil {
		t.Fatal("expected both dependent_v and base_v")
	}
	if dep.AnalyzedQuery == nil {
		t.Fatal("dependent_v must be ordered after base_v: its outer FROM base_v binds to the real view, not the nested CTE of the same name")
	}
}

// TestSDLRejectsCTAS verifies CREATE TABLE ... AS SELECT is rejected: it is an
// imperative construct the loader cannot faithfully materialize as a schema
// snapshot, and the catalog silently drops it otherwise.
func TestSDLRejectsCTAS(t *testing.T) {
	sql := `
CREATE DATABASE app;
USE app;
CREATE TABLE src (id int);
CREATE TABLE derived AS SELECT id FROM src;
`
	_, err := LoadSDL(sql)
	if err == nil {
		t.Fatal("expected CTAS to be rejected")
	}
	if !strings.Contains(err.Error(), "SELECT") {
		t.Fatalf("expected CTAS rejection error mentioning SELECT, got: %v", err)
	}
}

// TestSDLAdjacentStringLiterals covers MySQL's adjacent string-literal
// concatenation ('a' 'b' → 'ab') through the SDL loader. The stock MySQL 8.0
// sys schema relies on it: ps_trace_thread's body concatenates '\n' with the
// next segment across a line break, so a canonical dump of sys was unloadable
// before adjacency support. Bodies are stored verbatim; the loaded catalog must
// also self-diff empty so the adjacency never phantom-diffs.
func TestSDLAdjacentStringLiterals(t *testing.T) {
	sql := `
CREATE DATABASE app;
USE app;
CREATE TABLE t (a varchar(40) DEFAULT 'x' 'y');
CREATE PROCEDURE p()
BEGIN
    SELECT CONCAT('tmp disk tables: ', 3, '\n'
                  'select scan: ', 4, '\n');
END;
`
	db := loadSDLForDB(t, sql, "app")
	if db.Procedures["p"] == nil {
		t.Fatal("procedure p not loaded")
	}
	if !strings.Contains(db.Procedures["p"].Body, "'select scan: '") {
		t.Errorf("procedure body not stored verbatim: %q", db.Procedures["p"].Body)
	}

	c, err := LoadSDL(sql)
	if err != nil {
		t.Fatalf("LoadSDL error: %v", err)
	}
	if d := Diff(c, c); !d.IsEmpty() {
		t.Errorf("self-diff of SDL with adjacent literals not empty")
	}
}

// TestSDLIntroducerDefaultRendersBareValue pins that a charset-introduced
// string default — single (_utf8mb4'x') or folded (_utf8mb4'x' 'y') — is
// stored and rendered by VALUE, the way MySQL stores it (SHOW CREATE TABLE
// echoes DEFAULT 'x' / DEFAULT 'xy'; oracle 8.0.32 + 5.7.25). Deparsing the
// introducer into the default rendered DEFAULT '_utf8mb4'x” — invalid SQL
// the engine rejects with 1064 on apply.
func TestSDLIntroducerDefaultRendersBareValue(t *testing.T) {
	target, err := LoadSDL(`
CREATE DATABASE d1;
USE d1;
CREATE TABLE t (
    a varchar(10) CHARSET utf8mb4 DEFAULT _utf8mb4'x' 'y',
    b varchar(10) CHARSET utf8mb4 DEFAULT _utf8mb4'x'
);
`)
	if err != nil {
		t.Fatalf("LoadSDL error: %v", err)
	}
	empty, err := LoadSDL("CREATE DATABASE d1;\n")
	if err != nil {
		t.Fatalf("LoadSDL(empty) error: %v", err)
	}
	sql := GenerateMigration(empty, target, Diff(empty, target)).SQL()
	if !strings.Contains(sql, "DEFAULT 'xy'") || !strings.Contains(sql, "DEFAULT 'x'") {
		t.Errorf("generated DDL must carry bare folded defaults, got:\n%s", sql)
	}
	if strings.Contains(sql, "_utf8mb4") {
		t.Errorf("generated DDL leaked the charset introducer into a default:\n%s", sql)
	}
	if d := Diff(target, target); !d.IsEmpty() {
		t.Errorf("self-diff of introducer-default SDL not empty")
	}
}

// TestSDLParenSubqueryOperand covers a parenthesized scalar subquery used as
// the LEFT operand of a binary operator with the comparison itself inside
// parentheses — the stock MySQL 8.0 sys.metrics shape:
//
//	if(((select count(0) from performance_schema.setup_instruments
//	     where (... and ...))) = 0),'NO','YES')
//
// The paren-expression parser used to classify the leading '(' run as a
// parenthesized query expression and hard-expect ')' after the subquery, so
// the `= 0` was "unexpected token" and any dump containing sys.metrics was
// unloadable. The loaded catalog must also self-diff empty: the engine
// collapses redundant parens around a subquery to one pair in stored bodies
// (oracle 8.0.32 + 5.7.25), so the depth-2 user form and the collapsed stored
// form must canonicalize equal (see TestOracle_ViewIdempotence probes).
func TestSDLParenSubqueryOperand(t *testing.T) {
	sql := `
CREATE DATABASE app;
USE app;
CREATE TABLE t (a int);
CREATE VIEW v1 AS SELECT if(((SELECT count(0) FROM t WHERE (a > 0 AND a < 9)) = 0),'NO','YES') AS x;
CREATE VIEW v2 AS SELECT if((((SELECT count(0) FROM t)) = 0),'NO','YES') AS x;
CREATE VIEW v3 AS SELECT ((SELECT max(a) FROM t) + 1) AS x, (1 IN ((SELECT max(a) FROM t) + 1, 2)) AS y;
`
	db := loadSDLForDB(t, sql, "app")
	for _, name := range []string{"v1", "v2", "v3"} {
		if db.Views[name] == nil {
			t.Fatalf("view %s not loaded", name)
		}
	}
	// v1 and v2 differ only in redundant subquery parens; both canonicalize
	// to the single-pair form the engine stores.
	if b1, b2 := db.Views["v1"].Definition, db.Views["v2"].Definition; !strings.Contains(b2, "((select count(0) from `t`) = 0)") {
		t.Errorf("v2 wrapper parens not collapsed: %q (v1: %q)", b2, b1)
	}

	c, err := LoadSDL(sql)
	if err != nil {
		t.Fatalf("LoadSDL error: %v", err)
	}
	if d := Diff(c, c); !d.IsEmpty() {
		t.Errorf("self-diff of SDL with paren-subquery operands not empty")
	}
}
