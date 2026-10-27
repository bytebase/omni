package review

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"testing"

	mssql "github.com/microsoft/go-mssqldb"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/review/walk/walktest"
)

// walkBase is the shop database the curated cases start from: tables with
// a key, a default, a check, a unique constraint, a foreign key, and an
// index; a table with none; a view, a procedure, a function, a sequence,
// a synonym; schema s with a table; and a user whose default schema is s.
var walkBase = []string{
	"CREATE TABLE dbo.r (id int NOT NULL CONSTRAINT pk_r PRIMARY KEY)",
	"CREATE TABLE dbo.t (id int NOT NULL CONSTRAINT pk_t PRIMARY KEY, a int CONSTRAINT df_a DEFAULT 0, b int CONSTRAINT ck_b CHECK (b > 0), c int CONSTRAINT uq_c UNIQUE, d int CONSTRAINT fk_r REFERENCES dbo.r (id))",
	"CREATE INDEX ix_a ON dbo.t (a)",
	"CREATE TABLE dbo.nokey (a int)",
	"CREATE VIEW dbo.v AS SELECT id FROM dbo.r",
	"CREATE PROCEDURE dbo.pr AS SELECT 1",
	"CREATE FUNCTION dbo.fn() RETURNS int AS BEGIN RETURN 1 END",
	"CREATE SEQUENCE dbo.sq",
	"CREATE SYNONYM dbo.syn FOR dbo.t",
	"CREATE SCHEMA s",
	"CREATE TABLE s.st (a int)",
	"CREATE TABLE s.t (id int)",
	"CREATE USER u_s WITHOUT LOGIN WITH DEFAULT_SCHEMA = s",
	"GRANT CONTROL ON DATABASE::%s TO u_s",
}

// walkCase is a change on walkBase, after setup when it has one. With
// user, the change runs as that database user, whose default schema the
// target gives as defaultSchema.
type walkCase struct {
	name          string
	setup         string
	change        string
	complete      bool
	user          string
	defaultSchema string
}

