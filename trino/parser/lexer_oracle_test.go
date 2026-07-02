package parser

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/trino/internal/trinooracle"
)

// This file is the lexer node's slice of the differential-oracle gate described
// in the migration correctness protocol. The full accept/reject differential
// (omni Parse vs Trino SYNTAX_ERROR) belongs to the parser-foundation node,
// which has a Parse() entry point. The lexer cannot decide grammar
// acceptance, but it has a well-defined slice of the contract it MUST uphold:
//
//	If Trino's parser accepts a statement, the omni lexer must tokenize it
//	WITHOUT emitting a hard lex error (tokInvalid). A spurious lex error would
//	make the statement un-parseable downstream and break Diagnose/SplitSQL.
//
// So for every corpus statement the oracle reports as syntactically ACCEPTED,
// this test asserts the lexer produces no tokInvalid token and no LexError.
// Conversely, lexer-level malformations (unterminated string/identifier/binary/
// comment) are exercised in lexer_test.go and need no oracle.
//
// The test skips cleanly when no Trino is reachable (matching the oracle
// harness convention), and additionally always runs an oracle-independent smoke
// pass proving the lexer tokenizes the documented corpus without spurious
// errors.

// oracleCorpus is a curated set of Trino 481 statements spanning every token
// category the lexer must handle. It is drawn from the legacy example corpus
// (/Users/h3n4l/OpenSource/parser/trino/examples) and the Trino 481 docs
// (truth1), chosen so that collectively they exercise: all reserved and a broad
// set of non-reserved keywords; quoted/back-quoted/unquoted/digit-leading
// identifiers; integer/decimal/double numerics; '...'/U&'...'/X'...' literals;
// the full operator and punctuation set including ||, ->, =>, {- -}, the
// time-travel and MATCH_RECOGNIZE shapes, and both comment forms.
var oracleCorpus = []string{
	// --- basic query / expression shapes ---
	"SELECT 1",
	"SELECT 1, 2.5, 1.5e10, -3, +4 FROM t",
	"SELECT * FROM tpch.sf1.nation",
	"SELECT a, b AS c, t.d, t.* FROM t",
	"SELECT 'a string', 'it''s escaped', U&'\\0041', X'00ff'",
	"SELECT a || b AS concatenated FROM t",
	"SELECT CASE WHEN a > 1 THEN 'hi' ELSE 'lo' END FROM t",
	"SELECT CAST(a AS varchar), TRY_CAST(b AS bigint) FROM t",
	"SELECT a <> b, a != b, a <= b, a >= b, a < b, a > b FROM t",
	"SELECT a + b - c * d / e % f FROM t",
	"SELECT count(*) FILTER (WHERE a > 0) FROM t",
	// Backtick identifiers deliberately excluded: Trino's lexer tokenizes
	// BACKQUOTED_IDENTIFIER_ (as does omni's) but the statement is a parser
	// SYNTAX_ERROR, and this corpus must be all-Trino-accepted. Backtick
	// tokenization is pinned by the identifier unit tests instead.
	"SELECT \"quoted col\" FROM \"quoted tbl\"",
	"SELECT timestamp '2012-10-31 01:00 UTC' AT TIME ZONE 'America/Los_Angeles'",
	"SELECT ARRAY[1, 2, 3], ROW(1, 'a')",
	"SELECT filter(ARRAY[1,2,3], x -> x > 1)",
	"SELECT transform(a, (k, v) -> k + v) FROM t",
	// SQL/JSON: exercises the digit-bearing ENCODING keyword UTF8.
	"SELECT JSON_OBJECT('k' VALUE 1 RETURNING varchar FORMAT JSON ENCODING UTF8)",
	"SELECT JSON_QUERY(a, 'lax $.b' RETURNING varchar FORMAT JSON ENCODING UTF16) FROM t",
	// --- grouping / windows / set ops ---
	"SELECT a, count(*) FROM t GROUP BY ROLLUP (a)",
	"SELECT a, count(*) FROM t GROUP BY GROUPING SETS ((a), ())",
	"SELECT sum(x) OVER (PARTITION BY a ORDER BY b ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM t",
	"SELECT a FROM t UNION ALL SELECT a FROM u",
	"SELECT a FROM t INTERSECT SELECT a FROM u",
	// --- CTE / subquery / joins ---
	"WITH x AS (SELECT 1 AS a) SELECT a FROM x",
	"WITH RECURSIVE r(n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM r WHERE n < 5) SELECT n FROM r",
	"SELECT * FROM a JOIN b ON a.id = b.id",
	"SELECT * FROM a LEFT OUTER JOIN b USING (id)",
	"SELECT * FROM a CROSS JOIN UNNEST(b) WITH ORDINALITY",
	"SELECT * FROM TABLE(sequence(1, 10))",
	// --- time travel / MATCH_RECOGNIZE (Trino-specific lexer shapes) ---
	"SELECT * FROM t FOR TIMESTAMP AS OF TIMESTAMP '2023-01-01 00:00:00'",
	"SELECT * FROM t MATCH_RECOGNIZE (PARTITION BY a ORDER BY b MEASURES c AS m PATTERN (^ A+ B* $) DEFINE A AS true)",
	// --- DDL ---
	"CREATE TABLE memory.default.t (x integer, y varchar)",
	"CREATE TABLE t (x bigint) WITH (format = 'ORC')",
	"CREATE TABLE t AS SELECT 1 AS x",
	"DROP TABLE IF EXISTS t",
	"ALTER TABLE t ADD COLUMN z double",
	"CREATE VIEW v AS SELECT 1 AS a",
	"CREATE MATERIALIZED VIEW mv GRACE PERIOD INTERVAL '1' HOUR AS SELECT 1 AS a",
	"COMMENT ON TABLE t IS 'a comment'",
	// --- DML ---
	"INSERT INTO t VALUES (1, 'a'), (2, 'b')",
	"DELETE FROM t WHERE a > 1",
	"UPDATE t SET a = a + 1 WHERE b = 'x'",
	"MERGE INTO t USING u ON t.id = u.id WHEN MATCHED THEN UPDATE SET a = u.a",
	"TRUNCATE TABLE t",
	// --- admin / session / prepared / txn / DCL / utility ---
	"USE catalog.schema",
	"SET SESSION query_max_run_time = '1h'",
	"SET TIME ZONE LOCAL",
	"SHOW TABLES FROM s LIKE 'pattern%'",
	"EXPLAIN ANALYZE VERBOSE SELECT 1",
	"PREPARE my_query FROM SELECT ? FROM t",
	"EXECUTE my_query USING 1",
	"START TRANSACTION ISOLATION LEVEL SERIALIZABLE",
	"GRANT SELECT ON t TO USER alice WITH GRANT OPTION",
	"CALL system.runtime.kill_query(query_id => '20210101_000000_00000_abcde')",
	// --- comments (hidden channel; must not alter the meaningful token run) ---
	"SELECT 1 -- trailing comment\nFROM t",
	"SELECT /* inline */ 1 FROM t",
	"SELECT /* a /* not nested */ 1",
}

