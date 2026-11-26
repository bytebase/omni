package review

import (
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql"
	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/review/walk"
)

// The walk's spaces (docs/specs/2026-10-26-walk-through-mssql.md §2). A
// schema is a key of the root; an object, relation, or type a key of its
// schema's path; a column, index, constraint, or the primary key a key of
// its table's path, Path(schema class, table class). A temporary table is
// a key of the root in its own space, and its children are under
// Path(tempSchema, its class).
const (
	spaceSchema     = "schema"
	spaceObject     = "object"
	spaceRelation   = "relation"
	spaceType       = "type"
	spaceColumn     = "column"
	spaceIndex      = "index"
	spaceConstraint = "constraint"
	spacePrimary    = "primary"
	spaceTemp       = "temporary"
)

// tempSchema stands for tempdb in a temporary table's path. No class
// begins with it.
const tempSchema = "#"

// primaryName is the one key of a table's primary space: Present when the
// table has a primary key.
const primaryName = "key"

// Kinds. The object space holds every kind; the relation space tables,
// views, and sequences; the index space indexes, the PK's and UQ's of
// their kinds; the constraint space the constraint kinds.
const (
	kindTable     = "table"
	kindView      = "view"
	kindSequence  = "sequence"
	kindProcedure = "procedure"
	kindFunction  = "function"
	kindTrigger   = "trigger"
	kindSynonym   = "synonym"
	kindDefault   = "default"
	kindRule      = "rule"
	kindIndex     = "index"
	kindXMLIndex  = "xml"
	kindPK        = "pk"
	kindUQ        = "uq"
	kindFK        = "fk"
	kindCK        = "ck"
	kindDF        = "df"
)

// excludedSchemas are the schemas Bytebase's sync leaves out, which exist
// all the same; guest and db_owner can hold tables (verified).
//
// bytebase: backend/plugin/db/mssql/sync.go — getSchemas
var excludedSchemas = []string{"db_owner", "db_accessadmin", "db_backupoperator", "db_datareader", "db_datawriter", "db_ddladmin", "db_denydatareader", "db_denydatawriter", "db_securityadmin", "guest", "INFORMATION_SCHEMA", "sys"}

// snapName is one synced name of a class.
type snapName struct {
	spelling string
	kind     string
}

// snapSpace is one space under one parent: its names by class, and
// whether it holds a marked name.
type snapSpace struct {
	classes map[string][]snapName
	marked  bool
}

func (s *snapSpace) add(spelling, kind string) {
	n, ok := keyed(spelling)
	if !ok {
		s.marked = true
		return
	}
	if !n.plain {
		s.marked = true
	}
	if s.classes == nil {
		s.classes = make(map[string][]snapName)
	}
	s.classes[n.class] = append(s.classes[n.class], snapName{spelling, kind})
}

// lookup answers a class or exact key of the space. Absent is answered
// only where the sync lists the space completely (trusted) and the space
// holds no marked name.
func (s *snapSpace) lookup(key string, trusted bool) (walk.Fact, string) {
	if s == nil {
		if trusted {
			return walk.Absent, ""
		}
		return walk.Unknown, ""
	}
	if key == markedFlag {
		switch {
		case s.marked:
			return walk.Present, ""
		case trusted:
			return walk.Absent, ""
		}
		return walk.Unknown, ""
	}
	class, spelling, exact := strings.Cut(key, exactSep)
	names := s.classes[class]
	if exact {
		for _, n := range names {
			if n.spelling == spelling {
				return walk.Present, n.kind
			}
		}
		return walk.Unknown, ""
	}
	switch {
	case len(names) == 0:
		if trusted && !s.marked && !strings.HasPrefix(class, markedPrefix) {
			return walk.Absent, ""
		}
		return walk.Unknown, ""
	case len(names) == 1 && trusted:
		return walk.Present, classKind(names[0].kind, names[0].spelling)
	}
	return walk.Present, ""
}

