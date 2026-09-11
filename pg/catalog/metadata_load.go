package catalog

import (
	"cmp"
	"container/heap"
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/bytebase/omni/metadata"
	nodes "github.com/bytebase/omni/pg/ast"
	pgparser "github.com/bytebase/omni/pg/parser"
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
		if err := c.installGroup(ctx, group, opts.Full, report); err != nil {
			return report, err
		}
	}
	return report, nil
}

// --- Objects ---

// metadataObjectKind lists the kinds in the order they install where
// dependencies leave it open. Functions come before relations because
// CreateFunctionStmt does not analyze a body, while a default, an index
// expression, or a view body is analyzed as it installs and needs every
// function it calls. Sequences come before the defaults that draw from them,
// and constraints after the indexes an FK can reference.
type metadataObjectKind int

const (
	metadataSchema metadataObjectKind = iota
	metadataEnum
	metadataComposite
	metadataSequence
	metadataFunction
	metadataTable
	metadataView
	metadataMatView
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

	// procedure marks a function that is a procedure.
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

// key names the object in the dependency graph and the report.
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

// compareMetadataObjects orders the objects dependencies leave unordered, by
// kind and then by name.
func compareMetadataObjects(a, b *metadataObject) int {
	return cmp.Or(
		cmp.Compare(a.kind, b.kind),
		cmp.Compare(a.schema, b.schema),
		cmp.Compare(a.name, b.name),
		cmp.Compare(a.parent, b.parent),
		cmp.Compare(metadataFunctionIdentity(a.function), metadataFunctionIdentity(b.function)),
	)
}

func isViewObject(o *metadataObject) bool {
	return o.kind == metadataView || o.kind == metadataMatView
}

// --- Collecting ---

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
			out = appendPartitionTables(out, schema, t)
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
			// Procedures listed here wait for full, as the snapshot's own do.
			if o := functionObject(schema, fn, false); full || !o.procedure {
				out = append(out, o)
			}
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
			forEachPartition(t, func(p *metadata.TablePartitionMetadata) {
				out = appendConstraintObjects(out, schema, p.GetName(), p.GetIndexes(), p.GetCheckConstraints(), p.GetExcludeConstraints())
			})
		}
		for _, mv := range s.GetMaterializedViews() {
			out = appendIndexObjects(out, schema, mv.GetName(), mv.GetIndexes())
		}
	}
	return out
}

// appendPartitionTables adds t's partitions as plain tables with t's column
// names and types; t keeps the columns' defaults and sequences.
func appendPartitionTables(out []*metadataObject, schema string, t *metadata.TableMetadata) []*metadataObject {
	if len(t.GetPartitions()) == 0 {
		return out
	}
	columns := make([]*metadata.ColumnMetadata, 0, len(t.GetColumns()))
	for _, col := range t.GetColumns() {
		columns = append(columns, &metadata.ColumnMetadata{Name: col.GetName(), Type: col.GetType(), Nullable: col.GetNullable(), Collation: col.GetCollation()})
	}
	forEachPartition(t, func(p *metadata.TablePartitionMetadata) {
		out = append(out, &metadataObject{kind: metadataTable, schema: schema, name: p.GetName(), table: &metadata.TableMetadata{Name: p.GetName(), Columns: columns}})
	})
	return out
}

