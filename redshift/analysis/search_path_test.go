package analysis

import "testing"

// TestSplitSearchPath pins what the driver's row estimate splits off a
// statement: the SET LOCAL search_path a caller prefixed, and the statement it
// was prefixed for. Anything else is one statement, prefix-shaped or not.
func TestSplitSearchPath(t *testing.T) {
	statement := "DELETE FROM t WHERE a = 1;"
	prefixed := withSearchPath(statement, []string{"app", `we"ird`})
	for _, tc := range []struct {
		name      string
		in        string
		wantSetup string
		wantStmt  string
	}{
		{
			name:      "a prefixed statement",
			in:        prefixed,
			wantSetup: `SET LOCAL search_path TO "app", "we""ird"`,
			wantStmt:  statement,
		},
		{
			// A schema's name may hold the delimiter text inside its quotes.
			name:      "a prefix whose schema holds a semicolon",
			in:        withSearchPath(statement, []string{"app;\nname"}),
			wantSetup: "SET LOCAL search_path TO \"app;\nname\"",
			wantStmt:  statement,
		},
		{
			name:     "a statement without a prefix",
			in:       statement,
			wantStmt: statement,
		},
		{
			name:     "a statement that only looks prefixed",
			in:       "SET LOCAL search_path = app; DELETE FROM t;",
			wantStmt: "SET LOCAL search_path = app; DELETE FROM t;",
		},
		{
			name:     "a prefix with nothing after it",
			in:       searchPathPrefix,
			wantStmt: searchPathPrefix,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setup, rest := SplitSearchPath(tc.in)
			if setup != tc.wantSetup {
				t.Errorf("setup = %q, want %q", setup, tc.wantSetup)
			}
			if rest != tc.wantStmt {
				t.Errorf("statement = %q, want %q", rest, tc.wantStmt)
			}
		})
	}
}
