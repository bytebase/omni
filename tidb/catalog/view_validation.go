package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// checkViewBodyPreprocess enforces view-body restrictions that TiDB checks
// during preprocessing, before star expansion (1051) and before checking for
// view name collisions (1050):
//
//   - SELECT ... INTO (1350, ER_VIEW_SELECT_CLAUSE)
//   - Noop functions disabled by default: SQL_CALC_FOUND_ROWS and
//     GROUP BY expr ASC|DESC (1235, ER_NOT_SUPPORTED_YET)
//   - User or system variables @v / @@v (1351, ER_VIEW_SELECT_VARIABLE)
//   - Duplicate CTE names within a single WITH clause (1066, ER_NONUNIQ_TABLE)
func checkViewBodyPreprocess(sel *nodes.SelectStmt) error {
	if sel == nil {
		return nil
	}
	var refusal error
	var checkExprs func(rows [][]nodes.ExprNode)
	inspect := func(root nodes.Node) {
		nodes.Inspect(root, func(n nodes.Node) bool {
			if refusal != nil {
				return false
			}
			switch x := n.(type) {
			case *nodes.SelectStmt:
				if x.Into != nil {
					refusal = &Error{
						Code:     ErrViewSelectClause,
						SQLState: sqlState(ErrViewSelectClause),
						Message:  "View's SELECT contains a 'INTO' clause",
					}
					return false
				}
				if x.CalcFoundRows {
					refusal = &Error{
						Code:     ErrNotSupportedYet,
						SQLState: sqlState(ErrNotSupportedYet),
						Message:  "function SQL_CALC_FOUND_ROWS has only noop implementation in tidb now, use tidb_enable_noop_functions to enable these functions",
					}
					return false
				}
				for _, g := range x.GroupBy {
					if g != nil && (g.Desc || g.ExplicitAsc) {
						refusal = &Error{
							Code:     ErrNotSupportedYet,
							SQLState: sqlState(ErrNotSupportedYet),
							Message:  "function GROUP BY expr ASC|DESC has only noop implementation in tidb now, use tidb_enable_noop_functions to enable these functions",
						}
						return false
					}
				}
				if len(x.CTEs) > 1 {
					seen := make(map[string]bool, len(x.CTEs))
					for _, cte := range x.CTEs {
						k := toLower(cte.Name)
						if seen[k] {
							refusal = &Error{
								Code:     ErrNonuniqTable,
								SQLState: sqlState(ErrNonuniqTable),
								Message:  fmt.Sprintf("Not unique table/alias: '%s'", cte.Name),
							}
							return false
						}
						seen[k] = true
					}
				}
			case *nodes.VariableRef:
				refusal = &Error{
					Code:     ErrViewSelectVariable,
					SQLState: sqlState(ErrViewSelectVariable),
					Message:  "View's SELECT contains a variable or parameter",
				}
				return false
			case *nodes.ValuesStmt:
				checkExprs(x.Rows)
			}
			return refusal == nil
		})
	}
	checkExprs = func(rows [][]nodes.ExprNode) {
		for _, row := range rows {
			for _, e := range row {
				if refusal != nil {
					return
				}
				if e != nil {
					inspect(e)
				}
			}
		}
	}
	inspect(sel)
	return refusal
}

// checkViewBodySemantics validates static semantic rules that TiDB v8.5.5
// enforces when planning a CREATE VIEW body:
//
//   - Top-level ROW(a, b, ...) in a SELECT target list (1241)
//   - Unsupported view expressions: MATCH ... AGAINST (1105),
//     CAST(... AS ... ARRAY) (1235), and invalid DATE/TIME/TIMESTAMP literals (1292)
//   - Sequence functions (NEXTVAL / LASTVAL / SETVAL / NEXT VALUE FOR) naming
//     a table or view instead of a sequence (1347)
//   - Schema-qualified and unknown unqualified function calls (1305: TiDB has
//     no stored functions or GIS builtins)
//   - Builtin function argument counts for known native functions (1582, and
//     1105 for APPROX_PERCENTILE)
//   - Window function arguments (NTILE / NTH_VALUE, 1210), duplicate window
//     definitions (3591), undefined window references (3579), GROUPS frames
//     (1235), and illegal window frame bounds or ROWS INTERVAL frames
//     (3584, 3585, 3596)
//   - Unknown columns in blocks whose FROM scope is fully known (1054) and
//     non-JSON/non-string arguments to ->, ->>, and MEMBER OF (3146)
func (c *Catalog) checkViewBodySemantics(sel *nodes.SelectStmt, db *Database, source string) error {
	if sel == nil {
		return nil
	}
	if err := checkViewRowTargets(sel); err != nil {
		return err
	}
	if err := checkViewUnsupportedExprs(sel); err != nil {
		return err
	}
	if err := c.checkViewFunctions(sel, db, source); err != nil {
		return err
	}
	if err := checkViewWindows(sel, source); err != nil {
		return err
	}
	return c.checkViewUnknownColumns(sel, db, source, nil, nil)
}

