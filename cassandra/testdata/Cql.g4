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

// This file is derived from Apache Cassandra's src/antlr/Cql.g, Parser.g and
// Lexer.g, which are beside it as published, with the license and notice
// they come under (Cassandra.LICENSE, Cassandra.NOTICE). It was changed:
// cassandra/parser/cmd/antlr3to4 wrote it from them, converting ANTLR 3 to
// ANTLR 4 and keeping only what each rule derives. Actions, arguments,
// labels and predicates are gone; an alternative whose action refuses every
// input is gone with it; and what else was lost is noted beside its rule.
// Each rule is under the file and line it came from.
//
// Do not edit: rerun the converter.

grammar Cql;

// Cql.g:133
query
    : cqlStatement ( ';' )* EOF
    ;

// Parser.g:207
cqlStatement
    : selectStatement
    | insertStatement
    | updateStatement
    | batchStatement
    | deleteStatement
    | useStatement
    | truncateStatement
    | createKeyspaceStatement
    | createTableStatement
    | createIndexStatement
    | dropKeyspaceStatement
    | dropTableStatement
    | dropIndexStatement
    | alterTableStatement
    | alterKeyspaceStatement
    | grantPermissionsStatement
    | revokePermissionsStatement
    | listPermissionsStatement
    | createUserStatement
    | alterUserStatement
    | dropUserStatement
    | listUsersStatement
    | createTriggerStatement
    | dropTriggerStatement
    | createTypeStatement
    | alterTypeStatement
    | dropTypeStatement
    | createFunctionStatement
    | dropFunctionStatement
    | createAggregateStatement
    | dropAggregateStatement
    | createRoleStatement
    | alterRoleStatement
    | dropRoleStatement
    | listRolesStatement
    | grantRoleStatement
    | revokeRoleStatement
    | createMaterializedViewStatement
    | dropMaterializedViewStatement
    | alterMaterializedViewStatement
    | describeStatement
    | addIdentityStatement
    | dropIdentityStatement
    ;

// Parser.g:257
useStatement
    : K_USE keyspaceName
    ;

// Parser.g:267
selectStatement
    : K_SELECT ( K_JSON )? selectClause K_FROM columnFamilyName ( K_WHERE whereClause )? ( K_GROUP K_BY groupByClause ( ',' groupByClause )* )? ( K_ORDER K_BY orderByClause ( ',' orderByClause )* )? ( K_PER K_PARTITION K_LIMIT intValue )? ( K_LIMIT intValue )? ( K_ALLOW K_FILTERING )?
    ;

// Parser.g:297
selectClause
    : K_DISTINCT selectors
    | selectors
    ;

// Parser.g:304
selectors
    : selector ( ',' selector )*
    | '*'
    ;

// Parser.g:309
selector
    : unaliasedSelector ( K_AS noncol_ident )?
    ;

// Parser.g:314
unaliasedSelector
    : selectionAddition
    ;

// Parser.g:318
selectionAddition
    : selectionMultiplication ( '+' selectionMultiplication | '-' selectionMultiplication )*
    ;

// Parser.g:325
selectionMultiplication
    : selectionGroup ( '*' selectionGroup | '/' selectionGroup | '%' selectionGroup )*
    ;

// Parser.g:333
selectionGroup
    : selectionGroupWithField
    | selectionGroupWithoutField
    | '-' selectionGroup
    ;

// Parser.g:339
selectionGroupWithField
    : selectionGroupWithoutField selectorModifier
    ;

// Parser.g:343
selectorModifier
    : fieldSelectorModifier selectorModifier
    | '[' collectionSubSelection ']' selectorModifier
    | /* empty */
    ;

// Parser.g:349
fieldSelectorModifier
    : '.' fident
    ;

// Parser.g:353
collectionSubSelection
    : ( term ( RANGE ( term )? )? | RANGE term )
    ;

// Parser.g:364
selectionGroupWithoutField
    : simpleUnaliasedSelector
    | selectionTypeHint
    | selectionTupleOrNestedSelector
    | selectionList
    | selectionMapOrSet
    ;

// Parser.g:375
selectionTypeHint
    : '(' comparatorType ')' selectionGroupWithoutField
    ;

// Parser.g:379
selectionList
    : '[' ( unaliasedSelector ( ',' unaliasedSelector )* )? ']'
    ;

// Parser.g:385
selectionMapOrSet
    : '{' unaliasedSelector ( selectionMap | selectionSet ) '}'
    | '{' '}'
    ;

// Parser.g:390
selectionMap
    : ':' unaliasedSelector ( ',' unaliasedSelector ':' unaliasedSelector )*
    ;

// Parser.g:396
selectionSet
    : ( ',' unaliasedSelector )*
    ;

// Parser.g:402
selectionTupleOrNestedSelector
    : '(' unaliasedSelector ( ',' unaliasedSelector )* ')'
    ;

// Parser.g:412
simpleUnaliasedSelector
    : sident
    | selectionLiteral
    | selectionFunction
    ;

// Parser.g:418
selectionFunction
    : K_COUNT '(' '*' ')'
    | K_MAXWRITETIME '(' sident selectorModifier ')'
    | K_WRITETIME '(' sident selectorModifier ')'
    | K_TTL '(' sident selectorModifier ')'
    | K_CAST '(' unaliasedSelector K_AS native_type ')'
    | functionName selectionFunctionArgs
    ;

// Parser.g:427
selectionLiteral
    : constant
    | K_NULL
    | ':' noncol_ident
    | QMARK
    ;

// Parser.g:434
selectionFunctionArgs
    : '(' ( unaliasedSelector ( ',' unaliasedSelector )* )? ')'
    ;

// Parser.g:441
sident
    : IDENT
    | QUOTED_NAME
    | unreserved_keyword
    ;

// Parser.g:447
whereClause
    : relationOrExpression ( K_AND relationOrExpression )*
    ;

// Parser.g:452
relationOrExpression
    : relation
    | customIndexExpression
    ;

// Parser.g:457
customIndexExpression
    : 'expr(' idxName ',' term ')'
    ;

// Parser.g:462
orderByClause
    : cident ( K_ANN K_OF term )? ( K_ASC | K_DESC )?
    ;

// Parser.g:475
groupByClause
    : unaliasedSelector
    ;

// Parser.g:485
insertStatement
    : K_INSERT K_INTO columnFamilyName ( normalInsertStatement | K_JSON jsonInsertStatement )
    ;

// Parser.g:491
normalInsertStatement
    : '(' cident ( ',' cident )* ')' K_VALUES '(' term ( ',' term )* ')' ( K_IF K_NOT K_EXISTS )? ( usingClause )?
    ;

// Parser.g:508
jsonInsertStatement
    : jsonValue ( K_DEFAULT ( K_NULL | ( K_UNSET ) ) )? ( K_IF K_NOT K_EXISTS )? ( usingClause )?
    ;

// Parser.g:523
jsonValue
    : STRING_LITERAL
    | ':' noncol_ident
    | QMARK
    ;

