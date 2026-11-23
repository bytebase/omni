package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/oracle/internal/testcontainer"
	oracleparser "github.com/bytebase/omni/oracle/parser"
)

// TestQuerySimulation holds the catalog to testdata/queries/simulate.tsv:
// queries wrapped as CREATE VIEW v AS <query>, with what Oracle did and what
// Exec does. A row's outcome is the ORA- error number Oracle refuses the view
// with, or ok: and a hash of every view Oracle stores afterwards: its derived
// columns (name, data type, nullability, visibility) read from USER_TAB_COLS
// and its stored text from USER_VIEWS.
//
// The rows are every generated SELECT and WITH statement Oracle accepts in the
// parser's compat.tsv, on one_table@column, and a corpus of Oracle query forms
// (queryCorpus) on queries@emp.
//
// Oracle's side is on the test's tape, so the test starts no server.
// OMNI_RECORD=1 runs every row on the pinned Oracle again; -update-compat
// commits the rows a catalog fix moved; and -simulate-show=<text> prints both
// schemas of the rows that hold the text.
func TestQuerySimulation(t *testing.T) {
	rows, err := querySimulationRows(filepath.Join("..", "parser", "testdata", "compat.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "queries", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     rows,
		MadeUp:   (testcontainer.Codec{}).MadeUp(),
		Omni:     querySimulationOmni,
		Engine:   querySimulationEngine,
	})
}

const (
	queryCompatFixture = "one_table@column"
	queryCorpusFixture = "queries@emp"
)

func querySimulationRows(compatPath string) ([]simulate.Generated, error) {
	rows, err := simulate.FromCompat(compatPath, []string{queryCompatFixture}, func(statement string) bool {
		return strings.HasPrefix(statement, "SELECT ") || strings.HasPrefix(statement, "WITH ")
	})
	if err != nil {
		return nil, err
	}
	// Only the grammar's rows: compat.tsv's hand rows (suite:, issue:) are
	// corpora that grow on their own and name tables no fixture holds.
	rows = slices.DeleteFunc(rows, func(g simulate.Generated) bool { return handLabel(g.Rule) })
	seen := map[string]bool{}
	var out []simulate.Generated
	add := func(g simulate.Generated) {
		key := g.Fixture + "\x00" + g.Statement
		if seen[key] {
			return
		}
		seen[key] = true
		out = append(out, g)
	}
	for _, r := range rows {
		r.Statement = "CREATE VIEW v AS " + strings.TrimSpace(strings.TrimRight(strings.TrimSpace(r.Statement), ";"))
		add(r)
	}
	for _, q := range queryCorpus {
		add(simulate.Generated{Fixture: queryCorpusFixture, Rule: "corpus/" + q.name, Statement: "CREATE VIEW v AS " + q.query})
	}
	return out, nil
}

// querySimulationEngine reads every view the scratch schema holds: its
// columns in USER_TAB_COLS order and its stored text.
func querySimulationEngine(t testing.TB) func(simulate.Fixture, string) (simulate.Answer, error) {
	return simulationEngineWith(t, func(ctx context.Context, db *sql.DB) (simulate.Answer, error) {
		rows, err := db.QueryContext(ctx, `SELECT v.view_name, v.text_vc, v.read_only,
  c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
  c.char_used, c.char_length, c.nullable, c.hidden_column
FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
ORDER BY v.view_name, c.internal_column_id`)
		if err != nil {
			return simulate.Answer{}, fmt.Errorf("read stored views: %w", err)
		}
		defer rows.Close()
		var views []*queryView
		for rows.Next() {
			var (
				name, text, readOnly, col, typ, charUsed, nullable, hidden sql.NullString
				length, precision, scale, charLength                       sql.NullInt64
			)
			if err := rows.Scan(&name, &text, &readOnly, &col, &typ, &length, &precision, &scale, &charUsed, &charLength, &nullable, &hidden); err != nil {
				return simulate.Answer{}, err
			}
			if len(views) == 0 || views[len(views)-1].name != name.String {
				views = append(views, &queryView{name: name.String, text: text.String, readOnly: readOnly.String == "Y"})
			}
			v := views[len(views)-1]
			v.columns = append(v.columns, queryViewColumn{
				name:      col.String,
				typ:       dictionaryType(typ.String, length, precision, scale, charUsed.String, charLength),
				notNull:   nullable.String == "N",
				invisible: hidden.String == "YES",
			})
		}
		if err := rows.Err(); err != nil {
			return simulate.Answer{}, err
		}
		var forms []simulate.Form
		for _, v := range views {
			forms = append(forms, v.form())
		}
		return simulate.Answer{Forms: forms}, nil
	})
}

type queryView struct {
	name     string
	text     string
	readOnly bool
	columns  []queryViewColumn
}

type queryViewColumn struct {
	name, typ          string
	notNull, invisible bool
}

// form renders a view as one line: its columns, each with its type, NOT NULL
// and INVISIBLE, and its stored text.
func (v *queryView) form() simulate.Form {
	var b strings.Builder
	fmt.Fprintf(&b, "CREATE VIEW %s (", quoteIdent(v.name))
	for i, c := range v.columns {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%s %s", quoteIdent(c.name), c.typ)
		if c.invisible {
			b.WriteString(" INVISIBLE")
		}
		if c.notNull {
			b.WriteString(" NOT NULL")
		}
	}
	text := strings.TrimRight(strings.TrimSpace(v.text), "; \r\n")
	b.WriteString(") AS ")
	b.WriteString(text)
	if v.readOnly && !strings.Contains(strings.ToUpper(text), "READ ONLY") {
		b.WriteString(" WITH READ ONLY")
	}
	b.WriteString(";\n")
	return simulate.Form{Kind: "VIEW", Name: v.name, Text: b.String()}
}

