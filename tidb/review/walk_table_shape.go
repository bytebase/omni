package review

import (
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/tidb/ast"
)

// The walk follows a table's shape in keys under the table's path, each
// read only when Present: a column's type (its kind), the table's indexes
// with their columns (one key, its kind the list), its foreign keys with
// their columns and referenced table (one key each), the name the sync
// knew the table by, whether its AUTO_ID_CACHE is 1, and whether it is
// clustered by its primary key. The sync gives each; a statement that changes one in a way the walk does not follow
// leaves it Unknown, and a cleared table path leaves it Absent, which says
// nothing. The facts move with the table (walk.Move).
const (
	spaceColumnType = "column type"
	spaceIndexSet   = "index set"
	spaceFKDef      = "foreign key def"
	spaceOrigin     = "table origin"
	spaceIDCache    = "auto id cache"
	spaceClustered  = "clustered"
	shapeName       = "is"
)

func columnTypeKey(db, table, col string) walk.Key {
	return walk.Key{Space: spaceColumnType, Parent: tablePath(db, table), Name: col}
}

func indexSetKey(db, table string) walk.Key {
	return walk.Key{Space: spaceIndexSet, Parent: tablePath(db, table), Name: shapeName}
}

func fkDefKey(db, table, name string) walk.Key {
	return walk.Key{Space: spaceFKDef, Parent: tablePath(db, table), Name: name}
}

func originKey(db, table string) walk.Key {
	return walk.Key{Space: spaceOrigin, Parent: tablePath(db, table), Name: shapeName}
}

func idCacheKey(db, table string) walk.Key {
	return walk.Key{Space: spaceIDCache, Parent: tablePath(db, table), Name: shapeName}
}

// clusteredKey is Present, with kind clusteredYes, for a table the sync
// says is clustered by its primary key (TIDB_PK_TYPE), and nothing else: a
// CREATE TABLE never sets it, since whether TiDB clusters a key written
// without CLUSTERED or NONCLUSTERED is the session's
// tidb_enable_clustered_index. No ALTER TABLE makes a clustered table
// unclustered: TiDB refuses to drop its primary key.
func clusteredKey(db, table string) walk.Key {
	return walk.Key{Space: spaceClustered, Parent: tablePath(db, table), Name: shapeName}
}

const clusteredYes = "yes"

// shapeSpaces are the spaces of a table's shape, which DROP DATABASE
// clears.
var shapeSpaces = []string{spaceColumnType, spaceIndexSet, spaceFKDef, spaceOrigin, spaceIDCache, spaceCheckSet, spaceClustered, spacePartList}

// AUTO_ID_CACHE kinds: 1, or anything else (TiDB keeps the two apart).
const (
	idCacheOne   = "1"
	idCacheOther = "other"
)

// typeNames maps the type names the sync and the parser write to the base
// type the walk compares: an alias reads as the type TiDB stores.
var typeNames = map[string]string{
	"tinyint": "tinyint", "bool": "tinyint", "boolean": "tinyint",
	"smallint": "smallint", "mediumint": "mediumint", "int": "int", "integer": "int", "bigint": "bigint",
	"decimal": "decimal", "numeric": "decimal", "float": "float", "double": "double", "real": "double",
	"char": "char", "varchar": "varchar", "tinytext": "tinytext", "text": "text", "mediumtext": "mediumtext", "longtext": "longtext",
	"binary": "binary", "varbinary": "varbinary", "tinyblob": "tinyblob", "blob": "blob", "mediumblob": "mediumblob", "longblob": "longblob",
	"date": "date", "datetime": "datetime", "timestamp": "timestamp", "time": "time", "year": "year",
	"json": "json", "bit": "bit", "enum": "enum", "set": "set",
}

// intRanks orders the integer types by width.
var intRanks = map[string]int{"tinyint": 1, "smallint": 2, "mediumint": 3, "int": 4, "bigint": 5}

