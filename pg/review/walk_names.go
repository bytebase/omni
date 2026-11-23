package review

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// namedatalen is NAMEDATALEN - 1, the most bytes PostgreSQL keeps of a
// name: a longer one is truncated, quoted or not.
//
// pg: src/backend/parser/scansup.c — truncate_identifier
const namedatalen = 63

// walkName returns a name as PostgreSQL stores it, or false when the walk
// cannot tell. That is an empty name, or one with a byte outside ASCII,
// which PostgreSQL truncates at a character boundary (pg_mbcliplen) where
// this cuts at a byte.
func walkName(s string) (string, bool) {
	if s == "" {
		return "", false
	}
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 || s[i] == 0 {
			return "", false
		}
	}
	if len(s) > namedatalen {
		s = s[:namedatalen]
	}
	return s, true
}

func schemaKey(name string) walk.Key {
	return walk.Key{Space: spaceSchema, Name: name}
}

func relKey(schema, name string) walk.Key {
	return walk.Key{Space: spaceRelation, Parent: walk.Path(schema), Name: name}
}

func typeKey(schema, name string) walk.Key {
	return walk.Key{Space: spaceType, Parent: walk.Path(schema), Name: name}
}

func columnKey(t walk.Key, name string) walk.Key {
	return walk.Key{Space: spaceColumn, Parent: tableOf(t), Name: name}
}

func constraintKey(t walk.Key, name string) walk.Key {
	return walk.Key{Space: spaceConstraint, Parent: tableOf(t), Name: name}
}

func primaryKeyOf(t walk.Key) walk.Key {
	return walk.Key{Space: spacePrimaryKey, Parent: tableOf(t)}
}

func trigKey(t walk.Key, name string) walk.Key {
	return walk.Key{Space: spaceTrigger, Parent: tableOf(t), Name: name}
}

func polKey(t walk.Key, name string) walk.Key {
	return walk.Key{Space: spacePolicy, Parent: tableOf(t), Name: name}
}

func funcKey(schema, name string) walk.Key {
	return walk.Key{Space: spaceFunction, Parent: walk.Path(schema), Name: name}
}

func procKey(schema, name string) walk.Key {
	return walk.Key{Space: spaceProcedure, Parent: walk.Path(schema), Name: name}
}

// tableOf is the parent path of what a relation holds: its columns,
// constraints, and primary key.
func tableOf(rel walk.Key) string {
	return walk.Path(rel.Parent, rel.Name)
}

// schemaOfKey returns the schema holding a relation or a type.
func schemaOfKey(k walk.Key) string {
	return k.Parent
}

// hasRowType reports whether a relation of the kind has a row type, which
// takes its name in the schema's type namespace. Sequences and indexes do
// not (verified on PostgreSQL 17: CREATE TYPE q after CREATE SEQUENCE q
// succeeds).
func hasRowType(kind string) bool {
	switch kind {
	case relTable, relView, relMatview, relForeign, relComposite:
		return true
	}
	return false
}

// resolved is where a name resolves now: a fact and, when Present, the
// key it resolved to and its kind. Absent means the name resolves
// nowhere; missingSchema then names the schema the SQL wrote when it is
// the schema that certainly does not exist. fail is set when the name
// itself certainly fails, whatever exists: it names another database
// (CrossDatabaseReference). keys lists every key the answer read, which
// the within-statement rule compares against what the statement changes.
type resolved struct {
	fact          walk.Fact
	key           walk.Key
	kind          string
	missingSchema string
	fail          string
	builtin       bool
	keys          []walk.Key
}

// crossDatabase reports whether a name's catalog is certainly another
// database, which PostgreSQL refuses before it looks anything up, IF
// EXISTS or not (0A000, verified on PostgreSQL 17). A catalog naming the
// current database is allowed; the names compare as stored, so a longer
// spelling truncated to the current database's name is that database
// (verified on PostgreSQL 17).
//
// pg: src/backend/catalog/namespace.c — RangeVarGetRelidExtended,
// RangeVarGetCreationNamespace
func (w *walker) crossDatabase(catalog string) bool {
	if catalog == "" || w.database == "" {
		return false
	}
	n, ok := walkName(catalog)
	return ok && n != w.database
}

// currentDatabase reports whether a catalog the SQL wrote is certainly
// the current database, or there is none.
func (w *walker) currentDatabase(catalog string) bool {
	if catalog == "" {
		return true
	}
	n, ok := walkName(catalog)
	return ok && n == w.database
}

func crossDatabaseMessage(written string) string {
	return review.Code(written) + " refers to another database, which PostgreSQL doesn't allow."
}

