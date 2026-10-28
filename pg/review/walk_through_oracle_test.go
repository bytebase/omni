package review

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"testing"
	"unicode"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/review/walk/walktest"
)

// walkBase is the schema the curated WalkThrough cases start from: one of
// every relation kind, keys and constraints, a view reading a column, a
// partitioned table with a partition in another schema, inheritance, and
// an enum type. The partitioned table has an index, whose definition on
// the partition shows the sync the partition is no partitioned table: the
// sync may lose the partitions of one in another schema.
const walkBase = `
CREATE SCHEMA s;
CREATE TABLE p (id integer PRIMARY KEY);
CREATE TABLE t (id serial PRIMARY KEY, code text CONSTRAINT t_code_key UNIQUE, p_id integer CONSTRAINT t_p_fk REFERENCES p (id), n integer CONSTRAINT t_n_check CHECK (n > 0));
CREATE INDEX t_n_idx ON t (n);
CREATE TABLE s.t (id integer, a integer);
CREATE VIEW tv AS SELECT code FROM t;
CREATE MATERIALIZED VIEW mv AS SELECT id FROM p;
CREATE INDEX mv_id_idx ON mv (id);
CREATE SEQUENCE q;
CREATE TYPE mood AS ENUM ('a', 'b');
CREATE TYPE pair AS (x integer, y integer);
CREATE TABLE parted (k integer, v text) PARTITION BY RANGE (k);
CREATE TABLE s.parted_1 PARTITION OF parted FOR VALUES FROM (0) TO (10);
CREATE INDEX parted_k_idx ON parted (k);
CREATE TABLE ip (a integer CONSTRAINT ip_a_check CHECK (a > 0));
CREATE TABLE ic (b integer) INHERITS (ip);
CREATE TABLE nokey (a integer);
`

// walkCase is one curated change. complete asks that every statement the
// server rejects, after the statements it accepts, is found.
type walkCase struct {
	name     string
	change   string
	complete bool
	user     string
	// extra runs after walkBase.
	extra string
	// tenant runs the change in Bytebase's tenant mode, as walkOwner, the
	// database's owner, from a session of postgres.
	tenant bool
}

// walkOwner owns the oracle's database. The oracle runs a change as
// postgres, an admin who does not own the database, as Bytebase does
// outside tenant mode, or, for a tenant case, as walkOwner.
const walkOwner = "walk_owner"