// checkViewUnsupportedExprs rejects expressions TiDB v8.5.5 refuses during
// view expression building: MATCH ... AGAINST (1105), CAST(... AS ... ARRAY)
// outside functional indexes (1235), and malformed temporal literals (1292).
func checkViewUnsupportedExprs(sel *nodes.SelectStmt) error {
	var refusal error
	nodes.Inspect(sel, func(n nodes.Node) bool {
		if refusal != nil {
			return false
		}
		switch x := n.(type) {
		case *nodes.MatchExpr:
			refusal = &Error{
				Code:     ErrUnknownError,
				SQLState: sqlState(ErrUnknownError),
				Message:  "UnknownType: *ast.MatchAgainst",
			}
			return false
		case *nodes.CastExpr:
			if x.Array {
				refusal = &Error{
					Code:     ErrNotSupportedYet,
					SQLState: sqlState(ErrNotSupportedYet),
					Message:  "This version of TiDB doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'",
				}
				return false
			}
		case *nodes.TemporalLit:
			if x.Expr == nil && !isValidTemporalLiteral(x.Type, x.Value) {
				kind := strings.ToLower(x.Type)
				if kind == "timestamp" {
					kind = "datetime"
				}
				refusal = &Error{
					Code:     ErrTruncatedWrongValue,
					SQLState: sqlState(ErrTruncatedWrongValue),
					Message:  fmt.Sprintf("Incorrect %s value: '%s'", kind, x.Value),
				}
				return false
			}
		}
		return true
	})
	return refusal
}

// checkViewRowTargets rejects a SELECT target list item whose top-level
// expression is a multi-column ROW(...) constructor (1241, Operand should
// contain 1 column(s)).
func checkViewRowTargets(sel *nodes.SelectStmt) error {
	var refusal error
	nodes.Inspect(sel, func(n nodes.Node) bool {
		if refusal != nil {
			return false
		}
		block, ok := n.(*nodes.SelectStmt)
		if !ok {
			return true
		}
		for _, target := range block.TargetList {
			expr := target
			if rt, ok := target.(*nodes.ResTarget); ok {
				expr = rt.Val
			}
			for {
				p, ok := expr.(*nodes.ParenExpr)
				if !ok {
					break
				}
				expr = p.Expr
			}
			if row, ok := expr.(*nodes.RowExpr); ok && len(row.Items) != 1 {
				refusal = &Error{
					Code:     ErrOperandColumns,
					SQLState: sqlState(ErrOperandColumns),
					Message:  "Operand should contain 1 column(s)",
				}
				return false
			}
		}
		return true
	})
	return refusal
}

// isBacktickQuotedAt reports whether source[pos] is a backtick, which marks a
// quoted identifier rather than a builtin keyword.
func isBacktickQuotedAt(source string, pos int) bool {
	return pos >= 0 && pos < len(source) && source[pos] == '`'
}

// isSequenceBuiltinCall reports whether fc is an unquoted, unqualified call to
// NEXTVAL, LASTVAL, or SETVAL (including NEXT VALUE FOR, which parses as NEXTVAL).
func isSequenceBuiltinCall(fc *nodes.FuncCallExpr, source string) bool {
	if fc == nil || fc.Schema != "" || fc.SpacedParen || isBacktickQuotedAt(source, fc.Loc.Start) {
		return false
	}
	switch strings.ToUpper(fc.Name) {
	case "NEXTVAL", "LASTVAL", "SETVAL":
		return len(fc.Args) >= 1
	}
	return false
}

// builtinArity bounds the valid argument count for a native TiDB function.
// max == -1 means variadic (no upper bound).
type builtinArity struct {
	min, max int
}

// tidbBuiltinArity records exact argument count bounds for native TiDB v8.5.5
// functions that reject invalid parameter counts with 1582. Only unqualified,
// unquoted, non-spaced calls are checked against this table.
var tidbBuiltinArity = map[string]builtinArity{
	"ascii":       {1, 1},
	"charset":     {1, 1},
	"coalesce":    {1, -1},
	"collation":   {1, 1},
	"date":        {1, 1},
	"day":         {1, 1},
	"format":      {2, 3},
	"hour":        {1, 1},
	"if":          {3, 3},
	"insert":      {4, 4},
	"interval":    {2, -1},
	"left":        {2, 2},
	"log":         {1, 2},
	"microsecond": {1, 1},
	"minute":      {1, 1},
	"month":       {1, 1},
	"password":    {1, 1},
	"quarter":     {1, 1},
	"repeat":      {2, 2},
	"replace":     {3, 3},
	"reverse":     {1, 1},
	"right":       {2, 2},
	"second":      {1, 1},
	"time":        {1, 1},
	"timestamp":   {1, 2},
	"truncate":    {2, 2},
	"week":        {1, 2},
	"year":        {1, 1},
}

