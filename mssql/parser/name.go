// Package parser - name.go implements identifier and qualified name parsing.
package parser

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// isIdentLike returns true if the current token can be used as an identifier.
// In T-SQL, context keywords can be used as identifiers but Core keywords cannot
// (unless bracket-quoted). Use isKeywordOrIdent() for positions that accept all keywords.
// The tokIDENT type already includes [bracketed] and "quoted" identifiers.
func (p *Parser) isIdentLike() bool {
	if p.cur.Type == tokIDENT {
		return true
	}
	return isContextKeyword(p.cur.Type)
}

// isKeywordOrIdent returns true for ANY token that carries a name payload —
// tokIDENT or any registered keyword (including CoreKeyword). This is the
// correct predicate for true E-class positions where the grammar genuinely
// accepts reserved keywords as values or subcommand markers:
//   - multi-word permission names in GRANT (ALTER ANY DATABASE, CREATE TABLE, …)
//   - securable class prefix before :: (DATABASE::, SCHEMA::, OBJECT::, …)
//   - data type names (INT, VARCHAR, DATETIME are all keywords)
//   - audit action words (SELECT, INSERT, UPDATE, DELETE are CoreKeywords but
//     valid audit action names)
//
// For identifier / name / alias positions use isIdentLike, which correctly
// rejects CoreKeyword tokens unless bracketed.
func (p *Parser) isKeywordOrIdent() bool {
	if p.cur.Type == tokIDENT {
		return true
	}
	return p.cur.Type >= kwABSENT && p.cur.Str != ""
}

// parseIdentifier consumes and returns the current token as an identifier string.
// It accepts tokIDENT tokens and keywords used as identifiers.
// Returns ("", false) if the current token is not identifier-like.
//
//	identifier = regular_identifier | bracketed_identifier | quoted_identifier | keyword_as_identifier
func (p *Parser) parseIdentifier() (string, bool) {
	if p.cur.Type == tokIDENT {
		name := p.cur.Str
		p.advance()
		return name, true
	}
	if isContextKeyword(p.cur.Type) {
		name := p.cur.Str
		p.advance()
		return name, true
	}
	return "", false
}

// parseTableRef parses a qualified object name: [server.][database.][schema.]object
// Used for table names in DDL/DML contexts (FROM, CREATE TABLE, INSERT INTO, etc.).
//
// Ref: https://learn.microsoft.com/en-us/sql/relational-databases/databases/database-identifiers
//
//	qualified_name = [ server_name . [ database_name ] . [ schema_name ] . ]
//	                 | [ database_name . [ schema_name ] . ]
//	                 | [ schema_name . ]
//	                 object_name
func (p *Parser) parseTableRef() (*nodes.TableRef, error) {
	return p.parseObjectRef("table_ref")
}

