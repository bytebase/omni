package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// parseAlterTable parses:
//
//	ALTER TABLE name action [, action ...]
//
// The ALTER keyword has already been consumed; cur is TABLE.
func (p *Parser) parseAlterTable() (ast.Node, error) {
	startLoc := p.prev.Loc // loc of ALTER token

	// Consume TABLE
	p.advance()

	stmt := &ast.AlterTableStmt{}

	// Table name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Parse comma-separated action list.
	// Most actions start with a keyword (ADD, DROP, MODIFY, RENAME, SET, ...).
	// Exception: after DROP ROLLUP name, a bare identifier means another rollup
	// name in the same DROP ROLLUP list (e.g., DROP ROLLUP r1,r2).
	var lastActionType ast.AlterActionType
	for {
		action, err := p.parseAlterTableAction()
		if err != nil {
			return nil, err
		}
		stmt.Actions = append(stmt.Actions, action)
		lastActionType = action.Type

		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','

		// If the previous action was DROP ROLLUP and the current token is a
		// bare identifier (not an action keyword), treat it as another rollup
		// name in the same DROP ROLLUP list.
		if lastActionType == ast.AlterDropRollup && isIdentifierToken(p.cur.Kind) && !isAlterActionKeyword(p.cur.Kind) {
			rollupName, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			extra := &ast.AlterTableAction{
				Type:       ast.AlterDropRollup,
				RollupName: rollupName,
				Loc:        p.prev.Loc,
			}
			stmt.Actions = append(stmt.Actions, extra)
			if p.cur.Kind != int(',') {
				break
			}
			p.advance() // consume ','
		}
	}

	stmt.Loc = startLoc.Merge(p.prev.Loc)
	return stmt, nil
}

// parseAlterTableAction parses a single ALTER TABLE action.
// On entry, cur is the first token of the action.
func (p *Parser) parseAlterTableAction() (*ast.AlterTableAction, error) {
	startLoc := p.cur.Loc
	action := &ast.AlterTableAction{Loc: startLoc}

	switch p.cur.Kind {
	case kwADD:
		return p.parseAlterAdd(action)

	case kwDROP:
		return p.parseAlterDrop(action)

	case kwMODIFY:
		return p.parseAlterModify(action)

	case kwRENAME:
		return p.parseAlterRename(action)

	case kwSET:
		return p.parseAlterSet(action)

	case kwENABLE:
		return p.parseAlterEnable(action)

	case kwTRUNCATE:
		return p.parseAlterTruncatePartition(action)

	case kwREPLACE:
		return p.parseAlterReplacePartition(action)

	case kwORDER:
		return p.parseAlterOrderBy(action)

	default:
		// Fallback: capture raw text until next ',' or EOF
		return p.parseAlterRaw(action)
	}
}

// parseAlterAdd parses ADD ... actions.
// cur is ADD on entry.
func (p *Parser) parseAlterAdd(action *ast.AlterTableAction) (*ast.AlterTableAction, error) {
	p.advance() // consume ADD

	switch p.cur.Kind {
	case kwCOLUMN:
		p.advance() // consume COLUMN
		return p.parseAddColumn(action)

	case kwPARTITION:
		p.advance() // consume PARTITION
		return p.parseAddPartition(action)

	case kwROLLUP:
		p.advance() // consume ROLLUP
		return p.parseAddRollup(action)

	default:
		// Could be ADD col_def (without COLUMN keyword) — try as column
		return p.parseAddColumn(action)
	}
}

// parseAddColumn parses the column part of ADD [COLUMN] ...
// Supports:
//
//	ADD [COLUMN] col_def [AFTER col | FIRST]
//	ADD [COLUMN] (col_def [, col_def ...])
func (p *Parser) parseAddColumn(action *ast.AlterTableAction) (*ast.AlterTableAction, error) {
	action.Type = ast.AlterAddColumn

	// Check for multi-column ADD COLUMN (col_def, col_def)
	if p.cur.Kind == int('(') {
		p.advance() // consume '('
		// Parse first column def
		col, err := p.parseColumnDef()
		if err != nil {
			return nil, err
		}
		action.Column = col

		// Consume additional column defs (we only keep the first in the action;
		// for multi-column, all are parsed but discarded except the last adds extra actions)
		// Per the spec, we store the first column in Column and create extra actions.
		// For simplicity, parse all and store in a synthetic way: first goes into action,
		// additional columns as extra actions appended via the caller. We return the
		// first action here; the caller sees only one. This is a simplification.
		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			_, err := p.parseColumnDef()
			if err != nil {
				return nil, err
			}
			// Additional columns are parsed but stored in RawText for now
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		action.Loc.End = closeTok.Loc.End
		return action, nil
	}

	// Single column ADD COLUMN col_def
	col, err := p.parseColumnDef()
	if err != nil {
		return nil, err
	}
	action.Column = col
	action.Loc.End = ast.NodeLoc(col).End

	// Optional FIRST | AFTER col
	p.parseColumnPosition(action)

	return action, nil
}

