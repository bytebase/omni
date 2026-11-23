package catalog

import (
	"context"
	"fmt"
	"math/big"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle"
	"github.com/bytebase/omni/oracle/ast"
)

// ServerConfig models Oracle database instance and container settings that
// gate feature-specific DDL validation in Exec.
//
// The zero value keeps LoadSQL and GetDatabaseMetadata permissive for schema
// dumps from extended-collation, sharded, application-common, TDE-enabled, or
// cross-schema databases, while StandardOracle23ServerConfig matches a default
// Oracle 23ai Free PDB.
type ServerConfig struct {
	// MaxStringSize models Oracle's MAX_STRING_SIZE parameter. When set to
	// "STANDARD", column COLLATE and table/view DEFAULT COLLATION clauses are
	// refused with ORA-43929 before table/column/view existence checks.
	MaxStringSize string

	// DisableShardDDL models a non-sharded Oracle database (SHARD DDL disabled,
	// toggled via ALTER SESSION ENABLE/DISABLE SHARD DDL). When true, sequence
	// SHARD/NOSHARD clauses fail with ORA-02511.
	DisableShardDDL bool

	// NonApplicationRoot models a standard PDB outside an Oracle Application
	// Root container. When true, SHARING = METADATA | DATA | EXTENDED DATA
	// clauses fail with ORA-65021.
	NonApplicationRoot bool

	// KeystoreClosed models a database whose TDE keystore/wallet is not open.
	// When true, column ENCRYPT clauses fail with ORA-28365 after
	// table/column/foreign-key validation.
	KeystoreClosed bool

	// ValidateExternalSchemaRefs, when true, refuses foreign-key REFERENCES
	// targeting an external schema (other than the catalog owner) with
	// ORA-00942 instead of preserving unvalidated cross-schema foreign keys.
	ValidateExternalSchemaRefs bool

	// ApplySetUnused, when true, executes ALTER TABLE ... SET UNUSED COLUMN by
	// marking the columns unused (removing them and their dependent indexes
	// and constraints from visible table metadata). When false, LoadSQL keeps
	// its legacy refusal of SET UNUSED.
	ApplySetUnused bool
}

// StandardOracle23ServerConfig returns the instance configuration of a default
// Oracle 23ai Free pluggable database (MAX_STRING_SIZE=STANDARD, non-sharded,
// non-application-root, TDE keystore closed, isolated schema without external
// grants, and live SET UNUSED execution).
func StandardOracle23ServerConfig() ServerConfig {
	return ServerConfig{
		MaxStringSize:              "STANDARD",
		DisableShardDDL:            true,
		NonApplicationRoot:         true,
		KeystoreClosed:             true,
		ValidateExternalSchemaRefs: true,
		ApplySetUnused:             true,
	}
}

// Catalog is the in-memory Oracle catalog. Oracle has one schema per user, so a
// snapshot is normally one schema; the metadata type still carries the list.
type Catalog struct {
	meta                 *metadatapb.DatabaseSchemaMetadata
	serverConfig         ServerConfig
	namedNotNull         map[string]bool
	primaryKeyNotNull    map[string]bool
	syntheticConstraints map[string]bool
	viewConstraints      map[string][]viewConstraint
	editioningViews      map[string]bool
	encryptedColumns     map[string]bool
	domains              map[string]bool
	editions             map[string]bool
	sequenceCurrent      map[string]*big.Int
}

// New returns an empty catalog.
func New() *Catalog {
	return &Catalog{
		meta:                 &metadatapb.DatabaseSchemaMetadata{},
		namedNotNull:         map[string]bool{},
		primaryKeyNotNull:    map[string]bool{},
		syntheticConstraints: map[string]bool{},
		viewConstraints:      map[string][]viewConstraint{},
		editioningViews:      map[string]bool{},
		encryptedColumns:     map[string]bool{},
		domains:              map[string]bool{},
		editions:             map[string]bool{"ORA$BASE": true},
		sequenceCurrent:      map[string]*big.Int{},
	}
}

