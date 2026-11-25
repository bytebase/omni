package catalog

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	crparser "github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
	plpgsqlparser "github.com/cockroachdb/cockroachdb-parser/pkg/sql/plpgsql/parser"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/types"

	metadatapb "github.com/bytebase/omni/metadata"
)

// SDLMigrationOptions configures BuildSDLMigrationPlan.
type SDLMigrationOptions struct {
	// SkipSchemas lists schema names left out of the migration, such as a
	// backup schema.
	SkipSchemas []string
}

// BuildSDLMigrationPlan loads two SDL texts with LoadSDL and returns the
// migration from the first to the second, or nil when they hold the same
// schema.
func BuildSDLMigrationPlan(sourceSDL, targetSDL string, opts SDLMigrationOptions) (*MigrationPlan, error) {
	from, err := LoadSDL(sourceSDL)
	if err != nil {
		return nil, fmt.Errorf("failed to load source schema: %w", err)
	}
	to, err := LoadSDL(targetSDL)
	if err != nil {
		return nil, fmt.Errorf("failed to load target schema: %w", err)
	}
	plan, err := MetadataMigrationPlan(from, to, MetadataMigrationOptions(opts))
	if err != nil || len(plan.Ops) == 0 {
		return nil, err
	}
	return plan, nil
}

// LoadSDL reads a schema written as CockroachDB DDL, such as DatabaseDefinition
// renders, into the metadata a sync would state for it. Every type, default,
// check and index key passes through CockroachDB's formatter, so SDL written
// by hand and SDL rendered from a sync load alike where CockroachDB would
// store them alike. Session statements (SET, SELECT set_config) and owner
// changes are skipped; any other statement LoadSDL does not model is an error.
func LoadSDL(text string) (*metadatapb.DatabaseSchemaMetadata, error) {
	stmts, err := crparser.Parse(text)
	if err != nil {
		return nil, err
	}
	l := &sdlLoader{meta: &metadatapb.DatabaseSchemaMetadata{}, schemas: make(map[string]*metadatapb.SchemaMetadata)}
	for _, stmt := range stmts {
		if err := l.load(stmt.AST); err != nil {
			return nil, fmt.Errorf("%s: %w", stmt.SQL, err)
		}
	}
	for _, view := range l.starViews {
		for _, name := range l.starNames(view.query) {
			*view.columns = append(*view.columns, &metadatapb.ColumnMetadata{Name: name})
		}
	}
	for _, fk := range l.keyless {
		referenced := schemaTable(l.schema(fk.ReferencedSchema), fk.ReferencedTable)
		if referenced == nil {
			return nil, fmt.Errorf("foreign key %s references table %s.%s, which is not created", fk.Name, fk.ReferencedSchema, fk.ReferencedTable)
		}
		for _, index := range referenced.Indexes {
			if index.Primary {
				fk.ReferencedColumns = append(fk.ReferencedColumns, index.Expressions...)
			}
		}
		if len(fk.ReferencedColumns) == 0 {
			// Its key is the hidden rowid, which the sync does not read.
			return nil, fmt.Errorf("foreign key %s references table %s.%s, which declares no primary key", fk.Name, fk.ReferencedSchema, fk.ReferencedTable)
		}
	}
	return l.meta, nil
}

type sdlLoader struct {
	meta    *metadatapb.DatabaseSchemaMetadata
	schemas map[string]*metadatapb.SchemaMetadata
	// keyless holds the foreign keys that name no referenced columns.
	keyless []*metadatapb.ForeignKeyMetadata
	// starViews holds the views whose output names wait on the columns of
	// the relation they select a star from.
	starViews []starView
}

// star reports whether a select expression is * or relation.*.
func star(expr tree.Expr) bool {
	switch e := expr.(type) {
	case tree.UnqualifiedStar:
		return true
	case *tree.UnresolvedName:
		return e.Star
	}
	return false
}

type starView struct {
	columns *[]*metadatapb.ColumnMetadata
	query   *tree.Select
}

// starNames returns a query's output names where it selects from one stored
// relation the SDL declares, a star standing for that relation's visible
// columns, as CockroachDB expands it when it creates the view; nil when it
// cannot tell.
func (l *sdlLoader) starNames(query *tree.Select) []string {
	if query.With != nil {
		return nil
	}
	sel, ok := query.Select.(*tree.SelectClause)
	if !ok || len(sel.From.Tables) != 1 {
		return nil
	}
	aliased, ok := sel.From.Tables[0].(*tree.AliasedTableExpr)
	if !ok {
		return nil
	}
	tn, ok := aliased.Expr.(*tree.TableName)
	if !ok {
		return nil
	}
	s, ok := l.schemas[schemaOf(tn)]
	if !ok {
		return nil
	}
	name := string(tn.ObjectName)
	var relation []string
	if table := schemaTable(s, name); table != nil {
		for _, column := range table.Columns {
			if !column.IsInvisible {
				relation = append(relation, column.Name)
			}
		}
	} else if view := findView(s, name); view != nil {
		for _, column := range view.Columns {
			relation = append(relation, column.Name)
		}
	} else if view := findMaterializedView(s, name); view != nil {
		for _, column := range view.Columns {
			relation = append(relation, column.Name)
		}
	}
	if len(relation) == 0 {
		return nil
	}
	var names []string
	for _, expr := range sel.Exprs {
		if star(expr.Expr) && expr.As == "" {
			names = append(names, relation...)
			continue
		}
		if expr.As != "" {
			names = append(names, string(expr.As))
			continue
		}
		name, ok := outputName(expr.Expr)
		if !ok {
			return nil
		}
		names = append(names, name)
	}
	return names
}

