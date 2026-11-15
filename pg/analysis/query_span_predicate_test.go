package analysis

import (
	"context"
	"slices"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	omnipg "github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
	"github.com/stretchr/testify/require"
)

func predicateColumnNames(columns ColumnSet) []string {
	names := []string{}
	for column := range columns {
		if column.Column != "" {
			names = append(names, column.Table+"."+column.Column)
		}
	}
	slices.Sort(names)
	return names
}

func TestGetQuerySpanPredicateColumns(t *testing.T) {
	testCases := []struct {
		name      string
		statement string
		want      []string
	}{
		{"where", "SELECT id FROM t WHERE ssn = '1'", []string{"t.ssn"}},
		{"where function", "SELECT id FROM t WHERE ascii(substr(ssn, 1, 1)) > 70", []string{"t.ssn"}},
		{"row comparison", "SELECT id FROM t WHERE (ssn, id) > ('1', 0)", []string{"t.id", "t.ssn"}},
		{"join on", "SELECT o.amt FROM o JOIN t ON t.id = o.id AND t.ssn = '1'", []string{"o.id", "t.id", "t.ssn"}},
		{"left join on", "SELECT o.amt FROM o LEFT JOIN t ON t.id = o.id AND t.ssn = '1'", []string{"o.id", "t.id", "t.ssn"}},
		{"join using", "SELECT a.id FROM t a JOIN t b USING (ssn)", []string{"t.ssn"}},
		{"having", "SELECT id FROM t GROUP BY id HAVING max(ssn) > '1'", []string{"t.ssn"}},
		{"in subquery", "SELECT amt FROM o WHERE id IN (SELECT id FROM t WHERE ssn = '1')", []string{"o.id", "t.id", "t.ssn"}},
		{"subquery output", "SELECT amt FROM o WHERE '1' = (SELECT max(ssn) FROM t)", []string{"t.ssn"}},
		{"correlated exists", "SELECT amt FROM o WHERE EXISTS (SELECT 1 FROM t WHERE t.id = o.id AND t.ssn = '1')", []string{"o.id", "t.id", "t.ssn"}},
		{"subquery in select list", "SELECT (SELECT count(*) FROM t WHERE ssn = '1') AS c", []string{"t.ssn"}},
		{"derived table", "SELECT id FROM (SELECT * FROM t WHERE ssn = '1') s", []string{"t.ssn"}},
		{"derived table output", "SELECT s.id FROM (SELECT id, ssn AS x FROM t) s WHERE s.x = '1'", []string{"t.ssn"}},
		{"lateral subquery", "SELECT o.amt FROM o, LATERAL (SELECT 1 FROM t WHERE t.id = o.id AND t.ssn = '1') l", []string{"o.id", "t.id", "t.ssn"}},
		{"cte", "WITH c AS (SELECT id FROM t WHERE ssn = '1') SELECT id FROM c", []string{"t.ssn"}},
		{"cte output", "WITH c AS (SELECT id, ssn FROM t) SELECT id FROM c WHERE ssn = '1'", []string{"t.ssn"}},
		{"set operation arm", "SELECT id FROM t WHERE ssn = '1' UNION SELECT id FROM o", []string{"t.ssn"}},
		{"except filters by the other branch", "SELECT 'guess' AS x EXCEPT SELECT ssn FROM t", []string{"t.ssn"}},
		{"intersect filters by the other branch", "SELECT 'guess' AS x INTERSECT ALL SELECT ssn FROM t", []string{"t.ssn"}},
		{"union filters nothing", "SELECT 'guess' AS x UNION SELECT ssn FROM t", []string{}},
		{"explain analyze", "EXPLAIN ANALYZE SELECT id FROM t WHERE ssn = '1'", []string{"t.ssn"}},
		{"aggregate filter", "SELECT sum(id) FILTER (WHERE ssn = '1') AS x FROM t", []string{"t.ssn"}},
		{"natural join", "SELECT 1 FROM t NATURAL JOIN (SELECT '1'::text AS ssn) x", []string{"t.ssn"}},
		{"full join using", "SELECT id FROM t FULL JOIN o USING (id) WHERE id > 0", []string{"o.id", "t.id"}},
		{"limit subquery", "SELECT id FROM t LIMIT (SELECT count(*) FROM t WHERE ssn = '1')", []string{"t.ssn"}},
		{"limit value", "SELECT id FROM t LIMIT (SELECT CASE WHEN ssn = 'guess' THEN 1 ELSE 0 END FROM t LIMIT 1)", []string{"t.ssn"}},
		{"set operation limit", "SELECT 1 AS x UNION ALL SELECT 2 LIMIT (SELECT CASE WHEN ssn = 'guess' THEN 1 ELSE 0 END FROM t LIMIT 1)", []string{"t.ssn"}},
		{"set-returning function in from", "SELECT count(*) AS c FROM t CROSS JOIN LATERAL generate_series(1, ascii(substr(t.ssn, 1, 1)))", []string{"t.ssn"}},
		{"set-returning function in select list", "SELECT generate_series(1, ascii(substr(ssn, 1, 1))) AS g FROM t", []string{"t.ssn"}},
		{"scalar function in select list", "SELECT upper(ssn) AS u FROM t", []string{}},
		{"unknown function", "SELECT id FROM t WHERE levenshtein(ssn, 'a') > 0", []string{"t.ssn"}},
		{"unknown function in from", "SELECT t.id FROM t, LATERAL unknown_rows(t.ssn) f WHERE f.a = 1", []string{"t.ssn"}},
		{"select list only", "SELECT ssn FROM t", []string{}},
		{"condition in select list", "SELECT CASE WHEN ssn = '1' THEN 1 END AS x FROM t", []string{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			span := spanFor(t, tc.statement, healthySchema())
			require.Equal(t, tc.want, predicateColumnNames(span.PredicateColumns))
		})
	}
}

