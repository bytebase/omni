package catalog

import (
	"errors"
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/deparse"
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
	// path: mysql:8.0.46 answers 1054 for a SELECT naming a column no relation
	// of the body has, even when the name already belongs to a table (1050) or
	// a view. A synced view's body is the engine's already, so it is not judged
	// here.
	analyzedQuery, bodyErr := c.analyzeViewBodyFor(stmt.Select, source, db, stmt.Name.Name)
	if checkBody && bodyErr != nil {
		return bodyErr
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

	// A `?` parameter marker is not allowed in a view body; MySQL rejects
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
	definition, derivedCols := c.deparseViewSelect(stmt.Select, stmt.SelectText, source, db)

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
		AnalyzedQuery:   analyzedQuery,
		IsUpdatable:     inferViewUpdatable(stmt.Select, algorithm),
	}
	return nil
}

// analyzeViewBody analyzes a view body for the semantic IR lineage and SDL
// diff use, and reports the semantic failures MySQL's CREATE VIEW reports too.
// Failures in alias-bearing clauses (GROUP BY, HAVING, ORDER BY) and
// unsupported analyzer expressions leave the view without AnalyzedQuery rather
// than refusing a view the engine accepts.
func (c *Catalog) analyzeViewBody(sel *nodes.SelectStmt, source string) (*Query, error) {
	return c.analyzeViewBodyFor(sel, source, nil, "")
}