// dictionaryType renders a USER_TAB_COLS type as the catalog writes a
// column's type.
func dictionaryType(typ string, length, precision, scale sql.NullInt64, charUsed string, charLength sql.NullInt64) string {
	switch typ {
	case "VARCHAR2", "CHAR":
		if charUsed == "C" && charLength.Valid {
			return fmt.Sprintf("%s(%d CHAR)", typ, charLength.Int64)
		}
		return fmt.Sprintf("%s(%d BYTE)", typ, length.Int64)
	case "NVARCHAR2", "NCHAR":
		if charLength.Valid && charLength.Int64 > 0 {
			return fmt.Sprintf("%s(%d)", typ, charLength.Int64)
		}
		return fmt.Sprintf("%s(%d)", typ, length.Int64)
	case "RAW", "UROWID":
		return fmt.Sprintf("%s(%d)", typ, length.Int64)
	case "NUMBER":
		switch {
		case !precision.Valid && scale.Valid:
			return fmt.Sprintf("NUMBER(*,%d)", scale.Int64)
		case !precision.Valid:
			return "NUMBER"
		case scale.Valid && scale.Int64 != 0:
			return fmt.Sprintf("NUMBER(%d,%d)", precision.Int64, scale.Int64)
		default:
			return fmt.Sprintf("NUMBER(%d)", precision.Int64)
		}
	case "FLOAT":
		if precision.Valid && precision.Int64 > 0 {
			return fmt.Sprintf("FLOAT(%d)", precision.Int64)
		}
		return "FLOAT"
	}
	return typ
}

// querySimulationOmni runs the statement on the fixture in the catalog and
// renders every view it then holds as the engine's side does.
func querySimulationOmni(fx simulate.Fixture, statement string, engine simulate.Answer) simulate.Omni {
	c, err := LoadSQL(fx.SQL)
	if err != nil {
		panic(fmt.Sprintf("the catalog does not load fixture %s: %v", fx.Name, err))
	}
	c.meta.Name = simUser
	if len(c.meta.Schemas) > 0 {
		c.meta.Schemas[0].Name = simUser
	}
	c.SetServerConfig(StandardOracle23ServerConfig())
	results, err := c.Exec(statement, nil)
	if err != nil {
		var refusal *Error
		if errors.As(err, &refusal) {
			if code := oracleErrorCode(refusal.Code); code != "" {
				return simulate.Omni{Code: code, Why: err.Error()}
			}
		}
		var syntax *oracleparser.ParseError
		if errors.As(err, &syntax) && oracleSyntaxCode(engine.Code) {
			// omni's parser refuses the statement, and so does Oracle's.
			return simulate.Omni{Code: engine.Code, Why: err.Error()}
		}
		return simulate.Omni{Code: simulate.NoCode, Why: err.Error()}
	}
	for _, r := range results {
		var refusal *Error
		switch {
		case r.Error == nil:
		case errors.As(r.Error, &refusal):
			if code := oracleErrorCode(refusal.Code); code != "" {
				return simulate.Omni{Code: code, Why: r.Error.Error()}
			}
			return simulate.Omni{Code: simulate.NoCode, Why: r.Error.Error()}
		default:
			return simulate.Omni{Code: simulate.NoCode, Why: r.Error.Error()}
		}
	}
	if engine.Refused() {
		return simulate.Omni{}
	}
	var got []string
	for _, s := range c.Metadata().GetSchemas() {
		for _, v := range s.GetViews() {
			qv := &queryView{name: v.GetName(), text: v.GetDefinition()}
			for _, col := range v.GetColumns() {
				qv.columns = append(qv.columns, queryViewColumn{
					name:      col.GetName(),
					typ:       omniDictionaryType(col.GetType()),
					notNull:   !col.GetNullable(),
					invisible: col.GetIsInvisible(),
				})
			}
			got = append(got, qv.form().Text)
		}
	}
	slices.Sort(got)
	var want []string
	for _, f := range engine.Forms {
		want = append(want, f.Text)
	}
	slices.Sort(want)
	if slices.Equal(got, want) {
		return simulate.Omni{Same: true}
	}
	return simulate.Omni{Why: "the catalog stores:\n" + strings.Join(got, "") + "Oracle stores:\n" + strings.Join(want, "")}
}

// oracleSyntaxCode reports an error Oracle's SQL parser raises, as the
// compat judge classifies it.
func oracleSyntaxCode(code string) bool {
	n, err := strconv.Atoi(code)
	return err == nil && compat.OracleParserCode(n)
}

// omniDictionaryType renders a catalog column type as USER_TAB_COLS reads it
// back: INTEGER and its kin are NUMBER(*,0).
func omniDictionaryType(typ string) string {
	switch strings.ToUpper(strings.TrimSpace(typ)) {
	case "INTEGER", "INT", "SMALLINT":
		return "NUMBER(*,0)"
	}
	return typ
}

