// Amazon Redshift's SQL command reference as an ANTLR 4 grammar: the
// Syntax sections of the pages
// https://docs.aws.amazon.com/redshift/latest/dg/c_SQL_commands.html lists,
// converted by redshift/parser/cmd/docsgrammar. Do not edit it: correct a
// page's text in docs-fixes.tsv and run the converter again.
//
// Keywords are literals. ID, STRING, NUMBER, EXPR, QUERY, DATA_TYPE, CTE,
// ARGMODE and SCHEMA_ELEMENT are the values the docs write in italics or in lower case;
// tokens.tsv samples them. A form a page writes with ... for text it leaves out, or that does
// not read as notation, is not converted and is listed under its page. The
// pages a command page lists under Topics, such as SELECT's clauses, are
// not read: their syntax is a clause's, not a statement's.
//
// pin: default https://docs.aws.amazon.com/redshift/latest/dg/c_SQL_commands.html sha256:7f3b2ea3e1081db8a627bff3caec95c79f80a05bb6eb5e04d0c6197dce438007 2026-11-13

grammar RedshiftDocs;

statement
    : abort
    | alterDatabase
    | alterDatashare
    | alterDefaultPrivileges
    | alterExternalSchema
    | alterExternalView
    | alterFunction
    | alterGroup
    | alterIdentityProvider
    | alterMaskingPolicy
    | alterMaterializedView
    | alterRlsPolicy
    | alterRole
    | alterProcedure
    | alterSchema
    | alterSystem
    | alterTable
    | alterTableAppend
    | alterTemplate
    | alterUser
    | analyze
    | analyzeCompression
    | attachMaskingPolicy
    | attachRlsPolicy
    | begin
    | callProcedure
    | cancel
    | close
    | comment
    | commit
    | copy
    | createDatabase
    | createDatashare
    | createExternalFunction
    | createExternalModel
    | createExternalSchema
    | createExternalTable
    | createExternalView
    | createFunction
    | createGroup
    | createIdentityProvider
    | createLibrary
    | createMaskingPolicy
    | materializedViewCreateSqlCommand
    | createModel
    | createProcedure
    | createRlsPolicy
    | createRole
    | createSchema
    | createTableNew
    | createTableAs
    | createTemplate
    | createUser
    | createView
    | deallocate
    | declare
    | delete
    | descDatashare
    | descIdentityProvider
    | detachMaskingPolicy
    | detachRlsPolicy
    | dropDatabase
    | dropDatashare
    | dropExternalView
    | dropFunction
    | dropGroup
    | dropIdentityProvider
    | dropLibrary
    | dropMaskingPolicy
    | dropModel
    | materializedViewDropSqlCommand
    | dropProcedure
    | dropRlsPolicy
    | dropRole
    | dropSchema
    | dropTable
    | dropTemplate
    | dropUser
    | dropView
    | end
    | execute
    | explain
    | fetch
    | grant
    | insert30
    | insertExternalTable
    | lock
    | merge
    | prepare
    | materializedViewRefreshSqlCommand
    | reset
    | revoke
    | rollback
    | selectSynopsis
    | selectInto
    | set
    | setSessionAuthorization
    | show
    | showUserLockout
    | showColumnGrants
    | showColumns
    | showConstraints
    | showExternalTable
    | showDatabases
    | showFunctions
    | showGrants
    | showModel
    | showDatashares
    | showParameters
    | showPolicies
    | showProcedure
    | showProcedures
    | showSchemas
    | showTable
    | showTables
    | showTemplate
    | showTemplates
    | showView
    | truncate
    | unload
    | update
    | useCommand
    | vacuumCommand
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ABORT.html
abort
    : 'ABORT' ('WORK' | 'TRANSACTION')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_DATABASE.html
alterDatabase
    : 'ALTER' 'DATABASE' ID ('RENAME' 'TO' ID | 'OWNER' 'TO' ID | ('CONNECTION' 'LIMIT' (NUMBER | 'UNLIMITED'))? ('COLLATE' ('CASE_SENSITIVE' | 'CS' | 'CASE_INSENSITIVE' | 'CI'))? ('ISOLATION' 'LEVEL' ('SNAPSHOT' | 'SERIALIZABLE'))? | 'INTEGRATION' ('REFRESH' (('ALL' | 'INERROR' | 'REMEDIABLE') 'TABLES' ('IN' 'SCHEMA' ID (',' ID)*)? | 'TABLE' ID (',' ID)*) | 'SET' ('QUERY_ALL_STATES' '='? ('TRUE' | 'FALSE') ('ACCEPTINVCHARS' '='? ('TRUE' | 'FALSE'))? ('REFRESH_INTERVAL' NUMBER)? ('TRUNCATECOLUMNS' '='? ('TRUE' | 'FALSE'))? ('HISTORY_MODE' '='? ('TRUE' | 'FALSE') ('FOR' ('ALL' 'TABLES' ('IN' 'SCHEMA' ID (',' ID)*)? | 'TABLE' ID (',' ID)*))?)? ('AUTO_REMEDIATION' '='? ('TRUE' | 'FALSE'))? | 'ACCEPTINVCHARS' '='? ('TRUE' | 'FALSE') ('REFRESH_INTERVAL' NUMBER)? ('TRUNCATECOLUMNS' '='? ('TRUE' | 'FALSE'))? ('HISTORY_MODE' '='? ('TRUE' | 'FALSE') ('FOR' ('ALL' 'TABLES' ('IN' 'SCHEMA' ID (',' ID)*)? | 'TABLE' ID (',' ID)*))?)? ('AUTO_REMEDIATION' '='? ('TRUE' | 'FALSE'))? | 'REFRESH_INTERVAL' NUMBER ('TRUNCATECOLUMNS' '='? ('TRUE' | 'FALSE'))? ('HISTORY_MODE' '='? ('TRUE' | 'FALSE') ('FOR' ('ALL' 'TABLES' ('IN' 'SCHEMA' ID (',' ID)*)? | 'TABLE' ID (',' ID)*))?)? ('AUTO_REMEDIATION' '='? ('TRUE' | 'FALSE'))? | 'TRUNCATECOLUMNS' '='? ('TRUE' | 'FALSE') ('HISTORY_MODE' '='? ('TRUE' | 'FALSE') ('FOR' ('ALL' 'TABLES' ('IN' 'SCHEMA' ID (',' ID)*)? | 'TABLE' ID (',' ID)*))?)? ('AUTO_REMEDIATION' '='? ('TRUE' | 'FALSE'))? | 'HISTORY_MODE' '='? ('TRUE' | 'FALSE') ('FOR' ('ALL' 'TABLES' ('IN' 'SCHEMA' ID (',' ID)*)? | 'TABLE' ID (',' ID)*))? ('AUTO_REMEDIATION' '='? ('TRUE' | 'FALSE'))? | 'AUTO_REMEDIATION' '='? ('TRUE' | 'FALSE'))))
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_DATASHARE.html
// fix: the value is TRUE or FALSE; the choice is braced
// fix: the value is TRUE or FALSE; the choice is braced
alterDatashare
    : 'ALTER' 'DATASHARE' ID ('ADD' | 'REMOVE') ('TABLE' ID (',' ID)* | 'SCHEMA' ID (',' ID)* | 'FUNCTION' ID '(' DATA_TYPE (',' DATA_TYPE)* ')' (',' ID '(' DATA_TYPE (',' DATA_TYPE)* ')')* | 'ALL' 'TABLES' 'IN' 'SCHEMA' ID (',' ID)* | 'ALL' 'FUNCTIONS' 'IN' 'SCHEMA' ID (',' ID)*)
    | 'ALTER' 'DATASHARE' ID ('SET' 'PUBLICACCESSIBLE' '='? ('TRUE' | 'FALSE'))? ('SET' 'INCLUDENEW' '='? ('TRUE' | 'FALSE') 'FOR' 'SCHEMA' ID)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_DEFAULT_PRIVILEGES.html
alterDefaultPrivileges
    : 'ALTER' 'DEFAULT' 'PRIVILEGES' ('FOR' 'USER' ID (',' ID)*)? ('IN' 'SCHEMA' ID (',' ID)*)? alterDefaultPrivileges_grant_or_revoke_clause
    ;
alterDefaultPrivileges_grant_or_revoke_clause
    : 'GRANT' (('SELECT' | 'INSERT' | 'UPDATE' | 'DELETE' | 'DROP' | 'REFERENCES' | 'TRUNCATE') (',' ('SELECT' | 'INSERT' | 'UPDATE' | 'DELETE' | 'DROP' | 'REFERENCES' | 'TRUNCATE'))* | 'ALL' 'PRIVILEGES'?) 'ON' 'TABLES' 'TO' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))*
    | 'GRANT' ('EXECUTE' | 'ALL' 'PRIVILEGES'?) 'ON' 'FUNCTIONS' 'TO' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))*
    | 'GRANT' ('EXECUTE' | 'ALL' 'PRIVILEGES'?) 'ON' 'PROCEDURES' 'TO' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))*
    | 'REVOKE' ('GRANT' 'OPTION' 'FOR')? (('SELECT' | 'INSERT' | 'UPDATE' | 'DELETE' | 'REFERENCES' | 'TRUNCATE') (',' ('SELECT' | 'INSERT' | 'UPDATE' | 'DELETE' | 'REFERENCES' | 'TRUNCATE'))* | 'ALL' 'PRIVILEGES'?) 'ON' 'TABLES' 'FROM' ID (',' ID)* 'RESTRICT'?
    | 'REVOKE' (('SELECT' | 'INSERT' | 'UPDATE' | 'DELETE' | 'REFERENCES' | 'TRUNCATE') (',' ('SELECT' | 'INSERT' | 'UPDATE' | 'DELETE' | 'REFERENCES' | 'TRUNCATE'))* | 'ALL' 'PRIVILEGES'?) 'ON' 'TABLES' 'FROM' ('ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' ('ROLE' ID | 'GROUP' ID | 'PUBLIC'))* 'RESTRICT'?
    | 'REVOKE' ('GRANT' 'OPTION' 'FOR')? ('EXECUTE' | 'ALL' 'PRIVILEGES'?) 'ON' 'FUNCTIONS' 'FROM' ID (',' ID)* 'RESTRICT'?
    | 'REVOKE' ('EXECUTE' | 'ALL' 'PRIVILEGES'?) 'ON' 'FUNCTIONS' 'FROM' ('ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' ('ROLE' ID | 'GROUP' ID | 'PUBLIC'))* 'RESTRICT'?
    | 'REVOKE' ('GRANT' 'OPTION' 'FOR')? ('EXECUTE' | 'ALL' 'PRIVILEGES'?) 'ON' 'PROCEDURES' 'FROM' ID (',' ID)* 'RESTRICT'?
    | 'REVOKE' ('EXECUTE' | 'ALL' 'PRIVILEGES'?) 'ON' 'PROCEDURES' 'FROM' ('ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' ('ROLE' ID | 'GROUP' ID | 'PUBLIC'))* 'RESTRICT'?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_EXTERNAL_SCHEMA.html
