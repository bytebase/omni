package compat

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// oracleRaised is what OracleParse raises for an error, as go-ora reports it.
func oracleRaised(code, position int, compile string) error {
	return fmt.Errorf("ORA-20999: omni-parse code=%d position=%d compile=%s\nORA-06512: at line 52", code, position, compile)
}

// TestOracleVerdictReadsParseErrors pins how a verdict reads OracleParse's
// answers. Each answer, the probe's included, is what gvenzl/oracle-free:
// 23.26.3-slim-faststart gave, as a user that may create tables, views,
// sequences, synonyms, materialized views, types, triggers and procedures,
// in a schema that holds the table a.
func TestOracleVerdictReadsParseErrors(t *testing.T) {
	type answer struct {
		code, position int
		compile        string
	}
	for _, c := range []struct {
		name      string
		statement string
		first     *answer // nil: parsed
		probe     *answer // nil: not asked
		want      string
	}{
		{"parsed", "SELECT 1 FROM dual", nil, nil, Accept},
		{"a syntax error, which the probe meets first", "SELECT FROM dual", &answer{936, 7, ""}, &answer{936, 7, ""}, "reject@7"},
		{"23ai's naming of the word", "SELECT 1 FROM dual dual dual", &answer{3049, 24, ""}, &answer{3049, 24, ""}, "reject@24"},
		{"the lexer's", "SELECT 'a", &answer{1756, 7, ""}, &answer{1756, 7, ""}, "reject@7"},
		{"a byte offset, not a character's", "SELECT 'é' FROM dual dual dual", &answer{3049, 27, ""}, &answer{3049, 27, ""}, "reject@27"},
		{"the parser's error at the end", "SELECT 1 FROM", &answer{903, 13, ""}, nil, "reject@13"},
		{"a name looked up after the parse", "SELECT zz FROM dual", &answer{904, 7, ""}, &answer{911, 21, ""}, Accept},
		{"a name taken, after the parse", "CREATE TABLE a (a INT)", &answer{955, 13, ""}, &answer{911, 24, ""}, Accept},
		{"a missing clause, at the end and after the parse", "CREATE PLUGGABLE DATABASE a FROM b", &answer{65016, 34, ""}, nil, Accept},
		{"a statement's own error at the end", "ALTER TABLE a", &answer{2210, 13, ""}, nil, "reject@13"},
		{"a name missing at the end", "ALTER DATABASE DEFAULT TABLESPACE", &answer{2216, 33, ""}, nil, "reject@33"},
		{"a lookup's code at the end, where there is no name", "ALTER DOMAIN IF", &answer{904, 15, ""}, nil, "reject@15"},
		{"the server's setup, at the end", "CREATE BIGFILE TABLESPACE a", &answer{2199, 27, ""}, nil, Accept},
		{"a key the table lacks, at the end", "CREATE TABLE a ENABLE NOVALIDATE PRIMARY KEY", &answer{2432, 44, ""}, nil, Accept},
		{"a privilege", "CREATE USER a IDENTIFIED BY a", &answer{1031, 0, ""}, nil, Accept},
		{"a privilege, at a byte outside the statement", "AUDIT POLICY zz", &answer{1031, 27, ""}, nil, Accept},
		{"a trigger on a missing table", "CREATE TRIGGER a BEFORE INSERT ON zz BEGIN NULL; END;", &answer{942, 34, ""}, nil, Accept},
		{"a refusal while parsing", "AUDIT ALL ON a", &answer{46401, 6, ""}, &answer{46401, 6, ""}, "reject@6"},
		{"a parameter looked up while parsing", "ALTER SYSTEM SET a = 1", &answer{2065, 17, ""}, &answer{2065, 17, ""}, "reject@17"},
		{"a refusal at a byte outside the statement", "AUDIT a", &answer{46401, 30, ""}, &answer{46401, 30, ""}, Reject},
		{"a unit with a PL/SQL syntax error", "CREATE PROCEDURE a IS BEGIN NULL END;", &answer{24344, 0, "PLS-00103"}, nil, Reject},
		{"a unit with a name it cannot resolve", "CREATE PROCEDURE a IS BEGIN zz; END;", &answer{24344, 0, "PLS-00201,PLS-00000"}, nil, Accept},
		{"a unit with SQL on a missing table", "CREATE PROCEDURE a IS BEGIN INSERT INTO zz VALUES (1); END;", &answer{24344, 12, "ORA-00942,PLS-00000"}, nil, Accept},
		{"a unit with a SQL syntax error", "CREATE PROCEDURE a IS BEGIN INSERT INTO a VALUES (1) a a; END;", &answer{24344, 25, "ORA-03049,PLS-00000"}, nil, Reject},
		{"a block with SQL on a missing table", "BEGIN INSERT INTO zz VALUES (1); END;", &answer{6550, 18, "ORA-00942"}, nil, Accept},
		{"a block with SQL of mismatched types", "BEGIN INSERT INTO a SELECT 1 FROM dual UNION ALL SELECT 'x' FROM dual; END;", &answer{6550, 27, "ORA-01790"}, nil, Accept},
		{"a block calling a procedure named like an error", `BEGIN "PLS-00103"; END;`, &answer{6550, 7, "PLS-00201"}, nil, Accept},
		{"a block with a SQL syntax error", "BEGIN INSERT INTO a VALUES (1) a a; END;", &answer{6550, 31, "ORA-03049"}, nil, "reject@31"},
		{"a type with a PL/SQL syntax error", "CREATE TYPE a AS OBJECT (a INT", &answer{24344, 0, "PLS-00103"}, nil, Reject},
		{"a block with a PL/SQL syntax error", "BEGIN NULL END;", &answer{6550, 11, "PLS-00103"}, nil, "reject@11"},
		{"a block with a name it cannot resolve", "BEGIN zz; END;", &answer{6550, 6, "PLS-00201"}, nil, Accept},
		{"a probe answered at the tail by another error", "ALTER MATERIALIZED VIEW LOG ON a MERGE PARTITIONS a, a", &answer{12002, 0, ""}, &answer{14006, 56, ""}, Accept},
		{"a probe the tail stops early, after a lookup", "ALTER INDEX a UNUSABLE", &answer{1418, 0, ""}, &answer{14048, 14, ""}, Accept},
		{"a probe the tail stops early, after a syntax error", "SELECT 1 FROM dual a b", &answer{933, 21, ""}, &answer{923, 7, ""}, "reject@21"},
		{"a syntax error in a statement lexed whole first", "ALTER MATERIALIZED VIEW a zzz yyy", &answer{3049, 26, ""}, &answer{911, 35, ""}, "reject@26"},
		{"a statement's own syntax error, lexed whole first", "ALTER MATERIALIZED VIEW a PARAMETERS", &answer{14051, 26, ""}, &answer{911, 38, ""}, "reject@26"},
		{"a number not valid where it stands, lexed whole first", "ALTER MATERIALIZED VIEW a 1", &answer{3047, 26, ""}, &answer{911, 29, ""}, "reject@26"},
		{"a privilege no grammar knows, lexed whole first", "GRANT a a TO a", &answer{990, 6, ""}, &answer{911, 16, ""}, "reject@6"},
		{"no column list, lexed whole first", "CREATE TABLE a NO INMEMORY", &answer{906, 15, ""}, &answer{911, 28, ""}, "reject@15"},
		{"a role looked up, lexed whole first", "GRANT zzz TO a", &answer{1924, 6, ""}, &answer{911, 16, ""}, Accept},
		{"a data type looked up after the parse", "CREATE TABLE a(a UNITE)", &answer{902, 17, ""}, &answer{911, 25, ""}, Accept},
		{"a tablespace looked up after the parse", "ALTER FLASHBACK ARCHIVE a ADD TABLESPACE a", &answer{959, 41, ""}, &answer{911, 44, ""}, Accept},
		{"a cluster looked up after the parse", "ALTER CLUSTER a ALLOCATE EXTENT", &answer{943, 0, ""}, &answer{911, 33, ""}, Accept},
		{"a view column to name, after the parse", "CREATE VIEW a AS SELECT 1 FROM dual WITH READ ONLY", &answer{998, 24, ""}, &answer{911, 53, ""}, Accept},
		{"a database link looked up after the parse", "SELECT * FROM a@a zzz", &answer{2019, 16, ""}, &answer{911, 23, ""}, Accept},
		{"a 23ai expression syntax error, lexed whole first", "ALTER TABLE t MODIFY LOB(a(RETENTION)", &answer{3059, 25, ""}, &answer{911, 39, ""}, "reject@25"},
		{"MERGE missing USING, lexed whole first", "MERGE INTO a", &answer{38107, 0, ""}, &answer{911, 14, ""}, "reject@0"},
		{"ADMINISTER KEY MANAGEMENT SET missing clause", "ADMINISTER KEY MANAGEMENT SET", &answer{46601, 26, ""}, &answer{911, 31, ""}, "reject@26"},
		{"ADMINISTER KEY MANAGEMENT CREATE missing clause", "ADMINISTER KEY MANAGEMENT CREATE", &answer{46603, 26, ""}, &answer{911, 34, ""}, "reject@26"},
		{"CREATE SPFILE in PDB at the end", "CREATE SPFILE FROM PFILE", &answer{65168, 24, ""}, nil, Accept},
	} {
		var asked []string
		parse := func(text string) error {
			asked = append(asked, text)
			a := c.first
			if len(asked) == 2 {
				a = c.probe
				if text != c.statement+OracleProbeTail {
					t.Fatalf("%s: probed %q", c.name, text)
				}
			}
			if len(asked) > 2 || a == nil && len(asked) == 2 {
				t.Fatalf("%s: asked %q", c.name, asked)
			}
			if a == nil {
				return nil
			}
			return oracleRaised(a.code, a.position, a.compile)
		}
		got, err := OracleVerdict(c.statement, parse)
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}
		if got != c.want {
			t.Errorf("%s: got %s, want %s", c.name, got, c.want)
		}
		if probed := len(asked) == 2; probed != (c.probe != nil) {
			t.Errorf("%s: probed %v", c.name, probed)
		}
	}
}

