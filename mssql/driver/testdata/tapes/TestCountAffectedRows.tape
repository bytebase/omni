-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestCountAffectedRows default#1
> exec
| CREATE TABLE c_big (id INT PRIMARY KEY, s_id INT NOT NULL, v INT);
| INSERT INTO c_big SELECT TOP 1000 ROW_NUMBER() OVER (ORDER BY (SELECT NULL)), (ROW_NUMBER() OVER (ORDER BY (SELECT NULL)) % 100) + 1, 0 FROM sys.all_objects;
| CREATE TABLE c_archive (id INT PRIMARY KEY, s_id INT NOT NULL, v INT);
| UPDATE STATISTICS c_big;
< ok 1000 -
> exec
| SET SHOWPLAN_ALL ON;
< ok 0 -
> query
| DELETE FROM c_big WHERE s_id <= 50;
< result
< column "StmtText" "NVARCHAR"
< column "StmtId" "INT"
< column "NodeId" "INT"
< column "Parent" "INT"
< column "PhysicalOp" "NVARCHAR"
< column "LogicalOp" "NVARCHAR"
< column "Argument" "NVARCHAR"
< column "DefinedValues" "NVARCHAR"
< column "EstimateRows" "REAL"
< column "EstimateIO" "REAL"
< column "EstimateCPU" "REAL"
< column "AvgRowSize" "INT"
< column "TotalSubtreeCost" "REAL"
< column "OutputList" "NVARCHAR"
< column "Warnings" "NVARCHAR"
< column "Type" "NVARCHAR"
< column "Parallel" "BIT"
< column "EstimateExecutions" "REAL"
< row
: string
| DELETE FROM c_big WHERE s_id <= 50;
: int64 1
: int64 1
: int64 0
: null
: null
: string
| 1
: null
: float64 500
: null
: null
: null
: float64 0.019921859726309776
: null
: null
: string
| DELETE
: bool false
: null
< row
: string
|   |--Clustered Index Delete(OBJECT:([omni_TestCountAffectedRows].[dbo].[c_big].[PK__c_big__standin1]))
: int64 1
: int64 2
: int64 1
: string
| Clustered Index Delete
: string
| Delete
: string
| OBJECT:([omni_TestCountAffectedRows].[dbo].[c_big].[PK__c_big__standin1])
: null
: float64 500
: float64 0.013378378003835678
: float64 0.0005000000237487257
: int64 9
: float64 0.019921859726309776
: null
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|        |--Clustered Index Scan(OBJECT:([omni_TestCountAffectedRows].[dbo].[c_big].[PK__c_big__standin1]), WHERE:([omni_TestCountAffectedRows].[dbo].[c_big].[s_id]<=CONVERT_IMPLICIT(int,[@1],0)) ORDERED)
: int64 1
: int64 3
: int64 2
: string
| Clustered Index Scan
: string
| Clustered Index Scan
: string
| OBJECT:([omni_TestCountAffectedRows].[dbo].[c_big].[PK__c_big__standin1]), WHERE:([omni_TestCountAffectedRows].[dbo].[c_big].[s_id]<=CONVERT_IMPLICIT(int,[@1],0)) ORDERED
: string
| [omni_TestCountAffectedRows].[dbo].[c_big].[id]
: float64 500
: float64 0.004606481641530991
: float64 0.0012570000253617764
: int64 15
: float64 0.005863481666892767
: string
| [omni_TestCountAffectedRows].[dbo].[c_big].[id]
: null
: string
| PLAN_ROW
: bool false
: float64 1
> exec
| SET SHOWPLAN_ALL OFF;
< ok 0 -
> exec
| SET SHOWPLAN_ALL ON;
< ok 0 -
> query
| SELECT * FROM app.missing_table;
< error
| 208 1 16 1 Invalid object name 'app.missing_table'.
> exec
| SET SHOWPLAN_ALL OFF;
< ok 0 -
