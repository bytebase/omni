// Package parser implements an Oracle PL/SQL lexer and recursive descent parser.
package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// Token types returned by the lexer.
const (
	tokEOF = 0
)

// Non-keyword token types.
const (
	tokICONST = iota + 1000
	tokFCONST
	tokSCONST     // single-quoted string literal
	tokNCHARLIT   // N'...' national character literal
	tokIDENT      // unquoted identifier (lowercased)
	tokQIDENT     // "double-quoted" identifier (case-preserved)
	tokBIND       // :name or :1 bind variable
	tokCONCAT     // ||
	tokASSIGN     // :=
	tokASSOC      // =>
	tokDOTDOT     // ..
	tokEXPON      // **
	tokLESSEQ     // <=
	tokGREATEQ    // >=
	tokNOTEQ      // != or <> or ~= or ^=
	tokLABELOPEN  // <<
	tokLABELCLOSE // >>
	tokOP         // generic operator
	tokHINT       // /*+ hint text */
)

// Token represents a lexical token.
type Token struct {
	Type int    // token type
	Str  string // string value
	Ival int64  // integer value
	Loc  int    // byte offset in source
}

// Lexer implements an Oracle SQL/PL/SQL lexer.
type Lexer struct {
	input string
	pos   int
	start int

	// Error
	Err error
}

// NewLexer creates a new lexer for the given input.
func NewLexer(input string) *Lexer {
	return &Lexer{
		input: input,
		pos:   0,
	}
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	if l.pos >= len(l.input) {
		return Token{Type: tokEOF, Loc: l.pos}
	}

	l.start = l.pos
	return l.lexInitial()
}

// lexInitial handles tokens in the initial state.
func (l *Lexer) lexInitial() Token {
	l.skipWhitespace()
	if l.pos >= len(l.input) {
		return Token{Type: tokEOF, Loc: l.pos}
	}

	l.start = l.pos
	ch := l.input[l.pos]

	// Line comments: --
	if ch == '-' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '-' {
		l.skipLineComment()
		return l.NextToken()
	}

	// Block comments: /* */ and hints: /*+ */
	if ch == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '*' {
		return l.lexBlockCommentOrHint()
	}

	// National character literal: N'...'
	if (ch == 'N' || ch == 'n') && l.pos+1 < len(l.input) && l.input[l.pos+1] == '\'' {
		l.pos++ // skip N
		str := l.lexSingleQuotedString()
		return Token{Type: tokNCHARLIT, Str: str, Loc: l.start}
	}

	// Q-quote mechanism: q'[delim]...[delim]'
	if (ch == 'Q' || ch == 'q') && l.pos+1 < len(l.input) && l.input[l.pos+1] == '\'' {
		return l.lexQQuote()
	}

	// Single-quoted string: '...'
	if ch == '\'' {
		str := l.lexSingleQuotedString()
		return Token{Type: tokSCONST, Str: str, Loc: l.start}
	}

	// Double-quoted identifier: "..."
	if ch == '"' {
		return l.lexDoubleQuotedIdent()
	}

	// Bind variable: :name or :1
	if ch == ':' {
		if l.pos+1 < len(l.input) {
			next := l.input[l.pos+1]
			// := assignment
			if next == '=' {
				l.pos += 2
				return Token{Type: tokASSIGN, Str: ":=", Loc: l.start}
			}
			// Bind variable
			if isIdentStart(next) || isDigit(next) {
				l.pos++ // skip :
				return l.lexBindVariable()
			}
		}
		// Single : is returned as itself
		l.pos++
		return Token{Type: int(ch), Str: string(ch), Loc: l.start}
	}

	// Two-character tokens
	if l.pos+1 < len(l.input) {
		ch2 := l.input[l.pos : l.pos+2]
		switch ch2 {
		case "||":
			l.pos += 2
			return Token{Type: tokCONCAT, Str: "||", Loc: l.start}
		case "=>":
			l.pos += 2
			return Token{Type: tokASSOC, Str: "=>", Loc: l.start}
		case "..":
			l.pos += 2
			return Token{Type: tokDOTDOT, Str: "..", Loc: l.start}
		case "**":
			l.pos += 2
			return Token{Type: tokEXPON, Str: "**", Loc: l.start}
		case "<=":
			l.pos += 2
			return Token{Type: tokLESSEQ, Str: "<=", Loc: l.start}
		case ">=":
			l.pos += 2
			return Token{Type: tokGREATEQ, Str: ">=", Loc: l.start}
		case "!=":
			l.pos += 2
			return Token{Type: tokNOTEQ, Str: "!=", Loc: l.start}
		case "<>":
			l.pos += 2
			return Token{Type: tokNOTEQ, Str: "<>", Loc: l.start}
		case "~=":
			l.pos += 2
			return Token{Type: tokNOTEQ, Str: "~=", Loc: l.start}
		case "^=":
			l.pos += 2
			return Token{Type: tokNOTEQ, Str: "^=", Loc: l.start}
		case "<<":
			l.pos += 2
			return Token{Type: tokLABELOPEN, Str: "<<", Loc: l.start}
		case ">>":
			l.pos += 2
			return Token{Type: tokLABELCLOSE, Str: ">>", Loc: l.start}
		}
	}

	// Numbers
	if isDigit(ch) || (ch == '.' && l.pos+1 < len(l.input) && isDigit(l.input[l.pos+1])) {
		return l.lexNumber()
	}

	// Identifiers and keywords
	if isIdentStart(ch) {
		return l.lexIdentOrKeyword()
	}

	// Self-delimiting single characters
	l.pos++
	return Token{Type: int(ch), Str: string(ch), Loc: l.start}
}

