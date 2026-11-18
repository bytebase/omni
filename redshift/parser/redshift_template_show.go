package parser

import (
	"strconv"

	nodes "github.com/bytebase/omni/redshift/ast"
)

// redshiftReservedWords is the exact Redshift reserved-word set from
// https://docs.aws.amazon.com/redshift/latest/dg/r_pg_keywords.html
// (mirrored in redshift/testdata/tokens.tsv). An unquoted word in this set
// cannot stand as an identifier (ID) in Redshift DDL or SHOW statements.
var redshiftReservedWords = map[string]bool{
	"aes128": true, "aes256": true, "all": true, "allowoverwrite": true,
	"analyse": true, "analyze": true, "and": true, "any": true, "array": true,
	"as": true, "asc": true, "authorization": true, "az64": true,
	"backup": true, "between": true, "binary": true, "blanksasnull": true,
	"both": true, "bytedict": true, "bzip2": true, "case": true, "cast": true,
	"check": true, "collate": true, "column": true, "constraint": true,
	"create": true, "credentials": true, "cross": true, "current_date": true,
	"current_time": true, "current_timestamp": true, "current_user": true,
	"current_user_id": true, "default": true, "deferrable": true,
	"deflate": true, "defrag": true, "delta": true, "delta32k": true,
	"desc": true, "disable": true, "distinct": true, "do": true, "else": true,
	"emptyasnull": true, "enable": true, "encode": true, "encrypt": true,
	"encryption": true, "end": true, "except": true, "explicit": true,
	"false": true, "for": true, "foreign": true, "freeze": true, "from": true,
	"full": true, "globaldict256": true, "globaldict64k": true, "grant": true,
	"group": true, "gzip": true, "having": true, "identity": true,
	"ignore": true, "ilike": true, "in": true, "initially": true, "inner": true,
	"intersect": true, "interval": true, "into": true, "is": true,
	"isnull": true, "join": true, "leading": true, "left": true, "like": true,
	"limit": true, "localtime": true, "localtimestamp": true, "lun": true,
	"luns": true, "lzo": true, "lzop": true, "minus": true, "mostly16": true,
	"mostly32": true, "mostly8": true, "natural": true, "new": true,
	"not": true, "notnull": true, "null": true, "nulls": true, "off": true,
	"offline": true, "offset": true, "oid": true, "old": true, "on": true,
	"only": true, "open": true, "or": true, "order": true, "outer": true,
	"overlaps": true, "parallel": true, "partition": true, "percent": true,
	"permissions": true, "pivot": true, "placing": true, "primary": true,
	"raw": true, "readratio": true, "recover": true, "references": true,
	"rejectlog": true, "resort": true, "respect": true, "restore": true,
	"right": true, "select": true, "session_user": true, "similar": true,
	"snapshot": true, "some": true, "sysdate": true, "system": true,
	"table": true, "tag": true, "tdes": true, "text255": true, "text32k": true,
	"then": true, "timestamp": true, "to": true, "top": true, "trailing": true,
	"true": true, "truncatecolumns": true, "union": true, "unique": true,
	"unnest": true, "unpivot": true, "user": true, "using": true,
	"verbose": true, "wallet": true, "when": true, "where": true,
	"with": true, "without": true,
}

// redshiftCopyReservedOptionWords are COPY option and format words that
// Redshift reserves as SQL identifiers yet accepts as COPY options inside
// CREATE TEMPLATE ... FOR COPY and ALTER TEMPLATE ... ADD / DROP.
var redshiftCopyReservedOptionWords = map[string]bool{
	"aes128": true, "aes256": true, "allowoverwrite": true,
	"blanksasnull": true, "bzip2": true, "deflate": true, "emptyasnull": true,
	"encrypt": true, "encryption": true, "explicit": true, "freeze": true,
	"gzip": true, "ignore": true, "lzo": true, "lzop": true, "null": true,
	"off": true, "on": true, "parallel": true, "raw": true, "restore": true,
	"tdes": true, "truncatecolumns": true, "wallet": true,
}

