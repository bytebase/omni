package parser

import (
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
)

// parseAlterSystemStmt parses alterSystemStatement (StarRocks.g4:746):
//
//	ALTER SYSTEM alterClause
//
// ALTER has been consumed; cur is SYSTEM. alterClause (StarRocks.g4:874)
// holds the cluster clauses, read here into Verb, Object and the rest, and
// every table and partition clause too, read by parseAlterClause into
// Clause. An empty optimizeClause makes a bare ALTER SYSTEM a statement.
func (p *Parser) parseAlterSystemStmt(alterLoc ast.Loc) (ast.Node, error) {
	cp := p.save()
	p.advance() // consume SYSTEM
	stmt := &ast.SystemAlterStmt{}
	ok, err := p.parseAlterSystemClause(stmt)
	if err == nil && !ok {
		stmt.Clause, err = p.parseAlterClause()
	}
	if err != nil {
		if p.strictTrailing {
			return nil, err
		}
		// Best-effort parsing keeps the loose reader for what the grammar
		// does not derive.
		p.restore(cp)
		return p.parseSystemAlter(alterLoc)
	}
	stmt.Loc = ast.Loc{Start: alterLoc.Start, End: p.prev.Loc.End}
	return stmt, nil
}

// parseAlterSystemClause reads the cluster alternatives of alterClause
// (StarRocks.g4:874, the "Alter system clause" group). It reports false,
// consuming nothing, when cur does not begin one; DROP ALL TEMPORARY
// PARTITIONS, which shares its first two words with DROP ALL BROKER, it
// reads into stmt.Clause.
func (p *Parser) parseAlterSystemClause(stmt *ast.SystemAlterStmt) (bool, error) {
	verb := strings.ToUpper(p.cur.Str)
	var err error
	switch {
	case p.cur.Kind == kwADD && (p.alterNextAt("FOLLOWER") || p.alterNextAt("OBSERVER")),
		p.cur.Kind == kwDROP && (p.alterNextAt("FOLLOWER") || p.alterNextAt("OBSERVER")):
		// StarRocks.g4:931 addFrontendClause / StarRocks.g4:935
		// dropFrontendClause: (ADD | DROP) (FOLLOWER | OBSERVER) string
		p.advance()
		stmt.Object = strings.ToUpper(p.advance().Str)
		var tok Token
		if tok, err = p.alterExpectString(); err == nil {
			stmt.Hosts = []string{tok.Str}
		}
	case p.cur.Kind == kwADD && p.alterNextAt("BACKEND"):
		// StarRocks.g4:943 addBackendClause:
		// ADD BACKEND string (',' string)* (INTO WAREHOUSE identifierOrString)?
		p.advance()
		p.advance()
		stmt.Object = "BACKEND"
		err = p.parseAlterSystemHosts(stmt, kwINTO)
	case p.cur.Kind == kwDROP && p.alterNextAt("BACKEND"):
		// StarRocks.g4:947 dropBackendClause:
		// DROP BACKEND string (',' string)* (FROM WAREHOUSE identifierOrString)? FORCE?
		p.advance()
		p.advance()
		stmt.Object = "BACKEND"
		if err = p.parseAlterSystemHosts(stmt, kwFROM); err == nil && p.cur.Kind == kwFORCE {
			p.advance()
			stmt.Force = true
		}
	case p.cur.Kind == kwADD && p.alterNextAt("COMPUTE"),
		p.cur.Kind == kwDROP && p.alterNextAt("COMPUTE"):
		// StarRocks.g4:960 addComputeNodeClause / StarRocks.g4:964
		// dropComputeNodeClause:
		// (ADD | DROP) COMPUTE NODE string (',' string)* ((INTO | FROM) WAREHOUSE identifierOrString)?
		add := p.cur.Kind == kwADD
		p.advance()
		p.advance()
		if _, err = p.alterExpectWord("NODE"); err == nil {
			stmt.Object = "COMPUTE NODE"
			prep := kwFROM
			if add {
				prep = kwINTO
			}
			err = p.parseAlterSystemHosts(stmt, prep)
		}
	case p.cur.Kind == kwADD && p.alterNextAt("BROKER"),
		p.cur.Kind == kwDROP && p.alterNextAt("BROKER"):
		// StarRocks.g4:968 modifyBrokerClause:
		// (ADD | DROP) BROKER identifierOrString string (',' string)*
		p.advance()
		p.advance()
		stmt.Object = "BROKER"
		if stmt.BrokerName, _, err = p.parseIdentifierOrString(); err == nil {
			stmt.Hosts, err = p.alterStringList()
		}
	case p.cur.Kind == kwDROP && p.peekNext().Kind == kwALL:
		p.advance() // consume DROP; cur is ALL
		if p.peekNext().Kind == kwTEMPORARY {
			// dropPartitionClause DROP ALL TEMPORARY PARTITIONS FORCE?
			// (dropPartitionClause, StarRocks.g4:1157), a table clause.
			action := &ast.AlterTableAction{Loc: p.prev.Loc}
			if err := p.parseAlterDropPartition(action); err != nil {
				return true, err
			}
			action.Loc.End = p.prev.Loc.End
			stmt.Clause = action
			return true, nil
		}
		// modifyBrokerClause (StarRocks.g4:968): DROP ALL BROKER identifierOrString
		p.advance() // consume ALL
		stmt.Object = "BROKER"
		stmt.DropAll = true
		if _, err = p.alterExpectWord("BROKER"); err == nil {
			stmt.BrokerName, _, err = p.parseIdentifierOrString()
		}
	case p.cur.Kind == kwDECOMMISSION && p.alterNextAt("BACKEND"):
		// StarRocks.g4:951 decommissionBackendClause:
		// DECOMMISSION BACKEND string (',' string)*
		p.advance()
		p.advance()
		stmt.Object = "BACKEND"
		stmt.Hosts, err = p.alterStringList()
	case p.cur.Kind == kwMODIFY && (p.alterNextAt("FRONTEND") || p.alterNextAt("BACKEND")):
		// StarRocks.g4:939 modifyFrontendHostClause:
		//	MODIFY FRONTEND HOST string TO string
		// StarRocks.g4:955 modifyBackendClause:
		//	MODIFY BACKEND HOST string TO string
		//	MODIFY BACKEND string SET propertyList
		p.advance()
		stmt.Object = strings.ToUpper(p.advance().Str)
		if stmt.Object == "BACKEND" && !p.alterAt("HOST") {
			if stmt.Hosts, err = p.alterStringList(); err == nil {
				if _, err = p.expect(kwSET); err == nil {
					stmt.SetClause, err = p.parsePropertyList()
				}
			}
			break
		}
		if _, err = p.alterExpectWord("HOST"); err != nil {
			break
		}
		var tok Token
		if tok, err = p.alterExpectString(); err != nil {
			break
		}
		stmt.Hosts = []string{tok.Str}
		if _, err = p.expect(kwTO); err != nil {
			break
		}
		if tok, err = p.alterExpectString(); err == nil {
			stmt.NewHost = tok.Str
		}
	case p.cur.Kind == kwSET && p.peekNext().Kind == kwLOAD:
		// StarRocks.g4:974 alterLoadErrorUrlClause: SET LOAD ERRORS HUB properties?
		p.advance()
		p.advance()
		stmt.Object = "LOAD ERROR HUB"
		if _, err = p.alterExpectWord("ERRORS"); err == nil {
			if _, err = p.alterExpectWord("HUB"); err == nil && p.cur.Kind == kwPROPERTIES {
				stmt.Properties, err = p.parseProperties()
			}
		}
	case p.cur.Kind == kwCREATE && p.alterNextAt("IMAGE"):
		// StarRocks.g4:978 createImageClause: CREATE IMAGE
		p.advance()
		p.advance()
		stmt.Object = "IMAGE"
	case p.alterAt("CLEAN") && p.alterNextAt("TABLET"):
		// StarRocks.g4:982 cleanTabletSchedQClause: CLEAN TABLET SCHEDULER QUEUE
		p.advance()
		p.advance()
		stmt.Object = "TABLET SCHEDULER QUEUE"
		if _, err = p.alterExpectWord("SCHEDULER"); err == nil {
			_, err = p.alterExpectWord("QUEUE")
		}
	default:
		// decommissionDiskClause, cancelDecommissionDiskClause,
		// disableDiskClause and cancelDisableDiskClause (StarRocks.g4:986-998)
		// are grammatical, but the AST builder refuses each (ALTER SYSTEM
		// DISABLE DISK 'a' ON BACKEND 'a' is a reject), so they are left to
		// fail as table clauses.
		return false, nil
	}
	stmt.Verb = verb
	return true, err
}

// parseAlterSystemHosts reads string (',' string)* and an optional
// (INTO | FROM) WAREHOUSE identifierOrString, prep naming the preposition.
func (p *Parser) parseAlterSystemHosts(stmt *ast.SystemAlterStmt, prep int) error {
	hosts, err := p.alterStringList()
	if err != nil {
		return err
	}
	stmt.Hosts = hosts
	if p.cur.Kind == prep {
		p.advance()
		if _, err := p.alterExpectWord("WAREHOUSE"); err != nil {
			return err
		}
		if stmt.Warehouse, _, err = p.parseIdentifierOrString(); err != nil {
			return err
		}
	}
	return nil
}
