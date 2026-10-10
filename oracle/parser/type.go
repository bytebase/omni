package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseTypeName parses a data type specification.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/Data-Types.html
//
//	datatype ::=
//	    NUMBER [ ( precision [, scale] ) ]
//	  | FLOAT [ ( precision ) ]
//	  | INTEGER | SMALLINT | DECIMAL [ ( precision [, scale] ) ]
//	  | CHAR [ ( size [ BYTE | CHAR ] ) ]
//	  | VARCHAR2 ( size [ BYTE | CHAR ] )
//	  | VARCHAR ( size [ BYTE | CHAR ] )
//	  | NCHAR [ ( size ) ]
//	  | NVARCHAR2 ( size )
//	  | CLOB | BLOB | NCLOB
//	  | DATE
//	  | TIMESTAMP [ ( precision ) ] [ WITH [ LOCAL ] TIME ZONE ]
//	  | INTERVAL YEAR [ ( precision ) ] TO MONTH
//	  | INTERVAL DAY [ ( precision ) ] TO SECOND [ ( precision ) ]
//	  | RAW ( size )
//	  | LONG [ RAW ]
//	  | ROWID
//	  | ref%TYPE | ref%ROWTYPE
//	  | [ schema . ] type_name
func (p *Parser) parseTypeName() (*nodes.TypeName, error) {
	start := p.pos()
	tn := &nodes.TypeName{
		Names:    &nodes.List{},
		TypeMods: &nodes.List{},
		Loc:      nodes.Loc{Start: start},
	}

	switch p.cur.Type {
	case kwNUMBER, kwINTEGER, kwSMALLINT, kwDECIMAL, kwFLOAT:
		tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: p.cur.Str})
		p.advance()
		parseErr1116 := p.parseOptionalPrecisionScale(tn)
		if parseErr1116 != nil {
			return nil, parseErr1116
		}

	case kwCHAR, kwVARCHAR2, kwVARCHAR, kwNCHAR, kwNVARCHAR2:
		tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: p.cur.Str})
		p.advance()
		parseErr1117 := p.parseOptionalSizeWithSemantic(tn)
		if parseErr1117 != nil {
			return nil, parseErr1117
		}

	case kwCLOB, kwBLOB, kwNCLOB, kwJSON:
		tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: p.cur.Str})
		p.advance()

	case kwDATE:
		tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: "DATE"})
		p.advance()

	case kwTIMESTAMP:
		tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: "TIMESTAMP"})
		p.advance()
		parseErr1118 := p.parseOptionalPrecisionScale(tn)
		if parseErr1118 !=
			// WITH [ LOCAL ] TIME ZONE
			nil {
			return nil, parseErr1118
		}

		if p.cur.Type == kwWITH {
			tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: "WITH"})
			p.advance()
			if p.cur.Type == kwLOCAL {
				tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: "LOCAL"})
				p.advance()
			}
			// TIME
			if p.isIdentLikeStr("TIME") {
				tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: "TIME"})
				p.advance()
			}
			// ZONE
			if p.cur.Type == kwZONE {
				tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: "ZONE"})
				p.advance()
			}
		}

	case kwINTERVAL:
		tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: "INTERVAL"})
		p.advance()
		parseErr1119 := p.parseIntervalType(tn)
		if parseErr1119 != nil {
			return nil, parseErr1119
		}

	case kwRAW:
		tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: "RAW"})
		p.advance()
		parseErr1120 := p.parseOptionalPrecisionScale(tn)
		if parseErr1120 != nil {
			return nil, parseErr1120
		}

	case kwLONG:
		tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: "LONG"})
		p.advance()
		if p.cur.Type == kwRAW {
			tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: "RAW"})
			p.advance()
		}

	case kwROWID:
		tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: "ROWID"})
		p.advance()

	default:
		if p.isIdentLikeStr("NUMERIC") {
			tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: p.cur.Str})
			p.advance()
			parseErr1121 := p.parseOptionalPrecisionScale(tn)
			if parseErr1121 != nil {
				return nil, parseErr1121
			}
		} else {
			parseErr1121 :=
				// User-defined type or %TYPE/%ROWTYPE reference
				p.parseUserDefinedType(tn)
			if parseErr1121 != nil {
				return nil, parseErr1121
			}
		}
	}

	tn.Loc.End = p.prev.End
	return tn, nil
}