// checkViewFunctions validates function calls in a view body:
//   - Sequence calls naming a table or view instead of a sequence (1347)
//   - Schema-qualified calls (1305: TiDB has no stored functions)
//   - Builtin function parameter counts (1582)
func (c *Catalog) checkViewFunctions(sel *nodes.SelectStmt, db *Database, source string) error {
	var refusal error
	nodes.Inspect(sel, func(n nodes.Node) bool {
		if refusal != nil {
			return false
		}
		fc, ok := n.(*nodes.FuncCallExpr)
		if !ok {
			return true
		}
		if fc.Schema != "" {
			refusal = &Error{
				Code:     ErrNoSuchFunction,
				SQLState: sqlState(ErrNoSuchFunction),
				Message:  fmt.Sprintf("FUNCTION %s.%s does not exist", fc.Schema, fc.Name),
			}
			return false
		}
		if isSequenceBuiltinCall(fc, source) {
			if cr, ok := fc.Args[0].(*nodes.ColumnRef); ok {
				targetDB := db
				if cr.Table != "" {
					targetDB = c.GetDatabase(cr.Table)
				}
				if targetDB != nil {
					key := toLower(cr.Column)
					if targetDB.Sequences[key] == nil && (targetDB.Tables[key] != nil || targetDB.Views[key] != nil) {
						refusal = errWrongObjectKind(targetDB.Name, cr.Column, "SEQUENCE")
						return false
					}
				}
			}
			return true
		}
		lower := strings.ToLower(fc.Name)
		if !tidbBuiltinFunctions[lower] {
			dbName := ""
			if db != nil {
				dbName = db.Name
			}
			refusal = &Error{
				Code:     ErrNoSuchFunction,
				SQLState: sqlState(ErrNoSuchFunction),
				Message:  fmt.Sprintf("FUNCTION %s.%s does not exist", dbName, lower),
			}
			return false
		}
		if !fc.SpacedParen && !isBacktickQuotedAt(source, fc.Loc.Start) {
			if lower == "approx_percentile" && len(fc.Args) != 2 {
				refusal = &Error{
					Code:     ErrUnknownError,
					SQLState: sqlState(ErrUnknownError),
					Message:  "APPROX_PERCENTILE should take 2 arguments",
				}
				return false
			}
			if arity, ok := tidbBuiltinArity[lower]; ok {
				nArgs := len(fc.Args)
				if nArgs < arity.min || (arity.max >= 0 && nArgs > arity.max) {
					refusal = &Error{
						Code:     ErrWrongParamCount,
						SQLState: sqlState(ErrWrongParamCount),
						Message:  fmt.Sprintf("Incorrect parameter count in the call to native function '%s'", lower),
					}
					return false
				}
			}
		}
		return true
	})
	return refusal
}

// checkViewWindows validates window definitions, window function arguments,
// window references, and window frames for every SELECT block in sel.
func checkViewWindows(sel *nodes.SelectStmt, source string) error {
	var refusal error
	nodes.Inspect(sel, func(n nodes.Node) bool {
		if refusal != nil {
			return false
		}
		block, ok := n.(*nodes.SelectStmt)
		if !ok {
			return true
		}
		refusal = checkSelectBlockWindows(block, source)
		return refusal == nil
	})
	return refusal
}

func checkSelectBlockWindows(block *nodes.SelectStmt, source string) error {
	// 1. Check duplicate names in WINDOW clause (3591) and validate each window's frame.
	defined := make(map[string]bool, len(block.WindowClause))
	for _, wd := range block.WindowClause {
		if wd == nil {
			continue
		}
		if wd.Name != "" {
			key := toLower(wd.Name)
			if defined[key] {
				return &Error{
					Code:     ErrWindowDuplicateName,
					SQLState: sqlState(ErrWindowDuplicateName),
					Message:  fmt.Sprintf("Window '%s' is defined twice.", wd.Name),
				}
			}
			defined[key] = true
		}
	}
	for _, wd := range block.WindowClause {
		if err := checkWindowDef(wd, defined); err != nil {
			return err
		}
	}

	// 2. Check window function calls directly belonging to this block (do not
	// descend into nested SelectStmts; they are checked in their own block).
	var checkExpr func(nodes.Node) error
	checkExpr = func(root nodes.Node) error {
		if root == nil {
			return nil
		}
		var err error
		nodes.Inspect(root, func(n nodes.Node) bool {
			if err != nil {
				return false
			}
			if _, isSel := n.(*nodes.SelectStmt); isSel && n != root {
				return false
			}
			fc, ok := n.(*nodes.FuncCallExpr)
			if !ok || fc.Over == nil {
				return true
			}
			// TiDB validates NTILE / NTH_VALUE arguments (1210) before looking up
			// the referenced window name (3579).
			if fc.Schema == "" && !fc.SpacedParen && !isBacktickQuotedAt(source, fc.Loc.Start) {
				switch strings.ToUpper(fc.Name) {
				case "NTILE":
					if len(fc.Args) == 1 && !isValidWindowPositiveIntArg(fc.Args[0]) {
						err = &Error{
							Code:     ErrWrongArguments,
							SQLState: sqlState(ErrWrongArguments),
							Message:  "Incorrect arguments to ntile",
						}
						return false
					}
				case "NTH_VALUE":
					if len(fc.Args) == 2 && !isValidWindowPositiveIntArg(fc.Args[1]) {
						err = &Error{
							Code:     ErrWrongArguments,
							SQLState: sqlState(ErrWrongArguments),
							Message:  "Incorrect arguments to nth_value",
						}
						return false
					}
				}
			}
			if werr := checkWindowDef(fc.Over, defined); werr != nil {
				err = werr
				return false
			}
			return true
		})
		return err
	}

	for _, target := range block.TargetList {
		if err := checkExpr(target); err != nil {
			return err
		}
	}
	for _, item := range block.OrderBy {
		if item != nil {
			if err := checkExpr(item.Expr); err != nil {
				return err
			}
		}
	}
	return nil
}

