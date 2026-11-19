package ast

// This file holds the statement nodes the fork's tail added once the
// statements StarRocks kept from Doris were read: resource groups, pipes,
// data cache rules, the SQL and backend blacklists, prepared statements, plan
// advisor, TRANSLATE, EXECUTE AS and REFRESH EXTERNAL TABLE. Each cites the
// rule of StarRocks.g4 it comes from.

// ---------------------------------------------------------------------------
// REFRESH EXTERNAL TABLE
// ---------------------------------------------------------------------------

// RefreshExternalTableStmt represents refreshTableStatement (StarRocks.g4:550):
//
//	REFRESH EXTERNAL TABLE qualifiedName [(PARTITION ('p' [',' 'p']*))]
type RefreshExternalTableStmt struct {
	Name       *ObjectName
	Partitions []string
	Loc        Loc
}

// Tag implements Node.
func (n *RefreshExternalTableStmt) Tag() NodeTag { return T_RefreshExternalTableStmt }

var _ Node = (*RefreshExternalTableStmt)(nil)

// ExternalCatalogRef is one catalog of BACKUP's or RESTORE's optional scope
// (identifierList :2855, identifierWithAliasList :2858): a catalog, and the
// name RESTORE restores it as when AS names one.
type ExternalCatalogRef struct {
	Name  string
	Alias string
	Loc   Loc
}

// Tag implements Node.
func (n *ExternalCatalogRef) Tag() NodeTag { return T_ExternalCatalogRef }

var _ Node = (*ExternalCatalogRef)(nil)

// ---------------------------------------------------------------------------
// RESOURCE GROUP (:1358)
// ---------------------------------------------------------------------------

// CreateResourceGroupStmt represents createResourceGroupStatement
// (StarRocks.g4:1358):
//
//	CREATE RESOURCE GROUP [IF NOT EXISTS] [OR REPLACE] name
//	    [TO classifier [',' classifier]*] WITH '(' property [',' property]* ')'
//
// A classifier is a parenthesized expression list.
type CreateResourceGroupStmt struct {
	Name        *ObjectName
	IfNotExists bool
	OrReplace   bool
	Classifiers [][]Node
	Properties  []*Property
	Loc         Loc
}

// Tag implements Node.
func (n *CreateResourceGroupStmt) Tag() NodeTag { return T_CreateResourceGroupStmt }

var _ Node = (*CreateResourceGroupStmt)(nil)

// DropResourceGroupStmt represents dropResourceGroupStatement (StarRocks.g4:1363):
//
//	DROP RESOURCE GROUP identifier
type DropResourceGroupStmt struct {
	Name *ObjectName
	Loc  Loc
}

// Tag implements Node.
func (n *DropResourceGroupStmt) Tag() NodeTag { return T_DropResourceGroupStmt }

var _ Node = (*DropResourceGroupStmt)(nil)

// AlterResourceGroupStmt represents alterResourceGroupStatement
// (StarRocks.g4:1367):
//
//	ALTER RESOURCE GROUP name ADD classifier [',' classifier]*
//	ALTER RESOURCE GROUP name DROP '(' INTEGER_VALUE [',' INTEGER_VALUE]* ')'
//	ALTER RESOURCE GROUP name DROP ALL
//	ALTER RESOURCE GROUP name WITH '(' property [',' property]* ')'
type AlterResourceGroupStmt struct {
	Name        *ObjectName
	Action      string // "ADD" | "DROP" | "DROP ALL" | "WITH"
	Classifiers [][]Node
	ClassIDs    []string
	Properties  []*Property
	Loc         Loc
}

// Tag implements Node.
func (n *AlterResourceGroupStmt) Tag() NodeTag { return T_AlterResourceGroupStmt }

var _ Node = (*AlterResourceGroupStmt)(nil)

// ---------------------------------------------------------------------------
// DATA CACHE RULES (:1845)
// ---------------------------------------------------------------------------

// CreateDataCacheRuleStmt represents createDataCacheRuleStatement
// (StarRocks.g4:1851):
//
//	CREATE DATACACHE RULE target [(WHERE expression)] PRIORITY '=' ['-'] INTEGER_VALUE
//	    [properties]
//
// target is three parts — catalog, database and table — each an identifier, a
// string or the '*' wildcard, kept apart so that a quoted part holding a dot
// is not the same target as two parts.
type CreateDataCacheRuleStmt struct {
	Target     []string // catalog, database and table, each an identifier, a string or "*"
	Where      Node
	Priority   string // the integer as written
	Negative   bool
	Properties []*Property
	Loc        Loc
}

