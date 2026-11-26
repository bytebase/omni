package review

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/simulate"
	mysqlcatalog "github.com/bytebase/omni/mysql/catalog"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/review/walk/walktest"
)

// TestWalkSimulation holds MySQL's WALK_THROUGH to mysql/catalog's
// simulate.tsv, programs/simulate.tsv, queries/simulate.tsv and their
// mysql84/ copies offline: each fixture is built once in the in-memory
// catalog, and every row is checked with no server started.
func TestWalkSimulation(t *testing.T) {
	same := (&mysqlEngine{fold: false}).Same
	for _, suite := range []struct {
		name     string
		path     string
		minPct   float64
		minFound int
	}{
		// The 75 ddl refusals the walk misses: 51 rows the parser refuses
		// (49 give COLLATE utf8mb4, an unknown collation, 1273 while
		// parsing; 1504, 1064), which reach the review as a syntax finding.
		// 16 are ALTER TABLE ... TABLESPACE t and IMPORT TABLESPACE, on
		// tablespace state a schema doesn't show (3510, 1813). 6 are on
		// another database or a server object (1049, 1050, 3730; 1128 a
		// parser plugin, 3548 an SRID). 2 the walk doesn't judge: the row
		// size (1118) and what an engine supports (1178).
		{"ddl", filepath.Join("..", "catalog", "testdata", "simulate.tsv"), 97.6, 3081},
		// The 36 refusals the walk misses, none of which a schema shows: 9
		// drop a function the sync can't rule out: unqualified, a loadable
		// one of the name would go, and qualified, it is another database's
		// (1305). 2 are 8.0-only: mysql:8.0.46 answers 1064 to DISABLE ON
		// REPLICA, which 8.4 adds. 25 give a routine's parameter, return or
		// variable type COLLATE utf8mb4, an unknown collation the parser
		// refuses (1273), which reaches the review as a syntax finding.
		{"programs", filepath.Join("..", "catalog", "testdata", "programs", "simulate.tsv"), 88.8, 287},
		// compat.tsv's queries in a view on one_table@column. The 13
		// refusals the walk misses: 4 are syntax findings (1273, 1064), and
		// 9 call a function a loadable one of the name or another database
		// could hold (1305, 1584).
		{"queries", filepath.Join("..", "catalog", "testdata", "queries", "simulate.tsv"), 94.4, 222},
		// The ddl rows on mysql:8.4.13, with six 8.4-only refusals (6125,
		// 1063) the walk doesn't model: it judges every server alike.
		{"mysql84", filepath.Join("..", "catalog", "testdata", "mysql84", "simulate.tsv"), 97.4, 3082},
		// The programs and queries rows on mysql:8.4.13, which refuses the
		// same rows as 8.0.
		{"mysql84_programs", filepath.Join("..", "catalog", "testdata", "mysql84", "programs", "simulate.tsv"), 89.4, 287},
		{"mysql84_queries", filepath.Join("..", "catalog", "testdata", "mysql84", "queries", "simulate.tsv"), 94.4, 222},
	} {
		t.Run(suite.name, func(t *testing.T) {
			res := walktest.CheckSimulation(t, walktest.Simulation{
				Path:     suite.path,
				Fixtures: filepath.Join("..", "catalog", "testdata", "fixtures"),
				Setup:    mysqlSimFixture,
				Codes: func(message, stmt string) []string {
					for _, c := range walkCodes {
						if c.match(message, stmt) {
							return c.codes
						}
					}
					return nil
				},
				Same: same,
			})
			pct := 0.0
			if res.Refused > 0 {
				pct = 100 * float64(res.Found) / float64(res.Refused)
			}
			if res.Found < suite.minFound || pct < suite.minPct {
				t.Errorf("%s recall regressed: got %d/%d (%.1f%%), want at least %d (%.1f%%)",
					suite.name, res.Found, res.Refused, pct, suite.minFound, suite.minPct)
			}
		})
	}
}

