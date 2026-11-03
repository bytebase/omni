package completion

import (
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/oceanbase/catalog"
	"github.com/bytebase/omni/oceanbase/parser"
)

// cursorMarker marks the cursor in the scripts below.
const cursorMarker = "<|>"

// splitCursor removes cursorMarker from s and returns the text and the marker's
// byte offset.
func splitCursor(t *testing.T, s string) (string, int) {
	t.Helper()
	i := strings.Index(s, cursorMarker)
	if i < 0 {
		t.Fatalf("no cursor marker in %q", s)
	}
	return s[:i] + s[i+len(cursorMarker):], i
}

// refNames renders refs as "db.table alias" for comparison.
func refNames(refs []TableRef) []string {
	var names []string
	for _, r := range refs {
		n := r.Table
		if r.Database != "" {
			n = r.Database + "." + n
		}
		if r.Alias != "" {
			n += " " + r.Alias
		}
		names = append(names, n)
	}
	return names
}

// TestExtractTableRefsOpenToken pins the table scope of a statement whose last
// token is still open: a string, quoted identifier or comment the user has not
// closed, usually after the cursor. Parse rejects the statement, as the server
// does, and completion must still read its tables as it reads them once the
// token is closed. Every statement joins two tables with a comma, which the
// lexer fallback reads only the first of, so a pass here means the statement
// was parsed.
func TestExtractTableRefsOpenToken(t *testing.T) {
	tests := []struct {
		name   string
		open   string // the statement as typed, with the cursor marked
		closed string // the same statement with its open token closed
	}{
		{
			name:   "string after the cursor",
			open:   "SELECT <|>* FROM users u, orders o WHERE o.status = 'pend",
			closed: "SELECT <|>* FROM users u, orders o WHERE o.status = 'pend'",
		},
		{
			name:   "double-quoted string after the cursor",
			open:   `SELECT <|>* FROM users u, orders o WHERE o.status = "pend`,
			closed: `SELECT <|>* FROM users u, orders o WHERE o.status = "pend"`,
		},
		{
			name:   "string with a doubled quote",
			open:   "SELECT <|>* FROM users u, orders o WHERE u.name = 'it''s",
			closed: "SELECT <|>* FROM users u, orders o WHERE u.name = 'it''s'",
		},
		{
			name:   "string ending in a backslash",
			open:   `SELECT <|>* FROM users u, orders o WHERE u.name = 'ab\`,
			closed: `SELECT <|>* FROM users u, orders o WHERE u.name = 'ab'`,
		},
		{
			name:   "quoted identifier after the cursor",
			open:   "SELECT <|>* FROM users u, orders o WHERE o.`status",
			closed: "SELECT <|>* FROM users u, orders o WHERE o.`status`",
		},
		{
			name:   "quoted user variable after the cursor",
			open:   "SELECT <|>* FROM users u, orders o WHERE @`v",
			closed: "SELECT <|>* FROM users u, orders o WHERE @`v`",
		},
		{
			name:   "comment after the cursor",
			open:   "SELECT <|>* FROM users u, orders o /* draft",
			closed: "SELECT <|>* FROM users u, orders o /* draft */",
		},
		{
			name:   "executable comment after the cursor",
			open:   "SELECT <|>* FROM users u, orders o WHERE o.oid = 1 /*! AND u.id = 2",
			closed: "SELECT <|>* FROM users u, orders o WHERE o.oid = 1 /*! AND u.id = 2 */",
		},
		{
			name:   "string inside an open executable comment",
			open:   "SELECT <|>* FROM users u, orders o WHERE o.oid = 1 /*! AND u.name = 'ab",
			closed: "SELECT <|>* FROM users u, orders o WHERE o.oid = 1 /*! AND u.name = 'ab' */",
		},
		{
			name:   "executable comment ending in a line comment",
			open:   "SELECT <|>* FROM users u, orders o WHERE o.oid = 1 /*! AND u.id = 2 -- note",
			closed: "SELECT <|>* FROM users u, orders o WHERE o.oid = 1 /*! AND u.id = 2 -- note\n*/",
		},
		{
			name:   "cursor inside the open string",
			open:   "SELECT * FROM users u, orders o WHERE o.status = 'pe<|>",
			closed: "SELECT * FROM users u, orders o WHERE o.status = 'pe<|>'",
		},
		{
			name:   "cursor inside the open comment",
			open:   "SELECT * FROM users u, orders o /* dr<|>",
			closed: "SELECT * FROM users u, orders o /* dr<|> */",
		},
		{
			name:   "statement after a DELIMITER directive",
			open:   "DELIMITER //\nSELECT <|>* FROM users u, orders o WHERE o.status = 'pend",
			closed: "DELIMITER //\nSELECT <|>* FROM users u, orders o WHERE o.status = 'pend'",
		},
		{
			name:   "subquery tables stay out of scope",
			open:   "SELECT <|>* FROM users u WHERE u.id IN (SELECT user_id FROM orders) AND u.name = 'ab",
			closed: "SELECT <|>* FROM users u WHERE u.id IN (SELECT user_id FROM orders) AND u.name = 'ab'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			open, openCursor := splitCursor(t, tt.open)
			closed, closedCursor := splitCursor(t, tt.closed)
			if _, err := parser.Parse(closed); err != nil {
				t.Fatalf("Parse(%q): %v", closed, err)
			}
			want := refNames(extractTableRefs(closed, closedCursor))
			if len(want) == 0 {
				t.Fatalf("extractTableRefs(%q) found no tables", closed)
			}
			got := refNames(extractTableRefs(open, openCursor))
			if !slices.Equal(got, want) {
				t.Errorf("extractTableRefs(%q) = %q, want %q as for the closed statement", open, got, want)
			}
		})
	}
}