// searchPath returns the current search path with "$user" expanded, or
// false when the path, or the user "$user" stands for, is unknown.
func (w *walker) searchPath() ([]string, bool) {
	if w.path == nil {
		return nil, false
	}
	path := make([]string, 0, len(w.path))
	for _, name := range w.path {
		if name == "$user" {
			if w.user == "" {
				return nil, false
			}
			name = w.user
		}
		n, ok := walkName(name)
		if !ok {
			return nil, false
		}
		path = append(path, n)
	}
	return path, true
}

// relation resolves a relation name as a statement that refers to an
// existing relation does: a schema the SQL wrote, or else the session's
// temporary schema, pg_catalog, and the search path, in the order
// PostgreSQL searches them. A name that resolves nowhere for certain is
// Absent; the first place the walk cannot judge makes the answer Unknown.
//
// pg: src/backend/catalog/namespace.c — RangeVarGetRelidExtended,
// recomputeNamespacePath (pg_temp, then pg_catalog, unless the path lists
// them)
func (b *builder) relation(catalog, schema, name string) resolved {
	if b.w.crossDatabase(catalog) {
		return resolved{fail: crossDatabaseMessage(qualified([]string{catalog, schema, name}))}
	}
	rel, ok := walkName(name)
	if !ok || !b.w.currentDatabase(catalog) {
		return resolved{}
	}
	if schema != "" {
		s, ok := walkName(schema)
		if !ok || s == "pg_catalog" {
			// pg_catalog the SQL names is not judged; see inSchema.
			return resolved{}
		}
		return b.inSchema(s, rel)
	}
	path, ok := b.w.searchPath()
	if !ok {
		return resolved{}
	}
	var keys []walk.Key
	search := func(s string) (resolved, bool) {
		r := b.inSchema(s, rel)
		keys = append(keys, r.keys...)
		switch {
		case r.fact == walk.Present:
			r.keys = keys
			return r, true
		case r.fact == walk.Unknown:
			return resolved{}, true
		}
		return resolved{}, false
	}
	order := make([]string, 0, len(path)+2)
	for _, implicit := range []string{tempSchema, "pg_catalog"} {
		if !containsString(path, implicit) {
			order = append(order, implicit)
		}
	}
	order = append(order, path...)
	for _, s := range order {
		if r, done := search(s); done {
			return r
		}
	}
	return resolved{fact: walk.Absent, keys: keys}
}

// inSchema resolves a relation in one schema. pg_catalog is never synced.
// Every relation it holds has a name beginning pg_ (verified on
// PostgreSQL 17): PostgreSQL refuses to create any other there, to a
// superuser and an extension script alike, unless a superuser sets
// allow_system_table_mods. The walk assumes that developer setting is
// off, so a name the search path looks up there implicitly, without pg_,
// is not there. Judging every unqualified name Unknown instead would cost
// every unqualified "does not exist". A name the SQL qualifies with
// pg_catalog is not judged at all.
func (b *builder) inSchema(schema, rel string) resolved {
	if schema == tempSchema {
		k := relKey(tempSchema, rel)
		f, kind := b.l.Lookup(k)
		return resolved{fact: f, key: k, kind: kind, keys: []walk.Key{k}}
	}
	if schema == "pg_catalog" {
		if strings.HasPrefix(rel, "pg_") {
			return resolved{}
		}
		return resolved{fact: walk.Absent}
	}
	sk := schemaKey(schema)
	k := relKey(schema, rel)
	switch f, _ := b.l.Lookup(sk); f {
	case walk.Absent:
		return resolved{fact: walk.Absent, missingSchema: schema, keys: []walk.Key{sk}}
	case walk.Unknown:
		// A relation the change made there certainly exists, and so does
		// the schema; anything else is unknown.
		if f, kind := b.l.Lookup(k); f == walk.Present {
			return resolved{fact: f, key: k, kind: kind, keys: []walk.Key{sk, k}}
		}
		return resolved{keys: []walk.Key{sk}}
	}
	f, kind := b.l.Lookup(k)
	return resolved{fact: f, key: k, kind: kind, keys: []walk.Key{sk, k}}
}

// rangeVar resolves a relation the SQL names with a RangeVar.
func (b *builder) rangeVar(rv *ast.RangeVar) resolved {
	if rv == nil {
		return resolved{}
	}
	return b.relation(rv.Catalogname, rv.Schemaname, rv.Relname)
}

// nameList resolves a relation the SQL names with a list of names, as
// DROP and COMMENT do.
func (b *builder) nameList(parts []string) resolved {
	switch len(parts) {
	case 1:
		return b.relation("", "", parts[0])
	case 2:
		return b.relation("", parts[0], parts[1])
	case 3:
		return b.relation(parts[0], parts[1], parts[2])
	}
	return resolved{}
}