// skipWhitespace skips whitespace characters.
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

// skipLineComment skips a -- comment to end of line.
func (l *Lexer) skipLineComment() {
	l.pos += 2 // skip --
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		l.pos++
		if ch == '\n' {
			break
		}
	}
}

// lexBlockCommentOrHint handles /* ... */ comments and /*+ hints */.
func (l *Lexer) lexBlockCommentOrHint() Token {
	l.pos += 2 // skip /*

	// Check if it's a hint: /*+ ... */
	isHint := l.pos < len(l.input) && l.input[l.pos] == '+'
	if isHint {
		l.pos++ // skip +
	}

	var buf strings.Builder
	depth := 1

	for l.pos < len(l.input) && depth > 0 {
		if l.input[l.pos] == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '*' {
			depth++
			buf.WriteString("/*")
			l.pos += 2
			continue
		}
		if l.input[l.pos] == '*' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '/' {
			depth--
			if depth > 0 {
				buf.WriteString("*/")
			}
			l.pos += 2
			continue
		}
		buf.WriteByte(l.input[l.pos])
		l.pos++
	}

	if depth > 0 {
		l.Err = fmt.Errorf("unterminated block comment")
		return Token{Type: tokEOF, Loc: l.start}
	}

	if isHint {
		return Token{Type: tokHINT, Str: strings.TrimSpace(buf.String()), Loc: l.start}
	}

	// Regular comment -- skip and get next token
	return l.NextToken()
}

// lexSingleQuotedString handles '...' strings with ” escape.
func (l *Lexer) lexSingleQuotedString() string {
	l.pos++ // skip opening '
	var buf strings.Builder

	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '\'' {
			// Check for doubled quote ''
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '\'' {
				buf.WriteByte('\'')
				l.pos += 2
				continue
			}
			// End of string
			l.pos++
			return buf.String()
		}
		buf.WriteByte(ch)
		l.pos++
	}

	l.Err = fmt.Errorf("unterminated string literal")
	return buf.String()
}

// lexQQuote handles Oracle's alternative quoting mechanism: q'[delim]...[delim]'
func (l *Lexer) lexQQuote() Token {
	l.pos += 2 // skip q'

	if l.pos >= len(l.input) {
		l.Err = fmt.Errorf("unterminated q-quote string")
		return Token{Type: tokEOF, Loc: l.start}
	}

	openDelim := l.input[l.pos]
	l.pos++

	// Determine close delimiter
	var closeDelim byte
	switch openDelim {
	case '[':
		closeDelim = ']'
	case '{':
		closeDelim = '}'
	case '(':
		closeDelim = ')'
	case '<':
		closeDelim = '>'
	default:
		closeDelim = openDelim
	}

	var buf strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == closeDelim && l.pos+1 < len(l.input) && l.input[l.pos+1] == '\'' {
			l.pos += 2 // skip close delim and '
			return Token{Type: tokSCONST, Str: buf.String(), Loc: l.start}
		}
		buf.WriteByte(ch)
		l.pos++
	}

	l.Err = fmt.Errorf("unterminated q-quote string")
	return Token{Type: tokEOF, Loc: l.start}
}

