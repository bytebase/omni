package parser

import (
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
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

// parseStringList parses stringList (StarRocks.g4:2633):
// '(' string (',' string)* ')'.
func (p *Parser) parseStringList() ([]string, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	var strs []string
	for {
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		strs = append(strs, p.advance().Str)
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return strs, nil
}

// ---------------------------------------------------------------------------
// LOAD (T6.1)
// ---------------------------------------------------------------------------

// parseLoad parses loadStatement (StarRocks.g4:1427):
//
//	LOAD LABEL labelName dataDescList? brokerDesc?
//	    (BY identifierOrString)? (PROPERTIES propertyList)?
//	| LOAD LABEL labelName dataDescList? resourceDesc (PROPERTIES propertyList)?
//	brokerDesc   : WITH BROKER propertyList? | WITH BROKER identifierOrString propertyList?
//	resourceDesc : WITH RESOURCE identifierOrString propertyList?
//
// The LOAD keyword has NOT yet been consumed when this is called.
func (p *Parser) parseLoad() (ast.Node, error) {
	loadTok := p.advance() // consume LOAD
	stmt := &ast.LoadDataStmt{Loc: ast.Loc{Start: loadTok.Loc.Start}}

	if _, err := p.expect(kwLABEL); err != nil {
		return nil, err
	}
	label, err := p.parseLoadLabel()
	if err != nil {
		return nil, err
	}
	stmt.Label = label

	if p.cur.Kind == int('(') {
		// dataDescList (:1442): '(' dataDesc (',' dataDesc)* ')'
		p.advance() // consume '('
		for {
			desc, err := p.parseLoadDataDesc()
			if err != nil {
				return nil, err
			}
			stmt.DataDescs = append(stmt.DataDescs, desc)
			if p.cur.Kind != int(',') {
				break
			}
			p.advance() // consume ','
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
	}

	if p.cur.Kind == kwWITH {
		p.advance() // consume WITH
		switch p.cur.Kind {
		case kwBROKER:
			stmt.StorageType = "BROKER"
			p.advance()
			if p.cur.Kind != int('(') && p.cur.Kind != kwBY && p.cur.Kind != kwPROPERTIES &&
				p.cur.Kind != tokEOF && p.cur.Kind != int(';') {
				if stmt.BrokerName, _, err = p.parseIdentifierOrString(); err != nil {
					return nil, err
				}
			}
		case kwRESOURCE:
			stmt.StorageType = "RESOURCE"
			p.advance()
			if stmt.BrokerName, _, err = p.parseIdentifierOrString(); err != nil {
				return nil, err
			}
		default:
			return nil, p.syntaxErrorAtCur()
		}
		if p.cur.Kind == int('(') {
			if stmt.BrokerProperties, err = p.parsePropertyList(); err != nil {
				return nil, err
			}
		}
	}
	if p.cur.Kind == kwBY && stmt.StorageType != "RESOURCE" {
		p.advance() // consume BY
		if stmt.System, _, err = p.parseIdentifierOrString(); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind == kwPROPERTIES {
		if stmt.Properties, err = p.parseProperties(); err != nil {
			return nil, err
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseLoadLabel parses labelName (StarRocks.g4:1438): (identifier '.')? identifier.
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

// parseLoadDataDesc parses one dataDesc (StarRocks.g4:1447):
//
//	DATA INFILE stringList NEGATIVE? INTO TABLE identifier partitionNames?
//	    (COLUMNS TERMINATED BY string)? (ROWS TERMINATED BY string)?
//	    fileFormat? formatProps? columnAliases?
//	    (COLUMNS FROM PATH AS identifierList)? (SET classifier)? (WHERE expression)?
//	| DATA FROM TABLE identifier NEGATIVE? INTO TABLE identifier partitionNames?
//	    (SET classifier)? (WHERE expression)?
func (p *Parser) parseLoadDataDesc() (*ast.LoadDataDesc, error) {
	desc := &ast.LoadDataDesc{}
	startLoc := p.cur.Loc
	if _, err := p.expect(kwDATA); err != nil {
		return nil, err
	}

	fromTable := p.cur.Kind == kwFROM
	if fromTable {
		p.advance() // consume FROM
		if _, err := p.expect(kwTABLE); err != nil {
			return nil, err
		}
		src, err := p.parseLoadTableName()
		if err != nil {
			return nil, err
		}
		desc.SourceTable = src
	} else {
		if _, err := p.expect(kwINFILE); err != nil {
			return nil, err
		}
		files, err := p.parseStringList()
		if err != nil {
			return nil, err
		}
		desc.SourceFiles = files
	}
	if p.cur.Kind == kwNEGATIVE {
		desc.Negative = true
		p.advance()
	}
	if _, err := p.expect(kwINTO); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	target, err := p.parseLoadTableName()
	if err != nil {
		return nil, err
	}
	desc.Target = target

	if k := p.cur.Kind; k == kwTEMPORARY || k == kwPARTITION || k == kwPARTITIONS {
		var names ast.AlterTableAction
		if err := p.parseAlterPartitionNames(&names); err != nil {
			return nil, err
		}
		desc.TemporaryPart = names.Temporary
		desc.Partition = names.PartitionList
		if names.PartitionName != "" {
			desc.Partition = []string{names.PartitionName}
		}
	}

	if !fromTable {
		if desc.ColumnSeparator, err = p.parseTerminatedBy(kwCOLUMNS); err != nil {
			return nil, err
		}
		if desc.RowDelimiter, err = p.parseTerminatedBy(kwROWS); err != nil {
			return nil, err
		}
		if p.cur.Kind == kwFORMAT {
			// fileFormat (:2742): FORMAT AS (identifier | string)
			p.advance() // consume FORMAT
			if _, err := p.expect(kwAS); err != nil {
				return nil, err
			}
			if desc.Format, _, err = p.parseIdentifierOrString(); err != nil {
				return nil, err
			}
		}
		// formatProps? columnAliases?: both open with '(', and the token
		// after it decides (isLoadFormatPropsStart).
		if p.cur.Kind == int('(') {
			p.advance() // consume '('
			if p.isLoadFormatPropsStart() {
				if desc.FormatProps, err = p.parseLoadFormatProps(); err != nil {
					return nil, err
				}
				if p.cur.Kind == int('(') {
					p.advance() // consume '('
					if desc.ColumnList, err = p.parseLoadColumnAliases(); err != nil {
						return nil, err
					}
				}
			} else if desc.ColumnList, err = p.parseLoadColumnAliases(); err != nil {
				return nil, err
			}
		}
		if p.cur.Kind == kwCOLUMNS {
			p.advance() // consume COLUMNS
			for _, kw := range []int{kwFROM, kwPATH, kwAS} {
				if _, err := p.expect(kw); err != nil {
					return nil, err
				}
			}
			if desc.ColumnsFromPath, err = p.parseLoadIdentifierList(); err != nil {
				return nil, err
			}
		}
	}
	if p.cur.Kind == kwSET {
		// classifier (:1400): '(' expressionList ')'
		p.advance() // consume SET
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		if desc.ColumnMappings, err = p.parseExprList(); err != nil {
			return nil, err
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind == kwWHERE {
		p.advance() // consume WHERE
		if desc.Where, err = p.parseExpr(); err != nil {
			return nil, err
		}
	}

	desc.Loc = ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}
	return desc, nil
}

// isLoadFormatPropsStart reports whether the token after a data
// description's '(' begins formatProps (StarRocks.g4:1474) rather than
// columnAliases: ')' (every property is optional, and columnAliases needs
// one name), or SKIP_HEADER, TRIM_SPACE, ENCLOSE or ESCAPE followed by '='.
// A column may be named SKIP_HEADER, so the word alone does not decide.
func (p *Parser) isLoadFormatPropsStart() bool {
	if p.cur.Kind == int(')') {
		return true
	}
	return loadFormatPropWord(p.cur) != "" && p.peekNext().Kind == int('=')
}

// loadFormatPropWord returns the formatProps word tok spells, or "".
func loadFormatPropWord(tok Token) string {
	if tok.Kind == kwESCAPE {
		return "ESCAPE"
	}
	if tok.Kind == tokIdent {
		switch w := strings.ToUpper(tok.Str); w {
		case "SKIP_HEADER", "TRIM_SPACE", "ENCLOSE":
			return w
		}
	}
	return ""
}

// parseLoadFormatProps reads formatProps (StarRocks.g4:1474):
//
//	'(' (SKIP_HEADER '=' INTEGER_VALUE)? (TRIM_SPACE '=' booleanValue)?
//	    (ENCLOSE '=' string)? (ESCAPE '=' string)? ')'
//
// in that order, each at most once, with no commas. The '(' has been
// consumed.
func (p *Parser) parseLoadFormatProps() ([]*ast.Property, error) {
	var props []*ast.Property
	order := []string{"SKIP_HEADER", "TRIM_SPACE", "ENCLOSE", "ESCAPE"}
	next := 0
	for p.cur.Kind != int(')') {
		word := loadFormatPropWord(p.cur)
		i := next
		for i < len(order) && order[i] != word {
			i++
		}
		if word == "" || i == len(order) || p.peekNext().Kind != int('=') {
			return nil, p.syntaxErrorAtCur()
		}
		next = i + 1
		start := p.advance().Loc // consume the word
		p.advance()              // consume '='
		ok := false
		switch word {
		case "SKIP_HEADER":
			ok = p.cur.Kind == tokInt
		case "TRIM_SPACE":
			ok = p.cur.Kind == kwTRUE || p.cur.Kind == kwFALSE
		default:
			ok = p.cur.Kind == tokString
		}
		if !ok {
			return nil, p.syntaxErrorAtCur()
		}
		val := p.advance()
		value := val.Str
		if val.Kind != tokString {
			value = p.input[val.Loc.Start-p.baseOffset : val.Loc.End-p.baseOffset]
		}
		props = append(props, &ast.Property{
			Key:   word,
			Value: value,
			Loc:   ast.Loc{Start: start.Start, End: val.Loc.End},
		})
	}
	p.advance() // consume ')'
	return props, nil
}

// parseLoadColumnAliases reads the rest of columnAliases (StarRocks.g4:2263)
// after its '(': identifier (',' identifier)* ')'.
func (p *Parser) parseLoadColumnAliases() ([]string, error) {
	names, err := p.parseIdentifierList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return names, nil
}

// parseLoadTableName reads a data description's table: an identifier.
func (p *Parser) parseLoadTableName() (*ast.ObjectName, error) {
	name, loc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	return &ast.ObjectName{Parts: []string{name}, Loc: loc}, nil
}

// parseLoadIdentifierList reads '(' identifier (',' identifier)* ')'.
func (p *Parser) parseLoadIdentifierList() ([]string, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	names, err := p.parseIdentifierList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return names, nil
}

// parseTerminatedBy reads an optional (COLUMNS | ROWS) TERMINATED BY string
// whose first word is kw. COLUMNS FROM PATH is left for its own clause.
func (p *Parser) parseTerminatedBy(kw int) (string, error) {
	if p.cur.Kind != kw || p.peekNext().Kind != kwTERMINATED {
		return "", nil
	}
	p.advance() // consume COLUMNS / ROWS
	p.advance() // consume TERMINATED
	if _, err := p.expect(kwBY); err != nil {
		return "", err
	}
	if p.cur.Kind != tokString {
		return "", p.syntaxErrorAtCur()
	}
	return p.advance().Str, nil
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

		// Optional broker properties: brokerDesc's props=propertyList?, each
		// item a string on both sides.
		if p.cur.Kind == int('(') {
			props, err := p.parsePropertyList()
			if err != nil {
				return nil, err
			}
			stmt.BrokerProperties = props
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
