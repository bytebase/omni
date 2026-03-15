// Package parser - create_index.go implements T-SQL CREATE INDEX statement parsing.
package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateIndexStmt parses a CREATE INDEX statement.
//
// BNF: mssql/parser/bnf/create-index-transact-sql.bnf
//
//	CREATE [ UNIQUE ] [ CLUSTERED | NONCLUSTERED ] INDEX index_name
//	    ON <object> ( column [ ASC | DESC ] [ ,...n ] )
//	    [ INCLUDE ( column_name [ ,...n ] ) ]
//	    [ WHERE <filter_predicate> ]
//	    [ WITH ( <relational_index_option> [ ,...n ] ) ]
//	    [ ON { partition_scheme_name ( column_name )
//	         | filegroup_name
//	         | default
//	         }
//	    ]
//	    [ FILESTREAM_ON { filestream_filegroup_name | partition_scheme_name | "NULL" } ]
//
//	CREATE [ UNIQUE ] [ CLUSTERED | NONCLUSTERED ] COLUMNSTORE INDEX index_name
//	    ON <object>
//	    [ ( column [ ,...n ] ) ]
//	    [ ORDER ( column [ ,...n ] ) ]
//	    [ WHERE <filter_expression> ]
//	    [ WITH ( <with_option> [ ,...n ] ) ]
//	    [ ON { partition_scheme_name ( column_name )
//	         | filegroup_name
//	         | default
//	         }
//	    ]
//	    [ FILESTREAM_ON { filestream_filegroup_name | partition_scheme_name | "NULL" } ]
func (p *Parser) parseCreateIndexStmt(unique bool) *nodes.CreateIndexStmt {
	loc := p.pos()

	stmt := &nodes.CreateIndexStmt{
		Unique: unique,
		Loc:    nodes.Loc{Start: loc},
	}

	// CLUSTERED / NONCLUSTERED
	if p.cur.Type == kwCLUSTERED {
		p.advance()
		v := true
		stmt.Clustered = &v
	} else if p.cur.Type == kwNONCLUSTERED {
		p.advance()
		v := false
		stmt.Clustered = &v
	}

	// COLUMNSTORE (optional)
	if p.cur.Type == kwCOLUMNSTORE {
		p.advance()
		stmt.Columnstore = true
	}

	// INDEX keyword
	p.match(kwINDEX)

	// Index name
	name, _ := p.parseIdentifier()
	stmt.Name = name

	// ON table
	if _, ok := p.match(kwON); ok {
		stmt.Table = p.parseTableRef()
	}

	// Column list
	if p.cur.Type == '(' {
		stmt.Columns = p.parseIndexColumnList()
	}

	// INCLUDE (relational index only)
	if p.cur.Type == kwINCLUDE {
		p.advance()
		if p.cur.Type == '(' {
			stmt.IncludeCols = p.parseParenIdentList()
		}
	}

	// ORDER (columnstore index only)
	if p.cur.Type == kwORDER {
		p.advance()
		if p.cur.Type == '(' {
			stmt.OrderCols = p.parseParenIdentList()
		}
	}

	// WHERE (filtered index)
	if _, ok := p.match(kwWHERE); ok {
		stmt.WhereClause = p.parseExpr()
	}

	// WITH ( <relational_index_option> [ ,...n ] )
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == '(' {
			stmt.Options = p.parseAlterIndexOptions()
		}
	}

	// ON { partition_scheme_name ( column_name ) | filegroup_name | default }
	if _, ok := p.match(kwON); ok {
		if p.isIdentLike() || p.cur.Type == kwDEFAULT {
			stmt.OnFileGroup = p.cur.Str
			p.advance()
			// partition_scheme_name ( column_name )
			if p.cur.Type == '(' {
				p.advance()
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					p.advance()
				}
				p.match(')')
			}
		}
	}

	// FILESTREAM_ON { filestream_filegroup_name | partition_scheme_name | "NULL" }
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FILESTREAM_ON") {
		p.advance()
		if p.isIdentLike() || p.cur.Type == tokSCONST {
			stmt.FilestreamOn = p.cur.Str
			p.advance()
			// partition_scheme_name ( column_name )
			if p.cur.Type == '(' {
				p.advance()
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					p.advance()
				}
				p.match(')')
			}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseIndexColumnList parses (col [ASC|DESC], ...).
