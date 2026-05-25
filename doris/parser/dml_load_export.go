package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// TRUNCATE TABLE (T6.1)
// ---------------------------------------------------------------------------

// parseTruncateTable parses:
//
//	TRUNCATE TABLE [db.]name [PARTITION(p1, p2, ...)] [FORCE]
//
// The TRUNCATE keyword has NOT yet been consumed when this is called.
func (p *Parser) parseTruncateTable() (ast.Node, error) {
	truncateTok := p.advance() // consume TRUNCATE

	// Required TABLE keyword
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}

	stmt := &ast.TruncateTableStmt{
		Loc: ast.Loc{Start: truncateTok.Loc.Start},
	}

	// Table name
	target, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Target = target
	endLoc := ast.NodeLoc(target)

	// Optional PARTITION(p1, p2, ...) or PARTITIONS(...)
	if p.cur.Kind == kwPARTITION || p.cur.Kind == kwPARTITIONS {
		p.advance() // consume PARTITION / PARTITIONS
		partitions, err := p.parsePartitionNameList()
		if err != nil {
			return nil, err
		}
		stmt.Partition = partitions
		if p.prev.Loc.End > endLoc.End {
			endLoc = p.prev.Loc
		}
	}

	// Optional FORCE
	if p.cur.Kind == kwFORCE {
		endLoc = p.cur.Loc
		p.advance()
		stmt.Force = true
	}

	stmt.Loc.End = endLoc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// COPY INTO (T6.1)
// ---------------------------------------------------------------------------

