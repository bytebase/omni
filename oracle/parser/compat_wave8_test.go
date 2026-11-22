package parser

import (
	"errors"
	"testing"
)

// TestWave8CompatReadings pins the Wave 8 grammar and byte-offset families
// moved to Oracle 23.26.3's exact byte offsets, together with valid Oracle SQL
// controls beside each refusal.
func TestWave8CompatReadings(t *testing.T) {
	refusedAt := []struct {
		sql  string
		want int
	}{
		// DELETE ... FROM ... MATCH_RECOGNIZE(...) parses the clause and reports
		// syntax errors inside it at their exact offending token (matching SELECT),
		// while still rejecting a syntactically valid MATCH_RECOGNIZE(...) in DML FROM
		// at MATCH_RECOGNIZE (byte 16).
		{`DELETE a FROM a MATCH_RECOGNIZE PATTERN(a) DEFINE a AS 1 = 1)`, 32},
		{`DELETE a FROM a MATCH_RECOGNIZE((a) DEFINE a AS 1 = 1)`, 32},
		{`DELETE a FROM a MATCH_RECOGNIZE(PATTERN a) DEFINE a AS 1 = 1)`, 40},
		{`DELETE a FROM a MATCH_RECOGNIZE(PATTERN() DEFINE a AS 1 = 1)`, 40},
		{`DELETE a FROM a MATCH_RECOGNIZE(PATTERN(a DEFINE a AS 1 = 1)`, 51},
		{`DELETE a FROM a MATCH_RECOGNIZE(PATTERN(a) a AS 1 = 1)`, 43},
		{`DELETE a FROM a MATCH_RECOGNIZE(PATTERN(a) DEFINE AS 1 = 1)`, 50},
		{`DELETE a FROM a MATCH_RECOGNIZE(PATTERN(a) DEFINE a 1 = 1)`, 52},
		{`DELETE a FROM a MATCH_RECOGNIZE(PATTERN(a) DEFINE a AS)`, 54},
		{`DELETE a FROM a MATCH_RECOGNIZE(PARTITION BY a zzz PATTERN(a) DEFINE a AS 1 = 1)`, 47},
		{`DELETE a FROM a MATCH_RECOGNIZE(PATTERN(a) DEFINE a AS 1 = 1)`, 16},
		{`DELETE a FROM a MATCH_RECOGNIZE(PATTERN(a) DEFINE a AS 1 = 1) zzz`, 16},
		{`DELETE a FROM ANALYTIC VIEW(USING a) MATCH_RECOGNIZE((a) DEFINE a AS 1 = 1)`, 53},
		{`DELETE a FROM ANALYTIC VIEW(USING a) MATCH_RECOGNIZE(PATTERN(a) DEFINE a AS 1 = 1)`, 37},

		// CREATE VIEW ... OF XMLTYPE parses XMLType_view_clause per CREATE-VIEW.bnf.
		// XMLSCHEMA takes an identifier / double-quoted URL (refusing single-quoted 's').
		{`CREATE VIEW a OF XMLTYPE XMLSCHEMA 's' ELEMENT a AS SELECT 1 FROM dual`, 35},
		{`CREATE VIEW a OF XMLTYPE 's' ELEMENT a AS SELECT 1 FROM dual`, 25},
		{`CREATE VIEW a OF XMLTYPE XMLTYPE 's' ELEMENT a AS SELECT 1 FROM dual`, 25},
		{`CREATE VIEW a OF XMLTYPE XMLSCHEMA ELEMENT a AS SELECT 1 FROM dual`, 43},
		{`CREATE VIEW a OF XMLTYPE XMLSCHEMA "http://example.com/po.xsd" ELEMENT "po" WITH OBJECT IDENTIFIER () AS SELECT 1 FROM dual`, 100},
		{`CREATE VIEW a OF XMLTYPE XMLSCHEMA "http://example.com/po.xsd" ELEMENT "po" WITH OBJECT DEFAULT AS SELECT 1 FROM dual`, 88},
		{`CREATE VIEW a OF XMLTYPE XMLSCHEMA s ELEMENT e WITH OBJECT IDENTIFIER (1,) AS SELECT 1 FROM dual`, 73},

		// CREATE TABLE ENABLE CONSTRAINTS accepts the plural synonym and checks
		// constraint presence just like ENABLE CONSTRAINT.
		{`CREATE TABLE a ENABLE CONSTRAINTS a`, 35},
		{`CREATE TABLE a(id NUMBER) ENABLE CONSTRAINTS missing_ck`, 55},

		// CREATE/ALTER ROLE IDENTIFIED GLOBALLY AS requires a directory/domain string literal.
		{`CREATE ROLE a IDENTIFIED GLOBALLY AS`, 36},
		{`ALTER ROLE a IDENTIFIED GLOBALLY AS`, 35},
		{`CREATE ROLE a IDENTIFIED GLOBALLY AS foo`, 37},
		{`ALTER ROLE a IDENTIFIED GLOBALLY AS foo`, 36},
	}

	for _, tc := range refusedAt {
		_, err := Parse(tc.sql)
		if err == nil {
			t.Errorf("Parse(%q) accepted; want rejection at byte %d", tc.sql, tc.want)
			continue
		}
		var pe *ParseError
		if !errors.As(err, &pe) || pe.Position != tc.want {
			t.Errorf("Parse(%q) failed with %v; want ParseError at byte %d", tc.sql, err, tc.want)
		}
	}

	// Valid Oracle SQL controls beside the Wave 8 refusals.
	accepted := []string{
		`SELECT * FROM a MATCH_RECOGNIZE(PATTERN(a) DEFINE a AS 1 = 1)`,
		`SELECT * FROM a MATCH_RECOGNIZE`,
		`DELETE a FROM a MATCH_RECOGNIZE`,
		`CREATE VIEW v OF XMLTYPE XMLSCHEMA "http://example.com/po.xsd" ELEMENT "PurchaseOrder" AS SELECT 1 FROM dual`,
		`CREATE VIEW v OF XMLTYPE XMLSCHEMA "http://example.com/po.xsd" ELEMENT "PurchaseOrder" WITH OBJECT IDENTIFIER DEFAULT AS SELECT 1 FROM dual`,
		`CREATE VIEW v OF XMLTYPE XMLSCHEMA "http://example.com/po.xsd" ELEMENT "PurchaseOrder" WITH OBJECT IDENTIFIER (1, 2) AS SELECT 1 FROM dual`,
		`CREATE VIEW v OF XMLTYPE XMLSCHEMA s1 ELEMENT e1 WITH OBJECT IDENTIFIER (id) AS SELECT 1 FROM dual`,
		`CREATE VIEW v OF XMLTYPE ELEMENT "PurchaseOrder" AS SELECT 1 FROM dual`,
		`CREATE VIEW v OF XMLTYPE WITH OBJECT IDENTIFIER DEFAULT AS SELECT 1 FROM dual`,
		`CREATE TABLE t (id NUMBER, CONSTRAINT ck_t CHECK (id > 0)) ENABLE CONSTRAINTS ck_t`,
		`CREATE TABLE t (id NUMBER, CONSTRAINT ck_t CHECK (id > 0)) DISABLE CONSTRAINTS ck_t`,
		`CREATE ROLE r IDENTIFIED GLOBALLY`,
		`CREATE ROLE r IDENTIFIED GLOBALLY AS 'cn=app_role,ou=roles'`,
		`ALTER ROLE r IDENTIFIED GLOBALLY AS 'cn=app_role,ou=roles'`,
	}

	for _, sql := range accepted {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q) failed: %v; want accept", sql, err)
		}
	}
}
