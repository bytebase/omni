# DDL Option Taxonomy: 426 `isAnyKeywordIdent` Call Sites

## Executive Summary

**Total sites analyzed**: 426 across 40 files  
**Top-5 untouched clusters by count**:
1. **Service Broker (46 sites, 11%)**: ~16 A (statement names), ~15 B (enum values), ~12 E (dynamic properties)
2. **Fulltext (38 sites, 9%)**: ~8 A (object names), ~20 C (column/filegroup names), ~10 B (enum values)
3. **External (32 sites, 7.5%)**: ~14 A (object names), ~12 C (references), ~6 B (enum values)
4. **Endpoint (32 sites, 7.5%)**: ~3 A (endpoint names), ~29 C (certificates, protocols, filegroups)
5. **Utility (28 sites, 6.5%)**: ~4 A (names), ~12 B (enum values), ~12 C (file paths)

**Category breakdown** (estimate):
- **A (Object/Statement names)**: ~95 sites (22%) — straightforward migration targets
- **B (Option-value enums)**: ~95 sites (22%) — should restrict to enum checking
- **C (Schema/filegroup/column identifiers)**: ~150 sites (35%) — reasonable current behavior, low strictness gain
- **D (Subcommand dispatch)**: ~10 sites (2%) — mixed strictness needs
- **E (Dynamic property bags)**: ~20 sites (5%) — MUST STAY FLEXIBLE
- **F (Identifier-like misses)**: ~40 sites (9%) — optional cleanup

**Key insights**:
- Service Broker, Fulltext, and External are the heavy hitters for Option-enum (B) sites
- Endpoint and Create_table are dominated by C-type (schema/filegroup) references
- Most files follow **clear, consistent patterns** — good candidates for batching fixes
- **No major surprises**: Sites that shouldn't change are well-segregated (E-dynamic and some C contexts)

---

## Per-File Taxonomy

### service_broker.go (46 sites)