// Tag implements Node.
func (n *CreateDataCacheRuleStmt) Tag() NodeTag { return T_CreateDataCacheRuleStmt }

var _ Node = (*CreateDataCacheRuleStmt)(nil)

// DropDataCacheRuleStmt represents dropDataCacheRuleStatement (StarRocks.g4:1859):
//
//	DROP DATACACHE RULE INTEGER_VALUE
type DropDataCacheRuleStmt struct {
	ID  string
	Loc Loc
}

// Tag implements Node.
func (n *DropDataCacheRuleStmt) Tag() NodeTag { return T_DropDataCacheRuleStmt }

var _ Node = (*DropDataCacheRuleStmt)(nil)

// ClearDataCacheRulesStmt represents clearDataCacheRulesStatement
// (StarRocks.g4:1863): CLEAR DATACACHE RULES.
type ClearDataCacheRulesStmt struct {
	Loc Loc
}

// Tag implements Node.
func (n *ClearDataCacheRulesStmt) Tag() NodeTag { return T_ClearDataCacheRulesStmt }

var _ Node = (*ClearDataCacheRulesStmt)(nil)

// ---------------------------------------------------------------------------
// PIPE (:1913)
// ---------------------------------------------------------------------------

// CreatePipeStmt represents createPipeStatement (StarRocks.g4:1913):
//
//	CREATE [OR REPLACE] PIPE [IF NOT EXISTS] qualifiedName
//	    [properties] AS insertStatement
type CreatePipeStmt struct {
	Name        *ObjectName
	OrReplace   bool
	IfNotExists bool
	Properties  []*Property
	Insert      Node // the insertStatement body
	Loc         Loc
}

// Tag implements Node.
func (n *CreatePipeStmt) Tag() NodeTag { return T_CreatePipeStmt }

var _ Node = (*CreatePipeStmt)(nil)

// DropPipeStmt represents dropPipeStatement (StarRocks.g4:1919):
//
//	DROP PIPE [IF EXISTS] qualifiedName
type DropPipeStmt struct {
	Name     *ObjectName
	IfExists bool
	Loc      Loc
}

// Tag implements Node.
func (n *DropPipeStmt) Tag() NodeTag { return T_DropPipeStmt }

var _ Node = (*DropPipeStmt)(nil)

// AlterPipeStmt represents alterPipeStatement (StarRocks.g4:1923):
//
//	ALTER PIPE qualifiedName alterPipeClause
//
// where alterPipeClause (:1927) is SUSPEND, RESUME, RETRY ALL, RETRY FILE string
// or SET propertyList.
type AlterPipeStmt struct {
	Name       *ObjectName
	Action     string // "SUSPEND" | "RESUME" | "RETRY ALL" | "RETRY FILE" | "SET"
	File       string // RETRY FILE 'name'
	Properties []*Property
	Loc        Loc
}

// Tag implements Node.
func (n *AlterPipeStmt) Tag() NodeTag { return T_AlterPipeStmt }

var _ Node = (*AlterPipeStmt)(nil)

// DescPipeStmt represents descPipeStatement (StarRocks.g4:1935):
//
//	(DESC | DESCRIBE) PIPE qualifiedName
type DescPipeStmt struct {
	Name *ObjectName
	Loc  Loc
}

// Tag implements Node.
func (n *DescPipeStmt) Tag() NodeTag { return T_DescPipeStmt }

var _ Node = (*DescPipeStmt)(nil)

// ---------------------------------------------------------------------------
// PLAN ADVISOR (:2024)
// ---------------------------------------------------------------------------

// AlterPlanAdvisorStmt represents alterPlanAdvisorAddStatement
// (StarRocks.g4:2024) and alterPlanAdvisorDropStatement (:2030):
//
//	ALTER PLAN ADVISOR ADD queryStatement
//	ALTER PLAN ADVISOR DROP string
type AlterPlanAdvisorStmt struct {
	Action string // "ADD" | "DROP"
	Query  Node   // ADD
	ID     string // DROP's baseline id
	Loc    Loc
}

// Tag implements Node.
func (n *AlterPlanAdvisorStmt) Tag() NodeTag { return T_AlterPlanAdvisorStmt }