// asciiLower folds ASCII A-Z to a-z without touching non-ASCII bytes, matching
// the lexer's identifier and keyword folding (ident_fold_test.go).
func asciiLower(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			b := []byte(s)
			for ; i < len(b); i++ {
				if b[i] >= 'A' && b[i] <= 'Z' {
					b[i] += 'a' - 'A'
				}
			}
			return string(b)
		}
	}
	return s
}

func asciiEqualFold(s, lowerWord string) bool {
	if len(s) != len(lowerWord) {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		if c != lowerWord[i] {
			return false
		}
	}
	return true
}

// redshiftTokIsWord reports whether tok is the unquoted keyword or identifier
// word (never a string literal or delimited identifier), folded on ASCII only.
func (p *Parser) redshiftTokIsWord(tok Token, lowerWord string) bool {
	if tok.Type != IDENT && tok.Type < NOT_LA {
		return false
	}
	if p.lexer != nil && tok.End-tok.Loc != len(lowerWord) {
		return false
	}
	return asciiEqualFold(tok.Str, lowerWord)
}

func (p *Parser) redshiftCurIsWord(lowerWord string) bool {
	return p.redshiftTokIsWord(p.cur, lowerWord)
}

func (p *Parser) expectRedshiftUnquotedWord(lowerWord string) error {
	if !p.redshiftCurIsWord(lowerWord) {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

func (p *Parser) isDelimitedIdentToken(tok Token) bool {
	return tok.Type == IDENT && p.lexer != nil && tok.End-tok.Loc != len(tok.Str)
}

// isRedshiftDocIdent reports whether cur is a valid Redshift identifier (ID):
// a delimited identifier ("..."), or an unquoted identifier / non-reserved word.
func (p *Parser) isRedshiftDocIdent() bool {
	if p.isDelimitedIdentToken(p.cur) {
		return true
	}
	if p.cur.Type != IDENT && p.cur.Type < NOT_LA {
		return false
	}
	return !redshiftReservedWords[asciiLower(p.cur.Str)]
}

func (p *Parser) parseRedshiftDocIdent() (string, error) {
	if !p.isRedshiftDocIdent() {
		return "", p.syntaxErrorAtCur()
	}
	tok := p.advance()
	return tok.Str, nil
}

// parseRedshiftDottedDocIdent reads a dotted identifier with between minParts
// and maxParts parts (inclusive).
func (p *Parser) parseRedshiftDottedDocIdent(minParts, maxParts int) (*nodes.List, error) {
	first, err := p.parseRedshiftDocIdent()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{&nodes.String{Str: first}}
	for p.cur.Type == '.' {
		if len(items) >= maxParts {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume '.'
		part, err := p.parseRedshiftDocIdent()
		if err != nil {
			return nil, err
		}
		items = append(items, &nodes.String{Str: part})
	}
	if len(items) < minParts {
		return nil, p.syntaxErrorAtCur()
	}
	return &nodes.List{Items: items}, nil
}

// ---------------------------------------------------------------------------
// CREATE / ALTER / DROP TEMPLATE
// ---------------------------------------------------------------------------

// parseRedshiftCreateTemplateStmt parses:
//
//	CREATE [ OR REPLACE ] TEMPLATE [ database_name. ][ schema_name. ]template_name
//	  FOR COPY [ AS ] [ [ FORMAT ] [ AS ] data_format ]
//	  [ copy_option [ value ] [, ...] ]
//
// Ref: https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_TEMPLATE.html
//
// CREATE [ OR REPLACE ] has already been consumed; cur is TEMPLATE.
func (p *Parser) parseRedshiftCreateTemplateStmt(command string, loc int) (nodes.Node, error) {
	if err := p.expectRedshiftUnquotedWord("template"); err != nil {
		return nil, err
	}
	name, err := p.parseRedshiftDottedDocIdent(1, 3)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(FOR); err != nil {
		return nil, err
	}
	if err := p.expectRedshiftUnquotedWord("copy"); err != nil {
		return nil, err
	}

	var opts []nodes.Node
	if p.cur.Type == AS {
		p.advance() // optional AS after FOR COPY
		opts = append(opts, &nodes.String{Str: "as"})
	}

	if p.cur.Type != 0 && p.cur.Type != ';' {
		switch {
		case p.redshiftCurIsWord("format"):
			p.advance() // consume FORMAT
			opts = append(opts, &nodes.String{Str: "format"})
			if p.cur.Type == AS {
				p.advance()
				opts = append(opts, &nodes.String{Str: "as"})
			}
			fmtVal, _, err := p.parseRedshiftTemplateOption(false)
			if err != nil {
				return nil, err
			}
			opts = append(opts, &nodes.String{Str: fmtVal})
			if p.cur.Type != 0 && p.cur.Type != ';' {
				more, err := p.parseRedshiftTemplateCopyOptionList()
				if err != nil {
					return nil, err
				}
				opts = append(opts, more...)
			}
		case p.cur.Type == AS:
			p.advance() // AS before data_format when FORMAT was omitted
			opts = append(opts, &nodes.String{Str: "as"})
			fmtVal, _, err := p.parseRedshiftTemplateOption(false)
			if err != nil {
				return nil, err
			}
			opts = append(opts, &nodes.String{Str: fmtVal})
			if p.cur.Type != 0 && p.cur.Type != ';' {
				more, err := p.parseRedshiftTemplateCopyOptionList()
				if err != nil {
					return nil, err
				}
				opts = append(opts, more...)
			}
		default:
			// Without FORMAT or a second AS, the leading segment may hold
			// [data_format] followed by copy_option [value]: 1 to 3 tokens
			// before any comma, and 1 to 2 tokens after each comma.
			first, _, err := p.parseRedshiftTemplateOption(false)
			if err != nil {
				return nil, err
			}
			opts = append(opts, &nodes.String{Str: first})
			if p.cur.Type != 0 && p.cur.Type != ';' && p.cur.Type != ',' {
				second, secondIsWord, err := p.parseRedshiftTemplateOption(true)
				if err != nil {
					return nil, err
				}
				opts = append(opts, &nodes.String{Str: second})
				if secondIsWord && p.cur.Type != 0 && p.cur.Type != ';' && p.cur.Type != ',' {
					third, _, err := p.parseRedshiftTemplateOption(true)
					if err != nil {
						return nil, err
					}
					opts = append(opts, &nodes.String{Str: third})
				}
			}
			for p.cur.Type == ',' {
				p.advance() // consume ','
				opt, _, err := p.parseRedshiftTemplateOption(false)
				if err != nil {
					return nil, err
				}
				opts = append(opts, &nodes.String{Str: opt})
				if p.cur.Type != 0 && p.cur.Type != ';' && p.cur.Type != ',' {
					val, _, err := p.parseRedshiftTemplateOption(true)
					if err != nil {
						return nil, err
					}
					opts = append(opts, &nodes.String{Str: val})
				}
			}
		}
	}

	var optList *nodes.List
	if len(opts) > 0 {
		optList = &nodes.List{Items: opts}
	}
	return &nodes.RedshiftObjectStmt{
		Command:    command,
		ObjectType: "template",
		Name:       name,
		Options:    optList,
		Loc:        nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

func (p *Parser) parseRedshiftTemplateCopyOptionList() ([]nodes.Node, error) {
	var opts []nodes.Node
	for {
		opt, _, err := p.parseRedshiftTemplateOption(false)
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.String{Str: opt})
		if p.cur.Type != 0 && p.cur.Type != ';' && p.cur.Type != ',' {
			val, _, err := p.parseRedshiftTemplateOption(true)
			if err != nil {
				return nil, err
			}
			opts = append(opts, &nodes.String{Str: val})
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return opts, nil
}

// parseRedshiftTemplateOption reads a COPY option/format name or, when
// allowLiteral is true, an option value (string, integer, DEFAULT, TRUE/FALSE,
// or identifier/COPY keyword). It also reports whether the token was a word
// capable of standing as an option name.
func (p *Parser) parseRedshiftTemplateOption(allowLiteral bool) (string, bool, error) {
	if p.isRedshiftDocIdent() {
		tok := p.advance()
		return tok.Str, true, nil
	}
	if (p.cur.Type == IDENT || p.cur.Type >= NOT_LA) && !p.isDelimitedIdentToken(p.cur) &&
		redshiftCopyReservedOptionWords[asciiLower(p.cur.Str)] {
		tok := p.advance()
		return tok.Str, true, nil
	}
	if allowLiteral {
		switch p.cur.Type {
		case SCONST:
			tok := p.advance()
			return tok.Str, false, nil
		case ICONST:
			tok := p.advance()
			return strconv.FormatInt(tok.Ival, 10), false, nil
		case DEFAULT, TRUE_P, FALSE_P:
			tok := p.advance()
			return asciiLower(tok.Str), false, nil
		}
	}
	return "", false, p.syntaxErrorAtCur()
}

// parseRedshiftAlterTemplateStmt parses:
//
//	ALTER TEMPLATE [ database_name. ][ schema_name. ]template_name
//	  { RENAME TO new_name
//	  | OWNER TO new_owner
//	  | ADD copy_option [ AS ] [ value ]
//	  | DROP copy_option
//	  | SET parameter TO value [, parameter TO value ...] }
//
// Ref: https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_TEMPLATE.html
//
// ALTER has already been consumed; cur is TEMPLATE.
func (p *Parser) parseRedshiftAlterTemplateStmt(loc int) (nodes.Node, error) {
	if err := p.expectRedshiftUnquotedWord("template"); err != nil {
		return nil, err
	}
	name, err := p.parseRedshiftDottedDocIdent(1, 3)
	if err != nil {
		return nil, err
	}

	var opts []nodes.Node
	switch {
	case p.redshiftCurIsWord("rename"):
		p.advance() // consume RENAME
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		target, err := p.parseRedshiftDocIdent()
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.String{Str: "rename"}, &nodes.String{Str: "to"}, &nodes.String{Str: target})

	case p.redshiftCurIsWord("owner"):
		p.advance() // consume OWNER
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		owner, err := p.parseRedshiftDocIdent()
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.String{Str: "owner"}, &nodes.String{Str: "to"}, &nodes.String{Str: owner})

	case p.redshiftCurIsWord("add"):
		p.advance() // consume ADD
		opt, _, err := p.parseRedshiftTemplateOption(false)
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.String{Str: "add"}, &nodes.String{Str: opt})
		if p.cur.Type == AS {
			p.advance()
			opts = append(opts, &nodes.String{Str: "as"})
		}
		if p.cur.Type != 0 && p.cur.Type != ';' {
			val, _, err := p.parseRedshiftTemplateOption(true)
			if err != nil {
				return nil, err
			}
			opts = append(opts, &nodes.String{Str: val})
		}

	case p.redshiftCurIsWord("drop"):
		p.advance() // consume DROP
		opt, _, err := p.parseRedshiftTemplateOption(false)
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.String{Str: "drop"}, &nodes.String{Str: opt})

	case p.redshiftCurIsWord("set"):
		p.advance() // consume SET
		opts = append(opts, &nodes.String{Str: "set"})
		for {
			param, err := p.parseRedshiftDocIdent()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(TO); err != nil {
				return nil, err
			}
			val, _, err := p.parseRedshiftTemplateOption(true)
			if err != nil {
				return nil, err
			}
			opts = append(opts, &nodes.String{Str: param}, &nodes.String{Str: "to"}, &nodes.String{Str: val})
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	return &nodes.RedshiftObjectStmt{
		Command:    "alter",
		ObjectType: "template",
		Name:       name,
		Options:    &nodes.List{Items: opts},
		Loc:        nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// parseRedshiftDropTemplateStmt parses:
//
//	DROP TEMPLATE [ database_name. ][ schema_name. ]template_name
//
// Ref: https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_TEMPLATE.html
//
// DROP has already been consumed; cur is TEMPLATE.
func (p *Parser) parseRedshiftDropTemplateStmt(loc int) (nodes.Node, error) {
	if err := p.expectRedshiftUnquotedWord("template"); err != nil {
		return nil, err
	}
	name, err := p.parseRedshiftDottedDocIdent(1, 3)
	if err != nil {
		return nil, err
	}
	return &nodes.RedshiftObjectStmt{
		Command:    "drop",
		ObjectType: "template",
		Name:       name,
		Loc:        nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// ---------------------------------------------------------------------------
// Documented Redshift SHOW statements
// ---------------------------------------------------------------------------

func (p *Parser) parseRedshiftDocShowStmt(loc int) (nodes.Node, error) {
	switch {
	case p.cur.Type == USER && p.redshiftTokIsWord(p.peekNext(), "lockout"):
		// SHOW USER LOCKOUT
		// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_USER_LOCKOUT.html
		p.advance() // consume USER
		p.advance() // consume LOCKOUT
		return &nodes.RedshiftShowStmt{
			Command: "user lockout",
			Loc:     nodes.Loc{Start: loc, End: p.prev.End},
		}, nil

	case p.cur.Type == COLUMN && p.redshiftTokIsWord(p.peekNext(), "grants"):
		// SHOW COLUMN GRANTS ON TABLE { db.schema.table | schema.table }
		//   [ FOR { username | ROLE role_name | PUBLIC } ] [ LIMIT row_limit ]
		// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_COLUMN_GRANTS.html
		p.advance() // consume COLUMN
		p.advance() // consume GRANTS
		if _, err := p.expect(ON); err != nil {
			return nil, err
		}
		if _, err := p.expect(TABLE); err != nil {
			return nil, err
		}
		tbl, err := p.parseRedshiftDottedDocIdent(2, 3)
		if err != nil {
			return nil, err
		}
		opts := append([]nodes.Node{&nodes.String{Str: "on"}, &nodes.String{Str: "table"}}, tbl.Items...)
		if p.cur.Type == FOR {
			forOpts, err := p.parseRedshiftShowForPrincipal()
			if err != nil {
				return nil, err
			}
			opts = append(opts, forOpts...)
		}
		if p.cur.Type == LIMIT {
			limOpts, err := p.parseRedshiftShowLimit()
			if err != nil {
				return nil, err
			}
			opts = append(opts, limOpts...)
		}
		return &nodes.RedshiftShowStmt{
			Command: "column grants",
			Options: &nodes.List{Items: opts},
			Loc:     nodes.Loc{Start: loc, End: p.prev.End},
		}, nil

	case p.redshiftCurIsWord("constraints"):
		// SHOW CONSTRAINTS { PRIMARY KEYS | FOREIGN KEYS [ EXPORTED ] }
		//   FROM TABLE { db.schema.table | schema.table } [ LIMIT row_limit ]
		// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_CONSTRAINTS.html
		p.advance() // consume CONSTRAINTS
		var opts []nodes.Node
		switch p.cur.Type {
		case PRIMARY:
			p.advance()
			if err := p.expectRedshiftUnquotedWord("keys"); err != nil {
				return nil, err
			}
			opts = append(opts, &nodes.String{Str: "primary"}, &nodes.String{Str: "keys"})
		case FOREIGN:
			p.advance()
			if err := p.expectRedshiftUnquotedWord("keys"); err != nil {
				return nil, err
			}
			opts = append(opts, &nodes.String{Str: "foreign"}, &nodes.String{Str: "keys"})
			if p.redshiftCurIsWord("exported") {
				p.advance()
				opts = append(opts, &nodes.String{Str: "exported"})
			}
		default:
			return nil, p.syntaxErrorAtCur()
		}
		if _, err := p.expect(FROM); err != nil {
			return nil, err
		}
		if _, err := p.expect(TABLE); err != nil {
			return nil, err
		}
		tbl, err := p.parseRedshiftDottedDocIdent(2, 3)
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.String{Str: "from"}, &nodes.String{Str: "table"})
		opts = append(opts, tbl.Items...)
		if p.cur.Type == LIMIT {
			limOpts, err := p.parseRedshiftShowLimit()
			if err != nil {
				return nil, err
			}
			opts = append(opts, limOpts...)
		}
		return &nodes.RedshiftShowStmt{
			Command: "constraints",
			Options: &nodes.List{Items: opts},
			Loc:     nodes.Loc{Start: loc, End: p.prev.End},
		}, nil

	case p.redshiftCurIsWord("functions"):
		// SHOW FUNCTIONS FROM SCHEMA [ db. ]schema [ LIKE 'filter' ] [ LIMIT row_limit ]
		// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_FUNCTIONS.html
		return p.parseRedshiftShowFromSchema("functions", loc)

	case p.redshiftCurIsWord("procedures"):
		// SHOW PROCEDURES FROM SCHEMA [ db. ]schema [ LIKE 'filter' ] [ LIMIT row_limit ]
		// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_PROCEDURES.html
		return p.parseRedshiftShowFromSchema("procedures", loc)

	case p.redshiftCurIsWord("templates"):
		// SHOW TEMPLATES FROM SCHEMA [ db. ]schema [ LIKE 'filter' ] [ LIMIT row_limit ]
		// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_TEMPLATES.html
		return p.parseRedshiftShowFromSchema("templates", loc)

	case p.redshiftCurIsWord("parameters"):
		// SHOW PARAMETERS OF { FUNCTION | PROCEDURE } [ db. ]schema.routine ( arg_type [, ...] )
		//   [ LIKE 'filter' ]
		// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_PARAMETERS.html
		p.advance() // consume PARAMETERS
		if err := p.expectRedshiftUnquotedWord("of"); err != nil {
			return nil, err
		}
		var kind string
		switch {
		case p.redshiftCurIsWord("function"):
			p.advance()
			kind = "function"
		case p.redshiftCurIsWord("procedure"):
			p.advance()
			kind = "procedure"
		default:
			return nil, p.syntaxErrorAtCur()
		}
		routine, err := p.parseRedshiftDottedDocIdent(2, 3)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		opts := []nodes.Node{&nodes.String{Str: "of"}, &nodes.String{Str: kind}}
		opts = append(opts, routine.Items...)
		for {
			if p.cur.Type == SETOF {
				return nil, p.syntaxErrorAtCur()
			}
			tn, err := p.parseTypename()
			if err != nil {
				return nil, err
			}
			opts = append(opts, tn)
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		if p.cur.Type == LIKE {
			likeOpts, err := p.parseRedshiftShowLike()
			if err != nil {
				return nil, err
			}
			opts = append(opts, likeOpts...)
		}
		return &nodes.RedshiftShowStmt{
			Command: "parameters",
			Options: &nodes.List{Items: opts},
			Loc:     nodes.Loc{Start: loc, End: p.prev.End},
		}, nil

	case (p.redshiftCurIsWord("rls") || p.redshiftCurIsWord("masking")) &&
		p.redshiftTokIsWord(p.peekNext(), "policies"):
		// SHOW { RLS | MASKING } POLICIES
		//   [ ON { db.schema.rel | schema.rel } [ FOR { username | ROLE role_name | PUBLIC } ]
		//   | FROM DATABASE db_name ]
		//   [ LIMIT row_limit ]
		// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_POLICIES.html
		kind := asciiLower(p.advance().Str) // consume RLS / MASKING
		p.advance()                         // consume POLICIES
		var opts []nodes.Node
		switch p.cur.Type {
		case ON:
			p.advance() // consume ON
			rel, err := p.parseRedshiftDottedDocIdent(2, 3)
			if err != nil {
				return nil, err
			}
			opts = append(opts, &nodes.String{Str: "on"})
			opts = append(opts, rel.Items...)
			if p.cur.Type == FOR {
				forOpts, err := p.parseRedshiftShowForPrincipal()
				if err != nil {
					return nil, err
				}
				opts = append(opts, forOpts...)
			}
		case FROM:
			p.advance() // consume FROM
			if err := p.expectRedshiftUnquotedWord("database"); err != nil {
				return nil, err
			}
			db, err := p.parseRedshiftDocIdent()
			if err != nil {
				return nil, err
			}
			opts = append(opts, &nodes.String{Str: "from"}, &nodes.String{Str: "database"}, &nodes.String{Str: db})
		}
		if p.cur.Type == LIMIT {
			limOpts, err := p.parseRedshiftShowLimit()
			if err != nil {
				return nil, err
			}
			opts = append(opts, limOpts...)
		}
		var optList *nodes.List
		if len(opts) > 0 {
			optList = &nodes.List{Items: opts}
		}
		return &nodes.RedshiftShowStmt{
			Command: kind + " policies",
			Options: optList,
			Loc:     nodes.Loc{Start: loc, End: p.prev.End},
		}, nil

	case p.redshiftCurIsWord("template"):
		// SHOW TEMPLATE [ database_name. ][ schema_name. ]template_name
		// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_TEMPLATE.html
		p.advance() // consume TEMPLATE
		name, err := p.parseRedshiftDottedDocIdent(1, 3)
		if err != nil {
			return nil, err
		}
		return &nodes.RedshiftShowStmt{
			Command: "template",
			Options: name,
			Loc:     nodes.Loc{Start: loc, End: p.prev.End},
		}, nil

	default:
		return nil, p.syntaxErrorAtCur()
	}
}

func (p *Parser) parseRedshiftShowFromSchema(command string, loc int) (nodes.Node, error) {
	p.advance() // consume FUNCTIONS / PROCEDURES / TEMPLATES
	if _, err := p.expect(FROM); err != nil {
		return nil, err
	}
	if err := p.expectRedshiftUnquotedWord("schema"); err != nil {
		return nil, err
	}
	schema, err := p.parseRedshiftDottedDocIdent(1, 2)
	if err != nil {
		return nil, err
	}
	opts := append([]nodes.Node{&nodes.String{Str: "from"}, &nodes.String{Str: "schema"}}, schema.Items...)
	if p.cur.Type == LIKE {
		likeOpts, err := p.parseRedshiftShowLike()
		if err != nil {
			return nil, err
		}
		opts = append(opts, likeOpts...)
	}
	if p.cur.Type == LIMIT {
		limOpts, err := p.parseRedshiftShowLimit()
		if err != nil {
			return nil, err
		}
		opts = append(opts, limOpts...)
	}
	return &nodes.RedshiftShowStmt{
		Command: command,
		Options: &nodes.List{Items: opts},
		Loc:     nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// parseRedshiftShowForPrincipal reads FOR { username | ROLE role_name | PUBLIC }.
// As in GRANT/REVOKE (grant_role_name_test.go), ROLE is unreserved, so FOR ROLE
// with no following identifier (or before LIMIT) names the user "role".
func (p *Parser) parseRedshiftShowForPrincipal() ([]nodes.Node, error) {
	p.advance() // consume FOR
	if p.redshiftCurIsWord("public") {
		p.advance()
		return []nodes.Node{&nodes.String{Str: "for"}, &nodes.String{Str: "public"}}, nil
	}
	if p.redshiftCurIsWord("role") {
		p.advance() // consume ROLE
		if p.isRedshiftDocIdent() {
			roleName, err := p.parseRedshiftDocIdent()
			if err != nil {
				return nil, err
			}
			return []nodes.Node{&nodes.String{Str: "for"}, &nodes.String{Str: "role"}, &nodes.String{Str: roleName}}, nil
		}
		return []nodes.Node{&nodes.String{Str: "for"}, &nodes.String{Str: "role"}}, nil
	}
	user, err := p.parseRedshiftDocIdent()
	if err != nil {
		return nil, err
	}
	return []nodes.Node{&nodes.String{Str: "for"}, &nodes.String{Str: user}}, nil
}

func (p *Parser) parseRedshiftShowLike() ([]nodes.Node, error) {
	p.advance() // consume LIKE
	if p.cur.Type != SCONST {
		return nil, p.syntaxErrorAtCur()
	}
	pat := p.advance()
	return []nodes.Node{&nodes.String{Str: "like"}, &nodes.String{Str: pat.Str}}, nil
}

func (p *Parser) parseRedshiftShowLimit() ([]nodes.Node, error) {
	p.advance() // consume LIMIT
	if p.cur.Type == '+' {
		p.advance()
	}
	if p.cur.Type != ICONST {
		return nil, p.syntaxErrorAtCur()
	}
	n := p.advance()
	return []nodes.Node{&nodes.String{Str: "limit"}, &nodes.Integer{Ival: n.Ival}}, nil
}