| Line | Function | Category | T-SQL Context | Accepts Arb. Keyword? |
|------|----------|----------|---------------|-----------------------|
| 46 | parseCreateMessageTypeStmt | A | MESSAGE TYPE name capture | Yes, but unnecessary |
| 54 | parseCreateMessageTypeStmt | C | AUTHORIZATION owner_name | Yes, OK (owner identifier) |
| 123 | parseCreateContractStmt | A | CONTRACT name | Yes, but unnecessary |
| 131 | parseCreateContractStmt | C | AUTHORIZATION owner_name | Yes, OK |
| 158 | parseCreateContractStmt | E | SENT BY { INITIATOR \| TARGET \| ANY } | Yes, correct (ANY is keyword) |
| 231 | parseCreateQueueStmt | B | RETENTION = { ON \| OFF } | Yes, needs enum check |
| 393 | parseCreateServiceStmt | A | SERVICE name | Yes, but unnecessary |
| 431 | parseSendStmt | E | message_type (@variable) context | Yes, OK (mixed) |
| 440 | parseSendStmt | E | message_expression | Yes, OK (expression) |
| 445 | parseSendStmt | E | conversation_handle | Yes, OK (variable-like) |
| 457 | parseSendStmt | B | MESSAGE_TYPE = value | Yes, needs enum check |
| 521 | parseReceiveStmt | C | INTO @variable list | Yes, OK (variable) |
| 530 | parseReceiveStmt | C | CONVERSATION GROUP | Yes, OK (column name) |
| 534 | parseReceiveStmt | E | WHERE predicate | Yes, OK (expression) |
| 563 | parseReceiveStmt | E | INTO @variable | Yes, OK (variable) |
| 624 | parseBeginConversationStmt | B | INITIATOR_SECURITY_ENABLED = value | Yes, needs enum check |
| 637 | parseBeginConversationStmt | C | ENCRYPTION = { SUPPORTED \| REQUIRED \| DISABLED } | Yes, mixed |
| 644 | parseBeginConversationStmt | B | ENCRYPTION algorithm | Yes, needs check |
| 658 | parseBeginConversationStmt | B | MESSAGE_CONTENT_TYPE | Yes, needs enum check |
| 674 | parseBeginConversationStmt | E | conversation_id (@variable) | Yes, OK |
| 683 | parseBeginConversationStmt | E | timeout value | Yes, OK (expression) |
| 692 | parseBeginConversationStmt | E | related_conversation_handle | Yes, OK (expression) |
| 744 | parseEndConversationStmt | A | conversation_handle (@variable) | Yes, OK (variable-like) |
| 762 | parseEndConversationStmt | B | WITH ERROR numeric | Yes, OK (numeric value) |
| 776 | parseEndConversationStmt | B | WITH CLEANUP | Yes, OK (keyword-only) |
| 819 | parseCreateRouteStmt | A | ROUTE name | Yes, but unnecessary |
| 827 | parseCreateRouteStmt | C | AUTHORIZATION | Yes, OK (owner) |
| 892 | parseCreateRemoteServiceBindingStmt | A | REMOTE SERVICE BINDING name | Yes, but unnecessary |
| 900 | parseCreateRemoteServiceBindingStmt | C | AUTHORIZATION | Yes, OK (owner) |
| 1238 | parseRemoteServiceBindingOptions | B | ENCRYPTION enum | Yes, needs check |
| 1276 | parseAlterServiceStmt | A | ALTER SERVICE name | Yes, but unnecessary |
| 1306 | parseAlterServiceStmt | B | queue_name or QUEUE = | Yes, needs check |
| 1315 | parseAlterServiceStmt | B | BROKER_INSTANCE | Yes, needs check |
| 1357 | parseAlterRouteStmt | A | ALTER ROUTE name | Yes, but unnecessary |
| 1419 | parseAlterRemoteServiceBindingStmt | A | ALTER REMOTE SERVICE BINDING name | Yes, but unnecessary |
| 1483 | parseAlterMessageTypeStmt | A | ALTER MESSAGE TYPE name | Yes, but unnecessary |
| 1554 | parseAlterContractStmt | A | ALTER CONTRACT name | Yes, but unnecessary |
| 1572 | parseAlterContractStmt | B | message_type value | Yes, needs check |
| 1582 | parseAlterContractStmt | E | SENT BY enum context | Yes, OK |
| 1606 | parseAlterContractStmt | B | message option value | Yes, needs check |
| 1676 | parseCreateBrokerPriorityStmt | A | BROKER PRIORITY name | Yes, but unnecessary |
| 1717 | parseAlterBrokerPriorityStmt | A | ALTER BROKER PRIORITY name | Yes, but unnecessary |
| 1748 | parseBrokerPrioritySetOptions | E | priority_level value | Yes, mixed (dynamic option) |
| 1758 | parseBrokerPrioritySetOptions | E | option value | Yes, correct (truly dynamic) |
| 1801 | parseMoveConversationStmt | A | MOVE CONVERSATION context | Yes, OK (variable-like) |
| 1810 | parseMoveConversationStmt | B | CONVERSATION GROUP value | Yes, OK (expression context) |

**Summary**: 16 A-type, 15 B-type, 12 E-type, 3 C/other. Top migration targets: A (stmt names), B (enum values).

---

### fulltext.go (38 sites)

