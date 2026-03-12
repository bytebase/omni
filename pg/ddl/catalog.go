package ddl

import (
	"sort"
	"strings"
)

// typeKey identifies a type by namespace and name.
type typeKey struct {
	ns   uint32
	name string
}

// castKey identifies a cast by source and target type.
type castKey struct {
	source, target uint32
}

// operKey identifies an operator by name and operand types.
type operKey struct {
	name  string
	left  uint32
	right uint32
}

// Catalog is the in-memory PostgreSQL catalog.
type Catalog struct {
	oidGen *OIDGenerator

	// Schema indexes.
	schemas      map[uint32]*Schema
	schemaByName map[string]*Schema

	// Type indexes (built-in + user-created row types).
	typeByOID  map[uint32]*BuiltinType
	typeByName map[typeKey]*BuiltinType

	// Cast index.
	castIndex map[castKey]*BuiltinCast

	// Operator indexes.
	operByOID map[uint32]*BuiltinOperator
	operByKey map[operKey][]*BuiltinOperator

	// Proc indexes.
	procByOID  map[uint32]*BuiltinProc
	procByName map[string][]*BuiltinProc

	// Relation index.
	relationByOID map[uint32]*Relation

	// Constraint indexes.
	constraints map[uint32]*Constraint  // OID → Constraint
	consByRel   map[uint32][]*Constraint // relOID → constraints

	// Index indexes.
	indexes      map[uint32]*Index  // OID → Index
	indexesByRel map[uint32][]*Index // relOID → indexes

	// Sequence index.
	sequenceByOID map[uint32]*Sequence

	// Enum metadata.
	enumTypes map[uint32]*EnumType

	// Domain metadata.
	domainTypes map[uint32]*DomainType

	// Range type metadata.
	rangeTypes map[uint32]*RangeType

	// User-defined functions/procedures.
	userProcs map[uint32]*UserProc

	// Triggers.
	triggers      map[uint32]*Trigger
	triggersByRel map[uint32][]*Trigger

	// Comments.
	comments map[commentKey]string

	// Grants (privilege tracking).
	grants []Grant

	// Row-level security policies.
	policies      map[uint32]*Policy
	policiesByRel map[uint32][]*Policy

	// Extensions.
	extensions map[uint32]*Extension
	extByName  map[string]*Extension

	// Access methods, operator families, and operator classes.
	accessMethods      map[uint32]*AccessMethod
	accessMethodByName map[string]*AccessMethod
	opFamilies         map[uint32]*OpFamily
	opClasses          map[uint32]*OpClass
	opClassByKey       map[opClassKey]*OpClass // (amOID, typeOID) → default OpClass

	// Inheritance entries (pg_inherits).
	inhEntries []InhEntry

	// Dependencies.
	deps []DepEntry

	// Search path (schema names). pg_catalog is always searched implicitly.
	// Names are stored without validation — non-existent schemas are skipped
	// at lookup time, matching PostgreSQL's behavior.
	searchPath []string

	// Warnings collected during the most recent statement execution.
	warnings []Warning

	// Temporary CTE list for recursive CTE analysis. Set during analysis of
	// the recursive term so that nested analyzeSelectStmt calls can see the
	// partially-defined CTE.
	// pg: src/backend/parser/analyze.c — determineRecursiveColTypes
	visibleCTEs []*CommonTableExprQ
}

