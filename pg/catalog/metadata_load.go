package catalog

import (
	"cmp"
	"container/heap"
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/bytebase/omni/metadata"
	nodes "github.com/bytebase/omni/pg/ast"
)

// LoadMetadataOptions selects what LoadMetadata installs.
type LoadMetadataOptions struct {
	// Full also installs sequences, procedures, column defaults, identity and
	// generated columns, indexes, and constraints: what DDL simulation checks
	// against. Query analysis needs only schemas, types, relations, and
	// functions, and loads faster without the rest.
	Full bool
}

// LoadMetadataReport lists the snapshot objects LoadMetadata could not install
// as the snapshot defines them. Keys are "schema:<schema>", "type:<schema>.<name>",
// "seq:<schema>.<name>", "rel:<schema>.<name>", "func:<schema>.<name>|<signature>",
// "idx:<schema>.<relation>.<name>", and "con:<schema>.<table>.<name>".
type LoadMetadataReport struct {
	// Degraded maps each object replaced by a text-typed stand-in to the error
	// its definition failed with.
	Degraded map[string]error
	// Missing maps each object the catalog has nothing for to the error that
	// left it out.
	Missing map[string]error
}

// LoadMetadata installs the objects of a Bytebase schema snapshot into c in
// dependency order. An object that fails to install does not stop the load:
// an enum, composite type, table, view, or materialized view is replaced by a
// stand-in with the same name and column names, all typed text, and a function
// by one built from its signature that returns text, so the objects that
// depend on it still install as defined. Other kinds are left out.
//
// View bodies are analyzed as they install, so set the search path and session
// user before loading. The error is non-nil only when ctx is done.
func (c *Catalog) LoadMetadata(ctx context.Context, meta *metadata.DatabaseSchemaMetadata, opts LoadMetadataOptions) (*LoadMetadataReport, error) {
	report := &LoadMetadataReport{
		Degraded: make(map[string]error),
		Missing:  make(map[string]error),
	}
	if err := ctx.Err(); err != nil {
		return report, err
	}
	for _, group := range orderMetadataObjects(collectMetadataObjects(meta, opts.Full)) {
		// A group is one object or the members of a dependency cycle, which can
		// be false: a view waits for every overload of a function it calls. Each
		// pass retries the members that failed, since they may be waiting on one
		// another. When a pass installs nothing, one member is replaced by its
		// stand-in: the first view, whose stand-in needs nothing, if any.
		for len(group) > 0 {
			var failed []*metadataObject
			var errs []error
			for _, o := range group {
				if err := ctx.Err(); err != nil {
					return report, err
				}
				if err := c.installMetadataObject(o, opts.Full); err != nil {
					failed = append(failed, o)
					errs = append(errs, err)
				}
			}
			if len(failed) == len(group) {
				i := max(slices.IndexFunc(failed, isViewObject), 0)
				c.installStandIn(failed[i], errs[i], report)
				failed = slices.Delete(failed, i, i+1)
			}
			group = failed
		}
	}
	return report, nil
}

// installStandIn replaces o, whose definition failed with err, by its stand-in
// and records the outcome in report.
func (c *Catalog) installStandIn(o *metadataObject, err error, report *LoadMetadataReport) {
	standIn := metadataStandIn(o)
	if standIn == nil {
		report.Missing[o.key()] = err
		return
	}
	if standInErr := standIn(c); standInErr != nil {
		report.Missing[o.key()] = standInErr
		return
	}
	report.Degraded[o.key()] = err
}

type metadataObjectKind int

const (
	metadataSchema metadataObjectKind = iota
	metadataEnum
	metadataComposite
	metadataSequence
	metadataTable
	metadataView
	metadataMatView
	metadataFunction
	metadataIndex
	metadataConstraint
)