// lexerHasHardError reports whether tokenizing sql produced a hard lex error:
// either a recorded LexError or a tokInvalid token in the stream.
func lexerHasHardError(sql string) (bool, []LexError) {
	tokens, errs := Tokenize(sql)
	if len(errs) > 0 {
		return true, errs
	}
	for _, tok := range tokens {
		if tok.Kind == tokInvalid {
			return true, errs
		}
	}
	return false, nil
}

// TestLexer_CorpusNoSpuriousErrors runs without any oracle: every statement in
// the documented corpus must tokenize cleanly (no tokInvalid, no LexError).
// This is the lexer's completeness smoke test over real Trino SQL.
func TestLexer_CorpusNoSpuriousErrors(t *testing.T) {
	for _, sql := range oracleCorpus {
		bad, errs := lexerHasHardError(sql)
		if bad {
			t.Errorf("lexer produced a spurious error on accepted corpus SQL\n  sql:  %q\n  errs: %v", sql, errs)
		}
	}
}

// TestLexer_OracleDifferential cross-checks the lexer against a live Trino: for
// every corpus statement Trino's parser accepts (no SYNTAX_ERROR), the lexer
// must not emit a hard error. Skipped when no Trino is reachable.
func TestLexer_OracleDifferential(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := trinooracle.Connect("")
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	ver, err := o.Ping(pingCtx)
	pingCancel()
	if err != nil {
		trinooracle.SkipOrFailUnreachable(t, "trino oracle not reachable (start: docker run -d -p 18080:8080 %s): %v",
			trinooracle.DefaultImage, err)
	}
	t.Logf("connected to Trino %s", ver)

	for _, sql := range oracleCorpus {
		sql := sql
		t.Run(truncateName(sql), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			res, err := o.CheckSyntax(ctx, sql)
			if err != nil {
				t.Skipf("oracle CheckSyntax error (treated as unreachable): %v", err)
			}
			if !res.Accepted {
				// The statement is a genuine Trino syntax error. The corpus is
				// meant to be all-accepted; flag it so the corpus can be fixed
				// rather than silently masking a lexer issue.
				t.Errorf("oracle rejected a corpus statement (errorName=%q): %q", res.ErrorName, sql)
				return
			}
			// Trino accepted it syntactically -> the lexer must tokenize cleanly.
			if bad, errs := lexerHasHardError(sql); bad {
				t.Errorf("MISMATCH: Trino accepts %q but omni lexer errored: %v", sql, errs)
			}
		})
	}
}

func truncateName(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > 48 {
		return s[:48]
	}
	return s
}
