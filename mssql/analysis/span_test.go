package analysis

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/span"
	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/internal/testcontainer"
)

// TestQuerySpanFidelity holds ResolveQuerySpan to SQL Server: for each query
// of testdata/span/queries.sql, the output column names, each column's
// lineage and the tables read, as omni resolves them against the schema SQL
// Server stored for the fixture, and as SQL Server states them.
// testdata/README.md says where each of SQL Server's answers comes from.
func TestQuerySpanFidelity(t *testing.T) {
	queries, err := readSpanCorpus(filepath.Join("testdata", "span", "queries.sql"))
	if err != nil {
		t.Fatal(err)
	}
	// The schema SQL Server stored for each fixture, which omni resolves
	// against.
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
			ctx := context.Background()
			conn, err := testcontainer.Open(t).Conn(ctx)
			if err != nil {
				t.Fatalf("connect to SQL Server: %v", err)
			}
			t.Cleanup(func() { conn.Close() })
			e := &spanEngine{ctx: ctx, conn: conn, home: testcontainer.Database(t)}
			built := ""
			return func(fx span.Fixture, session, statement string) (span.Answer, error) {
				if built != fx.Name {
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

// omniSpan is ResolveQuerySpan's answer, resolved as Bytebase resolves it for
// SQL Server: the session's database is DefaultDatabase and its schema
// DefaultSchema, against the metadata a sync stores, case-insensitively
// (Bytebase's store.IsObjectCaseSensitive is false for SQL Server).
func omniSpan(fx *spanFixture, session, statement string) span.Answer {
	database, schema, _ := strings.Cut(session, ".")
	resolved, err := ResolveQuerySpan(statement, ResolveOptions{
		DefaultDatabase:     database,
		DefaultSchema:       schema,
		IgnoreCaseSensitive: true,
		Catalog:             fx.catalog(),
	})
	if err != nil || resolved.NotFoundError != nil {
		return span.Answer{Error: omniRefusal(err, resolved)}
	}
	var a span.Answer
	tables := map[string]bool{}
	for c := range resolved.SourceColumns {
		tables[spanName(c.Database, c.Schema, c.Table)] = true
	}
	for name := range tables {
		a.Tables = append(a.Tables, name)
	}
	for _, r := range resolved.Results {
		a.Results = append(a.Results, r.Name)
		var sources []string
		for c := range r.SourceColumns {
			sources = append(sources, spanName(c.Database, c.Schema, c.Table, c.Column))
		}
		a.Lineage = append(a.Lineage, sources)
	}
	return a
}

// omniRefusal names why omni gives no span: the number of the error SQL
// Server reports for the same refusal, where omni's error is one, and omni
// otherwise.
func omniRefusal(err error, resolved *ResolvedSpan) string {
	var (
		unknown   *unknownColumnError
		ambiguous *ambiguousColumnError
	)
	switch {
	case errors.As(err, &unknown) && unknown.table == "":
		return "207"
	case errors.As(err, &unknown):
		return "4104"
	case errors.As(err, &ambiguous):
		return "209"
	case err == nil && resolved.NotFoundError.Column == nil:
		return "208"
	}
	return "omni"
}

// cutLast is strings.Cut at the last sep.
func cutLast(s, sep string) (before, after string, found bool) {
	if i := strings.LastIndex(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}

// spanName joins the parts of a name with dots.
func spanName(parts ...string) string { return strings.Join(parts, ".") }

// spanFixture is what SQL Server stored for a fixture.
type spanFixture struct {
	// Databases is each database the fixture built, as Bytebase's sync
	// stores it: schemas, tables and views with their columns in order, and
	// each view's definition.
	Databases []*metadatapb.DatabaseSchemaMetadata
	// Views maps each view, database.schema.view, to the tables and views
	// its definition names, as sys.dm_sql_referenced_entities states them.
	Views map[string][]string
}

// catalog is the fixture's metadata as omni's Catalog, looked up as
// Bytebase's model looks it up for SQL Server: case-insensitively.
func (fx *spanFixture) catalog() Catalog {
	return Catalog{
		ListDatabaseNames: func() ([]string, error) {
			var names []string
			for _, d := range fx.Databases {
				names = append(names, d.Name)
			}
			return names, nil
		},
		GetDatabase: func(name string) (Database, error) {
			for _, d := range fx.Databases {
				if !strings.EqualFold(d.Name, name) {
					continue
				}
				db := &testDatabase{schemas: map[string]*metadatapb.SchemaMetadata{}}
				for _, s := range d.Schemas {
					db.schemas[s.Name] = s
					db.order = append(db.order, s.Name)
				}
				return db, nil
			}
			return nil, nil
		},
	}
}

// spanEngine asks SQL Server about queries, on one session of t's tape.
type spanEngine struct {
	ctx  context.Context
	conn *sql.Conn
	home string // the database of the test, which the session starts in
}

func (e *spanEngine) exec(stmt string) error {
	_, err := e.conn.ExecContext(e.ctx, stmt)
	return err
}

// query runs stmt and returns its rows, NULL as nil.
func (e *spanEngine) query(stmt string) ([][]*string, error) {
	rows, err := e.conn.QueryContext(e.ctx, stmt)
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

// literal quotes s as a T-SQL Unicode string.
func literal(s string) string { return "N'" + strings.ReplaceAll(s, "'", "''") + "'" }

// bracket quotes a name as a T-SQL bracketed identifier.
func bracket(s string) string { return "[" + strings.ReplaceAll(s, "]", "]]") + "]" }

var goLine = regexp.MustCompile(`(?mi)^\s*GO\s*$`)

// fixture drops every database but the test's own, builds the fixture, its
// batches separated by GO lines, and reads what SQL Server stored.
func (e *spanEngine) fixture(sqlText string) (*spanFixture, error) {
	if err := e.exec("USE " + bracket(e.home)); err != nil {
		return nil, err
	}
	others, err := e.databases()
	if err != nil {
		return nil, err
	}
	for _, d := range others {
		if err := e.exec("ALTER DATABASE " + bracket(d) + " SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE " + bracket(d)); err != nil {
			return nil, err
		}
	}
	for _, batch := range goLine.Split(sqlText, -1) {
		if strings.TrimSpace(batch) == "" {
			continue
		}
		if err := e.exec(batch); err != nil {
			return nil, fmt.Errorf("%s: %w", strings.TrimSpace(batch), err)
		}
	}
	if err := e.exec("USE " + bracket(e.home)); err != nil {
		return nil, err
	}
	names, err := e.databases()
	if err != nil {
		return nil, err
	}
	fx := &spanFixture{Views: map[string][]string{}}
	for _, name := range names {
		d, err := e.database(name, fx)
		if err != nil {
			return nil, err
		}
		fx.Databases = append(fx.Databases, d)
	}
	return fx, nil
}

// databases lists the databases but the system's and the test's own.
func (e *spanEngine) databases() ([]string, error) {
	rows, err := e.query("SELECT name FROM sys.databases WHERE database_id > 4 AND name <> " + literal(e.home) + " ORDER BY name")
	if err != nil {
		return nil, err
	}
	var out []string
	for _, r := range rows {
		out = append(out, *r[0])
	}
	return out, nil
}

// database reads one database as Bytebase's sync stores it (mssql/driver's
// sync.go): every user schema; its tables and views with their columns in
// column_id order; and each view's definition from sys.sql_modules. It also
// reads the tables and views each view names into fx.Views.
func (e *spanEngine) database(name string, fx *spanFixture) (*metadatapb.DatabaseSchemaMetadata, error) {
	db := bracket(name)
	rows, err := e.query(`SELECT s.name, o.name, o.type, c.name, m.definition
FROM ` + db + `.sys.schemas AS s
LEFT JOIN ` + db + `.sys.objects AS o ON o.schema_id = s.schema_id AND o.type IN ('U', 'V') AND o.is_ms_shipped = 0
LEFT JOIN ` + db + `.sys.columns AS c ON c.object_id = o.object_id
LEFT JOIN ` + db + `.sys.sql_modules AS m ON m.object_id = o.object_id
WHERE s.schema_id < 16384 AND s.name NOT IN ('guest', 'INFORMATION_SCHEMA', 'sys')
ORDER BY s.name, o.name, c.column_id`)
	if err != nil {
		return nil, err
	}
	meta := &metadatapb.DatabaseSchemaMetadata{Name: name}
	var (
		schema *metadatapb.SchemaMetadata
		table  *metadatapb.TableMetadata
		view   *metadatapb.ViewMetadata
	)
	for _, r := range rows {
		if schema == nil || schema.Name != *r[0] {
			schema = &metadatapb.SchemaMetadata{Name: *r[0]}
			meta.Schemas = append(meta.Schemas, schema)
			table, view = nil, nil
		}
		if r[1] == nil {
			continue
		}
		col := &metadatapb.ColumnMetadata{Name: *r[3], Nullable: true}
		switch strings.TrimSpace(*r[2]) {
		case "U":
			if table == nil || table.Name != *r[1] {
				table = &metadatapb.TableMetadata{Name: *r[1]}
				schema.Tables = append(schema.Tables, table)
			}
			table.Columns = append(table.Columns, col)
		case "V":
			if view == nil || view.Name != *r[1] {
				view = &metadatapb.ViewMetadata{Name: *r[1], Definition: *r[4]}
				schema.Views = append(schema.Views, view)
			}
			view.Columns = append(view.Columns, col)
		}
	}
	for _, s := range meta.Schemas {
		for _, v := range s.Views {
			if err := e.exec("USE " + bracket(name)); err != nil {
				return nil, err
			}
			reads, err := e.references(s.Name + "." + v.Name)
			if err != nil {
				return nil, fmt.Errorf("view %s.%s.%s: %w", name, s.Name, v.Name, err)
			}
			fx.Views[spanName(name, s.Name, v.Name)] = reads
		}
	}
	if err := e.exec("USE " + bracket(e.home)); err != nil {
		return nil, err
	}
	return meta, nil
}

// spanColumn is one row of sys.dm_exec_describe_first_result_set.
type spanColumn struct {
	name   string
	source string // database.schema.table.column, or "" where SQL Server names none
}

// describe asks sys.dm_exec_describe_first_result_set about query in the
// session's database with browse information on (browse_information_mode 1), and returns its
// columns but the hidden ones browse mode adds. A query SQL Server refuses
// returns its error number as failure.
func (e *spanEngine) describe(query string) (cols []spanColumn, failure string, err error) {
	rows, err := e.query(fmt.Sprintf(`SELECT name, source_database, source_schema, source_table, source_column, error_number
FROM sys.dm_exec_describe_first_result_set(%s, NULL, 1)
WHERE ISNULL(is_hidden, 0) = 0
ORDER BY column_ordinal`, literal(query)))
	if err != nil {
		return nil, "", err
	}
	for _, r := range rows {
		if r[5] != nil {
			return nil, *r[5], nil
		}
		c := spanColumn{}
		if r[0] != nil {
			c.name = *r[0]
		}
		if r[1] != nil && r[2] != nil && r[3] != nil && r[4] != nil {
			c.source = spanName(*r[1], *r[2], *r[3], *r[4])
		}
		cols = append(cols, c)
	}
	return cols, "", nil
}

// references lists the tables and views a module of the session's database
// names, as database.schema.name, from sys.dm_sql_referenced_entities.
func (e *spanEngine) references(module string) ([]string, error) {
	rows, err := e.query(`SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
FROM sys.dm_sql_referenced_entities(` + literal(module) + `, N'OBJECT')
WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL`)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, r := range rows {
		if r[1] == nil || r[2] == nil {
			return nil, fmt.Errorf("%s names an object SQL Server cannot name", module)
		}
		out = append(out, spanName(*r[0], *r[1], *r[2]))
	}
	return out, nil
}

// probe is the procedure span creates over a query, to read the tables and
// views it names.
const probe = "omni_span_probe"

// span asks SQL Server about one query in session, database.schema. The
// session's schema is the user's default schema, and sa's is dbo.
func (e *spanEngine) span(fx *spanFixture, session, statement string) (span.Answer, error) {
	database, schema, _ := strings.Cut(session, ".")
	if schema != "dbo" {
		return span.Answer{}, fmt.Errorf("session %s: the test's login is sa, whose default schema is dbo", session)
	}
	if err := e.exec("USE " + bracket(database)); err != nil {
		return span.Answer{}, err
	}
	defer e.exec("USE " + bracket(e.home))

	base, failure, err := e.describe(statement)
	if err != nil {
		return span.Answer{}, err
	}
	if failure != "" {
		return span.Answer{Error: failure}, nil
	}
	var a span.Answer
	for _, c := range base {
		a.Results = append(a.Results, c.name)
		// Browse mode names the table column a column's values are, read
		// through any view, and none for an expression. Where it cannot see
		// through a view, as for one over a UNION, it names the view's column,
		// which says nothing of the tables beneath.
		table, _, _ := cutLast(c.source, ".")
		if _, view := fx.Views[table]; c.source == "" || view {
			a.Lineage = append(a.Lineage, []string{span.Unstated})
		} else {
			a.Lineage = append(a.Lineage, []string{c.source})
		}
	}

	tables := map[string]bool{}
	read, err := e.showplan(statement)
	if err != nil {
		return span.Answer{}, err
	}
	for _, t := range read {
		tables[t] = true
	}
	if err := e.exec("CREATE PROCEDURE dbo." + probe + " AS\n" + statement); err != nil {
		return span.Answer{}, fmt.Errorf("create the probe procedure: %w", err)
	}
	named2, err := e.references("dbo." + probe)
	if err != nil {
		return span.Answer{}, err
	}
	if err := e.exec("DROP PROCEDURE dbo." + probe); err != nil {
		return span.Answer{}, err
	}
	for _, t := range named2 {
		fx.reads(t, tables)
	}
	for t := range tables {
		a.Tables = append(a.Tables, t)
	}
	return a, nil
}

// reads adds name, and the tables and views it reads if it is a view, to
// into.
func (fx *spanFixture) reads(name string, into map[string]bool) {
	if into[name] {
		return
	}
	into[name] = true
	for _, r := range fx.Views[name] {
		fx.reads(r, into)
	}
}

// showplan returns the tables the estimated plan of statement reads, the
// Object elements of its SET SHOWPLAN_XML plan, as database.schema.table.
func (e *spanEngine) showplan(statement string) ([]string, error) {
	if err := e.exec("SET SHOWPLAN_XML ON"); err != nil {
		return nil, err
	}
	rows, planErr := e.query(statement)
	if err := e.exec("SET SHOWPLAN_XML OFF"); err != nil {
		return nil, err
	}
	if planErr != nil {
		return nil, fmt.Errorf("showplan: %w", planErr)
	}
	var out []string
	for _, r := range rows {
		if r[0] == nil {
			continue
		}
		dec := xml.NewDecoder(strings.NewReader(*r[0]))
		for {
			tok, err := dec.Token()
			if err != nil {
				break
			}
			el, ok := tok.(xml.StartElement)
			if !ok || el.Name.Local != "Object" {
				continue
			}
			var db, schema, table string
			for _, at := range el.Attr {
				switch at.Name.Local {
				case "Database":
					db = unbracket(at.Value)
				case "Schema":
					schema = unbracket(at.Value)
				case "Table":
					table = unbracket(at.Value)
				}
			}
			// A worktable or a spool names no database.
			if db != "" && table != "" {
				out = append(out, spanName(db, schema, table))
			}
		}
	}
	return out, nil
}

// unbracket reads a showplan's bracketed name.
func unbracket(s string) string {
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		return strings.ReplaceAll(s[1:len(s)-1], "]]", "]")
	}
	return s
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
			out = append(out, span.Query{Fixture: fixture, Session: session, Statement: strings.TrimSuffix(line, ";")})
		}
	}
	return out, sc.Err()
}
