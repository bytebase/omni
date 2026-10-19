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
	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/ast"
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
CREATE VIEW pv AS SELECT id, note FROM p;
`

var mysqlOracle struct {
	once      sync.Once
	container testcontainers.Container
	db        *sql.DB
	err       error
}

func TestMain(m *testing.M) {
	code := m.Run()
	if mysqlOracle.container != nil {
		_ = mysqlOracle.db.Close()
		_ = testcontainers.TerminateContainer(mysqlOracle.container)
	}
	os.Exit(code)
}

// startMySQL returns the MySQL 8.0 the oracle tests of this package ask,
// started once.
func startMySQL(t *testing.T) *sql.DB {
	t.Helper()
	mysqlOracle.once.Do(func() {
		ctx := context.Background()
		c, err := tcmysql.Run(ctx, "mysql:8.0", tcmysql.WithDatabase("shop"), tcmysql.WithUsername("root"), tcmysql.WithPassword("test"))
		if err != nil {
			mysqlOracle.err = err
			return
		}
		mysqlOracle.container = c
		dsn, err := c.ConnectionString(ctx, "multiStatements=true")
		if err != nil {
			mysqlOracle.err = err
			return
		}
		mysqlOracle.db, mysqlOracle.err = sql.Open("mysql", dsn)
		if mysqlOracle.err == nil {
			mysqlOracle.err = mysqlOracle.db.PingContext(ctx)
		}
	})
	if mysqlOracle.err != nil {
		t.Fatalf("mysql oracle unavailable: %v", mysqlOracle.err)
	}
	return mysqlOracle.db
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
		// A case may have turned foreign key checks off on the pooled
		// connection; other may reference shop, which it then keeps from
		// being dropped.
		"SET foreign_key_checks = 1",
		"DROP DATABASE IF EXISTS other", "DROP DATABASE IF EXISTS shop", "DROP DATABASE IF EXISTS bbdataarchive",
		"CREATE DATABASE shop", "CREATE DATABASE bbdataarchive", "USE shop", oracleBase,
	} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	return conn
}

// syncSchema reads the shop database into metadata, the parts the rules
// read, as Bytebase syncs them: tables with their columns and generation
// expressions, indexes with their key parts, foreign keys with their
// columns, checks with their expressions, triggers, and partitions; and
// views.
//
// bytebase: backend/plugin/db/mysql/sync.go — SyncDBSchema
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
	query(`SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION`, func(r *sql.Rows) error {
		var tbl, col, gen, extra string
		if err := r.Scan(&tbl, &col, &gen, &extra); err != nil {
			return err
		}
		if t, ok := tables[tbl]; ok {
			c := &metadata.ColumnMetadata{Name: col}
			switch {
			case gen == "":
			case strings.Contains(strings.ToUpper(extra), "VIRTUAL GENERATED"):
				c.Generation = &metadata.GenerationMetadata{Type: metadata.GenerationMetadata_TYPE_VIRTUAL, Expression: gen}
			case strings.Contains(strings.ToUpper(extra), "STORED GENERATED"):
				c.Generation = &metadata.GenerationMetadata{Type: metadata.GenerationMetadata_TYPE_STORED, Expression: gen}
			}
			t.Columns = append(t.Columns, c)
		}
		return nil
	})
	// A key part is its column, or the expression of a functional one in
	// parentheses.
	indexes := make(map[[2]string]*metadata.IndexMetadata)
	query(`SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX`, func(r *sql.Rows) error {
		var tbl, idx, part string
		var nonUnique int
		if err := r.Scan(&tbl, &idx, &nonUnique, &part); err != nil {
			return err
		}
		t, ok := tables[tbl]
		if !ok {
			return nil
		}
		i, ok := indexes[[2]string{tbl, idx}]
		if !ok {
			i = &metadata.IndexMetadata{Name: idx, Primary: idx == "PRIMARY", Unique: nonUnique == 0}
			indexes[[2]string{tbl, idx}] = i
			t.Indexes = append(t.Indexes, i)
		}
		i.Expressions = append(i.Expressions, part)
		return nil
	})
	// A foreign key names the table it references without its database.
	foreignKeys := make(map[[2]string]*metadata.ForeignKeyMetadata)
	query(`SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION`, func(r *sql.Rows) error {
		var tbl, name, ref, col, refCol string
		if err := r.Scan(&tbl, &name, &ref, &col, &refCol); err != nil {
			return err
		}
		t, ok := tables[tbl]
		if !ok {
			return nil
		}
		fk, ok := foreignKeys[[2]string{tbl, name}]
		if !ok {
			fk = &metadata.ForeignKeyMetadata{Name: name, ReferencedTable: ref}
			foreignKeys[[2]string{tbl, name}] = fk
			t.ForeignKeys = append(t.ForeignKeys, fk)
		}
		fk.Columns = append(fk.Columns, col)
		fk.ReferencedColumns = append(fk.ReferencedColumns, refCol)
		return nil
	})
	query(`SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK'`, func(r *sql.Rows) error {
		var tbl, name, clause string
		if err := r.Scan(&tbl, &name, &clause); err != nil {
			return err
		}
		if t, ok := tables[tbl]; ok {
			t.CheckConstraints = append(t.CheckConstraints, &metadata.CheckConstraintMetadata{Name: name, Expression: clause})
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
	query(`SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION`, func(r *sql.Rows) error {
		var tbl, name, expr string
		if err := r.Scan(&tbl, &name, &expr); err != nil {
			return err
		}
		if t, ok := tables[tbl]; ok {
			t.Partitions = append(t.Partitions, &metadata.TablePartitionMetadata{Name: name, Expression: expr})
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

// TestPriorBackupAgainstMySQL builds each statement's backup query as
// Bytebase does and runs it on MySQL before the change. A PriorBackup
// finding on a statement must mean that query fails; for a complete case,
// every query that fails must have a finding.
func TestPriorBackupAgainstMySQL(t *testing.T) {
	db := startMySQL(t)
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
		{"a view as the target", "UPDATE pv SET note = 'z' WHERE id = 1;", true},
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

// TestScanRulesAgainstMySQL runs each change on MySQL, statement by
// statement, and derives from information_schema what the constraint
// rules should report: DisallowDropConstraint, a primary key, foreign key,
// unique key, or check constraint a statement removes; RequirePrimaryKey,
// a table a CREATE TABLE made, or one that had a primary key, that ends
// the change without one. Every finding must match; for a complete case,
// every expected finding must be reported.
func TestScanRulesAgainstMySQL(t *testing.T) {
	db := startMySQL(t)
	tests := []struct {
		name     string
		change   string
		complete bool
	}{
		{"constraint drops the SQL names", "ALTER TABLE t DROP PRIMARY KEY, DROP FOREIGN KEY fk_p, DROP CHECK ck_n;\nALTER TABLE t DROP CONSTRAINT uk_code;", true},
		{"a unique index dropped by name", "ALTER TABLE t DROP INDEX uk_code;\nDROP INDEX idx_n ON t;", true},
		{"a unique key the change makes, then drops", "ALTER TABLE nokey ADD UNIQUE KEY uk_a (a);\nALTER TABLE nokey DROP INDEX uk_a;", true},
		{"DROP INDEX of the primary key", "DROP INDEX `PRIMARY` ON plain;", true},
		{"a table created without a key", "CREATE TABLE a (id INT);\nCREATE TABLE b (id INT PRIMARY KEY);", true},
		{"CREATE TABLE AS SELECT, which copies no key", "CREATE TABLE c AS SELECT id, v FROM plain;\nCREATE TABLE d (PRIMARY KEY (id)) AS SELECT id FROM plain;", true},
		{"a key added later", "CREATE TABLE a (id INT);\nALTER TABLE a ADD PRIMARY KEY (id);", true},
		{"a primary key replaced", "ALTER TABLE plain DROP PRIMARY KEY, ADD PRIMARY KEY (id, v);", true},
		{"a primary key dropped, then added back", "ALTER TABLE plain DROP PRIMARY KEY;\nALTER TABLE plain ADD PRIMARY KEY (v);", true},
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
		if v, ok := s.node.(*ast.CreateTableStmt); ok && !v.Temporary && v.Like == nil {
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

// TestColumnDropAgainstMySQL runs ALTER TABLE ... DROP COLUMN on a table
// x's primary key on MySQL, and pins what the server does with it: the
// error it stops the change with, if any, and the key x has after. MySQL
// removes a dropped column from every index, and drops an index with the
// last of its columns, so a drop of every column of the key leaves x
// without one, and a drop of some leaves it on the rest. A RequirePrimaryKey
// finding must match: reported only when the change ran and x ended it
// without a key, except where the server refuses for something the synced
// schema does not hold.
//
// mysql: https://dev.mysql.com/doc/refman/8.0/en/alter-table.html#alter-table-add-drop-column
func TestColumnDropAgainstMySQL(t *testing.T) {
	db := startMySQL(t)
	tests := []struct {
		name   string
		setup  string
		change string
		// err is a part of the error the server stops the change with, ""
		// when the change runs, and key the columns of x's primary key after,
		// "" when it has none.
		err string
		key string
		// reported is whether the review reports x's key removed, and unseen
		// that the server refuses the change for something the synced schema
		// does not hold.
		reported bool
		unseen   bool
	}{
		// (a) A key of one column goes with it.
		{name: "a key of one column", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); INSERT INTO x VALUES (1, 1), (2, 2)", change: "ALTER TABLE x DROP COLUMN id", reported: true},
		{name: "a key of one AUTO_INCREMENT column", setup: "CREATE TABLE x (id INT AUTO_INCREMENT PRIMARY KEY, v INT); INSERT INTO x (v) VALUES (1), (2)", change: "ALTER TABLE x DROP id", reported: true},
		{name: "a key of a column prefix", setup: "CREATE TABLE x (s VARCHAR(50), v INT, PRIMARY KEY (s(10)))", change: "ALTER TABLE x DROP COLUMN s", reported: true},
		{name: "a key of one column, with a unique key on a NOT NULL column left", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT NOT NULL, UNIQUE KEY uv (v))", change: "ALTER TABLE x DROP COLUMN id", reported: true},
		// (b) A key of two columns shrinks to the other, which fails when its
		// values repeat, or goes when the statement drops both.
		{name: "one column of a key of two", setup: "CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b)); INSERT INTO x VALUES (1, 1, 1), (1, 2, 2)", change: "ALTER TABLE x DROP COLUMN a", key: "b"},
		{name: "one column of a key of two, the other repeating", setup: "CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b)); INSERT INTO x VALUES (1, 1, 1), (2, 1, 2)", change: "ALTER TABLE x DROP COLUMN a", err: "Duplicate entry", key: "a,b"},
		{name: "both columns of a key of two", setup: "CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b)); INSERT INTO x VALUES (1, 1, 1), (2, 1, 2)", change: "ALTER TABLE x DROP COLUMN b, DROP COLUMN a", reported: true},
		// The second statement drops the key, but its column is one the
		// first touched: not judged.
		{name: "both columns of a key of two, in two statements", setup: "CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b))", change: "ALTER TABLE x DROP COLUMN a;\nALTER TABLE x DROP COLUMN b"},
		// (c) What depends on the column decides whether the server drops it.
		{name: "a foreign key of another table references it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE ch (x_id INT, CONSTRAINT fk_ch FOREIGN KEY (x_id) REFERENCES x (id))", change: "ALTER TABLE x DROP COLUMN id", err: "needed in a foreign key constraint", key: "id"},
		{name: "a foreign key of another table references it, with foreign_key_checks off", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE ch (x_id INT, CONSTRAINT fk_ch FOREIGN KEY (x_id) REFERENCES x (id))", change: "SET foreign_key_checks = 0;\nALTER TABLE x DROP COLUMN id", err: "needed in a foreign key constraint", key: "id"},
		{name: "a foreign key of another database references it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE DATABASE other; CREATE TABLE other.ch (x_id INT, CONSTRAINT fk_och FOREIGN KEY (x_id) REFERENCES shop.x (id))", change: "ALTER TABLE x DROP COLUMN id", err: "needed in a foreign key constraint", key: "id", reported: true, unseen: true},
		{name: "a foreign key of another table references another column", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, UNIQUE KEY uv (v)); CREATE TABLE ch (x_v INT, CONSTRAINT fk_ch FOREIGN KEY (x_v) REFERENCES x (v))", change: "ALTER TABLE x DROP COLUMN id", reported: true},
		{name: "a foreign key of the table holds it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, CONSTRAINT fk_x FOREIGN KEY (id) REFERENCES p (id))", change: "ALTER TABLE x DROP COLUMN id", err: "needed in a foreign key constraint", key: "id"},
		{name: "a foreign key of the table holds another column", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, CONSTRAINT fk_x FOREIGN KEY (v) REFERENCES p (id))", change: "ALTER TABLE x DROP COLUMN id", reported: true},
		{name: "a generated column reads it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1))", change: "ALTER TABLE x DROP COLUMN id", err: "has a generated column dependency", key: "id"},
		{name: "a generated column reads it, and goes too", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1) STORED)", change: "ALTER TABLE x DROP COLUMN id, DROP COLUMN g", reported: true},
		{name: "a check reads it alone", setup: "CREATE TABLE x (id INT PRIMARY KEY CHECK (id > 0), v INT)", change: "ALTER TABLE x DROP COLUMN id", reported: true},
		{name: "a check reads it and another column", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, CONSTRAINT ck CHECK (id < v))", change: "ALTER TABLE x DROP COLUMN id", err: "hence column cannot be dropped", key: "id"},
		{name: "a check reads both columns of the key, both dropped", setup: "CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b), CONSTRAINT ck CHECK (a < b))", change: "ALTER TABLE x DROP COLUMN a, DROP COLUMN b", err: "hence column cannot be dropped", key: "a,b"},
		{name: "a functional index reads it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, KEY fx ((id + v)))", change: "ALTER TABLE x DROP COLUMN id", err: "has a functional index dependency", key: "id"},
		{name: "another index holds it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, KEY iv (v, id))", change: "ALTER TABLE x DROP COLUMN id", reported: true},
		{name: "a view or a trigger reads it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE VIEW xv AS SELECT id FROM x; CREATE TRIGGER xt BEFORE INSERT ON x FOR EACH ROW SET NEW.v = NEW.id", change: "ALTER TABLE x DROP COLUMN id", reported: true},
		{name: "the table is partitioned by it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT) PARTITION BY HASH (id) PARTITIONS 2", change: "ALTER TABLE x DROP COLUMN id", err: "has a partitioning function dependency", key: "id"},
		{name: "the table is partitioned by KEY (), its primary key", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT) PARTITION BY KEY () PARTITIONS 2", change: "ALTER TABLE x DROP COLUMN id", err: "Field in list of fields for partition function not found", key: "id"},
		{name: "it is the only column", setup: "CREATE TABLE x (id INT PRIMARY KEY)", change: "ALTER TABLE x DROP COLUMN id", err: "You can't delete all columns", key: "id"},
		// The key keeps a column of the dropped one's name the statement
		// adds back; a subcommand that names a column of the key is not
		// judged.
		{name: "it is added back under its name", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id, ADD COLUMN id INT", key: "id"},
		{name: "a key is added back", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id, ADD PRIMARY KEY (v)", key: "v"},
		{name: "a later statement adds a key back", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id;\nALTER TABLE x ADD PRIMARY KEY (v)", key: "v"},
		// What an earlier statement, or the statement itself, drops no longer
		// keeps the column.
		{name: "a generated column reading it, dropped before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1))", change: "ALTER TABLE x DROP COLUMN g;\nALTER TABLE x DROP COLUMN id", reported: true},
		{name: "a foreign key of the table holding it, dropped before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, CONSTRAINT fk_x FOREIGN KEY (id) REFERENCES p (id))", change: "ALTER TABLE x DROP FOREIGN KEY fk_x;\nALTER TABLE x DROP COLUMN id", reported: true},
		{name: "a foreign key of the table holding it, dropped by the statement", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, CONSTRAINT fk_x FOREIGN KEY (id) REFERENCES p (id))", change: "ALTER TABLE x DROP FOREIGN KEY fk_x, DROP COLUMN id", reported: true},
		{name: "a foreign key referencing it, dropped before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE ch (x_id INT, CONSTRAINT fk_ch FOREIGN KEY (x_id) REFERENCES x (id))", change: "ALTER TABLE ch DROP FOREIGN KEY fk_ch;\nALTER TABLE x DROP COLUMN id", reported: true},
		{name: "a table holding a foreign key referencing it, dropped before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE ch (x_id INT, CONSTRAINT fk_ch FOREIGN KEY (x_id) REFERENCES x (id))", change: "DROP TABLE ch;\nALTER TABLE x DROP COLUMN id", reported: true},
		{name: "a check reading it and another column, dropped before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, CONSTRAINT ck CHECK (id < v))", change: "ALTER TABLE x DROP CHECK ck;\nALTER TABLE x DROP COLUMN id", reported: true},
		{name: "a check reading it and another column, dropped by the statement", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, CONSTRAINT ck CHECK (id < v))", change: "ALTER TABLE x DROP CONSTRAINT ck, DROP COLUMN id", reported: true},
		{name: "a functional index reading it, dropped before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, KEY fx ((id + v)))", change: "DROP INDEX fx ON x;\nALTER TABLE x DROP COLUMN id", reported: true},
		{name: "a functional index reading it, dropped by the statement", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, KEY fx ((id + v)))", change: "ALTER TABLE x DROP INDEX fx, DROP COLUMN id", reported: true},
		// A temporary table of the name, dropped, leaves the table as it was.
		{name: "a temporary table of its name dropped before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "DROP TEMPORARY TABLE IF EXISTS x;\nALTER TABLE x DROP COLUMN id", reported: true},
		// A column the statement adds keeps the table from losing its last
		// one. One an earlier statement added is not followed: not judged.
		{name: "the only column, with another added by the statement", setup: "CREATE TABLE x (id INT PRIMARY KEY)", change: "ALTER TABLE x ADD COLUMN v INT, DROP COLUMN id", reported: true},
		{name: "the only column, with another added before", setup: "CREATE TABLE x (id INT PRIMARY KEY)", change: "ALTER TABLE x ADD COLUMN v INT;\nALTER TABLE x DROP COLUMN id"},
		// A drop of a column the table does not have fails the statement.
		{name: "with a column the table does not have", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id, DROP COLUMN missing", err: "Can't DROP 'missing'", key: "id"},
		{name: "twice", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id, DROP COLUMN id", err: "Can't DROP 'id'", key: "id"},
		// Another database dropped leaves the target's tables as they are.
		{name: "another database dropped before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "DROP DATABASE IF EXISTS other;\nALTER TABLE x DROP COLUMN id", reported: true},
		{name: "another database dropped after", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id;\nDROP DATABASE IF EXISTS other", reported: true},
		// What a statement does to a table of another database, or of a name
		// that differs in case where names compare with case, leaves the
		// target's table as it was.
		{name: "a table of its name in another database, dropped before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE DATABASE other; CREATE TABLE other.x (a INT PRIMARY KEY)", change: "DROP TABLE IF EXISTS other.x;\nALTER TABLE x DROP COLUMN id", reported: true},
		{name: "a table of its name in another database, made before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE DATABASE other", change: "CREATE TABLE other.x (a INT PRIMARY KEY);\nALTER TABLE x DROP COLUMN id", reported: true},
		{name: "a table of its name in another database, renamed before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE DATABASE other; CREATE TABLE other.x (a INT PRIMARY KEY)", change: "RENAME TABLE other.x TO other.y;\nALTER TABLE x DROP COLUMN id", reported: true},
		{name: "a table of its name in another database, given a key after", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE DATABASE other; CREATE TABLE other.x (a INT)", change: "ALTER TABLE x DROP COLUMN id;\nALTER TABLE other.x ADD PRIMARY KEY (a)", reported: true},
		{name: "a table of its name in another case, made before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "CREATE TABLE X (a INT PRIMARY KEY);\nALTER TABLE x DROP COLUMN id", reported: true},
		// Another subcommand of the statement the server refuses fails it.
		{name: "with a column the table has, added", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x ADD COLUMN v INT, DROP COLUMN id", err: "Duplicate column name 'v'", key: "id"},
		{name: "with one column added twice", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x ADD COLUMN w INT, ADD COLUMN w INT, DROP COLUMN id", err: "Duplicate column name 'w'", key: "id"},
		{name: "with a column the table does not have, modified", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x MODIFY missing INT, DROP COLUMN id", err: "Unknown column 'missing'", key: "id"},
		{name: "with a column changed to a name the table has", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, w INT)", change: "ALTER TABLE x CHANGE v w INT, DROP COLUMN id", err: "Duplicate column name 'w'", key: "id"},
		{name: "with a column changed or renamed", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, u INT)", change: "ALTER TABLE x CHANGE v w INT, RENAME COLUMN u TO s, DROP COLUMN id", reported: true},
		// What a statement the server refuses drops stays: the change stops
		// there.
		{name: "a generated column reading it, dropped by a statement that fails", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1))", change: "ALTER TABLE x DROP COLUMN g, DROP COLUMN missing;\nALTER TABLE x DROP COLUMN id", err: "Can't DROP 'missing'", key: "id"},
		// A DROP TABLE under IF EXISTS passes over a table the database lacks,
		// and the server refuses to drop a table a foreign key of a table it
		// keeps references.
		{name: "a table holding a foreign key referencing it, dropped with one the database lacks under IF EXISTS", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE ch (x_id INT, CONSTRAINT fk_ch FOREIGN KEY (x_id) REFERENCES x (id))", change: "DROP TABLE IF EXISTS ch, missing;\nALTER TABLE x DROP COLUMN id", reported: true},
		{name: "a table holding a foreign key referencing it, which another table references, dropped", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE ch (k INT PRIMARY KEY, x_id INT, CONSTRAINT fk_ch FOREIGN KEY (x_id) REFERENCES x (id)); CREATE TABLE gch (ch_k INT, CONSTRAINT fk_gch FOREIGN KEY (ch_k) REFERENCES ch (k))", change: "DROP TABLE ch;\nALTER TABLE x DROP COLUMN id", err: "Cannot drop table 'ch' referenced by a foreign key", key: "id"},
		{name: "a table holding a foreign key referencing it, dropped before the table referencing it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE ch (k INT PRIMARY KEY, x_id INT, CONSTRAINT fk_ch FOREIGN KEY (x_id) REFERENCES x (id)); CREATE TABLE gch (ch_k INT, CONSTRAINT fk_gch FOREIGN KEY (ch_k) REFERENCES ch (k))", change: "DROP TABLE ch, gch;\nALTER TABLE x DROP COLUMN id", reported: true},
		{name: "a table holding a foreign key referencing it, dropped after the table referencing it", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE ch (k INT PRIMARY KEY, x_id INT, CONSTRAINT fk_ch FOREIGN KEY (x_id) REFERENCES x (id)); CREATE TABLE gch (ch_k INT, CONSTRAINT fk_gch FOREIGN KEY (ch_k) REFERENCES ch (k))", change: "DROP TABLE gch, ch;\nALTER TABLE x DROP COLUMN id", reported: true},
		// The server refuses an index or a foreign key the statement adds
		// under a name in use, on a column the table lacks once it ran, or
		// with no index its referenced columns lead, and the rename of an
		// index the table lacks.
		{name: "with an index added under a name the table uses", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, KEY iv (v))", change: "ALTER TABLE x DROP COLUMN id, ADD INDEX iv (v)", err: "Duplicate key name 'iv'", key: "id"},
		{name: "with one index name added twice", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x ADD INDEX i1 (v), ADD INDEX i1 (v), DROP COLUMN id", err: "Duplicate key name 'i1'", key: "id"},
		{name: "with an index added on the column it drops", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x ADD INDEX i1 (id), DROP COLUMN id", err: "Key column 'id' doesn't exist", key: "id"},
		{name: "with a foreign key added under a name the database uses", setup: "CREATE TABLE par (id INT PRIMARY KEY); CREATE TABLE t2 (a INT, CONSTRAINT fk_dup FOREIGN KEY (a) REFERENCES par (id)); CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x ADD CONSTRAINT fk_dup FOREIGN KEY (v) REFERENCES par (id), DROP COLUMN id", err: "Duplicate foreign key constraint name", key: "id"},
		{name: "with a foreign key added whose referenced columns lead no index", setup: "CREATE TABLE par (id INT PRIMARY KEY, w INT); CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x ADD CONSTRAINT fk_w FOREIGN KEY (v) REFERENCES par (w), DROP COLUMN id", err: "Missing index for constraint", key: "id"},
		{name: "with an index the table lacks renamed", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x RENAME INDEX missing TO iw, DROP COLUMN id", err: "Key 'missing' doesn't exist", key: "id"},
		{name: "with an index dropped and added back, one renamed, and a unique key and a foreign key added", setup: "CREATE TABLE par (id INT PRIMARY KEY); CREATE TABLE x (id INT PRIMARY KEY, v INT, w INT, KEY iv (v), KEY iw (w))", change: "ALTER TABLE x DROP INDEX iv, ADD INDEX iv (v), RENAME INDEX iw TO iz, ADD UNIQUE KEY uv (v), ADD CONSTRAINT fk_ok FOREIGN KEY (v) REFERENCES par (id), DROP COLUMN id", reported: true},
		// A stored generated column the change makes an ordinary one reads
		// nothing; the server refuses to make a virtual one ordinary.
		{name: "a stored generated column reading it, made ordinary by the statement", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1) STORED)", change: "ALTER TABLE x MODIFY COLUMN g INT, DROP COLUMN id", reported: true},
		{name: "a stored generated column reading it, made ordinary before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1) STORED)", change: "ALTER TABLE x MODIFY COLUMN g INT;\nALTER TABLE x DROP COLUMN id", reported: true},
		{name: "a stored generated column reading it, changed to an ordinary one", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1) STORED)", change: "ALTER TABLE x CHANGE COLUMN g g2 INT, DROP COLUMN id", reported: true},
		{name: "a virtual generated column reading it, made ordinary by the statement", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1) VIRTUAL)", change: "ALTER TABLE x MODIFY COLUMN g INT, DROP COLUMN id", err: "'Changing the STORED status' is not supported", key: "id"},
		// A foreign key of the table to itself references an index the
		// statement keeps whole.
		{name: "with a foreign key to itself added, to a unique key it keeps", setup: "CREATE TABLE x (id INT PRIMARY KEY, p_id INT, code INT, UNIQUE KEY uk (code))", change: "ALTER TABLE x DROP COLUMN id, ADD CONSTRAINT fk_self FOREIGN KEY (p_id) REFERENCES x (code)", reported: true},
		{name: "with a foreign key to itself added, to a column with no index", setup: "CREATE TABLE x (id INT PRIMARY KEY, p_id INT, code INT)", change: "ALTER TABLE x DROP COLUMN id, ADD CONSTRAINT fk_self FOREIGN KEY (p_id) REFERENCES x (code)", err: "Missing index for constraint", key: "id"},
		{name: "with a foreign key to itself added, to a key the statement adds", setup: "CREATE TABLE x (id INT PRIMARY KEY, p_id INT, code INT)", change: "ALTER TABLE x ADD UNIQUE KEY uk (code), ADD CONSTRAINT fk_self FOREIGN KEY (p_id) REFERENCES x (code), DROP COLUMN id", err: "Cannot add foreign key constraint", key: "id"},
		{name: "with a foreign key to itself added, to an index the statement drops", setup: "CREATE TABLE x (id INT PRIMARY KEY, p_id INT, code INT, UNIQUE KEY uk (code))", change: "ALTER TABLE x DROP INDEX uk, DROP COLUMN id, ADD CONSTRAINT fk_self FOREIGN KEY (p_id) REFERENCES x (code)", err: "Missing index for constraint", key: "id"},
		{name: "with a foreign key to itself added, to an index the statement takes a column after the referenced ones from", setup: "CREATE TABLE x (id INT PRIMARY KEY, p_id INT, code INT, KEY ic (code, id))", change: "ALTER TABLE x DROP COLUMN id, ADD CONSTRAINT fk_self FOREIGN KEY (p_id) REFERENCES x (code)", reported: true},
		{name: "with a foreign key to itself added, to an index the statement takes a leading column from", setup: "CREATE TABLE x (id INT PRIMARY KEY, p_id INT, code INT, KEY ic (id, code))", change: "ALTER TABLE x DROP COLUMN id, ADD CONSTRAINT fk_self FOREIGN KEY (p_id) REFERENCES x (code)", err: "Cannot add foreign key constraint", key: "id"},
		// A new generation expression replaces what a column read; the
		// server keeps a column stored or not, and refuses an expression of
		// the column it drops.
		{name: "a stored generated column reading it, given one reading another column", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1) STORED)", change: "ALTER TABLE x MODIFY COLUMN g INT AS (v + 1) STORED, DROP COLUMN id", reported: true},
		{name: "a virtual generated column reading it, given one reading another column", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1) VIRTUAL)", change: "ALTER TABLE x MODIFY COLUMN g INT AS (v + 1) VIRTUAL, DROP COLUMN id", reported: true},
		{name: "a stored generated column reading it, changed to one of a new name reading another column", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1) STORED)", change: "ALTER TABLE x CHANGE COLUMN g g2 INT AS (v + 1) STORED, DROP COLUMN id", reported: true},
		{name: "an ordinary column made a stored generated one", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, c INT)", change: "ALTER TABLE x MODIFY COLUMN c INT AS (v + 1) STORED, DROP COLUMN id", reported: true},
		{name: "a generated column reading another column, added", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x ADD COLUMN g INT AS (v + 1) VIRTUAL, DROP COLUMN id", reported: true},
		{name: "a stored generated column made virtual", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1) STORED)", change: "ALTER TABLE x MODIFY COLUMN g INT AS (v + 1) VIRTUAL, DROP COLUMN id", err: "'Changing the STORED status' is not supported", key: "id"},
		{name: "an ordinary column made a virtual generated one", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, c INT)", change: "ALTER TABLE x MODIFY COLUMN c INT AS (v + 1) VIRTUAL, DROP COLUMN id", err: "'Changing the STORED status' is not supported", key: "id"},
		{name: "a stored generated column given one reading the column it drops", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1) STORED)", change: "ALTER TABLE x MODIFY COLUMN g INT AS (id + 2) STORED, DROP COLUMN id", err: "Unknown column 'id' in 'generated column function'", key: "id"},
		{name: "a generated column reading it, added", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x ADD COLUMN g INT AS (id + 1) STORED, DROP COLUMN id", err: "Unknown column 'id' in 'generated column function'", key: "id"},
		// A database holds a foreign key name once; another database may
		// hold it too.
		{name: "with a foreign key added under a name a statement before gave one in another database", setup: "CREATE TABLE x (id INT PRIMARY KEY, p_id INT); CREATE TABLE par (id INT PRIMARY KEY); CREATE DATABASE other; CREATE TABLE other.par (id INT PRIMARY KEY)", change: "CREATE TABLE other.c (a INT PRIMARY KEY, CONSTRAINT fk_new FOREIGN KEY (a) REFERENCES other.par (id));\nALTER TABLE x DROP COLUMN id, ADD CONSTRAINT fk_new FOREIGN KEY (p_id) REFERENCES par (id)", reported: true},
		// A column's REFERENCES, which MySQL parses but ignores, makes no
		// foreign key.
		{name: "a column's REFERENCES to it, made before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "CREATE TABLE c (k INT REFERENCES x (id), PRIMARY KEY (k));\nALTER TABLE x DROP COLUMN id", reported: true},
		{name: "a column's REFERENCES to it, added before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE c (k INT PRIMARY KEY)", change: "ALTER TABLE c ADD COLUMN r INT REFERENCES x (id);\nALTER TABLE x DROP COLUMN id", reported: true},
		// A check the change adds reads the columns it names; one of the
		// column alone, which the server takes with it, is not judged.
		{name: "a check of another column, added before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); INSERT INTO x VALUES (1, 1)", change: "ALTER TABLE x ADD CHECK (v > 0);\nALTER TABLE x DROP COLUMN id", reported: true},
		{name: "a check of another column, added by the statement", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); INSERT INTO x VALUES (1, 1)", change: "ALTER TABLE x ADD CHECK (v > 0), DROP COLUMN id", reported: true},
		{name: "a check of a column added with it, before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x ADD COLUMN w INT CHECK (w > 0);\nALTER TABLE x DROP COLUMN id", reported: true},
		{name: "a check of it and another column, added before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); INSERT INTO x VALUES (1, 2)", change: "ALTER TABLE x ADD CHECK (id < v);\nALTER TABLE x DROP COLUMN id", err: "hence column cannot be dropped", key: "id"},
		{name: "a check of it alone, added before", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT); INSERT INTO x VALUES (1, 2)", change: "ALTER TABLE x ADD CHECK (id > 0);\nALTER TABLE x DROP COLUMN id"},
		{name: "with a check added of a column the table lacks", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x ADD CHECK (missing > 0), DROP COLUMN id", err: "Unknown column 'missing'", key: "id"},
		{name: "with a check added under a name another table uses", setup: "CREATE TABLE par (id INT PRIMARY KEY); CREATE TABLE x (id INT PRIMARY KEY, v INT, w INT, CONSTRAINT ck1 CHECK (w > 0), CONSTRAINT fk1 FOREIGN KEY (w) REFERENCES par (id)); CREATE TABLE y (a INT, CONSTRAINT cky CHECK (a > 0))", change: "ALTER TABLE x ADD CONSTRAINT cky CHECK (v > 0), DROP COLUMN id", err: "Duplicate check constraint name 'cky'", key: "id"},
		{name: "with a check added under a name a foreign key of the table uses", setup: "CREATE TABLE par (id INT PRIMARY KEY); CREATE TABLE x (id INT PRIMARY KEY, v INT, w INT, CONSTRAINT ck1 CHECK (w > 0), CONSTRAINT fk1 FOREIGN KEY (w) REFERENCES par (id)); CREATE TABLE y (a INT, CONSTRAINT cky CHECK (a > 0))", change: "ALTER TABLE x ADD CONSTRAINT fk1 CHECK (v > 0), DROP COLUMN id", reported: true},
		{name: "with one check name added twice", setup: "CREATE TABLE par (id INT PRIMARY KEY); CREATE TABLE x (id INT PRIMARY KEY, v INT, w INT, CONSTRAINT ck1 CHECK (w > 0), CONSTRAINT fk1 FOREIGN KEY (w) REFERENCES par (id)); CREATE TABLE y (a INT, CONSTRAINT cky CHECK (a > 0))", change: "ALTER TABLE x ADD CONSTRAINT c2 CHECK (v > 0), ADD CONSTRAINT c2 CHECK (v < 9), DROP COLUMN id", err: "Duplicate check constraint name 'c2'", key: "id"},
		// A later statement that adds a key back but fails stops the change
		// before its end, which the rule, as pg/review does, does not judge.
		{name: "a later statement that adds a key back, but fails", setup: "CREATE TABLE x (id INT PRIMARY KEY, v INT)", change: "ALTER TABLE x DROP COLUMN id;\nALTER TABLE x DROP COLUMN missing, ADD PRIMARY KEY (v)", err: "Can't DROP 'missing'"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := resetOracle(t, db)
			ctx := context.Background()
			if _, err := conn.ExecContext(ctx, tt.setup); err != nil {
				t.Fatalf("setup: %v", err)
			}
			schema := syncSchema(t, conn)
			result, err := Review(ctx, tt.change, review.Options{Rules: []review.Rule{review.RequirePrimaryKey}},
				[]review.Target{{Schema: schema, LowerCaseTableNames: lowerCaseTableNames(t, conn)}})
			if err != nil {
				t.Fatalf("Review: %v", err)
			}
			reported := false
			for _, f := range result.Findings {
				if f.Rule == review.Syntax {
					t.Fatalf("syntax: %s", f.Message)
				}
				reported = true
			}
			var failed error
			for _, s := range parseChange(t, tt.change) {
				if _, failed = conn.ExecContext(ctx, tt.change[s.loc.Start:s.loc.End]); failed != nil {
					break
				}
			}
			switch {
			case tt.err == "" && failed != nil:
				t.Errorf("the change fails: %v", failed)
			case tt.err != "" && (failed == nil || !strings.Contains(failed.Error(), tt.err)):
				t.Errorf("the change ends with %v, want an error with %q", failed, tt.err)
			}
			if key := primaryKeyOf(t, conn, "x"); key != tt.key {
				t.Errorf("key of x after the change = %q, want %q", key, tt.key)
			}
			if reported != tt.reported {
				t.Errorf("reported = %v, want %v", reported, tt.reported)
			}
			if reported && (failed != nil || tt.key != "") && !tt.unseen {
				t.Errorf("reported, but x keeps its key")
			}
		})
	}
}

// primaryKeyOf reads the columns of a shop table's primary key, comma
// separated, "" when it has none.
func primaryKeyOf(t *testing.T, conn *sql.Conn, table string) string {
	t.Helper()
	var key sql.NullString
	if err := conn.QueryRowContext(context.Background(), `SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'`, table).Scan(&key); err != nil {
		t.Fatal(err)
	}
	return key.String
}
