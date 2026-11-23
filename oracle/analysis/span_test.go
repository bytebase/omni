package analysis

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/sijms/go-ora/v2/network"

	"github.com/bytebase/omni/enginetest/span"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/internal/testcontainer"
)

// TestQuerySpanFidelity holds ResolveQuerySpan to Oracle: for each query of
// testdata/span/queries.sql, the output column names, each column's lineage
// and the tables read, as omni resolves them against the schema Oracle
// stored for the fixture, and as Oracle states them. testdata/README.md says
// where each of Oracle's answers comes from.
func TestQuerySpanFidelity(t *testing.T) {
	queries, err := readSpanCorpus(filepath.Join("testdata", "span", "queries.sql"))
	if err != nil {
		t.Fatal(err)
	}
	// The schema Oracle stored for each fixture, which omni resolves against.
	stored := map[string]*spanFixture{}
	span.Check(t, span.Suite{
		Path:     filepath.Join("testdata", "span.tsv"),
		Pins:     testcontainer.Pins,
		Fixtures: filepath.Join("testdata", "span", "fixtures"),
		Queries:  queries,
		Omni: func(fx span.Fixture, session, statement string) span.Answer {
			return omniSpan(stored[fx.Name], session, statement)
		},
		Engine: func(t testing.TB) func(span.Fixture, string, string) (span.Answer, error) {
			e := &spanEngine{ctx: context.Background(), db: testcontainer.Open(t, "system")}
			built := ""
			return func(fx span.Fixture, session, statement string) (span.Answer, error) {
				if built != fx.Name {
					if built != "" {
						return span.Answer{}, fmt.Errorf("fixture %s: the test builds one fixture, and %s is built", fx.Name, built)
					}
					s, err := e.fixture(fx.SQL)
					if err != nil {
						return span.Answer{}, fmt.Errorf("fixture %s: %w", fx.Name, err)
					}
					stored[fx.Name], built = s, fx.Name
				}
				return e.span(stored[fx.Name], session, statement)
			}
		},
	})
}

// spanInstance is the instance the fixture's schemas are on, to omni's
// catalog.
const spanInstance = "oracle"

// spanLink is the fixture's public database link, a loopback to PEOPLE's schema
// on the same server.
const (
	spanLink     = "OMNI_LOOP"
	spanLinkUser = "PEOPLE"
)

// omniSpan is ResolveQuerySpan's answer, resolved as Bytebase resolves it for
// Oracle: the session's current schema is DefaultDatabase, against the
// metadata a sync stores, with names matched exactly.
func omniSpan(fx *spanFixture, session, statement string) span.Answer {
	resolved, err := ResolveQuerySpan(statement, ResolveOptions{
		Instance:        spanInstance,
		DefaultDatabase: session,
		Catalog:         fx.catalog(),
	})
	switch {
	case err != nil:
		return span.Answer{Error: "omni"}
	case resolved == nil:
		return span.Answer{Error: "omni: no span"}
	case resolved.NotFoundError != nil:
		return span.Answer{Error: omniRefusal(resolved.NotFoundError)}
	}
	var a span.Answer
	tables := map[string]bool{}
	for c := range resolved.SourceColumns {
		tables[spanName(c.Database, c.Table)] = true
	}
	for name := range tables {
		a.Tables = append(a.Tables, name)
	}
	for _, r := range resolved.Results {
		a.Results = append(a.Results, r.Name)
		var sources []string
		for c := range r.SourceColumns {
			sources = append(sources, spanName(c.Database, c.Table, c.Column))
		}
		a.Lineage = append(a.Lineage, sources)
	}
	return a
}

// omniRefusal names a missing resource by the number of the error Oracle
// gives the same refusal: ORA-00904 for a column, ORA-00942 for a table or a
// schema.
func omniRefusal(e *NotFoundError) string {
	switch {
	case e.Column != nil:
		return "904"
	case e.Table != nil || e.Database != nil:
		return "942"
	}
	return "omni: not found"
}

// spanName joins the parts of a name with dots.
func spanName(parts ...string) string { return strings.Join(parts, ".") }

