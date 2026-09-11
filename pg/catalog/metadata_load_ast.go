package catalog

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bytebase/omni/metadata"
	nodes "github.com/bytebase/omni/pg/ast"
	pgparser "github.com/bytebase/omni/pg/parser"
)

// The snapshot carries types, defaults, and bodies as the text PostgreSQL
// reports for them. The builders below turn that text into the statement ASTs
// the Define* functions take.

func metadataEnumStmt(schema string, e *metadata.EnumTypeMetadata) *nodes.CreateEnumStmt {
	vals := make([]nodes.Node, 0, len(e.GetValues()))
	for _, v := range e.GetValues() {
		vals = append(vals, &nodes.String{Str: v})
	}
	return &nodes.CreateEnumStmt{TypeName: metadataNameList(schema, e.GetName()), Vals: &nodes.List{Items: vals}}
}

func metadataCompositeTypeStmt(schema string, ct *metadata.CompositeTypeMetadata) (*nodes.CompositeTypeStmt, error) {
	cols := make([]nodes.Node, 0, len(ct.GetAttributes()))
	for _, a := range ct.GetAttributes() {
		if a.GetName() == "" {
			continue
		}
		tn, err := parseTypeName(a.GetType())
		if err != nil {
			return nil, fmt.Errorf("attribute %q: %w", a.GetName(), err)
		}
		cols = append(cols, &nodes.ColumnDef{
			Colname:    a.GetName(),
			TypeName:   tn,
			CollClause: collateClause(a.GetCollation()),
		})
	}
	return &nodes.CompositeTypeStmt{
		Typevar:    &nodes.RangeVar{Schemaname: schema, Relname: ct.GetName()},
		Coldeflist: &nodes.List{Items: cols},
	}, nil
}

// metadataSequenceStmt treats an empty or zero option as unset; sync reports
// zero for options the sequence leaves at their defaults.
func metadataSequenceStmt(schema string, seq *metadata.SequenceMetadata) *nodes.CreateSeqStmt {
	var opts []nodes.Node
	for _, o := range []struct{ name, value string }{
		{"increment", seq.GetIncrement()},
		{"minvalue", seq.GetMinValue()},
		{"maxvalue", seq.GetMaxValue()},
		{"start", seq.GetStart()},
		{"cache", seq.GetCacheSize()},
	} {
		if o.value == "" || o.value == "0" {
			continue
		}
		v, _ := strconv.ParseInt(o.value, 10, 64)
		opts = append(opts, &nodes.DefElem{Defname: o.name, Arg: &nodes.Integer{Ival: v}})
	}
	if seq.GetCycle() {
		opts = append(opts, &nodes.DefElem{Defname: "cycle", Arg: &nodes.Boolean{Boolval: true}})
	}
	stmt := &nodes.CreateSeqStmt{
		Sequence: &nodes.RangeVar{Schemaname: schema, Relname: seq.GetName(), Relpersistence: 'p'},
	}
	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}
	return stmt
}

// metadataTableStmt fails on any column whose type does not parse, so the
// table gets a stand-in rather than a partial definition. With full, a default
// or generation expression that does not parse is dropped instead.
func metadataTableStmt(schema string, t *metadata.TableMetadata, full bool) (*nodes.CreateStmt, error) {
	items := make([]nodes.Node, 0, len(t.GetColumns()))
	for _, col := range t.GetColumns() {
		if col.GetName() == "" {
			continue
		}
		if col.GetType() == "" {
			return nil, fmt.Errorf("column %q: empty type", col.GetName())
		}
		tn, err := parseTypeName(col.GetType())
		if err != nil {
			return nil, fmt.Errorf("column %q: %w", col.GetName(), err)
		}
		def := &nodes.ColumnDef{Colname: col.GetName(), TypeName: tn, IsNotNull: !col.GetNullable()}
		if full {
			setColumnValueProperties(def, col)
		}
		items = append(items, def)
	}
	return &nodes.CreateStmt{
		Relation:  &nodes.RangeVar{Schemaname: schema, Relname: t.GetName(), Relpersistence: 'p'},
		TableElts: &nodes.List{Items: items},
	}, nil
}

