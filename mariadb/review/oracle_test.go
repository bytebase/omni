package review

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	tcmariadb "github.com/testcontainers/testcontainers-go/modules/mariadb"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
)

// oracleBase is the schema every oracle case starts from.
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
`

var mariadbOracle struct {
	once      sync.Once
	container testcontainers.Container
	db        *sql.DB
	err       error
}

func TestMain(m *testing.M) {
	code := m.Run()
	if mariadbOracle.container != nil {
		_ = mariadbOracle.db.Close()
		_ = testcontainers.TerminateContainer(mariadbOracle.container)
	}
	os.Exit(code)
}

// startMariaDB returns the MariaDB 11.8 the oracle tests of this package
// ask, started once.
func startMariaDB(t *testing.T) *sql.DB {
	t.Helper()
	mariadbOracle.once.Do(func() {
		ctx := context.Background()
		c, err := tcmariadb.Run(ctx, "mariadb:11.8.8", tcmariadb.WithDatabase("shop"), tcmariadb.WithUsername("root"), tcmariadb.WithPassword("test"))
		if err != nil {
			mariadbOracle.err = err
			return
		}
		mariadbOracle.container = c
		dsn, err := c.ConnectionString(ctx, "multiStatements=true")
		if err != nil {
			mariadbOracle.err = err
			return
		}
		mariadbOracle.db, mariadbOracle.err = sql.Open("mysql", dsn)
		if mariadbOracle.err == nil {
			mariadbOracle.err = mariadbOracle.db.PingContext(ctx)
		}
	})
	if mariadbOracle.err != nil {
		t.Fatalf("mariadb oracle unavailable: %v", mariadbOracle.err)
	}
	return mariadbOracle.db
}

// resetOracle recreates the shop database from oracleBase, and an empty
// bbdataarchive, on a connection of its own with shop selected.
func resetOracle(t *testing.T, db *sql.DB) *sql.Conn {
	t.Helper()
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	for _, stmt := range []string{
		"DROP DATABASE IF EXISTS shop", "DROP DATABASE IF EXISTS bbdataarchive", "DROP DATABASE IF EXISTS other",
		"CREATE DATABASE shop", "CREATE DATABASE bbdataarchive", "USE shop", oracleBase,
	} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	return conn
}

// syncSchema reads the shop database into metadata, the parts the rules
// read: tables with their columns, indexes, foreign keys, checks, and
// triggers, and views.
func syncSchema(t *testing.T, conn *sql.Conn) *metadata.DatabaseSchemaMetadata {
	t.Helper()
	ctx := context.Background()
	query := func(q string, scan func(*sql.Rows) error) {
		rows, err := conn.QueryContext(ctx, q)
		if err != nil {
			t.Fatalf("%s: %v", q, err)
		}
		defer rows.Close()
		for rows.Next() {
			if err := scan(rows); err != nil {
				t.Fatal(err)
			}
		}
	}
	s := &metadata.SchemaMetadata{}
	tables := make(map[string]*metadata.TableMetadata)
	var order []string
	query(`SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME`, func(r *sql.Rows) error {
		var name, kind string
		if err := r.Scan(&name, &kind); err != nil {
			return err
		}
		if kind == "VIEW" {
			s.Views = append(s.Views, &metadata.ViewMetadata{Name: name})
			return nil
		}
		tables[name] = &metadata.TableMetadata{Name: name}
		order = append(order, name)
		return nil
	})
	query(`SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION`, func(r *sql.Rows) error {
		var tbl, col string
		if err := r.Scan(&tbl, &col); err != nil {
			return err
		}
		if t, ok := tables[tbl]; ok {
			t.Columns = append(t.Columns, &metadata.ColumnMetadata{Name: col})
		}
		return nil
	})
	query(`SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME`, func(r *sql.Rows) error {
		var tbl, idx string
		var nonUnique int
		if err := r.Scan(&tbl, &idx, &nonUnique); err != nil {
			return err
		}
		if t, ok := tables[tbl]; ok {
			t.Indexes = append(t.Indexes, &metadata.IndexMetadata{Name: idx, Primary: idx == "PRIMARY", Unique: nonUnique == 0})
		}
		return nil
	})
	query(`SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'`, func(r *sql.Rows) error {
		var tbl, name, ref string
		if err := r.Scan(&tbl, &name, &ref); err != nil {
			return err
		}
		if t, ok := tables[tbl]; ok {
			t.ForeignKeys = append(t.ForeignKeys, &metadata.ForeignKeyMetadata{Name: name, ReferencedTable: ref})
		}
		return nil
	})
	query(`SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' AND CONSTRAINT_TYPE = 'CHECK'`, func(r *sql.Rows) error {
		var tbl, name string
		if err := r.Scan(&tbl, &name); err != nil {
			return err
		}
		if t, ok := tables[tbl]; ok {
			t.CheckConstraints = append(t.CheckConstraints, &metadata.CheckConstraintMetadata{Name: name})
		}
		return nil
	})
	query(`SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'`, func(r *sql.Rows) error {
		var tbl, name string
		if err := r.Scan(&tbl, &name); err != nil {
			return err
		}
		if t, ok := tables[tbl]; ok {
			t.Triggers = append(t.Triggers, &metadata.TriggerMetadata{Name: name})
		}
		return nil
	})
	for _, name := range order {
		s.Tables = append(s.Tables, tables[name])
	}
	return &metadata.DatabaseSchemaMetadata{Name: "shop", Schemas: []*metadata.SchemaMetadata{s}}
}

