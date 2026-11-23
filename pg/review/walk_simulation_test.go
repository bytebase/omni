package review

import (
	"context"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/metadata"
	pgcatalog "github.com/bytebase/omni/pg/catalog"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/review/walk/walktest"
)

// TestWalkSimulation holds PostgreSQL's WALK_THROUGH to pg/catalog's
// simulate.tsv offline: each fixture is built once in the in-memory catalog,
// and every row is checked with no server started.
func TestWalkSimulation(t *testing.T) {
	res := walktest.CheckSimulation(t, walktest.Simulation{
		Path:     filepath.Join("..", "catalog", "testdata", "simulate.tsv"),
		Fixtures: filepath.Join("..", "catalog", "testdata", "fixtures"),
		Setup:    pgSimFixture,
		Codes:    func(msg, _ string) []string { return walkStates(msg) },
		Same:     func(claimed, stored walk.Key) bool { return claimed == stored },
	})
	const minRecall = 0.90
	if res.Refused > 0 {
		got := float64(res.Found) / float64(res.Refused)
		if got < minRecall {
			t.Errorf("WALK_THROUGH recall = %.1f%% (%d/%d), want >= %.1f%%", 100*got, res.Found, res.Refused, 100*minRecall)
		}
	}
}

func pgSimFixture(t testing.TB, fx simulate.Fixture) walktest.SimFixture {
	t.Helper()
	c, err := pgcatalog.LoadSQL(fx.SQL)
	if err != nil {
		t.Fatalf("pg catalog failed to load fixture %s: %v", fx.Name, err)
	}
	target := review.Target{
		Schema:      pgCatalogToMetadata(c),
		SessionUser: "postgres",
	}
	objects := pgCatalogInventory(c)
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
					// A parser gap tracked in compat.tsv yields no WALK_THROUGH finding.
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
			v := walkViolation(t, stmt, target)
			f := walkFinding(t, stmt, target)
			for _, prefix := range partitionedFindings {
				if f != nil && strings.HasPrefix(f.Message, prefix) {
					return partitionedClaim(v, strings.TrimSuffix(strings.TrimPrefix(f.Message, prefix), "`."))
				}
			}
			return v
		},
	}
}

func pgCatalogToMetadata(c *pgcatalog.Catalog) *metadata.DatabaseSchemaMetadata {
	db := &metadata.DatabaseSchemaMetadata{
		Name:       "simulation",
		SearchPath: `"$user", public`,
		Owner:      "postgres",
	}
	// Collect child partitions by parent OID so they appear under their parent
	// table's Partitions (matching Bytebase's PostgreSQL sync) instead of as
	// top-level tables.
	partChildren := make(map[uint32][]*pgcatalog.Relation)
	for _, s := range c.UserSchemas() {
		relNames := make([]string, 0, len(s.Relations))
		for name := range s.Relations {
			relNames = append(relNames, name)
		}
		slices.Sort(relNames)
		for _, name := range relNames {
			rel := s.Relations[name]
			if rel.PartitionOf != 0 {
				partChildren[rel.PartitionOf] = append(partChildren[rel.PartitionOf], rel)
			}
		}
	}
	indexMeta := func(rel *pgcatalog.Relation) []*metadata.IndexMetadata {
		var out []*metadata.IndexMetadata
		for _, idx := range c.IndexesOf(rel.OID) {
			only := ""
			if rel.RelKind == 'p' {
				only = "ONLY "
			}
			unique := ""
			if idx.IsUnique {
				unique = "UNIQUE "
			}
			schemaName := "public"
			if rel.Schema != nil && rel.Schema.Name != "" {
				schemaName = rel.Schema.Name
			}
			am := idx.AccessMethod
			if am == "" {
				am = "btree"
			}
			out = append(out, &metadata.IndexMetadata{
				Name:         idx.Name,
				Unique:       idx.IsUnique,
				Primary:      idx.IsPrimary,
				IsConstraint: idx.ConstraintOID != 0,
				Definition:   "CREATE " + unique + "INDEX " + ident(idx.Name) + " ON " + only + ident(schemaName) + "." + ident(rel.Name) + " USING " + ident(am) + " (a)",
			})
		}
		return out
	}
	var buildPartitions func(parentOID uint32) []*metadata.TablePartitionMetadata
	buildPartitions = func(parentOID uint32) []*metadata.TablePartitionMetadata {
		children := partChildren[parentOID]
		if len(children) == 0 {
			return nil
		}
		out := make([]*metadata.TablePartitionMetadata, 0, len(children))
		for _, child := range children {
			out = append(out, &metadata.TablePartitionMetadata{
				Name:          child.Name,
				Indexes:       indexMeta(child),
				Subpartitions: buildPartitions(child.OID),
			})
		}
		return out
	}
	for _, s := range c.UserSchemas() {
		sm := &metadata.SchemaMetadata{Name: s.Name}
		relNames := make([]string, 0, len(s.Relations))
		for name := range s.Relations {
			relNames = append(relNames, name)
		}
		slices.Sort(relNames)
		for _, name := range relNames {
			rel := s.Relations[name]
			if rel.PartitionOf != 0 {
				continue
			}
			switch rel.RelKind {
			case 'r', 'p':
				tbl := &metadata.TableMetadata{
					Name:       rel.Name,
					Indexes:    indexMeta(rel),
					Partitions: buildPartitions(rel.OID),
				}
				for _, col := range rel.Columns {
					tbl.Columns = append(tbl.Columns, &metadata.ColumnMetadata{Name: col.Name})
				}
				for _, con := range c.ConstraintsOf(rel.OID) {
					switch con.Type {
					case pgcatalog.ConstraintFK:
						tbl.ForeignKeys = append(tbl.ForeignKeys, &metadata.ForeignKeyMetadata{Name: con.Name})
					case pgcatalog.ConstraintCheck:
						tbl.CheckConstraints = append(tbl.CheckConstraints, &metadata.CheckConstraintMetadata{Name: con.Name})
					case pgcatalog.ConstraintExclude:
						tbl.ExcludeConstraints = append(tbl.ExcludeConstraints, &metadata.ExcludeConstraintMetadata{Name: con.Name})
					}
				}
				for _, tr := range c.TriggersOf(rel.OID) {
					tbl.Triggers = append(tbl.Triggers, &metadata.TriggerMetadata{Name: tr.Name})
				}
				sm.Tables = append(sm.Tables, tbl)
			case 'v':
				v := &metadata.ViewMetadata{Name: rel.Name}
				for _, col := range rel.Columns {
					v.Columns = append(v.Columns, &metadata.ColumnMetadata{Name: col.Name})
				}
				sm.Views = append(sm.Views, v)
			case 'm':
				mv := &metadata.MaterializedViewMetadata{Name: rel.Name}
				for _, idx := range c.IndexesOf(rel.OID) {
					mv.Indexes = append(mv.Indexes, &metadata.IndexMetadata{Name: idx.Name, Unique: idx.IsUnique})
				}
				sm.MaterializedViews = append(sm.MaterializedViews, mv)
			case 'f':
				ft := &metadata.ExternalTableMetadata{Name: rel.Name}
				for _, col := range rel.Columns {
					ft.Columns = append(ft.Columns, &metadata.ColumnMetadata{Name: col.Name})
				}
				sm.ExternalTables = append(sm.ExternalTables, ft)
			case 'c':
				sm.CompositeTypes = append(sm.CompositeTypes, &metadata.CompositeTypeMetadata{Name: rel.Name})
			}
		}
		seqNames := make([]string, 0, len(s.Sequences))
		for name := range s.Sequences {
			seqNames = append(seqNames, name)
		}
		slices.Sort(seqNames)
		for _, name := range seqNames {
			sm.Sequences = append(sm.Sequences, &metadata.SequenceMetadata{Name: name})
		}
		for _, typ := range c.QueryPgType(s.Name) {
			if typ.TypType == 'e' {
				sm.EnumTypes = append(sm.EnumTypes, &metadata.EnumTypeMetadata{Name: typ.TypName})
			}
		}
		for _, proc := range c.QueryPgProc(s.Name) {
			if builtinCalls[strings.ToLower(proc.ProName)] {
				continue
			}
			if proc.ProKind == 'p' {
				sm.Procedures = append(sm.Procedures, &metadata.ProcedureMetadata{Name: proc.ProName})
			} else {
				sm.Functions = append(sm.Functions, &metadata.FunctionMetadata{Name: proc.ProName})
			}
		}
		db.Schemas = append(db.Schemas, sm)
	}
	return db
}

