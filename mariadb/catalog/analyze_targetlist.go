package catalog

import (
	"fmt"
	"strconv"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// analyzeTargetList processes the SELECT list, populating q.TargetList.
func analyzeTargetList(c *Catalog, targetList []nodes.ExprNode, q *Query, scope *analyzerScope) error {
	resNo := 1
	for _, item := range targetList {
		entries, err := analyzeTargetEntry(c, item, q, scope, &resNo)
		if err != nil {
			return err
		}
		q.TargetList = append(q.TargetList, entries...)
	}
	return nil
}

// analyzeTargetEntry processes one item from the SELECT list. It may return
// multiple TargetEntryQ values when expanding star expressions.
func analyzeTargetEntry(c *Catalog, item nodes.ExprNode, q *Query, scope *analyzerScope, resNo *int) ([]*TargetEntryQ, error) {
	switch n := item.(type) {
	case *nodes.ResTarget:
		// Aliased expression: SELECT expr AS alias
		analyzed, err := analyzeExpr(c, n.Val, scope)
		if err != nil {
			return nil, err
		}
		if row, ok := unwrapParenExpr(n.Val).(*nodes.RowExpr); ok && len(row.Items) != 1 {
			return nil, errOperandColumns(1)
		}
		te := &TargetEntryQ{
			Expr:    analyzed,
			ResNo:   *resNo,
			ResName: n.Name,
		}
		fillProvenance(te, q)
		*resNo++
		return []*TargetEntryQ{te}, nil

	case *nodes.StarExpr:
		// SELECT *
		return expandStar("", "", q, scope, resNo)

	case *nodes.ColumnRef:
		if n.Star {
			// SELECT [db.]t.*
			return expandStar(n.Schema, n.Table, q, scope, resNo)
		}
		// Bare column reference: SELECT col
		analyzed, err := analyzeExpr(c, item, scope)
		if err != nil {
			return nil, err
		}
		te := &TargetEntryQ{
			Expr:    analyzed,
			ResNo:   *resNo,
			ResName: deriveColumnName(item, analyzed),
		}
		fillProvenance(te, q)
		*resNo++
		return []*TargetEntryQ{te}, nil

	default:
		// Bare expression (literal, function call, etc.)
		analyzed, err := analyzeExpr(c, item, scope)
		if err != nil {
			return nil, err
		}
		if row, ok := unwrapParenExpr(item).(*nodes.RowExpr); ok && len(row.Items) != 1 {
			return nil, errOperandColumns(1)
		}
		te := &TargetEntryQ{
			Expr:    analyzed,
			ResNo:   *resNo,
			ResName: deriveColumnName(item, analyzed),
		}
		fillProvenance(te, q)
		*resNo++
		return []*TargetEntryQ{te}, nil
	}
}

func unwrapParenExpr(e nodes.ExprNode) nodes.ExprNode {
	for {
		p, ok := e.(*nodes.ParenExpr)
		if !ok {
			return e
		}
		e = p.Expr
	}
}

// expandStar expands a star expression (SELECT * or SELECT [db.]t.*) into
// individual TargetEntryQ values.
func expandStar(schema, tableName string, q *Query, scope *analyzerScope, resNo *int) ([]*TargetEntryQ, error) {
	var result []*TargetEntryQ

	if tableName == "" {
		entries := scope.allEntries()
		if len(entries) == 0 {
			return nil, errNoTablesUsed()
		}
		// SELECT * — expand all tables in scope order.
		for _, e := range entries {
			entries, err := expandScopeEntry(e, q, scope, resNo)
			if err != nil {
				return nil, err
			}
			result = append(result, entries...)
		}
	} else {
		// SELECT [db.]t.* — expand only the named table.
		lower := toLower(tableName)
		if i := scope.entryIndex(lower); i >= 0 && (schema == "" || scope.entryDBName(i) == "" || eqFoldStr(scope.entryDBName(i), schema)) {
			entries, err := expandScopeEntry(scope.entry(i), q, scope, resNo)
			if err != nil {
				return nil, err
			}
			result = append(result, entries...)
		} else {
			if schema != "" {
				return nil, errUnknownTable(schema, tableName)
			}
			return nil, &Error{
				Code:     ErrUnknownTable,
				SQLState: sqlState(ErrUnknownTable),
				Message:  fmt.Sprintf("Unknown table '%s'", tableName),
			}
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no columns found for star expansion")
	}
	return result, nil
}

// expandScopeEntry expands all columns from a single scope entry,
// skipping columns marked as coalesced by USING/NATURAL.
func expandScopeEntry(e scopeEntry, q *Query, scope *analyzerScope, resNo *int) ([]*TargetEntryQ, error) {
	var result []*TargetEntryQ
	rte := q.RangeTable[e.rteIdx]
	for i, col := range e.columns {
		// Skip columns that were coalesced away by USING/NATURAL.
		if scope.isCoalesced(e.name, col.Name) {
			continue
		}
		te := &TargetEntryQ{
			Expr: &VarExprQ{
				RangeIdx: e.rteIdx,
				AttNum:   i + 1,
			},
			ResNo:        *resNo,
			ResName:      col.Name,
			ResOrigDB:    rte.DBName,
			ResOrigTable: rte.TableName,
			ResOrigCol:   col.Name,
		}
		*resNo++
		result = append(result, te)
	}
	return result, nil
}

// deriveColumnName generates the output column name for an unaliased expression.
func deriveColumnName(astNode nodes.ExprNode, _ AnalyzedExpr) string {
	switch n := astNode.(type) {
	case *nodes.ColumnRef:
		return n.Column
	case *nodes.IntLit:
		return strconv.FormatInt(n.Value, 10)
	case *nodes.StringLit:
		return n.Value
	case *nodes.FloatLit:
		return n.Value
	case *nodes.BoolLit:
		if n.Value {
			return "TRUE"
		}
		return "FALSE"
	case *nodes.NullLit:
		return "NULL"
	case *nodes.FuncCallExpr:
		return n.Name
	case *nodes.ParenExpr:
		return deriveColumnName(n.Expr, nil)
	default:
		return "?"
	}
}

// fillProvenance sets ResOrigDB/Table/Col when the expression is a VarExprQ.
func fillProvenance(te *TargetEntryQ, q *Query) {
	v, ok := te.Expr.(*VarExprQ)
	if !ok {
		return
	}
	if v.RangeIdx < 0 || v.RangeIdx >= len(q.RangeTable) {
		return
	}
	rte := q.RangeTable[v.RangeIdx]
	te.ResOrigDB = rte.DBName
	te.ResOrigTable = rte.TableName
	if v.AttNum >= 1 && v.AttNum <= len(rte.ColNames) {
		te.ResOrigCol = rte.ColNames[v.AttNum-1]
	}
}
