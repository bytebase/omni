package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Data protection objects, read as ast.ObjectStmt: BACKUP SET, BACKUP POLICY,
// SNAPSHOT SET, SNAPSHOT POLICY, SNAPSHOT (block storage), STORAGE LIFECYCLE
// POLICY, DATA MOVEMENT POLICY and DATA MOVEMENT RULE
// ---------------------------------------------------------------------------
//
// These objects are registered before the policy dispatch of create_table.go
// and database_schema.go (startsPolicyKeyword / consumePolicyKeywords), so
// BACKUP POLICY, SNAPSHOT POLICY and the rest never reach it.
//
// The pages of BACKUP SET, BACKUP POLICY, SNAPSHOT SET and SNAPSHOT POLICY
// write DROP <kind> <name> without IF EXISTS (dropNoIfExists).

func init() {
	registerObject(&objectForm{
		words:  []string{"BACKUP", "SET"},
		create: parseCreateBackupSet,
		alter:  parseAlterBackupSet,
		// https://docs.snowflake.com/en/sql-reference/sql/create-backup-set
		orReplace: true, orAlter: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-backup-set
		drop: true, dropNoIfExists: true,
		// The ALTER page writes no IF EXISTS.
		noIfExists: true,
		props:      dataProtectionCommentProps,
	})
	registerObject(&objectForm{
		words:  []string{"BACKUP", "POLICY"},
		create: parseCreateBackupPolicy,
		alter:  parseAlterBackupPolicy,
		// https://docs.snowflake.com/en/sql-reference/sql/create-backup-policy
		orReplace: true, orAlter: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-backup-policy
		drop: true, dropNoIfExists: true,
		// The ALTER page writes no IF EXISTS.
		noIfExists: true,
		props:      dataProtectionPolicyProps,
	})
	registerObject(&objectForm{
		words:  []string{"SNAPSHOT", "SET"},
		create: parseCreateSnapshotSet,
		alter:  parseAlterSnapshotSet,
		// https://docs.snowflake.com/en/sql-reference/sql/create-snapshot-set
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-snapshot-set
		drop: true, dropNoIfExists: true,
		// The ALTER page writes no IF EXISTS.
		noIfExists: true,
		props:      dataProtectionCommentProps,
	})
	registerObject(&objectForm{
		words:  []string{"SNAPSHOT", "POLICY"},
		create: parseCreateSnapshotPolicy,
		alter:  parseAlterSnapshotPolicy,
		// https://docs.snowflake.com/en/sql-reference/sql/create-snapshot-policy
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-snapshot-policy
		drop: true, dropNoIfExists: true,
		// The ALTER page writes no IF EXISTS.
		noIfExists: true,
		props:      dataProtectionPolicyProps,
	})
	registerObject(&objectForm{
		words:  []string{"SNAPSHOT"},
		create: parseCreateSnapshot,
		alter:  parseAlterSnapshot,
		// https://docs.snowflake.com/en/sql-reference/sql/create-snapshot
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-snapshot
		drop: true,
		// https://docs.snowflake.com/en/sql-reference/sql/undrop-snapshot
		// also writes RENAME TO <new_snapshot_name>.
		undrop:       true,
		undropRename: true,
		props:        dataProtectionCommentProps,
		unset:        []string{}, // ALTER SNAPSHOT documents no UNSET
	})
	// A STORAGE LIFECYCLE POLICY's and a DATA MOVEMENT RULE's body is read
	// by data_protection_body.go.
	registerObject(&objectForm{
		words:  []string{"STORAGE", "LIFECYCLE", "POLICY"},
		create: parseCreateStorageLifecyclePolicy,
		alter:  parseAlterStorageLifecyclePolicy,
		// https://docs.snowflake.com/en/sql-reference/sql/create-storage-lifecycle-policy
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-storage-lifecycle-policy
		drop: true,
		// https://docs.snowflake.com/en/sql-reference/sql/alter-storage-lifecycle-policy
		props: objectProps{
			"ARCHIVE_TIER":     propOneOf("COOL", "COLD"),
			"ARCHIVE_FOR_DAYS": propNumber,
			"COMMENT":          propString,
		},
		unset: []string{"ARCHIVE_FOR_DAYS", "COMMENT"},
	})
	registerObject(&objectForm{
		words:  []string{"DATA", "MOVEMENT", "POLICY"},
		create: parseCreateDataMovementPolicy,
		alter:  parseAlterDataMovementPolicy,
		// https://docs.snowflake.com/en/sql-reference/sql/create-data-movement-policy
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-data-movement-policy
		drop:  true,
		props: dataMovementPolicyProps,
		// https://docs.snowflake.com/en/sql-reference/sql/alter-data-movement-policy
		unset: []string{"COMMENT"},
	})
	registerObject(&objectForm{
		words:  []string{"DATA", "MOVEMENT", "RULE"},
		create: parseCreateDataMovementRule,
		alter:  parseAlterDataMovementRule,
		// https://docs.snowflake.com/en/sql-reference/sql/create-data-movement-rule
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-data-movement-rule
		drop: true,
		// CREATE takes TYPE = '<movement_type>' and COMMENT; ALTER SETs and
		// UNSETs COMMENT alone.
		props:      objectProps{"TYPE": propString, "COMMENT": propString},
		alterProps: dataProtectionCommentProps,
		unset:      []string{"COMMENT"},
	})
}