alterExternalSchema
    : 'ALTER' 'EXTERNAL' 'SCHEMA' ID ('IAM_ROLE' ('DEFAULT' | '\'SESSION\'' | STRING)?)? ('AUTHENTICATION' ('NONE' | 'IAM' | 'MTLS')?)? ('AUTHENTICATION_ARN' STRING | 'SECRET_ARN' STRING)? ('URI' STRING)?
    | 'ALTER' 'EXTERNAL' 'SCHEMA' ID 'AUTHENTICATION' 'MTLS' 'AUTHENTICATION_ARN' STRING
    | 'ALTER' 'EXTERNAL' 'SCHEMA' ID 'AUTHENTICATION' 'MTLS' 'SECRET_ARN' STRING
    | 'ALTER' 'EXTERNAL' 'SCHEMA' ID 'URI' STRING
    | 'ALTER' 'EXTERNAL' 'SCHEMA' ID 'IAM_ROLE' STRING
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_EXTERNAL_VIEW.html
alterExternalView
    : 'ALTER' 'EXTERNAL' 'VIEW' ID ID 'FORCE'? ('AS' '(' QUERY ')' | 'REMOVE' 'DEFINITION')
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_FUNCTION.html
// fix: the argument list opens { before [ and closes } before ]; swapped
alterFunction
    : 'ALTER' 'FUNCTION' ID '(' ((ID DATA_TYPE | DATA_TYPE) (',' (ID DATA_TYPE | DATA_TYPE))*)? ')' 'RENAME' 'TO' ID
    | 'ALTER' 'FUNCTION' ID '(' ((ID DATA_TYPE | DATA_TYPE) (',' (ID DATA_TYPE | DATA_TYPE))*)? ')' 'OWNER' 'TO' (ID | 'CURRENT_USER' | 'SESSION_USER')
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_GROUP.html
alterGroup
    : 'ALTER' 'GROUP' ID ('ADD' 'USER' ID (',' ID)* | 'DROP' 'USER' ID (',' ID)* | 'RENAME' 'TO' ID)
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_IDENTITY_PROVIDER.html
// fix: [AUTO_CREATE_ROLES is not closed; closed
alterIdentityProvider
    : 'ALTER' 'IDENTITY' 'PROVIDER' ID ('PARAMETERS' ID)? ('NAMESPACE' ID)? ('IAM_ROLE' ID)? ('AUTO_CREATE_ROLES' ('TRUE' (('INCLUDE' | 'EXCLUDE') 'GROUPS' 'LIKE' ID)? | 'FALSE')?)? ('DISABLE' | 'ENABLE')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_MASKING_POLICY.html
alterMaskingPolicy
    : 'ALTER' 'MASKING' 'POLICY' (ID | ID '.' ID) 'USING' '(' EXPR ')'
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_MATERIALIZED_VIEW.html
alterMaterializedView
    : 'ALTER' 'MATERIALIZED' 'VIEW' ID ('AUTO' 'REFRESH' ('YES' ('INTERVAL' NUMBER)? | 'NO') | 'ALTER' 'DISTKEY' ID | 'ALTER' 'DISTSTYLE' 'ALL' | 'ALTER' 'DISTSTYLE' 'EVEN' | 'ALTER' 'DISTSTYLE' 'KEY' 'DISTKEY' ID | 'ALTER' 'DISTSTYLE' 'AUTO' | 'ALTER' 'COMPOUND'? 'SORTKEY' '(' ID (',' ID)* ')' | 'ALTER' 'SORTKEY' 'AUTO' | 'ALTER' 'SORTKEY' 'NONE' | 'ROW' 'LEVEL' 'SECURITY' ('ON' | 'OFF') ('CONJUNCTION' 'TYPE' ('AND' | 'OR'))? ('FOR' 'DATASHARES')?)
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_RLS_POLICY.html
alterRlsPolicy
    : 'ALTER' 'RLS' 'POLICY' ID 'USING' '(' EXPR ')'
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_ROLE.html
alterRole
    : 'ALTER' 'ROLE' ID 'WITH'? ('RENAME' 'TO' ID | 'OWNER' 'TO' ID) (',' ('RENAME' 'TO' ID | 'OWNER' 'TO' ID))* ('EXTERNALID' 'TO' ID)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_PROCEDURE.html
alterProcedure
    : 'ALTER' 'PROCEDURE' ID ('(' (ID? ARGMODE? DATA_TYPE (',' DATA_TYPE)*)? ')')? 'RENAME' 'TO' ID
    | 'ALTER' 'PROCEDURE' ID ('(' (ID? ARGMODE? DATA_TYPE (',' DATA_TYPE)*)? ')')? 'OWNER' 'TO' (ID | 'CURRENT_USER' | 'SESSION_USER')
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_SCHEMA.html
alterSchema
    : 'ALTER' 'SCHEMA' ID ('RENAME' 'TO' ID | 'OWNER' 'TO' ID | 'QUOTA' (NUMBER ('MB' | 'GB' | 'TB')? | 'UNLIMITED'))
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_SYSTEM.html
alterSystem
    : 'ALTER' 'SYSTEM' 'SET' ID '=' ID
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_TABLE.html
// fix: two lines that say ALTER COLUMN ... ENCODE repeats with commas, written with .....; dropped, leaving the one-column form
// fix: ADD COLUMN: a | too many after it; dropped
// fix: the external table actions follow a sentence, "The following options apply only to external tables:"; read as the definition of an action ALTER TABLE takes
// fix: the external table actions are ALTER TABLE's too
// fix: external table actions: a | too many after SET FILE FORMAT; dropped
alterTable
    : 'ALTER' 'TABLE' ID ('ADD' alterTable_table_constraint | 'DROP' 'CONSTRAINT' ID ('RESTRICT' | 'CASCADE')? | 'OWNER' 'TO' ID | 'RENAME' 'TO' ID | 'RENAME' 'COLUMN' ID 'TO' ID | 'ALTER' 'COLUMN' ID 'TYPE' DATA_TYPE | 'ALTER' 'COLUMN' ID 'ENCODE' ID | 'ALTER' 'DISTKEY' ID | 'ALTER' 'DISTSTYLE' 'ALL' | 'ALTER' 'DISTSTYLE' 'EVEN' | 'ALTER' 'DISTSTYLE' 'KEY' 'DISTKEY' ID | 'ALTER' 'DISTSTYLE' 'AUTO' | 'ALTER' 'COMPOUND'? 'SORTKEY' '(' ID (',' ID)* ')' | 'ALTER' 'SORTKEY' 'AUTO' | 'ALTER' 'SORTKEY' 'NONE' | 'ALTER' 'ENCODE' 'AUTO' | 'ADD' 'COLUMN'? ID DATA_TYPE ('DEFAULT' EXPR)? ('ENCODE' ID)? ('NOT' 'NULL' | 'NULL')? ('COLLATE' ('CASE_SENSITIVE' | 'CS' | 'CASE_INSENSITIVE' | 'CI'))? | 'DROP' 'COLUMN'? ID ('RESTRICT' | 'CASCADE')? | 'ROW' 'LEVEL' 'SECURITY' ('ON' | 'OFF') ('CONJUNCTION' 'TYPE' ('AND' | 'OR'))? ('FOR' 'DATASHARES')? | 'MASKING' ('ON' | 'OFF') 'FOR' 'DATASHARES' | alterTable_external_table_action)
    | 'ALTER' 'TABLE' ID 'ALTER' 'SORTKEY' '(' ID ')' ',' 'ALTER' 'DISTKEY' ID
    | 'ALTER' 'TABLE' ID 'ALTER' 'DISTKEY' ID ',' 'ALTER' 'SORTKEY' '(' ID ')'
    | 'ALTER' 'TABLE' ID 'ALTER' 'SORTKEY' '(' ID ')' ',' 'ALTER' 'DISTSTYLE' 'ALL'
    | 'ALTER' 'TABLE' ID 'ALTER' 'DISTSTYLE' 'ALL' ',' 'ALTER' 'SORTKEY' '(' ID ')'
    ;
alterTable_table_constraint
    : ('CONSTRAINT' ID)? ('UNIQUE' '(' ID (',' ID)* ')' | 'PRIMARY' 'KEY' '(' ID (',' ID)* ')' | 'FOREIGN' 'KEY' '(' ID (',' ID)* ')' 'REFERENCES' ID ('(' ID ')')?)
    ;
alterTable_external_table_action
    : 'SET' 'LOCATION' STRING
    | 'SET' 'FILE' 'FORMAT' ID
    | 'SET' 'TABLE' 'PROPERTIES' '(' STRING '=' STRING ')'
    | 'PARTITION' '(' ID '=' ID (',' ID '=' ID)* ')' 'SET' 'LOCATION' STRING
    | 'ADD' ('IF' 'NOT' 'EXISTS')? 'PARTITION' '(' ID '=' ID (',' ID '=' ID)* ')' 'LOCATION' STRING (',' STRING)*
    | 'DROP' 'PARTITION' '(' ID '=' ID (',' ID '=' ID)* ')'
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_TABLE_APPEND.html
// fix: the source is required, as the Parameters section says; the choice is braced
alterTableAppend
    : 'ALTER' 'TABLE' ID 'APPEND' 'FROM' ID ('IGNOREEXTRA' | 'FILLTARGET')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_TEMPLATE.html
alterTemplate
    : 'ALTER' 'TEMPLATE' (ID '.')? (ID '.')? ID ('RENAME' 'TO' ID | 'OWNER' 'TO' ID | 'ADD' ID 'AS'? ID? | 'DROP' ID | 'SET' ID 'TO' ID (',' ID 'TO' ID)*)
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_USER.html
// fix: option: a | too many after RENAME TO; dropped
alterUser
    : 'ALTER' 'USER' ID 'WITH'? alterUser_option (',' alterUser_option)*
    ;
alterUser_option
    : 'CREATEDB'
    | 'NOCREATEDB'
    | 'CREATEUSER'
    | 'NOCREATEUSER'
    | 'SYSLOG' 'ACCESS' ('RESTRICTED' | 'UNRESTRICTED')
    | 'PASSWORD' (STRING | 'DISABLE') ('VALID' 'UNTIL' STRING)?
    | 'RENAME' 'TO' ID
    | 'CONNECTION' 'LIMIT' (NUMBER | 'UNLIMITED')
    | 'SESSION' 'TIMEOUT' NUMBER
    | 'RESET' 'SESSION' 'TIMEOUT'
    | 'SET' ID ('TO' | '=') (ID | 'DEFAULT')
    | 'RESET' ID
    | 'EXTERNALID' ID
    | 'NOLOGIN'
    | 'LOGIN' 'PASSWORD' STRING
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ANALYZE.html
// fix: a [ too many before the table name; dropped
analyze
    : 'ANALYZE' 'VERBOSE'? (ID ('(' ID (',' ID)* ')')?)? ('PREDICATE' 'COLUMNS' | 'ALL' 'COLUMNS')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ANALYZE_COMPRESSION.html
