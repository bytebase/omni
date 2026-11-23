package catalog

import (
	"errors"
	"strings"
	"testing"
)

// TestDropColumnReadByCheck pins what testdata/simulate.tsv's check:drop rows
// hold beside the codes: the SQLSTATE and text of the 1054, and which checks
// the table keeps, as mariadb:11.8.8 answers on features@check.
func TestDropColumnReadByCheck(t *testing.T) {
	const fixture = "CREATE DATABASE d; USE d;\n" +
		"CREATE TABLE bookings (id INT PRIMARY KEY, starts DATE, ends DATE, guests INT CHECK (guests > 0), note VARCHAR(100), " +
		"CONSTRAINT chk_dates CHECK (starts <= ends), CONSTRAINT chk_note CHECK (note <> 'starts'), CONSTRAINT chk_always CHECK (1 = 1));"
	cases := []struct {
		stmt    string
		message string   // the 1054 MariaDB answers, or empty
		checks  []string // the checks bookings keeps
	}{
		{"ALTER TABLE bookings DROP COLUMN starts", "Unknown column 'starts' in 'CHECK'", nil},
		{"ALTER TABLE bookings DROP COLUMN ends, DROP COLUMN guests", "Unknown column 'ends' in 'CHECK'", nil},
		{"ALTER TABLE bookings DROP COLUMN starts, ADD COLUMN starts DATE", "Unknown column 'starts' in 'CHECK'", nil},
		{"ALTER TABLE bookings DROP COLUMN starts, DROP COLUMN ends", "", []string{"chk_note"}},
		{"ALTER TABLE bookings DROP COLUMN starts, DROP COLUMN ends, DROP CONSTRAINT chk_dates", "", []string{"chk_note"}},
		{"ALTER TABLE bookings DROP COLUMN starts, DROP CONSTRAINT chk_dates", "", []string{"chk_note"}},
		{"ALTER TABLE bookings DROP COLUMN note", "", []string{"chk_dates"}},
		{"ALTER TABLE bookings DROP COLUMN IF EXISTS nothing", "", []string{"chk_dates", "chk_note", "chk_always"}},
		{"ALTER TABLE bookings ADD COLUMN floor INT", "", []string{"chk_dates", "chk_note", "chk_always"}},
	}
	for _, tc := range cases {
		t.Run(tc.stmt, func(t *testing.T) {
			c, err := LoadSQL(fixture)
			if err != nil {
				t.Fatal(err)
			}
			results, err := c.Exec(tc.stmt, nil)
			if err != nil {
				t.Fatal(err)
			}
			var got error
			for _, r := range results {
				if r.Error != nil {
					got = r.Error
				}
			}
			if tc.message != "" {
				var e *Error
				if !errors.As(got, &e) || e.Code != 1054 || e.SQLState != "42S22" || e.Message != tc.message {
					t.Fatalf("got %v, want 1054 (42S22) %s", got, tc.message)
				}
				return
			}
			if got != nil {
				t.Fatal(got)
			}
			var checks []string
			for _, con := range c.GetDatabase("d").GetTable("bookings").Constraints {
				if con.Type == ConCheck && !isColumnCheck(con) {
					checks = append(checks, con.Name)
				}
			}
			if strings.Join(checks, ",") != strings.Join(tc.checks, ",") {
				t.Fatalf("checks %v, want %v", checks, tc.checks)
			}
		})
	}
}