// dataProtectionCommentProps is the one property of BACKUP SET, SNAPSHOT SET
// and SNAPSHOT: CREATE takes COMMENT, and ALTER SETs and UNSETs it (on MODIFY
// <word> IDENTIFIER too). It is also the one property ALTER DATA MOVEMENT RULE
// SETs.
var dataProtectionCommentProps = objectProps{"COMMENT": propString}

// dataProtectionPolicyProps are BACKUP POLICY's and SNAPSHOT POLICY's
// properties, which CREATE takes and ALTER SETs and UNSETs.
var dataProtectionPolicyProps = objectProps{
	"SCHEDULE":          propString,
	"EXPIRE_AFTER_DAYS": propNumber,
	"COMMENT":           propString,
}

// dataMovementPolicyProps are DATA MOVEMENT POLICY's properties, which CREATE
// takes and ALTER SETs; ALTER also ADDs to and REMOVEs from the two rule lists.
var dataMovementPolicyProps = objectProps{
	"ENFORCE_RULES": propNames,
	"ALERT_RULES":   propNames,
	"COMMENT":       propString,
}

// parseCreateBackupSet reads what follows the name of
//
//	CREATE [ OR REPLACE ] BACKUP SET [ IF NOT EXISTS ] <name>
//	  FOR [ DYNAMIC ] [ ICEBERG ] TABLE <table_name>
//	  [ [ WITH ] TAG ( <tag_name> = '<tag_value>' [ , ... ] ) ]
//	  [ WITH BACKUP POLICY <policy_name> ]
//	  [ COMMENT = 'string_literal' ]
//
//	CREATE [ OR REPLACE ] BACKUP SET [ IF NOT EXISTS ] <name>
//	  FOR SCHEMA <schema_name> ...
//	CREATE [ OR REPLACE ] BACKUP SET [ IF NOT EXISTS ] <name>
//	  FOR DATABASE <database_name> ...
//
//	CREATE OR ALTER BACKUP SET <name>
//	  FOR { [ DYNAMIC ] [ ICEBERG ] TABLE <table_name> | SCHEMA <schema_name> | DATABASE <database_name> }
//	  [ WITH BACKUP POLICY <policy_name> ]
//	  [ COMMENT = 'string_literal' ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-backup-set
func parseCreateBackupSet(p *Parser, s *ast.ObjectStmt) error {
	return p.parseDataProtectionSet(s, "BACKUP", true)
}

// parseCreateSnapshotSet reads what follows the name of
//
//	CREATE [ OR REPLACE ] SNAPSHOT SET [ IF NOT EXISTS ] <name>
//	  FOR { [ DYNAMIC ] TABLE <table_name> | SCHEMA <schema_name> | DATABASE <database_name> }
//	  [ [ WITH ] TAG ( <tag_name> = '<tag_value>' [ , ... ] ) ]
//	  [ WITH SNAPSHOT POLICY <policy_name> ]
//	  [ COMMENT = 'string_literal' ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-snapshot-set
func parseCreateSnapshotSet(p *Parser, s *ast.ObjectStmt) error {
	return p.parseDataProtectionSet(s, "SNAPSHOT", false)
}