// metadataObject is one install step. Every kind but a schema sets the metadata
// pointer matching it; a constraint sets one of fk, check, or exclude.
type metadataObject struct {
	kind   metadataObjectKind
	schema string
	name   string
	// parent is the relation an index or constraint belongs to.
	parent string

	enum      *metadata.EnumTypeMetadata
	composite *metadata.CompositeTypeMetadata
	sequence  *metadata.SequenceMetadata
	table     *metadata.TableMetadata
	view      *metadata.ViewMetadata
	matView   *metadata.MaterializedViewMetadata
	function  *metadata.FunctionMetadata
	index     *metadata.IndexMetadata
	fk        *metadata.ForeignKeyMetadata
	check     *metadata.CheckConstraintMetadata
	exclude   *metadata.ExcludeConstraintMetadata
	// procedure marks a function that comes from the snapshot's procedures.
	procedure bool
	// createStmt is a function's parsed definition, or nil when the snapshot
	// has none that parses.
	createStmt *nodes.CreateFunctionStmt
	// body is a view's parsed query, or nil with bodyErr saying why.
	body    *nodes.SelectStmt
	bodyErr error
	// ownedSequences maps a table's columns to the sequences they own, which
	// install only with Full.
	ownedSequences map[string]*metadata.SequenceMetadata
}

func (o *metadataObject) key() string {
	switch o.kind {
	case metadataSchema:
		return "schema:" + o.schema
	case metadataEnum, metadataComposite:
		return metadataTypeKey(o.schema, o.name)
	case metadataSequence:
		return metadataSequenceKey(o.schema, o.name)
	case metadataTable, metadataView, metadataMatView:
		return metadataRelationKey(o.schema, o.name)
	case metadataFunction:
		return "func:" + o.schema + "." + o.name + "|" + metadataFunctionIdentity(o.function)
	case metadataIndex:
		return "idx:" + o.schema + "." + o.parent + "." + o.name
	default:
		return "con:" + o.schema + "." + o.parent + "." + o.name
	}
}

// sortKey orders objects the dependency graph leaves unordered, including the
// members of a cycle: by kind, then name.
func (o *metadataObject) sortKey() string {
	key := metadataKindLabels[o.kind] + "\x00" + o.schema + "\x00" + o.name
	switch o.kind {
	case metadataFunction:
		key += "\x00" + metadataFunctionIdentity(o.function)
	case metadataIndex, metadataConstraint:
		key += "\x00" + o.parent
	default:
	}
	return key
}

// metadataKindLabels orders kinds. Functions come before relations because
// CreateFunctionStmt does not analyze a body, while a default, an index
// expression, or a view body is analyzed as it installs and needs every
// function it calls. Sequences come before the defaults that draw from them,
// and constraints after the indexes an FK can reference.
var metadataKindLabels = map[metadataObjectKind]string{
	metadataSchema:     "0schema",
	metadataEnum:       "1enum",
	metadataComposite:  "1composite",
	metadataSequence:   "2seq",
	metadataFunction:   "3function",
	metadataTable:      "4table",
	metadataView:       "5view",
	metadataMatView:    "6matview",
	metadataIndex:      "7index",
	metadataConstraint: "8constraint",
}

func metadataTypeKey(schema, name string) string     { return "type:" + schema + "." + name }
func metadataSequenceKey(schema, name string) string { return "seq:" + schema + "." + name }
func metadataRelationKey(schema, name string) string { return "rel:" + schema + "." + name }

// metadataFunctionIdentity distinguishes overloads. A snapshot synced before
// sync recorded signatures has none, so the definition stands in.
func metadataFunctionIdentity(fn *metadata.FunctionMetadata) string {
	if fn.GetSignature() != "" {
		return fn.GetSignature()
	}
	return fn.GetDefinition()
}