var walkCases = []walkCase{
	// Relations that exist.
	{name: "create a table that exists", change: "CREATE TABLE t (id int);", complete: true},
	{name: "create a table named like an index", change: "CREATE TABLE t_n_idx (a int);", complete: true},
	{name: "create a table named like a view, a sequence, a composite type", change: "CREATE TABLE tv (a int);\nCREATE TABLE q (a int);\nCREATE TABLE pair (a int);", complete: true},
	{name: "create a table named like an enum", change: "CREATE TABLE mood (a int);", complete: true},
	{name: "create a type named like a table", change: "CREATE TYPE t AS ENUM ('x');\nCREATE TYPE q AS ENUM ('x');\nCREATE TYPE t_n_idx AS ENUM ('x');", complete: true},
	{name: "create in another schema", change: "CREATE TABLE s.t (a int);\nCREATE TABLE s.p (a int);\nCREATE TABLE s.p (a int);", complete: true},
	{name: "IF NOT EXISTS", change: "CREATE TABLE IF NOT EXISTS t (a int);\nCREATE TABLE IF NOT EXISTS tv (a int);\nCREATE TABLE IF NOT EXISTS mood (a int);\nCREATE INDEX IF NOT EXISTS t ON p (id);\nCREATE SEQUENCE IF NOT EXISTS t_n_idx;", complete: true},
	{name: "IF NOT EXISTS skips a missing reference", change: "CREATE TABLE IF NOT EXISTS t (a int REFERENCES nosuch);", complete: true},
	{name: "create twice", change: "CREATE TABLE x (a int);\nCREATE INDEX xi ON x (a);\nCREATE TABLE xi (a int);\nCREATE TABLE x (a int);", complete: true},
	{name: "create an index named like a relation", change: "CREATE INDEX q ON t (n);\nCREATE INDEX t_n_idx ON t (n);\nCREATE INDEX s ON t (n);", complete: true},
	{name: "a sequence and a view", change: "CREATE SEQUENCE t;\nCREATE VIEW q AS SELECT 1 AS x;\nCREATE VIEW tv AS SELECT 1 AS x;\nCREATE OR REPLACE VIEW tv AS SELECT code FROM t;", complete: true},
	{name: "a composite type", change: "CREATE TYPE p AS (a int);\nCREATE TYPE mood AS (a int);\nCREATE TYPE pair AS (a int);", complete: true},
	{name: "a domain", change: "CREATE DOMAIN mood AS int;\nCREATE DOMAIN t AS int;\nCREATE DOMAIN dom AS int;\nCREATE DOMAIN dom AS int;", complete: true},
	{name: "IF NOT EXISTS on a sequence leaves its name without a row type", change: "CREATE TABLE a (id serial PRIMARY KEY);\nCREATE TABLE IF NOT EXISTS a_id_seq (x int);\nCREATE TYPE a_id_seq AS ENUM ('a');\nCREATE TABLE IF NOT EXISTS q (x int);\nCREATE TYPE q AS ENUM ('a');"},
	{name: "a materialized view and CREATE TABLE AS", change: "CREATE MATERIALIZED VIEW mv AS SELECT 1 AS x;\nCREATE TABLE tv AS SELECT 1 AS x;\nCREATE TABLE IF NOT EXISTS t AS SELECT 1 AS x;\nSELECT 1 AS x INTO t;", complete: true},

	// Relations that do not exist.
	{name: "alter a missing table", change: "ALTER TABLE nosuch ADD COLUMN a int;\nALTER TABLE IF EXISTS nosuch ADD COLUMN a int;", complete: true},
	{name: "missing schema", change: "ALTER TABLE nosuch.t ADD COLUMN a int;\nCREATE TABLE nosuch.t (a int);\nDROP TABLE IF EXISTS nosuch.t;", complete: true},
	{name: "drop what is missing", change: "DROP TABLE nosuch;\nDROP TABLE IF EXISTS nosuch;\nDROP VIEW nosuch;\nDROP INDEX nosuch;\nDROP SEQUENCE nosuch;\nDROP MATERIALIZED VIEW nosuch;", complete: true},
	{name: "drop of another kind", change: "DROP TABLE tv;\nDROP TABLE IF EXISTS t_n_idx;\nDROP VIEW t;\nDROP INDEX q;\nDROP SEQUENCE IF EXISTS mv;\nDROP MATERIALIZED VIEW tv;", complete: true},
	{name: "drop then use", change: "DROP TABLE nokey;\nINSERT INTO nokey VALUES (1);\nALTER TABLE nokey ADD COLUMN b int;\nCREATE INDEX ON nokey (a);\nCOMMENT ON TABLE nokey IS 'x';\nCREATE TABLE nokey (a int);", complete: true},
	{name: "drop a table drops its indexes and sequences", change: "DROP TABLE t CASCADE;\nCREATE TABLE t_n_idx (a int);\nCREATE SEQUENCE t_id_seq;\nCREATE TABLE t_code_key (a int);"},
	{name: "dml targets", change: "INSERT INTO nosuch VALUES (1);\nUPDATE nosuch SET a = 1;\nDELETE FROM nosuch;\nWITH x AS (SELECT 1) UPDATE x SET a = 1;", complete: true},
	{name: "truncate", change: "TRUNCATE nosuch;\nTRUNCATE tv;\nTRUNCATE t, nokey;", complete: true},
	{name: "comment and grant", change: "COMMENT ON TABLE nosuch IS 'x';\nCOMMENT ON TABLE tv IS 'x';\nCOMMENT ON COLUMN t.nosuch IS 'x';\nCOMMENT ON COLUMN t.code IS 'x';\nGRANT SELECT ON nosuch TO PUBLIC;\nGRANT USAGE ON SCHEMA nosuch TO PUBLIC;\nCOMMENT ON SCHEMA nosuch IS 'x';", complete: true},
	{name: "index on what is not a table", change: "CREATE INDEX ON tv (code);\nCREATE INDEX ON q (last_value);\nCREATE INDEX ON nosuch (a);\nCREATE INDEX ON mv (id);", complete: true},
	{name: "a foreign key to the table being created, however written", change: "SET search_path = s, public;\nCREATE TABLE self2 (id int PRIMARY KEY, up int REFERENCES s.self2);\nCREATE TABLE public.self3 (id int PRIMARY KEY, up int REFERENCES self3);"},
	{name: "references", change: "CREATE TABLE x (a int REFERENCES nosuch);\nCREATE TABLE y (LIKE nosuch);\nCREATE TABLE z (a int) INHERITS (nosuch);\nCREATE TABLE self (id int PRIMARY KEY, up int REFERENCES self);\nALTER TABLE nokey ADD FOREIGN KEY (a) REFERENCES nosuch;", complete: true},

	// Columns.
	{name: "add a column that exists", change: "ALTER TABLE t ADD COLUMN code text;\nALTER TABLE t ADD COLUMN IF NOT EXISTS code text;\nALTER TABLE t ADD COLUMN xmin int;\nALTER TABLE t ADD COLUMN IF NOT EXISTS ctid int;\nALTER TABLE t ADD COLUMN oid int;", complete: true},
	{name: "drop a column that is missing", change: "ALTER TABLE t DROP COLUMN nosuch;\nALTER TABLE t DROP COLUMN IF EXISTS nosuch;\nALTER TABLE t DROP COLUMN ctid;\nALTER TABLE nokey DROP COLUMN a;\nALTER TABLE nokey DROP COLUMN a;", complete: true},
	{name: "alter a column", change: "ALTER TABLE t ALTER COLUMN nosuch TYPE text;\nALTER TABLE t ALTER COLUMN n SET DEFAULT 1;\nALTER TABLE t ALTER COLUMN xmin SET NOT NULL;\nALTER TABLE t ALTER COLUMN n SET NOT NULL;", complete: true},
	{name: "rename a column", change: "ALTER TABLE t RENAME COLUMN nosuch TO x;\nALTER TABLE t RENAME COLUMN n TO code;\nALTER TABLE t RENAME COLUMN n TO n2;\nALTER TABLE t DROP COLUMN n;\nALTER TABLE t ADD COLUMN n2 int;\nALTER TABLE t RENAME COLUMN ctid TO c;\nALTER TABLE t RENAME COLUMN n2 TO xmax;", complete: true},
	{name: "duplicate columns", change: "CREATE TABLE x (a int, b int, a text);\nCREATE TABLE y (a int PRIMARY KEY, b int, PRIMARY KEY (b));", complete: true},
	{name: "subcommands of one ALTER apply in passes", change: "ALTER TABLE nokey ADD COLUMN b int, DROP COLUMN b;\nALTER TABLE nokey DROP COLUMN a, ADD COLUMN a text;"},
	{name: "a column reaches inheritance children", change: "ALTER TABLE ip ADD COLUMN c int;\nALTER TABLE ic ADD COLUMN c int;\nALTER TABLE ip DROP COLUMN c;\nALTER TABLE ic DROP COLUMN c;\nALTER TABLE ip RENAME COLUMN a TO a2;\nALTER TABLE ic DROP COLUMN a;"},
	{name: "a table the change made has exactly its columns", change: "CREATE TABLE x (a int, b int);\nALTER TABLE x ADD COLUMN a int;\nALTER TABLE x DROP COLUMN c;\nCREATE TABLE y (LIKE t);\nALTER TABLE y ADD COLUMN code text;\nALTER TABLE y DROP COLUMN zz;", complete: false},

	// Constraints and keys.
	{name: "a second primary key", change: "ALTER TABLE t ADD PRIMARY KEY (n);\nALTER TABLE nokey ADD PRIMARY KEY (a);\nALTER TABLE nokey ADD PRIMARY KEY (a);\nALTER TABLE t ADD COLUMN k int PRIMARY KEY;", complete: true},
	{name: "a constraint name in use", change: "ALTER TABLE t ADD CONSTRAINT t_n_check CHECK (n > 1);\nALTER TABLE t ADD CONSTRAINT t_code_key UNIQUE (n);\nALTER TABLE t ADD CONSTRAINT t_p_fk CHECK (n > 1);\nALTER TABLE nokey ADD CONSTRAINT t_n_idx UNIQUE (a);\nALTER TABLE nokey ADD CONSTRAINT q CHECK (a > 0);", complete: true},
	{name: "drop a constraint", change: "ALTER TABLE t DROP CONSTRAINT nosuch;\nALTER TABLE t DROP CONSTRAINT t_n_check;\nALTER TABLE t DROP CONSTRAINT t_n_check;\nALTER TABLE t DROP CONSTRAINT IF EXISTS t_n_check;\nALTER TABLE t ADD CONSTRAINT t_n_check CHECK (n > 2);"},
	{name: "drop a key and its index", change: "ALTER TABLE t DROP CONSTRAINT t_code_key;\nCREATE TABLE t_code_key (a int);\nALTER TABLE t DROP CONSTRAINT t_pkey CASCADE;\nALTER TABLE t ADD PRIMARY KEY (n);", complete: true},
	{name: "rename a constraint", change: "ALTER TABLE t RENAME CONSTRAINT t_n_check TO t_n_check2;\nALTER TABLE t DROP CONSTRAINT t_n_check;\nALTER TABLE t RENAME CONSTRAINT t_code_key TO t_code_uk;\nCREATE TABLE t_code_key (a int);\nCREATE TABLE t_code_uk (a int);", complete: true},
	{name: "duplicate constraint names in one CREATE TABLE", change: "CREATE TABLE d1 (a int, CONSTRAINT c CHECK (true), CONSTRAINT c CHECK (false));\nCREATE TABLE d2 (a int, b int, CONSTRAINT c UNIQUE (a), CONSTRAINT c UNIQUE (b));\nCREATE TABLE d3 (a int, b int, CONSTRAINT c UNIQUE (a), CONSTRAINT c PRIMARY KEY (b));\nCREATE TABLE d4 (a int, CONSTRAINT c CHECK (a > 0), CONSTRAINT c UNIQUE (a));\nCREATE TABLE d5 (a int CONSTRAINT c REFERENCES p, CONSTRAINT c CHECK (a > 0));\nCREATE TABLE d6 (a int CONSTRAINT c REFERENCES p, b int CONSTRAINT c REFERENCES p);\nCREATE TABLE d7 (a int CONSTRAINT c CHECK (a > 0), b int CONSTRAINT c CHECK (b > 0));\nCREATE TABLE d1 (a int);", complete: true},
	{name: "a NOT NULL constraint's name is PostgreSQL 18's alone", change: "CREATE TABLE d8 (a int CONSTRAINT c NOT NULL, CONSTRAINT c CHECK (a > 0));\nCREATE TABLE d9 (a int CONSTRAINT nn NOT NULL);\nALTER TABLE d9 DROP CONSTRAINT nn;\nALTER TABLE nokey ADD COLUMN b int CONSTRAINT nn2 NOT NULL;\nALTER TABLE nokey DROP CONSTRAINT nn2;\nALTER TABLE nokey DROP CONSTRAINT nn2;"},
	{name: "PostgreSQL 18 names NOT NULL constraints", change: "ALTER TABLE nokey ALTER COLUMN a SET NOT NULL;\nALTER TABLE nokey DROP CONSTRAINT nokey_a_not_null;\nCREATE TABLE x (a int NOT NULL);\nALTER TABLE x DROP CONSTRAINT x_a_not_null;"},

	// Renames and moves.
	{name: "rename a table", change: "ALTER TABLE nokey RENAME TO t;\nALTER TABLE nokey RENAME TO mood;\nALTER TABLE nokey RENAME TO nk;\nALTER TABLE nk ADD COLUMN a int;\nALTER TABLE nokey ADD COLUMN b int;\nCREATE TABLE nokey (z int);", complete: true},
	{name: "rename an index", change: "ALTER INDEX t_n_idx RENAME TO t;\nALTER INDEX t_n_idx RENAME TO t_n_i;\nCREATE TABLE t_n_idx (a int);\nCREATE TABLE t_n_i (a int);", complete: true},
	{name: "set schema", change: "ALTER TABLE t SET SCHEMA s;\nALTER TABLE nokey SET SCHEMA nosuch;\nALTER TABLE nokey SET SCHEMA s;\nALTER TABLE s.nokey ADD COLUMN a int;\nCREATE TABLE nokey (a int);", complete: true},
	{name: "schemas", change: "CREATE SCHEMA s;\nCREATE SCHEMA IF NOT EXISTS s;\nCREATE SCHEMA pg_x;\nCREATE SCHEMA s2;\nCREATE TABLE s2.x (a int);\nCREATE TABLE s2.x (a int);\nALTER SCHEMA s2 RENAME TO s;\nALTER SCHEMA s2 RENAME TO s3;\nCREATE TABLE s3.x (a int);\nCREATE TABLE s2.y (a int);", complete: true},
	{name: "drop a schema", change: "DROP SCHEMA s;\nDROP SCHEMA nosuch;\nCREATE SCHEMA e;\nDROP SCHEMA e;\nCREATE SCHEMA f;\nCREATE TABLE f.x (a int);\nDROP SCHEMA f;\nDROP SCHEMA f CASCADE;\nCREATE TABLE f.y (a int);", complete: true},
	{name: "types", change: "DROP TYPE public.nosuch;\nCREATE TYPE public.e2 AS ENUM ('a');\nDROP TYPE public.e2;\nDROP TYPE public.e2;\nDROP TYPE IF EXISTS public.e2;"},

	// The search path and the session.
	{name: "search path", change: "SET search_path = s, public;\nCREATE TABLE t (a int);\nALTER TABLE t ADD COLUMN a int;\nALTER TABLE p ADD COLUMN id int;\nSET search_path = nosuch;\nCREATE TABLE z (a int);\nRESET search_path;\nCREATE TABLE z (a int);", complete: true},
	{name: "a view of a temporary relation needs no schema", change: "CREATE TEMP TABLE b (x int);\nSET search_path = nosuch;\nCREATE VIEW v AS SELECT x FROM b;\nCREATE TABLE z (a int);", complete: true},
	{name: "a schema the sync skips exists once something is created in it", change: "CREATE SCHEMA cron;\nCREATE TABLE cron.x (a int);\nCREATE TABLE cron.x (a int);", complete: true},
	{name: "temporary relations shadow", change: "CREATE TEMP TABLE t (x int);\nALTER TABLE t ADD COLUMN x int;\nALTER TABLE t ADD COLUMN code text;\nCREATE TEMP TABLE t (y int);\nCREATE VIEW tv2 AS SELECT * FROM t;\nCREATE TABLE tv2 (a int);", complete: true},
	{name: "pg_dump's search path", change: "SELECT pg_catalog.set_config('search_path', '', false);\nCREATE TABLE x (a int);\nCREATE TABLE public.t (a int);\nALTER TABLE t ADD COLUMN b int;\nSELECT pg_catalog.set_config('search_path', 's', false);\nCREATE TABLE t (a int);", complete: true},
	{name: "a foreign key naming another schema is not the new table", change: "CREATE TABLE public.x2 (id int REFERENCES missing.x2);\nCREATE TABLE public.x2 (id int);\nSET search_path = s, public;\nCREATE TABLE x3 (id int PRIMARY KEY, r int REFERENCES public.x3);\nCREATE TABLE x4 (id int PRIMARY KEY, r int REFERENCES s.x4);", complete: true},
	{name: "one column added or dropped twice in one ALTER", change: "ALTER TABLE nokey ADD COLUMN b int, ADD COLUMN b int;\nALTER TABLE nokey ADD COLUMN IF NOT EXISTS b int, ADD COLUMN IF NOT EXISTS b int;\nALTER TABLE t DROP COLUMN p_id, DROP COLUMN p_id;\nALTER TABLE t DROP COLUMN IF EXISTS n, DROP COLUMN IF EXISTS n;\nALTER TABLE t DROP COLUMN xmin, DROP COLUMN xmin;\nALTER TABLE t ADD COLUMN xmin int, ADD COLUMN xmin int;", complete: true},
	{name: "SET SCHEMA to the relation's own schema is a no-op", change: "ALTER TABLE public.nokey SET SCHEMA public;\nALTER TABLE nokey ADD COLUMN a int;", complete: true},
	{name: "ALTER SEQUENCE needs a sequence", change: "ALTER SEQUENCE public.nokey RESTART;\nALTER SEQUENCE IF EXISTS t_n_idx RESTART;\nALTER SEQUENCE q RESTART;", complete: true},
	{name: "column names compare truncated", change: "CREATE TABLE e1 (\"" + strings.Repeat("c", 63) + "A\" int, \"" + strings.Repeat("c", 63) + "B\" int);\nCREATE TABLE e2 (a int);", complete: true},
	{name: "a key named like its new table", change: "CREATE TABLE g1 (a int CONSTRAINT g1 UNIQUE);\nCREATE TABLE g2 (a int, CONSTRAINT g2 PRIMARY KEY (a));\nCREATE TABLE g3 (a int CONSTRAINT g3 CHECK (a > 0));\nCREATE TABLE g1 (a int);", complete: true},
	{name: "a dropped partition's name is no longer the parent's", change: "CREATE TABLE ph (k int) PARTITION BY RANGE (k);\nCREATE TABLE pc PARTITION OF ph FOR VALUES FROM (0) TO (10);\nDROP TABLE pc;\nCREATE TABLE pc (a int);\nDROP TABLE ph;\nALTER TABLE pc ADD COLUMN b int;\nALTER TABLE pc ADD COLUMN b int;", complete: true},
	{name: "a routine holds its schema, which the walk does not follow", change: "CREATE SCHEMA sf;\nCREATE FUNCTION sf.f() RETURNS int LANGUAGE sql AS 'SELECT 1';\nDROP SCHEMA sf;\nCREATE TABLE sf.t (a int);"},
	{name: "a foreign key into another database is not the new table", change: "CREATE TABLE public.x5 (id int PRIMARY KEY, r int REFERENCES other.public.x5);\nCREATE TABLE public.x5 (id int);", complete: true},
	{name: "constraint names compare truncated", change: "CREATE TABLE y1 (a int, CONSTRAINT \"" + strings.Repeat("k", 63) + "A\" CHECK (a > 0), CONSTRAINT \"" + strings.Repeat("k", 63) + "B\" CHECK (a > 1));\nCREATE TABLE y1 (a int);", complete: true},
	{name: "identical non-ASCII column names", change: "CREATE TABLE z1 (\"é\" int, \"é\" text);\nCREATE TABLE z2 (é int, é text);\nCREATE TABLE z1 (a int);", complete: true},
	{name: "USING INDEX attaches an existing index", change: "CREATE UNIQUE INDEX u ON nokey (a);\nALTER TABLE nokey ADD CONSTRAINT u UNIQUE USING INDEX u;\nALTER TABLE nokey ADD CONSTRAINT u UNIQUE USING INDEX u;\nCREATE UNIQUE INDEX u2 ON t (n);\nALTER TABLE t ADD CONSTRAINT u3 UNIQUE USING INDEX u2;\nCREATE TABLE u2 (a int);\nCREATE TABLE u3 (a int);\nALTER TABLE nokey ADD CONSTRAINT c9 UNIQUE USING INDEX nosuch;\nALTER TABLE nokey ADD CONSTRAINT c10 UNIQUE USING INDEX t;", complete: true},
	{name: "IF NOT EXISTS CREATE TABLE AS meets a type", change: "CREATE TABLE IF NOT EXISTS mood AS SELECT 1 AS x;\nCREATE MATERIALIZED VIEW IF NOT EXISTS mood AS SELECT 1 AS x;\nCREATE TABLE IF NOT EXISTS tv AS SELECT 1 AS x;", complete: true},
	{name: "CREATE INDEX CONCURRENTLY on a partitioned table", change: "CREATE INDEX CONCURRENTLY parted_kc ON parted (k);"},
	{name: "a view OR REPLACE creates holds its schema", change: "CREATE SCHEMA sv;\nCREATE OR REPLACE VIEW sv.v AS SELECT 1 AS x;\nDROP SCHEMA sv;\nCREATE TABLE sv.t (a int);", complete: true},
	{name: "another database", change: "CREATE TABLE other.public.x (a int);\nALTER TABLE other.public.t ADD COLUMN b int;\nDROP TABLE IF EXISTS other.public.t;\nINSERT INTO other.public.t VALUES (1);\nCOMMENT ON TABLE other.public.t IS 'x';\nCREATE TABLE omni_review.public.x (a int);\nCREATE TABLE omni_review.public.x (a int);\nCREATE TABLE x (a int);", complete: true},
	{name: "pg_catalog the SQL names is not judged", change: "DROP TABLE pg_catalog.foo;\nALTER TABLE pg_catalog.foo ADD COLUMN a int;\nCREATE TABLE pg_catalog.foo (a int);\nCREATE TABLE foo (a int);\nALTER TABLE nosuch ADD COLUMN a int;"},
	{name: "IF NOT EXISTS partitioning of a table that may exist", change: "CREATE TABLE ph (k int) PARTITION BY RANGE (k);\nCREATE TABLE IF NOT EXISTS s.parted_1 PARTITION OF ph FOR VALUES FROM (0) TO (10);\nDROP TABLE ph;\nCOMMENT ON TABLE s.parted_1 IS 'x';\nCREATE TABLE IF NOT EXISTS s.parted_1 (k int) PARTITION BY RANGE (k);\nCREATE INDEX CONCURRENTLY pk1 ON s.parted_1 (k);"},
	{name: "IF NOT EXISTS keys of a table that may exist", change: "CREATE TABLE t2 (id serial);\nCREATE TABLE IF NOT EXISTS t2x (a int CONSTRAINT uq UNIQUE, b int GENERATED ALWAYS AS IDENTITY (SEQUENCE NAME iseq));\nALTER INDEX uq RENAME TO uq2;\nALTER SEQUENCE iseq RESTART;"},
	{name: "ADD COLUMN IF NOT EXISTS of a column that may exist may skip its constraints", change: "CREATE TABLE x (LIKE t);\nALTER TABLE x ADD COLUMN IF NOT EXISTS code text CONSTRAINT uq3 UNIQUE;\nCREATE TABLE uq3 (a int);\nALTER TABLE x ADD COLUMN IF NOT EXISTS p_id int REFERENCES nosuch;\nALTER TABLE x ADD PRIMARY KEY (id);\nALTER TABLE x ADD COLUMN IF NOT EXISTS n int PRIMARY KEY;"},
	{name: "OR REPLACE VIEW of a relation that is no view", change: "CREATE OR REPLACE VIEW nokey AS SELECT 1 AS a;\nCREATE OR REPLACE VIEW q AS SELECT 1 AS a;\nCREATE OR REPLACE VIEW t_n_idx AS SELECT 1 AS a;\nCREATE OR REPLACE VIEW mv AS SELECT 1 AS a;\nCREATE OR REPLACE VIEW pair AS SELECT 1 AS a;\nCREATE OR REPLACE VIEW mood AS SELECT 1 AS a;\nCREATE OR REPLACE VIEW nokey AS SELECT * FROM nosuch;\nCREATE OR REPLACE VIEW tv AS SELECT code FROM t;\nDROP VIEW tv;", complete: true},
	{name: "one constraint name added twice in one ALTER", change: "ALTER TABLE nokey ADD CONSTRAINT c CHECK (true), ADD CONSTRAINT c CHECK (false);\nALTER TABLE nokey ADD CONSTRAINT c UNIQUE (a), ADD CONSTRAINT c PRIMARY KEY (a);\nALTER TABLE nokey ADD CONSTRAINT c CHECK (true), ADD CONSTRAINT c UNIQUE (a);\nALTER TABLE nokey ADD CONSTRAINT c FOREIGN KEY (a) REFERENCES p, ADD CONSTRAINT c CHECK (true);\nALTER TABLE nokey ADD COLUMN d int CONSTRAINT c UNIQUE, ADD CONSTRAINT c CHECK (true);\nALTER TABLE nokey ADD COLUMN IF NOT EXISTS a int CONSTRAINT c UNIQUE, ADD CONSTRAINT c CHECK (true);\nALTER TABLE nokey ADD CONSTRAINT c2 CHECK (zz > 0), ADD CONSTRAINT c2 CHECK (true);\nCREATE TABLE c (a int);", complete: true},
	{name: "a renamed or moved relation still holds its schema", change: "CREATE SCHEMA s2;\nALTER SCHEMA s RENAME TO s3;\nDROP SCHEMA s3;\nALTER TABLE s3.t RENAME TO u;\nDROP SCHEMA s3;\nALTER TABLE s3.u SET SCHEMA s2;\nDROP SCHEMA s2;", complete: true},
	{name: "a new table's system column names", change: "CREATE TABLE x (ctid int);\nCREATE TABLE x (a int, tableoid int);\nCREATE TABLE x (xmin int, xmax int);\nCREATE TABLE IF NOT EXISTS x (cmin int);\nCREATE TABLE x (cmax int) INHERITS (p);\nCREATE TABLE x (ctid nosuchtype);\nCREATE TABLE x (\"CTID\" int, oid int);\nCREATE TABLE IF NOT EXISTS x (ctid int);\nCREATE TABLE IF NOT EXISTS s.parted_1 (ctid int);", complete: true},
	{name: "an unqualified self-reference follows the search path", change: "SET search_path = s;\nCREATE TABLE public.x (id int PRIMARY KEY, r int REFERENCES x (id));\nSET search_path = s, public;\nCREATE TABLE public.x (id int PRIMARY KEY, r int REFERENCES x (id));\nCREATE TEMP TABLE y (id int PRIMARY KEY, r int REFERENCES y (id));", complete: true},
	{name: "an ALTER of a relation needs its kind", change: "ALTER VIEW nokey OWNER TO postgres;\nALTER MATERIALIZED VIEW tv OWNER TO postgres;\nALTER SEQUENCE nokey OWNER TO postgres;\nALTER INDEX nokey SET TABLESPACE pg_default;\nALTER INDEX IF EXISTS q SET TABLESPACE pg_default;\nALTER FOREIGN TABLE nokey ADD COLUMN zz int;\nALTER TABLE pair OWNER TO postgres;\nALTER TABLE tv OWNER TO postgres;\nALTER INDEX t_n_idx SET TABLESPACE pg_default;\nALTER VIEW nokey RENAME TO u;\nALTER SEQUENCE tv RENAME TO u;\nALTER TABLE pair RENAME TO u;\nALTER INDEX nokey RENAME TO u;\nALTER VIEW mv SET SCHEMA s;\nALTER TABLE t_n_idx SET SCHEMA s;\nALTER TABLE pair SET SCHEMA s;\nALTER TABLE q SET SCHEMA s;\nCREATE TABLE u (a int);\nCREATE TABLE s.q (a int);", complete: true},
	{name: "one ALTER repeats an unreadable column or a constraint drop", change: "ALTER TABLE nokey ADD COLUMN \"é\" int, ADD COLUMN \"é\" int;\nALTER TABLE nokey ADD COLUMN \"é\" int;\nALTER TABLE nokey DROP COLUMN \"é\", DROP COLUMN \"é\";\nALTER TABLE t DROP CONSTRAINT t_n_check, DROP CONSTRAINT t_n_check;\nALTER TABLE t DROP CONSTRAINT IF EXISTS t_n_check, DROP CONSTRAINT t_n_check;\nALTER TABLE t DROP CONSTRAINT t_n_check, DROP CONSTRAINT IF EXISTS t_n_check;\nALTER TABLE t DROP CONSTRAINT t_n_check;", complete: true},
	{
		name:   "an unqualified set_config may be a user function",
		extra:  "CREATE FUNCTION s.set_config(text, text, boolean) RETURNS text LANGUAGE sql AS $$ SELECT 'x'::text $$;",
		change: "SET search_path = s, pg_catalog, public;\nSELECT set_config('search_path', 'public', false);\nCREATE TABLE x (a int);\nCREATE TABLE public.x (a int);",
	},
	{name: "DROP INDEX CONCURRENTLY drops one index without CASCADE", change: "DROP INDEX CONCURRENTLY t_n_idx, mv_id_idx;\nDROP INDEX CONCURRENTLY t_n_idx CASCADE;\nDROP INDEX CONCURRENTLY IF EXISTS x, y;\nDROP INDEX CONCURRENTLY IF EXISTS x CASCADE;"},
	{name: "TRUNCATE ONLY of a partitioned table", change: "TRUNCATE ONLY parted;\nTRUNCATE nosuch, ONLY parted;\nTRUNCATE ONLY parted, nosuch;\nTRUNCATE ONLY public.parted;", complete: true},
	{name: "TRUNCATE ONLY of a partition", change: "TRUNCATE ONLY s.parted_1;\nTRUNCATE parted;"},
	{name: "one ALTER drops before it adds", change: "ALTER TABLE t ADD CONSTRAINT t_n_check CHECK (n > 1), DROP CONSTRAINT t_n_check;\nALTER TABLE t DROP CONSTRAINT t_n_check;\nALTER TABLE nokey ADD COLUMN a int, DROP COLUMN a;\nALTER TABLE nokey DROP COLUMN a;\nALTER TABLE t ADD CONSTRAINT t_code_key UNIQUE (n), DROP CONSTRAINT t_code_key;\nCREATE TABLE t_code_key (a int);"},
	{name: "a dropped schema takes its partition state", change: "CREATE SCHEMA ps;\nCREATE TABLE ps.p (k int) PARTITION BY RANGE (k);\nCREATE TABLE ps.p1 PARTITION OF ps.p FOR VALUES FROM (0) TO (10);\nDROP SCHEMA ps CASCADE;\nCREATE SCHEMA ps;\nCREATE TABLE ps.q (k int);\nALTER TABLE ps.q RENAME TO p;\nTRUNCATE ONLY ps.p;"},
	{name: "ADD COLUMN IF NOT EXISTS of an unreadable name may skip its constraints", change: "CREATE TABLE x (\"é\" int);\nALTER TABLE x ADD COLUMN IF NOT EXISTS \"é\" int CONSTRAINT uq UNIQUE;\nCREATE TABLE uq (a int);"},
	{name: "set_config of a setting named at run time", change: "SELECT pg_catalog.set_config('work' || '_mem', '4MB', false);\nCREATE TABLE x (a int);"},
	{name: "ADD COLUMN IF NOT EXISTS after the statement drops the column", change: "ALTER TABLE nokey DROP COLUMN a, ADD COLUMN IF NOT EXISTS a integer;\nALTER TABLE nokey DROP COLUMN a;\nALTER TABLE nokey ADD COLUMN IF NOT EXISTS a int CONSTRAINT nk_u UNIQUE, DROP COLUMN IF EXISTS a;\nCREATE TABLE nk_u (z int);", complete: true},
	{name: "ADD COLUMN IF NOT EXISTS after the statement adds the column", change: "ALTER TABLE nokey ADD COLUMN b int, ADD COLUMN IF NOT EXISTS b int CONSTRAINT uq UNIQUE;\nCREATE TABLE uq (a int);\nALTER TABLE nokey ADD COLUMN IF NOT EXISTS c int, ADD COLUMN IF NOT EXISTS c int CONSTRAINT uq2 UNIQUE;\nCREATE TABLE uq2 (a int);\nALTER TABLE nokey DROP COLUMN c, ADD COLUMN c int, ADD COLUMN IF NOT EXISTS c int CONSTRAINT uq3 UNIQUE;\nCREATE TABLE uq3 (a int);"},
	{name: "a partition proves its parent partitioned", extra: "CREATE TABLE ep (k int) PARTITION BY RANGE (k);", change: "CREATE TABLE ep1 PARTITION OF ep FOR VALUES FROM (0) TO (10);\nTRUNCATE ONLY ep;", complete: true},
	{name: "a multi-object DROP looks every object up first", change: "DROP TABLE nokey, q;\nDROP SEQUENCE q, nokey;\nDROP TABLE nokey, nosuch CASCADE;\nDROP TABLE nosuch, nokey;\nDROP TABLE nokey, nokey;\nCREATE TABLE nokey (a int);", complete: true},
	{name: "DROP DOMAIN of a type that is no domain", change: "DROP DOMAIN public.mood;\nDROP DOMAIN IF EXISTS public.pair;\nDROP DOMAIN public.nokey;\nCREATE DOMAIN public.dd AS int;\nDROP DOMAIN public.dd;\nCREATE TYPE public.e2 AS ENUM ('a');\nDROP DOMAIN public.e2;\nDROP TYPE public.e2;\nCREATE DOMAIN public.d2 AS int;\nDROP TYPE public.d2;", complete: true},
	{name: "REFRESH runs the view's query", change: "CREATE TABLE victim (a int);\nCREATE FUNCTION kill() RETURNS int LANGUAGE plpgsql AS $$ BEGIN DROP TABLE public.victim; RETURN 1; END $$;\nCREATE MATERIALIZED VIEW mvk AS SELECT public.kill() AS x WITH NO DATA;\nREFRESH MATERIALIZED VIEW mvk;\nCREATE TABLE victim (a int);"},
	{name: "a role's schema", change: "CREATE SCHEMA AUTHORIZATION bob;\nSET ROLE bob;\nCREATE TABLE x (a int);\nRESET ROLE;\nCREATE TABLE bob.x (a int);\nCREATE TABLE x (a int);", complete: true},

	// How Bytebase runs a change: CONCURRENTLY, VACUUM, and SET ROLE after
	// its transaction commits, or first in auto-commit mode.
	{name: "CONCURRENTLY runs after the transaction", change: "CREATE INDEX CONCURRENTLY ci ON t2 (a);\nCREATE TABLE t2 (a int);\nCREATE INDEX CONCURRENTLY ci ON t2 (a);\nDROP INDEX CONCURRENTLY t_n_idx;\nCREATE INDEX t_n_idx ON t (n);", complete: true},
	{name: "CONCURRENTLY on a table the transaction replaced", change: "CREATE INDEX CONCURRENTLY pk ON parted (k);\nDROP TABLE parted;\nCREATE TABLE parted (k int);"},
	{name: "auto-commit mode runs CONCURRENTLY first", change: "-- txn-mode = off\nCREATE TABLE t_n_idx (a int);\nDROP INDEX CONCURRENTLY t_n_idx;\nCREATE TABLE ix (a int);\nCREATE INDEX CONCURRENTLY ix ON t (n);", complete: true},
	{name: "SET ROLE runs again after the transaction", extra: "CREATE SCHEMA bob AUTHORIZATION bob;\nCREATE TABLE bob.bt (a int);", change: "SET ROLE bob;\nCREATE INDEX CONCURRENTLY bi ON bt (a);\nRESET ROLE;\nCREATE INDEX CONCURRENTLY bi2 ON bt (a);"},
	{name: "a CONCURRENTLY Bytebase does not set apart fails in its transaction", change: "CREATE /* x */ INDEX CONCURRENTLY t_n_idx ON t (n);\nDROP TABLE nosuch;"},
	{name: "views and materialized views the sync drops", extra: "CREATE VIEW hypopg_list_indexes AS SELECT 1 AS x;\nCREATE MATERIALIZED VIEW g_agg_stat_m AS SELECT 1 AS x;", change: "GRANT SELECT ON hypopg_list_indexes TO PUBLIC;\nCOMMENT ON MATERIALIZED VIEW g_agg_stat_m IS 'x';\nDROP VIEW hypopg_list_indexes;\nDROP VIEW hypopg_list_indexes;"},
	{name: "partitions of a partitioned partition in another schema the sync loses", extra: "CREATE SCHEMA s2;\nCREATE TABLE lr (k int, j int) PARTITION BY RANGE (k);\nCREATE TABLE s2.lp PARTITION OF lr FOR VALUES FROM (0) TO (10) PARTITION BY RANGE (j);\nCREATE TABLE s2.lsp PARTITION OF s2.lp FOR VALUES FROM (0) TO (10);", change: "ALTER TABLE s2.lsp ADD COLUMN z int;\nCOMMENT ON TABLE s2.lsp IS 'x';\nDROP TABLE s2.lsp;"},
	{name: "an index shows a partitioned partition in another schema", extra: "CREATE SCHEMA s3;\nCREATE TABLE r3 (k int, j int, PRIMARY KEY (k, j)) PARTITION BY RANGE (k);\nCREATE TABLE s3.r3p PARTITION OF r3 FOR VALUES FROM (0) TO (10) PARTITION BY RANGE (j);\nCREATE TABLE s3.r3sp PARTITION OF s3.r3p FOR VALUES FROM (0) TO (10);", change: "ALTER TABLE s3.r3sp ADD COLUMN z int;\nDROP TABLE s3.r3sp;"},
	{name: "a partitioned partition in its table's schema loses nothing", extra: "CREATE TABLE r2 (k int, j int, PRIMARY KEY (k, j)) PARTITION BY RANGE (k);\nCREATE TABLE r2p PARTITION OF r2 FOR VALUES FROM (0) TO (10) PARTITION BY RANGE (j);\nCREATE TABLE s.r2sp PARTITION OF r2p FOR VALUES FROM (0) TO (10);", change: "ALTER TABLE nosuch ADD COLUMN a int;\nCOMMENT ON TABLE s.r2sp IS 'x';\nALTER TABLE s.nosuch ADD COLUMN b int;", complete: true},
	{name: "tenant mode: RESET ROLE returns to a session user the walk does not know", tenant: true, extra: "CREATE SCHEMA walk_owner AUTHORIZATION walk_owner;\nCREATE TABLE walk_owner.x (a int);\nGRANT ALL ON SCHEMA public TO walk_owner;", change: "RESET ROLE;\nCREATE TABLE x (a int);\nCREATE TABLE x (a int);"},
	{name: "tenant mode: GRANT runs as the session user", tenant: true, extra: "CREATE SCHEMA postgres;\nCREATE TABLE postgres.g (a int);\nCREATE SCHEMA walk_owner AUTHORIZATION walk_owner;\nGRANT ALL ON SCHEMA public TO walk_owner;", change: "GRANT SELECT ON g TO PUBLIC;\nCREATE TABLE x (a int);\nCREATE TABLE x (a int);", complete: true},

	// Where the walk stops.
	{
		name: "an event trigger ends the walk where it may fire",
		extra: `CREATE FUNCTION et() RETURNS event_trigger LANGUAGE plpgsql AS $$ BEGIN
			IF to_regclass('public.x') IS NOT NULL THEN EXECUTE 'DROP TABLE public.x'; END IF; END $$;
			CREATE EVENT TRIGGER et ON ddl_command_end WHEN TAG IN ('CREATE TABLE') EXECUTE FUNCTION et();`,
		change: "ALTER TABLE nosuch ADD COLUMN a int;\nCREATE TABLE x (a int);\nCREATE TABLE x (a int);\nALTER TABLE nosuch ADD COLUMN a int;",
	},
	{name: "procedural code ends the walk", change: "DO $$ BEGIN EXECUTE 'DROP TABLE nokey'; END $$;\nCREATE TABLE nokey (a int);"},
	{name: "a call of a user function ends the walk", change: "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN EXECUTE 'DROP TABLE nokey'; RETURN 1; END $$;\nSELECT f();\nCREATE TABLE nokey (a int);"},
	{name: "a trigger ends the walk at DML", change: "CREATE FUNCTION tf() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN EXECUTE 'DROP TABLE p CASCADE'; RETURN NEW; END $$;\nCREATE TRIGGER tr BEFORE INSERT ON nokey FOR EACH ROW EXECUTE FUNCTION tf();\nINSERT INTO nokey VALUES (1);\nCREATE TABLE p (a int);"},
	{name: "a rollback ends the walk", change: "SAVEPOINT a;\nCREATE TABLE x (a int);\nROLLBACK TO SAVEPOINT a;\nCREATE TABLE x (a int);"},
	{name: "generated names", change: "CREATE TABLE g (id serial PRIMARY KEY, a int UNIQUE);\nCREATE TABLE g_pkey (a int);\nCREATE TABLE g_id_seq (a int);\nCREATE TABLE g_a_key (a int);\nCREATE INDEX ON g (a);\nCREATE TABLE g_a_idx (a int);\nCREATE TABLE gx (a int);\nCREATE TABLE gx (a int);"},
	{name: "partitions", change: "CREATE TABLE parted_2 PARTITION OF parted FOR VALUES FROM (10) TO (20);\nCREATE TABLE parted_2 (a int);\nCREATE INDEX parted_k ON parted (k);\nDROP TABLE parted;\nCREATE TABLE parted_2 (a int);\nCREATE TABLE s.parted_1 (a int);"},
	{name: "a partition from the synced schema", change: "ALTER TABLE s.parted_1 ADD COLUMN zz int;\nCREATE TABLE s.parted_1 (a int);\nDROP TABLE parted_1;"},
}

