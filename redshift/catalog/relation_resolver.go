package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/redshift/ast"
	pgparser "github.com/bytebase/omni/redshift/parser"
)

const (
	RelationKindTable            byte = 'r'
	RelationKindView             byte = 'v'
	RelationKindMaterializedView byte = 'm'
)

// RelationResolver lazily supplies relation metadata for catalog misses.
type RelationResolver interface {
	ResolveRelation(schemaName, relationName string, searchPath []string) (*RelationSpec, error)
}

// RelationSpec describes a relation returned by a RelationResolver.
type RelationSpec struct {
	SchemaName string
	Name       string
	Kind       byte
	Columns    []RelationColumnSpec

	// Definition is the SELECT definition, or a full CREATE VIEW /
	// CREATE MATERIALIZED VIEW statement, for a view or materialized view.
	Definition string
}

// RelationColumnSpec describes one relation column in a resolver result.
type RelationColumnSpec struct {
	Name string
	Type string
}

func (c *Catalog) resolveMissingRelation(schemaName, relName string) error {
	if c.relationResolver == nil {
		return nil
	}

	requestKey := relationResolutionKey(schemaName, relName)
	if c.relationResolutionStack[requestKey] {
		return errCyclicRelationResolution(requestKey)
	}
	c.relationResolutionStack[requestKey] = true
	defer delete(c.relationResolutionStack, requestKey)

	searchPath := append([]string(nil), c.searchPath...)
	spec, err := c.relationResolver.ResolveRelation(schemaName, relName, searchPath)
	if err != nil {
		return err
	}
	if spec == nil {
		return nil
	}

	spec = normalizeRelationSpec(spec, schemaName, relName)
	specKey := relationResolutionKey(spec.SchemaName, spec.Name)
	if specKey != requestKey {
		if c.relationResolutionStack[specKey] {
			return errCyclicRelationResolution(specKey)
		}
		c.relationResolutionStack[specKey] = true
		defer delete(c.relationResolutionStack, specKey)
	}

	return c.materializeRelationSpec(spec)
}

func normalizeRelationSpec(spec *RelationSpec, schemaName, relName string) *RelationSpec {
	normalized := *spec
	if normalized.SchemaName == "" {
		normalized.SchemaName = schemaName
	}
	if normalized.SchemaName == "" {
		normalized.SchemaName = "public"
	}
	if normalized.Name == "" {
		normalized.Name = relName
	}
	return &normalized
}

func relationResolutionKey(schemaName, relName string) string {
	return schemaName + "." + relName
}

func errCyclicRelationResolution(key string) error {
	return &Error{
		Code:    CodeInvalidObjectDefinition,
		Message: fmt.Sprintf("cyclic relation resolution detected for %q", key),
	}
}

func (c *Catalog) materializeRelationSpec(spec *RelationSpec) error {
	if spec.SchemaName == "" || spec.Name == "" {
		return &Error{Code: CodeInvalidObjectDefinition, Message: "relation resolver returned an incomplete relation name"}
	}
	if err := c.ensureResolverSchema(spec.SchemaName); err != nil {
		return err
	}

	switch spec.Kind {
	case RelationKindTable:
		return c.materializeResolvedTable(spec)
	case RelationKindView:
		stmt, err := parseRelationDefinitionView(spec)
		if err != nil {
			return err
		}
		return c.DefineView(stmt)
	case RelationKindMaterializedView:
		stmt, err := parseRelationDefinitionMaterializedView(spec)
		if err != nil {
			return err
		}
		return c.ExecCreateTableAs(stmt)
	default:
		return &Error{
			Code:    CodeFeatureNotSupported,
			Message: fmt.Sprintf("relation resolver returned unsupported relation kind %q", spec.Kind),
		}
	}
}

func (c *Catalog) ensureResolverSchema(schemaName string) error {
	if c.schemaByName[schemaName] != nil {
		return nil
	}
	return c.CreateSchemaCommand(&nodes.CreateSchemaStmt{Schemaname: schemaName})
}