func collectMetadataObjects(meta *metadata.DatabaseSchemaMetadata, full bool) []*metadataObject {
	var out []*metadataObject
	for _, s := range meta.GetSchemas() {
		if s.GetName() == "" {
			continue
		}
		schema := s.GetName()
		out = append(out, &metadataObject{kind: metadataSchema, schema: schema, name: schema})
		owned, identity := sequenceOwners(s)
		for _, e := range s.GetEnumTypes() {
			out = append(out, &metadataObject{kind: metadataEnum, schema: schema, name: e.GetName(), enum: e})
		}
		for _, ct := range s.GetCompositeTypes() {
			out = append(out, &metadataObject{kind: metadataComposite, schema: schema, name: ct.GetName(), composite: ct})
		}
		for _, t := range s.GetTables() {
			o := &metadataObject{kind: metadataTable, schema: schema, name: t.GetName(), table: t}
			if full {
				o.ownedSequences = owned[t.GetName()]
			}
			out = append(out, o)
			out = appendPartitionObjects(out, schema, t, full)
		}
		// A foreign table loads as a plain table: query analysis needs only its columns.
		for _, et := range s.GetExternalTables() {
			out = append(out, &metadataObject{kind: metadataTable, schema: schema, name: et.GetName(), table: &metadata.TableMetadata{
				Name:    et.GetName(),
				Columns: et.GetColumns(),
			}})
		}
		for _, v := range s.GetViews() {
			body, err := parseSelect(v.GetDefinition())
			out = append(out, &metadataObject{kind: metadataView, schema: schema, name: v.GetName(), view: v, body: body, bodyErr: err})
		}
		for _, mv := range s.GetMaterializedViews() {
			body, err := parseSelect(mv.GetDefinition())
			out = append(out, &metadataObject{kind: metadataMatView, schema: schema, name: mv.GetName(), matView: mv, body: body, bodyErr: err})
		}
		for _, fn := range s.GetFunctions() {
			out = append(out, functionObject(schema, fn, false))
		}
		if !full {
			continue
		}
		for _, seq := range s.GetSequences() {
			// An identity column creates its own sequence.
			if identity[seq.GetOwnerTable()+"\x00"+seq.GetOwnerColumn()] {
				continue
			}
			out = append(out, &metadataObject{kind: metadataSequence, schema: schema, name: seq.GetName(), sequence: seq})
		}
		for _, p := range s.GetProcedures() {
			out = append(out, functionObject(schema, &metadata.FunctionMetadata{
				Name:       p.GetName(),
				Definition: p.GetDefinition(),
				Signature:  p.GetSignature(),
			}, true))
		}
		for _, t := range s.GetTables() {
			out = appendConstraintObjects(out, schema, t.GetName(), t.GetIndexes(), t.GetCheckConstraints(), t.GetExcludeConstraints())
			for _, fk := range t.GetForeignKeys() {
				out = append(out, &metadataObject{kind: metadataConstraint, schema: schema, name: fk.GetName(), parent: t.GetName(), fk: fk})
			}
		}
		for _, mv := range s.GetMaterializedViews() {
			out = appendIndexObjects(out, schema, mv.GetName(), mv.GetIndexes())
		}
	}
	return out
}

// appendPartitionObjects adds t's partitions, and theirs, as plain tables with
// t's column names and types, and with full their own indexes and constraints.
// t keeps the columns' defaults and sequences.
func appendPartitionObjects(out []*metadataObject, schema string, t *metadata.TableMetadata, full bool) []*metadataObject {
	if len(t.GetPartitions()) == 0 {
		return out
	}
	columns := make([]*metadata.ColumnMetadata, 0, len(t.GetColumns()))
	for _, col := range t.GetColumns() {
		columns = append(columns, &metadata.ColumnMetadata{Name: col.GetName(), Type: col.GetType(), Nullable: col.GetNullable(), Collation: col.GetCollation()})
	}
	var add func(partitions []*metadata.TablePartitionMetadata)
	add = func(partitions []*metadata.TablePartitionMetadata) {
		for _, p := range partitions {
			out = append(out, &metadataObject{kind: metadataTable, schema: schema, name: p.GetName(), table: &metadata.TableMetadata{Name: p.GetName(), Columns: columns}})
			if full {
				out = appendConstraintObjects(out, schema, p.GetName(), p.GetIndexes(), p.GetCheckConstraints(), p.GetExcludeConstraints())
			}
			add(p.GetSubpartitions())
		}
	}
	add(t.GetPartitions())
	return out
}

