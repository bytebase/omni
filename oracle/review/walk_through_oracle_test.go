package review

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/review/walk/walktest"
)

// walkBase is the SHOP schema the curated cases start from: tables with
// keys, a named NOT NULL, a check, a foreign key, an index, and an
// invisible column; a table with none; a view, a sequence, a synonym, a
// function, a procedure, a package, a type, and a materialized view with
// an index.
var walkBase = []string{
	"CREATE TABLE p (id NUMBER CONSTRAINT pk_p PRIMARY KEY, note VARCHAR2(20))",
	"CREATE TABLE t (id NUMBER CONSTRAINT pk_t PRIMARY KEY, code VARCHAR2(20) CONSTRAINT nn_code NOT NULL CONSTRAINT uk_code UNIQUE, p_id NUMBER CONSTRAINT fk_p REFERENCES p (id), n NUMBER CONSTRAINT ck_n CHECK (n > 0), h NUMBER INVISIBLE)",
	"CREATE INDEX ix_n ON t (n)",
	"CREATE TABLE nokey (a NUMBER)",
	"CREATE VIEW pv AS SELECT id, note FROM p",
	"CREATE SEQUENCE sq",
	"CREATE SYNONYM syn FOR t",
	"CREATE OR REPLACE FUNCTION f RETURN NUMBER AS BEGIN RETURN 1; END;",
	"CREATE OR REPLACE PROCEDURE pr AS BEGIN NULL; END;",
	"CREATE OR REPLACE PACKAGE pkg AS x NUMBER; END;",
	"CREATE TYPE ty AS OBJECT (x NUMBER, y NUMBER)",
	"CREATE MATERIALIZED VIEW mv AS SELECT id, note FROM p",
	"CREATE INDEX mv_ix ON mv (note)",
}

// walkCase is a change on walkBase, after setup when it has one. With
// asShop, the change runs as SHOP, which owns the objects RENAME renames;
// otherwise as an administrator with CURRENT_SCHEMA = SHOP, as Bytebase
// runs it.
type walkCase struct {
	name     string
	setup    []string
	change   string
	complete bool
	asShop   bool
}

