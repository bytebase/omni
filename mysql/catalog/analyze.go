package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// AnalyzeSelectStmt performs semantic analysis on a parsed SELECT statement,
// returning a resolved Query IR.
func (c *Catalog) AnalyzeSelectStmt(stmt *nodes.SelectStmt) (*Query, error) {
	q := &Query{
		CommandType: CmdSelect,
		JoinTree:    &JoinTreeQ{},
	}

	scope := newScope()

	// Step 1: Analyze FROM clause → populate RangeTable and scope.
	if err := analyzeFromClause(c, stmt.From, q, scope); err != nil {
		return nil, err
	}

	// Step 2: Analyze target list (SELECT expressions).
	if err := analyzeTargetList(stmt.TargetList, q, scope); err != nil {
		return nil, err
	}

	// Step 3: Analyze WHERE clause.
	if stmt.Where != nil {
		analyzed, err := analyzeExpr(stmt.Where, scope)
		if err != nil {
			return nil, err
		}
		q.JoinTree.Quals = analyzed
	}

	return q, nil
}

// analyzeFromClause processes the FROM clause, populating the query's
// RangeTable, JoinTree.FromList, and the scope for column resolution.
func analyzeFromClause(c *Catalog, from []nodes.TableExpr, q *Query, scope *analyzerScope) error {
	for _, te := range from {
		switch ref := te.(type) {
		case *nodes.TableRef:
			rte, cols, err := analyzeTableRef(c, ref)
			if err != nil {
				return err
			}
			idx := len(q.RangeTable)
			q.RangeTable = append(q.RangeTable, rte)
			q.JoinTree.FromList = append(q.JoinTree.FromList, &RangeTableRefQ{RTIndex: idx})

			// Determine the scope name: alias if present, else table name.
			scopeName := rte.ERef
			scope.add(scopeName, idx, cols)
		default:
			return fmt.Errorf("unsupported FROM clause element: %T", te)
		}
	}
	return nil
}

// analyzeTableRef resolves a table reference from the FROM clause against
// the catalog, returning the RTE and the column list.
func analyzeTableRef(c *Catalog, ref *nodes.TableRef) (*RangeTableEntryQ, []*Column, error) {
	dbName := ref.Schema
	if dbName == "" {
		dbName = c.CurrentDatabase()
	}
	if dbName == "" {
		return nil, nil, errNoDatabaseSelected()
	}

	db := c.GetDatabase(dbName)
	if db == nil {
		return nil, nil, errUnknownDatabase(dbName)
	}

	// Check for a table first, then a view.
	tbl := db.GetTable(ref.Name)
	if tbl != nil {
		eref := ref.Name
		if ref.Alias != "" {
			eref = ref.Alias
		}
		colNames := make([]string, len(tbl.Columns))
		for i, col := range tbl.Columns {
			colNames[i] = col.Name
		}
		rte := &RangeTableEntryQ{
			Kind:      RTERelation,
			DBName:    db.Name,
			TableName: tbl.Name,
			Alias:     ref.Alias,
			ERef:      eref,
			ColNames:  colNames,
		}
		return rte, tbl.Columns, nil
	}

	// Check views.
	view := db.Views[toLower(ref.Name)]
	if view != nil {
		eref := ref.Name
		if ref.Alias != "" {
			eref = ref.Alias
		}
		// Build stub columns from view column names.
		cols := make([]*Column, len(view.Columns))
		colNames := make([]string, len(view.Columns))
		for i, name := range view.Columns {
			cols[i] = &Column{Position: i + 1, Name: name}
			colNames[i] = name
		}
		rte := &RangeTableEntryQ{
			Kind:          RTERelation,
			DBName:        db.Name,
			TableName:     view.Name,
			Alias:         ref.Alias,
			ERef:          eref,
			ColNames:      colNames,
			IsView:        true,
			ViewAlgorithm: viewAlgorithmFromString(view.Algorithm),
		}
		return rte, cols, nil
	}

	return nil, nil, errNoSuchTable(dbName, ref.Name)
}

// viewAlgorithmFromString converts a string algorithm value to ViewAlgorithm.
func viewAlgorithmFromString(s string) ViewAlgorithm {
	switch strings.ToUpper(s) {
	case "MERGE":
		return ViewAlgMerge
	case "TEMPTABLE":
		return ViewAlgTemptable
	case "UNDEFINED", "":
		return ViewAlgUndefined
	default:
		return ViewAlgUndefined
	}
}
