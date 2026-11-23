package catalog

import (
	"errors"
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/mariadb/deparse"
)

// createView applies CREATE VIEW. source is the text the statement was parsed
// from, or "" when there is none (see deparseViewSelect). checkBody refuses a
// view body the engine rejects (a `?` parameter marker); it is true on the
// CREATE VIEW path (Exec) and false on a DefineView sync, whose body the engine
// already accepted.
func (c *Catalog) createView(stmt *nodes.CreateViewStmt, source string, checkBody bool) error {
	db, err := c.resolveDatabase(stmt.Name.Schema)
	if err != nil {
		return err
	}
	// A view body is judged before the view's name is looked up on the Exec
	// path: mariadb:11.8.8 answers 1054 for a SELECT naming a column no
	// relation of the body has, even when the name already belongs to a table
	// (1050) or a view. A synced view's body is the engine's already, so it is
	// not judged here.
	if checkBody {
		if err := validateViewPreAnalysis(stmt.Select); err != nil {
			return err
		}
	}
	analyzedQuery, bodyErr := c.analyzeViewBodyFor(stmt.Select, source, db, checkBody)
	if checkBody && bodyErr != nil {
		return bodyErr
	}
	if checkBody {
		if err := c.validateViewPostAnalysis(stmt.Select, db); err != nil {
			return err
		}
	}
	key := toLower(stmt.Name.Name)
	// Tables and views share the same namespace in MySQL.
	if _, exists := db.Tables[key]; exists || db.Sequences[key] != nil {
		return errDupTable(stmt.Name.Name)
	}
	if _, exists := db.Views[key]; exists {
		if !stmt.OrReplace {
			return errDupTable(stmt.Name.Name)
		}
	}

	// A `?` parameter marker is not allowed in a view body; MariaDB rejects
	// CREATE VIEW with one. The grammar permits it (it is a LIMIT operand), so
	// the rejection lives here, not in the parser, and must fail Exec rather
	// than be swallowed as the analyzer's other failures are.
	if checkBody && stmt.Select != nil && containsParamMarker(stmt.Select) {
		return &Error{
			Code:     1351,
			SQLState: "HY000",
			Message:  "View's SELECT contains a '?' (parameter marker)",
		}
	}

	// MySQL always sets a definer. Default to `root`@`%` when not specified.
	definer := stmt.Definer
	if definer == "" {
		definer = "`root`@`%`"
	}

	// Resolve, rewrite, and deparse the SELECT to produce canonical definition.
	definition, derivedCols := c.deparseViewSelect(stmt.Select, stmt.SelectText, source, db, stmt.Columns)

	// Use explicit column list if provided, otherwise derive from SELECT target list.
	hasExplicit := len(stmt.Columns) > 0
	viewCols := stmt.Columns
	if !hasExplicit {
		viewCols = derivedCols
	}
	if checkBody {
		if err := checkDuplicateViewColumns(viewCols); err != nil {
			return err
		}
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
		SelectAST:       stmt.Select,
		AnalyzedQuery:   analyzedQuery,
		IsUpdatable:     inferViewUpdatable(stmt.Select, algorithm),
	}
	return nil
}

// analyzeViewBody analyzes a view body for the semantic IR lineage and SDL
// diff use, and reports semantic failures MariaDB's CREATE VIEW reports too.
// Alias-bearing clauses (GROUP BY, HAVING, ORDER BY) and unsupported analyzer
// expressions leave the view without AnalyzedQuery rather than refusing a view
// the engine accepts.
func (c *Catalog) analyzeViewBody(sel *nodes.SelectStmt, source string) (*Query, error) {
	return c.analyzeViewBodyFor(sel, source, c.GetDatabase(c.CurrentDatabase()), true)
}

