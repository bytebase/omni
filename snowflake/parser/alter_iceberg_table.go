package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// ALTER ICEBERG TABLE actions that ALTER TABLE does not have
//
//	https://docs.snowflake.com/en/sql-reference/sql/alter-iceberg-table
//	https://docs.snowflake.com/en/sql-reference/sql/alter-iceberg-table-refresh
//	https://docs.snowflake.com/en/sql-reference/sql/alter-iceberg-table-convert-to-managed
//	https://docs.snowflake.com/en/sql-reference/sql/alter-iceberg-table-optimize
//	https://docs.snowflake.com/en/sql-reference/sql/alter-iceberg-table-partition-evolution
//	https://docs.snowflake.com/en/sql-reference/sql/alter-iceberg-table-version-tag
// ---------------------------------------------------------------------------

// convertToManagedProperties are the properties CONVERT TO MANAGED takes.
var convertToManagedProperties = map[string]bool{
	"FILE_PATH":                    true,
	"BASE_LOCATION":                true,
	"DROP_FROM_SOURCE_CATALOG":     true,
	"STORAGE_SERIALIZATION_POLICY": true,
}

// seqIsWords reports whether cur and the tokens after it spell words, compared
// by their uppercased text (keywords and identifiers alike). The parser state is
// restored before returning.
func (p *Parser) seqIsWords(words ...string) bool {
	snap := p.snapshot()
	defer p.restore(snap)
	for _, w := range words {
		if !p.curIsWord(w) {
			return false
		}
		p.advance()
	}
	return true
}

// parseAlterIcebergOnlyAction parses the action at cur when it is one only ALTER
// ICEBERG TABLE has, and reports whether it did:
//
//	REFRESH [ '<relative_path>' ]
//	CONVERT TO MANAGED [ <property> = <value> ... ]
//	OPTIMIZE ( REMOVE_POSITION_DELETES )
//	{ ADD | DROP } PARTITION BY ( <partitionExpression> )
//	REPLACE PARTITION BY ( <partitionExpression> ) WITH ( <partitionExpression> )
//	CREATE [ OR REPLACE ] VERSION_TAG [ IF NOT EXISTS ] <name> [ AS OF ... [ RETAIN ... ] ]
//	ALTER VERSION_TAG <name> AS OF ...
//	REPLACE VERSION_TAG <name> AS OF ... [ RETAIN ... ]
//	DROP VERSION_TAG [ IF EXISTS ] <name>
//	UPGRADE TO ICEBERG VERSION <n>
func (p *Parser) parseAlterIcebergOnlyAction(action *ast.AlterTableAction) (bool, error) {
	switch {
	case p.cur.Type == kwREFRESH:
		p.advance() // consume REFRESH
		action.Kind = ast.AlterTableRefresh
		if p.cur.Type == tokString {
			s := p.advance().Str
			action.RefreshPath = &s
		}
		return true, nil

	case p.seqIsWords("CONVERT", "TO", "MANAGED"):
		p.advance() // consume CONVERT
		p.advance() // consume TO
		p.advance() // consume MANAGED
		action.Kind = ast.AlterTableConvertToManaged
		for p.startsAssignment() {
			if !convertToManagedProperties[strings.ToUpper(p.cur.Str)] {
				return true, p.syntaxErrorAtCur()
			}
			opt, err := p.parseAssignment()
			if err != nil {
				return true, err
			}
			action.Options = append(action.Options, opt)
		}
		return true, nil

	case p.curIsWord("OPTIMIZE") && p.peekNext().Type == '(':
		p.advance() // consume OPTIMIZE
		p.advance() // consume '('
		for {
			if !p.curIsWord("REMOVE_POSITION_DELETES") {
				return true, p.syntaxErrorAtCur()
			}
			p.advance()
			action.OptimizeOps = append(action.OptimizeOps, "REMOVE_POSITION_DELETES")
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
		}
		if _, err := p.expect(')'); err != nil {
			return true, err
		}
		action.Kind = ast.AlterTableOptimize
		return true, nil

	case p.seqIsWords("ADD", "PARTITION", "BY"), p.seqIsWords("DROP", "PARTITION", "BY"), p.seqIsWords("REPLACE", "PARTITION", "BY"):
		verb := p.advance() // consume ADD / DROP / REPLACE
		p.advance()         // consume PARTITION
		p.advance()         // consume BY
		part, err := p.parseIcebergPartitionParen()
		if err != nil {
			return true, err
		}
		action.Partition = part
		switch verb.Type {
		case kwADD:
			action.Kind = ast.AlterTableAddPartitionBy
		case kwDROP:
			action.Kind = ast.AlterTableDropPartitionBy
		default:
			if _, err := p.expect(kwWITH); err != nil {
				return true, err
			}
			newPart, err := p.parseIcebergPartitionParen()
			if err != nil {
				return true, err
			}
			action.Kind = ast.AlterTableReplacePartitionBy
			action.NewPartition = newPart
		}
		return true, nil

	case p.cur.Type == kwCREATE && (p.peekIsWord("VERSION_TAG") || p.peekNext().Type == kwOR):
		p.advance() // consume CREATE
		if p.cur.Type == kwOR {
			p.advance() // consume OR
			if _, err := p.expect(kwREPLACE); err != nil {
				return true, err
			}
			action.OrReplace = true
		}
		if !p.curIsWord("VERSION_TAG") {
			return true, p.syntaxErrorAtCur()
		}
		p.advance() // consume VERSION_TAG
		if !action.OrReplace && p.cur.Type == kwIF {
			p.advance() // consume IF
			if _, err := p.expect(kwNOT); err != nil {
				return true, err
			}
			if _, err := p.expect(kwEXISTS); err != nil {
				return true, err
			}
			action.IfNotExists = true
		}
		name, err := p.parseObjectName()
		if err != nil {
			return true, err
		}
		action.Kind = ast.AlterTableCreateVersionTag
		action.VersionTag = name
		if p.cur.Type == kwAS {
			if err := p.parseVersionTagAsOf(action, true); err != nil {
				return true, err
			}
		}
		return true, nil

	case (p.cur.Type == kwALTER || p.cur.Type == kwREPLACE) && p.peekIsWord("VERSION_TAG") && p.versionTagNameFollows(kwAS):
		verb := p.advance() // consume ALTER / REPLACE
		p.advance()         // consume VERSION_TAG
		name, err := p.parseObjectName()
		if err != nil {
			return true, err
		}
		action.VersionTag = name
		action.Kind = ast.AlterTableAlterVersionTag
		if verb.Type == kwREPLACE {
			action.Kind = ast.AlterTableReplaceVersionTag
		}
		if err := p.parseVersionTagAsOf(action, verb.Type == kwREPLACE); err != nil {
			return true, err
		}
		return true, nil

	case p.cur.Type == kwDROP && p.peekIsWord("VERSION_TAG") && p.versionTagNameFollows(0):
		p.advance() // consume DROP
		p.advance() // consume VERSION_TAG
		if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
			p.advance() // consume IF
			p.advance() // consume EXISTS
			action.IfExists = true
		}
		name, err := p.parseObjectName()
		if err != nil {
			return true, err
		}
		action.Kind = ast.AlterTableDropVersionTag
		action.VersionTag = name
		return true, nil

	case p.seqIsWords("UPGRADE", "TO", "ICEBERG", "VERSION"):
		for range 4 {
			p.advance() // consume UPGRADE TO ICEBERG VERSION
		}
		tok, err := p.expect(tokInt)
		if err != nil {
			return true, err
		}
		action.Kind = ast.AlterTableUpgradeIcebergVersion
		action.IcebergVersion = tok.Str
		return true, nil
	}
	return false, nil
}