// parseCopyInto parses:
//
//	COPY INTO target FROM @stage_or_path
//	    [FILES = ('f1', 'f2')]
//	    [PATTERN = 'glob']
//	    [PROPERTIES(...)]
//
// The COPY keyword has NOT yet been consumed when this is called.
func (p *Parser) parseCopyInto() (ast.Node, error) {
	copyTok := p.advance() // consume COPY

	// Required INTO keyword
	if _, err := p.expect(kwINTO); err != nil {
		return nil, err
	}

	stmt := &ast.CopyIntoStmt{
		Loc: ast.Loc{Start: copyTok.Loc.Start},
	}

	// Target table name
	target, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Target = target

	// Required FROM keyword
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}

	// Source: @stage, 'path', or bare identifier
	// Handle leading @ sign for stage references like @stage_name
	source, err := p.parseCopySource()
	if err != nil {
		return nil, err
	}
	stmt.Source = source

	// Optional FILES = ('f1', 'f2', ...)
	for p.cur.Kind != tokEOF {
		if isIdentText(p, "files") || p.cur.Kind == kwFILE {
			p.advance() // consume FILES/FILE
			if _, err := p.expect(int('=')); err != nil {
				return nil, err
			}
			files, err := p.parseStringList()
			if err != nil {
				return nil, err
			}
			stmt.Files = files
		} else if isIdentText(p, "pattern") || p.cur.Kind == kwPATH {
			// PATTERN = '...'
			p.advance() // consume PATTERN
			if _, err := p.expect(int('=')); err != nil {
				return nil, err
			}
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Pattern = p.cur.Str
			p.advance()
		} else if p.cur.Kind == kwPROPERTIES {
			props, err := p.parseProperties()
			if err != nil {
				return nil, err
			}
			stmt.Properties = props
		} else {
			break
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseCopySource parses the FROM source in COPY INTO.
// Handles: @stage_name, 'string_path', or bare identifier.
func (p *Parser) parseCopySource() (string, error) {
	// @stage_name: '@' is not a keyword, it's ASCII 64
	if p.cur.Kind == int('@') {
		p.advance() // consume '@'
		name, _, err := p.parseIdentifier()
		if err != nil {
			return "", err
		}
		// Handle @stage/path segments
		var sb strings.Builder
		sb.WriteByte('@')
		sb.WriteString(name)
		// Consume optional sub-path segments (slash-separated)
		for p.cur.Kind == int('/') {
			p.advance()
			sb.WriteByte('/')
			if p.cur.Kind == tokString || p.cur.Kind == tokIdent || isIdentifierToken(p.cur.Kind) {
				sb.WriteString(p.cur.Str)
				p.advance()
			}
		}
		return sb.String(), nil
	}

	// String literal path: 'path'
	if p.cur.Kind == tokString {
		s := p.cur.Str
		p.advance()
		return s, nil
	}

	// Bare identifier (stage name without @)
	if isIdentifierToken(p.cur.Kind) {
		name, _, err := p.parseIdentifier()
		if err != nil {
			return "", err
		}
		return name, nil
	}

	return "", p.syntaxErrorAtCur()
}

// isIdentText checks whether cur is an identifier with the given lowercase text.
// Used for pseudo-keywords like FILES and PATTERN that may appear as bare idents.
func isIdentText(p *Parser, lower string) bool {
	if p.cur.Kind == tokIdent || (p.cur.Kind >= 700 && !IsReserved(p.cur.Kind)) {
		return strings.ToLower(p.cur.Str) == lower
	}
	return false
}

// parseStringList parses ('str1', 'str2', ...).
func (p *Parser) parseStringList() ([]string, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	var strs []string
	for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		strs = append(strs, p.cur.Str)
		p.advance()
		if p.cur.Kind == int(',') {
			p.advance() // consume ','
		}
	}

	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return strs, nil
}

// ---------------------------------------------------------------------------
// LOAD (T6.1)
// ---------------------------------------------------------------------------

// parseLoad parses:
//
//	LOAD LABEL [db.]label_name
//	    (DATA INFILE('path') [NEGATIVE] INTO TABLE t [PARTITION(...)]
//	         [COLUMNS FROM PATH AS (col, ...)]
//	         [COLUMNS (c1, c2, ...)]
//	         [SET (...)]
//	         [WHERE expr]
//	    [, DATA INFILE...])
//	    [WITH BROKER broker_name [(key=val, ...)]]
//	    [PROPERTIES(...)]
//	    [COMMENT 'text']
//
// The LOAD keyword has NOT yet been consumed when this is called.
func (p *Parser) parseLoad() (ast.Node, error) {
	loadTok := p.advance() // consume LOAD

	stmt := &ast.LoadDataStmt{
		Loc: ast.Loc{Start: loadTok.Loc.Start},
	}

	// LABEL label_name (required)
	if _, err := p.expect(kwLABEL); err != nil {
		return nil, err
	}

	// Label: may be db_name.label_name
	labelParts, err := p.parseLoadLabel()
	if err != nil {
		return nil, err
	}
	stmt.Label = labelParts

	// Required '(' to open data desc list
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	// Parse one or more data descriptions separated by ','
	desc, err := p.parseLoadDataDesc()
	if err != nil {
		return nil, err
	}
	stmt.DataDescs = append(stmt.DataDescs, desc)

	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		desc, err = p.parseLoadDataDesc()
		if err != nil {
			return nil, err
		}
		stmt.DataDescs = append(stmt.DataDescs, desc)
	}

	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}

	// Optional WITH BROKER broker_name [(key=val, ...)]
	if p.cur.Kind == kwWITH {
		p.advance() // consume WITH
		if _, err := p.expect(kwBROKER); err != nil {
			return nil, err
		}
		brokerName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.BrokerName = brokerName

		// Optional broker properties (...)
		if p.cur.Kind == int('(') {
			if err := p.skipParenthesized(); err != nil {
				return nil, err
			}
		}
	}

	// Optional PROPERTIES(...)
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
	}

	// Optional COMMENT 'text'
	if p.cur.Kind == kwCOMMENT {
		p.advance() // consume COMMENT
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Comment = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseLoadLabel parses the load label, which may be db_name.label_name or
// just label_name.
func (p *Parser) parseLoadLabel() (string, error) {
	first, _, err := p.parseIdentifier()
	if err != nil {
		return "", err
	}
	if p.cur.Kind == int('.') {
		p.advance() // consume '.'
		second, _, err := p.parseIdentifier()
		if err != nil {
			return "", err
		}
		return first + "." + second, nil
	}
	return first, nil
}

// parseLoadDataDesc parses one DATA INFILE(...) INTO TABLE t ... clause.
func (p *Parser) parseLoadDataDesc() (*ast.LoadDataDesc, error) {
	desc := &ast.LoadDataDesc{}
	startLoc := p.cur.Loc

	// Optional NEGATIVE keyword before DATA
	if p.cur.Kind == kwNEGATIVE {
		p.advance() // consume NEGATIVE
		desc.Negative = true
	}

	// Required DATA keyword
	if _, err := p.expect(kwDATA); err != nil {
		return nil, err
	}

	// Required INFILE('path', ...)
	if _, err := p.expect(kwINFILE); err != nil {
		return nil, err
	}

	files, err := p.parseStringList()
	if err != nil {
		return nil, err
	}
	desc.SourceFiles = files

	// Optional FORMAT AS format_name
	if p.cur.Kind == kwFORMAT {
		p.advance() // consume FORMAT
		// optional AS keyword
		p.match(kwAS)
		if isIdentifierToken(p.cur.Kind) {
			desc.Format = p.cur.Str
			p.advance()
		}
	}

	// Required INTO TABLE target
	if _, err := p.expect(kwINTO); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	target, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	desc.Target = target

	// Parse optional modifiers after the table name
	for p.cur.Kind != tokEOF && p.cur.Kind != int(',') && p.cur.Kind != int(')') {
		switch p.cur.Kind {
		case kwPARTITION, kwPARTITIONS:
			p.advance() // consume PARTITION/PARTITIONS
			parts, err := p.parsePartitionNameList()
			if err != nil {
				return nil, err
			}
			desc.Partition = parts

		case kwCOLUMNS:
			p.advance() // consume COLUMNS
			if p.cur.Kind == kwFROM {
				// COLUMNS FROM PATH AS (col1, col2, ...)
				p.advance() // consume FROM
				if _, err := p.expect(kwPATH); err != nil {
					return nil, err
				}
				if _, err := p.expect(kwAS); err != nil {
					return nil, err
				}
				if _, err := p.expect(int('(')); err != nil {
					return nil, err
				}
				var pathCols []string
				for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
					col, _, err := p.parseIdentifier()
					if err != nil {
						return nil, err
					}
					pathCols = append(pathCols, col)
					if p.cur.Kind == int(',') {
						p.advance()
					}
				}
				if _, err := p.expect(int(')')); err != nil {
					return nil, err
				}
				desc.ColumnsFromPath = pathCols
			} else if p.cur.Kind == int('(') {
				// COLUMNS (c1, c2, ...)
				p.advance() // consume '('
				var cols []string
				for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
					col, _, err := p.parseIdentifier()
					if err != nil {
						return nil, err
					}
					cols = append(cols, col)
					if p.cur.Kind == int(',') {
						p.advance()
					}
				}
				if _, err := p.expect(int(')')); err != nil {
					return nil, err
				}
				desc.ColumnList = cols
			} else {
				return nil, p.syntaxErrorAtCur()
			}

		case kwSET:
			p.advance() // consume SET
			// Capture raw SET(...) as best-effort text
			raw, err := p.captureParenthesized()
			if err != nil {
				return nil, err
			}
			desc.SetExpr = raw

		case kwWHERE:
			p.advance() // consume WHERE
			// Capture raw WHERE expr as best-effort text (until next clause keyword)
			raw := p.captureUntilLoadClauseEnd()
			desc.Where = raw

		default:
			// Unknown modifier — stop parsing this desc
			goto done
		}
	}