// TestOracleVerdictAsksOnTheFixturesObject pins the stand-ins: Oracle looks
// up the table of ALTER TABLE a and the view of ALTER MATERIALIZED VIEW a
// before it reads the rest, so the verdict is the one on the fixture's t or
// v. Each answer is the pinned server's, with the fixture
// oracle/internal/testcontainer.JudgeFixture builds.
func TestOracleVerdictAsksOnTheFixturesObject(t *testing.T) {
	type answer struct{ code, position int }
	parsed := (*answer)(nil)
	for _, c := range []struct {
		name      string
		statement string
		answers   map[string]*answer // every text the verdict may ask
		want      string
	}{
		{"an option the table takes", "ALTER TABLE a LOGGING", map[string]*answer{
			"ALTER TABLE a LOGGING": {942, 0},
			"ALTER TABLE t LOGGING": parsed,
		}, Accept},
		{"a syntax error after the view's name", "ALTER MATERIALIZED VIEW a COMPILE zzz", map[string]*answer{
			"ALTER MATERIALIZED VIEW a COMPILE zzz":                   {12003, 0},
			"ALTER MATERIALIZED VIEW v COMPILE zzz":                   {3049, 34},
			"ALTER MATERIALIZED VIEW v COMPILE zzz" + OracleProbeTail: {911, 39},
		}, "reject@34"},
		{"no option, said at byte 0", "ALTER MATERIALIZED VIEW a NOMAPPING", map[string]*answer{
			"ALTER MATERIALIZED VIEW a NOMAPPING":                   {12003, 0},
			"ALTER MATERIALIZED VIEW v NOMAPPING":                   {2210, 0},
			"ALTER MATERIALIZED VIEW v NOMAPPING" + OracleProbeTail: {911, 37},
		}, "reject@0"},
		{"the table looked up at its name", "ALTER TABLE a ADD PARTITION zzz", map[string]*answer{
			"ALTER TABLE a ADD PARTITION zzz":                   {942, 12},
			"ALTER TABLE t ADD PARTITION zzz":                   {14501, 12},
			"ALTER TABLE t ADD PARTITION zzz" + OracleProbeTail: {911, 33},
		}, Accept},
		{"an option the table's kind does not take", "ALTER TABLE a OVERFLOW ALLOCATE EXTENT", map[string]*answer{
			"ALTER TABLE a OVERFLOW ALLOCATE EXTENT": {942, 0},
			"ALTER TABLE t OVERFLOW ALLOCATE EXTENT": {1735, 0},
		}, Accept},
		{"a syntax error before the lookup", "ALTER TABLE a zzz", map[string]*answer{
			"ALTER TABLE a zzz":                   {1735, 14},
			"ALTER TABLE a zzz" + OracleProbeTail: {1735, 14},
		}, "reject@14"},
		{"a view's partition looked up as a table", "ALTER MATERIALIZED VIEW a ADD PARTITION", map[string]*answer{
			"ALTER MATERIALIZED VIEW a ADD PARTITION": {942, 24},
		}, Accept},
		{"a materialized view log, which has no stand-in", "ALTER MATERIALIZED VIEW LOG ON a ADD(a)", map[string]*answer{
			"ALTER MATERIALIZED VIEW LOG ON a ADD(a)": {942, 0},
		}, Accept},
		{"ALTER TYPE on object fixture y", "ALTER TYPE a ADD ATTRIBUTE a NUMBER", map[string]*answer{
			"ALTER TYPE a ADD ATTRIBUTE a NUMBER": {4043, 0},
			"ALTER TYPE y ADD ATTRIBUTE a NUMBER": parsed,
		}, Accept},
		{"ALTER TYPE MODIFY LIMIT on varray fixture w", "ALTER TYPE a MODIFY LIMIT 1", map[string]*answer{
			"ALTER TYPE a MODIFY LIMIT 1": {4043, 0},
			"ALTER TYPE w MODIFY LIMIT 2": parsed,
		}, Accept},
		{"ALTER TYPE MODIFY ELEMENT TYPE on varray fixture w", "ALTER TYPE a MODIFY ELEMENT TYPE NUMBER", map[string]*answer{
			"ALTER TYPE a MODIFY ELEMENT TYPE NUMBER":       {4043, 0},
			"ALTER TYPE w MODIFY ELEMENT TYPE VARCHAR2(10)": parsed,
		}, Accept},
		{"ALTER DIMENSION ADD LEVEL on fixture d", "ALTER DIMENSION a ADD LEVEL a IS(a.a) zzz", map[string]*answer{
			"ALTER DIMENSION a ADD LEVEL a IS(a.a) zzz":                   {30333, 16},
			"ALTER DIMENSION d ADD LEVEL a IS(a.a) zzz":                   {3049, 38},
			"ALTER DIMENSION d ADD LEVEL a IS(a.a) zzz" + OracleProbeTail: {3049, 38},
		}, "reject@38"},
		{"ALTER DIMENSION DROP HIERARCHY on fixture e", "ALTER DIMENSION a DROP HIERARCHY a(a CHILD OF a)", map[string]*answer{
			"ALTER DIMENSION a DROP HIERARCHY a(a CHILD OF a)":                   {30333, 16},
			"ALTER DIMENSION e DROP HIERARCHY a(a CHILD OF a)":                   {3048, 34},
			"ALTER DIMENSION e DROP HIERARCHY a(a CHILD OF a)" + OracleProbeTail: {3048, 34},
		}, "reject@34"},
		{"INMEMORY JOIN GROUP multiple table stand-ins", "CREATE INMEMORY JOIN GROUP a(a(a), a(a))", map[string]*answer{
			"CREATE INMEMORY JOIN GROUP a(a(a), a(a))": {903, 29},
			"CREATE INMEMORY JOIN GROUP a(t(a), a(a))": {903, 35},
			"CREATE INMEMORY JOIN GROUP a(t(a), t(a))": parsed,
		}, Accept},
		{"CREATE PLUGGABLE DATABASE USING 's' accepted via 's.xml'", "CREATE PLUGGABLE DATABASE a USING 's'", map[string]*answer{
			"CREATE PLUGGABLE DATABASE a USING 's'":     {65125, 34},
			"CREATE PLUGGABLE DATABASE a USING 's.xml'": {65018, 41},
		}, Accept},
		{"CREATE PLUGGABLE DATABASE USING 's' zzz remaps byte offset", "CREATE PLUGGABLE DATABASE a USING 's' zzz", map[string]*answer{
			"CREATE PLUGGABLE DATABASE a USING 's' zzz":                       {65125, 34},
			"CREATE PLUGGABLE DATABASE a USING 's.xml' zzz":                   {922, 42},
			"CREATE PLUGGABLE DATABASE a USING 's.xml' zzz" + OracleProbeTail: {922, 42},
		}, "reject@38"},
		{"CREATE PLUGGABLE DATABASE chained USING and SOURCE_FILE_DIRECTORY", "CREATE PLUGGABLE DATABASE a USING 's' SOURCE_FILE_DIRECTORY = 's' zzz", map[string]*answer{
			"CREATE PLUGGABLE DATABASE a USING 's' SOURCE_FILE_DIRECTORY = 's' zzz":                        {65125, 34},
			"CREATE PLUGGABLE DATABASE a USING 's.xml' SOURCE_FILE_DIRECTORY = 's' zzz":                    {65155, 66},
			"CREATE PLUGGABLE DATABASE a USING 's.xml' SOURCE_FILE_DIRECTORY = '/s' zzz":                   {922, 71},
			"CREATE PLUGGABLE DATABASE a USING 's.xml' SOURCE_FILE_DIRECTORY = '/s' zzz" + OracleProbeTail: {922, 71},
		}, "reject@66"},
		{"SHARING = METADATA stand-in remaps byte offset", "CREATE ANALYTIC VIEW a SHARING = METADATA zzz USING a", map[string]*answer{
			"CREATE ANALYTIC VIEW a SHARING = METADATA zzz USING a":               {65021, 33},
			"CREATE ANALYTIC VIEW a SHARING = NONE zzz USING a":                   {2000, 38},
			"CREATE ANALYTIC VIEW a SHARING = NONE zzz USING a" + OracleProbeTail: {2000, 38},
		}, "reject@42"},
		{"CREATE PFILE FROM stand-in remaps byte offset", "CREATE PFILE FROM SPFILE zzz", map[string]*answer{
			"CREATE PFILE FROM SPFILE zzz":                         {65168, 13},
			"CREATE PFILE = 's' FROM SPFILE zzz":                   {3049, 31},
			"CREATE PFILE = 's' FROM SPFILE zzz" + OracleProbeTail: {3049, 31},
		}, "reject@25"},
		{"ALTER TABLE NO DELETE UNTIL 1 chained with table stand-in", "ALTER TABLE a NO DELETE UNTIL 1 DAY AFTER", map[string]*answer{
			"ALTER TABLE a NO DELETE UNTIL 1 DAY AFTER":  {942, 0},
			"ALTER TABLE t NO DELETE UNTIL 1 DAY AFTER":  {5741, 0},
			"ALTER TABLE t NO DELETE UNTIL 16 DAY AFTER": {2000, 42},
		}, "reject@41"},
		{"ALTER DATABASE RECOVER UNTIL TIME placeholder", "ALTER DATABASE RECOVER STANDBY DATABASE UNTIL TIME 'YYYY-MM-DD:HH24:MI:SS'", map[string]*answer{
			"ALTER DATABASE RECOVER STANDBY DATABASE UNTIL TIME 'YYYY-MM-DD:HH24:MI:SS'": {1841, 23},
			"ALTER DATABASE RECOVER STANDBY DATABASE UNTIL TIME '2026-01-01:00:00:00'":   {1031, 0},
		}, Accept},
		{"ALTER AUDIT POLICY CONDITION 's' stand-in", "ALTER AUDIT POLICY a CONDITION 's' EVALUATE PER STATEMENT", map[string]*answer{
			"ALTER AUDIT POLICY a CONDITION 's' EVALUATE PER STATEMENT":   {46368, 0},
			"ALTER AUDIT POLICY a CONDITION '1=1' EVALUATE PER STATEMENT": {1031, 0},
		}, Accept},
		{"CREATE AUDIT POLICY WHEN 's' stand-in", "CREATE AUDIT POLICY a WHEN 's' EVALUATE PER STATEMENT", map[string]*answer{
			"CREATE AUDIT POLICY a WHEN 's' EVALUATE PER STATEMENT":                                         {46373, 0},
			"CREATE AUDIT POLICY a ACTIONS SELECT ON t WHEN '1=1' EVALUATE PER STATEMENT":                   {41732, 40},
			"CREATE AUDIT POLICY a ACTIONS SELECT ON t WHEN '1=1' EVALUATE PER STATEMENT" + OracleProbeTail: {911, 0},
		}, Accept},
		{"CREATE AUDIT POLICY ONLY TOPLEVEL stand-in", "CREATE AUDIT POLICY a ONLY TOPLEVEL", map[string]*answer{
			"CREATE AUDIT POLICY a ONLY TOPLEVEL":                                       {46373, 0},
			"CREATE AUDIT POLICY a ACTIONS SELECT ON t ONLY TOPLEVEL":                   {41732, 40},
			"CREATE AUDIT POLICY a ACTIONS SELECT ON t ONLY TOPLEVEL" + OracleProbeTail: {911, 0},
		}, Accept},
		{"missing schema a. stand-in remaps byte offset", "CREATE MATERIALIZED VIEW a.zzz a AS SELECT 1 FROM dual", map[string]*answer{
			"CREATE MATERIALIZED VIEW a.zzz a AS SELECT 1 FROM dual":                        {1435, 0},
			"CREATE MATERIALIZED VIEW SYSTEM.zzz a AS SELECT 1 FROM dual":                   {2000, 36},
			"CREATE MATERIALIZED VIEW SYSTEM.zzz a AS SELECT 1 FROM dual" + OracleProbeTail: {2000, 36},
		}, "reject@31"},
		{"missing schema a. on column DOMAIN stand-in", "CREATE TABLE a(a DOMAIN a.zzz a)", map[string]*answer{
			"CREATE TABLE a(a DOMAIN a.zzz a)":                        {1917, 0},
			"CREATE TABLE a(a DOMAIN SYSTEM.zzz a)":                   {3099, 35},
			"CREATE TABLE a(a DOMAIN SYSTEM.zzz a)" + OracleProbeTail: {3099, 35},
		}, "reject@30"},
	} {
		parse := func(text string) error {
			a, ok := c.answers[text]
			if !ok {
				t.Fatalf("%s: asked %q", c.name, text)
			}
			if a == nil {
				return nil
			}
			return oracleRaised(a.code, a.position, "")
		}
		got, err := OracleVerdict(c.statement, parse)
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}
		if got != c.want {
			t.Errorf("%s: got %s, want %s", c.name, got, c.want)
		}
	}
}