// resolveType resolves a type name (1-, 2-, or 3-part) in the type namespace.
// Unqualified names follow PostgreSQL's search path (pg_temp, pg_catalog, then
// path), stopping at the first schema whose answer is Unknown so an unsynced
// domain, range, or base type in an earlier schema is never bypassed.
//
// pg: src/backend/catalog/namespace.c — TypenameGetTypidExtended
func (b *builder) resolveType(parts []string) resolved {
	var catalog, schema, name string
	switch len(parts) {
	case 1:
		name = parts[0]
	case 2:
		schema, name = parts[0], parts[1]
	case 3:
		catalog, schema, name = parts[0], parts[1], parts[2]
	default:
		return resolved{}
	}
	if b.w.crossDatabase(catalog) {
		return resolved{fail: crossDatabaseMessage(qualified(parts))}
	}
	typ, ok := walkName(name)
	if !ok || !b.w.currentDatabase(catalog) {
		return resolved{}
	}
	if schema != "" {
		s, ok := walkName(schema)
		if !ok {
			return resolved{}
		}
		if s == "pg_catalog" {
			if builtinCatalogTypes[typ] {
				return resolved{fact: walk.Present, kind: typOther, builtin: true}
			}
			return resolved{}
		}
		return b.typeInSchema(s, typ)
	}
	path, ok := b.w.searchPath()
	if !ok {
		return resolved{}
	}
	var keys []walk.Key
	order := make([]string, 0, len(path)+2)
	for _, implicit := range []string{tempSchema, "pg_catalog"} {
		if !containsString(path, implicit) {
			order = append(order, implicit)
		}
	}
	order = append(order, path...)
	for _, s := range order {
		if s == "pg_catalog" {
			if builtinCatalogTypes[typ] {
				return resolved{fact: walk.Present, kind: typOther, builtin: true, keys: keys}
			}
			if strings.HasPrefix(typ, "pg_") || strings.HasPrefix(typ, "_") {
				return resolved{}
			}
			continue
		}
		r := b.typeInSchema(s, typ)
		keys = append(keys, r.keys...)
		switch r.fact {
		case walk.Present:
			r.keys = keys
			return r
		case walk.Unknown:
			return resolved{}
		}
	}
	return resolved{fact: walk.Absent, keys: keys}
}

func (b *builder) typeInSchema(schema, typ string) resolved {
	if schema == tempSchema {
		k := typeKey(tempSchema, typ)
		f, kind := b.l.Lookup(k)
		return resolved{fact: f, key: k, kind: kind, keys: []walk.Key{k}}
	}
	sk := schemaKey(schema)
	k := typeKey(schema, typ)
	switch f, _ := b.l.Lookup(sk); f {
	case walk.Absent:
		return resolved{fact: walk.Absent, missingSchema: schema, keys: []walk.Key{sk}}
	case walk.Unknown:
		if f, kind := b.l.Lookup(k); f == walk.Present {
			return resolved{fact: f, key: k, kind: kind, keys: []walk.Key{sk, k}}
		}
		return resolved{keys: []walk.Key{sk}}
	}
	f, kind := b.l.Lookup(k)
	return resolved{fact: f, key: k, kind: kind, keys: []walk.Key{sk, k}}
}

// resolveRoutine resolves a function, procedure, or routine name (1-, 2-, or
// 3-part) in spaceFunction / spaceProcedure following PostgreSQL's search path.
// Unqualified lookup stops at the first Unknown schema so an unsynced or
// built-in routine is never bypassed.
//
// pg: src/backend/catalog/namespace.c — FuncnameGetCandidates
func (b *builder) resolveRoutine(parts []string, objType ast.ObjectType) resolved {
	var catalog, schema, name string
	switch len(parts) {
	case 1:
		name = parts[0]
	case 2:
		schema, name = parts[0], parts[1]
	case 3:
		catalog, schema, name = parts[0], parts[1], parts[2]
	default:
		if len(parts) > 3 {
			return resolved{fail: crossDatabaseMessage(qualified(parts))}
		}
		return resolved{}
	}
	if b.w.crossDatabase(catalog) {
		return resolved{fail: crossDatabaseMessage(qualified(parts))}
	}
	fn, ok := walkName(name)
	if !ok || !b.w.currentDatabase(catalog) {
		return resolved{}
	}
	if schema != "" {
		s, ok := walkName(schema)
		if !ok || s == "pg_catalog" {
			return resolved{}
		}
		return b.routineInSchema(s, fn, objType)
	}
	path, ok := b.w.searchPath()
	if !ok {
		return resolved{}
	}
	var keys []walk.Key
	order := make([]string, 0, len(path)+2)
	for _, implicit := range []string{tempSchema, "pg_catalog"} {
		if !containsString(path, implicit) {
			order = append(order, implicit)
		}
	}
	order = append(order, path...)
	for _, s := range order {
		r := b.routineInSchema(s, fn, objType)
		keys = append(keys, r.keys...)
		switch r.fact {
		case walk.Present:
			r.keys = keys
			return r
		case walk.Unknown:
			return resolved{}
		}
	}
	return resolved{fact: walk.Absent, keys: keys}
}

