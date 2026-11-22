/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// This file is derived from Apache Hive's parser/src/java/org/apache/hadoop/hive/ql/parse/
// ANTLR 3 grammar files (vendor/*.g), which are beside it as published, with the
// license and notice they come under (Hive.LICENSE, Hive.NOTICE). It was
// changed: hive/parser/cmd/antlr3to4 wrote it from them, converting ANTLR 3 to
// ANTLR 4 and keeping only what each rule derives. AST rewrite rules (-> ...),
// actions, arguments, labels, tree operators, and predicates are gone; empty
// lexer stubs overridden by HiveLexer.g are dropped; and what else was lost is
// noted beside its rule. Each rule is under the file and line it came from.
//
// Do not edit: rerun the converter.

grammar Hive;

// vendor/HiveParser.g:899
statement
    : explainStatement EOF
    | execStatement EOF
    ;

// vendor/HiveParser.g:904
explainStatement
    : KW_EXPLAIN ( explainOption* execStatement | KW_REWRITE queryStatementExpression )
    ;

// vendor/HiveParser.g:914
explainOption
    : KW_EXTENDED
    | KW_FORMATTED
    | KW_DEPENDENCY
    | KW_CBO ( KW_COST | KW_JOINCOST )?
    | KW_LOGICAL
    | KW_AUTHORIZATION
    | KW_ANALYZE
    | KW_REOPTIMIZATION
    | KW_LOCKS
    | KW_AST
    | ( KW_VECTORIZATION vectorizationOnly? vectorizatonDetail? )
    | KW_DEBUG
    | KW_DDL
    ;

// vendor/HiveParser.g:932
vectorizationOnly
    : KW_ONLY
    ;

// vendor/HiveParser.g:939
vectorizatonDetail
    : KW_SUMMARY
    | KW_OPERATOR
    | KW_EXPRESSION
    | KW_DETAIL
    ;

// vendor/HiveParser.g:952
execStatement
    : queryStatementExpression
    | loadStatement
    | exportStatement
    | importStatement
    | replDumpStatement
    | replLoadStatement
    | replStatusStatement
    | ddlStatement
    | deleteStatement
    | updateStatement
    | sqlTransactionStatement
    | mergeStatement
    | prepareStatement
    | executeStatement
    ;

// vendor/HiveParser.g:971
loadStatement
    : KW_LOAD KW_DATA ( KW_LOCAL )? KW_INPATH ( StringLiteral ) ( KW_OVERWRITE )? KW_INTO KW_TABLE ( tableOrPartition ) inputFileFormat?
    ;

// vendor/HiveParser.g:978
replicationClause
    : KW_FOR ( KW_METADATA )? KW_REPLICATION LPAREN ( StringLiteral ) RPAREN
    ;

// vendor/HiveParser.g:985
exportStatement
    : KW_EXPORT KW_TABLE ( tableOrPartition ) KW_TO ( StringLiteral ) replicationClause?
    ;

// vendor/HiveParser.g:995
importStatement
    : KW_IMPORT ( ( KW_EXTERNAL )? KW_TABLE ( tableOrPartition ) )? KW_FROM ( StringLiteral ) tableLocation?
    ;

// vendor/HiveParser.g:1005
replDumpStatement
    : KW_REPL KW_DUMP ( replDbPolicy ) ( KW_REPLACE replDbPolicy )? ( KW_WITH replConfigs )?
    ;

// vendor/HiveParser.g:1015
replDbPolicy
    : ( identifier ) ( DOT replTableLevelPolicy )?
    ;

// vendor/HiveParser.g:1022
replLoadStatement
    : KW_REPL KW_LOAD ( replDbPolicy ) ( KW_INTO identifier )? ( KW_WITH replConfigs )?
    ;

// vendor/HiveParser.g:1032
replConfigs
    : LPAREN replConfigsList RPAREN
    ;

// vendor/HiveParser.g:1039
replConfigsList
    : keyValueProperty ( COMMA keyValueProperty )*
    ;

// vendor/HiveParser.g:1046
replTableLevelPolicy
    : ( ( StringLiteral ) ( DOT StringLiteral )? )
    ;

// vendor/HiveParser.g:1054
replStatusStatement
    : KW_REPL KW_STATUS ( identifier ) ( KW_WITH replConfigs )?
    ;

// vendor/HiveParser.g:1063
ddlStatement
    : createDatabaseStatement
    | switchDatabaseStatement
    | dropDatabaseStatement
    | createTableStatement
    | dropTableStatement
    | truncateTableStatement
    | alterStatement
    | descStatement
    | showStatement
    | metastoreCheck
    | createViewStatement
    | createMaterializedViewStatement
    | createScheduledQueryStatement
    | alterScheduledQueryStatement
    | dropScheduledQueryStatement
    | dropViewStatement
    | dropMaterializedViewStatement
    | createFunctionStatement
    | createMacroStatement
    | dropFunctionStatement
    | reloadFunctionsStatement
    | dropMacroStatement
    | analyzeStatement
    | lockStatement
    | unlockStatement
    | lockDatabase
    | unlockDatabase
    | createRoleStatement
    | dropRoleStatement
    | grantPrivileges
    | revokePrivileges
    | showGrants
    | showRoleGrants
    | showRolePrincipals
    | showRoles
    | grantRole
    | revokeRole
    | setRole
    | showCurrentRole
    | abortTransactionStatement
    | abortCompactionStatement
    | killQueryStatement
    | resourcePlanDdlStatements
    | createDataConnectorStatement
    | dropDataConnectorStatement
    ;

// vendor/HiveParser.g:1113
ifExists
    : KW_IF KW_EXISTS
    ;

// vendor/HiveParser.g:1120
restrictOrCascade
    : KW_RESTRICT
    | KW_CASCADE
    ;

// vendor/HiveParser.g:1129
ifNotExists
    : KW_IF KW_NOT KW_EXISTS
    ;

// vendor/HiveParser.g:1136
force
    : KW_FORCE
    ;

// vendor/HiveParser.g:1143
rewriteEnabled
    : KW_ENABLE KW_REWRITE
    ;

// vendor/HiveParser.g:1150
rewriteDisabled
    : KW_DISABLE KW_REWRITE
    ;

// vendor/HiveParser.g:1157
storedAsDirs
    : KW_STORED KW_AS KW_DIRECTORIES
    ;

// vendor/HiveParser.g:1164
orReplace
    : KW_OR KW_REPLACE
    ;

// vendor/HiveParser.g:1171
createDatabaseStatement
    : KW_CREATE ( KW_DATABASE | KW_SCHEMA ) ifNotExists? identifier databaseComment? dbLocation? dbManagedLocation? ( KW_WITH KW_DBPROPERTIES dbProperties )?
    | KW_CREATE KW_REMOTE ( KW_DATABASE | KW_SCHEMA ) ifNotExists? identifier databaseComment? dbConnectorName ( KW_WITH KW_DBPROPERTIES dbProperties )?
    ;

// vendor/HiveParser.g:1192
dbLocation
    : KW_LOCATION StringLiteral
    ;

// vendor/HiveParser.g:1199
dbManagedLocation
    : KW_MANAGEDLOCATION StringLiteral
    ;

// vendor/HiveParser.g:1206
dbProperties
    : LPAREN dbPropertiesList RPAREN
    ;

// vendor/HiveParser.g:1213
dbPropertiesList
    : keyValueProperty ( COMMA keyValueProperty )*
    ;

// vendor/HiveParser.g:1220
dbConnectorName
    : KW_USING identifier
    ;

// vendor/HiveParser.g:1227
switchDatabaseStatement
    : KW_USE identifier
    ;

// vendor/HiveParser.g:1234
dropDatabaseStatement
    : KW_DROP ( KW_DATABASE | KW_SCHEMA ) ifExists? identifier restrictOrCascade?
    ;

// vendor/HiveParser.g:1241
databaseComment
    : KW_COMMENT StringLiteral
    ;

// vendor/HiveParser.g:1248
truncateTableStatement
    : KW_TRUNCATE KW_TABLE? tablePartitionPrefix ( KW_COLUMNS LPAREN columnNameList RPAREN )? force?
    ;

// vendor/HiveParser.g:1254
dropTableStatement
    : KW_DROP KW_TABLE ifExists? tableName KW_PURGE? replicationClause?
    ;

// vendor/HiveParser.g:1261
inputFileFormat
    : KW_INPUTFORMAT StringLiteral KW_SERDE StringLiteral
    ;

// vendor/HiveParser.g:1268
tabTypeExpr
    : identifier ( DOT identifier )? ( identifier ( DOT ( KW_ELEM_TYPE | KW_KEY_TYPE | KW_VALUE_TYPE | identifier ) )* )?
    ;

// vendor/HiveParser.g:1284
partTypeExpr
    : tabTypeExpr partitionSpec?
    ;

// vendor/HiveParser.g:1290
tabPartColTypeExpr
    : tableName partitionSpec? extColumnName?
    ;

// vendor/HiveParser.g:1296
descStatement
    : ( KW_DESCRIBE | KW_DESC ) ( ( KW_DATABASE | KW_SCHEMA ) KW_EXTENDED? ( identifier ) | ( KW_DATACONNECTOR ) KW_EXTENDED? ( identifier ) | KW_FUNCTION KW_EXTENDED? ( descFuncNames ) | ( ( KW_FORMATTED | KW_EXTENDED ) tabPartColTypeExpr ) | tabPartColTypeExpr )
    ;

// vendor/HiveParser.g:1314
analyzeStatement
    : KW_ANALYZE KW_TABLE ( tableOrPartition ) ( KW_COMPUTE KW_STATISTICS ( ( KW_NOSCAN ) | ( KW_FOR KW_COLUMNS ( columnNameList )? ) )? | KW_CACHE KW_METADATA )
    ;

// vendor/HiveParser.g:1327
showStatement
    : KW_SHOW ( KW_DATABASES | KW_SCHEMAS ) ( KW_LIKE showStmtIdentifier )?
    | KW_SHOW ( KW_EXTENDED )? KW_TABLES ( ( KW_FROM | KW_IN ) identifier )? ( showTablesFilterExpr )?
    | KW_SHOW KW_VIEWS ( ( KW_FROM | KW_IN ) identifier )? ( KW_LIKE showStmtIdentifier | showStmtIdentifier )?
    | KW_SHOW KW_MATERIALIZED KW_VIEWS ( ( KW_FROM | KW_IN ) identifier )? ( KW_LIKE showStmtIdentifier | showStmtIdentifier )?
    | KW_SHOW KW_SORTED? KW_COLUMNS ( KW_FROM | KW_IN ) tableName ( ( KW_FROM | KW_IN ) identifier )? ( KW_LIKE showStmtIdentifier | showStmtIdentifier )?
    | KW_SHOW KW_FUNCTIONS ( KW_LIKE showFunctionIdentifier )?
    | KW_SHOW KW_PARTITIONS tableName partitionSpec? whereClause? orderByClause? limitClause?
    | KW_SHOW KW_CREATE ( ( KW_DATABASE | KW_SCHEMA ) identifier | KW_TABLE tableName )
    | KW_SHOW KW_TABLE KW_EXTENDED ( ( KW_FROM | KW_IN ) identifier )? KW_LIKE showStmtIdentifier partitionSpec?
    | KW_SHOW KW_TBLPROPERTIES tableName ( LPAREN StringLiteral RPAREN )?
    | KW_SHOW KW_LOCKS ( ( KW_DATABASE | KW_SCHEMA ) ( identifier ) ( KW_EXTENDED )? | ( partTypeExpr )? ( KW_EXTENDED )? )
    | KW_SHOW KW_COMPACTIONS ( compactionId | ( KW_DATABASE | KW_SCHEMA ) ( identifier ) compactionPool? compactionType? compactionStatus? orderByClause? limitClause? | ( partTypeExpr )? compactionPool? compactionType? compactionStatus? orderByClause? limitClause? )
    | KW_SHOW KW_TRANSACTIONS
    | KW_SHOW KW_CONF StringLiteral
    | KW_SHOW KW_RESOURCE ( ( KW_PLAN identifier ) | ( KW_PLANS ) )
    | KW_SHOW ( KW_DATACONNECTORS )
    ;

// vendor/HiveParser.g:1371
showTablesFilterExpr
    : KW_WHERE identifier EQUAL StringLiteral
    | KW_LIKE showStmtIdentifier
    | showStmtIdentifier
    ;

// vendor/HiveParser.g:1380
lockStatement
    : KW_LOCK KW_TABLE tableName partitionSpec? lockMode
    ;

// vendor/HiveParser.g:1386
lockDatabase
    : KW_LOCK ( KW_DATABASE | KW_SCHEMA ) ( identifier ) lockMode
    ;

// vendor/HiveParser.g:1392
lockMode
    : KW_SHARED
    | KW_EXCLUSIVE
    ;

// vendor/HiveParser.g:1398
unlockStatement
    : KW_UNLOCK KW_TABLE tableName partitionSpec?
    ;

// vendor/HiveParser.g:1404
unlockDatabase
    : KW_UNLOCK ( KW_DATABASE | KW_SCHEMA ) ( identifier )
    ;

// vendor/HiveParser.g:1410
createRoleStatement
    : KW_CREATE KW_ROLE identifier
    ;

// vendor/HiveParser.g:1417
dropRoleStatement
    : KW_DROP KW_ROLE identifier
    ;

// vendor/HiveParser.g:1424
grantPrivileges
    : KW_GRANT privilegeList privilegeObject? KW_TO principalSpecification withGrantOption?
    ;

// vendor/HiveParser.g:1434
revokePrivileges
    : KW_REVOKE grantOptionFor? privilegeList privilegeObject? KW_FROM principalSpecification
    ;

// vendor/HiveParser.g:1441
grantRole
    : KW_GRANT KW_ROLE? identifier ( COMMA identifier )* KW_TO principalSpecification withAdminOption?
    ;

// vendor/HiveParser.g:1448
revokeRole
    : KW_REVOKE adminOptionFor? KW_ROLE? identifier ( COMMA identifier )* KW_FROM principalSpecification
    ;

// vendor/HiveParser.g:1455
showRoleGrants
    : KW_SHOW KW_ROLE KW_GRANT principalName
    ;

// vendor/HiveParser.g:1463
showRoles
    : KW_SHOW KW_ROLES
    ;

// vendor/HiveParser.g:1470
showCurrentRole
    : KW_SHOW KW_CURRENT KW_ROLES
    ;

// vendor/HiveParser.g:1477
setRole
    : KW_SET KW_ROLE ( ( KW_ALL ) | ( KW_NONE ) | identifier )
    ;

// vendor/HiveParser.g:1490
showGrants
    : KW_SHOW KW_GRANT principalName? ( KW_ON privilegeIncludeColObject )?
    ;

// vendor/HiveParser.g:1497
showRolePrincipals
    : KW_SHOW KW_PRINCIPALS identifier
    ;

// vendor/HiveParser.g:1505
privilegeIncludeColObject
    : KW_ALL
    | privObjectCols
    ;

// vendor/HiveParser.g:1512
privilegeObject
    : KW_ON privObject
    ;

// vendor/HiveParser.g:1519
privObject
    : ( KW_DATABASE | KW_SCHEMA ) identifier
    | KW_TABLE? tableName partitionSpec?
    | KW_URI ( StringLiteral )
    | KW_SERVER identifier
    ;

// vendor/HiveParser.g:1526
privObjectCols
    : ( KW_DATABASE | KW_SCHEMA ) identifier
    | KW_TABLE? tableName ( LPAREN columnNameList RPAREN )? partitionSpec?
    | KW_URI ( StringLiteral )
    | KW_SERVER identifier
    ;

// vendor/HiveParser.g:1533
privilegeList
    : privlegeDef ( COMMA privlegeDef )*
    ;

// vendor/HiveParser.g:1540
privlegeDef
    : privilegeType ( LPAREN columnNameList RPAREN )?
    ;

// vendor/HiveParser.g:1547
privilegeType
    : KW_ALL
    | KW_ALTER
    | KW_UPDATE
    | KW_CREATE
    | KW_DROP
    | KW_LOCK
    | KW_SELECT
    | KW_SHOW_DATABASE
    | KW_INSERT
    | KW_DELETE
    ;

// vendor/HiveParser.g:1562
principalSpecification
    : principalName ( COMMA principalName )*
    ;

// vendor/HiveParser.g:1568
principalName
    : KW_USER principalIdentifier
    | KW_GROUP principalIdentifier
    | KW_ROLE identifier
    ;

// vendor/HiveParser.g:1576
withGrantOption
    : KW_WITH KW_GRANT KW_OPTION
    ;

// vendor/HiveParser.g:1583
grantOptionFor
    : KW_GRANT KW_OPTION KW_FOR
    ;

// vendor/HiveParser.g:1590
adminOptionFor
    : KW_ADMIN KW_OPTION KW_FOR
    ;

// vendor/HiveParser.g:1597
withAdminOption
    : KW_WITH KW_ADMIN KW_OPTION
    ;

// vendor/HiveParser.g:1604
metastoreCheck
    : KW_MSCK ( KW_REPAIR )? ( KW_TABLE tableName ( ( KW_ADD | KW_DROP | KW_SYNC ) ( KW_PARTITIONS ) ( partitionSelectorSpec )? )? )
    ;

// vendor/HiveParser.g:1613
resourceList
    : resource ( COMMA resource )*
    ;

// vendor/HiveParser.g:1620
resource
    : resourceType StringLiteral
    ;

// vendor/HiveParser.g:1627
resourceType
    : KW_JAR
    | KW_FILE
    | KW_ARCHIVE
    ;

// vendor/HiveParser.g:1638
createFunctionStatement
    : KW_CREATE ( KW_TEMPORARY )? KW_FUNCTION functionIdentifier KW_AS StringLiteral ( KW_USING resourceList )?
    ;

// vendor/HiveParser.g:1647
dropFunctionStatement
    : KW_DROP ( KW_TEMPORARY )? KW_FUNCTION ifExists? functionIdentifier
    ;

// vendor/HiveParser.g:1655
reloadFunctionsStatement
    : KW_RELOAD ( KW_FUNCTIONS | KW_FUNCTION )
    ;

// vendor/HiveParser.g:1660
createMacroStatement
    : KW_CREATE KW_TEMPORARY KW_MACRO Identifier LPAREN columnNameTypeList? RPAREN expression
    ;

// vendor/HiveParser.g:1668
dropMacroStatement
    : KW_DROP KW_TEMPORARY KW_MACRO ifExists? Identifier
    ;