// columnType is a type as the walk keeps it: the base type, then
// " unsigned" for an unsigned number.
type columnType struct {
	base     string
	unsigned bool
}

func (t columnType) String() string {
	if t.unsigned {
		return t.base + " unsigned"
	}
	return t.base
}

func parseColumnType(kind string) (columnType, bool) {
	base, rest, _ := strings.Cut(kind, " ")
	if _, ok := typeNames[base]; !ok || base == "" {
		return columnType{}, false
	}
	return columnType{base: base, unsigned: rest == "unsigned"}, true
}

// syncedColumnType reads a synced column's type, information_schema's
// COLUMN_TYPE ("int(11) unsigned", "varchar(20)"); "" for one the walk
// does not read.
func syncedColumnType(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	end := strings.IndexAny(s, "( ")
	if end < 0 {
		end = len(s)
	}
	base, ok := typeNames[s[:end]]
	if !ok {
		return ""
	}
	t := columnType{base: base}
	if intRanks[base] > 0 || base == "decimal" || base == "float" || base == "double" {
		t.unsigned = strings.Contains(s[end:], "unsigned") || strings.Contains(s[end:], "zerofill")
	}
	return t.String()
}

// astColumnType reads a column definition's type, "" for one the walk does
// not read. SERIAL is BIGINT UNSIGNED.
func astColumnType(dt *ast.DataType) string {
	if dt == nil {
		return ""
	}
	name := strings.ToLower(dt.Name)
	if name == "serial" {
		return "bigint unsigned"
	}
	base, ok := typeNames[name]
	if !ok {
		return ""
	}
	t := columnType{base: base}
	if intRanks[base] > 0 || base == "decimal" || base == "float" || base == "double" {
		t.unsigned = dt.Unsigned || dt.Zerofill
	}
	return t.String()
}

// textTypes are the character string types.
var textTypes = map[string]bool{"char": true, "varchar": true, "tinytext": true, "text": true, "mediumtext": true, "longtext": true}

// fkIncompatible reports whether TiDB refuses a foreign key pairing
// columns of two types (3780, verified on v8.5.5): two integer types of
// different widths, whatever their sign or display width, or an integer
// with a decimal or a character string. The walk claims nothing of other
// pairs.
func fkIncompatible(a, b columnType) bool {
	ra, rb := intRanks[a.base], intRanks[b.base]
	switch {
	case ra > 0 && rb > 0:
		return ra != rb
	case ra > 0:
		return b.base == "decimal" || textTypes[b.base]
	case rb > 0:
		return a.base == "decimal" || textTypes[a.base]
	}
	return false
}

// rewritesData reports whether a column's change of type makes TiDB
// rewrite its data (needChangeColumnData in v8.5.5's ddl): an integer made
// narrower or given the other sign, or one made a character string. The
// walk claims nothing of other changes.
func rewritesData(from, to columnType) bool {
	rf, rt := intRanks[from.base], intRanks[to.base]
	switch {
	case rf > 0 && rt > 0:
		return rt < rf || from.unsigned != to.unsigned
	case rf > 0:
		return textTypes[to.base]
	}
	return false
}

// The kinds of the index set and of a foreign key are lists, written with
// separators no name the walk reads holds.
const (
	sepEntry = "\x1d"
	sepField = "\x1f"
	sepItem  = "\x1e"
)

// indexShape is one index of a table: its key, whether it is the primary
// key, unique, or global, and its key parts, "" for one that is no column.
type indexShape struct {
	name                    string
	primary, unique, global bool
	cols                    []string
}

func (i indexShape) has(col string) bool {
	for _, c := range i.cols {
		if c == col {
			return true
		}
	}
	return false
}

// hasAll reports whether every column of cols is one of the index's key
// parts.
func (i indexShape) hasAll(cols []string) bool {
	for _, c := range cols {
		if !i.has(c) {
			return false
		}
	}
	return true
}

