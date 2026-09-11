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
// stand-in with the same name and column names, all typed text, so the objects
// that depend on it still install as defined. Other kinds are left out.
//
// View bodies are analyzed as they install, so set the search path and session
// user before loading. The error is non-nil only when ctx is done.
func (c *Catalog) LoadMetadata(ctx context.Context, meta *metadata.DatabaseSchemaMetadata, opts LoadMetadataOptions) (*LoadMetadataReport, error) {
	report := &LoadMetadataReport{
		Degraded: make(map[string]error),
		Missing:  make(map[string]error),
	}
	for _, o := range orderMetadataObjects(collectMetadataObjects(meta, opts.Full)) {
		if err := ctx.Err(); err != nil {
			return report, err
		}
		err := c.installMetadataObject(o, opts.Full)
		if err == nil {
			continue
		}
		standIn := metadataStandIn(o)
		if standIn == nil {
			report.Missing[o.key()] = err
			continue
		}
		if standInErr := standIn(c); standInErr != nil {
			report.Missing[o.key()] = standInErr
			continue
		}
		report.Degraded[o.key()] = err
	}
	return report, nil
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

// metadataObject is one install step. Exactly one metadata pointer is set,
// matching kind; a constraint sets one of fk, check, or exclude.
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
// members of a cycle.
func (o *metadataObject) sortKey() string {
	key := o.schema + "\x00" + o.name + "\x00" + metadataKindLabels[o.kind]
	switch o.kind {
	case metadataFunction:
		key += "\x00" + metadataFunctionIdentity(o.function)
	case metadataIndex, metadataConstraint:
		key += "\x00" + o.parent
	default:
	}
	return key
}

var metadataKindLabels = map[metadataObjectKind]string{
	metadataSchema:     "0schema",
	metadataEnum:       "1enum",
	metadataComposite:  "1composite",
	metadataSequence:   "2seq",
	metadataTable:      "3table",
	metadataView:       "4view",
	metadataMatView:    "5matview",
	metadataFunction:   "6function",
	metadataIndex:      "7index",
	metadataConstraint: "8constraint",
}

func metadataTypeKey(schema, name string) string     { return "type:" + schema + "." + name }
func metadataSequenceKey(schema, name string) string { return "seq:" + schema + "." + name }
func metadataRelationKey(schema, name string) string { return "rel:" + schema + "." + name }

// metadataFunctionIdentity distinguishes overloads; sync leaves Signature empty for
// some functions, so the definition stands in.
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
		for _, e := range s.GetEnumTypes() {
			out = append(out, &metadataObject{kind: metadataEnum, schema: schema, name: e.GetName(), enum: e})
		}
		for _, ct := range s.GetCompositeTypes() {
			out = append(out, &metadataObject{kind: metadataComposite, schema: schema, name: ct.GetName(), composite: ct})
		}
		for _, t := range s.GetTables() {
			out = append(out, &metadataObject{kind: metadataTable, schema: schema, name: t.GetName(), table: t})
		}
		for _, v := range s.GetViews() {
			out = append(out, &metadataObject{kind: metadataView, schema: schema, name: v.GetName(), view: v})
		}
		for _, mv := range s.GetMaterializedViews() {
			out = append(out, &metadataObject{kind: metadataMatView, schema: schema, name: mv.GetName(), matView: mv})
		}
		for _, fn := range s.GetFunctions() {
			out = append(out, &metadataObject{kind: metadataFunction, schema: schema, name: fn.GetName(), function: fn})
		}
		if !full {
			continue
		}
		for _, seq := range s.GetSequences() {
			// DefineRelation creates the sequence behind an identity column.
			if ownedByIdentityColumn(seq, s.GetTables()) {
				continue
			}
			out = append(out, &metadataObject{kind: metadataSequence, schema: schema, name: seq.GetName(), sequence: seq})
		}
		for _, p := range s.GetProcedures() {
			out = append(out, &metadataObject{kind: metadataFunction, schema: schema, name: p.GetName(), function: &metadata.FunctionMetadata{
				Name:       p.GetName(),
				Definition: p.GetDefinition(),
				Signature:  p.GetSignature(),
			}})
		}
		for _, t := range s.GetTables() {
			out = appendIndexObjects(out, schema, t.GetName(), t.GetIndexes())
			for _, fk := range t.GetForeignKeys() {
				out = append(out, &metadataObject{kind: metadataConstraint, schema: schema, name: fk.GetName(), parent: t.GetName(), fk: fk})
			}
			for _, check := range t.GetCheckConstraints() {
				out = append(out, &metadataObject{kind: metadataConstraint, schema: schema, name: check.GetName(), parent: t.GetName(), check: check})
			}
			for _, exclude := range t.GetExcludeConstraints() {
				out = append(out, &metadataObject{kind: metadataConstraint, schema: schema, name: exclude.GetName(), parent: t.GetName(), exclude: exclude})
			}
		}
		for _, mv := range s.GetMaterializedViews() {
			out = appendIndexObjects(out, schema, mv.GetName(), mv.GetIndexes())
		}
	}
	return out
}