// parseObjectRef parses a qualified object name and emits completionRule after
// a dot in that qualified name. Use parseTableRef for normal table contexts.
func (p *Parser) parseObjectRef(completionRule string) (*nodes.TableRef, error) {
	loc := p.pos()

	// A name may start with dots, which name empty parts (identifierListElement,
	// TSql160.g:30332): SQL Server reads FROM .a, FROM ..a and FROM .a() (measured).
	var leading []string
	for p.cur.Type == '.' {
		// A name holds at most four parts, and a fifth is Msg 117: "The object
		// name '....a' contains more than the maximum number of prefixes"
		// (measured).
		if len(leading) == 3 {
			return nil, p.newParseError(p.cur.Loc, "an object name has at most four parts")
		}
		leading = append(leading, "")
		p.advance() // consume .
		if p.collectMode() {
			p.addRuleCandidate(completionRule)
			return nil, errCollecting
		}
	}

	name, ok := p.parseIdentifier()
	if !ok {
		if len(leading) > 0 {
			return nil, p.syntaxErrorAtCur()
		}
		return nil, nil
	}
	if p.collectMode() && p.cursorOff <= p.prev.End {
		p.addRuleCandidate("table_ref")
		return nil, errCollecting
	}

	ref := &nodes.TableRef{
		Object: name,
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	// Collect dot-separated parts
	parts := append(leading, name)
	for p.cur.Type == '.' {
		// Four parts are the most, and a fifth is Msg 117: "The object name
		// '...a.b' contains more than the maximum number of prefixes".
		if len(parts) == 4 {
			return nil, p.newParseError(p.cur.Loc, "an object name has at most four parts")
		}
		p.advance() // consume .
		// Completion: after dot in qualified name → caller-specific object rule.
		if p.collectMode() {
			p.addRuleCandidate(completionRule)
			return nil, errCollecting
		}
		part, ok := p.parseIdentifier()
		if !ok {
			// Handle trailing dot (e.g., "db..object" means db.dbo.object with empty schema)
			parts = append(parts, "")
			continue
		}
		parts = append(parts, part)
	}

	// Assign parts based on count: object, schema.object, db.schema.object, server.db.schema.object
	assignObjectRefParts(ref, parts)

	ref.Loc.End = p.prevEnd()
	return ref, nil
}

// assignObjectRefParts assigns the dot-separated parts of a name to a
// TableRef: object, schema.object, db.schema.object, server.db.schema.object.
func assignObjectRefParts(ref *nodes.TableRef, parts []string) {
	switch len(parts) {
	case 1:
		ref.Object = parts[0]
	case 2:
		ref.Schema = parts[0]
		ref.Object = parts[1]
	case 3:
		ref.Database = parts[0]
		ref.Schema = parts[1]
		ref.Object = parts[2]
	default: // 4+
		ref.Server = parts[0]
		ref.Database = parts[1]
		ref.Schema = parts[2]
		ref.Object = parts[3]
	}
}

// parseMultiPartIdentifier parses identifierList (TSql160.g:30325) into its
// parts, keeping the empty parts of a name such as `.a`, `..a` or `a..b`:
//
//	identifierList = identifier ( identifierListElement )*
//	               | identifierListElement ( identifierListElement )*
//	identifierListElement = Dot Dot* identifier
//
// Each dot after the first joins the part to the previous one; the rest name
// an empty part, so `a..b` is [a, "", b] and `..a` is ["", "", a].
//
// SQL Server reads at most four parts, and refuses the fifth: `a.b.c.d.e` is
// "The procedure name 'a.b.c.d.e' contains more than the maximum number of
// prefixes. The maximum is 3." — for an EXEC, for a call by name and for
// ALTER AUTHORIZATION's securable alike (measured on the pinned server).
//
// completionRule is what completion offers after a dot, as parseObjectRef
// does; "" offers nothing.
const maxNameParts = 4

func (p *Parser) parseMultiPartIdentifier(completionRule string) ([]string, error) {
	var parts []string
	for i := 0; ; i++ {
		dots := 0
		for p.cur.Type == '.' {
			dots++
			p.advance()
		}
		if i > 0 && dots == 0 {
			return parts, nil
		}
		if completionRule != "" && p.collectMode() {
			p.addRuleCandidate(completionRule)
			return nil, errCollecting
		}
		empties := dots
		if i > 0 {
			empties-- // the dot that joins this part to the previous one
		}
		if len(parts)+empties+1 > maxNameParts {
			return nil, p.newParseError(p.cur.Loc, "a name of more than four parts")
		}
		for j := 0; j < empties; j++ {
			parts = append(parts, "")
		}
		id, ok := p.parseIdentifier()
		if !ok {
			return nil, p.syntaxErrorAtCur()
		}
		parts = append(parts, id)
		if completionRule != "" && i == 0 && p.collectMode() && p.cursorOff <= p.prev.End {
			// The cursor ends the first part: a qualified name may follow.
			p.addRuleCandidate(completionRule)
			return nil, errCollecting
		}
	}
}

// parseSequenceRef parses a qualified sequence name for NEXT VALUE FOR.
// Shape is the same as other schema object names, but completion must surface
// sequence_ref rather than the generic table_ref rule.
func (p *Parser) parseSequenceRef() (*nodes.TableRef, error) {
	return p.parseObjectRef("sequence_ref")
}

// parseVariableTableSource parses a T-SQL table variable used as a table source
// in FROM / JOIN positions. Handles three shapes:
//
//	@t                            -> *nodes.TableVarRef
//	@t [AS] alias                 -> *nodes.TableVarRef with Alias
//	@v.Method(args) [alias [(c)]] -> *nodes.TableVarMethodCallRef
//
// Mirrors SqlScriptDOM variableTableReference + variableMethodCallTableReference.
// Caller must ensure p.cur.Type == tokVARIABLE.
func (p *Parser) parseVariableTableSource() (nodes.TableExpr, error) {
	loc := p.pos()
	name := p.cur.Str // includes leading '@'
	p.advance()

	// @v.Method(args) — variable method-call table reference (e.g. XML .nodes()).
	// Any '.' after @name is necessarily a method call (table variables cannot be
	// schema-qualified), so we can commit to this branch on sight of the dot.
	if p.cur.Type == '.' {
		p.advance() // consume .
		method, ok := p.parseIdentifier()
		if !ok {
			return nil, p.unexpectedToken()
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		var args []nodes.ExprNode
		if p.cur.Type != ')' {
			for {
				arg, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				args = append(args, arg)
				if _, ok := p.match(','); !ok {
					break
				}
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		mc := &nodes.TableVarMethodCallRef{
			Var:    name,
			Method: method,
			Args:   args,
			Loc:    nodes.Loc{Start: loc, End: p.prevEnd()},
		}
		// variableMethodCallTableReference takes an alias and its column list
		// (bnf/TSql170.g). SQL Server rejects @v.f()(c), whose column list has
		// no alias, @v.f() AS a, which has no column list, and @v.f() a(),
		// whose column list is empty.
		mc.Alias = p.parseOptionalAlias()
		if mc.Alias == "" || p.cur.Type != '(' {
			return nil, p.newParseError(loc, "a table variable method call takes an alias and its column list")
		}
		p.advance() // consume (
		for {
			col, ok := p.parseIdentifier()
			if !ok {
				return nil, p.newParseError(p.pos(), "a table variable method call takes a column list")
			}
			mc.Columns = append(mc.Columns, col)
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		mc.Loc.End = p.prevEnd()
		return mc, nil
	}

	tv := &nodes.TableVarRef{
		Name: name,
		Loc:  nodes.Loc{Start: loc, End: p.prevEnd()},
	}
	tv.Alias = p.parseOptionalAlias()
	tv.Loc.End = p.prevEnd()
	return tv, nil
}

// parseVariableDmlTarget parses a bare table variable as a DML target
// (INSERT/UPDATE/DELETE/MERGE). Does not accept alias or table hints —
// mirrors SqlScriptDOM variableDmlTarget.
// Caller must ensure p.cur.Type == tokVARIABLE.
func (p *Parser) parseVariableDmlTarget() *nodes.TableVarRef {
	loc := p.pos()
	name := p.cur.Str
	p.advance()
	return &nodes.TableVarRef{
		Name: name,
		Loc:  nodes.Loc{Start: loc, End: p.prevEnd()},
	}
}

// parseIdentExpr parses an identifier expression (column ref, function call, or qualified name).
// This handles both simple identifiers and dot-qualified references.
//
//	ident_expr = identifier [ '(' args ')' ]
//	           | identifier '.' identifier [ '.' identifier [ '.' identifier ] ]
//	           | identifier '.' '*'
func (p *Parser) parseIdentExpr() (nodes.ExprNode, error) {
	loc := p.pos()
	name := p.cur.Str
	delimited := p.isDelimited(p.cur)
	p.advance()
	if p.collectMode() && p.cursorOff <= p.prev.End {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, errCollecting
	}

	// Function call: ident(...). Without a schema it names a built-in, which
	// SQL Server's parser checks (Msg 195; see builtinFunctions).
	if p.cur.Type == '(' {
		// MATCH ( graph pattern ) is a predicate of a search condition
		// (graphMatchPredicate, TSql160.g). Elsewhere MATCH is not a
		// function, so the built-in check below refuses it (Msg 195).
		if !delimited && p.inSearchCondition() && p.isValidOptionToken(graphMatchWord, p.prev) {
			return p.parseGraphMatchPredicate(loc)
		}
		upper := strings.ToUpper(name)
		if delimited || !builtinFunctions[upper] {
			return nil, &ParseError{Message: fmt.Sprintf("%q is not a recognized built-in function name", name), Position: loc, Code: 195}
		}
		if strings.EqualFold(name, "PARSE") || strings.EqualFold(name, "TRY_PARSE") {
			return p.parseParseExpr(name, loc)
		}
		return p.parseFuncCall(name, loc)
	}

	// Static method call: type::Method(args)
	if p.cur.Type == tokCOLONCOLON {
		p.advance() // consume ::
		// A static method call names its method: SELECT a :: is Msg 102
		// (measured).
		if !p.isIdentLike() {
			return nil, p.syntaxErrorAtCur()
		}
		method := p.cur.Str
		p.advance()
		mc := &nodes.MethodCallExpr{
			Type:   &nodes.DataType{Name: name, Loc: nodes.Loc{Start: loc, End: -1}},
			Method: method,
			Loc:    nodes.Loc{Start: loc, End: -1},
		}
		if p.cur.Type == '(' {
			p.advance() // consume (
			var args []nodes.Node
			if p.cur.Type != ')' {
				for {
					arg, err := p.parseExpr()
					if err != nil {
						return nil, err
					}
					args = append(args, arg)
					if _, ok := p.match(','); !ok {
						break
					}
				}
			}
			mc.Args = &nodes.List{Items: args}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
		mc.Loc.End = p.prevEnd()
		return mc, nil
	}

	// Qualified name: ident.ident[.ident[.ident]] or ident.*
	if p.cur.Type == '.' {
		return p.parseQualifiedRef(name, loc)
	}

	// Simple column reference
	return &nodes.ColumnRef{
		Column: name,
		Loc:    nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
}

// knownPseudoColumns are the $-prefixed pseudo-columns SQL Server accepts as
// column references (lowercased): $action (MERGE OUTPUT), $IDENTITY, $ROWGUID,
// and the graph-table pseudo-columns. $PARTITION is handled separately as a
// partition-function call prefix. The engine rejects anything else with
// Msg 126 "Invalid pseudocolumn"; context restrictions (e.g. $action outside
// MERGE OUTPUT) are binding-time errors in the engine, not parse errors, so
// the parser does not enforce them.
//
// Deliberate divergence from SqlScriptDOM: ScriptDom also accepts $CUID
// (ColumnType.PseudoColumnCuid, present since its SQL 2008 grammar), but no
// shipped engine does — SQL Server 2022 rejects `SELECT $CUID` with Msg 126
// exactly like an unknown pseudo-column, and $CUID appears nowhere in the
// T-SQL documentation. We follow the engine.
var knownPseudoColumns = map[string]bool{
	"$action":   true,
	"$identity": true,
	"$rowguid":  true,
	"$node_id":  true,
	"$edge_id":  true,
	"$from_id":  true,
	"$to_id":    true,
}

// graphPseudoColumns is the subset valid in index key column lists
// (engine-verified: `CREATE INDEX ix ON Person ($node_id)` succeeds).
var graphPseudoColumns = map[string]bool{
	"$node_id": true,
	"$edge_id": true,
	"$from_id": true,
	"$to_id":   true,
}

// parseIndexColumnName parses one index key column name: a regular identifier
// or a graph pseudo-column.
func (p *Parser) parseIndexColumnName() (string, error) {
	if p.cur.Type == tokPSEUDOCOL {
		if !graphPseudoColumns[strings.ToLower(p.cur.Str)] {
			return "", p.newParseError(p.cur.Loc, fmt.Sprintf("invalid pseudocolumn %q", p.cur.Str))
		}
		name := p.cur.Str
		p.advance()
		return name, nil
	}
	name, ok := p.parseIdentifier()
	if !ok {
		return "", p.unexpectedToken()
	}
	return name, nil
}

// parsePseudoColumn parses a $-prefixed pseudo-column token in expression
// position: a known pseudo-column becomes a ColumnRef, $PARTITION.fn(args)
// becomes a schema-qualified function call, anything else is an error.
func (p *Parser) parsePseudoColumn() (nodes.ExprNode, error) {
	tok := p.advance()
	// Not keyword matching: pseudo-columns are tokPSEUDOCOL tokens, never
	// keywords, so the comparison is on the token's raw text.
	if strings.ToLower(tok.Str) == "$partition" {
		// $PARTITION.partition_function_name(expr)
		if _, err := p.expect('.'); err != nil {
			return nil, err
		}
		if !p.isIdentLike() {
			return nil, p.syntaxErrorAtCur()
		}
		fname := p.cur.Str
		p.advance()
		if p.cur.Type != '(' {
			return nil, p.syntaxErrorAtCur()
		}
		fc, err := p.parseFuncCallWithSchema(tok.Str, fname, tok.Loc)
		if err != nil {
			return nil, err
		}
		if err := requirePartitionArgs(p, fc); err != nil {
			return nil, err
		}
		return fc, nil
	}
	if !knownPseudoColumns[strings.ToLower(tok.Str)] {
		return nil, p.newParseError(tok.Loc, fmt.Sprintf("invalid pseudocolumn %q", tok.Str))
	}
	return &nodes.ColumnRef{
		Column: tok.Str,
		Loc:    nodes.Loc{Start: tok.Loc, End: p.prevEnd()},
	}, nil
}

// requirePartitionArgs rejects $PARTITION.fn() with no or malformed
// arguments — the engine requires an expression list and rejects empty,
// comma-only, and trailing-comma forms (all Msg 102, verified on SQL Server
// 2022). The generic function-call parser appends a nil item when an
// expression slot fails to parse, so nil items are rejected here too.
func requirePartitionArgs(p *Parser, fc nodes.ExprNode) error {
	f, ok := fc.(*nodes.FuncCallExpr)
	if !ok {
		return nil
	}
	if f.Star || f.Args == nil || len(f.Args.Items) == 0 {
		return p.newParseError(f.Loc.End, "$PARTITION function requires at least one argument")
	}
	for _, a := range f.Args.Items {
		if a == nil {
			return p.newParseError(f.Loc.End, "$PARTITION function has a malformed argument list")
		}
	}
	return nil
}

// parseInsertColumnName parses one name in an INSERT/MERGE column list: a
// regular identifier or a known pseudo-column — graph edge-table inserts name
// $from_id/$to_id as target columns (engine-verified:
// `INSERT INTO e ($from_id, $to_id) SELECT ...` executes).
func (p *Parser) parseInsertColumnName() (nodes.Node, error) {
	if p.cur.Type == tokPSEUDOCOL {
		if !knownPseudoColumns[strings.ToLower(p.cur.Str)] {
			return nil, p.newParseError(p.cur.Loc, fmt.Sprintf("invalid pseudocolumn %q", p.cur.Str))
		}
		name := p.cur.Str
		p.advance()
		return &nodes.String{Str: name}, nil
	}
	colName, ok := p.parseIdentifier()
	if !ok {
		return nil, p.unexpectedToken()
	}
	return &nodes.String{Str: colName}, nil
}

// parseQualifiedRef parses a dot-qualified column reference or star expression.
// The first part has already been consumed.
//
//	qualified_ref = first '.' ( '*' | ident [ '.' ( '*' | ident [ '.' ( '*' | ident ) ] ) ] )
func (p *Parser) parseQualifiedRef(first string, loc int) (nodes.ExprNode, error) {
	parts := []string{first}
	for p.cur.Type == '.' {
		p.advance() // consume .
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addTokenCandidate('*')
			return nil, errCollecting
		}

		// Check for table.* or schema.table.*
		if p.cur.Type == '*' {
			p.advance()
			// Build qualifier from collected parts
			qualifier := first
			if len(parts) > 1 {
				qualifier = parts[len(parts)-1]
			}
			return &nodes.StarExpr{
				Qualifier: qualifier,
				Loc:       nodes.Loc{Start: loc, End: p.prevEnd()},
			}, nil
		}

		// db.$PARTITION.fn(args) — the partition function accepts one optional
		// database qualifier (engine-verified: `db.$PARTITION.pf(10)` executes;
		// more than one qualifier is a syntax error).
		if p.cur.Type == tokPSEUDOCOL && strings.ToLower(p.cur.Str) == "$partition" {
			if len(parts) != 1 {
				return nil, p.syntaxErrorAtCur()
			}
			partName := p.cur.Str
			p.advance()
			if _, err := p.expect('.'); err != nil {
				return nil, err
			}
			if !p.isIdentLike() {
				return nil, p.syntaxErrorAtCur()
			}
			fname := p.cur.Str
			p.advance()
			if p.cur.Type != '(' {
				return nil, p.syntaxErrorAtCur()
			}
			fc, err := p.parseFuncCallWithSchema(partName, fname, loc)
			if err != nil {
				return nil, err
			}
			if err := requirePartitionArgs(p, fc); err != nil {
				return nil, err
			}
			if f, ok := fc.(*nodes.FuncCallExpr); ok && f.Name != nil {
				f.Name.Database = parts[0]
			}
			return fc, nil
		}

		// Qualified pseudo-column (t.$IDENTITY, Person.$node_id) or keyword
		// column (t.IDENTITYCOL, t.ROWGUIDCOL). Both terminate the chain.
		if p.cur.Type == tokPSEUDOCOL || p.cur.Type == kwIDENTITYCOL || p.cur.Type == kwROWGUIDCOL {
			if p.cur.Type == tokPSEUDOCOL && !knownPseudoColumns[strings.ToLower(p.cur.Str)] {
				return nil, p.newParseError(p.cur.Loc, fmt.Sprintf("invalid pseudocolumn %q", p.cur.Str))
			}
			parts = append(parts, p.cur.Str)
			p.advance()
			break
		}

		// Accept identifier or keyword-as-identifier after dot
		if p.isIdentLike() {
			partName := p.cur.Str
			p.advance()
			if p.collectMode() && p.cursorOff <= p.prev.End {
				p.addRuleCandidate("columnref")
				return nil, errCollecting
			}

			// Check if this part is followed by '(' -- meaning it's a function call
			// e.g., schema.function(args)
			if p.cur.Type == '(' {
				schema := first
				if len(parts) > 1 {
					schema = parts[0]
				}
				return p.parseFuncCallWithSchema(schema, partName, loc)
			}

			// Check for :: static method call: schema.type::Method(args)
			if p.cur.Type == tokCOLONCOLON {
				p.advance() // consume ::
				if !p.isIdentLike() {
					return nil, p.syntaxErrorAtCur()
				}
				method := p.cur.Str
				p.advance()
				dt := &nodes.DataType{Name: partName, Loc: nodes.Loc{Start: loc, End: -1}}
				if len(parts) > 0 {
					dt.Schema = parts[0]
				}
				mc := &nodes.MethodCallExpr{
					Type:   dt,
					Method: method,
					Loc:    nodes.Loc{Start: loc, End: -1},
				}
				if p.cur.Type == '(' {
					p.advance() // consume (
					var args []nodes.Node
					if p.cur.Type != ')' {
						for {
							arg, err := p.parseExpr()
							if err != nil {
								return nil, err
							}
							args = append(args, arg)
							if _, ok := p.match(','); !ok {
								break
							}
						}
					}
					mc.Args = &nodes.List{Items: args}
					if _, err := p.expect(')'); err != nil {
						return nil, err
					}
				}
				mc.Loc.End = p.prevEnd()
				return mc, nil
			}

			parts = append(parts, partName)
		} else {
			// A dot names the next part: SELECT a. and SELECT a.b. are Msg 102
			// (measured).
			return nil, p.syntaxErrorAtCur()
		}
	}

	ref := &nodes.ColumnRef{Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}
	switch len(parts) {
	case 1:
		ref.Column = parts[0]
	case 2:
		ref.Table = parts[0]
		ref.Column = parts[1]
	case 3:
		ref.Schema = parts[0]
		ref.Table = parts[1]
		ref.Column = parts[2]
	case 4:
		ref.Database = parts[0]
		ref.Schema = parts[1]
		ref.Table = parts[2]
		ref.Column = parts[3]
	default: // 5 parts: server.database.schema.table.column
		ref.Server = parts[0]
		ref.Database = parts[1]
		ref.Schema = parts[2]
		ref.Table = parts[3]
		ref.Column = parts[4]
	}
	return ref, nil
}

// parseFuncCallWithSchema parses a schema-qualified function call.
// schema.func(args)
func (p *Parser) parseFuncCallWithSchema(schema, funcName string, loc int) (nodes.ExprNode, error) {
	nameEnd := p.prevEnd()
	p.advance() // consume (

	fc := &nodes.FuncCallExpr{
		Name: &nodes.TableRef{Schema: schema, Object: funcName, Loc: nodes.Loc{Start: loc, End: nameEnd}},
		Loc:  nodes.Loc{Start: loc, End: -1},
	}

	// COUNT(*) special case
	if p.cur.Type == '*' {
		p.advance()
		fc.Star = true
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		fc.Loc.End = p.prevEnd()
		if p.cur.Type == kwOVER {
			var err error
			fc.Over, err = p.parseOverClause()
			if err != nil {
				return nil, err
			}
			fc.Loc.End = p.prevEnd()
		}
		return fc, nil
	}

	if p.collectMode() {
		p.addExpressionCandidates()
		return nil, errCollecting
	}

	if p.cur.Type == ')' {
		p.advance()
		fc.Loc.End = p.prevEnd()
		if p.cur.Type == kwOVER {
			var err error
			fc.Over, err = p.parseOverClause()
			if err != nil {
				return nil, err
			}
			fc.Loc.End = p.prevEnd()
		}
		return fc, nil
	}

	// Check for DISTINCT
	if _, ok := p.match(kwDISTINCT); ok {
		fc.Distinct = true
	}

	var args []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.collectMode() {
			p.addExpressionCandidates()
			return nil, errCollecting
		}
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if arg == nil {
			return nil, p.syntaxErrorAtCur()
		}
		args = append(args, arg)
		if _, ok := p.match(','); !ok {
			break
		}
		if p.collectMode() {
			p.addExpressionCandidates()
			return nil, errCollecting
		}
		// A trailing comma before ')' is a syntax error (engine: Msg 102),
		// matching the unqualified parseFuncCall path.
		if p.cur.Type == ')' {
			return nil, p.syntaxErrorAtCur()
		}
	}
	fc.Args = &nodes.List{Items: args}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	fc.Loc.End = p.prevEnd()

	// Check for OVER clause
	if p.cur.Type == kwOVER {
		var err error
		fc.Over, err = p.parseOverClause()
		if err != nil {
			return nil, err
		}
		fc.Loc.End = p.prevEnd()
	}

	return fc, nil
}