func TestWalkThroughAgainstPostgres(t *testing.T) {
	runWalkCases(t, startPostgres(t))
}

// TestWalkThroughAgainstPostgres18 runs the same cases on PostgreSQL 18,
// whose NOT NULL constraints have names: a fact the walk relies on must
// hold on every version Bytebase supports.
func TestWalkThroughAgainstPostgres18(t *testing.T) {
	runWalkCases(t, startPostgresImage(t, "postgres:18-alpine"))
}

func runWalkCases(t *testing.T, db *sql.DB) {
	t.Helper()
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	setupWalkRoles(t, conn)
	var found, rejected int
	for _, tt := range walkCases {
		t.Run(tt.name, func(t *testing.T) {
			resetWalk(t, conn, walkBase+tt.extra)
			schema := walkSync(t, conn)
			user := tt.user
			switch {
			case tt.tenant:
				user = walkOwner
			case user == "":
				user = "postgres"
			}
			f, r := judgeWalk(t, conn, schema, user, tt.change, tt.complete, tt.tenant)
			found += f
			rejected += r
		})
	}
	t.Logf("found %d of the %d statements the server rejects", found, rejected)
}

// setupWalkRoles creates bob, a role the cases set, and walkOwner, who
// then owns the database.
func setupWalkRoles(t *testing.T, conn *sql.Conn) {
	t.Helper()
	for _, q := range []string{
		"DO $$ BEGIN CREATE ROLE bob; EXCEPTION WHEN duplicate_object THEN NULL; END $$",
		"DO $$ BEGIN CREATE ROLE " + walkOwner + "; EXCEPTION WHEN duplicate_object THEN NULL; END $$",
		"ALTER DATABASE omni_review OWNER TO " + walkOwner,
	} {
		if _, err := conn.ExecContext(context.Background(), q); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
	}
}

