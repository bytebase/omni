package review

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	_ "github.com/sijms/go-ora/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/review"
)

// oracleBase is the schema every oracle case starts from, in schema SHOP,
// the database the change runs in.
var oracleBase = []string{
	"CREATE TABLE p (id NUMBER CONSTRAINT pk_p PRIMARY KEY, note VARCHAR2(20))",
	"CREATE TABLE t (id NUMBER CONSTRAINT pk_t PRIMARY KEY, code VARCHAR2(20) CONSTRAINT nn_code NOT NULL CONSTRAINT uk_code UNIQUE, p_id NUMBER CONSTRAINT fk_p REFERENCES p (id), n NUMBER CONSTRAINT ck_n CHECK (n > 0))",
	"CREATE INDEX ix_n ON t (n)",
	"CREATE TABLE nokey (a NUMBER)",
	"CREATE TABLE plain (id NUMBER CONSTRAINT pk_plain PRIMARY KEY, v NUMBER)",
	"CREATE UNIQUE INDEX ux_v ON plain (v)",
	"CREATE TABLE longt (id NUMBER CONSTRAINT pk_longt PRIMARY KEY, body LONG)",
	`CREATE TABLE "a'b" (id NUMBER PRIMARY KEY)`,
	"CREATE VIEW pv AS SELECT id, note FROM p",
	"INSERT INTO p VALUES (1, 'a')",
	"INSERT INTO p VALUES (2, 'b')",
	"INSERT INTO t VALUES (1, 'x', 1, 1)",
	"INSERT INTO t VALUES (2, 'y', 2, 2)",
	"INSERT INTO nokey VALUES (1)",
	"INSERT INTO plain VALUES (1, 1)",
	"INSERT INTO longt VALUES (1, 'x')",
	`INSERT INTO "a'b" VALUES (1)`,
}

const oraclePassword = "testpass"

var oracleOracle struct {
	once      sync.Once
	container testcontainers.Container
	admin     *sql.DB
	host      string
	port      string
	err       error
}

func TestMain(m *testing.M) {
	code := m.Run()
	if oracleOracle.container != nil {
		if oracleOracle.admin != nil {
			_ = oracleOracle.admin.Close()
		}
		_ = testcontainers.TerminateContainer(oracleOracle.container)
	}
	os.Exit(code)
}

// startOracle returns a SYSTEM connection to the Oracle the oracle tests
// of this package ask, started once, the image oracle/parser's tests pin.
// The users SHOP and BBDATAARCHIVE are made per case by resetOracle; the
// public synonym PUBT, for SYSTEM.PUBT, stays.
func startOracle(t *testing.T) *sql.DB {
	t.Helper()
	o := &oracleOracle
	o.once.Do(func() {
		ctx := context.Background()
		c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        "gvenzl/oracle-free:23.26.3-slim-faststart",
				ExposedPorts: []string{"1521/tcp"},
				Env:          map[string]string{"ORACLE_PASSWORD": oraclePassword},
				WaitingFor:   wait.ForLog("DATABASE IS READY TO USE!").WithStartupTimeout(5 * time.Minute),
			},
			Started: true,
		})
		if err != nil {
			o.err = err
			return
		}
		o.container = c
		if o.host, err = c.Host(ctx); err != nil {
			o.err = err
			return
		}
		port, err := c.MappedPort(ctx, "1521")
		if err != nil {
			o.err = err
			return
		}
		o.port = port.Port()
		if o.admin, o.err = openOracle("system"); o.err != nil {
			return
		}
		for _, stmt := range []string{
			"CREATE TABLE system.pubt (id NUMBER, a NUMBER)",
			"GRANT SELECT ON system.pubt TO PUBLIC",
			"CREATE PUBLIC SYNONYM pubt FOR system.pubt",
		} {
			if _, err := o.admin.ExecContext(ctx, stmt); err != nil {
				o.err = fmt.Errorf("%s: %w", stmt, err)
				return
			}
		}
	})
	if o.err != nil {
		t.Fatalf("Oracle oracle unavailable: %v", o.err)
	}
	return o.admin
}