func setColumnValueProperties(def *nodes.ColumnDef, col *metadata.ColumnMetadata) {
	if col.GetDefault() != "" && col.GetGeneration() == nil {
		if expr, err := parseScalar("SELECT " + col.GetDefault()); err == nil {
			def.RawDefault = expr
		}
	}
	if gen := col.GetGeneration(); gen.GetType() == metadata.GenerationMetadata_TYPE_STORED && gen.GetExpression() != "" {
		def.Generated = 's'
		if expr, err := parseScalar("SELECT (" + gen.GetExpression() + ")"); err == nil {
			def.Constraints = &nodes.List{Items: []nodes.Node{
				&nodes.Constraint{Contype: nodes.CONSTR_GENERATED, RawExpr: expr},
			}}
		}
	}
	if col.GetIsIdentity() {
		def.Identity = 'd'
		if col.GetIdentityGeneration() == metadata.ColumnMetadata_ALWAYS {
			def.Identity = 'a'
		}
	}
}

func metadataViewStmt(schema string, v *metadata.ViewMetadata) (*nodes.ViewStmt, error) {
	if v.GetDefinition() == "" {
		return nil, errors.New("empty view definition")
	}
	sel, err := parseSelect(v.GetDefinition())
	if err != nil {
		return nil, fmt.Errorf("view %q: %w", v.GetName(), err)
	}
	return &nodes.ViewStmt{
		View:  &nodes.RangeVar{Schemaname: schema, Relname: v.GetName(), Relpersistence: 'p'},
		Query: sel,
	}, nil
}

func metadataMatViewStmt(schema string, mv *metadata.MaterializedViewMetadata) (*nodes.CreateTableAsStmt, error) {
	if mv.GetDefinition() == "" {
		return nil, errors.New("empty materialized view definition")
	}
	sel, err := parseSelect(mv.GetDefinition())
	if err != nil {
		return nil, fmt.Errorf("materialized view %q: %w", mv.GetName(), err)
	}
	return matViewStmt(schema, mv.GetName(), sel), nil
}

func matViewStmt(schema, name string, sel *nodes.SelectStmt) *nodes.CreateTableAsStmt {
	return &nodes.CreateTableAsStmt{
		Query:   sel,
		Objtype: nodes.OBJECT_MATVIEW,
		Into: &nodes.IntoClause{
			Rel: &nodes.RangeVar{Schemaname: schema, Relname: name, Relpersistence: 'p'},
		},
	}
}

// metadataFunctionStmt parses the function's CREATE statement. When the
// snapshot has none that parses, it builds a SQL function from the signature,
// with the argument types and a text result, so calls by name still resolve.
func metadataFunctionStmt(schema string, fn *metadata.FunctionMetadata) (*nodes.CreateFunctionStmt, error) {
	if fn.GetDefinition() != "" {
		if list, err := pgparser.Parse(fn.GetDefinition()); err == nil && list != nil && len(list.Items) == 1 {
			if stmt, ok := unwrapRawStmt(list.Items[0]).(*nodes.CreateFunctionStmt); ok {
				return stmt, nil
			}
		}
	}
	argTypes, err := signatureArgTypes(fn.GetSignature())
	if err != nil {
		return nil, fmt.Errorf("function %q: %w", fn.GetName(), err)
	}
	stmt := &nodes.CreateFunctionStmt{
		Funcname:   metadataNameList(schema, fn.GetName()),
		ReturnType: textTypeName(),
		Options:    sqlFunctionBody("SELECT NULL::text"),
	}
	if len(argTypes) > 0 {
		params := make([]nodes.Node, 0, len(argTypes))
		for i, t := range argTypes {
			tn, err := parseTypeName(t)
			if err != nil {
				return nil, fmt.Errorf("function %q argument %d: %w", fn.GetName(), i+1, err)
			}
			params = append(params, &nodes.FunctionParameter{ArgType: tn, Mode: nodes.FUNC_PARAM_IN})
		}
		stmt.Parameters = &nodes.List{Items: params}
	}
	return stmt, nil
}

func sqlFunctionBody(body string) *nodes.List {
	return &nodes.List{Items: []nodes.Node{
		&nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: "sql"}},
		&nodes.DefElem{Defname: "as", Arg: &nodes.List{Items: []nodes.Node{&nodes.String{Str: body}}}},
	}}
}