// TestWalkThroughOracleChecksClaims holds the oracle itself to its job:
// a claim the catalog contradicts must not hold.
func TestWalkThroughOracleChecksClaims(t *testing.T) {
	db := startPostgres(t)
	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	resetWalk(t, conn, walkBase)
	nokey := relKey("public", "nokey")
	for _, tt := range []struct {
		name string
		v    walk.Violation
		msg  string
		want bool
	}{
		{"an existing relation claimed absent", walk.Violation{Keys: []walk.Key{schemaKey("public"), nokey}, Facts: []walk.Fact{walk.Present, walk.Absent}}, "relation nokey does not exist", false},
		{"an existing relation claimed present", walk.Violation{Keys: []walk.Key{nokey}, Facts: []walk.Fact{walk.Present}}, "", true},
		{"a missing schema claimed present", walk.Violation{Keys: []walk.Key{schemaKey("nosuch")}, Facts: []walk.Fact{walk.Present}}, "", false},
		{"a table claimed a view", walk.Violation{Check: walk.Check{Key: nokey}, Got: walk.Present, Kind: relView}, "nokey is not a table", false},
		{"a table claimed a table", walk.Violation{Check: walk.Check{Key: nokey}, Got: walk.Present, Kind: relTable}, "", true},
		{"an ordinary table claimed partitioned", walk.Violation{Keys: []walk.Key{nokey}, Facts: []walk.Fact{walk.Present}}, "cannot truncate only a partitioned table: nokey", false},
		{"a claim about the text", walk.Violation{}, "column a specified more than once", true},
	} {
		if got, why := verifyClaim(t, conn, &tt.v, tt.msg); got != tt.want {
			t.Errorf("%s: holds = %v (%s), want %v", tt.name, got, why, tt.want)
		}
	}
}