// appendConstraintObjects adds a table's indexes and its CHECK and EXCLUDE
// constraints.
func appendConstraintObjects(out []*metadataObject, schema, table string, indexes []*metadata.IndexMetadata, checks []*metadata.CheckConstraintMetadata, excludes []*metadata.ExcludeConstraintMetadata) []*metadataObject {
	out = appendIndexObjects(out, schema, table, indexes)
	for _, check := range checks {
		out = append(out, &metadataObject{kind: metadataConstraint, schema: schema, name: check.GetName(), parent: table, check: check})
	}
	for _, exclude := range excludes {
		out = append(out, &metadataObject{kind: metadataConstraint, schema: schema, name: exclude.GetName(), parent: table, exclude: exclude})
	}
	return out
}

// functionObject collects a function or procedure. PostgreSQL sync lists
// procedures among the functions, so a definition that parses says which.
func functionObject(schema string, fn *metadata.FunctionMetadata, procedure bool) *metadataObject {
	stmt, _ := parseStatement[*nodes.CreateFunctionStmt](fn.GetDefinition())
	if stmt != nil && createFunctionStmtIsProcedure(stmt) {
		procedure = true
	}
	return &metadataObject{kind: metadataFunction, schema: schema, name: fn.GetName(), function: fn, procedure: procedure, createStmt: stmt}
}

func appendIndexObjects(out []*metadataObject, schema, relation string, indexes []*metadata.IndexMetadata) []*metadataObject {
	for _, idx := range indexes {
		// The index behind an EXCLUDE constraint installs with the constraint.
		if idx.GetIsConstraint() && !isKeyIndex(idx) {
			continue
		}
		out = append(out, &metadataObject{kind: metadataIndex, schema: schema, name: idx.GetName(), parent: relation, index: idx})
	}
	return out
}

// sequenceOwners maps each table in s to the sequences its columns own, and
// returns the identity columns, keyed "table\x00column".
func sequenceOwners(s *metadata.SchemaMetadata) (owned map[string]map[string]*metadata.SequenceMetadata, identity map[string]bool) {
	identity = make(map[string]bool)
	for _, t := range s.GetTables() {
		for _, col := range t.GetColumns() {
			if col.GetIsIdentity() {
				identity[t.GetName()+"\x00"+col.GetName()] = true
			}
		}
	}
	owned = make(map[string]map[string]*metadata.SequenceMetadata)
	for _, seq := range s.GetSequences() {
		if seq.GetOwnerTable() == "" || seq.GetOwnerColumn() == "" {
			continue
		}
		if owned[seq.GetOwnerTable()] == nil {
			owned[seq.GetOwnerTable()] = make(map[string]*metadata.SequenceMetadata)
		}
		owned[seq.GetOwnerTable()][seq.GetOwnerColumn()] = seq
	}
	return owned, identity
}

func isViewObject(o *metadataObject) bool {
	return o.kind == metadataView || o.kind == metadataMatView
}

// isKeyIndex reports whether an index backs a PRIMARY KEY or UNIQUE
// constraint, which installs as the constraint rather than a bare index.
func isKeyIndex(idx *metadata.IndexMetadata) bool {
	return idx.GetPrimary() || (idx.GetIsConstraint() && idx.GetUnique())
}