func metadataKeyConstraint(idx *metadata.IndexMetadata) ConstraintDef {
	def := ConstraintDef{Name: idx.GetName(), Type: ConstraintUnique}
	if idx.GetPrimary() {
		def.Type = ConstraintPK
	}
	for _, expr := range idx.GetExpressions() {
		def.Columns = append(def.Columns, unquoteIdent(expr))
	}
	return def
}

// metadataIndexStmt reads each key as a column name unless it looks like an
// expression, in pg_get_indexdef's shape: a function call or a parenthesized,
// cast, or spaced expression.
func metadataIndexStmt(schema, relation string, idx *metadata.IndexMetadata) (*nodes.IndexStmt, error) {
	if idx.GetName() == "" {
		return nil, errors.New("index has no name")
	}
	params := make([]nodes.Node, 0, len(idx.GetExpressions()))
	for i, expr := range idx.GetExpressions() {
		if expr == "" {
			continue
		}
		elem := &nodes.IndexElem{}
		if i < len(idx.GetDescending()) && idx.GetDescending()[i] {
			elem.Ordering = nodes.SORTBY_DESC
		}
		name := unquoteIdent(expr)
		if strings.ContainsAny(name, "( ") || strings.Contains(name, "::") {
			if parsed, err := parseScalar("SELECT (" + expr + ")"); err == nil {
				elem.Expr = parsed
			} else {
				elem.Name = name
			}
		} else {
			elem.Name = name
		}
		params = append(params, elem)
	}
	return &nodes.IndexStmt{
		Idxname:      idx.GetName(),
		Relation:     &nodes.RangeVar{Schemaname: schema, Relname: relation},
		AccessMethod: idx.GetType(),
		IndexParams:  &nodes.List{Items: params},
		Unique:       idx.GetUnique(),
		Primary:      idx.GetPrimary(),
	}, nil
}

func metadataConstraintDef(o *metadataObject) ConstraintDef {
	switch {
	case o.fk != nil:
		refSchema := o.fk.GetReferencedSchema()
		if refSchema == "" {
			refSchema = o.schema
		}
		return ConstraintDef{
			Name:        o.fk.GetName(),
			Type:        ConstraintFK,
			Columns:     o.fk.GetColumns(),
			RefSchema:   refSchema,
			RefTable:    o.fk.GetReferencedTable(),
			RefColumns:  o.fk.GetReferencedColumns(),
			FKUpdAction: fkActionCode(o.fk.GetOnUpdate()),
			FKDelAction: fkActionCode(o.fk.GetOnDelete()),
			FKMatchType: fkMatchCode(o.fk.GetMatchType()),
		}
	case o.check != nil:
		return ConstraintDef{Name: o.check.GetName(), Type: ConstraintCheck, CheckExpr: o.check.GetExpression()}
	default:
		return ConstraintDef{Name: o.exclude.GetName(), Type: ConstraintExclude, CheckExpr: o.exclude.GetExpression()}
	}
}

func fkActionCode(action string) byte {
	switch strings.ToUpper(action) {
	case "RESTRICT":
		return 'r'
	case "CASCADE":
		return 'c'
	case "SET NULL":
		return 'n'
	case "SET DEFAULT":
		return 'd'
	default:
		return 'a'
	}
}

func fkMatchCode(match string) byte {
	switch strings.ToUpper(match) {
	case "FULL":
		return 'f'
	case "PARTIAL":
		return 'p'
	default:
		return 's'
	}
}

// Stand-ins keep an object's name and column names and type every column
// text, so they depend on no user object and objects that use them resolve.

func standInCompositeTypeStmt(schema, name string, ct *metadata.CompositeTypeMetadata) *nodes.CompositeTypeStmt {
	var cols []nodes.Node
	for _, a := range ct.GetAttributes() {
		if a.GetName() != "" {
			cols = append(cols, &nodes.ColumnDef{Colname: a.GetName(), TypeName: textTypeName()})
		}
	}
	return &nodes.CompositeTypeStmt{
		Typevar:    &nodes.RangeVar{Schemaname: schema, Relname: name},
		Coldeflist: &nodes.List{Items: cols},
	}
}

