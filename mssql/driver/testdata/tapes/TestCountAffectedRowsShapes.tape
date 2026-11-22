-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestCountAffectedRowsShapes default#1
> exec
| CREATE TABLE dbo.t (id INT IDENTITY(1,1) PRIMARY KEY, c1 INT NULL);
< ok 0 -
> exec
| INSERT INTO dbo.t (c1) VALUES (1),(2),(3),(4),(5);
< ok 5 -
> exec
| CREATE TABLE dbo.s (id INT PRIMARY KEY, flag BIT NOT NULL);
| CREATE TABLE dbo.big (id INT PRIMARY KEY, s_id INT NOT NULL, v INT NOT NULL);
| CREATE TABLE dbo.big_archive (id INT PRIMARY KEY, s_id INT NOT NULL, v INT NOT NULL);
< ok 0 -
> exec
| INSERT INTO dbo.s SELECT n, CASE WHEN n <= 10 THEN 1 ELSE 0 END FROM (SELECT TOP (100) ROW_NUMBER() OVER (ORDER BY (SELECT NULL)) AS n FROM sys.all_objects) AS seq;
| INSERT INTO dbo.big SELECT n, n % 100 + 1, 0 FROM (SELECT TOP (10000) ROW_NUMBER() OVER (ORDER BY (SELECT NULL)) AS n FROM sys.all_objects AS a CROSS JOIN sys.all_objects AS b) AS seq;
< ok 10100 -
> exec
| CREATE INDEX big_s_id ON dbo.big (s_id);
| CREATE STATISTICS s_flag ON dbo.s (flag);
< ok 0 -
> exec
| UPDATE STATISTICS dbo.t WITH FULLSCAN;
| UPDATE STATISTICS dbo.s WITH FULLSCAN;
| UPDATE STATISTICS dbo.big WITH FULLSCAN;
< ok 0 -
> exec
| SET SHOWPLAN_ALL ON;
< ok 0 -
> query
| INSERT INTO dbo.t DEFAULT VALUES;
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
| INSERT INTO dbo.t DEFAULT VALUES;
: int64 1
: int64 1
: int64 0
: null
: null
: string
| 1
: null
: float64 1
: null
: null
: null
: float64 0.01000225730240345
: null
: null
: string
| INSERT
: bool false
: null
< row
: string
|   |--Clustered Index Insert(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[t].[PK__t__standin1]), SET:([omni_TestCountAffectedRowsShapes].[dbo].[t].[id] = [Expr1002],[omni_TestCountAffectedRowsShapes].[dbo].[t].[c1] = NULL))
: int64 1
: int64 2
: int64 1
: string
| Clustered Index Insert
: string
| Insert
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[t].[PK__t__standin1]), SET:([omni_TestCountAffectedRowsShapes].[dbo].[t].[id] = [Expr1002],[omni_TestCountAffectedRowsShapes].[dbo].[t].[c1] = NULL)
: null
: float64 1
: float64 0.009999999776482582
: float64 9.999999974752427e-07
: int64 9
: float64 0.01000225730240345
: null
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|        |--Compute Scalar(DEFINE:([Expr1002]=getidentity((901578250),(5),NULL)))
: int64 1
: int64 3
: int64 2
: string
| Compute Scalar
: string
| Compute Scalar
: string
| DEFINE:([Expr1002]=getidentity((901578250),(5),NULL))
: string
| [Expr1002]=getidentity((901578250),(5),NULL)
: float64 1
: float64 0
: float64 1.0000000116860974e-07
: int64 11
: float64 1.257000008081377e-06
: string
| [Expr1002]
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|             |--Constant Scan
: int64 1
: int64 4
: int64 3
: string
| Constant Scan
: string
| Constant Scan
: null
: null
: float64 1
: float64 0
: float64 1.157000042439904e-06
: int64 9
: float64 1.157000042439904e-06
: null
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
| -- header comment
| UPDATE dbo.t SET c1 = 99 WHERE id = 1;
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
| -- header comment
| UPDATE dbo.t SET c1 = 99 WHERE id = 1;
: int64 1
: int64 1
: int64 0
: null
: null
: string
| 1
: null
: float64 1
: null
: null
: null
: float64 0.013284100219607353
: null
: null
: string
| UPDATE
: bool false
: null
< row
: string
|   |--Clustered Index Update(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[t].[PK__t__standin1]), SET:([omni_TestCountAffectedRowsShapes].[dbo].[t].[c1] = [@1]), WHERE:([omni_TestCountAffectedRowsShapes].[dbo].[t].[id]=CONVERT_IMPLICIT(int,[@2],0)))
: int64 1
: int64 2
: int64 1
: string
| Clustered Index Update
: string
| Update
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[t].[PK__t__standin1]), SET:([omni_TestCountAffectedRowsShapes].[dbo].[t].[c1] = [@1]), WHERE:([omni_TestCountAffectedRowsShapes].[dbo].[t].[id]=CONVERT_IMPLICIT(int,[@2],0))
: null
: float64 1
: float64 0.009999999776482582
: float64 9.999999974752427e-07
: int64 9
: float64 0.013284100219607353
: null
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
| UPDATE dbo.t SET c1 = 99 WHERE id = 1;
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
| UPDATE dbo.t SET c1 = 99 WHERE id = 1;
: int64 1
: int64 1
: int64 0
: null
: null
: string
| 1
: null
: float64 1
: null
: null
: null
: float64 0.013284100219607353
: null
: null
: string
| UPDATE
: bool false
: null
< row
: string
|   |--Clustered Index Update(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[t].[PK__t__standin1]), SET:([omni_TestCountAffectedRowsShapes].[dbo].[t].[c1] = [@1]), WHERE:([omni_TestCountAffectedRowsShapes].[dbo].[t].[id]=CONVERT_IMPLICIT(int,[@2],0)))
: int64 1
: int64 2
: int64 1
: string
| Clustered Index Update
: string
| Update
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[t].[PK__t__standin1]), SET:([omni_TestCountAffectedRowsShapes].[dbo].[t].[c1] = [@1]), WHERE:([omni_TestCountAffectedRowsShapes].[dbo].[t].[id]=CONVERT_IMPLICIT(int,[@2],0))
: null
: float64 1
: float64 0.009999999776482582
: float64 9.999999974752427e-07
: int64 9
: float64 0.013284100219607353
: null
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
| UPDATE dbo.t SET c1 = c1 + 1 WHERE id > 2;
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
| UPDATE dbo.t SET c1 = c1 + 1 WHERE id > 2;
: int64 1
: int64 1
: int64 0
: null
: null
: string
| 1
: null
: float64 3
: null
: null
: null
: float64 0.013288600370287895
: null
: null
: string
| UPDATE
: bool false
: null
< row
: string
|   |--Clustered Index Update(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[t].[PK__t__standin1]), SET:([omni_TestCountAffectedRowsShapes].[dbo].[t].[c1] = [Expr1002]), DEFINE:([Expr1002]=[omni_TestCountAffectedRowsShapes].[dbo].[t].[c1]+[@1]), WHERE:([omni_TestCountAffectedRowsShapes].[dbo].[t].[id] > CONVERT_IMPLICIT(int,[@2],0)))
: int64 1
: int64 2
: int64 1
: string
| Clustered Index Update
: string
| Update
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[t].[PK__t__standin1]), SET:([omni_TestCountAffectedRowsShapes].[dbo].[t].[c1] = [Expr1002]), DEFINE:([Expr1002]=[omni_TestCountAffectedRowsShapes].[dbo].[t].[c1]+[@1]), WHERE:([omni_TestCountAffectedRowsShapes].[dbo].[t].[id] > CONVERT_IMPLICIT(int,[@2],0))
: string
| [Expr1002]=[omni_TestCountAffectedRowsShapes].[dbo].[t].[c1]+[@1]
: float64 3
: float64 0.009999999776482582
: float64 3.000000106112566e-06
: int64 9
: float64 0.013288600370287895
: null
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
| DELETE FROM dbo.big WHERE s_id <= 3;
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
| DELETE FROM dbo.big WHERE s_id <= 3;
: int64 1
: int64 1
: int64 0
: null
: null
: string
| 1
: null
: float64 300
: null
: null
: null
: float64 0.024212000891566277
: null
: null
: string
| DELETE
: bool false
: null
< row
: string
|   |--Clustered Index Delete(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[PK__big__standin2]), OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[big_s_id]) WITH UNORDERED PREFETCH)
: int64 1
: int64 2
: int64 1
: string
| Clustered Index Delete
: string
| Delete
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[PK__big__standin2]), OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[big_s_id]) WITH UNORDERED PREFETCH
: null
: float64 300
: float64 0.019999999552965164
: float64 0.0006000000284984708
: int64 9
: float64 0.024212000891566277
: null
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|        |--Index Seek(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[big_s_id]), SEEK:([omni_TestCountAffectedRowsShapes].[dbo].[big].[s_id] <= CONVERT_IMPLICIT(int,[@1],0)) ORDERED FORWARD)
: int64 1
: int64 4
: int64 2
: string
| Index Seek
: string
| Index Seek
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[big_s_id]), SEEK:([omni_TestCountAffectedRowsShapes].[dbo].[big].[s_id] <= CONVERT_IMPLICIT(int,[@1],0)) ORDERED FORWARD
: string
| [omni_TestCountAffectedRowsShapes].[dbo].[big].[id]
: float64 300
: float64 0.0031250000465661287
: float64 0.0004870000120718032
: int64 11
: float64 0.0036120000295341015
: string
| [omni_TestCountAffectedRowsShapes].[dbo].[big].[id]
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
| DELETE b FROM dbo.big AS b JOIN dbo.s AS s ON b.s_id = s.id WHERE s.flag = 1;
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
| DELETE b FROM dbo.big AS b JOIN dbo.s AS s ON b.s_id = s.id WHERE s.flag = 1;
: int64 1
: int64 1
: int64 0
: null
: null
: string
| 1
: null
: float64 1000
: null
: null
: null
: float64 0.0711275041103363
: null
: null
: string
| DELETE
: bool false
: null
< row
: string
|   |--Clustered Index Delete(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[PK__big__standin2] AS [b]), OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[big_s_id] AS [b]) WITH UNORDERED PREFETCH)
: int64 1
: int64 2
: int64 1
: string
| Clustered Index Delete
: string
| Delete
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[PK__big__standin2] AS [b]), OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[big_s_id] AS [b]) WITH UNORDERED PREFETCH
: null
: float64 1000
: float64 0.036216218024492264
: float64 0.0020000000949949026
: int64 9
: float64 0.0711275041103363
: null
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|        |--Nested Loops(Inner Join, OUTER REFERENCES:([s].[id]))
: int64 1
: int64 4
: int64 2
: string
| Nested Loops
: string
| Inner Join
: string
| OUTER REFERENCES:([s].[id])
: null
: float64 1000
: float64 0
: float64 0.00418000016361475
: int64 11
: float64 0.032911285758018494
: string
| [b].[id]
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|             |--Clustered Index Scan(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[s].[PK__s__standin3] AS [s]), WHERE:([omni_TestCountAffectedRowsShapes].[dbo].[s].[flag] as [s].[flag]=(1)))
: int64 1
: int64 5
: int64 4
: string
| Clustered Index Scan
: string
| Clustered Index Scan
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[s].[PK__s__standin3] AS [s]), WHERE:([omni_TestCountAffectedRowsShapes].[dbo].[s].[flag] as [s].[flag]=(1))
: string
| [s].[id]
: float64 10
: float64 0.0031250000465661287
: float64 0.00026699999580159783
: int64 12
: float64 0.0033919999841600657
: string
| [s].[id]
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|             |--Index Seek(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[big_s_id] AS [b]), SEEK:([b].[s_id]=[omni_TestCountAffectedRowsShapes].[dbo].[s].[id] as [s].[id]) ORDERED FORWARD)
: int64 1
: int64 6
: int64 4
: string
| Index Seek
: string
| Index Seek
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[big_s_id] AS [b]), SEEK:([b].[s_id]=[omni_TestCountAffectedRowsShapes].[dbo].[s].[id] as [s].[id]) ORDERED FORWARD
: string
| [b].[id]
: float64 100
: float64 0.0031250000465661287
: float64 0.00026699999580159783
: int64 11
: float64 0.025291284546256065
: string
| [b].[id]
: null
: string
| PLAN_ROW
: bool false
: float64 10
> exec
| SET SHOWPLAN_ALL OFF;
< ok 0 -
> exec
| SET SHOWPLAN_ALL ON;
< ok 0 -
> query
| UPDATE b SET v = 1 FROM dbo.big AS b JOIN dbo.s AS s ON b.s_id = s.id WHERE s.flag = 1;
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
| UPDATE b SET v = 1 FROM dbo.big AS b JOIN dbo.s AS s ON b.s_id = s.id WHERE s.flag = 1;
: int64 1
: int64 1
: int64 0
: null
: null
: string
| 1
: null
: float64 1000
: null
: null
: null
: float64 0.05492209643125534
: null
: null
: string
| UPDATE
: bool false
: null
< row
: string
|   |--Clustered Index Update(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[PK__big__standin2] AS [b]), SET:([omni_TestCountAffectedRowsShapes].[dbo].[big].[v] as [b].[v] = [Expr1002]) WITH UNORDERED PREFETCH)
: int64 1
: int64 2
: int64 1
: string
| Clustered Index Update
: string
| Update
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[PK__big__standin2] AS [b]), SET:([omni_TestCountAffectedRowsShapes].[dbo].[big].[v] as [b].[v] = [Expr1002]) WITH UNORDERED PREFETCH
: null
: float64 1000
: float64 0.02081081084907055
: float64 0.0010000000474974513
: int64 9
: float64 0.05492209643125534
: null
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|        |--Compute Scalar(DEFINE:([Expr1002]=(1)))
: int64 1
: int64 4
: int64 2
: string
| Compute Scalar
: string
| Compute Scalar
: string
| DEFINE:([Expr1002]=(1))
: string
| [Expr1002]=(1)
: float64 1000
: float64 0
: float64 9.999999747378752e-05
: int64 15
: float64 0.03311128541827202
: string
| [b].[id], [Expr1002]
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|             |--Nested Loops(Inner Join, OUTER REFERENCES:([s].[id]))
: int64 1
: int64 6
: int64 4
: string
| Nested Loops
: string
| Inner Join
: string
| OUTER REFERENCES:([s].[id])
: null
: float64 1000
: float64 0
: float64 0.00418000016361475
: int64 11
: float64 0.032911285758018494
: string
| [b].[id]
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|                  |--Clustered Index Scan(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[s].[PK__s__standin3] AS [s]), WHERE:([omni_TestCountAffectedRowsShapes].[dbo].[s].[flag] as [s].[flag]=(1)))
: int64 1
: int64 7
: int64 6
: string
| Clustered Index Scan
: string
| Clustered Index Scan
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[s].[PK__s__standin3] AS [s]), WHERE:([omni_TestCountAffectedRowsShapes].[dbo].[s].[flag] as [s].[flag]=(1))
: string
| [s].[id]
: float64 10
: float64 0.0031250000465661287
: float64 0.00026699999580159783
: int64 12
: float64 0.0033919999841600657
: string
| [s].[id]
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|                  |--Index Seek(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[big_s_id] AS [b]), SEEK:([b].[s_id]=[omni_TestCountAffectedRowsShapes].[dbo].[s].[id] as [s].[id]) ORDERED FORWARD)
: int64 1
: int64 8
: int64 6
: string
| Index Seek
: string
| Index Seek
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[big_s_id] AS [b]), SEEK:([b].[s_id]=[omni_TestCountAffectedRowsShapes].[dbo].[s].[id] as [s].[id]) ORDERED FORWARD
: string
| [b].[id]
: float64 100
: float64 0.0031250000465661287
: float64 0.00026699999580159783
: int64 11
: float64 0.025291284546256065
: string
| [b].[id]
: null
: string
| PLAN_ROW
: bool false
: float64 10
> exec
| SET SHOWPLAN_ALL OFF;
< ok 0 -
> exec
| SET SHOWPLAN_ALL ON;
< ok 0 -
> query
| DELETE FROM dbo.big WHERE s_id IN (SELECT id FROM dbo.s WHERE flag = 1);
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
| DELETE FROM dbo.big WHERE s_id IN (SELECT id FROM dbo.s WHERE flag = 1);
: int64 1
: int64 1
: int64 0
: null
: null
: string
| 1
: null
: float64 1000
: null
: null
: null
: float64 0.0711275041103363
: null
: null
: string
| DELETE
: bool false
: null
< row
: string
|   |--Clustered Index Delete(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[PK__big__standin2]), OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[big_s_id]) WITH UNORDERED PREFETCH)
: int64 1
: int64 2
: int64 1
: string
| Clustered Index Delete
: string
| Delete
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[PK__big__standin2]), OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[big_s_id]) WITH UNORDERED PREFETCH
: null
: float64 1000
: float64 0.036216218024492264
: float64 0.0020000000949949026
: int64 9
: float64 0.0711275041103363
: null
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|        |--Nested Loops(Inner Join, OUTER REFERENCES:([omni_TestCountAffectedRowsShapes].[dbo].[s].[id]))
: int64 1
: int64 5
: int64 2
: string
| Nested Loops
: string
| Inner Join
: string
| OUTER REFERENCES:([omni_TestCountAffectedRowsShapes].[dbo].[s].[id])
: null
: float64 1000
: float64 0
: float64 0.00418000016361475
: int64 11
: float64 0.032911285758018494
: string
| [omni_TestCountAffectedRowsShapes].[dbo].[big].[id]
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|             |--Clustered Index Scan(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[s].[PK__s__standin3]), WHERE:([omni_TestCountAffectedRowsShapes].[dbo].[s].[flag]=(1)))
: int64 1
: int64 6
: int64 5
: string
| Clustered Index Scan
: string
| Clustered Index Scan
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[s].[PK__s__standin3]), WHERE:([omni_TestCountAffectedRowsShapes].[dbo].[s].[flag]=(1))
: string
| [omni_TestCountAffectedRowsShapes].[dbo].[s].[id]
: float64 10
: float64 0.0031250000465661287
: float64 0.00026699999580159783
: int64 12
: float64 0.0033919999841600657
: string
| [omni_TestCountAffectedRowsShapes].[dbo].[s].[id]
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|             |--Index Seek(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[big_s_id]), SEEK:([omni_TestCountAffectedRowsShapes].[dbo].[big].[s_id]=[omni_TestCountAffectedRowsShapes].[dbo].[s].[id]) ORDERED FORWARD)
: int64 1
: int64 7
: int64 5
: string
| Index Seek
: string
| Index Seek
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[big_s_id]), SEEK:([omni_TestCountAffectedRowsShapes].[dbo].[big].[s_id]=[omni_TestCountAffectedRowsShapes].[dbo].[s].[id]) ORDERED FORWARD
: string
| [omni_TestCountAffectedRowsShapes].[dbo].[big].[id]
: float64 100
: float64 0.0031250000465661287
: float64 0.00026699999580159783
: int64 11
: float64 0.025291284546256065
: string
| [omni_TestCountAffectedRowsShapes].[dbo].[big].[id]
: null
: string
| PLAN_ROW
: bool false
: float64 10
> exec
| SET SHOWPLAN_ALL OFF;
< ok 0 -
> exec
| SET SHOWPLAN_ALL ON;
< ok 0 -
> query
| INSERT INTO dbo.big_archive SELECT id, s_id, v FROM dbo.big WHERE s_id > 95;
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
| INSERT INTO dbo.big_archive SELECT id, s_id, v FROM dbo.big WHERE s_id > 95;
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
: float64 0.049179356545209885
: null
: null
: string
| INSERT
: bool false
: null
< row
: string
|   |--Clustered Index Insert(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[PK__big_arch__standin4]), SET:([omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[id] = [omni_TestCountAffectedRowsShapes].[dbo].[big].[id],[omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[s_id] = [omni_TestCountAffectedRowsShapes].[dbo].[big].[s_id],[omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[v] = [omni_TestCountAffectedRowsShapes].[dbo].[big].[v]))
: int64 1
: int64 2
: int64 1
: string
| Clustered Index Insert
: string
| Insert
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[PK__big_arch__standin4]), SET:([omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[id] = [omni_TestCountAffectedRowsShapes].[dbo].[big].[id],[omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[s_id] = [omni_TestCountAffectedRowsShapes].[dbo].[big].[s_id],[omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[v] = [omni_TestCountAffectedRowsShapes].[dbo].[big].[v])
: null
: float64 500
: float64 0.01107883732765913
: float64 0.0005000000237487257
: int64 9
: float64 0.049179356545209885
: null
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|        |--Clustered Index Scan(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[PK__big__standin2]), WHERE:([omni_TestCountAffectedRowsShapes].[dbo].[big].[s_id]>(95)))
: int64 1
: int64 3
: int64 2
: string
| Clustered Index Scan
: string
| Clustered Index Scan
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[PK__big__standin2]), WHERE:([omni_TestCountAffectedRowsShapes].[dbo].[big].[s_id]>(95))
: string
| [omni_TestCountAffectedRowsShapes].[dbo].[big].[id], [omni_TestCountAffectedRowsShapes].[dbo].[big].[s_id], [omni_TestCountAffectedRowsShapes].[dbo].[big].[v]
: float64 500
: float64 0.021643519401550293
: float64 0.011157000437378883
: int64 19
: float64 0.03280051797628403
: string
| [omni_TestCountAffectedRowsShapes].[dbo].[big].[id], [omni_TestCountAffectedRowsShapes].[dbo].[big].[s_id], [omni_TestCountAffectedRowsShapes].[dbo].[big].[v]
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
| INSERT INTO dbo.big_archive SELECT b.id, b.s_id, b.v FROM dbo.big AS b JOIN dbo.s AS s ON b.s_id = s.id WHERE s.flag = 1;
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
| INSERT INTO dbo.big_archive SELECT b.id, b.s_id, b.v FROM dbo.big AS b JOIN dbo.s AS s ON b.s_id = s.id WHERE s.flag = 1;
: int64 1
: int64 1
: int64 0
: null
: null
: string
| 1
: null
: float64 1000
: null
: null
: null
: float64 0.14742934703826904
: null
: null
: string
| INSERT
: bool false
: null
< row
: string
|   |--Clustered Index Insert(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[PK__big_arch__standin4]), SET:([omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[id] = [omni_TestCountAffectedRowsShapes].[dbo].[big].[id] as [b].[id],[omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[s_id] = [omni_TestCountAffectedRowsShapes].[dbo].[big].[s_id] as [b].[s_id],[omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[v] = [omni_TestCountAffectedRowsShapes].[dbo].[big].[v] as [b].[v]))
: int64 1
: int64 2
: int64 1
: string
| Clustered Index Insert
: string
| Insert
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[PK__big_arch__standin4]), SET:([omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[id] = [omni_TestCountAffectedRowsShapes].[dbo].[big].[id] as [b].[id],[omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[s_id] = [omni_TestCountAffectedRowsShapes].[dbo].[big].[s_id] as [b].[s_id],[omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[v] = [omni_TestCountAffectedRowsShapes].[dbo].[big].[v] as [b].[v])
: null
: float64 1000
: float64 0.01891443133354187
: float64 0.0010000000474974513
: int64 9
: float64 0.14742934703826904
: null
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|        |--Sort(ORDER BY:([b].[id] ASC))
: int64 1
: int64 3
: int64 2
: string
| Sort
: string
| Sort
: string
| ORDER BY:([b].[id] ASC)
: null
: float64 1000
: float64 0.011261261068284512
: float64 0.015646642073988914
: int64 19
: float64 0.12751492857933044
: string
| [b].[id], [b].[s_id], [b].[v]
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|             |--Hash Match(Inner Join, HASH:([s].[id])=([b].[s_id]))
: int64 1
: int64 4
: int64 3
: string
| Hash Match
: string
| Inner Join
: string
| HASH:([s].[id])=([b].[s_id])
: null
: float64 1000
: float64 0
: float64 0.0643635019659996
: int64 19
: float64 0.1006070151925087
: string
| [b].[id], [b].[s_id], [b].[v]
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|                  |--Clustered Index Scan(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[s].[PK__s__standin3] AS [s]), WHERE:([omni_TestCountAffectedRowsShapes].[dbo].[s].[flag] as [s].[flag]=(1)))
: int64 1
: int64 5
: int64 4
: string
| Clustered Index Scan
: string
| Clustered Index Scan
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[s].[PK__s__standin3] AS [s]), WHERE:([omni_TestCountAffectedRowsShapes].[dbo].[s].[flag] as [s].[flag]=(1))
: string
| [s].[id]
: float64 10
: float64 0.0031250000465661287
: float64 0.00026699999580159783
: int64 12
: float64 0.0033919999841600657
: string
| [s].[id]
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|                  |--Clustered Index Scan(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[PK__big__standin2] AS [b]))
: int64 1
: int64 6
: int64 4
: string
| Clustered Index Scan
: string
| Clustered Index Scan
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[PK__big__standin2] AS [b])
: string
| [b].[id], [b].[s_id], [b].[v]
: float64 10000
: float64 0.021643519401550293
: float64 0.011157000437378883
: int64 19
: float64 0.03280051797628403
: string
| [b].[id], [b].[s_id], [b].[v]
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
| MERGE dbo.big_archive AS tgt
| USING (SELECT b.id, b.s_id, b.v FROM dbo.big AS b JOIN dbo.s AS s ON b.s_id = s.id WHERE s.flag = 1) AS src
| ON tgt.id = src.id
| WHEN NOT MATCHED THEN INSERT (id, s_id, v) VALUES (src.id, src.s_id, src.v);
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
| MERGE dbo.big_archive AS tgt
| USING (SELECT b.id, b.s_id, b.v FROM dbo.big AS b JOIN dbo.s AS s ON b.s_id = s.id WHERE s.flag = 1) AS src
| ON tgt.id = src.id
| WHEN NOT MATCHED THEN INSERT (id, s_id, v) VALUES (src.id, src.s_id, src.v);
: int64 1
: int64 1
: int64 0
: null
: null
: string
| 1
: null
: float64 899.1900024414062
: null
: null
: null
: float64 6.216835021972656
: null
: null
: string
| MERGE
: bool false
: null
< row
: string
|   |--Clustered Index Merge(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[PK__big_arch__standin4] AS [tgt]), SET:([omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[v] as [tgt].[v] = [omni_TestCountAffectedRowsShapes].[dbo].[big].[v] as [b].[v],[omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[s_id] as [tgt].[s_id] = [omni_TestCountAffectedRowsShapes].[dbo].[big].[s_id] as [b].[s_id],[omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[id] as [tgt].[id] = [omni_TestCountAffectedRowsShapes].[dbo].[big].[id] as [b].[id]) ACTION:([Action1005]))
: int64 1
: int64 3
: int64 1
: string
| Clustered Index Merge
: string
| Merge
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[PK__big_arch__standin4] AS [tgt]), SET:([omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[v] as [tgt].[v] = [omni_TestCountAffectedRowsShapes].[dbo].[big].[v] as [b].[v],[omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[s_id] as [tgt].[s_id] = [omni_TestCountAffectedRowsShapes].[dbo].[big].[s_id] as [b].[s_id],[omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[id] as [tgt].[id] = [omni_TestCountAffectedRowsShapes].[dbo].[big].[id] as [b].[id]) ACTION:([Action1005])
: null
: float64 899.1900024414062
: float64 6.078851222991943
: float64 0.0008991899667307734
: int64 9
: float64 6.216835021972656
: null
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|        |--Filter(WHERE:([Action1005] IS NOT NULL))
: int64 1
: int64 4
: int64 3
: string
| Filter
: string
| Filter
: string
| WHERE:([Action1005] IS NOT NULL)
: null
: float64 899.1900024414062
: float64 0
: float64 0.0004795679997187108
: int64 27
: float64 0.1370847076177597
: string
| [b].[id], [b].[s_id], [b].[v], [tgt].[id], [Action1005]
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|             |--Compute Scalar(DEFINE:([Action1005]=ForceOrder(CASE WHEN [TrgPrb1003] IS NOT NULL THEN NULL ELSE (4) END)))
: int64 1
: int64 5
: int64 4
: string
| Compute Scalar
: string
| Compute Scalar
: string
| DEFINE:([Action1005]=ForceOrder(CASE WHEN [TrgPrb1003] IS NOT NULL THEN NULL ELSE (4) END))
: string
| [Action1005]=ForceOrder(CASE WHEN [TrgPrb1003] IS NOT NULL THEN NULL ELSE (4) END)
: float64 999.0999755859375
: float64 0
: float64 9.991000115405768e-05
: int64 27
: float64 0.1366051435470581
: string
| [b].[id], [b].[s_id], [b].[v], [tgt].[id], [Action1005]
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|                  |--Merge Join(Right Outer Join, MERGE:([tgt].[id])=([b].[id]), RESIDUAL:([omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[id] as [tgt].[id]=[omni_TestCountAffectedRowsShapes].[dbo].[big].[id] as [b].[id]))
: int64 1
: int64 6
: int64 5
: string
| Merge Join
: string
| Right Outer Join
: string
| MERGE:([tgt].[id])=([b].[id]), RESIDUAL:([omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[id] as [tgt].[id]=[omni_TestCountAffectedRowsShapes].[dbo].[big].[id] as [b].[id])
: null
: float64 999.0999755859375
: float64 0
: float64 0.005704110022634268
: int64 27
: float64 0.13650523126125336
: string
| [b].[id], [b].[s_id], [b].[v], [tgt].[id], [TrgPrb1003]
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|                       |--Compute Scalar(DEFINE:([TrgPrb1003]=(1)))
: int64 1
: int64 7
: int64 6
: string
| Compute Scalar
: string
| Compute Scalar
: string
| DEFINE:([TrgPrb1003]=(1))
: string
| [TrgPrb1003]=(1)
: float64 1
: float64 0
: float64 1.0000000116860974e-07
: int64 15
: float64 0.0032832000870257616
: string
| [tgt].[id], [TrgPrb1003]
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|                       |    |--Clustered Index Scan(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[PK__big_arch__standin4] AS [tgt]), ORDERED FORWARD)
: int64 1
: int64 8
: int64 7
: string
| Clustered Index Scan
: string
| Clustered Index Scan
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big_archive].[PK__big_arch__standin4] AS [tgt]), ORDERED FORWARD
: string
| [tgt].[id]
: float64 1
: float64 0.0031250000465661287
: float64 0.00015809999604243785
: int64 11
: float64 0.0032830999698489904
: string
| [tgt].[id]
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|                       |--Sort(ORDER BY:([b].[id] ASC))
: int64 1
: int64 9
: int64 6
: string
| Sort
: string
| Sort
: string
| ORDER BY:([b].[id] ASC)
: null
: float64 1000
: float64 0.011261261068284512
: float64 0.015646642073988914
: int64 19
: float64 0.12751492857933044
: string
| [b].[id], [b].[s_id], [b].[v]
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|                            |--Hash Match(Inner Join, HASH:([s].[id])=([b].[s_id]))
: int64 1
: int64 10
: int64 9
: string
| Hash Match
: string
| Inner Join
: string
| HASH:([s].[id])=([b].[s_id])
: null
: float64 1000
: float64 0
: float64 0.0643635019659996
: int64 19
: float64 0.1006070151925087
: string
| [b].[id], [b].[s_id], [b].[v]
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|                                 |--Clustered Index Scan(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[s].[PK__s__standin3] AS [s]), WHERE:([omni_TestCountAffectedRowsShapes].[dbo].[s].[flag] as [s].[flag]=(1)))
: int64 1
: int64 11
: int64 10
: string
| Clustered Index Scan
: string
| Clustered Index Scan
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[s].[PK__s__standin3] AS [s]), WHERE:([omni_TestCountAffectedRowsShapes].[dbo].[s].[flag] as [s].[flag]=(1))
: string
| [s].[id]
: float64 10
: float64 0.0031250000465661287
: float64 0.00026699999580159783
: int64 12
: float64 0.0033919999841600657
: string
| [s].[id]
: null
: string
| PLAN_ROW
: bool false
: float64 1
< row
: string
|                                 |--Clustered Index Scan(OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[PK__big__standin2] AS [b]))
: int64 1
: int64 12
: int64 10
: string
| Clustered Index Scan
: string
| Clustered Index Scan
: string
| OBJECT:([omni_TestCountAffectedRowsShapes].[dbo].[big].[PK__big__standin2] AS [b])
: string
| [b].[id], [b].[s_id], [b].[v]
: float64 10000
: float64 0.021643519401550293
: float64 0.011157000437378883
: int64 19
: float64 0.03280051797628403
: string
| [b].[id], [b].[s_id], [b].[v]
: null
: string
| PLAN_ROW
: bool false
: float64 1
> exec
| SET SHOWPLAN_ALL OFF;
< ok 0 -
