package catalog

import (
	"fmt"
	"io"
	"slices"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
)

// DatabaseDefinition renders a synced Redshift database snapshot as DDL.
func DatabaseDefinition(to *metadatapb.DatabaseSchemaMetadata) (string, error) {
	if to == nil {
		return "", nil
	}
	toState := convertToDatabaseState(to)

	var sb strings.Builder

	if err := writeTables(&sb, to, toState); err != nil {
		return "", err
	}
	if err := writeViews(&sb, to, toState); err != nil {
		return "", err
	}

	s := sb.String()
	// Trim the leading newline so the output starts directly with the first banner.
	s = strings.TrimLeft(s, "\n")
	return s, nil
}

// GetTableDefinition renders one table as DDL, matching DatabaseDefinition
// without the section banner.
func GetTableDefinition(schemaName string, table *metadatapb.TableMetadata) (string, error) {
	if table == nil {
		return "", nil
	}
	var sb strings.Builder
	if err := writeTable(&sb, convertToTableState(0, schemaName, table)); err != nil {
		return "", err
	}
	return sb.String(), nil
}

// GetViewDefinition renders one view as DDL, matching DatabaseDefinition
// without the section banner.
func GetViewDefinition(schemaName string, view *metadatapb.ViewMetadata) (string, error) {
	if view == nil {
		return "", nil
	}
	var sb strings.Builder
	if err := writeView(&sb, convertToViewState(0, schemaName, view)); err != nil {
		return "", err
	}
	return sb.String(), nil
}

func writeTables(w io.StringWriter, to *metadatapb.DatabaseSchemaMetadata, state *databaseState) error {
	// Follow the order of the input schemas.
	for _, schema := range to.Schemas {
		schemaState, ok := state.schemas[schema.Name]
		if !ok {
			continue
		}
		// Follow the order of the input tables.
		for _, table := range schema.Tables {
			tableState, ok := schemaState.tables[table.Name]
			if !ok {
				continue
			}
			if _, err := w.WriteString(getTableAnnouncement(tableState.name)); err != nil {
				return err
			}

			if err := writeTable(w, tableState); err != nil {
				return err
			}
			delete(schemaState.tables, tableState.name)
		}
	}
	return nil
}

func writeViews(w io.StringWriter, to *metadatapb.DatabaseSchemaMetadata, state *databaseState) error {
	// Follow the order of the input schemas.
	for _, schema := range to.Schemas {
		schemaState, ok := state.schemas[schema.Name]
		if !ok {
			continue
		}
		// Follow the order of the input views.
		for _, view := range schema.Views {
			viewState, ok := schemaState.views[view.Name]
			if !ok {
				continue
			}
			if _, err := w.WriteString(getViewAnnouncement(viewState.name)); err != nil {
				return err
			}

			if err := writeView(w, viewState); err != nil {
				return err
			}
			delete(schemaState.views, viewState.name)
		}
	}
	return nil
}

// writeTable writes one table and the comments that belong to it. Column
// comments follow the column order the CREATE TABLE body uses so output is
// deterministic across runs.
func writeTable(w io.StringWriter, table *tableState) error {
	buf := &strings.Builder{}
	if err := table.toString(buf); err != nil {
		return err
	}
	if table.comment != "" {
		if _, err := fmt.Fprintf(buf, "COMMENT ON TABLE %s IS '%s';\n", table.qualifiedName(), escapeSingleQuote(table.comment)); err != nil {
			return err
		}
	}
	for _, column := range sortedColumns(table.columns) {
		if column.comment == "" {
			continue
		}
		if _, err := fmt.Fprintf(buf, "COMMENT ON COLUMN %s.%s IS '%s';\n", table.qualifiedName(), quoteDefinitionIdentifier(column.name), escapeSingleQuote(column.comment)); err != nil {
			return err
		}
	}
	_, err := w.WriteString(buf.String())
	return err
}

// writeView writes one view and its comment.
func writeView(w io.StringWriter, view *viewState) error {
	buf := &strings.Builder{}
	if err := view.toString(buf); err != nil {
		return err
	}
	if view.comment != "" {
		if _, err := fmt.Fprintf(buf, "COMMENT ON VIEW %s IS '%s';\n", view.qualifiedName(), escapeSingleQuote(view.comment)); err != nil {
			return err
		}
	}
	_, err := w.WriteString(buf.String())
	return err
}

func getTableAnnouncement(name string) string {
	return fmt.Sprintf("\n--\n-- Table structure for `%s`\n--\n", name)
}

