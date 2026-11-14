package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateLogicalPartitionTrackingStmt parses a CREATE LOGICAL PARTITION
// TRACKING statement. It is called after CREATE LOGICAL PARTITION TRACKING.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/create-logical-partition-tracking.html
//
//	CREATE LOGICAL PARTITION TRACKING ON [ schema. ] table_name
//	    PARTITION BY RANGE ( column_name [, column_name ]... )
//	    ( partition_definition [, partition_definition ]... )
//
//	partition_definition:
//	    PARTITION [ partition_name ] VALUES LESS THAN ( value [, value ]... )
//
// The page lists INTERVAL beside RANGE, which Oracle 23ai rejects
// (ORA-14151), and names every partition; Oracle takes a partition without a
// name, a key of several columns, and nothing after a partition's bounds
// (ORA-00922). A database link on the table is refused at its '@'
// (ORA-00949). The options record the key columns and the partitions.
func (p *Parser) parseCreateLogicalPartitionTrackingStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_LOGICAL_PARTITION_TRACKING,
		Loc:        nodes.Loc{Start: start},
	}

	// ON [ schema. ] table_name
	if p.cur.Type != kwON {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume ON
	name, err := p.parsePlainObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// PARTITION BY RANGE ( column_name [, column_name ]... )
	optStart := p.pos()
	if p.cur.Type != kwPARTITION {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume PARTITION
	if p.cur.Type != kwBY {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume BY
	if p.cur.Type != kwRANGE && !p.isKeywordStr("RANGE") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume RANGE
	columns := &nodes.List{}
	if err := p.parseCreateTablePartitionColumns(columns); err != nil {
		return nil, err
	}
	opts := &nodes.List{}
	opts.Items = append(opts.Items, &nodes.DDLOption{
		Key:   "PARTITION BY",
		Value: "RANGE",
		Items: columns,
		Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
	})

	// ( partition_definition [, partition_definition ]... )
	optStart = p.pos()
	if err := p.expectCreateTableToken('('); err != nil {
		return nil, err
	}
	partitions := &nodes.List{}
	for {
		defStart := p.pos()
		if p.cur.Type != kwPARTITION {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume PARTITION
		def := &nodes.PartitionDef{Values: &nodes.List{}, Loc: nodes.Loc{Start: defStart}}
		if p.isCreateTablePartitionName(createTableRangeBound) {
			def.Name = p.cur.Str
			p.advance()
		}
		if err := p.parseCreateTablePartitionBound(createTableRangeBound, def); err != nil {
			return nil, err
		}
		def.Loc.End = p.prev.End
		partitions.Items = append(partitions.Items, def)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ,
	}
	if err := p.expectCreateTableToken(')'); err != nil {
		return nil, err
	}
	opts.Items = append(opts.Items, &nodes.DDLOption{
		Key:   "PARTITIONS",
		Items: partitions,
		Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
	})

	stmt.Options = opts
	stmt.Loc.End = p.prev.End
	return stmt, nil
}
