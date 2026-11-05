package completion

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/tidb/catalog"
)

// TestAsOfTableRefs checks the clause does not read as an alias named of, in the
// parsed path and in the lexer path a statement that does not parse falls back to.
// A wrong alias would hide the table's own name from `t.` completion.
func TestAsOfTableRefs(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		want []string
	}{
		{"select <|> from t as of timestamp now()", []string{"t"}},
		{"select <|> from t a as of timestamp now()", []string{"t a"}},
		{"select <|> from t as a as of timestamp now()", []string{"t a"}},
		{"select <|> from t as `of`", []string{"t of"}},
		{"select * from t as of <|>", []string{"t"}},
		{"select * from t a as of <|>", []string{"t a"}},
		{"select * from t as of timestamp now() join u on <|>", []string{"t", "u"}},
		// A statement the parser refuses, which leaves only the lexer.
		{"select * from t as of timestamp now() a join u on <|>", []string{"t", "u"}},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			sql, cursor := splitCursor(t, tc.sql)
			if got := refNames(extractTableRefs(sql, cursor)); !slices.Equal(got, tc.want) {
				t.Errorf("refs = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestAsOfCompletion checks what the clause offers at each word.
func TestAsOfCompletion(t *testing.T) {
	cat := catalog.New()
	if _, err := cat.Exec("CREATE DATABASE d; USE d; CREATE TABLE t (acol INT, bcol INT); CREATE TABLE u (ucol INT)", nil); err != nil {
		t.Fatal(err)
	}
	texts := func(sql string, cursor int) []string {
		var got []string
		for _, c := range Complete(sql, cursor, cat) {
			got = append(got, c.Text)
		}
		return got
	}
	for _, tc := range []struct {
		sql  string
		want []string
	}{
		{"select * from t as of <|>", []string{"TIMESTAMP"}},
		{"select * from t a as of <|>", []string{"TIMESTAMP"}},
		{"select * from t as of timestamp now() join u on <|>", []string{"acol", "bcol", "ucol"}},
		{"select a.<|> from t a as of timestamp now()", []string{"acol", "bcol"}},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			sql, cursor := splitCursor(t, tc.sql)
			got := texts(sql, cursor)
			slices.Sort(got)
			if !slices.Equal(got, tc.want) {
				t.Errorf("candidates = %v, want %v", got, tc.want)
			}
		})
	}
	t.Run("a column in the select list", func(t *testing.T) {
		sql, cursor := splitCursor(t, "select <|> from t as of timestamp now()")
		if got := texts(sql, cursor); !slices.Contains(got, "acol") {
			t.Errorf("candidates have no acol: %v", got)
		}
	})
}
