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

// parseViewConstraintState parses the constraint state allowed on view
// constraints. Verified against Oracle 23ai: DISABLE is mandatory
// (ORA-02000 "missing DISABLE keyword" otherwise), RELY|NORELY may precede
// it, NOVALIDATE may follow it, and VALIDATE is rejected outright
// (ORA-03082). Table-only clauses (DEFERRABLE, INITIALLY, USING INDEX,
// ENABLE) are not accepted.
//
//	[ RELY | NORELY ] DISABLE [ NOVALIDATE ]
func (p *Parser) parseViewConstraintState() error {
	if p.cur.Type == kwRELY || p.isIdentLikeStr("NORELY") {
		p.advance()
	}
	if p.cur.Type != kwDISABLE {
		return p.syntaxErrorAtCur()
	}
	p.advance() // consume DISABLE
	if p.isIdentLikeStr("NOVALIDATE") {
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
	"INDEXING":   true,
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

	// (create_index_statement) — validate the required structure CREATE
	// [UNIQUE|BITMAP] INDEX name ON table (...) before skipping the balanced
	// remainder (index attributes). Oracle raises ORA-02000 for a missing
	// CREATE keyword, ORA-00953 for a missing index name, and ORA-00969 for
	// a missing ON.
	if p.cur.Type == '(' {
		p.advance() // consume '('
		if p.cur.Type != kwCREATE {
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume CREATE
		if p.cur.Type == kwUNIQUE || p.cur.Type == kwBITMAP {
			p.advance()
		}
		if p.cur.Type != kwINDEX {
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume INDEX
		indexName, err := p.parseObjectName()
		if err != nil {
			return err
		}
		if indexName == nil || indexName.Name == "" {
			return p.syntaxErrorAtCur()
		}
		if p.cur.Type != kwON {
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume ON
		tableName, err := p.parseObjectName()
		if err != nil {
			return err
		}
		if tableName == nil || tableName.Name == "" {
			return p.syntaxErrorAtCur()
		}
		if p.cur.Type != '(' {
			return p.syntaxErrorAtCur()
		}
		depth := 1 // the outer '(' consumed above is still open
		for depth > 0 && p.cur.Type != tokEOF {
			switch p.cur.Type {
			case '(':
				depth++
			case ')':
				depth--
			}
			p.advance()
		}
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
			// LOCAL { STORE IN (tablespace [, ...]) | (partition specs) }.
			// The two forms are mutually exclusive — Oracle 23ai raises
			// ORA-14153 when both are given.
			if p.isIdentLikeStr("STORE") && p.peekNext().Type == kwIN {
				p.advance() // consume STORE
				p.advance() // consume IN
				if p.cur.Type != '(' {
					return p.syntaxErrorAtCur()
				}
				p.skipParenthesized()
			} else if p.cur.Type == '(' {
				p.skipParenthesized()
			}

		case p.cur.Type == kwGLOBAL:
			p.advance() // consume GLOBAL
			// GLOBAL [ PARTITION BY { RANGE | HASH } (cols)
			//   { (partition specs) | PARTITIONS n [ STORE IN (ts, ...) ] } ]
			if p.cur.Type == kwPARTITION {
				p.advance() // consume PARTITION
				if p.cur.Type != kwBY {
					return p.syntaxErrorAtCur()
				}
				p.advance() // consume BY
				isHash := p.cur.Type == kwHASH
				if p.cur.Type == kwRANGE || p.cur.Type == kwHASH {
					p.advance()
				}
				if p.cur.Type != '(' {
					return p.syntaxErrorAtCur()
				}
				p.skipParenthesized()
				if isHash && p.isIdentLikeStr("PARTITIONS") {
					// hash_partitions_by_quantity: PARTITIONS n [STORE IN (...)].
					// HASH only — Oracle rejects the quantity form for RANGE.
					p.advance() // consume PARTITIONS
					if p.cur.Type != tokICONST {
						return p.syntaxErrorAtCur()
					}
					p.advance()
					if p.isIdentLikeStr("STORE") && p.peekNext().Type == kwIN {
						p.advance() // consume STORE
						p.advance() // consume IN
						if p.cur.Type != '(' {
							return p.syntaxErrorAtCur()
						}
						p.skipParenthesized()
					}
				} else if p.cur.Type == '(' {
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
			// COMPRESS [ integer | ADVANCED [ LOW | HIGH ] ]
			p.advance()
			if p.cur.Type == tokICONST {
				p.advance()
			} else if p.isIdentLikeStr("ADVANCED") {
				p.advance()
				if p.isIdentLikeStr("LOW") || p.isIdentLikeStr("HIGH") {
					p.advance()
				}
			}

		case p.cur.Type == kwPARALLEL:
			p.advance()
			if p.cur.Type == tokICONST {
				p.advance()
			}

		case p.isIdentLikeStr("INDEXING"):
			// INDEXING { FULL | PARTIAL }
			p.advance()
			if !p.isIdentLikeStr("FULL") && !p.isIdentLikeStr("PARTIAL") {
				return p.syntaxErrorAtCur()
			}
			p.advance()

		case p.cur.Type == kwLOGGING || p.cur.Type == kwNOLOGGING ||
			p.cur.Type == kwONLINE || p.cur.Type == kwREVERSE ||
			p.cur.Type == kwNOPARALLEL:
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
