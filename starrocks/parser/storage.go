package parser

import (
	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// STORAGE POLICY
// ---------------------------------------------------------------------------

// parseCreateStoragePolicy parses:
//
//	CREATE STORAGE POLICY [IF NOT EXISTS] name PROPERTIES(...)
//
// CREATE and STORAGE have already been consumed; cur is POLICY.
func (p *Parser) parseCreateStoragePolicy(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume POLICY

	stmt := &ast.CreateStoragePolicyStmt{}

	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	endLoc := nameLoc

	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseAlterStoragePolicy parses:
//
//	ALTER STORAGE POLICY name PROPERTIES(...)
//
// ALTER and STORAGE have already been consumed; cur is POLICY.
func (p *Parser) parseAlterStoragePolicy(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume POLICY

	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	endLoc := nameLoc

	stmt := &ast.AlterStoragePolicyStmt{Name: name}

	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseDropStoragePolicy parses:
//
//	DROP STORAGE POLICY [IF EXISTS] name
//
// DROP and STORAGE have already been consumed; cur is POLICY.
func (p *Parser) parseDropStoragePolicy(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume POLICY

	stmt := &ast.DropStoragePolicyStmt{}

	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc = startLoc.Merge(nameLoc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// REPOSITORY
// ---------------------------------------------------------------------------

// parseCreateRepository parses:
//
//	CREATE [READ ONLY] REPOSITORY name
//	    WITH BROKER broker_name?
//	    ON LOCATION "uri"
//	    PROPERTIES(...)
//
// createRepositoryStatement (StarRocks.g4:1802) takes BROKER only. The
// Doris forms WITH S3 and WITH HDFS are not in the grammar, and the pinned
// 3.4.10 server refuses both.
//
// CREATE has already been consumed; cur is REPOSITORY (or READ).
// readOnly indicates whether READ ONLY was already consumed.
func (p *Parser) parseCreateRepository(startLoc ast.Loc, readOnly bool) (ast.Node, error) {
	p.advance() // consume REPOSITORY

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &ast.CreateRepositoryStmt{Name: name, ReadOnly: readOnly}

	// WITH BROKER broker_name?
	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwBROKER); err != nil {
		return nil, err
	}
	stmt.Type = "BROKER"
	if isIdentifierToken(p.cur.Kind) || p.cur.Kind == tokString {
		brokerName, _, err := p.parseIdentifierOrString()
		if err != nil {
			return nil, err
		}
		stmt.BrokerName = brokerName
	}

	// ON LOCATION "uri"
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwLOCATION); err != nil {
		return nil, err
	}
	// location is a string literal; consume it (we don't store it separately
	// since it will appear in PROPERTIES as "location" key in practice, but the
	// syntax places it here explicitly)
	if p.cur.Kind != tokString {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume location URI string

	endLoc := p.prev.Loc

	// PROPERTIES(...)
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseAlterRepository parses:
//
//	ALTER REPOSITORY name PROPERTIES(...)
//
// ALTER has already been consumed; cur is REPOSITORY.
func (p *Parser) parseAlterRepository(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume REPOSITORY

	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	endLoc := nameLoc

	stmt := &ast.AlterRepositoryStmt{Name: name}

	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseDropRepository parses:
//
//	DROP REPOSITORY name
//
// DROP has already been consumed; cur is REPOSITORY.
func (p *Parser) parseDropRepository(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume REPOSITORY

	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &ast.DropRepositoryStmt{
		Name: name,
		Loc:  startLoc.Merge(nameLoc),
	}
	return stmt, nil
}

// ---------------------------------------------------------------------------
// FILE
// ---------------------------------------------------------------------------

// parseCreateFile parses:
//
//	CREATE FILE file_name [IN db] PROPERTIES(...)
//
// CREATE has already been consumed; cur is FILE.
func (p *Parser) parseCreateFile(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume FILE

	// file_name is a string literal or identifier
	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	endLoc := nameLoc

	stmt := &ast.CreateFileStmt{Name: name}

	// Optional IN db_name
	if p.cur.Kind == kwIN {
		p.advance() // consume IN
		dbName, dbLoc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Database = dbName
		endLoc = dbLoc
	}

	// properties is required (StarRocks.g4 createFileStatement).
	if p.cur.Kind != kwPROPERTIES {
		return nil, p.syntaxErrorAtCur()
	}
	props, err := p.parseProperties()
	if err != nil {
		return nil, err
	}
	stmt.Properties = props
	if len(props) > 0 {
		endLoc = ast.NodeLoc(props[len(props)-1])
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseDropFile parses:
//
//	DROP FILE file_name [FROM db] PROPERTIES(...)
//
// DROP has already been consumed; cur is FILE.
func (p *Parser) parseDropFile(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume FILE

	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	endLoc := nameLoc

	stmt := &ast.DropFileStmt{Name: name}

	// Optional FROM db_name
	if p.cur.Kind == kwFROM {
		p.advance() // consume FROM
		dbName, dbLoc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Database = dbName
		endLoc = dbLoc
	}

	// properties is required (StarRocks.g4 dropFileStatement).
	if p.cur.Kind != kwPROPERTIES {
		return nil, p.syntaxErrorAtCur()
	}
	props, err := p.parseProperties()
	if err != nil {
		return nil, err
	}
	stmt.Properties = props
	if len(props) > 0 {
		endLoc = ast.NodeLoc(props[len(props)-1])
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// STORAGE VOLUME
// ---------------------------------------------------------------------------

// parseCreateStorageVolume parses createStorageVolumeStatement
// (StarRocks.g4:782):
//
//	CREATE STORAGE VOLUME [IF NOT EXISTS] name TYPE = type
//	    LOCATIONS = (uri, ...) [COMMENT 'text'] [PROPERTIES (...)]
//
// CREATE and STORAGE have already been consumed; cur is VOLUME. TYPE and
// LOCATIONS both sit outside any optional group, so a statement that leaves
// either one out is refused (the pinned server rejects
// CREATE STORAGE VOLUME a LOCATIONS = ('a') and
// CREATE STORAGE VOLUME a TYPE = a).
func (p *Parser) parseCreateStorageVolume(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume VOLUME

	stmt := &ast.CreateStorageVolumeStmt{}

	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	endLoc := nameLoc

	// typeDesc (StarRocks.g4:787): TYPE EQ identifier.
	if _, err := p.expect(kwTYPE); err != nil {
		return nil, err
	}
	if _, err := p.expect(int('=')); err != nil {
		return nil, err
	}
	typeName, typeLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Type = typeName
	endLoc = typeLoc

	// locationsDesc (StarRocks.g4:791): LOCATIONS EQ stringList.
	if _, err := p.expect(kwLOCATIONS); err != nil {
		return nil, err
	}
	if _, err := p.expect(int('=')); err != nil {
		return nil, err
	}
	locations, err := p.parseStringList()
	if err != nil {
		return nil, err
	}
	stmt.Locations = locations
	endLoc = p.prev.Loc

	// comment (StarRocks.g4:2734): COMMENT string.
	if p.cur.Kind == kwCOMMENT {
		p.advance() // consume COMMENT
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Comment = p.cur.Str
		endLoc = p.cur.Loc
		p.advance()
	}

	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		endLoc = p.prev.Loc
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseAlterStorageVolume parses alterStorageVolumeStatement
// (StarRocks.g4:803):
//
//	ALTER STORAGE VOLUME name clause (',' clause)*
//
// ALTER and STORAGE have already been consumed; cur is VOLUME.
func (p *Parser) parseAlterStorageVolume(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume VOLUME

	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}

	stmt := &ast.AlterStorageVolumeStmt{Name: name}
	endLoc := nameLoc

	for {
		clause, err := p.parseAlterStorageVolumeClause()
		if err != nil {
			return nil, err
		}
		stmt.Clauses = append(stmt.Clauses, clause)
		endLoc = clause.Loc
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseAlterStorageVolumeClause parses alterStorageVolumeClause
// (StarRocks.g4:807): modifyStorageVolumeCommentClause |
// modifyStorageVolumePropertiesClause.
func (p *Parser) parseAlterStorageVolumeClause() (*ast.AlterStorageVolumeClause, error) {
	startLoc := p.cur.Loc

	switch p.cur.Kind {
	case kwCOMMENT:
		// modifyStorageVolumeCommentClause (StarRocks.g4:816):
		// COMMENT '=' string, with the '=' obligatory.
		p.advance() // consume COMMENT
		if _, err := p.expect(int('=')); err != nil {
			return nil, err
		}
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		clause := &ast.AlterStorageVolumeClause{
			Type:    ast.StorageVolumeComment,
			Comment: p.cur.Str,
		}
		p.advance() // consume the comment string
		clause.Loc = startLoc.Merge(p.prev.Loc)
		return clause, nil

	case kwSET:
		// modifyStorageVolumePropertiesClause (StarRocks.g4:812):
		// SET propertyList.
		p.advance() // consume SET
		props, err := p.parsePropertyList()
		if err != nil {
			return nil, err
		}
		clause := &ast.AlterStorageVolumeClause{
			Type:       ast.StorageVolumeProperties,
			Properties: props,
		}
		clause.Loc = startLoc.Merge(p.prev.Loc)
		return clause, nil

	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseDropStorageVolume parses dropStorageVolumeStatement
// (StarRocks.g4:799):
//
//	DROP STORAGE VOLUME [IF EXISTS] name
//
// DROP and STORAGE have already been consumed; cur is VOLUME.
func (p *Parser) parseDropStorageVolume(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume VOLUME

	stmt := &ast.DropStorageVolumeStmt{}

	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc = startLoc.Merge(nameLoc)
	return stmt, nil
}

// parseDescStorageVolume parses descStorageVolumeStatement
// (StarRocks.g4:820):
//
//	(DESC | DESCRIBE) STORAGE VOLUME name
//
// DESC and STORAGE have been consumed and cur is VOLUME; the caller checked
// that VOLUME follows, since STORAGE and VOLUME are both nonReserved and DESC
// STORAGE alone describes a table named storage.
func (p *Parser) parseDescStorageVolume(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume VOLUME

	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}

	stmt := &ast.DescStorageVolumeStmt{Name: name, Loc: startLoc.Merge(nameLoc)}
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Shared dispatch helpers called from parser.go
// ---------------------------------------------------------------------------

// parseCreateStorage dispatches on the POLICY|VOLUME token after STORAGE.
// CREATE and STORAGE have already been consumed; cur is POLICY or VOLUME.
// StarRocks has no STORAGE VAULT: StarRocks.g4 has no VAULT token, and
// STORAGE VOLUME is its form.
func (p *Parser) parseCreateStorage(startLoc ast.Loc) (ast.Node, error) {
	switch p.cur.Kind {
	case kwPOLICY:
		return p.parseCreateStoragePolicy(startLoc)
	case kwVOLUME:
		return p.parseCreateStorageVolume(startLoc)
	default:
		return p.unsupported("CREATE STORAGE")
	}
}

// parseAlterStorage dispatches on POLICY|VOLUME after STORAGE.
// ALTER and STORAGE have already been consumed; cur is POLICY or VOLUME.
func (p *Parser) parseAlterStorage(startLoc ast.Loc) (ast.Node, error) {
	switch p.cur.Kind {
	case kwPOLICY:
		return p.parseAlterStoragePolicy(startLoc)
	case kwVOLUME:
		return p.parseAlterStorageVolume(startLoc)
	default:
		return p.unsupported("ALTER STORAGE")
	}
}

// parseDropStorage dispatches on POLICY|VOLUME after STORAGE.
// DROP and STORAGE have already been consumed; cur is POLICY or VOLUME.
func (p *Parser) parseDropStorage(startLoc ast.Loc) (ast.Node, error) {
	switch p.cur.Kind {
	case kwPOLICY:
		return p.parseDropStoragePolicy(startLoc)
	case kwVOLUME:
		return p.parseDropStorageVolume(startLoc)
	default:
		return p.unsupported("DROP STORAGE")
	}
}
