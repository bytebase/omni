-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestQuery default#1
> exec
| CREATE TABLE q_items (id INT PRIMARY KEY, note NVARCHAR(100), bin VARBINARY(16), flag BIT, amount DECIMAL(10,2), uid UNIQUEIDENTIFIER, created DATETIME2(3), offset_at DATETIMEOFFSET(0));
| INSERT INTO q_items VALUES (1, N'a', 0x01, 1, 1.50, 'A0EEBC99-9C0B-4EF8-BB6D-6BB9BD380A11', '2026-03-04 05:06:07.891', '2026-03-04 05:06:07.891 +02:00');
| INSERT INTO q_items VALUES (2, N'b', NULL, 0, 2.50, NULL, NULL, NULL);
| INSERT INTO q_items VALUES (3, N'c', 0x, 1, 3.50, NULL, NULL, NULL);
< ok 3 -
> exec
| CREATE PROCEDURE q_sp_items AS
| BEGIN
| 	SELECT id, note FROM q_items WHERE id = 1;
| 	SELECT COUNT(*) AS n FROM q_items;
| END;
< ok 0 -
> query
| SELECT
| 		42 AS i,
| 		3.5 AS f,
| 		N'hello' AS s,
| 		CAST(1 AS BIT) AS b,
| 		CAST(0 AS BIT) AS b0,
| 		CAST(12.3400 AS DECIMAL(10,4)) AS num,
| 		CAST('A0EEBC99-9C0B-4EF8-BB6D-6BB9BD380A11' AS UNIQUEIDENTIFIER) AS uid,
| 		CAST(0xDEADBEEF AS VARBINARY(4)) AS bin,
| 		CAST(0x AS VARBINARY(4)) AS empty_bin,
| 		CAST('12:34:56' AS TIME) AS t,
| 		CAST('2026-03-04' AS DATE) AS dt,
| 		CAST('2026-03-04 05:06:07.891' AS DATETIME2(3)) AS ts,
| 		CAST('2026-03-04 05:06:07 +02:00' AS DATETIMEOFFSET(0)) AS tstz,
| 		NULL AS nul;
< result
< column "i" "INT"
< column "f" "DECIMAL"
< column "s" "NVARCHAR"
< column "b" "BIT"
< column "b0" "BIT"
< column "num" "DECIMAL"
< column "uid" "UNIQUEIDENTIFIER"
< column "bin" "VARBINARY"
< column "empty_bin" "VARBINARY"
< column "t" "TIME"
< column "dt" "DATE"
< column "ts" "DATETIME2"
< column "tstz" "DATETIMEOFFSET"
< column "nul" "INT"
< row
: int64 42
: bytes
| 3.5
: string
| hello
: bool true
: bool false
: bytes
| 12.3400
: bytes
|q "\x99\xbc\xee\xa0\v\x9c\xf8N\xbbmk\xb9\xbd8\n\x11"
: bytes
|q "ޭ\xbe\xef"
: bytes
|
: time 0001-01-01T12:34:56Z
: time 2026-03-04T00:00:00Z
: time 2026-03-04T05:06:07.891Z
: time 2026-03-04T05:06:07+02:00
: null
> query
| SELECT id FROM q_items ORDER BY id OFFSET 0 ROWS FETCH NEXT 2 ROWS ONLY;
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
> exec
| INSERT INTO q_items (id, note, flag) VALUES (4, N'd', 1);
< ok 1 -
> query
| EXEC q_sp_items;
< result
< column "id" "INT"
< column "note" "NVARCHAR"
< row
: int64 1
: string
| a
< result
< column "n" "INT"
< row
: int64 4
> query ^
< result
< column "id" "INT"
< column "note" "NVARCHAR"
< row
: int64 1
: string
| a
< result
< column "n" "INT"
< row
: int64 4
> exec
| SET SHOWPLAN_ALL ON
< ok 0 -
> query
| SELECT note FROM q_items WHERE id > 1;
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
| SELECT note FROM q_items WHERE id > 1;
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
: float64 0.003285299986600876
: null
: null
: string
| SELECT
: bool false
: null
< row
: string
|   |--Clustered Index Seek(OBJECT:([omni_TestQuery].[dbo].[q_items].[PK__q_items__standin1]), SEEK:([omni_TestQuery].[dbo].[q_items].[id] > CONVERT_IMPLICIT(int,[@1],0)) ORDERED FORWARD)
: int64 1
: int64 2
: int64 1
: string
| Clustered Index Seek
: string
| Clustered Index Seek
: string
| OBJECT:([omni_TestQuery].[dbo].[q_items].[PK__q_items__standin1]), SEEK:([omni_TestQuery].[dbo].[q_items].[id] > CONVERT_IMPLICIT(int,[@1],0)) ORDERED FORWARD
: string
| [omni_TestQuery].[dbo].[q_items].[note]
: float64 3
: float64 0.0031250000465661287
: float64 0.00016029999824240804
: int64 111
: float64 0.003285299986600876
: string
| [omni_TestQuery].[dbo].[q_items].[note]
: null
: string
| PLAN_ROW
: bool false
: float64 1
> exec
| SET SHOWPLAN_ALL OFF
< ok 0 -
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT note FROM q_items WHERE id > 1;
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT note FROM q_items WHERE id &gt; 1" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032853" StatementEstRows="3" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x050D2D2347709AB2" QueryPlanHash="0x707ACD586232ED7C" CardinalityEstimationModelVersion="160" ParameterizedText="(@1 tinyint)SELECT [note] FROM [q_items] WHERE [id]&gt;@1"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="120"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="131420" EstimatedPagesCached="164275" EstimatedAvailableDegreeOfParallelism="10" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><OptimizerStatsUsage><StatisticsInfo LastUpdate="2026-01-01T00:00:00" ModificationCount="0" SamplingPercent="100" Statistics="[PK__q_items__standin1]" Table="[q_items]" Schema="[dbo]" Database="[omni_TestQuery]"></StatisticsInfo></OptimizerStatsUsage><RelOp NodeId="0" PhysicalOp="Clustered Index Seek" LogicalOp="Clustered Index Seek" EstimateRows="3" EstimatedRowsRead="3" EstimateIO="0.003125" EstimateCPU="0.0001603" AvgRowSize="111" EstimatedTotalSubtreeCost="0.0032853" TableCardinality="4" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[omni_TestQuery]" Schema="[dbo]" Table="[q_items]" Column="note"></ColumnReference></OutputList><IndexScan Ordered="1" ScanDirection="FORWARD" ForcedIndex="0" ForceSeek="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[omni_TestQuery]" Schema="[dbo]" Table="[q_items]" Column="note"></ColumnReference></DefinedValue></DefinedValues><Object Database="[omni_TestQuery]" Schema="[dbo]" Table="[q_items]" Index="[PK__q_items__standin1]" IndexKind="Clustered" Storage="RowStore"></Object><SeekPredicates><SeekPredicateNew><SeekKeys><StartRange ScanType="GT"><RangeColumns><ColumnReference Database="[omni_TestQuery]" Schema="[dbo]" Table="[q_items]" Column="id"></ColumnReference></RangeColumns><RangeExpressions><ScalarOperator ScalarString="CONVERT_IMPLICIT(int,[@1],0)"><Identifier><ColumnReference Column="ConstExpr1002"><ScalarOperator><Convert DataType="int" Style="0" Implicit="1"><ScalarOperator><Identifier><ColumnReference Column="@1"></ColumnReference></Identifier></ScalarOperator></Convert></ScalarOperator></ColumnReference></Identifier></ScalarOperator></RangeExpressions></StartRange></SeekKeys></SeekPredicateNew></SeekPredicates></IndexScan></RelOp><ParameterList><ColumnReference Column="@1" ParameterDataType="tinyint" ParameterCompiledValue="(1)"></ColumnReference></ParameterList></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> query
| DECLARE @id int = 1;
|
| SELECT @id AS id;
|
: message
< stream
< message nextresultset more
< message next
< message rowsaffected 1
< message nextresultset end
< result
< column "id" "INT"
< row
: int64 1
> query
| SELECT * FROM q_missing_table;
< error
| 208 1 16 1 Invalid object name 'q_missing_table'.
> query ^
< error
| 208 1 16 1 Invalid object name 'q_missing_table'.
> query
| SELECT COUNT(*) FROM q_items WHERE id = 1;
< result
< column "" "INT"
< row
: int64 1
> query
| CREATE TABLE q_scratch (id INT);
|
| INSERT INTO q_scratch VALUES (1);
|
: message
< stream
< message nextresultset more
< message rowsaffected 1
< message nextresultset more
< message nextresultset end
> exec
| CREATE TABLE q_archive (id INT, note NVARCHAR(100));
< ok 0 -
> exec
| SELECT id, note INTO q_archive2 FROM q_items;
< ok 4 -
> query
| SELECT COUNT(*) FROM q_archive2;
< result
< column "" "INT"
< row
: int64 4
> query
| SELECT id, REPLICATE(N'x', 200) AS payload FROM q_items ORDER BY id;
|
| SELECT 7 AS after_cap;
|
: message
< stream
< message next
< message rowsaffected 4
< message nextresultset more
< message next
< message rowsaffected 1
< message nextresultset end
< result
< column "id" "INT"
< column "payload" "NVARCHAR"
< row
: int64 1
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 2
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 3
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 4
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< result
< column "after_cap" "INT"
< row
: int64 7
> query
| SELECT 1 AS a;
< result
< column "a" "INT"
< row
: int64 1
> query
| SELECT 2 AS b;
< result
< column "b" "INT"
< row
: int64 2
> query
| INSERT INTO q_items (id, note, flag) VALUES (5, N'e', 1);
|
| SELECT 1 AS after_insert;
|
: message
< stream
< message rowsaffected 1
< message nextresultset more
< message next
< message rowsaffected 1
< message nextresultset end
< result
< column "after_insert" "INT"
< row
: int64 1
