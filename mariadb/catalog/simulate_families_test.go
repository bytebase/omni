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
		{"check:drop", "ALTER TABLE bookings RENAME COLUMN guests TO people"},
		{"check:drop", "ALTER TABLE bookings RENAME COLUMN note TO memo"},
		{"check:drop", "ALTER TABLE bookings CHANGE COLUMN starts begins DATE"},
		{"check:drop", "ALTER TABLE bookings RENAME COLUMN starts TO begins, DROP COLUMN note"},
		{"check:drop", "ALTER TABLE rooms DROP COLUMN beds"},
		{"check:drop", "ALTER TABLE rooms ADD COLUMN floor INT"},
	},
	// IF [NOT] EXISTS over names the fixture holds, which the generated rows,
	// whose statements fail first (SELECT a with no table), do not reach.
	"objects@view": {
		{"exists:view", "CREATE VIEW IF NOT EXISTS a AS SELECT 1"},
		{"exists:view", "CREATE VIEW IF NOT EXISTS t AS SELECT 1"},
		{"exists:view", "CREATE VIEW IF NOT EXISTS v AS SELECT 1"},
		{"exists:view", "CREATE VIEW IF NOT EXISTS a AS SELECT zzz FROM t"},
		{"exists:rename", "RENAME TABLE IF EXISTS a TO v2, m TO m2"},
		{"exists:rename", "RENAME TABLE IF EXISTS m TO m2, a TO t"},
		{"exists:column", "ALTER TABLE t CHANGE IF EXISTS x y INT, MODIFY IF EXISTS a BIGINT"},
		{"exists:column", "ALTER TABLE t RENAME COLUMN IF EXISTS x TO y, ALTER COLUMN IF EXISTS d SET DEFAULT 1"},
		{"exists:column", "ALTER TABLE t RENAME KEY IF EXISTS x TO y, RENAME KEY IF EXISTS idx_d TO idx_e"},
		{"exists:column", "ALTER TABLE t ALTER COLUMN IF EXISTS x SET DEFAULT (zzz), ADD COLUMN e INT"},
		{"exists:column", "ALTER TABLE t CHANGE IF EXISTS x y INT AS (zzz), ADD COLUMN e INT"},
	},
	"objects@partitioned": {
		{"exists:partition", "ALTER TABLE a ADD PARTITION IF NOT EXISTS (PARTITION p2 VALUES LESS THAN (300))"},
		{"exists:partition", "ALTER TABLE a ADD PARTITION IF NOT EXISTS (PARTITION a VALUES LESS THAN (300))"},
		{"exists:partition", "ALTER TABLE a ADD PARTITION IF NOT EXISTS (PARTITION a VALUES LESS THAN (200))"},
		{"exists:partition", "ALTER TABLE a ADD PARTITION IF NOT EXISTS (PARTITION p0 VALUES LESS THAN (150))"},
		{"exists:partition", "ALTER TABLE a ADD PARTITION IF NOT EXISTS (PARTITION a VALUES LESS THAN (300), PARTITION p3 VALUES LESS THAN (400))"},
		{"exists:partition", "ALTER TABLE a ADD PARTITION IF NOT EXISTS (PARTITION p3 VALUES LESS THAN (300), PARTITION a VALUES LESS THAN (400))"},
		{"exists:partition", "ALTER TABLE a DROP PARTITION IF EXISTS zzz"},
		{"exists:partition", "ALTER TABLE a DROP PARTITION IF EXISTS zzz, p0"},
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
