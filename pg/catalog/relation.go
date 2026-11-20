package catalog

// InhEntry represents a row in pg_inherits.
type InhEntry struct {
	InhRelID  uint32 // child relation OID
	InhParent uint32 // parent relation OID
	InhSeqNo  int32  // order in INHERITS list (1-based)
}

// PartitionInfo stores partition key information for a partitioned table.
type PartitionInfo struct {
	Strategy   byte    // 'l'=list, 'r'=range, 'h'=hash
	KeyAttNums []int16 // partition key column attnums
	NKeyAttrs  int16
}

// PartitionBound stores partition bound for a partition child.
type PartitionBound struct {
	Strategy   byte // must match parent's strategy
	IsDefault  bool
	ListValues []string // deparsed datum strings (LIST)
	LowerBound []string // deparsed lower bound values (RANGE)
	UpperBound []string // deparsed upper bound values (RANGE)
	Modulus    int      // HASH modulus
	Remainder  int      // HASH remainder
}

// Column represents a table column.
type Column struct {
	AttNum          int16
	Name            string
	TypeOID         uint32
	TypeMod         int32
	NotNull         bool
	HasDefault      bool
	Default         string
	Len             int16
	ByVal           bool
	Align           byte
	Storage         byte
	Collation       uint32
	DefaultAnalyzed AnalyzedExpr // analyzed form of default expression (for Tier 2 deparse)
	Generated       byte         // 's' = stored generated, 0 = none
	Identity        byte         // 'a' = ALWAYS, 'd' = BY DEFAULT, 0 = none
	IsLocal         bool         // true if defined locally (not only inherited)
	InhCount        int          // number of inheritance ancestors that define this column
	Compression     byte         // attcompression: 'p'=pglz, 'l'=lz4, 0=default
	GenerationExpr  string       // expression text for generated columns
	CollationName   string       // explicit collation name (empty = type default)
	Ndims           int16        // attndims — declared array dimensions
}

// Relation represents a table, view, or other relation.
type Relation struct {
	OID           uint32
	Name          string
	Schema        *Schema
	RelKind       byte // 'r'=table, 'v'=view, 'c'=composite, 'p'=partitioned, 'm'=matview
	Columns       []*Column
	colByName     map[string]int
	RowTypeOID    uint32
	ArrayOID      uint32
	AnalyzedQuery *Query // non-nil for views (analyzed form, for deparse)

	// Persistence (pg_class.relpersistence).
	Persistence     byte   // 'p'=permanent, 'u'=unlogged, 't'=temp (default 'p')
	ReplicaIdentity byte   // 'd'=default, 'f'=full, 'n'=nothing, 'i'=index
	OfTypeOID       uint32 // reloftype — 0 if not a typed table

	// Inheritance (pg_inherits).
	InhParents []uint32 // parent relation OIDs
	InhCount   int      // number of inheritance parents

	// Partitioning.
	PartitionInfo  *PartitionInfo  // non-nil if PARTITION BY (partitioned table)
	PartitionBound *PartitionBound // non-nil if PARTITION OF (partition)
	PartitionOf    uint32          // parent partitioned table OID (0 if not a partition)

	// Row-level security.
	RowSecurity      bool // relrowsecurity
	ForceRowSecurity bool // relforcerowsecurity

	// View options.
	CheckOption byte // 0=none, 'l'=LOCAL, 'c'=CASCADED (views only)

	// ON COMMIT behavior for temp tables.
	OnCommit byte // 0=noop, 'p'=preserve rows, 'd'=delete rows, 'D'=drop

	// Tablespace is the tablespace the relation was created in, or "" for
	// the database default. It is not part of the schema a sync reads.
	Tablespace string
	// AccessMethod is the table access method (pg_class.relam), or "" when
	// unset (e.g. default heap or partitioned table without explicit USING).
	AccessMethod string
	// OwnerName is the role that owns the relation, for ALTER ... OWNED BY.
	OwnerName string

	// Relation flags (pg_class).
	Owner       uint32          // relowner
	HasRules    bool            // relhasrules
	Rules       map[string]byte // rule name → firing mode ('O', 'D', 'R', 'A')
	HasSubclass bool            // relhassubclass — has child tables/partitions
	IsPopulated bool            // relispopulated — matview populated state
	IsPartition bool            // relispartition
}