// TestWalkThroughSweep runs random changes over a small pool of names, so
// creations, drops, renames, and generated names collide often, on
// random starting schemas. Precision must hold on every one; recall is
// logged. WALK_SWEEP_N sets the number of changes, WALK_SWEEP_SEED the
// seed, and WALK_SWEEP_IMAGE the PostgreSQL image.
func TestWalkThroughSweep(t *testing.T) {
	n, seed, image := walktest.SweepParams(t, 40, 20261019, "postgres:17-alpine")
	db := startPostgresImage(t, image)
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	var version string
	if err := conn.QueryRowContext(ctx, "SHOW server_version").Scan(&version); err != nil {
		t.Fatal(err)
	}
	t.Logf("PostgreSQL %s", version)
	setupWalkRoles(t, conn)
	rng := rand.New(rand.NewSource(seed))
	var found, rejected int
	for i := range n {
		base := sweepChange(rng, 8+rng.Intn(8))
		change := sweepChange(rng, 15+rng.Intn(15))
		if rng.Intn(6) == 0 {
			change = "-- txn-mode = off\n" + change
		}
		t.Run(fmt.Sprintf("change%d", i), func(t *testing.T) {
			resetWalk(t, conn, "")
			// The base keeps only the statements the server accepts. The
			// change then runs as Bytebase runs it, in a session of its
			// own: no temporary relation, the default settings.
			for _, stmt := range splitStatements(base) {
				_, _ = conn.ExecContext(ctx, stmt)
			}
			for _, stmt := range []string{"RESET ALL", "RESET ROLE", "DISCARD TEMP"} {
				if _, err := conn.ExecContext(ctx, stmt); err != nil {
					t.Fatalf("%s: %v", stmt, err)
				}
			}
			schema := walkSync(t, conn)
			f, r := judgeWalk(t, conn, schema, "postgres", change, false, false)
			found += f
			rejected += r
			if t.Failed() {
				t.Logf("base:\n%s\nchange:\n%s", base, change)
			}
		})
	}
	t.Logf("seed %d: found %d of the %d statements the server rejects", seed, found, rejected)
}

// splitStatements splits a change into its statements, each ending with a
// semicolon.
func splitStatements(sql string) []string {
	var out []string
	for _, seg := range pg.Split(sql) {
		if seg.Empty() {
			continue
		}
		text := strings.TrimSpace(sql[seg.ByteStart:seg.ByteEnd])
		if !strings.HasSuffix(text, ";") {
			text += ";"
		}
		out = append(out, text)
	}
	return out
}

// pgEngine is the oracle's PostgreSQL side of one change, which it runs
// as Bytebase's executor does (bytebaseExec). Its own transaction, which
// judgeWalk rolls back, stands for Bytebase's: there each statement runs
// in a savepoint, so a rejected one does nothing and the next runs as the
// change would after it if it had been left out. What Bytebase runs after
// committing its transaction, or in auto-commit mode, the oracle runs
// after committing, where PostgreSQL undoes a rejected statement itself.
type pgEngine struct {
	conn   *sql.Conn
	target review.Target
	tenant bool

	// The change as Bytebase runs it, and how far Exec has run it: units
	// in order, which of those run accepted, and the oracle's transaction.
	autocommit bool
	stmts      []string
	units      []bbUnit
	first      []int
	next       int
	accepted   []bool
	inTxn      bool
	// mismatch is how the walk's planRun differs from bytebaseExec.
	mismatch string
}

func (e *pgEngine) exec(t testing.TB, q string) {
	t.Helper()
	if _, err := e.conn.ExecContext(context.Background(), q); err != nil {
		t.Fatalf("%s: %v", q, err)
	}
}

// Statements runs nothing: it lays the change out as Bytebase runs it and
// checks that the walk lays it out alike.
func (e *pgEngine) Statements(change string) []string {
	e.autocommit, e.stmts, e.units = bytebaseExec(change, e.tenant, walkOwner)
	e.first = make([]int, len(e.stmts))
	for i := range e.first {
		e.first[i] = -1
	}
	texts := make([]string, len(e.units))
	for i, u := range e.units {
		if e.first[u.stmt] < 0 {
			e.first[u.stmt] = i
		}
		texts[i] = u.text
	}
	e.accepted = make([]bool, len(e.units))
	e.inTxn = true
	e.mismatch = e.samePlan(change)
	return texts
}

// samePlan compares the walk's planRun of the change with bytebaseExec's,
// and returns how they differ.
func (e *pgEngine) samePlan(change string) string {
	stmts, syntax := parse(change, statementRanges(change))
	if syntax != nil {
		return "syntax: " + syntax.Message
	}
	if len(stmts) != len(e.stmts) {
		return fmt.Sprintf("%d statements parsed, %d split", len(stmts), len(e.stmts))
	}
	plan := planRun(change, stmts)
	if plan == nil || plan.autocommit != e.autocommit || len(plan.steps) != len(e.units) {
		return fmt.Sprintf("planRun %+v, Bytebase runs %d units, auto-commit %v", plan, len(e.units), e.autocommit)
	}
	for i, u := range e.units {
		want := step{stmt: u.stmt, again: e.first[u.stmt] != i, committed: !e.autocommit && u.apart && (i == 0 || !e.units[i-1].apart), wrapped: u.superuser}
		if plan.steps[i] != want {
			return fmt.Sprintf("step %d: planRun %+v, Bytebase %+v", i, plan.steps[i], want)
		}
	}
	return ""
}

func (e *pgEngine) Exec(t testing.TB, _ string) walktest.Outcome {
	t.Helper()
	i := e.next
	e.next++
	u := e.units[i]
	if e.inTxn && (e.autocommit || u.apart) {
		e.exec(t, "COMMIT")
		e.inTxn = false
		if e.tenant {
			e.exec(t, "SET SESSION ROLE '"+walkOwner+"'")
		}
	} else if e.inTxn && i == 0 && e.tenant {
		e.exec(t, "SET LOCAL ROLE '"+walkOwner+"'")
	}
	// The oracle does not follow a transaction the change begins or ends.
	upper := strings.ToUpper(strings.TrimSpace(u.text))
	if strings.HasPrefix(upper, "BEGIN") || strings.HasPrefix(upper, "COMMIT") {
		e.accepted[i] = true
		return walktest.Outcome{}
	}
	if e.inTxn {
		e.exec(t, "SAVEPOINT walk_oracle")
		defer e.exec(t, "RELEASE SAVEPOINT walk_oracle")
	}
	if _, err := e.conn.ExecContext(context.Background(), u.text); err != nil {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			t.Fatalf("%s: %v", u.text, err)
		}
		// 42601 is not only a syntax error: PostgreSQL raises it for some
		// semantic refusals too (SET DEFAULT on an identity column). Text
		// the parser refuses shows as a Syntax finding.
		if e.inTxn {
			e.exec(t, "ROLLBACK TO SAVEPOINT walk_oracle")
		}
		return walktest.Outcome{Code: pgErr.Code, Message: pgErr.Message}
	}
	e.accepted[i] = true
	return walktest.Outcome{}
}

// replay returns the change as the walk is to see it where Judge asks:
// the statements whose first run the server accepted, and, unless text is
// those runs alone, the statement of the run Exec made last, as written,
// with its directives. Each statement maps to the position of its first
// run among those runs, which is how Judge counts.
func (e *pgEngine) replay(text string) (string, []int) {
	var runs []int
	var texts []string
	for i := 0; i < e.next; i++ {
		if e.accepted[i] {
			runs = append(runs, i)
			texts = append(texts, e.units[i].text)
		}
	}
	if text != strings.Join(texts, "\n") {
		runs = append(runs, e.next-1)
	}
	pos := make(map[int]int)
	for p, i := range runs {
		if s := e.units[i].stmt; e.first[s] == i {
			pos[s] = p
		}
	}
	var b strings.Builder
	if e.autocommit {
		b.WriteString("-- txn-mode = off\n")
	}
	var at []int
	for s, written := range e.stmts {
		if p, ok := pos[s]; ok {
			b.WriteString(written)
			at = append(at, p)
		}
	}
	return b.String(), at
}

func (e *pgEngine) Report(t testing.TB, change string) *walktest.Finding {
	t.Helper()
	sql, at := e.replay(change)
	f := walkFinding(t, sql, e.target)
	if f == nil {
		return nil
	}
	return &walktest.Finding{Statement: at[f.Statement], Message: f.Message}
}

func (e *pgEngine) Violation(t testing.TB, change string) *walk.Violation {
	t.Helper()
	sql, _ := e.replay(change)
	return walkViolation(t, sql, e.target)
}

func (e *pgEngine) Codes(message string) []string { return walkStates(message) }

func (e *pgEngine) Holds(t testing.TB, v *walk.Violation, message string) (bool, string) {
	t.Helper()
	return verifyClaim(t, e.conn, v, message)
}

// judgeWalk runs the change on the server as Bytebase would and holds the
// walk to it (walktest.Judge). It returns how many rejected statements the
// walk found and how many there were.
func judgeWalk(t *testing.T, conn *sql.Conn, schema *metadata.DatabaseSchemaMetadata, user, change string, complete, tenant bool) (found, rejected int) {
	t.Helper()
	e := &pgEngine{conn: conn, target: review.Target{Schema: schema, SessionUser: user}, tenant: tenant}
	e.exec(t, "BEGIN")
	defer e.exec(t, "ROLLBACK")
	found, rejected = walktest.Judge(t, e, change, complete)
	if e.mismatch != "" {
		t.Errorf("the walk does not run the change as Bytebase does: %s", e.mismatch)
	}
	return found, rejected
}

// bbUnit is one statement as Bytebase sends it to the server.
type bbUnit struct {
	stmt int
	text string
	// apart is true for a unit Bytebase runs outside its transaction;
	// superuser for one its tenant mode runs as the session user.
	apart, superuser bool
}

// The patterns of Bytebase's PostgreSQL executor, copied from it.
var (
	bbExecDropDatabase = regexp.MustCompile(`(?i)DROP DATABASE`)
	bbExecCreateIndex  = regexp.MustCompile(`(?i)CREATE(\s+(UNIQUE\s+)?)INDEX(\s+)CONCURRENTLY`)
	bbExecDropIndex    = regexp.MustCompile(`(?i)DROP(\s+)INDEX(\s+)CONCURRENTLY`)
	bbExecVacuum       = regexp.MustCompile(`(?i)^\s*VACUUM`)
	bbExecSetRole      = regexp.MustCompile(`(?i)SET\s+((SESSION|LOCAL)\s+)?ROLE`)
	bbExecTxnMode      = regexp.MustCompile(`(?i)^\s*--\s*txn-mode\s*=\s*(on|off)\s*$`)
	bbExecIsolation    = regexp.MustCompile(`(?i)^\s*--\s*txn-isolation\s*=\s*(READ\s+UNCOMMITTED|READ\s+COMMITTED|REPEATABLE\s+READ|SERIALIZABLE)\s*$`)
)

