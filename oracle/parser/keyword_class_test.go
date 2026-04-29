package parser

import (
	"path/filepath"
	"testing"
)

func lexOne(t *testing.T, s string) Token {
	t.Helper()
	lex := NewLexer(s)
	tok := lex.NextToken()
	if tok.Type == tokEOF {
		t.Fatalf("NewLexer(%q) returned EOF", s)
	}
	return tok
}

func TestOracleKeywordClassificationGolden(t *testing.T) {
	tests := []struct {
		word string
		want oracleKeywordCategory
	}{
		// Statement and control keywords.
		{"SELECT", oracleKeywordReserved},
		{"TABLE", oracleKeywordReserved},
		{"RESOURCE", oracleKeywordReserved},
		{"UID", oracleKeywordReserved},

		// Clause starters.
		{"FROM", oracleKeywordClauseStarter},
		{"WHERE", oracleKeywordClauseStarter},
		{"GROUP", oracleKeywordClauseStarter},
		{"ORDER", oracleKeywordClauseStarter},
		{"UNION", oracleKeywordClauseStarter},

		// Type names.
		{"DATE", oracleKeywordType},
		{"NUMBER", oracleKeywordType},
		{"VARCHAR2", oracleKeywordType},
		{"TIMESTAMP", oracleKeywordType},
		{"BLOB", oracleKeywordType},

		// Function-like keywords.
		{"CAST", oracleKeywordFunction},
		{"DECODE", oracleKeywordFunction},
		{"JSON_VALUE", oracleKeywordFunction},
		{"JSON_OBJECT", oracleKeywordFunction},
		{"XMLAGG", oracleKeywordFunction},

		// Pseudo-columns.
		{"ROWID", oracleKeywordPseudoColumn},
		{"ROWNUM", oracleKeywordPseudoColumn},
		{"LEVEL", oracleKeywordPseudoColumn},
		{"SYSDATE", oracleKeywordPseudoColumn},

		// Context-sensitive feature keywords.
		{"JSON_TABLE", oracleKeywordContext},
		{"LATERAL", oracleKeywordContext},
		{"XMLTABLE", oracleKeywordContext},
		{"GENERATED", oracleKeywordContext},
		{"IDENTITY", oracleKeywordContext},

		// Lexer keywords not assigned to a narrower category.
		{"CASE", oracleKeywordNonReserved},
		{"NOCACHE", oracleKeywordNonReserved},

		// Plain and quoted identifiers.
		{"EMPLOYEES", oracleKeywordIdentifier},
		{`"SELECT"`, oracleKeywordIdentifier},
	}

	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			got := oracleKeywordCategoryOf(lexOne(t, tt.word))
			if got != tt.want {
				t.Fatalf("keyword category for %q = %s, want %s", tt.word, got, tt.want)
			}
		})
	}
}

func TestOracleReservedIdentifierGolden(t *testing.T) {
	tests := []struct {
		word string
		want bool
	}{
		{"SELECT", true},
		{"FROM", true},
		{"WHERE", true},
		{"TABLE", true},
		{"DATE", true},
		{"NUMBER", true},
		{"ROWID", true},
		{"RESOURCE", true},
		{"UID", true},
		{"JSON_TABLE", false},
		{"LATERAL", false},
		{"GENERATED", false},
		{"IDENTITY", false},
		{"EMPLOYEES", false},
		{`"SELECT"`, false},
	}

	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			if got := isOracleSQLReservedKeyword(lexOne(t, tt.word)); got != tt.want {
				t.Fatalf("isOracleSQLReservedKeyword(%q) = %v, want %v", tt.word, got, tt.want)
			}
		})
	}
}

func TestOracleKeywordCategoryStringGolden(t *testing.T) {
	tests := map[oracleKeywordCategory]string{
		oracleKeywordIdentifier:    "identifier",
		oracleKeywordReserved:      "reserved",
		oracleKeywordNonReserved:   "nonreserved",
		oracleKeywordContext:       "context",
		oracleKeywordType:          "type",
		oracleKeywordFunction:      "function",
		oracleKeywordPseudoColumn:  "pseudo-column",
		oracleKeywordClauseStarter: "clause-starter",
	}
	for category, want := range tests {
		if got := category.String(); got != want {
			t.Fatalf("%v.String() = %q, want %q", int(category), got, want)
		}
	}
}

func TestOracleQuotedReservedIdentifiers(t *testing.T) {
	cases := []string{
		`CREATE TABLE "SELECT" (a NUMBER)`,
		`CREATE TABLE t ("SELECT" NUMBER)`,
		`CREATE INDEX "SELECT" ON t(a)`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q) error = %v, want success", sql, err)
			}
		})
	}
}

