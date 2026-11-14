package parser

import (
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
)

// This file reads StarRocks's alterClause (StarRocks.g4:874), which
// ALTER TABLE (alterTableStatement, StarRocks.g4:509) and ALTER SYSTEM
// (alterSystemStatement, StarRocks.g4:746) share. Each reader follows its
// grammar rule; where StarRocks's AST builder refuses a form the grammar
// derives, the comment says so and the compat rows
// (testdata/compat.tsv) are the evidence.

// alterWordTok reports whether t is the bare word w: a keyword or an
// unquoted identifier spelled w. Many of StarRocks's keywords are plain
// identifiers to omni's lexer (SWAP, IMAGE, NODE, ...), so the readers here
// match words by their text.
func alterWordTok(t Token, w string) bool {
	return (t.Kind == tokIdent || t.Kind >= 700) && strings.EqualFold(t.Str, w)
}

// alterAt reports whether cur is the word w.
func (p *Parser) alterAt(w string) bool { return alterWordTok(p.cur, w) }

// alterNextAt reports whether the token after cur is the word w.
func (p *Parser) alterNextAt(w string) bool { return alterWordTok(p.peekNext(), w) }

// alterExpectWord consumes the word w or fails at cur.
func (p *Parser) alterExpectWord(w string) (Token, error) {
	if p.alterAt(w) {
		return p.advance(), nil
	}
	return Token{}, p.syntaxErrorAtCur()
}

// alterExpectString consumes a string literal (StarRocks.g4:2746 string).
func (p *Parser) alterExpectString() (Token, error) {
	if p.cur.Kind == tokString {
		return p.advance(), nil
	}
	return Token{}, p.syntaxErrorAtCur()
}

// alterExpectInt consumes an INTEGER_VALUE.
func (p *Parser) alterExpectInt() (Token, error) {
	if p.cur.Kind == tokInt {
		return p.advance(), nil
	}
	return Token{}, p.syntaxErrorAtCur()
}

// alterStringList reads string (',' string)*, splitting a string that holds
// several comma-separated hosts the way Hosts always recorded them.
func (p *Parser) alterStringList() ([]string, error) {
	var hosts []string
	for {
		tok, err := p.alterExpectString()
		if err != nil {
			return nil, err
		}
		for _, h := range strings.Split(tok.Str, ",") {
			if h = strings.TrimSpace(h); h != "" {
				hosts = append(hosts, h)
			}
		}
		if p.cur.Kind != int(',') {
			return hosts, nil
		}
		p.advance() // consume ','
	}
}

// alterIdentifierList reads identifierList (StarRocks.g4:2873):
// '(' identifier (',' identifier)* ')'.
func (p *Parser) alterIdentifierList() ([]string, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	names, err := p.parseIdentifierList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return names, nil
}

// alterOptProperties reads an optional properties clause
// (StarRocks.g4:2699 PROPERTIES '(' property (',' property)* ')').
func (p *Parser) alterOptProperties(action *ast.AlterTableAction) error {
	if p.cur.Kind != kwPROPERTIES {
		return nil
	}
	props, err := p.parseProperties()
	if err != nil {
		return err
	}
	action.Properties = props
	return nil
}

// alterClauseEmpty reports whether cur ends the clause list where a clause
// would start: the empty optimizeClause (StarRocks.g4:1028), every part of
// which is optional.
func (p *Parser) alterClauseEmpty() bool {
	return p.cur.Kind == tokEOF || p.cur.Kind == int(',') || p.cur.Kind == int(';')
}

// parseAlterClause reads one table or partition alternative of alterClause
// (StarRocks.g4:874). cur is the clause's first token. An optimizeClause
// (StarRocks.g4:1028) may be empty, so a clause that ends where it starts is
// an AlterOptimize action with nothing in it.
func (p *Parser) parseAlterClause() (*ast.AlterTableAction, error) {
	action := &ast.AlterTableAction{Loc: p.cur.Loc}
	var err error
	switch {
	case p.alterClauseEmpty():
		// optimizeClause with every part absent.
		action.Type = ast.AlterOptimize
		action.Loc = ast.Loc{Start: p.cur.Loc.Start, End: p.cur.Loc.Start}
		return action, nil
	case p.cur.Kind == kwADD:
		err = p.parseAlterAddClause(action)
	case p.cur.Kind == kwDROP:
		err = p.parseAlterDropClause(action)
	case p.cur.Kind == kwMODIFY:
		err = p.parseAlterModifyClause(action)
	case p.cur.Kind == kwRENAME:
		err = p.parseAlterRenameClause(action)
	case p.alterAt("SWAP"):
		// StarRocks.g4:1016 swapTableClause: SWAP WITH identifier.
		p.advance()
		if _, err = p.alterExpectWord("WITH"); err == nil {
			action.Type = ast.AlterSwapTable
			action.NewName, _, err = p.parseIdentifier()
		}
	case p.cur.Kind == kwSET:
		// StarRocks.g4:1020 modifyPropertiesClause: SET propertyList.
		p.advance()
		action.Type = ast.AlterSetProperties
		action.Properties, err = p.parsePropertyList()
	case p.cur.Kind == kwCOMMENT:
		// StarRocks.g4:1024 modifyCommentClause: COMMENT '=' string.
		p.advance()
		if _, err = p.expect(int('=')); err == nil {
			var tok Token
			if tok, err = p.alterExpectString(); err == nil {
				action.Type = ast.AlterModifyComment
				action.Comment = tok.Str
			}
		}
	case p.alterAt("BASE") || p.alterAt("CUMULATIVE") || p.alterAt("COMPACT"):
		err = p.parseAlterCompactClause(action)
	case p.cur.Kind == kwCREATE || p.cur.Kind == kwREPLACE:
		err = p.parseAlterCreateOrReplaceClause(action)
	case p.cur.Kind == kwEXECUTE:
		// StarRocks.g4:1102 tableOperationClause: EXECUTE functionCall.
		err = p.parseAlterTableOperation(action)
	case p.cur.Kind == kwTRUNCATE:
		// StarRocks.g4:1164 truncatePartitionClause: TRUNCATE partitionNames.
		p.advance()
		action.Type = ast.AlterTruncatePartition
		err = p.parseAlterPartitionNames(action)
	case p.cur.Kind == kwORDER && p.peekNext().Kind == kwBY:
		err = p.parseAlterOrderByClause(action)
	default:
		err = p.parseAlterOptimizeClause(action)
	}
	if err != nil {
		return nil, err
	}
	action.Loc.End = p.prev.Loc.End
	return action, nil
}