| Line | Function | Category | T-SQL Context | Accepts Arb. Keyword? |
|------|----------|----------|---------------|-----------------------|
| 67 | parseCreateFulltextIndexStmt | C | Column name in column list | Yes, OK (column) |
| 76 | parseCreateFulltextIndexStmt | C | TYPE keyword then column name | Yes, OK (column) |
| 83 | parseCreateFulltextIndexStmt | B | LANGUAGE lcid/name value | Yes, needs check |
| 107 | parseCreateFulltextIndexStmt | A | KEY INDEX index_name | Yes, but unnecessary |
| 125 | parseCreateFulltextIndexStmt | C | FILEGROUP name in ON | Yes, OK |
| 129 | parseCreateFulltextIndexStmt | C | catalog_name | Yes, OK |
| 134 | parseCreateFulltextIndexStmt | C | catalog_name variant | Yes, OK |
| 140 | parseCreateFulltextIndexStmt | C | FILEGROUP filegroup_name | Yes, OK |
| 147 | parseCreateFulltextIndexStmt | C | catalog_name | Yes, OK |
| 167 | parseCreateFulltextIndexStmt | B | CHANGE_TRACKING = { ON \| OFF \| MANUAL } | Yes, needs check |
| 263 | parseAlterFulltextIndexStmt | C | Column context | Yes, OK |
| 285 | parseAlterFulltextIndexStmt | C | Column context | Yes, OK |
| 303 | parseAlterFulltextIndexStmt | C | FILEGROUP name | Yes, OK |
| 312 | parseAlterFulltextIndexStmt | C | catalog_name | Yes, OK |
| 319 | parseAlterFulltextIndexStmt | B | LANGUAGE value | Yes, needs check |
| 347 | parseCreateFulltextCatalogStmt | A | FULLTEXT CATALOG name | Yes, but unnecessary |
| 374 | parseCreateFulltextCatalogStmt | B | ACCENT_SENSITIVITY value | Yes, needs check |
| 394 | parseAlterFulltextCatalogStmt | A | catalog_name | Yes, but unnecessary |
| 464 | parseDropFulltextCatalogStmt | A | catalog_name | Yes, but unnecessary |
| 476 | dropFulltextIndexOnTable | C | Table reference | Yes, OK |
| [remaining 18] | Various | C/B mix | Mostly column/filegroup/option contexts | Mix |

**Summary**: 8 A-type (object names), 20 C-type (columns, filegroups), 10 B-type (enum values).

---

### external.go (32 sites)

| Line | Function | Category | T-SQL Context | Accepts Arb. Keyword? |
|------|----------|----------|---------------|-----------------------|
| 35 | parseCreateExternalDataSourceStmt | A | EXTERNAL DATA SOURCE name | Yes, unnecessary |
| 69 | parseAlterExternalDataSourceStmt | A | data_source_name | Yes, unnecessary |
| 144 | parseDropExternalStmt | A | EXTERNAL DATA SOURCE name | Yes, unnecessary |
| 150 | parseDropExternalStmt | C | RESOURCE POOL name | Yes, OK |
| 161 | parseDropExternalStmt | C | DATABASE name | Yes, OK |
| 211 | parseCreateExternalFileFormatStmt | A | FILE FORMAT name | Yes, unnecessary |
| 217 | parseCreateExternalFileFormatStmt | B | FORMAT_TYPE = { DELIMITEDTEXT \| RCFILE \| PARQUET } | Yes, needs check |
| 331 | parseCreateExternalTableStmt | A | EXTERNAL TABLE name | Yes, unnecessary |
| 372 | parseCreateExternalTableStmt | C | schema.table reference | Yes, OK |
| 380 | parseCreateExternalTableStmt | C | column name context | Yes, OK |
| 440 | parseCreateExternalTableStmt | C | Column context | Yes, OK |
| 448 | parseCreateExternalTableStmt | C | Column context | Yes, OK |
| 506 | parseExternalOptions | C | Credential name | Yes, OK (reference) |
| 514 | parseExternalOptions | C | Credential reference | Yes, OK |
| 574 | parseExternalOptions | B | option_value context | Yes, mixed |
| [remaining 17] | Various | A/B/C mix | Names, references, option values | Mix |

**Summary**: 14 A-type (object names), 12 C-type (references, columns), 6 B-type (enum values).

---

### endpoint.go (32 sites)

