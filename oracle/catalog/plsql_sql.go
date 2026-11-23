package catalog

import (
	"errors"
	"strings"

	"github.com/bytebase/omni/oracle/ast"
)

// plsqlSQLCodes are the errors of a static SQL statement the checker
// reports, as USER_ERRORS lists them ("PL/SQL: ORA-00942: ..."). The query
// analyzer finds others, which the program suite does not pin; the checker
// takes those statements as compiling.
var plsqlSQLCodes = map[string]bool{"00904": true, "00942": true, "00913": true, "00947": true}

// dictionaryPrefixes start the names of Oracle's dictionary views, which a
// unit reads through public synonyms the catalog does not hold.
var dictionaryPrefixes = []string{"USER_", "ALL_", "DBA_", "CDB_", "V$", "GV$", "DICT", "NLS_", "ROLE_", "SESSION_"}

// sql checks a static SQL statement in a unit, and returns the columns of a
// query's result, or nil.
func (c *plChecker) sql(n ast.StmtNode) map[string]bool {
	if n == nil || c.err != nil {
		return nil
	}
	if c.readsUnknown(n) {
		// The statement still reads the objects the catalog knows.
		c.readColumns(n)
		return nil
	}
	a := c.analyzer()
	var fields map[string]bool
	var err error
	switch s := n.(type) {
	case *ast.SelectStmt:
		var res *queryResult
		res, err = a.query(s, nil, nil, false)
		if err == nil {
			fields = map[string]bool{}
			for _, col := range res.cols {
				fields[col.name] = true
			}
		}
		c.targets(s.IntoVars)
	case *ast.InsertStmt:
		err = c.insert(a, s)
	case *ast.UpdateStmt:
		err = c.update(a, s)
	case *ast.DeleteStmt:
		err = c.delete(a, s)
	case *ast.MergeStmt:
		err = c.merge(a, s)
	}
	var refusal *Error
	if errors.As(err, &refusal) && plsqlSQLCodes[refusal.Code] {
		c.fail("ORA-"+refusal.Code, "%s", refusal.Message)
	}
	c.readColumns(n)
	return fields
}

// readColumns records the columns a statement names of each table it
// reads, and a table's every column for its * or an INSERT without a
// column list, so that a change to one column leaves only its readers to
// be compiled again.
func (c *plChecker) readColumns(n ast.Node) {
	tables := map[string][]string{}
	names := map[string]bool{}
	star := false
	ast.Inspect(n, func(x ast.Node) bool {
		switch x := x.(type) {
		case *ast.TableRef:
			if x.Name != nil {
				tables[x.Name.Name] = nil
			}
		case *ast.ColumnRef:
			if x.Column == "*" {
				star = true
			} else {
				names[x.Column] = true
			}
		}
		return true
	})
	for _, name := range dmlTables(n) {
		tables[name] = nil
	}
	for t := range tables {
		target := c.e.resolveSynonymTarget(t, "")
		var cols []string
		if tb := c.e.tables[target]; tb != nil {
			for _, col := range tb.Columns {
				cols = append(cols, col.Name)
			}
		} else if v := c.e.views[target]; v != nil {
			for _, col := range v.Columns {
				cols = append(cols, col.Name)
			}
		} else {
			continue
		}
		c.read(t)
		c.read(target)
		if star {
			c.read(target + ".*")
		}
		for _, col := range cols {
			if names[col] {
				c.read(target + "." + col)
			}
		}
	}
}

// dmlTables lists the tables a DML statement writes.
func dmlTables(n ast.Node) []string {
	switch s := n.(type) {
	case *ast.InsertStmt:
		if s.Table != nil {
			return []string{s.Table.Name}
		}
	case *ast.UpdateStmt:
		if s.Table != nil {
			return []string{s.Table.Name}
		}
	case *ast.DeleteStmt:
		if s.Table != nil {
			return []string{s.Table.Name}
		}
	case *ast.MergeStmt:
		if s.Target != nil {
			return []string{s.Target.Name}
		}
	}
	return nil
}