var walkCases = []walkCase{
	// Names a CREATE takes, in the shared object namespace.
	{name: "the object namespace", change: "CREATE TABLE t (a NUMBER);\nCREATE TABLE pv (a NUMBER);\nCREATE TABLE mv (a NUMBER);\nCREATE SEQUENCE f;\nCREATE VIEW sq AS SELECT 1 a FROM dual;\nCREATE SYNONYM pkg FOR t;\nCREATE TYPE pr AS OBJECT (a NUMBER);\n/\nCREATE OR REPLACE PROCEDURE f AS BEGIN NULL; END;\n/\nCREATE OR REPLACE VIEW t AS SELECT 1 a FROM dual;\nCREATE OR REPLACE SYNONYM pv FOR t;\nCREATE TABLE x (a NUMBER);\nCREATE TABLE x (a NUMBER);\nCREATE INDEX t ON nokey (a);\nCREATE TRIGGER t BEFORE INSERT ON nokey FOR EACH ROW BEGIN NULL; END;", complete: true},
	{name: "OR REPLACE of its own kind", change: "CREATE OR REPLACE VIEW pv AS SELECT 1 a FROM dual;\nCREATE OR REPLACE FUNCTION f RETURN NUMBER AS BEGIN RETURN 2; END;\n/\nCREATE OR REPLACE SYNONYM syn FOR p;\nCREATE OR REPLACE PACKAGE pkg AS y NUMBER; END;\n/\nCREATE OR REPLACE TYPE ty2 AS OBJECT (a NUMBER);\n/\nCREATE OR REPLACE TYPE ty2 AS OBJECT (b NUMBER);", complete: true},
	{name: "IF NOT EXISTS", change: "CREATE TABLE IF NOT EXISTS t (a NUMBER);\nCREATE TABLE IF NOT EXISTS pv (a NUMBER);\nCREATE TABLE IF NOT EXISTS t (a NUMBER, b NUMBER CONSTRAINT ck_n CHECK (b > 0));\nCREATE TABLE IF NOT EXISTS t (a NUMBER, a NUMBER);\nCREATE SEQUENCE IF NOT EXISTS sq;\nCREATE INDEX IF NOT EXISTS ix_n ON nokey (a);\nCREATE TABLE IF NOT EXISTS x (a NUMBER);\nCREATE TABLE x (a NUMBER);", complete: true},
	{name: "a new table's text", change: "CREATE TABLE x (a NUMBER, a NUMBER);\nCREATE TABLE x (a NUMBER CONSTRAINT c1 CHECK (a > 0), b NUMBER CONSTRAINT c1 NOT NULL);\nCREATE TABLE x (a NUMBER PRIMARY KEY, b NUMBER, PRIMARY KEY (b));\nCREATE TABLE x (a NUMBER CONSTRAINT ck_n CHECK (a > 0));\nCREATE TABLE x (a NUMBER CONSTRAINT nn_code NOT NULL);\nCREATE TABLE x (a NUMBER, CONSTRAINT fk_p FOREIGN KEY (a) REFERENCES p (id));", complete: true},
	{name: "a new table holds what it declares", change: "CREATE TABLE x (id NUMBER CONSTRAINT pk_x PRIMARY KEY, a NUMBER CONSTRAINT nn_a NOT NULL, b NUMBER CONSTRAINT ck_b CHECK (b > 0));\nALTER TABLE x ADD (a NUMBER);\nALTER TABLE x DROP COLUMN zz;\nALTER TABLE x ADD PRIMARY KEY (b);\nALTER TABLE nokey ADD CONSTRAINT nn_a CHECK (a > 0);\nALTER TABLE x DROP CONSTRAINT ck_b;\nALTER TABLE x DROP CONSTRAINT ck_b;\nCREATE TABLE y (a NUMBER);\nALTER TABLE y DROP PRIMARY KEY;", complete: true},
	{name: "names too long", change: "CREATE TABLE " + strings.Repeat("a", 129) + " (a NUMBER);\nDROP TABLE " + strings.Repeat("b", 129) + ";\nALTER TABLE nokey ADD (" + strings.Repeat("c", 129) + " NUMBER);\nCREATE TABLE " + strings.Repeat("d", 128) + " (a NUMBER);", complete: true},
	// ALTER TABLE.
	{name: "ALTER TABLE on what is no table", change: "ALTER TABLE nosuch ADD (x NUMBER);\nALTER TABLE pv ADD (x NUMBER);\nALTER TABLE syn ADD (x NUMBER);\nALTER TABLE sq ADD (x NUMBER);\nALTER TABLE f ADD (x NUMBER);\nALTER TABLE pubt ADD (x NUMBER);\nALTER TABLE mv LOGGING;", complete: true},
	{name: "columns", change: "ALTER TABLE t ADD (code NUMBER);\nALTER TABLE t ADD (h NUMBER);\nALTER TABLE t DROP COLUMN zz;\nALTER TABLE t DROP (n, zz);\nALTER TABLE t SET UNUSED COLUMN zz;\nALTER TABLE t MODIFY (zz NUMBER);\nALTER TABLE t RENAME COLUMN zz TO z2;\nALTER TABLE t RENAME COLUMN n TO code;\nALTER TABLE t RENAME COLUMN n TO n;\nALTER TABLE t DROP COLUMN h;\nALTER TABLE t DROP COLUMN h;", complete: true},
	{name: "ALTER TABLE sees the table as it was", change: "ALTER TABLE t ADD (x NUMBER) MODIFY (x NUMBER);\nALTER TABLE t ADD (x NUMBER, x NUMBER);\nALTER TABLE t DROP (n, n);\nALTER TABLE t ADD (CONSTRAINT c1 CHECK (n > 0), CONSTRAINT c1 CHECK (n > 1));\nALTER TABLE nokey ADD PRIMARY KEY (a) ADD (b NUMBER PRIMARY KEY);\nALTER TABLE t DROP CONSTRAINT ck_n DROP CONSTRAINT ck_n;\nALTER TABLE t DROP CONSTRAINT ck_n;", complete: true},
	{name: "constraints and keys", change: "ALTER TABLE nokey ADD CONSTRAINT ck_n CHECK (a > 0);\nALTER TABLE nokey ADD CONSTRAINT pk_p CHECK (a > 0);\nALTER TABLE nokey ADD (b NUMBER CONSTRAINT nn_code NOT NULL);\nALTER TABLE t ADD PRIMARY KEY (n);\nALTER TABLE t MODIFY (n PRIMARY KEY);\nALTER TABLE nokey DROP PRIMARY KEY;\nALTER TABLE p DROP CONSTRAINT ck_n;\nALTER TABLE p RENAME CONSTRAINT ck_n TO ck_m;\nALTER TABLE p ENABLE CONSTRAINT ck_n;\nALTER TABLE p DISABLE CONSTRAINT ck_n;\nALTER TABLE p MODIFY CONSTRAINT ck_n DISABLE;\nALTER TABLE t RENAME CONSTRAINT ck_n TO fk_p;\nALTER TABLE t RENAME CONSTRAINT ck_n TO ck_n;", complete: true},
	{name: "constraint effects", change: "ALTER TABLE t DROP CONSTRAINT ck_n;\nALTER TABLE t DROP CONSTRAINT ck_n;\nALTER TABLE nokey ADD CONSTRAINT ck_n CHECK (a > 0);\nALTER TABLE nokey ADD CONSTRAINT ck_n CHECK (a > 1);\nALTER TABLE nokey RENAME CONSTRAINT ck_n TO ck_m;\nALTER TABLE nokey DROP CONSTRAINT ck_n;\nALTER TABLE p ADD CONSTRAINT ck_n CHECK (id > 0);\nALTER TABLE nokey ADD CONSTRAINT pk_nokey PRIMARY KEY (a);\nALTER TABLE nokey ADD PRIMARY KEY (a);\nALTER TABLE nokey DROP PRIMARY KEY;\nALTER TABLE nokey DROP PRIMARY KEY;", complete: true},
	{name: "a dropped column takes its constraints", change: "ALTER TABLE t DROP COLUMN n;\nALTER TABLE nokey ADD CONSTRAINT ck_n CHECK (a > 0);\nALTER TABLE t DROP CONSTRAINT ck_n;"},
	{name: "a dropped key takes its index", change: "ALTER TABLE t DROP CONSTRAINT uk_code;\nCREATE INDEX uk_code ON nokey (a);\nALTER TABLE t DISABLE CONSTRAINT pk_t;\nCREATE INDEX pk_t ON nokey (a);\nALTER TABLE t ENABLE CONSTRAINT pk_t;"},
	{name: "RENAME TO", change: "ALTER TABLE t RENAME TO p;\nALTER TABLE t RENAME TO t;\nALTER TABLE t RENAME TO pv;\nALTER TABLE t RENAME TO t2;\nALTER TABLE t ADD (x NUMBER);\nALTER TABLE t2 ADD (code NUMBER);\nALTER TABLE t2 DROP COLUMN zz;\nCREATE TABLE t (a NUMBER);\nALTER TABLE t2 DROP CONSTRAINT ck_n;\nALTER TABLE t2 DROP CONSTRAINT ck_n;", complete: true},
	// DROP and TRUNCATE.
	{name: "DROP TABLE", change: "DROP TABLE nosuch;\nDROP TABLE pv;\nDROP TABLE syn;\nDROP TABLE sq;\nDROP TABLE ty;\nDROP TABLE mv;\nDROP TABLE IF EXISTS mv;\nDROP TABLE IF EXISTS pv;\nDROP TABLE IF EXISTS nosuch;\nDROP TABLE pubt;\nDROP TABLE nokey;\nDROP TABLE nokey;\nALTER TABLE nokey ADD (a NUMBER);", complete: true},
	{name: "a dropped table frees its names", change: "DROP TABLE t CASCADE CONSTRAINTS;\nCREATE INDEX ix_n ON nokey (a);\nALTER TABLE nokey ADD CONSTRAINT ck_n CHECK (a > 0);\nCREATE TABLE t (z NUMBER);\nALTER TABLE t DROP COLUMN code;\nALTER TABLE t ADD (z NUMBER);", complete: true},
	{name: "an index and a trigger in another schema than their table", setup: []string{"CREATE TABLE other.o (a NUMBER)"}, change: "CREATE INDEX oi ON other.o (a);\nCREATE TRIGGER otr BEFORE INSERT ON other.o FOR EACH ROW BEGIN NULL; END;\n/\nCREATE INDEX oi ON nokey (a);\nCREATE TRIGGER otr BEFORE INSERT ON nokey FOR EACH ROW BEGIN NULL; END;\n/\nDROP TABLE other.o;\nCREATE INDEX oi ON nokey (a);\nCREATE TRIGGER otr BEFORE INSERT ON nokey FOR EACH ROW BEGIN NULL; END;", complete: true},
	{name: "CASCADE CONSTRAINTS takes other tables' foreign keys", change: "DROP TABLE p CASCADE CONSTRAINTS;\nALTER TABLE nokey ADD CONSTRAINT fk_p CHECK (a > 0);"},
	{name: "TRUNCATE", change: "TRUNCATE TABLE nosuch;\nTRUNCATE TABLE pv;\nTRUNCATE TABLE sq;\nTRUNCATE TABLE syn;\nTRUNCATE TABLE f;\nTRUNCATE TABLE ty;\nTRUNCATE TABLE pubt;\nTRUNCATE TABLE mv;\nTRUNCATE TABLE nokey;", complete: true},
	{name: "DROP of other kinds", change: "DROP VIEW t;\nDROP VIEW nosuch;\nDROP VIEW syn;\nDROP VIEW IF EXISTS t;\nDROP SEQUENCE t;\nDROP SEQUENCE nosuch;\nDROP SEQUENCE IF EXISTS t;\nALTER SEQUENCE pv INCREMENT BY 2;\nDROP SYNONYM t;\nDROP PROCEDURE t;\nDROP FUNCTION pr;\nDROP PACKAGE f;\nDROP TYPE pkg;\nDROP VIEW pv;\nDROP VIEW pv;\nDROP SEQUENCE sq;\nDROP SEQUENCE sq;\nDROP SYNONYM syn;\nDROP SYNONYM syn;\nDROP FUNCTION f;\nDROP FUNCTION f;\nDROP PACKAGE pkg;\nDROP PACKAGE pkg;", complete: true},
	{name: "DROP of names the sync does not list", change: "DROP SYNONYM nosuch;\nDROP PROCEDURE nosuch;\nDROP TYPE nosuch;\nDROP INDEX nosuch;\nDROP TRIGGER nosuch;\nDROP PUBLIC SYNONYM nosuch;"},
	// Indexes and triggers.
	{name: "indexes", change: "CREATE INDEX ix_n ON nokey (a);\nCREATE INDEX mv_ix ON nokey (a);\nCREATE INDEX i ON pv (id);\nCREATE INDEX i ON sq (a);\nCREATE INDEX i ON f (a);\nCREATE INDEX i ON syn (p_id);\nCREATE INDEX i2 ON nokey (a);\nDROP INDEX i2;\nDROP INDEX i2;\nALTER INDEX i REBUILD;\nALTER INDEX i RENAME TO ix_n;\nALTER INDEX i RENAME TO i;\nALTER INDEX i RENAME TO i3;\nDROP INDEX i;\nCREATE INDEX i3 ON p (note);", complete: true},
	{name: "triggers", change: "CREATE TRIGGER tr BEFORE INSERT ON nokey FOR EACH ROW BEGIN NULL; END;\n/\nCREATE TRIGGER tr BEFORE INSERT ON p FOR EACH ROW BEGIN NULL; END;\n/\nCREATE OR REPLACE TRIGGER tr BEFORE INSERT ON p FOR EACH ROW BEGIN NULL; END;\n/\nCREATE OR REPLACE TRIGGER tr BEFORE UPDATE ON nokey FOR EACH ROW BEGIN NULL; END;\n/\nALTER TRIGGER tr DISABLE;\nALTER TRIGGER tr RENAME TO tr;\nALTER TRIGGER tr RENAME TO tr2;\nALTER TRIGGER tr ENABLE;\nDROP TRIGGER tr2;\nDROP TRIGGER tr2;", complete: true},
	{name: "a trigger with errors is created", change: "CREATE TRIGGER tr BEFORE INSERT ON nokey FOR EACH ROW BEGIN nosuch_call; END;\n/\nCREATE TRIGGER tr BEFORE INSERT ON nokey FOR EACH ROW BEGIN NULL; END;", complete: true},
	{name: "a view's triggers and constraints go with it", change: "CREATE VIEW v2 AS SELECT id FROM p WITH CHECK OPTION CONSTRAINT vck;\nCREATE TRIGGER vtr INSTEAD OF INSERT ON v2 FOR EACH ROW BEGIN NULL; END;\n/\nALTER TABLE nokey ADD CONSTRAINT vck CHECK (a > 0);\nCREATE TRIGGER vtr BEFORE INSERT ON nokey FOR EACH ROW BEGIN NULL; END;\n/\nCREATE OR REPLACE VIEW v2 AS SELECT id FROM p;\nALTER TABLE nokey ADD CONSTRAINT vck CHECK (a > 0);\nCREATE TRIGGER vtr BEFORE INSERT ON nokey FOR EACH ROW BEGIN NULL; END;", complete: true},
	// RENAME, as the owner.
	{name: "RENAME", asShop: true, change: "RENAME pv TO t;\nRENAME t TO t;\nRENAME f TO g;\nRENAME mv TO g;\nRENAME pv TO pv2;\nDROP VIEW pv;\nRENAME sq TO sq2;\nDROP SEQUENCE sq;\nRENAME syn TO syn2;\nDROP SYNONYM syn;\nRENAME nokey TO nk;\nALTER TABLE nk ADD (a NUMBER);\nCREATE TABLE nk (a NUMBER);", complete: true},
	{name: "RENAME by another user", change: "RENAME nosuch TO x;\nRENAME pv TO t;\nRENAME nokey TO nk;"},
	// What the walk leaves Unknown by design (spec §3, §7): names the sync
	// does not list, the targets of statements that follow synonyms, and
	// constraints a table may have that the sync does not show.
	{name: "unknown by design", setup: []string{"CREATE TABLE other.t (a NUMBER)"}, change: "ALTER SESSION SET CURRENT_SCHEMA = other;\nCREATE TABLE t (a NUMBER);\nALTER SESSION SET CURRENT_SCHEMA = shop;\nCREATE TABLE syn (a NUMBER);\nCREATE TABLE ty (a NUMBER);\nALTER TABLE t DROP CONSTRAINT ck_n DROP CONSTRAINT nosuch;\nCREATE INDEX i ON nosuch (a);\nCREATE INDEX i ON syn (n);\nRENAME nosuch TO x;\nRENAME ty TO g;\nCOMMENT ON COLUMN t.nosuch IS 'x';\nINSERT INTO nosuch VALUES (1);\nGRANT SELECT ON nosuch TO PUBLIC;\nDROP INDEX pk_t;"},
	// The session.
	{name: "CURRENT_SCHEMA", setup: []string{"CREATE TABLE other.t (a NUMBER)"}, change: "ALTER SESSION SET CURRENT_SCHEMA = other;\nCREATE TABLE x (a NUMBER);\nCREATE TABLE x (a NUMBER);\nCREATE TABLE shop.t (a NUMBER);\nALTER SESSION SET CURRENT_SCHEMA = shop;\nCREATE TABLE t (a NUMBER);\nCREATE TABLE other.x (a NUMBER);", complete: true},
	// DML, which may run code that changes a schema.
	{name: "a trigger may run DDL", setup: []string{
		"CREATE OR REPLACE PROCEDURE mk AS PRAGMA AUTONOMOUS_TRANSACTION; BEGIN EXECUTE IMMEDIATE 'CREATE TABLE shop.made (a NUMBER)'; END;",
		"CREATE TRIGGER mk_tr BEFORE INSERT ON nokey FOR EACH ROW BEGIN mk; END;",
	}, change: "INSERT INTO nokey VALUES (1);\nCREATE TABLE made (a NUMBER);"},
	{name: "a function may run DDL", setup: []string{
		"CREATE OR REPLACE FUNCTION mkf RETURN NUMBER AS PRAGMA AUTONOMOUS_TRANSACTION; BEGIN EXECUTE IMMEDIATE 'CREATE TABLE shop.made (a NUMBER)'; RETURN 1; END;",
	}, change: "SELECT mkf FROM dual;\nCREATE TABLE made (a NUMBER);\nCREATE TABLE x AS SELECT mkf() a FROM dual;\nDROP TABLE made;\nCREATE TABLE made (a NUMBER);"},
	{name: "DML on what the walk knows", change: "INSERT INTO nokey VALUES (1);\nUPDATE t SET n = 2 WHERE id = 1;\nDELETE FROM nokey;\nMERGE INTO nokey USING p ON (nokey.a = p.id) WHEN NOT MATCHED THEN INSERT (a) VALUES (p.id);\nSELECT upper(note) FROM p;\nCREATE TABLE t (a NUMBER);", complete: true},
}