// parseAlterAddClause reads the ADD alternatives of alterClause. cur is ADD.
func (p *Parser) parseAlterAddClause(action *ast.AlterTableAction) error {
	p.advance() // consume ADD
	switch {
	case p.cur.Kind == kwINDEX:
		return p.parseAlterAddIndex(action)
	case p.cur.Kind == kwCOLUMN:
		p.advance() // consume COLUMN
		return p.parseAlterAddColumn(action)
	case p.cur.Kind == kwTEMPORARY || p.cur.Kind == kwPARTITION || p.cur.Kind == kwPARTITIONS:
		return p.parseAlterAddPartition(action)
	}
	return p.syntaxErrorAtCur()
}

// parseAlterAddIndex reads createIndexClause (StarRocks.g4:1004):
//
//	ADD INDEX identifier identifierList (indexType propertyList?)? comment?
//
// with indexType (StarRocks.g4:525) USING (BITMAP | GIN | NGRAMBF | VECTOR).
// cur is INDEX.
func (p *Parser) parseAlterAddIndex(action *ast.AlterTableAction) error {
	idx := &ast.IndexDef{Loc: p.cur.Loc}
	p.advance() // consume INDEX
	name, _, err := p.parseIdentifier()
	if err != nil {
		return err
	}
	idx.Name = name
	if idx.Columns, err = p.alterIdentifierList(); err != nil {
		return err
	}
	if p.cur.Kind == kwUSING {
		p.advance()
		switch {
		case p.alterAt("BITMAP"), p.alterAt("GIN"), p.alterAt("NGRAMBF"), p.alterAt("VECTOR"):
			idx.IndexType = strings.ToUpper(p.cur.Str)
			p.advance()
		default:
			return p.syntaxErrorAtCur()
		}
		if p.cur.Kind == int('(') {
			if idx.Properties, err = p.parsePropertyList(); err != nil {
				return err
			}
		}
	}
	if p.cur.Kind == kwCOMMENT {
		p.advance()
		tok, err := p.alterExpectString()
		if err != nil {
			return err
		}
		idx.Comment = tok.Str
	}
	idx.Loc.End = p.prev.Loc.End
	action.Type = ast.AlterAddIndex
	action.Index = idx
	return nil
}

// parseAlterAddColumn reads addColumnClause (StarRocks.g4:1036) and
// addColumnsClause (StarRocks.g4:1040); ADD COLUMN has been consumed:
//
//	ADD COLUMN columnDesc (FIRST | AFTER identifier)? ((TO | IN) identifier)? properties?
//	ADD COLUMN '(' columnDesc (',' columnDesc)* ')' ((TO | IN) identifier)? properties?
func (p *Parser) parseAlterAddColumn(action *ast.AlterTableAction) error {
	action.Type = ast.AlterAddColumn
	if p.cur.Kind == int('(') {
		p.advance()
		for {
			col, err := p.parseAlterColumnDesc()
			if err != nil {
				return err
			}
			action.Columns = append(action.Columns, col)
			if p.cur.Kind != int(',') {
				break
			}
			p.advance()
		}
		if _, err := p.expect(int(')')); err != nil {
			return err
		}
		action.Column = action.Columns[0]
	} else {
		col, err := p.parseAlterColumnDesc()
		if err != nil {
			return err
		}
		action.Column = col
		action.Columns = []*ast.ColumnDef{col}
		if err := p.parseAlterColumnPosition(action); err != nil {
			return err
		}
	}
	if p.cur.Kind == kwTO || p.cur.Kind == kwIN {
		p.advance()
		name, _, err := p.parseIdentifier()
		if err != nil {
			return err
		}
		action.RollupName = name
	}
	return p.alterOptProperties(action)
}

// parseAlterColumnPosition reads (FIRST | AFTER identifier)?.
func (p *Parser) parseAlterColumnPosition(action *ast.AlterTableAction) error {
	switch p.cur.Kind {
	case kwFIRST:
		p.advance()
		action.First = true
	case kwAFTER:
		p.advance()
		name, _, err := p.parseIdentifier()
		if err != nil {
			return err
		}
		action.After = name
	}
	return nil
}

// parseAlterFromRollup reads (FROM rollupName=identifier)?.
func (p *Parser) parseAlterFromRollup(action *ast.AlterTableAction) error {
	if p.cur.Kind != kwFROM {
		return nil
	}
	p.advance()
	name, _, err := p.parseIdentifier()
	if err != nil {
		return err
	}
	action.RollupName = name
	return nil
}