// TestOracleVerdictRefusesWhatIsNoVerdict holds the verdict to failing on
// an answer it cannot read, rather than recording an agreement.
func TestOracleVerdictRefusesWhatIsNoVerdict(t *testing.T) {
	lost := errors.New("context deadline exceeded")
	for _, c := range []struct {
		name         string
		first, probe error
		want         string
	}{
		{"an error OracleParse did not raise", lost, nil, "deadline"},
		{"a probe answered so", oracleRaised(904, 7, ""), lost, "deadline"},
		{"a probe the server took", oracleRaised(904, 7, ""), nil, "took the probe's tail"},
		{"a unit with no errors listed", oracleRaised(24344, 0, ""), nil, "lists no PL/SQL error"},
		{"a block with no errors listed", oracleRaised(6550, 0, ""), nil, "lists no PL/SQL error"},
	} {
		calls := 0
		_, err := OracleVerdict("SELECT zz FROM dual", func(string) error {
			calls++
			if calls == 1 {
				return c.first
			}
			return c.probe
		})
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: got %v, want an error with %q", c.name, err, c.want)
		}
	}
}

// TestOracleVerdictReadsOraclesOwnFailures holds the verdict to reading a
// statement Oracle fails on in its own code, on two sessions, as a
// rejection, and asking again after one failure. The failures are those of
// the pinned server: CREATE JAVA SOURCE USING kills the session's process
// (ORA-07445 in the trace, end-of-file to the client), and a near-miss of
// ALTER MATERIALIZED VIEW LOG ... SPLIT PARTITION raises ORA-00600.
func TestOracleVerdictReadsOraclesOwnFailures(t *testing.T) {
	eof := errors.New("EOF")
	normal := errors.New("ORA-00000: normal, successful completion\nORA-06512: at line 33")
	internal := errors.New("ORA-00600: internal error code, arguments: [12333], [PTI], [prspo_set_up]\nORA-06512: at \"SYS.DBMS_SQL\", line 1251")
	for _, c := range []struct {
		name    string
		answers []error
		want    string
	}{
		{"a dropped session, twice", []error{eof, eof}, Reject},
		{"an internal error, twice", []error{internal, internal}, Reject},
		{"a session failing the block's next query, twice", []error{normal, normal}, Reject},
		{"a dropped session, then parsed", []error{eof, nil}, Accept},
		{"no answer, twice", []error{fmt.Errorf("after 1m: %w", ErrNoAnswer), ErrNoAnswer}, Reject},
		{"a dropped session, then a syntax error", []error{eof, oracleRaised(936, 7, ""), oracleRaised(936, 7, "")}, "reject@7"},
		{"a probe Oracle fails on, after a parser's error", []error{oracleRaised(936, 7, ""), internal}, "reject@7"},
		{"a probe Oracle fails on, after a lookup", []error{oracleRaised(904, 7, ""), internal}, Accept},
	} {
		calls := 0
		got, err := OracleVerdict("SELECT FROM zz", func(string) error {
			calls++
			if calls > len(c.answers) {
				t.Fatalf("%s: asked %d times", c.name, calls)
			}
			return c.answers[calls-1]
		})
		if err != nil || got != c.want {
			t.Errorf("%s: got %q, %v; want %s", c.name, got, err, c.want)
		}
		if calls != len(c.answers) {
			t.Errorf("%s: asked %d times, want %d", c.name, calls, len(c.answers))
		}
	}
}