// vendor/HiveParser.g:1675
createViewStatement
    : KW_CREATE ( orReplace )? KW_VIEW ( ifNotExists )? tableName ( LPAREN columnNameCommentList RPAREN )? tableComment? viewPartition? tablePropertiesPrefixed? KW_AS selectStatementWithCTE
    ;

// vendor/HiveParser.g:1695
viewPartition
    : KW_PARTITIONED KW_ON LPAREN columnNameList RPAREN
    | KW_PARTITIONED KW_ON KW_SPEC LPAREN ( partitionTransformSpec ) RPAREN
    ;

// vendor/HiveParser.g:1704
viewOrganization
    : viewClusterSpec
    | viewComplexSpec
    ;

// vendor/HiveParser.g:1711
viewClusterSpec
    : KW_CLUSTERED KW_ON LPAREN columnNameList RPAREN
    ;

// vendor/HiveParser.g:1718
viewComplexSpec
    : viewDistSpec viewSortSpec
    ;

// vendor/HiveParser.g:1724
viewDistSpec
    : KW_DISTRIBUTED KW_ON LPAREN columnNameList RPAREN
    ;

// vendor/HiveParser.g:1731
viewSortSpec
    : KW_SORTED KW_ON LPAREN columnNameList RPAREN
    ;

// vendor/HiveParser.g:1738
dropViewStatement
    : KW_DROP KW_VIEW ifExists? viewName
    ;

// vendor/HiveParser.g:1744
createMaterializedViewStatement
    : KW_CREATE KW_MATERIALIZED KW_VIEW ( ifNotExists )? tableName rewriteDisabled? tableComment? viewPartition? viewOrganization? tableRowFormat? tableFileFormat? tableLocation? tablePropertiesPrefixed? KW_AS selectStatementWithCTE
    ;

// vendor/HiveParser.g:1767
dropMaterializedViewStatement
    : KW_DROP KW_MATERIALIZED KW_VIEW ifExists? viewName
    ;

// vendor/HiveParser.g:1773
createScheduledQueryStatement
    : KW_CREATE KW_SCHEDULED KW_QUERY identifier scheduleSpec executedAsSpec? enableSpecification? definedAsSpec
    ;

// vendor/HiveParser.g:1790
dropScheduledQueryStatement
    : KW_DROP KW_SCHEDULED KW_QUERY identifier
    ;

// vendor/HiveParser.g:1800
alterScheduledQueryStatement
    : KW_ALTER KW_SCHEDULED KW_QUERY identifier alterScheduledQueryChange
    ;

// vendor/HiveParser.g:1811
alterScheduledQueryChange
    : scheduleSpec
    | executedAsSpec
    | enableSpecification
    | definedAsSpec
    | KW_EXECUTE
    ;

// vendor/HiveParser.g:1821
scheduleSpec
    : KW_CRON StringLiteral
    | KW_EVERY Number? intervalQualifiers ( ( KW_AT | KW_OFFSET KW_BY ) StringLiteral )?
    ;

// vendor/HiveParser.g:1829
executedAsSpec
    : KW_EXECUTED KW_AS StringLiteral
    ;

// vendor/HiveParser.g:1835
definedAsSpec
    : KW_DEFINED? KW_AS statement
    ;

// vendor/HiveParser.g:1841
showFunctionIdentifier
    : functionIdentifier
    | StringLiteral
    ;

// vendor/HiveParser.g:1848
showStmtIdentifier
    : identifier
    | StringLiteral
    ;

// vendor/HiveParser.g:1855
tableComment
    : KW_COMMENT StringLiteral
    ;

// vendor/HiveParser.g:1862
createTablePartitionSpec
    : KW_PARTITIONED KW_BY LPAREN ( createTablePartitionColumnTypeSpec | createTablePartitionColumnSpec ) RPAREN
    | KW_PARTITIONED KW_BY KW_SPEC LPAREN ( partitionTransformSpec ) RPAREN
    ;

// vendor/HiveParser.g:1872
createTablePartitionColumnTypeSpec
    : columnNameTypeConstraint ( COMMA columnNameTypeConstraint )*
    ;

// vendor/HiveParser.g:1879
createTablePartitionColumnSpec
    : columnName ( COMMA columnName )*
    ;

// vendor/HiveParser.g:1886
partitionTransformSpec
    : columnNameTransformConstraint ( COMMA columnNameTransformConstraint )*
    ;

// vendor/HiveParser.g:1893
columnNameTransformConstraint
    : partitionTransformType
    ;

// vendor/HiveParser.g:1900
partitionTransformType
    : columnName
    | KW_YEAR LPAREN columnName RPAREN
    | KW_MONTH LPAREN columnName RPAREN
    | KW_DAY LPAREN columnName RPAREN
    | KW_HOUR LPAREN columnName RPAREN
    | KW_TRUNCATE LPAREN Number COMMA columnName RPAREN
    | KW_BUCKET LPAREN Number COMMA columnName RPAREN
    ;

// vendor/HiveParser.g:1926
tableBuckets
    : KW_CLUSTERED KW_BY LPAREN columnNameList RPAREN ( KW_SORTED KW_BY LPAREN columnNameOrderList RPAREN )? KW_INTO Number KW_BUCKETS
    ;

// vendor/HiveParser.g:1934
tableImplBuckets
    : KW_CLUSTERED KW_INTO Number KW_BUCKETS
    ;

// vendor/HiveParser.g:1942
tableSkewed
    : KW_SKEWED KW_BY LPAREN columnNameList RPAREN KW_ON LPAREN ( skewedValueElement ) RPAREN ( storedAsDirs )?
    ;

// vendor/HiveParser.g:1950
rowFormat
    : rowFormatSerde
    | rowFormatDelimited
    | /* empty */
    ;

// vendor/HiveParser.g:1958
recordReader
    : KW_RECORDREADER StringLiteral
    | /* empty */
    ;

// vendor/HiveParser.g:1965
recordWriter
    : KW_RECORDWRITER StringLiteral
    | /* empty */
    ;

// vendor/HiveParser.g:1972
rowFormatSerde
    : KW_ROW KW_FORMAT KW_SERDE StringLiteral ( KW_WITH KW_SERDEPROPERTIES tableProperties )?
    ;

// vendor/HiveParser.g:1979
rowFormatDelimited
    : KW_ROW KW_FORMAT KW_DELIMITED tableRowFormatFieldIdentifier? tableRowFormatCollItemsIdentifier? tableRowFormatMapKeysIdentifier? tableRowFormatLinesIdentifier? tableRowNullFormat?
    ;

// vendor/HiveParser.g:1987
tableRowFormat
    : rowFormatDelimited
    | rowFormatSerde
    ;

// vendor/HiveParser.g:1997
tablePropertiesPrefixed
    : KW_TBLPROPERTIES tableProperties
    ;

// vendor/HiveParser.g:2004
tableProperties
    : LPAREN tablePropertiesList RPAREN
    ;

// vendor/HiveParser.g:2011
tablePropertiesList
    : keyValueProperty ( COMMA keyValueProperty )*
    | keyProperty ( COMMA keyProperty )*
    ;

// vendor/HiveParser.g:2020
keyValueProperty
    : StringLiteral EQUAL StringLiteral
    ;

// vendor/HiveParser.g:2027
keyProperty
    : StringLiteral
    ;

// vendor/HiveParser.g:2034
tableRowFormatFieldIdentifier
    : KW_FIELDS KW_TERMINATED KW_BY StringLiteral ( KW_ESCAPED KW_BY StringLiteral )?
    ;

// vendor/HiveParser.g:2042
tableRowFormatCollItemsIdentifier
    : KW_COLLECTION KW_ITEMS KW_TERMINATED KW_BY StringLiteral
    ;

// vendor/HiveParser.g:2050
tableRowFormatMapKeysIdentifier
    : KW_MAP KW_KEYS KW_TERMINATED KW_BY StringLiteral
    ;

// vendor/HiveParser.g:2058
tableRowFormatLinesIdentifier
    : KW_LINES KW_TERMINATED KW_BY StringLiteral
    ;

// vendor/HiveParser.g:2066
tableRowNullFormat
    : KW_NULL KW_DEFINED KW_AS StringLiteral
    ;

// vendor/HiveParser.g:2073
tableFileFormat
    : KW_STORED KW_AS KW_INPUTFORMAT StringLiteral KW_OUTPUTFORMAT StringLiteral ( KW_INPUTDRIVER StringLiteral KW_OUTPUTDRIVER StringLiteral )?
    | KW_STORED KW_BY StringLiteral ( KW_WITH KW_SERDEPROPERTIES tableProperties )? ( KW_STORED KW_AS identifier )?
    | KW_STORED KW_BY identifier ( KW_WITH KW_SERDEPROPERTIES tableProperties )? ( KW_STORED KW_AS identifier )?
    | KW_STORED KW_AS identifier
    ;

// vendor/HiveParser.g:2091
tableLocation
    : KW_LOCATION StringLiteral
    ;

// vendor/HiveParser.g:2098
columnNameTypeList
    : columnNameType ( COMMA columnNameType )*
    ;

// vendor/HiveParser.g:2104
columnNameTypeOrConstraintList
    : columnNameTypeOrConstraint ( COMMA columnNameTypeOrConstraint )*
    ;

// vendor/HiveParser.g:2110
columnNameColonTypeList
    : columnNameColonType ( COMMA columnNameColonType )*
    ;

// vendor/HiveParser.g:2116
columnNameList
    : columnName ( COMMA columnName )*
    ;

// vendor/HiveParser.g:2122
columnName
    : identifier
    ;

// vendor/HiveParser.g:2129
extColumnName
    : identifier ( DOT ( KW_ELEM_TYPE | KW_KEY_TYPE | KW_VALUE_TYPE | identifier ) )*
    ;

// vendor/HiveParser.g:2136
columnNameOrderList
    : columnNameOrder ( COMMA columnNameOrder )*
    ;

// vendor/HiveParser.g:2142
columnParenthesesList
    : LPAREN columnNameList RPAREN
    ;

// vendor/HiveParser.g:2148
enableValidateSpecification
    : enableSpecification validateSpecification?
    | enforcedSpecification
    ;

// vendor/HiveParser.g:2155
enableSpecification
    : KW_ENABLE
    | KW_DISABLE
    ;

// vendor/HiveParser.g:2162
validateSpecification
    : KW_VALIDATE
    | KW_NOVALIDATE
    ;

// vendor/HiveParser.g:2169
enforcedSpecification
    : KW_ENFORCED
    | KW_NOT KW_ENFORCED
    ;

// vendor/HiveParser.g:2176
relySpecification
    : KW_RELY
    | KW_NORELY
    ;

// vendor/HiveParser.g:2183
createConstraint
    : ( KW_CONSTRAINT identifier )? tableLevelConstraint constraintOptsCreate?
    ;

// vendor/HiveParser.g:2192
alterConstraintWithName
    : KW_CONSTRAINT identifier tableLevelConstraint constraintOptsAlter?
    ;

// vendor/HiveParser.g:2199
tableLevelConstraint
    : pkUkConstraint
    | checkConstraint
    ;

// vendor/HiveParser.g:2204
pkUkConstraint
    : tableConstraintType columnParenthesesList
    ;

// vendor/HiveParser.g:2211
checkConstraint
    : KW_CHECK LPAREN expression RPAREN
    ;

// vendor/HiveParser.g:2218
createForeignKey
    : ( KW_CONSTRAINT identifier )? KW_FOREIGN KW_KEY columnParenthesesList KW_REFERENCES tableName columnParenthesesList constraintOptsCreate?
    ;

// vendor/HiveParser.g:2227
alterForeignKeyWithName
    : KW_CONSTRAINT identifier KW_FOREIGN KW_KEY columnParenthesesList KW_REFERENCES tableName columnParenthesesList constraintOptsAlter?
    ;

// vendor/HiveParser.g:2234
skewedValueElement
    : skewedColumnValues
    | skewedColumnValuePairList
    ;

// vendor/HiveParser.g:2242
skewedColumnValuePairList
    : skewedColumnValuePair ( COMMA skewedColumnValuePair )*
    ;

// vendor/HiveParser.g:2248
skewedColumnValuePair
    : LPAREN skewedColumnValues RPAREN
    ;

// vendor/HiveParser.g:2256
skewedColumnValues
    : skewedColumnValue ( COMMA skewedColumnValue )*
    ;

// vendor/HiveParser.g:2262
skewedColumnValue
    : constant
    ;

// vendor/HiveParser.g:2269
skewedValueLocationElement
    : skewedColumnValue
    | skewedColumnValuePair
    ;

// vendor/HiveParser.g:2277
orderSpecification
    : KW_ASC
    | KW_DESC
    ;

// vendor/HiveParser.g:2282
nullOrdering
    : KW_NULLS KW_FIRST
    | KW_NULLS KW_LAST
    ;

// vendor/HiveParser.g:2289
columnNameOrder
    : identifier orderSpecification? nullOrdering?
    ;

// vendor/HiveParser.g:2306
columnNameCommentList
    : columnNameComment ( COMMA columnNameComment )*
    ;

// vendor/HiveParser.g:2312
columnNameComment
    : identifier ( KW_COMMENT StringLiteral )?
    ;

// vendor/HiveParser.g:2319
orderSpecificationRewrite
    : KW_ASC
    | KW_DESC
    ;

// vendor/HiveParser.g:2326
columnRefOrder
    : expression orderSpecificationRewrite? nullOrdering?
    ;

// vendor/HiveParser.g:2355
columnNameType
    : identifier colType ( KW_COMMENT StringLiteral )?
    ;

// vendor/HiveParser.g:2364
columnNameTypeOrConstraint
    : ( tableConstraint )
    | ( columnNameTypeConstraint )
    ;

// vendor/HiveParser.g:2371
tableConstraint
    : ( createForeignKey )
    | ( createConstraint )
    ;

// vendor/HiveParser.g:2378
columnNameTypeConstraint
    : identifier colType columnConstraint? ( KW_COMMENT StringLiteral )?
    ;

// vendor/HiveParser.g:2386
columnConstraint
    : ( foreignKeyConstraint )
    | ( colConstraint )
    ;

// vendor/HiveParser.g:2393
foreignKeyConstraint
    : ( KW_CONSTRAINT identifier )? KW_REFERENCES tableName LPAREN columnName RPAREN constraintOptsCreate?
    ;

// vendor/HiveParser.g:2402
colConstraint
    : ( KW_CONSTRAINT identifier )? columnConstraintType constraintOptsCreate?
    ;

// vendor/HiveParser.g:2411
alterColumnConstraint
    : ( alterForeignKeyConstraint )
    | ( alterColConstraint )
    ;

// vendor/HiveParser.g:2418
alterForeignKeyConstraint
    : ( KW_CONSTRAINT identifier )? KW_REFERENCES tableName LPAREN columnName RPAREN constraintOptsAlter?
    ;

// vendor/HiveParser.g:2427
alterColConstraint
    : ( KW_CONSTRAINT identifier )? columnConstraintType constraintOptsAlter?
    ;

// vendor/HiveParser.g:2436
columnConstraintType
    : KW_NOT KW_NULL
    | KW_DEFAULT defaultVal
    | checkConstraint
    | tableConstraintType
    ;

// vendor/HiveParser.g:2443
defaultVal
    : constant
    | function
    | castExpression
    ;

// vendor/HiveParser.g:2449
tableConstraintType
    : KW_PRIMARY KW_KEY
    | KW_UNIQUE
    ;

// vendor/HiveParser.g:2454
constraintOptsCreate
    : enableValidateSpecification relySpecification?
    ;

// vendor/HiveParser.g:2458
constraintOptsAlter
    : enableValidateSpecification relySpecification?
    ;

// vendor/HiveParser.g:2462
columnNameColonType
    : identifier COLON colType ( KW_COMMENT StringLiteral )?
    ;

// vendor/HiveParser.g:2470
colType
    : type
    ;

// vendor/HiveParser.g:2476
colTypeList
    : colType ( COMMA colType )*
    ;

// vendor/HiveParser.g:2482
type
    : primitiveType
    | listType
    | structType
    | mapType
    | unionType
    ;

// vendor/HiveParser.g:2489
primitiveType
    : KW_TINYINT
    | KW_SMALLINT
    | KW_INT
    | KW_BIGINT
    | KW_BOOLEAN
    | KW_FLOAT
    | KW_REAL
    | KW_DOUBLE KW_PRECISION?
    | KW_DATE
    | KW_DATETIME
    | KW_TIMESTAMP
    | KW_TIMESTAMPLOCALTZ
    | KW_TIMESTAMP KW_WITH KW_LOCAL KW_TIME KW_ZONE
    | KW_STRING
    | KW_BINARY
    | KW_DECIMAL ( LPAREN Number ( COMMA Number )? RPAREN )?
    | KW_VARCHAR LPAREN Number RPAREN
    | KW_CHAR LPAREN Number RPAREN
    ;

// vendor/HiveParser.g:2517
listType
    : KW_ARRAY LESSTHAN type GREATERTHAN
    ;

// vendor/HiveParser.g:2523
structType
    : KW_STRUCT LESSTHAN columnNameColonTypeList GREATERTHAN
    ;

// vendor/HiveParser.g:2529
mapType
    : KW_MAP LESSTHAN primitiveType COMMA type GREATERTHAN
    ;

// vendor/HiveParser.g:2536
unionType
    : KW_UNIONTYPE LESSTHAN colTypeList GREATERTHAN
    ;

// vendor/HiveParser.g:2542
setOperator
    : KW_UNION KW_ALL
    | KW_UNION KW_DISTINCT?
    | KW_INTERSECT KW_ALL
    | KW_INTERSECT KW_DISTINCT?
    | KW_EXCEPT KW_ALL
    | KW_EXCEPT KW_DISTINCT?
    | KW_MINUS KW_ALL
    | KW_MINUS KW_DISTINCT?
    ;

// vendor/HiveParser.g:2555
queryStatementExpression
    : ( withClause )? queryStatementExpressionBody
    ;

// vendor/HiveParser.g:2570
queryStatementExpressionBody
    : fromStatement
    | regularBody
    ;

// vendor/HiveParser.g:2576
withClause
    : KW_WITH cteStatement ( COMMA cteStatement )*
    ;

// vendor/HiveParser.g:2581
cteStatement
    : identifier ( LPAREN columnNameList RPAREN )? KW_AS LPAREN queryStatementExpression RPAREN
    ;

// vendor/HiveParser.g:2587
fromStatement
    : ( singleFromStatement ) ( setOperator singleFromStatement )*
    ;

// vendor/HiveParser.g:2608
singleFromStatement
    : fromClause ( body )+
    ;

// vendor/HiveParser.g:2620
regularBody
    : insertClause ( selectStatement )
    | selectStatement
    ;

