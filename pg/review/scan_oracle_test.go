package review

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
	"github.com/bytebase/omni/review"
)

// oracleBase is the schema every oracle case starts from: keys of every
// kind, a reference into p, a table without a key, a second schema with a
// table of the same name, and a schema named after the session user.
const oracleBase = `
CREATE TABLE p (id integer PRIMARY KEY);
CREATE TABLE t (id integer PRIMARY KEY, code text CONSTRAINT t_code_key UNIQUE, p_id integer CONSTRAINT t_p_fk REFERENCES p (id), n integer CONSTRAINT t_n_check CHECK (n > 0));
CREATE TABLE r (t_id integer CONSTRAINT r_t_fk REFERENCES t (id));
CREATE TABLE nokey (a integer);
CREATE SCHEMA s;
CREATE TABLE s.t (id integer CONSTRAINT t_pkey UNIQUE);
CREATE SCHEMA postgres;
CREATE TABLE postgres.u (id integer PRIMARY KEY, k integer CONSTRAINT u_k_check CHECK (k > 0));
CREATE TABLE "Quoted" (id integer CONSTRAINT "Quoted_pkey" PRIMARY KEY);
CREATE TYPE mood AS ENUM ('a', 'b');
CREATE TABLE e (id integer PRIMARY KEY, m mood CONSTRAINT e_m_check CHECK (m <> 'a'));
CREATE TABLE ek (m mood PRIMARY KEY);
CREATE INDEX t_n_idx ON t (n);
CREATE VIEW tv AS SELECT code FROM t;
CREATE TABLE ip (a integer CONSTRAINT ip_a_check CHECK (a > 0));
CREATE TABLE ic (b integer) INHERITS (ip);
CREATE TABLE cvt (a integer);
CREATE VIEW cv AS SELECT count(*) AS c FROM cvt;
`