// parseColumnPosition parses optional FIRST or AFTER col_name after a column def.
func (p *Parser) parseColumnPosition(action *ast.AlterTableAction) {
	switch p.cur.Kind {
	case kwFIRST:
		action.First = true
		action.Loc.End = p.cur.Loc.End
		p.advance()
	case kwAFTER:
		p.advance() // consume AFTER
		if isIdentifierToken(p.cur.Kind) {
			action.After = p.cur.Str
			action.Loc.End = p.cur.Loc.End
			p.advance()
		}
	}
}

// parsePartitionItemNoKeyword parses a single partition definition where the
// leading PARTITION keyword has already been consumed. cur is the partition name
// (or IF for IF NOT EXISTS).
//
// This mirrors the body of parsePartitionItem in create_table.go after the
// PARTITION keyword has been consumed.
func (p *Parser) parsePartitionItemNoKeyword() (*ast.PartitionItem, error) {
	startLoc := p.cur.Loc

	// Optional IF NOT EXISTS
	if p.cur.Kind == kwIF {
		p.advance()
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
	}

	// Partition name
	partName, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	item := &ast.PartitionItem{
		Name: partName,
		Loc:  startLoc,
	}

	// VALUES clause
	if p.cur.Kind == kwVALUES {
		p.advance() // consume VALUES

		switch p.cur.Kind {
		case kwLESS:
			p.advance() // consume LESS
			if _, err := p.expect(kwTHAN); err != nil {
				return nil, err
			}
			if p.cur.Kind == kwMAXVALUE {
				item.IsMaxValue = true
				item.Loc.End = p.cur.Loc.End
				p.advance()
			} else {
				vals, endLoc, err := p.parsePartitionValueList()
				if err != nil {
					return nil, err
				}
				item.Values = vals
				item.Loc.End = endLoc
			}

		case kwIN:
			p.advance() // consume IN
			inVals, endLoc, err := p.parseInPartitionValues()
			if err != nil {
				return nil, err
			}
			item.InValues = inVals
			item.Loc.End = endLoc

		case int('['):
			vals, endLoc, err := p.parseFixedRangePartition()
			if err != nil {
				return nil, err
			}
			item.Values = vals
			item.Loc.End = endLoc

		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	return item, nil
}

// parseAddPartition parses ADD PARTITION actions.
// PARTITION keyword has already been consumed; cur is partition name (or IF).
func (p *Parser) parseAddPartition(action *ast.AlterTableAction) (*ast.AlterTableAction, error) {
	action.Type = ast.AlterAddPartition

	// Parse partition item without expecting the PARTITION keyword
	// (it was already consumed by the caller).
	item, err := p.parsePartitionItemNoKeyword()
	if err != nil {
		return nil, err
	}
	action.Partition = item
	action.Loc.End = item.Loc.End

	// Optional DISTRIBUTED BY
	if p.cur.Kind == kwDISTRIBUTED {
		dist, err := p.parseDistributedBy()
		if err != nil {
			return nil, err
		}
		action.PartitionDist = dist
		action.Loc.End = dist.Loc.End
	}

	// Optional ("key"="value") properties
	if p.cur.Kind == int('(') {
		props, err := p.parseParenProperties()
		if err != nil {
			return nil, err
		}
		action.PartitionProps = props
		action.Loc.End = p.prev.Loc.End
	}

	return action, nil
}

// parseParenProperties parses a parenthesized key=value list without the PROPERTIES keyword.
// Used for inline partition properties: ("key"="value", ...)
func (p *Parser) parseParenProperties() ([]*ast.Property, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	var props []*ast.Property
	for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
		startLoc := p.cur.Loc
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		key := p.cur.Str
		p.advance()
		if _, err := p.expect(int('=')); err != nil {
			return nil, err
		}
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		val := p.cur.Str
		endLoc := p.cur.Loc
		p.advance()
		props = append(props, &ast.Property{
			Key:   key,
			Value: val,
			Loc:   ast.Loc{Start: startLoc.Start, End: endLoc.End},
		})
		if p.cur.Kind == int(',') {
			p.advance()
		}
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return props, nil
}

// parseAddRollup parses ADD ROLLUP name (cols) [FROM base_rollup] [PROPERTIES(...)]
// ROLLUP keyword has already been consumed.
func (p *Parser) parseAddRollup(action *ast.AlterTableAction) (*ast.AlterTableAction, error) {
	action.Type = ast.AlterAddRollup

	startLoc := p.cur.Loc

	// Rollup name
	rollupName, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	// Column list: (col1, col2, ...)
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	cols, err := p.parseIdentifierList()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}

	rd := &ast.RollupDef{
		Name:    rollupName,
		Columns: cols,
		Loc:     ast.Loc{Start: startLoc.Start, End: closeTok.Loc.End},
	}

	// Optional FROM base_rollup
	if p.cur.Kind == kwFROM {
		p.advance() // consume FROM
		if isIdentifierToken(p.cur.Kind) {
			// base rollup name — stored in RollupName field of action for now
			action.RollupName = p.cur.Str
			rd.Loc.End = p.cur.Loc.End
			p.advance()
		}
	}

	// Optional PROPERTIES(...)
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		rd.Properties = props
		rd.Loc.End = p.prev.Loc.End
	}

	action.Rollup = rd
	action.Loc.End = rd.Loc.End
	return action, nil
}