var walkCases = []walkCase{
	// Creations.
	{name: "creations over existing names", change: "CREATE TABLE dbo.t (a int);\nCREATE TABLE dbo.pk_t (a int);\nCREATE TABLE dbo.v (a int);\nCREATE TABLE dbo.pr (a int);\nCREATE TABLE dbo.sq (a int);\nCREATE TABLE dbo.df_a (a int);\nCREATE TABLE t (a int);\nCREATE TABLE nosuch.x (a int);\nCREATE TABLE s.t (a int);\nCREATE TABLE s.t2 (a int);\nCREATE TABLE dbo.x (a int);\nCREATE TABLE dbo.x (a int);", complete: true},
	{name: "objects share one space", change: "CREATE SEQUENCE dbo.t;\nCREATE SYNONYM dbo.pr FOR dbo.r;\nCREATE SEQUENCE dbo.ck_b;\nCREATE SEQUENCE nosuch.q;\nCREATE SEQUENCE dbo.q2;\nCREATE SYNONYM dbo.q2 FOR dbo.r;", complete: true},
	{name: "a new table's text", change: "CREATE TABLE dbo.x (a int, a int);\nCREATE TABLE dbo.x (a int PRIMARY KEY, b int PRIMARY KEY);\nCREATE TABLE dbo.x (a int CONSTRAINT c1 CHECK (a > 0), b int CONSTRAINT c1 CHECK (b > 0));\nCREATE TABLE dbo.x (a int CONSTRAINT x CHECK (a > 0));\nCREATE TABLE dbo.x (a int CONSTRAINT ck_b CHECK (a > 0));\nCREATE TABLE dbo.x (a int, INDEX i1 (a), INDEX i1 (a));\nCREATE TABLE dbo.x (a int CONSTRAINT u1 UNIQUE, INDEX u1 (a));", complete: true},
	{name: "foreign key targets", change: "CREATE TABLE dbo.x (a int REFERENCES dbo.nosuch (id));\nCREATE TABLE dbo.x (a int REFERENCES dbo.v (id));\nCREATE TABLE dbo.x (a int REFERENCES dbo.syn (id));\nCREATE TABLE dbo.x (id int NOT NULL PRIMARY KEY, p int REFERENCES dbo.x (id));\nALTER TABLE dbo.nokey ADD CONSTRAINT fk2 FOREIGN KEY (a) REFERENCES nosuch.r (id);\nALTER TABLE dbo.nokey ADD CONSTRAINT fk2 FOREIGN KEY (a) REFERENCES dbo.r (id);", complete: true},
	{name: "a new table holds what it declares", change: "CREATE TABLE dbo.x (id int NOT NULL CONSTRAINT pk_x PRIMARY KEY, a int CONSTRAINT df_x DEFAULT 0, b int UNIQUE, INDEX ix_x (a));\nALTER TABLE dbo.x ADD a int;\nALTER TABLE dbo.x DROP COLUMN zz;\nALTER TABLE dbo.x ADD CONSTRAINT pk2 PRIMARY KEY (id);\nCREATE INDEX ix_x ON dbo.x (b);\nCREATE INDEX pk_x ON dbo.x (b);\nALTER TABLE dbo.x DROP CONSTRAINT df_x;\nALTER TABLE dbo.x DROP CONSTRAINT df_x;\nCREATE TABLE dbo.df_x (a int);\nCREATE TABLE dbo.pk_x (a int);", complete: true},
	{name: "temporary tables", change: "CREATE TABLE #x (a int);\nCREATE TABLE #x (a int);\nALTER TABLE #x ADD a int;\nALTER TABLE #nosuch ADD b int;\nCREATE INDEX i ON #x (a);\nCREATE INDEX i ON #x (a);\nDROP INDEX nosuch ON #x;\nINSERT INTO #nosuch VALUES (1);\nDROP TABLE #x;\nDROP TABLE #x;\nSELECT 1 AS a INTO #y;\nSELECT 1 AS a INTO #y;\nTRUNCATE TABLE #nosuch;", complete: true},
	{name: "temporary tables and GO", change: "CREATE TABLE #x (a int);\nGO\nCREATE TABLE #x (a int);\nDROP TABLE #x;\nDROP TABLE #x;"},
	{name: "a temporary table's foreign key is skipped", change: "CREATE TABLE #x (a int CONSTRAINT pkx PRIMARY KEY, b int CONSTRAINT fkx REFERENCES dbo.nosuch (id));\nCREATE TABLE #y (a int CONSTRAINT pk_t PRIMARY KEY);", complete: true},
	{name: "names too long", change: "CREATE TABLE dbo." + strings.Repeat("a", 129) + " (a int);\nDROP TABLE IF EXISTS dbo." + strings.Repeat("b", 129) + ";\nCREATE TABLE #" + strings.Repeat("c", 116) + " (a int);\nCREATE TABLE #" + strings.Repeat("d", 115) + " (a int);\nCREATE TABLE dbo." + strings.Repeat("e", 128) + " (a int);", complete: true},
	{name: "SELECT INTO", change: "SELECT 1 AS a INTO dbo.t;\nSELECT 1 AS a INTO nosuch.n;\nSELECT 1 AS a INTO dbo.n;\nSELECT 1 AS a INTO dbo.n;\nALTER TABLE dbo.n ADD b int;\nCREATE INDEX i ON dbo.n (a);\nCREATE INDEX i ON dbo.n (a);\nDROP TABLE dbo.n;\nINSERT INTO dbo.n VALUES (1);", complete: true},
	// Drops.
	{name: "DROP TABLE", change: "DROP TABLE dbo.nosuch;\nDROP TABLE dbo.v;\nDROP TABLE dbo.syn;\nDROP TABLE dbo.pr;\nDROP TABLE dbo.sq;\nDROP TABLE dbo.pk_t;\nDROP TABLE IF EXISTS dbo.v;\nDROP TABLE IF EXISTS dbo.nosuch;\nDROP TABLE nosuch.t;\nDROP TABLE dbo.nokey;\nDROP TABLE dbo.nokey;\nDROP TABLE dbo.NOKEY;\nINSERT INTO dbo.nokey VALUES (1);\nDROP TABLE dbo.r, dbo.r;", complete: true},
	{name: "DROP TABLE runs its items in order", change: "DROP TABLE s.st, dbo.nosuch;\nDROP TABLE dbo.nosuch, s.st;\nDROP TABLE s.st, dbo.v;\nDROP TABLE s.st;\nDROP TABLE s.st;", complete: true},
	{name: "a dropped table frees its names", change: "DROP TABLE dbo.t;\nCREATE TABLE dbo.pk_t (a int);\nCREATE TABLE dbo.ix_a (a int);\nCREATE TABLE dbo.t (a int);\nCREATE TABLE dbo.t (a int);"},
	{name: "DROP of other kinds", change: "DROP VIEW dbo.t;\nDROP VIEW dbo.nosuch;\nDROP VIEW IF EXISTS dbo.t;\nDROP PROCEDURE dbo.t;\nDROP PROCEDURE dbo.fn;\nDROP PROCEDURE IF EXISTS dbo.t;\nDROP FUNCTION dbo.pr;\nDROP SEQUENCE dbo.nosuch;\nDROP SEQUENCE dbo.t;\nDROP SEQUENCE IF EXISTS dbo.t;\nDROP SYNONYM dbo.t;\nDROP SYNONYM IF EXISTS dbo.t;\nDROP TRIGGER dbo.t;\nDROP TRIGGER dbo.pk_t;\nDROP VIEW dbo.v;\nDROP VIEW dbo.v;\nDROP PROCEDURE dbo.pr;\nDROP PROCEDURE dbo.pr;\nDROP SYNONYM dbo.syn;\nDROP SYNONYM dbo.syn;\nDROP SEQUENCE dbo.sq;\nDROP SEQUENCE dbo.sq;", complete: true},
	{name: "DROP of unsynced names", change: "DROP PROCEDURE dbo.nosuch;\nDROP FUNCTION dbo.nosuch;\nDROP SYNONYM dbo.nosuch;\nDROP TRIGGER dbo.nosuch;\nDROP TYPE dbo.nosuch;"},
	{name: "schemas", change: "CREATE SCHEMA s;\nGO\nDROP SCHEMA nosuch;\nDROP SCHEMA IF EXISTS nosuch;\nDROP SCHEMA s;\nDROP TABLE s.st;\nDROP TABLE s.t;\nDROP SCHEMA s;\nDROP SCHEMA s;\nCREATE TABLE s.x (a int);\nGO\nCREATE SCHEMA s;\nGO\nCREATE TABLE s.x (a int);\nCREATE TABLE s.x (a int);\nDROP SCHEMA s;", complete: true},
	{name: "a schema a type holds", setup: "CREATE TYPE s.ty FROM int", change: "DROP TABLE s.st;\nDROP TABLE s.t;\nDROP SCHEMA s;"},
	// ALTER TABLE.
	{name: "ALTER TABLE", change: "ALTER TABLE dbo.nosuch ADD x int;\nALTER TABLE dbo.v ADD x int;\nALTER TABLE dbo.syn ADD x int;\nALTER TABLE dbo.pr ADD x int;\nALTER TABLE dbo.t ADD a int;\nALTER TABLE dbo.t ADD x int, x int;\nALTER TABLE dbo.t DROP COLUMN nosuch;\nALTER TABLE dbo.t DROP COLUMN IF EXISTS nosuch;\nALTER TABLE dbo.t ALTER COLUMN nosuch int;\nALTER TABLE dbo.t ADD CONSTRAINT pk2 PRIMARY KEY (id);\nALTER TABLE dbo.t ADD y int NOT NULL CONSTRAINT pk3 PRIMARY KEY;\nALTER TABLE dbo.nokey ADD CONSTRAINT pk_t CHECK (a > 0);\nALTER TABLE dbo.nokey ADD CONSTRAINT t CHECK (a > 0);\nALTER TABLE dbo.t ADD CONSTRAINT ix_a UNIQUE (b);\nALTER TABLE dbo.t DROP CONSTRAINT nosuch;\nALTER TABLE dbo.r DROP CONSTRAINT pk_t;\nALTER TABLE dbo.t DROP CONSTRAINT ix_a;\nALTER TABLE dbo.t DROP CONSTRAINT IF EXISTS nosuch;\nALTER TABLE dbo.t NOCHECK CONSTRAINT nosuch;\nALTER TABLE dbo.t ADD CONSTRAINT c1 CHECK (a > 0), CONSTRAINT c1 CHECK (b > 0);", complete: true},
	{name: "ALTER TABLE effects", change: "ALTER TABLE dbo.t DROP CONSTRAINT ck_b;\nALTER TABLE dbo.t DROP CONSTRAINT ck_b;\nALTER TABLE dbo.nokey ADD CONSTRAINT ck_b CHECK (a > 0);\nALTER TABLE dbo.nokey ADD z int;\nALTER TABLE dbo.nokey ADD z int;\nALTER TABLE dbo.nokey DROP COLUMN z;\nALTER TABLE dbo.nokey DROP COLUMN z;\nALTER TABLE dbo.t DROP CONSTRAINT df_a;\nALTER TABLE dbo.t DROP COLUMN a, COLUMN a;\nALTER TABLE dbo.t DROP CONSTRAINT ix_a;\nALTER TABLE dbo.nokey ADD CONSTRAINT pk_n PRIMARY KEY (z);\nALTER TABLE dbo.nokey ALTER COLUMN a int NOT NULL;\nALTER TABLE dbo.nokey ADD CONSTRAINT pk_n PRIMARY KEY (a);\nALTER TABLE dbo.nokey ADD CONSTRAINT pk_m PRIMARY KEY (a);\nALTER TABLE dbo.nokey DROP CONSTRAINT pk_n;\nALTER TABLE dbo.nokey ADD CONSTRAINT pk_m PRIMARY KEY (a);"},
	{name: "a list's constraints see its columns", change: "ALTER TABLE dbo.nokey ADD CONSTRAINT ckx CHECK (x > 0), x int;\nALTER TABLE dbo.nokey DROP CONSTRAINT ckx;\nALTER TABLE dbo.nokey DROP COLUMN x;\nALTER TABLE dbo.nokey DROP COLUMN x;", complete: true},
	{name: "TRUNCATE", change: "TRUNCATE TABLE dbo.nosuch;\nTRUNCATE TABLE dbo.v;\nTRUNCATE TABLE dbo.syn;\nTRUNCATE TABLE dbo.nokey;", complete: true},
	// Indexes.
	{name: "indexes", change: "CREATE INDEX ix_a ON dbo.t (b);\nCREATE INDEX pk_t ON dbo.t (b);\nCREATE INDEX pk_t ON dbo.r (id);\nCREATE INDEX i ON dbo.nosuch (a);\nCREATE INDEX i ON dbo.syn (a);\nCREATE INDEX nosuch ON dbo.t (a) WITH (DROP_EXISTING = ON);\nDROP INDEX nosuch ON dbo.t;\nDROP INDEX ix_a ON dbo.nosuch;\nDROP INDEX pk_t ON dbo.t;\nDROP INDEX IF EXISTS uq_c ON dbo.t;\nDROP INDEX IF EXISTS nosuch ON dbo.t;\nALTER INDEX nosuch ON dbo.t REBUILD;\nALTER INDEX ix_a ON dbo.nosuch REBUILD;\nDROP INDEX dbo.t.ix_a;\nDROP INDEX dbo.t.ix_a;\nCREATE INDEX ix_a ON dbo.t (a);\nDROP INDEX ix_a ON dbo.t, ix_a ON dbo.t;", complete: true},
	{name: "statistics share the index namespace", change: "CREATE STATISTICS st1 ON dbo.t (b);\nCREATE INDEX st1 ON dbo.t (b);"},
	// Views, routines, triggers, each beginning its batch.
	{name: "views", change: "CREATE VIEW dbo.t AS SELECT 1 AS a;\nGO\nCREATE VIEW dbo.v AS SELECT 1 AS a;\nGO\nALTER VIEW dbo.nosuch AS SELECT 1 AS a;\nGO\nALTER VIEW dbo.t AS SELECT 1 AS a;\nGO\nCREATE OR ALTER VIEW dbo.t AS SELECT 1 AS a;\nGO\nCREATE OR ALTER VIEW dbo.w AS SELECT 1 AS a;\nGO\nCREATE VIEW nosuch.w AS SELECT 1 AS a;\nGO\nCREATE VIEW dbo.w AS SELECT 1 AS a;\nGO\nDROP VIEW dbo.w;\nDROP VIEW dbo.w;", complete: true},
	{name: "routines", change: "CREATE PROCEDURE dbo.t AS SELECT 1;\nGO\nCREATE PROCEDURE dbo.pr AS SELECT 1;\nGO\nALTER PROCEDURE dbo.fn AS SELECT 1;\nGO\nCREATE OR ALTER PROCEDURE dbo.fn AS SELECT 1;\nGO\nCREATE OR ALTER PROCEDURE dbo.t AS SELECT 1;\nGO\nCREATE FUNCTION dbo.t() RETURNS int AS BEGIN RETURN 1 END;\nGO\nCREATE OR ALTER FUNCTION dbo.pr() RETURNS int AS BEGIN RETURN 1 END;\nGO\nALTER FUNCTION dbo.t() RETURNS int AS BEGIN RETURN 1 END;\nGO\nCREATE PROCEDURE nosuch.p AS SELECT 1;\nGO\nCREATE PROCEDURE dbo.p2 AS SELECT 1;\nGO\nALTER PROCEDURE dbo.p2 AS SELECT 2;\nGO\nDROP PROCEDURE dbo.p2;\nGO\nALTER PROCEDURE dbo.p2 AS SELECT 2;", complete: true},
	{name: "triggers", change: "CREATE TRIGGER dbo.t ON dbo.r AFTER INSERT AS SELECT 1;\nGO\nCREATE TRIGGER dbo.tr2 ON dbo.nosuch AFTER INSERT AS SELECT 1;\nGO\nCREATE TRIGGER dbo.tr2 ON dbo.syn AFTER INSERT AS SELECT 1;\nGO\nCREATE TRIGGER s.tr3 ON dbo.r AFTER INSERT AS SELECT 1;\nGO\nALTER TRIGGER dbo.t ON dbo.r AFTER INSERT AS SELECT 1;\nGO\nALTER TRIGGER dbo.pr ON dbo.r AFTER INSERT AS SELECT 1;\nGO\nCREATE OR ALTER TRIGGER dbo.t ON dbo.r AFTER INSERT AS SELECT 1;\nGO\nCREATE TRIGGER dbo.tr ON dbo.nokey AFTER INSERT AS SELECT 1;\nGO\nCREATE TRIGGER dbo.tr ON dbo.r AFTER INSERT AS SELECT 1;\nGO\nDROP TRIGGER dbo.tr;\nDROP TRIGGER dbo.tr;", complete: true},
	{name: "sequences and types", change: "ALTER SEQUENCE dbo.nosuch RESTART;\nALTER SEQUENCE dbo.t RESTART;\nALTER SEQUENCE dbo.sq RESTART;\nCREATE TYPE dbo.ty FROM int;\nCREATE TYPE dbo.ty FROM int;\nCREATE TYPE nosuch.ty FROM int;\nCREATE TABLE dbo.ty (a int);\nCREATE TYPE dbo.tt AS TABLE (a int);\nDROP TYPE dbo.tt;\nDROP TYPE dbo.tt;\nDROP TYPE IF EXISTS dbo.tt;", complete: true},
	// sp_rename.
	{name: "sp_rename", change: "EXEC sp_rename 'dbo.t', 'r';\nEXEC sp_rename 'dbo.t', 'pk_t';\nEXEC sp_rename 'dbo.t.a', 'b', 'COLUMN';\nEXEC sp_rename 'dbo.t.nosuch', 'z', 'COLUMN';\nEXEC sp_rename 'dbo.t.nosuch', 'ix_b', 'INDEX';\nEXEC sp_rename 'dbo.t.ix_a', 'pk_t', 'INDEX';\nEXEC sp_rename 'dbo.nokey', 'n2';\nALTER TABLE dbo.n2 ADD a int;\nALTER TABLE dbo.nokey ADD z int;\nCREATE TABLE dbo.nokey (a int);\nEXEC sp_rename 'dbo.t.a', 'z', 'COLUMN';\nALTER TABLE dbo.t DROP COLUMN a;\nALTER TABLE dbo.t ADD z int;\nEXEC sp_rename 'dbo.t.ix_a', 'ix_b', 'INDEX';\nDROP INDEX ix_a ON dbo.t;\nDROP INDEX ix_b ON dbo.t;\nDROP INDEX ix_b ON dbo.t;\nEXEC sp_rename @objname = N'dbo.v', @newname = N'v2', @objtype = N'OBJECT';\nDROP VIEW dbo.v;", complete: true},
	{name: "sp_rename keeps a table's children", change: "EXEC sp_rename 'dbo.t', 't2';\nALTER TABLE dbo.t2 DROP CONSTRAINT ck_b;\nALTER TABLE dbo.t2 DROP CONSTRAINT ck_b;\nCREATE TABLE dbo.ck_b (a int);\nDROP INDEX ix_a ON dbo.t2;\nDROP INDEX ix_a ON dbo.t2;\nEXEC sp_rename 'dbo.t2', 'T2';\nALTER TABLE dbo.T2 ADD id int;", complete: true},
	// DML.
	{name: "DML targets", change: "INSERT INTO dbo.nosuch VALUES (1);\nINSERT INTO dbo.syn (id) VALUES (5);\nDROP TABLE dbo.nokey;\nINSERT INTO dbo.nokey VALUES (1);"},
	{name: "DML on dropped names", change: "DROP TABLE dbo.nokey;\nUPDATE dbo.nokey SET a = 1;\nDELETE FROM dbo.nokey;\nMERGE dbo.nokey AS x USING dbo.r ON x.a = r.id WHEN MATCHED THEN DELETE;", complete: true},
	{name: "DML on a CTE or an alias", change: "WITH c AS (SELECT * FROM dbo.nokey) UPDATE c SET a = 1;\nUPDATE x SET a = 1 FROM dbo.nokey AS x;\nDELETE x FROM dbo.nokey x;\nCREATE TABLE dbo.t (a int);"},
	{name: "a trigger may run DDL", change: "CREATE TRIGGER dbo.mk ON dbo.nokey AFTER INSERT AS CREATE TABLE dbo.made (a int);\nGO\nINSERT INTO dbo.nokey VALUES (1);\nDROP TABLE dbo.made;\nDROP TABLE dbo.made;"},
	{name: "a cascade fires another table's trigger", setup: "CREATE TABLE dbo.par (id int NOT NULL PRIMARY KEY)\nCREATE TABLE dbo.chi (id int REFERENCES dbo.par (id) ON DELETE CASCADE)\nINSERT dbo.par VALUES (1)\nINSERT dbo.chi VALUES (1)\nCREATE TRIGGER dbo.chi_tr ON dbo.chi AFTER DELETE AS CREATE TABLE dbo.made2 (a int)", change: "DELETE FROM dbo.par;\nCREATE TABLE dbo.made2 (a int);"},
	{name: "an enabled trigger fires", setup: "CREATE TRIGGER dbo.mk ON dbo.nokey AFTER INSERT AS CREATE TABLE dbo.made (a int)\nDISABLE TRIGGER dbo.mk ON dbo.nokey", change: "ALTER TABLE dbo.nokey ENABLE TRIGGER ALL;\nINSERT INTO dbo.nokey VALUES (1);\nDROP TABLE dbo.made;"},
	// Batches and the session.
	{name: "a view after GO", change: "CREATE TABLE dbo.x (a int);\nGO\nCREATE VIEW dbo.w AS SELECT 1 AS a;"},
	{name: "USE", change: "USE master;\nCREATE TABLE dbo.t (a int);\nUSE shop;\nCREATE TABLE dbo.t (a int);\nCREATE TABLE shop.dbo.t (a int);", complete: true},
	{name: "USE and GO", change: "USE master;\nGO\nCREATE TABLE dbo.zz (a int);\nUSE shop;\nGO\nCREATE TABLE dbo.t (a int);", complete: true},
	{name: "IF", change: "IF OBJECT_ID('dbo.t') IS NOT NULL DROP TABLE dbo.t;\nCREATE TABLE dbo.t (a int);\nIF 1 = 0 DROP TABLE dbo.nokey;\nCREATE TABLE dbo.nokey (a int);\nIF 1 = 1 BEGIN CREATE TABLE dbo.x (a int); CREATE INDEX i ON dbo.x (a) END ELSE DROP TABLE dbo.r;\nCREATE TABLE dbo.x (a int);\nCREATE TABLE dbo.r (a int);"},
	{name: "the default schema", user: "u_s", defaultSchema: "s", change: "CREATE TABLE t (a int);\nCREATE TABLE st (a int);\nCREATE TABLE newt (a int);\nCREATE TABLE newt (a int);\nDROP TABLE nokey;\nDROP TABLE nokey;\nDROP TABLE t;\nDROP TABLE t;\nALTER TABLE r ADD x int;"},
	{name: "no default schema", defaultSchema: "-", change: "CREATE TABLE t (a int);\nDROP TABLE nokey;\nDROP TABLE nokey;\nCREATE TABLE dbo.t (a int);\nDROP TABLE nosuch;\nALTER TABLE nokey ADD a int;\nDROP TABLE dbo.nokey;"},
	{name: "unknown by design", change: "DROP TABLE IF EXISTS dbo.syn;\nALTER TABLE dbo.t ADD A int;\nCREATE INDEX IX_A ON dbo.t (b);\nEXEC sp_rename 'dbo.nosuch', 'x';\nEXEC sp_rename 'dbo.nokey', 'dbo.n3';\nDROP TABLE dbo.n3;\nCREATE TABLE master.dbo.r (a int);\nCREATE TABLE master.dbo.r (a int);\nDROP PROCEDURE dbo.nosuch;\nCREATE TABLE dbo.T (a int);\nINSERT INTO dbo.nosuch VALUES (1);\nDROP TABLE guest.nosuch;\nCREATE TABLE guest.x (a int);\nCREATE TABLE guest.x (a int);\nALTER TABLE dbo.t DROP COLUMN d;"},
}

