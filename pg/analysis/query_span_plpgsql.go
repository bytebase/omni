package analysis

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	omnipg "github.com/bytebase/omni/pg"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/catalog"
	plpgsqlast "github.com/bytebase/omni/pg/plpgsql/ast"
	plpgsqlparser "github.com/bytebase/omni/pg/plpgsql/parser"
)

// snapshotFunc is a function's language and body as the schema snapshot
// defines it.
type snapshotFunc struct {
	language string
	body     string
}

type funcBodyAnalysis struct {
	sourceColumns    []ColumnSet
	predicateColumns ColumnSet
}

// analyzeFunctionBody analyzes a user-defined function body and returns per-column
// source column sets. Results are cached by function OID on the extractor.
// A nil value in the cache means analysis is in progress (recursion sentinel).
// On parse errors, it degrades gracefully by returning empty sets.
func (e *querySpanExtractor) analyzeFunctionBody(proc *catalog.UserProc) []ColumnSet {
	analysis := e.analyzeFunctionBodyForProc(proc)
	e.mergeFunctionSourceTables(analysis.sourceColumns)
	mergeSourceColumnSetInto(e.predicateColumns, analysis.predicateColumns)
	return analysis.sourceColumns
}

func (e *querySpanExtractor) analyzeFunctionBodyForProc(proc *catalog.UserProc) *funcBodyAnalysis {
	// Check cache.
	if cached, ok := e.funcBodyCache[proc.OID]; ok {
		if cached == nil {
			// In-progress sentinel — recursive call. Return empty sets.
			return &funcBodyAnalysis{
				sourceColumns:    makeEmptySets(proc),
				predicateColumns: make(ColumnSet),
			}
		}
		return cached
	}

	// Store nil sentinel to detect recursion.
	e.funcBodyCache[proc.OID] = nil
	outerPredicateColumns := e.predicateColumns
	bodyPredicateColumns := make(ColumnSet)
	e.predicateColumns = bodyPredicateColumns
	defer func() {
		e.predicateColumns = outerPredicateColumns
	}()

	// The catalog keeps each overload's own body. A function it could not
	// install is a stand-in with a placeholder body, and the snapshot cannot
	// tell which overload it stands for, so every overload's body is analyzed
	// and merged.
	defs := []snapshotFunc{{language: proc.Language, body: proc.Body}}
	if proc.Schema != nil {
		key := strings.ToLower(proc.Schema.Name + "." + proc.Name)
		if e.standInFuncs[key] && len(e.snapshotFuncs[key]) > 0 {
			defs = e.snapshotFuncs[key]
		}
	}

	var result []ColumnSet
	var err error
	for _, def := range defs {
		defResult, defErr := e.analyzeBodyAs(proc, def)
		if defErr != nil {
			err = defErr
			continue
		}
		if len(defs) == 1 {
			result = defResult
			break
		}
		for i, set := range defResult {
			if i >= len(result) {
				result = append(result, make(ColumnSet))
			}
			mergeSourceColumnSetInto(result[i], set)
		}
	}

	if err != nil && result == nil {
		// On error, cache empty sets so we don't retry. The conditions found
		// before the error still run, so they stay.
		e.funcBodyCache[proc.OID] = &funcBodyAnalysis{
			sourceColumns:    makeEmptySets(proc),
			predicateColumns: cloneSourceColumnSet(bodyPredicateColumns),
		}
		return e.funcBodyCache[proc.OID]
	}

	e.funcBodyCache[proc.OID] = &funcBodyAnalysis{
		sourceColumns:    result,
		predicateColumns: cloneSourceColumnSet(bodyPredicateColumns),
	}
	return e.funcBodyCache[proc.OID]
}

func (e *querySpanExtractor) analyzeBodyAs(proc *catalog.UserProc, def snapshotFunc) ([]ColumnSet, error) {
	switch strings.ToLower(def.language) {
	case "plpgsql":
		return e.analyzePLpgSQLBody(proc, def.body)
	case "sql":
		return e.analyzeSQLBody(def.body)
	default:
		// Unsupported language (C, internal, etc.) — unknown lineage.
		return makeEmptySets(proc), nil
	}
}

func cloneSourceColumnSet(src ColumnSet) ColumnSet {
	dst := make(ColumnSet)
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func mergeSourceColumnSetInto(dst, src ColumnSet) {
	for k, v := range src {
		dst[k] = v
	}
}

func (e *querySpanExtractor) mergeFunctionSourceTables(bodySets []ColumnSet) {
	for _, colSet := range bodySets {
		for k := range colSet {
			e.funcSourceColumns[ColumnResource{
				Database: k.Database,
				Schema:   k.Schema,
				Table:    k.Table,
			}] = true
		}
	}
}

func makeEmptySets(proc *catalog.UserProc) []ColumnSet {
	n := countOutputParams(proc)
	if n == 0 {
		n = 1 // scalar return
	}
	sets := make([]ColumnSet, n)
	for i := range sets {
		sets[i] = make(ColumnSet)
	}
	return sets
}

func countOutputParams(proc *catalog.UserProc) int {
	count := 0
	for _, mode := range proc.ArgModes {
		if mode == 'o' || mode == 't' || mode == 'b' {
			count++
		}
	}
	return count
}

func (e *querySpanExtractor) analyzeSQLBody(body string) ([]ColumnSet, error) {
	if body == "" {
		return nil, nil
	}

	stmts, err := omnipg.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SQL function body: %w", err)
	}

	// Every statement in the body runs, so each one's conditions count; the
	// last one returns the function's result.
	last := -1
	for i, stmt := range stmts {
		if stmt.AST != nil {
			last = i
		}
	}
	var result []ColumnSet
	for i, stmt := range stmts {
		if stmt.AST == nil {
			continue
		}
		accesses := e.accessTablesOf(stmt.Text)
		// The caller reads what the body reads, whether or not the result shows it.
		mergeSourceColumnSetInto(e.funcSourceColumns, accesses)
		selStmt, ok := stmt.AST.(*ast.SelectStmt)
		if !ok {
			// The analyzer reads SELECT alone, so another statement's conditions
			// and RETURNING list are matched by name.
			e.collectFallbackPredicateColumns(stmt.AST, accesses)
			returning := e.returningLineage(stmt.AST, accesses)
			if i == last {
				result = returning
			} else {
				for _, set := range returning {
					mergeSourceColumnSetInto(e.predicateColumns, set)
				}
			}
			continue
		}
		query, err := e.analyzeSelect(selStmt)
		if err != nil {
			e.collectFallbackPredicateColumns(selStmt, accesses)
			if i != last {
				e.collectDiscardedOutputColumns(nil, selStmt, accesses)
				continue
			}
			// Fall back to parse-tree extraction when analysis fails
			// (e.g., nested function calls, parameter references).
			fallbackResults := e.extractFallbackColumns(selStmt)
			if len(fallbackResults) == 0 {
				return nil, fmt.Errorf("failed to analyze SQL function body: %w", err)
			}
			result = make([]ColumnSet, len(fallbackResults))
			for j, fr := range fallbackResults {
				result[j] = fr.SourceColumns
			}
			continue
		}
		e.collectStatementPredicateColumns(query, selStmt, accesses)
		if i == last {
			result = e.extractFuncLineage(query)
		} else {
			e.collectDiscardedOutputColumns(query, selStmt, accesses)
		}
	}
	return result, nil
}

func (e *querySpanExtractor) extractFuncLineage(q *catalog.Query) []ColumnSet {
	var result []ColumnSet
	for _, te := range q.TargetList {
		if te.ResJunk {
			continue
		}
		colSet := make(ColumnSet)
		e.walkExpr(q, te.Expr, colSet)
		result = append(result, colSet)
	}
	return result
}

