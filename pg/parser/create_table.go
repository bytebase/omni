package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
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
	relpersistence := p.parseOptTemp()

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
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rv := makeRangeVarFromNames(names)
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

	// OptInherit
	stmt.InhRelations = p.parseOptInherit()

	// OptPartitionSpec
	stmt.Partspec = p.parseOptPartitionSpec()

	// OptAccessMethod (USING name) -- but only for table AM, not constraint/index
	stmt.AccessMethod = p.parseOptAccessMethod()

	// OptWith
	stmt.Options = p.parseOptWith()

	// OnCommitOption
	stmt.OnCommit = p.parseOnCommitOption()

	// OptTableSpace
	stmt.Tablespacename = p.parseOptTableSpace()

	stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End}
	return stmt, nil
}

// parseOptTemp parses OptTemp production.
//
//	OptTemp: TEMPORARY | TEMP | LOCAL TEMPORARY | LOCAL TEMP | UNLOGGED | /* EMPTY */
func (p *Parser) parseOptTemp() byte {
	switch p.cur.Type {
	case TEMPORARY, TEMP:
		p.advance()
		return 't'
	case LOCAL:
		p.advance()
		if p.cur.Type == TEMPORARY || p.cur.Type == TEMP {
			p.advance()
		}
		return 't'
	case UNLOGGED:
		p.advance()
		return 'u'
	default:
		return 'p'
	}
}

// parseCreateStmtPartitionOf parses the PARTITION OF variant of CREATE TABLE.
func (p *Parser) parseCreateStmtPartitionOf(stmt *nodes.CreateStmt, relpersistence byte) (*nodes.CreateStmt, error) {
	// parent table name
	parentNames, err := p.parseQualifiedName()
	if err != nil {
		return stmt, err
	}
	parentRv := makeRangeVarFromNames(parentNames)
	parentRv.Relpersistence = relpersistence

	// Optional '(' TypedTableElementList ')'
	if p.cur.Type == '(' {
		p.advance()
		stmt.TableElts = p.parseTypedTableElementList()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	// ForValues
	stmt.InhRelations = &nodes.List{Items: []nodes.Node{parentRv}}
	stmt.Partbound = p.parseForValues()

	// OptPartitionSpec
	stmt.Partspec = p.parseOptPartitionSpec()

	// OptAccessMethod
	stmt.AccessMethod = p.parseOptAccessMethod()

	// OptWith
	stmt.Options = p.parseOptWith()

	// OnCommitOption
	stmt.OnCommit = p.parseOnCommitOption()

	// OptTableSpace
	stmt.Tablespacename = p.parseOptTableSpace()

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
		stmt.TableElts = p.parseTypedTableElementList()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	// OptAccessMethod
	stmt.AccessMethod = p.parseOptAccessMethod()

	// OptWith
	stmt.Options = p.parseOptWith()

	// OnCommitOption
	stmt.OnCommit = p.parseOnCommitOption()

	// OptTableSpace
	stmt.Tablespacename = p.parseOptTableSpace()

	return stmt, nil
}

// makeTypeNameFromNameList creates a TypeName from an any_name list.
func makeTypeNameFromNameList(names *nodes.List) *nodes.TypeName {
	return &nodes.TypeName{
		Names:    names,
		Loc: nodes.NoLoc(),
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
			break
		}
		items = append(items, elem)
	}
	return &nodes.List{Items: items}, nil
}

// parseTableElement parses a single table element.
//
//	TableElement: columnDef | TableConstraint | TableLikeClause
func (p *Parser) parseTableElement() (nodes.Node, error) {
	switch p.cur.Type {
	case CONSTRAINT:
		return p.parseTableConstraint(), nil
	case CHECK, UNIQUE, PRIMARY, FOREIGN, EXCLUDE:
		return p.parseTableConstraint(), nil
	case LIKE:
		return p.parseTableLikeClause()
	default:
		// columnDef: ColId Typename opt_column_constraints
		return p.parseColumnDef()
	}
}

