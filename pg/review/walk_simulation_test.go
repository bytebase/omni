package review

import (
	"context"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/ast"
	pgcatalog "github.com/bytebase/omni/pg/catalog"
	"github.com/bytebase/omni/pg/parser"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/review/walk/walktest"
)

// TestWalkSimulation holds PostgreSQL's WALK_THROUGH to each of
// pg/catalog's simulate.tsv suites offline (PostgreSQL 17's, PostgreSQL
// 18's, and the view queries and stored programs on each): each fixture is built once in the
// in-memory catalog, and every row is checked with no server started.
func TestWalkSimulation(t *testing.T) {
	for _, suite := range []struct {
		name      string
		path      string
		minRecall float64
	}{
		// The DDL the walk misses names what the sync does not list: a
		// type that is no table, enum, range or composite (a domain, a base
		// or shell type, an enum with no labels), a tablespace, a role, an
		// aggregate (which ROUTINE and COMMENT ON FUNCTION also name), a
		// policy, an operator class, operator, access method or collation,
		// PostgreSQL 18's NOT NULL constraint names, and a typed table's
		// type. parser/testdata/compat.tsv repeats each statement once per
		// optional clause, so these are most of the misses.
		{"pg17", filepath.Join("..", "catalog", "testdata", "simulate.tsv"), 0.9116},
		{"pg18", filepath.Join("..", "catalog", "testdata", "pg18", "simulate.tsv"), 0.9099},
		// The view queries the walk misses are what pg/catalog refuses with a
		// code catalogCodes leaves out: a missing type, function, tablesample
		// method or schema, and a coercion a cast could allow.
		{"queries", filepath.Join("..", "catalog", "testdata", "queries", "simulate.tsv"), 0.988},
		{"pg18_queries", filepath.Join("..", "catalog", "testdata", "pg18", "queries", "simulate.tsv"), 0.988},
		// Functions, procedures, triggers, event triggers, CALL and DO on
		// programs@objects, quoted type names on programs@quoted, and ALTER
		// COLUMN TYPE ... USING on programs@retype and programs@generated.
		// pg/catalog checks a routine's body; what is left is a missing name
		// in one (42P01, 42703, 42704, 42883), a missing column type such as
		// a quoted "integer" (42704), a retyped column, default or USING
		// result that does not cast to the new type (42804, 42846, 22P02),
		// what a USING expression may not hold (42703, 42704, 0A000, 42803,
		// 42P20, 42P02, 42611 on a generated column), dependencies on a
		// dropped table or column (2BP01), and what a catalog loaded from a
		// snapshot does not hold: a procedure's identity, an event trigger.
		{"programs", filepath.Join("..", "catalog", "testdata", "programs", "simulate.tsv"), 0.427},
		{"pg18_programs", filepath.Join("..", "catalog", "testdata", "pg18", "programs", "simulate.tsv"), 0.427},
	} {
		t.Run(suite.name, func(t *testing.T) {
			res := walktest.CheckSimulation(t, walktest.Simulation{
				Path:     suite.path,
				Fixtures: filepath.Join("..", "catalog", "testdata", "fixtures"),
				Setup:    pgSimFixture,
				Codes:    func(msg, _ string) []string { return walkStates(msg) },
				Same:     func(claimed, stored walk.Key) bool { return claimed == stored },
			})
			if res.Refused > 0 {
				got := float64(res.Found) / float64(res.Refused)
				if got < suite.minRecall {
					t.Errorf("WALK_THROUGH recall = %.1f%% (%d/%d), want >= %.1f%%", 100*got, res.Found, res.Refused, 100*suite.minRecall)
				}
			}
		})
	}
}

