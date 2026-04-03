// Package parser implements a recursive descent SQL parser for T-SQL (SQL Server).
package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// Token type constants for the T-SQL lexer.
const (
	tokEOF = 0

	// Literal tokens (non-keyword, non-single-char)
	tokICONST      = iota + 256 // integer constant
	tokFCONST                   // floating-point constant
	tokSCONST                   // string constant
	tokNSCONST                  // N'...' nvarchar string constant
	tokIDENT                    // identifier (regular or [bracketed])
	tokVARIABLE                 // @variable
	tokSYSVARIABLE              // @@system_variable

	// Multi-character operators
	tokNOTEQUAL   // <> or !=
	tokLESSEQUAL  // <=
	tokGREATEQUAL // >=
	tokNOTLESS    // !<
	tokNOTGREATER // !>
	tokCOLONCOLON // ::
	tokPLUSEQUAL  // +=
	tokMINUSEQUAL // -=
	tokMULEQUAL   // *=
	tokDIVEQUAL   // /=
	tokMODEQUAL   // %=
	tokANDEQUAL   // &=
	tokOREQUAL    // |=
	tokXOREQUAL   // ^=

	// Keywords start here. T-SQL keywords (case-insensitive).
	kwADD
	kwALL
	kwALTER
	kwAND
	kwANY
	kwAPPLY
	kwAS
	kwASC
	kwAT
	kwAUTHORIZATION
	kwBACKUP
	kwBEGIN
	kwBETWEEN
	kwBREAK
	kwBROWSE
	kwBULK
	kwBY
	kwCASCADE
	kwCASE
	kwCAST
	kwCATCH
	kwCHECK
	kwCHECKPOINT
	kwCLOSE
	kwCLUSTERED
	kwCOALESCE
	kwCOLLATE
	kwCOLUMN
	kwCOLUMNSTORE
	kwCOMMIT
	kwCOMPUTE
	kwCONSTRAINT
	kwCONTAINS
	kwCONTAINSTABLE
	kwCONTINUE
	kwCONVERT
	kwCREATE
	kwCROSS
	kwCURRENT
	kwCURRENT_DATE
	kwCURRENT_TIME
	kwCURRENT_TIMESTAMP
	kwCURRENT_USER
	kwCURSOR
	kwDATABASE
	kwDBCC
	kwDEALLOCATE
	kwDECLARE
	kwDEFAULT
	kwDELAY
	kwDELETE
	kwDENY
	kwDESC
	kwDISTINCT
	kwDISTRIBUTED
	kwDO
	kwDROP
	kwDUMP
	kwELSE
	kwEND
	kwERRLVL
	kwESCAPE
	kwEXCEPT
	kwEXEC
	kwEXECUTE
	kwEXISTS
	kwEXIT
	kwEXTERNAL
	kwFETCH
	kwFILE
	kwFILLFACTOR
	kwFOLLOWING
	kwFOR
	kwFOREIGN
	kwFREETEXT
	kwFREETEXTTABLE
	kwFROM
	kwFULL
	kwFUNCTION
	kwGO
	kwGOTO
	kwGRANT
	kwGROUP
	kwGROUPS
	kwHAVING
	kwHOLDLOCK
	kwIDENTITY
	kwIDENTITY_INSERT
	kwIDENTITYCOL
	kwIF
	kwIIF
	kwIN
	kwINCLUDE
	kwINDEX
	kwINNER
	kwINSERT
	kwINTERSECT
	kwINTO
	kwIS
	kwJOIN
	kwJSON
	kwKEY
	kwKILL
	kwLEFT
	kwLIKE
	kwLINENO
	kwLOAD
	kwLOGIN
	kwMATERIALIZED
	kwMAX
	kwMERGE
	kwNATIONAL
	kwNOCHECK
	kwNOCOUNT
	kwNOLOCK
	kwNONCLUSTERED
	kwNOT
	kwNOWAIT
	kwNULL
	kwNULLIF
	kwOF
	kwOFF
	kwOFFSET
	kwOFFSETS
	kwON
	kwOPEN
	kwOPENDATASOURCE
	kwOPENJSON
	kwOPENQUERY
	kwOPENROWSET
	kwOPENXML
	kwOPTION
	kwOR
	kwORDER
	kwOUTER
	kwOUTPUT
	kwOVER
	kwPARTITION
	kwPATH
	kwPRECEDING
	kwPERCENT
	kwPIVOT
	kwPLAN
	kwPRECISION
	kwPRIMARY
	kwPRINT
	kwPROC
	kwPROCEDURE
	kwPUBLIC
	kwRAISERROR
	kwRANGE
	kwRAW
	kwREAD
	kwREADONLY
	kwREADTEXT
	kwRECONFIGURE
	kwREFERENCES
	kwREPLICATION
	kwRESTORE
	kwRESTRICT
	kwRETURN
	kwRETURNS
	kwREVOKE
	kwRIGHT
	kwROLE
	kwROLLBACK
	kwROWCOUNT
	kwROWGUIDCOL
	kwROWS
	kwRULE
	kwSAVE
	kwSCHEMA
	kwSCHEMABINDING
	kwSECURITYAUDIT
	kwSELECT
	kwSEMIJOIN
	kwSESSION_USER
	kwSET
	kwSETUSER
	kwSHUTDOWN
	kwSOME
	kwSTATISTICS
	kwSYSTEM_USER
	kwTABLE
	kwTABLESAMPLE
	kwTEXTSIZE
	kwTHEN
	kwTHROW
	kwTIME
	kwTO
	kwTOP
	kwTRAN
	kwTRANSACTION
	kwTRIGGER
	kwTRUNCATE
	kwTRY
	kwTRY_CAST
	kwTRY_CONVERT
	kwTYPE
	kwUNBOUNDED
	kwUNION
	kwUNIQUE
	kwUNPIVOT
	kwUPDATE
	kwUPDATETEXT
	kwUSE
	kwUSER
	kwVALUES
	kwVARYING
	kwVECTOR
	kwVIEW
	kwWAITFOR
	kwWHEN
	kwWHERE
	kwWHILE
	kwWITH
	kwWITHIN
	kwWRITETEXT
	kwXML
	kwXACT_ABORT
	kwZONE
)