func standInTableStmt(schema, name string, columns []string) *nodes.CreateStmt {
	items := make([]nodes.Node, 0, len(columns))
	for _, col := range columns {
		items = append(items, &nodes.ColumnDef{Colname: col, TypeName: textTypeName()})
	}
	return &nodes.CreateStmt{
		Relation:  &nodes.RangeVar{Schemaname: schema, Relname: name, Relpersistence: 'p'},
		TableElts: &nodes.List{Items: items},
	}
}

func standInViewStmt(schema, name string, columns []string) (*nodes.ViewStmt, error) {
	sel, err := standInSelect(columns)
	if err != nil {
		return nil, err
	}
	return &nodes.ViewStmt{
		View:  &nodes.RangeVar{Schemaname: schema, Relname: name, Relpersistence: 'p'},
		Query: sel,
	}, nil
}

func standInMatViewStmt(schema, name string, columns []string) (*nodes.CreateTableAsStmt, error) {
	sel, err := standInSelect(columns)
	if err != nil {
		return nil, err
	}
	return matViewStmt(schema, name, sel), nil
}

// standInSelect returns SELECT NULL::text AS <column>, ... with one unnamed
// column when there are no names.
func standInSelect(columns []string) (*nodes.SelectStmt, error) {
	targets := make([]string, 0, len(columns))
	for _, col := range columns {
		targets = append(targets, "NULL::text AS "+quoteIdent(col))
	}
	if len(targets) == 0 {
		targets = []string{"NULL::text"}
	}
	return parseSelect("SELECT " + strings.Join(targets, ", "))
}

func tableColumnNames(t *metadata.TableMetadata) []string {
	names := make([]string, 0, len(t.GetColumns()))
	for _, col := range t.GetColumns() {
		names = appendUniqueName(names, col.GetName())
	}
	return names
}

// viewColumnNames prefers the view's own columns, which sync does not always
// report, over the columns it reads.
func viewColumnNames(v *metadata.ViewMetadata) []string {
	if len(v.GetColumns()) == 0 {
		return dependencyColumnNames(v.GetDependencyColumns())
	}
	names := make([]string, 0, len(v.GetColumns()))
	for _, col := range v.GetColumns() {
		names = appendUniqueName(names, col.GetName())
	}
	return names
}

func dependencyColumnNames(deps []*metadata.DependencyColumn) []string {
	var names []string
	for _, dep := range deps {
		names = appendUniqueName(names, dep.GetColumn())
	}
	return names
}

func appendUniqueName(names []string, name string) []string {
	for _, n := range names {
		if n == name {
			return names
		}
	}
	if name == "" {
		return names
	}
	return append(names, name)
}

// parseTypeName parses a type as PostgreSQL formats it, such as
// "character varying(255)" or "public.status[]".
func parseTypeName(typ string) (*nodes.TypeName, error) {
	expr, err := parseScalar("SELECT NULL::" + typ)
	if err != nil {
		return nil, fmt.Errorf("type %q: %w", typ, err)
	}
	cast, ok := expr.(*nodes.TypeCast)
	if !ok {
		return nil, fmt.Errorf("type %q: parsed as %T", typ, expr)
	}
	return cast.TypeName, nil
}

// parseScalar parses a single-target SELECT and returns its expression.
func parseScalar(sql string) (nodes.Node, error) {
	sel, err := parseSelect(sql)
	if err != nil {
		return nil, err
	}
	if sel.TargetList == nil || len(sel.TargetList.Items) != 1 {
		return nil, errors.New("expected one target")
	}
	rt, ok := sel.TargetList.Items[0].(*nodes.ResTarget)
	if !ok {
		return nil, fmt.Errorf("expected a target, got %T", sel.TargetList.Items[0])
	}
	return rt.Val, nil
}

// parseSelect parses exactly one SELECT statement.
func parseSelect(sql string) (*nodes.SelectStmt, error) {
	list, err := pgparser.Parse(sql)
	if err != nil {
		return nil, err
	}
	if list == nil || len(list.Items) != 1 {
		n := 0
		if list != nil {
			n = len(list.Items)
		}
		return nil, fmt.Errorf("expected one statement, got %d", n)
	}
	sel, ok := unwrapRawStmt(list.Items[0]).(*nodes.SelectStmt)
	if !ok {
		return nil, fmt.Errorf("expected SELECT, got %T", unwrapRawStmt(list.Items[0]))
	}
	return sel, nil
}

