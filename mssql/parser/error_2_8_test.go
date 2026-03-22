package parser

import "testing"

// TestCreateTableTruncations tests that truncated CREATE TABLE statements
// return errors instead of silently succeeding (section 2.8).
func TestCreateTableTruncations(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "CREATE TABLE with no name",
			sql:  "CREATE TABLE",
		},
		{
			name: "CREATE TABLE open paren no columns",
			sql:  "CREATE TABLE t (",
		},
		{
			name: "CREATE TABLE column name no type",
			sql:  "CREATE TABLE t (a",
		},
		{
			name: "CREATE TABLE trailing comma no next column",
			sql:  "CREATE TABLE t (a INT,",
		},
		{
			name: "CREATE TABLE CONSTRAINT consumed no name",
			sql:  "CREATE TABLE t (a INT CONSTRAINT",
		},
		{
			name: "CREATE TABLE PRIMARY consumed no KEY",
			sql:  "CREATE TABLE t (a INT PRIMARY",
		},
		{
			name: "CREATE TABLE REFERENCES consumed no table",
			sql:  "CREATE TABLE t (a INT REFERENCES",
		},
		{
			name: "CREATE TABLE DEFAULT consumed no value",
			sql:  "CREATE TABLE t (a INT DEFAULT",
		},
		{
			name: "CREATE TABLE CHECK open paren no condition",
			sql:  "CREATE TABLE t (a INT CHECK (",
		},
		{
			name: "CREATE TABLE PK open paren no columns",
			sql:  "CREATE TABLE t (CONSTRAINT pk PRIMARY KEY (",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.sql)
			if err == nil {
				t.Errorf("Parse(%q): expected error for truncated input, got nil", tt.sql)
			}
		})
	}
}