func openOracle(user string) (*sql.DB, error) {
	db, err := sql.Open("oracle", fmt.Sprintf("oracle://%s:%s@%s:%s/FREEPDB1", user, oraclePassword, oracleOracle.host, oracleOracle.port))
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(context.Background()); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// resetOracle recreates the users SHOP, from oracleBase, and, unless
// noArchive, an empty BBDATAARCHIVE, and returns a connection as SHOP:
// the session's current schema is SHOP, as Bytebase sets it for the task.
func resetOracle(t *testing.T, admin *sql.DB, noArchive bool) *sql.Conn {
	t.Helper()
	ctx := context.Background()
	stmts := []string{
		`BEGIN
  FOR u IN (SELECT username FROM dba_users WHERE username IN ('SHOP', 'BBDATAARCHIVE')) LOOP
    EXECUTE IMMEDIATE 'DROP USER ' || u.username || ' CASCADE';
  END LOOP;
END;`,
		"CREATE USER shop IDENTIFIED BY " + oraclePassword + " QUOTA UNLIMITED ON users",
		"GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE MATERIALIZED VIEW, CREATE SEQUENCE, CREATE SYNONYM, CREATE PROCEDURE, CREATE TRIGGER, CREATE ANY TABLE, COMMENT ANY TABLE TO shop",
	}
	if !noArchive {
		stmts = append(stmts, "CREATE USER bbdataarchive IDENTIFIED BY "+oraclePassword+" QUOTA UNLIMITED ON users")
	}
	for _, stmt := range stmts {
		if _, err := admin.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	db, err := openOracle("shop")
	if err != nil {
		t.Fatal(err)
	}
	conn, err := db.Conn(ctx)
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() {
		conn.Close()
		db.Close()
	})
	for _, stmt := range oracleBase {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	return conn
}

// syncSchema reads schema SHOP into metadata as Bytebase syncs an Oracle
// database, the parts the rules read: into one schema named "", the tables
// that are not materialized views, with their visible columns and types;
// the primary keys and unique constraints, as indexes; the foreign keys;
// the check constraints whose names are not system-generated; the other
// indexes; and the views.
//
// bytebase: backend/plugin/db/oracle/sync.go — SyncDBSchema, getTables, getTableColumns, getConstraints, getIndexesAndConstraints
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
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
	}
	s := &metadata.SchemaMetadata{}
	tables := make(map[string]*metadata.TableMetadata)
	query(`SELECT table_name FROM all_tables WHERE owner = 'SHOP' AND table_name NOT IN (SELECT mview_name FROM all_mviews WHERE owner = 'SHOP') ORDER BY table_name`, func(r *sql.Rows) error {
		var name string
		if err := r.Scan(&name); err != nil {
			return err
		}
		tables[name] = &metadata.TableMetadata{Name: name}
		s.Tables = append(s.Tables, tables[name])
		return nil
	})
	query(`SELECT table_name, column_name, data_type FROM all_tab_columns WHERE owner = 'SHOP' AND column_id IS NOT NULL ORDER BY table_name, column_id`, func(r *sql.Rows) error {
		var tbl, col, typ string
		if err := r.Scan(&tbl, &col, &typ); err != nil {
			return err
		}
		if t, ok := tables[tbl]; ok {
			t.Columns = append(t.Columns, &metadata.ColumnMetadata{Name: col, Type: typ})
		}
		return nil
	})
	constraints := make(map[string]bool)
	query(`SELECT table_name, constraint_name, constraint_type FROM all_constraints WHERE owner = 'SHOP' ORDER BY table_name, constraint_name`, func(r *sql.Rows) error {
		var tbl, name, kind string
		if err := r.Scan(&tbl, &name, &kind); err != nil {
			return err
		}
		t, ok := tables[tbl]
		if !ok {
			return nil
		}
		switch kind {
		case "P", "U":
			t.Indexes = append(t.Indexes, &metadata.IndexMetadata{Name: name, Primary: kind == "P", Unique: true, IsConstraint: true})
			constraints[tbl+"."+name] = true
		case "C":
			if !strings.HasPrefix(name, "SYS_C") {
				t.CheckConstraints = append(t.CheckConstraints, &metadata.CheckConstraintMetadata{Name: name})
			}
		case "R":
			t.ForeignKeys = append(t.ForeignKeys, &metadata.ForeignKeyMetadata{Name: name})
		}
		return nil
	})
	query(`SELECT table_name, index_name, uniqueness FROM all_indexes WHERE owner = 'SHOP' ORDER BY table_name, index_name`, func(r *sql.Rows) error {
		var tbl, name, uniqueness string
		if err := r.Scan(&tbl, &name, &uniqueness); err != nil {
			return err
		}
		if t, ok := tables[tbl]; ok && !constraints[tbl+"."+name] && !strings.HasPrefix(name, "SYS_") {
			t.Indexes = append(t.Indexes, &metadata.IndexMetadata{Name: name, Unique: uniqueness == "UNIQUE"})
		}
		return nil
	})
	query(`SELECT view_name FROM all_views WHERE owner = 'SHOP' ORDER BY view_name`, func(r *sql.Rows) error {
		var name string
		if err := r.Scan(&name); err != nil {
			return err
		}
		s.Views = append(s.Views, &metadata.ViewMetadata{Name: name})
		return nil
	})
	return &metadata.DatabaseSchemaMetadata{Name: "SHOP", Schemas: []*metadata.SchemaMetadata{s}}
}