func mysqlSimFixture(t testing.TB, fx simulate.Fixture) walktest.SimFixture {
	t.Helper()
	const simDB = "simulation"
	c, err := mysqlcatalog.LoadSQL("CREATE DATABASE " + simDB + " DEFAULT CHARSET=utf8mb4;\nUSE " + simDB + ";\n" + fx.SQL)
	if err != nil {
		t.Fatalf("mysql catalog failed to load fixture %s: %v", fx.Name, err)
	}
	target := review.Target{
		Schema:              c.DatabaseMetadata(simDB),
		LowerCaseTableNames: 0,
		ForeignKeyChecks:    review.SwitchOn,
	}
	objects := mysqlCatalogInventory(c, simDB)
	return walktest.SimFixture{
		Objects: objects,
		Report: func(t testing.TB, stmt string) *walktest.Finding {
			t.Helper()
			if strings.TrimSpace(stmt) == "" {
				return nil
			}
			result, err := Review(context.Background(), stmt, review.Options{Rules: []review.Rule{review.WalkThrough}}, []review.Target{target})
			if err != nil {
				t.Fatalf("Review: %v", err)
			}
			var out *walktest.Finding
			for _, f := range result.Findings {
				if err := review.CheckCodes(f); err != nil {
					t.Error(err)
				}
				switch f.Rule {
				case review.Syntax:
					return nil
				case review.WalkThrough:
					if out != nil {
						t.Errorf("more than one finding: %q and %q", out.Message, f.Message)
					}
					out = &walktest.Finding{Statement: f.Statement, Message: f.Message}
				}
			}
			return out
		},
		Violation: func(t testing.TB, stmt string) *walk.Violation {
			t.Helper()
			stmts, syntax := parse(stmt, statementRanges(stmt))
			if syntax != nil {
				t.Fatalf("syntax: %s", syntax.Message)
			}
			v, err := walkThrough(context.Background(), stmt, stmts, target)
			if err != nil || v == nil {
				t.Fatalf("no violation for a statement Review reports on (%v):\n%s", err, stmt)
			}
			return v
		},
	}
}

func mysqlCatalogInventory(c *mysqlcatalog.Catalog, dbName string) []walktest.Object {
	db := c.GetDatabase(dbName)
	if db == nil {
		return nil
	}
	objects := []walktest.Object{{Key: walk.Key{Space: spaceDatabase, Name: db.Name}}}
	for _, tbl := range db.Tables {
		if tbl.Temporary {
			continue
		}
		objects = append(objects, walktest.Object{
			Key:  walk.Key{Space: spaceRelation, Parent: walk.Path(db.Name), Name: tbl.Name},
			Kind: relTable,
		})
		for _, col := range tbl.Columns {
			objects = append(objects, walktest.Object{
				Key: walk.Key{Space: spaceColumn, Parent: walk.Path(db.Name, tbl.Name), Name: col.Name},
			})
		}
		for _, idx := range tbl.Indexes {
			name := idx.Name
			if idx.Primary {
				name = "PRIMARY"
			}
			objects = append(objects, walktest.Object{
				Key: walk.Key{Space: spaceIndex, Parent: walk.Path(db.Name, tbl.Name), Name: name},
			})
		}
		for _, con := range tbl.Constraints {
			switch con.Type {
			case mysqlcatalog.ConForeignKey:
				objects = append(objects,
					walktest.Object{Key: walk.Key{Space: spaceForeignKey, Parent: walk.Path(db.Name), Name: con.Name}},
					walktest.Object{Key: walk.Key{Space: spaceTableForeignKey, Parent: walk.Path(db.Name, tbl.Name), Name: con.Name}},
				)
			case mysqlcatalog.ConCheck:
				objects = append(objects,
					walktest.Object{Key: walk.Key{Space: spaceCheckConstraint, Parent: walk.Path(db.Name, tbl.Name), Name: con.Name}},
				)
			}
		}
	}
	for _, v := range db.Views {
		objects = append(objects, walktest.Object{
			Key:  walk.Key{Space: spaceRelation, Parent: walk.Path(db.Name), Name: v.Name},
			Kind: relView,
		})
	}
	for _, tr := range db.Triggers {
		objects = append(objects, walktest.Object{Key: walk.Key{Space: spaceTrigger, Parent: walk.Path(db.Name), Name: tr.Name}})
	}
	for _, p := range db.Procedures {
		objects = append(objects, walktest.Object{Key: walk.Key{Space: spaceProcedure, Parent: walk.Path(db.Name), Name: p.Name}})
	}
	for _, f := range db.Functions {
		objects = append(objects, walktest.Object{Key: walk.Key{Space: spaceFunction, Parent: walk.Path(db.Name), Name: f.Name}})
	}
	for _, ev := range db.Events {
		objects = append(objects, walktest.Object{Key: walk.Key{Space: spaceEvent, Parent: walk.Path(db.Name), Name: ev.Name}})
	}
	return objects
}
