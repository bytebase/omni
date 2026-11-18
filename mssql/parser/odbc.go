// Package parser - odbc.go implements the ODBC escape sequences of a scalar
// expression.
package parser

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// odbcFunction is the number of arguments an ODBC function escape takes.
type odbcFunction struct {
	min, max int
}

// odbcFunctions holds the ODBC scalar functions SQL Server's parser takes in
// { fn <name> ( <arguments> ) }, with the argument counts the pinned server
// accepts. SQL Server maps each one to the T-SQL function of the same meaning
// and refuses a name it does not know (Msg 195), so ScriptDOM's open
// `Identifier LeftParenthesis ...` alternative is narrower at the server.
// The counts were asked of the pinned SQL Server under SET PARSEONLY ON.
//
// Ref: https://learn.microsoft.com/en-us/sql/odbc/reference/develop-app/odbc-escape-sequences
var odbcFunctions = map[string]odbcFunction{
	// String functions.
	"ASCII":        {1, 1},
	"BIT_LENGTH":   {1, 1},
	"CHAR":         {1, 1},
	"CONCAT":       {2, 2},
	"DIFFERENCE":   {2, 2},
	"INSERT":       {4, 4},
	"LCASE":        {1, 1},
	"LEFT":         {2, 2},
	"LENGTH":       {1, 1},
	"LOCATE":       {2, 3},
	"LTRIM":        {1, 2},
	"OCTET_LENGTH": {1, 1},
	"REPEAT":       {2, 2},
	"REPLACE":      {3, 3},
	"RIGHT":        {2, 2},
	"RTRIM":        {1, 2},
	"SOUNDEX":      {1, 1},
	"SPACE":        {1, 1},
	"SUBSTRING":    {3, 3},
	"UCASE":        {1, 1},

	// Numeric functions.
	"ABS":      {1, 1},
	"ACOS":     {1, 1},
	"ASIN":     {1, 1},
	"ATAN":     {1, 1},
	"ATAN2":    {2, 2},
	"CEILING":  {1, 1},
	"COS":      {1, 1},
	"COT":      {1, 1},
	"DEGREES":  {1, 1},
	"EXP":      {1, 1},
	"FLOOR":    {1, 1},
	"LOG":      {1, 2},
	"LOG10":    {1, 1},
	"MOD":      {2, 2},
	"PI":       {0, 0},
	"POWER":    {2, 2},
	"RADIANS":  {1, 1},
	"RAND":     {0, 1},
	"ROUND":    {2, 2},
	"SIGN":     {1, 1},
	"SIN":      {1, 1},
	"SQRT":     {1, 1},
	"TAN":      {1, 1},
	"TRUNCATE": {2, 2},

	// Date and time functions.
	"CURRENT_DATE":      {0, 0},
	"CURRENT_TIME":      {0, 1},
	"CURRENT_TIMESTAMP": {0, 1},
	"CURDATE":           {0, 0},
	"CURTIME":           {0, 0},
	"DAYNAME":           {1, 1},
	"DAYOFMONTH":        {1, 1},
	"DAYOFWEEK":         {1, 1},
	"DAYOFYEAR":         {1, 1},
	"HOUR":              {1, 1},
	"MINUTE":            {1, 1},
	"MONTH":             {1, 1},
	"MONTHNAME":         {1, 1},
	"NOW":               {0, 0},
	"QUARTER":           {1, 1},
	"SECOND":            {1, 1},
	"WEEK":              {1, 1},
	"YEAR":              {1, 1},

	// System functions.
	"DATABASE": {0, 0},
	"IFNULL":   {2, 2},
	"USER":     {0, 0},
}

// odbcConvertTypes are the type names { fn CONVERT ( <expr> , <type> ) } takes
// (odbcConvertSpecification, bnf/TSql170.g). They are the ODBC SQL types, not
// T-SQL's: the pinned server refuses INT, varchar and SQL_SS_TIME2 there, and
// takes these.
var odbcConvertTypes = map[string]bool{
	"SQL_BIGINT":        true,
	"SQL_BINARY":        true,
	"SQL_BIT":           true,
	"SQL_CHAR":          true,
	"SQL_DATE":          true,
	"SQL_DECIMAL":       true,
	"SQL_DOUBLE":        true,
	"SQL_FLOAT":         true,
	"SQL_GUID":          true,
	"SQL_INTEGER":       true,
	"SQL_LONGVARBINARY": true,
	"SQL_LONGVARCHAR":   true,
	"SQL_NUMERIC":       true,
	"SQL_REAL":          true,
	"SQL_SMALLINT":      true,
	"SQL_TIME":          true,
	"SQL_TIMESTAMP":     true,
	"SQL_TINYINT":       true,
	"SQL_VARBINARY":     true,
	"SQL_VARCHAR":       true,
	"SQL_WCHAR":         true,
	"SQL_WLONGVARCHAR":  true,
	"SQL_WVARCHAR":      true,
}