// snapTable is a synced table's children.
type snapTable struct {
	columns, indexes, constraints snapSpace
	// names are the exact spellings of the columns, in order.
	names []string
	// colTypes maps exact synced column spelling to its lowercased base type
	// name (empty for computed or untyped columns).
	colTypes map[string]string
	// defaults maps the exact synced spelling of a column to the name of
	// its default constraint.
	defaults map[string]string
	primary  bool
	// edge is set for a graph edge table, whose edge constraints the sync
	// skips.
	edge bool
	// ledger is set for a table that may be a ledger table.
	ledger bool
	// clustered is the name of the table's clustered index, rowstore or
	// columnstore, as the sync's index type names it.
	clustered string
	// noPeriod is set when no column can be a SYSTEM_TIME period column:
	// each is computed or of a built-in type other than datetime2, which
	// a period column must be. Such a table is not system-versioned.
	noPeriod bool
	// noRowguid is set when no column can be the ROWGUIDCOL column, which
	// is of type uniqueidentifier: each is computed or of another
	// built-in type.
	noRowguid bool
	// noFilestream is set when no column can be a FILESTREAM column,
	// which is of type varbinary(max): each is computed or of another
	// built-in type.
	noFilestream bool
	// fileTable is false when the table lacks path_locator, a column of
	// FileTable's fixed schema.
	fileTable bool
	// identity holds the exact spellings of the identity columns.
	identity map[string]bool
	// nonclustered maps an index class to whether its one synced index
	// is a rowstore nonclustered index; a class two indexes share maps to
	// false.
	nonclustered map[string]bool
	// fkColumns maps the exact spelling of a column of one of the table's
	// foreign keys to that key's name.
	fkColumns map[string]string
	// referencedBy names the other synced tables whose foreign keys
	// reference the table, by path.
	referencedBy []fkSource
}

// fkSource is a table whose foreign key references another table.
type fkSource struct {
	path, written string
}

// noColumnOf reports whether no column of t can be of type typ: each is
// computed or of another built-in type.
func noColumnOf(t *metadata.TableMetadata, typ string) bool {
	for _, c := range t.GetColumns() {
		ct := strings.ToLower(strings.TrimSpace(c.GetType()))
		if strings.HasPrefix(ct, "as ") {
			continue
		}
		base, _, _ := strings.Cut(ct, "(")
		base = strings.TrimSpace(base)
		if base == "" || base == typ || !isBuiltinSystemType(base) {
			return false
		}
	}
	return true
}

// A ledger table has two or more NOT NULL ledger columns, of type bigint
// or an alias of it, which the sync does not mark (verified): any NOT NULL
// column of a type the sync names otherwise than these may be one.
var notLedgerTypes = map[string]bool{
	"bit": true, "tinyint": true, "smallint": true, "int": true, "decimal": true, "numeric": true, "money": true, "smallmoney": true,
	"float": true, "real": true, "date": true, "time": true, "datetime": true, "datetime2": true, "datetimeoffset": true, "smalldatetime": true,
	"char": true, "varchar": true, "nchar": true, "nvarchar": true, "text": true, "ntext": true, "binary": true, "varbinary": true, "image": true,
	"uniqueidentifier": true, "xml": true, "sql_variant": true, "hierarchyid": true, "geometry": true, "geography": true,
	"timestamp": true, "rowversion": true, "sysname": true, "json": true, "vector": true,
}

func ledgerTable(t *metadata.TableMetadata) bool {
	n := 0
	for _, c := range t.GetColumns() {
		typ := strings.ToLower(c.GetType())
		if c.GetNullable() || typ == "" || strings.HasPrefix(typ, "as ") {
			continue
		}
		if base, _, _ := strings.Cut(typ, "("); !notLedgerTypes[strings.TrimSpace(base)] {
			n++
		}
	}
	return n >= 2
}

