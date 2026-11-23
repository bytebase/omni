package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitStatements(t *testing.T) {
	tests := []struct {
		name       string
		script     string
		want       []Statement
		completion []Statement
	}{
		{
			name:   "trailing comment trimmed",
			script: "SELECT 1 FROM DUAL; -- c\nSELECT 2 FROM DUAL ",
			want: []Statement{
				{Text: "SELECT 1 FROM DUAL", Start: 0, End: 18},
				{Text: " -- c\nSELECT 2 FROM DUAL", Start: 19, End: 43},
			},
			completion: []Statement{
				{Text: "SELECT 1 FROM DUAL", Start: 0, End: 18},
				{Text: " -- c\nSELECT 2 FROM DUAL ", Start: 19, End: 44},
			},
		},
		{
			name:   "SQL*Plus command dropped",
			script: "SET DEFINE OFF\nSELECT 1 FROM DUAL",
			want:   []Statement{{Text: "SELECT 1 FROM DUAL", Start: 15, End: 33}},
		},
		{
			name:   "MATCH_RECOGNIZE DEFINE rejoined",
			script: "SELECT * FROM T MATCH_RECOGNIZE (ORDER BY A MEASURES A.B AS P PATTERN (A)\nDEFINE A AS B > 1)",
			want: []Statement{{
				Text:  "SELECT * FROM T MATCH_RECOGNIZE (ORDER BY A MEASURES A.B AS P PATTERN (A)\nDEFINE A AS B > 1)",
				Start: 0,
				End:   92,
			}},
		},
		{
			name:   "PL/SQL block ended by slash",
			script: "BEGIN NULL; END;\n/\nSELECT 1 FROM DUAL",
			want: []Statement{
				{Text: "BEGIN NULL; END;", Start: 0, End: 16},
				{Text: "\nSELECT 1 FROM DUAL", Start: 18, End: 37},
			},
		},
		{
			name:   "empty segments kept untrimmed",
			script: ";SELECT 1 FROM DUAL; -- c\n;\n",
			want: []Statement{
				{Text: "", Empty: true, Start: 0, End: 0},
				{Text: "SELECT 1 FROM DUAL", Start: 1, End: 19},
				{Text: " -- c\n", Empty: true, Start: 20, End: 26},
				{Text: "\n", Empty: true, Start: 27, End: 28},
			},
		},
		{
			name:   "MATCH_RECOGNIZE DEFINE rejoined up to its semicolon",
			script: "SELECT * FROM T MATCH_RECOGNIZE (ORDER BY A MEASURES A.B AS P PATTERN (A)\nDEFINE A AS B > 1)\n;\nSELECT 2 FROM DUAL;",
			want: []Statement{
				{Text: "SELECT * FROM T MATCH_RECOGNIZE (ORDER BY A MEASURES A.B AS P PATTERN (A)\nDEFINE A AS B > 1)", Start: 0, End: 92},
				{Text: "\nSELECT 2 FROM DUAL", Start: 94, End: 113},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, SplitStatements(test.script))
			if test.completion != nil {
				require.Equal(t, test.completion, SplitStatementsForCompletion(test.script))
			}
		})
	}
}
