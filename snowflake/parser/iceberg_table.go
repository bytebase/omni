package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// CREATE ICEBERG TABLE
// ---------------------------------------------------------------------------
//
// An Iceberg table is a table: it parses into ast.CreateTableStmt with Iceberg
// set, so table review rules see its columns and constraints like any other.
// Snowflake documents one CREATE ICEBERG TABLE page per catalog, and they share
// one shape:
//
//	CREATE [ OR REPLACE ] [ TRANSIENT ] ICEBERG TABLE [ IF NOT EXISTS ] <name>
//	  [ ( <col_def> [, ...] [, <out_of_line_constraint> ...] ) ]
//	  [ <properties> ]
//	  [ AS <query> ]
//	CREATE [ OR REPLACE ] [ TRANSIENT ] ICEBERG TABLE <name> LIKE <source> [ <properties> ]
//	CREATE [ OR REPLACE ] ICEBERG TABLE [ IF NOT EXISTS ] <name> CLONE <source> [ ... ] [ <properties> ]
//	CREATE ICEBERG TABLE <name> FROM BACKUP SET <set> IDENTIFIER '<id>'
//
// Snowflake as the catalog:
// https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-snowflake
// AWS Glue: https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-aws-glue
// Iceberg files in object storage:
// https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-iceberg-files
// Delta files: https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-delta
// Parquet files: https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-parquet
// Iceberg REST catalog: https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-rest
// The page that lists every catalog: https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table
//
// The properties are, in any order: PARTITION BY ( ... ), CLUSTER BY ( ... ),
// COPY GRANTS, COPY TAGS, COPY = ( ... ), PURGE, COMMENT = '...', [WITH] TAG
// ( ... ), [WITH] ROW ACCESS POLICY, [WITH] AGGREGATION POLICY, WITH CONTACT
// ( ... ), and the catalog and storage parameters (EXTERNAL_VOLUME, CATALOG,
// BASE_LOCATION, CATALOG_TABLE_NAME, METADATA_FILE_PATH, TARGET_FILE_SIZE, ...),
// which are `KEY = value` pairs kept open-ended as ast.CopyOption.