// mssqlSnapshot is a target's synced schema as WalkThrough trusts it
// (spec §3): Present for every name the sync lists; Absent for schemas,
// for the relations of a synced schema (tables, views, sequences, and
// synonyms), and for the columns, indexes, constraints, and primary key
// of a synced table, unless the space holds a marked name; never for the
// object space, which holds CLR modules, disabled triggers, and more that
// the sync skips, or for types.
// The session starts with no temporary table.
//
// bytebase: backend/plugin/db/mssql/sync.go — SyncDBSchema
type mssqlSnapshot struct {
	db       string
	schemas  snapSpace
	excluded map[string]bool
	// objects and relations by schema class; tables by table path. A
	// schema class or table path two synced names share is ambiguous,
	// and every key under it Unknown. The names an ambiguous table's
	// children hold in its schema are kept.
	objects   map[string]*snapSpace
	relations map[string]*snapSpace
	tables    map[string]*snapTable
	ambiguous map[string]bool
	// triggers is set when the sync lists a trigger, on any table or
	// view, of an ambiguous parent too.
	triggers bool
	// ledger is set when the sync lists a table that may be a ledger
	// table, whose drop renames its history table and ledger view, in any
	// schema (verified).
	ledger bool
	// cdc is set when the sync lists cdc.change_tables: change data
	// capture is on, and a dropped table's capture instance goes too.
	cdc bool
	// unbound holds the paths of views whose synced definition is not
	// schema-bound: such a view has no index.
	unbound map[string]bool
	// triggerOn maps a DML trigger's path, Path(schema class, trigger
	// class), to its table or view; "" for a path two triggers share.
	triggerOn map[string]fkSource
	// viewCols maps a view's path to the names of its columns, as its
	// synced definition gives them; nil when it does not.
	viewCols map[string][]string
	// functions maps a function's path to what its synced definition
	// shows; nil for a path two functions share or a definition the walk
	// cannot read.
	functions map[string]*snapFunc
}

// Function kinds, as a synced definition shows them.
const (
	funcScalar = "scalar"
	funcInline = "inline table-valued"
	funcMulti  = "multi-statement table-valued"
)

// snapFunc is what a synced function's definition shows.
type snapFunc struct {
	kind   string
	params int
	bound  bool
	// cols are the names of a table-valued function's columns; nil when
	// the definition does not give them.
	cols []string
}

// syncedFunc reads a synced function's definition. It reads a definition
// the server accepted, so one that does not parse is only left unread.
func syncedFunc(definition string) *snapFunc {
	stmts, err := mssql.Parse(definition)
	if err != nil || len(stmts) != 1 {
		return nil
	}
	v, ok := stmts[0].AST.(*ast.CreateFunctionStmt)
	if !ok || v.ExternalName != "" {
		return nil
	}
	f := &snapFunc{kind: funcKind(v), params: len(listNodes(v.Params)), bound: hasOption(v.Options, "SCHEMABINDING")}
	switch rt := v.ReturnsTable; f.kind {
	case funcMulti:
		for _, item := range listNodes(rt.Columns) {
			cd, ok := item.(*ast.ColumnDef)
			if !ok {
				f.cols = nil
				break
			}
			f.cols = append(f.cols, cd.Name)
		}
	case funcInline:
		if ret, ok := v.Body.(*ast.ReturnStmt); ok {
			if sub, ok := ret.Value.(*ast.SubqueryExpr); ok {
				if sel, ok := sub.Query.(*ast.SelectStmt); ok {
					f.cols, _ = outputNames(sel)
				}
			}
		}
	}
	return f
}

// funcKind returns the kind of function a definition makes.
func funcKind(v *ast.CreateFunctionStmt) string {
	switch rt := v.ReturnsTable; {
	case rt == nil:
		return funcScalar
	case rt.Variable != "":
		return funcMulti
	}
	return funcInline
}