// orderMetadataObjects groups objects into the strongly connected components
// of their dependencies, each sorted by sortKey, and orders the groups so each
// one follows the groups it depends on.
func orderMetadataObjects(objects []*metadataObject) [][]*metadataObject {
	if len(objects) == 0 {
		return nil
	}
	byKey := make(map[string]*metadataObject, len(objects))
	for _, o := range objects {
		byKey[o.key()] = o
	}
	edges := metadataDependencies(objects, byKey)
	sccs := stronglyConnectedObjects(objects, byKey, edges)

	sccOf := make(map[string]int, len(objects))
	for i, scc := range sccs {
		for _, o := range scc {
			sccOf[o.key()] = i
		}
	}
	// Kahn's algorithm over the condensed graph, taking ready components in
	// order of their smallest member.
	dependents := make([][]int, len(sccs))
	inDegree := make([]int, len(sccs))
	seen := make(map[[2]int]bool)
	for from, deps := range edges {
		fromSCC := sccOf[from]
		for _, dep := range deps {
			depSCC := sccOf[dep]
			if depSCC == fromSCC || seen[[2]int{depSCC, fromSCC}] {
				continue
			}
			seen[[2]int{depSCC, fromSCC}] = true
			dependents[depSCC] = append(dependents[depSCC], fromSCC)
			inDegree[fromSCC]++
		}
	}
	ready := &componentQueue{minKeys: make([]string, len(sccs))}
	for i, scc := range sccs {
		slices.SortFunc(scc, func(a, b *metadataObject) int { return cmp.Compare(a.sortKey(), b.sortKey()) })
		ready.minKeys[i] = scc[0].sortKey()
		if inDegree[i] == 0 {
			ready.components = append(ready.components, i)
		}
	}
	heap.Init(ready)
	ordered := make([][]*metadataObject, 0, len(sccs))
	for ready.Len() > 0 {
		next := heap.Pop(ready).(int)
		ordered = append(ordered, sccs[next])
		for _, d := range dependents[next] {
			inDegree[d]--
			if inDegree[d] == 0 {
				heap.Push(ready, d)
			}
		}
	}
	return ordered
}

// componentQueue is a min-heap of strongly connected components keyed by
// their smallest member's sort key.
type componentQueue struct {
	components []int
	minKeys    []string
}

func (q *componentQueue) Len() int { return len(q.components) }
func (q *componentQueue) Less(i, j int) bool {
	return q.minKeys[q.components[i]] < q.minKeys[q.components[j]]
}
func (q *componentQueue) Swap(i, j int) {
	q.components[i], q.components[j] = q.components[j], q.components[i]
}
func (q *componentQueue) Push(x any) { q.components = append(q.components, x.(int)) }
func (q *componentQueue) Pop() any {
	last := q.components[len(q.components)-1]
	q.components = q.components[:len(q.components)-1]
	return last
}

// metadataDependencies maps each object's key to the keys of the objects it
// needs installed first. It records only dependencies the snapshot contains.
func metadataDependencies(objects []*metadataObject, byKey map[string]*metadataObject) map[string][]string {
	// functionKeys maps "schema.name" to the keys of the function's overloads.
	functionKeys := make(map[string][]string)
	for _, o := range objects {
		if o.kind == metadataFunction {
			functionKeys[o.schema+"."+o.name] = append(functionKeys[o.schema+"."+o.name], o.key())
		}
	}

	edges := make(map[string][]string, len(objects))
	for _, o := range objects {
		if deps := o.dependencies(byKey, functionKeys); len(deps) > 0 {
			edges[o.key()] = deps
		}
	}
	return edges
}

