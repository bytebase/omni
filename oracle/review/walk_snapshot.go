package review

import (
	"slices"
	"strconv"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review/walk"
)

// The walk's spaces (docs/specs/2026-10-26-walk-through-oracle.md §2). The
// object space is a schema's object namespace, which tables, views,
// materialized views, sequences, private synonyms, routines, packages, and
// types share; its key's kind is what the object is. The table, view,
// sequence, and routine spaces repeat what the object space holds for those
// kinds, so that a statement of one kind can trust the sync where it is
// complete for that kind. The routine space holds functions, procedures,
// and packages, which the sync lists in full from ALL_SOURCE. A column, the
// primary key, and a constraint as DROP CONSTRAINT finds it on its table
// are keys of the table's path; a constraint name, an index, and a trigger
// are keys of the schema, each in a namespace of its own, of the kind of
// the table path they are on. The
// synonym space holds the synonyms a name resolves through, private ones
// in their schema and public ones in PUBLIC; only DML reads it (spec §3).
const (
	spaceObject          = "object"
	spaceTable           = "table"
	spaceView            = "view"
	spaceSequence        = "sequence"
	spaceRoutine         = "routine"
	spaceSynonym         = "synonym"
	spaceColumn          = "column"
	spacePrimary         = "primary key"
	spaceTableConstraint = "table constraint"
	spaceConstraint      = "constraint"
	spaceIndex           = "index"
	spaceTrigger         = "trigger"
)

// routineKinds are the kinds the routine space holds.
var routineKinds = []string{kindFunction, kindProcedure, kindPackage}

// publicSchema owns the public synonyms. No user can take its name
// (01920), so no schema of the walk is named so.
const publicSchema = "PUBLIC"

// Object kinds, as Oracle names them.
const (
	kindTable     = "TABLE"
	kindView      = "VIEW"
	kindMView     = "MATERIALIZED VIEW"
	kindSequence  = "SEQUENCE"
	kindSynonym   = "SYNONYM"
	kindProcedure = "PROCEDURE"
	kindFunction  = "FUNCTION"
	kindPackage   = "PACKAGE"
	kindType      = "TYPE"
)