// syncedViewCols returns the names of a synced view's columns: those the
// sync read from sys.columns, which sp_rename keeps current. Metadata
// without them, as the in-memory catalog's, gives the definition's column
// list or the names its query gives; nil when the walk cannot tell.
func syncedViewCols(v *metadata.ViewMetadata) []string {
	if cols := v.GetColumns(); len(cols) > 0 {
		names := make([]string, 0, len(cols))
		for _, c := range cols {
			names = append(names, c.GetName())
		}
		return names
	}
	stmts, err := mssql.Parse(v.GetDefinition())
	if err != nil || len(stmts) != 1 {
		return nil
	}
	cv, ok := stmts[0].AST.(*ast.CreateViewStmt)
	if !ok {
		return nil
	}
	if cv.Columns != nil {
		names, _ := listNames(cv.Columns)
		return names
	}
	sel, ok := cv.Query.(*ast.SelectStmt)
	if !ok {
		return nil
	}
	names, _ := outputNames(sel)
	return names
}

// viewUnbound reports whether a synced view's definition is a CREATE VIEW
// without SCHEMABINDING. It reads a definition the server accepted, so a
// definition that does not parse is only left unjudged.
func viewUnbound(definition string) bool {
	stmts, err := mssql.Parse(definition)
	if err != nil || len(stmts) != 1 {
		return false
	}
	v, ok := stmts[0].AST.(*ast.CreateViewStmt)
	return ok && !v.SchemaBinding && !hasOption(v.Options, "SCHEMABINDING")
}