func (l *sdlLoader) schema(name string) *metadatapb.SchemaMetadata {
	if name == "" {
		name = "public"
	}
	if s := l.schemas[name]; s != nil {
		return s
	}
	s := &metadatapb.SchemaMetadata{Name: name}
	l.schemas[name] = s
	l.meta.Schemas = append(l.meta.Schemas, s)
	return s
}

func (l *sdlLoader) load(stmt tree.Statement) error {
	switch n := stmt.(type) {
	case *tree.SetVar, *tree.Select, *tree.AlterTableOwner:
		return nil
	case *tree.AlterSchema:
		// A dump sets owners; a rename is a change, not a definition.
		if _, ok := n.Cmd.(*tree.AlterSchemaOwner); !ok {
			return fmt.Errorf("ALTER SCHEMA %s is not a schema definition", tree.AsString(n.Cmd))
		}
		return nil
	case *tree.CreateSchema:
		name := string(n.Schema.SchemaName)
		if name == "" {
			// CREATE SCHEMA AUTHORIZATION role names the schema after it.
			if n.AuthRole.RoleSpecType != tree.RoleName {
				return fmt.Errorf("CREATE SCHEMA AUTHORIZATION %s names no schema", n.AuthRole.Name)
			}
			name = n.AuthRole.Name
		}
		l.schema(name)
	case *tree.CreateExtension:
		l.meta.Extensions = append(l.meta.Extensions, &metadatapb.ExtensionMetadata{Name: n.Name.String()})
	case *tree.CreateType:
		return l.createType(n)
	case *tree.CreateSequence:
		if n.Persistence.IsTemporary() {
			return fmt.Errorf("sequence %s: a temporary sequence is not a schema definition", n.Name.String())
		}
		s := l.schema(schemaOf(&n.Name))
		if skip, err := repeated(findSequence(s, string(n.Name.ObjectName)) != nil, n.IfNotExists, n.Name.String()); skip || err != nil {
			return err
		}
		sequence := &metadatapb.SequenceMetadata{Name: string(n.Name.ObjectName)}
		if err := applySequenceOptions(sequence, n.Options, true); err != nil {
			return err
		}
		s.Sequences = append(s.Sequences, sequence)
	case *tree.AlterSequence:
		s := l.schema(n.Name.Schema())
		sequence := findSequence(s, n.Name.Object())
		if sequence == nil {
			return fmt.Errorf("sequence %s is not created", n.Name)
		}
		return applySequenceOptions(sequence, n.Options, false)
	case *tree.CreateTable:
		return l.createTable(n)
	case *tree.CreateIndex:
		schema := schemaOf(&n.Table)
		table := l.relationIndexes(schema, string(n.Table.ObjectName))
		if table == nil {
			return fmt.Errorf("table %s is not created", n.Table.String())
		}
		exists := n.Name != "" && slices.ContainsFunc(*table, func(i *metadatapb.IndexMetadata) bool { return i.Name == string(n.Name) })
		if skip, err := repeated(exists, n.IfNotExists, string(n.Name)); skip || err != nil {
			return err
		}
		*table = append(*table, sdlIndex(schema, string(n.Table.ObjectName), *table, n))
	case *tree.CreateView:
		return l.createView(n)
	case *tree.CreateRoutine:
		return l.createRoutine(n)
	case *tree.AlterTable:
		return l.alterTable(n)
	case *tree.CommentOnSchema:
		l.schema(string(n.Name.SchemaName)).Comment = comment(n.Comment)
	case *tree.CommentOnTable:
		table := schemaTable(l.schema(n.Table.Schema()), n.Table.Object())
		if table == nil {
			return fmt.Errorf("table %s is not created", n.Table)
		}
		table.Comment = comment(n.Comment)
	case *tree.CommentOnColumn:
		table := schemaTable(l.schema(n.TableName.Schema()), n.TableName.Object())
		if table == nil {
			return fmt.Errorf("table %s is not created", n.TableName)
		}
		column := findColumn(table, string(n.ColumnName))
		if column == nil {
			return fmt.Errorf("column %s is not created", n.ColumnItem)
		}
		column.Comment = comment(n.Comment)
	case *tree.CommentOnIndex:
		index := l.findIndex(schemaOf(&n.Index.Table), string(n.Index.Table.ObjectName), string(n.Index.Index))
		if index == nil {
			return fmt.Errorf("index %s is not created", n.Index.String())
		}
		index.Comment = comment(n.Comment)
	case *tree.CommentOnConstraint:
		return l.commentOnConstraint(n)
	case *tree.CommentOnType:
		s := l.schema(n.Name.Schema())
		for _, enum := range s.EnumTypes {
			if enum.Name == n.Name.Object() {
				enum.Comment = comment(n.Comment)
				return nil
			}
		}
		for _, composite := range s.CompositeTypes {
			if composite.Name == n.Name.Object() {
				composite.Comment = comment(n.Comment)
				return nil
			}
		}
		return fmt.Errorf("type %s is not created", n.Name)
	default:
		return fmt.Errorf("%s is not a schema definition statement", stmt.StatementTag())
	}
	return nil
}

