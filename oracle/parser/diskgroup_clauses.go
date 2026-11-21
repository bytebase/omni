package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// This file reads ALTER DISKGROUP by the grammar of
// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-DISKGROUP.html
// (bnf/ALTER-DISKGROUP.bnf), corrected where the pinned Oracle reads it
// otherwise: the compat rows and the probe decide. Oracle reads the clauses
// as the page did before 12c in many places: user, usergroup, attribute and
// file names are strings, add_disk_clause groups disks under [ QUORUM |
// REGULAR ] [ FAILGROUP name ] DISK, and the permission and owner clauses
// are SET PERMISSION and SET OWNERSHIP.

// dgClauseOption returns the clause the caller read from start, its text
// as its value.
func (p *Parser) dgClauseOption(key string, start int) *nodes.DDLOption {
	return &nodes.DDLOption{Key: key, Value: p.source[start:p.prev.End], Loc: nodes.Loc{Start: start, End: p.prev.End}}
}

// dgList reads item [, item ]..., each read by item.
func (p *Parser) dgList(item func() error) error {
	for {
		if err := item(); err != nil {
			return err
		}
		if p.cur.Type != ',' {
			return nil
		}
		p.advance()
	}
}

// dgString reads a string literal.
func (p *Parser) dgString() error {
	_, err := p.adbString()
	return err
}

// dgSize reads a size_clause; the clause keeps its text as its value.
func (p *Parser) dgSize() error {
	size, err := p.parseAlterDatabaseSize()
	if err == nil && size == "" {
		return p.syntaxErrorAtCur()
	}
	return err
}

// dgName reads an identifier Oracle does not reserve.
func (p *Parser) dgName() error {
	_, err := p.adbName()
	return err
}

// dgOptionalWord consumes one of words if it is the current token.
func (p *Parser) dgOptionalWord(words ...string) bool {
	for _, w := range words {
		if p.adbWord(w) {
			p.advance()
			return true
		}
	}
	return false
}

// parseAlterDiskgroupStmt parses ALTER DISKGROUP, ALTER DISKGROUP consumed.
//
//	ALTER DISKGROUP { diskgroup_name [, diskgroup_name ]... | ALL } clause
func (p *Parser) parseAlterDiskgroupStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "ALTER",
		ObjectType: nodes.OBJECT_DISKGROUP,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}
	if p.cur.Type == kwALL {
		allStart := p.pos()
		p.advance()
		opts.Items = append(opts.Items, p.adbItem("TARGET", "ALL", allStart))
	} else {
		nameStart := p.pos()
		name, err := p.adbName()
		if err != nil {
			return nil, err
		}
		stmt.Name = &nodes.ObjectName{Name: name, Loc: nodes.Loc{Start: nameStart, End: p.prev.End}}
		if p.cur.Type == ',' {
			p.advance()
			names := &nodes.List{}
			err := p.dgList(func() error {
				itemStart := p.pos()
				n, err := p.adbName()
				if err == nil {
					names.Items = append(names.Items, p.adbItem("DISKGROUP", n, itemStart))
				}
				return err
			})
			if err != nil {
				return nil, err
			}
			first := &nodes.DDLOption{Key: "DISKGROUP", Value: name, Loc: stmt.Name.Loc}
			names.Items = append([]nodes.Node{first}, names.Items...)
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "TARGET", Value: "LIST", Items: names,
				Loc: nodes.Loc{Start: nameStart, End: p.prev.End},
			})
		}
	}
	clauses, err := p.parseAlterDiskgroupClause()
	if err != nil {
		return nil, err
	}
	opts.Items = append(opts.Items, clauses...)
	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterDiskgroupClause reads the clause of ALTER DISKGROUP: one clause,
