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

// parseCopyInto parses Doris's copyInto (DorisParser.g4:936):
//
//	COPY INTO name=multipartIdentifier columns=identifierList? FROM
//	    (stageAndPattern | (LEFT_PAREN SELECT selectColumnClause
//	        FROM stageAndPattern whereClause? RIGHT_PAREN))
//	    properties=propertyClause?
//
// The source is a stage and what it names — @stage (stageAndPattern, :942) or
// a query of the stage's files — never the bare word it looks like the CUP
// parser's copy_from_param would take there: COPY INTO a FROM a is refused
// (compat.tsv's ..._FROM_(s...!drop-5 rows).
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

	target, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Target = target

	// identifierList? — the columns the load fills, in the target's order.
	if p.cur.Kind == int('(') {
		columns, err := p.parseLoadIdentifierList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = columns
	}

	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}

	if p.cur.Kind == int('(') {
		// (SELECT selectColumnClause FROM stageAndPattern whereClause?)
		p.advance() // consume '('
		if _, err := p.expect(kwSELECT); err != nil {
			return nil, err
		}
		items, err := p.parseCopySelectList()
		if err != nil {
			return nil, err
		}
		stmt.SelectItems = items
		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}
		if stmt.Stage, stmt.Pattern, err = p.parseCopyStage(); err != nil {
			return nil, err
		}
		if p.cur.Kind == kwWHERE {
			p.advance() // consume WHERE
			if stmt.Where, err = p.parseExpr(); err != nil {
				return nil, err
			}
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
	} else {
		if stmt.Stage, stmt.Pattern, err = p.parseCopyStage(); err != nil {
			return nil, err
		}
	}

	// propertyClause? — PROPERTIES(...), the form Doris configures a load with.
	if p.cur.Kind == kwPROPERTIES {
		if stmt.Properties, err = p.parseProperties(); err != nil {
			return nil, err
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseCopySelectList parses the selectColumnClause of a COPY INTO source
// (DorisParser.g4:1215), as the legacy parser's copy_select_expr_list
// (DorisLegacyParser.g4:625) also takes it. Doris parses a COPY INTO with both
// grammars, so the clause is a bare '*' or a list of expressions: the grammar
// takes an alias, and COPY INTO a FROM (SELECT 1 zzz FROM @a) is refused
// (compat.tsv's copy_select_expr_list/expr_list!insert row).
func (p *Parser) parseCopySelectList() ([]*ast.SelectItem, error) {
	if p.cur.Kind == int('*') {
		star := p.advance()
		return []*ast.SelectItem{{Star: true, Loc: star.Loc}}, nil
	}
	var items []*ast.SelectItem
	for {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		items = append(items, &ast.SelectItem{Expr: expr, Loc: ast.NodeLoc(expr)})
		if p.cur.Kind != int(',') {
			return items, nil
		}
		p.advance() // consume ','
	}
}

// parseCopyStage parses a stageAndPattern (DorisParser.g4:942):
//
//	ATSIGN (stage=identifier | TILDE) (LEFT_PAREN pattern=STRING_LITERAL RIGHT_PAREN)?
//
// returning the stage as written ("@stage" or "@~") and the glob it takes,
// or "" for one without a pattern. cur must be on the '@'.
func (p *Parser) parseCopyStage() (stage, pattern string, err error) {
	if p.cur.Kind != int('@') {
		return "", "", p.syntaxErrorAtCur()
	}
	p.advance() // consume '@'

	if p.cur.Kind == int('~') {
		p.advance()
		stage = "@~"
	} else {
		name, _, err := p.parseIdentifier()
		if err != nil {
			return "", "", err
		}
		stage = "@" + name
	}

	if p.cur.Kind == int('(') {
		p.advance() // consume '('
		if p.cur.Kind != tokString {
			return "", "", p.syntaxErrorAtCur()
		}
		pattern = p.advance().Str
		if _, err := p.expect(int(')')); err != nil {
			return "", "", err
		}
	}
	return stage, pattern, nil
}

// parseStringList parses '(' STRING_LITERAL (',' STRING_LITERAL)* ')': at
// least one string, comma-separated, as the grammars' file lists
// (DorisParser.g4:984, the CUP parser's string_list) have it.
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

// parseLoad parses the two LOAD statements of DorisParser.g4:
//
//	LOAD LABEL lableName=multipartIdentifier
//	    LEFT_PAREN dataDesc (COMMA dataDesc)* RIGHT_PAREN
//	    withRemoteStorageSystem? propertyClause? commentSpec?          (:153)
//	LOAD mysqlDataDesc (PROPERTIES LEFT_PAREN propertyItemList RIGHT_PAREN)?
//	    commentSpec?                                                   (:378)
//
// LOAD DATA is an unsupportedStatement, so Doris also parses it with the CUP
// parser's load_stmt (DorisLegacyParser.g4:588), and reads it here as both
// do. LOAD LABEL builds a LoadCommand, which hands its text to the CUP parser
// only for a broker load when it runs (testdata/README.md): Doris takes
// PROPERTIES (a = a) in a data description, which the CUP parser's string
// pairs refuse, but not WITH BROKER b (k = 'v').
//
// The LOAD keyword has NOT yet been consumed when this is called.
func (p *Parser) parseLoad() (ast.Node, error) {
	loadTok := p.advance() // consume LOAD
	stmt := &ast.LoadDataStmt{Loc: ast.Loc{Start: loadTok.Loc.Start}}

	if p.cur.Kind == kwDATA {
		desc, err := p.parseMysqlDataDesc()
		if err != nil {
			return nil, err
		}
		stmt.DataDescs = []*ast.LoadDataDesc{desc}
		if p.cur.Kind == kwPROPERTIES {
			if stmt.Properties, err = p.parseProperties(); err != nil {
				return nil, err
			}
		}
		if err := p.parseLoadComment(stmt); err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}

	if _, err := p.expect(kwLABEL); err != nil {
		return nil, err
	}
	label, labelParts, err := p.parseLoadLabel()
	if err != nil {
		return nil, err
	}
	stmt.Label = label
	stmt.LabelParts = labelParts

	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
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

	if p.cur.Kind == kwWITH {
		if err := p.parseLoadStorage(stmt); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind == kwPROPERTIES {
		if stmt.Properties, err = p.parseProperties(); err != nil {
			return nil, err
		}
	}
	if err := p.parseLoadComment(stmt); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseLoadComment reads an optional COMMENT STRING_LITERAL.
func (p *Parser) parseLoadComment(stmt *ast.LoadDataStmt) error {
	if p.cur.Kind != kwCOMMENT {
		return nil
	}
	p.advance() // consume COMMENT
	if p.cur.Kind != tokString {
		return p.syntaxErrorAtCur()
	}
	stmt.Comment = p.advance().Str
	return nil
}

// parseLoadStorage reads withRemoteStorageSystem for LOAD LABEL and stores it
// on stmt. The engine holds a broker load's list to strings on both sides, as
// the CUP parser's key_value_map does: (username = "u"), ("k" = 1) and
// ("k" = v) are refused by its legacy parser ("Syntax error ... Encountered:
// IDENTIFIER"), while EXPORT's list takes all three (container-verified both
// ways). So this reads the list with parseStringProperty, and parseExport
// reads it with parsePropertyItem.
func (p *Parser) parseLoadStorage(stmt *ast.LoadDataStmt) error {
	storageType, name, props, err := p.parseRemoteStorage(p.parseStringProperty)
	if err != nil {
		return err
	}
	stmt.StorageType = storageType
	stmt.BrokerName = name
	stmt.BrokerProperties = props
	return nil
}

// parseRemoteStorage reads withRemoteStorageSystem (DorisParser.g4:1121),
// which LOAD LABEL and EXPORT TABLE take:
//
//	WITH RESOURCE identifierOrText (LEFT_PAREN propertyItemList RIGHT_PAREN)?
//	| WITH (S3 | HDFS | LOCAL) LEFT_PAREN propertyItemList RIGHT_PAREN
//	| WITH BROKER identifierOrText (LEFT_PAREN propertyItemList RIGHT_PAREN)?
//
// A S3, HDFS or LOCAL clause needs its list; RESOURCE and BROKER name a
// resource or broker and may leave it out. cur must be kwWITH on entry; it is
// consumed here.
//
// parseItem reads one item of the list: the clauses differ in what they take
// for a key and a value, not in the shape of the list.
func (p *Parser) parseRemoteStorage(parseItem func() (*ast.Property, error)) (storageType, name string, props []*ast.Property, err error) {
	p.advance() // consume WITH
	needList := false
	switch p.cur.Kind {
	case kwRESOURCE, kwBROKER:
		storageType = strings.ToUpper(p.advance().Str)
		if name, _, err = p.parseIdentifierOrString(); err != nil {
			return "", "", nil, err
		}
	case kwS3, kwHDFS, kwLOCAL:
		storageType = strings.ToUpper(p.advance().Str)
		needList = true
	default:
		return "", "", nil, p.syntaxErrorAtCur()
	}
	if p.cur.Kind != int('(') {
		if needList {
			return "", "", nil, p.syntaxErrorAtCur()
		}
		return storageType, name, nil, nil
	}
	p.advance() // consume '('
	props, err = p.parsePropertyItems(parseItem)
	if err != nil {
		return "", "", nil, err
	}
	return storageType, name, props, nil
}

// parseLoadLabel parses the load label, which may be db_name.label_name or
// just label_name. It returns the joined label and its parts, so a caller can
// tell the database qualifier from a dot inside a quoted part.
func (p *Parser) parseLoadLabel() (string, []string, error) {
	first, _, err := p.parseIdentifier()
	if err != nil {
		return "", nil, err
	}
	if p.cur.Kind == int('.') {
		p.advance() // consume '.'
		second, _, err := p.parseIdentifier()
		if err != nil {
			return "", nil, err
		}
		return first + "." + second, []string{first, second}, nil
	}
	return first, []string{first}, nil
}

// parseLoadMergeType reads ((WITH)? mergeType)? before DATA (DorisParser.g4:984).
func (p *Parser) parseLoadMergeType(desc *ast.LoadDataDesc) error {
	withTok := p.cur.Kind == kwWITH
	if withTok {
		p.advance() // consume WITH
	}
	switch p.cur.Kind {
	case kwAPPEND, kwDELETE, kwMERGE:
		desc.MergeType = strings.ToUpper(p.advance().Str)
	default:
		if withTok {
			return p.syntaxErrorAtCur()
		}
	}
	return nil
}

// parseLoadDataDesc parses one dataDesc of LOAD LABEL (DorisParser.g4:983):
//
//	((WITH)? mergeType)? DATA INFILE LEFT_PAREN STRING_LITERAL (COMMA STRING_LITERAL)* RIGHT_PAREN
//	    INTO TABLE identifier partitionSpec?
//	    (COLUMNS TERMINATED BY STRING_LITERAL)? (LINES TERMINATED BY STRING_LITERAL)?
//	    (FORMAT AS identifierOrText)? (COMPRESS_TYPE AS identifierOrText)?
//	    identifierList? colFromPath? colMappingList? preFilterClause?
//	    whereClause? deleteOnClause? sequenceColClause? propertyClause?
//	| ((WITH)? mergeType)? DATA FROM TABLE identifier INTO TABLE identifier
//	    (PARTITION identifierList)? colMappingList? whereClause?
//	    deleteOnClause? propertyClause?
func (p *Parser) parseLoadDataDesc() (*ast.LoadDataDesc, error) {
	desc := &ast.LoadDataDesc{}
	startLoc := p.cur.Loc
	if err := p.parseLoadMergeType(desc); err != nil {
		return nil, err
	}
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

	if fromTable {
		if p.cur.Kind == kwPARTITION {
			p.advance() // consume PARTITION
			if desc.Partition, err = p.parseLoadIdentifierList(); err != nil {
				return nil, err
			}
		}
	} else {
		if err := p.parseLoadPartitionSpec(desc); err != nil {
			return nil, err
		}
		if desc.ColumnSeparator, err = p.parseTerminatedBy(kwCOLUMNS); err != nil {
			return nil, err
		}
		if desc.LineDelimiter, err = p.parseTerminatedBy(kwLINES); err != nil {
			return nil, err
		}
		if desc.Format, err = p.parseLoadAsName(kwFORMAT); err != nil {
			return nil, err
		}
		if desc.CompressType, err = p.parseLoadAsName(kwCOMPRESS_TYPE); err != nil {
			return nil, err
		}
		if p.cur.Kind == int('(') {
			if desc.ColumnList, err = p.parseLoadIdentifierList(); err != nil {
				return nil, err
			}
		}
		if p.cur.Kind == kwCOLUMNS {
			// colFromPath (:1109): COLUMNS FROM PATH AS identifierList
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
		// colMappingList (:1113): SET LEFT_PAREN mappingExpr (COMMA mappingExpr)* RIGHT_PAREN
		p.advance() // consume SET
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		if desc.ColumnMappings, err = p.parseLoadMappings(); err != nil {
			return nil, err
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
	}
	if !fromTable && p.cur.Kind == kwPRECEDING {
		p.advance() // consume PRECEDING
		if _, err := p.expect(kwFILTER); err != nil {
			return nil, err
		}
		if desc.PrecedingFilter, err = p.parseExpr(); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind == kwWHERE {
		p.advance() // consume WHERE
		if desc.Where, err = p.parseExpr(); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind == kwDELETE {
		p.advance() // consume DELETE
		if _, err := p.expect(kwON); err != nil {
			return nil, err
		}
		if desc.DeleteOn, err = p.parseExpr(); err != nil {
			return nil, err
		}
	}
	if !fromTable && p.cur.Kind == kwORDER {
		p.advance() // consume ORDER
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		if desc.SequenceColumn, _, err = p.parseIdentifier(); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind == kwPROPERTIES {
		// propertyClause (:1377), in the grammar's full item form.
		p.advance() // consume PROPERTIES
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		if desc.Properties, err = p.parsePropertyItems(p.parsePropertyItem); err != nil {
			return nil, err
		}
	}

	desc.Loc = ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}
	return desc, nil
}

// parseMysqlDataDesc parses LOAD DATA's mysqlDataDesc (DorisParser.g4:1142)
// as the CUP parser's mysql_data_desc (DorisLegacyParser.g4:663) also
// takes it:
//
//	DATA LOCAL? INFILE STRING_LITERAL INTO TABLE multipartIdentifier
//	    (PARTITION identifierList)?
//	    (COLUMNS TERMINATED BY STRING_LITERAL)? (LINES TERMINATED BY STRING_LITERAL)?
//	    (IGNORE INTEGER_VALUE (LINES | ROWS))? identifierList? colMappingList?
//	    propertyClause?
//
// whose properties are the CUP parser's 'k' = 'v' pairs.
func (p *Parser) parseMysqlDataDesc() (*ast.LoadDataDesc, error) {
	desc := &ast.LoadDataDesc{}
	startLoc := p.advance().Loc // consume DATA
	if p.cur.Kind == kwLOCAL {
		desc.Local = true
		p.advance()
	}
	if _, err := p.expect(kwINFILE); err != nil {
		return nil, err
	}
	if p.cur.Kind != tokString {
		return nil, p.syntaxErrorAtCur()
	}
	desc.SourceFiles = []string{p.advance().Str}
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
	if p.cur.Kind == kwPARTITION {
		p.advance() // consume PARTITION
		if desc.Partition, err = p.parseLoadIdentifierList(); err != nil {
			return nil, err
		}
	}
	if desc.ColumnSeparator, err = p.parseTerminatedBy(kwCOLUMNS); err != nil {
		return nil, err
	}
	if desc.LineDelimiter, err = p.parseTerminatedBy(kwLINES); err != nil {
		return nil, err
	}
	if p.cur.Kind == kwIGNORE {
		// skipLines (:1155): IGNORE INTEGER_VALUE (LINES | ROWS)
		p.advance() // consume IGNORE
		if p.cur.Kind != tokInt {
			return nil, p.syntaxErrorAtCur()
		}
		desc.SkipLines = int(p.advance().Ival)
		if p.cur.Kind != kwLINES && p.cur.Kind != kwROWS {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
	}
	if p.cur.Kind == int('(') {
		if desc.ColumnList, err = p.parseLoadIdentifierList(); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind == kwSET {
		p.advance() // consume SET
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		if desc.ColumnMappings, err = p.parseLoadMappings(); err != nil {
			return nil, err
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind == kwPROPERTIES {
		if desc.Properties, err = p.parseProperties(); err != nil {
			return nil, err
		}
	}
	desc.Loc = ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}
	return desc, nil
}

// parseLoadTableName reads a data description's table: an identifier, which
// both grammars take (DorisParser.g4:986; the CUP parser's ident).
func (p *Parser) parseLoadTableName() (*ast.ObjectName, error) {
	name, loc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	return &ast.ObjectName{Parts: []string{name}, Loc: loc}, nil
}

// parseLoadPartitionSpec reads a data description's partitionSpec.
func (p *Parser) parseLoadPartitionSpec(desc *ast.LoadDataDesc) error {
	names, temporary, err := p.parsePartitionSpec()
	desc.Partition, desc.TemporaryPart = names, temporary
	return err
}

// parsePartitionSpec reads an optional partitionSpec (DorisParser.g4:962):
//
//	TEMPORARY? (PARTITION | PARTITIONS) identifierList
//	| TEMPORARY? PARTITION errorCapturingIdentifier
//	| (PARTITION | PARTITIONS) LEFT_PAREN ASTERISK RIGHT_PAREN
//
// which the CUP parser's partition_names (DorisLegacyParser.g4:2311) also
// takes. It returns the names, ["*"] for (*), and whether TEMPORARY was
// given; nil when no PARTITION or PARTITIONS is at cur.
func (p *Parser) parsePartitionSpec() ([]string, bool, error) {
	temporary := false
	if p.cur.Kind == kwTEMPORARY {
		temporary = true
		p.advance()
		if p.cur.Kind != kwPARTITION && p.cur.Kind != kwPARTITIONS {
			return nil, true, p.syntaxErrorAtCur()
		}
	}
	if p.cur.Kind != kwPARTITION && p.cur.Kind != kwPARTITIONS {
		return nil, false, nil
	}
	plural := p.advance().Kind == kwPARTITIONS
	if p.cur.Kind == int('(') && p.peekNext().Kind == int('*') && !temporary {
		p.advance() // consume '('
		p.advance() // consume '*'
		_, err := p.expect(int(')'))
		return []string{"*"}, false, err
	}
	if p.cur.Kind == int('(') || plural {
		parts, err := p.parseLoadIdentifierList()
		return parts, temporary, err
	}
	name, _, err := p.parseIdentifier()
	return []string{name}, temporary, err
}

// parseLoadIdentifierList reads identifierList: '(' identifier (',' identifier)* ')'.
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

// parseTerminatedBy reads an optional (COLUMNS | LINES) TERMINATED BY
// STRING_LITERAL whose first word is kw. COLUMNS FROM PATH is left for its
// own clause.
func (p *Parser) parseTerminatedBy(kw int) (string, error) {
	if p.cur.Kind != kw || p.peekNext().Kind != kwTERMINATED {
		return "", nil
	}
	p.advance() // consume COLUMNS / LINES
	p.advance() // consume TERMINATED
	if _, err := p.expect(kwBY); err != nil {
		return "", err
	}
	if p.cur.Kind != tokString {
		return "", p.syntaxErrorAtCur()
	}
	return p.advance().Str, nil
}

// parseLoadAsName reads an optional (FORMAT | COMPRESS_TYPE) AS
// identifierOrText whose first word is kw.
func (p *Parser) parseLoadAsName(kw int) (string, error) {
	if p.cur.Kind != kw {
		return "", nil
	}
	p.advance() // consume FORMAT / COMPRESS_TYPE
	if _, err := p.expect(kwAS); err != nil {
		return "", err
	}
	name, _, err := p.parseIdentifierOrString()
	return name, err
}

// parseLoadMappings reads mappingExpr (COMMA mappingExpr)*, where
// mappingExpr is identifier EQ expression (DorisParser.g4:1117).
func (p *Parser) parseLoadMappings() ([]*ast.Assignment, error) {
	var out []*ast.Assignment
	for {
		start := p.cur.Loc
		name, loc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(int('=')); err != nil {
			return nil, err
		}
		val, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		out = append(out, &ast.Assignment{
			Column: &ast.ObjectName{Parts: []string{name}, Loc: loc},
			Value:  val,
			Loc:    ast.Loc{Start: start.Start, End: p.prev.Loc.End},
		})
		if p.cur.Kind != int(',') {
			return out, nil
		}
		p.advance() // consume ','
	}
}

// ---------------------------------------------------------------------------
// EXPORT TABLE (T6.1)
// ---------------------------------------------------------------------------

// parseExport parses EXPORT TABLE (DorisParser.g4:159-163):
//
//	EXPORT TABLE name
//	    [PARTITION(p1, p2, ...)] [WHERE expr]
//	    TO 'path'
//	    [PROPERTIES(...)]
//	    [withRemoteStorageSystem]
//
// DorisLegacyParser.g4's export_stmt (sql_parser.cup:1660) agrees on the
// PARTITION keyword: it takes `PARTITION`, never `PARTITIONS`, so
// `EXPORT TABLE a PARTITIONS(a) TO 'a'` is refused (container-verified). Its
// opt_broker takes the same WITH forms as the grammar's
// withRemoteStorageSystem, which parseRemoteStorage reads.
//
// Both property lists carry the grammar's full propertyItem, which is wider
// than the "key"="value" pair most clauses take: (a = a) and ('a' = 'a') are
// both accepted here (container-verified).
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
	if p.cur.Kind == kwPARTITION {
		p.advance() // consume PARTITION
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
		props, err := p.parseFullProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
	}

	// Optional withRemoteStorageSystem
	if p.cur.Kind == kwWITH {
		storageType, name, props, err := p.parseRemoteStorage(p.parsePropertyItem)
		if err != nil {
			return nil, err
		}
		stmt.StorageType = storageType
		stmt.BrokerName = name
		stmt.BrokerProperties = props
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