// bytebaseExec lays a change out as Bytebase's PostgreSQL executor runs
// it, transcribed from its code independently of the walk's planRun: the
// mode its directives choose, the statements its splitter cuts, and the
// units it sends, in order. In transaction mode it runs the statements it
// does not set apart in its transaction, then, after COMMIT, the SET ROLE
// statements again and the ones PostgreSQL refuses in a transaction
// block; in auto-commit mode those first, then the rest. In tenant mode
// it runs the statements a superuser may need as the session user.
//
// bytebase: backend/plugin/db/pg/pg.go — Execute, IsNonTransactionStatement,
// isSuperuserStatement; backend/plugin/parser/base/transaction_mode.go —
// ParseTransactionConfig (main, 2026-10-28)
func bytebaseExec(change string, tenant bool, owner string) (autocommit bool, stmts []string, units []bbUnit) {
	// ParseTransactionConfig.
	lines := strings.Split(change, "\n")
	directive := make(map[int]bool)
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if !strings.HasPrefix(trimmed, "--") {
			break
		}
		if m := bbExecTxnMode.FindStringSubmatch(line); len(m) == 2 {
			autocommit = strings.ToLower(m[1]) == "off"
			directive[i] = true
			continue
		}
		if m := bbExecIsolation.FindStringSubmatch(line); len(m) == 2 {
			directive[i] = true
			continue
		}
	}
	if len(directive) > 0 {
		var remaining []string
		for i, line := range lines {
			if !directive[i] {
				remaining = append(remaining, line)
			}
		}
		change = strings.TrimSpace(strings.Join(remaining, "\n"))
	}
	// Execute: SplitSQL and FilterEmptyStatements, then the loop.
	for _, seg := range pg.Split(change) {
		if !seg.Empty() {
			stmts = append(stmts, seg.Text)
		}
	}
	var commands, apart []bbUnit
	for i, text := range stmts {
		u := bbUnit{stmt: i, text: text}
		switch {
		case bbExecSetRole.MatchString(text):
			apart = append(apart, bbUnit{stmt: i, text: text, apart: true})
		case bbExecDropDatabase.MatchString(text) || bbExecCreateIndex.MatchString(text) || bbExecDropIndex.MatchString(text) || bbExecVacuum.MatchString(text):
			apart = append(apart, bbUnit{stmt: i, text: text, apart: true})
			continue
		case bbExecSuperuser(text):
			u.superuser = true
			if tenant {
				ct := text
				if !strings.HasSuffix(strings.TrimRightFunc(ct, unicode.IsSpace), ";") {
					ct += ";"
				}
				u.text = fmt.Sprintf("SET LOCAL ROLE NONE;%sSET LOCAL ROLE '%s';", ct, owner)
			}
		}
		commands = append(commands, u)
	}
	if autocommit {
		return autocommit, stmts, append(apart, commands...)
	}
	return autocommit, stmts, append(commands, apart...)
}

// bbExecSuperuser is Bytebase's isSuperuserStatement.
func bbExecSuperuser(stmt string) bool {
	upper := strings.ToUpper(strings.TrimLeftFunc(stmt, unicode.IsSpace))
	return strings.HasPrefix(upper, "GRANT") || strings.HasPrefix(upper, "CREATE EXTENSION") || strings.HasPrefix(upper, "CREATE EVENT TRIGGER") || strings.HasPrefix(upper, "COMMENT ON EVENT TRIGGER") || strings.HasPrefix(upper, "COMMENT ON EXTENSION")
}

// walkViolation returns the walk's violation for a change, with the key
// and fact it claims.
func walkViolation(t testing.TB, sql string, target review.Target) *walk.Violation {
	t.Helper()
	stmts, syntax := parse(sql, statementRanges(sql))
	if syntax != nil {
		t.Fatalf("syntax: %s", syntax.Message)
	}
	v, err := walkThrough(context.Background(), stmts, planRun(sql, stmts), target)
	if err != nil {
		t.Fatal(err)
	}
	if v == nil {
		t.Fatalf("no violation for a change Review reports on:\n%s", sql)
	}
	return v
}

// relkinds maps pg_class.relkind to the walk's relation kinds.
var relkinds = map[string]string{
	"r": relTable, "p": relTable, "v": relView, "m": relMatview, "S": relSequence,
	"i": relIndex, "I": relIndex, "f": relForeign, "c": relComposite,
}

// verifyClaim checks the facts a finding claims against the server's
// catalog, as it is now, key by key (walktest.Claims). A claim about the
// SQL's own text holds by itself. Two claims rest on what the walker
// keeps beside the ledger, that a table is partitioned, which the catalog
// confirms too.
func verifyClaim(t testing.TB, conn *sql.Conn, v *walk.Violation, msg string) (bool, string) {
	t.Helper()
	for _, prefix := range []string{"cannot truncate only a partitioned table: ", "cannot create index on partitioned table "} {
		if !strings.HasPrefix(msg, prefix) {
			continue
		}
		name := strings.TrimSuffix(strings.TrimPrefix(msg, prefix), " concurrently")
		if !catalogAsk(t, conn, `SELECT EXISTS (SELECT 1 FROM pg_class WHERE oid = to_regclass($1) AND relkind = 'p')`, name) {
			return false, name + " is not partitioned"
		}
	}
	for _, c := range walktest.Claims(v) {
		if ok, why := claimHolds(t, conn, c); !ok {
			return false, why
		}
	}
	return true, ""
}

func catalogAsk(t testing.TB, conn *sql.Conn, q string, args ...any) bool {
	t.Helper()
	var ok bool
	if err := conn.QueryRowContext(context.Background(), q, args...).Scan(&ok); err != nil {
		t.Fatalf("%s: %v", q, err)
	}
	return ok
}

// claimHolds checks one fact against the catalog. An Unknown fact claims
// nothing.
func claimHolds(t testing.TB, conn *sql.Conn, c walktest.Claim) (bool, string) {
	t.Helper()
	ctx := context.Background()
	ask := func(q string, args ...any) bool {
		t.Helper()
		return catalogAsk(t, conn, q, args...)
	}
	k := c.Key
	if c.Fact == walk.Unknown {
		return true, ""
	}
	present := c.Fact == walk.Present
	names := walk.Names(k.Parent)
	schemaOID := `(CASE WHEN $1 = 'pg_temp' THEN pg_my_temp_schema() ELSE (SELECT oid FROM pg_namespace WHERE nspname = $1) END)`
	switch k.Space {
	case spaceSchema:
		return ask(`SELECT EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = $1) = $2`, k.Name, present), fmt.Sprintf("schema %s present=%v", k.Name, !present)
	case spaceRelation:
		var kind sql.NullString
		if err := conn.QueryRowContext(ctx, `SELECT relkind::text FROM pg_class WHERE relnamespace = `+schemaOID+` AND relname = $2`, names[0], k.Name).Scan(&kind); err != nil && !errors.Is(err, sql.ErrNoRows) {
			t.Fatal(err)
		}
		switch {
		case !present:
			return !kind.Valid, fmt.Sprintf("relation %s.%s exists", names[0], k.Name)
		case !kind.Valid:
			return false, fmt.Sprintf("relation %s.%s does not exist", names[0], k.Name)
		case c.Kind != "" && relkinds[kind.String] != c.Kind:
			return false, fmt.Sprintf("relation %s.%s is %s, not %s", names[0], k.Name, relkinds[kind.String], c.Kind)
		}
		return true, ""
	case spaceType:
		if present && c.Kind == typOther {
			return ask(`SELECT EXISTS (SELECT 1 FROM pg_type WHERE typnamespace = `+schemaOID+` AND typname = $2 AND typtype <> 'd')`, names[0], k.Name), fmt.Sprintf("type %s.%s is a domain or missing", names[0], k.Name)
		}
		return ask(`SELECT EXISTS (SELECT 1 FROM pg_type WHERE typnamespace = `+schemaOID+` AND typname = $2) = $3`, names[0], k.Name, present), fmt.Sprintf("type %s.%s present=%v", names[0], k.Name, !present)
	case spaceColumn, spaceConstraint, spacePrimaryKey:
		rel := `(SELECT c.oid FROM pg_class c WHERE c.relnamespace = ` + schemaOID + ` AND c.relname = $2)`
		var q string
		switch k.Space {
		case spaceColumn:
			q = `SELECT EXISTS (SELECT 1 FROM pg_attribute WHERE attrelid = ` + rel + ` AND attname = $3 AND attnum > 0 AND NOT attisdropped) = $4`
		case spaceConstraint:
			q = `SELECT EXISTS (SELECT 1 FROM pg_constraint WHERE conrelid = ` + rel + ` AND conname = $3) = $4`
		default:
			q = `SELECT EXISTS (SELECT 1 FROM pg_constraint WHERE conrelid = ` + rel + ` AND contype = 'p' AND $3 = '') = $4`
		}
		return ask(q, names[0], names[1], k.Name, present), fmt.Sprintf("%s %s of %s.%s present=%v", k.Space, k.Name, names[0], names[1], !present)
	}
	return false, "unknown space " + k.Space
}

func walkFinding(t testing.TB, sql string, target review.Target) *review.Finding {
	t.Helper()
	if strings.TrimSpace(sql) == "" {
		return nil
	}
	result, err := Review(context.Background(), sql, review.Options{Rules: []review.Rule{review.WalkThrough}}, []review.Target{target})
	if err != nil {
		t.Fatalf("Review: %v", err)
	}
	var out *review.Finding
	for i, f := range result.Findings {
		switch f.Rule {
		case review.Syntax:
			t.Fatalf("syntax: %s\n%s", f.Message, sql)
		case review.WalkThrough:
			if out != nil {
				t.Errorf("more than one finding: %q and %q", out.Message, f.Message)
			}
			out = &result.Findings[i]
		}
	}
	return out
}

// walkStates lists the SQLSTATEs a finding's message stands for.
func walkStates(msg string) []string {
	switch {
	case strings.HasPrefix(msg, "relation ") && strings.Contains(msg, " already exists"):
		return []string{"42P07"}
	case strings.HasPrefix(msg, "type ") && strings.HasSuffix(msg, " already exists"),
		strings.HasPrefix(msg, "type ") && strings.Contains(msg, " already exists in schema"):
		return []string{"42710", "42P07"}
	case strings.HasPrefix(msg, "schema ") && strings.HasSuffix(msg, " already exists"):
		return []string{"42P06"}
	case strings.HasPrefix(msg, "schema ") && strings.HasSuffix(msg, " does not exist"),
		msg == "no schema has been selected to create in":
		return []string{"3F000"}
	case strings.HasPrefix(msg, "index ") && strings.HasSuffix(msg, " does not exist"),
		strings.HasPrefix(msg, "type ") && strings.HasSuffix(msg, " does not exist"),
		strings.HasPrefix(msg, "constraint ") && strings.HasSuffix(msg, " does not exist"):
		return []string{"42704"}
	case strings.HasPrefix(msg, "column ") && strings.HasSuffix(msg, " does not exist"):
		return []string{"42703"}
	case strings.HasSuffix(msg, " does not exist"):
		return []string{"42P01"}
	case strings.Contains(msg, " is not a"), strings.HasPrefix(msg, "cannot create index on relation"), strings.HasPrefix(msg, "cannot open relation"),
		strings.HasSuffix(msg, " is a composite type"), strings.HasPrefix(msg, "cannot change schema of index"),
		strings.HasPrefix(msg, "cannot truncate only a partitioned table"):
		return []string{"42809"}
	case strings.HasPrefix(msg, "column ") && (strings.HasSuffix(msg, " already exists") || strings.HasSuffix(msg, " specified more than once")),
		strings.HasSuffix(msg, " conflicts with a system column name"):
		return []string{"42701"}
	case strings.HasPrefix(msg, "cannot drop system column"), strings.HasPrefix(msg, "cannot alter system column"), strings.HasPrefix(msg, "cannot rename system column"):
		return []string{"0A000"}
	case strings.HasPrefix(msg, "multiple primary keys"):
		return []string{"42P16"}
	case strings.HasPrefix(msg, "constraint ") && strings.HasSuffix(msg, " already exists"),
		strings.HasPrefix(msg, "check constraint ") && strings.HasSuffix(msg, " already exists"):
		return []string{"42710"}
	case strings.HasPrefix(msg, "cross-database references are not implemented"), strings.HasPrefix(msg, "DROP INDEX CONCURRENTLY does not support"),
		strings.HasPrefix(msg, "cannot create index on partitioned table "):
		return []string{"0A000"}
	case strings.HasPrefix(msg, "cannot drop schema"):
		return []string{"2BP01"}
	case strings.HasPrefix(msg, "unacceptable schema name"):
		return []string{"42939"}
	}
	return nil
}

