package catalog

import (
	"fmt"
	"math/big"
	"slices"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/ast"
	oracleparser "github.com/bytebase/omni/oracle/parser"
)

type viewConstraint struct {
	name    string
	kind    ast.ConstraintType
	columns []string
}

func cloneViewConstraints(src []viewConstraint) []viewConstraint {
	if len(src) == 0 {
		return nil
	}
	out := make([]viewConstraint, len(src))
	for i, vc := range src {
		out[i] = viewConstraint{
			name:    vc.name,
			kind:    vc.kind,
			columns: append([]string(nil), vc.columns...),
		}
	}
	return out
}

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

func (e *oracleMetadataExtractor) findIndex(name string) *metadatapb.IndexMetadata {
	for _, table := range e.tables {
		for _, idx := range table.Indexes {
			if idx.Name == name {
				return idx
			}
		}
	}
	for _, mv := range e.materializedViews {
		for _, idx := range mv.Indexes {
			if idx.Name == name {
				return idx
			}
		}
	}
	return nil
}

func (e *oracleMetadataExtractor) explicitIndexExists(name string) bool {
	for _, table := range e.tables {
		for _, idx := range table.Indexes {
			if idx.Name != name {
				continue
			}
			if idx.IsConstraint && len(idx.Constraints) == 0 && e.syntheticConstraints[table.Name+"\x00"+idx.Name] {
				continue
			}
			return true
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

func (e *oracleMetadataExtractor) disambiguateSyntheticIndex(name string) {
	if name == "" {
		return
	}
	for _, table := range e.tables {
		if e.syntheticConstraints[table.Name+"\x00"+name] {
			e.disambiguateSyntheticConstraint(table, name)
		}
	}
}

func (e *oracleMetadataExtractor) constraintNameExists(name string) bool {
	if name == "" {
		return false
	}
	for _, table := range e.tables {
		if tableHasConstraint(table, name) && !e.syntheticConstraints[table.Name+"\x00"+name] {
			return true
		}
	}
	for _, constraints := range e.viewConstraints {
		for _, vc := range constraints {
			if vc.name == name {
				return true
			}
		}
	}
	return false
}

// resolveConstraintName returns explicitName when provided (disambiguating any
// same-table synthetic constraint that happened to use that fallback name), or
// generates a non-colliding synthetic name recorded in syntheticConstraints.
func (e *oracleMetadataExtractor) resolveConstraintName(table *metadatapb.TableMetadata, explicitName, fallback string) string {
	if explicitName != "" {
		e.disambiguateSyntheticConstraint(table, explicitName)
		return explicitName
	}
	name := fallback
	for suffix := 1; tableHasConstraint(table, name) || e.findIndex(name) != nil; suffix++ {
		name = fmt.Sprintf("%s_%d", fallback, suffix)
	}
	e.syntheticConstraints[table.Name+"\x00"+name] = true
	return name
}

func (e *oracleMetadataExtractor) disambiguateSyntheticConstraint(table *metadatapb.TableMetadata, explicitName string) {
	if table == nil || explicitName == "" {
		return
	}
	key := table.Name + "\x00" + explicitName
	if !e.syntheticConstraints[key] {
		return
	}
	newName := explicitName + "_SYS1"
	for suffix := 2; tableHasConstraint(table, newName) || e.findIndex(newName) != nil; suffix++ {
		newName = fmt.Sprintf("%s_SYS%d", explicitName, suffix)
	}
	e.renameConstraint(table, explicitName, newName)
	delete(e.syntheticConstraints, key)
	e.syntheticConstraints[table.Name+"\x00"+newName] = true
}

// constraintKeyColumns returns the referenced key columns if name identifies
// a primary key or unique constraint on table, and nil otherwise (including
// CHECK constraints, foreign keys, and standalone indexes that share the name).
func constraintKeyColumns(table *metadatapb.TableMetadata, name string) []string {
	for _, idx := range table.Indexes {
		if len(idx.Constraints) > 0 {
			for _, c := range idx.Constraints {
				if c.Name == name {
					if len(c.Columns) > 0 {
						return append([]string(nil), c.Columns...)
					}
					return append([]string(nil), idx.Expressions...)
				}
			}
			continue
		}
		if idx.IsConstraint && idx.Name == name {
			return append([]string(nil), idx.Expressions...)
		}
	}
	return nil
}

func tableHasConstraint(table *metadatapb.TableMetadata, name string) bool {
	for _, idx := range table.Indexes {
		if len(idx.Constraints) > 0 {
			for _, c := range idx.Constraints {
				if c.Name == name {
					return true
				}
			}
		} else if idx.IsConstraint && idx.Name == name {
			return true
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

func columnsMatchSet(a, wantSorted []string) bool {
	if len(a) != len(wantSorted) {
		return false
	}
	got := make([]string, len(a))
	for i, expr := range a {
		got[i] = strings.Trim(expr, `"`)
	}
	slices.Sort(got)
	return slices.Equal(got, wantSorted)
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
		if len(idx.Constraints) > 0 {
			for _, c := range idx.Constraints {
				constraintCols := c.Columns
				if len(constraintCols) == 0 {
					constraintCols = idx.Expressions
				}
				if columnsMatchSet(constraintCols, want) {
					return true
				}
			}
			continue
		}
		if (idx.Primary || idx.Unique) && idx.IsConstraint && columnsMatchSet(idx.Expressions, want) {
			return true
		}
	}
	return false
}

func tableHasUniqueConstraint(table *metadatapb.TableMetadata, cols []string) bool {
	if table == nil || len(cols) == 0 {
		return false
	}
	want := make([]string, len(cols))
	for i, c := range cols {
		want[i] = strings.Trim(c, `"`)
	}
	slices.Sort(want)
	for _, idx := range table.Indexes {
		if len(idx.Constraints) > 0 {
			for _, c := range idx.Constraints {
				if c.Primary {
					continue
				}
				constraintCols := c.Columns
				if len(constraintCols) == 0 {
					constraintCols = idx.Expressions
				}
				if columnsMatchSet(constraintCols, want) {
					return true
				}
			}
			continue
		}
		if idx.Unique && !idx.Primary && idx.IsConstraint && columnsMatchSet(idx.Expressions, want) {
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

func parseSignedBigInt(s string) (*big.Int, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, false
	}
	v, ok := new(big.Int).SetString(s, 10)
	return v, ok
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
	case "ADD_CONSTRAINT":
		if n.Constraint != nil {
			if n.Constraint.Name != "" && e.constraintNameExists(n.Constraint.Name) {
				e.setError(oraError("02264", "name %q already used by an existing constraint", n.Constraint.Name))
				return
			}
			e.viewConstraints[name] = append(e.viewConstraints[name], viewConstraint{
				name:    n.Constraint.Name,
				kind:    n.Constraint.Type,
				columns: stringList(n.Constraint.Columns),
			})
		}
	case "DROP_CONSTRAINT":
		if !e.removeViewConstraintByName(name, n.ConstraintName) {
			e.setError(oraError("02443", "cannot drop constraint - nonexistent constraint"))
		}
	case "DROP_PRIMARY_KEY":
		if !e.removeViewConstraintByKind(name, ast.CONSTRAINT_PRIMARY, nil) {
			e.setError(oraError("02443", "cannot drop constraint - nonexistent constraint"))
		}
	case "DROP_UNIQUE":
		cols := extractAlterViewUniqueColumns(e.nodeText(n))
		if !e.removeViewConstraintByKind(name, ast.CONSTRAINT_UNIQUE, cols) {
			e.setError(oraError("02443", "cannot drop constraint - nonexistent constraint"))
		}
	case "MODIFY_CONSTRAINT":
		if !e.hasViewConstraintByName(name, n.ConstraintName) {
			e.setError(oraError("25129", "cannot modify constraint - no such constraint"))
		}
	case "MODIFY_PRIMARY_KEY":
		if !e.hasViewConstraintByKind(name, ast.CONSTRAINT_PRIMARY, nil) {
			e.setError(oraError("25129", "cannot modify constraint - no such constraint"))
		}
	case "MODIFY_UNIQUE":
		cols := extractAlterViewUniqueColumns(e.nodeText(n))
		if !e.hasViewConstraintByKind(name, ast.CONSTRAINT_UNIQUE, cols) {
			e.setError(oraError("25129", "cannot modify constraint - no such constraint"))
		}
	case "READ_ONLY", "READ_WRITE":
		if !e.editioningViews[name] {
			e.setError(oraError("42310", "ALTER VIEW READ ONLY/READ WRITE is only supported on editioning views"))
		}
	default:
		// COMPILE, EDITIONABLE, NONEDITIONABLE, ANNOTATIONS.
	}
}

func (e *oracleMetadataExtractor) hasViewConstraintByName(viewName, constraintName string) bool {
	if constraintName == "" {
		return false
	}
	for _, vc := range e.viewConstraints[viewName] {
		if vc.name == constraintName {
			return true
		}
	}
	return false
}

func (e *oracleMetadataExtractor) removeViewConstraintByName(viewName, constraintName string) bool {
	if constraintName == "" {
		return false
	}
	constraints := e.viewConstraints[viewName]
	for i, vc := range constraints {
		if vc.name == constraintName {
			e.viewConstraints[viewName] = append(constraints[:i], constraints[i+1:]...)
			return true
		}
	}
	return false
}

func viewConstraintMatchesKind(vc viewConstraint, kind ast.ConstraintType, cols []string) bool {
	if vc.kind != kind {
		return false
	}
	if kind != ast.CONSTRAINT_UNIQUE || len(cols) == 0 || len(vc.columns) == 0 {
		return true
	}
	want := make([]string, len(cols))
	for i, c := range cols {
		want[i] = strings.Trim(c, `"`)
	}
	slices.Sort(want)
	return columnsMatchSet(vc.columns, want)
}

func (e *oracleMetadataExtractor) hasViewConstraintByKind(viewName string, kind ast.ConstraintType, cols []string) bool {
	for _, vc := range e.viewConstraints[viewName] {
		if viewConstraintMatchesKind(vc, kind, cols) {
			return true
		}
	}
	return false
}

func (e *oracleMetadataExtractor) removeViewConstraintByKind(viewName string, kind ast.ConstraintType, cols []string) bool {
	constraints := e.viewConstraints[viewName]
	for i, vc := range constraints {
		if viewConstraintMatchesKind(vc, kind, cols) {
			e.viewConstraints[viewName] = append(constraints[:i], constraints[i+1:]...)
			return true
		}
	}
	return false
}

func extractAlterViewUniqueColumns(text string) []string {
	lexer := oracleparser.NewLexer(text)
	seenUnique := false
	inParens := false
	var cols []string
	for {
		tok := lexer.NextToken()
		if tok.Type == 0 {
			break
		}
		if !seenUnique {
			if strings.EqualFold(tok.Str, "UNIQUE") {
				seenUnique = true
			}
			continue
		}
		if !inParens {
			if tok.Type == '(' {
				inParens = true
			}
			continue
		}
		if tok.Type == ')' {
			break
		}
		if tok.Type == ',' {
			continue
		}
		if tok.Str != "" {
			cols = append(cols, tok.Str)
		}
	}
	return cols
}

func extractViewInlineConstraints(schemaText string, n *ast.CreateViewStmt) []viewConstraint {
	if n == nil || n.Name == nil || n.Query == nil {
		return nil
	}
	start := n.Name.Loc.End
	end := n.Loc.End
	if qloc := ast.NodeLoc(n.Query); qloc.Start > start && qloc.Start <= len(schemaText) {
		end = qloc.Start
	}
	if start < 0 || end <= start || end > len(schemaText) {
		return nil
	}
	header := schemaText[start:end]
	openIdx := strings.IndexByte(header, '(')
	if openIdx < 0 {
		return nil
	}
	lexer := oracleparser.NewLexer(header[openIdx:])
	depth := 0
	started := false
	var itemTokens []oracleparser.Token
	var out []viewConstraint
	flush := func() {
		if len(itemTokens) == 0 {
			return
		}
		first := strings.ToUpper(itemTokens[0].Str)
		// Out-of-line constraints start with CONSTRAINT/PRIMARY/UNIQUE/FOREIGN and are already in n.Constraints.
		if first == "CONSTRAINT" || first == "PRIMARY" || first == "UNIQUE" || first == "FOREIGN" {
			itemTokens = itemTokens[:0]
			return
		}
		colName := itemTokens[0].Str
		pendingName := ""
		for i := 1; i < len(itemTokens); i++ {
			word := strings.ToUpper(itemTokens[i].Str)
			switch word {
			case "CONSTRAINT":
				if i+1 < len(itemTokens) {
					pendingName = itemTokens[i+1].Str
					i++
				}
			case "PRIMARY":
				out = append(out, viewConstraint{name: pendingName, kind: ast.CONSTRAINT_PRIMARY, columns: []string{colName}})
				pendingName = ""
			case "UNIQUE":
				out = append(out, viewConstraint{name: pendingName, kind: ast.CONSTRAINT_UNIQUE, columns: []string{colName}})
				pendingName = ""
			case "REFERENCES":
				out = append(out, viewConstraint{name: pendingName, kind: ast.CONSTRAINT_FOREIGN, columns: []string{colName}})
				pendingName = ""
			}
		}
		itemTokens = itemTokens[:0]
	}
	for {
		tok := lexer.NextToken()
		if tok.Type == 0 {
			break
		}
		if tok.Type == '(' {
			depth++
			started = true
			continue
		}
		if tok.Type == ')' {
			depth--
			if started && depth == 0 {
				flush()
				break
			}
			continue
		}
		if depth == 1 {
			if tok.Type == ',' {
				flush()
				continue
			}
			itemTokens = append(itemTokens, tok)
		}
	}
	return out
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
	if e.tables[name] == nil && e.materializedViews[name] == nil {
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
