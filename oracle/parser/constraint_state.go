package parser

// constraintState carries the fields of Oracle's constraint_state clause that
// are preserved in the AST. Subclauses that are parsed but not preserved
// (ENABLE/DISABLE, VALIDATE/NOVALIDATE, RELY/NORELY, exceptions_clause) leave
// no trace here.
type constraintState struct {
	Deferrable      bool
	Initially       string // "DEFERRED" or "IMMEDIATE"
	Tablespace      string // USING INDEX ... TABLESPACE
	UsingIndexLocal bool   // USING INDEX LOCAL
}

// parseConstraintState parses Oracle's constraint_state clause. Oracle
// enforces the documented slot order (verified against Oracle 23ai, which
// raises ORA-03075 for out-of-order subclauses such as ENABLE USING INDEX
// or DISABLE RELY):
//
//	[ [NOT] DEFERRABLE ] [ INITIALLY { DEFERRED | IMMEDIATE } ]  -- either order
//	[ RELY | NORELY ]
//	[ using_index_clause ]
//	[ ENABLE | DISABLE ]
//	[ VALIDATE | NOVALIDATE ]
//
// EXCEPTIONS INTO is not part of this clause: Oracle rejects it inside a
// CREATE TABLE constraint (ORA-00922) and accepts it only in ALTER TABLE
// enable/modify-constraint contexts — see parseExceptionsIntoClause.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/constraint.html
func (p *Parser) parseConstraintState(cs *constraintState) error {
	// [NOT] DEFERRABLE and INITIALLY {DEFERRED|IMMEDIATE}: one group, either
	// order within it (Oracle accepts INITIALLY DEFERRED DEFERRABLE), each at
	// most once.
	seenDeferrable, seenInitially := false, false
	for {
		if !seenDeferrable && p.cur.Type == kwDEFERRABLE {
			cs.Deferrable = true
			seenDeferrable = true
			p.advance()
			continue
		}
		if !seenDeferrable && p.cur.Type == kwNOT && p.peekNext().Type == kwDEFERRABLE {
			p.advance() // consume NOT
			p.advance() // consume DEFERRABLE
			cs.Deferrable = false
			seenDeferrable = true
			continue
		}
		if !seenInitially && p.cur.Type == kwINITIALLY {
			p.advance() // consume INITIALLY
			switch p.cur.Type {
			case kwDEFERRED:
				cs.Initially = "DEFERRED"
				p.advance()
			case kwIMMEDIATE:
				cs.Initially = "IMMEDIATE"
				p.advance()
			default:
				return p.syntaxErrorAtCur()
			}
			seenInitially = true
			continue
		}
		break
	}

	// [ RELY | NORELY ]
	if p.cur.Type == kwRELY || p.isIdentLikeStr("NORELY") {
		p.advance()
	}

	// [ using_index_clause ]
	if p.cur.Type == kwUSING && p.peekNext().Type == kwINDEX {
		if err := p.parseUsingIndexClause(cs); err != nil {
			return err
		}
	}

	// [ ENABLE | DISABLE ]
	if p.cur.Type == kwENABLE || p.cur.Type == kwDISABLE {
		p.advance()
	}

	// [ VALIDATE | NOVALIDATE ]
	if p.cur.Type == kwVALIDATE || p.isIdentLikeStr("NOVALIDATE") {
		p.advance()
	}

	return nil
}