// variableScope tracks PL/pgSQL variable source columns with nested scoping.
type variableScope struct {
	vars map[string]ColumnSet
	// declared maps each name the block declares to "", or for an ALIAS to
	// the lowercase name it stands for.
	declared map[string]string
	parent   *variableScope
}

func newVariableScope(parent *variableScope) *variableScope {
	return &variableScope{
		vars:     make(map[string]ColumnSet),
		declared: make(map[string]string),
		parent:   parent,
	}
}

func (s *variableScope) set(name string, sources ColumnSet) {
	s.vars[strings.ToLower(name)] = sources
}

func (s *variableScope) get(name string) (ColumnSet, bool) {
	lower := strings.ToLower(name)
	for scope := s; scope != nil; scope = scope.parent {
		if v, ok := scope.vars[lower]; ok {
			return v, true
		}
	}
	return nil, false
}

// merge adds sources to the variable name in the scope that holds it, so an
// assignment in a nested block or loop outlives it. A name no scope holds,
// such as an output parameter, lives as long as the function, in the
// outermost scope.
func (s *variableScope) merge(name string, sources ColumnSet) {
	lower := strings.ToLower(name)
	var owner *variableScope
	for scope := s; scope != nil; scope = scope.parent {
		if _, ok := scope.vars[lower]; ok {
			owner = scope
			break
		}
		owner = scope
	}
	merged := cloneSourceColumnSet(owner.vars[lower])
	mergeSourceColumnSetInto(merged, sources)
	owner.vars[lower] = merged
}

// plpgsqlAnalyzer holds state for analyzing a single PL/pgSQL function body.
type plpgsqlAnalyzer struct {
	extractor *querySpanExtractor
	scope     *variableScope
	// returnResults collects per-column source sets from all RETURN QUERY paths.
	returnResults [][]ColumnSet
	// outputs holds what each output parameter (OUT, INOUT, TABLE) was
	// assigned, keyed by $n and by its lowercase name, and outputNames lists
	// the $n keys in declared order.
	outputs     map[string]ColumnSet
	outputNames []string
	// returnsOutputs is set by a bare RETURN NEXT, which returns the output
	// parameters as a row.
	returnsOutputs bool
	// cteMap holds CTE definitions for resolving column refs through CTEs.
	cteMap       map[string]*ast.SelectStmt
	resolvingCTE map[cteResolveKey]bool
	// qualifiedFrom marks fromTables keys whose FROM item carried an
	// explicit schema qualifier, so fallback resolution does not probe the
	// search path for them. Reset per fallback statement; nested subquery
	// collection during one statement accumulates into the same scope.
	qualifiedFrom map[string]bool
	// fromItems maps a FROM name whose columns are not its table's, the
	// alias of a subquery or a join or a table alias with a column list, to
	// that FROM item. It is scoped as qualifiedFrom is.
	fromItems map[string]ast.Node
}

type cteResolveKey struct {
	sel *ast.SelectStmt
	col string
}

func (e *querySpanExtractor) analyzePLpgSQLBody(proc *catalog.UserProc, body string) ([]ColumnSet, error) {
	if body == "" {
		return makeEmptySets(proc), nil
	}

	block, err := plpgsqlparser.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PL/pgSQL function body: %w", err)
	}

	analyzer := &plpgsqlAnalyzer{
		extractor: e,
		scope:     newVariableScope(nil),
		outputs:   make(map[string]ColumnSet),
	}
	setReturning := proc.RetSet
	for i, mode := range proc.ArgModes {
		if mode != 'o' && mode != 't' && mode != 'b' {
			continue
		}
		setReturning = setReturning || mode == 't'
		// $n refers to a parameter whether or not it has a name.
		key := fmt.Sprintf("$%d", i+1)
		analyzer.outputs[key] = make(ColumnSet)
		if i < len(proc.ArgNames) && proc.ArgNames[i] != "" {
			analyzer.outputs[strings.ToLower(proc.ArgNames[i])] = analyzer.outputs[key]
		}
		analyzer.outputNames = append(analyzer.outputNames, key)
	}

	// Add function input parameters to scope so they can be substituted in
	// embedded SQL. Output/table params ('o', 't') are NOT added because their
	// names often collide with column names in the function body's queries
	// (e.g., RETURNS TABLE(a int, b int) where a, b are also column names).
	for i, name := range proc.ArgNames {
		if name == "" {
			continue
		}
		if i < len(proc.ArgModes) {
			mode := proc.ArgModes[i]
			if mode == 'o' || mode == 't' {
				continue
			}
		}
		analyzer.scope.set(name, make(ColumnSet))
	}

	analyzer.walkBlock(block)

	// A function with output parameters that returns no set returns them when
	// it ends; a set-returning one returns them at each bare RETURN NEXT.
	if len(analyzer.outputNames) > 0 && (analyzer.returnsOutputs || !setReturning) {
		row := make([]ColumnSet, len(analyzer.outputNames))
		for i, key := range analyzer.outputNames {
			row[i] = analyzer.outputs[key]
		}
		analyzer.returnResults = append(analyzer.returnResults, row)
	}

	nCols := countOutputParams(proc)
	if nCols == 0 {
		nCols = 1
	}

	return analyzer.mergeReturnResults(nCols), nil
}

func (a *plpgsqlAnalyzer) walkBlock(block *plpgsqlast.PLBlock) {
	// Push scope for declarations.
	a.scope = newVariableScope(a.scope)
	defer func() { a.scope = a.scope.parent }()

	// Process declarations.
	for _, decl := range block.Declarations {
		switch d := decl.(type) {
		case *plpgsqlast.PLDeclare:
			sources := make(ColumnSet)
			if d.Default != "" {
				sources = a.analyzeEmbeddedExpr(d.Default)
			}
			a.scope.set(d.Name, sources)
			a.scope.declared[strings.ToLower(d.Name)] = ""
		case *plpgsqlast.PLCursorDecl:
			// A cursor carries the columns of its query to what FETCH reads.
			a.scope.set(d.Name, a.analyzeQuerySources(d.Query))
			a.scope.declared[strings.ToLower(d.Name)] = ""
		case *plpgsqlast.PLAliasDecl:
			// Alias references another variable — copy its sources.
			if sources, ok := a.scope.get(d.RefName); ok {
				a.scope.set(d.Name, sources)
			} else {
				a.scope.set(d.Name, make(ColumnSet))
			}
			a.scope.declared[strings.ToLower(d.Name)] = strings.ToLower(d.RefName)
		default:
		}
	}

	// Process body statements.
	for _, stmt := range block.Body {
		a.walkStmt(stmt)
	}

	// Process exception handlers.
	for _, ex := range block.Exceptions {
		if ew, ok := ex.(*plpgsqlast.PLExceptionWhen); ok {
			for _, stmt := range ew.Body {
				a.walkStmt(stmt)
			}
		}
	}
}