var walkCollations = []string{"SQL_Latin1_General_CP1_CI_AS", "Latin1_General_100_CS_AS"}

func TestWalkThroughAgainstMSSQL(t *testing.T) {
	db := startMSSQL(t)
	for _, collation := range walkCollations {
		t.Run(collation, func(t *testing.T) {
			w := resetWalk(t, db, "shop", collation)
			var found, rejected int
			for _, tt := range walkCases {
				t.Run(tt.name, func(t *testing.T) {
					target := review.Target{DefaultSchema: "dbo"}
					switch tt.defaultSchema {
					case "-":
						target.DefaultSchema = ""
					case "":
					default:
						target.DefaultSchema = tt.defaultSchema
					}
					var setup []string
					if tt.setup != "" {
						setup = strings.Split(tt.setup, "\n")
					}
					f, r := judgeWalk(t, w, target, setup, tt.user, tt.change, tt.complete)
					found += f
					rejected += r
				})
			}
			t.Logf("%s: found %d of the %d statements the server rejects", collation, found, rejected)
		})
	}
}

// walkDB is one database the oracle runs changes on.
type walkDB struct {
	db        *sql.DB
	name      string
	collation string
}

// resetWalk recreates a database with a collation, from walkBase.
func resetWalk(t testing.TB, db *sql.DB, name, collation string) walkDB {
	t.Helper()
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	stmts := []string{
		"USE master",
		fmt.Sprintf("IF DB_ID('%[1]s') IS NOT NULL BEGIN ALTER DATABASE %[1]s SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE %[1]s; END", name),
		"CREATE DATABASE " + name + " COLLATE " + collation,
		"USE " + name,
	}
	for _, q := range walkBase {
		if strings.Contains(q, "%s") {
			q = fmt.Sprintf(q, name)
		}
		stmts = append(stmts, q)
	}
	for _, q := range stmts {
		if _, err := conn.ExecContext(ctx, q); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
	}
	return walkDB{db: db, name: name, collation: collation}
}