analyzeCompression
    : 'ANALYZE' 'COMPRESSION' (ID? ('(' ID (',' ID)* ')')?)? ('COMPROWS' NUMBER)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ATTACH_MASKING_POLICY.html
attachMaskingPolicy
    : 'ATTACH' 'MASKING' 'POLICY' (ID 'ON' ID | ID '.' ID 'ON' ID '.' ID '.' ID) '(' ID ')' ('USING' '(' ID ')')? 'TO' (ID | 'ROLE' ID | 'PUBLIC') ('PRIORITY' NUMBER)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ATTACH_RLS_POLICY.html
attachRlsPolicy
    : 'ATTACH' 'RLS' 'POLICY' ID 'ON' 'TABLE'? ID (',' ID)* 'TO' (ID | 'ROLE' ID | 'PUBLIC') (',' (ID | 'ROLE' ID | 'PUBLIC'))*
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_BEGIN.html
// fix: a note, not syntax; dropped
begin
    : 'BEGIN' ('WORK' | 'TRANSACTION')? ('ISOLATION' 'LEVEL' begin_option)? ('READ' 'WRITE' | 'READ' 'ONLY')?
    | 'START' 'TRANSACTION' ('ISOLATION' 'LEVEL' begin_option)? ('READ' 'WRITE' | 'READ' 'ONLY')?
    ;
begin_option
    : 'SERIALIZABLE'
    | 'READ' 'UNCOMMITTED'
    | 'READ' 'COMMITTED'
    | 'REPEATABLE' 'READ'
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CALL_procedure.html
callProcedure
    : 'CALL' ID '(' (ID (',' ID)*)? ')'
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CANCEL.html
cancel
    : 'CANCEL' NUMBER STRING?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/close.html
close
    : 'CLOSE' ID
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_COMMENT.html
// fix: IS takes a string or NULL; the choice is braced
comment
    : 'COMMENT' 'ON' ('TABLE' ID | 'COLUMN' ID '.' ID | 'CONSTRAINT' ID 'ON' ID | 'DATABASE' ID | 'VIEW' ID) 'IS' (STRING | 'NULL')
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_COMMIT.html
commit
    : 'COMMIT' ('WORK' | 'TRANSACTION')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_COPY.html
// fix: an example, not syntax; dropped
// fix: an example, not syntax; dropped
// fix: the data source is an S3 path in quotes, and authorization IAM_ROLE as r_copy-parameters-authorization.html writes it
// fix: the data source is an S3 path in quotes, and authorization IAM_ROLE as r_copy-parameters-authorization.html writes it
copy
    : 'COPY' ID ID? 'FROM' STRING 'IAM_ROLE' ('DEFAULT' | STRING) ('FORMAT'? 'AS'? ID)? (ID ID? (',' ID ID?)*)?
    | 'COPY' ID 'FROM' STRING 'IAM_ROLE' ('DEFAULT' | STRING)
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_DATABASE.html
// fix: the datashare form leaves its [ open; closed
// fix: the datashare form leaves its [ open; closed
createDatabase
    : 'CREATE' 'DATABASE' ID (('FROM' 'INTEGRATION' STRING ('DATABASE' STRING)? 'SET'? ('ACCEPTINVCHARS' '='? ('TRUE' | 'FALSE'))? ('QUERY_ALL_STATES' '='? ('TRUE' | 'FALSE'))? ('REFRESH_INTERVAL' NUMBER)? ('TRUNCATECOLUMNS' '='? ('TRUE' | 'FALSE'))? ('HISTORY_MODE' '='? ('TRUE' | 'FALSE'))? ('AUTO_REMEDIATION' '='? ('TRUE' | 'FALSE'))?)? 'WITH'? ('OWNER' '='? ID)? ('CONNECTION' 'LIMIT' (NUMBER | 'UNLIMITED'))? ('COLLATE' ('CASE_SENSITIVE' | 'CS' | 'CASE_INSENSITIVE' | 'CI'))? ('ISOLATION' 'LEVEL' ('SNAPSHOT' | 'SERIALIZABLE'))? | 'FROM' 'ARN' STRING ('WITH' 'DATA' 'CATALOG' 'SCHEMA' STRING | 'WITH' 'NO' 'DATA' 'CATALOG' 'SCHEMA') | 'IAM_ROLE' ('DEFAULT' | '\'SESSION\'' | STRING) | ('WITH' 'PERMISSIONS')? 'FROM' 'DATASHARE' ID 'OF' ('ACCOUNT' ID)? 'NAMESPACE' ID)?
    | 'CREATE' 'DATABASE' ID (('WITH' 'PERMISSIONS')? 'FROM' 'DATASHARE' ID 'OF' ('ACCOUNT' ID)? 'NAMESPACE' ID)?
    | 'CREATE' 'DATABASE' ID (('WITH' 'PERMISSIONS')? 'FROM' 'DATASHARE' ID 'OF' 'ACCOUNT' ID 'NAMESPACE' ID)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_DATASHARE.html
// fix: the value is TRUE or FALSE; the choice is braced
createDatashare
    : 'CREATE' 'DATASHARE' ID ('SET'? 'PUBLICACCESSIBLE' '='? ('TRUE' | 'FALSE'))?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_EXTERNAL_FUNCTION.html
// fix: IAM_ROLE { is not closed; closed
createExternalFunction
    : 'CREATE' ('OR' 'REPLACE')? 'EXTERNAL' 'FUNCTION' ID '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'RETURNS' DATA_TYPE ('VOLATILE' | 'STABLE') 'LAMBDA' STRING 'IAM_ROLE' ('DEFAULT' | STRING) 'RETRY_TIMEOUT' NUMBER 'MAX_BATCH_ROWS' NUMBER 'MAX_BATCH_SIZE' NUMBER ('KB' | 'MB')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_create_external_model.html
// fix: IAM_ROLE: the choice is written with / for |
// fix: the inference function called, not CREATE EXTERNAL MODEL; dropped
// fix: the inference function called, not CREATE EXTERNAL MODEL; dropped
createExternalModel
    : 'CREATE' 'EXTERNAL' 'MODEL' ID 'FUNCTION' ID 'IAM_ROLE' ('DEFAULT' | STRING) 'MODEL_TYPE' 'BEDROCK' 'SETTINGS' '(' 'MODEL_ID' ID (',' 'PROMPT' STRING)? (',' 'SUFFIX' STRING)? (',' 'REQUEST_TYPE' ('RAW' | 'UNIFIED'))? (',' 'RESPONSE_TYPE' ('VARCHAR' | 'SUPER'))? ')'
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_EXTERNAL_SCHEMA.html
createExternalSchema
    : 'CREATE' 'EXTERNAL' 'SCHEMA' ('IF' 'NOT' 'EXISTS')? ID 'FROM' (('DATA' 'CATALOG')? | 'HIVE' 'METASTORE' | 'POSTGRES' | 'MYSQL' | 'KINESIS' | 'MSK' | 'REDSHIFT' | 'KAFKA')? ('DATABASE' STRING)? ('SCHEMA' STRING)? ('REGION' STRING)? ('IAM_ROLE' ('DEFAULT' | '\'SESSION\'' | STRING)?)? ('AUTHENTICATION' ('NONE' | 'IAM' | 'MTLS')?)? ('AUTHENTICATION_ARN' STRING | 'SECRET_ARN' STRING)? ('URI' (STRING ('PORT' NUMBER)? | STRING)?)? ('CLUSTER_ARN' STRING)? ('CATALOG_ROLE' ('\'SESSION\'' | STRING)?)? ('CREATE' 'EXTERNAL' 'DATABASE' 'IF' 'NOT' 'EXISTS')? ('CATALOG_ID' STRING)?
    | 'CREATE' 'EXTERNAL' 'SCHEMA' ('IF' 'NOT' 'EXISTS')? ID 'FROM' 'POSTGRES' 'DATABASE' STRING ('SCHEMA' STRING)? 'URI' STRING ('PORT' NUMBER)? 'IAM_ROLE' ('DEFAULT' | STRING)? 'SECRET_ARN' STRING
    | 'CREATE' 'EXTERNAL' 'SCHEMA' ('IF' 'NOT' 'EXISTS')? ID 'FROM' 'MYSQL' 'DATABASE' STRING 'URI' STRING ('PORT' NUMBER)? 'IAM_ROLE' ('DEFAULT' | STRING)? 'SECRET_ARN' STRING
    | 'CREATE' 'EXTERNAL' 'SCHEMA' ('IF' 'NOT' 'EXISTS')? ID 'FROM' 'KINESIS' 'IAM_ROLE' ('DEFAULT' | STRING)?
    | 'CREATE' 'EXTERNAL' 'SCHEMA' ('IF' 'NOT' 'EXISTS')? ID 'FROM' 'KAFKA' ('IAM_ROLE' ('DEFAULT' | STRING)?)? 'URI' STRING 'AUTHENTICATION' ('NONE' | 'IAM' | 'MTLS')? ('AUTHENTICATION_ARN' STRING | 'SECRET_ARN' STRING)?
    | 'CREATE' 'EXTERNAL' 'SCHEMA' ID 'FROM' 'REDSHIFT' 'DATABASE' STRING 'SCHEMA' STRING
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_EXTERNAL_TABLE.html
createExternalTable
    : 'CREATE' 'EXTERNAL' 'TABLE' ID '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')' ('PARTITIONED' 'BY' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')')? (('ROW' 'FORMAT' 'DELIMITED' ID | 'ROW' 'FORMAT' 'SERDE' STRING ('WITH' 'SERDEPROPERTIES' '(' STRING '=' STRING (',' STRING '=' STRING)* ')')?))? 'STORED' 'AS' ID 'LOCATION' STRING ('TABLE' 'PROPERTIES' '(' STRING '=' STRING (',' STRING '=' STRING)* ')')?
    | 'CREATE' 'EXTERNAL' 'TABLE' ID ('PARTITIONED' 'BY' '(' ID (',' ID)* ')')? ('ROW' 'FORMAT' 'DELIMITED' ID)? 'STORED' 'AS' ID 'LOCATION' STRING ('TABLE' 'PROPERTIES' '(' STRING '=' STRING (',' STRING '=' STRING)* ')')? 'AS' QUERY
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_EXTERNAL_VIEW.html
createExternalView
    : 'CREATE' 'EXTERNAL' 'VIEW' ID ('IF' 'NOT' 'EXISTS')? ID 'AS' QUERY
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_FUNCTION.html
// fix: the argument list opens { before [ and closes } before ]; swapped
createFunction
    : 'CREATE' ('OR' 'REPLACE')? 'FUNCTION' ID '(' ((ID DATA_TYPE | DATA_TYPE) (',' (ID DATA_TYPE | DATA_TYPE))*)? ')' 'RETURNS' DATA_TYPE ('VOLATILE' | 'STABLE' | 'IMMUTABLE') 'AS' STRING 'LANGUAGE' ('PLPYTHONU' | 'SQL')
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_GROUP.html
createGroup
    : 'CREATE' 'GROUP' ID (('WITH' ('USER' ID)? (',' 'WITH' ('USER' ID)?)*)?)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_IDENTITY_PROVIDER.html