func (p *Parser) parseIndexColumnList() *nodes.List {
	p.advance() // consume (
	var items []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		loc := p.pos()
		name, ok := p.parseIdentifier()
		if !ok {
			break
		}
		dir := nodes.SortDefault
		if _, ok := p.match(kwASC); ok {
			dir = nodes.SortAsc
		} else if _, ok := p.match(kwDESC); ok {
			dir = nodes.SortDesc
		}
		items = append(items, &nodes.IndexColumn{
			Name:    name,
			SortDir: dir,
			Loc:     nodes.Loc{Start: loc},
		})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	_, _ = p.expect(')')
	return &nodes.List{Items: items}
}

// parseCreateXmlIndexStmt parses CREATE [PRIMARY] XML INDEX.
// Caller has consumed CREATE [PRIMARY] XML INDEX.
//
// BNF: mssql/parser/bnf/create-xml-index-transact-sql.bnf
//
//	CREATE [ PRIMARY ] XML INDEX index_name
//	    ON <object> ( xml_column_name )
//	    [ USING XML INDEX xml_index_name
//	        [ FOR { VALUE | PATH | PROPERTY } ]
//	    ]
//	    [ WITH ( <xml_index_option> [ ,...n ] ) ]
func (p *Parser) parseCreateXmlIndexStmt(primary bool) *nodes.CreateXmlIndexStmt {
	loc := p.pos()
	stmt := &nodes.CreateXmlIndexStmt{
		Primary: primary,
		Loc:     nodes.Loc{Start: loc},
	}

	// Index name
	name, _ := p.parseIdentifier()
	stmt.Name = name

	// ON table
	if _, ok := p.match(kwON); ok {
		stmt.Table = p.parseTableRef()
	}

	// (xml_column)
	if p.cur.Type == '(' {
		p.advance()
		col, _ := p.parseIdentifier()
		stmt.XmlColumn = col
		p.match(')')
	}

	// USING XML INDEX parent_index_name (secondary only)
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "USING") {
		p.advance() // consume USING
		p.match(kwXML)
		p.match(kwINDEX)
		idx, _ := p.parseIdentifier()
		stmt.UsingIndex = idx

		// FOR VALUE|PATH|PROPERTY
		if _, ok := p.match(kwFOR); ok {
			if p.isIdentLike() {
				stmt.SecondaryFor = p.cur.Str
				p.advance()
			}
		}
	}

	// WITH (options)
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == '(' {
			stmt.Options = p.parseOptionList()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateSelectiveXmlIndexStmt parses CREATE SELECTIVE XML INDEX.
// Caller has consumed CREATE SELECTIVE XML INDEX.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-selective-xml-index-transact-sql
func (p *Parser) parseCreateSelectiveXmlIndexStmt() *nodes.CreateSelectiveXmlIndexStmt {
	loc := p.pos()
	stmt := &nodes.CreateSelectiveXmlIndexStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Index name
	name, _ := p.parseIdentifier()
	stmt.Name = name

	// ON table
	if _, ok := p.match(kwON); ok {
		stmt.Table = p.parseTableRef()
	}

	// (xml_column)
	if p.cur.Type == '(' {
		p.advance()
		col, _ := p.parseIdentifier()
		stmt.XmlColumn = col
		p.match(')')
	}

	// WITH XMLNAMESPACES (...)
	if p.cur.Type == kwWITH {
		p.advance()
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "XMLNAMESPACES") {
			p.advance()
			if p.cur.Type == '(' {
				stmt.Namespaces = p.parseOptionList()
			}
		} else if p.cur.Type == '(' {
			// This is WITH (options), handle after FOR
			// Back up: we already consumed WITH, so parse options now
			stmt.Options = p.parseOptionList()
		}
	}

	// FOR (path_list)
	if _, ok := p.match(kwFOR); ok {
		if p.cur.Type == '(' {
			stmt.Paths = p.parseOptionList()
		}
	}

	// WITH (options) - if not already parsed
	if stmt.Options == nil && p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == '(' {
			stmt.Options = p.parseOptionList()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateSpatialIndexStmt parses CREATE SPATIAL INDEX.
// Caller has consumed CREATE SPATIAL INDEX.
//
// BNF: mssql/parser/bnf/create-spatial-index-transact-sql.bnf
//
//	CREATE SPATIAL INDEX index_name
//	    ON <object> ( spatial_column_name )
//	    { <geometry_tessellation> | <geography_tessellation> }
//	    [ ON { filegroup_name | "default" } ]
//
//	<geometry_tessellation> ::=
//	    [ USING { GEOMETRY_AUTO_GRID | GEOMETRY_GRID } ]
//	    WITH ( <bounding_box> [, <tessellation_grid>] [, <tessellation_cells_per_object>]
//	           [, <spatial_index_option> [,...n] ] )
//
//	<geography_tessellation> ::=
//	    [ USING { GEOGRAPHY_AUTO_GRID | GEOGRAPHY_GRID } ]
//	    [ WITH ( [<tessellation_grid>] [, <tessellation_cells_per_object>]
//	             [, <spatial_index_option> [,...n] ] ) ]
func (p *Parser) parseCreateSpatialIndexStmt() *nodes.CreateSpatialIndexStmt {
	loc := p.pos()
	stmt := &nodes.CreateSpatialIndexStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Index name
	name, _ := p.parseIdentifier()
	stmt.Name = name

	// ON table
	if _, ok := p.match(kwON); ok {
		stmt.Table = p.parseTableRef()
	}

	// (spatial_column)
	if p.cur.Type == '(' {
		p.advance()
		col, _ := p.parseIdentifier()
		stmt.SpatialColumn = col
		p.match(')')
	}

	// USING tessellation_scheme
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "USING") {
		p.advance() // consume USING
		if p.isIdentLike() {
			stmt.Using = p.cur.Str
			p.advance()
		}
	}

	// WITH (options) - spatial index uses the same option key=value format
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == '(' {
			stmt.Options = p.parseAlterIndexOptions()
		}
	}

	// ON filegroup
	if _, ok := p.match(kwON); ok {
		if p.isIdentLike() {
			stmt.OnFileGroup = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateAggregateStmt parses CREATE AGGREGATE.
// Caller has consumed CREATE AGGREGATE.
//
// BNF: mssql/parser/bnf/create-aggregate-transact-sql.bnf
//
//	CREATE AGGREGATE [ schema_name . ] aggregate_name
//	        (@param_name <input_sqltype>
//	        [ ,...n ] )
//	RETURNS <return_sqltype>
//	EXTERNAL NAME assembly_name [ .class_name ]
//
//	<input_sqltype> ::=
//	        system_scalar_type | { [ udt_schema_name. ] udt_type_name }
//
//	<return_sqltype> ::=
//	        system_scalar_type | { [ udt_schema_name. ] udt_type_name }
func (p *Parser) parseCreateAggregateStmt() *nodes.CreateAggregateStmt {
	loc := p.pos()
	stmt := &nodes.CreateAggregateStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// [schema.]aggregate_name
	stmt.Name = p.parseTableRef()

	// (@param_name type, ...)
	if p.cur.Type == '(' {
		p.advance()
		var params []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			param := &nodes.ParamDef{Loc: nodes.Loc{Start: p.pos()}}
			// @param_name
			if p.cur.Type == tokVARIABLE {
				param.Name = p.cur.Str
				p.advance()
			}
			// data type
			param.DataType = p.parseDataType()
			params = append(params, param)
			if _, ok := p.match(','); !ok {
				break
			}
		}
		p.match(')')
		stmt.Params = &nodes.List{Items: params}
	}

	// RETURNS return_type
	if p.cur.Type == kwRETURNS {
		p.advance()
		stmt.ReturnType = p.parseDataType()
	}

	// EXTERNAL NAME assembly_qualified_name
	if p.cur.Type == kwEXTERNAL {
		p.advance()
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "NAME") {
			p.advance()
		}
		// Read dotted name: assembly.class[.method]
		var parts []string
		if p.isIdentLike() {
			parts = append(parts, p.cur.Str)
			p.advance()
		}
		for {
			if _, ok := p.match('.'); ok {
				if p.isIdentLike() {
					parts = append(parts, p.cur.Str)
					p.advance()
				}
			} else {
				break
			}
		}
		stmt.ExternalName = joinDots(parts)
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseDropAggregateStmt parses DROP AGGREGATE [IF EXISTS] name.
// Caller has consumed DROP AGGREGATE.
//
// BNF: mssql/parser/bnf/drop-aggregate-transact-sql.bnf
//
//	DROP AGGREGATE [ IF EXISTS ] [ schema_name . ] aggregate_name
func (p *Parser) parseDropAggregateStmt() *nodes.DropAggregateStmt {
	loc := p.pos()
	stmt := &nodes.DropAggregateStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// IF EXISTS
	if _, ok := p.match(kwIF); ok {
		p.match(kwEXISTS)
		stmt.IfExists = true
	}

	stmt.Name = p.parseTableRef()
	stmt.Loc.End = p.pos()
	return stmt
}

// joinDots joins string parts with dots.
func joinDots(parts []string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += "."
		}
		result += p
	}
	return result
}

// parseCreateJsonIndexStmt parses a CREATE JSON INDEX statement.
// Caller has consumed CREATE JSON INDEX.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-json-index-transact-sql
//
//	CREATE JSON INDEX name ON table_name (json_column_name)
//	  [ FOR ( sql_json_path [ , ...n ] ) ]
//	  [ WITH ( <json_index_option> [ , ...n ] ) ]
//	  [ ON { filegroup_name | "default" } ]
//
//	<json_index_option> ::=
//	{
//	    OPTIMIZE_FOR_ARRAY_SEARCH = { ON | OFF }
//	  | FILLFACTOR = fillfactor
//	  | DROP_EXISTING = { ON | OFF }
//	  | ONLINE = OFF
//	  | ALLOW_ROW_LOCKS = { ON | OFF }
//	  | ALLOW_PAGE_LOCKS = { ON | OFF }
//	  | MAXDOP = max_degree_of_parallelism
//	  | DATA_COMPRESSION = { NONE | ROW | PAGE }
//	}
func (p *Parser) parseCreateJsonIndexStmt() *nodes.CreateJsonIndexStmt {
	loc := p.pos()
	stmt := &nodes.CreateJsonIndexStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Index name
	name, _ := p.parseIdentifier()
	stmt.Name = name

	// ON table
	if _, ok := p.match(kwON); ok {
		stmt.Table = p.parseTableRef()
	}

	// (json_column_name)
	if p.cur.Type == '(' {
		p.advance()
		col, _ := p.parseIdentifier()
		stmt.JsonColumn = col
		p.match(')')
	}

	// FOR (sql_json_path, ...)
	if _, ok := p.match(kwFOR); ok {
		if p.cur.Type == '(' {
			p.advance()
			var paths []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				// JSON paths are string literals like '$.name'
				expr := p.parseExpr()
				if expr != nil {
					paths = append(paths, expr)
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
			p.match(')')
			if len(paths) > 0 {
				stmt.ForPaths = &nodes.List{Items: paths}
			}
		}
	}

	// WITH (options)
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == '(' {
			stmt.Options = p.parseOptionList()
		}
	}

	// ON filegroup
	if _, ok := p.match(kwON); ok {
		if p.isIdentLike() {
			stmt.OnFileGroup = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateVectorIndexStmt parses a CREATE VECTOR INDEX statement.
// Caller has consumed CREATE VECTOR INDEX.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-vector-index-transact-sql
//
//	CREATE VECTOR INDEX index_name
//	ON object ( vector_column )
//	[ WITH (
//	    [ , ] METRIC = { 'cosine' | 'dot' | 'euclidean' }
//	    [ [ , ] TYPE = 'DiskANN' ]
//	    [ [ , ] MAXDOP = max_degree_of_parallelism ]
//	) ]
//	[ ON { filegroup_name | "default" } ]
func (p *Parser) parseCreateVectorIndexStmt() *nodes.CreateVectorIndexStmt {
	loc := p.pos()
	stmt := &nodes.CreateVectorIndexStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Index name
	name, _ := p.parseIdentifier()
	stmt.Name = name

	// ON table
	if _, ok := p.match(kwON); ok {
		stmt.Table = p.parseTableRef()
	}

	// (vector_column)
	if p.cur.Type == '(' {
		p.advance()
		col, _ := p.parseIdentifier()
		stmt.VectorCol = col
		p.match(')')
	}

	// WITH (options)
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == '(' {
			stmt.Options = p.parseOptionList()
		}
	}

	// ON filegroup
	if _, ok := p.match(kwON); ok {
		if p.isIdentLike() {
			stmt.OnFileGroup = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseOptionList parses (option = value, ...) used in WITH clauses.
func (p *Parser) parseOptionList() *nodes.List {
	p.advance() // consume (
	var items []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		expr := p.parseExpr()
		if expr != nil {
			items = append(items, expr)
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	_, _ = p.expect(')')
	return &nodes.List{Items: items}
}