func encodeIndexes(ix []indexShape) string {
	parts := make([]string, len(ix))
	for n, i := range ix {
		flags := ""
		if i.primary {
			flags += "p"
		}
		if i.unique {
			flags += "u"
		}
		if i.global {
			flags += "g"
		}
		parts[n] = i.name + sepField + flags + sepField + strings.Join(i.cols, sepItem)
	}
	return strings.Join(parts, sepEntry)
}

func decodeIndexes(kind string) []indexShape {
	if kind == "" {
		return nil
	}
	var out []indexShape
	for _, e := range strings.Split(kind, sepEntry) {
		f := strings.Split(e, sepField)
		if len(f) != 3 {
			return nil
		}
		out = append(out, indexShape{
			name:    f[0],
			primary: strings.Contains(f[1], "p"),
			unique:  strings.Contains(f[1], "u"),
			global:  strings.Contains(f[1], "g"),
			cols:    strings.Split(f[2], sepItem),
		})
	}
	return out
}

// plainName returns a synced name as the walk keys it, false for one that
// holds a separator or is no plain identifier (a functional key part).
func plainName(space, s string) (string, bool) {
	if s == "" || strings.ContainsAny(s, "\x1d\x1e\x1f`()' ") {
		return "", false
	}
	return nameKey(space, s, true)
}

// syncedIndexes reads a synced table's indexes, false when one of them has
// a name the walk does not read.
func syncedIndexes(t *metadata.TableMetadata) (string, bool) {
	var ix []indexShape
	for _, i := range t.GetIndexes() {
		s := indexShape{primary: i.GetPrimary(), unique: i.GetUnique() || i.GetPrimary(), global: i.GetIsGlobal()}
		if s.primary {
			s.name = primaryIndex
		} else {
			k, ok := plainName(spaceIndex, i.GetName())
			if !ok {
				return "", false
			}
			s.name = k
		}
		for _, e := range i.GetExpressions() {
			k, ok := plainName(spaceColumn, e)
			if !ok {
				k = ""
			}
			s.cols = append(s.cols, k)
		}
		ix = append(ix, s)
	}
	return encodeIndexes(ix), true
}

// fkShape is a foreign key: its columns, and the referenced table's
// database, its name as the sync knew it, and its columns, as the walk
// keys each.
type fkShape struct {
	cols            []string
	refDB, refTable string
	refCols         []string
}

func (f fkShape) encode() string {
	return strings.Join(f.cols, sepItem) + sepField + f.refDB + sepField + f.refTable + sepField + strings.Join(f.refCols, sepItem)
}

func decodeFK(kind string) (fkShape, bool) {
	p := strings.Split(kind, sepField)
	if len(p) != 4 {
		return fkShape{}, false
	}
	f := fkShape{cols: strings.Split(p[0], sepItem), refDB: p[1], refTable: p[2], refCols: strings.Split(p[3], sepItem)}
	return f, len(f.cols) == len(f.refCols) && len(f.cols) > 0
}

// syncedFK reads a synced foreign key of a table in database db, false
// when a name of it is one the walk does not read.
func syncedFK(db string, fk *metadata.ForeignKeyMetadata) (fkShape, bool) {
	f := fkShape{refDB: db}
	if s := fk.GetReferencedSchema(); s != "" {
		k, ok := plainName(spaceDatabase, s)
		if !ok {
			return fkShape{}, false
		}
		f.refDB = k
	}
	t, ok := plainName(spaceRelation, fk.GetReferencedTable())
	if !ok {
		return fkShape{}, false
	}
	f.refTable = t
	for _, c := range fk.GetColumns() {
		k, ok := plainName(spaceColumn, c)
		if !ok {
			return fkShape{}, false
		}
		f.cols = append(f.cols, k)
	}
	for _, c := range fk.GetReferencedColumns() {
		k, ok := plainName(spaceColumn, c)
		if !ok {
			return fkShape{}, false
		}
		f.refCols = append(f.refCols, k)
	}
	if len(f.cols) == 0 || len(f.cols) != len(f.refCols) {
		return fkShape{}, false
	}
	return f, true
}