// judgeWalk runs the change on a session of its own, in one transaction
// as Bytebase's transaction mode runs it, after setup, as user when given,
// and holds the walk to it (walktest.Judge). The target is the database
// as the sync reads it before the change. The transaction is rolled back
// at the end, which keeps walkBase for the next case.
func judgeWalk(t *testing.T, w walkDB, target review.Target, setup []string, user, change string, complete bool) (found, rejected int) {
	t.Helper()
	e := newWalkEngine(t, w, user)
	for _, q := range setup {
		if strings.TrimSpace(q) == "" {
			continue
		}
		if o := e.Exec(t, q); o.Code != "" {
			t.Fatalf("setup %s: %s %s", q, o.Code, o.Message)
		}
	}
	e.prefix = e.accepted
	e.accepted = nil
	target.Schema = walkSync(t, e.conn, w.name)
	e.target = target
	return walktest.Judge(t, e, change, complete)
}

// walkEngine is the oracle's SQL Server side of one change, on one
// connection in one transaction. A rejected statement is undone by
// rolling back to the savepoint before it; one that ends the transaction
// (2714, 1913, 3728, verified) is undone by replaying the statements
// accepted so far in a new one.
type walkEngine struct {
	walkDB
	conn     *sql.Conn
	user     string
	target   review.Target
	prefix   []string
	accepted []string
	// last is the statement run before, and how it ended, for a failure
	// that shows only later.
	last string
}

func newWalkEngine(t testing.TB, w walkDB, user string) *walkEngine {
	t.Helper()
	conn, err := w.db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = conn.ExecContext(context.Background(), "IF @@TRANCOUNT > 0 ROLLBACK")
		// The session is discarded, not pooled: one that ran as another
		// user does not survive the pool's reset.
		_ = conn.Raw(func(any) error { return driver.ErrBadConn })
		_ = conn.Close()
	})
	e := &walkEngine{walkDB: w, conn: conn, user: user}
	e.begin(t)
	return e
}

// begin starts the transaction on the shop database, as the user.
func (e *walkEngine) begin(t testing.TB) {
	t.Helper()
	stmts := []string{"IF @@TRANCOUNT > 0 ROLLBACK", "REVERT", "USE " + e.name, "BEGIN TRANSACTION"}
	if e.user != "" {
		stmts = append(stmts, "EXECUTE AS USER = '"+e.user+"'")
	}
	for _, q := range stmts {
		if _, err := e.conn.ExecContext(context.Background(), q); err != nil && q != "REVERT" {
			t.Fatalf("%s: %v", q, err)
		}
	}
}

func (e *walkEngine) Statements(change string) []string {
	_, stmts, syntax := parseReview(change)
	if syntax != nil {
		panic("syntax: " + syntax.Message)
	}
	out := make([]string, len(stmts))
	for i, s := range stmts {
		text := change[s.loc.Start:s.loc.End]
		if _, isGo := s.node.(*ast.GoStmt); !isGo {
			text += ";"
		}
		out[i] = text
	}
	return out
}

var goText = regexp.MustCompile(`(?i)^GO(\s+\d+)?$`)

func isGoText(s string) bool { return goText.MatchString(strings.TrimSpace(s)) }

func (e *walkEngine) Exec(t testing.TB, stmt string) walktest.Outcome {
	t.Helper()
	if isGoText(stmt) {
		e.accepted = append(e.accepted, stmt)
		return walktest.Outcome{}
	}
	ctx := context.Background()
	if _, err := e.conn.ExecContext(ctx, "SAVE TRANSACTION walk_sp"); err != nil {
		t.Fatalf("SAVE TRANSACTION before %q after %q: %v", stmt, e.last, err)
	}
	e.last = stmt
	_, err := e.conn.ExecContext(ctx, stmt)
	if err == nil {
		var n int
		if err := e.conn.QueryRowContext(ctx, "SELECT @@TRANCOUNT").Scan(&n); err != nil {
			t.Fatal(err)
		}
		if n == 0 {
			t.Fatalf("%s ended the transaction", stmt)
		}
		e.accepted = append(e.accepted, stmt)
		return walktest.Outcome{}
	}
	var me mssql.Error
	if !errors.As(err, &me) {
		t.Fatalf("%s: %v", stmt, err)
	}
	first := me
	if len(me.All) > 0 {
		first = me.All[0]
	}
	var state int
	if err := e.conn.QueryRowContext(ctx, "SELECT XACT_STATE()").Scan(&state); err != nil {
		t.Fatal(err)
	}
	e.last = fmt.Sprintf("%s (%d, xact %d)", stmt, first.Number, state)
	if state == 1 {
		if _, err := e.conn.ExecContext(ctx, "ROLLBACK TRANSACTION walk_sp"); err != nil {
			t.Fatalf("ROLLBACK TRANSACTION walk_sp: %v", err)
		}
	} else {
		e.replay(t)
	}
	return walktest.Outcome{Code: strconv.Itoa(int(first.Number)), Message: first.Message}
}

// replay redoes the setup and the accepted statements in a new
// transaction, after a rejection ended the old one.
func (e *walkEngine) replay(t testing.TB) {
	t.Helper()
	e.begin(t)
	for _, q := range append(append([]string(nil), e.prefix...), e.accepted...) {
		if isGoText(q) {
			continue
		}
		if _, err := e.conn.ExecContext(context.Background(), q); err != nil {
			t.Fatalf("replay %s: %v", q, err)
		}
	}
}

