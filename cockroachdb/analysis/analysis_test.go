package analysis

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSplit(t *testing.T) {
	for _, tc := range []struct {
		statement string
		want      []string
	}{
		{statement: "SELECT 1;", want: []string{"SELECT 1;"}},
		{statement: "SELECT\n\n\n\n1;", want: []string{"SELECT\n\n\n\n1;"}},
		{
			statement: "SELECT\n\n\n\n1;\n\n\nSELECT\n\n\n\n2;",
			want:      []string{"SELECT\n\n\n\n1;", "SELECT\n\n\n\n2;"},
		},
		{
			statement: "SELECT\n\n\n\n1;\n\n;\n\n\nSELECT\n\n\n\n2;",
			want:      []string{"SELECT\n\n\n\n1;", "SELECT\n\n\n\n2;"},
		},
	} {
		t.Run(tc.statement, func(t *testing.T) {
			got, err := Split(tc.statement)
			if err != nil {
				t.Fatalf("Split(%q): %v", tc.statement, err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Split() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseDoBlock(t *testing.T) {
	parsed, err := Parse(`DO $$ BEGIN DELETE FROM t; END $$`)
	if err != nil {
		t.Fatalf("Parse(DO block): %v", err)
	}
	if len(parsed) != 1 {
		t.Fatalf("Parse(DO block) returned %d statements, want 1", len(parsed))
	}
}
