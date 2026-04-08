# Snowflake Migration DAG

## Nodes

| Node | Package | Depends On | Can Parallelize With | Priority | Status |
|------|---------|------------|----------------------|----------|--------|
| ast | snowflake/ast | (none) | - | P0 | not started |
| parser | snowflake/parser | ast | - | P0 | not started |
| analysis | snowflake/analysis | parser | catalog | P0 | not started |
| catalog | snowflake/catalog | ast | analysis | P0 | not started |
| completion | snowflake/completion | parser, catalog | - | P0 | not started |
| semantic | snowflake/semantic | parser, catalog | quality | P1 | not started |
| quality | snowflake/quality | parser, analysis | semantic | P1 | not started |
| deparse | snowflake/deparse | ast | analysis, catalog | P1 | not started |

Priority: P0 = consumed by bytebase (critical path), P1 = legacy parser parity.

## Execution Order

```
Phase 1: ast (blocking — all packages depend on AST node types)
Phase 2: parser (blocking — most packages depend on parse tree)
Phase 3: analysis + catalog (parallel)
Phase 4: completion (after parser + catalog)
Phase 5: semantic + quality + deparse (parallel, P1)
```

### Dependency Graph

```
ast
 ├── parser
 │     ├── analysis ──────────┐
 │     │     └── quality      │
 │     ├── completion ◄───────┤
 │     ├── semantic ◄─────────┤
 │     └── deparse            │
 └── catalog ─────────────────┘
```

## Node Scope

### ast (P0)

AST node types for all Snowflake statement types. Based on analysis, this covers:

**DDL nodes:**
- CreateTable, AlterTable, DropTable (+ variants: AS SELECT, LIKE, CLONE)
- CreateView, AlterView, DropView, CreateMaterializedView, AlterMaterializedView, DropMaterializedView
- CreateDatabase, AlterDatabase, DropDatabase
- CreateSchema, AlterSchema, DropSchema
- CreateDynamicTable, AlterDynamicTable, DropDynamicTable
- CreateExternalTable, AlterExternalTable, DropExternalTable
- CreateEventTable
- CreateFunction, AlterFunction, DropFunction
- CreateProcedure, AlterProcedure, DropProcedure
- CreateExternalFunction, DropExternalFunction
- CreateStage, AlterStage, DropStage
- CreateFileFormat, AlterFileFormat, DropFileFormat
- CreatePipe, AlterPipe, DropPipe
- CreateStream, AlterStream, DropStream
- CreateSequence, AlterSequence, DropSequence
- CreateTag, AlterTag, DropTag
- CreateTask, AlterTask, DropTask
- CreateApiIntegration, AlterApiIntegration, DropApiIntegration (+ Security, Notification, Storage integrations)
- CreateMaskingPolicy, AlterMaskingPolicy, DropMaskingPolicy (+ RowAccess, Session, Password, Network policies)
- CreateRole, AlterRole, DropRole
- CreateUser, AlterUser, DropUser
- CreateWarehouse, AlterWarehouse, DropWarehouse (from ALTER ACCOUNT / resource monitor context)
- CreateAlert, AlterAlert, DropAlert
- CreateFailoverGroup, AlterFailoverGroup, DropFailoverGroup (+ ReplicationGroup)
- CreateGitRepository, AlterGitRepository, DropGitRepository
- CreateShare, AlterShare, DropShare
- CreateManagedAccount, AlterManagedAccount, DropManagedAccount
- CreateConnection, DropConnection
- CreateResourceMonitor, AlterResourceMonitor, DropResourceMonitor
- CreateSemanticView, DropSemanticView
- CreateDataset, DropDataset
- CreateAccount, AlterAccount
- UndropDatabase, UndropSchema, UndropTable, UndropDynamicTable, UndropTag
- Truncate (TABLE, MATERIALIZED VIEW)
- Comment

**DML nodes:**
- Select (full query model: CTEs, JOINs, set ops, PIVOT/UNPIVOT, MATCH_RECOGNIZE, time travel, window functions, SAMPLE, QUALIFY)
- Insert (standard, OVERWRITE, multi-table INSERT)
- Update
- Delete
- Merge
- CopyInto (TABLE and LOCATION variants)

**DCL nodes:**
- Grant, Revoke (all variants: object, schema, future, role, share)