// isValidWindowPositiveIntArg reports whether expr is acceptable as the
// positive integer argument of NTILE(n) or NTH_VALUE(x, n). It rejects column
// references, NULL, and signed integer literals folding to <= 0, while
// tracking unary-minus sign flips (e.g. -(-2) is +2) and allowing non-literal
// constant expressions so valid constant folds are not falsely refused.
func isValidWindowPositiveIntArg(expr nodes.ExprNode) bool {
	if expr == nil {
		return false
	}
	negated := false
	for {
		switch e := expr.(type) {
		case *nodes.ParenExpr:
			expr = e.Expr
			continue
		case *nodes.UnaryExpr:
			if e.Op == nodes.UnaryPlus {
				expr = e.Operand
				continue
			}
			if e.Op == nodes.UnaryMinus {
				negated = !negated
				expr = e.Operand
				continue
			}
		}
		break
	}
	switch e := expr.(type) {
	case *nodes.IntLit:
		if negated {
			return e.Value < 0
		}
		return e.Value > 0
	case *nodes.ColumnRef, *nodes.NullLit:
		return false
	}
	hasCol := false
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if _, ok := n.(*nodes.ColumnRef); ok {
			hasCol = true
			return false
		}
		return !hasCol
	})
	return !hasCol
}

func checkWindowDef(wd *nodes.WindowDef, defined map[string]bool) error {
	if wd == nil {
		return nil
	}
	if wd.RefName != "" && !defined[toLower(wd.RefName)] {
		return &Error{
			Code:     ErrWindowNoSuchWindow,
			SQLState: sqlState(ErrWindowNoSuchWindow),
			Message:  fmt.Sprintf("Window name '%s' is not defined.", wd.RefName),
		}
	}
	f := wd.Frame
	if f == nil {
		return nil
	}
	winName := wd.Name
	if winName == "" {
		winName = "<unnamed window>"
	}
	// Check frame type and bounds in TiDB's order: GROUPS unsupported (1235),
	// 3584 (start UNBOUNDED FOLLOWING), 3585 (end UNBOUNDED PRECEDING),
	// 3596 (INTERVAL with ROWS frame).
	if f.Type == nodes.FrameGroups {
		return &Error{
			Code:     ErrNotSupportedYet,
			SQLState: sqlState(ErrNotSupportedYet),
			Message:  "This version of TiDB doesn't yet support 'GROUPS'",
		}
	}
	if f.Start != nil && f.Start.Type == nodes.BoundUnboundedFollowing {
		return &Error{
			Code:     ErrWindowFrameStartIllegal,
			SQLState: sqlState(ErrWindowFrameStartIllegal),
			Message:  fmt.Sprintf("Window '%s': frame start cannot be UNBOUNDED FOLLOWING.", winName),
		}
	}
	if f.End != nil && f.End.Type == nodes.BoundUnboundedPreceding {
		return &Error{
			Code:     ErrWindowFrameEndIllegal,
			SQLState: sqlState(ErrWindowFrameEndIllegal),
			Message:  fmt.Sprintf("Window '%s': frame end cannot be UNBOUNDED PRECEDING.", winName),
		}
	}
	if f.Type == nodes.FrameRows && (isIntervalBound(f.Start) || isIntervalBound(f.End)) {
		return &Error{
			Code:     ErrWindowRowsIntervalUse,
			SQLState: sqlState(ErrWindowRowsIntervalUse),
			Message:  fmt.Sprintf("Window '%s': INTERVAL can only be used with RANGE frames.", winName),
		}
	}
	return nil
}

func isIntervalBound(b *nodes.WindowFrameBound) bool {
	if b == nil || b.Offset == nil {
		return false
	}
	_, ok := b.Offset.(*nodes.IntervalExpr)
	return ok
}

// viewColScope holds the known columns of every FROM item in one SELECT block.
type viewColScope struct {
	relations []viewColRelation
	opaque    bool
}

type viewColRelation struct {
	schema, name string
	cols         map[string]string // lowercase col name -> lowercase DataType ("" if unknown)
}

