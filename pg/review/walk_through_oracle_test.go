package review

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// walkBase is the schema the curated WalkThrough cases start from: one of
// every relation kind, keys and constraints, a view reading a column, a
// partitioned table with a partition in another schema, inheritance, and
// an enum type.
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
}

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
	{name: "a role's schema", change: "CREATE SCHEMA AUTHORIZATION bob;\nSET ROLE bob;\nCREATE TABLE x (a int);\nRESET ROLE;\nCREATE TABLE bob.x (a int);\nCREATE TABLE x (a int);", complete: true},

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
	if _, err := conn.ExecContext(ctx, "DO $$ BEGIN CREATE ROLE bob; EXCEPTION WHEN duplicate_object THEN NULL; END $$"); err != nil {
		t.Fatal(err)
	}
	var found, rejected int
	for _, tt := range walkCases {
		t.Run(tt.name, func(t *testing.T) {
			resetWalk(t, conn, walkBase+tt.extra)
			schema := walkSync(t, conn)
			user := tt.user
			if user == "" {
				user = "postgres"
			}
			f, r := judgeWalk(t, conn, schema, user, tt.change, tt.complete)
			found += f
			rejected += r
		})
	}
	t.Logf("found %d of the %d statements the server rejects", found, rejected)
}

// TestWalkThroughSweep runs random changes over a small pool of names, so
// creations, drops, renames, and generated names collide often, on
// random starting schemas. Precision must hold on every one; recall is
// logged. WALK_SWEEP_N sets the number of changes, WALK_SWEEP_SEED the
// seed, and WALK_SWEEP_IMAGE the PostgreSQL image.
func TestWalkThroughSweep(t *testing.T) {
	n := 40
	if v := os.Getenv("WALK_SWEEP_N"); v != "" {
		var err error
		if n, err = strconv.Atoi(v); err != nil {
			t.Fatalf("WALK_SWEEP_N: %v", err)
		}
	}
	seed := int64(20261019)
	if v := os.Getenv("WALK_SWEEP_SEED"); v != "" {
		seed, _ = strconv.ParseInt(v, 10, 64)
	}
	image := "postgres:17-alpine"
	if v := os.Getenv("WALK_SWEEP_IMAGE"); v != "" {
		image = v
	}
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
	rng := rand.New(rand.NewSource(seed))
	var found, rejected int
	for i := range n {
		base := sweepChange(rng, 8+rng.Intn(8))
		change := sweepChange(rng, 15+rng.Intn(15))
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
			f, r := judgeWalk(t, conn, schema, "postgres", change, false)
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

// outcome is how the server took one statement: accepted, or rejected
// with a SQLSTATE.
type outcome struct {
	text  string
	state string
	msg   string
}

// runChange runs a change on the server statement by statement, each in a
// savepoint of one transaction it rolls back at the end, so a rejected
// statement does nothing and the next one runs as the change would after
// it if it had been left out. onReject sees each rejected statement with
// the server as it was right before it, and the statements accepted
// before it.
func runChange(t *testing.T, conn *sql.Conn, change string, onReject func(accepted []string, o outcome)) []outcome {
	t.Helper()
	ctx := context.Background()
	exec := func(q string) {
		t.Helper()
		if _, err := conn.ExecContext(ctx, q); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
	}
	exec("BEGIN")
	defer exec("ROLLBACK")
	var out []outcome
	var accepted []string
	for _, stmt := range splitStatements(change) {
		o := outcome{text: stmt}
		upper := strings.ToUpper(stmt)
		if strings.HasPrefix(upper, "BEGIN") || strings.HasPrefix(upper, "COMMIT") {
			out = append(out, o)
			accepted = append(accepted, stmt)
			continue
		}
		exec("SAVEPOINT walk_oracle")
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			var pgErr *pgconn.PgError
			if !errors.As(err, &pgErr) {
				t.Fatalf("%s: %v", stmt, err)
			}
			o.state, o.msg = pgErr.Code, pgErr.Message
			// 42601 is not only a syntax error: PostgreSQL raises it for
			// some semantic refusals too (SET DEFAULT on an identity
			// column). Text the parser refuses shows as a Syntax finding.
			exec("ROLLBACK TO SAVEPOINT walk_oracle")
			onReject(accepted, o)
		} else {
			accepted = append(accepted, stmt)
		}
		exec("RELEASE SAVEPOINT walk_oracle")
		out = append(out, o)
	}
	return out
}

// judgeWalk runs the change on the server and holds the walk to it. For
// each statement the server rejects, the walk over the statements it
// accepted before it, then that statement, may report nothing or that
// statement. What it reports must be true: either the server rejects the
// statement with the SQLSTATE the message stands for, or the server's
// catalog, as it was right before the statement, shows the fact the
// finding claims (the statement then fails for a reason PostgreSQL checks
// first, and for this one too). The walk over the statements the server
// accepts must report nothing. With complete, every rejected statement
// must be found. It returns how many rejected statements the walk found
// and how many there were.
func judgeWalk(t *testing.T, conn *sql.Conn, schema *metadata.DatabaseSchemaMetadata, user, change string, complete bool) (found, rejected int) {
	t.Helper()
	target := review.Target{Schema: schema, SessionUser: user}
	outcomes := runChange(t, conn, change, func(accepted []string, o outcome) {
		rejected++
		text := strings.Join(append(append([]string{}, accepted...), o.text), "\n")
		f := walkFinding(t, text, target)
		switch {
		case f == nil:
			if complete {
				t.Errorf("missed the rejection of %q (%s: %s)", o.text, o.state, o.msg)
			}
			return
		case f.Statement != len(accepted):
			t.Errorf("reported %q on statement %d, which the server accepts; it rejects %q (%s: %s)", f.Message, f.Statement, o.text, o.state, o.msg)
			return
		}
		found++
		if containsString(walkStates(f.Message), o.state) {
			return
		}
		v := walkViolation(t, text, target)
		if ok, why := verifyClaim(t, conn, v, f.Message); !ok {
			t.Errorf("reported %q on %q, which the server rejects with %s (%s), and the claim does not hold: %s", f.Message, o.text, o.state, o.msg, why)
		}
	})
	var accepted []string
	for _, o := range outcomes {
		if o.state == "" {
			accepted = append(accepted, o.text)
		}
	}
	if f := walkFinding(t, strings.Join(accepted, "\n"), target); f != nil {
		t.Errorf("reported %q on statement %d of a change the server runs:\n%s", f.Message, f.Statement, strings.Join(accepted, "\n"))
	}
	return found, rejected
}

// walkViolation returns the walk's violation for a change, with the key
// and fact it claims.
func walkViolation(t *testing.T, sql string, target review.Target) *walk.Violation {
	t.Helper()
	stmts, syntax := parse(sql, statementRanges(sql))
	if syntax != nil {
		t.Fatalf("syntax: %s", syntax.Message)
	}
	v, err := walkThrough(context.Background(), stmts, target)
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

// verifyClaim checks the fact a finding claims against the server's
// catalog, as it is now. A claim about the SQL's own text holds by
// itself.
func verifyClaim(t *testing.T, conn *sql.Conn, v *walk.Violation, msg string) (bool, string) {
	t.Helper()
	ctx := context.Background()
	ask := func(q string, args ...any) bool {
		t.Helper()
		var ok bool
		if err := conn.QueryRowContext(ctx, q, args...).Scan(&ok); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
		return ok
	}
	k := v.Check.Key
	if k == (walk.Key{}) {
		switch {
		case strings.HasPrefix(msg, "schema ") && strings.HasSuffix(msg, " does not exist"):
			name := strings.TrimSuffix(strings.TrimPrefix(msg, "schema "), " does not exist")
			return ask(`SELECT to_regnamespace($1) IS NULL`, name), "schema " + name + " exists"
		case strings.Contains(msg, " of relation "):
			// A column or constraint one ALTER drops twice: a claim
			// about the text.
			return true, ""
		case strings.HasSuffix(msg, " does not exist"):
			words := strings.TrimSuffix(msg, " does not exist")
			name := words[strings.LastIndex(words, " ")+1:]
			return ask(`SELECT to_regclass($1) IS NULL`, name), name + " exists"
		case msg == "no schema has been selected to create in":
			return ask(`SELECT current_schema() IS NULL`), "the search path has a schema"
		case strings.HasPrefix(msg, "cannot drop schema "):
			name := strings.TrimSuffix(strings.TrimPrefix(msg, "cannot drop schema "), " because other objects depend on it")
			return ask(`SELECT EXISTS (SELECT 1 FROM pg_class WHERE relnamespace = to_regnamespace($1)) OR EXISTS (SELECT 1 FROM pg_type WHERE typnamespace = to_regnamespace($1))`, name), "schema " + name + " is empty"
		case strings.HasPrefix(msg, "cannot truncate only a partitioned table: "):
			name := strings.TrimPrefix(msg, "cannot truncate only a partitioned table: ")
			return ask(`SELECT EXISTS (SELECT 1 FROM pg_class WHERE oid = to_regclass($1) AND relkind = 'p')`, name), name + " is not partitioned"
		case strings.HasPrefix(msg, "cannot create index on partitioned table "):
			// The oracle runs in a transaction, where CONCURRENTLY fails
			// first for being in one.
			name := strings.TrimSuffix(strings.TrimPrefix(msg, "cannot create index on partitioned table "), " concurrently")
			return ask(`SELECT EXISTS (SELECT 1 FROM pg_class WHERE oid = to_regclass($1) AND relkind = 'p')`, name), name + " is not partitioned"
		}
		return true, ""
	}
	present := v.Got == walk.Present
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
		case v.Kind != "" && relkinds[kind.String] != v.Kind:
			return false, fmt.Sprintf("relation %s.%s is %s, not %s", names[0], k.Name, relkinds[kind.String], v.Kind)
		}
		return true, ""
	case spaceType:
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

func walkFinding(t *testing.T, sql string, target review.Target) *review.Finding {
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
	case strings.HasPrefix(msg, "cross-database references are not implemented"), strings.HasPrefix(msg, "DROP INDEX CONCURRENTLY does not support"):
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
// for what WalkThrough reads, and no more completely: schemas but the
// system ones, tables but partitions with their columns, indexes, and
// constraints, partitions under their parent by name with their indexes,
// views, materialized views with their indexes, sequences, foreign
// tables, enum and composite types, routines, table triggers, event
// triggers, and the search path.
//
// bytebase: backend/plugin/db/pg/sync.go
func walkSync(t *testing.T, conn *sql.Conn) *metadata.DatabaseSchemaMetadata {
	t.Helper()
	ctx := context.Background()
	db := &metadata.DatabaseSchemaMetadata{Name: "omni_review"}
	if err := conn.QueryRowContext(ctx, "SHOW search_path").Scan(&db.SearchPath); err != nil {
		t.Fatal(err)
	}
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
	const user = `n.nspname NOT LIKE 'pg\_%' AND n.nspname <> 'information_schema'`
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
		query(`SELECT i.relname, x.indisunique, x.indisprimary, EXISTS (SELECT 1 FROM pg_constraint c WHERE c.conindid = x.indexrelid AND c.contype IN ('p', 'u', 'x'))
			FROM pg_index x JOIN pg_class i ON i.oid = x.indexrelid WHERE x.indrelid = $1 ORDER BY 1`, func(rows *sql.Rows) error {
			i := &metadata.IndexMetadata{}
			if err := rows.Scan(&i.Name, &i.Unique, &i.Primary, &i.IsConstraint); err != nil {
				return err
			}
			out = append(out, i)
			return nil
		}, oid)
		return out
	}
	var partitionsOf func(oid int64) []*metadata.TablePartitionMetadata
	partitionsOf = func(oid int64) []*metadata.TablePartitionMetadata {
		type child struct {
			oid  int64
			name string
		}
		var children []child
		query(`SELECT c.oid, c.relname FROM pg_inherits i JOIN pg_class c ON c.oid = i.inhrelid WHERE i.inhparent = $1 AND c.relispartition ORDER BY 2`, func(rows *sql.Rows) error {
			var c child
			if err := rows.Scan(&c.oid, &c.name); err != nil {
				return err
			}
			children = append(children, c)
			return nil
		}, oid)
		var out []*metadata.TablePartitionMetadata
		for _, c := range children {
			out = append(out, &metadata.TablePartitionMetadata{Name: c.name, Indexes: indexesOf(c.oid), Subpartitions: partitionsOf(c.oid)})
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
			tb := &metadata.TableMetadata{Name: r.name, Columns: columnsOf(r.oid), Indexes: indexesOf(r.oid), Partitions: partitionsOf(r.oid)}
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
			s.Views = append(s.Views, &metadata.ViewMetadata{Name: r.name, Columns: columnsOf(r.oid)})
		case "m":
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
			stmt = fmt.Sprintf("DROP TYPE %s%s%s", maybe("IF EXISTS "), pick("public.", "s1.")+pick("a", "e", "v"), maybe(" CASCADE"))
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
				fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s%s int, ADD COLUMN %s%s int", table(), maybe("IF NOT EXISTS "), c, maybe("IF NOT EXISTS "), c),
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