// parseExceptionsIntoClause parses [ EXCEPTIONS INTO [schema.]table ] if
// present. Oracle allows this only in ALTER TABLE enable/modify-constraint
// contexts, not inside CREATE TABLE constraints.
func (p *Parser) parseExceptionsIntoClause() error {
	if !p.isIdentLikeStr("EXCEPTIONS") || p.peekNext().Type != kwINTO {
		return nil
	}
	p.advance() // consume EXCEPTIONS
	p.advance() // consume INTO
	name, err := p.parseObjectName()
	if err != nil {
		return err
	}
	if name == nil || name.Name == "" {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// identLikeUsingIndexProperties are index_properties that lex as plain
// identifiers rather than keywords. They must not be mistaken for an index
// name after USING INDEX.
var identLikeUsingIndexProperties = map[string]bool{
	"PCTUSED":    true,
	"INITRANS":   true,
	"MAXTRANS":   true,
	"COMPUTE":    true,
	"NOCOMPRESS": true,
	"SORT":       true,
	"NOSORT":     true,
	"VISIBLE":    true,
	"INVISIBLE":  true,
	"STORE":      true,
	"NORELY":     true,
	"NOVALIDATE": true,
	"EXCEPTIONS": true,
}

// parseUsingIndexClause parses using_index_clause:
//
//	USING INDEX { [schema.]index | (create_index_statement) | index_properties }
//
// The properties form is parsed with an explicit whitelist so the clause
// terminates safely inside a CREATE TABLE column list: it never consumes the
// table's closing ')' or a following ','. This is deliberately not shared
// with parseCreateIndexAttributes, whose option collector runs to the
// statement terminator.
func (p *Parser) parseUsingIndexClause(cs *constraintState) error {
	p.advance() // consume USING
	p.advance() // consume INDEX

	// (create_index_statement)
	if p.cur.Type == '(' {
		p.skipParenthesized()
		return nil
	}

	// [schema.]index — a plain identifier that is not an identifier-like
	// index property or a constraint_state continuation. Quoted identifiers
	// are always names.
	if p.cur.Type == tokQIDENT ||
		(p.cur.Type == tokIDENT && !identLikeUsingIndexProperties[p.cur.Str]) {
		name, err := p.parseObjectName()
		if err != nil {
			return err
		}
		if name == nil || name.Name == "" {
			return p.syntaxErrorAtCur()
		}
		return nil
	}

	// index_properties
	for {
		switch {
		case p.cur.Type == kwTABLESPACE:
			p.advance() // consume TABLESPACE
			name, err := p.parseIdentifier()
			if err != nil {
				return err
			}
			if name == "" {
				return p.syntaxErrorAtCur()
			}
			cs.Tablespace = name

		case p.cur.Type == kwLOCAL:
			p.advance() // consume LOCAL
			cs.UsingIndexLocal = true
			// LOCAL [ STORE IN (tablespace [, ...]) ] [ (partition specs) ]
			if p.isIdentLikeStr("STORE") && p.peekNext().Type == kwIN {
				p.advance() // consume STORE
				p.advance() // consume IN
				if p.cur.Type != '(' {
					return p.syntaxErrorAtCur()
				}
				p.skipParenthesized()
			}
			if p.cur.Type == '(' {
				p.skipParenthesized()
			}

		case p.cur.Type == kwGLOBAL:
			p.advance() // consume GLOBAL
			// GLOBAL [ PARTITION BY { RANGE | HASH } (cols) [ (partition specs) ] ]
			if p.cur.Type == kwPARTITION {
				p.advance() // consume PARTITION
				if p.cur.Type != kwBY {
					return p.syntaxErrorAtCur()
				}
				p.advance() // consume BY
				if p.cur.Type == kwRANGE || p.cur.Type == kwHASH {
					p.advance()
				}
				if p.cur.Type != '(' {
					return p.syntaxErrorAtCur()
				}
				p.skipParenthesized()
				if p.cur.Type == '(' {
					p.skipParenthesized()
				}
			}

		case p.cur.Type == kwPCTFREE || p.isIdentLikeStr("PCTUSED") ||
			p.isIdentLikeStr("INITRANS") || p.isIdentLikeStr("MAXTRANS"):
			p.advance()
			if p.cur.Type != tokICONST {
				return p.syntaxErrorAtCur()
			}
			p.advance()

		case p.cur.Type == kwSTORAGE:
			p.advance() // consume STORAGE
			if p.cur.Type != '(' {
				return p.syntaxErrorAtCur()
			}
			p.skipParenthesized()

		case p.isIdentLikeStr("COMPUTE"):
			// COMPUTE STATISTICS (emitted by DBMS_METADATA)
			p.advance() // consume COMPUTE
			if !p.isIdentLikeStr("STATISTICS") {
				return p.syntaxErrorAtCur()
			}
			p.advance() // consume STATISTICS

		case p.cur.Type == kwCOMPRESS:
			p.advance()
			if p.cur.Type == tokICONST {
				p.advance()
			}

		case p.cur.Type == kwLOGGING || p.cur.Type == kwNOLOGGING ||
			p.cur.Type == kwONLINE || p.cur.Type == kwREVERSE:
			p.advance()

		case p.isIdentLikeStr("NOCOMPRESS") || p.isIdentLikeStr("SORT") ||
			p.isIdentLikeStr("NOSORT") || p.isIdentLikeStr("VISIBLE") ||
			p.isIdentLikeStr("INVISIBLE"):
			p.advance()

		default:
			return nil
		}
	}
}