func TestGetQuerySpanPredicateColumnsOfFunctionBodies(t *testing.T) {
	schema := healthySchema()
	// * over an unknown function fails analysis of the body that uses it.
	unanalyzable := "(SELECT * FROM o CROSS JOIN unknown_rows() u) o JOIN t ON t.id = o.id AND t.ssn = guess"
	schema.Schemas[0].Functions = []*metadatapb.FunctionMetadata{
		{Name: "probe", Definition: "CREATE FUNCTION public.probe(guess text) RETURNS TABLE(id int) LANGUAGE plpgsql AS $$ BEGIN RETURN QUERY SELECT t.id FROM t WHERE t.ssn = guess; END; $$"},
		{Name: "probe_sql", Definition: "CREATE FUNCTION public.probe_sql(guess text) RETURNS SETOF int LANGUAGE sql AS $$ SELECT id FROM t WHERE ssn = guess $$"},
		{Name: "probe_fixed", Definition: "CREATE FUNCTION public.probe_fixed() RETURNS SETOF int LANGUAGE sql AS $$ SELECT id FROM t WHERE ssn = 'guess' $$"},
		{Name: "probe_multi", Definition: "CREATE FUNCTION public.probe_multi() RETURNS int LANGUAGE sql AS $$ SELECT pg_sleep(0) FROM t WHERE ssn = 'guess'; SELECT 1 $$"},
		{Name: "probe_perform", Definition: "CREATE FUNCTION public.probe_perform() RETURNS int LANGUAGE plpgsql AS $$ BEGIN PERFORM pg_sleep(0) FROM t WHERE ssn = 'guess'; RETURN 1; END; $$"},
		{Name: "probe_if", Definition: "CREATE FUNCTION public.probe_if() RETURNS int LANGUAGE plpgsql AS $$ BEGIN IF (SELECT ssn FROM t LIMIT 1) = 'guess' THEN RETURN 1; END IF; RETURN 0; END; $$"},
		{Name: "probe_variable", Definition: "CREATE FUNCTION public.probe_variable() RETURNS int LANGUAGE plpgsql AS $$ DECLARE v text; BEGIN SELECT ssn INTO v FROM t WHERE id = 1; IF v = 'guess' THEN RETURN 1; END IF; RETURN 0; END; $$"},
		{Name: "probe_update", Definition: "CREATE FUNCTION public.probe_update() RETURNS int LANGUAGE sql AS $$ UPDATE t SET id = id WHERE ssn = 'guess'; SELECT 1 $$"},
		{Name: "probe_exec", Definition: "CREATE FUNCTION public.probe_exec() RETURNS int LANGUAGE plpgsql AS $$ BEGIN DELETE FROM t WHERE ssn = 'guess'; RETURN 1; END; $$"},
		{Name: "probe_table", Definition: "CREATE FUNCTION public.probe_table() RETURNS TABLE(id int) LANGUAGE sql AS $$ SELECT id FROM t WHERE ssn = 'guess' $$"},
		{Name: "probe_record", Definition: "CREATE FUNCTION public.probe_record() RETURNS int LANGUAGE plpgsql AS $$ DECLARE r record; BEGIN SELECT ssn INTO r FROM t LIMIT 1; IF r.ssn = 'guess' THEN RETURN 1; END IF; RETURN 0; END; $$"},
		{Name: "probe_execute", Definition: "CREATE FUNCTION public.probe_execute() RETURNS int LANGUAGE plpgsql AS $$ BEGIN EXECUTE 'SELECT pg_sleep(0) FROM t ' || 'WHERE ssn = ''guess'''; RETURN 1; END; $$"},
		{Name: "probe_count", Definition: "CREATE FUNCTION public.probe_count() RETURNS int LANGUAGE sql AS $$ SELECT count(*)::int FROM t WHERE ssn = 'guess' $$"},
		{Name: "probe_bool", Definition: "CREATE FUNCTION public.probe_bool() RETURNS boolean LANGUAGE sql AS $$ SELECT EXISTS (SELECT 1 FROM t WHERE ssn = 'guess') $$"},
		{Name: "probe_raise", Definition: "CREATE FUNCTION public.probe_raise() RETURNS int LANGUAGE plpgsql AS $$ DECLARE v text; BEGIN SELECT ssn INTO v FROM t LIMIT 1; RAISE EXCEPTION '%', v; END; $$"},
		// The overload that reads ssn loads first, so a body chosen by name alone
		// would be the harmless one.
		{Name: "probe_over", Definition: "CREATE FUNCTION public.probe_over(x text) RETURNS boolean LANGUAGE sql AS $$ SELECT EXISTS (SELECT 1 FROM t WHERE ssn = x) $$"},
		{Name: "probe_over", Definition: "CREATE FUNCTION public.probe_over(x int) RETURNS boolean LANGUAGE sql AS $$ SELECT true $$"},
		// A missing return type makes the catalog install these as stand-ins.
		{Name: "probe_stand_in", Definition: "CREATE FUNCTION public.probe_stand_in(x text) RETURNS SETOF missing_tbl LANGUAGE sql AS $$ SELECT * FROM t WHERE ssn = x $$"},
		{Name: "probe_stand_in", Definition: "CREATE FUNCTION public.probe_stand_in(x int) RETURNS SETOF missing_tbl LANGUAGE sql AS $$ SELECT * FROM o $$"},
		{Name: "probe_stand_in_pl", Definition: "CREATE FUNCTION public.probe_stand_in_pl(x text) RETURNS SETOF missing_tbl LANGUAGE plpgsql AS $$ BEGIN RETURN QUERY SELECT * FROM t WHERE ssn = x; END; $$"},
		{Name: "probe_discard", Definition: "CREATE FUNCTION public.probe_discard() RETURNS int LANGUAGE sql AS $$ SELECT pg_sleep(CASE WHEN ssn = 'guess' THEN 0 ELSE 0 END) FROM t; SELECT 1 $$"},
		{Name: "probe_declare", Definition: "CREATE FUNCTION public.probe_declare() RETURNS int LANGUAGE plpgsql AS $$ DECLARE v text := (SELECT ssn FROM t LIMIT 1); BEGIN RAISE EXCEPTION '%', v; END; $$"},
		{Name: "proc_check", Definition: "CREATE PROCEDURE public.proc_check(x text) LANGUAGE plpgsql AS $$ BEGIN IF x = 'guess' THEN RAISE EXCEPTION 'hit'; END IF; END; $$"},
		{Name: "probe_call_subquery", Definition: "CREATE FUNCTION public.probe_call_subquery() RETURNS int LANGUAGE plpgsql AS $$ BEGIN CALL proc_check((SELECT ssn FROM t LIMIT 1)); RETURN 1; END; $$"},
		{Name: "probe_call", Definition: "CREATE FUNCTION public.probe_call() RETURNS int LANGUAGE plpgsql AS $$ DECLARE v text; BEGIN SELECT ssn INTO v FROM t LIMIT 1; CALL proc_check(v); RETURN 1; END; $$"},
		// The last SELECT reads a table the snapshot lacks, so the body fails analysis.
		{Name: "probe_failing", Definition: "CREATE FUNCTION public.probe_failing() RETURNS SETOF record LANGUAGE sql AS $$ SELECT pg_sleep(0) FROM t WHERE ssn = 'guess'; SELECT * FROM missing_tbl $$"},
		{Name: "probe_for", Definition: "CREATE FUNCTION public.probe_for() RETURNS SETOF int LANGUAGE plpgsql AS $$ BEGIN FOR i IN 1..(SELECT ascii(ssn) FROM t LIMIT 1) LOOP RETURN NEXT i; END LOOP; END; $$"},
		{Name: "probe_join", Definition: "CREATE FUNCTION public.probe_join(guess text) RETURNS TABLE(id int) LANGUAGE plpgsql AS $$ BEGIN RETURN QUERY SELECT o.id FROM " + unanalyzable + "; END; $$"},
		{Name: "probe_join_sql", Definition: "CREATE FUNCTION public.probe_join_sql(guess text) RETURNS SETOF int LANGUAGE sql AS $$ SELECT o.id FROM " + unanalyzable + " $$"},
	}
	testCases := []struct {
		statement string
		want      []string
	}{
		{"SELECT count(*) AS c FROM probe('guess')", []string{"t.ssn"}},
		{"SELECT count(*) AS c FROM probe_sql('guess')", []string{"t.ssn"}},
		{"SELECT count(*) AS c FROM probe_fixed()", []string{"t.ssn"}},
		{"SELECT probe_multi() AS x", []string{"t.ssn"}},
		{"SELECT probe_perform() AS x", []string{"t.ssn"}},
		{"SELECT probe_if() AS x", []string{"t.ssn"}},
		{"SELECT probe_variable() AS x", []string{"t.id", "t.ssn"}},
		{"SELECT count(*) AS c FROM probe_for()", []string{"t.ssn"}},
		{"SELECT probe_update() AS x", []string{"t.ssn"}},
		{"SELECT probe_record() AS x", []string{"t.ssn"}},
		{"SELECT probe_raise() AS x", []string{"t.ssn"}},
		{"SELECT count(*) AS c FROM probe_stand_in('guess')", []string{"t.ssn"}},
		{"SELECT count(*) AS c FROM probe_stand_in_pl('guess')", []string{"t.ssn"}},
		{"SELECT probe_discard() AS x", []string{"t.ssn"}},
		{"SELECT probe_declare() AS x", []string{"t.ssn"}},
		{"SELECT probe_call() AS x", []string{"t.ssn"}},
		{"SELECT probe_call_subquery() AS x", []string{"t.ssn"}},
		{"SELECT count(*) AS c FROM probe_failing() AS f(a int)", []string{"t.ssn"}},
		{"SELECT probe_over('guess') AS b", []string{"t.ssn"}},
		{"SELECT probe_over(1) AS b", []string{}},
		{"SELECT * FROM o CROSS JOIN unknown_rows() u WHERE probe_over(1)", []string{"t.ssn"}},
		// omni leaves a set operation's LIMIT out of its tree.
		{"SELECT 1 AS x UNION ALL SELECT 2 LIMIT probe_count()", []string{"t.ssn"}},
		{"SELECT * FROM o CROSS JOIN unknown_rows() u WHERE probe_bool()", []string{"t.ssn"}},
		{"SELECT * FROM o CROSS JOIN unknown_rows() u WHERE public.probe_bool()", []string{"t.ssn"}},
		{"SELECT probe_execute() AS x", []string{"t.ssn"}},
		{"SELECT probe_exec() AS x", []string{"t.ssn"}},
		// * over an unknown function sends the outer query to the fallback.
		{"SELECT * FROM probe_table() p CROSS JOIN unknown_rows() u", []string{"t.ssn"}},
		{"SELECT count(*) AS c FROM probe_join('guess')", []string{"o.id", "t.id", "t.ssn"}},
		{"SELECT count(*) AS c FROM probe_join_sql('guess')", []string{"o.id", "t.id", "t.ssn"}},
	}
	for _, tc := range testCases {
		t.Run(tc.statement, func(t *testing.T) {
			span := spanFor(t, tc.statement, schema)
			require.Equal(t, tc.want, predicateColumnNames(span.PredicateColumns))
		})
	}
}