// analyzer is a query analyzer whose names that resolve to no column may
// be the unit's own.
func (c *plChecker) analyzer() *queryAnalyzer {
	return &queryAnalyzer{
		e:          c.e,
		plsqlFuncs: map[string]qtype{},
		owner:      map[*ast.ColumnRef]*qfrom{},
		memo:       map[exprMemoKey]exprMemo{},
		queryMemo:  map[queryMemoKey]queryMemo{},
		blockMemo:  map[blockMemoKey]queryMemo{},
		hashMemo:   map[hashMemoKey]uint64{},
		plsqlName:  c.sqlName,
		plsqlCall:  c.sqlCall,
		plsqlBind: func(b *ast.BindVariable) error {
			// A bad row field is the unit's error, not the statement's.
			c.bind(b)
			return nil
		},
	}
}

// sqlName reports whether a name a SQL statement reads, which names no
// column, is one the unit or the schema declares: a variable, a parameter,
// a record's field, a package's item, a sequence's pseudocolumn, or a
// function called without arguments.
func (c *plChecker) sqlName(ref *ast.ColumnRef) bool {
	switch {
	case ref.Table == "":
		if c.lookup(ref.Column) != nil || plsqlValues[ref.Column] {
			return true
		}
		target := c.e.resolveSynonymTarget(ref.Column, "")
		if c.e.functions[target] != nil {
			c.readVia(ref.Column, target)
			return true
		}
		return false
	case ref.Schema == "":
		if sym := c.lookup(ref.Table); sym != nil {
			return true
		}
		target := c.e.resolveSynonymTarget(ref.Table, "")
		if c.e.sequences[target] != nil && (ref.Column == "NEXTVAL" || ref.Column == "CURRVAL") {
			c.readVia(ref.Table, target)
			return true
		}
		if c.e.packages[target] != nil {
			c.readVia(ref.Table, target)
			c.read(target + "." + ref.Column)
			return true
		}
		// A first part the catalog cannot place: another schema's
		// object, or a public synonym's.
		return !c.knownObject(ref.Table)
	}
	return true
}

// sqlCall checks a function a SQL statement calls: a built-in, a schema
// or package function, or a type's constructor. A unit's own local
// function cannot be called from SQL, nor can a name the schema does not
// hold (ORA-00904).
func (c *plChecker) sqlCall(n *ast.FuncCallExpr) error {
	if n.FuncName == nil || n.FuncName.Name == "" {
		return nil
	}
	if n.FuncName.Schema != "" {
		if target := c.e.resolveSynonymTarget(n.FuncName.Schema, ""); c.e.packages[target] != nil {
			c.readVia(n.FuncName.Schema, target)
			c.read(target + "." + n.FuncName.Name)
		}
		return nil
	}
	name := n.FuncName.Name
	if plsqlFunctions[name] || n.Over != nil || sqlQuantifiers[name] || rowPatternFunctions[name] {
		return nil
	}
	target := c.e.resolveSynonymTarget(name, "")
	if c.e.functions[target] != nil || c.e.hasType(target) {
		c.readVia(name, target)
		return nil
	}
	if c.lookup(name) != nil || c.knownObject(name) {
		return nil
	}
	return oraError("00904", "%q: invalid identifier", name)
}

// readsUnknown reports a statement that reads a relation the catalog
// cannot know: a dictionary view, a database link's table, or a
// materialized view whose columns it does not keep.
func (c *plChecker) readsUnknown(n ast.Node) bool {
	unknown := false
	ast.Inspect(n, func(x ast.Node) bool {
		ref, ok := x.(*ast.TableRef)
		if !ok || ref.Name == nil {
			return !unknown
		}
		if ref.Dblink != "" || ref.Name.DBLink != "" || ref.Name.Schema != "" && c.e.isForeignSchema(ref.Name.Schema) {
			unknown = true
		}
		if target := c.e.resolveSynonymTarget(ref.Name.Name, ""); c.e.materializedViews[target] != nil && c.e.tables[target] == nil {
			// A materialized view whose columns the catalog does not keep.
			c.readVia(ref.Name.Name, target)
			unknown = true
		}
		if ref.Name.Schema == "" && !c.e.objectExistsInNamespace(ref.Name.Name) {
			for _, p := range dictionaryPrefixes {
				if strings.HasPrefix(ref.Name.Name, p) {
					unknown = true
				}
			}
		}
		return !unknown
	})
	return unknown
}