func pgSimFixture(t testing.TB, fx simulate.Fixture) walktest.SimFixture {
	t.Helper()
	c, err := pgcatalog.LoadSQL(fx.SQL)
	if err != nil {
		t.Fatalf("pg catalog failed to load fixture %s: %v", fx.Name, err)
	}
	// pgCatalogToMetadata lists no rewrite rule, which the walk then takes
	// as absent from every table and view: a fixture must create none.
	if list, err := parser.Parse(fx.SQL); err == nil && list != nil {
		for _, item := range list.Items {
			if raw, ok := item.(*ast.RawStmt); ok {
				if _, rule := raw.Stmt.(*ast.RuleStmt); rule {
					t.Fatalf("fixture %s creates a rule, which pgCatalogToMetadata does not list", fx.Name)
				}
			}
		}
	}
	target := review.Target{
		Schema:      pgCatalogToMetadata(c, routineDefinitions(fx.SQL)),
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

// routineDefinitions returns, by schema and name in the order the fixture
// creates them, each routine's definition as the sync stores it:
// pg_get_functiondef's CREATE OR REPLACE FUNCTION or PROCEDURE, here the
// fixture's own statement.
func routineDefinitions(sql string) map[[2]string][]string {
	out := map[[2]string][]string{}
	list, err := parser.Parse(sql)
	if err != nil || list == nil {
		return out
	}
	for _, item := range list.Items {
		raw, ok := item.(*ast.RawStmt)
		if !ok {
			continue
		}
		fn, ok := raw.Stmt.(*ast.CreateFunctionStmt)
		if !ok || raw.Loc.Start < 0 || raw.Loc.End > len(sql) {
			continue
		}
		parts := nameParts(fn.Funcname)
		if len(parts) == 0 {
			continue
		}
		schema := "public"
		if len(parts) > 1 {
			schema = parts[len(parts)-2]
		}
		text := strings.TrimSpace(sql[raw.Loc.Start:raw.Loc.End])
		if !fn.IsOrReplace {
			text = "CREATE OR REPLACE" + strings.TrimPrefix(text, "CREATE")
		}
		key := [2]string{schema, parts[len(parts)-1]}
		out[key] = append(out[key], text)
	}
	return out
}

// viewDefinition is a view's definition as the sync reads it,
// pg_get_viewdef with an empty search_path, which qualifies every name; ""
// where the catalog cannot deparse it.
func viewDefinition(c *pgcatalog.Catalog, schema, name string, materialized bool) string {
	c.SetSearchPath(nil)
	defer c.SetSearchPath([]string{"public"})
	get := c.GetViewDefinition
	if materialized {
		get = c.GetMatViewDefinition
	}
	def, err := get(schema, name)
	if err != nil {
		return ""
	}
	return def
}

func pgCatalogToMetadata(c *pgcatalog.Catalog, definitions map[[2]string][]string) *metadata.DatabaseSchemaMetadata {
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
			var keyParts []string
			exprIdx := 0
			nkey := idx.NKeyColumns
			if nkey <= 0 || nkey > len(idx.Columns) {
				nkey = len(idx.Columns)
			}
			for _, attnum := range idx.Columns[:nkey] {
				if attnum == 0 {
					if exprIdx < len(idx.Exprs) && idx.Exprs[exprIdx] != "" {
						keyParts = append(keyParts, "("+idx.Exprs[exprIdx]+")")
					} else {
						keyParts = append(keyParts, "(1)")
					}
					exprIdx++
					continue
				}
				for _, col := range rel.Columns {
					if col.AttNum == attnum {
						keyParts = append(keyParts, ident(col.Name))
						break
					}
				}
			}
			if len(keyParts) == 0 {
				keyParts = []string{"a"}
			}
			out = append(out, &metadata.IndexMetadata{
				Name:         idx.Name,
				Unique:       idx.IsUnique,
				Primary:      idx.IsPrimary,
				IsConstraint: idx.ConstraintOID != 0,
				Definition:   "CREATE " + unique + "INDEX " + ident(idx.Name) + " ON " + only + ident(schemaName) + "." + ident(rel.Name) + " USING " + ident(am) + " (" + strings.Join(keyParts, ", ") + ")",
			})
		}
		return out
	}
	partKeyExpr := func(rel *pgcatalog.Relation) string {
		if rel.PartitionInfo == nil {
			return ""
		}
		var strat string
		switch rel.PartitionInfo.Strategy {
		case 'r':
			strat = "RANGE"
		case 'l':
			strat = "LIST"
		case 'h':
			strat = "HASH"
		default:
			return ""
		}
		var cols []string
		for _, attnum := range rel.PartitionInfo.KeyAttNums {
			for _, col := range rel.Columns {
				if col.AttNum == attnum {
					cols = append(cols, ident(col.Name))
					break
				}
			}
		}
		if len(cols) == 0 {
			return ""
		}
		return strat + " (" + strings.Join(cols, ", ") + ")"
	}
	var buildPartitions func(parent *pgcatalog.Relation) []*metadata.TablePartitionMetadata
	buildPartitions = func(parent *pgcatalog.Relation) []*metadata.TablePartitionMetadata {
		children := partChildren[parent.OID]
		if len(children) == 0 {
			return nil
		}
		expr := partKeyExpr(parent)
		out := make([]*metadata.TablePartitionMetadata, 0, len(children))
		for _, child := range children {
			out = append(out, &metadata.TablePartitionMetadata{
				Name:          child.Name,
				Expression:    expr,
				Indexes:       indexMeta(child),
				Subpartitions: buildPartitions(child),
			})
		}
		return out
	}
	// Map proc OID -> (schema, name) for synced trigger definitions.
	procByOID := make(map[uint32][2]string)
	for _, s := range c.UserSchemas() {
		for _, proc := range c.QueryPgProc(s.Name) {
			procByOID[proc.OID] = [2]string{s.Name, proc.ProName}
		}
	}
	colMeta := func(col *pgcatalog.Column, pos int32) *metadata.ColumnMetadata {
		cm := &metadata.ColumnMetadata{
			Name:     col.Name,
			Position: pos,
			Type:     c.FormatType(col.TypeOID, col.TypeMod),
			Nullable: !col.NotNull,
		}
		switch col.Identity {
		case 'a':
			cm.IdentityGeneration = metadata.ColumnMetadata_ALWAYS
		case 'd':
			cm.IdentityGeneration = metadata.ColumnMetadata_BY_DEFAULT
		}
		switch col.Generated {
		case 's':
			cm.Generation = &metadata.GenerationMetadata{Type: metadata.GenerationMetadata_TYPE_STORED, Expression: col.GenerationExpr}
		case 'v':
			cm.Generation = &metadata.GenerationMetadata{Type: metadata.GenerationMetadata_TYPE_VIRTUAL, Expression: col.GenerationExpr}
		}
		return cm
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
					Partitions: buildPartitions(rel),
				}
				for i, col := range rel.Columns {
					tbl.Columns = append(tbl.Columns, colMeta(col, int32(i+1)))
				}
				for _, con := range c.ConstraintsOf(rel.OID) {
					switch con.Type {
					case pgcatalog.ConstraintFK:
						fk := &metadata.ForeignKeyMetadata{Name: con.Name}
						if ref := c.GetRelationByOID(con.FRelOID); ref != nil {
							fk.ReferencedTable = ref.Name
							if ref.Schema != nil {
								fk.ReferencedSchema = ref.Schema.Name
							}
						}
						tbl.ForeignKeys = append(tbl.ForeignKeys, fk)
					case pgcatalog.ConstraintCheck:
						tbl.CheckConstraints = append(tbl.CheckConstraints, &metadata.CheckConstraintMetadata{Name: con.Name})
					case pgcatalog.ConstraintExclude:
						tbl.ExcludeConstraints = append(tbl.ExcludeConstraints, &metadata.ExcludeConstraintMetadata{Name: con.Name})
					}
				}
				for _, tr := range c.TriggersOf(rel.OID) {
					tm := &metadata.TriggerMetadata{Name: tr.Name}
					if fn, ok := procByOID[tr.FuncOID]; ok {
						tm.Body = "CREATE TRIGGER " + ident(tr.Name) + " BEFORE INSERT ON " + ident(s.Name) + "." + ident(rel.Name) + " FOR EACH ROW EXECUTE FUNCTION " + ident(fn[0]) + "." + ident(fn[1]) + "()"
					}
					tbl.Triggers = append(tbl.Triggers, tm)
				}
				sm.Tables = append(sm.Tables, tbl)
			case 'v':
				v := &metadata.ViewMetadata{Name: rel.Name, Definition: viewDefinition(c, s.Name, rel.Name, false)}
				for i, col := range rel.Columns {
					v.Columns = append(v.Columns, colMeta(col, int32(i+1)))
				}
				for _, tr := range c.TriggersOf(rel.OID) {
					v.Triggers = append(v.Triggers, &metadata.TriggerMetadata{Name: tr.Name})
				}
				sm.Views = append(sm.Views, v)
			case 'm':
				mv := &metadata.MaterializedViewMetadata{
					Name:       rel.Name,
					Definition: viewDefinition(c, s.Name, rel.Name, true),
					Indexes:    indexMeta(rel),
				}
				sm.MaterializedViews = append(sm.MaterializedViews, mv)
			case 'f':
				ft := &metadata.ExternalTableMetadata{Name: rel.Name}
				for i, col := range rel.Columns {
					ft.Columns = append(ft.Columns, colMeta(col, int32(i+1)))
				}
				sm.ExternalTables = append(sm.ExternalTables, ft)
			case 'c':
				ct := &metadata.CompositeTypeMetadata{Name: rel.Name}
				for _, col := range rel.Columns {
					ct.Attributes = append(ct.Attributes, &metadata.CompositeTypeAttribute{
						Name: col.Name,
						Type: c.FormatType(col.TypeOID, col.TypeMod),
					})
				}
				sm.CompositeTypes = append(sm.CompositeTypes, ct)
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
			switch typ.TypType {
			case 'e':
				var labels []string
				for _, l := range c.QueryPgEnum(typ.OID) {
					labels = append(labels, l.EnumLabel)
				}
				sm.EnumTypes = append(sm.EnumTypes, &metadata.EnumTypeMetadata{Name: typ.TypName, Values: labels})
			case 'r':
				sm.RangeTypes = append(sm.RangeTypes, &metadata.RangeTypeMetadata{Name: typ.TypName})
			}
		}
		for _, proc := range c.QueryPgProc(s.Name) {
			if builtinCalls[strings.ToLower(proc.ProName)] {
				continue
			}
			// The sync stores a procedure among the functions, each with
			// its definition.
			if defs := definitions[[2]string{s.Name, proc.ProName}]; len(defs) > 0 {
				definitions[[2]string{s.Name, proc.ProName}] = defs[1:]
				sm.Functions = append(sm.Functions, &metadata.FunctionMetadata{Name: proc.ProName, Definition: defs[0]})
				continue
			}
			if proc.ProKind == 'p' {
				sm.Procedures = append(sm.Procedures, &metadata.ProcedureMetadata{Name: proc.ProName})
			} else {
				var args []string
				for _, oid := range proc.ProArgTypes {
					args = append(args, c.FormatType(oid, -1))
				}
				ret := "void"
				if proc.ProRetType != 0 {
					ret = c.FormatType(proc.ProRetType, -1)
				}
				if proc.ProRetSet {
					ret = "SETOF " + ret
				}
				sm.Functions = append(sm.Functions, &metadata.FunctionMetadata{
					Name:       proc.ProName,
					Definition: "CREATE OR REPLACE FUNCTION " + ident(s.Name) + "." + ident(proc.ProName) + "(" + strings.Join(args, ", ") + ") RETURNS " + ret + " LANGUAGE sql AS ''",
				})
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
