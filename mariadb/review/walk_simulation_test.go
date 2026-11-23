package review

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/simulate"
	mariadbcatalog "github.com/bytebase/omni/mariadb/catalog"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/review/walk/walktest"
)

// TestWalkSimulation holds MariaDB's WALK_THROUGH to mariadb/catalog's
// simulate.tsv and queries/simulate.tsv offline: each fixture is built once
// in the in-memory catalog, and every row is checked with no server started.
func TestWalkSimulation(t *testing.T) {
	same := (&mariadbEngine{}).Same
	for _, suite := range []struct {
		name     string
		path     string
		minPct   float64
		minFound int
	}{
		// The partitioned@, features@ and schemas@ rows and the alter:algorithm
		// and event: ones hold 142 refusals the walk finds 100 of: it does not
		// model partition bounds and names, InnoDB's foreign key and row-size
		// refusals, ALGORITHM, or invisible and spatial column rules. The rows
		// before them keep 4,175/4,340 (96.2%).
		{"ddl", filepath.Join("..", "catalog", "testdata", "simulate.tsv"), 95.0, 4275},
		{"queries", filepath.Join("..", "catalog", "testdata", "queries", "simulate.tsv"), 68.0, 135},
	} {
		t.Run(suite.name, func(t *testing.T) {
			res := walktest.CheckSimulation(t, walktest.Simulation{
				Path:     suite.path,
				Fixtures: filepath.Join("..", "catalog", "testdata", "fixtures"),
				Setup:    mariadbSimFixture,
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

func mariadbSimFixture(t testing.TB, fx simulate.Fixture) walktest.SimFixture {
	t.Helper()
	const simDB = "omni_simulation"
	c, err := mariadbcatalog.LoadSQL("CREATE DATABASE " + simDB + ";\nUSE " + simDB + ";\n" + fx.SQL)
	if err != nil {
		t.Fatalf("mariadb catalog failed to load fixture %s: %v", fx.Name, err)
	}
	target := review.Target{
		Schema:           c.DatabaseMetadata(simDB),
		ForeignKeyChecks: review.SwitchOn,
	}
	objects := mariadbCatalogInventory(c, simDB)
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

func mariadbCatalogInventory(c *mariadbcatalog.Catalog, dbName string) []walktest.Object {
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
		for _, col := range tbl.VisibleColumns() {
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
			if con.Type == mariadbcatalog.ConForeignKey {
				objects = append(objects,
					walktest.Object{Key: walk.Key{Space: spaceForeignKey, Parent: walk.Path(db.Name), Name: con.Name}},
					walktest.Object{Key: walk.Key{Space: spaceTableForeignKey, Parent: walk.Path(db.Name, tbl.Name), Name: con.Name}},
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
	for _, q := range db.Sequences {
		objects = append(objects, walktest.Object{
			Key:  walk.Key{Space: spaceRelation, Parent: walk.Path(db.Name), Name: q.Name},
			Kind: relSequence,
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
