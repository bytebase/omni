package parser

import "testing"

// TestColumnAndTableOptionGrammar holds column, index, partition and database
// options to sql_yacc.yy, as MySQL 8.0.46 answers them.
func TestColumnAndTableOptionGrammar(t *testing.T) {
	for _, sql := range []string{
		"ALTER TABLE t ADD c INT NOT SECONDARY",
		"ALTER TABLE t ADD c INT NOT NULL NOT SECONDARY",
		"ALTER TABLE t ADD c INT STORAGE DEFAULT",
		"ALTER TABLE t ENGINE_ATTRIBUTE ''",
		"ALTER TABLE t ENGINE_ATTRIBUTE = '{\"a\": [1, 2]}'",
		"ALTER TABLE t ADD c INT SECONDARY_ENGINE_ATTRIBUTE '1'",
		"ALTER TABLE t ADD KEY (c) ENGINE_ATTRIBUTE '{}'",
		"ALTER TABLE t ADD FULLTEXT (c) WITH PARSER ngram",
		"ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN ((a), (MAXVALUE)))",
		"ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN (1, (a)))",
		"ALTER DATABASE d READ ONLY 0x1",
		"ALTER DATABASE d READ ONLY = 1.5",
		"ALTER DATABASE d READ ONLY 1e0",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
	for _, sql := range []string{
		"ALTER TABLE t ENGINE_ATTRIBUTE 's'",
		"ALTER TABLE t ADD c INT ENGINE_ATTRIBUTE '{'",
		"ALTER TABLE t ADD KEY (c) SECONDARY_ENGINE_ATTRIBUTE 'x'",
		"ALTER TABLE t ADD FULLTEXT (c) WITH PARSER p",
		"ALTER TABLE t ADD c INT COLUMN_FORMAT",
		"ALTER TABLE t ADD c INT STORAGE",
		"ALTER TABLE t ADD c INT SRID",
		"ALTER TABLE t ADD PRIMARY (c)",
		"ALTER TABLE t ADD FOREIGN (c) REFERENCES u (c)",
		"ALTER TABLE t ADD c INT GENERATED AS (1)",
		"ALTER TABLE t ADD KEY (c())",
		"ALTER TABLE t ALTER INDEX i",
		"ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN (1, MAXVALUE))",
		"ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN (MAXVALUE, a))",
		"ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN ((1)))",
		"ALTER DATABASE d READ ONLY 2",
		"ALTER DATABASE d READ ONLY 0x2",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q): want an error", sql)
		}
	}
}