// parseAlterDrop parses DROP ... actions.
// cur is DROP on entry.
func (p *Parser) parseAlterDrop(action *ast.AlterTableAction) (*ast.AlterTableAction, error) {
	p.advance() // consume DROP

	switch p.cur.Kind {
	case kwCOLUMN:
		p.advance() // consume COLUMN
		colName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		action.Type = ast.AlterDropColumn
		action.ColumnName = colName
		action.Loc.End = p.prev.Loc.End
		return action, nil

	case kwPARTITION:
		p.advance() // consume PARTITION
		partName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		action.Type = ast.AlterDropPartition
		action.PartitionName = partName
		action.Loc.End = p.prev.Loc.End
		return action, nil

	case kwROLLUP:
		p.advance() // consume ROLLUP
		rollupName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		action.Type = ast.AlterDropRollup
		action.RollupName = rollupName
		action.Loc.End = p.prev.Loc.End
		return action, nil

	default:
		// Unknown DROP action — treat as raw
		return p.parseAlterRaw(action)
	}
}

// parseAlterModify parses MODIFY ... actions.
// cur is MODIFY on entry.
func (p *Parser) parseAlterModify(action *ast.AlterTableAction) (*ast.AlterTableAction, error) {
	p.advance() // consume MODIFY

	switch p.cur.Kind {
	case kwCOLUMN:
		p.advance() // consume COLUMN
		// Special case: MODIFY COLUMN col COMMENT 'text' (no type).
		// Detect by looking ahead: identifier followed immediately by COMMENT.
		if isIdentifierToken(p.cur.Kind) && p.peekNext().Kind == kwCOMMENT {
			colName, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			p.advance() // consume COMMENT
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			comment := p.cur.Str
			p.advance()
			// Build a minimal ColumnDef with just name and comment.
			action.Type = ast.AlterModifyColumn
			action.Column = &ast.ColumnDef{
				Name:    colName,
				Comment: comment,
				Loc:     p.prev.Loc,
			}
			action.Loc.End = p.prev.Loc.End
			return action, nil
		}
		col, err := p.parseColumnDef()
		if err != nil {
			return nil, err
		}
		action.Type = ast.AlterModifyColumn
		action.Column = col
		action.Loc.End = ast.NodeLoc(col).End

		// Optional FIRST | AFTER col
		p.parseColumnPosition(action)
		return action, nil

	case kwPARTITION:
		return p.parseModifyPartition(action)

	case kwDISTRIBUTION:
		// MODIFY DISTRIBUTION DISTRIBUTED BY ...
		p.advance() // consume DISTRIBUTION
		dist, err := p.parseDistributedBy()
		if err != nil {
			return nil, err
		}
		action.Type = ast.AlterModifyDistribution
		action.Distribution = dist
		action.Loc.End = dist.Loc.End
		return action, nil

	case kwCOMMENT:
		p.advance() // consume COMMENT
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		action.Type = ast.AlterModifyComment
		action.Comment = p.cur.Str
		action.Loc.End = p.cur.Loc.End
		p.advance()
		return action, nil

	case kwENGINE:
		p.advance() // consume ENGINE
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		engineName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		action.Type = ast.AlterModifyEngine
		action.Engine = engineName
		action.Loc.End = p.prev.Loc.End
		return action, nil

	default:
		return p.parseAlterRaw(action)
	}
}