// target resolves the table a DML statement writes, through a synonym, and
// returns the scope its columns are read in.
func (c *plChecker) dmlTarget(a *queryAnalyzer, name *ast.ObjectName, alias *ast.Alias) (*qfrom, error) {
	if name == nil {
		return &qfrom{}, nil
	}
	f, err := a.tableRef(&ast.TableRef{Name: name, Alias: alias, Loc: ast.Loc{Start: -1, End: -1}}, &qscope{})
	if err == nil {
		c.read(name.Name)
	}
	return f, err
}

// relColumns lists a FROM item's columns by name.
func relColumns(f *qfrom) map[string]bool {
	cols := map[string]bool{}
	for _, r := range f.rels {
		for _, col := range r.cols {
			cols[col.name] = true
		}
	}
	return cols
}

func (c *plChecker) insert(a *queryAnalyzer, s *ast.InsertStmt) error {
	if s.Table == nil {
		return nil
	}
	f, err := c.dmlTarget(a, s.Table, s.Alias)
	if err != nil {
		return err
	}
	cols := relColumns(f)
	listed := 0
	for _, item := range listItems(s.Columns) {
		ref, ok := item.(*ast.ColumnRef)
		if !ok {
			continue
		}
		listed++
		if !cols[ref.Column] {
			return oraError("00904", "%q: invalid identifier", ref.Column)
		}
	}
	if listed == 0 {
		listed = visibleColumns(f)
		c.read(c.e.resolveSynonymTarget(s.Table.Name, "") + ".*")
	}
	sc := &qscope{}
	given := -1
	switch {
	case s.Values != nil:
		given = 0
		for _, item := range listItems(s.Values) {
			x, ok := item.(ast.ExprNode)
			if !ok {
				continue
			}
			given++
			if c.isDefault(x) {
				continue
			}
			if _, _, err := a.expr(x, sc, clauseSelect); err != nil {
				return err
			}
		}
	case s.Select != nil:
		res, err := a.query(s.Select, nil, nil, false)
		if err != nil {
			return err
		}
		given = len(res.cols)
	}
	if s.ValuesRecord != nil {
		// VALUES rec: a record the unit declares.
		c.expr(s.ValuesRecord)
	}
	switch {
	case given < 0, s.ValuesRecord != nil:
	case given > listed:
		return oraError("00913", "too many values")
	case given < listed:
		return oraError("00947", "not enough values")
	}
	return c.returning(a, s.Returning, f)
}

// returning checks a DML statement's RETURNING clause, which the parser
// keeps as the returned expressions followed by the INTO targets: as many
// of each, or Oracle refuses the statement. The expressions read the
// target's columns; the targets are the unit's own.
func (c *plChecker) returning(a *queryAnalyzer, list *ast.List, f *qfrom) error {
	items := listItems(list)
	if len(items) == 0 || len(items)%2 != 0 {
		return nil
	}
	half := len(items) / 2
	for _, item := range items[:half] {
		if x, ok := item.(ast.ExprNode); ok {
			if _, _, err := a.expr(x, &qscope{from: f}, clauseSelect); err != nil {
				return err
			}
		}
	}
	for _, item := range items[half:] {
		if x, ok := item.(ast.ExprNode); ok {
			c.target(x)
		}
	}
	return nil
}

func visibleColumns(f *qfrom) int {
	n := 0
	for _, r := range f.rels {
		for _, col := range r.cols {
			if !col.hidden {
				n++
			}
		}
	}
	return n
}