// versionTagNameFollows reports whether the VERSION_TAG after cur (cur is the
// verb) is followed by a tag name and then, when then is non-zero, by that
// token: `DROP VERSION_TAG [IF EXISTS] name`, `ALTER VERSION_TAG name AS`.
// Otherwise VERSION_TAG is a column name (DROP version_tag, ALTER version_tag
// SET NOT NULL).
func (p *Parser) versionTagNameFollows(then int) bool {
	snap := p.snapshot()
	defer p.restore(snap)
	p.advance() // consume the verb
	p.advance() // consume VERSION_TAG
	if then == 0 && p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		return true
	}
	if !p.isColumnNameToken(p.cur.Type) && p.cur.Type != tokQuotedIdent {
		return false
	}
	if _, err := p.parseObjectName(); err != nil {
		return false
	}
	if then == 0 {
		return p.cur.Type != ','
	}
	return p.cur.Type == then
}

// parseVersionTagAsOf parses AS OF { VERSION <id> | TIMESTAMP <n> } and, with
// retain, an optional RETAIN <n> { DAYS | HOURS | MINUTES }. cur is AS.
func (p *Parser) parseVersionTagAsOf(action *ast.AlterTableAction, retain bool) error {
	if _, err := p.expect(kwAS); err != nil {
		return err
	}
	if _, err := p.expect(kwOF); err != nil {
		return err
	}
	switch p.cur.Type {
	case kwVERSION:
		p.advance() // consume VERSION
		action.AsOfKind = "VERSION"
		if p.cur.Type == tokInt {
			action.AsOfValue = p.advance().Str
		} else {
			id, err := p.parseIdent()
			if err != nil {
				return err
			}
			action.AsOfValue = id.String()
		}
	case kwTIMESTAMP:
		p.advance() // consume TIMESTAMP
		tok, err := p.expect(tokInt)
		if err != nil {
			return err
		}
		action.AsOfKind = "TIMESTAMP"
		action.AsOfValue = tok.Str
	default:
		return p.syntaxErrorAtCur()
	}
	if retain && p.curIsWord("RETAIN") {
		p.advance() // consume RETAIN
		tok, err := p.expect(tokInt)
		if err != nil {
			return err
		}
		n := tok.Ival
		action.RetainCount = &n
		switch {
		case p.curIsWord("DAYS"), p.curIsWord("HOURS"), p.curIsWord("MINUTES"):
			action.RetainUnit = strings.ToUpper(p.advance().Str)
		default:
			return p.syntaxErrorAtCur()
		}
	}
	return nil
}

// parseIcebergPartitionParen parses ( <partitionExpression> ): exactly one
// expression, as the partition-evolution page's usage notes say. The
// expression is read by CREATE ICEBERG TABLE's parseIcebergPartitionExpr.
func (p *Parser) parseIcebergPartitionParen() (ast.Node, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	part, err := p.parseIcebergPartitionExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return part, nil
}
