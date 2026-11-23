package catalog

import (
	"fmt"
	"math/big"
	"slices"
	"strconv"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql"
	"github.com/bytebase/omni/mssql/ast"
)

// objectKind names a schema-scoped object kind as SQL Server identifies it in
// namespace and DROP/ALTER validation.
type objectKind string

const (
	kindNone       objectKind = ""
	kindTable      objectKind = "table"
	kindView       objectKind = "view"
	kindProcedure  objectKind = "procedure"
	kindFunction   objectKind = "function"
	kindSequence   objectKind = "sequence"
	kindSynonym    objectKind = "synonym"
	kindTrigger    objectKind = "trigger"
	kindConstraint objectKind = "constraint"
)

// resolveTableInSchema resolves name within schema using the same exact-first,
// reject-ambiguous-fold semantics as GetTable.
func (c *Catalog) resolveTableInSchema(schema *metadata.SchemaMetadata, name string) *metadata.TableMetadata {
	if schema == nil || name == "" {
		return nil
	}
	if table := c.tableByName[tableKey{schema: schema.Name, table: name}]; table != nil {
		return table
	}
	var match *metadata.TableMetadata
	for _, table := range schema.Tables {
		if strings.EqualFold(table.Name, name) {
			if match != nil {
				return nil
			}
			match = table
		}
	}
	return match
}

// lookupSchemaObject returns the kind of the object named name in schema,
// checking tables, views, procedures, functions, sequences, synonyms, triggers,
// and named constraints. When exactTableCase is true, table matching requires
// exact case so case-distinct tables in CS_AS collations can coexist.
func (c *Catalog) lookupSchemaObject(schema *metadata.SchemaMetadata, name string, exactTableCase bool) objectKind {
	if schema == nil || name == "" {
		return kindNone
	}
	// Pass 1: exact name match across the entire schema-object namespace so
	// case-distinct objects in case-sensitive databases remain independently
	// addressable regardless of object kind.
	for _, t := range schema.Tables {
		if t.Name == name {
			return kindTable
		}
	}
	for _, v := range schema.Views {
		if v.Name == name {
			return kindView
		}
	}
	for _, p := range schema.Procedures {
		if p.Name == name {
			return kindProcedure
		}
	}
	for _, fn := range schema.Functions {
		if fn.Name == name {
			return kindFunction
		}
	}
	for _, sq := range schema.Sequences {
		if sq.Name == name {
			return kindSequence
		}
	}
	for _, sy := range schema.Synonyms {
		if sy.Name == name {
			return kindSynonym
		}
	}
	for _, t := range schema.Tables {
		for _, tr := range t.Triggers {
			if tr.Name == name {
				return kindTrigger
			}
		}
	}
	for _, v := range schema.Views {
		for _, tr := range v.Triggers {
			if tr.Name == name {
				return kindTrigger
			}
		}
	}
	for _, t := range schema.Tables {
		for _, idx := range t.Indexes {
			if idx.IsConstraint && idx.Name != "" && idx.Name == name {
				return kindConstraint
			}
		}
		for _, fk := range t.ForeignKeys {
			if fk.Name != "" && fk.Name == name {
				return kindConstraint
			}
		}
		for _, ck := range t.CheckConstraints {
			if ck.Name != "" && ck.Name == name {
				return kindConstraint
			}
		}
		for _, col := range t.Columns {
			if col.DefaultConstraintName != "" && col.DefaultConstraintName == name {
				return kindConstraint
			}
		}
	}

	// Pass 2: case-folded fallback.
	if !exactTableCase && c.resolveTableInSchema(schema, name) != nil {
		return kindTable
	}
	for _, v := range schema.Views {
		if strings.EqualFold(v.Name, name) {
			return kindView
		}
	}
	for _, p := range schema.Procedures {
		if strings.EqualFold(p.Name, name) {
			return kindProcedure
		}
	}
	for _, fn := range schema.Functions {
		if strings.EqualFold(fn.Name, name) {
			return kindFunction
		}
	}
	for _, sq := range schema.Sequences {
		if strings.EqualFold(sq.Name, name) {
			return kindSequence
		}
	}
	for _, sy := range schema.Synonyms {
		if strings.EqualFold(sy.Name, name) {
			return kindSynonym
		}
	}
	for _, t := range schema.Tables {
		for _, tr := range t.Triggers {
			if strings.EqualFold(tr.Name, name) {
				return kindTrigger
			}
		}
	}
	for _, v := range schema.Views {
		for _, tr := range v.Triggers {
			if strings.EqualFold(tr.Name, name) {
				return kindTrigger
			}
		}
	}
	for _, t := range schema.Tables {
		for _, idx := range t.Indexes {
			if idx.IsConstraint && idx.Name != "" && strings.EqualFold(idx.Name, name) {
				return kindConstraint
			}
		}
		for _, fk := range t.ForeignKeys {
			if fk.Name != "" && strings.EqualFold(fk.Name, name) {
				return kindConstraint
			}
		}
		for _, ck := range t.CheckConstraints {
			if ck.Name != "" && strings.EqualFold(ck.Name, name) {
				return kindConstraint
			}
		}
		for _, col := range t.Columns {
			if col.DefaultConstraintName != "" && strings.EqualFold(col.DefaultConstraintName, name) {
				return kindConstraint
			}
		}
	}
	return kindNone
}

// isValidDataType reports whether dt names a SQL Server built-in data type.
func isValidDataType(dt *ast.DataType, allowCursor bool) bool {
	if dt == nil {
		return false
	}
	// Schema-qualified custom types (e.g. dbo.MyType) are not tracked as
	// built-in types; treat unqualified names against the built-in registry.
	base := strings.Join(strings.Fields(strings.ToLower(dt.Name)), "")
	base = strings.TrimPrefix(strings.TrimSuffix(base, "]"), "[")
	if base == "cursor" {
		return allowCursor
	}
	return isBuiltinSQLServerType(base)
}

func isValidSequenceType(dt *ast.DataType) bool {
	if dt == nil {
		return true
	}
	base := strings.Join(strings.Fields(strings.ToLower(dt.Name)), "")
	base = strings.TrimPrefix(strings.TrimSuffix(base, "]"), "[")
	switch base {
	case "tinyint", "smallint", "int", "integer", "bigint", "decimal", "dec", "numeric":
		return true
	default:
		return false
	}
}

func isValidFilegroup(name string) bool {
	if name == "" {
		return true
	}
	u := strings.ToUpper(strings.Trim(name, "[]\"'"))
	return u == "PRIMARY" || u == "DEFAULT"
}

func hasInboundFK(c *Catalog, targetSchema string, targetTable *metadata.TableMetadata) (string, bool) {
	return hasInboundFKExcludingTables(c, targetSchema, targetTable, nil)
}

func hasInboundFKExcludingTables(c *Catalog, targetSchema string, targetTable *metadata.TableMetadata, exclude map[*metadata.TableMetadata]bool) (string, bool) {
	if targetTable == nil {
		return "", false
	}
	for _, s := range c.schemas {
		for _, t := range s.Tables {
			if t == targetTable {
				continue
			}
			if exclude != nil && exclude[t] {
				continue
			}
			for _, fk := range t.ForeignKeys {
				refSchema := fk.ReferencedSchema
				if refSchema == "" {
					refSchema = c.defaultSchema
				}
				if c.GetTable(refSchema, fk.ReferencedTable) == targetTable {
					return fk.Name, true
				}
			}
		}
	}
	return "", false
}

// indexReferencedByFK reports whether dropping idx on targetSchema.targetTable
// would leave any foreign key (including self-referencing FKs) without a
// matching primary-key or unique index on the target table.
func indexReferencedByFK(c *Catalog, targetSchema string, targetTable *metadata.TableMetadata, idx *metadata.IndexMetadata) (string, string, bool) {
	return indexReferencedByFKExcluding(c, targetSchema, targetTable, idx, nil)
}

func indexReferencedByFKExcluding(c *Catalog, targetSchema string, targetTable *metadata.TableMetadata, idx *metadata.IndexMetadata, excluded map[*metadata.IndexMetadata]bool) (string, string, bool) {
	if idx == nil || (!idx.Primary && !idx.Unique) || indexWhereClause(idx.Definition) != "" {
		return "", "", false
	}
	for _, s := range c.schemas {
		for _, t := range s.Tables {
			for _, fk := range t.ForeignKeys {
				refSchema := fk.ReferencedSchema
				if refSchema == "" {
					refSchema = c.defaultSchema
				}
				if c.GetTable(refSchema, fk.ReferencedTable) != targetTable {
					continue
				}
				if !fkMatchesIndex(fk, idx) {
					continue
				}
				if hasOtherMatchingUniqueIndexExcluding(targetTable, idx, excluded, fk) {
					continue
				}
				return t.Name, fk.Name, true
			}
		}
	}
	return "", "", false
}

func schemaBoundTableConsumer(c *Catalog, targetSchema, targetTable string) (string, bool) {
	return schemaBoundObjectConsumer(c, targetSchema, targetTable, nil)
}

func schemaBoundObjectConsumer(c *Catalog, targetSchema, targetObj string, excluded func(schemaName, objName string) bool) (string, bool) {
	for _, s := range c.schemas {
		for _, v := range s.Views {
			if strings.EqualFold(s.Name, targetSchema) && strings.EqualFold(v.Name, targetObj) {
				continue
			}
			if excluded != nil && excluded(s.Name, v.Name) {
				continue
			}
			if isSchemaBound(v.Definition) && definitionReferencesTable(v.Definition, s.Name, targetSchema, targetObj) {
				return v.Name, true
			}
		}
		for _, fn := range s.Functions {
			if strings.EqualFold(s.Name, targetSchema) && strings.EqualFold(fn.Name, targetObj) {
				continue
			}
			if excluded != nil && excluded(s.Name, fn.Name) {
				continue
			}
			if isSchemaBound(fn.Definition) && definitionReferencesTable(fn.Definition, s.Name, targetSchema, targetObj) {
				return fn.Name, true
			}
		}
	}
	return "", false
}

func schemaBoundFunctionAlterConsumer(c *Catalog, targetSchema, targetFunc string) (string, bool) {
	for _, s := range c.schemas {
		for _, v := range s.Views {
			if isSchemaBound(v.Definition) && definitionReferencesFunction(v.Definition, s.Name, targetSchema, targetFunc) {
				return v.Name, true
			}
		}
		for _, fn := range s.Functions {
			if strings.EqualFold(s.Name, targetSchema) && strings.EqualFold(fn.Name, targetFunc) {
				continue
			}
			if isSchemaBound(fn.Definition) && definitionReferencesFunction(fn.Definition, s.Name, targetSchema, targetFunc) {
				return fn.Name, true
			}
		}
	}
	return "", false
}

func schemaBoundFunctionConsumer(c *Catalog, targetSchema, targetFunc string, excluded func(schemaName, objName string) bool) (string, bool) {
	for _, s := range c.schemas {
		for _, v := range s.Views {
			if excluded != nil && excluded(s.Name, v.Name) {
				continue
			}
			if isSchemaBound(v.Definition) && definitionReferencesFunction(v.Definition, s.Name, targetSchema, targetFunc) {
				return v.Name, true
			}
		}
		for _, fn := range s.Functions {
			if strings.EqualFold(s.Name, targetSchema) && strings.EqualFold(fn.Name, targetFunc) {
				continue
			}
			if excluded != nil && excluded(s.Name, fn.Name) {
				continue
			}
			if isSchemaBound(fn.Definition) && definitionReferencesFunction(fn.Definition, s.Name, targetSchema, targetFunc) {
				return fn.Name, true
			}
		}
		for _, t := range s.Tables {
			for _, ck := range t.CheckConstraints {
				if expressionReferencesFunction(ck.Expression, s.Name, targetSchema, targetFunc) {
					return ck.Name, true
				}
			}
			for _, col := range t.Columns {
				if isComputedColumnType(col.Type) && expressionReferencesFunction(col.Type, s.Name, targetSchema, targetFunc) {
					return t.Name, true
				}
				if col.Default != "" && expressionReferencesFunction(col.Default, s.Name, targetSchema, targetFunc) {
					name := col.DefaultConstraintName
					if name == "" {
						name = t.Name
					}
					return name, true
				}
			}
		}
	}
	return "", false
}

func schemaBoundColumnConsumer(c *Catalog, targetSchema, targetTable, colName string) (string, bool) {
	for _, s := range c.schemas {
		for _, v := range s.Views {
			if isSchemaBound(v.Definition) && definitionReferencesTable(v.Definition, s.Name, targetSchema, targetTable) &&
				definitionReferencesTableColumn(c, v.Definition, s.Name, targetSchema, targetTable, colName) {
				return v.Name, true
			}
		}
		for _, fn := range s.Functions {
			if isSchemaBound(fn.Definition) && definitionReferencesTable(fn.Definition, s.Name, targetSchema, targetTable) &&
				definitionReferencesTableColumn(c, fn.Definition, s.Name, targetSchema, targetTable, colName) {
				return fn.Name, true
			}
		}
	}
	return "", false
}

type tableSourceBinding struct {
	alias    string
	hasAlias bool
	schema   string
	table    string
	isTarget bool
}

type columnRefResolver struct {
	cat           *Catalog
	defaultSchema string
	targetSchema  string
	targetTable   string
	targetCol     string
	cteStack      []map[string]bool
	scopeStack    [][]tableSourceBinding
	found         bool
}

func (r *columnRefResolver) inCTEScope(name string) bool {
	lower := strings.ToLower(name)
	for i := len(r.cteStack) - 1; i >= 0; i-- {
		if r.cteStack[i][lower] {
			return true
		}
	}
	return false
}

func (r *columnRefResolver) collectFromSources(from *ast.List) []tableSourceBinding {
	if from == nil {
		return nil
	}
	var out []tableSourceBinding
	var walkExpr func(te ast.TableExpr)
	walkExpr = func(te ast.TableExpr) {
		switch x := te.(type) {
		case *ast.JoinClause:
			walkExpr(x.Left)
			walkExpr(x.Right)
		case *ast.TableRef:
			if x.Object == "" {
				return
			}
			inCTE := x.Schema == "" && r.inCTEScope(x.Object)
			schema := x.Schema
			if schema == "" {
				schema = r.defaultSchema
			}
			alias := x.Alias
			hasAlias := alias != ""
			if alias == "" {
				alias = x.Object
			}
			isTarget := !inCTE && strings.EqualFold(schema, r.targetSchema) && strings.EqualFold(x.Object, r.targetTable)
			out = append(out, tableSourceBinding{
				alias:    alias,
				hasAlias: hasAlias,
				schema:   schema,
				table:    x.Object,
				isTarget: isTarget,
			})
		case *ast.AliasedTableRef:
			if tr, ok := x.Table.(*ast.TableRef); ok && tr.Object != "" {
				inCTE := tr.Schema == "" && r.inCTEScope(tr.Object)
				schema := tr.Schema
				if schema == "" {
					schema = r.defaultSchema
				}
				alias := x.Alias
				if alias == "" {
					alias = tr.Alias
				}
				hasAlias := alias != ""
				if alias == "" {
					alias = tr.Object
				}
				isTarget := !inCTE && strings.EqualFold(schema, r.targetSchema) && strings.EqualFold(tr.Object, r.targetTable)
				out = append(out, tableSourceBinding{
					alias:    alias,
					hasAlias: hasAlias,
					schema:   schema,
					table:    tr.Object,
					isTarget: isTarget,
				})
				return
			}
			if x.Alias != "" {
				out = append(out, tableSourceBinding{alias: x.Alias, hasAlias: true})
			}
		case *ast.TableVarRef:
			alias := x.Alias
			if alias == "" {
				alias = x.Name
			}
			out = append(out, tableSourceBinding{alias: alias, hasAlias: x.Alias != ""})
		case *ast.FullTextTableRef:
			if x.Alias != "" {
				out = append(out, tableSourceBinding{alias: x.Alias, hasAlias: true})
			}
		case *ast.SemanticTableRef:
			if x.Alias != "" {
				out = append(out, tableSourceBinding{alias: x.Alias, hasAlias: true})
			}
		case *ast.PredictTableRef:
			if x.Alias != "" {
				out = append(out, tableSourceBinding{alias: x.Alias, hasAlias: true})
			}
		case *ast.AdHocTableRef:
			if x.Alias != "" {
				out = append(out, tableSourceBinding{alias: x.Alias, hasAlias: true})
			}
		}
	}
	for _, item := range from.Items {
		if te, ok := item.(ast.TableExpr); ok {
			walkExpr(te)
		}
	}
	return out
}

func (r *columnRefResolver) matchesColumnRef(cr *ast.ColumnRef) bool {
	if !strings.EqualFold(cr.Column, r.targetCol) {
		return false
	}
	if cr.Table != "" {
		for i := len(r.scopeStack) - 1; i >= 0; i-- {
			matchedAny := false
			matchedTarget := false
			for _, b := range r.scopeStack[i] {
				matched := false
				if cr.Schema != "" {
					matched = !b.hasAlias && strings.EqualFold(b.schema, cr.Schema) && strings.EqualFold(b.table, cr.Table)
				} else {
					matched = strings.EqualFold(b.alias, cr.Table)
				}
				if matched {
					matchedAny = true
					if b.isTarget {
						matchedTarget = true
					}
				}
			}
			if matchedAny {
				return matchedTarget
			}
		}
		return strings.EqualFold(cr.Table, r.targetTable) && (cr.Schema == "" || strings.EqualFold(cr.Schema, r.targetSchema))
	}
	for i := len(r.scopeStack) - 1; i >= 0; i-- {
		hasTarget := false
		otherHasCol := false
		for _, b := range r.scopeStack[i] {
			if b.isTarget {
				hasTarget = true
				continue
			}
			if r.cat != nil && b.table != "" {
				if tbl := r.cat.GetTable(b.schema, b.table); tbl != nil && columnByName(tbl, r.targetCol) != nil {
					otherHasCol = true
				}
			}
		}
		if hasTarget {
			return true
		}
		if otherHasCol {
			return false
		}
	}
	return false
}

type popColScopeVisitor struct {
	r       *columnRefResolver
	popCTE  bool
	popScope bool
}

func (p popColScopeVisitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		if p.popScope {
			p.r.scopeStack = p.r.scopeStack[:len(p.r.scopeStack)-1]
		}
		if p.popCTE {
			p.r.cteStack = p.r.cteStack[:len(p.r.cteStack)-1]
		}
		return nil
	}
	return p.r.Visit(n)
}

