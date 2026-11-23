package catalog

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/ast"
)

func oraError(code, format string, args ...any) *Error {
	return &Error{Code: code, Message: fmt.Sprintf(format, args...)}
}

func (e *oracleMetadataExtractor) setError(err error) {
	if e.err == nil && err != nil {
		e.err = err
	}
}

// checkSchema validates schema ownership against this catalog and, if another
// schema is named, records an Oracle error with oraCode.
func (e *oracleMetadataExtractor) checkSchema(schema, oraCode string) bool {
	if schema == "" {
		return true
	}
	if e.currentSchema == "" {
		if e.owner != "" && schema != e.owner {
			e.setError(oraError(oraCode, "schema %q does not match catalog schema %q", schema, e.owner))
			return false
		}
		e.currentSchema = schema
		e.explicitSchema = true
		return true
	}
	if schema == e.currentSchema {
		return true
	}
	e.setError(oraError(oraCode, "schema %q does not match catalog schema %q", schema, e.currentSchema))
	return false
}

// isForeignSchema reports whether schema explicitly names an owner other than
// the catalog's current owner.
func (e *oracleMetadataExtractor) isForeignSchema(schema string) bool {
	if schema == "" {
		return false
	}
	owner := e.currentSchema
	if owner == "" {
		owner = e.owner
	}
	return owner != "" && schema != owner
}

// objectExistsInNamespace reports whether name is already taken in Oracle's
// schema object namespace (tables, views, materialized views, sequences,
// private synonyms, functions, procedures, and packages). Indexes and
// constraints live in separate namespaces.
func (e *oracleMetadataExtractor) objectExistsInNamespace(name string) bool {
	return e.tables[name] != nil ||
		e.views[name] != nil ||
		e.materializedViews[name] != nil ||
		e.sequences[name] != nil ||
		e.synonyms[name] != nil ||
		e.functions[name] != nil ||
		e.procedures[name] != nil ||
		e.packages[name] != nil
}

func (e *oracleMetadataExtractor) nonViewObjectExists(name string) bool {
	return e.tables[name] != nil ||
		e.materializedViews[name] != nil ||
		e.sequences[name] != nil ||
		e.synonyms[name] != nil ||
		e.functions[name] != nil ||
		e.procedures[name] != nil ||
		e.packages[name] != nil
}

func (e *oracleMetadataExtractor) nonSynonymObjectExists(name string) bool {
	return e.tables[name] != nil ||
		e.views[name] != nil ||
		e.materializedViews[name] != nil ||
		e.sequences[name] != nil ||
		e.functions[name] != nil ||
		e.procedures[name] != nil ||
		e.packages[name] != nil
}

func (e *oracleMetadataExtractor) indexExists(name string) bool {
	for _, table := range e.tables {
		for _, idx := range table.Indexes {
			if idx.Name == name {
				return true
			}
		}
	}
	for _, mv := range e.materializedViews {
		for _, idx := range mv.Indexes {
			if idx.Name == name {
				return true
			}
		}
	}
	return false
}

func (e *oracleMetadataExtractor) constraintNameExists(name string) bool {
	if name == "" {
		return false
	}
	for _, table := range e.tables {
		if tableHasConstraint(table, name) {
			return true
		}
	}
	return false
}

func tableHasConstraint(table *metadatapb.TableMetadata, name string) bool {
	for _, idx := range table.Indexes {
		if idx.IsConstraint && idx.Name == name {
			return true
		}
		for _, c := range idx.Constraints {
			if c.Name == name {
				return true
			}
		}
	}
	for _, chk := range table.CheckConstraints {
		if chk.Name == name {
			return true
		}
	}
	for _, fk := range table.ForeignKeys {
		if fk.Name == name {
			return true
		}
	}
	return false
}

func tableFindColumn(table *metadatapb.TableMetadata, name string) *metadatapb.ColumnMetadata {
	for _, col := range table.Columns {
		if col.Name == name {
			return col
		}
	}
	return nil
}