// Parser.g:529
usingClause
    : K_USING usingClauseObjective ( K_AND usingClauseObjective )*
    ;

// Parser.g:533
usingClauseObjective
    : K_TIMESTAMP intValue
    | K_TTL intValue
    ;

// Parser.g:545
updateStatement
    : K_UPDATE columnFamilyName ( usingClause )? K_SET columnOperation ( ',' columnOperation )* K_WHERE whereClause ( K_IF ( K_EXISTS | updateConditions ) )?
    ;

// Parser.g:566
updateConditions
    : columnCondition ( K_AND columnCondition )*
    ;

// Parser.g:579
deleteStatement
    : K_DELETE ( deleteSelection )? K_FROM columnFamilyName ( usingClauseDelete )? K_WHERE whereClause ( K_IF ( K_EXISTS | updateConditions ) )?
    ;

// Parser.g:600
deleteSelection
    : deleteOp ( ',' deleteOp )*
    ;

// Parser.g:606
deleteOp
    : cident
    | cident '[' term ']'
    | cident '.' fident
    ;

// Parser.g:612
usingClauseDelete
    : K_USING K_TIMESTAMP intValue
    ;

// Parser.g:640
batchStatement
    : K_BEGIN ( K_UNLOGGED | K_COUNTER )? K_BATCH ( usingClause )? ( batchStatementObjective ';'? )* K_APPLY K_BATCH
    ;

// Parser.g:656
batchStatementObjective
    : insertStatement
    | updateStatement
    | deleteStatement
    ;

// Parser.g:662
createAggregateStatement
    : K_CREATE ( K_OR K_REPLACE )? K_AGGREGATE ( K_IF K_NOT K_EXISTS )? functionName '(' ( comparatorType ( ',' comparatorType )* )? ')' K_SFUNC allowedFunctionName K_STYPE comparatorType ( K_FINALFUNC allowedFunctionName )? ( K_INITCOND term )?
    ;

// Parser.g:690
dropAggregateStatement
    : K_DROP K_AGGREGATE ( K_IF K_EXISTS )? functionName ( '(' ( comparatorType ( ',' comparatorType )* )? ')' )?
    ;

// Parser.g:711
createFunctionStatement
    : K_CREATE ( K_OR K_REPLACE )? K_FUNCTION ( K_IF K_NOT K_EXISTS )? functionName '(' ( noncol_ident comparatorType ( ',' noncol_ident comparatorType )* )? ')' ( ( K_RETURNS K_NULL ) | ( K_CALLED ) ) K_ON K_NULL K_INPUT K_RETURNS comparatorType K_LANGUAGE IDENT K_AS STRING_LITERAL
    ;

// Parser.g:739
dropFunctionStatement
    : K_DROP K_FUNCTION ( K_IF K_EXISTS )? functionName ( '(' ( comparatorType ( ',' comparatorType )* )? ')' )?
    ;

// Parser.g:763
createKeyspaceStatement
    : K_CREATE K_KEYSPACE ( K_IF K_NOT K_EXISTS )? keyspaceName K_WITH properties
    ;

// Parser.g:779
createTableStatement
    : K_CREATE K_COLUMNFAMILY ( K_IF K_NOT K_EXISTS )? columnFamilyName tableDefinition
    ;

// Parser.g:786
tableDefinition
    : '(' tableColumns ( ',' tableColumns? )* ')' ( K_WITH tableProperty ( K_AND tableProperty )* )?
    ;

// Parser.g:791
tableColumns
    : ident comparatorType ( K_STATIC )? ( columnMask )? ( K_PRIMARY K_KEY )?
    | K_PRIMARY K_KEY '(' tablePartitionKey ( ',' ident )* ')'
    ;

// Parser.g:798
columnMask
    : K_MASKED K_WITH functionName columnMaskArguments
    | K_MASKED K_WITH K_DEFAULT
    ;

// Parser.g:804
columnMaskArguments
    : '(' ')'
    | '(' term ( ',' term )* ')'
    ;

// Parser.g:808
tablePartitionKey
    : ident
    | '(' ident ( ',' ident )* ')'
    ;

// Parser.g:815
tableProperty
    : property
    | K_COMPACT K_STORAGE
    | K_CLUSTERING K_ORDER K_BY '(' tableClusteringOrder ( ',' tableClusteringOrder )* ')'
    ;

// Parser.g:821
tableClusteringOrder
    : ident ( K_ASC | K_DESC )
    ;

// Parser.g:833
createTypeStatement
    : K_CREATE K_TYPE ( K_IF K_NOT K_EXISTS )? userTypeName '(' typeColumns ( ',' typeColumns? )* ')'
    ;

// Parser.g:840
typeColumns
    : fident comparatorType
    ;

// Parser.g:848
createIndexStatement
    : K_CREATE ( K_CUSTOM )? K_INDEX ( K_IF K_NOT K_EXISTS )? ( idxName )? K_ON columnFamilyName '(' ( indexIdent ( ',' indexIdent )* )? ')' ( K_USING STRING_LITERAL )? ( K_WITH properties )?
    ;

// Parser.g:862
indexIdent
    : cident
    | K_VALUES '(' cident ')'
    | K_KEYS '(' cident ')'
    | K_ENTRIES '(' cident ')'
    | K_FULL '(' cident ')'
    ;

// Parser.g:878
createMaterializedViewStatement
    : K_CREATE K_MATERIALIZED K_VIEW ( K_IF K_NOT K_EXISTS )? columnFamilyName K_AS K_SELECT selectors K_FROM columnFamilyName ( K_WHERE whereClause )? viewPrimaryKey ( K_WITH viewProperty ( K_AND viewProperty )* )?
    ;

// Parser.g:893
viewPrimaryKey
    : K_PRIMARY K_KEY '(' viewPartitionKey ( ',' ident )* ')'
    ;

// Parser.g:897
viewPartitionKey
    : ident
    | '(' ident ( ',' ident )* ')'
    ;

// Parser.g:904
// Parser.g:906: alternative dropped, its action refuses every input: K_COMPACT K_STORAGE {throw new SyntaxException("COMPACT STORAGE tables are not allowed starting with version 4.0");}
viewProperty
    : property
    | K_CLUSTERING K_ORDER K_BY '(' viewClusteringOrder ( ',' viewClusteringOrder )* ')'
    ;

// Parser.g:910
viewClusteringOrder
    : ident ( K_ASC | K_DESC )
    ;

// Parser.g:918
createTriggerStatement
    : K_CREATE K_TRIGGER ( K_IF K_NOT K_EXISTS )? ( ident ) K_ON columnFamilyName K_USING STRING_LITERAL
    ;

// Parser.g:930
dropTriggerStatement
    : K_DROP K_TRIGGER ( K_IF K_EXISTS )? ( ident ) K_ON columnFamilyName
    ;

// Parser.g:939
alterKeyspaceStatement
    : K_ALTER K_KEYSPACE ( K_IF K_EXISTS )? keyspaceName K_WITH properties
    ;