**Utility nodes:**
- Show (60+ variants — likely a single Show node with variant enum)
- Describe (27 variants — similar pattern)
- Use (DATABASE, SCHEMA, ROLE, WAREHOUSE, SECONDARY ROLES)
- Explain, Call, Execute (IMMEDIATE, TASK)
- Set, Unset
- Begin, Commit, Rollback
- Get, Put, List, Remove (stage operations)

**Expression/clause nodes:**
- Expression (arithmetic, comparison, logical, string, NULL, BETWEEN, IN, EXISTS, CASE, CAST, TRY_CAST, IFF, COALESCE, NULLIF, subquery)
- DataType (30+ types including VARIANT, OBJECT, ARRAY, GEOGRAPHY, GEOMETRY, VECTOR)
- ColumnDef, Constraint (PK, FK, UNIQUE, CHECK, NOT NULL)
- TableRef, JoinClause, FromClause, WhereClause, GroupByClause, HavingClause, QualifyClause, OrderByClause, LimitClause
- WindowSpec, WindowFrame
- ObjectName (database.schema.object — Snowflake's 3-part naming)

### parser (P0)

Hand-written recursive descent lexer + parser producing the AST nodes above. Must cover:

**Lexer:**
- Case-insensitive keywords (Snowflake behavior)
- 1000+ keywords from legacy SnowflakeLexer.g4
- Identifiers (unquoted, double-quoted with case sensitivity)
- String literals (single-quoted)
- Numeric literals (integer, decimal, float, scientific notation)
- Operators and punctuation
- Comments (single-line --, block /* */)
- Semi-structured access operators (colon, bracket, dot notation for VARIANT)

**Parser — statement-level:**
- All DDL, DML, DCL, utility statements from the AST scope
- Multi-statement parsing (semicolon-separated)
- Error recovery for best-effort parsing

**Critical for bytebase consumption:**
- Statement splitting (lexer-level, semicolon-based)
- Statement type classification (DDL/DML/DQL)
- Parse tree that supports listener/walker patterns for query span extraction and SQL review

### analysis (P0)

Statement analysis for bytebase consumption. Maps to these bytebase registrations:
- `ParseStatementsFunc` — parse + classify statements
- `QueryValidator` — validate if statement is executable, detect EXECUTE
- `QueryType` — DDL/DML/DQL classification
- `QuerySpan` — table/column dependency extraction for data masking
- `DiagnoseFunc` — syntax error reporting with positions
- `StatementRangesFunc` — statement position ranges (UTF-16)

### catalog (P0)

Schema metadata types for Snowflake. Needed by completion and semantic analysis. Covers:
- Database, Schema, Table, View, MaterializedView, DynamicTable, ExternalTable
- Column (with data types, constraints, defaults)
- Function/Procedure signatures
- Stage, FileFormat, Pipe, Stream, Sequence
- Tag, Policy (masking, row access, etc.)
- Role, User, Warehouse

### completion (P0)

Auto-complete candidates based on cursor position. No legacy Snowflake completion exists in bytebase — this is new functionality. Needs parser (for partial parse) + catalog (for schema-aware suggestions).

### semantic (P1)

Semantic validation beyond syntax: type checking, constraint validation, object existence. Not currently consumed by bytebase for Snowflake but needed for parity.

### quality (P1)

SQL lint/quality rules. Maps to the 15 bytebase advisor rules. In the omni architecture these move from bytebase into the omni quality package. Can be implemented incrementally per rule.

### deparse (P1)

AST to SQL string generation. Useful for query rewriting (e.g., limit injection). The legacy parser uses ANTLR token stream rewriting; omni will use AST-based deparsing.

## Notes

- **No AST layer in legacy**: The legacy parser exposes raw ANTLR4 parse tree. The omni migration introduces a proper AST for the first time, which is a significant design task.
- **Schema sync does not use parser**: Bytebase's Snowflake schema sync uses `GET_DDL()` DB function, not the parser. No parser dependency there.
- **Snowflake case sensitivity**: Unquoted identifiers are uppercased by Snowflake. The parser/AST must handle this correctly.
- **Semi-structured data access**: Snowflake's VARIANT/OBJECT/ARRAY types use colon, bracket, and dot notation for field access. This is a distinctive Snowflake feature that needs careful lexer/parser support.
- **3-part naming**: Snowflake uses database.schema.object naming throughout. ObjectName handling is critical.
- **Large grammar**: The legacy grammar is 141KB (parser) + 56KB (lexer). The hand-written parser will be substantial.
- **Completion is new**: No legacy completion exists — this is greenfield work in omni.
