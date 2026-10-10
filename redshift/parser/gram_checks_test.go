package parser

import (
	"errors"
	"testing"
)

// TestConstraintAttributes pins the constraint attributes each kind takes,
// as pg/parser's TestConstraintAttributesMatchPG does against PostgreSQL
// 17, from whose grammar this fork comes: gram.y's processCASbits and
// ConstraintAttributeSpec, and a column's constraint list, which takes no
// NOT VALID. Redshift documents none of these attributes for its
// constraints, so refusing them refuses nothing it runs; no Redshift
// server is reachable here to say more.
func TestConstraintAttributes(t *testing.T) {
	tests := []struct {
		sql     string
		message string // "" when Parse accepts the statement
	}{
		{"CREATE TABLE n (id int, PRIMARY KEY (id) NOT VALID)", "PRIMARY KEY constraints cannot be marked NOT VALID"},
		{"CREATE TABLE n (id int, UNIQUE (id) NO INHERIT)", "UNIQUE constraints cannot be marked NO INHERIT"},
		{"CREATE TABLE n (id int, EXCLUDE (id WITH =) NOT VALID)", "EXCLUDE constraints cannot be marked NOT VALID"},
		{"CREATE TABLE n (id int, CHECK (id > 0) DEFERRABLE)", "CHECK constraints cannot be marked DEFERRABLE"},
		{"CREATE TABLE n (id int, FOREIGN KEY (id) REFERENCES t NO INHERIT)", "FOREIGN KEY constraints cannot be marked NO INHERIT"},
		{"ALTER TABLE t ADD CONSTRAINT x PRIMARY KEY (x) NOT VALID", "PRIMARY KEY constraints cannot be marked NOT VALID"},
		{"ALTER TABLE t ALTER CONSTRAINT c NOT VALID", "FOREIGN KEY constraints cannot be marked NOT VALID"},
		{"CREATE CONSTRAINT TRIGGER tr AFTER INSERT ON t NO INHERIT FOR EACH ROW EXECUTE FUNCTION f()", "TRIGGER constraints cannot be marked NO INHERIT"},
		{"ALTER DOMAIN d ADD CONSTRAINT c NOT NULL NOT VALID", "NOT NULL constraints cannot be marked NOT VALID"},
		{"CREATE TABLE n (id int, PRIMARY KEY (id) NOT DEFERRABLE INITIALLY DEFERRED)", "constraint declared INITIALLY DEFERRED must be DEFERRABLE"},
		{"CREATE TABLE n (id int, PRIMARY KEY (id) DEFERRABLE NOT DEFERRABLE)", "conflicting constraint properties"},
		{"CREATE TABLE n (id int REFERENCES t NOT VALID)", `syntax error at or near "VALID"`},
		{"CREATE TABLE n (id int REFERENCES t NO INHERIT)", `syntax error at or near "NO"`},
		{"ALTER TABLE t ADD COLUMN id int REFERENCES t NOT VALID", `syntax error at or near "VALID"`},
		{"CREATE TABLE n (id int, PRIMARY KEY (id))", ""},
		{"CREATE TABLE n (id int REFERENCES t (id))", ""},
		{"CREATE TABLE n (id int, CHECK (id > 0) NOT VALID NO INHERIT)", ""},
		{"CREATE TABLE n (id int, FOREIGN KEY (id) REFERENCES t NOT VALID DEFERRABLE)", ""},
		{"CREATE TABLE n (id int REFERENCES t DEFERRABLE INITIALLY DEFERRED)", ""},
		{"ALTER TABLE t ALTER CONSTRAINT c DEFERRABLE INITIALLY DEFERRED", ""},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			_, err := Parse(tt.sql)
			if tt.message == "" {
				if err != nil {
					t.Fatalf("Parse: %v", err)
				}
				return
			}
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("Parse error = %v, want %q", err, tt.message)
			}
			if pe.Message != tt.message {
				t.Errorf("message = %q, want %q", pe.Message, tt.message)
			}
		})
	}
}
