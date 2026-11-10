package ast

// IsDDL reports whether n is a data-definition statement. It covers exactly the
// set pingcap/tidb's ast.DDLNode marker interface covers (the types embedding
// ddlNode in the bytebase/tidb fork's pkg/parser/ast), so Bytebase's prior-backup
// safety gate — which refuses a backup batch that mixes DDL and DML — behaves
// identically once it moves off the pingcap bridge onto omni.
//
// Several statements that read like admin or lock operations are DDLNode in
// pingcap and are kept here for that parity: LOCK/UNLOCK TABLES, OPTIMIZE TABLE,
// REPAIR TABLE, ADMIN CLEANUP TABLE LOCK, and ADMIN REPAIR TABLE. GRANT/REVOKE,
// CREATE USER/ROLE/EVENT/TRIGGER/FUNCTION, and the SET/SHOW/other-ADMIN
// statements are not DDLNode in pingcap, so IsDDL is false for them.
//
// Three statements are a dedicated DDLNode type in pingcap but a different shape
// in omni, so they are matched here by omni's shape rather than a same-named
// type (verified by parsing):
//   - DROP VIEW is a DropTableStmt (IsView=true) in pingcap but its own
//     DropViewStmt here; CREATE VIEW is a CreateViewStmt on both sides.
//   - ADMIN CLEANUP TABLE LOCK is pingcap's CleanupTableLockStmt but an AdminStmt
//     with Kind AdminCleanupTableLock here.
//   - ADMIN REPAIR TABLE t CREATE TABLE ... is pingcap's RepairTableStmt but an
//     AdminStmt with Kind AdminRepairTable here (omni's RepairTableStmt covers
//     the plain REPAIR TABLE form instead).
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
		// Two pingcap DDL types are a dedicated statement there but an AdminStmt
		// here, so they are matched by admin kind: CleanupTableLockStmt
		// (ADMIN CLEANUP TABLE LOCK) and RepairTableStmt (ADMIN REPAIR TABLE t
		// CREATE TABLE ...). Every other ADMIN variant (CHECK, SHOW DDL, …) is
		// not DDL.
		return s.Kind == AdminCleanupTableLock || s.Kind == AdminRepairTable
	}
	return false
}