func (e *walkEngine) Report(t testing.TB, change string) *walktest.Finding {
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

func (e *walkEngine) Violation(t testing.TB, change string) *walk.Violation {
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

// walkCodes maps a finding's message to the SQL Server errors it stands
// for, the most specific first.
var walkCodes = []struct {
	re    *regexp.Regexp
	codes []string
}{
	{regexp.MustCompile(`^foreign key target .* does not exist$`), []string{"1767"}},
	{regexp.MustCompile(`^foreign key target .* is not a table$`), []string{"1768"}},
	{regexp.MustCompile(`^schema .* already exists$`), []string{"2714"}},
	{regexp.MustCompile(`^schema .* still holds objects$`), []string{"3729"}},
	{regexp.MustCompile(`^schema .* does not exist$`), []string{"2760", "15151"}},
	{regexp.MustCompile(`^an object named .* already exists$`), []string{"2714", "15335"}},
	{regexp.MustCompile(`^table .* already has a primary key$`), []string{"1779"}},
	{regexp.MustCompile(`^table .* does not exist$`), []string{"3701", "4902", "4701", "1088", "8197", "15248"}},
	{regexp.MustCompile(`^(view|procedure|function|trigger|synonym) .* does not exist$`), []string{"3701", "208"}},
	{regexp.MustCompile(`^sequence .* does not exist$`), []string{"3701", "15151"}},
	{regexp.MustCompile(`^type .* does not exist$`), []string{"218"}},
	{regexp.MustCompile(`^type .* already exists$`), []string{"219"}},
	{regexp.MustCompile(`^object .* does not exist$`), []string{"208", "15225"}},
	{regexp.MustCompile(`^column .* already exists in `), []string{"2705", "15335"}},
	{regexp.MustCompile(`^column .* does not exist in `), []string{"4924", "15248"}},
	{regexp.MustCompile(`^column .* is named twice`), []string{"2705"}},
	{regexp.MustCompile(`^the name .* is given twice`), []string{"8168"}},
	{regexp.MustCompile(`^the statement adds two primary keys`), []string{"8110"}},
	{regexp.MustCompile(`^index .* of .* enforces a constraint$`), []string{"3723"}},
	{regexp.MustCompile(`^index .* already exists on `), []string{"1913", "15335"}},
	{regexp.MustCompile(`^index .* does not exist on `), []string{"3701", "7999", "2727", "15248"}},
	{regexp.MustCompile(`^index .* is dropped twice$`), []string{"3701"}},
	{regexp.MustCompile(` is not a constraint of `), []string{"3728", "3733", "4917"}},
	{regexp.MustCompile(`^trigger .* must be in the schema of its table`), []string{"2103"}},
	{regexp.MustCompile(`^the identifier .* is longer than `), []string{"103", "193"}},
	{regexp.MustCompile(` must be the first statement in its batch$`), []string{"111"}},
	{regexp.MustCompile(` is dropped twice$`), []string{"3701"}},
	{regexp.MustCompile(` is not a table or view$`), []string{"1914", "8197", "3701"}},
	{regexp.MustCompile(` is not a (table|view|procedure|function|trigger|synonym|sequence)$`), []string{"3705", "4909", "4708", "2010", "15151"}},
}

func (e *walkEngine) Codes(message string) []string {
	for _, c := range walkCodes {
		if c.re.MatchString(message) {
			return c.codes
		}
	}
	return nil
}

// Holds checks each fact the violation claims against the catalog, read
// on the change's own session, inside its transaction.
func (e *walkEngine) Holds(t testing.TB, v *walk.Violation, message string) (bool, string) {
	t.Helper()
	for _, c := range walktest.Claims(v) {
		if c.Fact == walk.Unknown {
			continue
		}
		objs, ok := e.catalog(t, c.Key)
		if !ok {
			return false, fmt.Sprintf("cannot read the parent of %+v", c.Key)
		}
		collation := e.collation
		if c.Key.Space == spaceTemp || strings.HasPrefix(c.Key.Parent, tempSchema) {
			collation = "DATABASE_DEFAULT"
		}
		equal := func(a, b string) bool {
			var n int
			q := "SELECT CASE WHEN CAST(@p1 AS nvarchar(128)) COLLATE " + collation + " = CAST(@p2 AS nvarchar(128)) THEN 1 ELSE 0 END"
			if collation == "DATABASE_DEFAULT" {
				q = "SELECT CASE WHEN CAST(@p1 AS nvarchar(128)) = CAST(@p2 AS nvarchar(128)) COLLATE " + serverCollation + " THEN 1 ELSE 0 END"
			}
			if err := e.conn.QueryRowContext(context.Background(), q, a, b).Scan(&n); err != nil {
				t.Fatal(err)
			}
			return n == 1
		}
		if ok, why := claimHolds(c, objs, equal); !ok {
			return false, why
		}
	}
	return true, ""
}

// catObject is one name of the catalog in a space under a parent.
type catObject struct {
	name string
	kind string
}

// claimHolds checks one claim against the names of its space and parent.
// An exact key claims an object of exactly its spelling (Present), or
// none the collation equates to it (Absent), which the walk claims only
// when no name of its class exists; a class key claims names of the
// class: none, exactly one of its spelling and kind, or some.
func claimHolds(c walktest.Claim, objs []catObject, equal func(a, b string) bool) (bool, string) {
	class, spelling, exact := strings.Cut(c.Key.Name, exactSep)
	var same []catObject
	for _, o := range objs {
		if n, ok := keyed(o.name); ok && n.class == class {
			same = append(same, o)
		}
	}
	switch {
	case exact && c.Fact == walk.Present:
		for _, o := range same {
			if o.name == spelling && (c.Kind == "" || o.kind == c.Kind) {
				return true, ""
			}
		}
		return false, fmt.Sprintf("no %s %q of kind %q in %q: %+v", c.Key.Space, spelling, c.Kind, walk.Names(c.Key.Parent), same)
	case exact:
		for _, o := range same {
			if equal(o.name, spelling) {
				return false, fmt.Sprintf("%s %q exists in %q: %+v", c.Key.Space, spelling, walk.Names(c.Key.Parent), same)
			}
		}
		return true, ""
	case c.Fact == walk.Absent:
		if len(same) > 0 {
			return false, fmt.Sprintf("%s class %q holds %+v", c.Key.Space, class, same)
		}
		return true, ""
	}
	if kind, sp, ok := single(c.Kind); ok {
		if len(same) != 1 || same[0].name != sp || kind != "" && same[0].kind != kind {
			return false, fmt.Sprintf("%s class %q holds %+v, not one %s %q", c.Key.Space, class, same, kind, sp)
		}
		return true, ""
	}
	if len(same) == 0 {
		return false, fmt.Sprintf("%s class %q is empty", c.Key.Space, class)
	}
	return true, ""
}

// objectKinds maps sys.objects types to the walk's kinds.
var objectKinds = map[string]string{
	"U": kindTable, "V": kindView, "SO": kindSequence, "P": kindProcedure, "FN": kindFunction, "IF": kindFunction, "TF": kindFunction,
	"TR": kindTrigger, "SN": kindSynonym, "PK": kindPK, "UQ": kindUQ, "F": kindFK, "C": kindCK, "D": kindDF,
}

// catalog reads the names of a key's space under its parent, in shop and
// the session's tempdb; false when the parent's path names no single
// schema or table.
func (e *walkEngine) catalog(t testing.TB, k walk.Key) ([]catObject, bool) {
	t.Helper()
	path := walk.Names(k.Parent)
	query := func(q string, args ...any) []catObject {
		t.Helper()
		q = strings.ReplaceAll(q, "shop.", e.name+".")
		rows, err := e.conn.QueryContext(context.Background(), q, args...)
		if err != nil {
			t.Fatalf("%s: %v", q, err)
		}
		defer rows.Close()
		var out []catObject
		for rows.Next() {
			var o catObject
			if err := rows.Scan(&o.name, &o.kind); err != nil {
				t.Fatal(err)
			}
			o.kind = strings.TrimSpace(o.kind)
			out = append(out, o)
		}
		return out
	}
	// one returns the id of the one catalog name of a class.
	one := func(names []catObject, class string) (string, bool) {
		id := ""
		for _, o := range names {
			if n, ok := keyed(o.name); ok && n.class == class {
				if id != "" {
					return "", false
				}
				id = o.kind
			}
		}
		return id, id != ""
	}
	schemaID := func(class string) (string, bool) {
		return one(query("SELECT name, CAST(schema_id AS varchar) FROM shop.sys.schemas"), class)
	}
	switch k.Space {
	case spaceSchema:
		return query("SELECT name, '' FROM shop.sys.schemas"), true
	case spaceTemp:
		return e.tempTables(t), true
	}
	if len(path) == 2 && path[0] == tempSchema {
		var id string
		for _, tt := range e.tempTables(t) {
			if n, ok := keyed(strings.TrimPrefix(tt.name, "#")); ok && n.class == path[1] {
				id = tt.kind
			}
		}
		if id == "" {
			return nil, false
		}
		return e.children(t, k.Space, "tempdb", id, query), true
	}
	if len(path) == 0 {
		return nil, false
	}
	sid, ok := schemaID(path[0])
	if !ok {
		// No schema of the class holds anything; several make the
		// parent ambiguous.
		return nil, sid == "" && !e.classShared(t, query("SELECT name, '' FROM shop.sys.schemas"), path[0])
	}
	switch k.Space {
	case spaceObject, spaceRelation:
		objs := query("SELECT name, type FROM shop.sys.objects WHERE schema_id = @p1", sid)
		var out []catObject
		for _, o := range objs {
			kind := objectKinds[o.kind]
			if kind == "" {
				kind = "other"
			}
			if k.Space == spaceRelation && !isRelation(kind) {
				continue
			}
			out = append(out, catObject{o.name, kind})
		}
		return out, true
	case spaceType:
		return query("SELECT name, '' FROM shop.sys.types WHERE is_user_defined = 1 AND schema_id = @p1", sid), true
	}
	if len(path) != 2 {
		return nil, false
	}
	tables := query("SELECT name, CAST(object_id AS varchar) FROM shop.sys.tables WHERE schema_id = @p1", sid)
	tid, ok := one(tables, path[1])
	if !ok {
		return nil, !e.classShared(t, tables, path[1])
	}
	return e.children(t, k.Space, e.name, tid, query), true
}

// children reads a table's columns, indexes, constraints, or primary key.
func (e *walkEngine) children(t testing.TB, space, dbName, id string, query func(string, ...any) []catObject) []catObject {
	switch space {
	case spaceColumn:
		return query("SELECT name, '' FROM "+dbName+".sys.columns WHERE object_id = @p1", id)
	case spaceIndex:
		return query("SELECT name, CASE WHEN is_primary_key = 1 THEN 'pk' WHEN is_unique_constraint = 1 THEN 'uq' ELSE 'index' END FROM "+dbName+".sys.indexes WHERE object_id = @p1 AND name IS NOT NULL", id)
	case spaceConstraint:
		objs := query("SELECT name, type FROM "+dbName+".sys.objects WHERE parent_object_id = @p1 AND type IN ('PK', 'UQ', 'F', 'C', 'D', 'EC')", id)
		for i := range objs {
			objs[i].kind = objectKinds[objs[i].kind]
		}
		return objs
	case spacePrimary:
		return query("SELECT '"+primaryName+"', '' FROM "+dbName+".sys.key_constraints WHERE parent_object_id = @p1 AND type = 'PK'", id)
	}
	t.Fatalf("unknown space %s", space)
	return nil
}

// serverCollation is the container's, which tempdb's names compare under.
const serverCollation = "SQL_Latin1_General_CP1_CI_AS"

// classShared reports whether several names of the catalog share a class,
// which leaves a parent the walk keys by the class ambiguous.
func (e *walkEngine) classShared(t testing.TB, names []catObject, class string) bool {
	n := 0
	for _, o := range names {
		if k, ok := keyed(o.name); ok && k.class == class {
			n++
		}
	}
	return n > 1
}

// tempTables lists the session's temporary tables by the names they were
// created with, each with its object id: tempdb pads a name with
// underscores and a suffix, and OBJECT_ID finds the session's own.
func (e *walkEngine) tempTables(t testing.TB) []catObject {
	t.Helper()
	rows, err := e.conn.QueryContext(context.Background(), "SELECT name, object_id FROM tempdb.sys.tables WHERE name LIKE '#%' AND name NOT LIKE '##%'")
	if err != nil {
		t.Fatal(err)
	}
	type row struct {
		name string
		id   int64
	}
	var all []row
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.name, &r.id); err != nil {
			t.Fatal(err)
		}
		all = append(all, r)
	}
	rows.Close()
	var out []catObject
	for _, r := range all {
		name := r.name
		if i := strings.Index(name, "___"); i > 0 {
			name = name[:i]
		}
		var id sql.NullInt64
		if err := e.conn.QueryRowContext(context.Background(), "SELECT OBJECT_ID('tempdb..'+@p1)", name).Scan(&id); err != nil {
			t.Fatal(err)
		}
		if id.Valid && id.Int64 == r.id {
			out = append(out, catObject{name, strconv.FormatInt(r.id, 10)})
		}
	}
	return out
}

