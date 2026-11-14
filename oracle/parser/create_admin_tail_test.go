package parser

import (
	"errors"
	"testing"
)

// TestAdminTailOracleAnswers pins CREATE and ALTER JAVA, CREATE and ALTER
// LIBRARY, CREATE MLE ENV and MLE MODULE, CREATE PMEM FILESTORE, and CREATE
// LOGICAL PARTITION TRACKING forms to Oracle 23ai's answers, probed on
// 23.26.3. The accepted forms fail on the server only past the parser:
// ORA-29538 or ORA-29548 (no Java VM), ORA-01031 (privilege), ORA-04043 (no
// such library), ORA-00439 (no PMEM), ORA-32373.
func TestAdminTailOracleAnswers(t *testing.T) {
	accepts := []string{
		"CREATE JAVA SOURCE USING BFILE(a, 's')",
		"CREATE JAVA SOURCE USING CLOB SELECT 1 FROM dual",
		"CREATE JAVA SOURCE USING CLOB (SELECT 1 FROM dual)",
		"CREATE JAVA SOURCE IF NOT EXISTS USING 'k'",
		"CREATE AND RESOLVE NOFORCE JAVA SOURCE USING 'k'",
		"CREATE JAVA SOURCE SHARING = NONE AUTHID DEFINER RESOLVER((a b)) USING 'k'",
		"CREATE JAVA SOURCE RESOLVER((a -)(* PUBLIC)) USING 'k'",
		"CREATE JAVA SOURCE RESOLVER() USING 'k'",
		"CREATE JAVA SOURCE NAMED USING USING 'k'",
		"ALTER JAVA SOURCE a RESOLVER((\"x/*\" b)(* PUBLIC)) RESOLVE",
		"ALTER JAVA SOURCE a RESOLVER((a b)) AUTHID DEFINER",
		"ALTER JAVA CLASS a COMPILE",
		"CREATE EDITIONABLE LIBRARY a IS 's'",
		"CREATE OR REPLACE NONEDITIONABLE LIBRARY a IS 's'",
		"ALTER LIBRARY IF EXISTS a EDITIONABLE",
		"ALTER LIBRARY COMPILE COMPILE",
		"ALTER LIBRARY a COMPILE plsql_optimize_level = 2 plsql_warnings = 'ENABLE:ALL' REUSE SETTINGS",
		"ALTER LIBRARY a COMPILE plsql_code_type = native",
		"CREATE MLE ENV IF NOT EXISTS a CLONE b",
		"CREATE MLE ENV a IMPORTS('a' MODULE a, N'b' MODULE b.c) LANGUAGE OPTIONS 's'",
		"CREATE MLE MODULE IF NOT EXISTS a LANGUAGE \"JAVASCRIPT\" VERSION 's' USING CLOB(SELECT 1 FROM dual x)",
		"CREATE MLE MODULE a LANGUAGE JAVASCRIPT USING CLOB((SELECT 1 FROM dual) UNION (SELECT 1 FROM dual))",
		"CREATE MLE MODULE a LANGUAGE JAVASCRIPT USING BFILE(\"a\", 's')",
		`CREATE PMEM FILESTORE "a"`,
		`CREATE PMEM FILESTORE "a" AUTOEXTEND ON NEXT 1 MOUNTPOINT 's' SIZE 1 M SIZE 2`,
		"CREATE LOGICAL PARTITION TRACKING ON a PARTITION BY RANGE(a, b)(PARTITION VALUES LESS THAN(1, 2), PARTITION \"p\" VALUES LESS THAN(MAXVALUE, 'x'))",
	}
	for _, sql := range accepts {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}

	rejects := []struct {
		sql string
		pos int
	}{
		{"CREATE IF NOT EXISTS JAVA SOURCE USING 'k'", 7},                                                                  // ORA-00901
		{"CREATE JAVA RESOURCE USING 'k'", 21},                                                                             // ORA-29502
		{"CREATE JAVA SOURCE RESOLVER((a b)) AUTHID DEFINER USING 'k'", 35},                                                // ORA-00922
		{"CREATE JAVA SOURCE RESOLVER((a b),(c d)) USING 'k'", 33},                                                         // ORA-29512
		{"CREATE JAVA SOURCE USING BFILE(a, 's') zzz", 39},                                                                 // ORA-03049
		{"ALTER JAVA IF EXISTS SOURCE a RESOLVE", 11},                                                                      // ORA-02000
		{"ALTER JAVA RESOURCE a RESOLVE", 11},                                                                              // ORA-02000
		{"ALTER JAVA SOURCE a RESOLVE AUTHID DEFINER", 28},                                                                 // ORA-03049
		{"ALTER JAVA SOURCE a AUTHID DEFINER COMPILE", 35},                                                                 // ORA-03049
		{"ALTER JAVA SOURCE a RESOLVER((a b c)) RESOLVE", 34},                                                              // ORA-29512
		{"ALTER JAVA SOURCE a RESOLVER((PUBLIC b)) RESOLVE", 30},                                                           // ORA-29512
		{"ALTER JAVA SOURCE a RESOLVER((a b))", 35},                                                                        // ORA-00922
		{"CREATE IF NOT EXISTS LIBRARY a IS 's'", 7},                                                                       // ORA-00901
		{"CREATE LIBRARY SELECT IS 's'", 15},                                                                               // ORA-04050
		{"CREATE LIBRARY a.IS IS 's'", 17},                                                                                 // ORA-04050
		{"ALTER LIBRARY a COMPILE EDITIONABLE", 24},                                                                        // ORA-03049
		{"ALTER LIBRARY a@b COMPILE", 15},                                                                                  // ORA-00922
		{"ALTER LIBRARY a COMPILE PERMIT_92_WRAP_FORMAT = 1", 24},                                                          // ORA-03049
		{"ALTER LIBRARY a COMPILE plsql_optimize_level = 1 + 1", 49},                                                       // ORA-03048
		{"ALTER LIBRARY a COMPILE plsql_optimize_level = 'x'", 47},                                                         // ORA-02017
		{"ALTER LIBRARY a COMPILE plsql_optimize_level = 1 DEBUG", 49},                                                     // ORA-03049
		{"CREATE PURE MLE ENV a", 7},                                                                                       // ORA-00901
		{"CREATE IF NOT EXISTS MLE MODULE a LANGUAGE JAVASCRIPT AS a", 7},                                                  // ORA-00901
		{"CREATE OR REPLACE MLE ENV IF NOT EXISTS a", 40},                                                                  // ORA-11541
		{"CREATE MLE ENV IF NOT a", 15},                                                                                    // ORA-11543
		{"CREATE MLE ENV a CLONE b IMPORTS('a' MODULE a)", 25},                                                             // ORA-03049
		{"CREATE MLE ENV a IMPORTS(\"a\" MODULE a)", 25},                                                                   // ORA-04114
		{"CREATE MLE ENV a LANGUAGE OPTIONS 's' 't'", 38},                                                                  // ORA-03046
		{"CREATE MLE MODULE a LANGUAGE JAVASCRIPT VERSION a AS a", 48},                                                     // ORA-04112
		{"CREATE MLE MODULE a LANGUAGE JAVASCRIPT USING CLOB(SELECT 1 FROM dual ORDER BY 1)", 46},                          // ORA-29506
		{"CREATE MLE MODULE a LANGUAGE JAVASCRIPT USING CLOB(WITH q AS (SELECT 1 c FROM dual) SELECT c FROM q)", 46},       // ORA-29506
		{"CREATE MLE MODULE a LANGUAGE JAVASCRIPT USING BLOB(SELECT 1 FROM dual) zzz", 46},                                 // ORA-29506
		{"CREATE PMEM FILESTORE a MOUNTPOINT 's'", 22},                                                                     // ORA-00931
		{`CREATE PMEM FILESTORE "a" AUTOEXTEND ON MAXSIZE 1G NEXT 1G`, 51},                                                 // ORA-00922
		{`CREATE PMEM FILESTORE "a" SIZE 1.5G`, 31},                                                                        // ORA-02017
		{"CREATE LOGICAL PARTITION TRACKING ON a PARTITION BY INTERVAL(a)(PARTITION a VALUES LESS THAN(1))", 52},           // ORA-14151
		{"CREATE LOGICAL PARTITION TRACKING ON a@b PARTITION BY RANGE(a)(PARTITION a VALUES LESS THAN(1))", 38},            // ORA-00949
		{"CREATE LOGICAL PARTITION TRACKING ON a PARTITION BY RANGE(a)(PARTITION a VALUES LESS THAN(1) TABLESPACE t)", 93}, // ORA-00922
	}
	for _, tt := range rejects {
		_, err := Parse(tt.sql)
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q) = %v, want a ParseError", tt.sql, err)
			continue
		}
		if pe.Position != tt.pos {
			t.Errorf("Parse(%q) failed at %d, want %d", tt.sql, pe.Position, tt.pos)
		}
	}
}