func (o *metadataObject) dependencies(byKey map[string]*metadataObject, functionKeys map[string][]string) []string {
	self := o.key()
	var deps []string
	add := func(key string) {
		if byKey[key] != nil && key != self && !slices.Contains(deps, key) {
			deps = append(deps, key)
		}
	}
	addRef := func(schema, name string) {
		add(metadataTypeKey(schema, name))
		// A column, attribute, or parameter may use a relation's row type.
		add(metadataRelationKey(schema, name))
	}
	addType := func(typeStr string) {
		if schema, name, ok := userTypeRef(typeStr); ok {
			addRef(schema, name)
		}
	}
	addTypeName := func(tn *nodes.TypeName) {
		if tn == nil || tn.Names == nil || len(tn.Names.Items) != 2 {
			return
		}
		schema, schemaOK := tn.Names.Items[0].(*nodes.String)
		name, nameOK := tn.Names.Items[1].(*nodes.String)
		if schemaOK && nameOK {
			addRef(schema.Str, name.Str)
		}
	}
	if o.kind != metadataSchema {
		add("schema:" + o.schema)
	}
	switch o.kind {
	case metadataComposite:
		for _, a := range o.composite.GetAttributes() {
			addType(a.GetType())
		}
	case metadataTable:
		for _, col := range o.table.GetColumns() {
			addType(col.GetType())
		}
	case metadataView, metadataMatView:
		// A view whose body does not parse installs as a stand-in, which needs nothing.
		if o.body == nil {
			break
		}
		// Sync qualifies every name in a view's body, which reads as "schema.name".
		funcRefs, relRefs, typeRefs := collectExprDeps(o.body)
		for _, ref := range relRefs {
			add("rel:" + ref)
		}
		for _, ref := range typeRefs {
			add("type:" + ref)
			add("rel:" + ref)
		}
		// A call does not say which overload it uses, so the view waits for
		// all of them; LoadMetadata retries the cycles that can close.
		for _, ref := range funcRefs {
			for _, key := range functionKeys[ref] {
				add(key)
			}
		}
	case metadataFunction:
		if o.createStmt == nil {
			argTypes, _ := signatureArgTypes(o.function.GetName(), o.function.GetSignature())
			for _, t := range argTypes {
				addType(t)
			}
			break
		}
		// Every parameter, including OUT and TABLE ones, and the result.
		if params := o.createStmt.Parameters; params != nil {
			for _, item := range params.Items {
				if p, ok := item.(*nodes.FunctionParameter); ok {
					addTypeName(p.ArgType)
				}
			}
		}
		addTypeName(o.createStmt.ReturnType)
	case metadataIndex:
		add(metadataRelationKey(o.schema, o.parent))
	case metadataConstraint:
		add(metadataRelationKey(o.schema, o.parent))
		if o.fk != nil {
			refSchema := o.fk.GetReferencedSchema()
			if refSchema == "" {
				refSchema = o.schema
			}
			add(metadataRelationKey(refSchema, o.fk.GetReferencedTable()))
		}
	default:
	}
	return deps
}