func (e *oracleMetadataExtractor) tableHasIncomingForeignKeys(tableName string) bool {
	for _, other := range e.tables {
		if other.Name == tableName {
			continue
		}
		for _, fk := range other.ForeignKeys {
			if fk.ReferencedTable == tableName && (fk.ReferencedSchema == "" || fk.ReferencedSchema == e.currentSchema) {
				return true
			}
		}
	}
	return false
}

func (e *oracleMetadataExtractor) columnsHaveIncomingForeignKeys(tableName string, cols []string) bool {
	if len(cols) == 0 {
		return e.tableHasIncomingForeignKeys(tableName)
	}
	colSet := make(map[string]bool, len(cols))
	for _, c := range cols {
		colSet[strings.Trim(c, `"`)] = true
	}
	for _, other := range e.tables {
		for _, fk := range other.ForeignKeys {
			if fk.ReferencedTable != tableName {
				continue
			}
			if fk.ReferencedSchema != "" && fk.ReferencedSchema != e.currentSchema {
				continue
			}
			if len(fk.ReferencedColumns) == 0 {
				return true
			}
			if len(fk.ReferencedColumns) != len(colSet) {
				continue
			}
			matched := true
			for _, rc := range fk.ReferencedColumns {
				if !colSet[strings.Trim(rc, `"`)] {
					matched = false
					break
				}
			}
			if matched {
				return true
			}
		}
	}
	return false
}

func tableHasUniqueKey(table *metadatapb.TableMetadata, cols []string) bool {
	if table == nil || len(cols) == 0 {
		return false
	}
	want := make([]string, len(cols))
	for i, c := range cols {
		want[i] = strings.Trim(c, `"`)
	}
	slices.Sort(want)
	for _, idx := range table.Indexes {
		if !idx.Unique || !idx.IsConstraint || len(idx.Expressions) != len(want) {
			continue
		}
		got := make([]string, len(idx.Expressions))
		for i, expr := range idx.Expressions {
			got[i] = strings.Trim(expr, `"`)
		}
		slices.Sort(got)
		if slices.Equal(got, want) {
			return true
		}
	}
	return false
}

// resolveLocalSynonym follows a private synonym whose target is in the same
// schema with no database link.
func (e *oracleMetadataExtractor) resolveLocalSynonym(name string) string {
	seen := map[string]bool{}
	cur := name
	for !seen[cur] {
		seen[cur] = true
		syn := e.synonyms[cur]
		if syn == nil || syn.TargetDbLink != "" || syn.TargetName == "" {
			return cur
		}
		if syn.TargetSchema != "" && syn.TargetSchema != e.currentSchema && syn.TargetSchema != e.owner {
			return cur
		}
		cur = syn.TargetName
	}
	return cur
}

func (e *oracleMetadataExtractor) renameSynonym(oldName, newName string) bool {
	syn := e.synonyms[oldName]
	if syn == nil {
		return false
	}
	delete(e.synonyms, oldName)
	syn.Name = newName
	e.synonyms[newName] = syn
	return true
}

func colHasInlineReferences(col *ast.ColumnDef) bool {
	for _, item := range listItems(col.Constraints) {
		if c, ok := item.(*ast.ColumnConstraint); ok && c.Type == ast.CONSTRAINT_FOREIGN {
			return true
		}
	}
	return false
}

func queryHasUnaliasedExpression(query ast.StmtNode) bool {
	sel, ok := query.(*ast.SelectStmt)
	if !ok || sel == nil || sel.TargetList == nil {
		return false
	}
	for _, item := range listItems(sel.TargetList) {
		target, ok := item.(*ast.ResTarget)
		if !ok || target.Name != "" {
			continue
		}
		switch target.Expr.(type) {
		case *ast.ColumnRef, *ast.Star:
			continue
		default:
			return true
		}
	}
	return false
}

func parseSignedIntText(s string) (int64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	v, err := strconv.ParseInt(s, 10, 64)
	return v, err == nil
}