// parseDataProtectionSet reads a BACKUP SET's or SNAPSHOT SET's FOR clause,
// then its tags, WITH <word> POLICY <name> and COMMENT. word is BACKUP or
// SNAPSHOT; iceberg says whether FOR [ DYNAMIC ] ICEBERG TABLE is documented.
func (p *Parser) parseDataProtectionSet(s *ast.ObjectStmt, word string, iceberg bool) error {
	if s.OrAlter && s.IfNotExists {
		// CREATE OR ALTER <word> SET is written without IF NOT EXISTS.
		return p.syntaxErrorAtCur()
	}
	start := p.cur.Loc.Start
	if err := p.expectWords("FOR"); err != nil {
		return err
	}
	keywords := "FOR"
	var name *ast.ObjectName
	var err error
	switch {
	case p.acceptWords("SCHEMA"):
		// FOR SCHEMA [ <db>. ]<schema>
		keywords += " SCHEMA"
		if name, err = p.parseObjectName(); err == nil && !name.Database.IsEmpty() {
			err = &ParseError{
				Position: name.Loc.Start, End: name.Loc.End,
				Message: "syntax error: a schema's name has at most 2 part(s)",
			}
		}
	case p.acceptWords("DATABASE"):
		// FOR DATABASE <database_name>
		keywords += " DATABASE"
		name, err = p.parseObjectSimpleName()
	default:
		if p.acceptWords("DYNAMIC") {
			keywords += " DYNAMIC"
		}
		if iceberg && p.acceptWords("ICEBERG") {
			keywords += " ICEBERG"
		}
		if err := p.expectWords("TABLE"); err != nil {
			return err
		}
		keywords += " TABLE"
		name, err = p.parseObjectName()
	}
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, keywords, start).Names = []*ast.ObjectName{name}

	withPolicy := false
	for {
		c, err := p.parseObjectProperties(s, "", false)
		if err != nil {
			return err
		}
		if c != nil && len(c.Tags) > 0 && s.OrAlter {
			return p.dataProtectionTagError(c)
		}
		start := p.cur.Loc.Start
		if withPolicy || !p.acceptWords("WITH", word, "POLICY") {
			return nil
		}
		withPolicy = true
		policy, err := p.parseObjectName()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "WITH "+word+" POLICY", start).Names = []*ast.ObjectName{policy}
	}
}

// parseCreateBackupPolicy reads what follows the name of
//
//	CREATE [ OR REPLACE ] BACKUP POLICY [ IF NOT EXISTS ] <name>
//	  [ [ WITH ] TAG ( <tag_name> = '<tag_value>' [ , ... ] ) ]
//	  [ WITH RETENTION LOCK ]
//	  [ SCHEDULE = '{ <num> MINUTE | USING CRON <expr> <time_zone> }' ]
//	  [ EXPIRE_AFTER_DAYS = <days_integer> ]
//	  [ COMMENT = '<string_literal>' ]
//
//	CREATE OR ALTER BACKUP POLICY <name>
//	  [ WITH RETENTION LOCK ] [ SCHEDULE = ... ] [ EXPIRE_AFTER_DAYS = ... ] [ COMMENT = ... ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-backup-policy
func parseCreateBackupPolicy(p *Parser, s *ast.ObjectStmt) error {
	return p.parseDataProtectionPolicy(s)
}

// parseCreateSnapshotPolicy reads what follows the name of
//
//	CREATE [ OR REPLACE ] SNAPSHOT POLICY [ IF NOT EXISTS ] <name>
//	  [ [ WITH ] TAG ( <tag_name> = '<tag_value>' [ , ... ] ) ]
//	  [ WITH RETENTION LOCK ]
//	  [ SCHEDULE = '{ <num> MINUTE | USING CRON <expr> <time_zone> }' ]
//	  [ EXPIRE_AFTER_DAYS = <days_integer> ]
//	  [ COMMENT = '<string_literal>' ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-snapshot-policy
func parseCreateSnapshotPolicy(p *Parser, s *ast.ObjectStmt) error {
	return p.parseDataProtectionPolicy(s)
}