func (a *plpgsqlAnalyzer) walkStmt(node plpgsqlast.Node) {
	switch stmt := node.(type) {
	case *plpgsqlast.PLAssign:
		a.handleAssign(stmt)
	case *plpgsqlast.PLReturnQuery:
		a.handleReturnQuery(stmt)
	case *plpgsqlast.PLReturn:
		a.handleReturn(stmt)
	case *plpgsqlast.PLReturnNext:
		a.handleReturnNext(stmt)
	case *plpgsqlast.PLExecSQL:
		a.handleExecSQL(stmt)
	case *plpgsqlast.PLIf:
		a.handleIf(stmt)
	case *plpgsqlast.PLCase:
		a.handleCase(stmt)
	case *plpgsqlast.PLLoop:
		for _, s := range stmt.Body {
			a.walkStmt(s)
		}
	case *plpgsqlast.PLWhile:
		a.handleCondition(stmt.Condition)
		for _, s := range stmt.Body {
			a.walkStmt(s)
		}
	case *plpgsqlast.PLForI:
		a.handleCondition(stmt.Lower)
		a.handleCondition(stmt.Upper)
		a.handleCondition(stmt.Step)
		// The loop declares its integer variable, which shadows any other
		// of its name inside the loop.
		a.scope = newVariableScope(a.scope)
		a.scope.set(stmt.Var, make(ColumnSet))
		a.scope.declared[strings.ToLower(stmt.Var)] = ""
		for _, s := range stmt.Body {
			a.walkStmt(s)
		}
		a.scope = a.scope.parent
	case *plpgsqlast.PLForS:
		a.handleForS(stmt)
	case *plpgsqlast.PLForC:
		// The loop declares its record variable, as an integer loop does.
		sources, _ := a.scope.get(stmt.CursorVar)
		a.scope = newVariableScope(a.scope)
		a.scope.set(stmt.Var, cloneSourceColumnSet(sources))
		a.scope.declared[strings.ToLower(stmt.Var)] = ""
		for _, s := range stmt.Body {
			a.walkStmt(s)
		}
		a.scope = a.scope.parent
	case *plpgsqlast.PLForDynS:
		for _, s := range stmt.Body {
			a.walkStmt(s)
		}
	case *plpgsqlast.PLForEachA:
		a.handleCondition(stmt.ArrayExpr)
		for _, s := range stmt.Body {
			a.walkStmt(s)
		}
	case *plpgsqlast.PLBlock:
		a.walkBlock(stmt)
	// PERFORM discards its rows but still runs, so how long it takes reveals
	// what it reads.
	case *plpgsqlast.PLPerform:
		a.handleCondition(stmt.Expr)
	case *plpgsqlast.PLDynExecute:
		a.handleDynExecute(stmt)
	case *plpgsqlast.PLOpen:
		if stmt.Query != "" {
			a.scope.merge(stmt.CursorVar, a.analyzeQuerySources(stmt.Query))
		}
	case *plpgsqlast.PLFetch:
		if sources, ok := a.scope.get(stmt.CursorVar); ok {
			for _, varName := range stmt.Into {
				a.assign(varName, sources)
			}
		}
	case *plpgsqlast.PLCall:
		a.handleCall(stmt)
	// RAISE puts the values of its parameters in the message the caller sees.
	case *plpgsqlast.PLRaise:
		for _, param := range stmt.Params {
			a.handleCondition(param)
		}
	case *plpgsqlast.PLExit:
		a.handleCondition(stmt.Condition)
	case *plpgsqlast.PLAssert:
		a.handleCondition(stmt.Condition)
	// Skip nodes that don't affect lineage:
	// PLNull, PLCommit, PLRollback, PLClose, PLGetDiag
	default:
	}
}

// handleCondition records the columns an expression that steers the body
// reads, directly or through variables: it decides which statements run, and
// so what the function returns and how long it takes.
func (a *plpgsqlAnalyzer) handleCondition(expr string) {
	if strings.TrimSpace(expr) == "" {
		return
	}
	mergeSourceColumnSetInto(a.extractor.predicateColumns, a.analyzeEmbeddedExpr(expr))
	stmts, err := omnipg.Parse("SELECT " + expr)
	if err != nil || len(stmts) != 1 {
		return
	}
	a.addScopeReferences(stmts[0].AST)
}

// addScopeReferences records, as predicates, the columns that the variables
// node reads carry.
func (a *plpgsqlAnalyzer) addScopeReferences(node ast.Node) {
	mergeSourceColumnSetInto(a.extractor.predicateColumns, a.scopeSources(node))
}

// scopeSources returns the columns that the variables node reads carry. A
// variable, or the record variable a qualified reference such as r.ssn names,
// carries the columns it was assigned from.
func (a *plpgsqlAnalyzer) scopeSources(node ast.Node) ColumnSet {
	sources := make(ColumnSet)
	ast.Inspect(node, func(n ast.Node) bool {
		if ref, ok := n.(*ast.ColumnRef); ok && ref.Fields != nil && len(ref.Fields.Items) > 0 {
			if name, ok := ref.Fields.Items[0].(*ast.String); ok {
				if varSources, ok := a.scope.get(name.Str); ok {
					mergeSourceColumnSetInto(sources, varSources)
				}
			}
		}
		return true
	})
	return sources
}

// analyzeQuerySources returns every column the output of sql reads, merged.
func (a *plpgsqlAnalyzer) analyzeQuerySources(sql string) ColumnSet {
	merged := make(ColumnSet)
	if strings.TrimSpace(sql) == "" {
		return merged
	}
	for _, set := range a.analyzeEmbeddedSQL(sql) {
		mergeSourceColumnSetInto(merged, set)
	}
	return merged
}

// handleCall analyzes the body of the procedure CALL invokes, whatever
// overload it is. The procedure runs for its effects, so the values passed to
// it count as predicates.
func (a *plpgsqlAnalyzer) handleCall(stmt *plpgsqlast.PLCall) {
	stmts, err := omnipg.Parse(stmt.SQLText)
	if err != nil || len(stmts) != 1 {
		return
	}
	call, ok := stmts[0].AST.(*ast.CallStmt)
	if !ok || call.Funccall == nil {
		return
	}
	for _, proc := range a.extractor.userProcsNamed(astStrings(call.Funccall.Funcname)) {
		a.extractor.analyzeFunctionBody(proc)
	}
	if call.Funccall.Args == nil {
		return
	}
	for _, arg := range call.Funccall.Args.Items {
		loc := ast.NodeLoc(arg)
		if loc.Start >= 0 && loc.End > loc.Start && loc.End <= len(stmt.SQLText) {
			a.handleCondition(stmt.SQLText[loc.Start:loc.End])
			continue
		}
		a.addScopeReferences(arg)
	}
}

func (a *plpgsqlAnalyzer) handleAssign(stmt *plpgsqlast.PLAssign) {
	varName := extractBaseVarName(stmt.Target)
	sources := a.analyzeEmbeddedExpr(stmt.Expr)
	a.assign(varName, sources)
}

// assign merges sources into the variable name, and into the output
// parameter it names, if any.
func (a *plpgsqlAnalyzer) assign(name string, sources ColumnSet) {
	a.scope.merge(name, sources)
	if out, ok := a.outputParam(name); ok {
		mergeSourceColumnSetInto(out, sources)
	}
}

// outputParam returns the output parameter that name refers to. A variable a
// block declares shadows a parameter of its name; an ALIAS refers to what it
// stands for.
func (a *plpgsqlAnalyzer) outputParam(name string) (ColumnSet, bool) {
	name = strings.ToLower(extractBaseVarName(name))
	scope := a.scope
	// The bound stops a cycle of aliases.
	for range 100 {
		for scope != nil {
			if _, ok := scope.declared[name]; ok {
				break
			}
			scope = scope.parent
		}
		if scope == nil {
			out, ok := a.outputs[name]
			return out, ok
		}
		target := scope.declared[name]
		if target == "" || target == name {
			return nil, false
		}
		name = target
	}
	return nil, false
}

func (a *plpgsqlAnalyzer) handleReturnQuery(stmt *plpgsqlast.PLReturnQuery) {
	if stmt.Query == "" {
		return
	}
	results := a.analyzeEmbeddedSQL(stmt.Query)
	if len(results) > 0 {
		a.returnResults = append(a.returnResults, results)
	}
}