// repeated decides a CREATE of an object the SDL already created, as the
// engine does: IF NOT EXISTS skips it, and without it the CREATE fails.
func repeated(exists, ifNotExists bool, name string) (bool, error) {
	if !exists {
		return false, nil
	}
	if ifNotExists {
		return true, nil
	}
	return false, fmt.Errorf("%s is created twice", name)
}

// schemaTable finds a table of a loaded schema.
func schemaTable(s *metadatapb.SchemaMetadata, name string) *metadatapb.TableMetadata {
	for _, table := range s.Tables {
		if table.Name == name {
			return table
		}
	}
	return nil
}

func comment(c *string) string {
	if c == nil {
		return ""
	}
	return *c
}

func (l *sdlLoader) createType(n *tree.CreateType) error {
	s := l.schema(n.TypeName.Schema())
	exists := slices.ContainsFunc(s.EnumTypes, func(e *metadatapb.EnumTypeMetadata) bool { return e.Name == n.TypeName.Object() }) ||
		slices.ContainsFunc(s.CompositeTypes, func(c *metadatapb.CompositeTypeMetadata) bool { return c.Name == n.TypeName.Object() })
	if skip, err := repeated(exists, n.IfNotExists, n.TypeName.String()); skip || err != nil {
		return err
	}
	switch n.Variety {
	case tree.Enum:
		enum := &metadatapb.EnumTypeMetadata{Name: n.TypeName.Object()}
		for _, label := range n.EnumLabels {
			enum.Values = append(enum.Values, string(label))
		}
		s.EnumTypes = append(s.EnumTypes, enum)
	case tree.Composite:
		composite := &metadatapb.CompositeTypeMetadata{Name: n.TypeName.Object()}
		for _, elem := range n.CompositeTypeList {
			typ, collation := sdlType(elem.Type)
			composite.Attributes = append(composite.Attributes, &metadatapb.CompositeTypeAttribute{Name: string(elem.Label), Type: typ, Collation: collation})
		}
		s.CompositeTypes = append(s.CompositeTypes, composite)
	default:
		return fmt.Errorf("type %s is neither an enum nor a composite type", n.TypeName)
	}
	return nil
}

// applySequenceOptions sets a sequence's options. A new sequence takes
// CockroachDB's defaults for the options it leaves out. An ALTER SEQUENCE
// sets a bound it resets with NO MINVALUE or NO MAXVALUE to INT8's default,
// and widens a bound it does not name that is its type's limit to INT8's,
// whatever the sequence's type: v25.2.3 does both.
func applySequenceOptions(sequence *metadatapb.SequenceMetadata, options tree.SequenceOptions, create bool) error {
	minStated, maxStated := false, false
	for _, option := range options {
		value := ""
		if option.IntVal != nil {
			value = strconv.FormatInt(*option.IntVal, 10)
		}
		switch option.Name {
		case tree.SeqOptAs:
			sequence.DataType = option.AsIntegerType.SQLString()
		case tree.SeqOptIncrement:
			sequence.Increment = value
		case tree.SeqOptMinValue:
			sequence.MinValue, minStated = value, true
		case tree.SeqOptMaxValue:
			sequence.MaxValue, maxStated = value, true
		case tree.SeqOptStart:
			sequence.Start = value
		case tree.SeqOptCache:
			sequence.CacheSize = value
		case tree.SeqOptCycle:
			sequence.Cycle = true
		case tree.SeqOptNoCycle:
			sequence.Cycle = false
		case tree.SeqOptOwnedBy:
			if option.ColumnItemVal == nil {
				sequence.OwnerTable, sequence.OwnerColumn = "", ""
			} else {
				sequence.OwnerTable = option.ColumnItemVal.TableName.Object()
				sequence.OwnerColumn = string(option.ColumnItemVal.ColumnName)
			}
		default:
			return fmt.Errorf("sequence option %s is not supported", option.Name)
		}
	}
	if !create {
		limit := sequenceLimit(sequence.DataType)
		if !minStated && sequence.MinValue == strconv.FormatInt(-limit-1, 10) {
			sequence.MinValue = strconv.FormatInt(math.MinInt64, 10)
		}
		if !maxStated && sequence.MaxValue == strconv.FormatInt(limit, 10) {
			sequence.MaxValue = strconv.FormatInt(math.MaxInt64, 10)
		}
		defaultBounds(sequence, math.MaxInt64)
		return nil
	}
	if sequence.DataType == "" {
		sequence.DataType = "INT8"
	}
	if sequence.Increment == "" {
		sequence.Increment = "1"
	}
	defaultBounds(sequence, sequenceLimit(sequence.DataType))
	if sequence.Start == "" {
		sequence.Start = sequence.MinValue
		if strings.HasPrefix(sequence.Increment, "-") {
			sequence.Start = sequence.MaxValue
		}
	}
	return nil
}

// sequenceLimit is the largest value of a sequence's type.
func sequenceLimit(dataType string) int64 {
	switch dataType {
	case "INT2":
		return math.MaxInt16
	case "INT4":
		return math.MaxInt32
	}
	return math.MaxInt64
}

// defaultBounds gives a sequence the bounds it does not state, for its
// direction, up to limit.
func defaultBounds(sequence *metadatapb.SequenceMetadata, limit int64) {
	descending := strings.HasPrefix(sequence.Increment, "-")
	if sequence.MinValue == "" {
		sequence.MinValue = "1"
		if descending {
			sequence.MinValue = strconv.FormatInt(-limit-1, 10)
		}
	}
	if sequence.MaxValue == "" {
		sequence.MaxValue = strconv.FormatInt(limit, 10)
		if descending {
			sequence.MaxValue = "-1"
		}
	}
}