// resetWalk empties the database and runs base.
func resetWalk(t *testing.T, conn *sql.Conn, base string) {
	t.Helper()
	ctx := context.Background()
	rows, err := conn.QueryContext(ctx, `SELECT nspname FROM pg_namespace WHERE nspname NOT LIKE 'pg\_%' AND nspname <> 'information_schema'`)
	if err != nil {
		t.Fatal(err)
	}
	var schemas []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			t.Fatal(err)
		}
		schemas = append(schemas, `"`+strings.ReplaceAll(s, `"`, `""`)+`"`)
	}
	rows.Close()
	stmts := []string{"RESET ALL", "RESET ROLE", "DISCARD TEMP"}
	rows, err = conn.QueryContext(ctx, `SELECT evtname FROM pg_event_trigger`)
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatal(err)
		}
		stmts = append(stmts, `DROP EVENT TRIGGER "`+name+`"`)
	}
	rows.Close()
	if len(schemas) > 0 {
		stmts = append(stmts, "DROP SCHEMA "+strings.Join(schemas, ", ")+" CASCADE")
	}
	stmts = append(stmts, "CREATE SCHEMA public")
	for _, s := range stmts {
		if _, err := conn.ExecContext(ctx, s); err != nil {
			t.Fatalf("%s: %v", s, err)
		}
	}
	if base != "" {
		if _, err := conn.ExecContext(ctx, base); err != nil {
			t.Fatalf("base: %v", err)
		}
	}
}

// walkSync reads the database into the metadata Bytebase's sync produces
// for what WalkThrough reads, and no more completely: the database's
// owner, schemas but the system ones, tables but partitions with their
// columns, indexes with their definitions, and constraints, partitions
// with their indexes under the table Bytebase finds them under,
// views, materialized views with their indexes, sequences, foreign
// tables, enum and composite types, routines, table triggers, event
// triggers, and the search path. It drops what the sync drops by name:
// every pg_ name stands for its system tables and views, which all begin
// with pg_, and the walk trusts no pg_ name's absence anyway; the view
// and materialized view prefixes of IsSystemView; the functions of
// builtinCalls, all of them on its list of system functions, and the
// extension prefixes of IsSystemFunction.
//
// bytebase: backend/plugin/db/pg/sync.go — SyncDBSchema, getTablePartitions,
// warpTablePartitions; backend/plugin/parser/pg/system_objects.go
func walkSync(t *testing.T, conn *sql.Conn) *metadata.DatabaseSchemaMetadata {
	t.Helper()
	ctx := context.Background()
	db := &metadata.DatabaseSchemaMetadata{Name: "omni_review"}
	if err := conn.QueryRowContext(ctx, "SHOW search_path").Scan(&db.SearchPath); err != nil {
		t.Fatal(err)
	}
	if err := conn.QueryRowContext(ctx, "SELECT pg_get_userbyid(datdba) FROM pg_database WHERE datname = current_database()").Scan(&db.Owner); err != nil {
		t.Fatal(err)
	}
	// The sync reads in a transaction on an empty search path, so the
	// definitions it reads name every schema.
	for _, q := range []string{"BEGIN", "SELECT set_config('search_path', '', true)"} {
		if _, err := conn.ExecContext(ctx, q); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
	}
	defer func() {
		if _, err := conn.ExecContext(ctx, "COMMIT"); err != nil {
			t.Fatal(err)
		}
	}()
	query := func(q string, scan func(*sql.Rows) error, args ...any) {
		t.Helper()
		rows, err := conn.QueryContext(ctx, q, args...)
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
	systemView := func(name string) bool {
		for _, prefix := range []string{"pg_", "g_columnar_", "google_db_advisor_", "g_agg_stat_", "hypopg"} {
			if strings.HasPrefix(name, prefix) {
				return true
			}
		}
		return false
	}
	const user = `n.nspname NOT LIKE 'pg\_%' AND n.nspname NOT IN ('information_schema', 'rw_catalog',
		'timescaledb_information', 'timescaledb_experimental', '_timescaledb_cache', '_timescaledb_catalog', '_timescaledb_internal', '_timescaledb_config',
		'squeeze', 'citus', 'citus_internal', 'columnar', 'columnar_internal', 'cron', 'partman')`
	schemas := make(map[string]*metadata.SchemaMetadata)
	query(`SELECT n.nspname FROM pg_namespace n WHERE `+user+` ORDER BY 1`, func(rows *sql.Rows) error {
		s := &metadata.SchemaMetadata{}
		if err := rows.Scan(&s.Name); err != nil {
			return err
		}
		schemas[s.Name] = s
		db.Schemas = append(db.Schemas, s)
		return nil
	})
	columnsOf := func(oid int64) []*metadata.ColumnMetadata {
		var cols []*metadata.ColumnMetadata
		query(`SELECT attname FROM pg_attribute WHERE attrelid = $1 AND attnum > 0 AND NOT attisdropped ORDER BY attnum`, func(rows *sql.Rows) error {
			c := &metadata.ColumnMetadata{}
			if err := rows.Scan(&c.Name); err != nil {
				return err
			}
			cols = append(cols, c)
			return nil
		}, oid)
		return cols
	}
	indexesOf := func(oid int64) []*metadata.IndexMetadata {
		var out []*metadata.IndexMetadata
		query(`SELECT i.relname, x.indisunique, x.indisprimary, EXISTS (SELECT 1 FROM pg_constraint c WHERE c.conindid = x.indexrelid AND c.contype IN ('p', 'u', 'x')), pg_get_indexdef(i.oid) || ';'
			FROM pg_index x JOIN pg_class i ON i.oid = x.indexrelid WHERE x.indrelid = $1 ORDER BY 1`, func(rows *sql.Rows) error {
			i := &metadata.IndexMetadata{}
			if err := rows.Scan(&i.Name, &i.Unique, &i.Primary, &i.IsConstraint, &i.Definition); err != nil {
				return err
			}
			out = append(out, i)
			return nil
		}, oid)
		return out
	}
	// Bytebase keys a partition by its parent's own schema and name, and
	// looks the partitions of a partition up under the schema of the table
	// they descend from (warpTablePartitions), each key once.
	type parentKey struct{ schema, name string }
	partitions := make(map[parentKey][]*metadata.TablePartitionMetadata)
	type part struct {
		oid                          int64
		name, parentSchema, parentOf string
	}
	var parts []part
	query(`SELECT c.oid, c.relname, pn.nspname, p.relname FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
		JOIN pg_inherits i ON i.inhrelid = c.oid JOIN pg_class p ON p.oid = i.inhparent JOIN pg_namespace pn ON pn.oid = p.relnamespace
		WHERE c.relkind IN ('r', 'f', 'p') AND c.relispartition AND `+user+` ORDER BY c.oid`, func(rows *sql.Rows) error {
		var p part
		if err := rows.Scan(&p.oid, &p.name, &p.parentSchema, &p.parentOf); err != nil {
			return err
		}
		parts = append(parts, p)
		return nil
	})
	for _, p := range parts {
		if strings.HasPrefix(p.name, "pg_") || strings.HasPrefix(p.parentOf, "pg_") {
			continue
		}
		k := parentKey{p.parentSchema, p.parentOf}
		partitions[k] = append(partitions[k], &metadata.TablePartitionMetadata{Name: p.name, Indexes: indexesOf(p.oid)})
	}
	var warp func(schema, name string) []*metadata.TablePartitionMetadata
	warp = func(schema, name string) []*metadata.TablePartitionMetadata {
		k := parentKey{schema, name}
		out, ok := partitions[k]
		if !ok {
			return nil
		}
		delete(partitions, k)
		for _, p := range out {
			p.Subpartitions = warp(schema, p.Name)
		}
		return out
	}
	type rel struct {
		oid    int64
		schema string
		name   string
		kind   string
	}
	var rels []rel
	query(`SELECT c.oid, n.nspname, c.relname, c.relkind::text FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'f', 'c') AND NOT c.relispartition AND `+user+` ORDER BY 2, 3`, func(rows *sql.Rows) error {
		var r rel
		if err := rows.Scan(&r.oid, &r.schema, &r.name, &r.kind); err != nil {
			return err
		}
		rels = append(rels, r)
		return nil
	})
	for _, r := range rels {
		s := schemas[r.schema]
		switch r.kind {
		case "r", "p":
			if strings.HasPrefix(r.name, "pg_") {
				continue
			}
			tb := &metadata.TableMetadata{Name: r.name, Columns: columnsOf(r.oid), Indexes: indexesOf(r.oid), Partitions: warp(r.schema, r.name)}
			query(`SELECT conname, contype::text FROM pg_constraint WHERE conrelid = $1 AND contype IN ('f', 'c', 'x') ORDER BY 1`, func(rows *sql.Rows) error {
				var name, kind string
				if err := rows.Scan(&name, &kind); err != nil {
					return err
				}
				switch kind {
				case "f":
					tb.ForeignKeys = append(tb.ForeignKeys, &metadata.ForeignKeyMetadata{Name: name})
				case "c":
					tb.CheckConstraints = append(tb.CheckConstraints, &metadata.CheckConstraintMetadata{Name: name})
				case "x":
					tb.ExcludeConstraints = append(tb.ExcludeConstraints, &metadata.ExcludeConstraintMetadata{Name: name})
				}
				return nil
			}, r.oid)
			query(`SELECT tgname FROM pg_trigger WHERE tgrelid = $1 AND NOT tgisinternal`, func(rows *sql.Rows) error {
				tr := &metadata.TriggerMetadata{}
				if err := rows.Scan(&tr.Name); err != nil {
					return err
				}
				tb.Triggers = append(tb.Triggers, tr)
				return nil
			}, r.oid)
			s.Tables = append(s.Tables, tb)
		case "v":
			if systemView(r.name) {
				continue
			}
			s.Views = append(s.Views, &metadata.ViewMetadata{Name: r.name, Columns: columnsOf(r.oid)})
		case "m":
			if systemView(r.name) {
				continue
			}
			s.MaterializedViews = append(s.MaterializedViews, &metadata.MaterializedViewMetadata{Name: r.name, Indexes: indexesOf(r.oid)})
		case "S":
			s.Sequences = append(s.Sequences, &metadata.SequenceMetadata{Name: r.name})
		case "f":
			s.ExternalTables = append(s.ExternalTables, &metadata.ExternalTableMetadata{Name: r.name, Columns: columnsOf(r.oid)})
		case "c":
			s.CompositeTypes = append(s.CompositeTypes, &metadata.CompositeTypeMetadata{Name: r.name})
		}
	}
	query(`SELECT n.nspname, t.typname FROM pg_type t JOIN pg_namespace n ON n.oid = t.typnamespace WHERE t.typtype = 'e' AND `+user+` ORDER BY 1, 2`, func(rows *sql.Rows) error {
		var schema string
		e := &metadata.EnumTypeMetadata{}
		if err := rows.Scan(&schema, &e.Name); err != nil {
			return err
		}
		schemas[schema].EnumTypes = append(schemas[schema].EnumTypes, e)
		return nil
	})
	query(`SELECT n.nspname, p.proname, p.prokind::text FROM pg_proc p JOIN pg_namespace n ON n.oid = p.pronamespace WHERE `+user+` ORDER BY 1, 2`, func(rows *sql.Rows) error {
		var schema, name, kind string
		if err := rows.Scan(&schema, &name, &kind); err != nil {
			return err
		}
		if builtinCalls[strings.ToLower(name)] {
			return nil
		}
		for _, prefix := range []string{"g_columnar_", "google_columnar_", "google_db_advisor_", "g_agg_stat_", "hypopg"} {
			if strings.HasPrefix(name, prefix) {
				return nil
			}
		}
		if kind == "p" {
			schemas[schema].Procedures = append(schemas[schema].Procedures, &metadata.ProcedureMetadata{Name: name})
		} else {
			schemas[schema].Functions = append(schemas[schema].Functions, &metadata.FunctionMetadata{Name: name})
		}
		return nil
	})
	query(`SELECT evtname, evtevent, evtenabled <> 'D', coalesce(array_to_json(evttags)::text, '[]') FROM pg_event_trigger ORDER BY 1`, func(rows *sql.Rows) error {
		e := &metadata.EventTriggerMetadata{}
		var tags string
		if err := rows.Scan(&e.Name, &e.Event, &e.Enabled, &tags); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(tags), &e.Tags); err != nil {
			return err
		}
		db.EventTriggers = append(db.EventTriggers, e)
		return nil
	})
	return db
}

