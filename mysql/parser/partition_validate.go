package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
)

// What MySQL's parser refuses in a partition clause after its productions read it. The
// grammar takes a looser shape than the server does and the actions, which run while
// the statement parses, turn the rest away: PARTITIONS 0 with 1504, and the others with
// a 1064 carrying their own message. mysql:8.0 raises each of them with a syntax error
// after the statement as well, which is what tells them from the checks it runs after
// parsing (1480, 1492, 1564, 1653, 1657), and those stay out of Parse.
//
// Ref: PT_partition, PT_part_definition, PT_part_values_in_item and
// PT_part_values_in_list in MySQL 8.0 sql/parse_tree_partitions.cc, opt_key_algo and
// part_values_in in sql/sql_yacc.yy.

// partitionBound is the kind of bound a partition definition wrote.
type partitionBound int

const (
	boundNone partitionBound = iota + 1
	boundLessThan
	boundIn
)

// parsePartitionKeyAlgorithm reads KEY's optional ALGORITHM = {1 | 2}, where the equals
// sign is required. It returns 0 when the clause is not there.
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

// parsePartitionCount reads the required number after PARTITIONS or SUBPARTITIONS,
// which PARTITION BY refuses with 1504, ErrNoParts, when it is 0.
func (p *Parser) parsePartitionCount(what string) (int, error) {
	if p.cur.Type != tokICONST {
		return 0, p.syntaxErrorAtCur()
	}
	if p.cur.Ival == 0 {
		return 0, &ParseError{Message: "Number of " + what + " = 0 is not an allowed value", Position: p.cur.Loc}
	}
	return int(p.advance().Ival), nil
}

// validatePartitionDefList holds a list of definitions to itself, which every path that
// reads one does: a partition takes as many subpartitions as the first one in the list.
func validatePartitionDefList(defs []*nodes.PartitionDef) error {
	for _, def := range defs[1:] {
		if len(def.SubPartitions) != len(defs[0].SubPartitions) {
			return &ParseError{
				Message:  "Wrong number of subpartitions defined, mismatch with previous setting",
				Position: def.Loc.Start,
			}
		}
	}
	return nil
}

// validatePartitionClause holds the definitions of PARTITION BY to the method above
// them. ADD PARTITION and REORGANIZE PARTITION have no method to hold them to. bounds
// runs beside part.Partitions.
func validatePartitionClause(part *nodes.PartitionClause, bounds []partitionBound) error {
	defs := part.Partitions
	if len(defs) == 0 {
		return nil
	}
	if part.NumParts != 0 && part.NumParts != len(defs) {
		return &ParseError{
			Message:  "Wrong number of partitions defined, mismatch with previous setting",
			Position: part.Loc.Start,
		}
	}
	if subs := len(defs[0].SubPartitions); part.NumSubParts != 0 && subs != 0 && part.NumSubParts != subs {
		return &ParseError{
			Message:  "Wrong number of subpartitions defined, mismatch with previous setting",
			Position: part.Loc.Start,
		}
	}
	for i, def := range defs {
		if err := validatePartitionBound(part, def, bounds[i]); err != nil {
			return err
		}
	}
	return nil
}

func validatePartitionBound(part *nodes.PartitionClause, def *nodes.PartitionDef, bound partitionBound) error {
	refuse := func(msg string) error {
		return &ParseError{Message: msg, Position: def.Loc.Start}
	}
	columns := len(part.Columns)
	switch bound {
	case boundNone:
		switch part.Type {
		case nodes.PartitionRange:
			return refuse("Syntax error: RANGE PARTITIONING requires definition of VALUES LESS THAN for each partition")
		case nodes.PartitionList:
			return refuse("Syntax error: LIST PARTITIONING requires definition of VALUES IN for each partition")
		}
	case boundLessThan:
		// A bound of one value, MAXVALUE or not, reads as a single column.
		if _, tuple := def.Values.(*nodes.List); part.Type == nodes.PartitionRange && columns > 1 && !tuple {
			return refuse("Inconsistency in usage of column lists for partitioning")
		}
	case boundIn:
		list, ok := def.Values.(*nodes.List)
		if !ok || part.Type != nodes.PartitionList {
			return nil
		}
		_, tuples := list.Items[0].(*nodes.List)
		if columns <= 1 && tuples {
			return refuse("Row expressions in VALUES IN only allowed for multi-field column partitioning")
		}
		// Flat values name no columns, unless a row among them, read as an expression,
		// gets to answer first with 1564 after parsing: `(10,(20,30))`.
		if columns > 1 && !tuples && !holdsRow(list.Items) {
			return refuse("Inconsistency in usage of column lists for partitioning")
		}
	}
	return nil
}

func holdsRow(items []nodes.Node) bool {
	for _, item := range items {
		if list, ok := item.(*nodes.List); ok && len(list.Items) > 1 {
			return true
		}
	}
	return false
}

// validatePartitionInRows is part_values_in: VALUES IN takes values, of which a later
// one may be parenthesized, or takes rows of at least two values and nothing else. A
// parenthesized lone value does not make a row, so `((1))` and `((1),(2))` are a 1064
// on every path, and so is `((1,2),3)`, where `(3,(1,2))` parses and is refused after.
func validatePartitionInRows(rows []nodes.Node, pos int) error {
	if _, tuples := rows[0].(*nodes.List); !tuples {
		return nil
	}
	for _, row := range rows {
		if list, ok := row.(*nodes.List); !ok || len(list.Items) < 2 {
			return &ParseError{Message: "Inconsistency in usage of column lists for partitioning", Position: pos}
		}
	}
	return nil
}

// startsDefaultKeyword reports a DEFAULT that is no call to DEFAULT(col). MySQL has no
// DEFAULT partition, so VALUES IN (DEFAULT) is a 1064, where the expression reader would
// take the word.
func (p *Parser) startsDefaultKeyword() bool {
	return p.cur.Type == kwDEFAULT && p.peekNext().Type != '('
}