// relationNamespaceKind returns the relkind a schema's relation namespace
// holds a name under, or false when it holds none. PostgreSQL keeps tables,
// views, materialized views, composite types, sequences and indexes in one
// namespace (pg_class), so a name it holds is never free: 'S' is a sequence
// and 'i' an index, the kinds another map in this catalog stores.
func relationNamespaceKind(s *Schema, name string) (byte, bool) {
	if rel := s.Relations[name]; rel != nil {
		return rel.RelKind, true
	}
	if s.Indexes[name] != nil {
		return 'i', true
	}
	if s.Sequences[name] != nil {
		return 'S', true
	}
	return 0, false
}

// findNamedIndex returns the index a schema's relation namespace holds a
// name under. An index name PostgreSQL resolves to a table, view or sequence
// is 42809, not a missing index.
func (c *Catalog) findNamedIndex(schema *Schema, name string) (*Index, error) {
	if idx := schema.Indexes[name]; idx != nil {
		return idx, nil
	}
	if _, ok := relationNamespaceKind(schema, name); ok {
		return nil, errWrongObjectType(name, "an index")
	}
	return nil, errUndefinedObject("index", name)
}

// findIndexRelation locates an index by schema (or search path) and name as a
// relation in pg_class (used by ALTER INDEX and RENAME): a missing name is
// 42P01 (relation does not exist), while another relation kind is 42809.
//
// pg: src/backend/catalog/namespace.c — RangeVarGetRelidExtended + relkind check
func (c *Catalog) findIndexRelation(schemaName, relName string) (*Schema, *Index, error) {
	if schemaName != "" {
		s := c.schemaByName[schemaName]
		if s == nil {
			return nil, nil, errUndefinedSchema(schemaName)
		}
		if idx := s.Indexes[relName]; idx != nil {
			return s, idx, nil
		}
		if _, ok := relationNamespaceKind(s, relName); ok {
			return nil, nil, errWrongObjectType(relName, "an index")
		}
		return nil, nil, errUndefinedTable(relName)
	}
	for _, nsOID := range c.searchPathWithCatalog() {
		s := c.schemas[nsOID]
		if s == nil {
			continue
		}
		if idx := s.Indexes[relName]; idx != nil {
			return s, idx, nil
		}
		if _, ok := relationNamespaceKind(s, relName); ok {
			return nil, nil, errWrongObjectType(relName, "an index")
		}
	}
	return nil, nil, errUndefinedTable(relName)
}

// findRelation locates a relation by schema (or search path) and name. A
// name the namespace holds as another kind is 42809, not a missing relation:
// PostgreSQL finds the index or sequence and then says it is not the object
// the statement needs.
func (c *Catalog) findRelation(schemaName, relName string) (*Schema, *Relation, error) {
	if schemaName != "" {
		s := c.schemaByName[schemaName]
		if s == nil {
			return nil, nil, errUndefinedSchema(schemaName)
		}
		r := s.Relations[relName]
		if r == nil {
			if _, ok := relationNamespaceKind(s, relName); ok {
				return nil, nil, errWrongObjectType(relName, "a table")
			}
			return nil, nil, errUndefinedTable(relName)
		}
		return s, r, nil
	}
	for _, nsOID := range c.searchPathWithCatalog() {
		s := c.schemas[nsOID]
		if s == nil {
			continue
		}
		if r := s.Relations[relName]; r != nil {
			return s, r, nil
		}
		if _, ok := relationNamespaceKind(s, relName); ok {
			return nil, nil, errWrongObjectType(relName, "a table")
		}
	}
	return nil, nil, errUndefinedTable(relName)
}

// findRelByOID locates a relation by its OID across all schemas.
// (pgddl helper — PG uses RelationIdGetRelation)
func (c *Catalog) findRelByOID(oid uint32) *Relation {
	for _, s := range c.schemas {
		for _, r := range s.Relations {
			if r.OID == oid {
				return r
			}
		}
	}
	return nil
}