// typeModMode says what the parenthesized modifiers of a datatype (length,
// precision, scale) may hold.
type typeModMode int

const (
	// typeModsSQL takes integer literals, as SQL does.
	typeModsSQL typeModMode = iota
	// typeModsPLSQL takes expressions: the datatype of a PL/SQL declaration.
	// Where SQL takes an integer literal, PL/SQL takes a static expression,
	// VARCHAR2(ORA_MAX_NAME_LEN + 2) or NUMBER(pkg.p, 2). That the expression
	// is static is checked when the unit compiles (PLS-00491 on Oracle 23ai),
	// not when it parses, so any expression parses.
	//
	// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/plsql-language-fundamentals.html (Static Expressions)
	typeModsPLSQL
	// typeModsNone takes no modifiers: the type of a parameter or of a
	// function result is unconstrained. Oracle 23ai rejects p VARCHAR2(10),
	// RETURN NUMBER(5), and INTERVAL DAY(2) TO SECOND there with PLS-00103,
	// in standalone, package, local, and object type subprograms and cursor
	// parameters alike.
	typeModsNone
)

// plsqlCharsetUse says which CHARACTER SET clause a PL/SQL datatype takes.
type plsqlCharsetUse int

const (
	// charsetNone takes no CHARACTER SET clause.
	charsetNone plsqlCharsetUse = iota
	// charsetFixed takes a named character set only: a record field or a
	// collection element (PLS-00552 for ANY_CS or item%CHARSET).
	charsetFixed
	// charsetFlexible also takes item%CHARSET: a declaration, a subtype, a
	// function result, or a cursor parameter.
	charsetFlexible
	// charsetAnyCS also takes ANY_CS: a subprogram parameter, the one place
	// Oracle allows it (PLS-00551 elsewhere).
	charsetAnyCS
)

// parsePLSQLTypeName parses the datatype of a PL/SQL declaration whose
// modifiers may be expressions and which takes no CHARACTER SET clause.
func (p *Parser) parsePLSQLTypeName() (*nodes.TypeName, error) {
	return p.parsePLSQLDatatype(typeModsPLSQL, charsetNone)
}

// parsePLSQLDatatype parses a PL/SQL datatype: mods says what its modifiers
// may hold, use which CHARACTER SET clause may follow it.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/datatype-attribute.html
//
//	datatype [ CHARACTER SET { character_set | item%CHARSET } ]
func (p *Parser) parsePLSQLDatatype(mods typeModMode, use plsqlCharsetUse) (*nodes.TypeName, error) {
	saved := p.typeMods
	p.typeMods = mods
	tn, err := p.parseTypeName()
	p.typeMods = saved
	if err != nil {
		return nil, err
	}
	if use != charsetNone && p.cur.Type == tokIDENT && p.cur.Str == "CHARACTER" && p.peekNext().Type == kwSET {
		if err := p.parsePLSQLCharacterSet(tn, use); err != nil {
			return nil, err
		}
	}
	return tn, nil
}

// plsqlNonCharacterTypes are the predefined datatypes, SQL and PL/SQL, that
// take no CHARACTER SET clause: Oracle 23ai rejects BOOLEAN or NUMBER
// CHARACTER SET ANY_CS with PLS-00550, and NCHAR, NVARCHAR2, NCLOB, and LONG
// ones, whose character set is fixed, with PLS-00554. CHAR, CHARACTER,
// VARCHAR2, VARCHAR, STRING, and CLOB take one. A name the text does not
// settle, such as a subtype, an object type, or t%TYPE, is left to the
// engine.
var plsqlNonCharacterTypes = map[string]bool{
	"NUMBER": true, "INTEGER": true, "INT": true, "SMALLINT": true, "DECIMAL": true,
	"DEC": true, "NUMERIC": true, "FLOAT": true, "REAL": true, "DOUBLE": true,
	"BINARY_FLOAT": true, "BINARY_DOUBLE": true, "SIMPLE_FLOAT": true, "SIMPLE_DOUBLE": true,
	"PLS_INTEGER": true, "BINARY_INTEGER": true, "SIMPLE_INTEGER": true, "NATURAL": true,
	"NATURALN": true, "POSITIVE": true, "POSITIVEN": true, "SIGNTYPE": true, "BOOLEAN": true,
	"BOOL": true,
	"DATE": true, "TIMESTAMP": true, "INTERVAL": true, "RAW": true, "LONG": true,
	"ROWID": true, "UROWID": true, "BLOB": true, "BFILE": true, "JSON": true, "VECTOR": true,
	"MLSLABEL": true, "SYS_REFCURSOR": true, "NCHAR": true, "NVARCHAR2": true, "NCLOB": true,
}