// odbcExtractParts are the parts { fn EXTRACT ( <part> FROM <expr> ) } takes.
// They are the six parts of the ODBC escape with T-SQL's abbreviations, and
// the pinned server refuses the rest: quarter, week, dayofyear, millisecond
// and nanosecond are all errors there.
var odbcExtractParts = map[string]bool{
	"YEAR": true, "YY": true, "YYYY": true,
	"MONTH": true, "MM": true, "M": true,
	"DAY": true, "DD": true, "D": true,
	"HOUR": true, "HH": true,
	"MINUTE": true, "MI": true, "N": true,
	"SECOND": true, "SS": true, "S": true,
}

// odbcTimestampIntervals are the interval names { fn TIMESTAMPADD (...) } and
// { fn TIMESTAMPDIFF (...) } take as their first argument.
var odbcTimestampIntervals = map[string]bool{
	"SQL_TSI_FRAC_SECOND": true,
	"SQL_TSI_SECOND":      true,
	"SQL_TSI_MINUTE":      true,
	"SQL_TSI_HOUR":        true,
	"SQL_TSI_DAY":         true,
	"SQL_TSI_WEEK":        true,
	"SQL_TSI_MONTH":       true,
	"SQL_TSI_QUARTER":     true,
	"SQL_TSI_YEAR":        true,
}

