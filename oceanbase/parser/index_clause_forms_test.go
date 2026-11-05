package parser

import "testing"

// The clause set CREATE INDEX and DROP INDEX read, which is where OceanBase
// departs from the MySQL grammar this tree forked, in four places:
//
//   - The index_type clause is spelled USING and only USING. MySQL also takes
//     TYPE, in the index_option list and in front of a key definition's column
//     list, and CE 4.3.5.6 answers 1064 to every position.
//   - Its names are BTREE and HASH. MySQL's third name, RTREE, is a 1064 in
//     every position, as is a name that is no index type and a backquoted
//     `BTREE`, which lexes as an identifier rather than as the keyword.
//   - CREATE INDEX reads no ALGORITHM and no LOCK clause, in any spelling,
//     where MySQL closes the statement with either or both.
//   - DROP INDEX ends at the table: no ALGORITHM, no LOCK, and no IF EXISTS,
//     which MariaDB and TiDB do read.
//
// The engine column is what 5.7.25-OceanBase_CE-v4.3.5.6 answered on
// 2026-11-04. TestIndexClauseFormsOceanBaseOracle reads that column back from a
// live server; TestIndexClauseForms holds the parser to the record without one.
//
// Ref: index_using_algorithm in OceanBase sql_parser_mysql_mode.y;
// index_type_clause, index_type_name, opt_index_name_and_type and
// opt_index_lock_and_algorithm in MySQL 8.0 sql/sql_yacc.yy.
type indexClauseCase struct {
	sql string
	// engine is whether OceanBase CE 4.3.5.6 parses the statement.
	engine bool
	// pending marks a statement OceanBase accepts that this parser does not
	// read yet. Implementing it means clearing the flag in the same change.
	pending bool
}

