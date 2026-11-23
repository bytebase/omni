package catalog

import (
	"errors"
	"strings"
	"testing"
)

// TestDropColumnReadByCheck pins what testdata/simulate.tsv's check:drop rows
// hold beside the codes: the SQLSTATE and text OceanBase 4.3.5.6 answers on
// features@check, and which checks the table keeps.
func TestDropColumnReadByCheck(t *testing.T) {
	const fixture = "CREATE DATABASE d; USE d;\n" +
		"CREATE TABLE bookings (id INT PRIMARY KEY, starts DATE, ends DATE, guests INT, nights INT, note VARCHAR(100), " +
		"CONSTRAINT chk_dates CHECK (starts <= ends), CONSTRAINT chk_nights CHECK (nights > 0), CONSTRAINT chk_always CHECK (1 = 1), " +
		"CONSTRAINT chk_text CHECK (note <> 'starts'));"
	cases := []struct {
		stmt   string
		msg    string // the 3959 OceanBase answers, or empty
		checks string // the checks bookings keeps
	}{
		{"ALTER TABLE bookings DROP COLUMN ENDS, DROP COLUMN starts", "Check constraint 'chk_dates' uses column 'ends', hence column cannot be dropped or renamed.", ""},
		{"ALTER TABLE bookings RENAME COLUMN Nights TO stays", "Check constraint 'chk_nights' uses column 'Nights', hence column cannot be dropped or renamed.", ""},
		{"ALTER TABLE bookings RENAME COLUMN starts TO begins, DROP COLUMN ends", "Check constraint 'chk_dates' uses column 'starts', hence column cannot be dropped or renamed.", ""},
		{"ALTER TABLE bookings DROP COLUMN starts, DROP COLUMN ends, DROP CHECK chk_dates", "", "chk_nights,chk_text"},
		{"ALTER TABLE bookings DROP COLUMN nights", "", "chk_dates,chk_text"},
		{"ALTER TABLE bookings DROP CHECK (chk_dates), DROP COLUMN starts", "", "chk_nights,chk_text"},
		{"ALTER TABLE bookings DROP CONSTRAINT (chk_nights, chk_dates), DROP COLUMN starts", "", "chk_text"},
		{"ALTER TABLE bookings ADD COLUMN floor INT", "", "chk_dates,chk_nights,chk_always,chk_text"},
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
			if tc.msg != "" {
				var e *Error
				if !errors.As(results[0].Error, &e) || e.Code != 3959 || e.SQLState != "HY000" || e.Message != tc.msg {
					t.Fatalf("got %v, want 3959 (HY000) %s", results[0].Error, tc.msg)
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
	t.Run("CHANGE COLUMN follows the rename into the check", func(t *testing.T) {
		c, err := LoadSQL(fixture + "\nALTER TABLE bookings CHANGE COLUMN starts begins DATE;")
		if err != nil {
			t.Fatal(err)
		}
		if got := c.ShowCreateTable("d", "bookings"); !strings.Contains(got, "CHECK ((`begins` <= `ends`))") || !strings.Contains(got, "'starts'") {
			t.Fatalf("chk_dates does not follow the rename, or the string moved:\n%s", got)
		}
	})
}