func (c *Catalog) buildViewColScope(sel *nodes.SelectStmt, db *Database, source string, ctes map[string]bool) *viewColScope {
	s := &viewColScope{}
	if readsOnlyNoTable(sel, source) {
		return s
	}
	var add func(nodes.TableExpr)
	add = func(te nodes.TableExpr) {
		switch t := te.(type) {
		case *nodes.TableRef:
			if t.Schema == "" && ctes[toLower(t.Name)] {
				s.opaque = true
				return
			}
			targetDB := db
			if t.Schema != "" {
				targetDB = c.GetDatabase(t.Schema)
			}
			if targetDB == nil {
				s.opaque = true
				return
			}
			var cols map[string]string
			if tbl := targetDB.GetTable(t.Name); tbl != nil {
				cols = make(map[string]string, len(tbl.Columns)+1)
				for _, col := range tbl.Columns {
					cols[toLower(col.Name)] = strings.ToLower(col.DataType)
				}
				if tableHasTiDBRowID(tbl) {
					cols["_tidb_rowid"] = "bigint"
				}
			} else if v := targetDB.Views[toLower(t.Name)]; v != nil && len(v.Columns) > 0 {
				cols = make(map[string]string, len(v.Columns))
				for _, colName := range v.Columns {
					cols[toLower(colName)] = ""
				}
			} else {
				s.opaque = true
				return
			}
			rel := viewColRelation{name: toLower(t.Name), cols: cols}
			switch {
			case t.Alias != "":
				rel.name = toLower(t.Alias)
			case t.Schema != "":
				rel.schema = toLower(t.Schema)
			default:
				rel.schema = toLower(targetDB.Name)
			}
			s.relations = append(s.relations, rel)
		case *nodes.JoinClause:
			add(t.Left)
			add(t.Right)
		default:
			s.opaque = true
		}
	}
	for _, from := range sel.From {
		add(from)
	}
	return s
}

func (s *viewColScope) lookupColumn(cr *nodes.ColumnRef) (string, bool) {
	col := toLower(cr.Column)
	if cr.Table != "" {
		tbl := toLower(cr.Table)
		sch := toLower(cr.Schema)
		for _, r := range s.relations {
			if r.name != tbl {
				continue
			}
			if sch != "" && r.schema != sch {
				continue
			}
			if dt, ok := r.cols[col]; ok {
				return dt, true
			}
		}
		return "", false
	}
	for _, r := range s.relations {
		if dt, ok := r.cols[col]; ok {
			return dt, true
		}
	}
	return "", false
}

func (s *viewColScope) hasColumn(cr *nodes.ColumnRef) bool {
	_, ok := s.lookupColumn(cr)
	return ok
}

// checkViewUnknownColumns rejects column references in SELECT target lists,
// WHERE clauses, and JOIN ON conditions that cannot resolve against any table
// in the block (or its enclosing query scopes) when all visible FROM items
// have known column lists.
func (c *Catalog) checkViewUnknownColumns(sel *nodes.SelectStmt, db *Database, source string, ctes map[string]bool, outer []*viewColScope) error {
	if sel == nil {
		return nil
	}
	ctes = withCTENames(sel, ctes)
	for _, cte := range sel.CTEs {
		if cte != nil && cte.Select != nil {
			// CTEs cannot correlate with outer query blocks.
			if err := c.checkViewUnknownColumns(cte.Select, db, source, ctes, nil); err != nil {
				return err
			}
		}
	}
	if sel.ParenSource != nil {
		return c.checkViewUnknownColumns(sel.ParenSource, db, source, ctes, outer)
	}
	if sel.Left != nil || sel.Right != nil {
		if err := c.checkViewUnknownColumns(sel.Left, db, source, ctes, outer); err != nil {
			return err
		}
		return c.checkViewUnknownColumns(sel.Right, db, source, ctes, outer)
	}
	if sel.TableSource != nil || sel.ValuesSource != nil {
		return nil
	}

	sc := c.buildViewColScope(sel, db, source, ctes)
	scopes := make([]*viewColScope, 0, len(outer)+1)
	scopes = append(scopes, sc)
	scopes = append(scopes, outer...)

	// Check derived tables in FROM (they cannot correlate unless LATERAL, so
	// keep outer nil unless needed, and check JOIN ON against this block's scope).
	var checkFrom func(nodes.TableExpr) error
	checkFrom = func(te nodes.TableExpr) error {
		switch t := te.(type) {
		case *nodes.JoinClause:
			if err := checkFrom(t.Left); err != nil {
				return err
			}
			if err := checkFrom(t.Right); err != nil {
				return err
			}
			if on, ok := t.Condition.(*nodes.OnCondition); ok && on.Expr != nil {
				if err := c.checkExprColumns(on.Expr, db, source, ctes, scopes, "on clause"); err != nil {
					return err
				}
				if err := checkExprJSONTypes(on.Expr, scopes); err != nil {
					return err
				}
			}
		case *nodes.SubqueryExpr:
			if t.Select != nil {
				// A derived table may be LATERAL; pass outer scopes conservatively
				// only when Lateral is true.
				var derivedOuter []*viewColScope
				if t.Lateral {
					derivedOuter = scopes
				}
				if err := c.checkViewUnknownColumns(t.Select, db, source, ctes, derivedOuter); err != nil {
					return err
				}
			}
		}
		return nil
	}
	for _, from := range sel.From {
		if err := checkFrom(from); err != nil {
			return err
		}
	}

	for _, target := range sel.TargetList {
		if err := c.checkExprColumns(target, db, source, ctes, scopes, "field list"); err != nil {
			return err
		}
		if err := checkExprJSONTypes(target, scopes); err != nil {
			return err
		}
	}
	if sel.Where != nil {
		if err := c.checkExprColumns(sel.Where, db, source, ctes, scopes, "where clause"); err != nil {
			return err
		}
		if err := checkExprJSONTypes(sel.Where, scopes); err != nil {
			return err
		}
	}
	return nil
}