// vendor/HiveParser.g:2631
atomSelectStatement
    : selectClause fromClause? whereClause? groupByClause? havingClause? window_clause? qualifyClause?
    | LPAREN selectStatement RPAREN
    | valuesSource
    ;

// vendor/HiveParser.g:2648
selectStatement
    : atomSelectStatement setOpSelectStatement? orderByClause? clusterByClause? distributeByClause? sortByClause? limitClause?
    ;

// vendor/HiveParser.g:2685
setOpSelectStatement
    : ( setOperator atomSelectStatement )+
    ;

// vendor/HiveParser.g:2738
selectStatementWithCTE
    : ( withClause )? selectStatement
    ;

// vendor/HiveParser.g:2749
body
    : insertClause selectClause lateralView? whereClause? groupByClause? havingClause? window_clause? qualifyClause? orderByClause? clusterByClause? distributeByClause? sortByClause? limitClause?
    | selectClause lateralView? whereClause? groupByClause? havingClause? window_clause? qualifyClause? orderByClause? clusterByClause? distributeByClause? sortByClause? limitClause?
    ;

// vendor/HiveParser.g:2783
insertClause
    : KW_INSERT KW_OVERWRITE destination ifNotExists?
    | KW_INSERT KW_INTO KW_TABLE? tableOrPartition ( LPAREN columnNameList RPAREN )?
    ;

// vendor/HiveParser.g:2792
destination
    : ( KW_LOCAL )? KW_DIRECTORY StringLiteral tableRowFormat? tableFileFormat?
    | KW_TABLE tableOrPartition
    ;

// vendor/HiveParser.g:2801
limitClause
    : KW_LIMIT ( ( Number COMMA )? Number )
    | KW_LIMIT Number KW_OFFSET Number
    ;

// vendor/HiveParser.g:2810
deleteStatement
    : KW_DELETE KW_FROM tableName ( whereClause )?
    ;

// vendor/HiveParser.g:2818
columnAssignmentClause
    : /* empty */
    | tableOrColumn EQUAL precedencePlusExpressionOrDefault
    ;

// vendor/HiveParser.g:2823
precedencePlusExpressionOrDefault
    : defaultValue
    | precedencePlusExpression
    ;

// vendor/HiveParser.g:2831
setColumnsClause
    : KW_SET columnAssignmentClause ( COMMA columnAssignmentClause )*
    ;

// vendor/HiveParser.g:2840
updateStatement
    : KW_UPDATE tableName setColumnsClause whereClause?
    ;

// vendor/HiveParser.g:2851
sqlTransactionStatement
    : startTransactionStatement
    | commitStatement
    | rollbackStatement
    | setAutoCommitStatement
    ;

// vendor/HiveParser.g:2861
startTransactionStatement
    : KW_START KW_TRANSACTION ( transactionMode ( COMMA transactionMode )* )?
    ;

// vendor/HiveParser.g:2866
transactionMode
    : isolationLevel
    | transactionAccessMode
    ;

// vendor/HiveParser.g:2872
transactionAccessMode
    : KW_READ KW_ONLY
    | KW_READ KW_WRITE
    ;

// vendor/HiveParser.g:2878
isolationLevel
    : KW_ISOLATION KW_LEVEL levelOfIsolation
    ;

// vendor/HiveParser.g:2884
levelOfIsolation
    : KW_SNAPSHOT
    ;

// vendor/HiveParser.g:2889
commitStatement
    : KW_COMMIT ( KW_WORK )?
    ;

// vendor/HiveParser.g:2894
rollbackStatement
    : KW_ROLLBACK ( KW_WORK )?
    ;

// vendor/HiveParser.g:2898
setAutoCommitStatement
    : KW_SET KW_AUTOCOMMIT booleanValueTok
    ;

// vendor/HiveParser.g:2906
abortTransactionStatement
    : KW_ABORT KW_TRANSACTIONS ( Number )+
    ;

// vendor/HiveParser.g:2913
abortCompactionStatement
    : KW_ABORT KW_COMPACTIONS ( Number )+
    ;

// vendor/HiveParser.g:2923
mergeStatement
    : KW_MERGE QUERY_HINT? KW_INTO tableName ( KW_AS? identifier )? KW_USING joinSourcePart KW_ON expression whenClauses
    ;

// vendor/HiveParser.g:2936
whenClauses
    : ( whenMatchedAndClause | whenMatchedThenClause )* whenNotMatchedClause?
    ;

// vendor/HiveParser.g:2940
whenNotMatchedClause
    : KW_WHEN KW_NOT KW_MATCHED ( KW_AND expression )? KW_THEN KW_INSERT ( columnParenthesesList )? KW_VALUES valueRowConstructor
    ;

// vendor/HiveParser.g:2947
whenMatchedAndClause
    : KW_WHEN KW_MATCHED KW_AND expression KW_THEN updateOrDelete
    ;

// vendor/HiveParser.g:2954
whenMatchedThenClause
    : KW_WHEN KW_MATCHED KW_THEN updateOrDelete
    ;

// vendor/HiveParser.g:2961
updateOrDelete
    : KW_UPDATE setColumnsClause
    | KW_DELETE
    ;

// vendor/HiveParser.g:2971
killQueryStatement
    : KW_KILL KW_QUERY ( StringLiteral )+
    ;

// vendor/HiveParser.g:2981
compactionId
    : KW_COMPACT_ID EQUAL Number
    ;

// vendor/HiveParser.g:2984
compactionPool
    : KW_POOL StringLiteral
    ;

// vendor/HiveParser.g:2987
compactionType
    : KW_TYPE StringLiteral
    ;

// vendor/HiveParser.g:2990
compactionStatus
    : KW_STATUS StringLiteral
    ;

// vendor/AlterClauseParser.g:47
alterStatement
    : KW_ALTER KW_TABLE tableName alterTableStatementSuffix
    | KW_ALTER KW_VIEW tableName KW_AS? alterViewStatementSuffix
    | KW_ALTER KW_MATERIALIZED KW_VIEW tableName alterMaterializedViewStatementSuffix
    | KW_ALTER ( KW_DATABASE | KW_SCHEMA ) alterDatabaseStatementSuffix
    | KW_ALTER KW_DATACONNECTOR alterDataConnectorStatementSuffix
    | KW_OPTIMIZE KW_TABLE tableName optimizeTableStatementSuffix
    ;

// vendor/AlterClauseParser.g:58
alterTableStatementSuffix
    : alterStatementSuffixRename
    | alterStatementSuffixDropPartitions
    | alterStatementSuffixAddPartitions
    | alterStatementSuffixTouch
    | alterStatementSuffixArchive
    | alterStatementSuffixUnArchive
    | alterStatementSuffixProperties
    | alterStatementSuffixSkewedby
    | alterStatementSuffixExchangePartition
    | alterStatementPartitionKeyType
    | alterStatementSuffixDropConstraint
    | alterStatementSuffixAddConstraint
    | alterTblPartitionStatementSuffix
    | partitionSpec alterTblPartitionStatementSuffix
    | alterStatementSuffixSetOwner
    | alterStatementSuffixSetPartSpec
    | alterStatementSuffixExecute
    | alterStatementSuffixCreateBranch
    | alterStatementSuffixDropBranch
    | alterStatementSuffixCreateTag
    | alterStatementSuffixDropTag
    | alterStatementSuffixConvert
    ;

// vendor/AlterClauseParser.g:85
alterTblPartitionStatementSuffix
    : alterStatementSuffixFileFormat
    | alterStatementSuffixLocation
    | alterStatementSuffixMergeFiles
    | alterStatementSuffixSerdeProperties
    | alterStatementSuffixRenamePart
    | alterStatementSuffixBucketNum
    | alterTblPartitionStatementSuffixSkewedLocation
    | alterStatementSuffixClusterbySortby
    | alterStatementSuffixCompact
    | alterStatementSuffixUpdateStatsCol
    | alterStatementSuffixUpdateStats
    | alterStatementSuffixRenameCol
    | alterStatementSuffixAddCol
    | alterStatementSuffixUpdateColumns
    ;

// vendor/AlterClauseParser.g:104
optimizeTableStatementSuffix
    : optimizeTblRewriteDataSuffix
    ;

// vendor/AlterClauseParser.g:110
optimizeTblRewriteDataSuffix
    : KW_REWRITE KW_DATA
    ;

// vendor/AlterClauseParser.g:117
alterStatementPartitionKeyType
    : KW_PARTITION KW_COLUMN LPAREN columnNameType RPAREN
    ;

// vendor/AlterClauseParser.g:124
alterViewStatementSuffix
    : alterViewSuffixProperties
    | alterStatementSuffixRename
    | alterStatementSuffixAddPartitions
    | alterStatementSuffixDropPartitions
    | selectStatementWithCTE
    ;

// vendor/AlterClauseParser.g:134
alterMaterializedViewStatementSuffix
    : alterMaterializedViewSuffixRewrite
    | alterMaterializedViewSuffixRebuild
    ;

// vendor/AlterClauseParser.g:141
alterMaterializedViewSuffixRewrite
    : ( rewriteEnabled | rewriteDisabled )
    ;

// vendor/AlterClauseParser.g:148
alterMaterializedViewSuffixRebuild
    : KW_REBUILD
    ;

// vendor/AlterClauseParser.g:154
alterDatabaseStatementSuffix
    : alterDatabaseSuffixProperties
    | alterDatabaseSuffixSetOwner
    | alterDatabaseSuffixSetLocation
    ;

// vendor/AlterClauseParser.g:162
alterDatabaseSuffixProperties
    : identifier KW_SET KW_DBPROPERTIES dbProperties
    ;

// vendor/AlterClauseParser.g:169
alterDatabaseSuffixSetOwner
    : identifier KW_SET KW_OWNER principalName
    ;

// vendor/AlterClauseParser.g:176
alterDatabaseSuffixSetLocation
    : identifier KW_SET KW_LOCATION StringLiteral
    | identifier KW_SET KW_MANAGEDLOCATION StringLiteral
    ;

// vendor/AlterClauseParser.g:185
alterDatabaseSuffixSetManagedLocation
    : identifier KW_SET KW_MANAGEDLOCATION StringLiteral
    ;

// vendor/AlterClauseParser.g:192
alterStatementSuffixRename
    : KW_RENAME KW_TO tableName
    ;

// vendor/AlterClauseParser.g:200
alterStatementSuffixAddCol
    : ( KW_ADD | KW_REPLACE ) KW_COLUMNS LPAREN columnNameTypeList RPAREN restrictOrCascade?
    ;

// vendor/AlterClauseParser.g:208
alterStatementSuffixAddConstraint
    : KW_ADD ( alterForeignKeyWithName | alterConstraintWithName )
    ;

// vendor/AlterClauseParser.g:216
alterStatementSuffixUpdateColumns
    : KW_UPDATE KW_COLUMNS restrictOrCascade?
    ;

// vendor/AlterClauseParser.g:223
alterStatementSuffixDropConstraint
    : KW_DROP KW_CONSTRAINT identifier
    ;

// vendor/AlterClauseParser.g:230
alterStatementSuffixRenameCol
    : KW_CHANGE KW_COLUMN? identifier identifier colType alterColumnConstraint? ( KW_COMMENT StringLiteral )? alterStatementChangeColPosition? restrictOrCascade?
    ;

// vendor/AlterClauseParser.g:237
alterStatementSuffixUpdateStatsCol
    : KW_UPDATE KW_STATISTICS KW_FOR KW_COLUMN? identifier KW_SET tableProperties ( KW_COMMENT StringLiteral )?
    ;

// vendor/AlterClauseParser.g:245
alterStatementSuffixUpdateStats
    : KW_UPDATE KW_STATISTICS KW_SET tableProperties
    ;

// vendor/AlterClauseParser.g:253
alterStatementChangeColPosition
    : KW_FIRST
    | KW_AFTER identifier
    ;

// vendor/AlterClauseParser.g:259
alterStatementSuffixAddPartitions
    : KW_ADD ifNotExists? alterStatementSuffixAddPartitionsElement+
    ;

// vendor/AlterClauseParser.g:267
alterStatementSuffixAddPartitionsElement
    : partitionSpec partitionLocation?
    ;

// vendor/AlterClauseParser.g:271
alterStatementSuffixTouch
    : KW_TOUCH ( partitionSpec )*
    ;

// vendor/AlterClauseParser.g:278
alterStatementSuffixArchive
    : KW_ARCHIVE ( partitionSpec )*
    ;

// vendor/AlterClauseParser.g:285
alterStatementSuffixUnArchive
    : KW_UNARCHIVE ( partitionSpec )*
    ;

// vendor/AlterClauseParser.g:292
partitionLocation
    : KW_LOCATION StringLiteral
    ;

// vendor/AlterClauseParser.g:299
alterStatementSuffixDropPartitions
    : KW_DROP ifExists? KW_PARTITION partitionSelectorSpec ( COMMA KW_PARTITION partitionSelectorSpec )* KW_PURGE? replicationClause?
    ;

// vendor/AlterClauseParser.g:307
alterStatementSuffixProperties
    : KW_SET KW_TBLPROPERTIES tableProperties
    | KW_UNSET KW_TBLPROPERTIES ifExists? tableProperties
    ;

// vendor/AlterClauseParser.g:316
alterViewSuffixProperties
    : KW_SET KW_TBLPROPERTIES tableProperties
    | KW_UNSET KW_TBLPROPERTIES ifExists? tableProperties
    ;

// vendor/AlterClauseParser.g:325
alterStatementSuffixSerdeProperties
    : KW_SET KW_SERDE StringLiteral ( KW_WITH KW_SERDEPROPERTIES tableProperties )?
    | KW_SET KW_SERDEPROPERTIES tableProperties
    | KW_UNSET KW_SERDEPROPERTIES tableProperties
    ;

// vendor/AlterClauseParser.g:339
tablePartitionPrefix
    : tableName partitionSpec?
    ;

// vendor/AlterClauseParser.g:346
alterStatementSuffixFileFormat
    : KW_SET KW_FILEFORMAT fileFormat
    ;

// vendor/AlterClauseParser.g:354
alterStatementSuffixClusterbySortby
    : KW_NOT KW_CLUSTERED
    | KW_NOT KW_SORTED
    | tableBuckets
    ;

// vendor/AlterClauseParser.g:362
alterTblPartitionStatementSuffixSkewedLocation
    : KW_SET KW_SKEWED KW_LOCATION skewedLocations
    ;

// vendor/AlterClauseParser.g:369
skewedLocations
    : LPAREN skewedLocationsList RPAREN
    ;

// vendor/AlterClauseParser.g:376
skewedLocationsList
    : skewedLocationMap ( COMMA skewedLocationMap )*
    ;

// vendor/AlterClauseParser.g:383
skewedLocationMap
    : skewedValueLocationElement EQUAL StringLiteral
    ;

// vendor/AlterClauseParser.g:390
alterStatementSuffixLocation
    : KW_SET KW_LOCATION StringLiteral
    ;

// vendor/AlterClauseParser.g:399
alterStatementSuffixSkewedby
    : tableSkewed
    | KW_NOT KW_SKEWED
    | KW_NOT storedAsDirs
    ;

// vendor/AlterClauseParser.g:412
alterStatementSuffixExchangePartition
    : KW_EXCHANGE partitionSpec KW_WITH KW_TABLE tableName
    ;

// vendor/AlterClauseParser.g:419
alterStatementSuffixRenamePart
    : KW_RENAME KW_TO partitionSpec
    ;

// vendor/AlterClauseParser.g:426
alterStatementSuffixStatsPart
    : KW_UPDATE KW_STATISTICS KW_FOR KW_COLUMN? identifier KW_SET tableProperties ( KW_COMMENT StringLiteral )?
    ;

// vendor/AlterClauseParser.g:433
alterStatementSuffixMergeFiles
    : KW_CONCATENATE
    ;

// vendor/AlterClauseParser.g:441
alterStatementSuffixBucketNum
    : KW_INTO Number KW_BUCKETS
    ;

// vendor/AlterClauseParser.g:449
blocking
    : KW_AND KW_WAIT
    ;

// vendor/AlterClauseParser.g:454
compactPool
    : KW_POOL StringLiteral
    ;

// vendor/AlterClauseParser.g:459
alterStatementSuffixCompact
    : KW_COMPACT StringLiteral tableImplBuckets? blocking? compactPool? ( KW_WITH KW_OVERWRITE KW_TBLPROPERTIES tableProperties )? orderByClause?
    ;

// vendor/AlterClauseParser.g:466
alterStatementSuffixSetOwner
    : KW_SET KW_OWNER principalName
    ;

// vendor/AlterClauseParser.g:473
alterStatementSuffixSetPartSpec
    : KW_SET KW_PARTITION KW_SPEC LPAREN ( partitionTransformSpec ) RPAREN
    ;

// vendor/AlterClauseParser.g:480
alterStatementSuffixConvert
    : KW_CONVERT KW_TO identifier tablePropertiesPrefixed?
    ;

// vendor/AlterClauseParser.g:487
alterStatementSuffixExecute
    : KW_EXECUTE KW_ROLLBACK LPAREN ( ( StringLiteral | Number ) ) RPAREN
    | KW_EXECUTE KW_EXPIRE_SNAPSHOTS ( LPAREN ( StringLiteral ) RPAREN )?
    | KW_EXECUTE KW_SET_CURRENT_SNAPSHOT LPAREN ( expression ) RPAREN
    | KW_EXECUTE KW_FAST_FORWARD StringLiteral ( StringLiteral )?
    | KW_EXECUTE KW_CHERRY_PICK Number
    | KW_EXECUTE KW_EXPIRE_SNAPSHOTS KW_BETWEEN ( StringLiteral ) KW_AND ( StringLiteral )
    | KW_EXECUTE KW_EXPIRE_SNAPSHOTS KW_RETAIN KW_LAST Number
    | KW_EXECUTE KW_DELETE KW_ORPHAN_FILES ( KW_OLDER KW_THAN LPAREN ( StringLiteral ) RPAREN )?
    ;

// vendor/AlterClauseParser.g:508
alterStatementSuffixDropBranch
    : KW_DROP KW_BRANCH ifExists? identifier
    ;

// vendor/AlterClauseParser.g:515
alterStatementSuffixCreateBranch
    : KW_CREATE KW_BRANCH identifier snapshotIdOfRef? refRetain? retentionOfSnapshots?
    ;

// vendor/AlterClauseParser.g:522
snapshotIdOfRef
    : KW_FOR KW_SYSTEM_VERSION KW_AS KW_OF Number
    | ( KW_FOR KW_SYSTEM_TIME KW_AS KW_OF StringLiteral )
    | ( KW_FOR KW_TAG KW_AS KW_OF identifier )
    ;

