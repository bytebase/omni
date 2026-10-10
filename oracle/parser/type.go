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

// parsePLSQLTypeName parses the datatype of a PL/SQL declaration. Where SQL
// takes an integer literal for a length, precision, or scale, PL/SQL takes a
// static expression: VARCHAR2(ORA_MAX_NAME_LEN + 2), NUMBER(pkg.p, 2). That
// the expression is static is checked when the unit compiles (PLS-00491 on
// Oracle 23ai), not when it parses, so any expression parses here.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/plsql-language-fundamentals.html (Static Expressions)
func (p *Parser) parsePLSQLTypeName() (*nodes.TypeName, error) {
	saved := p.plsqlTypeMods
	p.plsqlTypeMods = true
	tn, err := p.parseTypeName()
	p.plsqlTypeMods = saved
	return tn, err
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
	p.plsqlTypeMods = false
	expr, err := p.parseExpr()
	p.plsqlTypeMods = true
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
	if p.plsqlTypeMods {
		return p.parsePLSQLTypeMods(tn, false)
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
	if p.plsqlTypeMods {
		return p.parsePLSQLTypeMods(tn, true)
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