// Parser.g:957
alterTableStatement
    : K_ALTER K_COLUMNFAMILY ( K_IF K_EXISTS )? columnFamilyName ( K_ALTER cident K_TYPE comparatorType | K_ALTER ( K_IF K_EXISTS )? cident ( columnMask | K_DROP K_MASKED ) | K_ADD ( K_IF K_NOT K_EXISTS )? ( ident comparatorType isStaticColumn ( columnMask )? | ( '(' ident comparatorType isStaticColumn ( columnMask )? ( ',' ident comparatorType isStaticColumn ( columnMask )? )* ')' ) ) | K_DROP ( K_IF K_EXISTS )? ( ident | ( '(' ident ( ',' ident )* ')' ) ) ( K_USING K_TIMESTAMP INTEGER )? | K_RENAME ( K_IF K_EXISTS )? ( ident K_TO ident ( K_AND ident K_TO ident )* ) | K_DROP K_COMPACT K_STORAGE | K_WITH properties )
    ;

// Parser.g:989
isStaticColumn
    : ( K_STATIC )?
    ;

// Parser.g:994
alterMaterializedViewStatement
    : K_ALTER K_MATERIALIZED K_VIEW ( K_IF K_EXISTS )? columnFamilyName K_WITH properties
    ;

// Parser.g:1012
alterTypeStatement
    : K_ALTER K_TYPE ( K_IF K_EXISTS )? userTypeName ( K_ALTER fident K_TYPE comparatorType | K_ADD ( K_IF K_NOT K_EXISTS )? fident comparatorType | K_RENAME ( K_IF K_EXISTS )? fident K_TO fident ( K_AND fident K_TO fident )* )
    ;

// Parser.g:1030
dropKeyspaceStatement
    : K_DROP K_KEYSPACE ( K_IF K_EXISTS )? keyspaceName
    ;

// Parser.g:1038
dropTableStatement
    : K_DROP K_COLUMNFAMILY ( K_IF K_EXISTS )? columnFamilyName
    ;

// Parser.g:1046
dropTypeStatement
    : K_DROP K_TYPE ( K_IF K_EXISTS )? userTypeName
    ;

// Parser.g:1054
dropIndexStatement
    : K_DROP K_INDEX ( K_IF K_EXISTS )? indexName
    ;

// Parser.g:1063
dropMaterializedViewStatement
    : K_DROP K_MATERIALIZED K_VIEW ( K_IF K_EXISTS )? columnFamilyName
    ;

// Parser.g:1072
truncateStatement
    : K_TRUNCATE ( K_COLUMNFAMILY )? columnFamilyName
    ;

// Parser.g:1079
grantPermissionsStatement
    : K_GRANT permissionOrAll K_ON resource K_TO userOrRoleName
    ;

// Parser.g:1092
revokePermissionsStatement
    : K_REVOKE permissionOrAll K_ON resource K_FROM userOrRoleName
    ;

// Parser.g:1105
grantRoleStatement
    : K_GRANT userOrRoleName K_TO userOrRoleName
    ;

// Parser.g:1116
revokeRoleStatement
    : K_REVOKE userOrRoleName K_FROM userOrRoleName
    ;

// Parser.g:1124
listPermissionsStatement
    : K_LIST permissionOrAll ( K_ON resource )? ( K_OF roleName )? ( K_NORECURSIVE )?
    ;

// Parser.g:1138
permission
    : ( K_CREATE | K_ALTER | K_DROP | K_SELECT | K_MODIFY | K_AUTHORIZE | K_DESCRIBE | K_EXECUTE | K_UNMASK | K_SELECT_MASKED )
    ;

// Parser.g:1143
permissionOrAll
    : K_ALL ( K_PERMISSIONS )?
    | permission ( K_PERMISSION )? ( ',' permission ( K_PERMISSION )? )*
    ;

// Parser.g:1148
resource
    : dataResource
    | roleResource
    | functionResource
    | jmxResource
    ;

// Parser.g:1155
dataResource
    : K_ALL K_KEYSPACES
    | K_KEYSPACE keyspaceName
    | ( K_COLUMNFAMILY )? columnFamilyName
    | K_ALL K_TABLES K_IN K_KEYSPACE keyspaceName
    ;

// Parser.g:1162
jmxResource
    : K_ALL K_MBEANS
    | K_MBEAN mbean
    | K_MBEANS mbean
    ;

// Parser.g:1170
roleResource
    : K_ALL K_ROLES
    | K_ROLE userOrRoleName
    ;

// Parser.g:1175
functionResource
    : K_ALL K_FUNCTIONS
    | K_ALL K_FUNCTIONS K_IN K_KEYSPACE keyspaceName
    | K_FUNCTION functionName ( '(' ( comparatorType ( ',' comparatorType )* )? ')' )
    ;

// Parser.g:1197
// Parser.g:1208: action may refuse the input: {opts.setOption(IRoleManager.Option.SUPERUSER, superuser); if (opts.getPassword().isPresent() && opts.getHashedPassword().isPresent()) { throw new SyntaxException("Options 'password' and 'hashed password' are mutually exclusive"); } $stmt = new CreateRoleStatement(name, opts, DCPermissions.all(), CIDRPermissions.all(), ifNotExists);}
createUserStatement
    : K_CREATE K_USER ( K_IF K_NOT K_EXISTS )? username ( K_WITH userPassword )? ( K_SUPERUSER | K_NOSUPERUSER )?
    ;

// Parser.g:1219
// Parser.g:1229: action may refuse the input: {if (opts.getPassword().isPresent() && opts.getHashedPassword().isPresent()) { throw new SyntaxException("Options 'password' and 'hashed password' are mutually exclusive"); } $stmt = new AlterRoleStatement(name, opts, null, null, ifExists);}
alterUserStatement
    : K_ALTER K_USER ( K_IF K_EXISTS )? username ( K_WITH userPassword )? ( K_SUPERUSER | K_NOSUPERUSER )?
    ;

// Parser.g:1241
dropUserStatement
    : K_DROP K_USER ( K_IF K_EXISTS )? username
    ;

// Parser.g:1251
addIdentityStatement
    : K_ADD K_IDENTITY ( K_IF K_NOT K_EXISTS )? identity K_TO K_ROLE identity
    ;

// Parser.g:1263
dropIdentityStatement
    : K_DROP K_IDENTITY ( K_IF K_EXISTS )? identity
    ;

// Parser.g:1274
listUsersStatement
    : K_LIST K_USERS
    ;

// Parser.g:1291
// Parser.g:1300: action may refuse the input: {// set defaults if they weren't explictly supplied if (!opts.getLogin().isPresent()) { opts.setOption(IRoleManager.Option.LOGIN, false); } if (!opts.getSuperuser().isPresent()) { opts.setOption(IRoleManager.Option.SUPERUSER, false); } if (opts.getPassword().isPresent() && opts.getHashedPassword().isPresent()) { throw new SyntaxException("Options 'password' and 'hashed password' are mutually exclusive"); } $stmt = new CreateRoleStatement(name, opts, dcperms.build(), cidrperms.build(), ifNotExists);}
createRoleStatement
    : K_CREATE K_ROLE ( K_IF K_NOT K_EXISTS )? userOrRoleName ( K_WITH roleOptions )?
    ;