// TestExtractTableRefsCursorStatement pins that the table scope comes from the
// statement holding the cursor and from no other: a statement that does not
// parse cannot empty it, and the tables another statement names stay out.
func TestExtractTableRefsCursorStatement(t *testing.T) {
	tests := []struct {
		name   string
		script string
		want   []string
	}{
		{
			name:   "open string in a later statement",
			script: "SELECT <|>* FROM users u;\nSELECT * FROM orders WHERE status = 'pend",
			want:   []string{"users u"},
		},
		{
			name:   "open comment in a later statement",
			script: "SELECT <|>* FROM users u;\n/* draft",
			want:   []string{"users u"},
		},
		{
			name:   "syntax error in an earlier statement",
			script: "SELECT * FROM FROM orders;\nSELECT <|>* FROM users u",
			want:   []string{"users u"},
		},
		{
			name:   "tables of an earlier statement",
			script: "SELECT * FROM orders;\nSELECT <|>* FROM users u",
			want:   []string{"users u"},
		},
		{
			name:   "tables of a later statement",
			script: "SELECT <|>* FROM users u;\nSELECT * FROM orders",
			want:   []string{"users u"},
		},
		{
			name:   "cursor at the end of a statement, before its semicolon",
			script: "SELECT * FROM users u WHERE <|>;SELECT * FROM orders WHERE status = 'pend",
			want:   []string{"users u"},
		},
		{
			name:   "cursor at the start of a statement, after the semicolon",
			script: "SELECT * FROM orders;<|>SELECT * FROM users u",
			want:   []string{"users u"},
		},
		{
			name:   "empty statement before the cursor",
			script: "SELECT * FROM orders;;SELECT <|>* FROM users u WHERE u.name = 'ab",
			want:   []string{"users u"},
		},
		{
			name:   "closed string in an earlier statement",
			script: "SELECT * FROM orders WHERE status = 'pend';\nSELECT <|>* FROM users u",
			want:   []string{"users u"},
		},
		{
			name:   "custom delimiter",
			script: "DELIMITER //\nSELECT * FROM orders//\nSELECT <|>* FROM users u WHERE u.name = 'ab",
			want:   []string{"users u"},
		},
		{
			name:   "cursor past the last statement",
			script: "SELECT * FROM orders; <|>",
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, cursor := splitCursor(t, tt.script)
			got := refNames(extractTableRefs(sql, cursor))
			if !slices.Equal(got, tt.want) {
				t.Errorf("extractTableRefs(%q, %d) = %q, want %q", sql, cursor, got, tt.want)
			}
		})
	}
}

// TestExtractTableRefsLexerReadsStatement pins that when the statement does not
// parse even with its open token closed, the lexer fallback reads the whole
// statement, including the FROM clause after the cursor.
func TestExtractTableRefsLexerReadsStatement(t *testing.T) {
	tests := []struct {
		script string
		want   []string
	}{
		{
			script: "SELECT <|> FROM users u WHERE u.name = (",
			want:   []string{"users u"},
		},
		{
			script: "SELECT o.<|> FROM orders o WHERE o.oid IN (SELECT id FROM users WHERE name = 'ab",
			want:   []string{"orders o", "users"},
		},
	}
	for _, tt := range tests {
		sql, cursor := splitCursor(t, tt.script)
		got := refNames(extractTableRefs(sql, cursor))
		if !slices.Equal(got, tt.want) {
			t.Errorf("extractTableRefs(%q, %d) = %q, want %q", sql, cursor, got, tt.want)
		}
	}
}

// TestCompleteColumnsWithOpenToken pins the column candidates of the scripts
// that lost their table scope once the lexer rejected an unterminated token.
func TestCompleteColumnsWithOpenToken(t *testing.T) {
	cat := catalog.New()
	mustExec(t, cat, "CREATE DATABASE test")
	cat.SetCurrentDatabase("test")
	mustExec(t, cat, "CREATE TABLE users (id INT, name VARCHAR(100), email VARCHAR(200))")
	mustExec(t, cat, "CREATE TABLE orders (oid INT, user_id INT, status VARCHAR(20))")

	tests := []string{
		"SELECT u.<|> FROM users u WHERE u.name = 'ab",
		"SELECT u.<|> FROM users u;\nSELECT * FROM orders WHERE status = 'pend",
		"SELECT u.<|> FROM users u;\n/* draft",
		"SELECT <|> FROM users WHERE name = 'ab",
	}
	want := []string{"email", "id", "name"}
	for _, script := range tests {
		sql, cursor := splitCursor(t, script)
		var got []string
		for _, c := range Complete(sql, cursor, cat) {
			if c.Type == CandidateColumn {
				got = append(got, c.Text)
			}
		}
		slices.Sort(got)
		if !slices.Equal(got, want) {
			t.Errorf("Complete(%q, %d) columns = %q, want %q", sql, cursor, got, want)
		}
	}
}
