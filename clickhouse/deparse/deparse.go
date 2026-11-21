// Package deparse renders ClickHouse schema metadata back into SQL DDL.
package deparse

import (
	"fmt"
	"io"
	"slices"
	"strings"
	"unicode"

	metadatapb "github.com/bytebase/omni/metadata"
)

// DatabaseDefinition renders a ClickHouse database schema snapshot as DDL
// statements for all tables followed by all views, preserving snapshot order.
func DatabaseDefinition(metadata *metadatapb.DatabaseSchemaMetadata) (string, error) {
	if metadata == nil {
		return "", nil
	}
	toState := convertToDatabaseState(metadata)
	var sb strings.Builder

	if err := writeTables(&sb, metadata, toState); err != nil {
		return "", err
	}
	if err := writeViews(&sb, metadata, toState); err != nil {
		return "", err
	}

	// Trim leading newlines so YAML/block-scalar consumers receive clean DDL.
	return strings.TrimLeft(sb.String(), "\n"), nil
}

// GetTableDefinition renders one ClickHouse table as its section banner and
// CREATE TABLE statement.
func GetTableDefinition(_ string, table *metadatapb.TableMetadata) (string, error) {
	if table == nil {
		return "", nil
	}
	var buf strings.Builder
	tableState := convertToTableState(0, table)

	if _, err := buf.WriteString(getTableAnnouncement(table.Name)); err != nil {
		return "", err
	}
	if err := tableState.toString(&buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GetViewDefinition renders one ClickHouse view as its section banner and
// CREATE VIEW statement.
func GetViewDefinition(_ string, view *metadatapb.ViewMetadata) (string, error) {
	if view == nil {
		return "", nil
	}
	var buf strings.Builder
	viewState := convertToViewState(0, view)

	if _, err := buf.WriteString(getViewAnnouncement(view.Name)); err != nil {
		return "", err
	}
	if err := viewState.toString(&buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func writeTables(w io.StringWriter, to *metadatapb.DatabaseSchemaMetadata, state *databaseState) error {
	// Follow the order of the input schemas.
	for _, schema := range to.Schemas {
		if schema == nil {
			continue
		}
		schemaState, ok := state.schemas[schema.Name]
		if !ok {
			continue
		}
		// Follow the order of the input tables.
		for _, table := range schema.Tables {
			if table == nil {
				continue
			}
			ts, ok := schemaState.tables[table.Name]
			if !ok {
				continue
			}
			if _, err := w.WriteString(getTableAnnouncement(ts.name)); err != nil {
				return err
			}

			buf := &strings.Builder{}
			if err := ts.toString(buf); err != nil {
				return err
			}
			if _, err := w.WriteString(buf.String()); err != nil {
				return err
			}
			delete(schemaState.tables, ts.name)
		}
	}
	return nil
}

func writeViews(w io.StringWriter, to *metadatapb.DatabaseSchemaMetadata, state *databaseState) error {
	// Follow the order of the input schemas.
	for _, schema := range to.Schemas {
		if schema == nil {
			continue
		}
		schemaState, ok := state.schemas[schema.Name]
		if !ok {
			continue
		}
		// Follow the order of the input views.
		for _, view := range schema.Views {
			if view == nil {
				continue
			}
			vs, ok := schemaState.views[view.Name]
			if !ok {
				continue
			}
			if _, err := w.WriteString(getViewAnnouncement(vs.name)); err != nil {
				return err
			}

			buf := &strings.Builder{}
			if err := vs.toString(buf); err != nil {
				return err
			}
			if _, err := w.WriteString(buf.String()); err != nil {
				return err
			}
			delete(schemaState.views, vs.name)
		}
	}
	return nil
}

func getTableAnnouncement(name string) string {
	return fmt.Sprintf("\n--\n-- Table structure for %s\n--\n", quoteIdentifier(name))
}

func getViewAnnouncement(name string) string {
	return fmt.Sprintf("\n--\n-- View structure for %s\n--\n", quoteIdentifier(name))
}

// quoteIdentifier writes name in ClickHouse backquotes, escaping backslashes,
// backquotes, and control characters so identifiers with spaces, keywords,
// punctuation, or line breaks produce valid ClickHouse SQL and safe comments.
func quoteIdentifier(name string) string {
	var b strings.Builder
	b.WriteByte('`')
	for i := 0; i < len(name); i++ {
		switch c := name[i]; c {
		case '\\', '`':
			b.WriteByte('\\')
			b.WriteByte(c)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		case '\b':
			b.WriteString(`\b`)
		case '\f':
			b.WriteString(`\f`)
		case 0:
			b.WriteString(`\0`)
		default:
			b.WriteByte(c)
		}
	}
	b.WriteByte('`')
	return b.String()
}

// quoteStringLiteral writes s as a ClickHouse single-quoted string literal,
// escaping single quotes, backslashes, and control characters.
func quoteStringLiteral(s string) string {
	var b strings.Builder
	b.WriteByte('\'')
	for i := 0; i < len(s); i++ {
		switch c := s[i]; c {
		case '\\', '\'':
			b.WriteByte('\\')
			b.WriteByte(c)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		case '\b':
			b.WriteString(`\b`)
		case '\f':
			b.WriteString(`\f`)
		case 0:
			b.WriteString(`\0`)
		default:
			b.WriteByte(c)
		}
	}
	b.WriteByte('\'')
	return b.String()
}

type databaseState struct {
	name    string
	schemas map[string]*schemaState
}

func newDatabaseState() *databaseState {
	return &databaseState{
		schemas: make(map[string]*schemaState),
	}
}

func convertToDatabaseState(database *metadatapb.DatabaseSchemaMetadata) *databaseState {
	state := newDatabaseState()
	state.name = database.Name
	for _, schema := range database.Schemas {
		if schema == nil {
			continue
		}
		state.schemas[schema.Name] = convertToSchemaState(schema)
	}
	return state
}

type schemaState struct {
	name   string
	tables map[string]*tableState
	views  map[string]*viewState
}

func newSchemaState() *schemaState {
	return &schemaState{
		tables: make(map[string]*tableState),
		views:  make(map[string]*viewState),
	}
}

func convertToSchemaState(schema *metadatapb.SchemaMetadata) *schemaState {
	state := newSchemaState()
	state.name = schema.Name
	for i, table := range schema.Tables {
		if table == nil {
			continue
		}
		state.tables[table.Name] = convertToTableState(i, table)
	}
	for i, view := range schema.Views {
		if view == nil {
			continue
		}
		state.views[view.Name] = convertToViewState(i, view)
	}
	return state
}

type tableState struct {
	id          int
	name        string
	columns     map[string]*columnState
	sortingKeys []string
	primaryKeys []string
	comment     string
	engine      string
}

func (t *tableState) toString(buf *strings.Builder) error {
	if _, err := fmt.Fprintf(buf, "CREATE TABLE %s (\n  ", quoteIdentifier(t.name)); err != nil {
		return err
	}
	columns := make([]*columnState, 0, len(t.columns))
	for _, column := range t.columns {
		columns = append(columns, column)
	}
	slices.SortFunc(columns, func(x, y *columnState) int {
		if x.id < y.id {
			return -1
		} else if x.id > y.id {
			return 1
		}
		return 0
	})
	for i, column := range columns {
		if i > 0 {
			if _, err := fmt.Fprint(buf, ",\n  "); err != nil {
				return err
			}
		}
		if err := column.toString(buf); err != nil {
			return err
		}
	}
	if _, err := buf.WriteString("\n)"); err != nil {
		return err
	}
	if t.engine != "" {
		if _, err := fmt.Fprintf(buf, "\nENGINE = %s", t.engine); err != nil {
			return err
		}
	}
	if len(t.sortingKeys) == 1 {
		if _, err := fmt.Fprintf(buf, "\nORDER BY %s", t.sortingKeys[0]); err != nil {
			return err
		}
	} else if len(t.sortingKeys) > 1 {
		if _, err := fmt.Fprintf(buf, "\nORDER BY (%s)", strings.Join(t.sortingKeys, ", ")); err != nil {
			return err
		}
	} else {
		// MergeTree tables require an explicit ORDER BY clause; tuple() matches ClickHouse's empty-order form.
		if strings.ToLower(t.engine) == "mergetree" {
			if _, err := buf.WriteString("\nORDER BY tuple()"); err != nil {
				return err
			}
		}
	}
	if len(t.primaryKeys) == 1 {
		if _, err := fmt.Fprintf(buf, "\nPRIMARY KEY %s", t.primaryKeys[0]); err != nil {
			return err
		}
	} else if len(t.primaryKeys) > 1 {
		if _, err := fmt.Fprintf(buf, "\nPRIMARY KEY (%s)", strings.Join(t.primaryKeys, ", ")); err != nil {
			return err
		}
	}
	if t.comment != "" {
		if _, err := fmt.Fprintf(buf, "\nCOMMENT %s", quoteStringLiteral(t.comment)); err != nil {
			return err
		}
	}
	if _, err := buf.WriteString(";\n"); err != nil {
		return err
	}
	return nil
}

func newTableState(id int, name string) *tableState {
	return &tableState{
		id:      id,
		name:    name,
		columns: make(map[string]*columnState),
	}
}

func convertToTableState(id int, table *metadatapb.TableMetadata) *tableState {
	state := newTableState(id, table.Name)
	state.comment = table.Comment
	state.engine = table.Engine
	state.sortingKeys = table.SortingKeys
	for i, column := range table.Columns {
		if column == nil {
			continue
		}
		state.columns[column.Name] = convertToColumnState(i, column)
	}
	for _, index := range table.Indexes {
		if index != nil && index.Primary {
			state.primaryKeys = index.Expressions
			break
		}
	}
	return state
}

type columnState struct {
	id         int
	name       string
	tp         string
	defaultSQL string
	comment    string
	nullable   bool
}

func (c *columnState) toString(buf *strings.Builder) error {
	if _, err := fmt.Fprintf(buf, "%s %s", quoteIdentifier(c.name), c.tp); err != nil {
		return err
	}
	if c.defaultSQL != "" {
		if _, err := fmt.Fprintf(buf, " DEFAULT %s", c.defaultSQL); err != nil {
			return err
		}
	}
	if c.comment != "" {
		if _, err := fmt.Fprintf(buf, " COMMENT %s", quoteStringLiteral(c.comment)); err != nil {
			return err
		}
	}
	return nil
}

func convertToColumnState(id int, column *metadatapb.ColumnMetadata) *columnState {
	result := &columnState{
		id:       id,
		name:     column.Name,
		tp:       column.Type,
		nullable: column.Nullable,
		comment:  column.Comment,
	}
	def := strings.TrimSpace(column.Default)
	if strings.EqualFold(def, "NULL") {
		// ClickHouse sync stores "NULL" when information_schema.columns.column_default
		// is empty (even for non-nullable columns); only nullable columns accept DEFAULT NULL.
		if column.Nullable {
			result.defaultSQL = "NULL"
		}
	} else if def != "" {
		// information_schema.columns.column_default holds a ClickHouse SQL expression
		// (already quoted for string/date literals), so emit it verbatim.
		result.defaultSQL = def
	}
	return result
}

type viewState struct {
	id         int
	name       string
	definition string
	comment    string
}

func convertToViewState(id int, view *metadatapb.ViewMetadata) *viewState {
	return &viewState{
		id:         id,
		name:       view.Name,
		definition: view.Definition,
		comment:    view.Comment,
	}
}

func (v *viewState) toString(buf io.StringWriter) error {
	trimmed := trimStatement(v.definition)
	var stmt string
	if isCreateViewStatement(trimmed) {
		// clickhouse/driver.SyncDatabase populates ViewMetadata.Definition from
		// system.tables.create_table_query, which is already a complete CREATE VIEW
		// statement (including any COMMENT clause).
		stmt = trimmed + ";\n"
	} else {
		stmt = fmt.Sprintf("CREATE OR REPLACE VIEW %s AS %s", quoteIdentifier(v.name), trimmed)
		if v.comment != "" {
			stmt += " COMMENT " + quoteStringLiteral(v.comment)
		}
		stmt += ";\n"
	}
	if _, err := buf.WriteString(stmt); err != nil {
		return err
	}
	return nil
}

// isCreateViewStatement reports whether stmt is already a full CREATE ... VIEW
// statement (as returned by ClickHouse's system.tables.create_table_query) rather
// than a bare SELECT query body.
func isCreateViewStatement(stmt string) bool {
	words := strings.Fields(stmt)
	if len(words) < 2 || !strings.EqualFold(words[0], "CREATE") {
		return false
	}
	i := 1
	if i+1 < len(words) && strings.EqualFold(words[i], "OR") && strings.EqualFold(words[i+1], "REPLACE") {
		i += 2
	}
	if i < len(words) && (strings.EqualFold(words[i], "MATERIALIZED") ||
		strings.EqualFold(words[i], "LIVE") ||
		strings.EqualFold(words[i], "WINDOW")) {
		i++
	}
	return i < len(words) && strings.EqualFold(words[i], "VIEW")
}

func trimStatement(statement string) string {
	return strings.TrimLeftFunc(strings.TrimRightFunc(statement, func(r rune) bool {
		return unicode.IsSpace(r) || r == ';'
	}), unicode.IsSpace)
}