func getViewAnnouncement(name string) string {
	return fmt.Sprintf("\n--\n-- View structure for `%s`\n--\n", name)
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

// qualifyName prefixes the object with its schema, so DDL for analytics.orders
// cannot land in public. The prefix is dropped when no schema is known, which is
// how the single-object entry points behave when the caller passes none.
func qualifyName(schemaName, objectName string) string {
	if schemaName == "" {
		return quoteDefinitionIdentifier(objectName)
	}
	return quoteDefinitionIdentifier(schemaName) + "." + quoteDefinitionIdentifier(objectName)
}

// quoteDefinitionIdentifier renders a name as a quoted SQL identifier. Redshift
// folds an unquoted name to lower case and rejects one containing a space or a
// reserved word outright, so a name the catalog reports verbatim has to be
// quoted to come back as itself. An embedded double quote is doubled, as in
// PostgreSQL.
func quoteDefinitionIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func (t *tableState) qualifiedName() string {
	return qualifyName(t.schema, t.name)
}

func (v *viewState) qualifiedName() string {
	return qualifyName(v.schema, v.name)
}

// escapeSingleQuote encodes a comment for a SQL string literal; an apostrophe in
// "owner's orders" would otherwise close the literal and produce invalid DDL.
func escapeSingleQuote(s string) string {
	return strings.ReplaceAll(s, "'", "''")
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
		state.tables[table.Name] = convertToTableState(i, schema.Name, table)
	}
	for i, view := range schema.Views {
		if view == nil {
			continue
		}
		state.views[view.Name] = convertToViewState(i, schema.Name, view)
	}
	return state
}

type tableState struct {
	id          int
	schema      string
	name        string
	columns     map[string]*columnState
	indexes     map[string]*indexState
	foreignKeys map[string]*foreignKeyState
	comment     string
}

func (t *tableState) toString(buf *strings.Builder) error {
	if _, err := fmt.Fprintf(buf, "CREATE TABLE %s (\n  ", t.qualifiedName()); err != nil {
		return err
	}
	columns := sortedColumns(t.columns)
	for i, column := range columns {
		if i > 0 {
			if _, err := buf.WriteString(",\n  "); err != nil {
				return err
			}
		}
		if err := column.toString(buf); err != nil {
			return err
		}
	}

	indexes := make([]*indexState, 0, len(t.indexes))
	for _, index := range t.indexes {
		indexes = append(indexes, index)
	}
	slices.SortFunc(indexes, func(a, b *indexState) int {
		if a.primary && !b.primary {
			return -1
		}
		if !a.primary && b.primary {
			return 1
		}
		return strings.Compare(a.name, b.name)
	})

	for i, index := range indexes {
		if i+len(columns) > 0 {
			if _, err := buf.WriteString(",\n  "); err != nil {
				return err
			}
		}
		if err := index.toString(buf); err != nil {
			return err
		}
	}

	foreignKeys := make([]*foreignKeyState, 0, len(t.foreignKeys))
	for _, fk := range t.foreignKeys {
		foreignKeys = append(foreignKeys, fk)
	}
	slices.SortFunc(foreignKeys, func(a, b *foreignKeyState) int {
		return strings.Compare(a.name, b.name)
	})

	for i, fk := range foreignKeys {
		if i+len(columns)+len(indexes) > 0 {
			if _, err := buf.WriteString(",\n  "); err != nil {
				return err
			}
		}
		if err := fk.toString(buf); err != nil {
			return err
		}
	}

	if _, err := buf.WriteString("\n);\n"); err != nil {
		return err
	}
	return nil
}

func newTableState(id int, name string) *tableState {
	return &tableState{
		id:          id,
		name:        name,
		columns:     make(map[string]*columnState),
		indexes:     make(map[string]*indexState),
		foreignKeys: make(map[string]*foreignKeyState),
	}
}

// sortedColumns returns a table's columns in declaration order. Both the
// CREATE TABLE body and the trailing COMMENT ON COLUMN statements go through
// it, so the two agree and neither depends on map iteration order.
func sortedColumns(columns map[string]*columnState) []*columnState {
	sorted := make([]*columnState, 0, len(columns))
	for _, column := range columns {
		sorted = append(sorted, column)
	}
	slices.SortFunc(sorted, func(a, b *columnState) int {
		if a.id < b.id {
			return -1
		}
		if a.id > b.id {
			return 1
		}
		return 0
	})
	return sorted
}

func convertToTableState(id int, schemaName string, table *metadatapb.TableMetadata) *tableState {
	state := newTableState(id, table.Name)
	state.schema = schemaName
	state.comment = table.Comment
	for i, column := range table.Columns {
		if column == nil {
			continue
		}
		state.columns[column.Name] = convertToColumnState(i, column)
	}
	for i, index := range table.Indexes {
		if index == nil {
			continue
		}
		state.indexes[index.Name] = convertToIndexState(i, index)
	}
	for i, fk := range table.ForeignKeys {
		if fk == nil {
			continue
		}
		state.foreignKeys[fk.Name] = convertToForeignKeyState(i, fk)
	}
	return state
}

type foreignKeyState struct {
	id                int
	name              string
	columns           []string
	referencedSchema  string
	referencedTable   string
	referencedColumns []string
}