func (c *Catalog) checkExprColumns(root nodes.Node, db *Database, source string, ctes map[string]bool, scopes []*viewColScope, clause string) error {
	if root == nil {
		return nil
	}
	allKnown := true
	for _, s := range scopes {
		if s.opaque {
			allKnown = false
			break
		}
	}
	var refusal error
	nodes.Inspect(root, func(n nodes.Node) bool {
		if refusal != nil {
			return false
		}
		switch x := n.(type) {
		case *nodes.SelectStmt:
			refusal = c.checkViewUnknownColumns(x, db, source, ctes, scopes)
			return false
		case *nodes.FuncCallExpr:
			// Sequence builtins take a table/sequence name as Args[0], parsed
			// as a ColumnRef; do not treat it as a column reference.
			if isSequenceBuiltinCall(x, source) {
				for _, arg := range x.Args[1:] {
					if err := c.checkExprColumns(arg, db, source, ctes, scopes, clause); err != nil {
						refusal = err
						return false
					}
				}
				return false
			}
		case *nodes.ColumnRef:
			if x.Star || !allKnown {
				return false
			}
			for _, s := range scopes {
				if s.hasColumn(x) {
					return false
				}
			}
			refusal = errNoSuchColumn(formatColumnRefName(x), clause)
			return false
		}
		return true
	})
	return refusal
}

// tableHasTiDBRowID reports whether a base table exposes TiDB's hidden
// `_tidb_rowid` handle column: true when the table has no primary key or a
// NONCLUSTERED primary key, and false when it has a CLUSTERED primary key.
func tableHasTiDBRowID(tbl *Table) bool {
	for _, idx := range tbl.Indexes {
		if idx.Primary {
			return !primaryKeyClustered(tbl, idx)
		}
	}
	return true
}

func formatColumnRefName(cr *nodes.ColumnRef) string {
	switch {
	case cr.Schema != "":
		return cr.Schema + "." + cr.Table + "." + cr.Column
	case cr.Table != "":
		return cr.Table + "." + cr.Column
	default:
		return cr.Column
	}
}

// checkExprJSONTypes enforces TiDB v8.5.5's 3146 check on ->, ->>, and
// MEMBER OF operands whose static type is known to be neither JSON nor a
// character string type. Nested SELECT subqueries are skipped here because
// checkViewUnknownColumns validates each block in its own scope.
func checkExprJSONTypes(root nodes.Node, scopes []*viewColScope) error {
	if root == nil {
		return nil
	}
	var refusal error
	nodes.Inspect(root, func(n nodes.Node) bool {
		if refusal != nil {
			return false
		}
		if _, isSel := n.(*nodes.SelectStmt); isSel {
			return false
		}
		switch x := n.(type) {
		case *nodes.BinaryExpr:
			if x.Op == nodes.BinOpJsonExtract || x.Op == nodes.BinOpJsonUnquote {
				if !exprCanBeJSON(x.Left, scopes) {
					refusal = &Error{
						Code:     ErrInvalidTypeForJSON,
						SQLState: sqlState(ErrInvalidTypeForJSON),
						Message:  "Invalid data type for JSON data in argument 1 to function json_extract; a JSON string or JSON type is required.",
					}
					return false
				}
			}
		case *nodes.MemberOfExpr:
			if !exprCanBeJSON(x.Array, scopes) {
				refusal = &Error{
					Code:     ErrInvalidTypeForJSON,
					SQLState: sqlState(ErrInvalidTypeForJSON),
					Message:  "Invalid data type for JSON data in argument 2 to function member of; a JSON string or JSON type is required.",
				}
				return false
			}
		}
		return true
	})
	return refusal
}

// exprCanBeJSON reports whether expr could hold a JSON value or JSON string
// (or NULL, or an expression whose type is not statically known).
func exprCanBeJSON(expr nodes.ExprNode, scopes []*viewColScope) bool {
	for {
		p, ok := expr.(*nodes.ParenExpr)
		if !ok {
			break
		}
		expr = p.Expr
	}
	switch e := expr.(type) {
	case *nodes.IntLit, *nodes.FloatLit, *nodes.BoolLit, *nodes.TemporalLit:
		return false
	case *nodes.ColumnRef:
		for _, s := range scopes {
			if dt, ok := s.lookupColumn(e); ok {
				if dt == "" {
					return true
				}
				return isJSONOrStringDataType(dt)
			}
		}
		return true
	default:
		return true
	}
}