func TestGetQuerySpanResultLineage(t *testing.T) {
	testCases := []struct {
		name      string
		statement string
		want      []string
	}{
		{"set operation in subquery", "SELECT s.a FROM (SELECT ssn AS a FROM t UNION SELECT email FROM t) s", []string{"t.email", "t.ssn"}},
		{"set operation in cte", "WITH c AS (SELECT ssn FROM t UNION SELECT email FROM t) SELECT ssn FROM c", []string{"t.email", "t.ssn"}},
		{"except returns left branch rows", "SELECT s.a FROM (SELECT ssn AS a FROM t EXCEPT SELECT email FROM t) s", []string{"t.ssn"}},
		{"window before column in subquery", "SELECT s.x FROM (SELECT row_number() OVER (ORDER BY ssn) AS rn, email AS x FROM t) s", []string{"t.email"}},
		{"full join using column", "SELECT id FROM t FULL JOIN o USING (id)", []string{"o.id", "t.id"}},
		{"aggregate order by", "SELECT string_agg(email, ',' ORDER BY ssn) AS x FROM t", []string{"t.email", "t.ssn"}},
		{"unknown function in from", "SELECT f.a FROM t CROSS JOIN LATERAL unknown_rows(t.ssn) f", []string{"t.ssn"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			span := spanFor(t, tc.statement, healthySchema())
			require.Len(t, span.Results, 1)
			require.Equal(t, tc.want, predicateColumnNames(span.Results[0].SourceColumns))
		})
	}
}