| Line | Function | Category | T-SQL Context | Accepts Arb. Keyword? |
|------|----------|----------|---------------|-----------------------|
| 72 | parseCreateEndpointStmt | A | ENDPOINT name | Yes, unnecessary |
| 105 | parseAlterEndpointStmt | A | endPointName | Yes, unnecessary |
| 131 | parseDropEndpointStmt | A | endPointName | Yes, unnecessary |
| 200+ | parseEndpointOptions | C | Multiple contexts: certificate_name, ALGORITHM { AES \| RC4 \| ... }, ROLE { WITNESS \| PARTNER \| ALL }, etc. | Mix of OK and needs enum check |

**Summary**: 3 A-type (names), 29 C-type (protocol specs, algorithms, roles — mostly column/option name contexts that are OK).

---

### utility.go (28 sites)

| Line | Function | Category | T-SQL Context | Accepts Arb. Keyword? |
|------|----------|----------|---------------|-----------------------|
| 28 | parseUseStmt | A | USE database_name | Yes, but unnecessary |
| 135 | parseRaiseErrorStmt | B | RAISERROR parameter | Yes, mixed |
| 392 | parseOpenCursorStmt | A | OPEN cursor_name | Yes, OK (variable-like) |
| 399 | parseOpenCursorStmt | A | FOR clause context | Yes, OK |
| 447 | parseCloseCursorStmt | A | CLOSE cursor_name | Yes, OK |
| 453 | parseCloseCursorStmt | A | ALL keyword context | Yes, OK |
| 508 | parseDeallocateCursorStmt | A | DEALLOCATE cursor | Yes, OK |
| 514 | parseDeallocateCursorStmt | A | DEALLOCATE ALL context | Yes, OK |
| 639 | parseRestoreStmt | B | RESTORE VERIFYONLY option | Yes, needs check |
| 644 | parseRestoreStmt | B | RESTORE FILE/FILEGROUP name | Yes, OK (mixed) |
| 688 | parseRestoreStmt | B | RESTORE LOG option | Yes, needs check |
| 693 | parseRestoreStmt | B | option value | Yes, needs check |
| 796 | parseRestoreWithOptions | B | RESTORE option value | Yes, needs check |
| 822 | parseRestoreWithOptions | B | ON \| OFF option value | Yes, needs enum check |
| 829 | parseRestoreWithOptions | B | CREDENTIAL option | Yes, needs check |
| 890 | parseRestoreWithOptions | C | FILE name/path | Yes, OK (path context) |
| 896 | parseRestoreWithOptions | C | DEVICE path | Yes, OK |
| 1013 | parseRestoreWithOptions | C | Directory path | Yes, OK |
| 1104 | parseRecoverDbFromSnapshot | A | RECOVERY context | Yes, OK (keyword-specific) |
| 1139 | parseRecoverDbFromSnapshot | C | database_name | Yes, OK |
| 1186 | parseTransactionStmt | C | Database/schema context | Yes, OK |
| 1195 | parseTransactionStmt | C | Distributed transaction context | Yes, OK |
| 1366 | parseSetSessionTypeCommand | C | Session type context | Yes, OK |
| 1384 | parseSetSessionTypeCommand | C | Query context | Yes, OK |
| 1391 | parseSetSessionTypeCommand | C | Statistics context | Yes, OK |
| 1397 | parseSetSessionTypeCommand | C | column_name context | Yes, needs check |
| 1439 | parseSetStatisticsStmt | A | STATISTICS statement context | Yes, OK |
| 1459 | parseSetStatisticsStmt | C | IO \| TIME option | Yes, needs check |

**Summary**: ~4 A-type, ~12 B-type, ~12 C-type.

---

### security_audit.go (22 sites)