// lexDoubleQuotedIdent handles "..." identifiers.
func (l *Lexer) lexDoubleQuotedIdent() Token {
	l.pos++ // skip opening "
	var buf strings.Builder

	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '"' {
			// Check for "" (escaped quote)
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '"' {
				buf.WriteByte('"')
				l.pos += 2
				continue
			}
			l.pos++ // skip closing "
			str := buf.String()
			if str == "" {
				l.Err = fmt.Errorf("zero-length delimited identifier")
				return Token{Type: tokEOF, Loc: l.start}
			}
			return Token{Type: tokQIDENT, Str: str, Loc: l.start}
		}
		buf.WriteByte(ch)
		l.pos++
	}

	l.Err = fmt.Errorf("unterminated quoted identifier")
	return Token{Type: tokEOF, Loc: l.start}
}

// lexBindVariable handles :name or :1 bind variables.
func (l *Lexer) lexBindVariable() Token {
	start := l.pos

	if isDigit(l.input[l.pos]) {
		// Numeric bind: :1, :2, etc.
		for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
			l.pos++
		}
	} else {
		// Named bind: :name
		for l.pos < len(l.input) && isIdentCont(l.input[l.pos]) {
			l.pos++
		}
	}

	name := l.input[start:l.pos]
	return Token{Type: tokBIND, Str: name, Loc: l.start}
}

// lexNumber handles integer and floating-point literals.
func (l *Lexer) lexNumber() Token {
	start := l.pos
	isFloat := false

	// Integer part (may be empty for .5)
	for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
		l.pos++
	}

	// Decimal point
	if l.pos < len(l.input) && l.input[l.pos] == '.' {
		// Check for .. (range operator)
		if l.pos+1 < len(l.input) && l.input[l.pos+1] == '.' {
			goto done
		}
		l.pos++
		isFloat = true
		for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
			l.pos++
		}
	}

	// Exponent
	if l.pos < len(l.input) && (l.input[l.pos] == 'e' || l.input[l.pos] == 'E') {
		l.pos++
		if l.pos < len(l.input) && (l.input[l.pos] == '+' || l.input[l.pos] == '-') {
			l.pos++
		}
		if l.pos >= len(l.input) || !isDigit(l.input[l.pos]) {
			l.Err = fmt.Errorf("invalid numeric literal")
			return Token{Type: tokEOF, Loc: l.start}
		}
		for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
			l.pos++
		}
		isFloat = true
	}

	// Oracle allows d/D suffix for double, f/F for float (in some contexts)
	if l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == 'f' || ch == 'F' || ch == 'd' || ch == 'D' {
			l.pos++
			isFloat = true
		}
	}

done:
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

// lexIdentOrKeyword handles identifiers and keywords.
func (l *Lexer) lexIdentOrKeyword() Token {
	start := l.pos

	for l.pos < len(l.input) && isIdentCont(l.input[l.pos]) {
		l.pos++
	}

	// Check for @dblink suffix on identifiers is handled at parser level
	ident := l.input[start:l.pos]

	// Check if it's a keyword (case-insensitive)
	upper := strings.ToUpper(ident)
	if kw, ok := oracleKeywords[upper]; ok {
		return Token{Type: kw, Str: upper, Loc: l.start}
	}

	// Return as identifier (Oracle unquoted identifiers are case-insensitive;
	// we store them uppercased per Oracle convention)
	return Token{Type: tokIDENT, Str: upper, Loc: l.start}
}

// Character classification functions.

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isIdentStart(ch byte) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || ch == '_' || ch >= 128
}

func isIdentCont(ch byte) bool {
	return isIdentStart(ch) || isDigit(ch) || ch == '$' || ch == '#'
}

