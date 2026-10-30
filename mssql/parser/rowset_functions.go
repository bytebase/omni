// Package parser - rowset_functions.go implements T-SQL rowset function parsing.
// Rowset functions: OPENROWSET, OPENQUERY, OPENJSON, OPENDATASOURCE, OPENXML,
// and PREDICT. These appear as table sources in FROM clauses.
package parser

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseRowsetFunction parses OPENROWSET, OPENQUERY, OPENJSON, OPENDATASOURCE, OPENXML
// as table sources. These are treated as FuncCallExpr nodes with keyword names.
func (p *Parser) parseRowsetFunction() (nodes.TableExpr, error) {
	loc := p.pos()
	funcTok := p.cur.Type
	funcName := p.cur.Str
	p.advance() // consume the keyword
	nameEnd := p.prevEnd()

	if p.cur.Type == '(' {
		fc := &nodes.FuncCallExpr{
			Name: &nodes.TableRef{Object: funcName, Loc: nodes.Loc{Start: loc, End: nameEnd}},
			Loc:  nodes.Loc{Start: loc, End: -1},
		}
		p.advance() // consume (

		argIndex := 0
		args, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
			if isRowsetObjectArg(funcName, argIndex) && p.isIdentLike() {
				ref, err := p.parseTableRef()
				if err != nil {
					return nil, err
				}
				if ref == nil {
					return nil, p.unexpectedToken()
				}
				argIndex++
				return ref, nil
			}
			arg, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if arg == nil {
				return nil, p.unexpectedToken()
			}
			argIndex++
			return arg, nil
		})
		if err != nil {
			return nil, err
		}
		if len(args) > 0 {
			fc.Args = &nodes.List{Items: args}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		fc.Loc.End = p.prevEnd()

		if funcTok == kwOPENDATASOURCE {
			return p.parseAdHocTableRef(fc)
		}

		// WITH schema: OPENJSON and OPENXML only. SQL Server rejects WITH
		// after OPENQUERY and OPENROWSET. OPENXML also takes the schema of a
		// table: OPENXML(@h, '/r') WITH dbo.t.
		var withClause *nodes.List
		var withTable *nodes.TableRef
		isJSON := strings.EqualFold(funcName, "OPENJSON")
		isXML := strings.EqualFold(funcName, "OPENXML")
		if (isJSON || isXML) && p.cur.Type == kwWITH && p.peekNext().Type == '(' {
			p.advance() // consume WITH
			p.advance() // consume (
			withClause, err = p.parseRowsetWithClause(isJSON)
			if err != nil {
				return nil, err
			}
		} else if isXML && p.cur.Type == kwWITH {
			if next := p.peekNext().Type; next == tokIDENT || isContextKeyword(next) {
				p.advance() // consume WITH
				if withTable, err = p.parseOpenXMLWithTable(); err != nil {
					return nil, err
				}
			}
		}

		alias := p.parseOptionalAlias()
		aliasEnd := p.prevEnd()

		return &nodes.AliasedTableRef{
			Table:     fc,
			Alias:     alias,
			Columns:   withClause,
			WithTable: withTable,
			Loc:       nodes.Loc{Start: loc, End: aliasEnd},
		}, nil
	}

	// Shouldn't happen but handle gracefully
	return &nodes.AliasedTableRef{
		Table: &nodes.FuncCallExpr{
			Name: &nodes.TableRef{Object: funcName, Loc: nodes.Loc{Start: loc, End: nameEnd}},
			Loc:  nodes.Loc{Start: loc, End: nameEnd},
		},
		Loc: nodes.Loc{Start: loc, End: nameEnd},
	}, nil
}

// parseOpenXMLWithTable parses the table after OPENXML ... WITH whose columns
// OPENXML takes. SQL Server takes a name of one to three parts that ends with
// the table, and rejects a remote table: "OPENXML does not support retrieving
// schema from remote tables".
func (p *Parser) parseOpenXMLWithTable() (*nodes.TableRef, error) {
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if ref == nil || ref.Object == "" {
		// parseTableRef leaves Object empty for a name that ends with a dot,
		// as in WITH dbo.
		return nil, p.unexpectedToken()
	}
	if ref.Server != "" {
		return nil, p.newParseError(ref.Loc.Start, "OPENXML does not take the schema of a remote table")
	}
	return ref, nil
}

