// Package parser - returning.go implements the MariaDB RETURNING clause
// (BYT-9135 P2) on INSERT/REPLACE and single-table DELETE.
//
// MariaDB accepts RETURNING <select_list> on INSERT (all forms), REPLACE, and
// single-table DELETE; it has no UPDATE RETURNING and no multi-table DELETE
// RETURNING (both 1064). The list is a SELECT target-list (*, exprs, aliases,
// subqueries; aggregates parse and are semantic-rejected, i.e. in contract).
package parser

import (
	nodes "github.com/bytebase/omni/mariadb/ast"
)

// parseReturningClause parses RETURNING <select_list>. The caller has confirmed
// p.cur is RETURNING. It reuses the SELECT expr-list parser, so a bare RETURNING
// (no list) fails there — matching MariaDB's 1064.
func (p *Parser) parseReturningClause() ([]nodes.Node, error) {
	p.advance() // consume RETURNING
	exprs, err := p.parseSelectExprList()
	if err != nil {
		return nil, err
	}
	list := make([]nodes.Node, len(exprs))
	for i, e := range exprs {
		list[i] = e
	}
	return list, nil
}