// Parser.g:1331
// Parser.g:1340: action may refuse the input: {if (opts.getPassword().isPresent() && opts.getHashedPassword().isPresent()) { throw new SyntaxException("Options 'password' and 'hashed password' are mutually exclusive"); } $stmt = new AlterRoleStatement(name, opts, dcperms.isModified() ? dcperms.build() : null, cidrperms.isModified() ? cidrperms.build() : null, ifExists);}
alterRoleStatement
    : K_ALTER K_ROLE ( K_IF K_EXISTS )? userOrRoleName ( K_WITH roleOptions )?
    ;

// Parser.g:1352
dropRoleStatement
    : K_DROP K_ROLE ( K_IF K_EXISTS )? userOrRoleName
    ;

// Parser.g:1363
listRolesStatement
    : K_LIST K_ROLES ( K_OF roleName )? ( K_NORECURSIVE )?
    ;

// Parser.g:1374
roleOptions
    : roleOption ( K_AND roleOption )*
    ;

// Parser.g:1378
roleOption
    : K_PASSWORD '=' STRING_LITERAL
    | K_HASHED K_PASSWORD '=' STRING_LITERAL
    | K_OPTIONS '=' fullMapLiteral
    | K_SUPERUSER '=' BOOLEAN
    | K_LOGIN '=' BOOLEAN
    | K_ACCESS K_TO K_ALL K_DATACENTERS
    | K_ACCESS K_TO K_DATACENTERS '{' dcPermission ( ',' dcPermission )* '}'
    | K_ACCESS K_FROM K_ALL K_CIDRS
    | K_ACCESS K_FROM K_CIDRS '{' cidrPermission ( ',' cidrPermission )* '}'
    ;

// Parser.g:1390
dcPermission
    : STRING_LITERAL
    ;

// Parser.g:1394
cidrPermission
    : STRING_LITERAL
    ;

// Parser.g:1399
userPassword
    : K_PASSWORD STRING_LITERAL
    | K_HASHED K_PASSWORD STRING_LITERAL
    ;

// Parser.g:1410
describeStatement
    : ( K_DESCRIBE | K_DESC ) ( K_CLUSTER | ( K_FULL )? K_SCHEMA | K_KEYSPACES | ( K_ONLY )? K_KEYSPACE ( keyspaceName )? | K_TABLES | K_COLUMNFAMILY columnFamilyName | K_INDEX columnFamilyName | K_MATERIALIZED K_VIEW columnFamilyName | K_TYPES | K_TYPE userTypeName | K_FUNCTIONS | K_FUNCTION functionName | K_AGGREGATES | K_AGGREGATE functionName | ( ( IDENT | QUOTED_NAME | unreserved_keyword ) '.' )? ( IDENT | QUOTED_NAME | unreserved_keyword ) ) ( K_WITH K_INTERNALS )?
    ;

// Parser.g:1449
cident
    : EMPTY_QUOTED_NAME
    | ident
    ;

// Parser.g:1454
ident
    : IDENT
    | QUOTED_NAME
    | unreserved_keyword
    ;

// Parser.g:1460
fident
    : IDENT
    | QUOTED_NAME
    | unreserved_keyword
    ;

// Parser.g:1467
noncol_ident
    : IDENT
    | QUOTED_NAME
    | unreserved_keyword
    ;

// Parser.g:1474
keyspaceName
    : ksName
    ;

// Parser.g:1479
indexName
    : ( ksName '.' )? idxName
    ;

// Parser.g:1484
columnFamilyName
    : ( ksName '.' )? cfName
    ;

// Parser.g:1489
userTypeName
    : ( noncol_ident '.' )? non_type_ident
    ;

// Parser.g:1493
userOrRoleName
    : roleName
    ;

// Parser.g:1498
// Parser.g:1502: alternative dropped, its action refuses every input: QMARK {addRecognitionError("Bind variables cannot be used for keyspace names");}
ksName
    : IDENT
    | QUOTED_NAME
    | unreserved_keyword
    ;

// Parser.g:1505
// Parser.g:1509: alternative dropped, its action refuses every input: QMARK {addRecognitionError("Bind variables cannot be used for table names");}
cfName
    : IDENT
    | QUOTED_NAME
    | unreserved_keyword
    ;

// Parser.g:1512
// Parser.g:1516: alternative dropped, its action refuses every input: QMARK {addRecognitionError("Bind variables cannot be used for index names");}
idxName
    : IDENT
    | QUOTED_NAME
    | unreserved_keyword
    ;

// Parser.g:1519
// Parser.g:1524: alternative dropped, its action refuses every input: QMARK {addRecognitionError("Bind variables cannot be used for role names");}
roleName
    : IDENT
    | STRING_LITERAL
    | QUOTED_NAME
    | unreserved_keyword
    ;

// Parser.g:1527
constant
    : STRING_LITERAL
    | INTEGER
    | FLOAT
    | BOOLEAN
    | DURATION
    | UUID
    | HEXNUMBER
    | ( ( K_POSITIVE_NAN | K_NEGATIVE_NAN ) | K_POSITIVE_INFINITY | K_NEGATIVE_INFINITY )
    ;

// Parser.g:1540
fullMapLiteral
    : '{' ( term ':' term ( ',' term ':' term )* )? '}'
    ;

// Parser.g:1547
setOrMapLiteral
    : mapLiteral
    | setLiteral
    ;

// Parser.g:1552
setLiteral
    : ( ',' term )*
    ;

// Parser.g:1558
mapLiteral
    : ':' term ( ',' term ':' term )*
    ;

// Parser.g:1564
collectionLiteral
    : listLiteral
    | '{' term setOrMapLiteral '}'
    | '{' '}'
    ;

// Parser.g:1572
listLiteral
    : '[' ( term ( ',' term )* )? ']'
    ;

// Parser.g:1578
usertypeLiteral
    : '{' fident ':' term ( ',' fident ':' term )* '}'
    ;

// Parser.g:1585
tupleLiteral
    : '(' term ( ',' term )* ')'
    ;

// Parser.g:1591
value
    : constant
    | collectionLiteral
    | usertypeLiteral
    | tupleLiteral
    | K_NULL
    | ':' noncol_ident
    | QMARK
    ;

// Parser.g:1601
intValue
    : INTEGER
    | ':' noncol_ident
    | QMARK
    ;

// Parser.g:1607
functionName
    : ( keyspaceName '.' )? allowedFunctionName
    ;

// Parser.g:1613
allowedFunctionName
    : IDENT
    | QUOTED_NAME
    | unreserved_function_keyword
    | K_TOKEN
    | K_COUNT
    ;

// Parser.g:1621
function
    : functionName '(' ')'
    | functionName '(' functionArgs ')'
    ;

// Parser.g:1626
functionArgs
    : term ( ',' term )*
    ;

// Parser.g:1631
term
    : termAddition
    ;

