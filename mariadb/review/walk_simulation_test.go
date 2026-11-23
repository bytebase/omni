package review

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/simulate"
	mariadbcatalog "github.com/bytebase/omni/mariadb/catalog"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/review/walk/walktest"
)

// TestWalkSimulation holds MariaDB's WALK_THROUGH to mariadb/catalog's
// simulate.tsv offline: each fixture is built once in the in-memory catalog,
// and every row is checked with no server started.
func TestWalkSimulation(t *testing.T) {
	same := (&mariadbEngine{}).Same
	res := walktest.CheckSimulation(t, walktest.Simulation{
		Path:     filepath.Join("..", "catalog", "testdata", "simulate.tsv"),
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
	const minRecallPct = 90.0
	if res.Refused > 0 {
		pct := 100 * float64(res.Found) / float64(res.Refused)
		if pct < minRecallPct {
			t.Errorf("WALK_THROUGH refusal recall = %.1f%% (%d/%d), want >= %.1f%%", pct, res.Found, res.Refused, minRecallPct)
		}
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
		Schema:           mariadbCatalogToMetadata(c, simDB),
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

func mariadbCatalogToMetadata(c *mariadbcatalog.Catalog, dbName string) *metadata.DatabaseSchemaMetadata {
	out := &metadata.DatabaseSchemaMetadata{Name: dbName}
	db := c.GetDatabase(dbName)
	if db == nil {
		return out
	}
	sm := &metadata.SchemaMetadata{}
	tblKeys := make([]string, 0, len(db.Tables))
	for k := range db.Tables {
		tblKeys = append(tblKeys, k)
	}
	slices.Sort(tblKeys)
	tblByLower := make(map[string]*metadata.TableMetadata, len(db.Tables))
	for _, k := range tblKeys {
		t := db.Tables[k]
		if t.Temporary {
			continue
		}
		tm := &metadata.TableMetadata{
			Name:   t.Name,
			Engine: t.Engine,
		}
		if t.SystemVersioned {
			tm.CreateOptions = "WITH SYSTEM VERSIONING"
		}
		if t.Partitioning != nil {
			for _, p := range t.Partitioning.Partitions {
				tm.Partitions = append(tm.Partitions, &metadata.TablePartitionMetadata{
					Name:       p.Name,
					Expression: t.Partitioning.Expr,
				})
			}
			if len(tm.Partitions) == 0 {
				n := max(t.Partitioning.NumParts, 1)
				for i := range n {
					tm.Partitions = append(tm.Partitions, &metadata.TablePartitionMetadata{
						Name:       fmt.Sprintf("p%d", i),
						Expression: t.Partitioning.Expr,
					})
				}
			}
		}
		for _, col := range t.VisibleColumns() {
			tm.Columns = append(tm.Columns, &metadata.ColumnMetadata{Name: col.Name})
		}
		for _, idx := range t.Indexes {
			name := idx.Name
			if idx.Primary {
				name = "PRIMARY"
			}
			var exprs []string
			var keyLengths []int64
			for _, ic := range idx.Columns {
				if ic == nil {
					continue
				}
				if ic.Expr != "" {
					exprs = append(exprs, ic.Expr)
				} else {
					exprs = append(exprs, ic.Name)
				}
				if ic.Length > 0 {
					keyLengths = append(keyLengths, int64(ic.Length))
				} else {
					keyLengths = append(keyLengths, -1)
				}
			}
			tm.Indexes = append(tm.Indexes, &metadata.IndexMetadata{
				Name:        name,
				Expressions: exprs,
				KeyLength:   keyLengths,
				Type:        idx.IndexType,
				Unique:      idx.Unique,
				Primary:     idx.Primary,
			})
		}
		for _, con := range t.Constraints {
			if con.Type == mariadbcatalog.ConForeignKey {
				tm.ForeignKeys = append(tm.ForeignKeys, &metadata.ForeignKeyMetadata{
					Name:            con.Name,
					Columns:         append([]string(nil), con.Columns...),
					ReferencedTable: con.RefTable,
				})
			}
		}
		sm.Tables = append(sm.Tables, tm)
		tblByLower[strings.ToLower(t.Name)] = tm
	}
	for _, tr := range db.Triggers {
		if tm := tblByLower[strings.ToLower(tr.Table)]; tm != nil {
			tm.Triggers = append(tm.Triggers, &metadata.TriggerMetadata{Name: tr.Name})
		}
	}
	for _, v := range db.Views {
		sm.Views = append(sm.Views, &metadata.ViewMetadata{Name: v.Name})
	}
	for _, q := range db.Sequences {
		sm.Sequences = append(sm.Sequences, &metadata.SequenceMetadata{Name: q.Name})
	}
	for _, f := range db.Functions {
		sm.Functions = append(sm.Functions, &metadata.FunctionMetadata{Name: f.Name})
	}
	for _, p := range db.Procedures {
		sm.Procedures = append(sm.Procedures, &metadata.ProcedureMetadata{Name: p.Name})
	}
	for _, ev := range db.Events {
		sm.Events = append(sm.Events, &metadata.EventMetadata{Name: ev.Name})
	}
	out.Schemas = []*metadata.SchemaMetadata{sm}
	return out
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