createIdentityProvider
    : 'CREATE' 'IDENTITY' 'PROVIDER' ID 'TYPE' ID 'NAMESPACE' ID ('PARAMETERS' ID)? ('APPLICATION_ARN' ID)? ('IAM_ROLE' ID)? ('AUTO_CREATE_ROLES' ('TRUE' (('INCLUDE' | 'EXCLUDE') 'GROUPS' 'LIKE' ID)? | 'FALSE')?)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_LIBRARY.html
// fix: authorization names the IAM_ROLE clause after it; dropped
createLibrary
    : 'CREATE' ('OR' 'REPLACE')? 'LIBRARY' ID 'LANGUAGE' 'PLPYTHONU' 'FROM' (STRING | STRING ('REGION' 'AS'? STRING)? 'IAM_ROLE' ('DEFAULT' | STRING))
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_MASKING_POLICY.html
createMaskingPolicy
    : 'CREATE' 'MASKING' 'POLICY' (ID | ID '.' ID) ('IF' 'NOT' 'EXISTS')? 'WITH' '(' ID ')' 'USING' '(' EXPR ')'
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/materialized-view-create-sql-command.html
// fix: table_attributes are the DISTSTYLE, DISTKEY and SORTKEY the Parameters section lists
// fix: the block after the Iceberg form lists what partition_transform is; read as its definition
materializedViewCreateSqlCommand
    : 'CREATE' 'MATERIALIZED' 'VIEW' ID ('BACKUP' ('YES' | 'NO'))? ('DISTSTYLE' ('EVEN' | 'ALL' | 'KEY'))? ('DISTKEY' '(' ID ')')? ('SORTKEY' '(' ID (',' ID)* ')')? ('AUTO' 'REFRESH' ('YES' ('INTERVAL' NUMBER)? | 'NO'))? 'AS' QUERY
    | 'CREATE' 'MATERIALIZED' 'VIEW' ID 'USING' 'ICEBERG' ('LOCATION' STRING)? ('PARTITIONED' 'BY' '(' materializedViewCreateSqlCommand_partition_transform (',' materializedViewCreateSqlCommand_partition_transform)* ')')? ('TABLE' 'PROPERTIES' '(' STRING '=' STRING (',' STRING '=' STRING)* ')')? 'AS' QUERY
    ;
materializedViewCreateSqlCommand_partition_transform
    : ID
    | 'YEAR' '(' ID ')'
    | 'MONTH' '(' ID ')'
    | 'DAY' '(' ID ')'
    | 'HOUR' '(' ID ')'
    | 'BUCKET' '(' NUMBER ',' ID ')'
    | 'TRUNCATE' '(' NUMBER ',' ID ')'
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_MODEL.html
// fix: AUTO: the choice is written with / for |
// fix: the choice is written in parentheses, which are not SQL here (the examples write PROBLEM_TYPE BINARY_CLASSIFICATION)
// fix: the choice is written in parentheses, which are not SQL here (the examples write OBJECTIVE 'binary:logistic')
// fix: the choice is written in parentheses, which are not SQL here
// fix: OBJECTIVE: a , for a | between two choices
// fix: HYPERPARAMETERS: Key stands for a hyperparameter's name, and (,...) for a list of key-value pairs
// fix: SETTINGS: each setting is written with a comma and a | after it, and a comment; read as a list of settings, each as written
createModel
    : 'CREATE' 'MODEL' ID 'FROM' (ID | '(' QUERY ')' | STRING) ('TARGET' ID)? 'FUNCTION' ID ('(' DATA_TYPE (',' DATA_TYPE)* ')')? ('RETURNS' DATA_TYPE)? ('SAGEMAKER' STRING (':' STRING)?)? 'IAM_ROLE' ('DEFAULT' | STRING) ('AUTO' ('ON' | 'OFF'))? ('MODEL_TYPE' ('XGBOOST' | 'MLP' | 'LINEAR_LEARNER' | 'KMEANS' | 'FORECAST'))? ('PROBLEM_TYPE' ('REGRESSION' | 'BINARY_CLASSIFICATION' | 'MULTICLASS_CLASSIFICATION'))? ('OBJECTIVE' ('\'MSE\'' | STRING | '\'F1\'' | '\'AUC\'' | '\'RMSE\'' | '\'WAPE\'' | '\'MAPE\'' | '\'MASE\''))? ('PREPROCESSORS' STRING)? ('HYPERPARAMETERS' ('DEFAULT' | 'DEFAULT' 'EXCEPT' '(' ID STRING (',' ID STRING)* ')'))? ('SETTINGS' '(' ('S3_BUCKET' STRING | 'TAGS' STRING | 'KMS_KEY_ID' STRING | 'S3_GARBAGE_COLLECT' ('ON' | 'OFF') | 'MAX_CELLS' NUMBER | 'MAX_RUNTIME' NUMBER | 'HORIZON' NUMBER | 'FREQUENCY' NUMBER | 'PERCENTILES' STRING | 'MAX_BATCH_ROWS' NUMBER) (',' ('S3_BUCKET' STRING | 'TAGS' STRING | 'KMS_KEY_ID' STRING | 'S3_GARBAGE_COLLECT' ('ON' | 'OFF') | 'MAX_CELLS' NUMBER | 'MAX_RUNTIME' NUMBER | 'HORIZON' NUMBER | 'FREQUENCY' NUMBER | 'PERCENTILES' STRING | 'MAX_BATCH_ROWS' NUMBER))* ')')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_PROCEDURE.html
createProcedure
    : 'CREATE' ('OR' 'REPLACE')? 'PROCEDURE' ID '(' (ID? ARGMODE? DATA_TYPE (',' DATA_TYPE)*)? ')' 'NONATOMIC'? 'AS' STRING 'LANGUAGE' 'PLPGSQL' (('SECURITY' 'INVOKER' | 'SECURITY' 'DEFINER'))? ('SET' ID ('TO' ID | '=' ID))?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_RLS_POLICY.html
createRlsPolicy
    : 'CREATE' 'RLS' 'POLICY' ID ('WITH' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')' ('AS'? ID)?)? 'USING' '(' EXPR ')'
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_ROLE.html
createRole
    : 'CREATE' 'ROLE' ID ('EXTERNALID' ID)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_SCHEMA.html
// fix: the first form leaves [ <schema_element> open; closed
createSchema
    : 'CREATE' 'SCHEMA' ('IF' 'NOT' 'EXISTS')? ID ('AUTHORIZATION' ID)? ('QUOTA' (NUMBER ('MB' | 'GB' | 'TB')? | 'UNLIMITED'))? SCHEMA_ELEMENT+?
    | 'CREATE' 'SCHEMA' 'AUTHORIZATION' ID ('QUOTA' (NUMBER ('MB' | 'GB' | 'TB')? | 'UNLIMITED'))? SCHEMA_ELEMENT+?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_TABLE_NEW.html
// fix: table_constraints: the FOREIGN KEY line leaves its [ open, and an element is one of the three constraints the Parameters section lists, not any of them in a row; read as a choice
createTableNew
    : 'CREATE' ('LOCAL'? ('TEMPORARY' | 'TEMP'))? 'TABLE' ('IF' 'NOT' 'EXISTS')? ID '(' (ID DATA_TYPE createTableNew_column_attributes? createTableNew_column_constraints? | createTableNew_table_constraints | 'LIKE' ID (('INCLUDING' | 'EXCLUDING') 'DEFAULTS')?) (',' (ID DATA_TYPE createTableNew_column_attributes? createTableNew_column_constraints? | createTableNew_table_constraints | 'LIKE' ID (('INCLUDING' | 'EXCLUDING') 'DEFAULTS')?))* ')' ('BACKUP' ('YES' | 'NO'))? createTableNew_table_attributes?
    ;
createTableNew_column_attributes
    : ('DEFAULT' EXPR)? ('IDENTITY' '(' NUMBER ',' NUMBER ')')? ('GENERATED' 'BY' 'DEFAULT' 'AS' 'IDENTITY' '(' NUMBER ',' NUMBER ')')? ('ENCODE' ID)? 'DISTKEY'? 'SORTKEY'? ('COLLATE' ('CASE_SENSITIVE' | 'CS' | 'CASE_INSENSITIVE' | 'CI'))?
    ;
createTableNew_column_constraints
    : (('NOT' 'NULL' | 'NULL'))? (('UNIQUE' | 'PRIMARY' 'KEY'))? ('REFERENCES' ID ('(' ID ')')?)?
    ;
createTableNew_table_constraints
    : ('UNIQUE' '(' ID (',' ID)* ')' | 'PRIMARY' 'KEY' '(' ID (',' ID)* ')' | 'FOREIGN' 'KEY' '(' ID (',' ID)* ')' 'REFERENCES' ID ('(' ID ')')?)
    ;
createTableNew_table_attributes
    : ('DISTSTYLE' ('AUTO' | 'EVEN' | 'KEY' | 'ALL'))? ('DISTKEY' '(' ID ')')? (('COMPOUND' | 'INTERLEAVED')? 'SORTKEY' '(' ID (',' ID)* ')' | ('SORTKEY' 'AUTO')?)? ('ENCODE' 'AUTO')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_TABLE_AS.html
createTableAs
    : 'CREATE' ('LOCAL'? ('TEMPORARY' | 'TEMP'))? 'TABLE' ID ('(' ID (',' ID)* ')')? ('BACKUP' ('YES' | 'NO'))? createTableAs_table_attributes? 'AS' QUERY
    ;
createTableAs_table_attributes
    : ('DISTSTYLE' ('AUTO' | 'EVEN' | 'ALL' | 'KEY'))? ('DISTKEY' '(' ID ')')? (('COMPOUND' | 'INTERLEAVED')? 'SORTKEY' '(' ID (',' ID)* ')')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_TEMPLATE.html
// fix: COPY is in italics, but the command's keyword
createTemplate
    : 'CREATE' ('OR' 'REPLACE')? 'TEMPLATE' (ID '.')? (ID '.')? ID 'FOR' 'COPY' 'AS'? ('FORMAT'? 'AS'? ID)? (ID ID? (',' ID ID?)*)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_USER.html
createUser
    : 'CREATE' 'USER' ID 'WITH'? 'PASSWORD' (STRING | 'DISABLE') createUser_option+?
    ;
createUser_option
    : 'CREATEDB'
    | 'NOCREATEDB'
    | 'CREATEUSER'
    | 'NOCREATEUSER'
    | 'SYSLOG' 'ACCESS' ('RESTRICTED' | 'UNRESTRICTED')
    | 'IN' 'GROUP' ID (',' ID)*
    | 'VALID' 'UNTIL' STRING
    | 'CONNECTION' 'LIMIT' (NUMBER | 'UNLIMITED')
    | 'SESSION' 'TIMEOUT' NUMBER
    | 'EXTERNALID' ID
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_VIEW.html
createView
    : 'CREATE' ('OR' 'REPLACE')? 'VIEW' ID ('(' ID (',' ID)* ')')? 'AS' QUERY ('WITH' 'NO' 'SCHEMA' 'BINDING')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DEALLOCATE.html