// ServerConfig returns the catalog's active Oracle server configuration.
func (c *Catalog) ServerConfig() ServerConfig {
	if c == nil {
		return ServerConfig{}
	}
	return c.serverConfig
}

// SetServerConfig replaces the catalog's Oracle server configuration.
func (c *Catalog) SetServerConfig(cfg ServerConfig) {
	if c == nil {
		return
	}
	c.serverConfig = cfg
}

// Metadata returns the catalog's snapshot.
func (c *Catalog) Metadata() *metadatapb.DatabaseSchemaMetadata {
	if c == nil || c.meta == nil {
		return &metadatapb.DatabaseSchemaMetadata{}
	}
	return c.meta
}

// SetMetadata replaces the catalog's snapshot.
func (c *Catalog) SetMetadata(meta *metadatapb.DatabaseSchemaMetadata) {
	if meta == nil {
		meta = &metadatapb.DatabaseSchemaMetadata{}
	}
	c.meta = meta
	c.namedNotNull = map[string]bool{}
	c.primaryKeyNotNull = map[string]bool{}
	c.syntheticConstraints = map[string]bool{}
	c.viewConstraints = map[string][]viewConstraint{}
	c.editioningViews = map[string]bool{}
	c.encryptedColumns = map[string]bool{}
	c.domains = map[string]bool{}
	c.editions = map[string]bool{"ORA$BASE": true}
	c.sequenceCurrent = map[string]*big.Int{}
}

// LoadMetadataReport lists snapshot objects that could not be installed as
// defined. The current catalog stores the snapshot rather than re-executing it,
// so both maps are empty; they exist for the driver contract and for callers
// that inspect a report.
type LoadMetadataReport struct {
	Degraded map[string]error
	Missing  map[string]error
}

// LoadMetadata installs a synced snapshot into the catalog.
func (c *Catalog) LoadMetadata(ctx context.Context, meta *metadatapb.DatabaseSchemaMetadata) (*LoadMetadataReport, error) {
	report := &LoadMetadataReport{Degraded: map[string]error{}, Missing: map[string]error{}}
	if err := ctx.Err(); err != nil {
		return report, err
	}
	c.namedNotNull = map[string]bool{}
	c.primaryKeyNotNull = map[string]bool{}
	c.syntheticConstraints = map[string]bool{}
	c.viewConstraints = map[string][]viewConstraint{}
	c.editioningViews = map[string]bool{}
	c.encryptedColumns = map[string]bool{}
	c.domains = map[string]bool{}
	c.editions = map[string]bool{"ORA$BASE": true}
	c.sequenceCurrent = map[string]*big.Int{}
	if meta == nil {
		c.meta = &metadatapb.DatabaseSchemaMetadata{}
		return report, nil
	}
	c.meta = cloneMetadata(meta)
	return report, nil
}

// ExecOptions controls Exec.
type ExecOptions struct {
	ContinueOnError bool
}

// ExecResult is the outcome of one statement.
type ExecResult struct {
	Index   int
	SQL     string
	Line    int
	Skipped bool
	Error   error
}

// Exec parses and applies one or more DDL statements. DML and transaction
// statements are skipped.
func (c *Catalog) Exec(sql string, opts *ExecOptions) ([]ExecResult, error) {
	stmts, err := oracle.Parse(sql)
	if err != nil {
		return nil, err
	}
	continueOnError := opts != nil && opts.ContinueOnError
	results := make([]ExecResult, 0, len(stmts))
	for i, stmt := range stmts {
		result := ExecResult{Index: i, SQL: stmt.Text, Line: stmt.Start.Line}
		node, ok := stmt.AST.(ast.StmtNode)
		if !ok || node == nil {
			result.Skipped = true
			results = append(results, result)
			continue
		}
		if isDML(node) {
			result.Skipped = true
			results = append(results, result)
			continue
		}
		result.Error = c.applyStatement(node, sql)
		results = append(results, result)
		if result.Error != nil && !continueOnError {
			break
		}
	}
	return results, nil
}

