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
		{"a block with a SQL syntax error", "BEGIN INSERT INTO a VALUES (1) a a; END;", &answer{6550, 31, "ORA-03049"}, nil, "reject@31"},
		{"a type with a PL/SQL syntax error", "CREATE TYPE a AS OBJECT (a INT", &answer{24344, 0, "PLS-00103"}, nil, Reject},
		{"a block with a PL/SQL syntax error", "BEGIN NULL END;", &answer{6550, 11, "PLS-00103"}, nil, "reject@11"},
		{"a block with a name it cannot resolve", "BEGIN zz; END;", &answer{6550, 6, "PLS-00201"}, nil, Accept},
		{"a probe answered at the tail by another error", "ALTER MATERIALIZED VIEW LOG ON a MERGE PARTITIONS a, a", &answer{12002, 0, ""}, &answer{14006, 56, ""}, Accept},
		{"a probe the tail stops early, after a lookup", "ALTER INDEX a UNUSABLE", &answer{1418, 0, ""}, &answer{14048, 14, ""}, Accept},
		{"a probe the tail stops early, after a syntax error", "SELECT 1 FROM dual a b", &answer{933, 21, ""}, &answer{923, 7, ""}, "reject@21"},
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
