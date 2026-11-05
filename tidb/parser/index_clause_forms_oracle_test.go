package parser

import "testing"

// The clause set CREATE INDEX and DROP INDEX read: the index_type clause with
// its two spellings and its closed list of names, the ALGORITHM and LOCK
// clauses that close either statement, and the IF EXISTS that TiDB's DROP
// INDEX reads and MySQL's does not.
//
// The rejections are the point of this file. The corpus in tidb/parser pins the
// forms pingcap's own parser test writes, and it writes none of these: an
// ALGORITHM with no value, an IF EXISTS in the wrong place, a TYPE whose name is
// no index type, a second ALGORITHM, a trailing comma in the key_part list, and
// the clauses each engine does not have. Every case here was measured against
// pingcap/tidb:v8.5.5 on 2026-11-04, and the expectation is what that server answered, so
// an image upgrade that changes a verdict fails here instead of passing
// silently.
//
// Ref: index_type_clause, index_type_name, opt_index_name_and_type and
// opt_index_lock_and_algorithm in MySQL 8.0 sql/sql_yacc.yy.

// indexClauseCase is a statement and whether the engine's parser accepts it.
type indexClauseCase struct {
	sql    string
	accept bool
}

var indexClauseCases = []indexClauseCase{
	// The index_type clause, in both spellings and in every position the
	// grammar puts it: before ON in CREATE INDEX, in the index_option list
	// after the column list, and in front of a key definition's column list.
	{"CREATE INDEX i ON t (a) USING BTREE", true},
	{"CREATE INDEX i ON t (a) USING HASH", true},
	{"CREATE INDEX i ON t (a) USING RTREE", true},
	{"CREATE INDEX i ON t (a) TYPE BTREE", true},
	{"CREATE INDEX i ON t (a) TYPE HASH", true},
	{"CREATE INDEX i ON t (a) TYPE RTREE", true},
	{"CREATE INDEX i USING BTREE ON t (a)", true},
	{"CREATE INDEX i TYPE BTREE ON t (a)", true},
	{"CREATE UNIQUE INDEX i ON d.t (a, b ASC) TYPE BTREE", true},
	{"CREATE UNIQUE INDEX i ON d.t (a, b ASC) TYPE HASH", true},
	{"CREATE INDEX i ON t (a) TYPE BTREE COMMENT 'c'", true},
	{"CREATE INDEX i ON t (a) COMMENT 'c' TYPE BTREE", true},
	// index_option is a list, so a second index_type is allowed and the last
	// one stands.
	{"CREATE INDEX i ON t (a) USING BTREE TYPE HASH", true},
	{"CREATE INDEX i TYPE BTREE ON t (a) TYPE HASH", true},
	// A key definition's own index_type, which reads before the column list
	// and again among the options.
	{"CREATE TABLE ct11 (a INT, KEY k (a) USING BTREE)", true},
	{"CREATE TABLE ct12 (a INT, KEY k USING BTREE (a))", true},
	{"CREATE TABLE ct13 (a INT, PRIMARY KEY USING BTREE (a))", true},
	{"ALTER TABLE t ADD KEY k (a) USING HASH", true},
	{"CREATE TABLE ct1 (a INT, KEY k (a) TYPE BTREE)", true},
	{"CREATE TABLE ct2 (a INT, KEY k TYPE BTREE (a))", true},
	{"CREATE TABLE ct3 (a INT, PRIMARY KEY pk TYPE BTREE (a))", true},
	{"CREATE TABLE ct4 (a INT, UNIQUE KEY u TYPE BTREE (a))", true},
	{"ALTER TABLE t ADD KEY k (a) TYPE BTREE", true},
	{"ALTER TABLE t ADD KEY k TYPE BTREE (a)", true},
	{"ALTER TABLE t ADD UNIQUE KEY u (a) TYPE HASH", true},
	// In front of the column list the TYPE spelling needs a name before it,
	// where USING does not: opt_index_name_and_type reads `ident TYPE
	// index_type_name`, so a bare TYPE is read as the index's own name and the
	// type after it has nowhere to go.
	{"CREATE TABLE ct5 (a INT, KEY TYPE BTREE (a))", false},
	{"CREATE TABLE ct6 (a INT, PRIMARY KEY TYPE BTREE (a))", false},
	{"CREATE TABLE ct7 (a INT, UNIQUE KEY TYPE BTREE (a))", false},
	{"ALTER TABLE t ADD KEY TYPE BTREE (a)", false},
	{"ALTER TABLE t ADD PRIMARY KEY TYPE BTREE (a)", false},
	// index_type_name is a closed keyword list, so a name that is no index
	// type is a syntax error, and so is a backquoted `BTREE`, which lexes as an
	// identifier rather than as the keyword the grammar names.
	{"CREATE INDEX i ON t (a) USING nosuchtype", false},
	{"CREATE INDEX i ON t (a) TYPE nosuchtype", false},
	{"CREATE INDEX i ON t (a) USING `BTREE`", false},
	{"CREATE INDEX i USING nosuchtype ON t (a)", false},
	{"CREATE INDEX i ON t (a) USING", false},
	{"CREATE INDEX i ON t (a) TYPE", false},
	{"CREATE INDEX i ON t (a) TYPE 1", false},
	{"CREATE INDEX i ON t (a) TYPE 'BTREE'", false},
	{"CREATE TABLE ct8 (a INT, KEY k (a) USING nosuchtype)", false},
	{"CREATE TABLE ct9 (a INT, KEY k USING nosuchtype (a))", false},
	{"ALTER TABLE t ADD KEY k (a) USING nosuchtype", false},
	// HNSW is TiDB's fourth index type, for a vector index, and it is a token
	// there as the other three are: v8.5.5 reads it unquoted, in any case, and
	// answers 1064 for the backquoted spelling. CreateIndexStmt's action answers
	// 1149 for it on any index but a vector one, where a key definition may take
	// it. MySQL and MariaDB have no such name.
	{"CREATE INDEX i ON t (a) USING HNSW", false},
	{"CREATE INDEX i ON t (a) USING hnsw", false},
	{"CREATE INDEX i ON t (a) TYPE HNSW", false},
	{"CREATE INDEX i ON t (a) USING `HNSW`", false},
	{"CREATE TABLE ct10 (a INT, KEY k (a) USING HNSW)", true},
	// The ALGORITHM and LOCK clauses that close CREATE INDEX. The value is the
	// reserved word DEFAULT or an ident, which is why DEFAULT needs an arm of
	// its own, and an ident the engine does not know is a 1800 or a 1801 once
	// the statement runs rather than a 1064 while it parses.
	{"CREATE INDEX i ON t (a) ALGORITHM = DEFAULT", true},
	{"CREATE INDEX i ON t (a) ALGORITHM DEFAULT", true},
	{"CREATE INDEX i ON t (a) LOCK = DEFAULT", true},
	{"CREATE INDEX i ON t (a) LOCK DEFAULT", true},
	{"CREATE INDEX i ON t (a) ALGORITHM = INPLACE", true},
	{"CREATE INDEX i ON t (a) LOCK = NONE", true},
	{"CREATE INDEX i ON t (a) LOCK = EXCLUSIVE", true},
	{"CREATE INDEX i ON t (a) ALGORITHM = ident", false},
	{"CREATE INDEX i ON t (a) ALGORITHM = `default`", false},
	{"CREATE INDEX i ON t (a) ALGORITHM = DEFAULT LOCK = DEFAULT", true},
	{"CREATE INDEX i ON t (a) LOCK = DEFAULT ALGORITHM = DEFAULT", true},
	// Neither value is optional, neither is a number or a string, no reserved
	// word other than DEFAULT is one, and the grammar takes at most one clause
	// of each kind, so a repeat does not win -- it is a syntax error.
	{"CREATE INDEX i ON t (a) ALGORITHM", false},
	{"CREATE INDEX i ON t (a) LOCK", false},
	{"CREATE INDEX i ON t (a) ALGORITHM =", false},
	{"CREATE INDEX i ON t (a) ALGORITHM = 1", false},
	{"CREATE INDEX i ON t (a) ALGORITHM = 'INPLACE'", false},
	{"CREATE INDEX i ON t (a) ALGORITHM = select", false},
	{"CREATE INDEX i ON t (a) ALGORITHM = INPLACE ALGORITHM = COPY", false},
	{"CREATE INDEX i ON t (a) LOCK = SHARED LOCK = NONE", false},
	{"CREATE INDEX i ON t (a) ALGORITHM = INPLACE LOCK = NONE ALGORITHM = COPY", false},
	// DROP INDEX. IF EXISTS is MariaDB's and TiDB's, and MySQL has none; the
	// place it reads is after INDEX and before the name, in that one spelling.
	{"DROP INDEX idx ON t", true},
	{"DROP INDEX IF EXISTS idx ON t", true},
	{"DROP INDEX IF EXISTS idx ON d.t", true},
	{"DROP INDEX IF NOT EXISTS idx ON t", false},
	{"DROP INDEX idx IF EXISTS ON t", false},
	{"DROP INDEX IF EXISTS ON t", false},
	// The same ALGORITHM and LOCK clauses close DROP INDEX on MySQL and TiDB.
	// MariaDB's DROP INDEX reads neither, in any spelling, and answers 1064 to
	// the whole clause rather than only to DEFAULT.
	{"DROP INDEX idx ON t ALGORITHM = DEFAULT", true},
	{"DROP INDEX idx ON t ALGORITHM DEFAULT", true},
	{"DROP INDEX idx ON t LOCK = DEFAULT", true},
	{"DROP INDEX idx ON t LOCK DEFAULT", true},
	{"DROP INDEX idx ON t ALGORITHM = INPLACE", true},
	{"DROP INDEX idx ON t LOCK = SHARED", true},
	{"DROP INDEX idx ON t ALGORITHM = DEFAULT LOCK = DEFAULT", true},
	{"DROP INDEX idx ON t LOCK = DEFAULT ALGORITHM = DEFAULT", true},
	{"DROP INDEX IF EXISTS idx ON t ALGORITHM = DEFAULT", true},
	{"DROP INDEX idx ON t ALGORITHM", false},
	{"DROP INDEX idx ON t LOCK", false},
	{"DROP INDEX idx ON t ALGORITHM =", false},
	{"DROP INDEX idx ON t ALGORITHM = select", false},
	{"DROP INDEX idx ON t ALGORITHM = INPLACE ALGORITHM = COPY", false},
	{"DROP INDEX idx ON t LOCK = SHARED LOCK = NONE", false},
	{"DROP INDEX idx ON t, ALGORITHM = INPLACE", false},
	// A key_part list takes no trailing comma, wherever the clauses above sit.
	{"CREATE INDEX i ON t (a,)", false},
	{"CREATE UNIQUE INDEX i ON d.t (a, b ASC,) TYPE BTREE", false},
	// LOCK and UNLOCK as statements, which the DROP INDEX rejections above lean
	// on: the LOCK that ends a run-on must not be read as a LOCK TABLES of its
	// value. Both read the keyword that says what they lock, and neither takes
	// it as read.
	{"LOCK TABLES t READ", true},
	{"LOCK TABLE t READ", true},
	{"LOCK = t READ", false},
	{"UNLOCK TABLES", true},
	{"UNLOCK TABLE", true},
	{"UNLOCK nosuch", false},
	// LOCK STATS and UNLOCK STATS are TiDB's own statements, which MySQL has
	// not; stats_lock_forms_oracle_test.go holds the rest of their grammar.
	{"LOCK STATS t", true},
	{"UNLOCK STATS t", true},
	{"LOCK STATS", false},
	{"UNLOCK STATS", false},
}

// TestIndexClauseFormsOracle asserts omni agrees with the pinned TiDB on
// indexClauseCases (container test). The verdict comes from PREPARE, which
// parses the statement without running it: 1064 is the parser's rejection, and
// any other server error means the text parsed and failed later.
func TestIndexClauseFormsOracle(t *testing.T) {
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	for i, c := range indexClauseCases {
		t.Run(c.sql, func(t *testing.T) {
			previous := "<start of cases>"
			if i > 0 {
				previous = indexClauseCases[i-1].sql
			}
			accepts, _, err := tc.syntaxVerdict(t, c.sql, previous)
			if err != nil {
				t.Fatalf("no verdict: %v", err)
			}
			if accepts != c.accept {
				t.Fatalf("TiDB %s accepts=%v, want %v for %q", tidbCorpusVersion, accepts, c.accept, c.sql)
			}
			_, perr := Parse(c.sql)
			if omniAccepts := perr == nil; omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, TiDB %s accepts=%v for %q: %v",
					omniAccepts, tidbCorpusVersion, c.accept, c.sql, perr)
			}
		})
	}
}