// keywordMap maps lowercase keyword strings to token types.
var keywordMap map[string]int

func init() {
	keywordMap = map[string]int{
		"add": kwADD, "all": kwALL, "alter": kwALTER, "and": kwAND, "any": kwANY,
		"apply": kwAPPLY, "as": kwAS, "asc": kwASC, "at": kwAT, "authorization": kwAUTHORIZATION,
		"backup": kwBACKUP, "begin": kwBEGIN, "between": kwBETWEEN, "break": kwBREAK,
		"browse": kwBROWSE, "bulk": kwBULK, "by": kwBY,
		"cascade": kwCASCADE, "case": kwCASE, "cast": kwCAST, "catch": kwCATCH,
		"check": kwCHECK, "checkpoint": kwCHECKPOINT, "close": kwCLOSE,
		"clustered": kwCLUSTERED, "coalesce": kwCOALESCE, "collate": kwCOLLATE,
		"column": kwCOLUMN, "columnstore": kwCOLUMNSTORE, "commit": kwCOMMIT,
		"compute": kwCOMPUTE, "constraint": kwCONSTRAINT, "contains": kwCONTAINS,
		"containstable": kwCONTAINSTABLE, "continue": kwCONTINUE, "convert": kwCONVERT,
		"create": kwCREATE, "cross": kwCROSS, "current": kwCURRENT,
		"current_date": kwCURRENT_DATE, "current_time": kwCURRENT_TIME,
		"current_timestamp": kwCURRENT_TIMESTAMP, "current_user": kwCURRENT_USER,
		"cursor":   kwCURSOR,
		"database": kwDATABASE, "dbcc": kwDBCC, "deallocate": kwDEALLOCATE,
		"declare": kwDECLARE, "default": kwDEFAULT, "delay": kwDELAY,
		"delete": kwDELETE, "deny": kwDENY, "desc": kwDESC, "distinct": kwDISTINCT,
		"distributed": kwDISTRIBUTED, "do": kwDO, "drop": kwDROP, "dump": kwDUMP,
		"else": kwELSE, "end": kwEND, "errlvl": kwERRLVL, "escape": kwESCAPE,
		"except": kwEXCEPT, "exec": kwEXEC, "execute": kwEXECUTE, "exists": kwEXISTS,
		"exit": kwEXIT, "external": kwEXTERNAL,
		"fetch": kwFETCH, "file": kwFILE, "fillfactor": kwFILLFACTOR,
		"following": kwFOLLOWING, "for": kwFOR,
		"foreign": kwFOREIGN, "freetext": kwFREETEXT, "freetexttable": kwFREETEXTTABLE,
		"from": kwFROM, "full": kwFULL, "function": kwFUNCTION,
		"go": kwGO, "goto": kwGOTO, "grant": kwGRANT, "group": kwGROUP, "groups": kwGROUPS,
		"having": kwHAVING, "holdlock": kwHOLDLOCK,
		"identity": kwIDENTITY, "identity_insert": kwIDENTITY_INSERT,
		"identitycol": kwIDENTITYCOL, "if": kwIF, "iif": kwIIF, "in": kwIN,
		"include": kwINCLUDE, "index": kwINDEX, "inner": kwINNER, "insert": kwINSERT,
		"intersect": kwINTERSECT, "into": kwINTO, "is": kwIS,
		"join": kwJOIN, "json": kwJSON,
		"key": kwKEY, "kill": kwKILL,
		"left": kwLEFT, "like": kwLIKE, "lineno": kwLINENO, "load": kwLOAD,
		"login": kwLOGIN,
		"materialized": kwMATERIALIZED, "max": kwMAX, "merge": kwMERGE,
		"national": kwNATIONAL, "nocheck": kwNOCHECK, "nocount": kwNOCOUNT,
		"nolock": kwNOLOCK, "nonclustered": kwNONCLUSTERED, "not": kwNOT,
		"nowait": kwNOWAIT, "null": kwNULL, "nullif": kwNULLIF,
		"of": kwOF, "off": kwOFF, "offset": kwOFFSET, "offsets": kwOFFSETS,
		"on": kwON, "open": kwOPEN, "opendatasource": kwOPENDATASOURCE,
		"openjson": kwOPENJSON, "openquery": kwOPENQUERY, "openrowset": kwOPENROWSET,
		"openxml": kwOPENXML, "option": kwOPTION, "or": kwOR, "order": kwORDER,
		"outer": kwOUTER, "output": kwOUTPUT, "over": kwOVER,
		"partition": kwPARTITION, "path": kwPATH, "preceding": kwPRECEDING, "percent": kwPERCENT,
		"pivot": kwPIVOT, "plan": kwPLAN, "precision": kwPRECISION,
		"primary": kwPRIMARY, "print": kwPRINT, "proc": kwPROC,
		"procedure": kwPROCEDURE, "public": kwPUBLIC,
		"raiserror": kwRAISERROR, "range": kwRANGE, "raw": kwRAW, "read": kwREAD,
		"readonly": kwREADONLY, "readtext": kwREADTEXT, "reconfigure": kwRECONFIGURE,
		"references": kwREFERENCES, "replication": kwREPLICATION, "restore": kwRESTORE,
		"restrict": kwRESTRICT, "return": kwRETURN, "returns": kwRETURNS,
		"revoke": kwREVOKE, "right": kwRIGHT, "role": kwROLE, "rollback": kwROLLBACK,
		"rowcount": kwROWCOUNT, "rowguidcol": kwROWGUIDCOL, "rows": kwROWS, "rule": kwRULE,
		"save": kwSAVE, "schema": kwSCHEMA, "schemabinding": kwSCHEMABINDING,
		"securityaudit": kwSECURITYAUDIT, "select": kwSELECT, "semijoin": kwSEMIJOIN,
		"session_user": kwSESSION_USER, "set": kwSET, "setuser": kwSETUSER,
		"shutdown": kwSHUTDOWN, "some": kwSOME, "statistics": kwSTATISTICS,
		"system_user": kwSYSTEM_USER,
		"table":       kwTABLE, "tablesample": kwTABLESAMPLE, "textsize": kwTEXTSIZE,
		"then": kwTHEN, "throw": kwTHROW, "time": kwTIME, "to": kwTO, "top": kwTOP,
		"tran": kwTRAN, "transaction": kwTRANSACTION, "trigger": kwTRIGGER,
		"truncate": kwTRUNCATE, "try": kwTRY, "try_cast": kwTRY_CAST,
		"try_convert": kwTRY_CONVERT, "type": kwTYPE,
		"unbounded": kwUNBOUNDED, "union": kwUNION, "unique": kwUNIQUE, "unpivot": kwUNPIVOT,
		"update": kwUPDATE, "updatetext": kwUPDATETEXT, "use": kwUSE, "user": kwUSER,
		"values": kwVALUES, "varying": kwVARYING, "vector": kwVECTOR, "view": kwVIEW,
		"waitfor": kwWAITFOR, "when": kwWHEN, "where": kwWHERE, "while": kwWHILE,
		"with": kwWITH, "within": kwWITHIN, "writetext": kwWRITETEXT,
		"xml": kwXML, "xact_abort": kwXACT_ABORT,
		"zone": kwZONE,
	}
}