func (c *Catalog) analyzeViewBodyFor(sel *nodes.SelectStmt, source string, viewDB *Database, checkFuncs bool) (*Query, error) {
	if sel == nil {
		return nil, nil
	}
	savedValidating, savedDB := c.validatingViewBody, c.validatingViewDB
	if checkFuncs {
		c.validatingViewBody = true
		c.validatingViewDB = viewDB
	}
	q, err := c.analyzeSelectStmtWithSource(sel, source)
	c.validatingViewBody, c.validatingViewDB = savedValidating, savedDB
	if err == nil {
		return q, nil
	}
	var refusal *Error
	if errors.As(err, &refusal) {
		switch refusal.Code {
		case ErrNoTablesUsed, ErrUnknownTable, ErrOperandColumns,
			ErrNoSuchFunction, ErrFuncInexistentNameCollision, ErrFunctionNotDefined:
			return nil, refusal
		case ErrNoSuchTable:
			if !isSystemDatabaseTableErr(refusal) {
				return nil, refusal
			}
		case ErrNoSuchColumn:
			if !isAliasBearingColumnErr(refusal) {
				return nil, refusal
			}
		}
	}
	return nil, nil
}

// retagNoSuchColumn updates the clause context of an ErrNoSuchColumn error so
// analyzeViewBody can tell alias-bearing clauses (GROUP BY, HAVING, ORDER BY)
// apart from field-list / WHERE column lookups.
func retagNoSuchColumn(err error, context string) error {
	var ce *Error
	if errors.As(err, &ce) && ce.Code == ErrNoSuchColumn && strings.HasSuffix(ce.Message, "in 'field list'") {
		return &Error{
			Code:     ce.Code,
			SQLState: ce.SQLState,
			Message:  strings.TrimSuffix(ce.Message, "in 'field list'") + "in '" + context + "'",
		}
	}
	return err
}

func isAliasBearingColumnErr(ce *Error) bool {
	return strings.HasSuffix(ce.Message, "in 'group statement'") ||
		strings.HasSuffix(ce.Message, "in 'having clause'") ||
		strings.HasSuffix(ce.Message, "in 'order clause'")
}

func isSystemDatabaseTableErr(ce *Error) bool {
	lower := strings.ToLower(ce.Message)
	for _, sysDB := range []string{"'information_schema.", "'performance_schema.", "'mysql.", "'sys."} {
		if strings.Contains(lower, sysDB) {
			return true
		}
	}
	return false
}

func checkDuplicateViewColumns(cols []string) error {
	seen := make(map[string]bool, len(cols))
	for _, col := range cols {
		key := toLower(col)
		if seen[key] {
			return errDupColumn(col)
		}
		seen[key] = true
	}
	return nil
}

// validateViewPreAnalysis checks structural restrictions that MariaDB 11.8
// enforces on a view body before resolving tables or columns:
//   - [NO]WAIT on FOR UPDATE / LOCK IN SHARE MODE (1350)
//   - schema-qualified special-syntax TRIM/SUBSTRING outside mariadb_schema/oracle_schema (1128)
func validateViewPreAnalysis(sel *nodes.SelectStmt) error {
	if sel == nil {
		return nil
	}
	var err error
	nodes.Inspect(sel, func(n nodes.Node) bool {
		if err != nil {
			return false
		}
		switch x := n.(type) {
		case *nodes.SelectStmt:
			if x.ForUpdate != nil && (x.ForUpdate.HasWait || x.ForUpdate.NoWait) {
				err = errViewSelectClause("[NO]WAIT")
				return false
			}
		case *nodes.FuncCallExpr:
			if x.Schema != "" && x.SpecialSyntax &&
				!strings.EqualFold(x.Schema, "mariadb_schema") && !strings.EqualFold(x.Schema, "oracle_schema") {
				raw := x.RawName
				if raw == "" {
					raw = x.Name
				}
				err = errFunctionNotDefined(x.Schema + "." + raw)
				return false
			}
		}
		return true
	})
	return err
}