// parseColumnDef parses a column definition.
//
// Ref: gram.y columnDef
//
//	columnDef: ColId Typename opt_column_constraints
func (p *Parser) parseColumnDef() (*nodes.ColumnDef, error) {
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
		Loc: nodes.NoLoc(),
	}

	// opt_column_constraints
	qualList, err := p.parseOptColumnConstraints()
	if err != nil {
		return nil, err
	}
	splitColQualList(qualList, n)

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
func (p *Parser) parseColConstraint() (nodes.Node, error) {
	switch p.cur.Type {
	case CONSTRAINT:
		p.advance()
		name, _ := p.parseName()
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
				Contype:  nodes.CONSTR_NOTNULL,
				Loc: nodes.NoLoc(),
			}, nil
		}
		if next.Type == DEFERRABLE {
			p.advance() // NOT
			p.advance() // DEFERRABLE
			return &nodes.Constraint{
				Contype:  nodes.CONSTR_ATTR_NOT_DEFERRABLE,
				Loc: nodes.NoLoc(),
			}, nil
		}
		return nil, nil
	case NULL_P:
		p.advance()
		return &nodes.Constraint{
			Contype:  nodes.CONSTR_NULL,
			Loc: nodes.NoLoc(),
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
			Loc: nodes.NoLoc(),
		}, nil
	case COMPRESSION:
		p.advance()
		id, err := p.parseColId()
		if err != nil {
			return nil, err
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
		p.advance()
		id, err := p.parseColId()
		if err != nil {
			return nil, err
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
			Contype:  nodes.CONSTR_ATTR_DEFERRABLE,
			Loc: nodes.NoLoc(),
		}, nil
	case INITIALLY:
		p.advance()
		if p.cur.Type == DEFERRED {
			p.advance()
			return &nodes.Constraint{
				Contype:  nodes.CONSTR_ATTR_DEFERRED,
				Loc: nodes.NoLoc(),
			}, nil
		}
		if p.cur.Type == IMMEDIATE {
			p.advance()
			return &nodes.Constraint{
				Contype:  nodes.CONSTR_ATTR_IMMEDIATE,
				Loc: nodes.NoLoc(),
			}, nil
		}
		return nil, nil
	default:
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
		p.expect(NULL_P)
		return &nodes.Constraint{
			Contype:  nodes.CONSTR_NOTNULL,
			Loc: nodes.NoLoc(),
		}, nil
	case NULL_P:
		p.advance()
		return &nodes.Constraint{
			Contype:  nodes.CONSTR_NULL,
			Loc: nodes.NoLoc(),
		}, nil
	case UNIQUE:
		p.advance()
		nullsNotDistinct := p.parseOptUniqueNullTreatment()
		options := p.parseOptDefinition()
		indexspace := p.parseOptConsTableSpace()
		return &nodes.Constraint{
			Contype:          nodes.CONSTR_UNIQUE,
			NullsNotDistinct: nullsNotDistinct,
			Options:          options,
			Indexspace:        indexspace,
			Loc: nodes.NoLoc(),
		}, nil
	case PRIMARY:
		p.advance() // PRIMARY
		p.expect(KEY)
		options := p.parseOptDefinition()
		indexspace := p.parseOptConsTableSpace()
		return &nodes.Constraint{
			Contype:   nodes.CONSTR_PRIMARY,
			Options:   options,
			Indexspace: indexspace,
			Loc: nodes.NoLoc(),
		}, nil
	case CHECK:
		p.advance() // CHECK
		p.expect('(')
		expr, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		p.expect(')')
		noInherit := p.parseNoInherit()
		return &nodes.Constraint{
			Contype:        nodes.CONSTR_CHECK,
			RawExpr:        expr,
			IsNoInherit:    noInherit,
			Loc: nodes.NoLoc(),
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
			Contype:  nodes.CONSTR_DEFAULT,
			RawExpr:  expr,
			Loc: nodes.NoLoc(),
		}, nil
	case REFERENCES:
		p.advance()
		refNames, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && refNames == nil {
			return nil, p.syntaxErrorAtCur()
		}
		refRv := makeRangeVarFromNames(refNames)
		pkAttrs, _ := p.parseOptColumnList()
		matchType := p.parseKeyMatch()
		updAction, delAction, delSetCols := p.parseKeyActions()
		attrs := p.parseConstraintAttributeSpec()
		n := &nodes.Constraint{
			Contype:        nodes.CONSTR_FOREIGN,
			Pktable:        refRv,
			PkAttrs:        pkAttrs,
			FkMatchtype:    matchType,
			FkUpdaction:    updAction,
			FkDelaction:    delAction,
			FkDelsetcols:   delSetCols,
			Loc: nodes.NoLoc(),
			InitiallyValid: true,
		}
		applyConstraintAttrs(n, attrs)
		return n, nil
	case GENERATED:
		return p.parseGeneratedConstraint()
	default:
		return nil, nil
	}
}

// parseGeneratedConstraint parses GENERATED ... column constraints.
func (p *Parser) parseGeneratedConstraint() (*nodes.Constraint, error) {
	p.advance() // GENERATED

	if p.cur.Type == ALWAYS {
		p.advance() // ALWAYS
		p.expect(AS)
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
				Loc: nodes.NoLoc(),
			}, nil
		}
		// GENERATED ALWAYS AS '(' a_expr ')' STORED
		p.expect('(')
		expr, _ := p.parseAExpr(0)
		p.expect(')')
		p.expect(STORED)
		return &nodes.Constraint{
			Contype:       nodes.CONSTR_GENERATED,
			GeneratedWhen: 'a',
			RawExpr:       expr,
			Loc: nodes.NoLoc(),
		}, nil
	}

	// GENERATED BY DEFAULT AS ...
	p.expect(BY)
	p.expect(DEFAULT)
	p.expect(AS)
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
			Loc: nodes.NoLoc(),
		}, nil
	}
	// GENERATED BY DEFAULT AS '(' a_expr ')' STORED
	p.expect('(')
	expr, _ := p.parseAExpr(0)
	p.expect(')')
	p.expect(STORED)
	return &nodes.Constraint{
		Contype:       nodes.CONSTR_GENERATED,
		GeneratedWhen: 'd',
		RawExpr:       expr,
		Loc: nodes.NoLoc(),
	}, nil
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
	p.expect(')')
	return opts, nil
}