// parseCreateIcebergTableStmt parses CREATE ICEBERG TABLE. CREATE and the
// optional OR REPLACE / TRANSIENT have been consumed; cur is ICEBERG.
func (p *Parser) parseCreateIcebergTableStmt(start ast.Loc, orReplace, transient bool) (ast.Node, error) {
	p.advance() // consume ICEBERG
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}

	stmt := &ast.CreateTableStmt{
		OrReplace: orReplace,
		Transient: transient,
		Iceberg:   true,
		Loc:       ast.Loc{Start: start.Start},
	}
	if err := p.parseIfNotExistsInto(&stmt.IfNotExists); err != nil {
		return nil, err
	}
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	if p.curIs(kwFROM) {
		if orReplace || transient || stmt.IfNotExists {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume FROM
		backup, err := p.parseBackupSource()
		if err != nil {
			return nil, err
		}
		stmt.FromBackup = backup
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}

	asAllowed := true
	switch {
	case p.curIs(kwLIKE):
		p.advance() // consume LIKE
		src, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Like = src
		asAllowed = false
	case p.curIs(kwCLONE):
		clone, err := p.parseCloneSource()
		if err != nil {
			return nil, err
		}
		stmt.Clone = clone
		asAllowed = false
	case p.curIs('('):
		if err := p.parseColumnDeclItems(stmt); err != nil {
			return nil, err
		}
	}

	if err := p.parseIcebergTableProperties(stmt); err != nil {
		return nil, err
	}

	if asAllowed && p.curIs(kwAS) {
		for _, col := range stmt.Columns {
			if col.VirtualExpr != nil {
				return nil, &ParseError{
					Position: col.Loc.Start, End: col.Loc.End,
					Message: "syntax error: virtual columns are not documented with AS SELECT",
				}
			}
		}
		p.advance() // consume AS
		query, err := p.parseQueryExpr()
		if err != nil {
			return nil, err
		}
		stmt.AsSelect = query
	}
	for _, col := range stmt.Columns {
		if col.VirtualExpr == nil {
			continue
		}
		if stmt.Transient {
			return nil, &ParseError{
				Position: col.Loc.Start, End: col.Loc.End,
				Message: "syntax error: virtual columns are not documented on TRANSIENT ICEBERG TABLE",
			}
		}
		if col.Default != nil || col.Identity != nil || col.Collate != "" || col.ProjectionPolicy != nil {
			return nil, &ParseError{
				Position: col.Loc.Start, End: col.Loc.End,
				Message: "syntax error: an Iceberg virtual column takes no DEFAULT, IDENTITY, COLLATE or PROJECTION POLICY",
			}
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseIcebergTableProperties parses the run of Iceberg table properties, in
// any order, stopping at the first token that begins none of them.
func (p *Parser) parseIcebergTableProperties(stmt *ast.CreateTableStmt) error {
	for {
		switch {
		case p.curIs(kwPARTITION) && p.nextIs(kwBY):
			exprs, err := p.parseIcebergPartitionBy()
			if err != nil {
				return err
			}
			stmt.PartitionBy = exprs
		case p.curIs(kwCLUSTER):
			if err := p.parseClusterByInto(&stmt.ClusterBy, &stmt.Linear); err != nil {
				return err
			}
		case p.curIs(kwCOPY) && p.nextIs(kwGRANTS):
			p.advance() // consume COPY
			p.advance() // consume GRANTS
			stmt.CopyGrants = true
		case p.curIs(kwCOPY) && p.peekIsWord("TAGS"):
			p.advance() // consume COPY
			p.advance() // consume TAGS
			stmt.CopyTags = true
		case p.curIs(kwCOPY) && p.nextIs('='):
			opt, err := p.parseCopyWhatOption()
			if err != nil {
				return err
			}
			stmt.Options = append(stmt.Options, opt)
		case p.curIsWord("PURGE") && !p.nextIs('='):
			p.advance() // consume PURGE
			stmt.Purge = true
		case p.curIs(kwCOMMENT):
			p.advance() // consume COMMENT
			if _, err := p.expect('='); err != nil {
				return err
			}
			tok, err := p.expect(tokString)
			if err != nil {
				return err
			}
			s := tok.Str
			stmt.Comment = &s
		case p.curIs(kwTAG):
			tags, err := p.parseTagAssignments()
			if err != nil {
				return err
			}
			stmt.Tags = append(stmt.Tags, tags...)
		case p.curIs(kwROW) && p.nextIs(kwACCESS):
			rp, err := p.parseRowAccessPolicyClause()
			if err != nil {
				return err
			}
			stmt.RowPolicy = rp
		case p.startsAggregationPolicy():
			p.advance() // consume AGGREGATION
			p.advance() // consume POLICY
			name, err := p.parseObjectName()
			if err != nil {
				return err
			}
			stmt.AggregationPolicy = name
		case p.curIs(kwWITH):
			// WITH introduces TAG, ROW ACCESS POLICY, AGGREGATION POLICY or
			// CONTACT; anything else ends the properties.
			switch {
			case p.nextIs(kwTAG):
				p.advance() // consume WITH
				continue
			case p.nextIs(kwROW):
				p.advance() // consume WITH
				if !p.nextIs(kwACCESS) {
					return p.syntaxErrorAtCur()
				}
				continue
			case p.peekIsWord("AGGREGATION"):
				p.advance() // consume WITH
				if !p.startsAggregationPolicy() {
					return p.syntaxErrorAtCur()
				}
				continue
			case p.peekIsWord("CONTACT"):
				p.advance() // consume WITH
				contacts, err := p.parseWithContact()
				if err != nil {
					return err
				}
				stmt.Contacts = append(stmt.Contacts, contacts...)
			default:
				return nil
			}
		case p.startsCopyOption() && p.nextIs('='):
			// A catalog or storage parameter: EXTERNAL_VOLUME, CATALOG,
			// BASE_LOCATION, CATALOG_TABLE_NAME, CATALOG_NAMESPACE,
			// METADATA_FILE_PATH, TARGET_FILE_SIZE, AUTO_REFRESH, ...
			opt, err := p.parseCopyOption()
			if err != nil {
				return err
			}
			stmt.Options = append(stmt.Options, opt)
		default:
			return nil
		}
	}
}

// startsAggregationPolicy reports whether cur begins AGGREGATION POLICY.
// AGGREGATION is not a reserved keyword, so it is matched by text.
func (p *Parser) startsAggregationPolicy() bool {
	return p.curIsWord("AGGREGATION") && p.nextIs(kwPOLICY)
}

// parseIcebergPartitionBy parses PARTITION BY ( <partition_expr> [, ...] ).
// cur is PARTITION.
func (p *Parser) parseIcebergPartitionBy() ([]ast.Node, error) {
	p.advance() // consume PARTITION
	p.advance() // consume BY
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var exprs []ast.Node
	for {
		e, err := p.parseIcebergPartitionExpr()
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, e)
		if !p.curIs(',') {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return exprs, nil
}

// icebergPartitionTransforms maps each Iceberg partition transform to whether it
// takes a leading numeric argument (BUCKET / TRUNCATE) before the column.
var icebergPartitionTransforms = map[string]bool{
	"BUCKET":   true,
	"TRUNCATE": true,
	"YEAR":     false,
	"MONTH":    false,
	"DAY":      false,
	"HOUR":     false,
}

// parseIcebergPartitionExpr parses one Iceberg partition expression, as the
// CREATE ICEBERG TABLE pages give it:
//
//	<col_name>
//	| BUCKET ( <num_buckets> , <col_name> ) | TRUNCATE ( <width> , <col_name> )
//	| YEAR ( <col_name> ) | MONTH ( <col_name> ) | DAY ( <col_name> ) | HOUR ( <col_name> )
//
// A column is returned as a ColumnRef, a transform as a FuncCallExpr.
func (p *Parser) parseIcebergPartitionExpr() (ast.Node, error) {
	if p.nextIs('(') && p.cur.Str != "" {
		word := strings.ToUpper(p.cur.Str)
		numeric, ok := icebergPartitionTransforms[word]
		if !ok || p.curIs(tokQuotedIdent) {
			return nil, p.syntaxErrorAtCur()
		}
		nameTok := p.advance() // consume the transform name
		p.advance()            // consume '('
		call := &ast.FuncCallExpr{
			Name: ast.ObjectName{Name: ast.Ident{Name: word, Loc: nameTok.Loc}, Loc: nameTok.Loc},
			Loc:  ast.Loc{Start: nameTok.Loc.Start},
		}
		if numeric {
			numTok, err := p.expect(tokInt)
			if err != nil {
				return nil, err
			}
			call.Args = append(call.Args, &ast.Literal{Kind: ast.LitInt, Value: numTok.Str, Ival: numTok.Ival, Loc: numTok.Loc})
			if _, err := p.expect(','); err != nil {
				return nil, err
			}
		}
		col, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		call.Args = append(call.Args, &ast.ColumnRef{Parts: []ast.Ident{col}, Loc: col.Loc})
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		call.Loc.End = p.prev.Loc.End
		return call, nil
	}
	col, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	return &ast.ColumnRef{Parts: []ast.Ident{col}, Loc: col.Loc}, nil
}

// parseCopyWhatOption parses COPY = ( { COMMENTS | TAGS | GRANTS } [, ...] ),
// returned as a CopyOption named COPY whose Group holds one bare entry per
// item. cur is COPY.
func (p *Parser) parseCopyWhatOption() (*ast.CopyOption, error) {
	copyTok := p.advance() // consume COPY
	p.advance()            // consume '='
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	opt := &ast.CopyOption{Name: "COPY", Group: []*ast.CopyOption{}, Loc: ast.Loc{Start: copyTok.Loc.Start}}
	for {
		word := strings.ToUpper(p.cur.Str)
		if p.curIs(tokQuotedIdent) || (word != "COMMENTS" && word != "TAGS" && word != "GRANTS") {
			return nil, p.syntaxErrorAtCur()
		}
		tok := p.advance()
		opt.Group = append(opt.Group, &ast.CopyOption{Name: word, Bare: true, Loc: tok.Loc})
		if !p.curIs(',') {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	opt.Loc.End = p.prev.Loc.End
	return opt, nil
}

// parseWithContact parses CONTACT ( <purpose> = <contact> [, ...] ), after WITH;
// cur is CONTACT. Each pair is a CopyOption whose Words is the contact name.
func (p *Parser) parseWithContact() ([]*ast.CopyOption, error) {
	p.advance() // consume CONTACT
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var contacts []*ast.CopyOption
	for {
		purpose, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		contact, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, &ast.CopyOption{
			Name:  strings.ToUpper(purpose.Name),
			Words: contact.String(),
			Loc:   ast.Loc{Start: purpose.Loc.Start, End: contact.Loc.End},
		})
		if !p.curIs(',') {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return contacts, nil
}
