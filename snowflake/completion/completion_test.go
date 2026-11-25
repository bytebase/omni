package completion

import (
	"context"
	"strings"
	"testing"
)

// fuzzSeeds are statements whose every caret position the sweep and the fuzz
// corpus complete.
var fuzzSeeds = []string{
	"SELECT o.id, c.name FROM orders o JOIN customers c ON o.customer_id = c.id WHERE o.id > 1 QUALIFY ROW_NUMBER() OVER (PARTITION BY c.id ORDER BY o.id) = 1",
	"WITH RECURSIVE r (n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM r WHERE n < 3) SELECT * FROM r",
	"SELECT f.value:name::string FROM orders o, LATERAL FLATTEN(input => o.src:items) f",
	"SELECT $1, t.$2 FROM @mystage/path/ (FILE_FORMAT => 'csv') t",
	"SELECT src:customer.address[0].city FROM db.public.orders AT(OFFSET => -60) o",
	"INSERT INTO orders (id, customer_id) SELECT id, id FROM customers",
	"UPDATE orders o SET customer_id = c.id FROM customers c WHERE o.id = c.id",
	"MERGE INTO orders t USING customers s ON t.id = s.id WHEN MATCHED THEN UPDATE SET t.customer_id = s.id WHEN NOT MATCHED THEN INSERT (id) VALUES (s.id)",
	`SELECT "Order Items".qty FROM "Order Items" -- comment` + "\n" + `WHERE 'str' = $var`,
	"CREATE OR REPLACE TABLE t (a NUMBER(38, 0) NOT NULL, b VARIANT) CLUSTER BY (a)",
	"BEGIN LET x := 1; SELECT * FROM orders; END",
	"SELECT * FROM orders; SELECT * FROM customers WHERE",
	"((((((SELECT",
	"SELECT \"unterminated",
	"SELECT 'unterminated",
	"/* unterminated",
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

func TestCaretWordAt(t *testing.T) {
	tests := []struct {
		text   string // | is the caret
		ok     bool
		word   string // the bytes the placeholder replaces
		typed  string
		quoted bool
	}{
		{text: "SELECT |", ok: true},
		{text: "SELECT ab|", ok: true, word: "ab", typed: "ab"},
		{text: "SELECT a|bc FROM t", ok: true, word: "abc", typed: "a"},
		{text: "SELECT t.a|", ok: true, word: "a", typed: "a"},
		{text: `SELECT "a""b|`, ok: true, word: `"a""b`, typed: `a"b`, quoted: true},
		{text: `SELECT "ab|cd" FROM t`, ok: true, word: `"abcd"`, typed: "ab", quoted: true},
		{text: "SELECT 'ab|", ok: false},
		{text: "SELECT $ab|", ok: false},
		{text: "SELECT 1 -- a|", ok: false},
		{text: "SELECT 1 /* a|", ok: false},
	}
	for _, tt := range tests {
		i := strings.Index(tt.text, "|")
		stmt := tt.text[:i] + tt.text[i+1:]
		w, ok := caretWordAt(stmt, i)
		if ok != tt.ok {
			t.Errorf("caretWordAt(%q) ok = %v", tt.text, ok)
			continue
		}
		if !ok {
			continue
		}
		if got := stmt[w.start:w.end]; got != tt.word || w.typed != tt.typed || w.quoted != tt.quoted {
			t.Errorf("caretWordAt(%q) = word %q typed %q quoted %v; want %q %q %v", tt.text, got, w.typed, w.quoted, tt.word, tt.typed, tt.quoted)
		}
	}
}