// syncedFKRef is a foreign key the sync lists: its table and its key.
type syncedFKRef struct {
	child, name string
}

// readShape records what the sync says of one table's shape.
func (s *tidbSnapshot) readShape(tk string, t *metadata.TableMetadata) {
	types := make(map[string]string)
	for _, c := range t.GetColumns() {
		if k, ok := plainName(spaceColumn, c.GetName()); ok {
			if ct := syncedColumnType(c.GetType()); ct != "" {
				types[k] = ct
			}
		}
	}
	s.colTypes[tk] = types
	if ix, ok := syncedIndexes(t); ok {
		s.indexSets[tk] = ix
	}
	if cs, ok := syncedChecks(t); ok {
		s.checkSets[tk] = cs
	}
	defs := make(map[string]string)
	for _, fk := range t.GetForeignKeys() {
		name, ok := plainName(spaceTableForeignKey, fk.GetName())
		if !ok {
			continue
		}
		if f, ok := syncedFK(s.db, fk); ok {
			defs[name] = f.encode()
			s.fkRefs = append(s.fkRefs, syncedFKRef{child: tk, name: name})
		}
	}
	s.fkDefs[tk] = defs
	if t.GetAutoIdCache() == 1 {
		s.idCaches[tk] = idCacheOne
	} else {
		s.idCaches[tk] = idCacheOther
	}
	// A global temporary table fails a SHARD_ROW_ID_BITS for being
	// temporary first, with another error.
	if strings.EqualFold(t.GetPrimaryKeyType(), "CLUSTERED") && !t.GetIsGlobalTemporary() {
		s.clustered[tk] = true
	}
}

// lookupShape answers a key of a table's shape from the sync.
func (s *tidbSnapshot) lookupShape(k walk.Key, table string) (walk.Fact, string) {
	switch k.Space {
	case spaceColumnType:
		if t := s.colTypes[table][k.Name]; t != "" {
			return walk.Present, t
		}
	case spaceIndexSet:
		if ix, ok := s.indexSets[table]; ok && k.Name == shapeName {
			return walk.Present, ix
		}
	case spaceFKDef:
		if d := s.fkDefs[table][k.Name]; d != "" {
			return walk.Present, d
		}
	case spaceOrigin:
		if s.relations[table] == relTable && k.Name == shapeName {
			return walk.Present, table
		}
	case spaceIDCache:
		if c, ok := s.idCaches[table]; ok && k.Name == shapeName {
			return walk.Present, c
		}
	case spaceCheckSet:
		if cs, ok := s.checkSets[table]; ok && k.Name == shapeName {
			return walk.Present, cs
		}
	case spacePartList:
		if l, ok := s.partLists[table]; ok && k.Name == shapeName {
			return walk.Present, l
		}
	case spaceClustered:
		if s.clustered[table] && k.Name == shapeName {
			return walk.Present, clusteredYes
		}
	}
	return walk.Unknown, ""
}

// clustered reports whether the walk knows the table is clustered by its
// primary key before the statement.
func (b *builder) clustered(r tableRef) bool {
	f, kind := b.l.Lookup(clusteredKey(r.db, r.name))
	return f == walk.Present && kind == clusteredYes
}

// typeOf returns a column's type before the statement, false when the
// walk does not know it.
func (b *builder) typeOf(db, table, col string) (columnType, bool) {
	k, ok := b.key(spaceColumn, col)
	if !ok {
		return columnType{}, false
	}
	f, kind := b.l.Lookup(columnTypeKey(db, table, k))
	if f != walk.Present {
		return columnType{}, false
	}
	return parseColumnType(kind)
}