func (a *plpgsqlAnalyzer) handleReturn(stmt *plpgsqlast.PLReturn) {
	if stmt.Expr == "" {
		return
	}
	sources := a.analyzeEmbeddedExpr(stmt.Expr)
	a.returnResults = append(a.returnResults, []ColumnSet{sources})
}

func (a *plpgsqlAnalyzer) handleReturnNext(stmt *plpgsqlast.PLReturnNext) {
	if stmt.Expr == "" {
		a.returnsOutputs = true
		return
	}
	sources := a.analyzeEmbeddedExpr(stmt.Expr)
	a.returnResults = append(a.returnResults, []ColumnSet{sources})
}

func (a *plpgsqlAnalyzer) handleExecSQL(stmt *plpgsqlast.PLExecSQL) {
	stmts, err := omnipg.Parse(stmt.SQLText)
	if err != nil || len(stmts) != 1 {
		return
	}
	var results []ColumnSet
	if _, ok := stmts[0].AST.(*ast.SelectStmt); ok {
		results = a.analyzeEmbeddedSQL(stmt.SQLText)
	} else {
		// The analyzer reads SELECT alone, so another statement, which runs for
		// its effect, has its conditions and RETURNING list matched by name.
		accesses := a.extractor.accessTablesOf(stmt.SQLText)
		a.extractor.collectFallbackPredicateColumns(stmts[0].AST, accesses)
		results = a.extractor.returningLineage(stmts[0].AST, accesses)
	}
	for i, varName := range stmt.Into {
		if i < len(results) {
			a.assign(varName, results[i])
		}
	}
}

// handleDynExecute analyzes EXECUTE of a constant query like the statement it
// runs. A query built at run time cannot be analyzed.
func (a *plpgsqlAnalyzer) handleDynExecute(stmt *plpgsqlast.PLDynExecute) {
	stmts, err := omnipg.Parse("SELECT " + stmt.Query)
	if err != nil || len(stmts) != 1 {
		return
	}
	sel, ok := stmts[0].AST.(*ast.SelectStmt)
	if !ok || sel.TargetList == nil || len(sel.TargetList.Items) != 1 {
		return
	}
	rt, ok := sel.TargetList.Items[0].(*ast.ResTarget)
	if !ok {
		return
	}
	query, ok := constantString(rt.Val)
	if !ok {
		return
	}
	a.handleExecSQL(&plpgsqlast.PLExecSQL{SQLText: query, Into: stmt.Into})
}

// constantString returns the text of a string literal, a cast of one, or a
// concatenation of them.
func constantString(node ast.Node) (string, bool) {
	switch v := node.(type) {
	case *ast.A_Const:
		if s, ok := v.Val.(*ast.String); ok {
			return s.Str, true
		}
	case *ast.TypeCast:
		return constantString(v.Arg)
	case *ast.A_Expr:
		if names := astStrings(v.Name); len(names) == 1 && names[0] == "||" {
			left, ok := constantString(v.Lexpr)
			if !ok {
				return "", false
			}
			right, ok := constantString(v.Rexpr)
			if !ok {
				return "", false
			}
			return left + right, true
		}
	default:
	}
	return "", false
}

func (a *plpgsqlAnalyzer) handleForS(stmt *plpgsqlast.PLForS) {
	results := a.analyzeEmbeddedSQL(stmt.Query)
	// A list of scalar targets takes the columns in order; a single record
	// or row target takes them all.
	if targets := strings.Split(stmt.Var, ","); len(targets) > 1 {
		for i, target := range targets {
			if i < len(results) {
				a.assign(strings.TrimSpace(target), results[i])
			}
		}
	} else {
		merged := make(ColumnSet)
		for _, colSet := range results {
			mergeSourceColumnSetInto(merged, colSet)
		}
		a.assign(stmt.Var, merged)
	}
	for _, s := range stmt.Body {
		a.walkStmt(s)
	}
}

func (a *plpgsqlAnalyzer) handleIf(stmt *plpgsqlast.PLIf) {
	a.handleCondition(stmt.Condition)
	for _, s := range stmt.ThenBody {
		a.walkStmt(s)
	}
	for _, elsif := range stmt.ElsIfs {
		if ei, ok := elsif.(*plpgsqlast.PLElsIf); ok {
			a.handleCondition(ei.Condition)
			for _, s := range ei.Body {
				a.walkStmt(s)
			}
		}
	}
	for _, s := range stmt.ElseBody {
		a.walkStmt(s)
	}
}

func (a *plpgsqlAnalyzer) handleCase(stmt *plpgsqlast.PLCase) {
	a.handleCondition(stmt.TestExpr)
	for _, when := range stmt.Whens {
		if cw, ok := when.(*plpgsqlast.PLCaseWhen); ok {
			a.handleCondition(cw.Expr)
			for _, s := range cw.Body {
				a.walkStmt(s)
			}
		}
	}
	for _, s := range stmt.ElseBody {
		a.walkStmt(s)
	}
}

func (a *plpgsqlAnalyzer) analyzeEmbeddedSQL(sql string) []ColumnSet {
	substituted := a.substituteVariables(sql)

	// Always collect access tables from embedded SQL for the top-level QuerySpan,
	// regardless of whether full analysis succeeds.
	accesses := a.collectAccessTablesFromSQL(substituted)

	stmts, err := omnipg.Parse(substituted)
	if err != nil {
		return nil
	}
	if len(stmts) == 0 {
		return nil
	}

	selStmt, ok := stmts[0].AST.(*ast.SelectStmt)
	if !ok {
		return nil
	}

	query, err := a.extractor.analyzeSelect(selStmt)
	if err != nil {
		// When full analysis fails (e.g., type mismatches with aggregate functions),
		// fall back to extracting column references from the parse tree and resolving
		// them against known tables in the catalog.
		a.extractor.collectFallbackPredicateColumns(selStmt, accesses)
		return a.fallbackExtractLineage(sql, selStmt)
	}

	a.extractor.collectStatementPredicateColumns(query, selStmt, accesses)

	results := a.extractor.extractFuncLineage(query)

	// Map back variable references.
	origSelectItems := extractSelectItemNames(sql)
	for i, colSet := range results {
		if len(colSet) == 0 && i < len(origSelectItems) {
			varName := strings.TrimSpace(origSelectItems[i])
			if sources, ok := a.scope.get(varName); ok {
				results[i] = sources
			}
		}
	}

	return results
}

func extractSelectItemNames(sql string) []string {
	upper := strings.ToUpper(strings.TrimSpace(sql))
	idx := strings.Index(upper, "SELECT")
	if idx < 0 {
		return nil
	}
	rest := sql[idx+6:]

	fromIdx := -1
	depth := 0
	upperRest := strings.ToUpper(rest)
	for i := 0; i < len(upperRest); i++ {
		switch upperRest[i] {
		case '(':
			depth++
		case ')':
			depth--
		default:
			if depth == 0 && i+4 <= len(upperRest) && upperRest[i:i+4] == "FROM" {
				if (i == 0 || !isWordChar(rest[i-1])) && (i+4 >= len(rest) || !isWordChar(rest[i+4])) {
					fromIdx = i
				}
			}
		}
		if fromIdx >= 0 {
			break
		}
	}

	targetList := rest
	if fromIdx >= 0 {
		targetList = rest[:fromIdx]
	}

	var items []string
	depth = 0
	start := 0
	for i := 0; i < len(targetList); i++ {
		switch targetList[i] {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				items = append(items, strings.TrimSpace(targetList[start:i]))
				start = i + 1
			}
		default:
		}
	}
	items = append(items, strings.TrimSpace(targetList[start:]))

	return items
}

