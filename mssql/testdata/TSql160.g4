// MIT License
//
// Copyright (c) Microsoft Corporation.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// This file is derived from Microsoft ScriptDOM's TSql160.g, TSqlTokenTypes.g
// and CodeGenerationSupporter.cs, which are beside it as published, with their
// license (ScriptDom.LICENSE). mssql/parser/cmd/scriptdomg4 wrote it from
// them, converting ANTLR 2 to ANTLR 4 and keeping only what each rule derives:
// actions, arguments, labels and predicates are gone, a word an action matches
// an identifier against is written as that word, and what else was lost is
// noted beside its rule. It generates test statements; it is not a parser.

grammar TSql160;

options { caseInsensitive = true; }

// TSql160.g:118
script
    : batch ( Go batch )* EOF
    ;

// TSql160.g:158
batch
    : ( Semicolon )* ( ( lastStatementOptSemi ) | ( optSimpleExecute ( statementOptSemi )* ) )
    ;

// TSql160.g:252
statementOptSemi
    : statement optSemicolons
    ;

// TSql160.g:256
lastStatementOptSemi
    : lastStatement optSemicolons
    ;

// TSql160.g:260
optSemicolons
    : ( Semicolon )*
    ;

// TSql160.g:280
copyStatement
    : 'COPY' Into schemaObjectThreePartName ( copyColumnList )? From nonEmptyString ( Comma nonEmptyString )* ( copyWithClause )? ( optimizerHints )?
    ;

// TSql160.g:312
copyColumnList
    : LeftParenthesis copyColumnOption ( Comma copyColumnOption )* RightParenthesis
    ;

// TSql160.g:341
copyColumnOption
    : identifier ( Default defaultValueLiteral )? ( integer )?
    ;

// TSql160.g:363
copyWithClause
    : With LeftParenthesis copyOption ( Comma copyOption )* RightParenthesis
    ;

// TSql160.g:383
// TSql160.g:413: action may refuse the input: {if ((vResult.Kind != CopyOptionKind.Credential || vResult.Kind != CopyOptionKind.ErrorFileCredential) && CopyIdentifierOrValueOptionsHelper.Instance.ValidateCopyCredential((CopyCredentialOption)vValue)) { vResult.Value = (CopyCredentialOption) vValue; ...
copyOption
    : ( Identifier | IdentityInsert | Credential ) EqualsSign ( singleValueTypeCopyOption | copyCredentialOption )
    ;

// TSql160.g:427
singleValueTypeCopyOption
    : identifier
    | integer
    | stringLiteral
    ;

// TSql160.g:447
copyCredentialOption
    : LeftParenthesis Identity EqualsSign stringLiteral ( Comma 'SECRET' EqualsSign stringLiteral )? RightParenthesis
    ;

// TSql160.g:523
ctasCreateTableStatement
    : ( columnNameList )? ( withTableOptions ) As selectStatement
    ;

// TSql160.g:546
ctasCreateExternalTableStatement
    : ( withExternalTableOptions ) As selectStatement
    ;

// TSql160.g:567
// TSql160.g:572: predicate dropped: {!NextTokenMatches(CodeGenerationSupporter.Disable) && !NextTokenMatches(CodeGenerationSupporter.Enable) && !NextTokenMatches(CodeGenerationSupporter.Move) && !NextTokenMatches(CodeGenerationSupporter.Get) && !NextTokenMatches(CodeGenerationSupporter.Receive) && ...
optSimpleExecute
    : ( execProc optSemicolons )
    | /* empty */
    ;

// TSql160.g:588
statement
    : createTableStatement
    | alterTableStatement
    | createIndexStatement
    | copyStatement
    | declareStatements
    | setStatements
    | beginStatements
    | breakStatement
    | continueStatement
    | ifStatement
    | whileStatement
    | labelStatement
    | backupStatements
    | restoreStatements
    | gotoStatement
    | saveTransactionStatement
    | rollbackTransactionStatement
    | commitTransactionStatement
    | createStatisticsStatement
    | updateStatisticsStatement
    | alterDatabaseStatements
    | executeStatement
    | withCommonTableExpressionsAndXmlNamespacesStatements
    | raiseErrorStatement
    | alter2005Statements
    | create2005Statements
    | createDatabaseStatements
    | addStatements
    | identifierStatements
    | printStatement
    | waitForStatement
    | readTextStatement
    | updateTextStatement
    | writeTextStatement
    | lineNoStatement
    | useStatement
    | killStatements
    | bulkInsertStatement
    | insertBulkStatement
    | checkpointStatement
    | reconfigureStatement
    | shutdownStatement
    | setUserStatement
    | truncateTableStatement
    | grantStatement90
    | denyStatement90
    | revokeStatement90
    | returnStatement
    | openStatements
    | closeStatements
    | deallocateCursorStatement
    | fetchCursorStatement
    | dropStatements
    | dbccStatement
    | revertStatement
    | executeAsStatement
    | endConversationStatement
    ;

// TSql160.g:699
// TSql160.g:709: action may refuse the input: {// check for invalid combination of CHANGE_TRACKING_CONTEXT and Select statement if ((vWithCommonTableExpressionsAndXmlNamespaces != null) && (vWithCommonTableExpressionsAndXmlNamespaces.ChangeTrackingContext != null)) ThrowParseErrorException("SQL4607...
withCommonTableExpressionsAndXmlNamespacesStatements
    : ( withCommonTableExpressionsAndXmlNamespaces )? ( select | deleteStatement | insertStatement | updateStatement | mergeStatement )
    ;

// TSql160.g:728
lastStatement
    : createProcedureStatement
    | alterProcedureStatement
    | createTriggerStatement
    | alterTriggerStatement
    | createDefaultStatement
    | createRuleStatement
    | createViewStatement
    | alterViewStatement
    | createFunctionStatement
    | alterFunctionStatement
    | createSchemaStatement
    | createIdentifierStatement
    | alterIdentifierStatement
    | useFederationStatement
    | createOrAlterStatements
    ;

// TSql160.g:746
// TSql160.g:749: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Materialized)}?
// TSql160.g:752: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Federation)}?
createIdentifierStatement
    : Create ( createMaterializedViewStatement_MATERIALIZED | createFederationStatement_FEDERATION )
    ;

// TSql160.g:760
// TSql160.g:763: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Materialized)}?
// TSql160.g:766: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Federation)}?
alterIdentifierStatement
    : Alter ( alterMaterializedViewStatement_MATERIALIZED | alterFederationStatement_FEDERATION )
    ;

// TSql160.g:774
createOrAlterStatements
    : Create Or Alter ( createOrAlterFunctionStatement | createOrAlterProcedureStatement | createOrAlterTriggerStatement | createOrAlterViewStatement )
    ;

// TSql160.g:789
// TSql160.g:790: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Disable)}?
// TSql160.g:792: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Enable)}?
// TSql160.g:794: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Move)}?
// TSql160.g:796: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Get)}?
// TSql160.g:798: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Receive)}?
// TSql160.g:800: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Send)}?
// TSql160.g:802: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Throw)}?
// TSql160.g:804: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Rename)}?
identifierStatements
    : disableTriggerStatement_DISABLE
    | enableTriggerStatement_ENABLE
    | moveConversationStatement_MOVE
    | getConversationGroupStatement_GET
    | receiveStatement_RECEIVE
    | sendStatement_SEND
    | throwStatement_THROW
    | renameEntityStatement_RENAME
    ;

// TSql160.g:831
enableDisableTriggerBody
    : Trigger ( schemaObjectThreePartName ( Comma schemaObjectThreePartName )* | All ) On triggerObject
    ;

// TSql160.g:861
// TSql160.g:864: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Aggregate)}?
// TSql160.g:867: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Application)}?
// TSql160.g:870: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Assembly)}?
// TSql160.g:873: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Asymmetric)}?
// TSql160.g:876: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Availability)}?
// TSql160.g:879: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Broker)}?
// TSql160.g:882: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Certificate)}?
// TSql160.g:885: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.Column)}?
// TSql160.g:888: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.ColumnStore)}?
// TSql160.g:891: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Contract)}?
// TSql160.g:894: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Credential)}?
// TSql160.g:897: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Cryptographic)}?
// TSql160.g:900: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Endpoint)}?
// TSql160.g:903: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Event)}?
// TSql160.g:906: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.External)}?
// TSql160.g:909: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Fulltext)}?
// TSql160.g:914: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Selective)}?
// TSql160.g:917: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Xml)}?
// TSql160.g:920: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Login)}?
// TSql160.g:923: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Message)}?
// TSql160.g:926: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Master)}?
// TSql160.g:929: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Partition)}?
// TSql160.g:932: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Queue)}?
// TSql160.g:935: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Remote)}?
// TSql160.g:938: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Resource)}?
// TSql160.g:941: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Role)}?
// TSql160.g:944: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Route)}?
// TSql160.g:947: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Search)}?
// TSql160.g:950: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Service)}?
// TSql160.g:953: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Spatial)}?
// TSql160.g:956: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Symmetric)}?
// TSql160.g:959: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Synonym)}?
// TSql160.g:962: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Type)}?
// TSql160.g:965: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Server)}?
// TSql160.g:968: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Workload)}?
// TSql160.g:971: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Sequence)}?
// TSql160.g:974: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Security)}?
create2005Statements
    : Create ( createAggregateStatement_AGGREGATE | createApplicationRoleStatement_APPLICATION | createAssemblyStatement_ASSEMBLY | createAsymmetricKeyStatement_ASYMMETRIC | createAvailabilityGroupStatement_AVAILABILITY | createBrokerPriorityStatement_BROKER | createCertificateStatement_CERTIFICATE | createColumnStatements | createColumnStoreIndexStatement_COLUMNSTORE | createContractStatement_CONTRACT | createCredentialStatement_CREDENTIAL | createCryptographicProviderStatement_CRYPTOGRAPHIC | createEndpointStatement_ENDPOINT | createEventStatement_EVENT | createExternalStatements | createFulltextStatement_FULLTEXT | createPrimaryXmlIndexStatement | createSelectiveXmlIndexStatement_SELECTIVE | createXmlStatements_XML | createLoginStatement_LOGIN | createMessageTypeStatement_MESSAGE | createMasterKeyStatement_MASTER | createPartitionStatement_PARTITION | createQueueStatement_QUEUE | createRemoteServiceBindingStatement_REMOTE | createResourcePoolStatement_RESOURCE | createRoleStatement_ROLE | createRouteStatement_ROUTE | createSearchPropertyListStatement_SEARCH | createServiceStatement_SERVICE | createSpatialIndexStatement_SPATIAL | createSymmetricKeyStatement_SYMMETRIC | createSynonymStatement_SYNONYM | createTypeStatement_TYPE | createServerStatements_SERVER | createWorkloadStatements_WORKLOAD | createSequenceStatement_SEQUENCE | createSecurityPolicyStatement_SECURITY | createUserStatement )
    ;

// TSql160.g:1026
aggregateParameter
    : identifierVariable ( As )? scalarDataType ( nullableConstraint )?
    ;

// TSql160.g:1142
binaryOrString
    : binary
    | stringLiteral
    ;

// TSql160.g:1149
createExternalLibraryFileOption
    : LeftParenthesis 'CONTENT' EqualsSign binaryOrString ( Comma 'PLATFORM' EqualsSign identifier )* RightParenthesis
    ;

// TSql160.g:1168
alterExternalLibraryFileOption
    : LeftParenthesis 'CONTENT' EqualsSign binaryOrString ( Comma 'PLATFORM' EqualsSign identifier )* RightParenthesis
    ;

// TSql160.g:1210
// TSql160.g:1231: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Platform)}?
// TSql160.g:1240: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Parameters)}?
// TSql160.g:1249: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.EnvironmentVariables)}?
externalLanguageFileOption
    : LeftParenthesis 'CONTENT' EqualsSign binaryOrString Comma 'FILE_NAME' EqualsSign nonEmptyString ( Comma ( ( 'PLATFORM' EqualsSign identifier ) | ( 'PARAMETERS' EqualsSign nonEmptyString ) | ( 'ENVIRONMENT_VARIABLES' EqualsSign nonEmptyString ) ) )* RightParenthesis
    ;

// TSql160.g:1347
createAsymmetricKeyParams
    : From asymKeySource
    | With asymKeySpec
    ;

// TSql160.g:1358
// TSql160.g:1361: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Assembly)}?
asymKeySource
    : fileEncryptionSource
    | assemblyEncryptionSource_ASSEMBLY
    | providerEncryptionSource
    ;

// TSql160.g:1366
assemblyEncryptionSource
    : 'ASSEMBLY' identifier
    ;

// TSql160.g:1377
providerEncryptionSource
    : 'PROVIDER' identifier providerKeySourceOptions
    ;

// TSql160.g:1389
fileEncryptionSource
    : ( 'EXECUTABLE' )? File EqualsSign stringLiteral
    ;

// TSql160.g:1405
asymKeySpec
    : 'ALGORITHM' EqualsSign Identifier
    ;

// TSql160.g:1433
createCertificateParams
    : From certificateSource
    | ( ( encryptClause )? With certificateOption ( Comma certificateOption )* )
    ;

// TSql160.g:1460
createCertificateActivityFlag
    : 'ACTIVE' For 'BEGIN_DIALOG' EqualsSign optionOnOff
    ;

// TSql160.g:1472
certificateOption
    : Identifier EqualsSign stringLiteral
    ;

// TSql160.g:1484
certificateSource
    : ( fileEncryptionSource ( With privateKeySpec )? | assemblyEncryptionSource )
    ;

// TSql160.g:1504
encryptClause
    : 'ENCRYPTION' By 'PASSWORD' EqualsSign stringLiteral
    ;

// TSql160.g:1512
privateKeySpec
    : 'PRIVATE' Key LeftParenthesis certificatePrivateKeySpec ( Comma certificatePrivateKeySpec )* RightParenthesis
    ;

// TSql160.g:1520
// TSql160.g:1526: action may refuse the input: {if (vParent.PrivateKeyPath != null) throw GetUnexpectedTokenErrorException(tFile); else vParent.PrivateKeyPath = vFilePath;}
certificatePrivateKeySpec
    : passwordChangeOption
    | File EqualsSign stringLiteral
    ;

// TSql160.g:1534
// TSql160.g:1539: action may refuse the input: {if (TryMatch(tEncryptionDecryption,CodeGenerationSupporter.Encryption)) { if (vParent.EncryptionPassword != null) throw GetUnexpectedTokenErrorException(tEncryptionDecryption); else vParent.EncryptionPassword = vPassword; } else { Match(tEncryptionDec...
passwordChangeOption
    : ( 'ENCRYPTION' | 'DECRYPTION' ) By Identifier EqualsSign stringLiteral
    ;

// TSql160.g:1586
contractMessage
    : identifier 'SENT' By ( Any | ( 'INITIATOR' | 'TARGET' ) )
    ;

// TSql160.g:1641
credentialStatementBody
    : 'CREDENTIAL' identifier With Identity EqualsSign stringLiteral ( Comma 'SECRET' EqualsSign stringLiteral )?
    ;

// TSql160.g:1690
createServerAuditStatement
    : identifier auditTargetClause ( auditCreateWithClause )? ( Where eventBooleanExpression )?
    ;

// TSql160.g:1717
// TSql160.g:1762: action may refuse the input: {UpdateTokenInfo(vResult,tRParen); if (pathRequired && !pathOptionEncountered) { if(tFile != null) { ThrowParseErrorException("SQL46056", tFile, TSqlParserResource.SQL46056Message); } else { ThrowParseErrorException("SQL46126", tUrl, TSqlParserResource...
auditTargetClause
    : To ( ( File LeftParenthesis auditFileOption ( Comma auditFileOption )* | 'URL' LeftParenthesis auditUrlOption ( Comma auditUrlOption )? ) RightParenthesis | ( 'APPLICATION_LOG' | 'SECURITY_LOG' | 'EXTERNAL_MONITOR' ) )
    ;

// TSql160.g:1798
// TSql160.g:1800: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MaxSize)}?
// TSql160.g:1803: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MaxRolloverFiles)}?
// TSql160.g:1806: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.ReserveDiskSpace)}?
// TSql160.g:1809: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MaxFiles)}?
auditFileOption
    : maxSizeAuditFileOption_MAXSIZE
    | maxRolloverFilesAuditFileOption_MAX_ROLLOVER_FILES
    | reserveDiskSpaceAuditFileOption_RESERVE_DISK_SPACE
    | maxFilesAuditFileOption_MAX_FILES
    | pathAuditFileOption
    ;

// TSql160.g:1815
// TSql160.g:1817: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.RetentionDays)}?
auditUrlOption
    : retentionDaysAuditUrlOption_RETENTION_DAYS
    | pathAuditFileOption
    ;

// TSql160.g:1937
pathAuditFileOption
    : ( 'FILEPATH' | 'PATH' ) EqualsSign stringLiteral
    ;

// TSql160.g:1957
auditCreateWithClause
    : With LeftParenthesis auditCreateOption ( Comma auditCreateOption )* RightParenthesis
    ;

// TSql160.g:1976
auditWithClause
    : With LeftParenthesis auditOption ( Comma auditOption )* RightParenthesis
    ;

// TSql160.g:1996
auditCreateOption
    : Identifier EqualsSign ( queueDelayAuditOption | onFailureAuditOption | auditGuidAuditOption | operatorAuditOption )
    ;

// TSql160.g:2011
// TSql160.g:2019: predicate dropped: {TryMatch(tOption, CodeGenerationSupporter.OperatorAudit)}?
// TSql160.g:2022: predicate dropped: {TryMatch(tOption, CodeGenerationSupporter.State)}?
auditOption
    : Identifier EqualsSign ( queueDelayAuditOption | onFailureAuditOption | operatorAuditOption | stateAuditOption )
    ;

// TSql160.g:2027
queueDelayAuditOption
    : integer
    ;

// TSql160.g:2040
onFailureAuditOption
    : ( Continue | Shutdown | 'FAIL_OPERATION' )
    ;

// TSql160.g:2069
auditGuidAuditOption
    : stringLiteral
    ;

// TSql160.g:2083
operatorAuditOption
    : optionOnOff
    ;

// TSql160.g:2096
stateAuditOption
    : optionOnOff
    ;

// TSql160.g:2109
createServerAuditSpecificationStatement
    : 'SPECIFICATION' identifier auditSpecificationForClause ( createAuditSpecificationDetail ( Comma createAuditSpecificationDetail )* )? auditSpecificationStateOpt
    ;

// TSql160.g:2194
// TSql160.g:2231: action may refuse the input: {if(vTarget == null && (vResult.Options == null || vResult.Options.Count == 0) && vFilterPredicate == null) { ThrowIncorrectSyntaxErrorException(vAuditName); }}
alterServerAuditStatement
    : identifier ( 'MODIFY' 'NAME' EqualsSign identifier | ( auditTargetClause )? ( auditWithClause )? ( Where eventBooleanExpression )? | 'REMOVE' Where )
    ;

// TSql160.g:2293
// TSql160.g:2304: action may refuse the input: {// Additional options are only allowed when external authentication is set to ON if (vOptionValue.OptionState != OptionState.On) ThrowIncorrectSyntaxErrorException(vOptionValue);}
// TSql160.g:2316: action may refuse the input: {// Empty rule: setting external authorization to OFF is the only allowed if (vOptionValue.OptionState != OptionState.Off) ThrowIncorrectSyntaxErrorException(vOptionValue);}
alterServerConfigurationExternalAuthenticationContainerOption
    : onOffOptionValue ( LeftParenthesis ( alterServerConfigurationExternalAuthenticationOption ) RightParenthesis |  )
    ;

// TSql160.g:2324
// TSql160.g:2325: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.UseIdentity)}?
// TSql160.g:2328: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.CredentialName)}?
alterServerConfigurationExternalAuthenticationOption
    : alterServerConfigurationExternalAuthenticationUseIdentityOption_USE_IDENTITY
    | alterServerConfigurationExternalAuthenticationCredentialNameOption_CREDENTIAL_NAME
    ;

// TSql160.g:2368
alterServerConfigurationSoftNumaOption
    : onOffOptionValue
    ;

// TSql160.g:2396
// TSql160.g:2408: action may refuse the input: {// Additional options are only allowed when buffer pool extension is set to ON if (vOptionValue.OptionState != OptionState.On) ThrowIncorrectSyntaxErrorException(vOptionValue);}
// TSql160.g:2424: action may refuse the input: {// Empty rule: setting buffer pool extension to OFF is the only allowed if (vOptionValue.OptionState != OptionState.Off) ThrowIncorrectSyntaxErrorException(vOptionValue);}
alterServerConfigurationBufferPoolExtensionContainerOption
    : onOffOptionValue ( LeftParenthesis alterServerConfigurationBufferPoolExtensionFileNameOption Comma alterServerConfigurationBufferPoolExtensionSizeOption RightParenthesis |  )
    ;

// TSql160.g:2432
alterServerConfigurationBufferPoolExtensionFileNameOption
    : 'FILENAME' EqualsSign stringLiteralOptionValue
    ;

// TSql160.g:2444
// TSql160.g:2450: action may refuse the input: {Match(tSize, CodeGenerationSupporter.Size); if (vMemUnit != MemoryUnit.KB && vMemUnit != MemoryUnit.MB && vMemUnit != MemoryUnit.GB) ThrowIncorrectSyntaxErrorException(vSize); vResult.OptionKind = AlterServerConfigurationBufferPoolExtensionOptionHelpe...
alterServerConfigurationBufferPoolExtensionSizeOption
    : 'SIZE' EqualsSign integerLiteralOptionValue memUnit
    ;

// TSql160.g:2477
// TSql160.g:2495: predicate dropped: {vResult.OptionKind == AlterServerConfigurationDiagnosticsLogOptionKind.Path}?
// TSql160.g:2501: predicate dropped: {vResult.OptionKind == AlterServerConfigurationDiagnosticsLogOptionKind.MaxFiles}?
// TSql160.g:2487: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MaxUnderscoreSize)}?
alterServerConfigurationDiagnosticsLogOption
    : onOffOptionValue
    | alterServerConfigurationDiagnosticsLogMaxSizeOption_MAX_SIZE
    | Identifier EqualsSign ( stringOrDefaultLiteralOptionValue | integerOrDefaultLiteralOptionValue )
    ;

// TSql160.g:2551
// TSql160.g:2560: predicate dropped: {vResult.OptionKind == AlterServerConfigurationFailoverClusterPropertyOptionKind.SqlDumperDumpFlags}?
// TSql160.g:2566: predicate dropped: {vResult.OptionKind == AlterServerConfigurationFailoverClusterPropertyOptionKind.SqlDumperDumpPath}?
alterServerConfigurationFailoverClusterPropertyOption
    : Identifier EqualsSign ( binaryOrDefaultLiteralOptionValue | stringOrDefaultLiteralOptionValue | integerOrDefaultLiteralOptionValue )
    ;

// TSql160.g:2594
alterServerConfigurationHadrClusterOption
    : Identifier EqualsSign ( stringLiteralOptionValue | 'LOCAL' )
    ;

// TSql160.g:2651
affinityRangeList
    : affinityRange ( Comma affinityRange )*
    ;

// TSql160.g:2666
affinityRange
    : integer ( To integer )?
    ;

// TSql160.g:2684
// TSql160.g:2688: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Audit) && NextTokenMatches(CodeGenerationSupporter.Specification, 2)}?
// TSql160.g:2691: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Scoped) && NextTokenMatches(CodeGenerationSupporter.Credential, 2)}?
// TSql160.g:2694: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Scoped) && NextTokenMatches(CodeGenerationSupporter.Configuration, 2)}?
alterDatabaseStatements
    : Alter Database ( alterDatabaseAuditSpecification_AUDIT | alterDatabaseScopedCredentialStatement_SCOPED | alterDatabaseScopedConfigurationStatement_SCOPED | alterDatabase | alterDatabaseEncryptionKey )
    ;

// TSql160.g:2733
alterDatabaseScopedConfigClear
    : 'CLEAR' databaseConfigurationClearOption
    ;

// TSql160.g:2749
databaseConfigurationClearOption
    : Identifier ( binary )?
    ;

// TSql160.g:2767
// TSql160.g:2774: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MaxDop)}?
// TSql160.g:2777: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.QueryOptimizerHotFixes) || NextTokenMatches(CodeGenerationSupporter.ParameterSniffing) || NextTokenMatches(CodeGenerationSupporter.LegacyCardinalityEstimation)}?
alterDatabaseScopedConfigSet
    : Set ( alterDatabaseScopedMaxDopOption_MAXDOP | ( alterDatabaseScopedOnOffPrimaryOption_QUERY_OPTIMIZER_HOTFIXES | alterDatabaseScopedOnOffPrimaryOption_PARAMETER_SNIFFING | alterDatabaseScopedOnOffPrimaryOption_LEGACY_CARDINALITY_ESTIMATION ) | alterDatabaseScopedGenericOption )
    ;

// TSql160.g:2842
// TSql160.g:2858: action may refuse the input: {vValue = CreateIdentifierOrScalarExpressionFromIdentifier(CreateIdentifierFromToken(tPrimary)); UpdateTokenInfo(vValue, token); if (!forSecondary) { ThrowParseErrorException("SQL46115", vValue, TSqlParserResource.SQL46115Message); } vResult.GenericOpt...
// TSql160.g:2856: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.Primary)}?
alterDatabaseScopedGenericOption
    : identifier EqualsSign ( Primary | stringOrSignedIntegerOrIdentifier | onOff )
    ;

// TSql160.g:2880
// TSql160.g:2896: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Remove)}?
// TSql160.g:2898: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Modify)}?
alterDatabase
    : ( identifier | sqlCommandIdentifier | Current ) ( alterDbAdd | alterDbRemove_REMOVE | alterDbModify_MODIFY | alterDbSet | alterDbCollate | alterDbRebuild )
    ;

// TSql160.g:2926
alterDbCollate
    : collation
    ;

// TSql160.g:2930
alterDbRebuild
    : 'REBUILD' 'LOG' ( On fileDecl )?
    ;

// TSql160.g:2948
alterDbAdd
    : Add ( alterDbAddFile | alterDbAddFilegroup )
    ;

// TSql160.g:2958
alterDbAddFile
    : ( 'LOG' )? File fileDeclBodyList ( toFilegroup )?
    ;

// TSql160.g:2981
alterDbAddFilegroup
    : 'FILEGROUP' identifier ( Contains ( 'FILESTREAM' | 'MEMORY_OPTIMIZED_DATA' ) )?
    ;

// TSql160.g:3063
alterDbModifyAzureOptions
    : azureOptions
    ;

// TSql160.g:3069
alterDbModifyFile
    : File fileDecl
    ;

// TSql160.g:3083
alterDbModifyFilegroup
    : identifier ( ( 'NAME' EqualsSign identifier ) | Default | ( Identifier ( xactTermination )? ) )
    ;

// TSql160.g:3128
alterDbSet
    : Set dbOptionStateList ( xactTermination )?
    ;

// TSql160.g:3143
toFilegroup
    : To 'FILEGROUP' identifier
    ;

// TSql160.g:3150
xactTermination
    : With ( Rollback ( ( 'AFTER' integer ( 'SECONDS' )? ) | 'IMMEDIATE' ) | 'NO_WAIT' )
    ;

// TSql160.g:3190
dbOptionStateList
    : dbOptionStateItem ( Comma dbOptionStateItem )*
    ;

// TSql160.g:3206
// TSql160.g:3238: predicate dropped: {NextTokenMatchesOneOf(OptionValidForCreateDatabase())}?
// TSql160.g:3245: predicate dropped: {!NextTokenMatches(CodeGenerationSupporter.Ledger)}?
// TSql160.g:3208: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.CursorDefault)}?
// TSql160.g:3210: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Recovery)}?
// TSql160.g:3212: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.TargetRecoveryTime)}?
// TSql160.g:3214: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.PageVerify)}?
// TSql160.g:3216: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Partner)}?
// TSql160.g:3218: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Witness)}?
// TSql160.g:3220: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Parameterization)}?
// TSql160.g:3222: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.CompatibilityLevel)}?
// TSql160.g:3224: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.ChangeTracking)}?
// TSql160.g:3226: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.AcceleratedDatabaseRecovery)}?
// TSql160.g:3228: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Containment)}?
// TSql160.g:3230: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Hadr)}?
// TSql160.g:3232: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.DelayedDurability)}?
// TSql160.g:3234: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.AutoCreateStatistics)}?
// TSql160.g:3236: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.RemoteDataArchive)}?
// TSql160.g:3240: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.QueryStore)}?
// TSql160.g:3242: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.AutomaticTuning)}?
dbOptionStateItem
    : cursorDefaultDbOption_CURSOR_DEFAULT
    | recoveryDbOption_RECOVERY
    | targetRecoveryTimeDbOption_TARGET_RECOVERY_TIME
    | pageVerifyDbOption_PAGE_VERIFY
    | partnerOption_PARTNER
    | witnessOption_WITNESS
    | parameterizationOption_PARAMETERIZATION
    | compatibilityLevelDbOption_COMPATIBILITY_LEVEL
    | changeTrackingDbOption_CHANGE_TRACKING
    | acceleratedDatabaseRecoveryOption_ACCELERATED_DATABASE_RECOVERY
    | dbContainmentOption_CONTAINMENT
    | hadrDbOption_HADR
    | dbDelayedDurabilityOption_DELAYED_DURABILITY
    | autoCreateStatisticsDbOption_AUTO_CREATE_STATISTICS
    | remoteDataArchiveDbOption_REMOTE_DATA_ARCHIVE
    | createAlterDbOption
    | queryStoreDbOption_QUERY_STORE
    | automaticTuningDbOption_AUTOMATIC_TUNING
    | dbSingleIdentOption
    | alterDbOnOffOption
    ;

// TSql160.g:3249
dbSingleIdentOption
    : Identifier
    ;

// TSql160.g:3352
// TSql160.g:3368: action may refuse the input: {bool requiresEquals = OnOffSimpleDbOptionsHelper.Instance.RequiresEqualsSign(vResult.OptionKind); if((vHasEquals && !requiresEquals) || (!vHasEquals && requiresEquals)) { ThrowIncorrectSyntaxErrorException(tOption); } vResult.OptionState = vOptionState;}
alterDbOnOffOption
    : Identifier ( EqualsSign )? optionOnOff
    ;

// TSql160.g:3541
changeTrackingOnOptions
    : LeftParenthesis changeTrackingOneOption ( Comma changeTrackingOneOption )* RightParenthesis
    ;

// TSql160.g:3563
// TSql160.g:3580: action may refuse the input: {Match(tAutoCleanupChangeRetention, CodeGenerationSupporter.ChangeRetention); if (changeRetentionEncountered) { ThrowParseErrorException("SQL46050", tAutoCleanupChangeRetention, TSqlParserResource.SQL46050Message, tAutoCleanupChangeRetention.getText())...
changeTrackingOneOption
    : Identifier EqualsSign ( On | Off | integer Identifier )
    ;

// TSql160.g:3644
queryStoreOnOptions
    : LeftParenthesis queryStoreOneOption ( Comma queryStoreOneOption )* RightParenthesis
    ;

// TSql160.g:3663
// TSql160.g:3665: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.DesiredState)}?
// TSql160.g:3668: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.OperationMode)}?
// TSql160.g:3671: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.QueryCaptureMode)}?
// TSql160.g:3674: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.SizeBasedCleanupMode)}?
// TSql160.g:3677: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.FlushIntervalSeconds)}?
// TSql160.g:3680: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.FlushIntervalSecondsAlt)}?
// TSql160.g:3683: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.IntervalLengthMinutes)}?
// TSql160.g:3686: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MaxQdsSize)}?
// TSql160.g:3689: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MaxPlansPerQuery)}?
// TSql160.g:3692: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.CleanupPolicy)}?
// TSql160.g:3695: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.WaitStatsCaptureMode)}?
queryStoreOneOption
    : queryStoreDesiredStateOption_DESIRED_STATE
    | queryStoreDesiredStateOption_OPERATION_MODE
    | queryStoreCaptureModeOption_QUERY_CAPTURE_MODE
    | queryStoreSizeCleanupModeOption_SIZE_BASED_CLEANUP_MODE
    | queryStoreDataFlushIntervalOption_FLUSH_INTERVAL_SECONDS
    | queryStoreDataFlushIntervalOption_DATA_FLUSH_INTERVAL_SECONDS
    | queryStoreIntervalLengthOption_INTERVAL_LENGTH_MINUTES
    | queryStoreMaxQdsSizeOption_MAX_STORAGE_SIZE_MB
    | queryStoreMaxPlansPerQueryOption_MAX_PLANS_PER_QUERY
    | queryStoreTimeCleanupPolicy_CLEANUP_POLICY
    | queryStoreWaitStatsCaptureOption_WAIT_STATS_CAPTURE_MODE
    ;

// TSql160.g:3948
automaticTuningOptions
    : LeftParenthesis automaticTuningOneOption ( Comma automaticTuningOneOption )* RightParenthesis
    ;

// TSql160.g:3971
// TSql160.g:3974: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.ForceLastGoodPlan)}?
// TSql160.g:3977: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.CreateIndex)}?
// TSql160.g:3980: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.DropIndex)}?
// TSql160.g:3983: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MaintainIndex)}?
automaticTuningOneOption
    : ( automaticTuningForceLastGoodPlanOption_FORCE_LAST_GOOD_PLAN | automaticTuningCreateIndexOption_CREATE_INDEX | automaticTuningDropIndexOption_DROP_INDEX | automaticTuningMaintainIndexOption_MAINTAIN_INDEX )
    ;

// TSql160.g:4128
dbContainmentOption
    : 'CONTAINMENT' EqualsSign Identifier
    ;

// TSql160.g:4169
simpleHadrDbOption
    : ( 'RESUME' | 'SUSPEND' )
    | Off
    ;

// TSql160.g:4235
alterDatabaseEncryptionKey
    : 'ENCRYPTION' Key ( ( 'REGENERATE' With 'ALGORITHM' EqualsSign Identifier ( 'ENCRYPTION' By dekEncryptor )? ) | ( 'ENCRYPTION' By dekEncryptor ) )
    ;

// TSql160.g:4317
// TSql160.g:4324: action may refuse the input: {SensitivityClassification.OptionType optionType = SensitivityClassification.GetOptionTypeByName(tOption.getText()); switch (optionType) { case SensitivityClassification.OptionType.Undefined: ThrowParseErrorException("SQL46125", tOption, TSqlParserReso...
// TSql160.g:4347: action may refuse the input: {SensitivityClassification.OptionType optionType = SensitivityClassification.GetOptionTypeByName(tOption.getText()); switch (optionType) { case SensitivityClassification.OptionType.Rank: break; case SensitivityClassification.OptionType.Undefined: Throw...
sensitivityClassificationOption
    : Identifier EqualsSign ( nonEmptyString | Identifier )
    ;

// TSql160.g:4406
// TSql160.g:4411: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Audit) && NextTokenMatches(CodeGenerationSupporter.Specification, 2)}?
// TSql160.g:4414: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Scoped) && NextTokenMatches(CodeGenerationSupporter.Credential, 2)}?
createDatabaseStatements
    : Create Database ( createDatabaseAuditSpecification_AUDIT | createDatabaseScopedCredentialStatement_SCOPED | createDatabase | createDatabaseEncryptionKey )
    ;

// TSql160.g:4427
// TSql160.g:4439: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.Copy, 2)}?
// TSql160.g:4443: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.MaxSize,2) || NextTokenMatches(CodeGenerationSupporter.Edition,2) || NextTokenMatches(CodeGenerationSupporter.ServiceObjective,2)}?
// TSql160.g:4458: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.MaxSize,2) || NextTokenMatches(CodeGenerationSupporter.Edition,2) || NextTokenMatches(CodeGenerationSupporter.ServiceObjective,2)}?
createDatabase
    : identifier ( ( azureAsCopyOf )? ( azureOptions )? ( dbContainmentOption )? recoveryUnitList collationOpt ( azureOptions )? ( dbAddendums )? ( withCreateDbOptions )? )
    ;

// TSql160.g:4472
azureOptions
    : LeftParenthesis azureOption ( Comma azureOption )* RightParenthesis
    ;

// TSql160.g:4495
// TSql160.g:4501: predicate dropped: {LA(3) == Identifier}?
// TSql160.g:4497: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MaxSize)}?
// TSql160.g:4499: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Edition)}?
azureOption
    : azureMaxSizeDatabaseOption_MAXSIZE
    | azureEditionDatabaseOption_EDITION
    | azureElasticPoolServiceObjectiveDatabaseOption
    | azureServiceObjectiveDatabaseOption
    ;

// TSql160.g:4543
azureServiceObjectiveDatabaseOption
    : 'SERVICE_OBJECTIVE' EqualsSign stringLiteral
    ;

// TSql160.g:4556
azureElasticPoolServiceObjectiveDatabaseOption
    : 'SERVICE_OBJECTIVE' EqualsSign 'ELASTIC_POOL' LeftParenthesis 'NAME' EqualsSign identifier RightParenthesis
    ;

// TSql160.g:4571
azureAsCopyOf
    : As 'COPY' Of multiPartIdentifier
    ;

// TSql160.g:4582
recoveryUnitList
    : ( onDisk )? ( 'LOG' On fileDeclBodyList )?
    ;

// TSql160.g:4593
onDisk
    : On fileDecl ( Comma ( fileDecl | fileGroupDecl ) )*
    ;

// TSql160.g:4620
fileDeclBodyList
    : fileDeclBody ( Comma fileDeclBody )*
    ;

// TSql160.g:4635
fileDecl
    : Primary fileDeclBody
    | fileDeclBody
    ;

// TSql160.g:4644
fileGroupDecl
    : 'FILEGROUP' identifier ( Contains ( 'FILESTREAM' | 'MEMORY_OPTIMIZED_DATA' ) )? ( Default )? fileDeclBody
    ;

// TSql160.g:4680
fileDeclBody
    : LeftParenthesis fileOption ( Comma fileOption )* RightParenthesis
    ;

// TSql160.g:4703
// TSql160.g:4716: action may refuse the input: {if (!newNameAllowed) { ThrowParseErrorException("SQL46062", vResult, TSqlParserResource.SQL46062Message); }}
// TSql160.g:4704: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Name)}?
// TSql160.g:4706: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.FileName)}?
// TSql160.g:4708: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Size)}?
// TSql160.g:4710: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MaxSize)}?
// TSql160.g:4712: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.FileGrowth)}?
// TSql160.g:4714: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.NewName)}?
// TSql160.g:4723: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Offline)}?
fileOption
    : nameFileOption_NAME
    | fileNameFileOption_FILENAME
    | sizeFileOption_SIZE
    | maxSizeFileOption_MAXSIZE
    | fileGrowthFileOption_FILEGROWTH
    | newNameFileOption_NEWNAME
    | offlineFileOption_OFFLINE
    ;

// TSql160.g:4861
memUnit
    : Identifier
    ;

// TSql160.g:4869
// TSql160.g:4874: action may refuse the input: {vParent.AttachMode = AttachModeHelper.Instance.ParseOption(tAttach); if (vParent.AttachMode == AttachMode.Load) ThrowIncorrectSyntaxErrorException(tAttach); UpdateTokenInfo(vParent,tAttach);}
dbAddendums
    : ( For Identifier )
    | As 'SNAPSHOT' Of identifier
    ;

// TSql160.g:4889
withCreateDbOptions
    : With createAlterDbOption ( Comma createAlterDbOption )*
    ;

// TSql160.g:4906
createDbServiceBrokerOption
    : Identifier
    ;

// TSql160.g:4930
// TSql160.g:4933: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.RestrictedUser)}?
// TSql160.g:4936: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.DbChaining) || NextTokenMatches(CodeGenerationSupporter.Trustworthy) || NextTokenMatches(CodeGenerationSupporter.DataRetention)}?
// TSql160.g:4939: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.FileStream)}?
createAlterDbOption
    : ( restrictedUserCreateDbOption_RESTRICTED_USER | ( createDbOnOffOption_DB_CHAINING | createDbOnOffOption_TRUSTWORTHY | createDbOnOffOption_DATA_RETENTION ) | fileStreamCreateAlterDbOption_FILESTREAM | createDbServiceBrokerOption | createAlterDbEqualsSignOption )
    ;

// TSql160.g:4976
createAlterDbEqualsSignOption
    : Identifier EqualsSign ( createAlterDbLiteralOption | createAlterDbIdentifierOption | createAlterDbOnOffOption )
    ;

// TSql160.g:4985
createAlterDbOnOffOption
    : optionOnOff
    ;

// TSql160.g:5010
createAlterDbLiteralOption
    : integer
    ;

// TSql160.g:5036
createAlterDbIdentifierOption
    : 'SQL_LATIN1_GENERAL_CP1_CI_AS'
    ;

// TSql160.g:5067
databaseConfigurationOptionOnOffPrimary
    : Primary
    | On
    | Off
    ;

// TSql160.g:5084
optionOnOff
    : On
    | Off
    ;

// TSql160.g:5125
// TSql160.g:5148: action may refuse the input: {Match(tOption, CodeGenerationSupporter.NonTransactedAccess); if(vParent.NonTransactedAccess != null) { ThrowParseErrorException("SQL46049", tOption, TSqlParserResource.SQL46049Message, tOption.getText()); } vParent.NonTransactedAccess=vNonTransactedFi...
// TSql160.g:5159: action may refuse the input: {Match(tOption, CodeGenerationSupporter.DirectoryName); if(vParent.DirectoryName != null) { ThrowParseErrorException("SQL46049", tOption, TSqlParserResource.SQL46049Message, tOption.getText()); } vParent.DirectoryName=vDirectoryName;}
fileStreamOption
    : Identifier EqualsSign ( ( Off | Full | 'READ_ONLY' ) | ( stringLiteralOrNull ) )
    ;

// TSql160.g:5202
remoteDataArchiveDbSettingsList
    : LeftParenthesis remoteDataArchiveDbSetting ( Comma remoteDataArchiveDbSetting )* RightParenthesis
    ;

// TSql160.g:5225
// TSql160.g:5232: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Server)}?
// TSql160.g:5242: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Credential)}?
// TSql160.g:5252: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.FederatedServiceAccount)}?
remoteDataArchiveDbSetting
    : identifier_SERVER EqualsSign stringLiteral
    | identifier_CREDENTIAL EqualsSign identifier
    | identifier_FEDERATED_SERVICE_ACCOUNT EqualsSign onOffOptionValue
    ;

// TSql160.g:5295
auditSpecificationForClause
    : For 'SERVER' 'AUDIT' identifier
    ;

// TSql160.g:5307
createAuditSpecificationDetailDb
    : Add LeftParenthesis ( auditActionSpecification | databaseAuditActionGroup ) RightParenthesis
    ;

// TSql160.g:5326
createAuditSpecificationDetail
    : Add LeftParenthesis serverAuditActionGroup RightParenthesis
    ;

// TSql160.g:5339
auditSpecificationDetailDb
    : ( Add | Drop ) LeftParenthesis ( auditActionSpecification | databaseAuditActionGroup ) RightParenthesis
    ;

// TSql160.g:5369
auditSpecificationDetail
    : ( Add | Drop ) LeftParenthesis serverAuditActionGroup RightParenthesis
    ;

// TSql160.g:5393
databaseAuditActionGroup
    : Identifier
    ;

// TSql160.g:5402
serverAuditActionGroup
    : Identifier
    ;

// TSql160.g:5411
auditActionSpecification
    : actionWithQual ( Comma actionWithQual )* securityTargetObject By principal ( Comma principal )*
    ;

// TSql160.g:5445
actionWithQual
    : Select
    | Update
    | Insert
    | Delete
    | Execute
    | 'RECEIVE'
    | References
    ;

// TSql160.g:5491
auditSpecificationStateOpt
    : ( With LeftParenthesis 'STATE' EqualsSign optionOnOff RightParenthesis )?
    ;

// TSql160.g:5506
createDatabaseEncryptionKey
    : 'ENCRYPTION' Key With 'ALGORITHM' EqualsSign Identifier 'ENCRYPTION' By dekEncryptor
    ;

// TSql160.g:5522
dekEncryptor
    : dekEncryptorType identifier
    ;

// TSql160.g:5532
dekEncryptorType
    : 'SERVER' ( ( 'CERTIFICATE' ) | ( 'ASYMMETRIC' Key ) )
    ;

// TSql160.g:5557
// TSql160.g:5560: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Certificate)}?
// TSql160.g:5562: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Service)}?
backupStatements
    : Backup ( backupCertificateStatement_CERTIFICATE | backupServiceMasterKeyStatement_SERVICE | backupMasterKeyStatement | backupStatement )
    ;

// TSql160.g:5574
backupStatement
    : backupMain ( backupDevices )? ( backupOptions )?
    ;

// TSql160.g:5586
// TSql160.g:5588: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Service)}?
restoreStatements
    : Restore ( restoreServiceMasterKeyStatement_SERVICE | restoreMasterKeyStatement | restoreStatement )
    ;

// TSql160.g:5600
restoreStatement
    : ( ( restoreMain ( From devList )? ) | ( Identifier From devList ) ) ( restoreOptions )?
    ;

// TSql160.g:5623
backupMain
    : backupDatabase
    | backupTransactionLog
    ;

// TSql160.g:5629
backupDatabase
    : Database identifierOrVariable backupFileListOpt
    ;

// TSql160.g:5641
backupTransactionLog
    : 'LOG' identifierOrVariable
    ;

// TSql160.g:5654
backupFileListOpt
    : ( backupRestoreFile ( Comma backupRestoreFile )* )?
    ;

// TSql160.g:5670
restoreFileListOpt
    : ( backupRestoreFile ( Comma backupRestoreFile )* )?
    ;

// TSql160.g:5685
backupRestoreFile
    : ( File EqualsSign ( stringOrVariable | backupRestoreFileNameList ) )
    | ( Identifier EqualsSign ( stringOrVariable | backupRestoreFileNameList ) )
    | 'READ_WRITE_FILEGROUPS'
    ;

// TSql160.g:5730
backupRestoreFileNameList
    : LeftParenthesis stringOrVariable ( Comma stringOrVariable )* RightParenthesis
    ;

// TSql160.g:5749
backupDevices
    : To devList ( mirrorTo )*
    ;

// TSql160.g:5761
mirrorTo
    : 'MIRROR' To devList
    ;

// TSql160.g:5769
devList
    : deviceInfo ( Comma deviceInfo )*
    ;

// TSql160.g:5784
deviceInfo
    : identifierOrVariable
    | ( Identifier ) EqualsSign stringOrVariable
    ;

// TSql160.g:5806
backupOptions
    : With backupOption ( Comma backupOption )*
    ;

// TSql160.g:5821
// TSql160.g:5826: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Encryption)}?
backupOption
    : backupEncryptionOption_ENCRYPTION
    | Identifier
    | ( Identifier EqualsSign ( signedIntegerOrVariable | stringLiteral ) )
    ;

// TSql160.g:5867
backupEncrytor
    : dekEncryptorType EqualsSign identifier
    ;

// TSql160.g:5877
restoreMain
    : ( Database identifierOrVariable restoreFileListOpt )
    | ( 'LOG' identifierOrVariable restoreFileListOpt )
    ;

// TSql160.g:5901
restoreOptions
    : With restoreOptionsList
    ;

// TSql160.g:5905
restoreOptionsList
    : restoreOption ( Comma restoreOption )*
    ;

// TSql160.g:5920
// TSql160.g:5930: predicate dropped: {IsStopAtBeforeMarkRestoreOption(tOption)}?
// TSql160.g:5926: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.FileStream)}?
restoreOption
    : fileStreamRestoreOption_FILESTREAM
    | simpleRestoreOption
    | ( Identifier EqualsSign ( ( stringOrVariable afterClause ) | signedInteger | stringOrVariable ) )
    | moveRestoreOption
    | fileRestoreOption
    ;

// TSql160.g:5954
simpleRestoreOption
    : Identifier
    ;

// TSql160.g:5962
moveRestoreOption
    : 'MOVE' stringOrVariable To stringOrVariable
    ;

// TSql160.g:5975
fileRestoreOption
    : File EqualsSign signedIntegerOrVariable
    ;

// TSql160.g:5997
afterClause
    : 'AFTER' stringOrVariable
    ;

// TSql160.g:6046
restoreMasterKeyBody
    : 'MASTER' Key From File EqualsSign stringLiteral 'DECRYPTION' By 'PASSWORD' EqualsSign stringLiteral
    ;

// TSql160.g:6061
backupMasterKeyStatement
    : backupMasterKeyBody
    ;

// TSql160.g:6065
backupMasterKeyBody
    : 'MASTER' Key To File EqualsSign stringLiteral 'ENCRYPTION' By 'PASSWORD' EqualsSign stringLiteral
    ;

// TSql160.g:6079
restoreMasterKeyStatement
    : restoreMasterKeyBody 'ENCRYPTION' By 'PASSWORD' EqualsSign stringLiteral ( 'FORCE' )?
    ;

// TSql160.g:6104
bulkInsertStatement
    : Bulk Insert schemaObjectThreePartName From bulkInsertFrom ( bulkInsertOptions )?
    ;

// TSql160.g:6126
bulkInsertFrom
    : integer
    | stringOrIdentifier
    ;

// TSql160.g:6137
bulkInsertOptions
    : With LeftParenthesis bulkInsertOption ( Comma bulkInsertOption )* RightParenthesis
    ;

// TSql160.g:6160
bulkInsertOption
    : bulkInsertSortOrderOption
    | simpleBulkInsertOptionWithValue
    | simpleBulkInsertOptionNoValue
    ;

// TSql160.g:6166
insertBulkOption
    : bulkInsertSortOrderOption
    | simpleInsertBulkOption
    ;

// TSql160.g:6171
// TSql160.g:6194: action may refuse the input: {vResult=FragmentFactory.CreateFragment<BulkInsertOption>(); UpdateTokenInfo(vResult, tOption); vResult.OptionKind = BulkInsertFlagOptionsHelper.Instance.ParseOption(tOption); if (vResult.OptionKind == BulkInsertOptionKind.KeepIdentity) // For some rea...
simpleInsertBulkOption
    : Identifier ( ( EqualsSign integerOrNumeric ) |  )
    ;

// TSql160.g:6204
simpleBulkInsertOptionWithValue
    : Identifier EqualsSign ( integerOrNumeric | stringLiteral | ( 'TRUE' | 'FALSE' ) )
    ;

// TSql160.g:6263
simpleBulkInsertOptionNoValue
    : Identifier
    ;

// TSql160.g:6271
insertBulkStatement
    : Insert Bulk schemaObjectThreePartName ( coldefList )? ( insertBulkOptions )?
    ;

// TSql160.g:6290
coldefList
    : LeftParenthesis coldefItem ( Comma coldefItem )* RightParenthesis
    ;

// TSql160.g:6309
coldefItem
    : columnDefinitionEx ( nullNotNull )?
    ;

// TSql160.g:6325
insertBulkOptions
    : With LeftParenthesis insertBulkOption ( Comma insertBulkOption )* RightParenthesis
    ;

// TSql160.g:6347
bulkInsertSortOrderOption
    : Order LeftParenthesis columnWithSortOrder ( Comma columnWithSortOrder )* RightParenthesis
    ;

// TSql160.g:6368
openRowsetBulkOrderHint
    : bulkInsertSortOrderOption ( Unique )?
    ;

// TSql160.g:6381
// TSql160.g:6383: action may refuse the input: {DbccCommand command; if (DbccCommandsHelper.Instance.TryParseOption(tIdentifier,out command)) { if (command == DbccCommand.ConcurrencyViolation || command == DbccCommand.MemObjList || command == DbccCommand.MemoryMap) { ThrowIncorrectSyntaxErrorExcept...
dbccStatement
    : Dbcc Identifier ( dbccNamedLiteralList )? ( dbccOptions )?
    ;

// TSql160.g:6416
dbccOptions
    : With dbccOptionsList
    ;

// TSql160.g:6420
dbccOptionsList
    : dbccOptionsListItems
    | dbccOptionsJoin
    ;

// TSql160.g:6425
dbccOptionsListItems
    : dbccOption ( Comma dbccOption )*
    ;

// TSql160.g:6440
dbccOption
    : Identifier
    ;

// TSql160.g:6448
dbccOptionsJoin
    : dbccJoinOption ( Join dbccJoinOption )+
    ;

// TSql160.g:6464
dbccJoinOption
    : Identifier
    ;

// TSql160.g:6472
dbccNamedLiteralList
    : LeftParenthesis ( ( dbccNamedLiteral ( Comma dbccNamedLiteral )* ) |  ) RightParenthesis
    ;

// TSql160.g:6502
dbccNamedLiteral
    : ( Identifier EqualsSign )? possibleNegativeConstantOrIdentifier
    ;

// TSql160.g:6534
endpointOptions
    : endpointOption ( Comma endpointOption )*
    ;

// TSql160.g:6538
// TSql160.g:6539: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.State)}?
endpointOption
    : endpointState_STATE
    | endpointAffinity
    ;

// TSql160.g:6554
// TSql160.g:6560: action may refuse the input: {Match(tAffinity,CodeGenerationSupporter.Affinity); if (vParent.Affinity != null) throw GetUnexpectedTokenErrorException(tAffinity); UpdateTokenInfo(vAffinity, tAffinity);}
endpointAffinity
    : 'AFFINITY' EqualsSign ( ( 'NONE' | 'ADMIN' ) | integer )
    ;

// TSql160.g:6588
protocolInfo
    : As Identifier LeftParenthesis protocolOptionsItem ( Comma protocolOptionsItem )* RightParenthesis
    ;

// TSql160.g:6614
// TSql160.g:6617: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.ListenerIP)}?
// TSql160.g:6619: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.ListenerPort)}?
// TSql160.g:6621: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Path)}?
// TSql160.g:6623: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Site)}?
// TSql160.g:6625: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.ClearPort)}?
// TSql160.g:6627: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.SslPort)}?
// TSql160.g:6629: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Authentication)}?
// TSql160.g:6631: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Ports)}?
// TSql160.g:6633: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.AuthRealm)}?
// TSql160.g:6635: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.DefaultLogonDomain)}?
protocolOptionsItem
    : listenerIpProtocolOption_LISTENER_IP
    | listenerPortProtocolOption_LISTENER_PORT
    | pathProtocolOption_PATH
    | siteProtocolOption_SITE
    | clearPortProtocolOption_CLEAR_PORT
    | sslPortProtocolOption_SSL_PORT
    | authenticationProtocolOption_AUTHENTICATION
    | portsProtocolOption_PORTS
    | authenticationRealmProtocolOption_AUTH_REALM
    | defaultLogonDomainProtocolOption_DEFAULT_LOGON_DOMAIN
    | compressionProtocolOption
    ;

// TSql160.g:6785
compressionProtocolOption
    : Identifier EqualsSign ( 'ENABLED' | 'DISABLED' )
    ;

// TSql160.g:6839
ipAddressV4
    : Numeric ( Numeric ( Numeric ( integer |  ) | Dot integer | integer ipV4DotNumberTail ) | ipV4DotLastTwoPartsTail | integer ( Numeric ipV4DotNumberTail | ipV4DotLastTwoPartsTail ) )
    | integer ( Numeric ( Numeric ipV4DotNumberTail | ipV4DotLastTwoPartsTail ) | Dot ( Numeric ( Numeric ( integer |  ) | Dot integer | integer ipV4DotNumberTail ) | integer ( Numeric ipV4DotNumberTail | ipV4DotLastTwoPartsTail ) ) )
    ;

// TSql160.g:7018
ipV4DotLastTwoPartsTail
    : Dot ( ipV4NumericOrNumericIntegerTail | integer ipV4DotNumberTail )
    ;

// TSql160.g:7034
ipV4DotNumberTail
    : Numeric
    | Dot integer
    ;

// TSql160.g:7051
ipV4NumericOrNumericIntegerTail
    : Numeric ( integer |  )
    ;

// TSql160.g:7074
// TSql160.g:7099: action may refuse the input: {if ((vParent.EndpointType == EndpointType.DatabaseMirroring) && ((encountered & PayloadOptionKinds.Role) != PayloadOptionKinds.Role)) { ThrowParseErrorException("SQL46080", tEndpointType, TSqlParserResource.SQL46080Message); }}
payloadInfo
    : For 'DATA_MIRRORING' LeftParenthesis ( payloadOption ( Comma payloadOption )* )? RightParenthesis
    ;

// TSql160.g:7112
// TSql160.g:7114: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Batches)}?
// TSql160.g:7116: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Wsdl)}?
// TSql160.g:7118: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Sessions)}?
// TSql160.g:7120: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.LoginType)}?
// TSql160.g:7122: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.SessionTimeout)}?
// TSql160.g:7125: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Namespace)}?
// TSql160.g:7128: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.CharacterSet)}?
// TSql160.g:7130: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.HeaderLimit)}?
// TSql160.g:7132: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Authentication)}?
// TSql160.g:7134: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Encryption)}?
// TSql160.g:7136: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MessageForwarding)}?
// TSql160.g:7138: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MessageForwardSize)}?
payloadOption
    : soapMethod
    | batchesPayloadOption_BATCHES
    | wsdlPayloadOption_WSDL
    | sessionsPayloadOption_SESSIONS
    | loginTypePayloadOption_LOGIN_TYPE
    | sessionTimeoutPayloadOption_SESSION_TIMEOUT
    | databasePayloadOption
    | namespacePayloadOption_NAMESPACE
    | schemaPayloadOption
    | characterSetPayloadOption_CHARACTER_SET
    | headerLimitPayloadOption_HEADER_LIMIT
    | authenticationPayloadOption_AUTHENTICATION
    | encryptionPayloadOption_ENCRYPTION
    | messageForwardingPayloadOption_MESSAGE_FORWARDING
    | messageForwardSizePayloadOption_MESSAGE_FORWARD_SIZE
    | rolePayloadOption
    ;

// TSql160.g:7151
enabledDisabled
    : ( 'ENABLED' | 'DISABLED' )
    ;

// TSql160.g:7233
databasePayloadOption
    : Database EqualsSign ( stringLiteral | defaultLiteral )
    ;

// TSql160.g:7263
schemaPayloadOption
    : Schema EqualsSign ( 'STANDARD' | 'NONE' )
    ;

// TSql160.g:7433
rolePayloadOption
    : 'ROLE' EqualsSign ( ( 'WITNESS' | 'PARTNER' ) | All )
    ;

// TSql160.g:7457
soapMethod
    : ( ( Add | Alter )? 'WEBMETHOD' soapMethodAlias LeftParenthesis soapMethodOption ( Comma soapMethodOption )* RightParenthesis )
    | Drop 'WEBMETHOD' soapMethodAlias
    ;

// TSql160.g:7489
soapMethodAlias
    : stringLiteral
    | stringLiteral Dot stringLiteral
    ;

// TSql160.g:7504
// TSql160.g:7509: action may refuse the input: {if (vParent.Schema != SoapMethodSchemas.NotSpecified) throw GetUnexpectedTokenErrorException(tSchema);}
// TSql160.g:7530: action may refuse the input: {if (vParent.Format != SoapMethodFormat.NotSpecified) throw GetUnexpectedTokenErrorException(tFormatName); Match(tFormatName,CodeGenerationSupporter.Format); vParent.Format = SoapMethodFormatsHelper.Instance.ParseOption(tFormat);}
// TSql160.g:7537: action may refuse the input: {Match(tFormatName,CodeGenerationSupporter.Name); if (vParent.Name != null) throw GetUnexpectedTokenErrorException(tFormatName); vParent.Name = vName;}
soapMethodOption
    : Schema EqualsSign ( ( 'NONE' | 'STANDARD' ) | Default )
    | Identifier EqualsSign ( Identifier | stringLiteral )
    ;

// TSql160.g:7634
eventTypeGroupContainer
    : Identifier
    ;

// TSql160.g:7660
// TSql160.g:7672: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Queue)}?
eventNotificationObjectScope
    : On ( Database | 'QUEUE' schemaObjectThreePartName | 'SERVER' )
    ;

// TSql160.g:7688
eventSessionScope
    : On ( Database | 'SERVER' )
    ;

// TSql160.g:7797
eventDeclarationList
    : Add 'EVENT' eventDeclaration ( Comma Add 'EVENT' eventDeclaration )*
    ;

// TSql160.g:7816
targetDeclarationList
    : Add 'TARGET' targetDeclaration ( Comma Add 'TARGET' targetDeclaration )*
    ;

// TSql160.g:7835
sessionOptionList
    : ( With ( LeftParenthesis sessionOption ( Comma sessionOption )* RightParenthesis ) )
    ;

// TSql160.g:7862
optSessionOptionList
    : ( sessionOptionList )?
    ;

// TSql160.g:7870
eventDeclaration
    : eventSessionNonEmptyThreePartObjectName ( LeftParenthesis ( eventDeclarationSetParameters )? ( eventDeclarationActionParameters )? ( Where eventBooleanExpression )? RightParenthesis )?
    ;

// TSql160.g:7897
eventDeclarationSetParameters
    : ( Set eventDeclarationSetParameter ( Comma eventDeclarationSetParameter )* )
    ;

// TSql160.g:7916
eventDeclarationActionParameters
    : ( 'ACTION' LeftParenthesis eventSessionNonEmptyThreePartObjectName ( Comma eventSessionNonEmptyThreePartObjectName )* RightParenthesis )
    ;

// TSql160.g:7943
targetDeclaration
    : eventSessionNonEmptyThreePartObjectName ( LeftParenthesis ( Set eventDeclarationSetParameter ( Comma eventDeclarationSetParameter )* ) RightParenthesis )?
    ;

// TSql160.g:7974
// TSql160.g:7977: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MaxMemory)}?
// TSql160.g:7979: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MaxEventSize)}?
// TSql160.g:7981: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.TrackCausality)}?
// TSql160.g:7983: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.StartupState)}?
// TSql160.g:7985: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.EventRetentionMode)}?
// TSql160.g:7987: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MemoryPartitionMode)}?
// TSql160.g:7989: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MaxDispatchLatency)}?
sessionOption
    : maxMemorySessionOption_MAX_MEMORY
    | maxEventSizeSessionOption_MAX_EVENT_SIZE
    | trackCausalitySessionOption_TRACK_CAUSALITY
    | startupStateSessionOption_STARTUP_STATE
    | eventRetentionSessionOption_EVENT_RETENTION_MODE
    | memoryPartitionSessionOption_MEMORY_PARTITION_MODE
    | maxDispatchLatencySessionOption_MAX_DISPATCH_LATENCY
    ;

// TSql160.g:8000
integerSessionOptionValue
    : EqualsSign integer Identifier
    ;

// TSql160.g:8061
onOffSessionOption
    : EqualsSign optionOnOff
    ;

// TSql160.g:8089
eventDeclarationSetParameter
    : identifier EqualsSign eventDeclarationValue
    ;

// TSql160.g:8102
eventDeclarationValue
    : Minus integerOrRealOrNumeric
    | eventSessionExpressionParenthesis
    | eventSessionLiteral
    ;

// TSql160.g:8122
eventSessionLiteral
    : integer
    | numeric
    | real
    | stringLiteral
    ;

// TSql160.g:8129
eventSessionExpressionParenthesis
    : LeftParenthesis eventDeclarationValue RightParenthesis
    ;

// TSql160.g:8142
eventBooleanExpression
    : eventBooleanExpressionOr
    ;

// TSql160.g:8147
eventBooleanExpressionOr
    : eventBooleanExpressionAnd ( Or eventBooleanExpressionAnd )*
    ;

// TSql160.g:8161
eventBooleanExpressionAnd
    : eventBooleanExpressionUnary ( And eventBooleanExpressionUnary )*
    ;

// TSql160.g:8175
eventBooleanExpressionUnary
    : Not eventBooleanExpressionUnary
    | eventBooleanExpressionParenthesis
    | eventDeclarationPredicateParameter
    ;

// TSql160.g:8193
eventBooleanExpressionParenthesis
    : LeftParenthesis eventBooleanExpression RightParenthesis
    ;

// TSql160.g:8206
// TSql160.g:8220: action may refuse the input: {if (vName.MultiPartIdentifier.Count < 2) ThrowIncorrectSyntaxErrorException(tLparen); vSourceDeclaration.Value = vName; eventDeclarationCompareFunctionValue.Name = vName;}
eventDeclarationPredicateParameter
    : eventSessionOneOrTwoOrThreePartObjectName ( ( LeftParenthesis ( eventDeclarationCompareFunction ) RightParenthesis ) | ( eventDeclarationComparisonPredicate ) )
    ;

// TSql160.g:8246
eventDeclarationCompareFunction
    : eventSessionOneOrTwoOrThreePartObjectName Comma eventDeclarationValue
    ;

// TSql160.g:8260
// TSql160.g:8266: predicate dropped: {LA(2) == Like}?
eventDeclarationComparisonPredicate
    : ( comparisonOperator | Not Like ) eventDeclarationValue
    ;

// TSql160.g:8273
dropEventDeclarationList
    : Drop 'EVENT' eventSessionNonEmptyThreePartObjectName ( Comma Drop 'EVENT' eventSessionNonEmptyThreePartObjectName )*
    ;

// TSql160.g:8292
dropTargetDeclarationList
    : Drop 'TARGET' eventSessionNonEmptyThreePartObjectName ( Comma Drop 'TARGET' eventSessionNonEmptyThreePartObjectName )*
    ;

// TSql160.g:8322
createFulltextCatalogStatement
    : 'CATALOG' identifier ( On 'FILEGROUP' identifier )? ( In 'PATH' stringLiteral )? ( With accentSensitivity )? ( As Default )? authorizationOpt
    ;

// TSql160.g:8363
accentSensitivity
    : 'ACCENT_SENSITIVITY' EqualsSign optionOnOff
    ;

// TSql160.g:8375
// TSql160.g:8386: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Type)}?
// TSql160.g:8394: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Language)}?
fulltextIndexColumn
    : identifier ( 'TYPE' Column identifier )? ( languageTerm_LANGUAGE )? ( 'STATISTICAL_SEMANTICS' )?
    ;

// TSql160.g:8410
createFulltextIndexStatement
    : Index On schemaObjectFourPartName ( LeftParenthesis fulltextIndexColumn ( Comma fulltextIndexColumn )* RightParenthesis )? Key Index identifier ( fullTextCatalogAndFileGroup )? createFulltextOptions
    ;

// TSql160.g:8448
fullTextCatalogAndFileGroup
    : On ( identifier | LeftParenthesis ( ( identifier ( Comma 'FILEGROUP' identifier )? ) | ( 'FILEGROUP' identifier ( Comma identifier )? ) ) RightParenthesis )
    ;

// TSql160.g:8500
createFulltextOptions
    : ( With ( createFulltextOptionsList | LeftParenthesis createFulltextOptionsList RightParenthesis ) )?
    ;

// TSql160.g:8516
createFulltextOptionsList
    : fullTextIndexOption ( Comma fullTextIndexOption )*
    ;

// TSql160.g:8536
// TSql160.g:8538: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.ChangeTracking)}?
fullTextIndexOption
    : changeTrackingFulltextIndexOption_CHANGE_TRACKING
    | stoplistFulltextIndexOption
    | searchPropertyListFullTextIndexOption
    ;

// TSql160.g:8583
stoplistFulltextIndexOption
    : StopList ( EqualsSign )? ( identifier | Off )
    ;

// TSql160.g:8608
searchPropertyListFullTextIndexOption
    : 'SEARCH' 'PROPERTY' 'LIST' ( EqualsSign )? ( identifier | Off )
    ;

// TSql160.g:8634
createFulltextStoplistStatement
    : StopList identifier ( From ( ( identifier Dot )? identifier | 'SYSTEM' StopList ) )? authorizationOpt requiredSemicolon
    ;

// TSql160.g:8672
alterFulltextStoplistStatement
    : StopList identifier ( fulltextStoplistAction ) requiredSemicolon
    ;

// TSql160.g:8691
fulltextStoplistAction
    : addFulltextStoplistAction
    | dropFulltextStoplistAction
    ;

// TSql160.g:8698
languageTerm
    : 'LANGUAGE' ( stringLiteral | integer | identifier | binary )
    ;

// TSql160.g:8735
addFulltextStoplistAction
    : Add ( stringLiteral ) languageTerm
    ;

// TSql160.g:8757
dropFulltextStoplistAction
    : Drop ( stringLiteral languageTerm | ( All ( languageTerm )? ) )
    ;

// TSql160.g:8793
dropFulltextStoplistStatement
    : StopList identifier requiredSemicolon
    ;

// TSql160.g:8845
searchPropertyListAction
    : addSearchPropertyListAction
    | dropSearchPropertyListAction
    ;

// TSql160.g:8851
addSearchPropertyListAction
    : Add stringLiteral With LeftParenthesis 'PROPERTY_SET_GUID' EqualsSign stringLiteral Comma 'PROPERTY_INT_ID' EqualsSign integer ( Comma Identifier EqualsSign stringLiteral )? RightParenthesis
    ;

// TSql160.g:8885
dropSearchPropertyListAction
    : Drop stringLiteral
    ;

// TSql160.g:8910
// TSql160.g:8913: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.Table)}?
// TSql160.g:8916: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Data)}?
// TSql160.g:8919: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.File)}?
// TSql160.g:8922: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Resource)}?
// TSql160.g:8925: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Library)}?
// TSql160.g:8928: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Language)}?
// TSql160.g:8931: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Stream)}?
createExternalStatements
    : External ( createExternalTableStatement | createExternalDataSourceStatement_DATA | createExternalFileFormatStatement | createExternalResourcePoolStatement_RESOURCE | createExternalLibraryStatement_LIBRARY | createExternalLanguageStatement_LANGUAGE | createExternalStreamStatement_STREAM )
    ;

// TSql160.g:8936
// TSql160.g:8939: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Data)}?
// TSql160.g:8942: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Resource)}?
// TSql160.g:8945: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Library)}?
// TSql160.g:8948: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Language)}?
alterExternalStatements
    : External ( alterExternalDataSourceStatement_DATA | alterExternalResourcePoolStatement_RESOURCE | alterExternalLibraryStatement_LIBRARY | alterExternalLanguageStatement_LANGUAGE )
    ;

// TSql160.g:8953
createPrimaryXmlIndexStatement
    : Primary 'XML' createXmlIndexStatement
    ;

// TSql160.g:9035
createXmlSchemaCollectionStatement
    : Schema 'COLLECTION' schemaObjectNonEmptyTwoPartName As expression
    ;

// TSql160.g:9053
// TSql160.g:9063: predicate dropped: {vPrimary}?
createXmlIndexStatement
    : Index identifier On schemaObjectThreePartName LeftParenthesis identifier RightParenthesis (  | 'USING' 'XML' Index identifier ( secondaryXmlIndexStatementBody | secondarySelectiveXmlIndex ) ) ( With indexOptionList )?
    ;

// TSql160.g:9096
secondaryXmlIndexStatementBody
    : For Identifier
    ;

// TSql160.g:9112
secondarySelectiveXmlIndex
    : For LeftParenthesis identifier RightParenthesis
    ;

// TSql160.g:9133
// TSql160.g:9179: action may refuse the input: {if(vXQuery != null && vXQueryType == null && vMaxLength == null && !vSingleton) { ThrowIncorrectSyntaxErrorException(vXQuery); }}
// TSql160.g:9152: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.XQuery)}?
promotedSelectiveXmlIndexPath
    : identifier EqualsSign stringLiteral ( As ( identifier_XQUERY ( stringLiteral )? maxlengthOption | 'SQL' scalarDataType ) singletonOption )?
    ;

// TSql160.g:9187
// TSql160.g:9193: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Singleton)}?
singletonOption
    : ( identifier_SINGLETON )?
    ;

// TSql160.g:9201
// TSql160.g:9208: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MaxLength)}?
maxlengthOption
    : ( identifier_MAXLENGTH LeftParenthesis integer RightParenthesis )?
    ;

// TSql160.g:9220
promotedSelectiveXmlIndexPathInAlter
    : Add promotedSelectiveXmlIndexPath
    | 'REMOVE' identifier
    ;

// TSql160.g:9269
passwordLoginSource
    : With 'PASSWORD' EqualsSign loginPassword ( createLoginPasswordOption )* ( Comma createLoginParam )*
    ;

// TSql160.g:9289
loginPassword
    : binary
    | stringLiteral
    ;

// TSql160.g:9294
// TSql160.g:9296: action may refuse the input: {if (TryMatch(tHashedMustChange, CodeGenerationSupporter.Hashed)) { if (vParent.Hashed) throw GetUnexpectedTokenErrorException(tHashedMustChange); else vParent.Hashed = true; } else { Match(tHashedMustChange, CodeGenerationSupporter.MustChange); if (vP...
createLoginPasswordOption
    : ( 'HASHED' | 'MUST_CHANGE' )
    ;

// TSql160.g:9316
createLoginParam
    : Identifier EqualsSign ( onOffPrincipalOption | binary | identifier )
    ;

// TSql160.g:9345
onOffPrincipalOption
    : optionOnOff
    ;

// TSql160.g:9383
// TSql160.g:9388: action may refuse the input: {vResult.OptionKind = IdentifierCreateLoginOptionsHelper.Instance.ParseOption(tOption); if (vResult.OptionKind == PrincipalOptionKind.Credential) throw GetUnexpectedTokenErrorException(tOption); vResult.Identifier = vValue;}
createLoginParamWin
    : Identifier EqualsSign identifier
    ;

// TSql160.g:9396
// TSql160.g:9425: action may refuse the input: {if (typeSidOption == 1) ThrowParseErrorException("SQL46121", vResult, TSqlParserResource.SQL46121Message);}
externalLoginSource
    : External 'PROVIDER' ( With createLoginParamExt ( Comma createLoginParamExt )* )?
    ;

// TSql160.g:9432
// TSql160.g:9440: action may refuse the input: {if(!TryMatch(tOption,CodeGenerationSupporter.Type) &&!TryMatch(tOption, CodeGenerationSupporter.Sid) &&!TryMatch(tOption, CodeGenerationSupporter.DefaultLanguage) &&!TryMatch(tOption, CodeGenerationSupporter.DefaultDatabase)) ThrowParseErrorException(...
createLoginParamExt
    : ( 'TYPE' | 'SID' | 'DEFAULT_LANGUAGE' | 'DEFAULT_DATABASE' ) EqualsSign ( identifier | binary )
    ;

// TSql160.g:9475
certificateLoginSource
    : 'CERTIFICATE' identifier createLoginParamCertOpt
    ;

// TSql160.g:9488
createLoginParamCertOpt
    : ( With 'CREDENTIAL' EqualsSign identifier )?
    ;

// TSql160.g:9498
asymmetricKeyLoginSource
    : 'ASYMMETRIC' Key identifier createLoginParamCertOpt
    ;

// TSql160.g:9525
// TSql160.g:9539: action may refuse the input: {if (vParent.ValidationMethod != MessageValidationMethod.ValidXml) throw GetUnexpectedTokenErrorException(tWith); Match(tCollection, CodeGenerationSupporter.Collection); vParent.XmlSchemaCollectionName = vSchemaObjectName;}
// TSql160.g:9548: action may refuse the input: {if (vParent.ValidationMethod == MessageValidationMethod.ValidXml && vParent.XmlSchemaCollectionName == null) { ThrowIncorrectSyntaxErrorException(tOption); }}
messageTypeValidation
    : 'VALIDATION' EqualsSign Identifier ( With Schema 'COLLECTION' schemaObjectNonEmptyTwoPartName )? (  )
    ;

// TSql160.g:9631
alterAvailabilityGroupAddDatabase
    : Add Database identifier ( Comma identifier )*
    ;

// TSql160.g:9648
alterAvailabilityGroupRemoveDatabase
    : 'REMOVE' Database identifier ( Comma identifier )*
    ;

// TSql160.g:9666
alterAvailabilityGroupAddReplica
    : Add 'REPLICA' On availabilityReplica ( Comma availabilityReplica )*
    ;

// TSql160.g:9722
alterAvailabilityGroupSetOption
    : Set LeftParenthesis availabilityGroupOption RightParenthesis
    ;

// TSql160.g:9734
// TSql160.g:9740: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Failover)}?
alterAvailabilityGroupTakeAction
    : ( alterAvailabilityGroupFailoverAction_FAILOVER | alterAvailabilityGroupAction )
    ;

// TSql160.g:9771
alterAvailabilityGroupFailoverOption
    : 'TARGET' EqualsSign stringLiteral
    ;

// TSql160.g:9783
alterAvailabilityGroupAction
    : Identifier
    | Join
    ;

// TSql160.g:9797
availabilityGroupOption
    : 'REQUIRED_COPIES_TO_COMMIT' EqualsSign integer
    ;

// TSql160.g:9808
availabilityReplicaName
    : stringLiteral
    ;

// TSql160.g:9817
availabilityReplica
    : stringLiteral With LeftParenthesis ( availabilityReplicaOption ( Comma availabilityReplicaOption )* ) RightParenthesis
    ;

// TSql160.g:9848
// TSql160.g:9850: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.AvailabilityMode)}?
// TSql160.g:9853: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.FailoverMode)}?
// TSql160.g:9856: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.SecondaryRole)}?
// TSql160.g:9859: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.PrimaryRole)}?
availabilityReplicaOption
    : availabilityModeReplicaOption_AVAILABILITY_MODE
    | failoverModeReplicaOption_FAILOVER_MODE
    | secondaryRoleReplicaOption_SECONDARY_ROLE
    | primaryRoleReplicaOption_PRIMARY_ROLE
    | literalReplicaOption
    ;

// TSql160.g:9973
literalReplicaOption
    : Identifier EqualsSign ( stringLiteral | integer )
    ;

// TSql160.g:10046
createPartitionFunction
    : Function identifier LeftParenthesis partitionParameterType RightParenthesis As partitionFunctionTypes
    ;

// TSql160.g:10064
partitionParameterType
    : scalarDataType collationOpt
    ;

// TSql160.g:10075
partitionFunctionTypes
    : 'RANGE' partitionFunctionRange For Values LeftParenthesis ( expression ( Comma expression )* )? RightParenthesis
    ;

// TSql160.g:10100
partitionFunctionRange
    : /* empty */
    | Left
    | Right
    ;

// TSql160.g:10114
partitionOptionTypes
    : 'RANGE' partitionTableOptionRange For Values LeftParenthesis ( expression ( Comma expression )* )? RightParenthesis
    ;

// TSql160.g:10139
partitionTableOptionRange
    : /* empty */
    | Left
    | Right
    ;

// TSql160.g:10153
createPartitionScheme
    : 'SCHEME' identifier As 'PARTITION' identifier ( All )? To LeftParenthesis stringOrIdentifier ( Comma stringOrIdentifier )* RightParenthesis
    ;

// TSql160.g:10264
bindingUserOption
    : User EqualsSign identifier
    ;

// TSql160.g:10275
bindingAnonymousOption
    : 'ANONYMOUS' EqualsSign optionOnOff
    ;

// TSql160.g:10357
onQueueClause
    : On 'QUEUE' schemaObjectThreePartName
    ;

// TSql160.g:10369
serviceContract
    : identifier
    ;

// TSql160.g:10379
alterServiceStatement
    : identifier ( ( onQueueClause ( addDropServiceContractList )? ) | addDropServiceContractList )
    ;

// TSql160.g:10398
addDropServiceContractList
    : LeftParenthesis addDropServiceContract ( Comma addDropServiceContract )* RightParenthesis
    ;

// TSql160.g:10417
addDropServiceContract
    : ( Drop | Add ) Identifier identifier
    ;

// TSql160.g:10470
keySpecList
    : keySpec ( Comma keySpec )*
    ;

// TSql160.g:10488
// TSql160.g:10490: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.KeySource)}?
// TSql160.g:10493: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Algorithm)}?
keySpec
    : keySourceKeySpec_KEY_SOURCE
    | algorithmKeySpec_ALGORITHM
    | identityValueKeySpec
    ;

// TSql160.g:10515
// TSql160.g:10520: action may refuse the input: {Match(tIdentityValue, CodeGenerationSupporter.IdentityValue); vResult.OptionKind=KeyOptionKind.IdentityValue; if (identityValueSpecified) throw GetUnexpectedTokenErrorException(tIdentityValue); identityValueSpecified = true; UpdateTokenInfo(vResult, ...
identityValueKeySpec
    : 'IDENTITY_VALUE' EqualsSign stringLiteral
    ;

// TSql160.g:10562
// TSql160.g:10565: action may refuse the input: {Match(tCreationDisposition, CodeGenerationSupporter.CreationDisposition); vResult.OptionKind=KeyOptionKind.CreationDisposition; if (creationDispositionSpecified) throw GetUnexpectedTokenErrorException(tCreationDisposition); creationDispositionSpecifi...
creationDispositionSourceOption
    : 'CREATION_DISPOSITION' EqualsSign ( 'CREATE_NEW' | 'OPEN_EXISTING' )
    ;

// TSql160.g:10583
encryptorListWithNoPassword
    : 'ENCRYPTION' By cryptoWithNoPasswordOrJustPassword ( Comma cryptoWithNoPasswordOrJustPassword )*
    ;

// TSql160.g:10599
cryptoWithNoPasswordOrJustPassword
    : certificateCrypto
    | keyCrypto
    | passwordCrypto
    ;

// TSql160.g:10605
certificateCrypto
    : 'CERTIFICATE' identifier
    ;

// TSql160.g:10617
keyCrypto
    : ( 'SYMMETRIC' | 'ASYMMETRIC' ) Key identifier
    ;

// TSql160.g:10634
passwordCrypto
    : 'PASSWORD' EqualsSign stringLiteral
    ;

// TSql160.g:10646
providerKeySourceOptionsOpt
    : ( providerKeySourceOptions )?
    ;

// TSql160.g:10654
providerKeySourceOptions
    : With providerKeySourceOption ( Comma providerKeySourceOption )*
    ;

// TSql160.g:10672
// TSql160.g:10674: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Algorithm)}?
// TSql160.g:10677: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.ProviderKeyName)}?
providerKeySourceOption
    : algorithmKeySpec_ALGORITHM
    | providerKeyNameSourceOption_PROVIDER_KEY_NAME
    | creationDispositionSourceOption
    ;

// TSql160.g:10757
createTypeUddtStatement
    : From scalarDataType ( nullableConstraint )?
    ;

// TSql160.g:10780
createTypeUdtStatement
    : External assemblyName
    ;

// TSql160.g:10796
createTypeTableStatement
    : As Table LeftParenthesis tableDefinition RightParenthesis ( withTypeTableOptions )?
    ;

// TSql160.g:10815
withTypeTableOptions
    : With LeftParenthesis memoryOptimizedTableOption RightParenthesis
    ;

// TSql160.g:10829
userLoginOption
    : ( For | From ) ( ( ( 'ASYMMETRIC' Key | Identifier ) identifier ) | External 'PROVIDER' )
    | 'WITHOUT' 'LOGIN'
    ;

// TSql160.g:10875
createUserStatement
    : User identifier ( userLoginOption )? ( With createUserOption ( Comma createUserOption )* )?
    ;

// TSql160.g:10911
// TSql160.g:10917: action may refuse the input: {if(vHasUserLoginOption && vResult.OptionKind != PrincipalOptionKind.DefaultSchema && vResult.OptionKind != PrincipalOptionKind.Sid && vResult.OptionKind != PrincipalOptionKind.Type && vResult.OptionKind != PrincipalOptionKind.Object_ID) { ThrowParseE...
createUserOption
    : Identifier EqualsSign ( identifierCreateUserOption | literalCreateUserOption )
    ;

// TSql160.g:10925
identifierCreateUserOption
    : identifier
    ;

// TSql160.g:10948
literalCreateUserOption
    : binary
    | stringLiteral
    | integer
    | nullLiteral
    ;

// TSql160.g:10986
// TSql160.g:10990: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Application)}?
// TSql160.g:10993: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Assembly)}?
// TSql160.g:10996: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Asymmetric)}?
// TSql160.g:11001: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Availability)}?
// TSql160.g:11004: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Broker)}?
// TSql160.g:11007: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Cryptographic)}?
// TSql160.g:11010: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Event)}?
// TSql160.g:11013: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.External)}?
// TSql160.g:11016: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Remote)}?
// TSql160.g:11019: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Resource)}?
// TSql160.g:11022: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Workload)}?
// TSql160.g:11025: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Certificate)}?
// TSql160.g:11028: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Credential)}?
// TSql160.g:11031: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Endpoint)}?
// TSql160.g:11034: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Fulltext)}?
// TSql160.g:11039: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Login)}?
// TSql160.g:11042: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Master)}?
// TSql160.g:11045: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Message)}?
// TSql160.g:11048: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Partition)}?
// TSql160.g:11051: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Queue)}?
// TSql160.g:11054: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Role)}?
// TSql160.g:11057: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Route)}?
// TSql160.g:11060: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Search)}?
// TSql160.g:11065: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Sequence)}?
// TSql160.g:11068: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Service)}?
// TSql160.g:11071: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Symmetric)}?
// TSql160.g:11074: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Server)}?
// TSql160.g:11079: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Xml)}?
// TSql160.g:11082: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Security)}?
// TSql160.g:11085: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.Column)}?
// TSql160.g:11088: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Materialized)}?
alter2005Statements
    : Alter ( alterApplicationRoleStatement_APPLICATION | alterAssemblyStatement_ASSEMBLY | alterAsymmetricKeyStatement_ASYMMETRIC | alterAuthorizationStatement | alterAvailabilityGroupStatement_AVAILABILITY | alterBrokerPriorityStatement_BROKER | alterCryptographicProviderStatement_CRYPTOGRAPHIC | alterEventSessionStatement_EVENT | alterExternalStatements | alterRemoteServiceBindingStatement_REMOTE | alterResourceStatement_RESOURCE | alterWorkloadGroupStatement_WORKLOAD | alterCertificateStatement_CERTIFICATE | alterCredentialStatement_CREDENTIAL | alterEndpointStatement_ENDPOINT | alterFulltextStatement_FULLTEXT | alterIndexStatement | alterLoginStatement_LOGIN | alterMasterKeyStatement_MASTER | alterMessageTypeStatement_MESSAGE | alterPartitionStatement_PARTITION | alterQueueStatement_QUEUE | alterRoleStatement_ROLE | alterRouteStatement_ROUTE | alterSearchPropertyListStatement_SEARCH | alterSchemaStatement | alterSequenceStatement_SEQUENCE | alterServiceStatements_SERVICE | alterSymmetricKeyStatement_SYMMETRIC | alterServerStatements_SERVER | alterUserStatement | alterXmlSchemaCollectionStatement_XML | alterSecurityPolicyStatement_SECURITY | alterColumnEncryptionKeyStatement | alterMaterializedViewStatement_MATERIALIZED )
    ;

// TSql160.g:11097
// TSql160.g:11111: action may refuse the input: {vResult.OptionKind = ApplicationRoleOptionHelper.Instance.ParseOption(tOption); if (defaultSchemaAndPasswordOnly && vResult.OptionKind != ApplicationRoleOptionKind.DefaultSchema) { throw GetUnexpectedTokenErrorException(tOption); } vResult.Value = Id...
// TSql160.g:11103: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Password)}?
applicationRoleOption
    : 'PASSWORD' EqualsSign stringLiteral
    | Identifier EqualsSign identifier
    ;

// TSql160.g:11193
alterAssemblyFromClause
    : From expression
    ;

// TSql160.g:11197
alterAssemblyWith
    : With alterAssemblyWithItem ( Comma alterAssemblyWithItem )*
    ;

// TSql160.g:11216
alterAssemblyWithItem
    : assemblyPermissionSetOrVisibilityOption
    | assemblyUncheckedDataOption
    ;

// TSql160.g:11222
assemblyPermissionSetOrVisibilityOption
    : Identifier EqualsSign ( assemblyPermissionSetOption | assemblyVisibilityOption )
    ;

// TSql160.g:11230
assemblyPermissionSetOption
    : Identifier
    ;

// TSql160.g:11242
assemblyVisibilityOption
    : optionOnOff
    ;

// TSql160.g:11255
assemblyUncheckedDataOption
    : 'UNCHECKED' 'DATA'
    ;

// TSql160.g:11266
alterAssemblyDropFile
    : Drop File ( All | stringLiteral ( Comma stringLiteral )* )
    ;

// TSql160.g:11290
alterAssemblyAddFile
    : Add File From alterAssemblyAddFileSpec ( Comma alterAssemblyAddFileSpec )*
    ;

// TSql160.g:11305
alterAssemblyAddFileSpec
    : expression ( As stringLiteral )?
    ;

// TSql160.g:11512
alterFulltextCatalogStatement
    : Identifier identifier ( ( ( 'REORGANIZE' | 'REBUILD' ) ( With accentSensitivity )? ) | As Default )
    ;

// TSql160.g:11559
alterFulltextIndexStatement
    : Index On schemaObjectThreePartName alterFulltextIndexAction
    ;

// TSql160.g:11575
// TSql160.g:11576: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Enable) || NextTokenMatches(CodeGenerationSupporter.Disable)}?
// TSql160.g:11580: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Start)}?
alterFulltextIndexAction
    : ( enableDisableAlterFulltextIndexAction_ENABLE | enableDisableAlterFulltextIndexAction_DISABLE )
    | addAlterFulltextIndexAction
    | dropAlterFulltextIndexAction
    | startPopulationAlterFulltextIndexAction_START
    | otherPopulationAlterFulltextIndexAction
    | setAlterFullTextIndexAction
    | alterColumnAlterFullTextIndexAction
    ;

// TSql160.g:11587
// TSql160.g:11592: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.ChangeTracking)}?
// TSql160.g:11594: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Search)}?
setAlterFullTextIndexAction
    : Set ( setStoplistAlterFulltextIndexAction | setChangeTrackingAlterFulltextIndexAction_CHANGE_TRACKING | setSearchPropertyListAlterFullTextIndexAction_SEARCH )
    ;

// TSql160.g:11599
setStoplistAlterFulltextIndexAction
    : stoplistFulltextIndexOption populationOption
    ;

// TSql160.g:11664
addAlterFulltextIndexAction
    : Add LeftParenthesis fulltextIndexColumn ( Comma fulltextIndexColumn )* RightParenthesis populationOption
    ;

// TSql160.g:11688
populationOption
    : ( With 'NO' 'POPULATION' )?
    ;

// TSql160.g:11702
dropAlterFulltextIndexAction
    : Drop LeftParenthesis identifier ( Comma identifier )* RightParenthesis populationOption
    ;

// TSql160.g:11752
otherPopulationAlterFulltextIndexAction
    : ( 'STOP' | 'PAUSE' | 'RESUME' ) 'POPULATION'
    ;

// TSql160.g:11769
alterColumnAlterFullTextIndexAction
    : Alter Column alterFullTextIndexColumn populationOption
    ;

// TSql160.g:11788
alterFullTextIndexColumn
    : identifier ( Add | Drop ) 'STATISTICAL_SEMANTICS'
    ;

// TSql160.g:11812
// TSql160.g:11843: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Disable)}?
// TSql160.g:11850: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Pause)}?
// TSql160.g:11857: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Abort)}?
alterIndexStatement
    : Index ( identifier | All ) On schemaObjectFourPartName ( Set indexOptionList | 'DISABLE' | 'PAUSE' | 'ABORT' | Identifier ( partitionSpecifier )? ( With indexOptionList )? | ( With xmlNamespaces )? For LeftParenthesis promotedSelectiveXmlIndexPathInAlter ( Comma promotedSelectiveXmlIndexPathInAlter )* RightParenthesis )
    ;

// TSql160.g:11916
partitionSpecifier
    : 'PARTITION' EqualsSign ( expression | All )
    ;

// TSql160.g:11955
alterLoginEnableDisable
    : Identifier
    ;

// TSql160.g:11966
alterLoginOptions
    : With alterLoginParam ( Comma alterLoginParam )*
    ;

// TSql160.g:11989
alterLoginAddDropCredential
    : ( Add | Drop ) 'CREDENTIAL' identifier
    ;

// TSql160.g:12016
alterLoginParam
    : ( ( Identifier EqualsSign ( passwordAlterLoginOption | onOffPrincipalOption | identifier ) ) | 'NO' 'CREDENTIAL' )
    ;

// TSql160.g:12055
passwordAlterLoginOption
    : loginPassword ( 'OLD_PASSWORD' EqualsSign stringLiteral | ( Identifier )* )
    ;

// TSql160.g:12230
alterPartitionFunctionStatement
    : Function identifier LeftParenthesis RightParenthesis ( 'SPLIT' | Merge ) ( 'RANGE' LeftParenthesis expression RightParenthesis )?
    ;

// TSql160.g:12263
alterPartitionSchemeStatement
    : 'SCHEME' identifier 'NEXT' 'USED' ( stringOrIdentifier )?
    ;

// TSql160.g:12288
// TSql160.g:12313: action may refuse the input: {vResult.ExecuteAsOption = ExecuteAsOptionHelper.Instance.ParseOption(tOption); if ((vCallerProhibited && vResult.ExecuteAsOption == ExecuteAsOption.Caller) || (vOwnerProhibited && vResult.ExecuteAsOption == ExecuteAsOption.Owner)) { throw GetUnexpect...
executeAsClause
    : ( Exec | Execute ) As ( stringLiteral | Identifier )
    ;

// TSql160.g:12324
queueOptionList
    : With queueOption ( Comma queueOption )*
    ;

// TSql160.g:12332
// TSql160.g:12341: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Activation)}?
queueOption
    : stateQueueOption
    | queueActivationOption_ACTIVATION
    | queuePoisonMessageHandlingOption
    ;

// TSql160.g:12350
stateQueueOption
    : ( 'STATUS' | 'RETENTION' ) EqualsSign optionOnOff
    ;

// TSql160.g:12401
activationQueueOptionArgument
    : ( Identifier EqualsSign ( statusQueueOptionArgument | procedureNameQueueOptionArgument | maxQueueReadersQueueOptionArgument ) )
    | dropQueueOptionArgument
    | executeAsQueueOptionArgument
    ;

// TSql160.g:12427
statusQueueOptionArgument
    : optionOnOff
    ;

// TSql160.g:12440
procedureNameQueueOptionArgument
    : schemaObjectThreePartName
    ;

// TSql160.g:12453
maxQueueReadersQueueOptionArgument
    : integer
    ;

// TSql160.g:12466
// TSql160.g:12469: action may refuse the input: {if (!vDropAccepted) throw GetUnexpectedTokenErrorException(tDrop); UpdateTokenInfo(vResult,tDrop); vResult.OptionKind = QueueOptionKind.ActivationDrop;}
dropQueueOptionArgument
    : Drop
    ;

// TSql160.g:12477
executeAsQueueOptionArgument
    : executeAsClause
    ;

// TSql160.g:12488
queuePoisonMessageHandlingOption
    : 'POISON_MESSAGE_HANDLING' LeftParenthesis 'STATUS' EqualsSign optionOnOff RightParenthesis
    ;

// TSql160.g:12553
alterRoleAction
    : addMemberAlterRoleAction
    | dropMemberAlterRoleAction
    | renameAlterRoleAction
    ;

// TSql160.g:12560
renameAlterRoleAction
    : With 'NAME' EqualsSign identifier
    ;

// TSql160.g:12572
addMemberAlterRoleAction
    : Add 'MEMBER' identifier
    ;

// TSql160.g:12584
dropMemberAlterRoleAction
    : Drop 'MEMBER' identifier
    ;

// TSql160.g:12596
routeOptionList
    : With routeOption ( Comma routeOption )*
    ;

// TSql160.g:12615
// TSql160.g:12620: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.LifeTime)}?
routeOption
    : 'LIFETIME' EqualsSign integer
    | Identifier EqualsSign stringLiteral
    ;

// TSql160.g:12650
alterSchemaStatement
    : Schema identifier 'TRANSFER' ( ( 'OBJECT' | 'TYPE' ) DoubleColon | 'XML' Schema 'COLLECTION' DoubleColon )? schemaObjectTwoPartName
    ;

// TSql160.g:12705
alterSequenceOptionList
    : ( alterSequenceOptionListElement )+
    ;

// TSql160.g:12732
alterServiceMasterKeyStatement
    : 'MASTER' Key ( 'FORCE' 'REGENERATE' | 'REGENERATE' | With ( 'OLD_ACCOUNT' | 'NEW_ACCOUNT' ) EqualsSign stringLiteral Comma ( 'OLD_PASSWORD' | 'NEW_PASSWORD' ) EqualsSign stringLiteral )
    ;

// TSql160.g:12776
// TSql160.g:12779: predicate dropped: {TryMatch(tOption, CodeGenerationSupporter.Password)}?
alterUserOption
    : Identifier EqualsSign ( passwordAlterUserOption | identifierAlterUserOption | literalAlterUserOption )
    ;

// TSql160.g:12786
identifierAlterUserOption
    : identifier
    ;

// TSql160.g:12798
literalAlterUserOption
    : integer
    | nullLiteral
    ;

// TSql160.g:12817
passwordAlterUserOption
    : loginPassword ( 'OLD_PASSWORD' EqualsSign stringLiteral )?
    ;

// TSql160.g:12837
alterUserStatement
    : User identifier With alterUserOption ( Comma alterUserOption )*
    ;

// TSql160.g:12878
createSchemaStatement
    : Create Schema ( identifier authorizationOpt | authorization ) createSchemaElementList
    ;

// TSql160.g:12907
createSchemaElementList
    : ( createSchemaElement )*
    ;

// TSql160.g:12920
createSchemaElement
    : ( createViewStatement | createTableStatement | grantStatement90 | denyStatement90 | revokeStatement90 ) optSemicolons
    ;

// TSql160.g:12936
createFunctionStatement
    : Create functionStatementBody
    ;

// TSql160.g:12953
alterFunctionStatement
    : Alter functionStatementBody
    ;

// TSql160.g:12971
createOrAlterFunctionStatement
    : functionStatementBody
    ;

// TSql160.g:12984
functionStatementBody
    : Function schemaObjectThreePartName LeftParenthesis ( functionParameterList )? RightParenthesis 'RETURNS' functionReturnTypeAndBody
    ;

// TSql160.g:13002
functionParameterList
    : functionParameter ( Comma functionParameter )*
    ;

// TSql160.g:13018
functionParameter
    : identifierVariable ( As )? scalarProcedureParameter
    ;

// TSql160.g:13029
functionReturnTypeAndBody
    : scalarDataType ( functionAttributes )? ( As )? ( beginEndBlockStatement | methodSpecifier optSemicolons )
    | Table ( functionAttributesNoExecuteAs )? ( As )? Return functionReturnClauseRelational
    | declareTableBody ( functionAttributes )? ( As )? beginEndBlockStatement
    | Table declareTableBodyMain ( functionAttributes )? ( clrTableValuedFunctionOrderHint )? ( As )? methodSpecifier optSemicolons
    ;

// TSql160.g:13094
// TSql160.g:13105: action may refuse the input: {UpdateTokenInfo(vResult,tLParen); // check for invalid combination of CHANGE_TRACKING_CONTEXT and Select statement if ((vWithCTEAndXmlNamespaces != null) && (vWithCTEAndXmlNamespaces.ChangeTrackingContext != null)) ThrowParseErrorException("SQL46072"...
functionReturnClauseRelational
    : subqueryExpressionWithOptionalCTE
    | LeftParenthesis withCommonTableExpressionsAndXmlNamespaces subqueryExpression RightParenthesis
    ;

// TSql160.g:13118
functionAttributes
    : With functionAttribute ( Comma functionAttribute )*
    ;

// TSql160.g:13137
functionAttributesNoExecuteAs
    : With functionAttributeNoExecuteAs ( Comma functionAttributeNoExecuteAs )*
    ;

// TSql160.g:13156
// TSql160.g:13177: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Inline)}?
functionAttributeNoExecuteAs
    : Identifier
    | 'RETURNS' Null On Null 'INPUT'
    | 'CALLED' On Null 'INPUT'
    | inlineFunctionOption_INLINE
    ;

// TSql160.g:13203
functionAttribute
    : functionAttributeNoExecuteAs
    | functionExecuteAsOption
    ;

// TSql160.g:13208
functionExecuteAsOption
    : executeAsClause
    ;

// TSql160.g:13220
clrTableValuedFunctionOrderHint
    : bulkInsertSortOrderOption
    ;

// TSql160.g:13230
createStatisticsStatement
    : Create Statistics identifier On schemaObjectThreePartName identifierColumnList ( filterClause )? ( With createStatisticsStatementWithOption ( Comma createStatisticsStatementWithOption )* )?
    ;

// TSql160.g:13271
// TSql160.g:13277: action may refuse the input: {Match(tSample, CodeGenerationSupporter.Sample); if (isConflictingOption == true) ThrowParseErrorException("SQL46071", tSample, TSqlParserResource.SQL46071Message); else isConflictingOption = true; vResult.Literal = vLiteral;}
sampleStatisticsOption
    : 'SAMPLE' integer ( Identifier | Percent )
    ;

// TSql160.g:13299
statsStreamStatisticsOption
    : 'STATS_STREAM' EqualsSign binary
    ;

// TSql160.g:13312
// TSql160.g:13315: action may refuse the input: {if (TryMatch(tOption, CodeGenerationSupporter.FullScan)) { if (isConflictingOption == true) ThrowParseErrorException("SQL46071", tOption, TSqlParserResource.SQL46071Message); else isConflictingOption = true; } UpdateTokenInfo(vResult,tOption); vResul...
simpleStatisticsOption
    : 'FULLSCAN'
    ;

// TSql160.g:13328
statisticsPartitionRange
    : integer ( To integer )?
    ;

// TSql160.g:13412
// TSql160.g:13414: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Incremental)}?
// TSql160.g:13417: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.AutoDrop)}?
// TSql160.g:13420: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.PersistSamplePercent)}?
createStatisticsStatementWithOption
    : incrementalStatisticsOption_INCREMENTAL
    | autoDropStatisticsOption_AUTO_DROP
    | persistSamplePercentStatisticsOption_PERSIST_SAMPLE_PERCENT
    | sampleStatisticsOption
    | statsStreamStatisticsOption
    | simpleStatisticsOption
    ;

// TSql160.g:13430
updateStatisticsStatement
    : Update Statistics schemaObjectThreePartName ( columnNameList | identifier )? ( With updateStatisticsStatementWithOption ( Comma updateStatisticsStatementWithOption )* )?
    ;

// TSql160.g:13472
// TSql160.g:13474: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Resample)}?
// TSql160.g:13476: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Incremental)}?
// TSql160.g:13478: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.AutoDrop)}?
// TSql160.g:13480: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.PersistSamplePercent)}?
// TSql160.g:13484: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.StatsStream)}?
updateStatisticsStatementWithOption
    : resampleStatisticsOption_RESAMPLE
    | incrementalStatisticsOption_INCREMENTAL
    | autoDropStatisticsOption_AUTO_DROP
    | persistSamplePercentStatisticsOption_PERSIST_SAMPLE_PERCENT
    | sampleStatisticsOption
    | statsStreamStatisticsOption_STATS_STREAM
    | updateStatisticsLiteralOption
    | updateStatisticsSimpleOption
    ;

// TSql160.g:13492
updateStatisticsLiteralOption
    : 'PAGECOUNT' EqualsSign integerOrRealOrNumeric
    | RowCount EqualsSign integerOrRealOrNumeric
    ;

// TSql160.g:13513
// TSql160.g:13516: action may refuse the input: {if (TryMatch(tOption, CodeGenerationSupporter.FullScan)) { if (isConflictingOption) ThrowParseErrorException("SQL46071", tOption, TSqlParserResource.SQL46071Message); else isConflictingOption = true; } UpdateTokenInfo(vResult,tOption); vResult.Option...
updateStatisticsSimpleOption
    : 'FULLSCAN'
    | All
    | Index
    ;

// TSql160.g:13541
ifStatement
    : If booleanExpression ( statementOptSemi ) ( Else ( statementOptSemi ) )?
    ;

// TSql160.g:13594
whileStatement
    : While booleanExpression ( statementOptSemi )
    ;

// TSql160.g:13623
lineNoStatement
    : LineNo integer
    ;

// TSql160.g:13635
useStatement
    : Use identifier
    ;

// TSql160.g:13647
useFederationStatement
    : Use 'FEDERATION' ( 'ROOT' With | identifier LeftParenthesis identifier EqualsSign expression RightParenthesis With 'FILTERING' EqualsSign ( On | Off ) Comma ) 'RESET'
    ;

// TSql160.g:13706
// TSql160.g:13712: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Query)}?
killStatements
    : Kill ( killStatement | killQueryNotificationSubscriptionStatement_QUERY | killStatsJobStatement )
    ;

// TSql160.g:13723
killStatement
    : ( signedInteger | stringLiteral ) ( With 'STATUSONLY' )?
    ;

// TSql160.g:13779
killStatsJobStatement
    : 'STATS' 'JOB' signedInteger
    ;

// TSql160.g:13798
checkpointStatement
    : Checkpoint ( integer )?
    ;

// TSql160.g:13815
reconfigureStatement
    : Reconfigure ( With 'OVERRIDE' )?
    ;

// TSql160.g:13833
shutdownStatement
    : Shutdown ( With 'NOWAIT' )?
    ;

// TSql160.g:13851
setUserStatement
    : SetUser ( stringOrVariable ( With 'NORESET' )? )?
    ;

// TSql160.g:13878
truncateTableStatement
    : Truncate Table schemaObjectThreePartName ( With LeftParenthesis Identifier LeftParenthesis compressionPartitionRange ( Comma compressionPartitionRange )* RightParenthesis RightParenthesis )?
    ;

// TSql160.g:13904
// TSql160.g:13911: action may refuse the input: {if (vResult.Identifiers.Count >= 7) { throw GetUnexpectedTokenErrorException(vIdentifier); } AddAndUpdateTokenInfo(vResult, vResult.Identifiers, vIdentifier);}
permission
    : ( securityStatementPermission )+ ( columnNameList )?
    ;

// TSql160.g:13924
securityTargetObjectCommon
    : On ( ( securityStatementPermission | securityStatementPermission ( securityStatementPermission | securityStatementPermission securityStatementPermission ) ) DoubleColon )? multiPartIdentifier
    ;

// TSql160.g:13967
// TSql160.g:13972: action may refuse the input: {if(vColumnsDisallowed) { ThrowIncorrectSyntaxErrorException(vResult.Columns[0]); }}
securityTargetObject
    : securityTargetObjectCommon ( columnNameList )?
    ;

// TSql160.g:13981
authorizationTargetObject
    : securityTargetObjectCommon
    ;

// TSql160.g:13986
principal
    : Public
    | Null
    | identifier
    ;

// TSql160.g:14010
alterAuthorizationStatement
    : Authorization authorizationTargetObject To ( Schema 'OWNER' | identifier )
    ;

// TSql160.g:14036
grantStatement90
    : Grant permissionsList securityStatementTargetObjectOpt To securityStatementPrincipalList ( With Grant Option )? securityStatementAsClauseOpt
    ;

// TSql160.g:14059
denyStatement90
    : Deny permissionsList securityStatementTargetObjectOpt To securityStatementPrincipalList ( Cascade )? securityStatementAsClauseOpt
    ;

// TSql160.g:14080
revokeStatement90
    : Revoke ( Grant Option For )? permissionsList securityStatementTargetObjectOpt ( To | From ) securityStatementPrincipalList ( Cascade )? securityStatementAsClauseOpt
    ;

// TSql160.g:14107
securityStatementPrincipalList
    : principal ( Comma principal )*
    ;

// TSql160.g:14123
permissionsList
    : permission ( Comma permission )*
    ;

// TSql160.g:14141
securityStatementTargetObjectOpt
    : ( securityTargetObject )?
    ;

// TSql160.g:14152
securityStatementAsClauseOpt
    : ( As identifier )?
    ;

// TSql160.g:14163
printStatement
    : Print expression
    ;

// TSql160.g:14175
waitForStatement
    : WaitFor ( Identifier stringOrVariable | LeftParenthesis waitforInnerStatement RightParenthesis ( Comma 'TIMEOUT' signedIntegerOrVariable )? )
    ;

// TSql160.g:14207
// TSql160.g:14208: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Receive)}?
waitforInnerStatement
    : receiveStatement_RECEIVE
    | getConversationGroupStatement
    ;

// TSql160.g:14213
readTextStatement
    : ReadText column binaryOrVariable integerOrVariable integerOrVariable ( HoldLock )?
    ;

// TSql160.g:14249
updateTextStatement
    : UpdateText modificationTextStatement signedIntegerOrVariableOrNull signedIntegerOrVariableOrNull ( modificationTextStatementWithLog )? ( column binaryOrVariable | writeString )?
    ;

// TSql160.g:14294
writeTextStatement
    : WriteText modificationTextStatement ( modificationTextStatementWithLog )? writeString
    ;

// TSql160.g:14313
modificationTextStatementWithLog
    : With 'LOG'
    ;

// TSql160.g:14323
modificationTextStatement
    : ( Bulk )? column ( binaryOrVariable | integer ) ( 'TIMESTAMP' EqualsSign binary )?
    ;

// TSql160.g:14361
writeString
    : nullLiteral
    | stringLiteral
    | binary
    | variable
    ;

// TSql160.g:14372
returnStatement
    : Return ( expression |  )
    ;

// TSql160.g:14390
// TSql160.g:14394: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Master)}?
// TSql160.g:14397: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Symmetric)}?
openStatements
    : Open ( openMasterKeyStatement_MASTER | openSymmetricKeyStatement_SYMMETRIC | openCursorStatement )
    ;

// TSql160.g:14430
// TSql160.g:14444: action may refuse the input: {if (vResult.CryptoMechanismType == CryptoMechanismType.SymmetricKey) throw GetUnexpectedTokenErrorException(tWith);}
decryptionMechanism
    : ( passwordCrypto | certificateCrypto ( With decryptionMechanismPassword )? | keyCrypto ( With decryptionMechanismPassword )? )
    ;

// TSql160.g:14452
decryptionMechanismPassword
    : 'PASSWORD' EqualsSign stringLiteral
    ;

// TSql160.g:14487
openCursorStatement
    : cursorId
    ;

// TSql160.g:14498
// TSql160.g:14505: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.Symmetric) || LA(1) == All}?
// TSql160.g:14502: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Master)}?
closeStatements
    : Close ( closeMasterKeyStatement_MASTER | closeSymmetricKeyStatement | closeCursorStatement )
    ;

// TSql160.g:14527
closeSymmetricKeyStatement
    : ( 'SYMMETRIC' Key identifier | All 'SYMMETRIC' 'KEYS' )
    ;

// TSql160.g:14555
closeCursorStatement
    : cursorId
    ;

// TSql160.g:14566
deallocateCursorStatement
    : Deallocate cursorId
    ;

// TSql160.g:14577
fetchCursorStatement
    : Fetch rowSelector ( Into variable ( Comma variable )* )?
    ;

// TSql160.g:14597
rowSelector
    : ( cursorId )
    | ( fetchType From cursorId )
    | ( From cursorId )
    ;

// TSql160.g:14620
// TSql160.g:14630: action may refuse the input: {if (vResult.Orientation != FetchOrientation.Relative && vResult.Orientation != FetchOrientation.Absolute) throw GetUnexpectedTokenErrorException(tFetchOrientation); vResult.RowOffset = vOffset;}
// TSql160.g:14638: action may refuse the input: {if (vResult.Orientation != FetchOrientation.Relative && vResult.Orientation != FetchOrientation.Absolute) throw GetUnexpectedTokenErrorException(tFetchOrientation); vResult.RowOffset = vOffset;}
// TSql160.g:14645: action may refuse the input: {if (vResult.Orientation == FetchOrientation.Relative || vResult.Orientation == FetchOrientation.Absolute) throw GetUnexpectedTokenErrorException(tFetchOrientation);}
fetchType
    : Identifier ( ( signedInteger ) | ( variable ) |  )
    ;

// TSql160.g:14652
// TSql160.g:14657: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Global)}?
cursorId
    : ( 'GLOBAL' identifier )
    | identifierOrVariable
    ;

// TSql160.g:14675
// TSql160.g:14678: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Sensitivity)}?
addStatements
    : Add ( addSensitivityClassificationStatement_SENSITIVITY | addSignatureStatement )
    ;

// TSql160.g:14690
// TSql160.g:14703: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Aggregate)}?
// TSql160.g:14705: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Assembly)}?
// TSql160.g:14707: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Availability)}?
// TSql160.g:14709: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Partition)}?
// TSql160.g:14711: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Synonym)}?
// TSql160.g:14713: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Application)}?
// TSql160.g:14715: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Broker)}?
// TSql160.g:14717: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Cryptographic)}?
// TSql160.g:14719: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Fulltext)}?
// TSql160.g:14721: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Login)}?
// TSql160.g:14723: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Resource)}?
// TSql160.g:14725: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Workload)}?
// TSql160.g:14727: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Role)}?
// TSql160.g:14729: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Master)}?
// TSql160.g:14731: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Symmetric)}?
// TSql160.g:14733: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Asymmetric)}?
// TSql160.g:14735: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Certificate)}?
// TSql160.g:14737: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Credential)}?
// TSql160.g:14739: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Type)}?
// TSql160.g:14741: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Xml)}?
// TSql160.g:14743: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Contract)}?
// TSql160.g:14745: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Queue)}?
// TSql160.g:14747: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Service)}?
// TSql160.g:14749: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Route)}?
// TSql160.g:14751: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Message)}?
// TSql160.g:14753: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Remote)}?
// TSql160.g:14755: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Endpoint)}?
// TSql160.g:14757: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.Signature) || NextTokenMatches(CodeGenerationSupporter.Counter)}?
// TSql160.g:14759: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Event)}?
// TSql160.g:14761: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Search)}?
// TSql160.g:14763: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Sequence)}?
// TSql160.g:14765: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Federation)}?
// TSql160.g:14767: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Security)}?
// TSql160.g:14769: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.Column)}?
// TSql160.g:14771: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.External)}?
// TSql160.g:14773: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Sensitivity)}?
dropStatements
    : Drop ( dropDatabaseStatements | dropIndexStatement | dropStatisticsStatement | dropTableStatement | dropProcedureStatement | dropFunctionStatement | dropViewStatement | dropDefaultStatement | dropRuleStatement | dropTriggerStatement | dropSchemaStatement | dropAggregateStatement_AGGREGATE | dropAssemblyStatement_ASSEMBLY | dropAvailabilityGroupStatement_AVAILABILITY | dropPartitionStatements_PARTITION | dropSynonymStatement_SYNONYM | dropApplicationRoleStatement_APPLICATION | dropBrokerPriorityStatement_BROKER | dropCryptographicProviderStatement_CRYPTOGRAPHIC | dropFulltextStatements_FULLTEXT | dropLoginStatement_LOGIN | dropResourcePoolStatement_RESOURCE | dropWorkloadStatements_WORKLOAD | dropRoleStatement_ROLE | dropMasterKeyStatement_MASTER | dropSymmetricKeyStatement_SYMMETRIC | dropAsymmetricKeyStatement_ASYMMETRIC | dropCertificateStatement_CERTIFICATE | dropCredentialStatement_CREDENTIAL | dropTypeStatement_TYPE | dropXmlSchemaCollectionStatement_XML | dropContractStatement_CONTRACT | dropQueueStatement_QUEUE | dropServiceStatement_SERVICE | dropRouteStatement_ROUTE | dropMessageTypeStatement_MESSAGE | dropRemoteServiceBindingStatement_REMOTE | dropEndpointStatement_ENDPOINT | dropSignatureStatement | dropEventStatement_EVENT | dropSearchPropertyListStatement_SEARCH | dropSequenceStatement_SEQUENCE | dropFederationStatement_FEDERATION | dropSecurityPolicyStatement_SECURITY | dropColumnStatements | dropExternalStatement | dropSensitivityClassificationStatement_SENSITIVITY | dropServerStatements | dropUserStatement )
    ;

// TSql160.g:14784
// TSql160.g:14810: action may refuse the input: {// Schema name with more than one part must be followed by either Restrict or Cascade (since it is SQL 2000 syntax) if (vResult.Schema.SchemaIdentifier != null && vResult.DropBehavior == DropSchemaBehavior.None) throw GetUnexpectedTokenErrorException...
dropSchemaStatement
    : Schema ( If Exists )? schemaObjectThreePartName ( Cascade | Restrict )?
    ;

// TSql160.g:14817
// TSql160.g:14827: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Specification)}?
// TSql160.g:14823: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Audit)}?
dropServerStatements
    : 'SERVER' ( ( 'AUDIT' ( dropServerAuditSpecificationStatement_SPECIFICATION | dropServerAuditStatement ) ) | dropServerRoleStatement )
    ;

// TSql160.g:14837
dropServerRoleStatement
    : 'ROLE' identifier
    ;

// TSql160.g:14859
dropServerAuditStatement
    : identifier
    ;

// TSql160.g:14869
// TSql160.g:14872: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Audit)}?
// TSql160.g:14875: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Scoped) && NextTokenMatches(CodeGenerationSupporter.Credential, 2)}?
dropDatabaseStatements
    : Database ( dropDatabaseAuditSpecificationStatement_AUDIT | dropDatabaseScopedCredentialStatement_SCOPED | dropDatabaseEncryptionKeyStatement | dropDatabaseStatement )
    ;

// TSql160.g:14884
dropDatabaseStatement
    : ( If Exists )? identifier ( Comma identifier )*
    ;

// TSql160.g:14918
dropDatabaseEncryptionKeyStatement
    : 'ENCRYPTION' Key
    ;

// TSql160.g:14926
dropIndexStatement
    : Index ( If Exists )? dropIndexClause ( Comma dropIndexClause )*
    ;

// TSql160.g:14948
dropIndexClause
    : indexDropObject
    | indexDropObjectNewNameFormat
    ;

// TSql160.g:14953
indexDropObjectNewNameFormat
    : identifier On schemaObjectThreePartName ( dropClusteredIndexOptions )?
    ;

// TSql160.g:14967
dropClusteredIndexOptions
    : With LeftParenthesis dropClusteredIndexOption ( Comma dropClusteredIndexOption )* RightParenthesis
    ;

// TSql160.g:14977
// TSql160.g:14982: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MaxDop)}?
// TSql160.g:14985: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Online)}?
// TSql160.g:14988: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Move)}?
// TSql160.g:14991: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.FileStreamOn)}?
// TSql160.g:14994: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.DataCompression)}?
// TSql160.g:14997: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.XmlCompression)}?
// TSql160.g:15000: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.WaitAtLowPriority)}?
dropClusteredIndexOption
    : ( maxDopOption_MAXDOP | onlineIndexOption_ONLINE | dropIndexMoveToOption_MOVE | dropIndexFileStreamOnOption_FILESTREAM_ON | dataCompressionOption_DATA_COMPRESSION | xmlCompressionOption_XML_COMPRESSION | waitAtLowPriorityOption_WAIT_AT_LOW_PRIORITY )
    ;

// TSql160.g:15030
dropStatisticsStatement
    : Statistics statisticsDropObject ( Comma statisticsDropObject )*
    ;

// TSql160.g:15045
// TSql160.g:15047: action may refuse the input: {if (vResult.BaseIdentifier == null) ThrowParseErrorException("SQL46038", vResult, TSqlParserResource.SQL46038Message);}
statisticsDropObject
    : childObjectNameWithThreePrefixes
    ;

// TSql160.g:15053
childObjectNameWithThreePrefixes
    : identifierList
    ;

// TSql160.g:15063
dropTableStatement
    : Table ( If Exists )? dropObjectList
    ;

// TSql160.g:15074
dropProcedureStatement
    : ( Procedure | Proc ) ( If Exists )? dropObjectList
    ;

// TSql160.g:15085
dropFunctionStatement
    : Function ( If Exists )? dropObjectList
    ;

// TSql160.g:15096
dropViewStatement
    : View ( If Exists )? dropObjectList
    ;

// TSql160.g:15107
dropDefaultStatement
    : Default ( If Exists )? dropObjectList
    ;

// TSql160.g:15118
dropRuleStatement
    : Rule ( If Exists )? dropObjectList
    ;

// TSql160.g:15129
dropTriggerStatement
    : Trigger ( If Exists )? dropObjectList ( On triggerScope )?
    ;

// TSql160.g:15149
dropExternalTableStatement
    : Table dropObjectList
    ;

// TSql160.g:15153
dropObjectList
    : dropObject ( Comma dropObject )*
    ;

// TSql160.g:15168
dropObject
    : schemaObjectThreePartName
    ;

// TSql160.g:15176
// TSql160.g:15181: action may refuse the input: {if (vObjName.BaseIdentifier == null) ThrowParseErrorException("SQL46027", vObjName, TSqlParserResource.SQL46027Message); vResult.Index = vObjName;}
indexDropObject
    : childObjectNameWithThreePrefixes
    ;

// TSql160.g:15200
dropPartitionFunctionStatement
    : Function identifier
    ;

// TSql160.g:15259
dropPartitionSchemeStatement
    : 'SCHEME' identifier
    ;

// TSql160.g:15310
dropFulltextCatalogStatement
    : 'CATALOG' identifier
    ;

// TSql160.g:15321
dropFulltextIndexStatement
    : Index On dropObject
    ;

// TSql160.g:15400
removeProviderKeyOpt
    : ( 'REMOVE' 'PROVIDER' Key )?
    ;

// TSql160.g:15482
dropUserStatement
    : User ( If Exists )? identifier
    ;

// TSql160.g:15615
eventSessionDropScope
    : On ( Database | 'SERVER' )
    ;

// TSql160.g:15646
// TSql160.g:15649: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.Table)}?
// TSql160.g:15652: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Data)}?
// TSql160.g:15655: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.File)}?
// TSql160.g:15658: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Library)}?
// TSql160.g:15661: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Language)}?
// TSql160.g:15664: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Resource)}?
dropExternalStatement
    : External ( dropExternalTableStatement | dropExternalDataSourceStatement_DATA | dropExternalFileFormatStatement | dropExternalLibraryStatement_LIBRARY | dropExternalLanguageStatement_LANGUAGE | dropExternalResourcePoolStatement_RESOURCE )
    ;

// TSql160.g:15669
// TSql160.g:15672: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Master)}?
// TSql160.g:15675: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Encryption)}?
dropColumnStatements
    : Column ( dropColumnMasterKeyStatement_MASTER | dropColumnEncryptionKeyStatement_ENCRYPTION )
    ;

// TSql160.g:15747
addSignatureStatement
    : signatureType To signableElement By cryptoListWithOptionalPasswordSignature
    ;

// TSql160.g:15751
dropSignatureStatement
    : signatureType From signableElement By cryptoListWithOptionalPasswordSignature
    ;

// TSql160.g:15755
signatureType
    : ( 'COUNTER' )? 'SIGNATURE'
    ;

// TSql160.g:15768
signableElement
    : schemaObjectThreePartName
    | ( 'OBJECT' | 'ASSEMBLY' ) DoubleColon schemaObjectThreePartName
    | Database DoubleColon schemaObjectThreePartName
    ;

// TSql160.g:15794
cryptoListWithOptionalPasswordSignature
    : cryptoWithOptionalPasswordSignature ( Comma cryptoWithOptionalPasswordSignature )*
    ;

// TSql160.g:15809
cryptoWithOptionalPasswordSignature
    : certificateCrypto withSignatureOrPasswordOpt
    | keyCrypto withSignatureOrPasswordOpt
    | passwordCrypto
    ;

// TSql160.g:15826
withSignatureOrPasswordOpt
    : ( With Identifier EqualsSign ( stringLiteral | binary ) )?
    ;

// TSql160.g:15839
labelStatement
    : Label
    ;

// TSql160.g:15847
gotoStatement
    : GoTo nonQuotedIdentifier
    ;

// TSql160.g:15858
// TSql160.g:15863: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.Try, 2)}?
// TSql160.g:15865: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.Conversation, 2)}?
// TSql160.g:15867: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.Dialog, 2)}?
beginStatements
    : beginTransactionStatement
    | tryCatchStatement
    | beginConversationTimerStatement
    | beginDialogStatement
    | beginEndBlockStatement
    ;

// TSql160.g:15873
saveTransactionStatement
    : Save ( Tran | Transaction ) transactionName
    ;

// TSql160.g:15883
rollbackTransactionStatement
    : Rollback ( 'WORK' | ( ( Tran | Transaction ) ( transactionName )? ) )?
    ;

// TSql160.g:15917
commitTransactionStatement
    : Commit ( 'WORK' | ( ( Tran | Transaction ) ( transactionName )? ) )? ( With LeftParenthesis 'DELAYED_DURABILITY' EqualsSign optionOnOff RightParenthesis )?
    ;

// TSql160.g:15961
beginTransactionStatement
    : Begin ( Distributed )? ( Tran | Transaction ) ( transactionName )? ( With 'MARK' ( stringOrVariable )? )?
    ;

// TSql160.g:16007
transactionName
    : identifierOrVariable
    | weirdTransactionName
    ;

// TSql160.g:16022
weirdTransactionName
    : ( Minus )? Integer Colon tranIdentifier Dot tranIdentifier
    ;

// TSql160.g:16051
tranIdentifier
    : Identifier
    | QuotedIdentifier
    ;

// TSql160.g:16066
tryCatchStatement
    : Begin 'TRY' tryStatementList End 'TRY' Begin 'CATCH' catchStatementList End 'CATCH'
    ;

// TSql160.g:16103
// TSql160.g:16110: predicate dropped: {IsStatementIsNext()}?
tryStatementList
    : ( Semicolon )* ( statementOptSemi )+
    ;

// TSql160.g:16125
// TSql160.g:16132: predicate dropped: {IsStatementIsNext()}?
catchStatementList
    : ( Semicolon )* ( statementOptSemi )*
    ;

// TSql160.g:16146
// TSql160.g:16168: predicate dropped: {IsStatementIsNext()}?
beginEndBlockStatement
    : Begin ( beginEndAtomicBlock )? ( Semicolon )* ( statementOptSemi )+ End
    ;

// TSql160.g:16190
beginEndAtomicBlock
    : 'ATOMIC' atomicBlockOptions
    ;

// TSql160.g:16198
atomicBlockOptions
    : With LeftParenthesis atomicBlockOption ( Comma atomicBlockOption )* RightParenthesis
    ;

// TSql160.g:16221
// TSql160.g:16223: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.DateFirst)}?
// TSql160.g:16226: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.DateFormat)}?
// TSql160.g:16229: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.DelayedDurability)}?
// TSql160.g:16232: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Language)}?
atomicBlockOption
    : dateFirstOption_DATEFIRST
    | dateFormatOption_DATEFORMAT
    | delayedDurabilityOption_DELAYED_DURABILITY
    | languageOption_LANGUAGE
    | isolationLevelIdentifierAtomicBlockOption
    ;

// TSql160.g:16274
isolationLevelIdentifierAtomicBlockOption
    : Transaction 'ISOLATION' 'LEVEL' EqualsSign isolationLevel
    ;

// TSql160.g:16319
breakStatement
    : Break
    ;

// TSql160.g:16327
continueStatement
    : Continue
    ;

// TSql160.g:16335
declareVariableElement
    : identifierVariable ( As )? ( scalarDataType | cursorDataType ) ( nullableConstraint )? ( EqualsSign expression )?
    ;

// TSql160.g:16366
declareStatements
    : Declare ( declareTableBody | declareVariableStatement | declareCursorStatement )
    ;

// TSql160.g:16386
// TSql160.g:16389: predicate dropped: {!NextTokenMatches(CodeGenerationSupporter.FipsFlagger)}?
setStatements
    : Set ( predicateSetStatement | setVariableStatement | setStatisticsStatement | setRowcountStatement | setOffsetsStatement | setCommandStatement | setTransactionIsolationLevelStatement | setTextSizeStatement | setIdentityInsertStatement | setErrorLevelStatement )
    ;

// TSql160.g:16406
setVariableStatement
    : variable ( ( Dot | DoubleColon ) identifier ( parenthesizedOptExpressionWithDefaultList | assignmentWithOptOp expression ) | ( assignmentWithOptOp expression | EqualsSign cursorDefinition ) )
    ;

// TSql160.g:16470
predicateSetStatement
    : Identifier ( Comma Identifier )* setOnOff
    ;

// TSql160.g:16494
setRowcountStatement
    : RowCount integerOrVariable
    ;

// TSql160.g:16504
setOffsetsStatement
    : Offsets offsetItem ( Comma offsetItem )* setOnOff
    ;

// TSql160.g:16520
offsetItem
    : Select
    | From
    | Order
    | Compute
    | Table
    | ( Procedure | Proc )
    | ( Execute | Exec )
    | ( 'STATEMENT' | 'PARAM' )
    ;

// TSql160.g:16540
setCommandStatement
    : setCommand ( Comma setCommand )*
    ;

// TSql160.g:16555
// TSql160.g:16560: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.FipsFlagger)}?
setCommand
    : ( 'FIPS_FLAGGER' fipsFlaggerLevel )
    | ( Identifier possibleNegativeConstantOrIdentifier )
    ;

// TSql160.g:16573
fipsFlaggerLevel
    : Off
    | stringLiteral
    ;

// TSql160.g:16589
setTransactionIsolationLevelStatement
    : ( Transaction | Tran ) 'ISOLATION' 'LEVEL' isolationLevel
    ;

// TSql160.g:16601
isolationLevel
    : ( ( Read ( 'COMMITTED' | 'UNCOMMITTED' ) ) | ( 'REPEATABLE' Read ) | ( ( 'SNAPSHOT' | 'SERIALIZABLE' ) ) )
    ;

// TSql160.g:16640
setTextSizeStatement
    : TextSize signedInteger
    ;

// TSql160.g:16650
setIdentityInsertStatement
    : IdentityInsert schemaObjectThreePartName setOnOff
    ;

// TSql160.g:16661
setErrorLevelStatement
    : Errlvl signedInteger
    ;

// TSql160.g:16671
setStatisticsStatement
    : Statistics Identifier ( Comma Identifier )* setOnOff
    ;

// TSql160.g:16684
setOnOff
    : On
    | Off
    ;

// TSql160.g:16697
declareTableBody
    : identifierVariable ( As )? Table declareTableBodyMain
    ;

// TSql160.g:16716
declareTableBodyMain
    : LeftParenthesis tableDefinition RightParenthesis
    ;

// TSql160.g:16732
tableDefinition
    : ( tableElement ( Comma tableElement )* )
    ;

// TSql160.g:16755
declareVariableStatement
    : declareVariableElement ( Comma declareVariableElement )*
    ;

// TSql160.g:16770
declareCursorStatement
    : identifier cursorOpts cursorDefinitionOptions
    ;

// TSql160.g:16783
cursorDefinition
    : cursorDefinitionOptions
    ;

// TSql160.g:16790
cursorDefinitionOptions
    : Cursor cursorOpts For selectStatement
    ;

// TSql160.g:16801
// TSql160.g:16806: action may refuse the input: {// Only INSENSITIVE and SCROLL are allowed in old syntax, and INSENSITIVE is not allowed in new syntax if (oldSyntax) { if (vOption.OptionKind != CursorOptionKind.Insensitive && vOption.OptionKind != CursorOptionKind.Scroll) ThrowIncorrectSyntaxError...
cursorOpts
    : ( cursorOption )*
    ;

// TSql160.g:16823
cursorOption
    : Identifier
    ;

// TSql160.g:16831
createIndexStatement
    : Create ( Unique )? ( ( ( Clustered | NonClustered ) ( createRelationalIndexStatement | createColumnStoreIndexStatement ) ) | createRelationalIndexStatement )
    ;

// TSql160.g:16865
// TSql160.g:16887: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Include)}?
createRelationalIndexStatement
    : Index identifier On schemaObjectThreePartName LeftParenthesis indexKeyColumnList RightParenthesis ( 'INCLUDE' identifierColumnList )? ( filterClause )? ( With ( indexLegacyOptionList | indexOptionList ) )? ( On filegroupOrPartitionScheme )? fileStreamOnOpt
    ;

// TSql160.g:16918
// TSql160.g:16941: action may refuse the input: {UpdateTokenInfo(vResult, tRParen); if(isClustered.HasValue && isClustered==true && vResult.Columns.Count > 0) { ThrowIncorrectSyntaxErrorException(vResult.Columns[0]); }}
// TSql160.g:16953: action may refuse the input: {foreach (var col in vResult.OrderedColumns) { if (PseudoColumnHelper.IsGraphPseudoColumn(col)) { ThrowIncorrectSyntaxErrorException(col); } }}
// TSql160.g:16964: action may refuse the input: {if ((!isClustered.HasValue || isClustered==false) && vResult.Columns.Count == 0) { ThrowIncorrectSyntaxErrorException(vResult.OnName); }}
createColumnStoreIndexStatement
    : 'COLUMNSTORE' Index identifier On schemaObjectFourPartName ( LeftParenthesis columnstoreIndexKeyColumnList RightParenthesis )? ( Order ( identifierColumnList ) )? ( filterClause )? ( With ( indexOptionList ) )? ( On filegroupOrPartitionScheme )?
    ;

// TSql160.g:16990
indexKeyColumnList
    : indexKeyColumn ( Comma indexKeyColumn )*
    ;

// TSql160.g:17007
indexKeyColumn
    : columnWithSortOrder
    | ( graphPseudoColumn ( orderByOption )? )
    ;

// TSql160.g:17029
columnstoreIndexKeyColumnList
    : columnstoreIndexKeyColumn ( Comma columnstoreIndexKeyColumn )*
    ;

// TSql160.g:17046
columnstoreIndexKeyColumn
    : identifierColumnReferenceExpression
    | graphPseudoColumn
    ;

// TSql160.g:17052
// TSql160.g:17055: action may refuse the input: {switch (vResult.ColumnType) { case ColumnType.PseudoColumnGraphEdgeId: case ColumnType.PseudoColumnGraphNodeId: case ColumnType.PseudoColumnGraphFromId: case ColumnType.PseudoColumnGraphToId: break; default: ThrowIncorrectSyntaxErrorException(vResult...
graphPseudoColumn
    : systemColumn
    ;

// TSql160.g:17070
// TSql160.g:17072: action may refuse the input: {if (filterProhibited) { ThrowIncorrectSyntaxErrorException(tWhere); }}
filterClause
    : Where filterExpression
    ;

// TSql160.g:17080
filterExpressionPrimary
    : ( filterColumn ( Is ( filterNullPredicate | ( Not )? Distinct From ( filterNullPredicateFromDistinctPredicate | filterDistinctPredicate ) ) | filterComparisonPredicate | filterInPredicate ) | filterParenthesisExpression )
    ;

// TSql160.g:17113
filterExpression
    : filterExpressionPrimary ( And filterExpression )?
    ;

// TSql160.g:17129
filterColumn
    : ( specialColumn | multiPartIdentifier ( Dot specialColumn )? )
    ;

// TSql160.g:17149
filterNullPredicate
    : booleanExpressionNullPredicate
    ;

// TSql160.g:17153
filterDistinctPredicate
    : expression
    ;

// TSql160.g:17166
filterNullPredicateFromDistinctPredicate
    : Null
    ;

// TSql160.g:17176
filterComparisonPredicate
    : comparisonOperator expression
    ;

// TSql160.g:17190
filterInPredicate
    : In LeftParenthesis expression ( Comma expression )* RightParenthesis
    ;

// TSql160.g:17212
filterParenthesisExpression
    : LeftParenthesis filterExpression RightParenthesis
    ;

// TSql160.g:17224
filegroupOrPartitionScheme
    : stringOrIdentifier ( columnNameList )?
    ;

// TSql160.g:17238
indexLegacyOptionList
    : indexLegacyOption ( Comma indexLegacyOption )*
    ;

// TSql160.g:17256
fillFactorOption
    : FillFactor EqualsSign integer
    ;

// TSql160.g:17269
indexLegacyOption
    : fillFactorOption
    | Identifier
    ;

// TSql160.g:17282
// TSql160.g:17291: action may refuse the input: {VerifyAllowedIndexOption160(statementType, vIndexOption); // order option not supported in NCCI now if (statementType == IndexAffectingStatement.CreateColumnStoreIndex) { vStatement = (CreateColumnStoreIndexStatement)vParent; vIsNcci = vStatement.Clu...
indexOptionList
    : LeftParenthesis indexOption ( Comma indexOption )* RightParenthesis
    ;

// TSql160.g:17343
// TSql160.g:17349: action may refuse the input: {VerifyAllowedIndexOption160(statement, vAlterTableAlterIndexOption); if(vAlterTableAlterIndexOption.OptionKind == IndexOptionKind.BucketCount) { optionsList.Add(vAlterTableAlterIndexOption); } else { ThrowIncorrectSyntaxErrorException(vAlterTableAlte...
alterTableAlterIndexOptionList
    : LeftParenthesis indexOption RightParenthesis
    ;

// TSql160.g:17417
lowPriorityLockWaitOption
    : 'WAIT_AT_LOW_PRIORITY' LeftParenthesis 'MAX_DURATION' EqualsSign integer ( 'MINUTES' )? Comma 'ABORT_AFTER_WAIT' EqualsSign Identifier RightParenthesis
    ;

// TSql160.g:17496
suppressMessagesOption
    : 'SUPPRESS_MESSAGES' EqualsSign ( On | Off )
    ;

// TSql160.g:17514
truncateTargetOption
    : 'TRUNCATE_TARGET' ( ( EqualsSign Off ) | ( EqualsSign On ) )
    ;

// TSql160.g:17537
indexStateOption
    : Identifier EqualsSign optionOnOff
    ;

// TSql160.g:17564
orderIndexOption
    : Order identifierColumnList
    ;

// TSql160.g:17601
// TSql160.g:17605: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MaxDop)}?
// TSql160.g:17608: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.DataCompression)}?
// TSql160.g:17611: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.XmlCompression)}?
// TSql160.g:17614: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.BucketCount)}?
// TSql160.g:17617: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Online)}?
// TSql160.g:17620: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.Order)}?
// TSql160.g:17623: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.CompressionDelay)}?
// TSql160.g:17626: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.IgnoreDupKey)}?
// TSql160.g:17629: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MaxDuration)}?
// TSql160.g:17632: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.WaitAtLowPriority)}?
indexOption
    : fillFactorOption
    | maxDopOption_MAXDOP
    | dataCompressionOption_DATA_COMPRESSION
    | xmlCompressionOption_XML_COMPRESSION
    | bucketCountOption_BUCKET_COUNT
    | onlineIndexOption_ONLINE
    | orderIndexOption
    | compressionDelayOption_COMPRESSION_DELAY
    | ignoreDupKeyIndexOption_IGNORE_DUP_KEY
    | maxDurationOption_MAX_DURATION
    | waitAtLowPriorityOption_WAIT_AT_LOW_PRIORITY
    | indexStateOption
    ;

// TSql160.g:17638
// TSql160.g:17649: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.XmlNamespaces)}?
// TSql160.g:17656: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.ChangeTrackingContext)}?
withCommonTableExpressionsAndXmlNamespaces
    : With ( xmlNamespaces_XMLNAMESPACES | ctContext_CHANGE_TRACKING_CONTEXT | commonTableExpression ) ( Comma commonTableExpression )*
    ;

// TSql160.g:17685
commonTableExpression
    : identifier ( columnNameList )? As LeftParenthesis subqueryExpression RightParenthesis
    ;

// TSql160.g:17703
xmlNamespaces
    : 'XMLNAMESPACES' LeftParenthesis xmlNamespacesElement ( Comma xmlNamespacesElement )* RightParenthesis
    ;

// TSql160.g:17727
xmlNamespacesElement
    : xmlNamespacesDefaultElement
    | xmlNamespacesAliasElement
    ;

// TSql160.g:17734
xmlNamespacesDefaultElement
    : Default stringLiteral
    ;

// TSql160.g:17745
xmlNamespacesAliasElement
    : stringLiteral As identifier
    ;

// TSql160.g:17762
subqueryExpressionWithOptionalCTE
    : ( withCommonTableExpressionsAndXmlNamespaces )? subqueryExpression
    ;

// TSql160.g:17780
selectStatement
    : ( withCommonTableExpressionsAndXmlNamespaces )? select
    ;

// TSql160.g:17792
select
    : queryExpression ( orderByClause ( offsetClause )? )? ( forClause )? ( optimizerHints )?
    ;

// TSql160.g:17825
// TSql160.g:17833: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.For) && (LT(2).getText().Equals(CodeGenerationSupporter.Path, StringComparison.OrdinalIgnoreCase))}?
derivedTable
    : ( queryDerivedTable | inlinedDerivedTable ) ( forPathClause )? simpleTableReferenceAlias ( columnNameList )?
    ;

// TSql160.g:17843
queryDerivedTable
    : LeftParenthesis subqueryExpression RightParenthesis
    ;

// TSql160.g:17855
inlinedDerivedTable
    : LeftParenthesis Values rowValueExpression ( Comma rowValueExpression )* RightParenthesis
    ;

// TSql160.g:17875
// TSql160.g:17880: action may refuse the input: {if (ExpressionFlags.ScalarSubqueriesDisallowed == (expressionFlags & ExpressionFlags.ScalarSubqueriesDisallowed)) { ThrowParseErrorException("SQL46098", vQueryExpression, TSqlParserResource.SQL46098Message); } vResult.QueryExpression = vQueryExpressi...
subquery
    : LeftParenthesis subqueryExpression RightParenthesis
    ;

// TSql160.g:17891
subqueryExpression
    : subqueryExpressionUnit ( ( Union | Except | Intersect ) ( All )? subqueryExpressionUnit )*
    ;

// TSql160.g:17932
subqueryExpressionUnit
    : subquerySpecification
    | subqueryParenthesis
    ;

// TSql160.g:17939
subqueryParenthesis
    : LeftParenthesis subqueryExpression RightParenthesis
    ;

// TSql160.g:17951
queryExpression
    : queryExpressionUnit ( ( Union | Except | Intersect ) ( All )? queryExpressionUnit )*
    ;

// TSql160.g:17991
queryExpressionUnit
    : querySpecification
    | queryParenthesis
    ;

// TSql160.g:17998
queryParenthesis
    : LeftParenthesis queryExpression RightParenthesis
    ;

// TSql160.g:18010
// TSql160.g:18087: predicate dropped: {LA(1) == For && LA(2) == Browse}?
// TSql160.g:18099: action may refuse the input: {if(vResult.ForClause != null) { ThrowIncorrectSyntaxErrorException(tFor); } vResult.ForClause = vXmlForClause;}
// TSql160.g:18108: action may refuse the input: {if(vResult.ForClause != null) { ThrowIncorrectSyntaxErrorException(tFor); } vResult.ForClause = vJsonForClause;}
// TSql160.g:18117: action may refuse the input: {if (vResult.OrderByClause != null && vResult.TopRowFilter == null && vResult.ForClause == null && vResult.OffsetClause == null) { ThrowParseErrorException("SQL46047", vResult, TSqlParserResource.SQL46047Message); } if (vResult.TopRowFilter != null &&...
// TSql160.g:18097: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Xml)}?
// TSql160.g:18106: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Json)}?
subquerySpecification
    : Select uniqueRowFilterOpt ( topRowFilter )? selectColumnOrStarExpression ( Comma selectColumnOrStarExpression )* fromClauseOpt ( whereClause )? ( groupByClause )? ( havingClause )? ( windowClause )? ( orderByClause ( offsetClause )? )? ( For browseForClause |  ) ( For ( xmlForClause_XML | jsonForClause_JSON ) )?
    ;

// TSql160.g:18133
browseForClause
    : Browse
    ;

// TSql160.g:18141
// TSql160.g:18171: action may refuse the input: {if(vSelectStatement == null) { ThrowIncorrectSyntaxErrorException(tInto); } vSelectStatement.Into = vSchemaObjectName;}
// TSql160.g:18182: action may refuse the input: {if(vSelectStatement == null) { ThrowIncorrectSyntaxErrorException(tOn); } if(vSelectStatement.Into == null) { ThrowIncorrectSyntaxErrorException(tOn); } vSelectStatement.On = vFileGroupName;}
querySpecification
    : Select uniqueRowFilterOpt ( topRowFilter )? selectExpression ( Comma selectExpression )* ( Into schemaObjectThreePartName )? ( On identifier )? fromClauseOpt ( whereClause )? ( groupByClause )? ( havingClause )? ( windowClause )?
    ;

// TSql160.g:18224
uniqueRowFilterOpt
    : ( uniqueRowFilter )?
    ;

// TSql160.g:18236
uniqueRowFilter
    : All
    | Distinct
    ;

// TSql160.g:18254
offsetClause
    : 'OFFSET' expression ( 'ROW' | 'ROWS' ) ( Fetch ( 'FIRST' | 'NEXT' ) expression ( 'ROW' | 'ROWS' ) 'ONLY' )?
    ;

// TSql160.g:18289
topRowFilter
    : Top ( integerOrRealOrNumeric | parenthesisDisambiguatorForExpressions ) ( Percent )? ( With 'TIES' )?
    ;

// TSql160.g:18325
dmlTopRowFilter
    : Top parenthesisDisambiguatorForExpressions ( Percent )?
    ;

// TSql160.g:18343
dmlTopRowFilterOpt
    : ( dmlTopRowFilter )?
    ;

// TSql160.g:18355
selectExpression
    : selectSetVariable
    | selectStarExpression
    | selectColumn
    ;

// TSql160.g:18380
selectColumnOrStarExpression
    : selectStarExpression
    | selectColumn
    ;

// TSql160.g:18388
outputClauseSelectElement
    : selectStarExpression
    | outputClauseSelectColumn
    ;

// TSql160.g:18396
outputClauseSelectColumn
    : expression ( ( As )? ( stringOrIdentifier ) )?
    ;

// TSql160.g:18419
selectColumn
    : selectColumnExpression ( ( As )? ( stringOrIdentifier ) )?
    | ( stringOrIdentifier ) EqualsSign ( selectColumnExpression )
    ;

// TSql160.g:18456
selectColumnExpression
    : expression
    | identityFunction
    ;

// TSql160.g:18463
identityFunction
    : Identity LeftParenthesis scalarDataType ( Comma seedIncrement Comma seedIncrement )? RightParenthesis
    ;

// TSql160.g:18494
selectSetVariable
    : variable assignmentWithOptOp expression
    ;

// TSql160.g:18508
// TSql160.g:18546: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Repeatable)}?
tableSampleClause
    : TableSample ( 'SYSTEM' )? LeftParenthesis expression ( 'ROWS' | Percent )? RightParenthesis ( 'REPEATABLE' LeftParenthesis expression RightParenthesis |  )
    ;

// TSql160.g:18561
fromClauseOpt
    : ( fromClause |  )
    ;

// TSql160.g:18569
fromClause
    : From selectTableReferenceWithOdbc ( Comma selectTableReferenceWithOdbc )*
    ;

// TSql160.g:18586
forPathClause
    : ( For 'PATH' )
    ;

// TSql160.g:18614
// TSql160.g:18639: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Contained)}?
temporalClause
    : ( For 'SYSTEM_TIME' ( As Of stringLiteralOrVariable | From stringLiteralOrVariable To stringLiteralOrVariable | Between stringLiteralOrVariable And stringLiteralOrVariable | 'CONTAINED' In LeftParenthesis stringLiteralOrVariable Comma stringLiteralOrVariable RightParenthesis | All ) )
    ;

// TSql160.g:18669
selectTableReferenceWithOdbc
    : selectTableReference
    | odbcQualifiedJoin
    ;

// TSql160.g:18676
selectTableReference
    : selectTableReferenceElement ( selectTableReferenceAdditionalElement )*
    ;

// TSql160.g:18691
// TSql160.g:18712: action may refuse the input: {if (!(vSelectTableReference is QualifiedJoin)) { ThrowParseErrorException("SQL46035", tLCurly, TSqlParserResource.SQL46035Message); }}
// TSql160.g:18685: alternative dropped, its action refuses every input: tOdbcInitiator : OdbcInitiator {ThrowParseErrorException("SQL46036", tOdbcInitiator, TSqlParserResource.SQL46036Message);}
// TSql160.g:18684: rule dropped, every alternative refuses its input
odbcQualifiedJoin
    : ( LeftCurly ) 'OJ' ( odbcQualifiedJoin | selectTableReference ) RightCurly
    ;

// TSql160.g:18728
odbcConvertSpecification
    : nonQuotedIdentifier
    ;

// TSql160.g:18754
// TSql160.g:18835: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Extract)}?
odbcFunctionCall
    : LeftCurly 'FN' ( Convert LeftParenthesis expression Comma odbcConvertSpecification RightParenthesis | Truncate LeftParenthesis expression Comma expression RightParenthesis | ( Database | User | CurrentDate ) LeftParenthesis RightParenthesis | ( Insert | Left | Right ) LeftParenthesis expressionList RightParenthesis | ( CurrentTime | CurrentTimestamp ) ( LeftParenthesis ( expressionList )? RightParenthesis |  ) | extractFromExpression_EXTRACT | Identifier LeftParenthesis ( expressionList )? RightParenthesis ) RightCurly
    ;

// TSql160.g:18853
joinTableReference
    : ( joinParenthesis ) ( joinElement )*
    | selectTableReferenceElementWithoutJoinParenthesis ( joinElement )+
    ;

// TSql160.g:18877
selectTableReferenceAdditionalElement
    : joinElement
    | pivotedTableReference
    | unpivotedTableReference
    ;

// TSql160.g:18886
pivotedTableReference
    : Pivot LeftParenthesis multiPartIdentifier ( LeftParenthesis fixedColumn ( Comma fixedColumn )* RightParenthesis ) For fixedColumn In columnNameList RightParenthesis simpleTableReferenceAlias
    ;

// TSql160.g:18920
unpivotedTableReference
    : Unpivot LeftParenthesis identifier For identifier In columnListWithParenthesis RightParenthesis simpleTableReferenceAlias
    ;

// TSql160.g:18941
fixedColumn
    : multiPartIdentifier
    ;

// TSql160.g:18953
columnListWithParenthesis
    : LeftParenthesis fixedColumn ( Comma fixedColumn )* RightParenthesis
    ;

// TSql160.g:18972
identifierColumnList
    : LeftParenthesis ( identifierColumnReferenceExpression | graphPseudoColumn ) ( Comma ( identifierColumnReferenceExpression | graphPseudoColumn ) )* RightParenthesis
    ;

// TSql160.g:19005
joinElement
    : unqualifiedJoin
    | qualifiedJoin
    ;

// TSql160.g:19011
selectTableReferenceElement
    : joinParenthesis
    | selectTableReferenceElementWithoutJoinParenthesis
    ;

// TSql160.g:19018
// TSql160.g:19023: predicate dropped: {NextIdentifierMatchesOneOf(new string[] {CodeGenerationSupporter.StringSplit, CodeGenerationSupporter.GenerateSeries})}?
// TSql160.g:19020: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.ChangeTable)}?
// TSql160.g:19032: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.OpenJson)}?
// TSql160.g:19035: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Predict)}?
selectTableReferenceElementWithoutJoinParenthesis
    : changeTableTableReference_CHANGETABLE
    | builtInFunctionTableReference
    | globalFunctionTableReference
    | variableTableReference
    | variableMethodCallTableReference
    | derivedTable
    | openRowset
    | fulltextTableReference
    | semanticTableReference
    | openXmlTableReference
    | openJsonTableReference_OPENJSON
    | subDmlTableReference
    | predictTableReference_PREDICT
    | schemaObjectOrFunctionTableReference
    ;

// TSql160.g:19051
predictParams
    : ( 'MODEL' EqualsSign expression | Identifier EqualsSign subquery ) Comma Identifier EqualsSign mergeTarget ( Comma Identifier EqualsSign identifier )?
    ;

// TSql160.g:19083
predictWithClauseOpt
    : ( ( With LeftParenthesis predictSchemaItemList RightParenthesis ) )
    ;

// TSql160.g:19094
predictSchemaItemList
    : predictSchemaItem ( Comma predictSchemaItem )*
    ;

// TSql160.g:19109
predictSchemaItem
    : columnDefinitionBasic ( stringLiteral )? ( As 'PREDICT' )?
    ;

// TSql160.g:19130
// TSql160.g:19144: predicate dropped: {!NextTokenMatches(CodeGenerationSupporter.Using)}?
mergeTarget
    : dmlTarget ( ( As identifier ) | ( identifier ) |  )
    ;

// TSql160.g:19187
changesChangeTableParams
    : ( integerOrVariable | nullLiteral ) ( Comma 'FORCESEEK' )?
    ;

// TSql160.g:19210
versionChangeTableParams
    : columnNameList Comma LeftParenthesis expressionList RightParenthesis ( Comma 'FORCESEEK' )?
    ;

// TSql160.g:19229
// TSql160.g:19234: action may refuse the input: {UpdateTokenInfo(vResult, tLParen); if ((subDmlFlags & SubDmlFlags.InsideSubDml) == SubDmlFlags.InsideSubDml) ThrowParseErrorException("SQL46075", tLParen, TSqlParserResource.SQL46075Message);}
// TSql160.g:19240: action may refuse the input: {if ((subDmlFlags & SubDmlFlags.SelectNotForInsert) == SubDmlFlags.SelectNotForInsert) ThrowParseErrorException("SQL46076", vSpec, TSqlParserResource.SQL46076Message); if ((subDmlFlags & SubDmlFlags.UpdateDeleteFrom) == SubDmlFlags.UpdateDeleteFrom) T...
subDmlTableReference
    : LeftParenthesis innerDmlStatement RightParenthesis simpleTableReferenceAlias ( columnNameList )?
    ;

// TSql160.g:19259
innerDmlStatement
    : insertSpecification
    | updateSpecification
    | deleteSpecification
    | mergeSpecification
    ;

// TSql160.g:19266
unqualifiedJoin
    : ( Cross ( Join | 'APPLY' ) | Outer 'APPLY' ) selectTableReferenceElement
    ;

// TSql160.g:19301
qualifiedJoin
    : ( Join | ( Inner | Left ( Outer )? | Right ( Outer )? | Full ( Outer )? ) ( joinHint )? Join ) selectTableReferenceWithOdbc On booleanExpression
    ;

// TSql160.g:19349
// TSql160.g:19354: action may refuse the input: {Match(tId1,CodeGenerationSupporter.Local); vParent.JoinHint = JoinHintHelper.Instance.ParseOption(tId2); if (vParent.JoinHint == JoinHint.Remote) { ThrowIncorrectSyntaxErrorException(tId2); }}
joinHint
    : ( Identifier ( Identifier | Merge |  ) )
    | Merge
    ;

// TSql160.g:19382
joinParenthesis
    : LeftParenthesis joinTableReference RightParenthesis
    ;

// TSql160.g:19394
builtInFunctionTableReference
    : DoubleColon identifier LeftParenthesis ( expressionWithDefault ( Comma expressionWithDefault )* )? RightParenthesis simpleTableReferenceAliasOpt
    ;

// TSql160.g:19424
globalFunctionTableReference
    : identifier LeftParenthesis ( expressionWithDefault ( Comma expressionWithDefault )* )? RightParenthesis simpleTableReferenceAliasOpt
    ;

// TSql160.g:19454
variableTableReference
    : variable simpleTableReferenceAliasOpt
    ;

// TSql160.g:19465
variableMethodCallTableReference
    : variable Dot identifier parenthesizedOptExpressionWithDefaultList ( simpleTableReferenceAlias columnNameList )
    ;

// TSql160.g:19486
raiseErrorStatement
    : Raiserror LeftParenthesis signedIntegerOrStringOrVariable Comma signedIntegerOrVariable Comma signedIntegerOrVariable ( Comma possibleNegativeConstant )* RightParenthesis ( With Identifier ( Comma Identifier )* )?
    ;

// TSql160.g:19609
outputIntoClause
    : Into ( variableDmlTarget | intoSchemaObjectTable ) ( identifierColumnList )?
    ;

// TSql160.g:19653
intoSchemaObjectTable
    : schemaObjectFourPartName
    ;

// TSql160.g:19663
// TSql160.g:19669: action may refuse the input: {if ((subDmlFlags & SubDmlFlags.InsideSubDml) == SubDmlFlags.InsideSubDml) ThrowParseErrorException("SQL46079", vParent, TSqlParserResource.SQL46079Message);}
// TSql160.g:19665: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Output)}?
outputClauseOpt
    : ( outputIntoOutputClause_OUTPUT |  )
    ;

// TSql160.g:19676
deleteStatement
    : deleteSpecification ( optimizerHints )?
    ;

// TSql160.g:19690
deleteSpecification
    : Delete dmlTopRowFilterOpt ( From )? dmlTarget outputClauseOpt fromClauseOpt ( dmlWhereClause )?
    ;

// TSql160.g:19721
insertStatement
    : insertSpecification ( optimizerHints )?
    ;

// TSql160.g:19735
insertSpecification
    : Insert dmlTopRowFilterOpt ( Into | Over )? dmlTarget ( LeftParenthesis insertColumn ( Comma insertColumn )* RightParenthesis )? outputClauseOpt ( ( valuesInsertSource ) | executeInsertSource | selectInsertSource )
    ;

// TSql160.g:19794
updateStatement
    : updateSpecification ( optimizerHints )?
    ;

// TSql160.g:19808
updateSpecification
    : Update dmlTopRowFilterOpt dmlTarget setClausesList outputClauseOpt fromClauseOpt ( dmlWhereClause )?
    ;

// TSql160.g:19837
setClausesList
    : Set setClause ( Comma setClause )*
    ;

// TSql160.g:19854
setClause
    : setClauseStartingWithVariable
    | multiPartIdentifier ( setClauseColumnAssignment | setClauseFunctionCall )
    | setClauseSubItemSystemColumn
    ;

// TSql160.g:19871
setClauseStartingWithVariable
    : variable ( ( EqualsSign ( multiPartIdentifier assignmentWithOptOp expressionWithDefault | expression ) ) | ( assignmentWithOp expression ) )
    ;

// TSql160.g:19912
setClauseColumnAssignment
    : assignmentWithOptOp expressionWithDefault
    ;

// TSql160.g:19925
setClauseFunctionCall
    : parenthesizedOptExpressionWithDefaultList
    ;

// TSql160.g:19937
setClauseSubItemSystemColumn
    : systemColumn assignmentWithOptOp expressionWithDefault
    ;

// TSql160.g:19951
executeInsertSource
    : executeSpecification
    ;

// TSql160.g:19962
selectInsertSource
    : queryExpression ( orderByClause ( offsetClause )? )?
    ;

// TSql160.g:19987
valuesInsertSource
    : defaultValuesInsertSource
    | Values rowValueExpressionWithDefault ( Comma rowValueExpressionWithDefault )*
    ;

// TSql160.g:20006
defaultValuesInsertSource
    : Default Values
    ;

// TSql160.g:20016
rowValueExpression
    : LeftParenthesis expressionList RightParenthesis
    ;

// TSql160.g:20029
rowValueExpressionWithDefault
    : LeftParenthesis expressionWithDefaultList RightParenthesis
    ;

// TSql160.g:20042
insertColumn
    : multiPartIdentifier
    | systemColumn
    ;

// TSql160.g:20055
openRowsetColumn
    : columnDefinitionBasic ( integer | stringLiteral )?
    ;

// TSql160.g:20077
dmlTarget
    : schemaObjectDmlTarget
    | openRowset
    | variableDmlTarget
    ;

// TSql160.g:20084
mergeInsertDmlColumnListOpt
    : ( LeftParenthesis mergeInsertDmlColumn ( Comma mergeInsertDmlColumn )* RightParenthesis )?
    ;

// TSql160.g:20106
mergeInsertDmlColumn
    : multiPartIdentifier
    | systemColumn
    ;

// TSql160.g:20119
variableDmlTarget
    : variable
    ;

// TSql160.g:20130
schemaObjectDmlTarget
    : schemaObjectFunctionDmlTarget
    | schemaObjectTableDmlTarget
    ;

// TSql160.g:20137
schemaObjectFunctionDmlTarget
    : schemaObjectFourPartName parenthesizedOptExpressionWithDefaultList
    ;

// TSql160.g:20150
schemaObjectTableDmlTarget
    : schemaObjectFourPartName ( With tableHints )?
    ;

// TSql160.g:20164
// TSql160.g:20171: predicate dropped: {IsTableReference(false)}?
schemaObjectOrFunctionTableReference
    : schemaObjectFourPartName ( schemaObjectTableReference | schemaObjectFunctionTableReference )
    ;

// TSql160.g:20177
schemaObjectFunctionTableReference
    : ( parenthesizedOptExpressionWithDefaultList simpleTableReferenceAliasOpt ( columnNameList )? )
    ;

// TSql160.g:20193
// TSql160.g:20235: action may refuse the input: {if (vResult.TableSampleClause != null) { // If tableSampleClause was parsed, this is an error. throw GetUnexpectedTokenErrorException(); }}
// TSql160.g:20203: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.For) && (LT(2).getText().Equals(CodeGenerationSupporter.SystemTime, StringComparison.OrdinalIgnoreCase))}?
// TSql160.g:20210: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.For) && (LT(2).getText().Equals(CodeGenerationSupporter.Path, StringComparison.OrdinalIgnoreCase))}?
schemaObjectTableReference
    : ( temporalClause )? ( forPathClause )? ( ( singleOldStyleTableHint ( simpleTableReferenceAlias )? ) | ( tableSampleClause )? ( nonParameterTableHints | simpleTableReferenceAlias ( tableSampleClause )? ( oldForceIndex | nonParameterTableHints | singleOldStyleTableHint )? )? )
    ;

// TSql160.g:20264
simpleTableReferenceAlias
    : ( As )? identifier
    ;

// TSql160.g:20284
simpleTableReferenceAliasOpt
    : ( simpleTableReferenceAlias )?
    ;

// TSql160.g:20292
fulltextTableReference
    : fullTextTable LeftParenthesis schemaObjectFourPartName Comma fulltextTableColumnList Comma stringOrVariable ( fulltextTableOptions )? RightParenthesis simpleTableReferenceAliasOpt
    ;

// TSql160.g:20314
fullTextTable
    : ContainsTable
    | FreeTextTable
    ;

// TSql160.g:20327
fulltextTableColumnList
    : fulltextSemanticTableColumnList
    | 'PROPERTY' LeftParenthesis identifierColumnReferenceExpression Comma stringLiteral RightParenthesis
    ;

// TSql160.g:20346
fulltextSemanticTableColumnList
    : identifierColumnReferenceExpression
    | starColumnReferenceExpression
    | LeftParenthesis starColumnReferenceExpression RightParenthesis
    | LeftParenthesis identifierColumnReferenceExpression ( Comma identifierColumnReferenceExpression )* RightParenthesis
    ;

// TSql160.g:20373
fulltextTableOptions
    : Comma languageExpression ( Comma integerOrVariable )?
    | Comma integerOrVariable ( Comma languageExpression )?
    ;

// TSql160.g:20398
languageExpression
    : 'LANGUAGE' binaryOrIntegerOrStringOrVariable
    ;

// TSql160.g:20407
semanticTableReference
    : semanticKeyPhraseTableReference
    | semanticSimilarityTableReference
    | semanticSimilarityDetailsTableReference
    ;

// TSql160.g:20414
semanticKeyPhraseTableReference
    : SemanticKeyPhraseTable LeftParenthesis schemaObjectFourPartName Comma fulltextSemanticTableColumnList ( Comma possibleNegativeConstant )? RightParenthesis simpleTableReferenceAliasOpt
    ;

// TSql160.g:20439
semanticSimilarityTableReference
    : SemanticSimilarityTable LeftParenthesis schemaObjectFourPartName Comma fulltextSemanticTableColumnList Comma possibleNegativeConstant RightParenthesis simpleTableReferenceAliasOpt
    ;

// TSql160.g:20462
semanticSimilarityDetailsTableReference
    : SemanticSimilarityDetailsTable LeftParenthesis schemaObjectFourPartName Comma identifierColumnReferenceExpression Comma possibleNegativeConstant Comma identifierColumnReferenceExpression Comma possibleNegativeConstant RightParenthesis simpleTableReferenceAliasOpt
    ;

// TSql160.g:20499
openXmlTableReference
    : OpenXml LeftParenthesis openXmlParams RightParenthesis openXmlWithClauseOpt simpleTableReferenceAliasOpt
    ;

// TSql160.g:20510
openXmlWithClauseOpt
    : ( ( With LeftParenthesis openXmlSchemaItemList RightParenthesis ) | ( With schemaObjectThreePartName ) )
    | /* empty */
    ;

// TSql160.g:20531
openXmlParams
    : variable Comma stringOrVariable ( Comma integerOrVariable )?
    ;

// TSql160.g:20549
openXmlSchemaItemList
    : openXmlSchemaItem ( Comma openXmlSchemaItem )*
    ;

// TSql160.g:20564
openXmlSchemaItem
    : columnDefinitionBasic ( stringOrVariable )?
    ;

// TSql160.g:20590
openJsonWithClauseOpt
    : ( ( With LeftParenthesis openJsonSchemaItemList RightParenthesis ) )
    | /* empty */
    ;

// TSql160.g:20602
openJsonParams
    : expression ( Comma expression )?
    ;

// TSql160.g:20618
openJsonSchemaItemList
    : openJsonSchemaItem ( Comma openJsonSchemaItem )*
    ;

// TSql160.g:20633
openJsonSchemaItem
    : columnDefinitionBasic ( stringLiteral )? ( As 'JSON' )?
    ;

// TSql160.g:20656
openRowset
    : openRowsetRowset
    | openQueryRowset simpleTableReferenceAliasOpt
    | adhocRowset simpleTableReferenceAliasOpt
    ;

// TSql160.g:20666
openRowsetRowset
    : OpenRowSet LeftParenthesis ( openRowsetParams | internalOpenRowsetArgs | openRowsetBulk | openRowsetCosmos )
    ;

// TSql160.g:20681
internalOpenRowsetArgs
    : identifier ( Comma possibleNegativeConstant )* RightParenthesis simpleTableReferenceAliasOpt
    ;

// TSql160.g:20702
openRowsetCosmos
    : openRowsetCosmosOptionHint ( Comma openRowsetCosmosOptionHint )* RightParenthesis ( With ( LeftParenthesis openRowsetColumn ( Comma openRowsetColumn )* RightParenthesis ) )? simpleTableReferenceAliasOpt
    ;

// TSql160.g:20749
openRowsetCosmosOptionHint
    : Identifier EqualsSign nonEmptyString
    ;

// TSql160.g:20761
openRowsetBulk
    : Bulk ( ( stringLiteral ) | ( LeftParenthesis stringLiteral ( Comma stringLiteral )* ( Comma )? RightParenthesis ) ) ( Comma openRowsetBulkHint )* RightParenthesis ( With ( LeftParenthesis openRowsetColumn ( Comma openRowsetColumn )* RightParenthesis ) )? simpleTableReferenceAliasOpt ( columnNameList )?
    ;

// TSql160.g:20833
openRowsetBulkHint
    : openRowsetBulkHintNoValue
    | simpleBulkInsertOptionWithValue
    | openRowsetBulkOrderHint
    ;

// TSql160.g:20839
openRowsetBulkHintNoValue
    : Identifier
    ;

// TSql160.g:20847
openQueryRowset
    : OpenQuery LeftParenthesis identifier Comma stringLiteral RightParenthesis
    ;

// TSql160.g:20861
adhocRowset
    : adhocDataSource Dot objectOrString
    ;

// TSql160.g:20893
openRowsetParams
    : stringLiteral Comma ( ( stringLiteral Semicolon ( stringLiteral )? Semicolon ( stringLiteral )? ) | ( stringLiteral ) ) Comma ( schemaObjectThreePartName | stringLiteral ) RightParenthesis ( With ( LeftParenthesis openRowsetColumn ( Comma openRowsetColumn )* RightParenthesis ) )? simpleTableReferenceAliasOpt
    ;

// TSql160.g:20968
oldForceIndex
    : LeftParenthesis integer RightParenthesis
    ;

// TSql160.g:20981
nonParameterTableHints
    : HoldLock ( oldForceIndex | simpleTableHints )?
    | simpleTableHints
    ;

// TSql160.g:21009
simpleTableHints
    : With ( oldForceIndex | tableHints )
    ;

// TSql160.g:21027
tableHints
    : LeftParenthesis tableHintsBody RightParenthesis
    ;

// TSql160.g:21040
tableHintsBody
    : tableHint ( ( Comma )? tableHint )*
    ;

// TSql160.g:21056
// TSql160.g:21058: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.ForceSeek)}?
tableHint
    : forceSeekTableHint_FORCESEEK
    | simpleTableHint
    | indexTableHint
    | literalTableHint
    ;

// TSql160.g:21068
simpleTableHint
    : HoldLock
    | Identifier
    ;

// TSql160.g:21082
// TSql160.g:21087: action may refuse the input: {if (!indexHintAllowed) ThrowParseErrorException("SQL46074", tIndex, TSqlParserResource.SQL46074Message); UpdateTokenInfo(vResult,tIndex); vResult.HintKind=TableHintKind.Index;}
indexTableHint
    : Index ( EqualsSign identifierOrInteger | LeftParenthesis identifierOrInteger ( Comma identifierOrInteger )* RightParenthesis )
    ;

// TSql160.g:21117
literalTableHint
    : 'SPATIAL_WINDOW_MAX_CELLS' EqualsSign integer
    ;

// TSql160.g:21170
singleOldStyleTableHint
    : LeftParenthesis tableHint RightParenthesis
    ;

// TSql160.g:21185
whereClause
    : Where booleanExpressionWithFlags
    ;

// TSql160.g:21196
dmlWhereClause
    : whereClause
    | whereCurrentOfCursorClause
    ;

// TSql160.g:21203
// TSql160.g:21208: action may refuse the input: {if ((subDmlFlags & SubDmlFlags.InsideSubDml) == SubDmlFlags.InsideSubDml) ThrowParseErrorException("SQL46083", tWhere, TSqlParserResource.SQL46083Message); UpdateTokenInfo(vResult,tWhere); vResult.Cursor = vCursorId;}
whereCurrentOfCursorClause
    : Where Current Of cursorId
    ;

// TSql160.g:21220
// TSql160.g:21249: action may refuse the input: {if (vResult.All) ThrowParseErrorException("SQL46084", tOption, TSqlParserResource.SQL46084Message); if (encounteredCubeRollupGroupingSets) ThrowParseErrorException("SQL46085", tOption, TSqlParserResource.SQL46085Message); UpdateTokenInfo(vResult,tOpt...
groupByClause
    : Group By ( All )? groupByItem ( Comma groupByItem )* ( With Identifier )?
    ;

// TSql160.g:21262
// TSql160.g:21264: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Cube)}?
// TSql160.g:21269: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Rollup)}?
// TSql160.g:21274: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Grouping)}?
groupByItem
    : cubeSpec_CUBE
    | rollupSpec_ROLLUP
    | groupingSetsSpec_GROUPING
    | grandTotal
    | simpleGroupByItem
    ;

// TSql160.g:21309
groupingSet
    : grandTotal
    | composingGroupingSet
    | groupingSetItem
    ;

// TSql160.g:21320
composingGroupingSet
    : LeftParenthesis groupingSetItem ( Comma groupingSetItem )* RightParenthesis
    ;

// TSql160.g:21336
grandTotal
    : LeftParenthesis RightParenthesis
    ;

// TSql160.g:21345
// TSql160.g:21350: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Cube)}?
// TSql160.g:21352: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Rollup)}?
groupingSetItem
    : cubeSpec_CUBE
    | rollupSpec_ROLLUP
    | simpleGroupByItem
    ;

// TSql160.g:21386
cubeOrRollupArgumentList
    : cubeOrRollupArgument ( Comma cubeOrRollupArgument )*
    ;

// TSql160.g:21401
cubeOrRollupArgument
    : compositeGroupByItem
    | simpleGroupByItem
    ;

// TSql160.g:21412
compositeGroupByItem
    : LeftParenthesis simpleGroupByItem ( Comma simpleGroupByItem )* RightParenthesis
    ;

// TSql160.g:21433
// TSql160.g:21445: action may refuse the input: {Match(tDistributedAgg, CodeGenerationSupporter.DistributedAgg); if (alreadyEncounteredDistributedAggHint) ThrowParseErrorException("SQL46129", tDistributedAgg, TSqlParserResource.SQL46129Message); vResult.DistributedAggregation = true; UpdateTokenInf...
simpleGroupByItem
    : expression ( With LeftParenthesis 'DISTRIBUTED_AGG' RightParenthesis )?
    ;

// TSql160.g:21463
havingClause
    : Having booleanExpression
    ;

// TSql160.g:21478
orderByClause
    : Order By expressionWithSortOrder ( Comma expressionWithSortOrder )*
    ;

// TSql160.g:21499
expressionWithSortOrder
    : expression ( orderByOption )?
    ;

// TSql160.g:21516
// TSql160.g:21529: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Xml)}?
// TSql160.g:21534: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Json)}?
forClause
    : For ( browseForClause | ( Read 'ONLY' ) | xmlForClause_XML | jsonForClause_JSON | updateForClause )
    ;

// TSql160.g:21564
// TSql160.g:21575: action may refuse the input: {if (vResult.OptionKind == XmlForClauseOptions.Explicit || vResult.OptionKind == XmlForClauseOptions.Auto) throw GetUnexpectedTokenErrorException(tLParen); vResult.Value = vElementName; UpdateTokenInfo(vResult,tRParen);}
xmlForClauseMode
    : Identifier ( LeftParenthesis stringLiteral RightParenthesis )?
    ;

// TSql160.g:21586
// TSql160.g:21593: action may refuse the input: {if (TryMatch(tOption,CodeGenerationSupporter.XmlSchema) || TryMatch(tOption,CodeGenerationSupporter.Root)) { vResult.Value = vValue; UpdateTokenInfo(vResult,tRParen); } else throw GetUnexpectedTokenErrorException(tLParen);}
xmlParam
    : Identifier ( LeftParenthesis stringLiteral RightParenthesis | ( 'BASE64' | 'XSINIL' | 'ABSENT' ) )?
    ;

// TSql160.g:21654
jsonForClauseMode
    : Identifier
    ;

// TSql160.g:21662
// TSql160.g:21669: action may refuse the input: {if (TryMatch(tOption,CodeGenerationSupporter.Root)) { vResult.Value = vValue; UpdateTokenInfo(vResult,tRParen); } else throw GetUnexpectedTokenErrorException(tLParen);}
// TSql160.g:21679: action may refuse the input: {vResult.OptionKind = JsonForClauseOptionsHelper.Instance.ParseOption(tOption); UpdateTokenInfo(vResult,tOption); if (((encountered & vResult.OptionKind) != 0) || (((encountered & JsonForClauseOptions.WithoutArrayWrapper) != 0) && (vResult.OptionKind ...
jsonParam
    : Identifier ( LeftParenthesis stringLiteral RightParenthesis )?
    ;

// TSql160.g:21689
updateForClause
    : Update ( Of column ( Comma column )* )?
    ;

// TSql160.g:21709
optimizerHints
    : Option LeftParenthesis hint ( Comma hint )* RightParenthesis
    ;

// TSql160.g:21728
hint
    : literalOptimizerHint
    | simpleOptimizerHint
    | usePlanOptimizerHint
    | optimizeForOptimizerHint
    | tableHintsOptimizerHint
    | useHintClause
    ;

// TSql160.g:21737
tableHintsOptimizerHint
    : Table 'HINT' LeftParenthesis schemaObjectFourPartName ( Comma tableHintsBody )? RightParenthesis
    ;

// TSql160.g:21755
simpleOptimizerHint
    : Identifier Join
    | Merge Join
    | Identifier Union
    | Merge Union
    | 'FORCE' Order
    | 'HASH' Group
    | Order Group
    | Identifier Plan
    | Identifier Identifier (  | Union All )
    | Identifier
    ;

// TSql160.g:21840
literalOptimizerHint
    : Identifier integer
    | Identifier EqualsSign ( integerOrRealOrNumeric | stringLiteral )
    ;

// TSql160.g:21865
usePlanOptimizerHint
    : Use Plan stringLiteral
    ;

// TSql160.g:21876
optimizeForOptimizerHint
    : 'OPTIMIZE' For ( LeftParenthesis variableValuePair ( Comma variableValuePair )* RightParenthesis | 'UNKNOWN' )
    ;

// TSql160.g:21910
variableValuePair
    : variable ( EqualsSign possibleNegativeConstant | 'UNKNOWN' )
    ;

// TSql160.g:21934
useHintClause
    : Use 'HINT' LeftParenthesis stringLiteral ( Comma stringLiteral )* RightParenthesis
    ;

// TSql160.g:21958
createRuleStatement
    : Create Rule schemaObjectThreePartName As booleanExpression
    ;

// TSql160.g:21998
cryptographicProviderFile
    : From File EqualsSign stringLiteral
    ;

// TSql160.g:22127
resourcePoolStatementBody
    : ( With LeftParenthesis resourcePoolParameter ( Comma resourcePoolParameter )* RightParenthesis )?
    ;

// TSql160.g:22150
// TSql160.g:22168: action may refuse the input: {Int32 value; if (Int32.TryParse(vPercent.Value, NumberStyles.Integer, CultureInfo.InvariantCulture, out value)) { if ((vResult.ParameterType == ResourcePoolParameterType.MinCpuPercent) || (vResult.ParameterType == ResourcePoolParameterType.MinMemoryP...
// TSql160.g:22156: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Affinity)}?
resourcePoolParameter
    : resourcePoolAffinitySpecification_AFFINITY
    | Identifier EqualsSign integer
    ;

// TSql160.g:22242
poolAffinityRange
    : integer ( To integer )?
    ;

// TSql160.g:22327
externalResourcePoolStatementBody
    : ( With LeftParenthesis externalResourcePoolParameter ( Comma externalResourcePoolParameter )* RightParenthesis )?
    ;

// TSql160.g:22356
// TSql160.g:22375: action may refuse the input: {if ((vResult.ParameterType == ExternalResourcePoolParameterType.MaxCpuPercent) || (vResult.ParameterType == ExternalResourcePoolParameterType.MaxMemoryPercent)) { Int32 vPercentIntValue = Int32.Parse(vPercent.Value, NumberStyles.Integer, CultureInfo....
// TSql160.g:22362: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Affinity)}?
externalResourcePoolParameter
    : externalResourcePoolAffinitySpecification_AFFINITY
    | Identifier EqualsSign integer
    ;

// TSql160.g:22430
externalPoolAffinityRange
    : integer ( To integer )?
    ;

// TSql160.g:22515
// TSql160.g:22522: action may refuse the input: {Match(tExternal, CodeGenerationSupporter.External); UpdateTokenInfo(vParent, tExternal); if(vParent.ExternalPoolName != null) { ThrowParseErrorException("SQL46049", tExternal, TSqlParserResource.SQL46049Message, tExternal.getText()); } vParent.Extern...
// TSql160.g:22534: action may refuse the input: {if(vParent.PoolName != null) { ThrowParseErrorException("SQL46049", vIdentifier, TSqlParserResource.SQL46049Message, "pool_name"); } vParent.PoolName = vIdentifier;}
workloadGroupUsingOption
    : 'EXTERNAL' identifier
    | identifier
    ;

// TSql160.g:22544
// TSql160.g:22563: action may refuse the input: {WorkloadGroupResourceParameter workloadGroupResource = this.FragmentFactory.CreateFragment<WorkloadGroupResourceParameter>(); workloadGroupResource.ParameterType = WorkloadGroupResourceParameterHelper.Instance.ParseOption(tParameter, SqlVersionFlags....
workloadGroupParameter
    : Identifier EqualsSign ( Identifier | integerOrRealOrNumeric )
    ;

// TSql160.g:22647
withClassifierOptions
    : With LeftParenthesis workloadClassifierOption ( Comma workloadClassifierOption )* RightParenthesis
    ;

// TSql160.g:22674
// TSql160.g:22677: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.WorkloadGroup)}?
// TSql160.g:22680: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MemberName)}?
// TSql160.g:22683: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.WlmLabel)}?
// TSql160.g:22686: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.WlmContext)}?
// TSql160.g:22689: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Importance)}?
// TSql160.g:22692: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.StartTime)}?
// TSql160.g:22695: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.EndTime)}?
workloadClassifierOption
    : ( classifierWorkloadGroupOption_WORKLOAD_GROUP | classifierMemberNameOption_MEMBERNAME | classifierWlmLabelOption_WLM_LABEL | classifierWlmContextOption_WLM_CONTEXT | classifierImportanceOption_IMPORTANCE | classifierStartTimeOption_START_TIME | classifierEndTimeOption_END_TIME )
    ;

// TSql160.g:22794
wlmTimeLiteral
    : stringLiteral
    ;

// TSql160.g:22876
// TSql160.g:22891: action may refuse the input: {// should not be remote service name or priority_level if ((vResult.ParameterType == BrokerPriorityParameterType.ContractName) || (vResult.ParameterType == BrokerPriorityParameterType.LocalServiceName)) { vResult.ParameterValue = IdentifierOrValueExp...
// TSql160.g:22904: action may refuse the input: {if (vResult.ParameterType == BrokerPriorityParameterType.RemoteServiceName) { vResult.ParameterValue = IdentifierOrValueExpression(vParameter2); vResult.IsDefaultOrAny = BrokerPriorityParameterSpecialType.None; } else ThrowIncorrectSyntaxErrorExcepti...
// TSql160.g:22916: action may refuse the input: {if (vResult.ParameterType == BrokerPriorityParameterType.PriorityLevel) { vResult.ParameterValue = IdentifierOrValueExpression(vParameter); vResult.IsDefaultOrAny = BrokerPriorityParameterSpecialType.None; } else ThrowIncorrectSyntaxErrorException(tP...
// TSql160.g:22927: action may refuse the input: {if (vResult.ParameterType == BrokerPriorityParameterType.PriorityLevel) { vResult.IsDefaultOrAny = BrokerPriorityParameterSpecialType.Default; } else ThrowIncorrectSyntaxErrorException(tParameter);}
// TSql160.g:22937: action may refuse the input: {if (vResult.ParameterType != BrokerPriorityParameterType.PriorityLevel) { vResult.IsDefaultOrAny = BrokerPriorityParameterSpecialType.Any; } else ThrowIncorrectSyntaxErrorException(tParameter);}
brokerPriorityParameter
    : Identifier EqualsSign ( identifier | stringLiteral | integer | Default | Any )
    ;

// TSql160.g:22963
createSequenceOptionList
    : ( createSequenceOptionListElement )+
    ;

// TSql160.g:22979
// TSql160.g:22981: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Restart)}?
alterSequenceOptionListElement
    : sequenceRestartOptionListElement_RESTART
    | commonSequenceOptionListElement
    ;

// TSql160.g:22987
// TSql160.g:22990: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Start)}?
createSequenceOptionListElement
    : sequenceDatatypeOptionListElement
    | sequenceStartOptionListElement_START
    | commonSequenceOptionListElement
    ;

// TSql160.g:22997
// TSql160.g:23000: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Cache)}?
// TSql160.g:23003: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Cycle)}?
// TSql160.g:23006: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Increment)}?
// TSql160.g:23009: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.No)}?
commonSequenceOptionListElement
    : ( sequenceCacheOptionListElement_CACHE | sequenceCycleOptionListElement_CYCLE | sequenceIncrementOptionListElement_INCREMENT | sequenceNoOptionListElement_NO | sequenceMinMaxOptionListElement )
    ;

// TSql160.g:23016
sequenceDatatypeOptionListElement
    : As scalarDataType
    ;

// TSql160.g:23117
sequenceMinMaxOptionListElement
    : ( 'MINVALUE' | 'MAXVALUE' ) seedIncrement
    ;

// TSql160.g:23137
// TSql160.g:23140: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Master)}?
// TSql160.g:23143: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Encryption)}?
createColumnStatements
    : Column ( createColumnMasterKeyStatement_MASTER | createColumnEncryptionKeyStatement_ENCRYPTION )
    ;

// TSql160.g:23178
// TSql160.g:23179: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.KeyStoreProviderName)}?
// TSql160.g:23181: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.KeyPath)}?
// TSql160.g:23183: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.EnclaveComputations)}?
columnMasterkeyParameter
    : columnMasterKeyStoreProviderNameParameter_KEY_STORE_PROVIDER_NAME
    | columnMasterKeyPathParameter_KEY_PATH
    | columnMasterKeyEnclaveComputationsParameter_ENCLAVE_COMPUTATIONS
    ;

// TSql160.g:23259
createColumnEncryptionKeyValue
    : LeftParenthesis columnEncryptionKeyValueParameter Comma columnEncryptionKeyValueParameter Comma columnEncryptionKeyValueParameter RightParenthesis
    ;

// TSql160.g:23283
// TSql160.g:23284: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.ColumnMasterKey)}?
// TSql160.g:23286: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Algorithm)}?
// TSql160.g:23288: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.EncryptedValue)}?
columnEncryptionKeyValueParameter
    : columnMasterKeyNameParameter_COLUMN_MASTER_KEY
    | columnEncryptionAlgorithmNameParameter_ALGORITHM
    | encryptedValueParameter_ENCRYPTED_VALUE
    ;

// TSql160.g:23292
columnMasterKeyNameParameter
    : 'COLUMN_MASTER_KEY' EqualsSign identifier
    ;

// TSql160.g:23331
// TSql160.g:23341: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.Add)}?
// TSql160.g:23344: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.Drop)}?
alterColumnEncryptionKeyStatement
    : Column 'ENCRYPTION' Key identifier ( alterColumnEncryptionKeyAddStatement | alterColumnEncryptionKeyDropStatement )
    ;

// TSql160.g:23349
alterColumnEncryptionKeyAddStatement
    : Add 'VALUE' createColumnEncryptionKeyValue
    ;

// TSql160.g:23364
alterColumnEncryptionKeyDropStatement
    : Drop 'VALUE' LeftParenthesis columnMasterKeyNameParameter RightParenthesis
    ;

// TSql160.g:23428
createSecurityPredicateListElement
    : Add createOrAltersecurityPredicateListElement
    ;

// TSql160.g:23483
alterSecurityPredicateListElement
    : ( Add createOrAltersecurityPredicateListElement )
    | ( Alter createOrAltersecurityPredicateListElement )
    | ( Drop dropSecurityPredicateListElement )
    ;

// TSql160.g:23512
securityPolicyOptionsList
    : securityPolicyOptionListElement ( Comma securityPolicyOptionListElement )*
    ;

// TSql160.g:23531
// TSql160.g:23536: action may refuse the input: {if (TryMatch(tOption, CodeGenerationSupporter.State)) { vResult.OptionKind = SecurityPolicyOptionKind.State; } else if (TryMatch(tOption, CodeGenerationSupporter.SchemaBinding) && vParent.ActionType == SecurityPolicyActionType.Create) { vResult.Optio...
securityPolicyOptionListElement
    : ( 'STATE' | 'SCHEMABINDING' ) EqualsSign optionOnOff
    ;

// TSql160.g:23555
createOrAltersecurityPredicateListElement
    : Identifier 'PREDICATE' schemaObjectThreePartName userFunctionCall On schemaObjectTwoPartName ( securityPredicateOperation )?
    ;

// TSql160.g:23581
dropSecurityPredicateListElement
    : Identifier 'PREDICATE' On schemaObjectTwoPartName ( securityPredicateOperation )?
    ;

// TSql160.g:23600
// TSql160.g:23602: action may refuse the input: {if (vResult.SecurityPredicateType == SecurityPredicateType.Filter) { ThrowIncorrectSyntaxErrorException(tBeforeAfter); } else if (TryMatch(tBeforeAfter, CodeGenerationSupporter.Before)) { vResult.SecurityPredicateOperation = SecurityPredicateOperatio...
// TSql160.g:23624: action may refuse the input: {if (vResult.SecurityPredicateType == SecurityPredicateType.Filter) { ThrowIncorrectSyntaxErrorException(tAfter); } Match(tAfter, CodeGenerationSupporter.After); vResult.SecurityPredicateOperation = SecurityPredicateOperation.AfterInsert; UpdateTokenI...
// TSql160.g:23636: action may refuse the input: {if (vResult.SecurityPredicateType == SecurityPredicateType.Filter) { ThrowIncorrectSyntaxErrorException(tBefore); } Match(tBefore, CodeGenerationSupporter.Before); vResult.SecurityPredicateOperation = SecurityPredicateOperation.BeforeDelete; UpdateTo...
securityPredicateOperation
    : ( 'BEFORE' | 'AFTER' ) Update
    | 'AFTER' Insert
    | 'BEFORE' Delete
    ;

// TSql160.g:23795
// TSql160.g:23808: action may refuse the input: {if (vResult.OptionKind != ExternalDataSourceOptionKind.Credential) { throw GetUnexpectedTokenErrorException(tOption); } vResult.Value = IdentifierOrValueExpression(vIdentifier);}
// TSql160.g:23817: action may refuse the input: {if (vResult.OptionKind != ExternalDataSourceOptionKind.ResourceManagerLocation && vResult.OptionKind != ExternalDataSourceOptionKind.DatabaseName && vResult.OptionKind != ExternalDataSourceOptionKind.ShardMapName && vResult.OptionKind != ExternalData...
externalDataSourceLiteralOrIdentifierOption
    : Identifier EqualsSign ( identifier | stringLiteral )
    ;

// TSql160.g:23979
// TSql160.g:23992: action may refuse the input: {if (vResult.OptionKind != ExternalStreamOptionKind.DataSource && vResult.OptionKind != ExternalStreamOptionKind.FileFormat) { throw GetUnexpectedTokenErrorException(tOption); } vResult.Value = IdentifierOrValueExpression(vIdentifier);}
// TSql160.g:24002: action may refuse the input: {if (vResult.OptionKind != ExternalStreamOptionKind.DataSource && vResult.OptionKind != ExternalStreamOptionKind.FileFormat) { throw GetUnexpectedTokenErrorException(tOption); } vResult.Value = IdentifierOrValueExpression(vLiteral);}
externalStreamLiteralOrIdentifierOption
    : Identifier EqualsSign ( identifier | stringLiteral )
    ;

// TSql160.g:24013
// TSql160.g:24031: action may refuse the input: {if (TryMatch(tExternalFileFormatType, CodeGenerationSupporter.DelimitedText)) { UpdateTokenInfo(vResult, tExternalFileFormatType); vResult.FormatType = ExternalFileFormatType.DelimitedText; } else if (TryMatch(tExternalFileFormatType, CodeGenerationS...
createExternalFileFormatStatement
    : File 'FORMAT' identifier With LeftParenthesis 'FORMAT_TYPE' EqualsSign ( ( 'DELIMITEDTEXT' | 'RCFILE' | 'ORC' | 'PARQUET' | 'JSON' | 'DELTA' ) ) ( Comma externalFileFormatOptions )? RightParenthesis
    ;

// TSql160.g:24083
// TSql160.g:24086: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.FormatOptions)}?
// TSql160.g:24089: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.SerDeMethod) || NextTokenMatches(CodeGenerationSupporter.DataCompression)}?
// TSql160.g:24096: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.FormatOptions)}?
// TSql160.g:24099: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.SerDeMethod) || NextTokenMatches(CodeGenerationSupporter.DataCompression)}?
externalFileFormatOptions
    : ( externalFormatingOption_FORMAT_OPTIONS | ( externalFileFormatLiteralOptions_SERDE_METHOD | externalFileFormatLiteralOptions_DATA_COMPRESSION ) ) ( Comma ( externalFormatingOption_FORMAT_OPTIONS | ( externalFileFormatLiteralOptions_SERDE_METHOD | externalFileFormatLiteralOptions_DATA_COMPRESSION ) ) )*
    ;

// TSql160.g:24161
// TSql160.g:24163: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.UseTypeDefault)}?
externalDataFormatOption
    : externalDataFormatUseTypeDefaultOption_USE_TYPE_DEFAULT
    | externalDataFormatLiteralOption
    ;

// TSql160.g:24198
// TSql160.g:24206: action may refuse the input: {vResult.OptionKind = ExternalFileFormatOptionHelper.Instance.ParseOption(tOption); if (vResult.OptionKind != ExternalFileFormatOptionKind.FieldTerminator && vResult.OptionKind != ExternalFileFormatOptionKind.StringDelimiter && vResult.OptionKind != E...
// TSql160.g:24221: action may refuse the input: {vResult.OptionKind = ExternalFileFormatOptionHelper.Instance.ParseOption(tOption); if (vResult.OptionKind != ExternalFileFormatOptionKind.FirstRow) { throw GetUnexpectedTokenErrorException(tOption); } vResult.Value = vLiteral;}
externalDataFormatLiteralOption
    : Identifier EqualsSign ( stringLiteral | integer )
    ;

// TSql160.g:24256
dropExternalFileFormatStatement
    : File 'FORMAT' identifier
    ;

// TSql160.g:24307
spatialIndexingScheme
    : 'USING' Identifier
    ;

// TSql160.g:24319
spatialIndexOptionsList
    : LeftParenthesis spatialIndexOption ( Comma spatialIndexOption )* RightParenthesis
    ;

// TSql160.g:24341
// TSql160.g:24344: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.BoundingBox)}?
// TSql160.g:24347: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Grids)}?
// TSql160.g:24350: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.CellsPerObject)}?
spatialIndexOption
    : ( boundingBoxSpatialIndexOption_BOUNDING_BOX | gridsSpatialIndexOption_GRIDS | cellsPerObjectSpatialIndexOption_CELLS_PER_OBJECT | spatialIndexRegularOption )
    ;

// TSql160.g:24360
spatialIndexRegularOption
    : ( indexOption )
    ;

// TSql160.g:24391
boundingBoxParameter
    : LeftParenthesis ( boundingBoxParameterListByname | boundingBoxParameterListByord ) RightParenthesis
    ;

// TSql160.g:24404
boundingBoxParameterListByname
    : boundingBoxParameterByName ( Comma boundingBoxParameterByName )*
    ;

// TSql160.g:24423
boundingBoxParameterByName
    : Identifier EqualsSign signedIntegerOrReal
    ;

// TSql160.g:24435
boundingBoxParameterListByord
    : boundingBoxParameterByOrd ( Comma boundingBoxParameterByOrd )*
    ;

// TSql160.g:24451
boundingBoxParameterByOrd
    : signedIntegerOrReal
    ;

// TSql160.g:24473
gridParameter
    : LeftParenthesis ( gridParameterListByName | gridParameterListByOrd ) RightParenthesis
    ;

// TSql160.g:24486
gridParameterListByName
    : gridParameterByName ( Comma gridParameterByName )*
    ;

// TSql160.g:24502
gridParameterByName
    : Identifier EqualsSign Identifier
    ;

// TSql160.g:24511
gridParameterListByOrd
    : gridParameterByOrd ( Comma gridParameterByOrd )*
    ;

// TSql160.g:24527
gridParameterByOrd
    : Identifier
    ;

// TSql160.g:24552
createDefaultStatement
    : Create Default schemaObjectThreePartName As expression
    ;

// TSql160.g:24570
createViewStatement
    : Create viewStatementBody
    ;

// TSql160.g:24578
alterViewStatement
    : Alter viewStatementBody
    ;

// TSql160.g:24586
createOrAlterViewStatement
    : viewStatementBody
    ;

// TSql160.g:24590
viewStatementBody
    : View schemaObjectTwoPartName ( columnNameList )? ( With viewOption ( Comma viewOption )* )? As subqueryExpressionWithOptionalCTE ( With Check Option )?
    ;

// TSql160.g:24705
viewOption
    : Identifier
    ;

// TSql160.g:24714
// TSql160.g:24716: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Distribution)}?
// TSql160.g:24719: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.ForAppend)}?
createMaterializedViewOption
    : viewDistributionOption_DISTRIBUTION
    | viewForAppendOption_FOR_APPEND
    ;

// TSql160.g:24752
// TSql160.g:24754: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.RoundRobin)}?
// TSql160.g:24757: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Hash)}?
viewDistributionPolicy
    : viewRoundRobinDistributionPolicy_ROUND_ROBIN
    | viewHashDistributionPolicy_HASH
    ;

// TSql160.g:24795
columnNameList
    : LeftParenthesis identifier ( Comma identifier )* RightParenthesis
    ;

// TSql160.g:24814
methodSpecifier
    : External 'NAME' identifier Dot identifier Dot identifier
    ;

// TSql160.g:24837
triggerOption
    : Identifier
    | executeAsClause
    ;

// TSql160.g:24859
procedureOptions
    : With procedureOption ( Comma procedureOption )*
    ;

// TSql160.g:24878
procedureOption
    : Identifier
    | executeAsClause
    ;

// TSql160.g:24898
alterProcedureStatement
    : Alter procedureStatementBody
    ;

// TSql160.g:24915
createProcedureStatement
    : Create procedureStatementBody
    ;

// TSql160.g:24932
createOrAlterProcedureStatement
    : procedureStatementBody
    ;

// TSql160.g:24945
procedureStatementBody
    : ( Proc | Procedure ) procedureReference ( LeftParenthesis procedureParameterList RightParenthesis | ( procedureParameterList )? ) ( procedureOptions )? ( For Replication )? As ( optSemicolons ( statementList )? | methodSpecifier optSemicolons )
    ;

// TSql160.g:25013
statementList
    : ( statementOptSemi )+
    ;

// TSql160.g:25028
procedureParameterList
    : procedureParameter ( Comma procedureParameter )*
    ;

// TSql160.g:25044
procedureParameter
    : identifierVariable ( As )? ( cursorProcedureParameter | scalarProcedureParameter )
    ;

// TSql160.g:25059
// TSql160.g:25071: action may refuse the input: {if (nullableAllowed) vParent.Nullable=vNullableConstraintDefinition; else ThrowParseErrorException("SQL46039", vNullableConstraintDefinition, TSqlParserResource.SQL46039Message);}
// TSql160.g:25089: action may refuse the input: {if (TryMatch(tId2, CodeGenerationSupporter.ReadOnlyOld)) { vParent.Modifier = ParameterModifier.ReadOnly; } else { Match(tId2,CodeGenerationSupporter.Output,CodeGenerationSupporter.Out); if (outputAllowed) vParent.Modifier = ParameterModifier.Output;...
scalarProcedureParameter
    : scalarDataType ( nullableConstraint )? ( EqualsSign ( possibleNegativeConstantOrIdentifierWithDefault ) )? ( ( 'READONLY' | 'OUTPUT' | 'OUT' ) )?
    ;

// TSql160.g:25108
cursorProcedureParameter
    : cursorDataType ( Varying ) ( ( 'OUTPUT' | 'OUT' ) )
    ;

// TSql160.g:25132
possibleNegativeConstantOrIdentifier
    : possibleNegativeConstant
    | identifierLiteral
    ;

// TSql160.g:25138
identifierLiteral
    : Identifier
    | QuotedIdentifier
    ;

// TSql160.g:25155
possibleNegativeConstantOrIdentifierWithDefault
    : possibleNegativeConstantOrIdentifier
    | defaultLiteral
    ;

// TSql160.g:25160
possibleNegativeConstant
    : literal
    | negativeConstant
    ;

// TSql160.g:25165
negativeConstant
    : Minus subroutineParameterLiteral
    ;

// TSql160.g:25184
alterTriggerStatement
    : Alter triggerStatementBody
    ;

// TSql160.g:25201
createTriggerStatement
    : Create triggerStatementBody
    ;

// TSql160.g:25218
createOrAlterTriggerStatement
    : triggerStatementBody
    ;

// TSql160.g:25231
// TSql160.g:25269: predicate dropped: {vDdlTrigger}?
triggerStatementBody
    : Trigger schemaObjectThreePartName On triggerObject ( With triggerOption ( Comma triggerOption )* )? ( ddlTriggerMidSection | dmlTriggerMidSection ) As ( ( ( Semicolon )* statementList ) | methodSpecifier optSemicolons )
    ;

// TSql160.g:25292
dmlTriggerMidSection
    : ( For | ( 'INSTEAD' | 'AFTER' ) ( Of )? ) dmlTriggerAction ( Comma dmlTriggerAction )* ( Not For Replication )?
    ;

// TSql160.g:25346
// TSql160.g:25364: predicate dropped: {!NextTokenMatches(CodeGenerationSupporter.Logon)}?
// TSql160.g:25379: action may refuse the input: {Match(tLogon, CodeGenerationSupporter.Logon); if (vParent.TriggerObject.TriggerScope == TriggerScope.AllServer) { vAction = FragmentFactory.CreateFragment<TriggerAction>(); vAction.TriggerActionType = TriggerActionType.LogOn; AddAndUpdateTokenInfo(vP...
ddlTriggerMidSection
    : ( For | 'AFTER' ) ( ( ddlTriggerAction ( Comma ddlTriggerAction )* ) | 'LOGON' )
    ;

// TSql160.g:25394
triggerObject
    : schemaObjectThreePartName
    | triggerScope
    ;

// TSql160.g:25412
triggerScope
    : Database
    | All 'SERVER'
    ;

// TSql160.g:25429
dmlTriggerAction
    : Insert
    | Update
    | Delete
    ;

// TSql160.g:25450
ddlTriggerAction
    : Identifier
    ;

// TSql160.g:25471
endConversationStatement
    : End 'CONVERSATION' expression endConversationArgumentsOpt
    ;

// TSql160.g:25483
endConversationArgumentsOpt
    : ( With Identifier (  | EqualsSign integerOrVariable 'DESCRIPTION' EqualsSign stringOrVariable ) )?
    ;

// TSql160.g:25521
getConversationGroupStatement
    : 'GET' 'CONVERSATION' Group variable From schemaObjectThreePartName
    ;

// TSql160.g:25586
receiveSelectExpression
    : selectSetVariable
    | selectStarExpression
    | receiveColumnSelectExpression
    ;

// TSql160.g:25595
receiveColumnSelectExpression
    : expression ( ( As )? stringOrIdentifier )?
    ;

// TSql160.g:25660
beginConversationTimerStatement
    : Begin 'CONVERSATION' 'TIMER' LeftParenthesis expression RightParenthesis 'TIMEOUT' EqualsSign expression
    ;

// TSql160.g:25678
beginDialogStatement
    : Begin 'DIALOG' ( 'CONVERSATION' )? variable From 'SERVICE' identifierOrVariable To 'SERVICE' stringOrVariable ( Comma stringOrVariable )? ( On 'CONTRACT' identifierOrVariable )? ( With beginConversationArgument ( Comma beginConversationArgument )* )?
    ;

// TSql160.g:25731
beginConversationArgument
    : Identifier EqualsSign ( scalarExpressionBeginDialogConversationArgument | onOffBeginDialogConversationArgument )
    ;

// TSql160.g:25739
scalarExpressionBeginDialogConversationArgument
    : expression
    ;

// TSql160.g:25769
onOffBeginDialogConversationArgument
    : optionOnOff
    ;

// TSql160.g:25783
mergeStatement
    : mergeStatementWithoutSemicolon requiredSemicolon
    ;

// TSql160.g:25789
// TSql160.g:25798: action may refuse the input: {if(tSemi==null) { ThrowParseErrorException("SQL46097", vParent,TSqlParserResource.SQL46097Message, vStatement); }}
requiredSemicolon
    : ( Semicolon )?
    ;

// TSql160.g:25806
mergeStatementWithoutSemicolon
    : mergeSpecification ( optimizerHints )?
    ;

// TSql160.g:25820
// TSql160.g:25844: predicate dropped: {!NextTokenMatches(CodeGenerationSupporter.Using)}?
mergeSpecification
    : Merge dmlTopRowFilterOpt ( Into )? dmlTarget ( ( As identifier ) | ( identifier ) |  ) 'USING' selectTableReferenceWithOdbc On booleanExpressionWithFlags ( mergeActionClause )+ outputClauseOpt
    ;

// TSql160.g:25874
mergeActionClause
    : When mergeCondition ( And booleanExpression )? Then mergeAction
    ;

// TSql160.g:25897
mergeCondition
    : 'MATCHED'
    | ( Not 'MATCHED' ( By ( 'TARGET' | 'SOURCE' ) |  ) )
    ;

// TSql160.g:25930
mergeAction
    : updateMergeAction
    | insertMergeAction
    | deleteMergeAction
    ;

// TSql160.g:25936
// TSql160.g:25939: action may refuse the input: {if (condition == MergeCondition.NotMatched || condition == MergeCondition.NotMatchedByTarget) { ThrowParseErrorException("SQL46041", tUpdate,TSqlParserResource.SQL46041Message, "Update"); } UpdateTokenInfo(vResult,tUpdate);}
updateMergeAction
    : Update setClausesList
    ;

// TSql160.g:25948
// TSql160.g:25953: action may refuse the input: {if (condition == MergeCondition.Matched || condition == MergeCondition.NotMatchedBySource) { ThrowParseErrorException("SQL46040", tInsert,TSqlParserResource.SQL46040Message); } UpdateTokenInfo(vResult, tInsert); vResult.Source = vInsertSource;}
insertMergeAction
    : Insert mergeInsertDmlColumnListOpt mergeInsertSource
    ;

// TSql160.g:25964
// TSql160.g:25967: action may refuse the input: {if (condition == MergeCondition.NotMatched || condition == MergeCondition.NotMatchedByTarget) { ThrowParseErrorException("SQL46041", tDelete,TSqlParserResource.SQL46041Message, "Delete"); } UpdateTokenInfo(vResult, tDelete);}
deleteMergeAction
    : Delete
    ;

// TSql160.g:25976
mergeInsertSource
    : defaultValuesInsertSource
    | Values rowValueExpressionWithDefault
    ;

// TSql160.g:25991
executeAsStatement
    : execStart As executeAsStatementContext executeContextStatementOptionsOpt
    ;

// TSql160.g:26002
executeAsStatementContext
    : User EqualsSign expression
    | Identifier (  | EqualsSign expression )
    ;

// TSql160.g:26028
executeContextStatementOptionsOpt
    : ( With ( 'NO' Revert | 'COOKIE' Into variable ) )?
    ;

// TSql160.g:26055
executeStatement
    : executeSpecification ( With executeOption ( Comma executeOption )* )?
    ;

// TSql160.g:26084
// TSql160.g:26085: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Recompile)}?
executeOption
    : executeOptionRecompile_RECOMPILE
    | executeOptionResultSets
    ;

// TSql160.g:26099
executeOptionResultSets
    : 'RESULT' 'SETS' ( ( 'UNDEFINED' | 'NONE' ) | LeftParenthesis resultSetDefinition ( Comma resultSetDefinition )* RightParenthesis )
    ;

// TSql160.g:26143
resultSetDefinition
    : inlineResultSetDefinition
    | asForXmlResultSetDefinition
    | asSchemaObjectResultSetDefinition
    ;

// TSql160.g:26149
inlineResultSetDefinition
    : LeftParenthesis resultColumnDefinition ( Comma resultColumnDefinition )* RightParenthesis
    ;

// TSql160.g:26169
resultColumnDefinition
    : columnDefinitionBasic ( nullableConstraint )?
    ;

// TSql160.g:26186
asForXmlResultSetDefinition
    : As For 'XML'
    ;

// TSql160.g:26195
asSchemaObjectResultSetDefinition
    : As ( 'OBJECT' | 'TYPE' ) schemaObjectFourPartName
    ;

// TSql160.g:26214
execStart
    : Execute
    | Exec
    ;

// TSql160.g:26225
executeSpecification
    : execStart execTypes
    ;

// TSql160.g:26231
// TSql160.g:26249: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.At)}?
execTypes
    : LeftParenthesis execStrTypes RightParenthesis ( execStrExecCtxt )? ( linkedServer_AT |  )
    | execProcEx
    | variable EqualsSign execProcEx
    ;

// TSql160.g:26267
execStrExecCtxt
    : As execCtxtStmt
    ;

// TSql160.g:26271
execCtxtStmt
    : execCtxtStmtType ( execCtxtStmtPrincipal )?
    ;

// TSql160.g:26283
execCtxtStmtPrincipal
    : EqualsSign stringLiteral
    ;

// TSql160.g:26287
execStrTypes
    : execSqlList ( Comma setParamList )?
    ;

// TSql160.g:26297
execProcEx
    : execProc
    | adhocDataSourceExecproc
    ;

// TSql160.g:26302
adhocDataSourceExecproc
    : adhocDataSource Dot procObjectReference ( setParamList )?
    ;

// TSql160.g:26318
execProc
    : ( procObjectReference | varObjectReference ) ( setParamList )?
    ;

// TSql160.g:26332
procObjectReference
    : procedureReference
    ;

// TSql160.g:26342
varObjectReference
    : variable
    ;

// TSql160.g:26352
procedureReference
    : schemaObjectFourPartName procNumOpt
    ;

// TSql160.g:26371
execSqlList
    : stringOrGlobalVariableOrVariable ( Plus stringOrGlobalVariableOrVariable )*
    ;

// TSql160.g:26386
binaryOrVariable
    : binary
    | variable
    ;

// TSql160.g:26391
integerOrVariable
    : integer
    | variable
    ;

// TSql160.g:26396
integerOrRealOrNumeric
    : integer
    | real
    | numeric
    ;

// TSql160.g:26422
setParamList
    : setParam ( Comma setParam )*
    ;

// TSql160.g:26439
// TSql160.g:26454: action may refuse the input: {vResult.ParameterValue = vValue; if (vResult.Variable != null) { nameEqualsValueWasUsed = true; } else if (nameEqualsValueWasUsed) { ThrowParseErrorException("SQL46089", vValue, TSqlParserResource.SQL46089Message, parameterNumber.ToString(CultureInfo...
// TSql160.g:26470: action may refuse the input: {VariableReference parameterAsVariableReference = vValue as VariableReference; GlobalVariableExpression parameterAsGlobalVarExpr = vValue as GlobalVariableExpression; if (parameterAsVariableReference == null && parameterAsGlobalVarExpr == null) { Thro...
setParam
    : ( variable EqualsSign )? ( ( possibleNegativeConstantOrIdentifier ( ( 'OUTPUT' | 'OUT' ) )? ) | defaultLiteral )
    ;

// TSql160.g:26492
adhocDataSource
    : OpenDataSource LeftParenthesis stringLiteral Comma stringLiteral RightParenthesis
    ;

// TSql160.g:26505
procNumOpt
    : ProcNameSemicolon integer
    | /* empty */
    ;

// TSql160.g:26510
execCtxtStmtType
    : User
    | 'LOGIN'
    ;

// TSql160.g:26528
// TSql160.g:26543: predicate dropped: {(LA(1) == LeftParenthesis && (((LA(2) == Identifier || LA(2) == QuotedIdentifier) && LA(3) == Comma) || ((LA(2) == Identifier || LA(2) == QuotedIdentifier) && LA(3) == RightParenthesis)))}?
// TSql160.g:26594: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Federated)}?
createTableStatement
    : Create Table schemaObjectThreePartName ( ( ctasCreateTableStatement ) | ( LeftParenthesis tableDefinitionCreateTable RightParenthesis ( graphAsNodeOrEdge )? ) | ctasCreateTableStatement | As ( 'FILETABLE' | 'EDGE' ) | As ) ( On filegroupOrPartitionScheme )? ( federatedOn_FEDERATED )? largeDataOnOpt ( withTableOptions )?
    ;

// TSql160.g:26612
graphAsNodeOrEdge
    : As ( 'EDGE' | 'NODE' )
    ;

// TSql160.g:26628
tableDefinitionCreateTable
    : ( tableElement ( Comma tableElement )* ( Comma )? )
    ;

// TSql160.g:26670
largeDataOnOpt
    : ( largeDataOn ( largeDataOn )? )?
    ;

// TSql160.g:26682
// TSql160.g:26687: action may refuse the input: {if (TryMatch(tTextImageFileStreamOn, CodeGenerationSupporter.TextImageOn)) { if (vParent.TextImageOn != null) { ThrowParseErrorException("SQL46058", tTextImageFileStreamOn, TSqlParserResource.SQL46058Message, tTextImageFileStreamOn.getText()); } vPar...
largeDataOn
    : ( 'TEXTIMAGE_ON' | 'FILESTREAM_ON' ) stringOrIdentifier
    ;

// TSql160.g:26710
withTableOptions
    : With LeftParenthesis createTableOption ( Comma createTableOption )* RightParenthesis
    ;

// TSql160.g:26729
// TSql160.g:26743: predicate dropped: {NextTokenMatchesOneOf(new string[] {CodeGenerationSupporter.FileTablePrimaryKeyConstraintName, CodeGenerationSupporter.FileTableStreamIdUniqueConstraintName, CodeGenerationSupporter.FileTableFullPathUniqueConstraintName})}?
// TSql160.g:26731: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.DataCompression)}?
// TSql160.g:26734: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.XmlCompression)}?
// TSql160.g:26737: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.FileTableDirectory)}?
// TSql160.g:26740: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.FileTableCollateFileName)}?
// TSql160.g:26748: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MemoryOptimized)}?
// TSql160.g:26751: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Durability)}?
// TSql160.g:26754: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.RemoteDataArchive)}?
// TSql160.g:26757: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.SystemVersioning)}?
// TSql160.g:26760: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Ledger)}?
// TSql160.g:26763: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Distribution)}?
// TSql160.g:26768: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Partition)}?
// TSql160.g:26771: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Location)}?
// TSql160.g:26774: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.DataDeletion)}?
createTableOption
    : tableDataCompressionOption_DATA_COMPRESSION
    | tableXmlCompressionOption_XML_COMPRESSION
    | fileTableDirectoryTableOption_FILETABLE_DIRECTORY
    | fileTableCollateFileNameTableOption_FILETABLE_COLLATE_FILENAME
    | fileTableConstraintNameTableOption
    | memoryOptimizedTableOption_MEMORY_OPTIMIZED
    | durabilityTableOption_DURABILITY
    | remoteDataArchiveTableOption_REMOTE_DATA_ARCHIVE
    | systemVersioningTableOnOffOption_SYSTEM_VERSIONING
    | ledgerTableOption_LEDGER
    | tableDistributionOption_DISTRIBUTION
    | tableIndexOption
    | tablePartitionOption_PARTITION
    | locationOption_LOCATION
    | dataRetentionOption_DATA_DELETION
    ;

// TSql160.g:26819
dataRetentionSubOption
    : 'FILTER_COLUMN' EqualsSign identifier
    | 'RETENTION_PERIOD' EqualsSign ( 'INFINITE' | signedInteger Identifier )
    ;

// TSql160.g:26862
memoryOptimizedTableOption
    : 'MEMORY_OPTIMIZED' EqualsSign optionOnOff
    ;

// TSql160.g:26936
// TSql160.g:26938: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Replicate)}?
// TSql160.g:26941: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.RoundRobin)}?
// TSql160.g:26944: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Hash)}?
tableDistributionPolicy
    : tableReplicateDistributionPolicy_REPLICATE
    | tableRoundRobinDistributionPolicy_ROUND_ROBIN
    | tableHashDistributionPolicy_HASH
    ;

// TSql160.g:26991
tableIndexOption
    : ( tableIndexType )
    ;

// TSql160.g:27004
// TSql160.g:27006: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Heap)}?
// TSql160.g:27009: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.Clustered)}?
tableIndexType
    : tableNonClusteredIndexType_HEAP
    | tableClusteredIndexType
    ;

// TSql160.g:27021
tableClusteredIndexType
    : Clustered ( ( 'COLUMNSTORE' Index ( Order ( identifierColumnList ) )? ) | ( Index LeftParenthesis columnWithSortOrder ( Comma columnWithSortOrder )* RightParenthesis ) )
    ;

// TSql160.g:27100
// TSql160.g:27111: action may refuse the input: {Match(tParameterIdentifier, CodeGenerationSupporter.MigrationState); if(vIsMigrationStateSpecified) { ThrowIncorrectSyntaxErrorException(tParameterIdentifier); } vLocalMigrationState = MigrationStateHelper.Instance.ParseOption(tMig); vIsMigrationStat...
// TSql160.g:27123: action may refuse the input: {Match(tParameterIdentifier, CodeGenerationSupporter.FilterPredicate); CheckTwoPartNameForSchemaObjectName(vSchemaObjectName, CodeGenerationSupporter.FilterPredicate); if(vIsFilterPredicateSpecified) { ThrowIncorrectSyntaxErrorException(tParameterIden...
// TSql160.g:27140: action may refuse the input: {Match(tParameterIdentifier, CodeGenerationSupporter.FilterPredicate); if(vIsFilterPredicateSpecified) { ThrowIncorrectSyntaxErrorException(tParameterIdentifier); } vFilterPredicate = null; vIsFilterPredicateSpecified = true;}
remoteDataArchiveTableSettings
    : ( Identifier EqualsSign ( Identifier | schemaObjectThreePartName ( userFunctionCall ) | Null ) ) ( Comma remoteDataArchiveTableSettings )*
    ;

// TSql160.g:27218
// TSql160.g:27227: action may refuse the input: {Match(tHistoryTable, CodeGenerationSupporter.HistoryTable); if (historyTableSpecified) { ThrowParseErrorException("SQL46010", vParent, TSqlParserResource.SQL46010Message, CodeGenerationSupporter.HistoryTable); } historyTableSpecified = true;}
// TSql160.g:27244: action may refuse the input: {Match(tDataConsistencyCheck, CodeGenerationSupporter.DataConsistencyCheck); if (dataConsistencyCheckSpecified) { ThrowParseErrorException("SQL46010", vParent, TSqlParserResource.SQL46010Message, CodeGenerationSupporter.DataConsistencyCheck); } dataCo...
// TSql160.g:27269: action may refuse the input: {Match(tHistoryRetentionPeriod, CodeGenerationSupporter.HistoryRetentionPeriod); if (retentionPeriodSpecified) { ThrowParseErrorException("SQL46010", vParent, TSqlParserResource.SQL46010Message, CodeGenerationSupporter.HistoryRetentionPeriod); } reten...
// TSql160.g:27224: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.HistoryTable)}?
// TSql160.g:27241: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.DataConsistencyCheck)}?
// TSql160.g:27266: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.HistoryRetentionPeriod)}?
systemVersioningTableSuboption
    : ( 'HISTORY_TABLE' EqualsSign schemaObjectNonEmptyTwoPartName )
    | ( 'DATA_CONSISTENCY_CHECK' EqualsSign ( On | Off ) )
    | ( 'HISTORY_RETENTION_PERIOD' EqualsSign ( ( signedInteger Identifier ) | 'INFINITE' ) )
    ;

// TSql160.g:27336
// TSql160.g:27344: action may refuse the input: {Match(tLedger, CodeGenerationSupporter.LedgerView); if (ledgerViewSpecified) { ThrowParseErrorException("SQL46010", vParent, TSqlParserResource.SQL46010Message, CodeGenerationSupporter.LedgerView); } ledgerViewSpecified = true;}
// TSql160.g:27369: action may refuse the input: {Match(tAppendOnly, CodeGenerationSupporter.AppendOnly); if (appendOnlySpecified) { ThrowParseErrorException("SQL46010", vParent, TSqlParserResource.SQL46010Message, CodeGenerationSupporter.AppendOnly); } appendOnlySpecified = true;}
// TSql160.g:27341: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.LedgerView)}?
// TSql160.g:27366: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.AppendOnly)}?
ledgerTableSuboption
    : ( 'LEDGER_VIEW' EqualsSign schemaObjectNonEmptyTwoPartName ( LeftParenthesis ledgerViewSuboption ( Comma ledgerViewSuboption )* RightParenthesis )? )
    | ( 'APPEND_ONLY' EqualsSign ( Off | On ) )
    ;

// TSql160.g:27392
// TSql160.g:27403: action may refuse the input: {Match(tTransactionIdColumnName, CodeGenerationSupporter.TransactionIdColumnName); if (transactionIdSpecified) { ThrowParseErrorException("SQL46010", vParent, TSqlParserResource.SQL46010Message, CodeGenerationSupporter.TransactionIdColumnName); } tran...
// TSql160.g:27421: action may refuse the input: {Match(tSequenceNumberColumnName, CodeGenerationSupporter.SequenceNumberColumnName); if (sequenceNumberSpecified) { ThrowParseErrorException("SQL46010", vParent, TSqlParserResource.SQL46010Message, CodeGenerationSupporter.SequenceNumberColumnName); } ...
// TSql160.g:27439: action may refuse the input: {Match(tOperationTypeColumnName, CodeGenerationSupporter.OperationTypeColumnName); if (operationTypeSpecified) { ThrowParseErrorException("SQL46010", vParent, TSqlParserResource.SQL46010Message, CodeGenerationSupporter.OperationTypeColumnName); } oper...
// TSql160.g:27457: action may refuse the input: {Match(tOperationTypeDescColumnName, CodeGenerationSupporter.OperationTypeDescColumnName); if (operationTypeIdDescSpecified) { ThrowParseErrorException("SQL46010", vParent, TSqlParserResource.SQL46010Message, CodeGenerationSupporter.OperationTypeDescC...
// TSql160.g:27400: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.TransactionIdColumnName)}?
// TSql160.g:27418: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.SequenceNumberColumnName)}?
// TSql160.g:27436: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.OperationTypeColumnName)}?
// TSql160.g:27454: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.OperationTypeDescColumnName)}?
ledgerViewSuboption
    : ( 'TRANSACTION_ID_COLUMN_NAME' EqualsSign identifier )
    | ( 'SEQUENCE_NUMBER_COLUMN_NAME' EqualsSign identifier )
    | ( 'OPERATION_TYPE_COLUMN_NAME' EqualsSign identifier )
    | ( 'OPERATION_TYPE_DESC_COLUMN_NAME' EqualsSign identifier )
    ;

// TSql160.g:27473
// TSql160.g:27478: action may refuse the input: {vResult.Value=vConstraint; if (TryMatch(tFileTableConstraintName, CodeGenerationSupporter.FileTablePrimaryKeyConstraintName)) { vResult.OptionKind=TableOptionKind.FileTablePrimaryKeyConstraintName; } else if (TryMatch(tFileTableConstraintName, CodeGe...
fileTableConstraintNameTableOption
    : ( 'FILETABLE_PRIMARY_KEY_CONSTRAINT_NAME' | 'FILETABLE_STREAMID_UNIQUE_CONSTRAINT_NAME' | 'FILETABLE_FULLPATH_UNIQUE_CONSTRAINT_NAME' ) EqualsSign identifier
    ;

// TSql160.g:27623
compressionPartitionRange
    : expression ( To expression )?
    ;

// TSql160.g:27638
createExternalTableStatement
    : Table schemaObjectThreePartName ( ( LeftParenthesis externalTableColumnElement ( Comma externalTableColumnElement )* RightParenthesis withExternalTableOptions ) | ctasCreateExternalTableStatement )
    ;

// TSql160.g:27669
externalTableColumnElement
    : columnDefinitionCreateExternalTable
    ;

// TSql160.g:27679
// TSql160.g:27687: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.DataSource)}?
// TSql160.g:27700: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.DataSource)}?
withExternalTableOptions
    : With LeftParenthesis ( externalTableDataSource_DATA_SOURCE | createExternalTableOption ) ( Comma ( externalTableDataSource_DATA_SOURCE | createExternalTableOption ) )* RightParenthesis
    ;

// TSql160.g:27729
// TSql160.g:27731: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.RejectType)}?
// TSql160.g:27734: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Distribution)}?
// TSql160.g:27737: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.RejectSampleValue) || NextTokenMatches(CodeGenerationSupporter.RejectValue)}?
createExternalTableOption
    : externalTableRejectTypeOption_REJECT_TYPE
    | externalTableDistributionOption_DISTRIBUTION
    | ( externalTableNumericLiteralOption_REJECT_SAMPLE_VALUE | externalTableNumericLiteralOption_REJECT_VALUE )
    | externalTableLiteralOrIdentifierOption
    ;

// TSql160.g:27743
// TSql160.g:27757: action may refuse the input: {if (vResult.OptionKind != ExternalTableOptionKind.FileFormat) { throw GetUnexpectedTokenErrorException(tOption); } vResult.Value = IdentifierOrValueExpression(vIdentifier);}
// TSql160.g:27766: action may refuse the input: {if (vResult.OptionKind != ExternalTableOptionKind.Location && vResult.OptionKind != ExternalTableOptionKind.SchemaName && vResult.OptionKind != ExternalTableOptionKind.ObjectName && vResult.OptionKind != ExternalTableOptionKind.RejectedRowLocation &&...
externalTableLiteralOrIdentifierOption
    : Identifier EqualsSign ( identifier | stringLiteral )
    ;

// TSql160.g:27851
// TSql160.g:27853: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Replicated)}?
// TSql160.g:27856: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.RoundRobin)}?
// TSql160.g:27859: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Sharded)}?
externalTableDistributionPolicy
    : externalTableReplicatedDistributionPolicy_REPLICATED
    | externalTableRoundRobinDistributionPolicy_ROUND_ROBIN
    | externalTableShardedDistributionPolicy_SHARDED
    ;

// TSql160.g:27895
// TSql160.g:27922: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.ChangeTracking, 2)}?
// TSql160.g:27925: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.FileTableNamespace, 2)}?
// TSql160.g:27906: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Switch)}?
// TSql160.g:27916: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Rebuild)}?
// TSql160.g:27919: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.Merge) || NextTokenMatches(CodeGenerationSupporter.Split)}?
alterTableStatement
    : Alter Table schemaObjectThreePartName ( alterTableAlterColumnStatement | alterTableAlterIndexStatement | alterTableTriggerModificationStatement | alterTableDropTableElementStatement | alterTableSwitchStatement_SWITCH | ( ( With constraintEnforcement )? ( alterTableAddTableElementStatement | alterTableConstraintModificationStatement ) ) | alterTableRebuildStatement_REBUILD | alterTableAlterPartitionStatement | alterTableChangeTrackingModificationStatement | alterTableFileTableNamespaceStatement | alterTableSetStatement )
    ;

// TSql160.g:27945
constraintEnforcement
    : Check
    | NoCheck
    ;

// TSql160.g:27985
alterTableAlterPartitionStatement
    : ( 'SPLIT' | Merge ) 'RANGE' LeftParenthesis expression RightParenthesis
    ;

// TSql160.g:28016
// TSql160.g:28030: action may refuse the input: {Match(tTrackColumnsUpdated, CodeGenerationSupporter.TrackColumnsUpdated); if (!vResult.IsEnable) { ThrowIncorrectSyntaxErrorException(tEnableDisable); } vResult.TrackColumnsUpdated = vOptionState; UpdateTokenInfo(vResult, tRParen);}
alterTableChangeTrackingModificationStatement
    : Identifier 'CHANGE_TRACKING' ( With LeftParenthesis 'TRACK_COLUMNS_UPDATED' EqualsSign optionOnOff RightParenthesis )?
    ;

// TSql160.g:28042
alterTableFileTableNamespaceStatement
    : Identifier 'FILETABLE_NAMESPACE'
    ;

// TSql160.g:28051
alterTableSetStatement
    : Set LeftParenthesis tableOption ( Comma tableOption )* RightParenthesis
    ;

// TSql160.g:28070
// TSql160.g:28072: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.LockEscalation)}?
// TSql160.g:28075: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.FileStreamOn)}?
// TSql160.g:28078: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.FileTableDirectory)}?
// TSql160.g:28081: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.SystemVersioning)}?
// TSql160.g:28084: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.RemoteDataArchive)}?
// TSql160.g:28087: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.DataDeletion)}?
tableOption
    : lockEscalationTableOption_LOCK_ESCALATION
    | fileStreamOnTableOption_FILESTREAM_ON
    | fileTableDirectoryTableOption_FILETABLE_DIRECTORY
    | systemVersioningTableOnOffOption_SYSTEM_VERSIONING
    | remoteDataArchiveAlterTableOption_REMOTE_DATA_ARCHIVE
    | dataRetentionOption_DATA_DELETION
    ;

// TSql160.g:28125
alterTableAlterIndexStatement
    : Alter Index identifier Identifier With ( alterTableAlterIndexOptionList )
    ;

// TSql160.g:28172
remoteDtatArchiveTableOnOffOption
    : Off
    | On
    | 'OFF_WITHOUT_DATA_RECOVERY'
    ;

// TSql160.g:28194
// TSql160.g:28310: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.Masked) && !vAddDefined}?
// TSql160.g:28215: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Generated)}?
// TSql160.g:28219: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Masked)}?
// TSql160.g:28223: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Encrypted)}?
// TSql160.g:28231: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Sparse) || NextTokenMatches(CodeGenerationSupporter.FileStream) || NextTokenMatches(CodeGenerationSupporter.ColumnSet)}?
// TSql160.g:28239: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Hidden)}?
// TSql160.g:28282: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Persisted)}?
// TSql160.g:28292: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Sparse)}?
// TSql160.g:28302: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Masked) && vAddDefined}?
// TSql160.g:28319: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Hidden)}?
alterTableAlterColumnStatement
    : Alter Column identifier ( scalarDataType collationOpt ( generatedAlwaysClauseInAlter_GENERATED )? ( maskedClause_MASKED )? ( columnEncryptionClause_ENCRYPTED )? ( ( columnStorage_SPARSE | columnStorage_FILESTREAM | columnStorage_COLUMN_SET ) )? ( 'HIDDEN' (  ) )? ( nullNotNull )? | ( Add | Drop ) ( RowGuidColumn | Not For Replication | ( 'PERSISTED' ) | ( 'SPARSE' ) | ( maskedClause_MASKED ) | ( 'MASKED' ) | ( 'HIDDEN' ) ) ) ( With indexOptionList )?
    ;

// TSql160.g:28341
alterTableTriggerModificationStatement
    : Identifier Trigger ( All | identifier ( Comma identifier )* )
    ;

// TSql160.g:28373
alterTableDropTableElementStatement
    : Drop alterTableDropTableElement ( Comma alterTableDropTableElement )*
    ;

// TSql160.g:28393
alterTableDropTableElement
    : ( 'PERIOD' For 'SYSTEM_TIME' )
    | ( Constraint ( If Exists )? )? identifier ( alterTableDropTableElementWithOptions )?
    | Column ( If Exists )? identifier
    | Index identifier
    ;

// TSql160.g:28454
alterTableDropTableElementWithOptions
    : With LeftParenthesis dropClusteredConstraintOption ( Comma dropClusteredConstraintOption )* RightParenthesis
    ;

// TSql160.g:28476
// TSql160.g:28485: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.MaxDop)}?
// TSql160.g:28494: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Online)}?
// TSql160.g:28506: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Move)}?
dropClusteredConstraintOption
    : 'MAXDOP' EqualsSign integer
    | 'ONLINE' EqualsSign optionOnOff
    | 'MOVE' To filegroupOrPartitionScheme
    | lowPriorityLockWaitOption
    ;

// TSql160.g:28559
tableSwitchOptionList
    : lowPriorityLockWaitOption
    | truncateTargetOption
    ;

// TSql160.g:28579
alterTableAddTableElementStatement
    : Add tableDefinition
    ;

// TSql160.g:28592
alterTableConstraintModificationStatement
    : constraintEnforcement Constraint ( All | identifier ( Comma identifier )* )
    ;

// TSql160.g:28624
// TSql160.g:28658: action may refuse the input: {if(vIsHash) { if(vIsClustered.HasValue && vIsClustered.Value) { ThrowIncorrectSyntaxErrorException(tHashOrColumnStore); } vIndexTypeKind = IndexTypeKind.NonClusteredHash; } else if(vIsClustered.HasValue) { if(vIsClustered.Value) { vIndexTypeKind = vI...
inlineIndexType
    : ( Clustered | NonClustered )? ( ( 'HASH' | 'COLUMNSTORE' ) )?
    ;

// TSql160.g:28690
// TSql160.g:28739: action may refuse the input: {if(vIndexType != null && vIndexType.IndexTypeKind != null) { if((vIndexType.IndexTypeKind != IndexTypeKind.ClusteredColumnStore && !vHasColumnList) || (vIndexType.IndexTypeKind == IndexTypeKind.ClusteredColumnStore && vHasColumnList)) { ThrowIncorrec...
// TSql160.g:28735: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Include)}?
inlineIndexTableDefinition
    : Index identifier ( Unique )? inlineIndexType ( LeftParenthesis indexKeyColumn ( Comma indexKeyColumn )* RightParenthesis )? ( 'INCLUDE' identifierColumnList )? ( filterClause )? ( With ( indexOptionList ) )? ( On filegroupOrPartitionScheme )? fileStreamOnOpt
    ;

// TSql160.g:28778
tablePeriodDefinition
    : 'PERIOD' For 'SYSTEM_TIME' LeftParenthesis identifier Comma identifier RightParenthesis
    ;

// TSql160.g:28800
// TSql160.g:28826: action may refuse the input: {if(vIndexType != null && vIndexType.IndexTypeKind != null) { // Inline COLUMNSTORE indexes are not supported in the column definition switch(vIndexType.IndexTypeKind) { case IndexTypeKind.ClusteredColumnStore: case IndexTypeKind.NonClusteredColumnSto...
// TSql160.g:28857: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Include)}?
inlineIndexColumnDefinition
    : Index ( identifier ) ( Unique )? inlineIndexType ( LeftParenthesis columnWithSortOrder ( Comma columnWithSortOrder )* RightParenthesis )? ( 'INCLUDE' identifierColumnList )? ( filterClause )? ( With ( indexOptionList ) )? ( On filegroupOrPartitionScheme )? fileStreamOnOpt
    ;

// TSql160.g:28890
tableElement
    : tableConstraint
    | columnDefinition
    | tableConstraint
    | inlineIndexTableDefinition
    | tablePeriodDefinition
    ;

// TSql160.g:28920
columnDefinition
    : identifier ( ( computedColumnBody computedColumnConstraintListOpt ) | ( regularColumnBody columnConstraintListOpt ) )
    ;

// TSql160.g:28953
columnDefinitionCreateExternalTable
    : columnDefinitionBasic ( nullableConstraint )?
    ;

// TSql160.g:28970
// TSql160.g:28979: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Persisted)}?
computedColumnBody
    : As expressionWithFlags ( 'PERSISTED' )?
    ;

// TSql160.g:28988
// TSql160.g:29001: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Generated)}?
// TSql160.g:29005: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Masked)}?
// TSql160.g:29009: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Encrypted)}?
// TSql160.g:29018: predicate written as its words: {(NextTokenMatches(CodeGenerationSupporter.Sparse) && NextTokenMatches(CodeGenerationSupporter.FileStream, 2)) || (NextTokenMatches(CodeGenerationSupporter.FileStream) && NextTokenMatches(CodeGenerationSupporter.Sparse, 2))}?
// TSql160.g:29024: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Sparse)}?
// TSql160.g:29030: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Encrypted)}?
// TSql160.g:29036: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.ColumnSet)}?
// TSql160.g:29041: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.FileStream)}?
// TSql160.g:29049: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Masked)}?
// TSql160.g:29054: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Hidden)}?
regularColumnBody
    : ( scalarDataType collationOpt ( generatedAlwaysClause_GENERATED )? ( maskedClause_MASKED )? ( columnEncryptionClause_ENCRYPTED )? ( ( sparseAndFilestreamStorage_SPARSE | sparseAndFilestreamStorage_FILESTREAM ) | sparseStorage_SPARSE ( columnEncryptionClause_ENCRYPTED )? | allSparseStorage_COLUMN_SET | filestreamStorage_FILESTREAM )? ( maskedClause_MASKED )? )? ( 'HIDDEN' (  ) )?
    ;

// TSql160.g:29193
fileStreamOnOpt
    : ( fileStreamOn )?
    ;

// TSql160.g:29200
fileStreamOn
    : 'FILESTREAM_ON' stringOrIdentifier
    ;

// TSql160.g:29212
computedColumnConstraintListOpt
    : ( columnConstraint )*
    ;

// TSql160.g:29223
// TSql160.g:29232: action may refuse the input: {if (vResult.IdentityOptions != null) { ThrowParseErrorException("SQL46043", vIdentityOptions, TSqlParserResource.SQL46043Message); } vResult.IdentityOptions=vIdentityOptions;}
columnConstraintListOpt
    : ( rowguidcolConstraint | identityConstraint | columnConstraint | inlineIndexColumnDefinition )*
    ;

// TSql160.g:29404
// TSql160.g:29405: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.ColumnEncryptionKey)}?
// TSql160.g:29407: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.EncryptionType)}?
// TSql160.g:29409: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Algorithm)}?
columnEncryptionDefinitionParameter
    : columnEncryptionKeyNameParameter_COLUMN_ENCRYPTION_KEY
    | columnEncryptionTypeParameter_ENCRYPTION_TYPE
    | columnEncryptionAlgorithmParameter_ALGORITHM
    ;

// TSql160.g:29480
// TSql160.g:29482: action may refuse the input: {if (vParent.IsRowGuidCol) ThrowParseErrorException("SQL46042", tRowguidcol, TSqlParserResource.SQL46042Message); UpdateTokenInfo(vParent,tRowguidcol); vParent.IsRowGuidCol = true;}
rowguidcolConstraint
    : RowGuidColumn
    ;

// TSql160.g:29491
identityConstraint
    : Identity ( LeftParenthesis seedIncrement Comma seedIncrement RightParenthesis )? replicationClauseOpt
    ;

// TSql160.g:29522
columnDefinitionBasic
    : identifier scalarDataType collationOpt
    ;

// TSql160.g:29536
columnDefinitionEx
    : columnDefinitionBasic
    | 'TIMESTAMP'
    ;

// TSql160.g:29551
columnConstraint
    : ( Constraint identifier )? ( nullableConstraint | defaultColumnConstraint | uniqueColumnConstraint | foreignKeyColumnConstraint | checkConstraint )
    ;

// TSql160.g:29587
tableConstraint
    : ( Constraint identifier )? ( uniqueTableConstraint | defaultTableConstraint | foreignKeyTableConstraint | checkConstraint | graphConnectionConstraint )
    ;

// TSql160.g:29628
nullableConstraint
    : nullNotNull
    ;

// TSql160.g:29639
// TSql160.g:29652: action may refuse the input: {if (statementType != IndexAffectingStatement.AlterTableAddElement) ThrowParseErrorException("SQL46013", tDefault, TSqlParserResource.SQL46013Message); UpdateTokenInfo(vResult,tValues); vResult.WithValues = true;}
defaultColumnConstraint
    : Default expressionWithFlags ( With Values )?
    ;

// TSql160.g:29663
// TSql160.g:29669: action may refuse the input: {if (statementType != IndexAffectingStatement.AlterTableAddElement) ThrowParseErrorException("SQL46014", tDefault, TSqlParserResource.SQL46014Message); UpdateTokenInfo(vResult,tDefault); ThrowConstraintIfPhaseOne(vResult);}
defaultTableConstraint
    : Default expressionWithFlags For identifier ( With Values )?
    ;

// TSql160.g:29693
uniqueColumnConstraint
    : uniqueConstraintHeader ( LeftParenthesis columnWithSortOrder ( Comma columnWithSortOrder )* RightParenthesis )? ( uniqueConstraintEnforcement | uniqueConstraintTailOpt )
    ;

// TSql160.g:29723
uniqueTableConstraint
    : uniqueConstraintHeader LeftParenthesis columnWithSortOrder ( Comma columnWithSortOrder )* RightParenthesis ( uniqueConstraintEnforcement | uniqueConstraintTailOpt )
    ;

// TSql160.g:29749
// TSql160.g:29752: action may refuse the input: {if (TryMatch(tEnforced, CodeGenerationSupporter.Enforced)) { UpdateTokenInfo(vParent, tEnforced); vParent.IsEnforced = false; } else { ThrowParseErrorException("SQL46010", tEnforced, TSqlParserResource.SQL46010Message, tEnforced.getText()); }}
uniqueConstraintEnforcement
    : ( Not 'ENFORCED' )
    ;

// TSql160.g:29767
uniqueConstraintHeader
    : ( Unique | Primary Key ) ( Clustered | NonClustered 'HASH' | NonClustered )?
    ;

// TSql160.g:29806
// TSql160.g:29813: action may refuse the input: {if (statementType == IndexAffectingStatement.CreateType) { ThrowIncorrectSyntaxErrorException(tOn); } vParent.OnFileGroupOrPartitionScheme = vFileGroupOrPartitionScheme;}
uniqueConstraintTailOpt
    : uniqueConstraintIndexOptionsOpt ( On filegroupOrPartitionScheme )? fileStreamOnOpt
    ;

// TSql160.g:29825
// TSql160.g:29838: predicate dropped: {NextTokenMatchesOneOf(CodeGenerationSupporter.SortedData, CodeGenerationSupporter.SortedDataReorg)}?
uniqueConstraintIndexOptionsOpt
    : ( With ( fillFactorOption ( sortedDataOptions )? | sortedDataOptions ( fillFactorOption )? | indexOptionList ) )?
    ;

// TSql160.g:29855
sortedDataOptions
    : ( 'SORTED_DATA' | 'SORTED_DATA_REORG' )
    ;

// TSql160.g:29862
columnWithSortOrder
    : identifierColumnReferenceExpression ( orderByOption )?
    ;

// TSql160.g:29878
identifierColumnReferenceExpression
    : multiPartIdentifier
    ;

// TSql160.g:29890
deleteUpdateAction
    : 'NO' 'ACTION'
    | Cascade
    | Set ( Null | Default )
    ;

// TSql160.g:29924
// TSql160.g:29956: action may refuse the input: {if (vOnDeleteParsed) throw GetUnexpectedTokenErrorException(tOn); vParent.DeleteAction = vDeleteUpdateAction;}
// TSql160.g:29965: action may refuse the input: {if (TryMatch(tEnforced, CodeGenerationSupporter.Enforced)) { UpdateTokenInfo(vParent, tEnforced); vParent.IsEnforced = false; } else { ThrowParseErrorException("SQL46010", tEnforced, TSqlParserResource.SQL46010Message, tEnforced.getText()); }}
foreignKeyConstraintCommonEnd
    : References schemaObjectThreePartName ( columnNameList )? ( On Delete deleteUpdateAction )? ( On Update deleteUpdateAction ( On Delete deleteUpdateAction )? )? ( Not 'ENFORCED' )? replicationClauseOpt
    ;

// TSql160.g:29984
replicationClauseOpt
    : ( Not For Replication )?
    ;

// TSql160.g:29998
foreignKeyColumnConstraint
    : ( Foreign Key foreignConstraintColumnsOpt )? foreignKeyConstraintCommonEnd
    ;

// TSql160.g:30012
foreignKeyTableConstraint
    : Foreign Key foreignConstraintColumnsOpt foreignKeyConstraintCommonEnd
    ;

// TSql160.g:30024
// TSql160.g:30036: action may refuse the input: {vResult.DeleteAction = vDeleteUpdateAction; // The only two supported actions for edge constraints are NO ACTION // and ON DELETE CASCADE. If the user does not specify an action, the // action defaults to NO ACTION. // if (vResult.DeleteAction != Del...
graphConnectionConstraint
    : 'CONNECTION' graphConnectionFromNodeToNodeList ( On Delete deleteUpdateAction )?
    ;

// TSql160.g:30053
graphConnectionFromNodeToNodeList
    : LeftParenthesis graphConnectionFromToNodeItem ( Comma graphConnectionFromToNodeItem )* RightParenthesis
    ;

// TSql160.g:30072
graphConnectionFromToNodeItem
    : schemaObjectThreePartName To schemaObjectThreePartName
    ;

// TSql160.g:30091
foreignConstraintColumnsOpt
    : ( columnNameList )?
    ;

// TSql160.g:30099
checkConstraint
    : Check replicationClauseOpt LeftParenthesis booleanExpressionWithFlags RightParenthesis
    ;

// TSql160.g:30123
cursorDataType
    : Cursor
    ;

// TSql160.g:30131
xmlDataType
    : ( LeftParenthesis ( Identifier )? dataTypeSchemaObjectName RightParenthesis )?
    ;

// TSql160.g:30159
// TSql160.g:30177: predicate dropped: {isXmlDataType}?
// TSql160.g:30180: predicate dropped: {typeOption != SqlDataTypeOption.None}?
scalarDataType
    : ( dataTypeSchemaObjectName ( xmlDataType | sqlDataTypeWithoutNational | userDataType ) )
    | doubleDataType
    | sqlDataTypeWithNational
    ;

// TSql160.g:30192
// TSql160.g:30198: action may refuse the input: {vResult.SqlDataTypeOption = ParseDataType100(vName.BaseIdentifier.Value); if (vResult.SqlDataTypeOption == SqlDataTypeOption.None || (vName.SchemaIdentifier != null && !IsSys(vName.SchemaIdentifier))) { ThrowParseErrorException("SQL46003", tNational,...
sqlDataTypeWithNational
    : National dataTypeSchemaObjectName ( Varying )? dataTypeParametersOpt
    ;

// TSql160.g:30224
sqlDataTypeWithoutNational
    : ( Varying )? dataTypeParametersOpt
    ;

// TSql160.g:30245
userDataType
    : dataTypeParametersOpt
    ;

// TSql160.g:30253
dataTypeParametersOpt
    : ( LeftParenthesis ( max | integer ( Comma integer )? ) RightParenthesis )?
    ;

// TSql160.g:30283
doubleDataType
    : Double 'PRECISION' ( LeftParenthesis ( max | integer ) RightParenthesis )?
    ;

// TSql160.g:30314
multiPartIdentifier
    : identifierList
    ;

// TSql160.g:30325
identifierList
    : identifier ( identifierListElement )*
    | identifierListElement ( identifierListElement )*
    ;

// TSql160.g:30346
identifierListElement
    : Dot ( Dot )* identifier
    ;

// TSql160.g:30372
dataTypeSchemaObjectName
    : identifier ( Dot identifier )?
    ;

// TSql160.g:30390
schemaObjectNonEmptyTwoPartName
    : identifier ( Dot identifier )?
    ;

// TSql160.g:30407
schemaObjectTwoPartName
    : schemaObjectNonEmptyTwoPartName
    | Dot identifier
    ;

// TSql160.g:30422
eventSessionNonEmptyThreePartObjectName
    : nonEmptyThreePartObjectName
    ;

// TSql160.g:30433
nonEmptyThreePartObjectName
    : identifier Dot identifier ( Dot identifier )?
    ;

// TSql160.g:30454
schemaObjectThreePartName
    : identifierList
    ;

// TSql160.g:30464
eventSessionOneOrTwoOrThreePartObjectName
    : oneOrTwoOrThreePartObjectName
    ;

// TSql160.g:30475
oneOrTwoOrThreePartObjectName
    : identifier ( Dot identifier ( Dot identifier )? )?
    ;

// TSql160.g:30498
schemaObjectFourPartName
    : identifierList
    ;

// TSql160.g:30508
booleanExpression
    : booleanExpressionWithFlags
    ;

// TSql160.g:30512
booleanExpressionWithFlags
    : booleanExpressionOr
    ;

// TSql160.g:30516
booleanExpressionOr
    : booleanExpressionAnd ( Or booleanExpressionAnd )*
    ;

// TSql160.g:30529
booleanExpressionAnd
    : booleanExpressionUnary ( And booleanExpressionUnary )*
    ;

// TSql160.g:30542
booleanExpressionUnary
    : Not booleanExpressionUnary
    | booleanExpressionPrimary
    ;

// TSql160.g:30556
// TSql160.g:30566: predicate dropped: {IsNextRuleBooleanParenthesis()}?
// TSql160.g:30571: action may refuse the input: {if (ExpressionFlags.MatchClauseAllowed != (expressionFlags & ExpressionFlags.MatchClauseAllowed)) { ThrowParseErrorException("SQL46010", tMatch, TSqlParserResource.SQL46010Message, CodeGenerationSupporter.GraphMatch); } GraphMatchPredicate vMatchPred...
// TSql160.g:30569: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.GraphMatch)}?
booleanExpressionPrimary
    : booleanExpressionParenthesis
    | 'MATCH' LeftParenthesis graphMatchPattern RightParenthesis
    | expressionWithFlags ( comparisonOperator ( comparisonPredicate | subqueryComparisonPredicate ) | Is ( booleanExpressionNullPredicate | ( Not )? Distinct From ( booleanExpressionNullPredicateFromDistinctPredicate | booleanExpressionDistinctPredicate | subqueryComparisonDistinctPredicate ) ) | ( Not )? ( inPredicate | betweenPredicate | likePredicate ) )
    | fulltextPredicate
    | existsPredicate
    | tsEqualCall
    | updateCall
    ;

// TSql160.g:30641
graphMatchPattern
    : graphMatchExpressionPattern ( And graphMatchExpressionPattern )*
    ;

// TSql160.g:30649
// TSql160.g:30652: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.ShortestPath)}?
// TSql160.g:30655: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.LastNode)}?
graphMatchExpressionPattern
    : ( graphMatchRecursivePattern_SHORTEST_PATH | graphMatchLastNodePattern_LAST_NODE | graphCompositeMatchPattern )
    ;

// TSql160.g:30716
graphRecursiveMatchQuantifier
    : ( Plus | LeftCurly integer Comma integer RightCurly )
    ;

// TSql160.g:30804
// TSql160.g:30812: action may refuse the input: {Match(tLast_Node, CodeGenerationSupporter.LastNode); if (!vCanUseLastNode) { ThrowIncorrectSyntaxErrorException(tLast_Node); } UpdateTokenInfo(vLastNodeExpression,tLast_Node); vLastNodeExpression.UsesLastNode = true; vLastNodeExpression.Node = vNode;}
graphMatchNodeExpression
    : ( 'LAST_NODE' LeftParenthesis identifier RightParenthesis | identifier )
    ;

// TSql160.g:30831
graphCompositeMatchPattern
    : graphMatchNodeExpression graphSingleMatchPattern ( graphSingleMatchPattern )*
    ;

// TSql160.g:30859
graphMatchSPRightRecursivePattern
    : ( graphSingleMatchPattern )*
    ;

// TSql160.g:30873
graphMatchSPLeftRecursivePattern
    : ( graphSingleLeftMatchPattern )*
    ;

// TSql160.g:30887
graphSingleLeftMatchPattern
    : ( ( graphMatchNodeExpression Minus LeftParenthesis identifier RightParenthesis Minus GreaterThan ) | ( graphMatchNodeExpression LessThan Minus LeftParenthesis identifier RightParenthesis Minus ) )
    ;

// TSql160.g:30915
graphSingleMatchPattern
    : ( ( Minus LeftParenthesis identifier RightParenthesis Minus GreaterThan graphMatchNodeExpression ) | ( LessThan Minus LeftParenthesis identifier RightParenthesis Minus graphMatchNodeExpression ) )
    ;

// TSql160.g:30945
tsEqualCall
    : TSEqual LeftParenthesis expression Comma expression RightParenthesis
    ;

// TSql160.g:30962
updateCall
    : Update LeftParenthesis identifier RightParenthesis
    ;

// TSql160.g:30974
fulltextPredicate
    : ( Contains | FreeText ) LeftParenthesis ( fulltextColumn | LeftParenthesis ( starColumn | column ( Comma column )* ) RightParenthesis | 'PROPERTY' LeftParenthesis identifierColumnReferenceExpression Comma stringLiteral RightParenthesis ) Comma stringOrVariable ( Comma languageExpression )? RightParenthesis
    ;

// TSql160.g:31049
fulltextColumn
    : starColumn
    | column
    ;

// TSql160.g:31057
booleanExpressionParenthesis
    : LeftParenthesis booleanExpressionWithFlags RightParenthesis
    ;

// TSql160.g:31070
existsPredicate
    : Exists subquery
    ;

// TSql160.g:31080
comparisonPredicate
    : expressionWithFlags
    ;

// TSql160.g:31093
subqueryComparisonPredicate
    : subqueryComparisonPredicateType subquery
    ;

// TSql160.g:31109
subqueryComparisonDistinctPredicate
    : subqueryComparisonPredicateType subquery
    ;

// TSql160.g:31136
// TSql160.g:31148: predicate dropped: {IsNextRuleSelectParenthesis()}?
inPredicate
    : In ( subquery | LeftParenthesis expressionList RightParenthesis )
    ;

// TSql160.g:31161
betweenPredicate
    : Between expressionWithFlags And expressionWithFlags
    ;

// TSql160.g:31186
likePredicate
    : Like expressionWithFlags ( escapeExpression | LeftCurly escapeExpression RightCurly )?
    ;

// TSql160.g:31219
booleanExpressionDistinctPredicate
    : expressionWithFlags
    ;

// TSql160.g:31232
booleanExpressionNullPredicateFromDistinctPredicate
    : Null
    ;

// TSql160.g:31242
booleanExpressionNullPredicate
    : nullNotNull
    ;

// TSql160.g:31253
escapeExpression
    : Escape expressionWithFlags
    ;

// TSql160.g:31264
comparisonOperator
    : EqualsSign
    | GreaterThan ( EqualsSign )?
    | LessThan ( EqualsSign | GreaterThan )?
    | Bang ( EqualsSign | LessThan | GreaterThan )
    ;

// TSql160.g:31311
subqueryComparisonPredicateType
    : All
    | Any
    | Some
    ;

// TSql160.g:31328
expression
    : expressionWithFlags
    ;

// TSql160.g:31332
expressionWithFlags
    : expressionBinaryPri2
    ;

// TSql160.g:31336
// TSql160.g:31366: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.At) && NextTokenMatches(CodeGenerationSupporter.Time, 2) && NextTokenMatches(CodeGenerationSupporter.Zone, 3)}?
expressionBinaryPri2
    : expressionBinaryPri1 ( Plus expressionBinaryPri1 | Minus expressionBinaryPri1 | Ampersand expressionBinaryPri1 | VerticalLine expressionBinaryPri1 | Circumflex expressionBinaryPri1 | Concat expressionBinaryPri1 | 'AT' Identifier Identifier expressionBinaryPri1 )*
    ;

// TSql160.g:31377
expressionBinaryPri1
    : expressionUnary ( Star expressionUnary | Divide expressionUnary | PercentSign expressionUnary | LeftShift expressionUnary | RightShift expressionUnary )*
    ;

// TSql160.g:31406
expressionUnary
    : ( Plus | Minus | Tilde ) expressionUnary
    | expressionWithClrElements
    ;

// TSql160.g:31441
expressionWithClrElements
    : ( udtExpression | expressionPrimary ) ( Dot identifier ( expressionWithClrElementsFunctionCallPart collationOpt | collationOpt ) )*
    ;

// TSql160.g:31484
expressionWithClrElementsFunctionCallPart
    : parenthesizedOptExpressionWithDefaultList
    ;

// TSql160.g:31489
// TSql160.g:31494: predicate dropped: {LA(1) == LeftCurly && NextTokenMatches(CodeGenerationSupporter.Fn, 2)}?
// TSql160.g:31531: predicate dropped: {((LA(1) == Identifier || LA(1) == QuotedIdentifier) && LA(2) == Dot && LA(3) == DollarPartition) || LA(1) == DollarPartition}?
// TSql160.g:18685: alternative dropped, its action refuses every input: tOdbcInitiator : OdbcInitiator {ThrowParseErrorException("SQL46036", tOdbcInitiator, TSqlParserResource.SQL46036Message);}
// TSql160.g:18684: rule dropped, every alternative refuses its input
// TSql160.g:31499: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Cast) && (LA(2) == LeftParenthesis)}?
// TSql160.g:31502: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.TryCast) && (LA(2) == LeftParenthesis)}?
// TSql160.g:31505: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Parse) && (LA(2) == LeftParenthesis)}?
// TSql160.g:31508: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.TryParse) && (LA(2) == LeftParenthesis)}?
// TSql160.g:31511: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.IIf) && (LA(2) == LeftParenthesis)}?
// TSql160.g:31514: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.JsonObject) && (LA(2) == LeftParenthesis)}?
// TSql160.g:31517: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.JsonArray) && (LA(2) == LeftParenthesis)}?
// TSql160.g:31520: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Trim) && (LA(2) == LeftParenthesis) && (LA(3) != Star)}?
// TSql160.g:31534: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Next) && NextTokenMatches(CodeGenerationSupporter.Value, 2)}?
expressionPrimary
    : ( odbcFunctionCall | literal | castCall_CAST | tryCastCall_TRY_CAST | parseCall_PARSE | tryParseCall_TRY_PARSE | iIfCall_IIF | jsonObjectCall_JSON_OBJECT | jsonArrayCall_JSON_ARRAY | trimCall_TRIM | builtInFunctionCall | leftFunctionCall | rightFunctionCall | partitionFunctionCall | nextValueForCall_NEXT | columnOrFunctionCall | nullIfExpression | coalesceExpression | caseExpression | convertCall | tryConvertCall | parameterlessCall | parenthesisDisambiguatorForExpressions ) collationOpt
    ;

// TSql160.g:31556
// TSql160.g:31559: predicate dropped: {IsNextRuleSelectParenthesis()}?
parenthesisDisambiguatorForExpressions
    : subquery
    | expressionParenthesis
    ;

// TSql160.g:31565
udtExpression
    : userDefinedTypeCallTarget ( udtFunctionExpression | udtPropertyExpression )
    ;

// TSql160.g:31579
userDefinedTypeCallTarget
    : schemaObjectTwoPartName DoubleColon
    ;

// TSql160.g:31590
udtPropertyExpression
    : identifier collationOpt
    ;

// TSql160.g:31601
udtFunctionExpression
    : identifier parenthesizedOptExpressionWithDefaultList
    ;

// TSql160.g:31624
expressionParenthesis
    : LeftParenthesis expressionWithFlags RightParenthesis
    ;

// TSql160.g:31636
convertCall
    : Convert LeftParenthesis scalarDataType Comma expression ( Comma expression )? RightParenthesis
    ;

// TSql160.g:31660
tryConvertCall
    : TryConvert LeftParenthesis scalarDataType Comma expression ( Comma expression )? RightParenthesis
    ;

// TSql160.g:31776
parameterlessCall
    : User
    | CurrentUser
    | SessionUser
    | SystemUser
    | CurrentTimestamp
    | CurrentDate
    ;

// TSql160.g:31815
expressionList
    : expression ( Comma expression )*
    ;

// TSql160.g:31831
jsonObjectExpressionList
    : ( jsonKeyValueExpression ( Comma jsonKeyValueExpression )* |  ) ( jsonNullClauseFunction |  )
    ;

// TSql160.g:31857
jsonNullClauseFunction
    : ( Null On Null | 'ABSENT' On Null )
    ;

// TSql160.g:31882
jsonKeyValueExpression
    : ( expression Colon expression | Label expression )
    ;

// TSql160.g:31916
windowClause
    : 'WINDOW' windowDefinition ( Comma windowDefinition )*
    ;

// TSql160.g:31938
windowDefinition
    : identifier As LeftParenthesis ( identifier )? ( 'PARTITION' By expressionList )? ( orderByClause )? ( windowFrameClause )? RightParenthesis
    ;

// TSql160.g:31983
// TSql160.g:32003: action may refuse the input: {// Window name inside parentheses requires at least one clause (PARTITION BY, ORDER BY, or window frame) if (vResult.WindowName != null && vResult.Partitions.Count == 0 && vResult.OrderByClause == null && vResult.WindowFrameClause == null) { ThrowPar...
overClause
    : overClauseBeginning ( orderByClause )? ( windowFrameClause )? RightParenthesis
    | overClauseWithWindow
    ;

// TSql160.g:32014
overClauseWithWindow
    : Over identifier
    ;

// TSql160.g:32029
overClauseNoOrderBy
    : overClauseBeginningNoWindowName RightParenthesis
    ;

// TSql160.g:32038
overClauseBeginning
    : Over LeftParenthesis ( identifier )? ( 'PARTITION' By expressionList )?
    ;

// TSql160.g:32063
overClauseBeginningNoWindowName
    : Over LeftParenthesis ( 'PARTITION' By expressionList |  )
    ;

// TSql160.g:32082
windowFrameClause
    : ( 'ROWS' | 'RANGE' ) windowFrameExtent
    ;

// TSql160.g:32103
windowFrameExtent
    : Between windowFrameDelimiter And windowFrameDelimiter
    | windowFrameDelimiter
    ;

// TSql160.g:32121
windowFrameDelimiter
    : Current Identifier
    | integer ( 'PRECEDING' | 'FOLLOWING' )
    | 'UNBOUNDED' ( 'PRECEDING' | 'FOLLOWING' )
    ;

// TSql160.g:32170
// TSql160.g:32184: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Graph) && (LT(2).getText().Equals(CodeGenerationSupporter.Path, StringComparison.OrdinalIgnoreCase))}?
withinGroupClause
    : 'WITHIN' Group LeftParenthesis ( ( graphPathClause_GRAPH ) | orderByClause ) RightParenthesis
    ;

// TSql160.g:32202
// TSql160.g:32212: predicate dropped: {(vResult.FunctionName != null && vResult.FunctionName.Value.ToUpper(CultureInfo.InvariantCulture) == CodeGenerationSupporter.JsonArray)}?
// TSql160.g:32215: predicate dropped: {(vResult.FunctionName != null && vResult.FunctionName.Value.ToUpper(CultureInfo.InvariantCulture) == CodeGenerationSupporter.JsonObject)}?
// TSql160.g:32218: predicate dropped: {(vResult.FunctionName != null && vResult.FunctionName.Value.ToUpper(CultureInfo.InvariantCulture) == CodeGenerationSupporter.Trim) && (NextTokenMatches(CodeGenerationSupporter.Leading) | NextTokenMatches(CodeGenerationSupporter.Trailing) | NextTokenMatches(Cod...
// TSql160.g:32222: predicate dropped: {(vResult.FunctionName != null && vResult.FunctionName.Value.ToUpper(CultureInfo.InvariantCulture) == CodeGenerationSupporter.Trim)}?
// TSql160.g:32225: predicate dropped: {(vResult.FunctionName != null && (vResult.FunctionName.Value.ToUpper(CultureInfo.InvariantCulture) == CodeGenerationSupporter.LastValue || vResult.FunctionName.Value.ToUpper(CultureInfo.InvariantCulture) == CodeGenerationSupporter.FirstValue || vResult.Functio...
builtInFunctionCall
    : nonQuotedIdentifier LeftParenthesis ( jsonArrayBuiltInFunctionCall | jsonObjectBuiltInFunctionCall | trim3ArgsBuiltInFunctionCall | trimBuiltInFunctionCall | ignoreRespectNullsRegularBuiltInFunctionCall | regularBuiltInFunctionCall | aggregateBuiltInFunctionCall )
    ;

// TSql160.g:32237
jsonArrayBuiltInFunctionCall
    : ( expressionList |  ) ( jsonNullClauseFunction |  ) RightParenthesis
    ;

// TSql160.g:32256
jsonObjectBuiltInFunctionCall
    : ( jsonObjectExpressionList |  ) RightParenthesis
    ;

// TSql160.g:32272
// TSql160.g:32295: predicate dropped: {(LA(1)==Over && (LA(2)==LeftParenthesis || LA(2)==Identifier || LA(2)==QuotedIdentifier)) || (NextTokenMatches(CodeGenerationSupporter.Within) && LA(2)==Group && LA(3)==LeftParenthesis)}?
regularBuiltInFunctionCall
    : ( expressionList | starColumnReferenceExpression |  ) RightParenthesis ( withinGroupOrOverClause )?
    ;

// TSql160.g:32301
// TSql160.g:32321: predicate dropped: {((NextTokenMatches(CodeGenerationSupporter.Ignore) || NextTokenMatches(CodeGenerationSupporter.Respect) )&& LA(2)==Identifier)}?
// TSql160.g:32326: predicate dropped: {(LA(1)==Over && (LA(2)==LeftParenthesis || LA(2)==Identifier || LA(2)==QuotedIdentifier)) || (NextTokenMatches(CodeGenerationSupporter.Within) && LA(2)==Group && LA(3)==LeftParenthesis)}?
ignoreRespectNullsRegularBuiltInFunctionCall
    : ( expressionList | starColumnReferenceExpression |  ) RightParenthesis ( ignoreRespectNulls )? ( withinGroupOrOverClause )?
    ;

// TSql160.g:32332
ignoreRespectNulls
    : ( 'IGNORE' | 'RESPECT' ) 'NULLS'
    ;

// TSql160.g:32354
withinGroupOrOverClause
    : withinGroupClause ( overClauseNoOrderBy )?
    | overClause
    ;

// TSql160.g:32377
starColumnReferenceExpression
    : Star
    ;

// TSql160.g:32385
aggregateBuiltInFunctionCall
    : uniqueRowFilter expression RightParenthesis ( overClauseNoOrderBy )?
    ;

// TSql160.g:32407
trimBuiltInFunctionCall
    : expression ( From expression )? RightParenthesis
    ;

// TSql160.g:32427
trim3ArgsBuiltInFunctionCall
    : Identifier expression From expression RightParenthesis
    ;

// TSql160.g:32452
leftFunctionCall
    : Left reservedBuiltInFunctionCallParameters
    ;

// TSql160.g:32461
rightFunctionCall
    : Right reservedBuiltInFunctionCallParameters
    ;

// TSql160.g:32470
reservedBuiltInFunctionCallParameters
    : LeftParenthesis ( expressionList )? RightParenthesis
    ;

// TSql160.g:32480
expressionWithDefault
    : expression
    | defaultLiteral
    ;

// TSql160.g:32485
expressionWithDefaultList
    : expressionWithDefault ( Comma expressionWithDefault )*
    ;

// TSql160.g:32501
parenthesizedOptExpressionWithDefaultList
    : LeftParenthesis ( expressionWithDefaultList )? RightParenthesis
    ;

// TSql160.g:32511
nullIfExpression
    : NullIf LeftParenthesis expressionWithFlags Comma expressionWithFlags RightParenthesis
    ;

// TSql160.g:32613
coalesceExpression
    : Coalesce LeftParenthesis expressionWithFlags ( Comma expressionWithFlags )+ RightParenthesis
    ;

// TSql160.g:32640
simpleWhenClause
    : When expressionWithFlags Then expressionWithFlags
    ;

// TSql160.g:32659
searchedWhenClause
    : When booleanExpressionWithFlags Then expressionWithFlags
    ;

// TSql160.g:32679
caseExpression
    : Case ( expression simpleCaseExpression | searchedCaseExpression ) ( Else expression )? End
    ;

// TSql160.g:32704
simpleCaseExpression
    : ( simpleWhenClause )+
    ;

// TSql160.g:32718
searchedCaseExpression
    : ( searchedWhenClause )+
    ;

// TSql160.g:32731
systemColumn
    : PseudoColumn
    ;

// TSql160.g:32740
specialColumn
    : IdentityColumn
    | RowGuidColumn
    | systemColumn
    ;

// TSql160.g:32758
column
    : ( multiPartIdentifier ( Dot specialColumn )? | specialColumn )
    ;

// TSql160.g:32783
partitionFunctionCall
    : ( identifier Dot )? DollarPartition Dot identifier LeftParenthesis expressionList RightParenthesis
    ;

// TSql160.g:32831
userFunctionCall
    : LeftParenthesis ( ( expressionWithDefaultList )? | uniqueRowFilter expressionList ) RightParenthesis ( overClauseNoOrderBy )?
    ;

// TSql160.g:32863
columnOrFunctionCall
    : ( multiPartIdentifier ( Dot specialColumn | userFunctionCall |  ) | specialColumn )
    ;

// TSql160.g:32917
starColumn
    : ( multiPartIdentifier Dot Star | Star )
    ;

// TSql160.g:32949
selectStarExpression
    : ( multiPartIdentifier Dot Star | Star )
    ;

// TSql160.g:32977
revertStatement
    : Revert ( With 'COOKIE' EqualsSign expression )?
    ;

// TSql160.g:32997
nullNotNull
    : ( Not )? Null
    ;

// TSql160.g:33010
orderByOption
    : Asc
    | Desc
    ;

// TSql160.g:33023
signedIntegerOrVariable
    : signedInteger
    | variable
    ;

// TSql160.g:33028
signedIntegerOrStringOrVariable
    : signedInteger
    | stringOrVariable
    ;

// TSql160.g:33033
signedIntegerOrVariableOrNull
    : signedIntegerOrVariable
    | nullLiteral
    ;

// TSql160.g:33038
stringLiteralOrNull
    : stringLiteral
    | nullLiteral
    ;

// TSql160.g:33043
signedInteger
    : ( Minus )? integer
    ;

// TSql160.g:33070
signedIntegerOrReal
    : ( Minus )? integerOrRealOrNumeric
    ;

// TSql160.g:33097
identifierOrVariable
    : variable
    | identifier
    ;

// TSql160.g:33112
stringOrVariable
    : stringLiteral
    | variable
    ;

// TSql160.g:33117
stringOrGlobalVariableOrVariable
    : stringLiteral
    | globalVariableOrVariableReference
    ;

// TSql160.g:33137
objectOrString
    : schemaObjectThreePartName
    | stringLiteral
    ;

// TSql160.g:33152
authorizationOpt
    : ( authorization )?
    ;

// TSql160.g:33156
authorization
    : Authorization identifier
    ;

// TSql160.g:33166
collationOpt
    : /* empty */
    | collation
    ;

// TSql160.g:33171
collation
    : Collate nonQuotedIdentifier
    ;

// TSql160.g:33181
assemblyName
    : 'NAME' identifier ( Dot identifier )?
    ;

// TSql160.g:33206
identifierVariable
    : Variable
    ;

// TSql160.g:33214
identifierOrInteger
    : identifier
    | integer
    ;

// TSql160.g:33229
stringOrIdentifier
    : stringLiteral
    | identifier
    ;

// TSql160.g:33245
stringOrSignedIntegerOrIdentifier
    : stringLiteral
    | signedInteger
    | identifier
    ;

// TSql160.g:33265
nonEmptyStringOrIdentifier
    : nonEmptyString
    | identifier
    ;

// TSql160.g:33280
nonQuotedIdentifier
    : Identifier
    ;

// TSql160.g:33289
sqlCommandIdentifier
    : SqlCommandIdentifier
    ;

// TSql160.g:33298
identifier
    : Identifier
    | QuotedIdentifier
    ;

// TSql160.g:33315
seedIncrement
    : ( Minus | Plus )? integerOrNumeric
    ;

// TSql160.g:33348
subroutineParameterLiteral
    : integer
    | numeric
    | real
    | moneyLiteral
    ;

// TSql160.g:33355
max
    : 'MAX'
    ;

// TSql160.g:33365
binaryOrIntegerOrStringOrVariable
    : binary
    | stringLiteral
    | integer
    | variable
    ;

// TSql160.g:33376
integerOrNumeric
    : integer
    | numeric
    ;

// TSql160.g:33382
literal
    : integer
    | numeric
    | real
    | moneyLiteral
    | binary
    | stringLiteral
    | nullLiteral
    | globalVariableOrVariableReference
    | odbcLiteral
    ;

// TSql160.g:33394
globalVariableOrVariableReference
    : Variable
    ;

// TSql160.g:33417
variable
    : Variable
    ;

// TSql160.g:33425
integer
    : Integer
    ;

// TSql160.g:33433
numeric
    : Numeric
    ;

// TSql160.g:33441
real
    : Real
    ;

// TSql160.g:33449
// TSql160.g:33451: action may refuse the input: {if (vResult.Value == null || vResult.Value.Length == 0) { ThrowParseErrorException("SQL46063", vResult, TSqlParserResource.SQL46063Message); }}
nonEmptyString
    : stringLiteral
    ;

// TSql160.g:33460
defaultValueLiteral
    : literal
    | signedIntegerOrReal
    ;

// TSql160.g:33464
stringLiteral
    : AsciiStringLiteral
    | UnicodeStringLiteral
    ;

// TSql160.g:33480
stringLiteralOrVariable
    : stringLiteral
    | variable
    ;

// TSql160.g:33495
binary
    : HexLiteral
    ;

// TSql160.g:33505
nullLiteral
    : Null
    ;

// TSql160.g:33513
moneyLiteral
    : Money
    ;

// TSql160.g:33521
odbcLiteral
    : LeftCurly Identifier ( AsciiStringLiteral | UnicodeStringLiteral ) RightCurly
    ;

// TSql160.g:33549
defaultLiteral
    : Default
    ;

// TSql160.g:33557
onOffOptionValue
    : optionOnOff
    ;

// TSql160.g:33568
stringOrDefaultLiteralOptionValue
    : ( stringLiteral | defaultLiteral )
    ;

// TSql160.g:33586
stringLiteralOptionValue
    : stringLiteral
    ;

// TSql160.g:33597
integerOrDefaultLiteralOptionValue
    : ( integer | defaultLiteral )
    ;

// TSql160.g:33615
integerLiteralOptionValue
    : integer
    ;

// TSql160.g:33626
defaultLiteralOptionValue
    : defaultLiteral
    ;

// TSql160.g:33637
binaryOrDefaultLiteralOptionValue
    : ( binary | defaultLiteral )
    ;

// TSql160.g:33655
assignmentWithOptOp
    : EqualsSign
    | assignmentWithOp
    ;

// TSql160.g:33665
assignmentWithOp
    : AddEquals
    | SubtractEquals
    | MultiplyEquals
    | DivideEquals
    | ModEquals
    | BitwiseAndEquals
    | BitwiseOrEquals
    | BitwiseXorEquals
    | ConcatEquals
    ;

// TSql160.g:33704
onOff
    : On
    | Off
    ;

// TSql160.g:33716
securityStatementPermission
    : Identifier
    | Add
    | All
    | Alter
    | And
    | Any
    | As
    | Asc
    | Authorization
    | Backup
    | Begin
    | Between
    | Break
    | Browse
    | Bulk
    | By
    | Cascade
    | Case
    | Check
    | Checkpoint
    | Close
    | Clustered
    | Coalesce
    | Collate
    | Column
    | Commit
    | Compute
    | Connection
    | Constraint
    | Contains
    | ContainsTable
    | Continue
    | Convert
    | Create
    | Cross
    | Current
    | CurrentDate
    | CurrentTime
    | CurrentTimestamp
    | CurrentUser
    | Cursor
    | Database
    | Dbcc
    | Deallocate
    | Declare
    | Default
    | Delete
    | Deny
    | Desc
    | Disk
    | Distinct
    | Distributed
    | Double
    | Drop
    | Else
    | End
    | Errlvl
    | Escape
    | Except
    | Exec
    | Execute
    | External
    | Exists
    | Exit
    | Fetch
    | File
    | FillFactor
    | For
    | Foreign
    | FreeText
    | FreeTextTable
    | Full
    | Function
    | GoTo
    | Group
    | Having
    | HoldLock
    | Identity
    | IdentityInsert
    | IdentityColumn
    | If
    | In
    | Index
    | Inner
    | Insert
    | Intersect
    | Into
    | Is
    | Join
    | Key
    | Kill
    | Left
    | Like
    | LineNo
    | Merge
    | National
    | NoCheck
    | NonClustered
    | Not
    | Null
    | NullIf
    | Object
    | Of
    | Off
    | Offsets
    | Open
    | OpenDataSource
    | OpenQuery
    | OpenRowSet
    | OpenXml
    | Option
    | Or
    | Order
    | Outer
    | Over
    | Percent
    | Pivot
    | Plan
    | Precision
    | Primary
    | Print
    | Proc
    | Procedure
    | Public
    | Raiserror
    | Read
    | ReadText
    | Reconfigure
    | References
    | Replication
    | Restore
    | Restrict
    | Return
    | Revert
    | Revoke
    | Right
    | Rollback
    | RowCount
    | RowGuidColumn
    | Rule
    | Save
    | Schema
    | Select
    | SessionUser
    | Set
    | SetUser
    | Shutdown
    | Some
    | Statistics
    | StopList
    | SystemUser
    | Table
    | TableSample
    | TextSize
    | Then
    | Top
    | Tran
    | Transaction
    | Trigger
    | Truncate
    | TSEqual
    | Union
    | Unique
    | Unpivot
    | Update
    | UpdateText
    | Use
    | User
    | Values
    | Varying
    | View
    | WaitFor
    | When
    | While
    | With
    | WriteText
    ;

// TSql160.g:24634
// createMaterializedViewStatement where a predicate requires the next word to be MATERIALIZED
createMaterializedViewStatement_MATERIALIZED
    : 'MATERIALIZED' View schemaObjectTwoPartName With LeftParenthesis createMaterializedViewOption ( Comma createMaterializedViewOption )? RightParenthesis As subqueryExpressionWithOptionalCTE
    ;

// TSql160.g:10003
// createFederationStatement where a predicate requires the next word to be FEDERATION
createFederationStatement_FEDERATION
    : 'FEDERATION' identifier LeftParenthesis identifier scalarDataType 'RANGE' RightParenthesis
    ;

// TSql160.g:24674
// alterMaterializedViewStatement where a predicate requires the next word to be MATERIALIZED
alterMaterializedViewStatement_MATERIALIZED
    : 'MATERIALIZED' View schemaObjectTwoPartName ( 'REBUILD' | 'DISABLE' )
    ;

// TSql160.g:12166
// alterFederationStatement where a predicate requires the next word to be FEDERATION
alterFederationStatement_FEDERATION
    : 'FEDERATION' identifier ( 'SPLIT' 'AT' LeftParenthesis | Drop 'AT' LeftParenthesis ( 'LOW' | 'HIGH' ) ) identifier EqualsSign expression RightParenthesis
    ;

// TSql160.g:809
// disableTriggerStatement where a predicate requires the next word to be DISABLE
disableTriggerStatement_DISABLE
    : 'DISABLE' enableDisableTriggerBody
    ;

// TSql160.g:820
// enableTriggerStatement where a predicate requires the next word to be ENABLE
enableTriggerStatement_ENABLE
    : 'ENABLE' enableDisableTriggerBody
    ;

// TSql160.g:25507
// moveConversationStatement where a predicate requires the next word to be MOVE
moveConversationStatement_MOVE
    : 'MOVE' 'CONVERSATION' expression To expression
    ;

// TSql160.g:25521
// getConversationGroupStatement where a predicate requires the next word to be GET
getConversationGroupStatement_GET
    : 'GET' 'CONVERSATION' Group variable From schemaObjectThreePartName
    ;

// TSql160.g:25536
// receiveStatement where a predicate requires the next word to be RECEIVE
receiveStatement_RECEIVE
    : 'RECEIVE' ( Top LeftParenthesis expression RightParenthesis )? receiveSelectExpression ( Comma receiveSelectExpression )* From schemaObjectThreePartName ( Into variableDmlTarget )? ( Where ( 'CONVERSATION_GROUP_ID' | 'CONVERSATION_HANDLE' ) EqualsSign stringOrVariable )?
    ;

// TSql160.g:25611
// sendStatement where a predicate requires the next word to be SEND
sendStatement_SEND
    : 'SEND' On 'CONVERSATION' ( LeftParenthesis expression ( Comma expression )* RightParenthesis | expression ) ( 'MESSAGE' 'TYPE' identifierOrVariable )? ( LeftParenthesis expression RightParenthesis )?
    ;

// TSql160.g:19537
// throwStatement where a predicate requires the next word to be THROW
throwStatement_THROW
    : 'THROW' ( integerOrVariable Comma stringOrVariable Comma integerOrVariable )?
    ;

// TSql160.g:476
// renameEntityStatement where a predicate requires the next word to be RENAME
renameEntityStatement_RENAME
    : 'RENAME' securityStatementPermission ( DoubleColon )? schemaObjectThreePartName To ( identifier )
    ;

// TSql160.g:991
// createAggregateStatement where a predicate requires the next word to be AGGREGATE
createAggregateStatement_AGGREGATE
    : 'AGGREGATE' schemaObjectThreePartName LeftParenthesis aggregateParameter ( Comma aggregateParameter )* RightParenthesis 'RETURNS' scalarDataType External assemblyName
    ;

// TSql160.g:11121
// applicationRoleStatement where a predicate requires the next word to be APPLICATION
applicationRoleStatement_APPLICATION
    : 'APPLICATION' 'ROLE' identifier With applicationRoleOption ( Comma applicationRoleOption )*
    ;

// TSql160.g:1045
// createApplicationRoleStatement where a predicate requires the next word to be APPLICATION
createApplicationRoleStatement_APPLICATION
    : applicationRoleStatement_APPLICATION
    ;

// TSql160.g:1050
// createAssemblyStatement where a predicate requires the next word to be ASSEMBLY
createAssemblyStatement_ASSEMBLY
    : 'ASSEMBLY' identifier authorizationOpt From expressionList ( With Identifier EqualsSign assemblyPermissionSetOption )?
    ;

// TSql160.g:1324
// createAsymmetricKeyStatement where a predicate requires the next word to be ASYMMETRIC
createAsymmetricKeyStatement_ASYMMETRIC
    : 'ASYMMETRIC' Key identifier authorizationOpt createAsymmetricKeyParams ( encryptClause )?
    ;

// TSql160.g:9558
// createAvailabilityGroupStatement where a predicate requires the next word to be AVAILABILITY
createAvailabilityGroupStatement_AVAILABILITY
    : 'AVAILABILITY' Group identifier ( With LeftParenthesis availabilityGroupOption RightParenthesis )? For ( Database identifier ( Comma identifier )* )? 'REPLICA' On availabilityReplica ( Comma availabilityReplica )*
    ;

// TSql160.g:22834
// brokerPriorityStatementBody where a predicate requires the next word to be BROKER
brokerPriorityStatementBody_BROKER
    : 'BROKER' 'PRIORITY' identifier For 'CONVERSATION' ( Set LeftParenthesis brokerPriorityParameter ( Comma brokerPriorityParameter )* RightParenthesis )?
    ;

// TSql160.g:22807
// createBrokerPriorityStatement where a predicate requires the next word to be BROKER
createBrokerPriorityStatement_BROKER
    : brokerPriorityStatementBody_BROKER
    ;

// TSql160.g:1414
// createCertificateStatement where a predicate requires the next word to be CERTIFICATE
createCertificateStatement_CERTIFICATE
    : 'CERTIFICATE' identifier authorizationOpt createCertificateParams ( createCertificateActivityFlag )?
    ;

// TSql160.g:16918
// createColumnStoreIndexStatement where a predicate requires the next word to be COLUMNSTORE
createColumnStoreIndexStatement_COLUMNSTORE
    : 'COLUMNSTORE' Index identifier On schemaObjectFourPartName ( LeftParenthesis columnstoreIndexKeyColumnList RightParenthesis )? ( Order ( identifierColumnList ) )? ( filterClause )? ( With ( indexOptionList ) )? ( On filegroupOrPartitionScheme )?
    ;

// TSql160.g:1559
// createContractStatement where a predicate requires the next word to be CONTRACT
createContractStatement_CONTRACT
    : 'CONTRACT' identifier authorizationOpt LeftParenthesis contractMessage ( Comma contractMessage )* RightParenthesis
    ;

// TSql160.g:1641
// credentialStatementBody where a predicate requires the next word to be CREDENTIAL
credentialStatementBody_CREDENTIAL
    : 'CREDENTIAL' identifier With Identity EqualsSign stringLiteral ( Comma 'SECRET' EqualsSign stringLiteral )?
    ;

// TSql160.g:1625
// createCredentialStatement where a predicate requires the next word to be CREDENTIAL
createCredentialStatement_CREDENTIAL
    : credentialStatementBody_CREDENTIAL ( For 'CRYPTOGRAPHIC' 'PROVIDER' identifier )?
    ;

// TSql160.g:21980
// createCryptographicProviderStatement where a predicate requires the next word to be CRYPTOGRAPHIC
createCryptographicProviderStatement_CRYPTOGRAPHIC
    : 'CRYPTOGRAPHIC' 'PROVIDER' identifier cryptographicProviderFile
    ;

// TSql160.g:6520
// createEndpointStatement where a predicate requires the next word to be ENDPOINT
createEndpointStatement_ENDPOINT
    : 'ENDPOINT' identifier authorizationOpt ( endpointOptions )? protocolInfo payloadInfo
    ;

// TSql160.g:7572
// createEventStatement where a predicate requires the next word to be EVENT
// TSql160.g:7578: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Notification)}?
// TSql160.g:7581: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Session)}?
createEventStatement_EVENT
    : 'EVENT' ( createEventNotificationStatement_NOTIFICATION | createEventSessionStatement_SESSION )
    ;

// TSql160.g:8311
// createFulltextStatement where a predicate requires the next word to be FULLTEXT
createFulltextStatement_FULLTEXT
    : 'FULLTEXT' ( createFulltextCatalogStatement | createFulltextIndexStatement | createFulltextStoplistStatement )
    ;

// TSql160.g:8962
// createSelectiveXmlIndexStatement where a predicate requires the next word to be SELECTIVE
createSelectiveXmlIndexStatement_SELECTIVE
    : 'SELECTIVE' 'XML' Index identifier On schemaObjectThreePartName LeftParenthesis identifier RightParenthesis ( With xmlNamespaces )? For LeftParenthesis promotedSelectiveXmlIndexPath ( Comma promotedSelectiveXmlIndexPath )* RightParenthesis ( With indexOptionList )?
    ;

// TSql160.g:9023
// createXmlStatements where a predicate requires the next word to be XML
createXmlStatements_XML
    : 'XML' ( createXmlIndexStatement | createXmlSchemaCollectionStatement )
    ;

// TSql160.g:9240
// createLoginStatement where a predicate requires the next word to be LOGIN
// TSql160.g:9256: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Windows)}?
// TSql160.g:9258: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.External)}?
createLoginStatement_LOGIN
    : 'LOGIN' identifier ( passwordLoginSource | From ( windowsLoginSource_WINDOWS | externalLoginSource | certificateLoginSource | asymmetricKeyLoginSource ) )
    ;

// TSql160.g:9511
// createMessageTypeStatement where a predicate requires the next word to be MESSAGE
createMessageTypeStatement_MESSAGE
    : 'MESSAGE' 'TYPE' identifier authorizationOpt ( messageTypeValidation )?
    ;

// TSql160.g:7546
// createMasterKeyStatement where a predicate requires the next word to be MASTER
createMasterKeyStatement_MASTER
    : 'MASTER' Key ( 'ENCRYPTION' By 'PASSWORD' EqualsSign stringLiteral )?
    ;

// TSql160.g:10033
// createPartitionStatement where a predicate requires the next word to be PARTITION
createPartitionStatement_PARTITION
    : 'PARTITION' ( createPartitionFunction | createPartitionScheme )
    ;

// TSql160.g:10189
// createQueueStatement where a predicate requires the next word to be QUEUE
createQueueStatement_QUEUE
    : 'QUEUE' schemaObjectThreePartName ( queueOptionList )? ( On stringOrIdentifier )?
    ;

// TSql160.g:10214
// createRemoteServiceBindingStatement where a predicate requires the next word to be REMOTE
createRemoteServiceBindingStatement_REMOTE
    : 'REMOTE' 'SERVICE' 'BINDING' identifier authorizationOpt To 'SERVICE' stringLiteral With ( ( bindingUserOption ( Comma bindingAnonymousOption )? ) | ( bindingAnonymousOption ( Comma bindingUserOption )? ) )
    ;

// TSql160.g:22112
// createResourcePoolStatement where a predicate requires the next word to be RESOURCE
createResourcePoolStatement_RESOURCE
    : 'RESOURCE' 'POOL' identifier resourcePoolStatementBody
    ;

// TSql160.g:10287
// createRoleStatement where a predicate requires the next word to be ROLE
createRoleStatement_ROLE
    : 'ROLE' identifier authorizationOpt
    ;

// TSql160.g:10313
// createRouteStatement where a predicate requires the next word to be ROUTE
createRouteStatement_ROUTE
    : 'ROUTE' identifier authorizationOpt routeOptionList
    ;

// TSql160.g:8804
// createSearchPropertyListStatement where a predicate requires the next word to be SEARCH
createSearchPropertyListStatement_SEARCH
    : 'SEARCH' 'PROPERTY' 'LIST' identifier ( From multiPartIdentifier )? authorizationOpt requiredSemicolon
    ;

// TSql160.g:10327
// createServiceStatement where a predicate requires the next word to be SERVICE
createServiceStatement_SERVICE
    : 'SERVICE' identifier authorizationOpt onQueueClause ( LeftParenthesis serviceContract ( Comma serviceContract )* RightParenthesis )?
    ;

// TSql160.g:24267
// createSpatialIndexStatement where a predicate requires the next word to be SPATIAL
createSpatialIndexStatement_SPATIAL
    : 'SPATIAL' Index identifier On schemaObjectFourPartName LeftParenthesis identifier RightParenthesis ( spatialIndexingScheme )? ( With spatialIndexOptionsList )? ( On stringOrIdentifier )?
    ;

// TSql160.g:10438
// createSymmetricKeyStatement where a predicate requires the next word to be SYMMETRIC
createSymmetricKeyStatement_SYMMETRIC
    : 'SYMMETRIC' Key identifier authorizationOpt ( ( With keySpecList encryptorListWithNoPassword ) | ( From 'PROVIDER' identifier providerKeySourceOptionsOpt ( encryptorListWithNoPassword )? ) )
    ;

// TSql160.g:10706
// createSynonymStatement where a predicate requires the next word to be SYNONYM
createSynonymStatement_SYNONYM
    : 'SYNONYM' schemaObjectThreePartName For schemaObjectFourPartName
    ;

// TSql160.g:10727
// createTypeStatement where a predicate requires the next word to be TYPE
createTypeStatement_TYPE
    : 'TYPE' schemaObjectThreePartName ( createTypeUddtStatement | createTypeUdtStatement | createTypeTableStatement )
    ;

// TSql160.g:1664
// createServerStatements where a predicate requires the next word to be SERVER
// TSql160.g:1670: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Audit)}?
// TSql160.g:1673: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Role)}?
createServerStatements_SERVER
    : 'SERVER' ( createServerAuditStatements_AUDIT | createServerRoleStatement_ROLE )
    ;

// TSql160.g:22467
// workloadGroupStatementBody where a predicate requires the next word to be WORKLOAD
workloadGroupStatementBody_WORKLOAD
    : 'WORKLOAD' Group identifier ( With LeftParenthesis workloadGroupParameter ( Comma workloadGroupParameter )* RightParenthesis )? ( 'USING' workloadGroupUsingOption ( Comma workloadGroupUsingOption )* )?
    ;

// TSql160.g:22456
// createWorkloadGroupStatement where a predicate requires the next word to be WORKLOAD
createWorkloadGroupStatement_WORKLOAD
    : workloadGroupStatementBody_WORKLOAD
    ;

// TSql160.g:22634
// workloadClassifierStatementBody where a predicate requires the next word to be WORKLOAD
workloadClassifierStatementBody_WORKLOAD
    : 'WORKLOAD' 'CLASSIFIER' identifier withClassifierOptions
    ;

// TSql160.g:22629
// createWorkloadClassifierStatement where a predicate requires the next word to be WORKLOAD
createWorkloadClassifierStatement_WORKLOAD
    : workloadClassifierStatementBody_WORKLOAD
    ;

// TSql160.g:22447
// createWorkloadStatements where a predicate requires the next word to be WORKLOAD
createWorkloadStatements_WORKLOAD
    : ( createWorkloadGroupStatement_WORKLOAD | createWorkloadClassifierStatement_WORKLOAD )
    ;

// TSql160.g:22948
// createSequenceStatement where a predicate requires the next word to be SEQUENCE
createSequenceStatement_SEQUENCE
    : 'SEQUENCE' schemaObjectTwoPartName ( createSequenceOptionList )?
    ;

// TSql160.g:23395
// createSecurityPolicyStatement where a predicate requires the next word to be SECURITY
createSecurityPolicyStatement_SECURITY
    : 'SECURITY' 'POLICY' schemaObjectThreePartName ( createSecurityPredicateListElement ( Comma createSecurityPredicateListElement )* )? ( With LeftParenthesis securityPolicyOptionsList RightParenthesis )? ( Not For Replication )?
    ;

// TSql160.g:1366
// assemblyEncryptionSource where a predicate requires the next word to be ASSEMBLY
assemblyEncryptionSource_ASSEMBLY
    : 'ASSEMBLY' identifier
    ;

// TSql160.g:1678
// createServerAuditStatements where a predicate requires the next word to be AUDIT
createServerAuditStatements_AUDIT
    : 'AUDIT' ( createServerAuditSpecificationStatement | createServerAuditStatement )
    ;

// TSql160.g:10300
// createServerRoleStatement where a predicate requires the next word to be ROLE
createServerRoleStatement_ROLE
    : 'ROLE' identifier authorizationOpt
    ;

// TSql160.g:1823
// maxSizeAuditFileOption where a predicate requires the next word to be MAXSIZE
maxSizeAuditFileOption_MAXSIZE
    : 'MAXSIZE' EqualsSign ( integer ( 'GB' | 'TB' | 'MB' ) | 'UNLIMITED' )
    ;

// TSql160.g:1886
// maxRolloverFilesAuditFileOption where a predicate requires the next word to be MAX_ROLLOVER_FILES
maxRolloverFilesAuditFileOption_MAX_ROLLOVER_FILES
    : 'MAX_ROLLOVER_FILES' EqualsSign ( integer | 'UNLIMITED' )
    ;

// TSql160.g:1924
// reserveDiskSpaceAuditFileOption where a predicate requires the next word to be RESERVE_DISK_SPACE
reserveDiskSpaceAuditFileOption_RESERVE_DISK_SPACE
    : 'RESERVE_DISK_SPACE' EqualsSign optionOnOff
    ;

// TSql160.g:1911
// maxFilesAuditFileOption where a predicate requires the next word to be MAX_FILES
maxFilesAuditFileOption_MAX_FILES
    : 'MAX_FILES' EqualsSign integer
    ;

// TSql160.g:1868
// retentionDaysAuditUrlOption where a predicate requires the next word to be RETENTION_DAYS
retentionDaysAuditUrlOption_RETENTION_DAYS
    : 'RETENTION_DAYS' EqualsSign ( integer )
    ;

// TSql160.g:2154
// alterServerAuditStatements where a predicate requires the next word to be AUDIT
// TSql160.g:2160: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Specification)}?
alterServerAuditStatements_AUDIT
    : 'AUDIT' ( alterServerAuditSpecificationStatement_SPECIFICATION | alterServerAuditStatement )
    ;

// TSql160.g:2247
// alterServerConfigurationStatement where a predicate requires the next word to be CONFIGURATION
// TSql160.g:2254: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Process)}?
// TSql160.g:2257: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Buffer)}?
// TSql160.g:2260: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Diagnostics)}?
// TSql160.g:2263: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Failover)}?
// TSql160.g:2266: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Hadr)}?
// TSql160.g:2269: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.SoftNuma)}?
// TSql160.g:2272: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.External)}?
alterServerConfigurationStatement_CONFIGURATION
    : 'CONFIGURATION' Set ( alterServerConfigurationSetProcessAffinityStatement_PROCESS | alterServerConfigurationSetBufferPoolExtensionStatement_BUFFER | alterServerConfigurationSetDiagnosticsLogStatement_DIAGNOSTICS | alterServerConfigurationSetFailoverClusterPropertyStatement_FAILOVER | alterServerConfigurationSetHadrClusterStatement_HADR | alterServerConfigurationSetSoftNumaStatement_SOFTNUMA | alterServerConfigurationSetExternalAuthenticationStatement_EXTERNAL )
    ;

// TSql160.g:12536
// alterServerRoleStatement where a predicate requires the next word to be ROLE
alterServerRoleStatement_ROLE
    : 'ROLE' identifier alterRoleAction
    ;

// TSql160.g:2167
// alterServerAuditSpecificationStatement where a predicate requires the next word to be SPECIFICATION
alterServerAuditSpecificationStatement_SPECIFICATION
    : 'SPECIFICATION' identifier ( auditSpecificationForClause )? ( auditSpecificationDetail ( Comma auditSpecificationDetail )* )? auditSpecificationStateOpt
    ;

// TSql160.g:2617
// alterServerConfigurationSetProcessAffinityStatement where a predicate requires the next word to be PROCESS
alterServerConfigurationSetProcessAffinityStatement_PROCESS
    : 'PROCESS' 'AFFINITY' Identifier EqualsSign ( affinityRangeList | 'AUTO' )
    ;

// TSql160.g:2380
// alterServerConfigurationSetBufferPoolExtensionStatement where a predicate requires the next word to be BUFFER
alterServerConfigurationSetBufferPoolExtensionStatement_BUFFER
    : 'BUFFER' 'POOL' 'EXTENSION' alterServerConfigurationBufferPoolExtensionContainerOption
    ;

// TSql160.g:2462
// alterServerConfigurationSetDiagnosticsLogStatement where a predicate requires the next word to be DIAGNOSTICS
alterServerConfigurationSetDiagnosticsLogStatement_DIAGNOSTICS
    : 'DIAGNOSTICS' 'LOG' alterServerConfigurationDiagnosticsLogOption
    ;

// TSql160.g:2535
// alterServerConfigurationSetFailoverClusterPropertyStatement where a predicate requires the next word to be FAILOVER
alterServerConfigurationSetFailoverClusterPropertyStatement_FAILOVER
    : 'FAILOVER' 'CLUSTER' 'PROPERTY' alterServerConfigurationFailoverClusterPropertyOption
    ;

// TSql160.g:2579
// alterServerConfigurationSetHadrClusterStatement where a predicate requires the next word to be HADR
alterServerConfigurationSetHadrClusterStatement_HADR
    : 'HADR' 'CLUSTER' alterServerConfigurationHadrClusterOption
    ;

// TSql160.g:2354
// alterServerConfigurationSetSoftNumaStatement where a predicate requires the next word to be SOFTNUMA
alterServerConfigurationSetSoftNumaStatement_SOFTNUMA
    : 'SOFTNUMA' alterServerConfigurationSoftNumaOption
    ;

// TSql160.g:2278
// alterServerConfigurationSetExternalAuthenticationStatement where a predicate requires the next word to be EXTERNAL
alterServerConfigurationSetExternalAuthenticationStatement_EXTERNAL
    : 'EXTERNAL' 'AUTHENTICATION' alterServerConfigurationExternalAuthenticationContainerOption
    ;

// TSql160.g:2344
// alterServerConfigurationExternalAuthenticationUseIdentityOption where a predicate requires the next word to be USE_IDENTITY
alterServerConfigurationExternalAuthenticationUseIdentityOption_USE_IDENTITY
    : 'USE_IDENTITY'
    ;

// TSql160.g:2332
// alterServerConfigurationExternalAuthenticationCredentialNameOption where a predicate requires the next word to be CREDENTIAL_NAME
alterServerConfigurationExternalAuthenticationCredentialNameOption_CREDENTIAL_NAME
    : 'CREDENTIAL_NAME' EqualsSign stringLiteralOptionValue
    ;

// TSql160.g:2509
// alterServerConfigurationDiagnosticsLogMaxSizeOption where a predicate requires the next word to be MAX_SIZE
alterServerConfigurationDiagnosticsLogMaxSizeOption_MAX_SIZE
    : 'MAX_SIZE' EqualsSign ( integerLiteralOptionValue 'MB' | defaultLiteralOptionValue )
    ;

// TSql160.g:4206
// alterDatabaseAuditSpecification where a predicate requires the next word to be AUDIT
alterDatabaseAuditSpecification_AUDIT
    : 'AUDIT' 'SPECIFICATION' identifier ( auditSpecificationForClause )? ( auditSpecificationDetailDb ( Comma auditSpecificationDetailDb )* )? auditSpecificationStateOpt
    ;

// TSql160.g:11466
// alterDatabaseScopedCredentialStatement where a predicate requires the next word to be SCOPED
alterDatabaseScopedCredentialStatement_SCOPED
    : 'SCOPED' credentialStatementBody
    ;

// TSql160.g:2703
// alterDatabaseScopedConfigurationStatement where a predicate requires the next word to be SCOPED
alterDatabaseScopedConfigurationStatement_SCOPED
    : 'SCOPED' 'CONFIGURATION' ( For Identifier )? ( alterDatabaseScopedConfigSet | alterDatabaseScopedConfigClear )
    ;

// TSql160.g:2788
// alterDatabaseScopedMaxDopOption where a predicate requires the next word to be MAXDOP
// TSql160.g:2804: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.Primary)}?
alterDatabaseScopedMaxDopOption_MAXDOP
    : 'MAXDOP' EqualsSign ( integer | Primary )
    ;

// TSql160.g:2817
// alterDatabaseScopedOnOffPrimaryOption where a predicate requires the next word to be QUERY_OPTIMIZER_HOTFIXES
alterDatabaseScopedOnOffPrimaryOption_QUERY_OPTIMIZER_HOTFIXES
    : 'QUERY_OPTIMIZER_HOTFIXES' EqualsSign ( databaseConfigurationOptionOnOffPrimary )
    ;

// TSql160.g:2817
// alterDatabaseScopedOnOffPrimaryOption where a predicate requires the next word to be PARAMETER_SNIFFING
alterDatabaseScopedOnOffPrimaryOption_PARAMETER_SNIFFING
    : 'PARAMETER_SNIFFING' EqualsSign ( databaseConfigurationOptionOnOffPrimary )
    ;

// TSql160.g:2817
// alterDatabaseScopedOnOffPrimaryOption where a predicate requires the next word to be LEGACY_CARDINALITY_ESTIMATION
alterDatabaseScopedOnOffPrimaryOption_LEGACY_CARDINALITY_ESTIMATION
    : 'LEGACY_CARDINALITY_ESTIMATION' EqualsSign ( databaseConfigurationOptionOnOffPrimary )
    ;

// TSql160.g:3006
// alterDbRemove where a predicate requires the next word to be REMOVE
alterDbRemove_REMOVE
    : 'REMOVE' ( File identifier | 'FILEGROUP' identifier )
    ;

// TSql160.g:3032
// alterDbModify where a predicate requires the next word to be MODIFY
// TSql160.g:3041: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Name)}?
alterDbModify_MODIFY
    : 'MODIFY' ( ( 'NAME' EqualsSign identifier ) | ( 'FILEGROUP' alterDbModifyFilegroup ) | alterDbModifyFile | alterDbModifyAzureOptions )
    ;

// TSql160.g:3257
// cursorDefaultDbOption where a predicate requires the next word to be CURSOR_DEFAULT
cursorDefaultDbOption_CURSOR_DEFAULT
    : 'CURSOR_DEFAULT' ( 'LOCAL' | 'GLOBAL' )
    ;

// TSql160.g:3273
// recoveryDbOption where a predicate requires the next word to be RECOVERY
recoveryDbOption_RECOVERY
    : 'RECOVERY' ( Full | Identifier )
    ;

// TSql160.g:3294
// targetRecoveryTimeDbOption where a predicate requires the next word to be TARGET_RECOVERY_TIME
targetRecoveryTimeDbOption_TARGET_RECOVERY_TIME
    : 'TARGET_RECOVERY_TIME' EqualsSign integer Identifier
    ;

// TSql160.g:3308
// pageVerifyDbOption where a predicate requires the next word to be PAGE_VERIFY
pageVerifyDbOption_PAGE_VERIFY
    : 'PAGE_VERIFY' Identifier
    ;

// TSql160.g:3378
// partnerOption where a predicate requires the next word to be PARTNER
// TSql160.g:3388: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Timeout)}?
// TSql160.g:3396: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Safety)}?
partnerOption_PARTNER
    : 'PARTNER' ( ( ( 'TIMEOUT' integer ) | ( 'SAFETY' ( Full | Off ) ) | Identifier ) | Off | EqualsSign stringLiteral )
    ;

// TSql160.g:3434
// witnessOption where a predicate requires the next word to be WITNESS
witnessOption_WITNESS
    : 'WITNESS' ( Off | EqualsSign stringLiteral )
    ;

// TSql160.g:3457
// parameterizationOption where a predicate requires the next word to be PARAMETERIZATION
parameterizationOption_PARAMETERIZATION
    : 'PARAMETERIZATION' ( 'SIMPLE' | 'FORCED' )
    ;

// TSql160.g:3473
// compatibilityLevelDbOption where a predicate requires the next word to be COMPATIBILITY_LEVEL
compatibilityLevelDbOption_COMPATIBILITY_LEVEL
    : 'COMPATIBILITY_LEVEL' EqualsSign integer
    ;

// TSql160.g:3486
// changeTrackingDbOption where a predicate requires the next word to be CHANGE_TRACKING
changeTrackingDbOption_CHANGE_TRACKING
    : 'CHANGE_TRACKING' ( ( EqualsSign Off ) | ( EqualsSign On ( changeTrackingOnOptions )? ) | changeTrackingOnOptions )
    ;

// TSql160.g:3517
// acceleratedDatabaseRecoveryOption where a predicate requires the next word to be ACCELERATED_DATABASE_RECOVERY
acceleratedDatabaseRecoveryOption_ACCELERATED_DATABASE_RECOVERY
    : 'ACCELERATED_DATABASE_RECOVERY' ( ( EqualsSign Off ) | ( EqualsSign On ) )
    ;

// TSql160.g:4128
// dbContainmentOption where a predicate requires the next word to be CONTAINMENT
dbContainmentOption_CONTAINMENT
    : 'CONTAINMENT' EqualsSign Identifier
    ;

// TSql160.g:4142
// hadrDbOption where a predicate requires the next word to be HADR
// TSql160.g:4145: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Availability)}?
hadrDbOption_HADR
    : 'HADR' ( hadrAvailabilityDbOption_AVAILABILITY | simpleHadrDbOption )
    ;

// TSql160.g:4192
// dbDelayedDurabilityOption where a predicate requires the next word to be DELAYED_DURABILITY
dbDelayedDurabilityOption_DELAYED_DURABILITY
    : 'DELAYED_DURABILITY' EqualsSign Identifier
    ;

// TSql160.g:3318
// autoCreateStatisticsDbOption where a predicate requires the next word to be AUTO_CREATE_STATISTICS
autoCreateStatisticsDbOption_AUTO_CREATE_STATISTICS
    : 'AUTO_CREATE_STATISTICS' optionOnOff ( LeftParenthesis 'INCREMENTAL' EqualsSign optionOnOff RightParenthesis )?
    ;

// TSql160.g:5171
// remoteDataArchiveDbOption where a predicate requires the next word to be REMOTE_DATA_ARCHIVE
remoteDataArchiveDbOption_REMOTE_DATA_ARCHIVE
    : 'REMOTE_DATA_ARCHIVE' ( ( EqualsSign Off ) | ( EqualsSign On ( remoteDataArchiveDbSettingsList ) ) | ( remoteDataArchiveDbSettingsList ) )
    ;

// TSql160.g:3599
// queryStoreDbOption where a predicate requires the next word to be QUERY_STORE
queryStoreDbOption_QUERY_STORE
    : 'QUERY_STORE' ( 'CLEAR' | 'CLEAR' All | ( EqualsSign Off ) | ( EqualsSign On ( queryStoreOnOptions )? ) | queryStoreOnOptions )
    ;

// TSql160.g:3903
// automaticTuningDbOption where a predicate requires the next word to be AUTOMATIC_TUNING
// TSql160.g:3913: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Auto)}?
// TSql160.g:3923: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Inherit)}?
// TSql160.g:3933: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Custom)}?
automaticTuningDbOption_AUTOMATIC_TUNING
    : 'AUTOMATIC_TUNING' ( EqualsSign ( ( 'AUTO' ) | ( 'INHERIT' ) | ( 'CUSTOM' ) ) | automaticTuningOptions )
    ;

// TSql160.g:3699
// queryStoreDesiredStateOption where a predicate requires the next word to be DESIRED_STATE
queryStoreDesiredStateOption_DESIRED_STATE
    : ( ( 'DESIRED_STATE' ) ) EqualsSign Identifier
    ;

// TSql160.g:3699
// queryStoreDesiredStateOption where a predicate requires the next word to be OPERATION_MODE
queryStoreDesiredStateOption_OPERATION_MODE
    : ( ( 'OPERATION_MODE' ) ) EqualsSign Identifier
    ;

// TSql160.g:3730
// queryStoreCaptureModeOption where a predicate requires the next word to be QUERY_CAPTURE_MODE
queryStoreCaptureModeOption_QUERY_CAPTURE_MODE
    : 'QUERY_CAPTURE_MODE' EqualsSign ( ( All ) | ( Identifier ) )
    ;

// TSql160.g:3756
// queryStoreSizeCleanupModeOption where a predicate requires the next word to be SIZE_BASED_CLEANUP_MODE
queryStoreSizeCleanupModeOption_SIZE_BASED_CLEANUP_MODE
    : 'SIZE_BASED_CLEANUP_MODE' EqualsSign ( ( Off ) | ( Identifier ) )
    ;

// TSql160.g:3782
// queryStoreDataFlushIntervalOption where a predicate requires the next word to be FLUSH_INTERVAL_SECONDS
queryStoreDataFlushIntervalOption_FLUSH_INTERVAL_SECONDS
    : ( 'FLUSH_INTERVAL_SECONDS' ) EqualsSign integer
    ;

// TSql160.g:3782
// queryStoreDataFlushIntervalOption where a predicate requires the next word to be DATA_FLUSH_INTERVAL_SECONDS
queryStoreDataFlushIntervalOption_DATA_FLUSH_INTERVAL_SECONDS
    : ( 'DATA_FLUSH_INTERVAL_SECONDS' ) EqualsSign integer
    ;

// TSql160.g:3809
// queryStoreIntervalLengthOption where a predicate requires the next word to be INTERVAL_LENGTH_MINUTES
queryStoreIntervalLengthOption_INTERVAL_LENGTH_MINUTES
    : 'INTERVAL_LENGTH_MINUTES' EqualsSign integer
    ;

// TSql160.g:3825
// queryStoreMaxQdsSizeOption where a predicate requires the next word to be MAX_STORAGE_SIZE_MB
queryStoreMaxQdsSizeOption_MAX_STORAGE_SIZE_MB
    : 'MAX_STORAGE_SIZE_MB' EqualsSign integer
    ;

// TSql160.g:3841
// queryStoreMaxPlansPerQueryOption where a predicate requires the next word to be MAX_PLANS_PER_QUERY
queryStoreMaxPlansPerQueryOption_MAX_PLANS_PER_QUERY
    : 'MAX_PLANS_PER_QUERY' EqualsSign integer
    ;

// TSql160.g:3857
// queryStoreTimeCleanupPolicy where a predicate requires the next word to be CLEANUP_POLICY
queryStoreTimeCleanupPolicy_CLEANUP_POLICY
    : 'CLEANUP_POLICY' EqualsSign LeftParenthesis 'STALE_QUERY_THRESHOLD_DAYS' EqualsSign integer RightParenthesis
    ;

// TSql160.g:3879
// queryStoreWaitStatsCaptureOption where a predicate requires the next word to be WAIT_STATS_CAPTURE_MODE
queryStoreWaitStatsCaptureOption_WAIT_STATS_CAPTURE_MODE
    : 'WAIT_STATS_CAPTURE_MODE' ( ( EqualsSign Off ) | ( EqualsSign On ) )
    ;

// TSql160.g:4058
// automaticTuningForceLastGoodPlanOption where a predicate requires the next word to be FORCE_LAST_GOOD_PLAN
automaticTuningForceLastGoodPlanOption_FORCE_LAST_GOOD_PLAN
    : 'FORCE_LAST_GOOD_PLAN' EqualsSign ( On (  ) | Off (  ) | Default (  ) )
    ;

// TSql160.g:3988
// automaticTuningCreateIndexOption where a predicate requires the next word to be CREATE_INDEX
automaticTuningCreateIndexOption_CREATE_INDEX
    : 'CREATE_INDEX' EqualsSign ( On (  ) | Off (  ) | Default (  ) )
    ;

// TSql160.g:4023
// automaticTuningDropIndexOption where a predicate requires the next word to be DROP_INDEX
automaticTuningDropIndexOption_DROP_INDEX
    : 'DROP_INDEX' EqualsSign ( On (  ) | Off (  ) | Default (  ) )
    ;

// TSql160.g:4093
// automaticTuningMaintainIndexOption where a predicate requires the next word to be MAINTAIN_INDEX
automaticTuningMaintainIndexOption_MAINTAIN_INDEX
    : 'MAINTAIN_INDEX' EqualsSign ( On (  ) | Off (  ) | Default (  ) )
    ;

// TSql160.g:4157
// hadrAvailabilityDbOption where a predicate requires the next word to be AVAILABILITY
hadrAvailabilityDbOption_AVAILABILITY
    : 'AVAILABILITY' Group EqualsSign identifier
    ;

// TSql160.g:5267
// createDatabaseAuditSpecification where a predicate requires the next word to be AUDIT
createDatabaseAuditSpecification_AUDIT
    : 'AUDIT' 'SPECIFICATION' identifier auditSpecificationForClause ( createAuditSpecificationDetailDb ( Comma createAuditSpecificationDetailDb )* )? auditSpecificationStateOpt
    ;

// TSql160.g:1614
// createDatabaseScopedCredentialStatement where a predicate requires the next word to be SCOPED
createDatabaseScopedCredentialStatement_SCOPED
    : 'SCOPED' credentialStatementBody
    ;

// TSql160.g:4506
// azureMaxSizeDatabaseOption where a predicate requires the next word to be MAXSIZE
azureMaxSizeDatabaseOption_MAXSIZE
    : 'MAXSIZE' EqualsSign integer ( 'GB' | 'MB' )
    ;

// TSql160.g:4530
// azureEditionDatabaseOption where a predicate requires the next word to be EDITION
azureEditionDatabaseOption_EDITION
    : 'EDITION' EqualsSign stringLiteral
    ;

// TSql160.g:4727
// nameFileOption where a predicate requires the next word to be NAME
nameFileOption_NAME
    : 'NAME' EqualsSign nonEmptyStringOrIdentifier
    ;

// TSql160.g:4755
// fileNameFileOption where a predicate requires the next word to be FILENAME
fileNameFileOption_FILENAME
    : 'FILENAME' EqualsSign nonEmptyString
    ;

// TSql160.g:4768
// sizeFileOption where a predicate requires the next word to be SIZE
sizeFileOption_SIZE
    : 'SIZE' EqualsSign integer ( memUnit )?
    ;

// TSql160.g:4787
// maxSizeFileOption where a predicate requires the next word to be MAXSIZE
maxSizeFileOption_MAXSIZE
    : 'MAXSIZE' EqualsSign ( ( integer ( memUnit )? ) | 'UNLIMITED' )
    ;

// TSql160.g:4819
// fileGrowthFileOption where a predicate requires the next word to be FILEGROWTH
fileGrowthFileOption_FILEGROWTH
    : 'FILEGROWTH' EqualsSign ( integer ( ( memUnit ) | PercentSign )? )
    ;

// TSql160.g:4741
// newNameFileOption where a predicate requires the next word to be NEWNAME
newNameFileOption_NEWNAME
    : 'NEWNAME' EqualsSign nonEmptyStringOrIdentifier
    ;

// TSql160.g:4852
// offlineFileOption where a predicate requires the next word to be OFFLINE
offlineFileOption_OFFLINE
    : 'OFFLINE'
    ;

// TSql160.g:5097
// restrictedUserCreateDbOption where a predicate requires the next word to be RESTRICTED_USER
restrictedUserCreateDbOption_RESTRICTED_USER
    : 'RESTRICTED_USER'
    ;

// TSql160.g:4951
// createDbOnOffOption where a predicate requires the next word to be DB_CHAINING
createDbOnOffOption_DB_CHAINING
    : ( 'DB_CHAINING' ) optionOnOff
    ;

// TSql160.g:4951
// createDbOnOffOption where a predicate requires the next word to be TRUSTWORTHY
createDbOnOffOption_TRUSTWORTHY
    : ( 'TRUSTWORTHY' ) optionOnOff
    ;

// TSql160.g:4951
// createDbOnOffOption where a predicate requires the next word to be DATA_RETENTION
createDbOnOffOption_DATA_RETENTION
    : ( 'DATA_RETENTION' ) optionOnOff
    ;

// TSql160.g:5106
// fileStreamCreateAlterDbOption where a predicate requires the next word to be FILESTREAM
fileStreamCreateAlterDbOption_FILESTREAM
    : 'FILESTREAM' LeftParenthesis ( fileStreamOption ) ( Comma fileStreamOption )* RightParenthesis
    ;

// TSql160.g:33298
// identifier where a predicate requires the next word to be SERVER
identifier_SERVER
    : 'SERVER'
    ;

// TSql160.g:33298
// identifier where a predicate requires the next word to be CREDENTIAL
identifier_CREDENTIAL
    : 'CREDENTIAL'
    ;

// TSql160.g:33298
// identifier where a predicate requires the next word to be FEDERATED_SERVICE_ACCOUNT
identifier_FEDERATED_SERVICE_ACCOUNT
    : 'FEDERATED_SERVICE_ACCOUNT'
    ;

// TSql160.g:6005
// backupCertificateStatement where a predicate requires the next word to be CERTIFICATE
backupCertificateStatement_CERTIFICATE
    : 'CERTIFICATE' identifier To File EqualsSign stringLiteral ( With privateKeySpec )?
    ;

// TSql160.g:6022
// backupServiceMasterKeyStatement where a predicate requires the next word to be SERVICE
backupServiceMasterKeyStatement_SERVICE
    : 'SERVICE' backupMasterKeyBody
    ;

// TSql160.g:6030
// restoreServiceMasterKeyStatement where a predicate requires the next word to be SERVICE
restoreServiceMasterKeyStatement_SERVICE
    : 'SERVICE' restoreMasterKeyBody ( 'FORCE' )?
    ;

// TSql160.g:5849
// backupEncryptionOption where a predicate requires the next word to be ENCRYPTION
backupEncryptionOption_ENCRYPTION
    : 'ENCRYPTION' LeftParenthesis 'ALGORITHM' EqualsSign Identifier Comma backupEncrytor RightParenthesis
    ;

// TSql160.g:5986
// fileStreamRestoreOption where a predicate requires the next word to be FILESTREAM
fileStreamRestoreOption_FILESTREAM
    : fileStreamCreateAlterDbOption_FILESTREAM
    ;

// TSql160.g:6544
// endpointState where a predicate requires the next word to be STATE
endpointState_STATE
    : 'STATE' EqualsSign Identifier
    ;

// TSql160.g:6801
// listenerIpProtocolOption where a predicate requires the next word to be LISTENER_IP
listenerIpProtocolOption_LISTENER_IP
    : 'LISTENER_IP' EqualsSign ( ( LeftParenthesis ( ipAddressV4 ( Colon ipAddressV4 )? | stringLiteral ) RightParenthesis ) | All )
    ;

// TSql160.g:6688
// listenerPortProtocolOption where a predicate requires the next word to be LISTENER_PORT
listenerPortProtocolOption_LISTENER_PORT
    : 'LISTENER_PORT' EqualsSign integer
    ;

// TSql160.g:6640
// pathProtocolOption where a predicate requires the next word to be PATH
pathProtocolOption_PATH
    : 'PATH' EqualsSign stringLiteral
    ;

// TSql160.g:6652
// siteProtocolOption where a predicate requires the next word to be SITE
siteProtocolOption_SITE
    : 'SITE' EqualsSign stringLiteral
    ;

// TSql160.g:6664
// clearPortProtocolOption where a predicate requires the next word to be CLEAR_PORT
clearPortProtocolOption_CLEAR_PORT
    : 'CLEAR_PORT' EqualsSign integer
    ;

// TSql160.g:6676
// sslPortProtocolOption where a predicate requires the next word to be SSL_PORT
sslPortProtocolOption_SSL_PORT
    : 'SSL_PORT' EqualsSign integer
    ;

// TSql160.g:6701
// authenticationProtocolOption where a predicate requires the next word to be AUTHENTICATION
authenticationProtocolOption_AUTHENTICATION
    : 'AUTHENTICATION' EqualsSign LeftParenthesis Identifier ( Comma Identifier )* RightParenthesis
    ;

// TSql160.g:6723
// portsProtocolOption where a predicate requires the next word to be PORTS
portsProtocolOption_PORTS
    : 'PORTS' EqualsSign LeftParenthesis Identifier ( Comma Identifier )* RightParenthesis
    ;

// TSql160.g:6745
// authenticationRealmProtocolOption where a predicate requires the next word to be AUTH_REALM
authenticationRealmProtocolOption_AUTH_REALM
    : 'AUTH_REALM' EqualsSign ( nonEmptyString | 'NONE' )
    ;

// TSql160.g:6765
// defaultLogonDomainProtocolOption where a predicate requires the next word to be DEFAULT_LOGON_DOMAIN
defaultLogonDomainProtocolOption_DEFAULT_LOGON_DOMAIN
    : 'DEFAULT_LOGON_DOMAIN' EqualsSign ( nonEmptyString | 'NONE' )
    ;

// TSql160.g:7143
// batchesPayloadOption where a predicate requires the next word to be BATCHES
batchesPayloadOption_BATCHES
    : 'BATCHES' EqualsSign enabledDisabled
    ;

// TSql160.g:7164
// wsdlPayloadOption where a predicate requires the next word to be WSDL
wsdlPayloadOption_WSDL
    : 'WSDL' EqualsSign ( 'NONE' | defaultLiteral | stringLiteral )
    ;

// TSql160.g:7189
// sessionsPayloadOption where a predicate requires the next word to be SESSIONS
sessionsPayloadOption_SESSIONS
    : 'SESSIONS' EqualsSign enabledDisabled
    ;

// TSql160.g:7197
// loginTypePayloadOption where a predicate requires the next word to be LOGIN_TYPE
loginTypePayloadOption_LOGIN_TYPE
    : 'LOGIN_TYPE' EqualsSign ( 'WINDOWS' | 'MIXED' )
    ;

// TSql160.g:7212
// sessionTimeoutPayloadOption where a predicate requires the next word to be SESSION_TIMEOUT
sessionTimeoutPayloadOption_SESSION_TIMEOUT
    : 'SESSION_TIMEOUT' EqualsSign ( 'NEVER' | integer )
    ;

// TSql160.g:7248
// namespacePayloadOption where a predicate requires the next word to be NAMESPACE
namespacePayloadOption_NAMESPACE
    : 'NAMESPACE' EqualsSign ( stringLiteral | defaultLiteral )
    ;

// TSql160.g:7278
// characterSetPayloadOption where a predicate requires the next word to be CHARACTER_SET
characterSetPayloadOption_CHARACTER_SET
    : 'CHARACTER_SET' EqualsSign ( 'SQL' | 'XML' )
    ;

// TSql160.g:7293
// headerLimitPayloadOption where a predicate requires the next word to be HEADER_LIMIT
headerLimitPayloadOption_HEADER_LIMIT
    : 'HEADER_LIMIT' EqualsSign integer
    ;

// TSql160.g:7305
// authenticationPayloadOption where a predicate requires the next word to be AUTHENTICATION
authenticationPayloadOption_AUTHENTICATION
    : 'AUTHENTICATION' EqualsSign Identifier (  | identifier | ( identifier identifier (  | identifier ) ) )
    ;

// TSql160.g:7377
// encryptionPayloadOption where a predicate requires the next word to be ENCRYPTION
encryptionPayloadOption_ENCRYPTION
    : 'ENCRYPTION' EqualsSign Identifier ( identifier ( identifier | ( identifier identifier ) ) )?
    ;

// TSql160.g:7413
// messageForwardingPayloadOption where a predicate requires the next word to be MESSAGE_FORWARDING
messageForwardingPayloadOption_MESSAGE_FORWARDING
    : 'MESSAGE_FORWARDING' EqualsSign enabledDisabled
    ;

// TSql160.g:7421
// messageForwardSizePayloadOption where a predicate requires the next word to be MESSAGE_FORWARD_SIZE
messageForwardSizePayloadOption_MESSAGE_FORWARD_SIZE
    : 'MESSAGE_FORWARD_SIZE' EqualsSign integer
    ;

// TSql160.g:7586
// createEventNotificationStatement where a predicate requires the next word to be NOTIFICATION
createEventNotificationStatement_NOTIFICATION
    : 'NOTIFICATION' identifier eventNotificationObjectScope ( With 'FAN_IN' )? For eventTypeGroupContainer ( Comma eventTypeGroupContainer )* To 'SERVICE' stringLiteral Comma stringLiteral
    ;

// TSql160.g:7706
// createEventSessionStatement where a predicate requires the next word to be SESSION
createEventSessionStatement_SESSION
    : 'SESSION' identifier eventSessionScope eventDeclarationList ( targetDeclarationList )? optSessionOptionList
    ;

// TSql160.g:7993
// maxMemorySessionOption where a predicate requires the next word to be MAX_MEMORY
maxMemorySessionOption_MAX_MEMORY
    : 'MAX_MEMORY' integerSessionOptionValue
    ;

// TSql160.g:8012
// maxEventSizeSessionOption where a predicate requires the next word to be MAX_EVENT_SIZE
maxEventSizeSessionOption_MAX_EVENT_SIZE
    : 'MAX_EVENT_SIZE' integerSessionOptionValue
    ;

// TSql160.g:8045
// trackCausalitySessionOption where a predicate requires the next word to be TRACK_CAUSALITY
trackCausalitySessionOption_TRACK_CAUSALITY
    : 'TRACK_CAUSALITY' onOffSessionOption
    ;

// TSql160.g:8053
// startupStateSessionOption where a predicate requires the next word to be STARTUP_STATE
startupStateSessionOption_STARTUP_STATE
    : 'STARTUP_STATE' onOffSessionOption
    ;

// TSql160.g:8071
// eventRetentionSessionOption where a predicate requires the next word to be EVENT_RETENTION_MODE
eventRetentionSessionOption_EVENT_RETENTION_MODE
    : 'EVENT_RETENTION_MODE' EqualsSign Identifier
    ;

// TSql160.g:8080
// memoryPartitionSessionOption where a predicate requires the next word to be MEMORY_PARTITION_MODE
memoryPartitionSessionOption_MEMORY_PARTITION_MODE
    : 'MEMORY_PARTITION_MODE' EqualsSign Identifier
    ;

// TSql160.g:8019
// maxDispatchLatencySessionOption where a predicate requires the next word to be MAX_DISPATCH_LATENCY
maxDispatchLatencySessionOption_MAX_DISPATCH_LATENCY
    : 'MAX_DISPATCH_LATENCY' EqualsSign ( 'INFINITE' | integer 'SECONDS' )
    ;

// TSql160.g:8698
// languageTerm where a predicate requires the next word to be LANGUAGE
languageTerm_LANGUAGE
    : 'LANGUAGE' ( stringLiteral | integer | identifier | binary )
    ;

// TSql160.g:8544
// changeTrackingFulltextIndexOption where a predicate requires the next word to be CHANGE_TRACKING
changeTrackingFulltextIndexOption_CHANGE_TRACKING
    : 'CHANGE_TRACKING' ( EqualsSign )? ( Off ( Comma 'NO' 'POPULATION' )? | ( 'MANUAL' | 'AUTO' ) )
    ;

// TSql160.g:23663
// createExternalDataSourceStatement where a predicate requires the next word to be DATA
// TSql160.g:23680: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Type)}?
// TSql160.g:23683: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Location)}?
// TSql160.g:23686: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.PushdownOption)}?
// TSql160.g:23699: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Type)}?
// TSql160.g:23702: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Location)}?
// TSql160.g:23705: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.PushdownOption)}?
createExternalDataSourceStatement_DATA
    : 'DATA' 'SOURCE' identifier With LeftParenthesis ( externalDataSourceType_TYPE | externalDataSourceLocation_LOCATION | externalDataSourcePushdownOption_PUSHDOWN | externalDataSourceLiteralOrIdentifierOption ) ( Comma ( externalDataSourceType_TYPE | externalDataSourceLocation_LOCATION | externalDataSourcePushdownOption_PUSHDOWN | externalDataSourceLiteralOrIdentifierOption ) )* RightParenthesis
    ;

// TSql160.g:22285
// createExternalResourcePoolStatement where a predicate requires the next word to be RESOURCE
createExternalResourcePoolStatement_RESOURCE
    : 'RESOURCE' 'POOL' identifier externalResourcePoolStatementBody
    ;

// TSql160.g:1074
// createExternalLibraryStatement where a predicate requires the next word to be LIBRARY
createExternalLibraryStatement_LIBRARY
    : 'LIBRARY' identifier authorizationOpt From createExternalLibraryFileOption ( Comma createExternalLibraryFileOption )* With LeftParenthesis 'LANGUAGE' EqualsSign nonEmptyString RightParenthesis
    ;

// TSql160.g:1187
// createExternalLanguageStatement where a predicate requires the next word to be LANGUAGE
createExternalLanguageStatement_LANGUAGE
    : 'LANGUAGE' identifier authorizationOpt From externalLanguageFileOption ( Comma externalLanguageFileOption )*
    ;

// TSql160.g:23889
// createExternalStreamStatement where a predicate requires the next word to be STREAM
// TSql160.g:23904: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Location)}?
// TSql160.g:23907: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.InputOptions)}?
// TSql160.g:23910: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.OutputOptions)}?
// TSql160.g:23923: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Location)}?
// TSql160.g:23926: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.InputOptions)}?
// TSql160.g:23929: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.OutputOptions)}?
createExternalStreamStatement_STREAM
    : 'STREAM' identifier With LeftParenthesis ( externalStreamLocation_LOCATION | externalStreamInputOptions_INPUT_OPTIONS | externalStreamOutputOptions_OUTPUT_OPTIONS | externalStreamLiteralOrIdentifierOption ) ( Comma ( externalStreamLocation_LOCATION | externalStreamInputOptions_INPUT_OPTIONS | externalStreamOutputOptions_OUTPUT_OPTIONS | externalStreamLiteralOrIdentifierOption ) )* RightParenthesis
    ;

// TSql160.g:23830
// alterExternalDataSourceStatement where a predicate requires the next word to be DATA
// TSql160.g:23846: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Location)}?
// TSql160.g:23849: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.PushdownOption)}?
// TSql160.g:23862: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Location)}?
// TSql160.g:23865: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.PushdownOption)}?
alterExternalDataSourceStatement_DATA
    : 'DATA' 'SOURCE' identifier Set ( externalDataSourceLocation_LOCATION | externalDataSourcePushdownOption_PUSHDOWN | externalDataSourceLiteralOrIdentifierOption ) ( Comma ( externalDataSourceLocation_LOCATION | externalDataSourcePushdownOption_PUSHDOWN | externalDataSourceLiteralOrIdentifierOption ) )*
    ;

// TSql160.g:22300
// alterExternalResourcePoolStatement where a predicate requires the next word to be RESOURCE
alterExternalResourcePoolStatement_RESOURCE
    : 'RESOURCE' 'POOL' identifier externalResourcePoolStatementBody
    ;

// TSql160.g:1104
// alterExternalLibraryStatement where a predicate requires the next word to be LIBRARY
alterExternalLibraryStatement_LIBRARY
    : 'LIBRARY' identifier authorizationOpt Set alterExternalLibraryFileOption With LeftParenthesis 'LANGUAGE' EqualsSign nonEmptyString RightParenthesis
    ;

// TSql160.g:1263
// alterExternalLanguageStatement where a predicate requires the next word to be LANGUAGE
alterExternalLanguageStatement_LANGUAGE
    : 'LANGUAGE' identifier ( authorizationOpt ( ( Set externalLanguageFileOption ) | ( Add externalLanguageFileOption ) | ( 'REMOVE' Identifier identifier ) ) )
    ;

// TSql160.g:33298
// identifier where a predicate requires the next word to be XQUERY
identifier_XQUERY
    : 'XQUERY'
    ;

// TSql160.g:33298
// identifier where a predicate requires the next word to be SINGLETON
identifier_SINGLETON
    : 'SINGLETON'
    ;

// TSql160.g:33298
// identifier where a predicate requires the next word to be MAXLENGTH
identifier_MAXLENGTH
    : 'MAXLENGTH'
    ;

// TSql160.g:9356
// windowsLoginSource where a predicate requires the next word to be WINDOWS
windowsLoginSource_WINDOWS
    : 'WINDOWS' ( With createLoginParamWin ( Comma createLoginParamWin )* )?
    ;

// TSql160.g:9684
// alterAvailabilityGroupModifyReplica where a predicate requires the next word to be MODIFY
alterAvailabilityGroupModifyReplica_MODIFY
    : 'MODIFY' 'REPLICA' On availabilityReplica ( Comma availabilityReplica )*
    ;

// TSql160.g:9703
// alterAvailabilityGroupRemoveReplica where a predicate requires the next word to be REMOVE
alterAvailabilityGroupRemoveReplica_REMOVE
    : 'REMOVE' 'REPLICA' On availabilityReplicaName ( Comma availabilityReplicaName )*
    ;

// TSql160.g:9750
// alterAvailabilityGroupFailoverAction where a predicate requires the next word to be FAILOVER
alterAvailabilityGroupFailoverAction_FAILOVER
    : 'FAILOVER' ( With LeftParenthesis alterAvailabilityGroupFailoverOption RightParenthesis )?
    ;

// TSql160.g:9865
// availabilityModeReplicaOption where a predicate requires the next word to be AVAILABILITY_MODE
availabilityModeReplicaOption_AVAILABILITY_MODE
    : 'AVAILABILITY_MODE' EqualsSign ( 'SYNCHRONOUS_COMMIT' | 'ASYNCHRONOUS_COMMIT' )
    ;

// TSql160.g:9884
// failoverModeReplicaOption where a predicate requires the next word to be FAILOVER_MODE
failoverModeReplicaOption_FAILOVER_MODE
    : 'FAILOVER_MODE' EqualsSign ( 'AUTOMATIC' | 'MANUAL' )
    ;

// TSql160.g:9938
// secondaryRoleReplicaOption where a predicate requires the next word to be SECONDARY_ROLE
secondaryRoleReplicaOption_SECONDARY_ROLE
    : 'SECONDARY_ROLE' LeftParenthesis 'ALLOW_CONNECTIONS' EqualsSign ( ( 'NO' | 'READ_ONLY' ) | All ) RightParenthesis
    ;

// TSql160.g:9903
// primaryRoleReplicaOption where a predicate requires the next word to be PRIMARY_ROLE
primaryRoleReplicaOption_PRIMARY_ROLE
    : 'PRIMARY_ROLE' LeftParenthesis 'ALLOW_CONNECTIONS' EqualsSign ( ( 'NO' | 'READ_WRITE' ) | All ) RightParenthesis
    ;

// TSql160.g:10499
// keySourceKeySpec where a predicate requires the next word to be KEY_SOURCE
keySourceKeySpec_KEY_SOURCE
    : 'KEY_SOURCE' EqualsSign stringLiteral
    ;

// TSql160.g:10531
// algorithmKeySpec where a predicate requires the next word to be ALGORITHM
algorithmKeySpec_ALGORITHM
    : 'ALGORITHM' EqualsSign Identifier
    ;

// TSql160.g:10546
// providerKeyNameSourceOption where a predicate requires the next word to be PROVIDER_KEY_NAME
providerKeyNameSourceOption_PROVIDER_KEY_NAME
    : 'PROVIDER_KEY_NAME' EqualsSign stringLiteral
    ;

// TSql160.g:11145
// alterApplicationRoleStatement where a predicate requires the next word to be APPLICATION
alterApplicationRoleStatement_APPLICATION
    : applicationRoleStatement_APPLICATION
    ;

// TSql160.g:11150
// alterAssemblyStatement where a predicate requires the next word to be ASSEMBLY
alterAssemblyStatement_ASSEMBLY
    : 'ASSEMBLY' identifier ( alterAssemblyFromClause )? ( alterAssemblyWith )? ( alterAssemblyDropFile )? ( alterAssemblyAddFile )?
    ;

// TSql160.g:11321
// alterAsymmetricKeyStatement where a predicate requires the next word to be ASYMMETRIC
alterAsymmetricKeyStatement_ASYMMETRIC
    : 'ASYMMETRIC' Key identifier ( 'REMOVE' ( 'PRIVATE' Key | 'ATTESTED' Option ) | With 'PRIVATE' Key LeftParenthesis passwordChangeOption ( Comma passwordChangeOption )* RightParenthesis | Identifier By stringLiteral )
    ;

// TSql160.g:9606
// alterAvailabilityGroupStatement where a predicate requires the next word to be AVAILABILITY
alterAvailabilityGroupStatement_AVAILABILITY
    : 'AVAILABILITY' Group identifier ( alterAvailabilityGroupAddDatabase | alterAvailabilityGroupRemoveDatabase | alterAvailabilityGroupAddReplica | alterAvailabilityGroupModifyReplica_MODIFY | alterAvailabilityGroupRemoveReplica_REMOVE | alterAvailabilityGroupSetOption | alterAvailabilityGroupTakeAction )
    ;

// TSql160.g:22814
// alterBrokerPriorityStatement where a predicate requires the next word to be BROKER
alterBrokerPriorityStatement_BROKER
    : brokerPriorityStatementBody_BROKER
    ;

// TSql160.g:22008
// alterCryptographicProviderStatement where a predicate requires the next word to be CRYPTOGRAPHIC
alterCryptographicProviderStatement_CRYPTOGRAPHIC
    : 'CRYPTOGRAPHIC' 'PROVIDER' identifier ( cryptographicProviderFile | Identifier )
    ;

// TSql160.g:7725
// alterEventSessionStatement where a predicate requires the next word to be EVENT
alterEventSessionStatement_EVENT
    : 'EVENT' 'SESSION' identifier eventSessionScope ( ( ( eventDeclarationList | targetDeclarationList ) optSessionOptionList ) | ( ( dropEventDeclarationList | dropTargetDeclarationList ) optSessionOptionList ) | ( sessionOptionList ) | ( ( 'STATE' EqualsSign ( 'START' | 'STOP' ) ) ) )
    ;

// TSql160.g:11372
// alterRemoteServiceBindingStatement where a predicate requires the next word to be REMOTE
alterRemoteServiceBindingStatement_REMOTE
    : 'REMOTE' 'SERVICE' 'BINDING' identifier With ( ( bindingUserOption ( Comma bindingAnonymousOption )? ) | ( bindingAnonymousOption ( Comma bindingUserOption )? ) )
    ;

// TSql160.g:22046
// alterResourceStatement where a predicate requires the next word to be RESOURCE
// TSql160.g:22052: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Pool)}?
// TSql160.g:22055: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Governor)}?
alterResourceStatement_RESOURCE
    : 'RESOURCE' ( alterResourcePoolStatement_POOL | alterResourceGovernorStatement_GOVERNOR )
    ;

// TSql160.g:22461
// alterWorkloadGroupStatement where a predicate requires the next word to be WORKLOAD
alterWorkloadGroupStatement_WORKLOAD
    : workloadGroupStatementBody_WORKLOAD
    ;

// TSql160.g:11415
// alterCertificateStatement where a predicate requires the next word to be CERTIFICATE
alterCertificateStatement_CERTIFICATE
    : 'CERTIFICATE' identifier ( 'REMOVE' ( 'PRIVATE' Key | 'ATTESTED' Option ) | With ( privateKeySpec | createCertificateActivityFlag ) | 'ATTESTED' By stringLiteral )
    ;

// TSql160.g:11478
// alterCredentialStatement where a predicate requires the next word to be CREDENTIAL
alterCredentialStatement_CREDENTIAL
    : credentialStatementBody_CREDENTIAL
    ;

// TSql160.g:11485
// alterEndpointStatement where a predicate requires the next word to be ENDPOINT
alterEndpointStatement_ENDPOINT
    : 'ENDPOINT' identifier ( endpointOptions )? ( protocolInfo )? ( payloadInfo )?
    ;

// TSql160.g:11498
// alterFulltextStatement where a predicate requires the next word to be FULLTEXT
alterFulltextStatement_FULLTEXT
    : 'FULLTEXT' ( alterFulltextCatalogStatement | alterFulltextIndexStatement | alterFulltextStoplistStatement )
    ;

// TSql160.g:11938
// alterLoginStatement where a predicate requires the next word to be LOGIN
alterLoginStatement_LOGIN
    : 'LOGIN' identifier ( alterLoginEnableDisable | alterLoginOptions | alterLoginAddDropCredential )
    ;

// TSql160.g:12079
// alterMasterKeyStatement where a predicate requires the next word to be MASTER
alterMasterKeyStatement_MASTER
    : 'MASTER' Key ( 'REGENERATE' With 'ENCRYPTION' By 'PASSWORD' EqualsSign stringLiteral | 'FORCE' 'REGENERATE' With 'ENCRYPTION' By 'PASSWORD' EqualsSign stringLiteral | Add 'ENCRYPTION' By ( ( 'SERVICE' 'MASTER' Key ) | ( 'PASSWORD' EqualsSign stringLiteral ) ) | Drop 'ENCRYPTION' By ( ( 'SERVICE' 'MASTER' Key ) | ( 'PASSWORD' EqualsSign stringLiteral ) ) )
    ;

// TSql160.g:12152
// alterMessageTypeStatement where a predicate requires the next word to be MESSAGE
alterMessageTypeStatement_MESSAGE
    : 'MESSAGE' 'TYPE' identifier messageTypeValidation
    ;

// TSql160.g:12218
// alterPartitionStatement where a predicate requires the next word to be PARTITION
alterPartitionStatement_PARTITION
    : 'PARTITION' ( alterPartitionFunctionStatement | alterPartitionSchemeStatement )
    ;

// TSql160.g:12506
// alterQueueStatement where a predicate requires the next word to be QUEUE
alterQueueStatement_QUEUE
    : 'QUEUE' schemaObjectThreePartName queueOptionList
    ;

// TSql160.g:12519
// alterRoleStatement where a predicate requires the next word to be ROLE
alterRoleStatement_ROLE
    : 'ROLE' identifier alterRoleAction
    ;

// TSql160.g:12636
// alterRouteStatement where a predicate requires the next word to be ROUTE
alterRouteStatement_ROUTE
    : 'ROUTE' identifier routeOptionList
    ;

// TSql160.g:8826
// alterSearchPropertyListStatement where a predicate requires the next word to be SEARCH
alterSearchPropertyListStatement_SEARCH
    : 'SEARCH' 'PROPERTY' 'LIST' identifier searchPropertyListAction requiredSemicolon
    ;

// TSql160.g:12689
// alterSequenceStatement where a predicate requires the next word to be SEQUENCE
alterSequenceStatement_SEQUENCE
    : 'SEQUENCE' schemaObjectTwoPartName ( alterSequenceOptionList )?
    ;

// TSql160.g:12722
// alterServiceStatements where a predicate requires the next word to be SERVICE
alterServiceStatements_SERVICE
    : 'SERVICE' ( alterServiceMasterKeyStatement | alterServiceStatement )
    ;

// TSql160.g:10684
// alterSymmetricKeyStatement where a predicate requires the next word to be SYMMETRIC
alterSymmetricKeyStatement_SYMMETRIC
    : 'SYMMETRIC' Key identifier ( Add | Drop ) encryptorListWithNoPassword
    ;

// TSql160.g:2136
// alterServerStatements where a predicate requires the next word to be SERVER
alterServerStatements_SERVER
    : 'SERVER' ( alterServerAuditStatements_AUDIT | alterServerConfigurationStatement_CONFIGURATION | alterServerRoleStatement_ROLE )
    ;

// TSql160.g:12859
// alterXmlSchemaCollectionStatement where a predicate requires the next word to be XML
alterXmlSchemaCollectionStatement_XML
    : 'XML' Schema 'COLLECTION' schemaObjectNonEmptyTwoPartName Add expression
    ;

// TSql160.g:23439
// alterSecurityPolicyStatement where a predicate requires the next word to be SECURITY
alterSecurityPolicyStatement_SECURITY
    : 'SECURITY' 'POLICY' schemaObjectThreePartName ( ( With LeftParenthesis securityPolicyOptionsList RightParenthesis ) | ( alterSecurityPredicateListElement ( Comma alterSecurityPredicateListElement )* ) | Add Not For Replication | Drop Not For Replication )
    ;

// TSql160.g:11629
// enableDisableAlterFulltextIndexAction where a predicate requires the next word to be ENABLE
enableDisableAlterFulltextIndexAction_ENABLE
    : 'ENABLE'
    ;

// TSql160.g:11629
// enableDisableAlterFulltextIndexAction where a predicate requires the next word to be DISABLE
enableDisableAlterFulltextIndexAction_DISABLE
    : 'DISABLE'
    ;

// TSql160.g:11726
// startPopulationAlterFulltextIndexAction where a predicate requires the next word to be START
startPopulationAlterFulltextIndexAction_START
    : 'START' ( 'INCREMENTAL' | Update | Full ) 'POPULATION'
    ;

// TSql160.g:11638
// setChangeTrackingAlterFulltextIndexAction where a predicate requires the next word to be CHANGE_TRACKING
setChangeTrackingAlterFulltextIndexAction_CHANGE_TRACKING
    : 'CHANGE_TRACKING' ( EqualsSign )? ( ( 'MANUAL' | 'AUTO' ) | Off )
    ;

// TSql160.g:8608
// searchPropertyListFullTextIndexOption where a predicate requires the next word to be SEARCH
searchPropertyListFullTextIndexOption_SEARCH
    : 'SEARCH' 'PROPERTY' 'LIST' ( EqualsSign )? ( identifier | Off )
    ;

// TSql160.g:11614
// setSearchPropertyListAlterFullTextIndexAction where a predicate requires the next word to be SEARCH
setSearchPropertyListAlterFullTextIndexAction_SEARCH
    : searchPropertyListFullTextIndexOption_SEARCH populationOption
    ;

// TSql160.g:12369
// queueActivationOption where a predicate requires the next word to be ACTIVATION
queueActivationOption_ACTIVATION
    : 'ACTIVATION' LeftParenthesis activationQueueOptionArgument ( Comma activationQueueOptionArgument )* RightParenthesis
    ;

// TSql160.g:13181
// inlineFunctionOption where a predicate requires the next word to be INLINE
inlineFunctionOption_INLINE
    : 'INLINE' EqualsSign ( On | Off )
    ;

// TSql160.g:13373
// incrementalStatisticsOption where a predicate requires the next word to be INCREMENTAL
incrementalStatisticsOption_INCREMENTAL
    : 'INCREMENTAL' EqualsSign optionOnOff
    ;

// TSql160.g:13386
// autoDropStatisticsOption where a predicate requires the next word to be AUTO_DROP
autoDropStatisticsOption_AUTO_DROP
    : 'AUTO_DROP' EqualsSign optionOnOff
    ;

// TSql160.g:13399
// persistSamplePercentStatisticsOption where a predicate requires the next word to be PERSIST_SAMPLE_PERCENT
persistSamplePercentStatisticsOption_PERSIST_SAMPLE_PERCENT
    : 'PERSIST_SAMPLE_PERCENT' EqualsSign optionOnOff
    ;

// TSql160.g:13344
// resampleStatisticsOption where a predicate requires the next word to be RESAMPLE
resampleStatisticsOption_RESAMPLE
    : 'RESAMPLE' ( On 'PARTITIONS' LeftParenthesis statisticsPartitionRange ( Comma statisticsPartitionRange )* RightParenthesis )?
    ;

// TSql160.g:13299
// statsStreamStatisticsOption where a predicate requires the next word to be STATS_STREAM
statsStreamStatisticsOption_STATS_STREAM
    : 'STATS_STREAM' EqualsSign binary
    ;

// TSql160.g:13748
// killQueryNotificationSubscriptionStatement where a predicate requires the next word to be QUERY
killQueryNotificationSubscriptionStatement_QUERY
    : 'QUERY' 'NOTIFICATION' 'SUBSCRIPTION' ( All | integer )
    ;

// TSql160.g:14407
// openMasterKeyStatement where a predicate requires the next word to be MASTER
openMasterKeyStatement_MASTER
    : 'MASTER' Key 'DECRYPTION' By 'PASSWORD' EqualsSign stringLiteral
    ;

// TSql160.g:14463
// openSymmetricKeyStatement where a predicate requires the next word to be SYMMETRIC
openSymmetricKeyStatement_SYMMETRIC
    : 'SYMMETRIC' Key identifier 'DECRYPTION' By decryptionMechanism
    ;

// TSql160.g:14515
// closeMasterKeyStatement where a predicate requires the next word to be MASTER
closeMasterKeyStatement_MASTER
    : 'MASTER' Key
    ;

// TSql160.g:4277
// addSensitivityClassificationStatement where a predicate requires the next word to be SENSITIVITY
addSensitivityClassificationStatement_SENSITIVITY
    : 'SENSITIVITY' 'CLASSIFICATION' To ( column ( Comma column )* ) With LeftParenthesis sensitivityClassificationOption ( Comma sensitivityClassificationOption )* RightParenthesis
    ;

// TSql160.g:15210
// dropAggregateStatement where a predicate requires the next word to be AGGREGATE
dropAggregateStatement_AGGREGATE
    : 'AGGREGATE' ( If Exists )? dropObjectList
    ;

// TSql160.g:15224
// dropAssemblyStatement where a predicate requires the next word to be ASSEMBLY
dropAssemblyStatement_ASSEMBLY
    : 'ASSEMBLY' ( If Exists )? dropObjectList ( With 'NO' 'DEPENDENTS' )?
    ;

// TSql160.g:15248
// dropAvailabilityGroupStatement where a predicate requires the next word to be AVAILABILITY
dropAvailabilityGroupStatement_AVAILABILITY
    : 'AVAILABILITY' Group identifier
    ;

// TSql160.g:15188
// dropPartitionStatements where a predicate requires the next word to be PARTITION
dropPartitionStatements_PARTITION
    : 'PARTITION' ( dropPartitionFunctionStatement | dropPartitionSchemeStatement )
    ;

// TSql160.g:15270
// dropSynonymStatement where a predicate requires the next word to be SYNONYM
dropSynonymStatement_SYNONYM
    : 'SYNONYM' ( If Exists )? dropObjectList
    ;

// TSql160.g:15284
// dropApplicationRoleStatement where a predicate requires the next word to be APPLICATION
dropApplicationRoleStatement_APPLICATION
    : 'APPLICATION' 'ROLE' identifier
    ;

// TSql160.g:22821
// dropBrokerPriorityStatement where a predicate requires the next word to be BROKER
dropBrokerPriorityStatement_BROKER
    : 'BROKER' 'PRIORITY' identifier
    ;

// TSql160.g:22034
// dropCryptographicProviderStatement where a predicate requires the next word to be CRYPTOGRAPHIC
dropCryptographicProviderStatement_CRYPTOGRAPHIC
    : 'CRYPTOGRAPHIC' 'PROVIDER' identifier
    ;

// TSql160.g:15296
// dropFulltextStatements where a predicate requires the next word to be FULLTEXT
dropFulltextStatements_FULLTEXT
    : 'FULLTEXT' ( dropFulltextCatalogStatement | dropFulltextIndexStatement | dropFulltextStoplistStatement )
    ;

// TSql160.g:15331
// dropLoginStatement where a predicate requires the next word to be LOGIN
dropLoginStatement_LOGIN
    : 'LOGIN' identifier
    ;

// TSql160.g:22272
// dropResourcePoolStatement where a predicate requires the next word to be RESOURCE
dropResourcePoolStatement_RESOURCE
    : 'RESOURCE' 'POOL' identifier
    ;

// TSql160.g:15705
// dropWorkloadStatements where a predicate requires the next word to be WORKLOAD
// TSql160.g:15711: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Group)}?
// TSql160.g:15714: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Classifier)}?
dropWorkloadStatements_WORKLOAD
    : 'WORKLOAD' ( dropWorkloadGroupStatement_GROUP | dropWorkloadClassifierStatement_CLASSIFIER )
    ;

// TSql160.g:15342
// dropRoleStatement where a predicate requires the next word to be ROLE
dropRoleStatement_ROLE
    : 'ROLE' ( If Exists )? identifier
    ;

// TSql160.g:15360
// dropMasterKeyStatement where a predicate requires the next word to be MASTER
dropMasterKeyStatement_MASTER
    : 'MASTER' Key
    ;

// TSql160.g:15368
// dropSymmetricKeyStatement where a predicate requires the next word to be SYMMETRIC
dropSymmetricKeyStatement_SYMMETRIC
    : 'SYMMETRIC' Key identifier removeProviderKeyOpt
    ;

// TSql160.g:15384
// dropAsymmetricKeyStatement where a predicate requires the next word to be ASYMMETRIC
dropAsymmetricKeyStatement_ASYMMETRIC
    : 'ASYMMETRIC' Key identifier removeProviderKeyOpt
    ;

// TSql160.g:15414
// dropCertificateStatement where a predicate requires the next word to be CERTIFICATE
dropCertificateStatement_CERTIFICATE
    : 'CERTIFICATE' identifier
    ;

// TSql160.g:15439
// dropCredentialStatement where a predicate requires the next word to be CREDENTIAL
dropCredentialStatement_CREDENTIAL
    : 'CREDENTIAL' identifier
    ;

// TSql160.g:15451
// dropTypeStatement where a predicate requires the next word to be TYPE
dropTypeStatement_TYPE
    : 'TYPE' ( If Exists )? dropObject
    ;

// TSql160.g:15469
// dropXmlSchemaCollectionStatement where a predicate requires the next word to be XML
dropXmlSchemaCollectionStatement_XML
    : 'XML' Schema 'COLLECTION' schemaObjectNonEmptyTwoPartName
    ;

// TSql160.g:15499
// dropContractStatement where a predicate requires the next word to be CONTRACT
dropContractStatement_CONTRACT
    : 'CONTRACT' identifier
    ;

// TSql160.g:15510
// dropQueueStatement where a predicate requires the next word to be QUEUE
dropQueueStatement_QUEUE
    : 'QUEUE' schemaObjectThreePartName
    ;

// TSql160.g:15521
// dropServiceStatement where a predicate requires the next word to be SERVICE
dropServiceStatement_SERVICE
    : 'SERVICE' identifier
    ;

// TSql160.g:15532
// dropRouteStatement where a predicate requires the next word to be ROUTE
dropRouteStatement_ROUTE
    : 'ROUTE' identifier
    ;

// TSql160.g:15543
// dropMessageTypeStatement where a predicate requires the next word to be MESSAGE
dropMessageTypeStatement_MESSAGE
    : 'MESSAGE' 'TYPE' identifier
    ;

// TSql160.g:15555
// dropRemoteServiceBindingStatement where a predicate requires the next word to be REMOTE
dropRemoteServiceBindingStatement_REMOTE
    : 'REMOTE' 'SERVICE' 'BINDING' identifier
    ;

// TSql160.g:15568
// dropEndpointStatement where a predicate requires the next word to be ENDPOINT
dropEndpointStatement_ENDPOINT
    : 'ENDPOINT' identifier
    ;

// TSql160.g:15579
// dropEventStatement where a predicate requires the next word to be EVENT
// TSql160.g:15585: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Notification)}?
// TSql160.g:15588: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Session)}?
dropEventStatement_EVENT
    : 'EVENT' ( dropEventNotificationStatement_NOTIFICATION | dropEventSessionStatement_SESSION )
    ;

// TSql160.g:8896
// dropSearchPropertyListStatement where a predicate requires the next word to be SEARCH
dropSearchPropertyListStatement_SEARCH
    : 'SEARCH' 'PROPERTY' 'LIST' identifier requiredSemicolon
    ;

// TSql160.g:15680
// dropSequenceStatement where a predicate requires the next word to be SEQUENCE
dropSequenceStatement_SEQUENCE
    : 'SEQUENCE' ( If Exists )? dropObjectList
    ;

// TSql160.g:15694
// dropFederationStatement where a predicate requires the next word to be FEDERATION
dropFederationStatement_FEDERATION
    : 'FEDERATION' identifier
    ;

// TSql160.g:23648
// dropSecurityPolicyStatement where a predicate requires the next word to be SECURITY
dropSecurityPolicyStatement_SECURITY
    : 'SECURITY' 'POLICY' ( If Exists )? dropObjectList
    ;

// TSql160.g:4379
// dropSensitivityClassificationStatement where a predicate requires the next word to be SENSITIVITY
dropSensitivityClassificationStatement_SENSITIVITY
    : 'SENSITIVITY' 'CLASSIFICATION' From ( column ( Comma column )* )
    ;

// TSql160.g:14848
// dropServerAuditSpecificationStatement where a predicate requires the next word to be SPECIFICATION
dropServerAuditSpecificationStatement_SPECIFICATION
    : 'SPECIFICATION' identifier
    ;

// TSql160.g:14906
// dropDatabaseAuditSpecificationStatement where a predicate requires the next word to be AUDIT
dropDatabaseAuditSpecificationStatement_AUDIT
    : 'AUDIT' 'SPECIFICATION' identifier
    ;

// TSql160.g:15425
// dropDatabaseScopedCredentialStatement where a predicate requires the next word to be SCOPED
dropDatabaseScopedCredentialStatement_SCOPED
    : 'SCOPED' 'CREDENTIAL' identifier
    ;

// TSql160.g:17366
// maxDopOption where a predicate requires the next word to be MAXDOP
maxDopOption_MAXDOP
    : 'MAXDOP' EqualsSign integer
    ;

// TSql160.g:17381
// onlineIndexOption where a predicate requires the next word to be ONLINE
onlineIndexOption_ONLINE
    : 'ONLINE' EqualsSign ( On ( LeftParenthesis lowPriorityLockWaitOption RightParenthesis )? | Off )
    ;

// TSql160.g:15009
// dropIndexMoveToOption where a predicate requires the next word to be MOVE
dropIndexMoveToOption_MOVE
    : 'MOVE' To filegroupOrPartitionScheme
    ;

// TSql160.g:29200
// fileStreamOn where a predicate requires the next word to be FILESTREAM_ON
fileStreamOn_FILESTREAM_ON
    : 'FILESTREAM_ON' stringOrIdentifier
    ;

// TSql160.g:15022
// dropIndexFileStreamOnOption where a predicate requires the next word to be FILESTREAM_ON
dropIndexFileStreamOnOption_FILESTREAM_ON
    : fileStreamOn_FILESTREAM_ON
    ;

// TSql160.g:27549
// dataCompressionOption where a predicate requires the next word to be DATA_COMPRESSION
dataCompressionOption_DATA_COMPRESSION
    : 'DATA_COMPRESSION' EqualsSign Identifier ( On 'PARTITIONS' LeftParenthesis compressionPartitionRange ( Comma compressionPartitionRange )* RightParenthesis )?
    ;

// TSql160.g:27578
// xmlCompressionOption where a predicate requires the next word to be XML_COMPRESSION
xmlCompressionOption_XML_COMPRESSION
    : 'XML_COMPRESSION' EqualsSign ( On (  ) | Off (  ) ) ( On 'PARTITIONS' LeftParenthesis compressionPartitionRange ( Comma compressionPartitionRange )* RightParenthesis )?
    ;

// TSql160.g:17417
// lowPriorityLockWaitOption where a predicate requires the next word to be WAIT_AT_LOW_PRIORITY
lowPriorityLockWaitOption_WAIT_AT_LOW_PRIORITY
    : 'WAIT_AT_LOW_PRIORITY' LeftParenthesis 'MAX_DURATION' EqualsSign integer ( 'MINUTES' )? Comma 'ABORT_AFTER_WAIT' EqualsSign Identifier RightParenthesis
    ;

// TSql160.g:27542
// waitAtLowPriorityOption where a predicate requires the next word to be WAIT_AT_LOW_PRIORITY
waitAtLowPriorityOption_WAIT_AT_LOW_PRIORITY
    : lowPriorityLockWaitOption_WAIT_AT_LOW_PRIORITY
    ;

// TSql160.g:15594
// dropEventNotificationStatement where a predicate requires the next word to be NOTIFICATION
dropEventNotificationStatement_NOTIFICATION
    : 'NOTIFICATION' identifier ( Comma identifier )* eventNotificationObjectScope
    ;

// TSql160.g:15633
// dropEventSessionStatement where a predicate requires the next word to be SESSION
dropEventSessionStatement_SESSION
    : 'SESSION' identifier eventSessionDropScope
    ;

// TSql160.g:23877
// dropExternalDataSourceStatement where a predicate requires the next word to be DATA
dropExternalDataSourceStatement_DATA
    : 'DATA' 'SOURCE' identifier
    ;

// TSql160.g:1129
// dropExternalLibraryStatement where a predicate requires the next word to be LIBRARY
dropExternalLibraryStatement_LIBRARY
    : 'LIBRARY' identifier authorizationOpt
    ;

// TSql160.g:1311
// dropExternalLanguageStatement where a predicate requires the next word to be LANGUAGE
dropExternalLanguageStatement_LANGUAGE
    : 'LANGUAGE' identifier authorizationOpt
    ;

// TSql160.g:22314
// dropExternalResourcePoolStatement where a predicate requires the next word to be RESOURCE
dropExternalResourcePoolStatement_RESOURCE
    : 'RESOURCE' 'POOL' identifier
    ;

// TSql160.g:23224
// dropColumnMasterKeyStatement where a predicate requires the next word to be MASTER
dropColumnMasterKeyStatement_MASTER
    : 'MASTER' Key identifier
    ;

// TSql160.g:23384
// dropColumnEncryptionKeyStatement where a predicate requires the next word to be ENCRYPTION
dropColumnEncryptionKeyStatement_ENCRYPTION
    : 'ENCRYPTION' Key identifier
    ;

// TSql160.g:15719
// dropWorkloadGroupStatement where a predicate requires the next word to be GROUP
dropWorkloadGroupStatement_GROUP
    : 'GROUP' identifier
    ;

// TSql160.g:15730
// dropWorkloadClassifierStatement where a predicate requires the next word to be CLASSIFIER
dropWorkloadClassifierStatement_CLASSIFIER
    : 'CLASSIFIER' identifier
    ;

// TSql160.g:16238
// dateFirstOption where a predicate requires the next word to be DATEFIRST
dateFirstOption_DATEFIRST
    : 'DATEFIRST' EqualsSign integer
    ;

// TSql160.g:16250
// dateFormatOption where a predicate requires the next word to be DATEFORMAT
dateFormatOption_DATEFORMAT
    : 'DATEFORMAT' EqualsSign stringLiteral
    ;

// TSql160.g:16262
// delayedDurabilityOption where a predicate requires the next word to be DELAYED_DURABILITY
delayedDurabilityOption_DELAYED_DURABILITY
    : 'DELAYED_DURABILITY' EqualsSign optionOnOff
    ;

// TSql160.g:16307
// languageOption where a predicate requires the next word to be LANGUAGE
languageOption_LANGUAGE
    : 'LANGUAGE' EqualsSign stringLiteral
    ;

// TSql160.g:17552
// bucketCountOption where a predicate requires the next word to be BUCKET_COUNT
bucketCountOption_BUCKET_COUNT
    : 'BUCKET_COUNT' EqualsSign integer
    ;

// TSql160.g:17574
// compressionDelayOption where a predicate requires the next word to be COMPRESSION_DELAY
compressionDelayOption_COMPRESSION_DELAY
    : 'COMPRESSION_DELAY' EqualsSign ( integer ( Identifier )? )
    ;

// TSql160.g:17460
// ignoreDupKeyIndexOption where a predicate requires the next word to be IGNORE_DUP_KEY
ignoreDupKeyIndexOption_IGNORE_DUP_KEY
    : 'IGNORE_DUP_KEY' EqualsSign ( On ( LeftParenthesis suppressMessagesOption RightParenthesis )? | Off )
    ;

// TSql160.g:27521
// maxDurationOption where a predicate requires the next word to be MAX_DURATION
maxDurationOption_MAX_DURATION
    : 'MAX_DURATION' EqualsSign integer ( 'MINUTES' )?
    ;

// TSql160.g:17703
// xmlNamespaces where a predicate requires the next word to be XMLNAMESPACES
xmlNamespaces_XMLNAMESPACES
    : 'XMLNAMESPACES' LeftParenthesis xmlNamespacesElement ( Comma xmlNamespacesElement )* RightParenthesis
    ;

// TSql160.g:17672
// ctContext where a predicate requires the next word to be CHANGE_TRACKING_CONTEXT
ctContext_CHANGE_TRACKING_CONTEXT
    : 'CHANGE_TRACKING_CONTEXT' LeftParenthesis binaryOrVariable RightParenthesis
    ;

// TSql160.g:21546
// xmlForClause where a predicate requires the next word to be XML
xmlForClause_XML
    : 'XML' xmlForClauseMode ( Comma xmlParam )*
    ;

// TSql160.g:21636
// jsonForClause where a predicate requires the next word to be JSON
jsonForClause_JSON
    : 'JSON' jsonForClauseMode ( Comma jsonParam )*
    ;

// TSql160.g:18739
// extractFromExpression where a predicate requires the next word to be EXTRACT
extractFromExpression_EXTRACT
    : 'EXTRACT' LeftParenthesis identifier From expression RightParenthesis
    ;

// TSql160.g:19156
// changeTableTableReference where a predicate requires the next word to be CHANGETABLE
changeTableTableReference_CHANGETABLE
    : 'CHANGETABLE' LeftParenthesis Identifier schemaObjectFourPartName Comma ( changesChangeTableParams | versionChangeTableParams ) RightParenthesis simpleTableReferenceAlias ( columnNameList )?
    ;

// TSql160.g:20580
// openJsonTableReference where a predicate requires the next word to be OPENJSON
openJsonTableReference_OPENJSON
    : 'OPENJSON' LeftParenthesis openJsonParams RightParenthesis openJsonWithClauseOpt simpleTableReferenceAliasOpt
    ;

// TSql160.g:19040
// predictTableReference where a predicate requires the next word to be PREDICT
predictTableReference_PREDICT
    : 'PREDICT' LeftParenthesis predictParams RightParenthesis predictWithClauseOpt simpleTableReferenceAliasOpt
    ;

// TSql160.g:19631
// outputClause where a predicate requires the next word to be OUTPUT
outputClause_OUTPUT
    : 'OUTPUT' outputClauseSelectElement ( Comma outputClauseSelectElement )*
    ;

// TSql160.g:19562
// outputIntoOutputClause where a predicate requires the next word to be OUTPUT
outputIntoOutputClause_OUTPUT
    : 'OUTPUT' outputClauseSelectElement ( Comma outputClauseSelectElement )* ( ( outputIntoClause ( outputClause_OUTPUT )? ) |  )
    ;

// TSql160.g:21130
// forceSeekTableHint where a predicate requires the next word to be FORCESEEK
forceSeekTableHint_FORCESEEK
    : 'FORCESEEK' ( LeftParenthesis identifierOrInteger LeftParenthesis identifierColumnReferenceExpression ( Comma identifierColumnReferenceExpression )* RightParenthesis RightParenthesis )?
    ;

// TSql160.g:21372
// cubeSpec where a predicate requires the next word to be CUBE
cubeSpec_CUBE
    : 'CUBE' LeftParenthesis cubeOrRollupArgumentList RightParenthesis
    ;

// TSql160.g:21358
// rollupSpec where a predicate requires the next word to be ROLLUP
rollupSpec_ROLLUP
    : 'ROLLUP' LeftParenthesis cubeOrRollupArgumentList RightParenthesis
    ;

// TSql160.g:21284
// groupingSetsSpec where a predicate requires the next word to be GROUPING
groupingSetsSpec_GROUPING
    : 'GROUPING' 'SETS' LeftParenthesis groupingSet ( Comma groupingSet )* RightParenthesis
    ;

// TSql160.g:22259
// alterResourcePoolStatement where a predicate requires the next word to be POOL
alterResourcePoolStatement_POOL
    : 'POOL' identifier resourcePoolStatementBody
    ;

// TSql160.g:22060
// alterResourceGovernorStatement where a predicate requires the next word to be GOVERNOR
alterResourceGovernorStatement_GOVERNOR
    : 'GOVERNOR' ( 'DISABLE' | Reconfigure | With LeftParenthesis 'CLASSIFIER_FUNCTION' EqualsSign ( schemaObjectNonEmptyTwoPartName | Null ) RightParenthesis | 'RESET' Statistics )
    ;

// TSql160.g:22202
// resourcePoolAffinitySpecification where a predicate requires the next word to be AFFINITY
resourcePoolAffinitySpecification_AFFINITY
    : 'AFFINITY' Identifier EqualsSign ( 'AUTO' | LeftParenthesis poolAffinityRange ( Comma poolAffinityRange )* RightParenthesis )
    ;

// TSql160.g:22388
// externalResourcePoolAffinitySpecification where a predicate requires the next word to be AFFINITY
externalResourcePoolAffinitySpecification_AFFINITY
    : 'AFFINITY' Identifier EqualsSign ( 'AUTO' | LeftParenthesis externalPoolAffinityRange ( Comma externalPoolAffinityRange )* RightParenthesis )
    ;

// TSql160.g:22700
// classifierWorkloadGroupOption where a predicate requires the next word to be WORKLOAD_GROUP
classifierWorkloadGroupOption_WORKLOAD_GROUP
    : 'WORKLOAD_GROUP' EqualsSign stringLiteral
    ;

// TSql160.g:22713
// classifierMemberNameOption where a predicate requires the next word to be MEMBERNAME
classifierMemberNameOption_MEMBERNAME
    : 'MEMBERNAME' EqualsSign stringLiteral
    ;

// TSql160.g:22727
// classifierWlmLabelOption where a predicate requires the next word to be WLM_LABEL
classifierWlmLabelOption_WLM_LABEL
    : 'WLM_LABEL' EqualsSign stringLiteral
    ;

// TSql160.g:22742
// classifierWlmContextOption where a predicate requires the next word to be WLM_CONTEXT
classifierWlmContextOption_WLM_CONTEXT
    : 'WLM_CONTEXT' EqualsSign stringLiteral
    ;

// TSql160.g:22756
// classifierImportanceOption where a predicate requires the next word to be IMPORTANCE
classifierImportanceOption_IMPORTANCE
    : 'IMPORTANCE' EqualsSign Identifier
    ;

// TSql160.g:22767
// classifierStartTimeOption where a predicate requires the next word to be START_TIME
classifierStartTimeOption_START_TIME
    : 'START_TIME' EqualsSign wlmTimeLiteral
    ;

// TSql160.g:22780
// classifierEndTimeOption where a predicate requires the next word to be END_TIME
classifierEndTimeOption_END_TIME
    : 'END_TIME' EqualsSign wlmTimeLiteral
    ;

// TSql160.g:23098
// sequenceRestartOptionListElement where a predicate requires the next word to be RESTART
sequenceRestartOptionListElement_RESTART
    : 'RESTART' ( With seedIncrement )?
    ;

// TSql160.g:23086
// sequenceStartOptionListElement where a predicate requires the next word to be START
sequenceStartOptionListElement_START
    : 'START' With seedIncrement
    ;

// TSql160.g:23028
// sequenceCacheOptionListElement where a predicate requires the next word to be CACHE
sequenceCacheOptionListElement_CACHE
    : 'CACHE' ( seedIncrement )?
    ;

// TSql160.g:23040
// sequenceCycleOptionListElement where a predicate requires the next word to be CYCLE
sequenceCycleOptionListElement_CYCLE
    : 'CYCLE'
    ;

// TSql160.g:23048
// sequenceIncrementOptionListElement where a predicate requires the next word to be INCREMENT
sequenceIncrementOptionListElement_INCREMENT
    : 'INCREMENT' By seedIncrement
    ;

// TSql160.g:23060
// sequenceNoOptionListElement where a predicate requires the next word to be NO
sequenceNoOptionListElement_NO
    : 'NO' ( 'MINVALUE' | 'MAXVALUE' | 'CACHE' | 'CYCLE' )
    ;

// TSql160.g:23148
// createColumnMasterKeyStatement where a predicate requires the next word to be MASTER
createColumnMasterKeyStatement_MASTER
    : 'MASTER' Key identifier With LeftParenthesis columnMasterkeyParameter Comma columnMasterkeyParameter ( Comma columnMasterkeyParameter )* RightParenthesis
    ;

// TSql160.g:23235
// createColumnEncryptionKeyStatement where a predicate requires the next word to be ENCRYPTION
createColumnEncryptionKeyStatement_ENCRYPTION
    : 'ENCRYPTION' Key identifier With Values createColumnEncryptionKeyValue ( Comma createColumnEncryptionKeyValue )?
    ;

// TSql160.g:23187
// columnMasterKeyStoreProviderNameParameter where a predicate requires the next word to be KEY_STORE_PROVIDER_NAME
columnMasterKeyStoreProviderNameParameter_KEY_STORE_PROVIDER_NAME
    : 'KEY_STORE_PROVIDER_NAME' EqualsSign stringLiteral
    ;

// TSql160.g:23199
// columnMasterKeyPathParameter where a predicate requires the next word to be KEY_PATH
columnMasterKeyPathParameter_KEY_PATH
    : 'KEY_PATH' EqualsSign stringLiteral
    ;

// TSql160.g:23211
// columnMasterKeyEnclaveComputationsParameter where a predicate requires the next word to be ENCLAVE_COMPUTATIONS
columnMasterKeyEnclaveComputationsParameter_ENCLAVE_COMPUTATIONS
    : 'ENCLAVE_COMPUTATIONS' LeftParenthesis 'SIGNATURE' EqualsSign binary RightParenthesis
    ;

// TSql160.g:23292
// columnMasterKeyNameParameter where a predicate requires the next word to be COLUMN_MASTER_KEY
columnMasterKeyNameParameter_COLUMN_MASTER_KEY
    : 'COLUMN_MASTER_KEY' EqualsSign identifier
    ;

// TSql160.g:23305
// columnEncryptionAlgorithmNameParameter where a predicate requires the next word to be ALGORITHM
columnEncryptionAlgorithmNameParameter_ALGORITHM
    : 'ALGORITHM' EqualsSign stringLiteral
    ;

// TSql160.g:23318
// encryptedValueParameter where a predicate requires the next word to be ENCRYPTED_VALUE
encryptedValueParameter_ENCRYPTED_VALUE
    : 'ENCRYPTED_VALUE' EqualsSign binary
    ;

// TSql160.g:23722
// externalDataSourceType where a predicate requires the next word to be TYPE
externalDataSourceType_TYPE
    : 'TYPE' EqualsSign ( ( 'HADOOP' | 'RDBMS' | 'SHARD_MAP_MANAGER' | 'BLOB_STORAGE' ) )
    ;

// TSql160.g:23761
// externalDataSourceLocation where a predicate requires the next word to be LOCATION
externalDataSourceLocation_LOCATION
    : 'LOCATION' EqualsSign stringLiteral
    ;

// TSql160.g:23773
// externalDataSourcePushdownOption where a predicate requires the next word to be PUSHDOWN
externalDataSourcePushdownOption_PUSHDOWN
    : 'PUSHDOWN' EqualsSign ( On | Off )
    ;

// TSql160.g:23946
// externalStreamLocation where a predicate requires the next word to be LOCATION
externalStreamLocation_LOCATION
    : 'LOCATION' EqualsSign stringLiteral
    ;

// TSql160.g:23957
// externalStreamInputOptions where a predicate requires the next word to be INPUT_OPTIONS
externalStreamInputOptions_INPUT_OPTIONS
    : 'INPUT_OPTIONS' EqualsSign stringLiteral
    ;

// TSql160.g:23968
// externalStreamOutputOptions where a predicate requires the next word to be OUTPUT_OPTIONS
externalStreamOutputOptions_OUTPUT_OPTIONS
    : 'OUTPUT_OPTIONS' EqualsSign stringLiteral
    ;

// TSql160.g:24129
// externalFormatingContainerOption where a predicate requires the next word to be FORMAT_OPTIONS
externalFormatingContainerOption_FORMAT_OPTIONS
    : 'FORMAT_OPTIONS' LeftParenthesis externalDataFormatOption ( Comma externalDataFormatOption )* RightParenthesis
    ;

// TSql160.g:24118
// externalFormatingOption where a predicate requires the next word to be FORMAT_OPTIONS
externalFormatingOption_FORMAT_OPTIONS
    : externalFormatingContainerOption_FORMAT_OPTIONS
    ;

// TSql160.g:24233
// externalFileFormatLiteralOption where a predicate requires the next word to be SERDE_METHOD
externalFileFormatLiteralOption_SERDE_METHOD
    : 'SERDE_METHOD' EqualsSign ( stringLiteral )
    ;

// TSql160.g:24105
// externalFileFormatLiteralOptions where a predicate requires the next word to be SERDE_METHOD
externalFileFormatLiteralOptions_SERDE_METHOD
    : externalFileFormatLiteralOption_SERDE_METHOD
    ;

// TSql160.g:24233
// externalFileFormatLiteralOption where a predicate requires the next word to be DATA_COMPRESSION
externalFileFormatLiteralOption_DATA_COMPRESSION
    : 'DATA_COMPRESSION' EqualsSign ( stringLiteral )
    ;

// TSql160.g:24105
// externalFileFormatLiteralOptions where a predicate requires the next word to be DATA_COMPRESSION
externalFileFormatLiteralOptions_DATA_COMPRESSION
    : externalFileFormatLiteralOption_DATA_COMPRESSION
    ;

// TSql160.g:24169
// externalDataFormatUseTypeDefaultOption where a predicate requires the next word to be USE_TYPE_DEFAULT
externalDataFormatUseTypeDefaultOption_USE_TYPE_DEFAULT
    : 'USE_TYPE_DEFAULT' EqualsSign ( ( 'TRUE' | 'FALSE' ) )
    ;

// TSql160.g:24374
// boundingBoxSpatialIndexOption where a predicate requires the next word to be BOUNDING_BOX
boundingBoxSpatialIndexOption_BOUNDING_BOX
    : 'BOUNDING_BOX' EqualsSign boundingBoxParameter
    ;

// TSql160.g:24462
// gridsSpatialIndexOption where a predicate requires the next word to be GRIDS
gridsSpatialIndexOption_GRIDS
    : 'GRIDS' EqualsSign gridParameter
    ;

// TSql160.g:24536
// cellsPerObjectSpatialIndexOption where a predicate requires the next word to be CELLS_PER_OBJECT
cellsPerObjectSpatialIndexOption_CELLS_PER_OBJECT
    : 'CELLS_PER_OBJECT' EqualsSign integer
    ;

// TSql160.g:24733
// viewDistributionOption where a predicate requires the next word to be DISTRIBUTION
viewDistributionOption_DISTRIBUTION
    : 'DISTRIBUTION' EqualsSign ( viewDistributionPolicy )
    ;

// TSql160.g:24723
// viewForAppendOption where a predicate requires the next word to be FOR_APPEND
viewForAppendOption_FOR_APPEND
    : 'FOR_APPEND'
    ;

// TSql160.g:24786
// viewRoundRobinDistributionPolicy where a predicate requires the next word to be ROUND_ROBIN
viewRoundRobinDistributionPolicy_ROUND_ROBIN
    : 'ROUND_ROBIN'
    ;

// TSql160.g:24761
// viewHashDistributionPolicy where a predicate requires the next word to be HASH
viewHashDistributionPolicy_HASH
    : 'HASH' LeftParenthesis identifier ( Comma identifier )* RightParenthesis
    ;

// TSql160.g:26090
// executeOptionRecompile where a predicate requires the next word to be RECOMPILE
executeOptionRecompile_RECOMPILE
    : 'RECOMPILE'
    ;

// TSql160.g:26364
// linkedServer where a predicate requires the next word to be AT
linkedServer_AT
    : 'AT' identifier
    ;

// TSql160.g:33298
// identifier where a predicate requires the next word to be FEDERATED
identifier_FEDERATED
    : 'FEDERATED'
    ;

// TSql160.g:26650
// federatedOn where a predicate requires the next word to be FEDERATED
federatedOn_FEDERATED
    : identifier_FEDERATED On LeftParenthesis identifier EqualsSign identifier RightParenthesis
    ;

// TSql160.g:27499
// tableDataCompressionOption where a predicate requires the next word to be DATA_COMPRESSION
tableDataCompressionOption_DATA_COMPRESSION
    : dataCompressionOption_DATA_COMPRESSION
    ;

// TSql160.g:27510
// tableXmlCompressionOption where a predicate requires the next word to be XML_COMPRESSION
tableXmlCompressionOption_XML_COMPRESSION
    : xmlCompressionOption_XML_COMPRESSION
    ;

// TSql160.g:27156
// fileTableDirectoryTableOption where a predicate requires the next word to be FILETABLE_DIRECTORY
fileTableDirectoryTableOption_FILETABLE_DIRECTORY
    : 'FILETABLE_DIRECTORY' EqualsSign stringLiteralOrNull
    ;

// TSql160.g:27168
// fileTableCollateFileNameTableOption where a predicate requires the next word to be FILETABLE_COLLATE_FILENAME
fileTableCollateFileNameTableOption_FILETABLE_COLLATE_FILENAME
    : 'FILETABLE_COLLATE_FILENAME' EqualsSign nonQuotedIdentifier
    ;

// TSql160.g:26862
// memoryOptimizedTableOption where a predicate requires the next word to be MEMORY_OPTIMIZED
memoryOptimizedTableOption_MEMORY_OPTIMIZED
    : 'MEMORY_OPTIMIZED' EqualsSign optionOnOff
    ;

// TSql160.g:26874
// durabilityTableOption where a predicate requires the next word to be DURABILITY
durabilityTableOption_DURABILITY
    : 'DURABILITY' EqualsSign Identifier
    ;

// TSql160.g:26883
// remoteDataArchiveTableOption where a predicate requires the next word to be REMOTE_DATA_ARCHIVE
remoteDataArchiveTableOption_REMOTE_DATA_ARCHIVE
    : 'REMOTE_DATA_ARCHIVE' EqualsSign On ( LeftParenthesis remoteDataArchiveTableSettings RightParenthesis )
    ;

// TSql160.g:27180
// systemVersioningTableOnOffOption where a predicate requires the next word to be SYSTEM_VERSIONING
systemVersioningTableOnOffOption_SYSTEM_VERSIONING
    : 'SYSTEM_VERSIONING' EqualsSign ( Off | ( On ( LeftParenthesis systemVersioningTableSuboption ( Comma systemVersioningTableSuboption )* RightParenthesis )? ) )
    ;

// TSql160.g:27297
// ledgerTableOption where a predicate requires the next word to be LEDGER
ledgerTableOption_LEDGER
    : 'LEDGER' EqualsSign ( Off | ( On ( LeftParenthesis ledgerTableSuboption ( Comma ledgerTableSuboption )* RightParenthesis )? ) )
    ;

// TSql160.g:26917
// tableDistributionOption where a predicate requires the next word to be DISTRIBUTION
tableDistributionOption_DISTRIBUTION
    : 'DISTRIBUTION' EqualsSign ( tableDistributionPolicy )
    ;

// TSql160.g:27065
// tablePartitionOption where a predicate requires the next word to be PARTITION
tablePartitionOption_PARTITION
    : 'PARTITION' LeftParenthesis identifier partitionOptionTypes RightParenthesis
    ;

// TSql160.g:27085
// locationOption where a predicate requires the next word to be LOCATION
locationOption_LOCATION
    : 'LOCATION' EqualsSign 'USER_DB'
    ;

// TSql160.g:26778
// dataRetentionOption where a predicate requires the next word to be DATA_DELETION
dataRetentionOption_DATA_DELETION
    : 'DATA_DELETION' EqualsSign ( Off | On LeftParenthesis dataRetentionSubOption ( Comma dataRetentionSubOption )? RightParenthesis )
    ;

// TSql160.g:26973
// tableReplicateDistributionPolicy where a predicate requires the next word to be REPLICATE
tableReplicateDistributionPolicy_REPLICATE
    : 'REPLICATE'
    ;

// TSql160.g:26982
// tableRoundRobinDistributionPolicy where a predicate requires the next word to be ROUND_ROBIN
tableRoundRobinDistributionPolicy_ROUND_ROBIN
    : 'ROUND_ROBIN'
    ;

// TSql160.g:26948
// tableHashDistributionPolicy where a predicate requires the next word to be HASH
tableHashDistributionPolicy_HASH
    : 'HASH' LeftParenthesis identifier ( Comma identifier )* RightParenthesis
    ;

// TSql160.g:27013
// tableNonClusteredIndexType where a predicate requires the next word to be HEAP
tableNonClusteredIndexType_HEAP
    : 'HEAP'
    ;

// TSql160.g:27717
// externalTableDataSource where a predicate requires the next word to be DATA_SOURCE
externalTableDataSource_DATA_SOURCE
    : 'DATA_SOURCE' EqualsSign identifier
    ;

// TSql160.g:27803
// externalTableRejectTypeOption where a predicate requires the next word to be REJECT_TYPE
externalTableRejectTypeOption_REJECT_TYPE
    : 'REJECT_TYPE' EqualsSign ( ( 'VALUE' | 'PERCENTAGE' ) )
    ;

// TSql160.g:27832
// externalTableDistributionOption where a predicate requires the next word to be DISTRIBUTION
externalTableDistributionOption_DISTRIBUTION
    : 'DISTRIBUTION' EqualsSign ( externalTableDistributionPolicy )
    ;

// TSql160.g:27780
// externalTableNumericLiteralOption where a predicate requires the next word to be REJECT_SAMPLE_VALUE
externalTableNumericLiteralOption_REJECT_SAMPLE_VALUE
    : 'REJECT_SAMPLE_VALUE' EqualsSign ( integerOrNumeric )
    ;

// TSql160.g:27780
// externalTableNumericLiteralOption where a predicate requires the next word to be REJECT_VALUE
externalTableNumericLiteralOption_REJECT_VALUE
    : 'REJECT_VALUE' EqualsSign ( integerOrNumeric )
    ;

// TSql160.g:27863
// externalTableReplicatedDistributionPolicy where a predicate requires the next word to be REPLICATED
externalTableReplicatedDistributionPolicy_REPLICATED
    : 'REPLICATED'
    ;

// TSql160.g:27871
// externalTableRoundRobinDistributionPolicy where a predicate requires the next word to be ROUND_ROBIN
externalTableRoundRobinDistributionPolicy_ROUND_ROBIN
    : 'ROUND_ROBIN'
    ;

// TSql160.g:27879
// externalTableShardedDistributionPolicy where a predicate requires the next word to be SHARDED
externalTableShardedDistributionPolicy_SHARDED
    : 'SHARDED' LeftParenthesis identifier RightParenthesis
    ;

// TSql160.g:28521
// alterTableSwitchStatement where a predicate requires the next word to be SWITCH
// TSql160.g:28542: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Partition)}?
alterTableSwitchStatement_SWITCH
    : 'SWITCH' ( 'PARTITION' expression )? To schemaObjectThreePartName ( 'PARTITION' expression )? ( With LeftParenthesis tableSwitchOptionList RightParenthesis )?
    ;

// TSql160.g:27957
// alterTableRebuildStatement where a predicate requires the next word to be REBUILD
alterTableRebuildStatement_REBUILD
    : 'REBUILD' ( partitionSpecifier )? ( With indexOptionList )?
    ;

// TSql160.g:28091
// lockEscalationTableOption where a predicate requires the next word to be LOCK_ESCALATION
lockEscalationTableOption_LOCK_ESCALATION
    : 'LOCK_ESCALATION' EqualsSign ( Identifier | Table )
    ;

// TSql160.g:28113
// fileStreamOnTableOption where a predicate requires the next word to be FILESTREAM_ON
fileStreamOnTableOption_FILESTREAM_ON
    : 'FILESTREAM_ON' EqualsSign stringOrIdentifier
    ;

// TSql160.g:28143
// remoteDataArchiveAlterTableOption where a predicate requires the next word to be REMOTE_DATA_ARCHIVE
remoteDataArchiveAlterTableOption_REMOTE_DATA_ARCHIVE
    : 'REMOTE_DATA_ARCHIVE' EqualsSign remoteDtatArchiveTableOnOffOption ( LeftParenthesis remoteDataArchiveTableSettings RightParenthesis )
    ;

// TSql160.g:29325
// generatedAlwaysClauseInAlter where a predicate requires the next word to be GENERATED
// TSql160.g:29334: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Start)}?
// TSql160.g:29354: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.End)}?
generatedAlwaysClauseInAlter_GENERATED
    : 'GENERATED' 'ALWAYS' As Identifier ( ( 'START' ) | ( End ) )
    ;

// TSql160.g:29464
// maskedClause where a predicate requires the next word to be MASKED
maskedClause_MASKED
    : 'MASKED' With LeftParenthesis Function EqualsSign stringLiteral RightParenthesis
    ;

// TSql160.g:29376
// columnEncryptionClause where a predicate requires the next word to be ENCRYPTED
columnEncryptionClause_ENCRYPTED
    : 'ENCRYPTED' With LeftParenthesis columnEncryptionDefinitionParameter Comma columnEncryptionDefinitionParameter Comma columnEncryptionDefinitionParameter RightParenthesis
    ;

// TSql160.g:29131
// columnStorage where a predicate requires the next word to be SPARSE
columnStorage_SPARSE
    : 'SPARSE' ( ( 'FILESTREAM' | 'SPARSE' ) | For 'ALL_SPARSE_COLUMNS' |  )
    ;

// TSql160.g:29131
// columnStorage where a predicate requires the next word to be FILESTREAM
columnStorage_FILESTREAM
    : 'FILESTREAM' ( ( 'FILESTREAM' | 'SPARSE' ) | For 'ALL_SPARSE_COLUMNS' |  )
    ;

// TSql160.g:29131
// columnStorage where a predicate requires the next word to be COLUMN_SET
columnStorage_COLUMN_SET
    : 'COLUMN_SET' ( ( 'FILESTREAM' | 'SPARSE' ) | For 'ALL_SPARSE_COLUMNS' |  )
    ;

// TSql160.g:29250
// generatedAlwaysClause where a predicate requires the next word to be GENERATED
// TSql160.g:29259: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Start)}?
// TSql160.g:29291: predicate dropped: {NextTokenMatches(CodeGenerationSupporter.End)}?
generatedAlwaysClause_GENERATED
    : 'GENERATED' 'ALWAYS' As Identifier ( ( 'START' ) | ( End ) )
    ;

// TSql160.g:29098
// sparseAndFilestreamStorage where a predicate requires the next word to be SPARSE
sparseAndFilestreamStorage_SPARSE
    : ( 'SPARSE' ) ( 'FILESTREAM' | 'SPARSE' )
    ;

// TSql160.g:29098
// sparseAndFilestreamStorage where a predicate requires the next word to be FILESTREAM
sparseAndFilestreamStorage_FILESTREAM
    : ( 'FILESTREAM' ) ( 'FILESTREAM' | 'SPARSE' )
    ;

// TSql160.g:29067
// sparseStorage where a predicate requires the next word to be SPARSE
sparseStorage_SPARSE
    : 'SPARSE'
    ;

// TSql160.g:29077
// allSparseStorage where a predicate requires the next word to be COLUMN_SET
allSparseStorage_COLUMN_SET
    : 'COLUMN_SET' For 'ALL_SPARSE_COLUMNS'
    ;

// TSql160.g:29121
// filestreamStorage where a predicate requires the next word to be FILESTREAM
filestreamStorage_FILESTREAM
    : 'FILESTREAM'
    ;

// TSql160.g:29413
// columnEncryptionKeyNameParameter where a predicate requires the next word to be COLUMN_ENCRYPTION_KEY
columnEncryptionKeyNameParameter_COLUMN_ENCRYPTION_KEY
    : 'COLUMN_ENCRYPTION_KEY' EqualsSign identifier
    ;

// TSql160.g:29426
// columnEncryptionTypeParameter where a predicate requires the next word to be ENCRYPTION_TYPE
columnEncryptionTypeParameter_ENCRYPTION_TYPE
    : 'ENCRYPTION_TYPE' EqualsSign ( 'DETERMINISTIC' | 'RANDOMIZED' )
    ;

// TSql160.g:29451
// columnEncryptionAlgorithmParameter where a predicate requires the next word to be ALGORITHM
columnEncryptionAlgorithmParameter_ALGORITHM
    : 'ALGORITHM' EqualsSign stringLiteral
    ;

// TSql160.g:30662
// graphMatchRecursivePattern where a predicate requires the next word to be SHORTEST_PATH
graphMatchRecursivePattern_SHORTEST_PATH
    : 'SHORTEST_PATH' LeftParenthesis ( graphMatchNodeExpression LeftParenthesis graphMatchSPRightRecursivePattern RightParenthesis graphRecursiveMatchQuantifier | LeftParenthesis graphMatchSPLeftRecursivePattern RightParenthesis graphRecursiveMatchQuantifier graphMatchNodeExpression ) RightParenthesis
    ;

// TSql160.g:33298
// identifier where a predicate requires the next word to be LAST_NODE
identifier_LAST_NODE
    : 'LAST_NODE'
    ;

// TSql160.g:30804
// graphMatchNodeExpression where a predicate requires the next word to be LAST_NODE
graphMatchNodeExpression_LAST_NODE
    : ( 'LAST_NODE' LeftParenthesis identifier RightParenthesis | identifier_LAST_NODE )
    ;

// TSql160.g:30748
// graphMatchLastNodePattern where a predicate requires the next word to be LAST_NODE
graphMatchLastNodePattern_LAST_NODE
    : graphMatchNodeExpression_LAST_NODE ( EqualsSign graphMatchNodeExpression | graphSingleMatchPattern ( graphSingleMatchPattern )* )
    ;

// TSql160.g:31744
// castCall where a predicate requires the next word to be CAST
castCall_CAST
    : 'CAST' LeftParenthesis expression As scalarDataType RightParenthesis
    ;

// TSql160.g:31760
// tryCastCall where a predicate requires the next word to be TRY_CAST
tryCastCall_TRY_CAST
    : 'TRY_CAST' LeftParenthesis expression As scalarDataType RightParenthesis
    ;

// TSql160.g:31684
// parseCall where a predicate requires the next word to be PARSE
parseCall_PARSE
    : 'PARSE' LeftParenthesis expression As scalarDataType ( 'USING' expression )? RightParenthesis
    ;

// TSql160.g:31714
// tryParseCall where a predicate requires the next word to be TRY_PARSE
tryParseCall_TRY_PARSE
    : 'TRY_PARSE' LeftParenthesis expression As scalarDataType ( 'USING' expression )? RightParenthesis
    ;

// TSql160.g:32536
// iIfCall where a predicate requires the next word to be IIF
iIfCall_IIF
    : 'IIF' LeftParenthesis booleanExpression Comma expression Comma expression RightParenthesis
    ;

// TSql160.g:32569
// jsonObjectCall where a predicate requires the next word to be JSON_OBJECT
jsonObjectCall_JSON_OBJECT
    : 'JSON_OBJECT' LeftParenthesis jsonObjectBuiltInFunctionCall
    ;

// TSql160.g:32582
// jsonArrayCall where a predicate requires the next word to be JSON_ARRAY
jsonArrayCall_JSON_ARRAY
    : 'JSON_ARRAY' LeftParenthesis jsonArrayBuiltInFunctionCall
    ;

// TSql160.g:32595
// trimCall where a predicate requires the next word to be TRIM
// TSql160.g:32606: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Leading) || NextTokenMatches(CodeGenerationSupporter.Trailing) || NextTokenMatches(CodeGenerationSupporter.Both)}?
trimCall_TRIM
    : 'TRIM' LeftParenthesis ( ( trim3ArgsBuiltInFunctionCall_LEADING | trim3ArgsBuiltInFunctionCall_TRAILING | trim3ArgsBuiltInFunctionCall_BOTH ) | trimBuiltInFunctionCall )
    ;

// TSql160.g:32811
// nextValueForCall where a predicate requires the next word to be NEXT
nextValueForCall_NEXT
    : 'NEXT' 'VALUE' For schemaObjectThreePartName ( overClause )?
    ;

// TSql160.g:18598
// graphPathClause where a predicate requires the next word to be GRAPH
graphPathClause_GRAPH
    : ( 'GRAPH' 'PATH' )
    ;

// TSql160.g:32427
// trim3ArgsBuiltInFunctionCall where a predicate requires the next word to be LEADING
trim3ArgsBuiltInFunctionCall_LEADING
    : 'LEADING' expression From expression RightParenthesis
    ;

// TSql160.g:32427
// trim3ArgsBuiltInFunctionCall where a predicate requires the next word to be TRAILING
trim3ArgsBuiltInFunctionCall_TRAILING
    : 'TRAILING' expression From expression RightParenthesis
    ;

// TSql160.g:32427
// trim3ArgsBuiltInFunctionCall where a predicate requires the next word to be BOTH
trim3ArgsBuiltInFunctionCall_BOTH
    : 'BOTH' expression From expression RightParenthesis
    ;

// Keywords: the tokens { } blocks, which ScriptDOM's lexer matches an
// identifier's text against.

// TSqlTokenTypes.g:36
Add
    : 'ADD'
    ;

// TSqlTokenTypes.g:37
All
    : 'ALL'
    ;

// TSqlTokenTypes.g:38
Alter
    : 'ALTER'
    ;

// TSqlTokenTypes.g:39
And
    : 'AND'
    ;

// TSqlTokenTypes.g:40
Any
    : 'ANY'
    ;

// TSqlTokenTypes.g:41
As
    : 'AS'
    ;

// TSqlTokenTypes.g:42
Asc
    : 'ASC'
    ;

// TSqlTokenTypes.g:43
Authorization
    : 'AUTHORIZATION'
    ;

// TSqlTokenTypes.g:44
Backup
    : 'BACKUP'
    ;

// TSqlTokenTypes.g:45
Begin
    : 'BEGIN'
    ;

// TSqlTokenTypes.g:46
Between
    : 'BETWEEN'
    ;

// TSqlTokenTypes.g:47
Break
    : 'BREAK'
    ;

// TSqlTokenTypes.g:48
Browse
    : 'BROWSE'
    ;

// TSqlTokenTypes.g:49
Bulk
    : 'BULK'
    ;

// TSqlTokenTypes.g:50
By
    : 'BY'
    ;

// TSqlTokenTypes.g:51
Cascade
    : 'CASCADE'
    ;

// TSqlTokenTypes.g:52
Case
    : 'CASE'
    ;

// TSqlTokenTypes.g:53
Check
    : 'CHECK'
    ;

// TSqlTokenTypes.g:54
Checkpoint
    : 'CHECKPOINT'
    ;

// TSqlTokenTypes.g:55
Close
    : 'CLOSE'
    ;

// TSqlTokenTypes.g:56
Clustered
    : 'CLUSTERED'
    ;

// TSqlTokenTypes.g:57
Coalesce
    : 'COALESCE'
    ;

// TSqlTokenTypes.g:58
Collate
    : 'COLLATE'
    ;

// TSqlTokenTypes.g:59
Column
    : 'COLUMN'
    ;

// TSqlTokenTypes.g:60
Commit
    : 'COMMIT'
    ;

// TSqlTokenTypes.g:61
Compute
    : 'COMPUTE'
    ;

// TSqlTokenTypes.g:62
Constraint
    : 'CONSTRAINT'
    ;

// TSqlTokenTypes.g:63
Contains
    : 'CONTAINS'
    ;

// TSqlTokenTypes.g:64
ContainsTable
    : 'CONTAINSTABLE'
    ;

// TSqlTokenTypes.g:65
Continue
    : 'CONTINUE'
    ;

// TSqlTokenTypes.g:66
Convert
    : 'CONVERT'
    ;

// TSqlTokenTypes.g:67
Create
    : 'CREATE'
    ;

// TSqlTokenTypes.g:68
Cross
    : 'CROSS'
    ;

// TSqlTokenTypes.g:69
Current
    : 'CURRENT'
    ;

// TSqlTokenTypes.g:70
CurrentDate
    : 'CURRENT_DATE'
    ;

// TSqlTokenTypes.g:71
CurrentTime
    : 'CURRENT_TIME'
    ;

// TSqlTokenTypes.g:72
CurrentTimestamp
    : 'CURRENT_TIMESTAMP'
    ;

// TSqlTokenTypes.g:73
CurrentUser
    : 'CURRENT_USER'
    ;

// TSqlTokenTypes.g:74
Cursor
    : 'CURSOR'
    ;

// TSqlTokenTypes.g:75
Database
    : 'DATABASE'
    ;

// TSqlTokenTypes.g:76
Dbcc
    : 'DBCC'
    ;

// TSqlTokenTypes.g:77
Deallocate
    : 'DEALLOCATE'
    ;

// TSqlTokenTypes.g:78
Declare
    : 'DECLARE'
    ;

// TSqlTokenTypes.g:79
Default
    : 'DEFAULT'
    ;

// TSqlTokenTypes.g:80
Delete
    : 'DELETE'
    ;

// TSqlTokenTypes.g:81
Deny
    : 'DENY'
    ;

// TSqlTokenTypes.g:82
Desc
    : 'DESC'
    ;

// TSqlTokenTypes.g:83
Distinct
    : 'DISTINCT'
    ;

// TSqlTokenTypes.g:84
Distributed
    : 'DISTRIBUTED'
    ;

// TSqlTokenTypes.g:85
Double
    : 'DOUBLE'
    ;

// TSqlTokenTypes.g:86
Drop
    : 'DROP'
    ;

// TSqlTokenTypes.g:87
Else
    : 'ELSE'
    ;

// TSqlTokenTypes.g:88
End
    : 'END'
    ;

// TSqlTokenTypes.g:89
Errlvl
    : 'ERRLVL'
    ;

// TSqlTokenTypes.g:90
Escape
    : 'ESCAPE'
    ;

// TSqlTokenTypes.g:91
Except
    : 'EXCEPT'
    ;

// TSqlTokenTypes.g:92
Exec
    : 'EXEC'
    ;

// TSqlTokenTypes.g:93
Execute
    : 'EXECUTE'
    ;

// TSqlTokenTypes.g:94
Exists
    : 'EXISTS'
    ;

// TSqlTokenTypes.g:95
Exit
    : 'EXIT'
    ;

// TSqlTokenTypes.g:96
Fetch
    : 'FETCH'
    ;

// TSqlTokenTypes.g:97
File
    : 'FILE'
    ;

// TSqlTokenTypes.g:98
FillFactor
    : 'FILLFACTOR'
    ;

// TSqlTokenTypes.g:99
For
    : 'FOR'
    ;

// TSqlTokenTypes.g:100
Foreign
    : 'FOREIGN'
    ;

// TSqlTokenTypes.g:101
FreeText
    : 'FREETEXT'
    ;

// TSqlTokenTypes.g:102
FreeTextTable
    : 'FREETEXTTABLE'
    ;

// TSqlTokenTypes.g:103
From
    : 'FROM'
    ;

// TSqlTokenTypes.g:104
Full
    : 'FULL'
    ;

// TSqlTokenTypes.g:105
Function
    : 'FUNCTION'
    ;

// TSqlTokenTypes.g:106
GoTo
    : 'GOTO'
    ;

// TSqlTokenTypes.g:107
Grant
    : 'GRANT'
    ;

// TSqlTokenTypes.g:108
Group
    : 'GROUP'
    ;

// TSqlTokenTypes.g:109
Having
    : 'HAVING'
    ;

// TSqlTokenTypes.g:110
HoldLock
    : 'HOLDLOCK'
    ;

// TSqlTokenTypes.g:111
Identity
    : 'IDENTITY'
    ;

// TSqlTokenTypes.g:112
IdentityInsert
    : 'IDENTITY_INSERT'
    ;

// TSqlTokenTypes.g:113
IdentityColumn
    : 'IDENTITYCOL'
    ;

// TSqlTokenTypes.g:114
If
    : 'IF'
    ;

// TSqlTokenTypes.g:115
In
    : 'IN'
    ;

// TSqlTokenTypes.g:116
Index
    : 'INDEX'
    ;

// TSqlTokenTypes.g:117
Inner
    : 'INNER'
    ;

// TSqlTokenTypes.g:118
Insert
    : 'INSERT'
    ;

// TSqlTokenTypes.g:119
Intersect
    : 'INTERSECT'
    ;

// TSqlTokenTypes.g:120
Into
    : 'INTO'
    ;

// TSqlTokenTypes.g:121
Is
    : 'IS'
    ;

// TSqlTokenTypes.g:122
Join
    : 'JOIN'
    ;

// TSqlTokenTypes.g:123
Key
    : 'KEY'
    ;

// TSqlTokenTypes.g:124
Kill
    : 'KILL'
    ;

// TSqlTokenTypes.g:125
Left
    : 'LEFT'
    ;

// TSqlTokenTypes.g:126
Like
    : 'LIKE'
    ;

// TSqlTokenTypes.g:127
LineNo
    : 'LINENO'
    ;

// TSqlTokenTypes.g:128
National
    : 'NATIONAL'
    ;

// TSqlTokenTypes.g:129
NoCheck
    : 'NOCHECK'
    ;

// TSqlTokenTypes.g:130
NonClustered
    : 'NONCLUSTERED'
    ;

// TSqlTokenTypes.g:131
Not
    : 'NOT'
    ;

// TSqlTokenTypes.g:132
Null
    : 'NULL'
    ;

// TSqlTokenTypes.g:133
NullIf
    : 'NULLIF'
    ;

// TSqlTokenTypes.g:134
Of
    : 'OF'
    ;

// TSqlTokenTypes.g:135
Off
    : 'OFF'
    ;

// TSqlTokenTypes.g:136
Offsets
    : 'OFFSETS'
    ;

// TSqlTokenTypes.g:137
On
    : 'ON'
    ;

// TSqlTokenTypes.g:138
Open
    : 'OPEN'
    ;

// TSqlTokenTypes.g:139
OpenDataSource
    : 'OPENDATASOURCE'
    ;

// TSqlTokenTypes.g:140
OpenQuery
    : 'OPENQUERY'
    ;

// TSqlTokenTypes.g:141
OpenRowSet
    : 'OPENROWSET'
    ;

// TSqlTokenTypes.g:142
OpenXml
    : 'OPENXML'
    ;

// TSqlTokenTypes.g:143
Option
    : 'OPTION'
    ;

// TSqlTokenTypes.g:144
Or
    : 'OR'
    ;

// TSqlTokenTypes.g:145
Order
    : 'ORDER'
    ;

// TSqlTokenTypes.g:146
Outer
    : 'OUTER'
    ;

// TSqlTokenTypes.g:147
Over
    : 'OVER'
    ;

// TSqlTokenTypes.g:148
Percent
    : 'PERCENT'
    ;

// TSqlTokenTypes.g:149
Plan
    : 'PLAN'
    ;

// TSqlTokenTypes.g:150
Primary
    : 'PRIMARY'
    ;

// TSqlTokenTypes.g:151
Print
    : 'PRINT'
    ;

// TSqlTokenTypes.g:152
Proc
    : 'PROC'
    ;

// TSqlTokenTypes.g:153
Procedure
    : 'PROCEDURE'
    ;

// TSqlTokenTypes.g:154
Public
    : 'PUBLIC'
    ;

// TSqlTokenTypes.g:155
Raiserror
    : 'RAISERROR'
    ;

// TSqlTokenTypes.g:156
Read
    : 'READ'
    ;

// TSqlTokenTypes.g:157
ReadText
    : 'READTEXT'
    ;

// TSqlTokenTypes.g:158
Reconfigure
    : 'RECONFIGURE'
    ;

// TSqlTokenTypes.g:159
References
    : 'REFERENCES'
    ;

// TSqlTokenTypes.g:160
Replication
    : 'REPLICATION'
    ;

// TSqlTokenTypes.g:161
Restore
    : 'RESTORE'
    ;

// TSqlTokenTypes.g:162
Restrict
    : 'RESTRICT'
    ;

// TSqlTokenTypes.g:163
Return
    : 'RETURN'
    ;

// TSqlTokenTypes.g:164
Revoke
    : 'REVOKE'
    ;

// TSqlTokenTypes.g:165
Right
    : 'RIGHT'
    ;

// TSqlTokenTypes.g:166
Rollback
    : 'ROLLBACK'
    ;

// TSqlTokenTypes.g:167
RowCount
    : 'ROWCOUNT'
    ;

// TSqlTokenTypes.g:168
RowGuidColumn
    : 'ROWGUIDCOL'
    ;

// TSqlTokenTypes.g:169
Rule
    : 'RULE'
    ;

// TSqlTokenTypes.g:170
Save
    : 'SAVE'
    ;

// TSqlTokenTypes.g:171
Schema
    : 'SCHEMA'
    ;

// TSqlTokenTypes.g:172
Select
    : 'SELECT'
    ;

// TSqlTokenTypes.g:173
SessionUser
    : 'SESSION_USER'
    ;

// TSqlTokenTypes.g:174
Set
    : 'SET'
    ;

// TSqlTokenTypes.g:175
SetUser
    : 'SETUSER'
    ;

// TSqlTokenTypes.g:176
Shutdown
    : 'SHUTDOWN'
    ;

// TSqlTokenTypes.g:177
Some
    : 'SOME'
    ;

// TSqlTokenTypes.g:178
Statistics
    : 'STATISTICS'
    ;

// TSqlTokenTypes.g:179
SystemUser
    : 'SYSTEM_USER'
    ;

// TSqlTokenTypes.g:180
Table
    : 'TABLE'
    ;

// TSqlTokenTypes.g:181
TextSize
    : 'TEXTSIZE'
    ;

// TSqlTokenTypes.g:182
Then
    : 'THEN'
    ;

// TSqlTokenTypes.g:183
To
    : 'TO'
    ;

// TSqlTokenTypes.g:184
Top
    : 'TOP'
    ;

// TSqlTokenTypes.g:185
Tran
    : 'TRAN'
    ;

// TSqlTokenTypes.g:186
Transaction
    : 'TRANSACTION'
    ;

// TSqlTokenTypes.g:187
Trigger
    : 'TRIGGER'
    ;

// TSqlTokenTypes.g:188
Truncate
    : 'TRUNCATE'
    ;

// TSqlTokenTypes.g:189
TSEqual
    : 'TSEQUAL'
    ;

// TSqlTokenTypes.g:190
Union
    : 'UNION'
    ;

// TSqlTokenTypes.g:191
Unique
    : 'UNIQUE'
    ;

// TSqlTokenTypes.g:192
Update
    : 'UPDATE'
    ;

// TSqlTokenTypes.g:193
UpdateText
    : 'UPDATETEXT'
    ;

// TSqlTokenTypes.g:194
Use
    : 'USE'
    ;

// TSqlTokenTypes.g:195
User
    : 'USER'
    ;

// TSqlTokenTypes.g:196
Values
    : 'VALUES'
    ;

// TSqlTokenTypes.g:197
Varying
    : 'VARYING'
    ;

// TSqlTokenTypes.g:198
View
    : 'VIEW'
    ;

// TSqlTokenTypes.g:199
WaitFor
    : 'WAITFOR'
    ;

// TSqlTokenTypes.g:200
When
    : 'WHEN'
    ;

// TSqlTokenTypes.g:201
Where
    : 'WHERE'
    ;

// TSqlTokenTypes.g:202
While
    : 'WHILE'
    ;

// TSqlTokenTypes.g:203
With
    : 'WITH'
    ;

// TSqlTokenTypes.g:204
WriteText
    : 'WRITETEXT'
    ;

// TSql160.g:33926
External
    : 'EXTERNAL'
    ;

// TSql160.g:33927
Pivot
    : 'PIVOT'
    ;

// TSql160.g:33928
Revert
    : 'REVERT'
    ;

// TSql160.g:33929
TableSample
    : 'TABLESAMPLE'
    ;

// TSql160.g:33930
Unpivot
    : 'UNPIVOT'
    ;

// TSql160.g:33933
Merge
    : 'MERGE'
    ;

// TSql160.g:33934
StopList
    : 'STOPLIST'
    ;

// TSql160.g:33937
SemanticKeyPhraseTable
    : 'SEMANTICKEYPHRASETABLE'
    ;

// TSql160.g:33938
SemanticSimilarityTable
    : 'SEMANTICSIMILARITYTABLE'
    ;

// TSql160.g:33939
SemanticSimilarityDetailsTable
    : 'SEMANTICSIMILARITYDETAILSTABLE'
    ;

// TSql160.g:33940
TryConvert
    : 'TRY_CONVERT'
    ;

// TSql160.g:33950
Bang
    : '!'
    ;

// TSql160.g:33951
PercentSign
    : '%'
    ;

// TSql160.g:33952
Ampersand
    : '&'
    ;

// TSql160.g:33953
LeftParenthesis
    : '('
    ;

// TSql160.g:33954
RightParenthesis
    : ')'
    ;

// TSql160.g:33955
LeftCurly
    : '{'
    ;

// TSql160.g:33956
RightCurly
    : '}'
    ;

// TSql160.g:33957
Star
    : '*'
    ;

// TSql160.g:33958
MultiplyEquals
    : '*='
    ;

// TSql160.g:33959
Plus
    : '+'
    ;

// TSql160.g:33960
Comma
    : ','
    ;

// TSql160.g:33961
Minus
    : '-'
    ;

// TSql160.g:33964
Divide
    : '/'
    ;

// TSql160.g:33965
Colon
    : ':'
    ;

// TSql160.g:33966
DoubleColon
    : '::'
    ;

// TSql160.g:33967
LessThan
    : '<'
    ;

// TSql160.g:33968
EqualsSign
    : '='
    ;

// TSql160.g:33969
RightOuterJoin
    : '=*'
    ;

// TSql160.g:33970
GreaterThan
    : '>'
    ;

// TSql160.g:33971
Circumflex
    : '^'
    ;

// TSql160.g:33972
VerticalLine
    : '|'
    ;

// TSql160.g:33973
Tilde
    : '~'
    ;

// TSql160.g:33974
AddEquals
    : '+='
    ;

// TSql160.g:33975
SubtractEquals
    : '-='
    ;

// TSql160.g:33976
DivideEquals
    : '/='
    ;

// TSql160.g:33977
ModEquals
    : '%='
    ;

// TSql160.g:33978
BitwiseAndEquals
    : '&='
    ;

// TSql160.g:33979
BitwiseOrEquals
    : '|='
    ;

// TSql160.g:33980
BitwiseXorEquals
    : '^='
    ;

// TSql160.g:33981
LeftShift
    : '<<'
    ;

// TSql160.g:33982
RightShift
    : '>>'
    ;

// TSql160.g:33983
Concat
    : '||'
    ;

// TSql160.g:33984
ConcatEquals
    : '||='
    ;

// TSql160.g:33987
fragment Semicolon
    : ';'
    ;

// TSql160.g:33990
fragment Digit
    : '0'..'9'
    ;

// TSql160.g:33995
fragment FirstLetter
    : 'a'..'z'
    | '_'
    | '#'
    | CharHighNotWhitespace
    ;

// TSql160.g:34003
fragment Letter
    : 'a'..'z'
    | '_'
    | '#'
    | '@'
    | '$'
    | CharHighNotWhitespace
    ;

// TSql160.g:34013
fragment CharHighNotWhitespace
    : '\u0080'..'\u0084'
    | '\u0086'..'\u009F'
    | '\u00A1'..'\u1679'
    | '\u1681'..'\u1FFF'
    | '\u200C'..'\u2027'
    | '\u202A'..'\u202E'
    | '\u2030'..'\u205E'
    | '\u2060'..'\u2FFF'
    | '\u3001'..'\uFFFE'
    ;

// TSql160.g:34043
fragment MoneySign
    : '$'
    | '\u00A3'
    | '\u00A4'
    | '\u00A5'
    | '\u09F2'
    | '\u09F3'
    | '\u0E3F'
    | '\u20AC'
    | '\u20A1'
    | '\u20A2'
    | '\u20A3'
    | '\u20A4'
    | '\u20A6'
    | '\u20A7'
    | '\u20A8'
    | '\u20A9'
    | '\u20AA'
    | '\u20AB'
    ;

// TSql160.g:34064
// TSql160.g:34070: this alternative may emit Semicolon instead ($setType)
ProcNameSemicolon
    : Semicolon
    | Semicolon
    ;

// TSql160.g:34076
fragment WS_CHAR_WO_NEWLINE
    : '\u0000'..'\b'
    | '\t'
    | '\u000B'
    | '\f'
    | '\u000E'..'\u001F'
    | ' '
    | '\u0085'
    | '\u00A0'
    | '\u1680'
    | '\u2000'..'\u200A'
    | '\u200B'
    | '\u2028'
    | '\u2029'
    | '\u202F'
    | '\u205F'
    | '\u3000'
    ;

// TSql160.g:34105
WhiteSpace
    : ( WS_CHAR_WO_NEWLINE )+
    | EndOfLine
      -> channel(HIDDEN)
    ;

// TSql160.g:34125
// TSql160.g:34127: predicate dropped: {CurrentOffset==_acceptableGoOffset}?
// TSql160.g:34132: this alternative may emit Identifier instead ($setType)
// TSql160.g:34135: this alternative may emit Identifier instead ($setType)
// TSql160.g:34140: this alternative may emit Label instead ($setType)
Go
    : 'go' ( ( Letter | Digit )* ( Colon )? )
    ;

// TSql160.g:34163
// TSql160.g:34165: this alternative may emit Numeric instead ($setType)
// TSql160.g:34165: this alternative may emit Integer instead ($setType)
// TSql160.g:34177: this alternative may emit Numeric instead ($setType)
// TSql160.g:34183: this alternative may emit Real instead ($setType)
// TSql160.g:34188: this alternative may emit Dot instead ($setType)
// TSql160.g:34192: this alternative may emit Numeric instead ($setType)
// TSql160.g:34196: this alternative may emit Real instead ($setType)
// TSql160.g:34206: this alternative may emit HexLiteral instead ($setType)
Number
    : ( Digit )+ ( '.' ( Digit )* )? ( Exponent )?
    | '.' ( ( Digit )+ ( Exponent )? )?
    | '0x' ( Digit | 'a'..'f' | '\\' EndOfLine )*
    ;

// TSql160.g:34212
fragment Exponent
    : 'e' ( '+' | '-' )? ( Digit )*
    ;

// TSql160.g:34217
fragment EndOfLine
    : ( '\r\n' | '\r' | '\n' )
    ;

// TSql160.g:34251
AsciiStringLiteral
    : '\'' ( ~( '\'' | '\n' | '\r' ) | EndOfLine | '\'' '\'' )* '\''
    ;

// TSql160.g:34264
UnicodeStringLiteral
    : 'n' '\'' ( ~( '\'' | '\n' | '\r' ) | EndOfLine | '\'' '\'' )* '\''
    ;

// TSql160.g:34277
// TSql160.g:34289: this alternative may emit SqlCommandIdentifier instead ($setType)
// TSql160.g:34295: this alternative may emit DollarPartition instead ($setType)
// TSql160.g:34295: this alternative may emit PseudoColumn instead ($setType)
// TSql160.g:34308: this alternative may emit Money instead ($setType)
// TSql160.g:34316: this alternative may emit Label instead ($setType)
Identifier
    : '$' LeftParenthesis ( ~( ')' | '\n' | '\r' ) | EndOfLine )+ RightParenthesis
    | '$' ( '@' | FirstLetter ) ( Letter | Digit )*
    | MoneySign ( ' ' )* ( Minus | Plus )? ( Digit )+ ( '.' ( Digit )* ( Exponent )? | Exponent )?
    | FirstLetter ( Letter | Digit )* ( Colon |  )
    ;

// TSql160.g:34324
// TSql160.g:34346: this alternative may emit AsciiStringOrQuotedIdentifier instead ($setType)
// TSql160.g:34346: this alternative may emit AsciiStringLiteral instead ($setType)
QuotedIdentifier
    : '[' ( ~( ']' | '\n' | '\r' ) | EndOfLine | ']' ']' )+ ']'
    | '"' ( ( ~( '"' | '\n' | '\r' ) | EndOfLine | '"' '"' )* ) '"'
    ;

// TSql160.g:34360
Variable
    : '@' ( Letter | Digit )*
    ;

// TSql160.g:34367
// TSql160.g:34372: this alternative may emit OdbcInitiator instead ($setType)
// TSql160.g:34377: predicate dropped: {LA(1) != EOF_CHAR}?
SingleLineComment
    : '--' ( '(' '*' | ( ~( '\r' | '\n' ) )* )
      -> channel(HIDDEN)
    ;

// TSql160.g:34386
// TSql160.g:34396: predicate dropped: {LA(2)!='/'}?
// TSql160.g:34397: predicate dropped: {LA(2)!='*'}?
MultilineComment
    : '/*' ( EndOfLine | '*' | '/' | ~( '*' | '\n' | '\r' | '/' ) | MultilineComment )* '*/'
      -> channel(HIDDEN)
    ;

// Tokens with no rule of their own: an action of the rule named sets
// their type ($setType), and tokens.tsv samples them.
// Dot: TSql160.g:33963
// Label: TSql160.g:34148
// Integer: TSql160.g:34151
// Real: TSql160.g:34154
// Numeric: TSql160.g:34157
// HexLiteral: TSql160.g:34160
// Money: TSql160.g:34237
// SqlCommandIdentifier: TSql160.g:34240
// PseudoColumn: TSql160.g:34243
// DollarPartition: TSql160.g:34246
// AsciiStringOrQuotedIdentifier: TSql160.g:34249
// OdbcInitiator: TSql160.g:34365