// lookupKeyword returns the token type for a keyword (case-insensitive),
// or tokIDENT if not a keyword.
func lookupKeyword(ident string) int {
	if tok, ok := keywordMap[strings.ToLower(ident)]; ok {
		return tok
	}
	return tokIDENT
}

// Token represents a lexical token.
type Token struct {
	Type int    // Token type
	Str  string // String value for identifiers, string literals, operators
	Ival int64  // Integer value for tokICONST
	Loc  int    // Byte offset in the source text
	End  int    // Exclusive byte offset past token end
}

// Lexer implements a T-SQL lexer.
type Lexer struct {
	input string
	pos   int
	start int
	Err   error
}

// NewLexer creates a new T-SQL lexer.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	tok := l.nextTokenInner()
	tok.End = l.pos
	return tok
}

// nextTokenInner performs the actual lexing. NextToken wraps this to set tok.End.
func (l *Lexer) nextTokenInner() Token {
	l.skipWhitespace()
	if l.pos >= len(l.input) {
		return Token{Type: tokEOF, Loc: l.pos}
	}

	l.start = l.pos
	ch := l.input[l.pos]

	// Line comment: --
	if ch == '-' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '-' {
		l.skipLineComment()
		return l.nextTokenInner()
	}

	// Block comment: /* ... */
	if ch == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '*' {
		l.skipBlockComment()
		return l.nextTokenInner()
	}

	// N'...' nvarchar string literal
	if (ch == 'N' || ch == 'n') && l.pos+1 < len(l.input) && l.input[l.pos+1] == '\'' {
		l.pos++ // skip N
		return l.lexNString()
	}

	// '...' string literal
	if ch == '\'' {
		return l.lexString()
	}

	// [bracketed identifier]
	if ch == '[' {
		return l.lexBracketedIdent()
	}

	// "quoted identifier"
	if ch == '"' {
		return l.lexQuotedIdent()
	}

	// @variable or @@sysvariable
	if ch == '@' {
		return l.lexVariable()
	}

	// Two-character operators (check before single-char)
	if l.pos+1 < len(l.input) {
		ch2 := l.input[l.pos : l.pos+2]
		switch ch2 {
		case "<>":
			l.pos += 2
			return Token{Type: tokNOTEQUAL, Str: "<>", Loc: l.start}
		case "!=":
			l.pos += 2
			return Token{Type: tokNOTEQUAL, Str: "!=", Loc: l.start}
		case "<=":
			l.pos += 2
			return Token{Type: tokLESSEQUAL, Str: "<=", Loc: l.start}
		case ">=":
			l.pos += 2
			return Token{Type: tokGREATEQUAL, Str: ">=", Loc: l.start}
		case "!<":
			l.pos += 2
			return Token{Type: tokNOTLESS, Str: "!<", Loc: l.start}
		case "!>":
			l.pos += 2
			return Token{Type: tokNOTGREATER, Str: "!>", Loc: l.start}
		case "::":
			l.pos += 2
			return Token{Type: tokCOLONCOLON, Str: "::", Loc: l.start}
		case "+=":
			l.pos += 2
			return Token{Type: tokPLUSEQUAL, Str: "+=", Loc: l.start}
		case "-=":
			l.pos += 2
			return Token{Type: tokMINUSEQUAL, Str: "-=", Loc: l.start}
		case "*=":
			l.pos += 2
			return Token{Type: tokMULEQUAL, Str: "*=", Loc: l.start}
		case "/=":
			l.pos += 2
			return Token{Type: tokDIVEQUAL, Str: "/=", Loc: l.start}
		case "%=":
			l.pos += 2
			return Token{Type: tokMODEQUAL, Str: "%=", Loc: l.start}
		case "&=":
			l.pos += 2
			return Token{Type: tokANDEQUAL, Str: "&=", Loc: l.start}
		case "|=":
			l.pos += 2
			return Token{Type: tokOREQUAL, Str: "|=", Loc: l.start}
		case "^=":
			l.pos += 2
			return Token{Type: tokXOREQUAL, Str: "^=", Loc: l.start}
		}
	}

	// Numbers (digit or .digit)
	if isDigit(ch) || (ch == '.' && l.pos+1 < len(l.input) && isDigit(l.input[l.pos+1])) {
		return l.lexNumber()
	}

	// Single-character tokens
	if isSingleChar(ch) {
		l.pos++
		return Token{Type: int(ch), Str: string(ch), Loc: l.start}
	}

	// Identifiers and keywords
	if isIdentStart(ch) {
		return l.lexIdent()
	}

	// Unknown character - return as itself
	l.pos++
	return Token{Type: int(ch), Str: string(ch), Loc: l.start}
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' || ch == '\f' {
			l.pos++
		} else {
			break
		}
	}
}