// vendor/AlterClauseParser.g:535
refRetain
    : KW_RETAIN Number timeUnitQualifiers
    ;

// vendor/AlterClauseParser.g:542
retentionOfSnapshots
    : ( KW_WITH KW_SNAPSHOT KW_RETENTION Number KW_SNAPSHOTS ( Number timeUnitQualifiers )? )
    ;

// vendor/AlterClauseParser.g:549
alterStatementSuffixDropTag
    : KW_DROP KW_TAG ifExists? identifier
    ;

// vendor/AlterClauseParser.g:556
alterStatementSuffixCreateTag
    : KW_CREATE KW_TAG identifier snapshotIdOfRef? refRetain?
    ;

// vendor/AlterClauseParser.g:563
fileFormat
    : KW_INPUTFORMAT StringLiteral KW_OUTPUTFORMAT StringLiteral KW_SERDE StringLiteral ( KW_INPUTDRIVER StringLiteral KW_OUTPUTDRIVER StringLiteral )?
    | identifier
    ;

// vendor/AlterClauseParser.g:571
alterDataConnectorStatementSuffix
    : alterDataConnectorSuffixProperties
    | alterDataConnectorSuffixSetOwner
    | alterDataConnectorSuffixSetUrl
    ;

// vendor/AlterClauseParser.g:579
alterDataConnectorSuffixProperties
    : identifier KW_SET KW_DCPROPERTIES dcProperties
    ;

// vendor/AlterClauseParser.g:586
alterDataConnectorSuffixSetOwner
    : identifier KW_SET KW_OWNER principalName
    ;

// vendor/AlterClauseParser.g:593
alterDataConnectorSuffixSetUrl
    : identifier KW_SET KW_URL StringLiteral
    ;

// vendor/SelectClauseParser.g:48
selectClause
    : KW_SELECT QUERY_HINT? ( ( ( KW_ALL | KW_DISTINCT )? selectList ) | ( KW_TRANSFORM selectTrfmClause ) )
    | trfmClause
    ;

// vendor/SelectClauseParser.g:61
selectList
    : selectItem ( COMMA selectItem )*
    ;

// vendor/SelectClauseParser.g:68
selectTrfmClause
    : LPAREN selectExpressionList RPAREN rowFormat recordWriter KW_USING StringLiteral ( KW_AS ( ( LPAREN ( aliasList | columnNameTypeList ) RPAREN ) | ( aliasList | columnNameTypeList ) ) )? rowFormat recordReader
    ;

// vendor/SelectClauseParser.g:80
selectItem
    : tableAllColumns
    | ( expression ( ( KW_AS? identifier ) | ( KW_AS LPAREN identifier ( COMMA identifier )* RPAREN ) )? )
    ;

// vendor/SelectClauseParser.g:91
trfmClause
    : ( KW_MAP selectExpressionList | KW_REDUCE selectExpressionList ) rowFormat recordWriter KW_USING StringLiteral ( KW_AS ( ( LPAREN ( aliasList | columnNameTypeList ) RPAREN ) | ( aliasList | columnNameTypeList ) ) )? rowFormat recordReader
    ;

// vendor/SelectClauseParser.g:104
selectExpression
    : tableAllColumns
    | expression
    ;

// vendor/SelectClauseParser.g:113
selectExpressionList
    : selectExpression ( COMMA selectExpression )*
    ;

// vendor/SelectClauseParser.g:121
window_clause
    : KW_WINDOW window_defn ( COMMA window_defn )*
    ;

// vendor/SelectClauseParser.g:128
window_defn
    : identifier KW_AS window_specification
    ;

// vendor/SelectClauseParser.g:135
window_specification
    : ( identifier | ( LPAREN identifier? partitioningSpec? window_frame? RPAREN ) )
    ;

// vendor/SelectClauseParser.g:142
window_frame
    : window_range_expression
    | window_value_expression
    ;

// vendor/SelectClauseParser.g:147
window_range_expression
    : KW_ROWS window_frame_start_boundary
    | KW_ROWS KW_BETWEEN window_frame_boundary KW_AND window_frame_boundary
    ;

// vendor/SelectClauseParser.g:155
window_value_expression
    : KW_RANGE window_frame_start_boundary
    | KW_RANGE KW_BETWEEN window_frame_boundary KW_AND window_frame_boundary
    ;

// vendor/SelectClauseParser.g:163
window_frame_start_boundary
    : KW_UNBOUNDED KW_PRECEDING
    | KW_CURRENT KW_ROW
    | Number KW_PRECEDING
    ;

// vendor/SelectClauseParser.g:172
window_frame_boundary
    : KW_UNBOUNDED ( KW_PRECEDING | KW_FOLLOWING )
    | KW_CURRENT KW_ROW
    | Number ( KW_PRECEDING | KW_FOLLOWING )
    ;

// vendor/FromClauseParser.g:50
tableAllColumns
    : STAR
    | tableName DOT STAR
    ;

// vendor/FromClauseParser.g:58
tableOrColumn
    : identifier
    ;

// vendor/FromClauseParser.g:65
defaultValue
    : KW_DEFAULT
    ;

// vendor/FromClauseParser.g:72
expressionList
    : expression ( COMMA expression )*
    ;

// vendor/FromClauseParser.g:79
aliasList
    : identifier ( COMMA identifier )*
    ;

// vendor/FromClauseParser.g:88
fromClause
    : KW_FROM fromSource
    ;

// vendor/FromClauseParser.g:95
fromSource
    : uniqueJoinToken uniqueJoinSource ( COMMA uniqueJoinSource )+
    | joinSource
    ;

// vendor/FromClauseParser.g:105
atomjoinSource
    : tableSource ( lateralView )*
    | virtualTableSource ( lateralView )*
    | subQuerySource ( lateralView )*
    | subQuerySource ( lateralView )*
    | subQuerySource ( lateralView )*
    | partitionedTableFunction ( lateralView )*
    | LPAREN joinSource RPAREN
    ;

// vendor/FromClauseParser.g:117
// vendor/FromClauseParser.g:119: predicate dropped: {$joinToken.start.getType() != COMMA}?
// vendor/FromClauseParser.g:119: predicate dropped: {$joinToken.start.getType() != COMMA}?
joinSource
    : atomjoinSource ( joinToken joinSourcePart ( KW_ON expression | KW_USING columnParenthesesList )? )*
    ;

// vendor/FromClauseParser.g:122
joinSourcePart
    : ( tableSource | virtualTableSource | subQuerySource | partitionedTableFunction ) ( lateralView )*
    ;

// vendor/FromClauseParser.g:129
uniqueJoinSource
    : KW_PRESERVE? uniqueJoinTableSource uniqueJoinExpr
    ;

// vendor/FromClauseParser.g:135
uniqueJoinExpr
    : LPAREN expressionList RPAREN
    ;

// vendor/FromClauseParser.g:141
uniqueJoinToken
    : KW_UNIQUEJOIN
    ;

// vendor/FromClauseParser.g:146
joinToken
    : KW_JOIN
    | KW_INNER KW_JOIN
    | COMMA
    | KW_CROSS KW_JOIN
    | KW_LEFT ( KW_OUTER )? KW_JOIN
    | KW_RIGHT ( KW_OUTER )? KW_JOIN
    | KW_FULL ( KW_OUTER )? KW_JOIN
    | KW_LEFT KW_SEMI KW_JOIN
    | KW_LEFT KW_ANTI KW_JOIN
    ;

// vendor/FromClauseParser.g:161
lateralView
    : KW_LATERAL KW_VIEW KW_OUTER function tableAlias ( KW_AS identifier ( COMMA identifier )* )?
    | COMMA? KW_LATERAL KW_VIEW function tableAlias ( KW_AS identifier ( COMMA identifier )* )?
    | COMMA? KW_LATERAL KW_TABLE LPAREN valuesClause RPAREN KW_AS? tableAlias ( LPAREN identifier ( COMMA identifier )* RPAREN )?
    ;

// vendor/FromClauseParser.g:175
tableAlias
    : identifier
    ;

// vendor/FromClauseParser.g:182
tableBucketSample
    : KW_TABLESAMPLE LPAREN KW_BUCKET ( Number ) KW_OUT KW_OF ( Number ) ( KW_ON expression ( COMMA expression )* )? RPAREN
    ;

// vendor/FromClauseParser.g:189
splitSample
    : KW_TABLESAMPLE LPAREN ( Number ) ( KW_PERCENT | KW_ROWS ) RPAREN
    | KW_TABLESAMPLE LPAREN ( ByteLengthLiteral ) RPAREN
    ;

// vendor/FromClauseParser.g:201
tableSample
    : tableBucketSample
    | splitSample
    ;

// vendor/FromClauseParser.g:209
tableSource
    : tableName tableProperties? tableSample? ( asOfClause )? ( KW_AS? identifier )?
    ;

// vendor/FromClauseParser.g:216
asOfClause
    : ( KW_FOR KW_SYSTEM_TIME KW_AS KW_OF expression )
    | ( KW_FOR KW_SYSTEM_VERSION KW_AS KW_OF expression )
    ;

// vendor/FromClauseParser.g:227
uniqueJoinTableSource
    : tableName tableSample? ( KW_AS? identifier )?
    ;

// vendor/FromClauseParser.g:234
tableName
    : identifier DOT identifier ( DOT identifier )?
    | identifier
    ;

// vendor/FromClauseParser.g:247
viewName
    : ( identifier DOT )? identifier
    ;

// vendor/FromClauseParser.g:255
subQuerySource
    : LPAREN queryStatementExpression RPAREN KW_AS? identifier
    ;

// vendor/FromClauseParser.g:263
partitioningSpec
    : partitionByClause orderByClause?
    | orderByClause
    | distributeByClause sortByClause?
    | sortByClause
    | clusterByClause
    ;

// vendor/FromClauseParser.g:274
partitionTableFunctionSource
    : subQuerySource
    | tableSource
    | partitionedTableFunction
    ;

// vendor/FromClauseParser.g:283
partitionedTableFunction
    : identifier LPAREN KW_ON ( ( partitionTableFunctionSource partitioningSpec? ) ) ( Identifier LPAREN expression RPAREN ( COMMA Identifier LPAREN expression RPAREN )* )? ( ( RPAREN ) ) ( identifier )?
    ;

// vendor/FromClauseParser.g:296
whereClause
    : KW_WHERE searchCondition
    ;

// vendor/FromClauseParser.g:303
searchCondition
    : expression
    ;

// vendor/FromClauseParser.g:313
valuesSource
    : valuesClause
    ;

// vendor/FromClauseParser.g:330
valuesClause
    : KW_VALUES valuesTableConstructor
    ;

// vendor/FromClauseParser.g:337
valuesTableConstructor
    : ( valueRowConstructor ( COMMA valueRowConstructor )* )
    | firstValueRowConstructor ( COMMA valueRowConstructor )*
    ;

// vendor/FromClauseParser.g:346
valueRowConstructor
    : expressionsInParenthesis
    ;

// vendor/FromClauseParser.g:353
firstValueRowConstructor
    : LPAREN firstExpressionsWithAlias RPAREN
    ;

// vendor/FromClauseParser.g:364
virtualTableSource
    : KW_TABLE LPAREN valuesClause RPAREN KW_AS? tableAlias ( LPAREN identifier ( COMMA identifier )* )? RPAREN
    ;

// vendor/IdentifiersParser.g:69
groupByClause
    : KW_GROUP KW_BY groupby_expression
    ;

// vendor/IdentifiersParser.g:78
groupby_expression
    : rollupStandard
    | rollupOldSyntax
    | groupByEmpty
    ;

// vendor/IdentifiersParser.g:84
groupByEmpty
    : LPAREN RPAREN
    ;

// vendor/IdentifiersParser.g:90
rollupStandard
    : ( KW_ROLLUP | KW_CUBE ) LPAREN expression ( COMMA expression )* RPAREN
    ;

// vendor/IdentifiersParser.g:101
rollupOldSyntax
    : expressionsNotInParenthesis ( ( KW_WITH KW_ROLLUP ) | ( KW_WITH KW_CUBE ) )? ( KW_GROUPING KW_SETS LPAREN groupingSetExpression ( COMMA groupingSetExpression )* RPAREN )?
    ;

// vendor/IdentifiersParser.g:116
groupingSetExpression
    : groupingSetExpressionMultiple
    | groupingExpressionSingle
    ;

// vendor/IdentifiersParser.g:125
groupingSetExpressionMultiple
    : LPAREN expression? ( COMMA expression )* RPAREN
    ;

// vendor/IdentifiersParser.g:135
groupingExpressionSingle
    : expression
    ;

// vendor/IdentifiersParser.g:142
havingClause
    : KW_HAVING havingCondition
    ;

// vendor/IdentifiersParser.g:149
qualifyClause
    : KW_QUALIFY expression
    ;

// vendor/IdentifiersParser.g:156
havingCondition
    : expression
    ;

// vendor/IdentifiersParser.g:163
expressionsInParenthesis
    : LPAREN expressionsNotInParenthesis RPAREN
    ;

// vendor/IdentifiersParser.g:168
expressionsNotInParenthesis
    : expressionOrDefault expressionPart?
    ;

// vendor/IdentifiersParser.g:178
expressionPart
    : ( COMMA expressionOrDefault )+
    ;

// vendor/IdentifiersParser.g:185
expressionOrDefault
    : defaultValue
    | expression
    ;

// vendor/IdentifiersParser.g:193
firstExpressionsWithAlias
    : expression KW_AS? identifier? ( COMMA expressionWithAlias )*
    ;

// vendor/IdentifiersParser.g:203
expressionWithAlias
    : expression KW_AS? identifier?
    ;

// vendor/IdentifiersParser.g:210
expressions
    : expressionsInParenthesis
    | expressionsNotInParenthesis
    ;

// vendor/IdentifiersParser.g:217
columnRefOrderInParenthesis
    : LPAREN columnRefOrder ( COMMA columnRefOrder )* RPAREN
    ;

// vendor/IdentifiersParser.g:222
columnRefOrderNotInParenthesis
    : columnRefOrder ( COMMA columnRefOrder )*
    ;

// vendor/IdentifiersParser.g:228
orderByClause
    : KW_ORDER KW_BY columnRefOrder ( COMMA columnRefOrder )*
    ;

// vendor/IdentifiersParser.g:235
clusterByClause
    : KW_CLUSTER KW_BY expressions
    ;

// vendor/IdentifiersParser.g:242
partitionByClause
    : KW_PARTITION KW_BY expressions
    ;

// vendor/IdentifiersParser.g:249
distributeByClause
    : KW_DISTRIBUTE KW_BY expressions
    ;

// vendor/IdentifiersParser.g:256
sortByClause
    : KW_SORT KW_BY ( columnRefOrderInParenthesis | columnRefOrderNotInParenthesis )
    ;

// vendor/IdentifiersParser.g:269
trimFunction
    : KW_TRIM LPAREN ( KW_LEADING | KW_TRAILING | KW_BOTH )? ( selectExpression )? KW_FROM ( selectExpression ) RPAREN
    ;

// vendor/IdentifiersParser.g:278
function
    : ( trimFunction )
    | functionName LPAREN ( ( STAR ) | ( KW_DISTINCT | KW_ALL )? ( selectExpression ( COMMA selectExpression )* )? ) ( ( RPAREN ( KW_WITHIN KW_GROUP LPAREN orderByClause RPAREN ) ) | ( RPAREN ( ( null_treatment )? ( KW_OVER window_specification ) ) ) | ( null_treatment ) RPAREN ( KW_OVER window_specification ) | RPAREN )
    ;

// vendor/IdentifiersParser.g:306
null_treatment
    : KW_RESPECT KW_NULLS
    | KW_IGNORE KW_NULLS
    ;

// vendor/IdentifiersParser.g:313
functionName
    : functionIdentifier
    | sql11ReservedKeywordsUsedAsFunctionName
    ;

// vendor/IdentifiersParser.g:322
castExpression
    : KW_CAST LPAREN expression KW_AS type ( KW_FORMAT StringLiteral )? RPAREN
    ;

// vendor/IdentifiersParser.g:344
caseExpression
    : KW_CASE expression ( KW_WHEN expression KW_THEN expression )+ ( KW_ELSE expression )? KW_END
    ;

// vendor/IdentifiersParser.g:354
whenExpression
    : KW_CASE ( KW_WHEN expression KW_THEN expression )+ ( KW_ELSE expression )? KW_END
    ;

// vendor/IdentifiersParser.g:364
floorExpression
    : KW_FLOOR LPAREN expression ( KW_TO ( floorDateQualifiers ) )? RPAREN
    ;

// vendor/IdentifiersParser.g:376
floorDateQualifiers
    : KW_YEAR
    | KW_QUARTER
    | KW_MONTH
    | KW_WEEK
    | KW_DAY
    | KW_HOUR
    | KW_MINUTE
    | KW_SECOND
    ;

// vendor/IdentifiersParser.g:388
extractExpression
    : KW_EXTRACT LPAREN ( timeQualifiers ) KW_FROM expression RPAREN
    ;

// vendor/IdentifiersParser.g:398
timeQualifiers
    : KW_YEAR
    | KW_QUARTER
    | KW_MONTH
    | KW_WEEK
    | KW_DAY
    | KW_DOW
    | KW_HOUR
    | KW_MINUTE
    | KW_SECOND
    ;

// vendor/IdentifiersParser.g:411
timeUnitQualifiers
    : KW_DAY
    | KW_HOUR
    | KW_MINUTE
    ;

// vendor/IdentifiersParser.g:418
constant
    : intervalLiteral
    | Number
    | dateLiteral
    | timestampLiteral
    | timestampLocalTZLiteral
    | StringLiteral
    | stringLiteralSequence
    | IntegralLiteral
    | NumberLiteral
    | charSetStringLiteral
    | booleanValue
    | KW_NULL
    | prepareStmtParam
    ;

// vendor/IdentifiersParser.g:437
prepareStmtParam
    : parameterIdx
    ;

// vendor/IdentifiersParser.g:445
parameterIdx
    : QUESTION
    ;

// vendor/IdentifiersParser.g:455
stringLiteralSequence
    : StringLiteral StringLiteral+
    ;

// vendor/IdentifiersParser.g:460
charSetStringLiteral
    : CharSetName CharSetLiteral
    ;

// vendor/IdentifiersParser.g:467
dateLiteral
    : KW_DATE StringLiteral
    | KW_CURRENT_DATE
    ;