func (a *plpgsqlAnalyzer) analyzeEmbeddedExpr(expr string) ColumnSet {
	results := a.analyzeEmbeddedSQL("SELECT " + expr)
	merged := make(ColumnSet)
	for _, colSet := range results {
		for k, v := range colSet {
			merged[k] = v
		}
	}
	if len(merged) == 0 {
		if stmts, err := omnipg.Parse("SELECT " + expr); err == nil && len(stmts) == 1 {
			mergeSourceColumnSetInto(merged, a.scopeSources(stmts[0].AST))
		}
	}
	return merged
}

func (a *plpgsqlAnalyzer) substituteVariables(sql string) string {
	stmts, err := omnipg.Parse(sql)
	if err == nil && len(stmts) > 0 {
		if selStmt, ok := stmts[0].AST.(*ast.SelectStmt); ok {
			if _, err := a.extractor.analyzeSelect(selStmt); err == nil {
				return sql
			}
		}
	}

	result := sql
	for scope := a.scope; scope != nil; scope = scope.parent {
		for _, varName := range slices.Sorted(maps.Keys(scope.vars)) {
			result = replaceWholeWord(result, varName, "NULL::text")
		}
	}

	return result
}

func replaceWholeWord(s, old, replacement string) string {
	lower := strings.ToLower(s)
	oldLower := strings.ToLower(old)
	var result strings.Builder
	i := 0
	for i < len(lower) {
		idx := strings.Index(lower[i:], oldLower)
		if idx < 0 {
			result.WriteString(s[i:])
			break
		}
		pos := i + idx
		if pos > 0 && isWordChar(s[pos-1]) {
			result.WriteString(s[i : pos+len(old)])
			i = pos + len(old)
			continue
		}
		end := pos + len(old)
		if end < len(s) && isWordChar(s[end]) {
			result.WriteString(s[i : pos+len(old)])
			i = pos + len(old)
			continue
		}
		result.WriteString(s[i:pos])
		result.WriteString(replacement)
		i = end
	}
	return result.String()
}

func isWordChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

func extractBaseVarName(target string) string {
	for i, c := range target {
		if c == '.' || c == '[' {
			return target[:i]
		}
	}
	return target
}

// collectAccessTablesFromSQL extracts all table/column references from embedded SQL,
// adds them to the extractor's funcSourceColumns, and returns them. This works even
// when AnalyzeSelectStmt fails, because it uses the AST-level ExtractAccessTables.
func (a *plpgsqlAnalyzer) collectAccessTablesFromSQL(sql string) ColumnSet {
	accesses := a.extractor.accessTablesOf(sql)
	for resource := range accesses {
		a.extractor.funcSourceColumns[resource] = true
	}
	return accesses
}

// accessTablesOf returns the tables sql reads, or none when it cannot tell.
func (e *querySpanExtractor) accessTablesOf(sql string) ColumnSet {
	accesses := make(ColumnSet)
	accessTables, err := e.extractAccessTables(sql)
	if err != nil {
		return accesses
	}
	for _, resource := range accessTables {
		accesses[resource] = true
	}
	return accesses
}

func (a *plpgsqlAnalyzer) mergeReturnResults(nCols int) []ColumnSet {
	result := make([]ColumnSet, nCols)
	for i := range result {
		result[i] = make(ColumnSet)
	}

	for _, returnCols := range a.returnResults {
		for i := 0; i < nCols && i < len(returnCols); i++ {
			for k, v := range returnCols[i] {
				result[i][k] = v
			}
		}
	}

	return result
}

// fallbackExtractLineage is called when AnalyzeSelectStmt fails for embedded SQL
// (e.g., due to type mismatches with aggregate functions). It uses the parse tree
// to extract column references and resolves them against known tables and variables.
func (a *plpgsqlAnalyzer) fallbackExtractLineage(origSQL string, selStmt *ast.SelectStmt) []ColumnSet {
	if selStmt == nil || selStmt.TargetList == nil {
		return nil
	}

	// Collect table references from the FROM clause. The explicit-schema
	// markers are scoped to this statement — a previous fallback query's
	// qualified FROM must not pin a later query's unqualified reference.
	a.qualifiedFrom = make(map[string]bool)
	a.fromItems = make(map[string]ast.Node)
	fromTables := a.collectFromTables(selStmt)

	var results []ColumnSet
	origItems := extractSelectItemNames(origSQL)

	for i, item := range selStmt.TargetList.Items {
		rt, ok := item.(*ast.ResTarget)
		if !ok {
			results = append(results, make(ColumnSet))
			continue
		}

		colSet := make(ColumnSet)

		// Try to extract column references from the expression.
		a.extractColumnRefsFromExpr(rt.Val, fromTables, colSet)

		// If still empty, try the original item name as a variable reference.
		if len(colSet) == 0 && i < len(origItems) {
			varName := strings.TrimSpace(origItems[i])
			// Strip wrapping parentheses for subquery expressions.
			varName = strings.TrimPrefix(strings.TrimSuffix(varName, ")"), "(")
			varName = strings.TrimSpace(varName)
			if sources, ok := a.scope.get(varName); ok {
				for k, v := range sources {
					colSet[k] = v
				}
			}
		}

		results = append(results, colSet)
	}

	// Extract column refs from WHERE clause for predicate tracking.
	if selStmt.WhereClause != nil {
		whereColSet := make(ColumnSet)
		a.extractColumnRefsFromExpr(selStmt.WhereClause, fromTables, whereColSet)
		for k, v := range whereColSet {
			a.extractor.predicateColumns[k] = v
		}
	}

	return results
}

// collectFromTables collects table references from the FROM clause and maps
// table names/aliases to their schema-qualified names.
func (a *plpgsqlAnalyzer) collectFromTables(selStmt *ast.SelectStmt) map[string]ColumnResource {
	tables := make(map[string]ColumnResource)
	if selStmt.FromClause == nil {
		return tables
	}

	for _, item := range selStmt.FromClause.Items {
		a.collectTablesFromNode(item, tables)
	}
	return tables
}

// resolveFallbackRelation resolves a textual fallback resource against the
// catalog. Unqualified FROM items default to "public" in collectTablesFromNode,
// so on a miss the extractor's search path is probed the way PostgreSQL
// would resolve the name.
// pushMarkerScope opens a nested explicit-qualification scope seeded from
// the current one, so markers collected inside a subquery or CTE apply to
// that scope's resolution but never leak outward. The returned restore
// function reinstates the caller's scope; the returned map is the nested
// scope for selective exports.
func (a *plpgsqlAnalyzer) pushMarkerScope() (func(), map[string]bool) {
	outer, outerItems := a.qualifiedFrom, a.fromItems
	nested := make(map[string]bool, len(outer))
	for k, v := range outer {
		nested[k] = v
	}
	a.qualifiedFrom = nested
	a.fromItems = maps.Clone(outerItems)
	return func() { a.qualifiedFrom, a.fromItems = outer, outerItems }, nested
}

// setFromItem records the FROM item a name refers to when its columns are
// not its table's, and forgets any outer item of that name otherwise.
func (a *plpgsqlAnalyzer) setFromItem(name string, item ast.Node) {
	if a.fromItems == nil {
		a.fromItems = make(map[string]ast.Node)
	}
	if item == nil {
		delete(a.fromItems, name)
		return
	}
	a.fromItems[name] = item
}

