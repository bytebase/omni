package parser

import (
	"testing"

	ast "github.com/bytebase/omni/oracle/ast"
)

func TestP2CatalogAdminPositive(t *testing.T) {
	cases := []string{
		"ALTER DATABASE DICTIONARY ENCRYPT CREDENTIALS",
		"ALTER DATABASE LINK remote_db CONNECT TO admin IDENTIFIED BY pass",
		"ALTER DATABASE OPEN RESETLOGS",
		"ALTER DISKGROUP dg1 ADD DISK '/dev/sdc1' NAME disk3",
		"ALTER PLUGGABLE DATABASE pdb1 OPEN",
		"ALTER RESOURCE COST CPU_PER_SESSION 100",
		"CREATE CONTROLFILE REUSE DATABASE mydb NORESETLOGS NOARCHIVELOG",
		"CREATE DATABASE LINK remote_db CONNECT TO admin IDENTIFIED BY pass USING 'srv'",
		"CREATE DATABASE mydb USER SYS IDENTIFIED BY password",
		"CREATE DISKGROUP dg1 DISK '/dev/sda1' NAME disk1",
		"CREATE PFILE = '/tmp/init.ora' FROM SPFILE = '/tmp/spfile.ora'",
		"CREATE PLUGGABLE DATABASE pdb1 ADMIN USER pdbadmin IDENTIFIED BY pass",
		"CREATE RESTORE POINT rp1 GUARANTEE FLASHBACK DATABASE",
		"CREATE SPFILE = '/tmp/spfile.ora' FROM PFILE = '/tmp/init.ora'",
		"CREATE TABLESPACE SET ts_set",
		"CREATE TABLESPACE users DATAFILE '/u01/users01.dbf' SIZE 100M",
		"DROP DATABASE LINK remote_db",
		"DROP DATABASE",
		"DROP DISKGROUP dg1 INCLUDING CONTENTS",
		"DROP PLUGGABLE DATABASE pdb1 INCLUDING DATAFILES",
		"DROP RESTORE POINT rp1",
		"DROP TABLESPACE SET ts_set INCLUDING CONTENTS",
		"DROP TABLESPACE users INCLUDING CONTENTS",
		"FLASHBACK DATABASE TO RESTORE POINT rp1",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q): %v", sql, err)
			}
		})
	}
}

func TestP2CatalogAdminNegative(t *testing.T) {
	cases := []string{
		"ALTER DATABASE LINK CONNECT TO admin",
		"CREATE DATABASE LINK CONNECT TO admin",
		"CREATE DISKGROUP DISK '/dev/sda1'",
		"CREATE PFILE FROM",
		"DROP DISKGROUP",
		"FLASHBACK DATABASE TO RESTORE POINT",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			assertParseErrorContains(t, sql, "syntax error")
		})
	}
}

func TestP2CatalogAdminLoc(t *testing.T) {
	cases := []string{
		"ALTER DATABASE OPEN RESETLOGS",
		"CREATE DATABASE mydb USER SYS IDENTIFIED BY password",
		"CREATE DISKGROUP dg1 DISK '/dev/sda1' NAME disk1",
		"CREATE RESTORE POINT rp1 GUARANTEE FLASHBACK DATABASE",
		"FLASHBACK DATABASE TO RESTORE POINT rp1",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			CheckLocations(t, sql)
		})
	}
}

func TestP2AlterTablespaceEncryptionPreservesPayload(t *testing.T) {
	result := ParseAndCheck(t, "ALTER TABLESPACE users ENCRYPTION ONLINE USING 'AES256' ENCRYPT FILE_NAME_CONVERT = ('old', 'new') FLASHBACK ON")
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.AlterTablespaceStmt)
	if !ok {
		t.Fatalf("expected *AlterTablespaceStmt, got %T", raw.Stmt)
	}
	if stmt.Encryption != "ENCRYPTION ONLINE USING AES256 ENCRYPT FILE_NAME_CONVERT = ( old , new )" {
		t.Fatalf("Encryption = %q", stmt.Encryption)
	}
	if stmt.Flashback != "ON" {
		t.Fatalf("Flashback = %q, want ON", stmt.Flashback)
	}
}
