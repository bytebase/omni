package catalog

import "testing"

// TestExecSetCharsetForms checks NAMES and the charset forms reach the session wherever they
// sit in the list, in every spelling a server reads.
func TestExecSetCharsetForms(t *testing.T) {
	for _, tc := range []struct{ sql, want string }{
		{"SET NAMES latin1", "latin1"},
		{"SET NAMES 'latin1'", "latin1"},
		{"SET NAMES binary", "binary"},
		{"SET @a = 1, NAMES latin1", "latin1"},
		{"SET NAMES latin1, @a = 1", "latin1"},
		{"SET CHARSET latin1, @a = 1", "latin1"},
		{"SET @a = 1, CHARACTER SET latin1", "latin1"},
		{"SET CHAR SET latin1", "latin1"},
		{"SET @a = 1, CHAR SET 'latin1', @b = 2", "latin1"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			c := New()
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatal(err)
			}
			for _, r := range results {
				if r.Error != nil {
					t.Fatal(r.Error)
				}
			}
			if c.charsetClient != tc.want {
				t.Errorf("charsetClient = %q, want %q", c.charsetClient, tc.want)
			}
		})
	}
}