// walkDeadline bounds every statement and query of the oracle: the
// container hung once in CI.
const walkDeadline = 2 * time.Minute

func walkContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), walkDeadline)
}

func TestWalkThroughAgainstOracle(t *testing.T) {
	admin := startOracle(t)
	var found, rejected int
	for _, tt := range walkCases {
		t.Run(tt.name, func(t *testing.T) {
			w := resetWalk(t, admin, "SHOP", "OTHER", walkBase)
			for _, q := range tt.setup {
				w.exec(t, w.adminConn, q)
			}
			conn := w.adminConn
			if tt.asShop {
				conn = w.shopConn
			}
			f, r := judgeWalk(t, w, conn, tt.change, tt.complete)
			found += f
			rejected += r
		})
	}
	t.Logf("found %d of the %d statements the server rejects", found, rejected)
}

// walkSchemas is a schema of the change and the other schema a change may
// name, made afresh for one change, with their users' numbers.
type walkSchemas struct {
	admin     *sql.DB
	shop      string
	other     string
	users     string
	adminConn *sql.Conn
	shopConn  *sql.Conn
}

func (w *walkSchemas) exec(t testing.TB, conn *sql.Conn, q string) {
	t.Helper()
	ctx, cancel := walkContext()
	defer cancel()
	if _, err := conn.ExecContext(ctx, q); err != nil {
		t.Fatalf("%s: %v", q, err)
	}
}