deallocate
    : 'DEALLOCATE' 'PREPARE'? ID
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/declare.html
declare
    : 'DECLARE' ID 'CURSOR' 'FOR' QUERY
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DELETE.html
// fix: the WITH clause is the statement's, on the line before it; joined
// fix: USING takes a list of tables, not a repeated USING
delete
    : ('WITH' 'RECURSIVE'? CTE (',' CTE)*)? 'DELETE' 'FROM'? ID ('USING' ID (',' ID)*)? ('WHERE' EXPR)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DESC_DATASHARE.html
descDatashare
    : 'DESC' 'DATASHARE' ID ('OF' ('ACCOUNT' ID)? 'NAMESPACE' ID)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DESC_IDENTITY_PROVIDER.html
descIdentityProvider
    : 'DESC' 'IDENTITY' 'PROVIDER' ID
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DETACH_MASKING_POLICY.html
detachMaskingPolicy
    : 'DETACH' 'MASKING' 'POLICY' ID 'ON' ID '(' ID ')' 'FROM' (ID | 'ROLE' ID | 'PUBLIC')
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DETACH_RLS_POLICY.html
detachRlsPolicy
    : 'DETACH' 'RLS' 'POLICY' ID 'ON' 'TABLE'? ID (',' ID)* 'FROM' (ID | 'ROLE' ID | 'PUBLIC') (',' (ID | 'ROLE' ID | 'PUBLIC'))*
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_DATABASE.html
dropDatabase
    : 'DROP' 'DATABASE' ID 'FORCE'?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_DATASHARE.html
dropDatashare
    : 'DROP' 'DATASHARE' ID
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_EXTERNAL_VIEW.html
dropExternalView
    : 'DROP' 'EXTERNAL' 'VIEW' ID ('IF' 'EXISTS')? ID
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_FUNCTION.html
dropFunction
    : 'DROP' 'FUNCTION' ID '(' ID? DATA_TYPE (',' DATA_TYPE)* ')' ('CASCADE' | 'RESTRICT')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_GROUP.html
dropGroup
    : 'DROP' 'GROUP' ID
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_IDENTITY_PROVIDER.html
dropIdentityProvider
    : 'DROP' 'IDENTITY' 'PROVIDER' ID 'CASCADE'?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_LIBRARY.html
dropLibrary
    : 'DROP' 'LIBRARY' ID
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_MASKING_POLICY.html
dropMaskingPolicy
    : 'DROP' 'MASKING' 'POLICY' ID
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_MODEL.html
dropModel
    : 'DROP' 'MODEL' ('IF' 'EXISTS')? ID
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/materialized-view-drop-sql-command.html
materializedViewDropSqlCommand
    : 'DROP' 'MATERIALIZED' 'VIEW' ('IF' 'EXISTS')? ID (',' ID)* ('CASCADE' | 'RESTRICT')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_PROCEDURE.html
dropProcedure
    : 'DROP' 'PROCEDURE' ID '(' (ID? ARGMODE? DATA_TYPE (',' DATA_TYPE)*)? ')'
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_RLS_POLICY.html
dropRlsPolicy
    : 'DROP' 'RLS' 'POLICY' ('IF' 'EXISTS')? ID ('CASCADE' | 'RESTRICT')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_ROLE.html
dropRole
    : 'DROP' 'ROLE' ID ('FORCE' | 'RESTRICT')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_SCHEMA.html
dropSchema
    : 'DROP' 'SCHEMA' ('IF' 'EXISTS')? ID (',' ID)* ('DROP' 'EXTERNAL' 'DATABASE')? ('CASCADE' | 'RESTRICT')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_TABLE.html
dropTable
    : 'DROP' 'TABLE' ('IF' 'EXISTS')? ID (',' ID)* ('CASCADE' | 'RESTRICT')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_TEMPLATE.html
dropTemplate
    : 'DROP' 'TEMPLATE' (ID '.')? (ID '.')? ID
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_USER.html
dropUser
    : 'DROP' 'USER' ('IF' 'EXISTS')? ID (',' ID)*
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_VIEW.html
dropView
    : 'DROP' 'VIEW' ('IF' 'EXISTS')? ID (',' ID)* ('CASCADE' | 'RESTRICT')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_END.html
end
    : 'END' ('WORK' | 'TRANSACTION')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_EXECUTE.html
execute
    : 'EXECUTE' ID ('(' ID (',' ID)* ')')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_EXPLAIN.html