// parseDataProtectionPolicy reads a BACKUP POLICY's or SNAPSHOT POLICY's
// tags, WITH RETENTION LOCK and properties.
func (p *Parser) parseDataProtectionPolicy(s *ast.ObjectStmt) error {
	if s.OrAlter && s.IfNotExists {
		// CREATE OR ALTER BACKUP POLICY is written without IF NOT EXISTS.
		return p.syntaxErrorAtCur()
	}
	lock := false
	for {
		c, err := p.parseObjectProperties(s, "", false)
		if err != nil {
			return err
		}
		if c != nil && len(c.Tags) > 0 && s.OrAlter {
			return p.dataProtectionTagError(c)
		}
		start := p.cur.Loc.Start
		if lock || !p.acceptWords("WITH", "RETENTION", "LOCK") {
			return nil
		}
		lock = true
		p.addObjectKeywordsClause(s, "WITH RETENTION LOCK", start)
	}
}

// dataProtectionTagError reports the tags of clause c, which CREATE OR ALTER
// does not document.
func (p *Parser) dataProtectionTagError(c *ast.ObjectClause) error {
	return &ParseError{
		Position: c.Tags[0].Name.Loc.Start, End: c.Tags[0].Name.Loc.End,
		Message: "syntax error: CREATE OR ALTER does not take tags",
	}
}

// parseAlterBackupSet reads what follows the name of
//
//	ALTER BACKUP SET <name> ADD BACKUP
//	ALTER BACKUP SET <name> APPLY BACKUP POLICY <policy_name> [ FORCE ]
//	ALTER BACKUP SET <name> { SUSPEND | RESUME } BACKUP [ CREATION | EXPIRATION ] POLICY
//	ALTER BACKUP SET <name> DELETE BACKUP IDENTIFIER '<backup_id>'
//	ALTER BACKUP SET <name> MODIFY BACKUP IDENTIFIER '<backup_id>' { ADD | REMOVE } LEGAL HOLD
//	ALTER BACKUP SET <name> MODIFY BACKUP IDENTIFIER '<backup_id>' SET COMMENT = '<string_literal>'
//	ALTER BACKUP SET <name> MODIFY BACKUP IDENTIFIER '<backup_id>' UNSET COMMENT
//	ALTER BACKUP SET <name> RENAME TO <new_name>
//	ALTER BACKUP SET <name> SET COMMENT = '<string_literal>'
//	ALTER BACKUP SET <name> UNSET COMMENT
//	ALTER BACKUP SET <name> SET TAG <tag_name> = '<tag_value>' [ , ... ]
//	ALTER BACKUP SET <name> UNSET TAG <tag_name> [ , ... ]
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-backup-set
func parseAlterBackupSet(p *Parser, s *ast.ObjectStmt) error {
	return p.parseAlterDataProtectionSet(s, "BACKUP", actSet|actUnset|actTags|actRename)
}

// parseAlterSnapshotSet reads what follows the name of
//
//	ALTER SNAPSHOT SET <name> ADD SNAPSHOT
//	ALTER SNAPSHOT SET <name> APPLY SNAPSHOT POLICY <policy_name> [ FORCE ]
//	ALTER SNAPSHOT SET <name> { SUSPEND | RESUME } SNAPSHOT [ CREATION | EXPIRATION ] POLICY
//	ALTER SNAPSHOT SET <name> DELETE SNAPSHOT IDENTIFIER '<snapshot_id>'
//	ALTER SNAPSHOT SET <name> MODIFY SNAPSHOT IDENTIFIER '<snapshot_id>' { ADD | REMOVE } LEGAL HOLD
//	ALTER SNAPSHOT SET <name> MODIFY SNAPSHOT IDENTIFIER '<snapshot_id>' SET COMMENT = '<string_literal>'
//	ALTER SNAPSHOT SET <name> MODIFY SNAPSHOT IDENTIFIER '<snapshot_id>' UNSET COMMENT
//	ALTER SNAPSHOT SET <name> SET COMMENT = '<string_literal>'
//	ALTER SNAPSHOT SET <name> UNSET COMMENT
//	ALTER SNAPSHOT SET <name> SET TAG <tag_name> = '<tag_value>' [ , ... ]
//	ALTER SNAPSHOT SET <name> UNSET TAG <tag_name> [ , ... ]
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-snapshot-set
func parseAlterSnapshotSet(p *Parser, s *ast.ObjectStmt) error {
	return p.parseAlterDataProtectionSet(s, "SNAPSHOT", actSet|actUnset|actTags)
}