// userTypeRef returns the schema and name a type string refers to when it is
// schema-qualified, as sync reports user-defined types: "public.status",
// "public.status[]", or `"My Schema"."My Type"`.
func userTypeRef(typ string) (schema, name string, ok bool) {
	s := strings.TrimSpace(typ)
	for strings.HasSuffix(s, "[]") {
		s = strings.TrimSpace(strings.TrimSuffix(s, "[]"))
	}
	if before, _, found := strings.Cut(s, "("); found {
		s = strings.TrimSpace(before)
	}
	return splitQualifiedIdent(s)
}

var nextvalRef = regexp.MustCompile(`nextval\('([^']+)'`)

// nextvalSequence returns the sequence a nextval('...') default draws from,
// resolving an unqualified name to the table's schema.
func nextvalSequence(defaultExpr, tableSchema string) (schema, name string, ok bool) {
	m := nextvalRef.FindStringSubmatch(defaultExpr)
	if m == nil {
		return "", "", false
	}
	if schema, name, ok := splitQualifiedIdent(m[1]); ok {
		return schema, name, true
	}
	return tableSchema, unquoteIdent(m[1]), true
}

// signatureArgTypes splits the argument types out of a signature such as
// "fn(integer, numeric(10,2))".
func signatureArgTypes(signature string) ([]string, error) {
	open := strings.Index(signature, "(")
	if open < 0 {
		return nil, nil
	}
	closing := strings.LastIndex(signature, ")")
	if closing < open {
		return nil, fmt.Errorf("signature %q: unbalanced parentheses", signature)
	}
	inner := strings.TrimSpace(signature[open+1 : closing])
	if inner == "" {
		return nil, nil
	}
	var args []string
	depth, start := 0, 0
	for i := 0; i < len(inner); i++ {
		switch inner[i] {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				args = append(args, strings.TrimSpace(inner[start:i]))
				start = i + 1
			}
		default:
		}
	}
	return append(args, strings.TrimSpace(inner[start:])), nil
}

// splitQualifiedIdent splits "schema.name", where either part may be a quoted
// identifier with doubled quotes inside.
func splitQualifiedIdent(s string) (schema, name string, ok bool) {
	var parts []string
	for i := 0; ; {
		var part strings.Builder
		if i < len(s) && s[i] == '"' {
			closed := false
			for i++; i < len(s); i++ {
				if s[i] == '"' {
					if i+1 < len(s) && s[i+1] == '"' {
						part.WriteByte('"')
						i++
						continue
					}
					i++
					closed = true
					break
				}
				part.WriteByte(s[i])
			}
			if !closed {
				return "", "", false
			}
		} else {
			start := i
			for i < len(s) && s[i] != '.' {
				i++
			}
			part.WriteString(s[start:i])
		}
		parts = append(parts, part.String())
		if i >= len(s) {
			break
		}
		if s[i] != '.' {
			return "", "", false
		}
		i++
	}
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// collateClause turns a collation as sync stores it, quoted as needed and
// schema-qualified outside pg_catalog, into a COLLATE clause.
func collateClause(collation string) *nodes.CollateClause {
	if collation == "" {
		return nil
	}
	if schema, name, ok := splitQualifiedIdent(collation); ok {
		return &nodes.CollateClause{Collname: metadataNameList(schema, name)}
	}
	return &nodes.CollateClause{Collname: metadataNameList("", unquoteIdent(collation))}
}

func metadataNameList(schema, name string) *nodes.List {
	if schema == "" {
		return &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}}}
	}
	return &nodes.List{Items: []nodes.Node{&nodes.String{Str: schema}, &nodes.String{Str: name}}}
}

func textTypeName() *nodes.TypeName {
	return &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "text"}}}, Typemod: -1}
}

func quoteIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func unquoteIdent(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return strings.ReplaceAll(s[1:len(s)-1], `""`, `"`)
	}
	return s
}