// TestOracleVerdictOnly pins judge revision 5's coarse accept/reject verdicts,
// new stand-ins, error-code classifications, ORA-00600 MV-log stand-in,
// structural PL/SQL compiler errors, and unclosed-'(' detection.
func TestOracleVerdictOnly(t *testing.T) {
	type answer struct {
		code, position int
		compile        string
		rawErr         error
	}
	parsed := (*answer)(nil)
	internal := errors.New("ORA-00600: internal error code, arguments: [12333], [PTI], [prspo_set_up]\nORA-06512: at \"SYS.DBMS_SQL\", line 1251")

	for _, c := range []struct {
		name      string
		statement string
		answers   map[string]*answer
		want      string
	}{
		{"coarse strips byte offset", "SELECT FROM dual", map[string]*answer{
			"SELECT FROM dual":                   {936, 7, "", nil},
			"SELECT FROM dual" + OracleProbeTail: {936, 7, "", nil},
		}, Reject},
		{"DROP TABLESPACE SET stand-in accepts valid", "DROP TABLESPACE SET a INCLUDING CONTENTS", map[string]*answer{
			"DROP TABLESPACE SET a INCLUDING CONTENTS": {2557, 20, "", nil},
			"DROP TABLESPACE a INCLUDING CONTENTS":     {1031, 0, "", nil},
		}, Accept},
		{"DROP TABLESPACE SET stand-in rejects trailing zzz", "DROP TABLESPACE SET a INCLUDING CONTENTS zzz", map[string]*answer{
			"DROP TABLESPACE SET a INCLUDING CONTENTS zzz":               {2557, 20, "", nil},
			"DROP TABLESPACE a INCLUDING CONTENTS zzz":                   {2173, 37, "", nil},
			"DROP TABLESPACE a INCLUDING CONTENTS zzz" + OracleProbeTail: {2173, 37, "", nil},
		}, Reject},
		{"SEQUENCE SHARD stand-in accepts valid", "CREATE SEQUENCE a SHARD", map[string]*answer{
			"CREATE SEQUENCE a SHARD": {2511, 18, "", nil},
			"CREATE SEQUENCE a SCALE": parsed,
		}, Accept},
		{"SEQUENCE NOSHARD stand-in rejects trailing zzz", "CREATE SEQUENCE a NOSHARD zzz", map[string]*answer{
			"CREATE SEQUENCE a NOSHARD zzz":                   {2511, 18, "", nil},
			"CREATE SEQUENCE a NOSCALE zzz":                   {3049, 26, "", nil},
			"CREATE SEQUENCE a NOSCALE zzz" + OracleProbeTail: {3049, 26, "", nil},
		}, Reject},
		{"IDENTIFIED GLOBALLY stand-in accepts ROLE AS string", "CREATE ROLE a IDENTIFIED GLOBALLY AS 's'", map[string]*answer{
			"CREATE ROLE a IDENTIFIED GLOBALLY AS 's'":   {439, 0, "", nil},
			"CREATE USER a IDENTIFIED EXTERNALLY AS 's'": {1031, 0, "", nil},
		}, Accept},
		{"IDENTIFIED GLOBALLY stand-in rejects missing string", "CREATE ROLE a IDENTIFIED GLOBALLY AS", map[string]*answer{
			"CREATE ROLE a IDENTIFIED GLOBALLY AS":                     {439, 0, "", nil},
			"CREATE USER a IDENTIFIED EXTERNALLY AS":                   {28025, 38, "", nil},
			"CREATE USER a IDENTIFIED EXTERNALLY AS" + OracleProbeTail: {911, 40, "", nil},
		}, Reject},
		{"PDB SNAPSHOT stand-in accepts valid", "ALTER PLUGGABLE DATABASE SNAPSHOT EVERY 1 HOURS", map[string]*answer{
			"ALTER PLUGGABLE DATABASE SNAPSHOT EVERY 1 HOURS":     {12754, 0, "", nil},
			"ALTER PLUGGABLE DATABASE REFRESH MODE EVERY 1 HOURS": {1031, 0, "", nil},
		}, Accept},
		{"NUMBER COLLATE stand-in accepts valid", "ALTER TABLE a ADD(a NUMBER COLLATE a)", map[string]*answer{
			"ALTER TABLE a ADD(a NUMBER COLLATE a)":                         {43911, 27, "", nil},
			"ALTER TABLE a ADD(a VARCHAR2(10) COLLATE a)":                   {942, 0, "", nil},
			"ALTER TABLE t ADD(a VARCHAR2(10) COLLATE a)":                   {1430, 18, "", nil},
			"ALTER TABLE t ADD(a VARCHAR2(10) COLLATE a)" + OracleProbeTail: {911, 45, "", nil},
		}, Accept},
		{"duplicate column ORA-00957 stand-in accepts valid", "ALTER TABLE a ADD(a NUMBER, a NUMBER)", map[string]*answer{
			"ALTER TABLE a ADD(a NUMBER, a NUMBER)":                   {957, 28, "", nil},
			"ALTER TABLE a ADD(a NUMBER, b NUMBER)":                   {942, 0, "", nil},
			"ALTER TABLE t ADD(a NUMBER, b NUMBER)":                   {1430, 18, "", nil},
			"ALTER TABLE t ADD(a NUMBER, b NUMBER)" + OracleProbeTail: {911, 39, "", nil},
		}, Accept},
		{"duplicate column ORA-00957 stand-in rejects unclosed paren", "ALTER TABLE a ADD(a NUMBER, a", map[string]*answer{
			"ALTER TABLE a ADD(a NUMBER, a":                   {957, 28, "", nil},
			"ALTER TABLE a ADD(a NUMBER, b":                   {907, 29, "", nil},
			"ALTER TABLE a ADD(a NUMBER, b" + OracleProbeTail: {911, 31, "", nil},
		}, Reject},
		{"INMEMORY JOIN GROUP duplicate column ORA-00957 stand-in rejects missing comma", "CREATE INMEMORY JOIN GROUP a(a(a) a(a))", map[string]*answer{
			"CREATE INMEMORY JOIN GROUP a(a(a) a(a))":                   {903, 29, "", nil},
			"CREATE INMEMORY JOIN GROUP a(t(a) a(a))":                   {903, 34, "", nil},
			"CREATE INMEMORY JOIN GROUP a(t(a) t(a))":                   {957, 34, "", nil},
			"CREATE INMEMORY JOIN GROUP a(t(a) t(b))":                   {907, 34, "", nil},
			"CREATE INMEMORY JOIN GROUP a(t(a) t(b))" + OracleProbeTail: {911, 41, "", nil},
		}, Reject},
		{"SPLIT PARTITION NESTED TABLE stand-in accepts valid", "ALTER TABLE a SPLIT PARTITION a AT(1) NESTED TABLE a INTO(PARTITION a, PARTITION a)", map[string]*answer{
			"ALTER TABLE a SPLIT PARTITION a AT(1) NESTED TABLE a INTO(PARTITION a, PARTITION a)":     {22912, 51, "", nil},
			"ALTER TABLE p SPLIT PARTITION p1 AT(1) INTO(PARTITION a, PARTITION a)":                   {14011, 67, "", nil},
			"ALTER TABLE p SPLIT PARTITION p1 AT(1) INTO(PARTITION a, PARTITION a)" + OracleProbeTail: {911, 71, "", nil},
		}, Accept},
		{"MV LOG SPLIT crash ORA-00600 stand-in accepts PARALLEL", "ALTER MATERIALIZED VIEW LOG ON a SPLIT PARTITION a AT(1) PARALLEL", map[string]*answer{
			"ALTER MATERIALIZED VIEW LOG ON a SPLIT PARTITION a AT(1) PARALLEL": {rawErr: internal},
			"ALTER TABLE p SPLIT PARTITION p1 AT(1) PARALLEL":                   parsed,
		}, Accept},
		{"MV LOG SPLIT crash ORA-00600 stand-in rejects trailing zzz", "ALTER MATERIALIZED VIEW LOG ON a SPLIT PARTITION a AT(1) PARALLEL zzz", map[string]*answer{
			"ALTER MATERIALIZED VIEW LOG ON a SPLIT PARTITION a AT(1) PARALLEL zzz": {rawErr: internal},
			"ALTER TABLE p SPLIT PARTITION p1 AT(1) PARALLEL zzz":                   {14126, 48, "", nil},
			"ALTER TABLE p SPLIT PARTITION p1 AT(1) PARALLEL zzz" + OracleProbeTail: {14126, 48, "", nil},
		}, Reject},
		{"DIMENSION self-loop stand-in accepts valid", "ALTER DIMENSION a ADD HIERARCHY a(a CHILD OF a JOIN KEY a REFERENCES a)", map[string]*answer{
			"ALTER DIMENSION a ADD HIERARCHY a(a CHILD OF a JOIN KEY a REFERENCES a)":                   {30336, 70, "", nil},
			"ALTER DIMENSION d ADD HIERARCHY a(a CHILD OF b JOIN KEY a REFERENCES b)":                   {30342, 35, "", nil},
			"ALTER DIMENSION d ADD HIERARCHY a(a CHILD OF b JOIN KEY a REFERENCES b)" + OracleProbeTail: {911, 74, "", nil},
		}, Accept},
		{"USING INDEX a stand-in accepts valid", "ALTER TABLE a ENABLE PRIMARY KEY USING INDEX a", map[string]*answer{
			"ALTER TABLE a ENABLE PRIMARY KEY USING INDEX a":                      {3049, 45, "", nil},
			"ALTER TABLE a ENABLE PRIMARY KEY USING INDEX idx1":                   {942, 0, "", nil},
			"ALTER TABLE t ENABLE PRIMARY KEY USING INDEX idx1":                   {2432, 0, "", nil},
			"ALTER TABLE t ENABLE PRIMARY KEY USING INDEX idx1" + OracleProbeTail: {911, 51, "", nil},
		}, Accept},
		{"FLASHBACK ARCHIVE a stand-in accepts valid", "ALTER TABLE a FLASHBACK ARCHIVE a", map[string]*answer{
			"ALTER TABLE a FLASHBACK ARCHIVE a":                     {55643, 32, "", nil},
			"ALTER TABLE a FLASHBACK ARCHIVE fa1":                   {942, 0, "", nil},
			"ALTER TABLE t FLASHBACK ARCHIVE fa1":                   {55605, 0, "", nil},
			"ALTER TABLE t FLASHBACK ARCHIVE fa1" + OracleProbeTail: {911, 37, "", nil},
		}, Accept},
		{"CREATE MANDATORY PROFILE stand-in accepts valid", "CREATE MANDATORY PROFILE a LIMIT SESSIONS_PER_USER 1", map[string]*answer{
			"CREATE MANDATORY PROFILE a LIMIT SESSIONS_PER_USER 1": {2387, 51, "", nil},
			"CREATE PROFILE a LIMIT SESSIONS_PER_USER 1":           {1031, 0, "", nil},
		}, Accept},
		{"CREATE SPFILE = 's' FROM PFILE stand-in accepts valid", "CREATE SPFILE = 's' FROM PFILE", map[string]*answer{
			"CREATE SPFILE = 's' FROM PFILE": {65136, 16, "", nil},
			"CREATE PFILE = 's' FROM SPFILE": {1031, 0, "", nil},
		}, Accept},
		{"CREATE SPFILE FROM PFILE AS COPY stand-in accepts valid", "CREATE SPFILE FROM PFILE AS COPY", map[string]*answer{
			"CREATE SPFILE FROM PFILE AS COPY":       {65168, 25, "", nil},
			"CREATE SPFILE FROM PFILE = 's' AS COPY": {1031, 0, "", nil},
		}, Accept},
		{"CREATE TABLESPACE SET TEMPLATE stand-in accepts valid", "CREATE TABLESPACE SET a USING TEMPLATE(DATAFILE 's')", map[string]*answer{
			"CREATE TABLESPACE SET a USING TEMPLATE(DATAFILE 's')": {2537, 51, "", nil},
			"CREATE TABLESPACE a DATAFILE 's'":                     {1031, 0, "", nil},
		}, Accept},
		{"AUDIT -> NOAUDIT stand-in accepts valid traditional audit", "AUDIT SELECT ON t BY ACCESS WHENEVER NOT SUCCESSFUL", map[string]*answer{
			"AUDIT SELECT ON t BY ACCESS WHENEVER NOT SUCCESSFUL": {46401, 6, "", nil},
			"NOAUDIT SELECT ON t WHENEVER NOT SUCCESSFUL":         parsed,
		}, Accept},
		{"AUDIT -> NOAUDIT stand-in rejects invalid option", "AUDIT PORT a", map[string]*answer{
			"AUDIT PORT a":                     {46401, 6, "", nil},
			"NOAUDIT PORT a":                   {956, 8, "", nil},
			"NOAUDIT PORT a" + OracleProbeTail: {956, 8, "", nil},
		}, Reject},
		{"AUDIT -> NOAUDIT stand-in preserves duplicate BY mode clauses", "AUDIT SELECT ON t BY ACCESS BY SESSION", map[string]*answer{
			"AUDIT SELECT ON t BY ACCESS BY SESSION":          {46401, 6, "", nil},
			"NOAUDIT SELECT ON t BY ACCESS":                   {1718, 20, "", nil},
			"NOAUDIT SELECT ON t BY ACCESS" + OracleProbeTail: {1718, 20, "", nil},
		}, Reject},
		{"RECOVER MANAGED STANDBY FINISH FORCE stand-in accepts valid", "ALTER DATABASE RECOVER MANAGED STANDBY DATABASE FINISH FORCE", map[string]*answer{
			"ALTER DATABASE RECOVER MANAGED STANDBY DATABASE FINISH FORCE":     {274, 0, "", nil},
			"ALTER DATABASE RECOVER MANAGED STANDBY DATABASE CANCEL IMMEDIATE": {1031, 0, "", nil},
		}, Accept},
		{"RECOVER MANAGED STANDBY FINISH FORCE zzz rejects", "ALTER DATABASE RECOVER MANAGED STANDBY DATABASE FINISH FORCE zzz", map[string]*answer{
			"ALTER DATABASE RECOVER MANAGED STANDBY DATABASE FINISH FORCE zzz":                       {274, 0, "", nil},
			"ALTER DATABASE RECOVER MANAGED STANDBY DATABASE CANCEL IMMEDIATE zzz":                   {16036, 65, "", nil},
			"ALTER DATABASE RECOVER MANAGED STANDBY DATABASE CANCEL IMMEDIATE zzz" + OracleProbeTail: {16036, 65, "", nil},
		}, Reject},
		{"CREATE DATABASE accepted on ORA-01501", "CREATE DATABASE a", map[string]*answer{
			"CREATE DATABASE a": {1501, 0, "", nil},
		}, Accept},
		{"CREATE JAVA SOURCE AS accepted on ORA-29538", "CREATE JAVA SOURCE AS 's' zzz", map[string]*answer{
			"CREATE JAVA SOURCE AS 's' zzz": {29538, 0, "", nil},
		}, Accept},
		{"CREATE MLE MODULE accepted on ORA-00439", "CREATE MLE MODULE a LANGUAGE JAVASCRIPT AS a", map[string]*answer{
			"CREATE MLE MODULE a LANGUAGE JAVASCRIPT AS a": {439, 0, "", nil},
		}, Accept},
		{"CREATE MLE MODULE empty body rejected on ORA-04125", "CREATE MLE MODULE a LANGUAGE JAVASCRIPT AS", map[string]*answer{
			"CREATE MLE MODULE a LANGUAGE JAVASCRIPT AS":                   {4125, 0, "", nil},
			"CREATE MLE MODULE a LANGUAGE JAVASCRIPT AS" + OracleProbeTail: {911, 44, "", nil},
		}, Reject},
		{"GRANT program unit syntax error ORA-28700 rejected", "GRANT TO FUNCTION a", map[string]*answer{
			"GRANT TO FUNCTION a":                   {28700, 0, "", nil},
			"GRANT TO FUNCTION a" + OracleProbeTail: {911, 21, "", nil},
		}, Reject},
		{"GRANT program unit syntax error ORA-28701 rejected", "GRANT a TO FUNCTION a, a", map[string]*answer{
			"GRANT a TO FUNCTION a, a":                   {28701, 0, "", nil},
			"GRANT a TO FUNCTION a, a" + OracleProbeTail: {911, 26, "", nil},
		}, Reject},
		{"GRANT duplicate grantee ORA-01700 stand-in rejects trailing zzz", "GRANT a TO a, a zzz", map[string]*answer{
			"GRANT a TO a, a zzz":                   {1700, 14, "", nil},
			"GRANT a TO a, b zzz":                   {3049, 16, "", nil},
			"GRANT a TO a, b zzz" + OracleProbeTail: {911, 21, "", nil},
		}, Reject},
		{"GRANT IDENTIFIED BY ORA-00989 stand-in rejects trailing comma", "GRANT a TO a IDENTIFIED BY a,", map[string]*answer{
			"GRANT a TO a IDENTIFIED BY a,":                      {989, 28, "", nil},
			"GRANT a TO a, b IDENTIFIED BY a,":                   {988, 32, "", nil},
			"GRANT a TO a, b IDENTIFIED BY a," + OracleProbeTail: {911, 34, "", nil},
		}, Reject},
		{"ALTER MATERIALIZED VIEW SUBPARTITIONS ORA-14156 rejected", "ALTER MATERIALIZED VIEW a SET SUBPARTITION TEMPLATE SUBPARTITIONS 1 IN(a)", map[string]*answer{
			"ALTER MATERIALIZED VIEW a SET SUBPARTITION TEMPLATE SUBPARTITIONS 1 IN(a)":                   {14156, 52, "", nil},
			"ALTER MATERIALIZED VIEW a SET SUBPARTITION TEMPLATE SUBPARTITIONS 1 IN(a)" + OracleProbeTail: {911, 75, "", nil},
		}, Reject},
		{"LOCK TABLE a PARTITION(1) stand-in to p rejects ORA-14108", "LOCK TABLE a PARTITION(1) IN SHARE MODE", map[string]*answer{
			"LOCK TABLE a PARTITION(1) IN SHARE MODE":                   {942, 11, "", nil},
			"LOCK TABLE p PARTITION(1) IN SHARE MODE":                   {14108, 23, "", nil},
			"LOCK TABLE p PARTITION(1) IN SHARE MODE" + OracleProbeTail: {911, 41, "", nil},
		}, Reject},
		{"CREATE SCHEMA AUTHORIZATION a stand-in rejects trailing zzz", "CREATE SCHEMA AUTHORIZATION a CREATE TABLE a (a NUMBER) zzz", map[string]*answer{
			"CREATE SCHEMA AUTHORIZATION a CREATE TABLE a (a NUMBER) zzz":                            {2421, 0, "", nil},
			"CREATE SCHEMA AUTHORIZATION OMNI_JUDGE CREATE TABLE a (a NUMBER) zzz":                   {2425, 39, "", nil},
			"CREATE SCHEMA AUTHORIZATION OMNI_JUDGE CREATE TABLE a (a NUMBER) zzz" + OracleProbeTail: {911, 70, "", nil},
		}, Reject},
		{"CREATE FUNCTION a.zzz stand-in to OMNI_JUDGE rejects PLS-00103", "CREATE FUNCTION a.zzz a RETURN NUMBER IS BEGIN RETURN 1; END;", map[string]*answer{
			"CREATE FUNCTION a.zzz a RETURN NUMBER IS BEGIN RETURN 1; END;":          {1031, 0, "", nil},
			"CREATE FUNCTION OMNI_JUDGE.zzz a RETURN NUMBER IS BEGIN RETURN 1; END;": {24344, 0, "PLS-00103", nil},
		}, Reject},
		{"CREATE AUDIT POLICY a ONLY rejected on ORA-46373 without rewrite", "CREATE AUDIT POLICY a ONLY", map[string]*answer{
			"CREATE AUDIT POLICY a ONLY":                   {46373, 0, "", nil},
			"CREATE AUDIT POLICY a ONLY" + OracleProbeTail: {911, 28, "", nil},
		}, Reject},
		{"DROP CONTEXT rejected at EOF (1804 removed from checkedAtEnd)", "DROP CONTEXT", map[string]*answer{
			"DROP CONTEXT": {1804, 12, "", nil},
		}, Reject},
		{"SELECT * FOR UPDATE WAIT rejected at EOF (30020 removed from checkedAtEnd)", "SELECT * FOR UPDATE WAIT", map[string]*answer{
			"SELECT * FOR UPDATE WAIT": {30020, 24, "", nil},
		}, Reject},
		{"FOR UPDATE over OFFSET rejected on ORA-02014", "(SELECT dummy FROM dual OFFSET 1 ROWS) FOR UPDATE", map[string]*answer{
			"(SELECT dummy FROM dual OFFSET 1 ROWS) FOR UPDATE":                   {2014, 0, "", nil},
			"(SELECT dummy FROM dual OFFSET 1 ROWS) FOR UPDATE" + OracleProbeTail: {911, 51, "", nil},
		}, Reject},
		{"FOR UPDATE over WITH TIES rejected on ORA-02014", "SELECT dummy FROM dual FETCH FIRST 1 ROWS WITH TIES FOR UPDATE", map[string]*answer{
			"SELECT dummy FROM dual FETCH FIRST 1 ROWS WITH TIES FOR UPDATE":                   {2014, 0, "", nil},
			"SELECT dummy FROM dual FETCH FIRST 1 ROWS WITH TIES FOR UPDATE" + OracleProbeTail: {911, 64, "", nil},
		}, Reject},
		{"semantic FOR UPDATE over DISTINCT view stays accepted on ORA-02014", "SELECT * FROM (SELECT DISTINCT dummy FROM dual) FOR UPDATE", map[string]*answer{
			"SELECT * FROM (SELECT DISTINCT dummy FROM dual) FOR UPDATE":                   {2014, 0, "", nil},
			"SELECT * FROM (SELECT DISTINCT dummy FROM dual) FOR UPDATE" + OracleProbeTail: {911, 60, "", nil},
		}, Accept},
		{"ALTER DISKGROUP trailing comma rejected on ORA-15260", "ALTER DISKGROUP a REBALANCE WITH RESTORE,", map[string]*answer{
			"ALTER DISKGROUP a REBALANCE WITH RESTORE,": {15260, 0, "", nil},
		}, Reject},
		{"PL/SQL structural error PLS-00710 rejected", "BEGIN PRAGMA AUTONOMOUS_TRANSACTION; NULL; END;", map[string]*answer{
			"BEGIN PRAGMA AUTONOMOUS_TRANSACTION; NULL; END;": {6550, 13, "PLS-00710", nil},
		}, Reject},
		{"PL/SQL PLS-00306 with inquiry directive named arg rejected", "DECLARE PROCEDURE p(a NUMBER) IS BEGIN NULL; END; BEGIN p($$PLSQL_LINE => 1); END;", map[string]*answer{
			"DECLARE PROCEDURE p(a NUMBER) IS BEGIN NULL; END; BEGIN p($$PLSQL_LINE => 1); END;": {6550, 56, "PLS-00306", nil},
		}, Reject},
		{"PL/SQL ordinary PLS-00306 signature mismatch stays accepted", "DECLARE PROCEDURE p(x NUMBER) IS BEGIN NULL; END; BEGIN p; END;", map[string]*answer{
			"DECLARE PROCEDURE p(x NUMBER) IS BEGIN NULL; END; BEGIN p; END;": {6550, 56, "PLS-00306", nil},
		}, Accept},
		{"PL/SQL REF omni_missing_t stand-in rejects PLS-00550", "DECLARE x REF omni_missing_t CHARACTER SET NCHAR_CS; BEGIN NULL; END;", map[string]*answer{
			"DECLARE x REF omni_missing_t CHARACTER SET NCHAR_CS; BEGIN NULL; END;": {6550, 14, "PLS-00201", nil},
			"DECLARE x REF y CHARACTER SET NCHAR_CS; BEGIN NULL; END;":              {6550, 10, "PLS-00550", nil},
		}, Reject},
		{"PL/SQL ORA-32551 rejected", "DECLARE x NUMBER; BEGIN SELECT 1 INTO $$PLSQL_LINE FROM dual; END;", map[string]*answer{
			"DECLARE x NUMBER; BEGIN SELECT 1 INTO $$PLSQL_LINE FROM dual; END;": {6550, 39, "ORA-32551", nil},
		}, Reject},
		{"unclosed paren caught when Oracle aborts early on ORA-01031", "ADMINISTER KEY MANAGEMENT EXPORT KEYS WITH SECRET a TO 's' IDENTIFIED BY a WITH IDENTIFIER IN(a", map[string]*answer{
			"ADMINISTER KEY MANAGEMENT EXPORT KEYS WITH SECRET a TO 's' IDENTIFIED BY a WITH IDENTIFIER IN(a": {1031, 0, "", nil},
		}, Reject},
		{"unclosed paren caught on CREATE INMEMORY JOIN GROUP parser bug", "CREATE INMEMORY JOIN GROUP a(a(a)", map[string]*answer{
			"CREATE INMEMORY JOIN GROUP a(a(a)": {903, 29, "", nil},
			"CREATE INMEMORY JOIN GROUP a(t(a)": parsed,
		}, Reject},
		{"successful PL/SQL parse with skipped $IF branch preserves Accept", "BEGIN $IF FALSE $THEN ( $ELSE NULL; $END END;", map[string]*answer{
			"BEGIN $IF FALSE $THEN ( $ELSE NULL; $END END;": parsed,
		}, Accept},
		{"PL/SQL $IF branch with semantic PLS-00201 preserves Accept", "BEGIN $IF FALSE $THEN ( $ELSE zz; $END END;", map[string]*answer{
			"BEGIN $IF FALSE $THEN ( $ELSE zz; $END END;": {6550, 30, "PLS-00201", nil},
		}, Accept},
		{"stray rparen still accepted when Oracle accepts", "CREATE JAVA SOURCE USING BFILE a, 's')", map[string]*answer{
			"CREATE JAVA SOURCE USING BFILE a, 's')": parsed,
		}, Accept},
	} {
		t.Run(c.name, func(t *testing.T) {
			parse := func(text string) error {
				a, ok := c.answers[text]
				if !ok {
					t.Fatalf("%s: unexpected parse(%q)", c.name, text)
				}
				if a == nil {
					return nil
				}
				if a.rawErr != nil {
					return a.rawErr
				}
				return oracleRaised(a.code, a.position, a.compile)
			}
			got, err := OracleVerdictOnly(c.statement, parse)
			if err != nil {
				t.Fatalf("%s: err=%v", c.name, err)
			}
			if got != c.want {
				t.Errorf("%s: got %q, want %q", c.name, got, c.want)
			}
		})
	}
}