func TestGetQuerySpanResultLineageOfFunctionBodies(t *testing.T) {
	schema := healthySchema()
	schema.Schemas[0].Functions = []*metadatapb.FunctionMetadata{
		{Name: "read_cursor", Definition: "CREATE FUNCTION public.read_cursor() RETURNS text LANGUAGE plpgsql AS $$ DECLARE c CURSOR FOR SELECT ssn FROM t; v text; BEGIN OPEN c; FETCH c INTO v; RETURN v; END; $$"},
		{Name: "loop_cursor", Definition: "CREATE FUNCTION public.loop_cursor() RETURNS SETOF text LANGUAGE plpgsql AS $$ DECLARE c CURSOR FOR SELECT ssn FROM t; BEGIN FOR r IN c LOOP RETURN NEXT r.ssn; END LOOP; END; $$"},
		{Name: "returning_sql", Definition: "CREATE FUNCTION public.returning_sql() RETURNS SETOF text LANGUAGE sql AS $$ UPDATE t SET id = id RETURNING ssn $$"},
		{Name: "returning_table", Definition: "CREATE FUNCTION public.returning_table() RETURNS TABLE(i int, e text, s text) LANGUAGE sql AS $$ DELETE FROM t WHERE false RETURNING * $$"},
		{Name: "loop_outer", Definition: "CREATE FUNCTION public.loop_outer() RETURNS text LANGUAGE plpgsql AS $$ DECLARE v text; BEGIN FOR i IN 1..1 LOOP SELECT ssn INTO v FROM t LIMIT 1; END LOOP; RETURN v; END; $$"},
		{Name: "block_outer", Definition: "CREATE FUNCTION public.block_outer() RETURNS text LANGUAGE plpgsql AS $$ DECLARE v text; BEGIN BEGIN SELECT ssn INTO v FROM t LIMIT 1; END; RETURN v; END; $$"},
		{Name: "returning_pl", Definition: "CREATE FUNCTION public.returning_pl() RETURNS text LANGUAGE plpgsql AS $$ DECLARE v text; BEGIN UPDATE t SET id = id WHERE id = 1 RETURNING ssn INTO v; RETURN v; END; $$"},
	}
	testCases := []struct {
		statement string
		want      []string
	}{
		{"SELECT read_cursor() AS x", []string{"t.ssn"}},
		{"SELECT loop_cursor() AS x", []string{"t.ssn"}},
		{"SELECT returning_sql() AS x", []string{"t.ssn"}},
		{"SELECT returning_pl() AS x", []string{"t.ssn"}},
		{"SELECT loop_outer() AS x", []string{"t.ssn"}},
		{"SELECT block_outer() AS x", []string{"t.ssn"}},
		{"SELECT r.s FROM returning_table() r", []string{"t.ssn"}},
	}
	for _, tc := range testCases {
		t.Run(tc.statement, func(t *testing.T) {
			span := spanFor(t, tc.statement, schema)
			require.Len(t, span.Results, 1)
			require.Equal(t, tc.want, predicateColumnNames(span.Results[0].SourceColumns))
		})
	}
}

