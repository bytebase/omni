package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/deparse"
)

func (c *Catalog) createView(stmt *nodes.CreateViewStmt) error {
	db, err := c.resolveDatabase(stmt.Name.Schema)
	if err != nil {
		return err
	}
	key := toLower(stmt.Name.Name)
	// Tables and views share the same namespace in MySQL.
	if _, exists := db.Tables[key]; exists {
		return errDupTable(stmt.Name.Name)
	}
	if _, exists := db.Views[key]; exists {
		if !stmt.OrReplace {
			return errDupTable(stmt.Name.Name)
		}
	}

	// MySQL always sets a definer. Default to `root`@`%` when not specified.
	definer := stmt.Definer
	if definer == "" {
		definer = "`root`@`%`"
	}

	// Analyze the view body for semantic IR (used by lineage, SDL diff).
	// Done before deparseViewSelect because deparse mutates the AST.
	var analyzedQuery *Query
	if stmt.Select != nil {
		q, err := c.AnalyzeSelectStmt(stmt.Select)
		if err == nil {
			analyzedQuery = q
		}
		// Swallow error: view analysis may fail for complex views not yet
		// supported by the analyzer. The view still gets created with
		// Definition text but without AnalyzedQuery.
	}

	// Resolve, rewrite, and deparse the SELECT to produce canonical definition.
	definition, derivedCols := c.deparseViewSelect(stmt.Select, stmt.SelectText, db)

	// Use explicit column list if provided, otherwise derive from SELECT target list.
	hasExplicit := len(stmt.Columns) > 0
	viewCols := stmt.Columns
	if !hasExplicit {
		viewCols = derivedCols
	}
	columnMetadata := inferViewColumnMetadata(stmt.Select, viewCols, db)

	algorithm := stmt.Algorithm
	if algorithm == "" {
		algorithm = "UNDEFINED"
	}
	sqlSecurity := stmt.SqlSecurity
	if sqlSecurity == "" {
		sqlSecurity = "DEFINER"
	}

	db.Views[key] = &View{
		Name:            stmt.Name.Name,
		Database:        db,
		Definition:      definition,
		Algorithm:       algorithm,
		Definer:         definer,
		SqlSecurity:     sqlSecurity,
		CheckOption:     stmt.CheckOption,
		Columns:         viewCols,
		ColumnMetadata:  columnMetadata,
		ExplicitColumns: hasExplicit,
		AnalyzedQuery:   analyzedQuery,
		IsUpdatable:     inferViewUpdatable(stmt.Select, algorithm),
	}
	return nil
}

func (c *Catalog) alterView(stmt *nodes.AlterViewStmt) error {
	db, err := c.resolveDatabase(stmt.Name.Schema)
	if err != nil {
		return err
	}
	key := toLower(stmt.Name.Name)
	// ALTER VIEW requires the view to exist.
	if _, exists := db.Views[key]; !exists {
		return errUnknownTable(db.Name, stmt.Name.Name)
	}

	// MySQL always sets a definer. Default to `root`@`%` when not specified.
	definer := stmt.Definer
	if definer == "" {
		definer = "`root`@`%`"
	}

	// Analyze the view body for semantic IR (used by lineage, SDL diff).
	// Done before deparseViewSelect because deparse mutates the AST.
	var analyzedQuery *Query
	if stmt.Select != nil {
		q, err := c.AnalyzeSelectStmt(stmt.Select)
		if err == nil {
			analyzedQuery = q
		}
	}

	// Resolve, rewrite, and deparse the SELECT to produce canonical definition.
	definition, derivedCols := c.deparseViewSelect(stmt.Select, stmt.SelectText, db)

	// Use explicit column list if provided, otherwise derive from SELECT target list.
	hasExplicit := len(stmt.Columns) > 0
	viewCols := stmt.Columns
	if !hasExplicit {
		viewCols = derivedCols
	}
	columnMetadata := inferViewColumnMetadata(stmt.Select, viewCols, db)

	algorithm := stmt.Algorithm
	if algorithm == "" {
		algorithm = "UNDEFINED"
	}
	sqlSecurity := stmt.SqlSecurity
	if sqlSecurity == "" {
		sqlSecurity = "DEFINER"
	}

	db.Views[key] = &View{
		Name:            stmt.Name.Name,
		Database:        db,
		Definition:      definition,
		Algorithm:       algorithm,
		Definer:         definer,
		SqlSecurity:     sqlSecurity,
		CheckOption:     stmt.CheckOption,
		Columns:         viewCols,
		ColumnMetadata:  columnMetadata,
		ExplicitColumns: hasExplicit,
		AnalyzedQuery:   analyzedQuery,
		IsUpdatable:     inferViewUpdatable(stmt.Select, algorithm),
	}
	return nil
}