func TestHasUnclosedParen(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		want bool
	}{
		{"SELECT (1) FROM dual", false},
		{"SELECT (1 FROM dual", true},
		{"SELECT '(' FROM dual", false},
		{"SELECT q'[(]' FROM dual", false},
		{"SELECT Q'{(}' FROM dual", false},
		{"SELECT q'((a)' FROM dual", false},
		{`SELECT "(" FROM dual`, false},
		{"SELECT 1 -- (\nFROM dual", false},
		{"SELECT /* ( */ 1 FROM dual", false},
		{"CREATE JAVA SOURCE AS class X { void f( {", false},
		{"CREATE MLE MODULE a LANGUAGE JAVASCRIPT AS function f( {", false},
		{"CREATE JAVA SOURCE USING BFILE a, 's')", false}, // stray ')' is not unclosed '('
	} {
		if got := hasUnclosedParen(tc.sql); got != tc.want {
			t.Errorf("hasUnclosedParen(%q) = %v, want %v", tc.sql, got, tc.want)
		}
	}
}

func TestOracleVerdictOnlySchemaCustomScratchUser(t *testing.T) {
	type answer struct {
		code, position int
		compile        string
	}
	for _, c := range []struct {
		name      string
		statement string
		answers   map[string]*answer
		want      string
	}{
		{"PL/SQL schema stand-in uses connected scratch user", "CREATE FUNCTION a.zzz a RETURN NUMBER IS BEGIN RETURN 1; END;", map[string]*answer{
			"CREATE FUNCTION a.zzz a RETURN NUMBER IS BEGIN RETURN 1; END;":            {1031, 0, ""},
			"CREATE FUNCTION SCRATCH_USER.zzz a RETURN NUMBER IS BEGIN RETURN 1; END;": {24344, 0, "PLS-00103"},
		}, Reject},
		{"CREATE SCHEMA stand-in uses connected scratch user", "CREATE SCHEMA AUTHORIZATION a GRANT SELECT ON a TO a zzz", map[string]*answer{
			"CREATE SCHEMA AUTHORIZATION a GRANT SELECT ON a TO a zzz":                                         {2421, 0, ""},
			"CREATE SCHEMA AUTHORIZATION SCRATCH_USER GRANT SELECT ON t TO SCRATCH_USER zzz":                   {2426, 41, ""},
			"CREATE SCHEMA AUTHORIZATION SCRATCH_USER GRANT SELECT ON t TO SCRATCH_USER zzz" + OracleProbeTail: {911, 79, ""},
		}, Reject},
	} {
		t.Run(c.name, func(t *testing.T) {
			parse := func(text string) error {
				a, ok := c.answers[text]
				if !ok {
					t.Fatalf("%s: unexpected parse(%q)", c.name, text)
				}
				if a == nil {
					return nil
				}
				return oracleRaised(a.code, a.position, a.compile)
			}
			got, err := OracleVerdictOnlySchema(c.statement, "SCRATCH_USER", parse)
			if err != nil {
				t.Fatalf("%s: err=%v", c.name, err)
			}
			if got != c.want {
				t.Errorf("%s: got %q, want %q", c.name, got, c.want)
			}
		})
	}
}
