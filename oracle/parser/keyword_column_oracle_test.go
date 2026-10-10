package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

var oracleKeywordColumnGaps = filepath.Join("testdata", "coverage", "oracle_keyword_column_gaps.tsv")

// keywordColumnContexts reference a column named by the keyword after it is
// defined; %[1]s is the table and %[2]s the keyword.
var keywordColumnContexts = []struct {
	name string
	sql  string
}{
	{"select list", "SELECT %[2]s FROM %[1]s"},
	{"where", "SELECT a FROM %[1]s WHERE %[2]s IS NULL"},
	{"not in", "SELECT a FROM %[1]s WHERE %[2]s NOT IN (1, 2)"},
	{"order by", "SELECT a FROM %[1]s ORDER BY %[2]s"},
	{"update", "UPDATE %[1]s SET %[2]s = NULL WHERE %[2]s IS NULL"},
	{"insert columns", "INSERT INTO %[1]s (%[2]s) VALUES (1)"},
	{"index", "CREATE INDEX %[1]s_i ON %[1]s (%[2]s)"},
	{"outer join marker", "SELECT d.dummy FROM dual d, %[1]s WHERE %[2]s(+) = LENGTH(d.dummy)"},
	{"case selector", "SELECT CASE %[2]s WHEN 1 THEN 2 END FROM %[1]s"},
}

// TestOracleNonReservedKeywordsAsColumns checks every word omni lexes as a
// keyword that Oracle does not reserve (V$RESERVED_WORDS) as a column name:
// defined by CREATE TABLE, then referenced from a select list, WHERE
// predicates (IS NULL, NOT IN), ORDER BY, UPDATE, an INSERT column list, an
// index, with the legacy outer-join marker (+), and as a simple CASE
// selector. omni must
// accept or reject each statement exactly as Oracle does. The keyword
// manifest's column_name context covers only the definition; OFFSET, FETCH,
// JOIN, MODEL, and USING failed in the references alone.
//
// Pseudo-column keywords are left out: a reference resolves to the
// pseudo-column (SYSTIMESTAMP compares as a TIMESTAMP), not to the column.
// Disagreements still open are recorded in oracle_keyword_column_gaps.tsv;
// the list only shrinks, and -update regenerates it.
func TestOracleNonReservedKeywordsAsColumns(t *testing.T) {
	o := startOracleDB(t)

	reserved := map[string]bool{}
	rows, err := o.adminDB.QueryContext(o.ctx, `
SELECT keyword FROM v$reserved_words
 WHERE keyword IS NOT NULL
   AND (reserved = 'Y' OR res_type = 'Y' OR res_attr = 'Y' OR res_semi = 'Y')`)
	if err != nil {
		t.Fatalf("query v$reserved_words: %v", err)
	}
	for rows.Next() {
		var keyword string
		if err := rows.Scan(&keyword); err != nil {
			t.Fatalf("scan v$reserved_words: %v", err)
		}
		reserved[keyword] = true
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate v$reserved_words: %v", err)
	}
	rows.Close()

	var words []string
	for word, tokenType := range oracleKeywords {
		if reserved[word] || isOraclePseudoColumnKeyword(tokenType) {
			continue
		}
		words = append(words, word)
	}
	sort.Strings(words)
	if len(words) < 200 {
		t.Fatalf("checked %d non-reserved keywords, want at least 200", len(words))
	}

	known := map[string]coverageRow{}
	if !*update {
		gaps, _ := readCoverageTSV(t, oracleKeywordColumnGaps)
		for _, row := range gaps {
			known[row.Key] = row
		}
	}

	run := fmt.Sprintf("%d", time.Now().UnixNano()%1_000_000_000)
	var gaps []string
	var newGaps, fixed, checks int
	// check records a disagreement with the table renamed to t, so the
	// manifest does not change from run to run.
	check := func(word, context, table, sqlText string, oracleAccepts bool) {
		checks++
		_, omniErr := Parse(sqlText)
		key := word + ":" + context
		_, wasKnown := known[key]
		if oracleAccepts == (omniErr == nil) {
			if wasKnown {
				fixed++
				t.Errorf("FIXED %s now matches Oracle: remove it from %s", key, oracleKeywordColumnGaps)
			}
			return
		}
		gaps = append(gaps, fmt.Sprintf("%s\t%v\t%v\t%s", key, oracleAccepts, omniErr == nil, strings.ReplaceAll(sqlText, table, "t")))
		if !wasKnown && !*update {
			newGaps++
			t.Errorf("NEW GAP %s: oracle_accepts=%v omni_err=%v sql=%q", key, oracleAccepts, omniErr, sqlText)
		}
	}

	for i, word := range words {
		table := fmt.Sprintf("omni_kw_%s_%d", run, i)
		create := fmt.Sprintf("CREATE TABLE %s (a NUMBER, %s NUMBER)", table, word)
		_, createErr := o.db.ExecContext(o.ctx, create)
		check(word, "column definition", table, create, createErr == nil)
		if createErr != nil {
			continue
		}
		for _, c := range keywordColumnContexts {
			sqlText := fmt.Sprintf(c.sql, table, word)
			check(word, c.name, table, sqlText, oracleParseOnly(o.ctx, o.db, sqlText) == nil)
		}
		if _, err := o.db.ExecContext(o.ctx, "DROP TABLE "+table+" PURGE"); err != nil {
			t.Logf("drop %s: %v", table, err)
		}
	}

	t.Logf("non-reserved keywords as columns: %d words, %d checks, %d disagree with Oracle (%d new, %d fixed)",
		len(words), checks, len(gaps), newGaps, fixed)

	if *update {
		content := "case\toracle_accepts\tomni_accepts\tsql\n" + strings.Join(gaps, "\n")
		if len(gaps) > 0 {
			content += "\n"
		}
		if err := os.WriteFile(oracleKeywordColumnGaps, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", oracleKeywordColumnGaps, err)
		}
		t.Logf("wrote %s (%d gaps)", oracleKeywordColumnGaps, len(gaps))
	}
}