// oracleKeywords maps uppercase keyword strings to token types.
// Oracle keywords. This is a subset; more can be added as batches are implemented.
var oracleKeywords = map[string]int{
	// Reserved words
	"ACCESS":       kwACCESS,
	"ADD":          kwADD,
	"ALL":          kwALL,
	"ALTER":        kwALTER,
	"AND":          kwAND,
	"ANY":          kwANY,
	"AS":           kwAS,
	"ASC":          kwASC,
	"AUDIT":        kwAUDIT,
	"BETWEEN":      kwBETWEEN,
	"BY":           kwBY,
	"CASE":         kwCASE,
	"CAST":         kwCAST,
	"CHECK":        kwCHECK,
	"CLUSTER":      kwCLUSTER,
	"COLUMN":       kwCOLUMN,
	"COMMENT":      kwCOMMENT,
	"COMMIT":       kwCOMMIT,
	"COMPRESS":     kwCOMPRESS,
	"CONNECT":      kwCONNECT,
	"CONSTRAINT":   kwCONSTRAINT,
	"CREATE":       kwCREATE,
	"CROSS":        kwCROSS,
	"CURRENT":      kwCURRENT,
	"DATE":         kwDATE,
	"DECIMAL":      kwDECIMAL,
	"DECODE":       kwDECODE,
	"DEFAULT":      kwDEFAULT,
	"DELETE":       kwDELETE,
	"DESC":         kwDESC,
	"DISTINCT":     kwDISTINCT,
	"DROP":         kwDROP,
	"ELSE":         kwELSE,
	"ELSIF":        kwELSIF,
	"END":          kwEND,
	"EXCEPTION":    kwEXCEPTION,
	"EXCLUSIVE":    kwEXCLUSIVE,
	"EXECUTE":      kwEXECUTE,
	"EXISTS":       kwEXISTS,
	"FETCH":        kwFETCH,
	"FILE":         kwFILE,
	"FIRST":        kwFIRST,
	"FLOAT":        kwFLOAT,
	"FOR":          kwFOR,
	"FORCE":        kwFORCE,
	"FROM":         kwFROM,
	"FULL":         kwFULL,
	"FUNCTION":     kwFUNCTION,
	"GOTO":         kwGOTO,
	"GRANT":        kwGRANT,
	"GROUP":        kwGROUP,
	"HAVING":       kwHAVING,
	"IDENTIFIED":   kwIDENTIFIED,
	"IF":           kwIF,
	"IMMEDIATE":    kwIMMEDIATE,
	"IN":           kwIN,
	"INCREMENT":    kwINCREMENT,
	"INDEX":        kwINDEX,
	"INITIAL":      kwINITIAL,
	"INNER":        kwINNER,
	"INSERT":       kwINSERT,
	"INTEGER":      kwINTEGER,
	"INTERSECT":    kwINTERSECT,
	"INTO":         kwINTO,
	"IS":           kwIS,
	"JOIN":         kwJOIN,
	"KEEP":         kwKEEP,
	"LAST":         kwLAST,
	"LEFT":         kwLEFT,
	"LEVEL":        kwLEVEL,
	"LIKE":         kwLIKE,
	"LIKE2":        kwLIKE2,
	"LIKE4":        kwLIKE4,
	"LIKEC":        kwLIKEC,
	"LOCK":         kwLOCK,
	"LONG":         kwLONG,
	"LOOP":         kwLOOP,
	"MATERIALIZED": kwMATERIALIZED,
	"MAXEXTENTS":   kwMAXEXTENTS,
	"MERGE":        kwMERGE,
	"MINUS":        kwMINUS,
	"MODE":         kwMODE,
	"MODIFY":       kwMODIFY,
	"NATURAL":      kwNATURAL,
	"NEXT":         kwNEXT,
	"NOAUDIT":      kwNOAUDIT,
	"NOCACHE":      kwNOCACHE,
	"NOCOMPRESS":   kwNOCOMPRESS,
	"NOCYCLE":      kwNOCYCLE,
	"NOMAXVALUE":   kwNOMAXVALUE,
	"NOMINVALUE":   kwNOMINVALUE,
	"NOORDER":      kwNOORDER,
	"NOT":          kwNOT,
	"NOWAIT":       kwNOWAIT,
	"NULL":         kwNULL,
	"NUMBER":       kwNUMBER,
	"OF":           kwOF,
	"OFFLINE":      kwOFFLINE,
	"OFFSET":       kwOFFSET,
	"ON":           kwON,
	"ONLINE":       kwONLINE,
	"ONLY":         kwONLY,
	"OPEN":         kwOPEN,
	"OPTION":       kwOPTION,
	"OR":           kwOR,
	"ORDER":        kwORDER,
	"OUTER":        kwOUTER,
	"OVER":         kwOVER,
	"PACKAGE":      kwPACKAGE,
	"PARALLEL":     kwPARALLEL,
	"PARTITION":    kwPARTITION,
	"PCTFREE":      kwPCTFREE,
	"PERCENT":      kwPERCENT,
	"PIVOT":        kwPIVOT,
	"PRIOR":        kwPRIOR,
	"PRIVILEGES":   kwPRIVILEGES,
	"PROCEDURE":    kwPROCEDURE,
	"PUBLIC":       kwPUBLIC,
	"PURGE":        kwPURGE,
	"RANGE":        kwRANGE,
	"RAW":          kwRAW,
	"RENAME":       kwRENAME,
	"REPLACE":      kwREPLACE,
	"RETURN":       kwRETURN,
	"RETURNING":    kwRETURNING,
	"REVERSE":      kwREVERSE,
	"REVOKE":       kwREVOKE,
	"RIGHT":        kwRIGHT,
	"ROLLBACK":     kwROLLBACK,
	"ROW":          kwROW,
	"ROWID":        kwROWID,
	"ROWNUM":       kwROWNUM,
	"ROWS":         kwROWS,
	"SAMPLE":       kwSAMPLE,
	"SAVEPOINT":    kwSAVEPOINT,
	"SELECT":       kwSELECT,
	"SEQUENCE":     kwSEQUENCE,
	"SESSION":      kwSESSION,
	"SET":          kwSET,
	"SIZE":         kwSIZE,
	"SMALLINT":     kwSMALLINT,
	"START":        kwSTART,
	"SUBPARTITION": kwSUBPARTITION,
	"SYNONYM":      kwSYNONYM,
	"SYSDATE":      kwSYSDATE,
	"SYSTIMESTAMP": kwSYSTIMESTAMP,
	"SYSTEM":       kwSYSTEM,
	"TABLE":        kwTABLE,
	"THEN":         kwTHEN,
	"TIES":         kwTIES,
	"TO":           kwTO,
	"TRIGGER":      kwTRIGGER,
	"TRUNCATE":     kwTRUNCATE,
	"TYPE":         kwTYPE,
	"UNION":        kwUNION,
	"UNIQUE":       kwUNIQUE,
	"UNPIVOT":      kwUNPIVOT,
	"UPDATE":       kwUPDATE,
	"USING":        kwUSING,
	"VALUES":       kwVALUES,
	"VARCHAR":      kwVARCHAR,
	"VARCHAR2":     kwVARCHAR2,
	"VIEW":         kwVIEW,
	"WAIT":         kwWAIT,
	"WHEN":         kwWHEN,
	"WHERE":        kwWHERE,
	"WHILE":        kwWHILE,
	"WITH":         kwWITH,
	"WORK":         kwWORK,

	// Additional Oracle keywords for various features
	"ANALYZE":       kwANALYZE,
	"BEGIN":         kwBEGIN,
	"BITMAP":        kwBITMAP,
	"BLOB":          kwBLOB,
	"BODY":          kwBODY,
	"BULK":          kwBULK,
	"CACHE":         kwCACHE,
	"CASCADE":       kwCASCADE,
	"CHAR":          kwCHAR,
	"CLOB":          kwCLOB,
	"CLOSE":         kwCLOSE,
	"COLLECT":       kwCOLLECT,
	"COMPOUND":      kwCOMPOUND,
	"CONSTRAINTS":   kwCONSTRAINTS,
	"CONTINUE":      kwCONTINUE,
	"CURSOR":        kwCURSOR,
	"CYCLE":         kwCYCLE,
	"DATABASE":      kwDATABASE,
	"DECLARE":       kwDECLARE,
	"DEFERRABLE":    kwDEFERRABLE,
	"DEFERRED":      kwDEFERRED,
	"DETERMINISTIC": kwDETERMINISTIC,
	"EACH":          kwEACH,
	"ENABLE":        kwENABLE,
	"ESCAPE":        kwESCAPE,
	"EXPLAIN":       kwEXPLAIN,
	"FLASHBACK":     kwFLASHBACK,
	"FOREIGN":       kwFOREIGN,
	"GENERATED":     kwGENERATED,
	"GLOBAL":        kwGLOBAL,
	"HASH":          kwHASH,
	"IDENTITY":      kwIDENTITY,
	"INCLUDE":       kwINCLUDE,
	"INITIALLY":     kwINITIALLY,
	"INSTEAD":       kwINSTEAD,
	"INTERVAL":      kwINTERVAL,
	"INVISIBLE":     kwINVISIBLE,
	"ISOLATION":     kwISOLATION,
	"KEY":           kwKEY,
	"LIMIT":         kwLIMIT,
	"LINK":          kwLINK,
	"LIST":          kwLIST,
	"LOCAL":         kwLOCAL,
	"LOG":           kwLOG,
	"LOGGING":       kwLOGGING,
	"MATCHED":       kwMATCHED,
	"MAXVALUE":      kwMAXVALUE,
	"MINVALUE":      kwMINVALUE,
	"MODEL":         kwMODEL,
	"MULTISET":      kwMULTISET,
	"NCHAR":         kwNCHAR,
	"NCLOB":         kwNCLOB,
	"NOLOGGING":     kwNOLOGGING,
	"NOPARALLEL":    kwNOPARALLEL,
	"NVARCHAR2":     kwNVARCHAR2,
	"NULLS":         kwNULLS,
	"PIPELINED":     kwPIPELINED,
	"PLAN":          kwPLAN,
	"PRIMARY":       kwPRIMARY,
	"PRIVATE":       kwPRIVATE,
	"RAISE":         kwRAISE,
	"READ":          kwREAD,
	"RECORD":        kwRECORD,
	"RECURSIVE":     kwRECURSIVE,
	"REF":           kwREF,
	"REFERENCES":    kwREFERENCES,
	"REFRESH":       kwREFRESH,
	"REJECT":        kwREJECT,
	"RESTRICT":      kwRESTRICT,
	"RESULT_CACHE":  kwRESULT_CACHE,
	"ROLE":          kwROLE,
	"ROWTYPE":       kwROWTYPE,
	"SKIP":          kwSKIP,
	"STORAGE":       kwSTORAGE,
	"SUBTYPE":       kwSUBTYPE,
	"TABLESPACE":    kwTABLESPACE,
	"TEMPORARY":     kwTEMPORARY,
	"TIMESTAMP":     kwTIMESTAMP,
	"TRANSACTION":   kwTRANSACTION,
	"USER":          kwUSER,
	"VALIDATE":      kwVALIDATE,
	"VARRAY":        kwVARRAY,
	"VIRTUAL":       kwVIRTUAL,
	"WRITE":         kwWRITE,
	"ZONE":          kwZONE,

	// Additional Oracle-specific
	"AFTER":               kwAFTER,
	"ALWAYS":              kwALWAYS,
	"BEFORE":              kwBEFORE,
	"CONNECT_BY_ROOT":     kwCONNECT_BY_ROOT,
	"CONSTANT":            kwCONSTANT,
	"DENSE_RANK":          kwDENSE_RANK,
	"DISABLE":             kwDISABLE,
	"ERRORS":              kwERRORS,
	"EXIT":                kwEXIT,
	"FORALL":              kwFORALL,
	"LOCKED":              kwLOCKED,
	"NOCOPY":              kwNOCOPY,
	"OUT":                 kwOUT,
	"PARALLEL_ENABLE":     kwPARALLEL_ENABLE,
	"PRAGMA":              kwPRAGMA,
	"RELY":                kwRELY,
	"REWRITE":             kwREWRITE,
	"SCN":                 kwSCN,
	"SEED":                kwSEED,
	"SYS_CONNECT_BY_PATH": kwSYS_CONNECT_BY_PATH,
	"UNDER":               kwUNDER,

	// SAMPLE / FLASHBACK query keywords
	"BLOCK":    kwBLOCK,
	"VERSIONS": kwVERSIONS,

	// MODEL clause keywords
	"AUTOMATIC": kwAUTOMATIC,
	"DECREMENT": kwDECREMENT,
	"DIMENSION": kwDIMENSION,
	"ITERATE":   kwITERATE,
	"MAIN":      kwMAIN,
	"MEASURES":  kwMEASURES,
	"NAV":       kwNAV,
	"REFERENCE": kwREFERENCE,
	"RULES":     kwRULES,
	"SEQUENTIAL": kwSEQUENTIAL,
	"UNTIL":     kwUNTIL,
	"UPDATED":   kwUPDATED,
	"UPSERT":    kwUPSERT,

	// Analytic function keywords
	"FOLLOWING":  kwFOLLOWING,
	"GROUPS":     kwGROUPS,
	"PRECEDING":  kwPRECEDING,
	"UNBOUNDED":  kwUNBOUNDED,
	"WITHIN":     kwWITHIN,

	// Grouping extension keywords
	"CUBE":     kwCUBE,
	"GROUPING": kwGROUPING,
	"ROLLUP":   kwROLLUP,
	"SETS":     kwSETS,

	// LATERAL / XMLTABLE / JSON_TABLE keywords
	"COLUMNS":    kwCOLUMNS,
	"JSON_TABLE": kwJSON_TABLE,
	"LATERAL":    kwLATERAL,
	"NESTED":     kwNESTED,
	"ORDINALITY": kwORDINALITY,
	"PASSING":    kwPASSING,
	"PATH":       kwPATH,
	"XMLTABLE":   kwXMLTABLE,

	// Cursor / MULTISET / TREAT keywords
	"EXCEPT": kwEXCEPT,
	"TREAT":  kwTREAT,

	// XML expression keywords
	"CONTENT":      kwCONTENT,
	"NAME":         kwNAME,
	"XMLAGG":       kwXMLAGG,
	"XMLELEMENT":   kwXMLELEMENT,
	"XMLFOREST":    kwXMLFOREST,
	"XMLPARSE":     kwXMLPARSE,
	"XMLROOT":      kwXMLROOT,
	"XMLSERIALIZE": kwXMLSERIALIZE,

	// Utility statement keywords
	"CALL":  kwCALL,
	"SHARE": kwSHARE,

	// JSON expression keywords
	"FORMAT":         kwFORMAT,
	"JSON":           kwJSON,
	"JSON_ARRAY":     kwJSON_ARRAY,
	"JSON_EXISTS":    kwJSON_EXISTS,
	"JSON_MERGEPATCH": kwJSON_MERGEPATCH,
	"JSON_OBJECT":    kwJSON_OBJECT,
	"JSON_QUERY":     kwJSON_QUERY,
	"JSON_VALUE":     kwJSON_VALUE,
}