// sdlType returns a column type as CockroachDB's formatter writes it, with a
// collated string's collation apart, as the sync states them.
func sdlType(ref tree.ResolvableTypeReference) (string, string) {
	t, ok := tree.GetStaticallyKnownType(ref)
	if !ok {
		return schemaQualified(ref).SQLString(), ""
	}
	if t.Family() == types.CollatedStringFamily {
		locale := t.Locale()
		return strings.TrimSuffix(t.SQLString(), " COLLATE "+tree.NameString(locale)), locale
	}
	return t.SQLString(), ""
}

// schemaQualified names a user-defined type by its schema, as the sync does. An
// unqualified name is the one CockroachDB's default search path finds, in
// public.
func schemaQualified(ref tree.ResolvableTypeReference) tree.ResolvableTypeReference {
	switch r := ref.(type) {
	case *tree.UnresolvedObjectName:
		if r.NumParts == 1 {
			if qualified, err := tree.NewUnresolvedObjectName(2, [3]string{r.Parts[0], "public"}, tree.NoAnnotation); err == nil {
				return qualified
			}
		}
	case *tree.ArrayTypeReference:
		return &tree.ArrayTypeReference{ElementType: schemaQualified(r.ElementType)}
	}
	return ref
}

// sdlExpr formats an expression without its outer parentheses.
func sdlExpr(expr tree.Expr) string {
	for {
		paren, ok := expr.(*tree.ParenExpr)
		if !ok {
			break
		}
		expr = paren.Expr
	}
	return tree.AsString(expr)
}

func (l *sdlLoader) createTable(n *tree.CreateTable) error {
	if n.AsSource != nil {
		return fmt.Errorf("CREATE TABLE AS is not a schema definition")
	}
	if n.Persistence.IsTemporary() {
		return fmt.Errorf("table %s: a temporary table is not a schema definition", n.Table.String())
	}
	if n.Locality != nil {
		// The sync reads no locality, so a migration could not keep one.
		return fmt.Errorf("table %s: LOCALITY is not supported", n.Table.String())
	}
	// One column family, the default, is the only layout the sync reads.
	families := make(map[string]bool)
	for _, def := range n.Defs {
		switch d := def.(type) {
		case *tree.FamilyTableDef:
			families[string(d.Name)] = true
		case *tree.ColumnTableDef:
			if d.HasColumnFamily() {
				families[string(d.Family.Name)] = true
			}
		}
	}
	if len(families) > 1 {
		return fmt.Errorf("table %s: more than one column family is not supported", n.Table.String())
	}
	if n.PartitionByTable != nil {
		// Nor partitioning.
		return fmt.Errorf("table %s: PARTITION BY is not supported", n.Table.String())
	}
	if len(n.StorageParams) > 0 {
		// Nor storage parameters, row-level TTL among them.
		return fmt.Errorf("table %s: storage parameters are not supported", n.Table.String())
	}
	schema := schemaOf(&n.Table)
	s := l.schema(schema)
	name := string(n.Table.ObjectName)
	exists := schemaTable(s, name) != nil || findView(s, name) != nil || findMaterializedView(s, name) != nil
	if skip, err := repeated(exists, n.IfNotExists, n.Table.String()); skip || err != nil {
		return err
	}
	table := &metadatapb.TableMetadata{Name: string(n.Table.ObjectName)}
	s.Tables = append(s.Tables, table)
	// CockroachDB names a table's own checks and foreign keys before its
	// columns' ones.
	for _, def := range n.Defs {
		switch def.(type) {
		case *tree.CheckConstraintTableDef, *tree.ForeignKeyConstraintTableDef:
			if err := l.addTableDef(schema, table, def); err != nil {
				return err
			}
		}
	}
	for _, def := range n.Defs {
		if column, ok := def.(*tree.ColumnTableDef); ok {
			if err := l.addColumn(schema, table, column); err != nil {
				return err
			}
		}
	}
	for _, def := range n.Defs {
		switch def.(type) {
		case *tree.CheckConstraintTableDef, *tree.ForeignKeyConstraintTableDef:
			continue
		}
		if err := l.addTableDef(schema, table, def); err != nil {
			return err
		}
	}
	return nil
}