// Parser.g:1635
termAddition
    : termMultiplication ( '+' termMultiplication | '-' termMultiplication )*
    ;

// Parser.g:1642
termMultiplication
    : termGroup ( '*' termGroup | '/' termGroup | '%' termGroup )*
    ;

// Parser.g:1650
termGroup
    : simpleTerm
    | '-' simpleTerm
    ;

// Parser.g:1655
simpleTerm
    : value
    | function
    | '(' comparatorType ')' simpleTerm
    | K_CAST '(' simpleTerm K_AS native_type ')'
    ;

// Parser.g:1662
columnOperation
    : cident columnOperationDifferentiator
    ;

// Parser.g:1666
columnOperationDifferentiator
    : '=' normalColumnOperation
    | shorthandColumnOperation
    | '[' term ']' collectionColumnOperation
    | '.' fident udtColumnOperation
    ;

// Parser.g:1673
// Parser.g:1675: action may refuse the input: {if (c == null) { addRawUpdate(operations, key, new Operation.SetValue(t)); } else { if (!key.equals(c)) addRecognitionError("Only expressions of the form X = <value> + X are supported."); addRawUpdate(operations, key, new Operation.Prepend(t)); }}
// Parser.g:1688: action may refuse the input: {if (!key.equals(c)) addRecognitionError("Only expressions of the form X = X " + $sig.text + "<value> are supported."); addRawUpdate(operations, key, $sig.text.equals("+") ? new Operation.Addition(t) : new Operation.Substraction(t));}
// Parser.g:1694: action may refuse the input: {// Note that this production *is* necessary because X = X - 3 will in fact be lexed as [ X, '=', X, INTEGER]. if (!key.equals(c)) // We don't yet allow a '+' in front of an integer, but we could in the future really, so let's be future-proof in our error message addRecognitionError("Only expressions of the form X = X " + ($i.text.charAt(0) == '-' ? '-' : '+') + " <value> are supported."); addRawUpdate(operations, key, new Operation.Addition(Constants.Literal.integer($i.text)));}
normalColumnOperation
    : term ( '+' cident )?
    | cident ( '+' | '-' ) term
    | cident INTEGER
    ;

// Parser.g:1703
shorthandColumnOperation
    : ( '+=' | '-=' ) term
    ;

// Parser.g:1710
collectionColumnOperation
    : '=' term
    ;

// Parser.g:1717
udtColumnOperation
    : '=' term
    ;

// Parser.g:1724
columnCondition
    : cident ( relationType term | containsOperator term | K_IN ( singleColumnInValues | inMarker ) | '[' term ']' ( relationType term | K_IN ( singleColumnInValues | inMarker ) ) | '.' fident ( relationType term | K_IN ( singleColumnInValues | inMarker ) ) )
    ;

// Parser.g:1750
properties
    : property ( K_AND property )*
    ;

// Parser.g:1754
// Parser.g:1755: action may refuse the input: {try { $props.addProperty(k.toString(), simple); } catch (SyntaxException e) { addRecognitionError(e.getMessage()); }}
// Parser.g:1756: action may refuse the input: {try { $props.addProperty(k.toString(), convertPropertyMap(map)); } catch (SyntaxException e) { addRecognitionError(e.getMessage()); }}
property
    : noncol_ident '=' propertyValue
    | noncol_ident '=' fullMapLiteral
    ;

// Parser.g:1759
propertyValue
    : constant
    | unreserved_keyword
    ;

// Parser.g:1764
relationType
    : '='
    | '<'
    | '<='
    | '>'
    | '>='
    | '!='
    ;

// Parser.g:1773
relation
    : cident relationType term
    | cident K_LIKE term
    | cident K_IS K_NOT K_NULL
    | K_TOKEN tupleOfIdentifiers relationType term
    | cident K_IN inMarker
    | cident K_IN singleColumnInValues
    | cident containsOperator term
    | cident '[' term ']' relationType term
    | tupleOfIdentifiers ( K_IN ( '(' ')' | inMarkerForTuple | tupleOfTupleLiterals | tupleOfMarkersForTuples ) | relationType tupleLiteral | relationType markerForTuple )
    | '(' relation ')'
    ;

// Parser.g:1808
containsOperator
    : K_CONTAINS ( K_KEY )?
    ;

// Parser.g:1812
inMarker
    : QMARK
    | ':' noncol_ident
    ;

// Parser.g:1817
tupleOfIdentifiers
    : '(' cident ( ',' cident )* ')'
    ;

// Parser.g:1822
singleColumnInValues
    : '(' ( term ( ',' term )* )? ')'
    ;

// Parser.g:1827
tupleOfTupleLiterals
    : '(' tupleLiteral ( ',' tupleLiteral )* ')'
    ;

// Parser.g:1832
markerForTuple
    : QMARK
    | ':' noncol_ident
    ;

// Parser.g:1837
tupleOfMarkersForTuples
    : '(' markerForTuple ( ',' markerForTuple )* ')'
    ;

// Parser.g:1842
inMarkerForTuple
    : QMARK
    | ':' noncol_ident
    ;

// Parser.g:1847
// Parser.g:1854: action may refuse the input: {try { $t = f.freeze(); } catch (InvalidRequestException e) { addRecognitionError(e.getMessage()); }}
// Parser.g:1862: action may refuse the input: {try { $t = CQL3Type.Raw.from(new CQL3Type.Custom($s.text)); } catch (SyntaxException e) { addRecognitionError("Cannot parse type " + $s.text + ": " + e.getMessage()); } catch (ConfigurationException e) { addRecognitionError("Error setting type " + $s.text + ": " + e.getMessage()); }}
comparatorType
    : native_type
    | collection_type
    | tuple_type
    | vector_type
    | userTypeName
    | K_FROZEN '<' comparatorType '>'
    | STRING_LITERAL
    ;

// Parser.g:1873
native_type
    : K_ASCII
    | K_BIGINT
    | K_BLOB
    | K_BOOLEAN
    | K_COUNTER
    | K_DECIMAL
    | K_DOUBLE
    | K_DURATION
    | K_FLOAT
    | K_INET
    | K_INT
    | K_SMALLINT
    | K_TEXT
    | K_TIMESTAMP
    | K_TINYINT
    | K_UUID
    | K_VARCHAR
    | K_VARINT
    | K_TIMEUUID
    | K_DATE
    | K_TIME
    ;

// Parser.g:1897
collection_type
    : K_MAP '<' comparatorType ',' comparatorType '>'
    | K_LIST '<' comparatorType '>'
    | K_SET '<' comparatorType '>'
    ;

// Parser.g:1910
tuple_type
    : K_TUPLE '<' comparatorType ( ',' comparatorType )* '>'
    ;

// Parser.g:1916
vector_type
    : K_VECTOR '<' comparatorType ',' INTEGER '>'
    ;

// Parser.g:1921
// Parser.g:1924: alternative dropped, its action refuses every input: QUOTED_NAME {addRecognitionError("Quoted strings are are not supported for user names and USER is deprecated, please use ROLE");}
username
    : IDENT
    | STRING_LITERAL
    ;