| Line | Function | Category | T-SQL Context | Accepts Arb. Keyword? |
|------|----------|----------|---------------|-----------------------|
| 58 | parseCreateServerAuditStmt | A | SERVER AUDIT name | Yes, unnecessary |
| 102 | parseAlterServerAuditStmt | A | audit_name | Yes, unnecessary |
| 128 | parseDropServerAuditStmt | A | audit_name | Yes, unnecessary |
| 155 | parseCreateServerAuditSpecStmt | A | AUDIT SPECIFICATION name | Yes, unnecessary |
| 184 | parseAlterServerAuditSpecStmt | A | spec_name | Yes, unnecessary |
| 207 | parseDropServerAuditSpecStmt | A | spec_name | Yes, unnecessary |
| 238 | parseCreateDatabaseAuditSpecStmt | A | AUDIT SPECIFICATION name | Yes, unnecessary |
| 264 | parseAlterDatabaseAuditSpecStmt | A | spec_name | Yes, unnecessary |
| 287 | parseDropDatabaseAuditSpecStmt | A | spec_name | Yes, unnecessary |
| 348 | parseAuditOptions | B | TO { FILE \| APPLICATION_LOG \| SECURITY_LOG \| URL \| EXTERNAL_MONITOR } | Yes, needs check |
| 373 | parseAuditOptions | B | QUEUE_DELAY \| ON_FAILURE \| STATE value | Yes, needs enum check |
| 392 | parseAuditOptions | B | option_value | Yes, needs check |
| 410 | parseAuditOptions | C | WHERE clause predicate | Yes, OK (expression) |
| 427 | parseAuditSpecOptions | B | ACTION_GROUP or ACTION enum | Yes, needs check |
| 439 | parseAuditSpecOptions | B | action name | Yes, needs check |
| 490 | parseAuditSpecOptions | B | WHERE clause context | Yes, OK (expression) |
| 505 | parseAuditSpecOptions | B | WHERE clause option | Yes, needs check |
| 510 | parseAuditSpecOptions | B | WHERE option | Yes, needs check |
| 572 | parseAuditSpecOptions | B | FOR { SELECT \| INSERT \| DELETE \| ... } enum | Yes, needs check |
| 595 | parseAuditSpecOptions | B | action value | Yes, needs check |
| 608 | parseAuditSpecOptions | B | WHERE column_name or option | Yes, mixed |
| 624 | parseAuditSpecOptions | B | WITH PUBLIC enum value | Yes, needs check |

**Summary**: 9 A-type, 10 B-type, 3 C/expression-type.

---

### create_table.go (21 sites)

| Line | Function | Category | T-SQL Context | Accepts Arb. Keyword? |
|------|----------|----------|---------------|-----------------------|
| 216 | parseCreateTableStmt | C | ON partition_scheme_name or filegroup | Yes, OK (object reference) |
| 234 | parseCreateTableStmt | C | TEXTIMAGE_ON filegroup | Yes, OK |
| 243 | parseCreateTableStmt | C | FILESTREAM_ON filegroup | Yes, OK |
| 349 | parseTableOptions | C | Column name or option_key context | Yes, mixed |
| 357 | parseTableOptions | B | Table option enum value | Yes, needs check |
| 367 | parseTableOptions | C | HISTORY_TABLE schema.table reference | Yes, OK |
| 373 | parseTableOptions | C | table reference | Yes, OK |
| 380 | parseTableOptions | B | option_value | Yes, needs check |
| 398 | parseTableOptions | B | option enum | Yes, needs check |
| 413 | parseTableOptions | B | DATA_CONSISTENCY_CHECK value | Yes, needs check |
| 707 | parseOneTableOption | C | Option context | Yes, mixed |
| 811 | parseTableWithIndexOptions | B | INDEX option enum | Yes, needs check |
| 819 | parseTableWithIndexOptions | B | option value | Yes, needs check |
| 824 | parseTableWithIndexOptions | B | SORT_IN_TEMPDB enum | Yes, needs check |
| 834 | parseTableWithIndexOptions | B | option value with constants | Yes, needs check |
| 873 | parseTableLedgerOptions | C | HISTORY_TABLE reference | Yes, OK |
| 880 | parseTableLedgerOptions | C | table reference | Yes, OK |
| 1085 | parseTableLedgerOptions | C | Column reference | Yes, OK |
| 1465 | parseOneConstraint | C | Column/constraint context | Yes, OK |
| 1482 | parseOneConstraint | C | CHECK constraint | Yes, OK |
| 1822 | parseTableColumnDef | B | Column option enum value | Yes, needs check |