// routineInSchema resolves a routine in one schema. Like inSchema for
// relations, pg_catalog is never synced: PostgreSQL refuses to create user
// objects in pg_catalog unless allow_system_table_mods is enabled
// (pg: src/backend/catalog/namespace.c — QualifiedNameGetCreationNamespace).
// Assuming that developer setting is off, implicit pg_catalog lookup treats
// built-in/system routine names (pg_*, builtinCatalogFunctions, builtinCalls,
// IsSystemFunction) as Unknown and other names as Absent, while names
// explicitly qualified with pg_catalog are never judged (see resolveRoutine).
func (b *builder) routineInSchema(schema, name string, objType ast.ObjectType) resolved {
	if schema == "pg_catalog" {
		if strings.HasPrefix(name, "pg_") {
			return resolved{}
		}
		if objType == ast.OBJECT_PROCEDURE {
			return resolved{fact: walk.Absent}
		}
		if builtinCatalogFunctions[name] || !trustedRoutineName(name) {
			return resolved{}
		}
		return resolved{fact: walk.Absent}
	}
	if schema == tempSchema {
		return b.routineKeysInSchema(tempSchema, name, objType, nil)
	}
	sk := schemaKey(schema)
	switch f, _ := b.l.Lookup(sk); f {
	case walk.Absent:
		return resolved{fact: walk.Absent, missingSchema: schema, keys: []walk.Key{sk}}
	case walk.Unknown:
		r := b.routineKeysInSchema(schema, name, objType, []walk.Key{sk})
		if r.fact == walk.Present {
			return r
		}
		return resolved{keys: []walk.Key{sk}}
	}
	return b.routineKeysInSchema(schema, name, objType, []walk.Key{sk})
}

func (b *builder) routineKeysInSchema(schema, name string, objType ast.ObjectType, baseKeys []walk.Key) resolved {
	fk := funcKey(schema, name)
	pk := procKey(schema, name)
	switch objType {
	case ast.OBJECT_FUNCTION:
		f, _ := b.l.Lookup(fk)
		return resolved{fact: f, key: fk, keys: append(baseKeys, fk)}
	case ast.OBJECT_PROCEDURE:
		f, _ := b.l.Lookup(pk)
		return resolved{fact: f, key: pk, keys: append(baseKeys, pk)}
	default: // OBJECT_ROUTINE
		ff, _ := b.l.Lookup(fk)
		if ff == walk.Present {
			return resolved{fact: walk.Present, key: fk, keys: append(baseKeys, fk)}
		}
		pf, _ := b.l.Lookup(pk)
		if pf == walk.Present {
			return resolved{fact: walk.Present, key: pk, keys: append(baseKeys, pk)}
		}
		keys := append(baseKeys, fk, pk)
		if ff == walk.Absent && pf == walk.Absent {
			return resolved{fact: walk.Absent, keys: keys}
		}
		return resolved{keys: keys}
	}
}

