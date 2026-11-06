package parser

import (
	nodes "github.com/bytebase/omni/mariadb/ast"
)

// What MariaDB's parser refuses in a partition clause after its productions read it.
// The actions run while the statement parses, so unlike on MySQL the checks of a
// definition against the method above it are parse errors too: 1479, 1480, 1492, 1564,
// 1653, 1657, 4113 and 4128 all come back with a syntax error after the statement,
// beside the 1064s with messages of their own.
//
// Ref: partitioning, part_definition, part_values_in and opt_versioning_rotation in
// MariaDB 11.8 sql/sql_yacc.yy; add_partition_values and the partition_info checks in
// sql/sql_lex.cc and sql/partition_info.cc.

// partitionBound is the kind of bound a partition definition wrote.
type partitionBound int

const (
	boundNone partitionBound = iota + 1
	boundLessThan
	boundIn
	boundDefault
	boundSystemTime
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
// which PARTITION BY refuses with 1504, ER_NO_PARTS_ERROR, when it is 0.
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
		if part.Type == nodes.PartitionRange || part.Type == nodes.PartitionList {
			return &ParseError{Message: "For RANGE and LIST partitions each partition must be defined", Position: part.Loc.Start}
		}
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
			return refuse("Syntax error: RANGE PARTITIONING requires definition of VALUES LESS THAN for each partition") // 1479
		case nodes.PartitionList:
			return refuse("Syntax error: LIST PARTITIONING requires definition of VALUES IN for each partition") // 1479
		case nodes.PartitionSystemTime:
			return refuse("Wrong partitions for system_time: each partition must be HISTORY or CURRENT") // 4128
		}
	case boundSystemTime:
		if part.Type != nodes.PartitionSystemTime {
			return refuse("Only SYSTEM_TIME PARTITIONING can use HISTORY or CURRENT in partition definition") // 4113
		}
	case boundDefault:
		if part.Type != nodes.PartitionList {
			return refuse("Only LIST PARTITIONING can use DEFAULT in partition definition") // 1480
		}
	case boundLessThan:
		if part.Type != nodes.PartitionRange {
			return refuse("Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition") // 1480
		}
		list, tuple := def.Values.(*nodes.List)
		switch {
		case columns == 0 && tuple:
			return refuse("Cannot have more than one value for this type of RANGE partitioning") // 1657
		case columns > 1 && !tuple:
			return refuse("Inconsistency in usage of column lists for partitioning") // 1064
		case columns > 0 && tuple && len(list.Items) != columns:
			return refuse("Inconsistency in usage of column lists for partitioning") // 1653
		}
	case boundIn:
		if part.Type != nodes.PartitionList {
			return refuse("Only LIST PARTITIONING can use VALUES IN in partition definition") // 1480
		}
		rows := def.Values.(*nodes.List).Items
		first, tuples := rows[0].(*nodes.List)
		if columns > 1 && !tuples {
			return refuse("Inconsistency in usage of column lists for partitioning") // 1064
		}
		// A row holds at least two values, so this also turns one away with fewer columns, a
		// 1064 there.
		if tuples && len(first.Items) != columns {
			return refuse("Inconsistency in usage of column lists for partitioning") // 1653
		}
	}
	return nil
}

// validatePartitionInRows is part_values_in, on every path. VALUES IN takes values or
// rows, where a row holds at least two values. A parenthesized lone value does not make a
// row, so `((1))` and `((1),(2))` are a 1064, and so is `((1,2),3)`; `(3,(1,2))` is 1564,
// a row read as an expression, while `(3,(1))` is two values. Rows hold one width, 1657
// when there is no method to say which and 1653 when there is.
func validatePartitionInRows(rows []nodes.Node, pos int) error {
	refuse := func(msg string) error {
		return &ParseError{Message: msg, Position: pos}
	}
	first, tuples := rows[0].(*nodes.List)
	for _, row := range rows {
		list, ok := row.(*nodes.List)
		switch {
		case tuples && (!ok || len(list.Items) < 2):
			return refuse("Inconsistency in usage of column lists for partitioning")
		case tuples && len(list.Items) != len(first.Items):
			return refuse("Inconsistency in usage of column lists for partitioning")
		case !tuples && ok && len(list.Items) > 1:
			return refuse("Partition constant is out of partition function domain")
		}
	}
	return nil
}

// refusePartitionValue holds a partition's values to what the action evaluates as it reads them:
// constants and the arithmetic partitioning allows, and nothing that names a column, which is
// 1054, ER_BAD_FIELD_ERROR. Anything else is 1564, ER_PARTITION_FUNCTION_IS_NOT_ALLOWED: `1 | 2`
// and `(1 OR 2)` parse and are no partition value. refusePartitionFunction answers for a call's
// name.
func refusePartitionValue(value nodes.Node, pos int) error {
	refused := func() error {
		return &ParseError{Message: "This partition function is not allowed", Position: pos}
	}
	stack := []nodes.Node{value}
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		switch e := n.(type) {
		case *nodes.String, *nodes.IntLit, *nodes.FloatLit, *nodes.StringLit, *nodes.BoolLit, *nodes.NullLit,
			*nodes.HexLit, *nodes.BitLit, *nodes.TemporalLit:
		case *nodes.ColumnRef:
			return &ParseError{Message: "Unknown column '" + e.Column + "' in 'partition function'", Position: pos}
		case *nodes.DefaultExpr:
			return &ParseError{Message: "Unknown column in 'partition function'", Position: pos}
		case *nodes.List:
			stack = append(stack, e.Items...)
		case *nodes.BinaryExpr:
			switch e.Op {
			case nodes.BinOpAdd, nodes.BinOpSub, nodes.BinOpMul, nodes.BinOpDivInt, nodes.BinOpMod:
				stack = append(stack, e.Left, e.Right)
			default:
				return refused()
			}
		case *nodes.UnaryExpr:
			if e.Op != nodes.UnaryMinus && e.Op != nodes.UnaryPlus {
				return refused()
			}
			stack = append(stack, e.Operand)
		case *nodes.ParenExpr:
			stack = append(stack, e.Expr)
		case *nodes.ExtractExpr:
			stack = append(stack, e.Expr)
		case *nodes.FuncCallExpr:
			for _, arg := range e.Args {
				stack = append(stack, arg)
			}
		default:
			return refused()
		}
	}
	return nil
}

// startsDefaultKeyword reports a DEFAULT that is no call to DEFAULT(col). DEFAULT is a
// partition's own bound and no value, so VALUES IN (DEFAULT) is a 1064, where the
// expression reader would take the word.
func (p *Parser) startsDefaultKeyword() bool {
	return p.cur.Type == kwDEFAULT && p.peekNext().Type != '('
}