// addColumn adds a column a CREATE TABLE or ALTER TABLE ADD COLUMN declares,
// with the keys, checks and foreign key it declares inline.
func (l *sdlLoader) addColumn(schema string, table *metadatapb.TableMetadata, column *tree.ColumnTableDef) error {
	if findColumn(table, string(column.Name)) != nil {
		return fmt.Errorf("column %s is added twice", column.Name)
	}
	typ, collation := sdlType(column.Type)
	c := &metadatapb.ColumnMetadata{
		Name:        string(column.Name),
		Position:    int32(len(table.Columns) + 1),
		Type:        typ,
		Collation:   collation,
		Nullable:    column.Nullable.Nullability != tree.NotNull && !column.PrimaryKey.IsPrimaryKey,
		IsInvisible: column.Hidden,
	}
	if column.DefaultExpr.Expr != nil {
		c.Default = sdlExpr(column.DefaultExpr.Expr)
	}
	if column.GeneratedIdentity.IsGeneratedAsIdentity {
		// The sync reads no identity, so a migration could not keep one.
		return fmt.Errorf("column %s: an identity column is not supported", column.Name)
	}
	if column.OnUpdateExpr.Expr != nil {
		// The sync reads no ON UPDATE, so a migration could not keep one.
		return fmt.Errorf("column %s: ON UPDATE is not supported", column.Name)
	}
	if column.IsSerial {
		c.Default = "unique_rowid()"
	}
	if column.Computed.Computed {
		generation := &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_STORED, Expression: sdlExpr(column.Computed.Expr)}
		if column.Computed.Virtual {
			generation.Type = metadatapb.GenerationMetadata_TYPE_VIRTUAL
		}
		c.Generation = generation
	}
	table.Columns = append(table.Columns, c)
	if column.PrimaryKey.IsPrimaryKey {
		index := &tree.CreateIndex{Columns: tree.IndexElemList{{Column: column.Name}}, StorageParams: column.PrimaryKey.StorageParams}
		if column.PrimaryKey.Sharded {
			index.Sharded = &tree.ShardedIndexDef{ShardBuckets: column.PrimaryKey.ShardBuckets}
		}
		// An inline primary key's CONSTRAINT name is the parser's
		// Unique.ConstraintName.
		addPrimaryKey(schema, table, string(column.Unique.ConstraintName), index)
	}
	if column.Unique.WithoutIndex {
		return fmt.Errorf("column %s: UNIQUE WITHOUT INDEX is not supported", column.Name)
	}
	if column.Unique.IsUnique {
		index := &tree.CreateIndex{Unique: true, Columns: tree.IndexElemList{{Column: column.Name}}, Name: column.Unique.ConstraintName}
		table.Indexes = append(table.Indexes, sdlIndex(schema, table.Name, table.Indexes, index))
	}
	for _, check := range column.CheckExprs {
		table.CheckConstraints = append(table.CheckConstraints, sdlCheck(table, string(check.ConstraintName), check.Expr))
	}
	if column.References.Table != nil {
		fk := &tree.ForeignKeyConstraintTableDef{
			Name: column.References.ConstraintName, Table: *column.References.Table,
			FromCols: tree.NameList{column.Name}, Actions: column.References.Actions, Match: column.References.Match,
		}
		if column.References.Col != "" {
			fk.ToCols = tree.NameList{column.References.Col}
		}
		if err := l.addForeignKey(schema, table, fk); err != nil {
			return err
		}
	}
	return nil
}

// addTableDef adds a constraint or index a CREATE TABLE or ALTER TABLE ADD
// CONSTRAINT declares.
func (l *sdlLoader) addTableDef(schema string, table *metadatapb.TableMetadata, def tree.TableDef) error {
	switch d := def.(type) {
	case *tree.ColumnTableDef, *tree.FamilyTableDef:
	case *tree.UniqueConstraintTableDef:
		if d.WithoutIndex {
			// The sync reads indexes; a constraint without one is lost.
			return fmt.Errorf("UNIQUE WITHOUT INDEX is not supported")
		}
		index := indexOf(&d.IndexTableDef)
		if d.PrimaryKey {
			addPrimaryKey(schema, table, string(d.Name), index)
			return nil
		}
		index.Unique = true
		table.Indexes = append(table.Indexes, sdlIndex(schema, table.Name, table.Indexes, index))
	case *tree.IndexTableDef:
		table.Indexes = append(table.Indexes, sdlIndex(schema, table.Name, table.Indexes, indexOf(d)))
	case *tree.CheckConstraintTableDef:
		table.CheckConstraints = append(table.CheckConstraints, sdlCheck(table, string(d.Name), d.Expr))
	case *tree.ForeignKeyConstraintTableDef:
		return l.addForeignKey(schema, table, d)
	default:
		return fmt.Errorf("table element %s is not supported", tree.AsString(def))
	}
	return nil
}

func indexOf(d *tree.IndexTableDef) *tree.CreateIndex {
	return &tree.CreateIndex{
		Name: d.Name, Columns: d.Columns, Sharded: d.Sharded, Storing: d.Storing, Type: d.Type,
		PartitionByIndex: d.PartitionByIndex, StorageParams: d.StorageParams, Predicate: d.Predicate, Invisibility: d.Invisibility,
	}
}

// addPrimaryKey adds a primary key, named as CockroachDB names it when the
// SDL does not.
func addPrimaryKey(schema string, table *metadatapb.TableMetadata, name string, index *tree.CreateIndex) {
	if name == "" {
		name = table.Name + "_pkey"
	}
	index.Name = tree.Name(name)
	index.Unique = true
	key := sdlIndex(schema, table.Name, table.Indexes, index)
	key.Primary = true
	table.Indexes = append([]*metadatapb.IndexMetadata{key}, table.Indexes...)
	for _, elem := range index.Columns {
		if column := findColumn(table, string(elem.Column)); column != nil {
			column.Nullable = false
		}
	}
}