// forEachPartition calls f for each of t's partitions and their subpartitions.
func forEachPartition(t *metadata.TableMetadata, f func(*metadata.TablePartitionMetadata)) {
	var walk func([]*metadata.TablePartitionMetadata)
	walk = func(partitions []*metadata.TablePartitionMetadata) {
		for _, p := range partitions {
			f(p)
			walk(p.GetSubpartitions())
		}
	}
	walk(t.GetPartitions())
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

// isKeyIndex reports whether an index backs a PRIMARY KEY or UNIQUE
// constraint, which installs as the constraint rather than a bare index.
func isKeyIndex(idx *metadata.IndexMetadata) bool {
	return idx.GetPrimary() || (idx.GetIsConstraint() && idx.GetUnique())
}

// functionObject collects a function or procedure. PostgreSQL sync lists
// procedures among the functions, so the definition says which.
func functionObject(schema string, fn *metadata.FunctionMetadata, procedure bool) *metadataObject {
	stmt, _ := parseStatement[*nodes.CreateFunctionStmt](fn.GetDefinition())
	procedure = procedure || isProcedureDefinition(fn.GetDefinition())
	return &metadataObject{kind: metadataFunction, schema: schema, name: fn.GetName(), function: fn, procedure: procedure, createStmt: stmt}
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

// --- Ordering ---

// orderMetadataObjects groups objects into the strongly connected components
// of their dependencies, each sorted by compareMetadataObjects, and orders the
// groups so each one follows the groups it depends on.
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
	// order of their first member.
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
	ready := &componentQueue{first: make([]*metadataObject, len(sccs))}
	for i, scc := range sccs {
		slices.SortFunc(scc, compareMetadataObjects)
		ready.first[i] = scc[0]
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

// componentQueue is a min-heap of strongly connected components, keyed by
// their first member.
type componentQueue struct {
	components []int
	first      []*metadataObject
}

func (q *componentQueue) Len() int { return len(q.components) }
func (q *componentQueue) Less(i, j int) bool {
	return compareMetadataObjects(q.first[q.components[i]], q.first[q.components[j]]) < 0
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
		// all of them; installGroup retries the cycles that can close.
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

// --- Installing ---

// installGroup installs one object, or the members of a dependency cycle,
// which can be false: a view waits for every overload of a function it calls.
// Each pass retries the members that failed, since they may be waiting on one
// another. When a pass installs nothing, one member is replaced by its
// stand-in: the first view, whose stand-in needs nothing, if the group has one.
func (c *Catalog) installGroup(ctx context.Context, group []*metadataObject, full bool, report *LoadMetadataReport) error {
	for len(group) > 0 {
		var failed []*metadataObject
		var errs []error
		for _, o := range group {
			if err := ctx.Err(); err != nil {
				return err
			}
			if err := c.installMetadataObject(o, full); err != nil {
				failed = append(failed, o)
				errs = append(errs, err)
			}
		}
		if len(failed) == len(group) {
			i := max(slices.IndexFunc(failed, isViewObject), 0)
			c.replaceWithStandIn(failed[i], errs[i], report)
			failed = slices.Delete(failed, i, i+1)
		}
		group = failed
	}
	return nil
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

// --- Statements ---

// The snapshot carries types, defaults, and bodies as the text PostgreSQL
// reports for them. These builders turn that text into the statement ASTs the
// Define* functions take.

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

func metadataSequenceStmt(schema string, seq *metadata.SequenceMetadata) *nodes.CreateSeqStmt {
	var opts []nodes.Node
	if tn, err := parseTypeName(seq.GetDataType()); err == nil {
		opts = append(opts, &nodes.DefElem{Defname: "as", Arg: tn})
	}
	opts = append(opts, sequenceOptions(seq)...)
	stmt := &nodes.CreateSeqStmt{
		Sequence: &nodes.RangeVar{Schemaname: schema, Relname: seq.GetName(), Relpersistence: 'p'},
	}
	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}
	return stmt
}

// sequenceOptions returns a sequence's settings other than its type.
func sequenceOptions(seq *metadata.SequenceMetadata) []nodes.Node {
	var opts []nodes.Node
	for _, o := range []struct{ name, value string }{
		{"increment", seq.GetIncrement()},
		{"minvalue", seq.GetMinValue()},
		{"maxvalue", seq.GetMaxValue()},
		{"start", seq.GetStart()},
		{"cache", seq.GetCacheSize()},
	} {
		if o.value == "" {
			continue
		}
		v, _ := strconv.ParseInt(o.value, 10, 64)
		opts = append(opts, &nodes.DefElem{Defname: o.name, Arg: &nodes.Integer{Ival: v}})
	}
	if seq.GetCycle() {
		opts = append(opts, &nodes.DefElem{Defname: "cycle", Arg: &nodes.Boolean{Boolval: true}})
	}
	return opts
}

// metadataTableStmt fails on any column whose type does not parse, so the
// table gets a stand-in rather than a partial definition. With full, a default
// or generation expression that does not parse is dropped instead.
func metadataTableStmt(schema string, t *metadata.TableMetadata, full bool, ownedSequences map[string]*metadata.SequenceMetadata) (*nodes.CreateStmt, error) {
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
		if col.GetCollation() != "" {
			// Sync reports a column's collation by its bare name.
			def.CollClause = &nodes.CollateClause{Collname: metadataNameList("", col.GetCollation())}
		}
		if full {
			setColumnValueProperties(schema, def, col, ownedSequences[col.GetName()])
		}
		items = append(items, def)
	}
	return &nodes.CreateStmt{
		Relation:  &nodes.RangeVar{Schemaname: schema, Relname: t.GetName(), Relpersistence: 'p'},
		TableElts: &nodes.List{Items: items},
	}, nil
}

// setColumnValueProperties sets a column's default, generation, and identity;
// seq is the sequence the column owns, when the snapshot has one.
func setColumnValueProperties(schema string, def *nodes.ColumnDef, col *metadata.ColumnMetadata, seq *metadata.SequenceMetadata) {
	if col.GetDefault() != "" && col.GetGeneration() == nil {
		if expr, err := parseScalar("SELECT " + col.GetDefault()); err == nil {
			def.RawDefault = expr
		}
	}
	if gen := col.GetGeneration(); gen.GetType() == metadata.GenerationMetadata_TYPE_STORED && gen.GetExpression() != "" {
		if expr, err := parseScalar("SELECT (" + gen.GetExpression() + ")"); err == nil {
			def.Generated = 's'
			def.Constraints = &nodes.List{Items: []nodes.Node{
				&nodes.Constraint{Contype: nodes.CONSTR_GENERATED, RawExpr: expr},
			}}
		}
	}
	if col.GetIsIdentity() {
		identity := &nodes.Constraint{Contype: nodes.CONSTR_IDENTITY, GeneratedWhen: 'd'}
		if col.GetIdentityGeneration() == metadata.ColumnMetadata_ALWAYS {
			identity.GeneratedWhen = 'a'
		}
		if seq != nil {
			opts := append([]nodes.Node{&nodes.DefElem{Defname: "sequence_name", Arg: metadataNameList(schema, seq.GetName())}}, sequenceOptions(seq)...)
			identity.Options = &nodes.List{Items: opts}
		}
		if def.Constraints == nil {
			def.Constraints = &nodes.List{}
		}
		def.Constraints.Items = append(def.Constraints.Items, identity)
	}
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

// metadataIndexStmt parses the index's definition, which keeps a partial
// index's predicate. Without one that parses, it rebuilds the index from its
// keys: a quoted identifier is a column, and so is a bare key unless it looks
// like an expression in pg_get_indexdef's shape, a function call or a
// parenthesized, cast, or spaced expression.
func metadataIndexStmt(schema, relation string, idx *metadata.IndexMetadata) (*nodes.IndexStmt, error) {
	if stmt, err := parseStatement[*nodes.IndexStmt](idx.GetDefinition()); err == nil {
		return stmt, nil
	}
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
		if name == expr && (strings.ContainsAny(name, "( ") || strings.Contains(name, "::")) {
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

func metadataForeignKey(schema string, fk *metadata.ForeignKeyMetadata) ConstraintDef {
	refSchema := fk.GetReferencedSchema()
	if refSchema == "" {
		refSchema = schema
	}
	return ConstraintDef{
		Name:        fk.GetName(),
		Type:        ConstraintFK,
		Columns:     fk.GetColumns(),
		RefSchema:   refSchema,
		RefTable:    fk.GetReferencedTable(),
		RefColumns:  fk.GetReferencedColumns(),
		FKUpdAction: fkActionCode(fk.GetOnUpdate()),
		FKDelAction: fkActionCode(fk.GetOnDelete()),
		FKMatchType: fkMatchCode(fk.GetMatchType()),
	}
}

// metadataConstraintStmt parses a CHECK or EXCLUDE constraint from its
// pg_get_constraintdef text; sync strips the CHECK keyword from checks.
func metadataConstraintStmt(o *metadataObject) (*nodes.AlterTableStmt, error) {
	def := o.exclude.GetExpression()
	if o.check != nil {
		def = "CHECK " + o.check.GetExpression()
	}
	return parseStatement[*nodes.AlterTableStmt](fmt.Sprintf("ALTER TABLE %s.%s ADD CONSTRAINT %s %s",
		quoteIdent(o.schema), quoteIdent(o.parent), quoteIdent(o.name), def))
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

// --- Stand-ins ---

// A stand-in keeps an object's name and column names and types every column
// text, so it depends on no user object and the objects that use it resolve.

// replaceWithStandIn installs the stand-in of o, whose definition failed with
// err, and records the outcome in report.
func (c *Catalog) replaceWithStandIn(o *metadataObject, err error, report *LoadMetadataReport) {
	switch standInErr := c.installStandIn(o); {
	case standInErr == nil:
		report.Degraded[o.key()] = err
	case errors.Is(standInErr, errNoStandIn):
		report.Missing[o.key()] = err
	default:
		report.Missing[o.key()] = standInErr
	}
}

var errNoStandIn = errors.New("no stand-in for this kind of object")

// installStandIn installs o's stand-in, or returns errNoStandIn for a kind
// without one. A function's stand-in keeps its argument types, so it cannot
// win calls meant for another overload.
func (c *Catalog) installStandIn(o *metadataObject) error {
	switch o.kind {
	case metadataEnum:
		return c.DefineEnum(&nodes.CreateEnumStmt{TypeName: metadataNameList(o.schema, o.name), Vals: &nodes.List{}})
	case metadataComposite:
		return c.DefineCompositeType(standInCompositeTypeStmt(o.schema, o.name, o.composite))
	case metadataFunction:
		stmt, err := standInFunctionStmt(o.schema, o.function, o.procedure)
		if err != nil {
			return err
		}
		return c.CreateFunctionStmt(stmt)
	case metadataTable:
		if err := c.DefineRelation(standInTableStmt(o.schema, o.name, tableColumnNames(o.table)), 'r'); err != nil {
			return err
		}
		c.ownSerialSequences(o)
		return nil
	case metadataView:
		stmt, err := standInViewStmt(o.schema, o.name, standInColumnNames(o))
		if err != nil {
			return err
		}
		return c.DefineView(stmt)
	case metadataMatView:
		stmt, err := standInMatViewStmt(o.schema, o.name, standInColumnNames(o))
		if err != nil {
			return err
		}
		return c.ExecCreateTableAs(stmt)
	default:
		return errNoStandIn
	}
}

// standInFunctionStmt builds a function's stand-in from its signature: the
// argument types and, unless it is a procedure, a text result, so calls by
// name still resolve.
func standInFunctionStmt(schema string, fn *metadata.FunctionMetadata, procedure bool) (*nodes.CreateFunctionStmt, error) {
	argTypes, err := signatureArgTypes(fn.GetName(), fn.GetSignature())
	if err != nil {
		return nil, fmt.Errorf("function %q: %w", fn.GetName(), err)
	}
	stmt := &nodes.CreateFunctionStmt{
		Funcname:   metadataNameList(schema, fn.GetName()),
		ReturnType: textTypeName(),
		Options:    sqlFunctionBody("SELECT NULL::text"),
	}
	if procedure {
		stmt.ReturnType = nil
		stmt.Options.Items = append(stmt.Options.Items, &nodes.DefElem{Defname: "isProcedure", Arg: &nodes.Boolean{Boolval: true}})
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

// standInColumnNames returns the columns a view's or materialized view's
// stand-in exposes: the view's own columns when sync reports them, else the
// ones its body outputs, else, when the body does not parse, the ones it reads.
func standInColumnNames(o *metadataObject) []string {
	var names []string
	for _, col := range o.view.GetColumns() {
		names = appendUniqueName(names, col.GetName())
	}
	if len(names) > 0 {
		return names
	}
	if o.body != nil {
		return targetColumnNames(o.body)
	}
	deps := o.view.GetDependencyColumns()
	if o.kind == metadataMatView {
		deps = o.matView.GetDependencyColumns()
	}
	for _, dep := range deps {
		names = appendUniqueName(names, dep.GetColumn())
	}
	return names
}

// targetColumnNames returns the names of the columns a query outputs; a set
// operation takes them from its first branch.
func targetColumnNames(sel *nodes.SelectStmt) []string {
	for sel.Larg != nil {
		sel = sel.Larg
	}
	if sel.TargetList == nil {
		return nil
	}
	var names []string
	for _, item := range sel.TargetList.Items {
		if rt, ok := item.(*nodes.ResTarget); ok {
			name := rt.Name
			if name == "" {
				name = figureColName(rt.Val)
			}
			names = appendUniqueName(names, name)
		}
	}
	return names
}

func appendUniqueName(names []string, name string) []string {
	if name == "" || slices.Contains(names, name) {
		return names
	}
	return append(names, name)
}

// --- Parsing snapshot text ---

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

func parseSelect(sql string) (*nodes.SelectStmt, error) {
	return parseStatement[*nodes.SelectStmt](sql)
}

// parseStatement parses sql as exactly one statement of type T.
func parseStatement[T nodes.Node](sql string) (T, error) {
	var zero T
	list, err := pgparser.Parse(sql)
	if err != nil {
		return zero, err
	}
	if list == nil || len(list.Items) != 1 {
		n := 0
		if list != nil {
			n = len(list.Items)
		}
		return zero, fmt.Errorf("expected one statement, got %d", n)
	}
	stmt, ok := unwrapRawStmt(list.Items[0]).(T)
	if !ok {
		return zero, fmt.Errorf("expected %T, got %T", zero, unwrapRawStmt(list.Items[0]))
	}
	return stmt, nil
}

// userTypeRef returns the schema and name a type string refers to when it is
// schema-qualified, as sync reports user-defined types: "public.status",
// "public.status[]", or `"My Schema"."My Type"`.
func userTypeRef(typ string) (schema, name string, ok bool) {
	s := strings.TrimSpace(typ)
	for strings.HasSuffix(s, "[]") {
		s = strings.TrimSpace(strings.TrimSuffix(s, "[]"))
	}
	// Drop type modifiers, which start at the first parenthesis outside quotes.
	quoted := false
	for i := 0; i < len(s); i++ {
		if s[i] == '"' {
			quoted = !quoted
		} else if s[i] == '(' && !quoted {
			s = strings.TrimSpace(s[:i])
			break
		}
	}
	return splitQualifiedIdent(s)
}

// isProcedureDefinition reports whether a definition creates a procedure, even
// one that does not parse; pg_get_functiondef starts it CREATE OR REPLACE
// PROCEDURE.
func isProcedureDefinition(def string) bool {
	head := strings.Join(strings.Fields(strings.ToUpper(def[:min(len(def), 64)])), " ")
	return strings.HasPrefix(head, "CREATE PROCEDURE ") || strings.HasPrefix(head, "CREATE OR REPLACE PROCEDURE ")
}

// signatureArgTypes splits the argument types out of a signature, which sync
// writes as the function's name followed by its arguments in parentheses, as
// in "fn(integer, numeric(10,2))". The name can itself contain parentheses.
func signatureArgTypes(name, signature string) ([]string, error) {
	signature = strings.TrimPrefix(signature, name)
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
	depth, start, quoted := 0, 0, false
	for i := 0; i < len(inner); i++ {
		if inner[i] == '"' {
			quoted = !quoted
		}
		if quoted {
			continue
		}
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