// vendor/IdentifiersParser.g:479
timestampLiteral
    : KW_TIMESTAMP StringLiteral
    | KW_CURRENT_TIMESTAMP
    ;

// vendor/IdentifiersParser.g:489
timestampLocalTZLiteral
    : KW_TIMESTAMPLOCALTZ StringLiteral
    ;

// vendor/IdentifiersParser.g:497
intervalValue
    : StringLiteral
    | Number
    ;

// vendor/IdentifiersParser.g:502
intervalLiteral
    : intervalValue intervalQualifiers
    ;

// vendor/IdentifiersParser.g:508
intervalExpression
    : LPAREN intervalValue RPAREN intervalQualifiers
    | KW_INTERVAL intervalValue intervalQualifiers
    | KW_INTERVAL LPAREN expression RPAREN intervalQualifiers
    ;

// vendor/IdentifiersParser.g:520
intervalQualifiers
    : KW_YEAR KW_TO KW_MONTH
    | KW_DAY KW_TO KW_SECOND
    | KW_YEAR
    | KW_MONTH
    | KW_DAY
    | KW_HOUR
    | KW_MINUTE
    | KW_SECOND
    ;

// vendor/IdentifiersParser.g:532
expression
    : precedenceOrExpression
    ;

// vendor/IdentifiersParser.g:539
atomExpression
    : constant
    | intervalExpression
    | castExpression
    | extractExpression
    | floorExpression
    | caseExpression
    | whenExpression
    | ( subQueryExpression )
    | function
    | tableOrColumn
    | expressionsInParenthesis
    ;

// vendor/IdentifiersParser.g:555
precedenceFieldExpression
    : atomExpression ( ( LSQUARE expression RSQUARE ) | ( DOT identifier ) )*
    ;

// vendor/IdentifiersParser.g:560
precedenceUnaryOperator
    : PLUS
    | MINUS
    | TILDE
    ;

// vendor/IdentifiersParser.g:565
precedenceUnaryPrefixExpression
    : ( precedenceUnaryOperator )* precedenceFieldExpression
    ;

// vendor/IdentifiersParser.g:571
precedenceBitwiseXorOperator
    : BITWISEXOR
    ;

// vendor/IdentifiersParser.g:576
precedenceBitwiseXorExpression
    : precedenceUnaryPrefixExpression ( precedenceBitwiseXorOperator precedenceUnaryPrefixExpression )*
    ;

// vendor/IdentifiersParser.g:582
precedenceStarOperator
    : STAR
    | DIVIDE
    | MOD
    | DIV
    ;

// vendor/IdentifiersParser.g:587
precedenceStarExpression
    : precedenceBitwiseXorExpression ( precedenceStarOperator precedenceBitwiseXorExpression )*
    ;

// vendor/IdentifiersParser.g:593
precedencePlusOperator
    : PLUS
    | MINUS
    ;

// vendor/IdentifiersParser.g:598
precedencePlusExpression
    : precedenceStarExpression ( precedencePlusOperator precedenceStarExpression )*
    ;

// vendor/IdentifiersParser.g:603
precedenceConcatenateOperator
    : CONCATENATE
    ;

// vendor/IdentifiersParser.g:608
precedenceConcatenateExpression
    : ( precedencePlusExpression ) ( precedenceConcatenateOperator precedencePlusExpression )*
    ;

// vendor/IdentifiersParser.g:618
precedenceAmpersandOperator
    : AMPERSAND
    ;

// vendor/IdentifiersParser.g:623
precedenceAmpersandExpression
    : precedenceConcatenateExpression ( precedenceAmpersandOperator precedenceConcatenateExpression )*
    ;

// vendor/IdentifiersParser.g:629
precedenceBitwiseOrOperator
    : BITWISEOR
    ;

// vendor/IdentifiersParser.g:634
precedenceBitwiseOrExpression
    : precedenceAmpersandExpression ( precedenceBitwiseOrOperator precedenceAmpersandExpression )*
    ;

// vendor/IdentifiersParser.g:640
precedenceRegexpOperator
    : KW_LIKE
    | KW_RLIKE
    | KW_REGEXP
    ;

// vendor/IdentifiersParser.g:645
precedenceSimilarOperator
    : precedenceRegexpOperator
    | LESSTHANOREQUALTO
    | LESSTHAN
    | GREATERTHANOREQUALTO
    | GREATERTHAN
    ;

// vendor/IdentifiersParser.g:650
subQueryExpression
    : LPAREN selectStatement RPAREN
    ;

// vendor/IdentifiersParser.g:655
precedenceSimilarExpression
    : precedenceSimilarExpressionMain
    | KW_EXISTS subQueryExpression
    ;

// vendor/IdentifiersParser.g:662
precedenceSimilarExpressionMain
    : precedenceBitwiseOrExpression precedenceSimilarExpressionPart?
    ;

// vendor/IdentifiersParser.g:669
precedenceSimilarExpressionPart
    : ( precedenceSimilarOperator precedenceBitwiseOrExpression )
    | precedenceSimilarExpressionAtom
    | ( KW_NOT precedenceSimilarExpressionPartNot )
    ;

// vendor/IdentifiersParser.g:679
precedenceSimilarExpressionAtom
    : KW_IN precedenceSimilarExpressionIn
    | KW_BETWEEN ( precedenceBitwiseOrExpression ) KW_AND ( precedenceBitwiseOrExpression )
    | KW_LIKE KW_ANY ( expressionsInParenthesis )
    | KW_LIKE KW_ALL ( expressionsInParenthesis )
    | precedenceSimilarExpressionQuantifierPredicate
    ;

// vendor/IdentifiersParser.g:695
precedenceSimilarExpressionQuantifierPredicate
    : subQuerySelectorOperator quantifierType subQueryExpression
    ;

// vendor/IdentifiersParser.g:701
quantifierType
    : KW_ANY
    | KW_SOME
    | KW_ALL
    ;

// vendor/IdentifiersParser.g:710
precedenceSimilarExpressionIn
    : subQueryExpression
    | expressionsInParenthesis
    ;

// vendor/IdentifiersParser.g:718
precedenceSimilarExpressionPartNot
    : precedenceRegexpOperator precedenceBitwiseOrExpression
    | precedenceSimilarExpressionAtom
    ;

// vendor/IdentifiersParser.g:726
precedenceDistinctOperator
    : KW_IS KW_DISTINCT KW_FROM
    ;

// vendor/IdentifiersParser.g:731
precedenceEqualOperator
    : EQUAL
    | EQUAL_NS
    | NOTEQUAL
    | KW_IS KW_NOT KW_DISTINCT KW_FROM
    ;

// vendor/IdentifiersParser.g:736
precedenceEqualExpression
    : ( precedenceSimilarExpression ) ( precedenceEqualOperator precedenceSimilarExpression | precedenceDistinctOperator precedenceSimilarExpression )*
    ;

// vendor/IdentifiersParser.g:749
isCondition
    : KW_NULL
    | KW_TRUE
    | KW_FALSE
    | KW_UNKNOWN
    | KW_NOT KW_NULL
    | KW_NOT KW_TRUE
    | KW_NOT KW_FALSE
    | KW_NOT KW_UNKNOWN
    ;

// vendor/IdentifiersParser.g:760
precedenceUnarySuffixExpression
    : precedenceEqualExpression ( KW_IS isCondition )?
    ;

// vendor/IdentifiersParser.g:766
precedenceNotOperator
    : KW_NOT
    ;

// vendor/IdentifiersParser.g:771
precedenceNotExpression
    : ( precedenceNotOperator )* precedenceUnarySuffixExpression
    ;

// vendor/IdentifiersParser.g:777
precedenceAndOperator
    : KW_AND
    ;

// vendor/IdentifiersParser.g:782
precedenceAndExpression
    : precedenceNotExpression ( precedenceAndOperator precedenceNotExpression )*
    ;

// vendor/IdentifiersParser.g:788
precedenceOrOperator
    : KW_OR
    ;

// vendor/IdentifiersParser.g:793
precedenceOrExpression
    : precedenceAndExpression ( precedenceOrOperator precedenceAndExpression )*
    ;

// vendor/IdentifiersParser.g:799
booleanValue
    : KW_TRUE
    | KW_FALSE
    ;

// vendor/IdentifiersParser.g:804
booleanValueTok
    : KW_TRUE
    | KW_FALSE
    ;

// vendor/IdentifiersParser.g:810
tableOrPartition
    : tableName partitionSpec?
    ;

// vendor/IdentifiersParser.g:815
partitionSpec
    : KW_PARTITION LPAREN partitionVal ( COMMA partitionVal )* RPAREN
    ;

// vendor/IdentifiersParser.g:821
partitionVal
    : identifier ( EQUAL constant )?
    ;

// vendor/IdentifiersParser.g:826
partitionSelectorSpec
    : LPAREN partitionSelectorVal ( COMMA partitionSelectorVal )* RPAREN
    ;

// vendor/IdentifiersParser.g:831
partitionSelectorVal
    : identifier partitionSelectorOperator constant
    ;

// vendor/IdentifiersParser.g:836
partitionSelectorOperator
    : KW_LIKE
    | subQuerySelectorOperator
    ;

// vendor/IdentifiersParser.g:841
subQuerySelectorOperator
    : EQUAL
    | NOTEQUAL
    | LESSTHANOREQUALTO
    | LESSTHAN
    | GREATERTHANOREQUALTO
    | GREATERTHAN
    ;

// vendor/IdentifiersParser.g:846
sysFuncNames
    : KW_AND
    | KW_OR
    | KW_NOT
    | KW_LIKE
    | KW_IF
    | KW_CASE
    | KW_WHEN
    | KW_FLOOR
    | KW_TINYINT
    | KW_SMALLINT
    | KW_INT
    | KW_BIGINT
    | KW_FLOAT
    | KW_REAL
    | KW_DOUBLE
    | KW_BOOLEAN
    | KW_STRING
    | KW_BINARY
    | KW_ARRAY
    | KW_MAP
    | KW_STRUCT
    | KW_UNIONTYPE
    | EQUAL
    | EQUAL_NS
    | NOTEQUAL
    | LESSTHANOREQUALTO
    | LESSTHAN
    | GREATERTHANOREQUALTO
    | GREATERTHAN
    | DIVIDE
    | PLUS
    | MINUS
    | STAR
    | MOD
    | DIV
    | AMPERSAND
    | TILDE
    | BITWISEOR
    | BITWISEXOR
    | KW_RLIKE
    | KW_REGEXP
    | KW_IN
    | KW_BETWEEN
    ;

// vendor/IdentifiersParser.g:893
descFuncNames
    : sysFuncNames
    | StringLiteral
    | functionIdentifier
    ;

// vendor/IdentifiersParser.g:900
identifier
    : Identifier
    | nonReserved
    ;

// vendor/IdentifiersParser.g:906
functionIdentifier
    : identifier DOT identifier
    | identifier
    ;

// vendor/IdentifiersParser.g:915
principalIdentifier
    : identifier
    | QuotedIdentifier
    ;

// vendor/IdentifiersParser.g:931
nonReserved
    : KW_ABORT
    | KW_ADD
    | KW_ADMIN
    | KW_AFTER
    | KW_ANALYZE
    | KW_ARCHIVE
    | KW_ASC
    | KW_BEFORE
    | KW_BUCKET
    | KW_BUCKETS
    | KW_CASCADE
    | KW_CBO
    | KW_CHANGE
    | KW_CHECK
    | KW_CLUSTER
    | KW_CLUSTERED
    | KW_CLUSTERSTATUS
    | KW_COLLECTION
    | KW_COLUMNS
    | KW_COMMENT
    | KW_COMPACT
    | KW_COMPACTIONS
    | KW_COMPUTE
    | KW_CONCATENATE
    | KW_CONTINUE
    | KW_COST
    | KW_DATA
    | KW_DAY
    | KW_DATABASES
    | KW_DATETIME
    | KW_DBPROPERTIES
    | KW_DCPROPERTIES
    | KW_DEFERRED
    | KW_DEFINED
    | KW_DELIMITED
    | KW_DEPENDENCY
    | KW_DESC
    | KW_DIRECTORIES
    | KW_DIRECTORY
    | KW_DISABLE
    | KW_DISTRIBUTE
    | KW_DISTRIBUTED
    | KW_DOW
    | KW_ELEM_TYPE
    | KW_ENABLE
    | KW_ENFORCED
    | KW_ESCAPED
    | KW_EXCLUSIVE
    | KW_EXPLAIN
    | KW_EXPORT
    | KW_FIELDS
    | KW_FILE
    | KW_FILEFORMAT
    | KW_FIRST
    | KW_FORMAT
    | KW_FORMATTED
    | KW_FUNCTIONS
    | KW_HOLD_DDLTIME
    | KW_HOUR
    | KW_IDXPROPERTIES
    | KW_RESPECT
    | KW_IGNORE
    | KW_INDEX
    | KW_INDEXES
    | KW_INPATH
    | KW_INPUTDRIVER
    | KW_INPUTFORMAT
    | KW_ITEMS
    | KW_JAR
    | KW_JOINCOST
    | KW_KILL
    | KW_KEYS
    | KW_KEY_TYPE
    | KW_LAST
    | KW_LIMIT
    | KW_OFFSET
    | KW_LINES
    | KW_LOAD
    | KW_LOCATION
    | KW_LOCK
    | KW_LOCKS
    | KW_LOGICAL
    | KW_LONG
    | KW_MANAGED
    | KW_MANAGEDLOCATION
    | KW_MAPJOIN
    | KW_MATERIALIZED
    | KW_METADATA
    | KW_MINUTE
    | KW_MONTH
    | KW_MSCK
    | KW_NOSCAN
    | KW_NO_DROP
    | KW_NULLS
    | KW_OFFLINE
    | KW_OPTION
    | KW_OUTPUTDRIVER
    | KW_OUTPUTFORMAT
    | KW_OVERWRITE
    | KW_OWNER
    | KW_PARTITIONED
    | KW_PARTITIONS
    | KW_PLUS
    | KW_PRINCIPALS
    | KW_PROTECTION
    | KW_PURGE
    | KW_QUERY
    | KW_QUARTER
    | KW_READ
    | KW_READONLY
    | KW_REBUILD
    | KW_RECORDREADER
    | KW_RECORDWRITER
    | KW_RELOAD
    | KW_REMOTE
    | KW_RENAME
    | KW_REPAIR
    | KW_REPLACE
    | KW_REPLICATION
    | KW_RESTRICT
    | KW_REWRITE
    | KW_ROLE
    | KW_ROLES
    | KW_SCHEMA
    | KW_SCHEMAS
    | KW_SECOND
    | KW_SEMI
    | KW_SERDE
    | KW_SERDEPROPERTIES
    | KW_SERVER
    | KW_SETS
    | KW_SHARED
    | KW_SHOW
    | KW_SHOW_DATABASE
    | KW_SKEWED
    | KW_SORT
    | KW_SORTED
    | KW_SSL
    | KW_STATISTICS
    | KW_STORED
    | KW_AST
    | KW_STREAMTABLE
    | KW_STRING
    | KW_STRUCT
    | KW_TABLES
    | KW_TBLPROPERTIES
    | KW_TEMPORARY
    | KW_TERMINATED
    | KW_TINYINT
    | KW_TOUCH
    | KW_TRANSACTIONAL
    | KW_TRANSACTIONS
    | KW_TYPE
    | KW_UNARCHIVE
    | KW_UNDO
    | KW_UNIONTYPE
    | KW_UNLOCK
    | KW_UNSET
    | KW_UNSIGNED
    | KW_URI
    | KW_URL
    | KW_USE
    | KW_UTC
    | KW_UTCTIMESTAMP
    | KW_VALUE_TYPE
    | KW_VIEW
    | KW_WEEK
    | KW_WHILE
    | KW_YEAR
    | KW_WORK
    | KW_TRANSACTION
    | KW_WRITE
    | KW_ISOLATION
    | KW_LEVEL
    | KW_SNAPSHOT
    | KW_AUTOCOMMIT
    | KW_RELY
    | KW_NORELY
    | KW_VALIDATE
    | KW_NOVALIDATE
    | KW_KEY
    | KW_MATCHED
    | KW_REPL
    | KW_DUMP
    | KW_BATCH
    | KW_STATUS
    | KW_CACHE
    | KW_DAYOFWEEK
    | KW_VIEWS
    | KW_VECTORIZATION
    | KW_SUMMARY
    | KW_OPERATOR
    | KW_EXPRESSION
    | KW_DETAIL
    | KW_DEBUG
    | KW_WAIT
    | KW_ZONE
    | KW_TIMESTAMPTZ
    | KW_DEFAULT
    | KW_REOPTIMIZATION
    | KW_EXECUTED
    | KW_SCHEDULED
    | KW_CRON
    | KW_EVERY
    | KW_AT
    | KW_EXECUTE
    | KW_RESOURCE
    | KW_PLAN
    | KW_PLANS
    | KW_QUERY_PARALLELISM
    | KW_ACTIVATE
    | KW_MOVE
    | KW_DO
    | KW_POOL
    | KW_ALLOC_FRACTION
    | KW_SCHEDULING_POLICY
    | KW_PATH
    | KW_MAPPING
    | KW_WORKLOAD
    | KW_MANAGEMENT
    | KW_ACTIVE
    | KW_UNMANAGED
    | KW_UNKNOWN
    | KW_WITHIN
    | KW_TRIM
    | KW_SPEC
    | KW_SYSTEM_TIME
    | KW_SYSTEM_VERSION
    | KW_EXPIRE_SNAPSHOTS
    | KW_SET_CURRENT_SNAPSHOT
    | KW_BRANCH
    | KW_SNAPSHOTS
    | KW_RETAIN
    | KW_RETENTION
    | KW_TAG
    | KW_FAST_FORWARD
    | KW_OPTIMIZE
    ;

// vendor/IdentifiersParser.g:994
sql11ReservedKeywordsUsedAsFunctionName
    : KW_IF
    | KW_ARRAY
    | KW_MAP
    | KW_BIGINT
    | KW_BINARY
    | KW_BOOLEAN
    | KW_CURRENT_DATE
    | KW_CURRENT_TIMESTAMP
    | KW_DATE
    | KW_DOUBLE
    | KW_FLOAT
    | KW_REAL
    | KW_GROUPING
    | KW_INT
    | KW_SMALLINT
    | KW_TIMESTAMP
    ;

// vendor/ResourcePlanParser.g:46
resourcePlanDdlStatements
    : createResourcePlanStatement
    | alterResourcePlanStatement
    | dropResourcePlanStatement
    | globalWmStatement
    | replaceResourcePlanStatement
    | createTriggerStatement
    | alterTriggerStatement
    | dropTriggerStatement
    | createPoolStatement
    | alterPoolStatement
    | dropPoolStatement
    | createMappingStatement
    | alterMappingStatement
    | dropMappingStatement
    ;