func newMSSQLSnapshot(db *metadata.DatabaseSchemaMetadata) *mssqlSnapshot {
	s := &mssqlSnapshot{
		db:        db.GetName(),
		excluded:  make(map[string]bool),
		objects:   make(map[string]*snapSpace),
		relations: make(map[string]*snapSpace),
		tables:    make(map[string]*snapTable),
		ambiguous: make(map[string]bool),
		unbound:   make(map[string]bool),
		triggerOn: make(map[string]fkSource),
		viewCols:  make(map[string][]string),
		functions: make(map[string]*snapFunc),
	}
	for _, e := range excludedSchemas {
		if n, ok := keyed(e); ok {
			s.excluded[n.class] = true
		}
	}
	// The facts of the whole database count every table and view the sync
	// lists, before any is keyed.
	for _, schema := range db.GetSchemas() {
		for _, t := range schema.GetTables() {
			s.ledger = s.ledger || ledgerTable(t)
			s.triggers = s.triggers || len(t.GetTriggers()) > 0
			s.cdc = s.cdc || schema.GetName() == "cdc" && t.GetName() == "change_tables"
		}
		for _, v := range schema.GetViews() {
			s.triggers = s.triggers || len(v.GetTriggers()) > 0
		}
	}
	type fkRef struct {
		from          fkSource
		schema, table string
	}
	var refs []fkRef
	for _, schema := range db.GetSchemas() {
		s.schemas.add(schema.GetName(), "")
		sn, ok := keyed(schema.GetName())
		if !ok {
			continue
		}
		sc := sn.class
		if _, seen := s.objects[sc]; seen {
			s.ambiguous[sc] = true
			continue
		}
		objs, rels := &snapSpace{}, &snapSpace{}
		views := make(map[string]bool)
		s.objects[sc], s.relations[sc] = objs, rels
		for _, t := range schema.GetTables() {
			objs.add(t.GetName(), kindTable)
			rels.add(t.GetName(), kindTable)
			// A table on an ambiguous path keeps no children, but the
			// names they hold in the schema count all the same.
			st := &snapTable{ledger: ledgerTable(t), noPeriod: noColumnOf(t, "datetime2"), noRowguid: noColumnOf(t, "uniqueidentifier"), noFilestream: noColumnOf(t, "varbinary"), identity: make(map[string]bool), nonclustered: make(map[string]bool), fkColumns: make(map[string]string)}
			if tn, ok := keyed(t.GetName()); ok {
				path := walk.Path(sc, tn.class)
				if _, seen := s.tables[path]; seen {
					s.ambiguous[path] = true
				} else {
					s.tables[path] = st
				}
			}
			for _, c := range t.GetColumns() {
				st.columns.add(c.GetName(), "")
				st.names = append(st.names, c.GetName())
				if typ := strings.ToLower(strings.TrimSpace(c.GetType())); typ != "" && !strings.HasPrefix(typ, "as ") {
					base, _, _ := strings.Cut(typ, "(")
					if base = strings.TrimSpace(base); base != "" {
						if st.colTypes == nil {
							st.colTypes = make(map[string]string)
						}
						st.colTypes[c.GetName()] = base
					}
				}
				if strings.HasPrefix(c.GetName(), "$edge_id") {
					st.edge = true
				}
				if strings.EqualFold(c.GetName(), "path_locator") {
					st.fileTable = true
				}
				if c.GetIsIdentity() {
					st.identity[c.GetName()] = true
				}
				if df := c.GetDefaultConstraintName(); df != "" {
					st.constraints.add(df, kindDF)
					objs.add(df, kindDF)
				}
				if df := c.GetDefaultConstraintName(); df != "" {
					if st.defaults == nil {
						st.defaults = make(map[string]string)
					}
					st.defaults[c.GetName()] = df
				}
			}
			for _, i := range t.GetIndexes() {
				kind := kindIndex
				switch {
				case strings.EqualFold(i.GetType(), "XML"):
					// Dropping a primary XML index drops its secondary
					// ones (verified).
					kind = kindXMLIndex
				case i.GetPrimary():
					kind = kindPK
					st.primary = true
				case i.GetUnique() && i.GetIsConstraint():
					kind = kindUQ
				}
				if strings.EqualFold(i.GetType(), "CLUSTERED") || strings.EqualFold(i.GetType(), "CLUSTERED COLUMNSTORE") {
					st.clustered = i.GetName()
				}
				if in, ok := keyed(i.GetName()); ok {
					_, seen := st.nonclustered[in.class]
					st.nonclustered[in.class] = !seen && strings.EqualFold(i.GetType(), "NONCLUSTERED")
				}
				st.indexes.add(i.GetName(), kind)
				if kind == kindPK || kind == kindUQ {
					// A key's name is also a constraint's, in the schema;
					// an index's is the table's alone.
					st.constraints.add(i.GetName(), kind)
					objs.add(i.GetName(), kind)
				}
			}
			for _, fk := range t.GetForeignKeys() {
				st.constraints.add(fk.GetName(), kindFK)
				objs.add(fk.GetName(), kindFK)
				for _, c := range fk.GetColumns() {
					st.fkColumns[c] = fk.GetName()
				}
				if tn, ok := keyed(t.GetName()); ok {
					refs = append(refs, fkRef{from: fkSource{walk.Path(sc, tn.class), schema.GetName() + "." + t.GetName()}, schema: fk.GetReferencedSchema(), table: fk.GetReferencedTable()})
				}
			}
			for _, ck := range t.GetCheckConstraints() {
				st.constraints.add(ck.GetName(), kindCK)
				objs.add(ck.GetName(), kindCK)
			}
			for _, tr := range t.GetTriggers() {
				objs.add(tr.GetName(), kindTrigger)
				s.triggerOnRelation(sc, tr.GetName(), schema.GetName(), t.GetName())
			}
		}
		for _, v := range schema.GetViews() {
			objs.add(v.GetName(), kindView)
			rels.add(v.GetName(), kindView)
			if vn, ok := keyed(v.GetName()); ok {
				path := walk.Path(sc, vn.class)
				if _, seen := s.tables[path]; seen || views[path] {
					s.ambiguous[path] = true
				}
				views[path] = true
				if viewUnbound(v.GetDefinition()) {
					s.unbound[path] = true
				}
				s.viewCols[path] = syncedViewCols(v)
			}
			for _, tr := range v.GetTriggers() {
				objs.add(tr.GetName(), kindTrigger)
				s.triggerOnRelation(sc, tr.GetName(), schema.GetName(), v.GetName())
			}
		}
		for _, q := range schema.GetSequences() {
			objs.add(q.GetName(), kindSequence)
			rels.add(q.GetName(), kindSequence)
		}
		for _, sy := range schema.GetSynonyms() {
			objs.add(sy.GetName(), kindSynonym)
			rels.add(sy.GetName(), kindSynonym)
		}
		for _, f := range schema.GetFunctions() {
			objs.add(f.GetName(), kindFunction)
			if fn, ok := keyed(f.GetName()); ok {
				path := walk.Path(sc, fn.class)
				if _, seen := s.functions[path]; seen {
					s.functions[path] = nil
				} else {
					s.functions[path] = syncedFunc(f.GetDefinition())
				}
			}
		}
		for _, p := range schema.GetProcedures() {
			objs.add(p.GetName(), kindProcedure)
		}
	}
	for _, ref := range refs {
		sn, ok1 := keyed(ref.schema)
		tn, ok2 := keyed(ref.table)
		if !ok1 || !ok2 {
			continue
		}
		to := walk.Path(sn.class, tn.class)
		if t := s.tables[to]; t != nil && to != ref.from.path {
			t.referencedBy = append(t.referencedBy, ref.from)
		}
	}
	return s
}

