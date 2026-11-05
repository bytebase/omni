package parser

import (
	"errors"
	"slices"
	"sort"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// setNameVerdict asks the server whether sql parses, by preparing it: these rows name real
// variables, and running SET GLOBAL on the shared server would change it for the other tests.
func setNameVerdict(t *testing.T, o *parserOracle, sql string) bool {
	t.Helper()
	stmt, err := o.db.Prepare(sql)
	if err == nil {
		stmt.Close()
		return true
	}
	var my *gomysql.MySQLError
	if !errors.As(err, &my) {
		t.Fatalf("no verdict for %q: %v", sql, err)
	}
	return my.Number != 1064
}

// TestSetNameFormsMariaDBOracle asserts omni agrees with a live MariaDB 11.8.8 on setNameCases
// (container test).
func TestSetNameFormsMariaDBOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, c := range setNameCases {
		t.Run(c.sql, func(t *testing.T) {
			if accepts := setNameVerdict(t, o, c.sql); accepts != c.accept {
				t.Fatalf("MariaDB accepts=%v, want %v for %q", accepts, c.accept, c.sql)
			}
			omni, err := setNameParses(c.sql)
			if omni != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, MariaDB accepts=%v for %q: %v", omni, c.accept, c.sql, err)
			}
		})
	}
}

// TestSetNameSweepMariaDBOracle holds every registered keyword to a live MariaDB 11.8.8 in each
// position of a SET name. Only setNameGap may disagree, and only by omni refusing, in the
// positions without a dot, where the word's category decides.
func TestSetNameSweepMariaDBOracle(t *testing.T) {
	o := startMariaDB(t)
	words := make([]string, 0, len(keywords))
	for word := range keywords {
		words = append(words, word)
	}
	sort.Strings(words)
	for _, position := range setNamePositions {
		for _, word := range words {
			sql := strings.ReplaceAll(position, "%s", word)
			accepts := setNameVerdict(t, o, sql)
			omni, err := setNameParses(sql)
			inGap := slices.Contains(setNameGap, word) && !strings.Contains(position, ".")
			switch {
			case inGap && (!accepts || omni):
				t.Errorf("%q is in setNameGap, but MariaDB accepts=%v and omni accepts=%v: drop it from the list", sql, accepts, omni)
			case !inGap && accepts != omni:
				t.Errorf("DIVERGENCE omni accepts=%v, MariaDB accepts=%v for %q: %v", omni, accepts, sql, err)
			}
		}
	}
}