func (l *Lexer) skipLineComment() {
	l.pos += 2
	for l.pos < len(l.input) {
		if l.input[l.pos] == '\n' {
			l.pos++
			return
		}
		l.pos++
	}
}

func (l *Lexer) skipBlockComment() {
	l.pos += 2
	depth := 1
	for l.pos < len(l.input) && depth > 0 {
		if l.input[l.pos] == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '*' {
			depth++
			l.pos += 2
		} else if l.input[l.pos] == '*' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '/' {
			depth--
			l.pos += 2
		} else {
			l.pos++
		}
	}
	if depth > 0 {
		l.Err = fmt.Errorf("unterminated block comment")
	}
}

func (l *Lexer) lexString() Token {
	l.pos++ // skip opening quote
	var buf strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '\'' {
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '\'' {
				buf.WriteByte('\'')
				l.pos += 2
				continue
			}
			l.pos++
			return Token{Type: tokSCONST, Str: buf.String(), Loc: l.start}
		}
		buf.WriteByte(ch)
		l.pos++
	}
	l.Err = fmt.Errorf("unterminated string literal")
	return Token{Type: tokEOF, Loc: l.start}
}

func (l *Lexer) lexNString() Token {
	tok := l.lexString()
	if tok.Type == tokSCONST {
		tok.Type = tokNSCONST
	}
	return tok
}