// vendor/ResourcePlanParser.g:63
rpAssign
    : ( ( KW_QUERY_PARALLELISM EQUAL Number ) | ( KW_DEFAULT KW_POOL EQUAL poolPath ) )
    ;

// vendor/ResourcePlanParser.g:72
rpAssignList
    : rpAssign ( COMMA rpAssign )*
    ;

// vendor/ResourcePlanParser.g:78
rpUnassign
    : ( ( KW_QUERY_PARALLELISM ) | ( KW_DEFAULT KW_POOL ) )
    ;

// vendor/ResourcePlanParser.g:87
rpUnassignList
    : rpUnassign ( COMMA rpUnassign )*
    ;

// vendor/ResourcePlanParser.g:93
createResourcePlanStatement
    : KW_CREATE KW_RESOURCE KW_PLAN ifNotExists? ( ( identifier KW_LIKE identifier ) | ( identifier ( KW_WITH rpAssignList )? ) )
    ;

// vendor/ResourcePlanParser.g:103
withReplace
    : KW_WITH KW_REPLACE
    ;

// vendor/ResourcePlanParser.g:104
activate
    : KW_ACTIVATE withReplace?
    ;

// vendor/ResourcePlanParser.g:105
enable
    : KW_ENABLE
    ;

// vendor/ResourcePlanParser.g:106
disable
    : KW_DISABLE
    ;

// vendor/ResourcePlanParser.g:107
unmanaged
    : KW_UNMANAGED
    ;

// vendor/ResourcePlanParser.g:110
alterResourcePlanStatement
    : KW_ALTER KW_RESOURCE KW_PLAN identifier ( ( KW_VALIDATE ) | ( KW_DISABLE ) | ( KW_SET rpAssignList ) | ( KW_UNSET rpUnassignList ) | ( KW_RENAME KW_TO identifier ) | ( ( activate enable? | enable activate? ) ) )
    ;

// vendor/ResourcePlanParser.g:125
globalWmStatement
    : KW_ENABLE KW_WORKLOAD KW_MANAGEMENT
    | KW_DISABLE KW_WORKLOAD KW_MANAGEMENT
    ;

// vendor/ResourcePlanParser.g:132
replaceResourcePlanStatement
    : KW_REPLACE ( ( KW_ACTIVE KW_RESOURCE KW_PLAN KW_WITH identifier ) | ( KW_RESOURCE KW_PLAN identifier KW_WITH identifier ) )
    ;

// vendor/ResourcePlanParser.g:141
dropResourcePlanStatement
    : KW_DROP KW_RESOURCE KW_PLAN ifExists? identifier
    ;

// vendor/ResourcePlanParser.g:147
poolPath
    : identifier ( DOT identifier )*
    ;

// vendor/ResourcePlanParser.g:153
triggerExpression
    : triggerAtomExpression
    ;

// vendor/ResourcePlanParser.g:159
triggerExpressionStandalone
    : triggerExpression EOF
    ;

// vendor/ResourcePlanParser.g:165
triggerOrExpression
    : triggerAndExpression ( KW_OR triggerAndExpression )*
    ;

// vendor/ResourcePlanParser.g:171
triggerAndExpression
    : triggerAtomExpression ( KW_AND triggerAtomExpression )*
    ;

// vendor/ResourcePlanParser.g:177
triggerAtomExpression
    : identifier comparisionOperator triggerLiteral
    ;

// vendor/ResourcePlanParser.g:184
triggerLiteral
    : Number
    | StringLiteral
    ;

// vendor/ResourcePlanParser.g:191
comparisionOperator
    : GREATERTHAN
    ;

// vendor/ResourcePlanParser.g:197
triggerActionExpression
    : KW_KILL
    | ( KW_MOVE KW_TO poolPath )
    ;

// vendor/ResourcePlanParser.g:204
triggerActionExpressionStandalone
    : triggerActionExpression EOF
    ;

// vendor/ResourcePlanParser.g:206
createTriggerStatement
    : KW_CREATE KW_TRIGGER identifier DOT identifier KW_WHEN triggerExpression KW_DO triggerActionExpression
    ;

// vendor/ResourcePlanParser.g:214
alterTriggerStatement
    : KW_ALTER KW_TRIGGER identifier DOT identifier ( ( KW_WHEN triggerExpression KW_DO triggerActionExpression ) | ( KW_ADD KW_TO KW_POOL poolPath ) | ( KW_DROP KW_FROM KW_POOL poolPath ) | ( KW_ADD KW_TO KW_UNMANAGED ) | ( KW_DROP KW_FROM KW_UNMANAGED ) )
    ;

// vendor/ResourcePlanParser.g:228
dropTriggerStatement
    : KW_DROP KW_TRIGGER identifier DOT identifier
    ;

// vendor/ResourcePlanParser.g:235
poolAssign
    : ( ( KW_ALLOC_FRACTION EQUAL Number ) | ( KW_QUERY_PARALLELISM EQUAL Number ) | ( KW_SCHEDULING_POLICY EQUAL StringLiteral ) | ( KW_PATH EQUAL poolPath ) )
    ;

// vendor/ResourcePlanParser.g:246
poolAssignList
    : poolAssign ( COMMA poolAssign )*
    ;

// vendor/ResourcePlanParser.g:252
createPoolStatement
    : KW_CREATE KW_POOL identifier DOT poolPath KW_WITH poolAssignList
    ;

// vendor/ResourcePlanParser.g:260
alterPoolStatement
    : KW_ALTER KW_POOL identifier DOT poolPath ( ( KW_SET poolAssignList ) | ( KW_UNSET KW_SCHEDULING_POLICY ) | ( KW_ADD KW_TRIGGER identifier ) | ( KW_DROP KW_TRIGGER identifier ) )
    ;

// vendor/ResourcePlanParser.g:271
dropPoolStatement
    : KW_DROP KW_POOL identifier DOT poolPath
    ;

// vendor/ResourcePlanParser.g:278
createMappingStatement
    : ( KW_CREATE ( KW_USER | KW_GROUP | KW_APPLICATION ) KW_MAPPING StringLiteral KW_IN identifier ( ( KW_TO poolPath ) | unmanaged ) ( KW_WITH KW_ORDER Number )? )
    ;

// vendor/ResourcePlanParser.g:288
alterMappingStatement
    : ( KW_ALTER ( KW_USER | KW_GROUP | KW_APPLICATION ) KW_MAPPING StringLiteral KW_IN identifier ( ( KW_TO poolPath ) | unmanaged ) ( KW_WITH KW_ORDER Number )? )
    ;

// vendor/ResourcePlanParser.g:298
dropMappingStatement
    : KW_DROP ( KW_USER | KW_GROUP | KW_APPLICATION ) KW_MAPPING StringLiteral KW_IN identifier
    ;

// vendor/CreateDDLParser.g:46
likeTableOrFile
    : KW_LIKE KW_FILE
    | ( KW_LIKE KW_FILE identifier StringLiteral )
    | ( KW_LIKE tableName )
    ;

// vendor/CreateDDLParser.g:53
createTableStatement
    : KW_CREATE ( KW_TEMPORARY )? ( KW_TRANSACTIONAL )? ( KW_EXTERNAL )? KW_TABLE ifNotExists? tableName ( likeTableOrFile createTablePartitionSpec? tableRowFormat? tableFileFormat? tableLocation? tablePropertiesPrefixed? | ( LPAREN columnNameTypeOrConstraintList RPAREN )? tableComment? createTablePartitionSpec? tableBuckets? tableSkewed? tableRowFormat? tableFileFormat? tableLocation? tablePropertiesPrefixed? ( KW_AS selectStatementWithCTE )? )
    | KW_CREATE KW_MANAGED KW_TABLE ifNotExists? tableName ( likeTableOrFile tableRowFormat? tableFileFormat? tableLocation? tablePropertiesPrefixed? | ( LPAREN columnNameTypeOrConstraintList RPAREN )? tableComment? createTablePartitionSpec? tableBuckets? tableSkewed? tableRowFormat? tableFileFormat? tableLocation? tablePropertiesPrefixed? ( KW_AS selectStatementWithCTE )? )
    ;

// vendor/CreateDDLParser.g:119
createDataConnectorStatement
    : KW_CREATE KW_DATACONNECTOR ifNotExists? identifier dataConnectorType dataConnectorUrl dataConnectorComment? ( KW_WITH KW_DCPROPERTIES dcProperties )?
    ;

// vendor/CreateDDLParser.g:126
dataConnectorComment
    : KW_COMMENT StringLiteral
    ;

// vendor/CreateDDLParser.g:133
dataConnectorUrl
    : KW_URL StringLiteral
    ;

// vendor/CreateDDLParser.g:140
dataConnectorType
    : KW_TYPE StringLiteral
    ;

// vendor/CreateDDLParser.g:147
dcProperties
    : LPAREN dbPropertiesList RPAREN
    ;

// vendor/CreateDDLParser.g:154
dropDataConnectorStatement
    : KW_DROP ( KW_DATACONNECTOR ) ifExists? identifier
    ;

// vendor/PrepareStatementParser.g:47
prepareStatement
    : KW_PREPARE identifier KW_FROM queryStatementExpression
    ;

// vendor/PrepareStatementParser.g:54
executeStatement
    : KW_EXECUTE identifier KW_USING executeParamList
    ;

// vendor/PrepareStatementParser.g:62
executeParamList
    : constant ( COMMA constant )*
    ;

// vendor/HiveLexer.g:30
StringLiteral
    : ( '"' ( ~( '"' | '\\' ) | ( '\\' . ) )* '"' | '\'' ( ~( '\'' | '\\' ) | ( '\\' . ) )* '\'' )+
    ;

// vendor/HiveLexer.g:36
fragment QuotedIdentifier
    : ( '`' ( '``' | ~( '`' ) )+ '`' )
    ;

// vendor/HiveLexerParent.g:21
KW_TRUE
    : 'TRUE'
    ;

// vendor/HiveLexerParent.g:22
KW_FALSE
    : 'FALSE'
    ;

// vendor/HiveLexerParent.g:23
KW_UNKNOWN
    : 'UNKNOWN'
    ;

// vendor/HiveLexerParent.g:24
KW_ALL
    : 'ALL'
    ;

// vendor/HiveLexerParent.g:25
KW_SOME
    : 'SOME'
    ;

// vendor/HiveLexerParent.g:26
KW_NONE
    : 'NONE'
    ;

// vendor/HiveLexerParent.g:27
KW_AND
    : 'AND'
    ;

// vendor/HiveLexerParent.g:28
KW_OR
    : 'OR'
    ;

// vendor/HiveLexerParent.g:29
KW_NOT
    : 'NOT'
    | '!'
    ;

// vendor/HiveLexerParent.g:30
KW_LIKE
    : 'LIKE'
    ;

// vendor/HiveLexerParent.g:31
KW_ANY
    : 'ANY'
    ;

// vendor/HiveLexerParent.g:32
KW_IF
    : 'IF'
    ;

// vendor/HiveLexerParent.g:33
KW_EXISTS
    : 'EXISTS'
    ;

// vendor/HiveLexerParent.g:35
KW_ASC
    : 'ASC'
    ;

// vendor/HiveLexerParent.g:36
KW_DESC
    : 'DESC'
    ;

// vendor/HiveLexerParent.g:37
KW_NULLS
    : 'NULLS'
    ;

// vendor/HiveLexerParent.g:38
KW_LAST
    : 'LAST'
    ;

// vendor/HiveLexerParent.g:39
KW_ORDER
    : 'ORDER'
    ;

// vendor/HiveLexerParent.g:40
KW_GROUP
    : 'GROUP'
    ;

// vendor/HiveLexerParent.g:41
KW_BY
    : 'BY'
    ;

// vendor/HiveLexerParent.g:42
KW_HAVING
    : 'HAVING'
    ;

// vendor/HiveLexerParent.g:43
KW_QUALIFY
    : 'QUALIFY'
    ;

// vendor/HiveLexerParent.g:44
KW_WHERE
    : 'WHERE'
    ;

// vendor/HiveLexerParent.g:45
KW_FROM
    : 'FROM'
    ;

// vendor/HiveLexerParent.g:46
KW_AS
    : 'AS'
    ;

// vendor/HiveLexerParent.g:47
KW_SELECT
    : 'SELECT'
    ;

// vendor/HiveLexerParent.g:48
KW_DISTINCT
    : 'DISTINCT'
    ;

// vendor/HiveLexerParent.g:49
KW_INSERT
    : 'INSERT'
    ;

// vendor/HiveLexerParent.g:50
KW_OVERWRITE
    : 'OVERWRITE'
    ;

// vendor/HiveLexerParent.g:51
KW_OUTER
    : 'OUTER'
    ;

// vendor/HiveLexerParent.g:52
KW_UNIQUEJOIN
    : 'UNIQUEJOIN'
    ;

// vendor/HiveLexerParent.g:53
KW_PRESERVE
    : 'PRESERVE'
    ;

// vendor/HiveLexerParent.g:54
KW_JOIN
    : 'JOIN'
    ;

// vendor/HiveLexerParent.g:55
KW_LEFT
    : 'LEFT'
    ;

// vendor/HiveLexerParent.g:56
KW_RIGHT
    : 'RIGHT'
    ;

// vendor/HiveLexerParent.g:57
KW_FULL
    : 'FULL'
    ;

// vendor/HiveLexerParent.g:58
KW_ON
    : 'ON'
    ;

// vendor/HiveLexerParent.g:59
KW_PARTITION
    : 'PARTITION'
    ;

// vendor/HiveLexerParent.g:60
KW_PARTITIONS
    : 'PARTITIONS'
    ;

// vendor/HiveLexerParent.g:61
KW_TABLE
    : 'TABLE'
    ;

// vendor/HiveLexerParent.g:62
KW_TABLES
    : 'TABLES'
    ;

// vendor/HiveLexerParent.g:63
KW_COLUMNS
    : 'COLUMNS'
    ;

// vendor/HiveLexerParent.g:64
KW_INDEX
    : 'INDEX'
    ;

// vendor/HiveLexerParent.g:65
KW_INDEXES
    : 'INDEXES'
    ;

// vendor/HiveLexerParent.g:66
KW_REBUILD
    : 'REBUILD'
    ;

// vendor/HiveLexerParent.g:67
KW_FUNCTIONS
    : 'FUNCTIONS'
    ;

// vendor/HiveLexerParent.g:68
KW_SHOW
    : 'SHOW'
    ;

// vendor/HiveLexerParent.g:69
KW_MSCK
    : 'MSCK'
    ;

// vendor/HiveLexerParent.g:70
KW_REPAIR
    : 'REPAIR'
    ;

// vendor/HiveLexerParent.g:71
KW_DIRECTORY
    : 'DIRECTORY'
    ;

// vendor/HiveLexerParent.g:72
KW_LOCAL
    : 'LOCAL'
    ;

// vendor/HiveLexerParent.g:73
KW_TRANSFORM
    : 'TRANSFORM'
    ;

// vendor/HiveLexerParent.g:74
KW_USING
    : 'USING'
    ;

// vendor/HiveLexerParent.g:75
KW_COMPACT_ID
    : 'COMPACTIONID'
    ;

// vendor/HiveLexerParent.g:76
KW_CLUSTER
    : 'CLUSTER'
    ;

// vendor/HiveLexerParent.g:77
KW_DISTRIBUTE
    : 'DISTRIBUTE'
    ;

// vendor/HiveLexerParent.g:78
KW_SORT
    : 'SORT'
    ;

// vendor/HiveLexerParent.g:79
KW_UNION
    : 'UNION'
    ;

// vendor/HiveLexerParent.g:80
KW_EXCEPT
    : 'EXCEPT'
    ;

// vendor/HiveLexerParent.g:81
KW_LOAD
    : 'LOAD'
    ;

// vendor/HiveLexerParent.g:82
KW_PREPARE
    : 'PREPARE'
    ;

// vendor/HiveLexerParent.g:83
KW_EXPORT
    : 'EXPORT'
    ;

// vendor/HiveLexerParent.g:84
KW_IMPORT
    : 'IMPORT'
    ;

// vendor/HiveLexerParent.g:85
KW_REPLICATION
    : 'REPLICATION'
    ;

// vendor/HiveLexerParent.g:86
KW_METADATA
    : 'METADATA'
    ;

// vendor/HiveLexerParent.g:87
KW_DATA
    : 'DATA'
    ;

// vendor/HiveLexerParent.g:88
KW_INPATH
    : 'INPATH'
    ;

// vendor/HiveLexerParent.g:89
KW_IS
    : 'IS'
    ;

// vendor/HiveLexerParent.g:90
KW_NULL
    : 'NULL'
    ;

// vendor/HiveLexerParent.g:91
KW_CREATE
    : 'CREATE'
    ;

// vendor/HiveLexerParent.g:92
KW_EXTERNAL
    : 'EXTERNAL'
    ;

// vendor/HiveLexerParent.g:93
KW_ALTER
    : 'ALTER'
    ;

// vendor/HiveLexerParent.g:94
KW_CHANGE
    : 'CHANGE'
    ;

// vendor/HiveLexerParent.g:95
KW_COLUMN
    : 'COLUMN'
    ;

// vendor/HiveLexerParent.g:96
KW_FIRST
    : 'FIRST'
    ;

// vendor/HiveLexerParent.g:97
KW_AFTER
    : 'AFTER'
    ;

// vendor/HiveLexerParent.g:98
KW_DESCRIBE
    : 'DESCRIBE'
    ;

// vendor/HiveLexerParent.g:99
KW_DROP
    : 'DROP'
    ;

// vendor/HiveLexerParent.g:100
KW_RENAME
    : 'RENAME'
    ;

// vendor/HiveLexerParent.g:101
KW_TO
    : 'TO'
    ;

// vendor/HiveLexerParent.g:102
KW_COMMENT
    : 'COMMENT'
    ;

// vendor/HiveLexerParent.g:103
KW_BOOLEAN
    : 'BOOLEAN'
    ;

// vendor/HiveLexerParent.g:104
KW_TINYINT
    : 'TINYINT'
    ;

// vendor/HiveLexerParent.g:105
KW_SMALLINT
    : 'SMALLINT'
    ;

// vendor/HiveLexerParent.g:106
KW_INT
    : 'INT'
    | 'INTEGER'
    ;

// vendor/HiveLexerParent.g:107
KW_BIGINT
    : 'BIGINT'
    ;

