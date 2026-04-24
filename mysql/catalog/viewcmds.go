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
	columnMetadata, err := inferViewColumnMetadata(stmt.Select, viewCols, db)
	if err != nil {
		return err
	}

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
	columnMetadata, err := inferViewColumnMetadata(stmt.Select, viewCols, db)
	if err != nil {
		return err
	}

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

type viewCollationDerivation int

const (
	viewDerivationExplicit viewCollationDerivation = iota
	viewDerivationImplicit
	viewDerivationCoercible
)

type viewCollationInfo struct {
	charset    string
	collation  string
	derivation viewCollationDerivation
	valid      bool
}

func inferViewColumnMetadata(sel *nodes.SelectStmt, names []string, db *Database) ([]ViewColumn, error) {
	cols := make([]ViewColumn, len(names))
	for i, name := range names {
		cols[i].Name = name
	}
	if sel == nil {
		return cols, nil
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
		if ci, err := inferViewExprCollation(rt.Val, relations); err != nil {
			return nil, err
		} else if ci.valid {
			cols[i].Charset = ci.charset
			cols[i].Collation = ci.collation
		}
	}
	if err := validateViewSelectCollations(sel, relations); err != nil {
		return nil, err
	}
	return cols, nil
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

func validateViewSelectCollations(sel *nodes.SelectStmt, relations map[string]viewRelationInfo) error {
	if sel == nil {
		return nil
	}
	for _, expr := range []nodes.ExprNode{sel.Where, sel.Having} {
		if expr == nil {
			continue
		}
		if _, err := inferViewExprCollation(expr, relations); err != nil {
			return err
		}
	}
	for _, expr := range sel.GroupBy {
		if _, err := inferViewExprCollation(expr, relations); err != nil {
			return err
		}
	}
	for _, item := range sel.OrderBy {
		if item == nil {
			continue
		}
		if _, err := inferViewExprCollation(item.Expr, relations); err != nil {
			return err
		}
	}
	return nil
}

func inferViewExprCollation(expr nodes.ExprNode, relations map[string]viewRelationInfo) (viewCollationInfo, error) {
	switch e := expr.(type) {
	case *nodes.ColumnRef:
		return inferViewColumnRefCollation(e, relations), nil
	case *nodes.StringLit:
		charset := normalizeCharsetName(strings.TrimPrefix(toLower(e.Charset), "_"))
		if charset == "" {
			charset = "utf8mb4"
		}
		return newViewCollation(charset, "", viewDerivationCoercible), nil
	case *nodes.CollateExpr:
		child, err := inferViewExprCollation(e.Expr, relations)
		if err != nil {
			return viewCollationInfo{}, err
		}
		collation := toLower(e.Collation)
		charset := normalizeCharsetName(charsetForCollation(collation))
		if charset == "" && child.valid {
			charset = child.charset
		}
		return newViewCollation(charset, collation, viewDerivationExplicit), nil
	case *nodes.ConvertExpr:
		if e.Charset != "" {
			return newViewCollation(normalizeCharsetName(e.Charset), "", viewDerivationImplicit), nil
		}
		return inferViewExprCollation(e.Expr, relations)
	case *nodes.FuncCallExpr:
		return inferViewFuncCollation(e, relations)
	case *nodes.ParenExpr:
		return inferViewExprCollation(e.Expr, relations)
	case *nodes.BinaryExpr:
		return inferViewBinaryCollation(e, relations)
	case *nodes.UnaryExpr:
		return inferViewExprCollation(e.Operand, relations)
	default:
		return viewCollationInfo{}, nil
	}
}

func inferViewFuncCollation(fn *nodes.FuncCallExpr, relations map[string]viewRelationInfo) (viewCollationInfo, error) {
	name := strings.ToLower(fn.Name)
	switch name {
	case "concat", "concat_ws":
		var out viewCollationInfo
		for _, arg := range fn.Args {
			ci, err := inferViewExprCollation(arg, relations)
			if err != nil {
				return viewCollationInfo{}, err
			}
			agg, err := aggregateViewCollations(out, ci, name)
			if err != nil {
				return viewCollationInfo{}, err
			}
			out = agg
		}
		return out, nil
	case "repeat", "lpad", "rpad":
		if len(fn.Args) == 0 {
			return viewCollationInfo{}, nil
		}
		return inferViewExprCollation(fn.Args[0], relations)
	default:
		var out viewCollationInfo
		for _, arg := range fn.Args {
			ci, err := inferViewExprCollation(arg, relations)
			if err != nil {
				return viewCollationInfo{}, err
			}
			agg, err := aggregateViewCollations(out, ci, name)
			if err != nil {
				return viewCollationInfo{}, err
			}
			out = agg
		}
		return out, nil
	}
}

func inferViewBinaryCollation(expr *nodes.BinaryExpr, relations map[string]viewRelationInfo) (viewCollationInfo, error) {
	left, err := inferViewExprCollation(expr.Left, relations)
	if err != nil {
		return viewCollationInfo{}, err
	}
	right, err := inferViewExprCollation(expr.Right, relations)
	if err != nil {
		return viewCollationInfo{}, err
	}
	op := binaryOpToString(expr.Op)
	if isViewCollationComparisonOp(expr.Op) {
		if _, err := aggregateViewCollationsForComparison(left, right, op); err != nil {
			return viewCollationInfo{}, err
		}
		return viewCollationInfo{}, nil
	}
	return aggregateViewCollations(left, right, op)
}

func inferViewColumnRefCollation(cr *nodes.ColumnRef, relations map[string]viewRelationInfo) viewCollationInfo {
	if cr == nil || cr.Star {
		return viewCollationInfo{}
	}
	if cr.Table != "" {
		if rel, ok := relations[toLower(cr.Table)]; ok {
			return viewColumnCollation(rel.table.GetColumn(cr.Column))
		}
		return viewCollationInfo{}
	}
	var out viewCollationInfo
	found := false
	for _, rel := range relations {
		ci := viewColumnCollation(rel.table.GetColumn(cr.Column))
		if !ci.valid {
			continue
		}
		if found {
			return viewCollationInfo{}
		}
		found = true
		out = ci
	}
	return out
}

func viewColumnCollation(col *Column) viewCollationInfo {
	if col == nil || col.Charset == "" || !isStringType(col.DataType) {
		return viewCollationInfo{}
	}
	return newViewCollation(col.Charset, col.Collation, viewDerivationImplicit)
}

func newViewCollation(charset, collation string, derivation viewCollationDerivation) viewCollationInfo {
	charset = normalizeCharsetName(charset)
	collation = toLower(collation)
	if charset == "" && collation != "" {
		charset = normalizeCharsetName(charsetForCollation(collation))
	}
	if collation == "" && charset != "" {
		collation = defaultCollationForCharset[toLower(charset)]
	}
	if charset == "" || collation == "" {
		return viewCollationInfo{}
	}
	return viewCollationInfo{
		charset:    charset,
		collation:  collation,
		derivation: derivation,
		valid:      true,
	}
}

func aggregateViewCollations(left, right viewCollationInfo, op string) (viewCollationInfo, error) {
	if !left.valid {
		return right, nil
	}
	if !right.valid {
		return left, nil
	}
	if strings.EqualFold(left.collation, right.collation) {
		return left, nil
	}
	if left.derivation == viewDerivationExplicit && right.derivation == viewDerivationExplicit {
		return viewCollationInfo{}, errIllegalMixCollations(left, right, op)
	}
	if left.derivation != right.derivation {
		if left.derivation < right.derivation {
			return left, nil
		}
		return right, nil
	}
	if isCharsetSuperset(left.charset, right.charset) {
		return left, nil
	}
	if isCharsetSuperset(right.charset, left.charset) {
		return right, nil
	}
	// MySQL 8.0 permits several same-charset string function mixes that the
	// comparison path rejects. Preserve the left side for view metadata.
	if strings.EqualFold(left.charset, right.charset) {
		return left, nil
	}
	return viewCollationInfo{}, errIllegalMixCollations(left, right, op)
}

func aggregateViewCollationsForComparison(left, right viewCollationInfo, op string) (viewCollationInfo, error) {
	if !left.valid || !right.valid || strings.EqualFold(left.collation, right.collation) {
		if left.valid {
			return left, nil
		}
		return right, nil
	}
	if left.derivation == viewDerivationExplicit && right.derivation == viewDerivationExplicit {
		return viewCollationInfo{}, errIllegalMixCollations(left, right, op)
	}
	if left.derivation == right.derivation && strings.EqualFold(left.charset, right.charset) {
		return viewCollationInfo{}, errIllegalMixCollations(left, right, op)
	}
	return aggregateViewCollations(left, right, op)
}

func isViewCollationComparisonOp(op nodes.BinaryOp) bool {
	switch op {
	case nodes.BinOpEq, nodes.BinOpNe, nodes.BinOpLt, nodes.BinOpGt, nodes.BinOpLe, nodes.BinOpGe, nodes.BinOpNullSafeEq:
		return true
	default:
		return false
	}
}

func isCharsetSuperset(charset, other string) bool {
	charset = normalizeCharsetName(charset)
	other = normalizeCharsetName(other)
	if charset == other {
		return true
	}
	switch charset {
	case "utf8mb4":
		return other == "latin1" || other == "ascii" || other == "utf8mb3" || other == "utf8"
	case "utf8mb3":
		return other == "ascii"
	case "latin1":
		return other == "ascii"
	default:
		return false
	}
}

func errIllegalMixCollations(left, right viewCollationInfo, op string) error {
	return &Error{
		Code:     1267,
		SQLState: "HY000",
		Message: fmt.Sprintf("Illegal mix of collations (%s,%s) and (%s,%s) for operation '%s'",
			left.collation, viewDerivationName(left.derivation),
			right.collation, viewDerivationName(right.derivation), op),
	}
}

func viewDerivationName(derivation viewCollationDerivation) string {
	switch derivation {
	case viewDerivationExplicit:
		return "EXPLICIT"
	case viewDerivationImplicit:
		return "IMPLICIT"
	case viewDerivationCoercible:
		return "COERCIBLE"
	default:
		return "NONE"
	}
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