// predefinedTypeLead returns the name of the predefined datatype tn spells,
// "" for a name the text does not settle (a subtype, a schema-qualified
// type, t%TYPE, or a quoted name: "NUMBER" may be a user subtype of
// VARCHAR2, and Oracle 23ai compiles one as a SQL macro result). Multiword
// predefined types answer their first word: TIMESTAMP WITH [LOCAL] TIME
// ZONE, INTERVAL DAY|YEAR TO ..., LONG RAW.
func (p *Parser) predefinedTypeLead(tn *nodes.TypeName) string {
	if tn == nil || tn.IsPercType || tn.IsPercRowtype || tn.Names.Len() == 0 {
		return ""
	}
	if tn.Loc.Start >= 0 && tn.Loc.Start < len(p.source) && p.source[tn.Loc.Start] == '"' {
		return ""
	}
	word := func(i int) string {
		s, _ := tn.Names.Items[i].(*nodes.String)
		if s == nil {
			return ""
		}
		return s.Str
	}
	if tn.Names.Len() == 1 {
		return word(0)
	}
	switch first, second := word(0), word(1); {
	case first == "TIMESTAMP" && second == "WITH",
		first == "INTERVAL" && (second == "DAY" || second == "YEAR"),
		first == "LONG" && second == "RAW":
		return first
	}
	return ""
}