// lowerCaseTableNames reads the server's lower_case_table_names.
func lowerCaseTableNames(t *testing.T, conn *sql.Conn) int {
	t.Helper()
	var v int
	if err := conn.QueryRowContext(context.Background(), "SELECT @@lower_case_table_names").Scan(&v); err != nil {
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

// TestPriorBackupAgainstMariaDB builds each statement's backup query as
// Bytebase does and runs it on MariaDB before the change. A PriorBackup
// finding on a statement must mean that query fails; for a complete case,
// every query that fails must have a finding.
func TestPriorBackupAgainstMariaDB(t *testing.T) {
	db := startMariaDB(t)
	tests := []struct {
		name     string
		change   string
		complete bool
	}{
		{"a plain backup", "UPDATE t SET n = 5 WHERE id = 1;\nDELETE FROM nokey WHERE a = 1;", true},
		{"a table the change creates", "CREATE TABLE n2 (id INT PRIMARY KEY);\nUPDATE n2 SET id = 2 WHERE id = 1;", true},
		{"a table the change creates, joined", "CREATE TABLE n2 (id INT PRIMARY KEY);\nDELETE t FROM t JOIN n2 ON t.id = n2.id WHERE t.n = 1;", true},
		{"a table the change creates, in a subquery", "CREATE TABLE n2 (id INT PRIMARY KEY);\nDELETE FROM t WHERE id IN (SELECT id FROM n2);", true},
		{"a table the change renames to", "RENAME TABLE nokey TO nk;\nDELETE FROM nk WHERE a = 1;", true},
		{"a column the change adds, unqualified", "ALTER TABLE t ADD COLUMN gone INT;\nUPDATE t SET n = 1 WHERE gone = 1;", true},
		{"a column the change adds, qualified", "ALTER TABLE nokey ADD COLUMN gone INT;\nDELETE FROM nokey WHERE nokey.gone = 1;", true},
		{"a column the change adds, in ORDER BY", "ALTER TABLE nokey ADD COLUMN gone INT;\nDELETE FROM nokey WHERE a = 1 ORDER BY gone LIMIT 1;", true},
		{"a column the change renames to, through an alias", "ALTER TABLE p RENAME COLUMN note TO memo;\nUPDATE p AS x SET x.id = 3 WHERE x.memo = 'a';", true},
		{"a column the change changes the name of", "ALTER TABLE p CHANGE note memo VARCHAR(20);\nDELETE FROM p WHERE memo = 'zz';", true},
		{"a column of a joined table", "ALTER TABLE p ADD COLUMN gone INT;\nUPDATE t JOIN p ON p.gone = t.n SET t.n = 1 WHERE t.id = 3;", true},
		{"a column of a table the change recreates", "DROP TABLE nokey;\nCREATE TABLE nokey (b INT);\nDELETE FROM nokey WHERE b = 1;", true},
		{"a column in a subquery, not judged", "ALTER TABLE t ADD COLUMN gone INT;\nUPDATE plain SET v = 1 WHERE id IN (SELECT gone FROM t);", false},
		{"a column the SET clause names, not copied", "ALTER TABLE t ADD COLUMN gone INT;\nUPDATE t SET gone = 1 WHERE id = 1;", true},
		{"a column USING names", "ALTER TABLE t ADD COLUMN note INT;\nUPDATE t JOIN p USING (note) SET t.n = 1 WHERE t.id = 1;", true},
		{"a column USING names, a derived table on its side may hold", "ALTER TABLE t ADD COLUMN note INT;\nUPDATE t JOIN (SELECT 1 AS note) d ON d.note = t.id JOIN p USING (note) SET t.n = 1 WHERE t.id = 3;", true},
		{"statements before a USE", "ALTER TABLE t ADD COLUMN gone INT;\nUPDATE t SET n = 1 WHERE gone = 1;\nUSE bbdataarchive;\nUPDATE shop.t SET n = 2 WHERE id = 1;", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := resetOracle(t, db)
			schema := syncSchema(t, conn)
			result, err := Review(context.Background(), tt.change, review.Options{
				Rules:  []review.Rule{review.PriorBackup},
				Change: review.Change{PriorBackup: true, MaxBackupSize: 1 << 20},
			}, []review.Target{{Schema: schema, BackupDatabaseExists: true, LowerCaseTableNames: lowerCaseTableNames(t, conn)}})
			if err != nil {
				t.Fatalf("Review: %v", err)
			}
			reported := make(map[int]string)
			for _, f := range result.Findings {
				if f.Rule == review.Syntax {
					t.Fatalf("syntax: %s", f.Message)
				}
				reported[f.Statement] = f.Message
			}
			fails := backupFailures(t, conn, tt.change)
			for i, msg := range reported {
				if _, ok := fails[i]; !ok {
					t.Errorf("statement %d reported, but its backup runs: %s", i, msg)
				}
			}
			if tt.complete {
				for i, err := range fails {
					if _, ok := reported[i]; !ok {
						t.Errorf("statement %d not reported, but its backup fails: %v", i, err)
					}
				}
			}
		})
	}
}

// backupFailures runs each top-level UPDATE and DELETE's backup query, as
// Bytebase builds it, before the change, and returns the statements whose
// query fails, with the error.
//
// bytebase: backend/plugin/parser/mysql/backup.go — generateSQLForTable
func backupFailures(t *testing.T, conn *sql.Conn, change string) map[int]error {
	t.Helper()
	ctx := context.Background()
	fails := make(map[int]error)
	for n, s := range parseChange(t, change) {
		var queries []string
		for i, q := range backupQueries(change, s.node) {
			target := fmt.Sprintf("bk_%d_%d", n, i)
			queries = append(queries, fmt.Sprintf("CREATE TABLE `bbdataarchive`.`%s` LIKE `shop`.`%s`", target, q.table), fmt.Sprintf("INSERT INTO `bbdataarchive`.`%s` %s", target, q.query))
		}
		for _, q := range queries {
			if _, err := conn.ExecContext(ctx, q); err != nil {
				fails[s.index] = fmt.Errorf("%s: %w", q, err)
				break
			}
		}
	}
	return fails
}

type backupQuery struct {
	table string
	query string
}

// backupQueries returns the SELECTs Bytebase copies an UPDATE or DELETE's
// rows with, one per table it backs up: the tables, as written, and what
// follows the SET list or the tables, as written.
func backupQueries(sql string, n ast.Node) []backupQuery {
	text := func(from, to int) string { return sql[from:to] }
	exprStart := func(e ast.TableExpr) int { return nodeLoc(e).Start }
	exprEnd := func(e ast.TableExpr) int { return nodeLoc(e).End }
	var out []backupQuery
	add := func(name, alias, from string) {
		sel := name
		if alias != "" {
			sel = alias
		}
		out = append(out, backupQuery{table: name, query: "SELECT `" + sel + "`.* FROM " + from})
	}
	refs := func(exprs []ast.TableExpr) []*ast.TableRef {
		var list []*ast.TableRef
		for _, e := range exprs {
			ast.Inspect(e, func(n ast.Node) bool {
				if r, ok := n.(*ast.TableRef); ok {
					list = append(list, r)
				}
				_, isSub := n.(*ast.SubqueryExpr)
				return !isSub
			})
		}
		return list
	}
	switch v := n.(type) {
	case *ast.UpdateStmt:
		if len(v.Tables) == 0 || len(v.SetList) == 0 {
			return nil
		}
		last := v.SetList[len(v.SetList)-1]
		from := strings.TrimSpace(text(exprStart(v.Tables[0]), exprEnd(v.Tables[len(v.Tables)-1])))
		if suffix := strings.TrimSpace(text(nodeLoc(last).End, nodeLoc(v).End)); suffix != "" {
			from += " " + suffix
		}
		targets := make(map[string]bool)
		for _, a := range v.SetList {
			if a != nil && a.Column != nil && a.Column.Table != "" {
				targets[strings.ToLower(a.Column.Table)] = true
			}
		}
		for _, r := range refs(v.Tables) {
			key := strings.ToLower(r.Name)
			if r.Alias != "" {
				key = strings.ToLower(r.Alias)
			}
			if len(targets) == 0 || targets[key] {
				add(r.Name, r.Alias, from)
			}
		}
	case *ast.DeleteStmt:
		list := v.Tables
		if len(v.Using) > 0 {
			list = v.Using
		}
		if len(list) == 0 {
			return nil
		}
		from := strings.TrimSpace(text(exprStart(list[0]), nodeLoc(v).End))
		if len(v.Using) == 0 {
			for _, r := range refs(v.Tables) {
				add(r.Name, r.Alias, from)
				break
			}
			return out
		}
		using := refs(v.Using)
		for _, target := range refs(v.Tables) {
			for _, r := range using {
				if strings.EqualFold(r.Alias, target.Name) || r.Alias == "" && strings.EqualFold(r.Name, target.Name) {
					add(r.Name, r.Alias, from)
					break
				}
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].table < out[j].table })
	return out
}

// TestScanRulesAgainstMariaDB runs each change on MariaDB, statement by
// statement, and derives from information_schema what the constraint
// rules should report: DisallowDropConstraint, a primary key, foreign key,
// unique key, or check constraint a statement removes; RequirePrimaryKey,
// a table a CREATE TABLE made, or one that had a primary key, that ends
// the change without one. Every finding must match; for a complete case,
// every expected finding must be reported.
func TestScanRulesAgainstMariaDB(t *testing.T) {
	db := startMariaDB(t)
	tests := []struct {
		name     string
		change   string
		complete bool
	}{
		{"constraint drops the SQL names", "ALTER TABLE t DROP PRIMARY KEY, DROP FOREIGN KEY fk_p, DROP CONSTRAINT ck_n;\nALTER TABLE t DROP CONSTRAINT uk_code;", true},
		{"a unique index dropped by name", "ALTER TABLE t DROP INDEX uk_code;\nDROP INDEX idx_n ON t;", true},
		{"a unique key the change makes, then drops", "ALTER TABLE nokey ADD UNIQUE KEY uk_a (a);\nALTER TABLE nokey DROP INDEX uk_a;", true},
		{"DROP INDEX of the primary key", "DROP INDEX `PRIMARY` ON plain;", true},
		{"a table created without a key", "CREATE TABLE a (id INT);\nCREATE TABLE b (id INT PRIMARY KEY);", true},
		{"a key added later", "CREATE TABLE a (id INT);\nALTER TABLE a ADD PRIMARY KEY (id);", true},
		{"a primary key replaced", "ALTER TABLE plain DROP PRIMARY KEY, ADD PRIMARY KEY (id, v);", true},
		{"a primary key dropped, then added back", "ALTER TABLE plain DROP PRIMARY KEY;\nALTER TABLE plain ADD PRIMARY KEY (v);", true},
		{"CREATE OR REPLACE TABLE without a key", "CREATE OR REPLACE TABLE plain (id INT);", false},
		{"a table renamed after losing its key", "ALTER TABLE plain DROP PRIMARY KEY;\nRENAME TABLE plain TO plain2;", false},
		{"CREATE TABLE IF NOT EXISTS of a table with a key", "CREATE TABLE IF NOT EXISTS plain (id INT);", true},
		{"CREATE TABLE IF NOT EXISTS, not judged", "CREATE TABLE IF NOT EXISTS a (id INT);", false},
		{"a table whose database the change drops", "CREATE DATABASE d;\nCREATE TABLE d.e (id INT);\nDROP DATABASE d;", true},
		{"an index of a table CREATE TABLE IF NOT EXISTS leaves", "CREATE TABLE IF NOT EXISTS t (id INT);\nDROP INDEX uk_code ON t;", true},
		{"an index of a table the change recreates", "DROP TABLE t;\nCREATE TABLE t (code INT, KEY uk_code (code));\nALTER TABLE t DROP INDEX uk_code;", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := resetOracle(t, db)
			schema := syncSchema(t, conn)
			result, err := Review(context.Background(), tt.change, review.Options{
				Rules: []review.Rule{review.DisallowDropConstraint, review.RequirePrimaryKey},
			}, []review.Target{{Schema: schema, BackupDatabaseExists: true, LowerCaseTableNames: lowerCaseTableNames(t, conn)}})
			if err != nil {
				t.Fatalf("Review: %v", err)
			}
			got := make(map[string]bool)
			for _, f := range result.Findings {
				if f.Rule == review.Syntax {
					t.Fatalf("syntax: %s", f.Message)
				}
				got[fmt.Sprintf("%s %d", f.Rule, f.Statement)] = true
			}
			want := scanTruth(t, conn, tt.change)
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
func constraintSet(t *testing.T, conn *sql.Conn) map[string]bool {
	t.Helper()
	ctx := context.Background()
	out := make(map[string]bool)
	rows, err := conn.QueryContext(ctx, `SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tbl, kind, name, def string
		if err := rows.Scan(&tbl, &kind, &name, &def); err != nil {
			t.Fatal(err)
		}
		out[strings.ToLower(tbl+"/"+kind+"/"+name+"/"+def)] = true
	}
	return out
}

// primaryKeyed reads the tables of the shop database that have a primary
// key, and all its base tables.
func primaryKeyed(t *testing.T, conn *sql.Conn) (keyed, all map[string]bool) {
	t.Helper()
	ctx := context.Background()
	keyed, all = make(map[string]bool), make(map[string]bool)
	rows, err := conn.QueryContext(ctx, `SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var has bool
		if err := rows.Scan(&name, &has); err != nil {
			t.Fatal(err)
		}
		all[strings.ToLower(name)] = true
		if has {
			keyed[strings.ToLower(name)] = true
		}
	}
	return keyed, all
}

// scanTruth runs the change statement by statement and returns the
// findings the server's catalog calls for, as "RULE statement".
func scanTruth(t *testing.T, conn *sql.Conn, change string) map[string]bool {
	t.Helper()
	ctx := context.Background()
	want := make(map[string]bool)
	keyedBefore, _ := primaryKeyed(t, conn)
	// lostAt is the statement that removed a table's key, created the
	// statement that made a table.
	lostAt, createdAt := make(map[string]int), make(map[string]int)
	for _, s := range parseChange(t, change) {
		before := constraintSet(t, conn)
		keyed, _ := primaryKeyed(t, conn)
		if _, err := conn.ExecContext(ctx, change[s.loc.Start:s.loc.End]); err != nil {
			t.Fatalf("statement %d: %v", s.index, err)
		}
		after := constraintSet(t, conn)
		for k := range before {
			if !after[k] {
				want[fmt.Sprintf("%s %d", review.DisallowDropConstraint, s.index)] = true
			}
		}
		keyedNow, _ := primaryKeyed(t, conn)
		for name := range keyed {
			if !keyedNow[name] {
				lostAt[name] = s.index
			}
		}
		if v, ok := s.node.(*ast.CreateTableStmt); ok && !v.Temporary && v.Like == nil && v.Select == nil {
			createdAt[strings.ToLower(v.Table.Name)] = s.index
		}
	}
	keyedEnd, all := primaryKeyed(t, conn)
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
	return want
}
