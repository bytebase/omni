package parser

import "github.com/bytebase/omni/mongo/ast"

// dbMethodNames maps keyword token types to method name strings for known database methods.
var dbMethodNames = map[int]string{
	kwGetName:                       "getName",
	kwGetSiblingDB:                  "getSiblingDB",
	kwGetMongo:                      "getMongo",
	kwGetCollectionNames:            "getCollectionNames",
	kwGetCollectionInfos:            "getCollectionInfos",
	kwGetCollection:                 "getCollection",
	kwCreateCollection:              "createCollection",
	kwCreateView:                    "createView",
	kwDropDatabase:                  "dropDatabase",
	kwAdminCommand:                  "adminCommand",
	kwRunCommand:                    "runCommand",
	kwGetProfilingStatus:            "getProfilingStatus",
	kwSetProfilingLevel:             "setProfilingLevel",
	kwGetLogComponents:              "getLogComponents",
	kwSetLogLevel:                   "setLogLevel",
	kwFsyncLock:                     "fsyncLock",
	kwFsyncUnlock:                   "fsyncUnlock",
	kwCurrentOp:                     "currentOp",
	kwKillOp:                        "killOp",
	kwGetUser:                       "getUser",
	kwGetUsers:                      "getUsers",
	kwCreateUser:                    "createUser",
	kwUpdateUser:                    "updateUser",
	kwDropUser:                      "dropUser",
	kwDropAllUsers:                  "dropAllUsers",
	kwGrantRolesToUser:              "grantRolesToUser",
	kwRevokeRolesFromUser:           "revokeRolesFromUser",
	kwGetRole:                       "getRole",
	kwGetRoles:                      "getRoles",
	kwCreateRole:                    "createRole",
	kwUpdateRole:                    "updateRole",
	kwDropRole:                      "dropRole",
	kwDropAllRoles:                  "dropAllRoles",
	kwGrantPrivilegesToRole:         "grantPrivilegesToRole",
	kwRevokePrivilegesFromRole:      "revokePrivilegesFromRole",
	kwGrantRolesToRole:              "grantRolesToRole",
	kwRevokeRolesFromRole:           "revokeRolesFromRole",
	kwServerStatus:                  "serverStatus",
	kwIsMaster:                      "isMaster",
	kwHello:                         "hello",
	kwHostInfo:                      "hostInfo",
	kwAggregate:                     "aggregate",
	kwWatch:                         "watch",
	kwStats:                         "stats",
	kwVersion:                       "version",
	kwPrintReplicationInfo:          "printReplicationInfo",
	kwPrintSecondaryReplicationInfo: "printSecondaryReplicationInfo",
}

