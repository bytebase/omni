package parser

import (
	nodes "github.com/bytebase/omni/tidb/ast"
)

// What TiDB's grammar actions refuse in a PARTITION BY clause after its productions read
// it. PartitionNumOpt, SubPartitionNumOpt and PartitionKeyAlgorithmOpt check one value;
// PartitionOpt calls PartitionOptions.Validate, which holds the definitions to the
// method and to each other. The server raises these while it parses, each with its own
// code, and CREATE TABLE and ALTER TABLE ... PARTITION BY read the same clause.
// ADD PARTITION and REORGANIZE PARTITION go through neither: they read
// PartitionDefinitionList with no method to hold it to.
//
// Ref: PartitionOpt, PartitionNumOpt, SubPartitionNumOpt, PartitionKeyAlgorithmOpt in
// TiDB v8.5.5 pkg/parser/parser.y; PartitionOptions.Validate and the
// PartitionDefinitionClause Validate methods in pkg/parser/ast/ddl.go.

// partitionBound is the kind of bound a partition definition wrote.
type partitionBound int

const (
	boundNone partitionBound = iota + 1
	boundLessThan
	boundIn
)

// parsePartitionKeyAlgorithm reads KEY's optional ALGORITHM = NUM, where the equals sign
// is required and the action answers 1149 for a number that is neither 1 nor 2. It
// returns 0 when the clause is not there.
func (p *Parser) parsePartitionKeyAlgorithm() (int, error) {
	if p.cur.Type != kwALGORITHM {
		return 0, nil
	}
	p.advance()
	if _, err := p.expect('='); err != nil {
		return 0, err
	}
	if p.cur.Type != tokICONST {
		return 0, p.syntaxErrorAtCur()
	}
	if p.cur.Ival != 1 && p.cur.Ival != 2 {
		return 0, &ParseError{Message: "a KEY partition's ALGORITHM is 1 or 2", Position: p.cur.Loc}
	}
	return int(p.advance().Ival), nil
}

// parsePartitionCount reads the number after PARTITIONS or SUBPARTITIONS, which the
// action refuses with 1504, ErrNoParts, when it is 0.
func (p *Parser) parsePartitionCount(what string) (int, error) {
	if p.cur.Type != tokICONST {
		return 0, p.syntaxErrorAtCur()
	}
	if p.cur.Ival == 0 {
		return 0, &ParseError{Message: "Number of " + what + " = 0 is not an allowed value", Position: p.cur.Loc}
	}
	return int(p.advance().Ival), nil
}

// validatePartitionClause is PartitionOptions.Validate. bounds runs beside
// part.Partitions.
func validatePartitionClause(part *nodes.PartitionClause, bounds []partitionBound) error {
	defs := part.Partitions
	refuse := func(pos int, msg string) error {
		return &ParseError{Message: msg, Position: pos}
	}
	if part.NumParts != 0 && len(defs) != 0 && part.NumParts != len(defs) {
		return refuse(part.Loc.Start, "Wrong number of partitions defined, mismatch with previous setting") // 1484
	}
	if len(defs) > 0 {
		subs := len(defs[0].SubPartitions)
		for _, def := range defs[1:] {
			if len(def.SubPartitions) != subs {
				return refuse(def.Loc.Start, "Wrong number of subpartitions defined, mismatch with previous setting") // 1485
			}
		}
		if part.SubPartType != 0 {
			if part.NumSubParts != 0 && subs != 0 && part.NumSubParts != subs {
				return refuse(part.Loc.Start, "Wrong number of subpartitions defined, mismatch with previous setting") // 1485
			}
		} else if subs != 0 {
			return refuse(defs[0].Loc.Start, "It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning") // 1500
		}
	}
	if (part.Type == nodes.PartitionRange || part.Type == nodes.PartitionList) && len(defs) == 0 {
		return refuse(part.Loc.Start, "For RANGE and LIST partitions each partition must be defined") // 1492
	}
	for i, def := range defs {
		if err := validatePartitionBound(part, def, bounds[i]); err != nil {
			return err
		}
	}
	return nil
}

// validatePartitionBound is the Validate method of the definition's clause.
func validatePartitionBound(part *nodes.PartitionClause, def *nodes.PartitionDef, bound partitionBound) error {
	refuse := func(msg string) error {
		return &ParseError{Message: msg, Position: def.Loc.Start}
	}
	columns := len(part.Columns)
	switch bound {
	case boundNone:
		switch part.Type {
		case nodes.PartitionRange:
			return refuse("Syntax error: RANGE PARTITIONING requires definition of VALUES LESS THAN for each partition") // 1479
		case nodes.PartitionList:
			return refuse("Syntax error: LIST PARTITIONING requires definition of VALUES IN for each partition") // 1479
		}
	case boundLessThan:
		if part.Type != nodes.PartitionRange {
			return refuse("Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition") // 1480
		}
		values := 1
		if list, ok := def.Values.(*nodes.List); ok {
			values = len(list.Items)
		}
		if columns == 0 && values != 1 {
			return refuse("Cannot have more than one value for this type of RANGE partitioning") // 1657
		}
		if columns > 0 && values != columns {
			return refuse("Inconsistency in usage of column lists for partitioning") // 1653
		}
	case boundIn:
		if part.Type != nodes.PartitionList {
			return refuse("Only LIST PARTITIONING can use VALUES IN in partition definition") // 1480
		}
		list, ok := def.Values.(*nodes.List)
		if !ok {
			return nil // the bare DEFAULT, which has no values to count
		}
		return validateInValues(list.Items, columns, refuse)
	}
	return nil
}

// validateInValues is PartitionDefinitionClauseIn.Validate over the values of one
// VALUES IN. A row of values is a tuple or a lone value, and DEFAULT may stand alone in
// place of a row of any width.
func validateInValues(rows []nodes.Node, columns int, refuse func(string) error) error {
	next := 1
	expected := partitionRowWidth(rows[0])
	if isPartitionDefault(rows[0]) && len(rows) > 1 {
		expected = partitionRowWidth(rows[1])
		next = 2
	}
	for _, row := range rows[next:] {
		if partitionRowWidth(row) != expected && !isPartitionDefault(row) {
			return refuse("Inconsistency in usage of column lists for partitioning") // 1653
		}
	}
	if columns == 0 && expected != 1 {
		return refuse("Row expressions in VALUES IN only allowed for multi-field column partitioning") // 1658
	}
	if columns > 0 && expected != columns && !(len(rows) == 1 && isPartitionDefault(rows[0])) {
		return refuse("Inconsistency in usage of column lists for partitioning") // 1653
	}
	return nil
}

func partitionRowWidth(row nodes.Node) int {
	if list, ok := row.(*nodes.List); ok {
		return len(list.Items)
	}
	return 1
}

func isPartitionDefault(row nodes.Node) bool {
	s, ok := row.(*nodes.String)
	return ok && s.Str == "DEFAULT"
}