// TestScanRulesAgainstPostgres runs each change on PostgreSQL, statement
// by statement, and derives from the server's catalogs what
// DisallowDropConstraint and RequirePrimaryKey should report: the
// constraints of the schema before the change that an ALTER TABLE ...
// DROP CONSTRAINT removed, the tables a CREATE TABLE made, and the keyed
// tables an ALTER TABLE left without their key, still without a key at
// the end. The review, given the schema synced from the same server,
// must never report anything else. Where the scan follows the whole
// change (complete), it must report all of it.
//
// Up to the first statement the server rejects: the change stops there
// when it runs. That statement does nothing, so a finding on it is wrong,
// and so is a finding after it, or a RequirePrimaryKey finding about the
// end of a change that never gets there: the review must stop where the
// server does, or before.
func TestScanRulesAgainstPostgres(t *testing.T) {
	db := startPostgres(t)
	tests := []struct {
		name     string
		change   string
		complete bool
	}{
		{"every constraint kind", "ALTER TABLE t DROP CONSTRAINT t_pkey CASCADE, DROP CONSTRAINT t_code_key, DROP CONSTRAINT t_p_fk, DROP CONSTRAINT t_n_check;", true},
		{"names resolve along the search path", "ALTER TABLE s.t DROP CONSTRAINT t_pkey;\nSET search_path = s, public;\nALTER TABLE t DROP CONSTRAINT IF EXISTS t_pkey;\nALTER TABLE p DROP CONSTRAINT p_pkey CASCADE;", true},
		{"$user resolves to the session user's schema", "ALTER TABLE u DROP CONSTRAINT u_k_check;\nALTER TABLE u DROP CONSTRAINT u_pkey;", true},
		// A renamed table is not followed; a renamed constraint is.
		{"renamed table", "ALTER TABLE t RENAME TO t2;\nALTER TABLE t2 DROP CONSTRAINT t_n_check;\nALTER TABLE IF EXISTS t DROP CONSTRAINT t_code_key;", false},
		{"renamed constraint", "ALTER TABLE t RENAME CONSTRAINT t_n_check TO t_n_check2;\nALTER TABLE t DROP CONSTRAINT t_n_check2;\nALTER TABLE t DROP CONSTRAINT IF EXISTS t_n_check;", true},
		{"renamed key, then dropped", "ALTER TABLE e RENAME CONSTRAINT e_pkey TO e_k;\nALTER TABLE e DROP CONSTRAINT e_k;", true},
		{"renamed key's index takes the new name", "ALTER TABLE e RENAME CONSTRAINT e_pkey TO e_k;\nCREATE TABLE e_k (id int);", true},
		{"renamed key's index frees the old name", "ALTER TABLE e RENAME CONSTRAINT e_pkey TO e_k;\nCREATE TABLE e_pkey (id int);", true},
		{"renamed key still keys its table", "ALTER TABLE e RENAME CONSTRAINT e_pkey TO e_k;\nALTER TABLE e DROP CONSTRAINT e_m_check, ADD PRIMARY KEY (m);", true},
		{"renamed key's old name is gone", "ALTER TABLE e RENAME CONSTRAINT e_pkey TO e_k;\nALTER TABLE e DROP CONSTRAINT e_m_check, DROP CONSTRAINT e_pkey;", true},
		{"renamed table frees its name", "ALTER TABLE nokey RENAME TO nk;\nCREATE TABLE nokey (a int);", true},
		// The catalog renames the index alone, so it rejects the DROP.
		{"renamed index renames its constraint", "ALTER INDEX t_code_key RENAME TO t_code_uk;\nALTER TABLE t DROP CONSTRAINT t_code_uk;\nALTER TABLE t DROP CONSTRAINT IF EXISTS t_code_key;", false},
		{"DROP TYPE CASCADE drops the columns of the type", "DROP TYPE mood CASCADE;\nALTER TABLE e DROP CONSTRAINT IF EXISTS e_m_check;\nALTER TABLE ek ADD COLUMN x int;", true},
		{"constraint the change added", "ALTER TABLE nokey ADD CONSTRAINT nokey_a_check CHECK (a > 0);\nALTER TABLE nokey DROP CONSTRAINT nokey_a_check;", true},
		{"dropped and added again", "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT t_n_check CHECK (n > 1);\nALTER TABLE t DROP CONSTRAINT t_n_check;", true},
		{"missing constraint under IF EXISTS", "ALTER TABLE t DROP CONSTRAINT IF EXISTS nope;\nALTER TABLE IF EXISTS nope DROP CONSTRAINT x;", true},
		{"key referenced by a foreign key, without CASCADE", "ALTER TABLE t DROP CONSTRAINT t_pkey;", true},
		{"CASCADE drops the foreign keys on the key", "ALTER TABLE t DROP CONSTRAINT t_pkey CASCADE;\nALTER TABLE r DROP CONSTRAINT IF EXISTS r_t_fk;", true},
		{"DROP TABLE CASCADE drops the foreign keys on it", "DROP TABLE t CASCADE;\nALTER TABLE r DROP CONSTRAINT IF EXISTS r_t_fk;", true},
		{"DROP COLUMN CASCADE drops the foreign keys on it", "ALTER TABLE p DROP COLUMN id CASCADE;\nALTER TABLE t DROP CONSTRAINT IF EXISTS t_p_fk;", true},
		// A cascading drop of a schema ends the scan.
		{"DROP SCHEMA CASCADE", "DROP SCHEMA s CASCADE;\nSET search_path = s, public;\nALTER TABLE t DROP CONSTRAINT t_n_check;", false},
		{"renamed column keeps the check", "ALTER TABLE t RENAME COLUMN n TO m;\nALTER TABLE t DROP CONSTRAINT t_n_check;", true},
		// A moved table is not followed, and neither is a table the change
		// creates with a key.
		{"moved to another schema", "CREATE SCHEMA s2;\nALTER TABLE t SET SCHEMA s2;\nALTER TABLE s2.t DROP CONSTRAINT t_n_check;\nALTER TABLE s2.t DROP CONSTRAINT t_pkey CASCADE;", false},
		{"copy of a table", "CREATE TABLE c (LIKE t INCLUDING ALL);\nALTER TABLE c DROP CONSTRAINT c_pkey;\nALTER TABLE c DROP CONSTRAINT t_n_check;", false},
		{"shadowed by the search path", "SET search_path = s, public;\nALTER TABLE t DROP CONSTRAINT t_pkey;\nALTER TABLE public.t DROP CONSTRAINT t_code_key;", true},
		{"shadowed by a new table", "CREATE TABLE postgres.t (id int CONSTRAINT t_pkey UNIQUE, n int CONSTRAINT t_n_check CHECK (n > 0));\nALTER TABLE t DROP CONSTRAINT t_pkey;\nALTER TABLE t DROP CONSTRAINT t_n_check;\nALTER TABLE public.t DROP CONSTRAINT t_code_key;", true},
		{"shadowed by a new view", "CREATE VIEW postgres.nokey AS SELECT 1 AS a;\nCREATE TABLE IF NOT EXISTS nokey (a int);\nALTER TABLE public.t DROP CONSTRAINT t_n_check;", true},
		{"dropped with its column", "ALTER TABLE t DROP COLUMN n;\nALTER TABLE t DROP CONSTRAINT IF EXISTS t_n_check;", true},
		// After a column's drop the scan no longer follows the table's
		// constraints, nor a key column through a rename.
		{"dropped with its column in the same statement", "ALTER TABLE t DROP COLUMN n, DROP CONSTRAINT IF EXISTS t_n_check, DROP CONSTRAINT t_p_fk;", false},
		{"column a view reads, without CASCADE", "ALTER TABLE t DROP COLUMN code;", true},
		{"column a view reads, with CASCADE", "ALTER TABLE t DROP COLUMN code CASCADE;\nALTER TABLE t DROP CONSTRAINT IF EXISTS t_code_key;\nCREATE VIEW tv AS SELECT 1 AS x;", true},
		{"renamed column, then dropped", "ALTER TABLE t RENAME COLUMN id TO tid;\nALTER TABLE r DROP CONSTRAINT r_t_fk;\nALTER TABLE t DROP COLUMN tid;\nALTER TABLE t DROP COLUMN IF EXISTS id;", false},
		{"a missing column is refused", "ALTER TABLE t DROP CONSTRAINT t_n_check, DROP COLUMN nope;", true},
		{"a missing constraint is refused", "ALTER TABLE t DROP CONSTRAINT t_n_check, DROP CONSTRAINT nope;", true},
		{"dropped key added again in the same statement", "ALTER TABLE t ADD PRIMARY KEY (code), DROP CONSTRAINT t_pkey CASCADE;", true},
		{"IF NOT EXISTS in a new schema on the path", "CREATE SCHEMA x;\nSET search_path = x, public;\nCREATE TABLE IF NOT EXISTS t (id int);\nCREATE TABLE IF NOT EXISTS public.nokey (a int);\nCREATE TABLE IF NOT EXISTS public.fresh (a int);", false},
		{"same names in two schemas", "ALTER TABLE s.t DROP CONSTRAINT t_pkey;\nALTER TABLE public.t DROP CONSTRAINT t_pkey CASCADE;", true},
		{"a dropped table no longer holds its foreign keys", "DROP TABLE r;\nALTER TABLE t DROP CONSTRAINT t_pkey;", true},
		{"a dropped table frees its index names", "DROP TABLE e;\nCREATE TABLE IF NOT EXISTS e_pkey (id int);", true},
		{"a later subcommand fails", "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD COLUMN id int;", true},
		{"a later subcommand alters a missing column", "ALTER TABLE t DROP CONSTRAINT t_n_check, ALTER COLUMN nope SET DEFAULT 1;", true},
		{"a later subcommand reuses a name", "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT t_code_key CHECK (n > 1);", true},
		{"a later subcommand adds a second key", "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD PRIMARY KEY (code);", true},
		{"drops run first", "ALTER TABLE t ADD CONSTRAINT t_n_check CHECK (n > 1), DROP CONSTRAINT t_n_check, DROP COLUMN n, ADD COLUMN n int, ADD COLUMN IF NOT EXISTS id int;\nALTER TABLE s.t DROP CONSTRAINT t_pkey, ADD COLUMN x int, ALTER COLUMN id SET NOT NULL;", true},
		{"dropping a table others depend on", "DROP TABLE p;\nALTER TABLE t DROP CONSTRAINT t_n_check;", true},
		{"dropping a table a view reads", "DROP TABLE r;\nDROP TABLE t;\nALTER TABLE p DROP CONSTRAINT p_pkey CASCADE;", true},
		{"dropping both ends of a foreign key", "DROP VIEW tv;\nDROP TABLE r, t;\nALTER TABLE p DROP CONSTRAINT p_pkey;", true},
		{"an unrelated cascading column drop", "ALTER TABLE p ADD COLUMN note int;\nALTER TABLE p DROP COLUMN note CASCADE;\nALTER TABLE t DROP CONSTRAINT t_p_fk;", true},
		{"the same constraint dropped twice", "ALTER TABLE t DROP CONSTRAINT t_n_check, DROP CONSTRAINT t_n_check;", true},
		{"a constraint gone with its column", "ALTER TABLE t DROP CONSTRAINT t_code_key, DROP COLUMN n, DROP CONSTRAINT t_n_check;\nALTER TABLE s.t DROP CONSTRAINT t_pkey;", true},
		{"a new foreign key on one key", "CREATE TABLE c (t_id int REFERENCES t (id));\nALTER TABLE t DROP CONSTRAINT t_code_key;", true},
		{"a new foreign key on the primary key", "CREATE TABLE c2 (t_id int REFERENCES t);\nALTER TABLE r DROP CONSTRAINT r_t_fk;\nALTER TABLE t DROP CONSTRAINT t_pkey;", true},
		{"a view a new view reads", "CREATE VIEW tv2 AS SELECT code FROM tv;\nDROP VIEW tv;\nCREATE TABLE n (id int);", true},
		{"a qualified move of a new table", "CREATE TABLE public.n (id int);\nCREATE SCHEMA x;\nCREATE TABLE x.n (id int);\nALTER TABLE public.n SET SCHEMA s;\nALTER TABLE x.n ADD PRIMARY KEY (id);", true},
		{"a new foreign key dropped with its table", "ALTER TABLE r DROP CONSTRAINT r_t_fk;\nCREATE TABLE c (t_id int REFERENCES t (id));\nDROP TABLE c;\nALTER TABLE t DROP CONSTRAINT t_pkey;", true},
		{"a new foreign key dropped by name", "ALTER TABLE r DROP CONSTRAINT r_t_fk;\nCREATE TABLE c (t_id int CONSTRAINT c_fk REFERENCES t (id));\nALTER TABLE c DROP CONSTRAINT c_fk;\nALTER TABLE t DROP CONSTRAINT t_pkey;", true},
		{"a new view on a table of the same name elsewhere", "CREATE VIEW v AS SELECT * FROM s.t;\nDROP TABLE r;\nDROP VIEW tv;\nDROP TABLE public.t;\nCREATE TABLE n (id int);", true},
		{"dropping a missing table", "DROP TABLE nope;\nCREATE TABLE n (id int);", true},
		{"a column added twice", "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD COLUMN x int, ADD COLUMN x int;", true},
		{"tables of CREATE SCHEMA", "CREATE SCHEMA z CREATE TABLE zt (id int) CREATE TABLE zk (id int PRIMARY KEY);", true},
		{"a CTE named like a table", "CREATE VIEW v AS WITH nokey AS (SELECT 1 AS a) SELECT * FROM nokey;\nDROP TABLE nokey;\nCREATE TABLE n (id int);", true},
		{"a view the change drops", "CREATE VIEW v AS SELECT * FROM nokey;\nDROP VIEW v;\nDROP TABLE nokey;\nCREATE TABLE n (id int);", true},
		{"a function the change creates returns the table", "CREATE FUNCTION g() RETURNS SETOF nokey LANGUAGE sql AS 'SELECT * FROM nokey';\nDROP TABLE nokey;\nCREATE TABLE n (id int);", true},
		{"DROP TABLE of a view", "DROP TABLE tv;\nCREATE TABLE n (id int);", true},
		{"DROP INDEX of a key's index", "DROP INDEX t_pkey;\nCREATE TABLE n (id int);", true},
		{"CREATE TABLE of an existing name", "CREATE TABLE public.nokey (a int);\nCREATE TABLE n (id int);", true},
		{"CREATE TABLE in the session user's schema", "CREATE TABLE nokey (a int);", true},
		{"a CTE body reads the table it shadows", "CREATE VIEW v AS WITH nokey AS (SELECT * FROM nokey) SELECT * FROM nokey;\nDROP TABLE nokey;\nCREATE TABLE n (id int);", true},
		{"IF NOT EXISTS in the session user's schema", "CREATE TABLE IF NOT EXISTS nokey (a text REFERENCES t (code));\nALTER TABLE t DROP CONSTRAINT t_code_key;", true},
		{"ADD COLUMN IF NOT EXISTS of an existing column", "ALTER TABLE r DROP CONSTRAINT r_t_fk;\nALTER TABLE t DROP CONSTRAINT t_pkey;\nALTER TABLE t ADD COLUMN IF NOT EXISTS id int PRIMARY KEY;", true},
		// A name of the form the server generates is not reported after a
		// statement that may have generated it, wherever it would go.
		{"a table named like a generated index elsewhere", "CREATE INDEX ON nokey (a);\nCREATE TABLE nokey_a_idx (x int);", false},
		{"a table named like a generated index", "CREATE INDEX ON nokey (a);\nCREATE TABLE public.nokey_a_idx (x int);", true},
		{"a table named like a key's index", "CREATE TABLE a (id int CONSTRAINT foo PRIMARY KEY);\nCREATE TABLE foo (x int);", true},
		{"a table in a missing schema", "CREATE TABLE missing.n (id int);", true},
		{"IF EXISTS of a missing table, then created", "DROP TABLE IF EXISTS gone;\nCREATE TABLE gone (id int);", true},
		{"DROP FUNCTION IF EXISTS of another signature", "CREATE FUNCTION g(integer) RETURNS SETOF nokey LANGUAGE sql AS 'SELECT * FROM nokey';\nDROP FUNCTION IF EXISTS g(text);\nDROP TABLE nokey;\nCREATE TABLE n (id int);", true},
		{"a new view replaced", "CREATE VIEW v AS SELECT * FROM nokey;\nCREATE OR REPLACE VIEW v AS SELECT 1 AS a;\nDROP TABLE nokey;\nCREATE TABLE n (id int);", true},
		{"a table in another database", "CREATE TABLE other.public.n (id int);", true},
		{"ALTER TABLE of a missing table", "ALTER TABLE nope DROP CONSTRAINT x;\nCREATE TABLE n (id int);", true},
		{"ALTER TABLE IF EXISTS of a missing table", "ALTER TABLE IF EXISTS nope DROP CONSTRAINT x;\nCREATE TABLE n (id int);", true},
		{"CREATE SCHEMA of an existing schema", "CREATE SCHEMA s CREATE TABLE zt (id int);", true},
		{"DROP SCHEMA IF EXISTS of a missing schema", "DROP SCHEMA IF EXISTS missing;\nCREATE TABLE missing.n (id int);", true},
		{"a skipped ADD COLUMN", "ALTER TABLE t ADD COLUMN IF NOT EXISTS id int PRIMARY KEY;\nCREATE TABLE n (id int);", true},
		{"two column drops", "ALTER TABLE r DROP CONSTRAINT r_t_fk;\nALTER TABLE t DROP COLUMN n, DROP COLUMN id;", true},
		{"EXPLAIN ANALYZE of CREATE MATERIALIZED VIEW", "EXPLAIN ANALYZE CREATE MATERIALIZED VIEW mv AS SELECT * FROM nokey;\nDROP TABLE nokey;\nCREATE TABLE n (id int);", true},
		{"CREATE SCHEMA IF NOT EXISTS with elements", "CREATE SCHEMA IF NOT EXISTS s CREATE TABLE zt (id int);", true},
		{"CREATE VIEW of an existing name", "CREATE VIEW public.nokey AS SELECT 1 AS a;\nCREATE TABLE n (id int);", true},
		{"CREATE VIEW in the session user's schema", "CREATE VIEW nokey AS SELECT 1 AS a;\nCREATE TABLE n (id int);", true},
		{"a new function takes the table's row type", "CREATE FUNCTION g(x nokey) RETURNS integer LANGUAGE sql AS 'SELECT 1';\nDROP TABLE nokey;\nCREATE TABLE n (id int);", true},
		{"a definite rename of a new table", "CREATE TABLE public.a (id int);\nALTER TABLE public.a RENAME TO b;\nCREATE TABLE public.a (id int PRIMARY KEY);\nDROP TABLE public.a;", true},
		{"a column drop frees a constraint name", "ALTER TABLE t DROP COLUMN code CASCADE, ADD CONSTRAINT t_code_key UNIQUE (n);\nCREATE TABLE n2 (id int);", true},
		{"CREATE SCHEMA IF NOT EXISTS of an existing schema", "CREATE SCHEMA IF NOT EXISTS s;\nCREATE TABLE s.t (id int);", true},
		{"CREATE SEQUENCE of an existing name", "CREATE SEQUENCE public.nokey;\nCREATE TABLE n (id int);", true},
		{"DROP SCHEMA of a nonempty schema", "DROP SCHEMA s;\nCREATE TABLE s.n (id int);", true},
		{"a renamed new view", "CREATE VIEW public.v AS SELECT * FROM public.nokey;\nALTER VIEW public.v RENAME TO w;\nDROP VIEW public.w;\nDROP TABLE public.nokey;\nCREATE TABLE n (id int);", true},
		{"a new view dropped with what it reads", "CREATE VIEW w AS SELECT * FROM tv;\nDROP VIEW tv, w;\nCREATE TABLE n (id int);", true},
		{"a column renamed twice", "ALTER TABLE r DROP CONSTRAINT r_t_fk;\nALTER TABLE t RENAME COLUMN id TO x;\nALTER TABLE t RENAME COLUMN x TO y;\nALTER TABLE t DROP CONSTRAINT t_pkey;\nALTER TABLE t ADD COLUMN IF NOT EXISTS x int PRIMARY KEY;", true},
		{"a column added twice in one statement", "CREATE TABLE n (id int);\nALTER TABLE n ADD COLUMN x int, ADD COLUMN IF NOT EXISTS x int PRIMARY KEY;", true},
		{"DROP COLUMN IF EXISTS of a missing column", "ALTER TABLE t DROP COLUMN IF EXISTS nope;\nALTER TABLE t DROP CONSTRAINT t_n_check;", true},
		{"a new SQL-standard body reads the table", "CREATE FUNCTION g() RETURNS bigint LANGUAGE SQL RETURN (SELECT count(*) FROM nokey);\nDROP TABLE nokey;\nCREATE TABLE n (id int);", true},
		{"DROP of a name in another database", "DROP TABLE other.public.nokey;\nCREATE TABLE n (id int);", true},
		{"a foreign key a cascade took", "ALTER TABLE p DROP CONSTRAINT p_pkey CASCADE;\nALTER TABLE t DROP CONSTRAINT t_p_fk;\nCREATE TABLE n (id int);", true},
		{"a cascade past a replaced view", "CREATE VIEW tv2 AS SELECT code FROM tv;\nCREATE OR REPLACE VIEW tv AS SELECT 'x'::text AS code;\nDROP TABLE r;\nDROP TABLE t CASCADE;\nCREATE TABLE tv2 (id int);", true},
		{"DROP SCHEMA of a schema with a table, at the end", "CREATE TABLE n (id int);\nDROP SCHEMA s;", true},
		{"an added key on a missing column", "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT replacement UNIQUE (nope);", true},
		{"DROP FUNCTION of another signature", "CREATE FUNCTION g(integer) RETURNS SETOF nokey LANGUAGE sql AS 'SELECT * FROM nokey';\nDROP FUNCTION g(text);\nDROP TABLE nokey;", true},
		{"a no-op column drop and a reused name", "ALTER TABLE t DROP CONSTRAINT t_code_key, DROP COLUMN IF EXISTS nope, ADD CONSTRAINT t_n_check CHECK (n > 0);", true},
		{"a dropped column takes its foreign key", "ALTER TABLE t DROP COLUMN p_id;\nALTER TABLE p DROP CONSTRAINT p_pkey;", true},
		{"CREATE SCHEMA elements of one name", "CREATE SCHEMA z CREATE TABLE x (id int) CREATE VIEW x AS SELECT 1;", true},
		{"CREATE SCHEMA element named like a generated name", "CREATE SCHEMA z CREATE TABLE x (id int UNIQUE) CREATE TABLE x_id_key (id int);", true},
		{"CREATE SCHEMA element in another schema", "CREATE SCHEMA z CREATE TABLE x (id int) CREATE TABLE s.y (id int);", true},
		{"added foreign key to a missing column", "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT bad FOREIGN KEY (p_id) REFERENCES p (nope);", true},
		{"added foreign key to a table without a key", "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT bad FOREIGN KEY (n) REFERENCES nokey;", true},
		{"added foreign key of another column count", "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT bad FOREIGN KEY (p_id, n) REFERENCES p (id);", true},
		{"added foreign key to a view", "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT bad FOREIGN KEY (code) REFERENCES tv (code);", true},
		{"added foreign key to a unique index the change made", "CREATE UNIQUE INDEX ON nokey (a);\nALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT good FOREIGN KEY (n) REFERENCES nokey (a);", false},
		{"added foreign key to its own table without a key", "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT bad FOREIGN KEY (n) REFERENCES t (n);", true},
		{"added foreign key to its own table's new key", "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT good FOREIGN KEY (n) REFERENCES t (n), ADD CONSTRAINT t_n_key UNIQUE (n);", true},
		{"a column added after the constraint on it", "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT replacement UNIQUE (x), ADD COLUMN x int;", true},
		{"a cascade takes a view the change made", "CREATE VIEW w AS SELECT * FROM t;\nDROP TABLE r;\nDROP TABLE t CASCADE;\nCREATE TABLE w (id int);", true},
		{"a cascade takes views reading a view the change made", "CREATE VIEW w AS SELECT * FROM tv;\nCREATE VIEW w2 AS SELECT * FROM public.w;\nDROP TABLE r;\nDROP TABLE t CASCADE;\nCREATE TABLE w (id int);\nCREATE TABLE w2 (id int);", true},
		{"DROP FUNCTION of a name the target lacks", "DROP FUNCTION public.missing();\nCREATE TABLE n (id int);", true},
		{"DROP PROCEDURE", "CREATE PROCEDURE pr(a nokey) LANGUAGE sql AS 'SELECT 1';\nDROP PROCEDURE pr(nokey);\nDROP TABLE nokey;\nCREATE TABLE n (id int);", true},
		{"a renamed function", "CREATE FUNCTION g(nokey) RETURNS int LANGUAGE sql AS 'SELECT 1';\nALTER FUNCTION g(nokey) RENAME TO h;\nDROP FUNCTION h(nokey);\nDROP TABLE nokey;\nCREATE TABLE n (id int);", true},
		{"a renamed function holds the table", "CREATE FUNCTION g(nokey) RETURNS int LANGUAGE sql AS 'SELECT 1';\nALTER FUNCTION g(nokey) RENAME TO h;\nDROP TABLE nokey;\nCREATE TABLE n (id int);", true},
		{"a cascade takes a function the change made", "CREATE FUNCTION g(nokey) RETURNS SETOF r LANGUAGE sql AS 'SELECT * FROM r';\nDROP TABLE nokey CASCADE;\nDROP TABLE r;\nCREATE TABLE n (id int);", true},
		{"a plain table generates no names", "CREATE TABLE public.a (id int);\nCREATE TABLE public.n_idx (id int);\nDROP TABLE public.a;", true},
		{"LIKE of a keyless table the change made", "CREATE TABLE public.a (id int);\nCREATE TABLE public.b (LIKE public.a INCLUDING INDEXES);\nDROP TABLE public.a;", true},
		{"CREATE TABLE IF NOT EXISTS in a new schema", "CREATE SCHEMA z;\nCREATE TABLE IF NOT EXISTS z.n (id int);", true},
		{"CREATE TABLE of an enum type's name", "CREATE TABLE IF NOT EXISTS mood (id int);", true},
		{"CREATE TABLE of a domain the change made", "CREATE DOMAIN d AS int;\nCREATE TABLE d (id int);", true},
		{"DROP FUNCTION IF EXISTS of a function the change made", "CREATE FUNCTION public.f(public.nokey) RETURNS integer LANGUAGE sql AS 'SELECT 1';\nDROP FUNCTION IF EXISTS public.f(public.nokey);\nDROP TABLE public.nokey;\nCREATE TABLE n (id int);", true},
		{"a view follows a renamed table", "CREATE VIEW public.v AS SELECT * FROM public.nokey;\nALTER TABLE public.nokey RENAME TO u;\nDROP TABLE public.u;\nCREATE TABLE n (id int);", true},
		{"a cascading DROP TABLE IF EXISTS of a missing table", "ALTER TABLE t DROP CONSTRAINT t_pkey CASCADE;\nDROP TABLE IF EXISTS public.missing CASCADE;", true},
		{"USING INDEX of a missing index", "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT replacement UNIQUE USING INDEX missing_idx;", true},
		{"USING INDEX of an index that is not unique", "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT t_n_key UNIQUE USING INDEX t_n_idx;", true},
		{"a key named after a table", "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT nokey UNIQUE (n);", true},
		{"a key named after another table's index", "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT p_pkey UNIQUE (n);", true},
		{"a key named after the one the statement drops", "ALTER TABLE t DROP CONSTRAINT t_code_key, ADD CONSTRAINT t_code_key UNIQUE (code);", true},
		{"a foreign key to a column twice", "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT bad FOREIGN KEY (p_id, n) REFERENCES p (id, id);", true},
		{"CREATE TABLE with a column twice", "CREATE TABLE n (id int, id text);", true},
		{"a rename to a name the schema holds", "ALTER TABLE t RENAME TO nokey;\nCREATE TABLE n (id int);", true},
		{"a view follows a table to another schema", "CREATE VIEW public.v AS SELECT * FROM public.nokey;\nALTER TABLE public.nokey SET SCHEMA s;\nDROP TABLE s.nokey;\nCREATE TABLE n (id int);", true},
		{"dropping a key frees its index's name", "ALTER TABLE t DROP CONSTRAINT t_code_key;\nCREATE TABLE IF NOT EXISTS t_code_key (id int);", true},
		{"a qualified DROP of a table the change made", "CREATE TABLE public.n (id int PRIMARY KEY);\nDROP TABLE public.n;\nCREATE TABLE public.n (id int);", true},
		{"a rename of a relation the target lacks", "ALTER TABLE public.missing RENAME TO x;\nCREATE TABLE n (id int);", true},
		{"a rename under IF EXISTS of a relation the target lacks", "ALTER TABLE IF EXISTS public.missing RENAME TO x;\nCREATE TABLE n (id int);", true},
		{"DROP SCHEMA of a schema the change put a table in", "CREATE SCHEMA z;\nCREATE TABLE z.x (id int PRIMARY KEY);\nDROP SCHEMA z;\nCREATE TABLE n (id int);", true},
		{"DROP SCHEMA of a schema the change emptied again", "CREATE SCHEMA z;\nCREATE TABLE z.x (id int PRIMARY KEY);\nDROP TABLE z.x;\nDROP SCHEMA z;\nCREATE TABLE n (id int);", true},
		{"a foreign key of a table the same DROP drops", "CREATE TABLE public.c (id int PRIMARY KEY, n_a int REFERENCES public.p);\nALTER TABLE t DROP CONSTRAINT t_p_fk;\nDROP TABLE public.c, public.p;\nCREATE TABLE n (id int);", true},
		{"INHERITS of a parent the target lacks", "CREATE TABLE child (id int) INHERITS (public.missing);", true},
		{"INHERITS of a view", "CREATE TABLE child (id int) INHERITS (tv);", true},
		{"a second CREATE TABLE of a name the change made", "CREATE TABLE public.x (id int PRIMARY KEY);\nCREATE TABLE public.x (id int);\nCREATE TABLE n (id int);", true},
		{"ALTER TABLE of an index", "ALTER TABLE public.t_pkey DROP CONSTRAINT x;\nCREATE TABLE n (id int);", true},
		{"SET SCHEMA vacates the schema", "ALTER TABLE s.t SET SCHEMA postgres;\nDROP SCHEMA s;\nCREATE TABLE n (id int);", true},
		{"DROP SCHEMA after dropping the function made there", "CREATE SCHEMA z;\nCREATE FUNCTION z.f() RETURNS int LANGUAGE sql AS 'SELECT 1';\nDROP FUNCTION z.f();\nDROP SCHEMA z;\nCREATE TABLE n (id int);", true},
		{"an unqualified name resolved elsewhere", "CREATE TABLE postgres.t (id int);\nSET search_path = s, postgres;\nALTER TABLE t ADD CONSTRAINT t_pk PRIMARY KEY (id);", true},
		{"CREATE TABLE with a foreign key to a table the target lacks", "CREATE TABLE public.n (id int REFERENCES public.missing);", true},
		{"CREATE TABLE with a foreign key to a synced key", "CREATE TABLE n (id int, p_id int REFERENCES p);", true},
		{"CREATE TABLE with a foreign key to columns without a key", "CREATE TABLE n (id int, t_n int REFERENCES t (n));", true},
		{"CREATE TABLE with a foreign key to a table the change made", "CREATE TABLE a (id int PRIMARY KEY);\nCREATE TABLE b (id int, a_id int REFERENCES a);", true},
		{"CREATE TABLE with a foreign key to itself without the key", "CREATE TABLE n (id int, parent int REFERENCES n (id));", true},
		{"a view of CREATE SCHEMA dropped", "CREATE SCHEMA z CREATE VIEW v AS SELECT 1;\nDROP VIEW z.v;\nCREATE TABLE z.v (id int);", true},
		{"DROP TYPE of a type a column uses", "DROP TYPE mood;\nCREATE TABLE n (id int);", true},
		{"a generated name in another schema", "CREATE TABLE s.x (id serial PRIMARY KEY);\nCREATE TABLE postgres.audit_seq (id int);", true},
		{"DROP COLUMN read by a generated column", "ALTER TABLE nokey ADD COLUMN b int GENERATED ALWAYS AS (a + 1) STORED;\nALTER TABLE nokey DROP COLUMN a;\nCREATE TABLE n (id int);", false},
		{"a sequence the change made and dropped", "CREATE SEQUENCE public.q;\nDROP SEQUENCE public.q;\nCREATE TABLE public.q (id int);", true},
		{"LIKE of a table whose key the change dropped", "ALTER TABLE t DROP CONSTRAINT t_pkey CASCADE;\nCREATE TABLE n (LIKE t INCLUDING INDEXES);", true},
		{"LIKE of a table the change gave a key", "ALTER TABLE nokey ADD PRIMARY KEY (a);\nCREATE TABLE n (LIKE nokey INCLUDING INDEXES);", true},
		{"CREATE VIEW of a relation the target lacks", "CREATE VIEW public.v AS SELECT * FROM public.missing;\nCREATE TABLE n (id int);", true},
		{"CREATE VIEW of a CTE named like a missing relation", "CREATE VIEW public.v AS WITH missing AS (SELECT 1 AS a) SELECT * FROM missing;\nCREATE TABLE n (id int);", true},
		{"CREATE TABLE OF a type the target lacks", "CREATE TABLE n OF public.missing;", true},
		{"CREATE TABLE OF a table's row type", "CREATE TABLE n OF public.nokey;", true},
		{"CREATE TABLE OF a composite type the change made", "CREATE TYPE ct AS (a int);\nCREATE TABLE n OF ct;", true},
		{"CREATE TABLE naming a constraint twice", "CREATE TABLE n (a int CONSTRAINT c CHECK (a > 0), b int CONSTRAINT c CHECK (b > 0));", true},
		{"CREATE TABLE with a key named after a relation", "CREATE TABLE n (a int CONSTRAINT nokey UNIQUE);", true},
		{"CREATE TABLE with a key named after the table", "CREATE TABLE n (a int CONSTRAINT n UNIQUE);", true},
		{"CREATE TABLE with a foreign key naming a column twice", "CREATE TABLE n (a int, FOREIGN KEY (a, a) REFERENCES p (id, id));", true},
		{"RENAME COLUMN of a column the table lacks", "ALTER TABLE t RENAME COLUMN missing TO x;\nCREATE TABLE n (id int);", true},
		{"RENAME COLUMN to a column the table has", "ALTER TABLE t RENAME COLUMN n TO code;\nCREATE TABLE n (id int);", true},
		{"RENAME COLUMN of a column the table has", "ALTER TABLE t RENAME COLUMN n TO n2;\nCREATE TABLE n (id int);", true},
		{"CREATE FUNCTION in a schema the target lacks", "CREATE FUNCTION missing.f() RETURNS int LANGUAGE sql AS 'SELECT 1';\nCREATE TABLE n (id int);", true},
		{"CREATE TABLE AS of a name the schema holds", "CREATE TABLE public.nokey AS SELECT 1;\nCREATE TABLE n (id int);", true},
		{"CREATE TABLE AS IF NOT EXISTS of a name the schema holds", "CREATE TABLE IF NOT EXISTS public.nokey AS SELECT 1;\nCREATE TABLE n (id int);", true},
		{"SELECT INTO a name the schema holds", "SELECT 1 AS a INTO public.nokey;\nCREATE TABLE n (id int);", true},
		{"CREATE TABLE AS of a relation the target lacks", "CREATE TABLE x AS SELECT * FROM public.missing;\nCREATE TABLE n (id int);", true},
		{"RENAME CONSTRAINT to a name the table has", "ALTER TABLE t RENAME CONSTRAINT t_n_check TO t_pkey;\nCREATE TABLE n (id int);", true},
		{"RENAME CONSTRAINT of a key to a relation's name", "ALTER TABLE t RENAME CONSTRAINT t_code_key TO nokey;\nCREATE TABLE n (id int);", true},
		{"ALTER TABLE of a table the change dropped", "DROP TABLE public.nokey;\nALTER TABLE public.nokey ADD COLUMN x int;\nCREATE TABLE n (id int);", true},
		{"DROP TABLE IF EXISTS of a table the change dropped", "DROP TABLE public.nokey;\nDROP TABLE IF EXISTS public.nokey;\nCREATE TABLE n (id int);", true},
		{"CREATE FUNCTION with a body of a relation the target lacks", "CREATE FUNCTION public.f() RETURNS bigint LANGUAGE SQL RETURN (SELECT count(*) FROM public.missing);\nCREATE TABLE n (id int);", true},
		{"an added check on a column the table lacks", "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT bad CHECK (nope > 0);", true},
		{"DROP FUNCTION of a scalar signature with an array overload", "CREATE FUNCTION g(int[]) RETURNS SETOF nokey LANGUAGE sql AS 'SELECT * FROM nokey';\nDROP FUNCTION g(int);\nDROP TABLE nokey;\nCREATE TABLE n (id int);", true},
		{"DROP FUNCTION of an array signature", "CREATE FUNCTION g(int[]) RETURNS SETOF nokey LANGUAGE sql AS 'SELECT * FROM nokey';\nDROP FUNCTION g(int[]);\nDROP TABLE nokey;\nCREATE TABLE n (id int);", true},
		{"DROP FUNCTION of a built-in type written otherwise", "CREATE FUNCTION g(int4) RETURNS SETOF nokey LANGUAGE sql AS 'SELECT * FROM nokey';\nDROP FUNCTION g(integer);\nDROP TABLE nokey;\nCREATE TABLE n (id int);", true},
		{"CREATE FOREIGN TABLE of a name the schema holds", "CREATE FOREIGN TABLE public.nokey (id int) SERVER s;\nCREATE TABLE n (id int);", true},
		{"a cascade of a parent of the same name in another schema", "CREATE TABLE postgres.n (x int) INHERITS (s.t);\nDROP TABLE r;\nDROP TABLE public.t CASCADE;", true},
		{"a cascade of a renamed parent", "CREATE TABLE postgres.n (x int) INHERITS (nokey);\nALTER TABLE nokey RENAME TO nk2;\nDROP TABLE nk2 CASCADE;", true},
		{"dropping a renamed key column", "ALTER TABLE \"Quoted\" RENAME COLUMN id TO x;\nALTER TABLE \"Quoted\" DROP COLUMN x;", true},
		{"a renamed column's drop takes the constraint on it", "ALTER TABLE t RENAME COLUMN code TO c2;\nALTER TABLE t DROP COLUMN c2, DROP CONSTRAINT t_code_key;", true},
		{"a mixed ALTER with a subcommand certain to fail", "ALTER TABLE t DROP CONSTRAINT nope, ADD COLUMN id int;\nCREATE TABLE n (id int);", true},
		{"a drop of a statement that may not run", "ALTER TABLE t DROP CONSTRAINT t_p_fk, DROP CONSTRAINT nope;\nALTER TABLE p DROP CONSTRAINT p_pkey;", true},
		{"CREATE INDEX on a relation the target lacks", "CREATE INDEX i ON public.missing (id);\nCREATE TABLE n (id int);", true},
		{"CREATE INDEX on a view", "CREATE INDEX i ON tv (code);\nCREATE TABLE n (id int);", true},
		{"CREATE INDEX of a name the schema holds", "CREATE INDEX nokey ON t (n);\nCREATE TABLE n (id int);", true},
		{"CREATE INDEX IF NOT EXISTS of a name the schema holds", "CREATE INDEX IF NOT EXISTS nokey ON t (n);\nCREATE TABLE n (id int);", true},
		{"CREATE TYPE of a table's name", "CREATE TYPE public.nokey AS ENUM ('x');\nCREATE TABLE n (id int);", true},
		{"a foreign key to the table itself, qualified", "CREATE TABLE n (id int UNIQUE, parent int REFERENCES public.n (id));", true},
		{"a rename of a table the change made to a taken name", "CREATE TABLE public.x (id int PRIMARY KEY);\nALTER TABLE public.x RENAME TO nokey;\nCREATE TABLE n (id int);", true},
		{"DROP TABLE of a view the change made", "CREATE VIEW public.v AS SELECT 1;\nDROP TABLE public.v;\nCREATE TABLE n (id int);", true},
		{"a second CREATE SCHEMA", "CREATE SCHEMA z;\nCREATE SCHEMA z;\nCREATE TABLE n (id int);", true},
		{"CREATE TABLE with a check on a column it lacks", "CREATE TABLE x (a int CHECK (missing > 0));\nCREATE TABLE n (id int);", true},
		{"CREATE TABLE with a key on a column it lacks", "CREATE TABLE x (a int, PRIMARY KEY (missing));\nCREATE TABLE n (id int);", true},
		{"CREATE TABLE with a key naming a column twice", "CREATE TABLE x (a int, b int, UNIQUE (a, a));\nCREATE TABLE n (id int);", true},
		{"CREATE TABLE with two primary keys", "CREATE TABLE x (a int PRIMARY KEY, b int, PRIMARY KEY (b));\nCREATE TABLE n (id int);", true},
		{"CREATE FUNCTION of a signature the change made", "CREATE FUNCTION g() RETURNS int LANGUAGE sql AS 'SELECT 1';\nCREATE FUNCTION g() RETURNS int LANGUAGE sql AS 'SELECT 2';\nCREATE TABLE n (id int);", true},
		{"an added key naming a column twice", "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD UNIQUE (id, id);", true},
		{"a qualified DROP of another schema's table", "CREATE TABLE t2 (id int);\nDROP TABLE s.t;\nCREATE TABLE postgres.t2 (id int);", true},
		{"a table the search path puts in pg_temp", "SET search_path = pg_temp, public;\nCREATE TABLE n (id int);", true},
		{"CREATE VIEW reading an index", "CREATE VIEW v AS SELECT * FROM public.t_pkey;\nCREATE TABLE n (id int);", true},
		{"CREATE SCHEMA with an index on a table it lacks", "CREATE SCHEMA z CREATE INDEX i ON missing (id);\nCREATE TABLE n (id int);", true},
		{"CREATE SCHEMA with an index on its table", "CREATE SCHEMA z CREATE TABLE x (id int PRIMARY KEY) CREATE INDEX i ON x (id);\nCREATE TABLE n (id int);", true},
		{"CREATE FOREIGN TABLE with a column twice", "CREATE FOREIGN TABLE x (a int, a int) SERVER s;\nCREATE TABLE n (id int);", true},
		{"ALTER FUNCTION RENAME of a function the target lacks", "ALTER FUNCTION public.missing() RENAME TO f;\nCREATE TABLE n (id int);", true},
		{"ALTER PROCEDURE RENAME of a procedure the target lacks", "ALTER PROCEDURE public.missing() RENAME TO f;\nCREATE TABLE n (id int);", true},
		{"DROP NOT NULL of a primary key column", "ALTER TABLE t ALTER COLUMN id DROP NOT NULL;\nCREATE TABLE n (id int);", true},
		{"DROP NOT NULL of another column", "ALTER TABLE t ALTER COLUMN n DROP NOT NULL;\nCREATE TABLE n (id int);", true},
		{"dropping the key the change added", "CREATE TABLE n (id int);\nALTER TABLE n ADD PRIMARY KEY (id);\nALTER TABLE n DROP CONSTRAINT n_pkey;", true},
		{"dropping a named key the change added", "CREATE TABLE n (id int);\nALTER TABLE n ADD CONSTRAINT n_pk PRIMARY KEY (id);\nALTER TABLE n DROP CONSTRAINT n_pk;", true},
		{"a second CREATE TYPE of a name", "CREATE TYPE public.x AS ENUM ('a');\nCREATE TYPE public.x AS ENUM ('b');\nCREATE TABLE n (id int);", true},
		{"CREATE OR REPLACE VIEW of a table the change made", "CREATE TABLE public.x (id int PRIMARY KEY);\nCREATE OR REPLACE VIEW public.x AS SELECT 1;\nCREATE TABLE n (id int);", true},
		{"ALTER SEQUENCE of a sequence the target lacks", "ALTER SEQUENCE public.missing RESTART;\nCREATE TABLE n (id int);", true},
		{"ALTER SEQUENCE of a table", "ALTER SEQUENCE public.nokey RESTART;\nCREATE TABLE n (id int);", true},
		{"DROP TABLE of a table the change made that a new view reads", "CREATE TABLE public.x (id int PRIMARY KEY);\nCREATE VIEW public.w AS SELECT * FROM public.x;\nDROP TABLE public.x;\nCREATE TABLE n (id int);", true},
		{"a statement that may fail", "ALTER TABLE t DROP CONSTRAINT nope, ADD CONSTRAINT x CHECK (n > 0);\nALTER TABLE t DROP CONSTRAINT x;\nCREATE TABLE fresh (id int);", true},
		{"an added foreign key to a table the change made", "CREATE TABLE a (id int PRIMARY KEY);\nALTER TABLE t ADD CONSTRAINT t_a_fk FOREIGN KEY (n) REFERENCES a;\nCREATE TABLE b (id int);", true},
		{"CREATE FUNCTION of a signature the change made with arguments", "CREATE FUNCTION g(integer) RETURNS int LANGUAGE sql AS 'SELECT 1';\nCREATE FUNCTION g(int4) RETURNS int LANGUAGE sql AS 'SELECT 2';\nCREATE TABLE n (id int);", true},
		{"ALTER INDEX RENAME of a table", "ALTER INDEX public.t RENAME TO x;\nCREATE TABLE n (id int);", true},
		{"ALTER SEQUENCE RENAME of a table", "ALTER SEQUENCE public.t RENAME TO x;\nCREATE TABLE n (id int);", true},
		{"ALTER TABLE RENAME of a view", "ALTER TABLE tv RENAME TO tv2;\nCREATE TABLE n (id int);", true},
		{"ALTER VIEW SET SCHEMA of a table", "ALTER VIEW public.nokey SET SCHEMA s;\nCREATE TABLE n (id int);", true},
		{"CREATE INDEX on a column the table lacks", "CREATE INDEX i ON t (missing);\nCREATE TABLE n (id int);", true},
		{"DROP FUNCTION by a name several routines have", "CREATE FUNCTION g(integer) RETURNS int LANGUAGE sql AS 'SELECT 1';\nCREATE FUNCTION g(text) RETURNS int LANGUAGE sql AS 'SELECT 1';\nDROP FUNCTION g;\nCREATE TABLE n (id int);", true},
		{"ALTER TABLE IF EXISTS of a table the target lacks", "ALTER TABLE IF EXISTS public.missing ADD CONSTRAINT fk FOREIGN KEY (x) REFERENCES public.t (id);\nALTER TABLE public.t DROP CONSTRAINT t_code_key;", true},
		{"DROP CONSTRAINT of an inherited check", "ALTER TABLE ic DROP CONSTRAINT ip_a_check;", true},
		{"DROP TABLE a view reads without naming a column", "DROP TABLE cvt;\nCREATE TABLE n (id int);", true},
		{"a cascade takes a view reading a table without naming a column", "DROP TABLE cvt CASCADE;\nCREATE TABLE IF NOT EXISTS cv (id int);", true},
		{"OWNED BY a column the table lacks", "CREATE SEQUENCE public.q;\nALTER SEQUENCE public.q OWNED BY t.nope;\nCREATE TABLE n (id int);", true},
		{"CREATE TYPE AS with an attribute twice", "CREATE TYPE public.ct AS (a int, a text);\nCREATE TABLE n (id int);", true},
		{"RENAME CONSTRAINT of a constraint the table lacks", "ALTER TABLE public.t RENAME CONSTRAINT missing TO replacement;\nCREATE TABLE n (id int);", true},
		{"SET SCHEMA of a table the change made to a schema holding its name", "CREATE TABLE public.x (id int PRIMARY KEY);\nCREATE TABLE s.x (id int PRIMARY KEY);\nALTER TABLE public.x SET SCHEMA s;\nCREATE TABLE n (id int);", true},
		{"ALTER TABLE of a view the change made", "CREATE VIEW public.v AS SELECT 1;\nALTER TABLE public.v ADD COLUMN x int;\nCREATE TABLE n (id int);", true},
		{"CREATE OR REPLACE VIEW renaming a column", "CREATE OR REPLACE VIEW tv AS SELECT 'x'::text AS other;\nCREATE TABLE n (id int);", true},
		{"CREATE OR REPLACE VIEW adding a column", "CREATE OR REPLACE VIEW tv AS SELECT code, 1 AS extra FROM t;\nCREATE TABLE n (id int);", true},
		{"CREATE TRIGGER on a relation the target lacks", "CREATE TRIGGER tr BEFORE INSERT ON public.missing FOR EACH ROW EXECUTE FUNCTION public.tf();\nCREATE TABLE n (id int);", true},
		{"COMMENT ON a column the table lacks", "COMMENT ON COLUMN t.nope IS 'x';\nCREATE TABLE n (id int);", true},
		{"GRANT on a table the target lacks", "GRANT SELECT ON public.missing TO PUBLIC;\nCREATE TABLE n (id int);", true},
		{"COMMENT ON a table the target has", "COMMENT ON TABLE t IS 'x';\nCREATE TABLE n (id int);", true},
		{"ALTER FUNCTION RENAME of a signature the change lacks", "CREATE FUNCTION public.g(integer) RETURNS int LANGUAGE sql AS 'SELECT 1';\nALTER FUNCTION public.g(text) RENAME TO h;\nCREATE TABLE n (id int);", true},
		{"CREATE OR REPLACE VIEW changing a column's type", "CREATE OR REPLACE VIEW tv AS SELECT 1 AS code;\nCREATE TABLE n (id int);", true},
		{"CREATE VIEW naming a column twice", "CREATE VIEW v (a, a) AS SELECT 1, 2;\nCREATE TABLE n (id int);", true},
		{"ALTER TABLE DROP COLUMN of a column a table the change made lacks", "CREATE TABLE public.x (id int PRIMARY KEY);\nALTER TABLE public.x DROP COLUMN missing;\nCREATE TABLE n (id int);", true},
		{"OWNED BY a column a table the change made lacks", "CREATE TABLE public.t2 (id int PRIMARY KEY);\nCREATE SEQUENCE public.q;\nALTER SEQUENCE public.q OWNED BY public.t2.missing;\nCREATE TABLE n (id int);", true},
		{"CREATE TRIGGER on a materialized view", "CREATE MATERIALIZED VIEW public.mv AS SELECT 1 AS a;\nCREATE FUNCTION public.tf() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';\nCREATE TRIGGER tr BEFORE INSERT ON public.mv FOR EACH ROW EXECUTE FUNCTION public.tf();\nCREATE TABLE n (id int);", true},
		{"ALTER FUNCTION SET SCHEMA to a schema the target lacks", "CREATE FUNCTION public.g() RETURNS int LANGUAGE sql AS 'SELECT 1';\nALTER FUNCTION public.g() SET SCHEMA missing;\nCREATE TABLE n (id int);", true},
		{"CREATE TYPE AS ENUM with a label twice", "CREATE TYPE public.e2 AS ENUM ('a', 'a');\nCREATE TABLE n (id int);", true},
		{"CREATE OR REPLACE VIEW keeping a column's type", "CREATE OR REPLACE VIEW tv AS SELECT 'x'::text AS code;\nCREATE TABLE n (id int);", true},
		{"DROP CONSTRAINT of a name a table the change made did not give", "CREATE TABLE public.x (id int PRIMARY KEY, CONSTRAINT c CHECK (id > 0));\nALTER TABLE public.x DROP CONSTRAINT missing;\nCREATE TABLE n (id int);", true},
		{"DROP CONSTRAINT of a name a table the change made gave", "CREATE TABLE public.x (id int PRIMARY KEY, CONSTRAINT c CHECK (id > 0));\nALTER TABLE public.x DROP CONSTRAINT c;\nCREATE TABLE n (id int);", true},
		{"CREATE SCHEMA with a view naming an output twice", "CREATE SCHEMA z CREATE VIEW v AS SELECT 1 AS a, 2 AS a;\nCREATE TABLE n (id int);", true},
		{"CREATE INDEX on a table the change made under a taken name", "CREATE TABLE public.t2 (id int PRIMARY KEY);\nCREATE INDEX nokey ON public.t2 (id);\nCREATE TABLE n (id int);", true},
		{"a function the change made follows SET SCHEMA", "CREATE FUNCTION public.f() RETURNS public.nokey LANGUAGE sql AS 'SELECT NULL::public.nokey';\nALTER FUNCTION public.f() SET SCHEMA s;\nDROP FUNCTION s.f();\nDROP TABLE public.nokey;\nCREATE TABLE n (id int);", true},
		{"CREATE SCHEMA with a trigger on a relation the target lacks", "CREATE FUNCTION public.tf() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';\nCREATE SCHEMA z CREATE TRIGGER tr BEFORE INSERT ON public.missing FOR EACH ROW EXECUTE FUNCTION public.tf();\nCREATE TABLE n (id int);", true},
		{"CREATE VIEW with more column names than outputs", "CREATE VIEW v (a, b) AS SELECT 1;\nCREATE TABLE n (id int);", true},
		{"CREATE SEQUENCE OWNED BY a table the target lacks", "CREATE SEQUENCE public.q OWNED BY public.missing.id;\nCREATE TABLE n (id int);", true},
		{"TRUNCATE of a view", "TRUNCATE tv;\nCREATE TABLE n (id int);", true},
		{"a view the change made unqualified, dropped qualified", "CREATE VIEW v AS SELECT * FROM public.nokey;\nDROP VIEW public.v;\nDROP TABLE public.nokey;\nCREATE TABLE n (id int);", true},
		{"CREATE TABLE with a generated column of a column it lacks", "CREATE TABLE public.x (id int PRIMARY KEY, g int GENERATED ALWAYS AS (missing + 1) STORED);\nCREATE TABLE n (id int);", true},
		{"a second CREATE TRIGGER of a name", "CREATE FUNCTION public.tf() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';\nCREATE TRIGGER tr BEFORE INSERT ON nokey FOR EACH ROW EXECUTE FUNCTION public.tf();\nCREATE TRIGGER tr BEFORE INSERT ON nokey FOR EACH ROW EXECUTE FUNCTION public.tf();\nCREATE TABLE n (id int);", false},
		{"CREATE SCHEMA with a sequence owned by a table it lacks", "CREATE SCHEMA z CREATE SEQUENCE q OWNED BY z.missing.id;\nCREATE TABLE n (id int);", true},
		{"a second CREATE TRIGGER of a name on a table", "CREATE FUNCTION public.tf() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';\nCREATE TRIGGER tr BEFORE INSERT ON nokey FOR EACH ROW EXECUTE FUNCTION public.tf();\nCREATE TRIGGER tr BEFORE INSERT ON nokey FOR EACH ROW EXECUTE FUNCTION public.tf();\nCREATE TABLE n (id int);", true},
		{"dropping the key a table the change made was made with", "CREATE TABLE n (id int PRIMARY KEY);\nALTER TABLE n DROP CONSTRAINT n_pkey;", true},
		{"UPDATE of a table the target lacks", "UPDATE public.missing SET id = 1;\nCREATE TABLE n (id int);", true},
		{"SELECT of a table the target lacks", "SELECT * FROM public.missing;\nCREATE TABLE n (id int);", true},
		{"CREATE INDEX on an expression of a column the table lacks", "CREATE INDEX i ON public.t ((missing + 1));\nCREATE TABLE n (id int);", true},
		{"CREATE DOMAIN with a check naming anything but VALUE", "CREATE DOMAIN public.d AS int CHECK (missing > 0);\nCREATE TABLE n (id int);", true},
		{"CREATE DOMAIN with a check of VALUE", "CREATE DOMAIN public.d AS int CHECK (VALUE > 0);\nCREATE TABLE n (id int);", true},
		{"ALTER TYPE ADD VALUE of a type the target lacks", "ALTER TYPE public.missing ADD VALUE 'x';\nCREATE TABLE n (id int);", true},
		{"ALTER TYPE ADD VALUE of an enum", "ALTER TYPE mood ADD VALUE 'x';\nCREATE TABLE n (id int);", true},
		{"DROP INDEX of a renamed constraint index", "ALTER INDEX public.t_pkey RENAME TO x;\nDROP INDEX public.x;\nCREATE TABLE n (id int);", true},
		{"a trigger of CREATE SCHEMA on its own table", "CREATE FUNCTION public.tf() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';\nCREATE SCHEMA z CREATE TABLE t (id int) CREATE TRIGGER tr BEFORE INSERT ON z.t FOR EACH ROW EXECUTE FUNCTION public.tf();", true},
		{"ALTER TYPE SET SCHEMA to a schema the target lacks", "ALTER TYPE mood SET SCHEMA missing;\nCREATE TABLE n (id int);", true},
		{"a renamed key of a table the change made, then dropped", "CREATE TABLE n (id int PRIMARY KEY);\nALTER TABLE n RENAME CONSTRAINT n_pkey TO x;\nALTER TABLE n DROP CONSTRAINT x;", true},
		{"DROP EXTENSION IF EXISTS of an extension the target lacks", "DROP EXTENSION IF EXISTS missing;\nCREATE TABLE n (id int);", true},
		{"DROP EVENT TRIGGER of one the target lacks", "DROP EVENT TRIGGER missing;\nCREATE TABLE n (id int);", true},
		{"COMMENT ON a column a view lacks", "COMMENT ON COLUMN tv.missing IS 'x';\nCREATE TABLE n (id int);", true},
		{"ALTER SCHEMA OWNER of a schema the target lacks", "ALTER SCHEMA missing OWNER TO postgres;\nCREATE TABLE n (id int);", true},
		{"quoted names", "ALTER TABLE \"Quoted\" DROP CONSTRAINT \"Quoted_pkey\";\nCREATE TABLE \"New\" (\"Id\" int);", true},
		{"new tables", "CREATE TABLE n (id int);\nCREATE TABLE k (id int PRIMARY KEY);\nCREATE TABLE k2 (a int, b int, CONSTRAINT k2_pk PRIMARY KEY (a, b));\nCREATE UNLOGGED TABLE ul (a int);\nCREATE TABLE \"Mixed\" (a int);", true},
		{"keys added later", "CREATE TABLE n (id int);\nALTER TABLE n RENAME TO m;\nALTER TABLE m ADD PRIMARY KEY (id);\nCREATE TABLE o (id int NOT NULL);\nCREATE UNIQUE INDEX o_id ON o (id);\nALTER TABLE o ADD CONSTRAINT o_pkey PRIMARY KEY USING INDEX o_id;\nCREATE TABLE q (id int);\nALTER TABLE q ADD COLUMN k serial PRIMARY KEY;", true},
		{"key dropped at the end", "CREATE TABLE n (id int PRIMARY KEY);\nALTER TABLE n DROP CONSTRAINT n_pkey;", false},
		{"table dropped again", "CREATE TABLE n (id int);\nDROP TABLE n;", true},
		{"tables not made by CREATE TABLE", "CREATE TABLE a AS SELECT 1 AS x;\nSELECT 1 AS x INTO b;\nCREATE TABLE IF NOT EXISTS nokey (a int);\nCREATE MATERIALIZED VIEW mv AS SELECT 1 AS x;", true},
		{"partitions", "CREATE TABLE pt (id int) PARTITION BY RANGE (id);\nCREATE TABLE pt1 PARTITION OF pt FOR VALUES FROM (1) TO (10);\nCREATE TABLE pk2 (id int PRIMARY KEY) PARTITION BY LIST (id);\nCREATE TABLE pk2_1 PARTITION OF pk2 FOR VALUES IN (1);\nCREATE TABLE pk2_2 (id int NOT NULL);\nALTER TABLE pk2 ATTACH PARTITION pk2_2 FOR VALUES IN (2);", true},
		{"LIKE", "CREATE TABLE a (LIKE t INCLUDING ALL);\nCREATE TABLE b (LIKE t INCLUDING INDEXES);\nCREATE TABLE c (LIKE t INCLUDING CONSTRAINTS);\nCREATE TABLE d (LIKE t);", true},
		{"inheritance does not pass the key on", "CREATE TABLE child (extra int) INHERITS (t);", true},
		{"typed tables", "CREATE TYPE ty AS (id int);\nCREATE TABLE typed OF ty (PRIMARY KEY (id));\nCREATE TABLE typed2 OF ty;", true},
		{"identity and unique", "CREATE TABLE n (id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY, x int UNIQUE);\nCREATE TABLE m (id int GENERATED ALWAYS AS IDENTITY, x int UNIQUE NOT NULL);", true},
		{"table recreated", "DROP TABLE r;\nDROP TABLE t CASCADE;\nCREATE TABLE t (id int);", true},
		{"keyed tables losing their key", "ALTER TABLE t DROP CONSTRAINT t_pkey CASCADE;\nALTER TABLE p DROP COLUMN id CASCADE;", true},
		{"key dropped with its column", "ALTER TABLE r DROP CONSTRAINT r_t_fk;\nALTER TABLE t DROP COLUMN id;", true},
		{"key column a foreign key references, without CASCADE", "ALTER TABLE t DROP COLUMN id;", true},
		{"key a new foreign key references, without CASCADE", "CREATE TABLE r2 (t_code text REFERENCES t (code));\nALTER TABLE t DROP CONSTRAINT t_code_key;", true},
		{"column a new view reads, without CASCADE", "CREATE VIEW v2 AS SELECT n FROM t;\nALTER TABLE t DROP COLUMN n, DROP CONSTRAINT t_code_key;", true},
		{"key dropped and added back", "ALTER TABLE t DROP CONSTRAINT t_pkey CASCADE;\nALTER TABLE t ADD PRIMARY KEY (id);", true},
		{"column type change keeps the key", "ALTER TABLE r DROP CONSTRAINT r_t_fk;\nALTER TABLE t ALTER COLUMN id TYPE bigint;", true},
		{"no key before the change", "ALTER TABLE nokey ADD PRIMARY KEY (a);\nALTER TABLE nokey DROP CONSTRAINT nokey_pkey;", true},
		{"pg_dump", "SELECT pg_catalog.set_config('search_path', '', false);\nCREATE TABLE public.n (id integer NOT NULL);\nCREATE TABLE public.m (id integer NOT NULL);\nALTER TABLE ONLY public.n ADD CONSTRAINT n_pkey PRIMARY KEY (id);", true},
		{"new schema on the path", "CREATE SCHEMA x;\nCREATE TABLE x.t (id int);\nSET search_path = x, public;\nALTER TABLE t ADD PRIMARY KEY (id);\nCREATE TABLE y (id int);\nALTER TABLE public.t DROP CONSTRAINT t_n_check;", true},
		{"SET ROLE moves $user", "SET ROLE s;\nALTER TABLE t DROP CONSTRAINT t_pkey;\nRESET ROLE;\nALTER TABLE u DROP CONSTRAINT u_k_check;", true},
		{"procedural code", "CREATE TABLE n (id int);\nDO $$ BEGIN ALTER TABLE n ADD PRIMARY KEY (id); END $$;\nALTER TABLE t DROP CONSTRAINT t_n_check;", false},
		{"temporary tables", "CREATE TEMP TABLE t (id int);\nALTER TABLE t ADD CONSTRAINT t_n_check CHECK (id > 0);\nALTER TABLE t DROP CONSTRAINT t_n_check;", false},
		{"rolled back", "BEGIN;\nALTER TABLE t DROP CONSTRAINT t_n_check;\nCREATE TABLE n (id int);\nROLLBACK;", false},
		{"in a transaction", "BEGIN;\nALTER TABLE t DROP CONSTRAINT t_n_check;\nCREATE TABLE n (id int);\nCOMMIT;", true},
		// CREATE SCHEMA runs its sequences, tables, views, indexes, and
		// triggers in that order, whatever the order written.
		{"CREATE SCHEMA makes tables before indexes", "CREATE SCHEMA z CREATE INDEX i ON x (id) CREATE TABLE x (id int);", true},
		{"CREATE SCHEMA makes tables before views", "CREATE SCHEMA z CREATE VIEW v AS SELECT id FROM z.x CREATE TABLE x (id int);", true},
		{"CREATE SCHEMA makes sequences before tables", "CREATE SCHEMA z CREATE TABLE x (id int DEFAULT nextval('z.q')) CREATE SEQUENCE q;", true},
		{"COPY of a relation the target lacks", "COPY public.missing TO STDOUT;\nCREATE TABLE n (id int);", true},
		{"COPY of a query reading a relation the target lacks", "COPY (SELECT id FROM public.missing) TO STDOUT;\nCREATE TABLE n (id int);", true},
		{"COPY TO of a view", "COPY tv TO STDOUT;\nCREATE TABLE n (id int);", true},
		{"dropping a key column of a new table", "CREATE TABLE n (id int PRIMARY KEY, x int);\nALTER TABLE n DROP COLUMN id;", true},
		{"dropping a column of a new composite key", "CREATE TABLE n (a int, b int, PRIMARY KEY (a, b));\nALTER TABLE n DROP COLUMN b, DROP COLUMN IF EXISTS c;", true},
		{"dropping a key column its own foreign key references", "CREATE TABLE n (id int PRIMARY KEY, parent int REFERENCES n (id));\nALTER TABLE n DROP COLUMN id;\nCREATE TABLE m (id int);", true},
		{"dropping a key column a generated column reads", "CREATE TABLE n (id int PRIMARY KEY, g int GENERATED ALWAYS AS (id * 2) STORED);\nALTER TABLE n DROP COLUMN id;\nCREATE TABLE m (id int);", true},
		{"dropping a key column a policy reads", "CREATE TABLE n (id int PRIMARY KEY, x int);\nCREATE POLICY p ON n USING (id > 0);\nALTER TABLE n DROP COLUMN id;\nCREATE TABLE m (id int);", true},
		{"dropping a synced column a new policy reads", "CREATE POLICY p ON t USING (n > 0);\nALTER TABLE t DROP COLUMN n;\nCREATE TABLE m (id int);", true},
		{"a temporary table in a regular schema", "CREATE TEMP TABLE public.x (id int);\nCREATE TABLE n (id int);", true},
		{"a temporary view in a regular schema", "CREATE TEMP VIEW public.v AS SELECT 1;\nCREATE TABLE n (id int);", true},
		{"a temporary table in the temporary schema", "CREATE TEMP TABLE pg_temp.x (id int);\nCREATE TABLE n (id int);", true},
		{"a column type in a schema the target lacks", "CREATE TABLE public.n (x missing.foo);\nCREATE TABLE m (id int);", true},
		{"a column type in pg_catalog", "CREATE TABLE public.n (x pg_catalog.int4);", true},
		{"a function the change made, renamed, is gone under its old name", "CREATE FUNCTION public.g(integer) RETURNS int LANGUAGE sql AS 'SELECT 1';\nALTER FUNCTION public.g(integer) RENAME TO h;\nDROP FUNCTION public.g(integer);\nCREATE TABLE n (id int);", true},
		{"a new table shadowing another on the path", "SET search_path = public, s;\nCREATE TABLE s.n (id int);\nCREATE TABLE public.n (id int PRIMARY KEY);\nDROP TABLE n;", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			conn, err := db.Conn(ctx)
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			resetOracle(t, conn)
			schema := syncSchema(t, conn)

			result, err := Review(ctx, tt.change, review.Options{Rules: []review.Rule{review.DisallowDropConstraint, review.RequirePrimaryKey}}, []review.Target{{Schema: schema, SessionUser: "postgres"}})
			if err != nil {
				t.Fatalf("Review: %v", err)
			}
			want, failedAt := walkTruth(t, conn, tt.change)
			var got []string
			for _, f := range result.Findings {
				if f.Rule == review.Syntax {
					t.Fatalf("syntax: %s", f.Message)
				}
				if failedAt >= 0 && f.Statement == failedAt {
					t.Errorf("reported a finding on the statement the server rejects: %s", f.Message)
					continue
				}
				if failedAt >= 0 && (f.Statement > failedAt || f.Rule == review.RequirePrimaryKey) {
					t.Errorf("reported a finding the change never gets to, after the statement the server rejects: %s", f.Message)
					continue
				}
				got = append(got, fmt.Sprintf("%d %s %q: %s", f.Statement, f.Rule, tt.change[f.Range.Start:f.Range.End], f.Message))
			}
			for _, g := range got {
				if !slices.Contains(want, g) {
					t.Errorf("reported what the server contradicts: %s\nthe server's: %q", g, want)
				}
			}
			if tt.complete && failedAt < 0 {
				for _, w := range want {
					if !slices.Contains(got, w) {
						t.Errorf("missed: %s\ngot: %q", w, got)
					}
				}
			}
		})
	}
}

