package review

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/testcontainers/testcontainers-go"
	tcmssql "github.com/testcontainers/testcontainers-go/modules/mssql"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/review"
)

// oracleBase is the schema every oracle case starts from, in database shop.
var oracleBase = []string{
	"CREATE TABLE dbo.p (id INT NOT NULL CONSTRAINT pk_p PRIMARY KEY, note VARCHAR(20))",
	`CREATE TABLE dbo.t (id INT NOT NULL CONSTRAINT pk_t PRIMARY KEY, code VARCHAR(20) CONSTRAINT uk_code UNIQUE, p_id INT CONSTRAINT fk_p REFERENCES dbo.p (id), n INT CONSTRAINT ck_n CHECK (n > 0))`,
	"CREATE INDEX ix_n ON dbo.t (n)",
	"CREATE TABLE dbo.nokey (a INT)",
	"CREATE TABLE dbo.plain (id INT NOT NULL CONSTRAINT pk_plain PRIMARY KEY, v INT)",
	"INSERT INTO dbo.p VALUES (1, 'a'), (2, 'b')",
	"INSERT INTO dbo.t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2)",
	"INSERT INTO dbo.nokey VALUES (1), (2)",
	"INSERT INTO dbo.plain VALUES (1, 1), (2, 2)",
}

var mssqlOracle struct {
	once      sync.Once
	container *tcmssql.MSSQLServerContainer
	db        *sql.DB
	err       error
}

func TestMain(m *testing.M) {
	code := m.Run()
	if mssqlOracle.container != nil {
		_ = mssqlOracle.db.Close()
		_ = testcontainers.TerminateContainer(mssqlOracle.container)
	}
	os.Exit(code)
}

// startMSSQL returns the SQL Server the oracle tests of this package ask,
// started once.
func startMSSQL(t *testing.T) *sql.DB {
	t.Helper()
	mssqlOracle.once.Do(func() {
		ctx := context.Background()
		c, err := tcmssql.Run(ctx, "mcr.microsoft.com/mssql/server:2022-latest", tcmssql.WithAcceptEULA(), tcmssql.WithPassword("Str0ngPa$$w0rd!"))
		if err != nil {
			mssqlOracle.err = err
			return
		}
		mssqlOracle.container = c
		dsn, err := c.ConnectionString(ctx)
		if err != nil {
			mssqlOracle.err = err
			return
		}
		mssqlOracle.db, mssqlOracle.err = sql.Open("sqlserver", dsn)
		if mssqlOracle.err == nil {
			mssqlOracle.err = mssqlOracle.db.PingContext(ctx)
		}
	})
	if mssqlOracle.err != nil {
		t.Fatalf("SQL Server oracle unavailable: %v", mssqlOracle.err)
	}
	return mssqlOracle.db
}