// stronglyConnectedObjects runs Tarjan's algorithm, visiting keys in sorted
// order so the result is deterministic.
func stronglyConnectedObjects(objects []*metadataObject, byKey map[string]*metadataObject, edges map[string][]string) [][]*metadataObject {
	type state struct {
		index, low int
		onStack    bool
	}
	states := make(map[string]*state, len(objects))
	keys := make([]string, 0, len(byKey))
	for k := range byKey {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	var (
		next   int
		stack  []string
		result [][]*metadataObject
	)
	var visit func(v string)
	visit = func(v string) {
		states[v] = &state{index: next, low: next, onStack: true}
		next++
		stack = append(stack, v)
		for _, w := range edges[v] {
			if ws, ok := states[w]; !ok {
				visit(w)
				states[v].low = min(states[v].low, states[w].low)
			} else if ws.onStack {
				states[v].low = min(states[v].low, ws.index)
			}
		}
		if states[v].low != states[v].index {
			return
		}
		var scc []*metadataObject
		for {
			w := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			states[w].onStack = false
			scc = append(scc, byKey[w])
			if w == v {
				break
			}
		}
		result = append(result, scc)
	}
	for _, k := range keys {
		if _, ok := states[k]; !ok {
			visit(k)
		}
	}
	return result
}

func (c *Catalog) installMetadataObject(o *metadataObject, full bool) error {
	switch o.kind {
	case metadataSchema:
		err := c.CreateSchemaCommand(&nodes.CreateSchemaStmt{Schemaname: o.schema})
		// New() already created public, pg_catalog, and the other built-in schemas.
		var catErr *Error
		if errors.As(err, &catErr) && catErr.Code == CodeDuplicateSchema {
			return nil
		}
		return err
	case metadataEnum:
		return c.DefineEnum(metadataEnumStmt(o.schema, o.enum))
	case metadataComposite:
		stmt, err := metadataCompositeTypeStmt(o.schema, o.composite)
		if err != nil {
			return err
		}
		return c.DefineCompositeType(stmt)
	case metadataSequence:
		return c.DefineSequence(metadataSequenceStmt(o.schema, o.sequence))
	case metadataTable:
		stmt, err := metadataTableStmt(o.schema, o.table, full, o.ownedSequences)
		if err != nil {
			return err
		}
		if err := c.DefineRelation(stmt, 'r'); err != nil {
			return err
		}
		c.ownSerialSequences(o)
		return nil
	case metadataView:
		if o.body == nil {
			return fmt.Errorf("view %q: %w", o.name, o.bodyErr)
		}
		return c.DefineView(&nodes.ViewStmt{
			View:  &nodes.RangeVar{Schemaname: o.schema, Relname: o.name, Relpersistence: 'p'},
			Query: o.body,
		})
	case metadataMatView:
		if o.body == nil {
			return fmt.Errorf("materialized view %q: %w", o.name, o.bodyErr)
		}
		return c.ExecCreateTableAs(matViewStmt(o.schema, o.name, o.body))
	case metadataFunction:
		if o.createStmt == nil {
			return errors.New("no definition that parses")
		}
		return c.CreateFunctionStmt(o.createStmt)
	case metadataIndex:
		if isKeyIndex(o.index) {
			return c.AddConstraint(o.schema, o.parent, metadataKeyConstraint(o.index))
		}
		stmt, err := metadataIndexStmt(o.schema, o.parent, o.index)
		if err != nil {
			return err
		}
		return c.DefineIndex(stmt)
	case metadataConstraint:
		if o.fk != nil {
			return c.AddConstraint(o.schema, o.parent, metadataForeignKey(o.schema, o.fk))
		}
		stmt, err := metadataConstraintStmt(o)
		if err != nil {
			return err
		}
		return c.AlterTableStmt(stmt)
	default:
		return fmt.Errorf("unknown metadata object kind %d", o.kind)
	}
}

// ownSerialSequences gives the sequences the table's non-identity columns own
// their owner, as SERIAL does, so dropping the table, or its stand-in, drops
// them. Errors are ignored: the table itself installed.
func (c *Catalog) ownSerialSequences(o *metadataObject) {
	for _, col := range o.table.GetColumns() {
		seq := o.ownedSequences[col.GetName()]
		if seq == nil || col.GetIsIdentity() {
			continue
		}
		stmt, err := parseStatement[*nodes.AlterSeqStmt](fmt.Sprintf("ALTER SEQUENCE %s.%s OWNED BY %s.%s.%s",
			quoteIdent(o.schema), quoteIdent(seq.GetName()), quoteIdent(o.schema), quoteIdent(o.name), quoteIdent(col.GetName())))
		if err == nil {
			_ = c.AlterSequenceStmt(stmt)
		}
	}
}

// metadataStandIn returns the installer of an object's stand-in, or nil for
// kinds without one. A function's stand-in keeps its argument types, so it
// cannot win calls meant for another overload.
func metadataStandIn(o *metadataObject) func(*Catalog) error {
	switch o.kind {
	case metadataFunction:
		return func(c *Catalog) error {
			stmt, err := signatureFunctionStmt(o.schema, o.function, o.procedure)
			if err != nil {
				return err
			}
			return c.CreateFunctionStmt(stmt)
		}
	case metadataEnum:
		return func(c *Catalog) error {
			return c.DefineEnum(&nodes.CreateEnumStmt{TypeName: metadataNameList(o.schema, o.name), Vals: &nodes.List{}})
		}
	case metadataComposite:
		return func(c *Catalog) error {
			return c.DefineCompositeType(standInCompositeTypeStmt(o.schema, o.name, o.composite))
		}
	case metadataTable:
		return func(c *Catalog) error {
			if err := c.DefineRelation(standInTableStmt(o.schema, o.name, tableColumnNames(o.table)), 'r'); err != nil {
				return err
			}
			c.ownSerialSequences(o)
			return nil
		}
	case metadataView:
		return func(c *Catalog) error {
			stmt, err := standInViewStmt(o.schema, o.name, standInColumnNames(o))
			if err != nil {
				return err
			}
			return c.DefineView(stmt)
		}
	case metadataMatView:
		return func(c *Catalog) error {
			stmt, err := standInMatViewStmt(o.schema, o.name, standInColumnNames(o))
			if err != nil {
				return err
			}
			return c.ExecCreateTableAs(stmt)
		}
	default:
		return nil
	}
}