func (c *plChecker) setClauses(a *queryAnalyzer, list *ast.List, sc *qscope) error {
	cols := relColumns(sc.from)
	for _, item := range listItems(list) {
		set, ok := item.(*ast.SetClause)
		if !ok {
			continue
		}
		if set.Column != nil && !cols[set.Column.Column] {
			return oraError("00904", "%q: invalid identifier", set.Column.Column)
		}
		targets := 0
		for _, item := range listItems(set.Columns) {
			if ref, ok := item.(*ast.ColumnRef); ok {
				targets++
				if !cols[ref.Column] {
					return oraError("00904", "%q: invalid identifier", ref.Column)
				}
			}
		}
		switch {
		case set.Value == nil, c.isDefault(set.Value):
		case targets > 0:
			// SET (a, b) = (SELECT ...): the query gives a row.
			sq, ok := set.Value.(*ast.SubqueryExpr)
			if !ok {
				continue
			}
			sel, ok := sq.Subquery.(*ast.SelectStmt)
			if !ok {
				continue
			}
			res, err := a.query(sel, sc, nil, false)
			if err != nil {
				return err
			}
			if len(res.cols) > targets {
				return oraError("00913", "too many values")
			}
			if len(res.cols) < targets {
				return oraError("00947", "not enough values")
			}
		default:
			if _, _, err := a.expr(set.Value, sc, clauseSelect); err != nil {
				return err
			}
		}
	}
	return nil
}

// isDefault reports DML's DEFAULT keyword, which the parser keeps as a
// name: VALUES (DEFAULT), SET a = DEFAULT. A quoted "DEFAULT" is a name.
func (c *plChecker) isDefault(x ast.ExprNode) bool {
	ref, ok := x.(*ast.ColumnRef)
	return ok && ref.Table == "" && ref.Column == "DEFAULT" && !strings.Contains(c.e.nodeText(ref), `"`)
}

func (c *plChecker) update(a *queryAnalyzer, s *ast.UpdateStmt) error {
	if s.Table == nil || s.FromClause != nil {
		return nil
	}
	f, err := c.dmlTarget(a, s.Table, s.Alias)
	if err != nil {
		return err
	}
	sc := &qscope{from: f}
	if err := c.setClauses(a, s.SetClauses, sc); err != nil {
		return err
	}
	if _, _, err := a.expr(s.WhereClause, sc, clauseWhere); err != nil {
		return err
	}
	return c.returning(a, s.Returning, f)
}

func (c *plChecker) delete(a *queryAnalyzer, s *ast.DeleteStmt) error {
	if s.Table == nil || s.FromClause != nil {
		return nil
	}
	f, err := c.dmlTarget(a, s.Table, s.Alias)
	if err != nil {
		return err
	}
	if _, _, err := a.expr(s.WhereClause, &qscope{from: f}, clauseWhere); err != nil {
		return err
	}
	return c.returning(a, s.Returning, f)
}

func (c *plChecker) merge(a *queryAnalyzer, s *ast.MergeStmt) error {
	target, err := c.dmlTarget(a, s.Target, s.TargetAlias)
	if err != nil {
		return err
	}
	src, err := a.fromItem(s.Source, &qscope{})
	if err != nil {
		return err
	}
	if s.SourceAlias != nil && len(src.rels) == 1 {
		src.rels[0].name = s.SourceAlias.Name
	}
	both := &qfrom{rels: append(append([]*qrel{}, target.rels...), src.rels...)}
	both.items = append(append([]fromCol{}, target.items...), src.items...)
	sc := &qscope{from: both}
	if _, _, err := a.expr(s.On, sc, clauseWhere); err != nil {
		return err
	}
	for _, item := range listItems(s.Clauses) {
		m, ok := item.(*ast.MergeClause)
		if !ok {
			continue
		}
		if err := c.setClauses(a, m.UpdateSet, &qscope{from: target, parent: &qscope{from: src}}); err != nil {
			return err
		}
	}
	return nil
}

// sqlQuantifiers are the words the parser keeps as a call's name for a
// quantified comparison's list: x = ANY (1, 2).
var sqlQuantifiers = map[string]bool{"ALL": true, "ANY": true, "SOME": true}

// rowPatternFunctions are MATCH_RECOGNIZE's, which the query analyzer reads
// in a pattern's MEASURES and DEFINE.
var rowPatternFunctions = map[string]bool{"MATCH_NUMBER": true, "CLASSIFIER": true, "PREV": true, "NEXT": true, "FIRST": true, "LAST": true}