// queryCorpus is the Oracle query forms the grammar's rows do not reach with
// real tables: joins, outer joins in both notations, CTEs and recursive WITH,
// set operators, aggregates and analytic functions, CASE/DECODE/NVL and other
// built-ins, hierarchical queries, row limiting, PIVOT/UNPIVOT, and the
// refusals of each. They run on queries@emp, whose tables are dept, emp and
// misc and whose view is emp_v.
var queryCorpus = []struct{ name, query string }{
	// Select lists and naming.
	{"star", "SELECT * FROM emp"},
	{"qualified_star", "SELECT e.* FROM emp e"},
	{"star_two_tables_dup", "SELECT * FROM emp, dept"},
	{"star_and_column", "SELECT e.*, d.dname FROM emp e, dept d WHERE e.deptno = d.deptno"},
	{"table_name_qualifier", "SELECT emp.ename, dept.dname FROM emp, dept WHERE emp.deptno = dept.deptno"},
	{"schema_qualified", "SELECT omni_sim.emp.ename FROM omni_sim.emp"},
	{"alias_hides_table", "SELECT emp.ename FROM emp e"},
	{"alias_star_hides_table", "SELECT emp.* FROM emp e"},
	{"quoted_alias", `SELECT ename "Name", job "job title" FROM emp`},
	{"quoted_column_case", `SELECT "ename" FROM emp`},
	{"as_alias", "SELECT ename AS name, job AS title FROM emp"},
	{"dup_column", "SELECT ename, ename FROM emp"},
	{"dup_alias", "SELECT ename AS a, job AS a FROM emp"},
	{"dup_self_join", "SELECT e.ename, m.ename FROM emp e, emp m WHERE e.mgr = m.empno"},
	{"unaliased_literal", "SELECT 1 FROM dual"},
	{"unaliased_expression", "SELECT sal + 1 FROM emp"},
	{"unaliased_function", "SELECT UPPER(ename) FROM emp"},
	{"missing_column", "SELECT nosuch FROM emp"},
	{"missing_table", "SELECT ename FROM nosuch"},
	{"missing_qualified_column", "SELECT e.nosuch FROM emp e"},
	{"missing_qualifier", "SELECT x.ename FROM emp e"},
	{"missing_where_column", "SELECT ename FROM emp WHERE nosuch = 1"},
	{"missing_order_column", "SELECT ename FROM emp ORDER BY nosuch"},
	{"missing_group_column", "SELECT ename FROM emp GROUP BY nosuch"},
	{"alias_in_where", "SELECT ename e1 FROM emp WHERE e1 = 'X'"},
	{"order_by_position_out_of_range", "SELECT ename FROM emp ORDER BY 2"},
	{"ambiguous_column", "SELECT deptno FROM emp, dept"},
	{"ambiguous_self_join", "SELECT sal FROM emp e JOIN emp m ON e.mgr = m.empno"},
	{"view_over_view", "SELECT ename FROM emp_v"},
	{"view_star", "SELECT * FROM emp_v"},
	{"view_join_dup", "SELECT * FROM emp_v v JOIN dept d ON v.deptno = d.deptno"},
	{"rownum_alias", "SELECT ROWNUM rn, ename FROM emp"},
	{"rownum_bare", "SELECT ROWNUM, ename FROM emp"},
	{"rowid_alias", "SELECT ROWID rid, ename FROM emp"},
	{"distinct", "SELECT DISTINCT job FROM emp"},
	{"unique", "SELECT UNIQUE job FROM emp"},
	{"all_types", "SELECT * FROM misc"},
	{"all_types_listed", "SELECT id, code, title, name, note, ts, tz, bd, bf, r, fl, ym, ds, n, flag FROM misc"},
	{"bind_variable", "SELECT ename FROM emp WHERE deptno = :1"},
	{"sample", "SELECT ename FROM emp SAMPLE (10)"},
	{"read_only", "SELECT ename FROM emp WITH READ ONLY"},
	{"check_option", "SELECT ename, sal FROM emp WHERE deptno = 10 WITH CHECK OPTION"},
	{"check_option_named", "SELECT ename, sal FROM emp WHERE deptno = 10 WITH CHECK OPTION CONSTRAINT v_chk"},
	{"for_update", "SELECT ename FROM emp FOR UPDATE"},
	{"order_by", "SELECT ename, sal FROM emp ORDER BY sal DESC, ename"},

	// Joins.
	{"join_on", "SELECT e.ename, d.dname FROM emp e JOIN dept d ON e.deptno = d.deptno"},
	{"inner_join_paren", "SELECT e.ename FROM emp e INNER JOIN dept d ON (e.deptno = d.deptno)"},
	{"join_using", "SELECT ename, dname FROM emp JOIN dept USING (deptno)"},
	{"join_using_star", "SELECT * FROM emp JOIN dept USING (deptno)"},
	{"join_using_qualified", "SELECT e.deptno FROM emp e JOIN dept d USING (deptno)"},
	{"join_using_missing", "SELECT e.ename FROM emp e JOIN dept d USING (nosuch)"},
	{"natural_join", "SELECT * FROM emp NATURAL JOIN dept"},
	{"natural_join_view", "SELECT * FROM dept NATURAL JOIN emp_v"},
	{"natural_join_many", "SELECT empno FROM emp NATURAL JOIN emp_v"},
	{"left_join", "SELECT e.ename, d.dname FROM emp e LEFT OUTER JOIN dept d ON e.deptno = d.deptno"},
	{"right_join", "SELECT e.ename, d.dname FROM emp e RIGHT JOIN dept d ON e.deptno = d.deptno"},
	{"full_join", "SELECT e.ename, d.dname FROM emp e FULL OUTER JOIN dept d ON e.deptno = d.deptno"},
	{"full_join_using_star", "SELECT * FROM emp FULL JOIN dept USING (deptno)"},
	{"join_using_left_star", "SELECT * FROM emp LEFT JOIN dept USING (deptno)"},
	{"join_using_right_star", "SELECT * FROM emp RIGHT JOIN dept USING (deptno)"},
	{"join_using_reverse_star", "SELECT * FROM dept JOIN emp USING (deptno)"},
	{"join_using_types", "SELECT * FROM (SELECT CAST(1 AS NUMBER(3)) k FROM dual) JOIN (SELECT CAST(1 AS NUMBER(5,1)) k FROM dual) USING (k)"},
	{"join_using_left_types", "SELECT * FROM (SELECT CAST(1 AS NUMBER(3)) k FROM dual) LEFT JOIN (SELECT CAST(1 AS NUMBER(5,1)) k FROM dual) USING (k)"},
	{"self_left_join", "SELECT e.ename, m.ename AS mgr_name FROM emp e LEFT JOIN emp m ON e.mgr = m.empno"},
	{"cross_join", "SELECT e.deptno FROM emp e CROSS JOIN dept d"},
	{"join_three", "SELECT e.ename, d.dname, m.job FROM emp e JOIN dept d ON e.deptno = d.deptno JOIN emp m ON m.empno = e.mgr"},
	{"join_on_extra", "SELECT ename FROM emp e JOIN dept d ON e.deptno = d.deptno AND d.loc = 'DALLAS'"},
	{"comma_join", "SELECT e.ename, d.loc FROM emp e, dept d WHERE e.deptno = d.deptno"},
	{"plus_outer_join", "SELECT e.ename, d.dname FROM emp e, dept d WHERE e.deptno = d.deptno(+)"},
	{"plus_outer_join_right", "SELECT e.ename, d.dname, d.loc FROM emp e, dept d WHERE e.deptno(+) = d.deptno"},
	{"plus_outer_join_filter", "SELECT ename FROM emp e, dept d WHERE e.deptno = d.deptno(+) AND d.loc(+) = 'X'"},
	{"plus_both_sides", "SELECT e.ename FROM emp e, dept d WHERE e.deptno(+) = d.deptno(+)"},
	{"plus_in_ansi_join", "SELECT e.ename FROM emp e LEFT JOIN dept d ON e.deptno = d.deptno(+)"},
	{"cross_apply", "SELECT d.dname, e.ename FROM dept d CROSS APPLY (SELECT ename FROM emp WHERE emp.deptno = d.deptno) e"},
	{"outer_apply", "SELECT d.dname, x.cnt FROM dept d OUTER APPLY (SELECT COUNT(*) cnt FROM emp e WHERE e.deptno = d.deptno) x"},
	{"lateral", "SELECT d.dname, x.cnt FROM dept d, LATERAL (SELECT COUNT(*) cnt FROM emp e WHERE e.deptno = d.deptno) x"},
	{"lateral_star", "SELECT * FROM emp e, LATERAL (SELECT dname FROM dept d WHERE d.deptno = e.deptno)"},
	{"left_join_group", "SELECT d.dname, COUNT(e.empno) cnt FROM dept d LEFT JOIN emp e ON e.deptno = d.deptno GROUP BY d.dname"},

	// Subqueries.
	{"inline_view", "SELECT x.ename FROM (SELECT ename, sal FROM emp WHERE sal > 1000) x"},
	{"inline_view_unaliased_expr", "SELECT * FROM (SELECT ename, sal * 2 FROM emp)"},
	{"inline_view_aliased_expr", "SELECT * FROM (SELECT ename, sal * 2 dbl FROM emp)"},
	{"nested_star", "SELECT * FROM (SELECT * FROM dept)"},
	{"nested_star_qualified", "SELECT x.* FROM (SELECT d.* FROM dept d) x"},
	{"inline_view_union", "SELECT x FROM (SELECT 1 x FROM dual UNION ALL SELECT 2 FROM dual)"},
	{"in_subquery", "SELECT ename FROM emp WHERE deptno IN (SELECT deptno FROM dept WHERE loc = 'DALLAS')"},
	{"exists", "SELECT ename FROM emp e WHERE EXISTS (SELECT 1 FROM dept d WHERE d.deptno = e.deptno)"},
	{"not_exists", "SELECT d.deptno, d.dname FROM dept d WHERE NOT EXISTS (SELECT NULL FROM emp e WHERE e.deptno = d.deptno)"},
	{"scalar_subquery", "SELECT ename, (SELECT dname FROM dept d WHERE d.deptno = e.deptno) dname FROM emp e"},
	{"scalar_subquery_two_columns", "SELECT ename, (SELECT dname, loc FROM dept d WHERE d.deptno = e.deptno) x FROM emp e"},
	{"all_subquery", "SELECT ename FROM emp WHERE sal > ALL (SELECT sal FROM emp WHERE deptno = 30)"},
	{"any_list", "SELECT ename FROM emp WHERE deptno = ANY (10, 20)"},
	{"row_in_subquery", "SELECT ename FROM emp WHERE (deptno, job) IN (SELECT deptno, job FROM emp WHERE ename = 'SMITH')"},
	{"correlated", "SELECT ename, sal FROM emp e WHERE sal > (SELECT AVG(sal) FROM emp WHERE deptno = e.deptno)"},
	{"equals_subquery", "SELECT empno, ename FROM emp WHERE empno = (SELECT MAX(empno) FROM emp)"},
	{"equals_subquery_two_columns", "SELECT ename FROM emp WHERE empno = (SELECT empno, ename FROM emp)"},
	{"predicates", "SELECT ename FROM emp WHERE ename LIKE 'A%' AND sal BETWEEN 1000 AND 3000 AND deptno IN (10, 20) AND comm IS NOT NULL"},
	{"cursor_expression", "SELECT d.dname, CURSOR(SELECT ename FROM emp e WHERE e.deptno = d.deptno) c FROM dept d"},

	// WITH.
	{"cte", "WITH d AS (SELECT deptno, dname FROM dept) SELECT * FROM d"},
	{"cte_columns", "WITH d (no, name) AS (SELECT deptno, dname FROM dept) SELECT * FROM d"},
	{"cte_columns_mismatch", "WITH d (no) AS (SELECT deptno, dname FROM dept) SELECT * FROM d"},
	{"cte_chain", "WITH a AS (SELECT empno, mgr FROM emp), b AS (SELECT empno FROM a WHERE mgr IS NULL) SELECT * FROM b"},
	{"cte_join", "WITH d AS (SELECT * FROM dept) SELECT d.dname, e.ename FROM d JOIN emp e ON e.deptno = d.deptno"},
	{"recursive_cte", "WITH r (empno, ename, lvl) AS (SELECT empno, ename, 1 FROM emp WHERE mgr IS NULL UNION ALL SELECT e.empno, e.ename, r.lvl + 1 FROM emp e JOIN r ON e.mgr = r.empno) SELECT * FROM r"},
	{"recursive_cte_counter", "WITH r (n) AS (SELECT 1 FROM dual UNION ALL SELECT n + 1 FROM r WHERE n < 5) SELECT n FROM r"},
	{"recursive_cte_no_columns", "WITH r AS (SELECT 1 n FROM dual UNION ALL SELECT n + 1 FROM r WHERE n < 5) SELECT n FROM r"},
	{"recursive_cte_search_cycle", "WITH r (n) AS (SELECT 1 FROM dual UNION ALL SELECT n + 1 FROM r WHERE n < 5) SEARCH DEPTH FIRST BY n SET ord CYCLE n SET is_cycle TO 'Y' DEFAULT 'N' SELECT * FROM r"},
	{"cte_unused", "WITH d AS (SELECT deptno FROM dept) SELECT ename FROM emp"},
	{"cte_missing_column", "WITH d AS (SELECT deptno FROM dept) SELECT dname FROM d"},
	{"with_function", "WITH FUNCTION dbl(x NUMBER) RETURN NUMBER IS BEGIN RETURN x * 2; END; SELECT dbl(sal) d FROM emp"},

	// Set operators.
	{"union", "SELECT empno, ename FROM emp UNION SELECT deptno, dname FROM dept"},
	{"union_all", "SELECT empno FROM emp UNION ALL SELECT deptno FROM dept"},
	{"minus", "SELECT deptno FROM dept MINUS SELECT deptno FROM emp"},
	{"intersect", "SELECT deptno FROM dept INTERSECT SELECT deptno FROM emp"},
	{"except", "SELECT deptno FROM dept EXCEPT SELECT deptno FROM emp"},
	{"minus_all", "SELECT deptno FROM dept MINUS ALL SELECT deptno FROM emp"},
	{"intersect_all", "SELECT deptno FROM dept INTERSECT ALL SELECT deptno FROM emp"},
	{"union_column_count", "SELECT empno, ename FROM emp UNION SELECT deptno FROM dept"},
	{"union_type_mismatch", "SELECT empno FROM emp UNION SELECT dname FROM dept"},
	{"union_varchar_lengths", "SELECT ename FROM emp UNION SELECT dname FROM dept"},
	{"union_order_by_position", "SELECT job FROM emp UNION SELECT loc FROM dept ORDER BY 1"},
	{"union_order_by_name", "SELECT ename FROM emp UNION SELECT dname FROM dept ORDER BY ename"},
	{"union_order_by_second_name", "SELECT ename FROM emp UNION SELECT dname FROM dept ORDER BY dname"},
	{"union_parenthesized", "(SELECT ename FROM emp) UNION ALL (SELECT dname FROM dept)"},
	{"union_numbers", "SELECT sal FROM emp UNION ALL SELECT comm FROM emp"},
	{"union_number_scales", "SELECT sal FROM emp UNION SELECT deptno FROM dept"},
	{"union_char_varchar", "SELECT code FROM misc UNION SELECT title FROM misc"},
	{"union_null_first", "SELECT NULL x FROM dual UNION SELECT ename FROM emp"},
	{"union_date_timestamp", "SELECT hiredate FROM emp UNION SELECT ts FROM misc"},
	{"union_alias_first", "SELECT ename n FROM emp UNION SELECT dname FROM dept"},
	{"union_char_literals", "SELECT 'a' c FROM dual UNION SELECT 'bcd' FROM dual"},
	{"union_literal_column", "SELECT 'a' c FROM dual UNION ALL SELECT ename FROM emp"},
	{"union_number_literals", "SELECT 1 n FROM dual UNION SELECT 2.5 FROM dual"},
	{"union_star", "SELECT * FROM dept UNION SELECT * FROM dept"},
	{"union_same_char", "SELECT code FROM misc UNION SELECT code FROM misc"},
	{"union_clob", "SELECT note FROM misc UNION SELECT note FROM misc"},
	{"union_all_clob", "SELECT note FROM misc UNION ALL SELECT note FROM misc"},

	// Aggregates.
	{"group_by", "SELECT deptno, COUNT(*) cnt, SUM(sal) total, AVG(sal) avg_sal, MIN(hiredate) first_hire, MAX(ename) last_name FROM emp GROUP BY deptno"},
	{"min_max_number", "SELECT MAX(sal) mx, MIN(empno) mn, MAX(code) mc FROM misc, emp"},
	{"having", "SELECT deptno, COUNT(*) cnt FROM emp GROUP BY deptno HAVING COUNT(*) > 3"},
	{"not_group_by_expression", "SELECT deptno, ename, COUNT(*) cnt FROM emp GROUP BY deptno"},
	{"not_single_group", "SELECT ename, COUNT(*) cnt FROM emp"},
	{"aggregate_in_where", "SELECT deptno FROM emp WHERE SUM(sal) > 10 GROUP BY deptno"},
	{"rollup", "SELECT deptno, job, SUM(sal) s FROM emp GROUP BY ROLLUP (deptno, job)"},
	{"cube_grouping", "SELECT deptno, job, SUM(sal) s, GROUPING(job) g FROM emp GROUP BY CUBE (deptno, job)"},
	{"grouping_sets", "SELECT deptno, job, SUM(sal) s FROM emp GROUP BY GROUPING SETS ((deptno), (job), ())"},
	{"count_distinct", "SELECT COUNT(DISTINCT job) jobs FROM emp"},
	{"statistics", "SELECT deptno, STDDEV(sal) sd, VARIANCE(sal) v, MEDIAN(sal) med FROM emp GROUP BY deptno"},
	{"percentile_cont", "SELECT deptno, PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY sal) p50 FROM emp GROUP BY deptno"},
	{"aggregate_expression", "SELECT 1 + COUNT(*) c FROM emp"},
	{"group_by_alias", "SELECT deptno d, COUNT(*) c FROM emp GROUP BY d"},
	{"group_by_position", "SELECT deptno, COUNT(*) c FROM emp GROUP BY 1"},
	{"nested_aggregate", "SELECT MAX(AVG(sal)) m FROM emp GROUP BY deptno"},
	{"nested_aggregate_no_group", "SELECT MAX(AVG(sal)) m FROM emp"},
	{"nested_aggregate_with_column", "SELECT deptno, MAX(AVG(sal)) m FROM emp GROUP BY deptno"},
	{"listagg", "SELECT deptno, LISTAGG(ename, ',') WITHIN GROUP (ORDER BY ename) names FROM emp GROUP BY deptno"},
	{"keep_dense_rank", "SELECT deptno, MAX(sal) KEEP (DENSE_RANK LAST ORDER BY hiredate) last_sal FROM emp GROUP BY deptno"},
	{"aggregate_types", "SELECT SUM(bd) s, AVG(bf) a, MAX(ts) m, MIN(code) mc, COUNT(note) cn FROM misc"},
	{"group_by_clob", "SELECT COUNT(*) c FROM misc GROUP BY note"},
	{"distinct_clob", "SELECT DISTINCT note FROM misc"},
	{"order_by_clob", "SELECT note FROM misc ORDER BY note"},

	// Analytic functions.
	{"row_number", "SELECT ename, ROW_NUMBER() OVER (ORDER BY sal DESC) rn FROM emp"},
	{"rank", "SELECT ename, RANK() OVER (PARTITION BY deptno ORDER BY sal) rnk, DENSE_RANK() OVER (PARTITION BY deptno ORDER BY sal) drnk FROM emp"},
	{"sum_over_partition", "SELECT ename, SUM(sal) OVER (PARTITION BY deptno) dept_sal FROM emp"},
	{"moving_average", "SELECT ename, AVG(sal) OVER (ORDER BY hiredate ROWS BETWEEN 2 PRECEDING AND CURRENT ROW) mov FROM emp"},
	{"lag_lead", "SELECT ename, LAG(sal) OVER (ORDER BY hiredate) prev, LEAD(sal, 1, 0) OVER (ORDER BY hiredate) nxt FROM emp"},
	{"first_value", "SELECT ename, FIRST_VALUE(ename) OVER (PARTITION BY deptno ORDER BY sal) lowest FROM emp"},
	{"ntile", "SELECT ename, NTILE(4) OVER (ORDER BY sal) q FROM emp"},
	{"ratio_to_report", "SELECT ename, RATIO_TO_REPORT(sal) OVER () pct FROM emp"},
	{"count_over", "SELECT ename, COUNT(*) OVER () total FROM emp"},
	{"percent_rank", "SELECT ename, PERCENT_RANK() OVER (ORDER BY sal) pr, CUME_DIST() OVER (ORDER BY sal) cd FROM emp"},
	{"window_in_where", "SELECT ename, sal FROM emp WHERE ROW_NUMBER() OVER (ORDER BY sal) = 1"},
	{"window_clause", "SELECT ename, SUM(sal) OVER w s FROM emp WINDOW w AS (PARTITION BY deptno)"},
	{"window_arithmetic", "SELECT ename, sal, sal - AVG(sal) OVER (PARTITION BY deptno) diff FROM emp"},
	{"window_over_aggregate", "SELECT deptno, SUM(SUM(sal)) OVER () total FROM emp GROUP BY deptno"},
	{"window_order_by_clob", "SELECT ROW_NUMBER() OVER (ORDER BY note) r FROM misc"},

	// Hierarchical queries.
	{"connect_by_level", "SELECT LEVEL lvl, ename FROM emp CONNECT BY PRIOR empno = mgr START WITH mgr IS NULL"},
	{"connect_by_functions", "SELECT ename, SYS_CONNECT_BY_PATH(ename, '/') path, CONNECT_BY_ISLEAF leaf, CONNECT_BY_ROOT ename root FROM emp START WITH mgr IS NULL CONNECT BY NOCYCLE PRIOR empno = mgr"},
	{"connect_by_level_bare", "SELECT ename, LEVEL FROM emp CONNECT BY PRIOR empno = mgr ORDER SIBLINGS BY ename"},
	{"connect_by_dual", "SELECT LEVEL l FROM dual CONNECT BY LEVEL <= 10"},
	{"connect_by_start_with", "SELECT ename FROM emp START WITH mgr IS NULL CONNECT BY PRIOR empno = mgr"},
	{"connect_by_prior_column", "SELECT PRIOR ename p, ename FROM emp CONNECT BY PRIOR empno = mgr"},

	// Row limiting.
	{"fetch_first", "SELECT ename, sal FROM emp ORDER BY sal DESC FETCH FIRST 5 ROWS ONLY"},
	{"offset_fetch", "SELECT ename FROM emp ORDER BY sal OFFSET 2 ROWS FETCH NEXT 3 ROWS ONLY"},
	{"fetch_percent_ties", "SELECT ename FROM emp ORDER BY sal FETCH FIRST 10 PERCENT ROWS WITH TIES"},
	{"rownum_filter", "SELECT ename FROM emp WHERE ROWNUM <= 5"},
	{"rownum_order", "SELECT ename FROM emp WHERE ROWNUM < 3 ORDER BY ename"},

	// PIVOT and UNPIVOT.
	{"pivot", "SELECT * FROM (SELECT deptno, job, sal FROM emp) PIVOT (SUM(sal) FOR deptno IN (10 AS d10, 20 AS d20, 30 AS d30))"},
	{"pivot_two_aggregates", "SELECT * FROM (SELECT deptno, job, sal FROM emp) PIVOT (SUM(sal) s, COUNT(*) c FOR deptno IN (10, 20))"},
	{"pivot_multi_column", "SELECT * FROM (SELECT job, deptno, sal FROM emp) PIVOT (MAX(sal) FOR (job, deptno) IN (('CLERK', 10) AS c10, ('MANAGER', 20) AS m20))"},
	{"pivot_string_values", "SELECT * FROM (SELECT deptno, job, sal FROM emp) PIVOT (SUM(sal) FOR job IN ('CLERK', 'MANAGER'))"},
	{"pivot_table", "SELECT * FROM emp PIVOT (COUNT(*) FOR deptno IN (10 AS d10, 20 AS d20))"},
	{"unpivot", "SELECT * FROM (SELECT empno, sal, comm FROM emp) UNPIVOT (amount FOR kind IN (sal AS 'SAL', comm AS 'COMM'))"},
	{"unpivot_include_nulls", "SELECT * FROM (SELECT empno, sal, comm FROM emp) UNPIVOT INCLUDE NULLS (amount FOR kind IN (sal, comm))"},

	// CASE, DECODE, NVL and kin.
	{"case_searched", "SELECT CASE WHEN sal > 2000 THEN 'high' ELSE 'low' END band FROM emp"},
	{"case_simple", "SELECT CASE deptno WHEN 10 THEN 'ACCOUNTING' WHEN 20 THEN 'RESEARCH' END dname FROM emp"},
	{"case_numbers", "SELECT CASE WHEN sal > 2000 THEN sal ELSE comm END x FROM emp"},
	{"case_type_mismatch", "SELECT CASE WHEN sal > 2000 THEN 'x' ELSE 1 END x FROM emp"},
	{"case_char_varchar", "SELECT CASE WHEN id > 1 THEN code ELSE title END x FROM misc"},
	{"case_timestamp_null", "SELECT CASE WHEN id > 1 THEN ts ELSE NULL END x FROM misc"},
	{"decode_strings", "SELECT DECODE(deptno, 10, 'A', 20, 'BB', 'CCC') d FROM emp"},
	{"decode_numbers", "SELECT DECODE(deptno, 10, sal, comm) d FROM emp"},
	{"decode_number_results", "SELECT DECODE(job, 'CLERK', 1, 2) d FROM emp"},
	{"decode_null_first", "SELECT DECODE(id, 1, NULL, code) x FROM misc"},
	{"nvl_number", "SELECT NVL(comm, 0) c FROM emp"},
	{"nvl_string", "SELECT NVL(job, 'none') j FROM emp"},
	{"nvl_null_first", "SELECT NVL(NULL, 1) x FROM dual"},
	{"nvl_varchar_nvarchar", "SELECT NVL(title, name) x FROM misc"},
	{"nvl2", "SELECT NVL2(comm, 'Y', 'N') has_comm FROM emp"},
	{"coalesce", "SELECT COALESCE(comm, sal, 0) c FROM emp"},
	{"coalesce_null_first", "SELECT COALESCE(NULL, code) x FROM misc"},
	{"nullif", "SELECT NULLIF(sal, comm) n FROM emp"},
	{"greatest_least", "SELECT GREATEST(sal, comm) g, LEAST(ename, job) l FROM emp"},

	// Operators and built-in functions.
	{"concat_operator", "SELECT ename || ' ' || job label FROM emp"},
	{"concat_function", "SELECT CONCAT(ename, job) c FROM emp"},
	{"concat_number", "SELECT sal || 'x' x FROM emp"},
	{"arithmetic", "SELECT sal * 12 + NVL(comm, 0) annual FROM emp"},
	{"number_string_arithmetic", "SELECT sal + ename x FROM emp"},
	{"numeric_functions", "SELECT sal / 3 s, MOD(empno, 2) m, ROUND(sal, 1) r, TRUNC(sal) t, ABS(sal) a, CEIL(sal) c, FLOOR(sal) f, POWER(sal, 2) p, SQRT(sal) q FROM emp"},
	{"math_functions", "SELECT SIGN(sal) s, LN(sal) l, EXP(1) e, LOG(10, sal) lg FROM emp"},
	{"string_functions", "SELECT UPPER(ename) u, LOWER(ename) l, INITCAP(ename) i, LENGTH(ename) n, SUBSTR(ename, 1, 3) s, INSTR(ename, 'A') p FROM emp"},
	{"pad_trim_functions", "SELECT LPAD(ename, 15, '*') l, RPAD(job, 12) r, TRIM(ename) t, LTRIM(job) lt, RTRIM(job) rt, REPLACE(ename, 'A', 'a') rp, TRANSLATE(job, 'AB', 'ab') tr FROM emp"},
	{"regexp_functions", "SELECT REGEXP_SUBSTR(ename, '[A-Z]+') rs, REGEXP_REPLACE(job, 'A', 'B') rr, REGEXP_INSTR(ename, 'A') ri, REGEXP_COUNT(ename, 'A') rc FROM emp"},
	{"chr_ascii", "SELECT CHR(65) c, ASCII(ename) a FROM emp"},
	{"dump_vsize", "SELECT DUMP(ename) d, VSIZE(ename) v FROM emp"},
	{"to_char", "SELECT TO_CHAR(hiredate, 'YYYY-MM-DD') d, TO_CHAR(sal) s, TO_CHAR(sal, '9999.99') f FROM emp"},
	{"conversions", "SELECT TO_NUMBER('42') n, TO_DATE('2024-01-01', 'YYYY-MM-DD') d, TO_TIMESTAMP('2024-01-01', 'YYYY-MM-DD') t FROM dual"},
	{"current_time", "SELECT SYSDATE now, SYSTIMESTAMP ts, CURRENT_DATE cd, CURRENT_TIMESTAMP ct, LOCALTIMESTAMP lt FROM dual"},
	{"session_functions", "SELECT USER u, UID i, SYS_GUID() g, USERENV('LANG') l FROM dual"},
	{"date_functions", "SELECT ADD_MONTHS(hiredate, 1) a, MONTHS_BETWEEN(SYSDATE, hiredate) m, LAST_DAY(hiredate) l, NEXT_DAY(hiredate, 'MONDAY') n, TRUNC(hiredate, 'MM') t, EXTRACT(YEAR FROM hiredate) y FROM emp"},
	{"date_arithmetic", "SELECT hiredate + 1 d, hiredate - hiredate n, SYSDATE - hiredate days FROM emp"},
	{"date_number_compare", "SELECT ename FROM emp WHERE hiredate > 5"},
	{"cast", "SELECT CAST(sal AS VARCHAR2(20)) s, CAST(ename AS CHAR(12)) c, CAST(empno AS NUMBER(10,2)) n, CAST(hiredate AS TIMESTAMP) t FROM emp"},
	{"cast_literals", "SELECT CAST(NULL AS DATE) d, CAST('1' AS INTEGER) i, CAST(1 AS BINARY_DOUBLE) b FROM dual"},
	{"literals", "SELECT 'abc' s, 1 n, 1.5 f, N'x' ns, DATE '2024-01-01' d, TIMESTAMP '2024-01-01 00:00:00' t, NULL z FROM dual"},
	{"boolean_literal", "SELECT TRUE b FROM dual"},
	{"misc_expressions", "SELECT UPPER(code) a, code || title b, LENGTH(note) c, SUBSTR(note, 1, 10) d, TO_CHAR(ts) e, bd + 1 f, bf * 2 g, ts + 1 h, ts - ts i, ds * 2 j FROM misc"},
	{"lob_raw_conversions", "SELECT TO_CLOB(title) c, TO_NCHAR(title) nc, RAWTOHEX(r) h, HEXTORAW('AB') rr FROM misc"},
	{"json_object", "SELECT JSON_OBJECT('name' VALUE ename) j FROM emp"},
	{"json_value", `SELECT JSON_VALUE('{"a":1}', '$.a') j FROM dual`},
	{"json_arrayagg", "SELECT JSON_ARRAYAGG(ename) j FROM emp"},
	{"xmlelement", `SELECT XMLELEMENT("e", ename) x FROM emp`},
	{"json_table", `SELECT empno, xt.c FROM emp, JSON_TABLE('{"a":[1,2]}', '$.a[*]' COLUMNS (c NUMBER PATH '$')) xt`},
	{"table_collection", "SELECT column_value cv FROM TABLE(sys.odcinumberlist(1, 2, 3))"},
	{"table_collection_star", "SELECT * FROM TABLE(sys.odcivarchar2list('a', 'b'))"},
	{"model", "SELECT ename, sal FROM emp MODEL DIMENSION BY (empno) MEASURES (sal, ename) RULES (sal[7369] = 0)"},
	{"match_recognize", "SELECT * FROM emp MATCH_RECOGNIZE (ORDER BY hiredate MEASURES FIRST(ename) AS first_name PATTERN (a b*) DEFINE b AS sal > PREV(sal))"},
	{"flashback_query", "SELECT ename FROM emp AS OF TIMESTAMP SYSDATE - 1"},
}