// triggerOnRelation records that a synced trigger of schema class sc is
// on relation schema.rel.
func (s *mssqlSnapshot) triggerOnRelation(sc, trigger, schema, rel string) {
	tn, ok1 := keyed(trigger)
	sn, ok2 := keyed(schema)
	rn, ok3 := keyed(rel)
	if !ok1 || !ok2 || !ok3 || !sn.plain || !rn.plain {
		// Only plain names of different classes differ under every
		// collation.
		return
	}
	path := walk.Path(sc, tn.class)
	if _, seen := s.triggerOn[path]; seen {
		s.triggerOn[path] = fkSource{}
		return
	}
	s.triggerOn[path] = fkSource{walk.Path(sc, rn.class), schema + "." + rel}
}

// Lookup implements walk.Snapshot.
func (s *mssqlSnapshot) Lookup(k walk.Key) (walk.Fact, string) {
	path := walk.Names(k.Parent)
	switch k.Space {
	case spaceSchema:
		class, _, _ := strings.Cut(k.Name, exactSep)
		if s.excluded[class] {
			return walk.Unknown, ""
		}
		return s.schemas.lookup(k.Name, true)
	case spaceTemp:
		// The change runs on a session of its own.
		return walk.Absent, ""
	}
	if len(path) == 0 || path[0] == tempSchema || s.excluded[path[0]] || s.ambiguous[path[0]] {
		return walk.Unknown, ""
	}
	switch k.Space {
	case spaceObject:
		if len(path) != 1 {
			return walk.Unknown, ""
		}
		return s.objects[path[0]].lookup(k.Name, false)
	case spaceRelation:
		if len(path) != 1 {
			return walk.Unknown, ""
		}
		if _, synced := s.objects[path[0]]; !synced {
			// A schema the sync does not list exists only as one it
			// excludes or of a marked name, else it holds nothing.
			if f, _ := s.schemas.lookup(path[0], true); f == walk.Absent && !s.excluded[path[0]] {
				return walk.Absent, ""
			}
			return walk.Unknown, ""
		}
		return s.relations[path[0]].lookup(k.Name, true)
	case spaceColumn, spaceIndex, spaceConstraint, spacePrimary:
		if len(path) != 2 || s.ambiguous[k.Parent] {
			return walk.Unknown, ""
		}
		t := s.tables[k.Parent]
		if t == nil {
			if k.Space == spaceIndex && s.unbound[k.Parent] {
				return walk.Absent, ""
			}
			return walk.Unknown, ""
		}
		switch k.Space {
		case spaceColumn:
			return t.columns.lookup(k.Name, true)
		case spaceIndex:
			return t.indexes.lookup(k.Name, true)
		case spaceConstraint:
			return t.constraints.lookup(k.Name, !t.edge)
		default:
			if t.primary {
				return walk.Present, ""
			}
			return walk.Absent, ""
		}
	}
	return walk.Unknown, ""
}