// archiveExists reports whether the BBDATAARCHIVE user exists, which
// Bytebase reads from its own metadata as the backup database.
func archiveExists(t *testing.T, admin *sql.DB) bool {
	t.Helper()
	var n int
	if err := admin.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM dba_users WHERE username = 'BBDATAARCHIVE'").Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n > 0
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

// TestPriorBackupAgainstOracle builds the change's backup statements as
// Bytebase does and runs them on Oracle before the change. Each group of
// statements on one table succeeds or fails as a whole, as Bytebase runs
// it: a PriorBackup finding on a statement must mean its group fails or
// Bytebase refuses it, and a finding on the whole change must mean
// Bytebase refuses the change; for a complete case, every group that fails
// must have a finding on one of its statements, and every refusal of the
// change a finding on the change.
func TestPriorBackupAgainstOracle(t *testing.T) {
	admin := startOracle(t)
	tests := []struct {
		name      string
		change    string
		complete  bool
		noArchive bool
	}{
		{"a plain backup", "UPDATE t SET n = 5 WHERE id = 1;\nDELETE FROM nokey WHERE a = 1;\nDELETE FROM plain x;", true, false},
		{"two statements on one table", "UPDATE t SET n = 5 WHERE id = 1;\nUPDATE t x SET x.n = 6 WHERE x.id = 2;\nUPDATE t x SET (n, p_id) = (SELECT 1, 1 FROM dual) WHERE x.id = 1;", true, false},
		{"a view as the target", "UPDATE pv SET note = 'z' WHERE id = 1;\nUPDATE pv SET note = 'y' WHERE id = 2;", true, false},
		{"UPDATE and DELETE on one table", "UPDATE t SET n = 5 WHERE id = 1;\nDELETE FROM shop.t WHERE id = 2;", true, false},
		{"a LONG column", "DELETE FROM longt WHERE id = 1;", true, false},
		{"a quote in a table name", "DELETE FROM \"a'b\" WHERE id = 1;", true, false},
		{"a database link", "UPDATE t@lnk SET n = 1 WHERE id = 1;\nDELETE FROM nokey WHERE a = 1;", true, false},
		{"nothing to back up", "ALTER TABLE t ADD (m NUMBER);\nINSERT INTO nokey VALUES (3);\nBEGIN DELETE FROM nokey WHERE a = 3; END;\n/\nUPDATE (SELECT * FROM t) SET n = 1 WHERE id = 1;", true, false},
		{"no BBDATAARCHIVE", "UPDATE t SET n = 5 WHERE id = 1;", true, true},
		{"a table the change creates, named with its schema", "CREATE TABLE n2 (id NUMBER PRIMARY KEY);\nUPDATE shop.n2 SET id = 2 WHERE id = 1;", true, false},
		{"a table the change creates, in a subquery", "CREATE TABLE n2 (id NUMBER PRIMARY KEY);\nDELETE FROM t WHERE id IN (SELECT id FROM shop.n2);", true, false},
		{"a table the change creates, in a group of two", "CREATE TABLE n2 (id NUMBER PRIMARY KEY);\nUPDATE n2 SET id = 2 WHERE id = 1;\nUPDATE n2 SET id = 3 WHERE id = 2;", true, false},
		{"a table the change creates, unqualified, not judged", "CREATE TABLE n2 (id NUMBER PRIMARY KEY);\nUPDATE n2 SET id = 2 WHERE id = 1;", false, false},
		{"a table the change creates, which a public synonym names before", "CREATE TABLE pubt (id NUMBER PRIMARY KEY, a NUMBER);\nUPDATE pubt SET a = 1 WHERE id = 1;", true, false},
		{"a column the change adds, unqualified", "ALTER TABLE t ADD (gone NUMBER);\nUPDATE t SET n = 1 WHERE gone = 1;", true, false},
		{"a column the change renames to, through an alias", "ALTER TABLE p RENAME COLUMN note TO memo;\nUPDATE p x SET x.id = 3 WHERE x.memo = 'a';", true, false},
		{"a column of a table an UPDATE's FROM list joins", "ALTER TABLE p ADD (memo NUMBER);\nUPDATE t SET n = 1 FROM p WHERE t.p_id = p.id AND p.memo = 1;", true, false},
		{"a column of a table the change recreates", "DROP TABLE nokey;\nCREATE TABLE nokey (b NUMBER);\nDELETE FROM nokey WHERE b = 1;", true, false},
		{"a column the SET clause names, not copied", "ALTER TABLE t ADD (gone NUMBER);\nUPDATE t SET gone = 1 WHERE id = 1;", true, false},
		{"a column in a subquery, not judged", "ALTER TABLE t ADD (gone NUMBER);\nUPDATE plain SET v = 1 WHERE id IN (SELECT gone FROM t);", false, false},
		{"after the session moves, not judged", "ALTER SESSION SET CURRENT_SCHEMA = bbdataarchive;\nALTER TABLE shop.nokey ADD (gone NUMBER);\nDELETE FROM shop.nokey WHERE gone = 1;", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := resetOracle(t, admin, tt.noArchive)
			schema := syncSchema(t, conn)
			result, err := Review(context.Background(), tt.change, review.Options{
				Rules:  []review.Rule{review.PriorBackup},
				Change: review.Change{PriorBackup: true, MaxBackupSize: 1 << 20},
			}, []review.Target{{Schema: schema, BackupDatabaseExists: archiveExists(t, admin)}})
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
			outcome := runBackup(t, conn, tt.change, schema.Name, !tt.noArchive)
			for i, msg := range reported {
				if _, ok := outcome.failed[i]; !ok {
					t.Errorf("statement %d reported, but its backup runs: %s", i, msg)
				}
			}
			if tt.complete {
				if err, ok := outcome.failed[-1]; ok {
					if _, ok := reported[-1]; !ok {
						t.Errorf("the change not reported, but Bytebase refuses its backup: %v", err)
					}
				}
				for _, g := range outcome.groups {
					if g.err == nil {
						continue
					}
					hit := false
					for _, i := range g.statements {
						_, ok := reported[i]
						hit = hit || ok
					}
					if !hit {
						t.Errorf("statements %v not reported, but their backup fails: %v", g.statements, g.err)
					}
				}
			}
		})
	}
}