func TestGetQuerySpanKeepsSQLBodyAccesses(t *testing.T) {
	schema := degradedSchema()
	schema.Schemas[0].Functions = []*metadatapb.FunctionMetadata{
		{Name: "probe_sql", Definition: "CREATE FUNCTION public.probe_sql() RETURNS int LANGUAGE sql AS $$ SELECT 1 FROM t WHERE secret = 'guess' $$"},
	}
	span := spanFor(t, "SELECT probe_sql() AS x", schema)
	require.Contains(t, span.SourceColumns, ColumnResource{Database: "db", Schema: "public", Table: "t"})
	require.NotNil(t, span.UnresolvedColumnsError, "t, read only in the body's condition, was synced with no columns")
}

func TestGetQuerySpanFallbackKeepsFunctionAccesses(t *testing.T) {
	schema := degradedSchema()
	schema.Schemas[0].Functions = []*metadatapb.FunctionMetadata{
		{Name: "probe_t", Definition: "CREATE FUNCTION public.probe_t() RETURNS TABLE(id int) LANGUAGE plpgsql AS $$ BEGIN RETURN QUERY SELECT 1 FROM t WHERE secret = 'guess'; END; $$"},
	}
	// * over an unknown function sends the statement to the fallback.
	span := spanFor(t, "SELECT * FROM probe_t() p CROSS JOIN unknown_rows() u", schema)
	require.Contains(t, span.SourceColumns, ColumnResource{Database: "db", Schema: "public", Table: "t"})
	require.NotNil(t, span.UnresolvedColumnsError, "t, read only inside the function, was synced with no columns")
}

func newFallbackTestExtractor(t *testing.T) *querySpanExtractor {
	t.Helper()
	e := newQuerySpanExtractor("db", nil, testCatalog(healthySchema()))
	e.ctx = context.Background()
	require.NoError(t, e.initCatalog())
	return e
}

func TestCollectFallbackPredicateColumns(t *testing.T) {
	testCases := []struct {
		name      string
		statement string
		want      []string
	}{
		{"where", "SELECT id FROM t WHERE ssn = '1'", []string{"t.ssn"}},
		{"qualified reference matches its table", "SELECT o.amt FROM o JOIN t ON o.id = t.id WHERE o.amt > 0", []string{"o.amt", "o.id", "t.id"}},
		{"unqualified reference matches every table", "SELECT 1 FROM o JOIN t USING (id) WHERE id > 0", []string{"o.id", "t.id"}},
		{"reused alias matches every table it names", "SELECT 1 FROM t x WHERE EXISTS (SELECT 1 FROM o x WHERE x.id = 1)", []string{"o.id", "t.id"}},
		{"subquery in a condition reads every name in it", "SELECT 1 FROM o WHERE EXISTS (WITH c AS (SELECT id FROM t WHERE ssn = '1') SELECT 1 FROM c)", []string{"o.id", "t.id", "t.ssn"}},
		{"cte column alias", "WITH c(a, b) AS (SELECT id, ssn FROM t) SELECT a FROM c WHERE b = '1'", []string{"t.ssn"}},
		{"recursive cte column alias", "WITH RECURSIVE c(a, n) AS (SELECT ssn, 1 FROM t UNION ALL SELECT a, n + 1 FROM c WHERE n < 3) SELECT a FROM c", []string{}},
		{"derived column alias", "SELECT 1 FROM (SELECT id, ssn FROM t) s(a, b) WHERE s.b = '1'", []string{"t.ssn"}},
		{"output name", "SELECT 1 FROM (SELECT ssn AS x FROM t) s WHERE x = '1'", []string{"t.ssn"}},
		{"whole row", "SELECT 1 FROM t WHERE t::text LIKE '%1%'", []string{"t.email", "t.id", "t.ssn"}},
		{"aliased whole row", "SELECT 1 FROM t x WHERE ROW(x.*) IS NOT NULL", []string{"t.email", "t.id", "t.ssn"}},
		{"join using", "SELECT 1 FROM t a JOIN t b USING (ssn)", []string{"t.ssn"}},
		{"having", "SELECT id FROM t GROUP BY id HAVING max(ssn) > '1'", []string{"t.ssn"}},
		{"aggregate filter", "SELECT count(*) FILTER (WHERE ssn = '1') FROM t", []string{"t.ssn"}},
		{"natural join reads both sides", "SELECT 1 FROM o NATURAL JOIN (SELECT id FROM t) x", []string{"o.amt", "o.id", "t.id"}},
		{"select list only", "SELECT ssn FROM t", []string{}},
		{"except filters by the other branch", "SELECT 'guess' EXCEPT SELECT ssn FROM t", []string{"t.ssn"}},
		{"intersect over star", "SELECT 1, 'a', 'b' INTERSECT SELECT * FROM t", []string{"t.email", "t.id", "t.ssn"}},
		{"union filters nothing", "SELECT 'guess' UNION SELECT ssn FROM t", []string{}},
		{"limit value", "SELECT * FROM t LIMIT (SELECT CASE WHEN ssn = 'guess' THEN 1 ELSE 0 END FROM t LIMIT 1)", []string{"t.ssn"}},
		{"function in from", "SELECT * FROM t CROSS JOIN LATERAL generate_series(1, ascii(substr(t.ssn, 1, 1))) g", []string{"t.ssn"}},
		{"set-returning function in select list", "SELECT generate_series(1, ascii(substr(ssn, 1, 1))) FROM t", []string{"t.ssn"}},
		{"scalar function in select list", "SELECT upper(ssn) FROM t", []string{}},
	}
	accesses := ColumnSet{
		{Database: "db", Schema: "public", Table: "t"}: true,
		{Database: "db", Schema: "public", Table: "o"}: true,
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := newFallbackTestExtractor(t)
			stmts, err := omnipg.Parse(tc.statement)
			require.NoError(t, err)
			e.collectFallbackPredicateColumns(stmts[0].AST, accesses)
			require.Equal(t, tc.want, predicateColumnNames(e.predicateColumns))
		})
	}
}