// parseAlterDataProtectionSet reads one action of ALTER BACKUP SET or ALTER
// SNAPSHOT SET. word is BACKUP or SNAPSHOT; allow names the shared actions
// the page documents.
func (p *Parser) parseAlterDataProtectionSet(s *ast.ObjectStmt, word string, allow objectActions) error {
	start := p.cur.Loc.Start
	switch {
	case p.acceptWords("ADD", word):
		p.addObjectKeywordsClause(s, "ADD "+word, start)
		return nil
	case p.acceptWords("APPLY", word, "POLICY"):
		name, err := p.parseObjectName()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "APPLY "+word+" POLICY", start).Names = []*ast.ObjectName{name}
		if start := p.cur.Loc.Start; p.acceptWords("FORCE") {
			p.addObjectKeywordsClause(s, "FORCE", start)
		}
		return nil
	case p.atWords("SUSPEND", word), p.atWords("RESUME", word):
		keywords := strings.ToUpper(p.advance().Str) + " " + word
		p.advance() // consume word
		for _, w := range []string{"CREATION", "EXPIRATION"} {
			if p.acceptWords(w) {
				keywords += " " + w
				break
			}
		}
		if err := p.expectWords("POLICY"); err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, keywords+" POLICY", start)
		return nil
	case p.acceptWords("DELETE", word, "IDENTIFIER"):
		return p.parseDataProtectionIdentifier(s, "DELETE "+word+" IDENTIFIER", start)
	case p.acceptWords("MODIFY", word, "IDENTIFIER"):
		if err := p.parseDataProtectionIdentifier(s, "MODIFY "+word+" IDENTIFIER", start); err != nil {
			return err
		}
		start := p.cur.Loc.Start
		switch {
		case p.acceptWords("ADD", "LEGAL", "HOLD"):
			p.addObjectKeywordsClause(s, "ADD LEGAL HOLD", start)
			return nil
		case p.acceptWords("REMOVE", "LEGAL", "HOLD"):
			p.addObjectKeywordsClause(s, "REMOVE LEGAL HOLD", start)
			return nil
		}
		return p.parseAlterObjectActions(s, actSet|actUnset)
	}
	return p.parseAlterObjectActions(s, allow)
}

// parseDataProtectionIdentifier reads the '<id>' after { DELETE | MODIFY }
// <word> IDENTIFIER into a clause of keywords, whose one nameless option holds
// the string.
func (p *Parser) parseDataProtectionIdentifier(s *ast.ObjectStmt, keywords string, start int) error {
	if p.cur.Type != tokString {
		return p.syntaxErrorAtCur()
	}
	lit, err := p.parseCopyLiteral()
	if err != nil {
		return err
	}
	c := p.addObjectKeywordsClause(s, keywords, start)
	c.Options = []*ast.CopyOption{{Lit: lit, Loc: lit.Loc}}
	return nil
}

// parseAlterBackupPolicy reads what follows the name of
//
//	ALTER BACKUP POLICY <name> RENAME TO <new_name>
//	ALTER BACKUP POLICY <name> SET
//	  [ COMMENT = '<string_literal>' ] [ SCHEDULE = '...' ] [ EXPIRE_AFTER_DAYS = <days_integer> ]
//	ALTER BACKUP POLICY <name> UNSET { COMMENT | SCHEDULE | EXPIRE_AFTER_DAYS }
//	ALTER BACKUP POLICY <name> SET TAG <tag_name> = '<tag_value>' [ , ... ]
//	ALTER BACKUP POLICY <name> UNSET TAG <tag_name> [ , ... ]
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-backup-policy
func parseAlterBackupPolicy(p *Parser, s *ast.ObjectStmt) error {
	return p.parseAlterObjectActions(s, actSet|actUnset|actTags|actRename)
}

// parseAlterSnapshotPolicy reads what follows the name of
//
//	ALTER SNAPSHOT POLICY <name> RENAME TO <new_name>
//	ALTER SNAPSHOT POLICY <name> SET
//	  [ COMMENT = '<string_literal>' ] [ SCHEDULE = '...' ] [ EXPIRE_AFTER_DAYS = <days_integer> ]
//	ALTER SNAPSHOT POLICY <name> UNSET { COMMENT | SCHEDULE | EXPIRE_AFTER_DAYS }
//	ALTER SNAPSHOT POLICY <name> SET TAG <tag_name> = '<tag_value>' [ , ... ]
//	ALTER SNAPSHOT POLICY <name> UNSET TAG <tag_name> [ , ... ]
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-snapshot-policy
func parseAlterSnapshotPolicy(p *Parser, s *ast.ObjectStmt) error {
	return p.parseAlterObjectActions(s, actSet|actUnset|actTags|actRename)
}

