package catalog

import "testing"

// TestCommentOnTypeName holds COMMENT ON TYPE and DOMAIN, whose object the
// parser gives as a TypeName, to what PostgreSQL 17.12 answers.
func TestCommentOnTypeName(t *testing.T) {
	setup := "CREATE TYPE mood AS ENUM ('x'); CREATE DOMAIN d AS int; "
	for _, tc := range []struct {
		sql  string
		code string // "" for success
	}{
		{"COMMENT ON TYPE mood IS 'a'", ""},
		{"COMMENT ON TYPE int IS 'b'", ""},
		{"COMMENT ON TYPE mood[] IS 'c'", ""},
		{"COMMENT ON TYPE mood(1) IS 'd'", CodeSyntaxError},
		{"COMMENT ON DOMAIN mood IS 'e'", CodeWrongObjectType},
		{"COMMENT ON DOMAIN d IS 'f'", ""},
		{"COMMENT ON TYPE d IS 'g'", ""},
		{"COMMENT ON TYPE nope IS 'h'", CodeUndefinedObject},
		{"COMMENT ON TYPE IS IS 'i'", CodeUndefinedObject},
		{"COMMENT ON DOMAIN public.d IS NULL", ""},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			results, err := New().Exec(setup+tc.sql, nil)
			if err != nil {
				t.Fatal(err)
			}
			var got error
			for _, r := range results {
				if r.Error != nil {
					got = r.Error
					break
				}
			}
			if tc.code == "" {
				if got != nil {
					t.Fatalf("got %v, want success", got)
				}
				return
			}
			assertCode(t, got, tc.code)
		})
	}
}
