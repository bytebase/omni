package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiagnose(t *testing.T) {
	require.Nil(t, Diagnose(""))
	require.Nil(t, Diagnose("SELECT 1 FROM DUAL"))
	require.Nil(t, Diagnose("SELECT 1 FROM DUAL;\nSELECT 2 FROM DUAL  \n"))

	got := Diagnose("SELECT 1 FROM DUAL;\nSELECT * FROM")
	require.NotNil(t, got)
	require.Equal(t, "SELECT 1 FROM DUAL;\nSELECT * FROM;", got.Text)
	require.GreaterOrEqual(t, got.Offset, len("SELECT 1 FROM DUAL;\n"))
	require.NotEmpty(t, got.Message)

	// Trigger with REFERENCING, ALTER SYSTEM, and interval-partitioned tables
	// with USING INDEX LOCAL / DBMS_METADATA storage attributes diagnose clean.
	for _, valid := range []string{
		`CREATE OR REPLACE TRIGGER trg
BEFORE INSERT OR UPDATE OF col1, col2 ON tbl
REFERENCING OLD AS o NEW AS n
FOR EACH ROW
WHEN (n.col1 > 0)
BEGIN
  :n.col2 := :o.col2 + 1;
END;`,
		"alter system set dg_broker_start=true;",
		"alter system switch logfile;",
		`CREATE TABLE APP.ORDER_EVENT
(
    EVENT_CODE  VARCHAR2(250) NOT NULL,
    CREATE_DATE DATE DEFAULT SYSDATE NOT NULL,
    CONSTRAINT PK_ORDER_EVENT PRIMARY KEY (EVENT_CODE, CREATE_DATE) USING INDEX LOCAL
)
PARTITION BY RANGE (CREATE_DATE)
INTERVAL (NUMTOYMINTERVAL(1, 'MONTH'))
(
    PARTITION P202607 VALUES LESS THAN (DATE '2026-07-01')
)`,
		`CREATE TABLE "S"."T" ("A" NUMBER(*,0) NOT NULL ENABLE, CONSTRAINT "PK_T" PRIMARY KEY ("A") USING INDEX PCTFREE 10 INITRANS 2 MAXTRANS 255 COMPUTE STATISTICS TABLESPACE "USERS" ENABLE) PCTFREE 10 PCTUSED 40 INITRANS 1 MAXTRANS 255 TABLESPACE "USERS"`,
	} {
		require.Nil(t, Diagnose(valid), valid)
	}
}