explain
    : 'EXPLAIN' 'VERBOSE'? QUERY
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/fetch.html
fetch
    : 'FETCH' ('NEXT' | 'ALL' | 'FORWARD' (NUMBER | 'ALL')?)? 'FROM' ID
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_GRANT.html
// fix: FOR TABLES IN form: two } too many in its privileges; dropped
// fix: FOR SCHEMAS IN form: its privileges close a { they never open; opened
// fix: system privileges: no | after ACCESS SYSTEM TABLE; added
// fix: system privileges: no | after TRUNCATE TABLE; added
// fix: system privileges: ALL PRIVILEGES is chosen outside the braces that hold the choice; moved inside
// fix: datashare usage: the choice after TO is not braced
// fix: an example, not syntax; dropped
// fix: an example, not syntax; dropped
// fix: CONNECT: the choice after TO is not braced
grant
    : 'GRANT' (('SELECT' | 'INSERT' | 'UPDATE' | 'DELETE' | 'DROP' | 'REFERENCES' | 'ALTER' | 'TRUNCATE') (',' ('SELECT' | 'INSERT' | 'UPDATE' | 'DELETE' | 'DROP' | 'REFERENCES' | 'ALTER' | 'TRUNCATE'))* | 'ALL' 'PRIVILEGES'?) 'ON' ('TABLE'? ID (',' ID)* | 'ALL' 'TABLES' 'IN' 'SCHEMA' ID (',' ID)*) 'TO' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))*
    | 'GRANT' (('CREATE' | 'USAGE' | 'TEMPORARY' | 'TEMP' | 'ALTER') (',' ('CREATE' | 'USAGE' | 'TEMPORARY' | 'TEMP' | 'ALTER'))* | 'ALL' 'PRIVILEGES'?) 'ON' 'DATABASE' ID (',' ID)* 'TO' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))*
    | 'GRANT' (('CREATE' | 'USAGE' | 'ALTER' | 'DROP') (',' ('CREATE' | 'USAGE' | 'ALTER' | 'DROP'))* | 'ALL' 'PRIVILEGES'?) 'ON' 'SCHEMA' ID (',' ID)* 'TO' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))*
    | 'GRANT' ('EXECUTE' | 'ALL' 'PRIVILEGES'?) 'ON' ('FUNCTION' ID '(' (ID? DATA_TYPE (',' DATA_TYPE)*)? ')' (',' ID '(' (ID? DATA_TYPE (',' DATA_TYPE)*)? ')')* | 'ALL' 'FUNCTIONS' 'IN' 'SCHEMA' ID (',' ID)*) 'TO' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))*
    | 'GRANT' ('EXECUTE' | 'ALL' 'PRIVILEGES'?) 'ON' ('PROCEDURE' ID '(' (ID? DATA_TYPE (',' DATA_TYPE)*)? ')' (',' ID '(' (ID? DATA_TYPE (',' DATA_TYPE)*)? ')')* | 'ALL' 'PROCEDURES' 'IN' 'SCHEMA' ID (',' ID)*) 'TO' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))*
    | 'GRANT' 'USAGE' 'ON' 'LANGUAGE' ID (',' ID)* 'TO' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))*
    | 'GRANT' (('ALTER' | 'DROP') (',' ('ALTER' | 'DROP'))* | 'ALL' 'PRIVILEGES'?) 'ON' 'COPY' 'JOB' ID (',' ID)* 'TO' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))*
    | 'GRANT' (('ALTER' | 'DROP' | 'USAGE') (',' ('ALTER' | 'DROP' | 'USAGE'))* | 'ALL' 'PRIVILEGES'?) 'ON' 'TEMPLATE' (ID '.')? (ID '.')? ID (',' ID)* 'TO' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))*
    | 'GRANT' (('SELECT' | 'UPDATE') '(' ID (',' ID)* ')' (',' '(' ID (',' ID)* ')')* | 'ALL' 'PRIVILEGES'? '(' ID (',' ID)* ')') 'ON' 'TABLE'? ID (',' ID)* 'TO' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))*
    | 'GRANT' 'ASSUMEROLE' 'ON' (STRING (',' STRING)* | 'DEFAULT' | 'ALL') 'TO' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))* 'FOR' ('ALL' | 'COPY' | 'UNLOAD' | 'EXTERNAL' 'FUNCTION' | 'CREATE' 'MODEL') (',' ('ALL' | 'COPY' | 'UNLOAD' | 'EXTERNAL' 'FUNCTION' | 'CREATE' 'MODEL'))*
    | 'GRANT' ('SELECT' | 'ALL' 'PRIVILEGES'?) '(' ID ')' 'ON' 'EXTERNAL' 'TABLE' ID 'TO' 'IAM_ROLE' ID (',' 'IAM_ROLE' ID)* ('WITH' 'GRANT' 'OPTION')?
    | 'GRANT' (('SELECT' | 'ALTER' | 'DROP' | 'DELETE' | 'INSERT') (',' ('SELECT' | 'ALTER' | 'DROP' | 'DELETE' | 'INSERT'))* | 'ALL' 'PRIVILEGES'?) 'ON' 'EXTERNAL' 'TABLE' ID (',' ID)* 'TO' ('IAM_ROLE' ID (',' 'IAM_ROLE' ID)* | 'PUBLIC') ('WITH' 'GRANT' 'OPTION')?
    | 'GRANT' (('CREATE' | 'ALTER' | 'DROP') (',' ('CREATE' | 'ALTER' | 'DROP'))* | 'ALL' 'PRIVILEGES'?) 'ON' 'EXTERNAL' 'SCHEMA' ID (',' ID)* 'TO' 'IAM_ROLE' ID (',' 'IAM_ROLE' ID)* ('WITH' 'GRANT' 'OPTION')?
    | 'GRANT' ('ALTER' | 'SHARE') 'ON' 'DATASHARE' ID 'TO' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))*
    | 'GRANT' 'USAGE' 'ON' 'DATASHARE' ID 'TO' ('NAMESPACE' STRING | 'ACCOUNT' STRING ('VIA' 'DATA' 'CATALOG')?)
    | 'GRANT' 'EXECUTE' 'ON' 'PROCEDURE' ID '(' (ID? DATA_TYPE (',' DATA_TYPE)*)? ')' 'TO' 'DATASHARE' ID
    | 'GRANT' 'USAGE' 'ON' ('DATABASE' ID (',' ID)* | 'SCHEMA' ID) 'TO' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))*
    | 'GRANT' (('CREATE' | 'USAGE' | 'ALTER' | 'DROP') (',' ('CREATE' | 'USAGE' | 'ALTER' | 'DROP'))* | 'ALL' 'PRIVILEGES'?) 'FOR' 'SCHEMAS' 'IN' 'DATABASE' ID 'TO' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID) (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID))*
    | 'GRANT' (('SELECT' | 'INSERT' | 'UPDATE' | 'DELETE' | 'DROP' | 'ALTER' | 'TRUNCATE' | 'REFERENCES') (',' ('SELECT' | 'INSERT' | 'UPDATE' | 'DELETE' | 'DROP' | 'ALTER' | 'TRUNCATE' | 'REFERENCES'))* | 'ALL' 'PRIVILEGES'?) 'FOR' 'TABLES' 'IN' ('SCHEMA' ID ('DATABASE' ID)? | 'DATABASE' ID) 'TO' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID) (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID))*
    | 'GRANT' ('EXECUTE' | 'ALL' 'PRIVILEGES'?) 'FOR' 'FUNCTIONS' 'IN' ('SCHEMA' ID ('DATABASE' ID)? | 'DATABASE' ID) 'TO' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID) (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID))*
    | 'GRANT' ('EXECUTE' | 'ALL' 'PRIVILEGES'?) 'FOR' 'PROCEDURES' 'IN' ('SCHEMA' ID ('DATABASE' ID)? | 'DATABASE' ID) 'TO' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID) (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID))*
    | 'GRANT' 'USAGE' 'FOR' 'LANGUAGES' 'IN' 'DATABASE' ID 'TO' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID) (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID))*
    | 'GRANT' (('CREATE' | 'ALTER' | 'DROP') (',' ('CREATE' | 'ALTER' | 'DROP'))* | 'ALL' 'PRIVILEGES'?) 'FOR' 'COPY' 'JOBS' 'IN' 'DATABASE' ID 'TO' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID) (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID))*
    | 'GRANT' (('ALTER' | 'DROP' | 'USAGE') (',' ('ALTER' | 'DROP' | 'USAGE'))* | 'ALL' 'PRIVILEGES'?) 'FOR' 'TEMPLATES' 'IN' ('SCHEMA' ID ('DATABASE' ID)? | 'DATABASE' ID) 'TO' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID) (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID))*
    | 'GRANT' 'CREATE' 'MODEL' 'TO' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))*
    | 'GRANT' ('EXECUTE' | 'ALL' 'PRIVILEGES'?) 'ON' 'MODEL' ID (',' ID)* 'TO' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))*
    | 'GRANT' 'ROLE' ID (',' 'ROLE' ID)* 'TO' (ID ('WITH' 'ADMIN' 'OPTION')? | 'ROLE' ID) (',' (ID ('WITH' 'ADMIN' 'OPTION')? | 'ROLE' ID))*
    | 'GRANT' (('CREATE' 'USER' | 'DROP' 'USER' | 'ALTER' 'USER' | 'CREATE' 'SCHEMA' | 'DROP' 'SCHEMA' | 'ALTER' 'DEFAULT' 'PRIVILEGES' | 'ACCESS' 'CATALOG' | 'ACCESS' 'SYSTEM' 'TABLE' | 'CREATE' 'TABLE' | 'DROP' 'TABLE' | 'ALTER' 'TABLE' | 'CREATE' 'OR' 'REPLACE' 'FUNCTION' | 'CREATE' 'OR' 'REPLACE' 'EXTERNAL' 'FUNCTION' | 'DROP' 'FUNCTION' | 'CREATE' 'OR' 'REPLACE' 'PROCEDURE' | 'DROP' 'PROCEDURE' | 'CREATE' 'OR' 'REPLACE' 'VIEW' | 'DROP' 'VIEW' | 'CREATE' 'MODEL' | 'DROP' 'MODEL' | 'CREATE' 'DATASHARE' | 'ALTER' 'DATASHARE' | 'DROP' 'DATASHARE' | 'CREATE' 'LIBRARY' | 'DROP' 'LIBRARY' | 'CREATE' 'ROLE' | 'DROP' 'ROLE' | 'TRUNCATE' 'TABLE' | 'VACUUM' | 'ANALYZE' | 'CANCEL' | 'IGNORE' 'RLS' | 'EXPLAIN' 'RLS' | 'EXPLAIN' 'MASKING') (',' ('CREATE' 'USER' | 'DROP' 'USER' | 'ALTER' 'USER' | 'CREATE' 'SCHEMA' | 'DROP' 'SCHEMA' | 'ALTER' 'DEFAULT' 'PRIVILEGES' | 'ACCESS' 'CATALOG' | 'ACCESS' 'SYSTEM' 'TABLE' | 'CREATE' 'TABLE' | 'DROP' 'TABLE' | 'ALTER' 'TABLE' | 'CREATE' 'OR' 'REPLACE' 'FUNCTION' | 'CREATE' 'OR' 'REPLACE' 'EXTERNAL' 'FUNCTION' | 'DROP' 'FUNCTION' | 'CREATE' 'OR' 'REPLACE' 'PROCEDURE' | 'DROP' 'PROCEDURE' | 'CREATE' 'OR' 'REPLACE' 'VIEW' | 'DROP' 'VIEW' | 'CREATE' 'MODEL' | 'DROP' 'MODEL' | 'CREATE' 'DATASHARE' | 'ALTER' 'DATASHARE' | 'DROP' 'DATASHARE' | 'CREATE' 'LIBRARY' | 'DROP' 'LIBRARY' | 'CREATE' 'ROLE' | 'DROP' 'ROLE' | 'TRUNCATE' 'TABLE' | 'VACUUM' | 'ANALYZE' | 'CANCEL' | 'IGNORE' 'RLS' | 'EXPLAIN' 'RLS' | 'EXPLAIN' 'MASKING'))* | 'ALL' 'PRIVILEGES'?) 'TO' 'ROLE' ID (',' ID)*
    | 'GRANT' 'EXPLAIN' ('RLS' | 'MASKING') 'TO' 'ROLE' ID
    | 'GRANT' 'IGNORE' 'RLS' 'TO' 'ROLE' ID
    | 'GRANT' 'SELECT' 'ON' 'TABLE'? ID (',' ID)* 'TO' ('RLS' | 'MASKING') 'POLICY' ID (',' ID)*
    | 'GRANT' 'CONNECT' ('ON' 'WORKGROUP')? 'TO' (ID ':' ID | 'ROLE' ID ':' ID | 'PUBLIC')
    | 'GRANT' 'DEBUG' 'ON' 'DATABASE' ID 'TO' (ID ':' ID | 'ROLE' ID ':' ID | 'ADMINS' 'OF' 'ACCOUNT' STRING)
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_INSERT_30.html
insert30
    : 'INSERT' 'INTO' ID ('(' ID (',' ID)* ')')? ('DEFAULT' 'VALUES' | 'VALUES' '(' (EXPR | 'DEFAULT') (',' (EXPR | 'DEFAULT'))* ')' (',' '(' (EXPR | 'DEFAULT') (',' (EXPR | 'DEFAULT'))* ')')* | QUERY)
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_INSERT_external_table.html
insertExternalTable
    : 'INSERT' 'INTO' ID QUERY
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_LOCK.html
lock
    : 'LOCK' 'TABLE'? ID (',' ID)*
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_MERGE.html
// fix: SET repeats the whole assignment
merge
    : 'MERGE' 'INTO' ID 'USING' ID ('AS'? ID)? 'ON' EXPR ('WHEN' 'MATCHED' 'THEN' ('UPDATE' 'SET' ID '=' EXPR (',' ID '=' EXPR)* | 'DELETE') 'WHEN' 'NOT' 'MATCHED' 'THEN' 'INSERT' ('(' ID (',' ID)* ')')? 'VALUES' '(' EXPR (',' EXPR)* ')' | 'REMOVE' 'DUPLICATES')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_PREPARE.html
prepare
    : 'PREPARE' ID ('(' DATA_TYPE (',' DATA_TYPE)* ')')? 'AS' QUERY
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/materialized-view-refresh-sql-command.html
materializedViewRefreshSqlCommand
    : 'REFRESH' 'MATERIALIZED' 'VIEW' ID ('RESTRICT' | 'CASCADE')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_RESET.html