func (c *Catalog) analyzeViewBodyFor(sel *nodes.SelectStmt, source string, viewDB *Database, viewName string) (*Query, error) {
	if sel == nil {
		return nil, nil
	}
	body := collectViewBody(sel, source)
	blocks := body.reachableBlocks()
	if err := c.validateViewBodyPreAnalysis(blocks, viewDB); err != nil {
		return nil, err
	}
	for _, block := range blocks {
		if err := viewBlockStar(block, source); err != nil {
			return nil, err
		}
	}
	savedValidating := c.validatingViewBody
	c.validatingViewBody = true
	q, err := c.analyzeSelectStmtIn(sel, nil, nil, body.reads, source)
	c.validatingViewBody = savedValidating
	if err != nil {
		var refusal *Error
		if errors.As(err, &refusal) {
			if refusal.Code == ErrNoSuchTable && !isSystemDatabaseTableErr(refusal) {
				return nil, refusal
			}
			if refusal.Code == ErrNoSuchColumn && !isAliasBearingColumnErr(refusal) && !hasSchemaQualifiedFuncCall(sel) {
				return nil, refusal
			}
			switch refusal.Code {
			case ErrNoDefaultForField, ErrIllegalValueForType, ErrInvalidTypeForJSON, ErrInvalidCastToTimestampWithTimeZone:
				return nil, refusal
			}
		}
	}
	if postErr := validateViewBodyPostAnalysis(sel, blocks, viewName); postErr != nil {
		return nil, postErr
	}
	if err != nil {
		return nil, nil
	}
	return q, nil
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
	for _, sysDB := range []string{"'information_schema.", "'performance_schema.", "'mysql.", "'sys."} {
		if strings.Contains(strings.ToLower(ce.Message), sysDB) {
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

// collectViewBody walks sel and computes the CTEs reachable from the view body.
func collectViewBody(sel *nodes.SelectStmt, source string) *viewBody {
	body := &viewBody{edges: map[*nodes.CommonTableExpr]map[*nodes.CommonTableExpr]bool{}}
	body.walk(sel, nil, nil, source)
	body.reach()
	return body
}

// reachableBlocks returns the query blocks MySQL resolves for this view body:
// the body's own blocks plus any CTEs reached from it.
func (v *viewBody) reachableBlocks() []*nodes.SelectStmt {
	out := make([]*nodes.SelectStmt, 0, len(v.blocks))
	for _, b := range v.blocks {
		if b.owner != nil && !v.reads[b.owner] {
			continue
		}
		out = append(out, b.block)
	}
	return out
}

// inspectBlockOwnNodes visits nodes belonging directly to block without
// descending into nested SelectStmt query blocks or CTE definitions.
// Also walks ValuesSource.Rows, which nodes.Inspect skips.
func inspectBlockOwnNodes(block *nodes.SelectStmt, fn func(nodes.Node) bool) {
	walk := func(root nodes.Node) {
		nodes.Inspect(root, func(n nodes.Node) bool {
			switch e := n.(type) {
			case *nodes.SelectStmt:
				if e != block {
					return false
				}
			case *nodes.CommonTableExpr:
				return false
			}
			return fn(n)
		})
	}
	walk(block)
	if block.ValuesSource != nil {
		for _, row := range block.ValuesSource.Rows {
			for _, item := range row {
				if item != nil {
					walk(item)
				}
			}
		}
	}
}

// validateViewBodyPreAnalysis checks view body restrictions that mysql:8.0.46
// enforces ahead of column/table resolution (INTO clause, VALUES row errors,
// LIMIT/NTILE variables, user variables, temporal literals, unknown functions,
// locking-clause table names).
func (c *Catalog) validateViewBodyPreAnalysis(blocks []*nodes.SelectStmt, viewDB *Database) error {
	currentDB := c.CurrentDatabase()
	for _, block := range blocks {
		if block.Into != nil {
			return errViewSelectClause("INTO")
		}
		if block.ValuesSource != nil {
			for _, row := range block.ValuesSource.Rows {
				if len(row) == 0 {
					return errEmptyRowInTVC()
				}
				for _, item := range row {
					if d, ok := item.(*nodes.DefaultExpr); ok && d.Column == "" {
						return errDefaultAsValInTVC()
					}
				}
			}
		}
		if block.Limit != nil {
			for _, op := range []nodes.ExprNode{block.Limit.Count, block.Limit.Offset} {
				switch e := op.(type) {
				case *nodes.ColumnRef:
					return errSPUndeclaredVar(e.Column)
				case *nodes.VariableRef:
					if !e.System {
						return errViewSelectVariable()
					}
				}
			}
		}
		var blockErr error
		inspectBlockOwnNodes(block, func(n nodes.Node) bool {
			if blockErr != nil {
				return false
			}
			switch e := n.(type) {
			case *nodes.VariableRef:
				if !e.System {
					blockErr = errViewSelectVariable()
					return false
				}
			case *nodes.TemporalLit:
				if err := validateTemporalLiteral(e); err != nil {
					blockErr = err
					return false
				}
			case *nodes.FuncCallExpr:
				if err := c.validateViewFuncExists(e, viewDB); err != nil {
					blockErr = err
					return false
				}
				if strings.EqualFold(e.Name, "ntile") && !e.SpacedParen && e.Schema == "" && len(e.Args) >= 1 {
					switch arg := e.Args[0].(type) {
					case *nodes.ColumnRef:
						blockErr = errSPUndeclaredVar(arg.Column)
						return false
					case *nodes.VariableRef:
						if !arg.System {
							blockErr = errViewSelectVariable()
							return false
						}
					}
				}
			}
			return true
		})
		if blockErr != nil {
			return blockErr
		}
		if err := validateBlockLocking(block, currentDB); err != nil {
			return err
		}
	}
	return nil
}

func (c *Catalog) validateViewFuncExists(fc *nodes.FuncCallExpr, viewDB *Database) error {
	if fc == nil || fc.Name == "" || c.loadingSDL {
		return nil
	}
	lower := toLower(fc.Name)
	if fc.Schema != "" {
		if isSystemDatabaseName(fc.Schema) {
			return nil
		}
		targetDB := c.GetDatabase(fc.Schema)
		if targetDB == nil || targetDB.Functions[lower] == nil {
			return errNoSuchFunction(fc.Schema + "." + fc.Name)
		}
		return nil
	}
	if fc.SpacedParen {
		return nil
	}
	if isBuiltinMySQLFunction(lower) {
		return nil
	}
	curDB := c.GetDatabase(c.CurrentDatabase())
	if curDB == nil {
		curDB = viewDB
	}
	if curDB != nil && curDB.Functions[lower] != nil {
		return nil
	}
	qualName := fc.Name
	if curDB != nil && curDB.Name != "" {
		qualName = curDB.Name + "." + fc.Name
	}
	return errNoSuchFunction(qualName)
}

func isSystemDatabaseName(name string) bool {
	switch toLower(name) {
	case "information_schema", "performance_schema", "mysql", "sys":
		return true
	default:
		return false
	}
}

func validateTemporalLiteral(lit *nodes.TemporalLit) error {
	if lit == nil {
		return nil
	}
	switch strings.ToUpper(lit.Type) {
	case "DATE":
		s := strings.TrimSpace(lit.Value)
		if s == "" || strings.ContainsAny(s, " \t\n\rTt") {
			return errWrongValue("DATE", lit.Value)
		}
		runs := countDigitRuns(s)
		if (runs == 1 && len(s) > 8) || runs > 3 {
			return errWrongValue("DATE", lit.Value)
		}
		p, ok := parseMySQLStringDateTime("date", 0, s, false)
		if !ok || p.allZero || p.month < 1 || p.day < 1 || p.day > daysInMonth(p.year, p.month) {
			return errWrongValue("DATE", lit.Value)
		}
	case "TIME":
		if _, _, _, _, _, ok := parseTimeValue(6, lit.Value, false, false); !ok {
			return errWrongValue("TIME", lit.Value)
		}
	case "TIMESTAMP":
		s := strings.TrimSpace(lit.Value)
		if s == "" || !hasDateTimeTimePart(s) {
			return errWrongValue("DATETIME", lit.Value)
		}
		p, ok := parseMySQLStringDateTime("datetime", 6, s, false)
		if !ok || p.allZero || p.month < 1 || p.day < 1 || p.day > daysInMonth(p.year, p.month) {
			return errWrongValue("DATETIME", lit.Value)
		}
	}
	return nil
}

func hasDateTimeTimePart(s string) bool {
	if strings.ContainsAny(s, " \t\n\rTt") {
		return true
	}
	runs := countDigitRuns(s)
	if runs <= 2 {
		// Compact digit string with optional fractional seconds:
		// YYYYMMDDhhmmss[.frac] or YYMMDDhhmmss[.frac].
		head := s
		if dot := strings.IndexByte(head, '.'); dot >= 0 {
			head = head[:dot]
		}
		return countDigitRuns(head) == 1 && len(head) >= 12
	}
	// Delimited without whitespace/T: a date-only value has at most 3 digit runs.
	return runs > 3
}

func validateBlockLocking(block *nodes.SelectStmt, currentDB string) error {
	clauses := make([]*nodes.ForUpdate, 0, 1+len(block.MoreLocking))
	if block.ForUpdate != nil {
		clauses = append(clauses, block.ForUpdate)
	}
	clauses = append(clauses, block.MoreLocking...)
	if len(clauses) == 0 {
		return nil
	}
	var relations map[string]bool
	for _, fu := range clauses {
		if fu == nil || len(fu.Tables) == 0 {
			continue
		}
		if relations == nil {
			relations = blockRelationNames(block)
		}
		for _, tbl := range fu.Tables {
			name := tbl.Name
			if tbl.Schema != "" {
				name = tbl.Schema + "." + tbl.Name
			}
			matched := relations[name]
			if !matched && tbl.Schema != "" && currentDB != "" && strings.EqualFold(tbl.Schema, currentDB) {
				matched = relations[tbl.Name]
			}
			if !matched {
				return errUnresolvedTableLock(tbl.Name)
			}
		}
	}
	return nil
}

// validateViewBodyPostAnalysis checks view body restrictions enforced after
// column/table resolution: window names & frames, top-level ROW() in target
// list, top-level DEFAULT(col) in view SELECT list, GROUPING() without WITH
// ROLLUP, and non-constant ESCAPE / AGAINST args.
func validateViewBodyPostAnalysis(sel *nodes.SelectStmt, blocks []*nodes.SelectStmt, viewName string) error {
	if viewName != "" {
		for _, outer := range viewOuterSelectBlocks(sel) {
			for _, item := range outer.TargetList {
				expr := item
				alias := ""
				if rt, ok := item.(*nodes.ResTarget); ok {
					expr = rt.Val
					alias = rt.Name
				}
				for {
					p, ok := expr.(*nodes.ParenExpr)
					if !ok {
						break
					}
					expr = p.Expr
				}
				if d, ok := expr.(*nodes.DefaultExpr); ok && d.Column != "" {
					if alias == "" {
						alias = "DEFAULT(" + d.Column + ")"
					}
					return errColumnAccessDenied("create view", alias, viewName)
				}
			}
		}
	}
	for _, block := range blocks {
		if err := validateBlockWindows(block); err != nil {
			return err
		}
		for _, item := range block.TargetList {
			expr := item
			if rt, ok := item.(*nodes.ResTarget); ok {
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
				return errOperandColumns(1)
			}
		}
		var blockErr error
		inspectBlockOwnNodes(block, func(n nodes.Node) bool {
			if blockErr != nil {
				return false
			}
			switch e := n.(type) {
			case *nodes.FuncCallExpr:
				if !block.WithRollup && strings.EqualFold(e.Name, "grouping") && !e.SpacedParen && e.Schema == "" {
					blockErr = errInvalidGroupFuncUse()
					return false
				}
			case *nodes.LikeExpr:
				if e.Escape != nil && exprNamesColumn(e.Escape) {
					blockErr = errWrongArguments("ESCAPE")
					return false
				}
			case *nodes.MatchExpr:
				if e.Against != nil && exprNamesColumn(e.Against) {
					blockErr = errWrongArguments("AGAINST")
					return false
				}
			}
			return true
		})
		if blockErr != nil {
			return blockErr
		}
	}
	return nil
}

// viewOuterSelectBlocks returns the top-level SELECT blocks of a view body,
// unwrapping ParenSource and walking set-operation arms.
func viewOuterSelectBlocks(sel *nodes.SelectStmt) []*nodes.SelectStmt {
	if sel == nil {
		return nil
	}
	if sel.ParenSource != nil {
		return viewOuterSelectBlocks(sel.ParenSource)
	}
	if sel.SetOp != nodes.SetOpNone {
		return append(viewOuterSelectBlocks(sel.Left), viewOuterSelectBlocks(sel.Right)...)
	}
	return []*nodes.SelectStmt{sel}
}

func validateBlockWindows(block *nodes.SelectStmt) error {
	defined := make(map[string]bool, len(block.WindowClause))
	for _, wd := range block.WindowClause {
		if wd != nil && wd.Name != "" {
			defined[toLower(wd.Name)] = true
		}
	}
	var checkWin func(wd *nodes.WindowDef) error
	checkWin = func(wd *nodes.WindowDef) error {
		if wd == nil {
			return nil
		}
		if wd.RefName != "" && !defined[toLower(wd.RefName)] {
			return errWindowNoSuchWindow(wd.RefName)
		}
		if wd.Frame != nil {
			if wd.Frame.Type == nodes.FrameGroups {
				return errNotSupportedYet("GROUPS")
			}
			if wd.Frame.Exclude != "" {
				return errNotSupportedYet("EXCLUDE")
			}
			if wd.Frame.Type == nodes.FrameRows {
				if (wd.Frame.Start != nil && isIntervalExpr(wd.Frame.Start.Offset)) ||
					(wd.Frame.End != nil && isIntervalExpr(wd.Frame.End.Offset)) {
					return errWindowRangeFrameNumericType()
				}
			}
			if wd.Frame.Start != nil && wd.Frame.Start.Type == nodes.BoundUnboundedFollowing {
				return errWindowFrameStartIllegal()
			}
			if wd.Frame.End != nil && wd.Frame.End.Type == nodes.BoundUnboundedPreceding {
				return errWindowFrameEndIllegal()
			}
		}
		return nil
	}
	for _, wd := range block.WindowClause {
		if err := checkWin(wd); err != nil {
			return err
		}
	}
	var winErr error
	inspectBlockOwnNodes(block, func(n nodes.Node) bool {
		if winErr != nil {
			return false
		}
		if fc, ok := n.(*nodes.FuncCallExpr); ok && fc.Over != nil {
			if err := checkWin(fc.Over); err != nil {
				winErr = err
				return false
			}
		}
		return true
	})
	return winErr
}

func isIntervalExpr(expr nodes.ExprNode) bool {
	for {
		p, ok := expr.(*nodes.ParenExpr)
		if !ok {
			break
		}
		expr = p.Expr
	}
	_, ok := expr.(*nodes.IntervalExpr)
	return ok
}

// viewBody is a view body's query blocks and what its CTEs read. Every query
// block MySQL resolves is judged: the blocks of the body itself, of a set
// operation, of a parenthesis, of a derived table and of a subquery expression,
// and the body of a CTE the body reaches. MySQL resolves a CTE only where it is
// read, and only through the declarations its own scope holds — one declared
// later in the same WITH is not visible, and a nearer one hides an outer one —
// so the walk resolves each reference against the block's own scope and follows
// the reads from the body rather than from every CTE.
type viewBody struct {
	blocks []viewBlock
	// roots are the CTEs the body itself reads, which is where resolution
	// starts.
	roots map[*nodes.CommonTableExpr]bool
	// edges holds, for a CTE, the CTEs its own body reads.
	edges map[*nodes.CommonTableExpr]map[*nodes.CommonTableExpr]bool
	// reads is the set of CTEs the body reaches, filled by reach.
	reads map[*nodes.CommonTableExpr]bool
}

// viewBlock is one query block, with the CTE whose body holds it: nil for the
// blocks of the view body itself.
type viewBlock struct {
	block *nodes.SelectStmt
	owner *nodes.CommonTableExpr
}

// walk records block and the CTEs it reads, then the blocks nested in it. ctes
// are the declarations visible from the enclosing blocks, outermost first.
func (v *viewBody) walk(block *nodes.SelectStmt, ctes []*nodes.CommonTableExpr, owner *nodes.CommonTableExpr, source string) {
	if block == nil {
		return
	}
	visible := append(append([]*nodes.CommonTableExpr{}, ctes...), block.CTEs...)
	v.blocks = append(v.blocks, viewBlock{block: block, owner: owner})
	var children []*nodes.SelectStmt
	var names []string
	inspectBlock := func(root nodes.Node) {
		nodes.Inspect(root, func(n nodes.Node) bool {
			switch e := n.(type) {
			case *nodes.SelectStmt:
				if e == block {
					return true
				}
				// A nested block reads through its own scope: walked below.
				children = append(children, e)
				return false
			case *nodes.CommonTableExpr:
				// Walked below, with the CTE as the body's owner.
				return false
			case *nodes.TableRef:
				// The dummy DUAL reads no CTE, whatever a CTE of that name is.
				if e.Schema == "" && e.Name != "" && !isDualRef(e, source) {
					names = append(names, e.Name)
				}
			}
			return true
		})
	}
	inspectBlock(block)
	if block.ValuesSource != nil {
		for _, row := range block.ValuesSource.Rows {
			for _, item := range row {
				if item != nil {
					inspectBlock(item)
				}
			}
		}
	}
	for _, name := range names {
		cte := visibleCte(visible, name)
		if cte == nil {
			continue
		}
		if owner == nil {
			if v.roots == nil {
				v.roots = map[*nodes.CommonTableExpr]bool{}
			}
			v.roots[cte] = true
			continue
		}
		if v.edges[owner] == nil {
			v.edges[owner] = map[*nodes.CommonTableExpr]bool{}
		}
		v.edges[owner][cte] = true
	}
	for i, cte := range block.CTEs {
		// A CTE reads the declarations of the enclosing blocks and the siblings
		// written before it, not the ones after it, and not itself unless the
		// WITH says RECURSIVE (mysql:8.0.46).
		mine := append(append([]*nodes.CommonTableExpr{}, ctes...), block.CTEs[:i]...)
		if cte.Recursive {
			mine = append(mine, cte)
		}
		v.walk(cte.Select, mine, cte, source)
	}
	for _, child := range children {
		v.walk(child, visible, owner, source)
	}
}

// reach fills reads with the CTEs the body reads, and the CTEs those read, and
// so on: a CTE nothing in the body reaches is one MySQL never resolves.
func (v *viewBody) reach() {
	v.reads = map[*nodes.CommonTableExpr]bool{}
	queue := make([]*nodes.CommonTableExpr, 0, len(v.roots))
	for cte := range v.roots {
		queue = append(queue, cte)
	}
	for len(queue) > 0 {
		cte := queue[0]
		queue = queue[1:]
		if v.reads[cte] {
			continue
		}
		v.reads[cte] = true
		for next := range v.edges[cte] {
			if !v.reads[next] {
				queue = append(queue, next)
			}
		}
	}
}

// visibleCte is the CTE a name reads: the innermost declaration of it wins.
func visibleCte(ctes []*nodes.CommonTableExpr, name string) *nodes.CommonTableExpr {
	for i := len(ctes) - 1; i >= 0; i-- {
		if strings.EqualFold(ctes[i].Name, name) {
			return ctes[i]
		}
	}
	return nil
}

// viewBlockStar reports the star of one query block that names no relation it
// can expand from: a block with no table, or a qualified star whose qualifier is
// not one of the block's relations.
func viewBlockStar(sel *nodes.SelectStmt, source string) error {
	if readsOnlyNoTable(sel, source) {
		for _, item := range sel.TargetList {
			expr := item
			if target, ok := item.(*nodes.ResTarget); ok {
				expr = target.Val
			}
			switch e := expr.(type) {
			case *nodes.StarExpr:
				return errNoTablesUsed()
			case *nodes.ColumnRef:
				if e.Star {
					return errUnknownTableName(e.Table)
				}
			}
		}
		return nil
	}
	relations := blockRelationNames(sel)
	for _, item := range sel.TargetList {
		expr := item
		if target, ok := item.(*nodes.ResTarget); ok {
			expr = target.Val
		}
		if e, ok := expr.(*nodes.ColumnRef); ok && e.Star {
			if qualifier := starQualifier(e); !relations[qualifier] {
				return errUnknownTableName(qualifier)
			}
		}
	}
	return nil
}

// starQualifier is the name a qualified star asks for: its table, qualified by
// its schema when the star wrote one. A relation's name and alias are compared
// as written: the pinned server runs with lower_case_table_names=0, where
// `SELECT T.* FROM a AS t` is 1051 (mysql:8.0.46).
func starQualifier(e *nodes.ColumnRef) string {
	if e.Schema == "" {
		return e.Table
	}
	return e.Schema + "." + e.Table
}

// blockRelationNames is the names a qualified star may carry in the block: each
// relation's alias, or its own name when the statement gave it none.
func blockRelationNames(sel *nodes.SelectStmt) map[string]bool {
	names := make(map[string]bool, len(sel.From))
	var add func(nodes.TableExpr)
	add = func(te nodes.TableExpr) {
		switch e := te.(type) {
		case *nodes.TableRef:
			// An alias hides the name; without one the reference answers to its
			// table name, and to the schema-qualified name too.
			if e.Alias != "" {
				names[e.Alias] = true
				return
			}
			if e.Name != "" {
				names[e.Name] = true
			}
			if e.Schema != "" && e.Name != "" {
				names[e.Schema+"."+e.Name] = true
			}
		case *nodes.SubqueryExpr:
			if e.Alias != "" {
				names[e.Alias] = true
			}
		case *nodes.JsonTableExpr:
			if e.Alias != "" {
				names[e.Alias] = true
			}
		case *nodes.JoinClause:
			add(e.Left)
			add(e.Right)
		}
	}
	for _, te := range sel.From {
		add(te)
	}
	return names
}

// readsOnlyNoTable reports whether the block reads no relation at all, or only
// the dummy DUAL. MySQL expands a star over DUAL from nothing, so
// `CREATE VIEW v AS SELECT * FROM DUAL` is 1096 like a block with no FROM while
// `CREATE VIEW v AS SELECT 1 FROM DUAL` is a view it takes (read on
// mysql:8.0.46). A table really named dual is written quoted and is not the
// dummy: the same server takes the view over a backtick-quoted dual table.
func readsOnlyNoTable(sel *nodes.SelectStmt, source string) bool {
	if len(sel.From) > 1 {
		return false
	}
	if len(sel.From) == 0 {
		return true
	}
	return isDualRef(sel.From[0], source)
}

// isDualRef reports whether a FROM entry is the dummy DUAL. The parser reads the
// DUAL keyword and the quoted table name into the same TableRef, so the text the
// reference starts at decides: a backtick opens a quoted identifier, which names
// a real table. Without the statement's text (a synced definition) the reference
// is read as that table, which can only over-accept.
func isDualRef(expr nodes.TableExpr, source string) bool {
	ref, ok := expr.(*nodes.TableRef)
	if !ok || ref.Schema != "" || !strings.EqualFold(ref.Name, "dual") {
		return false
	}
	start := ref.Loc.Start
	if start < 0 || start >= len(source) {
		return false
	}
	return source[start] != '`'
}

// hasSchemaQualifiedFuncCall reports whether sel contains a schema-qualified
// function call, whose arguments may reference a stored function's context
// rather than a column in the outer SELECT.
func hasSchemaQualifiedFuncCall(sel *nodes.SelectStmt) bool {
	found := false
	nodes.Inspect(sel, func(n nodes.Node) bool {
		if found {
			return false
		}
		if fc, ok := n.(*nodes.FuncCallExpr); ok && fc.Schema != "" {
			found = true
			return false
		}
		return true
	})
	return found
}

// exprNamesColumn reports whether expr reads a column of this scope. A nested
// query has a scope of its own, and a schema-qualified call may name a stored
// function's column.
func exprNamesColumn(expr nodes.ExprNode) bool {
	if expr == nil {
		return false
	}
	found := false
	nodes.Inspect(expr, func(n nodes.Node) bool {
		switch e := n.(type) {
		case *nodes.ColumnRef:
			if !e.Star {
				found = true
			}
			return false
		case *nodes.SubqueryExpr, *nodes.ExistsExpr:
			return false
		case *nodes.FuncCallExpr:
			if e.Schema != "" {
				return false
			}
		}
		return !found
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
	// ALTER VIEW requires the view to exist, and answers by what holds the name:
	// mysql:8.0.46 says 1347 "'db.a' is not VIEW" for a table and 1146 "Table
	// 'db.a' doesn't exist" for nothing.
	if _, exists := db.Views[key]; !exists {
		if db.Tables[key] != nil {
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
	// the one failure MySQL's ALTER VIEW reports too.
	analyzedQuery, err := c.analyzeViewBodyFor(stmt.Select, source, db, stmt.Name.Name)
	if err != nil {
		return err
	}

	// Resolve, rewrite, and deparse the SELECT to produce canonical definition.
	definition, derivedCols := c.deparseViewSelect(stmt.Select, stmt.SelectText, source, db)

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
	case *nodes.CaseExpr:
		if e.Default == nil || inferViewExprNullable(e.Default, relations) {
			return true
		}
		for _, w := range e.Whens {
			if w != nil && inferViewExprNullable(w.Result, relations) {
				return true
			}
		}
		return false
	case *nodes.FuncCallExpr:
		switch strings.ToLower(e.Name) {
		case "nullif", "mod",
			"date_add", "date_sub", "adddate", "subdate",
			"timestampadd", "timestampdiff", "datediff", "timediff", "addtime", "subtime":
			return true
		}
		if (strings.EqualFold(e.Name, "coalesce") || strings.EqualFold(e.Name, "ifnull")) && len(e.Args) > 0 {
			for _, arg := range e.Args {
				if !inferViewExprNullable(arg, relations) {
					return false
				}
			}
			return true
		}
		if strings.EqualFold(e.Name, "if") && len(e.Args) == 3 {
			return inferViewExprNullable(e.Args[1], relations) || inferViewExprNullable(e.Args[2], relations)
		}
		if isViewStringMetadataNullableFunc(e.Name) || isViewNullableAggregateFunc(e.Name) {
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
		if _, ok := e.Left.(*nodes.IntervalExpr); ok {
			return true
		}
		if _, ok := e.Right.(*nodes.IntervalExpr); ok {
			return true
		}
		if e.Op == nodes.BinOpNullSafeEq {
			return false
		}
		// Division and modulo are NULL for a zero divisor (see
		// isAnalyzedExprNullable).
		if e.Op == nodes.BinOpDiv || e.Op == nodes.BinOpDivInt || e.Op == nodes.BinOpMod {
			return true
		}
		return inferViewExprNullable(e.Left, relations) || inferViewExprNullable(e.Right, relations)
	case *nodes.BetweenExpr:
		return inferViewExprNullable(e.Expr, relations) ||
			inferViewExprNullable(e.Low, relations) ||
			inferViewExprNullable(e.High, relations)
	case *nodes.LikeExpr:
		return inferViewExprNullable(e.Expr, relations) ||
			inferViewExprNullable(e.Pattern, relations) ||
			(e.Escape != nil && inferViewExprNullable(e.Escape, relations))
	case *nodes.IsExpr:
		return false
	case *nodes.SubqueryExpr:
		return true
	case *nodes.InExpr:
		if e.Select != nil || inferViewExprNullable(e.Expr, relations) {
			return true
		}
		for _, item := range e.List {
			if inferViewExprNullable(item, relations) {
				return true
			}
		}
		return false
	case *nodes.UnaryExpr:
		return inferViewExprNullable(e.Operand, relations)
	case *nodes.WeightStringExpr:
		// Propagate like the plain WEIGHT_STRING(expr) function-call path.
		return inferViewExprNullable(e.Expr, relations)
	case *nodes.ExtractExpr:
		return inferViewExprNullable(e.Expr, relations)
	case *nodes.IntervalExpr:
		return inferViewExprNullable(e.Value, relations)
	case *nodes.CollateExpr:
		// Nullable whatever its operand, as the engine reports it (see
		// isAnalyzedExprNullable).
		return true
	case *nodes.CastExpr:
		// Covers plain casts and the WEIGHT_STRING(... AS BINARY(n)) desugar,
		// whose WeightStringExpr.Expr is a CastExpr (oracle 8.0.32:
		// IS_NULLABLE=YES for both over a nullable column).
		return inferViewExprNullable(e.Expr, relations)
	case *nodes.PositionExpr:
		// mysql:8.0: IS_NULLABLE is YES for position('x' in s) over a
		// nullable s and for position('x' in NULL), NO over a NOT NULL column.
		return inferViewExprNullable(e.Substr, relations) ||
			inferViewExprNullable(e.Str, relations)
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
	if fn.JSONValueArm {
		for _, arg := range fn.Args {
			if _, err := inferViewExprCollation(arg, relations); err != nil {
				return viewCollationInfo{}, err
			}
		}
		return viewCollationInfo{}, nil
	}
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
	if strings.EqualFold(left.charset, "binary") {
		return left, nil
	}
	if strings.EqualFold(right.charset, "binary") {
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

func isViewNullableAggregateFunc(name string) bool {
	switch strings.ToLower(name) {
	case "sum", "avg", "min", "max", "group_concat":
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
			// A name no view holds is 1347 when a table holds it, and 1051
			// "Unknown table" otherwise (mysql:8.0.46).
			if db.Tables[key] != nil {
				return errWrongObject(db.Name, ref.Name, "VIEW")
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
//
// source is the text the statement was parsed from, which the AST's Loc
// offsets index into. An unaliased select item is named after its text in
// source, as the engine names it; with no source ("" from DefineView), the
// name is rebuilt from the AST.
func (c *Catalog) deparseViewSelect(sel *nodes.SelectStmt, rawText, source string, db *Database) (string, []string) {
	if sel == nil {
		return rawText, nil
	}

	// Build lookups for table/view metadata: Lookup resolves unqualified names
	// from the database, while SchemaLookup resolves schema-qualified names
	// against their named schema in the catalog.
	lookup := tableLookupForDB(db)
	schemaLookup := c.schemaTableLookup(db)

	// Determine the database charset for CAST resolution.
	charset := db.Charset
	if charset == "" {
		charset = c.defaultCharset
	}

	// Resolve: qualify columns, expand *, normalize JOINs.
	resolver := &deparse.Resolver{
		Lookup:              lookup,
		SchemaLookup:        schemaLookup,
		DefaultCharset:      charset,
		DefaultWeekFormat:   c.session.DefaultWeekFormat,
		NormalizeCastAsYear: !c.hasStoredFunction(db, "cast_as_year"),
		Source:              source,
		DefaultSchema:       db.Name,
	}
	resolver.Resolve(sel)

	// Extract column names from the resolved target list.
	derivedCols := extractViewColumns(sel, lookup, schemaLookup)

	// Rewrite: NOT folding, boolean context wrapping.
	deparse.RewriteSelectStmt(sel)

	// Deparse: AST → canonical SQL text.
	return deparse.DeparseSelect(sel), derivedCols
}

// extractViewColumns extracts column names from a resolved SELECT target list.
// This produces the column list that MySQL would derive for a view. lookup and
// schemaLookup resolve the table of a TABLE query primary, which has no target
// list to read.
func extractViewColumns(sel *nodes.SelectStmt, lookup deparse.TableLookup, schemaLookup deparse.SchemaTableLookup) []string {
	if sel == nil {
		return nil
	}
	// The result columns live on the leftmost output SELECT. Unwrap any
	// parenthesized-query wrappers and walk set-operation left arms to reach it,
	// so a parenthesized or set-op view body (CREATE VIEW v AS (SELECT a, b) or
	// (SELECT 1 AS a) UNION (SELECT 2 AS a)) derives columns from that SELECT.
	// A TABLE query primary is SELECT * FROM t and a VALUES primary a row
	// constructor, so neither has a target list to read. Their columns are
	// derived by the same code the resolver uses for a CTE or derived table
	// over one, which searches the query's own CTE scope before lookup's
	// tables (container-verified: CREATE VIEW v AS TABLE t declares the
	// columns of t, CREATE VIEW v AS VALUES ROW(1,2) declares column_0 and
	// column_1, and WITH x AS (SELECT 1 AS a) TABLE x declares a).
	if primary := deparse.PrimaryOutputColumnsWithSchema(sel, lookup, schemaLookup); primary != nil {
		cols := make([]string, 0, len(primary))
		for _, c := range primary {
			cols = append(cols, c.Name)
		}
		return cols
	}
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

// tableLookupForDB returns a deparse.TableLookup function that resolves table
// and view names from the given database's Tables and Views maps.
func tableLookupForDB(db *Database) deparse.TableLookup {
	return func(tableName string) *deparse.ResolverTable {
		key := toLower(tableName)
		// Try tables first (excluding system-hidden columns such as functional-index backing columns).
		tbl := db.Tables[key]
		if tbl != nil {
			visible := tbl.VisibleColumns()
			cols := make([]deparse.ResolverColumn, len(visible))
			for i, c := range visible {
				cols[i] = deparse.ResolverColumn{
					Name:      c.Name,
					Position:  c.Position,
					DataType:  toLower(c.DataType),
					Scale:     columnDeparseScale(c),
					Invisible: c.Invisible,
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

// schemaTableLookup returns a deparse.SchemaTableLookup that resolves
// schema-qualified table/view names in the catalog (or falls back to defaultDB
// when schema is empty).
func (c *Catalog) schemaTableLookup(defaultDB *Database) deparse.SchemaTableLookup {
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

func (c *Catalog) hasStoredFunction(viewDB *Database, lowerName string) bool {
	if c != nil {
		if curDB := c.GetDatabase(c.CurrentDatabase()); curDB != nil && curDB.Functions[lowerName] != nil {
			return true
		}
	}
	return viewDB != nil && viewDB.Functions[lowerName] != nil
}

// columnDeparseScale returns a column's fractional seconds precision (for
// TIME/DATETIME/TIMESTAMP) or decimal scale (for DECIMAL/NUMERIC) for view
// deparsing.
func columnDeparseScale(col *Column) int {
	if col == nil {
		return 0
	}
	switch toLower(col.DataType) {
	case "datetime", "timestamp", "time":
		return columnTemporalFSP(col)
	case "decimal", "numeric":
		if pt := parseColumnType(col.ColumnType); pt.scale > 0 {
			return pt.scale
		}
	}
	return 0
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
	b.WriteString(" VIEW " + backtickIDQ(v.Name))

	// Column list (only if explicitly specified by user in CREATE VIEW).
	if v.ExplicitColumns && len(v.Columns) > 0 {
		cols := make([]string, len(v.Columns))
		for i, c := range v.Columns {
			cols[i] = backtickIDQ(c)
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
