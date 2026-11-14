package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreatePmemFilestoreStmt parses a CREATE PMEM FILESTORE statement. It
// is called after CREATE PMEM FILESTORE.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/create-pmem-filestore.html
//
//	CREATE PMEM FILESTORE filestore_name
//	    MOUNTPOINT 'file_path'
//	    BACKINGFILE 'backing_file_path' SIZE size_value BLOCKSIZE blocksize_value
//	    [ AUTOEXTEND { ON | OFF }
//	        [ NEXT size_value ]
//	        [ MAXSIZE { UNLIMITED | size_value } ] ]
//
// Oracle 23ai takes the name only quoted (ORA-00931 for a word) and without a
// schema (ORA-00922 at the '.'). It takes the clauses in any order, each of
// them optional and repeatable, and NEXT and MAXSIZE, in that order, only
// after AUTOEXTEND ON (ORA-00922). A size is an integer with an optional
// K, M, G, T, P, or E (ORA-02017 for a decimal, ORA-02237 for a string).
func (p *Parser) parseCreatePmemFilestoreStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_PMEM_FILESTORE,
		Loc:        nodes.Loc{Start: start},
	}
	if p.cur.Type != tokQIDENT {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Name = &nodes.ObjectName{Name: p.cur.Str, Loc: nodes.Loc{Start: p.cur.Loc, End: p.cur.End}}
	p.advance()

	opts := &nodes.List{}
	add := func(key, value string, optStart int) {
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: key, Value: value, Loc: nodes.Loc{Start: optStart, End: p.prev.End}})
	}
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		optStart := p.pos()
		switch {
		case p.isKeywordStr("MOUNTPOINT"), p.isKeywordStr("BACKINGFILE"):
			key := p.cur.Str
			p.advance()
			if p.cur.Type != tokSCONST && p.cur.Type != tokNCHARLIT {
				return nil, p.syntaxErrorAtCur()
			}
			value := p.cur.Str
			p.advance()
			add(key, value, optStart)
		case p.isKeywordStr("SIZE"), p.isKeywordStr("BLOCKSIZE"):
			key := p.cur.Str
			p.advance()
			size, err := p.parsePmemSize()
			if err != nil {
				return nil, err
			}
			add(key, size, optStart)
		case p.isKeywordStr("AUTOEXTEND"):
			p.advance() // consume AUTOEXTEND
			switch {
			case p.cur.Type == kwON:
				p.advance() // consume ON
				add("AUTOEXTEND", "ON", optStart)
				if p.isKeywordStr("NEXT") {
					nextStart := p.pos()
					p.advance() // consume NEXT
					size, err := p.parsePmemSize()
					if err != nil {
						return nil, err
					}
					add("NEXT", size, nextStart)
				}
				if p.isKeywordStr("MAXSIZE") {
					maxStart := p.pos()
					p.advance() // consume MAXSIZE
					size := "UNLIMITED"
					if p.isKeywordStr("UNLIMITED") {
						p.advance()
					} else {
						var err error
						if size, err = p.parsePmemSize(); err != nil {
							return nil, err
						}
					}
					add("MAXSIZE", size, maxStart)
				}
			case p.isKeywordStr("OFF"):
				p.advance() // consume OFF
				add("AUTOEXTEND", "OFF", optStart)
			default:
				return nil, p.syntaxErrorAtCur()
			}
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parsePmemSize parses a size_value of CREATE PMEM FILESTORE: an integer
// with an optional K, M, G, T, P, or E.
func (p *Parser) parsePmemSize() (string, error) {
	if p.cur.Type != tokICONST {
		return "", p.syntaxErrorAtCur()
	}
	return p.parseSizeValue()
}
