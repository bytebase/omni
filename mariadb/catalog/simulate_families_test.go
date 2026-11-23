package catalog

import (
	"slices"

	"github.com/bytebase/omni/enginetest/simulate"
)

// familyRows are hand-written DDL rows on fixtures that name their objects as
// an application would, for the cases compat.tsv's statements reach only with
// one name for every object. A label is a word and a colon, so simulate keeps
// these rows after the generated ones.
//
// features@check holds a column CHECK (guests), a table CHECK on two columns
// (chk_dates), a table CHECK on one (chk_note), and one that reads no column
// (chk_always). MariaDB parses no DROP
// CHECK, so the rows drop a constraint with DROP CONSTRAINT.
var familyRows = families(map[string][]familyRow{
	"features@check": {
		{"check:drop", "ALTER TABLE bookings DROP COLUMN starts, DROP COLUMN ends"},
		{"check:drop", "ALTER TABLE bookings DROP COLUMN ends, DROP COLUMN starts"},
		{"check:drop", "ALTER TABLE bookings DROP COLUMN starts"},
		{"check:drop", "ALTER TABLE bookings DROP COLUMN ends"},
		{"check:drop", "ALTER TABLE bookings DROP COLUMN guests"},
		{"check:drop", "ALTER TABLE bookings DROP COLUMN note"},
		{"check:drop", "ALTER TABLE bookings DROP COLUMN guests, DROP COLUMN starts"},
		{"check:drop", "ALTER TABLE bookings DROP CONSTRAINT chk_dates, DROP COLUMN starts, DROP COLUMN ends"},
		{"check:drop", "ALTER TABLE bookings DROP COLUMN starts, DROP COLUMN ends, DROP CONSTRAINT chk_dates"},
		{"check:drop", "ALTER TABLE bookings DROP CONSTRAINT IF EXISTS chk_dates, DROP COLUMN starts, DROP COLUMN ends"},
		{"check:drop", "ALTER TABLE bookings DROP CONSTRAINT chk_dates, DROP COLUMN starts"},
		{"check:drop", "ALTER TABLE bookings DROP CHECK chk_dates, DROP COLUMN starts"},
		{"check:drop", "ALTER TABLE bookings RENAME COLUMN starts TO begins, DROP COLUMN ends"},
		{"check:drop", "ALTER TABLE bookings DROP COLUMN starts, ADD COLUMN starts DATE"},
		{"check:drop", "ALTER TABLE rooms DROP COLUMN beds"},
		{"check:drop", "ALTER TABLE rooms ADD COLUMN floor INT"},
	},
})

type familyRow struct{ label, statement string }

// families flattens the rows by fixture, in fixture name order.
func families(byFixture map[string][]familyRow) []simulate.Generated {
	var out []simulate.Generated
	for _, fx := range slices.Sorted(func(yield func(string) bool) {
		for fx := range byFixture {
			if !yield(fx) {
				return
			}
		}
	}) {
		for _, r := range byFixture[fx] {
			out = append(out, simulate.Generated{Fixture: fx, Rule: r.label, Statement: r.statement})
		}
	}
	return out
}
