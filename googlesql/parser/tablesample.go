package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// This file is part of the `parser-query-clauses` DAG node. It implements
// GoogleSQL's TABLESAMPLE table-source operator (GoogleSQLParser.g4 §2.14
// sample_clause / opt_sample_clause_suffix / sample_size / repeatable_clause),
// a hand-port of Google's open-source ZetaSQL reference.
//
// TABLESAMPLE is the OUTERMOST table-source suffix: the grammar models it as
// `table_primary sample_clause`, so it wraps a whole table_primary (which may
// itself already carry a PIVOT/UNPIVOT and an `[AS] alias`). It therefore binds
// AFTER everything else on the source, and takes NO trailing alias of its own
// (oracle: `t TABLESAMPLE BERNOULLI (10 PERCENT) AS x` rejects;
// `t PIVOT(...) TABLESAMPLE BERNOULLI (10 PERCENT)` accepts).
//
// DIALECT NOTE. TABLESAMPLE is a SHARED form — both BigQuery (TABLESAMPLE SYSTEM
// (n PERCENT)) and Spanner (TABLESAMPLE {BERNOULLI|RESERVOIR} (n {PERCENT|ROWS}))
// document it (truth1 QUERY-015 / QUERY-019), so the Spanner emulator oracle is
// authoritative for it (both polarities). The sample method is a bare
// identifier in the grammar, so omni accepts any method name (the union).

// atTableSample reports whether the current token begins a TABLESAMPLE operator.
func (p *Parser) atTableSample() bool {
	return p.cur.Type == kwTABLESAMPLE
}

// parseTableSample parses a TABLESAMPLE operator (sample_clause):
//
//	TABLESAMPLE <method> ( <size> {PERCENT|ROWS} [PARTITION BY <expr>[, …]] )
//	  [ <suffix> ]
//
// TABLESAMPLE is the current token. The method is a bare identifier (SYSTEM /
// BERNOULLI / RESERVOIR — the grammar reads an `identifier`). The suffix is an
// optional REPEATABLE(seed) or WITH WEIGHT [[AS] alias] [REPEATABLE(seed)]
// (opt_sample_clause_suffix). Returns the populated *SampleClause.
func (p *Parser) parseTableSample() (*ast.SampleClause, error) {
	sampleTok := p.advance() // TABLESAMPLE
	sc := &ast.SampleClause{Loc: sampleTok.Loc}

	// Sampling method identifier.
	methodTok, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	method, err := p.identifierText(methodTok)
	if err != nil {
		return nil, err
	}
	sc.Method = method

	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	// sample_size: sample_size_value sample_size_unit [PARTITION BY …].
	// sample_size_value is the STRICT possibly_cast_int_literal_or_parameter |
	// floating_point_literal — NOT a full expression: a bare identifier, an
	// arithmetic expression, or a function call before the unit is a syntax error
	// (oracle: `(1 + 1 ROWS)` / `(foo() ROWS)` / `(x ROWS)` all reject). This is
	// the same sample_size_value the graph TABLESAMPLE uses, so reuse that helper.
	size, err := p.parseGraphSampleSizeValue()
	if err != nil {
		return nil, err
	}
	sc.Size = size

	// Unit: PERCENT | ROWS (required — oracle: `TABLESAMPLE BERNOULLI (10)`
	// rejects with "Expected keyword PERCENT or keyword ROWS").
	switch p.cur.Type {
	case kwPERCENT:
		p.advance()
		sc.Unit = ast.SampleUnitPercent
	case kwROWS:
		p.advance()
		sc.Unit = ast.SampleUnitRows
	default:
		return nil, p.syntaxErrorAtCur()
	}

	// Optional PARTITION BY <expr>[, …] inside the size clause
	// (partition_by_clause_prefix_no_hint).
	if p.cur.Type == kwPARTITION {
		p.advance() // PARTITION
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		for {
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			sc.PartitionBy = append(sc.PartitionBy, expr)
			if p.cur.Type != int(',') {
				break
			}
			p.advance() // ','
		}
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	sc.Loc.End = closeTok.Loc.End

	// Optional suffix (opt_sample_clause_suffix): REPEATABLE(seed) | WITH WEIGHT …
	if err := p.parseSampleSuffix(sc); err != nil {
		return nil, err
	}
	return sc, nil
}

// parseSampleSuffix parses the optional TABLESAMPLE suffix
// (opt_sample_clause_suffix): either a bare REPEATABLE(seed), or
// `WITH WEIGHT [[AS] alias] [REPEATABLE(seed)]`. The current token is just past
// the size-clause ')'. A no-suffix position is a no-op.
func (p *Parser) parseSampleSuffix(sc *ast.SampleClause) error {
	switch p.cur.Type {
	case kwREPEATABLE:
		seed, err := p.parseRepeatableClause()
		if err != nil {
			return err
		}
		sc.Repeatable = seed
		sc.Loc.End = p.prev.Loc.End
	case kwWITH:
		// WITH WEIGHT [identifier|AS identifier] [REPEATABLE(seed)].
		p.advance() // WITH
		if _, err := p.expect(kwWEIGHT); err != nil {
			return err
		}
		sc.WithWeight = true
		sc.Loc.End = p.prev.Loc.End
		// Optional weight-column alias: `WEIGHT alias` or `WEIGHT AS alias`.
		if p.cur.Type == kwAS {
			p.advance()
			aliasTok, err := p.expectIdentifier()
			if err != nil {
				return err
			}
			alias, err := p.identifierText(aliasTok)
			if err != nil {
				return err
			}
			sc.WeightAlias = alias
			sc.Loc.End = p.prev.Loc.End
		} else if isIdentifierStart(p.cur.Type) && p.cur.Type != kwREPEATABLE {
			aliasTok := p.advance()
			alias, err := p.identifierText(aliasTok)
			if err != nil {
				return err
			}
			sc.WeightAlias = alias
			sc.Loc.End = p.prev.Loc.End
		}
		// Optional trailing REPEATABLE(seed).
		if p.cur.Type == kwREPEATABLE {
			seed, err := p.parseRepeatableClause()
			if err != nil {
				return err
			}
			sc.Repeatable = seed
			sc.Loc.End = p.prev.Loc.End
		}
	}
	return nil
}

// parseRepeatableClause parses `REPEATABLE ( <seed> )` (repeatable_clause). The
// seed is a possibly-cast integer literal or parameter; we parse it as a general
// expression (which subsumes CAST(...) and @param). REPEATABLE is the current
// token. Returns the seed expression.
func (p *Parser) parseRepeatableClause() (ast.Node, error) {
	p.advance() // REPEATABLE
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	seed, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return seed, nil
}