// Parser.g:1927
// Parser.g:1930: alternative dropped, its action refuses every input: QUOTED_NAME {addRecognitionError("Quoted strings are are not supported for identity");}
identity
    : IDENT
    | STRING_LITERAL
    ;

// Parser.g:1933
mbean
    : STRING_LITERAL
    ;

// Parser.g:1939
// Parser.g:1940: action may refuse the input: {if (reservedTypeNames.contains($t.text)) addRecognitionError("Invalid (reserved) user type name " + $t.text); $id = new ColumnIdentifier($t.text, false);}
non_type_ident
    : IDENT
    | QUOTED_NAME
    | basic_unreserved_keyword
    | K_KEY
    ;

// Parser.g:1946
unreserved_keyword
    : unreserved_function_keyword
    | ( K_TTL | K_COUNT | K_WRITETIME | K_MAXWRITETIME | K_KEY | K_CAST | K_JSON | K_DISTINCT )
    ;

// Parser.g:1951
unreserved_function_keyword
    : basic_unreserved_keyword
    | native_type
    ;

// Parser.g:1956
basic_unreserved_keyword
    : ( K_KEYS | K_AS | K_CLUSTER | K_CLUSTERING | K_COMPACT | K_STORAGE | K_TABLES | K_TYPE | K_TYPES | K_VALUES | K_MAP | K_LIST | K_FILTERING | K_PERMISSION | K_PERMISSIONS | K_KEYSPACES | K_ALL | K_USER | K_USERS | K_ROLE | K_ROLES | K_IDENTITY | K_SUPERUSER | K_NOSUPERUSER | K_LOGIN | K_NOLOGIN | K_OPTIONS | K_PASSWORD | K_HASHED | K_EXISTS | K_CUSTOM | K_TRIGGER | K_CONTAINS | K_INTERNALS | K_ONLY | K_STATIC | K_FROZEN | K_TUPLE | K_FUNCTION | K_FUNCTIONS | K_AGGREGATE | K_AGGREGATES | K_SFUNC | K_STYPE | K_FINALFUNC | K_INITCOND | K_RETURNS | K_LANGUAGE | K_CALLED | K_INPUT | K_LIKE | K_PER | K_PARTITION | K_GROUP | K_DATACENTERS | K_CIDRS | K_ACCESS | K_DEFAULT | K_MBEAN | K_MBEANS | K_REPLACE | K_UNSET | K_MASKED | K_UNMASK | K_SELECT_MASKED | K_VECTOR | K_ANN )
    ;

// Lexer.g:62
K_SELECT
    : S E L E C T
    ;

// Lexer.g:63
K_FROM
    : F R O M
    ;

// Lexer.g:64
K_AS
    : A S
    ;

// Lexer.g:65
K_WHERE
    : W H E R E
    ;

// Lexer.g:66
K_AND
    : A N D
    ;

// Lexer.g:67
K_KEY
    : K E Y
    ;

// Lexer.g:68
K_KEYS
    : K E Y S
    ;

// Lexer.g:69
K_ENTRIES
    : E N T R I E S
    ;

// Lexer.g:70
K_FULL
    : F U L L
    ;

// Lexer.g:71
K_INSERT
    : I N S E R T
    ;

// Lexer.g:72
K_UPDATE
    : U P D A T E
    ;

// Lexer.g:73
K_WITH
    : W I T H
    ;

// Lexer.g:74
K_LIMIT
    : L I M I T
    ;

// Lexer.g:75
K_PER
    : P E R
    ;

// Lexer.g:76
K_PARTITION
    : P A R T I T I O N
    ;

// Lexer.g:77
K_USING
    : U S I N G
    ;

// Lexer.g:78
K_USE
    : U S E
    ;

// Lexer.g:79
K_DISTINCT
    : D I S T I N C T
    ;

// Lexer.g:80
K_COUNT
    : C O U N T
    ;

// Lexer.g:81
K_SET
    : S E T
    ;

// Lexer.g:82
K_BEGIN
    : B E G I N
    ;

// Lexer.g:83
K_UNLOGGED
    : U N L O G G E D
    ;

// Lexer.g:84
K_BATCH
    : B A T C H
    ;

// Lexer.g:85
K_APPLY
    : A P P L Y
    ;

// Lexer.g:86
K_TRUNCATE
    : T R U N C A T E
    ;

// Lexer.g:87
K_DELETE
    : D E L E T E
    ;

// Lexer.g:88
K_IN
    : I N
    ;

// Lexer.g:89
K_CREATE
    : C R E A T E
    ;

// Lexer.g:90
K_SCHEMA
    : S C H E M A
    ;

// Lexer.g:91
K_KEYSPACE
    : ( K E Y S P A C E | K_SCHEMA )
    ;

// Lexer.g:93
K_KEYSPACES
    : K E Y S P A C E S
    ;

// Lexer.g:94
K_COLUMNFAMILY
    : ( C O L U M N F A M I L Y | T A B L E )
    ;

// Lexer.g:96
K_TABLES
    : ( C O L U M N F A M I L I E S | T A B L E S )
    ;

// Lexer.g:98
K_MATERIALIZED
    : M A T E R I A L I Z E D
    ;

// Lexer.g:99
K_VIEW
    : V I E W
    ;

// Lexer.g:100
K_INDEX
    : I N D E X
    ;

// Lexer.g:101
K_CUSTOM
    : C U S T O M
    ;

// Lexer.g:102
K_ON
    : O N
    ;

// Lexer.g:103
K_TO
    : T O
    ;

// Lexer.g:104
K_DROP
    : D R O P
    ;

// Lexer.g:105
K_PRIMARY
    : P R I M A R Y
    ;

// Lexer.g:106
K_INTO
    : I N T O
    ;

// Lexer.g:107
K_VALUES
    : V A L U E S
    ;

// Lexer.g:108
K_TIMESTAMP
    : T I M E S T A M P
    ;

// Lexer.g:109
K_TTL
    : T T L
    ;

// Lexer.g:110
K_CAST
    : C A S T
    ;

// Lexer.g:111
K_ALTER
    : A L T E R
    ;

// Lexer.g:112
K_RENAME
    : R E N A M E
    ;

// Lexer.g:113
K_ADD
    : A D D
    ;

// Lexer.g:114
K_TYPE
    : T Y P E
    ;

// Lexer.g:115
K_TYPES
    : T Y P E S
    ;

// Lexer.g:116
K_COMPACT
    : C O M P A C T
    ;

// Lexer.g:117
K_STORAGE
    : S T O R A G E
    ;

// Lexer.g:118
K_ORDER
    : O R D E R
    ;

// Lexer.g:119
K_BY
    : B Y
    ;

// Lexer.g:120
K_ASC
    : A S C
    ;

// Lexer.g:121
K_DESC
    : D E S C
    ;

// Lexer.g:122
K_ALLOW
    : A L L O W
    ;

// Lexer.g:123
K_FILTERING
    : F I L T E R I N G
    ;

// Lexer.g:124
K_IF
    : I F
    ;