// makeDefElem creates a DefElem with Loc=NoLoc().
func makeDefElem(name string, arg nodes.Node) *nodes.DefElem {
	return &nodes.DefElem{
		Defname:  name,
		Arg:      arg,
		Loc: nodes.NoLoc(),
	}
}

// parseTableConstraint parses a TableConstraint.
//
//	TableConstraint:
//	    CONSTRAINT name ConstraintElem
//	    | ConstraintElem
func (p *Parser) parseTableConstraint() *nodes.Constraint {
	if p.cur.Type == CONSTRAINT {
		p.advance()
		name, _ := p.parseName()
		c := p.parseConstraintElem()
		if c != nil {
			c.Conname = name
		}
		return c
	}
	return p.parseConstraintElem()
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
func (p *Parser) parseConstraintElem() *nodes.Constraint {
	switch p.cur.Type {
	case UNIQUE:
		p.advance()
		nullsNotDistinct := p.parseOptUniqueNullTreatment()

		// Check for ExistingIndex: USING INDEX name
		if p.cur.Type == USING {
			next := p.peekNext()
			if next.Type == INDEX {
				p.advance() // USING
				p.advance() // INDEX
				if p.collectMode() {
					p.addRuleCandidate("qualified_name")
					return nil
				}
				idxName, _ := p.parseName()
				attrs := p.parseConstraintAttributeSpec()
				n := &nodes.Constraint{
					Contype:        nodes.CONSTR_UNIQUE,
					Indexname:      idxName,
					Loc: nodes.NoLoc(),
					InitiallyValid: true,
				}
				applyConstraintAttrs(n, attrs)
				return n
			}
		}

		p.expect('(')
		keys := p.parseColumnList()
		p.expect(')')
		including := p.parseOptCInclude()
		options := p.parseOptDefinition()
		indexspace := p.parseOptConsTableSpace()
		attrs := p.parseConstraintAttributeSpec()
		n := &nodes.Constraint{
			Contype:          nodes.CONSTR_UNIQUE,
			NullsNotDistinct: nullsNotDistinct,
			Keys:             keys,
			Including:        including,
			Options:          options,
			Indexspace:        indexspace,
			Loc: nodes.NoLoc(),
			InitiallyValid:   true,
		}
		applyConstraintAttrs(n, attrs)
		return n

	case PRIMARY:
		p.advance() // PRIMARY
		p.expect(KEY)

		// Check for ExistingIndex: USING INDEX name
		if p.cur.Type == USING {
			next := p.peekNext()
			if next.Type == INDEX {
				p.advance() // USING
				p.advance() // INDEX
				if p.collectMode() {
					p.addRuleCandidate("qualified_name")
					return nil
				}
				idxName, _ := p.parseName()
				attrs := p.parseConstraintAttributeSpec()
				n := &nodes.Constraint{
					Contype:        nodes.CONSTR_PRIMARY,
					Indexname:      idxName,
					Loc: nodes.NoLoc(),
					InitiallyValid: true,
				}
				applyConstraintAttrs(n, attrs)
				return n
			}
		}

		p.expect('(')
		keys := p.parseColumnList()
		p.expect(')')
		including := p.parseOptCInclude()
		options := p.parseOptDefinition()
		indexspace := p.parseOptConsTableSpace()
		attrs := p.parseConstraintAttributeSpec()
		n := &nodes.Constraint{
			Contype:        nodes.CONSTR_PRIMARY,
			Keys:           keys,
			Including:      including,
			Options:        options,
			Indexspace:      indexspace,
			Loc: nodes.NoLoc(),
			InitiallyValid: true,
		}
		applyConstraintAttrs(n, attrs)
		return n

	case CHECK:
		p.advance()
		p.expect('(')
		expr, _ := p.parseAExpr(0)
		p.expect(')')
		attrs := p.parseConstraintAttributeSpec()
		n := &nodes.Constraint{
			Contype:        nodes.CONSTR_CHECK,
			RawExpr:        expr,
			Loc: nodes.NoLoc(),
			InitiallyValid: true,
		}
		applyConstraintAttrs(n, attrs)
		return n

	case FOREIGN:
		p.advance() // FOREIGN
		p.expect(KEY)
		p.expect('(')
		fkAttrs := p.parseColumnList()
		p.expect(')')
		p.expect(REFERENCES)
		refNames, _ := p.parseQualifiedName()
		refRv := makeRangeVarFromNames(refNames)
		pkAttrs, _ := p.parseOptColumnList()
		matchType := p.parseKeyMatch()
		updAction, delAction, delSetCols := p.parseKeyActions()
		attrs := p.parseConstraintAttributeSpec()
		n := &nodes.Constraint{
			Contype:        nodes.CONSTR_FOREIGN,
			FkAttrs:        fkAttrs,
			Pktable:        refRv,
			PkAttrs:        pkAttrs,
			FkMatchtype:    matchType,
			FkUpdaction:    updAction,
			FkDelaction:    delAction,
			FkDelsetcols:   delSetCols,
			Loc: nodes.NoLoc(),
			InitiallyValid: true,
		}
		applyConstraintAttrs(n, attrs)
		return n

	case EXCLUDE:
		return p.parseExclusionConstraint()
	}
	return nil
}

// parseExclusionConstraint parses EXCLUDE constraint.
func (p *Parser) parseExclusionConstraint() *nodes.Constraint {
	p.advance() // EXCLUDE
	am := p.parseAccessMethodClause()
	p.expect('(')
	excl := p.parseExclusionConstraintList()
	p.expect(')')
	including := p.parseOptCInclude()
	options := p.parseOptDefinition()
	indexspace := p.parseOptConsTableSpace()
	where := p.parseWhereClauseOpt()
	attrs := p.parseConstraintAttributeSpec()
	n := &nodes.Constraint{
		Contype:        nodes.CONSTR_EXCLUSION,
		AccessMethod:   am,
		Exclusions:     excl,
		Including:      including,
		Options:        options,
		Indexspace:      indexspace,
		WhereClause:    where,
		Loc: nodes.NoLoc(),
		InitiallyValid: true,
	}
	applyConstraintAttrs(n, attrs)
	return n
}

// parseAccessMethodClause parses access_method_clause.
//
//	access_method_clause: USING name | /* EMPTY */
func (p *Parser) parseAccessMethodClause() string {
	if p.cur.Type == USING {
		p.advance()
		name, _ := p.parseName()
		return name
	}
	return ""
}

// parseExclusionConstraintList parses ExclusionConstraintList.
func (p *Parser) parseExclusionConstraintList() *nodes.List {
	elem := p.parseExclusionConstraintElem()
	items := []nodes.Node{elem}
	for p.cur.Type == ',' {
		p.advance()
		elem = p.parseExclusionConstraintElem()
		items = append(items, elem)
	}
	return &nodes.List{Items: items}
}

// parseExclusionConstraintElem parses a single exclusion constraint element.
//
//	ExclusionConstraintElem: index_elem WITH any_operator | index_elem WITH OPERATOR '(' any_operator ')'
func (p *Parser) parseExclusionConstraintElem() *nodes.List {
	idxElem := p.parseIndexElem()
	p.expect(WITH)
	var op *nodes.List
	if p.cur.Type == OPERATOR {
		p.advance()
		p.expect('(')
		op, _ = p.parseAnyOperator()
		p.expect(')')
	} else {
		op, _ = p.parseAnyOperator()
	}
	return &nodes.List{Items: []nodes.Node{idxElem, op}}
}

// parseOptCollate parses opt_collate.
//
//	opt_collate: COLLATE any_name | /* EMPTY */
func (p *Parser) parseOptCollate() *nodes.List {
	if p.cur.Type == COLLATE {
		p.advance()
		name, _ := p.parseAnyName()
		return name
	}
	return nil
}

// parseOptQualifiedName parses opt_qualified_name.
//
//	opt_qualified_name: any_name | /* EMPTY */
func (p *Parser) parseOptQualifiedName() *nodes.List {
	if p.isColId() || p.isTypeFunctionName() {
		name, _ := p.parseAnyName()
		return name
	}
	return nil
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
func (p *Parser) parseWhereClauseOpt() nodes.Node {
	if p.cur.Type == WHERE {
		p.advance()
		result, _ := p.parseAExpr(0)
		return result
	}
	return nil
}

// parseTableLikeClause parses a LIKE clause in CREATE TABLE.
//
//	TableLikeClause: LIKE qualified_name TableLikeOptionList?
func (p *Parser) parseTableLikeClause() (*nodes.TableLikeClause, error) {
	likeLoc := p.pos()
	p.advance() // LIKE
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	if !p.collectMode() && names == nil {
		return nil, p.syntaxErrorAtCur()
	}
	rv := makeRangeVarFromNames(names)
	var opts uint32
	for p.cur.Type == INCLUDING || p.cur.Type == EXCLUDING {
		isIncluding := p.cur.Type == INCLUDING
		p.advance()
		opt := p.parseTableLikeOption()
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
func (p *Parser) parseTableLikeOption() uint32 {
	switch p.cur.Type {
	case ALL:
		p.advance()
		return 0xFFFFFFFF
	case COMMENTS:
		p.advance()
		return 1
	case COMPRESSION:
		p.advance()
		return 2
	case CONSTRAINTS:
		p.advance()
		return 4
	case DEFAULTS:
		p.advance()
		return 8
	case GENERATED:
		p.advance()
		return 16
	case IDENTITY_P:
		p.advance()
		return 32
	case INDEXES:
		p.advance()
		return 64
	case STATISTICS:
		p.advance()
		return 128
	case STORAGE:
		p.advance()
		return 256
	}
	return 0
}

// parseTypedTableElementList parses a TypedTableElementList.
func (p *Parser) parseTypedTableElementList() *nodes.List {
	elem := p.parseTypedTableElement()
	if elem == nil {
		return nil
	}
	items := []nodes.Node{elem}
	for p.cur.Type == ',' {
		p.advance()
		elem = p.parseTypedTableElement()
		if elem == nil {
			break
		}
		items = append(items, elem)
	}
	return &nodes.List{Items: items}
}

// parseTypedTableElement parses a single TypedTableElement.
//
//	TypedTableElement: columnOptions | TableConstraint
func (p *Parser) parseTypedTableElement() nodes.Node {
	switch p.cur.Type {
	case CONSTRAINT, CHECK, UNIQUE, PRIMARY, FOREIGN, EXCLUDE:
		return p.parseTableConstraint()
	default:
		return p.parseColumnOptions()
	}
}

// parseColumnOptions parses columnOptions.
//
//	columnOptions: ColId opt_column_constraints | ColId WITH OPTIONS opt_column_constraints
func (p *Parser) parseColumnOptions() *nodes.ColumnDef {
	colname, err := p.parseColId()
	if err != nil {
		return nil
	}

	n := &nodes.ColumnDef{
		Colname:  colname,
		TypeName: nil,
		IsLocal:  true,
		Loc: nodes.NoLoc(),
	}

	if p.cur.Type == WITH {
		p.advance()
		p.expect(OPTIONS)
	}

	qualList, err := p.parseOptColumnConstraints()
	if err != nil {
		return nil
	}
	splitColQualList(qualList, n)

	return n
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
	list := p.parseColumnList()
	if list == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return list, nil
}

// parseColumnList parses a comma-separated list of column names.
//
//	columnList: name | columnList ',' name
func (p *Parser) parseColumnList() *nodes.List {
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil
	}
	name, err := p.parseName()
	if err != nil {
		return nil
	}
	items := []nodes.Node{&nodes.String{Str: name}}
	for p.cur.Type == ',' {
		p.advance()
		name, err = p.parseName()
		if err != nil {
			break
		}
		items = append(items, &nodes.String{Str: name})
	}
	return &nodes.List{Items: items}
}

// parseKeyMatch parses key_match.
//
//	key_match: MATCH FULL | MATCH PARTIAL | MATCH SIMPLE | /* EMPTY */
func (p *Parser) parseKeyMatch() byte {
	if p.cur.Type != MATCH {
		return 's' // default SIMPLE
	}
	p.advance()
	switch p.cur.Type {
	case FULL:
		p.advance()
		return 'f'
	case PARTIAL:
		p.advance()
		return 'p'
	default:
		// SIMPLE
		if p.cur.Type == IDENT && strings.EqualFold(p.cur.Str, "simple") {
			p.advance()
		}
		return 's'
	}
}

// parseKeyActions parses key_actions, returning (update action, delete action, delete set cols).
//
//	key_actions:
//	    key_update | key_delete | key_update key_delete | key_delete key_update | /* EMPTY */
func (p *Parser) parseKeyActions() (updAction byte, delAction byte, delSetCols *nodes.List) {
	updAction = 'a' // default NO ACTION
	delAction = 'a' // default NO ACTION

	for p.cur.Type == ON {
		next := p.peekNext()
		if next.Type == UPDATE {
			p.advance() // ON
			p.advance() // UPDATE
			updAction = p.parseKeyActionType(nil)
		} else if next.Type == DELETE_P {
			p.advance() // ON
			p.advance() // DELETE
			delAction = p.parseKeyActionType(&delSetCols)
		} else {
			break
		}
	}
	return
}

// parseKeyActionType parses a key_action value.
// If setCols is non-nil, it can receive the column list for SET NULL/SET DEFAULT.
func (p *Parser) parseKeyActionType(setCols **nodes.List) byte {
	switch p.cur.Type {
	case NO:
		p.advance()
		p.expect(ACTION)
		return 'a'
	case RESTRICT:
		p.advance()
		return 'r'
	case CASCADE:
		p.advance()
		return 'c'
	case SET:
		p.advance()
		if p.cur.Type == NULL_P {
			p.advance()
			if setCols != nil {
				*setCols, _ = p.parseOptColumnList()
			}
			return 'n'
		}
		// SET DEFAULT
		p.expect(DEFAULT)
		if setCols != nil {
			*setCols, _ = p.parseOptColumnList()
		}
		return 'd'
	default:
		return 'a'
	}
}

// parseConstraintAttributeSpec parses ConstraintAttributeSpec.
//
//	ConstraintAttributeSpec: /* EMPTY */ | ConstraintAttributeSpec ConstraintAttributeElem
func (p *Parser) parseConstraintAttributeSpec() int64 {
	var attrs int64
	for {
		switch p.cur.Type {
		case NOT:
			next := p.peekNext()
			switch next.Type {
			case DEFERRABLE:
				p.advance()
				p.advance()
				attrs |= int64(nodes.CAS_NOT_DEFERRABLE)
			case VALID:
				p.advance()
				p.advance()
				attrs |= int64(nodes.CAS_NOT_VALID)
			default:
				return attrs
			}
		case DEFERRABLE:
			p.advance()
			attrs |= int64(nodes.CAS_DEFERRABLE)
		case INITIALLY:
			p.advance()
			if p.cur.Type == IMMEDIATE {
				p.advance()
				attrs |= int64(nodes.CAS_INITIALLY_IMMEDIATE)
			} else if p.cur.Type == DEFERRED {
				p.advance()
				attrs |= int64(nodes.CAS_INITIALLY_DEFERRED)
			}
		case NO:
			next := p.peekNext()
			if next.Type == INHERIT {
				p.advance()
				p.advance()
				attrs |= int64(nodes.CAS_NO_INHERIT)
			} else {
				return attrs
			}
		default:
			return attrs
		}
	}
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
//	opt_unique_null_treatment: NULLS_LA DISTINCT | NULLS_LA NOT DISTINCT | /* EMPTY */
func (p *Parser) parseOptUniqueNullTreatment() bool {
	if p.cur.Type == NULLS_LA {
		p.advance()
		if p.cur.Type == NOT {
			p.advance()
			p.expect(DISTINCT)
			return true
		}
		p.expect(DISTINCT)
		return false
	}
	return false
}

// parseOptCInclude parses opt_c_include.
//
//	opt_c_include: INCLUDE '(' columnList ')' | /* EMPTY */
func (p *Parser) parseOptCInclude() *nodes.List {
	if p.cur.Type == INCLUDE {
		p.advance()
		p.expect('(')
		list := p.parseColumnList()
		p.expect(')')
		return list
	}
	return nil
}

// parseOptDefinition parses opt_definition.
//
//	opt_definition: WITH definition | /* EMPTY */
//	definition: '(' def_list ')'
func (p *Parser) parseOptDefinition() *nodes.List {
	if p.cur.Type != WITH {
		return nil
	}
	// Need lookahead to distinguish WITH '(' from WITH OIDS etc.
	next := p.peekNext()
	if next.Type != '(' {
		return nil
	}
	p.advance() // WITH
	p.advance() // (
	list, _ := p.parseDefList()
	p.expect(')')
	return list
}

// parseDefList parses def_list.
//
//	def_list: def_elem | def_list ',' def_elem
// parseDefList, parseDefElem, parseDefArg are in define.go

// parseOptConsTableSpace parses OptConsTableSpace.
//
//	OptConsTableSpace: USING INDEX TABLESPACE name | /* EMPTY */
func (p *Parser) parseOptConsTableSpace() string {
	if p.cur.Type == USING {
		next := p.peekNext()
		if next.Type == INDEX {
			p.advance() // USING
			p.advance() // INDEX
			p.expect(TABLESPACE)
			name, _ := p.parseName()
			return name
		}
	}
	return ""
}

// parseNoInherit parses no_inherit.
//
//	no_inherit: NO INHERIT | /* EMPTY */
func (p *Parser) parseNoInherit() bool {
	if p.cur.Type == NO {
		next := p.peekNext()
		if next.Type == INHERIT {
			p.advance()
			p.advance()
			return true
		}
	}
	return false
}

// parseOptInherit parses OptInherit.
//
//	OptInherit: INHERITS '(' qualified_name_list ')' | /* EMPTY */
func (p *Parser) parseOptInherit() *nodes.List {
	if p.cur.Type == INHERITS {
		p.advance()
		p.expect('(')
		list, _ := p.parseQualifiedNameList()
		p.expect(')')
		return list
	}
	return nil
}

// parseOptPartitionSpec parses OptPartitionSpec.
//
//	OptPartitionSpec: PartitionSpec | /* EMPTY */
func (p *Parser) parseOptPartitionSpec() *nodes.PartitionSpec {
	if p.cur.Type == PARTITION {
		return p.parsePartitionSpec()
	}
	return nil
}

// parsePartitionSpec parses PartitionSpec.
//
//	PartitionSpec: PARTITION BY ColId '(' part_params ')'
func (p *Parser) parsePartitionSpec() *nodes.PartitionSpec {
	p.advance() // PARTITION
	p.expect(BY)
	strategy, _ := p.parseColId()
	p.expect('(')
	params := p.parsePartParams()
	p.expect(')')
	return &nodes.PartitionSpec{
		Strategy:   parsePartitionStrategy(strategy),
		PartParams: params,
		Loc: nodes.NoLoc(),
	}
}

// parsePartitionStrategy converts partition strategy name to internal code.
func parsePartitionStrategy(strategy string) string {
	switch strings.ToLower(strategy) {
	case "list":
		return "l"
	case "range":
		return "r"
	case "hash":
		return "h"
	default:
		return strategy
	}
}

// parsePartParams parses part_params.
//
//	part_params: part_elem | part_params ',' part_elem
func (p *Parser) parsePartParams() *nodes.List {
	elem := p.parsePartElem()
	items := []nodes.Node{elem}
	for p.cur.Type == ',' {
		p.advance()
		elem = p.parsePartElem()
		items = append(items, elem)
	}
	return &nodes.List{Items: items}
}

// parsePartElem parses a single partition key element.
//
//	part_elem:
//	    ColId opt_collate opt_qualified_name
//	    | func_expr_windowless opt_collate opt_qualified_name
//	    | '(' a_expr ')' opt_collate opt_qualified_name
func (p *Parser) parsePartElem() *nodes.PartitionElem {
	if p.cur.Type == '(' {
		p.advance()
		expr, _ := p.parseAExpr(0)
		p.expect(')')
		collation := p.parseOptCollate()
		opclass := p.parseOptQualifiedName()
		return &nodes.PartitionElem{
			Expr:      expr,
			Collation: collation,
			Opclass:   opclass,
			Loc: nodes.NoLoc(),
		}
	}

	// Try as ColId first (checking if next is not '(' which would make it a function)
	if p.isColId() {
		next := p.peekNext()
		if next.Type != '(' {
			name, _ := p.parseColId()
			collation := p.parseOptCollate()
			opclass := p.parseOptQualifiedName()
			return &nodes.PartitionElem{
				Name:      name,
				Collation: collation,
				Opclass:   opclass,
				Loc: nodes.NoLoc(),
			}
		}
	}

	// func_expr_windowless
	expr, _ := p.parseFuncExprWindowless()
	collation := p.parseOptCollate()
	opclass := p.parseOptQualifiedName()
	return &nodes.PartitionElem{
		Expr:      expr,
		Collation: collation,
		Opclass:   opclass,
		Loc: nodes.NoLoc(),
	}
}

// parseOptAccessMethod parses OptAccessMethod.
//
//	OptAccessMethod: USING name | /* EMPTY */
func (p *Parser) parseOptAccessMethod() string {
	if p.cur.Type == USING {
		p.advance()
		name, _ := p.parseName()
		return name
	}
	return ""
}

// parseOptWith parses OptWith for CREATE TABLE.
//
//	OptWith: WITH reloptions | WITH OIDS | WITHOUT OIDS | /* EMPTY */
func (p *Parser) parseOptWith() *nodes.List {
	if p.cur.Type == WITH {
		next := p.peekNext()
		if next.Type == '(' {
			p.advance() // WITH
			return p.parseReloptions()
		}
		if next.Type == OIDS {
			p.advance() // WITH
			p.advance() // OIDS
			return nil  // deprecated
		}
		return nil
	}
	if p.cur.Type == WITHOUT {
		next := p.peekNext()
		if next.Type == OIDS {
			p.advance() // WITHOUT
			p.advance() // OIDS
			return nil
		}
	}
	return nil
}

// parseReloptions parses reloptions.
//
//	reloptions: '(' reloption_list ')'
func (p *Parser) parseReloptions() *nodes.List {
	p.expect('(')
	list := p.parseReloptionList()
	p.expect(')')
	return list
}

// parseReloptionList parses reloption_list.
func (p *Parser) parseReloptionList() *nodes.List {
	elem := p.parseReloptionElem()
	items := []nodes.Node{elem}
	for p.cur.Type == ',' {
		p.advance()
		elem = p.parseReloptionElem()
		items = append(items, elem)
	}
	return &nodes.List{Items: items}
}

// parseReloptionElem parses a single reloption_elem.
//
//	reloption_elem:
//	    ColLabel '=' def_arg
//	    | ColLabel
//	    | ColLabel '.' ColLabel '=' def_arg
//	    | ColLabel '.' ColLabel
func (p *Parser) parseReloptionElem() *nodes.DefElem {
	label, _ := p.parseColLabel()

	if p.cur.Type == '.' {
		p.advance()
		sublabel, _ := p.parseColLabel()
		if p.cur.Type == '=' {
			p.advance()
			arg, _ := p.parseDefArg()
			return &nodes.DefElem{
				Defnamespace: label,
				Defname:      sublabel,
				Arg:          arg,
				Loc: nodes.NoLoc(),
			}
		}
		return &nodes.DefElem{
			Defnamespace: label,
			Defname:      sublabel,
			Loc: nodes.NoLoc(),
		}
	}

	if p.cur.Type == '=' {
		p.advance()
		arg, _ := p.parseDefArg()
		return &nodes.DefElem{
			Defname:  label,
			Arg:      arg,
			Loc: nodes.NoLoc(),
		}
	}

	return &nodes.DefElem{
		Defname:  label,
		Loc: nodes.NoLoc(),
	}
}

// parseOnCommitOption parses OnCommitOption.
//
//	OnCommitOption: ON COMMIT DROP | ON COMMIT DELETE ROWS | ON COMMIT PRESERVE ROWS | /* EMPTY */
func (p *Parser) parseOnCommitOption() nodes.OnCommitAction {
	if p.cur.Type != ON {
		return nodes.ONCOMMIT_NOOP
	}
	next := p.peekNext()
	if next.Type != COMMIT {
		return nodes.ONCOMMIT_NOOP
	}
	p.advance() // ON
	p.advance() // COMMIT
	switch p.cur.Type {
	case DROP:
		p.advance()
		return nodes.ONCOMMIT_DROP
	case DELETE_P:
		p.advance()
		p.expect(ROWS)
		return nodes.ONCOMMIT_DELETE_ROWS
	case PRESERVE:
		p.advance()
		p.expect(ROWS)
		return nodes.ONCOMMIT_PRESERVE_ROWS
	}
	return nodes.ONCOMMIT_NOOP
}

// parseOptTableSpace parses OptTableSpace.
//
//	OptTableSpace: TABLESPACE name | /* EMPTY */
func (p *Parser) parseOptTableSpace() string {
	if p.cur.Type == TABLESPACE {
		p.advance()
		name, _ := p.parseName()
		return name
	}
	return ""
}

// parseForValues parses ForValues (partition bound spec).
//
//	ForValues: PartitionBoundSpec | DEFAULT
func (p *Parser) parseForValues() nodes.Node {
	if p.cur.Type == DEFAULT {
		p.advance()
		return &nodes.PartitionBoundSpec{
			IsDefault: true,
			Loc: nodes.NoLoc(),
		}
	}
	return p.parsePartitionBoundSpec()
}

// parsePartitionBoundSpec parses PartitionBoundSpec.
//
//	PartitionBoundSpec:
//	    FOR VALUES IN '(' expr_list ')'
//	    | FOR VALUES FROM '(' expr_list ')' TO '(' expr_list ')'
//	    | FOR VALUES WITH '(' hash_partbound ')'
func (p *Parser) parsePartitionBoundSpec() *nodes.PartitionBoundSpec {
	p.expect(FOR)
	p.expect(VALUES)

	switch p.cur.Type {
	case IN_P:
		p.advance()
		p.expect('(')
		list := p.parseAExprList()
		p.expect(')')
		return &nodes.PartitionBoundSpec{
			Strategy:   'l',
			Listdatums: list,
			Loc: nodes.NoLoc(),
		}
	case FROM:
		p.advance()
		p.expect('(')
		lower := p.parseAExprList()
		p.expect(')')
		p.expect(TO)
		p.expect('(')
		upper := p.parseAExprList()
		p.expect(')')
		return &nodes.PartitionBoundSpec{
			Strategy:    'r',
			Lowerdatums: lower,
			Upperdatums: upper,
			Loc: nodes.NoLoc(),
		}
	case WITH:
		p.advance()
		p.expect('(')
		bounds := p.parseHashPartbound()
		p.expect(')')
		spec := &nodes.PartitionBoundSpec{
			Strategy:  'h',
			Modulus:   -1,
			Remainder: -1,
			Loc: nodes.NoLoc(),
		}
		for _, item := range bounds.Items {
			opt := item.(*nodes.DefElem)
			switch opt.Defname {
			case "modulus":
				spec.Modulus = int(opt.Arg.(*nodes.Integer).Ival)
			case "remainder":
				spec.Remainder = int(opt.Arg.(*nodes.Integer).Ival)
			}
		}
		return spec
	}
	return nil
}

// parseAExprList parses a comma-separated list of a_expr.
func (p *Parser) parseAExprList() *nodes.List {
	expr, _ := p.parseAExpr(0)
	items := []nodes.Node{expr}
	for p.cur.Type == ',' {
		p.advance()
		expr, _ = p.parseAExpr(0)
		items = append(items, expr)
	}
	return &nodes.List{Items: items}
}

// parseHashPartbound parses hash_partbound.
//
//	hash_partbound: hash_partbound_elem | hash_partbound ',' hash_partbound_elem
func (p *Parser) parseHashPartbound() *nodes.List {
	elem := p.parseHashPartboundElem()
	items := []nodes.Node{elem}
	for p.cur.Type == ',' {
		p.advance()
		elem = p.parseHashPartboundElem()
		items = append(items, elem)
	}
	return &nodes.List{Items: items}
}

// parseHashPartboundElem parses a single hash_partbound_elem.
//
//	hash_partbound_elem: NonReservedWord Iconst
func (p *Parser) parseHashPartboundElem() *nodes.DefElem {
	// NonReservedWord is similar to ColId
	name, _ := p.parseColId()
	val := p.cur.Ival
	p.advance() // ICONST
	return &nodes.DefElem{
		Defname:  name,
		Arg:      &nodes.Integer{Ival: val},
		Loc: nodes.NoLoc(),
	}
}
