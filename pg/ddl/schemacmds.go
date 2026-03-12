package ddl

import (
	"fmt"
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
		ownerName = roleSpecName(stmt.Authrole)
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
	oid := c.oidGen.Next()
	s := &Schema{
		OID:       oid,
		Name:      schemaName,
		Owner:     ownerName,
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
		// unqualified DDL targets land in the newly created schema.
		prevSearchPath := c.searchPath
		c.searchPath = append([]string{schemaName}, prevSearchPath...)
		defer func() { c.searchPath = prevSearchPath }()

		for _, elt := range stmt.SchemaElts.Items {
			if elt == nil {
				continue
			}
			if err := c.ProcessUtility(elt); err != nil {
				return err
			}
		}
	}

	return nil
}

// RemoveSchemas drops one or more schemas from the catalog.
// Called from RemoveObjects when the drop target is OBJECT_SCHEMA.
//
// pg: src/backend/commands/schemacmds.c — (drop case in RemoveObjects via dropcmds.c)
func (c *Catalog) RemoveSchemas(stmt *nodes.DropStmt) error {
	cascade := stmt.Behavior == int(nodes.DROP_CASCADE)

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

		if !cascade && c.schemaHasObjects(s) {
			return errSchemaNotEmpty(schemaName)
		}

		// CASCADE: drop all contained objects.
		if cascade {
			for name, rel := range s.Relations {
				c.removeRelation(s, name, rel)
			}
			for name, idx := range s.Indexes {
				c.removeIndex(s, name, idx)
			}
			for _, seq := range s.Sequences {
				c.removeSequence(s, seq)
			}
			c.removeFunctionsInSchema(s)
			c.removeTypesInSchema(s)
		}

		c.removeComments('n', s.OID)
		delete(c.schemas, s.OID)
		delete(c.schemaByName, s.Name)
	}
	return nil
}