func (e *oracleMetadataExtractor) validateAlterForeignKey(table *metadatapb.TableMetadata, refTable *ast.ObjectName, refColumns *ast.List) bool {
	if refTable == nil || objectName(refTable) == "" {
		return true
	}
	if e.isForeignSchema(refTable.Schema) {
		e.setError(oraError("00942", "table or view %q does not exist", refTable.Schema+"."+objectName(refTable)))
		return false
	}
	targetName := objectName(refTable)
	target := e.tables[targetName]
	if target == nil {
		e.setError(oraError("00942", "table or view %q does not exist", targetName))
		return false
	}
	refCols := stringList(refColumns)
	if len(refCols) == 0 {
		if !hasPrimaryIndex(target) {
			e.setError(oraError("02268", "referenced table %q does not have a primary key", targetName))
			return false
		}
		return true
	}
	if !tableHasUniqueKey(target, refCols) {
		e.setError(oraError("02270", "no matching unique or primary key for this column-list on %q", targetName))
		return false
	}
	return true
}

func (e *oracleMetadataExtractor) extractAlterView(n *ast.AlterViewStmt) {
	name := objectName(n.Name)
	if name == "" {
		return
	}
	if !e.checkSchema(n.Name.Schema, "01917") {
		return
	}
	if e.views[name] == nil {
		if n.IfExists {
			return
		}
		e.setError(oraError("00942", "table or view %q does not exist", name))
		return
	}
	switch strings.ToUpper(n.Action) {
	case "DROP_CONSTRAINT", "DROP_PRIMARY_KEY", "DROP_UNIQUE":
		e.setError(oraError("02443", "cannot drop constraint - nonexistent constraint"))
	case "MODIFY_CONSTRAINT", "MODIFY_PRIMARY_KEY", "MODIFY_UNIQUE":
		e.setError(oraError("25129", "cannot modify constraint - no such constraint"))
	case "READ_ONLY", "READ_WRITE":
		e.setError(oraError("42310", "ALTER VIEW READ ONLY/READ WRITE is only supported on editioning views"))
	default:
		// COMPILE, EDITIONABLE, NONEDITIONABLE, ADD_CONSTRAINT (declarative view constraints).
	}
}

func (e *oracleMetadataExtractor) extractTruncate(n *ast.TruncateStmt) {
	if n.Cluster || n.Table == nil {
		return
	}
	name := objectName(n.Table)
	if name == "" {
		return
	}
	if !e.checkSchema(n.Table.Schema, "00942") {
		return
	}
	if e.views[name] != nil {
		e.setError(oraError("01702", "a view is not appropriate here"))
		return
	}
	if e.sequences[name] != nil {
		e.setError(oraError("02201", "sequence not allowed here"))
		return
	}
	if e.tables[name] == nil {
		e.setError(oraError("00942", "table or view %q does not exist", name))
		return
	}
}

func (e *oracleMetadataExtractor) viewHasColumn(view *metadatapb.ViewMetadata, colName string) bool {
	for _, col := range view.Columns {
		if col.Name == colName {
			return true
		}
	}
	if len(view.Columns) > 0 {
		return false
	}
	inferred, ok := inferredViewQueryColumnNamesForSchema(view.GetDefinition(), e.currentSchemaMetadata())
	if !ok {
		return true
	}
	return slices.Contains(inferred, colName)
}

func (e *oracleMetadataExtractor) currentSchemaMetadata() *metadatapb.SchemaMetadata {
	tables := make([]*metadatapb.TableMetadata, 0, len(e.tables))
	for _, t := range e.tables {
		tables = append(tables, t)
	}
	return &metadatapb.SchemaMetadata{Name: e.currentSchema, Tables: tables}
}

func (e *oracleMetadataExtractor) columnInForeignKey(table *metadatapb.TableMetadata, colName string) bool {
	for _, fk := range table.ForeignKeys {
		if slices.Contains(fk.Columns, colName) {
			return true
		}
	}
	return e.columnsHaveIncomingForeignKeys(table.Name, []string{colName})
}