// New creates a fully initialized Catalog with all built-in data indexed.
func New() *Catalog {
	c := &Catalog{
		oidGen:        NewOIDGenerator(),
		schemas:       make(map[uint32]*Schema),
		schemaByName:  make(map[string]*Schema),
		typeByOID:     make(map[uint32]*BuiltinType, len(BuiltinTypes)),
		typeByName:    make(map[typeKey]*BuiltinType, len(BuiltinTypes)),
		castIndex:     make(map[castKey]*BuiltinCast, len(BuiltinCasts)),
		operByOID:     make(map[uint32]*BuiltinOperator, len(BuiltinOperators)),
		operByKey:     make(map[operKey][]*BuiltinOperator, len(BuiltinOperators)),
		procByOID:     make(map[uint32]*BuiltinProc, len(BuiltinProcs)),
		procByName:    make(map[string][]*BuiltinProc),
		relationByOID: make(map[uint32]*Relation),
		constraints:   make(map[uint32]*Constraint),
		consByRel:     make(map[uint32][]*Constraint),
		indexes:       make(map[uint32]*Index),
		indexesByRel:  make(map[uint32][]*Index),
		sequenceByOID: make(map[uint32]*Sequence),
		enumTypes:     make(map[uint32]*EnumType),
		domainTypes:   make(map[uint32]*DomainType),
		rangeTypes:    make(map[uint32]*RangeType),
		userProcs:     make(map[uint32]*UserProc),
		triggers:      make(map[uint32]*Trigger),
		triggersByRel:  make(map[uint32][]*Trigger),
		comments:       make(map[commentKey]string),
		policies:       make(map[uint32]*Policy),
		policiesByRel:  make(map[uint32][]*Policy),
		extensions:         make(map[uint32]*Extension),
		extByName:          make(map[string]*Extension),
		accessMethods:      make(map[uint32]*AccessMethod),
		accessMethodByName: make(map[string]*AccessMethod),
		opFamilies:         make(map[uint32]*OpFamily),
		opClasses:          make(map[uint32]*OpClass),
		opClassByKey:       make(map[opClassKey]*OpClass),
	}

	// Index types.
	for i := range BuiltinTypes {
		t := &BuiltinTypes[i]
		c.typeByOID[t.OID] = t
		c.typeByName[typeKey{ns: t.Namespace, name: t.TypeName}] = t
	}

	// Index casts.
	for i := range BuiltinCasts {
		cast := &BuiltinCasts[i]
		c.castIndex[castKey{source: cast.Source, target: cast.Target}] = cast
	}

	// Index operators.
	for i := range BuiltinOperators {
		op := &BuiltinOperators[i]
		c.operByOID[op.OID] = op
		key := operKey{name: op.Name, left: op.Left, right: op.Right}
		c.operByKey[key] = append(c.operByKey[key], op)
	}

	// Index procs.
	for i := range BuiltinProcs {
		p := &BuiltinProcs[i]
		c.procByOID[p.OID] = p
		c.procByName[p.Name] = append(c.procByName[p.Name], p)
	}

	// Create built-in schemas.
	c.addBuiltinSchema(PGCatalogNamespace, "pg_catalog")
	c.addBuiltinSchema(PGToastNamespace, "pg_toast")
	c.addBuiltinSchema(PublicNamespace, "public")

	// Register built-in access methods using PG's real OIDs.
	// pg: src/include/catalog/pg_am.dat
	for _, am := range []AccessMethod{
		{OID: 403, Name: "btree", Type: 'i'},
		{OID: 405, Name: "hash", Type: 'i'},
		{OID: 783, Name: "gist", Type: 'i'},
		{OID: 2742, Name: "gin", Type: 'i'},
		{OID: 4000, Name: "spgist", Type: 'i'},
		{OID: 3580, Name: "brin", Type: 'i'},
	} {
		a := am // copy for pointer stability
		c.accessMethods[a.OID] = &a
		c.accessMethodByName[a.Name] = &a
	}

	// Default search path: public, pg_catalog (implicit).
	c.searchPath = []string{"public"}

	return c
}

func (c *Catalog) addBuiltinSchema(oid uint32, name string) {
	s := &Schema{
		OID:       oid,
		Name:      name,
		Relations: make(map[string]*Relation),
		Indexes:   make(map[string]*Index),
		Sequences: make(map[string]*Sequence),
	}
	c.schemas[oid] = s
	c.schemaByName[name] = s
}

// ConstraintsOf returns all constraints on the given relation.
func (c *Catalog) ConstraintsOf(relOID uint32) []*Constraint {
	return c.consByRel[relOID]
}

// IndexesOf returns all indexes on the given relation.
func (c *Catalog) IndexesOf(relOID uint32) []*Index {
	return c.indexesByRel[relOID]
}

// GetIndexByOID returns the index with the given OID, or nil.
func (c *Catalog) GetIndexByOID(oid uint32) *Index {
	return c.indexes[oid]
}

// GetSequenceByOID returns the sequence with the given OID, or nil.
func (c *Catalog) GetSequenceByOID(oid uint32) *Sequence {
	return c.sequenceByOID[oid]
}

// TypeByOID returns the type with the given OID, or nil.
func (c *Catalog) TypeByOID(oid uint32) *BuiltinType {
	return c.typeByOID[oid]
}