type viewRelationInfo struct {
	table    *Table
	optional bool
}

func inferViewColumnMetadata(sel *nodes.SelectStmt, names []string, db *Database) []ViewColumn {
	cols := make([]ViewColumn, len(names))
	for i, name := range names {
		cols[i].Name = name
	}
	if sel == nil {
		return cols
	}
	relations := map[string]viewRelationInfo{}
	for _, from := range sel.From {
		collectViewRelations(from, db, false, relations)
	}
	for i, target := range sel.TargetList {
		if i >= len(cols) {
			break
		}
		rt, ok := target.(*nodes.ResTarget)
		if !ok {
			continue
		}
		cols[i].Nullable = inferViewExprNullable(rt.Val, relations)
	}
	return cols
}

func collectViewRelations(te nodes.TableExpr, db *Database, optional bool, out map[string]viewRelationInfo) {
	switch n := te.(type) {
	case *nodes.TableRef:
		if db == nil {
			return
		}
		tbl := db.GetTable(n.Name)
		if tbl == nil {
			return
		}
		key := n.Name
		if n.Alias != "" {
			key = n.Alias
		}
		out[toLower(key)] = viewRelationInfo{table: tbl, optional: optional}
	case *nodes.JoinClause:
		leftOptional := optional
		rightOptional := optional
		switch n.Type {
		case nodes.JoinLeft, nodes.JoinNaturalLeft:
			rightOptional = true
		case nodes.JoinRight, nodes.JoinNaturalRight:
			leftOptional = true
		}
		collectViewRelations(n.Left, db, leftOptional, out)
		collectViewRelations(n.Right, db, rightOptional, out)
	}
}

func inferViewExprNullable(expr nodes.ExprNode, relations map[string]viewRelationInfo) bool {
	switch e := expr.(type) {
	case *nodes.ColumnRef:
		return inferViewColumnRefNullable(e, relations)
	case *nodes.FuncCallExpr:
		if isViewStringMetadataNullableFunc(e.Name) {
			return true
		}
		for _, arg := range e.Args {
			if inferViewExprNullable(arg, relations) {
				return true
			}
		}
		return false
	case *nodes.ParenExpr:
		return inferViewExprNullable(e.Expr, relations)
	case *nodes.BinaryExpr:
		return inferViewExprNullable(e.Left, relations) || inferViewExprNullable(e.Right, relations)
	case *nodes.UnaryExpr:
		return inferViewExprNullable(e.Operand, relations)
	case *nodes.NullLit:
		return true
	default:
		return false
	}
}

func inferViewColumnRefNullable(cr *nodes.ColumnRef, relations map[string]viewRelationInfo) bool {
	if cr == nil || cr.Star {
		return true
	}
	if cr.Table != "" {
		if rel, ok := relations[toLower(cr.Table)]; ok {
			if col := rel.table.GetColumn(cr.Column); col != nil {
				return rel.optional || col.Nullable
			}
		}
		return true
	}
	found := false
	nullable := true
	for _, rel := range relations {
		col := rel.table.GetColumn(cr.Column)
		if col == nil {
			continue
		}
		if found {
			return true
		}
		found = true
		nullable = rel.optional || col.Nullable
	}
	return nullable
}

func isViewStringMetadataNullableFunc(name string) bool {
	switch strings.ToLower(name) {
	case "concat", "concat_ws", "ifnull", "coalesce":
		return true
	default:
		return false
	}
}

