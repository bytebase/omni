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