// LookupCast returns the cast from source to target, or nil.
func (c *Catalog) LookupCast(source, target uint32) *BuiltinCast {
	return c.castIndex[castKey{source: source, target: target}]
}

// LookupOperatorExact returns the operators matching the exact signature.
func (c *Catalog) LookupOperatorExact(name string, left, right uint32) []*BuiltinOperator {
	return c.operByKey[operKey{name: name, left: left, right: right}]
}

// LookupProcByName returns all procs with the given name.
func (c *Catalog) LookupProcByName(name string) []*BuiltinProc {
	return c.procByName[name]
}

// LookupProcByOID returns the proc with the given OID, or nil.
func (c *Catalog) LookupProcByOID(oid uint32) *BuiltinProc {
	return c.procByOID[oid]
}

// GetSchema returns the schema with the given name, or nil.
func (c *Catalog) GetSchema(name string) *Schema {
	return c.schemaByName[name]
}

// GetRelationByOID returns the relation with the given OID, or nil.
func (c *Catalog) GetRelationByOID(oid uint32) *Relation {
	return c.relationByOID[oid]
}

// GetRelation returns the relation in the given schema, or nil.
// If schema is empty, the search path is used.
func (c *Catalog) GetRelation(schema, name string) *Relation {
	if schema != "" {
		s := c.schemaByName[schema]
		if s == nil {
			return nil
		}
		return s.Relations[name]
	}
	for _, nsOID := range c.searchPathWithCatalog() {
		s := c.schemas[nsOID]
		if s == nil {
			continue
		}
		if r := s.Relations[name]; r != nil {
			return r
		}
	}
	return nil
}

// SetSearchPath sets the schema search path by name.
// Non-existent schemas are accepted and silently skipped at lookup time,
// matching PostgreSQL's behavior.
//
// pg: src/backend/utils/init/postinit.c — InitializeSearchPath
func (c *Catalog) SetSearchPath(schemas []string) {
	c.searchPath = schemas
}

// searchPathWithCatalog resolves the search path to namespace OIDs,
// appending pg_catalog if not explicitly listed. Non-existent schemas
// are silently skipped.
func (c *Catalog) searchPathWithCatalog() []uint32 {
	hasPGCatalog := false
	for _, name := range c.searchPath {
		if name == "pg_catalog" {
			hasPGCatalog = true
			break
		}
	}

	names := c.searchPath
	if !hasPGCatalog {
		names = append(names, "pg_catalog")
	}

	oids := make([]uint32, 0, len(names))
	for _, name := range names {
		if s := c.schemaByName[name]; s != nil {
			oids = append(oids, s.OID)
		}
	}
	return oids
}

// resolveTargetSchema determines the schema for a DDL operation.
// If schemaName is given, it must exist. Otherwise, the first existing schema in the search path is used.
func (c *Catalog) resolveTargetSchema(schemaName string) (*Schema, error) {
	if schemaName != "" {
		s := c.schemaByName[schemaName]
		if s == nil {
			return nil, errUndefinedSchema(schemaName)
		}
		return s, nil
	}
	// First existing entry in search path.
	for _, name := range c.searchPath {
		if s := c.schemaByName[name]; s != nil {
			return s, nil
		}
	}
	return nil, errUndefinedSchema("(none)")
}

// addWarning appends a warning to the internal buffer.
func (c *Catalog) addWarning(code, message string) {
	c.warnings = append(c.warnings, Warning{Code: code, Message: message})
}

// UserSchemas returns all user-created schemas (excludes pg_catalog, pg_toast).
// Schemas are returned sorted by OID.
func (c *Catalog) UserSchemas() []*Schema {
	result := make([]*Schema, 0, len(c.schemas))
	for _, s := range c.schemas {
		if s.OID == PGCatalogNamespace || s.OID == PGToastNamespace {
			continue
		}
		result = append(result, s)
	}
	// Sort by OID for deterministic output.
	sort.Slice(result, func(i, j int) bool {
		return result[i].OID < result[j].OID
	})
	return result
}

// FormatType formats a type OID and typmod into a human-readable type name.
func (c *Catalog) FormatType(typeOID uint32, typmod int32) string {
	return c.formatType(typeOID, typmod)
}

