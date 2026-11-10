package ast

// IsDDL reports whether n is a data-definition statement. It covers exactly the
// set pingcap/tidb's ast.DDLNode marker interface covers (the types embedding
// ddlNode in the bytebase/tidb fork's pkg/parser/ast), so Bytebase's prior-backup
// safety gate — which refuses a backup batch that mixes DDL and DML — behaves
// identically once it moves off the pingcap bridge onto omni.
//
// Several statements that read like admin or lock operations are DDLNode in
// pingcap and are kept here for that parity: LOCK/UNLOCK TABLES, OPTIMIZE TABLE,
// REPAIR TABLE, and ADMIN CLEANUP TABLE LOCK (pingcap's CleanupTableLockStmt,
// which omni parses as an AdminStmt). GRANT/REVOKE, CREATE USER/ROLE/EVENT/
// TRIGGER/FUNCTION, and the SET/SHOW/other-ADMIN statements are not DDLNode in
// pingcap, so IsDDL is false for them.
//
// DropViewStmt has no pingcap counterpart: pingcap parses DROP VIEW as a
// DropTableStmt (IsView=true), which is a DDLNode, so DROP VIEW counts as DDL
// there. omni parses it into its own DropViewStmt, so that type is included here
// to keep DROP VIEW classified as DDL exactly as pingcap does. CREATE VIEW is
// already a CreateViewStmt on both sides.
func IsDDL(n Node) bool {
	switch s := n.(type) {
	case *AlterDatabaseStmt, *AlterPlacementPolicyStmt, *AlterResourceGroupStmt,
		*AlterSequenceStmt, *AlterTableStmt,
		*CreateDatabaseStmt, *CreateIndexStmt, *CreatePlacementPolicyStmt,
		*CreateResourceGroupStmt, *CreateSequenceStmt, *CreateTableStmt, *CreateViewStmt,
		*DropDatabaseStmt, *DropIndexStmt, *DropPlacementPolicyStmt, *DropResourceGroupStmt,
		*DropSequenceStmt, *DropTableStmt, *DropViewStmt,
		*FlashbackDatabaseStmt, *FlashbackTableStmt, *FlashbackToTimestampStmt,
		*LockTablesStmt, *UnlockTablesStmt, *OptimizeTableStmt,
		*RecoverTableStmt, *RenameTableStmt, *RepairTableStmt, *TruncateStmt:
		return true
	case *AdminStmt:
		// pingcap's 29th DDL type is CleanupTableLockStmt (ADMIN CLEANUP TABLE
		// LOCK), which omni folds into AdminStmt. Only that admin kind is DDL;
		// every other ADMIN variant is not.
		return s.Kind == AdminCleanupTableLock
	}
	return false
}