// resetWalk recreates the users shop and other, and shop's objects from
// base, made by an administrator whose CURRENT_SCHEMA is shop. It returns
// that administrator's session, as Bytebase opens it for the change, and
// one as shop.
func resetWalk(t testing.TB, admin *sql.DB, shop, other string, base []string) *walkSchemas {
	t.Helper()
	ctx, cancel := walkContext()
	defer cancel()
	for _, q := range []string{
		fmt.Sprintf(`BEGIN
  FOR u IN (SELECT username FROM dba_users WHERE username IN ('%s', '%s')) LOOP
    EXECUTE IMMEDIATE 'DROP USER ' || u.username || ' CASCADE';
  END LOOP;
END;`, shop, other),
		"CREATE USER " + shop + " IDENTIFIED BY " + oraclePassword + " QUOTA UNLIMITED ON users",
		"GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE MATERIALIZED VIEW, CREATE SEQUENCE, CREATE SYNONYM, CREATE PROCEDURE, CREATE TRIGGER, CREATE TYPE TO " + shop,
		"CREATE USER " + other + " IDENTIFIED BY " + oraclePassword + " QUOTA UNLIMITED ON users",
	} {
		if _, err := admin.ExecContext(ctx, q); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
	}
	w := &walkSchemas{admin: admin, shop: shop, other: other}
	var ids []string
	rows, err := admin.QueryContext(ctx, "SELECT user_id FROM dba_users WHERE username IN (:1, :2)", shop, other)
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			t.Fatal(err)
		}
		ids = append(ids, id)
	}
	rows.Close()
	w.users = "(" + strings.Join(ids, ", ") + ")"
	w.adminConn = walkConn(t, admin, shop)
	shopDB, err := openOracle(shop)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { shopDB.Close() })
	w.shopConn = walkConn(t, shopDB, "")
	for _, q := range base {
		w.exec(t, w.adminConn, q)
	}
	return w
}

// walkConn takes a session of its own from a pool, with CURRENT_SCHEMA set
// when schema is not empty.
func walkConn(t testing.TB, db *sql.DB, schema string) *sql.Conn {
	t.Helper()
	ctx, cancel := walkContext()
	defer cancel()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	if schema != "" {
		if _, err := conn.ExecContext(ctx, "ALTER SESSION SET CURRENT_SCHEMA = "+schema); err != nil {
			t.Fatal(err)
		}
	}
	return conn
}

// judgeWalk runs the change on a session and holds the walk to it
// (walktest.Judge). The target is shop as Bytebase's sync reads it before
// the change.
func judgeWalk(t testing.TB, w *walkSchemas, conn *sql.Conn, change string, complete bool) (found, rejected int) {
	t.Helper()
	e := &oracleEngine{
		w:      w,
		conn:   conn,
		target: review.Target{Schema: walkSync(t, w.adminConn, w.shop)},
	}
	return walktest.Judge(t, e, change, complete)
}

// oracleEngine is the oracle's Oracle side of one change. Oracle commits
// DDL as it runs, so a rejected statement leaves the database as it was
// only through its own atomicity, which Digest checks.
type oracleEngine struct {
	w      *walkSchemas
	conn   *sql.Conn
	target review.Target
}

// Statements splits a change into its statements, each ended so that the
// statements joined again are a change: a SQL statement with its ';', and
// a PL/SQL unit, whose text ends with the ';' after its END, with a '/'
// line. Exec runs each without its end, as the driver takes it.
func (e *oracleEngine) Statements(change string) []string {
	_, stmts, syntax := parseReview(change)
	if syntax != nil {
		panic("syntax: " + syntax.Message + "\n" + change)
	}
	out := make([]string, len(stmts))
	for i, s := range stmts {
		text := change[s.loc.Start:s.loc.End]
		if strings.HasSuffix(text, ";") {
			out[i] = text + "\n/"
		} else {
			out[i] = text + ";"
		}
	}
	return out
}

var oraCode = regexp.MustCompile(`ORA-(\d{5})`)

func (e *oracleEngine) Exec(t testing.TB, stmt string) walktest.Outcome {
	t.Helper()
	ctx, cancel := walkContext()
	defer cancel()
	if text, ok := strings.CutSuffix(stmt, "\n/"); ok {
		stmt = text
	} else {
		stmt = strings.TrimSuffix(stmt, ";")
	}
	_, err := e.conn.ExecContext(ctx, stmt)
	if err == nil {
		return walktest.Outcome{}
	}
	m := oraCode.FindStringSubmatch(err.Error())
	if m == nil || errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("%s: %v", stmt, err)
	}
	return walktest.Outcome{Code: m[1], Message: err.Error()}
}