func isJSONOrStringDataType(dt string) bool {
	switch strings.ToLower(dt) {
	case "json", "char", "varchar", "tinytext", "text", "mediumtext", "longtext", "enum", "set":
		return true
	}
	return false
}

// isValidTemporalLiteral checks whether val is a valid TiDB DATE, TIME, or
// TIMESTAMP literal string.
func isValidTemporalLiteral(typ, val string) bool {
	s := strings.TrimSpace(val)
	if s == "" {
		return false
	}
	switch strings.ToUpper(typ) {
	case "DATE":
		return isValidDateOrDateTimeLiteral(s, false)
	case "TIMESTAMP":
		return isValidDateOrDateTimeLiteral(s, true)
	case "TIME":
		return isValidTimeLiteral(s)
	default:
		return true
	}
}

func isASCIIDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isASCIIAlnum(b byte) bool {
	return isASCIIDigit(b) || (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func isValidCalendarDate(year, month, day int) bool {
	if year < 0 || year > 9999 || month < 0 || month > 12 || day < 0 || day > 31 {
		return false
	}
	if month == 0 || day == 0 {
		return true
	}
	maxDay := 31
	switch month {
	case 4, 6, 9, 11:
		maxDay = 30
	case 2:
		if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
			maxDay = 29
		} else {
			maxDay = 28
		}
	}
	return day <= maxDay
}

func parseDigits(s string, pos, maxLen int) (val, digits, next int) {
	i := pos
	for i < len(s) && digits < maxLen && isASCIIDigit(s[i]) {
		val = val*10 + int(s[i]-'0')
		digits++
		i++
	}
	return val, digits, i
}

func isValidDateOrDateTimeLiteral(s string, allowTZ bool) bool {
	// Strip optional trailing timezone offset (Z or ±HH:MM) for TIMESTAMP.
	if allowTZ {
		if strings.HasSuffix(s, "Z") || strings.HasSuffix(s, "z") {
			s = s[:len(s)-1]
		} else if len(s) >= 6 {
			sign := s[len(s)-6]
			if (sign == '+' || sign == '-') && s[len(s)-3] == ':' &&
				isASCIIDigit(s[len(s)-5]) && isASCIIDigit(s[len(s)-4]) &&
				isASCIIDigit(s[len(s)-2]) && isASCIIDigit(s[len(s)-1]) {
				tzH := int(s[len(s)-5]-'0')*10 + int(s[len(s)-4]-'0')
				tzM := int(s[len(s)-2]-'0')*10 + int(s[len(s)-1]-'0')
				if tzH > 14 || tzM > 59 {
					return false
				}
				s = s[:len(s)-6]
			}
		}
	}

	// Check for packed numeric form: YYMMDD, YYYYMMDD, YYMMDDHHMMSS, YYYYMMDDHHMMSS[.frac].
	mainPart, fracPart, hasDot := strings.Cut(s, ".")
	allDigits := len(mainPart) > 0
	for i := 0; i < len(mainPart); i++ {
		if !isASCIIDigit(mainPart[i]) {
			allDigits = false
			break
		}
	}
	if allDigits {
		if hasDot {
			for i := 0; i < len(fracPart); i++ {
				if !isASCIIDigit(fracPart[i]) {
					return false
				}
			}
		}
		var y, m, d, hh, mm, ss int
		switch len(mainPart) {
		case 6:
			y, _, _ = parseDigits(mainPart, 0, 2)
			if y < 70 {
				y += 2000
			} else {
				y += 1900
			}
			m, _, _ = parseDigits(mainPart, 2, 2)
			d, _, _ = parseDigits(mainPart, 4, 2)
		case 8:
			y, _, _ = parseDigits(mainPart, 0, 4)
			m, _, _ = parseDigits(mainPart, 4, 2)
			d, _, _ = parseDigits(mainPart, 6, 2)
		case 12:
			y, _, _ = parseDigits(mainPart, 0, 2)
			if y < 70 {
				y += 2000
			} else {
				y += 1900
			}
			m, _, _ = parseDigits(mainPart, 2, 2)
			d, _, _ = parseDigits(mainPart, 4, 2)
			hh, _, _ = parseDigits(mainPart, 6, 2)
			mm, _, _ = parseDigits(mainPart, 8, 2)
			ss, _, _ = parseDigits(mainPart, 10, 2)
		case 14:
			y, _, _ = parseDigits(mainPart, 0, 4)
			m, _, _ = parseDigits(mainPart, 4, 2)
			d, _, _ = parseDigits(mainPart, 6, 2)
			hh, _, _ = parseDigits(mainPart, 8, 2)
			mm, _, _ = parseDigits(mainPart, 10, 2)
			ss, _, _ = parseDigits(mainPart, 12, 2)
		default:
			return false
		}
		return isValidCalendarDate(y, m, d) && hh <= 23 && mm <= 59 && ss <= 59
	}

	// Delimited form: YYYY<sep>MM<sep>DD[<sep|T>HH<sep>MM[<sep>SS][.frac]]
	y, yn, pos := parseDigits(s, 0, 4)
	if yn == 0 || pos >= len(s) || isASCIIAlnum(s[pos]) {
		return false
	}
	if yn <= 2 {
		if y < 70 {
			y += 2000
		} else {
			y += 1900
		}
	}
	for pos < len(s) && !isASCIIAlnum(s[pos]) {
		pos++
	}
	m, mn, pos := parseDigits(s, pos, 2)
	if mn == 0 || pos >= len(s) || isASCIIAlnum(s[pos]) {
		return false
	}
	for pos < len(s) && !isASCIIAlnum(s[pos]) {
		pos++
	}
	d, dn, pos := parseDigits(s, pos, 2)
	if dn == 0 || !isValidCalendarDate(y, m, d) {
		return false
	}
	if pos == len(s) {
		return true
	}
	// Separator between date and time: 'T', 't', or non-alphanumeric punctuation/space.
	if s[pos] == 'T' || s[pos] == 't' {
		pos++
	} else if !isASCIIAlnum(s[pos]) {
		for pos < len(s) && !isASCIIAlnum(s[pos]) {
			pos++
		}
	} else {
		return false
	}
	hh, hn, pos := parseDigits(s, pos, 2)
	if hn == 0 || hh > 23 {
		return false
	}
	if pos < len(s) && s[pos] != '.' {
		if isASCIIAlnum(s[pos]) {
			return false
		}
		for pos < len(s) && !isASCIIAlnum(s[pos]) {
			pos++
		}
		mm, mmn, nextPos := parseDigits(s, pos, 2)
		if mmn == 0 || mm > 59 {
			return false
		}
		pos = nextPos
		if pos < len(s) && s[pos] != '.' {
			if isASCIIAlnum(s[pos]) {
				return false
			}
			for pos < len(s) && !isASCIIAlnum(s[pos]) {
				pos++
			}
			ss, ssn, nextPos2 := parseDigits(s, pos, 2)
			if ssn == 0 || ss > 59 {
				return false
			}
			pos = nextPos2
		}
	}
	if pos < len(s) {
		if s[pos] != '.' {
			return false
		}
		pos++
		for pos < len(s) {
			if !isASCIIDigit(s[pos]) {
				return false
			}
			pos++
		}
	}
	return true
}

func isValidTimeLiteral(s string) bool {
	if strings.HasPrefix(s, "-") {
		s = s[1:]
		if s == "" {
			return false
		}
	}
	mainPart, fracPart, hasDot := strings.Cut(s, ".")
	if hasDot {
		for i := 0; i < len(fracPart); i++ {
			if !isASCIIDigit(fracPart[i]) {
				return false
			}
		}
	}
	if mainPart == "" {
		return false
	}
	days := 0
	if spIdx := strings.IndexByte(mainPart, ' '); spIdx >= 0 {
		dStr := mainPart[:spIdx]
		if dStr == "" || len(dStr) > 2 {
			return false
		}
		for i := 0; i < len(dStr); i++ {
			if !isASCIIDigit(dStr[i]) {
				return false
			}
			days = days*10 + int(dStr[i]-'0')
		}
		mainPart = strings.TrimLeft(mainPart[spIdx+1:], " \t")
		if mainPart == "" {
			return false
		}
	}
	if strings.Contains(mainPart, ":") {
		parts := strings.Split(mainPart, ":")
		if len(parts) < 2 || len(parts) > 3 {
			return false
		}
		nums := make([]int, len(parts))
		for i, p := range parts {
			if p == "" || (i > 0 && len(p) > 2) || (i == 0 && len(p) > 3) {
				return false
			}
			v := 0
			for j := 0; j < len(p); j++ {
				if !isASCIIDigit(p[j]) {
					return false
				}
				v = v*10 + int(p[j]-'0')
			}
			nums[i] = v
		}
		hours := days*24 + nums[0]
		mins := nums[1]
		secs := 0
		if len(nums) == 3 {
			secs = nums[2]
		}
		return hours <= 838 && mins <= 59 && secs <= 59
	}
	// Packed digits form: [[[H]HH]MM]SS (or D HH when days was specified).
	for i := 0; i < len(mainPart); i++ {
		if !isASCIIDigit(mainPart[i]) {
			return false
		}
	}
	if days > 0 {
		if len(mainPart) > 2 {
			return false
		}
		hh, _, _ := parseDigits(mainPart, 0, len(mainPart))
		return days*24+hh <= 838
	}
	if len(mainPart) > 7 {
		return false
	}
	v, _, _ := parseDigits(mainPart, 0, len(mainPart))
	secs := v % 100
	mins := (v / 100) % 100
	hours := v / 10000
	return hours <= 838 && mins <= 59 && secs <= 59
}
