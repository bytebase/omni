package catalog

import (
	"fmt"
	"strings"
	"testing"

	mysqlparser "github.com/bytebase/omni/mysql/parser"
)

// TestOracle_ReservedKeywordAcceptance systematically tests whether omni
// and MySQL 8.0 agree on which reserved keywords are accepted in various
// syntactic "name" positions. A mismatch (MySQL accepts, omni rejects)
// reveals a parser bug where isIdentToken/parseIdent is too restrictive.
//
// This is a diagnostic test — it reports all mismatches rather than failing
// on the first one, so we get a complete gap picture.
func TestOracle_ReservedKeywordAcceptance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
	defer cleanup()

	// All reserved keywords from mysql/parser/name.go.
	// We extract them programmatically from the parser's keyword table.
	reservedKeywords := getReservedKeywords()
	if len(reservedKeywords) == 0 {
		t.Fatal("no reserved keywords found")
	}
	t.Logf("testing %d reserved keywords", len(reservedKeywords))

	// Each template defines a syntactic position where an identifier might
	// be a reserved keyword. The %s placeholder is where the keyword goes.
	// We use backtick-quoting in setup SQL to avoid interfering.
	type namePosition struct {
		name     string
		setup    string // SQL to run before the test (on both oracle and omni)
		template string // SQL with %s for the keyword being tested
		cleanup  string // SQL to run after each keyword attempt
	}

	positions := []namePosition{
		{
			name:     "CHARACTER SET value",
			setup:    "CREATE DATABASE IF NOT EXISTS kw_test; USE kw_test",
			template: "CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET %s)",
			cleanup:  "DROP TABLE IF EXISTS kw_cs_test",
		},
		{
			name:     "COLLATE value",
			setup:    "CREATE DATABASE IF NOT EXISTS kw_test; USE kw_test",
			template: "CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE %s)",
			cleanup:  "DROP TABLE IF EXISTS kw_co_test",
		},
		{
			name:     "ENGINE value",
			setup:    "CREATE DATABASE IF NOT EXISTS kw_test; USE kw_test",
			template: "CREATE TABLE kw_eng_test (id INT) ENGINE=%s",
			cleanup:  "DROP TABLE IF EXISTS kw_eng_test",
		},
		{
			name:     "INDEX name",
			setup:    "CREATE DATABASE IF NOT EXISTS kw_test; USE kw_test; CREATE TABLE kw_idx_base (id INT, val INT)",
			template: "CREATE INDEX %s ON kw_idx_base (val)",
			cleanup:  "DROP INDEX %s ON kw_idx_base",
		},
		{
			name:     "CONSTRAINT name in CREATE TABLE",
			setup:    "CREATE DATABASE IF NOT EXISTS kw_test; USE kw_test",
			template: "CREATE TABLE kw_con_test (id INT, CONSTRAINT %s UNIQUE (id))",
			cleanup:  "DROP TABLE IF EXISTS kw_con_test",
		},
		{
			name:     "Column alias in SELECT",
			setup:    "CREATE DATABASE IF NOT EXISTS kw_test; USE kw_test",
			template: "SELECT 1 AS %s",
			cleanup:  "",
		},
		{
			name:     "Table alias in SELECT",
			setup:    "CREATE DATABASE IF NOT EXISTS kw_test; USE kw_test; CREATE TABLE kw_alias_base (id INT)",
			template: "SELECT * FROM kw_alias_base %s",
			cleanup:  "",
		},
	}

	for _, pos := range positions {
		t.Run(pos.name, func(t *testing.T) {
			// Setup
			if pos.setup != "" {
				if err := oracle.execSQL(pos.setup); err != nil {
					t.Fatalf("oracle setup: %v", err)
				}
			}

			var mismatches []string

			for _, kw := range reservedKeywords {
				sql := fmt.Sprintf(pos.template, kw)

				// Try on MySQL 8.0
				oracleErr := oracle.execSQL(sql)

				// Try on omni (parse-only — we just check if it parses)
				_, omniErr := mysqlparser.Parse(sql)

				// Cleanup on oracle
				if pos.cleanup != "" {
					cleanSQL := fmt.Sprintf(pos.cleanup, kw)
					oracle.execSQL(cleanSQL) //nolint:errcheck
				}

				oracleOK := oracleErr == nil
				omniOK := omniErr == nil

				if oracleOK && !omniOK {
					mismatches = append(mismatches, fmt.Sprintf(
						"  %s: MySQL accepts, omni rejects — %v", kw, omniErr))
				}
				// We don't care about the reverse (omni accepts, MySQL rejects)
				// — that's leniency, not a bug.
			}

			if len(mismatches) > 0 {
				t.Errorf("%d keywords accepted by MySQL but rejected by omni:\n%s",
					len(mismatches), strings.Join(mismatches, "\n"))
			} else {
				t.Logf("all %d keywords match between MySQL and omni", len(reservedKeywords))
			}
		})
	}
}

// getReservedKeywords returns all reserved keyword strings from the parser.
func getReservedKeywords() []string {
	// Use the parser's exported TokenName to reverse-map token types to names.
	// We check all keyword token types (>= 700) and filter for reserved ones.
	var keywords []string
	seen := make(map[string]bool)

	// Scan the full range of possible keyword token types.
	// MySQL keywords are in the range 700-1500 approximately.
	for tok := 700; tok < 1500; tok++ {
		name := mysqlparser.TokenName(tok)
		if name == "" {
			continue
		}
		// Check if this is a reserved keyword by trying to use it as an identifier.
		// If Parse rejects "SELECT <kw>" as a column alias (without AS), it's reserved.
		// But a simpler approach: just collect all keyword names and let the test filter.
		lower := strings.ToLower(name)
		if !seen[lower] {
			seen[lower] = true
			keywords = append(keywords, name)
		}
	}
	return keywords
}