done:
	desc.Loc = ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}
	return desc, nil
}

// captureParenthesized reads '(' ... ')' and returns the raw text including parens.
func (p *Parser) captureParenthesized() (string, error) {
	if _, err := p.expect(int('(')); err != nil {
		return "", err
	}
	var sb strings.Builder
	sb.WriteByte('(')
	depth := 1
	for depth > 0 && p.cur.Kind != tokEOF {
		switch p.cur.Kind {
		case int('('):
			depth++
			sb.WriteString(p.cur.Str)
			if p.cur.Str == "" {
				sb.WriteByte('(')
			}
		case int(')'):
			depth--
			if depth > 0 {
				sb.WriteByte(')')
			}
		default:
			if p.cur.Str != "" {
				sb.WriteString(p.cur.Str)
			}
		}
		p.advance()
	}
	sb.WriteByte(')')
	return sb.String(), nil
}

// captureUntilLoadClauseEnd captures raw tokens until the next LOAD clause
// boundary (PARTITION, COLUMNS, SET, WHERE, or end-of-desc markers).
func (p *Parser) captureUntilLoadClauseEnd() string {
	var sb strings.Builder
	for p.cur.Kind != tokEOF && p.cur.Kind != int(',') && p.cur.Kind != int(')') {
		switch p.cur.Kind {
		case kwPARTITION, kwPARTITIONS, kwCOLUMNS, kwSET, kwWHERE:
			return sb.String()
		}
		if sb.Len() > 0 {
			sb.WriteByte(' ')
		}
		if p.cur.Str != "" {
			sb.WriteString(p.cur.Str)
		} else {
			sb.WriteString(TokenName(p.cur.Kind))
		}
		p.advance()
	}
	return sb.String()
}

