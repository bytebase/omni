package catalog

import "testing"

// TestDropTypeTypeName holds DROP TYPE and DROP DOMAIN on a type name with
// SETOF, modifiers or array bounds to what PostgreSQL 17.12 answers.
func TestDropTypeTypeName(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		code string // "" for success
	}{
		{"CREATE TYPE mood AS ENUM ('x'); DROP TYPE IF EXISTS SETOF mood; CREATE TYPE mood AS ENUM ('y')", ""},
		{"CREATE TYPE mood AS ENUM ('x'); DROP TYPE mood(1)", CodeSyntaxError},
		{"CREATE TYPE mood AS ENUM ('x'); DROP TYPE mood[]", CodeDependentObjects},
		{"CREATE TYPE mood AS ENUM ('x'); DROP TYPE mood ARRAY", CodeDependentObjects},
		{"DROP TYPE IF EXISTS nope[]", ""},
		{"DROP TYPE nope[]", CodeUndefinedObject},
		{"CREATE DOMAIN d AS int; DROP DOMAIN d[]", CodeWrongObjectType},
		{"CREATE DOMAIN d AS int; DROP DOMAIN d(3)", CodeSyntaxError},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			results, err := New().Exec(tc.sql, nil)
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