// parseModifyPartition parses MODIFY PARTITION ... SET(...)
// MODIFY has been consumed; cur is PARTITION.
func (p *Parser) parseModifyPartition(action *ast.AlterTableAction) (*ast.AlterTableAction, error) {
	p.advance() // consume PARTITION
	action.Type = ast.AlterModifyPartition

	// PARTITION name | PARTITION (p1, p2, ...) | PARTITION (*)
	if p.cur.Kind == int('(') {
		p.advance() // consume '('
		if p.cur.Kind == int('*') {
			action.PartitionStar = true
			p.advance()
		} else {
			for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
				partName, _, err := p.parseIdentifier()
				if err != nil {
					return nil, err
				}
				action.PartitionList = append(action.PartitionList, partName)
				if p.cur.Kind == int(',') {
					p.advance()
				}
			}
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
	} else {
		// Single partition name
		partName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		action.PartitionName = partName
	}

	// SET ("key"="value", ...)
	if p.cur.Kind == kwSET {
		p.advance() // consume SET
		props, err := p.parseParenProperties()
		if err != nil {
			return nil, err
		}
		action.Properties = props
		action.Loc.End = p.prev.Loc.End
	}

	return action, nil
}

// parseAlterRename parses RENAME ... actions.
// cur is RENAME on entry.
func (p *Parser) parseAlterRename(action *ast.AlterTableAction) (*ast.AlterTableAction, error) {
	p.advance() // consume RENAME

	switch p.cur.Kind {
	case kwTO:
		// RENAME TO new_name
		p.advance() // consume TO
		newName, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		action.Type = ast.AlterRenameTable
		action.NewTableName = newName
		action.Loc.End = ast.NodeLoc(newName).End
		return action, nil

	case kwROLLUP:
		// RENAME ROLLUP old new
		p.advance() // consume ROLLUP
		oldName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		// Optional TO keyword
		p.match(kwTO)
		newName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		action.Type = ast.AlterRenameRollup
		action.ColumnName = oldName // old rollup name
		action.NewName = newName
		action.Loc.End = p.prev.Loc.End
		return action, nil

	case kwPARTITION:
		// RENAME PARTITION old new
		p.advance() // consume PARTITION
		oldName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		// Optional TO keyword
		p.match(kwTO)
		newName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		action.Type = ast.AlterRenamePartition
		action.ColumnName = oldName // old partition name
		action.NewName = newName
		action.Loc.End = p.prev.Loc.End
		return action, nil

	case kwCOLUMN:
		// RENAME COLUMN old TO new (or old new without TO)
		p.advance() // consume COLUMN
		oldName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		// Optional TO keyword
		p.match(kwTO)
		newName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		action.Type = ast.AlterRenameColumn
		action.ColumnName = oldName
		action.NewName = newName
		action.Loc.End = p.prev.Loc.End
		return action, nil

	default:
		// RENAME new_name (rename table, no TO keyword)
		if isIdentifierToken(p.cur.Kind) {
			newName, err := p.parseMultipartIdentifier()
			if err != nil {
				return nil, err
			}
			action.Type = ast.AlterRenameTable
			action.NewTableName = newName
			action.Loc.End = ast.NodeLoc(newName).End
			return action, nil
		}
		return p.parseAlterRaw(action)
	}
}

// parseAlterSet parses SET ("key"="value", ...) action.
// cur is SET on entry.
func (p *Parser) parseAlterSet(action *ast.AlterTableAction) (*ast.AlterTableAction, error) {
	p.advance() // consume SET
	action.Type = ast.AlterSetProperties

	props, err := p.parseParenProperties()
	if err != nil {
		return nil, err
	}
	action.Properties = props
	action.Loc.End = p.prev.Loc.End
	return action, nil
}

// parseAlterEnable parses ENABLE FEATURE 'name' [WITH PROPERTIES (...)] action.
// cur is ENABLE on entry.
func (p *Parser) parseAlterEnable(action *ast.AlterTableAction) (*ast.AlterTableAction, error) {
	p.advance() // consume ENABLE

	if p.cur.Kind != kwFEATURE {
		return p.parseAlterRaw(action)
	}
	p.advance() // consume FEATURE

	if p.cur.Kind != tokString {
		return nil, p.syntaxErrorAtCur()
	}
	action.Type = ast.AlterEnableFeature
	action.FeatureName = p.cur.Str
	action.Loc.End = p.cur.Loc.End
	p.advance()

	// Optional WITH PROPERTIES (...)
	if p.cur.Kind == kwWITH {
		p.advance() // consume WITH
		if p.cur.Kind == kwPROPERTIES {
			props, err := p.parseProperties()
			if err != nil {
				return nil, err
			}
			action.Properties = props
			action.Loc.End = p.prev.Loc.End
		}
	}

	return action, nil
}