// parseAlterAddPartition reads addPartitionClause (StarRocks.g4:1152):
//
//	ADD TEMPORARY? (singleRangePartition | PARTITIONS multiRangePartition) distributionDesc? properties?
//	ADD TEMPORARY? (singleItemListPartitionDesc | multiItemListPartitionDesc) distributionDesc? properties?
//
// ADD has been consumed.
func (p *Parser) parseAlterAddPartition(action *ast.AlterTableAction) error {
	action.Type = ast.AlterAddPartition
	if p.cur.Kind == kwTEMPORARY {
		p.advance()
		action.Temporary = true
	}
	var item *ast.PartitionItem
	var err error
	switch p.cur.Kind {
	case kwPARTITIONS:
		p.advance()
		item, err = p.parseAlterMultiRangePartition()
	case kwPARTITION:
		item, err = p.parseAlterPartitionDef()
	default:
		return p.syntaxErrorAtCur()
	}
	if err != nil {
		return err
	}
	action.Partition = item
	// The partition's own propertyList is the action's PartitionProps.
	action.PartitionProps, item.Properties = item.Properties, nil
	if p.cur.Kind == kwDISTRIBUTED {
		if action.PartitionDist, err = p.parseAlterDistributionDesc(); err != nil {
			return err
		}
	}
	return p.alterOptProperties(action)
}

