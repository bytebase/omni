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
//   - Sequence functions (NEXTVAL / LASTVAL / SETVAL / NEXT VALUE FOR) naming
//     a table or view instead of a sequence (1347)
//   - Schema-qualified function calls (1305: TiDB has no stored functions)
//   - Builtin function argument counts for known native functions (1582)
//   - Window function arguments (NTILE / NTH_VALUE, 1210), duplicate window
//     definitions (3591), undefined window references (3579), and illegal
//     window frame bounds or ROWS INTERVAL frames (3584, 3585, 3596)
//   - Unknown columns in blocks whose FROM scope is fully known (1054)
func (c *Catalog) checkViewBodySemantics(sel *nodes.SelectStmt, db *Database, source string) error {
	if sel == nil {
		return nil
	}
	if err := checkViewRowTargets(sel); err != nil {
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
		if !fc.SpacedParen && !isBacktickQuotedAt(source, fc.Loc.Start) {
			lower := strings.ToLower(fc.Name)
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
// references and known non-positive literals while allowing constant
// expressions so valid folds are not falsely refused.
func isValidWindowPositiveIntArg(expr nodes.ExprNode) bool {
	if expr == nil {
		return false
	}
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
				return false
			}
		}
		break
	}
	switch e := expr.(type) {
	case *nodes.IntLit:
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
	// Check frame bounds in TiDB's order: 3584 (start UNBOUNDED FOLLOWING),
	// 3585 (end UNBOUNDED PRECEDING), 3596 (INTERVAL with ROWS frame).
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
	cols         map[string]bool
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
			var cols map[string]bool
			if tbl := targetDB.GetTable(t.Name); tbl != nil {
				cols = make(map[string]bool, len(tbl.Columns))
				for _, col := range tbl.Columns {
					cols[toLower(col.Name)] = true
				}
			} else if v := targetDB.Views[toLower(t.Name)]; v != nil && len(v.Columns) > 0 {
				cols = make(map[string]bool, len(v.Columns))
				for _, colName := range v.Columns {
					cols[toLower(colName)] = true
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

func (s *viewColScope) hasColumn(cr *nodes.ColumnRef) bool {
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
			if r.cols[col] {
				return true
			}
		}
		return false
	}
	for _, r := range s.relations {
		if r.cols[col] {
			return true
		}
	}
	return false
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
	}
	if sel.Where != nil {
		if err := c.checkExprColumns(sel.Where, db, source, ctes, scopes, "where clause"); err != nil {
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