// parseAlterObjectActions reads one of the shared ALTER actions allow names,
// or fails. Without actSetTag, a SET that names tags (SET TAG ..., or TAG
// among SET's properties) fails too.
func (p *Parser) parseAlterObjectActions(s *ast.ObjectStmt, allow objectActions) error {
	if allow&actSetTag == 0 && p.cur.Type == kwSET && p.peekNext().Type == kwTAG {
		p.advance() // consume SET
		return p.syntaxErrorAtCur()
	}
	ok, err := p.parseObjectAlterCommon(s, allow)
	if err != nil {
		return err
	}
	if !ok {
		return p.syntaxErrorAtCur()
	}
	if c := s.Clauses[len(s.Clauses)-1]; allow&actSetTag == 0 && c.Keywords == "SET" && len(c.Tags) > 0 {
		return &ParseError{
			Position: c.Tags[0].Name.Loc.Start, End: c.Tags[0].Name.Loc.End,
			Message: "syntax error: " + s.Kind + " does not take tags",
		}
	}
	return nil
}

// parseCreateSnapshot reads what follows the name of
//
//	CREATE [ OR REPLACE ] SNAPSHOT [ IF NOT EXISTS ] <name>
//	  FROM SERVICE <service_name>
//	  VOLUME "<volume_name>"
//	  INSTANCE <instance_id>
//	  [ COMMENT = '<string_literal>']
//	  [ [ WITH ] TAG ( <tag_name> = '<tag_value>' [ , ... ] ) ]
//
// <instance_id> is the instance's index, an integer from 0 (the page's
// example writes INSTANCE 0); the clause INSTANCE holds it as its one nameless
// option.
//
// https://docs.snowflake.com/en/sql-reference/sql/create-snapshot
func parseCreateSnapshot(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	if err := p.expectWords("FROM", "SERVICE"); err != nil {
		return err
	}
	service, err := p.parseObjectName()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, "FROM SERVICE", start).Names = []*ast.ObjectName{service}

	start = p.cur.Loc.Start
	if err := p.expectWords("VOLUME"); err != nil {
		return err
	}
	volume, err := p.parseObjectSimpleName()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, "VOLUME", start).Names = []*ast.ObjectName{volume}

	start = p.cur.Loc.Start
	if err := p.expectWords("INSTANCE"); err != nil {
		return err
	}
	if p.cur.Type != tokInt {
		return p.syntaxErrorAtCur()
	}
	lit, err := p.parseCopyLiteral()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, "INSTANCE", start).Options = []*ast.CopyOption{{Lit: lit, Loc: lit.Loc}}

	_, err = p.parseObjectProperties(s, "", false)
	return err
}

// parseAlterSnapshot reads what follows the name of
//
//	ALTER SNAPSHOT [ IF EXISTS ] <name> SET COMMENT = '<string_literal>'
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-snapshot
func parseAlterSnapshot(p *Parser, s *ast.ObjectStmt) error {
	return p.parseAlterObjectActions(s, actSet)
}

// parseAlterStorageLifecyclePolicy reads what follows the name of
//
//	ALTER STORAGE LIFECYCLE POLICY [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER STORAGE LIFECYCLE POLICY [ IF EXISTS ] <name> SET ARCHIVE_TIER = { COOL | COLD }
//	ALTER STORAGE LIFECYCLE POLICY [ IF EXISTS ] <name> SET ARCHIVE_FOR_DAYS = <number>
//	ALTER STORAGE LIFECYCLE POLICY [ IF EXISTS ] <name> SET COMMENT = '<string_literal>'
//	ALTER STORAGE LIFECYCLE POLICY [ IF EXISTS ] <name> SET TAG <tag_name> = '<tag_value>' [ , ... ]
//	ALTER STORAGE LIFECYCLE POLICY [ IF EXISTS ] <name> UNSET { ARCHIVE_FOR_DAYS | COMMENT }
//	ALTER STORAGE LIFECYCLE POLICY [ IF EXISTS ] <name> UNSET TAG <tag_name> [ , ... ]
//	ALTER STORAGE LIFECYCLE POLICY [ IF EXISTS ] <name> SET BODY -> <expression_on_arg_name>
//
// SET BODY is a clause whose Body holds the expression.
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-storage-lifecycle-policy
func parseAlterStorageLifecyclePolicy(p *Parser, s *ast.ObjectStmt) error {
	if start := p.cur.Loc.Start; p.acceptWords("SET", "BODY") {
		body, err := p.parseArrowBody(false)
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "SET BODY", start).Body = body
		return nil
	}
	return p.parseAlterObjectActions(s, actSet|actUnset|actTags|actRename)
}