// parseDbStatement parses db.method() and db.collection.method() statements.
// Dispatches to collection, bulk, encryption, or plan cache parsers as needed.
func (p *Parser) parseDbStatement() (ast.Node, error) {
	stmtStart := p.cur.Loc
	p.advance() // consume 'db'

	// Check for bracket access: db["coll"] or db['coll']
	if p.cur.Type == '[' {
		p.advance() // consume '['
		if p.cur.Type != tokString {
			return nil, p.errorAtCur("expected string for collection name in bracket access")
		}
		collName := p.cur.Str
		collTok := p.advance()
		collLoc := ast.Loc{Start: collTok.Loc, End: collTok.End}
		if _, err := p.expect(']'); err != nil {
			return nil, err
		}
		if _, err := p.expect('.'); err != nil {
			return nil, err
		}
		return p.parseCollectionStatement(collName, collLoc, "bracket", stmtStart)
	}

	// Must have a dot after db
	if _, err := p.expect('.'); err != nil {
		return nil, err
	}

	// After db., get the identifier/keyword
	if !p.isIdentLike(p.cur.Type) && p.cur.Type != tokIdent {
		return nil, p.syntaxErrorAtCur()
	}

	identTok := p.cur
	identName := identTok.Str

	// Check for getMongo() — could be encryption chain, connection chain, or standalone
	if identTok.Type == kwGetMongo && p.peekNext().Type == '(' {
		p.advance() // consume getMongo
		// Parse getMongo()
		args, err := p.parseArguments()
		if err != nil {
			return nil, err
		}

		// Check what follows: .getKeyVault() or .getClientEncryption() → encryption
		if p.cur.Type == '.' {
			next := p.peekNext()
			if next.Type == kwGetKeyVault || next.Type == kwGetClientEncryption {
				return p.parseEncryptionStatement(stmtStart)
			}
		}

		// If followed by .method() chain, treat as connection statement
		if p.cur.Type == '.' {
			var chain []ast.MethodCall
			for p.cur.Type == '.' {
				p.advance() // consume '.'
				mc, err := p.parseMethodCall()
				if err != nil {
					return nil, err
				}
				chain = append(chain, *mc)
			}
			return &ast.ConnectionStatement{
				Constructor:    "db.getMongo",
				Args:           args,
				ChainedMethods: chain,
				Loc:            ast.Loc{Start: stmtStart, End: p.prev.End},
			}, nil
		}

		// Otherwise it's a database method call
		return &ast.DatabaseStatement{
			Method: "getMongo",
			Args:   args,
			Loc:    ast.Loc{Start: stmtStart, End: p.prev.End},
		}, nil
	}

	// Check for getCollection("name").method() → collection access, or standalone db.getCollection(...)
	if identTok.Type == kwGetCollection && p.peekNext().Type == '(' {
		p.advance() // consume getCollection
		args, err := p.parseArguments()
		if err != nil {
			return nil, err
		}
		if p.cur.Type == '.' {
			if len(args) != 1 {
				return nil, p.errorAtCur("expected string argument for getCollection()")
			}
			strArg, ok := args[0].(*ast.StringLiteral)
			if !ok {
				return nil, p.errorAtCur("expected string argument for getCollection()")
			}
			p.advance() // consume '.'
			return p.parseCollectionStatement(strArg.Value, strArg.Loc, "getCollection", stmtStart)
		}
		return &ast.DatabaseStatement{
			Method: "getCollection",
			Args:   args,
			Loc:    ast.Loc{Start: stmtStart, End: p.prev.End},
		}, nil
	}

	// Check if this is a known database method followed by '('
	if _, isDbMethod := dbMethodNames[identTok.Type]; isDbMethod && p.peekNext().Type == '(' {
		p.advance() // consume method name
		args, err := p.parseArguments()
		if err != nil {
			return nil, err
		}
		chain, err := p.parseDatabaseMethodChain(identName)
		if err != nil {
			return nil, err
		}
		return &ast.DatabaseStatement{
			Method:         identName,
			Args:           args,
			ChainedMethods: chain,
			Loc:            ast.Loc{Start: stmtStart, End: p.prev.End},
		}, nil
	}

	// Disambiguation: after db.identifier, if next is '(' it's a db method call,
	// if it's '.' it's a collection name.
	nextTok := p.peekNext()
	if nextTok.Type == '(' {
		// Treat as database method (handles methods not in dbMethodNames too,
		// like auth, changeUserPassword, etc. that may be identifiers)
		p.advance() // consume method name
		args, err := p.parseArguments()
		if err != nil {
			return nil, err
		}
		chain, err := p.parseDatabaseMethodChain(identName)
		if err != nil {
			return nil, err
		}
		return &ast.DatabaseStatement{
			Method:         identName,
			Args:           args,
			ChainedMethods: chain,
			Loc:            ast.Loc{Start: stmtStart, End: p.prev.End},
		}, nil
	}

	if nextTok.Type == '.' {
		// It's a collection name — route to collection/bulk/plancache
		collTok := p.advance() // consume collection name
		collLoc := ast.Loc{Start: collTok.Loc, End: collTok.End}
		p.advance() // consume '.'
		return p.parseCollectionStatement(identName, collLoc, "dot", stmtStart)
	}

	return nil, p.syntaxErrorAtCur()
}

// parseDatabaseMethodChain parses optional trailing .method(args)... calls on a
// database method result (e.g. db.aggregate([...]).toArray(),
// db.watch().getResumeToken(), db.getSiblingDB("other").getMongo(), or
// db.getSiblingDB("other").users.find()).
func (p *Parser) parseDatabaseMethodChain(dbMethod string) ([]ast.MethodCall, error) {
	var chain []ast.MethodCall
	for p.cur.Type == '.' {
		p.advance() // consume '.'
		name, tok, err := p.expectIdent()
		if err != nil {
			return nil, err
		}
		// After getSiblingDB(...), allow one collection property segment before a method call:
		//   db.getSiblingDB("other").users.find()
		if dbMethod == "getSiblingDB" && len(chain) == 0 && p.cur.Type == '.' {
			chain = append(chain, ast.MethodCall{
				Name: name,
				Loc:  ast.Loc{Start: tok.Loc, End: tok.End},
			})
			continue
		}
		args, err := p.parseArguments()
		if err != nil {
			return nil, err
		}
		chain = append(chain, ast.MethodCall{
			Name: name,
			Args: args,
			Loc:  ast.Loc{Start: tok.Loc, End: p.prev.End},
		})
	}
	return chain, nil
}