// Lexer.g:125
K_IS
    : I S
    ;

// Lexer.g:126
K_CONTAINS
    : C O N T A I N S
    ;

// Lexer.g:127
K_GROUP
    : G R O U P
    ;

// Lexer.g:128
K_CLUSTER
    : C L U S T E R
    ;

// Lexer.g:129
K_INTERNALS
    : I N T E R N A L S
    ;

// Lexer.g:130
K_ONLY
    : O N L Y
    ;

// Lexer.g:132
K_GRANT
    : G R A N T
    ;

// Lexer.g:133
K_ALL
    : A L L
    ;

// Lexer.g:134
K_PERMISSION
    : P E R M I S S I O N
    ;

// Lexer.g:135
K_PERMISSIONS
    : P E R M I S S I O N S
    ;

// Lexer.g:136
K_OF
    : O F
    ;

// Lexer.g:137
K_REVOKE
    : R E V O K E
    ;

// Lexer.g:138
K_MODIFY
    : M O D I F Y
    ;

// Lexer.g:139
K_AUTHORIZE
    : A U T H O R I Z E
    ;

// Lexer.g:140
K_DESCRIBE
    : D E S C R I B E
    ;

// Lexer.g:141
K_EXECUTE
    : E X E C U T E
    ;

// Lexer.g:142
K_NORECURSIVE
    : N O R E C U R S I V E
    ;

// Lexer.g:143
K_MBEAN
    : M B E A N
    ;

// Lexer.g:144
K_MBEANS
    : M B E A N S
    ;

// Lexer.g:146
K_USER
    : U S E R
    ;

// Lexer.g:147
K_USERS
    : U S E R S
    ;

// Lexer.g:148
K_ROLE
    : R O L E
    ;

// Lexer.g:149
K_ROLES
    : R O L E S
    ;

// Lexer.g:150
K_SUPERUSER
    : S U P E R U S E R
    ;

// Lexer.g:151
K_NOSUPERUSER
    : N O S U P E R U S E R
    ;

// Lexer.g:152
K_PASSWORD
    : P A S S W O R D
    ;

// Lexer.g:153
K_HASHED
    : H A S H E D
    ;

// Lexer.g:154
K_LOGIN
    : L O G I N
    ;

// Lexer.g:155
K_NOLOGIN
    : N O L O G I N
    ;

// Lexer.g:156
K_OPTIONS
    : O P T I O N S
    ;

// Lexer.g:157
K_ACCESS
    : A C C E S S
    ;

// Lexer.g:158
K_DATACENTERS
    : D A T A C E N T E R S
    ;

// Lexer.g:159
K_CIDRS
    : C I D R S
    ;

// Lexer.g:160
K_IDENTITY
    : I D E N T I T Y
    ;

// Lexer.g:162
K_CLUSTERING
    : C L U S T E R I N G
    ;

// Lexer.g:163
K_ASCII
    : A S C I I
    ;

// Lexer.g:164
K_BIGINT
    : B I G I N T
    ;

// Lexer.g:165
K_BLOB
    : B L O B
    ;

// Lexer.g:166
K_BOOLEAN
    : B O O L E A N
    ;

// Lexer.g:167
K_COUNTER
    : C O U N T E R
    ;

// Lexer.g:168
K_DECIMAL
    : D E C I M A L
    ;

// Lexer.g:169
K_DOUBLE
    : D O U B L E
    ;

// Lexer.g:170
K_DURATION
    : D U R A T I O N
    ;

// Lexer.g:171
K_FLOAT
    : F L O A T
    ;

// Lexer.g:172
K_INET
    : I N E T
    ;

// Lexer.g:173
K_INT
    : I N T
    ;

// Lexer.g:174
K_SMALLINT
    : S M A L L I N T
    ;

// Lexer.g:175
K_TINYINT
    : T I N Y I N T
    ;

// Lexer.g:176
K_TEXT
    : T E X T
    ;

// Lexer.g:177
K_UUID
    : U U I D
    ;

// Lexer.g:178
K_VARCHAR
    : V A R C H A R
    ;

// Lexer.g:179
K_VARINT
    : V A R I N T
    ;

// Lexer.g:180
K_TIMEUUID
    : T I M E U U I D
    ;

// Lexer.g:181
K_TOKEN
    : T O K E N
    ;

// Lexer.g:182
K_WRITETIME
    : W R I T E T I M E
    ;

// Lexer.g:183
K_MAXWRITETIME
    : M A X W R I T E T I M E
    ;

// Lexer.g:184
K_DATE
    : D A T E
    ;

// Lexer.g:185
K_TIME
    : T I M E
    ;

// Lexer.g:187
K_NULL
    : N U L L
    ;

// Lexer.g:188
K_NOT
    : N O T
    ;

// Lexer.g:189
K_EXISTS
    : E X I S T S
    ;

// Lexer.g:191
K_MAP
    : M A P
    ;

// Lexer.g:192
K_LIST
    : L I S T
    ;

// Lexer.g:193
K_POSITIVE_NAN
    : N A N
    ;

// Lexer.g:194
K_NEGATIVE_NAN
    : '-' N A N
    ;

// Lexer.g:195
K_POSITIVE_INFINITY
    : I N F I N I T Y
    ;

// Lexer.g:196
K_NEGATIVE_INFINITY
    : '-' I N F I N I T Y
    ;

// Lexer.g:197
K_TUPLE
    : T U P L E
    ;

// Lexer.g:199
K_TRIGGER
    : T R I G G E R
    ;

// Lexer.g:200
K_STATIC
    : S T A T I C
    ;

// Lexer.g:201
K_FROZEN
    : F R O Z E N
    ;

// Lexer.g:203
K_FUNCTION
    : F U N C T I O N
    ;

// Lexer.g:204
K_FUNCTIONS
    : F U N C T I O N S
    ;

// Lexer.g:205
K_AGGREGATE
    : A G G R E G A T E
    ;

// Lexer.g:206
K_AGGREGATES
    : A G G R E G A T E S
    ;

// Lexer.g:207
K_SFUNC
    : S F U N C
    ;

// Lexer.g:208
K_STYPE
    : S T Y P E
    ;

// Lexer.g:209
K_FINALFUNC
    : F I N A L F U N C
    ;

// Lexer.g:210
K_INITCOND
    : I N I T C O N D
    ;

// Lexer.g:211
K_RETURNS
    : R E T U R N S
    ;

// Lexer.g:212
K_CALLED
    : C A L L E D
    ;

// Lexer.g:213
K_INPUT
    : I N P U T
    ;

// Lexer.g:214
K_LANGUAGE
    : L A N G U A G E
    ;

// Lexer.g:215
K_OR
    : O R
    ;

// Lexer.g:216
K_REPLACE
    : R E P L A C E
    ;

// Lexer.g:218
K_JSON
    : J S O N
    ;

// Lexer.g:219
K_DEFAULT
    : D E F A U L T
    ;

// Lexer.g:220
K_UNSET
    : U N S E T
    ;