// parsePLSQLCharacterSet parses CHARACTER SET { character_set | item%CHARSET }
// after the datatype tn, the current token at CHARACTER. Each check below
// follows from the text alone and was confirmed on Oracle 23ai; whether a
// name is a known character set (PLS-00553) is left to the engine.
func (p *Parser) parsePLSQLCharacterSet(tn *nodes.TypeName, use plsqlCharsetUse) error {
	// A %ROWTYPE record is never a character type (PLS-00550).
	if tn.IsPercRowtype {
		return p.syntaxErrorAtCur()
	}
	if plsqlNonCharacterTypes[p.predefinedTypeLead(tn)] {
		return p.syntaxErrorAtCur()
	}
	p.advance() // consume CHARACTER
	p.advance() // consume SET
	nameTok := p.cur
	if !p.isIdentLike() {
		return p.syntaxErrorAtCur()
	}
	name, err := p.parseIdentifier()
	if err != nil {
		return err
	}
	qualified := false
	for p.cur.Type == '.' {
		p.advance()
		if !p.isIdentLike() {
			return p.syntaxErrorAtCur()
		}
		part, err := p.parseIdentifier()
		if err != nil {
			return err
		}
		name += "." + part
		qualified = true
	}
	if p.cur.Type == '%' {
		p.advance()
		// %CHARSET is an attribute keyword: Oracle 23ai rejects a%"CHARSET".
		if !p.isKeywordStr("CHARSET") {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		tn.IsPercCharset = true
		if qualified {
			if p.qualifiedCharsets == nil {
				p.qualifiedCharsets = map[*nodes.TypeName]bool{}
			}
			p.qualifiedCharsets[tn] = true
		}
	} else if qualified {
		return p.syntaxErrorAtCur()
	}
	switch {
	case tn.IsPercCharset && use == charsetFixed:
		// PLS-00552: flexible character set is not allowed on component element
		return p.syntaxErrorAtTok(nameTok)
	case !tn.IsPercCharset && nameTok.Type != tokQIDENT && name == "ANY_CS" && use != charsetAnyCS:
		// PLS-00551: character set ANY_CS is only allowed on a subprogram parameter
		return p.syntaxErrorAtTok(nameTok)
	}
	tn.CharacterSet = name
	tn.Loc.End = p.prev.End
	return nil
}

// parsePLSQLTypeMods parses the parenthesized modifiers of a datatype in a
// PL/SQL declaration, the current token at '(':
//
//	( { expr | * } [ , expr ] )    precision and scale (sized is false)
//	( expr [ BYTE | CHAR ] )       size and length semantics (sized is true)
//
// An integer literal is recorded as *Integer, as in SQL; any other modifier
// is recorded as its expression.
func (p *Parser) parsePLSQLTypeMods(tn *nodes.TypeName, sized bool) error {
	p.advance() // consume '('
	if !sized && p.cur.Type == '*' {
		tn.TypeMods.Items = append(tn.TypeMods.Items, &nodes.String{Str: "*"})
		p.advance()
	} else if err := p.parsePLSQLTypeMod(tn); err != nil {
		return err
	}
	switch {
	case sized && p.isIdentLikeStr("BYTE"):
		tn.TypeMods.Items = append(tn.TypeMods.Items, &nodes.String{Str: "BYTE"})
		p.advance()
	case sized && p.cur.Type == kwCHAR:
		tn.TypeMods.Items = append(tn.TypeMods.Items, &nodes.String{Str: "CHAR"})
		p.advance()
	case !sized && p.cur.Type == ',':
		p.advance()
		if err := p.parsePLSQLTypeMod(tn); err != nil {
			return err
		}
	}
	if p.cur.Type != ')' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// parsePLSQLTypeMod parses one modifier expression. A datatype nested in
// it, as in CAST(x AS NUMBER(5)), follows the SQL rules again.
func (p *Parser) parsePLSQLTypeMod(tn *nodes.TypeName) error {
	p.typeMods = typeModsSQL
	expr, err := p.parseExpr()
	p.typeMods = typeModsPLSQL
	if err != nil {
		return err
	}
	if expr == nil {
		return p.syntaxErrorAtCur()
	}
	if lit, ok := expr.(*nodes.NumberLiteral); ok && !lit.IsFloat {
		tn.TypeMods.Items = append(tn.TypeMods.Items, &nodes.Integer{Ival: lit.Ival})
	} else {
		tn.TypeMods.Items = append(tn.TypeMods.Items, expr)
	}
	return nil
}

// parseOptionalPrecisionScale parses optional ( precision [, scale ] ).
func (p *Parser) parseOptionalPrecisionScale(tn *nodes.TypeName) error {
	if p.cur.Type != '(' {
		return nil
	}
	switch p.typeMods {
	case typeModsPLSQL:
		return p.parsePLSQLTypeMods(tn, false)
	case typeModsNone:
		return p.syntaxErrorAtCur()
	}
	p.advance() // consume '('

	// precision
	if p.cur.Type == tokICONST {
		tn.TypeMods.Items = append(tn.TypeMods.Items, &nodes.Integer{Ival: p.cur.Ival})
		p.advance()
	} else if p.cur.Type == '*' {
		// NUMBER(*) or NUMBER(*,s)
		tn.TypeMods.Items = append(tn.TypeMods.Items, &nodes.String{Str: "*"})
		p.advance()
	}

	// optional scale
	if p.cur.Type == ',' {
		p.advance()
		if p.cur.Type == tokICONST {
			tn.TypeMods.Items = append(tn.TypeMods.Items, &nodes.Integer{Ival: p.cur.Ival})
			p.advance()
		}
	}

	if p.cur.Type == ')' {
		p.advance()
	}
	return nil
}

// parseOptionalSizeWithSemantic parses optional ( size [ BYTE | CHAR ] ).
func (p *Parser) parseOptionalSizeWithSemantic(tn *nodes.TypeName) error {
	if p.cur.Type != '(' {
		return nil
	}
	switch p.typeMods {
	case typeModsPLSQL:
		return p.parsePLSQLTypeMods(tn, true)
	case typeModsNone:
		return p.syntaxErrorAtCur()
	}
	p.advance() // consume '('

	// size
	if p.cur.Type == tokICONST {
		tn.TypeMods.Items = append(tn.TypeMods.Items, &nodes.Integer{Ival: p.cur.Ival})
		p.advance()
	}

	// optional BYTE or CHAR semantic
	if p.isIdentLikeStr("BYTE") {
		tn.TypeMods.Items = append(tn.TypeMods.Items, &nodes.String{Str: "BYTE"})
		p.advance()
	} else if p.cur.Type == kwCHAR {
		tn.TypeMods.Items = append(tn.TypeMods.Items, &nodes.String{Str: "CHAR"})
		p.advance()
	}

	if p.cur.Type == ')' {
		p.advance()
	}
	return nil
}

// parseIntervalType parses the rest of an INTERVAL type after the INTERVAL keyword.
//
//	INTERVAL YEAR [ ( precision ) ] TO MONTH
//	INTERVAL DAY [ ( precision ) ] TO SECOND [ ( precision ) ]
func (p *Parser) parseIntervalType(tn *nodes.TypeName) error {
	// YEAR or DAY
	if p.isIdentLikeStr("YEAR") || p.isIdentLikeStr("DAY") {
		tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: p.cur.Str})
		p.advance()
	}
	parseErr1122 :=

		// optional ( precision )
		p.parseOptionalPrecisionScale(tn)
	if parseErr1122 !=

		// TO
		nil {
		return parseErr1122
	}

	if p.cur.Type == kwTO {
		tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: "TO"})
		p.advance()
	}

	// MONTH or SECOND
	if p.isIdentLikeStr("MONTH") || p.isIdentLikeStr("SECOND") {
		tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: p.cur.Str})
		p.advance()
	}
	parseErr1123 :=

		// SECOND may have ( precision )
		p.parseOptionalPrecisionScale(tn)
	if parseErr1123 !=

		// parseUserDefinedType parses a user-defined type name, possibly schema-qualified,
		// and checks for %TYPE / %ROWTYPE suffixes.
		nil {
		return parseErr1123
	}
	return nil
}