// resetOracle empties the database down to oracleBase.
func resetOracle(t *testing.T, conn *sql.Conn) {
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
		schemas = append(schemas, `"`+s+`"`)
	}
	rows.Close()
	stmts := []string{
		"ROLLBACK", "RESET SESSION AUTHORIZATION", "RESET ALL", "DISCARD TEMP",
		"DO $$ BEGIN CREATE ROLE s SUPERUSER; EXCEPTION WHEN duplicate_object THEN NULL; END $$",
		"DROP TYPE IF EXISTS ty CASCADE", "DROP TYPE IF EXISTS mood CASCADE",
	}
	if len(schemas) > 0 {
		stmts = append(stmts, "DROP SCHEMA "+strings.Join(schemas, ", ")+" CASCADE")
	}
	stmts = append(stmts, "CREATE SCHEMA public", oracleBase)
	for _, s := range stmts {
		if _, err := conn.ExecContext(ctx, s); err != nil {
			t.Fatalf("%s: %v", s, err)
		}
	}
}

// syncSchema reads the database into the metadata Bytebase would sync:
// the search path, schemas, enum types, views with the columns they read,
// tables, columns, indexes, and keys, foreign keys, and check constraints.
func syncSchema(t *testing.T, conn *sql.Conn) *metadata.DatabaseSchemaMetadata {
	t.Helper()
	ctx := context.Background()
	db := &metadata.DatabaseSchemaMetadata{Name: "omni_review"}
	if err := conn.QueryRowContext(ctx, "SHOW search_path").Scan(&db.SearchPath); err != nil {
		t.Fatal(err)
	}
	query := func(q string, args []any, scan func(*sql.Rows) error) {
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
	}
	schemas := make(map[string]*metadata.SchemaMetadata)
	query(`SELECT nspname FROM pg_namespace WHERE nspname NOT LIKE 'pg\_%' AND nspname <> 'information_schema' ORDER BY nspname`, nil, func(rows *sql.Rows) error {
		s := &metadata.SchemaMetadata{}
		if err := rows.Scan(&s.Name); err != nil {
			return err
		}
		schemas[s.Name] = s
		db.Schemas = append(db.Schemas, s)
		return nil
	})
	query(`SELECT n.nspname, t.typname, ARRAY(SELECT e.enumlabel FROM pg_enum e WHERE e.enumtypid = t.oid ORDER BY e.enumsortorder)::text[]
		FROM pg_type t JOIN pg_namespace n ON n.oid = t.typnamespace WHERE t.typtype = 'e' AND n.nspname NOT LIKE 'pg\_%' AND n.nspname <> 'information_schema'`, nil, func(rows *sql.Rows) error {
		var schema string
		e := &metadata.EnumTypeMetadata{}
		if err := rows.Scan(&schema, &e.Name, (*textArray)(&e.Values)); err != nil {
			return err
		}
		schemas[schema].EnumTypes = append(schemas[schema].EnumTypes, e)
		return nil
	})
	type rel struct {
		oid    int64
		schema string
		table  *metadata.TableMetadata
	}
	var rels []rel
	query(`SELECT c.oid, n.nspname, c.relname FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE c.relkind IN ('r', 'p') AND NOT c.relispartition AND n.nspname NOT LIKE 'pg\_%' AND n.nspname <> 'information_schema' ORDER BY 2, 3`, nil, func(rows *sql.Rows) error {
		r := rel{table: &metadata.TableMetadata{}}
		if err := rows.Scan(&r.oid, &r.schema, &r.table.Name); err != nil {
			return err
		}
		rels = append(rels, r)
		schemas[r.schema].Tables = append(schemas[r.schema].Tables, r.table)
		return nil
	})
	query(`SELECT n.nspname, c.relname, pg_get_viewdef(c.oid) FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE c.relkind = 'v' AND n.nspname NOT LIKE 'pg\_%' AND n.nspname <> 'information_schema' ORDER BY 1, 2`, nil, func(rows *sql.Rows) error {
		var schema string
		v := &metadata.ViewMetadata{}
		if err := rows.Scan(&schema, &v.Name, &v.Definition); err != nil {
			return err
		}
		schemas[schema].Views = append(schemas[schema].Views, v)
		return nil
	})
	for _, s := range db.Schemas {
		for _, v := range s.Views {
			query(`SELECT a.attname, format_type(a.atttypid, a.atttypmod) FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN pg_namespace n ON n.oid = c.relnamespace
				WHERE n.nspname = $1 AND c.relname = $2 AND a.attnum > 0 AND NOT a.attisdropped ORDER BY a.attnum`, []any{s.Name, v.Name}, func(rows *sql.Rows) error {
				c := &metadata.ColumnMetadata{}
				if err := rows.Scan(&c.Name, &c.Type); err != nil {
					return err
				}
				v.Columns = append(v.Columns, c)
				return nil
			})
			query(`SELECT table_schema, table_name, column_name FROM information_schema.view_column_usage WHERE view_schema = $1 AND view_name = $2`, []any{s.Name, v.Name}, func(rows *sql.Rows) error {
				d := &metadata.DependencyColumn{}
				if err := rows.Scan(&d.Schema, &d.Table, &d.Column); err != nil {
					return err
				}
				v.DependencyColumns = append(v.DependencyColumns, d)
				return nil
			})
		}
	}
	for _, r := range rels {
		query(`SELECT i.relname, x.indisunique, pg_get_indexdef(x.indexrelid) FROM pg_index x JOIN pg_class i ON i.oid = x.indexrelid
			WHERE x.indrelid = $1 AND NOT EXISTS (SELECT 1 FROM pg_constraint c WHERE c.conindid = x.indexrelid) ORDER BY 1`, []any{r.oid}, func(rows *sql.Rows) error {
			i := &metadata.IndexMetadata{Type: "btree"}
			if err := rows.Scan(&i.Name, &i.Unique, &i.Definition); err != nil {
				return err
			}
			r.table.Indexes = append(r.table.Indexes, i)
			return nil
		})
		query(`SELECT a.attname, format_type(a.atttypid, a.atttypmod), NOT a.attnotnull FROM pg_attribute a
			WHERE a.attrelid = $1 AND a.attnum > 0 AND NOT a.attisdropped ORDER BY a.attnum`, []any{r.oid}, func(rows *sql.Rows) error {
			c := &metadata.ColumnMetadata{}
			if err := rows.Scan(&c.Name, &c.Type, &c.Nullable); err != nil {
				return err
			}
			r.table.Columns = append(r.table.Columns, c)
			return nil
		})
		columns := `ARRAY(SELECT a.attname FROM unnest(%[1]s) WITH ORDINALITY k(n, i) JOIN pg_attribute a ON a.attrelid = %[2]s AND a.attnum = k.n ORDER BY k.i)::text[]`
		query(`SELECT c.conname, c.contype::text, `+fmt.Sprintf(columns, "c.conkey", "c.conrelid")+`,
				coalesce(fn.nspname, ''), coalesce(f.relname, ''), coalesce(`+fmt.Sprintf(columns, "c.confkey", "c.confrelid")+`, '{}'), coalesce(pg_get_expr(c.conbin, c.conrelid), '')
			FROM pg_constraint c LEFT JOIN pg_class f ON f.oid = c.confrelid LEFT JOIN pg_namespace fn ON fn.oid = f.relnamespace
			WHERE c.conrelid = $1 ORDER BY c.conname`, []any{r.oid}, func(rows *sql.Rows) error {
			var name, kind, refSchema, refTable, check string
			var cols, refCols []string
			if err := rows.Scan(&name, &kind, (*textArray)(&cols), &refSchema, &refTable, (*textArray)(&refCols), &check); err != nil {
				return err
			}
			switch kind {
			case "p", "u":
				r.table.Indexes = append(r.table.Indexes, &metadata.IndexMetadata{Name: name, Expressions: cols, Primary: kind == "p", Unique: true, IsConstraint: true, Type: "btree"})
			case "f":
				r.table.ForeignKeys = append(r.table.ForeignKeys, &metadata.ForeignKeyMetadata{Name: name, Columns: cols, ReferencedSchema: refSchema, ReferencedTable: refTable, ReferencedColumns: refCols})
			case "c":
				r.table.CheckConstraints = append(r.table.CheckConstraints, &metadata.CheckConstraintMetadata{Name: name, Expression: check})
			}
			return nil
		})
	}
	return db
}