func TestGetQuerySpanFallbackKeepsPredicateColumns(t *testing.T) {
	// * over the unknown columns of a function fails omni's analysis.
	statement := "SELECT * FROM o JOIN t ON t.id = o.id AND t.ssn = '1' CROSS JOIN unknown_rows() f"
	stmts, err := omnipg.Parse(statement)
	require.NoError(t, err)
	_, err = newFallbackTestExtractor(t).analyzeSelect(stmts[0].AST.(*ast.SelectStmt))
	require.Error(t, err, "the statement must fail analysis to reach the fallback")

	span := spanFor(t, statement, healthySchema())
	require.Equal(t, []string{"o.id", "t.id", "t.ssn"}, predicateColumnNames(span.PredicateColumns))
}

func TestGetQuerySpanTableFunctionArgumentLineage(t *testing.T) {
	schema := healthySchema()
	schema.Schemas[0].Functions = []*metadatapb.FunctionMetadata{
		{Name: "echo", Definition: "CREATE FUNCTION public.echo(x text) RETURNS TABLE(value text) LANGUAGE sql AS $$ SELECT x $$"},
		// A missing column type makes the catalog install a stand-in without
		// output parameters, so the outputs come from the definition text.
		{Name: "echo_stand_in", Definition: "CREATE FUNCTION public.echo_stand_in(x text) RETURNS TABLE(value missing_type) LANGUAGE sql AS $$ SELECT x $$"},
		{Name: "echo_pl", Definition: "CREATE FUNCTION public.echo_pl(x text) RETURNS TABLE(value text) LANGUAGE plpgsql AS $$ BEGIN RETURN QUERY SELECT x; END; $$"},
	}
	testCases := []struct {
		statement string
		want      []string
	}{
		{"SELECT f.value FROM t CROSS JOIN LATERAL echo(t.ssn) f", []string{"t.ssn"}},
		{"SELECT f.value FROM t CROSS JOIN LATERAL echo_pl(t.ssn) f", []string{"t.ssn"}},
		{"SELECT f.value FROM t, echo(upper(t.email)) f", []string{"t.email"}},
		// * over an unknown function sends the statement to the fallback.
		{"SELECT * FROM echo('x') f CROSS JOIN LATERAL echo((SELECT ssn FROM t LIMIT 1)) g CROSS JOIN unknown_rows() u", []string{"t.ssn"}},
		{"SELECT * FROM echo((SELECT ssn FROM t LIMIT 1)) f, unknown_rows() u", []string{"t.ssn"}},
		{"SELECT * FROM echo_stand_in((SELECT ssn FROM t LIMIT 1)) f CROSS JOIN unknown_rows() u", []string{"t.ssn"}},
	}
	for _, tc := range testCases {
		t.Run(tc.statement, func(t *testing.T) {
			span := spanFor(t, tc.statement, schema)
			require.NotEmpty(t, span.Results)
			got := ColumnSet{}
			for _, result := range span.Results {
				mergeSourceColumnSetInto(got, result.SourceColumns)
			}
			require.Equal(t, tc.want, predicateColumnNames(got))
		})
	}
}

func TestGetQuerySpanReadsFunctionsOutsideTheOutput(t *testing.T) {
	schema := healthySchema()
	schema.Schemas[0].Functions = []*metadatapb.FunctionMetadata{
		{Name: "probe", Definition: "CREATE FUNCTION public.probe() RETURNS int LANGUAGE sql AS $$ SELECT count(*)::int FROM t WHERE ssn = 'guess' $$"},
	}
	for _, statement := range []string{
		"SELECT id FROM o ORDER BY probe()",
		"SELECT id FROM o GROUP BY id, probe()",
		"SELECT DISTINCT ON (probe()) id FROM o",
		"SELECT row_number() OVER (ORDER BY probe()) AS n FROM o",
		"SELECT row_number() OVER (PARTITION BY probe()) AS n FROM o",
		"SELECT s.id FROM (SELECT id FROM o ORDER BY probe() LIMIT 1) s",
	} {
		t.Run(statement, func(t *testing.T) {
			span := spanFor(t, statement, schema)
			require.Contains(t, span.SourceColumns, ColumnResource{Database: "db", Schema: "public", Table: "t"})
			require.Equal(t, []string{"t.ssn"}, predicateColumnNames(span.PredicateColumns))
		})
	}
}

