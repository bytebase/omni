package catalog

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"

	mysqldriver "github.com/go-sql-driver/mysql"

	obparser "github.com/bytebase/omni/oceanbase/parser"
)

// TestContainer_ReservedKeywordAcceptance holds omni's parser to OceanBase's
// on every keyword in seven name positions: wherever OceanBase reads the
// keyword as a name, omni must. A mismatch means isIdentToken or parseIdent
// is too restrictive for that position. The reverse, omni reading what
// OceanBase refuses, is leniency this test does not look for; the parser
// package's TestReservedWordSweepOceanBaseOracle holds both directions in its
// own positions.
//
// It reports every mismatch of a position rather than stopping at the first.
func TestContainer_ReservedKeywordAcceptance(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	keywords := getReservedKeywords()
	if len(keywords) == 0 {
		t.Fatal("no keywords found")
	}
	t.Logf("testing %d keywords", len(keywords))

	// Each template is a statement with %s where the keyword goes. setup
	// creates what the statement names, so that OceanBase's verdict is on the
	// keyword and nothing else.
	positions := []struct {
		name     string
		setup    string
		template string
	}{
		{
			name:     "CHARACTER SET value",
			template: "CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET %s)",
		},
		{
			name:     "COLLATE value",
			template: "CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE %s)",
		},
		{
			name:     "ENGINE value",
			template: "CREATE TABLE kw_eng_test (id INT) ENGINE=%s",
		},
		{
			name:     "INDEX name",
			setup:    "CREATE TABLE kw_idx_base (id INT, val INT)",
			template: "CREATE INDEX %s ON kw_idx_base (val)",
		},
		{
			name:     "CONSTRAINT name in CREATE TABLE",
			template: "CREATE TABLE kw_con_test (id INT, CONSTRAINT %s UNIQUE (id))",
		},
		{
			name:     "Column alias in SELECT",
			template: "SELECT 1 AS %s",
		},
		{
			name:     "Table alias in SELECT",
			setup:    "CREATE TABLE kw_alias_base (id INT)",
			template: "SELECT * FROM kw_alias_base %s",
		},
	}

	for _, pos := range positions {
		t.Run(pos.name, func(t *testing.T) {
			if pos.setup != "" {
				if err := ctr.execSQL(pos.setup); err != nil {
					t.Fatalf("container setup: %v", err)
				}
			}

			var mismatches []string
			for _, kw := range keywords {
				sql := fmt.Sprintf(pos.template, kw)
				engineReads := prepareParses(t, ctr, sql)
				if _, omniErr := obparser.Parse(sql); engineReads && omniErr != nil {
					mismatches = append(mismatches, fmt.Sprintf("  %q: OceanBase reads it, omni rejects it: %v", sql, omniErr))
				}
			}
			if len(mismatches) > 0 {
				t.Errorf("%d keywords read by OceanBase and rejected by omni:\n%s",
					len(mismatches), strings.Join(mismatches, "\n"))
			}
		})
	}
}

// prepareParses asks OceanBase whether it reads stmt, without running it: the
// statement travels as PREPARE's argument, one request, so the sweep creates
// nothing per keyword and its tape holds one short exchange for each.
// OceanBase answers 1064 for what its parser refuses; any other answer, such
// as an unknown character set (1115), means the statement was read. An error
// that is not the server's is no verdict and fails the test.
func prepareParses(t *testing.T, ctr *mysqlContainer, stmt string) bool {
	t.Helper()
	_, err := ctr.db.ExecContext(ctr.ctx, "PREPARE omni_kw_stmt FROM ?", stmt)
	if err == nil {
		return true
	}
	var me *mysqldriver.MySQLError
	if !errors.As(err, &me) {
		t.Fatalf("no verdict for %q: %v", stmt, err)
	}
	return me.Number != 1064
}

// getReservedKeywords returns every keyword the parser's lexer knows, in name
// order.
func getReservedKeywords() []string {
	var keywords []string
	seen := make(map[string]bool)
	// Keyword token types are in the range 700-1500.
	for tok := 700; tok < 1500; tok++ {
		name := obparser.TokenName(tok)
		if name == "" {
			continue
		}
		lower := strings.ToLower(name)
		if !seen[lower] {
			seen[lower] = true
			keywords = append(keywords, name)
		}
	}
	sort.Strings(keywords)
	return keywords
}