// DrainWarnings returns all accumulated warnings and clears the buffer.
func (c *Catalog) DrainWarnings() []Warning {
	w := c.warnings
	c.warnings = nil
	return w
}

// regclassout formats a regclass value (relation name) for output,
// stripping schema qualification when the relation is visible via search_path.
//
// pg: src/backend/utils/adt/regproc.c — regclassout
func (c *Catalog) regclassout(name string) string {
	// Parse schema.name from the stored value.
	var schemaName, relName string
	if idx := strings.LastIndex(name, "."); idx >= 0 {
		schemaName = name[:idx]
		relName = name[idx+1:]
	} else {
		return name // no schema prefix, already unqualified
	}

	// Look up the relation by schema and name.
	rel := c.GetRelation(schemaName, relName)
	if rel == nil {
		return name // can't find it, keep as-is
	}

	// pg: regclassout — use RelationIsVisible to decide qualification
	if c.RelationIsVisible(rel.OID) {
		return relName
	}
	return name
}

// TypeIsVisible checks whether a type is visible in the current search path.
// A type is visible if its namespace is in the search path and no type with
// the same name appears in an earlier search path entry.
//
// pg: src/backend/catalog/namespace.c — TypeIsVisibleExt
func (c *Catalog) TypeIsVisible(typOID uint32) bool {
	t := c.typeByOID[typOID]
	if t == nil {
		return false
	}

	typNamespace := t.Namespace

	// Quick check: pg_catalog types are always visible.
	if typNamespace == PGCatalogNamespace {
		return true
	}

	// Check if the type's namespace is in the active search path at all.
	searchOIDs := c.searchPathWithCatalog()
	inPath := false
	for _, nsOID := range searchOIDs {
		if nsOID == typNamespace {
			inPath = true
			break
		}
	}
	if !inPath {
		return false
	}

	// It's in the path, but might be hidden by a same-named type earlier.
	// pg: TypeIsVisibleExt — slow check for conflicting types
	typName := t.TypeName
	for _, nsOID := range searchOIDs {
		if nsOID == typNamespace {
			// Found it first in path.
			return true
		}
		// Check if a different type with the same name exists in this namespace.
		if other := c.typeByName[typeKey{ns: nsOID, name: typName}]; other != nil {
			// Found something else first in path — hidden.
			return false
		}
	}

	return false
}

// LoadSQL parses SQL and executes all statements into a new Catalog.
// Returns the catalog and the first error encountered (if any).
// On error, the catalog reflects all statements that succeeded before the failure.
func LoadSQL(sql string) (*Catalog, error) {
	c := New()
	results, err := c.Exec(sql, nil)
	if err != nil {
		return nil, err
	}
	for _, r := range results {
		if r.Error != nil {
			return c, r.Error
		}
	}
	return c, nil
}

// RelationIsVisible checks whether a relation is visible in the current search
// path. A relation is visible if its schema is in the search path and no
// relation with the same name appears in an earlier search path entry.
//
// pg: src/backend/catalog/namespace.c — RelationIsVisibleExt
func (c *Catalog) RelationIsVisible(relOID uint32) bool {
	rel := c.relationByOID[relOID]
	if rel == nil {
		return false
	}
	if rel.Schema == nil {
		return false
	}

	relNamespace := rel.Schema.OID

	// Quick check: pg_catalog relations are always visible.
	if relNamespace == PGCatalogNamespace {
		return true
	}

	// Check if the relation's namespace is in the active search path.
	searchOIDs := c.searchPathWithCatalog()
	inPath := false
	for _, nsOID := range searchOIDs {
		if nsOID == relNamespace {
			inPath = true
			break
		}
	}
	if !inPath {
		return false
	}

	// Slow check for shadowing.
	// pg: RelationIsVisibleExt — check for conflicting relations
	relName := rel.Name
	for _, nsOID := range searchOIDs {
		if nsOID == relNamespace {
			return true
		}
		// Check if a relation with the same name exists in this namespace.
		schema := c.schemas[nsOID]
		if schema == nil {
			continue
		}
		for _, r := range schema.Relations {
			if r.Name == relName {
				return false
			}
		}
		// Also check sequences (they appear in pg_class too).
		for _, s := range schema.Sequences {
			if s.Name == relName {
				return false
			}
		}
	}

	return false
}
