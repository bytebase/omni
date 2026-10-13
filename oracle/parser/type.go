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
//	  | { CHAR | CHARACTER } [ VARYING ] [ ( size [ BYTE | CHAR ] ) ]
//	  | NATIONAL { CHAR | CHARACTER } [ VARYING ] [ ( size ) ]
//	  | NCHAR VARYING ( size )
//	  | CHARACTER LARGE OBJECT          -- PL/SQL only
//	  | LONG VARCHAR                    -- SQL only
//	  | DOUBLE PRECISION
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
//	  | LONG [ RAW ] [ ( size ) ]         -- size in PL/SQL only
//	  | ROWID
//	  | ref%TYPE | ref%ROWTYPE
//	  | REF [ schema . ] object_type
//	  | [ schema . ] type_name
//
// The ANSI spellings (CHARACTER VARYING, NATIONAL CHARACTER, DOUBLE
// PRECISION) and REF types compile in SQL and PL/SQL alike on Oracle 23ai;
// DOUBLE takes PRECISION (ORA-02000, PLS-00103).
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
		varying := p.cur.Type == kwCHAR || p.cur.Type == kwNCHAR
		national := p.cur.Type == kwNCHAR || p.cur.Type == kwNVARCHAR2
		tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: p.cur.Str})
		p.advance()
		if varying && p.isKeywordStr("VARYING") {
			tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: "VARYING"})
			p.advance()
			if err := p.requireVaryingSize(); err != nil {
				return nil, err
			}
		} else if p.prev.Type == kwNCHAR && p.typeMods != typeModsSQL && p.isKeywordStr("LARGE") {
			// NCHAR LARGE OBJECT, PL/SQL's NCLOB.
			if err := p.parseLargeObjectWords(tn); err != nil {
				return nil, err
			}
			break
		}
		parseErr1117 := p.parseOptionalSizeWithSemantic(tn)
		if parseErr1117 != nil {
			return nil, parseErr1117
		}
		if national {
			if err := p.checkNationalLength(tn); err != nil {
				return nil, err
			}
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
		} else if p.cur.Type == kwVARCHAR && p.typeMods == typeModsSQL {
			// LONG VARCHAR is SQL's alone (PLS-00103 in PL/SQL).
			tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: "VARCHAR"})
			p.advance()
		}
		// PL/SQL takes one length, LONG(100) or LONG RAW(32767), with no
		// scale or BYTE/CHAR (PLS-00572); SQL takes none (ORA-03074), and
		// neither does an unconstrained parameter (PLS-00103).
		if p.typeMods != typeModsSQL && p.cur.Type == '(' {
			if p.typeMods == typeModsNone {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance() // consume '('
			size, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if size == nil || p.cur.Type != ')' {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance() // consume ')'
			tn.TypeMods.Items = append(tn.TypeMods.Items, size)
		}

	case kwROWID:
		tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: "ROWID"})
		p.advance()

	case kwREF:
		// REF [ schema. ] object_type, no attribute (ORA-00911, PLS-00103
		// for %TYPE) and no third component (ORA-03099).
		p.advance() // consume REF
		name, err := p.parseUnitName()
		if err != nil {
			return nil, err
		}
		if name.Schema != "" {
			tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: name.Schema})
		}
		tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: name.Name})
		tn.IsRef = true

	default:
		// The ANSI spellings, unless the word is a schema qualifier.
		ansi := p.peekNext().Type != '.'
		next := p.peekNext()
		// NATIONAL names no item to take %TYPE from (PLS-00103 on Oracle
		// 23ai), as CHARACTER and DOUBLE name none.
		if ansi && next.Type == '%' && p.isKeywordStr("NATIONAL") {
			return nil, p.syntaxErrorAtTok(next)
		}
		switch {
		case ansi && p.isKeywordStr("CHARACTER"),
			ansi && p.isKeywordStr("NATIONAL") && (next.Type == kwCHAR || next.Type != tokQIDENT && next.Str == "CHARACTER"):
			if p.cur.Str == "NATIONAL" {
				tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: "NATIONAL"})
				p.advance()
			}
			tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: p.cur.Str})
			p.advance()
			if p.isKeywordStr("VARYING") {
				tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: "VARYING"})
				p.advance()
				if err := p.requireVaryingSize(); err != nil {
					return nil, err
				}
			} else if p.prev.Str == "CHARACTER" && p.typeMods != typeModsSQL && p.isKeywordStr("LARGE") {
				// [ NATIONAL ] CHARACTER LARGE OBJECT is PL/SQL's alone
				// (ORA-03062 in SQL); NATIONAL CHAR LARGE OBJECT is no type
				// (PLS-00103).
				if err := p.parseLargeObjectWords(tn); err != nil {
					return nil, err
				}
				break
			}
			if err := p.parseOptionalSizeWithSemantic(tn); err != nil {
				return nil, err
			}
			if tn.Names.Items[0].(*nodes.String).Str == "NATIONAL" {
				if err := p.checkNationalLength(tn); err != nil {
					return nil, err
				}
			}
		case ansi && p.isKeywordStr("DOUBLE"):
			p.advance() // consume DOUBLE
			if !p.isKeywordStr("PRECISION") {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance() // consume PRECISION
			tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: "DOUBLE"}, &nodes.String{Str: "PRECISION"})
		case p.isIdentLikeStr("NUMERIC"), ansi && p.isKeywordStr("DEC"):
			tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: p.cur.Str})
			p.advance()
			parseErr1121 := p.parseOptionalPrecisionScale(tn)
			if parseErr1121 != nil {
				return nil, parseErr1121
			}
		default:
			// User-defined type or %TYPE/%ROWTYPE reference
			parseErr1121 := p.parseUserDefinedType(tn)
			if parseErr1121 != nil {
				return nil, parseErr1121
			}
		}
	}

	tn.Loc.End = p.prev.End
	p.noteShadowedType(tn)
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