func TestOracleKeywordExpressionGolden(t *testing.T) {
	cases := []string{
		`SELECT CAST(1 AS NUMBER) FROM dual`,
		`SELECT JSON_VALUE(payload, '$.id') FROM t`,
		`SELECT ROWNUM FROM dual`,
		`SELECT "ROWNUM" FROM dual`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q) error = %v, want success", sql, err)
			}
		})
	}
}

func TestOracleKeywordManifestExhaustive(t *testing.T) {
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "oracle_keywords.tsv"))
	manifest := make(map[string]coverageRow, len(rows))
	categoryCounts := make(map[string]int)
	for _, row := range rows {
		word := row.Fields["word"]
		if word == "" {
			t.Fatalf("%s: keyword manifest row has empty word", row.Key)
		}
		if _, exists := manifest[word]; exists {
			t.Fatalf("%s: duplicate keyword manifest row", word)
		}
		manifest[word] = row
		categoryCounts[row.Fields["category"]]++

		got := oracleKeywordCategoryOf(lexOne(t, word)).String()
		if got != row.Fields["category"] {
			t.Fatalf("%s: category = %s, want %s", word, got, row.Fields["category"])
		}
	}

	for word := range oracleKeywords {
		if _, ok := manifest[word]; !ok {
			t.Fatalf("lexer keyword %s is missing from oracle keyword manifest", word)
		}
	}
	if len(manifest) != len(oracleKeywords) {
		t.Fatalf("keyword manifest rows = %d, lexer keywords = %d", len(manifest), len(oracleKeywords))
	}
	t.Logf("Oracle keyword manifest rows=%d category=%v", len(manifest), categoryCounts)
}

func TestOracleKeywordOfficialSQLReserved26aiAudit(t *testing.T) {
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "oracle_keywords.tsv"))
	manifest := make(map[string]coverageRow, len(rows))
	for _, row := range rows {
		manifest[row.Fields["word"]] = row
	}

	officialRows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "oracle_sql_reserved_26ai.tsv"))
	for _, row := range officialRows {
		word := row.Fields["word"]
		kwRow, ok := manifest[word]
		if !ok {
			t.Fatalf("official Oracle SQL reserved word %s is missing from local keyword manifest", word)
		}
		tok := NewLexer(word).NextToken()
		if tok.Type == tokIDENT {
			t.Fatalf("official Oracle SQL reserved word %s is lexed as tokIDENT", word)
		}
		if row.Fields["generic_identifier"] == "reject" {
			if kwRow.Fields["table_name"] != "reject" || kwRow.Fields["column_name"] != "reject" {
				t.Fatalf("%s: official reserved word must reject generic identifiers, table=%s column=%s",
					word, kwRow.Fields["table_name"], kwRow.Fields["column_name"])
			}
		}
	}
	t.Logf("Oracle official SQL reserved 26ai audit rows=%d", len(officialRows))
}

func TestOracleKeywordManifestIdentifierContexts(t *testing.T) {
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "oracle_keywords.tsv"))
	for _, row := range rows {
		word := row.Fields["word"]
		for _, tc := range []struct {
			context string
			policy  string
			sql     string
		}{
			{context: "table_name", policy: row.Fields["table_name"], sql: "CREATE TABLE " + word + " (a NUMBER)"},
			{context: "column_name", policy: row.Fields["column_name"], sql: "CREATE TABLE t (" + word + " NUMBER)"},
			{context: "alias", policy: row.Fields["alias"], sql: "SELECT 1 " + word + " FROM dual"},
			{context: "dotted_name", policy: row.Fields["dotted_name"], sql: "SELECT 1 FROM sc." + word},
		} {
			t.Run(word+"/"+tc.context, func(t *testing.T) {
				assertKeywordPolicy(t, tc.sql, tc.policy)
			})
		}
	}
}

func TestOracleKeywordManifestExpressionContexts(t *testing.T) {
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "oracle_keywords.tsv"))
	for _, row := range rows {
		word := row.Fields["word"]
		if row.Fields["function_call"] != "skip" {
			t.Run(word+"/function_call", func(t *testing.T) {
				assertKeywordPolicy(t, "SELECT "+word+"(1) FROM dual", row.Fields["function_call"])
			})
		}
		if row.Fields["pseudocolumn"] != "skip" {
			t.Run(word+"/pseudocolumn", func(t *testing.T) {
				assertKeywordPolicy(t, "SELECT "+word+" FROM dual", row.Fields["pseudocolumn"])
			})
		}
	}
}

func assertKeywordPolicy(t *testing.T, sql string, policy string) {
	t.Helper()
	_, err := Parse(sql)
	switch policy {
	case "allow":
		if err != nil {
			t.Fatalf("Parse(%q) error = %v, want success", sql, err)
		}
	case "reject":
		if err == nil {
			t.Fatalf("Parse(%q) succeeded, want error", sql)
		}
	case "skip":
	default:
		t.Fatalf("unknown keyword policy %q", policy)
	}
}