// validateViewPostAnalysis checks semantic restrictions MariaDB 11.8 reports
// after resolving columns:
//   - sequence expressions (NEXT/PREVIOUS VALUE FOR, NEXTVAL/LASTVAL/SETVAL)
//     naming a regular table or view rather than a sequence (4089)
//   - non-constant ESCAPE operand in LIKE or AGAINST operand in MATCH (1210)
func (c *Catalog) validateViewPostAnalysis(sel *nodes.SelectStmt, defaultDB *Database) error {
	if sel == nil {
		return nil
	}
	var err error
	checkSeqRef := func(schema, name string) {
		targetDB := defaultDB
		if schema != "" {
			targetDB = c.GetDatabase(schema)
		}
		if targetDB == nil {
			return
		}
		key := toLower(name)
		if targetDB.Sequences[key] != nil {
			return
		}
		if targetDB.Tables[key] != nil || targetDB.Views[key] != nil {
			err = errNotSequence(4089, targetDB.Name, name)
		}
	}
	nodes.Inspect(sel, func(n nodes.Node) bool {
		if err != nil {
			return false
		}
		switch x := n.(type) {
		case *nodes.NextValueForExpr:
			if x.Sequence != nil {
				checkSeqRef(x.Sequence.Schema, x.Sequence.Name)
			}
		case *nodes.PreviousValueForExpr:
			if x.Sequence != nil {
				checkSeqRef(x.Sequence.Schema, x.Sequence.Name)
			}
		case *nodes.FuncCallExpr:
			if funcErr := c.validateViewFuncExists(x, defaultDB); funcErr != nil {
				err = funcErr
				return false
			}
			if x.Schema == "" && !x.SpacedParen && len(x.Args) >= 1 {
				switch strings.ToUpper(x.Name) {
				case "NEXTVAL", "LASTVAL", "SETVAL":
					if cr, ok := x.Args[0].(*nodes.ColumnRef); ok && !cr.Star && cr.Schema == "" {
						checkSeqRef(cr.Table, cr.Column)
					}
				}
			}
		case *nodes.LikeExpr:
			if x.Escape != nil && exprContainsColumnRef(x.Escape) {
				err = errWrongArguments("ESCAPE")
			}
		case *nodes.MatchExpr:
			if x.Against != nil && exprContainsColumnRef(x.Against) {
				err = errWrongArguments("AGAINST")
			}
		}
		return err == nil
	})
	return err
}

func (c *Catalog) validateViewFuncExists(fc *nodes.FuncCallExpr, viewDB *Database) error {
	if fc == nil || fc.Name == "" || c.loadingSDL {
		return nil
	}
	lower := toLower(fc.Name)
	displayName := fc.RawName
	if displayName == "" {
		displayName = fc.Name
	}
	// Three-part package routine call: db.package.routine(). MariaDB always
	// reports 1305 (ER_SP_DOES_NOT_EXIST), never 1630, when the routine is absent.
	if fc.Schema != "" && fc.Package != "" {
		return errNoSuchFunction(fc.Schema + "." + fc.Package + "." + displayName)
	}
	// Two-part schema-qualified call: schema.routine().
	if fc.Schema != "" {
		if strings.EqualFold(fc.Schema, "mariadb_schema") || strings.EqualFold(fc.Schema, "oracle_schema") {
			if isMariaDBSchemaBuiltinFunction(fc, lower) {
				return nil
			}
			return errFunctionNotDefined(displayName)
		}
		targetDB := c.GetDatabase(fc.Schema)
		if targetDB != nil && targetDB.Functions[lower] != nil {
			return nil
		}
		qualName := fc.Schema + "." + displayName
		if isMariaDBFuncCollisionName(lower) {
			return errFuncInexistentNameCollision(qualName)
		}
		return errNoSuchFunction(qualName)
	}
	// One-part unqualified call.
	if fc.SpacedParen {
		return nil
	}
	if isMariaDBBuiltinFunction(lower, fc.NameQuoted) {
		return nil
	}
	curDB := c.GetDatabase(c.CurrentDatabase())
	if curDB == nil {
		curDB = viewDB
	}
	if curDB != nil && curDB.Functions[lower] != nil {
		return nil
	}
	qualName := displayName
	if curDB != nil && curDB.Name != "" {
		qualName = curDB.Name + "." + displayName
	}
	if isMariaDBFuncCollisionName(lower) {
		return errFuncInexistentNameCollision(qualName)
	}
	return errNoSuchFunction(qualName)
}

