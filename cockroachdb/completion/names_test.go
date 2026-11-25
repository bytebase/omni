package completion

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClosers(t *testing.T) {
	require.Equal(t, "", closers("SELECT (1)"))
	require.Equal(t, ")", closers("INSERT INTO t (a"))
	require.Equal(t, "])", closers("SELECT * FROM (SELECT * FROM [SHOW x"))
	require.Equal(t, "", closers("SELECT 1)"))
}

func TestPatchedTexts(t *testing.T) {
	texts := patchedTexts("UPDATE t SET p", " WHERE")
	require.Equal(t, "UPDATE t SET p WHERE", texts[0])
	require.Contains(t, texts, "UPDATE t SET p = 1")
	require.Contains(t, patchedTexts("INSERT INTO t (p", ""), "INSERT INTO t (p) VALUES (1)")
	for _, text := range texts {
		require.Equal(t, 1, countOf(texts, text), text)
	}
}

func countOf(list []string, s string) int {
	n := 0
	for _, v := range list {
		if v == s {
			n++
		}
	}
	return n
}

func TestMatchesQualifier(t *testing.T) {
	src := source{ref: "o", schema: "public", name: "orders"}
	require.True(t, matchesQualifier(src, nil))
	require.True(t, matchesQualifier(src, []string{"O"}))
	require.False(t, matchesQualifier(src, []string{"orders"}))
	require.True(t, matchesQualifier(src, []string{"public", "orders"}))
	require.False(t, matchesQualifier(src, []string{"other", "orders"}))
	require.True(t, matchesQualifier(source{ref: "orders", name: "orders"}, []string{"anything", "orders"}))
	require.False(t, matchesQualifier(source{ref: "c", derived: true}, []string{"s", "c"}))
}

func TestNameCandidatesWithoutCatalog(t *testing.T) {
	for _, kind := range []nameKind{kindObject, kindSchema, kindExpr, kindIndex} {
		nc := &nameContext{kind: kind, target: &source{ref: "t", name: "t"}, scope: []source{{ref: "t", name: "t"}}}
		require.Empty(t, nameCandidates(nc, nil, false))
	}
	nc := &nameContext{kind: kindObject, ctes: []string{"c"}}
	require.Equal(t, []Candidate{{Text: "c", Type: CandidateTable}}, nameCandidates(nc, nil, false))
}
