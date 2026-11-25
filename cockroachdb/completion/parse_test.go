package completion

import (
	"testing"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/lexbase"
	"github.com/stretchr/testify/require"
)

func TestReachesPast(t *testing.T) {
	require.True(t, reachesPast("SELECT * FROM t WHERE", 16))
	require.False(t, reachesPast("SELECT * FROM t SELECT", 16))
	require.True(t, reachesPast("SELECT 1", 7))
	// Leading whitespace, comments and semicolons do not shift positions.
	require.True(t, reachesPast(" ; /* c */\n SELECT * FROM t WHERE", 32))
	require.False(t, reachesPast(" ; /* c */\n SELECT * FROM t SELECT", 32))
	// A token on a later line.
	require.True(t, reachesPast("SELECT *\nFROM t\nWHERE", 16))
	require.False(t, reachesPast("SELECT *\nFROM t\nSELECT", 16))
}

func TestErrorOffset(t *testing.T) {
	for _, tc := range []struct {
		text string
		want int
	}{
		{"SELECT * FROM t SELECT", 16},
		{"  SELECT * FROM t SELECT", 18},
		{"SELECT *\nFROM t SELECT", 16},
		{"SELECT * FROM", 13},
	} {
		_, err := parseOne(tc.text)
		require.Error(t, err, tc.text)
		require.Equal(t, tc.want, errorOffset(tc.text, err), tc.text)
	}
	require.Equal(t, -1, errorOffset("x", errNotOneStatement))
}

func TestParseOne(t *testing.T) {
	stmt, err := parseOne("SELECT 1")
	require.NoError(t, err)
	require.NotNil(t, stmt)
	_, err = parseOne("SELECT 1; SELECT 2")
	require.ErrorIs(t, err, errNotOneStatement)
	_, err = parseOne("")
	require.ErrorIs(t, err, errNotOneStatement)
	_, err = parseOne("DO $$ BEGIN END $$")
	require.NoError(t, err)
}

func TestTokenize(t *testing.T) {
	toks := tokenize(`SELECT "Ab" FROM t`)
	require.Len(t, toks, 4)
	require.Equal(t, token{id: lexbase.IDENT, start: 7, end: 12, str: "Ab"}, toks[1])

	toks = tokenize(`SELECT "Ab`)
	last := toks[len(toks)-1]
	require.Equal(t, int32(lexbase.ERROR), last.id)
	require.True(t, last.quotedIdent)
	require.Equal(t, "Ab", last.str)

	toks = tokenize(`SELECT 'ab`)
	last = toks[len(toks)-1]
	require.Equal(t, int32(lexbase.ERROR), last.id)
	require.False(t, last.quotedIdent)

	require.Empty(t, tokenize(""))
}