func inferViewUpdatable(sel *nodes.SelectStmt, algorithm string) bool {
	if sel == nil || strings.EqualFold(algorithm, "TEMPTABLE") {
		return false
	}
	if sel.DistinctKind != nodes.DistinctNone && sel.DistinctKind != nodes.DistinctAll {
		return false
	}
	if sel.SetOp != nodes.SetOpNone || len(sel.GroupBy) > 0 || sel.Having != nil {
		return false
	}
	hasAgg := false
	nodes.Inspect(sel, func(n nodes.Node) bool {
		if hasAgg {
			return false
		}
		if fn, ok := n.(*nodes.FuncCallExpr); ok && isAggregateFunc(fn.Name) {
			hasAgg = true
			return false
		}
		return true
	})
	return !hasAgg
}

func (c *Catalog) dropView(stmt *nodes.DropViewStmt) error {
	for _, ref := range stmt.Views {
		db, err := c.resolveDatabase(ref.Schema)
		if err != nil {
			if stmt.IfExists {
				continue
			}
			return err
		}
		key := toLower(ref.Name)
		if _, exists := db.Views[key]; !exists {
			if stmt.IfExists {
				continue
			}
			return errUnknownTable(db.Name, ref.Name)
		}
		delete(db.Views, key)
	}
	return nil
}

// deparseViewSelect resolves, rewrites, and deparses the SELECT AST for a view.
// If the AST is nil (parser didn't produce one), falls back to the raw SelectText.
// Returns the deparsed definition and the derived column names from the resolved
// SELECT target list (used when no explicit column list is specified).
func (c *Catalog) deparseViewSelect(sel *nodes.SelectStmt, rawText string, db *Database) (string, []string) {
	if sel == nil {
		return rawText, nil
	}

	// Build a TableLookup that resolves table names from this database.
	lookup := tableLookupForDB(db)

	// Determine the database charset for CAST resolution.
	charset := db.Charset
	if charset == "" {
		charset = c.defaultCharset
	}

	// Resolve: qualify columns, expand *, normalize JOINs.
	resolver := &deparse.Resolver{
		Lookup:         lookup,
		DefaultCharset: charset,
	}
	resolver.Resolve(sel)

	// Extract column names from the resolved target list.
	derivedCols := extractViewColumns(sel)

	// Rewrite: NOT folding, boolean context wrapping.
	deparse.RewriteSelectStmt(sel)

	// Deparse: AST → canonical SQL text.
	return deparse.DeparseSelect(sel), derivedCols
}

// extractViewColumns extracts column names from a resolved SELECT target list.
// This produces the column list that MySQL would derive for a view.
func extractViewColumns(sel *nodes.SelectStmt) []string {
	if sel == nil {
		return nil
	}
	var cols []string
	for _, target := range sel.TargetList {
		rt, ok := target.(*nodes.ResTarget)
		if !ok {
			continue
		}
		if rt.Name != "" {
			cols = append(cols, canonicalViewDerivedColumnName(rt.Name))
		} else if cr, ok := rt.Val.(*nodes.ColumnRef); ok {
			cols = append(cols, cr.Column)
		} else {
			cols = append(cols, canonicalViewDerivedColumnName(nodeToSQL(rt.Val)))
		}
	}
	return cols
}

func canonicalViewDerivedColumnName(name string) string {
	if strings.ContainsAny(name, "+-*/%") {
		name = strings.ReplaceAll(name, " ", "")
	}
	replacer := strings.NewReplacer(
		" + ", "+",
		" - ", "-",
		" * ", "*",
		" / ", "/",
		" % ", "%",
	)
	name = replacer.Replace(name)
	lower := strings.ToLower(name)
	if strings.HasPrefix(lower, "count(") {
		return "COUNT" + name[len("count"):]
	}
	return name
}