func TestGetQuerySpanOutputParameterLineage(t *testing.T) {
	schema := healthySchema()
	schema.Schemas[0].Functions = []*metadatapb.FunctionMetadata{
		{Name: "next_out", Definition: "CREATE FUNCTION public.next_out() RETURNS TABLE(value text) LANGUAGE plpgsql AS $$ BEGIN SELECT ssn INTO value FROM t LIMIT 1; RETURN NEXT; END; $$"},
		{Name: "next_loop", Definition: "CREATE FUNCTION public.next_loop() RETURNS TABLE(value text) LANGUAGE plpgsql AS $$ BEGIN FOR value IN SELECT ssn FROM t LOOP RETURN NEXT; END LOOP; END; $$"},
		{Name: "next_assign", Definition: "CREATE FUNCTION public.next_assign() RETURNS TABLE(i int, value text) LANGUAGE plpgsql AS $$ BEGIN i := 1; BEGIN value := (SELECT ssn FROM t LIMIT 1); END; RETURN NEXT; END; $$"},
		// A set-returning function's bare RETURN ends it without a row.
		{Name: "srf_return", Definition: "CREATE FUNCTION public.srf_return() RETURNS TABLE(value text) LANGUAGE plpgsql AS $$ BEGIN SELECT ssn INTO value FROM t LIMIT 1; RETURN; END; $$"},
		{Name: "next_shadowed", Definition: "CREATE FUNCTION public.next_shadowed() RETURNS TABLE(value text) LANGUAGE plpgsql AS $$ BEGIN value := 'public'; DECLARE value text; BEGIN SELECT ssn INTO value FROM t LIMIT 1; END; RETURN NEXT; END; $$"},
		{Name: "next_alias", Definition: "CREATE FUNCTION public.next_alias() RETURNS TABLE(value text) LANGUAGE plpgsql AS $$ DECLARE v ALIAS FOR value; BEGIN SELECT ssn INTO v FROM t LIMIT 1; RETURN NEXT; END; $$"},
		{Name: "next_targets", Definition: "CREATE FUNCTION public.next_targets() RETURNS TABLE(a text, b text) LANGUAGE plpgsql AS $$ BEGIN FOR a, b IN SELECT email, ssn FROM t LOOP RETURN NEXT; END LOOP; END; $$"},
		{Name: "next_int_loop", Definition: "CREATE FUNCTION public.next_int_loop() RETURNS TABLE(value int) LANGUAGE plpgsql AS $$ BEGIN value := 0; FOR value IN 1..1 LOOP value := (SELECT id FROM t LIMIT 1); END LOOP; RETURN NEXT; END; $$"},
		// The second iteration returns the value the first assigned.
		{Name: "next_before_assign", Definition: "CREATE FUNCTION public.next_before_assign() RETURNS TABLE(value text) LANGUAGE plpgsql AS $$ BEGIN FOR i IN 1..2 LOOP RETURN NEXT; SELECT ssn INTO value FROM t LIMIT 1; END LOOP; END; $$"},
		{Name: "query_after_loop", Definition: "CREATE FUNCTION public.query_after_loop() RETURNS TABLE(value text) LANGUAGE plpgsql AS $$ BEGIN FOR i IN 1..1 LOOP SELECT ssn INTO value FROM t LIMIT 1; END LOOP; RETURN QUERY SELECT value; END; $$"},
		{Name: "next_cursor_loop", Definition: "CREATE FUNCTION public.next_cursor_loop() RETURNS TABLE(value text) LANGUAGE plpgsql AS $$ DECLARE c CURSOR FOR SELECT ssn FROM t; BEGIN value := 'public'; FOR value IN c LOOP RETURN NEXT; END LOOP; END; $$"},
		{Name: "out_implicit", Definition: "CREATE FUNCTION public.out_implicit(OUT a text, OUT b text) LANGUAGE plpgsql AS $$ BEGIN SELECT email, ssn INTO a, b FROM t LIMIT 1; END; $$"},
		{Name: "out_return", Definition: "CREATE FUNCTION public.out_return(OUT a text, OUT b text) RETURNS record LANGUAGE plpgsql AS $$ BEGIN SELECT email, ssn INTO a, b FROM t LIMIT 1; RETURN; END; $$"},
	}
	testCases := []struct {
		statement string
		want      [][]string
	}{
		{"SELECT value FROM next_out()", [][]string{{"t.ssn"}}},
		{"SELECT value FROM next_loop()", [][]string{{"t.ssn"}}},
		{"SELECT i, value FROM next_assign()", [][]string{{}, {"t.ssn"}}},
		{"SELECT value FROM srf_return()", [][]string{{}}},
		{"SELECT value FROM next_shadowed()", [][]string{{}}},
		{"SELECT value FROM next_alias()", [][]string{{"t.ssn"}}},
		{"SELECT a, b FROM next_targets()", [][]string{{"t.email"}, {"t.ssn"}}},
		{"SELECT value FROM next_int_loop()", [][]string{{}}},
		{"SELECT value FROM next_before_assign()", [][]string{{"t.ssn"}}},
		{"SELECT value FROM query_after_loop()", [][]string{{"t.ssn"}}},
		{"SELECT value FROM next_cursor_loop()", [][]string{{}}},
		{"SELECT a, b FROM out_implicit()", [][]string{{"t.email"}, {"t.ssn"}}},
		{"SELECT a, b FROM out_return()", [][]string{{"t.email"}, {"t.ssn"}}},
	}
	for _, tc := range testCases {
		t.Run(tc.statement, func(t *testing.T) {
			span := spanFor(t, tc.statement, schema)
			var got [][]string
			for _, result := range span.Results {
				got = append(got, predicateColumnNames(result.SourceColumns))
			}
			require.Equal(t, tc.want, got)
		})
	}
}

