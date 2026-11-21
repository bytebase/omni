/*
 * Copyright MongoDB, Inc. (mongodb-js/mongosh, Apache-2.0)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// MongoDB Shell (mongosh) call-shape grammar as ANTLR 4, converted from the
// official mongodb-js/mongosh TypeScript shell-api and shell-bson definitions
// by mongo/parser/cmd/mongoshgrammar. Do not edit by hand: rerun the converter.
//
// Scope (§2.2 of docs/specs/2026-11-11-engine-compatibility-without-servers-design.md):
// covers mongosh call shapes (direct shell commands, ShellApi/ShellConfig calls,
// db.*, db.<coll>.*, cursor/explain/bulk/planCache chains, rs.*, sh.*, sp.*,
// Mongo/connect/session/encryption chains, BSON constructors & static helpers,
// and representative document/array/literal argument shapes) rather than
// arbitrary JavaScript syntax.
//
// pin: mongosh github.com/mongodb-js/mongosh@82c228b597c3221de59bdbfaff4c346622780cdd sha256:8a62b35159f7080fd1defc4fd63c4b0f75f91a62e50e56cc1855b9e62f27e087 v2.14.3

grammar MongoShell;

statement
    : directShellCommand
    | shellApiCall
    | shellConfigCall
    | dbStatement
    | collectionStatement
    | explainableStatement
    | bulkStatement
    | planCacheStatement
    | rsStatement
    | shStatement
    | spStatement
    | mongoConnectionStatement
    | encryptionStatement
    | sessionStatement
    | ejsonStatement
    ;

// packages/shell-api/src/shell-api.ts (@directShellCommand) and packages/shell-api/src/mongo.ts (Mongo.show)
directShellCommand
    : 'use' IDENT // packages/shell-api/src/shell-api.ts:248
    | 'show' showTarget // packages/shell-api/src/shell-api.ts:255
    | 'show' 'log' IDENT // packages/shell-api/src/shell-api.ts:255 (show(cmd, arg?))
    | 'exit' // packages/shell-api/src/shell-api.ts:267
    | 'exit' NUMBER // packages/shell-api/src/shell-api.ts:267 (exitCode?: number)
    | 'quit' // packages/shell-api/src/shell-api.ts:282
    | 'quit' NUMBER // packages/shell-api/src/shell-api.ts:282 (exitCode?: number)
    | 'it' // packages/shell-api/src/shell-api.ts:330
    | 'cls' // packages/shell-api/src/shell-api.ts:455
    ;

// packages/shell-api/src/mongo.ts (Mongo.show)
showTarget
    : 'databases' // packages/shell-api/src/mongo.ts:457
    | 'dbs' // packages/shell-api/src/mongo.ts:458
    | 'collections' // packages/shell-api/src/mongo.ts:466
    | 'tables' // packages/shell-api/src/mongo.ts:467
    | 'profile' // packages/shell-api/src/mongo.ts:473
    | 'users' // packages/shell-api/src/mongo.ts:485
    | 'roles' // packages/shell-api/src/mongo.ts:488
    | 'log' // packages/shell-api/src/mongo.ts:491
    | 'logs' // packages/shell-api/src/mongo.ts:494
    | 'startupWarnings' // packages/shell-api/src/mongo.ts:497
    | 'automationNotices' // packages/shell-api/src/mongo.ts:539
    | 'nonGenuineMongoDBCheck' // packages/shell-api/src/mongo.ts:560
    ;

// packages/shell-api/src/shell-api.ts (ShellApi)
shellApiCall
    : shellApiMethod '(' argList? ')'
    ;

shellApiMethod
    : 'use' // packages/shell-api/src/shell-api.ts:248
    | 'show' // packages/shell-api/src/shell-api.ts:255
    | 'exit' // packages/shell-api/src/shell-api.ts:267
    | 'quit' // packages/shell-api/src/shell-api.ts:282
    | 'it' // packages/shell-api/src/shell-api.ts:330
    | 'version' // packages/shell-api/src/shell-api.ts:337
    | 'load' // packages/shell-api/src/shell-api.ts:342
    | 'enableTelemetry' // packages/shell-api/src/shell-api.ts:379
    | 'disableTelemetry' // packages/shell-api/src/shell-api.ts:395
    | 'passwordPrompt' // packages/shell-api/src/shell-api.ts:411
    | 'sleep' // packages/shell-api/src/shell-api.ts:423
    | 'print' // packages/shell-api/src/shell-api.ts:439
    | 'printjson' // packages/shell-api/src/shell-api.ts:444
    | 'convertShardKeyToHashed' // packages/shell-api/src/shell-api.ts:449
    | 'cls' // packages/shell-api/src/shell-api.ts:455
    | 'isInteractive' // packages/shell-api/src/shell-api.ts:460
    ;

// packages/shell-api/src/shell-api.ts (ShellConfig)
shellConfigCall
    : 'config' '.' shellConfigMethod '(' argList? ')'
    ;

shellConfigMethod
    : 'set' // packages/shell-api/src/shell-api.ts:59
    | 'get' // packages/shell-api/src/shell-api.ts:88
    | 'reset' // packages/shell-api/src/shell-api.ts:97
    ;

// packages/shell-api/src/database.ts (Database)
dbStatement
    : 'db' '.' databaseMethod '(' argList? ')'
    | 'db' '.' 'getSiblingDB' '(' STRING ')' '.' databaseMethod '(' argList? ')'
    | 'db' '.' 'getSiblingDB' '(' STRING ')' '.' IDENT '.' collectionMethod '(' argList? ')'
    | 'db' '.' 'aggregate' '(' argList? ')' '.' cursorMethod '(' argList? ')'
    | 'db' '.' 'watch' '(' argList? ')' '.' changeStreamMethod '(' argList? ')'
    | 'db' '.' 'getMongo' '(' ')' '.' mongoMethod '(' argList? ')'
    ;

databaseMethod
    : 'getMongo' // packages/shell-api/src/database.ts:364
    | 'getName' // packages/shell-api/src/database.ts:368
    | 'getCollectionNames' // packages/shell-api/src/database.ts:379
    | 'getCollectionInfos' // packages/shell-api/src/database.ts:395
    | 'runCommand' // packages/shell-api/src/database.ts:412
    | 'adminCommand' // packages/shell-api/src/database.ts:445
    | 'aggregate' // packages/shell-api/src/database.ts:466
    | 'getSiblingDB' // packages/shell-api/src/database.ts:538
    | 'getCollection' // packages/shell-api/src/database.ts:548
    | 'dropDatabase' // packages/shell-api/src/database.ts:576
    | 'createUser' // packages/shell-api/src/database.ts:585
    | 'updateUser' // packages/shell-api/src/database.ts:633
    | 'changeUserPassword' // packages/shell-api/src/database.ts:681
    | 'logout' // packages/shell-api/src/database.ts:710
    | 'dropUser' // packages/shell-api/src/database.ts:718
    | 'dropAllUsers' // packages/shell-api/src/database.ts:733
    | 'auth' // packages/shell-api/src/database.ts:743
    | 'grantRolesToUser' // packages/shell-api/src/database.ts:788
    | 'revokeRolesFromUser' // packages/shell-api/src/database.ts:808
    | 'getUser' // packages/shell-api/src/database.ts:828
    | 'getUsers' // packages/shell-api/src/database.ts:855
    | 'createCollection' // packages/shell-api/src/database.ts:863
    | 'createEncryptedCollection' // packages/shell-api/src/database.ts:881
    | 'createView' // packages/shell-api/src/database.ts:896
    | 'createRole' // packages/shell-api/src/database.ts:930
    | 'updateRole' // packages/shell-api/src/database.ts:953
    | 'dropRole' // packages/shell-api/src/database.ts:980
    | 'dropAllRoles' // packages/shell-api/src/database.ts:994
    | 'grantRolesToRole' // packages/shell-api/src/database.ts:1005
    | 'revokeRolesFromRole' // packages/shell-api/src/database.ts:1025
    | 'grantPrivilegesToRole' // packages/shell-api/src/database.ts:1044
    | 'revokePrivilegesFromRole' // packages/shell-api/src/database.ts:1067
    | 'getRole' // packages/shell-api/src/database.ts:1090
    | 'getRoles' // packages/shell-api/src/database.ts:1117
    | 'currentOp' // packages/shell-api/src/database.ts:1178
    | 'killOp' // packages/shell-api/src/database.ts:1189
    | 'shutdownServer' // packages/shell-api/src/database.ts:1199
    | 'fsyncLock' // packages/shell-api/src/database.ts:1209
    | 'fsyncUnlock' // packages/shell-api/src/database.ts:1219
    | 'version' // packages/shell-api/src/database.ts:1228
    | 'serverBits' // packages/shell-api/src/database.ts:1246
    | 'isMaster' // packages/shell-api/src/database.ts:1264
    | 'hello' // packages/shell-api/src/database.ts:1276
    | 'serverBuildInfo' // packages/shell-api/src/database.ts:1296
    | 'serverStatus' // packages/shell-api/src/database.ts:1305
    | 'stats' // packages/shell-api/src/database.ts:1315
    | 'hostInfo' // packages/shell-api/src/database.ts:1334
    | 'serverCmdLineOpts' // packages/shell-api/src/database.ts:1343
    | 'rotateCertificates' // packages/shell-api/src/database.ts:1353
    | 'printCollectionStats' // packages/shell-api/src/database.ts:1363
    | 'getProfilingStatus' // packages/shell-api/src/database.ts:1385
    | 'setProfilingLevel' // packages/shell-api/src/database.ts:1394
    | 'setLogLevel' // packages/shell-api/src/database.ts:1417
    | 'getLogComponents' // packages/shell-api/src/database.ts:1450
    | 'cloneDatabase' // packages/shell-api/src/database.ts:1465
    | 'cloneCollection' // packages/shell-api/src/database.ts:1472
    | 'copyDatabase' // packages/shell-api/src/database.ts:1479
    | 'commandHelp' // packages/shell-api/src/database.ts:1487
    | 'listCommands' // packages/shell-api/src/database.ts:1508
    | 'getLastErrorObj' // packages/shell-api/src/database.ts:1540
    | 'getLastError' // packages/shell-api/src/database.ts:1560
    | 'printShardingStatus' // packages/shell-api/src/database.ts:1576
    | 'printSecondaryReplicationInfo' // packages/shell-api/src/database.ts:1588
    | 'getReplicationInfo' // packages/shell-api/src/database.ts:1669
    | 'printReplicationInfo' // packages/shell-api/src/database.ts:1734
    | 'printSlaveReplicationInfo' // packages/shell-api/src/database.ts:1766
    | 'setSecondaryOk' // packages/shell-api/src/database.ts:1774
    | 'watch' // packages/shell-api/src/database.ts:1782
    | 'sql' // packages/shell-api/src/database.ts:1818
    | 'checkMetadataConsistency' // packages/shell-api/src/database.ts:1861
    ;

// packages/shell-api/src/collection.ts (Collection) and cursor classes
collectionRef
    : 'db' '.' IDENT
    | 'db' '[' STRING ']'
    | 'db' '.' 'getCollection' '(' STRING ')'
    ;

collectionStatement
    : collectionRef '.' collectionMethod '(' argList? ')' cursorChain?
    ;

cursorChain
    : '.' cursorMethod '(' argList? ')'
    ;

collectionMethod
    : 'aggregate' // packages/shell-api/src/collection.ts:189
    | 'bulkWrite' // packages/shell-api/src/collection.ts:278
    | 'count' // packages/shell-api/src/collection.ts:318
    | 'countDocuments' // packages/shell-api/src/collection.ts:345
    | 'deleteMany' // packages/shell-api/src/collection.ts:372
    | 'deleteOne' // packages/shell-api/src/collection.ts:406
    | 'distinct' // packages/shell-api/src/collection.ts:438
    | 'estimatedDocumentCount' // packages/shell-api/src/collection.ts:476
    | 'find' // packages/shell-api/src/collection.ts:501
    | 'findAndModify' // packages/shell-api/src/collection.ts:552
    | 'findOne' // packages/shell-api/src/collection.ts:614
    | 'renameCollection' // packages/shell-api/src/collection.ts:639
    | 'findOneAndDelete' // packages/shell-api/src/collection.ts:684
    | 'findOneAndReplace' // packages/shell-api/src/collection.ts:721
    | 'findOneAndUpdate' // packages/shell-api/src/collection.ts:767
    | 'insert' // packages/shell-api/src/collection.ts:811
    | 'insertMany' // packages/shell-api/src/collection.ts:855
    | 'insertOne' // packages/shell-api/src/collection.ts:891
    | 'isCapped' // packages/shell-api/src/collection.ts:915
    | 'remove' // packages/shell-api/src/collection.ts:947
    | 'replaceOne' // packages/shell-api/src/collection.ts:990
    | 'update' // packages/shell-api/src/collection.ts:1018
    | 'updateMany' // packages/shell-api/src/collection.ts:1075
    | 'updateOne' // packages/shell-api/src/collection.ts:1118
    | 'compactStructuredEncryptionData' // packages/shell-api/src/collection.ts:1152
    | 'convertToCapped' // packages/shell-api/src/collection.ts:1175
    | 'createIndexes' // packages/shell-api/src/collection.ts:1237
    | 'createIndex' // packages/shell-api/src/collection.ts:1263
    | 'ensureIndex' // packages/shell-api/src/collection.ts:1298
    | 'getIndexes' // packages/shell-api/src/collection.ts:1316
    | 'getIndexSpecs' // packages/shell-api/src/collection.ts:1334
    | 'getIndices' // packages/shell-api/src/collection.ts:1351
    | 'getIndexKeys' // packages/shell-api/src/collection.ts:1368
    | 'dropIndexes' // packages/shell-api/src/collection.ts:1387
    | 'dropIndex' // packages/shell-api/src/collection.ts:1433
    | 'totalIndexSize' // packages/shell-api/src/collection.ts:1468
    | 'reIndex' // packages/shell-api/src/collection.ts:1489
    | 'getDB' // packages/shell-api/src/collection.ts:1502
    | 'getMongo' // packages/shell-api/src/collection.ts:1513
    | 'dataSize' // packages/shell-api/src/collection.ts:1525
    | 'storageSize' // packages/shell-api/src/collection.ts:1538
    | 'totalSize' // packages/shell-api/src/collection.ts:1551
    | 'drop' // packages/shell-api/src/collection.ts:1564
    | 'exists' // packages/shell-api/src/collection.ts:1625
    | 'getFullName' // packages/shell-api/src/collection.ts:1638
    | 'getName' // packages/shell-api/src/collection.ts:1643
    | 'runCommand' // packages/shell-api/src/collection.ts:1650
    | 'explain' // packages/shell-api/src/collection.ts:1688
    | 'stats' // packages/shell-api/src/collection.ts:1971
    | 'latencyStats' // packages/shell-api/src/collection.ts:2055
    | 'initializeOrderedBulkOp' // packages/shell-api/src/collection.ts:2065
    | 'initializeUnorderedBulkOp' // packages/shell-api/src/collection.ts:2079
    | 'getPlanCache' // packages/shell-api/src/collection.ts:2092
    | 'mapReduce' // packages/shell-api/src/collection.ts:2102
    | 'validate' // packages/shell-api/src/collection.ts:2142
    | 'getShardVersion' // packages/shell-api/src/collection.ts:2156
    | 'getShardDistribution' // packages/shell-api/src/collection.ts:2219
    | 'getShardLocation' // packages/shell-api/src/collection.ts:2348
    | 'watch' // packages/shell-api/src/collection.ts:2393
    | 'hideIndex' // packages/shell-api/src/collection.ts:2446
    | 'unhideIndex' // packages/shell-api/src/collection.ts:2454
    | 'analyzeShardKey' // packages/shell-api/src/collection.ts:2463
    | 'configureQueryAnalyzer' // packages/shell-api/src/collection.ts:2480
    | 'checkMetadataConsistency' // packages/shell-api/src/collection.ts:2491
    | 'getSearchIndexes' // packages/shell-api/src/collection.ts:2509
    | 'createSearchIndex' // packages/shell-api/src/collection.ts:2527
    | 'createSearchIndexes' // packages/shell-api/src/collection.ts:2596
    | 'dropSearchIndex' // packages/shell-api/src/collection.ts:2614
    | 'updateSearchIndex' // packages/shell-api/src/collection.ts:2627
    ;

// packages/shell-api/src/{abstract-cursor,aggregate-or-find-cursor,cursor,change-stream-cursor,explainable-cursor}.ts
cursorMethod
    : 'close' // packages/shell-api/src/abstract-cursor.ts:84
    | 'forEach' // packages/shell-api/src/abstract-cursor.ts:89
    | 'hasNext' // packages/shell-api/src/abstract-cursor.ts:101
    | 'tryNext' // packages/shell-api/src/abstract-cursor.ts:106
    | 'isClosed' // packages/shell-api/src/abstract-cursor.ts:141
    | 'isExhausted' // packages/shell-api/src/abstract-cursor.ts:145
    | 'itcount' // packages/shell-api/src/abstract-cursor.ts:150
    | 'pretty' // packages/shell-api/src/abstract-cursor.ts:160
    | 'map' // packages/shell-api/src/abstract-cursor.ts:166
    | 'next' // packages/shell-api/src/abstract-cursor.ts:177
    | 'disableBlockWarnings' // packages/shell-api/src/abstract-cursor.ts:186
    | 'maxTimeMS' // packages/shell-api/src/abstract-cursor.ts:259
    | 'objsLeftInBatch' // packages/shell-api/src/abstract-cursor.ts:264
    | 'projection' // packages/shell-api/src/aggregate-or-find-cursor.ts:25
    | 'skip' // packages/shell-api/src/aggregate-or-find-cursor.ts:32
    | 'sort' // packages/shell-api/src/aggregate-or-find-cursor.ts:39
    | 'explain' // packages/shell-api/src/aggregate-or-find-cursor.ts:46
    | 'toJSON' // packages/shell-api/src/cursor.ts:51
    | 'addOption' // packages/shell-api/src/cursor.ts:72
    | 'allowDiskUse' // packages/shell-api/src/cursor.ts:97
    | 'allowPartialResults' // packages/shell-api/src/cursor.ts:104
    | 'collation' // packages/shell-api/src/cursor.ts:112
    | 'comment' // packages/shell-api/src/cursor.ts:120
    | 'count' // packages/shell-api/src/cursor.ts:128
    | 'hint' // packages/shell-api/src/cursor.ts:146
    | 'limit' // packages/shell-api/src/cursor.ts:153
    | 'max' // packages/shell-api/src/cursor.ts:160
    | 'maxAwaitTimeMS' // packages/shell-api/src/cursor.ts:168
    | 'min' // packages/shell-api/src/cursor.ts:175
    | 'noCursorTimeout' // packages/shell-api/src/cursor.ts:194
    | 'oplogReplay' // packages/shell-api/src/cursor.ts:201
    | 'readPref' // packages/shell-api/src/cursor.ts:208
    | 'returnKey' // packages/shell-api/src/cursor.ts:232
    | 'size' // packages/shell-api/src/cursor.ts:238
    | 'tailable' // packages/shell-api/src/cursor.ts:246
    | 'maxScan' // packages/shell-api/src/cursor.ts:257
    | 'showRecordId' // packages/shell-api/src/cursor.ts:266
    | 'readConcern' // packages/shell-api/src/cursor.ts:273
    | 'getResumeToken' // packages/shell-api/src/change-stream-cursor.ts:80
    | 'finish' // packages/shell-api/src/explainable-cursor.ts:28
    ;

// packages/shell-api/src/change-stream-cursor.ts (ChangeStreamCursor)
changeStreamMethod
    : 'getResumeToken' // packages/shell-api/src/change-stream-cursor.ts:80
    ;

// packages/shell-api/src/explainable.ts (Explainable)
explainableStatement
    : collectionRef '.' 'explain' '(' argList? ')' ( '.' explainableMethod '(' argList? ')' )?
    ;

explainableMethod
    : 'getCollection' // packages/shell-api/src/explainable.ts:81
    | 'getVerbosity' // packages/shell-api/src/explainable.ts:86
    | 'setVerbosity' // packages/shell-api/src/explainable.ts:91
    | 'find' // packages/shell-api/src/explainable.ts:100
    | 'aggregate' // packages/shell-api/src/explainable.ts:111
    | 'count' // packages/shell-api/src/explainable.ts:135
    | 'distinct' // packages/shell-api/src/explainable.ts:150
    | 'findAndModify' // packages/shell-api/src/explainable.ts:173
    | 'findOneAndDelete' // packages/shell-api/src/explainable.ts:185
    | 'findOneAndReplace' // packages/shell-api/src/explainable.ts:198
    | 'findOneAndUpdate' // packages/shell-api/src/explainable.ts:212
    | 'remove' // packages/shell-api/src/explainable.ts:226
    | 'update' // packages/shell-api/src/explainable.ts:237
    | 'mapReduce' // packages/shell-api/src/explainable.ts:252
    ;

// packages/shell-api/src/bulk.ts (Bulk, BulkFindOp)
bulkStatement
    : collectionRef '.' ( 'initializeOrderedBulkOp' | 'initializeUnorderedBulkOp' ) '(' ')' '.' bulkMethod '(' argList? ')'
    ;

bulkMethod
    : 'execute' // packages/shell-api/src/bulk.ts:188
    | 'find' // packages/shell-api/src/bulk.ts:206
    | 'insert' // packages/shell-api/src/bulk.ts:213
    | 'toJSON' // packages/shell-api/src/bulk.ts:220
    | 'toString' // packages/shell-api/src/bulk.ts:232
    | 'getOperations' // packages/shell-api/src/bulk.ts:236
    | 'collation' // packages/shell-api/src/bulk.ts:46
    | 'arrayFilters' // packages/shell-api/src/bulk.ts:53
    | 'hint' // packages/shell-api/src/bulk.ts:60
    | 'delete' // packages/shell-api/src/bulk.ts:68
    | 'deleteOne' // packages/shell-api/src/bulk.ts:76
    | 'remove' // packages/shell-api/src/bulk.ts:85
    | 'removeOne' // packages/shell-api/src/bulk.ts:92
    | 'replaceOne' // packages/shell-api/src/bulk.ts:98
    | 'updateOne' // packages/shell-api/src/bulk.ts:108
    | 'update' // packages/shell-api/src/bulk.ts:117
    | 'upsert' // packages/shell-api/src/bulk.ts:126
    ;

// packages/shell-api/src/plan-cache.ts (PlanCache)
planCacheStatement
    : collectionRef '.' 'getPlanCache' '(' ')' ( '.' planCacheMethod '(' argList? ')' )?
    ;

planCacheMethod
    : 'clear' // packages/shell-api/src/plan-cache.ts:38
    | 'clearPlansByQuery' // packages/shell-api/src/plan-cache.ts:44
    | 'list' // packages/shell-api/src/plan-cache.ts:62
    | 'listQueryShapes' // packages/shell-api/src/plan-cache.ts:72
    | 'getPlansByQuery' // packages/shell-api/src/plan-cache.ts:79
    ;

// packages/shell-api/src/replica-set.ts (ReplicaSet)
rsStatement
    : 'rs' '.' rsMethod '(' argList? ')'
    ;

rsMethod
    : 'initiate' // packages/shell-api/src/replica-set.ts:60
    | 'config' // packages/shell-api/src/replica-set.ts:105
    | 'conf' // packages/shell-api/src/replica-set.ts:115
    | 'reconfig' // packages/shell-api/src/replica-set.ts:128
    | 'reconfigForPSASet' // packages/shell-api/src/replica-set.ts:152
    | 'status' // packages/shell-api/src/replica-set.ts:247
    | 'isMaster' // packages/shell-api/src/replica-set.ts:256
    | 'hello' // packages/shell-api/src/replica-set.ts:263
    | 'printSecondaryReplicationInfo' // packages/shell-api/src/replica-set.ts:270
    | 'printSlaveReplicationInfo' // packages/shell-api/src/replica-set.ts:277
    | 'printReplicationInfo' // packages/shell-api/src/replica-set.ts:285
    | 'add' // packages/shell-api/src/replica-set.ts:292
    | 'addArb' // packages/shell-api/src/replica-set.ts:339
    | 'remove' // packages/shell-api/src/replica-set.ts:346
    | 'freeze' // packages/shell-api/src/replica-set.ts:370
    | 'stepDown' // packages/shell-api/src/replica-set.ts:380
    | 'syncFrom' // packages/shell-api/src/replica-set.ts:404
    | 'secondaryOk' // packages/shell-api/src/replica-set.ts:414
    ;

// packages/shell-api/src/shard.ts (Shard)
shStatement
    : 'sh' '.' shMethod '(' argList? ')'
    ;

shMethod
    : 'enableSharding' // packages/shell-api/src/shard.ts:81
    | 'commitReshardCollection' // packages/shell-api/src/shard.ts:111
    | 'abortReshardCollection' // packages/shell-api/src/shard.ts:126
    | 'shardCollection' // packages/shell-api/src/shard.ts:140
    | 'reshardCollection' // packages/shell-api/src/shard.ts:158
    | 'status' // packages/shell-api/src/shard.ts:215
    | 'addShard' // packages/shell-api/src/shard.ts:228
    | 'addShardToZone' // packages/shell-api/src/shard.ts:240
    | 'addShardTag' // packages/shell-api/src/shard.ts:257
    | 'updateZoneKeyRange' // packages/shell-api/src/shard.ts:276
    | 'addTagRange' // packages/shell-api/src/shard.ts:301
    | 'removeRangeFromZone' // packages/shell-api/src/shard.ts:327
    | 'removeTagRange' // packages/shell-api/src/shard.ts:344
    | 'removeShardFromZone' // packages/shell-api/src/shard.ts:368
    | 'removeShardTag' // packages/shell-api/src/shard.ts:386
    | 'enableAutoSplit' // packages/shell-api/src/shard.ts:406
    | 'disableAutoSplit' // packages/shell-api/src/shard.ts:430
    | 'splitAt' // packages/shell-api/src/shard.ts:453
    | 'splitFind' // packages/shell-api/src/shard.ts:464
    | 'moveChunk' // packages/shell-api/src/shard.ts:475
    | 'moveRange' // packages/shell-api/src/shard.ts:495
    | 'balancerCollectionStatus' // packages/shell-api/src/shard.ts:521
    | 'enableBalancing' // packages/shell-api/src/shard.ts:531
    | 'disableBalancing' // packages/shell-api/src/shard.ts:546
    | 'enableMigrations' // packages/shell-api/src/shard.ts:582
    | 'disableMigrations' // packages/shell-api/src/shard.ts:588
    | 'getBalancerState' // packages/shell-api/src/shard.ts:594
    | 'isBalancerRunning' // packages/shell-api/src/shard.ts:608
    | 'startBalancer' // packages/shell-api/src/shard.ts:618
    | 'stopBalancer' // packages/shell-api/src/shard.ts:629
    | 'setBalancerState' // packages/shell-api/src/shard.ts:640
    | 'getShardedDataDistribution' // packages/shell-api/src/shard.ts:653
    | 'startAutoMerger' // packages/shell-api/src/shard.ts:677
    | 'stopAutoMerger' // packages/shell-api/src/shard.ts:692
    | 'isAutoMergerEnabled' // packages/shell-api/src/shard.ts:707
    | 'disableAutoMerger' // packages/shell-api/src/shard.ts:722
    | 'enableAutoMerger' // packages/shell-api/src/shard.ts:739
    | 'checkMetadataConsistency' // packages/shell-api/src/shard.ts:755
    | 'shardAndDistributeCollection' // packages/shell-api/src/shard.ts:767
    | 'moveCollection' // packages/shell-api/src/shard.ts:803
    | 'abortMoveCollection' // packages/shell-api/src/shard.ts:819
    | 'unshardCollection' // packages/shell-api/src/shard.ts:828
    | 'abortUnshardCollection' // packages/shell-api/src/shard.ts:847
    | 'listShards' // packages/shell-api/src/shard.ts:859
    | 'isConfigShardEnabled' // packages/shell-api/src/shard.ts:879
    | 'shardDrainingStatus' // packages/shell-api/src/shard.ts:900
    | 'getTransitionToDedicatedConfigServerStatus' // packages/shell-api/src/shard.ts:917
    ;

// packages/shell-api/src/streams.ts (Streams) and packages/shell-api/src/stream-processor.ts (StreamProcessor)
spStatement
    : 'sp' '.' streamsMethod '(' argList? ')'
    | 'sp' '.' IDENT '.' streamProcessorMethod '(' argList? ')'
    ;

streamsMethod
    : 'getProcessor' // packages/shell-api/src/streams.ts:59
    | 'process' // packages/shell-api/src/streams.ts:64
    | 'createStreamProcessor' // packages/shell-api/src/streams.ts:106
    | 'listStreamProcessors' // packages/shell-api/src/streams.ts:137
    | 'listConnections' // packages/shell-api/src/streams.ts:163
    | 'listWorkspaceDefaults' // packages/shell-api/src/streams.ts:177
    ;

streamProcessorMethod
    : 'start' // packages/shell-api/src/stream-processor.ts:61
    | 'stop' // packages/shell-api/src/stream-processor.ts:69
    | 'drop' // packages/shell-api/src/stream-processor.ts:77
    | 'stats' // packages/shell-api/src/stream-processor.ts:89
    | 'modify' // packages/shell-api/src/stream-processor.ts:105
    | 'sample' // packages/shell-api/src/stream-processor.ts:137
    ;

// packages/shell-api/src/shell-api.ts (Mongo, connect) and packages/shell-api/src/mongo.ts (Mongo)
mongoConnectionStatement
    : 'Mongo' '(' argList? ')' ( '.' mongoMethod '(' argList? ')' )?
    | 'connect' '(' argList? ')' ( '.' databaseMethod '(' argList? ')' )?
    ;

mongoMethod
    : 'connect' // packages/shell-api/src/mongo.ts:221
    | 'getDB' // packages/shell-api/src/mongo.ts:284
    | 'getCollection' // packages/shell-api/src/mongo.ts:291
    | 'getURI' // packages/shell-api/src/mongo.ts:311
    | 'use' // packages/shell-api/src/mongo.ts:315
    | 'getDBs' // packages/shell-api/src/mongo.ts:383
    | 'bulkWrite' // packages/shell-api/src/mongo.ts:398
    | 'getDBNames' // packages/shell-api/src/mongo.ts:431
    | 'show' // packages/shell-api/src/mongo.ts:442
    | 'close' // packages/shell-api/src/mongo.ts:589
    | 'getReadPrefMode' // packages/shell-api/src/mongo.ts:608
    | 'getReadPrefTagSet' // packages/shell-api/src/mongo.ts:612
    | 'getReadPref' // packages/shell-api/src/mongo.ts:616
    | 'getReadConcern' // packages/shell-api/src/mongo.ts:628
    | 'getWriteConcern' // packages/shell-api/src/mongo.ts:637
    | 'setReadPref' // packages/shell-api/src/mongo.ts:646
    | 'setReadConcern' // packages/shell-api/src/mongo.ts:662
    | 'setWriteConcern' // packages/shell-api/src/mongo.ts:668
    | 'startSession' // packages/shell-api/src/mongo.ts:752
    | 'setCausalConsistency' // packages/shell-api/src/mongo.ts:804
    | 'isCausalConsistency' // packages/shell-api/src/mongo.ts:812
    | 'setSlaveOk' // packages/shell-api/src/mongo.ts:821
    | 'setSecondaryOk' // packages/shell-api/src/mongo.ts:829
    | 'watch' // packages/shell-api/src/mongo.ts:851
    | 'getClientEncryption' // packages/shell-api/src/mongo.ts:879
    | 'getKeyVault' // packages/shell-api/src/mongo.ts:896
    | 'convertShardKeyToHashed' // packages/shell-api/src/mongo.ts:905
    ;

// packages/shell-api/src/field-level-encryption.ts (KeyVault, ClientEncryption)
encryptionStatement
    : 'db' '.' 'getMongo' '(' ')' '.' 'getKeyVault' '(' ')' ( '.' keyVaultMethod '(' argList? ')' )?
    | 'db' '.' 'getMongo' '(' ')' '.' 'getClientEncryption' '(' ')' ( '.' clientEncryptionMethod '(' argList? ')' )?
    ;

keyVaultMethod
    : 'createKey' // packages/shell-api/src/field-level-encryption.ts:286
    | 'getKey' // packages/shell-api/src/field-level-encryption.ts:378
    | 'getKeyByAltName' // packages/shell-api/src/field-level-encryption.ts:386
    | 'getKeys' // packages/shell-api/src/field-level-encryption.ts:396
    | 'deleteKey' // packages/shell-api/src/field-level-encryption.ts:405
    | 'addKeyAlternateName' // packages/shell-api/src/field-level-encryption.ts:412
    | 'removeKeyAlternateName' // packages/shell-api/src/field-level-encryption.ts:429
    | 'rewrapManyDataKey' // packages/shell-api/src/field-level-encryption.ts:446
    | 'createDataKey' // packages/shell-api/src/field-level-encryption.ts:459
    | 'removeKeyAltName' // packages/shell-api/src/field-level-encryption.ts:468
    | 'addKeyAltName' // packages/shell-api/src/field-level-encryption.ts:477
    ;

clientEncryptionMethod
    : 'encrypt' // packages/shell-api/src/field-level-encryption.ts:120
    | 'decrypt' // packages/shell-api/src/field-level-encryption.ts:148
    | 'encryptExpression' // packages/shell-api/src/field-level-encryption.ts:154
    | 'createEncryptedCollection' // packages/shell-api/src/field-level-encryption.ts:172
    ;

// packages/shell-api/src/session.ts (Session)
sessionStatement
    : 'db' '.' 'getMongo' '(' ')' '.' 'startSession' '(' argList? ')' ( '.' sessionMethod '(' argList? ')' )?
    ;

sessionMethod
    : 'getDatabase' // packages/shell-api/src/session.ts:60
    | 'advanceOperationTime' // packages/shell-api/src/session.ts:80
    | 'advanceClusterTime' // packages/shell-api/src/session.ts:84
    | 'endSession' // packages/shell-api/src/session.ts:89
    | 'hasEnded' // packages/shell-api/src/session.ts:93
    | 'getClusterTime' // packages/shell-api/src/session.ts:97
    | 'getOperationTime' // packages/shell-api/src/session.ts:101
    | 'getOptions' // packages/shell-api/src/session.ts:105
    | 'startTransaction' // packages/shell-api/src/session.ts:109
    | 'commitTransaction' // packages/shell-api/src/session.ts:114
    | 'abortTransaction' // packages/shell-api/src/session.ts:119
    | 'withTransaction' // packages/shell-api/src/session.ts:124
    | 'bulkWrite' // packages/shell-api/src/session.ts:137
    ;

// packages/shell-bson/src/shell-bson.ts (EJSON)
ejsonStatement
    : 'EJSON' '.' ejsonMethod '(' argList? ')'
    ;

ejsonMethod
    : 'parse' // packages/shell-bson/src/shell-bson.ts:523
    | 'serialize' // packages/shell-bson/src/shell-bson.ts:524
    | 'stringify' // packages/shell-bson/src/shell-bson.ts:525
    | 'deserialize' // packages/shell-bson/src/shell-bson.ts:526
    ;

// Representative argument, document, array, and literal value shapes
argList
    : value ( ',' value )* ','?
    ;

value
    : document
    | array
    | STRING
    | NUMBER
    | NEG_NUMBER
    | 'true'
    | 'false'
    | 'null'
    | REGEX
    | bsonHelperCall
    | bsonStaticCall
    | IDENT
    ;

document
    : '{' '}'
    | '{' pair ( ',' pair )* ','? '}'
    ;

pair
    : docKey ':' value
    ;

docKey
    : IDENT
    | STRING
    | FIELD_OP
    ;

array
    : '[' ']'
    | '[' value ( ',' value )* ','? ']'
    ;

// packages/shell-bson/src/shell-bson.ts (ShellBsonBase)
bsonHelperCall
    : bsonHelperName '(' argList? ')'
    ;

bsonHelperName
    : 'DBRef' // packages/shell-bson/src/shell-bson.ts:21
    | 'bsonsize' // packages/shell-bson/src/shell-bson.ts:28
    | 'MaxKey' // packages/shell-bson/src/shell-bson.ts:29
    | 'MinKey' // packages/shell-bson/src/shell-bson.ts:32
    | 'ObjectId' // packages/shell-bson/src/shell-bson.ts:35
    | 'Timestamp' // packages/shell-bson/src/shell-bson.ts:39
    | 'Code' // packages/shell-bson/src/shell-bson.ts:44
    | 'NumberDecimal' // packages/shell-bson/src/shell-bson.ts:46
    | 'NumberInt' // packages/shell-bson/src/shell-bson.ts:47
    | 'NumberLong' // packages/shell-bson/src/shell-bson.ts:48
    | 'ISODate' // packages/shell-bson/src/shell-bson.ts:49
    | 'BinData' // packages/shell-bson/src/shell-bson.ts:50
    | 'HexData' // packages/shell-bson/src/shell-bson.ts:54
    | 'UUID' // packages/shell-bson/src/shell-bson.ts:55
    | 'LegacyJavaUUID' // packages/shell-bson/src/shell-bson.ts:56
    | 'LegacyCSharpUUID' // packages/shell-bson/src/shell-bson.ts:57
    | 'LegacyPythonUUID' // packages/shell-bson/src/shell-bson.ts:58
    | 'MD5' // packages/shell-bson/src/shell-bson.ts:59
    | 'Decimal128' // packages/shell-bson/src/shell-bson.ts:60
    | 'BSONSymbol' // packages/shell-bson/src/shell-bson.ts:61
    | 'Int32' // packages/shell-bson/src/shell-bson.ts:62
    | 'Long' // packages/shell-bson/src/shell-bson.ts:63
    | 'Binary' // packages/shell-bson/src/shell-bson.ts:64
    | 'Double' // packages/shell-bson/src/shell-bson.ts:65
    | 'BSONRegExp' // packages/shell-bson/src/shell-bson.ts:67
    ;

// packages/shell-bson/src/shell-bson.ts (static BSON helpers)
bsonStaticCall
    : bsonStaticTarget '(' argList? ')'
    ;

bsonStaticTarget
    : 'ObjectId' '.' 'generate' // packages/shell-bson/src/shell-bson.ts:191
    | 'ObjectId' '.' 'createFromTime' // packages/shell-bson/src/shell-bson.ts:191
    | 'ObjectId' '.' 'createFromHexString' // packages/shell-bson/src/shell-bson.ts:191
    | 'ObjectId' '.' 'createFromBase64' // packages/shell-bson/src/shell-bson.ts:191
    | 'ObjectId' '.' 'isValid' // packages/shell-bson/src/shell-bson.ts:191
    | 'Timestamp' '.' 'fromInt' // packages/shell-bson/src/shell-bson.ts:214
    | 'Timestamp' '.' 'fromNumber' // packages/shell-bson/src/shell-bson.ts:214
    | 'Timestamp' '.' 'fromBits' // packages/shell-bson/src/shell-bson.ts:214
    | 'Timestamp' '.' 'fromString' // packages/shell-bson/src/shell-bson.ts:214
    | 'Decimal128' '.' 'fromString' // packages/shell-bson/src/shell-bson.ts:447
    | 'Decimal128' '.' 'fromStringWithRounding' // packages/shell-bson/src/shell-bson.ts:448
    | 'Int32' '.' 'fromString' // packages/shell-bson/src/shell-bson.ts:457
    | 'Binary' '.' 'createFromBase64' // packages/shell-bson/src/shell-bson.ts:492
    | 'Binary' '.' 'createFromHexString' // packages/shell-bson/src/shell-bson.ts:493
    | 'Binary' '.' 'fromInt8Array' // packages/shell-bson/src/shell-bson.ts:494
    | 'Binary' '.' 'fromFloat32Array' // packages/shell-bson/src/shell-bson.ts:495
    | 'Binary' '.' 'fromPackedBits' // packages/shell-bson/src/shell-bson.ts:496
    | 'Binary' '.' 'fromBits' // packages/shell-bson/src/shell-bson.ts:497
    | 'Double' '.' 'fromString' // packages/shell-bson/src/shell-bson.ts:515
    | 'BSONRegExp' '.' 'parseOptions' // packages/shell-bson/src/shell-bson.ts:519
    | 'EJSON' '.' 'parse' // packages/shell-bson/src/shell-bson.ts:523
    | 'EJSON' '.' 'serialize' // packages/shell-bson/src/shell-bson.ts:524
    | 'EJSON' '.' 'stringify' // packages/shell-bson/src/shell-bson.ts:525
    | 'EJSON' '.' 'deserialize' // packages/shell-bson/src/shell-bson.ts:526
    ;