// tableLookupForDB returns a deparse.TableLookup function that resolves table
// and view names from the given database's Tables and Views maps.
func tableLookupForDB(db *Database) deparse.TableLookup {
	return func(tableName string) *deparse.ResolverTable {
		key := toLower(tableName)
		// Try tables first.
		tbl := db.Tables[key]
		if tbl != nil {
			cols := make([]deparse.ResolverColumn, len(tbl.Columns))
			for i, c := range tbl.Columns {
				cols[i] = deparse.ResolverColumn{
					Name:     c.Name,
					Position: c.Position,
				}
			}
			return &deparse.ResolverTable{
				Name:    tbl.Name,
				Columns: cols,
			}
		}
		// Fall back to views.
		v := db.Views[key]
		if v != nil {
			cols := make([]deparse.ResolverColumn, len(v.Columns))
			for i, colName := range v.Columns {
				cols[i] = deparse.ResolverColumn{
					Name:     colName,
					Position: i + 1,
				}
			}
			return &deparse.ResolverTable{
				Name:    v.Name,
				Columns: cols,
			}
		}
		return nil
	}
}

// ShowCreateView produces MySQL 8.0-compatible SHOW CREATE VIEW output.
// Returns "" if the database or view does not exist.
func (c *Catalog) ShowCreateView(database, name string) string {
	db := c.GetDatabase(database)
	if db == nil {
		return ""
	}
	v := db.Views[toLower(name)]
	if v == nil {
		return ""
	}
	return showCreateView(v)
}

// formatDefiner ensures the definer string is backtick-quoted per MySQL 8.0 format.
// Input can be: `root`@`%`, root@%, 'root'@'%', etc.
// Output: `root`@`%`
func formatDefiner(definer string) string {
	// If already formatted with backticks, return as-is.
	if strings.HasPrefix(definer, "`") && strings.Contains(definer, "@") {
		return definer
	}
	// Split on @
	parts := strings.SplitN(definer, "@", 2)
	if len(parts) == 1 {
		// No @ — just backtick-quote the whole thing.
		return "`" + strings.Trim(parts[0], "`'") + "`"
	}
	user := strings.Trim(parts[0], "`'")
	host := strings.Trim(parts[1], "`'")
	return fmt.Sprintf("`%s`@`%s`", user, host)
}

// showCreateView produces the SHOW CREATE VIEW output for a view.
// MySQL 8.0 format:
//
//	CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `view_name` AS select_statement
//	WITH CASCADED CHECK OPTION
func showCreateView(v *View) string {
	var b strings.Builder

	b.WriteString("CREATE")

	// ALGORITHM — MySQL 8.0 always shows ALGORITHM, defaults to UNDEFINED.
	algorithm := v.Algorithm
	if algorithm == "" {
		algorithm = "UNDEFINED"
	}
	b.WriteString(fmt.Sprintf(" ALGORITHM=%s", strings.ToUpper(algorithm)))

	// DEFINER — MySQL 8.0 always shows DEFINER with backtick-quoted user@host.
	if v.Definer != "" {
		b.WriteString(fmt.Sprintf(" DEFINER=%s", formatDefiner(v.Definer)))
	}

	// SQL SECURITY — MySQL 8.0 always shows SQL SECURITY, defaults to DEFINER.
	sqlSecurity := v.SqlSecurity
	if sqlSecurity == "" {
		sqlSecurity = "DEFINER"
	}
	b.WriteString(fmt.Sprintf(" SQL SECURITY %s", strings.ToUpper(sqlSecurity)))

	// VIEW name
	b.WriteString(fmt.Sprintf(" VIEW `%s`", v.Name))

	// Column list (only if explicitly specified by user in CREATE VIEW).
	if v.ExplicitColumns && len(v.Columns) > 0 {
		cols := make([]string, len(v.Columns))
		for i, c := range v.Columns {
			cols[i] = fmt.Sprintf("`%s`", c)
		}
		b.WriteString(fmt.Sprintf(" (%s)", strings.Join(cols, ",")))
	}

	// AS select_statement
	b.WriteString(" AS ")
	b.WriteString(v.Definition)

	// WITH CHECK OPTION
	if v.CheckOption != "" {
		b.WriteString(fmt.Sprintf(" WITH %s CHECK OPTION", strings.ToUpper(v.CheckOption)))
	}

	return b.String()
}