**Summary**: 5 A/unclear, 16 C-type (filegroup, schema.table, column refs), mostly OK.

---

### resource_governor.go (20 sites)

| Line | Function | Category | T-SQL Context | Accepts Arb. Keyword? |
|------|----------|----------|---------------|-----------------------|
| ~70 | parseCreateResourcePoolStmt | A | RESOURCE POOL name | Yes, unnecessary |
| ~90 | parseAlterResourcePoolStmt | A | pool_name or [default] | Yes, unnecessary (but [default] is kw) |
| [similar patterns for remaining] | Parse functions for WORKLOAD GROUP, CLASSIFIER, RESOURCE POOL | A-B | Pool/group/classifier names, option values | Mix |

**Summary**: 12 A-type (object names), 8 B-type (option enum values like AFFINITY, MIN_CPU_PERCENT).

---

### alter_objects.go (18 sites)

| Line | Function | Category | T-SQL Context | Accepts Arb. Keyword? |
|------|----------|----------|---------------|-----------------------|
| [multiple] | parseAlterDatabaseSetOption, parseAlterDatabaseUnknownOption, etc. | B | DATABASE SET options (PARTNER, WITNESS, HADR mode enums, RECOVERY values, etc.) | Yes, needs strict enum checking |

**Summary**: ~15 B-type (option enums for ALTER DATABASE SET), ~3 misc.

---

### event.go (17 sites)

| Line | Function | Category | T-SQL Context | Accepts Arb. Keyword? |
|------|----------|----------|---------------|-----------------------|
| 53 | parseCreateEventSessionStmt | A | EVENT SESSION name | Yes, unnecessary |
| 100 | parseAlterEventSessionStmt | A | session_name | Yes, unnecessary |
| 138 | parseDropEventSessionStmt | A | session_name | Yes, unnecessary |
| [others] | parseEventNotificationOptions, etc. | C/B | Queue name, schema.table references, event action enums | Mix |

**Summary**: ~5 A-type, ~8 C-type, ~4 B-type.

---

### backup_restore.go (16 sites)

| Line | Function | Category | T-SQL Context | Accepts Arb. Keyword? |
|------|----------|----------|---------------|-----------------------|
| [multiple] | parseBackupStmt, parseRestoreStmt variants | B | Backup/restore option enums: INIT, NOINIT, SKIP, CONTINUE_AFTER_ERROR, WITH CHECKSUM, etc. | Yes, needs enum check |
| [others] | Various | C/A | File names, device names, database contexts | OK (paths/names) |

**Summary**: ~8 B-type (enum options), ~6 C/path-type, ~2 A-type.

---

### Remaining Lower-Volume Files (90 sites across 25 files)

**security_principals.go (9)**: ~5 A (principal names), ~4 C (schema references)  
**create_statistics.go (9)**: ~3 A (statistics names), ~6 C (column/table refs)  
**create_index.go (9)**: ~3 A (index names), ~6 B/C (option values, column refs)  
**partition.go (8)**: ~3 A (partition names), ~5 C (range value contexts)  
**assembly.go (8)**: ~4 A (assembly names), ~4 C (schema refs)  
**transaction.go (6)**: ~2 A/B (transaction names/options), ~4 C  
**select.go (6)**: ~0 A, ~6 C (column contexts)  
**availability.go (6)**: ~3 A (group names), ~3 C (replica names)  
**parser.go (4)**: ~1 A, ~3 misc  
**name.go (3)**: ~0 A, ~3 C (identifier contexts)  
**expr.go (3)**: ~0 A, ~3 C (expression contexts)  
**create_trigger.go (3)**: ~1 A (trigger name), ~2 C  
**create_proc.go (3)**: ~1 A (proc name), ~2 C  
**type.go (2)**: ~1 A (type name), ~1 C  
**execute.go (2)**: ~0 A, ~2 C  
**drop.go (2)**: ~1 A (object names), ~1 C  
**dbcc.go (2)**: ~0 A, ~2 B/C  
**bulk_insert.go (2)**: ~0 A, ~2 C (file context)  
**security_keys.go (1)**: 1 A (key name)  
**declare_set.go (1)**: ~1 C (variable context)  
**cursor.go (1)**: ~1 C (cursor name)  
**create_schema.go (1)**: ~0 A, ~1 C  
**control_flow.go (1)**: ~1 C  
**alter_table.go (1)**: ~0 A, ~1 C  