// indexesOf returns a table's indexes before the statement, false when the
// walk does not know them all.
func (b *builder) indexesOf(r tableRef) ([]indexShape, bool) {
	f, kind := b.l.Lookup(indexSetKey(r.db, r.name))
	if f != walk.Present {
		return nil, false
	}
	return decodeIndexes(kind), true
}

// createdShape records the shape a CREATE TABLE with a column list gives
// a table: its columns' types and its AUTO_ID_CACHE. Its indexes are
// left Absent, which the walk reads as unknown: TiDB names an unnamed key
// itself.
func (b *builder) createdShape(v *ast.CreateTableStmt, r tableRef) {
	for _, c := range v.Columns {
		k, ok := b.key(spaceColumn, c.Name)
		if !ok {
			continue
		}
		if ct := astColumnType(c.TypeName); ct != "" {
			b.Set(columnTypeKey(r.db, r.name, k), walk.Present, ct)
		}
	}
	cache := idCacheOther
	for _, opt := range v.Options {
		if opt != nil && strings.EqualFold(opt.Name, "AUTO_ID_CACHE") {
			cache = idCacheKind(opt.Value)
		}
	}
	if cache != "" {
		b.Set(idCacheKey(r.db, r.name), walk.Present, cache)
	}
}

// idCacheKind is the kind of an AUTO_ID_CACHE value, "" for one the walk
// does not read.
func idCacheKind(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	for i := 0; i < len(v); i++ {
		if v[i] < '0' || v[i] > '9' {
			return ""
		}
	}
	if strings.TrimLeft(v, "0") == "1" {
		return idCacheOne
	}
	return idCacheOther
}