// parseAdHocTableRef parses the object name and alias after
// OPENDATASOURCE ( provider_name , init_string ), whose call is fc.
//
//	OPENDATASOURCE ( provider_name , init_string ) . database . schema . object [ [ AS ] alias ]
//
// SQL Server takes two string literals as the arguments and exactly three
// name parts, of which database and schema may be empty (db..t, ...t). The
// call alone, a shorter or longer name, a column list, a table hint, and
// TABLESAMPLE are errors.
func (p *Parser) parseAdHocTableRef(fc *nodes.FuncCallExpr) (*nodes.AdHocTableRef, error) {
	var args []nodes.Node
	if fc.Args != nil {
		args = fc.Args.Items
	}
	for i, arg := range args {
		if lit, ok := arg.(*nodes.Literal); i >= 2 || !ok || lit.Type != nodes.LitString {
			return nil, p.newParseError(nodes.NodeLoc(arg).Start, "OPENDATASOURCE takes a provider name and an init string")
		}
	}
	if len(args) < 2 {
		return nil, p.newParseError(fc.Loc.End-1, "OPENDATASOURCE takes a provider name and an init string")
	}

	if _, err := p.expect('.'); err != nil {
		return nil, err
	}
	objLoc := p.pos()
	var parts [3]string
	for i := range parts {
		if i > 0 {
			if _, err := p.expect('.'); err != nil {
				return nil, err
			}
		}
		if name, ok := p.parseIdentifier(); ok {
			parts[i] = name
		}
	}
	if parts[2] == "" || p.cur.Type == '.' {
		return nil, p.unexpectedToken()
	}

	ref := &nodes.AdHocTableRef{
		DataSource: fc,
		Object: &nodes.TableRef{
			Database: parts[0],
			Schema:   parts[1],
			Object:   parts[2],
			Loc:      nodes.Loc{Start: objLoc, End: p.prevEnd()},
		},
		Loc: nodes.Loc{Start: fc.Loc.Start, End: -1},
	}
	ref.Alias = p.parseOptionalAlias()
	ref.Loc.End = p.prevEnd()
	return ref, nil
}

func isRowsetObjectArg(funcName string, argIndex int) bool {
	switch {
	case strings.EqualFold(funcName, "OPENQUERY"):
		return argIndex == 0
	case strings.EqualFold(funcName, "OPENROWSET"):
		return argIndex == 2
	default:
		return false
	}
}

// parseRowsetWithClause parses the WITH schema of OPENJSON (isJSON) or
// OPENXML. The opening '(' has already been consumed.
//
//	column_name data_type [ COLLATE collation_name ] [ column_path ] [ AS JSON ]
//
// SQL Server requires the data type, and takes the collation only
// undelimited, the path only as a string literal, and AS JSON only in
// OPENJSON, in this order.
func (p *Parser) parseRowsetWithClause(isJSON bool) (*nodes.List, error) {
	cols, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
		loc := p.pos()
		colName, ok := p.parseIdentifier()
		if !ok {
			return nil, p.unexpectedToken()
		}
		if p.collectMode() {
			p.addRuleCandidate("type_name")
			return nil, errCollecting
		}
		dt, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		if dt == nil {
			return nil, p.unexpectedToken()
		}
		col := &nodes.RowsetColumnDef{
			Name:     colName,
			DataType: dt,
			Loc:      nodes.Loc{Start: loc, End: -1},
		}
		if _, ok := p.match(kwCOLLATE); ok {
			collation, err := p.parseCollationName()
			if err != nil {
				return nil, err
			}
			col.Collation = collation
		}
		if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
			tok := p.advance()
			col.Path = &nodes.Literal{
				Type:    nodes.LitString,
				Str:     tok.Str,
				IsNChar: tok.Type == tokNSCONST,
				Loc:     nodes.Loc{Start: tok.Loc, End: p.prevEnd()},
			}
		}
		if isJSON && p.cur.Type == kwAS && p.peekNext().Type == kwJSON {
			p.advance() // consume AS
			p.advance() // consume JSON
			col.AsJSON = true
		}
		col.Loc.End = p.prevEnd()
		return col, nil
	})
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: cols}, nil
}

