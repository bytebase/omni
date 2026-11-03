package review

import (
	"strings"

	"github.com/bytebase/omni/metadata"
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
	primary                       bool
	// edge is set for a graph edge table, whose edge constraints the sync
	// skips.
	edge bool
	// ledger is set for a table that may be a ledger table.
	ledger bool
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
}

func newMSSQLSnapshot(db *metadata.DatabaseSchemaMetadata) *mssqlSnapshot {
	s := &mssqlSnapshot{
		db:        db.GetName(),
		excluded:  make(map[string]bool),
		objects:   make(map[string]*snapSpace),
		relations: make(map[string]*snapSpace),
		tables:    make(map[string]*snapTable),
		ambiguous: make(map[string]bool),
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
		s.objects[sc], s.relations[sc] = objs, rels
		for _, t := range schema.GetTables() {
			objs.add(t.GetName(), kindTable)
			rels.add(t.GetName(), kindTable)
			// A table on an ambiguous path keeps no children, but the
			// names they hold in the schema count all the same.
			st := &snapTable{ledger: ledgerTable(t)}
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
				if strings.HasPrefix(c.GetName(), "$edge_id") {
					st.edge = true
				}
				if df := c.GetDefaultConstraintName(); df != "" {
					st.constraints.add(df, kindDF)
					objs.add(df, kindDF)
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
			}
			for _, ck := range t.GetCheckConstraints() {
				st.constraints.add(ck.GetName(), kindCK)
				objs.add(ck.GetName(), kindCK)
			}
			for _, tr := range t.GetTriggers() {
				objs.add(tr.GetName(), kindTrigger)
			}
		}
		for _, v := range schema.GetViews() {
			objs.add(v.GetName(), kindView)
			rels.add(v.GetName(), kindView)
			for _, tr := range v.GetTriggers() {
				objs.add(tr.GetName(), kindTrigger)
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
		}
		for _, p := range schema.GetProcedures() {
			objs.add(p.GetName(), kindProcedure)
		}
	}
	return s
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