var indexClauseCases = []indexClauseCase{
	// The index_type clause, in both spellings and in every position the
	// grammar puts it: before ON in CREATE INDEX, in the index_option list
	// after the column list, and in front of a key definition's column list.
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) USING BTREE", engine: true},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) USING HASH", engine: true},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) USING RTREE"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) TYPE BTREE"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) TYPE HASH"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) TYPE RTREE"},
	{sql: "CREATE INDEX obidx_i USING BTREE ON obidx_t (a)", engine: true},
	{sql: "CREATE INDEX obidx_i TYPE BTREE ON obidx_t (a)"},
	{sql: "CREATE UNIQUE INDEX obidx_i ON obidx_d.obidx_t (a, b ASC) TYPE BTREE"},
	{sql: "CREATE UNIQUE INDEX obidx_i ON obidx_d.obidx_t (a, b ASC) TYPE HASH"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) TYPE BTREE COMMENT 'c'"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) COMMENT 'c' TYPE BTREE"},
	// index_option is a list, so a second index_type is allowed and the last
	// one stands.
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) USING BTREE TYPE HASH"},
	{sql: "CREATE INDEX obidx_i TYPE BTREE ON obidx_t (a) TYPE HASH"},
	// A key definition's own index_type, which reads before the column list
	// and again among the options.
	{sql: "CREATE TABLE obidx_c11 (a INT, KEY obidx_k (a) USING BTREE)", engine: true},
	{sql: "CREATE TABLE obidx_c12 (a INT, KEY obidx_k USING BTREE (a))", engine: true},
	{sql: "CREATE TABLE obidx_c13 (a INT, PRIMARY KEY USING BTREE (a))", engine: true},
	{sql: "ALTER TABLE obidx_t ADD KEY obidx_k (a) USING HASH", engine: true},
	{sql: "CREATE TABLE obidx_c1 (a INT, KEY obidx_k (a) TYPE BTREE)"},
	{sql: "CREATE TABLE obidx_c2 (a INT, KEY obidx_k TYPE BTREE (a))"},
	{sql: "CREATE TABLE obidx_c3 (a INT, PRIMARY KEY obidx_pk TYPE BTREE (a))"},
	{sql: "CREATE TABLE obidx_c4 (a INT, UNIQUE KEY obidx_u TYPE BTREE (a))"},
	{sql: "ALTER TABLE obidx_t ADD KEY obidx_k (a) TYPE BTREE"},
	{sql: "ALTER TABLE obidx_t ADD KEY obidx_k TYPE BTREE (a)"},
	{sql: "ALTER TABLE obidx_t ADD UNIQUE KEY obidx_u (a) TYPE HASH"},
	// In front of the column list the TYPE spelling needs a name before it,
	// where USING does not: opt_index_name_and_type reads `ident TYPE
	// index_type_name`, so a bare TYPE is read as the index's own name and the
	// type after it has nowhere to go.
	{sql: "CREATE TABLE obidx_c5 (a INT, KEY TYPE BTREE (a))"},
	{sql: "CREATE TABLE obidx_c6 (a INT, PRIMARY KEY TYPE BTREE (a))"},
	{sql: "CREATE TABLE obidx_c7 (a INT, UNIQUE KEY TYPE BTREE (a))"},
	{sql: "ALTER TABLE obidx_t ADD KEY TYPE BTREE (a)"},
	{sql: "ALTER TABLE obidx_t ADD PRIMARY KEY TYPE BTREE (a)"},
	// index_type_name is a closed keyword list, so a name that is no index
	// type is a syntax error, and so is a backquoted `BTREE`, which lexes as an
	// identifier rather than as the keyword the grammar names.
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) USING nosuchtype"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) TYPE nosuchtype"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) USING `BTREE`"},
	{sql: "CREATE INDEX obidx_i USING nosuchtype ON obidx_t (a)"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) USING"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) TYPE"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) TYPE 1"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) TYPE 'BTREE'"},
	{sql: "CREATE TABLE obidx_c8 (a INT, KEY obidx_k (a) USING nosuchtype)"},
	{sql: "CREATE TABLE obidx_c9 (a INT, KEY obidx_k USING nosuchtype (a))"},
	{sql: "ALTER TABLE obidx_t ADD KEY obidx_k (a) USING nosuchtype"},
	// HNSW is TiDB's fourth index type, for a vector index, and it is a token
	// there as the other three are: v8.5.5 reads it unquoted, in any case, and
	// answers 1064 for the backquoted spelling. MySQL and MariaDB have no such
	// name.
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) USING HNSW"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) USING hnsw"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) TYPE HNSW"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) USING `HNSW`"},
	{sql: "CREATE TABLE obidx_c10 (a INT, KEY obidx_k (a) USING HNSW)"},
	// The ALGORITHM and LOCK clauses that close CREATE INDEX. The value is the
	// reserved word DEFAULT or an ident, which is why DEFAULT needs an arm of
	// its own, and an ident the engine does not know is a 1800 or a 1801 once
	// the statement runs rather than a 1064 while it parses.
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM = DEFAULT"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM DEFAULT"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) LOCK = DEFAULT"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) LOCK DEFAULT"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM = INPLACE"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) LOCK = NONE"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) LOCK = EXCLUSIVE"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM = ident"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM = `default`"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM = DEFAULT LOCK = DEFAULT"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) LOCK = DEFAULT ALGORITHM = DEFAULT"},
	// Neither value is optional, neither is a number or a string, no reserved
	// word other than DEFAULT is one, and the grammar takes at most one clause
	// of each kind, so a repeat does not win -- it is a syntax error.
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) LOCK"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM ="},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM = 1"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM = 'INPLACE'"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM = select"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM = INPLACE ALGORITHM = COPY"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) LOCK = SHARED LOCK = NONE"},
	{sql: "CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM = INPLACE LOCK = NONE ALGORITHM = COPY"},
	// DROP INDEX. IF EXISTS is MariaDB's and TiDB's, and MySQL has none; the
	// place it reads is after INDEX and before the name, in that one spelling.
	{sql: "DROP INDEX obidx_ia ON obidx_t", engine: true},
	{sql: "DROP INDEX IF EXISTS obidx_ia ON obidx_t"},
	{sql: "DROP INDEX IF EXISTS obidx_ia ON obidx_d.obidx_t"},
	{sql: "DROP INDEX IF NOT EXISTS obidx_ia ON obidx_t"},
	{sql: "DROP INDEX obidx_ia IF EXISTS ON obidx_t"},
	{sql: "DROP INDEX IF EXISTS ON obidx_t"},
	// The same ALGORITHM and LOCK clauses close DROP INDEX on MySQL and TiDB.
	// MariaDB's DROP INDEX reads neither, in any spelling, and answers 1064 to
	// the whole clause rather than only to DEFAULT.
	{sql: "DROP INDEX obidx_ia ON obidx_t ALGORITHM = DEFAULT"},
	{sql: "DROP INDEX obidx_ia ON obidx_t ALGORITHM DEFAULT"},
	{sql: "DROP INDEX obidx_ia ON obidx_t LOCK = DEFAULT"},
	{sql: "DROP INDEX obidx_ia ON obidx_t LOCK DEFAULT"},
	{sql: "DROP INDEX obidx_ia ON obidx_t ALGORITHM = INPLACE"},
	{sql: "DROP INDEX obidx_ia ON obidx_t LOCK = SHARED"},
	{sql: "DROP INDEX obidx_ia ON obidx_t ALGORITHM = DEFAULT LOCK = DEFAULT"},
	{sql: "DROP INDEX obidx_ia ON obidx_t LOCK = DEFAULT ALGORITHM = DEFAULT"},
	{sql: "DROP INDEX IF EXISTS obidx_ia ON obidx_t ALGORITHM = DEFAULT"},
	{sql: "DROP INDEX obidx_ia ON obidx_t ALGORITHM"},
	{sql: "DROP INDEX obidx_ia ON obidx_t LOCK"},
	{sql: "DROP INDEX obidx_ia ON obidx_t ALGORITHM ="},
	{sql: "DROP INDEX obidx_ia ON obidx_t ALGORITHM = select"},
	{sql: "DROP INDEX obidx_ia ON obidx_t ALGORITHM = INPLACE ALGORITHM = COPY"},
	{sql: "DROP INDEX obidx_ia ON obidx_t LOCK = SHARED LOCK = NONE"},
	{sql: "DROP INDEX obidx_ia ON obidx_t, ALGORITHM = INPLACE"},
	// A key_part list takes no trailing comma, wherever the clauses above sit.
	{sql: "CREATE INDEX obidx_i ON obidx_t (a,)"},
	{sql: "CREATE UNIQUE INDEX obidx_i ON obidx_d.obidx_t (a, b ASC,) TYPE BTREE"},
	// LOCK as a statement, which the DROP INDEX rejections above lean on: the
	// LOCK that ends a run-on must not be read as a LOCK TABLES of its value.
	{sql: "LOCK TABLES obidx_t READ", engine: true},
	{sql: "LOCK TABLE obidx_t READ", engine: true},
	{sql: "LOCK = obidx_t READ"},
	{sql: "UNLOCK TABLES", engine: true},
	{sql: "UNLOCK TABLE", engine: true},
	{sql: "UNLOCK = obidx_t"},
	{sql: "UNLOCK obidx_t"},
	// LOCK STATS and UNLOCK STATS are TiDB's own statements, which OceanBase
	// has not: 1064 for both, with a target and without one.
	{sql: "LOCK STATS obidx_t"},
	{sql: "UNLOCK STATS obidx_t"},
	{sql: "UNLOCK STATS"},
	{sql: "LOCK STATS obidx_t PARTITION (p0)"},
}

// TestIndexClauseForms holds the parser to the recorded OceanBase verdicts. It
// needs no server: the parser must read exactly what CE 4.3.5.6 parses, minus
// the pending rows.
func TestIndexClauseForms(t *testing.T) {
	for _, c := range indexClauseCases {
		t.Run(c.sql, func(t *testing.T) {
			want := c.engine && !c.pending
			_, err := Parse(c.sql)
			if got := err == nil; got != want {
				t.Errorf("parser accepts=%v, want %v for %q (OceanBase accepts=%v, pending=%v): %v",
					got, want, c.sql, c.engine, c.pending, err)
			}
		})
	}
}