func (r *columnRefResolver) Visit(n ast.Node) ast.Visitor {
	if n == nil || r.found {
		return nil
	}
	switch x := n.(type) {
	case *ast.CommonTableExpr:
		// A CTE query does not inherit the enclosing SELECT's FROM scope.
		saved := r.scopeStack
		r.scopeStack = nil
		if x.Query != nil {
			ast.Walk(r, x.Query)
		}
		r.scopeStack = saved
		return nil
	case *ast.SelectStmt:
		popCTE := false
		if names := withClauseCTENames(x.WithClause); len(names) > 0 {
			r.cteStack = append(r.cteStack, names)
			popCTE = true
		}
		sources := r.collectFromSources(x.FromClause)
		popScope := false
		if len(sources) > 0 {
			r.scopeStack = append(r.scopeStack, sources)
			popScope = true
		}
		if popCTE || popScope {
			return popColScopeVisitor{r: r, popCTE: popCTE, popScope: popScope}
		}
	case *ast.ColumnRef:
		if r.matchesColumnRef(x) {
			r.found = true
			return nil
		}
	}
	return r
}

func definitionReferencesTableColumn(c *Catalog, definition, defaultSchemaName, targetSchema, targetTable, colName string) bool {
	if strings.TrimSpace(definition) == "" || targetTable == "" || colName == "" {
		return false
	}
	if defaultSchemaName == "" {
		defaultSchemaName = defaultSchema
	}
	if targetSchema == "" {
		targetSchema = defaultSchema
	}
	stmts, err := mssql.Parse(creationDefinition(definition))
	if err != nil || len(stmts) == 0 {
		return false
	}
	var root ast.Node
	for i := range stmts {
		switch n := stmts[i].AST.(type) {
		case *ast.CreateViewStmt:
			root = n.Query
		case *ast.CreateFunctionStmt:
			root = n.Body
		}
		if root != nil {
			break
		}
	}
	if root == nil {
		return false
	}
	r := &columnRefResolver{
		cat:           c,
		defaultSchema: defaultSchemaName,
		targetSchema:  targetSchema,
		targetTable:   targetTable,
		targetCol:     colName,
	}
	ast.Walk(r, root)
	return r.found
}

