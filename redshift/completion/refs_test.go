package completion

import "testing"

func TestExtractTableRefs(t *testing.T) {
	tests := []struct {
		sql    string
		offset int
		want   []string
	}{
		{"SELECT * FROM users WHERE ", 25, []string{"users"}},
		{"SELECT * FROM users u, orders o WHERE ", 37, []string{"users", "orders"}},
		// DDL: ALTER TABLE
		{"ALTER TABLE users RENAME COLUMN ", 31, []string{"users"}},
		{"ALTER TABLE IF EXISTS users DROP COLUMN ", 39, []string{"users"}},
		// DDL: CREATE INDEX
		{"CREATE INDEX idx ON orders (", 27, []string{"orders"}},
		{"CREATE UNIQUE INDEX idx ON orders (", 34, []string{"orders"}},
		// DDL: COMMENT ON COLUMN
		{"COMMENT ON COLUMN users.name IS ", 31, []string{"users"}},
		{"COMMENT ON COLUMN public.users.name IS ", 39, []string{"users"}},
	}
	for _, tt := range tests {
		refs := extractTableRefs(tt.sql, tt.offset)
		got := make(map[string]bool)
		for _, r := range refs {
			got[r.Table] = true
		}
		for _, w := range tt.want {
			if !got[w] {
				t.Errorf("extractTableRefs(%q, %d): missing %q", tt.sql, tt.offset, w)
			}
		}
	}
}