// fromItemColumns returns the columns of a subquery, a join or a table with
// an alias column list, named as the query sees them, with what each reads.
// fromTables are the tables of the query the item is in, which a LATERAL
// subquery can reference.
func (a *plpgsqlAnalyzer) fromItemColumns(item ast.Node, fromTables map[string]ColumnResource) []QuerySpanResult {
	cache := &a.extractor.fromColumns
	if columns, ok := (*cache)[item]; ok {
		return cloneResults(columns)
	}
	if *cache == nil {
		*cache = make(map[ast.Node][]QuerySpanResult)
	}
	(*cache)[item] = nil

	var columns []QuerySpanResult
	var alias *ast.Alias
	switch v := item.(type) {
	case *ast.RangeVar:
		columns, alias = a.extractor.extractColumnsFromRangeVar(v), v.Alias
	case *ast.RangeSubselect:
		// Only a LATERAL subquery sees the other items of its FROM.
		var outer map[string]ColumnResource
		if v.Lateral {
			outer = fromTables
		}
		if sel, ok := v.Subquery.(*ast.SelectStmt); ok {
			columns = a.selectOutputColumns(sel, outer)
		}
		alias = v.Alias
	case *ast.JoinExpr:
		columns, alias = a.extractor.extractColumnsFromFromItem(v, a, fromTables), v.Alias
	default:
	}
	// pg: src/backend/parser/parse_relation.c — buildRelationAliases
	if alias != nil {
		for i, name := range astStrings(alias.Colnames) {
			if i < len(columns) {
				columns[i].Name = name
			}
		}
	}
	(*cache)[item] = columns
	return cloneResults(columns)
}

func cloneResults(results []QuerySpanResult) []QuerySpanResult {
	clones := slices.Clone(results)
	for i := range clones {
		clones[i].SourceColumns = maps.Clone(clones[i].SourceColumns)
	}
	return clones
}

// selectOutputColumns returns the output columns of a query in FROM, with
// what each reads.
func (a *plpgsqlAnalyzer) selectOutputColumns(sel *ast.SelectStmt, outerTables map[string]ColumnResource) []QuerySpanResult {
	if sel == nil {
		return nil
	}
	if sel.Op != ast.SETOP_NONE {
		// A set operation's columns are named by its first branch and read
		// what every branch puts in that position.
		columns := a.selectOutputColumns(sel.Larg, outerTables)
		for i, column := range a.selectOutputColumns(sel.Rarg, outerTables) {
			if i < len(columns) {
				maps.Copy(columns[i].SourceColumns, column.SourceColumns)
			}
		}
		return columns
	}
	if sel.TargetList == nil {
		return nil
	}
	restore, _ := a.pushMarkerScope()
	defer restore()
	a.fromItems = make(map[string]ast.Node)
	tables := a.collectFromTables(sel)
	for k, v := range outerTables {
		if _, exists := tables[k]; !exists {
			tables[k] = v
		}
	}
	var columns []QuerySpanResult
	for _, item := range sel.TargetList.Items {
		rt, ok := item.(*ast.ResTarget)
		if !ok {
			continue
		}
		if isStarTarget(rt) {
			columns = append(columns, a.starColumns(sel, starTableName(rt), tables)...)
			continue
		}
		name := rt.Name
		if name == "" {
			name = figureResTargetName(rt)
		}
		if name == "" {
			name = "?column?"
		}
		sources := make(ColumnSet)
		a.extractColumnRefsFromExpr(rt.Val, tables, sources)
		columns = append(columns, QuerySpanResult{Name: name, SourceColumns: sources})
	}
	return columns
}

// starColumns returns the columns * expands to in sel, or those of the FROM
// item named table for table.*.
func (a *plpgsqlAnalyzer) starColumns(sel *ast.SelectStmt, table string, tables map[string]ColumnResource) []QuerySpanResult {
	if sel.FromClause == nil {
		return nil
	}
	var columns []QuerySpanResult
	var expand func(node ast.Node)
	expand = func(node ast.Node) {
		if table == "" || fromItemName(node) == table {
			columns = append(columns, a.extractor.extractColumnsFromFromItem(node, a, tables)...)
			return
		}
		if join, ok := node.(*ast.JoinExpr); ok {
			expand(join.Larg)
			expand(join.Rarg)
		}
	}
	for _, item := range sel.FromClause.Items {
		expand(item)
	}
	return columns
}

// fromItemName is the name a query refers to a FROM item by.
func fromItemName(node ast.Node) string {
	switch v := node.(type) {
	case *ast.RangeVar:
		if v.Alias != nil && v.Alias.Aliasname != "" {
			return v.Alias.Aliasname
		}
		return v.Relname
	case *ast.RangeSubselect:
		if v.Alias != nil {
			return v.Alias.Aliasname
		}
	case *ast.RangeFunction:
		if v.Alias != nil {
			return v.Alias.Aliasname
		}
	case *ast.JoinExpr:
		if v.Alias != nil {
			return v.Alias.Aliasname
		}
	default:
	}
	return ""
}

// fromItemColumn adds what the column named colName of a FROM item reads.
// It reports whether the item has such a column. Names compare exactly: the
// parser has folded unquoted names already.
func (a *plpgsqlAnalyzer) fromItemColumn(item ast.Node, colName string, fromTables map[string]ColumnResource, result ColumnSet) bool {
	for _, column := range a.fromItemColumns(item, fromTables) {
		if column.Name == colName {
			maps.Copy(result, column.SourceColumns)
			return true
		}
	}
	return false
}

func (a *plpgsqlAnalyzer) resolveFallbackRelation(fromKey, schemaName, tableName string) *catalog.Relation {
	// collectTablesFromNode defaults unqualified FROM items to "public"; for
	// those (and only those — explicitly qualified items resolve exactly),
	// the configured search path decides resolution order, exactly like
	// PostgreSQL. The literal public lookup remains the last resort for
	// search paths that do not contain public.
	if schemaName == "public" && !a.qualifiedFrom[fromKey] {
		for _, pathSchema := range a.extractor.searchPath {
			if rel := a.extractor.cat.GetRelation(pathSchema, tableName); rel != nil {
				return rel
			}
		}
	}
	return a.extractor.cat.GetRelation(schemaName, tableName)
}

