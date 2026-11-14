// Oracle's SQL Language Reference as an ANTLR 4 grammar: the syntax of the
// statement pages oracle/parser/bnf/ quotes, converted by
// oracle/parser/cmd/docsgrammar. Do not edit it: correct a page's text in
// docs-fixes.tsv, or the quote in oracle/parser/bnf/, and run the converter
// again.
//
// Keywords are literals. ID, STRING, NUMBER, EXPR, CONDITION, QUERY,
// DATA_TYPE and HINT are the values the pages write in lower case and
// define nowhere; tokens.tsv samples them. A fragment another page defines
// is copied under the page that uses it, except a query, which is QUERY. A
// clause no quoted page defines is left out where it is optional, and the
// form or definition that requires it is not converted. Those, and a form
// that does not read as notation, are listed under their page.
//
// pin: docs https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ sha256:c523850bf1fa8fa444fe7cf746e5fb2de3b63fe9c62b16b1bd1c5ab03d0b5b0e 2026-03-15

grammar OracleDocs;

statement
    : administerKeyManagement
    | alterAnalyticView
    | alterAttributeDimension
    | alterAuditPolicyUnifiedAuditing
    | alterCluster
    | alterDatabaseDictionary
    | alterDatabaseLink
    | alterDatabase
    | alterDimension
    | alterDiskgroup
    | alterFlashbackArchive
    | alterFunction
    | alterHierarchy
    | alterIndex
    | alterIndextype
    | alterInmemoryJoinGroup
    | alterJava
    | alterLibrary
    | alterLockdownProfile
    | alterMaterializedViewLog
    | alterMaterializedView
    | alterMaterializedZonemap
    | alterOperator
    | alterOutline
    | alterPackage
    | alterPluggableDatabase
    | alterProcedure
    | alterProfile
    | alterResourceCost
    | alterRole
    | alterRollbackSegment
    | alterSequence
    | alterSession
    | alterSynonym
    | alterSystem
    | alterTable
    | alterTablespaceSet
    | alterTablespace
    | alterTrigger
    | alterType
    | alterUser
    | alterView
    | analyze
    | associateStatistics
    | auditTraditionalAuditing
    | auditUnifiedAuditing
    | call
    | comment
    | commit
    | createAnalyticView
    | createAttributeDimension
    | createAuditPolicyUnifiedAuditing
    | createCluster
    | createContext
    | createControlfile
    | createDatabaseLink
    | createDatabase
    | createDimension
    | createDirectory
    | createDiskgroup
    | createEdition
    | createFlashbackArchive
    | createFunction
    | createHierarchy
    | createIndex
    | createIndextype
    | createInmemoryJoinGroup
    | createJava
    | createLibrary
    | createLockdownProfile
    | createMaterializedViewLog
    | createMaterializedView
    | createMaterializedZonemap
    | createOperator
    | createOutline
    | createPackageBody
    | createPackage
    | createPfile
    | createPluggableDatabase
    | createProcedure
    | createProfile
    | createRestorePoint
    | createRole
    | createRollbackSegment
    | createSchema
    | createSequence
    | createSpfile
    | createSynonym
    | createTable
    | createTablespaceSet
    | createTablespace
    | createTrigger
    | createTypeBody
    | createType
    | createUser
    | createView
    | delete
    | disassociateStatistics
    | dropAnalyticView
    | dropAttributeDimension
    | dropAuditPolicyUnifiedAuditing
    | dropCluster
    | dropContext
    | dropDatabaseLink
    | dropDatabase
    | dropDimension
    | dropDirectory
    | dropDiskgroup
    | dropEdition
    | dropFlashbackArchive
    | dropFunction
    | dropHierarchy
    | dropIndex
    | dropIndextype
    | dropInmemoryJoinGroup
    | dropJava
    | dropLibrary
    | dropLockdownProfile
    | dropMaterializedViewLog
    | dropMaterializedView
    | dropMaterializedZonemap
    | dropOperator
    | dropOutline
    | dropPackage
    | dropPluggableDatabase
    | dropProcedure
    | dropProfile
    | dropRestorePoint
    | dropRole
    | dropRollbackSegment
    | dropSequence
    | dropSynonym
    | dropTable
    | dropTablespaceSet
    | dropTablespace
    | dropTrigger
    | dropTypeBody
    | dropType
    | dropUser
    | dropView
    | explainPlan
    | flashbackDatabase
    | flashbackTable
    | grant
    | insert
    | lockTable
    | merge
    | noauditTraditionalAuditing
    | noauditUnifiedAuditing
    | purge
    | rename
    | revoke
    | rollback
    | savepoint
    | select
    | setConstraints
    | setRole
    | setTransaction
    | truncateCluster
    | truncateTable
    | update
    | alterDomain
    | alterJsonRelationalDualityView
    | createDomain
    | createJsonRelationalDualityView
    | createLogicalPartitionTracking
    | createMleEnv
    | createMleModule
    | createPmemFilestore
    | createPropertyGraph
    | createVectorIndex
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ADMINISTER-KEY-MANAGEMENT.html
administerKeyManagement
    : 'ADMINISTER' 'KEY' 'MANAGEMENT' (administerKeyManagement_keystore_management_clauses | administerKeyManagement_key_management_clauses | administerKeyManagement_secret_management_clauses | administerKeyManagement_zero_downtime_software_patching_clauses)
    ;
administerKeyManagement_keystore_management_clauses
    : administerKeyManagement_create_keystore
    | administerKeyManagement_open_keystore
    | administerKeyManagement_close_keystore
    | administerKeyManagement_backup_keystore
    | administerKeyManagement_alter_keystore_password
    | administerKeyManagement_merge_into_new_keystore
    | administerKeyManagement_merge_into_existing_keystore
    | administerKeyManagement_isolate_keystore
    | administerKeyManagement_unite_keystore
    ;
administerKeyManagement_key_management_clauses
    : administerKeyManagement_set_key
    | administerKeyManagement_create_key
    | administerKeyManagement_use_key
    | administerKeyManagement_set_key_tag
    | administerKeyManagement_export_keys
    | administerKeyManagement_import_keys
    | administerKeyManagement_migrate_key
    | administerKeyManagement_reverse_migrate_key
    | administerKeyManagement_move_keys
    ;
administerKeyManagement_secret_management_clauses
    : administerKeyManagement_add_update_secret
    | administerKeyManagement_delete_secret
    | administerKeyManagement_add_update_secret_seps
    | administerKeyManagement_delete_secret_seps
    ;
administerKeyManagement_zero_downtime_software_patching_clauses
    : 'SWITCHOVER' 'TO' 'LIBRARY' STRING 'FOR' 'ALL' 'CONTAINERS'
    ;
administerKeyManagement_create_keystore
    : 'CREATE' 'KEYSTORE' STRING 'IDENTIFIED' 'BY' ID
    | 'CREATE' 'LOCAL'? 'AUTO_LOGIN' 'KEYSTORE' 'FROM' 'KEYSTORE' STRING 'IDENTIFIED' 'BY' ID
    ;
administerKeyManagement_open_keystore
    : 'SET' 'KEYSTORE' 'OPEN' ('FORCE' 'KEYSTORE')? 'IDENTIFIED' 'BY' ID ('CONTAINER' '=' ('CURRENT' | 'ALL'))?
    ;
administerKeyManagement_close_keystore
    : 'SET' 'KEYSTORE' 'CLOSE' ('IDENTIFIED' 'BY' ID)? ('CONTAINER' '=' ('CURRENT' | 'ALL'))?
    ;
administerKeyManagement_backup_keystore
    : 'BACKUP' 'KEYSTORE' ('USING' ID)? ('FORCE' 'KEYSTORE')? 'IDENTIFIED' 'BY' ID ('TO' STRING)?
    ;
administerKeyManagement_alter_keystore_password
    : 'ALTER' 'KEYSTORE' 'PASSWORD' 'IDENTIFIED' 'BY' ID 'SET' ID ('WITH' 'BACKUP' ('USING' ID)?)? ('FORCE' 'KEYSTORE')?
    ;
administerKeyManagement_merge_into_new_keystore
    : 'MERGE' 'KEYSTORE' STRING ('IDENTIFIED' 'BY' ID)? 'AND' 'KEYSTORE' STRING ('IDENTIFIED' 'BY' ID)? 'INTO' 'NEW' 'KEYSTORE' STRING 'IDENTIFIED' 'BY' ID
    ;
administerKeyManagement_merge_into_existing_keystore
    : 'MERGE' 'KEYSTORE' STRING ('IDENTIFIED' 'BY' ID)? 'INTO' 'EXISTING' 'KEYSTORE' STRING 'IDENTIFIED' 'BY' ID ('WITH' 'BACKUP' ('USING' ID)?)?
    ;
administerKeyManagement_isolate_keystore
    : 'FORCE'? 'ISOLATE' 'KEYSTORE' 'IDENTIFIED' 'BY' ID 'FROM' 'ROOT' 'KEYSTORE' ('FORCE' 'KEYSTORE')? 'IDENTIFIED' 'BY' ('EXTERNAL' 'STORE' | ID)? ('WITH' 'BACKUP' ('USING' ID)?)?
    ;
administerKeyManagement_unite_keystore
    : 'UNITE' 'KEYSTORE' 'IDENTIFIED' 'BY' ID 'WITH' 'ROOT' 'KEYSTORE' ('FORCE' 'KEYSTORE')? 'IDENTIFIED' 'BY' ('EXTERNAL' 'STORE' | ID)? ('WITH' 'BACKUP' ('USING' ID)?)?
    ;
administerKeyManagement_set_key
    : 'SET' 'ENCRYPTION'? 'KEY' ('MKID' ':' 'MK')? ('USING' 'TAG' ID)? ('USING' 'ALGORITHM' ('AES256' | 'ARIA256' | 'GOST256' | 'SEED128'))? ('FORCE' 'KEYSTORE')? 'IDENTIFIED' 'BY' ID ('WITH' 'BACKUP' ('USING' ID)?)? ('CONTAINER' '=' ('CURRENT' | 'ALL'))?
    ;
administerKeyManagement_create_key
    : 'CREATE' 'ENCRYPTION'? 'KEY' ('MKID' ':' 'MK')? ('USING' 'TAG' ID)? ('USING' 'ALGORITHM' ('AES256' | 'ARIA256' | 'GOST256' | 'SEED128'))? ('FORCE' 'KEYSTORE')? 'IDENTIFIED' 'BY' ID ('WITH' 'BACKUP' ('USING' ID)?)? ('CONTAINER' '=' ('CURRENT' | 'ALL'))?
    ;
administerKeyManagement_use_key
    : 'USE' 'ENCRYPTION'? 'KEY' ID ('USING' 'TAG' ID)? ('FORCE' 'KEYSTORE')? 'IDENTIFIED' 'BY' ID ('WITH' 'BACKUP' ('USING' ID)?)?
    ;
administerKeyManagement_set_key_tag
    : 'SET' 'TAG' ID 'FOR' ID ('FORCE' 'KEYSTORE')? 'IDENTIFIED' 'BY' ID ('WITH' 'BACKUP' ('USING' ID)?)?
    ;
administerKeyManagement_export_keys
    : 'EXPORT' 'ENCRYPTION'? 'KEYS' 'WITH' 'SECRET' ID 'TO' STRING ('FORCE' 'KEYSTORE')? 'IDENTIFIED' 'BY' ID ('WITH' 'IDENTIFIER' 'IN' '(' (ID (',' ID)* | QUERY) ')')?
    ;
administerKeyManagement_import_keys
    : 'IMPORT' 'ENCRYPTION'? 'KEYS' 'WITH' 'SECRET' ID 'FROM' STRING ('FORCE' 'KEYSTORE')? 'IDENTIFIED' 'BY' ID ('WITH' 'BACKUP' ('USING' ID)?)?
    ;
administerKeyManagement_migrate_key
    : ('USE' 'ENCRYPTION' 'KEY' ID)? 'IDENTIFIED' 'BY' STRING ('FORCE' 'KEYSTORE')? 'MIGRATE' 'USING' ID ('WITH' 'BACKUP' ('USING' ID)?)?
    ;
administerKeyManagement_reverse_migrate_key
    : 'SET' 'ENCRYPTION'? 'KEY' 'IDENTIFIED' 'BY' ID ('FORCE' 'KEYSTORE')? 'REVERSE' 'MIGRATE' 'USING' STRING
    ;
administerKeyManagement_move_keys
    : 'MOVE' 'KEYS' 'TO' 'NEW' 'KEYSTORE' STRING 'IDENTIFIED' 'BY' ID 'FROM' 'FORCE'? 'KEYSTORE' 'IDENTIFIED' 'BY' ID ('WITH' 'BACKUP' ('USING' ID)?)? 'WITH' 'IDENTIFIER' 'IN' '(' (ID (',' ID)* | QUERY) ')'
    ;
administerKeyManagement_add_update_secret
    : ('ADD' | 'UPDATE') 'SECRET' ID 'FOR' 'CLIENT' ID ('USING' 'TAG' ID)? ('FORCE' 'KEYSTORE')? 'IDENTIFIED' 'BY' ID ('WITH' 'BACKUP' ('USING' ID)?)?
    ;
administerKeyManagement_delete_secret
    : 'DELETE' 'SECRET' 'FOR' 'CLIENT' ID ('FORCE' 'KEYSTORE')? 'IDENTIFIED' 'BY' ID ('WITH' 'BACKUP' ('USING' ID)?)?
    ;
administerKeyManagement_add_update_secret_seps
    : ('ADD' | 'UPDATE') 'SECRET' ID 'FOR' 'CLIENT' ID ('USING' 'TAG' ID)? ('FORCE' 'KEYSTORE')? 'IDENTIFIED' 'BY' ('EXTERNAL' 'STORE' | ID)? ('WITH' 'BACKUP' ('USING' ID)?)?
    ;
administerKeyManagement_delete_secret_seps
    : 'DELETE' 'SECRET' 'FOR' 'CLIENT' ID ('FORCE' 'KEYSTORE')? 'IDENTIFIED' 'BY' ('EXTERNAL' 'STORE' | ID)? ('WITH' 'BACKUP' ('USING' ID)?)?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-ANALYTIC-VIEW.html
alterAnalyticView
    : 'ALTER' 'ANALYTIC' 'VIEW' ('IF' 'EXISTS')? (ID '.')? ID ('RENAME' 'TO' ID | 'COMPILE' | alterAnalyticView_alter_add_cache_clause | alterAnalyticView_alter_drop_cache_clause)
    ;
alterAnalyticView_alter_add_cache_clause
    : 'ADD' 'CACHE' 'MEASURE' 'GROUP' '(' ID (',' ID)* ')' 'LEVELS' '(' alterAnalyticView_level_specification (',' alterAnalyticView_level_specification)* ')'
    ;
alterAnalyticView_alter_drop_cache_clause
    : 'DROP' 'CACHE' 'MEASURE' 'GROUP' '(' ID (',' ID)* ')' 'LEVELS' '(' alterAnalyticView_level_specification (',' alterAnalyticView_level_specification)* ')'
    ;
alterAnalyticView_level_specification
    : (ID '.')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-ATTRIBUTE-DIMENSION.html
alterAttributeDimension
    : 'ALTER' 'ATTRIBUTE' 'DIMENSION' ('IF' 'EXISTS')? (ID '.')? ID ('RENAME' 'TO' ID | 'COMPILE')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-AUDIT-POLICY-Unified-Auditing.html
alterAuditPolicyUnifiedAuditing
    : 'ALTER' 'AUDIT' 'POLICY' ID ('ADD' | 'DROP') alterAuditPolicyUnifiedAuditing_privilege_audit_clause? alterAuditPolicyUnifiedAuditing_action_audit_clause? alterAuditPolicyUnifiedAuditing_role_audit_clause?
    | 'ALTER' 'AUDIT' 'POLICY' ID 'CONDITION' ('DROP' | STRING ('EVALUATE' ('PER' 'STATEMENT' | 'PER' 'SESSION' | 'PER' 'INSTANCE'))?)
    | 'ALTER' 'AUDIT' 'POLICY' ID ('ADD' | 'DROP') 'ONLY' 'TOPLEVEL'
    ;
alterAuditPolicyUnifiedAuditing_privilege_audit_clause
    : 'PRIVILEGES' ID (',' ID)*
    ;
alterAuditPolicyUnifiedAuditing_action_audit_clause
    : 'ACTIONS' (alterAuditPolicyUnifiedAuditing_standard_actions | alterAuditPolicyUnifiedAuditing_component_actions) (',' (alterAuditPolicyUnifiedAuditing_standard_actions | alterAuditPolicyUnifiedAuditing_component_actions))*
    ;
alterAuditPolicyUnifiedAuditing_role_audit_clause
    : 'ROLES' ID (',' ID)*
    ;
alterAuditPolicyUnifiedAuditing_standard_actions
    : ID ('ON' (ID '.')? ID ('(' ID (',' ID)* ')')?)?
    ;
alterAuditPolicyUnifiedAuditing_component_actions
    : 'COMPONENT' '=' ID ID (',' ID)*
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-CLUSTER.html
// fix: left out where optional, since no quoted page defines it or what it needs: storage_clause
alterCluster
    : 'ALTER' 'CLUSTER' ('IF' 'EXISTS')? (ID '.')? ID (alterCluster_physical_attributes_clause | 'SIZE' NUMBER | 'MODIFY' 'PARTITION' ID alterCluster_allocate_extent_clause | alterCluster_allocate_extent_clause | alterCluster_deallocate_unused_clause | alterCluster_parallel_clause)
    ;
alterCluster_physical_attributes_clause
    : ('PCTUSED' NUMBER)? ('PCTFREE' NUMBER)? ('INITRANS' NUMBER)?
    ;
alterCluster_allocate_extent_clause
    : 'ALLOCATE' 'EXTENT' ('(' ('SIZE' alterCluster_size_clause | 'DATAFILE' STRING | 'INSTANCE' NUMBER)+ ')')?
    ;
alterCluster_deallocate_unused_clause
    : 'DEALLOCATE' 'UNUSED' ('KEEP' alterCluster_size_clause)?
    ;
alterCluster_parallel_clause
    : ('PARALLEL' NUMBER? | 'NOPARALLEL')
    ;
alterCluster_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T' | 'P')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-DATABASE-DICTIONARY.html
alterDatabaseDictionary
    : 'ALTER' 'DATABASE' 'DICTIONARY' ('ENCRYPT' 'CREDENTIALS' | 'REKEY' 'CREDENTIALS' | 'DELETE' 'CREDENTIALS' 'KEY')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-DATABASE-LINK.html
alterDatabaseLink
    : 'ALTER' 'SHARED'? 'PUBLIC'? 'DATABASE' 'LINK' ID ('IF' 'EXISTS')? alterDatabaseLink_dblink_authentication
    ;
alterDatabaseLink_dblink_authentication
    : 'CONNECT' 'TO' ID 'IDENTIFIED' 'BY' ID ('AUTHENTICATED' 'BY' ID 'IDENTIFIED' 'BY' ID)?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-DATABASE.html
alterDatabase
    : 'ALTER' 'DATABASE' ID? (alterDatabase_startup_clauses | alterDatabase_recovery_clauses | alterDatabase_database_file_clauses | alterDatabase_logfile_clauses | alterDatabase_controlfile_clauses | alterDatabase_standby_database_clauses | alterDatabase_default_settings_clauses | alterDatabase_instance_clauses | alterDatabase_security_clause | 'RENAME' 'GLOBAL_NAME' 'TO' ID | ('ENABLE' | 'DISABLE') 'BLOCK' 'CHANGE' 'TRACKING' ('USING' 'FILE' STRING)? | alterDatabase_cdb_fleet_clauses | alterDatabase_replay_upgrade_clause)
    ;
alterDatabase_startup_clauses
    : ('MOUNT' ('STANDBY' 'DATABASE' | 'CLONE' 'DATABASE')? | 'OPEN' ('READ' 'WRITE' | 'READ' 'ONLY')? ('RESETLOGS' | 'NORESETLOGS')? ('UPGRADE' | 'DOWNGRADE')?)
    ;
alterDatabase_recovery_clauses
    : ('BEGIN' 'BACKUP' | 'END' 'BACKUP' | alterDatabase_general_recovery | alterDatabase_managed_standby_recovery)
    ;
alterDatabase_database_file_clauses
    : ('RENAME' 'FILE' STRING 'TO' STRING (',' STRING 'TO' STRING)* | 'CREATE' 'DATAFILE' (STRING | NUMBER) 'AS' ('NEW' | alterDatabase_file_specification) | 'DATAFILE' (STRING | NUMBER) ('ONLINE' | 'OFFLINE' ('FOR' 'DROP')? | 'RESIZE' alterDatabase_size_clause | 'END' 'BACKUP' | 'ENCRYPT' | 'DECRYPT' | alterDatabase_autoextend_clause) | 'TEMPFILE' (STRING | NUMBER) ('RESIZE' alterDatabase_size_clause | alterDatabase_autoextend_clause | 'DROP' ('INCLUDING' 'DATAFILES')?) | 'MOVE' 'DATAFILE' (STRING | NUMBER) 'TO' STRING 'REUSE'? 'KEEP'?)
    ;
alterDatabase_logfile_clauses
    : ('ARCHIVELOG' | 'MANUAL' | 'NOARCHIVELOG' | 'NO'? 'FORCE' 'LOGGING' | 'SET' 'STANDBY' 'NOLOGGING' 'FOR' ('LOAD' 'PERFORMANCE' | 'DATA' 'AVAILABILITY') | 'RENAME' 'FILE' STRING 'TO' STRING (',' STRING 'TO' STRING)* | 'CLEAR' 'LOGFILE' alterDatabase_logfile_descriptor ('UNARCHIVED' | 'UNRECOVERABLE' 'DATAFILE')? | alterDatabase_add_logfile_clauses | alterDatabase_drop_logfile_clauses | alterDatabase_switch_logfile_clause | alterDatabase_supplemental_db_logging)
    ;
alterDatabase_controlfile_clauses
    : ('CREATE' ('PHYSICAL' | 'LOGICAL')? 'STANDBY' 'CONTROLFILE' 'AS' STRING 'REUSE'? | 'CREATE' 'FAR' 'SYNC' 'INSTANCE' 'CONTROLFILE' 'AS' STRING 'REUSE'? | 'BACKUP' 'CONTROLFILE' 'TO' STRING 'REUSE'? | 'BACKUP' 'CONTROLFILE' 'TO' 'TRACE' ('AS' STRING)? 'REUSE'? (('RESETLOGS' | 'NORESETLOGS'))?)
    ;
alterDatabase_standby_database_clauses
    : (alterDatabase_activate_standby_db_clause | alterDatabase_maximize_standby_db_clause | alterDatabase_register_logfile_clause | alterDatabase_commit_switchover_clause | alterDatabase_start_standby_clause | alterDatabase_stop_standby_clause | alterDatabase_convert_database_clause | alterDatabase_switchover_clause | alterDatabase_failover_clause)
    ;
alterDatabase_default_settings_clauses
    : ('DEFAULT' 'EDITION' '=' ID | 'DEFAULT' 'TABLESPACE' ID | 'DEFAULT' 'LOCAL'? 'TEMPORARY' 'TABLESPACE' ID | 'SET' 'DEFAULT' ('BIGFILE' | 'SMALLFILE') 'TABLESPACE' | alterDatabase_flashback_mode_clause | alterDatabase_undo_mode_clause | alterDatabase_set_time_zone_clause)
    ;
alterDatabase_instance_clauses
    : ('ENABLE' 'RESTRICTED' 'SESSION' | 'DISABLE' 'RESTRICTED' 'SESSION')
    ;
alterDatabase_security_clause
    : (alterDatabase_prepare_clause | alterDatabase_drop_mirror_copy | alterDatabase_lost_write_protection)
    ;
alterDatabase_cdb_fleet_clauses
    : (alterDatabase_lead_cdb_clause | alterDatabase_lead_cdb_uri_clause | alterDatabase_property_clause)
    ;
alterDatabase_replay_upgrade_clause
    : ('START' 'REPLAY' | 'STOP' 'REPLAY')
    ;
alterDatabase_general_recovery
    : 'RECOVER' ('AUTOMATIC' | 'FROM' STRING)? (alterDatabase_full_database_recovery | alterDatabase_partial_database_recovery) ('LOGFILE' STRING)? 'TEST'? ('ALLOW' NUMBER 'CORRUPTION')? ('CONTINUE' 'DEFAULT'? | 'CANCEL')? alterDatabase_parallel_clause?
    ;
alterDatabase_managed_standby_recovery
    : 'RECOVER' 'MANAGED' 'STANDBY' 'DATABASE' ('USING' 'ARCHIVED' 'LOGFILE' | 'USING' 'CURRENT' 'LOGFILE')? ('DISCONNECT' ('FROM' 'SESSION')?)? 'NODELAY'? ('UNTIL' 'CHANGE' NUMBER | 'UNTIL' 'CONSISTENT')? ('USING' 'INSTANCES' ('ALL' | NUMBER))? ('FINISH' ('FORCE' | 'WAIT' | 'NOWAIT')?)? ('CANCEL' ('IMMEDIATE' | 'WAIT' | 'NOWAIT')?)? ('TO' 'LOGICAL' 'STANDBY' ID ('KEEP' 'IDENTITY')?)? alterDatabase_parallel_clause?
    ;
alterDatabase_file_specification
    : STRING? ('SIZE' alterDatabase_size_clause)? 'REUSE'? alterDatabase_autoextend_clause?
    ;
alterDatabase_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T' | 'P')?
    ;
alterDatabase_autoextend_clause
    : 'AUTOEXTEND' ('OFF' | 'ON' ('NEXT' alterDatabase_size_clause)? alterDatabase_maxsize_clause?)
    ;
alterDatabase_logfile_descriptor
    : ('GROUP' NUMBER | '(' STRING (',' STRING)* ')')
    ;
alterDatabase_add_logfile_clauses
    : ('ADD' 'STANDBY'? 'LOGFILE' ('INSTANCE' STRING)? ('THREAD' NUMBER)? ('GROUP' NUMBER)? STRING (',' STRING)* | 'ADD' 'LOGFILE' 'MEMBER' 'STANDBY'? STRING (',' STRING)* 'TO' alterDatabase_logfile_descriptor)
    ;
alterDatabase_drop_logfile_clauses
    : ('DROP' 'LOGFILE' alterDatabase_logfile_descriptor | 'DROP' 'LOGFILE' 'MEMBER' STRING (',' STRING)*)
    ;
alterDatabase_switch_logfile_clause
    : 'SWITCH' 'LOGFILE' 'TO' 'BLOCK' 'SIZE' NUMBER
    ;
alterDatabase_supplemental_db_logging
    : ('ADD' 'SUPPLEMENTAL' 'LOG' ('DATA' | alterDatabase_supplemental_id_key_clause | alterDatabase_supplemental_plsql_clause | 'DATA' 'SUBSET' 'DATABASE' 'REPLICATION') | 'DROP' 'SUPPLEMENTAL' 'LOG' ('DATA' | alterDatabase_supplemental_id_key_clause | alterDatabase_supplemental_plsql_clause | 'DATA' 'SUBSET' 'DATABASE' 'REPLICATION'))
    ;
alterDatabase_activate_standby_db_clause
    : 'ACTIVATE' ('PHYSICAL' | 'LOGICAL')? 'STANDBY' 'DATABASE' ('FINISH' 'APPLY')?
    ;
alterDatabase_maximize_standby_db_clause
    : 'SET' 'STANDBY' 'DATABASE' 'TO' ('MAXIMIZE' 'PROTECTION' | 'MAXIMIZE' 'AVAILABILITY' | 'MAXIMIZE' 'PERFORMANCE')
    ;
alterDatabase_register_logfile_clause
    : 'RECOVER' ('MANAGED' 'STANDBY' 'DATABASE')? 'REGISTER' 'LOGFILE' STRING (',' STRING)* ('OR' 'REPLACE')? ('FOR' ID)?
    ;
alterDatabase_commit_switchover_clause
    : ('PREPARE' 'TO' 'SWITCHOVER' 'TO' ('LOGICAL' 'STANDBY' | 'PRIMARY' 'DATABASE') | 'COMMIT' 'TO' 'SWITCHOVER' 'TO' ('PHYSICAL'? 'STANDBY' | 'LOGICAL' 'STANDBY' | 'PRIMARY' 'DATABASE')) ('WITH' 'SESSION' 'SHUTDOWN' | 'WITHOUT' 'SESSION' 'SHUTDOWN')? ('WAIT' | 'NOWAIT')? 'CANCEL'?
    ;
alterDatabase_start_standby_clause
    : 'START' 'LOGICAL' 'STANDBY' 'APPLY' 'IMMEDIATE'? 'NODELAY'? 'INITIAL'? ('NEW' 'PRIMARY')? ('SKIP' 'FAILED' 'TRANSACTION'?)? 'FINISH'?
    ;
alterDatabase_stop_standby_clause
    : 'STOP' 'LOGICAL' 'STANDBY' 'APPLY'
    ;
alterDatabase_convert_database_clause
    : 'CONVERT' 'TO' ('PHYSICAL' 'STANDBY' | 'SNAPSHOT' 'STANDBY')
    ;
alterDatabase_switchover_clause
    : 'SWITCHOVER' 'TO' ('PRIMARY' | 'PHYSICAL' 'STANDBY' | 'LOGICAL' 'STANDBY') 'VERIFY'? 'FORCE'? ID
    ;
alterDatabase_failover_clause
    : 'FAILOVER' 'TO' ('PRIMARY' | 'PHYSICAL' 'STANDBY' | 'LOGICAL' 'STANDBY') 'FORCE'? ID
    ;
alterDatabase_flashback_mode_clause
    : 'FLASHBACK' ('ON' | 'OFF')
    ;
alterDatabase_undo_mode_clause
    : 'SET' 'UNDO' 'TABLESPACE' '=' ID
    ;
alterDatabase_set_time_zone_clause
    : 'SET' 'TIME_ZONE' '=' (STRING | 'DBTIMEZONE' | 'LOCAL')
    ;
alterDatabase_prepare_clause
    : 'PREPARE' 'MIRROR' 'COPY'
    ;
alterDatabase_drop_mirror_copy
    : 'DROP' 'MIRROR' 'COPY'
    ;
alterDatabase_lost_write_protection
    : ('ENABLE' 'LOST' 'WRITE' 'PROTECTION' | 'DISABLE' 'LOST' 'WRITE' 'PROTECTION')
    ;
alterDatabase_lead_cdb_clause
    : 'SET' 'LEAD_CDB' '='? (ID | 'NONE')
    ;
alterDatabase_lead_cdb_uri_clause
    : 'SET' 'LEAD_CDB_URI' '='? (STRING | 'NONE')
    ;
alterDatabase_property_clause
    : 'SET' 'PROPERTY' ID '=' EXPR
    ;
alterDatabase_full_database_recovery
    : ('DATABASE' | 'STANDBY' 'DATABASE' ('UNTIL' ('CANCEL' | 'TIME' '\'YYYY-MM-DD:HH24:MI:SS\'' | 'CHANGE' NUMBER | 'CONSISTENT'))? ('USING' 'BACKUP' 'CONTROLFILE')?)
    ;
alterDatabase_partial_database_recovery
    : ('TABLESPACE' ID (',' ID)* | 'DATAFILE' (STRING | NUMBER) (',' (STRING | NUMBER))*)
    ;
alterDatabase_parallel_clause
    : ('PARALLEL' NUMBER? | 'NOPARALLEL')
    ;
alterDatabase_maxsize_clause
    : 'MAXSIZE' (alterDatabase_size_clause | 'UNLIMITED')
    ;
alterDatabase_supplemental_id_key_clause
    : 'DATA' '(' ('PRIMARY' 'KEY' | 'UNIQUE' | 'FOREIGN' 'KEY' | 'ALL') 'COLUMNS' ')'
    ;
alterDatabase_supplemental_plsql_clause
    : 'DATA' '(' 'PL' '/' 'SQL' 'CALL' ')'
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-DIMENSION.html
alterDimension
    : 'ALTER' 'DIMENSION' (ID '.')? ID ('ADD' alterDimension_level_clause | 'ADD' alterDimension_hierarchy_clause | 'ADD' alterDimension_attribute_clause | 'ADD' alterDimension_extended_attribute_clause | 'DROP' alterDimension_level_clause | 'DROP' alterDimension_hierarchy_clause | 'DROP' alterDimension_attribute_clause | 'DROP' alterDimension_extended_attribute_clause | 'COMPILE')
    ;
alterDimension_level_clause
    : 'LEVEL' ID 'IS' '(' ID '.' ID (',' ID '.' ID)* ')' ('SKIP' 'WHEN' EXPR)?
    ;
alterDimension_hierarchy_clause
    : 'HIERARCHY' ID '(' ID 'CHILD' 'OF' ID ('JOIN' 'KEY' ID 'REFERENCES' ID)? (',' ID 'CHILD' 'OF' ID ('JOIN' 'KEY' ID 'REFERENCES' ID)?)* ')'
    ;
alterDimension_attribute_clause
    : 'ATTRIBUTE' ID 'DETERMINES' '(' ID (',' ID)* ')'
    ;
alterDimension_extended_attribute_clause
    : 'ATTRIBUTE' ID 'OF' ID 'DETERMINES' '(' ID (',' ID)* ')'
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-DISKGROUP.html
alterDiskgroup
    : 'ALTER' 'DISKGROUP' ID (alterDiskgroup_add_disk_clause | alterDiskgroup_drop_disk_clause | alterDiskgroup_resize_disk_clause | alterDiskgroup_replace_disk_clause | alterDiskgroup_rename_disk_clause | alterDiskgroup_disk_online_clause | alterDiskgroup_disk_offline_clause | alterDiskgroup_rebalance_diskgroup_clause | alterDiskgroup_check_diskgroup_clause | alterDiskgroup_diskgroup_template_clauses | alterDiskgroup_diskgroup_directory_clauses | alterDiskgroup_diskgroup_alias_clauses | alterDiskgroup_diskgroup_volume_clauses | alterDiskgroup_diskgroup_attributes | alterDiskgroup_modify_diskgroup_file | alterDiskgroup_drop_diskgroup_file_clause | alterDiskgroup_convert_redundancy_clause | alterDiskgroup_usergroup_clauses | alterDiskgroup_user_clauses | alterDiskgroup_file_permissions_clause | alterDiskgroup_file_owner_clause | alterDiskgroup_scrub_clause | alterDiskgroup_quotagroup_clauses | alterDiskgroup_filegroup_clauses | alterDiskgroup_undrop_disk_clause | alterDiskgroup_diskgroup_availability | alterDiskgroup_enable_disable_volume)
    ;
alterDiskgroup_add_disk_clause
    : 'ADD' 'DISK' alterDiskgroup_qualified_disk_clause (',' alterDiskgroup_qualified_disk_clause)*
    ;
alterDiskgroup_drop_disk_clause
    : 'DROP' 'DISK' ID (',' ID)* ('FORCE' | 'NOFORCE')? alterDiskgroup_rebalance_diskgroup_clause?
    | 'DROP' 'DISKS' 'IN' 'FAILGROUP' ID ('FORCE' | 'NOFORCE')? alterDiskgroup_rebalance_diskgroup_clause?
    ;
alterDiskgroup_resize_disk_clause
    : 'RESIZE' 'ALL' alterDiskgroup_size_clause
    ;
alterDiskgroup_replace_disk_clause
    : 'REPLACE' 'DISK' ID 'WITH' STRING ('FORCE' | 'NOFORCE')? ('POWER' NUMBER | 'POWER' 'LIMIT' NUMBER)? ('WAIT' | 'NOWAIT')?
    ;
alterDiskgroup_rename_disk_clause
    : 'RENAME' 'DISK' ID 'TO' ID (',' ID 'TO' ID)*
    | 'RENAME' 'DISKS' 'ALL'
    ;
alterDiskgroup_disk_online_clause
    : 'ONLINE' 'DISK' ID (',' ID)* ('POWER' NUMBER)? ('WAIT' | 'NOWAIT')?
    | 'ONLINE' 'DISKS' 'IN' 'FAILGROUP' ID ('POWER' NUMBER)? ('WAIT' | 'NOWAIT')?
    | 'ONLINE' 'ALL' ('POWER' NUMBER)? ('WAIT' | 'NOWAIT')?
    ;
alterDiskgroup_disk_offline_clause
    : 'OFFLINE' 'DISK' ID (',' ID)* alterDiskgroup_timeout_clause? ('WAIT' | 'NOWAIT')?
    | 'OFFLINE' 'DISKS' 'IN' 'FAILGROUP' ID alterDiskgroup_timeout_clause? ('WAIT' | 'NOWAIT')?
    ;
alterDiskgroup_rebalance_diskgroup_clause
    : 'REBALANCE' ('WITH' ('RESTORE' | 'BALANCE' | 'PREPARE' | 'COMPACT') (',' ('RESTORE' | 'BALANCE' | 'PREPARE' | 'COMPACT'))?)*
    | 'REBALANCE' ('WITHOUT' ('BALANCE' | 'PREPARE' | 'COMPACT') (',' ('BALANCE' | 'PREPARE' | 'COMPACT'))?)* ('POWER' NUMBER)? ('WAIT' | 'NOWAIT')?
    | 'MODIFY' 'POWER' NUMBER?
    ;
alterDiskgroup_check_diskgroup_clause
    : 'CHECK' ('REPAIR' | 'NOREPAIR')?
    | 'CHECK' 'ALL' ('REPAIR' | 'NOREPAIR')?
    | 'CHECK' 'DISK' ID (',' ID)* ('REPAIR' | 'NOREPAIR')?
    | 'CHECK' 'DISKS' 'IN' 'FAILGROUP' ID ('REPAIR' | 'NOREPAIR')?
    | 'CHECK' 'FILE' alterDiskgroup_file_specification (',' alterDiskgroup_file_specification)* ('REPAIR' | 'NOREPAIR')?
    ;
alterDiskgroup_diskgroup_template_clauses
    : 'ADD' 'TEMPLATE' ID 'ATTRIBUTES' '(' alterDiskgroup_qualified_template_clause ')'
    | 'MODIFY' 'TEMPLATE' ID 'ATTRIBUTES' '(' alterDiskgroup_qualified_template_clause ')'
    | 'ALTER' 'TEMPLATE' ID 'ATTRIBUTES' '(' alterDiskgroup_qualified_template_clause ')'
    | 'DROP' 'TEMPLATE' ID
    ;
alterDiskgroup_diskgroup_directory_clauses
    : 'ADD' 'DIRECTORY' STRING
    | 'DROP' 'DIRECTORY' STRING 'FORCE'?
    | 'RENAME' 'DIRECTORY' STRING 'TO' STRING
    ;
alterDiskgroup_diskgroup_alias_clauses
    : 'ADD' 'ALIAS' STRING 'FOR' STRING
    | 'DROP' 'ALIAS' STRING
    | 'RENAME' 'ALIAS' STRING 'TO' STRING
    ;
alterDiskgroup_diskgroup_volume_clauses
    : alterDiskgroup_add_volume_clause
    | alterDiskgroup_modify_volume_clause
    | 'RESIZE' 'VOLUME' ID alterDiskgroup_size_clause
    | 'DROP' 'VOLUME' ID
    ;
alterDiskgroup_diskgroup_attributes
    : 'SET' 'ATTRIBUTE' ID '=' EXPR (',' ID '=' EXPR)*
    ;
alterDiskgroup_modify_diskgroup_file
    : 'MODIFY' 'FILE' STRING alterDiskgroup_redundancy_clause?
    ;
alterDiskgroup_drop_diskgroup_file_clause
    : 'DROP' 'FILE' STRING
    | 'DROP' 'FILE' ID
    ;
alterDiskgroup_convert_redundancy_clause
    : 'CONVERT' 'REDUNDANCY'
    ;
alterDiskgroup_usergroup_clauses
    : 'ADD' 'USERGROUP' ID ('MEMBER' ID)?
    | 'MODIFY' 'USERGROUP' ID ('ADD' 'MEMBER' ID | 'DROP' 'MEMBER' ID)
    | 'DROP' 'USERGROUP' ID
    ;
alterDiskgroup_user_clauses
    : 'ADD' 'USER' ID (',' ID)*
    | 'DROP' 'USER' ID 'CASCADE'? (',' ID 'CASCADE'?)*
    | 'REPLACE' 'USER' ID 'WITH' ID
    ;
alterDiskgroup_file_permissions_clause
    : 'SET' 'PERMISSION' STRING ('FILE_PERM_OWNER' EXPR)? ('FILE_PERM_GROUP' EXPR)? ('FILE_PERM_OTHER' EXPR)?
    ;
alterDiskgroup_file_owner_clause
    : 'SET' 'OWNER' ID 'FILE' STRING
    | 'SET' 'GROUP' ID 'FILE' STRING
    ;
alterDiskgroup_scrub_clause
    : 'SCRUB' (('FILE' alterDiskgroup_file_specification)? | ('DISK' ID)?)? ('REPAIR' | 'NOREPAIR')? ('POWER' ('AUTO' | 'LOW' | 'HIGH' | 'MAX'))? ('WAIT' | 'NOWAIT')? ('FORCE' | 'NOFORCE')?
    | 'SCRUB' 'STOP'
    ;
alterDiskgroup_quotagroup_clauses
    : 'ADD' 'QUOTAGROUP' ID ('SET' 'QUOTA' '=' (alterDiskgroup_size_clause | 'UNLIMITED'))?
    | 'MODIFY' 'QUOTAGROUP' ID 'SET' 'QUOTA' '=' (alterDiskgroup_size_clause | 'UNLIMITED')
    | 'MOVE' 'FILEGROUP' ID 'TO' 'QUOTAGROUP' ID
    | 'DROP' 'QUOTAGROUP' ID
    ;
alterDiskgroup_filegroup_clauses
    : alterDiskgroup_add_filegroup_clause
    | alterDiskgroup_modify_filegroup_clause
    | alterDiskgroup_move_to_filegroup_clause
    | alterDiskgroup_drop_filegroup_clause
    ;
alterDiskgroup_undrop_disk_clause
    : 'UNDROP' 'DISKS'
    ;
alterDiskgroup_diskgroup_availability
    : 'MOUNT' 'ALL'? ('RESTRICTED' | 'NORMAL')? ('FORCE' | 'NOFORCE')?
    | 'DISMOUNT' 'ALL'? 'FORCE'?
    ;
alterDiskgroup_enable_disable_volume
    : 'ENABLE' 'VOLUME' ('ALL' | ID (',' ID)*)
    | 'DISABLE' 'VOLUME' ('ALL' | ID (',' ID)*)
    ;
alterDiskgroup_qualified_disk_clause
    : ('QUORUM' | 'REGULAR')? ('NAME' ID)? STRING alterDiskgroup_size_clause? ('FAILGROUP' ID)?
    ;
alterDiskgroup_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T' | 'P')?
    ;
alterDiskgroup_timeout_clause
    : 'DISK' 'REPAIR' 'TIME' (NUMBER 'MINUTE' | NUMBER 'HOUR')?
    ;
alterDiskgroup_file_specification
    : STRING? ('SIZE' alterDiskgroup_size_clause)? 'REUSE'? alterDiskgroup_autoextend_clause?
    ;
alterDiskgroup_qualified_template_clause
    : alterDiskgroup_redundancy_clause? alterDiskgroup_striping_clause?
    ;
alterDiskgroup_add_volume_clause
    : 'ADD' 'VOLUME' ID alterDiskgroup_size_clause alterDiskgroup_redundancy_clause? ('STRIPE_WIDTH' alterDiskgroup_size_clause)? ('STRIPE_COLUMNS' NUMBER)?
    ;
alterDiskgroup_modify_volume_clause
    : 'MODIFY' 'VOLUME' ID ('MOUNTPATH' STRING)? ('USAGE' STRING)?
    ;
alterDiskgroup_redundancy_clause
    : 'MIRROR'
    | 'HIGH'
    | 'UNPROTECTED'
    | 'PARITY'
    | 'DOUBLE'
    ;
alterDiskgroup_add_filegroup_clause
    : 'ADD' 'FILEGROUP' ID ('DATABASE' ID?)? ('CLUSTER' ID)? ('VOLUME' ID)? ('TEMPLATE' ID)? ('SET' ID '=' EXPR (',' ID '=' EXPR)?)*
    ;
alterDiskgroup_modify_filegroup_clause
    : 'MODIFY' 'FILEGROUP' ID ('SET' ID '=' EXPR (',' ID '=' EXPR)?)*
    ;
alterDiskgroup_move_to_filegroup_clause
    : 'MOVE' 'FILE' alterDiskgroup_file_specification 'TO' 'FILEGROUP' ID
    ;
alterDiskgroup_drop_filegroup_clause
    : 'DROP' 'FILEGROUP' ID ('CASCADE' ('FOR' 'DATABASE' ID ('FOR' 'PLUGGABLE' 'DATABASE' ID)?)?)?
    ;
alterDiskgroup_autoextend_clause
    : 'AUTOEXTEND' ('OFF' | 'ON' ('NEXT' alterDiskgroup_size_clause)? alterDiskgroup_maxsize_clause?)
    ;
alterDiskgroup_striping_clause
    : 'FINE'
    | 'COARSE'
    ;
alterDiskgroup_maxsize_clause
    : 'MAXSIZE' (alterDiskgroup_size_clause | 'UNLIMITED')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-FLASHBACK-ARCHIVE.html
alterFlashbackArchive
    : 'ALTER' 'FLASHBACK' 'ARCHIVE' ID ('SET' 'DEFAULT' | ('ADD' | 'MODIFY') 'TABLESPACE' ID alterFlashbackArchive_flashback_archive_quota? | 'REMOVE' 'TABLESPACE' ID | 'MODIFY' 'RETENTION' alterFlashbackArchive_flashback_archive_retention | 'PURGE' ('ALL' | 'BEFORE' 'SCN' NUMBER | 'BEFORE' 'TIMESTAMP' EXPR) | 'NO'? 'OPTIMIZE' 'DATA')
    ;
alterFlashbackArchive_flashback_archive_quota
    : 'QUOTA' NUMBER ('M' | 'G' | 'T' | 'P' | 'E')
    ;
alterFlashbackArchive_flashback_archive_retention
    : NUMBER ('YEAR' | 'MONTH' | 'DAY')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-FUNCTION.html
alterFunction
    : 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? (ID '.')? ID alterFunction_function_compile_clause ('EDITIONABLE' | 'NONEDITIONABLE')?
    ;
alterFunction_function_compile_clause
    : 'COMPILE' 'DEBUG'? alterFunction_compiler_parameters_clause? ('REUSE' 'SETTINGS')?
    ;
alterFunction_compiler_parameters_clause
    : ID '=' EXPR
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-HIERARCHY.html
alterHierarchy
    : 'ALTER' 'HIERARCHY' ('IF' 'EXISTS')? (ID '.')? ID ('RENAME' 'TO' ID | 'COMPILE')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-INDEX.html
// fix: left out where optional, since no quoted page defines it or what it needs: index_ilm_clause, storage_clause, xmlindex_parameters_clause
alterIndex
    : 'ALTER' 'INDEX' ('IF' 'EXISTS')? (ID '.')? ID (alterIndex_deallocate_unused_clause | alterIndex_allocate_extent_clause | alterIndex_shrink_clause | alterIndex_parallel_clause | alterIndex_physical_attributes_clause | alterIndex_logging_clause | alterIndex_partial_index_clause | alterIndex_rebuild_clause | alterIndex_alter_index_partitioning | 'PARAMETERS' '(' STRING ')' | ('DEFERRED' | 'IMMEDIATE') 'INVALIDATION' | 'COMPILE' | 'ENABLE' | 'DISABLE' | ('USABLE' | 'UNUSABLE') 'ONLINE'? | ('VISIBLE' | 'INVISIBLE') | 'RENAME' 'TO' ID | 'COALESCE' ('CLEANUP' 'ONLY'?)? alterIndex_parallel_clause? | 'MONITORING' 'USAGE' | 'NOMONITORING' 'USAGE' | 'UPDATE' 'BLOCK' 'REFERENCES' | alterIndex_annotations_clause)
    ;
alterIndex_deallocate_unused_clause
    : 'DEALLOCATE' 'UNUSED' ('KEEP' alterIndex_size_clause)?
    ;
alterIndex_allocate_extent_clause
    : 'ALLOCATE' 'EXTENT' ('(' ('SIZE' alterIndex_size_clause | 'DATAFILE' STRING | 'INSTANCE' NUMBER)+ ')')?
    ;
alterIndex_shrink_clause
    : 'SHRINK' 'SPACE' 'COMPACT'? 'CASCADE'?
    ;
alterIndex_parallel_clause
    : ('NOPARALLEL' | 'PARALLEL' NUMBER?)
    ;
alterIndex_physical_attributes_clause
    : ('PCTFREE' NUMBER)? ('PCTUSED' NUMBER)? ('INITRANS' NUMBER)? ('MAXTRANS' NUMBER)?
    ;
alterIndex_logging_clause
    : ('LOGGING' | 'NOLOGGING')
    ;
alterIndex_partial_index_clause
    : 'INDEXING' ('FULL' | 'PARTIAL')
    ;
alterIndex_rebuild_clause
    : 'REBUILD' ('PARTITION' ID | 'SUBPARTITION' ID)? 'ONLINE'? ('TABLESPACE' ID)? ('PARAMETERS' '(' STRING ')')? (('REVERSE' | 'NOREVERSE'))? alterIndex_parallel_clause? alterIndex_physical_attributes_clause? alterIndex_index_compression? alterIndex_logging_clause? alterIndex_partial_index_clause? (('DEFERRED' | 'IMMEDIATE') 'INVALIDATION')?
    ;
alterIndex_alter_index_partitioning
    : (alterIndex_modify_index_default_attrs | alterIndex_add_hash_index_partition | alterIndex_modify_index_partition | alterIndex_rename_index_partition | alterIndex_drop_index_partition | alterIndex_split_index_partition | alterIndex_coalesce_index_partition | alterIndex_modify_index_subpartition)
    ;
alterIndex_annotations_clause
    : 'ANNOTATIONS' '(' ('ADD' alterIndex_annotation | 'DROP' ID (',' ID)* | 'REMOVE' 'ALL') ')'
    ;
alterIndex_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T' | 'P')?
    ;
alterIndex_index_compression
    : (alterIndex_prefix_compression | alterIndex_advanced_index_compression | 'NOCOMPRESS')
    ;
alterIndex_modify_index_default_attrs
    : 'MODIFY' 'DEFAULT' 'ATTRIBUTES' ('FOR' 'PARTITION' ID)? alterIndex_physical_attributes_clause? ('TABLESPACE' (ID | 'DEFAULT'))? alterIndex_logging_clause?
    ;
alterIndex_add_hash_index_partition
    : 'ADD' 'PARTITION' ID? ('TABLESPACE' ID)? alterIndex_index_compression? alterIndex_parallel_clause?
    ;
alterIndex_modify_index_partition
    : 'MODIFY' 'PARTITION' ID ((alterIndex_deallocate_unused_clause | alterIndex_allocate_extent_clause | alterIndex_physical_attributes_clause | alterIndex_logging_clause | alterIndex_index_compression) | 'PARAMETERS' '(' STRING ')' | 'COALESCE' ('CLEANUP' 'ONLY'?)? | 'UPDATE' 'BLOCK' 'REFERENCES' | 'UNUSABLE')
    ;
alterIndex_rename_index_partition
    : 'RENAME' ('PARTITION' ID | 'SUBPARTITION' ID) 'TO' ID
    ;
alterIndex_drop_index_partition
    : 'DROP' 'PARTITION' ID
    ;
alterIndex_split_index_partition
    : 'SPLIT' 'PARTITION' ID 'AT' '(' EXPR (',' EXPR)* ')' ('INTO' '(' alterIndex_index_partition_description ',' alterIndex_index_partition_description ')')? alterIndex_parallel_clause?
    ;
alterIndex_coalesce_index_partition
    : 'COALESCE' 'PARTITION' alterIndex_parallel_clause?
    ;
alterIndex_modify_index_subpartition
    : 'MODIFY' 'SUBPARTITION' ID ('UNUSABLE' | alterIndex_allocate_extent_clause | alterIndex_deallocate_unused_clause)
    ;
alterIndex_annotation
    : ID STRING?
    ;
alterIndex_prefix_compression
    : 'COMPRESS' NUMBER?
    ;
alterIndex_advanced_index_compression
    : 'COMPRESS' 'ADVANCED' ('LOW' | 'HIGH')?
    ;
alterIndex_index_partition_description
    : 'PARTITION' ID? alterIndex_segment_attributes_clause? alterIndex_index_compression? (('USABLE' | 'UNUSABLE'))?
    ;
alterIndex_segment_attributes_clause
    : alterIndex_physical_attributes_clause? ('TABLESPACE' ID)? alterIndex_logging_clause?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-INDEXTYPE.html
alterIndextype
    : 'ALTER' 'INDEXTYPE' ('IF' 'EXISTS')? (ID '.')? ID (('ADD' | 'DROP') (ID '.')? ID '(' DATA_TYPE (',' DATA_TYPE)* ')' (',' ('ADD' | 'DROP') (ID '.')? ID '(' DATA_TYPE (',' DATA_TYPE)* ')')* alterIndextype_using_type_clause? | 'COMPILE') ('WITH' 'LOCAL' 'PARTITION')? alterIndextype_storage_table_clause?
    ;
alterIndextype_using_type_clause
    : 'USING' (ID '.')? ID alterIndextype_array_dml_clause?
    ;
alterIndextype_storage_table_clause
    : 'WITH' 'SYSTEM' 'MANAGED' 'STORAGE' 'TABLES'
    | 'WITH' 'USER' 'MANAGED' 'STORAGE' 'TABLES'
    ;
alterIndextype_array_dml_clause
    : 'WITH' 'ARRAY' 'DML' ('(' (ID '.')? ID ',' (ID '.')? ID ')' (',' '(' (ID '.')? ID ',' (ID '.')? ID ')')?)*
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-INMEMORY-JOIN-GROUP.html
alterInmemoryJoinGroup
    : 'ALTER' 'INMEMORY' 'JOIN' 'GROUP' ('IF' 'EXISTS')? (ID '.')? ID ('ADD' | 'REMOVE') '(' (ID '.')? ID '(' ID ')' ')'
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-JAVA.html
alterJava
    : 'ALTER' 'JAVA' ('IF' 'EXISTS')? ('SOURCE' | 'CLASS') (ID '.')? ID ('RESOLVER' '(' '(' STRING ID ')' (',' '(' STRING ID ')')* ')')? alterJava_invoker_rights_clause? ('RESOLVE' | 'COMPILE')
    ;
alterJava_invoker_rights_clause
    : 'AUTHID' ('CURRENT_USER' | 'DEFINER')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-LIBRARY.html
alterLibrary
    : 'ALTER' 'LIBRARY' ('IF' 'EXISTS')? (ID '.')? ID alterLibrary_library_compile_clause ('EDITIONABLE' | 'NONEDITIONABLE')?
    ;
alterLibrary_library_compile_clause
    : 'COMPILE' 'DEBUG'? alterLibrary_compiler_parameters_clause? ('REUSE' 'SETTINGS')?
    ;
alterLibrary_compiler_parameters_clause
    : ID '=' EXPR
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-LOCKDOWN-PROFILE.html
alterLockdownProfile
    : 'ALTER' 'LOCKDOWN' 'PROFILE' ID (alterLockdownProfile_lockdown_features | alterLockdownProfile_lockdown_options | alterLockdownProfile_lockdown_statements) ('USERS' '=' ('ALL' | 'COMMON' | 'LOCAL'))?
    ;
alterLockdownProfile_lockdown_features
    : ('DISABLE' | 'ENABLE') 'FEATURE' '=' '(' (ID | 'ALL' | 'ALL' 'EXCEPT' '=' '(' ID (',' ID)* ')') ')'
    ;
alterLockdownProfile_lockdown_options
    : ('DISABLE' | 'ENABLE') 'OPTION' '=' '(' (ID | 'ALL' | 'ALL' 'EXCEPT' '=' '(' ID (',' ID)* ')') ')'
    ;
alterLockdownProfile_lockdown_statements
    : ('DISABLE' | 'ENABLE') 'STATEMENT' '=' '(' STRING (',' STRING)* ')' alterLockdownProfile_statement_clauses?
    ;
alterLockdownProfile_statement_clauses
    : 'CLAUSE' '=' '(' (ID (',' ID)* | 'ALL' | 'ALL' 'EXCEPT' '=' '(' ID (',' ID)* ')') ')' alterLockdownProfile_clause_options?
    ;
alterLockdownProfile_clause_options
    : 'OPTION' '=' '(' ID (',' ID)* ')' alterLockdownProfile_option_values?
    ;
alterLockdownProfile_option_values
    : ('VALUE' '=' '(' EXPR (',' EXPR)* ')' | 'MINVALUE' '=' EXPR | 'MAXVALUE' '=' EXPR | ('MINVALUE' '=' EXPR)? ('MAXVALUE' '=' EXPR)?)
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-MATERIALIZED-VIEW-LOG.html
// fix: left out where optional, since no quoted page defines it or what it needs: external_part_subpart_data_props, lob_parameters, segment_attributes, storage_clause
alterMaterializedViewLog
    : 'ALTER' 'MATERIALIZED' 'VIEW' 'LOG' ('IF' 'EXISTS')? 'FORCE'? 'ON' (ID '.')? ID (alterMaterializedViewLog_physical_attributes_clause | alterMaterializedViewLog_add_mv_log_column_clause | alterMaterializedViewLog_alter_table_partitioning | alterMaterializedViewLog_parallel_clause | alterMaterializedViewLog_logging_clause | alterMaterializedViewLog_allocate_extent_clause | alterMaterializedViewLog_shrink_clause | alterMaterializedViewLog_move_mv_log_clause | alterMaterializedViewLog_mv_log_augmentation | alterMaterializedViewLog_mv_log_purge_clause | alterMaterializedViewLog_for_refresh_clause)
    ;
alterMaterializedViewLog_physical_attributes_clause
    : ('PCTFREE' NUMBER)? ('PCTUSED' NUMBER)? ('INITRANS' NUMBER)?
    ;
alterMaterializedViewLog_add_mv_log_column_clause
    : 'ADD' '(' ID (',' ID)* ')'
    ;
alterMaterializedViewLog_alter_table_partitioning
    : (alterMaterializedViewLog_modify_table_default_attrs | alterMaterializedViewLog_alter_automatic_partitioning | alterMaterializedViewLog_alter_interval_partitioning | alterMaterializedViewLog_set_subpartition_template | alterMaterializedViewLog_modify_table_partition | alterMaterializedViewLog_modify_table_subpartition | alterMaterializedViewLog_move_table_partition | alterMaterializedViewLog_move_table_subpartition | alterMaterializedViewLog_add_table_partition | alterMaterializedViewLog_coalesce_table_partition | alterMaterializedViewLog_coalesce_table_subpartition | alterMaterializedViewLog_drop_table_partition | alterMaterializedViewLog_drop_table_subpartition | alterMaterializedViewLog_rename_partition_subpart | alterMaterializedViewLog_truncate_partition_subpart | alterMaterializedViewLog_split_table_partition | alterMaterializedViewLog_split_table_subpartition | alterMaterializedViewLog_merge_table_partitions | alterMaterializedViewLog_merge_table_subpartitions | alterMaterializedViewLog_exchange_partition_subpart)
    ;
alterMaterializedViewLog_parallel_clause
    : ('PARALLEL' NUMBER? | 'NOPARALLEL')
    ;
alterMaterializedViewLog_logging_clause
    : ('LOGGING' | 'NOLOGGING' | 'FILESYSTEM_LIKE_LOGGING')
    ;
alterMaterializedViewLog_allocate_extent_clause
    : 'ALLOCATE' 'EXTENT' ('(' ('SIZE' alterMaterializedViewLog_size_clause)? ('DATAFILE' STRING)? ('INSTANCE' NUMBER)? ')')?
    ;
alterMaterializedViewLog_shrink_clause
    : 'SHRINK' 'SPACE' 'COMPACT'? 'CASCADE'?
    ;
alterMaterializedViewLog_move_mv_log_clause
    : 'MOVE' ('TABLESPACE' ID)?
    ;
alterMaterializedViewLog_mv_log_augmentation
    : 'ADD' ('PRIMARY' 'KEY' | 'ROWID' | 'OBJECT' 'ID' | 'SEQUENCE' | '(' ID (',' ID)* ')') alterMaterializedViewLog_new_values_clause?
    ;
alterMaterializedViewLog_mv_log_purge_clause
    : 'PURGE' ('IMMEDIATE' 'SYNCHRONOUS' | 'IMMEDIATE' 'ASYNCHRONOUS' | 'START' 'WITH' EXPR 'NEXT' EXPR)
    ;
alterMaterializedViewLog_for_refresh_clause
    : ('FOR' 'SYNCHRONOUS' 'REFRESH' | 'FOR' 'FAST' 'REFRESH')
    ;
alterMaterializedViewLog_modify_table_default_attrs
    : 'MODIFY' 'DEFAULT' 'ATTRIBUTES' ('FOR' 'PARTITION' ID)? alterMaterializedViewLog_deferred_segment_creation? alterMaterializedViewLog_read_only_clause? alterMaterializedViewLog_indexing_clause? alterMaterializedViewLog_segment_attributes_clause? alterMaterializedViewLog_table_compression? alterMaterializedViewLog_inmemory_clause? alterMaterializedViewLog_prefix_compression? ('OVERFLOW' alterMaterializedViewLog_segment_attributes_clause)?
    ;
alterMaterializedViewLog_alter_automatic_partitioning
    : 'SET' 'PARTITIONING' ('AUTOMATIC' | 'MANUAL')
    ;
alterMaterializedViewLog_alter_interval_partitioning
    : 'SET' 'INTERVAL' '(' EXPR? ')'
    ;
alterMaterializedViewLog_set_subpartition_template
    : 'SET' 'SUBPARTITION' 'TEMPLATE' ('(' (alterMaterializedViewLog_range_subpartition_desc | alterMaterializedViewLog_list_subpartition_desc | alterMaterializedViewLog_individual_hash_subparts) (',' (alterMaterializedViewLog_range_subpartition_desc | alterMaterializedViewLog_list_subpartition_desc | alterMaterializedViewLog_individual_hash_subparts))* ')' | alterMaterializedViewLog_hash_subparts_by_quantity)
    ;
alterMaterializedViewLog_modify_table_partition
    : 'MODIFY' 'PARTITION' ID (alterMaterializedViewLog_partition_attributes | (alterMaterializedViewLog_add_range_subpartition | alterMaterializedViewLog_add_hash_subpartition | alterMaterializedViewLog_add_list_subpartition) | 'COALESCE' 'SUBPARTITION' alterMaterializedViewLog_update_index_clauses? alterMaterializedViewLog_parallel_clause? | 'REBUILD'? 'UNUSABLE' 'LOCAL' 'INDEXES' | alterMaterializedViewLog_read_only_clause | alterMaterializedViewLog_indexing_clause)
    ;
alterMaterializedViewLog_modify_table_subpartition
    : 'MODIFY' 'SUBPARTITION' ID (alterMaterializedViewLog_allocate_extent_clause | alterMaterializedViewLog_deallocate_unused_clause | 'SHRINK' 'SPACE' 'CASCADE'? | ('LOB' '(' ID ')' | 'VARRAY' ID) '(' alterMaterializedViewLog_modify_lob_parameters ')' | 'REBUILD'? 'UNUSABLE' 'LOCAL' 'INDEXES' | ('ADD' | 'DROP') 'VALUES' '(' EXPR (',' EXPR)* ')' | alterMaterializedViewLog_read_only_clause | alterMaterializedViewLog_indexing_clause)
    ;
alterMaterializedViewLog_move_table_partition
    : 'MOVE' 'PARTITION' ID ('MAPPING' 'TABLE')? alterMaterializedViewLog_table_partition_description? alterMaterializedViewLog_filter_condition? alterMaterializedViewLog_update_index_clauses? alterMaterializedViewLog_parallel_clause? alterMaterializedViewLog_allow_disallow_clustering? 'ONLINE'?
    ;
alterMaterializedViewLog_move_table_subpartition
    : 'MOVE' 'SUBPARTITION' ID alterMaterializedViewLog_indexing_clause? alterMaterializedViewLog_partitioning_storage_clause? alterMaterializedViewLog_update_index_clauses? alterMaterializedViewLog_filter_condition? alterMaterializedViewLog_parallel_clause? alterMaterializedViewLog_allow_disallow_clustering? 'ONLINE'?
    ;
alterMaterializedViewLog_add_table_partition
    : 'ADD' (alterMaterializedViewLog_add_range_partition_clause | alterMaterializedViewLog_add_list_partition_clause | alterMaterializedViewLog_add_system_partition_clause | alterMaterializedViewLog_add_hash_partition_clause)
    ;
alterMaterializedViewLog_coalesce_table_partition
    : 'COALESCE' 'PARTITION' alterMaterializedViewLog_update_index_clauses? alterMaterializedViewLog_parallel_clause? alterMaterializedViewLog_allow_disallow_clustering?
    ;
alterMaterializedViewLog_coalesce_table_subpartition
    : 'COALESCE' 'SUBPARTITION' ID? alterMaterializedViewLog_update_index_clauses? alterMaterializedViewLog_parallel_clause? alterMaterializedViewLog_allow_disallow_clustering?
    ;
alterMaterializedViewLog_drop_table_partition
    : 'DROP' 'PARTITION' (ID | 'FOR' '(' EXPR (',' EXPR)* ')') (',' (ID | 'FOR' '(' EXPR (',' EXPR)* ')'))* alterMaterializedViewLog_update_index_clauses? alterMaterializedViewLog_parallel_clause?
    ;
alterMaterializedViewLog_drop_table_subpartition
    : 'DROP' 'SUBPARTITION' (ID | 'FOR' '(' EXPR (',' EXPR)* ')') (',' (ID | 'FOR' '(' EXPR (',' EXPR)* ')'))* alterMaterializedViewLog_update_index_clauses? alterMaterializedViewLog_parallel_clause?
    ;
alterMaterializedViewLog_rename_partition_subpart
    : 'RENAME' ('PARTITION' ID | 'SUBPARTITION' ID) 'TO' ID
    ;
alterMaterializedViewLog_truncate_partition_subpart
    : 'TRUNCATE' ('PARTITION' | 'SUBPARTITION') (ID | 'FOR' '(' EXPR (',' EXPR)* ')') (',' (ID | 'FOR' '(' EXPR (',' EXPR)* ')'))* (('DROP' 'ALL'? | 'REUSE') 'STORAGE')? alterMaterializedViewLog_update_index_clauses? alterMaterializedViewLog_parallel_clause? 'CASCADE'?
    ;
alterMaterializedViewLog_split_table_partition
    : 'SPLIT' 'PARTITION' ID ('AT' '(' EXPR (',' EXPR)* ')' ('INTO' '(' alterMaterializedViewLog_partition_spec ',' alterMaterializedViewLog_partition_spec ')')? | 'VALUES' '(' EXPR (',' EXPR)* ')' ('INTO' '(' alterMaterializedViewLog_partition_spec ',' alterMaterializedViewLog_partition_spec ')')? | 'INTO' '(' alterMaterializedViewLog_partition_spec (',' alterMaterializedViewLog_partition_spec)* ',' alterMaterializedViewLog_partition_spec ')') alterMaterializedViewLog_split_nested_table_part? alterMaterializedViewLog_filter_condition? alterMaterializedViewLog_dependent_tables_clause? alterMaterializedViewLog_update_index_clauses? alterMaterializedViewLog_parallel_clause? alterMaterializedViewLog_allow_disallow_clustering? 'ONLINE'?
    ;
alterMaterializedViewLog_split_table_subpartition
    : 'SPLIT' 'SUBPARTITION' ID ('AT' '(' EXPR (',' EXPR)* ')' ('INTO' '(' alterMaterializedViewLog_subpartition_spec ',' alterMaterializedViewLog_subpartition_spec ')')? | 'VALUES' '(' EXPR (',' EXPR)* ')' ('INTO' '(' alterMaterializedViewLog_subpartition_spec ',' alterMaterializedViewLog_subpartition_spec ')')? | 'INTO' '(' alterMaterializedViewLog_subpartition_spec (',' alterMaterializedViewLog_subpartition_spec)* ',' alterMaterializedViewLog_subpartition_spec ')') alterMaterializedViewLog_filter_condition? alterMaterializedViewLog_dependent_tables_clause? alterMaterializedViewLog_update_index_clauses? alterMaterializedViewLog_parallel_clause? alterMaterializedViewLog_allow_disallow_clustering? 'ONLINE'?
    ;
alterMaterializedViewLog_merge_table_partitions
    : 'MERGE' 'PARTITIONS' (alterMaterializedViewLog_partition_or_key_value ',' alterMaterializedViewLog_partition_or_key_value ('INTO' alterMaterializedViewLog_partition_spec)? | alterMaterializedViewLog_partition_or_key_value 'TO' alterMaterializedViewLog_partition_or_key_value ('INTO' alterMaterializedViewLog_partition_spec)?) alterMaterializedViewLog_filter_condition? alterMaterializedViewLog_dependent_tables_clause? alterMaterializedViewLog_update_index_clauses? alterMaterializedViewLog_parallel_clause? alterMaterializedViewLog_allow_disallow_clustering? 'ONLINE'?
    ;
alterMaterializedViewLog_merge_table_subpartitions
    : 'MERGE' 'SUBPARTITIONS' (alterMaterializedViewLog_subpartition_or_key_value ',' alterMaterializedViewLog_subpartition_or_key_value ('INTO' alterMaterializedViewLog_subpartition_spec)? | alterMaterializedViewLog_subpartition_or_key_value 'TO' alterMaterializedViewLog_subpartition_or_key_value ('INTO' alterMaterializedViewLog_subpartition_spec)?) alterMaterializedViewLog_filter_condition? alterMaterializedViewLog_dependent_tables_clause? alterMaterializedViewLog_update_index_clauses? alterMaterializedViewLog_parallel_clause? alterMaterializedViewLog_allow_disallow_clustering? 'ONLINE'?
    ;
alterMaterializedViewLog_exchange_partition_subpart
    : 'EXCHANGE' ('PARTITION' ID | 'SUBPARTITION' ID) 'WITH' 'TABLE' (ID '.')? ID (('INCLUDING' | 'EXCLUDING') 'INDEXES')? (('WITH' | 'WITHOUT') 'VALIDATION')? alterMaterializedViewLog_exceptions_clause? alterMaterializedViewLog_update_index_clauses? alterMaterializedViewLog_parallel_clause? 'CASCADE'? 'ONLINE'?
    ;
alterMaterializedViewLog_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T' | 'P')?
    ;
alterMaterializedViewLog_new_values_clause
    : ('INCLUDING' 'NEW' 'VALUES' | 'EXCLUDING' 'NEW' 'VALUES')
    ;
alterMaterializedViewLog_deferred_segment_creation
    : 'SEGMENT' 'CREATION' ('IMMEDIATE' | 'DEFERRED')
    ;
alterMaterializedViewLog_read_only_clause
    : ('READ' 'ONLY' | 'READ' 'WRITE')
    ;
alterMaterializedViewLog_indexing_clause
    : 'INDEXING' ('ON' | 'OFF')
    ;
alterMaterializedViewLog_segment_attributes_clause
    : alterMaterializedViewLog_physical_attributes_clause? ('TABLESPACE' ID)? alterMaterializedViewLog_logging_clause?
    ;
alterMaterializedViewLog_table_compression
    : ('COMPRESS' 'BASIC'? | 'COMPRESS' 'FOR' ('OLTP' | 'QUERY' ('LOW' | 'HIGH')? | 'ARCHIVE' ('LOW' | 'HIGH')?) | 'NOCOMPRESS')
    ;
alterMaterializedViewLog_inmemory_clause
    : ('INMEMORY' alterMaterializedViewLog_inmemory_attributes? | 'NO' 'INMEMORY')
    ;
alterMaterializedViewLog_prefix_compression
    : 'COMPRESS' NUMBER?
    ;
alterMaterializedViewLog_range_subpartition_desc
    : 'SUBPARTITION' ID? alterMaterializedViewLog_range_values_clause alterMaterializedViewLog_read_only_clause? alterMaterializedViewLog_indexing_clause? alterMaterializedViewLog_partitioning_storage_clause?
    ;
alterMaterializedViewLog_list_subpartition_desc
    : 'SUBPARTITION' ID? alterMaterializedViewLog_list_values_clause alterMaterializedViewLog_read_only_clause? alterMaterializedViewLog_indexing_clause? alterMaterializedViewLog_partitioning_storage_clause?
    ;
alterMaterializedViewLog_individual_hash_subparts
    : 'SUBPARTITION' ID? alterMaterializedViewLog_read_only_clause? alterMaterializedViewLog_indexing_clause? alterMaterializedViewLog_partitioning_storage_clause?
    ;
alterMaterializedViewLog_hash_subparts_by_quantity
    : 'SUBPARTITIONS' NUMBER ('STORE' 'IN' '(' ID (',' ID)* ')')?
    ;
alterMaterializedViewLog_partition_attributes
    : alterMaterializedViewLog_physical_attributes_clause? alterMaterializedViewLog_logging_clause? alterMaterializedViewLog_allocate_extent_clause? alterMaterializedViewLog_deallocate_unused_clause? ('SHRINK' 'SPACE' 'CASCADE'?)? ((alterMaterializedViewLog_table_compression | alterMaterializedViewLog_prefix_compression))? alterMaterializedViewLog_inmemory_clause? (('LOB' '(' ID ')' | 'VARRAY' ID) '(' alterMaterializedViewLog_modify_lob_parameters ')')* ('OVERFLOW' (alterMaterializedViewLog_physical_attributes_clause | alterMaterializedViewLog_logging_clause | alterMaterializedViewLog_allocate_extent_clause | alterMaterializedViewLog_deallocate_unused_clause | 'SHRINK' 'SPACE' 'CASCADE'?))?
    ;
alterMaterializedViewLog_add_range_subpartition
    : 'ADD' alterMaterializedViewLog_range_subpartition_desc (',' alterMaterializedViewLog_range_subpartition_desc)* alterMaterializedViewLog_dependent_tables_clause? alterMaterializedViewLog_update_index_clauses?
    ;
alterMaterializedViewLog_add_hash_subpartition
    : 'ADD' alterMaterializedViewLog_individual_hash_subparts (',' alterMaterializedViewLog_individual_hash_subparts)* alterMaterializedViewLog_dependent_tables_clause? alterMaterializedViewLog_update_index_clauses? alterMaterializedViewLog_parallel_clause?
    ;
alterMaterializedViewLog_add_list_subpartition
    : 'ADD' alterMaterializedViewLog_list_subpartition_desc (',' alterMaterializedViewLog_list_subpartition_desc)* alterMaterializedViewLog_dependent_tables_clause? alterMaterializedViewLog_update_index_clauses?
    ;
alterMaterializedViewLog_update_index_clauses
    : (alterMaterializedViewLog_update_global_index_clause | alterMaterializedViewLog_update_all_indexes_clause)
    ;
alterMaterializedViewLog_deallocate_unused_clause
    : 'DEALLOCATE' 'UNUSED' ('KEEP' alterMaterializedViewLog_size_clause)?
    ;
alterMaterializedViewLog_modify_lob_parameters
    : alterMaterializedViewLog_lob_retention_clause? alterMaterializedViewLog_lob_compression_clause? alterMaterializedViewLog_logging_clause? alterMaterializedViewLog_allocate_extent_clause? alterMaterializedViewLog_shrink_clause? alterMaterializedViewLog_deallocate_unused_clause? (('ENCRYPT' ID | 'DECRYPT'))?
    ;
alterMaterializedViewLog_table_partition_description
    : alterMaterializedViewLog_deferred_segment_creation? alterMaterializedViewLog_read_only_clause? alterMaterializedViewLog_indexing_clause? alterMaterializedViewLog_segment_attributes_clause? (alterMaterializedViewLog_table_compression | alterMaterializedViewLog_prefix_compression)? alterMaterializedViewLog_inmemory_clause? alterMaterializedViewLog_ilm_clause? ('OVERFLOW' alterMaterializedViewLog_segment_attributes_clause?)?
    ;
alterMaterializedViewLog_filter_condition
    : 'INCLUDING' 'ROWS' 'WHERE' CONDITION
    ;
alterMaterializedViewLog_allow_disallow_clustering
    : ('ALLOW' | 'DISALLOW') 'CLUSTERING'
    ;
alterMaterializedViewLog_partitioning_storage_clause
    : ('TABLESPACE' ID)? ('OVERFLOW' 'TABLESPACE' ID)? alterMaterializedViewLog_table_compression? alterMaterializedViewLog_inmemory_clause?
    ;
alterMaterializedViewLog_add_range_partition_clause
    : 'PARTITION' ID? alterMaterializedViewLog_range_values_clause alterMaterializedViewLog_table_partition_description? ('(' (alterMaterializedViewLog_range_subpartition_desc | alterMaterializedViewLog_list_subpartition_desc | alterMaterializedViewLog_individual_hash_subparts) (',' (alterMaterializedViewLog_range_subpartition_desc | alterMaterializedViewLog_list_subpartition_desc | alterMaterializedViewLog_individual_hash_subparts))* ')' | alterMaterializedViewLog_hash_subparts_by_quantity)? alterMaterializedViewLog_update_index_clauses?
    ;
alterMaterializedViewLog_add_list_partition_clause
    : 'PARTITION' ID? alterMaterializedViewLog_list_values_clause alterMaterializedViewLog_table_partition_description? ('(' (alterMaterializedViewLog_range_subpartition_desc | alterMaterializedViewLog_list_subpartition_desc | alterMaterializedViewLog_individual_hash_subparts) (',' (alterMaterializedViewLog_range_subpartition_desc | alterMaterializedViewLog_list_subpartition_desc | alterMaterializedViewLog_individual_hash_subparts))* ')' | alterMaterializedViewLog_hash_subparts_by_quantity)? alterMaterializedViewLog_update_index_clauses?
    ;
alterMaterializedViewLog_add_system_partition_clause
    : 'PARTITION' ID? alterMaterializedViewLog_table_partition_description? alterMaterializedViewLog_update_index_clauses?
    ;
alterMaterializedViewLog_add_hash_partition_clause
    : 'PARTITION' ID? alterMaterializedViewLog_partitioning_storage_clause? alterMaterializedViewLog_update_index_clauses? alterMaterializedViewLog_parallel_clause? alterMaterializedViewLog_read_only_clause? alterMaterializedViewLog_indexing_clause?
    ;
alterMaterializedViewLog_partition_spec
    : 'PARTITION' ID? alterMaterializedViewLog_table_partition_description?
    ;
alterMaterializedViewLog_split_nested_table_part
    : 'NESTED' 'TABLE' ID 'INTO' '(' alterMaterializedViewLog_nested_table_partition_spec ',' alterMaterializedViewLog_nested_table_partition_spec ')'
    ;
alterMaterializedViewLog_dependent_tables_clause
    : 'DEPENDENT' 'TABLES' '(' ID '(' alterMaterializedViewLog_partition_spec (',' alterMaterializedViewLog_partition_spec)* ')' (',' ID '(' alterMaterializedViewLog_partition_spec (',' alterMaterializedViewLog_partition_spec)* ')')* ')'
    ;
alterMaterializedViewLog_subpartition_spec
    : 'SUBPARTITION' ID? alterMaterializedViewLog_partitioning_storage_clause?
    ;
alterMaterializedViewLog_partition_or_key_value
    : ID
    | 'FOR' '(' EXPR (',' EXPR)* ')'
    ;
alterMaterializedViewLog_subpartition_or_key_value
    : ID
    | 'FOR' '(' EXPR (',' EXPR)* ')'
    ;
alterMaterializedViewLog_exceptions_clause
    : 'EXCEPTIONS' 'INTO' (ID '.')? ID
    ;
alterMaterializedViewLog_inmemory_attributes
    : alterMaterializedViewLog_inmemory_memcompress? alterMaterializedViewLog_inmemory_priority? alterMaterializedViewLog_inmemory_distribute? alterMaterializedViewLog_inmemory_duplicate?
    ;
alterMaterializedViewLog_range_values_clause
    : 'VALUES' 'LESS' 'THAN' '(' (EXPR | 'MAXVALUE') (',' (EXPR | 'MAXVALUE'))* ')'
    ;
alterMaterializedViewLog_list_values_clause
    : 'VALUES' '(' (EXPR (',' EXPR)* | 'DEFAULT') ')'
    ;
alterMaterializedViewLog_update_global_index_clause
    : ('UPDATE' | 'INVALIDATE') 'GLOBAL' 'INDEXES'
    ;
alterMaterializedViewLog_update_all_indexes_clause
    : 'UPDATE' 'INDEXES' ('(' ID '(' (alterMaterializedViewLog_index_partition_description | alterMaterializedViewLog_index_subpartition_clause) (',' (alterMaterializedViewLog_index_partition_description | alterMaterializedViewLog_index_subpartition_clause))* ')' (',' ID '(' (alterMaterializedViewLog_index_partition_description | alterMaterializedViewLog_index_subpartition_clause) (',' (alterMaterializedViewLog_index_partition_description | alterMaterializedViewLog_index_subpartition_clause))* ')')* ')')?
    ;
alterMaterializedViewLog_lob_retention_clause
    : 'RETENTION' ('MAX' | 'MIN' NUMBER | 'AUTO' | 'NONE')?
    ;
alterMaterializedViewLog_lob_compression_clause
    : ('COMPRESS' ('HIGH' | 'MEDIUM' | 'LOW')? | 'NOCOMPRESS')
    ;
alterMaterializedViewLog_ilm_clause
    : 'ILM' ('ADD' 'POLICY' alterMaterializedViewLog_ilm_policy_clause | 'DELETE' 'POLICY' ID | 'DELETE_ALL' | 'ENABLE_ALL' | 'DISABLE_ALL')
    ;
alterMaterializedViewLog_nested_table_partition_spec
    : 'PARTITION' ID alterMaterializedViewLog_segment_attributes_clause?
    ;
alterMaterializedViewLog_inmemory_memcompress
    : 'MEMCOMPRESS' 'FOR' ('QUERY' ('LOW' | 'HIGH')? | 'OLTP' | 'ARCHIVE' ('LOW' | 'HIGH')?)
    ;
alterMaterializedViewLog_inmemory_priority
    : 'PRIORITY' ('LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL')
    ;
alterMaterializedViewLog_inmemory_distribute
    : 'DISTRIBUTE' 'BY' ('ROWID' ('ACROSS' NUMBER?)? | 'PARTITION' | 'SUBPARTITION')
    ;
alterMaterializedViewLog_inmemory_duplicate
    : 'DUPLICATE' ('ALL' | 'FOR' 'SERVICE' ID (',' ID)* | 'NONE')
    ;
alterMaterializedViewLog_index_partition_description
    : 'PARTITION' ID? alterMaterializedViewLog_segment_attributes_clause? alterMaterializedViewLog_index_compression? (('USABLE' | 'UNUSABLE'))?
    ;
alterMaterializedViewLog_index_subpartition_clause
    : '(' 'SUBPARTITION' ID? ('TABLESPACE' ID)? alterMaterializedViewLog_index_compression? (('USABLE' | 'UNUSABLE'))? (',' 'SUBPARTITION' ID? ('TABLESPACE' ID)? alterMaterializedViewLog_index_compression? (('USABLE' | 'UNUSABLE'))?)* ')'
    | 'STORE' 'IN' '(' ID (',' ID)* ')'
    ;
alterMaterializedViewLog_ilm_policy_clause
    : (alterMaterializedViewLog_ilm_compression_policy | alterMaterializedViewLog_ilm_tiering_policy | alterMaterializedViewLog_ilm_inmemory_policy)
    ;
alterMaterializedViewLog_index_compression
    : (alterMaterializedViewLog_prefix_compression | alterMaterializedViewLog_advanced_index_compression | 'NOCOMPRESS')
    ;
alterMaterializedViewLog_ilm_compression_policy
    : (alterMaterializedViewLog_table_compression | 'ROW' 'STORE' 'COMPRESS' 'ADVANCED'? | 'COLUMN' 'STORE' 'COMPRESS' ('FOR' 'QUERY' ('LOW' | 'HIGH')?)?) ('SEGMENT' | 'GROUP' | 'ROW') 'AFTER' alterMaterializedViewLog_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION')
    ;
alterMaterializedViewLog_ilm_tiering_policy
    : 'TIER' 'TO' ID ('SEGMENT' | 'GROUP' | 'ROW')? ('ON' ID)?
    ;
alterMaterializedViewLog_ilm_inmemory_policy
    : ('SET' 'INMEMORY' alterMaterializedViewLog_inmemory_memcompress? alterMaterializedViewLog_inmemory_priority? | 'MODIFY' 'INMEMORY' alterMaterializedViewLog_inmemory_memcompress | 'NO' 'INMEMORY') 'SEGMENT' 'AFTER' alterMaterializedViewLog_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION')
    ;
alterMaterializedViewLog_advanced_index_compression
    : 'COMPRESS' 'ADVANCED' ('LOW' | 'HIGH')?
    ;
alterMaterializedViewLog_ilm_time_period
    : NUMBER ('DAY' | 'DAYS' | 'MONTH' | 'MONTHS' | 'YEAR' | 'YEARS')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-MATERIALIZED-VIEW.html
// fix: left out where optional, since no quoted page defines it or what it needs: external_part_subpart_data_props, lob_parameters, storage_clause
alterMaterializedView
    : 'ALTER' 'MATERIALIZED' 'VIEW' ('IF' 'EXISTS')? (ID '.')? ID (alterMaterializedView_physical_attributes_clause | alterMaterializedView_modify_mv_column_clause | alterMaterializedView_table_compression | alterMaterializedView_inmemory_table_clause | alterMaterializedView_lob_storage_clause | alterMaterializedView_modify_lob_storage_clause | alterMaterializedView_alter_table_partitioning | alterMaterializedView_parallel_clause | alterMaterializedView_logging_clause | alterMaterializedView_allocate_extent_clause | alterMaterializedView_deallocate_unused_clause | alterMaterializedView_shrink_clause | ('CACHE' | 'NOCACHE') | alterMaterializedView_annotations_clause | alterMaterializedView_alter_iot_clauses | 'MODIFY' alterMaterializedView_scoped_table_ref_constraint | alterMaterializedView_alter_mv_refresh | alterMaterializedView_evaluation_edition_clause | ('ENABLE' | 'DISABLE') 'ON' 'QUERY' 'COMPUTATION' | alterMaterializedView_alter_query_rewrite_clause | 'COMPILE' | 'CONSIDER' 'FRESH')
    ;
alterMaterializedView_physical_attributes_clause
    : ('PCTFREE' NUMBER)? ('PCTUSED' NUMBER)? ('INITRANS' NUMBER)?
    ;
alterMaterializedView_modify_mv_column_clause
    : 'MODIFY' '(' ID alterMaterializedView_encryption_spec (',' ID alterMaterializedView_encryption_spec)* ')'
    ;
alterMaterializedView_table_compression
    : ('COMPRESS' 'BASIC'? | 'COMPRESS' 'FOR' ('OLTP' | 'QUERY' ('LOW' | 'HIGH')? | 'ARCHIVE' ('LOW' | 'HIGH')?) | 'NOCOMPRESS')
    ;
alterMaterializedView_inmemory_table_clause
    : ('INMEMORY' alterMaterializedView_inmemory_attributes? alterMaterializedView_inmemory_column_clause? | 'NO' 'INMEMORY')
    ;
alterMaterializedView_lob_storage_clause
    : 'LOB' '(' ID (',' ID)* ')' 'STORE' 'AS' ID? '(' alterMaterializedView_lob_storage_parameters ')'
    ;
alterMaterializedView_modify_lob_storage_clause
    : 'MODIFY' 'LOB' '(' ID (',' ID)* ')' '(' alterMaterializedView_modify_lob_parameters ')'
    ;
alterMaterializedView_alter_table_partitioning
    : (alterMaterializedView_modify_table_default_attrs | alterMaterializedView_alter_automatic_partitioning | alterMaterializedView_alter_interval_partitioning | alterMaterializedView_set_subpartition_template | alterMaterializedView_modify_table_partition | alterMaterializedView_modify_table_subpartition | alterMaterializedView_move_table_partition | alterMaterializedView_move_table_subpartition | alterMaterializedView_add_table_partition | alterMaterializedView_coalesce_table_partition | alterMaterializedView_coalesce_table_subpartition | alterMaterializedView_drop_table_partition | alterMaterializedView_drop_table_subpartition | alterMaterializedView_rename_partition_subpart | alterMaterializedView_truncate_partition_subpart | alterMaterializedView_split_table_partition | alterMaterializedView_split_table_subpartition | alterMaterializedView_merge_table_partitions | alterMaterializedView_merge_table_subpartitions | alterMaterializedView_exchange_partition_subpart)
    ;
alterMaterializedView_parallel_clause
    : ('PARALLEL' NUMBER? | 'NOPARALLEL')
    ;
alterMaterializedView_logging_clause
    : ('LOGGING' | 'NOLOGGING')
    ;
alterMaterializedView_allocate_extent_clause
    : 'ALLOCATE' 'EXTENT' ('(' ('SIZE' alterMaterializedView_size_clause)? ('DATAFILE' STRING)? ('INSTANCE' NUMBER)? ')')?
    ;
alterMaterializedView_deallocate_unused_clause
    : 'DEALLOCATE' 'UNUSED' ('KEEP' alterMaterializedView_size_clause)?
    ;
alterMaterializedView_shrink_clause
    : 'SHRINK' 'SPACE' 'COMPACT'? 'CASCADE'?
    ;
alterMaterializedView_annotations_clause
    : 'ANNOTATIONS' '(' ('ADD' alterMaterializedView_annotation | 'DROP' ID (',' ID)* | 'REMOVE' 'ALL') ')'
    ;
alterMaterializedView_alter_iot_clauses
    : (alterMaterializedView_index_org_table_clause | alterMaterializedView_alter_overflow_clause)
    ;
alterMaterializedView_scoped_table_ref_constraint
    : 'SCOPE' 'FOR' '(' ID ')' 'IS' (ID '.')? ID
    ;
alterMaterializedView_alter_mv_refresh
    : 'REFRESH' ('FAST' | 'COMPLETE' | 'FORCE') (('ON' 'COMMIT' | 'ON' 'DEMAND'))? ('START' 'WITH' EXPR)? ('NEXT' EXPR)? ('USING' 'ROLLBACK' 'SEGMENT' ID)? ('WITH' 'PRIMARY' 'KEY')? (('ENABLE' | 'DISABLE') 'CONCURRENT' 'REFRESH')?
    ;
alterMaterializedView_evaluation_edition_clause
    : 'EVALUATE' 'USING' 'EDITION' ID ('CONSIDER' 'FRESH')?
    ;
alterMaterializedView_alter_query_rewrite_clause
    : 'QUERY' 'REWRITE' ('ENABLE' alterMaterializedView_unusable_editions_clause? | 'DISABLE' | alterMaterializedView_unusable_editions_clause)
    ;
alterMaterializedView_encryption_spec
    : ('USING' STRING)? ('IDENTIFIED' 'BY' ID)? STRING? ('NO'? 'SALT')?
    ;
alterMaterializedView_inmemory_attributes
    : alterMaterializedView_inmemory_memcompress? alterMaterializedView_inmemory_priority? alterMaterializedView_inmemory_distribute? alterMaterializedView_inmemory_duplicate?
    ;
alterMaterializedView_inmemory_column_clause
    : 'INMEMORY' 'MEMCOMPRESS' 'FOR' ('QUERY' ('LOW' | 'HIGH')? | 'OLTP' | 'ARCHIVE' ('LOW' | 'HIGH')?) '(' ID (',' ID)* ')'
    ;
alterMaterializedView_lob_storage_parameters
    : (('TABLESPACE' ID | alterMaterializedView_logging_clause)+ | ('ENABLE' | 'DISABLE') 'STORAGE' 'IN' 'ROW' | 'CHUNK' NUMBER | 'PCTVERSION' NUMBER | 'FREEPOOLS' NUMBER | alterMaterializedView_lob_retention_clause | alterMaterializedView_lob_deduplicate_clause | alterMaterializedView_lob_compression_clause | ('ENCRYPT' alterMaterializedView_encryption_spec | 'DECRYPT') | ('CACHE' | 'NOCACHE' | 'CACHE' 'READS') alterMaterializedView_logging_clause?)+
    ;
alterMaterializedView_modify_lob_parameters
    : alterMaterializedView_lob_retention_clause? alterMaterializedView_lob_compression_clause? alterMaterializedView_logging_clause? alterMaterializedView_allocate_extent_clause? alterMaterializedView_shrink_clause? alterMaterializedView_deallocate_unused_clause? (('ENCRYPT' ID | 'DECRYPT'))?
    ;
alterMaterializedView_modify_table_default_attrs
    : 'MODIFY' 'DEFAULT' 'ATTRIBUTES' ('FOR' 'PARTITION' ID)? alterMaterializedView_deferred_segment_creation? alterMaterializedView_read_only_clause? alterMaterializedView_indexing_clause? alterMaterializedView_segment_attributes_clause? alterMaterializedView_table_compression? alterMaterializedView_inmemory_clause? alterMaterializedView_prefix_compression? ('OVERFLOW' alterMaterializedView_segment_attributes_clause)?
    ;
alterMaterializedView_alter_automatic_partitioning
    : 'SET' 'PARTITIONING' ('AUTOMATIC' | 'MANUAL')
    ;
alterMaterializedView_alter_interval_partitioning
    : 'SET' 'INTERVAL' '(' EXPR? ')'
    ;
alterMaterializedView_set_subpartition_template
    : 'SET' 'SUBPARTITION' 'TEMPLATE' ('(' (alterMaterializedView_range_subpartition_desc | alterMaterializedView_list_subpartition_desc | alterMaterializedView_individual_hash_subparts) (',' (alterMaterializedView_range_subpartition_desc | alterMaterializedView_list_subpartition_desc | alterMaterializedView_individual_hash_subparts))* ')' | alterMaterializedView_hash_subparts_by_quantity)
    ;
alterMaterializedView_modify_table_partition
    : 'MODIFY' 'PARTITION' ID (alterMaterializedView_partition_attributes | (alterMaterializedView_add_range_subpartition | alterMaterializedView_add_hash_subpartition | alterMaterializedView_add_list_subpartition) | 'COALESCE' 'SUBPARTITION' alterMaterializedView_update_index_clauses? alterMaterializedView_parallel_clause? | 'REBUILD'? 'UNUSABLE' 'LOCAL' 'INDEXES' | alterMaterializedView_read_only_clause | alterMaterializedView_indexing_clause)
    ;
alterMaterializedView_modify_table_subpartition
    : 'MODIFY' 'SUBPARTITION' ID (alterMaterializedView_allocate_extent_clause | alterMaterializedView_deallocate_unused_clause | 'SHRINK' 'SPACE' 'CASCADE'? | ('LOB' '(' ID ')' | 'VARRAY' ID) '(' alterMaterializedView_modify_lob_parameters ')' | 'REBUILD'? 'UNUSABLE' 'LOCAL' 'INDEXES' | ('ADD' | 'DROP') 'VALUES' '(' EXPR (',' EXPR)* ')' | alterMaterializedView_read_only_clause | alterMaterializedView_indexing_clause)
    ;
alterMaterializedView_move_table_partition
    : 'MOVE' 'PARTITION' ID ('MAPPING' 'TABLE')? alterMaterializedView_table_partition_description? alterMaterializedView_filter_condition? alterMaterializedView_update_index_clauses? alterMaterializedView_parallel_clause? alterMaterializedView_allow_disallow_clustering? 'ONLINE'?
    ;
alterMaterializedView_move_table_subpartition
    : 'MOVE' 'SUBPARTITION' ID alterMaterializedView_indexing_clause? alterMaterializedView_partitioning_storage_clause? alterMaterializedView_update_index_clauses? alterMaterializedView_filter_condition? alterMaterializedView_parallel_clause? alterMaterializedView_allow_disallow_clustering? 'ONLINE'?
    ;
alterMaterializedView_add_table_partition
    : 'ADD' (alterMaterializedView_add_range_partition_clause | alterMaterializedView_add_list_partition_clause | alterMaterializedView_add_system_partition_clause | alterMaterializedView_add_hash_partition_clause)
    ;
alterMaterializedView_coalesce_table_partition
    : 'COALESCE' 'PARTITION' alterMaterializedView_update_index_clauses? alterMaterializedView_parallel_clause? alterMaterializedView_allow_disallow_clustering?
    ;
alterMaterializedView_coalesce_table_subpartition
    : 'COALESCE' 'SUBPARTITION' ID? alterMaterializedView_update_index_clauses? alterMaterializedView_parallel_clause? alterMaterializedView_allow_disallow_clustering?
    ;
alterMaterializedView_drop_table_partition
    : 'DROP' 'PARTITION' (ID | 'FOR' '(' EXPR (',' EXPR)* ')') (',' (ID | 'FOR' '(' EXPR (',' EXPR)* ')'))* alterMaterializedView_update_index_clauses? alterMaterializedView_parallel_clause?
    ;
alterMaterializedView_drop_table_subpartition
    : 'DROP' 'SUBPARTITION' (ID | 'FOR' '(' EXPR (',' EXPR)* ')') (',' (ID | 'FOR' '(' EXPR (',' EXPR)* ')'))* alterMaterializedView_update_index_clauses? alterMaterializedView_parallel_clause?
    ;
alterMaterializedView_rename_partition_subpart
    : 'RENAME' ('PARTITION' ID | 'SUBPARTITION' ID) 'TO' ID
    ;
alterMaterializedView_truncate_partition_subpart
    : 'TRUNCATE' ('PARTITION' | 'SUBPARTITION') (ID | 'FOR' '(' EXPR (',' EXPR)* ')') (',' (ID | 'FOR' '(' EXPR (',' EXPR)* ')'))* (('DROP' 'ALL'? | 'REUSE') 'STORAGE')? alterMaterializedView_update_index_clauses? alterMaterializedView_parallel_clause? 'CASCADE'?
    ;
alterMaterializedView_split_table_partition
    : 'SPLIT' 'PARTITION' ID ('AT' '(' EXPR (',' EXPR)* ')' ('INTO' '(' alterMaterializedView_partition_spec ',' alterMaterializedView_partition_spec ')')? | 'VALUES' '(' EXPR (',' EXPR)* ')' ('INTO' '(' alterMaterializedView_partition_spec ',' alterMaterializedView_partition_spec ')')? | 'INTO' '(' alterMaterializedView_partition_spec (',' alterMaterializedView_partition_spec)* ',' alterMaterializedView_partition_spec ')') alterMaterializedView_split_nested_table_part? alterMaterializedView_filter_condition? alterMaterializedView_dependent_tables_clause? alterMaterializedView_update_index_clauses? alterMaterializedView_parallel_clause? alterMaterializedView_allow_disallow_clustering? 'ONLINE'?
    ;
alterMaterializedView_split_table_subpartition
    : 'SPLIT' 'SUBPARTITION' ID ('AT' '(' EXPR (',' EXPR)* ')' ('INTO' '(' alterMaterializedView_subpartition_spec ',' alterMaterializedView_subpartition_spec ')')? | 'VALUES' '(' EXPR (',' EXPR)* ')' ('INTO' '(' alterMaterializedView_subpartition_spec ',' alterMaterializedView_subpartition_spec ')')? | 'INTO' '(' alterMaterializedView_subpartition_spec (',' alterMaterializedView_subpartition_spec)* ',' alterMaterializedView_subpartition_spec ')') alterMaterializedView_filter_condition? alterMaterializedView_dependent_tables_clause? alterMaterializedView_update_index_clauses? alterMaterializedView_parallel_clause? alterMaterializedView_allow_disallow_clustering? 'ONLINE'?
    ;
alterMaterializedView_merge_table_partitions
    : 'MERGE' 'PARTITIONS' (alterMaterializedView_partition_or_key_value ',' alterMaterializedView_partition_or_key_value ('INTO' alterMaterializedView_partition_spec)? | alterMaterializedView_partition_or_key_value 'TO' alterMaterializedView_partition_or_key_value ('INTO' alterMaterializedView_partition_spec)?) alterMaterializedView_filter_condition? alterMaterializedView_dependent_tables_clause? alterMaterializedView_update_index_clauses? alterMaterializedView_parallel_clause? alterMaterializedView_allow_disallow_clustering? 'ONLINE'?
    ;
alterMaterializedView_merge_table_subpartitions
    : 'MERGE' 'SUBPARTITIONS' (alterMaterializedView_subpartition_or_key_value ',' alterMaterializedView_subpartition_or_key_value ('INTO' alterMaterializedView_subpartition_spec)? | alterMaterializedView_subpartition_or_key_value 'TO' alterMaterializedView_subpartition_or_key_value ('INTO' alterMaterializedView_subpartition_spec)?) alterMaterializedView_filter_condition? alterMaterializedView_dependent_tables_clause? alterMaterializedView_update_index_clauses? alterMaterializedView_parallel_clause? alterMaterializedView_allow_disallow_clustering? 'ONLINE'?
    ;
alterMaterializedView_exchange_partition_subpart
    : 'EXCHANGE' ('PARTITION' ID | 'SUBPARTITION' ID) 'WITH' 'TABLE' (ID '.')? ID (('INCLUDING' | 'EXCLUDING') 'INDEXES')? (('WITH' | 'WITHOUT') 'VALIDATION')? alterMaterializedView_exceptions_clause? alterMaterializedView_update_index_clauses? alterMaterializedView_parallel_clause? 'CASCADE'? 'ONLINE'?
    ;
alterMaterializedView_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T' | 'P')?
    ;
alterMaterializedView_annotation
    : ID STRING?
    ;
alterMaterializedView_index_org_table_clause
    : (alterMaterializedView_mapping_table_clauses | 'PCTTHRESHOLD' NUMBER | alterMaterializedView_prefix_compression | alterMaterializedView_index_org_overflow_clause)
    ;
alterMaterializedView_alter_overflow_clause
    : ('OVERFLOW' (alterMaterializedView_allocate_extent_clause | alterMaterializedView_deallocate_unused_clause | alterMaterializedView_segment_attributes_clause) | 'ADD' 'OVERFLOW' alterMaterializedView_segment_attributes_clause?)
    ;
alterMaterializedView_unusable_editions_clause
    : 'FOR' 'EDITION' ID (',' ID)*
    ;
alterMaterializedView_inmemory_memcompress
    : 'MEMCOMPRESS' 'FOR' ('QUERY' ('LOW' | 'HIGH')? | 'OLTP' | 'ARCHIVE' ('LOW' | 'HIGH')?)
    ;
alterMaterializedView_inmemory_priority
    : 'PRIORITY' ('LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL')
    ;
alterMaterializedView_inmemory_distribute
    : 'DISTRIBUTE' 'BY' ('ROWID' ('ACROSS' NUMBER?)? | 'PARTITION' | 'SUBPARTITION')
    ;
alterMaterializedView_inmemory_duplicate
    : 'DUPLICATE' ('ALL' | 'FOR' 'SERVICE' ID (',' ID)* | 'NONE')
    ;
alterMaterializedView_lob_retention_clause
    : 'RETENTION' ('MAX' | 'MIN' NUMBER | 'AUTO' | 'NONE')?
    ;
alterMaterializedView_lob_deduplicate_clause
    : ('DEDUPLICATE' | 'KEEP_DUPLICATES')
    ;
alterMaterializedView_lob_compression_clause
    : ('COMPRESS' ('HIGH' | 'MEDIUM' | 'LOW')? | 'NOCOMPRESS')
    ;
alterMaterializedView_deferred_segment_creation
    : 'SEGMENT' 'CREATION' ('IMMEDIATE' | 'DEFERRED')
    ;
alterMaterializedView_read_only_clause
    : ('READ' 'ONLY' | 'READ' 'WRITE')
    ;
alterMaterializedView_indexing_clause
    : 'INDEXING' ('ON' | 'OFF')
    ;
alterMaterializedView_segment_attributes_clause
    : alterMaterializedView_physical_attributes_clause? ('TABLESPACE' ID)? alterMaterializedView_logging_clause?
    ;
alterMaterializedView_inmemory_clause
    : ('INMEMORY' alterMaterializedView_inmemory_attributes? | 'NO' 'INMEMORY')
    ;
alterMaterializedView_prefix_compression
    : 'COMPRESS' NUMBER?
    ;
alterMaterializedView_range_subpartition_desc
    : 'SUBPARTITION' ID? alterMaterializedView_range_values_clause alterMaterializedView_read_only_clause? alterMaterializedView_indexing_clause? alterMaterializedView_partitioning_storage_clause?
    ;
alterMaterializedView_list_subpartition_desc
    : 'SUBPARTITION' ID? alterMaterializedView_list_values_clause alterMaterializedView_read_only_clause? alterMaterializedView_indexing_clause? alterMaterializedView_partitioning_storage_clause?
    ;
alterMaterializedView_individual_hash_subparts
    : 'SUBPARTITION' ID? alterMaterializedView_read_only_clause? alterMaterializedView_indexing_clause? alterMaterializedView_partitioning_storage_clause?
    ;
alterMaterializedView_hash_subparts_by_quantity
    : 'SUBPARTITIONS' NUMBER ('STORE' 'IN' '(' ID (',' ID)* ')')?
    ;
alterMaterializedView_partition_attributes
    : alterMaterializedView_physical_attributes_clause? alterMaterializedView_logging_clause? alterMaterializedView_allocate_extent_clause? alterMaterializedView_deallocate_unused_clause? ('SHRINK' 'SPACE' 'CASCADE'?)? ((alterMaterializedView_table_compression | alterMaterializedView_prefix_compression))? alterMaterializedView_inmemory_clause? (('LOB' '(' ID ')' | 'VARRAY' ID) '(' alterMaterializedView_modify_lob_parameters ')')* ('OVERFLOW' (alterMaterializedView_physical_attributes_clause | alterMaterializedView_logging_clause | alterMaterializedView_allocate_extent_clause | alterMaterializedView_deallocate_unused_clause | 'SHRINK' 'SPACE' 'CASCADE'?))?
    ;
alterMaterializedView_add_range_subpartition
    : 'ADD' alterMaterializedView_range_subpartition_desc (',' alterMaterializedView_range_subpartition_desc)* alterMaterializedView_dependent_tables_clause? alterMaterializedView_update_index_clauses?
    ;
alterMaterializedView_add_hash_subpartition
    : 'ADD' alterMaterializedView_individual_hash_subparts (',' alterMaterializedView_individual_hash_subparts)* alterMaterializedView_dependent_tables_clause? alterMaterializedView_update_index_clauses? alterMaterializedView_parallel_clause?
    ;
alterMaterializedView_add_list_subpartition
    : 'ADD' alterMaterializedView_list_subpartition_desc (',' alterMaterializedView_list_subpartition_desc)* alterMaterializedView_dependent_tables_clause? alterMaterializedView_update_index_clauses?
    ;
alterMaterializedView_update_index_clauses
    : (alterMaterializedView_update_global_index_clause | alterMaterializedView_update_all_indexes_clause)
    ;
alterMaterializedView_table_partition_description
    : alterMaterializedView_deferred_segment_creation? alterMaterializedView_read_only_clause? alterMaterializedView_indexing_clause? alterMaterializedView_segment_attributes_clause? (alterMaterializedView_table_compression | alterMaterializedView_prefix_compression)? alterMaterializedView_inmemory_clause? alterMaterializedView_ilm_clause? ('OVERFLOW' alterMaterializedView_segment_attributes_clause?)?
    ;
alterMaterializedView_filter_condition
    : 'INCLUDING' 'ROWS' 'WHERE' CONDITION
    ;
alterMaterializedView_allow_disallow_clustering
    : ('ALLOW' | 'DISALLOW') 'CLUSTERING'
    ;
alterMaterializedView_partitioning_storage_clause
    : ('TABLESPACE' ID)? ('OVERFLOW' 'TABLESPACE' ID)? alterMaterializedView_table_compression? alterMaterializedView_inmemory_clause?
    ;
alterMaterializedView_add_range_partition_clause
    : 'PARTITION' ID? alterMaterializedView_range_values_clause alterMaterializedView_table_partition_description? ('(' (alterMaterializedView_range_subpartition_desc | alterMaterializedView_list_subpartition_desc | alterMaterializedView_individual_hash_subparts) (',' (alterMaterializedView_range_subpartition_desc | alterMaterializedView_list_subpartition_desc | alterMaterializedView_individual_hash_subparts))* ')' | alterMaterializedView_hash_subparts_by_quantity)? alterMaterializedView_update_index_clauses?
    ;
alterMaterializedView_add_list_partition_clause
    : 'PARTITION' ID? alterMaterializedView_list_values_clause alterMaterializedView_table_partition_description? ('(' (alterMaterializedView_range_subpartition_desc | alterMaterializedView_list_subpartition_desc | alterMaterializedView_individual_hash_subparts) (',' (alterMaterializedView_range_subpartition_desc | alterMaterializedView_list_subpartition_desc | alterMaterializedView_individual_hash_subparts))* ')' | alterMaterializedView_hash_subparts_by_quantity)? alterMaterializedView_update_index_clauses?
    ;
alterMaterializedView_add_system_partition_clause
    : 'PARTITION' ID? alterMaterializedView_table_partition_description? alterMaterializedView_update_index_clauses?
    ;
alterMaterializedView_add_hash_partition_clause
    : 'PARTITION' ID? alterMaterializedView_partitioning_storage_clause? alterMaterializedView_update_index_clauses? alterMaterializedView_parallel_clause? alterMaterializedView_read_only_clause? alterMaterializedView_indexing_clause?
    ;
alterMaterializedView_partition_spec
    : 'PARTITION' ID? alterMaterializedView_table_partition_description?
    ;
alterMaterializedView_split_nested_table_part
    : 'NESTED' 'TABLE' ID 'INTO' '(' alterMaterializedView_nested_table_partition_spec ',' alterMaterializedView_nested_table_partition_spec ')'
    ;
alterMaterializedView_dependent_tables_clause
    : 'DEPENDENT' 'TABLES' '(' ID '(' alterMaterializedView_partition_spec (',' alterMaterializedView_partition_spec)* ')' (',' ID '(' alterMaterializedView_partition_spec (',' alterMaterializedView_partition_spec)* ')')* ')'
    ;
alterMaterializedView_subpartition_spec
    : 'SUBPARTITION' ID? alterMaterializedView_partitioning_storage_clause?
    ;
alterMaterializedView_partition_or_key_value
    : ID
    | 'FOR' '(' EXPR (',' EXPR)* ')'
    ;
alterMaterializedView_subpartition_or_key_value
    : ID
    | 'FOR' '(' EXPR (',' EXPR)* ')'
    ;
alterMaterializedView_exceptions_clause
    : 'EXCEPTIONS' 'INTO' (ID '.')? ID
    ;
alterMaterializedView_mapping_table_clauses
    : ('MAPPING' 'TABLE' | 'NOMAPPING')
    ;
alterMaterializedView_index_org_overflow_clause
    : ('INCLUDING' ID)? 'OVERFLOW' alterMaterializedView_segment_attributes_clause?
    ;
alterMaterializedView_range_values_clause
    : 'VALUES' 'LESS' 'THAN' '(' (EXPR | 'MAXVALUE') (',' (EXPR | 'MAXVALUE'))* ')'
    ;
alterMaterializedView_list_values_clause
    : 'VALUES' '(' (EXPR (',' EXPR)* | 'DEFAULT') ')'
    ;
alterMaterializedView_update_global_index_clause
    : ('UPDATE' | 'INVALIDATE') 'GLOBAL' 'INDEXES'
    ;
alterMaterializedView_update_all_indexes_clause
    : 'UPDATE' 'INDEXES' ('(' ID '(' (alterMaterializedView_index_partition_description | alterMaterializedView_index_subpartition_clause) (',' (alterMaterializedView_index_partition_description | alterMaterializedView_index_subpartition_clause))* ')' (',' ID '(' (alterMaterializedView_index_partition_description | alterMaterializedView_index_subpartition_clause) (',' (alterMaterializedView_index_partition_description | alterMaterializedView_index_subpartition_clause))* ')')* ')')?
    ;
alterMaterializedView_ilm_clause
    : 'ILM' ('ADD' 'POLICY' alterMaterializedView_ilm_policy_clause | 'DELETE' 'POLICY' ID | 'DELETE_ALL' | 'ENABLE_ALL' | 'DISABLE_ALL')
    ;
alterMaterializedView_nested_table_partition_spec
    : 'PARTITION' ID alterMaterializedView_segment_attributes_clause?
    ;
alterMaterializedView_index_partition_description
    : 'PARTITION' ID? alterMaterializedView_segment_attributes_clause? alterMaterializedView_index_compression? (('USABLE' | 'UNUSABLE'))?
    ;
alterMaterializedView_index_subpartition_clause
    : '(' 'SUBPARTITION' ID? ('TABLESPACE' ID)? alterMaterializedView_index_compression? (('USABLE' | 'UNUSABLE'))? (',' 'SUBPARTITION' ID? ('TABLESPACE' ID)? alterMaterializedView_index_compression? (('USABLE' | 'UNUSABLE'))?)* ')'
    | 'STORE' 'IN' '(' ID (',' ID)* ')'
    ;
alterMaterializedView_ilm_policy_clause
    : (alterMaterializedView_ilm_compression_policy | alterMaterializedView_ilm_tiering_policy | alterMaterializedView_ilm_inmemory_policy)
    ;
alterMaterializedView_index_compression
    : (alterMaterializedView_prefix_compression | alterMaterializedView_advanced_index_compression | 'NOCOMPRESS')
    ;
alterMaterializedView_ilm_compression_policy
    : (alterMaterializedView_table_compression | 'ROW' 'STORE' 'COMPRESS' 'ADVANCED'? | 'COLUMN' 'STORE' 'COMPRESS' ('FOR' 'QUERY' ('LOW' | 'HIGH')?)?) ('SEGMENT' | 'GROUP' | 'ROW') 'AFTER' alterMaterializedView_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION')
    ;
alterMaterializedView_ilm_tiering_policy
    : 'TIER' 'TO' ID ('SEGMENT' | 'GROUP' | 'ROW')? ('ON' ID)?
    ;
alterMaterializedView_ilm_inmemory_policy
    : ('SET' 'INMEMORY' alterMaterializedView_inmemory_memcompress? alterMaterializedView_inmemory_priority? | 'MODIFY' 'INMEMORY' alterMaterializedView_inmemory_memcompress | 'NO' 'INMEMORY') 'SEGMENT' 'AFTER' alterMaterializedView_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION')
    ;
alterMaterializedView_advanced_index_compression
    : 'COMPRESS' 'ADVANCED' ('LOW' | 'HIGH')?
    ;
alterMaterializedView_ilm_time_period
    : NUMBER ('DAY' | 'DAYS' | 'MONTH' | 'MONTHS' | 'YEAR' | 'YEARS')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-MATERIALIZED-ZONEMAP.html
alterMaterializedZonemap
    : 'ALTER' 'MATERIALIZED' 'ZONEMAP' ('IF' 'EXISTS')? (ID '.')? ID (alterMaterializedZonemap_alter_zonemap_attributes | alterMaterializedZonemap_zonemap_refresh_clause | ('ENABLE' | 'DISABLE') 'PRUNING' | 'COMPILE' | 'REBUILD' | 'UNUSABLE')
    ;
alterMaterializedZonemap_alter_zonemap_attributes
    : ('PCTFREE' NUMBER)? ('PCTUSED' NUMBER)? (('CACHE' | 'NOCACHE'))?
    ;
alterMaterializedZonemap_zonemap_refresh_clause
    : 'REFRESH' (('FAST' | 'COMPLETE' | 'FORCE'))? (('ON' 'COMMIT' | 'ON' 'DEMAND' | 'ON' 'LOAD' | 'ON' 'DATA' 'MOVEMENT' | 'ON' 'STATEMENT'))?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-OPERATOR.html
alterOperator
    : 'ALTER' 'OPERATOR' ('IF' 'EXISTS')? (ID '.')? ID (alterOperator_add_binding_clause | alterOperator_drop_binding_clause | 'COMPILE')
    ;
alterOperator_add_binding_clause
    : 'ADD' 'BINDING' '(' DATA_TYPE (',' DATA_TYPE)* ')' 'RETURN' '(' DATA_TYPE ')' alterOperator_implementation_clause? alterOperator_using_function_clause
    ;
alterOperator_drop_binding_clause
    : 'DROP' 'BINDING' '(' DATA_TYPE (',' DATA_TYPE)* ')' 'FORCE'?
    ;
alterOperator_implementation_clause
    : ('ANCILLARY' 'TO' ID '(' DATA_TYPE (',' DATA_TYPE)* ')' | alterOperator_context_clause)
    ;
alterOperator_using_function_clause
    : 'USING' (ID '.')? (ID '.')? ID
    ;
alterOperator_context_clause
    : 'WITH' 'INDEX' 'CONTEXT' ',' 'SCAN' 'CONTEXT' ID ('COMPUTE' 'ANCILLARY' 'DATA')? ('WITH' 'COLUMN' 'CONTEXT')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-OUTLINE.html
alterOutline
    : 'ALTER' 'OUTLINE' ('PUBLIC' | 'PRIVATE')? ID ('REBUILD' | 'RENAME' 'TO' ID | 'CHANGE' 'CATEGORY' 'TO' ID | ('ENABLE' | 'DISABLE'))
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-PACKAGE.html
alterPackage
    : 'ALTER' 'PACKAGE' ('IF' 'EXISTS')? (ID '.')? ID alterPackage_package_compile_clause? (('EDITIONABLE' | 'NONEDITIONABLE'))?
    ;
alterPackage_package_compile_clause
    : 'COMPILE' 'DEBUG'? ('PACKAGE' | 'SPECIFICATION' | 'BODY')? alterPackage_compiler_parameters_clause? ('REUSE' 'SETTINGS')?
    ;
alterPackage_compiler_parameters_clause
    : ID '=' EXPR
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-PLUGGABLE-DATABASE.html
alterPluggableDatabase
    : 'ALTER' 'PLUGGABLE' 'DATABASE' ID? (alterPluggableDatabase_pdb_unplug_clause | alterPluggableDatabase_pdb_settings_clauses | alterPluggableDatabase_pdb_datafile_clause | alterPluggableDatabase_pdb_recovery_clauses | alterPluggableDatabase_pdb_change_state | alterPluggableDatabase_pdb_change_state_from_root | alterPluggableDatabase_application_clauses | alterPluggableDatabase_snapshot_clauses)
    ;
alterPluggableDatabase_pdb_unplug_clause
    : 'UNPLUG' 'INTO' STRING alterPluggableDatabase_pdb_unplug_encrypt?
    ;
alterPluggableDatabase_pdb_settings_clauses
    : ('DEFAULT' 'EDITION' ID | 'SET' 'DEFAULT' 'TABLESPACE' ('BIGFILE' | 'SMALLFILE') | 'DEFAULT' 'TABLESPACE' ID | 'DEFAULT' 'TEMPORARY' 'TABLESPACE' ID | 'RENAME' 'GLOBAL_NAME' 'TO' ID | alterPluggableDatabase_set_time_zone_clause | alterPluggableDatabase_database_file_clauses | alterPluggableDatabase_supplemental_db_logging | alterPluggableDatabase_pdb_storage_clause | alterPluggableDatabase_pdb_logging_clauses | alterPluggableDatabase_pdb_refresh_mode_clause | alterPluggableDatabase_pdb_refresh_switchover_clause | 'CONTAINERS' ('DEFAULT' 'TARGET' (ID | 'NONE') | 'HOST' '=' STRING | 'PORT' '=' NUMBER) | 'PRIORITY' (NUMBER | 'NONE'))
    ;
alterPluggableDatabase_pdb_datafile_clause
    : ID? 'DATAFILE' ('ALL' | STRING | NUMBER) ('ONLINE' | 'OFFLINE')
    ;
alterPluggableDatabase_pdb_recovery_clauses
    : (alterPluggableDatabase_pdb_general_recovery | 'BACKUP' ('BEGIN' | 'END') | ('ENABLE' | 'DISABLE') 'RECOVERY')
    ;
alterPluggableDatabase_pdb_change_state
    : (alterPluggableDatabase_pdb_open | alterPluggableDatabase_pdb_close | alterPluggableDatabase_pdb_save_or_discard_state)
    ;
alterPluggableDatabase_pdb_change_state_from_root
    : (ID | 'ALL' | 'ALL' 'EXCEPT' '(' ID (',' ID)* ')') (alterPluggableDatabase_pdb_open | alterPluggableDatabase_pdb_close | alterPluggableDatabase_pdb_save_or_discard_state)
    ;
alterPluggableDatabase_application_clauses
    : 'APPLICATION' (ID ('BEGIN' ('INSTALL' 'VERSION' STRING ('COMMENT' STRING)? | 'PATCH' NUMBER ('MINIMUM' 'VERSION' STRING)? ('COMMENT' STRING)? | 'UPGRADE' 'FROM' STRING 'TO' STRING ('COMMENT' STRING)? | 'UNINSTALL') | 'END' ('INSTALL' ('VERSION' STRING)? | 'PATCH' NUMBER? | 'UPGRADE' ('TO' STRING)? | 'UNINSTALL') | 'SET' 'PATCH' NUMBER | 'SET' 'VERSION' STRING | 'SET' 'COMPATIBILITY' 'VERSION' (STRING | 'CURRENT') | 'SYNC' ('TO' (STRING | 'PATCH' NUMBER))?) | ID (',' ID)* 'SYNC' | 'ALL' ('EXCEPT' '(' ID (',' ID)* ')')? 'SYNC')
    ;
alterPluggableDatabase_snapshot_clauses
    : (alterPluggableDatabase_pdb_snapshot_clause | 'MATERIALIZE' | 'CREATE' 'SNAPSHOT' ID | 'DROP' 'SNAPSHOT' ID | 'SET' 'MAX_PDB_SNAPSHOTS' '=' NUMBER | alterPluggableDatabase_prepare_clause | 'DROP' 'MIRROR' 'COPY' ID | ('ENABLE' | 'DISABLE') 'LOST' 'WRITE' 'PROTECTION' | ('ENABLE' | 'DISABLE') 'BACKUP')
    ;
alterPluggableDatabase_pdb_unplug_encrypt
    : 'ENCRYPT' 'USING' ID
    ;
alterPluggableDatabase_set_time_zone_clause
    : 'SET' 'TIME_ZONE' '=' STRING
    ;
alterPluggableDatabase_database_file_clauses
    : ('RENAME' 'FILE' STRING 'TO' STRING (',' STRING 'TO' STRING)* | 'CREATE' 'DATAFILE' (STRING | NUMBER) 'AS' ('NEW' | alterPluggableDatabase_file_specification) | 'DATAFILE' (STRING | NUMBER) ('ONLINE' | 'OFFLINE' ('FOR' 'DROP')? | 'RESIZE' alterPluggableDatabase_size_clause | 'END' 'BACKUP' | 'ENCRYPT' | 'DECRYPT' | alterPluggableDatabase_autoextend_clause) | 'TEMPFILE' (STRING | NUMBER) ('RESIZE' alterPluggableDatabase_size_clause | alterPluggableDatabase_autoextend_clause | 'DROP' ('INCLUDING' 'DATAFILES')?) | 'MOVE' 'DATAFILE' (STRING | NUMBER) 'TO' STRING 'REUSE'? 'KEEP'?)
    ;
alterPluggableDatabase_supplemental_db_logging
    : ('ADD' 'SUPPLEMENTAL' 'LOG' ('DATA' | alterPluggableDatabase_supplemental_id_key_clause | alterPluggableDatabase_supplemental_plsql_clause | 'DATA' 'SUBSET' 'DATABASE' 'REPLICATION') | 'DROP' 'SUPPLEMENTAL' 'LOG' ('DATA' | alterPluggableDatabase_supplemental_id_key_clause | alterPluggableDatabase_supplemental_plsql_clause | 'DATA' 'SUBSET' 'DATABASE' 'REPLICATION'))
    ;
alterPluggableDatabase_pdb_storage_clause
    : 'STORAGE' '(' ('MAXSIZE' (alterPluggableDatabase_size_clause | 'UNLIMITED'))? ('MAX_AUDIT_SIZE' (alterPluggableDatabase_size_clause | 'UNLIMITED'))? ('MAX_DIAG_SIZE' (alterPluggableDatabase_size_clause | 'UNLIMITED'))? ')'
    ;
alterPluggableDatabase_pdb_logging_clauses
    : (alterPluggableDatabase_logging_clause | alterPluggableDatabase_pdb_force_logging_clause)
    ;
alterPluggableDatabase_pdb_refresh_mode_clause
    : 'REFRESH' 'MODE' ('NONE' | 'MANUAL' | 'EVERY' NUMBER ('MINUTES' | 'HOURS'))
    ;
alterPluggableDatabase_pdb_refresh_switchover_clause
    : 'REFRESH' 'SWITCHOVER' 'TO' 'PRIMARY' ID
    ;
alterPluggableDatabase_pdb_general_recovery
    : 'RECOVER' 'STANDBY'? ('UNTIL' ('CANCEL' | 'TIME' EXPR | 'CHANGE' NUMBER | 'SEQUENCE' NUMBER))? ('USING' 'BACKUP' 'CONTROLFILE')?
    ;
alterPluggableDatabase_pdb_open
    : 'OPEN' (('READ' 'WRITE' | 'READ' 'ONLY' | 'HYBRID' 'READ' 'ONLY' | ('READ' 'WRITE')? 'UPGRADE'))? 'RESTRICTED'? 'FORCE'? 'RESETLOGS'? alterPluggableDatabase_instances_clause? alterPluggableDatabase_services_clause?
    ;
alterPluggableDatabase_pdb_close
    : 'CLOSE' (('IMMEDIATE' | 'ABORT'))? alterPluggableDatabase_instances_clause? alterPluggableDatabase_relocate_clause?
    ;
alterPluggableDatabase_pdb_save_or_discard_state
    : ('SAVE' | 'DISCARD') 'STATE' alterPluggableDatabase_instances_clause?
    ;
alterPluggableDatabase_pdb_snapshot_clause
    : 'SNAPSHOT' ('NONE' | 'MANUAL' | 'EVERY' NUMBER ('MINUTES' | 'HOURS'))
    ;
alterPluggableDatabase_prepare_clause
    : 'PREPARE' 'MIRROR' 'COPY' ID ('WITH' ('EXTERNAL' | 'NORMAL' | 'HIGH') 'REDUNDANCY')? ('FOR' 'DATABASE' ID)?
    ;
alterPluggableDatabase_file_specification
    : STRING? ('SIZE' alterPluggableDatabase_size_clause)? 'REUSE'? alterPluggableDatabase_autoextend_clause?
    ;
alterPluggableDatabase_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T' | 'P')?
    ;
alterPluggableDatabase_autoextend_clause
    : 'AUTOEXTEND' ('OFF' | 'ON' ('NEXT' alterPluggableDatabase_size_clause)? alterPluggableDatabase_maxsize_clause?)
    ;
alterPluggableDatabase_supplemental_id_key_clause
    : 'DATA' '(' ('PRIMARY' 'KEY' | 'UNIQUE' | 'FOREIGN' 'KEY' | 'ALL') 'COLUMNS' ')'
    ;
alterPluggableDatabase_supplemental_plsql_clause
    : 'DATA' '(' 'PL' '/' 'SQL' 'CALL' ')'
    ;
alterPluggableDatabase_logging_clause
    : ('LOGGING' | 'NOLOGGING')
    ;
alterPluggableDatabase_pdb_force_logging_clause
    : ('ENABLE' 'FORCE' 'LOGGING' | 'DISABLE' 'FORCE' 'LOGGING' | 'ENABLE' 'FORCE' 'NOLOGGING' | 'DISABLE' 'FORCE' 'NOLOGGING' | 'SET' 'STANDBY' 'NOLOGGING' 'FOR' 'LOAD' 'PERFORMANCE' | 'SET' 'STANDBY' 'NOLOGGING' 'FOR' 'DATA' 'AVAILABILITY')
    ;
alterPluggableDatabase_instances_clause
    : 'INSTANCES' '=' '(' (ID (',' ID)* | 'ALL' | 'ALL' 'EXCEPT' '(' ID (',' ID)* ')') ')'
    ;
alterPluggableDatabase_services_clause
    : 'SERVICES' '=' '(' (ID (',' ID)* | 'ALL' | 'ALL' 'EXCEPT' '(' ID (',' ID)* ')') ')'
    ;
alterPluggableDatabase_relocate_clause
    : ('RELOCATE' ('TO' STRING)? | 'NORELOCATE')
    ;
alterPluggableDatabase_maxsize_clause
    : 'MAXSIZE' (alterPluggableDatabase_size_clause | 'UNLIMITED')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-PROCEDURE.html
alterProcedure
    : 'ALTER' 'PROCEDURE' ('IF' 'EXISTS')? (ID '.')? ID alterProcedure_procedure_compile_clause? (('EDITIONABLE' | 'NONEDITIONABLE'))?
    ;
alterProcedure_procedure_compile_clause
    : 'COMPILE' 'DEBUG'? alterProcedure_compiler_parameters_clause? ('REUSE' 'SETTINGS')?
    ;
alterProcedure_compiler_parameters_clause
    : ID '=' EXPR
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-PROFILE.html
alterProfile
    : 'ALTER' 'PROFILE' ID 'LIMIT' (alterProfile_resource_parameters | alterProfile_password_parameters) ('CONTAINER' '=' ('ALL' | 'CURRENT'))?
    ;
alterProfile_resource_parameters
    : ('SESSIONS_PER_USER' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'CPU_PER_SESSION' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'CPU_PER_CALL' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'CONNECT_TIME' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'IDLE_TIME' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'LOGICAL_READS_PER_SESSION' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'LOGICAL_READS_PER_CALL' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'PRIVATE_SGA' (alterProfile_size_clause | 'UNLIMITED' | 'DEFAULT') | 'COMPOSITE_LIMIT' (NUMBER | 'UNLIMITED' | 'DEFAULT')) (('SESSIONS_PER_USER' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'CPU_PER_SESSION' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'CPU_PER_CALL' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'CONNECT_TIME' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'IDLE_TIME' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'LOGICAL_READS_PER_SESSION' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'LOGICAL_READS_PER_CALL' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'PRIVATE_SGA' (alterProfile_size_clause | 'UNLIMITED' | 'DEFAULT') | 'COMPOSITE_LIMIT' (NUMBER | 'UNLIMITED' | 'DEFAULT')))*
    ;
alterProfile_password_parameters
    : ('PASSWORD_LIFE_TIME' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'PASSWORD_GRACE_TIME' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'PASSWORD_REUSE_TIME' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'PASSWORD_REUSE_MAX' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'PASSWORD_LOCK_TIME' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'FAILED_LOGIN_ATTEMPTS' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'INACTIVE_ACCOUNT_TIME' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'PASSWORD_ROLLOVER_TIME' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'PASSWORD_VERIFY_FUNCTION' (ID | 'NULL' | 'DEFAULT')) (('PASSWORD_LIFE_TIME' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'PASSWORD_GRACE_TIME' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'PASSWORD_REUSE_TIME' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'PASSWORD_REUSE_MAX' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'PASSWORD_LOCK_TIME' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'FAILED_LOGIN_ATTEMPTS' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'INACTIVE_ACCOUNT_TIME' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'PASSWORD_ROLLOVER_TIME' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'PASSWORD_VERIFY_FUNCTION' (ID | 'NULL' | 'DEFAULT')))*
    ;
alterProfile_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T' | 'P')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-RESOURCE-COST.html
alterResourceCost
    : 'ALTER' 'RESOURCE' 'COST' ('CPU_PER_SESSION' NUMBER | 'CONNECT_TIME' NUMBER | 'LOGICAL_READS_PER_SESSION' NUMBER | 'PRIVATE_SGA' NUMBER) (('CPU_PER_SESSION' NUMBER | 'CONNECT_TIME' NUMBER | 'LOGICAL_READS_PER_SESSION' NUMBER | 'PRIVATE_SGA' NUMBER))*
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-ROLE.html
alterRole
    : 'ALTER' 'ROLE' ID ('NOT' 'IDENTIFIED' | 'IDENTIFIED' 'BY' ID | 'IDENTIFIED' 'EXTERNALLY' | 'IDENTIFIED' 'GLOBALLY' 'AS' STRING | 'IDENTIFIED' 'USING' (ID '.')? ID) ('CONTAINER' '=' ('ALL' | 'CURRENT'))?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-ROLLBACK-SEGMENT.html
// fix: left out where optional, since no quoted page defines it or what it needs: storage_clause
alterRollbackSegment
    : 'ALTER' 'ROLLBACK' 'SEGMENT' ID ('ONLINE' | 'OFFLINE' | 'SHRINK' ('TO' alterRollbackSegment_size_clause)?)
    ;
alterRollbackSegment_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T' | 'P')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-SEQUENCE.html
alterSequence
    : 'ALTER' 'SEQUENCE' ('IF' 'EXISTS')? (ID '.')? ID ('INCREMENT' 'BY' NUMBER)? (('MAXVALUE' NUMBER | 'NOMAXVALUE'))? (('MINVALUE' NUMBER | 'NOMINVALUE'))? (('CYCLE' | 'NOCYCLE'))? (('CACHE' NUMBER | 'NOCACHE'))? (('ORDER' | 'NOORDER'))? (('KEEP' | 'NOKEEP'))? (('SCALE' ('EXTEND' | 'NOEXTEND')? | 'NOSCALE'))? (('SHARD' ('EXTEND' | 'NOEXTEND')? | 'NOSHARD'))? ('RESTART' ('WITH' NUMBER)?)?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-SESSION.html
alterSession
    : 'ALTER' 'SESSION' ('ADVISE' ('COMMIT' | 'ROLLBACK' | 'NOTHING') | 'CLOSE' 'DATABASE' 'LINK' ID | ('ENABLE' | 'DISABLE') 'COMMIT' 'IN' 'PROCEDURE' | ('ENABLE' | 'DISABLE') 'GUARD' | ('ENABLE' | 'DISABLE') 'PARALLEL' ('DML' | 'DDL' | 'QUERY') | 'FORCE' 'PARALLEL' ('DML' | 'DDL' | 'QUERY') ('PARALLEL' NUMBER)? | ('ENABLE' | 'DISABLE') 'RESUMABLE' ('TIMEOUT' NUMBER)? ('NAME' STRING)? | ('ENABLE' | 'DISABLE') 'SHARD' 'DDL' | 'SYNC' 'WITH' 'PRIMARY' | alterSession_alter_session_set_clause)
    ;
alterSession_alter_session_set_clause
    : 'SET' (ID '=' EXPR (',' ID '=' EXPR)* | 'EDITION' '=' ID | 'CONTAINER' '=' ID ('SERVICE' '=' ID)? | 'ROW' 'ARCHIVAL' 'VISIBILITY' '=' ('ACTIVE' | 'ALL') | 'DEFAULT_COLLATION' '=' (ID | 'NONE') | ('CONSTRAINT' | 'CONSTRAINTS') '=' ('IMMEDIATE' | 'DEFERRED' | 'DEFAULT') | 'CURRENT_SCHEMA' '=' ID | 'ERROR_ON_OVERLAP_TIME' '=' ('TRUE' | 'FALSE') | 'FLAGGER' '=' ('ENTRY' | 'OFF') | 'INSTANCE' '=' NUMBER | 'ISOLATION_LEVEL' '=' ('SERIALIZABLE' | 'READ' 'COMMITTED') | 'STANDBY_MAX_DATA_DELAY' '=' (NUMBER | 'NONE') | 'TIME_ZONE' '=' (STRING | 'LOCAL' | 'DBTIMEZONE') | 'USE_PRIVATE_OUTLINES' '=' ('TRUE' | 'FALSE' | ID) | 'USE_STORED_OUTLINES' '=' ('TRUE' | 'FALSE' | ID))
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-SYNONYM.html
alterSynonym
    : 'ALTER' 'PUBLIC'? 'SYNONYM' ('IF' 'EXISTS')? (ID '.')? ID ('EDITIONABLE' | 'NONEDITIONABLE' | 'COMPILE')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-SYSTEM.html
alterSystem
    : 'ALTER' 'SYSTEM' (alterSystem_archive_log_clause | alterSystem_checkpoint_clause | alterSystem_check_datafiles_clause | alterSystem_distributed_recov_clauses | 'SWITCH' 'LOGFILE' | ('SUSPEND' | 'RESUME') | alterSystem_quiesce_clauses | alterSystem_rolling_migration_clauses | alterSystem_rolling_patch_clauses | alterSystem_security_clauses | alterSystem_shutdown_dispatcher_clause | 'REGISTER' | alterSystem_alter_system_set_clause | alterSystem_alter_system_reset_clause | alterSystem_cancel_sql_clause | alterSystem_flush_clause | 'RELOCATE' 'CLIENT' STRING | alterSystem_end_session_clauses)
    ;
alterSystem_archive_log_clause
    : 'ARCHIVE' 'LOG' ('INSTANCE' STRING)? ('SEQUENCE' NUMBER | 'CHANGE' NUMBER | 'CURRENT' 'NOSWITCH'? | 'GROUP' NUMBER | 'LOGFILE' STRING ('USING' 'BACKUP' 'CONTROLFILE')? | 'NEXT' | 'ALL') ('THREAD' NUMBER)? ('TO' STRING)?
    ;
alterSystem_checkpoint_clause
    : 'CHECKPOINT' (('GLOBAL' | 'LOCAL'))?
    ;
alterSystem_check_datafiles_clause
    : 'CHECK' 'DATAFILES' (('GLOBAL' | 'LOCAL'))?
    ;
alterSystem_distributed_recov_clauses
    : ('ENABLE' | 'DISABLE') 'DISTRIBUTED' 'RECOVERY'
    ;
alterSystem_quiesce_clauses
    : ('QUIESCE' 'RESTRICTED' | 'UNQUIESCE')
    ;
alterSystem_rolling_migration_clauses
    : ('START' 'ROLLING' 'MIGRATION' 'TO' STRING | 'STOP' 'ROLLING' 'MIGRATION')
    ;
alterSystem_rolling_patch_clauses
    : ('START' 'ROLLING' 'PATCH' | 'STOP' 'ROLLING' 'PATCH')
    ;
alterSystem_security_clauses
    : ('ENABLE' 'RESTRICTED' 'SESSION' | 'DISABLE' 'RESTRICTED' 'SESSION' | 'SET' 'ENCRYPTION' 'WALLET' 'OPEN' 'IDENTIFIED' 'BY' ID | 'SET' 'ENCRYPTION' 'WALLET' 'CLOSE' ('IDENTIFIED' 'BY' ID)? | 'SET' 'ENCRYPTION' 'KEY' ('IDENTIFIED' 'BY' ID)?)
    ;
alterSystem_shutdown_dispatcher_clause
    : 'SHUTDOWN' 'IMMEDIATE'? STRING
    ;
alterSystem_alter_system_set_clause
    : 'SET' ID '=' EXPR (',' EXPR)* ('COMMENT' '=' STRING)? 'DEFERRED'? ('SCOPE' '=' ('MEMORY' | 'SPFILE' | 'BOTH'))? ('SID' '=' (STRING | '\'*\''))? ('CONTAINER' '=' ('ALL' | 'CURRENT'))?
    ;
alterSystem_alter_system_reset_clause
    : 'RESET' ID ('SCOPE' '=' ('MEMORY' | 'SPFILE' | 'BOTH'))? ('SID' '=' (STRING | '\'*\''))?
    ;
alterSystem_cancel_sql_clause
    : 'CANCEL' 'SQL' STRING (',' '@' NUMBER)? ('SQL_ID' STRING)?
    ;
alterSystem_flush_clause
    : 'FLUSH' ('SHARED_POOL' | 'GLOBAL' 'CONTEXT' | 'BUFFER_CACHE' (('GLOBAL' | 'LOCAL'))? | 'FLASH_CACHE' (('GLOBAL' | 'LOCAL'))? | 'REDO' 'TO' ID (('NO' 'CONFIRM' 'APPLY' | 'CONFIRM' 'APPLY'))? | 'PASSWORDFILE_METADATA_CACHE')
    ;
alterSystem_end_session_clauses
    : ('DISCONNECT' 'SESSION' STRING 'POST_TRANSACTION'? 'IMMEDIATE'? | 'KILL' 'SESSION' STRING (',' '@' NUMBER)? ('IMMEDIATE' | 'FORCE')? 'NOREPLAY'? ('TIMEOUT' NUMBER)?)
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-TABLE.html
// fix: left out where optional, since no quoted page defines it or what it needs: create_index_statement, external_part_subpart_data_props, json_parameters, lob_parameters, opaque_format_spec, out_of_line_part_storage, storage_clause
alterTable
    : 'ALTER' 'TABLE' (ID '.')? ID (alterTable_alter_table_properties | alterTable_column_clauses | alterTable_constraint_clauses | alterTable_alter_table_partitioning | alterTable_alter_table_partitionset | alterTable_alter_external_table | alterTable_move_table_clause | alterTable_modify_to_partitioned | alterTable_modify_opaque_type | alterTable_immutable_table_clauses | alterTable_blockchain_table_clauses | alterTable_duplicated_table_refresh | alterTable_enable_disable_clause) alterTable_enable_disable_clause* (('ENABLE' | 'DISABLE') 'TABLE' 'LOCK')? (('ENABLE' | 'DISABLE') 'ALL' 'TRIGGERS')?
    ;
alterTable_alter_table_properties
    : (alterTable_physical_attributes_clause | alterTable_logging_clause | alterTable_table_compression | alterTable_inmemory_table_clause | alterTable_ilm_clause | alterTable_supplemental_table_logging | alterTable_allocate_extent_clause | alterTable_deallocate_unused_clause | 'RENAME' 'LOB' '(' ID ')' 'TO' ID | 'RESULT_CACHE' '(' 'MODE' ('DEFAULT' | 'FORCE') ')' | alterTable_upgrade_table_clause | 'RECORDS_PER_BLOCK' NUMBER | alterTable_parallel_clause | alterTable_row_movement_clause | alterTable_logical_replication_clause | alterTable_flashback_archive_clause | 'SHRINK' 'SPACE' 'CASCADE'? | 'MODIFY' 'CLUSTERING' alterTable_attribute_clustering_clause | alterTable_alter_iot_clauses | alterTable_alter_xmlschema_clause | 'MEMOPTIMIZE' 'FOR' 'READ' | 'NO' 'MEMOPTIMIZE' 'FOR' 'READ' | 'MEMOPTIMIZE' 'FOR' 'WRITE' | 'NO' 'MEMOPTIMIZE' 'FOR' 'WRITE' | alterTable_annotations_clause | 'RENAME' 'TO' ID)
    ;
alterTable_column_clauses
    : (alterTable_add_column_clause | alterTable_modify_column_clauses | alterTable_drop_column_clause | alterTable_add_period_clause | alterTable_drop_period_clause | alterTable_rename_column_clause | alterTable_modify_collection_retrieval | alterTable_modify_lob_storage_clause | alterTable_alter_varray_col_properties)
    ;
alterTable_constraint_clauses
    : ('ADD' (alterTable_out_of_line_constraint | alterTable_out_of_line_ref_constraint | alterTable_constraint_state) | 'MODIFY' ('CONSTRAINT' ID | 'PRIMARY' 'KEY' | 'UNIQUE' '(' ID (',' ID)* ')') alterTable_constraint_state 'CASCADE'? | 'RENAME' 'CONSTRAINT' ID 'TO' ID | alterTable_drop_constraint_clause)
    ;
alterTable_alter_table_partitioning
    : (alterTable_modify_table_default_attrs | alterTable_alter_automatic_partitioning | alterTable_alter_interval_partitioning | alterTable_set_subpartition_template | alterTable_modify_table_partition | alterTable_modify_table_subpartition | alterTable_move_table_partition | alterTable_move_table_subpartition | alterTable_add_table_partition | alterTable_coalesce_table_partition | alterTable_coalesce_table_subpartition | alterTable_drop_table_partition | alterTable_drop_table_subpartition | alterTable_rename_partition_subpart | alterTable_truncate_partition_subpart | alterTable_split_table_partition | alterTable_split_table_subpartition | alterTable_merge_table_partitions | alterTable_merge_table_subpartitions | alterTable_exchange_partition_subpart)
    ;
alterTable_alter_table_partitionset
    : 'MODIFY' 'PARTITIONSET'
    ;
alterTable_alter_external_table
    : (alterTable_add_column_clause | alterTable_modify_column_clauses | alterTable_drop_column_clause | alterTable_parallel_clause | alterTable_external_table_data_props | 'REJECT' 'LIMIT' (NUMBER | 'UNLIMITED') | 'PROJECT' 'COLUMN' ('ALL' | 'REFERENCED'))
    ;
alterTable_move_table_clause
    : 'MOVE' 'ONLINE'? alterTable_segment_attributes_clause? alterTable_table_compression? alterTable_index_org_table_clause? alterTable_parallel_clause? alterTable_allow_disallow_clustering? alterTable_update_index_clauses?
    ;
alterTable_modify_to_partitioned
    : 'MODIFY' alterTable_table_partitioning_clauses alterTable_filter_condition? 'ONLINE'? alterTable_update_index_clauses?
    ;
alterTable_modify_opaque_type
    : 'MODIFY' 'OPAQUE' 'TYPE' ID 'STORE' 'AS' ('LOB' | 'OBJECT' 'RELATIONAL')? '(' alterTable_column_properties ')'
    ;
alterTable_immutable_table_clauses
    : ('NO' 'DROP' ('UNTIL' NUMBER ('DAY' | 'DAYS') 'IDLE')? | 'NO' 'DELETE' 'LOCKED'? ('UNTIL' NUMBER ('DAY' | 'DAYS') 'AFTER' 'INSERT')?)
    ;
alterTable_blockchain_table_clauses
    : ('NO' 'DROP' ('UNTIL' NUMBER ('DAY' | 'DAYS') 'IDLE')? | 'NO' 'DELETE' 'LOCKED'? ('UNTIL' NUMBER ('DAY' | 'DAYS') 'AFTER' 'INSERT')? | 'HASHING' 'USING' STRING | 'VERSION' STRING)
    ;
alterTable_duplicated_table_refresh
    : ('REFRESH' | 'NO' 'REFRESH')
    ;
alterTable_enable_disable_clause
    : ('ENABLE' | 'DISABLE') ('VALIDATE' | 'NOVALIDATE')? ('UNIQUE' '(' ID (',' ID)* ')' | 'PRIMARY' 'KEY' | 'CONSTRAINT' ID) alterTable_using_index_clause? alterTable_exceptions_clause? 'CASCADE'? (('KEEP' | 'DROP') 'INDEX')?
    ;
alterTable_physical_attributes_clause
    : ('PCTFREE' NUMBER | 'PCTUSED' NUMBER | 'INITRANS' NUMBER)
    ;
alterTable_logging_clause
    : ('LOGGING' | 'NOLOGGING')
    ;
alterTable_table_compression
    : ('COMPRESS' 'BASIC'? | 'ROW' 'STORE' 'COMPRESS' ('BASIC' | 'ADVANCED')? | 'COLUMN' 'STORE' 'COMPRESS' ('FOR' ('QUERY' | 'ARCHIVE') ('LOW' | 'HIGH')?)? | 'NOCOMPRESS')
    ;
alterTable_inmemory_table_clause
    : ('INMEMORY' alterTable_inmemory_memcompress? alterTable_inmemory_priority? alterTable_inmemory_distribute? alterTable_inmemory_duplicate? | 'NO' 'INMEMORY') alterTable_inmemory_column_clause?
    ;
alterTable_ilm_clause
    : 'ILM' ('ADD' 'POLICY' alterTable_ilm_policy_clause | 'DELETE' 'POLICY' ID | 'DELETE_ALL' | 'ENABLE_ALL' | 'DISABLE_ALL')
    ;
alterTable_supplemental_table_logging
    : ('ADD' 'SUPPLEMENTAL' 'LOG' ('DATA' '(' ('ALL' | 'PRIMARY' 'KEY' | 'UNIQUE' | 'FOREIGN' 'KEY') (',' ('ALL' | 'PRIMARY' 'KEY' | 'UNIQUE' | 'FOREIGN' 'KEY'))* ')' 'COLUMNS' | 'GROUP' ID '(' ID ('NO' 'LOG')? (',' ID ('NO' 'LOG')?)* ')' 'ALWAYS'?) | 'DROP' 'SUPPLEMENTAL' 'LOG' ('DATA' '(' ('ALL' | 'PRIMARY' 'KEY' | 'UNIQUE' | 'FOREIGN' 'KEY') (',' ('ALL' | 'PRIMARY' 'KEY' | 'UNIQUE' | 'FOREIGN' 'KEY'))* ')' 'COLUMNS' | 'GROUP' ID))
    ;
alterTable_allocate_extent_clause
    : 'ALLOCATE' 'EXTENT' ('(' ('SIZE' alterTable_size_clause | 'DATAFILE' STRING | 'INSTANCE' NUMBER)+ ')')?
    ;
alterTable_deallocate_unused_clause
    : 'DEALLOCATE' 'UNUSED' ('KEEP' alterTable_size_clause)?
    ;
alterTable_upgrade_table_clause
    : 'UPGRADE' 'NOT'? 'INCLUDING' 'DATA' alterTable_column_properties?
    ;
alterTable_parallel_clause
    : ('NOPARALLEL' | 'PARALLEL' NUMBER?)
    ;
alterTable_row_movement_clause
    : ('ENABLE' | 'DISABLE') 'ROW' 'MOVEMENT'
    ;
alterTable_logical_replication_clause
    : ('ENABLE' | 'DISABLE') 'LOGICAL' 'REPLICATION' ('ALL' 'KEYS' | 'ALLOW' 'NOVALIDATE' 'KEYS')?
    ;
alterTable_flashback_archive_clause
    : 'FLASHBACK' 'ARCHIVE' ID?
    | 'NO' 'FLASHBACK' 'ARCHIVE'
    ;
alterTable_attribute_clustering_clause
    : 'CLUSTERING' ((ID '.')? ID 'JOIN' (ID '.')? ID 'ON' '(' CONDITION ')')* 'BY' ('LINEAR' 'ORDER' | 'INTERLEAVED' 'ORDER')? '(' (ID | EXPR) ('ASC' | 'DESC')? (',' (ID | EXPR) ('ASC' | 'DESC')?)* ')' ID? alterTable_zonemap_clause?
    ;
alterTable_alter_iot_clauses
    : (alterTable_index_org_table_clause | alterTable_alter_overflow_clause | alterTable_alter_mapping_table_clauses | 'COALESCE')
    ;
alterTable_alter_xmlschema_clause
    : 'ALLOW' ('ANYSCHEMA' | 'NONSCHEMA')
    | 'DISALLOW' 'NONSCHEMA'
    ;
alterTable_annotations_clause
    : 'ANNOTATIONS' '(' ('ADD' | 'DROP' | 'REPLACE') alterTable_annotation (',' ('ADD' | 'DROP' | 'REPLACE') alterTable_annotation)* ')'
    ;
alterTable_add_column_clause
    : 'ADD' '(' (alterTable_column_definition | alterTable_virtual_column_definition) (',' (alterTable_column_definition | alterTable_virtual_column_definition))* ')' alterTable_column_properties?
    ;
alterTable_modify_column_clauses
    : 'MODIFY' ('(' alterTable_modify_col_properties (',' alterTable_modify_col_properties)* ')' | '(' alterTable_modify_virtcol_properties (',' alterTable_modify_virtcol_properties)* ')' | alterTable_modify_col_visibility | alterTable_modify_col_substitutable)
    ;
alterTable_drop_column_clause
    : ('SET' 'UNUSED' ('COLUMN' ID | '(' ID (',' ID)* ')') (('CASCADE' 'CONSTRAINTS' | 'INVALIDATE'))* 'ONLINE'? | 'DROP' ('COLUMN' ID | '(' ID (',' ID)* ')') (('CASCADE' 'CONSTRAINTS' | 'INVALIDATE'))* ('CHECKPOINT' NUMBER)? | 'DROP' ('UNUSED' 'COLUMNS' | 'COLUMNS' 'CONTINUE') ('CHECKPOINT' NUMBER)?)
    ;
alterTable_add_period_clause
    : 'ADD' '(' 'PERIOD' 'FOR' ID ('(' ID ',' ID ')')? ')'
    ;
alterTable_drop_period_clause
    : 'DROP' '(' 'PERIOD' 'FOR' ID ')'
    ;
alterTable_rename_column_clause
    : 'RENAME' 'COLUMN' ID 'TO' ID
    ;
alterTable_modify_collection_retrieval
    : 'MODIFY' 'NESTED' 'TABLE' ID 'RETURN' 'AS' ('LOCATOR' | 'VALUE')
    ;
alterTable_modify_lob_storage_clause
    : 'MODIFY' 'LOB' '(' ID ')' '(' alterTable_modify_lob_parameters ')'
    ;
alterTable_alter_varray_col_properties
    : 'MODIFY' 'VARRAY' ID '(' alterTable_modify_lob_parameters ')'
    ;
alterTable_out_of_line_constraint
    : ('CONSTRAINT' ID)? ('UNIQUE' '(' ID (',' ID)* ')' | 'PRIMARY' 'KEY' '(' ID (',' ID)* ')' | 'FOREIGN' 'KEY' '(' ID (',' ID)* ')' 'REFERENCES' (ID '.')? ID ('(' ID (',' ID)* ')')? ('ON' 'DELETE' ('CASCADE' | 'SET' 'NULL'))? | 'CHECK' '(' CONDITION ')') alterTable_constraint_state?
    ;
alterTable_out_of_line_ref_constraint
    : ('SCOPE' 'FOR' '(' ID ')' 'IS' (ID '.')? ID | 'REF' '(' ID ')' 'WITH' 'ROWID' | ('CONSTRAINT' ID)? 'FOREIGN' 'KEY' '(' ID ')' 'REFERENCES' (ID '.')? ID ('ON' 'DELETE' ('CASCADE' | 'SET' 'NULL'))? alterTable_constraint_state?)
    ;
alterTable_constraint_state
    : ('NOT'? 'DEFERRABLE' ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? | 'INITIALLY' ('DEFERRED' | 'IMMEDIATE') ('NOT'? 'DEFERRABLE')?)? ('RELY' | 'NORELY')? alterTable_using_index_clause? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? alterTable_exceptions_clause?
    ;
alterTable_drop_constraint_clause
    : 'DROP' (('PRIMARY' 'KEY' | 'UNIQUE' '(' ID (',' ID)* ')') 'CASCADE'? (('KEEP' | 'DROP') 'INDEX')? 'ONLINE'? | 'CONSTRAINT' ID 'CASCADE'? 'ONLINE'?)
    ;
alterTable_modify_table_default_attrs
    : 'MODIFY' 'DEFAULT' 'ATTRIBUTES' ('FOR' 'PARTITION' ID)? alterTable_deferred_segment_creation? alterTable_read_only_clause? alterTable_indexing_clause? alterTable_segment_attributes_clause? alterTable_table_compression? alterTable_inmemory_clause? alterTable_prefix_compression? ('OVERFLOW' alterTable_segment_attributes_clause)?
    ;
alterTable_alter_automatic_partitioning
    : 'SET' 'PARTITIONING' ('AUTOMATIC' | 'MANUAL')
    ;
alterTable_alter_interval_partitioning
    : 'SET' 'INTERVAL' '(' EXPR? ')'
    ;
alterTable_set_subpartition_template
    : 'SET' 'SUBPARTITION' 'TEMPLATE' ('(' (alterTable_range_subpartition_desc | alterTable_list_subpartition_desc | alterTable_individual_hash_subparts) (',' (alterTable_range_subpartition_desc | alterTable_list_subpartition_desc | alterTable_individual_hash_subparts))* ')' | alterTable_hash_subparts_by_quantity)
    ;
alterTable_modify_table_partition
    : 'MODIFY' 'PARTITION' ID (alterTable_partition_attributes | (alterTable_add_range_subpartition | alterTable_add_hash_subpartition | alterTable_add_list_subpartition) | 'COALESCE' 'SUBPARTITION' alterTable_update_index_clauses? alterTable_parallel_clause? | 'REBUILD'? 'UNUSABLE' 'LOCAL' 'INDEXES' | alterTable_read_only_clause | alterTable_indexing_clause)
    ;
alterTable_modify_table_subpartition
    : 'MODIFY' 'SUBPARTITION' ID (alterTable_allocate_extent_clause | alterTable_deallocate_unused_clause | 'SHRINK' 'SPACE' 'CASCADE'? | ('LOB' '(' ID ')' | 'VARRAY' ID) '(' alterTable_modify_lob_parameters ')' | 'REBUILD'? 'UNUSABLE' 'LOCAL' 'INDEXES' | ('ADD' | 'DROP') 'VALUES' '(' EXPR (',' EXPR)* ')' | alterTable_read_only_clause | alterTable_indexing_clause)
    ;
alterTable_move_table_partition
    : 'MOVE' 'PARTITION' ID ('MAPPING' 'TABLE')? alterTable_table_partition_description? alterTable_filter_condition? alterTable_update_index_clauses? alterTable_parallel_clause? alterTable_allow_disallow_clustering? 'ONLINE'?
    ;
alterTable_move_table_subpartition
    : 'MOVE' 'SUBPARTITION' ID alterTable_indexing_clause? alterTable_partitioning_storage_clause? alterTable_update_index_clauses? alterTable_filter_condition? alterTable_parallel_clause? alterTable_allow_disallow_clustering? 'ONLINE'?
    ;
alterTable_add_table_partition
    : 'ADD' (alterTable_add_range_partition_clause | alterTable_add_list_partition_clause | alterTable_add_system_partition_clause | alterTable_add_hash_partition_clause)
    ;
alterTable_coalesce_table_partition
    : 'COALESCE' 'PARTITION' alterTable_update_index_clauses? alterTable_parallel_clause? alterTable_allow_disallow_clustering?
    ;
alterTable_coalesce_table_subpartition
    : 'COALESCE' 'SUBPARTITION' ID? alterTable_update_index_clauses? alterTable_parallel_clause? alterTable_allow_disallow_clustering?
    ;
alterTable_drop_table_partition
    : 'DROP' 'PARTITION' (ID | 'FOR' '(' EXPR (',' EXPR)* ')') (',' (ID | 'FOR' '(' EXPR (',' EXPR)* ')'))* alterTable_update_index_clauses? alterTable_parallel_clause?
    ;
alterTable_drop_table_subpartition
    : 'DROP' 'SUBPARTITION' (ID | 'FOR' '(' EXPR (',' EXPR)* ')') (',' (ID | 'FOR' '(' EXPR (',' EXPR)* ')'))* alterTable_update_index_clauses? alterTable_parallel_clause?
    ;
alterTable_rename_partition_subpart
    : 'RENAME' ('PARTITION' ID | 'SUBPARTITION' ID) 'TO' ID
    ;
alterTable_truncate_partition_subpart
    : 'TRUNCATE' ('PARTITION' | 'SUBPARTITION') (ID | 'FOR' '(' EXPR (',' EXPR)* ')') (',' (ID | 'FOR' '(' EXPR (',' EXPR)* ')'))* (('DROP' 'ALL'? | 'REUSE') 'STORAGE')? alterTable_update_index_clauses? alterTable_parallel_clause? 'CASCADE'?
    ;
alterTable_split_table_partition
    : 'SPLIT' 'PARTITION' ID ('AT' '(' EXPR (',' EXPR)* ')' ('INTO' '(' alterTable_partition_spec ',' alterTable_partition_spec ')')? | 'VALUES' '(' EXPR (',' EXPR)* ')' ('INTO' '(' alterTable_partition_spec ',' alterTable_partition_spec ')')? | 'INTO' '(' alterTable_partition_spec (',' alterTable_partition_spec)* ',' alterTable_partition_spec ')') alterTable_split_nested_table_part? alterTable_filter_condition? alterTable_dependent_tables_clause? alterTable_update_index_clauses? alterTable_parallel_clause? alterTable_allow_disallow_clustering? 'ONLINE'?
    ;
alterTable_split_table_subpartition
    : 'SPLIT' 'SUBPARTITION' ID ('AT' '(' EXPR (',' EXPR)* ')' ('INTO' '(' alterTable_subpartition_spec ',' alterTable_subpartition_spec ')')? | 'VALUES' '(' EXPR (',' EXPR)* ')' ('INTO' '(' alterTable_subpartition_spec ',' alterTable_subpartition_spec ')')? | 'INTO' '(' alterTable_subpartition_spec (',' alterTable_subpartition_spec)* ',' alterTable_subpartition_spec ')') alterTable_filter_condition? alterTable_dependent_tables_clause? alterTable_update_index_clauses? alterTable_parallel_clause? alterTable_allow_disallow_clustering? 'ONLINE'?
    ;
alterTable_merge_table_partitions
    : 'MERGE' 'PARTITIONS' (alterTable_partition_or_key_value ',' alterTable_partition_or_key_value ('INTO' alterTable_partition_spec)? | alterTable_partition_or_key_value 'TO' alterTable_partition_or_key_value ('INTO' alterTable_partition_spec)?) alterTable_filter_condition? alterTable_dependent_tables_clause? alterTable_update_index_clauses? alterTable_parallel_clause? alterTable_allow_disallow_clustering? 'ONLINE'?
    ;
alterTable_merge_table_subpartitions
    : 'MERGE' 'SUBPARTITIONS' (alterTable_subpartition_or_key_value ',' alterTable_subpartition_or_key_value ('INTO' alterTable_subpartition_spec)? | alterTable_subpartition_or_key_value 'TO' alterTable_subpartition_or_key_value ('INTO' alterTable_subpartition_spec)?) alterTable_filter_condition? alterTable_dependent_tables_clause? alterTable_update_index_clauses? alterTable_parallel_clause? alterTable_allow_disallow_clustering? 'ONLINE'?
    ;
alterTable_exchange_partition_subpart
    : 'EXCHANGE' ('PARTITION' ID | 'SUBPARTITION' ID) 'WITH' 'TABLE' (ID '.')? ID (('INCLUDING' | 'EXCLUDING') 'INDEXES')? (('WITH' | 'WITHOUT') 'VALIDATION')? alterTable_exceptions_clause? alterTable_update_index_clauses? alterTable_parallel_clause? 'CASCADE'? 'ONLINE'?
    ;
alterTable_external_table_data_props
    : ('DEFAULT' 'DIRECTORY' ID)? ('ACCESS' 'PARAMETERS' 'USING' 'CLOB' QUERY)? ('LOCATION' '(' (ID ':')? STRING (',' (ID ':')? STRING)* ')')?
    ;
alterTable_segment_attributes_clause
    : (alterTable_physical_attributes_clause | 'TABLESPACE' ID | alterTable_logging_clause)+
    ;
alterTable_index_org_table_clause
    : (alterTable_mapping_table_clauses | 'PCTTHRESHOLD' NUMBER | alterTable_prefix_compression | alterTable_index_org_overflow_clause)
    ;
alterTable_allow_disallow_clustering
    : ('ALLOW' | 'DISALLOW') 'CLUSTERING'
    ;
alterTable_update_index_clauses
    : (alterTable_update_global_index_clause | alterTable_update_all_indexes_clause)
    ;
alterTable_table_partitioning_clauses
    : (alterTable_range_partitions | alterTable_list_partitions | alterTable_hash_partitions | alterTable_composite_range_partitions | alterTable_composite_list_partitions | alterTable_composite_hash_partitions | alterTable_reference_partitioning | alterTable_system_partitioning | alterTable_consistent_hash_partitions | alterTable_consistent_hash_with_subpartitions | alterTable_partitionset_clauses)
    ;
alterTable_filter_condition
    : 'INCLUDING' 'ROWS' 'WHERE' CONDITION
    ;
alterTable_column_properties
    : (alterTable_object_type_col_properties | alterTable_nested_table_col_properties | (alterTable_varray_col_properties | alterTable_lob_storage_clause) ('(' alterTable_lob_partition_storage (',' alterTable_lob_partition_storage)* ')')? | alterTable_xmltype_column_properties | alterTable_json_storage_clause)+
    ;
alterTable_using_index_clause
    : 'USING' 'INDEX' ((ID '.')? ID | alterTable_index_attributes)
    ;
alterTable_exceptions_clause
    : 'EXCEPTIONS' 'INTO' (ID '.')? ID
    ;
alterTable_inmemory_memcompress
    : 'MEMCOMPRESS' 'FOR' ('DML' | 'QUERY' ('LOW' | 'HIGH')? | 'CAPACITY' ('LOW' | 'HIGH')?)
    | 'NO' 'MEMCOMPRESS'
    ;
alterTable_inmemory_priority
    : 'PRIORITY' ('NONE' | 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL')
    ;
alterTable_inmemory_distribute
    : 'DISTRIBUTE' ('AUTO' | 'BY' ('ROWID' 'RANGE' | 'PARTITION' | 'SUBPARTITION'))? ('FOR' 'SERVICE' ('DEFAULT' | 'ALL' | ID | 'NONE'))?
    ;
alterTable_inmemory_duplicate
    : 'DUPLICATE' 'ALL'?
    | 'NO' 'DUPLICATE'
    ;
alterTable_inmemory_column_clause
    : (('INMEMORY' alterTable_inmemory_memcompress? | 'NO' 'INMEMORY') '(' ID (',' ID)* ')')+
    ;
alterTable_ilm_policy_clause
    : (alterTable_ilm_compression_policy | alterTable_ilm_tiering_policy | alterTable_ilm_inmemory_policy)
    ;
alterTable_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T' | 'P')?
    ;
alterTable_zonemap_clause
    : ('WITH' 'MATERIALIZED' 'ZONEMAP' ('(' ID ')')? | 'WITHOUT' 'MATERIALIZED' 'ZONEMAP')
    ;
alterTable_alter_overflow_clause
    : ('OVERFLOW' (alterTable_allocate_extent_clause | alterTable_deallocate_unused_clause | alterTable_segment_attributes_clause) | 'ADD' 'OVERFLOW' alterTable_segment_attributes_clause?)
    ;
alterTable_alter_mapping_table_clauses
    : 'MAPPING' 'TABLE' (alterTable_allocate_extent_clause | alterTable_deallocate_unused_clause)
    ;
alterTable_annotation
    : ID STRING?
    ;
alterTable_column_definition
    : ID DATA_TYPE ('COLLATE' ID)? 'SORT'? ('VISIBLE' | 'INVISIBLE')? ('DEFAULT' ('ON' 'NULL')? EXPR | alterTable_identity_clause)? ('ENCRYPT' alterTable_encryption_spec)? alterTable_inline_constraint* alterTable_inline_ref_constraint? alterTable_annotations_clause?
    ;
alterTable_virtual_column_definition
    : ID DATA_TYPE? ('COLLATE' ID)? ('VISIBLE' | 'INVISIBLE')? ('GENERATED' 'ALWAYS')? 'AS' '(' EXPR ')' 'VIRTUAL'? alterTable_evaluation_edition_clause? alterTable_unusable_editions_clause? alterTable_inline_constraint*
    ;
alterTable_modify_col_properties
    : ID DATA_TYPE? ('COLLATE' ID)? ('DEFAULT' ('ON' 'NULL')? EXPR | alterTable_identity_clause | 'DROP' 'IDENTITY')? (('ENCRYPT' alterTable_encryption_spec | 'DECRYPT'))? alterTable_inline_constraint* alterTable_lob_storage_clause? alterTable_alter_xmlschema_clause? alterTable_annotations_clause?
    ;
alterTable_modify_virtcol_properties
    : ID DATA_TYPE? ('COLLATE' ID)? ('GENERATED' 'ALWAYS')? 'AS' '(' EXPR ')' 'VIRTUAL'? alterTable_evaluation_edition_clause? alterTable_unusable_editions_clause? alterTable_inline_constraint*
    ;
alterTable_modify_col_visibility
    : ID ('VISIBLE' | 'INVISIBLE')
    ;
alterTable_modify_col_substitutable
    : 'COLUMN' ID 'NOT'? 'SUBSTITUTABLE' 'AT' 'ALL' 'LEVELS' 'FORCE'?
    ;
alterTable_modify_lob_parameters
    : ('PCTVERSION' NUMBER | 'FREEPOOLS' NUMBER | 'REBUILD' 'FREEPOOLS' | alterTable_lob_retention_clause | alterTable_lob_deduplicate_clause | alterTable_lob_compression_clause | ('ENCRYPT' alterTable_encryption_spec | 'DECRYPT') | ('CACHE' | 'NOCACHE' | 'CACHE' 'READS') alterTable_logging_clause? | alterTable_allocate_extent_clause | 'SHRINK' 'SPACE' 'CASCADE'? | alterTable_deallocate_unused_clause)+
    ;
alterTable_deferred_segment_creation
    : 'SEGMENT' 'CREATION' ('IMMEDIATE' | 'DEFERRED')
    ;
alterTable_read_only_clause
    : ('READ' 'ONLY' | 'READ' 'WRITE')
    ;
alterTable_indexing_clause
    : 'INDEXING' ('ON' | 'OFF')
    ;
alterTable_inmemory_clause
    : ('INMEMORY' alterTable_inmemory_attributes? | 'NO' 'INMEMORY')
    ;
alterTable_prefix_compression
    : ('COMPRESS' NUMBER? | 'NOCOMPRESS')
    ;
alterTable_range_subpartition_desc
    : 'SUBPARTITION' ID? alterTable_range_values_clause alterTable_read_only_clause? alterTable_indexing_clause? alterTable_partitioning_storage_clause?
    ;
alterTable_list_subpartition_desc
    : 'SUBPARTITION' ID? alterTable_list_values_clause alterTable_read_only_clause? alterTable_indexing_clause? alterTable_partitioning_storage_clause?
    ;
alterTable_individual_hash_subparts
    : 'SUBPARTITION' ID? alterTable_read_only_clause? alterTable_indexing_clause? alterTable_partitioning_storage_clause?
    ;
alterTable_hash_subparts_by_quantity
    : 'SUBPARTITIONS' NUMBER ('STORE' 'IN' '(' ID (',' ID)* ')')?
    ;
alterTable_partition_attributes
    : alterTable_physical_attributes_clause? alterTable_logging_clause? alterTable_allocate_extent_clause? alterTable_deallocate_unused_clause? ('SHRINK' 'SPACE' 'CASCADE'?)? ((alterTable_table_compression | alterTable_prefix_compression))? alterTable_inmemory_clause? (('LOB' '(' ID ')' | 'VARRAY' ID) '(' alterTable_modify_lob_parameters ')')* ('OVERFLOW' (alterTable_physical_attributes_clause | alterTable_logging_clause | alterTable_allocate_extent_clause | alterTable_deallocate_unused_clause | 'SHRINK' 'SPACE' 'CASCADE'?))?
    ;
alterTable_add_range_subpartition
    : 'ADD' alterTable_range_subpartition_desc (',' alterTable_range_subpartition_desc)* alterTable_dependent_tables_clause? alterTable_update_index_clauses?
    ;
alterTable_add_hash_subpartition
    : 'ADD' alterTable_individual_hash_subparts (',' alterTable_individual_hash_subparts)* alterTable_dependent_tables_clause? alterTable_update_index_clauses? alterTable_parallel_clause?
    ;
alterTable_add_list_subpartition
    : 'ADD' alterTable_list_subpartition_desc (',' alterTable_list_subpartition_desc)* alterTable_dependent_tables_clause? alterTable_update_index_clauses?
    ;
alterTable_table_partition_description
    : alterTable_deferred_segment_creation? alterTable_read_only_clause? alterTable_indexing_clause? alterTable_segment_attributes_clause? (alterTable_table_compression | alterTable_prefix_compression)? alterTable_inmemory_clause? alterTable_ilm_clause? ('OVERFLOW' alterTable_segment_attributes_clause?)?
    ;
alterTable_partitioning_storage_clause
    : ('TABLESPACE' ID)? ('OVERFLOW' 'TABLESPACE' ID)? alterTable_table_compression? alterTable_inmemory_clause?
    ;
alterTable_add_range_partition_clause
    : 'PARTITION' ID? alterTable_range_values_clause alterTable_table_partition_description? ('(' (alterTable_range_subpartition_desc | alterTable_list_subpartition_desc | alterTable_individual_hash_subparts) (',' (alterTable_range_subpartition_desc | alterTable_list_subpartition_desc | alterTable_individual_hash_subparts))* ')' | alterTable_hash_subparts_by_quantity)? alterTable_update_index_clauses?
    ;
alterTable_add_list_partition_clause
    : 'PARTITION' ID? alterTable_list_values_clause alterTable_table_partition_description? ('(' (alterTable_range_subpartition_desc | alterTable_list_subpartition_desc | alterTable_individual_hash_subparts) (',' (alterTable_range_subpartition_desc | alterTable_list_subpartition_desc | alterTable_individual_hash_subparts))* ')' | alterTable_hash_subparts_by_quantity)? alterTable_update_index_clauses?
    ;
alterTable_add_system_partition_clause
    : 'PARTITION' ID? alterTable_table_partition_description? alterTable_update_index_clauses?
    ;
alterTable_add_hash_partition_clause
    : 'PARTITION' ID? alterTable_partitioning_storage_clause? alterTable_update_index_clauses? alterTable_parallel_clause? alterTable_read_only_clause? alterTable_indexing_clause?
    ;
alterTable_partition_spec
    : 'PARTITION' ID? alterTable_table_partition_description?
    ;
alterTable_split_nested_table_part
    : 'NESTED' 'TABLE' ID 'INTO' '(' alterTable_nested_table_partition_spec ',' alterTable_nested_table_partition_spec ')'
    ;
alterTable_dependent_tables_clause
    : 'DEPENDENT' 'TABLES' '(' ID '(' alterTable_partition_spec (',' alterTable_partition_spec)* ')' (',' ID '(' alterTable_partition_spec (',' alterTable_partition_spec)* ')')* ')'
    ;
alterTable_subpartition_spec
    : 'SUBPARTITION' ID? alterTable_partitioning_storage_clause?
    ;
alterTable_partition_or_key_value
    : ID
    | 'FOR' '(' EXPR (',' EXPR)* ')'
    ;
alterTable_subpartition_or_key_value
    : ID
    | 'FOR' '(' EXPR (',' EXPR)* ')'
    ;
alterTable_mapping_table_clauses
    : ('MAPPING' 'TABLE' | 'NOMAPPING')
    ;
alterTable_index_org_overflow_clause
    : ('INCLUDING' ID)? 'OVERFLOW' alterTable_segment_attributes_clause?
    ;
alterTable_update_global_index_clause
    : ('UPDATE' | 'INVALIDATE') 'GLOBAL' 'INDEXES'
    ;
alterTable_update_all_indexes_clause
    : 'UPDATE' 'INDEXES' ('(' ID '(' (alterTable_index_partition_description | alterTable_index_subpartition_clause) (',' (alterTable_index_partition_description | alterTable_index_subpartition_clause))* ')' (',' ID '(' (alterTable_index_partition_description | alterTable_index_subpartition_clause) (',' (alterTable_index_partition_description | alterTable_index_subpartition_clause))* ')')* ')')?
    ;
alterTable_range_partitions
    : 'PARTITION' 'BY' 'RANGE' '(' ID (',' ID)* ')' ('INTERVAL' '(' EXPR ')' ('STORE' 'IN' '(' ID (',' ID)* ')')?)? '(' 'PARTITION' ID? alterTable_range_values_clause alterTable_table_partition_description? (',' 'PARTITION' ID? alterTable_range_values_clause alterTable_table_partition_description?)* ')'
    ;
alterTable_list_partitions
    : 'PARTITION' 'BY' 'LIST' '(' ID (',' ID)* ')' ('AUTOMATIC' ('STORE' 'IN' '(' ID (',' ID)* ')')?)? '(' 'PARTITION' ID? alterTable_list_values_clause alterTable_table_partition_description? (',' 'PARTITION' ID? alterTable_list_values_clause alterTable_table_partition_description?)* ')'
    ;
alterTable_hash_partitions
    : 'PARTITION' 'BY' 'HASH' '(' ID (',' ID)* ')' (alterTable_individual_hash_partitions | alterTable_hash_partitions_by_quantity)
    ;
alterTable_composite_range_partitions
    : 'PARTITION' 'BY' 'RANGE' '(' ID (',' ID)* ')' ('INTERVAL' '(' EXPR ')' ('STORE' 'IN' '(' ID (',' ID)* ')')?)? (alterTable_subpartition_by_range | alterTable_subpartition_by_list | alterTable_subpartition_by_hash) '(' 'PARTITION' ID? alterTable_range_values_clause (alterTable_table_partition_description ((alterTable_range_subpartition_desc | alterTable_list_subpartition_desc | alterTable_individual_hash_subparts | alterTable_hash_subparts_by_quantity))? (',' alterTable_table_partition_description ((alterTable_range_subpartition_desc | alterTable_list_subpartition_desc | alterTable_individual_hash_subparts | alterTable_hash_subparts_by_quantity))?)*)* ')'
    ;
alterTable_composite_list_partitions
    : 'PARTITION' 'BY' 'LIST' '(' ID (',' ID)* ')' ('AUTOMATIC' ('STORE' 'IN' '(' ID (',' ID)* ')')?)? (alterTable_subpartition_by_range | alterTable_subpartition_by_list | alterTable_subpartition_by_hash) '(' 'PARTITION' ID? alterTable_list_values_clause (alterTable_table_partition_description ((alterTable_range_subpartition_desc | alterTable_list_subpartition_desc | alterTable_individual_hash_subparts | alterTable_hash_subparts_by_quantity))? (',' alterTable_table_partition_description ((alterTable_range_subpartition_desc | alterTable_list_subpartition_desc | alterTable_individual_hash_subparts | alterTable_hash_subparts_by_quantity))?)*)* ')'
    ;
alterTable_composite_hash_partitions
    : 'PARTITION' 'BY' 'HASH' '(' ID (',' ID)* ')' (alterTable_subpartition_by_range | alterTable_subpartition_by_list | alterTable_subpartition_by_hash) (alterTable_individual_hash_partitions | alterTable_hash_partitions_by_quantity)
    ;
alterTable_reference_partitioning
    : 'PARTITION' 'BY' 'REFERENCE' '(' ID ')' ('(' 'PARTITION' ID? alterTable_table_partition_description? (',' 'PARTITION' ID? alterTable_table_partition_description?)* ')')?
    ;
alterTable_system_partitioning
    : ('PARTITION' 'BY' 'SYSTEM')? ('PARTITIONS' NUMBER | '(' 'PARTITION' ID? alterTable_table_partition_description? (',' 'PARTITION' ID? alterTable_table_partition_description?)* ')')?
    ;
alterTable_consistent_hash_partitions
    : 'PARTITION' 'BY' 'CONSISTENT' 'HASH' '(' ID (',' ID)* ')' ('PARTITIONS' 'AUTO')? 'TABLESPACE' 'SET' ID
    ;
alterTable_consistent_hash_with_subpartitions
    : 'PARTITION' 'BY' 'CONSISTENT' 'HASH' '(' ID (',' ID)* ')' (alterTable_subpartition_by_range | alterTable_subpartition_by_list | alterTable_subpartition_by_hash) ('PARTITIONS' 'AUTO')? 'TABLESPACE' 'SET' ID
    ;
alterTable_partitionset_clauses
    : (alterTable_range_partitionset_clause | alterTable_list_partitionset_clause)
    ;
alterTable_object_type_col_properties
    : 'COLUMN' ID alterTable_substitutable_column_clause
    ;
alterTable_nested_table_col_properties
    : 'NESTED' 'TABLE' (ID | 'COLUMN_VALUE') alterTable_substitutable_column_clause? ('LOCAL' | 'GLOBAL')? 'STORE' 'AS' ID ('(' '(' alterTable_object_properties ')' alterTable_physical_properties? alterTable_column_properties? ')')? ('RETURN' 'AS'? ('LOCATOR' | 'VALUE'))?
    ;
alterTable_varray_col_properties
    : 'VARRAY' ID (alterTable_substitutable_column_clause? 'STORE' 'AS' ('SECUREFILE' | 'BASICFILE')? 'LOB' (ID? '(' alterTable_lob_storage_parameters ')' | ID) | alterTable_substitutable_column_clause)
    ;
alterTable_lob_storage_clause
    : 'LOB' '(' ID (',' ID)* ')' 'STORE' 'AS' ID? '(' alterTable_lob_storage_parameters ')'
    ;
alterTable_lob_partition_storage
    : 'PARTITION' ID (alterTable_lob_storage_clause | alterTable_varray_col_properties)+ ('(' 'SUBPARTITION' ID (alterTable_lob_storage_clause | alterTable_varray_col_properties)+ (')')+)?
    ;
alterTable_xmltype_column_properties
    : 'XMLTYPE' 'COLUMN'? ID alterTable_xmltype_storage? alterTable_xmlschema_spec?
    ;
alterTable_json_storage_clause
    : 'JSON' '(' ID ')' 'STORE' 'AS'
    ;
alterTable_index_attributes
    : (alterTable_physical_attributes_clause | 'TABLESPACE' (ID | 'DEFAULT') | alterTable_logging_clause | alterTable_index_compression | alterTable_partial_index_clause | alterTable_parallel_clause)+
    ;
alterTable_ilm_compression_policy
    : (alterTable_table_compression | 'ROW' 'STORE' 'COMPRESS' 'ADVANCED'? | 'COLUMN' 'STORE' 'COMPRESS' ('FOR' 'QUERY' ('LOW' | 'HIGH')?)?) ('SEGMENT' | 'GROUP' | 'ROW') 'AFTER' alterTable_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION')
    ;
alterTable_ilm_tiering_policy
    : 'TIER' 'TO' ID ('SEGMENT' | 'GROUP' | 'ROW')? ('ON' ID)?
    ;
alterTable_ilm_inmemory_policy
    : ('SET' 'INMEMORY' alterTable_inmemory_memcompress? alterTable_inmemory_priority? | 'MODIFY' 'INMEMORY' alterTable_inmemory_memcompress | 'NO' 'INMEMORY') 'SEGMENT' 'AFTER' alterTable_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION')
    ;
alterTable_identity_clause
    : 'GENERATED' ('ALWAYS' | 'BY' 'DEFAULT' ('ON' 'NULL')?) 'AS' 'IDENTITY' ('(' alterTable_identity_options ')')?
    ;
alterTable_encryption_spec
    : ('USING' STRING)? ('IDENTIFIED' 'BY' ID)? STRING? ('NO'? 'SALT')?
    ;
alterTable_inline_constraint
    : ('CONSTRAINT' ID)? ('NOT'? 'NULL' | 'UNIQUE' | 'PRIMARY' 'KEY' | 'REFERENCES' (ID '.')? ID ('(' ID ')')? ('ON' 'DELETE' ('CASCADE' | 'SET' 'NULL'))? | 'CHECK' '(' CONDITION ')') alterTable_constraint_state?
    ;
alterTable_inline_ref_constraint
    : ('SCOPE' 'IS' (ID '.')? ID | 'WITH' 'ROWID' | ('CONSTRAINT' ID)? 'REFERENCES' (ID '.')? ID ('(' ID (',' ID)* ')')? ('ON' 'DELETE' ('CASCADE' | 'SET' 'NULL'))? alterTable_constraint_state?)
    ;
alterTable_evaluation_edition_clause
    : 'EVALUATE' 'USING' ('CURRENT' 'EDITION' | 'EDITION' ID | 'NULL' 'EDITION')
    ;
alterTable_unusable_editions_clause
    : ('UNUSABLE' 'BEFORE' ('CURRENT' 'EDITION' | 'EDITION' ID))? ('UNUSABLE' 'BEGINNING' 'WITH' ('CURRENT' 'EDITION' | 'EDITION' ID | 'NULL' 'EDITION'))?
    ;
alterTable_lob_retention_clause
    : 'RETENTION' ('MAX' | 'MIN' NUMBER | 'AUTO' | 'NONE')?
    ;
alterTable_lob_deduplicate_clause
    : ('DEDUPLICATE' | 'KEEP_DUPLICATES')
    ;
alterTable_lob_compression_clause
    : ('COMPRESS' ('HIGH' | 'MEDIUM' | 'LOW')? | 'NOCOMPRESS')
    ;
alterTable_inmemory_attributes
    : alterTable_inmemory_memcompress? alterTable_inmemory_priority? alterTable_inmemory_distribute? alterTable_inmemory_duplicate?
    ;
alterTable_range_values_clause
    : 'VALUES' 'LESS' 'THAN' '(' (EXPR | 'MAXVALUE') (',' (EXPR | 'MAXVALUE'))* ')'
    ;
alterTable_list_values_clause
    : 'VALUES' '(' (EXPR (',' EXPR)* | 'DEFAULT') ')'
    ;
alterTable_nested_table_partition_spec
    : 'PARTITION' ID alterTable_segment_attributes_clause?
    ;
alterTable_index_partition_description
    : 'PARTITION' ID? alterTable_segment_attributes_clause? alterTable_index_compression? (('USABLE' | 'UNUSABLE'))?
    ;
alterTable_index_subpartition_clause
    : '(' 'SUBPARTITION' ID? ('TABLESPACE' ID)? alterTable_index_compression? (('USABLE' | 'UNUSABLE'))? (',' 'SUBPARTITION' ID? ('TABLESPACE' ID)? alterTable_index_compression? (('USABLE' | 'UNUSABLE'))?)* ')'
    | 'STORE' 'IN' '(' ID (',' ID)* ')'
    ;
alterTable_individual_hash_partitions
    : '(' 'PARTITION' ID? alterTable_partitioning_storage_clause? (',' 'PARTITION' ID? alterTable_partitioning_storage_clause?)* ')'
    ;
alterTable_hash_partitions_by_quantity
    : 'PARTITIONS' NUMBER ('STORE' 'IN' '(' ID (',' ID)* ')')?
    ;
alterTable_subpartition_by_range
    : 'SUBPARTITION' 'BY' 'RANGE' '(' ID (',' ID)* ')' alterTable_subpartition_template?
    ;
alterTable_subpartition_by_list
    : 'SUBPARTITION' 'BY' 'LIST' '(' ID (',' ID)* ')' alterTable_subpartition_template?
    ;
alterTable_subpartition_by_hash
    : 'SUBPARTITION' 'BY' 'HASH' '(' ID (',' ID)* ')' ('SUBPARTITIONS' NUMBER ('STORE' 'IN' '(' ID (',' ID)* ')')? | alterTable_subpartition_template)?
    ;
alterTable_range_partitionset_clause
    : 'PARTITIONSET' 'BY' 'RANGE' '(' ID (',' ID)* ')' 'PARTITION' 'BY' 'CONSISTENT' 'HASH' '(' ID (',' ID)* ')' (alterTable_subpartition_by_range | alterTable_subpartition_by_list | alterTable_subpartition_by_hash)? 'PARTITIONS' 'AUTO' '(' 'PARTITIONSET' ID alterTable_range_values_clause ('TABLESPACE' 'SET' ID)? alterTable_lob_storage_clause? alterTable_subpartition_template? (',' 'PARTITIONSET' ID alterTable_range_values_clause ('TABLESPACE' 'SET' ID)? alterTable_lob_storage_clause? alterTable_subpartition_template?)* ')'
    ;
alterTable_list_partitionset_clause
    : 'PARTITIONSET' 'BY' 'LIST' '(' ID ')' 'PARTITION' 'BY' 'CONSISTENT' 'HASH' '(' ID (',' ID)* ')' (alterTable_subpartition_by_range | alterTable_subpartition_by_list | alterTable_subpartition_by_hash)? 'PARTITIONS' 'AUTO' '(' 'PARTITIONSET' ID alterTable_list_values_clause ('TABLESPACE' 'SET' ID)? alterTable_lob_storage_clause? alterTable_subpartition_template? (',' 'PARTITIONSET' ID alterTable_list_values_clause ('TABLESPACE' 'SET' ID)? alterTable_lob_storage_clause? alterTable_subpartition_template?)* ')'
    ;
alterTable_substitutable_column_clause
    : ('ELEMENT' 'IS' 'OF' 'TYPE'? '(' 'ONLY'? ID ')' | 'NOT'? 'SUBSTITUTABLE' 'AT' 'ALL' 'LEVELS')
    ;
alterTable_object_properties
    : (ID ('DEFAULT' EXPR)? ((alterTable_inline_constraint | alterTable_inline_ref_constraint)+)? | alterTable_out_of_line_constraint | alterTable_out_of_line_ref_constraint | alterTable_supplemental_logging_props) (',' (ID ('DEFAULT' EXPR)? ((alterTable_inline_constraint | alterTable_inline_ref_constraint)+)? | alterTable_out_of_line_constraint | alterTable_out_of_line_ref_constraint | alterTable_supplemental_logging_props))*
    ;
alterTable_physical_properties
    : alterTable_deferred_segment_creation? alterTable_segment_attributes_clause? alterTable_table_compression? alterTable_inmemory_table_clause? ((alterTable_heap_org_table_clause | alterTable_index_org_table_clause))? ('CLUSTER' ID '(' ID (',' ID)* ')')?
    ;
alterTable_lob_storage_parameters
    : (('TABLESPACE' ID | alterTable_logging_clause)+ | ('ENABLE' | 'DISABLE') 'STORAGE' 'IN' 'ROW' | 'CHUNK' NUMBER | 'PCTVERSION' NUMBER | 'FREEPOOLS' NUMBER | alterTable_lob_retention_clause | alterTable_lob_deduplicate_clause | alterTable_lob_compression_clause | ('ENCRYPT' alterTable_encryption_spec | 'DECRYPT') | ('CACHE' | 'NOCACHE' | 'CACHE' 'READS') alterTable_logging_clause?)+
    ;
alterTable_xmltype_storage
    : 'STORE' 'AS' ('OBJECT' 'RELATIONAL' | ('SECUREFILE' | 'BASICFILE')? ('CLOB' | 'BINARY' 'XML') (ID ('(' alterTable_lob_storage_parameters ')')? | '(' alterTable_lob_storage_parameters ')')?)
    ;
alterTable_xmlschema_spec
    : ('XMLSCHEMA' STRING)? 'ELEMENT' ID ('STORE' 'ALL' 'VARRAYS' 'AS' ('LOBS' | 'TABLES'))? (('ALLOW' | 'DISALLOW') 'NONSCHEMA')? (('ALLOW' | 'DISALLOW') 'ANYSCHEMA')?
    ;
alterTable_index_compression
    : ('COMPRESS' NUMBER? | 'NOCOMPRESS')
    ;
alterTable_partial_index_clause
    : 'INDEXING' ('FULL' | 'PARTIAL')
    ;
alterTable_ilm_time_period
    : NUMBER ('DAY' | 'DAYS' | 'MONTH' | 'MONTHS' | 'YEAR' | 'YEARS')
    ;
alterTable_identity_options
    : ('START' 'WITH' (NUMBER | 'LIMIT' 'VALUE') | 'INCREMENT' 'BY' NUMBER | 'MAXVALUE' NUMBER | 'NOMAXVALUE' | 'MINVALUE' NUMBER | 'NOMINVALUE' | 'CYCLE' | 'NOCYCLE' | 'CACHE' NUMBER | 'NOCACHE' | 'ORDER' | 'NOORDER')+
    ;
alterTable_subpartition_template
    : 'SUBPARTITION' 'TEMPLATE' '(' (alterTable_range_subpartition_desc | alterTable_list_subpartition_desc | alterTable_individual_hash_subparts) (',' (alterTable_range_subpartition_desc | alterTable_list_subpartition_desc | alterTable_individual_hash_subparts))* ')'
    | 'SUBPARTITION' 'TEMPLATE' NUMBER
    ;
alterTable_supplemental_logging_props
    : 'SUPPLEMENTAL' 'LOG' (alterTable_supplemental_log_grp_clause | alterTable_supplemental_id_key_clause)
    ;
alterTable_heap_org_table_clause
    : (('CACHE' | 'NOCACHE'))?
    ;
alterTable_supplemental_log_grp_clause
    : 'GROUP' ID '(' ID ('NO' 'LOG')? (',' ID ('NO' 'LOG')?)* ')' 'ALWAYS'?
    ;
alterTable_supplemental_id_key_clause
    : 'DATA' '(' ('PRIMARY' 'KEY' | 'UNIQUE' | 'FOREIGN' 'KEY' | 'ALL') 'COLUMNS' ')'
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-TABLESPACE-SET.html
// fix: left out where optional, since no quoted page defines it or what it needs: storage_clause
alterTablespaceSet
    : 'ALTER' 'TABLESPACE' 'SET' ID alterTablespaceSet_alter_tablespace_attrs
    ;
alterTablespaceSet_alter_tablespace_attrs
    : (alterTablespaceSet_default_tablespace_params | 'RESIZE' alterTablespaceSet_size_clause | 'COALESCE' | 'RENAME' 'TO' ID | 'BEGIN' 'BACKUP' | 'END' 'BACKUP' | alterTablespaceSet_datafile_tempfile_clauses | alterTablespaceSet_tablespace_logging_clauses | alterTablespaceSet_tablespace_state_clauses | alterTablespaceSet_autoextend_clause | alterTablespaceSet_alter_tablespace_encryption | alterTablespaceSet_lost_write_protection)
    ;
alterTablespaceSet_default_tablespace_params
    : 'DEFAULT' alterTablespaceSet_default_table_compression? alterTablespaceSet_default_index_compression? alterTablespaceSet_inmemory_clause? alterTablespaceSet_ilm_clause?
    ;
alterTablespaceSet_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T' | 'P')?
    ;
alterTablespaceSet_datafile_tempfile_clauses
    : ('ADD' ('DATAFILE' | 'TEMPFILE') (alterTablespaceSet_file_specification (',' alterTablespaceSet_file_specification)?)* | 'DROP' ('DATAFILE' | 'TEMPFILE') (STRING | NUMBER) | 'SHRINK' 'TEMPFILE' (STRING | NUMBER) ('KEEP' alterTablespaceSet_size_clause)? | 'RENAME' 'DATAFILE' STRING (',' STRING)* 'TO' STRING (',' STRING)* | ('DATAFILE' | 'TEMPFILE') ('ONLINE' | 'OFFLINE'))
    ;
alterTablespaceSet_tablespace_logging_clauses
    : (alterTablespaceSet_logging_clause | 'NO'? 'FORCE' 'LOGGING')
    ;
alterTablespaceSet_tablespace_state_clauses
    : ('ONLINE' | 'OFFLINE' ('NORMAL' | 'TEMPORARY' | 'IMMEDIATE')? | 'READ' 'ONLY' | 'READ' 'WRITE' | 'PERMANENT' | 'TEMPORARY')
    ;
alterTablespaceSet_autoextend_clause
    : 'AUTOEXTEND' ('OFF' | 'ON' ('NEXT' alterTablespaceSet_size_clause)? alterTablespaceSet_maxsize_clause?)
    ;
alterTablespaceSet_alter_tablespace_encryption
    : 'ENCRYPTION' ('ONLINE' alterTablespaceSet_tablespace_encryption_spec? alterTablespaceSet_ts_file_name_convert? | 'OFFLINE' ('ENCRYPT' alterTablespaceSet_tablespace_encryption_spec? | 'DECRYPT') | 'FINISH' ('ENCRYPT' | 'DECRYPT')? alterTablespaceSet_ts_file_name_convert?)
    ;
alterTablespaceSet_lost_write_protection
    : ('ENABLE' 'LOST' 'WRITE' 'PROTECTION' | 'DISABLE' 'LOST' 'WRITE' 'PROTECTION')
    ;
alterTablespaceSet_default_table_compression
    : 'TABLE' ('COMPRESS' ('BASIC' | 'FOR' 'OLTP')? | 'ROW' 'STORE' 'COMPRESS' ('BASIC' | 'ADVANCED')? | 'COLUMN' 'STORE' 'COMPRESS' ('FOR' ('QUERY' | 'ARCHIVE') ('LOW' | 'HIGH')?)? | 'NOCOMPRESS')
    ;
alterTablespaceSet_default_index_compression
    : 'INDEX' ('COMPRESS' ('ADVANCED' ('LOW' | 'HIGH')?)? | 'NOCOMPRESS')
    ;
alterTablespaceSet_inmemory_clause
    : ('INMEMORY' alterTablespaceSet_inmemory_attributes? | 'NO' 'INMEMORY')
    ;
alterTablespaceSet_ilm_clause
    : 'ILM' ('ADD' 'POLICY' alterTablespaceSet_ilm_policy_clause | 'DELETE' 'POLICY' ID | 'DELETE_ALL' | 'ENABLE_ALL' | 'DISABLE_ALL')
    ;
alterTablespaceSet_file_specification
    : STRING? ('SIZE' alterTablespaceSet_size_clause)? 'REUSE'? alterTablespaceSet_autoextend_clause?
    ;
alterTablespaceSet_logging_clause
    : ('LOGGING' | 'NOLOGGING')
    ;
alterTablespaceSet_maxsize_clause
    : 'MAXSIZE' (alterTablespaceSet_size_clause | 'UNLIMITED')
    ;
alterTablespaceSet_tablespace_encryption_spec
    : 'USING' STRING
    ;
alterTablespaceSet_ts_file_name_convert
    : 'FILE_NAME_CONVERT' '=' '(' STRING ',' STRING (',' STRING ',' STRING)* ')' 'KEEP'?
    ;
alterTablespaceSet_inmemory_attributes
    : alterTablespaceSet_inmemory_memcompress? alterTablespaceSet_inmemory_priority? alterTablespaceSet_inmemory_distribute? alterTablespaceSet_inmemory_duplicate?
    ;
alterTablespaceSet_ilm_policy_clause
    : (alterTablespaceSet_ilm_compression_policy | alterTablespaceSet_ilm_tiering_policy | alterTablespaceSet_ilm_inmemory_policy)
    ;
alterTablespaceSet_inmemory_memcompress
    : 'MEMCOMPRESS' 'FOR' ('QUERY' ('LOW' | 'HIGH')? | 'OLTP' | 'ARCHIVE' ('LOW' | 'HIGH')?)
    ;
alterTablespaceSet_inmemory_priority
    : 'PRIORITY' ('LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL')
    ;
alterTablespaceSet_inmemory_distribute
    : 'DISTRIBUTE' 'BY' ('ROWID' ('ACROSS' NUMBER?)? | 'PARTITION' | 'SUBPARTITION')
    ;
alterTablespaceSet_inmemory_duplicate
    : 'DUPLICATE' ('ALL' | 'FOR' 'SERVICE' ID (',' ID)* | 'NONE')
    ;
alterTablespaceSet_ilm_compression_policy
    : (alterTablespaceSet_table_compression | 'ROW' 'STORE' 'COMPRESS' 'ADVANCED'? | 'COLUMN' 'STORE' 'COMPRESS' ('FOR' 'QUERY' ('LOW' | 'HIGH')?)?) ('SEGMENT' | 'GROUP' | 'ROW') 'AFTER' alterTablespaceSet_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION')
    ;
alterTablespaceSet_ilm_tiering_policy
    : 'TIER' 'TO' ID ('SEGMENT' | 'GROUP' | 'ROW')? ('ON' ID)?
    ;
alterTablespaceSet_ilm_inmemory_policy
    : ('SET' 'INMEMORY' alterTablespaceSet_inmemory_memcompress? alterTablespaceSet_inmemory_priority? | 'MODIFY' 'INMEMORY' alterTablespaceSet_inmemory_memcompress | 'NO' 'INMEMORY') 'SEGMENT' 'AFTER' alterTablespaceSet_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION')
    ;
alterTablespaceSet_table_compression
    : ('COMPRESS' 'BASIC'? | 'COMPRESS' 'FOR' ('OLTP' | 'QUERY' ('LOW' | 'HIGH')? | 'ARCHIVE' ('LOW' | 'HIGH')?) | 'NOCOMPRESS')
    ;
alterTablespaceSet_ilm_time_period
    : NUMBER ('DAY' | 'DAYS' | 'MONTH' | 'MONTHS' | 'YEAR' | 'YEARS')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-TABLESPACE.html
// fix: left out where optional, since no quoted page defines it or what it needs: storage_clause
alterTablespace
    : 'ALTER' 'TABLESPACE' ('IF' 'EXISTS')? ID alterTablespace_alter_tablespace_attrs
    ;
alterTablespace_alter_tablespace_attrs
    : (alterTablespace_default_tablespace_params | 'MINIMUM' 'EXTENT' alterTablespace_size_clause | 'RESIZE' alterTablespace_size_clause | 'COALESCE' | 'SHRINK' 'SPACE' ('KEEP' alterTablespace_size_clause)? | 'RENAME' 'TO' ID | 'BEGIN' 'BACKUP' | 'END' 'BACKUP' | alterTablespace_datafile_tempfile_clauses | alterTablespace_tablespace_logging_clauses | alterTablespace_tablespace_group_clause | alterTablespace_tablespace_state_clauses | alterTablespace_autoextend_clause | alterTablespace_flashback_mode_clause | alterTablespace_tablespace_retention_clause | alterTablespace_alter_tablespace_encryption | alterTablespace_lost_write_protection)
    ;
alterTablespace_default_tablespace_params
    : 'DEFAULT' alterTablespace_default_table_compression? alterTablespace_default_index_compression? alterTablespace_inmemory_clause? alterTablespace_ilm_clause?
    ;
alterTablespace_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T' | 'P')?
    ;
alterTablespace_datafile_tempfile_clauses
    : ('ADD' ('DATAFILE' | 'TEMPFILE') (alterTablespace_file_specification (',' alterTablespace_file_specification)?)* | 'DROP' ('DATAFILE' | 'TEMPFILE') (STRING | NUMBER) | 'SHRINK' 'TEMPFILE' (STRING | NUMBER) ('KEEP' alterTablespace_size_clause)? | 'RENAME' 'DATAFILE' STRING (',' STRING)* 'TO' STRING (',' STRING)* | ('DATAFILE' | 'TEMPFILE') ('ONLINE' | 'OFFLINE'))
    ;
alterTablespace_tablespace_logging_clauses
    : (alterTablespace_logging_clause | 'NO'? 'FORCE' 'LOGGING')
    ;
alterTablespace_tablespace_group_clause
    : 'GROUP' (ID | STRING)
    ;
alterTablespace_tablespace_state_clauses
    : ('ONLINE' | 'OFFLINE' ('NORMAL' | 'TEMPORARY' | 'IMMEDIATE')? | 'READ' 'ONLY' | 'READ' 'WRITE' | 'PERMANENT' | 'TEMPORARY')
    ;
alterTablespace_autoextend_clause
    : 'AUTOEXTEND' ('OFF' | 'ON' ('NEXT' alterTablespace_size_clause)? alterTablespace_maxsize_clause?)
    ;
alterTablespace_flashback_mode_clause
    : 'FLASHBACK' ('ON' | 'OFF')
    ;
alterTablespace_tablespace_retention_clause
    : 'RETENTION' ('GUARANTEE' | 'NOGUARANTEE')
    ;
alterTablespace_alter_tablespace_encryption
    : 'ENCRYPTION' ('ONLINE' alterTablespace_tablespace_encryption_spec? alterTablespace_ts_file_name_convert? | 'OFFLINE' ('ENCRYPT' alterTablespace_tablespace_encryption_spec? | 'DECRYPT') | 'FINISH' ('ENCRYPT' | 'DECRYPT')? alterTablespace_ts_file_name_convert?)
    ;
alterTablespace_lost_write_protection
    : ('ENABLE' | 'DISABLE' | 'REMOVE' | 'SUSPEND') 'LOST' 'WRITE' 'PROTECTION'
    ;
alterTablespace_default_table_compression
    : 'TABLE' ('COMPRESS' ('BASIC' | 'FOR' 'OLTP')? | 'ROW' 'STORE' 'COMPRESS' ('BASIC' | 'ADVANCED')? | 'COLUMN' 'STORE' 'COMPRESS' ('FOR' ('QUERY' | 'ARCHIVE') ('LOW' | 'HIGH')?)? | 'NOCOMPRESS')
    ;
alterTablespace_default_index_compression
    : 'INDEX' ('COMPRESS' ('ADVANCED' ('LOW' | 'HIGH')?)? | 'NOCOMPRESS')
    ;
alterTablespace_inmemory_clause
    : ('INMEMORY' alterTablespace_inmemory_attributes? | 'NO' 'INMEMORY')
    ;
alterTablespace_ilm_clause
    : 'ILM' ('ADD' 'POLICY' alterTablespace_ilm_policy_clause | 'DELETE' 'POLICY' ID | 'DELETE_ALL' | 'ENABLE_ALL' | 'DISABLE_ALL')
    ;
alterTablespace_file_specification
    : STRING? ('SIZE' alterTablespace_size_clause)? 'REUSE'? alterTablespace_autoextend_clause?
    ;
alterTablespace_logging_clause
    : ('LOGGING' | 'NOLOGGING')
    ;
alterTablespace_maxsize_clause
    : 'MAXSIZE' ('UNLIMITED' | alterTablespace_size_clause)
    ;
alterTablespace_tablespace_encryption_spec
    : 'USING' STRING
    ;
alterTablespace_ts_file_name_convert
    : 'FILE_NAME_CONVERT' '=' '(' STRING ',' STRING (',' STRING ',' STRING)* ')' 'KEEP'?
    ;
alterTablespace_inmemory_attributes
    : alterTablespace_inmemory_memcompress? alterTablespace_inmemory_priority? alterTablespace_inmemory_distribute? alterTablespace_inmemory_duplicate?
    ;
alterTablespace_ilm_policy_clause
    : (alterTablespace_ilm_compression_policy | alterTablespace_ilm_tiering_policy | alterTablespace_ilm_inmemory_policy)
    ;
alterTablespace_inmemory_memcompress
    : 'MEMCOMPRESS' 'FOR' ('QUERY' ('LOW' | 'HIGH')? | 'OLTP' | 'ARCHIVE' ('LOW' | 'HIGH')?)
    ;
alterTablespace_inmemory_priority
    : 'PRIORITY' ('LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL')
    ;
alterTablespace_inmemory_distribute
    : 'DISTRIBUTE' 'BY' ('ROWID' ('ACROSS' NUMBER?)? | 'PARTITION' | 'SUBPARTITION')
    ;
alterTablespace_inmemory_duplicate
    : 'DUPLICATE' ('ALL' | 'FOR' 'SERVICE' ID (',' ID)* | 'NONE')
    ;
alterTablespace_ilm_compression_policy
    : (alterTablespace_table_compression | 'ROW' 'STORE' 'COMPRESS' 'ADVANCED'? | 'COLUMN' 'STORE' 'COMPRESS' ('FOR' 'QUERY' ('LOW' | 'HIGH')?)?) ('SEGMENT' | 'GROUP' | 'ROW') 'AFTER' alterTablespace_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION')
    ;
alterTablespace_ilm_tiering_policy
    : 'TIER' 'TO' ID ('SEGMENT' | 'GROUP' | 'ROW')? ('ON' ID)?
    ;
alterTablespace_ilm_inmemory_policy
    : ('SET' 'INMEMORY' alterTablespace_inmemory_memcompress? alterTablespace_inmemory_priority? | 'MODIFY' 'INMEMORY' alterTablespace_inmemory_memcompress | 'NO' 'INMEMORY') 'SEGMENT' 'AFTER' alterTablespace_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION')
    ;
alterTablespace_table_compression
    : ('COMPRESS' 'BASIC'? | 'COMPRESS' 'FOR' ('OLTP' | 'QUERY' ('LOW' | 'HIGH')? | 'ARCHIVE' ('LOW' | 'HIGH')?) | 'NOCOMPRESS')
    ;
alterTablespace_ilm_time_period
    : NUMBER ('DAY' | 'DAYS' | 'MONTH' | 'MONTHS' | 'YEAR' | 'YEARS')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-TRIGGER.html
alterTrigger
    : 'ALTER' 'TRIGGER' ('IF' 'EXISTS')? (ID '.')? ID (alterTrigger_trigger_compile_clause | 'ENABLE' | 'DISABLE' | 'RENAME' 'TO' ID | 'EDITIONABLE' | 'NONEDITIONABLE')
    ;
alterTrigger_trigger_compile_clause
    : 'COMPILE' 'DEBUG'? (alterTrigger_compiler_parameters_clause alterTrigger_compiler_parameters_clause?)* ('REUSE' 'SETTINGS')?
    ;
alterTrigger_compiler_parameters_clause
    : ID '=' EXPR
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-TYPE.html
// fix: left out where optional, since no quoted page defines it or what it needs: alter_type_replace_clause
alterType
    : 'ALTER' 'TYPE' ('IF' 'EXISTS')? (ID '.')? ID alterType_alter_type_clause ('EDITIONABLE' | 'NONEDITIONABLE')?
    ;
alterType_alter_type_clause
    : (alterType_alter_type_compile_clause | (alterType_alter_method_spec | alterType_alter_attribute_definition | alterType_alter_collection_clauses | 'MODIFY' alterType_modifier_clause) alterType_dependent_handling_clause?)
    ;
alterType_alter_type_compile_clause
    : 'COMPILE' 'DEBUG'? ('SPECIFICATION' | 'BODY')? (alterType_compiler_parameters_clause alterType_compiler_parameters_clause?)* ('REUSE' 'SETTINGS')?
    ;
alterType_alter_method_spec
    : ('ADD' | 'DROP') (MAP_ORDER_SPEC | SUBPROGRAM_SPEC) (',' (MAP_ORDER_SPEC | SUBPROGRAM_SPEC))*
    ;
alterType_alter_attribute_definition
    : ('ADD' | 'DROP' | 'MODIFY') 'ATTRIBUTE' (ID DATA_TYPE | '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')')
    ;
alterType_alter_collection_clauses
    : 'MODIFY' ('LIMIT' NUMBER | 'ELEMENT' 'TYPE' DATA_TYPE)
    ;
alterType_modifier_clause
    : ('NOT' 'INSTANTIABLE' | 'INSTANTIABLE' | 'NOT' 'FINAL' | 'FINAL')
    ;
alterType_dependent_handling_clause
    : ('INVALIDATE' | 'CASCADE' (('NOT' 'INCLUDING' 'TABLE' 'DATA' | 'INCLUDING' 'TABLE' 'DATA'))? 'FORCE'? (('CONVERT' 'TO' 'SUBSTITUTABLE' | 'NOT' 'SUBSTITUTABLE'))? alterType_exceptions_clause?)
    ;
alterType_compiler_parameters_clause
    : ID '=' EXPR
    ;
alterType_exceptions_clause
    : 'EXCEPTIONS' 'INTO' (ID '.')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-USER.html
alterUser
    : 'ALTER' 'USER' ID ('IF' 'EXISTS')? (alterUser_alter_identified_clause | alterUser_proxy_clause | alterUser_alter_user_clauses)
    ;
alterUser_alter_identified_clause
    : 'IDENTIFIED' ('BY' ID ('REPLACE' ID)? | 'EXTERNALLY' ('AS' STRING)? | 'GLOBALLY' ('AS' STRING)?)
    | 'NO' 'AUTHENTICATION'
    ;
alterUser_proxy_clause
    : ('GRANT' 'CONNECT' 'THROUGH' ('ENTERPRISE' 'USERS' | ID alterUser_proxy_clause_options?) | 'REVOKE' 'CONNECT' 'THROUGH' ('ENTERPRISE' 'USERS' | ID))
    ;
alterUser_alter_user_clauses
    : ('DEFAULT' 'TABLESPACE' ID)? ('TEMPORARY' 'TABLESPACE' ID)? ('LOCAL' 'TEMPORARY' 'TABLESPACE' ID)? (alterUser_quota_clause (',' alterUser_quota_clause)?)* ('PROFILE' ID)? ('DEFAULT' 'ROLE' (ID (',' ID)* | 'ALL' ('EXCEPT' ID (',' ID)?)* | 'NONE'))? ('PASSWORD' 'EXPIRE')? ('EXPIRE' 'PASSWORD' 'ROLLOVER' 'PERIOD')? ('ACCOUNT' ('LOCK' | 'UNLOCK'))? ('ENABLE' 'EDITIONS' ('FOR' ID (',' ID)?)* 'FORCE'?)? ('HTTP'? 'DIGEST' ('ENABLE' | 'DISABLE'))? ('CONTAINER' '=' ('ALL' | 'CURRENT'))? (('ENABLE' | 'DISABLE') 'DICTIONARY' 'PROTECTION')? ('DEFAULT' 'COLLATION' ID)? alterUser_container_data_clause?
    ;
alterUser_proxy_clause_options
    : ('WITH' ('ROLE' (ID (',' ID)* | 'ALL' 'EXCEPT' ID (',' ID)*) | 'NO' 'ROLES'))? ('AUTHENTICATION' 'REQUIRED')?
    ;
alterUser_quota_clause
    : 'QUOTA' (alterUser_size_clause | 'UNLIMITED') 'ON' ID
    ;
alterUser_container_data_clause
    : ('SET' | 'ADD' | 'REMOVE') 'CONTAINER_DATA' ('(' ID (',' ID)* ')' | 'ALL' | 'DEFAULT') ('FOR' (ID '.')? ID)?
    ;
alterUser_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T' | 'P')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-VIEW.html
// fix: left out where optional, since no quoted page defines it or what it needs: create_index_statement, storage_clause
alterView
    : 'ALTER' 'VIEW' ('IF' 'EXISTS')? (ID '.')? ID ('COMPILE' | ('READ' 'ONLY' | 'READ' 'WRITE') | 'ADD' alterView_out_of_line_constraint | 'MODIFY' 'CONSTRAINT' ID ('RELY' | 'NORELY') | 'DROP' 'CONSTRAINT' ID | ('EDITIONABLE' | 'NONEDITIONABLE') | 'ANNOTATIONS' '(' alterView_annotations_clause ')')
    ;
alterView_out_of_line_constraint
    : ('CONSTRAINT' ID)? ('UNIQUE' '(' ID (',' ID)* ')' | 'PRIMARY' 'KEY' '(' ID (',' ID)* ')' | 'FOREIGN' 'KEY' '(' ID (',' ID)* ')' 'REFERENCES' (ID '.')? ID ('(' ID (',' ID)* ')')? ('ON' 'DELETE' ('CASCADE' | 'SET' 'NULL'))? | 'CHECK' '(' CONDITION ')') alterView_constraint_state?
    ;
alterView_annotations_clause
    : ('ADD' 'OR' 'REPLACE' | 'DROP') 'ANNOTATIONS' '(' alterView_annotation (',' alterView_annotation)* ')'
    ;
alterView_constraint_state
    : ('NOT'? 'DEFERRABLE' ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? | 'INITIALLY' ('DEFERRED' | 'IMMEDIATE') ('NOT'? 'DEFERRABLE')?)? ('RELY' | 'NORELY')? alterView_using_index_clause? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? alterView_exceptions_clause?
    ;
alterView_annotation
    : ID EXPR?
    ;
alterView_using_index_clause
    : 'USING' 'INDEX' ((ID '.')? ID | alterView_index_attributes)
    ;
alterView_exceptions_clause
    : 'EXCEPTIONS' 'INTO' (ID '.')? ID
    ;
alterView_index_attributes
    : (alterView_physical_attributes_clause | 'TABLESPACE' (ID | 'DEFAULT') | alterView_logging_clause | alterView_index_compression | alterView_partial_index_clause | alterView_parallel_clause)+
    ;
alterView_physical_attributes_clause
    : ('PCTUSED' NUMBER)? ('PCTFREE' NUMBER)? ('INITRANS' NUMBER)?
    ;
alterView_logging_clause
    : ('LOGGING' | 'NOLOGGING')
    ;
alterView_index_compression
    : (alterView_prefix_compression | alterView_advanced_index_compression | 'NOCOMPRESS')
    ;
alterView_partial_index_clause
    : 'INDEXING' ('FULL' | 'PARTIAL')
    ;
alterView_parallel_clause
    : ('PARALLEL' NUMBER? | 'NOPARALLEL')
    ;
alterView_prefix_compression
    : 'COMPRESS' NUMBER?
    ;
alterView_advanced_index_compression
    : 'COMPRESS' 'ADVANCED' ('LOW' | 'HIGH')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ANALYZE.html
analyze
    : 'ANALYZE' ('TABLE' (ID '.')? ID analyze_partition_extension_clause? | 'INDEX' (ID '.')? ID analyze_partition_extension_clause? | 'CLUSTER' (ID '.')? ID) (analyze_validation_clauses | 'DELETE' 'SYSTEM'? 'STATISTICS')
    ;
analyze_partition_extension_clause
    : ('PARTITION' '(' ID ')' | 'SUBPARTITION' '(' ID ')')
    ;
analyze_validation_clauses
    : ('VALIDATE' 'REF' 'UPDATE' ('SET' 'DANGLING' 'TO' 'NULL')? | 'VALIDATE' 'STRUCTURE' ('CASCADE' 'FAST'?)? ('ONLINE' | 'OFFLINE')? analyze_into_clause? | 'LIST' 'CHAINED' 'ROWS' analyze_into_clause?)
    ;
analyze_into_clause
    : 'INTO' (ID '.')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ASSOCIATE-STATISTICS.html
associateStatistics
    : 'ASSOCIATE' 'STATISTICS' 'WITH' (associateStatistics_column_association | associateStatistics_function_association) associateStatistics_using_statistics_type associateStatistics_default_cost_clause? associateStatistics_default_selectivity_clause? associateStatistics_storage_table_clause?
    ;
associateStatistics_column_association
    : 'COLUMNS' (ID '.')? ID '.' ID (',' (ID '.')? ID '.' ID)*
    ;
associateStatistics_function_association
    : ('FUNCTIONS' (ID '.')? ID (',' (ID '.')? ID)* | 'PACKAGES' (ID '.')? ID (',' (ID '.')? ID)* | 'TYPES' (ID '.')? ID (',' (ID '.')? ID)* | 'DOMAIN' 'INDEXES' (ID '.')? ID (',' (ID '.')? ID)* | 'INDEXTYPES' (ID '.')? ID (',' (ID '.')? ID)*)
    ;
associateStatistics_using_statistics_type
    : 'USING' ((ID '.')? ID | 'NULL')
    ;
associateStatistics_default_cost_clause
    : 'DEFAULT' 'COST' '(' NUMBER ',' NUMBER ',' NUMBER ')'
    ;
associateStatistics_default_selectivity_clause
    : 'DEFAULT' 'SELECTIVITY' NUMBER
    ;
associateStatistics_storage_table_clause
    : 'WITH' ('SYSTEM' 'MANAGED' 'STORAGE' 'TABLES' | 'USER' 'MANAGED' 'STORAGE' 'TABLES')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/19/sqlrf/AUDIT-Traditional-Auditing.html
// fix: left out where optional, since no quoted page defines it or what it needs: sql_statement_shortcut
auditTraditionalAuditing
    : 'AUDIT' (auditTraditionalAuditing_audit_operation_clause (auditTraditionalAuditing_auditing_by_clause | auditTraditionalAuditing_auditing_on_clause)? | auditTraditionalAuditing_audit_schema_object_clause) ('BY' ('SESSION' | 'ACCESS'))? ('WHENEVER' 'NOT'? 'SUCCESSFUL')? ('CONTAINER' '=' ('CURRENT' | 'ALL'))?
    ;
auditTraditionalAuditing_audit_operation_clause
    : (ID | 'ALL' | 'ALL' 'STATEMENTS' | 'ALL' 'PRIVILEGES') (',' ID)*
    ;
auditTraditionalAuditing_auditing_by_clause
    : 'BY' (ID (',' ID)* | 'SESSION' 'CURRENT')
    ;
auditTraditionalAuditing_auditing_on_clause
    : 'ON' (ID '.')? ID
    | 'ON' 'DEFAULT'
    | 'ON' 'DIRECTORY' ID
    | 'ON' 'MINING' 'MODEL' ID
    | 'ON' 'SQL' 'TRANSLATION' 'PROFILE' ID
    | 'NETWORK'
    | 'DIRECT_PATH' 'LOAD'
    ;
auditTraditionalAuditing_audit_schema_object_clause
    : ID (',' ID)* auditTraditionalAuditing_auditing_on_clause
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/AUDIT-Unified-Auditing.html
auditUnifiedAuditing
    : 'AUDIT' ('POLICY' ID (('BY' | 'EXCEPT') ID (',' ID)?)* (('BY' | 'EXCEPT') 'USERS' 'WITH' 'ROLE' ID (',' ID)?)* ('WHENEVER' 'NOT'? 'SUCCESSFUL')? | 'CONTEXT' 'NAMESPACE' ID 'ATTRIBUTES' ID (',' ID)* ('BY' ID (',' ID)?)*)
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CALL.html
call
    : 'CALL' (call_routine_clause | call_object_access_expression) ('INTO' ':' ID ('INDICATOR'? ':' ID)?)?
    ;
call_routine_clause
    : (ID '.')? (ID '.')? ID ('@' ID)? '(' (ID (',' ID)?)* ')'
    ;
call_object_access_expression
    : '(' EXPR ')' '.' ID '(' (ID (',' ID)?)* ')'
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/COMMENT.html
comment
    : 'COMMENT' 'ON' ('AUDIT' 'POLICY' ID | 'COLUMN' (ID '.')? ID '.' ID | 'EDITION' ID | 'INDEXTYPE' (ID '.')? ID | 'MATERIALIZED' 'VIEW' (ID '.')? ID | 'MINING' 'MODEL' (ID '.')? ID | 'OPERATOR' (ID '.')? ID | 'PROPERTY' 'GRAPH' (ID '.')? ID | 'TABLE' (ID '.')? ID) 'IS' STRING
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/COMMIT.html
commit
    : 'COMMIT' 'WORK'? ('COMMENT' STRING)? ('WRITE' ('WAIT' | 'NOWAIT')? ('IMMEDIATE' | 'BATCH')?)? ('FORCE' STRING (',' NUMBER)?)?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-ANALYTIC-VIEW.html
createAnalyticView
    : 'CREATE' ('OR' 'REPLACE')? (('FORCE' | 'NOFORCE'))? 'ANALYTIC' 'VIEW' (ID '.')? ID ('SHARING' '=' ('METADATA' | 'NONE'))? createAnalyticView_classification_clause? createAnalyticView_using_clause createAnalyticView_dim_by_clause createAnalyticView_measures_clause createAnalyticView_default_measure_clause? createAnalyticView_default_aggregate_clause? createAnalyticView_cache_clause? createAnalyticView_fact_columns_clause? createAnalyticView_qry_transform_clause?
    | 'CREATE' (('FORCE' | 'NOFORCE'))? 'ANALYTIC' 'VIEW' ('IF' 'NOT' 'EXISTS')? (ID '.')? ID ('SHARING' '=' ('METADATA' | 'NONE'))? createAnalyticView_classification_clause? createAnalyticView_using_clause createAnalyticView_dim_by_clause createAnalyticView_measures_clause createAnalyticView_default_measure_clause? createAnalyticView_default_aggregate_clause? createAnalyticView_cache_clause? createAnalyticView_fact_columns_clause? createAnalyticView_qry_transform_clause?
    ;
createAnalyticView_classification_clause
    : ('CAPTION' STRING | 'DESCRIPTION' STRING | 'CLASSIFICATION' ID ('LANGUAGE' ID)? 'VALUE' STRING) createAnalyticView_classification_clause?
    ;
createAnalyticView_using_clause
    : 'USING' createAnalyticView_source_clause
    ;
createAnalyticView_dim_by_clause
    : 'DIMENSION' 'BY' '(' createAnalyticView_dim_key (',' createAnalyticView_dim_key)* ')'
    ;
createAnalyticView_measures_clause
    : 'MEASURES' '(' createAnalyticView_av_measure (',' createAnalyticView_av_measure)* ')'
    ;
createAnalyticView_default_measure_clause
    : 'DEFAULT' 'MEASURE' ID
    ;
createAnalyticView_default_aggregate_clause
    : 'DEFAULT' 'AGGREGATE' 'BY' createAnalyticView_aggr_function
    ;
createAnalyticView_cache_clause
    : 'CACHE' '(' createAnalyticView_cache_specification (',' createAnalyticView_cache_specification)* ')'
    ;
createAnalyticView_fact_columns_clause
    : 'FACT' 'COLUMNS' '(' ID ('AS'? ID)? (',' ID ('AS'? ID)?)* ')'
    ;
createAnalyticView_qry_transform_clause
    : 'ENABLE' 'QUERY' 'TRANSFORM' (('RELY' | 'NORELY'))?
    ;
createAnalyticView_source_clause
    : 'REMOTE'? (ID '.')? ID ('AS'? ID)?
    ;
createAnalyticView_dim_key
    : createAnalyticView_dim_ref 'KEY' ID (',' ID)* 'REFERENCES' (ID | 'DISTINCT' ID) 'HIERARCHIES' '(' createAnalyticView_hier_ref (',' createAnalyticView_hier_ref)* ')'
    ;
createAnalyticView_av_measure
    : ID (createAnalyticView_base_meas_clause | createAnalyticView_calc_meas_clause) createAnalyticView_classification_clause?
    ;
createAnalyticView_aggr_function
    : ('APPROX_COUNT_DISTINCT' | 'APPROX_COUNT_DISTINCT_AGG' | 'AVG' | 'COUNT' | 'MAX' | 'MIN' | 'STDDEV' | 'STDDEV_POP' | 'STDDEV_SAMP' | 'SUM' | 'VAR_POP' | 'VAR_SAMP' | 'VARIANCE')
    ;
createAnalyticView_cache_specification
    : 'MEASURE' 'GROUP' ('ALL' | ID (',' ID)*) 'LEVELS' '(' createAnalyticView_level_grouping (',' createAnalyticView_level_grouping)* ')'
    ;
createAnalyticView_dim_ref
    : (ID '.')? ID ('AS'? ID)?
    ;
createAnalyticView_hier_ref
    : (ID '.')? ID ('AS'? ID)? 'DEFAULT'?
    ;
createAnalyticView_base_meas_clause
    : ('FACT' ID)? createAnalyticView_meas_aggregate_clause?
    ;
createAnalyticView_calc_meas_clause
    : 'AS' '(' EXPR ')'
    ;
createAnalyticView_level_grouping
    : createAnalyticView_level_group_type? '(' createAnalyticView_level_specification (',' createAnalyticView_level_specification)* ')'
    ;
createAnalyticView_meas_aggregate_clause
    : createAnalyticView_aggr_function
    ;
createAnalyticView_level_group_type
    : 'USING' (ID '.')? ID
    ;
createAnalyticView_level_specification
    : (ID '.')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-ATTRIBUTE-DIMENSION.html
createAttributeDimension
    : 'CREATE' ('OR' 'REPLACE')? (('FORCE' | 'NOFORCE'))? 'ATTRIBUTE' 'DIMENSION' (ID '.')? ID ('SHARING' '=' ('METADATA' | 'NONE'))? createAttributeDimension_classification_clause? ('DIMENSION' 'TYPE' ('STANDARD' | 'TIME'))? createAttributeDimension_attr_dim_using_clause createAttributeDimension_attributes_clause createAttributeDimension_attr_dim_level_clause createAttributeDimension_attr_dim_level_clause* createAttributeDimension_all_clause?
    | 'CREATE' (('FORCE' | 'NOFORCE'))? 'ATTRIBUTE' 'DIMENSION' ('IF' 'NOT' 'EXISTS')? (ID '.')? ID ('SHARING' '=' ('METADATA' | 'NONE'))? createAttributeDimension_classification_clause? ('DIMENSION' 'TYPE' ('STANDARD' | 'TIME'))? createAttributeDimension_attr_dim_using_clause createAttributeDimension_attributes_clause createAttributeDimension_attr_dim_level_clause createAttributeDimension_attr_dim_level_clause* createAttributeDimension_all_clause?
    ;
createAttributeDimension_classification_clause
    : ('CAPTION' STRING | 'DESCRIPTION' STRING | 'CLASSIFICATION' ID ('LANGUAGE' ID)? 'VALUE' STRING) createAttributeDimension_classification_clause?
    ;
createAttributeDimension_attr_dim_using_clause
    : 'USING' createAttributeDimension_source_clause (',' createAttributeDimension_source_clause)*
    ;
createAttributeDimension_attributes_clause
    : 'ATTRIBUTES' '(' createAttributeDimension_attr_dim_attribute_clause (',' createAttributeDimension_attr_dim_attribute_clause)* ')'
    ;
createAttributeDimension_attr_dim_level_clause
    : 'LEVEL' ID ('LEVEL' 'TYPE' ('STANDARD' | 'YEARS' | 'HALF_YEARS' | 'QUARTERS' | 'MONTHS' | 'WEEKS' | 'DAYS' | 'HOURS' | 'MINUTES' | 'SECONDS'))? createAttributeDimension_classification_clause? createAttributeDimension_key_clause createAttributeDimension_alternate_key_clause? ('MEMBER' 'NAME' EXPR)? ('MEMBER' 'CAPTION' EXPR)? ('MEMBER' 'DESCRIPTION' EXPR)? createAttributeDimension_dim_order_clause? ('DETERMINES' '(' ID (',' ID)* ')')?
    ;
createAttributeDimension_all_clause
    : 'ALL' ('MEMBER' 'NAME' EXPR)? ('MEMBER' 'CAPTION' EXPR)? ('MEMBER' 'DESCRIPTION' EXPR)?
    ;
createAttributeDimension_source_clause
    : 'REMOTE'? (ID '.')? ID ('AS'? ID)?
    | createAttributeDimension_join_path_clause
    ;
createAttributeDimension_attr_dim_attribute_clause
    : ID ('AS' ID)? createAttributeDimension_classification_clause?
    ;
createAttributeDimension_key_clause
    : 'KEY' (ID ('NOT' 'NULL' | 'SKIP' 'WHEN' 'NULL')? | '(' ID (',' ID)* ')')
    ;
createAttributeDimension_alternate_key_clause
    : 'ALTERNATE' 'KEY' (ID | '(' ID (',' ID)* ')')
    ;
createAttributeDimension_dim_order_clause
    : 'ORDER' 'BY' (ID ('ASC' | 'DESC')? | '(' ID ('ASC' | 'DESC')? (',' ID ('ASC' | 'DESC')?)* ')')
    ;
createAttributeDimension_join_path_clause
    : 'JOIN' 'PATH' ID 'ON' createAttributeDimension_join_condition
    ;
createAttributeDimension_join_condition
    : createAttributeDimension_join_condition_elem ('AND' createAttributeDimension_join_condition_elem)*
    ;
createAttributeDimension_join_condition_elem
    : (ID '.')? ID '=' (ID '.')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-AUDIT-POLICY-Unified-Auditing.html
createAuditPolicyUnifiedAuditing
    : 'CREATE' 'AUDIT' 'POLICY' ID createAuditPolicyUnifiedAuditing_privilege_audit_clause? createAuditPolicyUnifiedAuditing_action_audit_clause? createAuditPolicyUnifiedAuditing_role_audit_clause? ('WHEN' STRING 'EVALUATE' 'PER' ('STATEMENT' | 'SESSION' | 'INSTANCE'))? ('ONLY' 'TOPLEVEL')? ('CONTAINER' '=' ('ALL' | 'CURRENT'))?
    ;
createAuditPolicyUnifiedAuditing_privilege_audit_clause
    : 'PRIVILEGES' ID (',' ID)*
    ;
createAuditPolicyUnifiedAuditing_action_audit_clause
    : 'ACTIONS' createAuditPolicyUnifiedAuditing_standard_actions? createAuditPolicyUnifiedAuditing_component_actions?
    ;
createAuditPolicyUnifiedAuditing_role_audit_clause
    : 'ROLES' ID (',' ID)*
    ;
createAuditPolicyUnifiedAuditing_standard_actions
    : (ID ('(' ID (',' ID)* ')')? 'ON' ((ID '.')? ID | 'DIRECTORY' ID | 'MINING' 'MODEL' (ID '.')? ID) | 'ALL' 'ON' ((ID '.')? ID | 'DIRECTORY' ID | 'MINING' 'MODEL' (ID '.')? ID) | ID | 'ALL') (',' (ID ('(' ID (',' ID)* ')')? 'ON' ((ID '.')? ID | 'DIRECTORY' ID | 'MINING' 'MODEL' (ID '.')? ID) | 'ALL' 'ON' ((ID '.')? ID | 'DIRECTORY' ID | 'MINING' 'MODEL' (ID '.')? ID) | ID | 'ALL'))*
    ;
createAuditPolicyUnifiedAuditing_component_actions
    : 'COMPONENT' '=' ('DATAPUMP' (ID | 'ALL') (',' (ID | 'ALL'))* | 'DIRECT_LOAD' (ID | 'ALL') (',' (ID | 'ALL'))* | 'OLS' (ID | 'ALL') (',' (ID | 'ALL'))* | 'XS' (ID | 'ALL') (',' (ID | 'ALL'))* | 'DV' (ID ('ON' ID)? | 'ALL') (',' (ID ('ON' ID)? | 'ALL'))* | 'SQL_FIREWALL' | 'PROTOCOL' ('HTTP' | 'FTP' | 'AUTHENTICATION'))
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-CLUSTER.html
// fix: left out where optional, since no quoted page defines it or what it needs: lob_parameters, storage_clause
createCluster
    : 'CREATE' 'CLUSTER' ('IF' 'NOT' 'EXISTS')? (ID '.')? ID ('SHARING' '=' ('METADATA' | 'NONE'))? '(' ID DATA_TYPE ('COLLATE' ID)? 'SORT'? (',' ID DATA_TYPE ('COLLATE' ID)? 'SORT'?)* ')' createCluster_physical_attributes_clause? ('SIZE' createCluster_size_clause)? ('TABLESPACE' ID)? (('INDEX' | 'HASHKEYS' NUMBER ('HASH' 'IS' EXPR)?))? ('SINGLE' 'TABLE')? createCluster_parallel_clause? (('NOROWDEPENDENCIES' | 'ROWDEPENDENCIES'))? (('CACHE' | 'NOCACHE'))? createCluster_cluster_range_partitions?
    ;
createCluster_physical_attributes_clause
    : ('PCTFREE' NUMBER)? ('PCTUSED' NUMBER)? ('INITRANS' NUMBER)?
    ;
createCluster_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T' | 'P')?
    ;
createCluster_parallel_clause
    : ('NOPARALLEL' | 'PARALLEL' NUMBER?)
    ;
createCluster_cluster_range_partitions
    : 'PARTITION' 'BY' 'RANGE' '(' ID (',' ID)* ')' '(' 'PARTITION' ID? 'VALUES' 'LESS' 'THAN' '(' (EXPR | 'MAXVALUE') (',' (EXPR | 'MAXVALUE'))* ')' createCluster_table_partition_description? (',' 'PARTITION' ID? 'VALUES' 'LESS' 'THAN' '(' (EXPR | 'MAXVALUE') (',' (EXPR | 'MAXVALUE'))* ')' createCluster_table_partition_description?)* ')'
    ;
createCluster_table_partition_description
    : createCluster_deferred_segment_creation? createCluster_read_only_clause? createCluster_indexing_clause? createCluster_segment_attributes_clause? (createCluster_table_compression | createCluster_prefix_compression)? createCluster_inmemory_clause? createCluster_ilm_clause? ('OVERFLOW' createCluster_segment_attributes_clause?)?
    ;
createCluster_deferred_segment_creation
    : 'SEGMENT' 'CREATION' ('IMMEDIATE' | 'DEFERRED')
    ;
createCluster_read_only_clause
    : ('READ' 'ONLY' | 'READ' 'WRITE')
    ;
createCluster_indexing_clause
    : 'INDEXING' ('ON' | 'OFF')
    ;
createCluster_segment_attributes_clause
    : createCluster_physical_attributes_clause? ('TABLESPACE' ID)? createCluster_logging_clause?
    ;
createCluster_table_compression
    : ('COMPRESS' 'BASIC'? | 'COMPRESS' 'FOR' ('OLTP' | 'QUERY' ('LOW' | 'HIGH')? | 'ARCHIVE' ('LOW' | 'HIGH')?) | 'NOCOMPRESS')
    ;
createCluster_prefix_compression
    : 'COMPRESS' NUMBER?
    ;
createCluster_inmemory_clause
    : ('INMEMORY' createCluster_inmemory_attributes? | 'NO' 'INMEMORY')
    ;
createCluster_ilm_clause
    : 'ILM' ('ADD' 'POLICY' createCluster_ilm_policy_clause | 'DELETE' 'POLICY' ID | 'DELETE_ALL' | 'ENABLE_ALL' | 'DISABLE_ALL')
    ;
createCluster_logging_clause
    : ('LOGGING' | 'NOLOGGING')
    ;
createCluster_inmemory_attributes
    : createCluster_inmemory_memcompress? createCluster_inmemory_priority? createCluster_inmemory_distribute? createCluster_inmemory_duplicate?
    ;
createCluster_ilm_policy_clause
    : (createCluster_ilm_compression_policy | createCluster_ilm_tiering_policy | createCluster_ilm_inmemory_policy)
    ;
createCluster_inmemory_memcompress
    : 'MEMCOMPRESS' 'FOR' ('QUERY' ('LOW' | 'HIGH')? | 'OLTP' | 'ARCHIVE' ('LOW' | 'HIGH')?)
    ;
createCluster_inmemory_priority
    : 'PRIORITY' ('LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL')
    ;
createCluster_inmemory_distribute
    : 'DISTRIBUTE' 'BY' ('ROWID' ('ACROSS' NUMBER?)? | 'PARTITION' | 'SUBPARTITION')
    ;
createCluster_inmemory_duplicate
    : 'DUPLICATE' ('ALL' | 'FOR' 'SERVICE' ID (',' ID)* | 'NONE')
    ;
createCluster_ilm_compression_policy
    : (createCluster_table_compression | 'ROW' 'STORE' 'COMPRESS' 'ADVANCED'? | 'COLUMN' 'STORE' 'COMPRESS' ('FOR' 'QUERY' ('LOW' | 'HIGH')?)?) ('SEGMENT' | 'GROUP' | 'ROW') 'AFTER' createCluster_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION')
    ;
createCluster_ilm_tiering_policy
    : 'TIER' 'TO' ID ('SEGMENT' | 'GROUP' | 'ROW')? ('ON' ID)?
    ;
createCluster_ilm_inmemory_policy
    : ('SET' 'INMEMORY' createCluster_inmemory_memcompress? createCluster_inmemory_priority? | 'MODIFY' 'INMEMORY' createCluster_inmemory_memcompress | 'NO' 'INMEMORY') 'SEGMENT' 'AFTER' createCluster_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION')
    ;
createCluster_ilm_time_period
    : NUMBER ('DAY' | 'DAYS' | 'MONTH' | 'MONTHS' | 'YEAR' | 'YEARS')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-CONTEXT.html
createContext
    : 'CREATE' ('OR' 'REPLACE')? 'CONTEXT' ID 'USING' (ID '.')? ID ('SHARING' '=' ('METADATA' | 'NONE'))? ('INITIALIZED' ('EXTERNALLY' | 'GLOBALLY'))? ('ACCESSED' 'GLOBALLY')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-CONTROLFILE.html
createControlfile
    : 'CREATE' 'CONTROLFILE' 'REUSE'? ('SET' 'DATABASE' ID | 'DATABASE' ID) createControlfile_logfile_clause (('RESETLOGS' | 'NORESETLOGS'))? ('DATAFILE' createControlfile_file_specification (',' createControlfile_file_specification)?)* ('MAXLOGFILES' NUMBER)? ('MAXLOGMEMBERS' NUMBER)? ('MAXLOGHISTORY' NUMBER)? ('MAXDATAFILES' NUMBER)? ('MAXINSTANCES' NUMBER)? (('ARCHIVELOG' | 'NOARCHIVELOG'))? ('FORCE' 'LOGGING')? ('SET' 'STANDBY' 'NOLOGGING' 'FOR' ('DATA' 'AVAILABILITY' | 'LOAD' 'PERFORMANCE'))? createControlfile_character_set_clause?
    ;
createControlfile_logfile_clause
    : 'LOGFILE' ('GROUP' NUMBER)? createControlfile_file_specification (',' ('GROUP' NUMBER)? createControlfile_file_specification)*
    ;
createControlfile_file_specification
    : STRING? ('SIZE' createControlfile_size_clause)? 'REUSE'? createControlfile_autoextend_clause?
    ;
createControlfile_character_set_clause
    : 'CHARACTER' 'SET' ID
    ;
createControlfile_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T' | 'P')?
    ;
createControlfile_autoextend_clause
    : 'AUTOEXTEND' ('OFF' | 'ON' ('NEXT' createControlfile_size_clause)? createControlfile_maxsize_clause?)
    ;
createControlfile_maxsize_clause
    : 'MAXSIZE' (createControlfile_size_clause | 'UNLIMITED')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-DATABASE-LINK.html
createDatabaseLink
    : 'CREATE' 'SHARED'? 'PUBLIC'? 'DATABASE' 'LINK' ('IF' 'NOT' 'EXISTS')? ID ('CONNECT' 'TO' ('CURRENT_USER' | ID 'IDENTIFIED' 'BY' ID createDatabaseLink_dblink_authentication?) | 'CONNECT' 'WITH' ID)? createDatabaseLink_dblink_authentication? 'USING' STRING
    ;
createDatabaseLink_dblink_authentication
    : 'AUTHENTICATED' 'BY' ID 'IDENTIFIED' 'BY' ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-DATABASE.html
createDatabase
    : 'CREATE' 'DATABASE' ID? ('USER' 'SYS' 'IDENTIFIED' 'BY' ID)? ('USER' 'SYSTEM' 'IDENTIFIED' 'BY' ID)? ('CONTROLFILE' 'REUSE')? ('MAXDATAFILES' NUMBER)? ('MAXINSTANCES' NUMBER)? ('CHARACTER' 'SET' ID)? ('NATIONAL' 'CHARACTER' 'SET' ('AL16UTF16' | 'UTF8'))? ('SET' 'DEFAULT' ('BIGFILE' | 'SMALLFILE') 'TABLESPACE')? createDatabase_database_logging_clauses? createDatabase_tablespace_clauses? createDatabase_set_time_zone_clause? createDatabase_enable_pluggable_database? ('USING' 'MIRROR' 'COPY' ID)?
    ;
createDatabase_database_logging_clauses
    : ('LOGFILE' ('GROUP' NUMBER)? createDatabase_file_specification (',' ('GROUP' NUMBER)? createDatabase_file_specification)?)* ('MAXLOGFILES' NUMBER)? ('MAXLOGMEMBERS' NUMBER)? ('MAXLOGHISTORY' NUMBER)? (('ARCHIVELOG' | 'NOARCHIVELOG'))? ('FORCE' 'LOGGING')? ('SET' 'STANDBY' 'NOLOGGING' 'FOR' ('DATA' 'AVAILABILITY' | 'LOAD' 'PERFORMANCE'))?
    ;
createDatabase_tablespace_clauses
    : ('EXTENT' 'MANAGEMENT' 'LOCAL' (('AUTOALLOCATE' | 'UNIFORM' ('SIZE' createDatabase_size_clause)?))?)? ('DATAFILE' createDatabase_file_specification (',' createDatabase_file_specification)?)* ('SYSAUX' 'DATAFILE' createDatabase_file_specification (',' createDatabase_file_specification)?)* createDatabase_default_tablespace? createDatabase_default_temp_tablespace? createDatabase_undo_tablespace? ('USER_DATA' 'TABLESPACE' ID (('BIGFILE' | 'SMALLFILE'))? 'DATAFILE' createDatabase_file_specification (',' createDatabase_file_specification)?)*
    ;
createDatabase_set_time_zone_clause
    : 'SET' 'TIME_ZONE' '=' (STRING | ID)
    ;
createDatabase_enable_pluggable_database
    : 'ENABLE' 'PLUGGABLE' 'DATABASE' ('SEED' createDatabase_file_name_convert? ('SYSTEM' 'DATAFILES' 'SIZE' createDatabase_size_clause ('AUTOEXTEND' ('ON' ('NEXT' createDatabase_size_clause)? ('MAXSIZE' (createDatabase_size_clause | 'UNLIMITED'))? | 'OFF'))?)? ('SYSAUX' 'DATAFILES' 'SIZE' createDatabase_size_clause ('AUTOEXTEND' ('ON' ('NEXT' createDatabase_size_clause)? ('MAXSIZE' (createDatabase_size_clause | 'UNLIMITED'))? | 'OFF'))?)? createDatabase_undo_mode_clause?)?
    ;
createDatabase_file_specification
    : STRING? ('SIZE' createDatabase_size_clause)? 'REUSE'? createDatabase_autoextend_clause?
    ;
createDatabase_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T' | 'P')?
    ;
createDatabase_default_tablespace
    : 'DEFAULT' 'TABLESPACE' ID (('BIGFILE' | 'SMALLFILE'))? ('DATAFILE' createDatabase_file_specification (',' createDatabase_file_specification)?)* ('EXTENT' 'MANAGEMENT' 'LOCAL' (('AUTOALLOCATE' | 'UNIFORM' ('SIZE' createDatabase_size_clause)?))?)?
    ;
createDatabase_default_temp_tablespace
    : ('DEFAULT' 'TEMPORARY' 'TABLESPACE' | 'DEFAULT' 'LOCAL' 'TEMPORARY' 'TABLESPACE' ('FOR' ('ALL' | 'LEAF'))?) ID (('BIGFILE' | 'SMALLFILE'))? ('TEMPFILE' createDatabase_file_specification (',' createDatabase_file_specification)?)* ('EXTENT' 'MANAGEMENT' 'LOCAL' (('AUTOALLOCATE' | 'UNIFORM' ('SIZE' createDatabase_size_clause)?))?)?
    ;
createDatabase_undo_tablespace
    : 'UNDO' 'TABLESPACE' ID (('BIGFILE' | 'SMALLFILE'))? ('DATAFILE' createDatabase_file_specification (',' createDatabase_file_specification)?)*
    ;
createDatabase_file_name_convert
    : 'FILE_NAME_CONVERT' '=' ('(' STRING ',' STRING (',' STRING ',' STRING)* ')' | 'NONE')
    ;
createDatabase_undo_mode_clause
    : 'LOCAL' 'UNDO' ('ON' | 'OFF')
    ;
createDatabase_autoextend_clause
    : 'AUTOEXTEND' ('OFF' | 'ON' ('NEXT' createDatabase_size_clause)? createDatabase_maxsize_clause?)
    ;
createDatabase_maxsize_clause
    : 'MAXSIZE' (createDatabase_size_clause | 'UNLIMITED')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-DIMENSION.html
createDimension
    : 'CREATE' 'DIMENSION' (ID '.')? ID createDimension_level_clause createDimension_level_clause* (createDimension_hierarchy_clause | createDimension_dimension_join_clause | createDimension_attribute_clause | createDimension_extended_attribute_clause)+
    ;
createDimension_level_clause
    : 'LEVEL' ID 'IS' '(' (ID '.')? ID '.' ID (',' (ID '.')? ID '.' ID)* ')' ('SKIP' 'WHEN' 'NULL')?
    ;
createDimension_hierarchy_clause
    : 'HIERARCHY' ID '(' ID 'CHILD' 'OF' ID ('CHILD' 'OF' ID)* createDimension_dimension_join_clause* ')'
    ;
createDimension_dimension_join_clause
    : 'JOIN' 'KEY' '(' ((ID '.')? ID '.')? ID (',' ((ID '.')? ID '.')? ID)* ')' 'REFERENCES' ID
    ;
createDimension_attribute_clause
    : 'ATTRIBUTE' ID 'DETERMINES' '(' (ID '.')? ID '.' ID (',' (ID '.')? ID '.' ID)* ')'
    ;
createDimension_extended_attribute_clause
    : 'ATTRIBUTE' ID 'LEVEL' ID 'DETERMINES' '(' (ID '.')? ID '.' ID (',' (ID '.')? ID '.' ID)* ')'
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-DIRECTORY.html
createDirectory
    : 'CREATE' ('OR' 'REPLACE' | 'IF' 'NOT' 'EXISTS')? 'DIRECTORY' ID ('SHARING' '=' ('METADATA' | 'NONE'))? 'AS' STRING
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-DISKGROUP.html
createDiskgroup
    : 'CREATE' 'DISKGROUP' ID (('NORMAL' | 'HIGH' | 'FLEX' | 'EXTENDED' | 'EXTERNAL') 'REDUNDANCY')? (('QUORUM' | 'REGULAR') 'FAILGROUP' ID)? 'DISK' createDiskgroup_qualified_disk_clause (',' createDiskgroup_qualified_disk_clause)* (('QUORUM' | 'REGULAR') 'FAILGROUP' ID 'DISK' createDiskgroup_qualified_disk_clause (',' createDiskgroup_qualified_disk_clause)?)* ('ATTRIBUTE' STRING '=' STRING (',' STRING '=' STRING)?)*
    ;
createDiskgroup_qualified_disk_clause
    : STRING ('NAME' ID)? ('SIZE' createDiskgroup_size_clause)? ('FORCE' | 'NOFORCE')?
    ;
createDiskgroup_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T' | 'P')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-EDITION.html
createEdition
    : 'CREATE' 'EDITION' ('IF' 'NOT' 'EXISTS')? ID ('AS' 'CHILD' 'OF' ID)?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-FLASHBACK-ARCHIVE.html
createFlashbackArchive
    : 'CREATE' 'FLASHBACK' 'ARCHIVE' 'DEFAULT'? ID 'TABLESPACE' ID createFlashbackArchive_flashback_archive_quota? (('NO' 'OPTIMIZE' 'DATA' | 'OPTIMIZE' 'DATA'))? createFlashbackArchive_flashback_archive_retention
    ;
createFlashbackArchive_flashback_archive_quota
    : 'QUOTA' NUMBER ('K' | 'M' | 'G' | 'T' | 'P' | 'E')
    ;
createFlashbackArchive_flashback_archive_retention
    : 'RETENTION' NUMBER ('DAY' | 'DAYS' | 'MONTH' | 'MONTHS' | 'YEAR' | 'YEARS')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-FUNCTION.html
// fix: left out where optional, since no quoted page defines it or what it needs: accessible_by_clause, aggregate_clause, default_collation_clause, deterministic_clause, parallel_enable_clause, parameter_declaration, pipelined_clause, result_cache_clause, sql_macro_clause
createFunction
    : 'CREATE' ('OR' 'REPLACE')? ('EDITIONABLE' | 'NONEDITIONABLE')? 'FUNCTION' (ID '.')? ID 'RETURN' DATA_TYPE ('SHARING' '=' ('METADATA' | 'NONE'))? createFunction_invoker_rights_clause* ('IS' | 'AS') (FUNCTION_SOURCE | CALL_SPEC)
    | 'CREATE' ('EDITIONABLE' | 'NONEDITIONABLE')? 'FUNCTION' ('IF' 'NOT' 'EXISTS')? (ID '.')? ID 'RETURN' DATA_TYPE ('SHARING' '=' ('METADATA' | 'NONE'))? createFunction_invoker_rights_clause* ('IS' | 'AS') (FUNCTION_SOURCE | CALL_SPEC)
    ;
createFunction_invoker_rights_clause
    : 'AUTHID' ('CURRENT_USER' | 'DEFINER')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-HIERARCHY.html
createHierarchy
    : 'CREATE' ('OR' 'REPLACE')? (('FORCE' | 'NOFORCE'))? 'HIERARCHY' (ID '.')? ID ('SHARING' '=' ('METADATA' | 'NONE'))? createHierarchy_classification_clause* createHierarchy_hier_using_clause '(' createHierarchy_level_hier_clause ')' createHierarchy_hier_attrs_clause?
    | 'CREATE' (('FORCE' | 'NOFORCE'))? 'HIERARCHY' ('IF' 'NOT' 'EXISTS')? (ID '.')? ID ('SHARING' '=' ('METADATA' | 'NONE'))? createHierarchy_classification_clause* createHierarchy_hier_using_clause '(' createHierarchy_level_hier_clause ')' createHierarchy_hier_attrs_clause?
    ;
createHierarchy_classification_clause
    : ('CAPTION' STRING | 'DESCRIPTION' STRING | 'CLASSIFICATION' ID 'VALUE' STRING ('LANGUAGE' ID)?)
    ;
createHierarchy_hier_using_clause
    : 'USING' (ID '.')? ID
    ;
createHierarchy_level_hier_clause
    : ID createHierarchy_classification_clause* ('CHILD' 'OF' createHierarchy_level_hier_clause)?
    ;
createHierarchy_hier_attrs_clause
    : 'HIERARCHICAL' 'ATTRIBUTES' '(' createHierarchy_hier_attr_clause (',' createHierarchy_hier_attr_clause)* ')'
    ;
createHierarchy_hier_attr_clause
    : createHierarchy_hier_attr_name createHierarchy_classification_clause*
    ;
createHierarchy_hier_attr_name
    : ('HIER_ORDER' | 'DEPTH' | 'IS_LEAF' | 'IS_ROOT' | 'MEMBER_NAME' | 'MEMBER_UNIQUE_NAME' | 'MEMBER_CAPTION' | 'MEMBER_DESCRIPTION' | 'PARENT_LEVEL_NAME' | 'PARENT_UNIQUE_NAME')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-INDEX.html
// fix: left out where optional, since no quoted page defines it or what it needs: index_ilm_clause, storage_clause
createIndex
    : 'CREATE' 'UNIQUE'? 'BITMAP'? 'MULTIVALUE'? 'INDEX' ('IF' 'NOT' 'EXISTS')? (ID '.')? ID (createIndex_cluster_index_clause | createIndex_table_index_clause | createIndex_bitmap_join_index_clause) (('DEFERRED' | 'IMMEDIATE') 'INVALIDATION')?
    ;
createIndex_cluster_index_clause
    : 'ON' 'CLUSTER' (ID '.')? ID createIndex_index_attributes?
    ;
createIndex_table_index_clause
    : 'ON' (ID '.')? ID ID? '(' createIndex_index_expr ('ASC' | 'DESC')? (',' createIndex_index_expr ('ASC' | 'DESC')?)* ')' createIndex_index_properties?
    ;
createIndex_bitmap_join_index_clause
    : 'ON' (ID '.')? ID '(' (ID '.')? ID '.' ID (',' (ID '.')? ID '.' ID)* ')' 'FROM' (ID '.')? ID (',' (ID '.')? ID)* 'WHERE' createIndex_join_condition createIndex_local_partitioned_index? createIndex_index_attributes?
    ;
createIndex_index_attributes
    : createIndex_physical_attributes_clause? ('TABLESPACE' (ID | 'DEFAULT'))? createIndex_index_compression? createIndex_partial_index_clause? createIndex_parallel_clause? (('SORT' | 'NOSORT'))? 'REVERSE'? (('VISIBLE' | 'INVISIBLE'))? createIndex_logging_clause? 'ONLINE'? createIndex_annotations_clause?
    ;
createIndex_index_expr
    : (ID | EXPR)
    ;
createIndex_index_properties
    : (createIndex_global_partitioned_index | createIndex_local_partitioned_index)? createIndex_index_attributes? (createIndex_domain_index_clause | createIndex_xmlindex_clause)?
    ;
createIndex_join_condition
    : createIndex_join_condition_elem ('AND' createIndex_join_condition_elem)*
    ;
createIndex_local_partitioned_index
    : 'LOCAL' (createIndex_on_range_partitioned_table | createIndex_on_list_partitioned_table | createIndex_on_hash_partitioned_table | createIndex_on_comp_partitioned_table)?
    ;
createIndex_physical_attributes_clause
    : ('PCTFREE' NUMBER)? ('INITRANS' NUMBER)? ('MAXTRANS' NUMBER)?
    ;
createIndex_index_compression
    : (createIndex_prefix_compression | createIndex_advanced_index_compression)
    ;
createIndex_partial_index_clause
    : 'INDEXING' ('FULL' | 'PARTIAL')
    ;
createIndex_parallel_clause
    : ('PARALLEL' NUMBER? | 'NOPARALLEL')
    ;
createIndex_logging_clause
    : ('LOGGING' | 'NOLOGGING')
    ;
createIndex_annotations_clause
    : 'ANNOTATIONS' '(' ('ADD' createIndex_annotation | 'DROP' ID (',' ID)* | 'REMOVE' 'ALL') ')'
    ;
createIndex_global_partitioned_index
    : 'GLOBAL' 'PARTITION' 'BY' ('RANGE' '(' ID ')' '(' createIndex_index_partitioning_clause (',' createIndex_index_partitioning_clause)* ')' | 'HASH' '(' ID ')' (createIndex_individual_hash_partitions | createIndex_hash_partitions_by_quantity))
    ;
createIndex_domain_index_clause
    : 'INDEXTYPE' 'IS' (ID '.')? ID createIndex_local_domain_index_clause? createIndex_parallel_clause? ('PARAMETERS' '(' STRING ')')?
    ;
createIndex_xmlindex_clause
    : 'INDEXTYPE' 'IS' 'XDB' '.' 'XMLINDEX' createIndex_local_xmlindex_clause? createIndex_parallel_clause? ('PARAMETERS' '(' STRING ')')?
    ;
createIndex_join_condition_elem
    : (ID '.')? ID '=' (ID '.')? ID
    ;
createIndex_on_range_partitioned_table
    : '(' 'PARTITION' ID? createIndex_segment_attributes_clause? createIndex_index_compression? (('USABLE' | 'UNUSABLE'))? (',' 'PARTITION' ID? createIndex_segment_attributes_clause? createIndex_index_compression? (('USABLE' | 'UNUSABLE'))?)* ')'
    ;
createIndex_on_list_partitioned_table
    : '(' 'PARTITION' ID? createIndex_segment_attributes_clause? createIndex_index_compression? (('USABLE' | 'UNUSABLE'))? (',' 'PARTITION' ID? createIndex_segment_attributes_clause? createIndex_index_compression? (('USABLE' | 'UNUSABLE'))?)* ')'
    ;
createIndex_on_hash_partitioned_table
    : ('(' 'PARTITION' ID? ('TABLESPACE' ID)? createIndex_index_compression? (('USABLE' | 'UNUSABLE'))? (',' 'PARTITION' ID? ('TABLESPACE' ID)? createIndex_index_compression? (('USABLE' | 'UNUSABLE'))?)* ')')? ('STORE' 'IN' '(' ID (',' ID)* ')')?
    ;
createIndex_on_comp_partitioned_table
    : '(' 'PARTITION' ID? createIndex_segment_attributes_clause? createIndex_index_compression? createIndex_index_subpartition_clause? (('USABLE' | 'UNUSABLE'))? (',' 'PARTITION' ID? createIndex_segment_attributes_clause? createIndex_index_compression? createIndex_index_subpartition_clause? (('USABLE' | 'UNUSABLE'))?)* ')'
    ;
createIndex_prefix_compression
    : 'COMPRESS' NUMBER?
    | 'NOCOMPRESS'
    ;
createIndex_advanced_index_compression
    : 'COMPRESS' 'ADVANCED' ('LOW' | 'HIGH')?
    | 'NOCOMPRESS'
    ;
createIndex_annotation
    : ID STRING?
    ;
createIndex_index_partitioning_clause
    : 'PARTITION' ID? 'VALUES' 'LESS' 'THAN' '(' ID ')' createIndex_segment_attributes_clause?
    ;
createIndex_individual_hash_partitions
    : '(' 'PARTITION' ID? createIndex_partitioning_storage_clause? (',' 'PARTITION' ID? createIndex_partitioning_storage_clause?)* ')'
    ;
createIndex_hash_partitions_by_quantity
    : 'PARTITIONS' NUMBER ('STORE' 'IN' '(' ID (',' ID)* ')')?
    ;
createIndex_local_domain_index_clause
    : 'LOCAL' ('(' 'PARTITION' ID? ('PARAMETERS' '(' STRING ')')? (',' 'PARTITION' ID? ('PARAMETERS' '(' STRING ')')?)* ')')?
    ;
createIndex_local_xmlindex_clause
    : 'LOCAL' ('(' 'PARTITION' ID? ('PARAMETERS' '(' STRING ')')? (',' 'PARTITION' ID? ('PARAMETERS' '(' STRING ')')?)* ')')?
    ;
createIndex_segment_attributes_clause
    : createIndex_physical_attributes_clause? ('TABLESPACE' (ID | 'DEFAULT'))? createIndex_logging_clause?
    ;
createIndex_index_subpartition_clause
    : '(' 'SUBPARTITION' ID? ('TABLESPACE' ID)? createIndex_index_compression? (('USABLE' | 'UNUSABLE'))? (',' 'SUBPARTITION' ID? ('TABLESPACE' ID)? createIndex_index_compression? (('USABLE' | 'UNUSABLE'))?)* ')'
    | 'STORE' 'IN' '(' ID (',' ID)* ')'
    ;
createIndex_partitioning_storage_clause
    : ('TABLESPACE' ID)? createIndex_index_compression?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-INDEXTYPE.html
createIndextype
    : 'CREATE' ('OR' 'REPLACE')? 'INDEXTYPE' (ID '.')? ID ('SHARING' '=' ('METADATA' | 'NONE'))? 'FOR' (ID '.')? ID '(' DATA_TYPE (',' DATA_TYPE)* ')' (',' (ID '.')? ID '(' DATA_TYPE (',' DATA_TYPE)* ')')* createIndextype_using_type_clause ('WITH' 'LOCAL' 'RANGE'? 'PARTITION')? createIndextype_storage_table_clause? createIndextype_array_dml_clause?
    | 'CREATE' 'INDEXTYPE' ('IF' 'NOT' 'EXISTS')? (ID '.')? ID ('SHARING' '=' ('METADATA' | 'NONE'))? 'FOR' (ID '.')? ID '(' DATA_TYPE (',' DATA_TYPE)* ')' (',' (ID '.')? ID '(' DATA_TYPE (',' DATA_TYPE)* ')')* createIndextype_using_type_clause ('WITH' 'LOCAL' 'RANGE'? 'PARTITION')? createIndextype_storage_table_clause? createIndextype_array_dml_clause?
    ;
createIndextype_using_type_clause
    : 'USING' (ID '.')? ID
    ;
createIndextype_storage_table_clause
    : 'WITH' ('SYSTEM' | 'USER') 'MANAGED' 'STORAGE' 'TABLES'
    ;
createIndextype_array_dml_clause
    : 'WITH' 'ARRAY' 'DML' ('(' (ID '.')? ID ',' (ID '.')? ID ')' (',' '(' (ID '.')? ID ',' (ID '.')? ID ')')?)*
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-INMEMORY-JOIN-GROUP.html
createInmemoryJoinGroup
    : 'CREATE' 'INMEMORY' 'JOIN' 'GROUP' ('IF' 'NOT' 'EXISTS')? (ID '.')? ID '(' (ID '.')? ID '(' ID ')' (',' (ID '.')? ID '(' ID ')')* ')'
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-JAVA.html
createJava
    : 'CREATE' ('OR' 'REPLACE' | 'IF' 'NOT' 'EXISTS')? ('AND' ('RESOLVE' | 'COMPILE'))? 'NOFORCE'? 'JAVA' ('SOURCE' | 'CLASS' | 'RESOURCE') ('NAMED' (ID '.')? ID)? ('SHARING' '=' ('METADATA' | 'NONE'))? createJava_invoker_rights_clause? createJava_resolver_clause? ('USING' ('BFILE' '(' ID ',' STRING ')' | ('CLOB' | 'BLOB' | 'BFILE') QUERY | ID) | 'AS' STRING)
    ;
createJava_invoker_rights_clause
    : 'AUTHID' ('CURRENT_USER' | 'DEFINER')
    ;
createJava_resolver_clause
    : 'RESOLVER' '(' '(' STRING (ID | '-') ')' ('(' STRING (ID | '-') ')')* ')'
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-LIBRARY.html
createLibrary
    : 'CREATE' ('OR' 'REPLACE' | 'IF' 'NOT' 'EXISTS')? ('EDITIONABLE' | 'NONEDITIONABLE')? 'LIBRARY' (ID '.')? ID createLibrary_plsql_library_source
    ;
createLibrary_plsql_library_source
    : ('IS' | 'AS') STRING ('AGENT' ID)? ('CREDENTIAL' ID)? ('SHARING' '=' ('METADATA' | 'NONE'))?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-LOCKDOWN-PROFILE.html
createLockdownProfile
    : 'CREATE' 'LOCKDOWN' 'PROFILE' ID (createLockdownProfile_static_base_profile | createLockdownProfile_dynamic_base_profile)?
    ;
createLockdownProfile_static_base_profile
    : 'USING' ID
    ;
createLockdownProfile_dynamic_base_profile
    : 'INCLUDING' ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-MATERIALIZED-VIEW-LOG.html
// fix: left out where optional, since no quoted page defines it or what it needs: lob_parameters, storage_clause
createMaterializedViewLog
    : 'CREATE' 'MATERIALIZED' 'VIEW' 'LOG' ('IF' 'NOT' 'EXISTS')? 'ON' (ID '.')? ID ('SHARING' '=' ('METADATA' | 'NONE'))? createMaterializedViewLog_physical_attributes_clause? ('TABLESPACE' ID)? createMaterializedViewLog_logging_clause? (('CACHE' | 'NOCACHE'))? createMaterializedViewLog_parallel_clause? createMaterializedViewLog_table_partitioning_clauses? ('WITH' (('PRIMARY' 'KEY' | 'ROWID' | 'OBJECT' 'ID'))? 'SEQUENCE'? ('(' ID (',' ID)* ')')? ('COMMIT' 'SCN')?)? createMaterializedViewLog_new_values_clause? createMaterializedViewLog_mv_log_purge_clause? createMaterializedViewLog_for_refresh_clause?
    ;
createMaterializedViewLog_physical_attributes_clause
    : ('PCTFREE' NUMBER)? ('PCTUSED' NUMBER)? ('INITRANS' NUMBER)? ('MAXTRANS' NUMBER)?
    ;
createMaterializedViewLog_logging_clause
    : ('LOGGING' | 'NOLOGGING')
    ;
createMaterializedViewLog_parallel_clause
    : ('NOPARALLEL' | 'PARALLEL' NUMBER?)
    ;
createMaterializedViewLog_table_partitioning_clauses
    : (createMaterializedViewLog_range_partitions | createMaterializedViewLog_list_partitions | createMaterializedViewLog_hash_partitions | createMaterializedViewLog_composite_range_partitions | createMaterializedViewLog_composite_list_partitions | createMaterializedViewLog_composite_hash_partitions | createMaterializedViewLog_reference_partitioning | createMaterializedViewLog_system_partitioning | createMaterializedViewLog_consistent_hash_partitions | createMaterializedViewLog_consistent_hash_with_subpartitions | createMaterializedViewLog_partitionset_clauses)
    ;
createMaterializedViewLog_new_values_clause
    : ('INCLUDING' | 'EXCLUDING') 'NEW' 'VALUES'
    ;
createMaterializedViewLog_mv_log_purge_clause
    : 'PURGE' ('IMMEDIATE' ('SYNCHRONOUS' | 'ASYNCHRONOUS') | 'START' 'WITH' EXPR ('NEXT' EXPR)? ('REPEAT' 'INTERVAL' EXPR)?)
    ;
createMaterializedViewLog_for_refresh_clause
    : ('FOR' 'FAST' 'REFRESH' | 'FOR' 'SYNCHRONOUS' 'REFRESH' 'USING' ID)
    ;
createMaterializedViewLog_range_partitions
    : 'PARTITION' 'BY' 'RANGE' '(' ID (',' ID)* ')' ('INTERVAL' '(' EXPR ')' ('STORE' 'IN' '(' ID (',' ID)* ')')?)? '(' 'PARTITION' ID? createMaterializedViewLog_range_values_clause createMaterializedViewLog_table_partition_description? (',' 'PARTITION' ID? createMaterializedViewLog_range_values_clause createMaterializedViewLog_table_partition_description?)* ')'
    ;
createMaterializedViewLog_list_partitions
    : 'PARTITION' 'BY' 'LIST' '(' ID (',' ID)* ')' ('AUTOMATIC' ('STORE' 'IN' '(' ID (',' ID)* ')')?)? '(' 'PARTITION' ID? createMaterializedViewLog_list_values_clause createMaterializedViewLog_table_partition_description? (',' 'PARTITION' ID? createMaterializedViewLog_list_values_clause createMaterializedViewLog_table_partition_description?)* ')'
    ;
createMaterializedViewLog_hash_partitions
    : 'PARTITION' 'BY' 'HASH' '(' ID (',' ID)* ')' (createMaterializedViewLog_individual_hash_partitions | createMaterializedViewLog_hash_partitions_by_quantity)
    ;
createMaterializedViewLog_composite_range_partitions
    : 'PARTITION' 'BY' 'RANGE' '(' ID (',' ID)* ')' ('INTERVAL' '(' EXPR ')' ('STORE' 'IN' '(' ID (',' ID)* ')')?)? (createMaterializedViewLog_subpartition_by_range | createMaterializedViewLog_subpartition_by_list | createMaterializedViewLog_subpartition_by_hash) '(' 'PARTITION' ID? createMaterializedViewLog_range_values_clause (createMaterializedViewLog_table_partition_description ((createMaterializedViewLog_range_subpartition_desc | createMaterializedViewLog_list_subpartition_desc | createMaterializedViewLog_individual_hash_subparts | createMaterializedViewLog_hash_subparts_by_quantity))? (',' createMaterializedViewLog_table_partition_description ((createMaterializedViewLog_range_subpartition_desc | createMaterializedViewLog_list_subpartition_desc | createMaterializedViewLog_individual_hash_subparts | createMaterializedViewLog_hash_subparts_by_quantity))?)*)* ')'
    ;
createMaterializedViewLog_composite_list_partitions
    : 'PARTITION' 'BY' 'LIST' '(' ID (',' ID)* ')' ('AUTOMATIC' ('STORE' 'IN' '(' ID (',' ID)* ')')?)? (createMaterializedViewLog_subpartition_by_range | createMaterializedViewLog_subpartition_by_list | createMaterializedViewLog_subpartition_by_hash) '(' 'PARTITION' ID? createMaterializedViewLog_list_values_clause (createMaterializedViewLog_table_partition_description ((createMaterializedViewLog_range_subpartition_desc | createMaterializedViewLog_list_subpartition_desc | createMaterializedViewLog_individual_hash_subparts | createMaterializedViewLog_hash_subparts_by_quantity))? (',' createMaterializedViewLog_table_partition_description ((createMaterializedViewLog_range_subpartition_desc | createMaterializedViewLog_list_subpartition_desc | createMaterializedViewLog_individual_hash_subparts | createMaterializedViewLog_hash_subparts_by_quantity))?)*)* ')'
    ;
createMaterializedViewLog_composite_hash_partitions
    : 'PARTITION' 'BY' 'HASH' '(' ID (',' ID)* ')' (createMaterializedViewLog_subpartition_by_range | createMaterializedViewLog_subpartition_by_list | createMaterializedViewLog_subpartition_by_hash) (createMaterializedViewLog_individual_hash_partitions | createMaterializedViewLog_hash_partitions_by_quantity)
    ;
createMaterializedViewLog_reference_partitioning
    : 'PARTITION' 'BY' 'REFERENCE' '(' ID ')' ('(' 'PARTITION' ID? createMaterializedViewLog_table_partition_description? (',' 'PARTITION' ID? createMaterializedViewLog_table_partition_description?)* ')')?
    ;
createMaterializedViewLog_system_partitioning
    : ('PARTITION' 'BY' 'SYSTEM')? ('PARTITIONS' NUMBER | '(' 'PARTITION' ID? createMaterializedViewLog_table_partition_description? (',' 'PARTITION' ID? createMaterializedViewLog_table_partition_description?)* ')')?
    ;
createMaterializedViewLog_consistent_hash_partitions
    : 'PARTITION' 'BY' 'CONSISTENT' 'HASH' '(' ID (',' ID)* ')' ('PARTITIONS' 'AUTO')? 'TABLESPACE' 'SET' ID
    ;
createMaterializedViewLog_consistent_hash_with_subpartitions
    : 'PARTITION' 'BY' 'CONSISTENT' 'HASH' '(' ID (',' ID)* ')' (createMaterializedViewLog_subpartition_by_range | createMaterializedViewLog_subpartition_by_list | createMaterializedViewLog_subpartition_by_hash) ('PARTITIONS' 'AUTO')? 'TABLESPACE' 'SET' ID
    ;
createMaterializedViewLog_partitionset_clauses
    : (createMaterializedViewLog_range_partitionset_clause | createMaterializedViewLog_list_partitionset_clause)
    ;
createMaterializedViewLog_range_values_clause
    : 'VALUES' 'LESS' 'THAN' '(' (EXPR | 'MAXVALUE') (',' (EXPR | 'MAXVALUE'))* ')'
    ;
createMaterializedViewLog_table_partition_description
    : createMaterializedViewLog_deferred_segment_creation? createMaterializedViewLog_read_only_clause? createMaterializedViewLog_indexing_clause? createMaterializedViewLog_segment_attributes_clause? (createMaterializedViewLog_table_compression | createMaterializedViewLog_prefix_compression)? createMaterializedViewLog_inmemory_clause? createMaterializedViewLog_ilm_clause? ('OVERFLOW' createMaterializedViewLog_segment_attributes_clause?)?
    ;
createMaterializedViewLog_list_values_clause
    : 'VALUES' '(' (EXPR (',' EXPR)* | 'DEFAULT') ')'
    ;
createMaterializedViewLog_individual_hash_partitions
    : '(' 'PARTITION' ID? createMaterializedViewLog_partitioning_storage_clause? (',' 'PARTITION' ID? createMaterializedViewLog_partitioning_storage_clause?)* ')'
    ;
createMaterializedViewLog_hash_partitions_by_quantity
    : 'PARTITIONS' NUMBER ('STORE' 'IN' '(' ID (',' ID)* ')')?
    ;
createMaterializedViewLog_subpartition_by_range
    : 'SUBPARTITION' 'BY' 'RANGE' '(' ID (',' ID)* ')' createMaterializedViewLog_subpartition_template?
    ;
createMaterializedViewLog_subpartition_by_list
    : 'SUBPARTITION' 'BY' 'LIST' '(' ID (',' ID)* ')' createMaterializedViewLog_subpartition_template?
    ;
createMaterializedViewLog_subpartition_by_hash
    : 'SUBPARTITION' 'BY' 'HASH' '(' ID (',' ID)* ')' ('SUBPARTITIONS' NUMBER ('STORE' 'IN' '(' ID (',' ID)* ')')? | createMaterializedViewLog_subpartition_template)?
    ;
createMaterializedViewLog_range_subpartition_desc
    : 'SUBPARTITION' ID? createMaterializedViewLog_range_values_clause createMaterializedViewLog_read_only_clause? createMaterializedViewLog_indexing_clause? createMaterializedViewLog_partitioning_storage_clause?
    ;
createMaterializedViewLog_list_subpartition_desc
    : 'SUBPARTITION' ID? createMaterializedViewLog_list_values_clause createMaterializedViewLog_read_only_clause? createMaterializedViewLog_indexing_clause? createMaterializedViewLog_partitioning_storage_clause?
    ;
createMaterializedViewLog_individual_hash_subparts
    : 'SUBPARTITION' ID? createMaterializedViewLog_read_only_clause? createMaterializedViewLog_indexing_clause? createMaterializedViewLog_partitioning_storage_clause?
    ;
createMaterializedViewLog_hash_subparts_by_quantity
    : 'SUBPARTITIONS' NUMBER ('STORE' 'IN' '(' ID (',' ID)* ')')?
    ;
createMaterializedViewLog_range_partitionset_clause
    : 'PARTITIONSET' 'BY' 'RANGE' '(' ID (',' ID)* ')' 'PARTITION' 'BY' 'CONSISTENT' 'HASH' '(' ID (',' ID)* ')' (createMaterializedViewLog_subpartition_by_range | createMaterializedViewLog_subpartition_by_list | createMaterializedViewLog_subpartition_by_hash)? 'PARTITIONS' 'AUTO' '(' 'PARTITIONSET' ID createMaterializedViewLog_range_values_clause ('TABLESPACE' 'SET' ID)? createMaterializedViewLog_lob_storage_clause? createMaterializedViewLog_subpartition_template? (',' 'PARTITIONSET' ID createMaterializedViewLog_range_values_clause ('TABLESPACE' 'SET' ID)? createMaterializedViewLog_lob_storage_clause? createMaterializedViewLog_subpartition_template?)* ')'
    ;
createMaterializedViewLog_list_partitionset_clause
    : 'PARTITIONSET' 'BY' 'LIST' '(' ID ')' 'PARTITION' 'BY' 'CONSISTENT' 'HASH' '(' ID (',' ID)* ')' (createMaterializedViewLog_subpartition_by_range | createMaterializedViewLog_subpartition_by_list | createMaterializedViewLog_subpartition_by_hash)? 'PARTITIONS' 'AUTO' '(' 'PARTITIONSET' ID createMaterializedViewLog_list_values_clause ('TABLESPACE' 'SET' ID)? createMaterializedViewLog_lob_storage_clause? createMaterializedViewLog_subpartition_template? (',' 'PARTITIONSET' ID createMaterializedViewLog_list_values_clause ('TABLESPACE' 'SET' ID)? createMaterializedViewLog_lob_storage_clause? createMaterializedViewLog_subpartition_template?)* ')'
    ;
createMaterializedViewLog_deferred_segment_creation
    : 'SEGMENT' 'CREATION' ('IMMEDIATE' | 'DEFERRED')
    ;
createMaterializedViewLog_read_only_clause
    : ('READ' 'ONLY' | 'READ' 'WRITE')
    ;
createMaterializedViewLog_indexing_clause
    : 'INDEXING' ('ON' | 'OFF')
    ;
createMaterializedViewLog_segment_attributes_clause
    : createMaterializedViewLog_physical_attributes_clause? ('TABLESPACE' ID)? createMaterializedViewLog_logging_clause?
    ;
createMaterializedViewLog_table_compression
    : ('COMPRESS' 'BASIC'? | 'COMPRESS' 'FOR' ('OLTP' | 'QUERY' ('LOW' | 'HIGH')? | 'ARCHIVE' ('LOW' | 'HIGH')?) | 'NOCOMPRESS')
    ;
createMaterializedViewLog_prefix_compression
    : 'COMPRESS' NUMBER?
    ;
createMaterializedViewLog_inmemory_clause
    : ('INMEMORY' createMaterializedViewLog_inmemory_attributes? | 'NO' 'INMEMORY')
    ;
createMaterializedViewLog_ilm_clause
    : 'ILM' ('ADD' 'POLICY' createMaterializedViewLog_ilm_policy_clause | 'DELETE' 'POLICY' ID | 'DELETE_ALL' | 'ENABLE_ALL' | 'DISABLE_ALL')
    ;
createMaterializedViewLog_partitioning_storage_clause
    : ('TABLESPACE' ID)? ('OVERFLOW' 'TABLESPACE' ID)? createMaterializedViewLog_table_compression? createMaterializedViewLog_inmemory_clause?
    ;
createMaterializedViewLog_subpartition_template
    : 'SUBPARTITION' 'TEMPLATE' '(' (createMaterializedViewLog_range_subpartition_desc | createMaterializedViewLog_list_subpartition_desc | createMaterializedViewLog_individual_hash_subparts) (',' (createMaterializedViewLog_range_subpartition_desc | createMaterializedViewLog_list_subpartition_desc | createMaterializedViewLog_individual_hash_subparts))* ')'
    | 'SUBPARTITION' 'TEMPLATE' NUMBER
    ;
createMaterializedViewLog_lob_storage_clause
    : 'LOB' '(' ID (',' ID)* ')' 'STORE' 'AS' ID? '(' createMaterializedViewLog_lob_storage_parameters ')'
    ;
createMaterializedViewLog_inmemory_attributes
    : createMaterializedViewLog_inmemory_memcompress? createMaterializedViewLog_inmemory_priority? createMaterializedViewLog_inmemory_distribute? createMaterializedViewLog_inmemory_duplicate?
    ;
createMaterializedViewLog_ilm_policy_clause
    : (createMaterializedViewLog_ilm_compression_policy | createMaterializedViewLog_ilm_tiering_policy | createMaterializedViewLog_ilm_inmemory_policy)
    ;
createMaterializedViewLog_lob_storage_parameters
    : (('TABLESPACE' ID | createMaterializedViewLog_logging_clause)+ | ('ENABLE' | 'DISABLE') 'STORAGE' 'IN' 'ROW' | 'CHUNK' NUMBER | 'PCTVERSION' NUMBER | 'FREEPOOLS' NUMBER | createMaterializedViewLog_lob_retention_clause | createMaterializedViewLog_lob_deduplicate_clause | createMaterializedViewLog_lob_compression_clause | ('ENCRYPT' createMaterializedViewLog_encryption_spec | 'DECRYPT') | ('CACHE' | 'NOCACHE' | 'CACHE' 'READS') createMaterializedViewLog_logging_clause?)+
    ;
createMaterializedViewLog_inmemory_memcompress
    : 'MEMCOMPRESS' 'FOR' ('QUERY' ('LOW' | 'HIGH')? | 'OLTP' | 'ARCHIVE' ('LOW' | 'HIGH')?)
    ;
createMaterializedViewLog_inmemory_priority
    : 'PRIORITY' ('LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL')
    ;
createMaterializedViewLog_inmemory_distribute
    : 'DISTRIBUTE' 'BY' ('ROWID' ('ACROSS' NUMBER?)? | 'PARTITION' | 'SUBPARTITION')
    ;
createMaterializedViewLog_inmemory_duplicate
    : 'DUPLICATE' ('ALL' | 'FOR' 'SERVICE' ID (',' ID)* | 'NONE')
    ;
createMaterializedViewLog_ilm_compression_policy
    : (createMaterializedViewLog_table_compression | 'ROW' 'STORE' 'COMPRESS' 'ADVANCED'? | 'COLUMN' 'STORE' 'COMPRESS' ('FOR' 'QUERY' ('LOW' | 'HIGH')?)?) ('SEGMENT' | 'GROUP' | 'ROW') 'AFTER' createMaterializedViewLog_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION')
    ;
createMaterializedViewLog_ilm_tiering_policy
    : 'TIER' 'TO' ID ('SEGMENT' | 'GROUP' | 'ROW')? ('ON' ID)?
    ;
createMaterializedViewLog_ilm_inmemory_policy
    : ('SET' 'INMEMORY' createMaterializedViewLog_inmemory_memcompress? createMaterializedViewLog_inmemory_priority? | 'MODIFY' 'INMEMORY' createMaterializedViewLog_inmemory_memcompress | 'NO' 'INMEMORY') 'SEGMENT' 'AFTER' createMaterializedViewLog_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION')
    ;
createMaterializedViewLog_lob_retention_clause
    : 'RETENTION' ('MAX' | 'MIN' NUMBER | 'AUTO' | 'NONE')?
    ;
createMaterializedViewLog_lob_deduplicate_clause
    : ('DEDUPLICATE' | 'KEEP_DUPLICATES')
    ;
createMaterializedViewLog_lob_compression_clause
    : ('COMPRESS' ('HIGH' | 'MEDIUM' | 'LOW')? | 'NOCOMPRESS')
    ;
createMaterializedViewLog_encryption_spec
    : ('USING' STRING)? ('IDENTIFIED' 'BY' ID)? STRING? ('NO'? 'SALT')?
    ;
createMaterializedViewLog_ilm_time_period
    : NUMBER ('DAY' | 'DAYS' | 'MONTH' | 'MONTHS' | 'YEAR' | 'YEARS')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-MATERIALIZED-VIEW.html
// fix: left out where optional, since no quoted page defines it or what it needs: create_index_statement, json_parameters, lob_parameters, storage_clause
createMaterializedView
    : 'CREATE' 'MATERIALIZED' 'VIEW' ('IF' 'NOT' 'EXISTS')? (ID '.')? ID ('(' ID ('ENCRYPT' createMaterializedView_encryption_spec)? (',' ID ('ENCRYPT' createMaterializedView_encryption_spec)?)* ')')? ('OF' ID)? createMaterializedView_scoped_table_ref_constraint? ('DEFAULT' 'COLLATION' ID)? ('ON' 'PREBUILT' 'TABLE' (('WITH' 'REDUCED' 'PRECISION' | 'WITHOUT' 'REDUCED' 'PRECISION'))?)? createMaterializedView_physical_properties? createMaterializedView_materialized_view_props? createMaterializedView_create_mv_refresh? createMaterializedView_evaluation_edition_clause? (('ENABLE' | 'DISABLE') 'ON' 'QUERY' 'COMPUTATION')? createMaterializedView_query_rewrite_clause? createMaterializedView_unusable_editions_clause? 'AS' QUERY createMaterializedView_annotations_clause?
    ;
createMaterializedView_encryption_spec
    : ('USING' STRING)? ('IDENTIFIED' 'BY' ID)? STRING? ('NO'? 'SALT')?
    ;
createMaterializedView_scoped_table_ref_constraint
    : 'SCOPE' 'FOR' '(' ('REF' '(' ID ')' | ID) ')' 'IS' (ID '.')? ID
    ;
createMaterializedView_physical_properties
    : createMaterializedView_deferred_segment_creation? createMaterializedView_segment_attributes_clause? createMaterializedView_table_compression? createMaterializedView_inmemory_table_clause? ((createMaterializedView_heap_org_table_clause | createMaterializedView_index_org_table_clause))? ('CLUSTER' ID '(' ID (',' ID)* ')')?
    ;
createMaterializedView_materialized_view_props
    : createMaterializedView_column_properties? createMaterializedView_table_partitioning_clauses? (('CACHE' | 'NOCACHE'))? createMaterializedView_parallel_clause? createMaterializedView_build_clause?
    ;
createMaterializedView_create_mv_refresh
    : ('REFRESH' (('FAST' | 'COMPLETE' | 'FORCE'))? (('ON' 'COMMIT' | 'ON' 'DEMAND' | 'ON' 'STATEMENT'))? ('START' 'WITH' EXPR)? ('NEXT' EXPR)? (('WITH' 'PRIMARY' 'KEY' | 'WITH' 'ROWID'))? ('USING' 'DEFAULT' ('MASTER' | 'LOCAL')? 'ROLLBACK' 'SEGMENT' ID?)? ('USING' ('TRUSTED' | 'ENFORCED') 'CONSTRAINTS')? | 'NEVER' 'REFRESH')
    ;
createMaterializedView_evaluation_edition_clause
    : 'EVALUATE' 'USING' ('CURRENT' 'EDITION' | 'EDITION' ID | 'NULL' 'EDITION')
    ;
createMaterializedView_query_rewrite_clause
    : ('ENABLE' | 'DISABLE') 'QUERY' 'REWRITE'
    ;
createMaterializedView_unusable_editions_clause
    : 'UNUSABLE' 'EDITIONS' ('BEFORE' ('CURRENT' 'EDITION' | 'EDITION' ID))? ('BEGINNING' 'WITH' ('CURRENT' 'EDITION' | 'EDITION' ID | 'NULL' 'EDITION'))?
    ;
createMaterializedView_annotations_clause
    : 'ANNOTATIONS' '(' ('ADD' | 'REPLACE') ID EXPR? (',' ('ADD' | 'REPLACE') ID EXPR?)* ')'
    ;
createMaterializedView_deferred_segment_creation
    : 'SEGMENT' 'CREATION' ('IMMEDIATE' | 'DEFERRED')
    ;
createMaterializedView_segment_attributes_clause
    : ('PCTFREE' NUMBER)? ('PCTUSED' NUMBER)? ('INITRANS' NUMBER)? ('TABLESPACE' ID)? createMaterializedView_logging_clause?
    ;
createMaterializedView_table_compression
    : ('COMPRESS' ('BASIC' | 'FOR' 'OLTP' | 'FOR' 'QUERY' (('LOW' | 'HIGH'))? | 'FOR' 'ARCHIVE' (('LOW' | 'HIGH'))?)? | 'NOCOMPRESS')
    ;
createMaterializedView_inmemory_table_clause
    : ('INMEMORY' createMaterializedView_inmemory_attributes? (createMaterializedView_inmemory_column_clause (',' createMaterializedView_inmemory_column_clause)?)* | 'NO' 'INMEMORY')
    ;
createMaterializedView_heap_org_table_clause
    : (('CACHE' | 'NOCACHE'))?
    ;
createMaterializedView_index_org_table_clause
    : 'ORGANIZATION' 'INDEX' (('COMPRESS' NUMBER? | 'NOCOMPRESS'))? ('OVERFLOW' createMaterializedView_segment_attributes_clause?)?
    ;
createMaterializedView_column_properties
    : (createMaterializedView_object_type_col_properties | createMaterializedView_nested_table_col_properties | (createMaterializedView_varray_col_properties | createMaterializedView_lob_storage_clause) ('(' createMaterializedView_lob_partition_storage (',' createMaterializedView_lob_partition_storage)* ')')? | createMaterializedView_xmltype_column_properties | createMaterializedView_json_storage_clause)+
    ;
createMaterializedView_table_partitioning_clauses
    : (createMaterializedView_range_partitions | createMaterializedView_list_partitions | createMaterializedView_hash_partitions | createMaterializedView_composite_range_partitions | createMaterializedView_composite_list_partitions | createMaterializedView_composite_hash_partitions | createMaterializedView_reference_partitioning | createMaterializedView_system_partitioning | createMaterializedView_consistent_hash_partitions | createMaterializedView_consistent_hash_with_subpartitions | createMaterializedView_partitionset_clauses)
    ;
createMaterializedView_parallel_clause
    : ('NOPARALLEL' | 'PARALLEL' NUMBER?)
    ;
createMaterializedView_build_clause
    : 'BUILD' ('IMMEDIATE' | 'DEFERRED')
    ;
createMaterializedView_logging_clause
    : ('LOGGING' | 'NOLOGGING')
    ;
createMaterializedView_inmemory_attributes
    : createMaterializedView_inmemory_memcompress? createMaterializedView_inmemory_priority? createMaterializedView_inmemory_distribute? createMaterializedView_inmemory_duplicate?
    ;
createMaterializedView_inmemory_column_clause
    : '(' ID (',' ID)* ')' ('INMEMORY' createMaterializedView_inmemory_memcompress? | 'NO' 'INMEMORY')
    ;
createMaterializedView_object_type_col_properties
    : 'COLUMN' ID createMaterializedView_substitutable_column_clause
    ;
createMaterializedView_nested_table_col_properties
    : 'NESTED' 'TABLE' (ID | 'COLUMN_VALUE') createMaterializedView_substitutable_column_clause? ('LOCAL' | 'GLOBAL')? 'STORE' 'AS' ID ('(' '(' createMaterializedView_object_properties ')' createMaterializedView_physical_properties? createMaterializedView_column_properties? ')')? ('RETURN' 'AS'? ('LOCATOR' | 'VALUE'))?
    ;
createMaterializedView_varray_col_properties
    : 'VARRAY' ID (createMaterializedView_substitutable_column_clause? 'STORE' 'AS' ('SECUREFILE' | 'BASICFILE')? 'LOB' (ID? '(' createMaterializedView_lob_storage_parameters ')' | ID) | createMaterializedView_substitutable_column_clause)
    ;
createMaterializedView_lob_storage_clause
    : 'LOB' '(' ID (',' ID)* ')' 'STORE' 'AS' ID? '(' createMaterializedView_lob_storage_parameters ')'
    ;
createMaterializedView_lob_partition_storage
    : 'PARTITION' ID (createMaterializedView_lob_storage_clause | createMaterializedView_varray_col_properties)+ ('(' 'SUBPARTITION' ID (createMaterializedView_lob_storage_clause | createMaterializedView_varray_col_properties)+ (')')+)?
    ;
createMaterializedView_xmltype_column_properties
    : 'XMLTYPE' 'COLUMN'? ID createMaterializedView_xmltype_storage? createMaterializedView_xmlschema_spec?
    ;
createMaterializedView_json_storage_clause
    : 'JSON' '(' ID ')' 'STORE' 'AS'
    ;
createMaterializedView_range_partitions
    : 'PARTITION' 'BY' 'RANGE' '(' ID (',' ID)* ')' ('INTERVAL' '(' EXPR ')' ('STORE' 'IN' '(' ID (',' ID)* ')')?)? '(' 'PARTITION' ID? createMaterializedView_range_values_clause createMaterializedView_table_partition_description? (',' 'PARTITION' ID? createMaterializedView_range_values_clause createMaterializedView_table_partition_description?)* ')'
    ;
createMaterializedView_list_partitions
    : 'PARTITION' 'BY' 'LIST' '(' ID (',' ID)* ')' ('AUTOMATIC' ('STORE' 'IN' '(' ID (',' ID)* ')')?)? '(' 'PARTITION' ID? createMaterializedView_list_values_clause createMaterializedView_table_partition_description? (',' 'PARTITION' ID? createMaterializedView_list_values_clause createMaterializedView_table_partition_description?)* ')'
    ;
createMaterializedView_hash_partitions
    : 'PARTITION' 'BY' 'HASH' '(' ID (',' ID)* ')' (createMaterializedView_individual_hash_partitions | createMaterializedView_hash_partitions_by_quantity)
    ;
createMaterializedView_composite_range_partitions
    : 'PARTITION' 'BY' 'RANGE' '(' ID (',' ID)* ')' ('INTERVAL' '(' EXPR ')' ('STORE' 'IN' '(' ID (',' ID)* ')')?)? (createMaterializedView_subpartition_by_range | createMaterializedView_subpartition_by_list | createMaterializedView_subpartition_by_hash) '(' 'PARTITION' ID? createMaterializedView_range_values_clause (createMaterializedView_table_partition_description ((createMaterializedView_range_subpartition_desc | createMaterializedView_list_subpartition_desc | createMaterializedView_individual_hash_subparts | createMaterializedView_hash_subparts_by_quantity))? (',' createMaterializedView_table_partition_description ((createMaterializedView_range_subpartition_desc | createMaterializedView_list_subpartition_desc | createMaterializedView_individual_hash_subparts | createMaterializedView_hash_subparts_by_quantity))?)*)* ')'
    ;
createMaterializedView_composite_list_partitions
    : 'PARTITION' 'BY' 'LIST' '(' ID (',' ID)* ')' ('AUTOMATIC' ('STORE' 'IN' '(' ID (',' ID)* ')')?)? (createMaterializedView_subpartition_by_range | createMaterializedView_subpartition_by_list | createMaterializedView_subpartition_by_hash) '(' 'PARTITION' ID? createMaterializedView_list_values_clause (createMaterializedView_table_partition_description ((createMaterializedView_range_subpartition_desc | createMaterializedView_list_subpartition_desc | createMaterializedView_individual_hash_subparts | createMaterializedView_hash_subparts_by_quantity))? (',' createMaterializedView_table_partition_description ((createMaterializedView_range_subpartition_desc | createMaterializedView_list_subpartition_desc | createMaterializedView_individual_hash_subparts | createMaterializedView_hash_subparts_by_quantity))?)*)* ')'
    ;
createMaterializedView_composite_hash_partitions
    : 'PARTITION' 'BY' 'HASH' '(' ID (',' ID)* ')' (createMaterializedView_subpartition_by_range | createMaterializedView_subpartition_by_list | createMaterializedView_subpartition_by_hash) (createMaterializedView_individual_hash_partitions | createMaterializedView_hash_partitions_by_quantity)
    ;
createMaterializedView_reference_partitioning
    : 'PARTITION' 'BY' 'REFERENCE' '(' ID ')' ('(' 'PARTITION' ID? createMaterializedView_table_partition_description? (',' 'PARTITION' ID? createMaterializedView_table_partition_description?)* ')')?
    ;
createMaterializedView_system_partitioning
    : ('PARTITION' 'BY' 'SYSTEM')? ('PARTITIONS' NUMBER | '(' 'PARTITION' ID? createMaterializedView_table_partition_description? (',' 'PARTITION' ID? createMaterializedView_table_partition_description?)* ')')?
    ;
createMaterializedView_consistent_hash_partitions
    : 'PARTITION' 'BY' 'CONSISTENT' 'HASH' '(' ID (',' ID)* ')' ('PARTITIONS' 'AUTO')? 'TABLESPACE' 'SET' ID
    ;
createMaterializedView_consistent_hash_with_subpartitions
    : 'PARTITION' 'BY' 'CONSISTENT' 'HASH' '(' ID (',' ID)* ')' (createMaterializedView_subpartition_by_range | createMaterializedView_subpartition_by_list | createMaterializedView_subpartition_by_hash) ('PARTITIONS' 'AUTO')? 'TABLESPACE' 'SET' ID
    ;
createMaterializedView_partitionset_clauses
    : (createMaterializedView_range_partitionset_clause | createMaterializedView_list_partitionset_clause)
    ;
createMaterializedView_inmemory_memcompress
    : 'MEMCOMPRESS' 'FOR' ('DML' | 'QUERY' (('LOW' | 'HIGH'))? | 'CAPACITY' (('LOW' | 'HIGH'))?)
    | 'NO' 'MEMCOMPRESS'
    ;
createMaterializedView_inmemory_priority
    : 'PRIORITY' ('NONE' | 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL')
    ;
createMaterializedView_inmemory_distribute
    : 'DISTRIBUTE' (('AUTO' | 'BY' 'ROWID' 'RANGE' | 'BY' 'PARTITION' | 'BY' 'SUBPARTITION'))?
    ;
createMaterializedView_inmemory_duplicate
    : 'DUPLICATE' 'ALL'?
    | 'NO' 'DUPLICATE'
    ;
createMaterializedView_substitutable_column_clause
    : ('ELEMENT' 'IS' 'OF' 'TYPE'? '(' 'ONLY'? ID ')' | 'NOT'? 'SUBSTITUTABLE' 'AT' 'ALL' 'LEVELS')
    ;
createMaterializedView_object_properties
    : (ID ('DEFAULT' EXPR)? ((createMaterializedView_inline_constraint | createMaterializedView_inline_ref_constraint)+)? | createMaterializedView_out_of_line_constraint | createMaterializedView_out_of_line_ref_constraint | createMaterializedView_supplemental_logging_props) (',' (ID ('DEFAULT' EXPR)? ((createMaterializedView_inline_constraint | createMaterializedView_inline_ref_constraint)+)? | createMaterializedView_out_of_line_constraint | createMaterializedView_out_of_line_ref_constraint | createMaterializedView_supplemental_logging_props))*
    ;
createMaterializedView_lob_storage_parameters
    : (('TABLESPACE' ID | createMaterializedView_logging_clause)+ | ('ENABLE' | 'DISABLE') 'STORAGE' 'IN' 'ROW' | 'CHUNK' NUMBER | 'PCTVERSION' NUMBER | 'FREEPOOLS' NUMBER | createMaterializedView_lob_retention_clause | createMaterializedView_lob_deduplicate_clause | createMaterializedView_lob_compression_clause | ('ENCRYPT' createMaterializedView_encryption_spec | 'DECRYPT') | ('CACHE' | 'NOCACHE' | 'CACHE' 'READS') createMaterializedView_logging_clause?)+
    ;
createMaterializedView_xmltype_storage
    : 'STORE' 'AS' ('OBJECT' 'RELATIONAL' | ('SECUREFILE' | 'BASICFILE')? ('CLOB' | 'BINARY' 'XML') (ID ('(' createMaterializedView_lob_storage_parameters ')')? | '(' createMaterializedView_lob_storage_parameters ')')?)
    ;
createMaterializedView_xmlschema_spec
    : ('XMLSCHEMA' STRING)? 'ELEMENT' ID ('STORE' 'ALL' 'VARRAYS' 'AS' ('LOBS' | 'TABLES'))? (('ALLOW' | 'DISALLOW') 'NONSCHEMA')? (('ALLOW' | 'DISALLOW') 'ANYSCHEMA')?
    ;
createMaterializedView_range_values_clause
    : 'VALUES' 'LESS' 'THAN' '(' (EXPR | 'MAXVALUE') (',' (EXPR | 'MAXVALUE'))* ')'
    ;
createMaterializedView_table_partition_description
    : createMaterializedView_deferred_segment_creation? createMaterializedView_read_only_clause? createMaterializedView_indexing_clause? createMaterializedView_segment_attributes_clause? (createMaterializedView_table_compression | createMaterializedView_prefix_compression)? createMaterializedView_inmemory_clause? createMaterializedView_ilm_clause? ('OVERFLOW' createMaterializedView_segment_attributes_clause?)?
    ;
createMaterializedView_list_values_clause
    : 'VALUES' '(' (EXPR (',' EXPR)* | 'DEFAULT') ')'
    ;
createMaterializedView_individual_hash_partitions
    : '(' 'PARTITION' ID? createMaterializedView_partitioning_storage_clause? (',' 'PARTITION' ID? createMaterializedView_partitioning_storage_clause?)* ')'
    ;
createMaterializedView_hash_partitions_by_quantity
    : 'PARTITIONS' NUMBER ('STORE' 'IN' '(' ID (',' ID)* ')')?
    ;
createMaterializedView_subpartition_by_range
    : 'SUBPARTITION' 'BY' 'RANGE' '(' ID (',' ID)* ')' createMaterializedView_subpartition_template?
    ;
createMaterializedView_subpartition_by_list
    : 'SUBPARTITION' 'BY' 'LIST' '(' ID (',' ID)* ')' createMaterializedView_subpartition_template?
    ;
createMaterializedView_subpartition_by_hash
    : 'SUBPARTITION' 'BY' 'HASH' '(' ID (',' ID)* ')' ('SUBPARTITIONS' NUMBER ('STORE' 'IN' '(' ID (',' ID)* ')')? | createMaterializedView_subpartition_template)?
    ;
createMaterializedView_range_subpartition_desc
    : 'SUBPARTITION' ID? createMaterializedView_range_values_clause createMaterializedView_read_only_clause? createMaterializedView_indexing_clause? createMaterializedView_partitioning_storage_clause?
    ;
createMaterializedView_list_subpartition_desc
    : 'SUBPARTITION' ID? createMaterializedView_list_values_clause createMaterializedView_read_only_clause? createMaterializedView_indexing_clause? createMaterializedView_partitioning_storage_clause?
    ;
createMaterializedView_individual_hash_subparts
    : 'SUBPARTITION' ID? createMaterializedView_read_only_clause? createMaterializedView_indexing_clause? createMaterializedView_partitioning_storage_clause?
    ;
createMaterializedView_hash_subparts_by_quantity
    : 'SUBPARTITIONS' NUMBER ('STORE' 'IN' '(' ID (',' ID)* ')')?
    ;
createMaterializedView_range_partitionset_clause
    : 'PARTITIONSET' 'BY' 'RANGE' '(' ID (',' ID)* ')' 'PARTITION' 'BY' 'CONSISTENT' 'HASH' '(' ID (',' ID)* ')' (createMaterializedView_subpartition_by_range | createMaterializedView_subpartition_by_list | createMaterializedView_subpartition_by_hash)? 'PARTITIONS' 'AUTO' '(' 'PARTITIONSET' ID createMaterializedView_range_values_clause ('TABLESPACE' 'SET' ID)? createMaterializedView_lob_storage_clause? createMaterializedView_subpartition_template? (',' 'PARTITIONSET' ID createMaterializedView_range_values_clause ('TABLESPACE' 'SET' ID)? createMaterializedView_lob_storage_clause? createMaterializedView_subpartition_template?)* ')'
    ;
createMaterializedView_list_partitionset_clause
    : 'PARTITIONSET' 'BY' 'LIST' '(' ID ')' 'PARTITION' 'BY' 'CONSISTENT' 'HASH' '(' ID (',' ID)* ')' (createMaterializedView_subpartition_by_range | createMaterializedView_subpartition_by_list | createMaterializedView_subpartition_by_hash)? 'PARTITIONS' 'AUTO' '(' 'PARTITIONSET' ID createMaterializedView_list_values_clause ('TABLESPACE' 'SET' ID)? createMaterializedView_lob_storage_clause? createMaterializedView_subpartition_template? (',' 'PARTITIONSET' ID createMaterializedView_list_values_clause ('TABLESPACE' 'SET' ID)? createMaterializedView_lob_storage_clause? createMaterializedView_subpartition_template?)* ')'
    ;
createMaterializedView_inline_constraint
    : ('CONSTRAINT' ID)? ('NOT'? 'NULL' | 'UNIQUE' | 'PRIMARY' 'KEY' | 'REFERENCES' (ID '.')? ID ('(' ID ')')? ('ON' 'DELETE' ('CASCADE' | 'SET' 'NULL'))? | 'CHECK' '(' CONDITION ')') createMaterializedView_constraint_state?
    ;
createMaterializedView_inline_ref_constraint
    : ('SCOPE' 'IS' (ID '.')? ID | 'WITH' 'ROWID' | ('CONSTRAINT' ID)? 'REFERENCES' (ID '.')? ID ('(' ID (',' ID)* ')')? ('ON' 'DELETE' ('CASCADE' | 'SET' 'NULL'))? createMaterializedView_constraint_state?)
    ;
createMaterializedView_out_of_line_constraint
    : ('CONSTRAINT' ID)? ('UNIQUE' '(' ID (',' ID)* ')' | 'PRIMARY' 'KEY' '(' ID (',' ID)* ')' | 'FOREIGN' 'KEY' '(' ID (',' ID)* ')' 'REFERENCES' (ID '.')? ID ('(' ID (',' ID)* ')')? ('ON' 'DELETE' ('CASCADE' | 'SET' 'NULL'))? | 'CHECK' '(' CONDITION ')') createMaterializedView_constraint_state?
    ;
createMaterializedView_out_of_line_ref_constraint
    : ('SCOPE' 'FOR' '(' ID ')' 'IS' (ID '.')? ID | 'REF' '(' ID ')' 'WITH' 'ROWID' | ('CONSTRAINT' ID)? 'FOREIGN' 'KEY' '(' ID ')' 'REFERENCES' (ID '.')? ID ('ON' 'DELETE' ('CASCADE' | 'SET' 'NULL'))? createMaterializedView_constraint_state?)
    ;
createMaterializedView_supplemental_logging_props
    : 'SUPPLEMENTAL' 'LOG' (createMaterializedView_supplemental_log_grp_clause | createMaterializedView_supplemental_id_key_clause)
    ;
createMaterializedView_lob_retention_clause
    : 'RETENTION' ('MAX' | 'MIN' NUMBER | 'AUTO' | 'NONE')?
    ;
createMaterializedView_lob_deduplicate_clause
    : ('DEDUPLICATE' | 'KEEP_DUPLICATES')
    ;
createMaterializedView_lob_compression_clause
    : ('COMPRESS' ('HIGH' | 'MEDIUM' | 'LOW')? | 'NOCOMPRESS')
    ;
createMaterializedView_read_only_clause
    : ('READ' 'ONLY' | 'READ' 'WRITE')
    ;
createMaterializedView_indexing_clause
    : 'INDEXING' ('ON' | 'OFF')
    ;
createMaterializedView_prefix_compression
    : 'COMPRESS' NUMBER?
    ;
createMaterializedView_inmemory_clause
    : ('INMEMORY' createMaterializedView_inmemory_attributes? | 'NO' 'INMEMORY')
    ;
createMaterializedView_ilm_clause
    : 'ILM' ('ADD' 'POLICY' createMaterializedView_ilm_policy_clause | 'DELETE' 'POLICY' ID | 'DELETE_ALL' | 'ENABLE_ALL' | 'DISABLE_ALL')
    ;
createMaterializedView_partitioning_storage_clause
    : ('TABLESPACE' ID)? ('OVERFLOW' 'TABLESPACE' ID)? createMaterializedView_table_compression? createMaterializedView_inmemory_clause?
    ;
createMaterializedView_subpartition_template
    : 'SUBPARTITION' 'TEMPLATE' '(' (createMaterializedView_range_subpartition_desc | createMaterializedView_list_subpartition_desc | createMaterializedView_individual_hash_subparts) (',' (createMaterializedView_range_subpartition_desc | createMaterializedView_list_subpartition_desc | createMaterializedView_individual_hash_subparts))* ')'
    | 'SUBPARTITION' 'TEMPLATE' NUMBER
    ;
createMaterializedView_constraint_state
    : ('NOT'? 'DEFERRABLE' ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? | 'INITIALLY' ('DEFERRED' | 'IMMEDIATE') ('NOT'? 'DEFERRABLE')?)? ('RELY' | 'NORELY')? createMaterializedView_using_index_clause? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? createMaterializedView_exceptions_clause?
    ;
createMaterializedView_supplemental_log_grp_clause
    : 'GROUP' ID '(' ID ('NO' 'LOG')? (',' ID ('NO' 'LOG')?)* ')' 'ALWAYS'?
    ;
createMaterializedView_supplemental_id_key_clause
    : 'DATA' '(' ('PRIMARY' 'KEY' | 'UNIQUE' | 'FOREIGN' 'KEY' | 'ALL') 'COLUMNS' ')'
    ;
createMaterializedView_ilm_policy_clause
    : (createMaterializedView_ilm_compression_policy | createMaterializedView_ilm_tiering_policy | createMaterializedView_ilm_inmemory_policy)
    ;
createMaterializedView_using_index_clause
    : 'USING' 'INDEX' ((ID '.')? ID | createMaterializedView_index_attributes)
    ;
createMaterializedView_exceptions_clause
    : 'EXCEPTIONS' 'INTO' (ID '.')? ID
    ;
createMaterializedView_ilm_compression_policy
    : (createMaterializedView_table_compression | 'ROW' 'STORE' 'COMPRESS' 'ADVANCED'? | 'COLUMN' 'STORE' 'COMPRESS' ('FOR' 'QUERY' ('LOW' | 'HIGH')?)?) ('SEGMENT' | 'GROUP' | 'ROW') 'AFTER' createMaterializedView_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION')
    ;
createMaterializedView_ilm_tiering_policy
    : 'TIER' 'TO' ID ('SEGMENT' | 'GROUP' | 'ROW')? ('ON' ID)?
    ;
createMaterializedView_ilm_inmemory_policy
    : ('SET' 'INMEMORY' createMaterializedView_inmemory_memcompress? createMaterializedView_inmemory_priority? | 'MODIFY' 'INMEMORY' createMaterializedView_inmemory_memcompress | 'NO' 'INMEMORY') 'SEGMENT' 'AFTER' createMaterializedView_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION')
    ;
createMaterializedView_index_attributes
    : (createMaterializedView_physical_attributes_clause | 'TABLESPACE' (ID | 'DEFAULT') | createMaterializedView_logging_clause | createMaterializedView_index_compression | createMaterializedView_partial_index_clause | createMaterializedView_parallel_clause)+
    ;
createMaterializedView_ilm_time_period
    : NUMBER ('DAY' | 'DAYS' | 'MONTH' | 'MONTHS' | 'YEAR' | 'YEARS')
    ;
createMaterializedView_physical_attributes_clause
    : ('PCTUSED' NUMBER)? ('PCTFREE' NUMBER)? ('INITRANS' NUMBER)?
    ;
createMaterializedView_index_compression
    : (createMaterializedView_prefix_compression | createMaterializedView_advanced_index_compression | 'NOCOMPRESS')
    ;
createMaterializedView_partial_index_clause
    : 'INDEXING' ('FULL' | 'PARTIAL')
    ;
createMaterializedView_advanced_index_compression
    : 'COMPRESS' 'ADVANCED' ('LOW' | 'HIGH')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-MATERIALIZED-ZONEMAP.html
createMaterializedZonemap
    : 'CREATE' 'MATERIALIZED' 'ZONEMAP' ('IF' 'NOT' 'EXISTS')? (ID '.')? ID createMaterializedZonemap_zonemap_attributes? createMaterializedZonemap_zonemap_refresh_clause? (('ENABLE' | 'DISABLE') 'PRUNING')? (createMaterializedZonemap_create_zonemap_on_table | createMaterializedZonemap_create_zonemap_as_subquery)
    ;
createMaterializedZonemap_zonemap_attributes
    : ('TABLESPACE' ID)? ('SCALE' NUMBER)? ('PCTFREE' NUMBER)? ('PCTUSED' NUMBER)? (('CACHE' | 'NOCACHE'))?
    ;
createMaterializedZonemap_zonemap_refresh_clause
    : 'REFRESH' (('FAST' | 'COMPLETE' | 'FORCE'))? (('ON' 'DEMAND' | 'ON' 'COMMIT' | 'ON' 'LOAD' | 'ON' 'DATA' 'MOVEMENT' | 'ON' 'LOAD' 'DATA' 'MOVEMENT'))?
    ;
createMaterializedZonemap_create_zonemap_on_table
    : 'ON' (ID '.')? ID '(' ID (',' ID)* ')'
    ;
createMaterializedZonemap_create_zonemap_as_subquery
    : ('(' ID (',' ID)* ')')? 'AS' QUERY
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-OPERATOR.html
createOperator
    : 'CREATE' ('OR' 'REPLACE' | 'IF' 'NOT' 'EXISTS')? 'OPERATOR' (ID '.')? ID createOperator_binding_clause ('SHARING' '=' ('METADATA' | 'NONE'))?
    ;
createOperator_binding_clause
    : 'BINDING' '(' (DATA_TYPE (',' DATA_TYPE)?)* ')' 'RETURN' DATA_TYPE createOperator_implementation_clause? createOperator_using_function_clause
    ;
createOperator_implementation_clause
    : (createOperator_context_clause | 'ANCILLARY' 'TO' ID '(' DATA_TYPE (',' DATA_TYPE)* ')')
    ;
createOperator_using_function_clause
    : 'USING' (ID '.')? (ID '.')? ID
    ;
createOperator_context_clause
    : ('WITH' 'INDEX' 'CONTEXT' ',' 'SCAN' 'CONTEXT' ID)? ('COMPUTE' 'ANCILLARY' 'DATA')? ('WITH' 'COLUMN' 'CONTEXT')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-OUTLINE.html
// fix: left out where optional, since no quoted page defines it or what it needs: statement
createOutline
    : 'CREATE' ('OR' 'REPLACE')? ('PUBLIC' | 'PRIVATE')? 'OUTLINE' ID? ('FOR' 'CATEGORY' ID)? 'FROM' 'PRIVATE'? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-PACKAGE-BODY.html
createPackageBody
    : 'CREATE' ('OR' 'REPLACE' | 'IF' 'NOT' 'EXISTS')? ('EDITIONABLE' | 'NONEDITIONABLE')? 'PACKAGE' 'BODY' (ID '.')? ID ('IS' | 'AS') PACKAGE_BODY_SOURCE
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-PACKAGE.html
createPackage
    : 'CREATE' ('OR' 'REPLACE' | 'IF' 'NOT' 'EXISTS')? ('EDITIONABLE' | 'NONEDITIONABLE')? 'PACKAGE' (ID '.')? ID ('SHARING' '=' ('METADATA' | 'NONE'))? ('ACCESSIBLE' 'BY' '(' createPackage_accessor (',' createPackage_accessor)* ')')? createPackage_invoker_rights_clause? ('IS' | 'AS') PACKAGE_SOURCE
    ;
createPackage_accessor
    : createPackage_unit_kind? (ID '.')? ID
    ;
createPackage_invoker_rights_clause
    : 'AUTHID' ('CURRENT_USER' | 'DEFINER')
    ;
createPackage_unit_kind
    : ('FUNCTION' | 'PROCEDURE' | 'PACKAGE' | 'TRIGGER' | 'TYPE')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-PFILE.html
createPfile
    : 'CREATE' 'PFILE' ('=' STRING)? 'FROM' ('SPFILE' ('=' STRING)? | 'MEMORY')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-PLUGGABLE-DATABASE.html
// fix: left out where optional, since no quoted page defines it or what it needs: snapshot_clause
createPluggableDatabase
    : 'CREATE' 'PLUGGABLE' 'DATABASE' ID (createPluggableDatabase_create_pdb_from_seed | createPluggableDatabase_create_pdb_clone | createPluggableDatabase_create_pdb_from_xml | createPluggableDatabase_create_pdb_from_mirror_copy | createPluggableDatabase_create_pdb_decrypt_from_xml)
    ;
createPluggableDatabase_create_pdb_from_seed
    : ('AS' ('APPLICATION' 'CONTAINER' | 'SEED'))? ('ADMIN' 'USER' ID 'IDENTIFIED' 'BY' ID)? createPluggableDatabase_pdb_dba_roles? createPluggableDatabase_parallel_pdb_creation_clause? createPluggableDatabase_default_tablespace? createPluggableDatabase_file_name_convert? createPluggableDatabase_service_name_convert? createPluggableDatabase_pdb_storage_clause? createPluggableDatabase_path_prefix_clause? createPluggableDatabase_tempfile_reuse_clause? createPluggableDatabase_user_tablespaces_clause? createPluggableDatabase_standbys_clause? createPluggableDatabase_logging_clause? createPluggableDatabase_create_file_dest_clause? createPluggableDatabase_container_map_clause? createPluggableDatabase_pdb_snapshot_clause?
    ;
createPluggableDatabase_create_pdb_clone
    : ('FROM' ID ('@' ID)? | 'AS' 'PROXY' 'FROM' ID '@' ID ('HOST' STRING)? ('PORT' NUMBER)?) createPluggableDatabase_parallel_pdb_creation_clause? createPluggableDatabase_default_tablespace? createPluggableDatabase_pdb_storage_clause? createPluggableDatabase_file_name_convert? createPluggableDatabase_service_name_convert? createPluggableDatabase_path_prefix_clause? createPluggableDatabase_tempfile_reuse_clause? createPluggableDatabase_user_tablespaces_clause? createPluggableDatabase_standbys_clause? createPluggableDatabase_logging_clause? createPluggableDatabase_create_file_dest_clause? createPluggableDatabase_keystore_clause? createPluggableDatabase_pdb_refresh_mode_clause? createPluggableDatabase_container_map_clause? createPluggableDatabase_pdb_snapshot_clause? ('SNAPSHOT' 'COPY' ('NO' 'DATA')? | 'NO' 'DATA')? ('RELOCATE' ('KEEP' 'SOURCE')? ('REFRESH' 'MODE' ('MANUAL' | 'EVERY' NUMBER ('MINUTES' | 'HOURS') | 'NONE'))? ('AVAILABILITY' ('NORMAL' | 'MAX'))?)?
    ;
createPluggableDatabase_create_pdb_from_xml
    : 'USING' STRING createPluggableDatabase_source_file_name_convert? createPluggableDatabase_source_file_directory? (('COPY' | 'MOVE' | 'NOCOPY'))? createPluggableDatabase_file_name_convert? createPluggableDatabase_service_name_convert? createPluggableDatabase_default_tablespace? createPluggableDatabase_pdb_storage_clause? createPluggableDatabase_path_prefix_clause? createPluggableDatabase_tempfile_reuse_clause? createPluggableDatabase_user_tablespaces_clause? createPluggableDatabase_standbys_clause? createPluggableDatabase_logging_clause? createPluggableDatabase_create_file_dest_clause? ('AS' 'CLONE')? createPluggableDatabase_container_map_clause? createPluggableDatabase_pdb_snapshot_clause?
    ;
createPluggableDatabase_create_pdb_from_mirror_copy
    : 'FROM' ID ('@' ID)? createPluggableDatabase_using_snapshot_clause? createPluggableDatabase_parallel_pdb_creation_clause? createPluggableDatabase_default_tablespace? createPluggableDatabase_pdb_storage_clause? createPluggableDatabase_file_name_convert? createPluggableDatabase_service_name_convert? createPluggableDatabase_path_prefix_clause? createPluggableDatabase_tempfile_reuse_clause? createPluggableDatabase_user_tablespaces_clause? createPluggableDatabase_standbys_clause? createPluggableDatabase_logging_clause? createPluggableDatabase_create_file_dest_clause? createPluggableDatabase_container_map_clause? createPluggableDatabase_pdb_snapshot_clause?
    ;
createPluggableDatabase_create_pdb_decrypt_from_xml
    : 'USING' STRING createPluggableDatabase_source_file_name_convert? createPluggableDatabase_source_file_directory? (('COPY' | 'MOVE' | 'NOCOPY'))? createPluggableDatabase_file_name_convert? createPluggableDatabase_service_name_convert? createPluggableDatabase_default_tablespace? createPluggableDatabase_pdb_storage_clause? createPluggableDatabase_path_prefix_clause? createPluggableDatabase_tempfile_reuse_clause? createPluggableDatabase_user_tablespaces_clause? createPluggableDatabase_standbys_clause? createPluggableDatabase_logging_clause? createPluggableDatabase_create_file_dest_clause? ('DECRYPT' 'USING' ID)? ('KEYSTORE' 'IDENTIFIED' 'BY' ID)? createPluggableDatabase_container_map_clause? createPluggableDatabase_pdb_snapshot_clause?
    ;
createPluggableDatabase_pdb_dba_roles
    : 'ROLES' '=' '(' ID (',' ID)* ')'
    ;
createPluggableDatabase_parallel_pdb_creation_clause
    : 'PARALLEL' NUMBER?
    ;
createPluggableDatabase_default_tablespace
    : 'DEFAULT' 'TABLESPACE' ID STRING?
    ;
createPluggableDatabase_file_name_convert
    : 'FILE_NAME_CONVERT' '=' '(' STRING ',' STRING (',' STRING ',' STRING)* ')'
    | 'FILE_NAME_CONVERT' '=' 'NONE'
    ;
createPluggableDatabase_service_name_convert
    : 'SERVICE_NAME_CONVERT' '=' '(' STRING ',' STRING (',' STRING ',' STRING)* ')'
    | 'SERVICE_NAME_CONVERT' '=' 'NONE'
    ;
createPluggableDatabase_pdb_storage_clause
    : 'STORAGE' '(' ('MAXSIZE' (createPluggableDatabase_size_clause | 'UNLIMITED'))? ('MAX_AUDIT_SIZE' (createPluggableDatabase_size_clause | 'UNLIMITED'))? ('MAX_DIAG_SIZE' (createPluggableDatabase_size_clause | 'UNLIMITED'))? ')'
    | 'STORAGE' 'UNLIMITED'
    ;
createPluggableDatabase_path_prefix_clause
    : 'PATH_PREFIX' '=' (STRING | ID)
    | 'PATH_PREFIX' '=' 'NONE'
    ;
createPluggableDatabase_tempfile_reuse_clause
    : 'TEMPFILE' 'REUSE'
    ;
createPluggableDatabase_user_tablespaces_clause
    : 'USER_TABLESPACES' '=' '(' ID (',' ID)* ')'
    | 'USER_TABLESPACES' '=' 'ALL'
    | 'USER_TABLESPACES' '=' 'ALL' 'EXCEPT' '(' ID (',' ID)* ')'
    | 'USER_TABLESPACES' '=' 'NONE'
    ;
createPluggableDatabase_standbys_clause
    : 'STANDBYS' '=' '(' ID (',' ID)* ')'
    | 'STANDBYS' '=' 'ALL'
    | 'STANDBYS' '=' 'ALL' 'EXCEPT' '(' ID (',' ID)* ')'
    | 'STANDBYS' '=' 'NONE'
    ;
createPluggableDatabase_logging_clause
    : ('LOGGING' | 'NOLOGGING')
    ;
createPluggableDatabase_create_file_dest_clause
    : 'CREATE_FILE_DEST' '=' (STRING | ID)
    | 'CREATE_FILE_DEST' '=' 'NONE'
    ;
createPluggableDatabase_container_map_clause
    : 'CONTAINER_MAP' 'UPDATE' '(' (createPluggableDatabase_add_partition_clause | createPluggableDatabase_split_partition_clause) ')'
    ;
createPluggableDatabase_pdb_snapshot_clause
    : 'SNAPSHOT' '=' 'NONE'
    | 'SNAPSHOT' '=' 'MANUAL'
    | 'SNAPSHOT' '=' 'EVERY' NUMBER ('MINUTES' | 'HOURS')
    ;
createPluggableDatabase_keystore_clause
    : 'KEYSTORE' 'IDENTIFIED' 'BY' ID ('DECRYPT' 'USING' ID)?
    | 'KEYSTORE' 'EXTERNAL' 'STORE' ('DECRYPT' 'USING' ID)?
    ;
createPluggableDatabase_pdb_refresh_mode_clause
    : 'REFRESH' 'MODE' ('MANUAL' | 'EVERY' NUMBER ('MINUTES' | 'HOURS') | 'NONE')
    ;
createPluggableDatabase_source_file_name_convert
    : 'SOURCE_FILE_NAME_CONVERT' '=' '(' STRING ',' STRING (',' STRING ',' STRING)* ')'
    | 'SOURCE_FILE_NAME_CONVERT' '=' 'NONE'
    ;
createPluggableDatabase_source_file_directory
    : 'SOURCE_FILE_DIRECTORY' '=' (STRING | ID)
    | 'SOURCE_FILE_DIRECTORY' '=' 'NONE'
    ;
createPluggableDatabase_using_snapshot_clause
    : 'USING' 'SNAPSHOT' ID
    | 'USING' 'SNAPSHOT' 'AT' 'SCN' NUMBER
    | 'USING' 'SNAPSHOT' 'AT' 'TIME' EXPR
    ;
createPluggableDatabase_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T' | 'P')?
    ;
createPluggableDatabase_add_partition_clause
    : 'ADD' 'PARTITION' ID 'VALUES' 'LESS' 'THAN' '(' EXPR ')'
    ;
createPluggableDatabase_split_partition_clause
    : 'SPLIT' 'PARTITION' ID 'AT' '(' EXPR ')' 'INTO' '(' 'PARTITION' ID ',' 'PARTITION' ID ')'
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-PROCEDURE.html
createProcedure
    : 'CREATE' ('OR' 'REPLACE' | 'IF' 'NOT' 'EXISTS')? ('EDITIONABLE' | 'NONEDITIONABLE')? 'PROCEDURE' (ID '.')? ID ('SHARING' '=' ('METADATA' | 'NONE'))? PROCEDURE_SOURCE
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-PROFILE.html
createProfile
    : 'CREATE' 'MANDATORY'? 'PROFILE' ID 'LIMIT' (createProfile_resource_parameters | createProfile_password_parameters) ((createProfile_resource_parameters | createProfile_password_parameters))* ('CONTAINER' '=' ('ALL' | 'CURRENT'))?
    ;
createProfile_resource_parameters
    : ('SESSIONS_PER_USER' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'CPU_PER_SESSION' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'CPU_PER_CALL' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'CONNECT_TIME' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'IDLE_TIME' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'LOGICAL_READS_PER_SESSION' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'LOGICAL_READS_PER_CALL' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'PRIVATE_SGA' (createProfile_size_clause | 'UNLIMITED' | 'DEFAULT') | 'COMPOSITE_LIMIT' (NUMBER | 'UNLIMITED' | 'DEFAULT'))
    ;
createProfile_password_parameters
    : ('FAILED_LOGIN_ATTEMPTS' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'PASSWORD_LIFE_TIME' (EXPR | 'UNLIMITED' | 'DEFAULT') | 'PASSWORD_REUSE_TIME' (EXPR | 'UNLIMITED' | 'DEFAULT') | 'PASSWORD_REUSE_MAX' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'PASSWORD_LOCK_TIME' (EXPR | 'UNLIMITED' | 'DEFAULT') | 'PASSWORD_GRACE_TIME' (EXPR | 'UNLIMITED' | 'DEFAULT') | 'INACTIVE_ACCOUNT_TIME' (NUMBER | 'UNLIMITED' | 'DEFAULT') | 'PASSWORD_VERIFY_FUNCTION' (ID | 'NULL' | 'DEFAULT') | 'PASSWORD_ROLLOVER_TIME' (EXPR | 'UNLIMITED' | 'DEFAULT'))
    ;
createProfile_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T' | 'P' | 'E')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-RESTORE-POINT.html
createRestorePoint
    : 'CREATE' 'CLEAN'? 'RESTORE' 'POINT' ID ('FOR' 'PLUGGABLE' 'DATABASE' ID)? ('AS' 'OF' ('TIMESTAMP' | 'SCN') EXPR)? 'PRESERVE'? ('GUARANTEE' 'FLASHBACK' 'DATABASE')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-ROLE.html
createRole
    : 'CREATE' 'ROLE' ID ('NOT' 'IDENTIFIED' | 'IDENTIFIED' ('BY' ID | 'USING' (ID '.')? ID | 'EXTERNALLY' | 'GLOBALLY' ('AS' STRING)?))? ('CONTAINER' '=' ('ALL' | 'CURRENT'))?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-ROLLBACK-SEGMENT.html
// fix: left out where optional, since no quoted page defines it or what it needs: storage_clause
createRollbackSegment
    : 'CREATE' 'PUBLIC'? 'ROLLBACK' 'SEGMENT' ID ('TABLESPACE' ID)?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-SCHEMA.html
createSchema
    : 'CREATE' 'SCHEMA' 'AUTHORIZATION' ID (CREATE_TABLE | CREATE_VIEW | GRANT) ((CREATE_TABLE | CREATE_VIEW | GRANT))*
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-SEQUENCE.html
createSequence
    : 'CREATE' 'SEQUENCE' (ID '.')? ID ('IF' 'NOT' 'EXISTS')? ('SHARING' '=' ('METADATA' | 'DATA' | 'NONE'))? ('INCREMENT' 'BY' NUMBER)? ('START' 'WITH' NUMBER)? (('MAXVALUE' NUMBER | 'NOMAXVALUE'))? (('MINVALUE' NUMBER | 'NOMINVALUE'))? (('CYCLE' | 'NOCYCLE'))? (('CACHE' NUMBER | 'NOCACHE'))? (('ORDER' | 'NOORDER'))? (('KEEP' | 'NOKEEP'))? (('SCALE' ('EXTEND' | 'NOEXTEND') | 'NOSCALE'))? (('SESSION' | 'GLOBAL'))? 'SHARD'?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-SPFILE.html
createSpfile
    : 'CREATE' 'SPFILE' ('=' STRING)? 'FROM' ('PFILE' ('=' STRING)? ('AS' 'COPY')? | 'MEMORY')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-SYNONYM.html
createSynonym
    : 'CREATE' ('OR' 'REPLACE' | 'IF' 'NOT' 'EXISTS')? ('EDITIONABLE' | 'NONEDITIONABLE')? 'PUBLIC'? 'SYNONYM' (ID '.')? ID ('SHARING' '=' ('METADATA' | 'NONE'))? 'FOR' (ID '.')? ID ('@' ID)?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-TABLE.html
// fix: left out where optional, since no quoted page defines it or what it needs: create_index_statement, json_collection_table, json_parameters, oject_properties, opaque_format_spec, storage_clause, xmltype_virtual_columns
createTable
    : 'CREATE' (('GLOBAL' 'TEMPORARY' | 'PRIVATE' 'TEMPORARY' | 'SHARDED' | 'DUPLICATED'))? 'TABLE' ('IF' 'NOT' 'EXISTS')? (ID '.')? ID ('SHARING' '=' ('METADATA' | 'DATA' | 'EXTENDED' 'DATA' | 'NONE'))? (createTable_relational_table | createTable_object_table | createTable_xmltype_table) ('MEMOPTIMIZE' 'FOR' 'READ')? ('MEMOPTIMIZE' 'FOR' 'WRITE')? ('PARENT' (ID '.')? ID)?
    ;
createTable_relational_table
    : ('(' createTable_relational_properties ')')? createTable_immutable_table_clauses? createTable_blockchain_table_clauses? ('DEFAULT' 'COLLATION' ID)? ('ON' 'COMMIT' ('DELETE' | 'PRESERVE') 'ROWS')? createTable_physical_properties? createTable_table_properties?
    ;
createTable_object_table
    : 'OF' (ID '.')? ID createTable_object_table_substitution? ('(' createTable_object_properties ')')? ('ON' 'COMMIT' ('DELETE' | 'PRESERVE') 'ROWS')? createTable_oid_clause? createTable_oid_index_clause? createTable_physical_properties? createTable_table_properties?
    ;
createTable_xmltype_table
    : 'OF' 'XMLTYPE' createTable_xmltype_storage? createTable_xmlschema_spec? ('ON' 'COMMIT' ('DELETE' | 'PRESERVE') 'ROWS')? createTable_oid_clause? createTable_oid_index_clause? createTable_physical_properties? createTable_table_properties?
    ;
createTable_relational_properties
    : (createTable_column_definition | createTable_virtual_column_definition | createTable_period_definition | createTable_out_of_line_constraint | createTable_out_of_line_ref_constraint | createTable_supplemental_logging_props | createTable_domain_clause) (',' (createTable_column_definition | createTable_virtual_column_definition | createTable_period_definition | createTable_out_of_line_constraint | createTable_out_of_line_ref_constraint | createTable_supplemental_logging_props | createTable_domain_clause))*
    ;
createTable_immutable_table_clauses
    : createTable_immutable_table_no_drop_clause? createTable_immutable_table_no_delete_clause?
    ;
createTable_blockchain_table_clauses
    : createTable_blockchain_drop_table_clause? createTable_blockchain_row_retention_clause? createTable_blockchain_hash_and_data_format_clause?
    ;
createTable_physical_properties
    : (createTable_deferred_segment_creation | createTable_segment_attributes_clause | createTable_table_compression | createTable_inmemory_table_clause | createTable_ilm_clause | createTable_heap_org_table_clause | createTable_index_org_table_clause | createTable_external_table_clause)
    ;
createTable_table_properties
    : createTable_column_properties? createTable_read_only_clause? createTable_indexing_clause? createTable_table_partitioning_clauses? createTable_attribute_clustering_clause? ('CACHE' | 'NOCACHE')? ('RESULT_CACHE' '(' 'MODE' ('DEFAULT' | 'FORCE') ')')? createTable_parallel_clause? ('ROWDEPENDENCIES' | 'NOROWDEPENDENCIES')? (createTable_enable_disable_clause createTable_enable_disable_clause?)* createTable_row_movement_clause? createTable_logical_replication_clause? createTable_flashback_archive_clause? ('AS' QUERY)?
    ;
createTable_object_table_substitution
    : 'NOT'? 'SUBSTITUTABLE' 'AT' 'ALL' 'LEVELS'
    ;
createTable_object_properties
    : (ID ('DEFAULT' EXPR)? ((createTable_inline_constraint | createTable_inline_ref_constraint)+)? | createTable_out_of_line_constraint | createTable_out_of_line_ref_constraint | createTable_supplemental_logging_props) (',' (ID ('DEFAULT' EXPR)? ((createTable_inline_constraint | createTable_inline_ref_constraint)+)? | createTable_out_of_line_constraint | createTable_out_of_line_ref_constraint | createTable_supplemental_logging_props))*
    ;
createTable_oid_clause
    : 'OBJECT' 'IDENTIFIER' 'IS' ('SYSTEM' 'GENERATED' | 'PRIMARY' 'KEY')
    ;
createTable_oid_index_clause
    : 'OIDINDEX' ID? '(' (createTable_physical_attributes_clause | 'TABLESPACE' ID)* ')'
    ;
createTable_xmltype_storage
    : 'STORE' 'AS' ('OBJECT' 'RELATIONAL' | ('SECUREFILE' | 'BASICFILE')? ('CLOB' | 'BINARY' 'XML') (ID ('(' createTable_lob_storage_parameters ')')? | '(' createTable_lob_storage_parameters ')')?)
    ;
createTable_xmlschema_spec
    : ('XMLSCHEMA' STRING)? 'ELEMENT' ID ('STORE' 'ALL' 'VARRAYS' 'AS' ('LOBS' | 'TABLES'))? (('ALLOW' | 'DISALLOW') 'NONSCHEMA')? (('ALLOW' | 'DISALLOW') 'ANYSCHEMA')?
    ;
createTable_column_definition
    : ID (DATA_TYPE | 'DOMAIN' (ID '.')? ID)? 'SORT'? ('VISIBLE' | 'INVISIBLE')? ('DEFAULT' ('ON' 'NULL' ('FOR' 'INSERT' 'ONLY')?)? EXPR | createTable_identity_clause)? ('ENCRYPT' createTable_encryption_spec)? ('COLLATE' ID)? (createTable_inline_constraint createTable_inline_constraint?)* createTable_inline_ref_constraint? createTable_annotations_clause?
    ;
createTable_virtual_column_definition
    : ID DATA_TYPE? ('COLLATE' ID)? ('VISIBLE' | 'INVISIBLE')? ('GENERATED' 'ALWAYS')? 'AS' '(' EXPR ')' 'VIRTUAL'? createTable_evaluation_edition_clause? createTable_unusable_editions_clause? (createTable_inline_constraint createTable_inline_constraint?)* createTable_annotations_clause?
    ;
createTable_period_definition
    : 'PERIOD' 'FOR' ID ('(' ID ',' ID ')')?
    ;
createTable_out_of_line_constraint
    : ('CONSTRAINT' ID)? ('UNIQUE' '(' ID (',' ID)* ')' | 'PRIMARY' 'KEY' '(' ID (',' ID)* ')' | 'FOREIGN' 'KEY' '(' ID (',' ID)* ')' 'REFERENCES' (ID '.')? ID ('(' ID (',' ID)* ')')? ('ON' 'DELETE' ('CASCADE' | 'SET' 'NULL'))? | 'CHECK' '(' CONDITION ')') createTable_constraint_state?
    ;
createTable_out_of_line_ref_constraint
    : ('SCOPE' 'FOR' '(' ID ')' 'IS' (ID '.')? ID | 'REF' '(' ID ')' 'WITH' 'ROWID' | ('CONSTRAINT' ID)? 'FOREIGN' 'KEY' '(' ID ')' 'REFERENCES' (ID '.')? ID ('(' ID (',' ID)* ')')? ('ON' 'DELETE' ('CASCADE' | 'SET' 'NULL'))? createTable_constraint_state?)
    ;
createTable_supplemental_logging_props
    : 'SUPPLEMENTAL' 'LOG' (createTable_supplemental_log_grp_clause | createTable_supplemental_id_key_clause)
    ;
createTable_domain_clause
    : 'DOMAIN' (ID '.')? ID '(' ID (',' ID)* ')'
    ;
createTable_immutable_table_no_drop_clause
    : 'NO' 'DROP' ('UNTIL' NUMBER 'DAYS' 'IDLE')?
    ;
createTable_immutable_table_no_delete_clause
    : 'NO' 'DELETE' ('LOCKED' | 'UNTIL' NUMBER 'DAYS' 'AFTER' 'INSERT')?
    ;
createTable_blockchain_drop_table_clause
    : 'NO' 'DROP' ('UNTIL' NUMBER 'DAYS' 'IDLE')?
    ;
createTable_blockchain_row_retention_clause
    : 'NO' 'DELETE' ('LOCKED' | 'UNTIL' NUMBER 'DAYS' 'AFTER' 'INSERT')
    ;
createTable_blockchain_hash_and_data_format_clause
    : 'HASHING' 'USING' STRING 'VERSION' STRING
    ;
createTable_deferred_segment_creation
    : 'SEGMENT' 'CREATION' ('IMMEDIATE' | 'DEFERRED')
    ;
createTable_segment_attributes_clause
    : (createTable_physical_attributes_clause | 'TABLESPACE' ID | createTable_logging_clause)
    ;
createTable_table_compression
    : ('NOCOMPRESS' | 'ROW' 'STORE' 'COMPRESS' ('BASIC' | 'ADVANCED')? | 'COLUMN' 'STORE' 'COMPRESS' 'FOR' ('QUERY' | 'ARCHIVE') ('LOW' | 'HIGH')?)
    ;
createTable_inmemory_table_clause
    : ('INMEMORY' createTable_inmemory_attributes? createTable_inmemory_column_clause? | 'NO' 'INMEMORY')
    ;
createTable_ilm_clause
    : 'ILM' ('ADD' 'POLICY' createTable_ilm_policy_clause | ('DELETE' | 'ENABLE' | 'DISABLE') 'POLICY' ID | 'DELETE_ALL' | 'ENABLE_ALL' | 'DISABLE_ALL')
    ;
createTable_heap_org_table_clause
    : ('ORGANIZATION' 'HEAP')?
    ;
createTable_index_org_table_clause
    : 'ORGANIZATION' 'INDEX' createTable_segment_attributes_clause? createTable_index_org_table_compression? (('MAPPING' 'TABLE' | 'NOMAPPING'))? createTable_index_org_overflow_clause?
    ;
createTable_external_table_clause
    : 'ORGANIZATION' 'EXTERNAL' '(' 'TYPE' ID ('DEFAULT' 'DIRECTORY' ID)? ('ACCESS' 'PARAMETERS' 'USING' 'CLOB' QUERY)? ('LOCATION' '(' (ID ':')? STRING (',' (ID ':')? STRING)* ')')? ')' ('REJECT' 'LIMIT' (NUMBER | 'UNLIMITED'))? createTable_parallel_clause? ('AS' QUERY)?
    ;
createTable_column_properties
    : (createTable_object_type_col_properties | createTable_nested_table_col_properties | (createTable_varray_col_properties | createTable_lob_storage_clause) ('(' createTable_lob_partition_storage (',' createTable_lob_partition_storage)* ')')? | createTable_xmltype_column_properties | createTable_json_storage_clause)+
    ;
createTable_read_only_clause
    : ('READ' 'ONLY' | 'READ' 'WRITE')
    ;
createTable_indexing_clause
    : 'INDEXING' ('ON' | 'OFF')
    ;
createTable_table_partitioning_clauses
    : (createTable_range_partitions | createTable_list_partitions | createTable_hash_partitions | createTable_composite_range_partitions | createTable_composite_list_partitions | createTable_composite_hash_partitions | createTable_reference_partitioning | createTable_system_partitioning | createTable_consistent_hash_partitions | createTable_consistent_hash_with_subpartitions | createTable_partitionset_clauses)
    ;
createTable_attribute_clustering_clause
    : 'CLUSTERING' ((ID '.')? ID)? createTable_clustering_join 'BY' ('LINEAR' 'ORDER' | 'INTERLEAVED' 'ORDER')? '(' createTable_clustering_columns ')' createTable_clustering_when? createTable_zonemap_clause?
    ;
createTable_parallel_clause
    : ('NOPARALLEL' | 'PARALLEL' NUMBER?)
    ;
createTable_enable_disable_clause
    : ('ENABLE' | 'DISABLE') ('VALIDATE' | 'NOVALIDATE')? ('UNIQUE' '(' ID (',' ID)* ')' | 'PRIMARY' 'KEY' | 'CONSTRAINT' ID) createTable_using_index_clause? createTable_exceptions_clause? 'CASCADE'? (('KEEP' | 'DROP') 'INDEX')?
    ;
createTable_row_movement_clause
    : ('ENABLE' | 'DISABLE') 'ROW' 'MOVEMENT'
    ;
createTable_logical_replication_clause
    : (('ENABLE' | 'DISABLE') 'LOGICAL' 'REPLICATION')?
    ;
createTable_flashback_archive_clause
    : ('FLASHBACK' 'ARCHIVE' ID? | 'NO' 'FLASHBACK' 'ARCHIVE')
    ;
createTable_inline_constraint
    : ('CONSTRAINT' ID)? ('NOT'? 'NULL' | 'UNIQUE' | 'PRIMARY' 'KEY' | 'REFERENCES' (ID '.')? ID ('(' ID ')')? ('ON' 'DELETE' ('CASCADE' | 'SET' 'NULL'))? | 'CHECK' '(' CONDITION ')') createTable_constraint_state?
    ;
createTable_inline_ref_constraint
    : ('SCOPE' 'IS' (ID '.')? ID | 'WITH' 'ROWID' | ('CONSTRAINT' ID)? 'REFERENCES' (ID '.')? ID ('(' ID (',' ID)* ')')? ('ON' 'DELETE' ('CASCADE' | 'SET' 'NULL'))? createTable_constraint_state?)
    ;
createTable_physical_attributes_clause
    : ('PCTFREE' NUMBER)? ('PCTUSED' NUMBER)? ('INITRANS' NUMBER)?
    ;
createTable_lob_storage_parameters
    : (('TABLESPACE' ID | createTable_logging_clause)+ | ('ENABLE' | 'DISABLE') 'STORAGE' 'IN' 'ROW' | 'CHUNK' NUMBER | 'PCTVERSION' NUMBER | 'FREEPOOLS' NUMBER | createTable_lob_retention_clause | createTable_lob_deduplicate_clause | createTable_lob_compression_clause | ('ENCRYPT' createTable_encryption_spec | 'DECRYPT') | ('CACHE' | 'NOCACHE' | 'CACHE' 'READS') createTable_logging_clause?)+
    ;
createTable_identity_clause
    : 'GENERATED' ('ALWAYS' | 'BY' 'DEFAULT' ('ON' 'NULL')?) 'AS' 'IDENTITY' ('(' createTable_identity_options ')')?
    ;
createTable_encryption_spec
    : ('USING' STRING)? ('IDENTIFIED' 'BY' ID)? ('SALT' | 'NO' 'SALT')?
    ;
createTable_annotations_clause
    : 'ANNOTATIONS' '(' ('ADD' | 'DROP') ID EXPR? (',' ID EXPR?)* ')'
    ;
createTable_evaluation_edition_clause
    : 'EVALUATE' 'USING' ('CURRENT' 'EDITION' | 'EDITION' ID)
    | 'EVALUATE' 'USING' 'NULL' 'EDITION'
    ;
createTable_unusable_editions_clause
    : ('UNUSABLE' 'BEFORE' ('CURRENT' 'EDITION' | 'EDITION' ID))? ('UNUSABLE' 'BEGINNING' 'WITH' ('CURRENT' 'EDITION' | 'EDITION' ID | 'NULL' 'EDITION'))?
    ;
createTable_constraint_state
    : ('NOT'? 'DEFERRABLE' ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? | 'INITIALLY' ('DEFERRED' | 'IMMEDIATE') ('NOT'? 'DEFERRABLE')?)? ('RELY' | 'NORELY')? createTable_using_index_clause? ('ENABLE' | 'DISABLE')? ('VALIDATE' | 'NOVALIDATE')? createTable_exceptions_clause?
    ;
createTable_supplemental_log_grp_clause
    : 'GROUP' ID '(' ID ('NO' 'LOG')? (',' ID ('NO' 'LOG')?)* ')' 'ALWAYS'?
    ;
createTable_supplemental_id_key_clause
    : 'DATA' '(' ('ALL' | 'PRIMARY' 'KEY' | 'UNIQUE' | 'FOREIGN' 'KEY') (',' ('ALL' | 'PRIMARY' 'KEY' | 'UNIQUE' | 'FOREIGN' 'KEY'))* ')' 'COLUMNS'
    ;
createTable_logging_clause
    : ('LOGGING' | 'NOLOGGING' | 'FILESYSTEM_LIKE_LOGGING')
    ;
createTable_inmemory_attributes
    : createTable_inmemory_memcompress? createTable_inmemory_priority? createTable_inmemory_distribute? createTable_inmemory_duplicate?
    ;
createTable_inmemory_column_clause
    : ('INMEMORY' createTable_inmemory_memcompress? | 'NO' 'INMEMORY') '(' ID (',' ID)* ')' (('INMEMORY' createTable_inmemory_memcompress? | 'NO' 'INMEMORY') '(' ID (',' ID)* ')')*
    ;
createTable_ilm_policy_clause
    : (createTable_ilm_compression_policy | createTable_ilm_tiering_policy | createTable_ilm_inmemory_policy)
    ;
createTable_index_org_table_compression
    : ('NOCOMPRESS' | 'COMPRESS' NUMBER | 'COMPRESS' 'ADVANCED' ('LOW' | 'HIGH')?)
    ;
createTable_index_org_overflow_clause
    : ('INCLUDING' ID)? 'OVERFLOW' createTable_segment_attributes_clause?
    ;
createTable_object_type_col_properties
    : 'COLUMN' ID createTable_substitutable_column_clause
    ;
createTable_nested_table_col_properties
    : 'NESTED' 'TABLE' (ID | 'COLUMN_VALUE') createTable_substitutable_column_clause? ('LOCAL' | 'GLOBAL')? 'STORE' 'AS' ID ('(' '(' createTable_object_properties ')' createTable_physical_properties? createTable_column_properties? ')')? ('RETURN' 'AS'? ('LOCATOR' | 'VALUE'))?
    ;
createTable_varray_col_properties
    : 'VARRAY' ID (createTable_substitutable_column_clause? 'STORE' 'AS' ('SECUREFILE' | 'BASICFILE')? 'LOB' (ID? '(' createTable_lob_storage_parameters ')' | ID) | createTable_substitutable_column_clause)
    ;
createTable_lob_storage_clause
    : 'LOB' ('(' ID (',' ID)* ')' 'STORE' 'AS' ('SECUREFILE' | 'BASICFILE')? ('(' createTable_lob_storage_parameters ')' | ID ('(' createTable_lob_storage_parameters ')')?) | '(' ID ')' 'STORE' 'AS' ('SECUREFILE' | 'BASICFILE')? (ID | '(' createTable_lob_storage_parameters ')' | ID '(' createTable_lob_storage_parameters ')'))
    ;
createTable_lob_partition_storage
    : 'PARTITION' ID (createTable_lob_storage_clause | createTable_varray_col_properties)+ ('(' 'SUBPARTITION' ID (createTable_lob_storage_clause | createTable_varray_col_properties)+ (')')+)?
    ;
createTable_xmltype_column_properties
    : 'XMLTYPE' 'COLUMN'? ID createTable_xmltype_storage? createTable_xmlschema_spec?
    ;
createTable_json_storage_clause
    : 'JSON' '(' ID ')' 'STORE' 'AS'
    ;
createTable_range_partitions
    : 'PARTITION' 'BY' 'RANGE' '(' ID (',' ID)* ')' ('INTERVAL' '(' EXPR ')' ('STORE' 'IN' '(' ID (',' ID)* ')')?)? '(' 'PARTITION' ID? createTable_range_values_clause createTable_table_partition_description? (',' 'PARTITION' ID? createTable_range_values_clause createTable_table_partition_description?)* ')'
    ;
createTable_list_partitions
    : 'PARTITION' 'BY' 'LIST' '(' ID (',' ID)* ')' ('AUTOMATIC' ('STORE' 'IN' '(' ID (',' ID)* ')')?)? '(' 'PARTITION' ID? createTable_list_values_clause createTable_table_partition_description? (',' 'PARTITION' ID? createTable_list_values_clause createTable_table_partition_description?)* ')'
    ;
createTable_hash_partitions
    : 'PARTITION' 'BY' 'HASH' '(' ID (',' ID)* ')' (createTable_individual_hash_partitions | createTable_hash_partitions_by_quantity)
    ;
createTable_composite_range_partitions
    : 'PARTITION' 'BY' 'RANGE' '(' ID (',' ID)* ')' ('INTERVAL' '(' EXPR ')' ('STORE' 'IN' '(' ID (',' ID)* ')')?)? (createTable_subpartition_by_range | createTable_subpartition_by_list | createTable_subpartition_by_hash) '(' 'PARTITION' ID? createTable_range_values_clause (createTable_table_partition_description ((createTable_range_subpartition_desc | createTable_list_subpartition_desc | createTable_individual_hash_subparts | createTable_hash_subparts_by_quantity))? (',' createTable_table_partition_description ((createTable_range_subpartition_desc | createTable_list_subpartition_desc | createTable_individual_hash_subparts | createTable_hash_subparts_by_quantity))?)*)* ')'
    ;
createTable_composite_list_partitions
    : 'PARTITION' 'BY' 'LIST' '(' ID (',' ID)* ')' ('AUTOMATIC' ('STORE' 'IN' '(' ID (',' ID)* ')')?)? (createTable_subpartition_by_range | createTable_subpartition_by_list | createTable_subpartition_by_hash) '(' 'PARTITION' ID? createTable_list_values_clause (createTable_table_partition_description ((createTable_range_subpartition_desc | createTable_list_subpartition_desc | createTable_individual_hash_subparts | createTable_hash_subparts_by_quantity))? (',' createTable_table_partition_description ((createTable_range_subpartition_desc | createTable_list_subpartition_desc | createTable_individual_hash_subparts | createTable_hash_subparts_by_quantity))?)*)* ')'
    ;
createTable_composite_hash_partitions
    : 'PARTITION' 'BY' 'HASH' '(' ID (',' ID)* ')' (createTable_subpartition_by_range | createTable_subpartition_by_list | createTable_subpartition_by_hash) (createTable_individual_hash_partitions | createTable_hash_partitions_by_quantity)
    ;
createTable_reference_partitioning
    : 'PARTITION' 'BY' 'REFERENCE' '(' ID ')' ('(' 'PARTITION' ID? createTable_table_partition_description? (',' 'PARTITION' ID? createTable_table_partition_description?)* ')')?
    ;
createTable_system_partitioning
    : ('PARTITION' 'BY' 'SYSTEM')? ('PARTITIONS' NUMBER | '(' 'PARTITION' ID? createTable_table_partition_description? (',' 'PARTITION' ID? createTable_table_partition_description?)* ')')?
    ;
createTable_consistent_hash_partitions
    : 'PARTITION' 'BY' 'CONSISTENT' 'HASH' '(' ID (',' ID)* ')' ('PARTITIONS' 'AUTO')? 'TABLESPACE' 'SET' ID
    ;
createTable_consistent_hash_with_subpartitions
    : 'PARTITION' 'BY' 'CONSISTENT' 'HASH' '(' ID (',' ID)* ')' (createTable_subpartition_by_range | createTable_subpartition_by_list | createTable_subpartition_by_hash) ('PARTITIONS' 'AUTO')? 'TABLESPACE' 'SET' ID
    ;
createTable_partitionset_clauses
    : (createTable_range_partitionset_clause | createTable_list_partitionset_clause)
    ;
createTable_clustering_join
    : 'JOIN' (ID '.')? ID 'ON' '(' ID '=' ID ')' ('JOIN' (ID '.')? ID 'ON' '(' ID '=' ID ')')* ('YES' 'ON' 'LOAD' ('YES' 'ON' 'DATA' 'MOVEMENT')? | 'NO' 'ON' 'LOAD' ('NO' 'ON' 'DATA' 'MOVEMENT')?)?
    ;
createTable_clustering_columns
    : (ID | '(' ID (',' ID)* ')') (('ASC' | 'DESC'))? (',' (ID | '(' ID (',' ID)* ')') (('ASC' | 'DESC'))?)*
    ;
createTable_clustering_when
    : ('YES' | 'NO') 'ON' 'LOAD' (('YES' | 'NO') 'ON' 'DATA' 'MOVEMENT')?
    ;
createTable_zonemap_clause
    : ('WITH' 'MATERIALIZED' 'ZONEMAP' ('(' ID (',' ID)* ')')? | 'WITHOUT' 'MATERIALIZED' 'ZONEMAP')
    ;
createTable_using_index_clause
    : 'USING' 'INDEX' ((ID '.')? ID | createTable_index_properties)
    ;
createTable_exceptions_clause
    : 'EXCEPTIONS' 'INTO' (ID '.')? ID
    ;
createTable_lob_retention_clause
    : 'RETENTION' (('MAX' | 'MIN' NUMBER | 'AUTO' | 'NONE'))?
    ;
createTable_lob_deduplicate_clause
    : ('DEDUPLICATE' | 'KEEP_DUPLICATES')
    ;
createTable_lob_compression_clause
    : ('COMPRESS' ('HIGH' | 'MEDIUM' | 'LOW')? | 'NOCOMPRESS')
    ;
createTable_identity_options
    : ('START' 'WITH' (NUMBER | 'LIMIT' 'VALUE'))? ('INCREMENT' 'BY' NUMBER)? (('MAXVALUE' NUMBER | 'NOMAXVALUE'))? (('MINVALUE' NUMBER | 'NOMINVALUE'))? (('CYCLE' | 'NOCYCLE'))? (('CACHE' NUMBER | 'NOCACHE'))? (('ORDER' | 'NOORDER'))?
    ;
createTable_inmemory_memcompress
    : ('MEMCOMPRESS' 'FOR' ('DML' | 'QUERY' ('LOW' | 'HIGH')? | 'CAPACITY' ('LOW' | 'HIGH')?) | 'NO' 'MEMCOMPRESS')
    ;
createTable_inmemory_priority
    : 'PRIORITY' ('NONE' | 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL')
    ;
createTable_inmemory_distribute
    : 'DISTRIBUTE' ('AUTO' | 'BY' ('ROWID' 'RANGE' | 'PARTITION' | 'SUBPARTITION'))? ('FOR' 'SERVICE' ('DEFAULT' | 'ALL' | ID | 'NONE'))?
    ;
createTable_inmemory_duplicate
    : ('DUPLICATE' | 'DUPLICATE' 'ALL' | 'NO' 'DUPLICATE')
    ;
createTable_ilm_compression_policy
    : (createTable_table_compression | ('REBUILD' | 'MERGE')) ('SEGMENT' | 'GROUP' | 'ROW') ('AFTER' createTable_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION') | 'ON' ID)
    ;
createTable_ilm_tiering_policy
    : 'TIER' 'TO' ID ('READ' 'ONLY')? ('SEGMENT' | 'GROUP') ('AFTER' createTable_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION') | 'ON' ID)
    ;
createTable_ilm_inmemory_policy
    : ('SET' 'INMEMORY' createTable_inmemory_attributes? | 'MODIFY' 'INMEMORY' createTable_inmemory_memcompress | 'NO' 'INMEMORY') 'SEGMENT' ('AFTER' createTable_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION') | 'ON' ID)
    ;
createTable_substitutable_column_clause
    : ('ELEMENT' 'IS' 'OF' 'TYPE'? '(' 'ONLY'? ID ')' | 'NOT'? 'SUBSTITUTABLE' 'AT' 'ALL' 'LEVELS')
    ;
createTable_range_values_clause
    : 'VALUES' 'LESS' 'THAN' '(' (EXPR | 'MAXVALUE') (',' (EXPR | 'MAXVALUE'))* ')'
    ;
createTable_table_partition_description
    : createTable_deferred_segment_creation? createTable_read_only_clause? createTable_indexing_clause? createTable_segment_attributes_clause? (createTable_table_compression | createTable_index_compression)? createTable_inmemory_clause? createTable_ilm_clause? ('OVERFLOW' createTable_segment_attributes_clause?)? ((createTable_lob_storage_clause | createTable_varray_col_properties | createTable_nested_table_col_properties)+)?
    ;
createTable_list_values_clause
    : 'VALUES' '(' (createTable_list_values | 'DEFAULT') ')'
    ;
createTable_individual_hash_partitions
    : '(' 'PARTITION' ID? createTable_read_only_clause? createTable_indexing_clause? createTable_partitioning_storage_clause? (',' 'PARTITION' ID? createTable_read_only_clause? createTable_indexing_clause? createTable_partitioning_storage_clause?)* ')'
    ;
createTable_hash_partitions_by_quantity
    : 'PARTITIONS' NUMBER ('STORE' 'IN' '(' ID (',' ID)* ')')? (createTable_table_compression | createTable_index_compression)? ('OVERFLOW' 'STORE' 'IN' '(' ID (',' ID)* ')')?
    ;
createTable_subpartition_by_range
    : 'SUBPARTITION' 'BY' 'RANGE' '(' ID (',' ID)* ')' createTable_subpartition_template?
    ;
createTable_subpartition_by_list
    : 'SUBPARTITION' 'BY' 'LIST' '(' ID (',' ID)* ')' createTable_subpartition_template?
    ;
createTable_subpartition_by_hash
    : 'SUBPARTITION' 'BY' 'HASH' '(' ID (',' ID)* ')' ('SUBPARTITIONS' NUMBER ('STORE' 'IN' '(' ID (',' ID)* ')')? | createTable_subpartition_template)?
    ;
createTable_range_subpartition_desc
    : 'SUBPARTITION' ID? createTable_range_values_clause createTable_read_only_clause? createTable_indexing_clause? createTable_partitioning_storage_clause?
    ;
createTable_list_subpartition_desc
    : 'SUBPARTITION' ID? createTable_list_values_clause createTable_read_only_clause? createTable_indexing_clause? createTable_partitioning_storage_clause?
    ;
createTable_individual_hash_subparts
    : 'SUBPARTITION' ID? createTable_read_only_clause? createTable_indexing_clause? createTable_partitioning_storage_clause?
    ;
createTable_hash_subparts_by_quantity
    : 'SUBPARTITIONS' NUMBER ('STORE' 'IN' '(' ID (',' ID)* ')')?
    ;
createTable_range_partitionset_clause
    : 'PARTITIONSET' 'BY' 'RANGE' '(' ID (',' ID)* ')' 'PARTITION' 'BY' 'CONSISTENT' 'HASH' '(' ID (',' ID)* ')' (createTable_subpartition_by_range | createTable_subpartition_by_list | createTable_subpartition_by_hash)? 'PARTITIONS' 'AUTO' '(' 'PARTITIONSET' ID createTable_range_values_clause ('TABLESPACE' 'SET' ID)? createTable_lob_storage_clause? createTable_subpartition_template? (',' 'PARTITIONSET' ID createTable_range_values_clause ('TABLESPACE' 'SET' ID)? createTable_lob_storage_clause? createTable_subpartition_template?)* ')'
    ;
createTable_list_partitionset_clause
    : 'PARTITIONSET' 'BY' 'LIST' '(' ID ')' 'PARTITION' 'BY' 'CONSISTENT' 'HASH' '(' ID (',' ID)* ')' (createTable_subpartition_by_range | createTable_subpartition_by_list | createTable_subpartition_by_hash)? 'PARTITIONS' 'AUTO' '(' 'PARTITIONSET' ID createTable_list_values_clause ('TABLESPACE' 'SET' ID)? createTable_lob_storage_clause? createTable_subpartition_template? (',' 'PARTITIONSET' ID createTable_list_values_clause ('TABLESPACE' 'SET' ID)? createTable_lob_storage_clause? createTable_subpartition_template?)* ')'
    ;
createTable_index_properties
    : (createTable_global_partitioned_index | createTable_local_partitioned_index)? createTable_index_attributes? (createTable_domain_index_clause | createTable_xmlindex_clause)?
    ;
createTable_ilm_time_period
    : NUMBER ('DAY' | 'DAYS' | 'MONTH' | 'MONTHS' | 'YEAR' | 'YEARS')
    ;
createTable_index_compression
    : (createTable_prefix_compression | createTable_advanced_index_compression | 'NOCOMPRESS')
    ;
createTable_inmemory_clause
    : ('INMEMORY' createTable_inmemory_attributes? | 'NO' 'INMEMORY')
    ;
createTable_list_values
    : (EXPR | 'NULL') (',' (EXPR | 'NULL'))*
    ;
createTable_partitioning_storage_clause
    : ('TABLESPACE' ID)? ('OVERFLOW' ('TABLESPACE' ID)?)? createTable_table_compression? createTable_index_compression? createTable_inmemory_clause? createTable_ilm_clause? (createTable_lob_partitioning_storage (',' createTable_lob_partitioning_storage)?)* ('VARRAY' ID 'STORE' 'AS' ('SECUREFILE' | 'BASICFILE')? 'LOB' ID)?
    ;
createTable_subpartition_template
    : 'SUBPARTITION' 'TEMPLATE' '(' (createTable_range_subpartition_desc | createTable_list_subpartition_desc | createTable_individual_hash_subparts) (',' (createTable_range_subpartition_desc | createTable_list_subpartition_desc | createTable_individual_hash_subparts))* ')'
    | 'SUBPARTITION' 'TEMPLATE' NUMBER
    ;
createTable_global_partitioned_index
    : 'GLOBAL' 'PARTITION' 'BY' ('RANGE' '(' ID ')' '(' createTable_index_partitioning_clause (',' createTable_index_partitioning_clause)* ')' | 'HASH' '(' ID ')' (createTable_individual_hash_partitions | createTable_hash_partitions_by_quantity))
    ;
createTable_local_partitioned_index
    : 'LOCAL' (createTable_on_range_partitioned_table | createTable_on_list_partitioned_table | createTable_on_hash_partitioned_table | createTable_on_comp_partitioned_table)?
    ;
createTable_index_attributes
    : (createTable_physical_attributes_clause | 'TABLESPACE' (ID | 'DEFAULT') | createTable_logging_clause | createTable_index_compression | createTable_partial_index_clause | createTable_parallel_clause)+
    ;
createTable_domain_index_clause
    : 'INDEXTYPE' 'IS' (ID '.')? ID createTable_local_domain_index_clause? createTable_parallel_clause? ('PARAMETERS' '(' STRING ')')?
    ;
createTable_xmlindex_clause
    : 'INDEXTYPE' 'IS' 'XDB' '.' 'XMLINDEX' createTable_local_xmlindex_clause? createTable_parallel_clause? ('PARAMETERS' '(' STRING ')')?
    ;
createTable_prefix_compression
    : 'COMPRESS' NUMBER?
    ;
createTable_advanced_index_compression
    : 'COMPRESS' 'ADVANCED' ('LOW' | 'HIGH')?
    ;
createTable_lob_partitioning_storage
    : 'LOB' '(' ID ')' 'STORE' 'AS' ('SECUREFILE' | 'BASICFILE')? (ID ('(' 'TABLESPACE' ID (')')+)? | '(' 'TABLESPACE' ID (')')+)?
    ;
createTable_index_partitioning_clause
    : 'PARTITION' ID? 'VALUES' 'LESS' 'THAN' '(' ID ')' createTable_segment_attributes_clause?
    ;
createTable_on_range_partitioned_table
    : '(' 'PARTITION' ID? createTable_segment_attributes_clause? createTable_index_compression? (('USABLE' | 'UNUSABLE'))? (',' 'PARTITION' ID? createTable_segment_attributes_clause? createTable_index_compression? (('USABLE' | 'UNUSABLE'))?)* ')'
    ;
createTable_on_list_partitioned_table
    : '(' 'PARTITION' ID? createTable_segment_attributes_clause? createTable_index_compression? (('USABLE' | 'UNUSABLE'))? (',' 'PARTITION' ID? createTable_segment_attributes_clause? createTable_index_compression? (('USABLE' | 'UNUSABLE'))?)* ')'
    ;
createTable_on_hash_partitioned_table
    : ('(' 'PARTITION' ID? ('TABLESPACE' ID)? createTable_index_compression? (('USABLE' | 'UNUSABLE'))? (',' 'PARTITION' ID? ('TABLESPACE' ID)? createTable_index_compression? (('USABLE' | 'UNUSABLE'))?)* ')')? ('STORE' 'IN' '(' ID (',' ID)* ')')?
    ;
createTable_on_comp_partitioned_table
    : '(' 'PARTITION' ID? createTable_segment_attributes_clause? createTable_index_compression? createTable_index_subpartition_clause? (('USABLE' | 'UNUSABLE'))? (',' 'PARTITION' ID? createTable_segment_attributes_clause? createTable_index_compression? createTable_index_subpartition_clause? (('USABLE' | 'UNUSABLE'))?)* ')'
    ;
createTable_partial_index_clause
    : 'INDEXING' ('FULL' | 'PARTIAL')
    ;
createTable_local_domain_index_clause
    : 'LOCAL' ('(' 'PARTITION' ID? ('PARAMETERS' '(' STRING ')')? (',' 'PARTITION' ID? ('PARAMETERS' '(' STRING ')')?)* ')')?
    ;
createTable_local_xmlindex_clause
    : 'LOCAL' ('(' 'PARTITION' ID? ('PARAMETERS' '(' STRING ')')? (',' 'PARTITION' ID? ('PARAMETERS' '(' STRING ')')?)* ')')?
    ;
createTable_index_subpartition_clause
    : '(' 'SUBPARTITION' ID? ('TABLESPACE' ID)? createTable_index_compression? (('USABLE' | 'UNUSABLE'))? (',' 'SUBPARTITION' ID? ('TABLESPACE' ID)? createTable_index_compression? (('USABLE' | 'UNUSABLE'))?)* ')'
    | 'STORE' 'IN' '(' ID (',' ID)* ')'
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-TABLESPACE-SET.html
// fix: left out where optional, since no quoted page defines it or what it needs: storage_clause
createTablespaceSet
    : 'CREATE' 'TABLESPACE' 'SET' ID ('IN' 'SHARDSPACE' ID)? ('USING' 'TEMPLATE' '(' ('DATAFILE' createTablespaceSet_file_specification?)? createTablespaceSet_permanent_tablespace_attrs? ')')?
    ;
createTablespaceSet_file_specification
    : STRING? ('SIZE' createTablespaceSet_size_clause)? 'REUSE'? createTablespaceSet_autoextend_clause?
    ;
createTablespaceSet_permanent_tablespace_attrs
    : createTablespaceSet_logging_clause? createTablespaceSet_tablespace_encryption_clause? ('DEFAULT' createTablespaceSet_default_tablespace_params)? createTablespaceSet_extent_management_clause? createTablespaceSet_segment_management_clause? createTablespaceSet_flashback_mode_clause?
    ;
createTablespaceSet_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T' | 'P')?
    ;
createTablespaceSet_autoextend_clause
    : 'AUTOEXTEND' ('OFF' | 'ON' ('NEXT' createTablespaceSet_size_clause)? createTablespaceSet_maxsize_clause?)
    ;
createTablespaceSet_logging_clause
    : ('LOGGING' | 'NOLOGGING')
    ;
createTablespaceSet_tablespace_encryption_clause
    : 'ENCRYPTION' createTablespaceSet_tablespace_encryption_spec? ('ENCRYPT' | 'DECRYPT')
    ;
createTablespaceSet_default_tablespace_params
    : 'DEFAULT' createTablespaceSet_default_table_compression? createTablespaceSet_default_index_compression? createTablespaceSet_inmemory_clause? createTablespaceSet_ilm_clause?
    ;
createTablespaceSet_extent_management_clause
    : 'EXTENT' 'MANAGEMENT' ('LOCAL' (('AUTOALLOCATE' | 'UNIFORM' ('SIZE' createTablespaceSet_size_clause)?))? | 'DICTIONARY')
    ;
createTablespaceSet_segment_management_clause
    : 'SEGMENT' 'SPACE' 'MANAGEMENT' ('AUTO' | 'MANUAL')
    ;
createTablespaceSet_flashback_mode_clause
    : 'FLASHBACK' ('ON' | 'OFF')
    ;
createTablespaceSet_maxsize_clause
    : 'MAXSIZE' (createTablespaceSet_size_clause | 'UNLIMITED')
    ;
createTablespaceSet_tablespace_encryption_spec
    : 'USING' STRING
    ;
createTablespaceSet_default_table_compression
    : 'TABLE' ('COMPRESS' ('BASIC' | 'FOR' 'OLTP')? | 'ROW' 'STORE' 'COMPRESS' ('BASIC' | 'ADVANCED')? | 'COLUMN' 'STORE' 'COMPRESS' ('FOR' ('QUERY' | 'ARCHIVE') ('LOW' | 'HIGH')?)? | 'NOCOMPRESS')
    ;
createTablespaceSet_default_index_compression
    : 'INDEX' ('COMPRESS' ('ADVANCED' ('LOW' | 'HIGH')?)? | 'NOCOMPRESS')
    ;
createTablespaceSet_inmemory_clause
    : ('INMEMORY' createTablespaceSet_inmemory_attributes? | 'NO' 'INMEMORY')
    ;
createTablespaceSet_ilm_clause
    : 'ILM' ('ADD' 'POLICY' createTablespaceSet_ilm_policy_clause | 'DELETE' 'POLICY' ID | 'DELETE_ALL' | 'ENABLE_ALL' | 'DISABLE_ALL')
    ;
createTablespaceSet_inmemory_attributes
    : createTablespaceSet_inmemory_memcompress? createTablespaceSet_inmemory_priority? createTablespaceSet_inmemory_distribute? createTablespaceSet_inmemory_duplicate?
    ;
createTablespaceSet_ilm_policy_clause
    : (createTablespaceSet_ilm_compression_policy | createTablespaceSet_ilm_tiering_policy | createTablespaceSet_ilm_inmemory_policy)
    ;
createTablespaceSet_inmemory_memcompress
    : 'MEMCOMPRESS' 'FOR' ('QUERY' ('LOW' | 'HIGH')? | 'OLTP' | 'ARCHIVE' ('LOW' | 'HIGH')?)
    ;
createTablespaceSet_inmemory_priority
    : 'PRIORITY' ('LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL')
    ;
createTablespaceSet_inmemory_distribute
    : 'DISTRIBUTE' 'BY' ('ROWID' ('ACROSS' NUMBER?)? | 'PARTITION' | 'SUBPARTITION')
    ;
createTablespaceSet_inmemory_duplicate
    : 'DUPLICATE' ('ALL' | 'FOR' 'SERVICE' ID (',' ID)* | 'NONE')
    ;
createTablespaceSet_ilm_compression_policy
    : (createTablespaceSet_table_compression | 'ROW' 'STORE' 'COMPRESS' 'ADVANCED'? | 'COLUMN' 'STORE' 'COMPRESS' ('FOR' 'QUERY' ('LOW' | 'HIGH')?)?) ('SEGMENT' | 'GROUP' | 'ROW') 'AFTER' createTablespaceSet_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION')
    ;
createTablespaceSet_ilm_tiering_policy
    : 'TIER' 'TO' ID ('SEGMENT' | 'GROUP' | 'ROW')? ('ON' ID)?
    ;
createTablespaceSet_ilm_inmemory_policy
    : ('SET' 'INMEMORY' createTablespaceSet_inmemory_memcompress? createTablespaceSet_inmemory_priority? | 'MODIFY' 'INMEMORY' createTablespaceSet_inmemory_memcompress | 'NO' 'INMEMORY') 'SEGMENT' 'AFTER' createTablespaceSet_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION')
    ;
createTablespaceSet_table_compression
    : ('COMPRESS' 'BASIC'? | 'COMPRESS' 'FOR' ('OLTP' | 'QUERY' ('LOW' | 'HIGH')? | 'ARCHIVE' ('LOW' | 'HIGH')?) | 'NOCOMPRESS')
    ;
createTablespaceSet_ilm_time_period
    : NUMBER ('DAY' | 'DAYS' | 'MONTH' | 'MONTHS' | 'YEAR' | 'YEARS')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-TABLESPACE.html
// fix: left out where optional, since no quoted page defines it or what it needs: storage_clause
createTablespace
    : 'CREATE' ('BIGFILE' | 'SMALLFILE')? (createTablespace_permanent_tablespace_clause | createTablespace_temporary_tablespace_clause | createTablespace_undo_tablespace_clause)
    ;
createTablespace_permanent_tablespace_clause
    : 'TABLESPACE' ID ('IF' 'NOT' 'EXISTS')? ('DATAFILE' createTablespace_file_specification (',' createTablespace_file_specification)?)* createTablespace_permanent_tablespace_attrs?
    ;
createTablespace_temporary_tablespace_clause
    : 'LOCAL'? 'TEMPORARY' 'TABLESPACE' ID ('IF' 'NOT' 'EXISTS')? ('TEMPFILE' createTablespace_file_specification (',' createTablespace_file_specification)?)* createTablespace_tablespace_group_clause? createTablespace_extent_management_clause? createTablespace_tablespace_encryption_clause? ('FOR' ('ALL' | 'LEAF'))?
    ;
createTablespace_undo_tablespace_clause
    : 'UNDO' 'TABLESPACE' ID ('IF' 'NOT' 'EXISTS')? ('DATAFILE' createTablespace_file_specification (',' createTablespace_file_specification)?)* createTablespace_extent_management_clause? createTablespace_tablespace_retention_clause? createTablespace_tablespace_encryption_clause?
    ;
createTablespace_file_specification
    : STRING? ('SIZE' createTablespace_size_clause)? 'REUSE'? createTablespace_autoextend_clause?
    ;
createTablespace_permanent_tablespace_attrs
    : ('MINIMUM' 'EXTENT' createTablespace_size_clause)? ('BLOCKSIZE' NUMBER 'K'?)? createTablespace_logging_clause? ('FORCE' 'LOGGING')? createTablespace_tablespace_encryption_clause? ('DEFAULT' createTablespace_default_tablespace_params?)? (('ONLINE' | 'OFFLINE'))? createTablespace_extent_management_clause? createTablespace_segment_management_clause? createTablespace_flashback_mode_clause? createTablespace_lost_write_protection?
    ;
createTablespace_tablespace_group_clause
    : 'TABLESPACE' 'GROUP' (ID | STRING)
    ;
createTablespace_extent_management_clause
    : 'EXTENT' 'MANAGEMENT' ('LOCAL' (('AUTOALLOCATE' | 'UNIFORM' ('SIZE' createTablespace_size_clause)?))? | 'DICTIONARY')
    ;
createTablespace_tablespace_encryption_clause
    : 'ENCRYPTION' createTablespace_tablespace_encryption_spec? ('ENCRYPT' | 'DECRYPT')
    ;
createTablespace_tablespace_retention_clause
    : 'RETENTION' ('GUARANTEE' | 'NOGUARANTEE')
    ;
createTablespace_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T' | 'P' | 'E')?
    ;
createTablespace_autoextend_clause
    : ('AUTOEXTEND' 'OFF' | 'AUTOEXTEND' 'ON' ('NEXT' createTablespace_size_clause)? ('MAXSIZE' ('UNLIMITED' | createTablespace_size_clause))?)
    ;
createTablespace_logging_clause
    : ('LOGGING' | 'NOLOGGING' | 'FILESYSTEM_LIKE_LOGGING')
    ;
createTablespace_default_tablespace_params
    : (createTablespace_default_table_compression | createTablespace_default_index_compression | createTablespace_inmemory_clause | createTablespace_ilm_clause)+
    ;
createTablespace_segment_management_clause
    : 'SEGMENT' 'SPACE' 'MANAGEMENT' ('AUTO' | 'MANUAL')
    ;
createTablespace_flashback_mode_clause
    : 'FLASHBACK' ('ON' | 'OFF')
    ;
createTablespace_lost_write_protection
    : ('ENABLE' | 'DISABLE' | 'SUSPEND' | 'REMOVE') 'LOST' 'WRITE' 'PROTECTION'
    ;
createTablespace_tablespace_encryption_spec
    : 'USING' STRING
    ;
createTablespace_default_table_compression
    : 'TABLE' ('COMPRESS' 'FOR' ('OLTP' | 'QUERY' ('LOW' | 'HIGH')? | 'ARCHIVE' ('LOW' | 'HIGH')?) | 'NOCOMPRESS' | 'ROW' 'STORE' 'COMPRESS' ('BASIC' | 'ADVANCED')? | 'COLUMN' 'STORE' 'COMPRESS' 'FOR' ('QUERY' | 'ARCHIVE') ('LOW' | 'HIGH')?)
    ;
createTablespace_default_index_compression
    : 'INDEX' ('COMPRESS' ('ADVANCED' ('LOW' | 'HIGH')?)? | 'NOCOMPRESS')
    ;
createTablespace_inmemory_clause
    : ('INMEMORY' createTablespace_inmemory_attributes? | 'NO' 'INMEMORY')
    ;
createTablespace_ilm_clause
    : 'ILM' ('ADD' 'POLICY' createTablespace_ilm_policy_clause | 'DELETE' 'POLICY' ID | 'DELETE_ALL' | 'ENABLE_ALL' | 'DISABLE_ALL')
    ;
createTablespace_inmemory_attributes
    : createTablespace_inmemory_memcompress? createTablespace_inmemory_priority? createTablespace_inmemory_distribute? createTablespace_inmemory_duplicate?
    ;
createTablespace_ilm_policy_clause
    : (createTablespace_ilm_compression_policy | createTablespace_ilm_tiering_policy | createTablespace_ilm_inmemory_policy)
    ;
createTablespace_inmemory_memcompress
    : ('MEMCOMPRESS' 'FOR' ('DML' | 'QUERY' ('LOW' | 'HIGH')? | 'CAPACITY' ('LOW' | 'HIGH')?) | 'NO' 'MEMCOMPRESS')
    ;
createTablespace_inmemory_priority
    : 'PRIORITY' ('NONE' | 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL')
    ;
createTablespace_inmemory_distribute
    : 'DISTRIBUTE' ('AUTO' | 'BY' ('ROWID' 'RANGE' | 'PARTITION' | 'SUBPARTITION'))?
    ;
createTablespace_inmemory_duplicate
    : ('DUPLICATE' | 'DUPLICATE' 'ALL' | 'NO' 'DUPLICATE')
    ;
createTablespace_ilm_compression_policy
    : (createTablespace_table_compression | 'ROW' 'STORE' 'COMPRESS' 'ADVANCED'? | 'COLUMN' 'STORE' 'COMPRESS' ('FOR' 'QUERY' ('LOW' | 'HIGH')?)?) ('SEGMENT' | 'GROUP' | 'ROW') 'AFTER' createTablespace_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION')
    ;
createTablespace_ilm_tiering_policy
    : 'TIER' 'TO' ID ('SEGMENT' | 'GROUP' | 'ROW')? ('ON' ID)?
    ;
createTablespace_ilm_inmemory_policy
    : ('SET' 'INMEMORY' createTablespace_inmemory_memcompress? createTablespace_inmemory_priority? | 'MODIFY' 'INMEMORY' createTablespace_inmemory_memcompress | 'NO' 'INMEMORY') 'SEGMENT' 'AFTER' createTablespace_ilm_time_period 'OF' ('NO' 'ACCESS' | 'NO' 'MODIFICATION' | 'CREATION')
    ;
createTablespace_table_compression
    : ('COMPRESS' 'BASIC'? | 'COMPRESS' 'FOR' ('OLTP' | 'QUERY' ('LOW' | 'HIGH')? | 'ARCHIVE' ('LOW' | 'HIGH')?) | 'NOCOMPRESS')
    ;
createTablespace_ilm_time_period
    : NUMBER ('DAY' | 'DAYS' | 'MONTH' | 'MONTHS' | 'YEAR' | 'YEARS')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-TRIGGER.html
createTrigger
    : 'CREATE' ('OR' 'REPLACE' | 'IF' 'NOT' 'EXISTS')? ('EDITIONABLE' | 'NONEDITIONABLE')? 'TRIGGER' TRIGGER_SOURCE
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-TYPE-BODY.html
createTypeBody
    : 'CREATE' ('OR' 'REPLACE' | 'IF' 'NOT' 'EXISTS')? ('EDITIONABLE' | 'NONEDITIONABLE')? 'TYPE' 'BODY' (ID '.')? ID ('IS' | 'AS') TYPE_BODY_SOURCE 'END'
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-TYPE.html
createType
    : 'CREATE' ('OR' 'REPLACE' | 'IF' 'NOT' 'EXISTS')? ('EDITIONABLE' | 'NONEDITIONABLE')? 'TYPE' TYPE_SOURCE
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-USER.html
createUser
    : 'CREATE' 'USER' ('IF' 'NOT' 'EXISTS')? ID 'IDENTIFIED' ('BY' ID | 'EXTERNALLY' ('AS' STRING)? | 'GLOBALLY' ('AS' STRING)? | 'NO' 'AUTHENTICATION') ('AND' 'FACTOR' STRING 'AS' STRING)? ('HTTP'? 'DIGEST' ('ENABLE' | 'DISABLE'))? ('DEFAULT' 'COLLATION' ID)? ('DEFAULT' 'TABLESPACE' ID)? ('LOCAL'? 'TEMPORARY' 'TABLESPACE' ID)? ('QUOTA' (createUser_size_clause | 'UNLIMITED') 'ON' ID)* ('PROFILE' ID)? ('PASSWORD' 'EXPIRE')? ('ACCOUNT' ('LOCK' | 'UNLOCK'))? ('ENABLE' 'EDITIONS')? ('CONTAINER' '=' ('ALL' | 'CURRENT'))? ('READ' 'ONLY' | 'READ' 'WRITE')?
    ;
createUser_size_clause
    : NUMBER ('K' | 'M' | 'G' | 'T')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-VIEW.html
createView
    : 'CREATE' ('OR' 'REPLACE' | 'IF' 'NOT' 'EXISTS')? (('FORCE' | 'NOFORCE'))? 'EDITIONING'? ('EDITIONABLE' | 'NONEDITIONABLE')? ('SHARING' '=' ('METADATA' | 'DATA' | 'EXTENDED' 'DATA' | 'NONE'))? 'VIEW' (ID '.')? ID ('(' ID ('VISIBLE' | 'INVISIBLE')? createView_inline_constraint* (',' ID ('VISIBLE' | 'INVISIBLE')? createView_inline_constraint*)* (',' createView_out_of_line_constraint)* ')')? (createView_object_view_clause | createView_xmltype_view_clause)? createView_annotations_clause? ('DEFAULT' 'COLLATION' ID)? ('BEQUEATH' ('CURRENT_USER' | 'DEFINER'))? 'AS' QUERY createView_subquery_restriction_clause? ('CONTAINER_MAP' | 'CONTAINERS_DEFAULT')?
    ;
createView_inline_constraint
    : ('CONSTRAINT' ID)? ('NOT' 'NULL' | 'UNIQUE' | 'PRIMARY' 'KEY' | 'REFERENCES' (ID '.')? ID ('(' ID ')')? ('ON' 'DELETE' ('CASCADE' | 'SET' 'NULL'))? | 'CHECK' '(' CONDITION ')') ('RELY' | 'NORELY')? ('ENABLE' | 'DISABLE')? ('VALIDATE' | 'NOVALIDATE')?
    ;
createView_out_of_line_constraint
    : ('CONSTRAINT' ID)? ('UNIQUE' '(' ID (',' ID)* ')' | 'PRIMARY' 'KEY' '(' ID (',' ID)* ')' | 'FOREIGN' 'KEY' '(' ID (',' ID)* ')' 'REFERENCES' (ID '.')? ID ('(' ID (',' ID)* ')')? ('ON' 'DELETE' ('CASCADE' | 'SET' 'NULL'))? | 'CHECK' '(' CONDITION ')') ('RELY' | 'NORELY')? ('ENABLE' | 'DISABLE')? ('VALIDATE' | 'NOVALIDATE')?
    ;
createView_object_view_clause
    : 'OF' (ID '.')? ID ('WITH' 'OBJECT' 'IDENTIFIER' ('DEFAULT' | '(' ID (',' ID)* ')') | 'UNDER' (ID '.')? ID)
    ;
createView_xmltype_view_clause
    : 'OF' 'XMLTYPE' createView_xmlschema_spec ('WITH' 'OBJECT' 'IDENTIFIER' ('DEFAULT' | '(' EXPR (',' EXPR)* ')'))?
    ;
createView_annotations_clause
    : 'ANNOTATIONS' '(' ID EXPR? (',' ID EXPR?)* ')'
    ;
createView_subquery_restriction_clause
    : ('WITH' 'READ' 'ONLY' ('CONSTRAINT' ID)? | 'WITH' 'CHECK' 'OPTION' ('CONSTRAINT' ID)?)
    ;
createView_xmlschema_spec
    : 'XMLSCHEMA' STRING 'ELEMENT' ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DELETE.html
// fix: left out where optional, since no quoted page defines it or what it needs: create_index_statement, opaque_format_spec, storage_clause
delete
    : 'DELETE' HINT? 'FROM'? delete_dml_table_expression_clause delete_from_using_clause? delete_where_clause? delete_returning_clause? delete_error_logging_clause?
    ;
delete_dml_table_expression_clause
    : ((ID '.')? ID delete_partition_extension_clause? ('@' ID)? ID? | '(' QUERY delete_subquery_restriction_clause? ')' ID? | 'TABLE' '(' EXPR ')' ('(' '+' ')')? ID?)
    ;
delete_from_using_clause
    : 'FROM' (delete_table_reference | 'LATERAL'? QUERY delete_join_clause? | delete_join_clause)
    ;
delete_where_clause
    : 'WHERE' CONDITION
    ;
delete_returning_clause
    : 'RETURNING' EXPR (',' EXPR)* 'INTO' ID (',' ID)*
    ;
delete_error_logging_clause
    : 'LOG' 'ERRORS' ('INTO' (ID '.')? ID)? ('(' EXPR ')')? ('REJECT' 'LIMIT' (NUMBER | 'UNLIMITED'))?
    ;
delete_partition_extension_clause
    : ('PARTITION' '(' ID ')' | 'SUBPARTITION' '(' ID ')')
    ;
delete_subquery_restriction_clause
    : ('WITH' 'READ' 'ONLY' | 'WITH' 'CHECK' 'OPTION' ('CONSTRAINT' ID)?)
    ;
delete_table_reference
    : (delete_query_table_expression | '(' delete_join_clause ')' | delete_inline_external_table) ID?
    ;
delete_join_clause
    : delete_table_reference (delete_inner_cross_join_clause | delete_outer_join_clause | delete_cross_outer_apply_clause)+
    ;
delete_query_table_expression
    : ((ID '.')? ID delete_partition_extension_clause? ('@' ID)? | '(' QUERY delete_subquery_restriction_clause? ')' | delete_table_collection_expression | delete_inline_analytic_view) delete_flashback_query_clause? delete_sample_clause? (delete_pivot_clause | delete_unpivot_clause)? delete_row_pattern_clause? delete_containers_clause? delete_shards_clause?
    ;
delete_inline_external_table
    : 'EXTERNAL' '(' '(' delete_column_definition (',' delete_column_definition)* ')' delete_inline_external_table_properties ')'
    ;
delete_inner_cross_join_clause
    : ('INNER'? 'JOIN' delete_table_reference ('ON' CONDITION | 'USING' '(' ID (',' ID)* ')') | 'CROSS' 'JOIN' delete_table_reference | 'NATURAL' 'INNER'? 'JOIN' delete_table_reference)
    ;
delete_outer_join_clause
    : (delete_query_partition_clause delete_outer_join_type 'JOIN' delete_table_reference ('ON' CONDITION | 'USING' '(' ID (',' ID)* ')') | delete_outer_join_type 'JOIN' delete_query_partition_clause delete_table_reference ('ON' CONDITION | 'USING' '(' ID (',' ID)* ')') | 'NATURAL' delete_outer_join_type 'JOIN' delete_table_reference)
    ;
delete_cross_outer_apply_clause
    : ('CROSS' | 'OUTER') 'APPLY' (delete_table_reference | EXPR)
    ;
delete_table_collection_expression
    : 'TABLE' '(' EXPR ')' ('(' '+' ')')?
    ;
delete_inline_analytic_view
    : 'ANALYTIC' 'VIEW' '(' delete_sub_av_clause ')' ID?
    ;
delete_flashback_query_clause
    : ('VERSIONS' 'BETWEEN' ('SCN' | 'TIMESTAMP') EXPR 'AND' EXPR | 'VERSIONS' 'PERIOD' 'FOR' ID 'BETWEEN' EXPR 'AND' EXPR | 'AS' 'OF' ('SCN' | 'TIMESTAMP') EXPR | 'AS' 'OF' 'PERIOD' 'FOR' ID EXPR)
    ;
delete_sample_clause
    : 'SAMPLE' 'BLOCK'? '(' NUMBER ')' ('SEED' '(' NUMBER ')')?
    ;
delete_pivot_clause
    : 'PIVOT' 'XML'? '(' ID '(' EXPR ')' ('AS'? ID)? (',' ID '(' EXPR ')' ('AS'? ID)?)* delete_pivot_for_clause delete_pivot_in_clause ')'
    ;
delete_unpivot_clause
    : 'UNPIVOT' (('INCLUDE' | 'EXCLUDE') 'NULLS')? '(' (ID | '(' ID (',' ID)* ')') delete_pivot_for_clause delete_unpivot_in_clause ')'
    ;
delete_row_pattern_clause
    : 'MATCH_RECOGNIZE' '(' ('PARTITION' 'BY' ID (',' ID)?)* ('ORDER' 'BY' ID ('ASC' | 'DESC')? (',' ID ('ASC' | 'DESC')?)?)* ('MEASURES' delete_row_pattern_measure_column (',' delete_row_pattern_measure_column)?)* ('ONE' 'ROW' 'PER' 'MATCH' | 'ALL' 'ROWS' 'PER' 'MATCH' (('SHOW' | 'OMIT') 'EMPTY' 'MATCHES')?)? ('AFTER' 'MATCH' 'SKIP' ('PAST' 'LAST' 'ROW' | 'TO' 'NEXT' 'ROW' | 'TO' 'FIRST' ID | 'TO' 'LAST' ID | 'TO' ID))? 'PATTERN' '(' delete_row_pattern ')' ('SUBSET' delete_subset_item (',' delete_subset_item)?)* 'DEFINE' delete_row_pattern_definition (',' delete_row_pattern_definition)* ')'
    ;
delete_containers_clause
    : 'CONTAINERS' '(' (ID '.')? ID ')'
    ;
delete_shards_clause
    : 'SHARDS' '(' (ID '.')? ID ')'
    ;
delete_column_definition
    : ID DATA_TYPE ('COLLATE' ID)? 'SORT'? ('VISIBLE' | 'INVISIBLE')? ('DEFAULT' ('ON' 'NULL')? EXPR | delete_identity_clause)? ('ENCRYPT' delete_encryption_spec)? delete_inline_constraint* delete_inline_ref_constraint? delete_annotations_clause?
    ;
delete_inline_external_table_properties
    : 'TYPE' ID 'DEFAULT' 'DIRECTORY' ID ('ACCESS' 'PARAMETERS' 'USING' 'CLOB' QUERY)? 'LOCATION' '(' STRING ')' ('REJECT' 'LIMIT' (NUMBER | 'UNLIMITED'))?
    ;
delete_query_partition_clause
    : 'PARTITION' 'BY' (EXPR (',' EXPR)* | '(' EXPR (',' EXPR)* ')')
    ;
delete_outer_join_type
    : ('LEFT' | 'RIGHT' | 'FULL') 'OUTER'?
    ;
delete_sub_av_clause
    : 'ANALYTIC' 'VIEW' 'USING' ID delete_hierarchies_clause? delete_filter_clauses? delete_add_meas_clause?
    ;
delete_pivot_for_clause
    : 'FOR' (ID | '(' ID (',' ID)* ')')
    ;
delete_pivot_in_clause
    : 'IN' '(' ((EXPR | '(' EXPR (',' EXPR)* ')') ('AS'? ID)? (',' (EXPR | '(' EXPR (',' EXPR)* ')') ('AS'? ID)?)* | QUERY | 'ANY' (',' 'ANY')*) ')'
    ;
delete_unpivot_in_clause
    : 'IN' '(' (ID | '(' ID (',' ID)* ')') ('AS' (EXPR | '(' EXPR (',' EXPR)* ')'))? (',' (ID | '(' ID (',' ID)* ')') ('AS' (EXPR | '(' EXPR (',' EXPR)* ')'))?)* ')'
    ;
delete_row_pattern_measure_column
    : EXPR 'AS' ID
    ;
delete_row_pattern
    : delete_row_pattern_term ('\'|\'' delete_row_pattern_term)*
    ;
delete_subset_item
    : ID '=' '(' ID (',' ID)* ')'
    ;
delete_row_pattern_definition
    : ID 'AS' CONDITION
    ;
delete_identity_clause
    : 'GENERATED' ('ALWAYS' | 'BY' 'DEFAULT' ('ON' 'NULL')?) 'AS' 'IDENTITY' ('(' delete_identity_options ')')?
    ;
delete_encryption_spec
    : ('USING' STRING)? ('IDENTIFIED' 'BY' ID)? STRING? ('NO'? 'SALT')?
    ;
delete_inline_constraint
    : ('CONSTRAINT' ID)? ('NOT'? 'NULL' | 'UNIQUE' | 'PRIMARY' 'KEY' | 'REFERENCES' (ID '.')? ID ('(' ID ')')? ('ON' 'DELETE' ('CASCADE' | 'SET' 'NULL'))? | 'CHECK' '(' CONDITION ')') delete_constraint_state?
    ;
delete_inline_ref_constraint
    : ('SCOPE' 'IS' (ID '.')? ID | 'WITH' 'ROWID' | ('CONSTRAINT' ID)? 'REFERENCES' (ID '.')? ID ('(' ID (',' ID)* ')')? ('ON' 'DELETE' ('CASCADE' | 'SET' 'NULL'))? delete_constraint_state?)
    ;
delete_annotations_clause
    : 'ANNOTATIONS' '(' ('ADD' delete_annotation | 'DROP' ID (',' ID)* | 'REMOVE' 'ALL') ')'
    ;
delete_hierarchies_clause
    : 'HIERARCHIES' '(' ID (',' ID)* ')'
    ;
delete_filter_clauses
    : delete_filter_clause delete_filter_clause*
    ;
delete_add_meas_clause
    : 'ADD' 'MEASURES' '(' delete_calc_meas_clause (',' delete_calc_meas_clause)* ')'
    ;
delete_row_pattern_term
    : delete_row_pattern_factor delete_row_pattern_factor*
    ;
delete_identity_options
    : ('START' 'WITH' (NUMBER | 'LIMIT' 'VALUE') | 'INCREMENT' 'BY' NUMBER | 'MAXVALUE' NUMBER | 'NOMAXVALUE' | 'MINVALUE' NUMBER | 'NOMINVALUE' | 'CYCLE' | 'NOCYCLE' | 'CACHE' NUMBER | 'NOCACHE' | 'ORDER' | 'NOORDER')+
    ;
delete_constraint_state
    : ('NOT'? 'DEFERRABLE' ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? | 'INITIALLY' ('DEFERRED' | 'IMMEDIATE') ('NOT'? 'DEFERRABLE')?)? ('RELY' | 'NORELY')? delete_using_index_clause? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? delete_exceptions_clause?
    ;
delete_annotation
    : ID STRING?
    ;
delete_filter_clause
    : 'FILTER' 'FACT' '(' ID 'TO' (ID '.')? CONDITION ')'
    ;
delete_calc_meas_clause
    : 'AS' '(' EXPR ')'
    ;
delete_row_pattern_factor
    : delete_row_pattern_primary delete_row_pattern_quantifier?
    ;
delete_using_index_clause
    : 'USING' 'INDEX' ((ID '.')? ID | delete_index_attributes)
    ;
delete_exceptions_clause
    : 'EXCEPTIONS' 'INTO' (ID '.')? ID
    ;
delete_row_pattern_primary
    : ID
    | '\'$\''
    | '\'^\''
    | '(' delete_row_pattern? ')'
    | '\'{-\'' delete_row_pattern '\'-}\''
    | 'PERMUTE' '(' delete_row_pattern (',' delete_row_pattern)* ')'
    ;
delete_row_pattern_quantifier
    : ('\'*\'' | '\'+\'' | '\'?\'') '\'?\''?
    | '\'{\'' NUMBER ('\',\'' NUMBER?)? '\'}\'' '\'?\''?
    | '\'{\'' '\',\'' NUMBER '\'}\'' '\'?\''?
    ;
delete_index_attributes
    : (delete_physical_attributes_clause | 'TABLESPACE' (ID | 'DEFAULT') | delete_logging_clause | delete_index_compression | delete_partial_index_clause | delete_parallel_clause)+
    ;
delete_physical_attributes_clause
    : ('PCTUSED' NUMBER)? ('PCTFREE' NUMBER)? ('INITRANS' NUMBER)?
    ;
delete_logging_clause
    : ('LOGGING' | 'NOLOGGING')
    ;
delete_index_compression
    : (delete_prefix_compression | delete_advanced_index_compression | 'NOCOMPRESS')
    ;
delete_partial_index_clause
    : 'INDEXING' ('FULL' | 'PARTIAL')
    ;
delete_parallel_clause
    : ('PARALLEL' NUMBER? | 'NOPARALLEL')
    ;
delete_prefix_compression
    : 'COMPRESS' NUMBER?
    ;
delete_advanced_index_compression
    : 'COMPRESS' 'ADVANCED' ('LOW' | 'HIGH')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DISASSOCIATE-STATISTICS.html
disassociateStatistics
    : 'DISASSOCIATE' 'STATISTICS' 'FROM' ('COLUMNS' | 'FUNCTIONS' | 'PACKAGES' | 'TYPES' | 'INDEXES' | 'INDEXTYPES') (ID '.')? ID (',' (ID '.')? ID)* 'FORCE'?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-ANALYTIC-VIEW.html
dropAnalyticView
    : 'DROP' 'ANALYTIC' 'VIEW' ('IF' 'EXISTS')? (ID '.')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-ATTRIBUTE-DIMENSION.html
dropAttributeDimension
    : 'DROP' 'ATTRIBUTE' 'DIMENSION' ('IF' 'EXISTS')? (ID '.')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-AUDIT-POLICY-Unified-Auditing.html
dropAuditPolicyUnifiedAuditing
    : 'DROP' 'AUDIT' 'POLICY' ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-CLUSTER.html
dropCluster
    : 'DROP' 'CLUSTER' ('IF' 'EXISTS')? (ID '.')? ID ('INCLUDING' 'TABLES' ('CASCADE' 'CONSTRAINTS')?)?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-CONTEXT.html
dropContext
    : 'DROP' 'CONTEXT' ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-DATABASE-LINK.html
dropDatabaseLink
    : 'DROP' 'PUBLIC'? 'DATABASE' 'LINK' ('IF' 'EXISTS')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-DATABASE.html
dropDatabase
    : 'DROP' 'DATABASE'
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-DIMENSION.html
dropDimension
    : 'DROP' 'DIMENSION' (ID '.')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-DIRECTORY.html
dropDirectory
    : 'DROP' 'DIRECTORY' ('IF' 'EXISTS')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-DISKGROUP.html
dropDiskgroup
    : 'DROP' 'DISKGROUP' ID ('INCLUDING' 'CONTENTS' 'FORCE'? | 'EXCLUDING' 'CONTENTS')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-EDITION.html
dropEdition
    : 'DROP' 'EDITION' ('IF' 'EXISTS')? ID 'CASCADE'?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-FLASHBACK-ARCHIVE.html
dropFlashbackArchive
    : 'DROP' 'FLASHBACK' 'ARCHIVE' ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-FUNCTION.html
dropFunction
    : 'DROP' 'FUNCTION' ('IF' 'EXISTS')? (ID '.')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-HIERARCHY.html
dropHierarchy
    : 'DROP' 'HIERARCHY' ('IF' 'EXISTS')? (ID '.')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-INDEX.html
dropIndex
    : 'DROP' 'INDEX' ('IF' 'EXISTS')? (ID '.')? ID 'ONLINE'? 'FORCE'? (('DEFERRED' | 'IMMEDIATE') 'INVALIDATION')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-INDEXTYPE.html
dropIndextype
    : 'DROP' 'INDEXTYPE' ('IF' 'EXISTS')? (ID '.')? ID 'FORCE'?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-INMEMORY-JOIN-GROUP.html
dropInmemoryJoinGroup
    : 'DROP' 'INMEMORY' 'JOIN' 'GROUP' ('IF' 'EXISTS')? (ID '.')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-JAVA.html
dropJava
    : 'DROP' 'JAVA' ('SOURCE' | 'CLASS' | 'RESOURCE') ('IF' 'EXISTS')? (ID '.')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-LIBRARY.html
dropLibrary
    : 'DROP' 'LIBRARY' ('IF' 'EXISTS')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-LOCKDOWN-PROFILE.html
dropLockdownProfile
    : 'DROP' 'LOCKDOWN' 'PROFILE' ('IF' 'EXISTS')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-MATERIALIZED-VIEW-LOG.html
dropMaterializedViewLog
    : 'DROP' 'MATERIALIZED' 'VIEW' 'LOG' ('IF' 'EXISTS')? 'ON' (ID '.')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-MATERIALIZED-VIEW.html
dropMaterializedView
    : 'DROP' 'MATERIALIZED' 'VIEW' ('IF' 'EXISTS')? (ID '.')? ID ('PRESERVE' 'TABLE')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-MATERIALIZED-ZONEMAP.html
dropMaterializedZonemap
    : 'DROP' 'MATERIALIZED' 'ZONEMAP' ('IF' 'EXISTS')? (ID '.')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-OPERATOR.html
dropOperator
    : 'DROP' 'OPERATOR' ('IF' 'EXISTS')? (ID '.')? ID 'FORCE'?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-OUTLINE.html
dropOutline
    : 'DROP' 'OUTLINE' ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-PACKAGE.html
dropPackage
    : 'DROP' 'PACKAGE' 'BODY'? ('IF' 'EXISTS')? (ID '.')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-PLUGGABLE-DATABASE.html
dropPluggableDatabase
    : 'DROP' 'PLUGGABLE' 'DATABASE' ('IF' 'EXISTS')? ID 'FORCE'? ('KEEP' 'DATAFILES' | 'INCLUDING' 'DATAFILES')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-PROCEDURE.html
dropProcedure
    : 'DROP' 'PROCEDURE' ('IF' 'EXISTS')? (ID '.')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-PROFILE.html
dropProfile
    : 'DROP' 'PROFILE' ID 'CASCADE'?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-RESTORE-POINT.html
dropRestorePoint
    : 'DROP' 'RESTORE' 'POINT' ID ('FOR' 'PLUGGABLE' 'DATABASE' ID)?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-ROLE.html
dropRole
    : 'DROP' 'ROLE' ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-ROLLBACK-SEGMENT.html
dropRollbackSegment
    : 'DROP' 'ROLLBACK' 'SEGMENT' ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-SEQUENCE.html
dropSequence
    : 'DROP' 'SEQUENCE' ('IF' 'EXISTS')? (ID '.')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-SYNONYM.html
dropSynonym
    : 'DROP' 'PUBLIC'? 'SYNONYM' ('IF' 'EXISTS')? (ID '.')? ID 'FORCE'?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-TABLE.html
dropTable
    : 'DROP' 'TABLE' ('IF' 'EXISTS')? (ID '.')? ID ('CASCADE' 'CONSTRAINTS')? 'PURGE'?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-TABLESPACE-SET.html
dropTablespaceSet
    : 'DROP' 'TABLESPACE' 'SET' ID ('INCLUDING' 'CONTENTS' (('AND' 'DATAFILES' | 'KEEP' 'DATAFILES'))? ('CASCADE' 'CONSTRAINTS')?)?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-TABLESPACE.html
dropTablespace
    : 'DROP' 'TABLESPACE' ('IF' 'EXISTS')? ID (('DROP' | 'KEEP') 'QUOTA')? ('INCLUDING' 'CONTENTS' (('AND' 'DATAFILES' | 'KEEP' 'DATAFILES'))? ('CASCADE' 'CONSTRAINTS')?)?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-TRIGGER.html
dropTrigger
    : 'DROP' 'TRIGGER' ('IF' 'EXISTS')? (ID '.')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-TYPE-BODY.html
dropTypeBody
    : 'DROP' 'TYPE' 'BODY' ('IF' 'EXISTS')? (ID '.')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-TYPE.html
dropType
    : 'DROP' 'TYPE' ('IF' 'EXISTS')? (ID '.')? ID ('FORCE' | 'VALIDATE')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-USER.html
dropUser
    : 'DROP' 'USER' ('IF' 'EXISTS')? ID 'CASCADE'?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DROP-VIEW.html
dropView
    : 'DROP' 'VIEW' ('IF' 'EXISTS')? (ID '.')? ID ('CASCADE' 'CONSTRAINTS')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/EXPLAIN-PLAN.html
explainPlan
    : 'EXPLAIN' 'PLAN' ('SET' 'STATEMENT_ID' '=' STRING)? ('INTO' (ID '.')? ID ('@' ID)?)? 'FOR' ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/FLASHBACK-DATABASE.html
flashbackDatabase
    : 'FLASHBACK' ('STANDBY' | 'PLUGGABLE')? 'DATABASE' ID? ('TO' 'SCN' NUMBER | 'TO' 'BEFORE' 'SCN' NUMBER | 'TO' 'TIMESTAMP' EXPR | 'TO' 'BEFORE' 'TIMESTAMP' EXPR | 'TO' 'RESTORE' 'POINT' ID | 'TO' 'BEFORE' 'RESETLOGS')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/FLASHBACK-TABLE.html
flashbackTable
    : 'FLASHBACK' 'TABLE' (ID '.')? ID (',' (ID '.')? ID)* ('TO' 'SCN' EXPR | 'TO' 'TIMESTAMP' EXPR | 'TO' 'RESTORE' 'POINT' ID | 'TO' 'BEFORE' 'DROP' ('RENAME' 'TO' ID)?) (('ENABLE' | 'DISABLE') 'TRIGGERS')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/GRANT.html
grant
    : grant_grant_system_privileges
    | grant_grant_schema_privileges
    | grant_grant_object_privileges
    | grant_grant_roles_to_programs
    ;
grant_grant_system_privileges
    : 'GRANT' (ID (',' ID)* | 'ALL' 'PRIVILEGES') 'TO' grant_grantee_clause (',' grant_grantee_clause)* ('IDENTIFIED' 'BY' ID (',' ID)?)* ('WITH' ('ADMIN' | 'DELEGATE') 'OPTION')? ('CONTAINER' '=' ('CURRENT' | 'ALL'))?
    ;
grant_grant_schema_privileges
    : 'GRANT' (ID (',' ID)* | 'ALL' 'PRIVILEGES'?) 'ON' 'SCHEMA' ID 'TO' grant_grantee_clause (',' grant_grantee_clause)* ('WITH' 'ADMIN' 'OPTION')? ('CONTAINER' '=' ('CURRENT' | 'ALL'))?
    ;
grant_grant_object_privileges
    : 'GRANT' (ID ('(' ID (',' ID)* ')')? (',' ID ('(' ID (',' ID)* ')')?)* | 'ALL' 'PRIVILEGES'?) grant_on_object_clause 'TO' grant_grantee_clause (',' grant_grantee_clause)* ('WITH' ('GRANT' | 'HIERARCHY') 'OPTION')? ('CONTAINER' '=' ('CURRENT' | 'ALL'))?
    ;
grant_grant_roles_to_programs
    : 'GRANT' ID (',' ID)* 'TO' grant_program_unit (',' grant_program_unit)* ('WITH' ('ADMIN' | 'DELEGATE') 'OPTION')? ('CONTAINER' '=' ('CURRENT' | 'ALL'))?
    ;
grant_grantee_clause
    : (ID | 'PUBLIC')
    ;
grant_on_object_clause
    : 'ON' (ID '.')? ID
    | 'ON' 'USER' ID (',' ID)*
    | 'ON' 'DIRECTORY' ID
    | 'ON' 'EDITION' ID
    | 'ON' 'MINING' 'MODEL' (ID '.')? ID
    | 'ON' 'JAVA' ('SOURCE' | 'RESOURCE') (ID '.')? ID
    | 'ON' 'SQL' 'TRANSLATION' 'PROFILE' (ID '.')? ID
    ;
grant_program_unit
    : (ID '.')? ('FUNCTION' ID | 'PROCEDURE' ID | 'PACKAGE' ID)
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/INSERT.html
insert
    : 'INSERT' HINT? (insert_single_table_insert | insert_multi_table_insert)
    ;
insert_single_table_insert
    : insert_insert_into_clause (insert_values_clause | QUERY | insert_insert_set_clause | insert_by_name_position_clause) insert_returning_clause? insert_error_logging_clause?
    ;
insert_multi_table_insert
    : ('ALL' insert_into_clause insert_values_clause? (',' insert_into_clause insert_values_clause?)* | insert_conditional_insert_clause) QUERY
    ;
insert_insert_into_clause
    : 'INTO' (ID '.')? ID insert_partition_extension_clause? ('@' ID)? ID? ('(' ID (',' ID)* ')')?
    ;
insert_values_clause
    : 'VALUES' '(' EXPR (',' EXPR)* ')'
    ;
insert_insert_set_clause
    : 'SET' ID '=' (EXPR | 'DEFAULT') (',' ID '=' (EXPR | 'DEFAULT'))*
    ;
insert_by_name_position_clause
    : 'BY' ('NAME' | 'POSITION') QUERY
    ;
insert_returning_clause
    : ('RETURNING' | 'RETURN') EXPR (',' EXPR)* 'INTO' ID (',' ID)*
    ;
insert_error_logging_clause
    : 'LOG' 'ERRORS' ('INTO' (ID '.')? ID)? ('(' EXPR ')')? ('REJECT' 'LIMIT' (NUMBER | 'UNLIMITED'))?
    ;
insert_into_clause
    : 'INTO' (ID '.')? ID ('(' ID (',' ID)* ')')?
    ;
insert_conditional_insert_clause
    : ('ALL' | 'FIRST')? 'WHEN' CONDITION 'THEN' insert_into_clause insert_values_clause? (',' insert_into_clause insert_values_clause?)* ('WHEN' CONDITION 'THEN' insert_into_clause insert_values_clause? (',' insert_into_clause insert_values_clause?)?)* ('ELSE' insert_into_clause insert_values_clause? (',' insert_into_clause insert_values_clause?)?)*
    ;
insert_partition_extension_clause
    : 'PARTITION' '(' ID ')'
    | 'PARTITION' 'FOR' '(' EXPR (',' EXPR)* ')'
    | 'SUBPARTITION' '(' ID ')'
    | 'SUBPARTITION' 'FOR' '(' EXPR (',' EXPR)* ')'
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/LOCK-TABLE.html
lockTable
    : 'LOCK' 'TABLE' (ID '.')? ID ('@' ID)? lockTable_partition_extension_clause? (',' (ID '.')? ID ('@' ID)? lockTable_partition_extension_clause?)* 'IN' lockTable_lockmode 'MODE' ('NOWAIT' | 'WAIT' NUMBER)?
    ;
lockTable_partition_extension_clause
    : 'PARTITION' '(' ID ')'
    | 'PARTITION' 'FOR' '(' EXPR ')'
    | 'SUBPARTITION' '(' ID ')'
    | 'SUBPARTITION' 'FOR' '(' EXPR ')'
    ;
lockTable_lockmode
    : 'ROW' 'SHARE'
    | 'ROW' 'EXCLUSIVE'
    | 'SHARE' 'UPDATE'
    | 'SHARE'
    | 'SHARE' 'ROW' 'EXCLUSIVE'
    | 'EXCLUSIVE'
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/MERGE.html
merge
    : 'MERGE' HINT? 'INTO' (ID '.')? ID ID? 'USING' ((ID '.')? ID | QUERY) ID? 'ON' '(' CONDITION ')' merge_merge_update_clause? merge_merge_insert_clause? merge_error_logging_clause?
    ;
merge_merge_update_clause
    : 'WHEN' 'MATCHED' 'THEN' 'UPDATE' 'SET' ID '=' (EXPR | 'DEFAULT') (',' ID '=' (EXPR | 'DEFAULT'))* ('WHERE' CONDITION)? ('DELETE' 'WHERE' CONDITION)?
    ;
merge_merge_insert_clause
    : 'WHEN' 'NOT' 'MATCHED' 'THEN' 'INSERT' ('(' ID (',' ID)* ')')? 'VALUES' '(' (EXPR | 'DEFAULT') (',' (EXPR | 'DEFAULT'))* ')' ('WHERE' CONDITION)?
    ;
merge_error_logging_clause
    : 'LOG' 'ERRORS' ('INTO' (ID '.')? ID)? ('(' EXPR ')')? ('REJECT' 'LIMIT' (NUMBER | 'UNLIMITED'))?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/NOAUDIT-Traditional-Auditing.html
noauditTraditionalAuditing
    : 'NOAUDIT' (noauditTraditionalAuditing_audit_operation_clause noauditTraditionalAuditing_auditing_by_clause? | noauditTraditionalAuditing_audit_schema_object_clause) ('WHENEVER' 'NOT'? 'SUCCESSFUL')? ('CONTAINER' '=' ('CURRENT' | 'ALL'))?
    ;
noauditTraditionalAuditing_audit_operation_clause
    : (ID (',' ID)* | 'ALL' | 'ALL' 'STATEMENTS' | 'ALL' 'PRIVILEGES')
    ;
noauditTraditionalAuditing_auditing_by_clause
    : 'BY' ID (',' ID)*
    ;
noauditTraditionalAuditing_audit_schema_object_clause
    : (ID (',' ID)* | 'ALL') 'ON' noauditTraditionalAuditing_auditing_on_clause
    ;
noauditTraditionalAuditing_auditing_on_clause
    : (ID '.')? ID
    | 'DIRECTORY' ID
    | 'SQL' 'TRANSLATION' 'PROFILE' (ID '.')? ID
    | 'DEFAULT'
    | 'NETWORK'
    | 'DIRECT_PATH' 'LOAD'
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/NOAUDIT-Unified-Auditing.html
noauditUnifiedAuditing
    : 'NOAUDIT' 'POLICY' ID ('BY' ID (',' ID)* ('WITH' 'ROLE' ID (',' ID)?)*)?
    | 'NOAUDIT' 'CONTEXT' 'NAMESPACE' ID 'ATTRIBUTES' ID (',' ID)* ('BY' ID (',' ID)* ('WITH' 'ROLE' ID (',' ID)?)*)?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/PURGE.html
purge
    : 'PURGE' ('TABLE' (ID '.')? ID | 'INDEX' (ID '.')? ID | 'TABLESPACE' ID ('USER' ID)? | 'TABLESPACE' 'SET' ID ('USER' ID)? | 'RECYCLEBIN' | 'DBA_RECYCLEBIN')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/RENAME.html
rename
    : 'RENAME' ID 'TO' ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/REVOKE.html
revoke
    : 'REVOKE' (revoke_revoke_system_privileges | revoke_revoke_schema_privileges | revoke_revoke_object_privileges | revoke_revoke_roles_from_programs) ('CONTAINER' '=' ('CURRENT' | 'ALL'))?
    ;
revoke_revoke_system_privileges
    : 'REVOKE' (ID (',' ID)* | 'ALL' 'PRIVILEGES') 'FROM' revoke_revokee_clause
    ;
revoke_revoke_schema_privileges
    : 'REVOKE' (ID (',' ID)* | 'ALL' 'PRIVILEGES'?) 'ON' 'SCHEMA' ID 'FROM' revoke_revokee_clause
    ;
revoke_revoke_object_privileges
    : 'REVOKE' (ID (',' ID)* | 'ALL' 'PRIVILEGES'?) revoke_on_object_clause 'FROM' revoke_revokee_clause ('CASCADE' 'CONSTRAINTS')? 'FORCE'?
    ;
revoke_revoke_roles_from_programs
    : 'REVOKE' (ID (',' ID)* | 'ALL') 'FROM' ('PROCEDURE' | 'FUNCTION' | 'PACKAGE') (ID '.')? revoke_program_unit
    ;
revoke_revokee_clause
    : (ID | 'PUBLIC') (',' (ID | 'PUBLIC'))*
    ;
revoke_on_object_clause
    : 'ON' ((ID '.')? ID | 'USER' ID (',' ID)* | 'DIRECTORY' ID | 'EDITION' ID | 'MINING' 'MODEL' (ID '.')? ID | 'JAVA' ('SOURCE' | 'RESOURCE') (ID '.')? ID | 'SQL' 'TRANSLATION' 'PROFILE' (ID '.')? ID)
    ;
revoke_program_unit
    : (ID '.')? ('FUNCTION' ID | 'PROCEDURE' ID | 'PACKAGE' ID)
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ROLLBACK.html
rollback
    : 'ROLLBACK' 'WORK'? ('TO' 'SAVEPOINT'? ID | 'FORCE' STRING)?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/SAVEPOINT.html
savepoint
    : 'SAVEPOINT' ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/SELECT.html
// fix: left out where optional, since no quoted page defines it or what it needs: create_index_statement, opaque_format_spec, storage_clause
select
    : select_with_clause? select_subquery select_for_update_clause?
    ;
select_with_clause
    : 'WITH' (select_plsql_declarations | select_subquery_factoring_clause | select_subav_factoring_clause)+
    ;
select_subquery
    : select_query_block (('UNION' 'ALL'? | 'INTERSECT' | 'MINUS') select_query_block)* select_order_by_clause? select_row_limiting_clause?
    ;
select_for_update_clause
    : 'FOR' 'UPDATE' ('OF' ((ID '.')? ID '.')? ID (',' ((ID '.')? ID '.')? ID)?)* ('NOWAIT' | 'WAIT' NUMBER | 'SKIP' 'LOCKED')?
    ;
select_plsql_declarations
    : (FUNCTION_DECLARATION | PROCEDURE_DECLARATION)+
    ;
select_subquery_factoring_clause
    : ID ('(' ID (',' ID)* ')')? 'AS' '(' select_subquery ')' select_search_clause? select_cycle_clause? (',' ID ('(' ID (',' ID)* ')')? 'AS' '(' select_subquery ')' select_search_clause? select_cycle_clause?)*
    ;
select_subav_factoring_clause
    : ID 'AS' '(' select_sub_av_clause ')' (',' ID 'AS' '(' select_sub_av_clause ')')*
    ;
select_query_block
    : 'SELECT' HINT? ('DISTINCT' | 'UNIQUE' | 'ALL')? select_select_list ('FROM' (select_table_reference | select_join_clause) (',' (select_table_reference | select_join_clause))?)* select_where_clause? select_hierarchical_query_clause? select_group_by_clause? select_model_clause? select_window_clause? select_qualify_clause?
    ;
select_order_by_clause
    : 'ORDER' 'SIBLINGS'? 'BY' (EXPR | NUMBER | ID) ('ASC' | 'DESC')? ('NULLS' 'FIRST' | 'NULLS' 'LAST')? (',' (EXPR | NUMBER | ID) ('ASC' | 'DESC')? ('NULLS' 'FIRST' | 'NULLS' 'LAST')?)*
    ;
select_row_limiting_clause
    : ('OFFSET' NUMBER ('ROW' | 'ROWS'))? ('FETCH' ('FIRST' | 'NEXT') ((NUMBER | NUMBER 'PERCENT'))? ('ROW' | 'ROWS') ('ONLY' | 'WITH' 'TIES'))?
    ;
select_search_clause
    : 'SEARCH' ('BREADTH' 'FIRST' | 'DEPTH' 'FIRST') 'BY' ID ('ASC' | 'DESC')? ('NULLS' 'FIRST' | 'NULLS' 'LAST')? (',' ID ('ASC' | 'DESC')? ('NULLS' 'FIRST' | 'NULLS' 'LAST')?)* 'SET' ID
    ;
select_cycle_clause
    : 'CYCLE' ID (',' ID)* 'SET' ID 'TO' EXPR 'DEFAULT' EXPR
    ;
select_sub_av_clause
    : 'ANALYTIC' 'VIEW' 'USING' ID select_hierarchies_clause? select_filter_clauses? select_add_meas_clause?
    ;
select_select_list
    : ('*' | (ID | (ID '.')? ID) '.' '*' | EXPR ('AS'? ID)?) (',' ('*' | (ID | (ID '.')? ID) '.' '*' | EXPR ('AS'? ID)?))*
    ;
select_table_reference
    : (select_query_table_expression | '(' select_join_clause ')' | select_inline_external_table) ID?
    ;
select_join_clause
    : select_table_reference (select_inner_cross_join_clause | select_outer_join_clause | select_cross_outer_apply_clause)+
    ;
select_where_clause
    : 'WHERE' CONDITION
    ;
select_hierarchical_query_clause
    : ('START' 'WITH' CONDITION 'CONNECT' 'BY' 'NOCYCLE'? CONDITION | 'CONNECT' 'BY' 'NOCYCLE'? CONDITION ('START' 'WITH' CONDITION)?)
    ;
select_group_by_clause
    : 'GROUP' 'BY' (EXPR | select_rollup_cube_clause | select_grouping_sets_clause) (',' (EXPR | select_rollup_cube_clause | select_grouping_sets_clause))* ('HAVING' CONDITION)?
    ;
select_model_clause
    : 'MODEL' select_cell_reference_options? select_return_rows_clause? select_reference_model* select_main_model
    ;
select_window_clause
    : 'WINDOW' ID 'AS' '(' select_window_specification ')' (',' ID 'AS' '(' select_window_specification ')')*
    ;
select_qualify_clause
    : 'QUALIFY' CONDITION
    ;
select_hierarchies_clause
    : 'HIERARCHIES' '(' ID (',' ID)* ')'
    ;
select_filter_clauses
    : select_filter_clause select_filter_clause*
    ;
select_add_meas_clause
    : 'ADD' 'MEASURES' '(' select_calc_meas_clause (',' select_calc_meas_clause)* ')'
    ;
select_query_table_expression
    : ((ID '.')? ID select_partition_extension_clause? ('@' ID)? | '(' select_subquery select_subquery_restriction_clause? ')' | select_table_collection_expression | select_inline_analytic_view) select_flashback_query_clause? select_sample_clause? (select_pivot_clause | select_unpivot_clause)? select_row_pattern_clause? select_containers_clause? select_shards_clause?
    ;
select_inline_external_table
    : 'EXTERNAL' '(' '(' select_column_definition (',' select_column_definition)* ')' select_inline_external_table_properties ')'
    ;
select_inner_cross_join_clause
    : ('INNER'? 'JOIN' select_table_reference ('ON' CONDITION | 'USING' '(' ID (',' ID)* ')') | 'CROSS' 'JOIN' select_table_reference | 'NATURAL' 'INNER'? 'JOIN' select_table_reference)
    ;
select_outer_join_clause
    : (select_query_partition_clause select_outer_join_type 'JOIN' select_table_reference ('ON' CONDITION | 'USING' '(' ID (',' ID)* ')') | select_outer_join_type 'JOIN' select_query_partition_clause select_table_reference ('ON' CONDITION | 'USING' '(' ID (',' ID)* ')') | 'NATURAL' select_outer_join_type 'JOIN' select_table_reference)
    ;
select_cross_outer_apply_clause
    : ('CROSS' | 'OUTER') 'APPLY' (select_table_reference | EXPR)
    ;
select_rollup_cube_clause
    : ('ROLLUP' | 'CUBE') '(' select_grouping_expression_list ')'
    ;
select_grouping_sets_clause
    : 'GROUPING' 'SETS' '(' (select_rollup_cube_clause | select_grouping_expression_list) (',' (select_rollup_cube_clause | select_grouping_expression_list))* ')'
    ;
select_cell_reference_options
    : ('IGNORE' 'NAV' | 'KEEP' 'NAV') 'UNIQUE' ('DIMENSION' | 'SINGLE' 'REFERENCE')
    ;
select_return_rows_clause
    : 'RETURN' ('ALL' | 'UPDATED') 'ROWS'
    ;
select_reference_model
    : 'REFERENCE' ID 'ON' '(' select_subquery ')' select_model_column_clauses select_cell_reference_options?
    ;
select_main_model
    : ('MAIN' ID)? select_model_column_clauses select_cell_reference_options? select_model_rules_clause
    ;
select_window_specification
    : ID? select_query_partition_clause? select_order_by_clause? select_windowing_clause?
    ;
select_filter_clause
    : 'FILTER' 'FACT' '(' ID 'TO' (ID '.')? CONDITION ')'
    ;
select_calc_meas_clause
    : ID 'AS' '(' EXPR ')'
    ;
select_partition_extension_clause
    : ('PARTITION' '(' ID ')' | 'PARTITION' 'FOR' '(' EXPR (',' EXPR)* ')' | 'SUBPARTITION' '(' ID ')' | 'SUBPARTITION' 'FOR' '(' EXPR (',' EXPR)* ')')
    ;
select_subquery_restriction_clause
    : 'WITH' ('READ' 'ONLY' | 'CHECK' 'OPTION' ('CONSTRAINT' ID)?)
    ;
select_table_collection_expression
    : 'TABLE' '(' EXPR ')' ('(' '+' ')')?
    ;
select_inline_analytic_view
    : 'ANALYTIC' 'VIEW' '(' select_sub_av_clause ')' ID?
    ;
select_flashback_query_clause
    : ('VERSIONS' 'BETWEEN' ('SCN' | 'TIMESTAMP') EXPR 'AND' EXPR | 'VERSIONS' 'PERIOD' 'FOR' ID 'BETWEEN' EXPR 'AND' EXPR | 'AS' 'OF' ('SCN' | 'TIMESTAMP') EXPR | 'AS' 'OF' 'PERIOD' 'FOR' ID EXPR)
    ;
select_sample_clause
    : 'SAMPLE' 'BLOCK'? '(' NUMBER ')' ('SEED' '(' NUMBER ')')?
    ;
select_pivot_clause
    : 'PIVOT' 'XML'? '(' ID '(' EXPR ')' ('AS'? ID)? (',' ID '(' EXPR ')' ('AS'? ID)?)* select_pivot_for_clause select_pivot_in_clause ')'
    ;
select_unpivot_clause
    : 'UNPIVOT' (('INCLUDE' | 'EXCLUDE') 'NULLS')? '(' (ID | '(' ID (',' ID)* ')') select_pivot_for_clause select_unpivot_in_clause ')'
    ;
select_row_pattern_clause
    : 'MATCH_RECOGNIZE' '(' ('PARTITION' 'BY' ID (',' ID)?)* ('ORDER' 'BY' ID ('ASC' | 'DESC')? (',' ID ('ASC' | 'DESC')?)?)* ('MEASURES' select_row_pattern_measure_column (',' select_row_pattern_measure_column)?)* ('ONE' 'ROW' 'PER' 'MATCH' | 'ALL' 'ROWS' 'PER' 'MATCH' (('SHOW' | 'OMIT') 'EMPTY' 'MATCHES')?)? ('AFTER' 'MATCH' 'SKIP' ('PAST' 'LAST' 'ROW' | 'TO' 'NEXT' 'ROW' | 'TO' 'FIRST' ID | 'TO' 'LAST' ID | 'TO' ID))? 'PATTERN' '(' select_row_pattern ')' ('SUBSET' select_subset_item (',' select_subset_item)?)* 'DEFINE' select_row_pattern_definition (',' select_row_pattern_definition)* ')'
    ;
select_containers_clause
    : 'CONTAINERS' '(' (ID '.')? ID ')'
    ;
select_shards_clause
    : 'SHARDS' '(' (ID '.')? ID ')'
    ;
select_column_definition
    : ID DATA_TYPE ('COLLATE' ID)? 'SORT'? ('VISIBLE' | 'INVISIBLE')? ('DEFAULT' ('ON' 'NULL')? EXPR | select_identity_clause)? ('ENCRYPT' select_encryption_spec)? select_inline_constraint* select_inline_ref_constraint? select_annotations_clause?
    ;
select_inline_external_table_properties
    : 'TYPE' ID 'DEFAULT' 'DIRECTORY' ID ('ACCESS' 'PARAMETERS' 'USING' 'CLOB' select_subquery)? 'LOCATION' '(' STRING ')' ('REJECT' 'LIMIT' (NUMBER | 'UNLIMITED'))?
    ;
select_query_partition_clause
    : 'PARTITION' 'BY' (EXPR (',' EXPR)* | '(' EXPR (',' EXPR)* ')')
    ;
select_outer_join_type
    : ('LEFT' | 'RIGHT' | 'FULL') 'OUTER'?
    ;
select_grouping_expression_list
    : select_expression_list (',' select_expression_list)*
    ;
select_model_column_clauses
    : select_query_partition_clause? 'DIMENSION' 'BY' '(' select_model_column (',' select_model_column)* ')' 'MEASURES' '(' select_model_column (',' select_model_column)* ')'
    ;
select_model_rules_clause
    : ('RULES' (('UPDATE' | 'UPSERT' 'ALL'?))? (('AUTOMATIC' | 'SEQUENTIAL') 'ORDER')? select_model_iterate_clause?)? '(' (('UPDATE' | 'UPSERT' 'ALL'?))? select_cell_assignment select_order_by_clause? '=' EXPR (',' (('UPDATE' | 'UPSERT' 'ALL'?))? select_cell_assignment select_order_by_clause? '=' EXPR)* ')'
    ;
select_windowing_clause
    : ('ROWS' | 'RANGE' | 'GROUPS') ('BETWEEN' ('UNBOUNDED' 'PRECEDING' | 'CURRENT' 'ROW' | EXPR ('PRECEDING' | 'FOLLOWING')) 'AND' ('UNBOUNDED' 'FOLLOWING' | 'CURRENT' 'ROW' | EXPR ('PRECEDING' | 'FOLLOWING')) | ('UNBOUNDED' 'PRECEDING' | 'CURRENT' 'ROW' | EXPR 'PRECEDING')) (('EXCLUDE' 'CURRENT' 'ROW' | 'EXCLUDE' 'GROUP' | 'EXCLUDE' 'TIES' | 'EXCLUDE' 'NO' 'OTHERS'))?
    ;
select_pivot_for_clause
    : 'FOR' (ID | '(' ID (',' ID)* ')')
    ;
select_pivot_in_clause
    : 'IN' '(' ((EXPR | '(' EXPR (',' EXPR)* ')') ('AS'? ID)? (',' (EXPR | '(' EXPR (',' EXPR)* ')') ('AS'? ID)?)* | select_subquery | 'ANY' (',' 'ANY')*) ')'
    ;
select_unpivot_in_clause
    : 'IN' '(' (ID | '(' ID (',' ID)* ')') ('AS' (EXPR | '(' EXPR (',' EXPR)* ')'))? (',' (ID | '(' ID (',' ID)* ')') ('AS' (EXPR | '(' EXPR (',' EXPR)* ')'))?)* ')'
    ;
select_row_pattern_measure_column
    : EXPR 'AS' ID
    ;
select_row_pattern
    : select_row_pattern_term ('\'|\'' select_row_pattern_term)*
    ;
select_subset_item
    : ID '=' '(' ID (',' ID)* ')'
    ;
select_row_pattern_definition
    : ID 'AS' CONDITION
    ;
select_identity_clause
    : 'GENERATED' ('ALWAYS' | 'BY' 'DEFAULT' ('ON' 'NULL')?) 'AS' 'IDENTITY' ('(' select_identity_options ')')?
    ;
select_encryption_spec
    : ('USING' STRING)? ('IDENTIFIED' 'BY' ID)? STRING? ('NO'? 'SALT')?
    ;
select_inline_constraint
    : ('CONSTRAINT' ID)? ('NOT'? 'NULL' | 'UNIQUE' | 'PRIMARY' 'KEY' | 'REFERENCES' (ID '.')? ID ('(' ID ')')? ('ON' 'DELETE' ('CASCADE' | 'SET' 'NULL'))? | 'CHECK' '(' CONDITION ')') select_constraint_state?
    ;
select_inline_ref_constraint
    : ('SCOPE' 'IS' (ID '.')? ID | 'WITH' 'ROWID' | ('CONSTRAINT' ID)? 'REFERENCES' (ID '.')? ID ('(' ID (',' ID)* ')')? ('ON' 'DELETE' ('CASCADE' | 'SET' 'NULL'))? select_constraint_state?)
    ;
select_annotations_clause
    : 'ANNOTATIONS' '(' ('ADD' select_annotation | 'DROP' ID (',' ID)* | 'REMOVE' 'ALL') ')'
    ;
select_expression_list
    : (EXPR (',' EXPR)* | '(' (EXPR (',' EXPR)?)* ')')
    ;
select_model_column
    : EXPR ('AS'? ID)?
    ;
select_model_iterate_clause
    : 'ITERATE' '(' NUMBER ')' ('UNTIL' '(' CONDITION ')')?
    ;
select_cell_assignment
    : ID ((CONDITION | EXPR | select_single_column_for_loop | select_multi_column_for_loop) (',' (CONDITION | EXPR | select_single_column_for_loop | select_multi_column_for_loop))?)*
    ;
select_row_pattern_term
    : select_row_pattern_factor select_row_pattern_factor*
    ;
select_identity_options
    : ('START' 'WITH' (NUMBER | 'LIMIT' 'VALUE') | 'INCREMENT' 'BY' NUMBER | 'MAXVALUE' NUMBER | 'NOMAXVALUE' | 'MINVALUE' NUMBER | 'NOMINVALUE' | 'CYCLE' | 'NOCYCLE' | 'CACHE' NUMBER | 'NOCACHE' | 'ORDER' | 'NOORDER')+
    ;
select_constraint_state
    : ('NOT'? 'DEFERRABLE' ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? | 'INITIALLY' ('DEFERRED' | 'IMMEDIATE') ('NOT'? 'DEFERRABLE')?)? ('RELY' | 'NORELY')? select_using_index_clause? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? select_exceptions_clause?
    ;
select_annotation
    : ID STRING?
    ;
select_single_column_for_loop
    : 'FOR' ID ('IN' '(' (EXPR (',' EXPR)* | select_subquery) ')' | ('LIKE' ID)? 'FROM' EXPR 'TO' EXPR ('INCREMENT' | 'DECREMENT') EXPR)
    ;
select_multi_column_for_loop
    : 'FOR' '(' ID (',' ID)* ')' 'IN' '(' ('(' EXPR (',' EXPR)* ')' (',' '(' EXPR (',' EXPR)* ')')* | select_subquery) ')'
    ;
select_row_pattern_factor
    : select_row_pattern_primary select_row_pattern_quantifier?
    ;
select_using_index_clause
    : 'USING' 'INDEX' ((ID '.')? ID | select_index_attributes)
    ;
select_exceptions_clause
    : 'EXCEPTIONS' 'INTO' (ID '.')? ID
    ;
select_row_pattern_primary
    : ID
    | '\'$\''
    | '\'^\''
    | '(' select_row_pattern? ')'
    | '\'{-\'' select_row_pattern '\'-}\''
    | 'PERMUTE' '(' select_row_pattern (',' select_row_pattern)* ')'
    ;
select_row_pattern_quantifier
    : ('\'*\'' | '\'+\'' | '\'?\'') '\'?\''?
    | '\'{\'' NUMBER ('\',\'' NUMBER?)? '\'}\'' '\'?\''?
    | '\'{\'' '\',\'' NUMBER '\'}\'' '\'?\''?
    ;
select_index_attributes
    : (select_physical_attributes_clause | 'TABLESPACE' (ID | 'DEFAULT') | select_logging_clause | select_index_compression | select_partial_index_clause | select_parallel_clause)+
    ;
select_physical_attributes_clause
    : ('PCTUSED' NUMBER)? ('PCTFREE' NUMBER)? ('INITRANS' NUMBER)?
    ;
select_logging_clause
    : ('LOGGING' | 'NOLOGGING')
    ;
select_index_compression
    : (select_prefix_compression | select_advanced_index_compression | 'NOCOMPRESS')
    ;
select_partial_index_clause
    : 'INDEXING' ('FULL' | 'PARTIAL')
    ;
select_parallel_clause
    : ('PARALLEL' NUMBER? | 'NOPARALLEL')
    ;
select_prefix_compression
    : 'COMPRESS' NUMBER?
    ;
select_advanced_index_compression
    : 'COMPRESS' 'ADVANCED' ('LOW' | 'HIGH')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/SET-CONSTRAINTS.html
setConstraints
    : 'SET' ('CONSTRAINT' | 'CONSTRAINTS') (ID (',' ID)* | 'ALL') ('IMMEDIATE' | 'DEFERRED')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/SET-ROLE.html
setRole
    : 'SET' 'ROLE' (ID ('IDENTIFIED' 'BY' ID)? (',' ID ('IDENTIFIED' 'BY' ID)?)* | 'ALL' ('EXCEPT' ID (',' ID)?)* | 'NONE')
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/SET-TRANSACTION.html
setTransaction
    : 'SET' 'TRANSACTION' (('READ' 'ONLY' | 'READ' 'WRITE') | 'ISOLATION' 'LEVEL' ('SERIALIZABLE' | 'READ' 'COMMITTED') | 'USE' 'ROLLBACK' 'SEGMENT' ID | 'NAME' STRING)+
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/TRUNCATE-CLUSTER.html
truncateCluster
    : 'TRUNCATE' 'CLUSTER' (ID '.')? ID (('DROP' 'STORAGE' | 'REUSE' 'STORAGE'))?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/TRUNCATE-TABLE.html
truncateTable
    : 'TRUNCATE' 'TABLE' (ID '.')? ID (('PRESERVE' | 'PURGE') 'MATERIALIZED' 'VIEW' 'LOG')? (('DROP' 'STORAGE' | 'DROP' 'ALL' 'STORAGE' | 'REUSE' 'STORAGE'))? 'CASCADE'?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/UPDATE.html
update
    : 'UPDATE' HINT? ((ID '.')? ID ID? update_partition_extension_clause? ('@' ID)? | '(' QUERY update_subquery_restriction_clause? ')' | update_table_collection_expression) update_update_set_clause update_from_using_clause? update_where_clause? update_returning_clause? update_error_logging_clause?
    ;
update_partition_extension_clause
    : ('PARTITION' '(' ID ')' | 'PARTITION' 'FOR' '(' EXPR (',' EXPR)* ')' | 'SUBPARTITION' '(' ID ')' | 'SUBPARTITION' 'FOR' '(' EXPR (',' EXPR)* ')')
    ;
update_subquery_restriction_clause
    : 'WITH' ('READ' 'ONLY' | 'CHECK' 'OPTION' ('CONSTRAINT' ID)?)
    ;
update_table_collection_expression
    : 'TABLE' '(' EXPR ')' ('(' '+' ')')?
    ;
update_update_set_clause
    : 'SET' ((ID '=' (EXPR | '(' QUERY ')' | 'DEFAULT') | '(' ID (',' ID)* ')' '=' '(' QUERY ')') (',' (ID '=' (EXPR | '(' QUERY ')' | 'DEFAULT') | '(' ID (',' ID)* ')' '=' '(' QUERY ')'))* | 'VALUE' '(' ID ')' '=' (EXPR | '(' QUERY ')'))
    ;
update_from_using_clause
    : 'FROM' (update_table_reference | update_join_clause) (',' (update_table_reference | update_join_clause))*
    ;
update_where_clause
    : 'WHERE' CONDITION
    ;
update_returning_clause
    : ('RETURN' | 'RETURNING') (EXPR | 'OLD' ID | 'NEW' ID) (',' (EXPR | 'OLD' ID | 'NEW' ID))* 'INTO' ID (',' ID)*
    ;
update_error_logging_clause
    : 'LOG' 'ERRORS' ('INTO' (ID '.')? ID)? ('(' EXPR ')')? ('REJECT' 'LIMIT' (NUMBER | 'UNLIMITED'))?
    ;
update_table_reference
    : ((ID '.')? ID | '(' QUERY ')') ID?
    ;
update_join_clause
    : update_table_reference (update_inner_cross_join_clause | update_outer_join_clause | update_cross_outer_apply_clause)+
    ;
update_inner_cross_join_clause
    : ('INNER'? 'JOIN' update_table_reference ('ON' CONDITION | 'USING' '(' ID (',' ID)* ')') | 'CROSS' 'JOIN' update_table_reference | 'NATURAL' 'INNER'? 'JOIN' update_table_reference)
    ;
update_outer_join_clause
    : (update_query_partition_clause update_outer_join_type 'JOIN' update_table_reference ('ON' CONDITION | 'USING' '(' ID (',' ID)* ')') | update_outer_join_type 'JOIN' update_query_partition_clause update_table_reference ('ON' CONDITION | 'USING' '(' ID (',' ID)* ')') | 'NATURAL' update_outer_join_type 'JOIN' update_table_reference)
    ;
update_cross_outer_apply_clause
    : ('CROSS' | 'OUTER') 'APPLY' (update_table_reference | EXPR)
    ;
update_query_partition_clause
    : 'PARTITION' 'BY' (EXPR (',' EXPR)* | '(' EXPR (',' EXPR)* ')')
    ;
update_outer_join_type
    : ('LEFT' | 'RIGHT' | 'FULL') 'OUTER'?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/alter-domain.html
alterDomain
    : 'ALTER' 'USECASE'? 'DOMAIN' ('IF' 'EXISTS')? (ID '.')? ID ('ADD' 'DISPLAY' EXPR | 'MODIFY' 'DISPLAY' EXPR | 'DROP' 'DISPLAY' | 'ADD' 'ORDER' EXPR | 'MODIFY' 'ORDER' EXPR | 'DROP' 'ORDER' | alterDomain_annotations_clause)
    ;
alterDomain_annotations_clause
    : 'ANNOTATIONS' '(' ('ADD' alterDomain_annotation | 'DROP' ID (',' ID)* | 'REMOVE' 'ALL') ')'
    ;
alterDomain_annotation
    : ID STRING?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/alter-json-relational-duality-view.html
alterJsonRelationalDualityView
    : 'ALTER' 'JSON' 'RELATIONAL' 'DUALITY' 'VIEW' (ID '.')? ID alterJsonRelationalDualityView_duality_view_replication_clause
    ;
alterJsonRelationalDualityView_duality_view_replication_clause
    : 'ENABLE' 'LOGICAL' 'REPLICATION'
    | 'DISABLE' 'LOGICAL' 'REPLICATION'
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/create-domain.html
// fix: left out where optional, since no quoted page defines it or what it needs: create_index_statement, storage_clause
createDomain
    : 'CREATE' ('OR' 'REPLACE')? 'USECASE'? 'DOMAIN' (ID '.')? ID 'AS' (createDomain_create_single_column_domain | createDomain_create_multi_column_domain | createDomain_create_flexible_domain)
    | 'CREATE' 'USECASE'? 'DOMAIN' ('IF' 'NOT' 'EXISTS')? (ID '.')? ID 'AS' (createDomain_create_single_column_domain | createDomain_create_multi_column_domain | createDomain_create_flexible_domain)
    ;
createDomain_create_single_column_domain
    : 'ENUM' '(' createDomain_enum_list ')' 'STRICT'?
    | ('DEFAULT' ('ON' 'NULL')? EXPR)?
    | createDomain_constraint_clause*
    | ('COLLATE' ID)?
    | ('DISPLAY' EXPR)?
    | ('ORDER' EXPR)?
    | createDomain_annotations_clause?
    ;
createDomain_create_multi_column_domain
    : '(' ID 'AS' DATA_TYPE createDomain_annotations_clause? (',' ID 'AS' DATA_TYPE createDomain_annotations_clause?)* ')' createDomain_constraint_clause* ('COLLATE' ID)? ('DISPLAY' EXPR)? ('ORDER' EXPR)? createDomain_annotations_clause?
    ;
createDomain_create_flexible_domain
    : 'FLEXIBLE' 'DOMAIN' (ID '.')? ID '(' ID (',' ID)* ')' 'CHOOSE' 'DOMAIN' 'USING' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')' 'FROM' ('DECODE' '(' EXPR ',' EXPR ',' createDomain_result_expr (',' EXPR ',' createDomain_result_expr)* ')' | 'CASE' ('WHEN' CONDITION 'THEN' createDomain_result_expr)+ 'END')
    ;
createDomain_enum_list
    : createDomain_enum_item_list (',' createDomain_enum_item_list)*
    ;
createDomain_constraint_clause
    : ('CONSTRAINT' ID)? ('NOT' 'NULL' | 'NULL' | 'CHECK' '(' CONDITION ')') createDomain_constraint_state?
    ;
createDomain_annotations_clause
    : 'ANNOTATIONS' '(' createDomain_annotation (',' createDomain_annotation)* ')'
    ;
createDomain_result_expr
    : (ID '.')? ID '(' ID (',' ID)* ')'
    ;
createDomain_enum_item_list
    : ID ('=' createDomain_enum_alias_list)? ('=' EXPR)?
    ;
createDomain_constraint_state
    : ('NOT'? 'DEFERRABLE' ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? | 'INITIALLY' ('DEFERRED' | 'IMMEDIATE') ('NOT'? 'DEFERRABLE')?)? ('RELY' | 'NORELY')? createDomain_using_index_clause? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? createDomain_exceptions_clause?
    ;
createDomain_annotation
    : ID STRING?
    ;
createDomain_enum_alias_list
    : ID ('=' ID)*
    ;
createDomain_using_index_clause
    : 'USING' 'INDEX' ((ID '.')? ID | createDomain_index_attributes)
    ;
createDomain_exceptions_clause
    : 'EXCEPTIONS' 'INTO' (ID '.')? ID
    ;
createDomain_index_attributes
    : (createDomain_physical_attributes_clause | 'TABLESPACE' (ID | 'DEFAULT') | createDomain_logging_clause | createDomain_index_compression | createDomain_partial_index_clause | createDomain_parallel_clause)+
    ;
createDomain_physical_attributes_clause
    : ('PCTUSED' NUMBER)? ('PCTFREE' NUMBER)? ('INITRANS' NUMBER)?
    ;
createDomain_logging_clause
    : ('LOGGING' | 'NOLOGGING')
    ;
createDomain_index_compression
    : (createDomain_prefix_compression | createDomain_advanced_index_compression | 'NOCOMPRESS')
    ;
createDomain_partial_index_clause
    : 'INDEXING' ('FULL' | 'PARTIAL')
    ;
createDomain_parallel_clause
    : ('PARALLEL' NUMBER? | 'NOPARALLEL')
    ;
createDomain_prefix_compression
    : 'COMPRESS' NUMBER?
    ;
createDomain_advanced_index_compression
    : 'COMPRESS' 'ADVANCED' ('LOW' | 'HIGH')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/create-json-relational-duality-view.html
createJsonRelationalDualityView
    : 'CREATE' ('OR' 'REPLACE' | 'IF' 'NOT' 'EXISTS')? 'JSON' 'RELATIONAL' 'DUALITY' 'VIEW' (ID '.')? ID createJsonRelationalDualityView_duality_view_replication_clause? 'AS' createJsonRelationalDualityView_duality_view_subquery
    ;
createJsonRelationalDualityView_duality_view_replication_clause
    : ('ENABLE' | 'DISABLE') 'LOGICAL' 'REPLICATION'
    ;
createJsonRelationalDualityView_duality_view_subquery
    : 'SELECT' 'JSON' createJsonRelationalDualityView_key_value_clause (',' createJsonRelationalDualityView_key_value_clause)* 'FROM' ID createJsonRelationalDualityView_table_tags_clause? ('WHERE' CONDITION)?
    ;
createJsonRelationalDualityView_key_value_clause
    : STRING ':' ID
    | ID ':' ID createJsonRelationalDualityView_column_tags_clause?
    | ID ':' '(' QUERY ')' createJsonRelationalDualityView_column_tags_clause?
    | ID ':' 'JSON_ARRAYAGG' '(' 'JSON' createJsonRelationalDualityView_key_value_clause (',' createJsonRelationalDualityView_key_value_clause)* ')'
    ;
createJsonRelationalDualityView_table_tags_clause
    : 'WITH' 'INSERT'? 'UPDATE'? 'DELETE'? (('CHECK' | 'NOCHECK') 'ETAG')?
    ;
createJsonRelationalDualityView_column_tags_clause
    : 'WITH' (('UPDATE' | 'NOUPDATE'))? (('CHECK' | 'NOCHECK') 'ETAG')?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/create-logical-partition-tracking.html
createLogicalPartitionTracking
    : 'CREATE' 'LOGICAL' 'PARTITION' 'TRACKING' 'ON' (ID '.')? ID 'PARTITION' 'BY' ('RANGE' | 'INTERVAL') '(' ID ')' '(' createLogicalPartitionTracking_partition_definition (',' createLogicalPartitionTracking_partition_definition)* ')'
    ;
createLogicalPartitionTracking_partition_definition
    : 'PARTITION' ID 'VALUES' 'LESS' 'THAN' '(' EXPR ')'
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/create-mle-env.html
createMleEnv
    : 'CREATE' ('OR' 'REPLACE' | 'IF' 'NOT' 'EXISTS')? 'PURE'? 'MLE' 'ENV' (ID '.')? ID ('CLONE' (ID '.')? ID)? createMleEnv_imports_clause? createMleEnv_language_options_clause?
    ;
createMleEnv_imports_clause
    : 'IMPORTS' '(' createMleEnv_import_item (',' createMleEnv_import_item)* ')'
    ;
createMleEnv_language_options_clause
    : 'LANGUAGE' 'OPTIONS' STRING
    ;
createMleEnv_import_item
    : ID 'MODULE' (ID '.')? ID
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/create-mle-module.html
createMleModule
    : 'CREATE' ('OR' 'REPLACE' | 'IF' 'NOT' 'EXISTS')? 'MLE' 'MODULE' (ID '.')? ID 'LANGUAGE' 'JAVASCRIPT' ('VERSION' STRING)? ('USING' ('CLOB' '(' QUERY ')' | 'BLOB' '(' QUERY ')' | 'BFILE' '(' QUERY ')' | 'BFILE' '(' ID ',' STRING ')') | 'AS' ID)
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/create-pmem-filestore.html
createPmemFilestore
    : 'CREATE' 'PMEM' 'FILESTORE' ID 'MOUNTPOINT' STRING 'BACKINGFILE' STRING 'SIZE' NUMBER 'BLOCKSIZE' NUMBER ('AUTOEXTEND' ('ON' | 'OFF') ('NEXT' NUMBER)? ('MAXSIZE' ('UNLIMITED' | NUMBER))?)?
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/create-property-graph.html
createPropertyGraph
    : 'CREATE' ('OR' 'REPLACE')? 'PROPERTY' 'GRAPH' (ID '.')? ID createPropertyGraph_vertex_tables_clause createPropertyGraph_edge_tables_clause? createPropertyGraph_graph_options?
    | 'CREATE' 'PROPERTY' 'GRAPH' ('IF' 'NOT' 'EXISTS')? (ID '.')? ID createPropertyGraph_vertex_tables_clause createPropertyGraph_edge_tables_clause? createPropertyGraph_graph_options?
    ;
createPropertyGraph_vertex_tables_clause
    : 'VERTEX' 'TABLES' '(' createPropertyGraph_vertex_table_definition (',' createPropertyGraph_vertex_table_definition)* ')'
    ;
createPropertyGraph_edge_tables_clause
    : 'EDGE' 'TABLES' '(' createPropertyGraph_edge_table_definition (',' createPropertyGraph_edge_table_definition)* ')'
    ;
createPropertyGraph_graph_options
    : 'OPTIONS' '(' (('ENFORCED' 'MODE' | 'TRUSTED' 'MODE'))? (',' ('ALLOW' 'MIXED' 'PROPERTY' 'TYPES' | 'DISALLOW' 'MIXED' 'PROPERTY' 'TYPES'))? ')'
    ;
createPropertyGraph_vertex_table_definition
    : createPropertyGraph_graph_element_name_and_key createPropertyGraph_graph_table_label_and_properties?
    ;
createPropertyGraph_edge_table_definition
    : (ID '.')? createPropertyGraph_graph_element_object_name ('AS' ID)? createPropertyGraph_graph_element_key? 'SOURCE' ('KEY' '(' ID (',' ID)* ')' 'REFERENCES')? createPropertyGraph_vertex_table_reference 'DESTINATION' ('KEY' '(' ID (',' ID)* ')' 'REFERENCES')? createPropertyGraph_vertex_table_reference createPropertyGraph_graph_table_label_and_properties?
    ;
createPropertyGraph_graph_element_name_and_key
    : (ID '.')? createPropertyGraph_graph_element_object_name ('AS' ID)? createPropertyGraph_graph_element_key?
    ;
createPropertyGraph_graph_table_label_and_properties
    : createPropertyGraph_graph_table_label_clause? createPropertyGraph_graph_table_label_properties_clause? ('DEFAULT' 'LABEL')?
    ;
createPropertyGraph_graph_element_object_name
    : ID
    ;
createPropertyGraph_graph_element_key
    : 'KEY' '(' ID (',' ID)* ')'
    ;
createPropertyGraph_vertex_table_reference
    : ID '(' ID (',' ID)* ')'
    ;
createPropertyGraph_graph_table_label_clause
    : 'LABEL' ID ('LABEL' ID)*
    ;
createPropertyGraph_graph_table_label_properties_clause
    : 'PROPERTIES' 'ARE'? 'ALL' 'COLUMNS'
    | 'PROPERTIES' 'ARE'? 'ALL' 'COLUMNS' 'EXCEPT' '(' ID (',' ID)* ')'
    | 'PROPERTIES' '(' ID ('AS' ID)? (',' ID ('AS' ID)?)* ')'
    | 'PROPERTIES' '(' EXPR 'AS' ID (',' EXPR 'AS' ID)* ')'
    | 'NO' 'PROPERTIES'
    ;

// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/create-vector-index.html
createVectorIndex
    : 'CREATE' 'VECTOR' 'INDEX' ('IF' 'NOT' 'EXISTS')? (ID '.')? ID 'ON' (ID '.')? ID '(' ID ')' ('INCLUDE' '(' ID (',' ID)* ')')? createVectorIndex_vector_index_organization_clause? ('DISTANCE' ID)? ('WITH' 'TARGET' 'ACCURACY' NUMBER ('PARAMETERS' '(' createVectorIndex_vector_index_parameters_clause ')')?)? createVectorIndex_vector_index_hnsw_replication_clause? 'ONLINE'? ('PARALLEL' NUMBER)?
    ;
createVectorIndex_vector_index_organization_clause
    : 'ORGANIZATION' ('INMEMORY' 'NEIGHBOR' 'GRAPH' | 'NEIGHBOR' 'PARTITIONS')
    ;
createVectorIndex_vector_index_parameters_clause
    : (createVectorIndex_vector_index_parameters_hnsw_clause | createVectorIndex_vector_index_parameters_ivf_clause)
    ;
createVectorIndex_vector_index_hnsw_replication_clause
    : ('DUPLICATE' 'ALL' | 'DISTRIBUTE' ('AUTO' | 'BY' 'ROWID' 'RANGE' | 'BY' 'PARTITION' | 'BY' 'SUBPARTITION')?)
    ;
createVectorIndex_vector_index_parameters_hnsw_clause
    : 'TYPE' 'HNSW' (',' ('NEIGHBORS' | 'M') NUMBER)? (',' 'EFCONSTRUCTION' NUMBER)? (',' 'RESCORE_FACTOR' NUMBER)? (',' 'QUANTIZATION' ID)?
    ;
createVectorIndex_vector_index_parameters_ivf_clause
    : 'TYPE' 'IVF' (',' 'NEIGHBOR' 'PARTITIONS' NUMBER)? (',' 'SAMPLES_PER_PARTITION' NUMBER)? (',' 'MIN_VECTORS_PER_PARTITION' NUMBER)? (',' 'QUANTIZATION' 'SCALAR' NUMBER?)?
    ;