// spanFixture is what Oracle stored for a fixture.
type spanFixture struct {
	// Schemas is each schema the fixture built, as Bytebase's sync stores
	// it: one database per schema, its objects in the unnamed schema.
	Schemas map[string]*metadata.DatabaseSchemaMetadata
	// Domains maps each domain the test gave a table column, owner.domain,
	// to that column, owner.table.column.
	Domains map[string]string
}

// catalog is the fixture's metadata as omni's Catalog. The link resolves to
// the schema it names, or to its user's.
func (fx *spanFixture) catalog() Catalog {
	database := func(name string) Database {
		if d, ok := fx.Schemas[name]; ok {
			return newSpanDatabase(d)
		}
		return nil
	}
	return Catalog{
		GetDatabase: func(instance, name string) (Database, error) {
			if instance != spanInstance {
				return nil, nil
			}
			return database(name), nil
		},
		GetLinkedDatabase: func(_, link, schema string) (string, string, Database, error) {
			if !strings.EqualFold(link, spanLink) {
				return "", "", nil, nil
			}
			if schema == "" {
				schema = spanLinkUser
			}
			return spanInstance, schema, database(schema), nil
		},
	}
}

// spanDatabase is one schema's metadata, looked up by exact name.
type spanDatabase struct {
	tables            map[string]*metadata.TableMetadata
	views             map[string]*metadata.ViewMetadata
	materializedViews map[string]*metadata.MaterializedViewMetadata
	synonyms, public  map[string]*metadata.SynonymMetadata
}

func newSpanDatabase(d *metadata.DatabaseSchemaMetadata) *spanDatabase {
	db := &spanDatabase{
		tables:            map[string]*metadata.TableMetadata{},
		views:             map[string]*metadata.ViewMetadata{},
		materializedViews: map[string]*metadata.MaterializedViewMetadata{},
		synonyms:          map[string]*metadata.SynonymMetadata{},
		public:            map[string]*metadata.SynonymMetadata{},
	}
	for _, s := range d.PublicSynonyms {
		db.public[s.Name] = s
	}
	for _, s := range d.Schemas {
		for _, syn := range s.Synonyms {
			db.synonyms[syn.Name] = syn
		}
		for _, t := range s.Tables {
			db.tables[t.Name] = t
		}
		for _, v := range s.Views {
			db.views[v.Name] = v
		}
		for _, v := range s.MaterializedViews {
			db.materializedViews[v.Name] = v
		}
	}
	return db
}

func (*spanDatabase) HasSchema() bool                                { return true }
func (d *spanDatabase) GetTable(name string) *metadata.TableMetadata { return d.tables[name] }
func (d *spanDatabase) GetView(name string) *metadata.ViewMetadata   { return d.views[name] }
func (d *spanDatabase) GetMaterializedView(name string) *metadata.MaterializedViewMetadata {
	return d.materializedViews[name]
}
func (*spanDatabase) GetExternalTable(string) *metadata.ExternalTableMetadata  { return nil }
func (d *spanDatabase) GetSynonym(name string) *metadata.SynonymMetadata       { return d.synonyms[name] }
func (d *spanDatabase) GetPublicSynonym(name string) *metadata.SynonymMetadata { return d.public[name] }

// spanEngine asks Oracle about queries, as SYSTEM, on one session of t's tape.
type spanEngine struct {
	ctx context.Context
	db  *sql.DB
}

func (e *spanEngine) exec(stmt string) error {
	_, err := e.db.ExecContext(e.ctx, stmt)
	return err
}

