package review

import (
	"context"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/simulate"
	obcatalog "github.com/bytebase/omni/oceanbase/catalog"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/review/walk/walktest"
)

// TestWalkSimulation holds OceanBase's WALK_THROUGH to oceanbase/catalog's
// simulate.tsv, programs/simulate.tsv, and queries/simulate.tsv offline: each
// fixture is built once in the in-memory catalog, and every row is checked
// with no server started.
func TestWalkSimulation(t *testing.T) {
	same := (&oceanBaseEngine{}).Same
	for _, suite := range []struct {
		name     string
		path     string
		minFound int
	}{
		// What the walk misses a synced schema can't show: a primary key
		// beside a secondary index the sync may take for it (1068),
		// tablespaces and tablegroups (3510, 5151), column groups (4185),
		// other databases (1049, 1146, 1360, 1305), events, which
		// OceanBase lists for dropped databases too (1537), and rows omni's
		// parser refuses first. Nor do the rows settle when OceanBase takes
		// STORAGE_FORMAT_VERSION (4179), DYNAMIC_PARTITION_POLICY,
		// repartitioning or REORGANIZE PARTITION of a partitioned table, or
		// a vector ORDER BY ... APPROX (1235).
		{"ddl", filepath.Join("..", "catalog", "testdata", "simulate.tsv"), 3599},
		{"programs", filepath.Join("..", "catalog", "testdata", "programs", "simulate.tsv"), 280},
		{"queries", filepath.Join("..", "catalog", "testdata", "queries", "simulate.tsv"), 4069},
		// A sequence is not synced, so the walk finds only what OceanBase
		// refuses of the statement itself, never a sequence that exists or not
		// (4297, 4298), a database that does not (1049), or what ALTER SEQUENCE
		// makes of the sequence's own options.
		{"sequences", filepath.Join("..", "catalog", "testdata", "sequences", "simulate.tsv"), 280},
	} {
		t.Run(suite.name, func(t *testing.T) {
			res := walktest.CheckSimulation(t, walktest.Simulation{
				Path:     suite.path,
				Fixtures: filepath.Join("..", "catalog", "testdata", "fixtures"),
				Setup:    oceanBaseSimFixture,
				Codes: func(message, stmt string) []string {
					var codes []string
					for _, c := range walkCodes {
						if c.match(message, stmt) {
							codes = c.codes
							break
						}
					}
					return append(slices.Clone(codes), oceanBaseCodes(message, stmt)...)
				},
				Same: same,
			})
			if res.Found < suite.minFound {
				t.Errorf("%s WALK_THROUGH recall regression: found %d/%d refusals, want >= %d", suite.name, res.Found, res.Refused, suite.minFound)
			}
		})
	}
}

func oceanBaseSimFixture(t testing.TB, fx simulate.Fixture) walktest.SimFixture {
	t.Helper()
	const simDB = "simulation"
	c, err := obcatalog.LoadSQL("CREATE DATABASE " + simDB + " DEFAULT CHARSET=utf8mb4;\nUSE " + simDB + ";\n" + fx.SQL)
	if err != nil {
		t.Fatalf("oceanbase catalog failed to load fixture %s: %v", fx.Name, err)
	}
	target := review.Target{
		Schema:              c.DatabaseMetadata(simDB),
		LowerCaseTableNames: 1,
		ForeignKeyChecks:    review.SwitchOn,
	}
	objects := oceanBaseCatalogInventory(c, simDB)
	return walktest.SimFixture{
		Objects: objects,
		Report: func(t testing.TB, stmt string) *walktest.Finding {
			t.Helper()
			if strings.TrimSpace(stmt) == "" {
				return nil
			}
			result, err := ReviewOceanBase(context.Background(), stmt, review.Options{Rules: []review.Rule{review.WalkThrough}}, []review.Target{target})
			if err != nil {
				t.Fatalf("ReviewOceanBase: %v", err)
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
			v, err := walkThrough(context.Background(), stmt, stmts, target, dialect{oceanBase: true})
			if err != nil || v == nil {
				t.Fatalf("no violation for a statement ReviewOceanBase reports on (%v):\n%s", err, stmt)
			}
			return v
		},
	}
}

func oceanBaseCatalogInventory(c *obcatalog.Catalog, dbName string) []walktest.Object {
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
			if con.Type == obcatalog.ConCheck {
				objects = append(objects, walktest.Object{
					Key: walk.Key{Space: spaceTableCheck, Parent: walk.Path(db.Name, tbl.Name), Name: con.Name},
				})
			}
			if con.Type == obcatalog.ConForeignKey {
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
