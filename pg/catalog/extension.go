package catalog

import (
	"fmt"

	nodes "github.com/bytebase/omni/pg/ast"
	pgparser "github.com/bytebase/omni/pg/parser"
)

// Extension represents a registered extension.
//
// pg: src/backend/commands/extension.c — ExtensionControlFile
type Extension struct {
	OID         uint32
	Name        string
	SchemaOID   uint32
	Relocatable bool
	// Version is the VERSION the extension was created with, or "" when the
	// statement named none.
	Version string
}

// CreateExtension loads and executes a bundled extension DDL script.
//
// pg: src/backend/commands/extension.c — CreateExtension
func (c *Catalog) CreateExtension(stmt *nodes.CreateExtensionStmt) error {
	// IF NOT EXISTS check.
	if c.extByName[stmt.Extname] != nil {
		if stmt.IfNotExists {
			c.addWarning(CodeWarningSkip,
				fmt.Sprintf("extension %q already exists, skipping", stmt.Extname))
			return nil
		}
		return errDuplicateObject("extension", stmt.Extname)
	}

	// Parse Options for "schema" DefElem.
	// pg: src/backend/commands/extension.c — CreateExtension (option parsing)
	var schemaName, version string
	if stmt.Options != nil {
		for _, item := range stmt.Options.Items {
			de, ok := item.(*nodes.DefElem)
			if !ok {
				continue
			}
			sv, ok := de.Arg.(*nodes.String)
			if !ok {
				continue
			}
			switch de.Defname {
			case "schema":
				schemaName = sv.Str
			case "new_version":
				version = sv.Str
			}
		}
	}

	// Resolve target schema.
	// Default: first in search path (typically "public").
	targetSchema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
	}

	// An extension without a bundled script is registered with no objects.
	sql, ok := extensionScripts[stmt.Extname]
	if !ok {
		c.addWarning(CodeWarning,
			fmt.Sprintf("extension %q is not bundled; its objects are unknown", stmt.Extname))
	}

	// Parse the extension DDL script.
	list, err := pgparser.Parse(sql)
	if err != nil {
		return fmt.Errorf("extension %q: parse error: %w", stmt.Extname, err)
	}
	firstMember, firstGrant := c.oidGen.next, len(c.grants)

	// Temporarily prepend target schema to search path so unqualified types
	// and functions land in the extension's schema.
	// pg: src/backend/commands/extension.c — execute_extension_script (line 973)
	origSearchPath := c.searchPath
	if targetSchema.Name != "" {
		found := false
		for _, s := range c.searchPath {
			if s == targetSchema.Name {
				found = true
				break
			}
		}
		if !found {
			c.searchPath = append([]string{targetSchema.Name}, c.searchPath...)
		}
	}

	// Execute each statement through ProcessUtility.
	if list != nil {
		for _, item := range list.Items {
			var node nodes.Node
			if raw, ok := item.(*nodes.RawStmt); ok {
				node = raw.Stmt.(nodes.Node)
			} else {
				node = item.(nodes.Node)
			}
			if err := c.ProcessUtility(node); err != nil {
				c.searchPath = origSearchPath
				return fmt.Errorf("extension %q: %w", stmt.Extname, err)
			}
		}
	}

	// Restore original search path.
	c.searchPath = origSearchPath
	c.markExtensionMembers(firstMember, c.oidGen.next)
	for i := firstGrant; i < len(c.grants); i++ {
		c.grants[i].FromExtension = true
	}

	// Register extension metadata.
	ext := &Extension{
		OID:         c.oidGen.Next(),
		Name:        stmt.Extname,
		SchemaOID:   targetSchema.OID,
		Relocatable: true,
		Version:     version,
	}
	c.extensions[ext.OID] = ext
	c.extByName[ext.Name] = ext

	return nil
}

// markExtensionMembers records the objects whose OIDs fall in [first, end)
// as belonging to an extension, as pg_depend's 'e' entries do. Every object
// an extension script creates takes its OID while the script runs.
func (c *Catalog) markExtensionMembers(first, end uint32) {
	for oid := first; oid < end; oid++ {
		c.extensionMembers[oid] = true
	}
}

// isExtensionMember reports whether the object with this OID belongs to an
// extension. Such objects come and go with their extension, so Diff leaves
// them out.
func (c *Catalog) isExtensionMember(oid uint32) bool {
	return c.extensionMembers[oid]
}

// RegisterExtensionSQL allows external code to register extension DDL scripts.
// Once registered, CREATE EXTENSION <name> will execute the provided SQL.
func RegisterExtensionSQL(name, sql string) {
	extensionScripts[name] = sql
}