// vendor/HiveLexerParent.g:108
KW_FLOAT
    : 'FLOAT'
    ;

// vendor/HiveLexerParent.g:109
KW_REAL
    : 'REAL'
    ;

// vendor/HiveLexerParent.g:110
KW_DOUBLE
    : 'DOUBLE'
    ;

// vendor/HiveLexerParent.g:111
KW_PRECISION
    : 'PRECISION'
    ;

// vendor/HiveLexerParent.g:112
KW_DATE
    : 'DATE'
    ;

// vendor/HiveLexerParent.g:113
KW_DATETIME
    : 'DATETIME'
    ;

// vendor/HiveLexerParent.g:114
KW_TIMESTAMP
    : 'TIMESTAMP'
    ;

// vendor/HiveLexerParent.g:115
KW_TIMESTAMPLOCALTZ
    : 'TIMESTAMPLOCALTZ'
    ;

// vendor/HiveLexerParent.g:116
KW_TIME
    : 'TIME'
    ;

// vendor/HiveLexerParent.g:117
KW_ZONE
    : 'ZONE'
    ;

// vendor/HiveLexerParent.g:118
KW_INTERVAL
    : 'INTERVAL'
    ;

// vendor/HiveLexerParent.g:119
KW_DECIMAL
    : 'DECIMAL'
    | 'DEC'
    | 'NUMERIC'
    ;

// vendor/HiveLexerParent.g:120
KW_STRING
    : 'STRING'
    ;

// vendor/HiveLexerParent.g:121
KW_CHAR
    : 'CHAR'
    ;

// vendor/HiveLexerParent.g:122
KW_VARCHAR
    : 'VARCHAR'
    ;

// vendor/HiveLexerParent.g:123
KW_ARRAY
    : 'ARRAY'
    ;

// vendor/HiveLexerParent.g:124
KW_STRUCT
    : 'STRUCT'
    ;

// vendor/HiveLexerParent.g:125
KW_MAP
    : 'MAP'
    ;

// vendor/HiveLexerParent.g:126
KW_UNIONTYPE
    : 'UNIONTYPE'
    ;

// vendor/HiveLexerParent.g:127
KW_REDUCE
    : 'REDUCE'
    ;

// vendor/HiveLexerParent.g:128
KW_PARTITIONED
    : 'PARTITIONED'
    ;

// vendor/HiveLexerParent.g:129
KW_CLUSTERED
    : 'CLUSTERED'
    ;

// vendor/HiveLexerParent.g:130
KW_DISTRIBUTED
    : 'DISTRIBUTED'
    ;

// vendor/HiveLexerParent.g:131
KW_SORTED
    : 'SORTED'
    ;

// vendor/HiveLexerParent.g:132
KW_INTO
    : 'INTO'
    ;

// vendor/HiveLexerParent.g:133
KW_BUCKETS
    : 'BUCKETS'
    ;

// vendor/HiveLexerParent.g:134
KW_ROW
    : 'ROW'
    ;

// vendor/HiveLexerParent.g:135
KW_ROWS
    : 'ROWS'
    ;

// vendor/HiveLexerParent.g:136
KW_FORMAT
    : 'FORMAT'
    ;

// vendor/HiveLexerParent.g:137
KW_DELIMITED
    : 'DELIMITED'
    ;

// vendor/HiveLexerParent.g:138
KW_FIELDS
    : 'FIELDS'
    ;

// vendor/HiveLexerParent.g:139
KW_TERMINATED
    : 'TERMINATED'
    ;

// vendor/HiveLexerParent.g:140
KW_ESCAPED
    : 'ESCAPED'
    ;

// vendor/HiveLexerParent.g:141
KW_COLLECTION
    : 'COLLECTION'
    ;

// vendor/HiveLexerParent.g:142
KW_ITEMS
    : 'ITEMS'
    ;

// vendor/HiveLexerParent.g:143
KW_KEYS
    : 'KEYS'
    ;

// vendor/HiveLexerParent.g:144
KW_KEY_TYPE
    : '$KEY$'
    ;

// vendor/HiveLexerParent.g:145
KW_KILL
    : 'KILL'
    ;

// vendor/HiveLexerParent.g:146
KW_LINES
    : 'LINES'
    ;

// vendor/HiveLexerParent.g:147
KW_STORED
    : 'STORED'
    ;

// vendor/HiveLexerParent.g:148
KW_FILEFORMAT
    : 'FILEFORMAT'
    ;

// vendor/HiveLexerParent.g:149
KW_INPUTFORMAT
    : 'INPUTFORMAT'
    ;

// vendor/HiveLexerParent.g:150
KW_OUTPUTFORMAT
    : 'OUTPUTFORMAT'
    ;

// vendor/HiveLexerParent.g:151
KW_INPUTDRIVER
    : 'INPUTDRIVER'
    ;

// vendor/HiveLexerParent.g:152
KW_OUTPUTDRIVER
    : 'OUTPUTDRIVER'
    ;

// vendor/HiveLexerParent.g:153
KW_ENABLE
    : 'ENABLE'
    | 'ENABLED'
    ;

// vendor/HiveLexerParent.g:154
KW_DISABLE
    : 'DISABLE'
    | 'DISABLED'
    ;

// vendor/HiveLexerParent.g:155
KW_EXECUTED
    : 'EXECUTED'
    ;

// vendor/HiveLexerParent.g:156
KW_EXECUTE
    : 'EXECUTE'
    ;

// vendor/HiveLexerParent.g:157
KW_LOCATION
    : 'LOCATION'
    ;

// vendor/HiveLexerParent.g:158
KW_MANAGED
    : 'MANAGED'
    ;

// vendor/HiveLexerParent.g:159
KW_MANAGEDLOCATION
    : 'MANAGEDLOCATION'
    ;

// vendor/HiveLexerParent.g:160
KW_TABLESAMPLE
    : 'TABLESAMPLE'
    ;

// vendor/HiveLexerParent.g:161
KW_BUCKET
    : 'BUCKET'
    ;

// vendor/HiveLexerParent.g:162
KW_OUT
    : 'OUT'
    ;

// vendor/HiveLexerParent.g:163
KW_OF
    : 'OF'
    ;

// vendor/HiveLexerParent.g:164
KW_PERCENT
    : 'PERCENT'
    ;

// vendor/HiveLexerParent.g:165
KW_CAST
    : 'CAST'
    ;

// vendor/HiveLexerParent.g:166
KW_ADD
    : 'ADD'
    ;

// vendor/HiveLexerParent.g:167
KW_REPLACE
    : 'REPLACE'
    ;

// vendor/HiveLexerParent.g:168
KW_RLIKE
    : 'RLIKE'
    ;

// vendor/HiveLexerParent.g:169
KW_REGEXP
    : 'REGEXP'
    ;

// vendor/HiveLexerParent.g:170
KW_TEMPORARY
    : 'TEMPORARY'
    ;

// vendor/HiveLexerParent.g:171
KW_FUNCTION
    : 'FUNCTION'
    ;

// vendor/HiveLexerParent.g:172
KW_MACRO
    : 'MACRO'
    ;

// vendor/HiveLexerParent.g:173
KW_FILE
    : 'FILE'
    ;

// vendor/HiveLexerParent.g:174
KW_JAR
    : 'JAR'
    ;

// vendor/HiveLexerParent.g:175
KW_EXPLAIN
    : 'EXPLAIN'
    ;

// vendor/HiveLexerParent.g:176
KW_DDL
    : 'DDL'
    ;

// vendor/HiveLexerParent.g:177
KW_EXTENDED
    : 'EXTENDED'
    ;

// vendor/HiveLexerParent.g:178
KW_DEBUG
    : 'DEBUG'
    ;

// vendor/HiveLexerParent.g:179
KW_FORMATTED
    : 'FORMATTED'
    ;

// vendor/HiveLexerParent.g:180
KW_DEPENDENCY
    : 'DEPENDENCY'
    ;

// vendor/HiveLexerParent.g:181
KW_LOGICAL
    : 'LOGICAL'
    ;

// vendor/HiveLexerParent.g:182
KW_CBO
    : 'CBO'
    ;

// vendor/HiveLexerParent.g:183
KW_SERDE
    : 'SERDE'
    ;

// vendor/HiveLexerParent.g:184
KW_WITH
    : 'WITH'
    ;

// vendor/HiveLexerParent.g:185
KW_DEFERRED
    : 'DEFERRED'
    ;

// vendor/HiveLexerParent.g:186
KW_SERDEPROPERTIES
    : 'SERDEPROPERTIES'
    ;

// vendor/HiveLexerParent.g:187
KW_DBPROPERTIES
    : 'DBPROPERTIES'
    ;

// vendor/HiveLexerParent.g:188
KW_DCPROPERTIES
    : 'DCPROPERTIES'
    ;

// vendor/HiveLexerParent.g:189
KW_LIMIT
    : 'LIMIT'
    ;

// vendor/HiveLexerParent.g:190
KW_OFFSET
    : 'OFFSET'
    ;

// vendor/HiveLexerParent.g:191
KW_SET
    : 'SET'
    ;

// vendor/HiveLexerParent.g:192
KW_UNSET
    : 'UNSET'
    ;

// vendor/HiveLexerParent.g:193
KW_TBLPROPERTIES
    : 'TBLPROPERTIES'
    ;

// vendor/HiveLexerParent.g:194
KW_IDXPROPERTIES
    : 'IDXPROPERTIES'
    ;

// vendor/HiveLexerParent.g:195
KW_VALUE_TYPE
    : '$VALUE$'
    ;

// vendor/HiveLexerParent.g:196
KW_ELEM_TYPE
    : '$ELEM$'
    ;

// vendor/HiveLexerParent.g:197
KW_DEFINED
    : 'DEFINED'
    ;

// vendor/HiveLexerParent.g:198
KW_CASE
    : 'CASE'
    ;

// vendor/HiveLexerParent.g:199
KW_WHEN
    : 'WHEN'
    ;

// vendor/HiveLexerParent.g:200
KW_THEN
    : 'THEN'
    ;

// vendor/HiveLexerParent.g:201
KW_ELSE
    : 'ELSE'
    ;

// vendor/HiveLexerParent.g:202
KW_END
    : 'END'
    ;

// vendor/HiveLexerParent.g:203
KW_MAPJOIN
    : 'MAPJOIN'
    ;

// vendor/HiveLexerParent.g:204
KW_STREAMTABLE
    : 'STREAMTABLE'
    ;

// vendor/HiveLexerParent.g:205
KW_CLUSTERSTATUS
    : 'CLUSTERSTATUS'
    ;

// vendor/HiveLexerParent.g:206
KW_UTC
    : 'UTC'
    ;

// vendor/HiveLexerParent.g:207
KW_UTCTIMESTAMP
    : 'UTC_TMESTAMP'
    ;

// vendor/HiveLexerParent.g:208
KW_LONG
    : 'LONG'
    ;

// vendor/HiveLexerParent.g:209
KW_DELETE
    : 'DELETE'
    ;

// vendor/HiveLexerParent.g:210
KW_PLUS
    : 'PLUS'
    ;

// vendor/HiveLexerParent.g:211
KW_MINUS
    : 'MINUS'
    ;

// vendor/HiveLexerParent.g:212
KW_FETCH
    : 'FETCH'
    ;

// vendor/HiveLexerParent.g:213
KW_INTERSECT
    : 'INTERSECT'
    ;

// vendor/HiveLexerParent.g:214
KW_VIEW
    : 'VIEW'
    ;

// vendor/HiveLexerParent.g:215
KW_VIEWS
    : 'VIEWS'
    ;

// vendor/HiveLexerParent.g:216
KW_IN
    : 'IN'
    ;

// vendor/HiveLexerParent.g:217
KW_DATABASE
    : 'DATABASE'
    ;

// vendor/HiveLexerParent.g:218
KW_DATABASES
    : 'DATABASES'
    ;

// vendor/HiveLexerParent.g:219
KW_MATERIALIZED
    : 'MATERIALIZED'
    ;

// vendor/HiveLexerParent.g:220
KW_SCHEMA
    : 'SCHEMA'
    ;

// vendor/HiveLexerParent.g:221
KW_SCHEMAS
    : 'SCHEMAS'
    ;

// vendor/HiveLexerParent.g:222
KW_GRANT
    : 'GRANT'
    ;

// vendor/HiveLexerParent.g:223
KW_REVOKE
    : 'REVOKE'
    ;

// vendor/HiveLexerParent.g:224
KW_SSL
    : 'SSL'
    ;

// vendor/HiveLexerParent.g:225
KW_UNDO
    : 'UNDO'
    ;

// vendor/HiveLexerParent.g:226
KW_LOCK
    : 'LOCK'
    ;

// vendor/HiveLexerParent.g:227
KW_LOCKS
    : 'LOCKS'
    ;

// vendor/HiveLexerParent.g:228
KW_UNLOCK
    : 'UNLOCK'
    ;

// vendor/HiveLexerParent.g:229
KW_SHARED
    : 'SHARED'
    ;

// vendor/HiveLexerParent.g:230
KW_EXCLUSIVE
    : 'EXCLUSIVE'
    ;

// vendor/HiveLexerParent.g:231
KW_PROCEDURE
    : 'PROCEDURE'
    ;

// vendor/HiveLexerParent.g:232
KW_UNSIGNED
    : 'UNSIGNED'
    ;

// vendor/HiveLexerParent.g:233
KW_WHILE
    : 'WHILE'
    ;

// vendor/HiveLexerParent.g:234
KW_READ
    : 'READ'
    ;

// vendor/HiveLexerParent.g:235
KW_READS
    : 'READS'
    ;

// vendor/HiveLexerParent.g:236
KW_PURGE
    : 'PURGE'
    ;

// vendor/HiveLexerParent.g:237
KW_RANGE
    : 'RANGE'
    ;

// vendor/HiveLexerParent.g:238
KW_ANALYZE
    : 'ANALYZE'
    ;

// vendor/HiveLexerParent.g:239
KW_BEFORE
    : 'BEFORE'
    ;

// vendor/HiveLexerParent.g:240
KW_BETWEEN
    : 'BETWEEN'
    ;

// vendor/HiveLexerParent.g:241
KW_BOTH
    : 'BOTH'
    ;

// vendor/HiveLexerParent.g:242
KW_BINARY
    : 'BINARY'
    ;

// vendor/HiveLexerParent.g:243
KW_CROSS
    : 'CROSS'
    ;

// vendor/HiveLexerParent.g:244
KW_CONTINUE
    : 'CONTINUE'
    ;

// vendor/HiveLexerParent.g:245
KW_CONVERT
    : 'CONVERT'
    ;

// vendor/HiveLexerParent.g:246
KW_CURSOR
    : 'CURSOR'
    ;

// vendor/HiveLexerParent.g:247
KW_TRIGGER
    : 'TRIGGER'
    ;

// vendor/HiveLexerParent.g:248
KW_RECORDREADER
    : 'RECORDREADER'
    ;

// vendor/HiveLexerParent.g:249
KW_RECORDWRITER
    : 'RECORDWRITER'
    ;

// vendor/HiveLexerParent.g:250
KW_SEMI
    : 'SEMI'
    ;

// vendor/HiveLexerParent.g:251
KW_LATERAL
    : 'LATERAL'
    ;

// vendor/HiveLexerParent.g:252
KW_TOUCH
    : 'TOUCH'
    ;

// vendor/HiveLexerParent.g:253
KW_ARCHIVE
    : 'ARCHIVE'
    ;

// vendor/HiveLexerParent.g:254
KW_UNARCHIVE
    : 'UNARCHIVE'
    ;

// vendor/HiveLexerParent.g:255
KW_COMPUTE
    : 'COMPUTE'
    ;

// vendor/HiveLexerParent.g:256
KW_STATISTICS
    : 'STATISTICS'
    ;

// vendor/HiveLexerParent.g:257
KW_USE
    : 'USE'
    ;

// vendor/HiveLexerParent.g:258
KW_OPTION
    : 'OPTION'
    ;

// vendor/HiveLexerParent.g:259
KW_CONCATENATE
    : 'CONCATENATE'
    ;

// vendor/HiveLexerParent.g:260
KW_SHOW_DATABASE
    : 'SHOW_DATABASE'
    ;

// vendor/HiveLexerParent.g:261
KW_UPDATE
    : 'UPDATE'
    ;

// vendor/HiveLexerParent.g:262
KW_RESTRICT
    : 'RESTRICT'
    ;

// vendor/HiveLexerParent.g:263
KW_CASCADE
    : 'CASCADE'
    ;

// vendor/HiveLexerParent.g:264
KW_SKEWED
    : 'SKEWED'
    ;

// vendor/HiveLexerParent.g:265
KW_ROLLUP
    : 'ROLLUP'
    ;

// vendor/HiveLexerParent.g:266
KW_CUBE
    : 'CUBE'
    ;

// vendor/HiveLexerParent.g:267
KW_DIRECTORIES
    : 'DIRECTORIES'
    ;

// vendor/HiveLexerParent.g:268
KW_FOR
    : 'FOR'
    ;

// vendor/HiveLexerParent.g:269
KW_WINDOW
    : 'WINDOW'
    ;

// vendor/HiveLexerParent.g:270
KW_UNBOUNDED
    : 'UNBOUNDED'
    ;

// vendor/HiveLexerParent.g:271
KW_PRECEDING
    : 'PRECEDING'
    ;

// vendor/HiveLexerParent.g:272
KW_FOLLOWING
    : 'FOLLOWING'
    ;

// vendor/HiveLexerParent.g:273
KW_CURRENT
    : 'CURRENT'
    ;

// vendor/HiveLexerParent.g:274
KW_CURRENT_DATE
    : 'CURRENT_DATE'
    ;

// vendor/HiveLexerParent.g:275
KW_CURRENT_TIMESTAMP
    : 'CURRENT_TIMESTAMP'
    ;

// vendor/HiveLexerParent.g:276
KW_LESS
    : 'LESS'
    ;

// vendor/HiveLexerParent.g:277
KW_MORE
    : 'MORE'
    ;

// vendor/HiveLexerParent.g:278
KW_OVER
    : 'OVER'
    ;

// vendor/HiveLexerParent.g:279
KW_GROUPING
    : 'GROUPING'
    ;

// vendor/HiveLexerParent.g:280
KW_SETS
    : 'SETS'
    ;

// vendor/HiveLexerParent.g:281
KW_TRUNCATE
    : 'TRUNCATE'
    ;

// vendor/HiveLexerParent.g:282
KW_NOSCAN
    : 'NOSCAN'
    ;

// vendor/HiveLexerParent.g:283
KW_USER
    : 'USER'
    ;

