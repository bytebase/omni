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
