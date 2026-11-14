package parser

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/redshift/ast"
)

// parseCreateStmt parses a CREATE TABLE statement.
//
// Ref: https://www.postgresql.org/docs/17/sql-createtable.html
//
//	CreateStmt:
//	    CREATE OptTemp TABLE qualified_name '(' OptTableElementList ')' OptInherit OptPartitionSpec OptAccessMethod OptWith OnCommitOption OptTableSpace
//	    | CREATE OptTemp TABLE IF NOT EXISTS qualified_name '(' OptTableElementList ')' OptInherit OptPartitionSpec OptAccessMethod OptWith OnCommitOption OptTableSpace
//	    | CREATE OptTemp TABLE qualified_name PARTITION OF qualified_name ...
//	    | CREATE OptTemp TABLE qualified_name OF any_name ...
func (p *Parser) parseCreateStmt() (*nodes.CreateStmt, error) {
	loc := p.pos()
	p.advance() // consume CREATE

	// OptTemp
	relpersistence, err := p.parseOptTemp()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(TABLE); err != nil {
		return nil, err
	}

	// IF NOT EXISTS
	ifNotExists := false
	if p.cur.Type == IF_P {
		p.advance()
		if _, err := p.expect(NOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		ifNotExists = true
	}

	// qualified_name
	nameStart := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rv := makeRangeVarFromNames(names, nodes.Loc{Start: nameStart, End: p.prev.End})
	rv.Relpersistence = relpersistence

	stmt := &nodes.CreateStmt{
		Relation:    rv,
		IfNotExists: ifNotExists,
	}

	// Check for PARTITION OF, OF
	switch p.cur.Type {
	case PARTITION:
		p.advance() // PARTITION
		if _, err := p.expect(OF); err != nil {
			return nil, err
		}
		return p.parseCreateStmtPartitionOf(stmt, relpersistence)
	case OF:
		p.advance()
		return p.parseCreateStmtOf(stmt)
	}

	// '(' OptTableElementList ')'
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	tableElts, err := p.parseOptTableElementList()
	if err != nil {
		return nil, err
	}
	stmt.TableElts = tableElts
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	if _, err := p.finishCreateStmtTail(stmt); err != nil {
		return nil, err
	}
	stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End}
	return stmt, nil
}

// finishCreateStmtTail parses what follows CreateStmt's
// '(' OptTableElementList ')':
//
//	OptInherit OptPartitionSpec table_access_method_clause OptWith
//	OnCommitOption OptTableSpace
//
// and Redshift's table attributes after them.
func (p *Parser) finishCreateStmtTail(stmt *nodes.CreateStmt) (*nodes.CreateStmt, error) {
	var err error
	if stmt.InhRelations, err = p.parseOptInherit(); err != nil {
		return nil, err
	}
	if stmt.Partspec, err = p.parseOptPartitionSpec(); err != nil {
		return nil, err
	}
	if stmt.AccessMethod, err = p.parseOptAccessMethod(); err != nil {
		return nil, err
	}
	if stmt.Options, err = p.parseOptWith(); err != nil {
		return nil, err
	}
	if stmt.OnCommit, err = p.parseOnCommitOption(); err != nil {
		return nil, err
	}
	if stmt.Tablespacename, err = p.parseOptTableSpace(); err != nil {
		return nil, err
	}
	if stmt.Options, err = p.parseRedshiftCreateTableAttributesInto(stmt.Options); err != nil {
		return nil, err
	}
	return stmt, nil
}

// parseOptTemp parses OptTemp production.
//
//	OptTemp: TEMPORARY | TEMP | LOCAL TEMPORARY | LOCAL TEMP | UNLOGGED | /* EMPTY */
func (p *Parser) parseOptTemp() (byte, error) {
	switch p.cur.Type {
	case TEMPORARY, TEMP:
		p.advance()
		return 't', nil
	case LOCAL, GLOBAL:
		// GLOBAL is accepted and ignored, as gram.y does with a warning.
		p.advance()
		if p.cur.Type != TEMPORARY && p.cur.Type != TEMP {
			return 0, p.syntaxErrorAtCur()
		}
		p.advance()
		return 't', nil
	case UNLOGGED:
		p.advance()
		return 'u', nil
	default:
		return 'p', nil
	}
}