// removeRelation removes a relation and its associated types, constraints, indexes, and deps.
func (c *Catalog) removeRelation(schema *Schema, name string, rel *Relation) {
	// Remove own constraints and their backing indexes.
	c.removeConstraintsForRelation(rel.OID, schema)

	// Remove own indexes.
	c.removeIndexesForRelation(rel.OID, schema)

	// Remove triggers on this relation.
	c.removeTriggersForRelation(rel.OID)

	// Remove RLS policies on this relation.
	c.removePoliciesForRelation(rel.OID)

	// Remove grants on this relation.
	c.removeGrantsForObject('r', rel.OID)

	// Remove owned sequences (DepAuto from sequence to this relation).
	c.removeOwnedSequences(rel.OID)

	// Remove comments on this relation and its columns.
	c.removeComments('r', rel.OID)

	// Remove inheritance entries (both as child and as parent).
	c.removeInhEntries(rel.OID)

	// Remove all dependency entries for this relation.
	c.removeDepsOf('r', rel.OID)
	c.removeDepsOn('r', rel.OID)

	delete(schema.Relations, name)
	delete(c.relationByOID, rel.OID)
	delete(c.extensionMembers, rel.OID)

	// Remove row type.
	if rel.RowTypeOID != 0 {
		if rt := c.typeByOID[rel.RowTypeOID]; rt != nil {
			delete(c.typeByName, typeKey{ns: rt.Namespace, name: rt.TypeName})
		}
		for key, cast := range c.castIndex {
			if cast.Source == rel.RowTypeOID || cast.Target == rel.RowTypeOID {
				delete(c.castIndex, key)
				delete(c.extensionCasts, key)
			}
		}
		delete(c.typeByOID, rel.RowTypeOID)
		delete(c.typeOwners, rel.RowTypeOID)
		delete(c.extensionMembers, rel.RowTypeOID)
		c.removeGrantsForObject('T', rel.RowTypeOID)
		c.removeComments('t', rel.RowTypeOID)
		c.removeDepsOf('t', rel.RowTypeOID)
		c.removeDepsOn('t', rel.RowTypeOID)
	}

	// Remove array type.
	if rel.ArrayOID != 0 {
		if at := c.typeByOID[rel.ArrayOID]; at != nil {
			delete(c.typeByName, typeKey{ns: at.Namespace, name: at.TypeName})
		}
		for key, cast := range c.castIndex {
			if cast.Source == rel.ArrayOID || cast.Target == rel.ArrayOID {
				delete(c.castIndex, key)
				delete(c.extensionCasts, key)
			}
		}
		delete(c.typeByOID, rel.ArrayOID)
		delete(c.typeOwners, rel.ArrayOID)
		delete(c.extensionMembers, rel.ArrayOID)
		c.removeGrantsForObject('T', rel.ArrayOID)
		c.removeComments('t', rel.ArrayOID)
		c.removeDepsOf('t', rel.ArrayOID)
		c.removeDepsOn('t', rel.ArrayOID)
	}
}

// removeTriggersForRelation removes all triggers on a relation (and any
// constraint triggers on other relations referencing it via FROM).
func (c *Catalog) removeTriggersForRelation(relOID uint32) {
	for _, trig := range c.triggersByRel[relOID] {
		delete(c.triggers, trig.OID)
		delete(c.extensionMembers, trig.OID)
		c.removeComments('g', trig.OID)
		c.removeDepsOf('g', trig.OID)
		c.removeDepsOn('g', trig.OID)
	}
	delete(c.triggersByRel, relOID)

	var refTrigs []*Trigger
	for _, trig := range c.triggers {
		if trig.ConstraintRelOID == relOID && trig.RelOID != relOID {
			refTrigs = append(refTrigs, trig)
		}
	}
	for _, trig := range refTrigs {
		c.removeTrigger(trig)
	}
}

// removeInhEntries removes all pg_inherits entries where relOID is child or parent.
func (c *Catalog) removeInhEntries(relOID uint32) {
	n := 0
	for _, e := range c.inhEntries {
		if e.InhRelID == relOID || e.InhParent == relOID {
			continue
		}
		c.inhEntries[n] = e
		n++
	}
	c.inhEntries = c.inhEntries[:n]
}

// removeOwnedSequences removes sequences auto-dependent on a relation.
func (c *Catalog) removeOwnedSequences(relOID uint32) {
	for _, seq := range c.sequenceByOID {
		if seq.OwnerRelOID == relOID {
			c.removeSequence(seq.Schema, seq)
		}
	}
}