// predictAhead reports whether the table source at the current token is
// PREDICT. PREDICT is not reserved: FROM predict, FROM predict(1), and FROM
// dbo.predict(1) name a table or a function, and SQL Server takes [PREDICT]
// ( ... ) WITH ( ... ) as PREDICT too. A table-valued function takes no
// name = value argument, so PREDICT ( followed by MODEL, DATA, or RUNTIME
// and = is PREDICT.
func (p *Parser) predictAhead() bool {
	return p.isValidOption(predictName) && p.peekNext().Type == '(' &&
		p.isValidOptionToken(predictArguments, p.peekSecond()) && p.peekThird().Type == '='
}

var (
	predictName      = newOptionSet(kwPREDICT)
	predictArguments = newOptionSet(kwMODEL, kwDATA).withIdents("RUNTIME")
	predictRuntimes  = newOptionSet().withIdents("ONNX")
)

// parsePredictTableRef parses PREDICT as a table source. The current token
// is PREDICT (predictAhead).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/predict-transact-sql
//
//	PREDICT ( MODEL = model , DATA = table_source [ [ AS ] table_alias ]
//	          [ , RUNTIME = ONNX ] )
//	WITH ( { column_name data_type [ COLLATE collation_name ] [ NULL | NOT NULL ] } [ ,...n ] )
//	[ [ AS ] table_alias ]
//
// PREDICT is a table source only: SQL Server rejects it as a statement. It
// takes the arguments in any order, MODEL and DATA once each (Msg 39036 when
// one is missing, Msg 39035 when one repeats). RUNTIME = ONNX is Azure
// Synapse Analytics syntax, which SQL Server 2022 rejects; omni keeps it, as
// it keeps the other Synapse statements. MODEL is a scalar expression, a
// scalar subquery included (Synapse). The WITH list is required and the
// alias after it optional; SQL Server rejects a column alias list, table
// hints, and TABLESAMPLE after the alias.
func (p *Parser) parsePredictTableRef() (nodes.TableExpr, error) {
	loc := p.pos()
	p.advance() // PREDICT
	p.advance() // (
	ref := &nodes.PredictTableRef{Loc: nodes.Loc{Start: loc, End: -1}}
	seen := map[string]bool{}
	for {
		if p.collectMode() {
			p.addTokenCandidate(kwMODEL)
			p.addTokenCandidate(kwDATA)
			return nil, errCollecting
		}
		if !p.isIdentLike() {
			return nil, p.unexpectedToken()
		}
		argLoc := p.pos()
		name, err := p.expectOption(predictArguments)
		if err != nil {
			return nil, err
		}
		if seen[name] {
			return nil, p.newParseError(argLoc, fmt.Sprintf("PREDICT argument %s is specified more than once", name))
		}
		seen[name] = true
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		switch name {
		case "MODEL":
			model, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if model == nil {
				return nil, p.unexpectedToken()
			}
			ref.Model = model
		case "DATA":
			data, err := p.parsePredictData()
			if err != nil {
				return nil, err
			}
			ref.Data = data
		case "RUNTIME":
			if !p.isIdentLike() || !p.isValidOption(predictRuntimes) {
				return nil, p.unexpectedToken()
			}
			p.advance()
			ref.Runtime = "ONNX"
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if ref.Model == nil || ref.Data == nil {
		missing := "MODEL"
		if ref.Model != nil {
			missing = "DATA"
		}
		return nil, p.newParseError(p.pos(), "PREDICT expects argument "+missing)
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	cols, err := p.parseCommaList(')', commaListStrict, p.parsePredictColumn)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	ref.Columns = &nodes.List{Items: cols}

	ref.Alias = p.parseOptionalAlias()
	ref.Loc.End = p.prevEnd()
	return ref, nil
}

// parsePredictData parses the table source of PREDICT's DATA argument, with
// its optional alias: a table, view, or common table expression, a
// schema-qualified table-valued function, a table variable, or a derived
// table.
//
// The documentation names a table, a table alias, a CTE alias, a view, and a
// table-valued function, and SQL Server also accepts a derived table. It
// rejects a column alias list, table hints, FOR SYSTEM_TIME, TABLESAMPLE, a
// join, VALUES, and the rowset functions there, and an unqualified function
// name with Msg 195, as no built-in function has that name. SQL Server's
// parser reads DATA as an expression, so it also lets DATA = 1 through, which
// is no table source; omni does not.
func (p *Parser) parsePredictData() (nodes.TableExpr, error) {
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, errCollecting
	}
	loc := p.pos()
	switch p.cur.Type {
	case '(':
		if p.peekNext().Type != kwSELECT {
			return nil, p.unexpectedToken()
		}
		p.advance() // (
		query, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		if query == nil {
			return nil, p.unexpectedToken()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		sub := &nodes.SubqueryExpr{Query: query, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}
		alias := p.parseOptionalAlias()
		if alias == "" {
			return sub, nil
		}
		return &nodes.AliasedTableRef{Table: sub, Alias: alias, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
	case tokVARIABLE:
		tok := p.advance()
		tv := &nodes.TableVarRef{Name: tok.Str, Loc: nodes.Loc{Start: loc, End: -1}}
		tv.Alias = p.parseOptionalAlias()
		tv.Loc.End = p.prevEnd()
		return tv, nil
	}
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if ref == nil {
		return nil, p.unexpectedToken()
	}
	if p.cur.Type != '(' {
		ref.Alias = p.parseOptionalAlias()
		ref.Loc.End = p.prevEnd()
		return ref, nil
	}
	if ref.Schema == "" {
		return nil, p.newParseError(loc, fmt.Sprintf("'%s' is not a recognized built-in function name", ref.Object))
	}
	p.advance() // (
	fc := &nodes.FuncCallExpr{Name: ref, Loc: nodes.Loc{Start: loc, End: -1}}
	args, err := p.parseCommaList(')', commaListAllowEmpty, func() (nodes.Node, error) {
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if arg == nil {
			return nil, p.unexpectedToken()
		}
		return arg, nil
	})
	if err != nil {
		return nil, err
	}
	if len(args) > 0 {
		fc.Args = &nodes.List{Items: args}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	fc.Loc.End = p.prevEnd()
	alias := p.parseOptionalAlias()
	return &nodes.AliasedTableRef{Table: fc, Alias: alias, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
}

// parsePredictColumn parses a column of PREDICT's WITH list.
//
//	column_name data_type [ COLLATE collation_name ] [ NULL | NOT NULL ]
//
// SQL Server's parser also takes some storage options of a table column
// there (SPARSE, FILESTREAM, MASKED WITH, GENERATED ALWAYS ... HIDDEN), which
// mean nothing for the columns PREDICT returns; omni takes the documented
// form.
func (p *Parser) parsePredictColumn() (nodes.Node, error) {
	loc := p.pos()
	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.unexpectedToken()
	}
	if p.collectMode() {
		p.addRuleCandidate("type_name")
		return nil, errCollecting
	}
	dt, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	if dt == nil {
		return nil, p.unexpectedToken()
	}
	col := &nodes.RowsetColumnDef{
		Name:     name,
		DataType: dt,
		Loc:      nodes.Loc{Start: loc, End: -1},
	}
	if _, ok := p.match(kwCOLLATE); ok {
		collation, err := p.parseCollationName()
		if err != nil {
			return nil, err
		}
		col.Collation = collation
	}
	switch {
	case p.cur.Type == kwNULL:
		nullLoc := p.pos()
		p.advance()
		col.Nullable = &nodes.NullableSpec{Loc: nodes.Loc{Start: nullLoc, End: p.prevEnd()}}
	case p.cur.Type == kwNOT && p.peekNext().Type == kwNULL:
		nullLoc := p.pos()
		p.advance() // NOT
		p.advance() // NULL
		col.Nullable = &nodes.NullableSpec{NotNull: true, Loc: nodes.Loc{Start: nullLoc, End: p.prevEnd()}}
	}
	col.Loc.End = p.prevEnd()
	return col, nil
}