// parseCreateStmtPartitionOf parses the PARTITION OF variant of CREATE TABLE.
func (p *Parser) parseCreateStmtPartitionOf(stmt *nodes.CreateStmt, relpersistence byte) (*nodes.CreateStmt, error) {
	// parent table name
	parentStart := p.pos()
	parentNames, err := p.parseQualifiedName()
	if err != nil {
		return stmt, err
	}
	parentRv := makeRangeVarFromNames(parentNames, nodes.Loc{Start: parentStart, End: p.prev.End})
	parentRv.Relpersistence = relpersistence

	// Optional '(' TypedTableElementList ')'
	if p.cur.Type == '(' {
		p.advance()
		elts, err := p.parseTypedTableElementList()
		if err != nil {
			return nil, err
		}
		stmt.TableElts = elts
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	// ForValues
	stmt.InhRelations = &nodes.List{Items: []nodes.Node{parentRv}}
	partbound, err := p.parseForValues()
	if err != nil {
		return nil, err
	}
	stmt.Partbound = partbound

	// OptPartitionSpec
	partspec, err := p.parseOptPartitionSpec()
	if err != nil {
		return nil, err
	}
	stmt.Partspec = partspec

	// OptAccessMethod
	if stmt.AccessMethod, err = p.parseOptAccessMethod(); err != nil {
		return nil, err
	}

	// OptWith
	if stmt.Options, err = p.parseOptWith(); err != nil {
		return nil, err
	}

	// OnCommitOption
	if stmt.OnCommit, err = p.parseOnCommitOption(); err != nil {
		return nil, err
	}

	// OptTableSpace
	if stmt.Tablespacename, err = p.parseOptTableSpace(); err != nil {
		return nil, err
	}
	stmt.Options, err = p.parseRedshiftCreateTableAttributesInto(stmt.Options)
	if err != nil {
		return nil, err
	}

	return stmt, nil
}

// parseCreateStmtOf parses the OF typename variant of CREATE TABLE.
func (p *Parser) parseCreateStmtOf(stmt *nodes.CreateStmt) (*nodes.CreateStmt, error) {
	// any_name -> TypeName
	anyName, err := p.parseAnyName()
	if err != nil {
		return stmt, err
	}
	stmt.OfTypename = makeTypeNameFromNameList(anyName)

	// OptTypedTableElementList: '(' TypedTableElementList ')' | EMPTY
	if p.cur.Type == '(' {
		p.advance()
		elts, err := p.parseTypedTableElementList()
		if err != nil {
			return nil, err
		}
		stmt.TableElts = elts
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	// OptAccessMethod
	if stmt.AccessMethod, err = p.parseOptAccessMethod(); err != nil {
		return nil, err
	}

	// OptWith
	if stmt.Options, err = p.parseOptWith(); err != nil {
		return nil, err
	}

	// OnCommitOption
	if stmt.OnCommit, err = p.parseOnCommitOption(); err != nil {
		return nil, err
	}

	// OptTableSpace
	if stmt.Tablespacename, err = p.parseOptTableSpace(); err != nil {
		return nil, err
	}
	stmt.Options, err = p.parseRedshiftCreateTableAttributesInto(stmt.Options)
	if err != nil {
		return nil, err
	}

	return stmt, nil
}

// makeTypeNameFromNameList creates a TypeName from an any_name list.
func makeTypeNameFromNameList(names *nodes.List) *nodes.TypeName {
	return &nodes.TypeName{
		Names: names,
		Loc:   nodes.NoLoc(),
	}
}

// parseOptTableElementList parses OptTableElementList.
//
//	OptTableElementList: TableElementList | /* EMPTY */
func (p *Parser) parseOptTableElementList() (*nodes.List, error) {
	if p.cur.Type == ')' {
		return nil, nil
	}
	return p.parseTableElementList()
}

// parseTableElementList parses a comma-separated list of TableElements.
//
//	TableElementList: TableElement | TableElementList ',' TableElement
func (p *Parser) parseTableElementList() (*nodes.List, error) {
	elem, err := p.parseTableElement()
	if err != nil {
		return nil, err
	}
	if elem == nil {
		return nil, nil
	}
	items := []nodes.Node{elem}
	for p.cur.Type == ',' {
		p.advance()
		elem, err = p.parseTableElement()
		if err != nil {
			return nil, err
		}
		if elem == nil {
			return nil, p.syntaxErrorAtCur()
		}
		items = append(items, elem)
	}
	return &nodes.List{Items: items}, nil
}

// parseTableElement parses a single table element.
//
//	TableElement: columnDef | TableConstraint | TableLikeClause
func (p *Parser) parseTableElement() (nodes.Node, error) {
	if p.startsTableConstraint() {
		return p.parseTableConstraint()
	}
	if p.cur.Type == LIKE {
		return p.parseTableLikeClause()
	}
	// columnDef: ColId Typename opt_column_constraints
	return p.parseColumnDef()
}

// parseColumnDef parses a column definition.
//
// Ref: gram.y columnDef
//
//	columnDef: ColId Typename opt_column_constraints
func (p *Parser) parseColumnDef() (*nodes.ColumnDef, error) {
	start := p.pos()
	colname, err := p.parseColId()
	if err != nil {
		return nil, err
	}

	tn, err := p.parseTypename()
	if err != nil {
		return nil, err
	}

	n := &nodes.ColumnDef{
		Colname:  colname,
		TypeName: tn,
		IsLocal:  true,
		Loc:      nodes.NoLoc(),
	}

	n.Fdwoptions, err = p.parseCreateGenericOptions()
	if err != nil {
		return nil, err
	}

	// opt_column_constraints
	qualList, err := p.parseOptColumnConstraints()
	if err != nil {
		return nil, err
	}
	splitColQualList(qualList, n)
	n.Loc = nodes.Loc{Start: start, End: p.prev.End}

	return n, nil
}

// parseOptColumnConstraints parses opt_column_constraints.
//
//	opt_column_constraints: ColConstraintList | /* EMPTY */
func (p *Parser) parseOptColumnConstraints() (*nodes.List, error) {
	var items []nodes.Node
	for {
		c, err := p.parseColConstraint()
		if err != nil {
			return nil, err
		}
		if c == nil {
			break
		}
		items = append(items, c)
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: items}, nil
}

// parseColConstraint parses a single ColConstraint.
//
//	ColConstraint:
//	    CONSTRAINT name ColConstraintElem
//	    | ColConstraintElem
//	    | COLLATE any_name
//	    | COMPRESSION ColId
//	    | STORAGE ColId
//	    | ConstraintAttr
//
// A constraint spans from its first token, CONSTRAINT when it is named, to
// its last (gram.y sets location = @1).
func (p *Parser) parseColConstraint() (nodes.Node, error) {
	start := p.pos()
	n, err := p.parseColConstraintBody()
	if c, ok := n.(*nodes.Constraint); ok && err == nil {
		c.Loc = nodes.Loc{Start: start, End: p.prev.End}
	}
	return n, err
}

// parseColConstraintBody parses the ColConstraint parseColConstraint
// locates.
func (p *Parser) parseColConstraintBody() (nodes.Node, error) {
	switch p.cur.Type {
	case CONSTRAINT:
		p.advance()
		name, err := p.parseName()
		if err != nil {
			return nil, err
		}
		c, err := p.parseColConstraintElem()
		if err != nil {
			return nil, err
		}
		if c == nil {
			return nil, nil
		}
		c.(*nodes.Constraint).Conname = name
		return c, nil
	case NOT:
		// NOT NULL or NOT DEFERRABLE
		next := p.peekNext()
		if next.Type == NULL_P {
			p.advance() // NOT
			p.advance() // NULL
			return &nodes.Constraint{
				Contype: nodes.CONSTR_NOTNULL,
				Loc:     nodes.NoLoc(),
			}, nil
		}
		if next.Type == DEFERRABLE {
			p.advance() // NOT
			p.advance() // DEFERRABLE
			return &nodes.Constraint{
				Contype: nodes.CONSTR_ATTR_NOT_DEFERRABLE,
				Loc:     nodes.NoLoc(),
			}, nil
		}
		// Nothing else follows NOT in a column's constraint list, which
		// takes no NOT VALID.
		if !p.collectMode() {
			return nil, p.syntaxErrorAtTok(next)
		}
		return nil, nil
	case NULL_P:
		p.advance()
		return &nodes.Constraint{
			Contype: nodes.CONSTR_NULL,
			Loc:     nodes.NoLoc(),
		}, nil
	case UNIQUE:
		return p.parseColConstraintElem()
	case PRIMARY:
		return p.parseColConstraintElem()
	case CHECK:
		return p.parseColConstraintElem()
	case DEFAULT:
		return p.parseColConstraintElem()
	case REFERENCES:
		return p.parseColConstraintElem()
	case GENERATED:
		return p.parseColConstraintElem()
	case IDENTITY_P:
		return p.parseRedshiftIdentityColumnConstraint()
	case COLLATE:
		p.advance()
		collname, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && collname == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return &nodes.CollateClause{
			Collname: collname,
			Loc:      nodes.NoLoc(),
		}, nil
	case COMPRESSION:
		// column_compression: COMPRESSION ColId | COMPRESSION DEFAULT
		p.advance()
		var id string
		if p.cur.Type == DEFAULT {
			p.advance()
			id = "default"
		} else {
			var err error
			id, err = p.parseColId()
			if err != nil {
				return nil, err
			}
		}
		if !p.collectMode() && id == "" {
			return nil, p.syntaxErrorAtCur()
		}
		return &nodes.DefElem{
			Defname: "compression",
			Arg:     &nodes.String{Str: id},
			Loc:     nodes.NoLoc(),
		}, nil
	case STORAGE:
		// column_storage: STORAGE ColId | STORAGE DEFAULT
		p.advance()
		var id string
		if p.cur.Type == DEFAULT {
			p.advance()
			id = "default"
		} else {
			var err error
			id, err = p.parseColId()
			if err != nil {
				return nil, err
			}
		}
		if !p.collectMode() && id == "" {
			return nil, p.syntaxErrorAtCur()
		}
		return &nodes.DefElem{
			Defname: "storage",
			Arg:     &nodes.String{Str: id},
			Loc:     nodes.NoLoc(),
		}, nil
	case DEFERRABLE:
		p.advance()
		return &nodes.Constraint{
			Contype: nodes.CONSTR_ATTR_DEFERRABLE,
			Loc:     nodes.NoLoc(),
		}, nil
	case INITIALLY:
		p.advance()
		if p.cur.Type == DEFERRED {
			p.advance()
			return &nodes.Constraint{
				Contype: nodes.CONSTR_ATTR_DEFERRED,
				Loc:     nodes.NoLoc(),
			}, nil
		}
		if p.cur.Type == IMMEDIATE {
			p.advance()
			return &nodes.Constraint{
				Contype: nodes.CONSTR_ATTR_IMMEDIATE,
				Loc:     nodes.NoLoc(),
			}, nil
		}
		if p.collectMode() {
			p.addTokenCandidate(DEFERRED)
			p.addTokenCandidate(IMMEDIATE)
			return nil, errCollecting
		}
		return nil, p.syntaxErrorAtCur()
	// optional-probe: ColConstraint loop terminator — parseOptColumnConstraints breaks on nil
	default:
		if attr, err := p.parseRedshiftColumnAttribute(); attr != nil || err != nil {
			return attr, err
		}
		return nil, nil
	}
}

// parseColConstraintElem parses a ColConstraintElem.
//
//	ColConstraintElem:
//	    NOT NULL | NULL | UNIQUE | PRIMARY KEY | CHECK '(' a_expr ')' no_inherit
//	    | DEFAULT b_expr | REFERENCES qualified_name opt_column_list key_match key_actions ConstraintAttributeSpec
//	    | GENERATED ALWAYS AS IDENTITY_P OptParenthesizedSeqOptList
//	    | GENERATED BY DEFAULT AS IDENTITY_P OptParenthesizedSeqOptList
//	    | GENERATED ALWAYS AS '(' a_expr ')' STORED
//	    | GENERATED BY DEFAULT AS '(' a_expr ')' STORED
func (p *Parser) parseColConstraintElem() (nodes.Node, error) {
	switch p.cur.Type {
	case NOT:
		p.advance() // NOT
		if _, err := p.expect(NULL_P); err != nil {
			return nil, err
		}
		return &nodes.Constraint{
			Contype: nodes.CONSTR_NOTNULL,
			Loc:     nodes.NoLoc(),
		}, nil
	case NULL_P:
		p.advance()
		return &nodes.Constraint{
			Contype: nodes.CONSTR_NULL,
			Loc:     nodes.NoLoc(),
		}, nil
	case UNIQUE:
		p.advance()
		nullsNotDistinct, err := p.parseOptUniqueNullTreatment()
		if err != nil {
			return nil, err
		}
		options, indexspace, err := p.parseIndexOptions()
		if err != nil {
			return nil, err
		}
		return &nodes.Constraint{
			Contype:          nodes.CONSTR_UNIQUE,
			NullsNotDistinct: nullsNotDistinct,
			Options:          options,
			Indexspace:       indexspace,
			Loc:              nodes.NoLoc(),
		}, nil
	case PRIMARY:
		p.advance() // PRIMARY
		if _, err := p.expect(KEY); err != nil {
			return nil, err
		}
		options, indexspace, err := p.parseIndexOptions()
		if err != nil {
			return nil, err
		}
		return &nodes.Constraint{
			Contype:    nodes.CONSTR_PRIMARY,
			Options:    options,
			Indexspace: indexspace,
			Loc:        nodes.NoLoc(),
		}, nil
	case CHECK:
		// CHECK '(' a_expr ')' opt_no_inherit: the parentheses belong to
		// the production, not to the expression.
		p.advance() // CHECK
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		expr, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		noInherit, err := p.parseNoInherit()
		if err != nil {
			return nil, err
		}
		return &nodes.Constraint{
			Contype:        nodes.CONSTR_CHECK,
			RawExpr:        expr,
			IsNoInherit:    noInherit,
			Loc:            nodes.NoLoc(),
			InitiallyValid: true,
		}, nil
	case DEFAULT:
		p.advance()
		expr, err := p.parseBExpr(0)
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return &nodes.Constraint{
			Contype: nodes.CONSTR_DEFAULT,
			RawExpr: expr,
			Loc:     nodes.NoLoc(),
		}, nil
	case REFERENCES:
		p.advance()
		refStart := p.pos()
		refNames, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && refNames == nil {
			return nil, p.syntaxErrorAtCur()
		}
		refRv := makeRangeVarFromNames(refNames, nodes.Loc{Start: refStart, End: p.prev.End})
		pkAttrs, err := p.parseOptColumnList()
		if err != nil {
			return nil, err
		}
		matchType, err := p.parseKeyMatch()
		if err != nil {
			return nil, err
		}
		updAction, delAction, delSetCols, err := p.parseKeyActions()
		if err != nil {
			return nil, err
		}
		// The column's DEFERRABLE and INITIALLY clauses that follow are
		// kept on the key; NOT VALID and NO INHERIT do not follow a
		// column's constraint.
		attrs, err := p.parseConstraintAttrs(false)
		if err != nil {
			return nil, err
		}
		n := &nodes.Constraint{
			Contype:        nodes.CONSTR_FOREIGN,
			Pktable:        refRv,
			PkAttrs:        pkAttrs,
			FkMatchtype:    matchType,
			FkUpdaction:    updAction,
			FkDelaction:    delAction,
			FkDelsetcols:   delSetCols,
			Loc:            nodes.NoLoc(),
			InitiallyValid: true,
		}
		applyConstraintAttrs(n, attrs)
		return n, nil
	case GENERATED:
		return p.parseGeneratedConstraint()
	// exhaustive: gram.y:3901 — ColConstraintElem enumerates NOT NULL /
	// NULL / UNIQUE / PRIMARY KEY / CHECK / DEFAULT / REFERENCES /
	// GENERATED. A `CONSTRAINT name <bad>` should raise rather than be
	// silently dropped at the caller.
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseGeneratedConstraint parses GENERATED ... column constraints.
func (p *Parser) parseGeneratedConstraint() (*nodes.Constraint, error) {
	p.advance() // GENERATED

	if p.cur.Type == ALWAYS {
		p.advance() // ALWAYS
		if _, err := p.expect(AS); err != nil {
			return nil, err
		}
		if p.cur.Type == IDENTITY_P {
			// GENERATED ALWAYS AS IDENTITY OptParenthesizedSeqOptList
			p.advance() // IDENTITY
			opts, err := p.parseOptParenthesizedSeqOptList()
			if err != nil {
				return nil, err
			}
			return &nodes.Constraint{
				Contype:       nodes.CONSTR_IDENTITY,
				GeneratedWhen: 'a',
				Options:       opts,
				Loc:           nodes.NoLoc(),
			}, nil
		}
		// GENERATED ALWAYS AS '(' a_expr ')' STORED
		expr, err := p.parseGeneratedExpr()
		if err != nil {
			return nil, err
		}
		return &nodes.Constraint{
			Contype:       nodes.CONSTR_GENERATED,
			GeneratedWhen: 'a',
			RawExpr:       expr,
			Loc:           nodes.NoLoc(),
		}, nil
	}

	// GENERATED BY DEFAULT AS ...
	for _, tok := range []int{BY, DEFAULT, AS} {
		if _, err := p.expect(tok); err != nil {
			return nil, err
		}
	}
	if p.cur.Type == IDENTITY_P {
		p.advance()
		opts, err := p.parseOptParenthesizedSeqOptList()
		if err != nil {
			return nil, err
		}
		return &nodes.Constraint{
			Contype:       nodes.CONSTR_IDENTITY,
			GeneratedWhen: 'd',
			Options:       opts,
			Loc:           nodes.NoLoc(),
		}, nil
	}
	// GENERATED BY DEFAULT AS '(' a_expr ')' STORED
	expr, err := p.parseGeneratedExpr()
	if err != nil {
		return nil, err
	}
	return &nodes.Constraint{
		Contype:       nodes.CONSTR_GENERATED,
		GeneratedWhen: 'd',
		RawExpr:       expr,
		Loc:           nodes.NoLoc(),
	}, nil
}

// parseGeneratedExpr parses the '(' a_expr ')' STORED that ends a
// generated column's constraint (gram.y ColConstraintElem).
func (p *Parser) parseGeneratedExpr() (nodes.Node, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	expr, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	if _, err := p.expect(STORED); err != nil {
		return nil, err
	}
	return expr, nil
}

// parseOptParenthesizedSeqOptList parses OptParenthesizedSeqOptList.
//
//	OptParenthesizedSeqOptList: '(' SeqOptList ')' | /* EMPTY */
func (p *Parser) parseOptParenthesizedSeqOptList() (*nodes.List, error) {
	if p.cur.Type != '(' {
		return nil, nil
	}
	p.advance()
	opts, err := p.parseSeqOptList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return opts, nil
}

// makeDefElem creates a DefElem with Loc=NoLoc().
func makeDefElem(name string, arg nodes.Node) *nodes.DefElem {
	return &nodes.DefElem{
		Defname: name,
		Arg:     arg,
		Loc:     nodes.NoLoc(),
	}
}

// parseTableConstraint parses a TableConstraint.
//
//	TableConstraint:
//	    CONSTRAINT name ConstraintElem
//	    | ConstraintElem
//
// The constraint spans from its first token, CONSTRAINT when it is named,
// to its last (gram.y sets location = @1).
func (p *Parser) parseTableConstraint() (*nodes.Constraint, error) {
	start := p.pos()
	named := p.cur.Type == CONSTRAINT
	var name string
	if named {
		p.advance()
		var err error
		if name, err = p.parseName(); err != nil {
			return nil, err
		}
	}
	c, err := p.parseConstraintElem()
	if err != nil {
		return nil, err
	}
	if c != nil {
		if named {
			c.Conname = name
		}
		c.Loc = nodes.Loc{Start: start, End: p.prev.End}
	}
	return c, nil
}

// parseConstraintElem parses a ConstraintElem.
//
//	ConstraintElem:
//	    UNIQUE opt_unique_null_treatment '(' columnList ')' opt_c_include opt_definition OptConsTableSpace ConstraintAttributeSpec
//	    | UNIQUE ExistingIndex ConstraintAttributeSpec
//	    | PRIMARY KEY '(' columnList ')' opt_c_include opt_definition OptConsTableSpace ConstraintAttributeSpec
//	    | PRIMARY KEY ExistingIndex ConstraintAttributeSpec
//	    | CHECK '(' a_expr ')' ConstraintAttributeSpec
//	    | FOREIGN KEY '(' columnList ')' REFERENCES qualified_name opt_column_list key_match key_actions ConstraintAttributeSpec
//	    | EXCLUDE ...
func (p *Parser) parseConstraintElem() (*nodes.Constraint, error) {
	switch p.cur.Type {
	case UNIQUE:
		p.advance()
		nullsNotDistinct, err := p.parseOptUniqueNullTreatment()
		if err != nil {
			return nil, err
		}

		if p.cur.Type == USING {
			return p.parseExistingIndexConstraint(nodes.CONSTR_UNIQUE, "UNIQUE")
		}

		keys, err := p.parseParenColumnList()
		if err != nil {
			return nil, err
		}
		including, err := p.parseOptCInclude()
		if err != nil {
			return nil, err
		}
		options, indexspace, err := p.parseIndexOptions()
		if err != nil {
			return nil, err
		}
		attrs, err := p.parseKindAttributes("UNIQUE", true, false, false)
		if err != nil {
			return nil, err
		}
		n := &nodes.Constraint{
			Contype:          nodes.CONSTR_UNIQUE,
			NullsNotDistinct: nullsNotDistinct,
			Keys:             keys,
			Including:        including,
			Options:          options,
			Indexspace:       indexspace,
			Loc:              nodes.NoLoc(),
			InitiallyValid:   true,
		}
		applyConstraintAttrs(n, attrs)
		return n, nil

	case PRIMARY:
		p.advance() // PRIMARY
		if _, err := p.expect(KEY); err != nil {
			return nil, err
		}

		if p.cur.Type == USING {
			return p.parseExistingIndexConstraint(nodes.CONSTR_PRIMARY, "PRIMARY KEY")
		}

		keys, err := p.parseParenColumnList()
		if err != nil {
			return nil, err
		}
		including, err := p.parseOptCInclude()
		if err != nil {
			return nil, err
		}
		options, indexspace, err := p.parseIndexOptions()
		if err != nil {
			return nil, err
		}
		attrs, err := p.parseKindAttributes("PRIMARY KEY", true, false, false)
		if err != nil {
			return nil, err
		}
		n := &nodes.Constraint{
			Contype:        nodes.CONSTR_PRIMARY,
			Keys:           keys,
			Including:      including,
			Options:        options,
			Indexspace:     indexspace,
			Loc:            nodes.NoLoc(),
			InitiallyValid: true,
		}
		applyConstraintAttrs(n, attrs)
		return n, nil

	case CHECK:
		// CHECK '(' a_expr ')' ConstraintAttributeSpec: the parentheses
		// belong to the production, not to the expression.
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		expr, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		attrs, err := p.parseKindAttributes("CHECK", false, true, true)
		if err != nil {
			return nil, err
		}
		n := &nodes.Constraint{
			Contype:        nodes.CONSTR_CHECK,
			RawExpr:        expr,
			Loc:            nodes.NoLoc(),
			InitiallyValid: true,
		}
		applyConstraintAttrs(n, attrs)
		return n, nil

	case FOREIGN:
		p.advance() // FOREIGN
		if _, err := p.expect(KEY); err != nil {
			return nil, err
		}
		fkAttrs, err := p.parseParenColumnList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(REFERENCES); err != nil {
			return nil, err
		}
		refStart := p.pos()
		refNames, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		refRv := makeRangeVarFromNames(refNames, nodes.Loc{Start: refStart, End: p.prev.End})
		pkAttrs, err := p.parseOptColumnList()
		if err != nil {
			return nil, err
		}
		matchType, err := p.parseKeyMatch()
		if err != nil {
			return nil, err
		}
		updAction, delAction, delSetCols, err := p.parseKeyActions()
		if err != nil {
			return nil, err
		}
		attrs, err := p.parseKindAttributes("FOREIGN KEY", true, true, false)
		if err != nil {
			return nil, err
		}
		n := &nodes.Constraint{
			Contype:        nodes.CONSTR_FOREIGN,
			FkAttrs:        fkAttrs,
			Pktable:        refRv,
			PkAttrs:        pkAttrs,
			FkMatchtype:    matchType,
			FkUpdaction:    updAction,
			FkDelaction:    delAction,
			FkDelsetcols:   delSetCols,
			Loc:            nodes.NoLoc(),
			InitiallyValid: true,
		}
		applyConstraintAttrs(n, attrs)
		return n, nil

	case EXCLUDE:
		return p.parseExclusionConstraint()
	}
	// exhaustive: gram.y ConstraintElem is CHECK, UNIQUE, PRIMARY KEY,
	// EXCLUDE, or FOREIGN KEY, so CONSTRAINT name followed by anything
	// else is refused here rather than dropped by the caller.
	if p.collectMode() {
		return nil, nil
	}
	return nil, p.syntaxErrorAtCur()
}

// parseExclusionConstraint parses EXCLUDE constraint.
func (p *Parser) parseExclusionConstraint() (*nodes.Constraint, error) {
	p.advance() // EXCLUDE
	am, err := p.parseAccessMethodClause()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	excl, err := p.parseExclusionConstraintList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	including, err := p.parseOptCInclude()
	if err != nil {
		return nil, err
	}
	options, indexspace, err := p.parseIndexOptions()
	if err != nil {
		return nil, err
	}
	where, err := p.parseWhereClauseOpt()
	if err != nil {
		return nil, err
	}
	attrs, err := p.parseKindAttributes("EXCLUDE", true, false, false)
	if err != nil {
		return nil, err
	}
	n := &nodes.Constraint{
		Contype:        nodes.CONSTR_EXCLUSION,
		AccessMethod:   am,
		Exclusions:     excl,
		Including:      including,
		Options:        options,
		Indexspace:     indexspace,
		WhereClause:    where,
		Loc:            nodes.NoLoc(),
		InitiallyValid: true,
	}
	applyConstraintAttrs(n, attrs)
	return n, nil
}

// parseAccessMethodClause parses access_method_clause.
//
//	access_method_clause: USING name | /* EMPTY */
func (p *Parser) parseAccessMethodClause() (string, error) {
	if p.cur.Type != USING {
		return "", nil
	}
	p.advance()
	return p.parseName()
}

// parseExclusionConstraintList parses ExclusionConstraintList.
//
//	ExclusionConstraintList: ExclusionConstraintElem
//	    | ExclusionConstraintList ',' ExclusionConstraintElem
func (p *Parser) parseExclusionConstraintList() (*nodes.List, error) {
	elem, err := p.parseExclusionConstraintElem()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{elem}
	for p.cur.Type == ',' {
		p.advance()
		elem, err = p.parseExclusionConstraintElem()
		if err != nil {
			return nil, err
		}
		items = append(items, elem)
	}
	return &nodes.List{Items: items}, nil
}

// parseExclusionConstraintElem parses a single exclusion constraint element.
//
//	ExclusionConstraintElem: index_elem WITH any_operator | index_elem WITH OPERATOR '(' any_operator ')'
func (p *Parser) parseExclusionConstraintElem() (*nodes.List, error) {
	idxElem, err := p.parseIndexElem()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(WITH); err != nil {
		return nil, err
	}
	var op *nodes.List
	if p.cur.Type == OPERATOR {
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		if op, err = p.parseAnyOperator(); err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	} else if op, err = p.parseAnyOperator(); err != nil {
		return nil, err
	}
	return &nodes.List{Items: []nodes.Node{idxElem, op}}, nil
}

// parseOptCollate parses opt_collate.
//
//	opt_collate: COLLATE any_name | /* EMPTY */
func (p *Parser) parseOptCollate() (*nodes.List, error) {
	if p.cur.Type != COLLATE {
		return nil, nil
	}
	p.advance()
	return p.parseAnyName()
}

// parseOptQualifiedName parses opt_qualified_name.
//
//	opt_qualified_name: any_name | /* EMPTY */
//
// Predicate is isColId only — parseAnyName calls parseColId
// downstream, which rejects TypeFuncNameKeyword. PG's any_name
// production is `ColId | ColId attrs`, also using ColId. The
// previous `|| p.isTypeFunctionName()` was dead code.
func (p *Parser) parseOptQualifiedName() (*nodes.List, error) {
	if !p.isColId() {
		return nil, nil
	}
	return p.parseAnyName()
}

// parseOptAscDesc parses opt_asc_desc.
func (p *Parser) parseOptAscDesc() int {
	switch p.cur.Type {
	case ASC:
		p.advance()
		return int(nodes.SORTBY_ASC)
	case DESC:
		p.advance()
		return int(nodes.SORTBY_DESC)
	default:
		return int(nodes.SORTBY_DEFAULT)
	}
}

// parseOptNullsOrder parses opt_nulls_order.
func (p *Parser) parseOptNullsOrder() int {
	if p.cur.Type == NULLS_LA {
		p.advance()
		switch p.cur.Type {
		case FIRST_P:
			p.advance()
			return int(nodes.SORTBY_NULLS_FIRST)
		case LAST_P:
			p.advance()
			return int(nodes.SORTBY_NULLS_LAST)
		}
	}
	return int(nodes.SORTBY_NULLS_DEFAULT)
}

// parseWhereClauseOpt parses an optional WHERE clause.
func (p *Parser) parseWhereClauseOpt() (nodes.Node, error) {
	if p.cur.Type != WHERE {
		return nil, nil
	}
	p.advance()
	result, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return result, nil
}

// parseTableLikeClause parses a LIKE clause in CREATE TABLE.
//
//	TableLikeClause: LIKE qualified_name TableLikeOptionList?
func (p *Parser) parseTableLikeClause() (*nodes.TableLikeClause, error) {
	likeLoc := p.pos()
	p.advance() // LIKE
	nameStart := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	if !p.collectMode() && names == nil {
		return nil, p.syntaxErrorAtCur()
	}
	rv := makeRangeVarFromNames(names, nodes.Loc{Start: nameStart, End: p.prev.End})
	var opts uint32
	for p.cur.Type == INCLUDING || p.cur.Type == EXCLUDING {
		isIncluding := p.cur.Type == INCLUDING
		p.advance()
		opt, err := p.parseTableLikeOption()
		if err != nil {
			return nil, err
		}
		if isIncluding {
			opts |= opt
		} else {
			opts &^= opt
		}
	}
	return &nodes.TableLikeClause{
		Relation: rv,
		Options:  opts,
		Loc:      nodes.Loc{Start: likeLoc, End: p.prev.End},
	}, nil
}

// parseTableLikeOption parses a single TableLikeOption keyword.
//
//	TableLikeOption: COMMENTS | COMPRESSION | CONSTRAINTS | DEFAULTS
//	    | IDENTITY_P | GENERATED | INDEXES | STATISTICS | STORAGE | ALL
func (p *Parser) parseTableLikeOption() (uint32, error) {
	var opt uint32
	switch p.cur.Type {
	case ALL:
		opt = 0xFFFFFFFF
	case COMMENTS:
		opt = 1
	case COMPRESSION:
		opt = 2
	case CONSTRAINTS:
		opt = 4
	case DEFAULTS:
		opt = 8
	case GENERATED:
		opt = 16
	case IDENTITY_P:
		opt = 32
	case INDEXES:
		opt = 64
	case STATISTICS:
		opt = 128
	case STORAGE:
		opt = 256
	default:
		return 0, p.syntaxErrorAtCur()
	}
	p.advance()
	return opt, nil
}

// parseTypedTableElementList parses a TypedTableElementList.
func (p *Parser) parseTypedTableElementList() (*nodes.List, error) {
	elem, err := p.parseTypedTableElement()
	if err != nil || elem == nil {
		return nil, err
	}
	items := []nodes.Node{elem}
	for p.cur.Type == ',' {
		p.advance()
		elem, err = p.parseTypedTableElement()
		if err != nil {
			return nil, err
		}
		if elem == nil {
			return nil, p.syntaxErrorAtCur()
		}
		items = append(items, elem)
	}
	return &nodes.List{Items: items}, nil
}

// parseTypedTableElement parses a single TypedTableElement.
//
//	TypedTableElement: columnOptions | TableConstraint
func (p *Parser) parseTypedTableElement() (nodes.Node, error) {
	if p.startsTableConstraint() {
		c, err := p.parseTableConstraint()
		if err != nil || c == nil {
			return nil, err
		}
		return c, nil
	}
	c, err := p.parseColumnOptions()
	if err != nil || c == nil {
		return nil, err
	}
	return c, nil
}

// parseColumnOptions parses columnOptions.
//
//	columnOptions: ColId opt_column_constraints | ColId WITH OPTIONS opt_column_constraints
func (p *Parser) parseColumnOptions() (*nodes.ColumnDef, error) {
	start := p.pos()
	colname, err := p.parseColId()
	if err != nil {
		return nil, err
	}

	n := &nodes.ColumnDef{
		Colname:  colname,
		TypeName: nil,
		IsLocal:  true,
		Loc:      nodes.NoLoc(),
	}

	if p.cur.Type == WITH {
		p.advance()
		if _, err := p.expect(OPTIONS); err != nil {
			return nil, err
		}
	}

	qualList, err := p.parseOptColumnConstraints()
	if err != nil {
		return nil, err
	}
	splitColQualList(qualList, n)
	n.Loc = nodes.Loc{Start: start, End: p.prev.End}

	return n, nil
}

// splitColQualList distributes constraint items to a ColumnDef.
func splitColQualList(qualList *nodes.List, coldef *nodes.ColumnDef) {
	if qualList == nil {
		return
	}
	var constraints []nodes.Node
	for _, item := range qualList.Items {
		if cc, ok := item.(*nodes.CollateClause); ok {
			coldef.CollClause = cc
		} else if de, ok := item.(*nodes.DefElem); ok && de.Defname == "compression" {
			if s, ok := de.Arg.(*nodes.String); ok {
				coldef.Compression = s.Str
			}
		} else if de, ok := item.(*nodes.DefElem); ok && de.Defname == "encode" {
			if s, ok := de.Arg.(*nodes.String); ok {
				coldef.Compression = s.Str
			}
		} else if de, ok := item.(*nodes.DefElem); ok && de.Defname == "storage" {
			if s, ok := de.Arg.(*nodes.String); ok {
				coldef.StorageName = s.Str
			}
		} else if de, ok := item.(*nodes.DefElem); ok && de.Defname == "fdwoptions" {
			if l, ok := de.Arg.(*nodes.List); ok {
				coldef.Fdwoptions = l
			}
		} else {
			constraints = append(constraints, item)
		}
	}
	if len(constraints) > 0 {
		coldef.Constraints = &nodes.List{Items: constraints}
	}
}

// parseOptColumnList parses opt_column_list.
//
//	opt_column_list: '(' columnList ')' | /* EMPTY */
func (p *Parser) parseOptColumnList() (*nodes.List, error) {
	if p.cur.Type != '(' {
		return nil, nil
	}
	p.advance()
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil, nil
	}
	list, err := p.parseColumnList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return list, nil
}

// parseParenColumnList parses '(' columnList ')'.
func (p *Parser) parseParenColumnList() (*nodes.List, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	list, err := p.parseColumnList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return list, nil
}

// parseColumnList parses a comma-separated list of column names.
//
//	columnList: columnElem | columnList ',' columnElem
//
// Every element is a name, so the list is never empty and a comma with no
// name after it is a syntax error at the token after the comma. In
// collect mode it offers column names and returns no list.
func (p *Parser) parseColumnList() (*nodes.List, error) {
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil, nil
	}
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{&nodes.String{Str: name}}
	for p.cur.Type == ',' {
		p.advance()
		name, err = p.parseName()
		if err != nil {
			return nil, err
		}
		items = append(items, &nodes.String{Str: name})
	}
	return &nodes.List{Items: items}, nil
}

// parseKeyMatch parses key_match.
//
//	key_match: MATCH FULL | MATCH PARTIAL | MATCH SIMPLE | /* EMPTY */
func (p *Parser) parseKeyMatch() (byte, error) {
	if p.cur.Type != MATCH {
		return 's', nil // default SIMPLE
	}
	matchLoc := p.pos()
	p.advance()
	switch p.cur.Type {
	case FULL:
		p.advance()
		return 'f', nil
	case PARTIAL:
		// gram.y refuses MATCH PARTIAL in its action.
		p.advance()
		return 0, &ParseError{Code: "0A000", Message: "MATCH PARTIAL not yet implemented", Position: matchLoc}
	case SIMPLE:
		p.advance()
		return 's', nil
	default:
		return 0, p.syntaxErrorAtCur()
	}
}

// parseKeyActions parses key_actions, returning (update action, delete action, delete set cols).
//
//	key_actions:
//	    key_update | key_delete | key_update key_delete | key_delete key_update | /* EMPTY */
//	key_update: ON UPDATE key_action
//	key_delete: ON DELETE_P key_action
func (p *Parser) parseKeyActions() (updAction byte, delAction byte, delSetCols *nodes.List, err error) {
	updAction = 'a' // default NO ACTION
	delAction = 'a' // default NO ACTION

	// Nothing else after a key's references starts with ON, so ON leads
	// a key action until both are read.
	var sawUpdate, sawDelete bool
	for p.cur.Type == ON && !(sawUpdate && sawDelete) {
		onLoc := p.pos()
		p.advance() // ON
		switch {
		case p.cur.Type == UPDATE && !sawUpdate:
			p.advance() // UPDATE
			var cols *nodes.List
			if updAction, err = p.parseKeyActionType(&cols); err != nil {
				return
			}
			if cols != nil {
				// gram.y's key_update refuses a column list.
				what := "SET NULL"
				if updAction == 'd' {
					what = "SET DEFAULT"
				}
				err = &ParseError{Code: "0A000", Message: "a column list with " + what + " is only supported for ON DELETE actions", Position: onLoc}
				return
			}
			sawUpdate = true
		case p.cur.Type == DELETE_P && !sawDelete:
			p.advance() // DELETE
			if delAction, err = p.parseKeyActionType(&delSetCols); err != nil {
				return
			}
			sawDelete = true
		default:
			if p.collectMode() {
				if !sawUpdate {
					p.addTokenCandidate(UPDATE)
				}
				if !sawDelete {
					p.addTokenCandidate(DELETE_P)
				}
				err = errCollecting
				return
			}
			err = p.syntaxErrorAtCur()
			return
		}
	}
	return
}

// parseKeyActionType parses a key_action, and the column list of SET NULL
// or SET DEFAULT into setCols.
//
//	key_action: NO ACTION | RESTRICT | CASCADE
//	    | SET NULL_P opt_column_list | SET DEFAULT opt_column_list
func (p *Parser) parseKeyActionType(setCols **nodes.List) (byte, error) {
	switch p.cur.Type {
	case NO:
		p.advance()
		if _, err := p.expect(ACTION); err != nil {
			return 0, err
		}
		return 'a', nil
	case RESTRICT:
		p.advance()
		return 'r', nil
	case CASCADE:
		p.advance()
		return 'c', nil
	case SET:
		p.advance()
		var action byte
		switch p.cur.Type {
		case NULL_P:
			action = 'n'
		case DEFAULT:
			action = 'd'
		default:
			if p.collectMode() {
				p.addTokenCandidate(NULL_P)
				p.addTokenCandidate(DEFAULT)
				return 0, errCollecting
			}
			return 0, p.syntaxErrorAtCur()
		}
		p.advance()
		cols, err := p.parseOptColumnList()
		if err != nil {
			return 0, err
		}
		*setCols = cols
		return action, nil
	}
	if p.collectMode() {
		for _, t := range []int{NO, RESTRICT, CASCADE, SET} {
			p.addTokenCandidate(t)
		}
		return 0, errCollecting
	}
	return 0, p.syntaxErrorAtCur()
}

// parseConstraintAttributeSpec parses ConstraintAttributeSpec and returns
// its CAS bits.
//
//	ConstraintAttributeSpec: /* EMPTY */ | ConstraintAttributeSpec ConstraintAttributeElem
func (p *Parser) parseConstraintAttributeSpec() (int64, error) {
	return p.parseConstraintAttrs(true)
}

// parseConstraintAttrs parses constraint attributes. validity tells whether
// NOT VALID and NO INHERIT belong to the production; when they do not, they
// end it; when they do, NOT and NO start nothing else. Like gram.y's
// ConstraintAttributeSpec, it refuses INITIALLY DEFERRED with NOT
// DEFERRABLE, and DEFERRABLE with NOT DEFERRABLE or INITIALLY IMMEDIATE
// with INITIALLY DEFERRED, at the attribute that conflicts.
//
// pg: src/backend/parser/gram.y — ConstraintAttributeSpec
func (p *Parser) parseConstraintAttrs(validity bool) (int64, error) {
	var attrs int64
	for {
		loc := p.cur.Loc
		var bit int64
		switch p.cur.Type {
		case NOT:
			switch p.peekNext().Type {
			case DEFERRABLE:
				bit = int64(nodes.CAS_NOT_DEFERRABLE)
			case VALID:
				if !validity {
					return attrs, nil
				}
				bit = int64(nodes.CAS_NOT_VALID)
			default:
				// In ConstraintAttributeSpec, NOT starts an attribute and
				// nothing else, so the error is at the token after it.
				if !validity {
					return attrs, nil
				}
				p.advance()
				if p.collectMode() {
					p.addTokenCandidate(DEFERRABLE)
					p.addTokenCandidate(VALID)
					return attrs, nil
				}
				return 0, p.syntaxErrorAtCur()
			}
			p.advance()
			p.advance()
		case DEFERRABLE:
			p.advance()
			bit = int64(nodes.CAS_DEFERRABLE)
		case INITIALLY:
			p.advance()
			switch p.cur.Type {
			case IMMEDIATE:
				bit = int64(nodes.CAS_INITIALLY_IMMEDIATE)
			case DEFERRED:
				bit = int64(nodes.CAS_INITIALLY_DEFERRED)
			default:
				if p.collectMode() {
					return attrs, nil
				}
				return 0, p.syntaxErrorAtCur()
			}
			p.advance()
		case NO:
			if !validity {
				return attrs, nil
			}
			p.advance()
			if _, err := p.expect(INHERIT); err != nil {
				return 0, err
			}
			bit = int64(nodes.CAS_NO_INHERIT)
		default:
			return attrs, nil
		}
		attrs |= bit
		notDeferrable, deferrable := int64(nodes.CAS_NOT_DEFERRABLE), int64(nodes.CAS_DEFERRABLE)
		immediate, deferred := int64(nodes.CAS_INITIALLY_IMMEDIATE), int64(nodes.CAS_INITIALLY_DEFERRED)
		switch {
		case attrs&(notDeferrable|deferred) == notDeferrable|deferred:
			return 0, &ParseError{Message: "constraint declared INITIALLY DEFERRED must be DEFERRABLE", Position: loc}
		case attrs&(notDeferrable|deferrable) == notDeferrable|deferrable, attrs&(immediate|deferred) == immediate|deferred:
			return 0, &ParseError{Message: "conflicting constraint properties", Position: loc}
		}
	}
}

// processCASbits refuses the attributes a constraint of the kind cannot
// take: deferrable, notValid, and noInherit tell which it can. The server
// reports no position for these: the location gram.y passes is
// ConstraintAttributeSpec's, which begins with the empty alternative and so
// is always -1.
//
// pg: src/backend/parser/gram.y — processCASbits
func (p *Parser) processCASbits(attrs int64, constrType string, deferrable, notValid, noInherit bool) error {
	if p.collectMode() {
		return nil
	}
	refuse := func(what string) error {
		return &ParseError{Code: "0A000", Message: constrType + " constraints cannot be marked " + what, Position: -1}
	}
	switch {
	case !deferrable && attrs&int64(nodes.CAS_DEFERRABLE|nodes.CAS_INITIALLY_DEFERRED) != 0:
		return refuse("DEFERRABLE")
	case !notValid && attrs&int64(nodes.CAS_NOT_VALID) != 0:
		return refuse("NOT VALID")
	case !noInherit && attrs&int64(nodes.CAS_NO_INHERIT) != 0:
		return refuse("NO INHERIT")
	}
	return nil
}

// parseKindAttributes parses the ConstraintAttributeSpec of a constraint of
// the kind and refuses the attributes the kind cannot take.
func (p *Parser) parseKindAttributes(constrType string, deferrable, notValid, noInherit bool) (int64, error) {
	attrs, err := p.parseConstraintAttributeSpec()
	if err != nil {
		return 0, err
	}
	if err := p.processCASbits(attrs, constrType, deferrable, notValid, noInherit); err != nil {
		return 0, err
	}
	return attrs, nil
}

// parseCreateAssertionStmt reads CREATE ASSERTION, which the grammar takes
// and its action then refuses, with no position, as not implemented. It
// returns that refusal, or the syntax error the grammar finds first.
// CREATE has been consumed; the current token is ASSERTION.
//
//	CreateAssertionStmt:
//	    CREATE ASSERTION any_name CHECK '(' a_expr ')' ConstraintAttributeSpec
//
// pg: src/backend/parser/gram.y — CreateAssertionStmt
func (p *Parser) parseCreateAssertionStmt() error {
	p.advance() // consume ASSERTION
	if _, err := p.parseAnyName(); err != nil {
		return err
	}
	if _, err := p.expect(CHECK); err != nil {
		return err
	}
	if _, err := p.parseParenRequiredAExpr(); err != nil {
		return err
	}
	if _, err := p.parseConstraintAttributeSpec(); err != nil {
		return err
	}
	if p.collectMode() {
		return errCollecting
	}
	return &ParseError{Code: "0A000", Message: "CREATE ASSERTION is not yet implemented", Position: -1}
}

// applyConstraintAttrs applies constraint attribute flags to a Constraint node.
func applyConstraintAttrs(n *nodes.Constraint, attrs int64) {
	if attrs&int64(nodes.CAS_DEFERRABLE) != 0 {
		n.Deferrable = true
	}
	if attrs&int64(nodes.CAS_INITIALLY_DEFERRED) != 0 {
		n.Initdeferred = true
		n.Deferrable = true
	}
	if attrs&int64(nodes.CAS_NOT_VALID) != 0 {
		n.SkipValidation = true
		n.InitiallyValid = false
	}
	if attrs&int64(nodes.CAS_NO_INHERIT) != 0 {
		n.IsNoInherit = true
	}
}

// parseOptUniqueNullTreatment parses opt_unique_null_treatment.
//
//	opt_unique_null_treatment: NULLS_P DISTINCT | NULLS_P NOT DISTINCT | /* EMPTY */
//
// PG's gram.y keeps the plain NULLS_P token here. base_yylex only
// reclassifies NULLS_P -> NULLS_LA when followed by FIRST_P/LAST_P
// (for ORDER BY x NULLS FIRST/LAST), which matches omni's
// reclassification rule. UNIQUE NULLS [NOT] DISTINCT therefore sees
// the raw NULLS_P token and must probe for it directly.
func (p *Parser) parseOptUniqueNullTreatment() (bool, error) {
	if p.cur.Type != NULLS_P {
		return false, nil
	}
	p.advance()
	notDistinct := p.cur.Type == NOT
	if notDistinct {
		p.advance()
	}
	if _, err := p.expect(DISTINCT); err != nil {
		return false, err
	}
	return notDistinct, nil
}

// parseOptCInclude parses opt_c_include.
//
//	opt_c_include: INCLUDE '(' columnList ')' | /* EMPTY */
func (p *Parser) parseOptCInclude() (*nodes.List, error) {
	if p.cur.Type != INCLUDE {
		return nil, nil
	}
	p.advance()
	return p.parseParenColumnList()
}

// parseOptDefinition parses opt_definition.
//
//	opt_definition: WITH definition | /* EMPTY */
//	definition: '(' def_list ')'
//
// Nowhere opt_definition stands does anything else start with WITH, so it
// is read as the clause.
func (p *Parser) parseOptDefinition() (*nodes.List, error) {
	if p.cur.Type != WITH {
		return nil, nil
	}
	p.advance() // WITH
	return p.parseDefinition()
}

// parseDefList parses def_list.
//
//	def_list: def_elem | def_list ',' def_elem
// parseDefList, parseDefElem, parseDefArg are in define.go

// parseOptConsTableSpace parses OptConsTableSpace.
//
//	OptConsTableSpace: USING INDEX TABLESPACE name | /* EMPTY */
//
// Nothing else after a key's options starts with USING, so the clause is
// read from it.
func (p *Parser) parseOptConsTableSpace() (string, error) {
	if p.cur.Type != USING {
		return "", nil
	}
	p.advance() // USING
	if _, err := p.expect(INDEX); err != nil {
		return "", err
	}
	if _, err := p.expect(TABLESPACE); err != nil {
		return "", err
	}
	return p.parseName()
}

// parseExistingIndexConstraint parses the rest of a UNIQUE or PRIMARY KEY
// table constraint on an existing index. The current token is USING, which
// starts nothing else there.
//
//	UNIQUE ExistingIndex ConstraintAttributeSpec
//	PRIMARY KEY ExistingIndex ConstraintAttributeSpec
//	ExistingIndex: USING INDEX name
func (p *Parser) parseExistingIndexConstraint(contype nodes.ConstrType, constrType string) (*nodes.Constraint, error) {
	p.advance() // USING
	if _, err := p.expect(INDEX); err != nil {
		return nil, err
	}
	if p.collectMode() {
		p.addRuleCandidate("qualified_name")
		return nil, nil
	}
	idxName, err := p.parseName()
	if err != nil {
		return nil, err
	}
	attrs, err := p.parseKindAttributes(constrType, true, false, false)
	if err != nil {
		return nil, err
	}
	n := &nodes.Constraint{
		Contype:        contype,
		Indexname:      idxName,
		Loc:            nodes.NoLoc(),
		InitiallyValid: true,
	}
	applyConstraintAttrs(n, attrs)
	return n, nil
}

// parseIndexOptions parses the opt_definition OptConsTableSpace of a key
// or exclusion constraint.
func (p *Parser) parseIndexOptions() (*nodes.List, string, error) {
	options, err := p.parseOptDefinition()
	if err != nil {
		return nil, "", err
	}
	indexspace, err := p.parseOptConsTableSpace()
	if err != nil {
		return nil, "", err
	}
	return options, indexspace, nil
}

// parseNoInherit parses opt_no_inherit. Nothing else after a column's
// CHECK starts with NO, so the clause is read from it.
//
//	opt_no_inherit: NO INHERIT | /* EMPTY */
func (p *Parser) parseNoInherit() (bool, error) {
	if p.cur.Type != NO {
		return false, nil
	}
	p.advance()
	if _, err := p.expect(INHERIT); err != nil {
		return false, err
	}
	return true, nil
}

// parseOptInherit parses OptInherit.
//
//	OptInherit: INHERITS '(' qualified_name_list ')' | /* EMPTY */
func (p *Parser) parseOptInherit() (*nodes.List, error) {
	if p.cur.Type != INHERITS {
		return nil, nil
	}
	p.advance()
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	list, err := p.parseQualifiedNameList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return list, nil
}

// parseOptPartitionSpec parses OptPartitionSpec.
//
//	OptPartitionSpec: PartitionSpec | /* EMPTY */
func (p *Parser) parseOptPartitionSpec() (*nodes.PartitionSpec, error) {
	if p.cur.Type == PARTITION {
		return p.parsePartitionSpec()
	}
	return nil, nil
}

// parsePartitionSpec parses PartitionSpec.
//
//	PartitionSpec: PARTITION BY ColId '(' part_params ')'
func (p *Parser) parsePartitionSpec() (*nodes.PartitionSpec, error) {
	p.advance() // PARTITION
	if _, err := p.expect(BY); err != nil {
		return nil, err
	}
	strategy, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	params, err := p.parsePartParams()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	code, err := parsePartitionStrategy(strategy)
	if err != nil {
		return nil, err
	}
	return &nodes.PartitionSpec{
		Strategy:   code,
		PartParams: params,
		Loc:        nodes.NoLoc(),
	}, nil
}

// parsePartitionStrategy converts partition strategy name to internal code.
// gram.y's parsePartitionStrategy refuses any other name once the
// PartitionSpec is read, with no position.
func parsePartitionStrategy(strategy string) (string, error) {
	switch strings.ToLower(strategy) {
	case "list":
		return "l", nil
	case "range":
		return "r", nil
	case "hash":
		return "h", nil
	}
	return "", &ParseError{Code: "22023", Message: fmt.Sprintf("unrecognized partitioning strategy \"%s\"", strategy), Position: -1}
}

// parsePartParams parses part_params.
//
//	part_params: part_elem | part_params ',' part_elem
func (p *Parser) parsePartParams() (*nodes.List, error) {
	elem, err := p.parsePartElem()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{elem}
	for p.cur.Type == ',' {
		p.advance()
		elem, err = p.parsePartElem()
		if err != nil {
			return nil, err
		}
		items = append(items, elem)
	}
	return &nodes.List{Items: items}, nil
}

// parsePartElem parses a single partition key element.
//
//	part_elem:
//	    ColId opt_collate opt_qualified_name
//	    | func_expr_windowless opt_collate opt_qualified_name
//	    | '(' a_expr ')' opt_collate opt_qualified_name
func (p *Parser) parsePartElem() (*nodes.PartitionElem, error) {
	if p.cur.Type == '(' {
		p.advance()
		expr, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		collation, opclass, err := p.parsePartElemTail()
		if err != nil {
			return nil, err
		}
		return &nodes.PartitionElem{
			Expr:      expr,
			Collation: collation,
			Opclass:   opclass,
			Loc:       nodes.NoLoc(),
		}, nil
	}

	// Try as ColId first (checking if next is not '(' which would make it a function)
	if p.isColId() {
		next := p.peekNext()
		if next.Type != '(' {
			name, err := p.parseColId()
			if err != nil {
				return nil, err
			}
			collation, opclass, err := p.parsePartElemTail()
			if err != nil {
				return nil, err
			}
			return &nodes.PartitionElem{
				Name:      name,
				Collation: collation,
				Opclass:   opclass,
				Loc:       nodes.NoLoc(),
			}, nil
		}
	}

	// func_expr_windowless
	expr, err := p.parseFuncExprWindowless()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}
	collation, opclass, err := p.parsePartElemTail()
	if err != nil {
		return nil, err
	}
	return &nodes.PartitionElem{
		Expr:      expr,
		Collation: collation,
		Opclass:   opclass,
		Loc:       nodes.NoLoc(),
	}, nil
}

// parsePartElemTail parses the opt_collate opt_qualified_name that ends a
// part_elem.
func (p *Parser) parsePartElemTail() (*nodes.List, *nodes.List, error) {
	collation, err := p.parseOptCollate()
	if err != nil {
		return nil, nil, err
	}
	opclass, err := p.parseOptQualifiedName()
	if err != nil {
		return nil, nil, err
	}
	return collation, opclass, nil
}

// parseOptAccessMethod parses OptAccessMethod.
//
//	OptAccessMethod: USING name | /* EMPTY */
func (p *Parser) parseOptAccessMethod() (string, error) {
	if p.cur.Type == USING {
		p.advance()
		return p.parseName()
	}
	return "", nil
}

// parseOptWith parses OptWith for CREATE TABLE. WITH starts nothing else
// there, so it is read as the clause; WITH OIDS has been a syntax error
// since PostgreSQL 12.
//
//	OptWith: WITH reloptions | WITHOUT OIDS | /* EMPTY */
func (p *Parser) parseOptWith() (*nodes.List, error) {
	switch p.cur.Type {
	case WITH:
		p.advance() // WITH
		return p.parseReloptions()
	case WITHOUT:
		p.advance() // WITHOUT
		if _, err := p.expect(OIDS); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// parseReloptions parses reloptions.
//
//	reloptions: '(' reloption_list ')'
func (p *Parser) parseReloptions() (*nodes.List, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	list, err := p.parseReloptionList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return list, nil
}

// parseReloptionList parses reloption_list.
func (p *Parser) parseReloptionList() (*nodes.List, error) {
	list := &nodes.List{}
	for {
		elem, err := p.parseReloptionElem()
		if err != nil {
			return nil, err
		}
		list.Items = append(list.Items, elem)
		if p.cur.Type != ',' {
			return list, nil
		}
		p.advance()
	}
}

// parseReloptionElem parses a single reloption_elem.
//
//	reloption_elem:
//	    ColLabel '=' def_arg
//	    | ColLabel
//	    | ColLabel '.' ColLabel '=' def_arg
//	    | ColLabel '.' ColLabel
func (p *Parser) parseReloptionElem() (*nodes.DefElem, error) {
	elem := &nodes.DefElem{Loc: nodes.NoLoc()}
	label, err := p.parseColLabel()
	if err != nil {
		return nil, err
	}
	elem.Defname = label
	if p.cur.Type == '.' {
		p.advance()
		sublabel, err := p.parseColLabel()
		if err != nil {
			return nil, err
		}
		elem.Defnamespace, elem.Defname = label, sublabel
	}
	if p.cur.Type == '=' {
		p.advance()
		arg, err := p.parseDefArg()
		if err != nil {
			return nil, err
		}
		elem.Arg = arg
	}
	return elem, nil
}

// parseOnCommitOption parses OnCommitOption.
//
//	OnCommitOption: ON COMMIT DROP | ON COMMIT DELETE ROWS | ON COMMIT PRESERVE ROWS | /* EMPTY */
//
// ON starts nothing else where the option stands, so it is read from ON.
func (p *Parser) parseOnCommitOption() (nodes.OnCommitAction, error) {
	if p.cur.Type != ON {
		return nodes.ONCOMMIT_NOOP, nil
	}
	p.advance() // ON
	if _, err := p.expect(COMMIT); err != nil {
		return 0, err
	}
	switch p.cur.Type {
	case DROP:
		p.advance()
		return nodes.ONCOMMIT_DROP, nil
	case DELETE_P:
		p.advance()
		if _, err := p.expect(ROWS); err != nil {
			return 0, err
		}
		return nodes.ONCOMMIT_DELETE_ROWS, nil
	case PRESERVE:
		p.advance()
		if _, err := p.expect(ROWS); err != nil {
			return 0, err
		}
		return nodes.ONCOMMIT_PRESERVE_ROWS, nil
	}
	return 0, p.syntaxErrorAtCur()
}

// parseOptTableSpace parses OptTableSpace.
//
//	OptTableSpace: TABLESPACE name | /* EMPTY */
func (p *Parser) parseOptTableSpace() (string, error) {
	if p.cur.Type == TABLESPACE {
		p.advance()
		return p.parseName()
	}
	return "", nil
}

// parseForValues parses ForValues (partition bound spec).
//
//	ForValues: PartitionBoundSpec | DEFAULT
func (p *Parser) parseForValues() (*nodes.PartitionBoundSpec, error) {
	if p.cur.Type == DEFAULT {
		p.advance()
		return &nodes.PartitionBoundSpec{
			IsDefault: true,
			Loc:       nodes.NoLoc(),
		}, nil
	}
	return p.parsePartitionBoundSpec()
}

// parsePartitionBoundSpec parses PartitionBoundSpec.
//
//	PartitionBoundSpec:
//	    FOR VALUES IN '(' expr_list ')'
//	    | FOR VALUES FROM '(' expr_list ')' TO '(' expr_list ')'
//	    | FOR VALUES WITH '(' hash_partbound ')'
func (p *Parser) parsePartitionBoundSpec() (*nodes.PartitionBoundSpec, error) {
	if _, err := p.expect(FOR); err != nil {
		return nil, err
	}
	if _, err := p.expect(VALUES); err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case IN_P:
		p.advance()
		list, err := p.parseParenExprList()
		if err != nil {
			return nil, err
		}
		return &nodes.PartitionBoundSpec{
			Strategy:   'l',
			Listdatums: list,
			Loc:        nodes.NoLoc(),
		}, nil
	case FROM:
		p.advance()
		lower, err := p.parseParenExprList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		upper, err := p.parseParenExprList()
		if err != nil {
			return nil, err
		}
		return &nodes.PartitionBoundSpec{
			Strategy:    'r',
			Lowerdatums: lower,
			Upperdatums: upper,
			Loc:         nodes.NoLoc(),
		}, nil
	case WITH:
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		spec := &nodes.PartitionBoundSpec{
			Strategy:  'h',
			Modulus:   -1,
			Remainder: -1,
			Loc:       nodes.NoLoc(),
		}
		// hash_partbound: hash_partbound_elem [, ...]
		var unknown *ParseError
		for {
			// hash_partbound_elem: NonReservedWord Iconst
			nameTok := p.cur
			name, err := p.parseNonReservedWord()
			if err != nil {
				return nil, err
			}
			val, err := p.expect(ICONST)
			if err != nil {
				return nil, err
			}
			switch name {
			case "modulus":
				spec.Modulus = int(val.Ival)
			case "remainder":
				spec.Remainder = int(val.Ival)
			default:
				if unknown == nil {
					unknown = p.syntaxErrorAtTok(nameTok)
				}
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		// gram.y's action rejects any other name as a syntax error at it,
		// once the whole bound has been read.
		if unknown != nil {
			return nil, unknown
		}
		return spec, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseParenExprList parses '(' expr_list ')'.
func (p *Parser) parseParenExprList() (*nodes.List, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	list, err := p.parseExprListFull()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return list, nil
}