// backupOutcome is what Bytebase's backup of a change does: each group's
// error, nil when it runs, and the statements whose group fails or is
// refused, -1 for a refusal of the whole change.
type backupOutcome struct {
	groups []*backupGroup
	failed map[int]error
}

type backupGroup struct {
	key        string
	table      *ast.ObjectName
	schema     string
	statements []int
	dmls       []oracleDML
	err        error
}

// oracleDML is a statement the backup copies, as Bytebase extracts it.
type oracleDML struct {
	index int
	kind  string
	node  ast.StmtNode
	table *ast.ObjectName
	alias *ast.Alias
}

// runBackup builds the backup statements of a change as Bytebase builds
// them, and runs each group's: CREATE TABLE ... AS SELECT into
// BBDATAARCHIVE, then the COMMENT on the copy. Each group is refused or run
// as if it were the change's only one, so every reason the backup fails
// shows; Bytebase stops at the first.
//
// bytebase: backend/plugin/parser/plsql/backup.go — TransformDMLToSelect, prepareTransformation, generateSQL, generateSQLForTable, writeBranchSelect, writeSuffixSelectClause
// bytebase: backend/runner/taskrun/database_migrate_executor.go — backupData
func runBackup(t *testing.T, conn *sql.Conn, change, db string, archive bool) backupOutcome {
	t.Helper()
	out := backupOutcome{failed: make(map[int]error)}
	if !archive {
		out.failed[-1] = fmt.Errorf("backup database %q not found", backupDatabase)
		return out
	}
	var dmls []oracleDML
	for _, s := range parseChange(t, change) {
		switch v := s.node.(type) {
		case *ast.UpdateStmt:
			if v.Table != nil && v.Table.Name != "" {
				dmls = append(dmls, oracleDML{s.index, "UPDATE", v, v.Table, v.Alias})
			}
		case *ast.DeleteStmt:
			if v.Table != nil && v.Table.Name != "" {
				dmls = append(dmls, oracleDML{s.index, "DELETE", v, v.Table, v.Alias})
			}
		}
	}
	if len(dmls) == 0 {
		out.failed[-1] = errors.New("no parse results")
		return out
	}
	byKey := make(map[string]*backupGroup)
	for _, d := range dmls {
		schema := d.table.Schema
		if schema == "" {
			schema = db
		}
		key := fmt.Sprintf("%s.%s@%s", schema, d.table.Name, d.table.DBLink)
		g, ok := byKey[key]
		if !ok {
			g = &backupGroup{key: key, table: d.table, schema: schema}
			byKey[key] = g
			out.groups = append(out.groups, g)
		}
		g.statements = append(g.statements, d.index)
		g.dmls = append(g.dmls, d)
	}
	ctx := context.Background()
	for n, g := range out.groups {
		for _, d := range g.dmls {
			if d.kind != g.dmls[0].kind {
				g.err = fmt.Errorf("prior backup cannot handle statements with different types on the same table: %s", g.key)
			}
		}
		if g.err == nil && g.table.DBLink != "" {
			g.err = fmt.Errorf("prior backup does not support tables referenced through a database link: %s", g.key)
		}
		if g.err == nil {
			target := fmt.Sprintf("_%d_%s_%s", n, g.table.Name, g.schema)
			stmts := []string{
				backupQuery(change, g, target),
				fmt.Sprintf(`COMMENT ON TABLE "%s"."%s" IS 'task 1, source table (%s, %s)'`, backupDatabase, target, db, g.table.Name),
			}
			for _, stmt := range stmts {
				if _, err := conn.ExecContext(ctx, stmt); err != nil {
					g.err = fmt.Errorf("%s: %w", stmt, err)
					break
				}
			}
		}
		if g.err != nil {
			for _, i := range g.statements {
				out.failed[i] = g.err
			}
		}
	}
	return out
}