// parseOdbcEscape parses an ODBC escape at the current '{'.
//
//	odbcFunctionCall: { fn <function> }
//	odbcLiteral:      { <kind> '<literal>' }
//
// SQL Server reads both in a scalar expression: SELECT { fn NOW() } and
// SELECT { d '2020-01-01' } are expressions. ScriptDOM's odbcQualifiedJoin,
// the form that wraps a join, refuses every input in its action, so a '{' in
// a FROM clause is an error here too.
//
// Ref: https://learn.microsoft.com/en-us/sql/odbc/reference/develop-app/odbc-escape-sequences
func (p *Parser) parseOdbcEscape() (nodes.ExprNode, error) {
	loc := p.pos()
	p.advance() // consume {

	// Completion: after { in an expression → the escape's introducer
	if p.collectMode() {
		return nil, errCollecting
	}

	if p.atOdbcFnIntroducer() {
		return p.parseOdbcFunctionCall(loc)
	}

	// odbcLiteral: { d|t|ts '<literal>' }. The grammar takes any identifier
	// as the kind; the escape the ODBC specification defines takes d, t and
	// ts, and SQL Server refuses the rest (Msg 102).
	kindLoc := p.pos()
	kind, ok := p.parseIdentifier()
	if !ok {
		return nil, p.unexpectedToken()
	}
	switch strings.ToLower(kind) {
	case "d", "t", "ts":
	default:
		return nil, p.newParseError(kindLoc, "an ODBC literal escape takes d, t or ts")
	}
	if p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST {
		return nil, p.unexpectedToken()
	}
	value, err := p.parseStringLiteral()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect('}'); err != nil {
		return nil, err
	}
	return &nodes.OdbcLiteralExpr{Kind: kind, Value: value, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
}

// odbcFnWord is FN, the introducer of an ODBC function escape. It is an
// identifier to ScriptDOM's lexer, not a keyword.
var odbcFnWord = newOptionSet().withIdents("FN")

// atOdbcFnIntroducer reports whether the current token introduces an ODBC
// function escape.
func (p *Parser) atOdbcFnIntroducer() bool {
	return !p.isDelimited(p.cur) && p.isValidOption(odbcFnWord)
}

// parseOdbcFunctionCall parses { fn <function> } after the introducer.
//
//	odbcFunctionCall, bnf/TSql170.g:
//	    { fn CONVERT ( <expr> , <odbc-type> ) }
//	    { fn EXTRACT ( <part> FROM <expr> ) }
//	    { fn TRUNCATE ( <expr> , <expr> ) }
//	    { fn DATABASE | USER | CURRENT_DATE ( ) }
//	    { fn INSERT | LEFT | RIGHT ( <expression-list> ) }
//	    { fn CURRENT_TIME | CURRENT_TIMESTAMP [ ( <expression-list> ) ] }
//	    { fn <identifier> ( [ <expression-list> ] ) }
//
// The last alternative is a name SQL Server knows (odbcFunctions), and the
// other forms are the named alternatives, each with its own argument count.
func (p *Parser) parseOdbcFunctionCall(loc int) (nodes.ExprNode, error) {
	p.advance() // consume fn
	if p.cur.Str == "" || p.isDelimited(p.cur) || p.cur.Type == tokEOF {
		return nil, p.unexpectedToken()
	}
	nameLoc := p.pos()
	name := strings.ToUpper(p.cur.Str)
	p.advance()

	fc := &nodes.OdbcFunctionCallExpr{Name: name, Loc: nodes.Loc{Start: loc, End: -1}}

	switch name {
	case "CONVERT":
		// CONVERT ( expression , odbcConvertSpecification ), whose second
		// argument is a type name without quotes.
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(','); err != nil {
			return nil, err
		}
		typeLoc := p.pos()
		typ, ok := p.parseIdentifier()
		if !ok {
			return nil, p.unexpectedToken()
		}
		if !odbcConvertTypes[strings.ToUpper(typ)] {
			return nil, p.newParseError(typeLoc, fmt.Sprintf("%q is not an ODBC type name", typ))
		}
		fc.Convert = typ
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		fc.Args = &nodes.List{Items: []nodes.Node{arg}}

	case "EXTRACT":
		// EXTRACT ( <part> FROM <expr> )
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		partLoc := p.pos()
		part, ok := p.parseIdentifier()
		if !ok {
			return nil, p.unexpectedToken()
		}
		if !odbcExtractParts[strings.ToUpper(part)] {
			return nil, p.newParseError(partLoc, fmt.Sprintf("%q is not a part EXTRACT takes", part))
		}
		fc.Extract = part
		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		fc.Args = &nodes.List{Items: []nodes.Node{arg}}

	case "TIMESTAMPADD", "TIMESTAMPDIFF":
		// TIMESTAMPADD | TIMESTAMPDIFF ( SQL_TSI_<unit> , <expr> , <expr> )
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		unitLoc := p.pos()
		unit, ok := p.parseIdentifier()
		if !ok {
			return nil, p.unexpectedToken()
		}
		if !odbcTimestampIntervals[strings.ToUpper(unit)] {
			return nil, p.newParseError(unitLoc, fmt.Sprintf("%q is not an interval %s takes", unit, name))
		}
		if _, err := p.expect(','); err != nil {
			return nil, err
		}
		args, err := p.parseOdbcExprList()
		if err != nil {
			return nil, err
		}
		if len(args) != 2 {
			return nil, p.newParseError(nameLoc, name+" takes an interval and two expressions")
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		fc.Args = &nodes.List{Items: args}

	default:
		spec, ok := odbcFunctions[name]
		if !ok {
			return nil, p.newParseError(nameLoc, fmt.Sprintf("%q is not a recognized ODBC function name", name))
		}
		// Only CURRENT_TIME and CURRENT_TIMESTAMP may be written without
		// parentheses (odbcFunctionCall, bnf/TSql170.g).
		optionalParens := name == "CURRENT_TIME" || name == "CURRENT_TIMESTAMP"
		var args []nodes.Node
		if optionalParens && p.cur.Type != '(' {
			break
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		var err error
		args, err = p.parseOdbcOptExprList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		if len(args) < spec.min || len(args) > spec.max {
			return nil, p.newParseError(nameLoc, fmt.Sprintf("%s takes %d to %d arguments", name, spec.min, spec.max))
		}
		if len(args) > 0 {
			fc.Args = &nodes.List{Items: args}
		}
	}

	if _, err := p.expect('}'); err != nil {
		return nil, err
	}
	fc.Loc.End = p.prevEnd()
	return fc, nil
}

// parseOdbcExprList parses one or more expressions up to, not including, ')'.
func (p *Parser) parseOdbcExprList() ([]nodes.Node, error) {
	var args []nodes.Node
	for {
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if arg == nil {
			return nil, p.unexpectedToken()
		}
		args = append(args, arg)
		if _, ok := p.match(','); !ok {
			return args, nil
		}
	}
}

// parseOdbcOptExprList parses zero or more expressions up to, not including,
// ')'.
func (p *Parser) parseOdbcOptExprList() ([]nodes.Node, error) {
	if p.cur.Type == ')' {
		return nil, nil
	}
	return p.parseOdbcExprList()
}