// parseCreateDataMovementPolicy reads what follows the name of
//
//	CREATE [ OR REPLACE ] DATA MOVEMENT POLICY [ IF NOT EXISTS ] <name>
//	  [ ENFORCE_RULES = ( <rule_name> [ , ... ] ) ]
//	  [ ALERT_RULES = ( <rule_name> [ , ... ] ) ]
//	  [ COMMENT = '<string_literal>' ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-data-movement-policy
func parseCreateDataMovementPolicy(p *Parser, s *ast.ObjectStmt) error {
	c, err := p.parseObjectProperties(s, "", false)
	if err != nil || c == nil {
		return err
	}
	if len(c.Tags) > 0 {
		return &ParseError{
			Position: c.Tags[0].Name.Loc.Start, End: c.Tags[0].Name.Loc.End,
			Message: "syntax error: DATA MOVEMENT POLICY does not take tags",
		}
	}
	return nil
}

// parseAlterDataMovementPolicy reads what follows the name of
//
//	ALTER DATA MOVEMENT POLICY [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER DATA MOVEMENT POLICY [ IF EXISTS ] <name> SET { ENFORCE_RULES | ALERT_RULES } = ( <rule_name> [ , ... ] )
//	ALTER DATA MOVEMENT POLICY [ IF EXISTS ] <name> { ADD | REMOVE } { ENFORCE_RULES | ALERT_RULES } = ( <rule_name> [ , ... ] )
//	ALTER DATA MOVEMENT POLICY [ IF EXISTS ] <name> SET COMMENT = '<string_literal>'
//	ALTER DATA MOVEMENT POLICY [ IF EXISTS ] <name> UNSET COMMENT
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-data-movement-policy
func parseAlterDataMovementPolicy(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	for _, verb := range []string{"ADD", "REMOVE"} {
		if !p.acceptWords(verb) {
			continue
		}
		if !p.curIsWord("ENFORCE_RULES") && !p.curIsWord("ALERT_RULES") || !p.startsObjectProperty() {
			return p.syntaxErrorAtCur()
		}
		opt, err := p.parseObjectProperty()
		if err != nil {
			return err
		}
		// The registry's property check covers SET and the CREATE run only.
		if !propNames.matches(opt) {
			return &ParseError{
				Position: opt.Loc.Start, End: opt.Loc.End,
				Message: "syntax error: " + opt.Name + " takes ( <rule_name> [ , ... ] )",
			}
		}
		p.addObjectKeywordsClause(s, verb, start).Options = []*ast.CopyOption{opt}
		return nil
	}
	return p.parseAlterObjectActions(s, actSet|actUnset|actRename)
}

// parseAlterDataMovementRule reads what follows the name of
//
//	ALTER DATA MOVEMENT RULE [ IF EXISTS ] <name> SET COMMENT = '<string_literal>'
//	ALTER DATA MOVEMENT RULE [ IF EXISTS ] <name> UNSET COMMENT
//	ALTER DATA MOVEMENT RULE [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER DATA MOVEMENT RULE [ IF EXISTS ] <name> SET MAX_ROWS AS () RETURNS INTEGER -> ( <expression> )
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-data-movement-rule
func parseAlterDataMovementRule(p *Parser, s *ast.ObjectStmt) error {
	if p.cur.Type == kwSET && isWordToken(p.peekNext(), "MAX_ROWS") {
		start := p.advance().Loc.Start // consume SET
		return p.parseMaxRowsBody(s, "SET MAX_ROWS AS", start)
	}
	return p.parseAlterObjectActions(s, actSet|actUnset|actRename)
}