// walkSync reads a database, on the change's session, into the metadata
// Bytebase's MSSQL sync produces for what WalkThrough reads: the schemas
// it lists, their tables with columns (and their default constraints'
// names), indexes, keys, foreign keys, checks, and enabled triggers, and
// views, sequences, functions, and procedures.
//
// bytebase: backend/plugin/db/mssql/sync.go — SyncDBSchema
func walkSync(t testing.TB, conn *sql.Conn, dbName string) *metadata.DatabaseSchemaMetadata {
	t.Helper()
	ctx := context.Background()
	query := func(q string, scan func(*sql.Rows) error) {
		t.Helper()
		q = strings.ReplaceAll(q, "shop.", dbName+".")
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
	db := &metadata.DatabaseSchemaMetadata{Name: dbName}
	schemas := make(map[string]*metadata.SchemaMetadata)
	query("SELECT schema_name FROM shop.INFORMATION_SCHEMA.SCHEMATA WHERE schema_name NOT IN ('db_owner', 'db_accessadmin', 'db_backupoperator', 'db_datareader', 'db_datawriter', 'db_ddladmin', 'db_denydatareader', 'db_denydatawriter', 'db_securityadmin', 'guest', 'INFORMATION_SCHEMA', 'sys') ORDER BY schema_name", func(r *sql.Rows) error {
		s := &metadata.SchemaMetadata{}
		if err := r.Scan(&s.Name); err != nil {
			return err
		}
		schemas[s.Name] = s
		db.Schemas = append(db.Schemas, s)
		return nil
	})
	tables := make(map[string]*metadata.TableMetadata)
	views := make(map[string]*metadata.ViewMetadata)
	key := func(schema, name string) string { return schema + "\x00" + name }
	query("SELECT s.name, t.name FROM shop.sys.tables t JOIN shop.sys.schemas s ON s.schema_id = t.schema_id", func(r *sql.Rows) error {
		var schema, name string
		if err := r.Scan(&schema, &name); err != nil {
			return err
		}
		if s, ok := schemas[schema]; ok {
			tbl := &metadata.TableMetadata{Name: name}
			s.Tables = append(s.Tables, tbl)
			tables[key(schema, name)] = tbl
		}
		return nil
	})
	query("SELECT s.name, v.name FROM shop.sys.views v JOIN shop.sys.sql_modules m ON m.object_id = v.object_id JOIN shop.sys.schemas s ON s.schema_id = v.schema_id", func(r *sql.Rows) error {
		var schema, name string
		if err := r.Scan(&schema, &name); err != nil {
			return err
		}
		if s, ok := schemas[schema]; ok {
			view := &metadata.ViewMetadata{Name: name}
			s.Views = append(s.Views, view)
			views[key(schema, name)] = view
		}
		return nil
	})
	query(`SELECT s.name, o.name, c.name, ISNULL(dc.name, '') FROM shop.sys.columns c
		JOIN shop.sys.objects o ON o.object_id = c.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
		LEFT JOIN shop.sys.default_constraints dc ON dc.object_id = c.default_object_id
		ORDER BY s.name, c.object_id, c.column_id`, func(r *sql.Rows) error {
		var schema, tbl, col, df string
		if err := r.Scan(&schema, &tbl, &col, &df); err != nil {
			return err
		}
		if t, ok := tables[key(schema, tbl)]; ok {
			t.Columns = append(t.Columns, &metadata.ColumnMetadata{Name: col, DefaultConstraintName: df})
		}
		return nil
	})
	// getKeys, then getIndexes: indexes with a key column.
	query(`SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint FROM shop.sys.indexes i
		JOIN shop.sys.objects o ON o.object_id = i.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
		WHERE i.index_id > 0 AND i.name IS NOT NULL AND o.type IN ('U', 'S', 'V')
		AND (i.is_primary_key = 1 OR i.is_unique_constraint = 1 OR EXISTS (SELECT 1 FROM shop.sys.index_columns ic JOIN shop.sys.all_columns c ON c.object_id = ic.object_id AND c.column_id = ic.column_id WHERE ic.object_id = i.object_id AND ic.index_id = i.index_id))`, func(r *sql.Rows) error {
		var schema, tbl, name string
		var primary, unique, constraint bool
		if err := r.Scan(&schema, &tbl, &name, &primary, &unique, &constraint); err != nil {
			return err
		}
		if t, ok := tables[key(schema, tbl)]; ok {
			t.Indexes = append(t.Indexes, &metadata.IndexMetadata{Name: name, Primary: primary, Unique: unique || primary, IsConstraint: primary || constraint})
		}
		return nil
	})
	query(`SELECT s.name, o.name, f.name FROM shop.sys.foreign_keys f JOIN shop.sys.objects o ON o.object_id = f.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
		WHERE EXISTS (SELECT 1 FROM shop.sys.foreign_key_columns fkc WHERE fkc.constraint_object_id = f.object_id)`, func(r *sql.Rows) error {
		var schema, tbl, name string
		if err := r.Scan(&schema, &tbl, &name); err != nil {
			return err
		}
		if t, ok := tables[key(schema, tbl)]; ok {
			t.ForeignKeys = append(t.ForeignKeys, &metadata.ForeignKeyMetadata{Name: name})
		}
		return nil
	})
	query(`SELECT s.name, o.name, c.name FROM shop.sys.check_constraints c JOIN shop.sys.objects o ON o.object_id = c.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U'`, func(r *sql.Rows) error {
		var schema, tbl, name string
		if err := r.Scan(&schema, &tbl, &name); err != nil {
			return err
		}
		if t, ok := tables[key(schema, tbl)]; ok {
			t.CheckConstraints = append(t.CheckConstraints, &metadata.CheckConstraintMetadata{Name: name})
		}
		return nil
	})
	query(`SELECT ss.name, so.name, so.type, st.name FROM shop.sys.triggers st JOIN shop.sys.sql_modules m ON m.object_id = st.object_id
		JOIN shop.sys.objects so ON so.object_id = st.parent_id JOIN shop.sys.schemas ss ON ss.schema_id = so.schema_id
		WHERE st.is_disabled = 0 AND st.is_ms_shipped = 0 AND st.parent_id <> 0 AND so.type IN ('U', 'V')`, func(r *sql.Rows) error {
		var schema, parent, kind, name string
		if err := r.Scan(&schema, &parent, &kind, &name); err != nil {
			return err
		}
		if strings.TrimSpace(kind) == "V" {
			if v, ok := views[key(schema, parent)]; ok {
				v.Triggers = append(v.Triggers, &metadata.TriggerMetadata{Name: name})
			}
		} else if t, ok := tables[key(schema, parent)]; ok {
			t.Triggers = append(t.Triggers, &metadata.TriggerMetadata{Name: name})
		}
		return nil
	})
	query("SELECT s.name, q.name FROM shop.sys.sequences q JOIN shop.sys.schemas s ON s.schema_id = q.schema_id", func(r *sql.Rows) error {
		var schema, name string
		if err := r.Scan(&schema, &name); err != nil {
			return err
		}
		if s, ok := schemas[schema]; ok {
			s.Sequences = append(s.Sequences, &metadata.SequenceMetadata{Name: name})
		}
		return nil
	})
	query(`SELECT s.name, ao.name, ao.type FROM shop.sys.all_objects ao JOIN shop.sys.all_sql_modules m ON m.object_id = ao.object_id JOIN shop.sys.schemas s ON s.schema_id = ao.schema_id
		WHERE ao.type IN ('FN', 'IF', 'TF', 'P', 'RF') AND ao.is_ms_shipped = 0
		AND NOT EXISTS (SELECT 1 FROM shop.sys.extended_properties ep WHERE ep.major_id = ao.object_id AND ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'microsoft_database_tools_support')`, func(r *sql.Rows) error {
		var schema, name, kind string
		if err := r.Scan(&schema, &name, &kind); err != nil {
			return err
		}
		s, ok := schemas[schema]
		if !ok {
			return nil
		}
		switch strings.TrimSpace(kind) {
		case "P", "RF":
			s.Procedures = append(s.Procedures, &metadata.ProcedureMetadata{Name: name})
		default:
			s.Functions = append(s.Functions, &metadata.FunctionMetadata{Name: name})
		}
		return nil
	})
	return db
}

// TestWalkThroughBatchesAgainstMSSQL runs changes the way Bytebase does,
// batch by batch at its GO lines, in one transaction, and holds the walk
// to what the batches show: a finding needs a failing batch, and a change
// whose batches all run gets none. It covers what statement-by-statement
// judging cannot: the statement that must begin its batch (111), and the
// procedure and trigger bodies that run to the end of theirs.
func TestWalkThroughBatchesAgainstMSSQL(t *testing.T) {
	w := resetWalk(t, startMSSQL(t), "shop", walkCollations[0])
	tests := []struct {
		name  string
		sql   string
		fails bool // the change fails
		found bool // the walk reports it
	}{
		{"a view after a statement", "CREATE TABLE dbo.x1 (a int);\nCREATE VIEW dbo.v1 AS SELECT 1 AS a;", true, true},
		{"a view after SET", "SET NOCOUNT ON;\nCREATE VIEW dbo.v1 AS SELECT 1 AS a;", true, true},
		{"a procedure after a statement", "SELECT 1;\nCREATE PROCEDURE dbo.p1 AS SELECT 1;", true, true},
		{"ALTER PROCEDURE after a statement", "SELECT 1;\nALTER PROCEDURE dbo.pr AS SELECT 2;", true, true},
		{"a function after a statement", "SELECT 1;\nCREATE FUNCTION dbo.f1() RETURNS int AS BEGIN RETURN 1 END;", true, true},
		{"a trigger after a statement", "SELECT 1;\nCREATE TRIGGER dbo.tr1 ON dbo.nokey AFTER INSERT AS SELECT 1;", true, true},
		{"a schema after a statement", "SELECT 1;\nCREATE SCHEMA s9;", true, true},
		{"a view after GO", "CREATE TABLE dbo.x1 (a int);\nGO\nCREATE VIEW dbo.v1 AS SELECT 1 AS a;", false, false},
		{"a view after a comment", "-- note\n/* note */\nCREATE VIEW dbo.v1 AS SELECT 1 AS a;", false, false},
		{"a procedure body swallows the batch", "CREATE PROCEDURE dbo.p1 AS SELECT 1;\nCREATE TABLE dbo.t (a int);\nDROP TABLE dbo.nosuch;", false, false},
		{"a procedure body with BEGIN END swallows the batch", "CREATE PROCEDURE dbo.p1 AS BEGIN SELECT 1 END;\nCREATE TABLE dbo.t (a int);", false, false},
		{"a trigger body swallows the batch", "CREATE TRIGGER dbo.tr1 ON dbo.nokey AFTER UPDATE AS SELECT 1;\nCREATE TABLE dbo.t (a int);", false, false},
		{"a procedure body ends at GO", "CREATE PROCEDURE dbo.p1 AS SELECT 1;\nGO\nCREATE TABLE dbo.t (a int);", true, true},
		{"a statement after a view", "CREATE VIEW dbo.v1 AS SELECT 1 AS a;\nCREATE TABLE dbo.t (a int);", true, false},
		{"a GO with a count", "CREATE TABLE dbo.g1 (a int)\nGO 2\nCREATE TABLE dbo.t (a int);", true, false},
		{"a GO with a comment is no GO", "SELECT 1\nGO -- note\nCREATE TABLE dbo.t (a int);", true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newWalkEngine(t, w, "")
			target := review.Target{Schema: walkSync(t, e.conn, w.name), DefaultSchema: "dbo"}
			fails := false
			for _, batch := range bytebaseBatches(tt.sql) {
				for range batch.count {
					if _, err := e.conn.ExecContext(context.Background(), batch.text); err != nil {
						fails = true
					}
				}
				if fails {
					break
				}
			}
			e.target = target
			f := e.Report(t, tt.sql)
			if fails != tt.fails {
				t.Errorf("the change fails: %v, want %v", fails, tt.fails)
			}
			switch {
			case f != nil && !fails:
				t.Errorf("reported %q on statement %d of a change the server runs", f.Message, f.Statement)
			case (f != nil) != tt.found:
				t.Errorf("finding %+v, want one: %v", f, tt.found)
			}
		})
	}
}

