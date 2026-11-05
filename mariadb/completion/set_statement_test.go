package completion

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/mariadb/catalog"
)

// TestSetStatementCompletion checks SET STATEMENT offers what follows each of its words.
func TestSetStatementCompletion(t *testing.T) {
	cat := catalog.New()
	for _, tc := range []struct {
		sql  string
		want string
	}{
		{"set statement a = 1 <|>", "FOR"},
		{"set statement a = 1 for <|>", "SELECT"},
		{"set statement a = 1 for <|>", "DROP"},
	} {
		t.Run(tc.sql+tc.want, func(t *testing.T) {
			sql, cursor := splitCursor(t, tc.sql)
			var got []string
			for _, c := range Complete(sql, cursor, cat) {
				got = append(got, c.Text)
			}
			if !slices.Contains(got, tc.want) {
				t.Errorf("candidates = %v, want %s among them", got, tc.want)
			}
		})
	}

	// An option takes no @@ name, which is what the variable candidates are.
	sql, cursor := splitCursor(t, "set statement a = 1, <|>")
	for _, c := range Complete(sql, cursor, cat) {
		if c.Text != "" && c.Text[0] == '@' {
			t.Errorf("offered %s as an option", c.Text)
		}
	}
}
