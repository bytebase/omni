package parser

// The Doris fork brought Doris's job scheduler statements with it: CREATE,
// ALTER, DROP, PAUSE and RESUME JOB, CANCEL TASK, SHOW JOB(S), and the
// materialized view forms PAUSE / RESUME MATERIALIZED VIEW JOB and CANCEL
// MATERIALIZED VIEW TASK. StarRocks has none of them. A 3.4.10 container
// rejects every one as a syntax error at the word after the verb ("No viable
// statement for input 'CREATE JOB'", "Unexpected input 'JOB', the most similar
// input is {'ROUTINE'}"), so this parser rejects them there too.
//
// StarRocks's own scheduled work goes through SUBMIT TASK and DROP TASK, a
// materialized view's refresh through ALTER / REFRESH / CANCEL REFRESH
// MATERIALIZED VIEW, and statistics jobs through the ANALYZE statements.

// rejectJobStatement reports a Doris job statement at cur, the word after its
// verb: JOB, JOBS, TASK, or the MATERIALIZED of a materialized view job.
//
// It is a syntax error rather than the dispatch's "not yet supported" stub,
// which says the parser lacks a statement the engine has; here the engine has
// none. SHOW and CANCEL need it explicitly, since their fallbacks take any
// tail and would accept these forms.
func (p *Parser) rejectJobStatement() *ParseError {
	return p.syntaxErrorAtCur()
}
