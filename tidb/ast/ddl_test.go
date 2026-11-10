package ast

import "testing"

// TestIsDDL pins the exact statement set IsDDL classifies as DDL, matching
// pingcap/tidb's ast.DDLNode marker so Bytebase's prior-backup safety gate
// behaves identically off the pingcap bridge.
func TestIsDDL(t *testing.T) {
	ddl := []Node{
		&AlterDatabaseStmt{}, &AlterPlacementPolicyStmt{}, &AlterResourceGroupStmt{},
		&AlterSequenceStmt{}, &AlterTableStmt{},
		&CreateDatabaseStmt{}, &CreateIndexStmt{}, &CreatePlacementPolicyStmt{},
		&CreateResourceGroupStmt{}, &CreateSequenceStmt{}, &CreateTableStmt{}, &CreateViewStmt{},
		&DropDatabaseStmt{}, &DropIndexStmt{}, &DropPlacementPolicyStmt{}, &DropResourceGroupStmt{},
		&DropSequenceStmt{}, &DropTableStmt{}, &DropViewStmt{},
		&FlashbackDatabaseStmt{}, &FlashbackTableStmt{}, &FlashbackToTimestampStmt{},
		&LockTablesStmt{}, &UnlockTablesStmt{}, &OptimizeTableStmt{},
		&RecoverTableStmt{}, &RenameTableStmt{}, &RepairTableStmt{}, &TruncateStmt{},
		&AdminStmt{Kind: AdminCleanupTableLock},
	}
	for _, n := range ddl {
		if !IsDDL(n) {
			t.Errorf("IsDDL(%T) = false, want true", n)
		}
	}

	notDDL := []Node{
		// ADMIN kinds other than CLEANUP TABLE LOCK are not DDL.
		&AdminStmt{Kind: AdminShowDDL}, &AdminStmt{Kind: AdminCheckTable},
		// Privilege / account / routine / event / trigger statements: DDLNode in
		// neither pingcap nor here.
		&GrantStmt{}, &RevokeStmt{}, &CreateUserStmt{}, &AlterUserStmt{}, &DropUserStmt{},
		&CreateRoleStmt{}, &CreateEventStmt{}, &CreateTriggerStmt{}, &CreateFunctionStmt{},
		// SET / SHOW / USE and DML.
		&SetStmt{}, &ShowStmt{}, &UseStmt{},
		&SelectStmt{}, &InsertStmt{}, &UpdateStmt{}, &DeleteStmt{},
	}
	for _, n := range notDDL {
		if IsDDL(n) {
			t.Errorf("IsDDL(%T) = true, want false", n)
		}
	}
	// A nil Node is not DDL and must not panic.
	if IsDDL(nil) {
		t.Errorf("IsDDL(nil) = true, want false")
	}
}