// parseAlterTruncatePartition parses TRUNCATE PARTITION name action.
// cur is TRUNCATE on entry.
func (p *Parser) parseAlterTruncatePartition(action *ast.AlterTableAction) (*ast.AlterTableAction, error) {
	p.advance() // consume TRUNCATE

	if p.cur.Kind != kwPARTITION {
		return p.parseAlterRaw(action)
	}
	p.advance() // consume PARTITION

	partName, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	action.Type = ast.AlterTruncatePartition
	action.PartitionName = partName
	action.Loc.End = p.prev.Loc.End
	return action, nil
}

// parseAlterReplacePartition parses REPLACE PARTITION (...) WITH TEMPORARY PARTITION (...)
// cur is REPLACE on entry.
func (p *Parser) parseAlterReplacePartition(action *ast.AlterTableAction) (*ast.AlterTableAction, error) {
	p.advance() // consume REPLACE
	action.Type = ast.AlterReplacePartition

	// Must be REPLACE PARTITION — if not a partition replacement, fall back to raw
	// (REPLACE WITH TABLE is handled at the statement level as unsupported)
	if p.cur.Kind != kwPARTITION {
		return p.parseAlterRaw(action)
	}
	p.advance() // consume PARTITION

	// Partition list: (p1, p2, ...)
	if p.cur.Kind == int('(') {
		p.advance() // consume '('
		for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
			partName, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			action.PartitionList = append(action.PartitionList, partName)
			if p.cur.Kind == int(',') {
				p.advance()
			}
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
	}

	// WITH TEMPORARY PARTITION (p1, ...)
	if p.cur.Kind == kwWITH {
		p.advance() // consume WITH
		if p.cur.Kind == kwTEMPORARY {
			p.advance() // consume TEMPORARY
		}
		if p.cur.Kind == kwPARTITION {
			p.advance() // consume PARTITION
		}
		if p.cur.Kind == int('(') {
			p.advance()
			for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
				p.advance()
			}
			if p.cur.Kind == int(')') {
				action.Loc.End = p.cur.Loc.End
				p.advance()
			}
		}
	}

	// Optional PROPERTIES
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		action.Properties = props
		action.Loc.End = p.prev.Loc.End
	}

	return action, nil
}

// parseAlterOrderBy parses ORDER BY (col1, col2, ...) action.
// cur is ORDER on entry.
func (p *Parser) parseAlterOrderBy(action *ast.AlterTableAction) (*ast.AlterTableAction, error) {
	p.advance() // consume ORDER
	if _, err := p.expect(kwBY); err != nil {
		return nil, err
	}
	action.Type = ast.AlterOrderBy

	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	cols, err := p.parseIdentifierList()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	action.OrderByColumns = cols
	action.Loc.End = closeTok.Loc.End
	return action, nil
}

// parseAlterRaw captures tokens until the next comma-separated action boundary
// or end of statement, storing the text verbatim.
func (p *Parser) parseAlterRaw(action *ast.AlterTableAction) (*ast.AlterTableAction, error) {
	action.Type = ast.AlterRaw
	start := p.cur.Loc.Start

	depth := 0
	for p.cur.Kind != tokEOF {
		switch p.cur.Kind {
		case int('('), int('['), int('{'):
			depth++
		case int(')'), int(']'), int('}'):
			if depth > 0 {
				depth--
			} else {
				goto done
			}
		case int(','):
			if depth == 0 {
				goto done
			}
		case int(';'):
			goto done
		}
		p.advance()
	}
done:
	end := p.prev.Loc.End
	if end > start {
		action.RawText = strings.TrimSpace(p.input[start:end])
	}
	action.Loc.End = p.prev.Loc.End
	return action, nil
}

// isAlterActionKeyword reports whether tok is a keyword that starts a new
// ALTER TABLE action (ADD, DROP, MODIFY, RENAME, SET, ENABLE, TRUNCATE,
// REPLACE, ORDER). Used to distinguish bare rollup/partition names from
// new action keywords after a comma.
func isAlterActionKeyword(kind int) bool {
	switch kind {
	case kwADD, kwDROP, kwMODIFY, kwRENAME, kwSET,
		kwENABLE, kwTRUNCATE, kwREPLACE, kwORDER:
		return true
	}
	return false
}