func convertToForeignKeyState(id int, foreignKey *metadatapb.ForeignKeyMetadata) *foreignKeyState {
	return &foreignKeyState{
		id:                id,
		name:              foreignKey.Name,
		columns:           foreignKey.Columns,
		referencedSchema:  foreignKey.ReferencedSchema,
		referencedTable:   foreignKey.ReferencedTable,
		referencedColumns: foreignKey.ReferencedColumns,
	}
}

func (f *foreignKeyState) toString(buf *strings.Builder) error {
	for i, column := range f.columns {
		if i >= len(f.referencedColumns) {
			break
		}
		if _, err := buf.WriteString("FOREIGN KEY ("); err != nil {
			return err
		}
		if _, err := buf.WriteString(quoteDefinitionIdentifier(column)); err != nil {
			return err
		}
		if _, err := buf.WriteString(") REFERENCES "); err != nil {
			return err
		}
		referencedColumn := f.referencedColumns[i]
		// The sync strips the qualifier off ReferencedTable and keeps it in
		// ReferencedSchema, so an unqualified target here would bind to whatever
		// the search path resolves rather than the table the constraint names.
		referenced := qualifyName(f.referencedSchema, f.referencedTable)
		if _, err := fmt.Fprintf(buf, "%s(%s)", referenced, quoteDefinitionIdentifier(referencedColumn)); err != nil {
			return err
		}
	}
	return nil
}

type indexState struct {
	id      int
	name    string
	keys    []string
	lengths []int64
	primary bool
	unique  bool
	tp      string
}

func convertToIndexState(id int, index *metadatapb.IndexMetadata) *indexState {
	return &indexState{
		id:      id,
		name:    index.Name,
		keys:    index.Expressions,
		lengths: index.KeyLength,
		primary: index.Primary,
		unique:  index.Unique,
		tp:      index.Type,
	}
}

func (i *indexState) toString(buf *strings.Builder) error {
	var keyword string
	switch {
	case i.primary:
		keyword = "PRIMARY KEY ("
	case i.unique:
		keyword = "UNIQUE ("
	default:
		return nil
	}
	if _, err := buf.WriteString(keyword); err != nil {
		return err
	}
	for idx, key := range i.keys {
		if idx > 0 {
			if _, err := buf.WriteString(", "); err != nil {
				return err
			}
		}
		// Written verbatim: these come from pg_get_indexdef(.., true), which
		// already quotes what needs it and can be an expression rather than
		// a column, so quoting again would name something else.
		if _, err := buf.WriteString(key); err != nil {
			return err
		}
	}
	_, err := buf.WriteString(")")
	return err
}

type defaultValue interface {
	toString() string
}

type defaultValueNull struct{}

func (*defaultValueNull) toString() string {
	return "NULL"
}

type defaultValueString struct {
	value string
}

func (d *defaultValueString) toString() string {
	return fmt.Sprintf("'%s'", d.value)
}

type defaultValueExpression struct {
	value string
}

func (d *defaultValueExpression) toString() string {
	return d.value
}

type columnState struct {
	id           int
	name         string
	tp           string
	defaultValue defaultValue
	comment      string
	nullable     bool
}

func (c *columnState) toString(buf *strings.Builder) error {
	if _, err := fmt.Fprintf(buf, "%s ", quoteDefinitionIdentifier(c.name)); err != nil {
		return err
	}
	if _, err := buf.WriteString(c.tp); err != nil {
		return err
	}
	if !c.nullable {
		if _, err := buf.WriteString(" NOT NULL"); err != nil {
			return err
		}
	}
	if c.defaultValue != nil {
		if _, err := fmt.Fprintf(buf, " DEFAULT %s", c.defaultValue.toString()); err != nil {
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
	if column.Default == "NULL" {
		result.defaultValue = &defaultValueNull{}
	} else if column.Default != "" {
		if isExpression(column.Default) {
			result.defaultValue = &defaultValueExpression{value: column.Default}
		} else {
			result.defaultValue = &defaultValueString{value: column.Default}
		}
	}
	return result
}

func isExpression(value string) bool {
	if strings.Contains(value, "(") || strings.Contains(value, ")") ||
		strings.Contains(value, "::") || strings.Contains(value, "NOW") ||
		strings.Contains(value, "CURRENT_") || strings.Contains(value, "GETDATE") ||
		strings.Contains(value, "TIMESTAMP") || strings.Contains(value, "SEQUENCE") {
		return true
	}
	return false
}

type viewState struct {
	id         int
	schema     string
	name       string
	definition string
	comment    string
}

func convertToViewState(id int, schemaName string, view *metadatapb.ViewMetadata) *viewState {
	return &viewState{
		id:         id,
		schema:     schemaName,
		name:       view.Name,
		definition: view.Definition,
		comment:    view.Comment,
	}
}

func (v *viewState) toString(buf io.StringWriter) error {
	stmt := fmt.Sprintf("CREATE OR REPLACE VIEW %s AS %s", v.qualifiedName(), v.definition)
	if !strings.HasSuffix(stmt, ";") {
		stmt += ";"
	}
	stmt += "\n"
	_, err := buf.WriteString(stmt)
	return err
}