// sweepChange writes a random change of n statements over the sweep's
// name pool.
func sweepChange(rng *rand.Rand, n int) string {
	pick := func(xs ...string) string { return xs[rng.Intn(len(xs))] }
	maybe := func(s string) string {
		if rng.Intn(3) == 0 {
			return s
		}
		return ""
	}
	schema := func() string { return pick("", "", "s1.", "s2.") }
	table := func() string { return schema() + pick("a", "b", "c", "a_pkey", "a_id_seq", "b_x_key", "v", "i") }
	col := func() string { return pick("id", "x", "y", "z", "xmin") }
	idx := func() string { return pick("i", "a_x_idx", "a", "b_x_key", "j") }
	con := func() string { return pick("k", "a_pkey", "b_x_key", "chk", "a_x_check") }
	var b strings.Builder
	for range n {
		var stmt string
		switch rng.Intn(50) {
		case 0, 1:
			stmt = fmt.Sprintf("CREATE TABLE %s%s (id serial PRIMARY KEY, x int%s, y text)", maybe("IF NOT EXISTS "), table(), maybe(" UNIQUE"))
		case 2:
			stmt = fmt.Sprintf("CREATE TABLE %s (x int, y int%s)", table(), maybe(", x text"))
		case 3:
			stmt = fmt.Sprintf("CREATE TABLE %s (LIKE %s INCLUDING ALL)", table(), table())
		case 4:
			stmt = fmt.Sprintf("CREATE TABLE %s%s AS SELECT 1 AS x", maybe("IF NOT EXISTS "), table())
		case 5, 6:
			stmt = fmt.Sprintf("DROP TABLE %s%s%s", maybe("IF EXISTS "), table(), maybe(" CASCADE"))
		case 7:
			stmt = fmt.Sprintf("DROP %s %s%s", pick("VIEW", "INDEX", "SEQUENCE", "MATERIALIZED VIEW"), maybe("IF EXISTS "), table())
		case 8, 9:
			stmt = fmt.Sprintf("ALTER TABLE %s%s ADD COLUMN %s%s int", maybe("IF EXISTS "), table(), maybe("IF NOT EXISTS "), col())
		case 10:
			stmt = fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s%s%s", table(), maybe("IF EXISTS "), col(), maybe(" CASCADE"))
		case 11:
			stmt = fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s", table(), col(), col())
		case 12:
			stmt = fmt.Sprintf("ALTER TABLE %s RENAME TO %s", table(), pick("a", "b", "c", "a_pkey", "v"))
		case 13:
			stmt = fmt.Sprintf("ALTER TABLE %s SET SCHEMA %s", table(), pick("public", "s1", "s2", "s3"))
		case 14:
			stmt = fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s %s", table(), con(), pick("UNIQUE (x)", "CHECK (x > 0)", "PRIMARY KEY (x)"))
		case 15:
			stmt = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s%s%s", table(), maybe("IF EXISTS "), con(), maybe(" CASCADE"))
		case 16:
			stmt = fmt.Sprintf("ALTER TABLE %s RENAME CONSTRAINT %s TO %s", table(), con(), con())
		case 17, 18:
			name := pick(idx(), idx(), "")
			if name != "" {
				name = maybe("IF NOT EXISTS ") + name
			}
			stmt = fmt.Sprintf("CREATE %sINDEX %s ON %s (%s)", maybe("UNIQUE "), name, table(), pick("x", "y", "id"))
		case 19:
			stmt = fmt.Sprintf("ALTER INDEX %s%s RENAME TO %s", maybe("IF EXISTS "), schema()+idx(), idx())
		case 20:
			stmt = fmt.Sprintf("CREATE %sVIEW %s AS SELECT 1 AS x", maybe("OR REPLACE "), table())
		case 21:
			stmt = fmt.Sprintf("CREATE SEQUENCE %s%s", maybe("IF NOT EXISTS "), table())
		case 22:
			stmt = fmt.Sprintf("CREATE TYPE %s AS ENUM ('a')", table())
		case 23:
			stmt = fmt.Sprintf("DROP %s %s%s%s", pick("TYPE", "TYPE", "DOMAIN"), maybe("IF EXISTS "), pick("", "public.", "s1.")+pick("a", "e", "v", "i"), maybe(" CASCADE"))
		case 24:
			stmt = fmt.Sprintf("CREATE SCHEMA %s%s", maybe("IF NOT EXISTS "), pick("s1", "s2", "s3"))
		case 25:
			stmt = fmt.Sprintf("DROP SCHEMA %s%s%s", maybe("IF EXISTS "), pick("s1", "s2", "s3"), maybe(" CASCADE"))
		case 26:
			stmt = fmt.Sprintf("ALTER SCHEMA %s RENAME TO %s", pick("s1", "s2", "s3"), pick("s1", "s2", "s3"))
		case 27:
			stmt = pick("SET search_path = s1, public", "SET search_path = public", "RESET search_path", "SET search_path = s2")
		case 28:
			stmt = fmt.Sprintf("CREATE TEMP TABLE %s (x int)", pick("a", "b", "c"))
		case 29:
			stmt = pick(fmt.Sprintf("COMMENT ON TABLE %s IS 'c'", table()), fmt.Sprintf("COMMENT ON COLUMN %s.%s IS 'c'", table(), col()))
		case 30:
			stmt = fmt.Sprintf("INSERT INTO %s DEFAULT VALUES", table())
		case 31:
			stmt = fmt.Sprintf("TRUNCATE %s%s", maybe("ONLY "), table())
		case 32:
			stmt = pick("CREATE TABLE pt (x int) PARTITION BY RANGE (x)", "CREATE TABLE pt_1 PARTITION OF pt FOR VALUES FROM (0) TO (10)", "CREATE TABLE s1.pt_2 PARTITION OF pt FOR VALUES FROM (10) TO (20)", "DROP TABLE pt", "DROP TABLE pt_1", "CREATE TABLE pt_1 (x int)", "ALTER TABLE pt_1 ADD COLUMN y int")
		case 33:
			stmt = fmt.Sprintf("CREATE TABLE %s (z int) INHERITS (%s)", table(), table())
		case 34:
			stmt = fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s %s", table(), col(), pick("TYPE text", "SET NOT NULL", "SET DEFAULT 1", "DROP NOT NULL"))
		case 35:
			stmt = pick(fmt.Sprintf("GRANT SELECT ON %s TO PUBLIC", table()), fmt.Sprintf("LOCK TABLE %s", table()))
		case 36:
			c := col()
			stmt = pick(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s int, DROP COLUMN %s", table(), col(), col()),
				fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s%s int, ADD COLUMN %s%s int%s", table(), maybe("IF NOT EXISTS "), c, maybe("IF NOT EXISTS "), c, maybe(" CONSTRAINT k UNIQUE")),
				fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s%s, DROP COLUMN %s%s", table(), maybe("IF EXISTS "), c, maybe("IF EXISTS "), c),
				fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s%s, ADD COLUMN IF NOT EXISTS %s int%s", table(), maybe("IF EXISTS "), c, c, maybe(" CONSTRAINT k UNIQUE")))
		case 37:
			stmt = pick(fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (x) REFERENCES %s (id)", table(), con(), table()),
				fmt.Sprintf("ALTER TABLE %s ADD %sUNIQUE USING INDEX %s", table(), pick("", "CONSTRAINT "+con()+" ", "CONSTRAINT "+idx()+" "), idx()))
		case 38:
			stmt = fmt.Sprintf("CREATE TABLE %s (id int CONSTRAINT %s PRIMARY KEY, x int CONSTRAINT %s UNIQUE, y int REFERENCES %s)", table(), con(), con(), table())
		case 39:
			stmt = fmt.Sprintf("CREATE VIEW %s AS SELECT x FROM %s", table(), table())
		case 40:
			stmt = fmt.Sprintf("CREATE MATERIALIZED VIEW %s AS SELECT 1 AS x", table())
		case 41:
			stmt = fmt.Sprintf("CREATE TABLE %s (id int GENERATED ALWAYS AS IDENTITY, x int)", table())
		case 42:
			stmt = fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s ADD GENERATED ALWAYS AS IDENTITY", table(), col())
		case 43:
			stmt = pick(fmt.Sprintf("CREATE TYPE %s AS (a int)", table()), fmt.Sprintf("CREATE DOMAIN %s AS int", table()))
		case 44:
			stmt = fmt.Sprintf("CREATE INDEX %s ON %s ((x + 1))", pick("", idx()), table())
		case 45:
			stmt = fmt.Sprintf("CREATE TEMP TABLE %s (x int) ON COMMIT DROP", pick("a", "b", "j"))
		case 46:
			stmt = pick(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (x int) PARTITION BY RANGE (x)", table()),
				fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s PARTITION OF pt FOR VALUES FROM (20) TO (30)", table()),
				fmt.Sprintf("CREATE INDEX %s%s ON %s (x)", maybe("CONCURRENTLY "), idx(), table()))
		case 47:
			kinds := []string{"UNIQUE (x)", "CHECK (x > 0)", "PRIMARY KEY (x)", "FOREIGN KEY (x) REFERENCES a (id)"}
			stmt = pick(fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s %s, ADD CONSTRAINT %s %s", table(), con(), pick(kinds...), con(), pick(kinds...)),
				fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s%s int CONSTRAINT %s %s", table(), maybe("IF NOT EXISTS "), col(), con(), pick("UNIQUE", "PRIMARY KEY", "CHECK (true)")),
				fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s %s, DROP CONSTRAINT %s%s", table(), con(), pick(kinds...), maybe("IF EXISTS "), con()))
		case 48:
			n := pick("a", "b", "c")
			stmt = pick(fmt.Sprintf("CREATE TABLE %s%s (%s int, y int)", maybe("IF NOT EXISTS "), table(), pick("ctid", "tableoid", "cmin", "x")),
				fmt.Sprintf("CREATE TABLE %s%s (id int PRIMARY KEY, r int REFERENCES %s (id))", schema(), n, n))
		case 49:
			obj := func() string { return pick("TABLE", "VIEW", "MATERIALIZED VIEW", "SEQUENCE") }
			stmt = pick(fmt.Sprintf("ALTER %s %s%s RENAME TO %s", pick(obj(), "INDEX"), maybe("IF EXISTS "), table(), pick("a", "b", "v", "i")),
				fmt.Sprintf("ALTER %s %s%s SET SCHEMA %s", obj(), maybe("IF EXISTS "), table(), pick("public", "s1", "s2")),
				fmt.Sprintf("ALTER %s %s OWNER TO postgres", obj(), table()),
				fmt.Sprintf("ALTER INDEX %s SET TABLESPACE pg_default", table()),
				fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s%s, DROP CONSTRAINT %s%s", table(), maybe("IF EXISTS "), con(), maybe("IF EXISTS "), con()),
				fmt.Sprintf("DROP INDEX CONCURRENTLY %s%s%s", maybe("IF EXISTS "), table(), pick("", ", "+idx(), " CASCADE")))
		}
		b.WriteString(stmt)
		b.WriteString(";\n")
	}
	return b.String()
}