func appendIndexObjects(out []*metadataObject, schema, relation string, indexes []*metadata.IndexMetadata) []*metadataObject {
	for _, idx := range indexes {
		out = append(out, &metadataObject{kind: metadataIndex, schema: schema, name: idx.GetName(), parent: relation, index: idx})
	}
	return out
}

func ownedByIdentityColumn(seq *metadata.SequenceMetadata, tables []*metadata.TableMetadata) bool {
	if seq.GetOwnerTable() == "" || seq.GetOwnerColumn() == "" {
		return false
	}
	for _, t := range tables {
		if t.GetName() != seq.GetOwnerTable() {
			continue
		}
		for _, col := range t.GetColumns() {
			if col.GetName() == seq.GetOwnerColumn() && col.GetIsIdentity() {
				return true
			}
		}
	}
	return false
}

// isKeyIndex reports whether an index backs a PRIMARY KEY or UNIQUE
// constraint, which installs as the constraint rather than a bare index.
func isKeyIndex(idx *metadata.IndexMetadata) bool {
	return idx.GetPrimary() || (idx.GetIsConstraint() && idx.GetUnique())
}

// orderMetadataObjects returns objects in an order that installs each one's
// dependencies first. A dependency cycle installs its lexically first member
// first; that member fails and gets a stand-in, and the rest install against it.
func orderMetadataObjects(objects []*metadataObject) []*metadataObject {
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
	ordered := make([]*metadataObject, 0, len(objects))
	for ready.Len() > 0 {
		next := heap.Pop(ready).(int)
		ordered = append(ordered, sccs[next]...)
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
	// keyIndexes maps "schema.table" to the keys of the table's PRIMARY KEY
	// and UNIQUE constraint indexes, which a referencing FK needs.
	keyIndexes := make(map[string][]string)
	for _, o := range objects {
		if o.kind == metadataIndex && isKeyIndex(o.index) {
			keyIndexes[o.schema+"."+o.parent] = append(keyIndexes[o.schema+"."+o.parent], o.key())
		}
	}

	edges := make(map[string][]string, len(objects))
	for _, o := range objects {
		self := o.key()
		var deps []string
		add := func(key string) {
			if byKey[key] != nil && key != self && !slices.Contains(deps, key) {
				deps = append(deps, key)
			}
		}
		addType := func(typeStr string) {
			if schema, name, ok := userTypeRef(typeStr); ok {
				add(metadataTypeKey(schema, name))
				// A composite attribute or column may use a relation's row type.
				add(metadataRelationKey(schema, name))
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
				if schema, name, ok := nextvalSequence(col.GetDefault(), o.schema); ok {
					add(metadataSequenceKey(schema, name))
				}
			}
		case metadataView:
			for _, dep := range o.view.GetDependencyColumns() {
				add(metadataRelationKey(dep.GetSchema(), dep.GetTable()))
			}
		case metadataMatView:
			for _, dep := range o.matView.GetDependencyColumns() {
				add(metadataRelationKey(dep.GetSchema(), dep.GetTable()))
			}
		case metadataFunction:
			argTypes, _ := signatureArgTypes(o.function.GetSignature())
			for _, t := range argTypes {
				addType(t)
			}
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
				for _, k := range keyIndexes[refSchema+"."+o.fk.GetReferencedTable()] {
					add(k)
				}
			}
		default:
		}
		if len(deps) > 0 {
			edges[self] = deps
		}
	}
	return edges
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
		stmt, err := metadataTableStmt(o.schema, o.table, full)
		if err != nil {
			return err
		}
		return c.DefineRelation(stmt, 'r')
	case metadataView:
		stmt, err := metadataViewStmt(o.schema, o.view)
		if err != nil {
			return err
		}
		return c.DefineView(stmt)
	case metadataMatView:
		stmt, err := metadataMatViewStmt(o.schema, o.matView)
		if err != nil {
			return err
		}
		return c.ExecCreateTableAs(stmt)
	case metadataFunction:
		stmt, err := metadataFunctionStmt(o.schema, o.function)
		if err != nil {
			return err
		}
		return c.CreateFunctionStmt(stmt)
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
		return c.AddConstraint(o.schema, o.parent, metadataConstraintDef(o))
	default:
		return fmt.Errorf("unknown metadata object kind %d", o.kind)
	}
}

// metadataStandIn returns the installer of an object's text-typed stand-in, or
// nil for kinds without one. A function has none because a text-typed overload
// would win overload resolution calls meant for other overloads.
func metadataStandIn(o *metadataObject) func(*Catalog) error {
	switch o.kind {
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
			return c.DefineRelation(standInTableStmt(o.schema, o.name, tableColumnNames(o.table)), 'r')
		}
	case metadataView:
		return func(c *Catalog) error {
			stmt, err := standInViewStmt(o.schema, o.name, viewColumnNames(o.view))
			if err != nil {
				return err
			}
			return c.DefineView(stmt)
		}
	case metadataMatView:
		return func(c *Catalog) error {
			stmt, err := standInMatViewStmt(o.schema, o.name, dependencyColumnNames(o.matView.GetDependencyColumns()))
			if err != nil {
				return err
			}
			return c.ExecCreateTableAs(stmt)
		}
	default:
		return nil
	}
}
