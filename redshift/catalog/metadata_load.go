package catalog

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
	nodes "github.com/bytebase/omni/redshift/ast"
	pgparser "github.com/bytebase/omni/redshift/parser"
)

// LoadMetadata installs the schemas, tables, views, materialized views, and
// sequences of a Bytebase schema snapshot into c, as completion needs them.
// Columns take a coarse type and a view selects NULL under its column names,
// so no object depends on another. A materialized view installs from its
// definition when that runs, and like a view otherwise. An object that still
// fails to install is left out.
func (c *Catalog) LoadMetadata(meta *metadata.DatabaseSchemaMetadata) {
	defer c.SetSearchPath(slices.Clone(c.searchPath))
	for _, s := range meta.GetSchemas() {
		schema := s.GetName()
		c.execMetadataDDL(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", metadataQuoteIdent(schema)))
		for _, t := range s.GetTables() {
			c.execMetadataDDL(metadataTableDDL(schema, t.GetName(), t.GetColumns()))
		}
		for _, t := range s.GetExternalTables() {
			c.execMetadataDDL(metadataTableDDL(schema, t.GetName(), t.GetColumns()))
		}
		for _, v := range s.GetViews() {
			c.execMetadataDDL(metadataViewDDL("VIEW", schema, v.GetName(), v.GetColumns()))
		}
		for _, mv := range s.GetMaterializedViews() {
			// The definition may name relations of its own schema unqualified.
			c.SetSearchPath([]string{schema, "public"})
			if mv.GetDefinition() == "" || !c.execMetadataDDL(metadataMatViewDDL(schema, mv.GetName(), mv.GetDefinition())) {
				c.execMetadataDDL(metadataViewDDL("MATERIALIZED VIEW", schema, mv.GetName(), nil))
			}
		}
		for _, seq := range s.GetSequences() {
			c.execMetadataDDL(fmt.Sprintf("CREATE SEQUENCE %s.%s;", metadataQuoteIdent(schema), metadataQuoteIdent(seq.GetName())))
		}
	}
}

// UseMetadata makes c resolve a relation it does not have from a Bytebase
// schema snapshot, when a statement first names it, and installs the
// snapshot's functions. A view resolves to its definition, with each
// unqualified relation in it qualified by the schema it resolves to from the
// view's own schema, so the definition reads the same under c's search path.
func (c *Catalog) UseMetadata(meta *metadata.DatabaseSchemaMetadata) {
	c.SetRelationResolver(newMetadataResolver(meta))
	for _, s := range meta.GetSchemas() {
		if len(s.GetFunctions()) == 0 {
			continue
		}
		// A function installs only into a schema c has.
		c.execMetadataDDL(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", metadataQuoteIdent(s.GetName())))
		for _, fn := range s.GetFunctions() {
			if def := strings.TrimSpace(fn.GetDefinition()); def != "" {
				// A function that does not install is left out.
				c.execMetadataDDL(def)
			}
		}
	}
}

// execMetadataDDL runs sql and reports whether every statement in it ran.
// Exec fails only on SQL that does not parse; a statement that fails reports
// its error in its result.
func (c *Catalog) execMetadataDDL(sql string) bool {
	results, err := c.Exec(sql, nil)
	if err != nil {
		return false
	}
	for _, r := range results {
		if r.Error != nil {
			return false
		}
	}
	return true
}

// metadataPlaceholderColumn names the one column of a view or resolved table
// the snapshot lists no columns for, since each needs one.
const metadataPlaceholderColumn = "__stand_in"

func metadataTableDDL(schema, name string, columns []*metadata.ColumnMetadata) string {
	defs := make([]string, 0, len(columns))
	for _, col := range columns {
		if col.GetName() != "" {
			defs = append(defs, metadataQuoteIdent(col.GetName())+" "+coarseType(col.GetType()))
		}
	}
	return fmt.Sprintf("CREATE TABLE %s.%s (%s);", metadataQuoteIdent(schema), metadataQuoteIdent(name), strings.Join(defs, ", "))
}

func metadataViewDDL(kind, schema, name string, columns []*metadata.ColumnMetadata) string {
	targets := make([]string, 0, len(columns))
	for _, col := range columns {
		if col.GetName() != "" {
			targets = append(targets, fmt.Sprintf("CAST(NULL AS %s) AS %s", coarseType(col.GetType()), metadataQuoteIdent(col.GetName())))
		}
	}
	if len(targets) == 0 {
		targets = append(targets, "1 AS "+metadataPlaceholderColumn)
	}
	return fmt.Sprintf("CREATE %s %s.%s AS SELECT %s;", kind, metadataQuoteIdent(schema), metadataQuoteIdent(name), strings.Join(targets, ", "))
}

// metadataMatViewDDL takes a definition that is either the full CREATE
// statement or its query.
func metadataMatViewDDL(schema, name, definition string) string {
	definition = strings.TrimSuffix(strings.TrimSpace(definition), ";")
	if strings.HasPrefix(strings.ToUpper(definition), "CREATE ") {
		return definition + ";"
	}
	return fmt.Sprintf("CREATE MATERIALIZED VIEW %s.%s AS %s;", metadataQuoteIdent(schema), metadataQuoteIdent(name), definition)
}

// coarseType maps a column type to the few types completion and lineage need,
// so no type the catalog lacks can fail an install.
func coarseType(typ string) string {
	lower := strings.ToLower(strings.TrimSpace(typ))
	switch {
	case strings.Contains(lower, "interval"):
		return "interval"
	case strings.Contains(lower, "bigint"):
		return "bigint"
	case strings.Contains(lower, "int"):
		return "integer"
	case strings.Contains(lower, "bool"):
		return "boolean"
	case strings.Contains(lower, "char"), strings.Contains(lower, "text"), strings.Contains(lower, "string"):
		return "text"
	case strings.Contains(lower, "numeric"), strings.Contains(lower, "decimal"):
		return "numeric"
	case strings.Contains(lower, "double"):
		return "double precision"
	case strings.Contains(lower, "float"), strings.Contains(lower, "real"):
		return "real"
	case strings.Contains(lower, "date"):
		return "date"
	case strings.Contains(lower, "time"):
		return "timestamp"
	default:
		return "text"
	}
}

func metadataQuoteIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// metadataResolver resolves relations from a snapshot, matching names exactly
// as Redshift sync reports them.
type metadataResolver struct {
	// schemas maps a schema name to its relations and sequences by name.
	schemas map[string]map[string]metadataRelation
}

type metadataRelation struct {
	kind       byte
	columns    []*metadata.ColumnMetadata
	definition string
	sequence   bool
}

func newMetadataResolver(meta *metadata.DatabaseSchemaMetadata) *metadataResolver {
	r := &metadataResolver{schemas: make(map[string]map[string]metadataRelation)}
	for _, s := range meta.GetSchemas() {
		// Relations share one namespace; on a clash, the kind added last wins,
		// tables first.
		relations := make(map[string]metadataRelation)
		for _, seq := range s.GetSequences() {
			relations[seq.GetName()] = metadataRelation{kind: RelationKindTable, sequence: true}
		}
		for _, mv := range s.GetMaterializedViews() {
			relations[mv.GetName()] = metadataRelation{kind: RelationKindMaterializedView, definition: mv.GetDefinition()}
		}
		for _, v := range s.GetViews() {
			relations[v.GetName()] = metadataRelation{kind: RelationKindView, columns: v.GetColumns(), definition: v.GetDefinition()}
		}
		for _, t := range s.GetExternalTables() {
			relations[t.GetName()] = metadataRelation{kind: RelationKindTable, columns: t.GetColumns()}
		}
		for _, t := range s.GetTables() {
			relations[t.GetName()] = metadataRelation{kind: RelationKindTable, columns: t.GetColumns()}
		}
		r.schemas[s.GetName()] = relations
	}
	return r
}

func (r *metadataResolver) ResolveRelation(schemaName, relationName string, searchPath []string) (*RelationSpec, error) {
	if schemaName == "" {
		schemaName = r.search(searchPath, relationName)
	}
	rel, ok := r.schemas[schemaName][relationName]
	if !ok {
		return nil, nil
	}
	spec := &RelationSpec{SchemaName: schemaName, Name: relationName, Kind: rel.kind}
	switch definition := strings.TrimSpace(rel.definition); {
	case rel.sequence:
		spec.Columns = []RelationColumnSpec{{Name: "last_value", Type: "bigint"}, {Name: "log_cnt", Type: "bigint"}, {Name: "is_called", Type: "boolean"}}
	case rel.kind != RelationKindTable && definition != "":
		spec.Definition = r.qualify(definition, schemaName)
	default:
		// A view without a definition resolves as a table of its columns.
		spec.Kind = RelationKindTable
		spec.Columns = metadataColumnSpecs(rel.columns)
	}
	return spec, nil
}

// search returns the first schema on searchPath with a relation or sequence
// named name.
func (r *metadataResolver) search(searchPath []string, name string) string {
	for _, schema := range searchPath {
		if _, ok := r.schemas[schema][name]; ok {
			return schema
		}
	}
	return ""
}

func metadataColumnSpecs(columns []*metadata.ColumnMetadata) []RelationColumnSpec {
	specs := make([]RelationColumnSpec, 0, len(columns))
	for _, col := range columns {
		if col.GetName() != "" {
			specs = append(specs, RelationColumnSpec{Name: col.GetName(), Type: coarseType(col.GetType())})
		}
	}
	if len(specs) == 0 {
		specs = append(specs, RelationColumnSpec{Name: metadataPlaceholderColumn, Type: "text"})
	}
	return specs
}

// qualify schema-qualifies each relation a view definition names without a
// schema, resolving it from the view's own schema. A definition that does not
// parse stays as it is.
func (r *metadataResolver) qualify(definition, schema string) string {
	list, err := pgparser.Parse(definition)
	if err != nil || list == nil || len(list.Items) != 1 {
		return definition
	}
	var query nodes.Node
	switch stmt := unwrapRawStmt(list.Items[0]).(type) {
	case *nodes.ViewStmt:
		query = stmt.Query
	case *nodes.CreateTableAsStmt:
		query = stmt.Query
	case *nodes.SelectStmt:
		query = stmt
	default:
	}
	if query == nil {
		return definition
	}
	searchPath := []string{"public"}
	if schema != "" {
		searchPath = []string{schema, "public"}
	}
	var edits []qualification
	r.collectQualifications(searchPath, query, map[string]bool{}, &edits)
	// Insert from the end so earlier offsets stay valid.
	slices.SortFunc(edits, func(a, b qualification) int { return b.offset - a.offset })
	for _, edit := range edits {
		if edit.offset >= 0 && edit.offset <= len(definition) {
			definition = definition[:edit.offset] + metadataQuoteIdent(edit.schema) + "." + definition[edit.offset:]
		}
	}
	return definition
}

// qualification inserts schema before the relation name at offset.
type qualification struct {
	offset int
	schema string
}

// collectQualifications finds the unqualified relation names in node that are
// not CTE references and resolve along searchPath. CTE names match as parsed:
// unquoted names folded to lower case, quoted ones as written.
func (r *metadataResolver) collectQualifications(searchPath []string, node nodes.Node, ctes map[string]bool, edits *[]qualification) {
	if node == nil {
		return
	}
	if sel, ok := node.(*nodes.SelectStmt); ok {
		r.collectSelectQualifications(searchPath, sel, ctes, edits)
		return
	}
	nodes.Inspect(node, func(n nodes.Node) bool {
		if sel, ok := n.(*nodes.SelectStmt); ok && n != node {
			r.collectSelectQualifications(searchPath, sel, ctes, edits)
			return false
		}
		rv, ok := n.(*nodes.RangeVar)
		if !ok || rv == nil || rv.Relname == "" || rv.Schemaname != "" || rv.Catalogname != "" || rv.Loc.Start < 0 || ctes[rv.Relname] {
			return true
		}
		if schema := r.search(searchPath, rv.Relname); schema != "" {
			*edits = append(*edits, qualification{offset: rv.Loc.Start, schema: schema})
		}
		return true
	})
}

// collectSelectQualifications scopes the CTEs a SELECT declares: a recursive
// WITH sees all of them in every CTE, and otherwise each CTE sees only the
// ones before it.
func (r *metadataResolver) collectSelectQualifications(searchPath []string, sel *nodes.SelectStmt, ctes map[string]bool, edits *[]qualification) {
	local := maps.Clone(ctes)
	if sel.WithClause != nil && sel.WithClause.Ctes != nil {
		scope := maps.Clone(ctes)
		if sel.WithClause.Recursive {
			for _, item := range sel.WithClause.Ctes.Items {
				if cte, ok := item.(*nodes.CommonTableExpr); ok && cte.Ctename != "" {
					scope[cte.Ctename] = true
				}
			}
		}
		for _, item := range sel.WithClause.Ctes.Items {
			cte, ok := item.(*nodes.CommonTableExpr)
			if !ok || cte.Ctename == "" {
				continue
			}
			r.collectQualifications(searchPath, cte.Ctequery, scope, edits)
			scope[cte.Ctename] = true
			local[cte.Ctename] = true
		}
	}
	for _, n := range []nodes.Node{sel.WhereClause, sel.HavingClause, sel.QualifyClause, sel.LimitOffset, sel.LimitCount} {
		if n != nil {
			r.collectQualifications(searchPath, n, local, edits)
		}
	}
	if sel.IntoClause != nil {
		r.collectQualifications(searchPath, sel.IntoClause, local, edits)
	}
	if sel.Larg != nil {
		r.collectQualifications(searchPath, sel.Larg, local, edits)
	}
	if sel.Rarg != nil {
		r.collectQualifications(searchPath, sel.Rarg, local, edits)
	}
	for _, list := range []*nodes.List{sel.DistinctClause, sel.TargetList, sel.FromClause, sel.GroupClause, sel.WindowClause, sel.ExcludeList, sel.ValuesLists, sel.SortClause, sel.LockingClause} {
		if list == nil {
			continue
		}
		for _, item := range list.Items {
			r.collectQualifications(searchPath, item, local, edits)
		}
	}
}
