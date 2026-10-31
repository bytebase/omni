package parser

import "testing"

// ============================================================================
// Temporal keyword-argument conformance oracle.
//
// TIMESTAMPDIFF, TIMESTAMPADD and GET_FORMAT take a bare keyword in their
// first argument position. The engine lexes it as a keyword token, never an
// identifier, so a quoted unit is a syntax error; and the unit set is narrower
// than the interval-unit set — only the nine simple units are accepted, with
// the SQL_TSI_ synonyms folding to their canonical names.
//
// The reject arm is what makes this durable. Without the dedicated production
// the generic argument path parses the unit as a column reference, which
// accepts every statement here and produces a semantically wrong tree: a
// consumer resolving that reference looks for a column named SECOND.
// ============================================================================

var temporalKeywordArgCorpus = []string{
	// --- accept: TIMESTAMPDIFF / TIMESTAMPADD, the nine simple units ---
	"SELECT TIMESTAMPDIFF(MICROSECOND, a, b) FROM tu",
	"SELECT TIMESTAMPDIFF(SECOND, a, b) FROM tu",
	"SELECT TIMESTAMPDIFF(MINUTE, a, b) FROM tu",
	"SELECT TIMESTAMPDIFF(HOUR, a, b) FROM tu",
	"SELECT TIMESTAMPDIFF(DAY, a, b) FROM tu",
	"SELECT TIMESTAMPDIFF(WEEK, a, b) FROM tu",
	"SELECT TIMESTAMPDIFF(MONTH, a, b) FROM tu",
	"SELECT TIMESTAMPDIFF(QUARTER, a, b) FROM tu",
	"SELECT TIMESTAMPDIFF(YEAR, a, b) FROM tu",
	"SELECT TIMESTAMPADD(SECOND, 1, a) FROM tu",
	"SELECT TIMESTAMPADD(YEAR, 1, a) FROM tu",
	// --- accept: the SQL_TSI_ synonyms, which fold to the canonical unit ---
	"SELECT TIMESTAMPDIFF(SQL_TSI_SECOND, a, b) FROM tu",
	"SELECT TIMESTAMPADD(SQL_TSI_YEAR, 1, a) FROM tu",
	"SELECT a + INTERVAL 1 SQL_TSI_DAY FROM tu",
	"SELECT a + INTERVAL 1 DAY FROM tu",
	// --- accept: GET_FORMAT's four types, TIMESTAMP being a DATETIME synonym ---
	"SELECT GET_FORMAT(DATE, 'USA') FROM tu",
	"SELECT GET_FORMAT(TIME, 'USA') FROM tu",
	"SELECT GET_FORMAT(DATETIME, 'USA') FROM tu",
	"SELECT GET_FORMAT(TIMESTAMP, 'USA') FROM tu",
	// --- accept: the unit is a keyword, so it is also still a column name ---
	"SELECT second FROM tu2",
	"CREATE TABLE tu3 (second INT, quarter INT)",
	// --- accept: a schema-qualified name is a stored function, not the builtin ---
	"SELECT test.my_diff(1, 2)",

	// --- reject (1064): compound units are not in interval_time_stamp ---
	"SELECT TIMESTAMPDIFF(DAY_HOUR, a, b) FROM tu",
	"SELECT TIMESTAMPDIFF(YEAR_MONTH, a, b) FROM tu",
	"SELECT TIMESTAMPADD(MINUTE_SECOND, 1, a) FROM tu",
	// --- reject (1064): the unit is lexed as a keyword, never an identifier ---
	"SELECT TIMESTAMPDIFF(`SECOND`, a, b) FROM tu",
	"SELECT TIMESTAMPDIFF('SECOND', a, b) FROM tu",
	"SELECT GET_FORMAT(`DATE`, 'USA') FROM tu",
	"SELECT GET_FORMAT('DATE', 'USA') FROM tu",
	// --- reject (1064): SQL_TSI_MICROSECOND is not one of the eight synonyms ---
	"SELECT TIMESTAMPDIFF(SQL_TSI_MICROSECOND, a, b) FROM tu",
	"SELECT a + INTERVAL 1 SQL_TSI_MICROSECOND FROM tu",
	// --- reject (1064): not a unit at all ---
	"SELECT TIMESTAMPDIFF(NOTAUNIT, a, b) FROM tu",
	"SELECT GET_FORMAT(YEAR, 'USA') FROM tu",
	// --- reject (1064): the arity is fixed ---
	"SELECT TIMESTAMPDIFF(SECOND, a) FROM tu",
	"SELECT TIMESTAMPDIFF(SECOND, a, b, a) FROM tu",
	"SELECT GET_FORMAT(DATE) FROM tu",
}

// TestTiDBTemporalKeywordArgOracle asserts omni agrees with the package's
// MySQL 8.0 oracle on every statement above (container test). The syntax is
// MySQL's, so MySQL is the right reference for it.
func TestTiDBTemporalKeywordArgOracle(t *testing.T) {
	o := startParserOracle(t)

	for _, ddl := range []string{
		"CREATE TABLE IF NOT EXISTS tu (a DATETIME, b DATETIME)",
		"CREATE TABLE IF NOT EXISTS tu2 (second INT, quarter INT)",
		"CREATE FUNCTION IF NOT EXISTS my_diff(x INT, y INT) RETURNS INT DETERMINISTIC RETURN x - y",
	} {
		if _, err := o.db.ExecContext(o.ctx, ddl); err != nil {
			t.Fatalf("seed %q: %v", ddl, err)
		}
	}

	for _, sqlStr := range temporalKeywordArgCorpus {
		_, perr := Parse(sqlStr)
		omniAccepts := perr == nil
		engineAccepts := o.checkSyntax(t, sqlStr)
		if omniAccepts != engineAccepts {
			t.Errorf("DIVERGENCE omni=%v mysql=%v for %q", omniAccepts, engineAccepts, sqlStr)
		}
	}
}