func (l *Lexer) lexBracketedIdent() Token {
	l.pos++ // skip [
	var buf strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == ']' {
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == ']' {
				buf.WriteByte(']')
				l.pos += 2
				continue
			}
			l.pos++
			return Token{Type: tokIDENT, Str: buf.String(), Loc: l.start}
		}
		buf.WriteByte(ch)
		l.pos++
	}
	l.Err = fmt.Errorf("unterminated bracketed identifier")
	return Token{Type: tokEOF, Loc: l.start}
}

func (l *Lexer) lexQuotedIdent() Token {
	l.pos++ // skip "
	var buf strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '"' {
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '"' {
				buf.WriteByte('"')
				l.pos += 2
				continue
			}
			l.pos++
			str := buf.String()
			if str == "" {
				l.Err = fmt.Errorf("zero-length delimited identifier")
				return Token{Type: tokEOF, Loc: l.start}
			}
			return Token{Type: tokIDENT, Str: str, Loc: l.start}
		}
		buf.WriteByte(ch)
		l.pos++
	}
	l.Err = fmt.Errorf("unterminated quoted identifier")
	return Token{Type: tokEOF, Loc: l.start}
}

func (l *Lexer) lexVariable() Token {
	l.pos++ // skip first @
	if l.pos < len(l.input) && l.input[l.pos] == '@' {
		// @@sysvariable
		l.pos++
		start := l.pos
		for l.pos < len(l.input) && isIdentCont(l.input[l.pos]) {
			l.pos++
		}
		name := "@@" + l.input[start:l.pos]
		return Token{Type: tokSYSVARIABLE, Str: name, Loc: l.start}
	}
	// @variable
	start := l.pos
	for l.pos < len(l.input) && isIdentCont(l.input[l.pos]) {
		l.pos++
	}
	name := "@" + l.input[start:l.pos]
	return Token{Type: tokVARIABLE, Str: name, Loc: l.start}
}