func pgCatalogInventory(c *pgcatalog.Catalog) []walktest.Object {
	objects := []walktest.Object{
		{Key: schemaKey("pg_catalog")},
		{Key: schemaKey("information_schema")},
	}
	for _, s := range c.UserSchemas() {
		objects = append(objects, walktest.Object{Key: schemaKey(s.Name)})
		for _, rel := range s.Relations {
			k := relKey(s.Name, rel.Name)
			walkKind := relkinds[string(rel.RelKind)]
			objects = append(objects, walktest.Object{Key: k, Kind: walkKind})
			if rel.RelKind == 'p' {
				objects = append(objects, walktest.Object{Key: k, Kind: relPartitioned})
			}
			for _, col := range rel.Columns {
				objects = append(objects, walktest.Object{Key: columnKey(k, col.Name)})
			}
			for _, idx := range c.IndexesOf(rel.OID) {
				objects = append(objects, walktest.Object{Key: relKey(s.Name, idx.Name), Kind: relIndex})
			}
			for _, con := range c.ConstraintsOf(rel.OID) {
				objects = append(objects, walktest.Object{Key: constraintKey(k, con.Name)})
				if con.Type == pgcatalog.ConstraintPK {
					objects = append(objects, walktest.Object{Key: primaryKeyOf(k)})
				}
			}
			for _, tr := range c.TriggersOf(rel.OID) {
				objects = append(objects, walktest.Object{Key: trigKey(k, tr.Name)})
			}
		}
		for _, typ := range c.QueryPgType(s.Name) {
			typKind := typOther
			if typ.TypType == 'd' {
				typKind = typDomain
			}
			objects = append(objects, walktest.Object{Key: typeKey(s.Name, typ.TypName), Kind: typKind})
		}
		for _, seq := range s.Sequences {
			objects = append(objects, walktest.Object{Key: relKey(s.Name, seq.Name), Kind: relSequence})
		}
		for _, proc := range c.QueryPgProc(s.Name) {
			if proc.ProKind == 'p' {
				objects = append(objects, walktest.Object{Key: procKey(s.Name, proc.ProName)})
			} else {
				objects = append(objects, walktest.Object{Key: funcKey(s.Name, proc.ProName)})
			}
		}
	}
	return objects
}