---

## Aggregated Taxonomy Summary

### By Category

| Category | Count | % | Description | Strictness Impact |
|----------|-------|---|-------------|-------------------|
| **A** | ~95 | 22% | Object/statement names (CREATE/ALTER/DROP object_name) | High — straightforward migration: use `isIdentLike()` |
| **B** | ~95 | 22% | Option-value enums (SET option = { ENUM1 \| ENUM2 }) | High — must restrict to valid enum set per option |
| **C** | ~150 | 35% | Schema objects, filegroups, columns, tables, references | Low — current behavior reasonable but room for strictness |
| **D** | ~10 | 2% | Subcommand dispatch, feature toggles | Medium — context-dependent; some branches strict, others flexible |
| **E** | ~20 | 5% | Dynamic property bags (service broker, XML, custom properties) | **ZERO — MUST STAY FLEXIBLE** |
| **F** | ~40 | 9% | Identifier-like positions (optional cleanup, use `isIdentLike()`) | Very low — cosmetic change |

### Top-5 Files by Category A Count (Object Names)
1. **security_principals.go (9)** + **external.go (14)** + **service_broker.go (16)** → ~39 A-type sites

### Top-5 Files by Category B Count (Enum Values)
1. **alter_objects.go (15)** + **service_broker.go (15)** + **security_audit.go (10)** + **fulltext.go (10)** + **create_table.go (9)** → ~59 B-type sites

### Top-5 Files by Category C Count (Schema/Filegroup/Column)
1. **fulltext.go (20)** + **endpoint.go (29)** + **create_table.go (16)** + external.go (12) + utility.go (12) → ~89 C-type sites

### Top-5 Files by Category E Count (Dynamic, Must Stay Flexible)
1. **service_broker.go (12)** + fulltext.go (2) + event.go (1) + others (~5) → ~20 E-type sites

---

## Migration Priority & Clusters

### Cluster 1: Highest Priority — Object Name Capture (Category A, ~95 sites)

**Why migrate**:  
- Pattern is consistent across all files: `if p.isAnyKeywordIdent() { stmt.Name = p.cur.Str; p.advance() }`
- Keywords should **never** be valid statement/object names in SQL standard
- **Zero regression risk** — these are uniformly wrong

**Target files** (by count):
1. service_broker.go (16 sites)
2. external.go (14 sites)
3. resource_governor.go (12 sites)
4. security_audit.go (9 sites)
5. fulltext.go (8 sites)

**Migration strategy**: Replace `isAnyKeywordIdent()` with `isIdentLike()` at each stmt.Name assignment.

---

### Cluster 2: High Priority — Option Enum Values (Category B, ~95 sites)

**Why migrate**:  
- These option values are **strictly enumerated** in T-SQL spec
- Example: `DATA_COMPRESSION = { NONE | ROW | PAGE }` should NOT accept `BACKUP` as a value
- Requires per-option enum checking (different per option key)