// or a sequence of add_disk_clause and drop_disk_clause, or resize_disk_clause,
// with a rebalance_diskgroup_clause after them.
func (p *Parser) parseAlterDiskgroupClause() ([]nodes.Node, error) {
	var out []nodes.Node
	if p.adbEnd() || p.cur.Type == tokQIDENT || !p.isIdentLike() {
		return nil, p.syntaxErrorAtCur()
	}
	diskClauses, resized := false, false
	for {
		start := p.pos()
		var key string
		var err error
		switch {
		case !resized && p.cur.Type == kwADD && p.dgAddsDisks():
			key, err = "ADD_DISK", p.parseDiskgroupAddDisk()
		case !resized && p.cur.Type == kwDROP && p.dgDropsDisks():
			key, err = "DROP_DISK", p.parseDiskgroupDropDisk()
		case !diskClauses && p.adbWord("RESIZE") && !p.isIdentLikeStrAt(p.peekNext(), "VOLUME"):
			key, err = "RESIZE_DISK", p.parseDiskgroupResize()
		default:
			if diskClauses {
				if p.adbWord("REBALANCE") {
					start := p.pos()
					if err := p.parseDiskgroupRebalance(); err != nil {
						return nil, err
					}
					out = append(out, p.dgClauseOption("REBALANCE", start))
				}
				if !p.adbEnd() {
					return nil, p.syntaxErrorAtCur()
				}
				return out, nil
			}
			opt, err := p.parseDiskgroupOtherClause()
			if err != nil {
				return nil, err
			}
			if !p.adbEnd() {
				return nil, p.syntaxErrorAtCur()
			}
			return []nodes.Node{opt}, nil
		}
		if err != nil {
			return nil, err
		}
		out = append(out, p.dgClauseOption(key, start))
		diskClauses = true
		resized = key == "RESIZE_DISK"
	}
}

// dgAddsDisks reports whether ADD opens add_disk_clause.
func (p *Parser) dgAddsDisks() bool {
	next := p.peekNext()
	return p.isIdentLikeStrAt(next, "DISK") || p.isIdentLikeStrAt(next, "FAILGROUP") ||
		p.isIdentLikeStrAt(next, "QUORUM") || p.isIdentLikeStrAt(next, "REGULAR")
}

// dgDropsDisks reports whether DROP opens drop_disk_clause.
func (p *Parser) dgDropsDisks() bool {
	next := p.peekNext()
	return p.isIdentLikeStrAt(next, "DISK") || p.isIdentLikeStrAt(next, "DISKS") ||
		p.isIdentLikeStrAt(next, "QUORUM") || p.isIdentLikeStrAt(next, "REGULAR")
}