// query runs stmt and returns its rows, NULL as nil.
func (e *spanEngine) query(stmt string) ([][]*string, error) {
	rows, err := e.db.QueryContext(e.ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var out [][]*string
	for rows.Next() {
		cells := make([]sql.NullString, len(cols))
		ptrs := make([]any, len(cols))
		for i := range cells {
			ptrs[i] = &cells[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		row := make([]*string, len(cols))
		for i, c := range cells {
			if c.Valid {
				row[i] = &c.String
			}
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

// literal quotes s as a SQL string.
func literal(s string) string { return "'" + strings.ReplaceAll(s, "'", "''") + "'" }

// quote quotes s as an Oracle identifier.
func quote(s string) string { return `"` + strings.ReplaceAll(s, `"`, `""`) + `"` }

// fixtureStatements splits a fixture into its statements: each ends with a
// semicolon at the end of a line, and a line starting with -- is a comment.
func fixtureStatements(text string) []string {
	var (
		out []string
		cur strings.Builder
	)
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "--") {
			continue
		}
		cur.WriteString(line + "\n")
		if strings.HasSuffix(strings.TrimSpace(line), ";") {
			stmt := strings.TrimSuffix(strings.TrimSpace(cur.String()), ";")
			out = append(out, stmt)
			cur.Reset()
		}
	}
	if s := strings.TrimSpace(cur.String()); s != "" {
		out = append(out, s)
	}
	return out
}

// fixture builds the fixture as SYSTEM and reads what Oracle stored. It then
// gives every table column a domain of its own, which Oracle carries to a
// view column that reads the column (span). The fixture grants each schema
// SELECT ANY TABLE, which the view over a query needs.
func (e *spanEngine) fixture(text string) (*spanFixture, error) {
	for _, stmt := range fixtureStatements(text) {
		if err := e.exec(stmt); err != nil {
			return nil, fmt.Errorf("%s: %w", stmt, err)
		}
	}
	users, err := e.query("SELECT username FROM all_users WHERE oracle_maintained = 'N' ORDER BY username")
	if err != nil {
		return nil, err
	}
	fx := &spanFixture{Schemas: map[string]*metadata.DatabaseSchemaMetadata{}, Domains: map[string]string{}}
	for _, u := range users {
		owner := *u[0]
		d, err := e.schema(owner)
		if err != nil {
			return nil, fmt.Errorf("schema %s: %w", owner, err)
		}
		fx.Schemas[owner] = d
	}
	// Bytebase's sync stores every public synonym with each database; the
	// fixture's are those that name its schemas, and the image's thousands
	// name none.
	var owners []string
	for _, u := range users {
		owners = append(owners, literal(*u[0]))
	}
	public, err := e.synonyms("owner = 'PUBLIC' AND table_owner IN (" + strings.Join(owners, ", ") + ")")
	if err != nil {
		return nil, err
	}
	for _, d := range fx.Schemas {
		d.PublicSynonyms = public
	}
	for _, u := range users {
		owner := *u[0]
		if err := e.domains(owner, fx); err != nil {
			return nil, fmt.Errorf("domains of %s: %w", owner, err)
		}
	}
	for _, u := range users {
		if err := e.exec("BEGIN DBMS_UTILITY.COMPILE_SCHEMA(" + literal(*u[0]) + ", FALSE); END;"); err != nil {
			return nil, err
		}
	}
	return fx, nil
}

// schema reads one schema as Bytebase's sync stores it (oracle/driver's
// sync.go): its tables, but the containers of materialized views, with the
// columns ALL_TAB_COLS lists as user-generated, in COLUMN_ID order with the
// invisible ones last; its views with their TEXT; its materialized views
// with their QUERY; and its synonyms.
func (e *spanEngine) schema(owner string) (*metadata.DatabaseSchemaMetadata, error) {
	s := &metadata.SchemaMetadata{}
	rows, err := e.query(`SELECT c.table_name, c.column_name, c.hidden_column
FROM all_tab_cols c JOIN all_tables t ON t.owner = c.owner AND t.table_name = c.table_name
WHERE c.owner = ` + literal(owner) + ` AND c.user_generated = 'YES'
AND NOT EXISTS (SELECT 1 FROM all_mviews m WHERE m.owner = c.owner AND m.mview_name = c.table_name)
ORDER BY c.table_name, c.column_id NULLS LAST, c.column_name`)
	if err != nil {
		return nil, err
	}
	var table *metadata.TableMetadata
	for _, r := range rows {
		if table == nil || table.Name != *r[0] {
			table = &metadata.TableMetadata{Name: *r[0]}
			s.Tables = append(s.Tables, table)
		}
		table.Columns = append(table.Columns, &metadata.ColumnMetadata{Name: *r[1], Nullable: true, IsInvisible: *r[2] == "YES"})
	}
	rows, err = e.query("SELECT view_name, text FROM all_views WHERE owner = " + literal(owner) + " ORDER BY view_name")
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		s.Views = append(s.Views, &metadata.ViewMetadata{Name: *r[0], Definition: strings.TrimRight(*r[1], "\x00")})
	}
	rows, err = e.query("SELECT mview_name, query FROM all_mviews WHERE owner = " + literal(owner) + " ORDER BY mview_name")
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		s.MaterializedViews = append(s.MaterializedViews, &metadata.MaterializedViewMetadata{Name: *r[0], Definition: strings.TrimRight(*r[1], "\x00")})
	}
	if s.Synonyms, err = e.synonyms("owner = " + literal(owner)); err != nil {
		return nil, err
	}
	return &metadata.DatabaseSchemaMetadata{Name: owner, Schemas: []*metadata.SchemaMetadata{s}}, nil
}

// synonyms reads the synonyms where filter holds, as oracle/driver's sync
// reads them. The sync reads ALL_SYNONYMS as the schema's owner, who sees
// every synonym of the schema; SYSTEM sees one through a database link only
// in DBA_SYNONYMS.
func (e *spanEngine) synonyms(filter string) ([]*metadata.SynonymMetadata, error) {
	rows, err := e.query("SELECT synonym_name, table_owner, table_name, db_link FROM dba_synonyms WHERE " + filter + " ORDER BY synonym_name")
	if err != nil {
		return nil, err
	}
	var out []*metadata.SynonymMetadata
	for _, r := range rows {
		syn := &metadata.SynonymMetadata{Name: *r[0]}
		if r[1] != nil {
			syn.TargetSchema = *r[1]
		}
		if r[2] != nil {
			syn.TargetName = *r[2]
		}
		if r[3] != nil {
			syn.TargetDbLink = *r[3]
		}
		out = append(out, syn)
	}
	return out, nil
}

// domainPrefix names the domains the test makes.
const domainPrefix = "OMNI_SPAN_D"

// domains gives each column of owner's tables a domain of its own, of the
// column's type, and records which column each domain marks.
func (e *spanEngine) domains(owner string, fx *spanFixture) error {
	rows, err := e.query(`SELECT c.table_name, c.column_name, c.data_type, c.data_precision, c.data_scale, c.char_length, c.char_used
FROM all_tab_cols c JOIN all_tables t ON t.owner = c.owner AND t.table_name = c.table_name
WHERE c.owner = ` + literal(owner) + ` AND c.user_generated = 'YES'
AND NOT EXISTS (SELECT 1 FROM all_mviews m WHERE m.owner = c.owner AND m.mview_name = c.table_name)
ORDER BY c.table_name, c.column_name`)
	if err != nil {
		return err
	}
	for _, r := range rows {
		name := fmt.Sprintf("%s%d", domainPrefix, len(fx.Domains)+1)
		domain := quote(owner) + "." + quote(name)
		if err := e.exec("CREATE DOMAIN " + domain + " AS " + columnType(r[2:])); err != nil {
			return err
		}
		if err := e.exec("ALTER TABLE " + quote(owner) + "." + quote(*r[0]) + " MODIFY (" + quote(*r[1]) + " DOMAIN " + domain + ")"); err != nil {
			return err
		}
		fx.Domains[spanName(owner, name)] = spanName(owner, *r[0], *r[1])
	}
	return nil
}

// columnType renders a column's type from ALL_TAB_COLS's data_type,
// data_precision, data_scale, char_length and char_used.
func columnType(c []*string) string {
	switch *c[0] {
	case "NUMBER":
		if c[1] == nil {
			return "NUMBER"
		}
		scale := "0"
		if c[2] != nil {
			scale = *c[2]
		}
		return fmt.Sprintf("NUMBER(%s,%s)", *c[1], scale)
	case "VARCHAR2", "CHAR", "NVARCHAR2", "NCHAR":
		unit := ""
		if *c[0] == "VARCHAR2" || *c[0] == "CHAR" {
			unit = " BYTE"
			if c[4] != nil && *c[4] == "C" {
				unit = " CHAR"
			}
		}
		return fmt.Sprintf("%s(%s%s)", *c[0], *c[3], unit)
	}
	return *c[0]
}

// probeView is the view span creates over a query in the session's schema,
// so that unqualified names in it resolve there.
const probeView = "OMNI_SPAN_PROBE"

// planID marks the rows EXPLAIN PLAN writes for a query.
const planID = "omni_span"

// span asks Oracle about one query, with the session's current schema set:
//
//   - its output column names, from the cursor's describe information;
//   - each column's source, from the domain Oracle gives the column of a
//     view over the query: a column that reads a table column carries its
//     domain, and the test gave each table column a domain of its own;
//   - the tables and views it reads: the objects of its EXPLAIN PLAN, and
//     what the view over it depends on, through synonyms and views.
func (e *spanEngine) span(fx *spanFixture, session, statement string) (span.Answer, error) {
	if err := e.exec("ALTER SESSION SET CURRENT_SCHEMA = " + quote(session)); err != nil {
		return span.Answer{}, err
	}
	defer e.exec("ALTER SESSION SET CURRENT_SCHEMA = SYSTEM")

	rows, err := e.db.QueryContext(e.ctx, statement)
	if err != nil {
		var oe *network.OracleError
		if errors.As(err, &oe) {
			return span.Answer{Error: strconv.Itoa(oe.ErrCode)}, nil
		}
		return span.Answer{}, err
	}
	names, err := rows.Columns()
	rows.Close()
	if err != nil {
		return span.Answer{}, err
	}
	a := span.Answer{Results: names}

	view := quote(session) + "." + probeView
	cols := make([]string, len(names))
	for i := range names {
		cols[i] = fmt.Sprintf("C%d", i+1)
	}
	if err := e.exec("CREATE VIEW " + view + " (" + strings.Join(cols, ", ") + ") AS " + statement); err != nil {
		return span.Answer{}, fmt.Errorf("create the probe view: %w", err)
	}
	domains, err := e.query("SELECT domain_owner, domain_name FROM all_tab_cols WHERE owner = " + literal(session) + " AND table_name = '" + probeView + "' ORDER BY column_id")
	if err != nil {
		return span.Answer{}, err
	}
	for _, d := range domains {
		sources := []string{span.Unstated}
		if d[0] != nil && d[1] != nil {
			column, ok := fx.Domains[spanName(*d[0], *d[1])]
			if !ok {
				return span.Answer{}, fmt.Errorf("a column of domain %s.%s, which the test did not make", *d[0], *d[1])
			}
			sources = append(sources, column)
		}
		a.Lineage = append(a.Lineage, sources)
	}
	tables := map[string]bool{}
	if err := e.dependencies(session, probeView, tables); err != nil {
		return span.Answer{}, err
	}
	if err := e.exec("DROP VIEW " + view); err != nil {
		return span.Answer{}, err
	}
	if err := e.plan(statement, tables); err != nil {
		return span.Answer{}, err
	}
	for t := range tables {
		a.Tables = append(a.Tables, t)
	}
	return a, nil
}

// dependencies adds the tables, views and materialized views the view
// owner.name depends on to into, reading through synonyms and, in turn,
// through the views it names. An object through a database link is another
// database's, which the test cannot name: it adds span.Unstated.
func (e *spanEngine) dependencies(owner, name string, into map[string]bool) error {
	rows, err := e.query(`SELECT referenced_owner, referenced_name, referenced_type, referenced_link_name
FROM all_dependencies WHERE owner = ` + literal(owner) + ` AND name = ` + literal(name) + ` AND type = 'VIEW'
ORDER BY 1, 2, 3`)
	if err != nil {
		return err
	}
	for _, r := range rows {
		if r[3] != nil {
			into[span.Unstated] = true
			continue
		}
		if err := e.read(*r[0], *r[1], *r[2], into); err != nil {
			return err
		}
	}
	return nil
}

// read adds the object owner.name of kind to into, as dependencies does.
func (e *spanEngine) read(owner, name, kind string, into map[string]bool) error {
	switch kind {
	case "TABLE", "MATERIALIZED VIEW":
		into[spanName(owner, name)] = true
	case "VIEW":
		if into[spanName(owner, name)] {
			return nil
		}
		into[spanName(owner, name)] = true
		return e.dependencies(owner, name, into)
	case "SYNONYM":
		rows, err := e.query("SELECT table_owner, table_name, db_link FROM dba_synonyms WHERE owner = " + literal(owner) + " AND synonym_name = " + literal(name))
		if err != nil {
			return err
		}
		if len(rows) != 1 {
			return fmt.Errorf("synonym %s.%s: %d rows", owner, name, len(rows))
		}
		if rows[0][2] != nil {
			into[span.Unstated] = true
			return nil
		}
		target, err := e.query("SELECT object_type FROM all_objects WHERE owner = " + literal(*rows[0][0]) + " AND object_name = " + literal(*rows[0][1]) + " AND object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW', 'SYNONYM') ORDER BY object_type")
		if err != nil {
			return err
		}
		if len(target) == 0 {
			return fmt.Errorf("synonym %s.%s names no table or view", owner, name)
		}
		// A materialized view is a table too; read it as the materialized view.
		return e.read(*rows[0][0], *rows[0][1], *target[0][0], into)
	}
	return nil
}

// plan adds the objects of statement's EXPLAIN PLAN to into: the tables,
// views and materialized views its operations read, and an index's table. A
// remote operation reads another database's objects, which the test cannot
// name: it adds span.Unstated. The plan also names objects of its own, which
// no dictionary view lists: the temporary table of a materialized WITH
// element, and the views of a transformation (VW_FOJ_0 for a FULL JOIN,
// VW_DCL_… for a decorrelated lateral view). They are not read.
func (e *spanEngine) plan(statement string, into map[string]bool) error {
	if err := e.exec("EXPLAIN PLAN SET STATEMENT_ID = '" + planID + "' FOR " + statement); err != nil {
		return fmt.Errorf("explain: %w", err)
	}
	rows, err := e.query(`SELECT p.operation, p.options, p.object_owner, p.object_name, p.object_type, p.object_node, i.table_owner, i.table_name,
	(SELECT COUNT(*) FROM all_objects o WHERE o.owner = p.object_owner AND o.object_name = p.object_name AND o.object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW'))
FROM plan_table p LEFT JOIN all_indexes i ON i.owner = p.object_owner AND i.index_name = p.object_name
WHERE p.statement_id = '` + planID + `' ORDER BY p.id`)
	if err != nil {
		return err
	}
	if err := e.exec("DELETE FROM plan_table WHERE statement_id = '" + planID + "'"); err != nil {
		return err
	}
	for _, r := range rows {
		op, options := *r[0], ""
		if r[1] != nil {
			options = *r[1]
		}
		if op == "REMOTE" || options == "REMOTE" || r[5] != nil {
			into[span.Unstated] = true
			continue
		}
		if r[2] == nil || r[3] == nil || r[4] == nil {
			continue
		}
		kind, listed := *r[4], r[8] != nil && *r[8] != "0"
		switch {
		case strings.HasPrefix(kind, "INDEX") && r[6] != nil && r[7] != nil:
			into[spanName(*r[6], *r[7])] = true
		case listed && (kind == "TABLE" || kind == "VIEW" || kind == "MAT_VIEW"):
			into[spanName(*r[2], *r[3])] = true
		}
	}
	return nil
}

// readSpanCorpus reads queries.sql: one query per line, in the session and
// on the fixture the last "-- session:" and "-- fixture:" lines name. Other
// lines starting "--" are comments.
func readSpanCorpus(path string) ([]span.Query, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var (
		out              []span.Query
		fixture, session string
	)
	sc := bufio.NewScanner(f)
	for n := 1; sc.Scan(); n++ {
		line := strings.TrimSpace(sc.Text())
		switch {
		case line == "":
		case strings.HasPrefix(line, "-- fixture:"):
			fixture = strings.TrimSpace(strings.TrimPrefix(line, "-- fixture:"))
		case strings.HasPrefix(line, "-- session:"):
			session = strings.TrimSpace(strings.TrimPrefix(line, "-- session:"))
		case strings.HasPrefix(line, "--"):
		case fixture == "" || session == "":
			return nil, fmt.Errorf("%s:%d: a query before the -- fixture: and -- session: lines", path, n)
		default:
			out = append(out, span.Query{Fixture: fixture, Session: session, Statement: line})
		}
	}
	return out, sc.Err()
}
