package ast

// This file holds AST node types for Doris utility / admin commands:
//   BACKUP, RESTORE, KILL, LOCK/UNLOCK TABLES, INSTALL/UNINSTALL PLUGIN,
//   WARM UP, CLEAN, CANCEL (generic), and RECOVER.
//
// All nodes carry a Loc field and implement Node via Tag().

// ---------------------------------------------------------------------------
// BACKUP / RESTORE SNAPSHOT
// ---------------------------------------------------------------------------

// BackupStmt represents:
//
//	BACKUP SNAPSHOT [db.]label TO repo_name
//	    [ON (tbl [PARTITION(p1,...)], ...)]
//	    [PROPERTIES("key"="value", ...)]
type BackupStmt struct {
	Label      string        // snapshot label (may be db.label)
	LabelParts []string      // raw parts: [db, label] or [label]
	Repo       string        // repository name
	Tables     []*ObjectName // optional ON clause table list
	Properties []*Property
	Loc        Loc
}

// Tag implements Node.
func (n *BackupStmt) Tag() NodeTag { return T_BackupStmt }

var _ Node = (*BackupStmt)(nil)

// RestoreStmt represents:
//
//	RESTORE SNAPSHOT [db.]label FROM repo_name
//	    [ON (tbl [PARTITION(p1,...)], ...)]
//	    [PROPERTIES("key"="value", ...)]
type RestoreStmt struct {
	Label      string
	LabelParts []string
	Repo       string
	Tables     []*ObjectName
	Properties []*Property
	Loc        Loc
}

// Tag implements Node.
func (n *RestoreStmt) Tag() NodeTag { return T_RestoreStmt }

var _ Node = (*RestoreStmt)(nil)

// ---------------------------------------------------------------------------
// KILL
// ---------------------------------------------------------------------------

// KillStmt represents:
//
//	KILL [CONNECTION | QUERY] id_or_query_id
type KillStmt struct {
	Kind   string // "" | "CONNECTION" | "QUERY"
	Target string // connection id (numeric string) or query id (string literal)
	Loc    Loc
}

// Tag implements Node.
func (n *KillStmt) Tag() NodeTag { return T_KillStmt }

var _ Node = (*KillStmt)(nil)

// ---------------------------------------------------------------------------
// LOCK / UNLOCK TABLES
// ---------------------------------------------------------------------------

// LockTablesStmt represents:
//
//	LOCK TABLES tbl [AS alias] {READ [LOCAL] | [LOW_PRIORITY] WRITE} [, ...]
type LockTablesStmt struct {
	Items []*LockItem
	Loc   Loc
}

// Tag implements Node.
func (n *LockTablesStmt) Tag() NodeTag { return T_LockTablesStmt }

var _ Node = (*LockTablesStmt)(nil)

// LockItem is one entry in a LOCK TABLES statement.
type LockItem struct {
	Table *ObjectName
	Alias string // optional AS alias
	Mode  string // READ | READ LOCAL | WRITE | LOW_PRIORITY WRITE
	Loc   Loc
}

// Tag implements Node.
func (n *LockItem) Tag() NodeTag { return T_LockItem }

var _ Node = (*LockItem)(nil)

// UnlockTablesStmt represents UNLOCK TABLES.
type UnlockTablesStmt struct {
	Loc Loc
}

// Tag implements Node.
func (n *UnlockTablesStmt) Tag() NodeTag { return T_UnlockTablesStmt }

var _ Node = (*UnlockTablesStmt)(nil)

// ---------------------------------------------------------------------------
// INSTALL / UNINSTALL PLUGIN
// ---------------------------------------------------------------------------

// InstallPluginStmt represents:
//
//	INSTALL PLUGIN FROM 'source_path_or_url'
//	    [PROPERTIES("key"="value", ...)]
//
//	INSTALL PLUGIN FROM SONAME 'library_name'
//	    [PROPERTIES("key"="value", ...)]
type InstallPluginStmt struct {
	Source     string      // path/URL or library name
	IsSoname   bool        // true when SONAME keyword was present
	Properties []*Property
	Loc        Loc
}