// textArray scans a text[] in PostgreSQL's array syntax for plain names.
type textArray []string

func (a *textArray) Scan(src any) error {
	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fmt.Errorf("textArray: %T", src)
	}
	s = strings.TrimSuffix(strings.TrimPrefix(s, "{"), "}")
	*a = nil
	if s != "" {
		*a = strings.Split(s, ",")
	}
	return nil
}

// oracleTable is a table as the server's catalogs describe it.
type oracleTable struct {
	schema    string
	name      string
	kind      string
	partition bool
	keyName   string
	keyCols   []string
}

// oracleConstraint is a constraint as the server's catalogs describe it.
type oracleConstraint struct {
	name  string
	kind  string
	table int64
}

// walkTruth runs the change and returns the findings the server's state
// supports, formatted as TestScanRulesAgainstPostgres formats findings,
// and the index of the first statement the server rejects, or -1.
func walkTruth(t *testing.T, conn *sql.Conn, change string) ([]string, int) {
	t.Helper()
	ctx := context.Background()
	tables := func() map[int64]oracleTable {
		out := make(map[int64]oracleTable)
		rows, err := conn.QueryContext(ctx, `SELECT c.oid, n.nspname, c.relname, c.relkind::text, c.relispartition, coalesce(k.conname, ''),
				coalesce(ARRAY(SELECT a.attname FROM pg_attribute a WHERE a.attrelid = c.oid AND a.attnum = ANY(k.conkey))::text[], '{}')
			FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace LEFT JOIN pg_constraint k ON k.conrelid = c.oid AND k.contype = 'p'
			WHERE c.relkind IN ('r', 'p') AND n.nspname NOT LIKE 'pg\_%' AND n.nspname <> 'information_schema'`)
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var oid int64
			var tb oracleTable
			if err := rows.Scan(&oid, &tb.schema, &tb.name, &tb.kind, &tb.partition, &tb.keyName, (*textArray)(&tb.keyCols)); err != nil {
				t.Fatal(err)
			}
			out[oid] = tb
		}
		return out
	}
	constraints := func() map[int64]oracleConstraint {
		out := make(map[int64]oracleConstraint)
		rows, err := conn.QueryContext(ctx, `SELECT c.oid, c.conname, c.contype::text, c.conrelid FROM pg_constraint c
			JOIN pg_class r ON r.oid = c.conrelid JOIN pg_namespace n ON n.oid = r.relnamespace
			WHERE n.nspname NOT LIKE 'pg\_%' AND n.nspname <> 'information_schema'`)
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var oid int64
			var c oracleConstraint
			if err := rows.Scan(&oid, &c.name, &c.kind, &c.table); err != nil {
				t.Fatal(err)
			}
			out[oid] = c
		}
		return out
	}
	resolve := func(rv *ast.RangeVar) int64 {
		name := `"` + strings.ReplaceAll(rv.Relname, `"`, `""`) + `"`
		if rv.Schemaname != "" {
			name = `"` + strings.ReplaceAll(rv.Schemaname, `"`, `""`) + `".` + name
		}
		var oid sql.NullInt64
		if err := conn.QueryRowContext(ctx, "SELECT to_regclass($1)::oid::int8", name).Scan(&oid); err != nil {
			t.Fatal(err)
		}
		return oid.Int64
	}
	keyed := func(tb oracleTable) bool {
		return (tb.kind == "r" || tb.kind == "p") && !tb.partition && tb.keyName != ""
	}

	initial := constraints()
	keyedAtStart := make(map[int64]bool)
	for oid, tb := range tables() {
		keyedAtStart[oid] = keyed(tb)
	}
	var want []string
	created := make(map[int64]string)
	unkeyed := make(map[int64]string)
	index := 0
	for _, seg := range pg.Split(change) {
		if seg.Empty() {
			continue
		}
		i := index
		index++
		list, err := parser.ParseRange(change, seg.ByteStart, seg.ByteEnd)
		if err != nil {
			t.Fatalf("parse: %v", err)
		}
		node := list.Items[0].(*ast.RawStmt).Stmt
		alter, _ := node.(*ast.AlterTableStmt)
		var altered int64
		if alter != nil && ast.ObjectType(alter.ObjType) == ast.OBJECT_TABLE {
			altered = resolve(alter.Relation)
		}
		tablesBefore, consBefore := tables(), constraints()
		if _, err := conn.ExecContext(ctx, change[seg.ByteStart:seg.ByteEnd]); err != nil {
			t.Logf("statement %d rejected: %v", i, err)
			return want, i
		}
		tablesAfter, consAfter := tables(), constraints()
		at := func(rng review.Range) string { return change[rng.Start:rng.End] }

		if altered != 0 {
			tb := tablesBefore[altered]
			for _, item := range alter.Cmds.Items {
				cmd := item.(*ast.AlterTableCmd)
				if ast.AlterTableType(cmd.Subtype) != ast.AT_DropConstraint || (tb.kind != "r" && tb.kind != "p") || tb.partition {
					continue
				}
				for oid, c := range consBefore {
					_, gone := consAfter[oid]
					_, existed := initial[oid]
					kind := map[string]string{"p": "primary key", "f": "foreign key", "u": "unique constraint", "c": "check constraint"}[c.kind]
					if c.table == altered && c.name == cmd.Name && existed && !gone && kind != "" {
						want = append(want, fmt.Sprintf("%d %s %q: drops %s %s of %s", i, review.DisallowDropConstraint, at(rangeOf(cmd.Loc)), kind, ident(c.name), relation(alter.Relation)))
					}
				}
			}
			if before := tablesBefore[altered]; keyedAtStart[altered] && keyed(before) {
				if after, ok := tablesAfter[altered]; ok && after.keyName == "" {
					rng := keyDropAnchor(alter, before.keyName, before.keyCols)
					unkeyed[altered] = fmt.Sprintf("%d %s %q: removes the primary key of %s, and the change adds none back", i, review.RequirePrimaryKey, at(rng), relation(alter.Relation))
				}
			}
		}
		if create, ok := node.(*ast.CreateStmt); ok && create.Partbound == nil {
			for oid := range tablesAfter {
				if _, existed := tablesBefore[oid]; !existed {
					created[oid] = fmt.Sprintf("%d %s %q: creates table %s without a primary key", i, review.RequirePrimaryKey, at(rangeOf(create.Loc)), relation(create.Relation))
				}
			}
		}
		// CREATE SCHEMA ... CREATE TABLE names the table in the new schema.
		if schema, ok := node.(*ast.CreateSchemaStmt); ok && schema.SchemaElts != nil {
			for oid, tb := range tablesAfter {
				if _, existed := tablesBefore[oid]; existed {
					continue
				}
				for _, elt := range schema.SchemaElts.Items {
					if create, ok := elt.(*ast.CreateStmt); ok && create.Partbound == nil && create.Relation.Relname == tb.name {
						created[oid] = fmt.Sprintf("%d %s %q: creates table %s without a primary key", i, review.RequirePrimaryKey, at(rangeOf(create.Loc)), relation(&ast.RangeVar{Schemaname: tb.schema, Relname: tb.name}))
					}
				}
			}
		}
	}
	end := tables()
	for _, found := range []map[int64]string{created, unkeyed} {
		for oid, finding := range found {
			if tb, ok := end[oid]; ok && (tb.kind == "r" || tb.kind == "p") && !tb.partition && tb.keyName == "" {
				want = append(want, finding)
			}
		}
	}
	return want, -1
}

// keyDropAnchor is the subcommand that drops the primary key: the first
// DROP CONSTRAINT naming it or DROP COLUMN of one of its columns, or else
// the statement.
func keyDropAnchor(v *ast.AlterTableStmt, name string, columns []string) review.Range {
	for _, item := range v.Cmds.Items {
		cmd := item.(*ast.AlterTableCmd)
		switch ast.AlterTableType(cmd.Subtype) {
		case ast.AT_DropConstraint:
			if cmd.Name == name {
				return rangeOf(cmd.Loc)
			}
		case ast.AT_DropColumn:
			if slices.Contains(columns, cmd.Name) {
				return rangeOf(cmd.Loc)
			}
		}
	}
	return rangeOf(v.Loc)
}
