// Package deparse renders Trino schema metadata back into SQL data-definition
// language (SDL). Its entry point, GetDatabaseDefinition, is the omni-side
// implementation of the bytebase schema.GetDatabaseDefinition contract for the
// Trino engine: given a snapshot of a database's schemas, tables, views and
// columns it produces the DDL that recreates them.
//
// # Why a deparse package, not the catalog
//
// trino/catalog also models the catalog -> schema -> table -> column hierarchy,
// but it is built for name resolution and completion: it folds identifiers to a
// normalized form and returns names in sorted order. SDL generation must do the
// opposite — it must preserve the metadata snapshot's original table/column
// order and emit names verbatim (the bytebase source quotes column.Name and
// table.Name exactly as given, never normalized). So deparse takes its own
// ordered, verbatim metadata model rather than reading from the catalog.
//
// # The metadata model
//
// The types below mirror, field-for-field, the subset of bytebase's
// storepb.DatabaseSchemaMetadata that the Trino schema dump consumes
// (DatabaseSchemaMetadata.Schemas, SchemaMetadata.{Name,Tables,Views},
// TableMetadata.{Name,Columns}, ViewMetadata.{Name,Definition},
// ColumnMetadata.{Name,Type,Nullable}). Keeping an omni-native copy lets the
// deparser live in the omni module without importing bytebase's generated
// protobuf; the bytebase-switch node maps the storepb structs onto these
// one-to-one when it wires GetDatabaseDefinition into plugin/schema/trino.
package deparse

// DatabaseSchemaMetadata is a snapshot of one Trino database (catalog) used as
// the input to GetDatabaseDefinition. Only the fields the Trino SDL dump reads
// are modeled; the storepb message carries more (collation, owner, extensions,
// …) that the Trino dump ignores.
type DatabaseSchemaMetadata struct {
	// Name is the catalog name. The Trino SDL dump does not emit it (CREATE
	// SCHEMA, CREATE TABLE and CREATE VIEW names are schema-relative), but it
	// is kept to mirror the storepb shape and identify the snapshot.
	Name string
	// Schemas are the database's schemas, in snapshot order.
	Schemas []*SchemaMetadata
}

// SchemaMetadata is one schema (the second level Trino exposes as the leading
// part of a CREATE TABLE name).
type SchemaMetadata struct {
	// Name is the schema name, emitted double-quoted in CREATE SCHEMA and as
	// the first qualifier of every CREATE TABLE and CREATE VIEW.
	Name string
	// Functions are the schema's user-defined functions, in snapshot order;
	// each becomes one CREATE FUNCTION statement before tables and views.
	Functions []*FunctionMetadata
	// Tables are the schema's tables, in snapshot order; each becomes one
	// CREATE TABLE statement.
	Tables []*TableMetadata
	// Views are the schema's views, in snapshot order; each becomes one
	// CREATE OR REPLACE VIEW statement after every schema's tables.
	Views []*ViewMetadata
	// MaterializedViews are the schema's materialized views, in snapshot order;
	// each becomes one CREATE MATERIALIZED VIEW statement in dependency order.
	MaterializedViews []*ViewMetadata
}

// FunctionMetadata is one user-defined function; Definition holds its CREATE FUNCTION DDL.
type FunctionMetadata struct {
	Name       string
	Definition string
}

// ViewMetadata is one view; it renders to a single CREATE OR REPLACE VIEW
// (or CREATE MATERIALIZED VIEW IF NOT EXISTS) statement.
type ViewMetadata struct {
	// Name is the view name, emitted double-quoted.
	Name string
	// Definition is the SELECT text the view is created AS, emitted as-is.
	Definition string
	// Comment is the optional view/materialized-view comment, emitted as COMMENT '...'.
	Comment string
}

// TableMetadata is one table; it renders to a single CREATE TABLE statement.
type TableMetadata struct {
	// Name is the table name, emitted double-quoted.
	Name string
	// Columns are the table's columns, in ordinal order; the order is preserved
	// in the emitted column list.
	Columns []*ColumnMetadata
	// Comment is the optional table comment, emitted as COMMENT '...'.
	Comment string
}

// ColumnMetadata is one column of a table.
type ColumnMetadata struct {
	// Name is the column name, emitted double-quoted.
	Name string
	// Type is the Trino type rendered as text (e.g. "bigint", "varchar",
	// "array(integer)"). It is emitted as-is and not re-parsed.
	Type string
	// Nullable reports whether the column accepts NULL. A non-nullable column
	// gets a trailing NOT NULL.
	Nullable bool
	// Default is the optional column default expression/literal.
	Default string
	// Comment is the optional column comment, emitted as COMMENT '...'.
	Comment string
}