// skipParenthesized skips a '(' ... ')' block without capturing content.
func (p *Parser) skipParenthesized() error {
	if _, err := p.expect(int('(')); err != nil {
		return err
	}
	depth := 1
	for depth > 0 && p.cur.Kind != tokEOF {
		switch p.cur.Kind {
		case int('('):
			depth++
		case int(')'):
			depth--
		}
		p.advance()
	}
	return nil
}

// ---------------------------------------------------------------------------
// EXPORT TABLE (T6.1)
// ---------------------------------------------------------------------------

// parseExport parses:
//
//	EXPORT TABLE name [PARTITION(p1, p2, ...)] [WHERE expr]
//	    TO 'path'
//	    [PROPERTIES(...)]
//	    [WITH BROKER name [(key=val, ...)]]
//
// The EXPORT keyword has NOT yet been consumed when this is called.
func (p *Parser) parseExport() (ast.Node, error) {
	exportTok := p.advance() // consume EXPORT

	stmt := &ast.ExportStmt{
		Loc: ast.Loc{Start: exportTok.Loc.Start},
	}

	// Required TABLE keyword
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}

	// Table name
	target, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Target = target

	// Optional PARTITION(p1, p2, ...)
	if p.cur.Kind == kwPARTITION || p.cur.Kind == kwPARTITIONS {
		p.advance() // consume PARTITION / PARTITIONS
		parts, err := p.parsePartitionNameList()
		if err != nil {
			return nil, err
		}
		stmt.Partition = parts
	}

	// Optional WHERE expr
	if p.cur.Kind == kwWHERE {
		p.advance() // consume WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}

	// Required TO 'path'
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	if p.cur.Kind != tokString {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Path = p.cur.Str
	p.advance()

	// Optional PROPERTIES(...)
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
	}

	// Optional WITH BROKER broker_name [(key=val, ...)]
	if p.cur.Kind == kwWITH {
		p.advance() // consume WITH
		if _, err := p.expect(kwBROKER); err != nil {
			return nil, err
		}
		brokerName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.BrokerName = brokerName

		// Optional broker properties (...)
		if p.cur.Kind == int('(') {
			if err := p.skipParenthesized(); err != nil {
				return nil, err
			}
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