func TestGetQuerySpanValuesLineage(t *testing.T) {
	testCases := []struct {
		statement  string
		results    map[string][]string
		predicates []string
	}{
		{"VALUES ((SELECT ssn FROM t LIMIT 1))", map[string][]string{"column1": {"t.ssn"}}, []string{}},
		{
			"VALUES (1, (SELECT email FROM t WHERE ssn = 'guess')), ((SELECT id FROM t LIMIT 1), 'x')",
			map[string][]string{"column1": {"t.id"}, "column2": {"t.email"}},
			[]string{"t.ssn"},
		},
		{"VALUES ('a'), ('b'), ('c'), ((SELECT ssn FROM t LIMIT 1)), ('e')", map[string][]string{"column1": {"t.ssn"}}, []string{}},
		{"SELECT * FROM (VALUES ((SELECT ssn FROM t LIMIT 1))) v", map[string][]string{"column1": {"t.ssn"}}, []string{}},
		{"WITH c(x) AS (VALUES ((SELECT ssn FROM t LIMIT 1))) SELECT x FROM c", map[string][]string{"x": {"t.ssn"}}, []string{}},
		{"VALUES (NULL), (NULL), (1), ((SELECT id FROM t LIMIT 1))", map[string][]string{"column1": {"t.id"}}, []string{}},
		{"SELECT s.ssn FROM (SELECT ssn FROM t) s WHERE EXISTS (VALUES (NULL), (NULL), (1), (2))", map[string][]string{"ssn": {"t.ssn"}}, []string{}},
		{"SELECT 'a' AS x UNION ALL VALUES ((SELECT ssn FROM t LIMIT 1))", map[string][]string{"x": {"t.ssn"}}, []string{}},
		{"SELECT (VALUES ((SELECT ssn FROM t WHERE id = 1))) AS x", map[string][]string{"x": {"t.ssn"}}, []string{"t.id"}},
	}
	for _, tc := range testCases {
		t.Run(tc.statement, func(t *testing.T) {
			span := spanFor(t, tc.statement, healthySchema())
			got := map[string][]string{}
			for _, result := range span.Results {
				got[result.Name] = predicateColumnNames(result.SourceColumns)
			}
			require.Equal(t, tc.results, got)
			require.Equal(t, tc.predicates, predicateColumnNames(span.PredicateColumns))
		})
	}

	span := spanFor(t, "EXPLAIN ANALYZE VALUES ((SELECT email FROM t WHERE ssn = 'guess'))", healthySchema())
	require.Equal(t, []string{"t.ssn"}, predicateColumnNames(span.PredicateColumns))
}

func TestGetQuerySpanMultiFunctionLineage(t *testing.T) {
	schema := healthySchema()
	schema.Schemas[0].Functions = []*metadatapb.FunctionMetadata{
		{Name: "echo", Definition: "CREATE FUNCTION public.echo(x text) RETURNS TABLE(value text) LANGUAGE sql AS $$ SELECT x $$"},
		{Name: "pair", Definition: "CREATE FUNCTION public.pair(k int) RETURNS TABLE(p text, q text) LANGUAGE sql AS $$ SELECT email, email FROM t WHERE id = k $$"},
	}
	testCases := []struct {
		statement string
		want      [][]string
	}{
		{"SELECT f.a, f.b FROM t CROSS JOIN LATERAL unnest(ARRAY[t.email], ARRAY[t.ssn]) AS f(a, b)", [][]string{{"t.email"}, {"t.ssn"}}},
		{"SELECT f.a, f.b, f.n FROM t CROSS JOIN LATERAL unnest(ARRAY[t.email], ARRAY[t.ssn]) WITH ORDINALITY AS f(a, b, n)", [][]string{{"t.email"}, {"t.ssn"}, {}}},
		{"SELECT f.a, f.b FROM t CROSS JOIN LATERAL ROWS FROM (unnest(ARRAY[t.email]), echo(t.ssn)) AS f(a, b)", [][]string{{"t.email"}, {"t.ssn"}}},
		{"SELECT f.x, f.y, f.z FROM t CROSS JOIN LATERAL ROWS FROM (pair(t.id), unnest(ARRAY[t.ssn])) AS f(x, y, z)", [][]string{{"t.email", "t.id"}, {"t.email", "t.id"}, {"t.ssn"}}},
	}
	for _, tc := range testCases {
		t.Run(tc.statement, func(t *testing.T) {
			span := spanFor(t, tc.statement, schema)
			var got [][]string
			for _, result := range span.Results {
				got = append(got, predicateColumnNames(result.SourceColumns))
			}
			require.Equal(t, tc.want, got)
		})
	}
}