func (e *oracleEngine) Report(t testing.TB, change string) *walktest.Finding {
	t.Helper()
	if strings.TrimSpace(change) == "" {
		return nil
	}
	result, err := Review(context.Background(), change, review.Options{Rules: []review.Rule{review.WalkThrough}}, []review.Target{e.target})
	if err != nil {
		t.Fatalf("Review: %v", err)
	}
	var out *walktest.Finding
	for _, f := range result.Findings {
		switch f.Rule {
		case review.Syntax:
			t.Fatalf("syntax: %s\n%s", f.Message, change)
		case review.WalkThrough:
			if out != nil {
				t.Errorf("more than one finding: %q and %q", out.Message, f.Message)
			}
			out = &walktest.Finding{Statement: f.Statement, Message: f.Message}
		}
	}
	return out
}

func (e *oracleEngine) Violation(t testing.TB, change string) *walk.Violation {
	t.Helper()
	_, stmts, syntax := parseReview(change)
	if syntax != nil {
		t.Fatalf("syntax: %s", syntax.Message)
	}
	v, err := walkThrough(context.Background(), change, stmts, e.target)
	if err != nil || v == nil {
		t.Fatalf("no violation for a change Review reports on (%v):\n%s", err, change)
	}
	return v
}

// walkCodes maps a finding's message to the Oracle errors it stands for.
var walkCodes = []struct {
	match func(string) bool
	codes []string
}{
	{func(m string) bool { return strings.HasPrefix(m, "Table or view ") }, []string{"00942"}},
	{func(m string) bool { return strings.HasSuffix(m, " is already used by an existing object") }, []string{"00955"}},
	{func(m string) bool { return strings.HasSuffix(m, " is already used by an existing constraint") }, []string{"02264"}},
	{func(m string) bool { return strings.HasPrefix(m, "Must use DROP MATERIALIZED VIEW") }, []string{"12083"}},
	{func(m string) bool { return strings.HasPrefix(m, "A view is not appropriate here") }, []string{"01702"}},
	{func(m string) bool { return strings.HasPrefix(m, "A sequence is not allowed here") }, []string{"02201"}},
	{func(m string) bool { return strings.HasPrefix(m, "A procedure, function, package, or type") }, []string{"04044"}},
	{func(m string) bool { return strings.HasPrefix(m, "Sequence ") }, []string{"02289"}},
	{func(m string) bool { return strings.HasPrefix(m, "Private synonym ") }, []string{"01434"}},
	{func(m string) bool { return strings.HasPrefix(m, "Object ") }, []string{"04043"}},
	{func(m string) bool { return strings.HasPrefix(m, "RENAME cannot rename") }, []string{"03001"}},
	{func(m string) bool { return strings.HasPrefix(m, "Cannot rename materialized view") }, []string{"32318"}},
	{func(m string) bool { return strings.Contains(m, " being added already exists") }, []string{"01430"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "Column ") && strings.Contains(m, " already exists in table ")
	}, []string{"00957"}},
	{func(m string) bool { return strings.HasPrefix(m, "Column ") && strings.Contains(m, " is named twice") }, []string{"00957"}},
	{func(m string) bool { return strings.HasPrefix(m, "Column ") }, []string{"00904"}},
	{func(m string) bool { return strings.HasSuffix(m, " can have only one primary key") }, []string{"02260"}},
	{func(m string) bool { return strings.HasSuffix(m, " has no primary key to drop") }, []string{"02441"}},
	{func(m string) bool { return strings.HasPrefix(m, "Cannot drop constraint") }, []string{"02443"}},
	{func(m string) bool { return strings.HasPrefix(m, "Cannot rename constraint") }, []string{"23292"}},
	{func(m string) bool { return strings.HasPrefix(m, "Cannot enable constraint") }, []string{"02430"}},
	{func(m string) bool { return strings.HasPrefix(m, "Cannot disable constraint") }, []string{"02431"}},
	{func(m string) bool { return strings.HasPrefix(m, "Index ") }, []string{"01418"}},
	{func(m string) bool {
		return strings.HasSuffix(m, " already exists on another table, cannot replace it")
	}, []string{"04095"}},
	{func(m string) bool {
		return strings.HasPrefix(m, "Trigger ") && strings.HasSuffix(m, " already exists")
	}, []string{"04081"}},
	{func(m string) bool { return strings.HasPrefix(m, "Trigger ") }, []string{"04080"}},
	{func(m string) bool { return strings.HasPrefix(m, "The identifier ") }, []string{"00972"}},
}

func (e *oracleEngine) Codes(message string) []string {
	for _, c := range walkCodes {
		if c.match(message) {
			return c.codes
		}
	}
	return nil
}

// Holds checks each fact the violation claims against the data
// dictionary. A claim about the text holds by itself.
func (e *oracleEngine) Holds(t testing.TB, v *walk.Violation, message string) (bool, string) {
	t.Helper()
	for _, c := range walktest.Claims(v) {
		if c.Fact == walk.Unknown {
			continue
		}
		present, kind := e.exists(t, c.Key)
		switch {
		case present != (c.Fact == walk.Present):
			return false, fmt.Sprintf("%s %q in %q is present: %v", c.Key.Space, c.Key.Name, walk.Names(c.Key.Parent), present)
		case present && c.Kind != "" && kind != c.Kind:
			return false, fmt.Sprintf("%s %q in %q is of kind %q, not %q", c.Key.Space, c.Key.Name, walk.Names(c.Key.Parent), kind, c.Kind)
		}
	}
	return true, ""
}