// resetOracle recreates the shop database from oracleBase, and an empty
// bbdataarchive, and returns a connection with shop selected.
func resetOracle(t *testing.T, db *sql.DB) *sql.Conn {
	t.Helper()
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	stmts := []string{"USE master"}
	for _, name := range []string{"shop", "bbdataarchive"} {
		stmts = append(stmts, fmt.Sprintf("IF DB_ID('%s') IS NOT NULL BEGIN ALTER DATABASE %s SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE %s; END", name, name, name), "CREATE DATABASE "+name)
	}
	stmts = append(stmts, "USE shop")
	stmts = append(stmts, oracleBase...)
	for _, stmt := range stmts {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	return conn
}

// syncSchema reads the shop database into metadata, the parts the rules
// read: tables with their columns, indexes, foreign keys, and checks, and
// views, by schema.
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
	schemas := make(map[string]*metadata.SchemaMetadata)
	var order []string
	schemaOf := func(name string) *metadata.SchemaMetadata {
		s, ok := schemas[name]
		if !ok {
			s = &metadata.SchemaMetadata{Name: name}
			schemas[name] = s
			order = append(order, name)
		}
		return s
	}
	tables := make(map[string]*metadata.TableMetadata)
	query(`SELECT s.name, o.name, o.type FROM sys.objects o JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type IN ('U', 'V') AND o.is_ms_shipped = 0 ORDER BY s.name, o.name`, func(r *sql.Rows) error {
		var schema, name, kind string
		if err := r.Scan(&schema, &name, &kind); err != nil {
			return err
		}
		s := schemaOf(schema)
		if strings.TrimSpace(kind) == "V" {
			s.Views = append(s.Views, &metadata.ViewMetadata{Name: name})
			return nil
		}
		tbl := &metadata.TableMetadata{Name: name}
		s.Tables = append(s.Tables, tbl)
		tables[schema+"."+name] = tbl
		return nil
	})
	query(`SELECT s.name, o.name, c.name FROM sys.columns c JOIN sys.objects o ON o.object_id = c.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.column_id`, func(r *sql.Rows) error {
		var schema, tbl, col string
		if err := r.Scan(&schema, &tbl, &col); err != nil {
			return err
		}
		if t, ok := tables[schema+"."+tbl]; ok {
			t.Columns = append(t.Columns, &metadata.ColumnMetadata{Name: col})
		}
		return nil
	})
	query(`SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint FROM sys.indexes i JOIN sys.objects o ON o.object_id = i.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' AND i.name IS NOT NULL`, func(r *sql.Rows) error {
		var schema, tbl, name string
		var primary, unique, constraint bool
		if err := r.Scan(&schema, &tbl, &name, &primary, &unique, &constraint); err != nil {
			return err
		}
		if t, ok := tables[schema+"."+tbl]; ok {
			t.Indexes = append(t.Indexes, &metadata.IndexMetadata{Name: name, Primary: primary, Unique: unique, IsConstraint: constraint || primary})
		}
		return nil
	})
	query(`SELECT s.name, o.name, f.name FROM sys.foreign_keys f JOIN sys.objects o ON o.object_id = f.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id`, func(r *sql.Rows) error {
		var schema, tbl, name string
		if err := r.Scan(&schema, &tbl, &name); err != nil {
			return err
		}
		if t, ok := tables[schema+"."+tbl]; ok {
			t.ForeignKeys = append(t.ForeignKeys, &metadata.ForeignKeyMetadata{Name: name})
		}
		return nil
	})
	query(`SELECT s.name, o.name, c.name FROM sys.check_constraints c JOIN sys.objects o ON o.object_id = c.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id`, func(r *sql.Rows) error {
		var schema, tbl, name string
		if err := r.Scan(&schema, &tbl, &name); err != nil {
			return err
		}
		if t, ok := tables[schema+"."+tbl]; ok {
			t.CheckConstraints = append(t.CheckConstraints, &metadata.CheckConstraintMetadata{Name: name})
		}
		return nil
	})
	db := &metadata.DatabaseSchemaMetadata{Name: "shop"}
	for _, name := range order {
		db.Schemas = append(db.Schemas, schemas[name])
	}
	return db
}

// parseChange parses a change the way Review does.
func parseChange(t *testing.T, change string) []statement {
	t.Helper()
	_, stmts, syntax := parseReview(change)
	if syntax != nil {
		t.Fatalf("syntax: %s", syntax.Message)
	}
	return stmts
}

