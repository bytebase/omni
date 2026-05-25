package parser

import (
	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// STORAGE VAULT
// ---------------------------------------------------------------------------

// parseCreateStorageVault parses:
//
//	CREATE STORAGE VAULT [IF NOT EXISTS] name PROPERTIES(...)
//
// CREATE and STORAGE have already been consumed; cur is VAULT.
func (p *Parser) parseCreateStorageVault(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume VAULT

	stmt := &ast.CreateStorageVaultStmt{}

	// Optional IF NOT EXISTS
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

// parseAlterStorageVault parses:
//
//	ALTER STORAGE VAULT name PROPERTIES(...)
//
// ALTER and STORAGE have already been consumed; cur is VAULT.
func (p *Parser) parseAlterStorageVault(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume VAULT

	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	endLoc := nameLoc

	stmt := &ast.AlterStorageVaultStmt{Name: name}

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

// parseDropStorageVault parses:
//
//	DROP STORAGE VAULT [IF EXISTS] name
//
// DROP and STORAGE have already been consumed; cur is VAULT.
func (p *Parser) parseDropStorageVault(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume VAULT

	stmt := &ast.DropStorageVaultStmt{}

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

// parseSetDefaultStorageVault parses:
//
//	SET DEFAULT STORAGE VAULT name
//
// SET has already been consumed; cur is DEFAULT.
func (p *Parser) parseSetDefaultStorageVault(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume DEFAULT
	if _, err := p.expect(kwSTORAGE); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwVAULT); err != nil {
		return nil, err
	}

	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &ast.SetDefaultStorageVaultStmt{
		Name: name,
		Loc:  startLoc.Merge(nameLoc),
	}
	return stmt, nil
}

// parseUnsetDefaultStorageVault parses:
//
//	UNSET DEFAULT STORAGE VAULT
//
// UNSET has already been consumed; cur is DEFAULT.
func (p *Parser) parseUnsetDefaultStorageVault(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume DEFAULT
	if _, err := p.expect(kwSTORAGE); err != nil {
		return nil, err
	}
	vaultTok, err := p.expect(kwVAULT)
	if err != nil {
		return nil, err
	}

	stmt := &ast.UnsetDefaultStorageVaultStmt{
		Loc: startLoc.Merge(vaultTok.Loc),
	}
	return stmt, nil
}

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
//	    WITH {S3 | HDFS | BROKER broker_name}
//	    ON LOCATION "uri"
//	    PROPERTIES(...)
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

	// WITH S3 | HDFS | BROKER broker_name
	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}

	switch p.cur.Kind {
	case kwS3:
		stmt.Type = "S3"
		p.advance()
	case kwHDFS:
		stmt.Type = "HDFS"
		p.advance()
	case kwBROKER:
		stmt.Type = "BROKER"
		p.advance()
		// broker_name is an identifier following BROKER
		brokerName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.BrokerName = brokerName
	default:
		return nil, p.syntaxErrorAtCur()
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
// STAGE
// ---------------------------------------------------------------------------

// parseCreateStage parses:
//
//	CREATE STAGE [IF NOT EXISTS] name PROPERTIES(...)
//
// CREATE has already been consumed; cur is STAGE.
func (p *Parser) parseCreateStage(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume STAGE

	stmt := &ast.CreateStageStmt{}

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

// parseDropStage parses:
//
//	DROP STAGE [IF EXISTS] name
//
// DROP has already been consumed; cur is STAGE.
func (p *Parser) parseDropStage(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume STAGE

	stmt := &ast.DropStageStmt{}

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

// ---------------------------------------------------------------------------
// Shared dispatch helpers called from parser.go
// ---------------------------------------------------------------------------

// parseCreateStorage dispatches on the VAULT|POLICY token after STORAGE.
// CREATE and STORAGE have already been consumed; cur is VAULT or POLICY.
func (p *Parser) parseCreateStorage(startLoc ast.Loc) (ast.Node, error) {
	switch p.cur.Kind {
	case kwVAULT:
		return p.parseCreateStorageVault(startLoc)
	case kwPOLICY:
		return p.parseCreateStoragePolicy(startLoc)
	default:
		return p.unsupported("CREATE STORAGE")
	}
}

// parseAlterStorage dispatches on VAULT|POLICY after STORAGE.
// ALTER and STORAGE have already been consumed; cur is VAULT or POLICY.
func (p *Parser) parseAlterStorage(startLoc ast.Loc) (ast.Node, error) {
	switch p.cur.Kind {
	case kwVAULT:
		return p.parseAlterStorageVault(startLoc)
	case kwPOLICY:
		return p.parseAlterStoragePolicy(startLoc)
	default:
		return p.unsupported("ALTER STORAGE")
	}
}

// parseDropStorage dispatches on VAULT|POLICY after STORAGE.
// DROP and STORAGE have already been consumed; cur is VAULT or POLICY.
func (p *Parser) parseDropStorage(startLoc ast.Loc) (ast.Node, error) {
	switch p.cur.Kind {
	case kwVAULT:
		return p.parseDropStorageVault(startLoc)
	case kwPOLICY:
		return p.parseDropStoragePolicy(startLoc)
	default:
		return p.unsupported("DROP STORAGE")
	}
}
