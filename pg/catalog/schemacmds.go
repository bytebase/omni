package catalog

import (
	"fmt"
	"sort"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// CreateSchemaCommand creates a new schema in the catalog.
//
// pg: src/backend/commands/schemacmds.c — CreateSchemaCommand
func (c *Catalog) CreateSchemaCommand(stmt *nodes.CreateSchemaStmt) error {
	schemaName := stmt.Schemaname

	// If AUTHORIZATION is specified and no explicit schema name, use the role name.
	// pg: src/backend/commands/schemacmds.c — CreateSchemaCommand (authId / schemaName logic)
	var ownerName string
	if stmt.Authrole != nil {
		ownerName = c.ownerRoleName(stmt.Authrole)
		if schemaName == "" {
			schemaName = ownerName
		}
	}

	// Reject reserved schema name prefix.
	// pg: src/backend/commands/schemacmds.c — CreateSchemaCommand (pg_ prefix check)
	if strings.HasPrefix(schemaName, "pg_") {
		return &Error{
			Code:    CodeReservedName,
			Message: "unacceptable schema name \"" + schemaName + "\"\nDetail: The prefix \"pg_\" is reserved for system schemas.",
		}
	}

	if _, exists := c.schemaByName[schemaName]; exists {
		if stmt.IfNotExists {
			c.addWarning(CodeWarningSkip, fmt.Sprintf("schema %q already exists, skipping", schemaName))
			return nil
		}
		return errDuplicateSchema(schemaName)
	}
	var saved *Catalog
	if stmt.SchemaElts != nil {
		saved = c.Clone()
	}

	effectiveOwner := ownerName
	if effectiveOwner == "" {
		effectiveOwner = c.currentOwner()
	}
	oid := c.oidGen.Next()
	s := &Schema{
		OID:       oid,
		Name:      schemaName,
		Owner:     ownerName,
		OwnerName: effectiveOwner,
		Relations: make(map[string]*Relation),
		Indexes:   make(map[string]*Index),
		Sequences: make(map[string]*Sequence),
	}
	c.schemas[oid] = s
	c.schemaByName[schemaName] = s

	// Process schema elements (CREATE TABLE, CREATE VIEW, etc. inside CREATE SCHEMA).
	// pg: src/backend/commands/schemacmds.c — CreateSchemaCommand (schemaElts loop)
	if stmt.SchemaElts != nil {
		// Temporarily prepend the new schema to the search path so that
		// unqualified DDL targets land in the newly created schema, and run
		// under AUTHORIZATION role when specified.
		prevSearchPath := c.searchPath
		prevCurrentUser := c.currentUser
		c.searchPath = append([]string{schemaName}, prevSearchPath...)
		if ownerName != "" {
			c.currentUser = ownerName
		}

		for _, elt := range stmt.SchemaElts.Items {
			if elt == nil {
				continue
			}
			if err := c.ProcessUtility(elt); err != nil {
				*c = *saved
				return err
			}
		}
		c.searchPath = prevSearchPath
		c.currentUser = prevCurrentUser
	}

	// Roles are cluster-wide and a dump does not create them, so an
	// AUTHORIZATION role is recorded as observed rather than required — only
	// once the whole statement succeeded, as a rollback would discard it.
	if ownerName != "" {
		c.roles[ownerName] = true
	}

	return nil
}

// RemoveSchemas drops one or more schemas from the catalog.
// Called from RemoveObjects when the drop target is OBJECT_SCHEMA.
//
// pg: src/backend/commands/schemacmds.c — (drop case in RemoveObjects via dropcmds.c)
func (c *Catalog) RemoveSchemas(stmt *nodes.DropStmt) error {
	cascade := stmt.Behavior == int(nodes.DROP_CASCADE)
	if stmt.Objects == nil {
		return nil
	}

	var targets []*Schema
	seen := map[uint32]bool{}
	for _, obj := range stmt.Objects.Items {
		// Schema names may be a *nodes.String or a *nodes.List wrapping one.
		_, schemaName := extractDropObjectName(obj)

		s := c.schemaByName[schemaName]
		if s == nil {
			if stmt.Missing_ok {
				c.addWarning(CodeWarningSkip, fmt.Sprintf("schema %q does not exist, skipping", schemaName))
				continue
			}
			return errUndefinedSchema(schemaName)
		}

		// Prevent dropping built-in schemas.
		if s.OID == PGCatalogNamespace || s.OID == PGToastNamespace {
			return errInvalidParameterValue("cannot drop schema " + schemaName)
		}
		if c.isExtensionMember(s.OID) {
			return errDependentObjects("extension", c.extensionMembers[s.OID])
		}

		if !seen[s.OID] {
			seen[s.OID] = true
			targets = append(targets, s)
		}
	}

	if !cascade {
		for _, s := range targets {
			if c.schemaHasObjects(s) {
				return errSchemaNotEmpty(s.Name)
			}
		}
	}

	for _, s := range targets {
		// CASCADE: drop all contained objects. Extensions go first, through the
		// extension CASCADE path, so their dependents elsewhere (an index on
		// their operator class, a table using their type) are found before the
		// generic removal erases the member and its dependency edges.
		if cascade {
			c.cascadeDropSchemaContents(s)
		}
		c.deleteSchemaRecord(s)
	}
	return nil
}

func (c *Catalog) cascadeDropSchemaContents(s *Schema) {
	if s == nil || c.schemas[s.OID] == nil {
		return
	}
	var exts []*Extension
	for _, ext := range c.extensions {
		if ext.SchemaOID == s.OID {
			exts = append(exts, ext)
		}
	}
	sort.Slice(exts, func(i, j int) bool { return exts[i].OID < exts[j].OID })
	for _, ext := range exts {
		_ = c.removeExtensionObjects(&nodes.DropStmt{
			RemoveType: int(nodes.OBJECT_EXTENSION),
			Objects:    &nodes.List{Items: []nodes.Node{&nodes.String{Str: ext.Name}}},
			Behavior:   int(nodes.DROP_CASCADE),
		})
	}
	if c.schemas[s.OID] == nil {
		return
	}
	var rels []*Relation
	for _, rel := range s.Relations {
		rels = append(rels, rel)
	}
	sort.Slice(rels, func(i, j int) bool { return rels[i].OID < rels[j].OID })
	for _, rel := range rels {
		if c.relationByOID[rel.OID] != nil {
			c.cascadeDropRelation(rel)
		}
	}
	var idxs []*Index
	for _, idx := range s.Indexes {
		idxs = append(idxs, idx)
	}
	sort.Slice(idxs, func(i, j int) bool { return idxs[i].OID < idxs[j].OID })
	for _, idx := range idxs {
		if c.indexes[idx.OID] == idx {
			c.dropDependentObject('i', idx.OID)
		}
	}
	var seqs []*Sequence
	for _, seq := range s.Sequences {
		seqs = append(seqs, seq)
	}
	sort.Slice(seqs, func(i, j int) bool { return seqs[i].OID < seqs[j].OID })
	for _, seq := range seqs {
		if c.sequenceByOID[seq.OID] != nil {
			c.dropDependents('s', seq.OID)
			if c.sequenceByOID[seq.OID] != nil {
				c.removeSequence(s, seq)
			}
		}
	}
	c.removeFunctionsInSchema(s)
	c.removeTypesInSchema(s)
	var opOIDs []uint32
	for oid, nsOID := range c.operatorSchemas {
		if nsOID == s.OID {
			opOIDs = append(opOIDs, oid)
		}
	}
	sort.Slice(opOIDs, func(i, j int) bool { return opOIDs[i] < opOIDs[j] })
	for _, oid := range opOIDs {
		c.dropDependents('O', oid)
		c.removeOperator(oid)
	}
	opcOIDs := make([]uint32, 0, len(c.opClasses))
	for oid, opc := range c.opClasses {
		if opc.Schema == s.OID {
			opcOIDs = append(opcOIDs, oid)
		}
	}
	sort.Slice(opcOIDs, func(i, j int) bool { return opcOIDs[i] < opcOIDs[j] })
	for _, oid := range opcOIDs {
		c.cascadeDropOpClass(oid)
	}
	famOIDs := make([]uint32, 0, len(c.opFamilies))
	for oid, fam := range c.opFamilies {
		if fam.Schema == s.OID {
			famOIDs = append(famOIDs, oid)
		}
	}
	sort.Slice(famOIDs, func(i, j int) bool { return famOIDs[i] < famOIDs[j] })
	for _, oid := range famOIDs {
		c.cascadeDropOpFamily(oid)
	}
	if len(s.Collations) > 0 {
		type collEntry struct {
			name string
			oid  uint32
		}
		colls := make([]collEntry, 0, len(s.Collations))
		for name, uc := range s.Collations {
			if uc != nil {
				colls = append(colls, collEntry{name: name, oid: uc.OID})
			}
		}
		sort.Slice(colls, func(i, j int) bool { return colls[i].oid < colls[j].oid })
		for _, ce := range colls {
			if uc := s.Collations[ce.name]; uc != nil && uc.OID == ce.oid {
				c.cascadeDropCollation(s, ce.name, ce.oid)
			}
		}
	}
}

func (c *Catalog) deleteSchemaRecord(s *Schema) {
	if s == nil {
		return
	}
	delete(c.extensionMembers, s.OID)
	c.removeGrantsForObject('n', s.OID)
	c.removeComments('n', s.OID)
	c.removeDepsOf('n', s.OID)
	c.removeDepsOn('n', s.OID)
	delete(c.schemas, s.OID)
	if c.schemaByName[s.Name] == s {
		delete(c.schemaByName, s.Name)
	}
}