func (c *Catalog) materializeResolvedTable(spec *RelationSpec) error {
	if len(spec.Columns) == 0 {
		return &Error{Code: CodeInvalidTableDefinition, Message: "relation resolver table spec requires at least one column"}
	}

	parts := make([]string, 0, len(spec.Columns))
	for _, col := range spec.Columns {
		if col.Name == "" || strings.TrimSpace(col.Type) == "" {
			return &Error{Code: CodeInvalidColumnDefinition, Message: "relation resolver table column requires name and type"}
		}
		parts = append(parts, fmt.Sprintf("%s %s", quoteIdentAlways(col.Name), col.Type))
	}
	ddl := fmt.Sprintf("CREATE TABLE %s.%s (%s)", quoteIdentAlways(spec.SchemaName), quoteIdentAlways(spec.Name), strings.Join(parts, ", "))

	list, err := pgparser.Parse(ddl)
	if err != nil {
		return err
	}
	if list == nil || len(list.Items) != 1 {
		return &Error{Code: CodeInvalidTableDefinition, Message: "relation resolver table definition did not parse to one statement"}
	}
	raw, ok := list.Items[0].(*nodes.RawStmt)
	if !ok {
		return &Error{Code: CodeInvalidTableDefinition, Message: "relation resolver table definition did not parse to a raw statement"}
	}
	stmt, ok := raw.Stmt.(*nodes.CreateStmt)
	if !ok {
		return &Error{Code: CodeInvalidTableDefinition, Message: "relation resolver table definition did not parse to CREATE TABLE"}
	}
	return c.DefineRelation(stmt, RelationKindTable)
}

func parseRelationDefinitionView(spec *RelationSpec) (*nodes.ViewStmt, error) {
	stmt, err := parseRelationDefinitionStatement(spec)
	if err != nil {
		return nil, err
	}
	switch n := stmt.(type) {
	case *nodes.SelectStmt:
		return &nodes.ViewStmt{
			View:  &nodes.RangeVar{Schemaname: spec.SchemaName, Relname: spec.Name},
			Query: n,
		}, nil
	case *nodes.ViewStmt:
		copy := *n
		copy.View = &nodes.RangeVar{Schemaname: spec.SchemaName, Relname: spec.Name}
		return &copy, nil
	default:
		return nil, &Error{Code: CodeInvalidObjectDefinition, Message: "relation resolver view requires a SELECT definition"}
	}
}

func parseRelationDefinitionMaterializedView(spec *RelationSpec) (*nodes.CreateTableAsStmt, error) {
	stmt, err := parseRelationDefinitionStatement(spec)
	if err != nil {
		return nil, err
	}
	switch n := stmt.(type) {
	case *nodes.SelectStmt:
		return &nodes.CreateTableAsStmt{
			Query:   n,
			Into:    &nodes.IntoClause{Rel: &nodes.RangeVar{Schemaname: spec.SchemaName, Relname: spec.Name}},
			Objtype: nodes.OBJECT_MATVIEW,
		}, nil
	case *nodes.CreateTableAsStmt:
		if n.Objtype != nodes.OBJECT_MATVIEW {
			return nil, &Error{Code: CodeInvalidObjectDefinition, Message: "relation resolver view requires a SELECT definition"}
		}
		copy := *n
		if copy.Into == nil {
			copy.Into = &nodes.IntoClause{}
		} else {
			into := *copy.Into
			copy.Into = &into
		}
		copy.Into.Rel = &nodes.RangeVar{Schemaname: spec.SchemaName, Relname: spec.Name}
		copy.Objtype = nodes.OBJECT_MATVIEW
		return &copy, nil
	default:
		return nil, &Error{Code: CodeInvalidObjectDefinition, Message: "relation resolver view requires a SELECT definition"}
	}
}

func parseRelationDefinitionStatement(spec *RelationSpec) (nodes.Node, error) {
	definition := strings.TrimSpace(spec.Definition)
	if definition == "" {
		return nil, &Error{Code: CodeInvalidObjectDefinition, Message: "relation resolver view requires a SELECT definition"}
	}
	list, err := pgparser.Parse(definition)
	if err != nil {
		return nil, err
	}
	if list == nil || len(list.Items) != 1 {
		return nil, &Error{Code: CodeInvalidObjectDefinition, Message: "relation resolver view requires a SELECT definition"}
	}
	raw, ok := list.Items[0].(*nodes.RawStmt)
	if !ok {
		return nil, &Error{Code: CodeInvalidObjectDefinition, Message: "relation resolver view requires a SELECT definition"}
	}
	return raw.Stmt, nil
}