reset
    : 'RESET' (ID | 'ALL')
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_REVOKE.html
// fix: FOR TABLES IN form: two } too many in its privileges; dropped
// fix: FOR SCHEMAS IN form: its privileges close a { they never open; opened
// fix: system privileges: no | after DROP ROLE and TRUNCATE TABLE; added
// fix: system privileges: ALL PRIVILEGES is chosen outside the braces that hold the choice; moved inside
// fix: datashare usage: the choice after FROM is not braced
// fix: an example, not syntax; dropped
// fix: an example, not syntax; dropped
revoke
    : 'REVOKE' ('GRANT' 'OPTION' 'FOR')? (('SELECT' | 'INSERT' | 'UPDATE' | 'DELETE' | 'DROP' | 'REFERENCES' | 'ALTER' | 'TRUNCATE') (',' ('SELECT' | 'INSERT' | 'UPDATE' | 'DELETE' | 'DROP' | 'REFERENCES' | 'ALTER' | 'TRUNCATE'))* | 'ALL' 'PRIVILEGES'?) 'ON' ('TABLE'? ID (',' ID)* | 'ALL' 'TABLES' 'IN' 'SCHEMA' ID (',' ID)*) 'FROM' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))* 'RESTRICT'?
    | 'REVOKE' ('GRANT' 'OPTION' 'FOR')? (('CREATE' | 'TEMPORARY' | 'TEMP' | 'ALTER') (',' ('CREATE' | 'TEMPORARY' | 'TEMP' | 'ALTER'))* | 'ALL' 'PRIVILEGES'?) 'ON' 'DATABASE' ID (',' ID)* 'FROM' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))* 'RESTRICT'?
    | 'REVOKE' ('GRANT' 'OPTION' 'FOR')? (('CREATE' | 'USAGE' | 'ALTER' | 'DROP') (',' ('CREATE' | 'USAGE' | 'ALTER' | 'DROP'))* | 'ALL' 'PRIVILEGES'?) 'ON' 'SCHEMA' ID (',' ID)* 'FROM' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))* 'RESTRICT'?
    | 'REVOKE' ('GRANT' 'OPTION' 'FOR')? 'EXECUTE' 'ON' 'FUNCTION' ID '(' (ID? DATA_TYPE (',' DATA_TYPE)*)? ')' (',' ID '(' (ID? DATA_TYPE (',' DATA_TYPE)*)? ')')* 'FROM' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))* 'RESTRICT'?
    | 'REVOKE' ('GRANT' 'OPTION' 'FOR')? ('EXECUTE' (',' 'EXECUTE')* | 'ALL' 'PRIVILEGES'?) 'ON' 'PROCEDURE' ID '(' (ID? DATA_TYPE (',' DATA_TYPE)*)? ')' (',' ID '(' (ID? DATA_TYPE (',' DATA_TYPE)*)? ')')* 'FROM' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))* 'RESTRICT'?
    | 'REVOKE' ('GRANT' 'OPTION' 'FOR')? 'USAGE' 'ON' 'LANGUAGE' ID (',' ID)* 'FROM' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))* 'RESTRICT'?
    | 'REVOKE' ('GRANT' 'OPTION' 'FOR')? (('ALTER' | 'DROP') (',' ('ALTER' | 'DROP'))* | 'ALL' 'PRIVILEGES'?) 'ON' 'COPY' 'JOB' ID (',' ID)* 'FROM' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))*
    | 'REVOKE' ('GRANT' 'OPTION' 'FOR')? (('ALTER' | 'DROP' | 'USAGE') (',' ('ALTER' | 'DROP' | 'USAGE'))* | 'ALL' 'PRIVILEGES'?) 'ON' 'TEMPLATE' ID (',' ID)* 'FROM' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))*
    | 'REVOKE' (('SELECT' | 'UPDATE') '(' ID (',' ID)* ')' (',' '(' ID (',' ID)* ')')* | 'ALL' 'PRIVILEGES'? '(' ID (',' ID)* ')') 'ON' 'TABLE'? ID (',' ID)* 'FROM' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))* 'RESTRICT'?
    | 'REVOKE' 'ASSUMEROLE' 'ON' (STRING (',' STRING)* | 'DEFAULT' | 'ALL') 'FROM' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))* 'FOR' ('ALL' | 'COPY' | 'UNLOAD' | 'EXTERNAL' 'FUNCTION' | 'CREATE' 'MODEL')
    | 'REVOKE' ('GRANT' 'OPTION' 'FOR')? ('SELECT' | 'ALL' 'PRIVILEGES'?) '(' ID ')' 'ON' 'EXTERNAL' 'TABLE' ID 'FROM' 'IAM_ROLE' ID (',' 'IAM_ROLE' ID)*
    | 'REVOKE' ('GRANT' 'OPTION' 'FOR')? (('SELECT' | 'ALTER' | 'DROP' | 'DELETE' | 'INSERT') (',' ('SELECT' | 'ALTER' | 'DROP' | 'DELETE' | 'INSERT'))* | 'ALL' 'PRIVILEGES'?) 'ON' 'EXTERNAL' 'TABLE' ID (',' ID)* 'FROM' ('IAM_ROLE' ID (',' 'IAM_ROLE' ID)* | 'PUBLIC')
    | 'REVOKE' ('GRANT' 'OPTION' 'FOR')? (('CREATE' | 'ALTER' | 'DROP') (',' ('CREATE' | 'ALTER' | 'DROP'))* | 'ALL' 'PRIVILEGES'?) 'ON' 'EXTERNAL' 'SCHEMA' ID (',' ID)* 'FROM' 'IAM_ROLE' ID (',' 'IAM_ROLE' ID)*
    | 'REVOKE' ('ALTER' | 'SHARE') 'ON' 'DATASHARE' ID 'FROM' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))*
    | 'REVOKE' 'USAGE' 'ON' 'DATASHARE' ID 'FROM' ('NAMESPACE' STRING (',' STRING)* | 'ACCOUNT' STRING ('VIA' 'DATA' 'CATALOG')? (',' STRING ('VIA' 'DATA' 'CATALOG')?)*)
    | 'REVOKE' 'EXECUTE' 'ON' 'PROCEDURE' ID '(' (ID? DATA_TYPE (',' DATA_TYPE)*)? ')' 'FROM' 'DATASHARE' ID
    | 'REVOKE' 'USAGE' 'ON' ('DATABASE' ID (',' ID)* | 'SCHEMA' ID) 'FROM' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))*
    | 'REVOKE' ('GRANT' 'OPTION')? (('CREATE' | 'USAGE' | 'ALTER' | 'DROP') (',' ('CREATE' | 'USAGE' | 'ALTER' | 'DROP'))* | 'ALL' 'PRIVILEGES'?) 'FOR' 'SCHEMAS' 'IN' 'DATABASE' ID 'FROM' (ID | 'ROLE' ID) (',' (ID | 'ROLE' ID))*
    | 'REVOKE' ('GRANT' 'OPTION')? (('SELECT' | 'INSERT' | 'UPDATE' | 'DELETE' | 'DROP' | 'ALTER' | 'TRUNCATE' | 'REFERENCES') (',' ('SELECT' | 'INSERT' | 'UPDATE' | 'DELETE' | 'DROP' | 'ALTER' | 'TRUNCATE' | 'REFERENCES'))* | 'ALL' 'PRIVILEGES'?) 'FOR' 'TABLES' 'IN' ('SCHEMA' ID ('DATABASE' ID)? | 'DATABASE' ID) 'FROM' (ID | 'ROLE' ID) (',' (ID | 'ROLE' ID))*
    | 'REVOKE' ('GRANT' 'OPTION')? ('EXECUTE' | 'ALL' 'PRIVILEGES'?) 'FOR' 'FUNCTIONS' 'IN' ('SCHEMA' ID ('DATABASE' ID)? | 'DATABASE' ID) 'FROM' (ID | 'ROLE' ID) (',' (ID | 'ROLE' ID))*
    | 'REVOKE' ('GRANT' 'OPTION')? ('EXECUTE' | 'ALL' 'PRIVILEGES'?) 'FOR' 'PROCEDURES' 'IN' ('SCHEMA' ID ('DATABASE' ID)? | 'DATABASE' ID) 'FROM' (ID | 'ROLE' ID) (',' (ID | 'ROLE' ID))*
    | 'REVOKE' ('GRANT' 'OPTION')? 'USAGE' 'FOR' 'LANGUAGES' 'IN' 'DATABASE' ID 'FROM' (ID | 'ROLE' ID) (',' (ID | 'ROLE' ID))*
    | 'REVOKE' 'GRANT_OPTION'? (('CREATE' | 'ALTER' | 'DROP') (',' ('CREATE' | 'ALTER' | 'DROP'))* | 'ALL' 'PRIVILEGES'?) 'FOR' 'COPY' 'JOBS' 'IN' 'DATABASE' ID 'FROM' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID) (',' (ID ('WITH' 'GRANT' 'OPTION')? | 'ROLE' ID))*
    | 'REVOKE' ('GRANT' 'OPTION')? (('ALTER' | 'DROP' | 'USAGE') (',' ('ALTER' | 'DROP' | 'USAGE'))* | 'ALL' 'PRIVILEGES'?) 'FOR' 'TEMPLATES' 'IN' ('SCHEMA' ID ('DATABASE' ID)? | 'DATABASE' ID) 'FROM' (ID | 'ROLE' ID) (',' (ID | 'ROLE' ID))*
    | 'REVOKE' ('GRANT' 'OPTION' 'FOR')? 'CREATE' 'MODEL' 'FROM' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))* 'RESTRICT'?
    | 'REVOKE' ('GRANT' 'OPTION' 'FOR')? ('EXECUTE' | 'ALL' 'PRIVILEGES'?) 'ON' 'MODEL' ID (',' ID)* 'FROM' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC') (',' (ID | 'ROLE' ID | 'GROUP' ID | 'PUBLIC'))* 'RESTRICT'?
    | 'REVOKE' ('ADMIN' 'OPTION' 'FOR')? 'ROLE' ID (',' 'ROLE' ID)* 'FROM' ID (',' ID)*
    | 'REVOKE' 'ROLE' ID (',' 'ROLE' ID)* 'FROM' 'ROLE' ID (',' 'ROLE' ID)*
    | 'REVOKE' (('CREATE' 'USER' | 'DROP' 'USER' | 'ALTER' 'USER' | 'CREATE' 'SCHEMA' | 'DROP' 'SCHEMA' | 'ALTER' 'DEFAULT' 'PRIVILEGES' | 'ACCESS' 'CATALOG' | 'CREATE' 'TABLE' | 'DROP' 'TABLE' | 'ALTER' 'TABLE' | 'CREATE' 'OR' 'REPLACE' 'FUNCTION' | 'CREATE' 'OR' 'REPLACE' 'EXTERNAL' 'FUNCTION' | 'DROP' 'FUNCTION' | 'CREATE' 'OR' 'REPLACE' 'PROCEDURE' | 'DROP' 'PROCEDURE' | 'CREATE' 'OR' 'REPLACE' 'VIEW' | 'DROP' 'VIEW' | 'CREATE' 'MODEL' | 'DROP' 'MODEL' | 'CREATE' 'DATASHARE' | 'ALTER' 'DATASHARE' | 'DROP' 'DATASHARE' | 'CREATE' 'LIBRARY' | 'DROP' 'LIBRARY' | 'CREATE' 'ROLE' | 'DROP' 'ROLE' | 'TRUNCATE' 'TABLE' | 'VACUUM' | 'ANALYZE' | 'CANCEL') (',' ('CREATE' 'USER' | 'DROP' 'USER' | 'ALTER' 'USER' | 'CREATE' 'SCHEMA' | 'DROP' 'SCHEMA' | 'ALTER' 'DEFAULT' 'PRIVILEGES' | 'ACCESS' 'CATALOG' | 'CREATE' 'TABLE' | 'DROP' 'TABLE' | 'ALTER' 'TABLE' | 'CREATE' 'OR' 'REPLACE' 'FUNCTION' | 'CREATE' 'OR' 'REPLACE' 'EXTERNAL' 'FUNCTION' | 'DROP' 'FUNCTION' | 'CREATE' 'OR' 'REPLACE' 'PROCEDURE' | 'DROP' 'PROCEDURE' | 'CREATE' 'OR' 'REPLACE' 'VIEW' | 'DROP' 'VIEW' | 'CREATE' 'MODEL' | 'DROP' 'MODEL' | 'CREATE' 'DATASHARE' | 'ALTER' 'DATASHARE' | 'DROP' 'DATASHARE' | 'CREATE' 'LIBRARY' | 'DROP' 'LIBRARY' | 'CREATE' 'ROLE' | 'DROP' 'ROLE' | 'TRUNCATE' 'TABLE' | 'VACUUM' | 'ANALYZE' | 'CANCEL'))* | 'ALL' 'PRIVILEGES'?) 'FROM' 'ROLE' ID (',' 'ROLE' ID)*
    | 'REVOKE' 'EXPLAIN' ('RLS' | 'MASKING') 'FROM' 'ROLE' ID
    | 'REVOKE' 'IGNORE' 'RLS' 'FROM' 'ROLE' ID
    | 'REVOKE' 'SELECT' 'ON' 'TABLE'? ID (',' ID)* 'FROM' ('RLS' | 'MASKING') 'POLICY' ID (',' ID)*
    | 'REVOKE' 'DEBUG' 'ON' 'DATABASE' ID 'FROM' (ID ':' ID | 'ROLE' ID ':' ID | 'ADMINS' 'OF' 'ACCOUNT' STRING)
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_ROLLBACK.html
rollback
    : 'ROLLBACK' ('WORK' | 'TRANSACTION')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SELECT_synopsis.html