// Tag implements Node.
func (n *InstallPluginStmt) Tag() NodeTag { return T_InstallPluginStmt }

var _ Node = (*InstallPluginStmt)(nil)

// UninstallPluginStmt represents:
//
//	UNINSTALL PLUGIN name
type UninstallPluginStmt struct {
	Name string
	Loc  Loc
}

// Tag implements Node.
func (n *UninstallPluginStmt) Tag() NodeTag { return T_UninstallPluginStmt }

var _ Node = (*UninstallPluginStmt)(nil)

// ---------------------------------------------------------------------------
// WARM UP
// ---------------------------------------------------------------------------

// WarmUpStmt represents:
//
//	WARM UP CLUSTER cluster_name FROM cluster_or_compute_group
//	WARM UP COMPUTE GROUP cg WITH TABLE tbl [PARTITION ...]
//
// For best-effort coverage, verb captures the object keyword (CLUSTER /
// COMPUTE GROUP) and the remaining tokens are stored as raw text in Target.
type WarmUpStmt struct {
	Verb   string // "CLUSTER" | "COMPUTE GROUP"
	Target string // raw remaining text after the verb
	Loc    Loc
}

// Tag implements Node.
func (n *WarmUpStmt) Tag() NodeTag { return T_WarmUpStmt }

var _ Node = (*WarmUpStmt)(nil)

// ---------------------------------------------------------------------------
// CLEAN
// ---------------------------------------------------------------------------

// CleanStmt represents:
//
//	CLEAN ALL PROFILE
//	CLEAN LABEL [label] [FROM db]
//	CLEAN QUERY STATS [FROM db] [ALL]
//	... (other CLEAN variants)
//
// Verb captures the identifying keyword(s) after CLEAN; Target is the
// remaining raw text.
type CleanStmt struct {
	Verb   string // ALL PROFILE | LABEL | QUERY STATS | etc.
	Target string // raw remaining text
	Loc    Loc
}

// Tag implements Node.
func (n *CleanStmt) Tag() NodeTag { return T_CleanStmt }

var _ Node = (*CleanStmt)(nil)

// ---------------------------------------------------------------------------
// CANCEL (generic, non-MTMV)
// ---------------------------------------------------------------------------

// CancelStmt represents generic CANCEL commands:
//
//	CANCEL LOAD [FROM db] WHERE ...
//	CANCEL EXPORT [FROM db] WHERE ...
//	CANCEL ALTER TABLE COLUMN|ROLLUP FROM table
//	CANCEL BACKUP [FROM db]
//	CANCEL RESTORE [FROM db]
//	CANCEL BUILD INDEX ON table
//
// Verb identifies the operation (LOAD, EXPORT, ALTER TABLE, BACKUP, RESTORE,
// BUILD INDEX); Target holds the raw remaining text.
type CancelStmt struct {
	Verb   string // LOAD | EXPORT | ALTER TABLE | BACKUP | RESTORE | BUILD INDEX
	Target string // raw remaining text
	Loc    Loc
}

// Tag implements Node.
func (n *CancelStmt) Tag() NodeTag { return T_CancelStmt }

var _ Node = (*CancelStmt)(nil)

// ---------------------------------------------------------------------------
// RECOVER
// ---------------------------------------------------------------------------

// RecoverStmt represents:
//
//	RECOVER DATABASE [name | dbid id] [AS new_name]
//	RECOVER TABLE [db.]name [tblid id] [AS new_name]
//	RECOVER PARTITION [name | pid id] FROM [db.]table [AS new_name]
type RecoverStmt struct {
	Verb      string      // DATABASE | TABLE | PARTITION
	Name      *ObjectName // object name (db/table/partition)
	ID        string      // optional numeric id (dbid/tblid/pid value)
	FromTable *ObjectName // for PARTITION: the FROM table name
	NewName   string      // optional AS new_name
	Loc       Loc
}

// Tag implements Node.
func (n *RecoverStmt) Tag() NodeTag { return T_RecoverStmt }

var _ Node = (*RecoverStmt)(nil)