func (p *Parser) parseUserDefinedType(tn *nodes.TypeName) error {
	if !p.isIdentLike() {
		return nil
	}

	name1, parseErr1124 := p.parseIdentifier()
	if parseErr1124 != nil {
		return parseErr1124
	}
	tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: name1})

	// Check for schema.type or ref.column (before %TYPE)
	if p.cur.Type == '.' {
		p.advance()
		// Check for %TYPE / %ROWTYPE after the dot-separated names
		name2, parseErr1125 := p.parseIdentifier()
		if parseErr1125 != nil {
			return parseErr1125
		}
		if name2 != "" {
			tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: name2})

			// Could be schema.table.column%TYPE
			if p.cur.Type == '.' {
				p.advance()
				name3, parseErr1126 := p.parseIdentifier()
				if parseErr1126 != nil {
					return parseErr1126
				}
				if name3 != "" {
					tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: name3})
				}
			}
		}
	}

	// %TYPE or %ROWTYPE
	if p.cur.Type == '%' {
		p.advance()
		if p.cur.Type == kwTYPE {
			tn.IsPercType = true
			p.advance()
		} else if p.cur.Type == kwROWTYPE {
			tn.IsPercRowtype = true
			p.advance()
		}
	}
	return nil
}

// isIdentLikeStr checks if the current token is an identifier-like token with the given uppercase string.
func (p *Parser) isIdentLikeStr(s string) bool {
	return p.isIdentLike() && p.cur.Str == s
}

// isKeywordStr reports whether the current token is the unquoted word s.
// Unlike isIdentLikeStr it never matches a quoted identifier: Oracle reads
// "NOKEEP" as an identifier, not as the NOKEEP keyword (ORA-03049).
func (p *Parser) isKeywordStr(s string) bool {
	return p.cur.Type != tokQIDENT && p.isIdentLikeStr(s)
}