// sdlIndex states an index as the sync does: its key list, and its definition
// as CockroachDB formats it, with each key's direction written. An unnamed
// index is named as CockroachDB names it: the table, its keys and _idx or
// _key, then 1, 2 and on until no index in existing has the name.
func sdlIndex(schema, table string, existing []*metadatapb.IndexMetadata, n *tree.CreateIndex) *metadatapb.IndexMetadata {
	index := &tree.CreateIndex{}
	*index = *n
	index.Table = tree.MakeTableNameWithSchema("", tree.Name(schema), tree.Name(table))
	index.Table.ExplicitCatalog = false
	index.IfNotExists = false
	index.Concurrently = false
	index.Columns = make(tree.IndexElemList, len(n.Columns))
	copy(index.Columns, n.Columns)
	var names []string
	metadata := &metadatapb.IndexMetadata{Unique: index.Unique}
	for i := range index.Columns {
		elem := &index.Columns[i]
		if elem.Direction == tree.DefaultDirection {
			elem.Direction = tree.Ascending
		}
		key := string(elem.Column)
		if elem.Expr != nil {
			key = sdlExpr(elem.Expr)
			names = append(names, "expr")
		} else {
			names = append(names, string(elem.Column))
		}
		metadata.Expressions = append(metadata.Expressions, key)
		metadata.Descending = append(metadata.Descending, elem.Direction == tree.Descending)
	}
	if index.Name == "" {
		suffix := "_idx"
		if index.Unique {
			suffix = "_key"
		}
		base := table + "_" + strings.Join(names, "_") + suffix
		name := base
		for i := 1; slices.ContainsFunc(existing, func(x *metadatapb.IndexMetadata) bool { return x.Name == name }); i++ {
			name = base + strconv.Itoa(i)
		}
		index.Name = tree.Name(name)
	}
	metadata.Name = string(index.Name)
	metadata.Definition = tree.AsString(index)
	return metadata
}

// sdlCheck states a check as the sync does, parenthesized, named as
// CockroachDB names it when the SDL does not: check_ and the columns it
// reads, then 1, 2 and on until no other constraint of the table has the
// name.
func sdlCheck(table *metadatapb.TableMetadata, name string, expr tree.Expr) *metadatapb.CheckConstraintMetadata {
	if name == "" {
		var columns []string
		seen := make(map[string]bool)
		tree.WalkExprConst(&columnCollector{visit: func(column string) {
			if !seen[column] {
				seen[column] = true
				columns = append(columns, column)
			}
		}}, expr)
		base := strings.Join(append([]string{"check"}, columns...), "_")
		name = base
		for i := 1; constraintTaken(table, name); i++ {
			name = base + strconv.Itoa(i)
		}
	}
	return &metadatapb.CheckConstraintMetadata{Name: name, Expression: "((" + sdlExpr(expr) + "))"}
}

// constraintTaken reports whether a check or a foreign key of a table has a
// name.
func constraintTaken(table *metadatapb.TableMetadata, name string) bool {
	return slices.ContainsFunc(table.CheckConstraints, func(c *metadatapb.CheckConstraintMetadata) bool { return c.Name == name }) ||
		slices.ContainsFunc(table.ForeignKeys, func(fk *metadatapb.ForeignKeyMetadata) bool { return fk.Name == name })
}

type columnCollector struct{ visit func(string) }

func (v *columnCollector) VisitPre(expr tree.Expr) (bool, tree.Expr) {
	if name, ok := expr.(*tree.UnresolvedName); ok && name.NumParts == 1 && !name.Star {
		v.visit(name.Parts[0])
	}
	return true, expr
}

func (*columnCollector) VisitPost(expr tree.Expr) tree.Expr { return expr }

func (l *sdlLoader) addForeignKey(schema string, table *metadatapb.TableMetadata, d *tree.ForeignKeyConstraintTableDef) error {
	fk := &metadatapb.ForeignKeyMetadata{
		Name:             string(d.Name),
		ReferencedSchema: schemaOf(&d.Table),
		ReferencedTable:  string(d.Table.ObjectName),
		OnDelete:         referenceAction(d.Actions.Delete),
		OnUpdate:         referenceAction(d.Actions.Update),
		MatchType:        "SIMPLE",
	}
	if d.Match == tree.MatchFull {
		fk.MatchType = "FULL"
	}
	for _, column := range d.FromCols {
		fk.Columns = append(fk.Columns, string(column))
	}
	for _, column := range d.ToCols {
		fk.ReferencedColumns = append(fk.ReferencedColumns, string(column))
	}
	if len(fk.ReferencedColumns) == 0 {
		// REFERENCES without columns names the referenced primary key, which
		// SDL may declare later.
		l.keyless = append(l.keyless, fk)
	}
	if fk.Name == "" {
		// CockroachDB names it the table, its columns and _fkey, then _1,
		// _2 and on until no constraint of the table has the name.
		base := table.Name + "_" + strings.Join(fk.Columns, "_") + "_fkey"
		fk.Name = base
		for i := 1; constraintTaken(table, fk.Name); i++ {
			fk.Name = base + "_" + strconv.Itoa(i)
		}
	}
	table.ForeignKeys = append(table.ForeignKeys, fk)
	return nil
}

// referenceAction states a foreign key action as the sync does.
func referenceAction(action tree.ReferenceAction) string {
	if action == tree.NoAction {
		return "NO ACTION"
	}
	return action.String()
}