// Lexer.g:221
K_LIKE
    : L I K E
    ;

// Lexer.g:223
K_MASKED
    : M A S K E D
    ;

// Lexer.g:224
K_UNMASK
    : U N M A S K
    ;

// Lexer.g:225
K_SELECT_MASKED
    : S E L E C T '_' M A S K E D
    ;

// Lexer.g:227
K_VECTOR
    : V E C T O R
    ;

// Lexer.g:228
K_ANN
    : A N N
    ;

// Lexer.g:231
fragment A
    : ( 'a' | 'A' )
    ;

// Lexer.g:232
fragment B
    : ( 'b' | 'B' )
    ;

// Lexer.g:233
fragment C
    : ( 'c' | 'C' )
    ;

// Lexer.g:234
fragment D
    : ( 'd' | 'D' )
    ;

// Lexer.g:235
fragment E
    : ( 'e' | 'E' )
    ;

// Lexer.g:236
fragment F
    : ( 'f' | 'F' )
    ;

// Lexer.g:237
fragment G
    : ( 'g' | 'G' )
    ;

// Lexer.g:238
fragment H
    : ( 'h' | 'H' )
    ;

// Lexer.g:239
fragment I
    : ( 'i' | 'I' )
    ;

// Lexer.g:240
fragment J
    : ( 'j' | 'J' )
    ;

// Lexer.g:241
fragment K
    : ( 'k' | 'K' )
    ;

// Lexer.g:242
fragment L
    : ( 'l' | 'L' )
    ;

// Lexer.g:243
fragment M
    : ( 'm' | 'M' )
    ;

// Lexer.g:244
fragment N
    : ( 'n' | 'N' )
    ;

// Lexer.g:245
fragment O
    : ( 'o' | 'O' )
    ;

// Lexer.g:246
fragment P
    : ( 'p' | 'P' )
    ;

// Lexer.g:247
fragment Q
    : ( 'q' | 'Q' )
    ;

// Lexer.g:248
fragment R
    : ( 'r' | 'R' )
    ;

// Lexer.g:249
fragment S
    : ( 's' | 'S' )
    ;

// Lexer.g:250
fragment T
    : ( 't' | 'T' )
    ;

// Lexer.g:251
fragment U
    : ( 'u' | 'U' )
    ;

// Lexer.g:252
fragment V
    : ( 'v' | 'V' )
    ;

// Lexer.g:253
fragment W
    : ( 'w' | 'W' )
    ;

// Lexer.g:254
fragment X
    : ( 'x' | 'X' )
    ;

// Lexer.g:255
fragment Y
    : ( 'y' | 'Y' )
    ;

// Lexer.g:256
fragment Z
    : ( 'z' | 'Z' )
    ;

// Lexer.g:258
// Lexer.g:268: predicate dropped: {(input.size() - input.index() > 1) && !"$$".equals(input.substring(input.index(), input.index() + 1))}?
STRING_LITERAL
    : ( '$' '$' .*? '$' '$' )
    | ( '\'' ( ~( '\'' ) | '\'' '\'' )* '\'' )
    ;

// Lexer.g:281
QUOTED_NAME
    : '"' ( ~( '"' ) | '"' '"' )+ '"'
    ;

// Lexer.g:287
EMPTY_QUOTED_NAME
    : '"' '"'
    ;

// Lexer.g:291
fragment DIGIT
    : '0'..'9'
    ;

// Lexer.g:295
fragment LETTER
    : ( 'A'..'Z' | 'a'..'z' )
    ;

// Lexer.g:299
fragment HEX
    : ( 'A'..'F' | 'a'..'f' | '0'..'9' )
    ;

// Lexer.g:303
fragment EXPONENT
    : E ( '+' | '-' )? DIGIT+
    ;

// Lexer.g:307
fragment DURATION_ISO_8601_PERIOD_DESIGNATORS
    : '-'? 'P' DIGIT+ 'Y' ( DIGIT+ 'M' )? ( DIGIT+ 'D' )?
    | '-'? 'P' DIGIT+ 'M' ( DIGIT+ 'D' )?
    | '-'? 'P' DIGIT+ 'D'
    ;

// Lexer.g:313
fragment DURATION_ISO_8601_TIME_DESIGNATORS
    : 'T' DIGIT+ 'H' ( DIGIT+ 'M' )? ( DIGIT+ 'S' )?
    | 'T' DIGIT+ 'M' ( DIGIT+ 'S' )?
    | 'T' DIGIT+ 'S'
    ;

// Lexer.g:319
fragment DURATION_ISO_8601_WEEK_PERIOD_DESIGNATOR
    : '-'? 'P' DIGIT+ 'W'
    ;

// Lexer.g:323
fragment DURATION_UNIT
    : Y
    | M O
    | W
    | D
    | H
    | M
    | S
    | M S
    | U S
    | '\u00B5' S
    | N S
    ;

// Lexer.g:337
INTEGER
    : '-'? DIGIT+
    ;

// Lexer.g:341
QMARK
    : '?'
    ;

// Lexer.g:345
RANGE
    : '..'
    ;

// Lexer.g:353
// Lexer.g:355: action dropped, it emits another token type: {$type = INTEGER;}
FLOAT
    : INTEGER '.'
    | INTEGER
    | INTEGER ( '.' DIGIT* )? EXPONENT?
    ;

// Lexer.g:362
BOOLEAN
    : T R U E
    | F A L S E
    ;

// Lexer.g:366
DURATION
    : '-'? DIGIT+ DURATION_UNIT ( DIGIT+ DURATION_UNIT )*
    | '-'? 'P' DIGIT DIGIT DIGIT DIGIT '-' DIGIT DIGIT '-' DIGIT DIGIT 'T' DIGIT DIGIT ':' DIGIT DIGIT ':' DIGIT DIGIT
    | '-'? 'P' DURATION_ISO_8601_TIME_DESIGNATORS
    | DURATION_ISO_8601_WEEK_PERIOD_DESIGNATOR
    | DURATION_ISO_8601_PERIOD_DESIGNATORS DURATION_ISO_8601_TIME_DESIGNATORS?
    ;

// Lexer.g:374
IDENT
    : LETTER ( LETTER | DIGIT | '_' )*
    ;

// Lexer.g:378
HEXNUMBER
    : '0' X HEX*
    ;

// Lexer.g:382
UUID
    : HEX HEX HEX HEX HEX HEX HEX HEX '-' HEX HEX HEX HEX '-' HEX HEX HEX HEX '-' HEX HEX HEX HEX '-' HEX HEX HEX HEX HEX HEX HEX HEX HEX HEX HEX HEX
    ;

// Lexer.g:390
WS
    : ( ' ' | '\t' | '\n' | '\r' )+
      -> channel(HIDDEN)
    ;

// Lexer.g:394
COMMENT
    : ( '--' | '//' ) .*? ( '\n' | '\r' )
      -> channel(HIDDEN)
    ;

// Lexer.g:398
MULTILINE_COMMENT
    : '/*' .*? '*/'
      -> channel(HIDDEN)
    ;
