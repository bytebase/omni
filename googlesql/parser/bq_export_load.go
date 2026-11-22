package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// ---------------------------------------------------------------------------
// BigQuery data-movement statements — EXPORT DATA / EXPORT MODEL / LOAD DATA /
// CLONE DATA (parser-dml-ext node).
// ---------------------------------------------------------------------------
//
// Ports the §2.8 data-movement rules of the legacy ANTLR GoogleSQLParser.g4
// (a hand-port of ZetaSQL):
//
//	export_data_statement:  export_data_no_query as_query
//	export_data_no_query:   EXPORT DATA with_connection_clause? opt_options_list?
//	export_model_statement: EXPORT MODEL path_expression with_connection_clause? opt_options_list?
//	aux_load_data_statement: LOAD DATA append_or_overwrite maybe_dashed_path_expression_with_scope
//	    table_element_list? load_data_partitions_clause? collate_clause?
//	    partition_by_clause_prefix_no_hint? cluster_by_clause_prefix_no_hint? opt_options_list?
//	    aux_load_data_from_files_options_list opt_external_table_with_clauses?
//	clone_data_statement:   CLONE DATA INTO maybe_dashed_path_expression FROM clone_data_source_list
//	clone_data_source_list: clone_data_source (UNION ALL clone_data_source)*
//	clone_data_source:      maybe_dashed_path_expression opt_at_system_time? where_clause?
//	export_metadata_statement: EXPORT table_or_table_function METADATA FROM
//	    maybe_dashed_path_expression opt_with_connection_clause opt_options_list
//
// ORACLE NOTE — all four are BigQuery-only at the GoogleSQL union level, and the
// emulator decides no row of theirs: it answers "Statement not supported:
// ExportDataStmt" / "AuxLoadDataStmt" / "CloneDataStmt" from statement dispatch
// BEFORE it parses a body, so `LOAD DATA INTO FROM FILES()` and a truncated
// `CLONE DATA INTO a FROM` come back as accepts too. The union parser follows
// the grammar (testdata/README.md, "Most of those families are not parse
// verdicts either") and the BigQuery corpus (OTHER-001 EXPORT DATA, OTHER-003
// LOAD DATA).
//
// EXPORT … METADATA (export_metadata_statement) is the third EXPORT form; it is
// parsed here, beside its siblings, because it shares their WITH CONNECTION /
// OPTIONS clauses.

// parseExportStmt dispatches the three EXPORT statements on the keyword that
// follows EXPORT:
//
//	EXPORT DATA …                         → export_data_statement  (this node)
//	EXPORT MODEL …                        → export_model_statement (this node)
//	EXPORT (TABLE|TABLE FUNCTION) METADATA → export_metadata_statement (this node)
//
// cur is at EXPORT.
func (p *Parser) parseExportStmt() (ast.Node, error) {
	switch p.peekNext().Type {
	case kwDATA:
		return p.parseStmtWithSubqueries(p.parseExportData)
	case kwMODEL:
		// EXPORT MODEL's OPTIONS values are full expressions that may embed
		// subqueries; wrap for fillSubqueries like the other statements here.
		return p.parseStmtWithSubqueries(p.parseExportModel)
	default:
		// EXPORT TABLE [FUNCTION] METADATA FROM … (export_metadata_statement). Its
		// OPTIONS values may embed subqueries, so wrap for fillSubqueries like the
		// DATA and MODEL branches.
		return p.parseStmtWithSubqueries(p.parseExportMetadata)
	}
}

// parseExportMetadata parses an export_metadata_statement:
//
//	EXPORT (TABLE | TABLE FUNCTION) METADATA FROM maybe_dashed_path_expression
//	  with_connection_clause? opt_options_list?
//
// cur is at EXPORT. The grammar reports `EXPORT TABLE FUNCTION METADATA` as a
// syntax error (a table-valued function's metadata is not exportable), so a
// FUNCTION keyword here is rejected.
func (p *Parser) parseExportMetadata() (ast.Node, error) {
	export := p.advance() // EXPORT
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	if p.cur.Type == kwFUNCTION {
		// export_metadata_statement's action: "EXPORT TABLE FUNCTION METADATA is not
		// supported". Report it at the keyword, as ZetaSQL does.
		return nil, &ParseError{Position: p.cur.Loc.Start, End: p.cur.Loc.End, Message: "syntax error: EXPORT TABLE FUNCTION METADATA is not supported"}
	}
	if _, err := p.expect(kwMETADATA); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}

	stmt := &ast.ExportMetadataStmt{Loc: export.Loc}
	path, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Path = path
	stmt.Loc.End = path.Loc.End

	// with_connection_clause?
	if p.cur.Type == kwWITH && p.peekNext().Type == kwCONNECTION {
		conn, err := p.parseConnectionClause()
		if err != nil {
			return nil, err
		}
		stmt.HasConnection = true
		stmt.ConnectionName = conn
		stmt.Loc.End = p.prev.Loc.End
	}

	// opt_options_list?
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
		stmt.Loc.End = p.prev.Loc.End
	}
	return stmt, nil
}