// generated reports whether a name is one the walk never claims in a
// space: a name Oracle generates, SYS_ for constraints, indexes, LOB
// segments, hidden columns, and IOT tables, ISEQ$$_ for identity
// sequences, BIN$ for the recycle bin, and ORA$ for private temporary
// tables; and, among columns, the hidden columns Oracle adds under names
// of its own, ORA_ARCHIVE_STATE for row archival and ORABCTAB_..$ for
// blockchain and immutable tables (verified: DDL on them fails with
// 38398 and 05715), which the sync does not list.
func generated(space, name string) bool {
	for _, p := range []string{"SYS_", "ISEQ$$_", "BIN$", "ORA$"} {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return space == spaceColumn && (strings.HasPrefix(name, "ORA_") || strings.HasPrefix(name, "ORABCTAB_"))
}

// tablePath is the parent of a table's columns, primary key, and
// constraints, and the kind of its indexes, constraint names, and
// triggers.
func tablePath(schema, name string) string { return walk.Path(schema, name) }

// oracleSnapshot is a target's synced schema as WalkThrough trusts it: one
// schema, the database the target names (spec §3). It answers Present for
// what Bytebase's sync lists, and Absent where that sync lists every object
// of a space: tables and materialized views, views, sequences, and
// functions, procedures, and packages, whose source ALL_SOURCE holds; the
// columns of a table whose columns it lists in full; whether a table has a
// primary key; and the synonyms a name resolves through, the schema's and
// the public ones. Constraint names, indexes, and triggers it lists only
// in part, and the object namespace holds types and other objects the
// sync does not list, so their absence is Unknown. Anything of another
// schema is Unknown.
type oracleSnapshot struct {
	schema  string
	objects map[string]string
	tables  map[string]*syncedTable
	// synonyms and publicSynonyms hold the names of the schema's private
	// synonyms and of the public synonyms. synonymTargets and
	// publicSynonymTargets map a local (no DBLink) synonym to its target
	// (schema, name) when synced.
	synonyms             map[string]bool
	publicSynonyms       map[string]bool
	synonymTargets       map[string][2]string
	publicSynonymTargets map[string][2]string
	// constraints, indexes, and triggers map a name to the table or view
	// it is on. constraintIndexes holds the names of the primary and unique
	// constraints, whose indexes may have those names, another, or none.
	constraints       map[string]string
	indexes           map[string]string
	indexTypes        map[string]string
	constraintIndexes map[string]bool
	triggers          map[string]string
	// domainTables holds the tables and materialized views with a domain
	// index, and domainIndexes those indexes: an indextype's code runs on
	// the table's DDL and DML, and can create and drop tables of its own,
	// as Oracle Text's DR$ tables. storage holds the tables the sync lists
	// no column of, a nested table's storage tables: DROP TABLE of the
	// parent, and DROP COLUMN of the nested table, drop them (verified).
	// logs holds the tables named MLOG$_, a materialized view log's
	// tables: DROP TABLE of the master drops its log's table and the
	// table's index (verified).
	domainTables  map[string]bool
	domainIndexes map[string]bool
	storage       []string
	logs          []string
}

// syncedCol holds column-level metadata populated by oracle/driver/sync.go.
type syncedCol struct {
	typ        string
	nullable   bool
	invisible  bool
	identity   bool
	virtual    bool
	hasDefault bool
}

// syncedTable is what the snapshot holds of a table or materialized view.
type syncedTable struct {
	mview    bool
	skipDump bool
	columns  map[string]bool
	// complete is set when the sync lists every column of the table.
	complete    bool
	primary     bool
	constraints map[string]bool

	// Rich metadata populated by oracle/driver/sync.go (and walkSync).
	// hasRichSync is true when synced columns carry non-empty Type strings,
	// distinguishing real Oracle sync output from minimal name-only test stubs.
	hasRichSync      bool
	colMeta          map[string]syncedCol
	visibleInsCols   int // count of non-invisible, non-virtual columns
	partType         metadata.TablePartitionMetadata_Type
	hasSubparts      bool
	hasInterval      bool
	hasForeignKey    bool
	partCount        int
	partNames        map[string]bool
	partKeyCols      map[string]bool
	partKeyCount     int
	allRangeInt      bool
	maxRangeHighVal  int64
	partRangeHigh    map[string]int64
	partRangePrev    map[string]int64
	partRangeHasPrev map[string]bool
	indexedCols      [][]string
}

func newOracleSnapshot(db *metadata.DatabaseSchemaMetadata) *oracleSnapshot {
	s := &oracleSnapshot{
		schema:               db.GetName(),
		objects:              make(map[string]string),
		tables:               make(map[string]*syncedTable),
		constraints:          make(map[string]string),
		indexes:              make(map[string]string),
		indexTypes:           make(map[string]string),
		constraintIndexes:    make(map[string]bool),
		triggers:             make(map[string]string),
		domainTables:         make(map[string]bool),
		domainIndexes:        make(map[string]bool),
		synonyms:             make(map[string]bool),
		publicSynonyms:       make(map[string]bool),
		synonymTargets:       make(map[string][2]string),
		publicSynonymTargets: make(map[string][2]string),
	}
	for _, syn := range db.GetPublicSynonyms() {
		s.publicSynonyms[syn.GetName()] = true
		if syn.GetTargetDbLink() == "" && syn.GetTargetName() != "" {
			s.publicSynonymTargets[syn.GetName()] = [2]string{syn.GetTargetSchema(), syn.GetTargetName()}
		}
	}
	for _, schema := range db.GetSchemas() {
		if name := schema.GetName(); name != "" && name != s.schema {
			continue
		}
		for _, t := range schema.GetTables() {
			s.objects[t.GetName()] = kindTable
			st := &syncedTable{
				skipDump:      t.GetSkipDump(),
				columns:       make(map[string]bool),
				constraints:   make(map[string]bool),
				colMeta:       make(map[string]syncedCol),
				hasForeignKey: len(t.GetForeignKeys()) > 0,
			}
			allTyped := len(t.GetColumns()) > 0
			for _, c := range t.GetColumns() {
				st.columns[c.GetName()] = true
				typ := baseOracleType(c.GetType())
				if typ == "" {
					allTyped = false
				}
				isVirt := c.GetGeneration() != nil && c.GetGeneration().GetType() != metadata.GenerationMetadata_TYPE_UNSPECIFIED
				sc := syncedCol{
					typ:        typ,
					nullable:   c.GetNullable(),
					invisible:  c.GetIsInvisible(),
					identity:   c.GetIsIdentity(),
					virtual:    isVirt,
					hasDefault: c.GetDefault() != "" || c.GetDefaultOnNull(),
				}
				st.colMeta[c.GetName()] = sc
				if !sc.invisible && !sc.virtual {
					st.visibleInsCols++
				}
			}
			st.hasRichSync = allTyped
			st.complete = columnsComplete(t, st.columns)
			if len(st.columns) == 0 {
				s.storage = append(s.storage, t.GetName())
			}
			if strings.HasPrefix(t.GetName(), "MLOG$_") {
				s.logs = append(s.logs, t.GetName())
			}
			s.partitionsOf(st, t.GetPartitions())
			s.indexesOf(t.GetName(), st, t.GetIndexes())
			for _, fk := range t.GetForeignKeys() {
				st.constraints[fk.GetName()] = true
				s.constraints[fk.GetName()] = t.GetName()
			}
			for _, ck := range t.GetCheckConstraints() {
				st.constraints[ck.GetName()] = true
				s.constraints[ck.GetName()] = t.GetName()
			}
			s.triggersOf(t.GetName(), t.GetTriggers())
			s.tables[t.GetName()] = st
		}
		for _, v := range schema.GetViews() {
			s.objects[v.GetName()] = kindView
			s.triggersOf(v.GetName(), v.GetTriggers())
		}
		for _, mv := range schema.GetMaterializedViews() {
			s.objects[mv.GetName()] = kindMView
			st := &syncedTable{mview: true, constraints: make(map[string]bool)}
			s.indexesOf(mv.GetName(), st, mv.GetIndexes())
			s.triggersOf(mv.GetName(), mv.GetTriggers())
			s.tables[mv.GetName()] = st
		}
		for _, seq := range schema.GetSequences() {
			s.objects[seq.GetName()] = kindSequence
		}
		for _, f := range schema.GetFunctions() {
			s.objects[f.GetName()] = kindFunction
		}
		for _, p := range schema.GetProcedures() {
			s.objects[p.GetName()] = kindProcedure
		}
		for _, p := range schema.GetPackages() {
			s.objects[p.GetName()] = kindPackage
		}
		for _, syn := range schema.GetSynonyms() {
			s.objects[syn.GetName()] = kindSynonym
			s.synonyms[syn.GetName()] = true
			if syn.GetTargetDbLink() == "" && syn.GetTargetName() != "" {
				s.synonymTargets[syn.GetName()] = [2]string{syn.GetTargetSchema(), syn.GetTargetName()}
			}
		}
	}
	return s
}

// baseOracleType normalizes a synced column type string (e.g. "VARCHAR2(20)",
// "NUMBER(10, 2)", "TIMESTAMP(6) WITH TIME ZONE") to its uppercased base name.
func baseOracleType(t string) string {
	t = strings.TrimSpace(t)
	if t == "" {
		return ""
	}
	if i := strings.IndexByte(t, '('); i >= 0 {
		t = strings.TrimSpace(t[:i])
	}
	return strings.ToUpper(t)
}

func (s *oracleSnapshot) partitionsOf(st *syncedTable, parts []*metadata.TablePartitionMetadata) {
	if len(parts) == 0 {
		return
	}
	st.partType = parts[0].GetType()
	st.partCount = len(parts)
	st.partNames = make(map[string]bool, len(parts))
	st.partKeyCols = make(map[string]bool)
	if expr := strings.TrimSpace(parts[0].GetExpression()); expr != "" {
		for _, col := range splitOracleIdentList(expr) {
			st.partKeyCols[col] = true
			st.partKeyCount++
		}
	}
	allInt := st.partType == metadata.TablePartitionMetadata_RANGE && st.partKeyCount == 1
	isSingleRange := allInt
	var maxVal int64
	first := true
	var prevVal int64
	hasPrev := false
	for _, p := range parts {
		if p.GetName() != "" {
			st.partNames[p.GetName()] = true
		}
		if len(p.GetSubpartitions()) > 0 {
			st.hasSubparts = true
		}
		if p.GetIntervalClause() != "" {
			st.hasInterval = true
		}
		if isSingleRange {
			v, err := strconv.ParseInt(strings.TrimSpace(p.GetValue()), 10, 64)
			if err == nil {
				if p.GetName() != "" {
					if st.partRangeHigh == nil {
						st.partRangeHigh = make(map[string]int64)
						st.partRangePrev = make(map[string]int64)
						st.partRangeHasPrev = make(map[string]bool)
					}
					st.partRangeHigh[p.GetName()] = v
					if hasPrev {
						st.partRangePrev[p.GetName()] = prevVal
						st.partRangeHasPrev[p.GetName()] = true
					}
				}
				if allInt && (first || v > maxVal) {
					maxVal = v
					first = false
				}
				prevVal = v
				hasPrev = true
			} else {
				if p.GetName() != "" && hasPrev {
					if st.partRangePrev == nil {
						st.partRangePrev = make(map[string]int64)
						st.partRangeHasPrev = make(map[string]bool)
					}
					st.partRangePrev[p.GetName()] = prevVal
					st.partRangeHasPrev[p.GetName()] = true
				}
				allInt = false
				hasPrev = false
			}
		}
	}
	st.allRangeInt = allInt && !first
	st.maxRangeHighVal = maxVal
}

// splitOracleIdentList splits a comma-separated list of Oracle identifiers
// (such as a synced partition key expression), keeping commas inside
// double-quoted identifiers intact and unescaping "" pairs.
func splitOracleIdentList(expr string) []string {
	var out []string
	var cur strings.Builder
	inQuote := false
	flush := func() {
		col := strings.TrimSpace(cur.String())
		cur.Reset()
		if len(col) >= 2 && col[0] == '"' && col[len(col)-1] == '"' {
			col = strings.ReplaceAll(col[1:len(col)-1], `""`, `"`)
		}
		if col != "" {
			out = append(out, col)
		}
	}
	for i := 0; i < len(expr); i++ {
		ch := expr[i]
		if ch == '"' {
			cur.WriteByte(ch)
			if inQuote && i+1 < len(expr) && expr[i+1] == '"' {
				cur.WriteByte(expr[i+1])
				i++
				continue
			}
			inQuote = !inQuote
			continue
		}
		if ch == ',' && !inQuote {
			flush()
			continue
		}
		cur.WriteByte(ch)
	}
	flush()
	return out
}

// indexesOf records a table's indexes. A primary key or unique constraint
// and the index that enforces it are two objects with two names, which the
// sync records as one entry under the index's (getConstraints,
// getIndexesAndConstraints), so the constraint is recorded under the name a
// DROP CONSTRAINT gives it.
func (s *oracleSnapshot) indexesOf(table string, st *syncedTable, indexes []*metadata.IndexMetadata) {
	for _, idx := range indexes {
		if idx.GetType() != "" {
			s.indexTypes[idx.GetName()] = idx.GetType()
		}
		if idx.GetType() == "DOMAIN" {
			s.domainTables[table] = true
			s.domainIndexes[idx.GetName()] = true
		} else if idx.GetVisible() && !strings.Contains(idx.GetType(), "REV") && !slices.Contains(idx.GetDescending(), true) {
			if exprs := idx.GetExpressions(); len(exprs) > 0 {
				cols := make([]string, 0, len(exprs))
				plain := true
				for _, e := range exprs {
					c := strings.TrimSpace(e)
					if len(c) >= 2 && c[0] == '"' && c[len(c)-1] == '"' {
						c = strings.ReplaceAll(c[1:len(c)-1], `""`, `"`)
					}
					if c == "" || strings.ContainsAny(c, "() ") || strings.HasPrefix(c, "SYS_NC") {
						plain = false
						break
					}
					cols = append(cols, c)
				}
				if plain {
					st.indexedCols = append(st.indexedCols, cols)
				}
			}
		}
		if !idx.GetIsConstraint() {
			s.indexes[idx.GetName()] = table
			continue
		}
		s.constraintIndexes[idx.GetName()] = true
		for _, constraint := range indexConstraints(idx) {
			name := constraint.GetName()
			st.constraints[name] = true
			s.constraints[name] = table
			if constraint.GetPrimary() {
				st.primary = true
			}
		}
	}
}

func (s *oracleSnapshot) triggersOf(table string, triggers []*metadata.TriggerMetadata) {
	for _, tr := range triggers {
		s.triggers[tr.GetName()] = table
	}
}

// columnsComplete reports whether the sync lists every column of a table
// whose columns are listed: the table is no object or XMLType table, which
// keep hidden columns of user names (XMLDATA, verified); it has a column,
// which the storage table of a nested table does not show (verified); and
// no check constraint names a column the sync does not list. A valid-time
// period adds hidden columns named after it, which the sync does not list,
// and a check constraint on them, which it does (verified: PERIOD FOR vt
// adds VT, VT_START, and VT_END, and a check "VT_START" < "VT_END").
func columnsComplete(t *metadata.TableMetadata, columns map[string]bool) bool {
	if t.GetSkipDump() || len(columns) == 0 {
		return false
	}
	for _, ck := range t.GetCheckConstraints() {
		for _, name := range quotedNames(ck.GetExpression()) {
			if !columns[name] {
				return false
			}
		}
	}
	return true
}

// quotedNames returns the double-quoted identifiers of a condition as
// Oracle stores it, outside its string literals.
func quotedNames(cond string) []string {
	var out []string
	for i := 0; i < len(cond); i++ {
		switch cond[i] {
		case '\'':
			for i++; i < len(cond); i++ {
				if cond[i] == '\'' {
					if i+1 < len(cond) && cond[i+1] == '\'' {
						i++
						continue
					}
					break
				}
			}
		case '"':
			end := strings.IndexByte(cond[i+1:], '"')
			if end < 0 {
				return out
			}
			out = append(out, cond[i+1:i+1+end])
			i += end + 1
		}
	}
	return out
}

// Lookup implements walk.Snapshot.
func (s *oracleSnapshot) Lookup(k walk.Key) (walk.Fact, string) {
	path := walk.Names(k.Parent)
	if k.Space == spaceSynonym && len(path) == 1 && path[0] == publicSchema {
		return synonymFact(s.publicSynonyms, k.Name)
	}
	if len(path) == 0 || path[0] != s.schema {
		return walk.Unknown, ""
	}
	switch k.Space {
	case spaceSynonym:
		return synonymFact(s.synonyms, k.Name)
	case spaceObject:
		if kind, ok := s.objects[k.Name]; ok {
			return walk.Present, kind
		}
		return walk.Unknown, ""
	case spaceTable, spaceView, spaceSequence, spaceRoutine:
		kind := s.objects[k.Name]
		switch {
		case k.Space == spaceTable && (kind == kindTable || kind == kindMView),
			k.Space == spaceView && kind == kindView,
			k.Space == spaceSequence && kind == kindSequence,
			k.Space == spaceRoutine && slices.Contains(routineKinds, kind):
			return walk.Present, kind
		case generated(k.Space, k.Name):
			return walk.Unknown, ""
		}
		return walk.Absent, ""
	case spaceColumn, spacePrimary, spaceTableConstraint:
		if len(path) != 2 {
			return walk.Unknown, ""
		}
		t, ok := s.tables[path[1]]
		if !ok || t.mview {
			return walk.Unknown, ""
		}
		switch k.Space {
		case spaceColumn:
			switch {
			case t.columns[k.Name]:
				return walk.Present, ""
			case t.complete && !generated(spaceColumn, k.Name):
				return walk.Absent, ""
			}
		case spacePrimary:
			if t.primary {
				return walk.Present, ""
			}
			return walk.Absent, ""
		case spaceTableConstraint:
			if t.constraints[k.Name] {
				return walk.Present, ""
			}
		}
		return walk.Unknown, ""
	case spaceConstraint:
		if table, ok := s.constraints[k.Name]; ok {
			return walk.Present, tablePath(s.schema, table)
		}
	case spaceIndex:
		if table, ok := s.indexes[k.Name]; ok && !s.constraintIndexes[k.Name] {
			return walk.Present, tablePath(s.schema, table)
		}
	case spaceTrigger:
		if table, ok := s.triggers[k.Name]; ok {
			return walk.Present, tablePath(s.schema, table)
		}
	}
	return walk.Unknown, ""
}

// synonymFact answers the synonym space from a list of synonyms. The sync
// leaves out a synonym whose object is gone or that the syncing user
// cannot reach, and DML through it fails too (00980, 00942, verified), so
// a name the list lacks resolves through no synonym.
func synonymFact(synonyms map[string]bool, name string) (walk.Fact, string) {
	switch {
	case synonyms[name]:
		return walk.Present, kindSynonym
	case generated(spaceSynonym, name):
		return walk.Unknown, ""
	}
	return walk.Absent, ""
}