// bytebaseBatch is one batch as Bytebase sends it, with its count.
type bytebaseBatch struct {
	text  string
	count int
}

// bytebaseBatches splits a change at the lines Bytebase's batcher takes
// for GO, outside strings and comments, which the cases do not hold.
//
// bytebase: backend/plugin/parser/tsql/batch/batch.go
func bytebaseBatches(sql string) []bytebaseBatch {
	var out []bytebaseBatch
	var cur []string
	for _, line := range strings.Split(sql, "\n") {
		if m := goLine.FindStringSubmatch(line); m != nil {
			n := 1
			if m[1] != "" {
				var err error
				if n, err = strconv.Atoi(m[1]); err != nil || n < 1 {
					cur = append(cur, line)
					continue
				}
			}
			out = append(out, bytebaseBatch{strings.Join(cur, "\n"), n})
			cur = nil
			continue
		}
		cur = append(cur, line)
	}
	if len(cur) > 0 {
		out = append(out, bytebaseBatch{strings.Join(cur, "\n"), 1})
	}
	return out
}

// TestWalkThroughSweepMSSQL runs random changes over a small pool of names
// that collide often, case, collation, and marked variants included, on a
// case-insensitive, a case-sensitive, and a Hungarian database in turn,
// with the default schema known or not. Precision must hold on every one;
// recall is logged. WALK_SWEEP_N sets the number of changes in all, and
// WALK_SWEEP_SEED the seed. The databases run one after another: sessions
// run at once on the server lost their transactions while idle.
func TestWalkThroughSweepMSSQL(t *testing.T) {
	n, seed, _ := walktest.SweepParams(t, 60, 20261027, "")
	db := startMSSQL(t)
	collations := append(walkCollations, "Hungarian_CI_AS")
	for i, collation := range collations {
		t.Run(collation, func(t *testing.T) {
			w := resetWalk(t, db, fmt.Sprintf("shop%d", i), collation)
			// Another database, which a USE selects.
			other := fmt.Sprintf("other%d", i)
			resetOther(t, db, other)
			rng := rand.New(rand.NewSource(seed + int64(i)))
			var found, rejected int
			for j := range (n + len(collations) - 1 - i) / len(collations) {
				base := walkSweepChange(rng, w.name, other, 6+rng.Intn(8))
				change := walkSweepChange(rng, w.name, other, 12+rng.Intn(12))
				target := review.Target{DefaultSchema: "dbo"}
				if rng.Intn(3) == 0 {
					target.DefaultSchema = ""
				}
				t.Run(fmt.Sprintf("change%d", j), func(t *testing.T) {
					e := newWalkEngine(t, w, "")
					// The base's accepted statements are the setup.
					for _, q := range e.Statements(base) {
						e.Exec(t, q)
					}
					e.prefix = append(e.prefix, e.accepted...)
					e.accepted = nil
					q := "USE " + w.name
					if _, err := e.conn.ExecContext(context.Background(), q); err != nil {
						t.Fatal(err)
					}
					e.prefix = append(e.prefix, q)
					// The base's temporary tables end with its session.
					e.dropTemporary(t)
					target.Schema = walkSync(t, e.conn, w.name)
					e.target = target
					f, r := walktest.Judge(t, e, change, false)
					found += f
					rejected += r
					if t.Failed() {
						t.Logf("default schema %q, base:\n%s\nchange:\n%s", target.DefaultSchema, base, change)
					}
				})
			}
			t.Logf("%s seed %d: found %d of the %d statements the server rejects", collation, seed, found, rejected)
		})
	}
}