func (l *sdlLoader) alterTable(n *tree.AlterTable) error {
	schema := n.Table.Schema()
	if schema == "" {
		schema = "public"
	}
	table := schemaTable(l.schema(schema), n.Table.Object())
	if table == nil {
		return fmt.Errorf("table %s is not created", n.Table)
	}
	for _, cmd := range n.Cmds {
		switch c := cmd.(type) {
		case *tree.AlterTableAddConstraint:
			if c.ValidationBehavior == tree.ValidationSkip {
				// The sync reads no validation state, so a migration
				// would validate it.
				return fmt.Errorf("NOT VALID is not supported")
			}
			if err := l.addTableDef(schema, table, c.ConstraintDef); err != nil {
				return err
			}
		case *tree.AlterTableAddColumn:
			if c.IfNotExists && findColumn(table, string(c.ColumnDef.Name)) != nil {
				// The engine skips it.
				continue
			}
			if err := l.addColumn(schema, table, c.ColumnDef); err != nil {
				return err
			}
		case *tree.AlterTableValidateConstraint:
			// A dump validates the foreign keys it added; SDL holds them valid.
		case *tree.AlterTableSetDefault:
			column := findColumn(table, string(c.Column))
			if column == nil {
				return fmt.Errorf("column %s is not created", c.Column)
			}
			column.Default = ""
			if c.Default != nil {
				column.Default = sdlExpr(c.Default)
			}
		default:
			return fmt.Errorf("ALTER TABLE %s is not a schema definition", tree.AsString(cmd))
		}
	}
	return nil
}

func (l *sdlLoader) createView(n *tree.CreateView) error {
	if n.Persistence.IsTemporary() {
		return fmt.Errorf("view %s: a temporary view is not a schema definition", n.Name.String())
	}
	schema := schemaOf(&n.Name)
	s := l.schema(schema)
	if !n.Replace {
		name := string(n.Name.ObjectName)
		exists := findView(s, name) != nil || findMaterializedView(s, name) != nil || schemaTable(s, name) != nil
		if skip, err := repeated(exists, n.IfNotExists, n.Name.String()); skip || err != nil {
			return err
		}
	}
	dropCatalogs(n.AsSource)
	definition := tree.AsString(n.AsSource)
	var dependencies []*metadatapb.DependencyColumn
	seen := make(map[string]bool)
	storedRelations(n.AsSource, func(tn *tree.TableName) {
		dep := &metadatapb.DependencyColumn{Schema: schemaOf(tn), Table: string(tn.ObjectName)}
		if id := getObjectID(dep.Schema, dep.Table); !seen[id] {
			seen[id] = true
			dependencies = append(dependencies, dep)
		}
	})
	names := make([]string, len(n.ColumnNames))
	for i, name := range n.ColumnNames {
		names[i] = string(name)
	}
	if len(names) == 0 {
		names = outputNames(n.AsSource)
	}
	var columns []*metadatapb.ColumnMetadata
	for _, name := range names {
		columns = append(columns, &metadatapb.ColumnMetadata{Name: name})
	}
	if n.Materialized {
		view := &metadatapb.MaterializedViewMetadata{Name: string(n.Name.ObjectName), Definition: definition, DependencyColumns: dependencies, Columns: columns}
		s.MaterializedViews = append(s.MaterializedViews, view)
		if names == nil {
			l.starViews = append(l.starViews, starView{&view.Columns, n.AsSource})
		}
		return nil
	}
	view := &metadatapb.ViewMetadata{Name: string(n.Name.ObjectName), Definition: definition, DependencyColumns: dependencies, Columns: columns}
	if names == nil {
		l.starViews = append(l.starViews, starView{&view.Columns, n.AsSource})
	}
	// CREATE OR REPLACE VIEW replaces a view the SDL created before.
	if i := slices.IndexFunc(s.Views, func(v *metadatapb.ViewMetadata) bool { return v.Name == view.Name }); i >= 0 && n.Replace {
		s.Views[i] = view
		return nil
	}
	s.Views = append(s.Views, view)
	return nil
}

// outputNames returns the names CockroachDB gives a query's output columns,
// or nil when the query has a form whose names it does not know (a star).
func outputNames(query *tree.Select) []string {
	var clause tree.SelectStatement = query.Select
	for {
		switch c := clause.(type) {
		case *tree.ParenSelect:
			clause = c.Select.Select
			continue
		case *tree.UnionClause:
			clause = c.Left.Select
			continue
		}
		break
	}
	sel, ok := clause.(*tree.SelectClause)
	if !ok {
		return nil
	}
	names := make([]string, len(sel.Exprs))
	for i, expr := range sel.Exprs {
		if expr.As != "" {
			names[i] = string(expr.As)
			continue
		}
		name, ok := outputName(expr.Expr)
		if !ok {
			return nil
		}
		names[i] = name
	}
	return names
}

// outputName is the name CockroachDB gives an output column it is not told
// one for, as v25.2.3 names them.
func outputName(expr tree.Expr) (string, bool) {
	switch e := expr.(type) {
	case *tree.UnresolvedName:
		if e.Star {
			return "", false
		}
		return e.Parts[0], true
	case *tree.ColumnItem:
		return string(e.ColumnName), true
	case *tree.ParenExpr:
		return outputName(e.Expr)
	case *tree.CastExpr:
		// A cast keeps a column's name; a cast constant is named after its
		// type, which is left unknown here.
		return columnName(e.Expr)
	case *tree.AnnotateTypeExpr:
		return columnName(e.Expr)
	case *tree.FuncExpr:
		if e.Func.FunctionReference == nil {
			return "", false
		}
		name, ok := e.Func.FunctionReference.(*tree.UnresolvedName)
		if !ok {
			return "", false
		}
		return strings.ToLower(name.Parts[0]), true
	case *tree.CaseExpr:
		return "case", true
	case *tree.DBool:
		return "bool", true
	case *tree.StrVal, *tree.NumVal, *tree.BinaryExpr, *tree.UnaryExpr:
		return "?column?", true
	}
	return "", false
}

