package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// This file reads ALTER TABLE by the grammar. Doris's ANTLR parser takes
// ALTER TABLE as an unsupportedAlterStatement (DorisParser.g4:556) and hands
// the text to its legacy CUP parser (DorisLegacyParser.g4:176 alter_stmt),
// so a statement must satisfy both: each reader below follows the narrower
// of the two rules, cited as DorisParser.g4 (Nereids) and
// DorisLegacyParser.g4 (CUP). The compat rows (testdata/compat.tsv) are the
// evidence where the two disagree.

// alterIsWord reports whether cur is the bare word w, a keyword or an
// unquoted identifier spelled w.
func (p *Parser) alterIsWord(w string) bool {
	return (p.cur.Kind == tokIdent || p.cur.Kind >= 700) && strings.EqualFold(p.cur.Str, w)
}

// alterExpectString consumes a string literal.
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

// alterIdentList reads identifierList (DorisParser.g4 identifierList;
// DorisLegacyParser.g4 '(' ident_list ')'): '(' ident (',' ident)* ')'.
func (p *Parser) alterIdentList() ([]string, error) {
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

// alterParenIdentList is alterIdentList returning the closing ')' too.
func (p *Parser) alterParenIdentList() ([]string, ast.Loc, error) {
	names, err := p.alterIdentList()
	if err != nil {
		return nil, ast.Loc{}, err
	}
	return names, p.prev.Loc, nil
}

// alterOptProperties reads an optional PROPERTIES '(' key_value_map ')'
// (DorisLegacyParser.g4:1175 opt_properties, whose key_value_map takes
// string keys and values only).
func (p *Parser) alterOptProperties() ([]*ast.Property, error) {
	if p.cur.Kind != kwPROPERTIES {
		return nil, nil
	}
	return p.parseProperties()
}

// alterOptIfExists reads an optional IF EXISTS.
func (p *Parser) alterOptIfExists() (bool, error) {
	if p.cur.Kind != kwIF {
		return false, nil
	}
	p.advance()
	if _, err := p.expect(kwEXISTS); err != nil {
		return false, err
	}
	return true, nil
}

// alterOptIfNotExists reads an optional IF NOT EXISTS.
func (p *Parser) alterOptIfNotExists() (bool, error) {
	if p.cur.Kind != kwIF {
		return false, nil
	}
	p.advance()
	if _, err := p.expect(kwNOT); err != nil {
		return false, err
	}
	if _, err := p.expect(kwEXISTS); err != nil {
		return false, err
	}
	return true, nil
}

// parseAlterTable parses ALTER TABLE. The ALTER keyword has been consumed;
// cur is TABLE. Parse reads it by the grammar (parseAlterTableStrict); a
// best-effort parse that the grammar refuses falls back to the loose reader.
func (p *Parser) parseAlterTable() (ast.Node, error) {
	startLoc := p.prev.Loc
	cp := p.save()
	node, err := p.parseAlterTableStrict(startLoc)
	if err == nil || p.strictTrailing {
		return node, err
	}
	p.restore(cp)
	return p.parseAlterTableLoose()
}

// parseAlterTableStrict reads the ALTER TABLE forms of
// DorisLegacyParser.g4:176 alter_stmt, each also an alternative of
// DorisParser.g4:556 unsupportedAlterStatement or :750
// unsupportedStatsStatement, plus ADD / DROP CONSTRAINT
// (DorisParser.g4:127, which Doris runs itself). cur is TABLE.
func (p *Parser) parseAlterTableStrict(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume TABLE

	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt := &ast.AlterTableStmt{Name: name}

	switch {
	case p.cur.Kind == kwADD && p.peekNext().Kind == kwCONSTRAINT:
		p.advance() // consume ADD
		return p.parseAddConstraint(startLoc, name)
	case p.cur.Kind == kwDROP && p.peekNext().Kind == kwCONSTRAINT:
		return p.parseDropConstraint(startLoc, name)
	case p.cur.Kind == kwADD && p.peekNext().Kind == kwROLLUP:
		err = p.parseAlterAddRollupList(stmt)
	case p.cur.Kind == kwDROP && p.peekNext().Kind == kwROLLUP:
		err = p.parseAlterDropRollupList(stmt)
	case p.cur.Kind == kwSET:
		err = p.parseAlterTableSet(stmt)
	case p.cur.Kind == kwINDEX || p.alterAtColumnStats():
		err = p.parseAlterColumnStats(stmt)
	default:
		err = p.parseAlterClauseList(stmt)
	}
	if err != nil {
		return nil, err
	}
	stmt.Loc = startLoc.Merge(p.prev.Loc)
	return stmt, nil
}

// parseAlterClauseList reads alter_table_clause (',' alter_table_clause)*
// (DorisLegacyParser.g4:272; DorisParser.g4:558).
func (p *Parser) parseAlterClauseList(stmt *ast.AlterTableStmt) error {
	for {
		action, err := p.parseAlterTableClause()
		if err != nil {
			return err
		}
		stmt.Actions = append(stmt.Actions, action)
		if p.cur.Kind != int(',') {
			return nil
		}
		p.advance() // consume ','
	}
}

// parseAlterAddRollupList reads ADD ROLLUP add_rollup_clause (','
// add_rollup_clause)* (DorisLegacyParser.g4:177; DorisParser.g4:560), each
// clause (DorisLegacyParser.g4:250; DorisParser.g4:621):
//
//	ident '(' ident_list ')' [DUPLICATE KEY '(' ident_list ')'] [FROM ident] [PROPERTIES (...)]
func (p *Parser) parseAlterAddRollupList(stmt *ast.AlterTableStmt) error {
	p.advance() // consume ADD
	p.advance() // consume ROLLUP
	for {
		action := &ast.AlterTableAction{Type: ast.AlterAddRollup, Loc: p.cur.Loc}
		rd := &ast.RollupDef{Loc: p.cur.Loc}
		var err error
		if rd.Name, _, err = p.parseIdentifier(); err != nil {
			return err
		}
		if rd.Columns, err = p.alterIdentList(); err != nil {
			return err
		}
		if p.cur.Kind == kwDUPLICATE {
			p.advance()
			if _, err := p.expect(kwKEY); err != nil {
				return err
			}
			if rd.DupKeys, err = p.alterIdentList(); err != nil {
				return err
			}
		}
		if p.cur.Kind == kwFROM {
			p.advance()
			if action.RollupName, _, err = p.parseIdentifier(); err != nil {
				return err
			}
		}
		if rd.Properties, err = p.alterOptProperties(); err != nil {
			return err
		}
		rd.Loc.End = p.prev.Loc.End
		action.Rollup = rd
		action.Loc.End = rd.Loc.End
		stmt.Actions = append(stmt.Actions, action)
		if p.cur.Kind != int(',') {
			return nil
		}
		p.advance() // consume ','
	}
}

// parseAlterDropRollupList reads DROP ROLLUP drop_rollup_clause (','
// drop_rollup_clause)* (DorisLegacyParser.g4:178; DorisParser.g4:562), each
// clause ident [PROPERTIES (...)] (DorisLegacyParser.g4:261;
// DorisParser.g4:617).
func (p *Parser) parseAlterDropRollupList(stmt *ast.AlterTableStmt) error {
	p.advance() // consume DROP
	p.advance() // consume ROLLUP
	for {
		action := &ast.AlterTableAction{Type: ast.AlterDropRollup, Loc: p.cur.Loc}
		var err error
		if action.RollupName, _, err = p.parseIdentifier(); err != nil {
			return err
		}
		if action.Properties, err = p.alterOptProperties(); err != nil {
			return err
		}
		action.Loc.End = p.prev.Loc.End
		stmt.Actions = append(stmt.Actions, action)
		if p.cur.Kind != int(',') {
			return nil
		}
		p.advance() // consume ','
	}
}

// parseAlterTableSet reads the two SET forms of alter_stmt, neither of which
// joins a clause list:
//
//	SET '(' key_value_map ')'                          (DorisLegacyParser.g4:199; DorisParser.g4:583)
//	SET STATS '(' key_value_map ')' opt_partition_names (DorisLegacyParser.g4:196; DorisParser.g4:756)
func (p *Parser) parseAlterTableSet(stmt *ast.AlterTableStmt) error {
	action := &ast.AlterTableAction{Type: ast.AlterSetProperties, Loc: p.cur.Loc}
	p.advance() // consume SET
	if p.cur.Kind == kwSTATS {
		p.advance()
		action.Type = ast.AlterSetStats
		if err := p.parseAlterStatsBody(action); err != nil {
			return err
		}
	} else {
		props, err := p.parseParenProperties()
		if err != nil {
			return err
		}
		action.Properties = props
	}
	action.Loc.End = p.prev.Loc.End
	stmt.Actions = append(stmt.Actions, action)
	return nil
}

// alterAtColumnStats reports whether cur begins MODIFY COLUMN ident SET,
// the column statistics form, rather than a modify-column clause.
func (p *Parser) alterAtColumnStats() bool {
	if p.cur.Kind != kwMODIFY || p.peekNext().Kind != kwCOLUMN {
		return false
	}
	cp := p.save()
	defer p.restore(cp)
	p.advance() // MODIFY
	p.advance() // COLUMN
	if _, _, err := p.parseIdentifier(); err != nil {
		return false
	}
	return p.cur.Kind == kwSET
}

// parseAlterColumnStats reads (DorisLegacyParser.g4:197-198;
// DorisParser.g4:758 alterColumnStats):
//
//	[INDEX ident] MODIFY COLUMN ident SET STATS '(' key_value_map ')' opt_partition_names
func (p *Parser) parseAlterColumnStats(stmt *ast.AlterTableStmt) error {
	action := &ast.AlterTableAction{Type: ast.AlterSetStats, Loc: p.cur.Loc}
	var err error
	if p.cur.Kind == kwINDEX {
		p.advance()
		if action.IndexName, _, err = p.parseIdentifier(); err != nil {
			return err
		}
	}
	for _, kw := range []int{kwMODIFY, kwCOLUMN} {
		if _, err := p.expect(kw); err != nil {
			return err
		}
	}
	if action.ColumnName, _, err = p.parseIdentifier(); err != nil {
		return err
	}
	for _, kw := range []int{kwSET, kwSTATS} {
		if _, err := p.expect(kw); err != nil {
			return err
		}
	}
	if err := p.parseAlterStatsBody(action); err != nil {
		return err
	}
	action.Loc.End = p.prev.Loc.End
	stmt.Actions = append(stmt.Actions, action)
	return nil
}

// parseAlterStatsBody reads '(' key_value_map ')' opt_partition_names after
// SET STATS.
func (p *Parser) parseAlterStatsBody(action *ast.AlterTableAction) error {
	props, err := p.parseParenProperties()
	if err != nil {
		return err
	}
	action.Properties = props
	if p.alterAtPartitionSpec() {
		temp, names, star, err := p.parseAlterPartitionSpec()
		if err != nil {
			return err
		}
		action.Temporary, action.PartitionList, action.PartitionStar = temp, names, star
	}
	return nil
}

// alterAtPartitionSpec reports whether cur begins a partition list.
func (p *Parser) alterAtPartitionSpec() bool {
	switch p.cur.Kind {
	case kwPARTITION, kwPARTITIONS:
		return true
	case kwTEMPORARY:
		next := p.peekNext().Kind
		return next == kwPARTITION || next == kwPARTITIONS
	}
	return false
}

// parseAlterPartitionSpec reads the partition lists that both
// DorisLegacyParser.g4:2294 partition_names and DorisParser.g4 partitionSpec
// derive (CUP's PARTITIONS WITH RECENT n is not in partitionSpec):
//
//	[TEMPORARY] (PARTITION | PARTITIONS) '(' ident_list ')'
//	[TEMPORARY] PARTITION ident
//	(PARTITION | PARTITIONS) '(' '*' ')'
func (p *Parser) parseAlterPartitionSpec() (temporary bool, names []string, star bool, err error) {
	if p.cur.Kind == kwTEMPORARY {
		temporary = true
		p.advance()
	}
	plural := p.cur.Kind == kwPARTITIONS
	if !plural && p.cur.Kind != kwPARTITION {
		return false, nil, false, p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Kind == int('(') {
		if !temporary && p.peekNext().Kind == int('*') {
			p.advance() // '('
			p.advance() // '*'
			if _, err := p.expect(int(')')); err != nil {
				return false, nil, false, err
			}
			return false, nil, true, nil
		}
		names, err = p.alterIdentList()
		return temporary, names, false, err
	}
	if plural {
		return false, nil, false, p.syntaxErrorAtCur()
	}
	name, _, err := p.parseIdentifier()
	if err != nil {
		return false, nil, false, err
	}
	return temporary, []string{name}, false, nil
}

// parseAlterTableClause reads one alter_table_clause (DorisLegacyParser.g4:304;
// DorisParser.g4:627 alterTableClause). cur is the clause's first token.
func (p *Parser) parseAlterTableClause() (*ast.AlterTableAction, error) {
	action := &ast.AlterTableAction{Loc: p.cur.Loc}
	var err error
	switch p.cur.Kind {
	case kwADD:
		err = p.parseAlterAddClause(action)
	case kwDROP:
		err = p.parseAlterDropClause(action)
	case kwMODIFY:
		err = p.parseAlterModifyClause(action)
	case kwORDER:
		// ORDER BY '(' ident_list ')' opt_from_rollup opt_properties
		// (DorisLegacyParser.g4:309; DorisParser.g4:635 reorderColumnsClause)
		p.advance()
		if _, err = p.expect(kwBY); err == nil {
			action.Type = ast.AlterOrderBy
			if action.OrderByColumns, err = p.alterIdentList(); err == nil {
				err = p.parseAlterRollupAndProps(action, false)
			}
		}
	case kwREPLACE:
		err = p.parseAlterReplaceClause(action)
	case kwRENAME:
		err = p.parseAlterRenameClause(action)
	case kwENABLE:
		// ENABLE FEATURE STRING_LITERAL ['WITH' properties]
		// (DorisLegacyParser.g4:323; DorisParser.g4:656 enableFeatureClause)
		p.advance()
		if _, err = p.expect(kwFEATURE); err == nil {
			var tok Token
			if tok, err = p.alterExpectString(); err == nil {
				action.Type = ast.AlterEnableFeature
				action.FeatureName = tok.Str
				if p.cur.Kind == kwWITH {
					p.advance()
					if p.cur.Kind != kwPROPERTIES {
						err = p.syntaxErrorAtCur()
					} else {
						action.Properties, err = p.parseProperties()
					}
				}
			}
		}
	case kwCREATE:
		err = p.parseAlterCreateRefClause(action)
	default:
		err = p.syntaxErrorAtCur()
	}
	if err != nil {
		return nil, err
	}
	action.Loc.End = p.prev.Loc.End
	return action, nil
}

// parseAlterRollupAndProps reads the tail shared by the column clauses:
// opt_to_rollup (TO | IN ident) when to is set, else opt_from_rollup
// (FROM ident), then opt_properties (DorisLegacyParser.g4:278, :285).
func (p *Parser) parseAlterRollupAndProps(action *ast.AlterTableAction, to bool) error {
	var err error
	if (to && (p.cur.Kind == kwTO || p.cur.Kind == kwIN)) || (!to && p.cur.Kind == kwFROM) {
		p.advance()
		if action.RollupName, _, err = p.parseIdentifier(); err != nil {
			return err
		}
	}
	action.Properties, err = p.alterOptProperties()
	return err
}

// parseAlterColumnPosition reads opt_col_pos: FIRST | AFTER ident
// (DorisLegacyParser.g4:291; DorisParser.g4 columnPosition).
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

// parseAlterAddClause reads the ADD alternatives of alter_table_clause.
// cur is ADD.
func (p *Parser) parseAlterAddClause(action *ast.AlterTableAction) error {
	p.advance() // consume ADD
	switch p.cur.Kind {
	case kwCOLUMN:
		p.advance()
		return p.parseAlterAddColumn(action)
	case kwINDEX:
		return p.parseAlterAddIndex(action)
	case kwTEMPORARY, kwPARTITION, kwPARTITIONS:
		if p.cur.Kind == kwTEMPORARY {
			action.Temporary = true
			p.advance()
		}
		if p.cur.Kind == kwPARTITIONS {
			return p.parseAlterAddPartitions(action)
		}
		if _, err := p.expect(kwPARTITION); err != nil {
			return err
		}
		return p.parseAlterAddPartition(action)
	}
	return p.syntaxErrorAtCur()
}

// parseAlterAddColumn reads, after ADD COLUMN (DorisLegacyParser.g4:305-306;
// DorisParser.g4:628-630 addColumnClause, addColumnsClause):
//
//	column_definition opt_col_pos opt_to_rollup opt_properties
//	'(' column_definition_list ')' opt_to_rollup opt_properties
func (p *Parser) parseAlterAddColumn(action *ast.AlterTableAction) error {
	action.Type = ast.AlterAddColumn
	if p.cur.Kind == int('(') {
		p.advance()
		for {
			col, err := p.parseAlterColumnDef()
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
		return p.parseAlterRollupAndProps(action, true)
	}
	col, err := p.parseAlterColumnDef()
	if err != nil {
		return err
	}
	action.Column = col
	action.Columns = []*ast.ColumnDef{col}
	if err := p.parseAlterColumnPosition(action); err != nil {
		return err
	}
	return p.parseAlterRollupAndProps(action, true)
}

// parseAlterColumnDef reads a column definition by what both grammars derive:
// DorisLegacyParser.g4 column_definition and DorisParser.g4 columnDef.
//
//	ident type [KEY] [agg_type] [NULL | NOT NULL] [AUTO_INCREMENT ['(' INTEGER ')']]
//	    [DEFAULT default_value] [COMMENT STRING]
//	ident type [KEY] [GENERATED ALWAYS] AS '(' expr ')' [NULL | NOT NULL] [COMMENT STRING]
//
// CUP has no ON UPDATE and no generated column with an aggregate type,
// AUTO_INCREMENT or DEFAULT; DorisParser.g4 has no generated column without
// ALWAYS after GENERATED.
func (p *Parser) parseAlterColumnDef() (*ast.ColumnDef, error) {
	startLoc := p.cur.Loc
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	typ, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	col := &ast.ColumnDef{Name: name, Type: typ, Loc: startLoc}
	if p.cur.Kind == kwKEY {
		p.advance()
		col.IsKey = true
	}
	if p.cur.Kind == kwGENERATED || p.cur.Kind == kwAS {
		if p.cur.Kind == kwGENERATED {
			p.advance()
			if _, err := p.expect(kwALWAYS); err != nil {
				return nil, err
			}
		}
		if _, err := p.expect(kwAS); err != nil {
			return nil, err
		}
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		if col.Generated, err = p.parseExpr(); err != nil {
			return nil, err
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
		if err := p.parseAlterColumnNullable(col); err != nil {
			return nil, err
		}
	} else {
		switch p.cur.Kind {
		case kwMAX, kwMIN, kwSUM, kwREPLACE, kwREPLACE_IF_NOT_NULL,
			kwHLL_UNION, kwBITMAP_UNION, kwQUANTILE_UNION, kwGENERIC:
			col.AggType = strings.ToUpper(p.cur.Str)
			p.advance()
		}
		if err := p.parseAlterColumnNullable(col); err != nil {
			return nil, err
		}
		if p.cur.Kind == kwAUTO_INCREMENT {
			p.advance()
			col.AutoInc = true
			if p.cur.Kind == int('(') {
				p.advance()
				if _, err := p.alterExpectInt(); err != nil {
					return nil, err
				}
				if _, err := p.expect(int(')')); err != nil {
					return nil, err
				}
			}
		}
		if p.cur.Kind == kwDEFAULT {
			p.advance()
			if col.Default, err = p.parseAlterDefaultValue(); err != nil {
				return nil, err
			}
		}
	}
	if p.cur.Kind == kwCOMMENT {
		p.advance()
		tok, err := p.alterExpectString()
		if err != nil {
			return nil, err
		}
		col.Comment = tok.Str
	}
	col.Loc.End = p.prev.Loc.End
	return col, nil
}

// parseAlterColumnNullable reads opt_nullable_type: NULL | NOT NULL.
func (p *Parser) parseAlterColumnNullable(col *ast.ColumnDef) error {
	switch p.cur.Kind {
	case kwNULL:
		p.advance()
		b := true
		col.Nullable = &b
	case kwNOT:
		p.advance()
		if _, err := p.expect(kwNULL); err != nil {
			return err
		}
		b := false
		col.Nullable = &b
	}
	return nil
}

// parseAlterDefaultValue reads the value after DEFAULT that both
// DorisLegacyParser.g4 opt_default_value and DorisParser.g4 columnDef take:
// STRING, NULL, an unsigned integer or decimal, BITMAP_EMPTY, or
// CURRENT_TIMESTAMP ['(' INTEGER ')'].
func (p *Parser) parseAlterDefaultValue() (ast.Node, error) {
	switch p.cur.Kind {
	case tokString, kwNULL, tokInt, tokFloat, kwBITMAP_EMPTY:
		return p.parseDefaultValue()
	case kwCURRENT_TIMESTAMP:
		if p.peekNext().Kind == int('(') {
			cp := p.save()
			p.advance() // CURRENT_TIMESTAMP
			p.advance() // '('
			if p.cur.Kind != tokInt {
				return nil, p.syntaxErrorAtCur()
			}
			p.restore(cp)
		}
		return p.parseDefaultValue()
	}
	return nil, p.syntaxErrorAtCur()
}

// parseAlterAddIndex reads ADD index_definition (DorisLegacyParser.g4:319;
// DorisParser.g4 indexDef):
//
//	INDEX [IF NOT EXISTS] ident '(' ident_list ')' [USING (BITMAP | NGRAM_BF | INVERTED)]
//	    [PROPERTIES (...)] [COMMENT STRING]
//
// cur is INDEX.
func (p *Parser) parseAlterAddIndex(action *ast.AlterTableAction) error {
	idx := &ast.IndexDef{Loc: p.cur.Loc}
	p.advance() // consume INDEX
	var err error
	if idx.IfNotExists, err = p.alterOptIfNotExists(); err != nil {
		return err
	}
	if idx.Name, _, err = p.parseIdentifier(); err != nil {
		return err
	}
	if idx.Columns, err = p.alterIdentList(); err != nil {
		return err
	}
	if p.cur.Kind == kwUSING {
		p.advance()
		if !p.alterIsWord("BITMAP") && !p.alterIsWord("NGRAM_BF") && !p.alterIsWord("INVERTED") {
			return p.syntaxErrorAtCur()
		}
		idx.IndexType = strings.ToUpper(p.advance().Str)
	}
	if idx.Properties, err = p.alterOptProperties(); err != nil {
		return err
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
	action.IndexName = idx.Name
	action.IfNotExists = idx.IfNotExists
	return nil
}

// parseAlterAddPartition reads, after ADD [TEMPORARY] PARTITION
// (DorisLegacyParser.g4:310 with single_partition_desc; DorisParser.g4:636
// addPartitionClause):
//
//	[IF NOT EXISTS] ident partition_values ['(' key_value_map ')'] opt_distribution opt_properties
//
// with partition_values one of
//
//	VALUES LESS THAN (MAXVALUE | '(' partition_key_list ')')
//	VALUES '[' '(' partition_key_list ')' ',' '(' partition_key_list ')' ')'
//	[VALUES IN] '(' list_partition_values ')'
//	(nothing)
func (p *Parser) parseAlterAddPartition(action *ast.AlterTableAction) error {
	action.Type = ast.AlterAddPartition
	item := &ast.PartitionItem{Loc: p.cur.Loc}
	var err error
	if action.IfNotExists, err = p.alterOptIfNotExists(); err != nil {
		return err
	}
	if item.Name, _, err = p.parseIdentifier(); err != nil {
		return err
	}
	switch {
	case p.cur.Kind == kwVALUES && p.peekNext().Kind == kwLESS:
		p.advance() // VALUES
		p.advance() // LESS
		if _, err := p.expect(kwTHAN); err != nil {
			return err
		}
		if p.cur.Kind == kwMAXVALUE {
			p.advance()
			item.IsMaxValue = true
		} else if item.Values, err = p.parseAlterPartitionKeyList(); err != nil {
			return err
		}
	case p.cur.Kind == kwVALUES && p.peekNext().Kind == int('['):
		p.advance() // VALUES
		p.advance() // '['
		lower, err := p.parseAlterPartitionKeyList()
		if err != nil {
			return err
		}
		if _, err := p.expect(int(',')); err != nil {
			return err
		}
		upper, err := p.parseAlterPartitionKeyList()
		if err != nil {
			return err
		}
		if _, err := p.expect(int(')')); err != nil {
			return err
		}
		item.Values = append(lower, upper...)
	case p.cur.Kind == kwVALUES:
		p.advance()
		if _, err := p.expect(kwIN); err != nil {
			return err
		}
		if item.InValues, err = p.parseAlterListPartitionValues(); err != nil {
			return err
		}
	case p.cur.Kind == int('(') && !p.alterAtStringProperty():
		if item.InValues, err = p.parseAlterListPartitionValues(); err != nil {
			return err
		}
	}
	if p.cur.Kind == int('(') {
		p.advance()
		if item.Properties, err = p.parsePropertyItems(p.parseStringProperty); err != nil {
			return err
		}
		action.PartitionProps = item.Properties
	}
	item.Loc.End = p.prev.Loc.End
	action.Partition = item
	if p.cur.Kind == kwDISTRIBUTED {
		if action.PartitionDist, err = p.parseDistributedBy(); err != nil {
			return err
		}
	}
	action.Properties, err = p.alterOptProperties()
	return err
}

// alterAtStringProperty reports whether cur is the '(' of a property list:
// '(' STRING '=' ..., not a list of partition values.
func (p *Parser) alterAtStringProperty() bool {
	if p.cur.Kind != int('(') {
		return false
	}
	cp := p.save()
	defer p.restore(cp)
	p.advance()
	if p.cur.Kind != tokString {
		return false
	}
	p.advance()
	return p.cur.Kind == int('=')
}

// parseAlterPartitionKeyList reads '(' partition_key_list ')' as both
// DorisLegacyParser.g4 partition_key_list and DorisParser.g4
// partitionValueList take it: at least one value, comma-separated, each a
// string, an unsigned integer or MAXVALUE, with NULL allowed first.
func (p *Parser) parseAlterPartitionKeyList() ([]string, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	var vals []string
	for {
		switch {
		case p.cur.Kind == tokString, p.cur.Kind == tokInt:
			vals = append(vals, p.advance().Str)
		case p.cur.Kind == kwMAXVALUE, p.cur.Kind == kwNULL && len(vals) == 0:
			vals = append(vals, strings.ToUpper(p.advance().Str))
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

// parseAlterListPartitionValues reads the values of a list partition as both
// grammars take them: DorisLegacyParser.g4 '(' list_partition_values_list ')'
// and DorisParser.g4 inPartitionDef, a list of parenthesized value lists or
// one value list. Either every item is a string, an unsigned integer or NULL,
// or every item is '(' value (',' value)* ')' whose values may also be
// MAXVALUE.
func (p *Parser) parseAlterListPartitionValues() ([][]string, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	var out [][]string
	nested := p.cur.Kind == int('(')
	for {
		if nested {
			if _, err := p.expect(int('(')); err != nil {
				return nil, err
			}
			var vals []string
			for {
				switch p.cur.Kind {
				case tokString, tokInt:
					vals = append(vals, p.advance().Str)
				case kwMAXVALUE, kwNULL:
					vals = append(vals, strings.ToUpper(p.advance().Str))
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
			out = append(out, vals)
		} else {
			switch p.cur.Kind {
			case tokString, tokInt:
				out = append(out, []string{p.advance().Str})
			case kwNULL:
				out = append(out, []string{strings.ToUpper(p.advance().Str)})
			default:
				return nil, p.syntaxErrorAtCur()
			}
		}
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

// parseAlterAddPartitions reads (DorisLegacyParser.g4:327-328;
// DorisParser.g4:663 alterMultiPartitionClause):
//
//	PARTITIONS FROM '(' partition_key_list ')' TO '(' partition_key_list ')'
//	    INTERVAL INTEGER [ident] opt_properties
//
// cur is PARTITIONS.
func (p *Parser) parseAlterAddPartitions(action *ast.AlterTableAction) error {
	action.Type = ast.AlterAddPartitions
	item := &ast.PartitionItem{IsStep: true, Loc: p.cur.Loc}
	p.advance() // consume PARTITIONS
	var err error
	if _, err = p.expect(kwFROM); err != nil {
		return err
	}
	if item.FromValues, err = p.parseAlterPartitionKeyList(); err != nil {
		return err
	}
	if _, err = p.expect(kwTO); err != nil {
		return err
	}
	if item.ToValues, err = p.parseAlterPartitionKeyList(); err != nil {
		return err
	}
	if _, err = p.expect(kwINTERVAL); err != nil {
		return err
	}
	tok, err := p.alterExpectInt()
	if err != nil {
		return err
	}
	item.Interval = tok.Str
	if p.cur.Kind != kwPROPERTIES && isIdentifierToken(p.cur.Kind) {
		item.IntervalUnit = strings.ToUpper(p.advance().Str)
	}
	item.Loc.End = p.prev.Loc.End
	action.Partition = item
	action.Properties, err = p.alterOptProperties()
	return err
}

// parseAlterDropClause reads the DROP alternatives of alter_table_clause.
// cur is DROP.
func (p *Parser) parseAlterDropClause(action *ast.AlterTableAction) error {
	p.advance() // consume DROP
	var err error
	switch p.cur.Kind {
	case kwCOLUMN:
		// DROP COLUMN ident opt_from_rollup opt_properties
		// (DorisLegacyParser.g4:307; DorisParser.g4:631 dropColumnClause)
		p.advance()
		action.Type = ast.AlterDropColumn
		if action.ColumnName, _, err = p.parseIdentifier(); err != nil {
			return err
		}
		return p.parseAlterRollupAndProps(action, false)
	case kwTEMPORARY, kwPARTITION:
		// DROP opt_tmp PARTITION opt_if_exists ident opt_force ['FROM' 'INDEX' ident]
		// (DorisLegacyParser.g4:311-312; DorisParser.g4:641 dropPartitionClause)
		if p.cur.Kind == kwTEMPORARY {
			p.advance()
			action.Temporary = true
		}
		if _, err = p.expect(kwPARTITION); err != nil {
			return err
		}
		action.Type = ast.AlterDropPartition
		if action.IfExists, err = p.alterOptIfExists(); err != nil {
			return err
		}
		if action.PartitionName, _, err = p.parseIdentifier(); err != nil {
			return err
		}
		if p.cur.Kind == kwFORCE {
			p.advance()
			action.Force = true
		}
		if p.cur.Kind == kwFROM {
			p.advance()
			if _, err = p.expect(kwINDEX); err != nil {
				return err
			}
			action.IndexName, _, err = p.parseIdentifier()
		}
		return err
	case kwINDEX, kwBRANCH, kwTAG:
		// DROP INDEX opt_if_exists ident (DorisLegacyParser.g4:320;
		// DorisParser.g4:655 dropIndexClause), DROP BRANCH / TAG
		// opt_if_exists ident (DorisLegacyParser.g4:335-336;
		// DorisParser.g4:706 dropBranchClause, :710 dropTagClause)
		switch p.cur.Kind {
		case kwINDEX:
			action.Type = ast.AlterDropIndex
		case kwBRANCH:
			action.Type = ast.AlterDropBranch
		default:
			action.Type = ast.AlterDropTag
		}
		p.advance()
		if action.IfExists, err = p.alterOptIfExists(); err != nil {
			return err
		}
		name, _, err := p.parseIdentifier()
		if err != nil {
			return err
		}
		if action.Type == ast.AlterDropIndex {
			action.IndexName = name
		} else {
			action.RefName = name
		}
		return nil
	}
	return p.syntaxErrorAtCur()
}

// parseAlterModifyClause reads the MODIFY alternatives of alter_table_clause.
// cur is MODIFY.
func (p *Parser) parseAlterModifyClause(action *ast.AlterTableAction) error {
	p.advance() // consume MODIFY
	var err error
	switch p.cur.Kind {
	case kwCOLUMN:
		p.advance()
		action.Type = ast.AlterModifyColumn
		if isIdentifierToken(p.cur.Kind) && p.peekNext().Kind == kwCOMMENT {
			// MODIFY COLUMN ident COMMENT STRING_LITERAL
			// (DorisLegacyParser.g4:326; DorisParser.g4:660)
			col := &ast.ColumnDef{Loc: p.cur.Loc}
			col.Name, _, _ = p.parseIdentifier()
			p.advance() // consume COMMENT
			tok, err := p.alterExpectString()
			if err != nil {
				return err
			}
			col.Comment = tok.Str
			col.Loc.End = tok.Loc.End
			action.Column = col
			return nil
		}
		// MODIFY COLUMN column_definition opt_col_pos opt_from_rollup opt_properties
		// (DorisLegacyParser.g4:308; DorisParser.g4:632 modifyColumnClause)
		if action.Column, err = p.parseAlterColumnDef(); err != nil {
			return err
		}
		if err := p.parseAlterColumnPosition(action); err != nil {
			return err
		}
		return p.parseAlterRollupAndProps(action, false)
	case kwTEMPORARY, kwPARTITION:
		// MODIFY opt_tmp PARTITION (ident | '(' ident_list ')' | '(' '*' ')')
		//     SET '(' key_value_map ')'
		// (DorisLegacyParser.g4:313-315; DorisParser.g4:643 modifyPartitionClause)
		if p.cur.Kind == kwTEMPORARY {
			p.advance()
			action.Temporary = true
		}
		if _, err = p.expect(kwPARTITION); err != nil {
			return err
		}
		action.Type = ast.AlterModifyPartition
		switch {
		case p.cur.Kind == int('(') && p.peekNext().Kind == int('*'):
			p.advance()
			p.advance()
			if _, err = p.expect(int(')')); err != nil {
				return err
			}
			action.PartitionStar = true
		case p.cur.Kind == int('('):
			if action.PartitionList, err = p.alterIdentList(); err != nil {
				return err
			}
		default:
			if action.PartitionName, _, err = p.parseIdentifier(); err != nil {
				return err
			}
		}
		if _, err = p.expect(kwSET); err != nil {
			return err
		}
		action.Properties, err = p.parseParenProperties()
		return err
	case kwDISTRIBUTION:
		// MODIFY DISTRIBUTION opt_distribution (DorisLegacyParser.g4:324;
		// DorisParser.g4:657 modifyDistributionClause)
		p.advance()
		action.Type = ast.AlterModifyDistribution
		if p.cur.Kind == kwDISTRIBUTED {
			action.Distribution, err = p.parseDistributedBy()
		}
		return err
	case kwCOMMENT:
		// MODIFY COMMENT STRING_LITERAL (DorisLegacyParser.g4:325;
		// DorisParser.g4:659)
		p.advance()
		tok, err := p.alterExpectString()
		if err != nil {
			return err
		}
		action.Type = ast.AlterModifyComment
		action.Comment = tok.Str
		return nil
	case kwENGINE:
		// MODIFY ENGINE TO ident opt_properties (DorisLegacyParser.g4:327;
		// DorisParser.g4:661)
		p.advance()
		if _, err = p.expect(kwTO); err != nil {
			return err
		}
		action.Type = ast.AlterModifyEngine
		if action.Engine, _, err = p.parseIdentifier(); err != nil {
			return err
		}
		action.Properties, err = p.alterOptProperties()
		return err
	}
	return p.syntaxErrorAtCur()
}

// parseAlterReplaceClause reads the REPLACE alternatives of
// alter_table_clause. cur is REPLACE.
//
//	REPLACE WITH TABLE ident opt_properties                 (DorisLegacyParser.g4:317; DorisParser.g4:648)
//	REPLACE partition_names WITH [partition_names] [FORCE] opt_properties
//	                                                        (DorisLegacyParser.g4:316; DorisParser.g4:646)
//	REPLACE BRANCH ident branch_options                     (DorisLegacyParser.g4:330; DorisParser.g4:675)
//	REPLACE TAG ident tag_options                           (DorisLegacyParser.g4:333; DorisParser.g4:670)
//
// CUP shifts WITH after REPLACE into REPLACE WITH TABLE, so the first
// partition list is not optional there although both grammars write it so.
func (p *Parser) parseAlterReplaceClause(action *ast.AlterTableAction) error {
	p.advance() // consume REPLACE
	switch p.cur.Kind {
	case kwWITH:
		p.advance()
		if _, err := p.expect(kwTABLE); err != nil {
			return err
		}
		action.Type = ast.AlterReplaceTable
		var err error
		if action.NewName, _, err = p.parseIdentifier(); err != nil {
			return err
		}
		action.Properties, err = p.alterOptProperties()
		return err
	case kwBRANCH, kwTAG:
		action.OrReplace = true
		return p.parseAlterRefBody(action)
	}
	if !p.alterAtPartitionSpec() {
		return p.syntaxErrorAtCur()
	}
	action.Type = ast.AlterReplacePartition
	temp, names, star, err := p.parseAlterPartitionSpec()
	if err != nil {
		return err
	}
	action.Temporary, action.PartitionList, action.PartitionStar = temp, names, star
	if _, err := p.expect(kwWITH); err != nil {
		return err
	}
	if p.alterAtPartitionSpec() {
		if _, action.TempPartitionList, _, err = p.parseAlterPartitionSpec(); err != nil {
			return err
		}
	}
	if p.cur.Kind == kwFORCE {
		p.advance()
		action.Force = true
	}
	action.Properties, err = p.alterOptProperties()
	return err
}

// parseAlterRenameClause reads the RENAME alternatives of alter_table_clause
// (DorisLegacyParser.g4:318-322; DorisParser.g4:650-653):
//
//	RENAME ident | RENAME (ROLLUP | PARTITION | COLUMN) ident ident
//
// cur is RENAME.
func (p *Parser) parseAlterRenameClause(action *ast.AlterTableAction) error {
	p.advance() // consume RENAME
	var typ ast.AlterActionType
	switch p.cur.Kind {
	case kwROLLUP:
		typ = ast.AlterRenameRollup
	case kwPARTITION:
		typ = ast.AlterRenamePartition
	case kwCOLUMN:
		typ = ast.AlterRenameColumn
	}
	if typ != 0 && isIdentifierToken(p.peekNext().Kind) {
		p.advance()
		oldName, _, err := p.parseIdentifier()
		if err != nil {
			return err
		}
		newName, _, err := p.parseIdentifier()
		if err != nil {
			return err
		}
		action.Type = typ
		action.ColumnName = oldName
		action.NewName = newName
		if typ == ast.AlterRenameRollup {
			action.RollupName = oldName
		}
		return nil
	}
	name, loc, err := p.parseIdentifier()
	if err != nil {
		return err
	}
	action.Type = ast.AlterRenameTable
	action.NewTableName = &ast.ObjectName{Parts: []string{name}, Loc: loc}
	return nil
}

// parseAlterCreateRefClause reads the CREATE alternatives of
// alter_table_clause (DorisLegacyParser.g4:329-334; DorisParser.g4:670
// createOrReplaceTagClause, :675 createOrReplaceBranchClause):
//
//	CREATE (BRANCH | TAG) [IF NOT EXISTS] ident options
//	CREATE OR REPLACE (BRANCH | TAG) ident options
//
// cur is CREATE.
func (p *Parser) parseAlterCreateRefClause(action *ast.AlterTableAction) error {
	p.advance() // consume CREATE
	if p.cur.Kind == kwOR {
		p.advance()
		if _, err := p.expect(kwREPLACE); err != nil {
			return err
		}
		if p.cur.Kind != kwBRANCH && p.cur.Kind != kwTAG {
			return p.syntaxErrorAtCur()
		}
		action.OrReplace = true
		return p.parseAlterRefBody(action)
	}
	if p.cur.Kind != kwBRANCH && p.cur.Kind != kwTAG {
		return p.syntaxErrorAtCur()
	}
	return p.parseAlterRefBody(action)
}

// parseAlterRefBody reads (BRANCH | TAG) [IF NOT EXISTS] ident and the
// options: tag_options (DorisLegacyParser.g4:1888; DorisParser.g4:680) or
// branch_options (DorisLegacyParser.g4:1828; DorisParser.g4:684).
//
//	[AS OF VERSION INTEGER] [RETAIN INTEGER (DAYS | HOURS | MINUTES)]
//	[WITH SNAPSHOT RETENTION (INTEGER SNAPSHOTS [time] | time)]   -- branches only
//
// IF NOT EXISTS follows CREATE only. cur is BRANCH or TAG.
func (p *Parser) parseAlterRefBody(action *ast.AlterTableAction) error {
	branch := p.cur.Kind == kwBRANCH
	if branch {
		action.Type = ast.AlterCreateBranch
	} else {
		action.Type = ast.AlterCreateTag
	}
	p.advance() // consume BRANCH / TAG
	var err error
	if !action.OrReplace {
		if action.IfNotExists, err = p.alterOptIfNotExists(); err != nil {
			return err
		}
	}
	if action.RefName, _, err = p.parseIdentifier(); err != nil {
		return err
	}
	if p.cur.Kind == kwAS {
		p.advance()
		for _, kw := range []int{kwOF, kwVERSION} {
			if _, err := p.expect(kw); err != nil {
				return err
			}
		}
		tok, err := p.alterExpectInt()
		if err != nil {
			return err
		}
		action.AsOfVersion = tok.Str
	}
	if p.cur.Kind == kwRETAIN {
		p.advance()
		if action.RetainTime, err = p.parseAlterTimeValue(); err != nil {
			return err
		}
	}
	if branch && p.cur.Kind == kwWITH {
		p.advance()
		for _, kw := range []int{kwSNAPSHOT, kwRETENTION} {
			if _, err := p.expect(kw); err != nil {
				return err
			}
		}
		if p.cur.Kind == tokInt && p.peekNext().Kind == kwSNAPSHOTS {
			action.MinSnapshots = p.advance().Str
			p.advance() // consume SNAPSHOTS
			if p.cur.Kind == tokInt {
				if action.MaxSnapshotAge, err = p.parseAlterTimeValue(); err != nil {
					return err
				}
			}
		} else if action.MaxSnapshotAge, err = p.parseAlterTimeValue(); err != nil {
			return err
		}
	}
	return nil
}

// parseAlterTimeValue reads time_value_with_unit: INTEGER (DAYS | HOURS |
// MINUTES) (DorisLegacyParser.g4:1855; DorisParser.g4:702).
func (p *Parser) parseAlterTimeValue() (string, error) {
	tok, err := p.alterExpectInt()
	if err != nil {
		return "", err
	}
	switch p.cur.Kind {
	case kwDAYS, kwHOURS, kwMINUTES:
		return tok.Str + " " + strings.ToUpper(p.advance().Str), nil
	}
	return "", p.syntaxErrorAtCur()
}