// LoadSQL returns a catalog built from sql.
func LoadSQL(sql string) (*Catalog, error) {
	c := New()
	results, err := c.Exec(sql, nil)
	if err != nil {
		return nil, err
	}
	for _, result := range results {
		if result.Error != nil {
			return c, result.Error
		}
	}
	return c, nil
}

// LoadSDL is LoadSQL: Oracle's declarative form is its DDL.
func LoadSDL(sql string) (*Catalog, error) {
	return LoadSQL(sql)
}

func isDML(stmt ast.StmtNode) bool {
	switch stmt.(type) {
	case *ast.SelectStmt, *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt,
		*ast.CommitStmt, *ast.RollbackStmt, *ast.SavepointStmt, *ast.SetTransactionStmt,
		*ast.ExplainPlanStmt, *ast.LockTableStmt, *ast.CallStmt:
		return true
	default:
		return false
	}
}

func (c *Catalog) applyStatement(stmt ast.StmtNode, source string) error {
	extractor := newMetadataExtractorFromMetadata(c.meta, source)
	extractor.serverConfig = c.serverConfig
	for key := range c.namedNotNull {
		extractor.namedNotNull[key] = true
	}
	for key := range c.primaryKeyNotNull {
		extractor.primaryKeyNotNull[key] = true
	}
	for key := range c.syntheticConstraints {
		extractor.syntheticConstraints[key] = true
	}
	for viewName, constraints := range c.viewConstraints {
		extractor.viewConstraints[viewName] = cloneViewConstraints(constraints)
	}
	for viewName := range c.editioningViews {
		extractor.editioningViews[viewName] = true
	}
	for key := range c.encryptedColumns {
		extractor.encryptedColumns[key] = true
	}
	for dom := range c.domains {
		extractor.domains[dom] = true
	}
	if len(c.editions) > 0 {
		extractor.editions = make(map[string]bool, len(c.editions))
		for ed := range c.editions {
			extractor.editions[ed] = true
		}
	}
	for seqName, val := range c.sequenceCurrent {
		if val != nil {
			extractor.sequenceCurrent[seqName] = new(big.Int).Set(val)
		}
	}
	extractor.extractStatement(stmt)
	if extractor.err != nil {
		return extractor.err
	}
	meta := extractor.databaseMetadata()
	// Keep any additional schemas the catalog already carried; Oracle's DDL
	// names one schema at a time.
	result := cloneMetadata(c.meta)
	result.Name = meta.GetName()
	if len(meta.GetSchemas()) > 0 {
		if len(result.GetSchemas()) == 0 {
			result.Schemas = meta.Schemas
		} else {
			oldSchema, newSchema := result.Schemas[0], meta.Schemas[0]
			merged := cloneSchema(oldSchema)
			merged.Name = oldSchema.GetName()
			if merged.Name == "" && extractor.explicitSchema {
				merged.Name = newSchema.GetName()
			}
			merged.Tables = newSchema.GetTables()
			merged.Views = newSchema.GetViews()
			merged.MaterializedViews = newSchema.GetMaterializedViews()
			merged.Functions = newSchema.GetFunctions()
			merged.Procedures = newSchema.GetProcedures()
			merged.Sequences = newSchema.GetSequences()
			merged.Packages = newSchema.GetPackages()
			merged.Synonyms = newSchema.GetSynonyms()
			result.Schemas[0] = merged
		}
	}
	result.PublicSynonyms = meta.GetPublicSynonyms()
	c.meta = result
	c.serverConfig = extractor.serverConfig
	c.namedNotNull = extractor.namedNotNull
	c.primaryKeyNotNull = extractor.primaryKeyNotNull
	c.syntheticConstraints = extractor.syntheticConstraints
	c.viewConstraints = extractor.viewConstraints
	c.editioningViews = extractor.editioningViews
	c.encryptedColumns = extractor.encryptedColumns
	c.domains = extractor.domains
	c.editions = extractor.editions
	c.sequenceCurrent = extractor.sequenceCurrent
	return nil
}

// Error is an Oracle catalog error.
type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	if e.Code == "" {
		return e.Message
	}
	return fmt.Sprintf("ORA-%s: %s", e.Code, e.Message)
}
