package completion

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	metadatapb "github.com/bytebase/omni/metadata"
)

// fuzzSeeds are statements whose every caret position the sweep and the fuzz
// corpus complete.
var fuzzSeeds = []string{
	"SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON o.customer_id = c.id WHERE o.amount > 1 ORDER BY o.id LIMIT 10",
	"WITH r AS (SELECT id FROM orders) SELECT * FROM r JOIN (SELECT id FROM customers) c ON r.id = c.id",
	"SELECT * FROM orders TABLET(10001) TABLESAMPLE(10 PERCENT) WHERE id IN (SELECT customer_id FROM db.orders)",
	"INSERT INTO orders (id, customer_id) SELECT id, id FROM customers",
	"UPDATE orders o SET amount = c.id FROM customers c WHERE o.id = c.id",
	"DELETE FROM orders PARTITION (p1) WHERE id = 1",
	"CREATE TABLE t (a INT, b VARCHAR(10)) DUPLICATE KEY(a) DISTRIBUTED BY HASH(a) BUCKETS 10 PROPERTIES (\"replication_num\" = \"1\")",
	"SELECT `order items`.qty FROM `order items` -- comment\nWHERE 'str' = \"dq\"",
	"SELECT * FROM orders; SELECT * FROM customers WHERE",
	"((((((SELECT",
	"SELECT `unterminated",
	"SELECT 'unterminated",
	"/* unterminated",
	"SELECT /*+ SET_VAR(a=1) */ id FROM orders",
	"",
}

// TestCompleteEveryCaret completes each seed at every byte offset: no input
// panics completion.
func TestCompleteEveryCaret(t *testing.T) {
	cfg := testConfig()
	for _, seed := range fuzzSeeds {
		for i := 0; i <= len(seed); i++ {
			lines := strings.Split(seed[:i], "\n")
			CompleteSession(context.Background(), cfg, seed, len(lines), len(lines[len(lines)-1]))
		}
	}
}

func FuzzCompleteSession(f *testing.F) {
	for _, seed := range fuzzSeeds {
		f.Add(seed, 1, len(seed)/2)
	}
	cfg := testConfig()
	f.Fuzz(func(t *testing.T, stmt string, line, offset int) {
		if len(stmt) > 4096 {
			return
		}
		CompleteSession(context.Background(), cfg, stmt, line, offset)
	})
}

// TestCompleteSessionLatency bounds one completion of a 2 KB statement over
// a catalog of 3,000 tables. Each keyword probe parses the whole statement,
// spread over every CPU; the median holds under a loaded machine, where
// single runs stall.
func TestCompleteSessionLatency(t *testing.T) {
	var tables []*metadatapb.TableMetadata
	for i := range 3000 {
		tables = append(tables, &metadatapb.TableMetadata{Name: fmt.Sprintf("t%d", i), Columns: []*metadatapb.ColumnMetadata{col("a", "INT"), col("b", "INT")}})
	}
	db := &metadatapb.DatabaseSchemaMetadata{Name: "big", Schemas: []*metadatapb.SchemaMetadata{{Tables: tables}}}
	cfg := SessionConfig{
		DefaultDatabase: "big",
		LoadDatabase: func(context.Context, string) (*metadatapb.DatabaseSchemaMetadata, error) {
			return db, nil
		},
	}
	var cols []string
	for i := 0; len(strings.Join(cols, ", ")) < 2048; i++ {
		cols = append(cols, fmt.Sprintf("t1.a + %d AS c%d", i, i))
	}
	stmt := "SELECT " + strings.Join(cols, ", ") + " FROM t1 WHERE t1.a > 1 AND |"
	var times []time.Duration
	for range 20 {
		start := time.Now()
		completeAt(t, cfg, stmt)
		times = append(times, time.Since(start))
	}
	slices.Sort(times)
	p50, p95 := times[len(times)/2], times[len(times)*95/100]
	t.Logf("p50 %v, p95 %v", p50, p95)
	if p50 > 500*time.Millisecond {
		t.Errorf("p50 %v over 500ms", p50)
	}
}

func TestCaretWordAt(t *testing.T) {
	tests := []struct {
		text   string // | is the caret
		ok     bool
		word   string // the bytes the placeholder replaces
		typed  string
		quoted bool
	}{
		{text: "SELECT |", ok: true},
		{text: "SELECT ab|cd FROM t", ok: true, word: "abcd", typed: "ab"},
		{text: "SELECT `a``b|", ok: true, word: "`a``b", typed: "a`b", quoted: true},
		{text: "SELECT `ab|c` FROM t", ok: true, word: "`abc`", typed: "ab", quoted: true},
		{text: "SELECT 'ab|", ok: false},
		{text: "SELECT /* ab|", ok: false},
		{text: "SELECT 1 -- ab|", ok: false},
	}
	for _, tc := range tests {
		i := strings.Index(tc.text, "|")
		stmt := tc.text[:i] + tc.text[i+1:]
		w, ok := caretWordAt(stmt, i)
		if ok != tc.ok {
			t.Errorf("%q: ok = %v, want %v", tc.text, ok, tc.ok)
			continue
		}
		if !ok {
			continue
		}
		if got := stmt[w.start:w.end]; got != tc.word || w.typed != tc.typed || w.quoted != tc.quoted {
			t.Errorf("%q: word %q typed %q quoted %v, want %q %q %v", tc.text, got, w.typed, w.quoted, tc.word, tc.typed, tc.quoted)
		}
	}
}

func TestInnermostQuery(t *testing.T) {
	tests := []struct {
		text string // | is the caret
		want string // the text completed, with the caret
	}{
		{"SELECT | FROM t", "SELECT | FROM t"},
		{"SELECT * FROM (SELECT | FROM t) x", "SELECT | FROM t"},
		{"SELECT * FROM (SELECT |", "SELECT |"},
		{"SELECT * FROM t WHERE a IN (SELECT b FROM u WHERE c IN (1, |)) AND d", "SELECT b FROM u WHERE c IN (1, |)"},
		{"SELECT (SELECT 1) + |", "SELECT (SELECT 1) + |"},
		{"SELECT count(|) FROM t", "SELECT count(|) FROM t"},
		{"CREATE TABLE t AS SELECT | FROM u", "SELECT | FROM u"},
		{"CREATE TABLE t AS SELECT * FROM u WHERE a IN (SELECT | FROM v)", "SELECT | FROM v"},
		{"CREATE TABLE t (a INT) COMMENT 'x' AS |", "CREATE TABLE t (a INT) COMMENT 'x' AS |"},
		{"WITH c AS (WITH d AS (SELECT 1) SELECT | FROM d) SELECT 1", "WITH d AS (SELECT 1) SELECT | FROM d"},
	}
	for _, tc := range tests {
		i := strings.Index(tc.text, "|")
		stmt := tc.text[:i] + tc.text[i+1:]
		body, caret := innermostQuery(stmt, i)
		if got := body[:caret] + "|" + body[caret:]; got != tc.want {
			t.Errorf("%q: got %q, want %q", tc.text, got, tc.want)
		}
	}
}