func fkMatchesIndex(fk *metadata.ForeignKeyMetadata, idx *metadata.IndexMetadata) bool {
	if len(fk.ReferencedColumns) == 0 {
		return idx.Primary
	}
	if len(fk.ReferencedColumns) != len(idx.Expressions) {
		return false
	}
	for _, rc := range fk.ReferencedColumns {
		found := false
		for _, ic := range idx.Expressions {
			if strings.EqualFold(rc, ic) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func hasOtherMatchingUniqueIndex(table *metadata.TableMetadata, exclude *metadata.IndexMetadata, fk *metadata.ForeignKeyMetadata) bool {
	return hasOtherMatchingUniqueIndexExcluding(table, exclude, nil, fk)
}

func hasOtherMatchingUniqueIndexExcluding(table *metadata.TableMetadata, exclude *metadata.IndexMetadata, excluded map[*metadata.IndexMetadata]bool, fk *metadata.ForeignKeyMetadata) bool {
	for _, other := range table.Indexes {
		if other == exclude || excluded[other] || (!other.Primary && !other.Unique) {
			continue
		}
		if indexWhereClause(other.Definition) != "" {
			continue
		}
		if fkMatchesIndex(fk, other) {
			return true
		}
	}
	return false
}

func tableHasClusteredIndex(table *metadata.TableMetadata) bool {
	return tableHasClusteredIndexExcluding(table, nil)
}

func tableHasClusteredIndexExcluding(table *metadata.TableMetadata, exclude *metadata.IndexMetadata) bool {
	if table == nil {
		return false
	}
	for _, idx := range table.Indexes {
		if idx == exclude {
			continue
		}
		if strings.EqualFold(idx.Type, "CLUSTERED") || strings.EqualFold(idx.Type, "CLUSTERED COLUMNSTORE") {
			return true
		}
		if idx.Type == "" && idx.IsConstraint && idx.Primary {
			return true
		}
	}
	return false
}

func tableHasIdentityColumn(table *metadata.TableMetadata) bool {
	if table == nil {
		return false
	}
	for _, col := range table.Columns {
		if col.IsIdentity {
			return true
		}
	}
	return false
}

func selectStarWithoutFrom(node ast.Node) bool {
	if node == nil {
		return false
	}
	switch s := node.(type) {
	case *ast.SelectStmt:
		if s.FromClause == nil || len(s.FromClause.Items) == 0 {
			if s.TargetList != nil {
				for _, item := range s.TargetList.Items {
					if _, ok := item.(*ast.StarExpr); ok {
						return true
					}
					if res, ok := item.(*ast.ResTarget); ok {
						if _, ok := res.Val.(*ast.StarExpr); ok {
							return true
						}
					}
				}
			}
		}
	case *ast.ReturnStmt:
		return selectStarWithoutFrom(s.Value)
	case *ast.SubqueryExpr:
		return selectStarWithoutFrom(s.Query)
	case *ast.BeginEndStmt:
		if s.Stmts != nil {
			for _, item := range s.Stmts.Items {
				if selectStarWithoutFrom(item) {
					return true
				}
			}
		}
	}
	return false
}

func isPureAlterStmt(node ast.Node, source string) bool {
	text := strings.TrimSpace(nodeText(node, source))
	return len(text) >= 5 && strings.EqualFold(text[:5], "ALTER")
}

// --- CREATE / ALTER SCHEMA ---

func (c *Catalog) execCreateSchema(n *ast.CreateSchemaStmt, source string) error {
	// Legacy `CREATE SCHEMA AUTHORIZATION owner` with no schema name and no
	// elements is accepted as a no-op by SQL Server 2022.
	if n.Name == "" {
		if n.Elements == nil || len(n.Elements.Items) == 0 {
			return nil
		}
	}
	if !c.lenientMetadataLoad && n.Authorization != "" && !strings.EqualFold(n.Authorization, "dbo") {
		return mssqlErr(ErrSecurableNotFound, "Cannot find the user '%s', because it does not exist or you do not have permission.", n.Authorization)
	}
	if c.GetSchema(n.Name) != nil {
		return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", n.Name)
	}
	schema := c.getOrCreateSchema(n.Name)
	if n.Elements != nil {
		savedIndexCounter := c.indexCounter
		previous := c.defaultSchema
		if n.Name != "" {
			c.defaultSchema = n.Name
		}
		for _, item := range n.Elements.Items {
			if err := c.execStatement(item, source); err != nil {
				c.defaultSchema = previous
				c.indexCounter = savedIndexCounter
				c.removeSchema(schema)
				return err
			}
		}
		c.defaultSchema = previous
	}
	return nil
}

func (c *Catalog) removeSchema(target *metadata.SchemaMetadata) {
	if target == nil {
		return
	}
	for i, s := range c.schemas {
		if s == target {
			fold := normalize(s.Name)
			for _, t := range s.Tables {
				delete(c.tableByName, tableKey{schema: s.Name, table: t.Name})
				periodKey := tableKey{schema: normalizeOrDefault(s.Name, c.defaultSchema), table: t.Name}
				delete(c.systemTimeTables, periodKey)
				delete(c.systemVersionedTables, periodKey)
			}
			delete(c.schemaByName, s.Name)
			var kept []*metadata.SchemaMetadata
			for _, fs := range c.schemasByFold[fold] {
				if fs != s {
					kept = append(kept, fs)
				}
			}
			if len(kept) == 0 {
				delete(c.schemasByFold, fold)
			} else {
				c.schemasByFold[fold] = kept
			}
			c.schemas = append(c.schemas[:i], c.schemas[i+1:]...)
			return
		}
	}
}

func (c *Catalog) execAlterSchema(n *ast.AlterSchemaStmt) error {
	targetSchema := c.GetSchema(n.Name)
	if targetSchema == nil {
		return mssqlErr(ErrSecurableNotFound, "Cannot find the schema '%s', because it does not exist or you do not have permission.", n.Name)
	}
	if strings.EqualFold(n.TransferType, "TYPE") || strings.EqualFold(n.TransferType, "XML SCHEMA COLLECTION") {
		return mssqlErr(ErrSecurableNotFound, "Cannot find the %s '%s', because it does not exist or you do not have permission.", strings.ToLower(n.TransferType), n.TransferEntity)
	}
	srcSchemaName := c.defaultSchema
	objName := n.TransferObject
	if n.TransferSchema != "" {
		srcSchemaName = n.TransferSchema
	} else if objName == "" {
		objName = n.TransferEntity
		if dot := strings.LastIndex(objName, "."); dot >= 0 {
			srcSchemaName = objName[:dot]
			objName = objName[dot+1:]
		}
	}
	srcSchema := c.GetSchema(srcSchemaName)
	if srcSchema == nil {
		return mssqlErr(ErrSecurableNotFound, "Cannot find the object '%s', because it does not exist or you do not have permission.", n.TransferEntity)
	}
	kind := c.lookupSchemaObject(srcSchema, objName, false)
	if kind == kindNone || kind == kindTrigger || kind == kindConstraint {
		return mssqlErr(ErrSecurableNotFound, "Cannot find the object '%s', because it does not exist or you do not have permission.", n.TransferEntity)
	}
	if srcSchema == targetSchema {
		return nil
	}
	if c.lookupSchemaObject(targetSchema, objName, false) != kindNone {
		return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", objName)
	}
	switch kind {
	case kindTable:
		table := c.resolveTableInSchema(srcSchema, objName)
		if table == nil {
			return mssqlErr(ErrSecurableNotFound, "Cannot find the object '%s', because it does not exist or you do not have permission.", n.TransferEntity)
		}
		if _, bound := schemaBoundTableConsumer(c, srcSchema.Name, table.Name); bound {
			return mssqlErr(ErrTransferSchemaBound, "Cannot transfer a schemabound object.")
		}
		if collision := c.findTransferredTableCollision(targetSchema, table); collision != "" {
			return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", collision)
		}
		for i, t := range srcSchema.Tables {
			if t == table {
				var inboundFKs []*metadata.ForeignKeyMetadata
				for _, s := range c.schemas {
					for _, tbl := range s.Tables {
						for _, fk := range tbl.ForeignKeys {
							refSchema := fk.ReferencedSchema
							if refSchema == "" {
								refSchema = c.defaultSchema
							}
							if c.GetTable(refSchema, fk.ReferencedTable) == t {
								inboundFKs = append(inboundFKs, fk)
							}
						}
					}
				}
				srcSchema.Tables = append(srcSchema.Tables[:i], srcSchema.Tables[i+1:]...)
				delete(c.tableByName, tableKey{schema: srcSchema.Name, table: t.Name})
				oldPeriodKey := tableKey{schema: normalizeOrDefault(srcSchema.Name, c.defaultSchema), table: t.Name}
				newPeriodKey := tableKey{schema: normalizeOrDefault(targetSchema.Name, c.defaultSchema), table: t.Name}
				if periodCols, ok := c.systemTimeTables[oldPeriodKey]; ok {
					delete(c.systemTimeTables, oldPeriodKey)
					if c.systemTimeTables == nil {
						c.systemTimeTables = make(map[tableKey][2]string)
					}
					c.systemTimeTables[newPeriodKey] = periodCols
				}
				if c.systemVersionedTables[oldPeriodKey] {
					delete(c.systemVersionedTables, oldPeriodKey)
					if c.systemVersionedTables == nil {
						c.systemVersionedTables = make(map[tableKey]bool)
					}
					c.systemVersionedTables[newPeriodKey] = true
				}
				targetSchema.Tables = append(targetSchema.Tables, t)
				c.indexTable(targetSchema.Name, t)
				for _, tr := range t.Triggers {
					tr.Body = qualifyTriggerDefinition(tr.Body, targetSchema.Name, t.Name, true)
				}
				for _, fk := range inboundFKs {
					fk.ReferencedSchema = targetSchema.Name
				}
				return nil
			}
		}
	case kindView:
		i, ok := resolveExactOrFoldIndex(srcSchema.Views, func(v *metadata.ViewMetadata) string { return v.Name }, objName)
		if !ok {
			return mssqlErr(ErrSecurableNotFound, "Cannot find the object '%s', because it does not exist or you do not have permission.", n.TransferEntity)
		}
		v := srcSchema.Views[i]
		if isSchemaBound(v.Definition) {
			return mssqlErr(ErrTransferSchemaBound, "Cannot transfer a schemabound object.")
		}
		if _, bound := schemaBoundObjectConsumer(c, srcSchema.Name, v.Name, nil); bound {
			return mssqlErr(ErrTransferSchemaBound, "Cannot transfer a schemabound object.")
		}
		for _, tr := range v.Triggers {
			if tr.Name != "" && c.lookupSchemaObject(targetSchema, tr.Name, false) != kindNone {
				return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", tr.Name)
			}
		}
		srcSchema.Views = append(srcSchema.Views[:i], srcSchema.Views[i+1:]...)
		v.Definition = qualifyModuleDefinition(v.Definition, targetSchema.Name, true)
		for _, tr := range v.Triggers {
			tr.Body = qualifyTriggerDefinition(tr.Body, targetSchema.Name, v.Name, true)
		}
		targetSchema.Views = append(targetSchema.Views, v)
		return nil
	case kindProcedure:
		i, ok := resolveExactOrFoldIndex(srcSchema.Procedures, func(p *metadata.ProcedureMetadata) string { return p.Name }, objName)
		if !ok {
			return mssqlErr(ErrSecurableNotFound, "Cannot find the object '%s', because it does not exist or you do not have permission.", n.TransferEntity)
		}
		p := srcSchema.Procedures[i]
		srcSchema.Procedures = append(srcSchema.Procedures[:i], srcSchema.Procedures[i+1:]...)
		p.Definition = qualifyModuleDefinition(p.Definition, targetSchema.Name, true)
		targetSchema.Procedures = append(targetSchema.Procedures, p)
		return nil
	case kindFunction:
		i, ok := resolveExactOrFoldIndex(srcSchema.Functions, func(fn *metadata.FunctionMetadata) string { return fn.Name }, objName)
		if !ok {
			return mssqlErr(ErrSecurableNotFound, "Cannot find the object '%s', because it does not exist or you do not have permission.", n.TransferEntity)
		}
		fn := srcSchema.Functions[i]
		if isSchemaBound(fn.Definition) {
			return mssqlErr(ErrTransferSchemaBound, "Cannot transfer a schemabound object.")
		}
		if _, bound := schemaBoundFunctionConsumer(c, srcSchema.Name, fn.Name, nil); bound {
			return mssqlErr(ErrTransferSchemaBound, "Cannot transfer a schemabound object.")
		}
		srcSchema.Functions = append(srcSchema.Functions[:i], srcSchema.Functions[i+1:]...)
		fn.Definition = qualifyModuleDefinition(fn.Definition, targetSchema.Name, true)
		targetSchema.Functions = append(targetSchema.Functions, fn)
		return nil
	case kindSequence:
		i, ok := resolveExactOrFoldIndex(srcSchema.Sequences, func(sq *metadata.SequenceMetadata) string { return sq.Name }, objName)
		if !ok {
			return mssqlErr(ErrSecurableNotFound, "Cannot find the object '%s', because it does not exist or you do not have permission.", n.TransferEntity)
		}
		sq := srcSchema.Sequences[i]
		srcSchema.Sequences = append(srcSchema.Sequences[:i], srcSchema.Sequences[i+1:]...)
		targetSchema.Sequences = append(targetSchema.Sequences, sq)
		return nil
	case kindSynonym:
		i, ok := resolveExactOrFoldIndex(srcSchema.Synonyms, func(sy *metadata.SynonymMetadata) string { return sy.Name }, objName)
		if !ok {
			return mssqlErr(ErrSecurableNotFound, "Cannot find the object '%s', because it does not exist or you do not have permission.", n.TransferEntity)
		}
		sy := srcSchema.Synonyms[i]
		srcSchema.Synonyms = append(srcSchema.Synonyms[:i], srcSchema.Synonyms[i+1:]...)
		targetSchema.Synonyms = append(targetSchema.Synonyms, sy)
		return nil
	}
	return nil
}

// --- CREATE TABLE / CTAS ---

func (c *Catalog) execCreateTableAsSelect(n *ast.CreateTableAsSelectStmt) error {
	if n.Columns != nil && len(n.Columns.Items) > 0 {
		first := ""
		if s, ok := n.Columns.Items[0].(*ast.String); ok {
			first = s.Str
		}
		return mssqlErr(ErrInvalidTableDef, "The definition for column '%s' must include a data type.", first)
	}
	return mssqlErr(ErrSyntax, "Incorrect syntax near 'HEAP'.")
}

func (c *Catalog) execCreateTable(n *ast.CreateTableStmt, source string) error {
	// Validate column data types first (compile-time check in SQL Server).
	if !c.lenientMetadataLoad && n.Columns != nil {
		for idx, item := range n.Columns.Items {
			col, ok := item.(*ast.ColumnDef)
			if !ok {
				continue
			}
			if col.Computed == nil && !isValidDataType(col.DataType, false) {
				typeName := ""
				if col.DataType != nil {
					typeName = col.DataType.Name
				}
				return mssqlErr(ErrUnknownColumnType, "Column, parameter, or variable #%d: Cannot find data type %s.", idx+1, typeName)
			}
		}
	}

	schemaName, tableName := tableRefSchemaObject(n.Name, c.defaultSchema)
	if schema := c.schemaByName[schemaName]; schema != nil {
		if c.lookupSchemaObject(schema, tableName, true) != kindNone {
			return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", tableName)
		}
	}
	if !c.lenientMetadataLoad && n.IsFileTable {
		return mssqlErr(ErrFileTableNotSupported, "Default FILESTREAM filegroup is not available.")
	}
	wantsSystemVersioning := false
	if n.TableOptions != nil {
		for _, item := range n.TableOptions.Items {
			if opt, ok := item.(*ast.TableOption); ok && strings.EqualFold(opt.Name, "SYSTEM_VERSIONING") && strings.EqualFold(opt.Value, "ON") {
				wantsSystemVersioning = true
				break
			}
		}
	}
	if wantsSystemVersioning && !c.lenientMetadataLoad {
		if n.PeriodStartCol == "" || n.PeriodEndCol == "" {
			return mssqlErr(ErrSystemVersioningNoPeriod, "Cannot set SYSTEM_VERSIONING to ON when SYSTEM_TIME period is not defined.")
		}
		if !createTableHasPrimaryKey(n) {
			return mssqlErr(ErrDropPKSystemVersioned, "System versioned temporal table '%s.%s' must have primary key defined.", schemaName, tableName)
		}
	}

	table := c.getOrCreateTable(schemaName, tableName)
	if n.PeriodStartCol != "" && n.PeriodEndCol != "" {
		if c.systemTimeTables == nil {
			c.systemTimeTables = make(map[tableKey][2]string)
		}
		periodKey := c.temporalKey(schemaName, table.Name)
		c.systemTimeTables[periodKey] = [2]string{n.PeriodStartCol, n.PeriodEndCol}
		if wantsSystemVersioning {
			if c.systemVersionedTables == nil {
				c.systemVersionedTables = make(map[tableKey]bool)
			}
			c.systemVersionedTables[periodKey] = true
		}
	}

	if n.Columns != nil {
		for _, item := range n.Columns.Items {
			column, ok := item.(*ast.ColumnDef)
			if !ok {
				continue
			}
			c.extractColumn(column, table, schemaName, source)
		}
	}
	if n.Constraints != nil {
		for _, item := range n.Constraints.Items {
			constraint, ok := item.(*ast.ConstraintDef)
			if !ok {
				continue
			}
			c.extractTableConstraint(constraint, table, schemaName, source)
		}
	}
	if n.Indexes != nil {
		for _, item := range n.Indexes.Items {
			index, ok := item.(*ast.InlineIndexDef)
			if !ok {
				continue
			}
			table.Indexes = append(table.Indexes, inlineIndexMetadata(index, source))
		}
	}
	for _, index := range table.Indexes {
		if !index.Primary {
			continue
		}
		for _, name := range index.Expressions {
			if column := columnByName(table, name); column != nil {
				column.Nullable = false
			}
		}
	}
	return nil
}

// --- CREATE / ALTER INDEX ---

func checkDuplicateIndexColumns(cols, includeCols *ast.List) error {
	seen := make(map[string]bool)
	for _, list := range []*ast.List{cols, includeCols} {
		if list == nil {
			continue
		}
		for _, item := range list.Items {
			var name string
			switch v := item.(type) {
			case *ast.IndexColumn:
				name = v.Name
			case *ast.String:
				name = v.Str
			}
			if name == "" {
				continue
			}
			low := strings.ToLower(name)
			if seen[low] {
				return mssqlErr(ErrDupIndexColumn, "Cannot use duplicate column names in index. Column name '%s' listed more than once.", name)
			}
			seen[low] = true
		}
	}
	return nil
}

func (c *Catalog) execCreateIndex(n *ast.CreateIndexStmt, source string) error {
	if err := checkDuplicateIndexColumns(n.Columns, n.IncludeCols); err != nil {
		return err
	}
	schemaName, tableName := tableRefSchemaObject(n.Table, c.defaultSchema)
	schema := c.GetSchema(schemaName)
	kind := c.lookupSchemaObject(schema, tableName, false)

	// Filtered indexes check object kind differently in SQL Server.
	if n.WhereClause != nil && kind != kindTable {
		if kind == kindView {
			if c.skipViewIndexes {
				return nil
			}
			return mssqlErr(ErrFilteredIndexOnView, "Filtered index '%s' cannot be created on object '%s' because it is not a user table.", n.Name, tableName)
		}
		return mssqlErr(ErrInvalidObjectName, "Invalid object name '%s'.", tableName)
	}

	switch kind {
	case kindNone:
		return mssqlErr(ErrObjectNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	case kindView:
		if c.skipViewIndexes {
			return nil
		}
		if n.Columnstore {
			if n.Clustered != nil && *n.Clustered {
				return mssqlErr(ErrClusteredCSOnView, "The statement failed because a clustered columnstore index cannot be created on a view.")
			}
			return mssqlErr(ErrIndexOnViewNotAllowed, "Cannot create index on view '%s'. It does not have a unique clustered index.", tableName)
		}
		view := c.GetView(schemaName, tableName)
		if view != nil && !isSchemaBound(view.Definition) {
			return mssqlErr(ErrViewNotSchemaBound, "Cannot create index on view '%s' because the view is not schema bound.", tableName)
		}
		return fmt.Errorf("indexes on views are not supported: %s.%s", schemaName, tableName)
	case kindTable:
		// proceed
	default:
		return mssqlErr(ErrIndexInvalidObject, "Index cannot be created on object '%s' because the object is not a user table or view.", tableName)
	}

	table := c.resolveTableInSchema(schema, tableName)
	if table == nil {
		return mssqlErr(ErrObjectNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	}
	if n.WhereClause != nil {
		if exprReferencesColumn(n.WhereClause, "IDENTITYCOL") && columnByName(table, "IDENTITYCOL") == nil && !tableHasIdentityColumn(table) {
			return mssqlErr(ErrInvalidColumnName, "Invalid column name 'IDENTITYCOL'.")
		}
	}

	dropExisting := hasOptionOn(n.Options, "DROP_EXISTING")
	var existingIdx *metadata.IndexMetadata
	existingPos := -1
	if n.Name != "" {
		existingIdx, existingPos = resolveTableIndex(table, n.Name, false)
	}
	if dropExisting {
		if existingIdx == nil {
			return mssqlErr(ErrIndexNotFound, "Could not find any index named '%s' for table '%s'.", n.Name, tableName)
		}
	} else if existingIdx != nil {
		return mssqlErr(ErrDupIndex, "The operation failed because an index or statistics with name '%s' already exists on table '%s'.", n.Name, tableName)
	}

	// Validate key, include, and order columns exist on the table.
	for _, list := range []*ast.List{n.Columns, n.IncludeCols, n.OrderCols} {
		if list == nil {
			continue
		}
		for _, item := range list.Items {
			var colName string
			switch v := item.(type) {
			case *ast.IndexColumn:
				colName = v.Name
			case *ast.String:
				colName = v.Str
			}
			if colName != "" && columnByName(table, colName) == nil {
				return mssqlErr(ErrIndexColumnNotFound, "Column name '%s' does not exist in the target table or view.", colName)
			}
		}
	}

	if !c.lenientMetadataLoad {
		if n.Columnstore && n.Clustered != nil && *n.Clustered {
			if tableHasClusteredIndexExcluding(table, existingIdx) {
				return mssqlErr(ErrClusteredCSExists, "You cannot create more than one clustered index on table '%s'.", tableName)
			}
		} else if n.Clustered != nil && *n.Clustered {
			if tableHasClusteredIndexExcluding(table, existingIdx) {
				return mssqlErr(ErrMultipleClusteredIndex, "Cannot create more than one clustered index on table '%s'.", tableName)
			}
		}
	}

	if !c.lenientMetadataLoad && !isValidFilegroup(n.OnFileGroup) {
		return mssqlErr(ErrInvalidFilegroup, "Invalid filegroup '%s' specified.", n.OnFileGroup)
	}

	index := &metadata.IndexMetadata{
		Name:         n.Name,
		Unique:       n.Unique,
		Type:         indexType(n.Columnstore, n.Clustered),
		Expressions:  []string{},
		Descending:   []bool{},
		IsConstraint: false,
	}
	appendIndexColumns(index, n.Columns)
	index.Definition = indexDefinition(n.IncludeCols, n.WhereClause, source)
	if dropExisting && existingPos >= 0 {
		if existingIdx.IsConstraint {
			if !index.Unique || n.Columnstore || n.WhereClause != nil ||
				(n.IncludeCols != nil && len(n.IncludeCols.Items) > 0) ||
				!sameIndexKeyColumns(existingIdx.Expressions, index.Expressions) {
				return mssqlErr(ErrDropExistingConstraintMismatch, "Cannot recreate index '%s'. The new index definition does not match the constraint being enforced by the existing index.", existingIdx.Name)
			}
			if existingIdx.Primary && isClusteredPrimaryKey(existingIdx) && !strings.EqualFold(index.Type, "CLUSTERED") && tableHasXMLIndex(table) {
				return mssqlErr(ErrDropPKWithXMLIndex, "Cannot drop the clustered primary key '%s' on table '%s' because the table has an XML index.", existingIdx.Name, table.Name)
			}
		} else if _, fkName, referenced := indexReferencedByFK(c, schemaName, table, existingIdx); referenced {
			if !index.Unique || n.Columnstore || n.WhereClause != nil || !sameIndexKeyColumns(existingIdx.Expressions, index.Expressions) {
				return mssqlErr(ErrDropIndexReferencedByFK, "An explicit DROP INDEX is not allowed on index '%s.%s'. It is being used for FOREIGN KEY constraint enforcement (%s).", tableName, existingIdx.Name, fkName)
			}
		}
		index.Primary = existingIdx.Primary
		index.IsConstraint = existingIdx.IsConstraint
		table.Indexes[existingPos] = index
	} else {
		table.Indexes = append(table.Indexes, index)
	}
	return nil
}

func sameIndexKeyColumns(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !strings.EqualFold(a[i], b[i]) {
			return false
		}
	}
	return true
}

func (c *Catalog) findTransferredTableCollision(targetSchema *metadata.SchemaMetadata, table *metadata.TableMetadata) string {
	for _, idx := range table.Indexes {
		if idx.IsConstraint && idx.Name != "" && c.lookupSchemaObject(targetSchema, idx.Name, false) != kindNone {
			return idx.Name
		}
	}
	for _, fk := range table.ForeignKeys {
		if fk.Name != "" && c.lookupSchemaObject(targetSchema, fk.Name, false) != kindNone {
			return fk.Name
		}
	}
	for _, ck := range table.CheckConstraints {
		if ck.Name != "" && c.lookupSchemaObject(targetSchema, ck.Name, false) != kindNone {
			return ck.Name
		}
	}
	for _, col := range table.Columns {
		if col.DefaultConstraintName != "" && c.lookupSchemaObject(targetSchema, col.DefaultConstraintName, false) != kindNone {
			return col.DefaultConstraintName
		}
	}
	for _, tr := range table.Triggers {
		if tr.Name != "" && c.lookupSchemaObject(targetSchema, tr.Name, false) != kindNone {
			return tr.Name
		}
	}
	return ""
}

func (c *Catalog) execCreateXMLIndex(n *ast.CreateXmlIndexStmt) error {
	schemaName, tableName := tableRefSchemaObject(n.Table, c.defaultSchema)
	schema := c.GetSchema(schemaName)
	kind := c.lookupSchemaObject(schema, tableName, false)
	switch kind {
	case kindNone:
		return mssqlErr(ErrObjectNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	case kindTable:
	default:
		return mssqlErr(ErrXmlIndexNotTable, "Could not create the XML or spatial index on object '%s' because that object is not a table.", tableName)
	}
	table := c.resolveTableInSchema(schema, tableName)
	if table == nil {
		return mssqlErr(ErrObjectNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	}
	col := columnByName(table, n.XmlColumn)
	if col == nil {
		return mssqlErr(ErrIndexColumnNotFound, "Column name '%s' does not exist in the target table or view.", n.XmlColumn)
	}
	if !n.Primary {
		parent := indexByName(table, n.UsingIndex)
		if parent == nil || !strings.EqualFold(parent.Type, "XML") {
			return mssqlErr(ErrXmlIndexMissingPrimary, "Could not find SELECTIVE XML index named '%s' on table '%s'.", n.UsingIndex, tableName)
		}
	}
	if n.Name != "" && indexByName(table, n.Name) != nil {
		return mssqlErr(ErrDupIndex, "The operation failed because an index or statistics with name '%s' already exists on table '%s'.", n.Name, tableName)
	}
	if normalizeType(col.Type) != "xml" {
		return mssqlErr(ErrXmlIndexNonXmlCol, "Column '%s' on table '%s' is not of type XML.", n.XmlColumn, tableName)
	}

	index := &metadata.IndexMetadata{
		Name:         n.Name,
		Type:         "XML",
		Expressions:  []string{n.XmlColumn},
		Descending:   []bool{false},
		IsConstraint: false,
		Definition:   "PRIMARY",
	}
	if !n.Primary {
		index.Definition = "USING XML INDEX " + n.UsingIndex
		if n.SecondaryFor != "" {
			index.Definition += " FOR " + n.SecondaryFor
		}
	}
	table.Indexes = append(table.Indexes, index)
	return nil
}

func (c *Catalog) execCreateSelectiveXMLIndex(n *ast.CreateSelectiveXmlIndexStmt) error {
	schemaName, tableName := tableRefSchemaObject(n.Table, c.defaultSchema)
	schema := c.GetSchema(schemaName)
	kind := c.lookupSchemaObject(schema, tableName, false)
	switch kind {
	case kindNone:
		return mssqlErr(ErrObjectNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	case kindTable:
	default:
		return mssqlErr(ErrXmlIndexNotTable, "Could not create the XML or spatial index on object '%s' because that object is not a table.", tableName)
	}
	table := c.resolveTableInSchema(schema, tableName)
	if table == nil {
		return mssqlErr(ErrObjectNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	}
	col := columnByName(table, n.XmlColumn)
	if col == nil {
		return mssqlErr(ErrIndexColumnNotFound, "Column name '%s' does not exist in the target table or view.", n.XmlColumn)
	}
	if n.Name != "" && indexByName(table, n.Name) != nil {
		return mssqlErr(ErrDupIndex, "The operation failed because an index or statistics with name '%s' already exists on table '%s'.", n.Name, tableName)
	}
	if normalizeType(col.Type) != "xml" {
		return mssqlErr(ErrXmlIndexNonXmlCol, "Column '%s' on table '%s' is not of type XML.", n.XmlColumn, tableName)
	}
	table.Indexes = append(table.Indexes, &metadata.IndexMetadata{
		Name:        n.Name,
		Type:        "XML",
		Expressions: []string{n.XmlColumn},
		Descending:  []bool{false},
		Definition:  selectiveXMLIndexDefinition(n),
	})
	return nil
}

func selectiveXMLIndexDefinition(n *ast.CreateSelectiveXmlIndexStmt) string {
	var parts []string
	if n.Namespaces != nil && len(n.Namespaces.Items) > 0 {
		var ns []string
		for _, item := range n.Namespaces.Items {
			decl, ok := item.(*ast.XmlNamespaceDecl)
			if !ok {
				continue
			}
			uri := "'" + strings.ReplaceAll(decl.URI, "'", "''") + "'"
			if decl.IsDefault {
				ns = append(ns, "DEFAULT "+uri)
			} else {
				ns = append(ns, uri+" AS "+quote(decl.Prefix))
			}
		}
		if len(ns) > 0 {
			parts = append(parts, "WITH XMLNAMESPACES ("+strings.Join(ns, ", ")+")")
		}
	}
	if n.Paths != nil && len(n.Paths.Items) > 0 {
		var paths []string
		for _, item := range n.Paths.Items {
			if s, ok := item.(*ast.String); ok && strings.TrimSpace(s.Str) != "" {
				paths = append(paths, strings.TrimSpace(s.Str))
			}
		}
		if len(paths) > 0 {
			parts = append(parts, "FOR ("+strings.Join(paths, ", ")+")")
		}
	}
	if len(parts) == 0 {
		return "SELECTIVE"
	}
	return "SELECTIVE " + strings.Join(parts, " ")
}

func (c *Catalog) execCreateSpatialIndex(n *ast.CreateSpatialIndexStmt) error {
	schemaName, tableName := tableRefSchemaObject(n.Table, c.defaultSchema)
	schema := c.GetSchema(schemaName)
	kind := c.lookupSchemaObject(schema, tableName, false)
	switch kind {
	case kindNone:
		return mssqlErr(ErrObjectNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	case kindTable:
	default:
		return mssqlErr(ErrXmlIndexNotTable, "Could not create the XML or spatial index on object '%s' because that object is not a table.", tableName)
	}
	table := c.resolveTableInSchema(schema, tableName)
	if table == nil {
		return mssqlErr(ErrObjectNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	}
	col := columnByName(table, n.SpatialColumn)
	if col == nil {
		return mssqlErr(ErrIndexColumnNotFound, "Column name '%s' does not exist in the target table or view.", n.SpatialColumn)
	}
	colType := normalizeType(col.Type)
	if colType != "geometry" && colType != "geography" {
		return mssqlErr(ErrSpatialIndexNonSpatialCol, "The requested spatial index on column '%s' of table '%s' could not be created because the column type is not geometry or geography.", n.SpatialColumn, tableName)
	}
	if n.Name != "" && indexByName(table, n.Name) != nil {
		return mssqlErr(ErrDupIndex, "The operation failed because an index or statistics with name '%s' already exists on table '%s'.", n.Name, tableName)
	}

	index := &metadata.IndexMetadata{
		Name:         n.Name,
		Type:         "SPATIAL",
		Expressions:  []string{},
		Descending:   []bool{},
		IsConstraint: false,
		SpatialConfig: &metadata.SpatialIndexConfig{
			Method: "SPATIAL",
			Tessellation: &metadata.TessellationConfig{
				Scheme: strings.ToUpper(n.Using),
			},
			Dimensional: &metadata.DimensionalConfig{
				Dimensions: 2,
			},
		},
	}
	if n.SpatialColumn != "" {
		index.Expressions = append(index.Expressions, n.SpatialColumn)
		index.Descending = append(index.Descending, false)
	}
	if strings.Contains(index.SpatialConfig.Tessellation.Scheme, "GEOGRAPHY") {
		index.SpatialConfig.Dimensional.DataType = "GEOGRAPHY"
	} else {
		index.SpatialConfig.Dimensional.DataType = "GEOMETRY"
	}
	parseSpatialOptions(n.Options, index)
	table.Indexes = append(table.Indexes, index)
	return nil
}

func (c *Catalog) execAlterIndex(n *ast.AlterIndexStmt) error {
	targetCat := c
	if n.Table != nil && n.Table.Database != "" {
		targetCat = c.DatabaseCatalog(n.Table.Database)
		if targetCat == nil {
			return mssqlErr(ErrObjectNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", n.Table.Object)
		}
	}
	schemaName, tableName := tableRefSchemaObject(n.Table, targetCat.defaultSchema)
	schema := targetCat.GetSchema(schemaName)
	kind := targetCat.lookupSchemaObject(schema, tableName, false)
	switch kind {
	case kindNone:
		return mssqlErr(ErrObjectNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	case kindTable, kindView:
	default:
		return mssqlErr(ErrIndexInvalidObject, "Index cannot be altered on object '%s' because the object is not a user table or view.", tableName)
	}
	var table *metadata.TableMetadata
	if kind == kindTable {
		table = targetCat.resolveTableInSchema(schema, tableName)
		if table == nil {
			return mssqlErr(ErrObjectNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
		}
	}
	if strings.EqualFold(n.IndexName, "ALL") {
		return nil
	}
	var idx *metadata.IndexMetadata
	if table != nil {
		idx, _ = resolveTableIndex(table, n.IndexName, false)
	}
	if idx == nil {
		return mssqlErr(ErrIndexNotFound, "Cannot find index '%s'.", n.IndexName)
	}
	switch strings.ToUpper(n.Action) {
	case "FOR":
		if !strings.EqualFold(idx.Type, "XML") {
			return mssqlErr(ErrXmlIndexMissingPrimary, "Could not find SELECTIVE XML index named '%s' on table '%s'.", n.IndexName, tableName)
		}
	case "PAUSE", "ABORT":
		return mssqlErr(ErrNoResumableIndexOp, "ALTER INDEX '%s' failed. There is no pending resumable index operation for the index '%s' on '%s'.", strings.ToUpper(n.Action), n.IndexName, tableName)
	}
	return nil
}

// --- CREATE / ALTER VIEW, PROCEDURE, FUNCTION, TRIGGER ---

func (c *Catalog) execCreateOrAlterView(n *ast.CreateViewStmt, source string) error {
	if selectStarWithoutFrom(n.Query) {
		return mssqlErr(ErrSelectStarWithoutFrom, "Must specify table to select from.")
	}
	schemaName, name := tableRefSchemaObject(n.Name, c.defaultSchema)
	schema := c.GetSchema(schemaName)
	kind := c.lookupSchemaObject(schema, name, false)
	pureAlter := isPureAlterStmt(n, source)

	switch {
	case pureAlter:
		if kind == kindNone {
			return mssqlErr(ErrInvalidObjectName, "Invalid object name '%s'.", name)
		}
		if kind != kindView {
			return mssqlErr(ErrIncompatibleAlterType, "Cannot perform alter on '%s' because it is an incompatible object type.", name)
		}
	case n.OrAlter:
		if kind != kindNone && kind != kindView {
			return mssqlErr(ErrIncompatibleAlterType, "Cannot perform alter on '%s' because it is an incompatible object type.", name)
		}
	default:
		if kind != kindNone {
			return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", name)
		}
	}

	schema = c.getOrCreateSchema(schemaName)
	definition := definitionText(n, source)
	if kind == kindView {
		if idx, ok := resolveExactOrFoldIndex(schema.Views, func(v *metadata.ViewMetadata) string { return v.Name }, name); ok {
			v := schema.Views[idx]
			if consumer, bound := schemaBoundObjectConsumer(c, schema.Name, v.Name, nil); bound {
				return mssqlErr(ErrDropSchemaNotEmpty, "Cannot ALTER VIEW '%s' because it is being referenced by object '%s'.", v.Name, consumer)
			}
			v.Definition = creationDefinition(definition)
			return nil
		}
		if pureAlter {
			return mssqlErr(ErrInvalidObjectName, "Invalid object name '%s'.", name)
		}
	}
	schema.Views = append(schema.Views, &metadata.ViewMetadata{Name: name, Definition: definition})
	return nil
}

func (c *Catalog) execCreateOrAlterProcedure(n *ast.CreateProcedureStmt, source string) error {
	if !c.lenientMetadataLoad && n.Params != nil {
		for idx, item := range n.Params.Items {
			p, ok := item.(*ast.ParamDef)
			if !ok {
				continue
			}
			if !isValidDataType(p.DataType, true) {
				typeName := ""
				if p.DataType != nil {
					typeName = p.DataType.Name
				}
				return mssqlErr(ErrUnknownColumnType, "Column, parameter, or variable #%d: Cannot find data type %s.", idx+1, typeName)
			}
		}
	}
	if !c.lenientMetadataLoad && n.ExternalName != "" {
		return mssqlErr(ErrAssemblyNotFound, "Assembly '%s' was not found in the SQL catalog.", n.ExternalName)
	}
	schemaName, name := tableRefSchemaObject(n.Name, c.defaultSchema)
	schema := c.GetSchema(schemaName)
	kind := c.lookupSchemaObject(schema, name, false)
	pureAlter := isPureAlterStmt(n, source)

	switch {
	case pureAlter:
		if kind == kindNone {
			return mssqlErr(ErrInvalidObjectName, "Invalid object name '%s'.", name)
		}
		if kind != kindProcedure {
			return mssqlErr(ErrIncompatibleAlterType, "Cannot perform alter on '%s' because it is an incompatible object type.", name)
		}
	case n.OrAlter:
		if kind != kindNone && kind != kindProcedure {
			return mssqlErr(ErrIncompatibleAlterType, "Cannot perform alter on '%s' because it is an incompatible object type.", name)
		}
	default:
		if kind != kindNone {
			return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", name)
		}
	}

	if !c.lenientMetadataLoad {
		if err := checkExecuteAsOptions(n.Options); err != nil {
			return err
		}
	}

	schema = c.getOrCreateSchema(schemaName)
	definition := definitionText(n, source)
	if n.OrAlter {
		definition = normalizeProcDefinition(definition)
	}
	if kind == kindProcedure {
		if idx, ok := resolveExactOrFoldIndex(schema.Procedures, func(p *metadata.ProcedureMetadata) string { return p.Name }, name); ok {
			schema.Procedures[idx].Definition = definition
			return nil
		}
		if pureAlter {
			return mssqlErr(ErrInvalidObjectName, "Invalid object name '%s'.", name)
		}
	}
	schema.Procedures = append(schema.Procedures, &metadata.ProcedureMetadata{Name: name, Definition: definition})
	return nil
}

// normalizeProcDefinition mirrors SQL Server's OBJECT_DEFINITION storage:
// CREATE OR ALTER PROC ... is stored as CREATE PROC ... and ALTER PROC ... as
// CREATE PROC ... (or keeps the normalized CREATE form for diffing).
func normalizeProcDefinition(def string) string {
	trimmed := strings.TrimSpace(def)
	upper := strings.ToUpper(trimmed)
	if strings.HasPrefix(upper, "CREATE OR ALTER ") {
		return "CREATE " + strings.TrimSpace(trimmed[len("CREATE OR ALTER "):])
	}
	if strings.HasPrefix(upper, "ALTER ") {
		return "CREATE " + strings.TrimSpace(trimmed[len("ALTER "):])
	}
	return trimmed
}

func (c *Catalog) execCreateOrAlterFunction(n *ast.CreateFunctionStmt, source string) error {
	if !c.lenientMetadataLoad && n.ExternalName != "" {
		return mssqlErr(ErrSecurableNotFound, "Cannot find the assembly '%s', because it does not exist or you do not have permission.", n.ExternalName)
	}
	if !c.lenientMetadataLoad && n.Params != nil {
		for idx, item := range n.Params.Items {
			p, ok := item.(*ast.ParamDef)
			if !ok {
				continue
			}
			if !isValidDataType(p.DataType, false) {
				typeName := ""
				if p.DataType != nil {
					typeName = p.DataType.Name
				}
				return mssqlErr(ErrUnknownColumnType, "Column, parameter, or variable #%d: Cannot find data type %s.", idx+1, typeName)
			}
		}
	}
	if !c.lenientMetadataLoad && n.ReturnType != nil && !isValidDataType(n.ReturnType, false) {
		return mssqlErr(ErrUnknownColumnType, "Column, parameter, or variable #0: Cannot find data type %s.", n.ReturnType.Name)
	}
	if n.ReturnsTable != nil && selectStarWithoutFrom(n.Body) {
		return mssqlErr(ErrSelectStarWithoutFrom, "Must specify table to select from.")
	}

	schemaName, name := tableRefSchemaObject(n.Name, c.defaultSchema)
	schema := c.GetSchema(schemaName)
	kind := c.lookupSchemaObject(schema, name, false)
	pureAlter := isPureAlterStmt(n, source)

	switch {
	case pureAlter:
		if kind == kindNone {
			return mssqlErr(ErrInvalidObjectName, "Invalid object name '%s'.", name)
		}
		if kind != kindFunction {
			return mssqlErr(ErrIncompatibleAlterType, "Cannot perform alter on '%s' because it is an incompatible object type.", name)
		}
	case n.OrAlter:
		if kind != kindNone && kind != kindFunction {
			return mssqlErr(ErrIncompatibleAlterType, "Cannot perform alter on '%s' because it is an incompatible object type.", name)
		}
	default:
		if kind != kindNone {
			return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", name)
		}
	}

	if !c.lenientMetadataLoad {
		if err := checkExecuteAsOptions(n.Options); err != nil {
			return err
		}
	}

	schema = c.getOrCreateSchema(schemaName)
	definition := definitionText(n, source)
	if kind == kindFunction {
		if idx, ok := resolveExactOrFoldIndex(schema.Functions, func(fn *metadata.FunctionMetadata) string { return fn.Name }, name); ok {
			fn := schema.Functions[idx]
			if consumer, bound := schemaBoundFunctionAlterConsumer(c, schema.Name, fn.Name); bound {
				return mssqlErr(ErrDropSchemaNotEmpty, "Cannot ALTER FUNCTION '%s' because it is being referenced by object '%s'.", fn.Name, consumer)
			}
			oldCat := definitionFunctionCategory(fn.Definition)
			newCat := functionCategory(n)
			if oldCat != "" && newCat != "" && oldCat != newCat {
				return mssqlErr(ErrIncompatibleAlterType, "Cannot perform alter on '%s' because it is an incompatible object type.", fn.Name)
			}
			fn.Definition = creationDefinition(definition)
			return nil
		}
		if pureAlter {
			return mssqlErr(ErrInvalidObjectName, "Invalid object name '%s'.", name)
		}
	}
	schema.Functions = append(schema.Functions, &metadata.FunctionMetadata{Name: name, Definition: definition})
	return nil
}

func functionCategory(n *ast.CreateFunctionStmt) string {
	if n == nil {
		return ""
	}
	if n.ReturnsTable == nil {
		return "scalar"
	}
	if n.ReturnsTable.Variable != "" {
		return "mstvf"
	}
	return "itvf"
}

func definitionFunctionCategory(def string) string {
	stmts, err := mssql.Parse(creationDefinition(def))
	if err != nil || len(stmts) == 0 {
		return ""
	}
	if fn, ok := stmts[0].AST.(*ast.CreateFunctionStmt); ok {
		return functionCategory(fn)
	}
	return ""
}

func (c *Catalog) execCreateOrAlterTrigger(n *ast.CreateTriggerStmt, source string) error {
	if n.Table == nil || n.Name == nil {
		return nil
	}
	schemaName, parentName := tableRefSchemaObject(n.Table, c.defaultSchema)
	schema := c.GetSchema(schemaName)
	parentKind := c.lookupSchemaObject(schema, parentName, false)
	switch parentKind {
	case kindTable:
	case kindView:
		if !strings.EqualFold(n.TriggerType, "INSTEAD OF") {
			return mssqlErr(ErrTriggerTargetInvalid, "The object '%s' does not exist or is invalid for this operation.", parentName)
		}
	default:
		return mssqlErr(ErrTriggerTargetInvalid, "The object '%s' does not exist or is invalid for this operation.", parentName)
	}

	trigSchemaName, trigName := tableRefSchemaObject(n.Name, schemaName)
	trigSchema := c.GetSchema(trigSchemaName)
	existingKind := c.lookupSchemaObject(trigSchema, trigName, false)
	if !n.OrAlter && existingKind != kindNone {
		return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", trigName)
	}

	var events []string
	if n.Events != nil {
		for _, item := range n.Events.Items {
			switch event := item.(type) {
			case *ast.String:
				events = append(events, event.Str)
			case *ast.TriggerEvent:
				events = append(events, event.Name)
			}
		}
	}
	timing := n.TriggerType
	if strings.EqualFold(timing, "FOR") {
		timing = "AFTER"
	}
	trigger := &metadata.TriggerMetadata{
		Name:   trigName,
		Timing: timing,
		Event:  strings.Join(events, ","),
		Body:   definitionText(n, source),
	}
	if parentKind == kindTable {
		table := c.resolveTableInSchema(schema, parentName)
		if table == nil {
			return mssqlErr(ErrTriggerTargetInvalid, "The object '%s' does not exist or is invalid for this operation.", parentName)
		}
		table.Triggers = append(table.Triggers, trigger)
		return nil
	}
	if view := c.GetView(schemaName, parentName); view != nil {
		view.Triggers = append(view.Triggers, trigger)
	}
	return nil
}

// --- CREATE / ALTER SEQUENCE ---

func (c *Catalog) execCreateSequence(n *ast.CreateSequenceStmt, source string) error {
	if !c.lenientMetadataLoad && n.DataType != nil && !isValidSequenceType(n.DataType) {
		return mssqlErr(ErrUnknownType, "Type '%s' is not a defined system type.", n.DataType.Name)
	}
	schemaName, sequenceName := tableRefSchemaObject(n.Name, c.defaultSchema)
	if isZeroNumericExpr(n.Increment) {
		return mssqlErr(ErrSequenceZeroIncrement, "The increment value for sequence object '%s' cannot be zero.", sequenceName)
	}
	if schema := c.GetSchema(schemaName); schema != nil {
		if c.lookupSchemaObject(schema, sequenceName, false) != kindNone {
			return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", sequenceName)
		}
	}
	schema := c.getOrCreateSchema(schemaName)
	sequence := &metadata.SequenceMetadata{Name: sequenceName}
	if n.DataType != nil {
		sequence.DataType = nodeText(n.DataType, source)
	}
	if n.Start != nil {
		sequence.Start = strings.TrimSpace(nodeText(n.Start, source))
	}
	if n.Increment != nil {
		sequence.Increment = strings.TrimSpace(nodeText(n.Increment, source))
	}
	if n.MinValue != nil {
		sequence.MinValue = strings.TrimSpace(nodeText(n.MinValue, source))
	}
	if n.MaxValue != nil {
		sequence.MaxValue = strings.TrimSpace(nodeText(n.MaxValue, source))
	}
	if n.Cycle != nil {
		sequence.Cycle = *n.Cycle
	}
	if n.NoCache {
		sequence.CacheSize = "NO CACHE"
	} else if n.Cache != nil {
		sequence.CacheSize = strings.TrimSpace(nodeText(n.Cache, source))
	}
	if !c.lenientMetadataLoad {
		if err := validateSequenceBounds(sequence, sequenceName, true, ""); err != nil {
			return err
		}
	}
	schema.Sequences = append(schema.Sequences, sequence)
	return nil
}

func (c *Catalog) execAlterSequence(n *ast.AlterSequenceStmt, source string) error {
	targetCat := c
	if n.Name != nil && n.Name.Database != "" {
		targetCat = c.DatabaseCatalog(n.Name.Database)
		if targetCat == nil {
			return mssqlErr(ErrSecurableNotFound, "Cannot find the object '%s', because it does not exist or you do not have permission.", n.Name.Object)
		}
	}
	schemaName, sequenceName := tableRefSchemaObject(n.Name, targetCat.defaultSchema)
	schema := targetCat.GetSchema(schemaName)
	var seq *metadata.SequenceMetadata
	if schema != nil {
		if idx, ok := resolveExactOrFoldIndex(schema.Sequences, func(sq *metadata.SequenceMetadata) string { return sq.Name }, sequenceName); ok {
			seq = schema.Sequences[idx]
		}
	}
	if seq == nil {
		return mssqlErr(ErrSecurableNotFound, "Cannot find the object '%s', because it does not exist or you do not have permission.", sequenceName)
	}
	if isZeroNumericExpr(n.Increment) {
		return mssqlErr(ErrSequenceZeroIncrement, "The increment value for sequence object '%s' cannot be zero.", sequenceName)
	}
	staged := &metadata.SequenceMetadata{
		Name:        seq.Name,
		DataType:    seq.DataType,
		Start:       seq.Start,
		Increment:   seq.Increment,
		MinValue:    seq.MinValue,
		MaxValue:    seq.MaxValue,
		Cycle:       seq.Cycle,
		CacheSize:   seq.CacheSize,
		OwnerTable:  seq.OwnerTable,
		OwnerColumn: seq.OwnerColumn,
		Comment:     seq.Comment,
	}
	if strings.TrimSpace(staged.Start) == "" && (n.Increment != nil || n.NoMinVal || n.MinValue != nil || n.NoMaxVal || n.MaxValue != nil) {
		staged.Start = normalizeSequenceMetadata(staged).start
	}
	if n.Increment != nil {
		staged.Increment = strings.TrimSpace(nodeText(n.Increment, source))
	}
	if n.NoMinVal {
		staged.MinValue = ""
	} else if n.MinValue != nil {
		staged.MinValue = strings.TrimSpace(nodeText(n.MinValue, source))
	}
	if n.NoMaxVal {
		staged.MaxValue = ""
	} else if n.MaxValue != nil {
		staged.MaxValue = strings.TrimSpace(nodeText(n.MaxValue, source))
	}
	if n.Cycle != nil {
		staged.Cycle = *n.Cycle
	}
	if n.NoCache {
		staged.CacheSize = "NO CACHE"
	} else if n.Cache != nil {
		staged.CacheSize = strings.TrimSpace(nodeText(n.Cache, source))
	} else if n.HasCache {
		staged.CacheSize = ""
	}
	restartWith := ""
	if n.RestartWith != nil {
		restartWith = strings.TrimSpace(nodeText(n.RestartWith, source))
	}
	if err := validateSequenceBounds(staged, sequenceName, false, restartWith); err != nil {
		return err
	}
	seq.Start = staged.Start
	seq.Increment = staged.Increment
	seq.MinValue = staged.MinValue
	seq.MaxValue = staged.MaxValue
	seq.Cycle = staged.Cycle
	seq.CacheSize = staged.CacheSize
	return nil
}

func parseSequenceBigInt(raw string) (*big.Int, bool) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return nil, false
	}
	s = strings.TrimPrefix(s, "+")
	v, ok := new(big.Int).SetString(s, 10)
	return v, ok
}

func validateSequenceBounds(seq *metadata.SequenceMetadata, seqName string, checkStart bool, restartWith string) error {
	norm := normalizeSequenceMetadata(seq)
	if inc, ok := parseSequenceBigInt(norm.increment); ok && inc.Sign() == 0 {
		return mssqlErr(ErrSequenceZeroIncrement, "The increment value for sequence object '%s' cannot be zero.", seqName)
	}
	minVal, hasMin := parseSequenceBigInt(norm.minValue)
	maxVal, hasMax := parseSequenceBigInt(norm.maxValue)
	typeMinStr, typeMaxStr := sequenceTypeBounds(norm.dataType)
	if typeMin, ok := parseSequenceBigInt(typeMinStr); ok && hasMin && minVal.Cmp(typeMin) < 0 {
		return mssqlErr(ErrSequenceInvalidBounds, "The minimum value for sequence object '%s' is out of range for its data type.", seqName)
	}
	if typeMax, ok := parseSequenceBigInt(typeMaxStr); ok && hasMax && maxVal.Cmp(typeMax) > 0 {
		return mssqlErr(ErrSequenceInvalidBounds, "The maximum value for sequence object '%s' is out of range for its data type.", seqName)
	}
	if hasMin && hasMax && minVal.Cmp(maxVal) >= 0 {
		return mssqlErr(ErrSequenceMinGreaterThanMax, "The minimum value for sequence object '%s' must be less than its maximum value.", seqName)
	}
	if checkStart {
		if startVal, ok := parseSequenceBigInt(norm.start); ok {
			if (hasMin && startVal.Cmp(minVal) < 0) || (hasMax && startVal.Cmp(maxVal) > 0) {
				return mssqlErr(ErrSequenceStartOutOfRange, "The start value for sequence object '%s' must be between the minimum and maximum value of the sequence object.", seqName)
			}
		}
	}
	if restartWith != "" {
		if rv, ok := parseSequenceBigInt(restartWith); ok {
			if (hasMin && rv.Cmp(minVal) < 0) || (hasMax && rv.Cmp(maxVal) > 0) {
				return mssqlErr(ErrSequenceStartOutOfRange, "The start value for sequence object '%s' must be between the minimum and maximum value of the sequence object.", seqName)
			}
		}
	}
	return nil
}

// --- CREATE SYNONYM ---

func (c *Catalog) execCreateSynonym(n *ast.CreateSynonymStmt) error {
	if n.Name == nil || n.Target == nil {
		return nil
	}
	schemaName, name := tableRefSchemaObject(n.Name, c.defaultSchema)
	if schema := c.GetSchema(schemaName); schema != nil {
		if c.lookupSchemaObject(schema, name, false) != kindNone {
			return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", name)
		}
	}
	schema := c.getOrCreateSchema(schemaName)
	schema.Synonyms = append(schema.Synonyms, &metadata.SynonymMetadata{
		Name:           name,
		TargetServer:   n.Target.Server,
		TargetDatabase: n.Target.Database,
		TargetSchema:   n.Target.Schema,
		TargetName:     n.Target.Object,
	})
	return nil
}

// --- ALTER TABLE ---

func (c *Catalog) execAlterTable(n *ast.AlterTableStmt, source string) error {
	if n.Actions == nil {
		return nil
	}
	// Pre-lookup compile-time checks SQL Server performs before resolving the
	// table name.
	if err := checkAlterTablePreLookup(n); err != nil {
		return err
	}

	if n.Name != nil && n.Name.Database != "" {
		targetCat := c.DatabaseCatalog(n.Name.Database)
		if targetCat == nil {
			return mssqlErr(ErrAlterTableNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", n.Name.Object)
		}
		if targetCat != c {
			return targetCat.execAlterTableInCatalog(n, source)
		}
	}
	return c.execAlterTableInCatalog(n, source)
}

func (c *Catalog) execAlterTableInCatalog(n *ast.AlterTableStmt, source string) error {
	schemaName, tableName := tableRefSchemaObject(n.Name, c.defaultSchema)
	schema := c.GetSchema(schemaName)
	kind := c.lookupSchemaObject(schema, tableName, false)
	switch kind {
	case kindNone:
		return mssqlErr(ErrAlterTableNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	case kindTable:
	default:
		return mssqlErr(ErrAlterTableNotTable, "Cannot alter '%s' because it is not a table.", tableName)
	}
	table := c.resolveTableInSchema(schema, tableName)
	if table == nil {
		return mssqlErr(ErrAlterTableNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	}

	// Stage actions against a clone of the table so a failing later action
	// leaves the catalog untouched, matching SQL Server statement atomicity.
	tableIdx := -1
	for i, t := range schema.Tables {
		if t == table {
			tableIdx = i
			break
		}
	}
	staged := cloneTable(table)
	savedIndexCounter := c.indexCounter
	key := tableKey{schema: schema.Name, table: table.Name}
	periodKey := tableKey{schema: normalizeOrDefault(schema.Name, c.defaultSchema), table: table.Name}
	savedPeriodCols, savedSystemTime := c.systemTimeTables[periodKey]
	savedSystemVersioned := c.systemVersionedTables[periodKey]
	if tableIdx >= 0 {
		schema.Tables[tableIdx] = staged
		c.tableByName[key] = staged
	}

	for _, item := range n.Actions.Items {
		action, ok := item.(*ast.AlterTableAction)
		if !ok {
			continue
		}
		if err := c.execAlterTableAction(staged, schema.Name, action, source); err != nil {
			if tableIdx >= 0 {
				schema.Tables[tableIdx] = table
				c.tableByName[key] = table
			}
			c.indexCounter = savedIndexCounter
			if savedSystemTime {
				if c.systemTimeTables == nil {
					c.systemTimeTables = make(map[tableKey][2]string)
				}
				c.systemTimeTables[periodKey] = savedPeriodCols
			} else {
				delete(c.systemTimeTables, periodKey)
			}
			if savedSystemVersioned {
				if c.systemVersionedTables == nil {
					c.systemVersionedTables = make(map[tableKey]bool)
				}
				c.systemVersionedTables[periodKey] = true
			} else {
				delete(c.systemVersionedTables, periodKey)
			}
			return err
		}
	}
	return nil
}

func checkAlterTablePreLookup(n *ast.AlterTableStmt) error {
	seenDrop := make(map[string]bool)
	for _, item := range n.Actions.Items {
		action, ok := item.(*ast.AlterTableAction)
		if !ok {
			continue
		}
		switch action.Type {
		case ast.ATAddIndex:
			if action.Index == nil || action.Index.Columns == nil || len(action.Index.Columns.Items) == 0 {
				return mssqlErr(ErrAlterTableAddIndexSyntax, "Table level constraint or index does not specify column list.")
			}
		case ast.ATDropConstraint:
			if action.Constraint != nil && action.Constraint.Name != "" {
				low := strings.ToLower(action.Constraint.Name)
				if seenDrop[low] {
					return mssqlErr(ErrDupTargetInAlter, "Cannot create, drop, enable, or disable more than one constraint, column, index, or trigger named '%s' in this context.", action.Constraint.Name)
				}
				seenDrop[low] = true
			}
		case ast.ATCheckConstraint, ast.ATNocheckConstraint:
			if action.Names != nil {
				seen := make(map[string]bool)
				for _, ni := range action.Names.Items {
					if s, ok := ni.(*ast.String); ok {
						low := strings.ToLower(s.Str)
						if seen[low] {
							return mssqlErr(ErrDupTargetInAlter, "Cannot create, drop, enable, or disable more than one constraint, column, index, or trigger named '%s' in this context.", s.Str)
						}
						seen[low] = true
					}
				}
			}
		}
	}
	return nil
}

func (c *Catalog) execAlterTableAction(table *metadata.TableMetadata, schemaName string, action *ast.AlterTableAction, source string) error {
	switch action.Type {
	case ast.ATAddColumn:
		if action.Column == nil {
			return nil
		}
		if columnByName(table, action.Column.Name) != nil {
			return mssqlErr(ErrDupColumn, "Column names in each table must be unique. Column name '%s' in table '%s' is specified more than once.", action.Column.Name, table.Name)
		}
		if !c.lenientMetadataLoad && action.Column.Computed == nil && !isValidDataType(action.Column.DataType, false) {
			typeName := ""
			if action.Column.DataType != nil {
				typeName = action.Column.DataType.Name
			}
			return mssqlErr(ErrUnknownColumnType, "Column, parameter, or variable #1: Cannot find data type %s.", typeName)
		}
		if action.Column.Identity != nil && tableHasIdentityColumn(table) {
			return mssqlErr(ErrMultipleIdentityColumns, "Multiple identity columns specified for table '%s'. Only one identity column per table is allowed.", table.Name)
		}
		if !c.lenientMetadataLoad && action.Column.Identity != nil {
			explicitNull := action.Column.Nullable != nil && !action.Column.Nullable.NotNull
			if !isValidIdentityDataType(action.Column.DataType) || explicitNull {
				return mssqlErr(ErrIdentityColumnType, "Identity column '%s' must be of data type int, bigint, smallint, tinyint, or decimal or numeric with a scale of 0, unencrypted, and constrained to be nonnullable.", action.Column.Name)
			}
		}
		if action.Column.Constraints != nil {
			schema := c.GetSchema(schemaName)
			addingClustered := false
			for _, ci := range action.Column.Constraints.Items {
				cd, ok := ci.(*ast.ConstraintDef)
				if !ok {
					continue
				}
				if cd.Name != "" && schema != nil && c.lookupSchemaObject(schema, cd.Name, false) != kindNone {
					return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", cd.Name)
				}
				if cd.Type == ast.ConstraintPrimaryKey && tableHasPrimaryKey(table) {
					return mssqlErr(ErrMultiplePrimaryKeys, "Table '%s' already has a primary key defined on it.", table.Name)
				}
				if !c.lenientMetadataLoad && (cd.Type == ast.ConstraintPrimaryKey || cd.Type == ast.ConstraintUnique) && cd.Clustered != nil && *cd.Clustered {
					if addingClustered || tableHasClusteredIndex(table) {
						return mssqlErr(ErrMultipleClusteredIndex, "Cannot create more than one clustered index on table '%s'.", table.Name)
					}
					addingClustered = true
				}
				if !c.lenientMetadataLoad && cd.Type == ast.ConstraintForeignKey {
					if err := c.validateForeignKeyTarget(table, []string{action.Column.Name}, cd); err != nil {
						return err
					}
				}
			}
		}
		c.extractColumn(action.Column, table, schemaName, source)
		return nil

	case ast.ATDropColumn:
		col := columnByName(table, action.ColName)
		if col == nil {
			if action.IfExists {
				return nil
			}
			return mssqlErr(ErrAlterColNotFound, "ALTER TABLE DROP COLUMN %s failed because column '%s' does not exist in table '%s'.", action.ColName, action.ColName, table.Name)
		}
		if c.isPeriodColumn(schemaName, table.Name, col.Name) {
			return mssqlErr(ErrDropPeriodColumn, "Cannot drop column '%s' because it is used in the SYSTEM_TIME period definition.", col.Name)
		}
		if dep, ok := columnDependency(c, schemaName, table, col); ok {
			return mssqlErr(ErrColumnDependentObject, "The object '%s' is dependent on column '%s'.", dep, col.Name)
		}
		if len(table.Columns) <= 1 {
			return mssqlErr(ErrDropLastColumn, "ALTER TABLE DROP COLUMN failed because '%s' is the only data column in table '%s'. A table must have at least one data column.", col.Name, table.Name)
		}
		for i, cur := range table.Columns {
			if cur == col {
				table.Columns = append(table.Columns[:i], table.Columns[i+1:]...)
				break
			}
		}
		for i, cur := range table.Columns {
			cur.Position = int32(i + 1)
		}
		return nil

	case ast.ATAlterColumn:
		col := columnByName(table, action.ColName)
		if col == nil {
			return mssqlErr(ErrAlterColNotFound, "ALTER TABLE ALTER COLUMN %s failed because column '%s' does not exist in table '%s'.", action.ColName, action.ColName, table.Name)
		}
		if isComputedColumnType(col.Type) {
			return mssqlErr(ErrAlterComputedColumn, "Cannot alter column '%s' because it is 'computed'.", col.Name)
		}
		if c.isPeriodColumn(schemaName, table.Name, col.Name) {
			return mssqlErr(ErrDropPeriodColumn, "Cannot alter column '%s' because it is used in the SYSTEM_TIME period definition.", col.Name)
		}
		if !c.lenientMetadataLoad && !isValidDataType(action.DataType, false) {
			typeName := ""
			if action.DataType != nil {
				typeName = action.DataType.Name
			}
			return mssqlErr(ErrUnknownColumnType, "Column, parameter, or variable #1: Cannot find data type %s.", typeName)
		}
		newNullable := true
		if action.Column != nil && action.Column.Nullable != nil {
			newNullable = !action.Column.Nullable.NotNull
		} else if col.IsIdentity {
			newNullable = false
		}
		if !c.lenientMetadataLoad && col.IsIdentity {
			if !isValidIdentityDataType(action.DataType) || newNullable {
				return mssqlErr(ErrIdentityColumnType, "Identity column '%s' must be of data type int, bigint, smallint, tinyint, or decimal or numeric with a scale of 0, unencrypted, and constrained to be nonnullable.", col.Name)
			}
		}
		newType := nodeText(action.DataType, source)
		collationChanged := action.Collation != "" && !strings.EqualFold(action.Collation, col.Collation)
		if !c.lenientMetadataLoad && (!typesEqual(col.Type, newType) || collationChanged) {
			if dep, blocked := alterColumnTypeDependency(c, schemaName, table, col, newType, collationChanged); blocked {
				return mssqlErr(ErrColumnDependentObject, "The object '%s' is dependent on column '%s'.", dep, col.Name)
			}
		}
		if !c.lenientMetadataLoad && newNullable != col.Nullable {
			if consumer, bound := schemaBoundColumnConsumer(c, schemaName, table.Name, col.Name); bound {
				return mssqlErr(ErrColumnDependentObject, "The object '%s' is dependent on column '%s'.", consumer, col.Name)
			}
		}
		if newNullable {
			if pkName, inPK := columnInPrimaryKey(table, col.Name); inPK {
				return mssqlErr(ErrColumnDependentObject, "The object '%s' is dependent on column '%s'.", pkName, col.Name)
			}
		}
		col.Type = newType
		col.Nullable = newNullable
		if action.Collation != "" {
			col.Collation = action.Collation
		} else if !isCharacterDataType(action.DataType) {
			col.Collation = ""
		}
		return nil

	case ast.ATAlterColumnAddDrop:
		return c.execAlterColumnAddDrop(table, action)

	case ast.ATAddConstraint:
		if action.Constraint == nil {
			return nil
		}
		return c.execAlterTableAddConstraint(table, schemaName, action.Constraint, source)

	case ast.ATDropConstraint:
		if action.Constraint == nil {
			if hasOnlineOption(action.Options) {
				return mssqlErr(ErrOnlineDropNonClustered, "Only a clustered index can be dropped online.")
			}
			return nil
		}
		name := action.Constraint.Name
		if hasOnlineOption(action.Options) {
			isClustered := false
			for _, idx := range table.Indexes {
				if idx.IsConstraint && strings.EqualFold(idx.Name, name) &&
					(strings.EqualFold(idx.Type, "CLUSTERED") || strings.EqualFold(idx.Type, "CLUSTERED COLUMNSTORE") || (idx.Type == "" && idx.Primary)) {
					isClustered = true
					break
				}
			}
			if !isClustered {
				return mssqlErr(ErrOnlineDropNonClustered, "Only a clustered index can be dropped online.")
			}
		}
		dropped, err := c.dropTableConstraint(schemaName, table, name)
		if err != nil {
			return err
		}
		if dropped {
			return nil
		}
		if action.IfExists {
			return nil
		}
		return mssqlErr(ErrConstraintNotFound, "'%s' is not a constraint.", name)

	case ast.ATAddIndex:
		if action.Index != nil {
			if err := checkDuplicateIndexColumns(action.Index.Columns, action.Index.IncludeCols); err != nil {
				return err
			}
		}
		return mssqlErr(ErrAlterTableIndexDiskTable, "The operation 'ALTER TABLE ADD INDEX' is supported only with memory optimized tables.")

	case ast.ATDropIndex:
		return mssqlErr(ErrAlterTableIndexDiskTable, "The operation 'ALTER TABLE DROP INDEX' is supported only with memory optimized tables.")

	case ast.ATCheckConstraint, ast.ATNocheckConstraint:
		if action.Names == nil {
			return nil
		}
		for _, item := range action.Names.Items {
			s, ok := item.(*ast.String)
			if !ok {
				continue
			}
			if strings.EqualFold(s.Str, "ALL") {
				return nil
			}
			if hasCheckOrFKConstraint(table, s.Str) {
				continue
			}
			if schema := c.GetSchema(schemaName); schema != nil && c.lookupSchemaObject(schema, s.Str, false) != kindNone {
				return mssqlErr(ErrCheckConstraintWrongKind, "Object '%s' cannot be disabled or enabled. This action applies only to foreign key and check constraints.", s.Str)
			}
			return mssqlErr(ErrCheckConstraintNotFound, "Constraint '%s' does not exist.", s.Str)
		}
		return nil

	case ast.ATRebuild:
		if action.Partition != nil {
			if _, ok := action.Partition.(*ast.ColumnRef); ok {
				return mssqlErr(ErrInvalidColumnName, "Invalid column name '%s'.", nodeText(action.Partition, source))
			}
			return mssqlErr(ErrRebuildPartitionInvalid, "Cannot specify partition number in the alter index statement as the table '%s' is not partitioned.", table.Name)
		}
		return nil

	case ast.ATSwitchPartition:
		if action.TargetName != nil {
			targetSchema, targetTable := tableRefSchemaObject(action.TargetName, c.defaultSchema)
			if strings.EqualFold(targetSchema, schemaName) && strings.EqualFold(targetTable, table.Name) {
				return mssqlErr(ErrSwitchSameTable, "ALTER TABLE SWITCH statement failed. The source table '%s' and target table '%s' are same.", table.Name, targetTable)
			}
		}
		return nil

	case ast.ATSet:
		return c.execAlterTableSet(schemaName, table, action.Options)

	case ast.ATAddPeriod:
		if action.Names == nil || len(action.Names.Items) != 2 {
			return mssqlErr(ErrSystemTimeInvalidCols, "Temporal generated always column has invalid data type.")
		}
		var periodCols [2]*metadata.ColumnMetadata
		for i, item := range action.Names.Items {
			s, ok := item.(*ast.String)
			if !ok {
				return mssqlErr(ErrSystemTimeInvalidCols, "Temporal generated always column has invalid data type.")
			}
			col := columnByName(table, s.Str)
			if col == nil {
				return mssqlErr(ErrAlterColNotFound, "Column '%s' does not exist in table '%s'.", s.Str, table.Name)
			}
			periodCols[i] = col
		}
		if strings.EqualFold(periodCols[0].Name, periodCols[1].Name) {
			return mssqlErr(ErrSystemTimeInvalidCols, "Temporal generated always column '%s' has invalid data type.", periodCols[0].Name)
		}
		p0, ok0 := datetime2Precision(periodCols[0].Type)
		p1, ok1 := datetime2Precision(periodCols[1].Type)
		if !ok0 || !ok1 || periodCols[0].Nullable || periodCols[1].Nullable || p0 != p1 {
			return mssqlErr(ErrSystemTimeInvalidCols, "Temporal generated always column '%s' has invalid data type.", periodCols[0].Name)
		}
		periodKey := tableKey{schema: normalizeOrDefault(schemaName, c.defaultSchema), table: table.Name}
		if _, exists := c.systemTimeTables[periodKey]; exists {
			return mssqlErr(ErrPeriodAlreadyDefined, "Cannot define SYSTEM_TIME period more than once on table '%s'.", table.Name)
		}
		if c.systemTimeTables == nil {
			c.systemTimeTables = make(map[tableKey][2]string)
		}
		c.systemTimeTables[periodKey] = [2]string{periodCols[0].Name, periodCols[1].Name}
		return nil

	case ast.ATDropPeriod:
		periodKey := tableKey{schema: normalizeOrDefault(schemaName, c.defaultSchema), table: table.Name}
		if _, exists := c.systemTimeTables[periodKey]; !exists {
			return mssqlErr(ErrDropPeriodNotDefined, "SYSTEM_TIME period cannot be dropped from table '%s', because SYSTEM_TIME period does not exist.", table.Name)
		}
		if c.systemVersionedTables[periodKey] {
			return mssqlErr(ErrDropPeriodSystemVersioned, "Cannot drop SYSTEM_TIME period on table '%s' while SYSTEM_VERSIONING is ON.", table.Name)
		}
		delete(c.systemTimeTables, periodKey)
		return nil
	}
	return nil
}

func (c *Catalog) execAlterColumnAddDrop(table *metadata.TableMetadata, action *ast.AlterTableAction) error {
	var opt *ast.AlterColumnOption
	if action.Options != nil {
		for _, item := range action.Options.Items {
			if o, ok := item.(*ast.AlterColumnOption); ok {
				opt = o
				break
			}
		}
	}
	if opt != nil && strings.EqualFold(opt.Action, "DROP") && strings.EqualFold(opt.Option, "ROWGUIDCOL") {
		return mssqlErr(ErrDropRowguidcolMissing, "ALTER TABLE ALTER COLUMN DROP ROWGUIDCOL failed because a column does not exist in table '%s' with ROWGUIDCOL property.", table.Name)
	}
	col := columnByName(table, action.ColName)
	if col == nil {
		return mssqlErr(ErrAlterColNotFound, "ALTER TABLE ALTER COLUMN %s failed because column '%s' does not exist in table '%s'.", action.ColName, action.ColName, table.Name)
	}
	if opt == nil {
		return nil
	}
	if strings.EqualFold(opt.Action, "ADD") {
		switch strings.ToUpper(opt.Option) {
		case "ROWGUIDCOL":
			if normalizeType(col.Type) != "uniqueidentifier" {
				return mssqlErr(ErrRowguidcolType, "The ROWGUIDCOL property can only be specified on the uniqueidentifier data type.")
			}
		case "NOT FOR REPLICATION":
			if !col.IsIdentity {
				return mssqlErr(ErrAlterColNFRNotIdentity, "Cannot alter NOT FOR REPLICATION attribute on column '%s' in table '%s' because this column is not an identity column.", col.Name, table.Name)
			}
		case "PERSISTED":
			if !isComputedColumnType(col.Type) {
				return mssqlErr(ErrAlterColPersistedNonComp, "PERSISTED attribute cannot be altered on column '%s' because this column is not computed.", col.Name)
			}
			trimmed := strings.TrimSpace(col.Type)
			if !strings.HasSuffix(strings.ToUpper(trimmed), " PERSISTED") {
				col.Type = trimmed + " PERSISTED"
			}
		case "MASKED":
			if !isValidMaskingFunction(opt.MaskFunction) {
				return mssqlErr(ErrInvalidMaskingFunction, "Invalid data masking format for function '%s' in column '%s'.", opt.MaskFunction, col.Name)
			}
		case "HIDDEN":
			return mssqlErr(ErrAlterColHiddenNotPeriod, "Cannot alter HIDDEN attribute on column '%s' in table '%s' because this column is not a generated always column.", col.Name, table.Name)
		case "SPARSE":
			return nil
		}
	} else if strings.EqualFold(opt.Action, "DROP") {
		switch strings.ToUpper(opt.Option) {
		case "PERSISTED":
			if !isComputedColumnType(col.Type) {
				return mssqlErr(ErrAlterColPersistedNonComp, "PERSISTED attribute cannot be altered on column '%s' because this column is not computed.", col.Name)
			}
			trimmed := strings.TrimSpace(col.Type)
			if len(trimmed) > len(" PERSISTED") && strings.EqualFold(trimmed[len(trimmed)-len(" PERSISTED"):], " PERSISTED") {
				col.Type = strings.TrimSpace(trimmed[:len(trimmed)-len(" PERSISTED")])
			}
		}
	}
	return nil
}

func isValidMaskingFunction(fn string) bool {
	low := strings.ToLower(strings.TrimSpace(fn))
	return low == "default()" || low == "email()" || strings.HasPrefix(low, "partial(") || strings.HasPrefix(low, "random(") || strings.HasPrefix(low, "datetime(")
}

func (c *Catalog) execAlterTableAddConstraint(table *metadata.TableMetadata, schemaName string, n *ast.ConstraintDef, source string) error {
	if n.Type == ast.ConstraintEdge {
		return mssqlErr(ErrEdgeConstraintNonEdge, "Edge constraint cannot be created on table '%s'. The table is not an edge table.", table.Name)
	}
	switch n.Type {
	case ast.ConstraintPrimaryKey, ast.ConstraintUnique:
		if err := checkDuplicateIndexColumns(n.Columns, nil); err != nil {
			return err
		}
		if n.Type == ast.ConstraintPrimaryKey && tableHasPrimaryKey(table) {
			return mssqlErr(ErrMultiplePrimaryKeys, "Table '%s' already has a primary key defined on it.", table.Name)
		}
		if !c.lenientMetadataLoad && n.Clustered != nil && *n.Clustered && tableHasClusteredIndex(table) {
			return mssqlErr(ErrMultipleClusteredIndex, "Cannot create more than one clustered index on table '%s'.", table.Name)
		}
		for _, colName := range stringNodeList(n.Columns) {
			if columnByName(table, colName) == nil {
				return mssqlErr(ErrIndexColumnNotFound, "Column name '%s' does not exist in the target table or view.", colName)
			}
		}
		if n.Name != "" {
			if schema := c.GetSchema(schemaName); schema != nil && c.lookupSchemaObject(schema, n.Name, false) != kindNone {
				return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", n.Name)
			}
		}
		c.extractTableConstraint(n, table, schemaName, source)
		return nil

	case ast.ConstraintDefault:
		cols := stringNodeList(n.Columns)
		if len(cols) == 0 {
			return nil
		}
		col := columnByName(table, cols[0])
		if col == nil {
			return mssqlErr(ErrAddDefaultColumnNotFound, "Cannot create default on column '%s' because it does not exist in table '%s'.", cols[0], table.Name)
		}
		if col.Default != "" {
			return mssqlErr(ErrColumnAlreadyHasDefault, "Column already has a DEFAULT bound to it.")
		}
		if n.Name != "" {
			if schema := c.GetSchema(schemaName); schema != nil && c.lookupSchemaObject(schema, n.Name, false) != kindNone {
				return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", n.Name)
			}
		}
		col.Default = nodeText(n.Expr, source)
		if n.Name != "" {
			col.DefaultConstraintName = n.Name
		} else if !c.omitUnnamedDefaultConstraintNames {
			col.DefaultConstraintName = "DF_" + table.Name + "_" + col.Name
		}
		return nil

	case ast.ConstraintForeignKey:
		cols := stringNodeList(n.Columns)
		if len(cols) == 0 {
			refName := ""
			if n.RefTable != nil {
				refName = n.RefTable.Object
			}
			return mssqlErr(ErrFKColumnCountMismatch, "The number of columns in the referencing column list does not match those of the primary key in the referenced table '%s'.", refName)
		}
		for _, colName := range cols {
			if columnByName(table, colName) == nil {
				return mssqlErr(ErrFKColumnNotFound, "Foreign key references invalid column '%s' in referencing table '%s'.", colName, table.Name)
			}
		}
		if n.Name != "" {
			if schema := c.GetSchema(schemaName); schema != nil && c.lookupSchemaObject(schema, n.Name, false) != kindNone {
				return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", n.Name)
			}
		}
		if !c.lenientMetadataLoad {
			if err := c.validateForeignKeyTarget(table, cols, n); err != nil {
				return err
			}
		}
		c.extractTableConstraint(n, table, schemaName, source)
		return nil

	case ast.ConstraintCheck:
		if n.Name != "" {
			if schema := c.GetSchema(schemaName); schema != nil && c.lookupSchemaObject(schema, n.Name, false) != kindNone {
				return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", n.Name)
			}
		}
		c.extractTableConstraint(n, table, schemaName, source)
		return nil
	}
	return nil
}

func (c *Catalog) validateForeignKeyTarget(srcTable *metadata.TableMetadata, srcCols []string, n *ast.ConstraintDef) error {
	if n == nil || n.RefTable == nil {
		return nil
	}
	if n.RefTable.Database != "" && (c.database == "" || !strings.EqualFold(n.RefTable.Database, c.database)) {
		return mssqlErr(ErrFKCrossDatabase, "Cross-database foreign key references are not supported. Foreign key '%s'.", n.Name)
	}
	refSchema, refTableName := tableRefSchemaObject(n.RefTable, c.defaultSchema)
	refTable := c.GetTable(refSchema, refTableName)
	if refTable == nil {
		return mssqlErr(ErrFKInvalidRefTable, "Foreign key '%s' references invalid table '%s'.", n.Name, refTableName)
	}
	refCols := stringNodeList(n.RefColumns)
	if len(refCols) == 0 {
		var pk *metadata.IndexMetadata
		for _, idx := range refTable.Indexes {
			if idx.Primary {
				pk = idx
				break
			}
		}
		if pk == nil {
			return mssqlErr(ErrFKNoCandidateKey, "There are no primary or candidate keys in the referenced table '%s' that match the referencing column list in the foreign key '%s'.", refTableName, n.Name)
		}
		if len(srcCols) != len(pk.Expressions) {
			return mssqlErr(ErrFKColumnCountMismatch, "The number of columns in the referencing column list does not match those of the primary key in the referenced table '%s'.", refTableName)
		}
		return nil
	}
	if len(srcCols) != len(refCols) {
		return mssqlErr(ErrFKColumnCountMismatch, "The number of columns in the referencing column list does not match those of the primary key in the referenced table '%s'.", refTableName)
	}
	for _, rc := range refCols {
		if columnByName(refTable, rc) == nil {
			return mssqlErr(ErrFKColumnNotFound, "Foreign key '%s' references invalid column '%s' in referenced table '%s'.", n.Name, rc, refTableName)
		}
	}
	hasCandidateKey := false
	for _, idx := range refTable.Indexes {
		if !idx.Primary && !idx.Unique {
			continue
		}
		if indexWhereClause(idx.Definition) != "" {
			continue
		}
		if len(idx.Expressions) != len(refCols) {
			continue
		}
		allMatch := true
		for _, rc := range refCols {
			found := false
			for _, ic := range idx.Expressions {
				if strings.EqualFold(rc, ic) {
					found = true
					break
				}
			}
			if !found {
				allMatch = false
				break
			}
		}
		if allMatch {
			hasCandidateKey = true
			break
		}
	}
	if !hasCandidateKey {
		return mssqlErr(ErrFKNoCandidateKey, "There are no primary or candidate keys in the referenced table '%s' that match the referencing column list in the foreign key '%s'.", refTableName, n.Name)
	}
	return nil
}

func (c *Catalog) dropTableConstraint(schemaName string, table *metadata.TableMetadata, name string) (bool, error) {
	dropIndexConstraint := func(i int) (bool, error) {
		idx := table.Indexes[i]
		if refTable, fkName, referenced := indexReferencedByFK(c, schemaName, table, idx); referenced {
			return false, mssqlErr(ErrDropConstraintReferencedByFK, "The constraint '%s' is being referenced by table '%s', foreign key constraint '%s'.", idx.Name, refTable, fkName)
		}
		if idx.Primary && isClusteredPrimaryKey(idx) && tableHasXMLIndex(table) {
			return false, mssqlErr(ErrDropPKWithXMLIndex, "Cannot drop the clustered primary key '%s' on table '%s' because the table has an XML index.", idx.Name, table.Name)
		}
		if idx.Primary && c.systemVersionedTables[c.temporalKey(schemaName, table.Name)] {
			return false, mssqlErr(ErrDropPKSystemVersioned, "Cannot drop primary key constraint '%s' on system-versioned temporal table '%s.%s'.", idx.Name, schemaName, table.Name)
		}
		table.Indexes = append(table.Indexes[:i], table.Indexes[i+1:]...)
		return true, nil
	}

	// Pass 1: exact name match across all constraint kinds.
	for i, idx := range table.Indexes {
		if idx.IsConstraint && idx.Name == name {
			return dropIndexConstraint(i)
		}
	}
	for i, fk := range table.ForeignKeys {
		if fk.Name == name {
			table.ForeignKeys = append(table.ForeignKeys[:i], table.ForeignKeys[i+1:]...)
			return true, nil
		}
	}
	for i, ck := range table.CheckConstraints {
		if ck.Name == name {
			table.CheckConstraints = append(table.CheckConstraints[:i], table.CheckConstraints[i+1:]...)
			return true, nil
		}
	}
	for _, col := range table.Columns {
		if col.DefaultConstraintName != "" && col.DefaultConstraintName == name {
			col.Default = ""
			col.DefaultConstraintName = ""
			return true, nil
		}
	}

	// Pass 2: unambiguous case-folded fallback across all constraint kinds.
	type constraintMatch struct {
		kind int // 1=idx, 2=fk, 3=ck, 4=df
		pos  int
	}
	var match constraintMatch
	matches := 0
	for i, idx := range table.Indexes {
		if idx.IsConstraint && strings.EqualFold(idx.Name, name) {
			match = constraintMatch{kind: 1, pos: i}
			matches++
		}
	}
	for i, fk := range table.ForeignKeys {
		if strings.EqualFold(fk.Name, name) {
			match = constraintMatch{kind: 2, pos: i}
			matches++
		}
	}
	for i, ck := range table.CheckConstraints {
		if strings.EqualFold(ck.Name, name) {
			match = constraintMatch{kind: 3, pos: i}
			matches++
		}
	}
	for i, col := range table.Columns {
		if col.DefaultConstraintName != "" && strings.EqualFold(col.DefaultConstraintName, name) {
			match = constraintMatch{kind: 4, pos: i}
			matches++
		}
	}
	if matches != 1 {
		return false, nil
	}
	switch match.kind {
	case 1:
		return dropIndexConstraint(match.pos)
	case 2:
		table.ForeignKeys = append(table.ForeignKeys[:match.pos], table.ForeignKeys[match.pos+1:]...)
		return true, nil
	case 3:
		table.CheckConstraints = append(table.CheckConstraints[:match.pos], table.CheckConstraints[match.pos+1:]...)
		return true, nil
	case 4:
		table.Columns[match.pos].Default = ""
		table.Columns[match.pos].DefaultConstraintName = ""
		return true, nil
	}
	return false, nil
}

func isCharacterDataType(dt *ast.DataType) bool {
	if dt == nil || dt.Schema != "" {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(dt.Name)) {
	case "char", "character", "varchar", "character varying", "char varying",
		"nchar", "national char", "national character",
		"nvarchar", "national char varying", "national character varying",
		"text", "ntext", "sysname":
		return true
	}
	return false
}

func isValidIdentityDataType(dt *ast.DataType) bool {
	if dt == nil || dt.Schema != "" {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(dt.Name)) {
	case "tinyint", "smallint", "int", "integer", "bigint":
		return true
	case "decimal", "dec", "numeric":
		if dt.Scale == nil {
			return true
		}
		return isZeroNumericExpr(dt.Scale)
	}
	return false
}

func (c *Catalog) temporalKey(schemaName, tableName string) tableKey {
	if s := c.GetSchema(schemaName); s != nil {
		schemaName = s.Name
	}
	return tableKey{schema: normalizeOrDefault(schemaName, c.defaultSchema), table: tableName}
}

func (c *Catalog) isPeriodColumn(schemaName, tableName, colName string) bool {
	periodKey := c.temporalKey(schemaName, tableName)
	cols, ok := c.systemTimeTables[periodKey]
	if !ok {
		return false
	}
	return strings.EqualFold(cols[0], colName) || strings.EqualFold(cols[1], colName)
}

func isClusteredPrimaryKey(idx *metadata.IndexMetadata) bool {
	if idx == nil || !idx.Primary {
		return false
	}
	return idx.Type == "" || strings.EqualFold(idx.Type, "CLUSTERED")
}

func tableHasXMLIndex(table *metadata.TableMetadata) bool {
	if table == nil {
		return false
	}
	for _, idx := range table.Indexes {
		if strings.EqualFold(idx.Type, "XML") {
			return true
		}
	}
	return false
}

func createTableHasPrimaryKey(n *ast.CreateTableStmt) bool {
	if n == nil {
		return false
	}
	if n.Columns != nil {
		for _, item := range n.Columns.Items {
			col, ok := item.(*ast.ColumnDef)
			if !ok || col.Constraints == nil {
				continue
			}
			for _, cItem := range col.Constraints.Items {
				if cd, ok := cItem.(*ast.ConstraintDef); ok && cd.Type == ast.ConstraintPrimaryKey {
					return true
				}
			}
		}
	}
	if n.Constraints != nil {
		for _, item := range n.Constraints.Items {
			if cd, ok := item.(*ast.ConstraintDef); ok && cd.Type == ast.ConstraintPrimaryKey {
				return true
			}
		}
	}
	return false
}

func tableHasPrimaryKey(table *metadata.TableMetadata) bool {
	if table == nil {
		return false
	}
	for _, idx := range table.Indexes {
		if idx.Primary {
			return true
		}
	}
	return false
}

// qualifyModuleDefinition ensures a view, procedure, or function CREATE header
// uses schemaName. When forceQualify is false, unqualified headers and headers
// already matching schemaName are returned unchanged so existing dumps keep
// exact formatting parity; when forceQualify is true (used after ALTER SCHEMA
// TRANSFER), even unqualified headers are qualified with schemaName.
func qualifyModuleDefinition(definition, schemaName string, forceQualify bool) string {
	if schemaName == "" || strings.TrimSpace(definition) == "" {
		return definition
	}
	stmts, err := mssql.Parse(definition)
	if err != nil || len(stmts) == 0 {
		return definition
	}
	var nameRef *ast.TableRef
	switch n := stmts[0].AST.(type) {
	case *ast.CreateViewStmt:
		nameRef = n.Name
	case *ast.CreateProcedureStmt:
		nameRef = n.Name
	case *ast.CreateFunctionStmt:
		nameRef = n.Name
	}
	if nameRef == nil || nameRef.Loc.Start < 0 || nameRef.Loc.End <= nameRef.Loc.Start || nameRef.Loc.End > len(definition) {
		return definition
	}
	if nameRef.Schema == "" && !forceQualify {
		return definition
	}
	if nameRef.Schema != "" && strings.EqualFold(nameRef.Schema, schemaName) {
		return definition
	}
	qualified := quote(schemaName) + "." + quote(nameRef.Object)
	return definition[:nameRef.Loc.Start] + qualified + definition[nameRef.Loc.End:]
}

// qualifyTriggerDefinition updates a trigger CREATE statement's trigger schema
// and ON target table/view schema. When forceQualify is false, already-matching
// or unqualified references are preserved verbatim.
func qualifyTriggerDefinition(definition, schemaName, parentName string, forceQualify bool) string {
	if schemaName == "" || strings.TrimSpace(definition) == "" {
		return definition
	}
	stmts, err := mssql.Parse(definition)
	if err != nil || len(stmts) == 0 {
		return definition
	}
	trig, ok := stmts[0].AST.(*ast.CreateTriggerStmt)
	if !ok {
		return definition
	}
	type replacement struct {
		start int
		end   int
		text  string
	}
	var reps []replacement
	addRef := func(ref *ast.TableRef, targetObj string) {
		if ref == nil || ref.Loc.Start < 0 || ref.Loc.End <= ref.Loc.Start || ref.Loc.End > len(definition) {
			return
		}
		if ref.Schema == "" && !forceQualify {
			return
		}
		if ref.Schema != "" && strings.EqualFold(ref.Schema, schemaName) && (targetObj == "" || strings.EqualFold(ref.Object, targetObj)) {
			return
		}
		obj := ref.Object
		if targetObj != "" {
			obj = targetObj
		}
		reps = append(reps, replacement{
			start: ref.Loc.Start,
			end:   ref.Loc.End,
			text:  quote(schemaName) + "." + quote(obj),
		})
	}
	addRef(trig.Name, "")
	addRef(trig.Table, parentName)
	if len(reps) == 0 {
		return definition
	}
	// Apply replacements from right to left so earlier byte offsets remain valid.
	if len(reps) == 2 && reps[0].start < reps[1].start {
		reps[0], reps[1] = reps[1], reps[0]
	}
	out := definition
	for _, r := range reps {
		out = out[:r.start] + r.text + out[r.end:]
	}
	return out
}

func isZeroNumericExpr(expr ast.ExprNode) bool {
	if expr == nil {
		return false
	}
	switch e := expr.(type) {
	case *ast.ParenExpr:
		return isZeroNumericExpr(e.Expr)
	case *ast.UnaryExpr:
		return isZeroNumericExpr(e.Operand)
	case *ast.Literal:
		s := strings.TrimSpace(e.Str)
		if s == "" {
			return false
		}
		hasDigit := false
		for i := 0; i < len(s); i++ {
			c := s[i]
			if c == '0' {
				hasDigit = true
			} else if c != '.' {
				return false
			}
		}
		return hasDigit
	}
	return false
}

func columnInPrimaryKey(table *metadata.TableMetadata, colName string) (string, bool) {
	for _, idx := range table.Indexes {
		if !idx.Primary {
			continue
		}
		for _, expr := range idx.Expressions {
			if strings.EqualFold(expr, colName) {
				return idx.Name, true
			}
		}
	}
	return "", false
}

func alterColumnTypeDependency(c *Catalog, schemaName string, table *metadata.TableMetadata, col *metadata.ColumnMetadata, newType string, collationChanged bool) (string, bool) {
	if col.Default != "" {
		name := col.DefaultConstraintName
		if name == "" {
			name = "DEFAULT"
		}
		return name, true
	}
	widenOK := !collationChanged && isCompatibleVarLengthWidening(col.Type, newType)
	for _, idx := range table.Indexes {
		for _, expr := range idx.Expressions {
			if strings.EqualFold(expr, col.Name) {
				if !widenOK {
					return idx.Name, true
				}
			}
		}
		for _, inc := range indexIncludeColumns(idx.Definition) {
			if strings.EqualFold(inc, col.Name) {
				if !widenOK {
					return idx.Name, true
				}
			}
		}
		if where := indexWhereClause(idx.Definition); where != "" && checkExprReferencesColumn(where, table.Name, col.Name) {
			return idx.Name, true
		}
	}
	for _, fk := range table.ForeignKeys {
		for _, fc := range fk.Columns {
			if strings.EqualFold(fc, col.Name) {
				return fk.Name, true
			}
		}
	}
	for _, s := range c.schemas {
		for _, t := range s.Tables {
			for _, fk := range t.ForeignKeys {
				refSchema := fk.ReferencedSchema
				if refSchema == "" {
					refSchema = c.defaultSchema
				}
				if c.GetTable(refSchema, fk.ReferencedTable) == table {
					for _, rc := range fk.ReferencedColumns {
						if strings.EqualFold(rc, col.Name) {
							return fk.Name, true
						}
					}
				}
			}
		}
	}
	for _, ck := range table.CheckConstraints {
		if checkExprReferencesColumn(ck.Expression, table.Name, col.Name) {
			name := ck.Name
			if name == "" {
				name = "CHECK"
			}
			return name, true
		}
	}
	for _, other := range table.Columns {
		if other == col || !isComputedColumnType(other.Type) {
			continue
		}
		if computedExprReferencesColumn(other.Type, table.Name, col.Name) {
			return other.Name, true
		}
	}
	if consumer, bound := schemaBoundColumnConsumer(c, schemaName, table.Name, col.Name); bound {
		return consumer, true
	}
	return "", false
}

// isCompatibleVarLengthWidening reports whether altering oldType to newType is
// a finite length increase on the same variable-length string/binary base type,
// which SQL Server permits in place on indexed columns.
func isCompatibleVarLengthWidening(oldType, newType string) bool {
	oldBase, oldLen, ok1 := parseVarLengthType(oldType)
	newBase, newLen, ok2 := parseVarLengthType(newType)
	if !ok1 || !ok2 || oldBase != newBase {
		return false
	}
	return newLen >= oldLen
}

func parseVarLengthType(t string) (string, int, bool) {
	s := strings.ToLower(strings.TrimSpace(t))
	open := strings.IndexByte(s, '(')
	close := strings.LastIndexByte(s, ')')
	if open <= 0 || close != len(s)-1 {
		return "", 0, false
	}
	base := strings.TrimSpace(s[:open])
	switch base {
	case "varchar", "nvarchar", "varbinary":
	default:
		return "", 0, false
	}
	arg := strings.TrimSpace(s[open+1 : close])
	if arg == "" || arg == "max" {
		return "", 0, false
	}
	n := 0
	for i := 0; i < len(arg); i++ {
		if arg[i] < '0' || arg[i] > '9' {
			return "", 0, false
		}
		n = n*10 + int(arg[i]-'0')
	}
	if n <= 0 {
		return "", 0, false
	}
	return base, n, true
}

func hasCheckOrFKConstraint(table *metadata.TableMetadata, name string) bool {
	for _, fk := range table.ForeignKeys {
		if strings.EqualFold(fk.Name, name) {
			return true
		}
	}
	for _, ck := range table.CheckConstraints {
		if strings.EqualFold(ck.Name, name) {
			return true
		}
	}
	return false
}

func columnDependency(c *Catalog, schemaName string, table *metadata.TableMetadata, col *metadata.ColumnMetadata) (string, bool) {
	if col.Default != "" {
		name := col.DefaultConstraintName
		if name == "" {
			name = "DEFAULT"
		}
		return name, true
	}
	for _, idx := range table.Indexes {
		for _, expr := range idx.Expressions {
			if strings.EqualFold(expr, col.Name) {
				return idx.Name, true
			}
		}
		for _, inc := range indexIncludeColumns(idx.Definition) {
			if strings.EqualFold(inc, col.Name) {
				return idx.Name, true
			}
		}
		if where := indexWhereClause(idx.Definition); where != "" && checkExprReferencesColumn(where, table.Name, col.Name) {
			return idx.Name, true
		}
	}
	for _, fk := range table.ForeignKeys {
		for _, fc := range fk.Columns {
			if strings.EqualFold(fc, col.Name) {
				return fk.Name, true
			}
		}
	}
	for _, ck := range table.CheckConstraints {
		if checkExprReferencesColumn(ck.Expression, table.Name, col.Name) {
			name := ck.Name
			if name == "" {
				name = "CHECK"
			}
			return name, true
		}
	}
	for _, other := range table.Columns {
		if other == col || !isComputedColumnType(other.Type) {
			continue
		}
		if computedExprReferencesColumn(other.Type, table.Name, col.Name) {
			return other.Name, true
		}
	}
	for _, s := range c.schemas {
		for _, t := range s.Tables {
			for _, fk := range t.ForeignKeys {
				refSchema := fk.ReferencedSchema
				if refSchema == "" {
					refSchema = c.defaultSchema
				}
				if c.GetTable(refSchema, fk.ReferencedTable) == table {
					for _, rc := range fk.ReferencedColumns {
						if strings.EqualFold(rc, col.Name) {
							return fk.Name, true
						}
					}
				}
			}
		}
	}
	if consumer, bound := schemaBoundColumnConsumer(c, schemaName, table.Name, col.Name); bound {
		return consumer, true
	}
	return "", false
}

func exprReferencesColumn(node ast.Node, colName string) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if found {
			return false
		}
		if ref, ok := n.(*ast.ColumnRef); ok && strings.EqualFold(ref.Column, colName) {
			found = true
			return false
		}
		return true
	})
	return found
}

func exprReferencesTableColumn(node ast.Node, tableName, colName string) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if found {
			return false
		}
		ref, ok := n.(*ast.ColumnRef)
		if !ok {
			return true
		}
		if strings.EqualFold(ref.Column, colName) && (ref.Table == "" || strings.EqualFold(ref.Table, tableName)) {
			found = true
			return false
		}
		return true
	})
	return found
}

func checkExprReferencesColumn(expr, tableName, colName string) bool {
	if strings.TrimSpace(expr) == "" {
		return false
	}
	stmts, err := mssql.Parse("SELECT 1 WHERE (" + expr + ")")
	if err != nil || len(stmts) == 0 {
		return false
	}
	sel, ok := stmts[0].AST.(*ast.SelectStmt)
	if !ok || sel.WhereClause == nil {
		return false
	}
	return exprReferencesTableColumn(sel.WhereClause, tableName, colName)
}

func computedExprReferencesColumn(colType, tableName, colName string) bool {
	if !isComputedColumnType(colType) {
		return false
	}
	stmts, err := mssql.Parse("CREATE TABLE _ (_ " + colType + ")")
	if err != nil || len(stmts) == 0 {
		return false
	}
	ct, ok := stmts[0].AST.(*ast.CreateTableStmt)
	if !ok || ct.Columns == nil || len(ct.Columns.Items) == 0 {
		return false
	}
	cd, ok := ct.Columns.Items[0].(*ast.ColumnDef)
	if !ok || cd.Computed == nil || cd.Computed.Expr == nil {
		return false
	}
	return exprReferencesTableColumn(cd.Computed.Expr, tableName, colName)
}

// hasOptionOn reports whether opts contains optionName set to ON.
// Parser option lists use either a single "NAME=VALUE" entry (CREATE/DROP INDEX)
// or paired "NAME", "VALUE" entries (ALTER TABLE DROP CONSTRAINT).
func hasOptionOn(opts *ast.List, optionName string) bool {
	if opts == nil {
		return false
	}
	for i := 0; i < len(opts.Items); i++ {
		s, ok := opts.Items[i].(*ast.String)
		if !ok {
			continue
		}
		raw := strings.TrimSpace(s.Str)
		if k, v, hasEq := strings.Cut(raw, "="); hasEq {
			if strings.EqualFold(strings.TrimSpace(k), optionName) {
				return strings.EqualFold(strings.TrimSpace(v), "ON")
			}
			continue
		}
		if strings.EqualFold(raw, optionName) && i+1 < len(opts.Items) {
			if next, ok := opts.Items[i+1].(*ast.String); ok {
				return strings.EqualFold(strings.TrimSpace(next.Str), "ON")
			}
		}
	}
	return false
}

func hasOnlineOption(opts *ast.List) bool {
	return hasOptionOn(opts, "ONLINE")
}

func (c *Catalog) execAlterTableSet(schemaName string, table *metadata.TableMetadata, opts *ast.List) error {
	if opts == nil {
		return nil
	}
	for i := 0; i < len(opts.Items); i++ {
		s, ok := opts.Items[i].(*ast.String)
		if !ok {
			continue
		}
		name := strings.ToUpper(strings.TrimSpace(s.Str))
		val := ""
		if i+1 < len(opts.Items) {
			if vs, ok := opts.Items[i+1].(*ast.String); ok {
				val = strings.ToUpper(strings.TrimSpace(vs.Str))
			}
		}
		switch name {
		case "SYSTEM_VERSIONING":
			periodKey := c.temporalKey(schemaName, table.Name)
			if val == "ON" {
				if _, exists := c.systemTimeTables[periodKey]; !exists {
					return mssqlErr(ErrSystemVersioningNoPeriod, "Cannot set SYSTEM_VERSIONING to ON when SYSTEM_TIME period is not defined.")
				}
				if !tableHasPrimaryKey(table) {
					return mssqlErr(ErrDropPKSystemVersioned, "System versioned temporal table '%s.%s' must have primary key defined.", schemaName, table.Name)
				}
				if c.systemVersionedTables == nil {
					c.systemVersionedTables = make(map[tableKey]bool)
				}
				c.systemVersionedTables[periodKey] = true
				return nil
			}
			if val == "OFF" {
				if !c.systemVersionedTables[periodKey] {
					return mssqlErr(ErrSystemVersioningNotEnabled, "SYSTEM_VERSIONING is not turned ON for table '%s'.", table.Name)
				}
				delete(c.systemVersionedTables, periodKey)
				return nil
			}
		case "FILESTREAM_ON":
			if !c.lenientMetadataLoad && !isValidFilegroup(val) {
				return mssqlErr(ErrInvalidFilegroup, "Invalid filegroup '%s' specified.", val)
			}
		case "FILETABLE_DIRECTORY":
			return mssqlErr(ErrFileTableDirectoryNonFT, "The option 'FILETABLE_DIRECTORY' is only valid when used on a FileTable.")
		}
	}
	return nil
}

func checkExecuteAsOptions(opts *ast.List) error {
	if opts == nil {
		return nil
	}
	for _, item := range opts.Items {
		s, ok := item.(*ast.String)
		if !ok {
			continue
		}
		upper := strings.ToUpper(strings.TrimSpace(s.Str))
		if strings.HasPrefix(upper, "EXECUTE AS ") {
			target := strings.TrimSpace(s.Str[len("EXECUTE AS "):])
			unquoted := strings.Trim(target, "'\"[]")
			switch strings.ToUpper(unquoted) {
			case "CALLER", "SELF", "OWNER", "DBO":
			default:
				return mssqlErr(ErrSecurableNotFound, "Cannot execute as the database principal because the principal \"%s\" does not exist.", unquoted)
			}
		}
	}
	return nil
}

func isNonNullLiteralExpr(expr ast.ExprNode) bool {
	if expr == nil {
		return false
	}
	switch e := expr.(type) {
	case *ast.ParenExpr:
		return isNonNullLiteralExpr(e.Expr)
	case *ast.Literal:
		return e.Type != ast.LitNull
	}
	return false
}

// --- DROP ---

func (c *Catalog) execDrop(n *ast.DropStmt) error {
	switch n.ObjectType {
	case ast.DropSchema:
		return c.execDropSchema(n)
	case ast.DropIndex:
		return c.execDropIndex(n)
	case ast.DropTable:
		return c.execDropObjects(n, kindTable, "TABLE")
	case ast.DropView:
		return c.execDropObjects(n, kindView, "VIEW")
	case ast.DropProcedure:
		return c.execDropObjects(n, kindProcedure, "PROCEDURE")
	case ast.DropFunction:
		return c.execDropObjects(n, kindFunction, "FUNCTION")
	case ast.DropSequence:
		return c.execDropObjects(n, kindSequence, "SEQUENCE")
	case ast.DropSynonym:
		return c.execDropObjects(n, kindSynonym, "SYNONYM")
	case ast.DropTrigger:
		if n.OnDatabase || n.OnAllServer {
			if n.IfExists {
				return nil
			}
			name := ""
			if n.Names != nil && len(n.Names.Items) > 0 {
				if ref, ok := n.Names.Items[0].(*ast.TableRef); ok {
					name = ref.Object
				}
			}
			return mssqlErr(ErrCannotDropNotFound, "Cannot drop the trigger '%s', because it does not exist or you do not have permission.", name)
		}
		return c.execDropObjects(n, kindTrigger, "TRIGGER")
	}
	return nil
}

func (c *Catalog) execDropSchema(n *ast.DropStmt) error {
	if n.Names == nil {
		return nil
	}
	for _, item := range n.Names.Items {
		ref, ok := item.(*ast.TableRef)
		if !ok {
			continue
		}
		name := ref.Object
		if isBuiltInSchema(name) {
			return mssqlErr(ErrSecurableNotFound, "Cannot drop the schema '%s', because it does not exist or you do not have permission.", name)
		}
		schema := c.GetSchema(name)
		if schema == nil {
			if n.IfExists {
				continue
			}
			return mssqlErr(ErrSecurableNotFound, "Cannot drop the schema '%s', because it does not exist or you do not have permission.", name)
		}
		if len(schema.Tables) > 0 || len(schema.Views) > 0 || len(schema.Procedures) > 0 ||
			len(schema.Functions) > 0 || len(schema.Sequences) > 0 || len(schema.Synonyms) > 0 {
			return mssqlErr(ErrDropSchemaNotEmpty, "Cannot drop schema '%s' because it is being referenced by an object.", name)
		}
		c.removeSchema(schema)
	}
	return nil
}

func datetime2Precision(colType string) (int, bool) {
	norm := normalizeType(colType)
	if norm == "datetime2" {
		return 7, true
	}
	if strings.HasPrefix(norm, "datetime2(") && strings.HasSuffix(norm, ")") {
		inner := strings.TrimSpace(norm[len("datetime2(") : len(norm)-1])
		n, err := strconv.Atoi(inner)
		if err == nil && n >= 0 && n <= 7 {
			return n, true
		}
	}
	return 0, false
}

func isBuiltInSchema(name string) bool {
	switch normalize(name) {
	case "dbo", "guest", "information_schema", "sys",
		"db_owner", "db_accessadmin", "db_securityadmin", "db_ddladmin",
		"db_backupoperator", "db_datareader", "db_datawriter",
		"db_denydatareader", "db_denydatawriter":
		return true
	}
	return false
}

func (c *Catalog) execDropIndex(n *ast.DropStmt) error {
	if n.Names == nil {
		return nil
	}
	type dropIndexTarget struct {
		cat        *Catalog
		schemaName string
		tableName  string
		table      *metadata.TableMetadata
		idx        *metadata.IndexMetadata
	}
	seen := make(map[*metadata.IndexMetadata]bool)
	var targets []dropIndexTarget
	for i, item := range n.Names.Items {
		ref, ok := item.(*ast.TableRef)
		if !ok {
			continue
		}
		var dbName, schemaName, tableName, indexName string
		if n.OnTables != nil && i < len(n.OnTables.Items) && n.OnTables.Items[i] != nil {
			onRef, _ := n.OnTables.Items[i].(*ast.TableRef)
			if onRef != nil {
				dbName = onRef.Database
			}
			schemaName, tableName = tableRefSchemaObject(onRef, c.defaultSchema)
			indexName = ref.Object
		} else {
			dbName = ref.Server
			schemaName = ref.Database
			tableName = ref.Schema
			indexName = ref.Object
			if schemaName == "" {
				schemaName = c.defaultSchema
			}
		}

		targetCat := c
		if dbName != "" {
			targetCat = c.DatabaseCatalog(dbName)
		}
		var table *metadata.TableMetadata
		if targetCat != nil {
			table = targetCat.GetTable(schemaName, tableName)
		}
		var foundIdx *metadata.IndexMetadata
		if table != nil {
			foundIdx, _ = resolveTableIndex(table, indexName, true)
		}
		if foundIdx == nil || seen[foundIdx] {
			if n.IfExists {
				continue
			}
			return mssqlErr(ErrCannotDropNotFound, "Cannot drop the index '%s.%s', because it does not exist or you do not have permission.", tableName, indexName)
		}

		isClustered := strings.EqualFold(foundIdx.Type, "CLUSTERED") || strings.EqualFold(foundIdx.Type, "CLUSTERED COLUMNSTORE")
		if err := checkDropIndexOptions(n.Options, tableName, indexName, isClustered); err != nil {
			return err
		}

		seen[foundIdx] = true
		targets = append(targets, dropIndexTarget{cat: targetCat, schemaName: schemaName, tableName: tableName, table: table, idx: foundIdx})
		if isPrimaryXMLIndex(foundIdx) {
			for _, sec := range table.Indexes {
				if sec != foundIdx && !seen[sec] && isSecondaryXMLIndexOf(sec, foundIdx.Name) {
					seen[sec] = true
					targets = append(targets, dropIndexTarget{cat: targetCat, schemaName: schemaName, tableName: tableName, table: table, idx: sec})
				}
			}
		}
	}
	for _, target := range targets {
		if _, fkName, referenced := indexReferencedByFKExcluding(target.cat, target.schemaName, target.table, target.idx, seen); referenced {
			return mssqlErr(ErrDropIndexReferencedByFK, "An explicit DROP INDEX is not allowed on index '%s.%s'. It is being used for FOREIGN KEY constraint enforcement (%s).", target.tableName, target.idx.Name, fkName)
		}
	}
	for _, target := range targets {
		target.table.Indexes = slices.DeleteFunc(target.table.Indexes, func(idx *metadata.IndexMetadata) bool {
			return idx == target.idx
		})
	}
	return nil
}

func resolveTableIndex(table *metadata.TableMetadata, name string, excludeConstraints bool) (*metadata.IndexMetadata, int) {
	if table == nil || name == "" {
		return nil, -1
	}
	for i, idx := range table.Indexes {
		if excludeConstraints && idx.IsConstraint {
			continue
		}
		if idx.Name == name {
			return idx, i
		}
	}
	var match *metadata.IndexMetadata
	matchPos := -1
	for i, idx := range table.Indexes {
		if excludeConstraints && idx.IsConstraint {
			continue
		}
		if strings.EqualFold(idx.Name, name) {
			if match != nil {
				return nil, -1
			}
			match = idx
			matchPos = i
		}
	}
	return match, matchPos
}

func isSecondaryXMLIndexOf(idx *metadata.IndexMetadata, primaryName string) bool {
	if idx == nil || !strings.EqualFold(idx.Type, "XML") || !strings.HasPrefix(idx.Definition, "USING XML INDEX ") {
		return false
	}
	rest := strings.TrimPrefix(idx.Definition, "USING XML INDEX ")
	parent := rest
	if before, _, ok := strings.Cut(rest, " FOR "); ok {
		parent = before
	}
	return strings.EqualFold(strings.TrimSpace(parent), primaryName)
}

func checkDropIndexOptions(opts *ast.List, tableName, indexName string, isClustered bool) error {
	if opts == nil {
		return nil
	}
	if hasOptionOn(opts, "ONLINE") && !isClustered {
		return mssqlErr(ErrOnlineDropNonClustered, "Only a clustered index can be dropped online.")
	}
	for _, item := range opts.Items {
		s, ok := item.(*ast.String)
		if !ok {
			continue
		}
		upper := strings.ToUpper(strings.TrimSpace(s.Str))
		switch {
		case strings.HasPrefix(upper, "FILESTREAM_ON"):
			return mssqlErr(ErrDropIndexFilestreamOn, "FILESTREAM_ON cannot be specified when a table has no FILESTREAM columns.")
		case strings.HasPrefix(upper, "MAXDOP") || strings.HasPrefix(upper, "MOVE"):
			if !isClustered {
				return mssqlErr(ErrDropIndexClusteredOpt, "Cannot drop non-clustered index '%s.%s' using drop clustered index clause.", tableName, indexName)
			}
		}
	}
	return nil
}

func (c *Catalog) execDropObjects(n *ast.DropStmt, wantKind objectKind, sqlKind string) error {
	if n.Names == nil {
		return nil
	}
	type dropTarget struct {
		cat    *Catalog
		schema *metadata.SchemaMetadata
		name   string
		table  *metadata.TableMetadata
	}
	type otherKey struct {
		cat    *Catalog
		schema string
		obj    string
	}
	seenTables := make(map[*metadata.TableMetadata]bool)
	seenOther := make(map[otherKey]bool)
	var targets []dropTarget
	for _, item := range n.Names.Items {
		ref, ok := item.(*ast.TableRef)
		if !ok {
			continue
		}
		targetCat := c
		if ref.Database != "" {
			targetCat = c.DatabaseCatalog(ref.Database)
		}
		schemaName, objName := tableRefSchemaObject(ref, c.defaultSchema)
		if targetCat != nil && ref.Schema == "" {
			schemaName = targetCat.defaultSchema
		}
		var schema *metadata.SchemaMetadata
		kind := kindNone
		if targetCat != nil {
			schema = targetCat.GetSchema(schemaName)
			kind = targetCat.lookupSchemaObject(schema, objName, false)
		}
		var resolvedTable *metadata.TableMetadata
		resolvedName := ""
		isDup := false
		if kind == kindTable {
			resolvedTable = targetCat.resolveTableInSchema(schema, objName)
			if resolvedTable == nil {
				kind = kindNone
			} else if seenTables[resolvedTable] {
				isDup = true
			}
		} else if kind != kindNone {
			if kind == wantKind {
				var ok bool
				resolvedName, ok = resolveSchemaObjectName(schema, wantKind, objName)
				if !ok {
					kind = kindNone
				} else if seenOther[otherKey{cat: targetCat, schema: schema.Name, obj: resolvedName}] {
					isDup = true
				}
			}
		}
		if kind == kindNone || isDup {
			if n.IfExists {
				continue
			}
			return mssqlErr(ErrCannotDropNotFound, "Cannot drop the %s '%s', because it does not exist or you do not have permission.", strings.ToLower(sqlKind), objName)
		}
		if kind != wantKind {
			return mssqlErr(ErrCannotDropWrongKind, "Cannot use DROP %s with '%s' because '%s' is a %s.", sqlKind, objName, objName, string(kind))
		}
		if wantKind == kindTable {
			periodKey := tableKey{schema: normalizeOrDefault(schema.Name, targetCat.defaultSchema), table: resolvedTable.Name}
			if targetCat.systemVersionedTables[periodKey] {
				return mssqlErr(ErrDropSystemVersionedTable, "Drop table operation failed on table '%s.%s' because it is not a supported operation on system-versioned temporal tables.", schema.Name, resolvedTable.Name)
			}
			if fkName, referenced := hasInboundFKExcludingTables(targetCat, schema.Name, resolvedTable, seenTables); referenced {
				return mssqlErr(ErrDropTableReferencedByFK, "Could not drop object '%s' because it is referenced by a FOREIGN KEY constraint '%s'.", objName, fkName)
			}
			if consumer, bound := schemaBoundTableConsumer(targetCat, schema.Name, resolvedTable.Name); bound {
				return mssqlErr(ErrDropSchemaNotEmpty, "Cannot DROP TABLE '%s' because it is being referenced by object '%s'.", objName, consumer)
			}
			seenTables[resolvedTable] = true
			targets = append(targets, dropTarget{cat: targetCat, schema: schema, name: resolvedTable.Name, table: resolvedTable})
		} else {
			if wantKind == kindView {
				if consumer, bound := schemaBoundObjectConsumer(targetCat, schema.Name, resolvedName, func(sName, oName string) bool {
					return seenOther[otherKey{cat: targetCat, schema: sName, obj: oName}]
				}); bound {
					return mssqlErr(ErrDropSchemaNotEmpty, "Cannot DROP VIEW '%s' because it is being referenced by object '%s'.", objName, consumer)
				}
			} else if wantKind == kindFunction {
				if consumer, bound := schemaBoundFunctionConsumer(targetCat, schema.Name, resolvedName, func(sName, oName string) bool {
					return seenOther[otherKey{cat: targetCat, schema: sName, obj: oName}]
				}); bound {
					return mssqlErr(ErrDropSchemaNotEmpty, "Cannot DROP FUNCTION '%s' because it is being referenced by object '%s'.", objName, consumer)
				}
			} else if wantKind == kindSequence {
				if consumer, bound := sequenceDefaultConsumer(targetCat, schema.Name, resolvedName); bound {
					return mssqlErr(ErrDropSchemaNotEmpty, "Cannot DROP SEQUENCE '%s' because it is being referenced by object '%s'.", objName, consumer)
				}
			}
			seenOther[otherKey{cat: targetCat, schema: schema.Name, obj: resolvedName}] = true
			targets = append(targets, dropTarget{cat: targetCat, schema: schema, name: resolvedName})
		}
	}
	for _, target := range targets {
		if target.table != nil {
			target.cat.removeTable(target.schema, target.table)
		} else {
			target.cat.removeObject(target.schema, wantKind, target.name)
		}
	}
	return nil
}

func resolveExactOrFoldIndex[T any](items []T, nameOf func(T) string, target string) (int, bool) {
	for i, item := range items {
		if nameOf(item) == target {
			return i, true
		}
	}
	match := -1
	for i, item := range items {
		if strings.EqualFold(nameOf(item), target) {
			if match >= 0 {
				return -1, false
			}
			match = i
		}
	}
	return match, match >= 0
}

func resolveSchemaObjectName(schema *metadata.SchemaMetadata, kind objectKind, name string) (string, bool) {
	if schema == nil || name == "" {
		return "", false
	}
	switch kind {
	case kindView:
		if idx, ok := resolveExactOrFoldIndex(schema.Views, func(v *metadata.ViewMetadata) string { return v.Name }, name); ok {
			return schema.Views[idx].Name, true
		}
	case kindProcedure:
		if idx, ok := resolveExactOrFoldIndex(schema.Procedures, func(p *metadata.ProcedureMetadata) string { return p.Name }, name); ok {
			return schema.Procedures[idx].Name, true
		}
	case kindFunction:
		if idx, ok := resolveExactOrFoldIndex(schema.Functions, func(fn *metadata.FunctionMetadata) string { return fn.Name }, name); ok {
			return schema.Functions[idx].Name, true
		}
	case kindSequence:
		if idx, ok := resolveExactOrFoldIndex(schema.Sequences, func(sq *metadata.SequenceMetadata) string { return sq.Name }, name); ok {
			return schema.Sequences[idx].Name, true
		}
	case kindSynonym:
		if idx, ok := resolveExactOrFoldIndex(schema.Synonyms, func(sy *metadata.SynonymMetadata) string { return sy.Name }, name); ok {
			return schema.Synonyms[idx].Name, true
		}
	case kindTrigger:
		var foldedName string
		foldedCount := 0
		for _, t := range schema.Tables {
			for _, tr := range t.Triggers {
				if tr.Name == name {
					return tr.Name, true
				}
				if strings.EqualFold(tr.Name, name) {
					foldedName = tr.Name
					foldedCount++
				}
			}
		}
		for _, v := range schema.Views {
			for _, tr := range v.Triggers {
				if tr.Name == name {
					return tr.Name, true
				}
				if strings.EqualFold(tr.Name, name) {
					foldedName = tr.Name
					foldedCount++
				}
			}
		}
		if foldedCount == 1 {
			return foldedName, true
		}
	}
	return "", false
}

func (c *Catalog) removeTable(schema *metadata.SchemaMetadata, target *metadata.TableMetadata) {
	for i, t := range schema.Tables {
		if t == target {
			delete(c.tableByName, tableKey{schema: schema.Name, table: t.Name})
			periodKey := tableKey{schema: normalizeOrDefault(schema.Name, c.defaultSchema), table: t.Name}
			delete(c.systemTimeTables, periodKey)
			delete(c.systemVersionedTables, periodKey)
			schema.Tables = append(schema.Tables[:i], schema.Tables[i+1:]...)
			return
		}
	}
}

func (c *Catalog) removeObject(schema *metadata.SchemaMetadata, kind objectKind, name string) {
	switch kind {
	case kindTable:
		target := c.resolveTableInSchema(schema, name)
		for i, t := range schema.Tables {
			if t == target || (target == nil && strings.EqualFold(t.Name, name)) {
				delete(c.tableByName, tableKey{schema: schema.Name, table: t.Name})
				periodKey := tableKey{schema: normalizeOrDefault(schema.Name, c.defaultSchema), table: t.Name}
				delete(c.systemTimeTables, periodKey)
				delete(c.systemVersionedTables, periodKey)
				schema.Tables = append(schema.Tables[:i], schema.Tables[i+1:]...)
				return
			}
		}
	case kindView:
		if idx, ok := resolveExactOrFoldIndex(schema.Views, func(v *metadata.ViewMetadata) string { return v.Name }, name); ok {
			schema.Views = append(schema.Views[:idx], schema.Views[idx+1:]...)
		}
	case kindProcedure:
		if idx, ok := resolveExactOrFoldIndex(schema.Procedures, func(p *metadata.ProcedureMetadata) string { return p.Name }, name); ok {
			schema.Procedures = append(schema.Procedures[:idx], schema.Procedures[idx+1:]...)
		}
	case kindFunction:
		if idx, ok := resolveExactOrFoldIndex(schema.Functions, func(fn *metadata.FunctionMetadata) string { return fn.Name }, name); ok {
			schema.Functions = append(schema.Functions[:idx], schema.Functions[idx+1:]...)
		}
	case kindSequence:
		if idx, ok := resolveExactOrFoldIndex(schema.Sequences, func(sq *metadata.SequenceMetadata) string { return sq.Name }, name); ok {
			schema.Sequences = append(schema.Sequences[:idx], schema.Sequences[idx+1:]...)
		}
	case kindSynonym:
		if idx, ok := resolveExactOrFoldIndex(schema.Synonyms, func(sy *metadata.SynonymMetadata) string { return sy.Name }, name); ok {
			schema.Synonyms = append(schema.Synonyms[:idx], schema.Synonyms[idx+1:]...)
		}
	case kindTrigger:
		for _, t := range schema.Tables {
			for i, tr := range t.Triggers {
				if tr.Name == name {
					t.Triggers = append(t.Triggers[:i], t.Triggers[i+1:]...)
					return
				}
			}
		}
		for _, v := range schema.Views {
			for i, tr := range v.Triggers {
				if tr.Name == name {
					v.Triggers = append(v.Triggers[:i], v.Triggers[i+1:]...)
					return
				}
			}
		}
		// Unambiguous folded fallback.
		var matchTable *metadata.TableMetadata
		var matchView *metadata.ViewMetadata
		matchIdx := -1
		matches := 0
		for _, t := range schema.Tables {
			for i, tr := range t.Triggers {
				if strings.EqualFold(tr.Name, name) {
					matchTable = t
					matchView = nil
					matchIdx = i
					matches++
				}
			}
		}
		for _, v := range schema.Views {
			for i, tr := range v.Triggers {
				if strings.EqualFold(tr.Name, name) {
					matchTable = nil
					matchView = v
					matchIdx = i
					matches++
				}
			}
		}
		if matches == 1 {
			if matchTable != nil {
				matchTable.Triggers = append(matchTable.Triggers[:matchIdx], matchTable.Triggers[matchIdx+1:]...)
			} else if matchView != nil {
				matchView.Triggers = append(matchView.Triggers[:matchIdx], matchView.Triggers[matchIdx+1:]...)
			}
		}
	}
}

func sequenceDefaultConsumer(c *Catalog, targetSchema, targetSeq string) (string, bool) {
	for _, s := range c.schemas {
		for _, t := range s.Tables {
			for _, col := range t.Columns {
				if col.Default == "" {
					continue
				}
				if defaultReferencesSequenceInSchema(col.Default, s.Name, c.defaultSchema, targetSchema, targetSeq) {
					name := col.DefaultConstraintName
					if name == "" {
						name = t.Name
					}
					return name, true
				}
			}
		}
	}
	return "", false
}

func defaultReferencesSequenceInSchema(expression, tableSchema, defaultSchemaName, targetSchema, targetSeq string) bool {
	for _, ref := range defaultSequenceRefs(expression) {
		if !strings.EqualFold(ref.Object, targetSeq) {
			continue
		}
		if ref.Schema == "" {
			if strings.EqualFold(targetSchema, tableSchema) || strings.EqualFold(targetSchema, defaultSchemaName) {
				return true
			}
			continue
		}
		if strings.EqualFold(ref.Schema, targetSchema) {
			return true
		}
	}
	return false
}

// --- TRUNCATE TABLE ---

func (c *Catalog) execTruncateTable(n *ast.TruncateStmt) error {
	targetCat := c
	if n.Table != nil && n.Table.Database != "" {
		targetCat = c.DatabaseCatalog(n.Table.Database)
		if targetCat == nil {
			return mssqlErr(ErrTruncateNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", n.Table.Object)
		}
	}
	schemaName, tableName := tableRefSchemaObject(n.Table, targetCat.defaultSchema)
	schema := targetCat.GetSchema(schemaName)
	kind := targetCat.lookupSchemaObject(schema, tableName, false)
	switch kind {
	case kindNone:
		return mssqlErr(ErrTruncateNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	case kindTable:
	default:
		return mssqlErr(ErrTruncateNotTable, "Cannot truncate object '%s' because it is not a table.", tableName)
	}
	table := targetCat.resolveTableInSchema(schema, tableName)
	if table == nil {
		return mssqlErr(ErrTruncateNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	}
	periodKey := tableKey{schema: normalizeOrDefault(schema.Name, targetCat.defaultSchema), table: table.Name}
	if targetCat.systemVersionedTables[periodKey] {
		return mssqlErr(ErrTruncateSystemVersionedTable, "Truncate failed on table '%s.%s' because it is not a supported operation on system-versioned tables.", schema.Name, table.Name)
	}
	if fkName, referenced := hasInboundFK(targetCat, schema.Name, table); referenced {
		return mssqlErr(ErrTruncateReferencedByFK, "Cannot truncate table '%s' because it is being referenced by a FOREIGN KEY constraint '%s'.", tableName, fkName)
	}
	return nil
}