var _ Node = (*AlterPlanAdvisorStmt)(nil)

// TruncatePlanAdvisorStmt represents truncatePlanAdvisorStatement
// (StarRocks.g4:2027): TRUNCATE PLAN ADVISOR.
type TruncatePlanAdvisorStmt struct {
	Loc Loc
}

// Tag implements Node.
func (n *TruncatePlanAdvisorStmt) Tag() NodeTag { return T_TruncatePlanAdvisorStmt }

var _ Node = (*TruncatePlanAdvisorStmt)(nil)

// ---------------------------------------------------------------------------
// TRANSLATE (:2073)
// ---------------------------------------------------------------------------

// TranslateStmt represents translateStatement (StarRocks.g4:2073):
//
//	TRANSLATE dialect translateSQL
//
// translateSQL (:2077) is `.+`: the rest of the statement, which StarRocks
// hands to the SQL converter as text rather than reading as SQL.
type TranslateStmt struct {
	Dialect string
	SQL     string
	Loc     Loc
}

// Tag implements Node.
func (n *TranslateStmt) Tag() NodeTag { return T_TranslateStmt }

var _ Node = (*TranslateStmt)(nil)

// ---------------------------------------------------------------------------
// PREPARE / EXECUTE / DEALLOCATE (:2282)
// ---------------------------------------------------------------------------

// PrepareStmt represents prepareStatement (StarRocks.g4:2282):
//
//	PREPARE identifier FROM prepareSql
//
// prepareSql (:2300) is a statement or a string literal. Query holds the
// parsed statement, which for a string literal is the tree its text makes: the
// string is SQL too, and its escapes are re-encoded so the tree's Locs are the
// statement's. SQL holds the string literal's text.
type PrepareStmt struct {
	Name  *ObjectName
	Query Node
	SQL   string
	Loc   Loc
}

// Tag implements Node.
func (n *PrepareStmt) Tag() NodeTag { return T_PrepareStmt }

var _ Node = (*PrepareStmt)(nil)

// ExecuteStmt represents executeStatement (StarRocks.g4:2291):
//
//	EXECUTE identifier [USING '@'identifierOrString [',' '@'identifierOrString]*]
//
// ExecuteAsStmt is the other statement EXECUTE opens, told apart by AS.
type ExecuteStmt struct {
	Name  *ObjectName
	Using []string
	Loc   Loc
}

// Tag implements Node.
func (n *ExecuteStmt) Tag() NodeTag { return T_ExecuteStmt }

var _ Node = (*ExecuteStmt)(nil)

// DeallocateStmt represents deallocateStatement (StarRocks.g4:2295):
//
//	(DEALLOCATE | DROP) PREPARE identifier
type DeallocateStmt struct {
	Name *ObjectName
	Loc  Loc
}

// Tag implements Node.
func (n *DeallocateStmt) Tag() NodeTag { return T_DeallocateStmt }

var _ Node = (*DeallocateStmt)(nil)

// ExecuteAsStmt represents executeAsStatement (StarRocks.g4:1646):
//
//	EXECUTE AS user [WITH NO REVERT]
type ExecuteAsStmt struct {
	User         *UserIdentity
	WithNoRevert bool
	Loc          Loc
}

// Tag implements Node.
func (n *ExecuteAsStmt) Tag() NodeTag { return T_ExecuteAsStmt }

var _ Node = (*ExecuteAsStmt)(nil)

// ---------------------------------------------------------------------------
// BLACKLISTS (:1815, :1833)
// ---------------------------------------------------------------------------

// BlacklistStmt represents the four blacklist statements (StarRocks.g4:1815,
// :1819, :1833, :1837):
//
//	ADD SQLBLACKLIST string
//	DELETE SQLBLACKLIST INTEGER_VALUE [',' INTEGER_VALUE]*
//	ADD BACKEND BLACKLIST INTEGER_VALUE [',' INTEGER_VALUE]*
//	DELETE BACKEND BLACKLIST INTEGER_VALUE [',' INTEGER_VALUE]*
type BlacklistStmt struct {
	Add     bool
	Backend bool
	SQL     string
	IDs     []string
	Loc     Loc
}

// Tag implements Node.
func (n *BlacklistStmt) Tag() NodeTag { return T_BlacklistStmt }

var _ Node = (*BlacklistStmt)(nil)
