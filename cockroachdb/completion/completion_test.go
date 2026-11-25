package completion

import (
	"context"
	"math/rand/v2"
	"path/filepath"
	"testing"

	"github.com/bytebase/omni/enginetest/compat"
	"github.com/stretchr/testify/require"
)

func TestCaretWordAt(t *testing.T) {
	tests := []struct {
		input string
		want  caretWord
		ok    bool
	}{
		{input: "|", want: caretWord{}, ok: true},
		{input: "SEL|", want: caretWord{start: 0, end: 3, typed: "SEL"}, ok: true},
		{input: "SELECT * FROM ord|ers", want: caretWord{start: 14, end: 20, typed: "ord"}, ok: true},
		{input: "SELECT * FROM t.|", want: caretWord{start: 16, end: 16}, ok: true},
		{input: `SELECT * FROM "Ord|er" x`, want: caretWord{start: 14, end: 21, typed: "Ord", quoted: true}, ok: true},
		{input: `SELECT * FROM "ord"|`, want: caretWord{start: 19, end: 19}, ok: true},
		{input: "SELECT 'abc|", ok: false},
		{input: "SELECT 1 -- note|", ok: false},
		{input: "SELECT 'x'-- note|", ok: false},
		{input: "SELECT 1 /* note|", ok: false},
		{input: "SELECT 1 /* -- */ |", want: caretWord{start: 18, end: 18}, ok: true},
		{input: "SELECT 1 -- note\n|", want: caretWord{start: 17, end: 17}, ok: true},
	}
	for _, tc := range tests {
		stmt, line, col := caret(t, tc.input)
		got, ok := caretWordAt(stmt, lineOffsetToBytePos(stmt, line, col))
		require.Equal(t, tc.ok, ok, tc.input)
		if ok {
			require.Equal(t, tc.want, got, tc.input)
		}
	}
}

func TestFilterByPrefix(t *testing.T) {
	cands := []Candidate{{Text: "orders"}, {Text: "Order Items"}, {Text: "events"}}
	require.Equal(t, cands, filterByPrefix(cands, ""))
	require.Equal(t, []Candidate{{Text: "orders"}, {Text: "Order Items"}}, filterByPrefix(cands, "ORD"))
}

// FuzzComplete holds complete to its contract on any text and caret: it
// returns, without panicking, for statements of CockroachDB's grammar cut
// anywhere and for text that is not SQL at all.
func FuzzComplete(f *testing.F) {
	file, err := compat.ReadFile(filepath.Join("..", "parser", "testdata", "compat.tsv"))
	require.NoError(f, err)
	rng := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < len(file.Rows); i += len(file.Rows) / 120 {
		stmt := file.Rows[i].Statement
		f.Add(stmt, uint(rng.IntN(len(stmt)+1)))
	}
	for _, seed := range []string{
		"", "\x00\xff\xfe", "SELECT '", `SELECT "`, "SELECT /*", "SELECT $$", "SELECT e'\\", "??", ";;;",
		"SELECT * FROM t@", "SELECT ((((((((((", "WITH RECURSIVE r AS (SELECT * FROM r) SELECT * FROM r",
		"CREATE FUNCTION f() RETURNS INT LANGUAGE SQL BEGIN ATOMIC SELECT 1; END",
		"DO $$ BEGIN END $$", "SELECT * FROM [SHOW TABLES]",
	} {
		f.Add(seed, uint(0))
		f.Add(seed, uint(len(seed)))
		f.Add(seed, uint(len(seed)/2))
	}
	cat := newCatalog(testMetadata(), "public")
	f.Fuzz(func(t *testing.T, stmt string, pos uint) {
		_ = complete(context.Background(), stmt, int(pos%uint(len(stmt)+1)), cat)
	})
}