// fix: the WITH clause is the statement's, on the line before it; joined
// fix: the select list is bracketed with TOP, ALL and DISTINCT as if optional; read as r_SELECT_list.html writes it, TOP or ALL or DISTINCT before a required list
// fix: GROUP BY takes ALL or a list; the choice is braced
// fix: ALL is UNION ALL's, as r_UNION.html writes it
selectSynopsis
    : ('WITH' CTE (',' CTE)*)? 'SELECT' ('TOP' NUMBER | 'ALL' | 'DISTINCT')? ('*' | EXPR ('AS' ID)? (',' EXPR ('AS' ID)?)*) ('EXCLUDE' ID)? ('FROM' ID (',' ID)*)? ('WHERE' EXPR)? (('START' 'WITH' EXPR)? 'CONNECT' 'BY' EXPR)? ('GROUP' 'BY' ('ALL' | EXPR (',' EXPR)*))? ('HAVING' EXPR)? ('QUALIFY' EXPR)? (('UNION' 'ALL'? | 'INTERSECT' | 'EXCEPT' | 'MINUS') QUERY)? ('ORDER' 'BY' EXPR ('ASC' | 'DESC')?)? ('LIMIT' (NUMBER | 'ALL'))? ('OFFSET' NUMBER)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SELECT_INTO.html
// fix: the WITH clause is the statement's, on the line before it; joined
// fix: the select list is bracketed with TOP, ALL and DISTINCT as if optional; read as r_SELECT_list.html writes it, TOP or ALL or DISTINCT before a required list
// fix: GROUP BY takes ALL or a list; the choice is braced
// fix: ALL is UNION ALL's, as r_UNION.html writes it
selectInto
    : ('WITH' CTE (',' CTE)*)? 'SELECT' ('TOP' NUMBER | 'ALL' | 'DISTINCT')? ('*' | EXPR ('AS' ID)? (',' EXPR ('AS' ID)?)*) ('EXCLUDE' ID)? 'INTO' ('TEMPORARY' | 'TEMP')? 'TABLE'? ID ('FROM' ID (',' ID)*)? ('WHERE' EXPR)? (('START' 'WITH' EXPR)? 'CONNECT' 'BY' EXPR)? ('GROUP' 'BY' ('ALL' | EXPR (',' EXPR)*))? ('HAVING' EXPR)? ('QUALIFY' EXPR)? (('UNION' 'ALL'? | 'INTERSECT' | 'EXCEPT' | 'MINUS') QUERY)? ('ORDER' 'BY' EXPR ('ASC' | 'DESC')?)? ('LIMIT' (NUMBER | 'ALL'))? ('OFFSET' NUMBER)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SET.html
// fix: the variable form opens a { it never closes; dropped
set
    : 'SET' (('SESSION' | 'LOCAL')? ('SEED' | ID) ('TO' | '=') (ID | STRING | 'DEFAULT') | 'SEED' 'TO' ID)
    | 'SET' ('SESSION' | 'LOCAL')? ID ('TO' | '=') (ID | STRING)
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SET_SESSION_AUTHORIZATION.html
setSessionAuthorization
    : 'SET' 'LOCAL'? 'SESSION' 'AUTHORIZATION' (ID | 'DEFAULT')
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SET_SESSION_CHARACTERISTICS.html
// no syntax

// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW.html
show
    : 'SHOW' (ID | 'ALL')
    | 'SHOW' ID
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_USER_LOCKOUT.html
showUserLockout
    : 'SHOW' 'USER' 'LOCKOUT'
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_COLUMN_GRANTS.html
showColumnGrants
    : 'SHOW' 'COLUMN' 'GRANTS' 'ON' 'TABLE' (ID '.' ID '.' ID | ID '.' ID) ('FOR' (ID | 'ROLE' ID | 'PUBLIC'))? ('LIMIT' NUMBER)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_COLUMNS.html
showColumns
    : 'SHOW' 'COLUMNS' 'FROM' 'TABLE' ID '.' ID '.' ID ('LIKE' STRING)? ('LIMIT' NUMBER)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_CONSTRAINTS.html
showConstraints
    : 'SHOW' 'CONSTRAINTS' ('PRIMARY' 'KEYS' | 'FOREIGN' 'KEYS' 'EXPORTED'?) 'FROM' 'TABLE' (ID '.' ID '.' ID | ID '.' ID) ('LIMIT' NUMBER)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_EXTERNAL_TABLE.html
showExternalTable
    : 'SHOW' 'EXTERNAL' 'TABLE' ID? '.' ID '.' ID 'PARTITION'?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_DATABASES.html
// fix: IAM_ROLE takes one of three; the choice is braced
showDatabases
    : 'SHOW' 'DATABASES' ('LIKE' STRING)? ('LIMIT' NUMBER)?
    | 'SHOW' 'DATABASES' 'FROM' 'DATA' 'CATALOG' ('ACCOUNT' STRING ',' STRING (',' 'ACCOUNT' STRING ',' STRING)*)? ('LIKE' STRING)? ('IAM_ROLE' ('DEFAULT' | '\'SESSION\'' | STRING))? ('LIMIT' NUMBER)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_FUNCTIONS.html
showFunctions
    : 'SHOW' 'FUNCTIONS' 'FROM' 'SCHEMA' (ID '.')? ID ('LIKE' STRING)? ('LIMIT' NUMBER)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_GRANTS.html
// fix: no | before TEMPLATE; added
showGrants
    : 'SHOW' 'GRANTS' 'ON' ('DATABASE' ID | 'FUNCTION' ID '(' (ID? DATA_TYPE (',' DATA_TYPE)*)? ')' | 'FUNCTION' ID | 'SCHEMA' ID | ('TABLE' ID | ID) | 'TEMPLATE' ID) ('FOR' (ID | 'ROLE' ID | 'PUBLIC'))? ('LIMIT' NUMBER)?
    | 'SHOW' 'GRANTS' 'FOR' (ID | 'ROLE' ID) ('FROM' 'DATABASE' ID)? ('LIMIT' NUMBER)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_MODEL.html
// fix: the choice is written in parentheses, which are not SQL here (the examples write SHOW MODEL ALL)
showModel
    : 'SHOW' 'MODEL' ('ALL' | ID)
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_DATASHARES.html
showDatashares
    : 'SHOW' 'DATASHARES' ('LIKE' STRING)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_PARAMETERS.html
showParameters
    : 'SHOW' 'PARAMETERS' 'OF' ('FUNCTION' | 'PROCEDURE') (ID '.')? ID '.' ID '(' DATA_TYPE (',' DATA_TYPE)* ')' ('LIKE' STRING)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_POLICIES.html
showPolicies
    : 'SHOW' ('RLS' | 'MASKING') 'POLICIES' ('ON' (ID '.' ID '.' ID | ID '.' ID) ('FOR' (ID | 'ROLE' ID | 'PUBLIC'))? | 'FROM' 'DATABASE' ID)? ('LIMIT' NUMBER)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_PROCEDURE.html
showProcedure
    : 'SHOW' 'PROCEDURE' ID ('(' (ID? ARGMODE? DATA_TYPE (',' DATA_TYPE)*)? ')')?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_PROCEDURES.html
showProcedures
    : 'SHOW' 'PROCEDURES' 'FROM' 'SCHEMA' (ID '.')? ID ('LIKE' STRING)? ('LIMIT' NUMBER)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_SCHEMAS.html
showSchemas
    : 'SHOW' 'SCHEMAS' 'FROM' 'DATABASE' ID ('LIKE' STRING)? ('LIMIT' NUMBER)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_TABLE.html
showTable
    : 'SHOW' 'TABLE' (ID '.')? ID
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_TABLES.html
showTables
    : 'SHOW' 'TABLES' 'FROM' 'SCHEMA' ID '.' ID ('LIKE' STRING)? ('LIMIT' NUMBER)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_TEMPLATE.html
showTemplate
    : 'SHOW' 'TEMPLATE' (ID '.')? (ID '.')? ID
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_TEMPLATES.html
showTemplates
    : 'SHOW' 'TEMPLATES' 'FROM' 'SCHEMA' (ID '.')? ID ('LIKE' STRING)? ('LIMIT' NUMBER)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_VIEW.html
showView
    : 'SHOW' 'VIEW' (ID '.')? ID
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_START_TRANSACTION.html
// no syntax

// https://docs.aws.amazon.com/redshift/latest/dg/r_TRUNCATE.html
truncate
    : 'TRUNCATE' 'TABLE'? ID
    | 'TRUNCATE' ID
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_UNLOAD.html
// fix: option: a | before the first choice, and the formats are not braced
// fix: option: a } after REGION closes no bracket; dropped
unload
    : 'UNLOAD' '(' STRING ')' 'TO' STRING unload_authorization (unload_option (',' unload_option)*)?
    ;
unload_authorization
    : 'IAM_ROLE' ('DEFAULT' | '\'SESSION\'' | STRING)
    ;
unload_option
    : ('FORMAT' 'AS'?)? ('CSV' | 'PARQUET' | 'JSON')
    | 'PARTITION' 'BY' '(' ID (',' ID)* ')' 'INCLUDE'?
    | 'MANIFEST' 'VERBOSE'?
    | 'HEADER'
    | 'DELIMITER' 'AS'? STRING
    | 'FIXEDWIDTH' 'AS'? STRING
    | 'ENCRYPTED' 'AUTO'?
    | 'BZIP2'
    | 'GZIP'
    | 'ZSTD'
    | 'ADDQUOTES'
    | 'NULL' 'AS'? STRING
    | 'ESCAPE'
    | 'ALLOWOVERWRITE'
    | 'CLEANPATH'
    | 'PARALLEL' (('ON' | 'TRUE') | ('OFF' | 'FALSE'))?
    | 'MAXFILESIZE' 'AS'? NUMBER ('MB' | 'GB')?
    | 'ROWGROUPSIZE' 'AS'? NUMBER ('MB' | 'GB')?
    | 'REGION' 'AS'? STRING
    | 'EXTENSION' STRING
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_UPDATE.html
// fix: SET repeats the whole assignment
update
    : ('WITH' 'RECURSIVE'? CTE (',' CTE)*)? 'UPDATE' ID ('AS'? ID)? 'SET' ID '=' (EXPR | 'DEFAULT') (',' ID '=' (EXPR | 'DEFAULT'))* ('FROM' ID)? ('WHERE' EXPR)?
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_USE_command.html
useCommand
    : 'USE' ID
    ;

// https://docs.aws.amazon.com/redshift/latest/dg/r_VACUUM_command.html
vacuumCommand
    : 'VACUUM' ('FULL' | 'SORT' 'ONLY' | 'DELETE' 'ONLY' | 'REINDEX' | 'RECLUSTER')? (ID? ('TO' NUMBER 'PERCENT')? 'BOOST'? 'FORCE'?)?
    ;
