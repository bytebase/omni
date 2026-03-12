package ddl

import (
	"fmt"

	nodes "github.com/bytebase/omni/pg/ast"
)

// ExecCreateTableAs creates a table or materialized view from a query
// (CREATE TABLE ... AS SELECT or CREATE MATERIALIZED VIEW ... AS SELECT).
//
// pg: src/backend/commands/createas.c — ExecCreateTableAs
func (c *Catalog) ExecCreateTableAs(stmt *nodes.CreateTableAsStmt) error {
	into := stmt.Into
	if into == nil || into.Rel == nil {
		return errInvalidParameterValue("CREATE MATERIALIZED VIEW requires a target name")
	}

	rv := into.Rel
	schemaName := rv.Schemaname
	mvName := rv.Relname

	schema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
	}

	// Check for existing relation.
	if _, exists := schema.Relations[mvName]; exists {
		if stmt.IfNotExists {
			c.addWarning(CodeWarningSkip, fmt.Sprintf("relation %q already exists, skipping", mvName))
			return nil
		}
		return errDuplicateTable(mvName)
	}

	// Check conflict with existing type name.
	if c.typeByName[typeKey{ns: schema.OID, name: mvName}] != nil {
		if stmt.IfNotExists {
			c.addWarning(CodeWarningSkip, fmt.Sprintf("relation %q already exists, skipping", mvName))
			return nil
		}
		return errDuplicateTable(mvName)
	}

	// Analyze the query to infer column types (same approach as DefineView).
	selStmt, ok := stmt.Query.(*nodes.SelectStmt)
	if !ok {
		return errInvalidParameterValue("CREATE MATERIALIZED VIEW requires a SELECT query")
	}

	// Analyze the query — single pipeline.
	// pg: src/backend/commands/createas.c — ExecCreateTableAs
	analyzedQuery, err := c.analyzeSelectStmt(selStmt)
	if err != nil {
		return err
	}

	// Extract column types from analyzed query target entries.
	var cols []*ResultColumn
	if analyzedQuery != nil {
		for _, te := range analyzedQuery.TargetList {
			if te.ResJunk {
				continue
			}
			cols = append(cols, &ResultColumn{
				Name:      te.ResName,
				TypeOID:   te.Expr.exprType(),
				TypeMod:   te.Expr.exprTypMod(),
				Collation: te.Expr.exprCollation(),
			})
		}
	}

	// Apply explicit column names from INTO clause if provided.
	if into.ColNames != nil {
		columnNames := stringListItems(into.ColNames)
		if len(columnNames) != len(cols) {
			// pg: src/backend/commands/createas.c — ExecCreateTableAs (column count mismatch)
			return &Error{
				Code:    CodeSyntaxError,
				Message: fmt.Sprintf("CREATE MATERIALIZED VIEW specifies %d column names, but query produces %d columns", len(columnNames), len(cols)),
			}
		}
		for i, name := range columnNames {
			cols[i].Name = name
		}
	}

	// Build Column objects.
	columns := make([]*Column, len(cols))
	colByName := make(map[string]int, len(cols))
	for i, rc := range cols {
		if _, dup := colByName[rc.Name]; dup {
			return errDuplicateColumn(rc.Name)
		}
		colByName[rc.Name] = i

		typ := c.typeByOID[rc.TypeOID]
		coll := rc.Collation
		if coll == 0 {
			coll = c.typeCollation(rc.TypeOID)
		}
		col := &Column{
			AttNum:    int16(i + 1),
			Name:      rc.Name,
			TypeOID:   rc.TypeOID,
			TypeMod:   rc.TypeMod,
			IsLocal:   true,
			Collation: coll,
		}
		if typ != nil {
			col.Len = typ.Len
			col.ByVal = typ.ByVal
			col.Align = typ.Align
			col.Storage = typ.Storage
		}
		columns[i] = col
	}

	// Allocate OIDs.
	relOID := c.oidGen.Next()
	rowTypeOID := c.oidGen.Next()
	arrayOID := c.oidGen.Next()

	// Apply column names to analyzed query target list.
	if analyzedQuery != nil {
		for i, te := range analyzedQuery.TargetList {
			if i < len(cols) {
				te.ResName = cols[i].Name
			}
		}
	}

	relkind := byte('m')
	if nodes.ObjectType(stmt.Objtype) == nodes.OBJECT_TABLE {
		relkind = 'r'
	}

	rel := &Relation{
		OID:           relOID,
		Name:          mvName,
		Schema:        schema,
		RelKind:       relkind,
		Columns:       columns,
		colByName:     colByName,
		RowTypeOID:    rowTypeOID,
		ArrayOID:      arrayOID,
		AnalyzedQuery: analyzedQuery,
	}

	schema.Relations[mvName] = rel
	c.relationByOID[relOID] = rel

	// Register composite row type.
	rowType := &BuiltinType{
		OID:       rowTypeOID,
		TypeName:  mvName,
		Namespace: schema.OID,
		Len:       -1,
		ByVal:     false,
		Type:      'c',
		Category:  'C',
		IsDefined: true,
		Delim:     ',',
		RelID:     relOID,
		Array:     arrayOID,
		Align:     'd',
		Storage:   'x',
		TypeMod:   -1,
	}
	c.typeByOID[rowTypeOID] = rowType
	c.typeByName[typeKey{ns: schema.OID, name: mvName}] = rowType

	// Register array type.
	arrayType := &BuiltinType{
		OID:       arrayOID,
		TypeName:  fmt.Sprintf("_%s", mvName),
		Namespace: schema.OID,
		Len:       -1,
		ByVal:     false,
		Type:      'b',
		Category:  'A',
		IsDefined: true,
		Delim:     ',',
		Elem:      rowTypeOID,
		Align:     'd',
		Storage:   'x',
		TypeMod:   -1,
	}
	c.typeByOID[arrayOID] = arrayType
	c.typeByName[typeKey{ns: schema.OID, name: arrayType.TypeName}] = arrayType

	// Record dependencies on source tables.
	// pg: src/backend/catalog/dependency.c — recordDependencyOnExpr
	c.recordViewDepsQ(relOID, analyzedQuery)

	return nil
}

// ExecRefreshMatView refreshes a materialized view.
// For pgddl, this is a no-op (no physical data to refresh) — just verifies the matview exists.
//
// pg: src/backend/commands/matview.c — ExecRefreshMatView
func (c *Catalog) ExecRefreshMatView(stmt *nodes.RefreshMatViewStmt) error {
	if stmt.Relation == nil {
		return errInvalidParameterValue("REFRESH MATERIALIZED VIEW requires a target name")
	}
	_, rel, err := c.findRelation(stmt.Relation.Schemaname, stmt.Relation.Relname)
	if err != nil {
		return err
	}
	if rel.RelKind != 'm' {
		return errWrongObjectType(rel.Name, "a materialized view")
	}
	return nil
}