// exists reads whether a key's object exists in the data dictionary, and
// its kind as the walk keys it.
func (e *oracleEngine) exists(t testing.TB, k walk.Key) (bool, string) {
	t.Helper()
	path := walk.Names(k.Parent)
	ctx, cancel := walkContext()
	defer cancel()
	rows := func(q string, args ...any) [][2]string {
		r, err := e.w.admin.QueryContext(ctx, q, args...)
		if err != nil {
			t.Fatalf("%s: %v", q, err)
		}
		defer r.Close()
		var out [][2]string
		for r.Next() {
			var a, b sql.NullString
			if err := r.Scan(&a, &b); err != nil {
				t.Fatal(err)
			}
			out = append(out, [2]string{a.String, b.String})
		}
		return out
	}
	first := func(rs [][2]string) (bool, string) {
		if len(rs) == 0 {
			return false, ""
		}
		return true, rs[0][1]
	}
	owner := path[0]
	switch k.Space {
	case spaceObject:
		rs := rows(`SELECT object_name, object_type FROM dba_objects WHERE owner = :1 AND object_name = :2 AND namespace = 1 ORDER BY CASE object_type WHEN 'MATERIALIZED VIEW' THEN 0 ELSE 1 END`, owner, k.Name)
		return first(rs)
	case spaceTable:
		rs := rows(`SELECT table_name, CASE WHEN EXISTS (SELECT 1 FROM dba_mviews m WHERE m.owner = x.owner AND m.mview_name = x.table_name) THEN 'MATERIALIZED VIEW' ELSE 'TABLE' END FROM (SELECT owner, table_name FROM dba_tables UNION ALL SELECT owner, table_name FROM dba_object_tables) x WHERE owner = :1 AND table_name = :2`, owner, k.Name)
		return first(rs)
	case spaceView:
		ok, _ := first(rows(`SELECT view_name, 'VIEW' FROM dba_views WHERE owner = :1 AND view_name = :2`, owner, k.Name))
		return ok, kindView
	case spaceSequence:
		ok, _ := first(rows(`SELECT sequence_name, 'SEQUENCE' FROM dba_sequences WHERE sequence_owner = :1 AND sequence_name = :2`, owner, k.Name))
		return ok, kindSequence
	case spaceColumn:
		return first(rows(`SELECT column_name, '' FROM dba_tab_cols WHERE owner = :1 AND table_name = :2 AND column_name = :3`, owner, path[1], k.Name))
	case spacePrimary:
		return first(rows(`SELECT constraint_name, '' FROM dba_constraints WHERE owner = :1 AND table_name = :2 AND constraint_type = 'P'`, owner, path[1]))
	case spaceTableConstraint:
		return first(rows(`SELECT constraint_name, '' FROM dba_constraints WHERE owner = :1 AND table_name = :2 AND constraint_name = :3 UNION ALL SELECT log_group_name, '' FROM dba_log_groups WHERE owner = :4 AND table_name = :5 AND log_group_name = :6`, owner, path[1], k.Name, owner, path[1], k.Name))
	case spaceConstraint:
		rs := rows(`SELECT owner, table_name FROM dba_constraints WHERE owner = :1 AND constraint_name = :2 UNION ALL SELECT owner, table_name FROM dba_log_groups WHERE owner = :3 AND log_group_name = :4`, owner, k.Name, owner, k.Name)
		if len(rs) == 0 {
			return false, ""
		}
		return true, tablePath(rs[0][0], rs[0][1])
	case spaceIndex:
		rs := rows(`SELECT table_owner, table_name FROM dba_indexes WHERE owner = :1 AND index_name = :2`, owner, k.Name)
		if len(rs) == 0 {
			return false, ""
		}
		return true, tablePath(rs[0][0], rs[0][1])
	case spaceTrigger:
		rs := rows(`SELECT table_owner, table_name FROM dba_triggers WHERE owner = :1 AND trigger_name = :2`, owner, k.Name)
		if len(rs) == 0 {
			return false, ""
		}
		return true, tablePath(rs[0][0], rs[0][1])
	}
	t.Fatalf("unknown space %s", k.Space)
	return false, ""
}