// parseExportData parses an export_data_statement:
//
//	EXPORT DATA with_connection_clause? opt_options_list? AS query_after_as
//
// cur is at EXPORT. query_after_as is a query or a GQL query, so `EXPORT DATA AS
// GRAPH g RETURN …` is accepted here.
func (p *Parser) parseExportData() (ast.Node, error) {
	stmt, err := p.parseExportDataNoQuery()
	if err != nil {
		return nil, err
	}

	// as_query — AS query_after_as (required).
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	if p.cur.Type == kwGRAPH {
		gql, err := p.parseGQLStmt()
		if err != nil {
			return nil, err
		}
		stmt.Query = gql
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}
	q, err := p.parseQuery()
	if err != nil {
		return nil, err
	}
	stmt.Query = q

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseExportDataNoQuery parses export_data_no_query:
//
//	EXPORT DATA with_connection_clause? opt_options_list?
//
// cur is at EXPORT. The caller adds the AS query (parseExportData) or, for the
// pipe EXPORT DATA operator, nothing: the exported rows are the pipe input.
func (p *Parser) parseExportDataNoQuery() (*ast.ExportDataStmt, error) {
	export := p.advance() // EXPORT
	if p.isBigQuery() {
		if _, err := p.expect(kwDATA); err != nil {
			return nil, err
		}
	} else {
		p.advance() // DATA
	}

	stmt := &ast.ExportDataStmt{}
	stmt.Loc.Start = export.Loc.Start

	// with_connection_clause? — WITH CONNECTION connection_clause.
	if p.cur.Type == kwWITH && p.peekNext().Type == kwCONNECTION {
		name, err := p.parseConnectionClause()
		if err != nil {
			return nil, err
		}
		stmt.HasConnection = true
		stmt.ConnectionName = name
	}

	// opt_options_list?
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseExportModel parses an export_model_statement:
//
//	EXPORT MODEL path_expression with_connection_clause? opt_options_list?
//
// cur is at EXPORT.
func (p *Parser) parseExportModel() (ast.Node, error) {
	export := p.advance() // EXPORT
	p.advance()           // MODEL

	stmt := &ast.ExportModelStmt{}
	stmt.Loc.Start = export.Loc.Start

	// export_model_statement uses `path_expression` (identifier (DOT identifier)*),
	// NOT maybe_dashed_path_expression — a dashed name like `my-project.ds.m` is
	// outside the rule (must be backtick-quoted). Use parsePathExpr, not
	// parseTablePath (which would fold dashed BigQuery components).
	name, err := p.parsePathExpr()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// with_connection_clause?
	if p.cur.Type == kwWITH && p.peekNext().Type == kwCONNECTION {
		conn, err := p.parseConnectionClause()
		if err != nil {
			return nil, err
		}
		stmt.HasConnection = true
		stmt.ConnectionName = conn
	}

	// opt_options_list?
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseLoadData parses an aux_load_data_statement:
//
//	LOAD DATA (INTO|OVERWRITE) maybe_dashed_path_expression_with_scope
//	  table_element_list? load_data_partitions_clause? collate_clause?
//	  partition_by_clause_prefix_no_hint? cluster_by_clause_prefix_no_hint?
//	  opt_options_list? aux_load_data_from_files_options_list
//	  opt_external_table_with_clauses?
//
// cur is at LOAD. The FROM FILES (options) clause is REQUIRED by the grammar.
func (p *Parser) parseLoadData() (ast.Node, error) {
	load := p.advance() // LOAD
	if _, err := p.expect(kwDATA); err != nil {
		return nil, err
	}

	stmt := &ast.LoadDataStmt{}
	stmt.Loc.Start = load.Loc.Start

	// append_or_overwrite: INTO | OVERWRITE.
	switch p.cur.Type {
	case kwINTO:
		p.advance()
		stmt.Overwrite = false
	case kwOVERWRITE:
		p.advance()
		stmt.Overwrite = true
	default:
		return nil, p.syntaxErrorAtCur()
	}

	// maybe_dashed_path_expression_with_scope:
	//   (TEMP|TEMPORARY) TABLE maybe_dashed_path_expression | maybe_dashed_path_expression
	//
	// TEMP / TEMPORARY are non-reserved (common_keyword_as_identifier), so they are
	// ALSO legal path components: `LOAD DATA INTO temp.t …` targets a table named
	// `temp`. Only consume the scope prefix when TABLE follows; otherwise fall
	// through to parse the leading TEMP/TEMPORARY as the path's first component.
	if (p.cur.Type == kwTEMP || p.cur.Type == kwTEMPORARY) && p.peekNext().Type == kwTABLE {
		scopeTok := p.advance() // TEMP | TEMPORARY
		stmt.Temp = true
		stmt.TempKeyword = TokenName(scopeTok.Type)
		p.advance() // TABLE
	}
	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// table_element_list? — an explicit `( col defs / constraints )`.
	if p.cur.Type == int('(') {
		cols, cons, err := p.parseTableElementList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
		stmt.Constraints = cons
	}

	// load_data_partitions_clause? — OVERWRITE? PARTITIONS ( expr ).
	if p.cur.Type == kwOVERWRITE || p.cur.Type == kwPARTITIONS {
		ow, expr, err := p.parseLoadDataPartitions()
		if err != nil {
			return nil, err
		}
		stmt.PartitionsOverwrite = ow
		stmt.Partitions = expr
	}

	// collate_clause? — COLLATE string_literal_or_parameter.
	if p.cur.Type == kwCOLLATE {
		coll, _, err := p.parseCollateClause()
		if err != nil {
			return nil, err
		}
		stmt.Collate = coll
	}

	// partition_by_clause_prefix_no_hint? — PARTITION BY expr (, expr)*.
	if p.cur.Type == kwPARTITION {
		exprs, err := p.parsePartitionByNoHint()
		if err != nil {
			return nil, err
		}
		stmt.PartitionBy = exprs
	}

	// cluster_by_clause_prefix_no_hint? — CLUSTER BY expr (, expr)*.
	if p.cur.Type == kwCLUSTER {
		exprs, err := p.parseClusterByNoHint()
		if err != nil {
			return nil, err
		}
		stmt.ClusterBy = exprs
	}

	// opt_options_list?
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	// aux_load_data_from_files_options_list — FROM FILES options_list (REQUIRED).
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFILES); err != nil {
		return nil, err
	}
	if p.cur.Type != int('(') {
		return nil, p.syntaxErrorAtCur()
	}
	// options_list is `( options_entry (, options_entry)* )`; parseOptionsList
	// expects cur at the OPTIONS keyword, so parse the bare parenthesized body.
	files, err := p.parseLoadDataFromFilesOptions()
	if err != nil {
		return nil, err
	}
	stmt.FromFiles = files

	// opt_external_table_with_clauses?:
	//   with_partition_columns_clause with_connection_clause
	//   | with_partition_columns_clause
	//   | with_connection_clause
	if p.cur.Type == kwWITH {
		switch p.peekNext().Type {
		case kwPARTITION:
			cols, err := p.parseWithPartitionColumns()
			if err != nil {
				return nil, err
			}
			stmt.HasPartitionColumns = true
			stmt.PartitionColumns = cols
			// optionally followed by WITH CONNECTION.
			if p.cur.Type == kwWITH && p.peekNext().Type == kwCONNECTION {
				conn, err := p.parseConnectionClause()
				if err != nil {
					return nil, err
				}
				stmt.HasConnection = true
				stmt.ConnectionName = conn
			}
		case kwCONNECTION:
			conn, err := p.parseConnectionClause()
			if err != nil {
				return nil, err
			}
			stmt.HasConnection = true
			stmt.ConnectionName = conn
		default:
			// A leading WITH followed by neither PARTITION nor CONNECTION is a
			// syntax error (no other external-table WITH clause exists). cur is the
			// WITH token; report there.
			return nil, p.syntaxErrorAtCur()
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseLoadDataPartitions parses a load_data_partitions_clause:
//
//	OVERWRITE? PARTITIONS '(' expression ')'
//
// cur is at OVERWRITE or PARTITIONS. Returns whether OVERWRITE was present and
// the partition expression.
func (p *Parser) parseLoadDataPartitions() (bool, ast.Node, error) {
	overwrite := false
	if p.cur.Type == kwOVERWRITE {
		p.advance() // OVERWRITE
		overwrite = true
	}
	if _, err := p.expect(kwPARTITIONS); err != nil {
		return false, nil, err
	}
	if _, err := p.expect(int('(')); err != nil {
		return false, nil, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return false, nil, err
	}
	if _, err := p.expect(int(')')); err != nil {
		return false, nil, err
	}
	return overwrite, expr, nil
}

// parseLoadDataFromFilesOptions parses the `options_list` of an
// aux_load_data_from_files_options_list (the `FROM FILES` keywords are already
// consumed). It is a bare `( options_entry (, options_entry)* )` — the same body
// parseOptionsList parses after the OPTIONS keyword. cur is at '('.
func (p *Parser) parseLoadDataFromFilesOptions() ([]*ast.OptionsEntry, error) {
	p.advance() // '('
	var entries []*ast.OptionsEntry
	if p.cur.Type == int(')') {
		p.advance()
		return entries, nil
	}
	for {
		entry, err := p.parseOptionsEntry()
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
		if _, ok := p.match(int(',')); ok {
			continue
		}
		break
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return entries, nil
}

// parseWithPartitionColumns parses a with_partition_columns_clause:
//
//	WITH PARTITION COLUMNS table_element_list?
//
// cur is at WITH (with PARTITION next). Returns the optional explicit column
// list (nil when the columns are inferred, i.e. no `( … )` follows).
func (p *Parser) parseWithPartitionColumns() ([]*ast.ColumnDef, error) {
	p.advance() // WITH
	if _, err := p.expect(kwPARTITION); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwCOLUMNS); err != nil {
		return nil, err
	}
	if p.cur.Type != int('(') {
		// No explicit list — partition columns are inferred from the files.
		return nil, nil
	}
	cols, _, err := p.parseTableElementList()
	if err != nil {
		return nil, err
	}
	return cols, nil
}

// parseCloneData parses a clone_data_statement:
//
//	CLONE DATA INTO maybe_dashed_path_expression FROM clone_data_source_list
//	clone_data_source_list: clone_data_source (UNION ALL clone_data_source)*
//	clone_data_source:      maybe_dashed_path_expression opt_at_system_time? where_clause?
//
// cur is at CLONE.
func (p *Parser) parseCloneData() (ast.Node, error) {
	clone := p.advance() // CLONE
	if _, err := p.expect(kwDATA); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwINTO); err != nil {
		return nil, err
	}

	stmt := &ast.CloneDataStmt{}
	stmt.Loc.Start = clone.Loc.Start

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}

	// clone_data_source_list: at least one source, then (UNION ALL source)*.
	for {
		src, err := p.parseCloneDataSource()
		if err != nil {
			return nil, err
		}
		stmt.Sources = append(stmt.Sources, src)
		if p.cur.Type == kwUNION && p.peekNext().Type == kwALL {
			p.advance() // UNION
			p.advance() // ALL
			continue
		}
		break
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseCloneDataSource parses one clone_data_source:
//
//	maybe_dashed_path_expression opt_at_system_time? where_clause?
//
// cur is at the source path.
func (p *Parser) parseCloneDataSource() (*ast.CloneDataSource, error) {
	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	src := &ast.CloneDataSource{Name: name}
	src.Loc = name.Loc

	// opt_at_system_time? — FOR SYSTEM_TIME AS OF expr.
	if p.cur.Type == kwFOR {
		ts, err := p.parseForSystemTimeExpr()
		if err != nil {
			return nil, err
		}
		src.ForSystemTime = ts
	}

	// where_clause? — WHERE expr.
	if p.cur.Type == kwWHERE {
		p.advance() // WHERE
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		src.Where = expr
	}

	src.Loc.End = p.prev.Loc.End
	return src, nil
}
