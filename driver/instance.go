package driver

import "github.com/bytebase/omni/metadata"

// Instance is what Syncer.SyncInstance reads from a server.
type Instance struct {
	// Version is the server version, such as 17.12.0.
	Version string
	// Databases are the server's user databases, each with its name and the
	// attributes the engine reports for a database (pg: character set,
	// collation, owner) and no schemas.
	Databases []*metadata.DatabaseSchemaMetadata
	// Roles are the server's user roles.
	Roles []*Role
	// LowerCaseTableNames is MySQL's lower_case_table_names and its forks':
	// 0, 1 or 2. Other engines leave it 0.
	LowerCaseTableNames int32
}

// Role is one server role or user.
type Role struct {
	Name string
	// ConnectionLimit is the role's connection limit, nil when the engine
	// reports none.
	ConnectionLimit *int32
	// ValidUntil is when the role's password expires, as the engine writes
	// it, nil when it does not.
	ValidUntil *string
	// Attribute describes the role's privileges in the engine's terms
	// (pg: its role attributes, MySQL: its global GRANT statements), nil when
	// the engine reports none.
	Attribute *string
}
