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
		// A partitioned table holds no foreign key, so DROP FOREIGN KEY IF
		// EXISTS is always dropped as a note before 1506 is checked.
		{"exists:foreign_key", "ALTER TABLE a DROP FOREIGN KEY IF EXISTS a, ADD COLUMN e INT"},
		{"exists:foreign_key", "ALTER TABLE a DROP FOREIGN KEY IF EXISTS a, DROP FOREIGN KEY zzz"},
		{"exists:foreign_key", "ALTER TABLE a DROP FOREIGN KEY IF EXISTS a, DROP COLUMN zzz"},
		{"exists:foreign_key", "ALTER TABLE a ADD FOREIGN KEY IF NOT EXISTS (a) REFERENCES a (b)"},
	},
	// ADD FOREIGN KEY IF NOT EXISTS looks for a foreign key named after the
	// name after FOREIGN KEY, the symbol or the first column; the fixture's is
	// CONSTRAINT a FOREIGN KEY (a).
	"objects@fk": {
		{"exists:foreign_key", "ALTER TABLE a ADD FOREIGN KEY IF NOT EXISTS (d) REFERENCES parent (b)"},
		{"exists:foreign_key", "ALTER TABLE a ADD FOREIGN KEY IF NOT EXISTS (A) REFERENCES parent (b)"},
		{"exists:foreign_key", "ALTER TABLE a ADD FOREIGN KEY IF NOT EXISTS (a, d) REFERENCES a (b, a)"},
		{"exists:foreign_key", "ALTER TABLE a ADD CONSTRAINT zz FOREIGN KEY IF NOT EXISTS (a) REFERENCES parent (b)"},
		{"exists:foreign_key", "ALTER TABLE a ADD FOREIGN KEY IF NOT EXISTS zz (a) REFERENCES parent (b)"},
		{"exists:foreign_key", "ALTER TABLE a ADD FOREIGN KEY IF NOT EXISTS (a) REFERENCES parent (b), ADD COLUMN e INT"},
		{"exists:foreign_key", "ALTER TABLE a ADD FOREIGN KEY IF NOT EXISTS (d) REFERENCES parent (b), ADD FOREIGN KEY IF NOT EXISTS (d) REFERENCES parent (b)"},
		{"exists:foreign_key", "ALTER TABLE a ADD FOREIGN KEY IF NOT EXISTS (d) REFERENCES parent (b), ADD CONSTRAINT d FOREIGN KEY IF NOT EXISTS (b) REFERENCES parent (b)"},
		{"exists:foreign_key", "ALTER TABLE a ADD FOREIGN KEY (d) REFERENCES parent (b), ADD FOREIGN KEY IF NOT EXISTS (d) REFERENCES parent (b)"},
		{"exists:foreign_key", "ALTER TABLE a ADD FOREIGN KEY IF NOT EXISTS (d) REFERENCES parent (b), ADD FOREIGN KEY (d) REFERENCES parent (b)"},
		{"exists:foreign_key", "ALTER TABLE a ADD FOREIGN KEY IF NOT EXISTS (b) REFERENCES parent (b), ADD FOREIGN KEY IF NOT EXISTS (d) REFERENCES parent (b)"},
		{"exists:foreign_key", "ALTER TABLE a ADD FOREIGN KEY (b) REFERENCES parent (b), ADD KEY IF NOT EXISTS (b)"},
		{"exists:foreign_key", "ALTER TABLE a ADD FOREIGN KEY (b) REFERENCES parent (b), ADD UNIQUE IF NOT EXISTS (b)"},
		{"exists:foreign_key", "ALTER TABLE a ADD CONSTRAINT zz FOREIGN KEY (b) REFERENCES parent (b), ADD KEY IF NOT EXISTS (b)"},
		{"exists:foreign_key", "ALTER TABLE a ADD CONSTRAINT zz FOREIGN KEY yy (b) REFERENCES parent (b), ADD KEY IF NOT EXISTS yy (b)"},
		{"exists:foreign_key", "ALTER TABLE a ADD CONSTRAINT zz FOREIGN KEY yy (b) REFERENCES parent (b), ADD KEY IF NOT EXISTS zz (b)"},
		// The fixture's KEY a came with its foreign key, and a second foreign
		// key over a replaces it with its own.
		{"fk:generated_key", "ALTER TABLE a ADD CONSTRAINT zz FOREIGN KEY (a) REFERENCES parent (b)"},
		{"fk:generated_key", "ALTER TABLE a ADD FOREIGN KEY (a) REFERENCES parent (b)"},
		{"fk:generated_key", "ALTER TABLE a ADD CONSTRAINT idx_d FOREIGN KEY (a) REFERENCES parent (b)"},
		{"fk:generated_key", "ALTER TABLE a ADD FOREIGN KEY idx_d (a) REFERENCES parent (b)"},
		{"fk:generated_key", "ALTER TABLE a ADD CONSTRAINT idx_d FOREIGN KEY (a) REFERENCES nosuch (b)"},
		{"fk:generated_key", "ALTER TABLE a ADD CONSTRAINT idx_d FOREIGN KEY (b, d) REFERENCES parent (b)"},
		{"fk:generated_key", "ALTER TABLE a ADD CONSTRAINT idx_d FOREIGN KEY (d) REFERENCES parent (b)"},
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