// Keyword token constants.
const (
	kwACCESS = iota + 2000
	kwADD
	kwAFTER
	kwALL
	kwALTER
	kwALWAYS
	kwANALYZE
	kwAND
	kwANY
	kwAS
	kwASC
	kwAUDIT
	kwBEFORE
	kwBEGIN
	kwBETWEEN
	kwBITMAP
	kwBLOB
	kwBODY
	kwBULK
	kwBY
	kwCACHE
	kwCASCADE
	kwCASE
	kwCAST
	kwCHAR
	kwCHECK
	kwCLOB
	kwCLOSE
	kwCLUSTER
	kwCOLLECT
	kwCOLUMN
	kwCOMMENT
	kwCOMMIT
	kwCOMPOUND
	kwCOMPRESS
	kwCONNECT
	kwCONNECT_BY_ROOT
	kwCONSTANT
	kwCONSTRAINT
	kwCONSTRAINTS
	kwCONTINUE
	kwCREATE
	kwCROSS
	kwCURRENT
	kwCURSOR
	kwCYCLE
	kwDATABASE
	kwDATE
	kwDECIMAL
	kwDECLARE
	kwDECODE
	kwDEFAULT
	kwDEFERRABLE
	kwDEFERRED
	kwDELETE
	kwDENSE_RANK
	kwDESC
	kwDETERMINISTIC
	kwDISABLE
	kwDISTINCT
	kwDROP
	kwEACH
	kwELSE
	kwELSIF
	kwENABLE
	kwEND
	kwERRORS
	kwESCAPE
	kwEXCEPTION
	kwEXCLUSIVE
	kwEXECUTE
	kwEXISTS
	kwEXIT
	kwEXPLAIN
	kwFETCH
	kwFILE
	kwFIRST
	kwFLASHBACK
	kwFLOAT
	kwFOR
	kwFORALL
	kwFORCE
	kwFOREIGN
	kwFROM
	kwFULL
	kwFUNCTION
	kwGENERATED
	kwGLOBAL
	kwGOTO
	kwGRANT
	kwGROUP
	kwHASH
	kwHAVING
	kwIDENTIFIED
	kwIDENTITY
	kwIF
	kwIMMEDIATE
	kwIN
	kwINCLUDE
	kwINCREMENT
	kwINDEX
	kwINITIAL
	kwINITIALLY
	kwINNER
	kwINSERT
	kwINSTEAD
	kwINTEGER
	kwINTERSECT
	kwINTERVAL
	kwINTO
	kwINVISIBLE
	kwIS
	kwISOLATION
	kwJOIN
	kwKEEP
	kwKEY
	kwLAST
	kwLEFT
	kwLEVEL
	kwLIKE
	kwLIKE2
	kwLIKE4
	kwLIKEC
	kwLIMIT
	kwLINK
	kwLIST
	kwLOCAL
	kwLOCK
	kwLOCKED
	kwLOG
	kwLOGGING
	kwLONG
	kwLOOP
	kwMATCHED
	kwMATERIALIZED
	kwMAXEXTENTS
	kwMAXVALUE
	kwMERGE
	kwMINUS
	kwMINVALUE
	kwMODE
	kwMODEL
	kwMODIFY
	kwMULTISET
	kwNATURAL
	kwNCHAR
	kwNCLOB
	kwNEXT
	kwNOAUDIT
	kwNOCACHE
	kwNOCOMPRESS
	kwNOCOPY
	kwNOCYCLE
	kwNOLOGGING
	kwNOMAXVALUE
	kwNOMINVALUE
	kwNOORDER
	kwNOPARALLEL
	kwNOT
	kwNOWAIT
	kwNULL
	kwNULLS
	kwNUMBER
	kwNVARCHAR2
	kwOF
	kwOFFLINE
	kwOFFSET
	kwON
	kwONLINE
	kwONLY
	kwOPEN
	kwOPTION
	kwOR
	kwORDER
	kwOUT
	kwOUTER
	kwOVER
	kwPACKAGE
	kwPARALLEL
	kwPARALLEL_ENABLE
	kwPARTITION
	kwPCTFREE
	kwPERCENT
	kwPIPELINED
	kwPIVOT
	kwPLAN
	kwPRAGMA
	kwPRIMARY
	kwPRIOR
	kwPRIVATE
	kwPRIVILEGES
	kwPROCEDURE
	kwPUBLIC
	kwPURGE
	kwRAISE
	kwRANGE
	kwRAW
	kwREAD
	kwRECORD
	kwRECURSIVE
	kwREF
	kwREFERENCES
	kwREFRESH
	kwREJECT
	kwRELY
	kwRENAME
	kwREPLACE
	kwRESTRICT
	kwRESULT_CACHE
	kwRETURN
	kwRETURNING
	kwREVERSE
	kwREVOKE
	kwREWRITE
	kwRIGHT
	kwROLE
	kwROLLBACK
	kwROW
	kwROWID
	kwROWNUM
	kwROWS
	kwROWTYPE
	kwSAMPLE
	kwSAVEPOINT
	kwSCN
	kwSEED
	kwSELECT
	kwSEQUENCE
	kwSESSION
	kwSET
	kwSIZE
	kwSKIP
	kwSMALLINT
	kwSTART
	kwSTORAGE
	kwSUBPARTITION
	kwSUBTYPE
	kwSYNONYM
	kwSYS_CONNECT_BY_PATH
	kwSYSDATE
	kwSYSTEM
	kwSYSTIMESTAMP
	kwTABLE
	kwTABLESPACE
	kwTEMPORARY
	kwTHEN
	kwTIES
	kwTIMESTAMP
	kwTO
	kwTRANSACTION
	kwTRIGGER
	kwTRUNCATE
	kwTYPE
	kwUNDER
	kwUNION
	kwUNIQUE
	kwUNPIVOT
	kwUPDATE
	kwUSER
	kwUSING
	kwVALIDATE
	kwVALUES
	kwVARCHAR
	kwVARCHAR2
	kwVARRAY
	kwVIEW
	kwVIRTUAL
	kwWAIT
	kwWHEN
	kwWHERE
	kwWHILE
	kwWITH
	kwWORK
	kwWRITE
	kwZONE

	// SAMPLE / FLASHBACK query keywords
	kwBLOCK
	kwVERSIONS

	// MODEL clause keywords
	kwAUTOMATIC
	kwDECREMENT
	kwDIMENSION
	kwITERATE
	kwMAIN
	kwMEASURES
	kwNAV
	kwREFERENCE
	kwRULES
	kwSEQUENTIAL
	kwUNTIL
	kwUPDATED
	kwUPSERT

	// Analytic function keywords
	kwFOLLOWING
	kwGROUPS
	kwPRECEDING
	kwUNBOUNDED
	kwWITHIN

	// Grouping extension keywords
	kwCUBE
	kwGROUPING
	kwROLLUP
	kwSETS

	// LATERAL / XMLTABLE / JSON_TABLE keywords
	kwCOLUMNS
	kwJSON_TABLE
	kwLATERAL
	kwNESTED
	kwORDINALITY
	kwPASSING
	kwPATH
	kwXMLTABLE

	// Cursor / MULTISET / TREAT keywords
	kwEXCEPT
	kwTREAT

	// XML expression keywords
	kwCONTENT
	kwNAME
	kwXMLAGG
	kwXMLELEMENT
	kwXMLFOREST
	kwXMLPARSE
	kwXMLROOT
	kwXMLSERIALIZE

	// Utility statement keywords
	kwCALL
	kwSHARE

	// JSON expression keywords
	kwFORMAT
	kwJSON
	kwJSON_ARRAY
	kwJSON_EXISTS
	kwJSON_MERGEPATCH
	kwJSON_OBJECT
	kwJSON_QUERY
	kwJSON_VALUE
)