// vendor/HiveLexerParent.g:284
KW_ROLE
    : 'ROLE'
    ;

// vendor/HiveLexerParent.g:285
KW_ROLES
    : 'ROLES'
    ;

// vendor/HiveLexerParent.g:286
KW_INNER
    : 'INNER'
    ;

// vendor/HiveLexerParent.g:287
KW_EXCHANGE
    : 'EXCHANGE'
    ;

// vendor/HiveLexerParent.g:288
KW_URI
    : 'URI'
    ;

// vendor/HiveLexerParent.g:289
KW_SERVER
    : 'SERVER'
    ;

// vendor/HiveLexerParent.g:290
KW_ADMIN
    : 'ADMIN'
    ;

// vendor/HiveLexerParent.g:291
KW_OWNER
    : 'OWNER'
    ;

// vendor/HiveLexerParent.g:292
KW_PRINCIPALS
    : 'PRINCIPALS'
    ;

// vendor/HiveLexerParent.g:293
KW_COMPACT
    : 'COMPACT'
    ;

// vendor/HiveLexerParent.g:294
KW_COMPACTIONS
    : 'COMPACTIONS'
    ;

// vendor/HiveLexerParent.g:295
KW_TRANSACTIONS
    : 'TRANSACTIONS'
    ;

// vendor/HiveLexerParent.g:296
KW_TRANSACTIONAL
    : 'TRANSACTIONAL'
    ;

// vendor/HiveLexerParent.g:297
KW_REWRITE
    : 'REWRITE'
    ;

// vendor/HiveLexerParent.g:298
KW_AUTHORIZATION
    : 'AUTHORIZATION'
    ;

// vendor/HiveLexerParent.g:299
KW_REOPTIMIZATION
    : 'REOPTIMIZATION'
    ;

// vendor/HiveLexerParent.g:300
KW_CONF
    : 'CONF'
    ;

// vendor/HiveLexerParent.g:301
KW_VALUES
    : 'VALUES'
    ;

// vendor/HiveLexerParent.g:302
KW_RELOAD
    : 'RELOAD'
    ;

// vendor/HiveLexerParent.g:303
KW_YEAR
    : 'YEAR'
    | 'YEARS'
    ;

// vendor/HiveLexerParent.g:304
KW_QUERY
    : 'QUERY'
    ;

// vendor/HiveLexerParent.g:305
KW_QUARTER
    : 'QUARTER'
    ;

// vendor/HiveLexerParent.g:306
KW_MONTH
    : 'MONTH'
    | 'MONTHS'
    ;

// vendor/HiveLexerParent.g:307
KW_WEEK
    : 'WEEK'
    | 'WEEKS'
    ;

// vendor/HiveLexerParent.g:308
KW_DAY
    : 'DAY'
    | 'DAYS'
    ;

// vendor/HiveLexerParent.g:309
KW_DOW
    : 'DAYOFWEEK'
    ;

// vendor/HiveLexerParent.g:310
KW_HOUR
    : 'HOUR'
    | 'HOURS'
    ;

// vendor/HiveLexerParent.g:311
KW_MINUTE
    : 'MINUTE'
    | 'MINUTES'
    ;

// vendor/HiveLexerParent.g:312
KW_SECOND
    : 'SECOND'
    | 'SECONDS'
    ;

// vendor/HiveLexerParent.g:313
KW_START
    : 'START'
    ;

// vendor/HiveLexerParent.g:314
KW_TRANSACTION
    : 'TRANSACTION'
    ;

// vendor/HiveLexerParent.g:315
KW_COMMIT
    : 'COMMIT'
    ;

// vendor/HiveLexerParent.g:316
KW_ROLLBACK
    : 'ROLLBACK'
    ;

// vendor/HiveLexerParent.g:317
KW_WORK
    : 'WORK'
    ;

// vendor/HiveLexerParent.g:318
KW_ONLY
    : 'ONLY'
    ;

// vendor/HiveLexerParent.g:319
KW_WRITE
    : 'WRITE'
    ;

// vendor/HiveLexerParent.g:320
KW_ISOLATION
    : 'ISOLATION'
    ;

// vendor/HiveLexerParent.g:321
KW_LEVEL
    : 'LEVEL'
    ;

// vendor/HiveLexerParent.g:322
KW_SNAPSHOT
    : 'SNAPSHOT'
    ;

// vendor/HiveLexerParent.g:323
KW_AUTOCOMMIT
    : 'AUTOCOMMIT'
    ;

// vendor/HiveLexerParent.g:324
KW_CACHE
    : 'CACHE'
    ;

// vendor/HiveLexerParent.g:325
KW_PRIMARY
    : 'PRIMARY'
    ;

// vendor/HiveLexerParent.g:326
KW_FOREIGN
    : 'FOREIGN'
    ;

// vendor/HiveLexerParent.g:327
KW_REFERENCES
    : 'REFERENCES'
    ;

// vendor/HiveLexerParent.g:328
KW_CONSTRAINT
    : 'CONSTRAINT'
    ;

// vendor/HiveLexerParent.g:329
KW_FORCE
    : 'FORCE'
    ;

// vendor/HiveLexerParent.g:330
KW_ENFORCED
    : 'ENFORCED'
    ;

// vendor/HiveLexerParent.g:331
KW_VALIDATE
    : 'VALIDATE'
    ;

// vendor/HiveLexerParent.g:332
KW_NOVALIDATE
    : 'NOVALIDATE'
    ;

// vendor/HiveLexerParent.g:333
KW_RELY
    : 'RELY'
    ;

// vendor/HiveLexerParent.g:334
KW_NORELY
    : 'NORELY'
    ;

// vendor/HiveLexerParent.g:335
KW_UNIQUE
    : 'UNIQUE'
    ;

// vendor/HiveLexerParent.g:336
KW_KEY
    : 'KEY'
    ;

// vendor/HiveLexerParent.g:337
KW_ABORT
    : 'ABORT'
    ;

// vendor/HiveLexerParent.g:338
KW_EXTRACT
    : 'EXTRACT'
    ;

// vendor/HiveLexerParent.g:339
KW_FLOOR
    : 'FLOOR'
    ;

// vendor/HiveLexerParent.g:340
KW_MERGE
    : 'MERGE'
    ;

// vendor/HiveLexerParent.g:341
KW_MATCHED
    : 'MATCHED'
    ;

// vendor/HiveLexerParent.g:342
KW_REPL
    : 'REPL'
    ;

// vendor/HiveLexerParent.g:343
KW_DUMP
    : 'DUMP'
    ;

// vendor/HiveLexerParent.g:344
KW_STATUS
    : 'STATUS'
    ;

// vendor/HiveLexerParent.g:345
KW_VECTORIZATION
    : 'VECTORIZATION'
    ;

// vendor/HiveLexerParent.g:346
KW_SUMMARY
    : 'SUMMARY'
    ;

// vendor/HiveLexerParent.g:347
KW_OPERATOR
    : 'OPERATOR'
    ;

// vendor/HiveLexerParent.g:348
KW_EXPRESSION
    : 'EXPRESSION'
    ;

// vendor/HiveLexerParent.g:349
KW_DETAIL
    : 'DETAIL'
    ;

// vendor/HiveLexerParent.g:350
KW_WAIT
    : 'WAIT'
    ;

// vendor/HiveLexerParent.g:351
KW_RESOURCE
    : 'RESOURCE'
    ;

// vendor/HiveLexerParent.g:352
KW_PLAN
    : 'PLAN'
    ;

// vendor/HiveLexerParent.g:353
KW_QUERY_PARALLELISM
    : 'QUERY_PARALLELISM'
    ;

// vendor/HiveLexerParent.g:354
KW_PLANS
    : 'PLANS'
    ;

// vendor/HiveLexerParent.g:355
KW_ACTIVATE
    : 'ACTIVATE'
    ;

// vendor/HiveLexerParent.g:356
KW_DEFAULT
    : 'DEFAULT'
    ;

// vendor/HiveLexerParent.g:357
KW_CHECK
    : 'CHECK'
    ;

// vendor/HiveLexerParent.g:358
KW_POOL
    : 'POOL'
    ;

// vendor/HiveLexerParent.g:359
KW_MOVE
    : 'MOVE'
    ;

// vendor/HiveLexerParent.g:360
KW_DO
    : 'DO'
    ;

// vendor/HiveLexerParent.g:361
KW_ALLOC_FRACTION
    : 'ALLOC_FRACTION'
    ;

// vendor/HiveLexerParent.g:362
KW_SCHEDULING_POLICY
    : 'SCHEDULING_POLICY'
    ;

// vendor/HiveLexerParent.g:363
KW_SCHEDULED
    : 'SCHEDULED'
    ;

// vendor/HiveLexerParent.g:364
KW_EVERY
    : 'EVERY'
    ;

// vendor/HiveLexerParent.g:365
KW_AT
    : 'AT'
    ;

// vendor/HiveLexerParent.g:366
KW_CRON
    : 'CRON'
    ;

// vendor/HiveLexerParent.g:367
KW_PATH
    : 'PATH'
    ;

// vendor/HiveLexerParent.g:368
KW_MAPPING
    : 'MAPPING'
    ;

// vendor/HiveLexerParent.g:369
KW_WORKLOAD
    : 'WORKLOAD'
    ;

// vendor/HiveLexerParent.g:370
KW_MANAGEMENT
    : 'MANAGEMENT'
    ;

// vendor/HiveLexerParent.g:371
KW_ACTIVE
    : 'ACTIVE'
    ;

// vendor/HiveLexerParent.g:372
KW_UNMANAGED
    : 'UNMANAGED'
    ;

// vendor/HiveLexerParent.g:373
KW_APPLICATION
    : 'APPLICATION'
    ;

// vendor/HiveLexerParent.g:374
KW_SYNC
    : 'SYNC'
    ;

// vendor/HiveLexerParent.g:375
KW_AST
    : 'AST'
    ;

// vendor/HiveLexerParent.g:376
KW_COST
    : 'COST'
    ;

// vendor/HiveLexerParent.g:377
KW_JOINCOST
    : 'JOINCOST'
    ;

// vendor/HiveLexerParent.g:378
KW_WITHIN
    : 'WITHIN'
    ;

// vendor/HiveLexerParent.g:379
KW_PKFK_JOIN
    : 'PKFK_JOIN'
    ;

// vendor/HiveLexerParent.g:380
KW_LEADING
    : 'LEADING'
    ;

// vendor/HiveLexerParent.g:381
KW_TRAILING
    : 'TRAILING'
    ;

// vendor/HiveLexerParent.g:382
KW_TRIM
    : 'TRIM'
    ;

// vendor/HiveLexerParent.g:383
KW_IGNORE
    : 'IGNORE'
    ;

// vendor/HiveLexerParent.g:384
KW_RESPECT
    : 'RESPECT'
    ;

// vendor/HiveLexerParent.g:385
KW_DATACONNECTOR
    : 'CONNECTOR'
    ;

// vendor/HiveLexerParent.g:386
KW_DATACONNECTORS
    : 'CONNECTORS'
    ;

// vendor/HiveLexerParent.g:387
KW_TYPE
    : 'TYPE'
    ;

// vendor/HiveLexerParent.g:388
KW_URL
    : 'URL'
    ;

// vendor/HiveLexerParent.g:389
KW_REMOTE
    : 'REMOTE'
    ;

// vendor/HiveLexerParent.g:390
KW_SPEC
    : 'SPEC'
    ;

// vendor/HiveLexerParent.g:391
KW_SYSTEM_TIME
    : 'SYSTEM_TIME'
    ;

// vendor/HiveLexerParent.g:392
KW_SYSTEM_VERSION
    : 'SYSTEM_VERSION'
    ;

// vendor/HiveLexerParent.g:393
KW_EXPIRE_SNAPSHOTS
    : 'EXPIRE_SNAPSHOTS'
    ;

// vendor/HiveLexerParent.g:394
KW_SET_CURRENT_SNAPSHOT
    : 'SET_CURRENT_SNAPSHOT'
    ;

// vendor/HiveLexerParent.g:395
KW_BRANCH
    : 'BRANCH'
    ;

// vendor/HiveLexerParent.g:396
KW_SNAPSHOTS
    : 'SNAPSHOTS'
    ;

// vendor/HiveLexerParent.g:397
KW_RETAIN
    : 'RETAIN'
    ;

// vendor/HiveLexerParent.g:398
KW_RETENTION
    : 'RETENTION'
    ;

// vendor/HiveLexerParent.g:399
KW_TAG
    : 'TAG'
    ;

// vendor/HiveLexerParent.g:400
KW_FAST_FORWARD
    : 'FAST-FORWARD'
    ;

// vendor/HiveLexerParent.g:401
KW_CHERRY_PICK
    : 'CHERRY-PICK'
    ;

// vendor/HiveLexerParent.g:402
KW_ORPHAN_FILES
    : 'ORPHAN-FILES'
    ;

// vendor/HiveLexerParent.g:403
KW_OLDER
    : 'OLDER'
    ;

// vendor/HiveLexerParent.g:404
KW_THAN
    : 'THAN'
    ;

// vendor/HiveLexerParent.g:405
KW_OPTIMIZE
    : 'OPTIMIZE'
    ;

// vendor/HiveLexerParent.g:410
DOT
    : '.'
    ;

// vendor/HiveLexerParent.g:411
COLON
    : ':'
    ;

// vendor/HiveLexerParent.g:412
COMMA
    : ','
    ;

// vendor/HiveLexerParent.g:413
SEMICOLON
    : ';'
    ;

// vendor/HiveLexerParent.g:415
LPAREN
    : '('
    ;

// vendor/HiveLexerParent.g:416
RPAREN
    : ')'
    ;

// vendor/HiveLexerParent.g:417
LSQUARE
    : '['
    ;

// vendor/HiveLexerParent.g:418
RSQUARE
    : ']'
    ;

// vendor/HiveLexerParent.g:419
LCURLY
    : '{'
    ;

// vendor/HiveLexerParent.g:420
RCURLY
    : '}'
    ;

// vendor/HiveLexerParent.g:422
EQUAL
    : '='
    | '=='
    ;

// vendor/HiveLexerParent.g:423
EQUAL_NS
    : '<=>'
    ;

// vendor/HiveLexerParent.g:424
NOTEQUAL
    : '<>'
    | '!='
    ;

// vendor/HiveLexerParent.g:425
LESSTHANOREQUALTO
    : '<='
    ;

// vendor/HiveLexerParent.g:426
LESSTHAN
    : '<'
    ;

// vendor/HiveLexerParent.g:427
GREATERTHANOREQUALTO
    : '>='
    ;

// vendor/HiveLexerParent.g:428
GREATERTHAN
    : '>'
    ;

// vendor/HiveLexerParent.g:430
DIVIDE
    : '/'
    ;

// vendor/HiveLexerParent.g:431
PLUS
    : '+'
    ;

// vendor/HiveLexerParent.g:432
MINUS
    : '-'
    ;

// vendor/HiveLexerParent.g:433
STAR
    : '*'
    ;

// vendor/HiveLexerParent.g:434
MOD
    : '%'
    ;

// vendor/HiveLexerParent.g:435
DIV
    : 'DIV'
    ;

// vendor/HiveLexerParent.g:437
AMPERSAND
    : '&'
    ;

// vendor/HiveLexerParent.g:438
TILDE
    : '~'
    ;

// vendor/HiveLexerParent.g:439
BITWISEOR
    : '|'
    ;

// vendor/HiveLexerParent.g:440
CONCATENATE
    : '||'
    ;

// vendor/HiveLexerParent.g:441
BITWISEXOR
    : '^'
    ;

// vendor/HiveLexerParent.g:442
QUESTION
    : '?'
    ;

// vendor/HiveLexerParent.g:443
DOLLAR
    : '$'
    ;

// vendor/HiveLexerParent.g:447
fragment Letter
    : 'a'..'z'
    | 'A'..'Z'
    ;

// vendor/HiveLexerParent.g:452
fragment HexDigit
    : 'a'..'f'
    | 'A'..'F'
    ;

// vendor/HiveLexerParent.g:457
fragment Digit
    : '0'..'9'
    ;

// vendor/HiveLexerParent.g:463
fragment Exponent
    : ( 'e' | 'E' ) ( PLUS | MINUS )? ( Digit )+
    ;

// vendor/HiveLexerParent.g:469
fragment RegexComponent
    : 'a'..'z'
    | 'A'..'Z'
    | '0'..'9'
    | '_'
    | PLUS
    | STAR
    | QUESTION
    | MINUS
    | DOT
    | LPAREN
    | RPAREN
    | LSQUARE
    | RSQUARE
    | LCURLY
    | RCURLY
    | BITWISEXOR
    | BITWISEOR
    | DOLLAR
    | '!'
    ;

// vendor/HiveLexerParent.g:481
CharSetLiteral
    : StringLiteral
    | '0' 'X' ( HexDigit | Digit )+
    ;

// vendor/HiveLexerParent.g:487
IntegralLiteral
    : ( Digit )+ ( 'L' | 'S' | 'Y' )
    ;

// vendor/HiveLexerParent.g:492
NumberLiteral
    : Number ( 'D' | 'B' 'D' )
    ;

// vendor/HiveLexerParent.g:497
ByteLengthLiteral
    : ( Digit )+ ( 'b' | 'B' | 'k' | 'K' | 'm' | 'M' | 'g' | 'G' )
    ;

// vendor/HiveLexerParent.g:502
Number
    : ( Digit )+ ( DOT ( Digit )* ( Exponent )? | Exponent )?
    ;

// vendor/HiveLexerParent.g:528
// vendor/HiveLexerParent.g:531: predicate dropped: {allowQuotedId() != Quotation.NONE}?
Identifier
    : ( Letter | Digit ) ( Letter | Digit | '_' )*
    | QuotedIdentifier
    | '`' RegexComponent+ '`'
    ;

// vendor/HiveLexerParent.g:541
CharSetName
    : '_' ( Letter | Digit | '_' | '-' | '.' | ':' )+
    ;

// vendor/HiveLexerParent.g:546
WS
    : ( ' ' | '\r' | '\t' | '\n' )
      -> channel(HIDDEN)
    ;

// vendor/HiveLexerParent.g:549
LINE_COMMENT
    : '--' ( ~( '\n' | '\r' ) )*
      -> channel(HIDDEN)
    ;

// vendor/HiveLexerParent.g:553
// vendor/HiveLexerParent.g:554: conditional hidden-channel action dropped: {if(getText().charAt(2) != '+') { $channel=HIDDEN; } else { setText(getText().substring(3, getText().length() - 2)); }}
QUERY_HINT
    : '/*' ( QUERY_HINT | . )* '*/'
    ;