func (a *plpgsqlAnalyzer) collectTablesFromNode(node ast.Node, tables map[string]ColumnResource) {
	switch v := node.(type) {
	case *ast.RangeVar:
		schema := "public"
		if v.Schemaname != "" {
			schema = v.Schemaname
		}
		resource := ColumnResource{
			Database: a.extractor.defaultDatabase,
			Schema:   schema,
			Table:    v.Relname,
		}
		// The newest FROM item owns its key's qualification marker in this
		// scope: set it when explicitly qualified, clear any inherited one
		// when not — an inner unqualified FROM must not be pinned by a
		// same-named outer explicit FROM.
		if a.qualifiedFrom == nil {
			a.qualifiedFrom = make(map[string]bool)
		}
		tables[strings.ToLower(v.Relname)] = resource
		if v.Schemaname != "" {
			a.qualifiedFrom[strings.ToLower(v.Relname)] = true
		} else {
			delete(a.qualifiedFrom, strings.ToLower(v.Relname))
		}
		a.setFromItem(strings.ToLower(v.Relname), nil)
		if v.Alias != nil && v.Alias.Aliasname != "" {
			tables[strings.ToLower(v.Alias.Aliasname)] = resource
			if v.Schemaname != "" {
				a.qualifiedFrom[strings.ToLower(v.Alias.Aliasname)] = true
			} else {
				delete(a.qualifiedFrom, strings.ToLower(v.Alias.Aliasname))
			}
			var item ast.Node
			if v.Alias.Colnames != nil {
				item = v
			}
			a.setFromItem(strings.ToLower(v.Alias.Aliasname), item)
		}
	case *ast.RangeSubselect:
		// For subquery aliases like (SELECT * FROM t) result, collect inner tables
		// and map them under the alias so column refs can be resolved.
		if v.Subquery != nil {
			if subSel, ok := v.Subquery.(*ast.SelectStmt); ok {
				restore, nested := a.pushMarkerScope()
				innerTables := a.collectFromTables(subSel)
				restore()
				if a.qualifiedFrom == nil {
					a.qualifiedFrom = make(map[string]bool)
				}
				// Add inner tables to the outer fromTables. Every key this
				// branch writes takes ownership of its qualification marker:
				// set from the borrowed resource's inner marker, cleared
				// otherwise — never inherited from the outer scope.
				for k, res := range innerTables {
					if _, exists := tables[k]; !exists {
						tables[k] = res
						if nested[k] {
							a.qualifiedFrom[k] = true
						} else {
							delete(a.qualifiedFrom, k)
						}
					}
				}
				// If the subquery has an alias, also map all inner tables'
				// columns as accessible through the alias.
				if v.Alias != nil && v.Alias.Aliasname != "" {
					alias := strings.ToLower(v.Alias.Aliasname)
					a.setFromItem(alias, v)
					for k, res := range innerTables {
						tables[alias] = res
						if nested[k] {
							a.qualifiedFrom[alias] = true
						} else {
							delete(a.qualifiedFrom, alias)
						}
						break
					}
				}
			}
		}
	case *ast.JoinExpr:
		if v.Larg != nil {
			a.collectTablesFromNode(v.Larg, tables)
		}
		if v.Rarg != nil {
			a.collectTablesFromNode(v.Rarg, tables)
		}
		if v.Alias != nil && v.Alias.Aliasname != "" {
			a.setFromItem(strings.ToLower(v.Alias.Aliasname), v)
		}
	default:
	}
}

// extractColumnRefsFromExpr recursively walks an AST expression and resolves
// column references against known tables.
func (a *plpgsqlAnalyzer) extractColumnRefsFromExpr(node ast.Node, fromTables map[string]ColumnResource, result ColumnSet) {
	if node == nil {
		return
	}

	switch v := node.(type) {
	case *ast.ColumnRef:
		a.resolveColumnRef(v, fromTables, result)
	case *ast.FuncCall:
		if v.Args != nil {
			for _, arg := range v.Args.Items {
				a.extractColumnRefsFromExpr(arg, fromTables, result)
			}
		}
		// FILTER (WHERE ...) clause on aggregates → predicate columns.
		if v.AggFilter != nil {
			filterColSet := make(ColumnSet)
			a.extractColumnRefsFromExpr(v.AggFilter, fromTables, filterColSet)
			for k, val := range filterColSet {
				a.extractor.predicateColumns[k] = val
			}
		}
	case *ast.A_Expr:
		a.extractColumnRefsFromExpr(v.Lexpr, fromTables, result)
		a.extractColumnRefsFromExpr(v.Rexpr, fromTables, result)
	case *ast.SubLink:
		// The left operand of IN, ANY, ALL or a row comparison.
		a.extractColumnRefsFromExpr(v.Testexpr, fromTables, result)
		if subSel, ok := v.Subselect.(*ast.SelectStmt); ok {
			a.extractSubqueryRefs(subSel, fromTables, result)
		}
	case *ast.TypeCast:
		a.extractColumnRefsFromExpr(v.Arg, fromTables, result)
	case *ast.BoolExpr:
		if v.Args != nil {
			for _, arg := range v.Args.Items {
				a.extractColumnRefsFromExpr(arg, fromTables, result)
			}
		}
	case *ast.CaseExpr:
		a.extractColumnRefsFromExpr(v.Arg, fromTables, result)
		if v.Args != nil {
			for _, w := range v.Args.Items {
				a.extractColumnRefsFromExpr(w, fromTables, result)
			}
		}
		a.extractColumnRefsFromExpr(v.Defresult, fromTables, result)
	case *ast.CaseWhen:
		a.extractColumnRefsFromExpr(v.Expr, fromTables, result)
		a.extractColumnRefsFromExpr(v.Result, fromTables, result)
	case *ast.CoalesceExpr:
		if v.Args != nil {
			for _, arg := range v.Args.Items {
				a.extractColumnRefsFromExpr(arg, fromTables, result)
			}
		}
	case *ast.NullTest:
		a.extractColumnRefsFromExpr(v.Arg, fromTables, result)
	case *ast.BooleanTest:
		a.extractColumnRefsFromExpr(v.Arg, fromTables, result)
	case *ast.A_Indirection:
		// Array subscript or field selection — the source is the base expression.
		a.extractColumnRefsFromExpr(v.Arg, fromTables, result)
		if v.Indirection != nil {
			for _, ind := range v.Indirection.Items {
				a.extractColumnRefsFromExpr(ind, fromTables, result)
			}
		}
	case *ast.JsonObjectConstructor:
		if v.Exprs != nil {
			for _, item := range v.Exprs.Items {
				if kv, ok := item.(*ast.JsonKeyValue); ok {
					a.extractColumnRefsFromExpr(kv.Key, fromTables, result)
					if kv.Value != nil {
						a.extractColumnRefsFromExpr(kv.Value.RawExpr, fromTables, result)
					}
				}
			}
		}
	case *ast.JsonArrayConstructor:
		if v.Exprs != nil {
			for _, item := range v.Exprs.Items {
				if jve, ok := item.(*ast.JsonValueExpr); ok {
					a.extractColumnRefsFromExpr(jve.RawExpr, fromTables, result)
				}
			}
		}
	case *ast.JsonValueExpr:
		a.extractColumnRefsFromExpr(v.RawExpr, fromTables, result)
	case *ast.List:
		for _, item := range v.Items {
			a.extractColumnRefsFromExpr(item, fromTables, result)
		}
	default:
		// Any other expression (GREATEST, ROW, ARRAY, a collation, ...)
		// reads what its children read.
		ast.Inspect(node, func(child ast.Node) bool {
			if child == node {
				return true
			}
			a.extractColumnRefsFromExpr(child, fromTables, result)
			return false
		})
	}
}

// extractSubqueryRefs adds what the outputs of a subquery in an expression
// read, from every branch of a set operation, and collects its WHERE
// columns as predicates.
func (a *plpgsqlAnalyzer) extractSubqueryRefs(subSel *ast.SelectStmt, fromTables map[string]ColumnResource, result ColumnSet) {
	if subSel == nil {
		return
	}
	if subSel.Op != ast.SETOP_NONE {
		a.extractSubqueryRefs(subSel.Larg, fromTables, result)
		a.extractSubqueryRefs(subSel.Rarg, fromTables, result)
		return
	}
	// The subquery's qualification markers apply only within its own
	// resolution scope.
	restore, _ := a.pushMarkerScope()
	defer restore()
	subTables := a.collectFromTables(subSel)
	for k, val := range fromTables {
		if _, exists := subTables[k]; !exists {
			subTables[k] = val
		}
	}
	if subSel.TargetList != nil {
		for _, item := range subSel.TargetList.Items {
			if rt, ok := item.(*ast.ResTarget); ok {
				a.extractColumnRefsFromExpr(rt.Val, subTables, result)
			}
		}
	}
	if subSel.WhereClause != nil {
		whereColSet := make(ColumnSet)
		a.extractColumnRefsFromExpr(subSel.WhereClause, subTables, whereColSet)
		for k, val := range whereColSet {
			a.extractor.predicateColumns[k] = val
		}
	}
}

