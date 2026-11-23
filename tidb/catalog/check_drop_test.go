package catalog

import (
	"errors"
	"strings"
	"testing"
)

// TestDropColumnReadByCheck pins what testdata/simulate.tsv's check:drop rows
// hold beside the codes: the SQLSTATE and text TiDB v8.5.5 answers on
// features@check, and which checks the table keeps.
func TestDropColumnReadByCheck(t *testing.T) {
	const fixture = "SET GLOBAL tidb_enable_check_constraint = ON; CREATE DATABASE d; USE d;\n" +
		"CREATE TABLE bookings (id INT PRIMARY KEY, starts DATE, ends DATE, guests INT CHECK (guests > 0), nights INT, note VARCHAR(100), " +
		"CONSTRAINT chk_dates CHECK (starts <= ends), CONSTRAINT chk_nights CHECK (nights > 0), CONSTRAINT chk_always CHECK (1 = 1), " +
		"CONSTRAINT chk_text CHECK (note <> 'starts'));"
	cases := []struct {
		stmt   string
		code   int
		msg    string
		checks string // the checks bookings keeps
	}{
		{"ALTER TABLE bookings DROP COLUMN Ends, DROP COLUMN starts", 3959, "Check constraint 'chk_dates' uses column 'Ends', hence column cannot be dropped or renamed.", ""},
		{"ALTER TABLE bookings RENAME COLUMN NIGHTS TO nights", 3959, "Check constraint 'chk_nights' uses column 'nights', hence column cannot be dropped or renamed.", ""},
		{"ALTER TABLE bookings DROP CHECK chk_dates, DROP COLUMN starts", 8200, "Unsupported multi schema change for drop check constraint", ""},
		{"ALTER TABLE bookings DROP CHECK chk_dates", 0, "", "chk_nights,chk_always,chk_text,bookings_chk_1"},
		{"ALTER TABLE bookings DROP COLUMN nights", 0, "", "chk_dates,chk_always,chk_text,bookings_chk_1"},
		{"ALTER TABLE bookings CHANGE COLUMN nights stays INT, DROP COLUMN guests", 0, "", "chk_dates,chk_always,chk_text"},
		{"ALTER TABLE bookings MODIFY COLUMN nights BIGINT", 0, "", "chk_dates,chk_nights,chk_always,chk_text,bookings_chk_1"},
	}
	for _, tc := range cases {
		t.Run(tc.stmt, func(t *testing.T) {
			c, err := LoadSQL(fixture)
			if err != nil {
				t.Fatal(err)
			}
			results, err := c.Exec(tc.stmt, nil)
			if err != nil || len(results) != 1 {
				t.Fatalf("%v %v", results, err)
			}
			if tc.code != 0 {
				var e *Error
				if !errors.As(results[0].Error, &e) || e.Code != tc.code || e.SQLState != "HY000" || e.Message != tc.msg {
					t.Fatalf("got %v, want %d (HY000) %s", results[0].Error, tc.code, tc.msg)
				}
				return
			}
			if results[0].Error != nil {
				t.Fatal(results[0].Error)
			}
			var checks []string
			for _, con := range c.GetDatabase("d").GetTable("bookings").Constraints {
				if con.Type == ConCheck {
					checks = append(checks, con.Name)
				}
			}
			if got := strings.Join(checks, ","); got != tc.checks {
				t.Fatalf("checks %s, want %s", got, tc.checks)
			}
		})
	}
}