// selfTypeImpossible reports whether a statement creating an object whose type
// occupies (schema, name) while also resolving tn as an existing type is
// impossible regardless of whether (schema, name) exists: if (schema, name)
// exists, the creation conflicts (42710/42P07/42809); if it does not exist, tn
// resolves nowhere (42704). It returns true only when every other schema where
// tn could be found is certainly devoid of it.
func (b *builder) selfTypeImpossible(schema, name string, tn *ast.TypeName) (bool, []walk.Key) {
	if tn == nil {
		return false, nil
	}
	parts := nameParts(tn.Names)
	var catalog, tSchema, tName string
	switch len(parts) {
	case 1:
		tName = parts[0]
	case 2:
		tSchema, tName = parts[0], parts[1]
	case 3:
		catalog, tSchema, tName = parts[0], parts[1], parts[2]
	default:
		return false, nil
	}
	if !b.w.currentDatabase(catalog) {
		return false, nil
	}
	norm, ok := walkName(tName)
	if !ok || norm != name {
		return false, nil
	}
	if tSchema != "" {
		s, ok := walkName(tSchema)
		return ok && s == schema, nil
	}
	path, ok := b.w.searchPath()
	if !ok {
		return false, nil
	}
	var keys []walk.Key
	order := make([]string, 0, len(path)+2)
	for _, implicit := range []string{tempSchema, "pg_catalog"} {
		if !containsString(path, implicit) {
			order = append(order, implicit)
		}
	}
	order = append(order, path...)
	for _, s := range order {
		if s == schema {
			continue
		}
		if s == "pg_catalog" {
			if builtinCatalogTypes[norm] || strings.HasPrefix(norm, "pg_") || strings.HasPrefix(norm, "_") {
				return false, nil
			}
			continue
		}
		if s == tempSchema {
			k := typeKey(tempSchema, norm)
			keys = append(keys, k)
			if f, _ := b.l.Lookup(k); f != walk.Absent {
				return false, nil
			}
			continue
		}
		sk := schemaKey(s)
		keys = append(keys, sk)
		if sf, _ := b.l.Lookup(sk); sf == walk.Absent {
			continue
		}
		tk := typeKey(s, norm)
		keys = append(keys, tk)
		if tf, _ := b.l.Lookup(tk); tf != walk.Absent {
			return false, nil
		}
	}
	return true, keys
}

// created is the schema a new object goes to. fail is set when the
// statement certainly fails for want of one, with the failure's code;
// fact is Unknown when the walk cannot tell where it goes.
type created struct {
	fact   walk.Fact
	schema string
	fail   string
	code   review.FindingCode
	keys   []walk.Key
}

// creationSchema returns the schema a new relation or type, name as
// written, goes to: the schema the SQL wrote, the session's temporary
// schema for a temporary relation, or else the first schema of the search
// path that exists. PostgreSQL never creates in the implicitly searched
// pg_temp or pg_catalog; a path that lists either is one the walk does
// not follow, as is a temporary relation the SQL puts in a schema of its
// own, which PostgreSQL refuses.
//
// pg: src/backend/catalog/namespace.c — RangeVarGetCreationNamespace,
// QualifiedNameGetCreationNamespace
func (b *builder) creationSchema(schema, name string, temp bool) created {
	if schema != "" {
		s, ok := walkName(schema)
		switch {
		case !ok, temp && s != tempSchema, s == "pg_catalog":
			return created{}
		case s == tempSchema:
			return created{fact: walk.Present, schema: tempSchema}
		}
		sk := schemaKey(s)
		switch f, _ := b.l.Lookup(sk); f {
		case walk.Absent:
			return created{fact: walk.Absent, fail: schemaNotFound(schema), code: review.ObjectNotFound, keys: []walk.Key{sk}}
		case walk.Unknown:
			// The object goes there if the statement succeeds, and then
			// the schema exists.
			b.set(sk, walk.Present, "")
		}
		return created{fact: walk.Present, schema: s, keys: []walk.Key{sk}}
	}
	if temp {
		return created{fact: walk.Present, schema: tempSchema}
	}
	path, ok := b.w.searchPath()
	if !ok {
		return created{}
	}
	var keys []walk.Key
	for _, s := range path {
		if s == tempSchema || s == "pg_catalog" {
			return created{}
		}
		sk := schemaKey(s)
		keys = append(keys, sk)
		switch f, _ := b.l.Lookup(sk); f {
		case walk.Present:
			return created{fact: walk.Present, schema: s, keys: keys}
		case walk.Unknown:
			return created{keys: keys}
		}
	}
	return created{fact: walk.Absent, fail: "No schema on the search path exists to create " + review.Code(ident(name)) + " in.", code: review.NoCurrentSchema, keys: keys}
}

// typeName splits a qualified type or object name list into schema and
// name.
func splitName(parts []string) (schema, name string, ok bool) {
	switch len(parts) {
	case 1:
		return "", parts[0], true
	case 2:
		return parts[0], parts[1], true
	}
	return "", "", false
}

// generatedPrefix is the start every name PostgreSQL generates for an
// object of a relation shares: an index, a key's index, a serial or
// identity column's sequence, or a constraint. makeObjectName truncates
// the longer of the relation and column parts first, so the relation's
// part keeps at least 25 bytes when its name is longer; 16 leaves room.
//
// pg: src/backend/commands/indexcmds.c — makeObjectName, ChooseRelationName
func generatedPrefix(rel string) string {
	if len(rel) > 16 {
		return rel[:16]
	}
	return rel
}

func containsString(list []string, s string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}
	return false
}