// resolveColumnRef resolves a ColumnRef to a source column using known tables.
// If the table is a CTE, it traces through the CTE to find real source columns.
func (a *plpgsqlAnalyzer) resolveColumnRef(cr *ast.ColumnRef, fromTables map[string]ColumnResource, result ColumnSet) {
	if cr == nil || cr.Fields == nil || len(cr.Fields.Items) == 0 {
		return
	}

	var parts []string
	for _, f := range cr.Fields.Items {
		if s, ok := f.(*ast.String); ok {
			parts = append(parts, s.Str)
		}
	}

	switch len(parts) {
	case 1:
		colName := parts[0]
		if _, ok := a.scope.get(colName); ok {
			return
		}
		// Without scopes the table a name binds to is unknown, so it reads
		// the column of every FROM table that has one.
		found := false
		for _, fromKey := range slices.Sorted(maps.Keys(fromTables)) {
			resource := fromTables[fromKey]
			rel := a.resolveFallbackRelation(fromKey, resource.Schema, resource.Table)
			// Standalone composite types are not FROM-able.
			if rel != nil && rel.RelKind != 'c' && relationHasColumn(rel, colName) {
				// Emit the schema the relation actually resolved to — an
				// unqualified item may have resolved through the search
				// path, not the public default.
				emitSchema, emitTable := resource.Schema, resource.Table
				if rel.Schema != nil {
					emitSchema, emitTable = rel.Schema.Name, rel.Name
				}
				result[ColumnResource{
					Database: resource.Database,
					Schema:   emitSchema,
					Table:    emitTable,
					Column:   colName,
				}] = true
				found = true
			}
		}
		if found {
			return
		}
		// An unqualified FROM name may be a CTE, whose columns are its
		// query's outputs.
		for _, fromKey := range slices.Sorted(maps.Keys(fromTables)) {
			resource := fromTables[fromKey]
			cteSel, ok := a.cteMap[strings.ToLower(resource.Table)]
			if !ok || a.qualifiedFrom[fromKey] {
				continue
			}
			sources := make(ColumnSet)
			a.resolveThroughCTE(cteSel, colName, fromTables, sources)
			if len(sources) > 0 {
				maps.Copy(result, sources)
				found = true
			}
		}
		if found {
			return
		}
		for _, name := range slices.Sorted(maps.Keys(a.fromItems)) {
			if a.fromItemColumn(a.fromItems[name], colName, fromTables, result) {
				found = true
			}
		}
		if found {
			return
		}
		// Check if it's a whole-row reference (table alias used as value, e.g., to_jsonb(t1)).
		if resource, ok := fromTables[strings.ToLower(colName)]; ok {
			rel := a.resolveFallbackRelation(strings.ToLower(colName), resource.Schema, resource.Table)
			if rel != nil && rel.RelKind != 'c' {
				emitSchema, emitTable := resource.Schema, resource.Table
				if rel.Schema != nil {
					emitSchema, emitTable = rel.Schema.Name, rel.Name
				}
				for _, col := range rel.Columns {
					result[ColumnResource{
						Database: resource.Database,
						Schema:   emitSchema,
						Table:    emitTable,
						Column:   col.Name,
					}] = true
				}
				return
			}
		}
	case 2:
		tableName := strings.ToLower(parts[0])
		colName := parts[1]
		if item, ok := a.fromItems[tableName]; ok {
			a.fromItemColumn(item, colName, fromTables, result)
			return
		}
		if resource, ok := fromTables[tableName]; ok {
			// Check if this table is actually a CTE — trace through to real sources.
			// The CTE name is the original table name (resource.Table), not the alias.
			if a.cteMap != nil {
				if cteSel, isCTE := a.cteMap[strings.ToLower(resource.Table)]; isCTE {
					a.resolveThroughCTE(cteSel, colName, fromTables, result)
					return
				}
			}
			// A known standalone composite type is not FROM-able; suppress
			// lineage for it (unknown relations keep the textual fallback);
			// a resolved relation supplies the schema it actually lives in.
			rel := a.resolveFallbackRelation(tableName, resource.Schema, resource.Table)
			if rel != nil && rel.RelKind == 'c' {
				return
			}
			emitSchema, emitTable := resource.Schema, resource.Table
			if rel != nil && rel.Schema != nil {
				emitSchema, emitTable = rel.Schema.Name, rel.Name
			}
			result[ColumnResource{
				Database: resource.Database,
				Schema:   emitSchema,
				Table:    emitTable,
				Column:   colName,
			}] = true
		}
	default:
	}
}

// resolveThroughCTE traces a column reference through a CTE to find real source columns.
func (a *plpgsqlAnalyzer) resolveThroughCTE(cteSel *ast.SelectStmt, colName string, _ map[string]ColumnResource, result ColumnSet) {
	if cteSel == nil || cteSel.TargetList == nil {
		return
	}
	key := cteResolveKey{
		sel: cteSel,
		col: strings.ToLower(colName),
	}
	if a.resolvingCTE == nil {
		a.resolvingCTE = make(map[cteResolveKey]bool)
	}
	if a.resolvingCTE[key] {
		return
	}
	a.resolvingCTE[key] = true
	defer delete(a.resolvingCTE, key)

	// Build FROM tables for the CTE's query; its qualification markers are
	// scoped to this resolution.
	restoreMarkers, _ := a.pushMarkerScope()
	defer restoreMarkers()
	cteTables := a.collectFromTables(cteSel)
	// Include parent CTEs for nested CTE references.
	if a.cteMap != nil {
		for name := range a.cteMap {
			if _, exists := cteTables[name]; !exists {
				cteTables[name] = ColumnResource{
					Database: a.extractor.defaultDatabase,
					Schema:   "public",
					Table:    name,
				}
			}
		}
	}

	// Find the target entry matching the column name.
	for _, item := range cteSel.TargetList.Items {
		rt, ok := item.(*ast.ResTarget)
		if !ok {
			continue
		}
		// Match by alias first, then by expression name.
		rtName := rt.Name
		if rtName == "" {
			rtName = figureResTargetName(rt)
		}
		if !strings.EqualFold(rtName, colName) {
			continue
		}
		// Found matching column — extract source refs from its expression.
		a.extractColumnRefsFromExpr(rt.Val, cteTables, result)

		// Also collect predicate columns from the CTE's WHERE clause.
		if cteSel.WhereClause != nil {
			whereColSet := make(ColumnSet)
			a.extractColumnRefsFromExpr(cteSel.WhereClause, cteTables, whereColSet)
			for k, v := range whereColSet {
				a.extractor.predicateColumns[k] = v
			}
		}
		return
	}
}

// collectCTEDefinitions extracts CTE names and their SelectStmt definitions from a query.
func collectCTEDefinitions(selStmt *ast.SelectStmt) map[string]*ast.SelectStmt {
	if selStmt == nil || selStmt.WithClause == nil {
		return nil
	}

	cteMap := make(map[string]*ast.SelectStmt)
	for _, item := range selStmt.WithClause.Ctes.Items {
		cte, ok := item.(*ast.CommonTableExpr)
		if !ok {
			continue
		}
		if sel, ok := cte.Ctequery.(*ast.SelectStmt); ok {
			cteMap[strings.ToLower(cte.Ctename)] = sel
		}
	}
	return cteMap
}

// relationHasColumn checks if a catalog relation has a column with the given name.
func relationHasColumn(rel *catalog.Relation, colName string) bool {
	lowerName := strings.ToLower(colName)
	for _, col := range rel.Columns {
		if strings.ToLower(col.Name) == lowerName {
			return true
		}
	}
	return false
}