// TestPriorBackupAgainstMSSQL builds each statement's backup query as
// Bytebase does and runs it on SQL Server before the change. A PriorBackup
// finding on a statement must mean that query fails, unless Bytebase
// refuses the statement before building it; for a complete case, every
// query that fails must have a finding.
func TestPriorBackupAgainstMSSQL(t *testing.T) {
	db := startMSSQL(t)
	tests := []struct {
		name     string
		change   string
		complete bool
	}{
		{"a plain backup", "UPDATE t SET n = 5 WHERE id = 1;\nDELETE FROM nokey WHERE a = 1;", true},
		{"TOP", "UPDATE TOP (1) p SET note = '' WHERE id > 0;", true},
		{"a table the change creates", "CREATE TABLE n2 (id INT PRIMARY KEY);\nUPDATE n2 SET id = 2 WHERE id = 1;", true},
		{"a table the change creates, joined", "CREATE TABLE n2 (id INT PRIMARY KEY);\nDELETE x FROM t AS x JOIN n2 ON x.id = n2.id WHERE x.n = 1;", true},
		{"a table the change creates, in a subquery", "CREATE TABLE n2 (id INT PRIMARY KEY);\nDELETE FROM t WHERE id IN (SELECT id FROM n2);", true},
		{"a column the change adds, unqualified", "ALTER TABLE t ADD gone INT;\nUPDATE t SET n = 1 WHERE gone = 1;", true},
		{"a column the change adds, qualified", "ALTER TABLE nokey ADD gone INT;\nDELETE FROM nokey WHERE nokey.gone = 1;", true},
		{"a column of a joined table", "ALTER TABLE p ADD memo VARCHAR(10);\nUPDATE x SET x.n = 1 FROM t AS x JOIN p ON p.memo = x.code WHERE x.id = 3;", true},
		{"a column of a table the change recreates", "DROP TABLE nokey;\nCREATE TABLE nokey (b INT);\nDELETE FROM nokey WHERE b = 1;", true},
		{"a column in a subquery, not judged", "ALTER TABLE t ADD gone INT;\nUPDATE plain SET v = 1 WHERE id IN (SELECT gone FROM t);", false},
		{"a column the SET clause names, not copied", "ALTER TABLE t ADD gone INT;\nUPDATE t SET gone = 1 WHERE id = 1;", true},
		{"a column sp_rename renames to", "EXEC sp_rename 'dbo.nokey.a', 'b', 'COLUMN';\nDELETE FROM dbo.nokey WHERE b = 1;", true},
		{"a FROM list without the target", "UPDATE dbo.t SET n = p.id FROM dbo.p AS p;", true},
		{"a table a CREATE SCHEMA makes", "CREATE SCHEMA s3 CREATE TABLE s3.n6 (id INT);\nUPDATE s3.n6 SET id = 2 WHERE id = 1;", true},
		{"a synonym the change creates", "CREATE SYNONYM n5 FOR dbo.t;\nUPDATE n5 SET n = 2 WHERE id = 1;", true},
		{"a table SELECT INTO makes", "SELECT 1 AS id INTO n4;\nUPDATE n4 SET id = 2 WHERE id = 1;", true},
		{"statements before a USE", "ALTER TABLE nokey ADD gone INT;\nDELETE FROM nokey WHERE gone = 1;\nUSE master;\nDELETE FROM shop.dbo.nokey WHERE a = 2;", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := resetOracle(t, db)
			schema := syncSchema(t, conn)
			result, err := Review(context.Background(), tt.change, review.Options{
				Rules:  []review.Rule{review.PriorBackup},
				Change: review.Change{PriorBackup: true, MaxBackupSize: 1 << 20},
			}, []review.Target{{Schema: schema, BackupDatabaseExists: true}})
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
// bytebase: backend/plugin/parser/tsql/backup.go — generateSQLForTable
func backupFailures(t *testing.T, conn *sql.Conn, change string) map[int]error {
	t.Helper()
	ctx := context.Background()
	fails := make(map[int]error)
	for n, s := range parseChange(t, change) {
		d, ok := dmlOf(s.node)
		if !ok {
			continue
		}
		ref, isRef := d.target.(*ast.TableRef)
		if !isRef {
			continue
		}
		table := backedUp(ref, d.from)
		sel := "[shop].[dbo].[" + table.Object + "].*"
		if table.Alias != "" {
			sel = "[" + table.Alias + "].*"
		}
		if d.top != nil {
			sel += " " + textOf(change, d.top.Loc)
		}
		var from string
		if d.from != nil && len(d.from.Items) > 0 {
			first, last := d.from.Items[0], d.from.Items[len(d.from.Items)-1]
			from = change[nodeLoc(first).Start:nodeLoc(last).End]
		} else {
			from = textOf(change, ref.Loc)
		}
		query := "SELECT " + sel + " FROM " + from
		if d.where != nil {
			query += " WHERE " + change[nodeLoc(d.where).Start:s.loc.End]
		}
		backup := fmt.Sprintf("SELECT * INTO [bbdataarchive].[dbo].[bk_%d] FROM (%s) AS backup_table", n, query)
		if d.with != nil {
			backup = textOf(change, d.with.Loc) + "\n" + backup
		}
		if _, err := conn.ExecContext(ctx, backup); err != nil {
			fails[s.index] = fmt.Errorf("%s: %w", backup, err)
		}
	}
	return fails
}

// TestScanRulesAgainstMSSQL runs each change on SQL Server, statement by
// statement, and derives from the catalog what the constraint rules
// should report: DisallowDropConstraint, a primary key, unique, foreign
// key, or check constraint a statement removes; RequirePrimaryKey, a table
// a CREATE TABLE made, or one that had a primary key, that ends the change
// without one. Every finding must match; for a complete case, every
// expected finding must be reported.
func TestScanRulesAgainstMSSQL(t *testing.T) {
	db := startMSSQL(t)
	tests := []struct {
		name     string
		change   string
		complete bool
	}{
		{"constraints of the synced schema", "ALTER TABLE t DROP CONSTRAINT fk_p, ck_n;\nALTER TABLE t DROP CONSTRAINT uk_code;", true},
		{"the primary key", "ALTER TABLE plain DROP CONSTRAINT pk_plain;", true},
		{"a constraint the change defines, then drops", "ALTER TABLE nokey ADD CONSTRAINT ck_a CHECK (a > 0);\nALTER TABLE nokey DROP CONSTRAINT ck_a;", true},
		{"a default constraint", "ALTER TABLE nokey ADD CONSTRAINT df_a DEFAULT 0 FOR a;\nALTER TABLE nokey DROP CONSTRAINT df_a;", true},
		{"a table created without a key", "CREATE TABLE a (id INT);\nCREATE TABLE b (id INT PRIMARY KEY);", true},
		{"a key added later", "CREATE TABLE a (id INT NOT NULL);\nALTER TABLE a ADD CONSTRAINT pk_a PRIMARY KEY (id);", true},
		{"a primary key dropped, then added back", "ALTER TABLE plain DROP CONSTRAINT pk_plain;\nALTER TABLE plain ADD CONSTRAINT pk_plain2 PRIMARY KEY (id);", true},
		{"a primary key dropped in a transaction rolled back", "BEGIN TRAN;\nALTER TABLE plain DROP CONSTRAINT pk_plain;\nROLLBACK;", true},
		{"a drop inside IF", "IF 1 = 1 ALTER TABLE t DROP CONSTRAINT ck_n;", true},
		{"a table whose database the change drops", "CREATE DATABASE d;\nCREATE TABLE d.dbo.e (id INT);\nDROP DATABASE d;", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := resetOracle(t, db)
			schema := syncSchema(t, conn)
			result, err := Review(context.Background(), tt.change, review.Options{
				Rules: []review.Rule{review.DisallowDropConstraint, review.RequirePrimaryKey},
			}, []review.Target{{Schema: schema, BackupDatabaseExists: true}})
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

// constraintSet reads every primary key, unique, foreign key, and check
// constraint of the shop database, as table/kind/name/definition.
func constraintSet(t *testing.T, conn *sql.Conn) map[string]bool {
	t.Helper()
	ctx := context.Background()
	out := make(map[string]bool)
	rows, err := conn.QueryContext(ctx, `
		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tbl, kind, name, def string
		if err := rows.Scan(&tbl, &kind, &name, &def); err != nil {
			t.Fatal(err)
		}
		out[strings.ToLower(tbl+"/"+strings.TrimSpace(kind)+"/"+name+"/"+def)] = true
	}
	return out
}

// primaryKeyed reads the user tables of the shop database that have a
// primary key, and all of them.
func primaryKeyed(t *testing.T, conn *sql.Conn) (keyed, all map[string]bool) {
	t.Helper()
	ctx := context.Background()
	keyed, all = make(map[string]bool), make(map[string]bool)
	rows, err := conn.QueryContext(ctx, `SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var has int
		if err := rows.Scan(&name, &has); err != nil {
			t.Fatal(err)
		}
		all[strings.ToLower(name)] = true
		if has == 1 {
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
		if v, ok := s.node.(*ast.CreateTableStmt); ok && judgedTable(v.Name) {
			createdAt[strings.ToLower(v.Name.Object)] = s.index
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

// nodeLoc reads a node's location.
func nodeLoc(n ast.Node) ast.Loc { return ast.NodeLoc(n) }