func exprContainsColumnRef(expr nodes.ExprNode) bool {
	found := false
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if found {
			return false
		}
		if cr, ok := n.(*nodes.ColumnRef); ok && !cr.Star {
			found = true
			return false
		}
		return true
	})
	return found
}

// containsParamMarker reports whether the subtree holds a `?` parameter
// marker, at any depth including nested subqueries. It also inspects the rows
// of INSERT ... VALUES and of a VALUES statement, which nodes.Inspect's walker
// skips, so a marker in a stored routine or trigger body cannot slip past.
func containsParamMarker(n nodes.Node) bool {
	found := false
	nodes.Inspect(n, func(x nodes.Node) bool {
		if found {
			return false
		}
		switch e := x.(type) {
		case *nodes.ParamMarkerExpr:
			found = true
		case *nodes.InsertStmt:
			found = rowsHaveParamMarker(e.Values)
		case *nodes.ValuesStmt:
			found = rowsHaveParamMarker(e.Rows)
		}
		return !found
	})
	return found
}

// rowsHaveParamMarker reports whether any expression in the rows of an
// INSERT ... VALUES or a VALUES statement holds a `?` marker. Those rows are
// not reached by nodes.Inspect, so containsParamMarker checks them here.
func rowsHaveParamMarker(rows [][]nodes.ExprNode) bool {
	for _, row := range rows {
		for _, v := range row {
			if v != nil && containsParamMarker(v) {
				return true
			}
		}
	}
	return false
}