// backupQuery writes a group's CREATE TABLE ... AS SELECT as Bytebase
// writes it, without the final semicolon, which the driver does not take.
func backupQuery(sql string, g *backupGroup, target string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "CREATE TABLE \"%s\".\"%s\" AS\n", backupDatabase, target)
	if len(g.dmls) == 1 {
		b.WriteString(branchSelect(sql, g.dmls[0], "*"))
		return b.String()
	}
	fmt.Fprintf(&b, `SELECT "%s"."%s".* FROM "%s"."%s" WHERE ROWID IN (`, g.schema, g.table.Name, g.schema, g.table.Name)
	for i, d := range g.dmls {
		if i > 0 {
			b.WriteString(" UNION ")
		}
		b.WriteString(branchSelect(sql, d, "ROWID"))
	}
	b.WriteString(")")
	return b.String()
}

// branchSelect writes one statement's SELECT: its alias, or its table,
// qualified as written, then the statement's target and WHERE clause, and
// an UPDATE's FROM list when it has a WHERE clause, as written.
func branchSelect(sql string, d oracleDML, projection string) string {
	var head string
	switch {
	case d.alias != nil:
		head = fmt.Sprintf(`SELECT "%s".%s FROM `, d.alias.Name, projection)
	case d.table.Schema != "":
		head = fmt.Sprintf(`SELECT "%s"."%s".%s FROM `, d.table.Schema, d.table.Name, projection)
	default:
		head = fmt.Sprintf(`SELECT "%s".%s FROM `, d.table.Name, projection)
	}
	text := func(start, end int) string {
		if start < 0 || end > len(sql) || end <= start {
			return ""
		}
		return strings.TrimSpace(sql[start:end])
	}
	join := func(parts ...string) string {
		var keep []string
		for _, p := range parts {
			if p != "" {
				keep = append(keep, p)
			}
		}
		return strings.Join(keep, " ")
	}
	targetText := func(partition *ast.PartitionExtClause) (string, int) {
		end := d.table.Loc.End
		if partition != nil {
			end = partition.Loc.End
		}
		if d.alias != nil {
			end = d.alias.Loc.End
		}
		return text(d.table.Loc.Start, end), end
	}
	switch v := d.node.(type) {
	case *ast.UpdateStmt:
		suffix, _ := targetText(v.PartitionExt)
		if v.WhereClause != nil && v.SetClauses != nil && len(v.SetClauses.Items) > 0 {
			whereStart := ast.NodeLoc(v.SetClauses.Items[len(v.SetClauses.Items)-1]).End
			if v.FromClause != nil && len(v.FromClause.Items) > 0 {
				first, last := ast.NodeLoc(v.FromClause.Items[0]), ast.NodeLoc(v.FromClause.Items[len(v.FromClause.Items)-1])
				suffix += ", " + text(first.Start, last.End)
				whereStart = last.End
			}
			suffix = join(suffix, text(whereStart, ast.NodeLoc(v.WhereClause).End))
		}
		return head + suffix
	case *ast.DeleteStmt:
		target, end := targetText(v.PartitionExt)
		where := ""
		if v.WhereClause != nil {
			where = text(end, ast.NodeLoc(v.WhereClause).End)
		}
		return head + join(target, where)
	}
	return head
}