// checkNationalLength rejects the length semantics a national character
// type does not take: BYTE anywhere (ORA-00907 in SQL, PLS-00639 in
// PL/SQL), and CHAR in SQL (ORA-00907), which PL/SQL and a type attribute
// take, on Oracle 23ai.
func (p *Parser) checkNationalLength(tn *nodes.TypeName) error {
	for _, mod := range tn.TypeMods.Items {
		if s, ok := mod.(*nodes.String); ok && (s.Str == "BYTE" || s.Str == "CHAR" && p.typeMods == typeModsSQL) {
			return p.syntaxErrorAtType(tn)
		}
	}
	return nil
}

// parseLargeObjectWords parses the LARGE OBJECT of a PL/SQL [ NATIONAL ]
// CHARACTER LARGE OBJECT or NCHAR LARGE OBJECT, at LARGE.
func (p *Parser) parseLargeObjectWords(tn *nodes.TypeName) error {
	p.advance() // consume LARGE
	if !p.isKeywordStr("OBJECT") {
		return p.syntaxErrorAtCur()
	}
	p.advance() // consume OBJECT
	tn.Names.Items = append(tn.Names.Items, &nodes.String{Str: "LARGE"}, &nodes.String{Str: "OBJECT"})
	return nil
}

// requireVaryingSize checks that the size of a VARYING character type
// follows, as SQL requires (ORA-00906 on Oracle 23ai for CHARACTER VARYING
// alone); PL/SQL, a type attribute, and a collection element take none.
func (p *Parser) requireVaryingSize() error {
	if p.typeMods == typeModsSQL && p.cur.Type != '(' {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// predefinedTypeLead returns the name of the predefined datatype tn spells,
// "" for a name the text does not settle (a subtype, a schema-qualified
// type, t%TYPE, or a quoted name: "NUMBER" may be a user subtype of
// VARCHAR2, and Oracle 23ai compiles one as a SQL macro result). Multiword
// predefined types answer their first word: TIMESTAMP WITH [LOCAL] TIME
// ZONE, INTERVAL DAY|YEAR TO ..., LONG RAW. A spelling a TYPE or SUBTYPE
// declaration may shadow answers "" too; see noteShadowedType.
func (p *Parser) predefinedTypeLead(tn *nodes.TypeName) string {
	if p.shadowedTypes[tn] {
		return ""
	}
	return p.predefinedTypeSpelling(tn)
}

// noteShadowedType records tn, just parsed, as shadowed when it spells a
// predefined type whose name a TYPE or SUBTYPE declaration in an open
// scope has redeclared, or any predefined type in a package body, whose
// specification may redeclare one. Oracle 23ai then resolves the name to
// the user's type: after SUBTYPE sys_refcursor IS NUMBER, a function taking
// a sys_refcursor may be RESULT_CACHE, and after SUBTYPE number IS
// VARCHAR2, a number parameter takes CHARACTER SET ANY_CS. A schema-level
// type of such a name does not shadow STANDARD's.
func (p *Parser) noteShadowedType(tn *nodes.TypeName) {
	// A declaration names one word: DOUBLE PRECISION or LONG RAW is the
	// predefined type wherever it stands (PLS-00550 on Oracle 23ai for
	// DOUBLE PRECISION CHARACTER SET in a package body).
	if len(p.typeScopes) == 0 && !p.inPackageBody || tn.Names.Len() != 1 {
		return
	}
	lead := p.predefinedTypeSpelling(tn)
	if lead == "" {
		return
	}
	shadowed := p.inPackageBody
	for _, scope := range p.typeScopes {
		shadowed = shadowed || scope[lead]
	}
	if shadowed {
		if p.shadowedTypes == nil {
			p.shadowedTypes = map[*nodes.TypeName]bool{}
		}
		p.shadowedTypes[tn] = true
	}
}

// declareTypeName enters the name of a TYPE or SUBTYPE declaration in the
// innermost open declaration scope.
func (p *Parser) declareTypeName(name string) {
	if n := len(p.typeScopes); n > 0 {
		p.typeScopes[n-1][name] = true
	}
}

// predefinedTypeSpelling is predefinedTypeLead without regard to scope.
func (p *Parser) predefinedTypeSpelling(tn *nodes.TypeName) string {
	if tn == nil || tn.IsPercType || tn.IsPercRowtype || tn.IsRef || tn.Names.Len() == 0 {
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
		first == "LONG" && second == "RAW",
		first == "LONG" && second == "VARCHAR",
		first == "DOUBLE" && second == "PRECISION",
		(first == "CHAR" || first == "CHARACTER") && second == "VARYING":
		return first
	case first == "CHARACTER" && second == "LARGE":
		return "CLOB"
	case first == "NCHAR" && second == "LARGE",
		first == "NATIONAL" && tn.Names.Len() > 2 && word(2) == "LARGE":
		return "NCLOB"
	case first == "NCHAR" && second == "VARYING",
		first == "NATIONAL" && (second == "CHAR" || second == "CHARACTER"):
		return "NCHAR"
	}
	return ""
}

// parsePLSQLCharacterSet parses CHARACTER SET { character_set | item%CHARSET }
// after the datatype tn, the current token at CHARACTER. Each check below
// follows from the text alone and was confirmed on Oracle 23ai; whether a
// name is a known character set (PLS-00553) is left to the engine.
func (p *Parser) parsePLSQLCharacterSet(tn *nodes.TypeName, use plsqlCharsetUse) error {
	// Neither a %ROWTYPE record nor an object reference is a character type
	// (PLS-00550).
	if tn.IsPercRowtype || tn.IsRef {
		return p.syntaxErrorAtCur()
	}
	// Neither is a predefined non-character type, nor is the TABLE or
	// COLUMNS pseudo-type of a polymorphic table function (PLS-00550).
	if plsqlNonCharacterTypes[p.predefinedTypeLead(tn)] || p.pseudoType(tn) != "" {
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
	parts := []string{name}
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
		parts = append(parts, part)
	}
	qualified := len(parts) > 1
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
				p.qualifiedCharsets = map[*nodes.TypeName][]string{}
			}
			p.qualifiedCharsets[tn] = parts
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

// parseOptionalPrecisionScale parses optional ( precision [, scale ] ) after
// the type word last added to tn. In SQL the type decides the shape, as on
// Oracle 23ai: NUMBER, DECIMAL, and NUMERIC take a precision or * and a
// scale, which may be negative; FLOAT takes a precision or *; TIMESTAMP,
// RAW, and the leading INTERVAL field take a number; INTEGER, SMALLINT, and
// MONTH take nothing. Anything else, an empty modifier among them, fails
// there (ORA-00907, or the range errors ORA-00910, ORA-01727, ORA-01728,
// and ORA-30088). The ranges are checked as Oracle 23ai does: a NUMBER
// precision is 1 to 38 (ORA-01727) and its scale -84 to 127 (ORA-01728), a
// FLOAT precision is 1 to 126 (ORA-01724), and a fractional-second or
// leading-field precision is 0 to 9 (ORA-30088).
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
	star, scale := false, false
	var minPrecision, maxPrecision int64 = 0, -1
	if word, ok := tn.Names.Items[tn.Names.Len()-1].(*nodes.String); ok {
		switch word.Str {
		case "NUMBER", "DECIMAL", "DEC", "NUMERIC":
			star, scale = true, true
			minPrecision, maxPrecision = 1, 38
		case "FLOAT":
			star = true
			minPrecision, maxPrecision = 1, 126
		case "TIMESTAMP", "YEAR", "DAY", "SECOND":
			minPrecision, maxPrecision = 0, 9
		case "INTEGER", "SMALLINT", "MONTH":
			return p.syntaxErrorAtCur()
		}
	}
	p.advance() // consume '('

	// precision
	switch {
	case p.cur.Type == tokICONST:
		if maxPrecision >= 0 && (p.cur.Ival < minPrecision || p.cur.Ival > maxPrecision) {
			return p.syntaxErrorAtCur()
		}
		tn.TypeMods.Items = append(tn.TypeMods.Items, &nodes.Integer{Ival: p.cur.Ival})
	case star && p.cur.Type == '*':
		// NUMBER(*) or NUMBER(*,s)
		tn.TypeMods.Items = append(tn.TypeMods.Items, &nodes.String{Str: "*"})
	default:
		return p.syntaxErrorAtCur()
	}
	p.advance()

	// optional scale
	if scale && p.cur.Type == ',' {
		p.advance()
		negative := p.cur.Type == '-'
		if negative {
			p.advance()
		}
		if p.cur.Type != tokICONST {
			return p.syntaxErrorAtCur()
		}
		ival := p.cur.Ival
		if negative {
			ival = -ival
		}
		if ival < -84 || ival > 127 {
			return p.syntaxErrorAtCur()
		}
		tn.TypeMods.Items = append(tn.TypeMods.Items, &nodes.Integer{Ival: ival})
		p.advance()
	}

	if p.cur.Type != ')' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
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

	// size, which SQL requires (ORA-00910 on Oracle 23ai for VARCHAR2()
	// and VARCHAR2(BYTE))
	if p.cur.Type != tokICONST {
		return p.syntaxErrorAtCur()
	}
	tn.TypeMods.Items = append(tn.TypeMods.Items, &nodes.Integer{Ival: p.cur.Ival})
	p.advance()

	// optional BYTE or CHAR semantic
	if p.isIdentLikeStr("BYTE") {
		tn.TypeMods.Items = append(tn.TypeMods.Items, &nodes.String{Str: "BYTE"})
		p.advance()
	} else if p.cur.Type == kwCHAR {
		tn.TypeMods.Items = append(tn.TypeMods.Items, &nodes.String{Str: "CHAR"})
		p.advance()
	}

	if p.cur.Type != ')' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
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