func (l *Lexer) lexNumber() Token {
	start := l.pos
	isFloat := false

	// Handle hex: 0x...
	if l.input[l.pos] == '0' && l.pos+1 < len(l.input) && (l.input[l.pos+1] == 'x' || l.input[l.pos+1] == 'X') {
		l.pos += 2
		for l.pos < len(l.input) && isHexDigit(l.input[l.pos]) {
			l.pos++
		}
		numStr := l.input[start:l.pos]
		val, err := strconv.ParseInt(numStr[2:], 16, 64)
		if err != nil {
			return Token{Type: tokFCONST, Str: numStr, Loc: l.start}
		}
		return Token{Type: tokICONST, Ival: val, Str: numStr, Loc: l.start}
	}

	// Leading dot: .123
	if l.input[l.pos] == '.' {
		isFloat = true
		l.pos++
	}

	for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
		l.pos++
	}

	// Decimal point
	if !isFloat && l.pos < len(l.input) && l.input[l.pos] == '.' {
		isFloat = true
		l.pos++
		for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
			l.pos++
		}
	}

	// Exponent
	if l.pos < len(l.input) && (l.input[l.pos] == 'e' || l.input[l.pos] == 'E') {
		isFloat = true
		l.pos++
		if l.pos < len(l.input) && (l.input[l.pos] == '+' || l.input[l.pos] == '-') {
			l.pos++
		}
		for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
			l.pos++
		}
	}

	numStr := l.input[start:l.pos]
	if isFloat {
		return Token{Type: tokFCONST, Str: numStr, Loc: l.start}
	}

	val, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return Token{Type: tokFCONST, Str: numStr, Loc: l.start}
	}
	return Token{Type: tokICONST, Ival: val, Str: numStr, Loc: l.start}
}

func (l *Lexer) lexIdent() Token {
	start := l.pos
	for l.pos < len(l.input) && isIdentCont(l.input[l.pos]) {
		l.pos++
	}
	ident := l.input[start:l.pos]

	// Check for multi-word keywords like TRY_CAST, TRY_CONVERT, etc.
	tok := lookupKeyword(ident)
	if tok != tokIDENT {
		return Token{Type: tok, Str: ident, Loc: l.start}
	}

	return Token{Type: tokIDENT, Str: ident, Loc: l.start}
}

// Character classification functions

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isHexDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

func isIdentStart(ch byte) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || ch == '_' || ch == '#' || ch >= 128
}

func isIdentCont(ch byte) bool {
	return isIdentStart(ch) || (ch >= '0' && ch <= '9') || ch == '$'
}

func isSingleChar(ch byte) bool {
	return ch == '(' || ch == ')' || ch == ',' || ch == ';' || ch == '.' ||
		ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '%' ||
		ch == '=' || ch == '<' || ch == '>' || ch == '~' || ch == '&' ||
		ch == '|' || ch == '^' || ch == ':' || ch == '!'
}