// TestScanRulesAgainstOracle runs each change on Oracle, statement by
// statement, and derives from the data dictionary what the constraint
// rules should report: DisallowDropConstraint, a primary key, unique,
// foreign key, or check constraint a statement removes from a table that
// stays; RequirePrimaryKey, a table that ends the change without a primary
// key, at the statement that made it so: the CREATE TABLE that made it
// without one, or the last statement that removed one it had before the
// change or since the change made it. A statement Oracle refuses ends the
// change, as it ends Bytebase's task. Every finding must match; for a
// complete case, every expected finding must be reported.
func TestScanRulesAgainstOracle(t *testing.T) {
	admin := startOracle(t)
	tests := []struct {
		name     string
		change   string
		complete bool
	}{
		{"constraints of the synced schema", "ALTER TABLE t DROP CONSTRAINT fk_p DROP CONSTRAINT ck_n;\nALTER TABLE t DROP CONSTRAINT uk_code;", true},
		{"a named NOT NULL constraint", "ALTER TABLE t DROP CONSTRAINT nn_code;", true},
		{"the primary key by name", "ALTER TABLE plain DROP CONSTRAINT pk_plain;", true},
		{"DROP PRIMARY KEY and DROP UNIQUE", "ALTER TABLE plain DROP PRIMARY KEY;\nALTER TABLE t DROP UNIQUE (code);", true},
		{"a constraint the change defines, then drops", "ALTER TABLE nokey ADD CONSTRAINT ck_a CHECK (a > 0);\nALTER TABLE nokey DROP CONSTRAINT ck_a;", true},
		{"a unique index that backs no constraint", "ALTER TABLE plain DROP CONSTRAINT ux_v;", true},
		{"a table created without a key", "CREATE TABLE a (id NUMBER);\nCREATE TABLE b (id NUMBER PRIMARY KEY);\nCREATE GLOBAL TEMPORARY TABLE g (id NUMBER);", true},
		{"CREATE TABLE AS SELECT, which copies no key", "CREATE TABLE c AS SELECT * FROM plain;\nCREATE TABLE d (id PRIMARY KEY) AS SELECT id FROM plain;", true},
		{"a key added later", "CREATE TABLE a (id NUMBER);\nALTER TABLE a ADD CONSTRAINT pk_a PRIMARY KEY (id);\nCREATE TABLE b (id NUMBER);\nALTER TABLE b MODIFY (id PRIMARY KEY);", true},
		{"a primary key dropped, then added back", "ALTER TABLE plain DROP PRIMARY KEY;\nALTER TABLE plain ADD CONSTRAINT pk_plain2 PRIMARY KEY (id);", true},
		{"a key a table without one gains and loses", "ALTER TABLE nokey ADD PRIMARY KEY (a);\nALTER TABLE nokey DROP PRIMARY KEY;", true},
		{"a key the change makes, then drops", "CREATE TABLE a (id NUMBER CONSTRAINT pk_a PRIMARY KEY);\nALTER TABLE a DROP CONSTRAINT pk_a;", true},
		{"a table renamed after losing its key", "ALTER TABLE plain DROP PRIMARY KEY;\nRENAME plain TO plain2;", true},
		{"a primary key dropped with CASCADE", "ALTER TABLE p DROP PRIMARY KEY CASCADE;\nALTER TABLE t DROP CONSTRAINT fk_p;", true},
		{"a name already taken fails its statement", "ALTER TABLE nokey ADD CONSTRAINT ck_n CHECK (a > 0);\nALTER TABLE t DROP CONSTRAINT uk_code;", true},
		{"a name the change frees first", "ALTER TABLE t DROP CONSTRAINT ck_n;\nALTER TABLE nokey ADD CONSTRAINT ck_n CHECK (a > 0);\nALTER TABLE nokey DROP CONSTRAINT ck_n;", true},
		{"a table the change recreates", "CREATE TABLE a (id NUMBER CONSTRAINT ck_a CHECK (id > 0));\nDROP TABLE a;\nCREATE TABLE a (id NUMBER);\nALTER TABLE a DROP CONSTRAINT ck_a;", true},
		{"a synced table the change recreates", "DROP TABLE t;\nCREATE TABLE t (id NUMBER CONSTRAINT ck_n CHECK (id > 0));\nALTER TABLE t DROP CONSTRAINT ck_n;", true},
		{"a column dropped with its constraint, not judged", "ALTER TABLE t DROP COLUMN n;\nALTER TABLE t DROP CONSTRAINT uk_code;", false},
		{"a renamed table, not judged", "RENAME t TO t2;\nALTER TABLE t2 DROP CONSTRAINT ck_n;", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := resetOracle(t, admin, false)
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
// constraint of the tables of schema SHOP, as table/type/name.
func constraintSet(t *testing.T, conn *sql.Conn) map[string]string {
	t.Helper()
	out := make(map[string]string)
	rows, err := conn.QueryContext(context.Background(), `SELECT c.table_name, c.constraint_type, c.constraint_name FROM user_constraints c JOIN user_tables t ON t.table_name = c.table_name WHERE c.constraint_type IN ('P', 'U', 'R', 'C') AND t.dropped = 'NO'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tbl, kind, name string
		if err := rows.Scan(&tbl, &kind, &name); err != nil {
			t.Fatal(err)
		}
		out[tbl+"/"+kind+"/"+name] = tbl
	}
	return out
}

// primaryKeyed reads the tables of schema SHOP that have a primary key,
// and all of them.
func primaryKeyed(t *testing.T, conn *sql.Conn) (keyed, all map[string]bool) {
	t.Helper()
	keyed, all = make(map[string]bool), make(map[string]bool)
	rows, err := conn.QueryContext(context.Background(), `SELECT t.table_name, (SELECT COUNT(*) FROM user_constraints c WHERE c.table_name = t.table_name AND c.constraint_type = 'P') FROM user_tables t WHERE t.dropped = 'NO' AND t.temporary = 'N'`)
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
		all[name] = true
		if has > 0 {
			keyed[name] = true
		}
	}
	return keyed, all
}

// scanTruth runs the change statement by statement, until Oracle refuses
// one, and returns the findings the data dictionary calls for, as "RULE
// statement".
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
			t.Logf("statement %d refused, which ends the change: %v", s.index, err)
			break
		}
		after := constraintSet(t, conn)
		keyedNow, allNow := primaryKeyed(t, conn)
		for k, tbl := range before {
			if _, ok := after[k]; !ok && allNow[tbl] {
				want[fmt.Sprintf("%s %d", review.DisallowDropConstraint, s.index)] = true
			}
		}
		for name := range keyed {
			if allNow[name] && !keyedNow[name] {
				lostAt[name] = s.index
			}
		}
		if v, ok := s.node.(*ast.CreateTableStmt); ok && judgedTable(v) {
			createdAt[v.Name.Name] = s.index
		}
	}
	keyedEnd, all := primaryKeyed(t, conn)
	for name := range all {
		if keyedEnd[name] {
			continue
		}
		lost, hasLost := lostAt[name]
		if made, ok := createdAt[name]; ok {
			if hasLost && lost > made {
				made = lost
			}
			want[fmt.Sprintf("%s %d", review.RequirePrimaryKey, made)] = true
			continue
		}
		if keyedBefore[name] && hasLost {
			want[fmt.Sprintf("%s %d", review.RequirePrimaryKey, lost)] = true
		}
	}
	return want
}

// TestRenameAgainstOracle runs each change on Oracle, statement by
// statement, and reads before each RENAME what the data dictionary says
// the old name is. A DisallowRename finding on a RENAME must mean it
// renamed a table or a view; for a complete case, every RENAME of one must
// be reported. A statement Oracle refuses ends the change.
func TestRenameAgainstOracle(t *testing.T) {
	admin := startOracle(t)
	tests := []struct {
		name     string
		change   string
		complete bool
	}{
		{"synced objects", "RENAME t TO t2;\nRENAME pv TO pv2;\nCREATE SEQUENCE sq;\nRENAME sq TO sq2;", true},
		{"what the change makes", "CREATE TABLE n1 (id NUMBER);\nRENAME n1 TO n2;\nRENAME n2 TO n3;\nCREATE VIEW w AS SELECT 1 a FROM dual;\nRENAME w TO w2;\nCREATE SYNONYM sy FOR p;\nRENAME sy TO sy2;", true},
		{"a name dropped and made again as another kind", "DROP TABLE nokey;\nCREATE VIEW nokey AS SELECT 1 a FROM dual;\nRENAME nokey TO nk;\nDROP VIEW pv;\nCREATE SEQUENCE pv;\nRENAME pv TO pv2;", true},
		{"a materialized view, which RENAME refuses", "CREATE MATERIALIZED VIEW mv AS SELECT * FROM p;\nRENAME mv TO mv2;", true},
		{"a table renamed by ALTER TABLE, then by RENAME", "ALTER TABLE plain RENAME TO plain2;\nRENAME plain2 TO plain3;", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := resetOracle(t, admin, false)
			schema := syncSchema(t, conn)
			result, err := Review(context.Background(), tt.change, review.Options{Rules: []review.Rule{review.DisallowRename}}, []review.Target{{Schema: schema}})
			if err != nil {
				t.Fatalf("Review: %v", err)
			}
			got := make(map[int]bool)
			stmts := parseChange(t, tt.change)
			for _, f := range result.Findings {
				if f.Rule == review.Syntax {
					t.Fatalf("syntax: %s", f.Message)
				}
				if _, ok := stmts[f.Statement].node.(*ast.RenameStmt); ok {
					got[f.Statement] = true
				}
			}
			want := make(map[int]bool)
			ctx := context.Background()
			for _, s := range stmts {
				var kinds []string
				if v, ok := s.node.(*ast.RenameStmt); ok {
					rows, err := conn.QueryContext(ctx, "SELECT object_type FROM user_objects WHERE object_name = :1", v.OldName.Name)
					if err != nil {
						t.Fatal(err)
					}
					for rows.Next() {
						var kind string
						if err := rows.Scan(&kind); err != nil {
							t.Fatal(err)
						}
						kinds = append(kinds, kind)
					}
					rows.Close()
				}
				if _, err := conn.ExecContext(ctx, tt.change[s.loc.Start:s.loc.End]); err != nil {
					t.Logf("statement %d refused, which ends the change: %v", s.index, err)
					break
				}
				if len(kinds) == 1 && (kinds[0] == "TABLE" || kinds[0] == "VIEW") {
					want[s.index] = true
				}
			}
			for i := range got {
				if !want[i] {
					t.Errorf("statement %d reported, but it renames no table or view", i)
				}
			}
			if tt.complete {
				for i := range want {
					if !got[i] {
						t.Errorf("statement %d renames a table or view, and is not reported", i)
					}
				}
			}
		})
	}
}

// TestRequireIsNullAgainstOracle confirms on Oracle that a comparison the
// rule reports is never true, and one it does not report can be: each
// condition is evaluated over rows holding NULL, the empty string, and
// values, in SQL and in PL/SQL. A reported condition must hold for no row,
// and neither may its negation, unknown for every row; an unreported one
// must hold for some row.
func TestRequireIsNullAgainstOracle(t *testing.T) {
	admin := startOracle(t)
	conn := resetOracle(t, admin, false)
	ctx := context.Background()
	for _, stmt := range []string{
		"CREATE TABLE probe (n NUMBER, s VARCHAR2(10), c CHAR(3))",
		"INSERT INTO probe VALUES (NULL, NULL, NULL)",
		"INSERT INTO probe VALUES (1, 'a', 'a')",
		"INSERT INTO probe VALUES (0, ' ', ' ')",
	} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	tests := []struct {
		cond   string
		report bool
	}{
		{"n = NULL", true},
		{"NULL = n", true},
		{"n <> NULL", true},
		{"n != NULL", true},
		{"n ^= NULL", true},
		{"(n) = (NULL)", true},
		{"s = ''", true},
		{"s <> ''", true},
		{"c = ''", true},
		{"c != ''", true},
		{"'' = ''", true},
		{"s = N''", true},
		{"s = q'[]'", true},
		{"NULL = NULL", true},
		{"n = 1", false},
		{"s = ' '", false},
		{"c = ' '", false},
		{"n IS NULL", false},
		{"s IS NULL", false},
		{"n <> 1", false},
	}
	for _, tt := range tests {
		t.Run(tt.cond, func(t *testing.T) {
			sql := "SELECT * FROM probe WHERE " + tt.cond
			result, err := Review(context.Background(), sql, review.Options{Rules: []review.Rule{review.RequireIsNull}}, nil)
			if err != nil {
				t.Fatal(err)
			}
			if got := len(result.Findings) > 0; got != tt.report {
				t.Fatalf("reported = %v, want %v: %v", got, tt.report, result.Findings)
			}
			count := func(where string) int {
				var n int
				if err := conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM probe WHERE "+where).Scan(&n); err != nil {
					t.Fatalf("%s: %v", where, err)
				}
				return n
			}
			pass, negated := count(tt.cond), count("NOT ("+tt.cond+")")
			plsql := plsqlTrueCount(t, conn, tt.cond)
			switch {
			case tt.report && (pass > 0 || negated > 0 || plsql > 0):
				t.Errorf("reported, but Oracle finds it true for %d rows, its negation for %d, and PL/SQL for %d", pass, negated, plsql)
			case !tt.report && pass == 0 && plsql == 0:
				t.Errorf("not reported, and Oracle finds it true for no row")
			}
		})
	}
}

// plsqlTrueCount counts the probe rows for which PL/SQL's IF finds the
// condition true, the columns read through a cursor record.
func plsqlTrueCount(t *testing.T, conn *sql.Conn, cond string) int {
	t.Helper()
	for _, col := range []string{"n", "s", "c"} {
		cond = replaceWord(cond, col, "r."+col)
	}
	block := `DECLARE k NUMBER := 0;
BEGIN
  FOR r IN (SELECT * FROM probe) LOOP
    IF ` + cond + ` THEN k := k + 1; END IF;
  END LOOP;
  :1 := k;
END;`
	var n int64
	if _, err := conn.ExecContext(context.Background(), block, sql.Out{Dest: &n}); err != nil {
		t.Fatalf("%s: %v", block, err)
	}
	return int(n)
}

// replaceWord replaces a whole word outside quotes.
func replaceWord(s, word, with string) string {
	var b strings.Builder
	inQuote := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\'' {
			inQuote = !inQuote
		}
		isWord := func(c byte) bool {
			return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_' || c >= '0' && c <= '9'
		}
		if !inQuote && strings.HasPrefix(s[i:], word) && (i == 0 || !isWord(s[i-1])) && (i+len(word) == len(s) || !isWord(s[i+len(word)])) {
			b.WriteString(with)
			i += len(word) - 1
			continue
		}
		b.WriteByte(c)
	}
	return b.String()
}