// columnName is the name of an expression that is a column, in parentheses
// or not.
func columnName(expr tree.Expr) (string, bool) {
	switch e := expr.(type) {
	case *tree.UnresolvedName, *tree.ColumnItem:
		return outputName(e)
	case *tree.ParenExpr:
		return columnName(e.Expr)
	}
	return "", false
}

func (l *sdlLoader) createRoutine(n *tree.CreateRoutine) error {
	replace := n.Replace
	n.Replace = false
	for i, option := range n.Options {
		if body, ok := option.(tree.RoutineBodyStr); ok {
			n.Options[i] = tree.RoutineBodyStr(bodyWithoutCatalogs(string(body), n.Options))
		}
	}
	schema := "public"
	if n.Name.ExplicitSchema {
		schema = string(n.Name.SchemaName)
	}
	// A user-defined type is named by its schema, as in a table's columns.
	for i := range n.Params {
		if _, ok := tree.GetStaticallyKnownType(n.Params[i].Type); !ok {
			n.Params[i].Type = schemaQualified(n.Params[i].Type)
		}
	}
	if n.ReturnType != nil {
		if _, ok := tree.GetStaticallyKnownType(n.ReturnType.Type); !ok {
			n.ReturnType.Type = schemaQualified(n.ReturnType.Type)
		}
	}
	var args []string
	for _, param := range n.Params {
		if param.Class == tree.RoutineParamOut {
			continue
		}
		args = append(args, param.Type.SQLString())
	}
	name := string(n.Name.ObjectName)
	function := &metadatapb.FunctionMetadata{
		Name:       name,
		Signature:  name + "(" + strings.Join(args, ", ") + ")",
		Definition: tree.AsString(n),
	}
	s := l.schema(schema)
	// CREATE OR REPLACE replaces a routine of the same signature the SDL
	// created before.
	if i := slices.IndexFunc(s.Functions, func(f *metadatapb.FunctionMetadata) bool { return f.Signature == function.Signature }); i >= 0 {
		if !replace {
			return fmt.Errorf("%s is created twice", function.Signature)
		}
		s.Functions[i] = function
		return nil
	}
	s.Functions = append(s.Functions, function)
	return nil
}

func (l *sdlLoader) relationIndexes(schema, name string) *[]*metadatapb.IndexMetadata {
	if table := schemaTable(l.schema(schema), name); table != nil {
		return &table.Indexes
	}
	if view := findMaterializedView(l.schema(schema), name); view != nil {
		return &view.Indexes
	}
	return nil
}

func (l *sdlLoader) findIndex(schema, table, name string) *metadatapb.IndexMetadata {
	s := l.schema(schema)
	var indexes []*metadatapb.IndexMetadata
	for _, t := range s.Tables {
		if table == "" || t.Name == table {
			indexes = append(indexes, t.Indexes...)
		}
	}
	for _, v := range s.MaterializedViews {
		if table == "" || v.Name == table {
			indexes = append(indexes, v.Indexes...)
		}
	}
	for _, index := range indexes {
		if index.Name == name {
			return index
		}
	}
	return nil
}

func (l *sdlLoader) commentOnConstraint(n *tree.CommentOnConstraint) error {
	table := schemaTable(l.schema(n.Table.Schema()), n.Table.Object())
	if table == nil {
		return fmt.Errorf("table %s is not created", n.Table)
	}
	name := string(n.Constraint)
	for _, check := range table.CheckConstraints {
		if check.Name == name {
			check.Comment = comment(n.Comment)
			return nil
		}
	}
	for _, fk := range table.ForeignKeys {
		if fk.Name == name {
			fk.Comment = comment(n.Comment)
			return nil
		}
	}
	for _, index := range table.Indexes {
		if index.Name == name {
			index.Comment = comment(n.Comment)
			return nil
		}
	}
	return fmt.Errorf("constraint %s is not created", name)
}

// dropCatalogs unqualifies every table name in node that names a database.
// SDL describes one database, and a dump (SHOW CREATE) names it in every
// relation a view reads, so the name is dropped as DatabaseDefinition drops it.
func dropCatalogs(node tree.NodeFormatter) bool {
	changed := false
	ctx := tree.NewFmtCtx(tree.FmtSimple, tree.FmtReformatTableNames(func(ctx *tree.FmtCtx, tn *tree.TableName) {
		if tn.ExplicitCatalog {
			tn.ExplicitCatalog = false
			tn.CatalogName = ""
			changed = true
		}
		ctx.WithReformatTableNames(nil, func() { ctx.FormatNode(tn) })
	}))
	ctx.FormatNode(node)
	ctx.Close()
	return changed
}

// bodyWithoutCatalogs is a routine body with dropCatalogs applied, or the body
// as written when it names no database or does not parse.
func bodyWithoutCatalogs(body string, options tree.RoutineOptions) string {
	language := ""
	for _, option := range options {
		if l, ok := option.(tree.RoutineLanguage); ok {
			language = strings.ToLower(string(l))
		}
	}
	switch language {
	case "sql":
		stmts, err := crparser.Parse(body)
		if err != nil {
			return body
		}
		changed := false
		parts := make([]string, len(stmts))
		for i, s := range stmts {
			changed = dropCatalogs(s.AST) || changed
			parts[i] = tree.AsString(s.AST) + ";"
		}
		if changed {
			return "\n" + strings.Join(parts, "\n") + "\n"
		}
	case "plpgsql":
		block, err := plpgsqlparser.Parse(body)
		if err == nil && dropCatalogs(block.AST) {
			return "\n" + tree.AsString(block.AST) + "\n"
		}
	}
	return body
}
