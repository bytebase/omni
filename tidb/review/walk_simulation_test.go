package review

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/review/walk/walktest"
	tidbcatalog "github.com/bytebase/omni/tidb/catalog"
)

// TestWalkSimulation holds TiDB's WALK_THROUGH to tidb/catalog's
// simulate.tsv, sequences/simulate.tsv and queries/simulate.tsv offline: each fixture is built once
// in the in-memory catalog, and every row is checked with no server started.
func TestWalkSimulation(t *testing.T) {
	same := (&tidbEngine{}).Same
	for _, suite := range []struct {
		name      string
		path      string
		minRecall float64
	}{
		// 30 of the ddl refusals the walk misses are ALTER TABLE ADD/DROP
		// STATS_EXTENDED, 1105 while tidb_enable_extended_stats is off: a
		// disabled feature, which the walk does not model (createStatistics).
		{"ddl", filepath.Join("..", "catalog", "testdata", "simulate.tsv"), 97.9},
		{"sequences", filepath.Join("..", "catalog", "testdata", "sequences", "simulate.tsv"), 100.0},
		{"queries", filepath.Join("..", "catalog", "testdata", "queries", "simulate.tsv"), 96.5},
	} {
		t.Run(suite.name, func(t *testing.T) {
			res := walktest.CheckSimulation(t, walktest.Simulation{
				Path:     suite.path,
				Fixtures: filepath.Join("..", "catalog", "testdata", "fixtures"),
				Setup:    tidbSimFixture,
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
			if res.Refused > 0 {
				pct := 100 * float64(res.Found) / float64(res.Refused)
				if pct < suite.minRecall {
					t.Errorf("%s recall = %.1f%% (%d/%d), want >= %.1f%%",
						suite.name, pct, res.Found, res.Refused, suite.minRecall)
				}
			}
		})
	}
}

func tidbSimFixture(t testing.TB, fx simulate.Fixture) walktest.SimFixture {
	t.Helper()
	const simDB = "simulation"
	c, err := tidbcatalog.LoadSQL("CREATE DATABASE " + simDB + ";\nUSE " + simDB + ";\n" + fx.SQL)
	if err != nil {
		t.Fatalf("tidb catalog failed to load fixture %s: %v", fx.Name, err)
	}
	target := review.Target{
		Schema:              c.DatabaseMetadata(simDB),
		LowerCaseTableNames: 2,
		ForeignKeyChecks:    review.SwitchOn,
	}
	objects := tidbCatalogInventory(c, simDB)
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

func tidbCatalogInventory(c *tidbcatalog.Catalog, dbName string) []walktest.Object {
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
			if con.Type == tidbcatalog.ConForeignKey {
				objects = append(objects, walktest.Object{
					Key: walk.Key{Space: spaceTableForeignKey, Parent: walk.Path(db.Name, tbl.Name), Name: con.Name},
				})
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
	return objects
}