// parseDiskgroupAddDisk reads
//
//	ADD { [ QUORUM | REGULAR ] [ FAILGROUP failgroup_name ] DISK qualified_disk_clause
//	      [, qualified_disk_clause ]... }...
//	qualified_disk_clause ::= 'path_name' [ NAME disk_name ] [ SIZE size_clause ] [ FORCE | NOFORCE ]
func (p *Parser) parseDiskgroupAddDisk() error {
	p.advance()
	for first := true; first || p.adbWord("QUORUM") || p.adbWord("REGULAR") || p.adbWord("FAILGROUP") || p.adbWord("DISK"); first = false {
		p.dgOptionalWord("QUORUM", "REGULAR")
		if p.adbWord("FAILGROUP") {
			p.advance()
			if err := p.dgName(); err != nil {
				return err
			}
		}
		if err := p.adbExpect("DISK"); err != nil {
			return err
		}
		err := p.dgList(func() error {
			if err := p.dgString(); err != nil {
				return err
			}
			if p.adbWord("NAME") {
				p.advance()
				if err := p.dgName(); err != nil {
					return err
				}
			}
			if p.cur.Type == kwSIZE {
				p.advance()
				if err := p.dgSize(); err != nil {
					return err
				}
			}
			p.dgOptionalWord("FORCE", "NOFORCE")
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// parseDiskgroupDropDisk reads
//
//	DROP { [ QUORUM | REGULAR ] DISK disk_name [ FORCE | NOFORCE ] [, disk_name [ FORCE | NOFORCE ] ]...
//	     | DISKS IN [ QUORUM | REGULAR ] FAILGROUP failgroup_name [ FORCE | NOFORCE ] [, ... ]... }
func (p *Parser) parseDiskgroupDropDisk() error {
	p.advance()
	p.dgOptionalWord("QUORUM", "REGULAR")
	if p.adbWord("DISKS") {
		p.advance()
		if p.cur.Type != kwIN {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		p.dgOptionalWord("QUORUM", "REGULAR")
		if err := p.adbExpect("FAILGROUP"); err != nil {
			return err
		}
	} else if err := p.adbExpect("DISK"); err != nil {
		return err
	}
	return p.dgList(func() error {
		if err := p.dgName(); err != nil {
			return err
		}
		p.dgOptionalWord("FORCE", "NOFORCE")
		return nil
	})
}

// parseDiskgroupResize reads
//
//	RESIZE { ALL [ SIZE size_clause ]
//	       | [ QUORUM | REGULAR ] DISK disk_name [ SIZE size_clause ] [, ... ]...
//	       | DISKS IN [ QUORUM | REGULAR ] FAILGROUP failgroup_name [ SIZE size_clause ] [, ... ]... }
func (p *Parser) parseDiskgroupResize() error {
	p.advance()
	optionalSize := func() error {
		if p.cur.Type == kwSIZE {
			p.advance()
			err := p.dgSize()
			return err
		}
		return nil
	}
	if p.cur.Type == kwALL {
		p.advance()
		return optionalSize()
	}
	p.dgOptionalWord("QUORUM", "REGULAR")
	if p.adbWord("DISKS") {
		p.advance()
		if p.cur.Type != kwIN {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		p.dgOptionalWord("QUORUM", "REGULAR")
		if err := p.adbExpect("FAILGROUP"); err != nil {
			return err
		}
	} else if err := p.adbExpect("DISK"); err != nil {
		return err
	}
	return p.dgList(func() error {
		if err := p.dgName(); err != nil {
			return err
		}
		return optionalSize()
	})
}

// parseDiskgroupRebalance reads
//
//	REBALANCE [ { WITH | WITHOUT } phase_list ] [ [ MODIFY ] POWER integer ] [ WAIT | NOWAIT ]
//
// Oracle reads the phase list (RESTORE, BALANCE, PREPARE, COMPACT) as the
// text up to POWER, WAIT or NOWAIT and checks it after the parse: WITH BASIC
// and WITH ( zzz pass its parser, failing only on the disk group (ORA-15260
// here), while WITH BALANCE POWER zzz fails at zzz (ORA-15102).
func (p *Parser) parseDiskgroupRebalance() error {
	p.advance()
	if p.cur.Type == kwWITH || p.adbWord("WITHOUT") {
		p.advance()
		if p.adbEnd() {
			return p.syntaxErrorAtCur()
		}
		for !p.adbEnd() && !p.isDiskgroupRebalanceOptionStart() {
			p.advance()
		}
	}
	if p.cur.Type == kwMODIFY {
		p.advance()
		if err := p.adbExpect("POWER"); err != nil {
			return err
		}
		if _, err := p.adbInteger(); err != nil {
			return err
		}
	} else if p.adbWord("POWER") {
		p.advance()
		if _, err := p.adbInteger(); err != nil {
			return err
		}
	}
	if p.cur.Type == kwWAIT || p.cur.Type == kwNOWAIT {
		p.advance()
	}
	return nil
}

// isDiskgroupRebalanceOptionStart reports whether the current token opens
// what follows a rebalance's phase list: [ MODIFY ] POWER, WAIT or NOWAIT.
func (p *Parser) isDiskgroupRebalanceOptionStart() bool {
	return p.cur.Type == kwMODIFY || p.adbWord("POWER") || p.cur.Type == kwWAIT || p.cur.Type == kwNOWAIT
}

// parseDiskgroupOtherClause reads a clause of ALTER DISKGROUP other than
// the disk clauses.
func (p *Parser) parseDiskgroupOtherClause() (*nodes.DDLOption, error) {
	start := p.pos()
	word := p.cur.Str
	var key string
	var err error
	switch word {
	case "ADD":
		key, err = p.parseDiskgroupAdd()
	case "DROP":
		key, err = p.parseDiskgroupDrop()
	case "RESIZE":
		p.advance()
		p.advance()
		key = "RESIZE_VOLUME"
		if err = p.dgName(); err == nil {
			if p.cur.Type != kwSIZE {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			err = p.dgSize()
		}
	case "REPLACE":
		key, err = p.parseDiskgroupReplace()
	case "REPAIR":
		// Oracle consumes REPAIR before refusing the unsupported top-level
		// REPAIR clause at the next token.
		p.advance()
		return nil, p.syntaxErrorAtCur()
	case "RENAME":
		key, err = p.parseDiskgroupRename()
	case "ONLINE", "OFFLINE":
		key, err = word, p.parseDiskgroupOnlineOffline(word)
	case "REBALANCE":
		key, err = "REBALANCE", p.parseDiskgroupRebalance()
	case "CHECK":
		key, err = "CHECK", p.parseDiskgroupCheck()
	case "MODIFY", "ALTER":
		key, err = p.parseDiskgroupModify(word)
	case "SET":
		key, err = p.parseDiskgroupSet()
	case "CONVERT":
		p.advance()
		key = "CONVERT_REDUNDANCY"
		if err = p.adbExpect("REDUNDANCY"); err == nil {
			if p.cur.Type != kwTO {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			err = p.adbExpect("FLEX")
		}
	case "SCRUB":
		key, err = "SCRUB", p.parseDiskgroupScrub()
	case "MOVE":
		key, err = p.parseDiskgroupMove()
	case "UNDROP":
		p.advance()
		key, err = "UNDROP_DISKS", p.adbExpect("DISKS")
	case "MOUNT":
		p.advance()
		key = "MOUNT"
		p.dgOptionalWord("RESTRICTED", "NORMAL")
		p.dgOptionalWord("FORCE", "NOFORCE")
	case "DISMOUNT":
		p.advance()
		key = "DISMOUNT"
		p.dgOptionalWord("FORCE", "NOFORCE")
	case "ENABLE", "DISABLE":
		p.advance()
		key = word + "_VOLUME"
		if err = p.adbExpect("VOLUME"); err == nil {
			if p.cur.Type == kwALL {
				p.advance()
			} else {
				err = p.dgList(p.dgName)
			}
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if err != nil {
		return nil, err
	}
	return p.dgClauseOption(key, start), nil
}

// parseDiskgroupAdd reads the ADD clauses other than add_disk_clause:
//
//	ADD TEMPLATE template_name ATTRIBUTES ( template_attributes )
//	ADD DIRECTORY 'path' [, 'path' ]...
//	ADD ALIAS 'alias' FOR 'file' [, 'alias' FOR 'file' ]...
//	ADD VOLUME asm_volume SIZE size_clause [ redundancy ] [ STRIPE_WIDTH size_clause ] [ STRIPE_COLUMNS integer ]
//	ADD USERGROUP 'usergroup' [ WITH MEMBER 'user' [, 'user' ]... ]
//	ADD USER 'user' [, 'user' ]...
//	ADD QUOTAGROUP quotagroup_name [ SET QUOTA = { size_clause | UNLIMITED } ]
//	ADD FILEGROUP filegroup_name { DATABASE db_name | CLUSTER cluster_name | VOLUME volume_name | TEMPLATE }
//	  [ FROM TEMPLATE template_name ] [ SET 'property' = 'value' ]
func (p *Parser) parseDiskgroupAdd() (string, error) {
	p.advance()
	word := p.cur.Str
	if p.cur.Type == tokQIDENT {
		word = ""
	}
	switch word {
	case "TEMPLATE":
		p.advance()
		if err := p.dgName(); err != nil {
			return "", err
		}
		return "ADD_TEMPLATE", p.parseDiskgroupTemplateAttributes("ATTRIBUTES")
	case "DIRECTORY":
		p.advance()
		return "ADD_DIRECTORY", p.dgList(p.dgString)
	case "ALIAS":
		p.advance()
		return "ADD_ALIAS", p.dgList(func() error {
			if err := p.dgString(); err != nil {
				return err
			}
			if p.cur.Type != kwFOR {
				return p.syntaxErrorAtCur()
			}
			p.advance()
			return p.dgString()
		})
	case "VOLUME":
		p.advance()
		if err := p.dgName(); err != nil {
			return "", err
		}
		if p.cur.Type != kwSIZE {
			return "", p.syntaxErrorAtCur()
		}
		p.advance()
		if err := p.dgSize(); err != nil {
			return "", err
		}
		p.dgOptionalWord("MIRROR", "HIGH", "UNPROTECTED", "PARITY", "DOUBLE")
		if p.adbWord("STRIPE_WIDTH") {
			p.advance()
			if err := p.dgSize(); err != nil {
				return "", err
			}
		}
		if p.adbWord("STRIPE_COLUMNS") {
			p.advance()
			if _, err := p.adbInteger(); err != nil {
				return "", err
			}
		}
		return "ADD_VOLUME", nil
	case "USERGROUP":
		p.advance()
		if err := p.dgString(); err != nil {
			return "", err
		}
		if p.cur.Type == kwWITH {
			p.advance()
			if err := p.adbExpect("MEMBER"); err != nil {
				return "", err
			}
			return "ADD_USERGROUP", p.dgList(p.dgString)
		}
		return "ADD_USERGROUP", nil
	case "USER":
		p.advance()
		return "ADD_USER", p.dgList(p.dgString)
	case "QUOTAGROUP":
		p.advance()
		if err := p.dgName(); err != nil {
			return "", err
		}
		if p.cur.Type == kwSET {
			return "ADD_QUOTAGROUP", p.parseDiskgroupSetQuota()
		}
		return "ADD_QUOTAGROUP", nil
	case "FILEGROUP":
		p.advance()
		if err := p.dgName(); err != nil {
			return "", err
		}
		switch {
		case p.cur.Type == kwDATABASE, p.adbWord("CLUSTER"), p.adbWord("VOLUME"):
			p.advance()
			if err := p.dgName(); err != nil {
				return "", err
			}
		case p.adbWord("TEMPLATE"):
			p.advance()
		default:
			return "", p.syntaxErrorAtCur()
		}
		if p.cur.Type == kwFROM {
			p.advance()
			if err := p.adbExpect("TEMPLATE"); err != nil {
				return "", err
			}
			if err := p.dgName(); err != nil {
				return "", err
			}
		}
		if p.cur.Type == kwSET {
			return "ADD_FILEGROUP", p.parseDiskgroupSetProperty()
		}
		return "ADD_FILEGROUP", nil
	}
	return "", p.syntaxErrorAtCur()
}

// parseDiskgroupSetQuota reads SET QUOTA = { size_clause | UNLIMITED }.
func (p *Parser) parseDiskgroupSetQuota() error {
	p.advance()
	if err := p.adbExpect("QUOTA"); err != nil {
		return err
	}
	if err := p.pdbExpectEquals(); err != nil {
		return err
	}
	if p.adbWord("UNLIMITED") {
		p.advance()
		return nil
	}
	err := p.dgSize()
	return err
}

// parseDiskgroupSetProperty reads SET 'property' = 'value', one property.
func (p *Parser) parseDiskgroupSetProperty() error {
	p.advance()
	if err := p.dgString(); err != nil {
		return err
	}
	if err := p.pdbExpectEquals(); err != nil {
		return err
	}
	return p.dgString()
}

// parseDiskgroupTemplateAttributes reads, at keyword (ATTRIBUTES or
// ATTRIBUTE),
//
//	keyword ( [ MIRROR | HIGH | UNPROTECTED | PARITY | DOUBLE ] [ FINE | COARSE ]
//	          [ HOT | COLD ] [ MIRRORHOT | MIRRORCOLD ] )
//
// Oracle takes the attributes in any order, each kind at most once.
func (p *Parser) parseDiskgroupTemplateAttributes(keyword string) error {
	if err := p.adbExpect(keyword); err != nil {
		return err
	}
	if p.cur.Type != '(' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	kinds := [][]string{
		{"MIRROR", "HIGH", "UNPROTECTED", "PARITY", "DOUBLE"},
		{"FINE", "COARSE"},
		{"HOT", "COLD"},
		{"MIRRORHOT", "MIRRORCOLD"},
	}
	seen := make([]bool, len(kinds))
	for p.cur.Type != ')' {
		matched := false
		for i, words := range kinds {
			if !seen[i] && p.dgOptionalWord(words...) {
				seen[i], matched = true, true
				break
			}
		}
		if !matched {
			return p.syntaxErrorAtCur()
		}
	}
	p.advance()
	return nil
}

// parseDiskgroupDrop reads the DROP clauses other than drop_disk_clause:
//
//	DROP TEMPLATE template_name [, template_name ]...
//	DROP DIRECTORY 'path' [ FORCE | NOFORCE ] [, ... ]...
//	DROP ALIAS 'alias' [, 'alias' ]...
//	DROP VOLUME asm_volume
//	DROP FILE 'file' [, 'file' ]...
//	DROP USERGROUP 'usergroup'
//	DROP USER 'user' [ CASCADE ] [, ... ]...
//	DROP QUOTAGROUP quotagroup_name
//	DROP FILEGROUP filegroup_name [ CASCADE [ FOR DATABASE db_unique_name ] ]
//
// Oracle does not take the page's FOR PLUGGABLE DATABASE (ORA-00905).
func (p *Parser) parseDiskgroupDrop() (string, error) {
	p.advance()
	word := p.cur.Str
	if p.cur.Type == tokQIDENT {
		word = ""
	}
	key := "DROP_" + word
	switch word {
	case "TEMPLATE":
		p.advance()
		return key, p.dgList(p.dgName)
	case "DIRECTORY":
		p.advance()
		return key, p.dgList(func() error {
			if err := p.dgString(); err != nil {
				return err
			}
			p.dgOptionalWord("FORCE", "NOFORCE")
			return nil
		})
	case "ALIAS", "FILE":
		p.advance()
		return key, p.dgList(p.dgString)
	case "VOLUME", "QUOTAGROUP":
		p.advance()
		return key, p.dgName()
	case "USERGROUP":
		p.advance()
		return key, p.dgString()
	case "USER":
		p.advance()
		return key, p.dgList(func() error {
			if err := p.dgString(); err != nil {
				return err
			}
			p.dgOptionalWord("CASCADE")
			return nil
		})
	case "FILEGROUP":
		p.advance()
		if err := p.dgName(); err != nil {
			return "", err
		}
		if p.dgOptionalWord("CASCADE") && p.cur.Type == kwFOR {
			p.advance()
			if p.cur.Type != kwDATABASE {
				return "", p.syntaxErrorAtCur()
			}
			p.advance()
			return key, p.dgName()
		}
		return key, nil
	}
	return "", p.syntaxErrorAtCur()
}

// parseDiskgroupReplace reads
//
//	REPLACE DISK disk_name WITH 'path' [ FORCE | NOFORCE ] [ POWER integer ] [ WAIT | NOWAIT ]
//	REPLACE USER 'old_user' WITH 'new_user'
func (p *Parser) parseDiskgroupReplace() (string, error) {
	p.advance()
	switch {
	case p.adbWord("DISK"):
		p.advance()
		if err := p.dgName(); err != nil {
			return "", err
		}
		if p.cur.Type != kwWITH {
			return "", p.syntaxErrorAtCur()
		}
		p.advance()
		if err := p.dgString(); err != nil {
			return "", err
		}
		p.dgOptionalWord("FORCE", "NOFORCE")
		if p.adbWord("POWER") {
			p.advance()
			if _, err := p.adbInteger(); err != nil {
				return "", err
			}
		}
		if p.cur.Type == kwWAIT || p.cur.Type == kwNOWAIT {
			p.advance()
		}
		return "REPLACE_DISK", nil
	case p.cur.Type == kwUSER || p.adbWord("USER"):
		p.advance()
		if err := p.dgString(); err != nil {
			return "", err
		}
		if p.cur.Type != kwWITH {
			return "", p.syntaxErrorAtCur()
		}
		p.advance()
		return "REPLACE_USER", p.dgString()
	}
	return "", p.syntaxErrorAtCur()
}

// parseDiskgroupRename reads
//
//	RENAME DISK 'old_disk_name' TO 'new_disk_name' [, ... ]...
//	RENAME DISKS ALL
//	RENAME DIRECTORY 'old_path' TO 'new_path' [, ... ]...
//	RENAME ALIAS 'old_alias' TO 'new_alias' [, ... ]...
//
// Oracle names the disks with strings here (ORA-15107).
func (p *Parser) parseDiskgroupRename() (string, error) {
	p.advance()
	if p.adbWord("DISKS") {
		p.advance()
		if p.cur.Type != kwALL {
			return "", p.syntaxErrorAtCur()
		}
		p.advance()
		return "RENAME_DISKS_ALL", nil
	}
	if !p.adbWord("DISK") && !p.adbWord("DIRECTORY") && !p.adbWord("ALIAS") {
		return "", p.syntaxErrorAtCur()
	}
	key := "RENAME_" + p.advance().Str
	return key, p.dgList(func() error {
		if err := p.dgString(); err != nil {
			return err
		}
		if p.cur.Type != kwTO {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		return p.dgString()
	})
}

// parseDiskgroupOnlineOffline reads
//
//	ONLINE { { [ QUORUM | REGULAR ] DISK disk_name [, disk_name ]...
//	         | DISKS IN [ QUORUM | REGULAR ] FAILGROUP failgroup_name [, ... ]... }...
//	       | ALL } [ POWER integer ] [ WAIT | NOWAIT ]
//	OFFLINE { [ QUORUM | REGULAR ] DISK disk_name [, disk_name ]...
//	        | DISKS IN [ QUORUM | REGULAR ] FAILGROUP failgroup_name [, ... ]... }...
//	        [ DROP AFTER integer { M | H } ] [ WAIT | NOWAIT ]
func (p *Parser) parseDiskgroupOnlineOffline(verb string) error {
	p.advance()
	if verb == "ONLINE" && p.cur.Type == kwALL {
		p.advance()
	} else {
		for first := true; first || p.adbWord("QUORUM") || p.adbWord("REGULAR") || p.adbWord("DISK") || p.adbWord("DISKS"); first = false {
			p.dgOptionalWord("QUORUM", "REGULAR")
			if p.adbWord("DISKS") {
				p.advance()
				if p.cur.Type != kwIN {
					return p.syntaxErrorAtCur()
				}
				p.advance()
				p.dgOptionalWord("QUORUM", "REGULAR")
				if err := p.adbExpect("FAILGROUP"); err != nil {
					return err
				}
			} else if err := p.adbExpect("DISK"); err != nil {
				return err
			}
			if err := p.dgList(p.dgName); err != nil {
				return err
			}
		}
	}
	if verb == "ONLINE" && p.adbWord("POWER") {
		p.advance()
		if _, err := p.adbInteger(); err != nil {
			return err
		}
	}
	if verb == "OFFLINE" && p.cur.Type == kwDROP {
		p.advance()
		if err := p.adbExpect("AFTER"); err != nil {
			return err
		}
		if _, err := p.adbInteger(); err != nil {
			return err
		}
		if !p.dgOptionalWord("M", "H") {
			return p.syntaxErrorAtCur()
		}
	}
	if p.cur.Type == kwWAIT || p.cur.Type == kwNOWAIT {
		p.advance()
	}
	return nil
}

// parseDiskgroupCheck reads
//
//	CHECK [ ALL | DISK | FILE [ 'file' [, 'file' ]... ] ] [ REPAIR | NOREPAIR ]
//
// Oracle takes DISK and FILE without the names the page gives them; it does
// not take disk names or CHECK DISKS IN FAILGROUP (ORA-00905).
func (p *Parser) parseDiskgroupCheck() error {
	p.advance()
	switch {
	case p.cur.Type == kwALL:
		p.advance()
	case p.cur.Type == kwFILE:
		p.advance()
		if p.cur.Type == tokSCONST {
			if err := p.dgList(p.dgString); err != nil {
				return err
			}
		}
	case p.adbWord("DISK"):
		p.advance()
	}
	p.dgOptionalWord("REPAIR", "NOREPAIR")
	return nil
}

// parseDiskgroupModify reads the clauses that begin with MODIFY or ALTER:
//
//	{ MODIFY | ALTER } TEMPLATE template_name ATTRIBUTES ( template_attributes )
//	MODIFY VOLUME asm_volume [ MOUNTPATH 'path' ] [ USAGE 'usage' ]
//	MODIFY FILE 'file' ATTRIBUTE ( template_attributes ) [, ... ]...
//	MODIFY USERGROUP 'usergroup' { ADD | DROP } MEMBER 'user' [, 'user' ]...
//	MODIFY QUOTAGROUP quotagroup_name SET QUOTA = { size_clause | UNLIMITED }
//	MODIFY FILEGROUP filegroup_name SET 'property' = 'value'
func (p *Parser) parseDiskgroupModify(verb string) (string, error) {
	p.advance()
	if p.adbWord("TEMPLATE") {
		p.advance()
		if err := p.dgName(); err != nil {
			return "", err
		}
		return verb + "_TEMPLATE", p.parseDiskgroupTemplateAttributes("ATTRIBUTES")
	}
	if verb != "MODIFY" {
		return "", p.syntaxErrorAtCur()
	}
	switch {
	case p.adbWord("VOLUME"):
		p.advance()
		if err := p.dgName(); err != nil {
			return "", err
		}
		if !p.adbWord("MOUNTPATH") && !p.adbWord("USAGE") {
			return "", p.syntaxErrorAtCur()
		}
		for _, w := range []string{"MOUNTPATH", "USAGE"} {
			if p.adbWord(w) {
				p.advance()
				if err := p.dgString(); err != nil {
					return "", err
				}
			}
		}
		return "MODIFY_VOLUME", nil
	case p.cur.Type == kwFILE:
		p.advance()
		return "MODIFY_FILE", p.dgList(func() error {
			if err := p.dgString(); err != nil {
				return err
			}
			return p.parseDiskgroupTemplateAttributes("ATTRIBUTE")
		})
	case p.adbWord("USERGROUP"):
		p.advance()
		if err := p.dgString(); err != nil {
			return "", err
		}
		if p.cur.Type != kwADD && p.cur.Type != kwDROP {
			return "", p.syntaxErrorAtCur()
		}
		p.advance()
		if err := p.adbExpect("MEMBER"); err != nil {
			return "", err
		}
		return "MODIFY_USERGROUP", p.dgList(p.dgString)
	case p.adbWord("QUOTAGROUP"):
		p.advance()
		if err := p.dgName(); err != nil {
			return "", err
		}
		if p.cur.Type != kwSET {
			return "", p.syntaxErrorAtCur()
		}
		return "MODIFY_QUOTAGROUP", p.parseDiskgroupSetQuota()
	case p.adbWord("FILEGROUP"):
		p.advance()
		if err := p.dgName(); err != nil {
			return "", err
		}
		if p.cur.Type != kwSET {
			return "", p.syntaxErrorAtCur()
		}
		return "MODIFY_FILEGROUP", p.parseDiskgroupSetProperty()
	}
	return "", p.syntaxErrorAtCur()
}

// parseDiskgroupSet reads
//
//	SET ATTRIBUTE 'attribute_name' = 'attribute_value'
//	SET PERMISSION { OWNER | GROUP | OTHER } = { NONE | READ ONLY | READ WRITE } [, ... ]...
//	  FOR FILE 'file' [, 'file' ]...
//	SET OWNERSHIP { OWNER | GROUP } = 'name' [, ... ]... FOR FILE 'file' [, 'file' ]...
//
// Oracle takes one attribute (ORA-15213), and the permission and owner
// clauses in these forms, not the page's (ORA-00905).
func (p *Parser) parseDiskgroupSet() (string, error) {
	p.advance()
	forFiles := func() error {
		if p.cur.Type != kwFOR {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		if p.cur.Type != kwFILE {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		return p.dgList(p.dgString)
	}
	switch {
	case p.adbWord("ATTRIBUTE"):
		p.advance()
		if err := p.dgString(); err != nil {
			return "", err
		}
		if err := p.pdbExpectEquals(); err != nil {
			return "", err
		}
		return "SET_ATTRIBUTE", p.dgString()
	case p.adbWord("PERMISSION"):
		p.advance()
		err := p.dgList(func() error {
			if !p.dgOptionalWord("OWNER", "GROUP", "OTHER") {
				return p.syntaxErrorAtCur()
			}
			if err := p.pdbExpectEquals(); err != nil {
				return err
			}
			if p.dgOptionalWord("NONE") {
				return nil
			}
			if p.cur.Type != kwREAD {
				return p.syntaxErrorAtCur()
			}
			p.advance()
			if p.cur.Type != kwONLY && p.cur.Type != kwWRITE {
				return p.syntaxErrorAtCur()
			}
			p.advance()
			return nil
		})
		if err != nil {
			return "", err
		}
		return "SET_PERMISSION", forFiles()
	case p.adbWord("OWNERSHIP"):
		p.advance()
		err := p.dgList(func() error {
			if !p.dgOptionalWord("OWNER", "GROUP") {
				return p.syntaxErrorAtCur()
			}
			if err := p.pdbExpectEquals(); err != nil {
				return err
			}
			return p.dgString()
		})
		if err != nil {
			return "", err
		}
		return "SET_OWNERSHIP", forFiles()
	}
	return "", p.syntaxErrorAtCur()
}

// parseDiskgroupScrub reads
//
//	SCRUB [ FILE 'file' | DISK disk_name ] [ REPAIR | NOREPAIR ]
//	  [ POWER { AUTO | LOW | HIGH | MAX } ] [ WAIT | NOWAIT ] [ FORCE | NOFORCE ]
//	SCRUB STOP
func (p *Parser) parseDiskgroupScrub() error {
	p.advance()
	if p.adbWord("STOP") {
		p.advance()
		return nil
	}
	switch {
	case p.cur.Type == kwFILE:
		p.advance()
		if err := p.dgString(); err != nil {
			return err
		}
	case p.adbWord("DISK"):
		p.advance()
		if err := p.dgName(); err != nil {
			return err
		}
	}
	p.dgOptionalWord("REPAIR", "NOREPAIR")
	if p.adbWord("POWER") {
		p.advance()
		if !p.dgOptionalWord("AUTO", "LOW", "HIGH", "MAX") {
			return p.syntaxErrorAtCur()
		}
	}
	if p.cur.Type == kwWAIT || p.cur.Type == kwNOWAIT {
		p.advance()
	}
	p.dgOptionalWord("FORCE", "NOFORCE")
	return nil
}

// parseDiskgroupMove reads
//
//	MOVE FILEGROUP filegroup_name TO quotagroup_name
//	MOVE FILE 'file' TO FILEGROUP filegroup_name
//
// Oracle names the quota group without the page's QUOTAGROUP.
func (p *Parser) parseDiskgroupMove() (string, error) {
	p.advance()
	switch {
	case p.adbWord("FILEGROUP"):
		p.advance()
		if err := p.dgName(); err != nil {
			return "", err
		}
		if p.cur.Type != kwTO {
			return "", p.syntaxErrorAtCur()
		}
		p.advance()
		return "MOVE_FILEGROUP", p.dgName()
	case p.cur.Type == kwFILE:
		p.advance()
		if err := p.dgString(); err != nil {
			return "", err
		}
		if p.cur.Type != kwTO {
			return "", p.syntaxErrorAtCur()
		}
		p.advance()
		if err := p.adbExpect("FILEGROUP"); err != nil {
			return "", err
		}
		return "MOVE_FILE", p.dgName()
	}
	return "", p.syntaxErrorAtCur()
}
