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

// isUnresolvedExternalSchema reports whether schema explicitly names a schema
// outside this catalog (including when the catalog has no owner yet).
func (e *oracleMetadataExtractor) isUnresolvedExternalSchema(schema string) bool {
	if schema == "" {
		return false
	}
	owner := e.currentSchema
	if owner == "" {
		owner = e.owner
	}
	return owner == "" || schema != owner
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

// tableFindColumnOrPath looks up an exact column name first; if name is a
// dotted object-attribute path (such as O.X), it resolves against the root
// column.
func tableFindColumnOrPath(table *metadatapb.TableMetadata, name string) *metadatapb.ColumnMetadata {
	if col := tableFindColumn(table, name); col != nil {
		return col
	}
	if dot := strings.IndexByte(name, '.'); dot > 0 {
		root := strings.Trim(name[:dot], `"`)
		return tableFindColumn(table, root)
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

func findStandaloneUniqueIndex(table *metadatapb.TableMetadata, cols []string) *metadatapb.IndexMetadata {
	if table == nil || len(cols) == 0 {
		return nil
	}
	want := make([]string, len(cols))
	for i, c := range cols {
		want[i] = strings.Trim(c, `"`)
	}
	slices.Sort(want)
	for _, idx := range table.Indexes {
		if idx.IsConstraint || len(idx.Constraints) > 0 || !idx.Unique || idx.Primary {
			continue
		}
		if columnsMatchSet(idx.Expressions, want) {
			return idx
		}
	}
	return nil
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

// resolveSynonymTarget resolves name with Oracle precedence:
//  1. A direct non-synonym object in this schema wins over any synonym.
//  2. A private synonym in this schema whose target is in this schema (no DB link).
//  3. When schema == "" (unqualified), a public synonym whose target is in this
//     schema (no DB link). Schema-qualified references never resolve public synonyms.
func (e *oracleMetadataExtractor) resolveSynonymTarget(name, schema string) string {
	if e.nonSynonymObjectExists(name) {
		return name
	}
	if resolved := e.resolveLocalSynonym(name); resolved != name || e.synonyms[name] != nil {
		return resolved
	}
	if schema != "" {
		return name
	}
	syn := e.publicSynonyms[name]
	if syn == nil || syn.TargetDbLink != "" || syn.TargetName == "" {
		return name
	}
	if syn.TargetSchema != "" && syn.TargetSchema != e.currentSchema && syn.TargetSchema != e.owner {
		return name
	}
	return e.resolveLocalSynonym(syn.TargetName)
}

// alterActionsFollowSynonym reports whether every clause in an ALTER TABLE
// statement is one that follows private and public synonyms in Oracle 23
// (MOVE and MODIFY PARTITION BY).
func alterActionsFollowSynonym(actions *ast.List) bool {
	items := listItems(actions)
	if len(items) == 0 {
		return false
	}
	for _, item := range items {
		cmd, ok := item.(*ast.AlterTableCmd)
		if !ok || (cmd.Action != ast.AT_MOVE && cmd.Action != ast.AT_MODIFY_TO_PARTITIONED) {
			return false
		}
	}
	return true
}

func (e *oracleMetadataExtractor) extractCreateDomain(n *ast.CreateDomainStmt) {
	if n == nil || n.Name == nil {
		return
	}
	name := objectName(n.Name)
	if name == "" {
		return
	}
	if !e.checkSchema(n.Name.Schema, "01031") {
		return
	}
	if n.IfNotExists && e.domains[name] {
		return
	}
	if !n.IfNotExists && !n.OrReplace && (e.domains[name] || e.objectExistsInNamespace(name)) {
		e.setError(oraError("00955", "name %q is already used by an existing object", name))
		return
	}
	e.domains[name] = true
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
	if !ok || sel == nil {
		return false
	}
	for sel.Larg != nil {
		sel = sel.Larg
	}
	if sel.TargetList == nil {
		return false
	}
	for _, item := range listItems(sel.TargetList) {
		target, ok := item.(*ast.ResTarget)
		if !ok || target.Name != "" || target.Expr == nil {
			continue
		}
		expr := target.Expr
		for {
			p, ok := expr.(*ast.ParenExpr)
			if !ok || p == nil {
				break
			}
			expr = p.Expr
		}
		switch expr.(type) {
		case *ast.ColumnRef, *ast.Star, *ast.PseudoColumn:
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

func indexColumnRefKey(colRef *ast.ColumnRef) string {
	var parts []string
	if colRef.Schema != "" {
		parts = append(parts, colRef.Schema)
	}
	if colRef.Table != "" {
		parts = append(parts, colRef.Table)
	}
	if colRef.Column != "" {
		parts = append(parts, colRef.Column)
	}
	return strings.Join(parts, ".")
}

// indexColumnRefExists validates a ColumnRef in CREATE INDEX against table,
// supporting both table/alias-qualified columns (t.col) and object-attribute
// paths (o.x or o.nested.x).
func indexColumnRefExists(table *metadatapb.TableMetadata, alias string, colRef *ast.ColumnRef) bool {
	if colRef.Schema == "" && colRef.Table == "" {
		return tableFindColumn(table, colRef.Column) != nil
	}
	if colRef.Schema == "" {
		if tableFindColumn(table, colRef.Table) != nil {
			return true
		}
		if colRef.Table == table.Name || (alias != "" && colRef.Table == alias) {
			return tableFindColumn(table, colRef.Column) != nil
		}
		return false
	}
	return tableFindColumn(table, colRef.Schema) != nil
}

func (e *oracleMetadataExtractor) validateAlterForeignKey(table *metadatapb.TableMetadata, refTable *ast.ObjectName, refColumns *ast.List) bool {
	if refTable == nil || objectName(refTable) == "" {
		return true
	}
	if e.serverConfig.ValidateExternalSchemaRefs && e.isUnresolvedExternalSchema(refTable.Schema) {
		e.setError(oraError("00942", "table or view %q.%q does not exist", refTable.Schema, objectName(refTable)))
		return false
	}
	if e.isForeignSchema(refTable.Schema) {
		// The catalog models a single schema, so it cannot assert absence of
		// an external schema's table; preserve cross-schema foreign keys.
		return true
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

func (e *oracleMetadataExtractor) isStandardMaxStringSize() bool {
	return strings.EqualFold(strings.TrimSpace(e.serverConfig.MaxStringSize), "STANDARD")
}

func isApplicationSharing(sharing string) bool {
	s := strings.TrimSpace(sharing)
	return s != "" && !strings.EqualFold(s, "NONE")
}

func isColumnEncrypt(encrypt string) bool {
	s := strings.TrimSpace(encrypt)
	return s != "" && !strings.EqualFold(s, "DECRYPT")
}

func (e *oracleMetadataExtractor) columnReferencedByIncomingForeignKey(table *metadatapb.TableMetadata, colName string) bool {
	if table == nil || colName == "" {
		return false
	}
	pkCols := primaryKeyColumns(table)
	for _, other := range e.tables {
		for _, fk := range other.ForeignKeys {
			if fk.ReferencedTable != table.Name {
				continue
			}
			if fk.ReferencedSchema != "" && fk.ReferencedSchema != e.currentSchema {
				continue
			}
			refCols := fk.ReferencedColumns
			if len(refCols) == 0 {
				refCols = pkCols
			}
			for _, rc := range refCols {
				if strings.Trim(rc, `"`) == colName {
					return true
				}
			}
		}
	}
	return false
}

func alterTableHasCollationClause(actions *ast.List) bool {
	for _, item := range listItems(actions) {
		action, ok := item.(*ast.AlterTableCmd)
		if !ok {
			continue
		}
		switch action.Action {
		case ast.AT_ADD_COLUMN, ast.AT_MODIFY_COLUMN:
			if action.ColumnDef != nil && action.ColumnDef.Collation != "" {
				return true
			}
			for _, citem := range listItems(action.ColumnDefs) {
				if col, ok := citem.(*ast.ColumnDef); ok && col.Collation != "" {
					return true
				}
			}
		case ast.AT_ALTER_PROPERTY:
			if strings.EqualFold(action.Subtype, "DEFAULT COLLATION") {
				return true
			}
		}
	}
	return false
}

// inferVirtualColumnType infers the Oracle data type of an untyped virtual
// column numeric literal or arithmetic expression using Oracle's numeric type
// precedence (BINARY_DOUBLE > BINARY_FLOAT > NUMBER). Non-numeric expressions
// (such as UPPER(C)) stay untyped so DatabaseDefinition preserves AS (expr).
func inferVirtualColumnType(expr ast.ExprNode) string {
	for expr != nil {
		switch v := expr.(type) {
		case *ast.ParenExpr:
			expr = v.Expr
		case *ast.UnaryExpr:
			if v.Op == "+" || v.Op == "-" {
				expr = v.Operand
				continue
			}
			return ""
		case *ast.NumberLiteral:
			return inferNumberLiteralType(v.Val)
		case *ast.BinaryExpr:
			switch v.Op {
			case "+", "-", "*", "/":
				lt := inferVirtualColumnType(v.Left)
				rt := inferVirtualColumnType(v.Right)
				if lt == "" || rt == "" {
					return ""
				}
				if lt == "BINARY_DOUBLE" || rt == "BINARY_DOUBLE" {
					return "BINARY_DOUBLE"
				}
				if lt == "BINARY_FLOAT" || rt == "BINARY_FLOAT" {
					return "BINARY_FLOAT"
				}
				return "NUMBER"
			}
			return ""
		default:
			return ""
		}
	}
	return ""
}

func inferNumberLiteralType(val string) string {
	val = strings.TrimSpace(val)
	if val == "" {
		return ""
	}
	switch val[len(val)-1] {
	case 'f', 'F':
		return "BINARY_FLOAT"
	case 'd', 'D':
		return "BINARY_DOUBLE"
	default:
		return "NUMBER"
	}
}

// virtualColumnEditionRefs extracts edition names referenced in virtual-column
// evaluation/usability clauses (EVALUATE USING EDITION <name>, UNUSABLE BEFORE
// EDITION <name>, UNUSABLE BEGINNING WITH EDITION <name>) at top-level depth.
// CURRENT EDITION and NULL EDITION are not named editions and are ignored.
func virtualColumnEditionRefs(colText string) []string {
	if colText == "" {
		return nil
	}
	type lexTok struct {
		str    string
		typ    int
		quoted bool
	}
	lexer := oracleparser.NewLexer(colText)
	depth := 0
	var top []lexTok
	for {
		tok := lexer.NextToken()
		if tok.Type == 0 {
			break
		}
		if tok.Type == '(' {
			depth++
			continue
		}
		if tok.Type == ')' {
			if depth > 0 {
				depth--
			}
			continue
		}
		if depth != 0 {
			continue
		}
		quoted := tok.Loc >= 0 && tok.Loc < len(colText) && colText[tok.Loc] == '"'
		top = append(top, lexTok{str: tok.Str, typ: tok.Type, quoted: quoted})
	}
	isWord := func(idx int, word string) bool {
		return idx >= 0 && idx < len(top) && !top[idx].quoted && strings.EqualFold(top[idx].str, word)
	}
	var editions []string
	for i := 0; i+1 < len(top); i++ {
		if !isWord(i, "EDITION") {
			continue
		}
		// Named edition clauses always follow USING, BEFORE, or WITH:
		//   EVALUATE USING EDITION <name>
		//   UNUSABLE BEFORE EDITION <name>
		//   UNUSABLE BEGINNING WITH EDITION <name>
		// In contrast, CURRENT EDITION and NULL EDITION put CURRENT/NULL
		// immediately before EDITION.
		if isWord(i-1, "USING") || isWord(i-1, "BEFORE") || isWord(i-1, "WITH") {
			name := top[i+1].str
			if !top[i+1].quoted {
				name = strings.ToUpper(name)
			}
			if name != "" {
				editions = append(editions, name)
			}
			i++
		}
	}
	return editions
}

func (e *oracleMetadataExtractor) validateVirtualColumnEditions(col *ast.ColumnDef) bool {
	if col == nil || col.Virtual == nil {
		return true
	}
	for _, ed := range virtualColumnEditionRefs(e.nodeText(col)) {
		if !e.editions[ed] {
			e.setError(oraError("38802", "edition %q does not exist", ed))
			return false
		}
	}
	return true
}

func (e *oracleMetadataExtractor) extractAdminDDL(n *ast.AdminDDLStmt) {
	if n == nil {
		return
	}
	if n.ObjectType == ast.OBJECT_EDITION && strings.EqualFold(n.Action, "CREATE") {
		name := objectName(n.Name)
		if name == "" {
			return
		}
		if e.editions[name] {
			if n.IfNotExists {
				return
			}
			e.setError(oraError("00955", "name %q is already used by an existing object", name))
			return
		}
		e.editions[name] = true
	}
}

func (e *oracleMetadataExtractor) extractAlterSession(n *ast.AlterSessionStmt) {
	if n == nil {
		return
	}
	if strings.EqualFold(n.Feature, "SHARD_DDL") {
		switch strings.ToUpper(n.Action) {
		case "ENABLE":
			e.serverConfig.DisableShardDDL = false
		case "DISABLE":
			e.serverConfig.DisableShardDDL = true
		}
	}
}

// parsePeriodForClause parses `[... PERIOD FOR] valid_time_col [(start_col, end_col)]`.
// When hasPrefix is true, it scans forward for the PERIOD FOR keywords first
// (as in ALTER TABLE ... ADD (PERIOD FOR ...)); when false, it parses from the
// start of the text (as in CreateTableStmt's DDLOption value after PERIOD FOR).
func parsePeriodForClause(text string, hasPrefix bool) (period, startCol, endCol string, ok bool) {
	lexer := oracleparser.NewLexer(text)
	nextTok := func() (string, oracleparser.Token, bool) {
		tok := lexer.NextToken()
		if tok.Type == 0 {
			return "", tok, false
		}
		quoted := tok.Loc >= 0 && tok.Loc < len(text) && text[tok.Loc] == '"'
		name := tok.Str
		if !quoted {
			name = strings.ToUpper(name)
		}
		return name, tok, quoted
	}
	if hasPrefix {
		found := false
		prevPeriod := false
		for {
			word, tok, quoted := nextTok()
			if tok.Type == 0 {
				break
			}
			if prevPeriod && !quoted && word == "FOR" {
				found = true
				break
			}
			prevPeriod = !quoted && word == "PERIOD"
		}
		if !found {
			return "", "", "", false
		}
	}
	var tok oracleparser.Token
	period, tok, _ = nextTok()
	if tok.Type == 0 || period == "" {
		return "", "", "", false
	}
	peek := lexer.NextToken()
	if peek.Type == '(' {
		s, st, _ := nextTok()
		if st.Type == 0 || s == "" {
			return period, "", "", true
		}
		comma := lexer.NextToken()
		if comma.Type != ',' {
			return period, "", "", true
		}
		en, et, _ := nextTok()
		if et.Type == 0 || en == "" {
			return period, "", "", true
		}
		return period, s, en, true
	}
	return period, "", "", true
}

func escapeOracleDoubleQuote(ident string) string {
	return strings.ReplaceAll(ident, `"`, `""`)
}

// addPeriodConstraint records the temporal-validity check constraint Oracle
// synthesizes for PERIOD FOR <period> [(<start>, <end>)].
func (e *oracleMetadataExtractor) addPeriodConstraint(table *metadatapb.TableMetadata, period, startCol, endCol string) {
	if table == nil || period == "" {
		return
	}
	var expr string
	if startCol == "" || endCol == "" {
		p := escapeOracleDoubleQuote(period)
		expr = fmt.Sprintf(`"%s_START" < "%s_END"`, p, p)
	} else {
		s := escapeOracleDoubleQuote(startCol)
		en := escapeOracleDoubleQuote(endCol)
		p := escapeOracleDoubleQuote(period)
		expr = fmt.Sprintf(`("%s" < "%s") and ("%s" > 0)`, s, en, p)
	}
	checkName := e.checkConstraintName(table, "", fmt.Sprintf("CHK_%s_%s", table.Name, period))
	e.syntheticConstraints[table.Name+"\x00"+checkName] = true
	e.reserveCheckConstraintName(table, checkName)
	table.CheckConstraints = append(table.CheckConstraints, &metadatapb.CheckConstraintMetadata{
		Name:       checkName,
		Expression: expr,
	})
}

// periodForCheckConstraint recognizes an Oracle temporal-validity check
// constraint whose synthetic period column(s) are not present in table.Columns
// and returns the corresponding PERIOD FOR clause so DatabaseDefinition and
// GenerateMigration emit valid PERIOD FOR DDL instead of a CHECK referencing
// nonexistent hidden columns.
func periodForCheckConstraint(table *metadatapb.TableMetadata, check *metadatapb.CheckConstraintMetadata) (periodName, clause string, ok bool) {
	if table == nil || check == nil {
		return "", "", false
	}
	expr := strings.TrimSpace(check.GetExpression())
	if expr == "" {
		return "", "", false
	}
	type tokItem struct {
		str    string
		typ    int
		quoted bool
	}
	lexer := oracleparser.NewLexer(expr)
	var toks []tokItem
	for {
		tok := lexer.NextToken()
		if tok.Type == 0 {
			break
		}
		quoted := tok.Loc >= 0 && tok.Loc < len(expr) && expr[tok.Loc] == '"'
		toks = append(toks, tokItem{str: tok.Str, typ: tok.Type, quoted: quoted})
	}
	// Pattern 1 (implicit endpoints): "<P>_START" < "<P>_END"
	if len(toks) == 3 && toks[0].quoted && toks[1].typ == '<' && toks[2].quoted {
		startName, endName := toks[0].str, toks[2].str
		if p1, ok1 := strings.CutSuffix(startName, "_START"); ok1 && p1 != "" && endName == p1+"_END" {
			if tableFindColumn(table, startName) == nil && tableFindColumn(table, endName) == nil {
				return p1, fmt.Sprintf(`PERIOD FOR "%s"`, escapeOracleDoubleQuote(p1)), true
			}
		}
	}
	// Pattern 2 (explicit endpoints): ("<S>" < "<E>") and ("<P>" > 0)
	if len(toks) == 11 &&
		toks[0].typ == '(' && toks[1].quoted && toks[2].typ == '<' && toks[3].quoted && toks[4].typ == ')' &&
		!toks[5].quoted && strings.EqualFold(toks[5].str, "AND") &&
		toks[6].typ == '(' && toks[7].quoted && toks[8].typ == '>' && !toks[9].quoted && toks[9].str == "0" && toks[10].typ == ')' {
		startCol, endCol, period := toks[1].str, toks[3].str, toks[7].str
		if period != "" && tableFindColumn(table, period) == nil {
			return period, fmt.Sprintf(`PERIOD FOR "%s" ("%s", "%s")`,
				escapeOracleDoubleQuote(period),
				escapeOracleDoubleQuote(startCol),
				escapeOracleDoubleQuote(endCol)), true
		}
	}
	return "", "", false
}

func (e *oracleMetadataExtractor) dropPeriodConstraint(table *metadatapb.TableMetadata, periodName string) {
	if table == nil || periodName == "" {
		return
	}
	var remaining []*metadatapb.CheckConstraintMetadata
	for _, check := range table.CheckConstraints {
		if p, _, ok := periodForCheckConstraint(table, check); ok && p == periodName {
			delete(e.syntheticConstraints, table.Name+"\x00"+check.Name)
			continue
		}
		remaining = append(remaining, check)
	}
	table.CheckConstraints = remaining
}

// extractSetUnusedSpec parses SET UNUSED { COLUMN col | (col [, col ...]) }
// [CASCADE CONSTRAINTS] [INVALIDATE] [ONLINE] from the action slice text.
func extractSetUnusedSpec(actionText, fallbackCol string) (columns []string, cascade bool) {
	if strings.TrimSpace(actionText) == "" {
		if fallbackCol != "" {
			return []string{fallbackCol}, false
		}
		return nil, false
	}
	lexer := oracleparser.NewLexer(actionText)
	// Skip leading SET UNUSED.
	for {
		tok := lexer.NextToken()
		if tok.Type == 0 {
			break
		}
		quoted := tok.Loc >= 0 && tok.Loc < len(actionText) && actionText[tok.Loc] == '"'
		if !quoted && strings.EqualFold(tok.Str, "UNUSED") {
			break
		}
	}
	tok := lexer.NextToken()
	if tok.Type == '(' {
		for {
			item := lexer.NextToken()
			if item.Type == 0 || item.Type == ')' {
				break
			}
			if item.Type == ',' {
				continue
			}
			name := item.Str
			quoted := item.Loc >= 0 && item.Loc < len(actionText) && actionText[item.Loc] == '"'
			if !quoted {
				name = strings.ToUpper(name)
			}
			if name != "" {
				columns = append(columns, name)
			}
		}
	} else {
		quoted := tok.Loc >= 0 && tok.Loc < len(actionText) && actionText[tok.Loc] == '"'
		if !quoted && strings.EqualFold(tok.Str, "COLUMN") {
			tok = lexer.NextToken()
			quoted = tok.Loc >= 0 && tok.Loc < len(actionText) && actionText[tok.Loc] == '"'
		}
		if tok.Type != 0 && tok.Str != "" {
			name := tok.Str
			if !quoted {
				name = strings.ToUpper(name)
			}
			columns = append(columns, name)
		}
	}
	for {
		t := lexer.NextToken()
		if t.Type == 0 {
			break
		}
		quoted := t.Loc >= 0 && t.Loc < len(actionText) && actionText[t.Loc] == '"'
		if !quoted && strings.EqualFold(t.Str, "CASCADE") {
			cascade = true
		}
	}
	if len(columns) == 0 && fallbackCol != "" {
		columns = []string{fallbackCol}
	}
	return columns, cascade
}

func (e *oracleMetadataExtractor) sequenceInitialCurrent(inc, startVal, minVal, maxVal *big.Int, hasStart, hasMin, hasMax bool) *big.Int {
	var start *big.Int
	switch {
	case hasStart && startVal != nil:
		start = new(big.Int).Set(startVal)
	case inc.Sign() > 0:
		if hasMin && minVal != nil {
			start = new(big.Int).Set(minVal)
		} else {
			start = big.NewInt(1)
		}
	default:
		if hasMax && maxVal != nil {
			start = new(big.Int).Set(maxVal)
		} else {
			start = big.NewInt(-1)
		}
	}
	return new(big.Int).Sub(start, inc)
}