// alterView applies ALTER VIEW; source is as for createView.
func (c *Catalog) alterView(stmt *nodes.AlterViewStmt, source string) error {
	db, err := c.resolveDatabase(stmt.Name.Schema)
	if err != nil {
		return err
	}
	key := toLower(stmt.Name.Name)
	// ALTER VIEW requires the view to exist, and answers differently when the
	// name belongs to a table: mariadb:11.8.8 says 1347 "'db.a' is not of type
	// 'VIEW'" for a table and 1146 "Table 'db.a' doesn't exist" for no object.
	if _, exists := db.Views[key]; !exists {
		if db.Tables[key] != nil || db.Sequences[key] != nil {
			return errWrongObject(db.Name, stmt.Name.Name, "VIEW")
		}
		return errNoSuchTable(db.Name, stmt.Name.Name)
	}

	// A `?` parameter marker is not allowed in a view body; the engine rejects
	// ALTER VIEW with one just as it rejects CREATE VIEW. ALTER VIEW comes only
	// from Exec, so its body is always checked.
	if stmt.Select != nil && containsParamMarker(stmt.Select) {
		return &Error{
			Code:     1351,
			SQLState: "HY000",
			Message:  "View's SELECT contains a '?' (parameter marker)",
		}
	}

	// MySQL always sets a definer. Default to `root`@`%` when not specified.
	definer := stmt.Definer
	if definer == "" {
		definer = "`root`@`%`"
	}

	// The body is analyzed for semantic IR (used by lineage, SDL diff) and for
	// the semantic failures MariaDB's ALTER VIEW reports too.
	if err := validateViewPreAnalysis(stmt.Select); err != nil {
		return err
	}
	analyzedQuery, err := c.analyzeViewBodyFor(stmt.Select, source, db, true)
	if err != nil {
		return err
	}
	if err := c.validateViewPostAnalysis(stmt.Select, db); err != nil {
		return err
	}

	// Resolve, rewrite, and deparse the SELECT to produce canonical definition.
	definition, derivedCols := c.deparseViewSelect(stmt.Select, stmt.SelectText, source, db, stmt.Columns)

	// Use explicit column list if provided, otherwise derive from SELECT target list.
	hasExplicit := len(stmt.Columns) > 0
	viewCols := stmt.Columns
	if !hasExplicit {
		viewCols = derivedCols
	}
	if err := checkDuplicateViewColumns(viewCols); err != nil {
		return err
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
		SelectAST:       stmt.Select,
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
	case *nodes.PositionExpr:
		// mariadb:11.8.8: IS_NULLABLE is YES for position('x' in s)
		// over a nullable s and for position('x' in NULL), NO over a NOT NULL
		// column.
		return inferViewExprNullable(e.Substr, relations) ||
			inferViewExprNullable(e.Str, relations)
	case *nodes.ExtractExpr:
		// mariadb:11.8.8: IS_NULLABLE is YES for extract(day from a) over a
		// nullable a.
		return inferViewExprNullable(e.Expr, relations)
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
	for _, item := range sel.GroupBy {
		if _, err := inferViewExprCollation(item.Expr, relations); err != nil {
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
	case "repeat":
		// One string argument, whose charset the result keeps.
		if len(fn.Args) == 0 {
			return viewCollationInfo{}, nil
		}
		return inferViewExprCollation(fn.Args[0], relations)
	case "lpad", "rpad":
		// Two string arguments. MariaDB coerces the result to the wider of
		// them, as it does for CONCAT, where MySQL 8.0 pins it to the first.
		// Measured on mariadb:11.8.8: `LPAD(a, 20, b)` over a latin1 `a` and a
		// utf8mb4 `b` reports utf8mb4, either way round.
		var out viewCollationInfo
		for i, arg := range fn.Args {
			// The middle argument is the length, not a string.
			if i == 1 {
				continue
			}
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
		// A spaced gated keyword is a generic scalar call, not the aggregate, so
		// it does not make the view non-updatable (see analyzeFuncCall).
		if fn, ok := n.(*nodes.FuncCallExpr); ok && isAggregateFunc(fn.Name) && !fn.SpacedParen {
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
			// A name no view holds is 4092 "Unknown VIEW", whether a table
			// holds it or nothing does (mariadb:11.8.8).
			return errUnknownView(db.Name, ref.Name)
		}
		delete(db.Views, key)
	}
	return nil
}

// deparseViewSelect resolves, rewrites, and deparses the SELECT AST for a view.
// If the AST is nil (parser didn't produce one), falls back to the raw SelectText.
// Returns the deparsed definition and the derived column names from the resolved
// SELECT target list (used when no explicit column list is specified).
//
// source is the text the statement was parsed from, which the AST's Loc
// offsets index into. An unaliased select item is named after its text in
// source, as the engine names it; with no source ("" from DefineView), the
// name is rebuilt from the AST.
func (c *Catalog) deparseViewSelect(sel *nodes.SelectStmt, rawText, source string, db *Database, explicitCols []string) (string, []string) {
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
		Lookup:              lookup,
		SchemaLookup:        c.schemaTableLookup(db),
		DefaultDB:           db.Name,
		DefaultCharset:      charset,
		ConnectionCollation: c.session.CollationConnection,
		Source:              source,
	}
	resolver.Resolve(sel)

	// Extract column names from the resolved target list.
	derivedCols := extractViewColumns(sel)

	// Rewrite: NOT folding, boolean context wrapping.
	deparse.RewriteSelectStmt(sel)

	// A view that names its columns carries those names into the body it
	// stores: `CREATE VIEW v (x, y) AS SELECT a, b FROM t` comes back as
	// `select ... AS \`x\`, ... AS \`y\``. Measured on mariadb:11.8.8.
	applyExplicitViewColumns(sel, explicitCols)

	// Deparse: AST → canonical SQL text.
	// MariaDB stores a table or column of the view's own database without the
	// database, and a function call with it. Measured on mariadb:11.8.8: `select
	// test.t.a, test.f(1) from test.t` in database test is stored select
	// `t`.`a` AS `a`,`test`.`f`(1) AS `test.f(1)` from `t`.
	return stripBodyQualifier(deparse.DeparseSelect(sel), db.Name, db.Name), derivedCols
}

// applyExplicitViewColumns renames the output aliases of the SELECT that
// produces the view's rows to the names the view declares. A name past the end
// of the target list is left to the column-count check, which rejects the
// statement.
func applyExplicitViewColumns(sel *nodes.SelectStmt, cols []string) {
	if sel == nil || len(cols) == 0 {
		return
	}
	leaf := nodes.LeftmostQueryLeaf(sel)
	if leaf == nil {
		return
	}
	renamed := make(map[string]string, len(cols))
	for i, target := range leaf.TargetList {
		if i >= len(cols) {
			break
		}
		if rt, ok := target.(*nodes.ResTarget); ok {
			if rt.Name != "" && !strings.EqualFold(rt.Name, cols[i]) {
				renamed[toLower(rt.Name)] = cols[i]
			}
			rt.Name = cols[i]
		}
	}
	// A HAVING reference to an output column is kept by name, so the name the
	// view declares replaces the one the target carried: measured on
	// mariadb:11.8.8, `CREATE VIEW v (x) AS SELECT a AS z FROM t HAVING z > 0`
	// comes back having `x` > 0.
	//
	// An ORDER BY item keeps the name too wherever it is not the bare alias,
	// which is every target but a column: `v (x) AS SELECT a + 1 AS z FROM t
	// ORDER BY z + 0` comes back order by `x` + 0, and so do ABS(z), -z, a
	// constant target and RAND(). A GROUP BY item carries no name to rename,
	// its mention being resolved to the target's expression, and a set
	// operation's own ORDER BY keeps the alias either way — a body the server
	// itself cannot execute.
	if len(renamed) > 0 {
		deparse.RenameOutputRefs(leaf.Having, renamed)
		for _, item := range leaf.OrderBy {
			if item != nil {
				deparse.RenameOutputRefs(item.Expr, renamed)
			}
		}
	}
}

// extractViewColumns extracts column names from a resolved SELECT target list.
// This produces the column list that MySQL would derive for a view.
func extractViewColumns(sel *nodes.SelectStmt) []string {
	if sel == nil {
		return nil
	}
	// The result columns live on the leftmost output SELECT. Unwrap any
	// parenthesized-query wrappers and walk set-operation left arms to reach it,
	// so a parenthesized or set-op view body (CREATE VIEW v AS (SELECT a, b) or
	// (SELECT 1 AS a) UNION (SELECT 2 AS a)) derives columns from that SELECT.
	sel = nodes.LeftmostQueryLeaf(sel)
	var cols []string
	for _, target := range sel.TargetList {
		if rt, ok := target.(*nodes.ResTarget); ok {
			if rt.Name != "" {
				cols = append(cols, rt.Name)
			} else if cr, ok := rt.Val.(*nodes.ColumnRef); ok {
				cols = append(cols, cr.Column)
			} else {
				cols = append(cols, nodeToSQL(rt.Val))
			}
		} else if cr, ok := target.(*nodes.ColumnRef); ok {
			cols = append(cols, cr.Column)
		} else {
			cols = append(cols, nodeToSQL(target))
		}
	}
	return cols
}

// schemaTableLookup returns a schema-aware lookup function that resolves
// schema-qualified references against the referenced database in c.
func (c *Catalog) schemaTableLookup(defaultDB *Database) func(schema, tableName string) *deparse.ResolverTable {
	return func(schema, tableName string) *deparse.ResolverTable {
		targetDB := defaultDB
		if schema != "" {
			targetDB = c.GetDatabase(schema)
		}
		if targetDB == nil {
			return nil
		}
		return tableLookupForDB(targetDB)(tableName)
	}
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
				colType := c.ColumnType
				if colType == "" {
					colType = c.DataType
				}
				cols[i] = deparse.ResolverColumn{
					Name:       c.Name,
					Position:   c.Position,
					ColumnType: colType,
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
		// Fall back to sequences (MariaDB exposes sequences as 8-column tables).
		if seq := db.Sequences[key]; seq != nil {
			cols := make([]deparse.ResolverColumn, len(sequenceColumnNames))
			for i, colName := range sequenceColumnNames {
				cols[i] = deparse.ResolverColumn{
					Name:       colName,
					Position:   i + 1,
					ColumnType: "bigint(21)",
				}
			}
			return &deparse.ResolverTable{
				Name:    seq.Name,
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

// formatDefiner ensures the definer string is backtick-quoted per MySQL/MariaDB
// SHOW CREATE format, preserving quoted account boundaries (e.g. 'a@b'@'localhost'
// becomes `a@b`@`localhost`). Special keywords CURRENT_USER[()] and
// CURRENT_ROLE[()] are kept as written.
func formatDefiner(definer string) string {
	switch strings.ToUpper(definer) {
	case "CURRENT_USER", "CURRENT_USER()", "CURRENT_ROLE", "CURRENT_ROLE()":
		return definer
	}
	userRaw, hostRaw, hasAt := splitDefinerAccount(definer)
	user := unquoteDefinerPart(userRaw)
	if !hasAt {
		return backtickIDQ(user)
	}
	host := unquoteDefinerPart(hostRaw)
	return backtickIDQ(user) + "@" + backtickIDQ(host)
}

// splitDefinerAccount splits a definer string into its user and host parts
// while respecting quoted user names that may contain '@'.
func splitDefinerAccount(s string) (user, host string, hasAt bool) {
	if s == "" {
		return "", "", false
	}
	switch q := s[0]; q {
	case '`':
		if end, ok := scanAccountPart(s, 0); ok {
			if end == len(s) {
				return s, "", false
			}
			if s[end] == '@' {
				return s[:end], s[end+1:], true
			}
		}
	case '\'', '"':
		// parseDefinerValue wraps already-decoded string literal values in
		// single quotes without re-escaping inner quotes or backslashes.
		// Hostnames cannot contain '@', so the last '@' preceded by q marks
		// the user/host boundary even when the user contains '@', '\', or q.
		for i := len(s) - 1; i > 0; i-- {
			if s[i] == '@' && s[i-1] == q {
				return s[:i], s[i+1:], true
			}
		}
		if len(s) >= 2 && s[len(s)-1] == q {
			return s, "", false
		}
	}
	// Hostnames cannot contain '@', so splitting on the last '@' handles both
	// plain user@host and backtick-unquoted usernames containing '@' (e.g. a@b@localhost).
	if at := strings.LastIndexByte(s, '@'); at >= 0 {
		return s[:at], s[at+1:], true
	}
	return s, "", false
}

// unquoteDefinerPart removes matching outer quotes (`, ', ") from a definer
// user or host component. Backtick-quoted components unescape doubled
// backticks; single/double-quoted components come from parseDefinerValue with
// lexer escapes already decoded, so only their outer quote pair is stripped.
func unquoteDefinerPart(s string) string {
	if len(s) < 2 {
		return s
	}
	q := s[0]
	if (q != '`' && q != '\'' && q != '"') || s[len(s)-1] != q {
		return s
	}
	inner := s[1 : len(s)-1]
	if q == '`' {
		return strings.ReplaceAll(inner, "``", "`")
	}
	return inner
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
	b.WriteString(" VIEW " + backtickIDQ(v.Name))

	// No column list: MariaDB carries the names a view declares into its
	// stored body's aliases and writes no list of its own, where MySQL 8.0
	// writes the list back. Measured on mariadb:11.8.8 through SHOW CREATE
	// VIEW; applyExplicitViewColumns puts the names where they belong.

	// AS select_statement
	b.WriteString(" AS ")
	b.WriteString(v.Definition)

	// WITH CHECK OPTION
	if v.CheckOption != "" {
		b.WriteString(fmt.Sprintf(" WITH %s CHECK OPTION", strings.ToUpper(v.CheckOption)))
	}

	return b.String()
}
