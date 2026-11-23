package review

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/metadata"
	mssqlcatalog "github.com/bytebase/omni/mssql/catalog"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/review/walk/walktest"
)

// walkCatalogSimulation holds WALK_THROUGH to one of mssql/catalog's
// simulate.tsv suites offline: each fixture is built once in the in-memory
// catalog, whose Metadata stands for Bytebase's sync, and every row is
// checked with no server started. A finding on an accepted row fails; a
// finding on a refused row must name SQL Server's error (walkCodes) or
// claim only facts the fixture's catalog holds.
func walkCatalogSimulation(t *testing.T, path, database string) walktest.SimulationResult {
	t.Helper()
	f, err := simulate.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	fixtures := make(map[string]catalogFixture)
	we := &walkEngine{}
	var res walktest.SimulationResult
	for _, r := range f.Rows {
		res.Total++
		where := fmt.Sprintf("%s:%d [%s]", path, r.Line(), r.Fixture)
		if r.Expected == simulate.Unjudged {
			t.Errorf("%s (%s): unjudged row:\n    %s", where, r.Rule, r.Statement)
			continue
		}
		fx, ok := fixtures[r.Fixture]
		if !ok {
			fx = prepareCatalogFixture(t, r.Fixture, database)
			fixtures[r.Fixture] = fx
		}
		finding, _ := reviewSimStatement(t, r.Statement, fx.target)
		if r.Accepted() {
			res.OK++
			if finding != nil {
				t.Errorf("%s (%s): WALK_THROUGH false positive %q on a statement SQL Server 2022 accepts:\n    %s",
					where, r.Rule, finding.Message, r.Statement)
			}
			continue
		}
		res.Refused++
		if finding == nil {
			continue
		}
		res.Found++
		if slices.Contains(we.Codes(finding.Message), r.Expected) {
			continue
		}
		v := simViolation(t, r.Statement, fx.target)
		if v == nil {
			t.Errorf("%s (%s): reported %q on SQL Server error %s with no violation to verify:\n    %s",
				where, r.Rule, finding.Message, r.Expected, r.Statement)
			continue
		}
		if ok, why := fx.cat.holds(v); !ok {
			t.Errorf("%s (%s): reported %q on SQL Server error %s, and the claim does not hold (%s):\n    %s",
				where, r.Rule, finding.Message, r.Expected, why, r.Statement)
		}
	}
	pct := 0.0
	if res.Refused > 0 {
		pct = 100 * float64(res.Found) / float64(res.Refused)
	}
	t.Logf("%s: %d rows (%d ok, %d refused); WALK_THROUGH found %d/%d refusals (%.1f%%)",
		path, res.Total, res.OK, res.Refused, res.Found, res.Refused, pct)
	return res
}

// catalogFixture is one prepared fixture: the target the walk reviews on,
// and the catalog its claims are checked against.
type catalogFixture struct {
	target review.Target
	cat    *simCatalog
}

// prepareCatalogFixture loads a mssql/catalog fixture into the in-memory
// catalog and returns the target the walk reviews on and the catalog its
// claims are checked against.
func prepareCatalogFixture(t *testing.T, fixture, database string) catalogFixture {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "catalog", "testdata", "fixtures", fixture+".sql"))
	if err != nil {
		t.Fatalf("read fixture %s: %v", fixture, err)
	}
	c, err := mssqlcatalog.LoadSQL(strings.TrimSpace(string(b)))
	if err != nil {
		t.Fatalf("mssql catalog failed to load fixture %s: %v", fixture, err)
	}
	meta := c.Metadata()
	meta.Name = database
	return catalogFixture{
		target: review.Target{Schema: meta, DefaultSchema: "dbo"},
		cat:    metadataSimCatalog(meta),
	}
}

// metadataSimCatalog reads a simCatalog off synced metadata, with the
// schemas every database holds that the sync leaves out.
func metadataSimCatalog(meta *metadata.DatabaseSchemaMetadata) *simCatalog {
	c := &simCatalog{
		schemaByCls: make(map[string]string),
		objects:     make(map[string][]catObject),
		relations:   make(map[string][]catObject),
		types:       make(map[string][]catObject),
		tableByPath: make(map[string]string),
		columns:     make(map[string][]catObject),
		indexes:     make(map[string][]catObject),
		constraints: make(map[string][]catObject),
		primaries:   make(map[string][]catObject),
	}
	schemaNames := append([]string{"dbo"}, excludedSchemas...)
	for _, s := range meta.GetSchemas() {
		if !slices.Contains(schemaNames, s.GetName()) {
			schemaNames = append(schemaNames, s.GetName())
		}
	}
	for _, name := range schemaNames {
		c.schemas = append(c.schemas, catObject{name: name})
		if n, ok := keyed(name); ok {
			if prev, exists := c.schemaByCls[n.class]; exists && prev != name {
				c.schemaByCls[n.class] = "*"
			} else {
				c.schemaByCls[n.class] = name
			}
		}
	}
	for _, s := range meta.GetSchemas() {
		sname := s.GetName()
		sn, sok := keyed(sname)
		object := func(name, kind string) {
			c.objects[sname] = append(c.objects[sname], catObject{name: name, kind: kind})
			if isRelation(kind) {
				c.relations[sname] = append(c.relations[sname], catObject{name: name, kind: kind})
			}
		}
		for _, tbl := range s.GetTables() {
			object(tbl.GetName(), kindTable)
			tkey := sname + "\x00" + tbl.GetName()
			if n, ok := keyed(tbl.GetName()); ok && sok {
				p := walk.Path(sn.class, n.class)
				if prev, exists := c.tableByPath[p]; exists && prev != tkey {
					c.tableByPath[p] = "*"
				} else {
					c.tableByPath[p] = tkey
				}
			}
			constraint := func(name, kind string) {
				object(name, kind)
				c.constraints[tkey] = append(c.constraints[tkey], catObject{name: name, kind: kind})
			}
			for _, col := range tbl.GetColumns() {
				c.columns[tkey] = append(c.columns[tkey], catObject{name: col.GetName()})
				if df := col.GetDefaultConstraintName(); df != "" {
					constraint(df, kindDF)
				}
			}
			for _, idx := range tbl.GetIndexes() {
				kind := "index"
				switch {
				case idx.GetPrimary():
					kind = kindPK
					constraint(idx.GetName(), kindPK)
					c.primaries[tkey] = append(c.primaries[tkey], catObject{name: primaryName})
				case idx.GetUnique() && idx.GetIsConstraint():
					kind = kindUQ
					constraint(idx.GetName(), kindUQ)
				}
				c.indexes[tkey] = append(c.indexes[tkey], catObject{name: idx.GetName(), kind: kind})
			}
			for _, fk := range tbl.GetForeignKeys() {
				constraint(fk.GetName(), kindFK)
			}
			for _, ck := range tbl.GetCheckConstraints() {
				constraint(ck.GetName(), kindCK)
			}
			for _, tr := range tbl.GetTriggers() {
				object(tr.GetName(), kindTrigger)
			}
		}
		for _, v := range s.GetViews() {
			object(v.GetName(), kindView)
			for _, tr := range v.GetTriggers() {
				object(tr.GetName(), kindTrigger)
			}
		}
		for _, q := range s.GetSequences() {
			object(q.GetName(), kindSequence)
		}
		for _, sy := range s.GetSynonyms() {
			object(sy.GetName(), kindSynonym)
		}
		for _, fn := range s.GetFunctions() {
			object(fn.GetName(), kindFunction)
		}
		for _, p := range s.GetProcedures() {
			object(p.GetName(), kindProcedure)
		}
	}
	return c
}