// parseAlterPartitionDef reads singleRangePartition (StarRocks.g4:2646),
// singleItemListPartitionDesc (StarRocks.g4:2608) or
// multiItemListPartitionDesc (StarRocks.g4:2612):
//
//	PARTITION (IF NOT EXISTS)? identifier VALUES partitionKeyDesc propertyList?
//	PARTITION (IF NOT EXISTS)? identifier VALUES IN singleListPartitionValues propertyList?
//	PARTITION (IF NOT EXISTS)? identifier VALUES IN multiListPartitionValues propertyList?
//
// cur is PARTITION.
func (p *Parser) parseAlterPartitionDef() (*ast.PartitionItem, error) {
	item := &ast.PartitionItem{Loc: p.cur.Loc}
	p.advance() // consume PARTITION
	if p.cur.Kind == kwIF {
		p.advance()
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
	}
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	item.Name = name
	if _, err := p.expect(kwVALUES); err != nil {
		return nil, err
	}
	switch p.cur.Kind {
	case kwLESS:
		// StarRocks.g4:2659 partitionKeyDesc: LESS THAN (MAXVALUE | partitionValueList)
		p.advance()
		if _, err := p.expect(kwTHAN); err != nil {
			return nil, err
		}
		if p.cur.Kind == kwMAXVALUE {
			p.advance()
			item.IsMaxValue = true
		} else if item.Values, err = p.parseAlterPartitionValueList(); err != nil {
			return nil, err
		}
	case int('['):
		// StarRocks.g4:2661 '[' partitionValueList ',' partitionValueList ')'
		p.advance()
		lower, err := p.parseAlterPartitionValueList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(int(',')); err != nil {
			return nil, err
		}
		upper, err := p.parseAlterPartitionValueList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
		item.Values = append(lower, upper...)
	case kwIN:
		p.advance()
		if item.InValues, err = p.parseAlterListPartitionValues(); err != nil {
			return nil, err
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if p.cur.Kind == int('(') {
		if item.Properties, err = p.parsePropertyList(); err != nil {
			return nil, err
		}
	}
	item.Loc.End = p.prev.Loc.End
	return item, nil
}

// parseAlterPartitionValueList reads partitionValueList (StarRocks.g4:2664):
// '(' partitionValue (',' partitionValue)* ')', partitionValue being
// MAXVALUE | string.
func (p *Parser) parseAlterPartitionValueList() ([]string, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	var vals []string
	for {
		switch {
		case p.cur.Kind == kwMAXVALUE:
			vals = append(vals, "MAXVALUE")
			p.advance()
		case p.cur.Kind == tokString:
			vals = append(vals, p.cur.Str)
			p.advance()
		default:
			return nil, p.syntaxErrorAtCur()
		}
		if p.cur.Kind != int(',') {
			break
		}
		p.advance()
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return vals, nil
}

// parseAlterListPartitionValues reads singleListPartitionValues
// (StarRocks.g4:2620) or multiListPartitionValues (StarRocks.g4:2616):
//
//	'(' listPartitionValue (',' listPartitionValue)* ')'
//	'(' singleListPartitionValues (',' singleListPartitionValues)* ')'
//
// with listPartitionValue (StarRocks.g4:2629) NULL | string.
func (p *Parser) parseAlterListPartitionValues() ([][]string, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	var out [][]string
	if p.cur.Kind == int('(') {
		for {
			one, err := p.parseAlterSingleListValues()
			if err != nil {
				return nil, err
			}
			out = append(out, one)
			if p.cur.Kind != int(',') {
				break
			}
			p.advance()
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
		return out, nil
	}
	var one []string
	for {
		v, err := p.parseAlterListPartitionValue()
		if err != nil {
			return nil, err
		}
		one = append(one, v)
		if p.cur.Kind != int(',') {
			break
		}
		p.advance()
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	for _, v := range one {
		out = append(out, []string{v})
	}
	return out, nil
}

// parseAlterSingleListValues reads singleListPartitionValues.
func (p *Parser) parseAlterSingleListValues() ([]string, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	var vals []string
	for {
		v, err := p.parseAlterListPartitionValue()
		if err != nil {
			return nil, err
		}
		vals = append(vals, v)
		if p.cur.Kind != int(',') {
			break
		}
		p.advance()
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return vals, nil
}

// parseAlterListPartitionValue reads listPartitionValue: NULL | string.
func (p *Parser) parseAlterListPartitionValue() (string, error) {
	switch p.cur.Kind {
	case kwNULL:
		p.advance()
		return "NULL", nil
	case tokString:
		return p.advance().Str, nil
	}
	return "", p.syntaxErrorAtCur()
}

// parseAlterMultiRangePartition reads multiRangePartition
// (StarRocks.g4:2650):
//
//	START '(' string ')' END '(' string ')' EVERY '(' interval ')'
//	START '(' string ')' END '(' string ')' EVERY '(' INTEGER_VALUE ')'
//
// The AST builder takes an interval only with an integer value (EVERY
// (INTERVAL CURRENT_USER YEAR) is refused while EVERY (INTERVAL 1 MONTH)
// is accepted), so this reads INTERVAL INTEGER_VALUE unitIdentifier.
func (p *Parser) parseAlterMultiRangePartition() (*ast.PartitionItem, error) {
	item := &ast.PartitionItem{IsStep: true, Loc: p.cur.Loc}
	bound := func(word string) (string, error) {
		if _, err := p.alterExpectWord(word); err != nil {
			return "", err
		}
		if _, err := p.expect(int('(')); err != nil {
			return "", err
		}
		tok, err := p.alterExpectString()
		if err != nil {
			return "", err
		}
		if _, err := p.expect(int(')')); err != nil {
			return "", err
		}
		return tok.Str, nil
	}
	from, err := bound("START")
	if err != nil {
		return nil, err
	}
	to, err := bound("END")
	if err != nil {
		return nil, err
	}
	item.FromValues, item.ToValues = []string{from}, []string{to}
	if _, err := p.alterExpectWord("EVERY"); err != nil {
		return nil, err
	}
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	if p.cur.Kind == kwINTERVAL {
		p.advance()
		n, err := p.alterExpectInt()
		if err != nil {
			return nil, err
		}
		item.Interval = n.Str
		if !alterIsUnitIdentifier(p.cur) {
			return nil, p.syntaxErrorAtCur()
		}
		item.IntervalUnit = strings.ToUpper(p.advance().Str)
	} else {
		n, err := p.alterExpectInt()
		if err != nil {
			return nil, err
		}
		item.Interval = n.Str
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	item.Loc.End = p.prev.Loc.End
	return item, nil
}

// alterIsUnitIdentifier reports whether t is a unitIdentifier
// (StarRocks.g4:2776).
func alterIsUnitIdentifier(t Token) bool {
	for _, w := range []string{"YEAR", "MONTH", "WEEK", "DAY", "HOUR", "MINUTE", "SECOND", "QUARTER", "MILLISECOND", "MICROSECOND"} {
		if alterWordTok(t, w) {
			return true
		}
	}
	return false
}

// parseAlterDistributionDesc reads distributionDesc (StarRocks.g4:2681):
//
//	DISTRIBUTED BY HASH identifierList (BUCKETS INTEGER_VALUE)?
//	DISTRIBUTED BY RANDOM (BUCKETS INTEGER_VALUE)?
//
// cur is DISTRIBUTED.
func (p *Parser) parseAlterDistributionDesc() (*ast.DistributionDesc, error) {
	dd := &ast.DistributionDesc{Loc: p.cur.Loc}
	p.advance() // consume DISTRIBUTED
	if _, err := p.expect(kwBY); err != nil {
		return nil, err
	}
	switch p.cur.Kind {
	case kwHASH:
		p.advance()
		dd.Type = "HASH"
		cols, err := p.alterIdentifierList()
		if err != nil {
			return nil, err
		}
		dd.Columns = cols
	case kwRANDOM:
		p.advance()
		dd.Type = "RANDOM"
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if p.cur.Kind == kwBUCKETS {
		p.advance()
		n, err := p.alterExpectInt()
		if err != nil {
			return nil, err
		}
		dd.Buckets = int(n.Ival)
	}
	dd.Loc.End = p.prev.Loc.End
	return dd, nil
}

// parseAlterDropClause reads the DROP alternatives of alterClause. cur is
// DROP.
func (p *Parser) parseAlterDropClause(action *ast.AlterTableAction) error {
	p.advance() // consume DROP
	switch {
	case p.cur.Kind == kwINDEX:
		// StarRocks.g4:1008 dropIndexClause: DROP INDEX identifier.
		p.advance()
		name, _, err := p.parseIdentifier()
		if err != nil {
			return err
		}
		action.Type = ast.AlterDropIndex
		action.IndexName = name
		return nil
	case p.cur.Kind == kwCOLUMN:
		// StarRocks.g4:1044 dropColumnClause:
		// DROP COLUMN identifier (FROM rollupName=identifier)? properties?
		p.advance()
		name, _, err := p.parseIdentifier()
		if err != nil {
			return err
		}
		action.Type = ast.AlterDropColumn
		action.ColumnName = name
		if err := p.parseAlterFromRollup(action); err != nil {
			return err
		}
		return p.alterOptProperties(action)
	case p.alterAt("BRANCH") || p.alterAt("TAG"):
		// StarRocks.g4:1094 dropBranchClause / StarRocks.g4:1098
		// dropTagClause: DROP (BRANCH | TAG) (IF EXISTS)? identifier.
		action.Type = ast.AlterDropBranch
		if p.alterAt("TAG") {
			action.Type = ast.AlterDropTag
		}
		p.advance()
		if p.cur.Kind == kwIF {
			p.advance()
			if _, err := p.expect(kwEXISTS); err != nil {
				return err
			}
			action.IfExists = true
		}
		name, _, err := p.parseIdentifier()
		if err != nil {
			return err
		}
		action.RefName = name
		return nil
	case p.alterAt("PERSISTENT"):
		// StarRocks.g4:1146 dropPersistentIndexClause:
		// DROP PERSISTENT INDEX ON TABLETS integer_list.
		p.advance()
		if _, err := p.expect(kwINDEX); err != nil {
			return err
		}
		if _, err := p.expect(kwON); err != nil {
			return err
		}
		if _, err := p.alterExpectWord("TABLETS"); err != nil {
			return err
		}
		if _, err := p.expect(int('(')); err != nil {
			return err
		}
		for {
			n, err := p.alterExpectInt()
			if err != nil {
				return err
			}
			action.TabletIDs = append(action.TabletIDs, n.Ival)
			if p.cur.Kind != int(',') {
				break
			}
			p.advance()
		}
		if _, err := p.expect(int(')')); err != nil {
			return err
		}
		action.Type = ast.AlterDropPersistentIndex
		return nil
	case p.cur.Kind == kwTEMPORARY || p.cur.Kind == kwPARTITION || p.cur.Kind == kwPARTITIONS || p.cur.Kind == kwALL:
		return p.parseAlterDropPartition(action)
	}
	return p.syntaxErrorAtCur()
}

// parseAlterDropPartition reads dropPartitionClause (StarRocks.g4:1157);
// DROP has been consumed:
//
//	DROP TEMPORARY? (PARTITION (IF EXISTS)? identifier | PARTITIONS (IF EXISTS)? identifierList) FORCE?
//	DROP TEMPORARY? PARTITIONS (IF EXISTS)? multiRangePartition FORCE?
//	DROP TEMPORARY? PARTITIONS (IF EXISTS)? WHERE expression FORCE?
//	DROP ALL TEMPORARY PARTITIONS FORCE?
func (p *Parser) parseAlterDropPartition(action *ast.AlterTableAction) error {
	action.Type = ast.AlterDropPartition
	if p.cur.Kind == kwALL {
		p.advance()
		if _, err := p.expect(kwTEMPORARY); err != nil {
			return err
		}
		if _, err := p.expect(kwPARTITIONS); err != nil {
			return err
		}
		action.Temporary = true
		action.PartitionStar = true
	} else {
		if p.cur.Kind == kwTEMPORARY {
			p.advance()
			action.Temporary = true
		}
		plural := p.cur.Kind == kwPARTITIONS
		if !plural && p.cur.Kind != kwPARTITION {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		if p.cur.Kind == kwIF {
			p.advance()
			if _, err := p.expect(kwEXISTS); err != nil {
				return err
			}
			action.IfExists = true
		}
		var err error
		switch {
		case !plural:
			action.PartitionName, _, err = p.parseIdentifier()
		case p.cur.Kind == int('('):
			action.PartitionList, err = p.alterIdentifierList()
		case p.cur.Kind == kwWHERE:
			p.advance()
			action.Where, err = p.parseExpr()
		default:
			action.Partition, err = p.parseAlterMultiRangePartition()
		}
		if err != nil {
			return err
		}
	}
	if p.cur.Kind == kwFORCE {
		p.advance()
		action.Force = true
	}
	return nil
}

// parseAlterModifyClause reads the MODIFY alternatives of alterClause. cur
// is MODIFY.
func (p *Parser) parseAlterModifyClause(action *ast.AlterTableAction) error {
	p.advance() // consume MODIFY
	switch p.cur.Kind {
	case kwCOLUMN:
		p.advance() // consume COLUMN
		// StarRocks.g4:1076 addFieldClause / StarRocks.g4:1080
		// dropFieldClause: MODIFY COLUMN identifier (ADD | DROP) FIELD ...
		if isIdentifierToken(p.cur.Kind) && (p.peekNext().Kind == kwADD || p.peekNext().Kind == kwDROP) {
			return p.parseAlterModifyColumnField(action)
		}
		// StarRocks.g4:1048 modifyColumnClause:
		// MODIFY COLUMN columnDesc (FIRST | AFTER identifier)? (FROM rollupName=identifier)? properties?
		col, err := p.parseAlterColumnDesc()
		if err != nil {
			return err
		}
		action.Type = ast.AlterModifyColumn
		action.Column = col
		action.Columns = []*ast.ColumnDef{col}
		if err := p.parseAlterColumnPosition(action); err != nil {
			return err
		}
		if err := p.parseAlterFromRollup(action); err != nil {
			return err
		}
		return p.alterOptProperties(action)
	case kwPARTITION:
		// StarRocks.g4:1168 modifyPartitionClause:
		//	MODIFY PARTITION (identifier | identifierList | '(' ASTERISK_SYMBOL ')') SET propertyList
		//	MODIFY PARTITION distributionDesc
		p.advance()
		action.Type = ast.AlterModifyPartition
		if p.cur.Kind == kwDISTRIBUTED {
			dist, err := p.parseAlterDistributionDesc()
			if err != nil {
				return err
			}
			action.Distribution = dist
			return nil
		}
		if p.cur.Kind == int('(') && p.peekNext().Kind == int('*') {
			p.advance()
			p.advance()
			if _, err := p.expect(int(')')); err != nil {
				return err
			}
			action.PartitionStar = true
		} else if p.cur.Kind == int('(') {
			names, err := p.alterIdentifierList()
			if err != nil {
				return err
			}
			action.PartitionList = names
		} else {
			name, _, err := p.parseIdentifier()
			if err != nil {
				return err
			}
			action.PartitionName = name
		}
		if _, err := p.expect(kwSET); err != nil {
			return err
		}
		props, err := p.parsePropertyList()
		if err != nil {
			return err
		}
		action.Properties = props
		return nil
	}
	return p.syntaxErrorAtCur()
}

// parseAlterModifyColumnField reads addFieldClause (StarRocks.g4:1076) and
// dropFieldClause (StarRocks.g4:1080); MODIFY COLUMN has been consumed:
//
//	MODIFY COLUMN identifier ADD FIELD subfieldDesc (FIRST | AFTER identifier)? properties?
//	MODIFY COLUMN identifier DROP FIELD nestedFieldName properties?
func (p *Parser) parseAlterModifyColumnField(action *ast.AlterTableAction) error {
	name, _, err := p.parseIdentifier()
	if err != nil {
		return err
	}
	action.ColumnName = name
	if p.cur.Kind == kwADD {
		p.advance()
		if _, err := p.expect(kwFIELD); err != nil {
			return err
		}
		action.Type = ast.AlterModifyColumnAddField
		// StarRocks.g4:2800 subfieldDesc: (identifier | nestedFieldName) type
		if action.FieldName, err = p.parseFieldPath(); err != nil {
			return err
		}
		if action.FieldType, err = p.parseAlterType(); err != nil {
			return err
		}
		if err := p.parseAlterColumnPosition(action); err != nil {
			return err
		}
		return p.alterOptProperties(action)
	}
	if _, err := p.expect(kwDROP); err != nil {
		return err
	}
	if _, err := p.expect(kwFIELD); err != nil {
		return err
	}
	action.Type = ast.AlterModifyColumnDropField
	if action.FieldName, err = p.parseFieldPath(); err != nil {
		return err
	}
	return p.alterOptProperties(action)
}

// parseAlterRenameClause reads the RENAME alternatives of alterClause. cur
// is RENAME.
func (p *Parser) parseAlterRenameClause(action *ast.AlterTableAction) error {
	p.advance() // consume RENAME
	var err error
	switch {
	case p.cur.Kind == kwCOLUMN:
		// StarRocks.g4:1052 columnRenameClause:
		// RENAME COLUMN oldColumn=identifier TO newColumn=identifier
		p.advance()
		action.Type = ast.AlterRenameColumn
		if action.ColumnName, _, err = p.parseIdentifier(); err != nil {
			return err
		}
		if _, err = p.expect(kwTO); err != nil {
			return err
		}
		action.NewName, _, err = p.parseIdentifier()
		return err
	case p.cur.Kind == kwROLLUP || p.cur.Kind == kwPARTITION:
		// StarRocks.g4:1060 rollupRenameClause / StarRocks.g4:1177
		// partitionRenameClause: RENAME (ROLLUP | PARTITION) identifier identifier
		action.Type = ast.AlterRenameRollup
		if p.cur.Kind == kwPARTITION {
			action.Type = ast.AlterRenamePartition
		}
		p.advance()
		if action.ColumnName, _, err = p.parseIdentifier(); err != nil {
			return err
		}
		if action.Type == ast.AlterRenameRollup {
			action.RollupName = action.ColumnName
		} else {
			action.PartitionName = action.ColumnName
		}
		action.NewName, _, err = p.parseIdentifier()
		return err
	}
	// StarRocks.g4:1012 tableRenameClause: RENAME identifier
	name, loc, err := p.parseIdentifier()
	if err != nil {
		return err
	}
	action.Type = ast.AlterRenameTable
	action.NewName = name
	action.NewTableName = &ast.ObjectName{Parts: []string{name}, Loc: loc}
	return nil
}

// parseAlterCompactClause reads compactionClause (StarRocks.g4:1064):
// (BASE | CUMULATIVE)? COMPACT (identifier | identifierList)?.
func (p *Parser) parseAlterCompactClause(action *ast.AlterTableAction) error {
	action.Type = ast.AlterCompact
	if p.alterAt("BASE") || p.alterAt("CUMULATIVE") {
		action.CompactType = strings.ToUpper(p.advance().Str)
	}
	if _, err := p.alterExpectWord("COMPACT"); err != nil {
		return err
	}
	if p.cur.Kind == int('(') {
		names, err := p.alterIdentifierList()
		if err != nil {
			return err
		}
		action.PartitionList = names
	} else if isIdentifierToken(p.cur.Kind) {
		name, _, err := p.parseIdentifier()
		if err != nil {
			return err
		}
		action.PartitionName = name
	}
	return nil
}

// parseAlterCreateOrReplaceClause reads createOrReplaceTagClause
// (StarRocks.g4:1084) and createOrReplaceBranchClause (StarRocks.g4:1089),
// or, for REPLACE followed by partitionNames, replacePartitionClause
// (StarRocks.g4:1173):
//
//	(CREATE OR)? REPLACE (TAG | BRANCH) identifier (tagOptions | branchOptions)
//	CREATE (TAG | BRANCH) (IF NOT EXISTS)? identifier (tagOptions | branchOptions)
//	REPLACE partitionNames WITH partitionNames properties?
func (p *Parser) parseAlterCreateOrReplaceClause(action *ast.AlterTableAction) error {
	if p.cur.Kind == kwREPLACE && !p.alterNextAt("TAG") && !p.alterNextAt("BRANCH") {
		p.advance() // consume REPLACE
		action.Type = ast.AlterReplacePartition
		if err := p.parseAlterPartitionNames(action); err != nil {
			return err
		}
		if _, err := p.expect(kwWITH); err != nil {
			return err
		}
		with := &ast.AlterTableAction{}
		if err := p.parseAlterPartitionNames(with); err != nil {
			return err
		}
		action.TempPartitionList = with.PartitionList
		if with.PartitionName != "" {
			action.TempPartitionList = []string{with.PartitionName}
		}
		return p.alterOptProperties(action)
	}
	create := p.cur.Kind == kwCREATE
	p.advance() // consume CREATE or REPLACE
	if create && p.cur.Kind == kwOR {
		p.advance()
		if _, err := p.expect(kwREPLACE); err != nil {
			return err
		}
		create = false
	}
	action.OrReplace = !create
	switch {
	case p.alterAt("TAG"):
		action.Type = ast.AlterCreateTag
	case p.alterAt("BRANCH"):
		action.Type = ast.AlterCreateBranch
	default:
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if create && p.cur.Kind == kwIF {
		p.advance()
		if _, err := p.expect(kwNOT); err != nil {
			return err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return err
		}
		action.IfNotExists = true
	}
	name, _, err := p.parseIdentifier()
	if err != nil {
		return err
	}
	action.RefName = name
	return p.parseAlterRefOptions(action, action.Type == ast.AlterCreateBranch)
}

// parseAlterRefOptions reads tagOptions (StarRocks.g4:1106) or branchOptions
// (StarRocks.g4:1110):
//
//	(AS OF VERSION snapshotId)? refRetain? snapshotRetention?   -- branch
//	(AS OF VERSION snapshotId)? refRetain?                       -- tag
//
// snapshotId and number are StarRocks.g4:2906 number, but the AST builder
// reads each as a long, so only an INTEGER_VALUE is accepted (RETAIN 1.5
// DAYS is refused).
func (p *Parser) parseAlterRefOptions(action *ast.AlterTableAction, branch bool) error {
	var err error
	if p.cur.Kind == kwAS {
		p.advance()
		if _, err := p.alterExpectWord("OF"); err != nil {
			return err
		}
		if _, err := p.alterExpectWord("VERSION"); err != nil {
			return err
		}
		tok, err := p.alterExpectInt()
		if err != nil {
			return err
		}
		action.AsOfVersion = tok.Str
	}
	if p.alterAt("RETAIN") {
		// StarRocks.g4:1120 refRetain: RETAIN number timeUnit
		p.advance()
		if action.RetainTime, err = p.parseAlterTimeValue(); err != nil {
			return err
		}
	}
	if branch && p.cur.Kind == kwWITH {
		// StarRocks.g4:1114 snapshotRetention:
		//	WITH SNAPSHOT RETENTION minSnapshotsToKeep
		//	WITH SNAPSHOT RETENTION maxSnapshotAge
		//	WITH SNAPSHOT RETENTION minSnapshotsToKeep maxSnapshotAge
		p.advance()
		if _, err := p.alterExpectWord("SNAPSHOT"); err != nil {
			return err
		}
		if _, err := p.alterExpectWord("RETENTION"); err != nil {
			return err
		}
		if p.cur.Kind == tokInt && alterWordTok(p.peekNext(), "SNAPSHOTS") {
			action.MinSnapshots = p.advance().Str
			p.advance() // consume SNAPSHOTS
			if p.cur.Kind == tokInt {
				action.MaxSnapshotAge, err = p.parseAlterTimeValue()
			}
			return err
		}
		action.MaxSnapshotAge, err = p.parseAlterTimeValue()
		return err
	}
	return nil
}

// parseAlterTimeValue reads number timeUnit (StarRocks.g4:1124
// maxSnapshotAge, timeUnit StarRocks.g4:1136: DAYS | HOURS | MINUTES) and
// returns it as written, "7 DAYS".
func (p *Parser) parseAlterTimeValue() (string, error) {
	n, err := p.alterExpectInt()
	if err != nil {
		return "", err
	}
	if p.alterAt("DAYS") || p.alterAt("HOURS") || p.alterAt("MINUTES") {
		return n.Str + " " + strings.ToUpper(p.advance().Str), nil
	}
	return "", p.syntaxErrorAtCur()
}

// parseAlterTableOperation reads tableOperationClause (StarRocks.g4:1102):
// EXECUTE functionCall. cur is EXECUTE.
func (p *Parser) parseAlterTableOperation(action *ast.AlterTableAction) error {
	p.advance() // consume EXECUTE
	if !isIdentifierToken(p.cur.Kind) || p.peekNext().Kind != int('(') {
		return p.syntaxErrorAtCur()
	}
	call, err := p.parseExpr()
	if err != nil {
		return err
	}
	if _, ok := call.(*ast.FuncCallExpr); !ok {
		return &ParseError{Position: action.Loc.Start, End: p.prev.Loc.End, Message: "EXECUTE needs a function call"}
	}
	action.Type = ast.AlterTableOperation
	action.Call = call
	return nil
}

// parseAlterPartitionNames reads partitionNames (StarRocks.g4:2268) into
// action.Temporary, PartitionName (one name) or PartitionList:
//
//	TEMPORARY? (PARTITION | PARTITIONS) '(' identifierOrString (',' identifierOrString)* ')'
//	TEMPORARY? (PARTITION | PARTITIONS) identifierOrString
//	keyPartitions: PARTITION '(' keyPartition (',' keyPartition)* ')'
func (p *Parser) parseAlterPartitionNames(action *ast.AlterTableAction) error {
	temporary := false
	if p.cur.Kind == kwTEMPORARY {
		p.advance()
		temporary = true
		action.Temporary = true
	}
	single := p.cur.Kind == kwPARTITION
	if !single && p.cur.Kind != kwPARTITIONS {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Kind != int('(') {
		name, _, err := p.parseIdentifierOrString()
		if err != nil {
			return err
		}
		action.PartitionName = name
		return nil
	}
	p.advance() // consume '('
	// StarRocks.g4:2668 keyPartition: identifier '=' literalExpression
	keyed := single && !temporary && isIdentifierToken(p.cur.Kind) && p.peekNext().Kind == int('=')
	for {
		if keyed {
			start := p.cur.Loc.Start
			if _, _, err := p.parseIdentifier(); err != nil {
				return err
			}
			if _, err := p.expect(int('=')); err != nil {
				return err
			}
			if err := p.parseAlterLiteralExpression(); err != nil {
				return err
			}
			action.PartitionList = append(action.PartitionList, p.input[start-p.baseOffset:p.prev.Loc.End-p.baseOffset])
		} else {
			name, _, err := p.parseIdentifierOrString()
			if err != nil {
				return err
			}
			action.PartitionList = append(action.PartitionList, name)
		}
		if p.cur.Kind != int(',') {
			break
		}
		p.advance()
	}
	_, err := p.expect(int(')'))
	return err
}

// parseAlterLiteralExpression reads literalExpression (StarRocks.g4:2428):
// NULL | booleanValue | number | (DATE | DATETIME) string | string |
// interval | unitBoundary | binary | PARAMETER.
func (p *Parser) parseAlterLiteralExpression() error {
	switch {
	case p.cur.Kind == kwNULL, p.cur.Kind == kwTRUE, p.cur.Kind == kwFALSE,
		p.cur.Kind == tokInt, p.cur.Kind == tokFloat, p.cur.Kind == tokString,
		p.cur.Kind == tokHexLiteral, p.cur.Kind == tokBitLiteral, p.cur.Kind == tokPlaceholder,
		p.alterAt("FLOOR"), p.alterAt("CEIL"):
		p.advance()
		return nil
	case p.cur.Kind == kwDATE || p.cur.Kind == kwDATETIME:
		p.advance()
		_, err := p.alterExpectString()
		return err
	case p.cur.Kind == kwINTERVAL:
		// interval: INTERVAL expression unitIdentifier (StarRocks.g4:2776)
		_, err := p.parseExpr()
		return err
	}
	return p.syntaxErrorAtCur()
}

// parseAlterOrderByClause reads reorderColumnsClause (StarRocks.g4:1056):
// ORDER BY identifierList (FROM rollupName=identifier)? properties?, or an
// optimizeClause whose first part is orderByDesc. cur is ORDER.
func (p *Parser) parseAlterOrderByClause(action *ast.AlterTableAction) error {
	p.advance() // consume ORDER
	p.advance() // consume BY
	cols, err := p.alterIdentifierList()
	if err != nil {
		return err
	}
	action.Type = ast.AlterOrderBy
	action.OrderByColumns = cols
	if p.cur.Kind == kwDISTRIBUTED {
		action.Type = ast.AlterOptimize
		action.Distribution, err = p.parseAlterDistributionDesc()
		return err
	}
	if err := p.parseAlterFromRollup(action); err != nil {
		return err
	}
	return p.alterOptProperties(action)
}

// parseAlterOptimizeClause reads optimizeClause (StarRocks.g4:1028):
//
//	partitionNames? keyDesc? partitionDesc? orderByDesc? distributionDesc?
//
// The AST builder refuses a partitionDesc here (every ALTER SYSTEM
// PARTITION BY ... row is a reject, whatever follows BY), so a PARTITION BY
// is a syntax error.
func (p *Parser) parseAlterOptimizeClause(action *ast.AlterTableAction) error {
	action.Type = ast.AlterOptimize
	start := p.cur.Loc.Start
	if p.cur.Kind == kwTEMPORARY || ((p.cur.Kind == kwPARTITION || p.cur.Kind == kwPARTITIONS) && p.peekNext().Kind != kwBY) {
		if err := p.parseAlterPartitionNames(action); err != nil {
			return err
		}
	}
	// StarRocks.g4:432 keyDesc: (AGGREGATE | UNIQUE | PRIMARY | DUPLICATE) KEY identifierList
	switch p.cur.Kind {
	case kwAGGREGATE, kwUNIQUE, kwPRIMARY, kwDUPLICATE:
		kd := &ast.KeyDesc{Type: strings.ToUpper(p.cur.Str), Loc: p.cur.Loc}
		p.advance()
		if _, err := p.expect(kwKEY); err != nil {
			return err
		}
		cols, err := p.alterIdentifierList()
		if err != nil {
			return err
		}
		kd.Columns = cols
		kd.Loc.End = p.prev.Loc.End
		action.KeyDesc = kd
	}
	if p.cur.Kind == kwORDER {
		// StarRocks.g4:436 orderByDesc: ORDER BY identifierList
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return err
		}
		cols, err := p.alterIdentifierList()
		if err != nil {
			return err
		}
		action.OrderByColumns = cols
	}
	if p.cur.Kind == kwDISTRIBUTED {
		dist, err := p.parseAlterDistributionDesc()
		if err != nil {
			return err
		}
		action.Distribution = dist
	}
	if p.cur.Loc.Start == start {
		// Nothing of an optimizeClause, and no other clause starts here.
		return p.syntaxErrorAtCur()
	}
	return nil
}