// shapeEffects records what an ALTER TABLE does to the table's shape. A
// subcommand that may change the table's keys leaves its index set
// Unknown; one that renames a column leaves the foreign keys naming it
// Unknown, the table's own and, since the walk cannot find a referencing
// table once it moves, every other's (renamedColumns).
func (b *builder) shapeEffects(v *ast.AlterTableStmt, r tableRef) {
	// What the table's CHECK constraints read is followed no further than
	// the sync.
	b.Set(checkSetKey(r.db, r.name), walk.Unknown, "")
	b.Set(partListKey(r.db, r.name), walk.Unknown, "")
	unknownIndexes := func() { b.Set(indexSetKey(r.db, r.name), walk.Unknown, "") }
	setType := func(name string, dt *ast.DataType) {
		if k, ok := b.key(spaceColumn, name); ok {
			if ct := astColumnType(dt); ct != "" {
				b.Set(columnTypeKey(r.db, r.name, k), walk.Present, ct)
			} else {
				b.Set(columnTypeKey(r.db, r.name, k), walk.Unknown, "")
			}
		}
	}
	dropType := func(name string) {
		if k, ok := b.key(spaceColumn, name); ok {
			b.Set(columnTypeKey(r.db, r.name, k), walk.Unknown, "")
		}
	}
	renamed := false
	for _, c := range v.Commands {
		// MODIFY or CHANGE IF EXISTS of a column the walk doesn't know is
		// there may change nothing (skipMissingColumn): its columns' types
		// are left Unknown.
		if (c.Type == ast.ATModifyColumn || c.Type == ast.ATChangeColumn) && c.IfExists && c.Column != nil {
			source := c.Name
			if c.Type == ast.ATModifyColumn {
				source = c.Column.Name
			}
			if k, ok := b.key(spaceColumn, source); !ok || b.lookupFact(colKey(r.db, r.name, k)) != walk.Present {
				dropType(source)
				dropType(c.Column.Name)
				unknownIndexes()
				continue
			}
		}
		switch c.Type {
		case ast.ATAddColumn:
			cols := c.Columns
			if c.Column != nil {
				cols = append([]*ast.ColumnDef{c.Column}, cols...)
			}
			for _, col := range cols {
				if c.IfNotExists {
					dropType(col.Name)
				} else {
					setType(col.Name, col.TypeName)
				}
				if len(col.Constraints) > 0 {
					unknownIndexes()
				}
			}
			if len(c.Constraints) > 0 {
				unknownIndexes()
			}
		case ast.ATModifyColumn:
			if c.Column != nil {
				setType(c.Column.Name, c.Column.TypeName)
				if len(c.Column.Constraints) > 0 {
					unknownIndexes()
				}
			}
		case ast.ATChangeColumn:
			if c.Column == nil {
				continue
			}
			dropType(c.Name)
			setType(c.Column.Name, c.Column.TypeName)
			if from, ok1 := b.key(spaceColumn, c.Name); ok1 {
				if to, ok2 := b.key(spaceColumn, c.Column.Name); ok2 && from == to {
					if len(c.Column.Constraints) > 0 {
						unknownIndexes()
					}
					continue
				}
			}
			unknownIndexes()
			renamed = true
			b.forgetOwnFKs(r, c.Name)
		case ast.ATRenameColumn:
			from, ok1 := b.key(spaceColumn, c.Name)
			to, ok2 := b.key(spaceColumn, c.NewName)
			if ok1 && ok2 {
				f, kind := b.l.Lookup(columnTypeKey(r.db, r.name, from))
				b.Set(columnTypeKey(r.db, r.name, from), walk.Unknown, "")
				if f == walk.Present {
					b.Set(columnTypeKey(r.db, r.name, to), walk.Present, kind)
				} else {
					b.Set(columnTypeKey(r.db, r.name, to), walk.Unknown, "")
				}
			}
			unknownIndexes()
			renamed = true
			b.forgetOwnFKs(r, c.Name)
		case ast.ATDropColumn:
			dropType(c.Name)
			unknownIndexes()
			b.forgetOwnFKs(r, c.Name)
		case ast.ATAddIndex, ast.ATAddConstraint:
			if c.Constraint != nil && c.Constraint.Type == ast.ConstrCheck {
				continue
			}
			unknownIndexes()
		case ast.ATDropIndex, ast.ATRenameIndex, ast.ATPartitionBy, ast.ATRemovePartitioning:
			// Partitioning may make a key global or local.
			unknownIndexes()
		case ast.ATDropConstraint:
			k, ok := b.key(spaceTableForeignKey, c.Name)
			switch dropKindOf(b.w.sql, c) {
			case dropForeignKey:
				// The foreign key goes; its index stays (verified).
				if ok {
					b.Set(fkDefKey(r.db, r.name, k), walk.Absent, "")
				}
			case dropCheck:
			default:
				unknownIndexes()
				if ok {
					b.Set(fkDefKey(r.db, r.name, k), walk.Unknown, "")
				}
			}
		case ast.ATTableOption:
			if c.Option != nil && strings.EqualFold(c.Option.Name, "AUTO_ID_CACHE") {
				if kind := idCacheKind(c.Option.Value); kind != "" {
					b.Set(idCacheKey(r.db, r.name), walk.Present, kind)
				} else {
					b.Set(idCacheKey(r.db, r.name), walk.Unknown, "")
				}
			}
		}
	}
	if renamed {
		b.Session(func() { b.w.renamedColumns = true })
	}
}

// forgetOwnFKs leaves Unknown every foreign key of the table the sync
// lists under a name that may now name a foreign key of the column.
func (b *builder) forgetOwnFKs(r tableRef, col string) {
	k, ok := b.key(spaceColumn, col)
	if !ok {
		return
	}
	for _, ref := range b.w.snap.fkRefs {
		dk := fkDefKey(r.db, r.name, ref.name)
		f, kind := b.l.Lookup(dk)
		if f != walk.Present {
			continue
		}
		if fk, ok := decodeFK(kind); !ok || contains(fk.cols, k) {
			b.Set(dk, walk.Unknown, "")
		}
	}
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

// lookupFact returns a key's fact before the statement.
func (b *builder) lookupFact(k walk.Key) walk.Fact {
	f, _ := b.l.Lookup(k)
	return f
}