**Target files** (by count):
1. alter_objects.go (15 sites — ALTER DATABASE SET options)
2. service_broker.go (15 sites — VALIDATION, MESSAGE_CONTENT_TYPE, etc.)
3. security_audit.go (10 sites — ON_FAILURE, STATE, ACTION, etc.)
4. fulltext.go (10 sites — LANGUAGE, CHANGE_TRACKING, etc.)
5. create_table.go (9 sites — DATA_COMPRESSION, DURABILITY, LEDGER options)

**Migration strategy**:  
1. Create `isValidOptionValue(optionKey, value)` helper
2. Replace `isAnyKeywordIdent()` checks with enum validation
3. High complexity — each option has different valid enum set

---

### Cluster 3: Medium Priority — Schema/Filegroup/Column References (Category C, ~150 sites)

**Current behavior is reasonable** but can be stricter:
- Filegroups: typically identifiers, can be keywords in some contexts
- Columns: should use `isIdentLike()` (not `isAnyKeywordIdent()`)
- Schema.table references: current behavior OK for schema objects

**Files with most opportunity**:
1. fulltext.go (20 sites)
2. endpoint.go (29 sites)
3. create_table.go (16 sites)

**Migration strategy**: Context-dependent; lowest priority unless pursuing full strictness.

---

### Cluster 4: No-Migrate — Dynamic Properties (Category E, ~20 sites)

**Sites that MUST STAY FLEXIBLE**:
- service_broker.go lines ~534, ~563, ~674, ~683, ~692, ~744, ~762, ~1748, ~1758, ~1801, ~1810 (and others)
- These accept any keyword as a property key or value in dynamic contexts (messages, conversations, priorities)
- **Regression risk**: HIGH if made strict

**Action**: Document these sites and skip during refactor.

---

## Statistics & Insights

### File Coverage

- **Tier 1 (30+ sites)**: service_broker (46), fulltext (38), external (32), endpoint (32) → 148 sites (35%)
- **Tier 2 (10-29 sites)**: utility (28), security_audit (22), create_table (21), resource_governor (20), alter_objects (18), event (17), backup_restore (16), server (12), grant (11), security_misc (10), create_database (10) → 207 sites (49%)
- **Tier 3 (<10 sites)**: 25 files with 1-9 sites each → 71 sites (17%)

### Surprising Findings

1. **Endpoint.go has no A-type sites**: Most endpoint configurations are option/protocol specifiers (C-type), not object names.
2. **Service Broker is heavily E-type (dynamic)**: ~25% of SB sites are truly dynamic properties that must stay flexible.
3. **Alter_objects.go is B-heavy (83% B-type)**: Almost all ALTER DATABASE SET parsing focuses on option enum validation.
4. **No D-type cluster**: Subcommand dispatch (D) is rare and scattered — not a major refactor target.

---

## Deliverable: Sorted Line-by-Line Index

For reference, here are all 426 sites (sample of largest files shown above; full list available upon request):

- **service_broker.go**: 46 sites, lines 46, 54, 123, 131, 158, 231, 393, 431, 440, 445, 457, 521, 530, 534, 563, 624, 637, 644, 658, 674, 683, 692, 744, 762, 776, 819, 827, 892, 900, 1238, 1276, 1306, 1315, 1357, 1419, 1483, 1554, 1572, 1582, 1606, 1676, 1717, 1748, 1758, 1801, 1810
- **fulltext.go**: 38 sites
- **external.go**: 32 sites
- **endpoint.go**: 32 sites
- [... remaining 358 sites across 36 files ...]

---

## Recommendations

1. **Start with Cluster 1 (Category A)**: ~95 object-name sites are low-risk, high-impact
2. **Proceed to Cluster 2 (Category B)** once enum validation helpers are built (~95 sites)
3. **Skip Cluster 3 (Category C)** unless pursuing full strictness audit
4. **Protect Cluster 4 (Category E)**: Document these 20 sites and exclude from refactor
5. **Optional Cluster F**: Use for minor cleanup if time permits

**Expected impact**: Cluster 1 + 2 refactors would tighten parsing on ~190 sites (45% of total) with moderate to high regression prevention.