// resetOther recreates an empty database with schema s.
func resetOther(t testing.TB, db *sql.DB, name string) {
	t.Helper()
	for _, q := range []string{
		fmt.Sprintf("IF DB_ID('%[1]s') IS NOT NULL BEGIN ALTER DATABASE %[1]s SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE %[1]s; END", name),
		"CREATE DATABASE " + name,
		"EXEC('USE " + name + "; EXEC(''CREATE SCHEMA s'')')",
	} {
		if _, err := db.ExecContext(context.Background(), q); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
	}
}

// dropTemporary drops the session's temporary tables, and records the
// drops in the prefix a replay redoes.
func (e *walkEngine) dropTemporary(t testing.TB) {
	t.Helper()
	for _, tt := range e.tempTables(t) {
		q := "DROP TABLE [" + tt.name + "]"
		if _, err := e.conn.ExecContext(context.Background(), q); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
		e.prefix = append(e.prefix, q)
	}
}

func walkSweepChange(rng *rand.Rand, dbName, other string, n int) string {
	pick := func(xs ...string) string { return xs[rng.Intn(len(xs))] }
	maybe := func(s string) string {
		if rng.Intn(3) == 0 {
			return s
		}
		return ""
	}
	schema := func() string { return pick("dbo.", "dbo.", "s.", "", dbName+".dbo.") }
	name := func() string { return pick("a", "b", "A", "ssz", "szsz", "pk_a", "v", "q", "[a b]", "[é]") }
	tbl := func() string { return schema() + name() }
	temp := func() string { return pick("#a", "#A", "#b") }
	col := func() string { return pick("x", "y", "X", "z") }
	idx := func() string { return pick("i", "I", "j", "pk_a") }
	con := func() string { return pick("pk_a", "c", "C", "ssz", "a") }
	var b strings.Builder
	for range n {
		var stmt string
		switch rng.Intn(30) {
		case 0, 1:
			stmt = fmt.Sprintf("CREATE TABLE %s (id int NOT NULL%s, x int%s, y int)", tbl(), maybe(" CONSTRAINT "+con()+" PRIMARY KEY"), maybe(" UNIQUE"))
		case 2:
			stmt = fmt.Sprintf("CREATE TABLE %s (x int, y int CONSTRAINT %s CHECK (y > 0)%s)", tbl(), con(), maybe(", INDEX "+idx()+" (x)"))
		case 3:
			stmt = fmt.Sprintf("CREATE TABLE %s (id int, x int)", temp())
		case 4:
			stmt = fmt.Sprintf("SELECT 1 AS x INTO %s", pick(tbl(), temp()))
		case 5, 6:
			stmt = fmt.Sprintf("DROP TABLE %s%s%s", maybe("IF EXISTS "), pick(tbl(), temp()), maybe(", "+tbl()))
		case 7:
			stmt = fmt.Sprintf("ALTER TABLE %s ADD %s int%s", pick(tbl(), temp()), col(), maybe(", "+col()+" int"))
		case 8:
			stmt = fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s%s", pick(tbl(), temp()), maybe("IF EXISTS "), col())
		case 9:
			stmt = fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s %s", tbl(), con(), pick("CHECK (x > 0)", "UNIQUE (x)", "DEFAULT 0 FOR y"))
		case 10:
			stmt = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s%s", tbl(), maybe("IF EXISTS "), con())
		case 11:
			stmt = fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s bigint", tbl(), col())
		case 12:
			stmt = fmt.Sprintf("CREATE INDEX %s ON %s (x)", idx(), pick(tbl(), temp()))
		case 13:
			stmt = fmt.Sprintf("DROP INDEX %s%s ON %s", maybe("IF EXISTS "), idx(), tbl())
		case 14:
			stmt = fmt.Sprintf("GO\nCREATE %sVIEW %s AS SELECT 1 AS x\nGO", maybe("OR ALTER "), tbl())
		case 15:
			stmt = fmt.Sprintf("DROP VIEW %s%s", maybe("IF EXISTS "), tbl())
		case 16:
			stmt = fmt.Sprintf("CREATE SEQUENCE %s", tbl())
		case 17:
			stmt = fmt.Sprintf("DROP SEQUENCE %s%s", maybe("IF EXISTS "), tbl())
		case 18:
			stmt = fmt.Sprintf("EXEC sp_rename '%s', '%s'", pick("dbo.", "s.")+name(), name())
		case 19:
			stmt = fmt.Sprintf("EXEC sp_rename '%s.%s', '%s', 'COLUMN'", pick("dbo.", "s.")+name(), col(), col())
		case 20:
			stmt = fmt.Sprintf("TRUNCATE TABLE %s", pick(tbl(), temp()))
		case 21:
			stmt = fmt.Sprintf("INSERT INTO %s (x) VALUES (1)", pick(tbl(), temp()))
		case 22:
			stmt = pick("GO\nCREATE SCHEMA s2\nGO", "GO\nCREATE SCHEMA s\nGO", "DROP SCHEMA s", "DROP SCHEMA IF EXISTS s2")
		case 23:
			stmt = pick("GO", "USE "+other, "USE "+dbName, "USE "+strings.ToUpper(dbName))
		case 24:
			stmt = fmt.Sprintf("IF 1 = %d DROP TABLE %s", rng.Intn(2), tbl())
		case 25:
			stmt = fmt.Sprintf("GO\nCREATE PROCEDURE %s AS SELECT 1\nGO", tbl())
		case 26:
			stmt = fmt.Sprintf("DROP PROCEDURE %s%s", maybe("IF EXISTS "), tbl())
		case 27:
			stmt = fmt.Sprintf("CREATE SYNONYM %s FOR dbo.a", tbl())
		case 28:
			stmt = fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s PRIMARY KEY (id)", tbl(), con())
		case 29:
			stmt = fmt.Sprintf("CREATE TYPE %s FROM int", tbl())
		}
		b.WriteString(stmt)
		b.WriteString(";\n")
	}
	return strings.ReplaceAll(b.String(), "GO;", "GO")
}