// Digest reads the names of both schemas' objects, columns, and
// constraints, log groups among them, from the dictionary's own tables,
// which answer in a few milliseconds where its views take hundreds, so a
// rejected statement that changed one fails the case. It leaves out the
// columns of a generated name: a rejected ADD of a column with a NOT NULL
// constraint of a taken name leaves the column behind as an unused hidden
// one of a generated name (verified), which no claim names. An object of
// type 10 is a placeholder for a name a dependency missed.
func (e *oracleEngine) Digest(t testing.TB) string {
	t.Helper()
	ctx, cancel := walkContext()
	defer cancel()
	in := e.w.users
	q := `SELECT 'O' || o.owner# || '.' || o.name || '.' || o.type# FROM sys.obj$ o WHERE o.owner# IN ` + in + ` AND o.type# <> 10
UNION ALL SELECT 'C' || o.owner# || '.' || o.name || '.' || c.name FROM sys.col$ c JOIN sys.obj$ o ON o.obj# = c.obj# WHERE o.owner# IN ` + in + ` AND c.name NOT LIKE 'SYS\_%' ESCAPE '\'
UNION ALL SELECT 'K' || c.owner# || '.' || c.name FROM sys.con$ c WHERE c.owner# IN ` + in + `
ORDER BY 1`
	rows, err := e.w.admin.QueryContext(ctx, q)
	if err != nil {
		t.Fatalf("digest: %v", err)
	}
	defer rows.Close()
	var b strings.Builder
	for rows.Next() {
		var line sql.NullString
		if err := rows.Scan(&line); err != nil {
			t.Fatal(err)
		}
		b.WriteString(line.String)
		b.WriteByte('\n')
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

// walkSync reads a schema into the metadata Bytebase's Oracle sync makes
// of it, for what WalkThrough reads: the tables, object tables among them,
// that are not materialized views' containers, with every column the user
// made, invisible ones too, and their primary and unique constraints as
// indexes, their foreign keys, and their check constraints not named
// SYS_C; the indexes of the schema's tables and materialized views that
// back no constraint of their name, are not named SYS_, and hold a column
// the user made; the triggers of the schema on its own tables and views;
// and its views, materialized views, sequences but identity ones, and
// functions, procedures, and packages.
//
// bytebase: backend/plugin/db/oracle/sync.go — getTables, getTableColumns, getConstraints, getIndexesAndConstraints, getTriggers, getViews, getMaterializedViews, getSequences, getRoutines
func walkSync(t testing.TB, conn *sql.Conn, schema string) *metadata.DatabaseSchemaMetadata {
	t.Helper()
	ctx, cancel := walkContext()
	defer cancel()
	query := func(q string, scan func(vals []sql.NullString)) {
		t.Helper()
		q = strings.ReplaceAll(q, ":1", "'"+schema+"'")
		rows, err := conn.QueryContext(ctx, q)
		if err != nil {
			t.Fatalf("%s: %v", q, err)
		}
		defer rows.Close()
		cols, _ := rows.Columns()
		for rows.Next() {
			vals := make([]sql.NullString, len(cols))
			ptrs := make([]any, len(cols))
			for i := range vals {
				ptrs[i] = &vals[i]
			}
			if err := rows.Scan(ptrs...); err != nil {
				t.Fatal(err)
			}
			scan(vals)
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
	}
	s := &metadata.SchemaMetadata{}
	tables := make(map[string]*metadata.TableMetadata)
	query(`SELECT table_name, 0 FROM all_tables WHERE owner = :1 AND table_name NOT IN (SELECT mview_name FROM all_mviews WHERE owner = all_tables.owner)
UNION ALL SELECT table_name, 1 FROM all_object_tables WHERE owner = :1 AND table_name NOT IN (SELECT mview_name FROM all_mviews WHERE owner = all_object_tables.owner)`, func(v []sql.NullString) {
		tables[v[0].String] = &metadata.TableMetadata{Name: v[0].String, SkipDump: v[1].String == "1"}
		s.Tables = append(s.Tables, tables[v[0].String])
	})
	query(`SELECT table_name, column_name FROM all_tab_cols WHERE owner = :1 AND user_generated = 'YES' ORDER BY table_name, column_id NULLS LAST, column_name`, func(v []sql.NullString) {
		if tb, ok := tables[v[0].String]; ok {
			tb.Columns = append(tb.Columns, &metadata.ColumnMetadata{Name: v[1].String})
		}
	})
	mviews := make(map[string]*metadata.MaterializedViewMetadata)
	query(`SELECT mview_name FROM all_mviews WHERE owner = :1`, func(v []sql.NullString) {
		mviews[v[0].String] = &metadata.MaterializedViewMetadata{Name: v[0].String}
		s.MaterializedViews = append(s.MaterializedViews, mviews[v[0].String])
	})
	indexed := make(map[string]bool)
	addIndex := func(table string, idx *metadata.IndexMetadata) {
		if tb, ok := tables[table]; ok {
			tb.Indexes = append(tb.Indexes, idx)
		} else if mv, ok := mviews[table]; ok {
			mv.Indexes = append(mv.Indexes, idx)
		}
	}
	query(`SELECT table_name, constraint_name, constraint_type, search_condition_vc FROM all_constraints WHERE owner = :1`, func(v []sql.NullString) {
		table, name := v[0].String, v[1].String
		switch v[2].String {
		case "P", "U":
			addIndex(table, &metadata.IndexMetadata{Name: name, Primary: v[2].String == "P", Unique: true, IsConstraint: true})
			indexed[table+"."+name] = true
		case "C":
			if tb, ok := tables[table]; ok && !strings.HasPrefix(name, "SYS_C") {
				tb.CheckConstraints = append(tb.CheckConstraints, &metadata.CheckConstraintMetadata{Name: name, Expression: v[3].String})
			}
		case "R":
			if tb, ok := tables[table]; ok {
				tb.ForeignKeys = append(tb.ForeignKeys, &metadata.ForeignKeyMetadata{Name: name})
			}
		}
	})
	// An index holds a column the user made unless every column it lists
	// is a hidden SYS_NC column and it has no expression.
	userColumns := make(map[string]bool)
	query(`SELECT c.index_name, c.column_name, (SELECT COUNT(*) FROM all_ind_expressions e WHERE e.index_owner = c.index_owner AND e.index_name = c.index_name) FROM all_ind_columns c WHERE c.table_owner = :1`, func(v []sql.NullString) {
		if !strings.HasPrefix(v[1].String, "SYS_NC") || v[2].String != "0" {
			userColumns[v[0].String] = true
		}
	})
	query(`SELECT table_name, index_name FROM all_indexes WHERE owner = :1`, func(v []sql.NullString) {
		table, name := v[0].String, v[1].String
		if indexed[table+"."+name] || strings.HasPrefix(name, "SYS_") || !userColumns[name] {
			return
		}
		addIndex(table, &metadata.IndexMetadata{Name: name})
	})
	views := make(map[string]*metadata.ViewMetadata)
	query(`SELECT view_name FROM all_views WHERE owner = :1`, func(v []sql.NullString) {
		views[v[0].String] = &metadata.ViewMetadata{Name: v[0].String}
		s.Views = append(s.Views, views[v[0].String])
	})
	query(`SELECT trigger_name, table_name, base_object_type FROM all_triggers WHERE owner = :1 AND table_owner = owner AND table_name IS NOT NULL AND base_object_type IN ('TABLE', 'VIEW')`, func(v []sql.NullString) {
		tr := &metadata.TriggerMetadata{Name: v[0].String}
		switch {
		case v[2].String == "VIEW":
			if vw, ok := views[v[1].String]; ok {
				vw.Triggers = append(vw.Triggers, tr)
			}
		case tables[v[1].String] != nil:
			tables[v[1].String].Triggers = append(tables[v[1].String].Triggers, tr)
		case mviews[v[1].String] != nil:
			mviews[v[1].String].Triggers = append(mviews[v[1].String].Triggers, tr)
		}
	})
	query(`SELECT sequence_name FROM all_sequences WHERE sequence_owner = :1 AND sequence_name NOT LIKE 'ISEQ$$\_%' ESCAPE '\'`, func(v []sql.NullString) {
		s.Sequences = append(s.Sequences, &metadata.SequenceMetadata{Name: v[0].String})
	})
	query(`SELECT DISTINCT name, type FROM all_source WHERE owner = :1 AND type IN ('FUNCTION', 'PROCEDURE', 'PACKAGE')`, func(v []sql.NullString) {
		switch v[1].String {
		case "FUNCTION":
			s.Functions = append(s.Functions, &metadata.FunctionMetadata{Name: v[0].String})
		case "PROCEDURE":
			s.Procedures = append(s.Procedures, &metadata.ProcedureMetadata{Name: v[0].String})
		case "PACKAGE":
			s.Packages = append(s.Packages, &metadata.PackageMetadata{Name: v[0].String})
		}
	})
	return &metadata.DatabaseSchemaMetadata{Name: schema, Schemas: []*metadata.SchemaMetadata{s}}
}

// walkSweepPerBase is how many changes of the sweep run on one base.
const walkSweepPerBase = 10

// TestWalkThroughSweepOracle runs random changes over a small pool of
// names that collide often, quoted lower-case ones included, on bases of
// random statements. Precision must hold on every one; recall is logged.
// WALK_SWEEP_N sets the number of changes, WALK_SWEEP_SEED the seed.
//
// The sync is the costly part, a second or so of dictionary views, so each
// base is synced once and serves walkSweepPerBase changes, each on the base
// made again: made again, it is the same but for the names Oracle
// generates, which the walk never claims. Two workers run at once, each on
// schemas of its own.
func TestWalkThroughSweepOracle(t *testing.T) {
	n, seed, _ := walktest.SweepParams(t, 40, 20261028, "")
	admin := startOracle(t)
	const workers = 2
	for wk := range workers {
		t.Run(fmt.Sprintf("worker%d", wk), func(t *testing.T) {
			t.Parallel()
			shop, other := fmt.Sprintf("SHOP%d", wk+1), fmt.Sprintf("OTHER%d", wk+1)
			rng := rand.New(rand.NewSource(seed + int64(wk)))
			changes := (n + workers - 1 - wk) / workers
			var found, rejected int
			for i := 0; i < changes; i += walkSweepPerBase {
				base := walkSweepChange(rng, shop, other, 6+rng.Intn(8))
				var target *metadata.DatabaseSchemaMetadata
				for j := i; j < min(changes, i+walkSweepPerBase); j++ {
					change := walkSweepChange(rng, shop, other, 10+rng.Intn(12))
					t.Run(fmt.Sprintf("change%d", j), func(t *testing.T) {
						w := resetWalk(t, admin, shop, other, nil)
						e := &oracleEngine{w: w, conn: w.adminConn}
						for _, q := range e.Statements(base) {
							e.Exec(t, q)
						}
						// The base may have moved the session.
						w.exec(t, w.adminConn, "ALTER SESSION SET CURRENT_SCHEMA = "+shop)
						if target == nil {
							target = walkSync(t, w.adminConn, shop)
						}
						e.target = review.Target{Schema: target}
						f, r := walktest.Judge(t, e, change, false)
						found += f
						rejected += r
						if t.Failed() {
							t.Logf("base:\n%s\nchange:\n%s", base, change)
						}
					})
				}
			}
			t.Logf("%s seed %d: found %d of the %d statements the server rejects", shop, seed+int64(wk), found, rejected)
		})
	}
}

// walkSweepChange writes n random statements on schema shop, naming other
// at times.
func walkSweepChange(rng *rand.Rand, shop, other string, n int) string {
	pick := func(xs ...string) string { return xs[rng.Intn(len(xs))] }
	maybe := func(s string) string {
		if rng.Intn(3) == 0 {
			return s
		}
		return ""
	}
	name := func() string { return pick("a", "b", "c", `"a"`, "v", "s") }
	tbl := func() string { return pick("", "", "", "", strings.ToLower(other)+".") + name() }
	col := func() string { return pick("x", "y", "z", `"x"`) }
	con := func() string { return pick("ck", "pk", "uq", "fk", "a") }
	idx := func() string { return pick("i", "j", "a", "ck") }
	trg := func() string { return pick("tr", "tr2", "a") }
	var b strings.Builder
	for range n {
		var stmt string
		plsql := false
		switch rng.Intn(34) {
		case 0, 1:
			stmt = fmt.Sprintf("CREATE TABLE %s%s (id NUMBER%s PRIMARY KEY, x NUMBER%s, y NUMBER)", maybe("IF NOT EXISTS "), tbl(), maybe(" CONSTRAINT "+con()), pick("", " CONSTRAINT "+con()+" CHECK (x > 0)", " CONSTRAINT "+con()+" NOT NULL"))
		case 2:
			stmt = fmt.Sprintf("CREATE TABLE %s (x NUMBER, y NUMBER, CONSTRAINT %s UNIQUE (y))", tbl(), con())
		case 3:
			stmt = fmt.Sprintf("CREATE TABLE %s AS SELECT 1 x, 2 y FROM dual", tbl())
		case 4, 5:
			stmt = fmt.Sprintf("DROP TABLE %s%s%s", maybe("IF EXISTS "), tbl(), maybe(" CASCADE CONSTRAINTS"))
		case 6:
			stmt = fmt.Sprintf("ALTER TABLE %s ADD (%s NUMBER%s)", tbl(), col(), pick("", " CONSTRAINT "+con()+" NOT NULL", " CONSTRAINT "+con()+" CHECK ("+col()+" > 0)"))
		case 7:
			stmt = fmt.Sprintf("ALTER TABLE %s %s", tbl(), pick("DROP COLUMN "+col(), "DROP ("+col()+", "+col()+")", "SET UNUSED COLUMN "+col()))
		case 8:
			stmt = fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s", tbl(), col(), col())
		case 9:
			stmt = fmt.Sprintf("ALTER TABLE %s MODIFY (%s %s)", tbl(), col(), pick("NUMBER(10)", "NULL", "NOT NULL"))
		case 10, 11:
			stmt = fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s %s", tbl(), con(), pick("CHECK (x > 0)", "UNIQUE (y)", "PRIMARY KEY (x)", "FOREIGN KEY (x) REFERENCES a (id)"))
		case 12:
			stmt = fmt.Sprintf("ALTER TABLE %s %s", tbl(), pick("DROP CONSTRAINT "+con(), "DROP PRIMARY KEY", "DROP UNIQUE (y)", "DROP CONSTRAINT "+con()+" CASCADE"))
		case 13:
			stmt = fmt.Sprintf("ALTER TABLE %s RENAME CONSTRAINT %s TO %s", tbl(), con(), con())
		case 14:
			stmt = fmt.Sprintf("ALTER TABLE %s RENAME TO %s", tbl(), name())
		case 15:
			stmt = fmt.Sprintf("ALTER TABLE %s %s CONSTRAINT %s", tbl(), pick("ENABLE", "DISABLE"), con())
		case 16, 17:
			stmt = fmt.Sprintf("CREATE %sINDEX %s%s ON %s (%s)", maybe("UNIQUE "), maybe("IF NOT EXISTS "), idx(), tbl(), col())
		case 18:
			stmt = fmt.Sprintf("DROP INDEX %s%s", maybe("IF EXISTS "), idx())
		case 19:
			stmt = fmt.Sprintf("ALTER INDEX %s RENAME TO %s", idx(), idx())
		case 20:
			stmt = fmt.Sprintf("CREATE %sVIEW %s AS SELECT 1 x FROM dual%s", maybe("OR REPLACE "), tbl(), maybe(" WITH READ ONLY CONSTRAINT "+con()))
		case 21:
			stmt = fmt.Sprintf("DROP VIEW %s%s", maybe("IF EXISTS "), tbl())
		case 22:
			stmt = pick("CREATE SEQUENCE "+maybe("IF NOT EXISTS ")+tbl(), "DROP SEQUENCE "+maybe("IF EXISTS ")+tbl(), "ALTER SEQUENCE "+tbl()+" INCREMENT BY 2")
		case 23:
			stmt = pick("CREATE "+maybe("OR REPLACE ")+"SYNONYM "+tbl()+" FOR a", "DROP SYNONYM "+tbl())
		case 24:
			stmt = "TRUNCATE TABLE " + tbl()
		case 25:
			stmt, plsql = fmt.Sprintf("CREATE %sTRIGGER %s BEFORE INSERT ON %s FOR EACH ROW BEGIN NULL; END;", maybe("OR REPLACE "), trg(), tbl()), true
		case 26:
			stmt = pick("DROP TRIGGER "+maybe("IF EXISTS ")+trg(), "ALTER TRIGGER "+trg()+" RENAME TO "+trg(), "ALTER TRIGGER "+trg()+" DISABLE")
		case 27:
			stmt = fmt.Sprintf("RENAME %s TO %s", name(), name())
		case 28:
			stmt = fmt.Sprintf("INSERT INTO %s (x) VALUES (1)", tbl())
		case 29:
			stmt = "ALTER SESSION SET CURRENT_SCHEMA = " + strings.ToLower(pick(shop, shop, other))
		case 30:
			stmt, plsql = fmt.Sprintf("CREATE OR REPLACE PROCEDURE %s AS BEGIN NULL; END;", tbl()), true
		case 31:
			stmt = pick("DROP PROCEDURE "+tbl(), "DROP FUNCTION "+tbl(), "COMMENT ON TABLE "+tbl()+" IS 'x'")
		case 32:
			stmt, plsql = fmt.Sprintf("CREATE %sTYPE %s AS OBJECT (x NUMBER);", maybe("OR REPLACE "), tbl()), true
		case 33:
			stmt = fmt.Sprintf("TRUNCATE TABLE %s", tbl())
		}
		b.WriteString(stmt)
		if plsql {
			b.WriteString("\n/\n")
		} else {
			b.WriteString(";\n")
		}
	}
	return b.String()
}
