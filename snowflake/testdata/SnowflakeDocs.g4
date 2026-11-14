// Snowflake's SQL command reference as an ANTLR 4 grammar: the Syntax
// sections of the pages https://docs.snowflake.com/en/sql-reference/sql-all
// lists, converted by snowflake/parser/cmd/docsgrammar. Do not edit it:
// correct a page's text in docs-fixes.tsv and run the converter again.
//
// Keywords are literals. ID, QNAME, STRING, NUMBER, EXPR, QUERY,
// COPY_STATEMENT, DATA_TYPE, STAGE, FILE, POSITION, VARIABLE, OBJECT_TYPE and
// OBJECT_TYPES are the values the docs write as <placeholders>; tokens.tsv
// samples them. A
// form a page writes with ... for text it leaves out, or with a word it
// never defines, is not converted and is listed under its page.
//
// pin: default https://docs.snowflake.com/en/sql-reference/sql-all sha256:334be655f4130196272f45c96ed77fd2bc1399b40ce5c958c53bea02f6ad9177 2026-11-13

grammar SnowflakeDocs;

statement
    : alterAccount
    | alterAccountSetIngressPolicy
    | alterAgent
    | alterAggregationPolicy
    | alterAlert
    | alterApiIntegration
    | alterApplication
    | alterApplicationDropAppSpec
    | alterApplicationDropConfigurationDefinition
    | alterApplicationPackage
    | alterApplicationPackageReleaseChannel
    | alterApplicationPackageReleaseDirective
    | alterApplicationPackageVersion
    | alterApplicationRole
    | alterApplicationSequenceNumber
    | alterApplicationService
    | createApplicationService
    | alterApplicationSetAppSpec
    | alterApplicationSetConfigurationDefinition
    | alterApplicationSetConfigurationValue
    | alterApplicationUnsetConfiguration
    | alterArtifactRepository
    | createArtifactRepository
    | alterAuthenticationPolicy
    | alterBackupPolicy
    | alterBackupSet
    | alterCatalogIntegration
    | alterComputePool
    | alterConnection
    | alterContact
    | alterCortexSearch
    | alterDataMovementPolicy
    | alterDataMovementRule
    | alterDatabase
    | alterDatabaseCatalogLinked
    | alterDatabaseRole
    | alterDatasetAddVersion
    | alterDatasetDropVersion
    | alterDbtProject
    | alterDcmProject
    | alterDynamicTable
    | alterEventRoutingTable
    | alterExperiment
    | alterExternalAgent
    | alterExternalConsumer
    | alterExternalConsumerAddProgrammaticAccessToken
    | alterExternalConsumerModifyProgrammaticAccessToken
    | alterExternalConsumerRemoveProgrammaticAccessToken
    | alterExternalConsumerRotateProgrammaticAccessToken
    | alterExternalAccessIntegration
    | alterExternalTable
    | alterExternalVolume
    | alterFailoverGroup
    | alterFeaturePolicy
    | alterFileFormat
    | alterFunction
    | alterFunctionDmf
    | alterFunctionSpcs
    | alterGateway
    | alterGitRepository
    | alterIcebergTable
    | alterIcebergTableAlterColumnSetDataType
    | alterIcebergTableConvertToManaged
    | alterIcebergTableOptimize
    | alterIcebergTablePartitionEvolution
    | alterIcebergTableRefresh
    | alterIcebergTableVersionTag
    | alterImageRepository
    | alterIngressPolicy
    | alterJoinPolicy
    | alterListing
    | alterMaintenancePolicy
    | alterMaskingPolicy
    | alterMaterializedView
    | alterModel
    | alterModelAddVersion
    | alterModelDropVersion
    | alterModelModifyVersion
    | alterModelMonitor
    | alterMultiPartyApprovalPolicy
    | alterNetworkPolicy
    | alterNetworkRule
    | alterNotebook
    | alterNotificationIntegrationEmail
    | alterNotificationIntegrationQueueInboundAzure
    | alterNotificationIntegrationQueueInboundGcp
    | alterNotificationIntegrationMultiQueue
    | alterNotificationIntegrationQueueOutboundAws
    | alterNotificationIntegrationQueueOutboundAzure
    | alterNotificationIntegrationQueueOutboundGcp
    | alterNotificationIntegrationWebhooks
    | alterOflowDataPlane
    | alterOnlineFeatureTable
    | createOnlineFeatureTable
    | alterOpenflowDeployment
    | alterOpenflowRuntime
    | alterOrganization
    | alterOrganizationAccount
    | alterOrganizationProfile
    | alterOrganizationSetEventRoutingTable
    | alterOrganizationUnsetEventRoutingTable
    | alterOrganizationUser
    | alterOrganizationUserGroup
    | alterPackagesPolicy
    | alterPasswordPolicy
    | alterPipe
    | alterPostgresInstance
    | alterPrivacyPolicy
    | alterProcedure
    | alterProjectionPolicy
    | alterReplicationGroup
    | alterResourceMonitor
    | alterRestrictedSessionScope
    | alterRole
    | alterRowAccessPolicy
    | alterSchema
    | alterSecret
    | alterSecurityIntegration
    | alterSecurityIntegrationApiAuth
    | alterSecurityIntegrationApiAuthExternalSecretProvider
    | alterSecurityIntegrationAwsIam
    | alterSecurityIntegrationOauthExternal
    | alterSecurityIntegrationOauthSnowflake
    | alterSecurityIntegrationOidc
    | alterSecurityIntegrationSaml2
    | alterSecurityIntegrationScim
    | alterSemanticView
    | alterSequence
    | alterService
    | alterSession
    | alterSessionPolicy
    | alterShare
    | alterSnapshot
    | alterSnapshotPolicy
    | alterSnapshotSet
    | alterStage
    | alterStorageIntegration
    | alterStorageLifecyclePolicy
    | alterStream
    | alterStreamlit
    | alterTable
    | alterTableColumn
    | alterTableEventTable
    | alterTag
    | alterTask
    | alterTaskGraph
    | alterType
    | alterUser
    | alterUserAddKeyPair
    | alterUserAddProgrammaticAccessToken
    | alterUserAddWorkloadIdentity
    | alterUserModifyKeyPair
    | alterUserModifyProgrammaticAccessToken
    | alterUserModifyWorkloadIdentity
    | alterUserRemoveKeyPair
    | alterUserRemoveProgrammaticAccessToken
    | alterUserRemoveWorkloadIdentity
    | alterUserRotateKeyPair
    | alterUserRotateProgrammaticAccessToken
    | alterUserSetIngressPolicy
    | alterVariable
    | alterView
    | alterWarehouse
    | begin
    | call
    | callWith
    | comment
    | commit
    | copyFiles
    | copyIntoLocation
    | copyIntoTable
    | createAccount
    | createAgent
    | createAggregationPolicy
    | createAiFunction
    | createAlert
    | createApiIntegration
    | createApplication
    | createApplicationPackage
    | createApplicationRole
    | createAuthenticationPolicy
    | createBackupPolicy
    | createBackupSet
    | createCatalogIntegrationDeltaSharing
    | createIcebergTableDelta
    | createCatalogIntegrationGlue
    | createCatalogIntegrationObjectStorage
    | createCatalogIntegrationHorizon
    | createCatalogIntegrationOpenCatalog
    | createCatalogIntegrationRest
    | createCatalogIntegrationSnowflakePostgres
    | createClone
    | createComputePool
    | createConnection
    | createContact
    | createCortexSearch
    | createDataMetricFunction
    | createDataMovementPolicy
    | createDataMovementRule
    | createDatabase
    | createDatabaseCatalogLinked
    | createDatabaseRole
    | createDataset
    | createDbtProject
    | createDcmProject
    | createDynamicTable
    | createEventTable
    | createExperiment
    | createExternalAgent
    | createExternalConsumer
    | createExternalAccessIntegration
    | createExternalFunction
    | createExternalTable
    | createFailoverGroup
    | createFeaturePolicy
    | createFileFormat
    | createFunction
    | createFunctionSpcs
    | createGateway
    | createGitRepository
    | createHybridTable
    | createIcebergTable
    | createIcebergTableAwsGlue
    | createIcebergTableIcebergFiles
    | createIcebergTableParquet
    | createIcebergTableRest
    | createIcebergTableSnowflake
    | createImageRepository
    | createIndex
    | createIngressPolicy
    | createInteractiveTable
    | createInteractiveWarehouse
    | createJoinPolicy
    | createListing
    | createMaintenancePolicy
    | createManagedAccount
    | createMaskingPolicy
    | createMaterializedView
    | createMcpServer
    | createModel
    | createModelMonitor
    | createMultiPartyApprovalPolicy
    | createNetworkPolicy
    | createNetworkRule
    | createNotebook
    | createNotebookProject
    | createCodeBundle
    | createNotificationIntegrationMultiQueue
    | createNotificationIntegrationQueueInboundAzure
    | createNotificationIntegrationQueueInboundGcp
    | createNotificationIntegrationQueueOutboundAws
    | createNotificationIntegrationQueueOutboundAzure
    | createNotificationIntegrationQueueOutboundGcp
    | createNotificationIntegrationWebhooks
    | createOpenflowConnector
    | createOpenflowDeployment
    | createOpenflowRuntime
    | createOrganizationAccount
    | createOrganizationListing
    | createOrganizationProfile
    | createOrganizationUser
    | createOrganizationUserGroup
    | createPackagesPolicy
    | createPasswordPolicy
    | createPipe
    | createPostgresInstance
    | createPrivacyPolicy
    | createProcedure
    | createProjectionPolicy
    | createProvisionedThroughput
    | createReplicationGroup
    | createResourceMonitor
    | createRestrictedSessionScope
    | createRole
    | createRowAccessPolicy
    | createSchema
    | createSecret
    | createSecurityIntegrationApiAuth
    | createSecurityIntegrationApiAuthExternalSecretProvider
    | createSecurityIntegrationAwsIam
    | createSecurityIntegrationOauthExternal
    | createSecurityIntegrationOauthSnowflake
    | createSecurityIntegrationOidc
    | createSecurityIntegrationSaml2
    | createSecurityIntegrationScim
    | createSemanticView
    | createSequence
    | createService
    | createSessionPolicy
    | createShare
    | createSnapshot
    | createSnapshotPolicy
    | createSnapshotSet
    | createStage
    | createStorageIntegration
    | createStorageIntegrationPostgresInternal
    | createStorageLifecyclePolicy
    | createStream
    | createStreamlit
    | createTable
    | createTableConstraint
    | createTag
    | createTask
    | createTaskGraph
    | createType
    | createUser
    | createVariable
    | createVersionedSchema
    | createView
    | createWarehouse
    | delete
    | descAgent
    | descAggregationPolicy
    | descAlert
    | descApplication
    | descApplicationPackage
    | descApplicationService
    | descArtifactRepository
    | descAuthenticationPolicy
    | descAvailableListing
    | descAvailableOrganizationProfile
    | descBackupPolicy
    | descBackupSet
    | descCatalogIntegration
    | descCommittedTask
    | descComputePool
    | descConfiguration
    | descCortexSearch
    | descDataMovementPolicy
    | descDataMovementRule
    | descDatabase
    | descDbtProject
    | descDcmProject
    | descDynamicTable
    | descEventTable
    | descExternalAgent
    | descExternalConsumer
    | descExternalTable
    | descExternalVolume
    | descFeaturePolicy
    | descFileFormat
    | descFunction
    | descFunctionDmf
    | descFunctionSpcs
    | descGateway
    | descGitRepository
    | descIcebergTable
    | descIngressPolicy
    | descIntegration
    | descJoinPolicy
    | descListing
    | descMaintenancePolicy
    | descMaskingPolicy
    | descMaterializedView
    | descMcpServer
    | descModelMonitor
    | descMultiPartyApprovalPolicy
    | descNetworkPolicy
    | descNetworkRule
    | descNotebook
    | descNotificationIntegration
    | descOflowDataPlaneIntegration
    | descOnlineFeatureTable
    | descOpenflowConnector
    | descOpenflowDeployment
    | descOpenflowRuntime
    | descOrganizationProfile
    | descPackagesPolicy
    | descPasswordPolicy
    | descPipe
    | descPostgresInstance
    | descPrivacyPolicy
    | descProcedure
    | descProjectionPolicy
    | descResult
    | descRestrictedSessionScope
    | descRowAccessPolicy
    | descSchema
    | descSearchOptimization
    | descSecret
    | descSemanticView
    | descSequence
    | descService
    | descSessionPolicy
    | descShare
    | descSnapshot
    | descSnapshotPolicy
    | descSnapshotSet
    | descSpecification
    | descStage
    | descStorageLifecyclePolicy
    | descStream
    | descStreamlit
    | descTable
    | descTask
    | descTaskGraph
    | descTransaction
    | descType
    | descUser
    | descVariable
    | descView
    | descWarehouse
    | dropAccount
    | dropAgent
    | dropAggregationPolicy
    | dropAlert
    | dropApplication
    | dropApplicationPackage
    | dropApplicationRole
    | dropApplicationService
    | dropArtifactRepository
    | dropAuthenticationPolicy
    | dropBackupPolicy
    | dropBackupSet
    | dropCatalogIntegration
    | dropComputePool
    | dropConnection
    | dropContact
    | dropCortexSearch
    | dropDataMovementPolicy
    | dropDataMovementRule
    | dropDatabase
    | dropDatabaseRole
    | dropDbtProject
    | dropDcmProject
    | dropDynamicTable
    | dropEventRoutingTable
    | dropExperiment
    | dropExternalAgent
    | dropExternalConsumer
    | dropExternalTable
    | dropExternalVolume
    | undropExternalVolume
    | dropFailoverGroup
    | dropFeaturePolicy
    | dropFileFormat
    | dropFunction
    | dropFunctionDmf
    | dropFunctionSpcs
    | dropGateway
    | dropGitRepository
    | dropIcebergTable
    | undropIcebergTable
    | dropImageRepository
    | dropIndex
    | dropIngressPolicy
    | dropIntegration
    | dropJoinPolicy
    | dropListing
    | dropMaintenancePolicy
    | dropManagedAccount
    | dropMaskingPolicy
    | dropMaterializedView
    | dropMcpServer
    | dropModel
    | dropModelMonitor
    | dropMultiPartyApprovalPolicy
    | dropNetworkPolicy
    | dropNetworkRule
    | dropNotebook
    | undropNotebook
    | dropOnlineFeatureTable
    | dropOpenflowConnector
    | dropOpenflowDeployment
    | dropOpenflowRuntime
    | dropOrganizationProfile
    | dropOrganizationUser
    | dropOrganizationUserGroup
    | dropPackagesPolicy
    | dropPasswordPolicy
    | dropPipe
    | dropPostgresInstance
    | dropPrivacyPolicy
    | dropProcedure
    | dropProjectionPolicy
    | dropReplicationGroup
    | dropResourceMonitor
    | dropRestrictedSessionScope
    | dropRole
    | dropRowAccessPolicy
    | dropSchema
    | dropSecret
    | dropSemanticView
    | dropSequence
    | dropService
    | dropSessionPolicy
    | dropShare
    | dropSnapshot
    | dropSnapshotPolicy
    | dropSnapshotSet
    | dropStage
    | dropStorageLifecyclePolicy
    | dropStream
    | dropStreamlit
    | dropTable
    | undropTable
    | dropTag
    | dropTask
    | dropTaskGraph
    | dropType
    | dropUser
    | dropVariable
    | dropView
    | dropWarehouse
    | executeAlert
    | executeDbtProject
    | executeDcmProject
    | executeImmediate
    | executeImmediateFrom
    | executeInferenceJobService
    | executeJobService
    | executeNotebook
    | executeNotebookProject
    | executeCodeBundle
    | executeTask
    | executeTaskGraph
    | explain
    | get
    | grantApplicationRole
    | grantCaller
    | grantDatabaseRole
    | grantDatabaseRoleShare
    | grantOwnership
    | grantPrivilege
    | grantPrivilegeApplication
    | grantPrivilegeApplicationRole
    | grantPrivilegeShare
    | grantPrivilegeUser
    | grantRole
    | grantServiceRole
    | insertMultiTable
    | list
    | merge
    | put
    | remove
    | revokeApplicationRole
    | revokeCaller
    | revokeDatabaseRole
    | revokeDatabaseRoleShare
    | revokePrivilege
    | revokePrivilegeApplication
    | revokePrivilegeApplicationRole
    | revokePrivilegeShare
    | revokePrivilegeUser
    | revokeRole
    | revokeServiceRole
    | rollback
    | select
    | set
    | showAccounts
    | showAgents
    | showAggregationPolicies
    | showAlerts
    | showApplicationPackages
    | showApplicationRoles
    | showApplicationServices
    | showApplications
    | showArtifactRepositories
    | showAuthenticationPolicies
    | showAvailableListings
    | showAvailableOffers
    | showAvailableOrganizationProfiles
    | showBackupPolicies
    | showBackupSets
    | showBackupsInBackupSet
    | showCallerGrants
    | showCatalogIntegrations
    | showChannels
    | showClasses
    | showColumns
    | showCommittedTasks
    | showComputePoolInstanceFamilies
    | showComputePools
    | showConfigurations
    | showConnections
    | showContacts
    | showCortexBaseModels
    | showCortexSearch
    | showDataMetricFunctions
    | showDataMovementPolicies
    | showDataMovementRules
    | showDatabaseRoles
    | showDatabases
    | showDatabasesInFailoverGroup
    | showDatabasesInReplicationGroup
    | showDatasets
    | showDbtProjects
    | showDcmProjects
    | showDelegatedAuthorizations
    | showDeploymentsInDcmProject
    | showDynamicTables
    | showEndpoints
    | showEntitiesInDcmProject
    | showEventRoutingTableOnOrganization
    | showEventRoutingTables
    | showEventTables
    | showExperiments
    | showExternalAgents
    | showExternalConsumerProgrammaticAccessTokens
    | showExternalConsumers
    | showExternalFunctions
    | showExternalTables
    | showExternalVolumes
    | showFailoverGroups
    | showFeaturePolicies
    | showFileFormats
    | showFilesInArtifactRepository
    | showFunctions
    | showFunctionsInModel
    | showGateways
    | showGitBranches
    | showGitRepositories
    | showGitTags
    | showGlobalAccounts
    | showGrants
    | showGrantsInDcmProject
    | showHybridTables
    | showIcebergTables
    | showImageRepositories
    | showImagesInImageRepository
    | showIndexes
    | showIndexesStandardTables
    | showIngressPolicies
    | showIntegrations
    | showInteractiveTables
    | showJoinPolicies
    | showListings
    | showListingsInFailoverGroup
    | showLocks
    | showMaintenancePolicies
    | showManagedAccounts
    | showMaskingPolicies
    | showMaterializedViews
    | showMcpServers
    | showMfaMethods
    | showModelMonitors
    | showModels
    | showMultiPartyApprovalPolicies
    | showNetworkPolicies
    | showNetworkRules
    | showNodesComputePool
    | showNotebookProjects
    | showCodeBundles
    | showNotebooks
    | showNotificationIntegrations
    | showObjects
    | showObjectsOwnedByApplication
    | showOffers
    | showOflowDataPlaneIntegration
    | showOnlineFeatureTables
    | showOpenflowConnectorDefinitions
    | showOpenflowConnectors
    | showOpenflowDeployments
    | showOpenflowRuntimes
    | showOrganizationAccounts
    | showOrganizationContracts
    | showOrganizationProfiles
    | showOrganizationUserGroups
    | showOrganizationUsers
    | showPackagesInArtifactRepository
    | showPackagesPolicies
    | showParameters
    | showPasswordPolicies
    | showPipes
    | showPostgresInstances
    | showPricingPlans
    | showPrimaryKeys
    | showPrivacyPolicies
    | showPrivileges
    | showProcedures
    | showProjectionPolicies
    | showReferences
    | showRegions
    | showReleaseChannels
    | showReleaseDirectives
    | showReplicationAccounts
    | showReplicationDatabases
    | showReplicationGroups
    | showResourceMonitors
    | showRestrictedSessionScopes
    | showRoles
    | showRolesInService
    | showRowAccessPolicies
    | showRulesInEventRoutingTable
    | showRunInExperiment
    | showRunsInExperiment
    | showSchemaVariables
    | showSchemas
    | showSecrets
    | showSemanticDimensions
    | showSemanticDimensionsForMetric
    | showSemanticFacts
    | showSemanticMetrics
    | showSemanticViews
    | showSequences
    | showServiceContainersInService
    | showServiceInstancesInService
    | showServiceVolumesInService
    | showServices
    | showSessionPolicies
    | showSharedContent
    | showShares
    | showSharesInFailoverGroup
    | showSharesInReplicationGroup
    | showSnapshotPolicies
    | showSnapshotSets
    | showSnapshots
    | showSnapshotsInSnapshotSet
    | showSpecifications
    | showStages
    | showStorageLifecyclePolicies
    | showStreamlits
    | showStreams
    | showTables
    | showTags
    | showTaskGraphs
    | showTasks
    | showTelemetryEventDefinitions
    | showTransactions
    | showTypes
    | showUserFunctions
    | showUserKeyPairs
    | showUserProcedures
    | showUserProgrammaticAccessTokens
    | showUserWorkloadIdentityAuthenticationMethods
    | showUsers
    | showVariables
    | showVersions
    | showVersionsInArtifactRepository
    | showVersionsInDataset
    | showVersionsInDbtProject
    | showVersionsInListing
    | showVersionsInModel
    | showVersionsInOpenflowConnector
    | showVersionsInOrganizationProfile
    | showViews
    | showWarehouses
    | showWorkspaces
    | truncateMaterializedView
    | truncateTable
    | undropAccount
    | undropAgent
    | undropDatabase
    | undropDynamicTable
    | undropSchema
    | undropSnapshot
    | undropStage
    | undropStreamlit
    | undropTag
    | undropType
    | undropVariable
    | unset
    | update
    | useDatabase
    | useRole
    | useSchema
    | useSecondaryRoles
    | useWarehouse
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter
// no syntax

// https://docs.snowflake.com/en/sql-reference/sql/alter-account
// fix: the schedule is a quoted string, a minute interval or a cron schedule, as the parameter says
alterAccount
    : 'ALTER' 'ACCOUNT' 'SET' (alterAccount_accountProperties | alterAccount_accountParams | alterAccount_objectParams | alterAccount_sessionParams)
    | 'ALTER' 'ACCOUNT' 'UNSET' ID (',' ID)*
    | 'ALTER' 'ACCOUNT' 'SET' 'RESOURCE_MONITOR' '=' ID
    | 'ALTER' 'ACCOUNT' 'ADD' 'ORGANIZATION' 'USER' 'GROUP' ID
    | 'ALTER' 'ACCOUNT' 'REMOVE' 'ORGANIZATION' 'USER' 'GROUP' ID
    | 'ALTER' 'ACCOUNT' 'SET' ('AUTHENTICATION' | 'SESSION') 'POLICY' QNAME (('FOR' 'ALL' 'PERSON' 'USERS' | 'FOR' 'ALL' 'SERVICE' 'USERS'))? 'FORCE'?
    | 'ALTER' 'ACCOUNT' 'UNSET' ('AUTHENTICATION' | 'SESSION') 'POLICY' (('FOR' 'ALL' 'PERSON' 'USERS' | 'FOR' 'ALL' 'SERVICE' 'USERS'))?
    | 'ALTER' 'ACCOUNT' 'SET' 'FEATURE' 'POLICY' QNAME 'FOR' 'ALL' 'APPLICATIONS' 'FORCE'?
    | 'ALTER' 'ACCOUNT' 'UNSET' 'FEATURE' 'POLICY' 'FOR' 'ALL' 'APPLICATIONS'
    | 'ALTER' 'ACCOUNT' 'SET' 'FEATURE' 'POLICY' QNAME 'FOR' 'ALL' 'DATABASES' 'FORCE'?
    | 'ALTER' 'ACCOUNT' 'UNSET' 'FEATURE' 'POLICY' 'FOR' 'ALL' 'DATABASES'
    | 'ALTER' 'ACCOUNT' 'SET' 'FEATURE' 'POLICY' QNAME 'FOR' 'ALL' 'PERSONAL' 'DATABASES' 'FORCE'?
    | 'ALTER' 'ACCOUNT' 'UNSET' 'FEATURE' 'POLICY' 'FOR' 'ALL' 'PERSONAL' 'DATABASES'
    | 'ALTER' 'ACCOUNT' 'SET' 'MAINTENANCE' 'POLICY' QNAME 'FORCE'? 'FOR' 'ALL' 'APPLICATIONS'
    | 'ALTER' 'ACCOUNT' 'UNSET' 'MAINTENANCE' 'POLICY' 'FOR' 'ALL' 'APPLICATIONS'
    | 'ALTER' 'ACCOUNT' 'SET' 'DATA' 'MOVEMENT' 'POLICY' QNAME 'FORCE'?
    | 'ALTER' 'ACCOUNT' 'UNSET' 'DATA' 'MOVEMENT' 'POLICY'
    | 'ALTER' 'ACCOUNT' 'SET' ('PACKAGES' | 'PASSWORD' | 'MULTI' 'PARTY' 'APPROVAL') 'POLICY' QNAME 'FORCE'?
    | 'ALTER' 'ACCOUNT' 'UNSET' ('PACKAGES' | 'PASSWORD' | 'MULTI' 'PARTY' 'APPROVAL') 'POLICY'
    | 'ALTER' 'ACCOUNT' 'SET' 'CONTACT' ID '=' ID (',' ID '=' ID)*
    | 'ALTER' 'ACCOUNT' 'UNSET' 'CONTACT' ID
    | 'ALTER' 'ACCOUNT' 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'ACCOUNT' 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'ACCOUNT' ID 'SET' 'IS_ORG_ADMIN' '=' ('TRUE' | 'FALSE')
    | 'ALTER' 'ACCOUNT' ID 'SET' 'EDITION' '=' ('\'STANDARD\'' | '\'ENTERPRISE\'' | '\'BUSINESS_CRITICAL\'')
    | 'ALTER' 'ACCOUNT' ID 'SET' 'TENANT_TYPE' '=' ('INTERNAL' | 'EXTERNAL')
    | 'ALTER' 'ACCOUNT' ID 'SET' 'DOMAIN_NAMES' '=' '(' STRING (',' STRING)* ')'
    | 'ALTER' 'ACCOUNT' ID 'UNSET' 'DOMAIN_NAMES'
    | 'ALTER' 'ACCOUNT' ID 'RENAME' 'TO' ID ('SAVE_OLD_URL' '=' ('TRUE' | 'FALSE'))?
    | 'ALTER' 'ACCOUNT' ID 'DROP' 'OLD' 'URL'
    | 'ALTER' 'ACCOUNT' ID 'DROP' 'OLD' 'ORGANIZATION' 'URL'
    ;
alterAccount_accountProperties
    : ('LOGIN_IDP_REDIRECT' '=' '(' ID '=' ID (',' ID '=' ID)* ')' | 'OBJECT_VISIBILITY' '=' (ID | 'PRIVILEGED'))+
    ;
alterAccount_accountParams
    : ('ALLOW_ID_TOKEN' '=' ('TRUE' | 'FALSE') | 'ALLOWED_SPCS_WORKLOAD_TYPES' '=' (STRING | '\'ALL\'') | 'CLIENT_ENCRYPTION_KEY_SIZE' '=' NUMBER | 'CORTEX_ENABLED_CROSS_REGION' '=' ('\'DISABLED\'' | '\'ANY_REGION\'' | STRING) | 'DISALLOWED_SPCS_WORKLOAD_TYPES' '=' (STRING | '\'ALL\'') | 'DISABLE_EXECUTE_DATA_METRIC_FUNCTION_GRANT_TO_PUBLIC' '=' ('TRUE' | 'FALSE') | 'DISABLE_USER_PRIVILEGE_GRANTS' '=' ('TRUE' | 'FALSE') | 'DEFAULT_DBT_VERSION' '=' STRING | 'ENABLE_EGRESS_COST_OPTIMIZER' '=' ('TRUE' | 'FALSE') | 'ENABLE_INTERNAL_STAGES_PRIVATELINK' '=' ('TRUE' | 'FALSE') | 'ENABLE_SNOWFLAKE_MANAGED_STORAGE_VOLUME_PRIVATELINK' '=' ('TRUE' | 'FALSE') | 'ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES' '=' ('TRUE' | 'FALSE') | 'ENFORCE_NETWORK_RULES_FOR_SNOWFLAKE_MANAGED_STORAGE_VOLUME' '=' ('TRUE' | 'FALSE') | 'ENFORCE_TAG_PROPAGATION_FOR_DATA_MOVEMENT_POLICIES' '=' ('TRUE' | 'FALSE') | 'ENABLE_NOTEBOOK_CREATION_IN_PERSONAL_DB' '=' ('TRUE' | 'FALSE') | 'ENABLE_SPCS_BLOCK_STORAGE_SNOWFLAKE_FULL_ENCRYPTION_ENFORCEMENT' '=' ('TRUE' | 'FALSE') | 'EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST' '=' ('TRUE' | 'FALSE') | 'INITIAL_REPLICATION_SIZE_LIMIT_IN_TB' '=' NUMBER | 'LISTING_AUTO_FULFILLMENT_INITIAL_REFRESH_SIZE_LIMIT_IN_TB' '=' NUMBER | 'LISTING_AUTO_FULFILLMENT_REPLICATION_REFRESH_SCHEDULE' '=' STRING | 'LLM_INFERENCE_PARSE_DOCUMENT_PRESIGNED_URL_EXPIRY_SECONDS' '=' NUMBER | 'MIN_STAGE_UNDROP_RETENTION_TIME_IN_DAYS' '=' NUMBER | 'NETWORK_POLICY' '=' STRING | 'OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST' '=' ('TRUE' | 'FALSE') | 'PERIODIC_DATA_REKEYING' '=' ('TRUE' | 'FALSE') | 'READ_CONSISTENCY_MODE' '=' ('\'SESSION\'' | '\'GLOBAL\'') | 'REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION' '=' ('TRUE' | 'FALSE') | 'REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION' '=' ('TRUE' | 'FALSE') | 'SAML_IDENTITY_PROVIDER' '=' STRING | 'SQL_TRACE_QUERY_TEXT' '=' ('ON' | 'OFF') | 'SSO_LOGIN_PAGE' '=' ('TRUE' | 'FALSE') | 'USE_WORKSPACES_FOR_SQL' '=' ('\'always\'' | '\'never\''))+
    ;
alterAccount_objectParams
    : ('BASE_LOCATION_PREFIX' '=' STRING | 'CATALOG' '=' ID | 'CATALOG_SYNC' '=' STRING | 'CORTEX_MODELS_ALLOWLIST' '=' (STRING | '\'ALL\'' | '\'NONE\'') | 'DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER | 'DEFAULT_DDL_COLLATION' '=' STRING | 'DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' ID | 'DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' ID | 'DEFAULT_STREAMLIT_COMPUTE_POOL' '=' ID | 'DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID | 'ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE') | 'ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\'') | 'ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE') | 'ENABLE_TAG_PROPAGATION_EVENT_LOGGING' '=' ('TRUE' | 'FALSE') | 'ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR' '=' ('TRUE' | 'FALSE') | 'ENABLE_UNREDACTED_SECURE_OBJECT_ERROR' '=' ('TRUE' | 'FALSE') | 'EVENT_TABLE' '=' STRING | 'EXTERNAL_VOLUME' '=' ID | 'ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING | 'ICEBERG_VERSION_DEFAULT' '=' NUMBER | 'LOG_LEVEL' '=' STRING | 'MAX_CONCURRENCY_LEVEL' '=' NUMBER | 'MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER | 'METRIC_LEVEL' '=' STRING | 'NETWORK_POLICY' '=' STRING | 'OAUTH_AUTHORIZATION_SERVER' '=' ID | 'OAUTH_SCOPES_SUPPORTED' '=' STRING | 'PIPE_EXECUTION_PAUSED' '=' ('TRUE' | 'FALSE') | 'PREVENT_UNLOAD_TO_INLINE_URL' '=' ('TRUE' | 'FALSE') | 'PREVENT_UNLOAD_TO_INTERNAL_STAGES' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'STAGE_UNDROP_RETENTION_TIME_IN_DAYS' '=' NUMBER | 'STATEMENT_QUEUED_TIMEOUT_IN_SECONDS' '=' NUMBER | 'STATEMENT_TIMEOUT_IN_SECONDS' '=' NUMBER | 'STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED') | 'TRACE_LEVEL' '=' STRING)+
    ;
alterAccount_sessionParams
    : ('ABORT_DETACHED_QUERY' '=' ('TRUE' | 'FALSE') | 'AUTOCOMMIT' '=' ('TRUE' | 'FALSE') | 'BINARY_INPUT_FORMAT' '=' STRING | 'BINARY_OUTPUT_FORMAT' '=' STRING | 'DATE_INPUT_FORMAT' '=' STRING | 'DATE_OUTPUT_FORMAT' '=' STRING | 'DEFAULT_NULL_ORDERING' '=' STRING | 'ENABLE_GET_DDL_USE_DATA_TYPE_ALIAS' '=' ('TRUE' | 'FALSE') | 'ERROR_ON_NONDETERMINISTIC_MERGE' '=' ('TRUE' | 'FALSE') | 'ERROR_ON_NONDETERMINISTIC_UPDATE' '=' ('TRUE' | 'FALSE') | 'JSON_INDENT' '=' NUMBER | 'LOCK_TIMEOUT' '=' NUMBER | 'OPT_OUT_ERROR_LOGGING' '=' ('TRUE' | 'FALSE') | 'QUERY_TAG' '=' STRING | 'ROWS_PER_RESULTSET' '=' NUMBER | 'S3_STAGE_VPCE_DNS_NAME' '=' STRING | 'SEARCH_PATH' '=' STRING | 'SIMULATED_DATA_SHARING_CONSUMER' '=' STRING | 'STATEMENT_TIMEOUT_IN_SECONDS' '=' NUMBER | 'STRICT_JSON_OUTPUT' '=' ('TRUE' | 'FALSE') | 'TIMESTAMP_DAY_IS_ALWAYS_24H' '=' ('TRUE' | 'FALSE') | 'TIMESTAMP_INPUT_FORMAT' '=' STRING | 'TIMESTAMP_LTZ_OUTPUT_FORMAT' '=' STRING | 'TIMESTAMP_NTZ_OUTPUT_FORMAT' '=' STRING | 'TIMESTAMP_OUTPUT_FORMAT' '=' STRING | 'TIMESTAMP_TYPE_MAPPING' '=' STRING | 'TIMESTAMP_TZ_OUTPUT_FORMAT' '=' STRING | 'TIMEZONE' '=' STRING | 'TIME_INPUT_FORMAT' '=' STRING | 'TIME_OUTPUT_FORMAT' '=' STRING | 'TRANSACTION_DEFAULT_ISOLATION_LEVEL' '=' STRING | 'TWO_DIGIT_CENTURY_START' '=' NUMBER | 'UNSUPPORTED_DDL_ACTION' '=' STRING | 'USE_CACHED_RESULT' '=' ('TRUE' | 'FALSE') | 'WEEK_OF_YEAR_POLICY' '=' NUMBER | 'WEEK_START' '=' NUMBER)+
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-account-set-ingress-policy
alterAccountSetIngressPolicy
    : 'ALTER' 'ACCOUNT' 'SET' 'INGRESS' 'POLICY' QNAME 'FORCE'?
    | 'ALTER' 'ACCOUNT' 'UNSET' 'INGRESS' 'POLICY'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-agent
// fix: the specification is a quoted (or dollar-quoted) YAML or JSON string, as the parameters say
alterAgent
    : 'ALTER' 'AGENT' QNAME 'SET' ('COMMENT' '=' STRING ('PROFILE' '=' STRING)? ('SECURE' '=' ('TRUE' | 'FALSE'))? | 'PROFILE' '=' STRING ('SECURE' '=' ('TRUE' | 'FALSE'))? | 'SECURE' '=' ('TRUE' | 'FALSE'))
    | 'ALTER' 'AGENT' QNAME 'MODIFY' 'LIVE' 'VERSION' 'SET' 'SPECIFICATION' '=' STRING
    | 'ALTER' 'AGENT' QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'AGENT' QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-aggregation-policy
alterAggregationPolicy
    : 'ALTER' 'AGGREGATION' 'POLICY' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'AGGREGATION' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'BODY' '->' EXPR
    | 'ALTER' 'AGGREGATION' 'POLICY' QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'AGGREGATION' 'POLICY' QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'AGGREGATION' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'AGGREGATION' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-alert
// fix: the condition is a query, as the parameters say and CREATE ALERT reads it
// fix: the action is a SQL statement, as the parameters say and CREATE ALERT reads it
// fix: the warehouse is a name, as the parameters say
alterAlert
    : 'ALTER' 'ALERT' ('IF' 'EXISTS')? QNAME ('RESUME' | 'SUSPEND')
    | 'ALTER' 'ALERT' ('IF' 'EXISTS')? QNAME 'SET' ('WAREHOUSE' '=' ID ('SCHEDULE' '=' STRING)? ('COMMENT' '=' STRING)? ('CONFIG' '=' STRING)? ('RUNBOOK' '=' STRING)? ('SUSPEND_ALERT_AFTER_NUM_FAILURES' '=' NUMBER)? | 'SCHEDULE' '=' STRING ('COMMENT' '=' STRING)? ('CONFIG' '=' STRING)? ('RUNBOOK' '=' STRING)? ('SUSPEND_ALERT_AFTER_NUM_FAILURES' '=' NUMBER)? | 'COMMENT' '=' STRING ('CONFIG' '=' STRING)? ('RUNBOOK' '=' STRING)? ('SUSPEND_ALERT_AFTER_NUM_FAILURES' '=' NUMBER)? | 'CONFIG' '=' STRING ('RUNBOOK' '=' STRING)? ('SUSPEND_ALERT_AFTER_NUM_FAILURES' '=' NUMBER)? | 'RUNBOOK' '=' STRING ('SUSPEND_ALERT_AFTER_NUM_FAILURES' '=' NUMBER)? | 'SUSPEND_ALERT_AFTER_NUM_FAILURES' '=' NUMBER)
    | 'ALTER' 'ALERT' ('IF' 'EXISTS')? QNAME 'FROM' 'TEMPLATE' ID 'SET' ('WAREHOUSE' '=' ID ('SCHEDULE' '=' STRING)? ('COMMENT' '=' STRING)? ('RUNBOOK' '=' STRING)? ('SUSPEND_ALERT_AFTER_NUM_FAILURES' '=' NUMBER)? ('TEMPLATE_PARAMS' '=' STRING)? | 'SCHEDULE' '=' STRING ('COMMENT' '=' STRING)? ('RUNBOOK' '=' STRING)? ('SUSPEND_ALERT_AFTER_NUM_FAILURES' '=' NUMBER)? ('TEMPLATE_PARAMS' '=' STRING)? | 'COMMENT' '=' STRING ('RUNBOOK' '=' STRING)? ('SUSPEND_ALERT_AFTER_NUM_FAILURES' '=' NUMBER)? ('TEMPLATE_PARAMS' '=' STRING)? | 'RUNBOOK' '=' STRING ('SUSPEND_ALERT_AFTER_NUM_FAILURES' '=' NUMBER)? ('TEMPLATE_PARAMS' '=' STRING)? | 'SUSPEND_ALERT_AFTER_NUM_FAILURES' '=' NUMBER ('TEMPLATE_PARAMS' '=' STRING)? | 'TEMPLATE_PARAMS' '=' STRING)
    | 'ALTER' 'ALERT' ('IF' 'EXISTS')? QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'ALERT' ('IF' 'EXISTS')? QNAME 'UNSET' ('WAREHOUSE' | 'COMMENT' | 'CONFIG' | 'RUNBOOK' | 'SUSPEND_ALERT_AFTER_NUM_FAILURES') (',' ('WAREHOUSE' | 'COMMENT' | 'CONFIG' | 'RUNBOOK' | 'SUSPEND_ALERT_AFTER_NUM_FAILURES'))*
    | 'ALTER' 'ALERT' QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'ALERT' ('IF' 'EXISTS')? QNAME 'MODIFY' 'CONDITION' 'EXISTS' '(' QUERY ')'
    | 'ALTER' 'ALERT' ('IF' 'EXISTS')? QNAME 'MODIFY' 'ACTION' QUERY
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-api-integration
// fix: read: ('<...>') after _PREFIXES = is a list of URLs, as the parameters say
alterApiIntegration
    : 'ALTER' 'API'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'SET' ('API_AWS_ROLE_ARN' '=' STRING ('AZURE_AD_APPLICATION_ID' '=' STRING)? ('API_KEY' '=' STRING)? ('ENABLED' '=' ('TRUE' | 'FALSE'))? ('API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')')? ('API_BLOCKED_PREFIXES' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_AUTHENTICATION_SECRETS' '=' '(' QNAME (',' QNAME)* ')' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'all' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'none')? ('COMMENT' '=' STRING)? | 'AZURE_AD_APPLICATION_ID' '=' STRING ('API_KEY' '=' STRING)? ('ENABLED' '=' ('TRUE' | 'FALSE'))? ('API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')')? ('API_BLOCKED_PREFIXES' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_AUTHENTICATION_SECRETS' '=' '(' QNAME (',' QNAME)* ')' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'all' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'none')? ('COMMENT' '=' STRING)? | 'API_KEY' '=' STRING ('ENABLED' '=' ('TRUE' | 'FALSE'))? ('API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')')? ('API_BLOCKED_PREFIXES' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_AUTHENTICATION_SECRETS' '=' '(' QNAME (',' QNAME)* ')' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'all' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'none')? ('COMMENT' '=' STRING)? | 'ENABLED' '=' ('TRUE' | 'FALSE') ('API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')')? ('API_BLOCKED_PREFIXES' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_AUTHENTICATION_SECRETS' '=' '(' QNAME (',' QNAME)* ')' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'all' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'none')? ('COMMENT' '=' STRING)? | 'API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')' ('API_BLOCKED_PREFIXES' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_AUTHENTICATION_SECRETS' '=' '(' QNAME (',' QNAME)* ')' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'all' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'none')? ('COMMENT' '=' STRING)? | 'API_BLOCKED_PREFIXES' '=' '(' STRING (',' STRING)* ')' ('ALLOWED_AUTHENTICATION_SECRETS' '=' '(' QNAME (',' QNAME)* ')' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'all' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'none')? ('COMMENT' '=' STRING)? | ('ALLOWED_AUTHENTICATION_SECRETS' '=' '(' QNAME (',' QNAME)* ')' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'all' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'none') ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'API'? 'INTEGRATION' ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'API'? 'INTEGRATION' ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'API'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'UNSET' ('API_KEY' | 'ENABLED' | 'API_BLOCKED_PREFIXES' | 'COMMENT') (',' ('API_KEY' | 'ENABLED' | 'API_BLOCKED_PREFIXES' | 'COMMENT'))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-application
alterApplication
    : 'ALTER' 'APPLICATION' ('IF' 'EXISTS')? ID 'SET' ('COMMENT' '=' STRING ('SHARE_EVENTS_WITH_PROVIDER' '=' ('TRUE' | 'FALSE'))? ('DEBUG_MODE' '=' ('TRUE' | 'FALSE'))? ('AUTHORIZE_RESTRICTED_PROVIDER_REMOTE_OPERATIONS_UNTIL' '=' (STRING | '\'INDEFINITE\'' | '\'NEVER\''))? | 'SHARE_EVENTS_WITH_PROVIDER' '=' ('TRUE' | 'FALSE') ('DEBUG_MODE' '=' ('TRUE' | 'FALSE'))? ('AUTHORIZE_RESTRICTED_PROVIDER_REMOTE_OPERATIONS_UNTIL' '=' (STRING | '\'INDEFINITE\'' | '\'NEVER\''))? | 'DEBUG_MODE' '=' ('TRUE' | 'FALSE') ('AUTHORIZE_RESTRICTED_PROVIDER_REMOTE_OPERATIONS_UNTIL' '=' (STRING | '\'INDEFINITE\'' | '\'NEVER\''))? | 'AUTHORIZE_RESTRICTED_PROVIDER_REMOTE_OPERATIONS_UNTIL' '=' (STRING | '\'INDEFINITE\'' | '\'NEVER\''))
    | 'ALTER' 'APPLICATION' ('IF' 'EXISTS')? ID 'UNSET' ('COMMENT' | 'SHARE_EVENTS_WITH_PROVIDER' | 'DEBUG_MODE' | 'AUTHORIZE_RESTRICTED_PROVIDER_REMOTE_OPERATIONS_UNTIL') (',' ('COMMENT' | 'SHARE_EVENTS_WITH_PROVIDER' | 'DEBUG_MODE' | 'AUTHORIZE_RESTRICTED_PROVIDER_REMOTE_OPERATIONS_UNTIL'))*
    | 'ALTER' 'APPLICATION' ('IF' 'EXISTS')? ID 'RENAME' 'TO' ID
    | 'ALTER' 'APPLICATION' ID 'SET' 'FEATURE' 'POLICY' QNAME 'FORCE'?
    | 'ALTER' 'APPLICATION' ID 'UNSET' 'FEATURE' 'POLICY'
    | 'ALTER' 'APPLICATION' ID 'SET' 'MAINTENANCE' 'POLICY' QNAME 'FORCE'?
    | 'ALTER' 'APPLICATION' ID 'UNSET' 'MAINTENANCE' 'POLICY'
    | 'ALTER' 'APPLICATION' ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'APPLICATION' ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'APPLICATION' ID 'SET' 'SHARED' 'TELEMETRY' 'EVENTS' '(' STRING (',' ID)* ')'
    | 'ALTER' 'APPLICATION' ID 'SET' 'AUTHORIZE_TELEMETRY_EVENT_SHARING' '=' ('TRUE' | 'FALSE')
    | 'ALTER' 'APPLICATION' ID 'UNSET' 'REFERENCES' ('(' STRING (',' STRING)? ')')?
    | 'ALTER' 'APPLICATION' ID 'UPGRADE'
    | 'ALTER' 'APPLICATION' ID 'UPGRADE' 'USING' 'VERSION' ID ('PATCH' NUMBER)?
    | 'ALTER' 'APPLICATION' ID 'UPGRADE' 'USING' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-application-drop-app-spec
alterApplicationDropAppSpec
    : 'ALTER' 'APPLICATION' 'DROP' 'SPECIFICATION' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-application-drop-configuration-definition
alterApplicationDropConfigurationDefinition
    : 'ALTER' 'APPLICATION' 'DROP' 'CONFIGURATION' 'DEFINITION' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-application-package
alterApplicationPackage
    : 'ALTER' 'APPLICATION' 'PACKAGE' ('IF' 'EXISTS')? ID 'SET' ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('COMMENT' '=' STRING)? ('DISTRIBUTION' '=' ('INTERNAL' | 'EXTERNAL'))? ('MULTIPLE_INSTANCES' '=' 'TRUE')? ('ENABLE_RELEASE_CHANNELS' '=' 'TRUE')? ('LISTING_AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('AUTOMATIC_APPLICATION_MAINTENANCE' '=' ('TRUE' | 'FALSE'))? | 'MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER ('DEFAULT_DDL_COLLATION' '=' STRING)? ('COMMENT' '=' STRING)? ('DISTRIBUTION' '=' ('INTERNAL' | 'EXTERNAL'))? ('MULTIPLE_INSTANCES' '=' 'TRUE')? ('ENABLE_RELEASE_CHANNELS' '=' 'TRUE')? ('LISTING_AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('AUTOMATIC_APPLICATION_MAINTENANCE' '=' ('TRUE' | 'FALSE'))? | 'DEFAULT_DDL_COLLATION' '=' STRING ('COMMENT' '=' STRING)? ('DISTRIBUTION' '=' ('INTERNAL' | 'EXTERNAL'))? ('MULTIPLE_INSTANCES' '=' 'TRUE')? ('ENABLE_RELEASE_CHANNELS' '=' 'TRUE')? ('LISTING_AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('AUTOMATIC_APPLICATION_MAINTENANCE' '=' ('TRUE' | 'FALSE'))? | 'COMMENT' '=' STRING ('DISTRIBUTION' '=' ('INTERNAL' | 'EXTERNAL'))? ('MULTIPLE_INSTANCES' '=' 'TRUE')? ('ENABLE_RELEASE_CHANNELS' '=' 'TRUE')? ('LISTING_AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('AUTOMATIC_APPLICATION_MAINTENANCE' '=' ('TRUE' | 'FALSE'))? | 'DISTRIBUTION' '=' ('INTERNAL' | 'EXTERNAL') ('MULTIPLE_INSTANCES' '=' 'TRUE')? ('ENABLE_RELEASE_CHANNELS' '=' 'TRUE')? ('LISTING_AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('AUTOMATIC_APPLICATION_MAINTENANCE' '=' ('TRUE' | 'FALSE'))? | 'MULTIPLE_INSTANCES' '=' 'TRUE' ('ENABLE_RELEASE_CHANNELS' '=' 'TRUE')? ('LISTING_AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('AUTOMATIC_APPLICATION_MAINTENANCE' '=' ('TRUE' | 'FALSE'))? | 'ENABLE_RELEASE_CHANNELS' '=' 'TRUE' ('LISTING_AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('AUTOMATIC_APPLICATION_MAINTENANCE' '=' ('TRUE' | 'FALSE'))? | 'LISTING_AUTO_REFRESH' '=' ('TRUE' | 'FALSE') ('AUTOMATIC_APPLICATION_MAINTENANCE' '=' ('TRUE' | 'FALSE'))? | 'AUTOMATIC_APPLICATION_MAINTENANCE' '=' ('TRUE' | 'FALSE'))
    | 'ALTER' 'APPLICATION' 'PACKAGE' ('IF' 'EXISTS')? ID 'UNSET' ('DATA_RETENTION_TIME_IN_DAYS' | 'MAX_DATA_EXTENSION_TIME_IN_DAYS' | 'DEFAULT_DDL_COLLATION' | 'COMMENT' | 'DISTRIBUTION') (',' ('DATA_RETENTION_TIME_IN_DAYS' | 'MAX_DATA_EXTENSION_TIME_IN_DAYS' | 'DEFAULT_DDL_COLLATION' | 'COMMENT' | 'DISTRIBUTION'))*
    | 'ALTER' 'APPLICATION' 'PACKAGE' ('IF' 'EXISTS')? ID 'RENAME' 'TO' ID
    | 'ALTER' 'APPLICATION' 'PACKAGE' ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'APPLICATION' 'PACKAGE' ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-application-package-release-channel
alterApplicationPackageReleaseChannel
    : 'ALTER' 'APPLICATION' 'PACKAGE' QNAME 'REGISTER' 'VERSION' ID? 'USING' STRING ('LABEL' '=' STRING)?
    | 'ALTER' 'APPLICATION' 'PACKAGE' QNAME 'DEREGISTER' 'VERSION' ID
    | 'ALTER' 'APPLICATION' 'PACKAGE' QNAME 'MODIFY' 'RELEASE' 'CHANNEL' ID 'ADD' 'VERSION' ID
    | 'ALTER' 'APPLICATION' 'PACKAGE' QNAME 'MODIFY' 'RELEASE' 'CHANNEL' ID 'DROP' 'VERSION' ID
    | 'ALTER' 'APPLICATION' 'PACKAGE' QNAME 'MODIFY' 'RELEASE' 'CHANNEL' ID 'ADD' 'ACCOUNTS' '=' '(' ID '.' ID (',' ID '.' ID)* ')'
    | 'ALTER' 'APPLICATION' 'PACKAGE' QNAME 'MODIFY' 'RELEASE' 'CHANNEL' ID 'REMOVE' 'ACCOUNTS' '=' '(' ID '.' ID (',' ID '.' ID)* ')'
    | 'ALTER' 'APPLICATION' 'PACKAGE' QNAME 'MODIFY' 'RELEASE' 'CHANNEL' ID 'SET' 'ACCOUNTS' '=' '(' ID '.' ID (',' ID '.' ID)* ')'
    | 'ALTER' 'APPLICATION' 'PACKAGE' QNAME 'MODIFY' 'RELEASE' 'CHANNEL' ID 'SET' 'DEFAULT' 'RELEASE' 'DIRECTIVE' 'VERSION' '=' ID 'PATCH' '=' NUMBER ('UPGRADE_AFTER' '=' STRING)? ('UPGRADE_IN_MAINTENANCE_WINDOW' '=' ('TRUE' | 'FALSE'))? ('UPGRADE_DEADLINE' '=' STRING)?
    | 'ALTER' 'APPLICATION' 'PACKAGE' QNAME 'MODIFY' 'RELEASE' 'CHANNEL' ID 'SET' 'RELEASE' 'DIRECTIVE' ID 'ACCOUNTS' '=' '(' ID '.' ID (',' ID '.' ID)* ')' 'VERSION' '=' ID 'PATCH' '=' NUMBER ('UPGRADE_AFTER' '=' STRING)? ('UPGRADE_IN_MAINTENANCE_WINDOW' '=' ('TRUE' | 'FALSE'))? ('UPGRADE_DEADLINE' '=' STRING)?
    | 'ALTER' 'APPLICATION' 'PACKAGE' QNAME 'MODIFY' 'RELEASE' 'CHANNEL' ID 'MODIFY' 'RELEASE' 'DIRECTIVE' ID 'VERSION' '=' ID 'PATCH' '=' NUMBER ('UPGRADE_AFTER' '=' STRING)? ('UPGRADE_IN_MAINTENANCE_WINDOW' '=' ('TRUE' | 'FALSE'))? ('UPGRADE_DEADLINE' '=' STRING)?
    | 'ALTER' 'APPLICATION' 'PACKAGE' QNAME 'MODIFY' 'RELEASE' 'CHANNEL' ID 'UNSET' 'RELEASE' 'DIRECTIVE' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-application-package-release-directive
alterApplicationPackageReleaseDirective
    : 'ALTER' 'APPLICATION' 'PACKAGE' QNAME 'MODIFY' 'RELEASE' 'DIRECTIVE' ID 'VERSION' '=' ID 'PATCH' '=' NUMBER ('UPGRADE_AFTER' '=' STRING)? ('UPGRADE_IN_MAINTENANCE_WINDOW' '=' ('TRUE' | 'FALSE'))? ('UPGRADE_DEADLINE' '=' STRING)?
    | 'ALTER' 'APPLICATION' 'PACKAGE' QNAME ('MODIFY' 'RELEASE' 'CHANNEL' ID)? 'MODIFY' 'RELEASE' 'DIRECTIVE' ID 'ADD' 'ACCOUNTS' '=' '(' ID '.' ID (',' ID '.' ID)* ')' ('VERSION' '=' ID 'PATCH' '=' NUMBER)? 'FORCE'?
    | 'ALTER' 'APPLICATION' 'PACKAGE' QNAME ('MODIFY' 'RELEASE' 'CHANNEL' ID)? 'MODIFY' 'RELEASE' 'DIRECTIVE' ID 'REMOVE' 'ACCOUNTS' '=' '(' ID '.' ID (',' ID '.' ID)* ')' ('VERSION' '=' ID 'PATCH' '=' NUMBER)?
    | 'ALTER' 'APPLICATION' 'PACKAGE' QNAME 'SET' 'DEFAULT' 'RELEASE' 'DIRECTIVE' 'VERSION' '=' ID 'PATCH' '=' NUMBER ('UPGRADE_AFTER' '=' STRING)? ('UPGRADE_IN_MAINTENANCE_WINDOW' '=' ('TRUE' | 'FALSE'))? ('UPGRADE_DEADLINE' '=' STRING)?
    | 'ALTER' 'APPLICATION' 'PACKAGE' QNAME 'SET' 'RELEASE' 'DIRECTIVE' ID 'ACCOUNTS' '=' '(' ID '.' ID (',' ID '.' ID)* ')' 'VERSION' '=' ID 'PATCH' '=' NUMBER ('UPGRADE_AFTER' '=' STRING)? ('UPGRADE_IN_MAINTENANCE_WINDOW' '=' ('TRUE' | 'FALSE'))? ('UPGRADE_DEADLINE' '=' STRING)? 'FORCE'?
    | 'ALTER' 'APPLICATION' 'PACKAGE' QNAME 'UNSET' 'RELEASE' 'DIRECTIVE' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-application-package-version
alterApplicationPackageVersion
    : 'ALTER' 'APPLICATION' 'PACKAGE' QNAME 'ADD' 'VERSION' ID? 'USING' STRING ('LABEL' '=' STRING)?
    | 'ALTER' 'APPLICATION' 'PACKAGE' QNAME 'DROP' 'VERSION' ID 'ASYNC'?
    | 'ALTER' 'APPLICATION' 'PACKAGE' QNAME 'ADD' 'PATCH' NUMBER? 'FOR' 'VERSION' ID? 'USING' STRING ('LABEL' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-application-role
alterApplicationRole
    : 'ALTER' 'APPLICATION' 'ROLE' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'APPLICATION' 'ROLE' ('IF' 'EXISTS')? QNAME 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'APPLICATION' 'ROLE' ('IF' 'EXISTS')? QNAME 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-application-sequence-number
alterApplicationSequenceNumber
    : 'ALTER' 'APPLICATION' ID ('APPROVE' | 'DECLINE') 'SPECIFICATION' ID 'SEQUENCE_NUMBER' '=' NUMBER
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-application-service
alterApplicationService
    : 'ALTER' 'APPLICATION' 'SERVICE' ('IF' 'EXISTS')? QNAME ('SUSPEND' | 'RESUME')
    | 'ALTER' 'APPLICATION' 'SERVICE' ('IF' 'EXISTS')? QNAME 'UPGRADE' ('TO' 'VERSION' ID)?
    | 'ALTER' 'APPLICATION' 'SERVICE' ('IF' 'EXISTS')? QNAME 'SET' ('AUTO_RESUME' '=' ('TRUE' | 'FALSE') ('AUTO_SUSPEND_SECS' '=' NUMBER)? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('QUERY_WAREHOUSE' '=' ID)? ('MIN_INSTANCES' '=' NUMBER)? ('MAX_INSTANCES' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'AUTO_SUSPEND_SECS' '=' NUMBER ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('QUERY_WAREHOUSE' '=' ID)? ('MIN_INSTANCES' '=' NUMBER)? ('MAX_INSTANCES' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')' ('QUERY_WAREHOUSE' '=' ID)? ('MIN_INSTANCES' '=' NUMBER)? ('MAX_INSTANCES' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'QUERY_WAREHOUSE' '=' ID ('MIN_INSTANCES' '=' NUMBER)? ('MAX_INSTANCES' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'MIN_INSTANCES' '=' NUMBER ('MAX_INSTANCES' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'MAX_INSTANCES' '=' NUMBER ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'APPLICATION' 'SERVICE' ('IF' 'EXISTS')? QNAME 'UNSET' ('AUTO_RESUME' | 'AUTO_SUSPEND_SECS' | 'EXTERNAL_ACCESS_INTEGRATIONS' | 'QUERY_WAREHOUSE' | 'MIN_INSTANCES' | 'MAX_INSTANCES' | 'COMMENT') (',' ('AUTO_RESUME' | 'AUTO_SUSPEND_SECS' | 'EXTERNAL_ACCESS_INTEGRATIONS' | 'QUERY_WAREHOUSE' | 'MIN_INSTANCES' | 'MAX_INSTANCES' | 'COMMENT'))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-application-service
createApplicationService
    : 'CREATE' 'APPLICATION' 'SERVICE' ('IF' 'NOT' 'EXISTS')? QNAME 'FROM' 'ARTIFACT' 'REPOSITORY' QNAME 'PACKAGE' ID ('VERSION' ID)? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('QUERY_WAREHOUSE' '=' ID)? ('EXECUTE_AS_ROLE' '=' ID)? ('AUTO_RESUME' '=' ('TRUE' | 'FALSE'))? ('AUTO_SUSPEND_SECS' '=' NUMBER)? ('MIN_INSTANCES' '=' NUMBER)? ('MAX_INSTANCES' '=' NUMBER)? ('COMMENT' '=' STRING)? ('SPECIFICATION' '=' STRING)?
    | 'CREATE' 'OR' 'ALTER' 'APPLICATION' 'SERVICE' QNAME 'FROM' 'ARTIFACT' 'REPOSITORY' QNAME 'PACKAGE' ID ('VERSION' ID)? ('COMMENT' '=' STRING)? 'SPECIFICATION' '=' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-application-set-app-spec
alterApplicationSetAppSpec
    : 'ALTER' 'APPLICATION' 'SET' 'SPECIFICATION' ID 'TYPE' '=' 'EXTERNAL_ACCESS' 'LABEL' '=' STRING 'DESCRIPTION' '=' STRING ('HOST_PORTS' | 'PRIVATE_HOST_PORTS') '=' '(' STRING (',' STRING)* ')'
    | 'ALTER' 'APPLICATION' 'SET' 'SPECIFICATION' ID 'TYPE' '=' 'SECURITY_INTEGRATION' 'LABEL' '=' STRING 'DESCRIPTION' '=' STRING 'OAUTH_TYPE' '=' '\'CLIENT_CREDENTIALS\'' 'OAUTH_TOKEN_ENDPOINT' '=' STRING 'OAUTH_ALLOWED_SCOPES' '=' '(' STRING (',' STRING)* ')'
    | 'ALTER' 'APPLICATION' 'SET' 'SPECIFICATION' ID 'TYPE' '=' 'SECURITY_INTEGRATION' 'LABEL' '=' STRING 'DESCRIPTION' '=' STRING 'OAUTH_TYPE' '=' '\'AUTHORIZATION_CODE\'' 'OAUTH_TOKEN_ENDPOINT' '=' STRING ('OAUTH_AUTHORIZATION_ENDPOINT' '=' STRING)? ('OAUTH_ALLOWED_SCOPES' '=' '(' STRING (',' STRING)* ')')?
    | 'ALTER' 'APPLICATION' 'SET' 'SPECIFICATION' ID 'TYPE' '=' 'SECURITY_INTEGRATION' 'LABEL' '=' STRING 'DESCRIPTION' '=' STRING 'OAUTH_TYPE' '=' '\'JWT_BEARER\'' 'OAUTH_TOKEN_ENDPOINT' '=' STRING ('OAUTH_AUTHORIZATION_ENDPOINT' '=' STRING)? ('OAUTH_ALLOWED_SCOPES' '=' '(' STRING (',' STRING)* ')')?
    | 'ALTER' 'APPLICATION' 'SET' 'SPECIFICATION' ID 'TYPE' '=' 'LISTING' 'LABEL' '=' STRING 'DESCRIPTION' '=' STRING 'TARGET_ACCOUNTS' '=' STRING 'LISTING' '=' ID ('AUTO_FULFILLMENT_REFRESH_SCHEDULE' '=' STRING)?
    | 'ALTER' 'APPLICATION' 'SET' 'SPECIFICATION' ID 'TYPE' '=' 'CONNECTION' 'LABEL' '=' STRING 'DESCRIPTION' '=' STRING 'SERVER_APPLICATION' '=' ID 'SERVER_APPLICATION_ROLES' '=' '(' ID (',' ID)* ')'
    | 'ALTER' 'APPLICATION' 'SET' 'SPECIFICATION' ID 'TYPE' '=' 'SETTING' 'LABEL' '=' STRING 'DESCRIPTION' '=' STRING 'SETTING' '=' ID ('VALUE' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-application-set-configuration-definition
alterApplicationSetConfigurationDefinition
    : 'ALTER' 'APPLICATION' 'SET' 'CONFIGURATION' 'DEFINITION' ID 'TYPE' '=' 'APPLICATION_NAME' 'LABEL' '=' STRING 'DESCRIPTION' '=' STRING 'APPLICATION_ROLES' '=' '(' ID (',' ID)* ')'
    | 'ALTER' 'APPLICATION' 'SET' 'CONFIGURATION' 'DEFINITION' ID 'TYPE' '=' 'STRING' 'LABEL' '=' STRING 'DESCRIPTION' '=' STRING 'APPLICATION_ROLES' '=' '(' ID (',' ID)* ')' 'SENSITIVE' '=' ('TRUE' | 'FALSE')
    | 'ALTER' 'APPLICATION' 'SET' 'CONFIGURATION' 'DEFINITION' ID 'TYPE' '=' 'SECRET_AUTHORIZATION' 'SECRET' '=' QNAME '.' ID 'LABEL' '=' STRING 'DESCRIPTION' '=' STRING 'APPLICATION_ROLES' '=' '(' ID (',' ID)* ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-application-set-configuration-value
alterApplicationSetConfigurationValue
    : 'ALTER' 'APPLICATION' ID 'SET' 'CONFIGURATION' ID 'VALUE' '=' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-application-unset-configuration
alterApplicationUnsetConfiguration
    : 'ALTER' 'APPLICATION' ID 'UNSET' 'CONFIGURATION' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-artifact-repository
alterArtifactRepository
    : 'ALTER' 'ARTIFACT' 'REPOSITORY' ('IF' 'EXISTS')? QNAME 'SET' ('COMMENT' '=' STRING ('INDEX_URL' '=' STRING)? ('AUTHENTICATION_SECRET' '=' QNAME)? ('PACKAGE_RETENTION' '=' 'ON_OBJECT_CREATION')? | 'INDEX_URL' '=' STRING ('AUTHENTICATION_SECRET' '=' QNAME)? ('PACKAGE_RETENTION' '=' 'ON_OBJECT_CREATION')? | 'AUTHENTICATION_SECRET' '=' QNAME ('PACKAGE_RETENTION' '=' 'ON_OBJECT_CREATION')? | 'PACKAGE_RETENTION' '=' 'ON_OBJECT_CREATION')
    | 'ALTER' 'ARTIFACT' 'REPOSITORY' ('IF' 'EXISTS')? QNAME 'UNSET' ('COMMENT' | 'AUTHENTICATION_SECRET' | 'PACKAGE_RETENTION')
    | 'ALTER' 'ARTIFACT' 'REPOSITORY' ('IF' 'EXISTS')? QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'ARTIFACT' 'REPOSITORY' ('IF' 'EXISTS')? QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-artifact-repository
createArtifactRepository
    : 'CREATE' ('OR' 'REPLACE')? 'ARTIFACT' 'REPOSITORY' QNAME 'TYPE' '=' ('APPLICATION' | 'PYPI') ('API_INTEGRATION' '=' STRING)? ('INDEX_URL' '=' STRING)? ('AUTHENTICATION_SECRET' '=' QNAME)? ('PACKAGE_RETENTION' '=' 'ON_OBJECT_CREATION')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' 'ARTIFACT' 'REPOSITORY' ('IF' 'NOT' 'EXISTS')? QNAME 'TYPE' '=' ('APPLICATION' | 'PYPI') ('API_INTEGRATION' '=' STRING)? ('INDEX_URL' '=' STRING)? ('AUTHENTICATION_SECRET' '=' QNAME)? ('PACKAGE_RETENTION' '=' 'ON_OBJECT_CREATION')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-authentication-policy
// fix: MFA_POLICY: the parameters list the properties, space-delimited
// fix: PAT_POLICY: the parameters list the properties, space-delimited
// fix: WORKLOAD_IDENTITY_POLICY: the parameters list the properties, space-delimited
alterAuthenticationPolicy
    : 'ALTER' 'AUTHENTICATION' 'POLICY' QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'AUTHENTICATION' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' ('AUTHENTICATION_METHODS' '=' '(' STRING (',' STRING)* ')' ('CLIENT_TYPES' '=' '(' STRING (',' STRING)* ')')? ('CLIENT_POLICY' '=' '(' ID '=' '(' 'MINIMUM_VERSION' '=' STRING ')' (',' ID '=' '(' 'MINIMUM_VERSION' '=' STRING ')')* ')')? ('SECURITY_INTEGRATIONS' '=' '(' STRING (',' STRING)* ')')? ('MFA_ENROLLMENT' '=' ('\'REQUIRED\'' | '\'REQUIRED_PASSWORD_ONLY\'' | '\'OPTIONAL\''))? ('MFA_POLICY' '=' '(' ('ALLOWED_METHODS' '=' '(' ('\'ALL\'' | '\'PASSKEY\'' | '\'TOTP\'' | '\'OTP\'' | '\'DUO\'') (',' ('\'ALL\'' | '\'PASSKEY\'' | '\'TOTP\'' | '\'OTP\'' | '\'DUO\''))* ')' | 'ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION' '=' ('\'ALL\'' | '\'NONE\''))+ ')')? ('PAT_POLICY' '=' '(' ('DEFAULT_EXPIRY_IN_DAYS' '=' NUMBER | 'MAX_EXPIRY_IN_DAYS' '=' NUMBER | 'NETWORK_POLICY_EVALUATION' '=' ('ENFORCED_REQUIRED' | 'ENFORCED_NOT_REQUIRED' | 'NOT_ENFORCED') | 'REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS' '=' ('TRUE' | 'FALSE'))+ ')')? ('WORKLOAD_IDENTITY_POLICY' '=' '(' ('ALLOWED_PROVIDERS' '=' '(' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC') (',' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC'))* ')' | 'ALLOWED_AWS_ACCOUNTS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_AZURE_ISSUERS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_OIDC_ISSUERS' '=' '(' STRING (',' STRING)* ')')+ ')')? ('COMMENT' '=' STRING)? | 'CLIENT_TYPES' '=' '(' STRING (',' STRING)* ')' ('CLIENT_POLICY' '=' '(' ID '=' '(' 'MINIMUM_VERSION' '=' STRING ')' (',' ID '=' '(' 'MINIMUM_VERSION' '=' STRING ')')* ')')? ('SECURITY_INTEGRATIONS' '=' '(' STRING (',' STRING)* ')')? ('MFA_ENROLLMENT' '=' ('\'REQUIRED\'' | '\'REQUIRED_PASSWORD_ONLY\'' | '\'OPTIONAL\''))? ('MFA_POLICY' '=' '(' ('ALLOWED_METHODS' '=' '(' ('\'ALL\'' | '\'PASSKEY\'' | '\'TOTP\'' | '\'OTP\'' | '\'DUO\'') (',' ('\'ALL\'' | '\'PASSKEY\'' | '\'TOTP\'' | '\'OTP\'' | '\'DUO\''))* ')' | 'ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION' '=' ('\'ALL\'' | '\'NONE\''))+ ')')? ('PAT_POLICY' '=' '(' ('DEFAULT_EXPIRY_IN_DAYS' '=' NUMBER | 'MAX_EXPIRY_IN_DAYS' '=' NUMBER | 'NETWORK_POLICY_EVALUATION' '=' ('ENFORCED_REQUIRED' | 'ENFORCED_NOT_REQUIRED' | 'NOT_ENFORCED') | 'REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS' '=' ('TRUE' | 'FALSE'))+ ')')? ('WORKLOAD_IDENTITY_POLICY' '=' '(' ('ALLOWED_PROVIDERS' '=' '(' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC') (',' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC'))* ')' | 'ALLOWED_AWS_ACCOUNTS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_AZURE_ISSUERS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_OIDC_ISSUERS' '=' '(' STRING (',' STRING)* ')')+ ')')? ('COMMENT' '=' STRING)? | 'CLIENT_POLICY' '=' '(' ID '=' '(' 'MINIMUM_VERSION' '=' STRING ')' (',' ID '=' '(' 'MINIMUM_VERSION' '=' STRING ')')* ')' ('SECURITY_INTEGRATIONS' '=' '(' STRING (',' STRING)* ')')? ('MFA_ENROLLMENT' '=' ('\'REQUIRED\'' | '\'REQUIRED_PASSWORD_ONLY\'' | '\'OPTIONAL\''))? ('MFA_POLICY' '=' '(' ('ALLOWED_METHODS' '=' '(' ('\'ALL\'' | '\'PASSKEY\'' | '\'TOTP\'' | '\'OTP\'' | '\'DUO\'') (',' ('\'ALL\'' | '\'PASSKEY\'' | '\'TOTP\'' | '\'OTP\'' | '\'DUO\''))* ')' | 'ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION' '=' ('\'ALL\'' | '\'NONE\''))+ ')')? ('PAT_POLICY' '=' '(' ('DEFAULT_EXPIRY_IN_DAYS' '=' NUMBER | 'MAX_EXPIRY_IN_DAYS' '=' NUMBER | 'NETWORK_POLICY_EVALUATION' '=' ('ENFORCED_REQUIRED' | 'ENFORCED_NOT_REQUIRED' | 'NOT_ENFORCED') | 'REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS' '=' ('TRUE' | 'FALSE'))+ ')')? ('WORKLOAD_IDENTITY_POLICY' '=' '(' ('ALLOWED_PROVIDERS' '=' '(' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC') (',' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC'))* ')' | 'ALLOWED_AWS_ACCOUNTS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_AZURE_ISSUERS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_OIDC_ISSUERS' '=' '(' STRING (',' STRING)* ')')+ ')')? ('COMMENT' '=' STRING)? | 'SECURITY_INTEGRATIONS' '=' '(' STRING (',' STRING)* ')' ('MFA_ENROLLMENT' '=' ('\'REQUIRED\'' | '\'REQUIRED_PASSWORD_ONLY\'' | '\'OPTIONAL\''))? ('MFA_POLICY' '=' '(' ('ALLOWED_METHODS' '=' '(' ('\'ALL\'' | '\'PASSKEY\'' | '\'TOTP\'' | '\'OTP\'' | '\'DUO\'') (',' ('\'ALL\'' | '\'PASSKEY\'' | '\'TOTP\'' | '\'OTP\'' | '\'DUO\''))* ')' | 'ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION' '=' ('\'ALL\'' | '\'NONE\''))+ ')')? ('PAT_POLICY' '=' '(' ('DEFAULT_EXPIRY_IN_DAYS' '=' NUMBER | 'MAX_EXPIRY_IN_DAYS' '=' NUMBER | 'NETWORK_POLICY_EVALUATION' '=' ('ENFORCED_REQUIRED' | 'ENFORCED_NOT_REQUIRED' | 'NOT_ENFORCED') | 'REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS' '=' ('TRUE' | 'FALSE'))+ ')')? ('WORKLOAD_IDENTITY_POLICY' '=' '(' ('ALLOWED_PROVIDERS' '=' '(' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC') (',' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC'))* ')' | 'ALLOWED_AWS_ACCOUNTS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_AZURE_ISSUERS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_OIDC_ISSUERS' '=' '(' STRING (',' STRING)* ')')+ ')')? ('COMMENT' '=' STRING)? | 'MFA_ENROLLMENT' '=' ('\'REQUIRED\'' | '\'REQUIRED_PASSWORD_ONLY\'' | '\'OPTIONAL\'') ('MFA_POLICY' '=' '(' ('ALLOWED_METHODS' '=' '(' ('\'ALL\'' | '\'PASSKEY\'' | '\'TOTP\'' | '\'OTP\'' | '\'DUO\'') (',' ('\'ALL\'' | '\'PASSKEY\'' | '\'TOTP\'' | '\'OTP\'' | '\'DUO\''))* ')' | 'ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION' '=' ('\'ALL\'' | '\'NONE\''))+ ')')? ('PAT_POLICY' '=' '(' ('DEFAULT_EXPIRY_IN_DAYS' '=' NUMBER | 'MAX_EXPIRY_IN_DAYS' '=' NUMBER | 'NETWORK_POLICY_EVALUATION' '=' ('ENFORCED_REQUIRED' | 'ENFORCED_NOT_REQUIRED' | 'NOT_ENFORCED') | 'REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS' '=' ('TRUE' | 'FALSE'))+ ')')? ('WORKLOAD_IDENTITY_POLICY' '=' '(' ('ALLOWED_PROVIDERS' '=' '(' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC') (',' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC'))* ')' | 'ALLOWED_AWS_ACCOUNTS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_AZURE_ISSUERS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_OIDC_ISSUERS' '=' '(' STRING (',' STRING)* ')')+ ')')? ('COMMENT' '=' STRING)? | 'MFA_POLICY' '=' '(' ('ALLOWED_METHODS' '=' '(' ('\'ALL\'' | '\'PASSKEY\'' | '\'TOTP\'' | '\'OTP\'' | '\'DUO\'') (',' ('\'ALL\'' | '\'PASSKEY\'' | '\'TOTP\'' | '\'OTP\'' | '\'DUO\''))* ')' | 'ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION' '=' ('\'ALL\'' | '\'NONE\''))+ ')' ('PAT_POLICY' '=' '(' ('DEFAULT_EXPIRY_IN_DAYS' '=' NUMBER | 'MAX_EXPIRY_IN_DAYS' '=' NUMBER | 'NETWORK_POLICY_EVALUATION' '=' ('ENFORCED_REQUIRED' | 'ENFORCED_NOT_REQUIRED' | 'NOT_ENFORCED') | 'REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS' '=' ('TRUE' | 'FALSE'))+ ')')? ('WORKLOAD_IDENTITY_POLICY' '=' '(' ('ALLOWED_PROVIDERS' '=' '(' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC') (',' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC'))* ')' | 'ALLOWED_AWS_ACCOUNTS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_AZURE_ISSUERS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_OIDC_ISSUERS' '=' '(' STRING (',' STRING)* ')')+ ')')? ('COMMENT' '=' STRING)? | 'PAT_POLICY' '=' '(' ('DEFAULT_EXPIRY_IN_DAYS' '=' NUMBER | 'MAX_EXPIRY_IN_DAYS' '=' NUMBER | 'NETWORK_POLICY_EVALUATION' '=' ('ENFORCED_REQUIRED' | 'ENFORCED_NOT_REQUIRED' | 'NOT_ENFORCED') | 'REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS' '=' ('TRUE' | 'FALSE'))+ ')' ('WORKLOAD_IDENTITY_POLICY' '=' '(' ('ALLOWED_PROVIDERS' '=' '(' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC') (',' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC'))* ')' | 'ALLOWED_AWS_ACCOUNTS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_AZURE_ISSUERS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_OIDC_ISSUERS' '=' '(' STRING (',' STRING)* ')')+ ')')? ('COMMENT' '=' STRING)? | 'WORKLOAD_IDENTITY_POLICY' '=' '(' ('ALLOWED_PROVIDERS' '=' '(' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC') (',' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC'))* ')' | 'ALLOWED_AWS_ACCOUNTS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_AZURE_ISSUERS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_OIDC_ISSUERS' '=' '(' STRING (',' STRING)* ')')+ ')' ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'AUTHENTICATION' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' ('AUTHENTICATION_METHODS' | 'CLIENT_TYPES' | 'CLIENT_POLICY' | 'SECURITY_INTEGRATIONS' | 'MFA_ENROLLMENT' | 'MFA_POLICY' | 'PAT_POLICY' | 'WORKLOAD_IDENTITY_POLICY' | 'COMMENT') (',' ('AUTHENTICATION_METHODS' | 'CLIENT_TYPES' | 'CLIENT_POLICY' | 'SECURITY_INTEGRATIONS' | 'MFA_ENROLLMENT' | 'MFA_POLICY' | 'PAT_POLICY' | 'WORKLOAD_IDENTITY_POLICY' | 'COMMENT'))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-backup-policy
alterBackupPolicy
    : 'ALTER' 'BACKUP' 'POLICY' QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'BACKUP' 'POLICY' QNAME 'SET' ('COMMENT' '=' STRING ('SCHEDULE' '=' STRING)? ('EXPIRE_AFTER_DAYS' '=' NUMBER)? | 'SCHEDULE' '=' STRING ('EXPIRE_AFTER_DAYS' '=' NUMBER)? | 'EXPIRE_AFTER_DAYS' '=' NUMBER)
    | 'ALTER' 'BACKUP' 'POLICY' QNAME 'UNSET' ('COMMENT' | 'SCHEDULE' | 'EXPIRE_AFTER_DAYS')
    | 'ALTER' 'BACKUP' 'POLICY' QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'BACKUP' 'POLICY' QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-backup-set
alterBackupSet
    : 'ALTER' 'BACKUP' 'SET' QNAME 'ADD' 'BACKUP'
    | 'ALTER' 'BACKUP' 'SET' QNAME 'APPLY' 'BACKUP' 'POLICY' QNAME 'FORCE'?
    | 'ALTER' 'BACKUP' 'SET' QNAME 'SUSPEND' 'BACKUP' (('CREATION' | 'EXPIRATION'))? 'POLICY'
    | 'ALTER' 'BACKUP' 'SET' QNAME 'RESUME' 'BACKUP' (('CREATION' | 'EXPIRATION'))? 'POLICY'
    | 'ALTER' 'BACKUP' 'SET' QNAME 'DELETE' 'BACKUP' 'IDENTIFIER' STRING
    | 'ALTER' 'BACKUP' 'SET' QNAME 'MODIFY' 'BACKUP' 'IDENTIFIER' STRING ('ADD' | 'REMOVE') 'LEGAL' 'HOLD'
    | 'ALTER' 'BACKUP' 'SET' QNAME 'MODIFY' 'BACKUP' 'IDENTIFIER' STRING 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'BACKUP' 'SET' QNAME 'MODIFY' 'BACKUP' 'IDENTIFIER' STRING 'UNSET' 'COMMENT'
    | 'ALTER' 'BACKUP' 'SET' QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'BACKUP' 'SET' QNAME 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'BACKUP' 'SET' QNAME 'UNSET' 'COMMENT'
    | 'ALTER' 'BACKUP' 'SET' QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'BACKUP' 'SET' QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-catalog-integration
alterCatalogIntegration
    : 'ALTER' 'CATALOG' 'INTEGRATION' ('IF' 'EXISTS')? ID 'SET' ('REST_AUTHENTICATION' '=' '(' alterCatalogIntegration_restAuthenticationParams ')' ('DEFAULT_STORAGE_CONFIG' '=' '(' 'USE_PRIVATELINK_ENDPOINT' '=' ('TRUE' | 'FALSE') ')')? ('REFRESH_INTERVAL_SECONDS' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'DEFAULT_STORAGE_CONFIG' '=' '(' 'USE_PRIVATELINK_ENDPOINT' '=' ('TRUE' | 'FALSE') ')' ('REFRESH_INTERVAL_SECONDS' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'REFRESH_INTERVAL_SECONDS' '=' NUMBER ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    ;
alterCatalogIntegration_restAuthenticationParams
    : 'OAUTH_CLIENT_SECRET' '=' STRING
    | 'BEARER_TOKEN' '=' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-compute-pool
alterComputePool
    : 'ALTER' 'COMPUTE' 'POOL' ('IF' 'EXISTS')? ID ('SUSPEND' | 'RESUME')
    | 'ALTER' 'COMPUTE' 'POOL' ('IF' 'EXISTS')? ID 'STOP' 'ALL' ('OF' 'TYPE' ID (',' ID)*)?
    | 'ALTER' 'COMPUTE' 'POOL' ('IF' 'EXISTS')? ID 'SET' ('MIN_NODES' '=' NUMBER ('MAX_NODES' '=' NUMBER)? ('AUTO_RESUME' '=' ('TRUE' | 'FALSE'))? ('AUTO_SUSPEND_SECS' '=' NUMBER)? ('PLACEMENT_GROUP' '=' STRING)? ('INSTANCE_FAMILY' '=' ID)? ('BACKUP_INSTANCE_FAMILIES' '=' '(' STRING (',' STRING)* ')')? ('TAG' QNAME '=' STRING (',' QNAME '=' STRING)*)? ('COMMENT' '=' STRING)? | 'MAX_NODES' '=' NUMBER ('AUTO_RESUME' '=' ('TRUE' | 'FALSE'))? ('AUTO_SUSPEND_SECS' '=' NUMBER)? ('PLACEMENT_GROUP' '=' STRING)? ('INSTANCE_FAMILY' '=' ID)? ('BACKUP_INSTANCE_FAMILIES' '=' '(' STRING (',' STRING)* ')')? ('TAG' QNAME '=' STRING (',' QNAME '=' STRING)*)? ('COMMENT' '=' STRING)? | 'AUTO_RESUME' '=' ('TRUE' | 'FALSE') ('AUTO_SUSPEND_SECS' '=' NUMBER)? ('PLACEMENT_GROUP' '=' STRING)? ('INSTANCE_FAMILY' '=' ID)? ('BACKUP_INSTANCE_FAMILIES' '=' '(' STRING (',' STRING)* ')')? ('TAG' QNAME '=' STRING (',' QNAME '=' STRING)*)? ('COMMENT' '=' STRING)? | 'AUTO_SUSPEND_SECS' '=' NUMBER ('PLACEMENT_GROUP' '=' STRING)? ('INSTANCE_FAMILY' '=' ID)? ('BACKUP_INSTANCE_FAMILIES' '=' '(' STRING (',' STRING)* ')')? ('TAG' QNAME '=' STRING (',' QNAME '=' STRING)*)? ('COMMENT' '=' STRING)? | 'PLACEMENT_GROUP' '=' STRING ('INSTANCE_FAMILY' '=' ID)? ('BACKUP_INSTANCE_FAMILIES' '=' '(' STRING (',' STRING)* ')')? ('TAG' QNAME '=' STRING (',' QNAME '=' STRING)*)? ('COMMENT' '=' STRING)? | 'INSTANCE_FAMILY' '=' ID ('BACKUP_INSTANCE_FAMILIES' '=' '(' STRING (',' STRING)* ')')? ('TAG' QNAME '=' STRING (',' QNAME '=' STRING)*)? ('COMMENT' '=' STRING)? | 'BACKUP_INSTANCE_FAMILIES' '=' '(' STRING (',' STRING)* ')' ('TAG' QNAME '=' STRING (',' QNAME '=' STRING)*)? ('COMMENT' '=' STRING)? | 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)* ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'COMPUTE' 'POOL' ('IF' 'EXISTS')? ID 'UNSET' ('AUTO_SUSPEND_SECS' | 'AUTO_RESUME' | 'PLACEMENT_GROUP' | 'BACKUP_INSTANCE_FAMILIES' | 'COMMENT') (',' ('AUTO_SUSPEND_SECS' | 'AUTO_RESUME' | 'PLACEMENT_GROUP' | 'BACKUP_INSTANCE_FAMILIES' | 'COMMENT'))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-connection
alterConnection
    : 'ALTER' 'CONNECTION' ('IF' 'EXISTS')? ID 'ENABLE' 'FAILOVER' 'TO' 'ACCOUNTS' ID '.' ID (',' ID '.' ID)* ('IGNORE' 'EDITION' 'CHECK')?
    | 'ALTER' 'CONNECTION' ('IF' 'EXISTS')? ID 'DISABLE' 'FAILOVER' ('TO' 'ACCOUNTS' ID '.' ID (',' ID '.' ID)*)?
    | 'ALTER' 'CONNECTION' ('IF' 'EXISTS')? ID 'PRIMARY'
    | 'ALTER' 'CONNECTION' ('IF' 'EXISTS')? ID 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'CONNECTION' ('IF' 'EXISTS')? ID 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-contact
alterContact
    : 'ALTER' 'CONTACT' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'CONTACT' ('IF' 'EXISTS')? QNAME 'SET' (('USERS' '=' '(' STRING (',' STRING)* ')' | 'EMAIL_DISTRIBUTION_LIST' '=' STRING | 'URL' '=' STRING) ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-cortex-search
// fix: the parameters say the scoring profile is a JSON string
alterCortexSearch
    : 'ALTER' 'CORTEX' 'SEARCH' 'SERVICE' ('IF' 'EXISTS')? QNAME ('SUSPEND' | 'RESUME') (('INDEXING' | 'SERVING'))?
    | 'ALTER' 'CORTEX' 'SEARCH' 'SERVICE' ('IF' 'EXISTS')? QNAME 'REFRESH'
    | 'ALTER' 'CORTEX' 'SEARCH' 'SERVICE' ('IF' 'EXISTS')? QNAME 'SET' ('TARGET_LAG' '=' STRING ('WAREHOUSE' '=' ID)? ('PRIMARY' 'KEY' '=' '(' ID (',' ID)* ')')? ('FULL_INDEX_BUILD_INTERVAL_DAYS' '=' NUMBER)? ('REQUEST_LOGGING' '=' ('TRUE' | 'FALSE'))? ('AUTO_SUSPEND' '=' (NUMBER | 'NULL'))? ('COMMENT' '=' STRING)? | 'WAREHOUSE' '=' ID ('PRIMARY' 'KEY' '=' '(' ID (',' ID)* ')')? ('FULL_INDEX_BUILD_INTERVAL_DAYS' '=' NUMBER)? ('REQUEST_LOGGING' '=' ('TRUE' | 'FALSE'))? ('AUTO_SUSPEND' '=' (NUMBER | 'NULL'))? ('COMMENT' '=' STRING)? | 'PRIMARY' 'KEY' '=' '(' ID (',' ID)* ')' ('FULL_INDEX_BUILD_INTERVAL_DAYS' '=' NUMBER)? ('REQUEST_LOGGING' '=' ('TRUE' | 'FALSE'))? ('AUTO_SUSPEND' '=' (NUMBER | 'NULL'))? ('COMMENT' '=' STRING)? | 'FULL_INDEX_BUILD_INTERVAL_DAYS' '=' NUMBER ('REQUEST_LOGGING' '=' ('TRUE' | 'FALSE'))? ('AUTO_SUSPEND' '=' (NUMBER | 'NULL'))? ('COMMENT' '=' STRING)? | 'REQUEST_LOGGING' '=' ('TRUE' | 'FALSE') ('AUTO_SUSPEND' '=' (NUMBER | 'NULL'))? ('COMMENT' '=' STRING)? | 'AUTO_SUSPEND' '=' (NUMBER | 'NULL') ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'CORTEX' 'SEARCH' 'SERVICE' ('IF' 'EXISTS')? QNAME 'UNSET' 'PRIMARY' 'KEY'
    | 'ALTER' 'CORTEX' 'SEARCH' 'SERVICE' ('IF' 'EXISTS')? QNAME 'SET' 'ATTRIBUTES' '(' ID (',' ID)* ')'
    | 'ALTER' 'CORTEX' 'SEARCH' 'SERVICE' ('IF' 'EXISTS')? QNAME 'UNSET' 'ATTRIBUTES'
    | 'ALTER' 'CORTEX' 'SEARCH' 'SERVICE' QNAME 'ADD' 'SCORING' 'PROFILE' ('IF' 'NOT' 'EXISTS')? ID STRING
    | 'ALTER' 'CORTEX' 'SEARCH' 'SERVICE' QNAME 'DROP' 'SCORING' 'PROFILE' ('IF' 'EXISTS')? ID
    | 'ALTER' 'CORTEX' 'SEARCH' 'SERVICE' ('IF' 'EXISTS')? QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'CORTEX' 'SEARCH' 'SERVICE' ('IF' 'EXISTS')? QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-data-movement-policy
alterDataMovementPolicy
    : 'ALTER' 'DATA' 'MOVEMENT' 'POLICY' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'DATA' 'MOVEMENT' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'ENFORCE_RULES' '=' '(' ID (',' ID)* ')'
    | 'ALTER' 'DATA' 'MOVEMENT' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'ALERT_RULES' '=' '(' ID (',' ID)* ')'
    | 'ALTER' 'DATA' 'MOVEMENT' 'POLICY' ('IF' 'EXISTS')? QNAME 'ADD' 'ENFORCE_RULES' '=' '(' ID (',' ID)* ')'
    | 'ALTER' 'DATA' 'MOVEMENT' 'POLICY' ('IF' 'EXISTS')? QNAME 'ADD' 'ALERT_RULES' '=' '(' ID (',' ID)* ')'
    | 'ALTER' 'DATA' 'MOVEMENT' 'POLICY' ('IF' 'EXISTS')? QNAME 'REMOVE' 'ENFORCE_RULES' '=' '(' ID (',' ID)* ')'
    | 'ALTER' 'DATA' 'MOVEMENT' 'POLICY' ('IF' 'EXISTS')? QNAME 'REMOVE' 'ALERT_RULES' '=' '(' ID (',' ID)* ')'
    | 'ALTER' 'DATA' 'MOVEMENT' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'DATA' 'MOVEMENT' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-data-movement-rule
alterDataMovementRule
    : 'ALTER' 'DATA' 'MOVEMENT' 'RULE' ('IF' 'EXISTS')? QNAME 'SET' 'MAX_ROWS' 'AS' '(' ')' 'RETURNS' 'INTEGER' '->' '(' EXPR ')'
    | 'ALTER' 'DATA' 'MOVEMENT' 'RULE' ('IF' 'EXISTS')? QNAME 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'DATA' 'MOVEMENT' 'RULE' ('IF' 'EXISTS')? QNAME 'UNSET' 'COMMENT'
    | 'ALTER' 'DATA' 'MOVEMENT' 'RULE' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-database
// fix: UNSET names a property without its value, as the page's parameters say
// fix: the settings are a quoted (dollar-quoted) YAML string, as the parameters say
alterDatabase
    : 'ALTER' 'DATABASE' ('IF' 'EXISTS')? ID 'RENAME' 'TO' ID
    | 'ALTER' 'DATABASE' ('IF' 'EXISTS')? ID 'SWAP' 'WITH' ID
    | 'ALTER' 'DATABASE' ('IF' 'EXISTS')? ID 'SET' ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('STAGE_UNDROP_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('EXTERNAL_VOLUME' '=' ID)? ('CATALOG' '=' ID)? ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_VERSION_DEFAULT' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('OBJECT_VISIBILITY' '=' (ID | 'PRIVILEGED'))? ('LOG_LEVEL' '=' STRING)? ('METRIC_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('EVENT_TABLE' '=' QNAME)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER ('STAGE_UNDROP_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('EXTERNAL_VOLUME' '=' ID)? ('CATALOG' '=' ID)? ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_VERSION_DEFAULT' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('OBJECT_VISIBILITY' '=' (ID | 'PRIVILEGED'))? ('LOG_LEVEL' '=' STRING)? ('METRIC_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('EVENT_TABLE' '=' QNAME)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'STAGE_UNDROP_RETENTION_TIME_IN_DAYS' '=' NUMBER ('EXTERNAL_VOLUME' '=' ID)? ('CATALOG' '=' ID)? ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_VERSION_DEFAULT' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('OBJECT_VISIBILITY' '=' (ID | 'PRIVILEGED'))? ('LOG_LEVEL' '=' STRING)? ('METRIC_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('EVENT_TABLE' '=' QNAME)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'EXTERNAL_VOLUME' '=' ID ('CATALOG' '=' ID)? ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_VERSION_DEFAULT' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('OBJECT_VISIBILITY' '=' (ID | 'PRIVILEGED'))? ('LOG_LEVEL' '=' STRING)? ('METRIC_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('EVENT_TABLE' '=' QNAME)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'CATALOG' '=' ID ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_VERSION_DEFAULT' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('OBJECT_VISIBILITY' '=' (ID | 'PRIVILEGED'))? ('LOG_LEVEL' '=' STRING)? ('METRIC_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('EVENT_TABLE' '=' QNAME)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING ('ICEBERG_VERSION_DEFAULT' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('OBJECT_VISIBILITY' '=' (ID | 'PRIVILEGED'))? ('LOG_LEVEL' '=' STRING)? ('METRIC_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('EVENT_TABLE' '=' QNAME)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'ICEBERG_VERSION_DEFAULT' '=' NUMBER ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('OBJECT_VISIBILITY' '=' (ID | 'PRIVILEGED'))? ('LOG_LEVEL' '=' STRING)? ('METRIC_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('EVENT_TABLE' '=' QNAME)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\'') ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('OBJECT_VISIBILITY' '=' (ID | 'PRIVILEGED'))? ('LOG_LEVEL' '=' STRING)? ('METRIC_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('EVENT_TABLE' '=' QNAME)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE') ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('OBJECT_VISIBILITY' '=' (ID | 'PRIVILEGED'))? ('LOG_LEVEL' '=' STRING)? ('METRIC_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('EVENT_TABLE' '=' QNAME)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') ('DEFAULT_DDL_COLLATION' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('OBJECT_VISIBILITY' '=' (ID | 'PRIVILEGED'))? ('LOG_LEVEL' '=' STRING)? ('METRIC_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('EVENT_TABLE' '=' QNAME)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'DEFAULT_DDL_COLLATION' '=' STRING ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('OBJECT_VISIBILITY' '=' (ID | 'PRIVILEGED'))? ('LOG_LEVEL' '=' STRING)? ('METRIC_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('EVENT_TABLE' '=' QNAME)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('OBJECT_VISIBILITY' '=' (ID | 'PRIVILEGED'))? ('LOG_LEVEL' '=' STRING)? ('METRIC_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('EVENT_TABLE' '=' QNAME)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING ('OBJECT_VISIBILITY' '=' (ID | 'PRIVILEGED'))? ('LOG_LEVEL' '=' STRING)? ('METRIC_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('EVENT_TABLE' '=' QNAME)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'OBJECT_VISIBILITY' '=' (ID | 'PRIVILEGED') ('LOG_LEVEL' '=' STRING)? ('METRIC_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('EVENT_TABLE' '=' QNAME)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'LOG_LEVEL' '=' STRING ('METRIC_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('EVENT_TABLE' '=' QNAME)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'METRIC_LEVEL' '=' STRING ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('EVENT_TABLE' '=' QNAME)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'TRACE_LEVEL' '=' STRING ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('EVENT_TABLE' '=' QNAME)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED') ('EVENT_TABLE' '=' QNAME)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'EVENT_TABLE' '=' QNAME ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'COMMENT' '=' STRING ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'CATALOG_SYNC' '=' STRING ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\'') ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'BASE_LOCATION_PREFIX' '=' STRING ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID)? ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' ID ('CLASSIFICATION_PROFILE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'CLASSIFICATION_PROFILE' '=' STRING ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'CONTACT' ID '=' ID (',' ID '=' ID)* ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE') ('DATA_QUALITY_MONITORING_SETTINGS' '=' STRING)? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'DATA_QUALITY_MONITORING_SETTINGS' '=' STRING ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'OAUTH_AUTHORIZATION_SERVER' '=' ID ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'OAUTH_SCOPES_SUPPORTED' '=' STRING)
    | 'ALTER' 'DATABASE' ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'DATABASE' ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'DATABASE' ('IF' 'EXISTS')? ID 'SET' 'FEATURE' 'POLICY' QNAME 'FORCE'?
    | 'ALTER' 'DATABASE' ('IF' 'EXISTS')? ID 'UNSET' 'FEATURE' 'POLICY'
    | 'ALTER' 'DATABASE' ('IF' 'EXISTS')? ID 'UNSET' ('DATA_RETENTION_TIME_IN_DAYS' | 'MAX_DATA_EXTENSION_TIME_IN_DAYS' | 'STAGE_UNDROP_RETENTION_TIME_IN_DAYS' | 'EXTERNAL_VOLUME' | 'CATALOG' | 'ICEBERG_DEFAULT_DDL_COLLATION' | 'ICEBERG_VERSION_DEFAULT' | 'ICEBERG_MERGE_ON_READ_BEHAVIOR' | 'ENABLE_ICEBERG_MERGE_ON_READ' | 'DEFAULT_DDL_COLLATION' | 'DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' | 'DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' | 'OBJECT_VISIBILITY' | 'STORAGE_SERIALIZATION_POLICY' | 'EVENT_TABLE' | 'COMMENT' | 'CATALOG_SYNC' | 'REPLICABLE_WITH_FAILOVER_GROUPS' | 'BASE_LOCATION_PREFIX' | 'DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' | 'CLASSIFICATION_PROFILE' | 'CONTACT' ID | 'ENABLE_DATA_COMPACTION' | 'OAUTH_AUTHORIZATION_SERVER' | 'OAUTH_SCOPES_SUPPORTED') (',' ('DATA_RETENTION_TIME_IN_DAYS' | 'MAX_DATA_EXTENSION_TIME_IN_DAYS' | 'STAGE_UNDROP_RETENTION_TIME_IN_DAYS' | 'EXTERNAL_VOLUME' | 'CATALOG' | 'ICEBERG_DEFAULT_DDL_COLLATION' | 'ICEBERG_VERSION_DEFAULT' | 'ICEBERG_MERGE_ON_READ_BEHAVIOR' | 'ENABLE_ICEBERG_MERGE_ON_READ' | 'DEFAULT_DDL_COLLATION' | 'DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' | 'DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' | 'OBJECT_VISIBILITY' | 'STORAGE_SERIALIZATION_POLICY' | 'EVENT_TABLE' | 'COMMENT' | 'CATALOG_SYNC' | 'REPLICABLE_WITH_FAILOVER_GROUPS' | 'BASE_LOCATION_PREFIX' | 'DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' | 'CLASSIFICATION_PROFILE' | 'CONTACT' ID | 'ENABLE_DATA_COMPACTION' | 'OAUTH_AUTHORIZATION_SERVER' | 'OAUTH_SCOPES_SUPPORTED'))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-database-catalog-linked
alterDatabaseCatalogLinked
    : 'ALTER' 'DATABASE' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'DATABASE' ('IF' 'EXISTS')? QNAME 'SUSPEND' 'DISCOVERY'
    | 'ALTER' 'DATABASE' ('IF' 'EXISTS')? QNAME 'RESUME' 'DISCOVERY'
    | 'ALTER' 'DATABASE' ('IF' 'EXISTS')? QNAME 'UPDATE' 'LINKED_CATALOG' 'ADD' '(' STRING (',' STRING)* ')' 'TO' 'ALLOWED_NAMESPACES'
    | 'ALTER' 'DATABASE' ('IF' 'EXISTS')? QNAME 'UPDATE' 'LINKED_CATALOG' 'REMOVE' '(' STRING (',' STRING)* ')' 'FROM' 'ALLOWED_NAMESPACES'
    | 'ALTER' 'DATABASE' ('IF' 'EXISTS')? QNAME 'UPDATE' 'LINKED_CATALOG' 'UNSET' 'ALLOWED_NAMESPACES'
    | 'ALTER' 'DATABASE' ('IF' 'EXISTS')? QNAME 'UPDATE' 'LINKED_CATALOG' 'ADD' '(' STRING (',' STRING)* ')' 'TO' 'BLOCKED_NAMESPACES'
    | 'ALTER' 'DATABASE' ('IF' 'EXISTS')? QNAME 'UPDATE' 'LINKED_CATALOG' 'REMOVE' '(' STRING (',' STRING)* ')' 'FROM' 'BLOCKED_NAMESPACES'
    | 'ALTER' 'DATABASE' ('IF' 'EXISTS')? QNAME 'UPDATE' 'LINKED_CATALOG' 'UNSET' 'BLOCKED_NAMESPACES'
    | 'ALTER' 'DATABASE' ('IF' 'EXISTS')? QNAME 'UPDATE' 'LINKED_CATALOG' 'SET' 'SYNC_INTERVAL_SECONDS' '=' NUMBER
    | 'ALTER' 'DATABASE' ('IF' 'EXISTS')? QNAME 'UPDATE' 'LINKED_CATALOG' 'SET' 'ALLOWED_WRITE_OPERATIONS' '=' ('NONE' | 'ALL')
    | 'ALTER' 'DATABASE' ('IF' 'EXISTS')? QNAME 'SET' ('BASE_LOCATION_PREFIX' '=' STRING ('COMMENT' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ICEBERG_VERSION_DEFAULT' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? | 'COMMENT' '=' STRING ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('ICEBERG_VERSION_DEFAULT' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? | 'CONTACT' ID '=' ID (',' ID '=' ID)* ('ICEBERG_VERSION_DEFAULT' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? | 'ICEBERG_VERSION_DEFAULT' '=' NUMBER ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? | 'ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\'') ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? | 'ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))
    | 'ALTER' 'DATABASE' ('IF' 'EXISTS')? QNAME 'UNSET' ('BASE_LOCATION_PREFIX' | 'COMMENT' | 'CONTACT' | 'ICEBERG_VERSION_DEFAULT' | 'ICEBERG_MERGE_ON_READ_BEHAVIOR' | 'ENABLE_ICEBERG_MERGE_ON_READ')
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-database-role
alterDatabaseRole
    : 'ALTER' 'DATABASE' 'ROLE' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'DATABASE' 'ROLE' ('IF' 'EXISTS')? QNAME 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'DATABASE' 'ROLE' ('IF' 'EXISTS')? QNAME 'UNSET' 'COMMENT'
    | 'ALTER' 'DATABASE' 'ROLE' ('IF' 'EXISTS')? QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'DATABASE' 'ROLE' ('IF' 'EXISTS')? QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-dataset
// no syntax

// https://docs.snowflake.com/en/sql-reference/sql/alter-dataset-add-version
alterDatasetAddVersion
    : 'ALTER' 'DATASET' QNAME 'ADD' 'VERSION' ID 'FROM' QUERY ('PARTITION' 'BY' EXPR)? ('COMMENT' '=' STRING)? ('METADATA' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-dataset-drop-version
alterDatasetDropVersion
    : 'ALTER' 'DATASET' ('IF' 'EXISTS')? QNAME 'DROP' 'VERSION' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-dbt-project
alterDbtProject
    : 'ALTER' 'DBT' 'PROJECT' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'DBT' 'PROJECT' QNAME 'DEPLOY' 'FROM' STRING
    | 'ALTER' 'DBT' 'PROJECT' ('IF' 'EXISTS')? QNAME 'SET' ('DBT_VERSION' '=' STRING ('DEFAULT_TARGET' '=' STRING)? ('DEFAULT_ENVIRONMENT' '=' STRING)? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('AUTO_COMPILE' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_WRITEBACK' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'DEFAULT_TARGET' '=' STRING ('DEFAULT_ENVIRONMENT' '=' STRING)? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('AUTO_COMPILE' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_WRITEBACK' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'DEFAULT_ENVIRONMENT' '=' STRING ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('AUTO_COMPILE' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_WRITEBACK' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')' ('AUTO_COMPILE' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_WRITEBACK' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'AUTO_COMPILE' '=' ('TRUE' | 'FALSE') ('DEFAULT_WRITEBACK' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'DEFAULT_WRITEBACK' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'DBT' 'PROJECT' ('IF' 'EXISTS')? QNAME 'UNSET' ('DBT_VERSION' | 'DEFAULT_TARGET' | 'DEFAULT_ENVIRONMENT' | 'EXTERNAL_ACCESS_INTEGRATIONS' | 'AUTO_COMPILE' | 'DEFAULT_WRITEBACK' | 'COMMENT') (',' ('DBT_VERSION' | 'DEFAULT_TARGET' | 'DEFAULT_ENVIRONMENT' | 'EXTERNAL_ACCESS_INTEGRATIONS' | 'AUTO_COMPILE' | 'DEFAULT_WRITEBACK' | 'COMMENT'))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-dcm-project
alterDcmProject
    : 'ALTER' 'DCM' 'PROJECT' ('IF' 'EXISTS')? QNAME 'SET' ('LOG_LEVEL' '=' ID ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'DCM' 'PROJECT' ('IF' 'EXISTS')? QNAME 'UNSET' ('LOG_LEVEL' | 'COMMENT') (',' ('LOG_LEVEL' | 'COMMENT'))*
    | 'ALTER' 'DCM' 'PROJECT' QNAME 'DROP' 'DEPLOYMENT' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-dynamic-table
// fix: dataGovnPolicyTagAction: the UNSET TAG variant ends in } but opens no {; added
// fix: UNSET names a property without its value, as the page's parameters say
// fix: a search method with its targets, as the parameters write it
// not converted (form): undefined fragment "tableColumnCommentAction": the page does not define it, or it did not convert
//   ALTER DYNAMIC TABLE [ IF EXISTS ] <name> { tableColumnCommentAction }
alterDynamicTable
    : 'ALTER' 'DYNAMIC' 'TABLE' ('IF' 'EXISTS')? QNAME ('SUSPEND' | 'RESUME')
    | 'ALTER' 'DYNAMIC' 'TABLE' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'DYNAMIC' 'TABLE' ('IF' 'EXISTS')? QNAME 'SWAP' 'WITH' QNAME
    | 'ALTER' 'DYNAMIC' 'TABLE' ('IF' 'EXISTS')? QNAME (',' QNAME)* 'REFRESH' ('COPY' 'SESSION')?
    | 'ALTER' 'DYNAMIC' 'TABLE' ('IF' 'EXISTS')? QNAME alterDynamicTable_clusteringAction
    | 'ALTER' 'DYNAMIC' 'TABLE' QNAME ('SET' 'COMMENT' '=' STRING | 'UNSET' 'COMMENT')
    | 'ALTER' 'DYNAMIC' 'TABLE' ('IF' 'EXISTS')? QNAME alterDynamicTable_dataGovnPolicyTagAction
    | 'ALTER' 'DYNAMIC' 'TABLE' ('IF' 'EXISTS')? QNAME alterDynamicTable_searchOptimizationAction
    | 'ALTER' 'DYNAMIC' 'TABLE' ('IF' 'EXISTS')? QNAME alterDynamicTable_storageLifecyclePolicyAction
    | 'ALTER' 'DYNAMIC' 'TABLE' ('IF' 'EXISTS')? QNAME 'SET' ('TARGET_LAG' '=' (STRING | 'DOWNSTREAM') | 'SCHEDULER' '=' 'DISABLE' | 'SCHEDULER' '=' 'ENABLE' | 'WAREHOUSE' '=' ID | 'INITIALIZATION_WAREHOUSE' '=' ID | 'DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER | 'MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER | 'DEFAULT_DDL_COLLATION' '=' STRING | 'LOG_LEVEL' '=' STRING | 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')' | 'FROZEN' 'WHERE' '(' EXPR ')' | 'EXECUTE' 'AS' 'USER' ID ('USE' 'SECONDARY' 'ROLES' ('ALL' | 'NONE' | ID (',' ID)*))? | 'ROW_TIMESTAMP' '=' ('TRUE' | 'FALSE')) (',' ('TARGET_LAG' '=' (STRING | 'DOWNSTREAM') | 'SCHEDULER' '=' 'DISABLE' | 'SCHEDULER' '=' 'ENABLE' | 'WAREHOUSE' '=' ID | 'INITIALIZATION_WAREHOUSE' '=' ID | 'DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER | 'MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER | 'DEFAULT_DDL_COLLATION' '=' STRING | 'LOG_LEVEL' '=' STRING | 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')' | 'FROZEN' 'WHERE' '(' EXPR ')' | 'EXECUTE' 'AS' 'USER' ID ('USE' 'SECONDARY' 'ROLES' ('ALL' | 'NONE' | ID (',' ID)*))? | 'ROW_TIMESTAMP' '=' ('TRUE' | 'FALSE')))*
    | 'ALTER' 'DYNAMIC' 'TABLE' ('IF' 'EXISTS')? QNAME 'UNSET' ('INITIALIZATION_WAREHOUSE' | 'DATA_RETENTION_TIME_IN_DAYS' | 'MAX_DATA_EXTENSION_TIME_IN_DAYS' | 'DEFAULT_DDL_COLLATION' | 'LOG_LEVEL' | 'CONTACT' ID | 'FROZEN' 'WHERE' | 'EXECUTE' 'AS' 'USER' | 'ROW_TIMESTAMP') (',' ('INITIALIZATION_WAREHOUSE' | 'DATA_RETENTION_TIME_IN_DAYS' | 'MAX_DATA_EXTENSION_TIME_IN_DAYS' | 'DEFAULT_DDL_COLLATION' | 'LOG_LEVEL' | 'CONTACT' ID | 'FROZEN' 'WHERE' | 'EXECUTE' 'AS' 'USER' | 'ROW_TIMESTAMP'))*
    ;
alterDynamicTable_clusteringAction
    : ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')' | ('SUSPEND' | 'RESUME') 'RECLUSTER' | 'DROP' 'CLUSTERING' 'KEY')
    ;
alterDynamicTable_dataGovnPolicyTagAction
    : ('ADD' 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')' | 'DROP' 'ROW' 'ACCESS' 'POLICY' QNAME | 'DROP' 'ROW' 'ACCESS' 'POLICY' QNAME ',' 'ADD' 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')' | 'DROP' 'ALL' 'ROW' 'ACCESS' 'POLICIES')
    | ('SET' 'AGGREGATION' 'POLICY' QNAME ('ENTITY' 'KEY' '(' ID (',' ID)* ')')? 'FORCE'? | 'UNSET' 'AGGREGATION' 'POLICY')
    | (('ALTER' | 'MODIFY') 'COLUMN'? ID 'SET' 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')? 'FORCE'? | ('ALTER' | 'MODIFY') 'COLUMN'? ID 'UNSET' 'MASKING' 'POLICY')
    | ('ALTER' | 'MODIFY') 'COLUMN'? ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)* ',' 'COLUMN'? ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | (('ALTER' | 'MODIFY') 'COLUMN'? ID 'SET' 'PROJECTION' 'POLICY' QNAME 'FORCE'? | ('ALTER' | 'MODIFY') 'COLUMN'? ID 'UNSET' 'PROJECTION' 'POLICY')
    | ('ALTER' | 'MODIFY') 'COLUMN'? ID 'UNSET' 'TAG' QNAME (',' QNAME)* ',' 'COLUMN'? ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    | ('SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)* | 'UNSET' 'TAG' QNAME (',' QNAME)*)
    ;
alterDynamicTable_searchOptimizationAction
    : ('ADD' 'SEARCH' 'OPTIMIZATION' ('ON' ('FULL_TEXT' | 'EQUALITY' | 'SUBSTRING' | 'GEO') '(' ID (',' ID)* (',' 'ANALYZER' '=>' STRING)? ')' (',' ('FULL_TEXT' | 'EQUALITY' | 'SUBSTRING' | 'GEO') '(' ID (',' ID)* (',' 'ANALYZER' '=>' STRING)? ')')* 'EQUALITY'?)? | 'DROP' 'SEARCH' 'OPTIMIZATION' ('ON' (('FULL_TEXT' | 'EQUALITY' | 'SUBSTRING' | 'GEO') '(' ID (',' ID)* (',' 'ANALYZER' '=>' STRING)? ')' | ID | EXPR) 'EQUALITY'? (',' (('FULL_TEXT' | 'EQUALITY' | 'SUBSTRING' | 'GEO') '(' ID (',' ID)* (',' 'ANALYZER' '=>' STRING)? ')' | ID | EXPR) 'EQUALITY'?)*)? | 'SUSPEND' 'SEARCH' 'OPTIMIZATION' ('ON' (('FULL_TEXT' | 'EQUALITY' | 'SUBSTRING' | 'GEO') '(' ID (',' ID)* (',' 'ANALYZER' '=>' STRING)? ')' | ID | EXPR) (',' (('FULL_TEXT' | 'EQUALITY' | 'SUBSTRING' | 'GEO') '(' ID (',' ID)* (',' 'ANALYZER' '=>' STRING)? ')' | ID | EXPR))*)? | 'RESUME' 'SEARCH' 'OPTIMIZATION' ('ON' (('FULL_TEXT' | 'EQUALITY' | 'SUBSTRING' | 'GEO') '(' ID (',' ID)* (',' 'ANALYZER' '=>' STRING)? ')' | ID | EXPR) (',' (('FULL_TEXT' | 'EQUALITY' | 'SUBSTRING' | 'GEO') '(' ID (',' ID)* (',' 'ANALYZER' '=>' STRING)? ')' | ID | EXPR))*)?)
    ;
alterDynamicTable_storageLifecyclePolicyAction
    : ('ADD' 'STORAGE' 'LIFECYCLE' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')' | 'DROP' 'STORAGE' 'LIFECYCLE' 'POLICY')
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-event-routing-table
alterEventRoutingTable
    : 'ALTER' 'EVENT' 'ROUTING' 'TABLE' QNAME 'FORCE'? 'SET' 'RULES' ID '=' '(' 'REGION_GROUP' '=' ID ',' 'REGIONS' '=' '(' STRING (',' STRING)* ')' ',' 'DESTINATION_ACCOUNT' '=' ID '.' ID ')' (',' ID '=' '(' 'REGION_GROUP' '=' ID ',' 'REGIONS' '=' '(' STRING (',' STRING)* ')' ',' 'DESTINATION_ACCOUNT' '=' ID '.' ID ')')*
    | 'ALTER' 'EVENT' 'ROUTING' 'TABLE' QNAME 'FORCE'? 'SET' 'RULE' ID 'REGION_GROUP' '=' ID 'REGIONS' '=' '(' STRING (',' STRING)* ')' 'DESTINATION_ACCOUNT' '=' ID '.' ID
    | 'ALTER' 'EVENT' 'ROUTING' 'TABLE' QNAME 'FORCE'? 'UNSET' 'RULE' ID
    | 'ALTER' 'EVENT' 'ROUTING' 'TABLE' QNAME 'MODIFY' 'RULE' ID 'RENAME' 'TO' ID
    | 'ALTER' 'EVENT' 'ROUTING' 'TABLE' QNAME 'RENAME' 'TO' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-experiment
alterExperiment
    : 'ALTER' 'EXPERIMENT' QNAME 'ADD' 'RUN' ID
    | 'ALTER' 'EXPERIMENT' QNAME 'COMMIT' 'RUN' ID ('WITH' 'STATUS' '=' ('\'FINISHED\'' | '\'FAILED\'' | '\'CANCELLED\''))?
    | 'ALTER' 'EXPERIMENT' QNAME 'MODIFY' 'RUN' ID 'SET' 'STATUS' '=' ('\'RUNNING\'' | '\'PAUSED\'')
    | 'ALTER' 'EXPERIMENT' QNAME 'DROP' 'RUN' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-external-agent
alterExternalAgent
    : 'ALTER' 'EXTERNAL' 'AGENT' ('IF' 'EXISTS')? QNAME 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'EXTERNAL' 'AGENT' ('IF' 'EXISTS')? QNAME 'ADD' 'VERSION' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-external-consumer
alterExternalConsumer
    : 'ALTER' 'EXTERNAL' 'CONSUMER' ('IF' 'EXISTS')? QNAME 'SET' ('COMMENT' '=' STRING ('EMAIL' '=' STRING)? | 'EMAIL' '=' STRING)
    | 'ALTER' 'EXTERNAL' 'CONSUMER' ('IF' 'EXISTS')? QNAME 'UNSET' ('COMMENT' | 'EMAIL') (',' ('COMMENT' | 'EMAIL'))*
    | 'ALTER' 'EXTERNAL' 'CONSUMER' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-external-consumer-add-programmatic-access-token
alterExternalConsumerAddProgrammaticAccessToken
    : 'ALTER' 'EXTERNAL' 'CONSUMER' QNAME 'ADD' ('PROGRAMMATIC' 'ACCESS' 'TOKEN' | 'PAT') ('IF' 'NOT' 'EXISTS')? ID ('DAYS_TO_EXPIRY' '=' NUMBER)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-external-consumer-modify-programmatic-access-token
alterExternalConsumerModifyProgrammaticAccessToken
    : 'ALTER' 'EXTERNAL' 'CONSUMER' QNAME 'MODIFY' ('PROGRAMMATIC' 'ACCESS' 'TOKEN' | 'PAT') ('IF' 'EXISTS')? ID ('SET' ('DAYS_TO_EXPIRY' '=' NUMBER ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING))? ('UNSET' 'COMMENT')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-external-consumer-remove-programmatic-access-token
alterExternalConsumerRemoveProgrammaticAccessToken
    : 'ALTER' 'EXTERNAL' 'CONSUMER' ('IF' 'EXISTS')? QNAME 'REMOVE' ('PROGRAMMATIC' 'ACCESS' 'TOKEN' | 'PAT') ('IF' 'EXISTS')? ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-external-consumer-rotate-programmatic-access-token
alterExternalConsumerRotateProgrammaticAccessToken
    : 'ALTER' 'EXTERNAL' 'CONSUMER' QNAME 'ROTATE' ('PROGRAMMATIC' 'ACCESS' 'TOKEN' | 'PAT') ('IF' 'EXISTS')? ID ('EXPIRE_ROTATED_TOKEN_AFTER_HOURS' '=' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-external-access-integration
alterExternalAccessIntegration
    : 'ALTER' 'EXTERNAL' 'ACCESS' 'INTEGRATION' ('IF' 'EXISTS')? ID 'SET' ('ALLOWED_NETWORK_RULES' '=' '(' ID (',' ID)* ')' ('ALLOWED_API_AUTHENTICATION_INTEGRATIONS' '=' ('(' ID (',' ID)* ')' | 'none'))? ('ALLOWED_AUTHENTICATION_SECRETS' '=' ('(' QNAME (',' QNAME)* ')' | 'all' | 'none'))? ('ENABLED' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('TAG' QNAME '=' STRING (',' QNAME '=' STRING)*)? | 'ALLOWED_API_AUTHENTICATION_INTEGRATIONS' '=' ('(' ID (',' ID)* ')' | 'none') ('ALLOWED_AUTHENTICATION_SECRETS' '=' ('(' QNAME (',' QNAME)* ')' | 'all' | 'none'))? ('ENABLED' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('TAG' QNAME '=' STRING (',' QNAME '=' STRING)*)? | 'ALLOWED_AUTHENTICATION_SECRETS' '=' ('(' QNAME (',' QNAME)* ')' | 'all' | 'none') ('ENABLED' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('TAG' QNAME '=' STRING (',' QNAME '=' STRING)*)? | 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)? ('TAG' QNAME '=' STRING (',' QNAME '=' STRING)*)? | 'COMMENT' '=' STRING ('TAG' QNAME '=' STRING (',' QNAME '=' STRING)*)? | 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*)
    | 'ALTER' 'EXTERNAL' 'ACCESS' 'INTEGRATION' ('IF' 'EXISTS')? ID 'UNSET' ('ALLOWED_NETWORK_RULES' | 'ALLOWED_API_AUTHENTICATION_INTEGRATIONS' | 'ALLOWED_AUTHENTICATION_SECRETS' | 'COMMENT' | 'TAG' QNAME) (',' ('ALLOWED_NETWORK_RULES' | 'ALLOWED_API_AUTHENTICATION_INTEGRATIONS' | 'ALLOWED_AUTHENTICATION_SECRETS' | 'COMMENT' | 'TAG' QNAME))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-external-table
// fix: any number of partition columns, as the example with three shows
// not converted (form): unbalanced "]"
//   ALTER EXTERNAL TABLE [ IF EXISTS ] <name> ADD FILES ( '<path>/[<filename>]' [ , '<path>/[<filename>'] ] )
alterExternalTable
    : 'ALTER' 'EXTERNAL' 'TABLE' ('IF' 'EXISTS')? QNAME 'REFRESH' STRING?
    | 'ALTER' 'EXTERNAL' 'TABLE' ('IF' 'EXISTS')? QNAME 'REMOVE' 'FILES' '(' STRING (',' STRING)? ')'
    | 'ALTER' 'EXTERNAL' 'TABLE' ('IF' 'EXISTS')? QNAME 'SET' 'AUTO_REFRESH' '=' ('TRUE' | 'FALSE')
    | 'ALTER' 'EXTERNAL' 'TABLE' QNAME ('IF' 'EXISTS')? 'ADD' 'PARTITION' '(' ID '=' STRING (',' ID '=' STRING)* ')' 'LOCATION' STRING
    | 'ALTER' 'EXTERNAL' 'TABLE' QNAME ('IF' 'EXISTS')? 'DROP' 'PARTITION' 'LOCATION' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-external-volume
alterExternalVolume
    : 'ALTER' 'EXTERNAL' 'VOLUME' ('IF' 'EXISTS')? ID 'ADD' 'STORAGE_LOCATION' '=' '(' 'NAME' '=' STRING alterExternalVolume_cloudProviderParams ')'
    | 'ALTER' 'EXTERNAL' 'VOLUME' ('IF' 'EXISTS')? ID 'REMOVE' 'STORAGE_LOCATION' STRING
    | 'ALTER' 'EXTERNAL' 'VOLUME' ('IF' 'EXISTS')? ID 'UPDATE' 'STORAGE_LOCATION' '=' STRING 'CREDENTIALS' '=' '(' 'AWS_KEY_ID' '=' STRING 'AWS_SECRET_KEY' '=' STRING ')'
    | 'ALTER' 'EXTERNAL' 'VOLUME' ('IF' 'EXISTS')? ID 'SET' 'ALLOW_WRITES' '=' ('TRUE' | 'FALSE')
    | 'ALTER' 'EXTERNAL' 'VOLUME' ('IF' 'EXISTS')? ID 'SET' 'COMMENT' '=' STRING
    ;
alterExternalVolume_cloudProviderParams
    : 'STORAGE_PROVIDER' '=' ('\'S3\'' | '\'S3CHINA\'' | '\'S3GOV\'') 'STORAGE_AWS_ROLE_ARN' '=' STRING 'STORAGE_BASE_URL' '=' STRING ('STORAGE_AWS_ACCESS_POINT_ARN' '=' STRING)? ('ENCRYPTION' '=' '(' (('TYPE' '=' '\'AWS_SSE_S3\'')? | ('TYPE' '=' '\'AWS_SSE_KMS\'' ('KMS_KEY_ID' '=' STRING)?)? | ('TYPE' '=' '\'NONE\'')?) ')')? ('USE_PRIVATELINK_ENDPOINT' '=' ('TRUE' | 'FALSE'))?
    | 'STORAGE_PROVIDER' '=' '\'GCS\'' 'STORAGE_BASE_URL' '=' STRING ('ENCRYPTION' '=' '(' (('TYPE' '=' '\'GCS_SSE_KMS\'')? ('KMS_KEY_ID' '=' STRING)? | ('TYPE' '=' '\'NONE\'')?) ')')?
    | 'STORAGE_PROVIDER' '=' '\'AZURE\'' 'AZURE_TENANT_ID' '=' STRING ('USE_PRIVATELINK_ENDPOINT' '=' ('TRUE' | 'FALSE'))? 'STORAGE_BASE_URL' '=' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-failover-group
// fix: OBJECT_TYPES takes the replication object types the parameters list, not any object type
// fix: ALLOWED_INTEGRATION_TYPES takes the integration types the parameters list
alterFailoverGroup
    : 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'RENAME' 'TO' ID
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'SET' ('OBJECT_TYPES' '=' ('ACCOUNT' 'PARAMETERS' | 'DATABASES' | 'EXTERNAL' 'VOLUMES' | 'INTEGRATIONS' | 'LISTINGS' | 'NETWORK' 'POLICIES' | 'PROFILES' | 'RESOURCE' 'MONITORS' | 'ROLES' | 'SHARES' | 'USERS' | 'WAREHOUSES') (',' ('ACCOUNT' 'PARAMETERS' | 'DATABASES' | 'EXTERNAL' 'VOLUMES' | 'INTEGRATIONS' | 'LISTINGS' | 'NETWORK' 'POLICIES' | 'PROFILES' | 'RESOURCE' 'MONITORS' | 'ROLES' | 'SHARES' | 'USERS' | 'WAREHOUSES'))* ('ALLOWED_DATABASES' '=' ID (',' ID)*)? ('ALLOWED_EXTERNAL_VOLUMES' '=' ID (',' ID)*)? ('ALLOWED_SHARES' '=' ID (',' ID)*)? | 'ALLOWED_DATABASES' '=' ID (',' ID)* ('ALLOWED_EXTERNAL_VOLUMES' '=' ID (',' ID)*)? ('ALLOWED_SHARES' '=' ID (',' ID)*)? | 'ALLOWED_EXTERNAL_VOLUMES' '=' ID (',' ID)* ('ALLOWED_SHARES' '=' ID (',' ID)*)? | 'ALLOWED_SHARES' '=' ID (',' ID)*)
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'SET' 'OBJECT_TYPES' '=' 'INTEGRATIONS' (',' ('ACCOUNT' 'PARAMETERS' | 'DATABASES' | 'EXTERNAL' 'VOLUMES' | 'INTEGRATIONS' | 'LISTINGS' | 'NETWORK' 'POLICIES' | 'PROFILES' | 'RESOURCE' 'MONITORS' | 'ROLES' | 'SHARES' | 'USERS' | 'WAREHOUSES'))* 'ALLOWED_INTEGRATION_TYPES' '=' ('SECURITY' 'INTEGRATIONS' | 'API' 'INTEGRATIONS' | 'STORAGE' 'INTEGRATIONS' | 'EXTERNAL' 'ACCESS' 'INTEGRATIONS' | 'NOTIFICATION' 'INTEGRATIONS') (',' ('SECURITY' 'INTEGRATIONS' | 'API' 'INTEGRATIONS' | 'STORAGE' 'INTEGRATIONS' | 'EXTERNAL' 'ACCESS' 'INTEGRATIONS' | 'NOTIFICATION' 'INTEGRATIONS'))*
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'SET' 'REPLICATION_SCHEDULE' '=' STRING
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'SET' 'OPTIMIZED_REFRESH' '=' ('TRUE' | 'FALSE')
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'SET' 'RPO_ASSURANCE' '=' ('TRUE' | 'FALSE')
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'SET' 'ERROR_INTEGRATION' '=' ID
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'SET' 'READ_ONLY' '=' 'TRUE'
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'UNSET' ('COMMENT' | 'REPLICATION_SCHEDULE' | 'OPTIMIZED_REFRESH' | 'RPO_ASSURANCE' | 'ERROR_INTEGRATION' | 'READ_ONLY') (',' ('COMMENT' | 'REPLICATION_SCHEDULE' | 'OPTIMIZED_REFRESH' | 'RPO_ASSURANCE' | 'ERROR_INTEGRATION' | 'READ_ONLY'))*
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'ADD' ID (',' ID)* 'TO' 'ALLOWED_DATABASES'
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'MOVE' 'DATABASES' ID (',' ID)* 'TO' 'FAILOVER' 'GROUP' ID
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'REMOVE' ID (',' ID)* 'FROM' 'ALLOWED_DATABASES'
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'ADD' ID (',' ID)* 'TO' 'ALLOWED_EXTERNAL_VOLUMES'
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'MOVE' 'EXTERNAL' 'VOLUMES' ID (',' ID)* 'TO' 'FAILOVER' 'GROUP' ID
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'REMOVE' ID (',' ID)* 'FROM' 'ALLOWED_EXTERNAL_VOLUMES'
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'ADD' ID (',' ID)* 'TO' 'ALLOWED_SHARES'
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'MOVE' 'SHARES' ID (',' ID)* 'TO' 'FAILOVER' 'GROUP' ID
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'REMOVE' ID (',' ID)* 'FROM' 'ALLOWED_SHARES'
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'ADD' ID '.' ID (',' ID '.' ID)* 'TO' 'ALLOWED_ACCOUNTS' ('IGNORE' 'EDITION' 'CHECK')?
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'REMOVE' ID '.' ID (',' ID '.' ID)* 'FROM' 'ALLOWED_ACCOUNTS'
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'REFRESH'
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'PRIMARY'
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'SUSPEND' 'IMMEDIATE'?
    | 'ALTER' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID 'RESUME'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-feature-policy
// fix: <type> is an object type in the plural, DATABASES or TABLES, as the page's parameters list
alterFeaturePolicy
    : 'ALTER' 'FEATURE' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' ('BLOCKED_OBJECT_TYPES_FOR_CREATION' '=' '(' (OBJECT_TYPES (',' OBJECT_TYPES)*)? ')' ('COMMENT' '=' STRING)? ('AS' STRING)? | 'COMMENT' '=' STRING ('AS' STRING)? | 'AS' STRING)
    | 'ALTER' 'FEATURE' 'POLICY' ('IF' 'EXISTS')? QNAME 'AS' STRING
    | 'ALTER' 'FEATURE' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' ('BLOCKED_OBJECT_TYPES_FOR_CREATION' | 'COMMENT' | 'POLICY_DEFINITION') (',' ('BLOCKED_OBJECT_TYPES_FOR_CREATION' | 'COMMENT' | 'POLICY_DEFINITION'))*
    | 'ALTER' 'FEATURE' 'POLICY' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'FEATURE' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'FEATURE' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-file-format
// fix: read: each -- If TYPE = X section of formatTypeOptions is the options of TYPE = X alone, and TYPE = { ... } [ formatTypeOptions ] chooses a TYPE with its own options
alterFileFormat
    : 'ALTER' 'FILE' 'FORMAT' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'FILE' 'FORMAT' ('IF' 'EXISTS')? QNAME 'SET' ((alterFileFormat_formatTypeOptions_CSV | alterFileFormat_formatTypeOptions_JSON | alterFileFormat_formatTypeOptions_AVRO | alterFileFormat_formatTypeOptions_ORC | alterFileFormat_formatTypeOptions_PARQUET | alterFileFormat_formatTypeOptions_XML) ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    ;
alterFileFormat_formatTypeOptions_CSV
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'RECORD_DELIMITER' '=' (STRING | 'NONE') | 'FIELD_DELIMITER' '=' (STRING | 'NONE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'FILE_EXTENSION' '=' STRING | 'PARSE_HEADER' '=' ('TRUE' | 'FALSE') | 'SKIP_HEADER' '=' NUMBER | 'SKIP_BLANK_LINES' '=' ('TRUE' | 'FALSE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'ESCAPE' '=' (STRING | 'NONE') | 'ESCAPE_UNENCLOSED_FIELD' '=' (STRING | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'FIELD_OPTIONALLY_ENCLOSED_BY' '=' (STRING | 'NONE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')' | 'ERROR_ON_COLUMN_COUNT_MISMATCH' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'EMPTY_FIELD_AS_NULL' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE') | 'ENCODING' '=' (STRING | 'UTF8')) (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'RECORD_DELIMITER' '=' (STRING | 'NONE') | 'FIELD_DELIMITER' '=' (STRING | 'NONE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'FILE_EXTENSION' '=' STRING | 'PARSE_HEADER' '=' ('TRUE' | 'FALSE') | 'SKIP_HEADER' '=' NUMBER | 'SKIP_BLANK_LINES' '=' ('TRUE' | 'FALSE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'ESCAPE' '=' (STRING | 'NONE') | 'ESCAPE_UNENCLOSED_FIELD' '=' (STRING | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'FIELD_OPTIONALLY_ENCLOSED_BY' '=' (STRING | 'NONE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')' | 'ERROR_ON_COLUMN_COUNT_MISMATCH' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'EMPTY_FIELD_AS_NULL' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE') | 'ENCODING' '=' (STRING | 'UTF8')))*
    ;
alterFileFormat_formatTypeOptions_JSON
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')' | 'FILE_EXTENSION' '=' STRING | 'ENABLE_OCTAL' '=' ('TRUE' | 'FALSE') | 'ALLOW_DUPLICATE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ARRAY' '=' ('TRUE' | 'FALSE') | 'STRIP_NULL_VALUES' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'IGNORE_UTF8_ERRORS' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE')) (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')' | 'FILE_EXTENSION' '=' STRING | 'ENABLE_OCTAL' '=' ('TRUE' | 'FALSE') | 'ALLOW_DUPLICATE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ARRAY' '=' ('TRUE' | 'FALSE') | 'STRIP_NULL_VALUES' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'IGNORE_UTF8_ERRORS' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE')))*
    ;
alterFileFormat_formatTypeOptions_AVRO
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')') (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')'))*
    ;
alterFileFormat_formatTypeOptions_ORC
    : ('TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')') (','? ('TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')'))*
    ;
alterFileFormat_formatTypeOptions_PARQUET
    : ('COMPRESSION' '=' ('AUTO' | 'LZO' | 'SNAPPY' | 'NONE') | 'SNAPPY_COMPRESSION' '=' ('TRUE' | 'FALSE') | 'BINARY_AS_TEXT' '=' ('TRUE' | 'FALSE') | 'USE_LOGICAL_TYPE' '=' ('TRUE' | 'FALSE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'USE_VECTORIZED_SCANNER' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')') (','? ('COMPRESSION' '=' ('AUTO' | 'LZO' | 'SNAPPY' | 'NONE') | 'SNAPPY_COMPRESSION' '=' ('TRUE' | 'FALSE') | 'BINARY_AS_TEXT' '=' ('TRUE' | 'FALSE') | 'USE_LOGICAL_TYPE' '=' ('TRUE' | 'FALSE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'USE_VECTORIZED_SCANNER' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')'))*
    ;
alterFileFormat_formatTypeOptions_XML
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'IGNORE_UTF8_ERRORS' '=' ('TRUE' | 'FALSE') | 'PRESERVE_SPACE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ELEMENT' '=' ('TRUE' | 'FALSE') | 'DISABLE_AUTO_CONVERT' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE')) (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'IGNORE_UTF8_ERRORS' '=' ('TRUE' | 'FALSE') | 'PRESERVE_SPACE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ELEMENT' '=' ('TRUE' | 'FALSE') | 'DISABLE_AUTO_CONVERT' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE')))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-function
alterFunction
    : 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'RENAME' 'TO' QNAME
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'SET' 'SECURE'
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'UNSET' ('SECURE' | 'LOG_LEVEL' | 'TRACE_LEVEL' | 'COMMENT')
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'SET' ('LOG_LEVEL' '=' STRING ('TRACE_LEVEL' '=' STRING)? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')')? ('COMMENT' '=' STRING)? | 'TRACE_LEVEL' '=' STRING ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')')? ('COMMENT' '=' STRING)? | 'EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')' ('SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')')? ('COMMENT' '=' STRING)? | 'SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')' ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'SET' ('LOG_LEVEL' '=' STRING ('TRACE_LEVEL' '=' STRING)? ('COMMENT' '=' STRING)? | 'TRACE_LEVEL' '=' STRING ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'SET' 'API_INTEGRATION' '=' ID
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'SET' 'HEADERS' '=' '(' (STRING '=' STRING (',' STRING '=' STRING)*)? ')'
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'SET' 'CONTEXT_HEADERS' '=' '(' (ID (',' ID)*)? ')'
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'SET' 'MAX_BATCH_ROWS' '=' NUMBER
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'SET' 'COMPRESSION' '=' ID
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'SET' ('REQUEST_TRANSLATOR' | 'RESPONSE_TRANSLATOR') '=' ID
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'UNSET' ('COMMENT' | 'HEADERS' | 'CONTEXT_HEADERS' | 'MAX_BATCH_ROWS' | 'COMPRESSION' | 'SECURE' | 'REQUEST_TRANSLATOR' | 'RESPONSE_TRANSLATOR')
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-function-dmf
alterFunctionDmf
    : 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' 'TABLE' '(' DATA_TYPE (',' DATA_TYPE)* ')' (',' 'TABLE' '(' DATA_TYPE (',' DATA_TYPE)* ')')? ')' 'RENAME' 'TO' QNAME
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' 'TABLE' '(' DATA_TYPE (',' DATA_TYPE)* ')' (',' 'TABLE' '(' DATA_TYPE (',' DATA_TYPE)* ')')? ')' 'SET' 'SECURE'
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' 'TABLE' '(' DATA_TYPE (',' DATA_TYPE)* ')' (',' 'TABLE' '(' DATA_TYPE (',' DATA_TYPE)* ')')? ')' 'UNSET' 'SECURE'
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' 'TABLE' '(' DATA_TYPE (',' DATA_TYPE)* ')' (',' 'TABLE' '(' DATA_TYPE (',' DATA_TYPE)* ')')? ')' 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' 'TABLE' '(' DATA_TYPE (',' DATA_TYPE)* ')' (',' 'TABLE' '(' DATA_TYPE (',' DATA_TYPE)* ')')? ')' 'UNSET' 'COMMENT'
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' 'TABLE' '(' DATA_TYPE (',' DATA_TYPE)* ')' (',' 'TABLE' '(' DATA_TYPE (',' DATA_TYPE)* ')')? ')' 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' 'TABLE' '(' DATA_TYPE (',' DATA_TYPE)* ')' (',' 'TABLE' '(' DATA_TYPE (',' DATA_TYPE)* ')')? ')' 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-function-spcs
alterFunctionSpcs
    : 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'RENAME' 'TO' QNAME
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'SET' 'CONTEXT_HEADERS' '=' '(' ID (',' ID)* ')'
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'SET' 'MAX_BATCH_ROWS' '=' NUMBER
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'SET' 'MAX_BATCH_RETRIES' '=' NUMBER
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'SET' 'ON_BATCH_FAILURE' '=' ('ABORT' | 'RETURN_NULL')
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'SET' 'BATCH_TIMEOUT_SECS' '=' NUMBER
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'SET' 'SERVICE' '=' STRING 'ENDPOINT' '=' STRING
    | 'ALTER' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'UNSET' ('CONTEXT_HEADERS' | 'MAX_BATCH_ROWS' | 'MAX_BATCH_RETRIES' | 'ON_BATCH_FAILURE' | 'BATCH_TIMEOUT_SECS' | 'COMMENT')
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-gateway
alterGateway
    : 'ALTER' 'GATEWAY' ('IF' 'EXISTS')? QNAME 'FROM' 'SPECIFICATION' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-git-repository
alterGitRepository
    : 'ALTER' 'GIT' 'REPOSITORY' ('IF' 'EXISTS')? QNAME 'SET' ('GIT_CREDENTIALS' '=' QNAME ('API_INTEGRATION' '=' ID)? ('COMMENT' '=' STRING)? | 'API_INTEGRATION' '=' ID ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'GIT' 'REPOSITORY' ('IF' 'EXISTS')? QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'GIT' 'REPOSITORY' ('IF' 'EXISTS')? QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'GIT' 'REPOSITORY' ('IF' 'EXISTS')? QNAME 'UNSET' ('GIT_CREDENTIALS' | 'COMMENT') (',' ('GIT_CREDENTIALS' | 'COMMENT'))*
    | 'ALTER' 'GIT' 'REPOSITORY' ('IF' 'EXISTS')? QNAME 'FETCH'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-iceberg-table
// fix: the governance form names no table; the page's other forms take <table_name>
// fix: tableColumnAction: ALTER and MODIFY both take the column actions after them
// fix: the column actions are a comma-separated list of actions, each one of those the page lists, as alter-table's are
// fix: a search method with its targets, as the parameters write it
// fix: read: <constraint_properties> is the properties CREATE TABLE's constraint page lists
// fix: read: a foreign key's REFERENCES takes MATCH and ON UPDATE/DELETE, the properties the constraint page gives foreign keys
alterIcebergTable
    : 'ALTER' 'ICEBERG' 'TABLE' ('IF' 'EXISTS')? QNAME (alterIcebergTable_clusteringAction | alterIcebergTable_tableColumnAction)
    | 'ALTER' 'ICEBERG' 'TABLE' ('IF' 'EXISTS')? QNAME 'SET' ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') ('CATALOG_SYNC' '=' STRING)? ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('TARGET_FILE_SIZE' '=' ('AUTO' | '16' 'MB' | '32' 'MB' | '64' 'MB' | '128' 'MB'))? ('PATH_LAYOUT' '=' ('FLAT' | 'HIERARCHICAL'))? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('LOG_EVENT_LEVEL' '=' ('ERROR' | 'WARN' | 'INFO' | 'DEBUG'))? ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('MATCH_BY_COLUMN_NAME' '=' ('CASE_INSENSITIVE' | 'CASE_SENSITIVE'))? ('ALLOW_V3_EXPRESSION_DEFAULTS' '=' ('TRUE' | 'FALSE'))? ('TIME_TRAVEL_TIMESTAMP_SOURCE' '=' ('\'SNOWFLAKE_METADATA\'' | '\'ICEBERG_METADATA\''))? | 'CATALOG_SYNC' '=' STRING ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('TARGET_FILE_SIZE' '=' ('AUTO' | '16' 'MB' | '32' 'MB' | '64' 'MB' | '128' 'MB'))? ('PATH_LAYOUT' '=' ('FLAT' | 'HIERARCHICAL'))? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('LOG_EVENT_LEVEL' '=' ('ERROR' | 'WARN' | 'INFO' | 'DEBUG'))? ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('MATCH_BY_COLUMN_NAME' '=' ('CASE_INSENSITIVE' | 'CASE_SENSITIVE'))? ('ALLOW_V3_EXPRESSION_DEFAULTS' '=' ('TRUE' | 'FALSE'))? ('TIME_TRAVEL_TIMESTAMP_SOURCE' '=' ('\'SNOWFLAKE_METADATA\'' | '\'ICEBERG_METADATA\''))? | 'DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('TARGET_FILE_SIZE' '=' ('AUTO' | '16' 'MB' | '32' 'MB' | '64' 'MB' | '128' 'MB'))? ('PATH_LAYOUT' '=' ('FLAT' | 'HIERARCHICAL'))? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('LOG_EVENT_LEVEL' '=' ('ERROR' | 'WARN' | 'INFO' | 'DEBUG'))? ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('MATCH_BY_COLUMN_NAME' '=' ('CASE_INSENSITIVE' | 'CASE_SENSITIVE'))? ('ALLOW_V3_EXPRESSION_DEFAULTS' '=' ('TRUE' | 'FALSE'))? ('TIME_TRAVEL_TIMESTAMP_SOURCE' '=' ('\'SNOWFLAKE_METADATA\'' | '\'ICEBERG_METADATA\''))? | 'AUTO_REFRESH' '=' ('TRUE' | 'FALSE') ('TARGET_FILE_SIZE' '=' ('AUTO' | '16' 'MB' | '32' 'MB' | '64' 'MB' | '128' 'MB'))? ('PATH_LAYOUT' '=' ('FLAT' | 'HIERARCHICAL'))? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('LOG_EVENT_LEVEL' '=' ('ERROR' | 'WARN' | 'INFO' | 'DEBUG'))? ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('MATCH_BY_COLUMN_NAME' '=' ('CASE_INSENSITIVE' | 'CASE_SENSITIVE'))? ('ALLOW_V3_EXPRESSION_DEFAULTS' '=' ('TRUE' | 'FALSE'))? ('TIME_TRAVEL_TIMESTAMP_SOURCE' '=' ('\'SNOWFLAKE_METADATA\'' | '\'ICEBERG_METADATA\''))? | 'TARGET_FILE_SIZE' '=' ('AUTO' | '16' 'MB' | '32' 'MB' | '64' 'MB' | '128' 'MB') ('PATH_LAYOUT' '=' ('FLAT' | 'HIERARCHICAL'))? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('LOG_EVENT_LEVEL' '=' ('ERROR' | 'WARN' | 'INFO' | 'DEBUG'))? ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('MATCH_BY_COLUMN_NAME' '=' ('CASE_INSENSITIVE' | 'CASE_SENSITIVE'))? ('ALLOW_V3_EXPRESSION_DEFAULTS' '=' ('TRUE' | 'FALSE'))? ('TIME_TRAVEL_TIMESTAMP_SOURCE' '=' ('\'SNOWFLAKE_METADATA\'' | '\'ICEBERG_METADATA\''))? | 'PATH_LAYOUT' '=' ('FLAT' | 'HIERARCHICAL') ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('LOG_EVENT_LEVEL' '=' ('ERROR' | 'WARN' | 'INFO' | 'DEBUG'))? ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('MATCH_BY_COLUMN_NAME' '=' ('CASE_INSENSITIVE' | 'CASE_SENSITIVE'))? ('ALLOW_V3_EXPRESSION_DEFAULTS' '=' ('TRUE' | 'FALSE'))? ('TIME_TRAVEL_TIMESTAMP_SOURCE' '=' ('\'SNOWFLAKE_METADATA\'' | '\'ICEBERG_METADATA\''))? | 'STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED') ('CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('LOG_EVENT_LEVEL' '=' ('ERROR' | 'WARN' | 'INFO' | 'DEBUG'))? ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('MATCH_BY_COLUMN_NAME' '=' ('CASE_INSENSITIVE' | 'CASE_SENSITIVE'))? ('ALLOW_V3_EXPRESSION_DEFAULTS' '=' ('TRUE' | 'FALSE'))? ('TIME_TRAVEL_TIMESTAMP_SOURCE' '=' ('\'SNOWFLAKE_METADATA\'' | '\'ICEBERG_METADATA\''))? | 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')' ('LOG_EVENT_LEVEL' '=' ('ERROR' | 'WARN' | 'INFO' | 'DEBUG'))? ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('MATCH_BY_COLUMN_NAME' '=' ('CASE_INSENSITIVE' | 'CASE_SENSITIVE'))? ('ALLOW_V3_EXPRESSION_DEFAULTS' '=' ('TRUE' | 'FALSE'))? ('TIME_TRAVEL_TIMESTAMP_SOURCE' '=' ('\'SNOWFLAKE_METADATA\'' | '\'ICEBERG_METADATA\''))? | 'LOG_EVENT_LEVEL' '=' ('ERROR' | 'WARN' | 'INFO' | 'DEBUG') ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('MATCH_BY_COLUMN_NAME' '=' ('CASE_INSENSITIVE' | 'CASE_SENSITIVE'))? ('ALLOW_V3_EXPRESSION_DEFAULTS' '=' ('TRUE' | 'FALSE'))? ('TIME_TRAVEL_TIMESTAMP_SOURCE' '=' ('\'SNOWFLAKE_METADATA\'' | '\'ICEBERG_METADATA\''))? | 'ERROR_LOGGING' '=' ('TRUE' | 'FALSE') ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('MATCH_BY_COLUMN_NAME' '=' ('CASE_INSENSITIVE' | 'CASE_SENSITIVE'))? ('ALLOW_V3_EXPRESSION_DEFAULTS' '=' ('TRUE' | 'FALSE'))? ('TIME_TRAVEL_TIMESTAMP_SOURCE' '=' ('\'SNOWFLAKE_METADATA\'' | '\'ICEBERG_METADATA\''))? | 'ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE') ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('MATCH_BY_COLUMN_NAME' '=' ('CASE_INSENSITIVE' | 'CASE_SENSITIVE'))? ('ALLOW_V3_EXPRESSION_DEFAULTS' '=' ('TRUE' | 'FALSE'))? ('TIME_TRAVEL_TIMESTAMP_SOURCE' '=' ('\'SNOWFLAKE_METADATA\'' | '\'ICEBERG_METADATA\''))? | 'ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\'') ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('MATCH_BY_COLUMN_NAME' '=' ('CASE_INSENSITIVE' | 'CASE_SENSITIVE'))? ('ALLOW_V3_EXPRESSION_DEFAULTS' '=' ('TRUE' | 'FALSE'))? ('TIME_TRAVEL_TIMESTAMP_SOURCE' '=' ('\'SNOWFLAKE_METADATA\'' | '\'ICEBERG_METADATA\''))? | 'ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE') ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('MATCH_BY_COLUMN_NAME' '=' ('CASE_INSENSITIVE' | 'CASE_SENSITIVE'))? ('ALLOW_V3_EXPRESSION_DEFAULTS' '=' ('TRUE' | 'FALSE'))? ('TIME_TRAVEL_TIMESTAMP_SOURCE' '=' ('\'SNOWFLAKE_METADATA\'' | '\'ICEBERG_METADATA\''))? | 'ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE') ('MATCH_BY_COLUMN_NAME' '=' ('CASE_INSENSITIVE' | 'CASE_SENSITIVE'))? ('ALLOW_V3_EXPRESSION_DEFAULTS' '=' ('TRUE' | 'FALSE'))? ('TIME_TRAVEL_TIMESTAMP_SOURCE' '=' ('\'SNOWFLAKE_METADATA\'' | '\'ICEBERG_METADATA\''))? | 'MATCH_BY_COLUMN_NAME' '=' ('CASE_INSENSITIVE' | 'CASE_SENSITIVE') ('ALLOW_V3_EXPRESSION_DEFAULTS' '=' ('TRUE' | 'FALSE'))? ('TIME_TRAVEL_TIMESTAMP_SOURCE' '=' ('\'SNOWFLAKE_METADATA\'' | '\'ICEBERG_METADATA\''))? | 'ALLOW_V3_EXPRESSION_DEFAULTS' '=' ('TRUE' | 'FALSE') ('TIME_TRAVEL_TIMESTAMP_SOURCE' '=' ('\'SNOWFLAKE_METADATA\'' | '\'ICEBERG_METADATA\''))? | 'TIME_TRAVEL_TIMESTAMP_SOURCE' '=' ('\'SNOWFLAKE_METADATA\'' | '\'ICEBERG_METADATA\''))
    | 'ALTER' 'ICEBERG' 'TABLE' ('IF' 'EXISTS')? QNAME 'UNSET' ('REPLACE_INVALID_CHARACTERS' | 'PATH_LAYOUT' | 'LOG_EVENT_LEVEL' | 'ERROR_LOGGING' | 'ENABLE_DATA_COMPACTION' | 'ICEBERG_MERGE_ON_READ_BEHAVIOR' | 'ENABLE_ICEBERG_MERGE_ON_READ' | 'ENABLE_SCHEMA_EVOLUTION' | 'ALLOW_V3_EXPRESSION_DEFAULTS') (',' ('REPLACE_INVALID_CHARACTERS' | 'PATH_LAYOUT' | 'LOG_EVENT_LEVEL' | 'ERROR_LOGGING' | 'ENABLE_DATA_COMPACTION' | 'ICEBERG_MERGE_ON_READ_BEHAVIOR' | 'ENABLE_ICEBERG_MERGE_ON_READ' | 'ENABLE_SCHEMA_EVOLUTION' | 'ALLOW_V3_EXPRESSION_DEFAULTS'))*
    | 'ALTER' 'ICEBERG' 'TABLE' ('IF' 'EXISTS')? QNAME 'UPGRADE' 'TO' 'ICEBERG' 'VERSION' '3'
    | 'ALTER' 'ICEBERG' 'TABLE' ('IF' 'EXISTS')? QNAME alterIcebergTable_dataGovnPolicyTagAction
    | 'ALTER' 'ICEBERG' 'TABLE' ('IF' 'EXISTS')? QNAME alterIcebergTable_searchOptimizationAction
    ;
alterIcebergTable_clusteringAction
    : ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')' | ('SUSPEND' | 'RESUME') 'RECLUSTER' | 'DROP' 'CLUSTERING' 'KEY')
    ;
alterIcebergTable_tableColumnAction
    : ('ADD' 'COLUMN'? ('IF' 'NOT' 'EXISTS')? ID DATA_TYPE ('DEFAULT' EXPR)? alterIcebergTable_inlineConstraint? ('COLLATE' STRING)? | 'RENAME' 'COLUMN' ID 'TO' ID | ('ALTER' | 'MODIFY') ('(' 'COLUMN'? ID (('SET'? 'NOT' 'NULL' | 'DROP' 'NOT' 'NULL') | (('SET' 'DATA')? 'TYPE')? DATA_TYPE | 'COMMENT' STRING | 'UNSET' 'COMMENT' | 'SET' 'WRITE' 'DEFAULT' EXPR | 'DROP' 'WRITE' 'DEFAULT') (',' 'COLUMN'? ID (('SET'? 'NOT' 'NULL' | 'DROP' 'NOT' 'NULL') | (('SET' 'DATA')? 'TYPE')? DATA_TYPE | 'COMMENT' STRING | 'UNSET' 'COMMENT' | 'SET' 'WRITE' 'DEFAULT' EXPR | 'DROP' 'WRITE' 'DEFAULT'))* ')' | 'COLUMN'? ID (('SET'? 'NOT' 'NULL' | 'DROP' 'NOT' 'NULL') | (('SET' 'DATA')? 'TYPE')? DATA_TYPE | 'COMMENT' STRING | 'UNSET' 'COMMENT' | 'SET' 'WRITE' 'DEFAULT' EXPR | 'DROP' 'WRITE' 'DEFAULT') (',' 'COLUMN'? ID (('SET'? 'NOT' 'NULL' | 'DROP' 'NOT' 'NULL') | (('SET' 'DATA')? 'TYPE')? DATA_TYPE | 'COMMENT' STRING | 'UNSET' 'COMMENT' | 'SET' 'WRITE' 'DEFAULT' EXPR | 'DROP' 'WRITE' 'DEFAULT'))*) | 'DROP' 'COLUMN'? ('IF' 'EXISTS')? ID (',' ID)*)
    ;
alterIcebergTable_inlineConstraint
    : ('NOT' 'NULL')? ('CONSTRAINT' ID)? ('UNIQUE' | 'PRIMARY' 'KEY' | ('FOREIGN' 'KEY')? 'REFERENCES' QNAME ('(' ID ')')? ('MATCH' ('FULL' | 'SIMPLE' | 'PARTIAL'))? ('ON' ('UPDATE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))? ('DELETE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))?)? | 'CHECK' '(' EXPR ')') ('NOT'? 'ENFORCED')? ('NOT'? 'DEFERRABLE')? ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? (('RELY' | 'NORELY'))?
    ;
alterIcebergTable_dataGovnPolicyTagAction
    : ('SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)* | 'UNSET' 'TAG' QNAME (',' QNAME)*)
    | ('ADD' 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')' | 'DROP' 'ROW' 'ACCESS' 'POLICY' QNAME | 'DROP' 'ROW' 'ACCESS' 'POLICY' QNAME ',' 'ADD' 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')' | 'DROP' 'ALL' 'ROW' 'ACCESS' 'POLICIES')
    | ('SET' 'AGGREGATION' 'POLICY' QNAME ('ENTITY' 'KEY' '(' ID (',' ID)* ')')? 'FORCE'? | 'UNSET' 'AGGREGATION' 'POLICY')
    | ('SET' 'JOIN' 'POLICY' QNAME 'FORCE'? | 'UNSET' 'JOIN' 'POLICY')
    | 'ADD' 'COLUMN'? ('IF' 'NOT' 'EXISTS')? ID DATA_TYPE ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | (('ALTER' | 'MODIFY') 'COLUMN'? ID 'SET' 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')? 'FORCE'? | ('ALTER' | 'MODIFY') 'COLUMN'? ID 'UNSET' 'MASKING' 'POLICY')
    | (('ALTER' | 'MODIFY') 'COLUMN'? ID 'SET' 'PROJECTION' 'POLICY' QNAME 'FORCE'? | ('ALTER' | 'MODIFY') 'COLUMN'? ID 'UNSET' 'PROJECTION' 'POLICY')
    | ('ALTER' | 'MODIFY') 'COLUMN'? ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)* ',' 'COLUMN'? ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | ('ALTER' | 'MODIFY') 'COLUMN'? ID 'UNSET' 'TAG' QNAME (',' QNAME)* ',' 'COLUMN'? ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;
alterIcebergTable_searchOptimizationAction
    : ('ADD' 'SEARCH' 'OPTIMIZATION' ('ON' ('FULL_TEXT' | 'EQUALITY' | 'SUBSTRING' | 'GEO') '(' ID (',' ID)* (',' 'ANALYZER' '=>' STRING)? ')' (',' ('FULL_TEXT' | 'EQUALITY' | 'SUBSTRING' | 'GEO') '(' ID (',' ID)* (',' 'ANALYZER' '=>' STRING)? ')')*)? | 'DROP' 'SEARCH' 'OPTIMIZATION' ('ON' (('FULL_TEXT' | 'EQUALITY' | 'SUBSTRING' | 'GEO') '(' ID (',' ID)* (',' 'ANALYZER' '=>' STRING)? ')' | ID | EXPR) (',' (('FULL_TEXT' | 'EQUALITY' | 'SUBSTRING' | 'GEO') '(' ID (',' ID)* (',' 'ANALYZER' '=>' STRING)? ')' | ID | EXPR))*)?)
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-iceberg-table-alter-column-set-data-type
alterIcebergTableAlterColumnSetDataType
    : 'ALTER' 'ICEBERG' 'TABLE' ('IF' 'EXISTS')? QNAME 'ALTER' 'COLUMN' ID 'SET' 'DATA' 'TYPE' DATA_TYPE
    | 'ALTER' 'ICEBERG' 'TABLE' ('IF' 'EXISTS')? QNAME 'ALTER' 'COLUMN' ID 'SET' 'DATA' 'TYPE' DATA_TYPE 'RENAME' 'FIELDS'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-iceberg-table-convert-to-managed
alterIcebergTableConvertToManaged
    : 'ALTER' 'ICEBERG' 'TABLE' ('IF' 'EXISTS')? QNAME 'CONVERT' 'TO' 'MANAGED' (('FILE_PATH' | 'BASE_LOCATION') '=' STRING)? ('DROP_FROM_SOURCE_CATALOG' '=' ('TRUE' | 'FALSE'))? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-iceberg-table-optimize
alterIcebergTableOptimize
    : 'ALTER' 'ICEBERG' 'TABLE' ('IF' 'EXISTS')? QNAME 'OPTIMIZE' '(' 'REMOVE_POSITION_DELETES' ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-iceberg-table-partition-evolution
// fix: the page says oldPartitionExpression uses partitionExpression's syntax
// fix: partitionExpression is defined under Parameters, outside the Syntax block; newPartitionExpression uses its syntax
alterIcebergTablePartitionEvolution
    : 'ALTER' 'ICEBERG' 'TABLE' ('IF' 'EXISTS')? QNAME 'ADD' 'PARTITION' 'BY' '(' alterIcebergTablePartitionEvolution_partitionExpression ')'
    | 'ALTER' 'ICEBERG' 'TABLE' ('IF' 'EXISTS')? QNAME 'DROP' 'PARTITION' 'BY' '(' alterIcebergTablePartitionEvolution_partitionExpression ')'
    | 'ALTER' 'ICEBERG' 'TABLE' ('IF' 'EXISTS')? QNAME 'REPLACE' 'PARTITION' 'BY' '(' alterIcebergTablePartitionEvolution_partitionExpression ')' 'WITH' '(' alterIcebergTablePartitionEvolution_partitionExpression ')'
    ;
alterIcebergTablePartitionEvolution_partitionExpression
    : ID
    | 'BUCKET' '(' NUMBER ',' ID ')'
    | 'TRUNCATE' '(' NUMBER ',' ID ')'
    | 'YEAR' '(' ID ')'
    | 'MONTH' '(' ID ')'
    | 'DAY' '(' ID ')'
    | 'HOUR' '(' ID ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-iceberg-table-refresh
alterIcebergTableRefresh
    : 'ALTER' 'ICEBERG' 'TABLE' ('IF' 'EXISTS')? QNAME 'REFRESH' STRING?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-iceberg-table-version-tag
alterIcebergTableVersionTag
    : 'ALTER' 'ICEBERG' 'TABLE' ('IF' 'EXISTS')? QNAME 'CREATE' ('OR' 'REPLACE')? 'VERSION_TAG' QNAME ('AS' 'OF' ('VERSION' ID | 'TIMESTAMP' NUMBER) alterIcebergTableVersionTag_retainClause?)?
    | 'ALTER' 'ICEBERG' 'TABLE' ('IF' 'EXISTS')? QNAME 'CREATE' 'VERSION_TAG' ('IF' 'NOT' 'EXISTS')? QNAME ('AS' 'OF' ('VERSION' ID | 'TIMESTAMP' NUMBER) alterIcebergTableVersionTag_retainClause?)?
    | 'ALTER' 'ICEBERG' 'TABLE' ('IF' 'EXISTS')? QNAME 'ALTER' 'VERSION_TAG' QNAME 'AS' 'OF' ('VERSION' ID | 'TIMESTAMP' NUMBER)
    | 'ALTER' 'ICEBERG' 'TABLE' ('IF' 'EXISTS')? QNAME 'REPLACE' 'VERSION_TAG' QNAME 'AS' 'OF' ('VERSION' ID | 'TIMESTAMP' NUMBER) alterIcebergTableVersionTag_retainClause?
    | 'ALTER' 'ICEBERG' 'TABLE' ('IF' 'EXISTS')? QNAME 'DROP' 'VERSION_TAG' ('IF' 'EXISTS')? QNAME
    ;
alterIcebergTableVersionTag_retainClause
    : 'RETAIN' NUMBER ('DAYS' | 'HOURS' | 'MINUTES')
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-image-repository
alterImageRepository
    : 'ALTER' 'IMAGE' 'REPOSITORY' ('IF' 'EXISTS')? QNAME 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'IMAGE' 'REPOSITORY' ('IF' 'EXISTS')? QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'IMAGE' 'REPOSITORY' ('IF' 'EXISTS')? QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'IMAGE' 'REPOSITORY' ('IF' 'EXISTS')? QNAME 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-ingress-policy
alterIngressPolicy
    : 'ALTER' 'INGRESS' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'INGRESS' 'POLICY' ('IF' 'EXISTS')? QNAME 'AS' STRING
    | 'ALTER' 'INGRESS' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'COMMENT' '=' STRING 'AS' STRING
    | 'ALTER' 'INGRESS' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-integration
// not converted (form): <integration_type> stands for a type's keyword, whose page gives its syntax
//   ALTER [ <integration_type> ] INTEGRATION <object_name> <actions>

// https://docs.snowflake.com/en/sql-reference/sql/alter-join-policy
alterJoinPolicy
    : 'ALTER' 'JOIN' 'POLICY' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'JOIN' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'BODY' '->' EXPR
    | 'ALTER' 'JOIN' 'POLICY' QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'JOIN' 'POLICY' QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'JOIN' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'JOIN' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-listing
alterListing
    : 'ALTER' 'LISTING' ('IF' 'EXISTS')? ID (('PUBLISH' | 'UNPUBLISH' | 'REVIEW'))?
    | 'ALTER' 'LISTING' ('IF' 'EXISTS')? ID 'AS' STRING ('PUBLISH' '=' ('TRUE' | 'FALSE'))? ('REVIEW' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)?
    | 'ALTER' 'LISTING' ID 'ADD' 'VERSION' (('IF' 'NOT' 'EXISTS')? ID)? 'FROM' ID ('COMMENT' '=' STRING)?
    | 'ALTER' 'LISTING' ('IF' 'EXISTS')? ID ('ADD' | 'REMOVE') 'TARGETS' ID
    | 'ALTER' 'LISTING' ('IF' 'EXISTS')? ID 'RENAME' 'TO' ID
    | 'ALTER' 'LISTING' ('IF' 'EXISTS')? ID 'SET' 'COMMENT' '=' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-maintenance-policy
alterMaintenancePolicy
    : 'ALTER' 'MAINTENANCE' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' ('SCHEDULE' '=' STRING ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'MAINTENANCE' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' 'COMMENT'
    | 'ALTER' 'MAINTENANCE' 'POLICY' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-masking-policy
alterMaskingPolicy
    : 'ALTER' 'MASKING' 'POLICY' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'MASKING' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'BODY' '->' EXPR
    | 'ALTER' 'MASKING' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'MASKING' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'MASKING' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'MASKING' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-materialized-view
// fix: UNSET takes SECURE, CONTACT or COMMENT, one at a time, as the parameters say; the | after SECURE is missing
// fix: the view's name is missing; the parameters name it
// fix: the view's name is missing; the parameters name it
alterMaterializedView
    : 'ALTER' 'MATERIALIZED' 'VIEW' QNAME ('RENAME' 'TO' QNAME | 'CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')' | 'DROP' 'CLUSTERING' 'KEY' | 'SUSPEND' 'RECLUSTER' | 'RESUME' 'RECLUSTER' | 'SUSPEND' | 'RESUME' | 'SET' ('SECURE' ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('COMMENT' '=' STRING)? | 'CONTACT' ID '=' ID (',' ID '=' ID)* ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING) | 'UNSET' ('SECURE' | 'CONTACT' ID | 'COMMENT'))
    | 'ALTER' 'MATERIALIZED' 'VIEW' QNAME 'SET' 'DATA_METRIC_SCHEDULE' '=' STRING
    | 'ALTER' 'MATERIALIZED' 'VIEW' QNAME 'UNSET' 'DATA_METRIC_SCHEDULE'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-model
alterModel
    : 'ALTER' 'MODEL' ('IF' 'EXISTS')? QNAME 'SET' ('COMMENT' '=' STRING ('DEFAULT_VERSION' '=' STRING)? | 'DEFAULT_VERSION' '=' STRING)
    | 'ALTER' 'MODEL' ('IF' 'EXISTS')? QNAME 'SET' 'TAG' QNAME '=' STRING
    | 'ALTER' 'MODEL' ('IF' 'EXISTS')? QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'MODEL' ('IF' 'EXISTS')? QNAME 'VERSION' ID 'SET' 'ALIAS' '=' ID
    | 'ALTER' 'MODEL' ('IF' 'EXISTS')? QNAME 'VERSION' ID 'UNSET' 'ALIAS'
    | 'ALTER' 'MODEL' QNAME 'RENAME' 'TO' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-model-add-version
alterModelAddVersion
    : 'ALTER' 'MODEL' ('IF' 'EXISTS')? QNAME 'ADD' 'VERSION' ID 'FROM' 'MODEL' QNAME ('VERSION' ID)?
    | 'ALTER' 'MODEL' ('IF' 'EXISTS')? QNAME 'ADD' 'VERSION' ID 'FROM' alterModelAddVersion_internalStage
    ;
alterModelAddVersion_internalStage
    : STAGE
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-model-drop-version
alterModelDropVersion
    : 'ALTER' 'MODEL' ('IF' 'EXISTS')? QNAME 'DROP' 'VERSION' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-model-modify-version
alterModelModifyVersion
    : 'ALTER' 'MODEL' ('IF' 'EXISTS')? QNAME 'MODIFY' 'VERSION' ID 'SET' ('COMMENT' '=' STRING ('METADATA' '=' STRING)? | 'METADATA' '=' STRING)
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-model-monitor
alterModelMonitor
    : 'ALTER' 'MODEL' 'MONITOR' ('IF' 'EXISTS')? ID ('SUSPEND' | 'RESUME')
    | 'ALTER' 'MODEL' 'MONITOR' ('IF' 'EXISTS')? ID 'SET' ('REFRESH_INTERVAL' '=' STRING ('WAREHOUSE' '=' ID)? ('COMMENT' '=' STRING)? | 'WAREHOUSE' '=' ID ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'MODEL' 'MONITOR' ('IF' 'EXISTS')? ID 'UNSET' 'COMMENT'
    | 'ALTER' 'MODEL' 'MONITOR' ('IF' 'EXISTS')? ID 'SET' 'BASELINE' '=' STRING
    | 'ALTER' 'MODEL' 'MONITOR' ('IF' 'EXISTS')? ID 'UNSET' 'BASELINE'
    | 'ALTER' 'MODEL' 'MONITOR' ('IF' 'EXISTS')? ID 'ADD' 'segment_column' '=' STRING
    | 'ALTER' 'MODEL' 'MONITOR' ('IF' 'EXISTS')? ID 'DROP' 'segment_column' '=' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-multi-party-approval-policy
alterMultiPartyApprovalPolicy
    : 'ALTER' 'MULTI' 'PARTY' 'APPROVAL' 'POLICY' QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'MULTI' 'PARTY' 'APPROVAL' 'POLICY' QNAME 'AS' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-network-policy
alterNetworkPolicy
    : 'ALTER' 'NETWORK' 'POLICY' ('IF' 'EXISTS')? ID 'SET' ('ALLOWED_NETWORK_RULE_LIST' '=' '(' STRING (',' STRING)* ')' ('BLOCKED_NETWORK_RULE_LIST' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_IP_LIST' '=' '(' (STRING (',' STRING)*)? ')')? ('BLOCKED_IP_LIST' '=' '(' (STRING (',' STRING)*)? ')')? ('COMMENT' '=' STRING)? | 'BLOCKED_NETWORK_RULE_LIST' '=' '(' STRING (',' STRING)* ')' ('ALLOWED_IP_LIST' '=' '(' (STRING (',' STRING)*)? ')')? ('BLOCKED_IP_LIST' '=' '(' (STRING (',' STRING)*)? ')')? ('COMMENT' '=' STRING)? | 'ALLOWED_IP_LIST' '=' '(' (STRING (',' STRING)*)? ')' ('BLOCKED_IP_LIST' '=' '(' (STRING (',' STRING)*)? ')')? ('COMMENT' '=' STRING)? | 'BLOCKED_IP_LIST' '=' '(' (STRING (',' STRING)*)? ')' ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'NETWORK' 'POLICY' ('IF' 'EXISTS')? ID 'UNSET' 'COMMENT'
    | 'ALTER' 'NETWORK' 'POLICY' ID 'ADD' ('ALLOWED_NETWORK_RULE_LIST' '=' '(' STRING (',' STRING)* ')' | 'BLOCKED_NETWORK_RULE_LIST' '=' '(' STRING (',' STRING)* ')')
    | 'ALTER' 'NETWORK' 'POLICY' ID 'REMOVE' ('ALLOWED_NETWORK_RULE_LIST' '=' '(' STRING (',' STRING)* ')' | 'BLOCKED_NETWORK_RULE_LIST' '=' '(' STRING (',' STRING)* ')')
    | 'ALTER' 'NETWORK' 'POLICY' ID 'RENAME' 'TO' ID
    | 'ALTER' 'NETWORK' 'POLICY' ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'NETWORK' 'POLICY' ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-network-rule
alterNetworkRule
    : 'ALTER' 'NETWORK' 'RULE' ('IF' 'EXISTS')? QNAME 'SET' 'VALUE_LIST' '=' '(' STRING (',' STRING)* ')' ('COMMENT' '=' STRING)?
    | 'ALTER' 'NETWORK' 'RULE' ('IF' 'EXISTS')? QNAME 'UNSET' ('VALUE_LIST' | 'COMMENT')
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-notebook
// fix: SECRETS: the mappings repeat inside the parentheses, as the page's parameters and CREATE FUNCTION write them
alterNotebook
    : 'ALTER' 'NOTEBOOK' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'NOTEBOOK' ('IF' 'EXISTS')? QNAME 'SET' ('COMMENT' '=' STRING ('QUERY_WAREHOUSE' '=' ID)? ('IDLE_AUTO_SHUTDOWN_TIME_SECONDS' '=' NUMBER)? ('SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')')? | 'QUERY_WAREHOUSE' '=' ID ('IDLE_AUTO_SHUTDOWN_TIME_SECONDS' '=' NUMBER)? ('SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')')? | 'IDLE_AUTO_SHUTDOWN_TIME_SECONDS' '=' NUMBER ('SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')')? | 'SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')')
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-notification-integration
// no syntax

// https://docs.snowflake.com/en/sql-reference/sql/alter-notification-integration-email
// not converted (form): ... stands for text the docs leave out
//   ALTER [ NOTIFICATION ] INTEGRATION [ IF EXISTS ] <name> SET
//     [ ENABLED = { TRUE | FALSE } ]
//     [ ALLOWED_RECIPIENTS = ( '<email_address>' [ , ... '<email_address>' ] ) ]
//     [ DEFAULT_RECIPIENTS = ( '<email_address>' [ , ... '<email_address>' ] ) ]
//     [ DEFAULT_SUBJECT = '<subject_line>' ]
//     [ COMMENT = '<string_literal>' ]
alterNotificationIntegrationEmail
    : 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'UNSET' 'ENABLED'
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'UNSET' 'ALLOWED_RECIPIENTS'
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'UNSET' 'DEFAULT_RECIPIENTS'
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'UNSET' 'DEFAULT_SUBJECT'
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-notification-integration-queue-inbound-azure
alterNotificationIntegrationQueueInboundAzure
    : 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'SET' ('ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-notification-integration-queue-inbound-gcp
alterNotificationIntegrationQueueInboundGcp
    : 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'SET' ('ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-notification-integration-multi-queue
alterNotificationIntegrationMultiQueue
    : 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'SET' 'ACTIVE' '=' STRING ('ENABLED' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)?
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-notification-integration-queue-outbound-aws
alterNotificationIntegrationQueueOutboundAws
    : 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'SET' ('ENABLED' '=' ('TRUE' | 'FALSE'))? 'AWS_SNS_TOPIC_ARN' '=' STRING 'AWS_SNS_ROLE_ARN' '=' STRING ('COMMENT' '=' STRING)?
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-notification-integration-queue-outbound-azure
alterNotificationIntegrationQueueOutboundAzure
    : 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'SET' ('ENABLED' '=' ('TRUE' | 'FALSE') ('AZURE_EVENT_GRID_TOPIC_ENDPOINT' '=' STRING)? ('AZURE_TENANT_ID' '=' STRING)? ('COMMENT' '=' STRING)? | 'AZURE_EVENT_GRID_TOPIC_ENDPOINT' '=' STRING ('AZURE_TENANT_ID' '=' STRING)? ('COMMENT' '=' STRING)? | 'AZURE_TENANT_ID' '=' STRING ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-notification-integration-queue-outbound-gcp
alterNotificationIntegrationQueueOutboundGcp
    : 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'SET' ('ENABLED' '=' ('TRUE' | 'FALSE'))? 'GCP_PUBSUB_SUBSCRIPTION_NAME' '=' STRING ('COMMENT' '=' STRING)?
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-notification-integration-webhooks
alterNotificationIntegrationWebhooks
    : 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'SET' ('ENABLED' '=' ('TRUE' | 'FALSE') ('WEBHOOK_URL' '=' STRING)? ('WEBHOOK_SECRET' '=' QNAME)? ('WEBHOOK_BODY_TEMPLATE' '=' STRING)? ('WEBHOOK_HEADERS' '=' '(' STRING '=' STRING (',' STRING '=' STRING)* ')')? ('COMMENT' '=' STRING)? | 'WEBHOOK_URL' '=' STRING ('WEBHOOK_SECRET' '=' QNAME)? ('WEBHOOK_BODY_TEMPLATE' '=' STRING)? ('WEBHOOK_HEADERS' '=' '(' STRING '=' STRING (',' STRING '=' STRING)* ')')? ('COMMENT' '=' STRING)? | 'WEBHOOK_SECRET' '=' QNAME ('WEBHOOK_BODY_TEMPLATE' '=' STRING)? ('WEBHOOK_HEADERS' '=' '(' STRING '=' STRING (',' STRING '=' STRING)* ')')? ('COMMENT' '=' STRING)? | 'WEBHOOK_BODY_TEMPLATE' '=' STRING ('WEBHOOK_HEADERS' '=' '(' STRING '=' STRING (',' STRING '=' STRING)* ')')? ('COMMENT' '=' STRING)? | 'WEBHOOK_HEADERS' '=' '(' STRING '=' STRING (',' STRING '=' STRING)* ')' ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'NOTIFICATION'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'UNSET' ('ENABLED' | 'WEBHOOK_SECRET' | 'WEBHOOK_BODY_TEMPLATE' | 'WEBHOOK_HEADERS' | 'COMMENT')
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-oflow-data-plane
alterOflowDataPlane
    : 'ALTER' 'OPENFLOW' 'DATA' 'PLANE' 'INTEGRATION' QNAME 'SET' 'EVENT_TABLE' '=' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-online-feature-table
// fix: tagAction is defined under Parameters, outside the Syntax block
alterOnlineFeatureTable
    : 'ALTER' 'ONLINE' 'FEATURE' 'TABLE' ('IF' 'EXISTS')? QNAME ('SUSPEND' | 'RESUME')
    | 'ALTER' 'ONLINE' 'FEATURE' 'TABLE' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'ONLINE' 'FEATURE' 'TABLE' ('IF' 'EXISTS')? QNAME 'REFRESH'
    | 'ALTER' 'ONLINE' 'FEATURE' 'TABLE' ('IF' 'EXISTS')? QNAME 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'ONLINE' 'FEATURE' 'TABLE' ('IF' 'EXISTS')? QNAME 'SET' ('TARGET_LAG' '=' STRING ('WAREHOUSE' '=' ID)? | 'WAREHOUSE' '=' ID)
    | 'ALTER' 'ONLINE' 'FEATURE' 'TABLE' ('IF' 'EXISTS')? QNAME ('SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)* | 'UNSET' 'TAG' QNAME (',' QNAME)*)
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-online-feature-table
createOnlineFeatureTable
    : 'CREATE' ('OR' 'REPLACE')? 'ONLINE' 'FEATURE' 'TABLE' QNAME 'PRIMARY' 'KEY' '(' ID (',' ID)* ')' 'TARGET_LAG' '=' STRING 'WAREHOUSE' '=' ID ('REFRESH_MODE' '=' ('AUTO' | 'FULL' | 'INCREMENTAL'))? ('TIMESTAMP_COLUMN' '=' ID)? ('WITH'? 'COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? 'FROM' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-openflow-connector
// no syntax

// https://docs.snowflake.com/en/sql-reference/sql/alter-openflow-deployment
alterOpenflowDeployment
    : 'ALTER' 'OPENFLOW' 'DEPLOYMENT' ('IF' 'EXISTS')? QNAME 'SET' ('DISPLAY_NAME' '=' STRING ('COMMENT' '=' STRING)? ('EVENT_TABLE' '=' (STRING | 'NONE'))? | 'COMMENT' '=' STRING ('EVENT_TABLE' '=' (STRING | 'NONE'))? | 'EVENT_TABLE' '=' (STRING | 'NONE'))
    | 'ALTER' 'OPENFLOW' 'DEPLOYMENT' ('IF' 'EXISTS')? QNAME 'UNSET' ('DISPLAY_NAME' | 'COMMENT' | 'EVENT_TABLE') (',' ('DISPLAY_NAME' | 'COMMENT' | 'EVENT_TABLE'))*
    | 'ALTER' 'OPENFLOW' 'DEPLOYMENT' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'OPENFLOW' 'DEPLOYMENT' ('IF' 'EXISTS')? QNAME 'UPGRADE'
    | 'ALTER' 'OPENFLOW' 'DEPLOYMENT' ('IF' 'EXISTS')? QNAME 'TERMINATE'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-openflow-runtime
alterOpenflowRuntime
    : 'ALTER' 'OPENFLOW' 'RUNTIME' ('IF' 'EXISTS')? QNAME 'SET' ('DISPLAY_NAME' '=' STRING ('MIN_NODES' '=' NUMBER)? ('MAX_NODES' '=' NUMBER)? ('NODE_TYPE_TIER' '=' ('\'S1\'' | '\'S2\'' | '\'S3\'' | '\'M4\'' | '\'M6\'' | '\'L8\''))? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('EXECUTE_AS_ROLE' '=' ID)? ('COMMENT' '=' STRING)? | 'MIN_NODES' '=' NUMBER ('MAX_NODES' '=' NUMBER)? ('NODE_TYPE_TIER' '=' ('\'S1\'' | '\'S2\'' | '\'S3\'' | '\'M4\'' | '\'M6\'' | '\'L8\''))? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('EXECUTE_AS_ROLE' '=' ID)? ('COMMENT' '=' STRING)? | 'MAX_NODES' '=' NUMBER ('NODE_TYPE_TIER' '=' ('\'S1\'' | '\'S2\'' | '\'S3\'' | '\'M4\'' | '\'M6\'' | '\'L8\''))? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('EXECUTE_AS_ROLE' '=' ID)? ('COMMENT' '=' STRING)? | 'NODE_TYPE_TIER' '=' ('\'S1\'' | '\'S2\'' | '\'S3\'' | '\'M4\'' | '\'M6\'' | '\'L8\'') ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('EXECUTE_AS_ROLE' '=' ID)? ('COMMENT' '=' STRING)? | 'EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')' ('EXECUTE_AS_ROLE' '=' ID)? ('COMMENT' '=' STRING)? | 'EXECUTE_AS_ROLE' '=' ID ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'OPENFLOW' 'RUNTIME' ('IF' 'EXISTS')? QNAME 'UNSET' ('DISPLAY_NAME' | 'EXTERNAL_ACCESS_INTEGRATIONS' | 'COMMENT') (',' ('DISPLAY_NAME' | 'EXTERNAL_ACCESS_INTEGRATIONS' | 'COMMENT'))*
    | 'ALTER' 'OPENFLOW' 'RUNTIME' ('IF' 'EXISTS')? QNAME 'ADD' 'EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')'
    | 'ALTER' 'OPENFLOW' 'RUNTIME' ('IF' 'EXISTS')? QNAME 'REMOVE' 'EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')'
    | 'ALTER' 'OPENFLOW' 'RUNTIME' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'OPENFLOW' 'RUNTIME' ('IF' 'EXISTS')? QNAME 'UPGRADE' ('RECOVERY' | 'FORCE')?
    | 'ALTER' 'OPENFLOW' 'RUNTIME' ('IF' 'EXISTS')? QNAME 'SUSPEND'
    | 'ALTER' 'OPENFLOW' 'RUNTIME' ('IF' 'EXISTS')? QNAME 'RESUME' 'RECOVERY'?
    | 'ALTER' 'OPENFLOW' 'RUNTIME' ('IF' 'EXISTS')? QNAME 'RESTART' 'RECOVERY'?
    | 'ALTER' 'OPENFLOW' 'RUNTIME' ('IF' 'EXISTS')? QNAME 'TERMINATE' 'CASCADE'?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-organization
alterOrganization
    : 'ALTER' 'ORGANIZATION' 'SET' 'DEFAULT_TENANT_TYPE' '=' ('INTERNAL' | 'EXTERNAL')
    | 'ALTER' 'ORGANIZATION' 'SET' 'DOMAIN_NAMES' '=' '(' STRING (',' STRING)* ')'
    | 'ALTER' 'ORGANIZATION' 'UNSET' 'DOMAIN_NAMES'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-organization-account
// not converted (form): undefined fragment "accountParams": the page does not define it, or it did not convert
//   ALTER ORGANIZATION ACCOUNT SET { [ accountParams ] | [ objectParams ] | [ sessionParams ] }
alterOrganizationAccount
    : 'ALTER' 'ORGANIZATION' 'ACCOUNT' 'UNSET' ID (',' ID)*
    | 'ALTER' 'ORGANIZATION' 'ACCOUNT' 'SET' 'RESOURCE_MONITOR' '=' ID
    | 'ALTER' 'ORGANIZATION' 'ACCOUNT' 'SET' ('PASSWORD' | 'SESSION') 'POLICY' QNAME
    | 'ALTER' 'ORGANIZATION' 'ACCOUNT' 'UNSET' ('PASSWORD' | 'SESSION') 'POLICY'
    | 'ALTER' 'ORGANIZATION' 'ACCOUNT' 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'ORGANIZATION' 'ACCOUNT' 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-organization-profile
alterOrganizationProfile
    : 'ALTER' 'ORGANIZATION' 'PROFILE' ('IF' 'EXISTS')? ID 'AS' STRING
    | 'ALTER' 'ORGANIZATION' 'PROFILE' ('IF' 'EXISTS')? ID 'RENAME' 'TO' ID
    | 'ALTER' 'ORGANIZATION' 'PROFILE' ('IF' 'EXISTS')? ID 'PUBLISH'
    | 'ALTER' 'ORGANIZATION' 'PROFILE' ID 'ADD' 'VERSION' (('IF' 'NOT' 'EXISTS')? ID)? 'FROM' STAGE
    | 'ALTER' 'ORGANIZATION' 'PROFILE' ID 'ADD' 'LIVE' 'VERSION' (('IF' 'NOT' 'EXISTS')? ID)? 'FROM' 'LAST'
    | 'ALTER' 'ORGANIZATION' 'PROFILE' ID 'COMMIT'
    | 'ALTER' 'ORGANIZATION' 'PROFILE' ID 'ABORT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-organization-set-event-routing-table
alterOrganizationSetEventRoutingTable
    : 'ALTER' 'ORGANIZATION' 'SET' 'EVENT' 'ROUTING' 'TABLE' QNAME 'FOR' 'ALL' 'APPLICATION' 'LISTINGS'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-organization-unset-event-routing-table
alterOrganizationUnsetEventRoutingTable
    : 'ALTER' 'ORGANIZATION' 'UNSET' 'EVENT' 'ROUTING' 'TABLE' 'FOR' 'ALL' 'APPLICATION' 'LISTINGS'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-organization-user
// fix: UNSET names the objectProperties without their values, as the example UNSET TYPE does
alterOrganizationUser
    : 'ALTER' 'ORGANIZATION' 'USER' ('IF' 'EXISTS')? ID 'SET' alterOrganizationUser_objectProperties
    | 'ALTER' 'ORGANIZATION' 'USER' ID 'UNSET' ('EMAIL' | 'DISPLAY_NAME' | 'FIRST_NAME' | 'MIDDLE_NAME' | 'LAST_NAME' | 'TYPE' | 'COMMENT') (',' ('EMAIL' | 'DISPLAY_NAME' | 'FIRST_NAME' | 'MIDDLE_NAME' | 'LAST_NAME' | 'TYPE' | 'COMMENT'))*
    ;
alterOrganizationUser_objectProperties
    : ('EMAIL' '=' STRING | 'DISPLAY_NAME' '=' STRING | 'FIRST_NAME' '=' STRING | 'MIDDLE_NAME' '=' STRING | 'LAST_NAME' '=' STRING | 'TYPE' '=' ('PERSON' | 'SERVICE') | 'COMMENT' '=' STRING)+
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-organization-user-group
alterOrganizationUserGroup
    : 'ALTER' 'ORGANIZATION' 'USER' 'GROUP' ('IF' 'EXISTS')? ID 'ADD' 'ORGANIZATION' 'USERS' ID (',' ID)*
    | 'ALTER' 'ORGANIZATION' 'USER' 'GROUP' ('IF' 'EXISTS')? ID 'REMOVE' 'ORGANIZATION' 'USERS' ID (',' ID)*
    | 'ALTER' 'ORGANIZATION' 'USER' 'GROUP' ('IF' 'EXISTS')? ID 'SET' 'VISIBILITY' '=' ('ALL' | 'ACCOUNTS' ID (',' ID)* | 'REGION' 'GROUPS' STRING (',' STRING)*)
    | 'ALTER' 'ORGANIZATION' 'USER' 'GROUP' ('IF' 'EXISTS')? ID 'SET' 'IS_GRANTABLE' '=' ('TRUE' | 'FALSE')
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-packages-policy
alterPackagesPolicy
    : 'ALTER' 'PACKAGES' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' ('ALLOWLIST' '=' '(' (STRING (',' STRING)*)? ')' ('BLOCKLIST' '=' '(' (STRING (',' STRING)*)? ')')? ('ADDITIONAL_CREATION_BLOCKLIST' '=' '(' (STRING (',' STRING)*)? ')')? ('COMMENT' '=' STRING)? | 'BLOCKLIST' '=' '(' (STRING (',' STRING)*)? ')' ('ADDITIONAL_CREATION_BLOCKLIST' '=' '(' (STRING (',' STRING)*)? ')')? ('COMMENT' '=' STRING)? | 'ADDITIONAL_CREATION_BLOCKLIST' '=' '(' (STRING (',' STRING)*)? ')' ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'PACKAGES' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' ('ALLOWLIST' | 'BLOCKLIST' | 'ADDITIONAL_CREATION_BLOCKLIST' | 'COMMENT') (',' ('ALLOWLIST' | 'BLOCKLIST' | 'ADDITIONAL_CREATION_BLOCKLIST' | 'COMMENT'))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-password-policy
alterPasswordPolicy
    : 'ALTER' 'PASSWORD' 'POLICY' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'PASSWORD' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' ('PASSWORD_MIN_LENGTH' '=' NUMBER ('PASSWORD_MAX_LENGTH' '=' NUMBER)? ('PASSWORD_MIN_UPPER_CASE_CHARS' '=' NUMBER)? ('PASSWORD_MIN_LOWER_CASE_CHARS' '=' NUMBER)? ('PASSWORD_MIN_NUMERIC_CHARS' '=' NUMBER)? ('PASSWORD_MIN_SPECIAL_CHARS' '=' NUMBER)? ('PASSWORD_MIN_AGE_DAYS' '=' NUMBER)? ('PASSWORD_MAX_AGE_DAYS' '=' NUMBER)? ('PASSWORD_MAX_RETRIES' '=' NUMBER)? ('PASSWORD_LOCKOUT_TIME_MINS' '=' NUMBER)? ('PASSWORD_HISTORY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'PASSWORD_MAX_LENGTH' '=' NUMBER ('PASSWORD_MIN_UPPER_CASE_CHARS' '=' NUMBER)? ('PASSWORD_MIN_LOWER_CASE_CHARS' '=' NUMBER)? ('PASSWORD_MIN_NUMERIC_CHARS' '=' NUMBER)? ('PASSWORD_MIN_SPECIAL_CHARS' '=' NUMBER)? ('PASSWORD_MIN_AGE_DAYS' '=' NUMBER)? ('PASSWORD_MAX_AGE_DAYS' '=' NUMBER)? ('PASSWORD_MAX_RETRIES' '=' NUMBER)? ('PASSWORD_LOCKOUT_TIME_MINS' '=' NUMBER)? ('PASSWORD_HISTORY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'PASSWORD_MIN_UPPER_CASE_CHARS' '=' NUMBER ('PASSWORD_MIN_LOWER_CASE_CHARS' '=' NUMBER)? ('PASSWORD_MIN_NUMERIC_CHARS' '=' NUMBER)? ('PASSWORD_MIN_SPECIAL_CHARS' '=' NUMBER)? ('PASSWORD_MIN_AGE_DAYS' '=' NUMBER)? ('PASSWORD_MAX_AGE_DAYS' '=' NUMBER)? ('PASSWORD_MAX_RETRIES' '=' NUMBER)? ('PASSWORD_LOCKOUT_TIME_MINS' '=' NUMBER)? ('PASSWORD_HISTORY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'PASSWORD_MIN_LOWER_CASE_CHARS' '=' NUMBER ('PASSWORD_MIN_NUMERIC_CHARS' '=' NUMBER)? ('PASSWORD_MIN_SPECIAL_CHARS' '=' NUMBER)? ('PASSWORD_MIN_AGE_DAYS' '=' NUMBER)? ('PASSWORD_MAX_AGE_DAYS' '=' NUMBER)? ('PASSWORD_MAX_RETRIES' '=' NUMBER)? ('PASSWORD_LOCKOUT_TIME_MINS' '=' NUMBER)? ('PASSWORD_HISTORY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'PASSWORD_MIN_NUMERIC_CHARS' '=' NUMBER ('PASSWORD_MIN_SPECIAL_CHARS' '=' NUMBER)? ('PASSWORD_MIN_AGE_DAYS' '=' NUMBER)? ('PASSWORD_MAX_AGE_DAYS' '=' NUMBER)? ('PASSWORD_MAX_RETRIES' '=' NUMBER)? ('PASSWORD_LOCKOUT_TIME_MINS' '=' NUMBER)? ('PASSWORD_HISTORY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'PASSWORD_MIN_SPECIAL_CHARS' '=' NUMBER ('PASSWORD_MIN_AGE_DAYS' '=' NUMBER)? ('PASSWORD_MAX_AGE_DAYS' '=' NUMBER)? ('PASSWORD_MAX_RETRIES' '=' NUMBER)? ('PASSWORD_LOCKOUT_TIME_MINS' '=' NUMBER)? ('PASSWORD_HISTORY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'PASSWORD_MIN_AGE_DAYS' '=' NUMBER ('PASSWORD_MAX_AGE_DAYS' '=' NUMBER)? ('PASSWORD_MAX_RETRIES' '=' NUMBER)? ('PASSWORD_LOCKOUT_TIME_MINS' '=' NUMBER)? ('PASSWORD_HISTORY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'PASSWORD_MAX_AGE_DAYS' '=' NUMBER ('PASSWORD_MAX_RETRIES' '=' NUMBER)? ('PASSWORD_LOCKOUT_TIME_MINS' '=' NUMBER)? ('PASSWORD_HISTORY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'PASSWORD_MAX_RETRIES' '=' NUMBER ('PASSWORD_LOCKOUT_TIME_MINS' '=' NUMBER)? ('PASSWORD_HISTORY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'PASSWORD_LOCKOUT_TIME_MINS' '=' NUMBER ('PASSWORD_HISTORY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'PASSWORD_HISTORY' '=' NUMBER ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'PASSWORD' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'PASSWORD' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'PASSWORD' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' ('PASSWORD_MIN_LENGTH' | 'PASSWORD_MAX_LENGTH' | 'PASSWORD_MIN_UPPER_CASE_CHARS' | 'PASSWORD_MIN_LOWER_CASE_CHARS' | 'PASSWORD_MIN_NUMERIC_CHARS' | 'PASSWORD_MIN_SPECIAL_CHARS' | 'PASSWORD_MIN_AGE_DAYS' | 'PASSWORD_MAX_AGE_DAYS' | 'PASSWORD_MAX_RETRIES' | 'PASSWORD_LOCKOUT_TIME_MINS' | 'PASSWORD_HISTORY' | 'COMMENT') (',' ('PASSWORD_MIN_LENGTH' | 'PASSWORD_MAX_LENGTH' | 'PASSWORD_MIN_UPPER_CASE_CHARS' | 'PASSWORD_MIN_LOWER_CASE_CHARS' | 'PASSWORD_MIN_NUMERIC_CHARS' | 'PASSWORD_MIN_SPECIAL_CHARS' | 'PASSWORD_MIN_AGE_DAYS' | 'PASSWORD_MAX_AGE_DAYS' | 'PASSWORD_MAX_RETRIES' | 'PASSWORD_LOCKOUT_TIME_MINS' | 'PASSWORD_HISTORY' | 'COMMENT'))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-pipe
// fix: the start time is a quoted ISO 8601 timestamp, as the parameters say
// fix: the properties may also be separated by commas, as the SET parameter says; a second form lists them so
// fix: the log event levels are those the parameter lists
alterPipe
    : 'ALTER' 'PIPE' ('IF' 'EXISTS')? QNAME 'SET' (alterPipe_objectProperties ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'PIPE' ('IF' 'EXISTS')? QNAME 'SET' ('ERROR_INTEGRATION' '=' ID | 'PIPE_EXECUTION_PAUSED' '=' ('TRUE' | 'FALSE') | 'LOG_EVENT_LEVEL' '=' ('TRACE' | 'DEBUG' | 'INFO' | 'WARN' | 'ERROR' | 'FATAL' | 'OFF') | 'COMMENT' '=' STRING) (',' ('ERROR_INTEGRATION' '=' ID | 'PIPE_EXECUTION_PAUSED' '=' ('TRUE' | 'FALSE') | 'LOG_EVENT_LEVEL' '=' ('TRACE' | 'DEBUG' | 'INFO' | 'WARN' | 'ERROR' | 'FATAL' | 'OFF') | 'COMMENT' '=' STRING))*
    | 'ALTER' 'PIPE' QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'PIPE' QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'PIPE' ('IF' 'EXISTS')? QNAME 'UNSET' (ID | 'COMMENT') (',' (ID | 'COMMENT'))*
    | 'ALTER' 'PIPE' ('IF' 'EXISTS')? QNAME 'REFRESH' ('PREFIX' '=' STRING)? ('MODIFIED_AFTER' '=' STRING)?
    ;
alterPipe_objectProperties
    : ('ERROR_INTEGRATION' '=' ID | 'PIPE_EXECUTION_PAUSED' '=' ('TRUE' | 'FALSE') | 'LOG_EVENT_LEVEL' '=' ('TRACE' | 'DEBUG' | 'INFO' | 'WARN' | 'ERROR' | 'FATAL' | 'OFF'))+
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-postgres-instance
// fix: the properties are ISSUER, AUDIENCE and ACCEPTED_TOKENS, each token MAPPING_CLAIM, SCOPES_CLAIM and SCOPES, as the parameters and example write them; SCOPES is a list or one string
alterPostgresInstance
    : 'ALTER' 'POSTGRES' 'INSTANCE' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'POSTGRES' 'INSTANCE' ('IF' 'EXISTS')? QNAME 'SET' ('NETWORK_POLICY' '=' STRING ('AUTHENTICATION_AUTHORITY' '=' ('POSTGRES' | 'POSTGRES_OR_SNOWFLAKE' | 'EXTERNAL_OAUTH'))? ('AUTHENTICATION_PROPERTIES' '=' '(' 'ISSUER' '=' STRING ',' 'AUDIENCE' '=' STRING ',' 'ACCEPTED_TOKENS' '=' '(' '(' 'MAPPING_CLAIM' '=' STRING ',' 'SCOPES_CLAIM' '=' STRING ',' 'SCOPES' '=' ('(' STRING (',' STRING)* ')' | STRING) ')' (',' '(' 'MAPPING_CLAIM' '=' STRING ',' 'SCOPES_CLAIM' '=' STRING ',' 'SCOPES' '=' ('(' STRING (',' STRING)* ')' | STRING) ')')* ')' ')')? ('COMMENT' '=' STRING)? ('HIGH_AVAILABILITY' '=' ('TRUE' | 'FALSE'))? ('COMPUTE_FAMILY' '=' STRING)? ('STORAGE_SIZE_GB' '=' NUMBER)? ('STORAGE_INTEGRATION' '=' STRING)? ('POSTGRES_VERSION' '=' ('16' | '17' | '18'))? ('MAINTENANCE_WINDOW_START' '=' NUMBER)? ('POSTGRES_SETTINGS' '=' STRING)? ('APPLY' ('IMMEDIATELY' | 'ON' STRING))? | 'AUTHENTICATION_AUTHORITY' '=' ('POSTGRES' | 'POSTGRES_OR_SNOWFLAKE' | 'EXTERNAL_OAUTH') ('AUTHENTICATION_PROPERTIES' '=' '(' 'ISSUER' '=' STRING ',' 'AUDIENCE' '=' STRING ',' 'ACCEPTED_TOKENS' '=' '(' '(' 'MAPPING_CLAIM' '=' STRING ',' 'SCOPES_CLAIM' '=' STRING ',' 'SCOPES' '=' ('(' STRING (',' STRING)* ')' | STRING) ')' (',' '(' 'MAPPING_CLAIM' '=' STRING ',' 'SCOPES_CLAIM' '=' STRING ',' 'SCOPES' '=' ('(' STRING (',' STRING)* ')' | STRING) ')')* ')' ')')? ('COMMENT' '=' STRING)? ('HIGH_AVAILABILITY' '=' ('TRUE' | 'FALSE'))? ('COMPUTE_FAMILY' '=' STRING)? ('STORAGE_SIZE_GB' '=' NUMBER)? ('STORAGE_INTEGRATION' '=' STRING)? ('POSTGRES_VERSION' '=' ('16' | '17' | '18'))? ('MAINTENANCE_WINDOW_START' '=' NUMBER)? ('POSTGRES_SETTINGS' '=' STRING)? ('APPLY' ('IMMEDIATELY' | 'ON' STRING))? | 'AUTHENTICATION_PROPERTIES' '=' '(' 'ISSUER' '=' STRING ',' 'AUDIENCE' '=' STRING ',' 'ACCEPTED_TOKENS' '=' '(' '(' 'MAPPING_CLAIM' '=' STRING ',' 'SCOPES_CLAIM' '=' STRING ',' 'SCOPES' '=' ('(' STRING (',' STRING)* ')' | STRING) ')' (',' '(' 'MAPPING_CLAIM' '=' STRING ',' 'SCOPES_CLAIM' '=' STRING ',' 'SCOPES' '=' ('(' STRING (',' STRING)* ')' | STRING) ')')* ')' ')' ('COMMENT' '=' STRING)? ('HIGH_AVAILABILITY' '=' ('TRUE' | 'FALSE'))? ('COMPUTE_FAMILY' '=' STRING)? ('STORAGE_SIZE_GB' '=' NUMBER)? ('STORAGE_INTEGRATION' '=' STRING)? ('POSTGRES_VERSION' '=' ('16' | '17' | '18'))? ('MAINTENANCE_WINDOW_START' '=' NUMBER)? ('POSTGRES_SETTINGS' '=' STRING)? ('APPLY' ('IMMEDIATELY' | 'ON' STRING))? | 'COMMENT' '=' STRING ('HIGH_AVAILABILITY' '=' ('TRUE' | 'FALSE'))? ('COMPUTE_FAMILY' '=' STRING)? ('STORAGE_SIZE_GB' '=' NUMBER)? ('STORAGE_INTEGRATION' '=' STRING)? ('POSTGRES_VERSION' '=' ('16' | '17' | '18'))? ('MAINTENANCE_WINDOW_START' '=' NUMBER)? ('POSTGRES_SETTINGS' '=' STRING)? ('APPLY' ('IMMEDIATELY' | 'ON' STRING))? | 'HIGH_AVAILABILITY' '=' ('TRUE' | 'FALSE') ('COMPUTE_FAMILY' '=' STRING)? ('STORAGE_SIZE_GB' '=' NUMBER)? ('STORAGE_INTEGRATION' '=' STRING)? ('POSTGRES_VERSION' '=' ('16' | '17' | '18'))? ('MAINTENANCE_WINDOW_START' '=' NUMBER)? ('POSTGRES_SETTINGS' '=' STRING)? ('APPLY' ('IMMEDIATELY' | 'ON' STRING))? | 'COMPUTE_FAMILY' '=' STRING ('STORAGE_SIZE_GB' '=' NUMBER)? ('STORAGE_INTEGRATION' '=' STRING)? ('POSTGRES_VERSION' '=' ('16' | '17' | '18'))? ('MAINTENANCE_WINDOW_START' '=' NUMBER)? ('POSTGRES_SETTINGS' '=' STRING)? ('APPLY' ('IMMEDIATELY' | 'ON' STRING))? | 'STORAGE_SIZE_GB' '=' NUMBER ('STORAGE_INTEGRATION' '=' STRING)? ('POSTGRES_VERSION' '=' ('16' | '17' | '18'))? ('MAINTENANCE_WINDOW_START' '=' NUMBER)? ('POSTGRES_SETTINGS' '=' STRING)? ('APPLY' ('IMMEDIATELY' | 'ON' STRING))? | 'STORAGE_INTEGRATION' '=' STRING ('POSTGRES_VERSION' '=' ('16' | '17' | '18'))? ('MAINTENANCE_WINDOW_START' '=' NUMBER)? ('POSTGRES_SETTINGS' '=' STRING)? ('APPLY' ('IMMEDIATELY' | 'ON' STRING))? | 'POSTGRES_VERSION' '=' ('16' | '17' | '18') ('MAINTENANCE_WINDOW_START' '=' NUMBER)? ('POSTGRES_SETTINGS' '=' STRING)? ('APPLY' ('IMMEDIATELY' | 'ON' STRING))? | 'MAINTENANCE_WINDOW_START' '=' NUMBER ('POSTGRES_SETTINGS' '=' STRING)? ('APPLY' ('IMMEDIATELY' | 'ON' STRING))? | 'POSTGRES_SETTINGS' '=' STRING ('APPLY' ('IMMEDIATELY' | 'ON' STRING))? | 'APPLY' ('IMMEDIATELY' | 'ON' STRING))
    | 'ALTER' 'POSTGRES' 'INSTANCE' ('IF' 'EXISTS')? QNAME 'UNSET' ('COMMENT' | 'POSTGRES_SETTINGS' | 'NETWORK_POLICY' | 'MAINTENANCE_WINDOW_START' | 'STORAGE_INTEGRATION') (',' ('COMMENT' | 'POSTGRES_SETTINGS' | 'NETWORK_POLICY' | 'MAINTENANCE_WINDOW_START' | 'STORAGE_INTEGRATION'))*
    | 'ALTER' 'POSTGRES' 'INSTANCE' ('IF' 'EXISTS')? QNAME 'SUSPEND'
    | 'ALTER' 'POSTGRES' 'INSTANCE' ('IF' 'EXISTS')? QNAME 'RESUME'
    | 'ALTER' 'POSTGRES' 'INSTANCE' ('IF' 'EXISTS')? QNAME 'RESET' 'ACCESS' 'FOR' ('\'snowflake_admin\'' | '\'application\'')
    | 'ALTER' 'POSTGRES' 'INSTANCE' ('IF' 'EXISTS')? QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'POSTGRES' 'INSTANCE' ('IF' 'EXISTS')? QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-privacy-policy
alterPrivacyPolicy
    : 'ALTER' 'PRIVACY' 'POLICY' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'PRIVACY' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'BODY' '->' EXPR
    | 'ALTER' 'PRIVACY' 'POLICY' QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'PRIVACY' 'POLICY' QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'PRIVACY' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'PRIVACY' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-procedure
// fix: SECRETS' mappings are in parentheses, as the page's parameters write them
// fix: the integrations are a parenthesized list of names, as the parameters and the other pages write it
alterProcedure
    : 'ALTER' 'PROCEDURE' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'RENAME' 'TO' QNAME
    | 'ALTER' 'PROCEDURE' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'SET' ('LOG_LEVEL' '=' STRING ('TRACE_LEVEL' '=' STRING)? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')')? ('COMMENT' '=' STRING)? | 'TRACE_LEVEL' '=' STRING ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')')? ('COMMENT' '=' STRING)? | 'EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')' ('SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')')? ('COMMENT' '=' STRING)? | 'SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')' ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'PROCEDURE' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'UNSET' 'COMMENT'
    | 'ALTER' 'PROCEDURE' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'PROCEDURE' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'PROCEDURE' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'EXECUTE' 'AS' ('OWNER' | 'CALLER' | 'RESTRICTED' 'CALLER')
    | 'ALTER' 'PROCEDURE' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'SET' ('LOG_LEVEL' '=' STRING ('TRACE_LEVEL' '=' STRING)? ('COMMENT' '=' STRING)? | 'TRACE_LEVEL' '=' STRING ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'PROCEDURE' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' 'SET' ('AUTO_EVENT_LOGGING' '=' STRING ('LOG_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('COMMENT' '=' STRING)? | 'LOG_LEVEL' '=' STRING ('TRACE_LEVEL' '=' STRING)? ('COMMENT' '=' STRING)? | 'TRACE_LEVEL' '=' STRING ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-projection-policy
alterProjectionPolicy
    : 'ALTER' 'PROJECTION' 'POLICY' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'PROJECTION' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'BODY' '->' EXPR
    | 'ALTER' 'PROJECTION' 'POLICY' QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'PROJECTION' 'POLICY' QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'PROJECTION' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'PROJECTION' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-replication-group
// fix: OBJECT_TYPES takes the replication object types the parameters list, not any object type
// fix: ALLOWED_INTEGRATION_TYPES takes the integration types the parameters list
alterReplicationGroup
    : 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'RENAME' 'TO' ID
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'SET' ('OBJECT_TYPES' '=' ('ACCOUNT' 'PARAMETERS' | 'DATABASES' | 'EXTERNAL' 'VOLUMES' | 'INTEGRATIONS' | 'NETWORK' 'POLICIES' | 'RESOURCE' 'MONITORS' | 'ROLES' | 'SHARES' | 'USERS' | 'WAREHOUSES') (',' ('ACCOUNT' 'PARAMETERS' | 'DATABASES' | 'EXTERNAL' 'VOLUMES' | 'INTEGRATIONS' | 'NETWORK' 'POLICIES' | 'RESOURCE' 'MONITORS' | 'ROLES' | 'SHARES' | 'USERS' | 'WAREHOUSES'))* ('ALLOWED_DATABASES' '=' ID (',' ID)*)? ('ALLOWED_EXTERNAL_VOLUMES' '=' ID (',' ID)*)? ('ALLOWED_SHARES' '=' ID (',' ID)*)? | 'ALLOWED_DATABASES' '=' ID (',' ID)* ('ALLOWED_EXTERNAL_VOLUMES' '=' ID (',' ID)*)? ('ALLOWED_SHARES' '=' ID (',' ID)*)? | 'ALLOWED_EXTERNAL_VOLUMES' '=' ID (',' ID)* ('ALLOWED_SHARES' '=' ID (',' ID)*)? | 'ALLOWED_SHARES' '=' ID (',' ID)*)
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'SET' 'OBJECT_TYPES' '=' 'INTEGRATIONS' (',' ('ACCOUNT' 'PARAMETERS' | 'DATABASES' | 'EXTERNAL' 'VOLUMES' | 'INTEGRATIONS' | 'NETWORK' 'POLICIES' | 'RESOURCE' 'MONITORS' | 'ROLES' | 'SHARES' | 'USERS' | 'WAREHOUSES'))* 'ALLOWED_INTEGRATION_TYPES' '=' ('SECURITY' 'INTEGRATIONS' | 'API' 'INTEGRATIONS' | 'STORAGE' 'INTEGRATIONS' | 'EXTERNAL' 'ACCESS' 'INTEGRATIONS' | 'NOTIFICATION' 'INTEGRATIONS') (',' ('SECURITY' 'INTEGRATIONS' | 'API' 'INTEGRATIONS' | 'STORAGE' 'INTEGRATIONS' | 'EXTERNAL' 'ACCESS' 'INTEGRATIONS' | 'NOTIFICATION' 'INTEGRATIONS'))*
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'SET' 'REPLICATION_SCHEDULE' '=' STRING
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'SET' 'ERROR_INTEGRATION' '=' ID
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'UNSET' ('COMMENT' | 'REPLICATION_SCHEDULE' | 'ERROR_INTEGRATION') (',' ('COMMENT' | 'REPLICATION_SCHEDULE' | 'ERROR_INTEGRATION'))*
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'ADD' ID (',' ID)* 'TO' 'ALLOWED_DATABASES'
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'MOVE' 'DATABASES' ID (',' ID)* 'TO' 'REPLICATION' 'GROUP' ID
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'REMOVE' ID (',' ID)* 'FROM' 'ALLOWED_DATABASES'
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'ADD' ID (',' ID)* 'TO' 'ALLOWED_EXTERNAL_VOLUMES'
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'MOVE' 'EXTERNAL' 'VOLUMES' ID (',' ID)* 'TO' 'REPLICATION' 'GROUP' ID
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'REMOVE' ID (',' ID)* 'FROM' 'ALLOWED_EXTERNAL_VOLUMES'
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'ADD' ID (',' ID)* 'TO' 'ALLOWED_SHARES'
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'MOVE' 'SHARES' ID (',' ID)* 'TO' 'REPLICATION' 'GROUP' ID
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'REMOVE' ID (',' ID)* 'FROM' 'ALLOWED_SHARES'
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'ADD' ID '.' ID (',' ID '.' ID)* 'TO' 'ALLOWED_ACCOUNTS' ('IGNORE' 'EDITION' 'CHECK')?
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'REMOVE' ID '.' ID (',' ID '.' ID)* 'FROM' 'ALLOWED_ACCOUNTS'
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'REFRESH'
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'SUSPEND' 'IMMEDIATE'?
    | 'ALTER' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID 'RESUME'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-resource-monitor
alterResourceMonitor
    : 'ALTER' 'RESOURCE' 'MONITOR' ('IF' 'EXISTS')? ID ('SET' ('CREDIT_QUOTA' '=' NUMBER ('FREQUENCY' '=' ('MONTHLY' | 'DAILY' | 'WEEKLY' | 'YEARLY' | 'NEVER'))? ('START_TIMESTAMP' '=' (STRING | 'IMMEDIATELY'))? ('END_TIMESTAMP' '=' STRING)? ('NOTIFY_USERS' '=' '(' ID (',' ID)* ')')? | 'FREQUENCY' '=' ('MONTHLY' | 'DAILY' | 'WEEKLY' | 'YEARLY' | 'NEVER') ('START_TIMESTAMP' '=' (STRING | 'IMMEDIATELY'))? ('END_TIMESTAMP' '=' STRING)? ('NOTIFY_USERS' '=' '(' ID (',' ID)* ')')? | 'START_TIMESTAMP' '=' (STRING | 'IMMEDIATELY') ('END_TIMESTAMP' '=' STRING)? ('NOTIFY_USERS' '=' '(' ID (',' ID)* ')')? | 'END_TIMESTAMP' '=' STRING ('NOTIFY_USERS' '=' '(' ID (',' ID)* ')')? | 'NOTIFY_USERS' '=' '(' ID (',' ID)* ')'))? ('TRIGGERS' alterResourceMonitor_triggerDefinition alterResourceMonitor_triggerDefinition*)?
    ;
alterResourceMonitor_triggerDefinition
    : 'ON' NUMBER 'PERCENT' 'DO' ('SUSPEND' | 'SUSPEND_IMMEDIATE' | 'NOTIFY')
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-restricted-session-scope
alterRestrictedSessionScope
    : 'ALTER' ('RESTRICTED' 'SESSION' 'SCOPE' | 'RSS') ('IF' 'EXISTS')? QNAME 'AS' STRING
    | 'ALTER' ('RESTRICTED' 'SESSION' 'SCOPE' | 'RSS') ('IF' 'EXISTS')? QNAME 'ADD' 'AS' STRING
    | 'ALTER' ('RESTRICTED' 'SESSION' 'SCOPE' | 'RSS') ('IF' 'EXISTS')? QNAME 'REMOVE' 'AS' STRING
    | 'ALTER' ('RESTRICTED' 'SESSION' 'SCOPE' | 'RSS') ('IF' 'EXISTS')? QNAME 'SET' 'COMMENT' '=' STRING
    | 'ALTER' ('RESTRICTED' 'SESSION' 'SCOPE' | 'RSS') ('IF' 'EXISTS')? QNAME 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-role
alterRole
    : 'ALTER' 'ROLE' ('IF' 'EXISTS')? ID 'RENAME' 'TO' ID
    | 'ALTER' 'ROLE' ('IF' 'EXISTS')? ID 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'ROLE' ('IF' 'EXISTS')? ID 'UNSET' 'COMMENT'
    | 'ALTER' 'ROLE' ('IF' 'EXISTS')? ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'ROLE' ('IF' 'EXISTS')? ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-row-access-policy
alterRowAccessPolicy
    : 'ALTER' 'ROW' 'ACCESS' 'POLICY' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'ROW' 'ACCESS' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'BODY' '->' EXPR
    | 'ALTER' 'ROW' 'ACCESS' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'ROW' 'ACCESS' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'ROW' 'ACCESS' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'ROW' 'ACCESS' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-schema
// fix: a } inside [ OBJECT_VISIBILITY = PRIVILEGED ] closes no bracket; dropped
alterSchema
    : 'ALTER' 'SCHEMA' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'SCHEMA' ('IF' 'EXISTS')? QNAME 'SWAP' 'WITH' QNAME
    | 'ALTER' 'SCHEMA' ('IF' 'EXISTS')? QNAME 'SET' ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('STAGE_UNDROP_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('EXTERNAL_VOLUME' '=' ID)? ('CATALOG' '=' ID)? ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_VERSION_DEFAULT' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('LOG_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CLASSIFICATION_PROFILE' '=' STRING)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER ('STAGE_UNDROP_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('EXTERNAL_VOLUME' '=' ID)? ('CATALOG' '=' ID)? ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_VERSION_DEFAULT' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('LOG_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CLASSIFICATION_PROFILE' '=' STRING)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'STAGE_UNDROP_RETENTION_TIME_IN_DAYS' '=' NUMBER ('EXTERNAL_VOLUME' '=' ID)? ('CATALOG' '=' ID)? ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_VERSION_DEFAULT' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('LOG_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CLASSIFICATION_PROFILE' '=' STRING)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'EXTERNAL_VOLUME' '=' ID ('CATALOG' '=' ID)? ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_VERSION_DEFAULT' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('LOG_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CLASSIFICATION_PROFILE' '=' STRING)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'CATALOG' '=' ID ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_VERSION_DEFAULT' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('LOG_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CLASSIFICATION_PROFILE' '=' STRING)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING ('ICEBERG_VERSION_DEFAULT' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('LOG_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CLASSIFICATION_PROFILE' '=' STRING)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'ICEBERG_VERSION_DEFAULT' '=' NUMBER ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('LOG_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CLASSIFICATION_PROFILE' '=' STRING)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\'') ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('LOG_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CLASSIFICATION_PROFILE' '=' STRING)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE') ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('LOG_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CLASSIFICATION_PROFILE' '=' STRING)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') ('DEFAULT_DDL_COLLATION' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('LOG_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CLASSIFICATION_PROFILE' '=' STRING)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'DEFAULT_DDL_COLLATION' '=' STRING ('DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING)? ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('LOG_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CLASSIFICATION_PROFILE' '=' STRING)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU' '=' STRING ('DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING)? ('LOG_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CLASSIFICATION_PROFILE' '=' STRING)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU' '=' STRING ('LOG_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CLASSIFICATION_PROFILE' '=' STRING)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'LOG_LEVEL' '=' STRING ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CLASSIFICATION_PROFILE' '=' STRING)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'TRACE_LEVEL' '=' STRING ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CLASSIFICATION_PROFILE' '=' STRING)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED') ('CLASSIFICATION_PROFILE' '=' STRING)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'CLASSIFICATION_PROFILE' '=' STRING ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'COMMENT' '=' STRING ('CATALOG_SYNC' '=' STRING)? ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'CATALOG_SYNC' '=' STRING ('REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\''))? ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'REPLICABLE_WITH_FAILOVER_GROUPS' '=' ('\'YES\'' | '\'NO\'') ('BASE_LOCATION_PREFIX' '=' STRING)? ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'BASE_LOCATION_PREFIX' '=' STRING ('DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' '=' STRING ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'CONTACT' ID '=' ID (',' ID '=' ID)* ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'OBJECT_VISIBILITY' '=' 'PRIVILEGED' ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE') ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'OAUTH_AUTHORIZATION_SERVER' '=' ID ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? | 'OAUTH_SCOPES_SUPPORTED' '=' STRING)
    | 'ALTER' 'SCHEMA' ('IF' 'EXISTS')? QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'SCHEMA' ('IF' 'EXISTS')? QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'SCHEMA' ('IF' 'EXISTS')? QNAME 'UNSET' ('DATA_RETENTION_TIME_IN_DAYS' | 'MAX_DATA_EXTENSION_TIME_IN_DAYS' | 'STAGE_UNDROP_RETENTION_TIME_IN_DAYS' | 'EXTERNAL_VOLUME' | 'CATALOG' | 'ICEBERG_DEFAULT_DDL_COLLATION' | 'ICEBERG_VERSION_DEFAULT' | 'ICEBERG_MERGE_ON_READ_BEHAVIOR' | 'ENABLE_ICEBERG_MERGE_ON_READ' | 'REPLACE_INVALID_CHARACTERS' | 'DEFAULT_DDL_COLLATION' | 'LOG_LEVEL' | 'TRACE_LEVEL' | 'STORAGE_SERIALIZATION_POLICY' | 'COMMENT' | 'CATALOG_SYNC' | 'REPLICABLE_WITH_FAILOVER_GROUPS' | 'BASE_LOCATION_PREFIX' | 'DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' | 'CONTACT' ID 'CLASSIFICATION_PROFILE' 'OBJECT_VISIBILITY' | 'CONTACT' ID | 'CLASSIFICATION_PROFILE' | 'ENABLE_DATA_COMPACTION' | 'OAUTH_AUTHORIZATION_SERVER' | 'OAUTH_SCOPES_SUPPORTED') (',' ('DATA_RETENTION_TIME_IN_DAYS' | 'MAX_DATA_EXTENSION_TIME_IN_DAYS' | 'STAGE_UNDROP_RETENTION_TIME_IN_DAYS' | 'EXTERNAL_VOLUME' | 'CATALOG' | 'ICEBERG_DEFAULT_DDL_COLLATION' | 'ICEBERG_VERSION_DEFAULT' | 'ICEBERG_MERGE_ON_READ_BEHAVIOR' | 'ENABLE_ICEBERG_MERGE_ON_READ' | 'REPLACE_INVALID_CHARACTERS' | 'DEFAULT_DDL_COLLATION' | 'LOG_LEVEL' | 'TRACE_LEVEL' | 'STORAGE_SERIALIZATION_POLICY' | 'COMMENT' | 'CATALOG_SYNC' | 'REPLICABLE_WITH_FAILOVER_GROUPS' | 'BASE_LOCATION_PREFIX' | 'DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE' | 'CONTACT' ID 'CLASSIFICATION_PROFILE' 'OBJECT_VISIBILITY' | 'CONTACT' ID | 'CLASSIFICATION_PROFILE' | 'ENABLE_DATA_COMPACTION' | 'OAUTH_AUTHORIZATION_SERVER' | 'OAUTH_SCOPES_SUPPORTED'))*
    | 'ALTER' 'SCHEMA' ('IF' 'EXISTS')? QNAME ('ENABLE' | 'DISABLE') 'MANAGED' 'ACCESS'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-secret
alterSecret
    : 'ALTER' 'SECRET' ('IF' 'EXISTS')? QNAME 'SET' ('OAUTH_SCOPES' '=' '(' STRING (',' STRING)* ')' ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'SECRET' ('IF' 'EXISTS')? QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'SECRET' ('IF' 'EXISTS')? QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'SECRET' ('IF' 'EXISTS')? QNAME 'SET' ('OAUTH_REFRESH_TOKEN' '=' STRING ('OAUTH_REFRESH_TOKEN_EXPIRY_TIME' '=' STRING)? ('COMMENT' '=' STRING)? | 'OAUTH_REFRESH_TOKEN_EXPIRY_TIME' '=' STRING ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'SECRET' ('IF' 'EXISTS')? QNAME 'SET' ('API_AUTHENTICATION' '=' STRING ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'SECRET' ('IF' 'EXISTS')? QNAME 'SET' ('USERNAME' '=' STRING ('PASSWORD' '=' STRING)? ('COMMENT' '=' STRING)? | 'PASSWORD' '=' STRING ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'SECRET' ('IF' 'EXISTS')? QNAME 'SET' ('SECRET_STRING' '=' STRING ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration
// fix: the parameters are each integration type's, given on its own page; {parameters} leaves the form unconverted
// not converted (form): undefined fragment "parameters": the page does not define it
//   ALTER [ SECURITY ] INTEGRATION [ IF EXISTS ] <name> SET {parameters}
alterSecurityIntegration
    : 'ALTER' 'SECURITY'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'UNSET' ID
    | 'ALTER' 'SECURITY'? 'INTEGRATION' ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'SECURITY'? 'INTEGRATION' ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-api-auth
// not converted (form): ... stands for text the docs leave out
//   ALTER [ SECURITY ] INTEGRATION [ IF EXISTS ] <name> UNSET {
//     ENABLED | [ , ... ]
//   }
// not converted (form): ... stands for text the docs leave out
//   ALTER [ SECURITY ] INTEGRATION [ IF EXISTS ] <name> UNSET {
//     ENABLED | [ , ... ]
//   }
// not converted (form): ... stands for text the docs leave out
//   ALTER [ SECURITY ] INTEGRATION [ IF EXISTS ] <name> UNSET {
//     ENABLED | [ , ... ]
//   }
alterSecurityIntegrationApiAuth
    : 'ALTER' 'SECURITY' 'INTEGRATION' ID 'SET' ('ENABLED' '=' ('TRUE' | 'FALSE') ('OAUTH_TOKEN_ENDPOINT' '=' STRING)? ('OAUTH_CLIENT_AUTH_METHOD' '=' ('CLIENT_SECRET_BASIC' | 'CLIENT_SECRET_POST'))? ('OAUTH_CLIENT_ID' '=' STRING)? ('OAUTH_CLIENT_SECRET' '=' STRING)? ('OAUTH_GRANT' '=' '\'CLIENT_CREDENTIALS\'')? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_ALLOWED_SCOPES' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'OAUTH_TOKEN_ENDPOINT' '=' STRING ('OAUTH_CLIENT_AUTH_METHOD' '=' ('CLIENT_SECRET_BASIC' | 'CLIENT_SECRET_POST'))? ('OAUTH_CLIENT_ID' '=' STRING)? ('OAUTH_CLIENT_SECRET' '=' STRING)? ('OAUTH_GRANT' '=' '\'CLIENT_CREDENTIALS\'')? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_ALLOWED_SCOPES' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'OAUTH_CLIENT_AUTH_METHOD' '=' ('CLIENT_SECRET_BASIC' | 'CLIENT_SECRET_POST') ('OAUTH_CLIENT_ID' '=' STRING)? ('OAUTH_CLIENT_SECRET' '=' STRING)? ('OAUTH_GRANT' '=' '\'CLIENT_CREDENTIALS\'')? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_ALLOWED_SCOPES' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'OAUTH_CLIENT_ID' '=' STRING ('OAUTH_CLIENT_SECRET' '=' STRING)? ('OAUTH_GRANT' '=' '\'CLIENT_CREDENTIALS\'')? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_ALLOWED_SCOPES' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'OAUTH_CLIENT_SECRET' '=' STRING ('OAUTH_GRANT' '=' '\'CLIENT_CREDENTIALS\'')? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_ALLOWED_SCOPES' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'OAUTH_GRANT' '=' '\'CLIENT_CREDENTIALS\'' ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_ALLOWED_SCOPES' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER ('OAUTH_ALLOWED_SCOPES' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'OAUTH_ALLOWED_SCOPES' '=' '(' STRING (',' STRING)* ')' ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'SECURITY'? 'INTEGRATION' ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'SECURITY'? 'INTEGRATION' ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'SECURITY' 'INTEGRATION' ID 'SET' ('ENABLED' '=' ('TRUE' | 'FALSE') ('OAUTH_AUTHORIZATION_ENDPOINT' '=' STRING)? ('OAUTH_TOKEN_ENDPOINT' '=' STRING)? ('OAUTH_CLIENT_AUTH_METHOD' '=' ('CLIENT_SECRET_BASIC' | 'CLIENT_SECRET_POST'))? ('OAUTH_CLIENT_ID' '=' STRING)? ('OAUTH_CLIENT_SECRET' '=' STRING)? ('OAUTH_GRANT' '=' '\'AUTHORIZATION_CODE\'')? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'OAUTH_AUTHORIZATION_ENDPOINT' '=' STRING ('OAUTH_TOKEN_ENDPOINT' '=' STRING)? ('OAUTH_CLIENT_AUTH_METHOD' '=' ('CLIENT_SECRET_BASIC' | 'CLIENT_SECRET_POST'))? ('OAUTH_CLIENT_ID' '=' STRING)? ('OAUTH_CLIENT_SECRET' '=' STRING)? ('OAUTH_GRANT' '=' '\'AUTHORIZATION_CODE\'')? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'OAUTH_TOKEN_ENDPOINT' '=' STRING ('OAUTH_CLIENT_AUTH_METHOD' '=' ('CLIENT_SECRET_BASIC' | 'CLIENT_SECRET_POST'))? ('OAUTH_CLIENT_ID' '=' STRING)? ('OAUTH_CLIENT_SECRET' '=' STRING)? ('OAUTH_GRANT' '=' '\'AUTHORIZATION_CODE\'')? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'OAUTH_CLIENT_AUTH_METHOD' '=' ('CLIENT_SECRET_BASIC' | 'CLIENT_SECRET_POST') ('OAUTH_CLIENT_ID' '=' STRING)? ('OAUTH_CLIENT_SECRET' '=' STRING)? ('OAUTH_GRANT' '=' '\'AUTHORIZATION_CODE\'')? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'OAUTH_CLIENT_ID' '=' STRING ('OAUTH_CLIENT_SECRET' '=' STRING)? ('OAUTH_GRANT' '=' '\'AUTHORIZATION_CODE\'')? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'OAUTH_CLIENT_SECRET' '=' STRING ('OAUTH_GRANT' '=' '\'AUTHORIZATION_CODE\'')? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'OAUTH_GRANT' '=' '\'AUTHORIZATION_CODE\'' ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'SECURITY' 'INTEGRATION' ID 'SET' ('ENABLED' '=' ('TRUE' | 'FALSE') ('OAUTH_AUTHORIZATION_ENDPOINT' '=' STRING)? ('OAUTH_TOKEN_ENDPOINT' '=' STRING)? ('OAUTH_CLIENT_AUTH_METHOD' '=' ('CLIENT_SECRET_BASIC' | 'CLIENT_SECRET_POST'))? ('OAUTH_CLIENT_ID' '=' STRING)? ('OAUTH_CLIENT_SECRET' '=' STRING)? ('OAUTH_ASSERTION_ISSUER' '=' STRING)? ('OAUTH_GRANT' '=' '\'JWT_BEARER\'')? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'OAUTH_AUTHORIZATION_ENDPOINT' '=' STRING ('OAUTH_TOKEN_ENDPOINT' '=' STRING)? ('OAUTH_CLIENT_AUTH_METHOD' '=' ('CLIENT_SECRET_BASIC' | 'CLIENT_SECRET_POST'))? ('OAUTH_CLIENT_ID' '=' STRING)? ('OAUTH_CLIENT_SECRET' '=' STRING)? ('OAUTH_ASSERTION_ISSUER' '=' STRING)? ('OAUTH_GRANT' '=' '\'JWT_BEARER\'')? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'OAUTH_TOKEN_ENDPOINT' '=' STRING ('OAUTH_CLIENT_AUTH_METHOD' '=' ('CLIENT_SECRET_BASIC' | 'CLIENT_SECRET_POST'))? ('OAUTH_CLIENT_ID' '=' STRING)? ('OAUTH_CLIENT_SECRET' '=' STRING)? ('OAUTH_ASSERTION_ISSUER' '=' STRING)? ('OAUTH_GRANT' '=' '\'JWT_BEARER\'')? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'OAUTH_CLIENT_AUTH_METHOD' '=' ('CLIENT_SECRET_BASIC' | 'CLIENT_SECRET_POST') ('OAUTH_CLIENT_ID' '=' STRING)? ('OAUTH_CLIENT_SECRET' '=' STRING)? ('OAUTH_ASSERTION_ISSUER' '=' STRING)? ('OAUTH_GRANT' '=' '\'JWT_BEARER\'')? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'OAUTH_CLIENT_ID' '=' STRING ('OAUTH_CLIENT_SECRET' '=' STRING)? ('OAUTH_ASSERTION_ISSUER' '=' STRING)? ('OAUTH_GRANT' '=' '\'JWT_BEARER\'')? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'OAUTH_CLIENT_SECRET' '=' STRING ('OAUTH_ASSERTION_ISSUER' '=' STRING)? ('OAUTH_GRANT' '=' '\'JWT_BEARER\'')? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'OAUTH_ASSERTION_ISSUER' '=' STRING ('OAUTH_GRANT' '=' '\'JWT_BEARER\'')? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'OAUTH_GRANT' '=' '\'JWT_BEARER\'' ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-api-auth-external-secret-provider
alterSecurityIntegrationApiAuthExternalSecretProvider
    : 'ALTER' 'SECURITY' 'INTEGRATION' ID 'SET' ('ENABLED' '=' ('TRUE' | 'FALSE') ('AWS_ROLE_ARN' '=' STRING)? ('AWS_REGION' '=' STRING)? ('COMMENT' '=' STRING)? | 'AWS_ROLE_ARN' '=' STRING ('AWS_REGION' '=' STRING)? ('COMMENT' '=' STRING)? | 'AWS_REGION' '=' STRING ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'SECURITY' 'INTEGRATION' ID 'SET' ('ENABLED' '=' ('TRUE' | 'FALSE') ('AZURE_TENANT_ID' '=' STRING)? ('AZURE_AD_APPLICATION_ID' '=' STRING)? ('AZURE_KEY_VAULT_URI' '=' STRING)? ('COMMENT' '=' STRING)? | 'AZURE_TENANT_ID' '=' STRING ('AZURE_AD_APPLICATION_ID' '=' STRING)? ('AZURE_KEY_VAULT_URI' '=' STRING)? ('COMMENT' '=' STRING)? | 'AZURE_AD_APPLICATION_ID' '=' STRING ('AZURE_KEY_VAULT_URI' '=' STRING)? ('COMMENT' '=' STRING)? | 'AZURE_KEY_VAULT_URI' '=' STRING ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'SECURITY' 'INTEGRATION' ID 'SET' ('ENABLED' '=' ('TRUE' | 'FALSE') ('GCP_SECRET_MANAGER_PROJECT' '=' STRING)? ('GCP_WIF_PROVIDER' '=' STRING)? ('GCP_SERVICE_ACCOUNT_EMAIL' '=' STRING)? ('COMMENT' '=' STRING)? | 'GCP_SECRET_MANAGER_PROJECT' '=' STRING ('GCP_WIF_PROVIDER' '=' STRING)? ('GCP_SERVICE_ACCOUNT_EMAIL' '=' STRING)? ('COMMENT' '=' STRING)? | 'GCP_WIF_PROVIDER' '=' STRING ('GCP_SERVICE_ACCOUNT_EMAIL' '=' STRING)? ('COMMENT' '=' STRING)? | 'GCP_SERVICE_ACCOUNT_EMAIL' '=' STRING ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-aws-iam
alterSecurityIntegrationAwsIam
    : 'ALTER' 'SECURITY'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'SET' ('TYPE' '=' 'AWS_IAM' ('AWS_ROLE_ARN' '=' STRING)? ('ENABLED' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'AWS_ROLE_ARN' '=' STRING ('ENABLED' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'SECURITY'? 'INTEGRATION' ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'SECURITY'? 'INTEGRATION' ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-oauth-external
alterSecurityIntegrationOauthExternal
    : 'ALTER' 'SECURITY'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'SET' ('TYPE' '=' 'EXTERNAL_OAUTH' ('ENABLED' '=' ('TRUE' | 'FALSE'))? ('EXTERNAL_OAUTH_TYPE' '=' ('OKTA' | 'AZURE' | 'PING_FEDERATE' | 'CUSTOM'))? ('EXTERNAL_OAUTH_ISSUER' '=' STRING)? ('EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM' '=' STRING | 'EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM' '=' '(' STRING ',' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE' '=' ('\'LOGIN_NAME\'' | '\'EMAIL_ADDRESS\''))? ('EXTERNAL_OAUTH_JWS_KEYS_URL' '=' STRING)? ('EXTERNAL_OAUTH_JWS_KEYS_URL' '=' STRING | 'EXTERNAL_OAUTH_JWS_KEYS_URL' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_RSA_PUBLIC_KEY' '=' STRING)? ('EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2' '=' STRING)? ('EXTERNAL_OAUTH_BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_AUDIENCE_LIST' '=' '(' STRING ')')? ('EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'DISABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE_FOR_PRIVILEGE')? ('EXTERNAL_OAUTH_SCOPE_DELIMITER' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('NETWORK_POLICY' '=' STRING)? ('COMMENT' '=' STRING)? | 'ENABLED' '=' ('TRUE' | 'FALSE') ('EXTERNAL_OAUTH_TYPE' '=' ('OKTA' | 'AZURE' | 'PING_FEDERATE' | 'CUSTOM'))? ('EXTERNAL_OAUTH_ISSUER' '=' STRING)? ('EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM' '=' STRING | 'EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM' '=' '(' STRING ',' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE' '=' ('\'LOGIN_NAME\'' | '\'EMAIL_ADDRESS\''))? ('EXTERNAL_OAUTH_JWS_KEYS_URL' '=' STRING)? ('EXTERNAL_OAUTH_JWS_KEYS_URL' '=' STRING | 'EXTERNAL_OAUTH_JWS_KEYS_URL' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_RSA_PUBLIC_KEY' '=' STRING)? ('EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2' '=' STRING)? ('EXTERNAL_OAUTH_BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_AUDIENCE_LIST' '=' '(' STRING ')')? ('EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'DISABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE_FOR_PRIVILEGE')? ('EXTERNAL_OAUTH_SCOPE_DELIMITER' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('NETWORK_POLICY' '=' STRING)? ('COMMENT' '=' STRING)? | 'EXTERNAL_OAUTH_TYPE' '=' ('OKTA' | 'AZURE' | 'PING_FEDERATE' | 'CUSTOM') ('EXTERNAL_OAUTH_ISSUER' '=' STRING)? ('EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM' '=' STRING | 'EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM' '=' '(' STRING ',' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE' '=' ('\'LOGIN_NAME\'' | '\'EMAIL_ADDRESS\''))? ('EXTERNAL_OAUTH_JWS_KEYS_URL' '=' STRING)? ('EXTERNAL_OAUTH_JWS_KEYS_URL' '=' STRING | 'EXTERNAL_OAUTH_JWS_KEYS_URL' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_RSA_PUBLIC_KEY' '=' STRING)? ('EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2' '=' STRING)? ('EXTERNAL_OAUTH_BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_AUDIENCE_LIST' '=' '(' STRING ')')? ('EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'DISABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE_FOR_PRIVILEGE')? ('EXTERNAL_OAUTH_SCOPE_DELIMITER' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('NETWORK_POLICY' '=' STRING)? ('COMMENT' '=' STRING)? | 'EXTERNAL_OAUTH_ISSUER' '=' STRING ('EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM' '=' STRING | 'EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM' '=' '(' STRING ',' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE' '=' ('\'LOGIN_NAME\'' | '\'EMAIL_ADDRESS\''))? ('EXTERNAL_OAUTH_JWS_KEYS_URL' '=' STRING)? ('EXTERNAL_OAUTH_JWS_KEYS_URL' '=' STRING | 'EXTERNAL_OAUTH_JWS_KEYS_URL' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_RSA_PUBLIC_KEY' '=' STRING)? ('EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2' '=' STRING)? ('EXTERNAL_OAUTH_BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_AUDIENCE_LIST' '=' '(' STRING ')')? ('EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'DISABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE_FOR_PRIVILEGE')? ('EXTERNAL_OAUTH_SCOPE_DELIMITER' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('NETWORK_POLICY' '=' STRING)? ('COMMENT' '=' STRING)? | ('EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM' '=' STRING | 'EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM' '=' '(' STRING ',' STRING (',' STRING)* ')') ('EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE' '=' ('\'LOGIN_NAME\'' | '\'EMAIL_ADDRESS\''))? ('EXTERNAL_OAUTH_JWS_KEYS_URL' '=' STRING)? ('EXTERNAL_OAUTH_JWS_KEYS_URL' '=' STRING | 'EXTERNAL_OAUTH_JWS_KEYS_URL' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_RSA_PUBLIC_KEY' '=' STRING)? ('EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2' '=' STRING)? ('EXTERNAL_OAUTH_BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_AUDIENCE_LIST' '=' '(' STRING ')')? ('EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'DISABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE_FOR_PRIVILEGE')? ('EXTERNAL_OAUTH_SCOPE_DELIMITER' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('NETWORK_POLICY' '=' STRING)? ('COMMENT' '=' STRING)? | 'EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE' '=' ('\'LOGIN_NAME\'' | '\'EMAIL_ADDRESS\'') ('EXTERNAL_OAUTH_JWS_KEYS_URL' '=' STRING)? ('EXTERNAL_OAUTH_JWS_KEYS_URL' '=' STRING | 'EXTERNAL_OAUTH_JWS_KEYS_URL' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_RSA_PUBLIC_KEY' '=' STRING)? ('EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2' '=' STRING)? ('EXTERNAL_OAUTH_BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_AUDIENCE_LIST' '=' '(' STRING ')')? ('EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'DISABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE_FOR_PRIVILEGE')? ('EXTERNAL_OAUTH_SCOPE_DELIMITER' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('NETWORK_POLICY' '=' STRING)? ('COMMENT' '=' STRING)? | 'EXTERNAL_OAUTH_JWS_KEYS_URL' '=' STRING ('EXTERNAL_OAUTH_JWS_KEYS_URL' '=' STRING | 'EXTERNAL_OAUTH_JWS_KEYS_URL' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_RSA_PUBLIC_KEY' '=' STRING)? ('EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2' '=' STRING)? ('EXTERNAL_OAUTH_BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_AUDIENCE_LIST' '=' '(' STRING ')')? ('EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'DISABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE_FOR_PRIVILEGE')? ('EXTERNAL_OAUTH_SCOPE_DELIMITER' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('NETWORK_POLICY' '=' STRING)? ('COMMENT' '=' STRING)? | ('EXTERNAL_OAUTH_JWS_KEYS_URL' '=' STRING | 'EXTERNAL_OAUTH_JWS_KEYS_URL' '=' '(' STRING (',' STRING)* ')') ('EXTERNAL_OAUTH_RSA_PUBLIC_KEY' '=' STRING)? ('EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2' '=' STRING)? ('EXTERNAL_OAUTH_BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_AUDIENCE_LIST' '=' '(' STRING ')')? ('EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'DISABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE_FOR_PRIVILEGE')? ('EXTERNAL_OAUTH_SCOPE_DELIMITER' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('NETWORK_POLICY' '=' STRING)? ('COMMENT' '=' STRING)? | 'EXTERNAL_OAUTH_RSA_PUBLIC_KEY' '=' STRING ('EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2' '=' STRING)? ('EXTERNAL_OAUTH_BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_AUDIENCE_LIST' '=' '(' STRING ')')? ('EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'DISABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE_FOR_PRIVILEGE')? ('EXTERNAL_OAUTH_SCOPE_DELIMITER' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('NETWORK_POLICY' '=' STRING)? ('COMMENT' '=' STRING)? | 'EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2' '=' STRING ('EXTERNAL_OAUTH_BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_AUDIENCE_LIST' '=' '(' STRING ')')? ('EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'DISABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE_FOR_PRIVILEGE')? ('EXTERNAL_OAUTH_SCOPE_DELIMITER' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('NETWORK_POLICY' '=' STRING)? ('COMMENT' '=' STRING)? | 'EXTERNAL_OAUTH_BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')' ('EXTERNAL_OAUTH_ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_AUDIENCE_LIST' '=' '(' STRING ')')? ('EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'DISABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE_FOR_PRIVILEGE')? ('EXTERNAL_OAUTH_SCOPE_DELIMITER' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('NETWORK_POLICY' '=' STRING)? ('COMMENT' '=' STRING)? | 'EXTERNAL_OAUTH_ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')' ('EXTERNAL_OAUTH_AUDIENCE_LIST' '=' '(' STRING ')')? ('EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'DISABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE_FOR_PRIVILEGE')? ('EXTERNAL_OAUTH_SCOPE_DELIMITER' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('NETWORK_POLICY' '=' STRING)? ('COMMENT' '=' STRING)? | 'EXTERNAL_OAUTH_AUDIENCE_LIST' '=' '(' STRING ')' ('EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'DISABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE_FOR_PRIVILEGE')? ('EXTERNAL_OAUTH_SCOPE_DELIMITER' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('NETWORK_POLICY' '=' STRING)? ('COMMENT' '=' STRING)? | ('EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'DISABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE' | 'EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' 'ENABLE_FOR_PRIVILEGE') ('EXTERNAL_OAUTH_SCOPE_DELIMITER' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('NETWORK_POLICY' '=' STRING)? ('COMMENT' '=' STRING)? | 'EXTERNAL_OAUTH_SCOPE_DELIMITER' '=' STRING ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('NETWORK_POLICY' '=' STRING)? ('COMMENT' '=' STRING)? | 'IS_AGENTIC' '=' ('TRUE' | 'FALSE') ('NETWORK_POLICY' '=' STRING)? ('COMMENT' '=' STRING)? | 'NETWORK_POLICY' '=' STRING ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'SECURITY'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'UNSET' ('ENABLED' | 'EXTERNAL_OAUTH_AUDIENCE_LIST' | 'IS_AGENTIC') (',' ('ENABLED' | 'EXTERNAL_OAUTH_AUDIENCE_LIST' | 'IS_AGENTIC'))*
    | 'ALTER' 'SECURITY'? 'INTEGRATION' ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'SECURITY'? 'INTEGRATION' ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-oauth-snowflake
alterSecurityIntegrationOauthSnowflake
    : 'ALTER' 'SECURITY'? 'INTEGRATION' ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'SECURITY'? 'INTEGRATION' ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'SECURITY'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'SET' ('ENABLED' '=' ('TRUE' | 'FALSE') ('OAUTH_ISSUE_REFRESH_TOKENS' '=' ('TRUE' | 'FALSE'))? ('OAUTH_REDIRECT_URI' '=' STRING)? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('NETWORK_POLICY' '=' STRING)? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'OAUTH_ISSUE_REFRESH_TOKENS' '=' ('TRUE' | 'FALSE') ('OAUTH_REDIRECT_URI' '=' STRING)? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('NETWORK_POLICY' '=' STRING)? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'OAUTH_REDIRECT_URI' '=' STRING ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('NETWORK_POLICY' '=' STRING)? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('NETWORK_POLICY' '=' STRING)? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('NETWORK_POLICY' '=' STRING)? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE') ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('NETWORK_POLICY' '=' STRING)? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE') ('ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('NETWORK_POLICY' '=' STRING)? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')' ('BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('NETWORK_POLICY' '=' STRING)? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')' ('NETWORK_POLICY' '=' STRING)? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'NETWORK_POLICY' '=' STRING ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'SECURITY'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'REFRESH' ('OAUTH_CLIENT_SECRET' | 'OAUTH_CLIENT_SECRET_2')
    | 'ALTER' 'SECURITY'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'UNSET' ('ENABLED' | 'COMMENT') (',' ('ENABLED' | 'COMMENT'))*
    | 'ALTER' 'SECURITY'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'SET' ('ENABLED' '=' ('TRUE' | 'FALSE') ('OAUTH_REDIRECT_URI' '=' STRING)? ('OAUTH_ALLOW_NON_TLS_REDIRECT_URI' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ALTERNATE_REDIRECT_URIS' '=' '(' STRING (',' STRING)* ')')? ('OAUTH_ENFORCE_PKCE' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ENABLE_ROLE_SELECTION' '=' ('TRUE' | 'FALSE'))? ('PRE_AUTHORIZED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('OAUTH_ANY_ROLE_MODE' '=' ('DISABLE' | 'ENABLE' | 'ENABLE_FOR_PRIVILEGE'))? ('OAUTH_ISSUE_REFRESH_TOKENS' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('NETWORK_POLICY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY_2' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'OAUTH_REDIRECT_URI' '=' STRING ('OAUTH_ALLOW_NON_TLS_REDIRECT_URI' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ALTERNATE_REDIRECT_URIS' '=' '(' STRING (',' STRING)* ')')? ('OAUTH_ENFORCE_PKCE' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ENABLE_ROLE_SELECTION' '=' ('TRUE' | 'FALSE'))? ('PRE_AUTHORIZED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('OAUTH_ANY_ROLE_MODE' '=' ('DISABLE' | 'ENABLE' | 'ENABLE_FOR_PRIVILEGE'))? ('OAUTH_ISSUE_REFRESH_TOKENS' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('NETWORK_POLICY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY_2' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'OAUTH_ALLOW_NON_TLS_REDIRECT_URI' '=' ('TRUE' | 'FALSE') ('OAUTH_ALTERNATE_REDIRECT_URIS' '=' '(' STRING (',' STRING)* ')')? ('OAUTH_ENFORCE_PKCE' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ENABLE_ROLE_SELECTION' '=' ('TRUE' | 'FALSE'))? ('PRE_AUTHORIZED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('OAUTH_ANY_ROLE_MODE' '=' ('DISABLE' | 'ENABLE' | 'ENABLE_FOR_PRIVILEGE'))? ('OAUTH_ISSUE_REFRESH_TOKENS' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('NETWORK_POLICY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY_2' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'OAUTH_ALTERNATE_REDIRECT_URIS' '=' '(' STRING (',' STRING)* ')' ('OAUTH_ENFORCE_PKCE' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ENABLE_ROLE_SELECTION' '=' ('TRUE' | 'FALSE'))? ('PRE_AUTHORIZED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('OAUTH_ANY_ROLE_MODE' '=' ('DISABLE' | 'ENABLE' | 'ENABLE_FOR_PRIVILEGE'))? ('OAUTH_ISSUE_REFRESH_TOKENS' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('NETWORK_POLICY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY_2' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'OAUTH_ENFORCE_PKCE' '=' ('TRUE' | 'FALSE') ('OAUTH_ENABLE_ROLE_SELECTION' '=' ('TRUE' | 'FALSE'))? ('PRE_AUTHORIZED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('OAUTH_ANY_ROLE_MODE' '=' ('DISABLE' | 'ENABLE' | 'ENABLE_FOR_PRIVILEGE'))? ('OAUTH_ISSUE_REFRESH_TOKENS' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('NETWORK_POLICY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY_2' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'OAUTH_ENABLE_ROLE_SELECTION' '=' ('TRUE' | 'FALSE') ('PRE_AUTHORIZED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('OAUTH_ANY_ROLE_MODE' '=' ('DISABLE' | 'ENABLE' | 'ENABLE_FOR_PRIVILEGE'))? ('OAUTH_ISSUE_REFRESH_TOKENS' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('NETWORK_POLICY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY_2' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'PRE_AUTHORIZED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')' ('ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('OAUTH_ANY_ROLE_MODE' '=' ('DISABLE' | 'ENABLE' | 'ENABLE_FOR_PRIVILEGE'))? ('OAUTH_ISSUE_REFRESH_TOKENS' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('NETWORK_POLICY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY_2' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')' ('BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('OAUTH_ANY_ROLE_MODE' '=' ('DISABLE' | 'ENABLE' | 'ENABLE_FOR_PRIVILEGE'))? ('OAUTH_ISSUE_REFRESH_TOKENS' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('NETWORK_POLICY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY_2' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')' ('OAUTH_ANY_ROLE_MODE' '=' ('DISABLE' | 'ENABLE' | 'ENABLE_FOR_PRIVILEGE'))? ('OAUTH_ISSUE_REFRESH_TOKENS' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('NETWORK_POLICY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY_2' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'OAUTH_ANY_ROLE_MODE' '=' ('DISABLE' | 'ENABLE' | 'ENABLE_FOR_PRIVILEGE') ('OAUTH_ISSUE_REFRESH_TOKENS' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('NETWORK_POLICY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY_2' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'OAUTH_ISSUE_REFRESH_TOKENS' '=' ('TRUE' | 'FALSE') ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('NETWORK_POLICY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY_2' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('NETWORK_POLICY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY_2' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('NETWORK_POLICY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY_2' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE') ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('NETWORK_POLICY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY_2' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE') ('NETWORK_POLICY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY_2' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'NETWORK_POLICY' '=' STRING ('OAUTH_CLIENT_RSA_PUBLIC_KEY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY_2' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'OAUTH_CLIENT_RSA_PUBLIC_KEY' '=' STRING ('OAUTH_CLIENT_RSA_PUBLIC_KEY_2' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'OAUTH_CLIENT_RSA_PUBLIC_KEY_2' '=' STRING ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'IS_AGENTIC' '=' ('TRUE' | 'FALSE') ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'SECURITY'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'UNSET' ('ENABLED' | 'NETWORK_POLICY' | 'OAUTH_CLIENT_RSA_PUBLIC_KEY' | 'OAUTH_CLIENT_RSA_PUBLIC_KEY_2' | 'IS_AGENTIC' | 'OAUTH_USE_SECONDARY_ROLES' | 'COMMENT') (',' ('ENABLED' | 'NETWORK_POLICY' | 'OAUTH_CLIENT_RSA_PUBLIC_KEY' | 'OAUTH_CLIENT_RSA_PUBLIC_KEY_2' | 'IS_AGENTIC' | 'OAUTH_USE_SECONDARY_ROLES' | 'COMMENT'))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-oidc
// not converted (form): ... stands for text the docs leave out
//   ALTER [ SECURITY ] INTEGRATION [ IF EXISTS ] <name> UNSET {
//       ALLOWED_EMAIL_PATTERNS |
//       ALLOWED_USER_DOMAINS |
//       COMMENT |
//       ENABLED |
//       OIDC_ENABLE_SSO_LOGIN_PAGE |
//       OIDC_LOGIN_PAGE_LABEL |
//       OIDC_SCOPES |
//       OIDC_USERINFO_ENDPOINT |
//       [ , ... ]
//       }
alterSecurityIntegrationOidc
    : 'ALTER' 'SECURITY'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'SET' ('ENABLED' '=' ('TRUE' | 'FALSE') ('OIDC_CLIENT_SECRET' '=' STRING)? ('OIDC_SCOPES' '=' '(' STRING (',' STRING)* ')')? ('OIDC_USERINFO_ENDPOINT' '=' STRING)? ('OIDC_ENABLE_SSO_LOGIN_PAGE' '=' ('TRUE' | 'FALSE'))? ('OIDC_LOGIN_PAGE_LABEL' '=' STRING)? ('OIDC_TOKEN_USER_MAPPING_CLAIM' '=' STRING)? ('OIDC_SNOWFLAKE_USER_MAPPING_ATTRIBUTE' '=' ('\'LOGIN_NAME\'' | '\'EMAIL_ADDRESS\''))? ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('OIDC_ALLOWED_GOOGLE_HOSTED_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('OIDC_ALLOWED_MSFT_TENANTS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'OIDC_CLIENT_SECRET' '=' STRING ('OIDC_SCOPES' '=' '(' STRING (',' STRING)* ')')? ('OIDC_USERINFO_ENDPOINT' '=' STRING)? ('OIDC_ENABLE_SSO_LOGIN_PAGE' '=' ('TRUE' | 'FALSE'))? ('OIDC_LOGIN_PAGE_LABEL' '=' STRING)? ('OIDC_TOKEN_USER_MAPPING_CLAIM' '=' STRING)? ('OIDC_SNOWFLAKE_USER_MAPPING_ATTRIBUTE' '=' ('\'LOGIN_NAME\'' | '\'EMAIL_ADDRESS\''))? ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('OIDC_ALLOWED_GOOGLE_HOSTED_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('OIDC_ALLOWED_MSFT_TENANTS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'OIDC_SCOPES' '=' '(' STRING (',' STRING)* ')' ('OIDC_USERINFO_ENDPOINT' '=' STRING)? ('OIDC_ENABLE_SSO_LOGIN_PAGE' '=' ('TRUE' | 'FALSE'))? ('OIDC_LOGIN_PAGE_LABEL' '=' STRING)? ('OIDC_TOKEN_USER_MAPPING_CLAIM' '=' STRING)? ('OIDC_SNOWFLAKE_USER_MAPPING_ATTRIBUTE' '=' ('\'LOGIN_NAME\'' | '\'EMAIL_ADDRESS\''))? ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('OIDC_ALLOWED_GOOGLE_HOSTED_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('OIDC_ALLOWED_MSFT_TENANTS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'OIDC_USERINFO_ENDPOINT' '=' STRING ('OIDC_ENABLE_SSO_LOGIN_PAGE' '=' ('TRUE' | 'FALSE'))? ('OIDC_LOGIN_PAGE_LABEL' '=' STRING)? ('OIDC_TOKEN_USER_MAPPING_CLAIM' '=' STRING)? ('OIDC_SNOWFLAKE_USER_MAPPING_ATTRIBUTE' '=' ('\'LOGIN_NAME\'' | '\'EMAIL_ADDRESS\''))? ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('OIDC_ALLOWED_GOOGLE_HOSTED_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('OIDC_ALLOWED_MSFT_TENANTS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'OIDC_ENABLE_SSO_LOGIN_PAGE' '=' ('TRUE' | 'FALSE') ('OIDC_LOGIN_PAGE_LABEL' '=' STRING)? ('OIDC_TOKEN_USER_MAPPING_CLAIM' '=' STRING)? ('OIDC_SNOWFLAKE_USER_MAPPING_ATTRIBUTE' '=' ('\'LOGIN_NAME\'' | '\'EMAIL_ADDRESS\''))? ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('OIDC_ALLOWED_GOOGLE_HOSTED_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('OIDC_ALLOWED_MSFT_TENANTS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'OIDC_LOGIN_PAGE_LABEL' '=' STRING ('OIDC_TOKEN_USER_MAPPING_CLAIM' '=' STRING)? ('OIDC_SNOWFLAKE_USER_MAPPING_ATTRIBUTE' '=' ('\'LOGIN_NAME\'' | '\'EMAIL_ADDRESS\''))? ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('OIDC_ALLOWED_GOOGLE_HOSTED_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('OIDC_ALLOWED_MSFT_TENANTS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'OIDC_TOKEN_USER_MAPPING_CLAIM' '=' STRING ('OIDC_SNOWFLAKE_USER_MAPPING_ATTRIBUTE' '=' ('\'LOGIN_NAME\'' | '\'EMAIL_ADDRESS\''))? ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('OIDC_ALLOWED_GOOGLE_HOSTED_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('OIDC_ALLOWED_MSFT_TENANTS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'OIDC_SNOWFLAKE_USER_MAPPING_ATTRIBUTE' '=' ('\'LOGIN_NAME\'' | '\'EMAIL_ADDRESS\'') ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('OIDC_ALLOWED_GOOGLE_HOSTED_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('OIDC_ALLOWED_MSFT_TENANTS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')' ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('OIDC_ALLOWED_GOOGLE_HOSTED_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('OIDC_ALLOWED_MSFT_TENANTS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')' ('OIDC_ALLOWED_GOOGLE_HOSTED_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('OIDC_ALLOWED_MSFT_TENANTS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'OIDC_ALLOWED_GOOGLE_HOSTED_DOMAINS' '=' '(' STRING (',' STRING)* ')' ('OIDC_ALLOWED_MSFT_TENANTS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'OIDC_ALLOWED_MSFT_TENANTS' '=' '(' STRING (',' STRING)* ')' ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'SECURITY'? 'INTEGRATION' ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'SECURITY'? 'INTEGRATION' ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-saml2
// not converted (form): ... stands for text the docs leave out
//   ALTER [ SECURITY ] INTEGRATION [ IF EXISTS ] <name> UNSET {
//       ENABLED |
//       [ , ... ]
//       }
alterSecurityIntegrationSaml2
    : 'ALTER' 'SECURITY'? 'INTEGRATION' ('IF' 'EXISTS')? QNAME 'SET' ('TYPE' '=' 'SAML2' ('ENABLED' '=' ('TRUE' | 'FALSE'))? ('METADATA_URL' '=' STRING)? ('SAML2_ISSUER' '=' STRING)? ('SAML2_SSO_URL' '=' STRING)? ('SAML2_PROVIDER' '=' STRING)? ('SAML2_X509_CERT' '=' STRING)? ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('SAML2_SP_INITIATED_LOGIN_PAGE_LABEL' '=' STRING)? ('SAML2_ENABLE_SP_INITIATED' '=' 'TRUE' | 'SAML2_ENABLE_SP_INITIATED' '=' 'FALSE')? ('SAML2_SNOWFLAKE_X509_CERT' '=' STRING)? ('SAML2_SIGN_REQUEST' '=' 'TRUE' | 'SAML2_SIGN_REQUEST' '=' 'FALSE')? ('SAML2_REQUESTED_NAMEID_FORMAT' '=' STRING)? ('SAML2_POST_LOGOUT_REDIRECT_URL' '=' STRING)? ('SAML2_FORCE_AUTHN' '=' 'TRUE' | 'SAML2_FORCE_AUTHN' '=' 'FALSE')? ('SAML2_SNOWFLAKE_ISSUER_URL' '=' STRING)? ('SAML2_SNOWFLAKE_ACS_URL' '=' STRING)? ('COMMENT' '=' STRING)? | 'ENABLED' '=' ('TRUE' | 'FALSE') ('METADATA_URL' '=' STRING)? ('SAML2_ISSUER' '=' STRING)? ('SAML2_SSO_URL' '=' STRING)? ('SAML2_PROVIDER' '=' STRING)? ('SAML2_X509_CERT' '=' STRING)? ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('SAML2_SP_INITIATED_LOGIN_PAGE_LABEL' '=' STRING)? ('SAML2_ENABLE_SP_INITIATED' '=' 'TRUE' | 'SAML2_ENABLE_SP_INITIATED' '=' 'FALSE')? ('SAML2_SNOWFLAKE_X509_CERT' '=' STRING)? ('SAML2_SIGN_REQUEST' '=' 'TRUE' | 'SAML2_SIGN_REQUEST' '=' 'FALSE')? ('SAML2_REQUESTED_NAMEID_FORMAT' '=' STRING)? ('SAML2_POST_LOGOUT_REDIRECT_URL' '=' STRING)? ('SAML2_FORCE_AUTHN' '=' 'TRUE' | 'SAML2_FORCE_AUTHN' '=' 'FALSE')? ('SAML2_SNOWFLAKE_ISSUER_URL' '=' STRING)? ('SAML2_SNOWFLAKE_ACS_URL' '=' STRING)? ('COMMENT' '=' STRING)? | 'METADATA_URL' '=' STRING ('SAML2_ISSUER' '=' STRING)? ('SAML2_SSO_URL' '=' STRING)? ('SAML2_PROVIDER' '=' STRING)? ('SAML2_X509_CERT' '=' STRING)? ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('SAML2_SP_INITIATED_LOGIN_PAGE_LABEL' '=' STRING)? ('SAML2_ENABLE_SP_INITIATED' '=' 'TRUE' | 'SAML2_ENABLE_SP_INITIATED' '=' 'FALSE')? ('SAML2_SNOWFLAKE_X509_CERT' '=' STRING)? ('SAML2_SIGN_REQUEST' '=' 'TRUE' | 'SAML2_SIGN_REQUEST' '=' 'FALSE')? ('SAML2_REQUESTED_NAMEID_FORMAT' '=' STRING)? ('SAML2_POST_LOGOUT_REDIRECT_URL' '=' STRING)? ('SAML2_FORCE_AUTHN' '=' 'TRUE' | 'SAML2_FORCE_AUTHN' '=' 'FALSE')? ('SAML2_SNOWFLAKE_ISSUER_URL' '=' STRING)? ('SAML2_SNOWFLAKE_ACS_URL' '=' STRING)? ('COMMENT' '=' STRING)? | 'SAML2_ISSUER' '=' STRING ('SAML2_SSO_URL' '=' STRING)? ('SAML2_PROVIDER' '=' STRING)? ('SAML2_X509_CERT' '=' STRING)? ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('SAML2_SP_INITIATED_LOGIN_PAGE_LABEL' '=' STRING)? ('SAML2_ENABLE_SP_INITIATED' '=' 'TRUE' | 'SAML2_ENABLE_SP_INITIATED' '=' 'FALSE')? ('SAML2_SNOWFLAKE_X509_CERT' '=' STRING)? ('SAML2_SIGN_REQUEST' '=' 'TRUE' | 'SAML2_SIGN_REQUEST' '=' 'FALSE')? ('SAML2_REQUESTED_NAMEID_FORMAT' '=' STRING)? ('SAML2_POST_LOGOUT_REDIRECT_URL' '=' STRING)? ('SAML2_FORCE_AUTHN' '=' 'TRUE' | 'SAML2_FORCE_AUTHN' '=' 'FALSE')? ('SAML2_SNOWFLAKE_ISSUER_URL' '=' STRING)? ('SAML2_SNOWFLAKE_ACS_URL' '=' STRING)? ('COMMENT' '=' STRING)? | 'SAML2_SSO_URL' '=' STRING ('SAML2_PROVIDER' '=' STRING)? ('SAML2_X509_CERT' '=' STRING)? ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('SAML2_SP_INITIATED_LOGIN_PAGE_LABEL' '=' STRING)? ('SAML2_ENABLE_SP_INITIATED' '=' 'TRUE' | 'SAML2_ENABLE_SP_INITIATED' '=' 'FALSE')? ('SAML2_SNOWFLAKE_X509_CERT' '=' STRING)? ('SAML2_SIGN_REQUEST' '=' 'TRUE' | 'SAML2_SIGN_REQUEST' '=' 'FALSE')? ('SAML2_REQUESTED_NAMEID_FORMAT' '=' STRING)? ('SAML2_POST_LOGOUT_REDIRECT_URL' '=' STRING)? ('SAML2_FORCE_AUTHN' '=' 'TRUE' | 'SAML2_FORCE_AUTHN' '=' 'FALSE')? ('SAML2_SNOWFLAKE_ISSUER_URL' '=' STRING)? ('SAML2_SNOWFLAKE_ACS_URL' '=' STRING)? ('COMMENT' '=' STRING)? | 'SAML2_PROVIDER' '=' STRING ('SAML2_X509_CERT' '=' STRING)? ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('SAML2_SP_INITIATED_LOGIN_PAGE_LABEL' '=' STRING)? ('SAML2_ENABLE_SP_INITIATED' '=' 'TRUE' | 'SAML2_ENABLE_SP_INITIATED' '=' 'FALSE')? ('SAML2_SNOWFLAKE_X509_CERT' '=' STRING)? ('SAML2_SIGN_REQUEST' '=' 'TRUE' | 'SAML2_SIGN_REQUEST' '=' 'FALSE')? ('SAML2_REQUESTED_NAMEID_FORMAT' '=' STRING)? ('SAML2_POST_LOGOUT_REDIRECT_URL' '=' STRING)? ('SAML2_FORCE_AUTHN' '=' 'TRUE' | 'SAML2_FORCE_AUTHN' '=' 'FALSE')? ('SAML2_SNOWFLAKE_ISSUER_URL' '=' STRING)? ('SAML2_SNOWFLAKE_ACS_URL' '=' STRING)? ('COMMENT' '=' STRING)? | 'SAML2_X509_CERT' '=' STRING ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('SAML2_SP_INITIATED_LOGIN_PAGE_LABEL' '=' STRING)? ('SAML2_ENABLE_SP_INITIATED' '=' 'TRUE' | 'SAML2_ENABLE_SP_INITIATED' '=' 'FALSE')? ('SAML2_SNOWFLAKE_X509_CERT' '=' STRING)? ('SAML2_SIGN_REQUEST' '=' 'TRUE' | 'SAML2_SIGN_REQUEST' '=' 'FALSE')? ('SAML2_REQUESTED_NAMEID_FORMAT' '=' STRING)? ('SAML2_POST_LOGOUT_REDIRECT_URL' '=' STRING)? ('SAML2_FORCE_AUTHN' '=' 'TRUE' | 'SAML2_FORCE_AUTHN' '=' 'FALSE')? ('SAML2_SNOWFLAKE_ISSUER_URL' '=' STRING)? ('SAML2_SNOWFLAKE_ACS_URL' '=' STRING)? ('COMMENT' '=' STRING)? | 'ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')' ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('SAML2_SP_INITIATED_LOGIN_PAGE_LABEL' '=' STRING)? ('SAML2_ENABLE_SP_INITIATED' '=' 'TRUE' | 'SAML2_ENABLE_SP_INITIATED' '=' 'FALSE')? ('SAML2_SNOWFLAKE_X509_CERT' '=' STRING)? ('SAML2_SIGN_REQUEST' '=' 'TRUE' | 'SAML2_SIGN_REQUEST' '=' 'FALSE')? ('SAML2_REQUESTED_NAMEID_FORMAT' '=' STRING)? ('SAML2_POST_LOGOUT_REDIRECT_URL' '=' STRING)? ('SAML2_FORCE_AUTHN' '=' 'TRUE' | 'SAML2_FORCE_AUTHN' '=' 'FALSE')? ('SAML2_SNOWFLAKE_ISSUER_URL' '=' STRING)? ('SAML2_SNOWFLAKE_ACS_URL' '=' STRING)? ('COMMENT' '=' STRING)? | 'ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')' ('SAML2_SP_INITIATED_LOGIN_PAGE_LABEL' '=' STRING)? ('SAML2_ENABLE_SP_INITIATED' '=' 'TRUE' | 'SAML2_ENABLE_SP_INITIATED' '=' 'FALSE')? ('SAML2_SNOWFLAKE_X509_CERT' '=' STRING)? ('SAML2_SIGN_REQUEST' '=' 'TRUE' | 'SAML2_SIGN_REQUEST' '=' 'FALSE')? ('SAML2_REQUESTED_NAMEID_FORMAT' '=' STRING)? ('SAML2_POST_LOGOUT_REDIRECT_URL' '=' STRING)? ('SAML2_FORCE_AUTHN' '=' 'TRUE' | 'SAML2_FORCE_AUTHN' '=' 'FALSE')? ('SAML2_SNOWFLAKE_ISSUER_URL' '=' STRING)? ('SAML2_SNOWFLAKE_ACS_URL' '=' STRING)? ('COMMENT' '=' STRING)? | 'SAML2_SP_INITIATED_LOGIN_PAGE_LABEL' '=' STRING ('SAML2_ENABLE_SP_INITIATED' '=' 'TRUE' | 'SAML2_ENABLE_SP_INITIATED' '=' 'FALSE')? ('SAML2_SNOWFLAKE_X509_CERT' '=' STRING)? ('SAML2_SIGN_REQUEST' '=' 'TRUE' | 'SAML2_SIGN_REQUEST' '=' 'FALSE')? ('SAML2_REQUESTED_NAMEID_FORMAT' '=' STRING)? ('SAML2_POST_LOGOUT_REDIRECT_URL' '=' STRING)? ('SAML2_FORCE_AUTHN' '=' 'TRUE' | 'SAML2_FORCE_AUTHN' '=' 'FALSE')? ('SAML2_SNOWFLAKE_ISSUER_URL' '=' STRING)? ('SAML2_SNOWFLAKE_ACS_URL' '=' STRING)? ('COMMENT' '=' STRING)? | ('SAML2_ENABLE_SP_INITIATED' '=' 'TRUE' | 'SAML2_ENABLE_SP_INITIATED' '=' 'FALSE') ('SAML2_SNOWFLAKE_X509_CERT' '=' STRING)? ('SAML2_SIGN_REQUEST' '=' 'TRUE' | 'SAML2_SIGN_REQUEST' '=' 'FALSE')? ('SAML2_REQUESTED_NAMEID_FORMAT' '=' STRING)? ('SAML2_POST_LOGOUT_REDIRECT_URL' '=' STRING)? ('SAML2_FORCE_AUTHN' '=' 'TRUE' | 'SAML2_FORCE_AUTHN' '=' 'FALSE')? ('SAML2_SNOWFLAKE_ISSUER_URL' '=' STRING)? ('SAML2_SNOWFLAKE_ACS_URL' '=' STRING)? ('COMMENT' '=' STRING)? | 'SAML2_SNOWFLAKE_X509_CERT' '=' STRING ('SAML2_SIGN_REQUEST' '=' 'TRUE' | 'SAML2_SIGN_REQUEST' '=' 'FALSE')? ('SAML2_REQUESTED_NAMEID_FORMAT' '=' STRING)? ('SAML2_POST_LOGOUT_REDIRECT_URL' '=' STRING)? ('SAML2_FORCE_AUTHN' '=' 'TRUE' | 'SAML2_FORCE_AUTHN' '=' 'FALSE')? ('SAML2_SNOWFLAKE_ISSUER_URL' '=' STRING)? ('SAML2_SNOWFLAKE_ACS_URL' '=' STRING)? ('COMMENT' '=' STRING)? | ('SAML2_SIGN_REQUEST' '=' 'TRUE' | 'SAML2_SIGN_REQUEST' '=' 'FALSE') ('SAML2_REQUESTED_NAMEID_FORMAT' '=' STRING)? ('SAML2_POST_LOGOUT_REDIRECT_URL' '=' STRING)? ('SAML2_FORCE_AUTHN' '=' 'TRUE' | 'SAML2_FORCE_AUTHN' '=' 'FALSE')? ('SAML2_SNOWFLAKE_ISSUER_URL' '=' STRING)? ('SAML2_SNOWFLAKE_ACS_URL' '=' STRING)? ('COMMENT' '=' STRING)? | 'SAML2_REQUESTED_NAMEID_FORMAT' '=' STRING ('SAML2_POST_LOGOUT_REDIRECT_URL' '=' STRING)? ('SAML2_FORCE_AUTHN' '=' 'TRUE' | 'SAML2_FORCE_AUTHN' '=' 'FALSE')? ('SAML2_SNOWFLAKE_ISSUER_URL' '=' STRING)? ('SAML2_SNOWFLAKE_ACS_URL' '=' STRING)? ('COMMENT' '=' STRING)? | 'SAML2_POST_LOGOUT_REDIRECT_URL' '=' STRING ('SAML2_FORCE_AUTHN' '=' 'TRUE' | 'SAML2_FORCE_AUTHN' '=' 'FALSE')? ('SAML2_SNOWFLAKE_ISSUER_URL' '=' STRING)? ('SAML2_SNOWFLAKE_ACS_URL' '=' STRING)? ('COMMENT' '=' STRING)? | ('SAML2_FORCE_AUTHN' '=' 'TRUE' | 'SAML2_FORCE_AUTHN' '=' 'FALSE') ('SAML2_SNOWFLAKE_ISSUER_URL' '=' STRING)? ('SAML2_SNOWFLAKE_ACS_URL' '=' STRING)? ('COMMENT' '=' STRING)? | 'SAML2_SNOWFLAKE_ISSUER_URL' '=' STRING ('SAML2_SNOWFLAKE_ACS_URL' '=' STRING)? ('COMMENT' '=' STRING)? | 'SAML2_SNOWFLAKE_ACS_URL' '=' STRING ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'SECURITY'? 'INTEGRATION' QNAME 'REFRESH' 'SAML2_SNOWFLAKE_PRIVATE_KEY'? 'METADATA_URL'?
    | 'ALTER' 'SECURITY'? 'INTEGRATION' QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'SECURITY'? 'INTEGRATION' QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-scim
// not converted (form): ... stands for text the docs leave out
//   ALTER [ SECURITY ] INTEGRATION [ IF EXISTS ] <name>  UNSET {
//                                                               NETWORK_POLICY |
//                                                               [ , ... ]
//                                                               }
alterSecurityIntegrationScim
    : 'ALTER' 'SECURITY'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'SET' ('ENABLED' '=' ('TRUE' | 'FALSE') ('NETWORK_POLICY' '=' STRING)? ('REJECT_TOKENS_ISSUED_BEFORE' '=' STRING)? ('SYNC_PASSWORD' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'NETWORK_POLICY' '=' STRING ('REJECT_TOKENS_ISSUED_BEFORE' '=' STRING)? ('SYNC_PASSWORD' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'REJECT_TOKENS_ISSUED_BEFORE' '=' STRING ('SYNC_PASSWORD' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? | 'SYNC_PASSWORD' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'SECURITY'? 'INTEGRATION' ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'SECURITY'? 'INTEGRATION' ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-semantic-view
alterSemanticView
    : 'ALTER' 'SEMANTIC' 'VIEW' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'SEMANTIC' 'VIEW' ('IF' 'EXISTS')? QNAME 'SET' ('COMMENT' '=' STRING | 'MAX_STALENESS' '=' NUMBER)
    | 'ALTER' 'SEMANTIC' 'VIEW' ('IF' 'EXISTS')? QNAME 'UNSET' ('COMMENT' | 'MAX_STALENESS')
    | 'ALTER' 'SEMANTIC' 'VIEW' QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'SEMANTIC' 'VIEW' QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'SEMANTIC' 'VIEW' QNAME 'ADD' 'MATERIALIZATION' ID 'WAREHOUSE' '=' ID ('REFRESH_MODE' '=' ('AUTO' | 'FULL' | 'INCREMENTAL'))? ('IMMUTABLE' 'WHERE' '(' EXPR ')')? 'AS' 'DIMENSIONS' ID (',' ID)* 'METRICS' ID (',' ID)* ('WHERE' '(' EXPR ')')?
    | 'ALTER' 'SEMANTIC' 'VIEW' QNAME 'DROP' 'MATERIALIZATION' ID
    | 'ALTER' 'SEMANTIC' 'VIEW' QNAME 'SUSPEND' 'MATERIALIZATION' ID
    | 'ALTER' 'SEMANTIC' 'VIEW' QNAME 'RESUME' 'MATERIALIZATION' ID
    | 'ALTER' 'SEMANTIC' 'VIEW' QNAME 'REFRESH' 'MATERIALIZATION' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-sequence
// fix: the increment may be negative, as the parameters say
alterSequence
    : 'ALTER' 'SEQUENCE' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'SEQUENCE' ('IF' 'EXISTS')? QNAME 'SET'? ('INCREMENT' 'BY'? '='? '-'? NUMBER)?
    | 'ALTER' 'SEQUENCE' ('IF' 'EXISTS')? QNAME 'SET' (('ORDER' | 'NOORDER') ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'SEQUENCE' ('IF' 'EXISTS')? QNAME 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-service
// fix: read: a template argument's <value> is a constant, a string or JSON value, as the USING parameters say
alterService
    : 'ALTER' 'SERVICE' ('IF' 'EXISTS')? QNAME ('SUSPEND' | 'RESUME')
    | 'ALTER' 'SERVICE' ('IF' 'EXISTS')? QNAME (alterService_fromSpecification | alterService_fromSpecificationTemplate)
    | 'ALTER' 'SERVICE' ('IF' 'EXISTS')? QNAME 'RESTORE' 'VOLUME' ID 'INSTANCES' ID 'FROM' 'SNAPSHOT' ID
    | 'ALTER' 'SERVICE' ('IF' 'EXISTS')? QNAME 'SET' ('MIN_INSTANCES' '=' NUMBER ('MAX_INSTANCES' '=' NUMBER)? ('LOG_LEVEL' '=' STRING)? ('AUTO_SUSPEND_SECS' '=' NUMBER)? ('MIN_READY_INSTANCES' '=' NUMBER)? ('QUERY_WAREHOUSE' '=' ID)? ('AUTO_RESUME' '=' ('TRUE' | 'FALSE'))? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('COMMENT' '=' STRING)? | 'MAX_INSTANCES' '=' NUMBER ('LOG_LEVEL' '=' STRING)? ('AUTO_SUSPEND_SECS' '=' NUMBER)? ('MIN_READY_INSTANCES' '=' NUMBER)? ('QUERY_WAREHOUSE' '=' ID)? ('AUTO_RESUME' '=' ('TRUE' | 'FALSE'))? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('COMMENT' '=' STRING)? | 'LOG_LEVEL' '=' STRING ('AUTO_SUSPEND_SECS' '=' NUMBER)? ('MIN_READY_INSTANCES' '=' NUMBER)? ('QUERY_WAREHOUSE' '=' ID)? ('AUTO_RESUME' '=' ('TRUE' | 'FALSE'))? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('COMMENT' '=' STRING)? | 'AUTO_SUSPEND_SECS' '=' NUMBER ('MIN_READY_INSTANCES' '=' NUMBER)? ('QUERY_WAREHOUSE' '=' ID)? ('AUTO_RESUME' '=' ('TRUE' | 'FALSE'))? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('COMMENT' '=' STRING)? | 'MIN_READY_INSTANCES' '=' NUMBER ('QUERY_WAREHOUSE' '=' ID)? ('AUTO_RESUME' '=' ('TRUE' | 'FALSE'))? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('COMMENT' '=' STRING)? | 'QUERY_WAREHOUSE' '=' ID ('AUTO_RESUME' '=' ('TRUE' | 'FALSE'))? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('COMMENT' '=' STRING)? | 'AUTO_RESUME' '=' ('TRUE' | 'FALSE') ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('COMMENT' '=' STRING)? | 'EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')' ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'SERVICE' ('IF' 'EXISTS')? QNAME 'UNSET' ('MIN_INSTANCES' | 'AUTO_SUSPEND_SECS' | 'MAX_INSTANCES' | 'LOG_LEVEL' | 'MIN_READY_INSTANCES' | 'QUERY_WAREHOUSE' | 'AUTO_RESUME' | 'EXTERNAL_ACCESS_INTEGRATIONS' | 'COMMENT') (',' ('MIN_INSTANCES' | 'AUTO_SUSPEND_SECS' | 'MAX_INSTANCES' | 'LOG_LEVEL' | 'MIN_READY_INSTANCES' | 'QUERY_WAREHOUSE' | 'AUTO_RESUME' | 'EXTERNAL_ACCESS_INTEGRATIONS' | 'COMMENT'))*
    | 'ALTER' 'SERVICE' ('IF' 'EXISTS')? QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    ;
alterService_fromSpecification
    : ('FROM' 'SPECIFICATION_FILE' '=' STRING | 'FROM' STAGE 'SPECIFICATION_FILE' '=' STRING | 'FROM' 'SPECIFICATION' STRING)
    ;
alterService_fromSpecificationTemplate
    : ('FROM' 'SPECIFICATION_TEMPLATE_FILE' '=' STRING | 'FROM' STAGE 'SPECIFICATION_TEMPLATE_FILE' '=' STRING | 'FROM' 'SPECIFICATION_TEMPLATE' STRING) 'USING' '(' ID '=>' EXPR (',' ID '=>' EXPR)* ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-session
// fix: the parameters may be separated by spaces or commas, as the SET parameter says
alterSession
    : 'ALTER' 'SESSION' 'SET' alterSession_sessionParams (',' alterSession_sessionParams)*
    | 'ALTER' 'SESSION' 'UNSET' ID (',' ID)*
    ;
alterSession_sessionParams
    : ('ABORT_DETACHED_QUERY' '=' ('TRUE' | 'FALSE') | 'ACTIVE_PYTHON_PROFILER' '=' ('\'LINE\'' | '\'MEMORY\'') | 'AUTOCOMMIT' '=' ('TRUE' | 'FALSE') | 'BINARY_INPUT_FORMAT' '=' STRING | 'BINARY_OUTPUT_FORMAT' '=' STRING | 'DATE_INPUT_FORMAT' '=' STRING | 'DATE_OUTPUT_FORMAT' '=' STRING | 'ERROR_ON_NONDETERMINISTIC_MERGE' '=' ('TRUE' | 'FALSE') | 'ERROR_ON_NONDETERMINISTIC_UPDATE' '=' ('TRUE' | 'FALSE') | 'GEOGRAPHY_OUTPUT_FORMAT' '=' ('\'GeoJSON\'' | '\'WKT\'' | '\'WKB\'' | '\'EWKT\'' | '\'EWKB\'') | 'HYBRID_TABLE_LOCK_TIMEOUT' '=' NUMBER | 'JSON_INDENT' '=' NUMBER | 'LOG_LEVEL' '=' STRING | 'LOCK_TIMEOUT' '=' NUMBER | 'OPT_OUT_ERROR_LOGGING' '=' ('TRUE' | 'FALSE') | 'PYTHON_PROFILER_TARGET_STAGE' '=' STRING | 'PYTHON_PROFILER_MODULES' '=' STRING | 'QUERY_TAG' '=' STRING | 'ROWS_PER_RESULTSET' '=' NUMBER | 'S3_STAGE_VPCE_DNS_NAME' '=' STRING | 'SEARCH_PATH' '=' STRING | 'SIMULATED_DATA_SHARING_CONSUMER' '=' STRING | 'STATEMENT_TIMEOUT_IN_SECONDS' '=' NUMBER | 'STRICT_JSON_OUTPUT' '=' ('TRUE' | 'FALSE') | 'TIMESTAMP_DAY_IS_ALWAYS_24H' '=' ('TRUE' | 'FALSE') | 'TIMESTAMP_INPUT_FORMAT' '=' STRING | 'TIMESTAMP_LTZ_OUTPUT_FORMAT' '=' STRING | 'TIMESTAMP_NTZ_OUTPUT_FORMAT' '=' STRING | 'TIMESTAMP_OUTPUT_FORMAT' '=' STRING | 'TIMESTAMP_TYPE_MAPPING' '=' STRING | 'TIMESTAMP_TZ_OUTPUT_FORMAT' '=' STRING | 'TIMEZONE' '=' STRING | 'TIME_INPUT_FORMAT' '=' STRING | 'TIME_OUTPUT_FORMAT' '=' STRING | 'TRACE_LEVEL' '=' STRING | 'TRANSACTION_DEFAULT_ISOLATION_LEVEL' '=' STRING | 'TWO_DIGIT_CENTURY_START' '=' NUMBER | 'UNSUPPORTED_DDL_ACTION' '=' STRING | 'USE_CACHED_RESULT' '=' ('TRUE' | 'FALSE') | 'WEEK_OF_YEAR_POLICY' '=' NUMBER | 'WEEK_START' '=' NUMBER)+
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-session-policy
alterSessionPolicy
    : 'ALTER' 'SESSION' 'POLICY' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'SESSION' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' ('SESSION_IDLE_TIMEOUT_MINS' '=' NUMBER ('SESSION_UI_IDLE_TIMEOUT_MINS' '=' NUMBER)? ('SESSION_MAX_LIFESPAN_MINS' '=' NUMBER)? ('SESSION_UI_MAX_LIFESPAN_MINS' '=' NUMBER)? ('ALLOWED_SECONDARY_ROLES' '=' '(' (('\'ALL\'' | ID (',' ID)*))? ')')? ('BLOCKED_SECONDARY_ROLES' '=' '(' (('\'ALL\'' | ID (',' ID)*))? ')')? ('AGENT_RESTRICTED_SESSION_SCOPE' '=' STRING)? ('COMMENT' '=' STRING)? | 'SESSION_UI_IDLE_TIMEOUT_MINS' '=' NUMBER ('SESSION_MAX_LIFESPAN_MINS' '=' NUMBER)? ('SESSION_UI_MAX_LIFESPAN_MINS' '=' NUMBER)? ('ALLOWED_SECONDARY_ROLES' '=' '(' (('\'ALL\'' | ID (',' ID)*))? ')')? ('BLOCKED_SECONDARY_ROLES' '=' '(' (('\'ALL\'' | ID (',' ID)*))? ')')? ('AGENT_RESTRICTED_SESSION_SCOPE' '=' STRING)? ('COMMENT' '=' STRING)? | 'SESSION_MAX_LIFESPAN_MINS' '=' NUMBER ('SESSION_UI_MAX_LIFESPAN_MINS' '=' NUMBER)? ('ALLOWED_SECONDARY_ROLES' '=' '(' (('\'ALL\'' | ID (',' ID)*))? ')')? ('BLOCKED_SECONDARY_ROLES' '=' '(' (('\'ALL\'' | ID (',' ID)*))? ')')? ('AGENT_RESTRICTED_SESSION_SCOPE' '=' STRING)? ('COMMENT' '=' STRING)? | 'SESSION_UI_MAX_LIFESPAN_MINS' '=' NUMBER ('ALLOWED_SECONDARY_ROLES' '=' '(' (('\'ALL\'' | ID (',' ID)*))? ')')? ('BLOCKED_SECONDARY_ROLES' '=' '(' (('\'ALL\'' | ID (',' ID)*))? ')')? ('AGENT_RESTRICTED_SESSION_SCOPE' '=' STRING)? ('COMMENT' '=' STRING)? | 'ALLOWED_SECONDARY_ROLES' '=' '(' (('\'ALL\'' | ID (',' ID)*))? ')' ('BLOCKED_SECONDARY_ROLES' '=' '(' (('\'ALL\'' | ID (',' ID)*))? ')')? ('AGENT_RESTRICTED_SESSION_SCOPE' '=' STRING)? ('COMMENT' '=' STRING)? | 'BLOCKED_SECONDARY_ROLES' '=' '(' (('\'ALL\'' | ID (',' ID)*))? ')' ('AGENT_RESTRICTED_SESSION_SCOPE' '=' STRING)? ('COMMENT' '=' STRING)? | 'AGENT_RESTRICTED_SESSION_SCOPE' '=' STRING ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'SESSION' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'SESSION' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'SESSION' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' ('SESSION_IDLE_TIMEOUT_MINS' | 'SESSION_UI_IDLE_TIMEOUT_MINS' | 'SESSION_MAX_LIFESPAN_MINS' | 'SESSION_UI_MAX_LIFESPAN_MINS' | 'ALLOWED_SECONDARY_ROLES' | 'BLOCKED_SECONDARY_ROLES' | 'AGENT_RESTRICTED_SESSION_SCOPE' | 'COMMENT') (',' ('SESSION_IDLE_TIMEOUT_MINS' | 'SESSION_UI_IDLE_TIMEOUT_MINS' | 'SESSION_MAX_LIFESPAN_MINS' | 'SESSION_UI_MAX_LIFESPAN_MINS' | 'ALLOWED_SECONDARY_ROLES' | 'BLOCKED_SECONDARY_ROLES' | 'AGENT_RESTRICTED_SESSION_SCOPE' | 'COMMENT'))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-share
alterShare
    : 'ALTER' 'SHARE' ('IF' 'EXISTS')? ID ('ADD' | 'REMOVE') 'ACCOUNTS' '=' ID (',' ID)* ('SHARE_RESTRICTIONS' '=' ('TRUE' | 'FALSE'))?
    | 'ALTER' 'SHARE' ('IF' 'EXISTS')? ID 'SET' ('ACCOUNTS' '=' ID (',' ID)* ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'SHARE' ('IF' 'EXISTS')? ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'SHARE' ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'SHARE' ('IF' 'EXISTS')? ID 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-snapshot
alterSnapshot
    : 'ALTER' 'SNAPSHOT' ('IF' 'EXISTS')? QNAME 'SET' 'COMMENT' '=' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-snapshot-policy
alterSnapshotPolicy
    : 'ALTER' 'SNAPSHOT' 'POLICY' QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'SNAPSHOT' 'POLICY' QNAME 'SET' ('COMMENT' '=' STRING ('SCHEDULE' '=' STRING)? ('EXPIRE_AFTER_DAYS' '=' NUMBER)? | 'SCHEDULE' '=' STRING ('EXPIRE_AFTER_DAYS' '=' NUMBER)? | 'EXPIRE_AFTER_DAYS' '=' NUMBER)
    | 'ALTER' 'SNAPSHOT' 'POLICY' QNAME 'UNSET' ('COMMENT' | 'SCHEDULE' | 'EXPIRE_AFTER_DAYS')
    | 'ALTER' 'SNAPSHOT' 'POLICY' QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'SNAPSHOT' 'POLICY' QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-snapshot-set
alterSnapshotSet
    : 'ALTER' 'SNAPSHOT' 'SET' QNAME 'ADD' 'SNAPSHOT'
    | 'ALTER' 'SNAPSHOT' 'SET' QNAME 'APPLY' 'SNAPSHOT' 'POLICY' QNAME 'FORCE'?
    | 'ALTER' 'SNAPSHOT' 'SET' QNAME 'SUSPEND' 'SNAPSHOT' (('CREATION' | 'EXPIRATION'))? 'POLICY'
    | 'ALTER' 'SNAPSHOT' 'SET' QNAME 'RESUME' 'SNAPSHOT' (('CREATION' | 'EXPIRATION'))? 'POLICY'
    | 'ALTER' 'SNAPSHOT' 'SET' QNAME 'DELETE' 'SNAPSHOT' 'IDENTIFIER' STRING
    | 'ALTER' 'SNAPSHOT' 'SET' QNAME 'MODIFY' 'SNAPSHOT' 'IDENTIFIER' STRING ('ADD' | 'REMOVE') 'LEGAL' 'HOLD'
    | 'ALTER' 'SNAPSHOT' 'SET' QNAME 'MODIFY' 'SNAPSHOT' 'IDENTIFIER' STRING 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'SNAPSHOT' 'SET' QNAME 'MODIFY' 'SNAPSHOT' 'IDENTIFIER' STRING 'UNSET' 'COMMENT'
    | 'ALTER' 'SNAPSHOT' 'SET' QNAME 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'SNAPSHOT' 'SET' QNAME 'UNSET' 'COMMENT'
    | 'ALTER' 'SNAPSHOT' 'SET' QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'SNAPSHOT' 'SET' QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-stage
// fix: externalStageParams (for Google Cloud Storage): a } inside [ STORAGE_INTEGRATION = ... ] closes no bracket; dropped
// fix: an ALTER may leave the URL as it is and change the credentials or encryption, as the page's examples do
// fix: read: each -- If TYPE = X section of formatTypeOptions is the options of TYPE = X alone, and TYPE = { ... } [ formatTypeOptions ] chooses a TYPE with its own options
alterStage
    : 'ALTER' 'STAGE' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'STAGE' ('IF' 'EXISTS')? QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'STAGE' QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'STAGE' ('IF' 'EXISTS')? QNAME 'SET' ('FILE_FORMAT' '=' '(' ('FORMAT_NAME' '=' STRING | ('TYPE' '=' 'CSV' alterStage_formatTypeOptions_CSV? | 'TYPE' '=' 'JSON' alterStage_formatTypeOptions_JSON? | 'TYPE' '=' 'AVRO' alterStage_formatTypeOptions_AVRO? | 'TYPE' '=' 'ORC' alterStage_formatTypeOptions_ORC? | 'TYPE' '=' 'PARQUET' alterStage_formatTypeOptions_PARQUET? | 'TYPE' '=' 'XML' alterStage_formatTypeOptions_XML? | 'TYPE' '=' 'CUSTOM')) ')')? ('STAGE_UNDROP_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('COMMENT' '=' STRING)?
    | 'ALTER' 'STAGE' ('IF' 'EXISTS')? QNAME 'UNSET' 'STAGE_UNDROP_RETENTION_TIME_IN_DAYS'
    | 'ALTER' 'STAGE' ('IF' 'EXISTS')? QNAME 'SET' (alterStage_externalStageParams ('FILE_FORMAT' '=' '(' ('FORMAT_NAME' '=' STRING | ('TYPE' '=' 'CSV' alterStage_formatTypeOptions_CSV? | 'TYPE' '=' 'JSON' alterStage_formatTypeOptions_JSON? | 'TYPE' '=' 'AVRO' alterStage_formatTypeOptions_AVRO? | 'TYPE' '=' 'ORC' alterStage_formatTypeOptions_ORC? | 'TYPE' '=' 'PARQUET' alterStage_formatTypeOptions_PARQUET? | 'TYPE' '=' 'XML' alterStage_formatTypeOptions_XML? | 'TYPE' '=' 'CUSTOM')) ')')? ('COMMENT' '=' STRING)? | 'FILE_FORMAT' '=' '(' ('FORMAT_NAME' '=' STRING | ('TYPE' '=' 'CSV' alterStage_formatTypeOptions_CSV? | 'TYPE' '=' 'JSON' alterStage_formatTypeOptions_JSON? | 'TYPE' '=' 'AVRO' alterStage_formatTypeOptions_AVRO? | 'TYPE' '=' 'ORC' alterStage_formatTypeOptions_ORC? | 'TYPE' '=' 'PARQUET' alterStage_formatTypeOptions_PARQUET? | 'TYPE' '=' 'XML' alterStage_formatTypeOptions_XML? | 'TYPE' '=' 'CUSTOM')) ')' ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    ;
alterStage_externalStageParams
    : ('URL' '=' STRING)? ('AWS_ACCESS_POINT_ARN' '=' STRING)? ('STORAGE_INTEGRATION' '=' ID | 'CREDENTIALS' '=' '(' ('AWS_KEY_ID' '=' STRING 'AWS_SECRET_KEY' '=' STRING ('AWS_TOKEN' '=' STRING)? | 'AWS_ROLE' '=' STRING) ')')? ('ENCRYPTION' '=' '(' (('TYPE' '=' '\'AWS_CSE\'')? 'MASTER_KEY' '=' STRING | 'TYPE' '=' '\'AWS_SSE_S3\'' | 'TYPE' '=' '\'AWS_SSE_KMS\'' ('KMS_KEY_ID' '=' STRING)? | 'TYPE' '=' '\'NONE\'') ')')? ('USE_PRIVATELINK_ENDPOINT' '=' ('TRUE' | 'FALSE'))?
    | ('URL' '=' STRING)? ('STORAGE_INTEGRATION' '=' ID)? ('ENCRYPTION' '=' '(' ('TYPE' '=' '\'GCS_SSE_KMS\'' ('KMS_KEY_ID' '=' STRING)? | 'TYPE' '=' '\'NONE\'') ')')?
    | ('URL' '=' STRING)? ('STORAGE_INTEGRATION' '=' ID | 'CREDENTIALS' '=' '(' ('AZURE_SAS_TOKEN' '=' STRING)? ')')? ('ENCRYPTION' '=' '(' ('TYPE' '=' '\'AZURE_CSE\'' ('MASTER_KEY' '=' STRING)? | 'TYPE' '=' '\'NONE\'') ')')? ('USE_PRIVATELINK_ENDPOINT' '=' ('TRUE' | 'FALSE'))?
    ;
alterStage_formatTypeOptions_CSV
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'RECORD_DELIMITER' '=' (STRING | 'NONE') | 'FIELD_DELIMITER' '=' (STRING | 'NONE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'FILE_EXTENSION' '=' STRING | 'PARSE_HEADER' '=' ('TRUE' | 'FALSE') | 'SKIP_HEADER' '=' NUMBER | 'SKIP_BLANK_LINES' '=' ('TRUE' | 'FALSE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'ESCAPE' '=' (STRING | 'NONE') | 'ESCAPE_UNENCLOSED_FIELD' '=' (STRING | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'FIELD_OPTIONALLY_ENCLOSED_BY' '=' (STRING | 'NONE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')' | 'ERROR_ON_COLUMN_COUNT_MISMATCH' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'EMPTY_FIELD_AS_NULL' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE') | 'ENCODING' '=' (STRING | 'UTF8')) (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'RECORD_DELIMITER' '=' (STRING | 'NONE') | 'FIELD_DELIMITER' '=' (STRING | 'NONE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'FILE_EXTENSION' '=' STRING | 'PARSE_HEADER' '=' ('TRUE' | 'FALSE') | 'SKIP_HEADER' '=' NUMBER | 'SKIP_BLANK_LINES' '=' ('TRUE' | 'FALSE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'ESCAPE' '=' (STRING | 'NONE') | 'ESCAPE_UNENCLOSED_FIELD' '=' (STRING | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'FIELD_OPTIONALLY_ENCLOSED_BY' '=' (STRING | 'NONE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')' | 'ERROR_ON_COLUMN_COUNT_MISMATCH' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'EMPTY_FIELD_AS_NULL' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE') | 'ENCODING' '=' (STRING | 'UTF8')))*
    ;
alterStage_formatTypeOptions_JSON
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')' | 'FILE_EXTENSION' '=' STRING | 'ENABLE_OCTAL' '=' ('TRUE' | 'FALSE') | 'ALLOW_DUPLICATE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ARRAY' '=' ('TRUE' | 'FALSE') | 'STRIP_NULL_VALUES' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'IGNORE_UTF8_ERRORS' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE')) (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')' | 'FILE_EXTENSION' '=' STRING | 'ENABLE_OCTAL' '=' ('TRUE' | 'FALSE') | 'ALLOW_DUPLICATE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ARRAY' '=' ('TRUE' | 'FALSE') | 'STRIP_NULL_VALUES' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'IGNORE_UTF8_ERRORS' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE')))*
    ;
alterStage_formatTypeOptions_AVRO
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')') (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')'))*
    ;
alterStage_formatTypeOptions_ORC
    : ('TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')') (','? ('TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')'))*
    ;
alterStage_formatTypeOptions_PARQUET
    : ('COMPRESSION' '=' ('AUTO' | 'LZO' | 'SNAPPY' | 'NONE') | 'SNAPPY_COMPRESSION' '=' ('TRUE' | 'FALSE') | 'BINARY_AS_TEXT' '=' ('TRUE' | 'FALSE') | 'USE_LOGICAL_TYPE' '=' ('TRUE' | 'FALSE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'USE_VECTORIZED_SCANNER' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')') (','? ('COMPRESSION' '=' ('AUTO' | 'LZO' | 'SNAPPY' | 'NONE') | 'SNAPPY_COMPRESSION' '=' ('TRUE' | 'FALSE') | 'BINARY_AS_TEXT' '=' ('TRUE' | 'FALSE') | 'USE_LOGICAL_TYPE' '=' ('TRUE' | 'FALSE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'USE_VECTORIZED_SCANNER' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')'))*
    ;
alterStage_formatTypeOptions_XML
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'IGNORE_UTF8_ERRORS' '=' ('TRUE' | 'FALSE') | 'PRESERVE_SPACE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ELEMENT' '=' ('TRUE' | 'FALSE') | 'DISABLE_AUTO_CONVERT' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE')) (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'IGNORE_UTF8_ERRORS' '=' ('TRUE' | 'FALSE') | 'PRESERVE_SPACE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ELEMENT' '=' ('TRUE' | 'FALSE') | 'DISABLE_AUTO_CONVERT' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE')))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-storage-integration
alterStorageIntegration
    : 'ALTER' 'STORAGE'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'SET' (alterStorageIntegration_cloudProviderParams ('ENABLED' '=' ('TRUE' | 'FALSE'))? ('ACTIVE' '=' STRING)? ('STORAGE_ALLOWED_LOCATIONS' '=' '(' STRING (',' STRING)* ')')? ('STORAGE_BLOCKED_LOCATIONS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'ENABLED' '=' ('TRUE' | 'FALSE') ('ACTIVE' '=' STRING)? ('STORAGE_ALLOWED_LOCATIONS' '=' '(' STRING (',' STRING)* ')')? ('STORAGE_BLOCKED_LOCATIONS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'ACTIVE' '=' STRING ('STORAGE_ALLOWED_LOCATIONS' '=' '(' STRING (',' STRING)* ')')? ('STORAGE_BLOCKED_LOCATIONS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'STORAGE_ALLOWED_LOCATIONS' '=' '(' STRING (',' STRING)* ')' ('STORAGE_BLOCKED_LOCATIONS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? | 'STORAGE_BLOCKED_LOCATIONS' '=' '(' STRING (',' STRING)* ')' ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'STORAGE'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'STORAGE'? 'INTEGRATION' ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'STORAGE'? 'INTEGRATION' ('IF' 'EXISTS')? ID 'UNSET' ('ENABLED' | 'STORAGE_BLOCKED_LOCATIONS' | 'COMMENT') (',' ('ENABLED' | 'STORAGE_BLOCKED_LOCATIONS' | 'COMMENT'))*
    ;
alterStorageIntegration_cloudProviderParams
    : 'STORAGE_AWS_ROLE_ARN' '=' STRING ('STORAGE_AWS_OBJECT_ACL' '=' '\'bucket-owner-full-control\'')? ('STORAGE_AWS_EXTERNAL_ID' '=' STRING)? ('USE_PRIVATELINK_ENDPOINT' '=' ('TRUE' | 'FALSE'))?
    | 'AZURE_TENANT_ID' '=' STRING ('USE_PRIVATELINK_ENDPOINT' '=' ('TRUE' | 'FALSE'))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-storage-lifecycle-policy
alterStorageLifecyclePolicy
    : 'ALTER' 'STORAGE' 'LIFECYCLE' 'POLICY' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'STORAGE' 'LIFECYCLE' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'BODY' '->' EXPR
    | 'ALTER' 'STORAGE' 'LIFECYCLE' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'ARCHIVE_TIER' '=' ('COOL' | 'COLD')
    | 'ALTER' 'STORAGE' 'LIFECYCLE' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'ARCHIVE_FOR_DAYS' '=' NUMBER
    | 'ALTER' 'STORAGE' 'LIFECYCLE' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'STORAGE' 'LIFECYCLE' 'POLICY' ('IF' 'EXISTS')? QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'STORAGE' 'LIFECYCLE' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' 'ARCHIVE_FOR_DAYS'
    | 'ALTER' 'STORAGE' 'LIFECYCLE' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' 'COMMENT'
    | 'ALTER' 'STORAGE' 'LIFECYCLE' 'POLICY' ('IF' 'EXISTS')? QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-stream
alterStream
    : 'ALTER' 'STREAM' ('IF' 'EXISTS')? QNAME 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'STREAM' ('IF' 'EXISTS')? QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'STREAM' QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'STREAM' ('IF' 'EXISTS')? QNAME 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-streamlit
alterStreamlit
    : 'ALTER' 'STREAMLIT' ('IF' 'EXISTS')? QNAME 'SET' ('MAIN_FILE' '=' STRING ('QUERY_WAREHOUSE' '=' ID)? ('RUNTIME_NAME' '=' STRING)? ('COMPUTE_POOL' '=' ID)? ('COMMENT' '=' STRING)? ('TITLE' '=' STRING)? ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' ID (',' STRING '=' ID)* ')')? | 'QUERY_WAREHOUSE' '=' ID ('RUNTIME_NAME' '=' STRING)? ('COMPUTE_POOL' '=' ID)? ('COMMENT' '=' STRING)? ('TITLE' '=' STRING)? ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' ID (',' STRING '=' ID)* ')')? | 'RUNTIME_NAME' '=' STRING ('COMPUTE_POOL' '=' ID)? ('COMMENT' '=' STRING)? ('TITLE' '=' STRING)? ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' ID (',' STRING '=' ID)* ')')? | 'COMPUTE_POOL' '=' ID ('COMMENT' '=' STRING)? ('TITLE' '=' STRING)? ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' ID (',' STRING '=' ID)* ')')? | 'COMMENT' '=' STRING ('TITLE' '=' STRING)? ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' ID (',' STRING '=' ID)* ')')? | 'TITLE' '=' STRING ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' ID (',' STRING '=' ID)* ')')? | 'IMPORTS' '=' '(' STRING (',' STRING)* ')' ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' ID (',' STRING '=' ID)* ')')? | 'EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')' ('SECRETS' '=' '(' STRING '=' ID (',' STRING '=' ID)* ')')? | 'SECRETS' '=' '(' STRING '=' ID (',' STRING '=' ID)* ')')
    | 'ALTER' 'STREAMLIT' ('IF' 'EXISTS')? QNAME 'UNSET' ('SECRETS' | 'EXTERNAL_ACCESS_INTEGRATIONS' | 'QUERY_WAREHOUSE' | 'TITLE' | 'COMMENT') (',' ('SECRETS' | 'EXTERNAL_ACCESS_INTEGRATIONS' | 'QUERY_WAREHOUSE' | 'TITLE' | 'COMMENT'))*
    | 'ALTER' 'STREAMLIT' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'STREAMLIT' QNAME 'COMMIT'
    | 'ALTER' 'STREAMLIT' QNAME 'PUSH' ('TO' ID)? (('GIT_CREDENTIALS' '=' ID | 'USERNAME' '=' ID 'PASSWORD' '=' STRING) 'NAME' '=' ID 'EMAIL' '=' ID)? ('COMMENT' '=' ID)?
    | 'ALTER' 'STREAMLIT' QNAME 'ABORT'
    | 'ALTER' 'STREAMLIT' QNAME 'PULL'
    | 'ALTER' 'STREAMLIT' QNAME 'ADD' 'LIVE' 'VERSION' 'FROM' 'LAST'
    | 'ALTER' 'STREAMLIT' ('IF' 'EXISTS')? QNAME 'SET' ('ROOT_LOCATION' '=' STRING ('MAIN_FILE' '=' STRING)? ('QUERY_WAREHOUSE' '=' ID)? ('COMMENT' '=' STRING)? ('TITLE' '=' STRING)? ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? | 'MAIN_FILE' '=' STRING ('QUERY_WAREHOUSE' '=' ID)? ('COMMENT' '=' STRING)? ('TITLE' '=' STRING)? ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? | 'QUERY_WAREHOUSE' '=' ID ('COMMENT' '=' STRING)? ('TITLE' '=' STRING)? ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? | 'COMMENT' '=' STRING ('TITLE' '=' STRING)? ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? | 'TITLE' '=' STRING ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? | 'IMPORTS' '=' '(' STRING (',' STRING)* ')' ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? | 'EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-table
// fix: constraintAction: a } after ALTER | MODIFY's constraint list closes no bracket; dropped
// fix: tableColumnAction: ALTER and MODIFY both take the column actions after them
// fix: the column actions are a comma-separated list of actions, each one of those the page lists, as its examples show
// fix: list_of_properties is a space-delimited list of the properties the parameters name
// fix: a search method with its targets, as the parameters write it
// fix: read: <constraint_properties> is the properties CREATE TABLE's constraint page lists
// fix: read: <constraint_properties> is the properties CREATE TABLE's constraint page lists
// fix: read: an out-of-line key constraint names its columns, as the constraint page says
// fix: read: a foreign key's REFERENCES takes MATCH and ON UPDATE/DELETE, the properties the constraint page gives foreign keys
alterTable
    : 'ALTER' 'TABLE' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'TABLE' ('IF' 'EXISTS')? QNAME 'SWAP' 'WITH' QNAME
    | 'ALTER' 'TABLE' ('IF' 'EXISTS')? QNAME (alterTable_clusteringAction | alterTable_tableColumnAction | alterTable_constraintAction)
    | 'ALTER' 'TABLE' ('IF' 'EXISTS')? QNAME alterTable_dataMetricFunctionAction
    | 'ALTER' 'TABLE' ('IF' 'EXISTS')? QNAME alterTable_dataGovnPolicyTagAction
    | 'ALTER' 'TABLE' ('IF' 'EXISTS')? QNAME alterTable_extTableColumnAction
    | 'ALTER' 'TABLE' ('IF' 'EXISTS')? QNAME alterTable_searchOptimizationAction
    | 'ALTER' 'TABLE' ('IF' 'EXISTS')? QNAME 'ADD' 'STORAGE' 'LIFECYCLE' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')'
    | 'ALTER' 'TABLE' ('IF' 'EXISTS')? QNAME 'DROP' 'STORAGE' 'LIFECYCLE' 'POLICY'
    | 'ALTER' 'TABLE' ('IF' 'EXISTS')? QNAME 'SET' ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('CHANGE_TRACKING' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('COMMENT' '=' STRING)? ('ROW_TIMESTAMP' '=' ('TRUE' | 'FALSE'))? | 'MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER ('CHANGE_TRACKING' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('COMMENT' '=' STRING)? ('ROW_TIMESTAMP' '=' ('TRUE' | 'FALSE'))? | 'CHANGE_TRACKING' '=' ('TRUE' | 'FALSE') ('DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('COMMENT' '=' STRING)? ('ROW_TIMESTAMP' '=' ('TRUE' | 'FALSE'))? | 'DEFAULT_DDL_COLLATION' '=' STRING ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('COMMENT' '=' STRING)? ('ROW_TIMESTAMP' '=' ('TRUE' | 'FALSE'))? | 'ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('COMMENT' '=' STRING)? ('ROW_TIMESTAMP' '=' ('TRUE' | 'FALSE'))? | 'ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE') ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('COMMENT' '=' STRING)? ('ROW_TIMESTAMP' '=' ('TRUE' | 'FALSE'))? | 'ERROR_LOGGING' '=' ('TRUE' | 'FALSE') ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('COMMENT' '=' STRING)? ('ROW_TIMESTAMP' '=' ('TRUE' | 'FALSE'))? | 'CONTACT' ID '=' ID (',' ID '=' ID)* ('COMMENT' '=' STRING)? ('ROW_TIMESTAMP' '=' ('TRUE' | 'FALSE'))? | 'COMMENT' '=' STRING ('ROW_TIMESTAMP' '=' ('TRUE' | 'FALSE'))? | 'ROW_TIMESTAMP' '=' ('TRUE' | 'FALSE'))
    | 'ALTER' 'TABLE' ('IF' 'EXISTS')? QNAME 'UNSET' ('DATA_RETENTION_TIME_IN_DAYS' | 'MAX_DATA_EXTENSION_TIME_IN_DAYS' | 'CHANGE_TRACKING' | 'DEFAULT_DDL_COLLATION' | 'ICEBERG_DEFAULT_DDL_COLLATION' | 'ENABLE_SCHEMA_EVOLUTION' | 'ERROR_LOGGING' | 'CONTACT' ID | 'COMMENT' | 'ROW_TIMESTAMP') (',' ('DATA_RETENTION_TIME_IN_DAYS' | 'MAX_DATA_EXTENSION_TIME_IN_DAYS' | 'CHANGE_TRACKING' | 'DEFAULT_DDL_COLLATION' | 'ICEBERG_DEFAULT_DDL_COLLATION' | 'ENABLE_SCHEMA_EVOLUTION' | 'ERROR_LOGGING' | 'CONTACT' ID | 'COMMENT' | 'ROW_TIMESTAMP'))*
    ;
alterTable_clusteringAction
    : ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')' | 'RECLUSTER' ('MAX_SIZE' '=' NUMBER)? ('WHERE' EXPR)? | ('SUSPEND' | 'RESUME') 'RECLUSTER' | 'DROP' 'CLUSTERING' 'KEY')
    ;
alterTable_tableColumnAction
    : ('ADD' 'COLUMN'? ('IF' 'NOT' 'EXISTS')? ID DATA_TYPE (('GENERATED' 'ALWAYS')? 'AS' '(' EXPR ')' 'VIRTUAL'?)? (('DEFAULT' EXPR | ('AUTOINCREMENT' | 'IDENTITY') (('(' NUMBER ',' NUMBER ')' | 'START' NUMBER 'INCREMENT' NUMBER))? (('ORDER' | 'NOORDER'))?))? alterTable_inlineConstraint? ('COLLATE' STRING)? | 'RENAME' 'COLUMN' ID 'TO' ID | ('ALTER' | 'MODIFY') ('(' 'COLUMN'? ID ('DROP' 'DEFAULT' | 'SET' 'DEFAULT' QNAME '.' 'NEXTVAL' | ('SET'? 'NOT' 'NULL' | 'DROP' 'NOT' 'NULL') | (('SET' 'DATA')? 'TYPE')? DATA_TYPE | 'COMMENT' STRING | 'UNSET' 'COMMENT') (',' 'COLUMN'? ID ('DROP' 'DEFAULT' | 'SET' 'DEFAULT' QNAME '.' 'NEXTVAL' | ('SET'? 'NOT' 'NULL' | 'DROP' 'NOT' 'NULL') | (('SET' 'DATA')? 'TYPE')? DATA_TYPE | 'COMMENT' STRING | 'UNSET' 'COMMENT'))* ')' | 'COLUMN'? ID ('DROP' 'DEFAULT' | 'SET' 'DEFAULT' QNAME '.' 'NEXTVAL' | ('SET'? 'NOT' 'NULL' | 'DROP' 'NOT' 'NULL') | (('SET' 'DATA')? 'TYPE')? DATA_TYPE | 'COMMENT' STRING | 'UNSET' 'COMMENT') (',' 'COLUMN'? ID ('DROP' 'DEFAULT' | 'SET' 'DEFAULT' QNAME '.' 'NEXTVAL' | ('SET'? 'NOT' 'NULL' | 'DROP' 'NOT' 'NULL') | (('SET' 'DATA')? 'TYPE')? DATA_TYPE | 'COMMENT' STRING | 'UNSET' 'COMMENT'))*) | 'DROP' 'COLUMN'? ('IF' 'EXISTS')? ID (',' ID)*)
    ;
alterTable_inlineConstraint
    : ('NOT' 'NULL')? ('CONSTRAINT' ID)? ('UNIQUE' | 'PRIMARY' 'KEY' | ('FOREIGN' 'KEY')? 'REFERENCES' QNAME ('(' ID ')')? ('MATCH' ('FULL' | 'SIMPLE' | 'PARTIAL'))? ('ON' ('UPDATE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))? ('DELETE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))?)? | 'CHECK' '(' EXPR ')') ('NOT'? 'ENFORCED')? ('NOT'? 'DEFERRABLE')? ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? (('RELY' | 'NORELY'))?
    ;
alterTable_dataMetricFunctionAction
    : 'SET' 'DATA_METRIC_SCHEDULE' '=' (STRING | '\'TRIGGER_ON_CHANGES\'')
    | 'UNSET' 'DATA_METRIC_SCHEDULE'
    | 'ADD' 'DATA' 'METRIC' 'FUNCTION' ID 'ON' '(' ID (',' ID)* (',' 'TABLE' QNAME '(' ID (',' ID)* ')')? ')' ('WITHIN' 'GROUP' '(' ID (',' ID)* ')' ('GROUP' 'LIMIT' NUMBER)?)? ('FILTER' '(' EXPR ')')? ('EXPECTATION' ID '(' EXPR ')' (',' ID '(' EXPR ')')*)? ('EXECUTE' 'AS' 'ROLE' ID)? ('ANOMALY_DETECTION' '=' ('TRUE' | 'FALSE'))? ('SENSITIVITY' '=' ('\'LOW\'' | '\'MEDIUM\'' | '\'HIGH\''))? (',' ID 'ON' '(' ID (',' ID)* (',' 'TABLE' QNAME '(' ID (',' ID)* ')')? ')')? ('WITHIN' 'GROUP' '(' ID (',' ID)* ')' ('GROUP' 'LIMIT' NUMBER)?)? ('FILTER' '(' EXPR ')')? ('EXPECTATION' ID '(' EXPR ')' (',' ID '(' EXPR ')')*)? ('EXECUTE' 'AS' 'ROLE' ID)? ('ANOMALY_DETECTION' '=' ('TRUE' | 'FALSE'))? ('SENSITIVITY' '=' ('\'LOW\'' | '\'MEDIUM\'' | '\'HIGH\''))?
    | 'DROP' 'DATA' 'METRIC' 'FUNCTION' ID 'ON' '(' ID (',' ID)* (',' 'TABLE' QNAME '(' ID (',' ID)* ')')? ')' (',' ID 'ON' '(' ID (',' ID)* (',' 'TABLE' QNAME '(' ID (',' ID)* ')')? ')')?
    | 'MODIFY' 'DATA' 'METRIC' 'FUNCTION' ID 'ON' '(' ID (',' ID)* (',' 'TABLE' QNAME '(' ID (',' ID)* ')')? ')' ('SUSPEND' | 'RESUME') (',' ID 'ON' '(' ID (',' ID)* (',' 'TABLE' QNAME '(' ID (',' ID)* ')')? ')' ('SUSPEND' | 'RESUME'))?
    | 'MODIFY' 'DATA' 'METRIC' 'FUNCTION' ID 'ON' '(' ID (',' ID)* (',' 'TABLE' QNAME '(' ID (',' ID)* ')')? ')' ('ADD' | 'MODIFY') 'EXPECTATION' ID '(' EXPR ')' (',' ID '(' EXPR ')')*
    | 'MODIFY' 'DATA' 'METRIC' 'FUNCTION' ID 'ON' '(' ID (',' ID)* (',' 'TABLE' QNAME '(' ID (',' ID)* ')')? ')' 'DROP' 'EXPECTATION' ID (',' ID)*
    | 'MODIFY' 'DATA' 'METRIC' 'FUNCTION' ID 'ON' '(' ID (',' ID)* (',' 'TABLE' QNAME '(' ID (',' ID)* ')')? ')' 'FILTER' '(' EXPR? ')' (',' ID 'ON' '(' ID (',' ID)* (',' 'TABLE' QNAME '(' ID (',' ID)* ')')? ')' 'FILTER' '(' EXPR? ')')?
    | 'MODIFY' 'DATA' 'METRIC' 'FUNCTION' ID 'SET' ('ANOMALY_DETECTION' '=' ('TRUE' | 'FALSE') | 'SENSITIVITY' '=' ('\'LOW\'' | '\'MEDIUM\'' | '\'HIGH\'') | 'DATA_QUALITY_NOTIFICATION' '=' ('TRUE' | 'FALSE'))+
    ;
alterTable_dataGovnPolicyTagAction
    : ('SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)* | 'UNSET' 'TAG' QNAME (',' QNAME)*)
    | ('ADD' 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')' | 'DROP' 'ROW' 'ACCESS' 'POLICY' QNAME | 'DROP' 'ROW' 'ACCESS' 'POLICY' QNAME ',' 'ADD' 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')' | 'DROP' 'ALL' 'ROW' 'ACCESS' 'POLICIES')
    | ('SET' 'AGGREGATION' 'POLICY' QNAME ('ENTITY' 'KEY' '(' ID (',' ID)* ')')? 'FORCE'? | 'UNSET' 'AGGREGATION' 'POLICY')
    | ('SET' 'JOIN' 'POLICY' QNAME 'FORCE'? | 'UNSET' 'JOIN' 'POLICY')
    | 'ADD' 'COLUMN'? ('IF' 'NOT' 'EXISTS')? ID DATA_TYPE ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | (('ALTER' | 'MODIFY') 'COLUMN'? ID 'SET' 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')? 'FORCE'? | ('ALTER' | 'MODIFY') 'COLUMN'? ID 'UNSET' 'MASKING' 'POLICY')
    | (('ALTER' | 'MODIFY') 'COLUMN'? ID 'SET' 'PROJECTION' 'POLICY' QNAME 'FORCE'? | ('ALTER' | 'MODIFY') 'COLUMN'? ID 'UNSET' 'PROJECTION' 'POLICY')
    | ('ALTER' | 'MODIFY') 'COLUMN'? ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)* ',' 'COLUMN'? ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | ('ALTER' | 'MODIFY') 'COLUMN'? ID 'UNSET' 'TAG' QNAME (',' QNAME)* ',' 'COLUMN'? ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;
alterTable_extTableColumnAction
    : ('ADD' 'COLUMN'? ('IF' 'NOT' 'EXISTS')? ID DATA_TYPE 'AS' '(' EXPR ')' | 'RENAME' 'COLUMN' ID 'TO' ID | 'DROP' 'COLUMN'? ('IF' 'EXISTS')? ID (',' ID)*)
    ;
alterTable_constraintAction
    : ('ADD' alterTable_outoflineConstraint | 'RENAME' 'CONSTRAINT' ID 'TO' ID | ('ALTER' | 'MODIFY') ('CONSTRAINT' ID | 'PRIMARY' 'KEY' | 'UNIQUE' | 'FOREIGN' 'KEY') '(' ID (',' ID)* ')' ('NOT'? 'ENFORCED')? ('VALIDATE' | 'NOVALIDATE')? ('RELY' | 'NORELY')? | 'DROP' ('CONSTRAINT' ID | 'PRIMARY' 'KEY' | 'UNIQUE' | 'FOREIGN' 'KEY') '(' ID (',' ID)* ')' ('CASCADE' | 'RESTRICT')?)
    ;
alterTable_outoflineConstraint
    : ('CONSTRAINT' ID)? ('UNIQUE' '(' ID (',' ID)* ')' | 'PRIMARY' 'KEY' '(' ID (',' ID)* ')' | ('FOREIGN' 'KEY')? '(' ID (',' ID)* ')' 'REFERENCES' QNAME ('(' ID (',' ID)* ')')? ('MATCH' ('FULL' | 'SIMPLE' | 'PARTIAL'))? ('ON' ('UPDATE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))? ('DELETE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))?)? | 'CHECK' '(' EXPR ')') ('NOT'? 'ENFORCED')? ('NOT'? 'DEFERRABLE')? ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? (('RELY' | 'NORELY'))?
    ;
alterTable_searchOptimizationAction
    : ('ADD' 'SEARCH' 'OPTIMIZATION' ('ON' ('FULL_TEXT' | 'EQUALITY' | 'SUBSTRING' | 'GEO') '(' ID (',' ID)* (',' 'ANALYZER' '=>' STRING)? ')' (',' ('FULL_TEXT' | 'EQUALITY' | 'SUBSTRING' | 'GEO') '(' ID (',' ID)* (',' 'ANALYZER' '=>' STRING)? ')')*)? | 'DROP' 'SEARCH' 'OPTIMIZATION' ('ON' (('FULL_TEXT' | 'EQUALITY' | 'SUBSTRING' | 'GEO') '(' ID (',' ID)* (',' 'ANALYZER' '=>' STRING)? ')' | ID | EXPR) (',' (('FULL_TEXT' | 'EQUALITY' | 'SUBSTRING' | 'GEO') '(' ID (',' ID)* (',' 'ANALYZER' '=>' STRING)? ')' | ID | EXPR))*)?)
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-table-column
// fix: the column actions are a comma-separated list of actions, each one of those the page lists, as its examples show
// not converted (form): undefined fragment "dataGovnPolicyTagAction": the page does not define it, or it did not convert
//   ALTER TABLE <name> { ALTER | MODIFY } [ COLUMN ] dataGovnPolicyTagAction
alterTableColumn
    : 'ALTER' 'TABLE' QNAME ('ALTER' | 'MODIFY') ('(' 'COLUMN'? ID ('DROP' 'DEFAULT' | 'SET' 'DEFAULT' QNAME '.' 'NEXTVAL' | ('SET'? 'NOT' 'NULL' | 'DROP' 'NOT' 'NULL') | (('SET' 'DATA')? 'TYPE')? DATA_TYPE | 'COMMENT' STRING | 'UNSET' 'COMMENT') (',' 'COLUMN'? ID ('DROP' 'DEFAULT' | 'SET' 'DEFAULT' QNAME '.' 'NEXTVAL' | ('SET'? 'NOT' 'NULL' | 'DROP' 'NOT' 'NULL') | (('SET' 'DATA')? 'TYPE')? DATA_TYPE | 'COMMENT' STRING | 'UNSET' 'COMMENT'))* ')' | 'COLUMN'? ID ('DROP' 'DEFAULT' | 'SET' 'DEFAULT' QNAME '.' 'NEXTVAL' | ('SET'? 'NOT' 'NULL' | 'DROP' 'NOT' 'NULL') | (('SET' 'DATA')? 'TYPE')? DATA_TYPE | 'COMMENT' STRING | 'UNSET' 'COMMENT') (',' 'COLUMN'? ID ('DROP' 'DEFAULT' | 'SET' 'DEFAULT' QNAME '.' 'NEXTVAL' | ('SET'? 'NOT' 'NULL' | 'DROP' 'NOT' 'NULL') | (('SET' 'DATA')? 'TYPE')? DATA_TYPE | 'COMMENT' STRING | 'UNSET' 'COMMENT'))*)
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-table-event-table
// fix: a search method with its targets, as the parameters write it
alterTableEventTable
    : 'ALTER' 'TABLE' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'TABLE' ('IF' 'EXISTS')? QNAME alterTableEventTable_clusteringAction
    | 'ALTER' 'TABLE' ('IF' 'EXISTS')? QNAME alterTableEventTable_dataGovnPolicyTagAction
    | 'ALTER' 'TABLE' ('IF' 'EXISTS')? QNAME alterTableEventTable_searchOptimizationAction
    | 'ALTER' 'TABLE' ('IF' 'EXISTS')? QNAME 'SET' ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('CHANGE_TRACKING' '=' ('TRUE' | 'FALSE'))? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('COMMENT' '=' STRING)? | 'MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER ('CHANGE_TRACKING' '=' ('TRUE' | 'FALSE'))? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('COMMENT' '=' STRING)? | 'CHANGE_TRACKING' '=' ('TRUE' | 'FALSE') ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('COMMENT' '=' STRING)? | 'CONTACT' ID '=' ID (',' ID '=' ID)* ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'TABLE' ('IF' 'EXISTS')? QNAME 'UNSET' ('DATA_RETENTION_TIME_IN_DAYS' | 'MAX_DATA_EXTENSION_TIME_IN_DAYS' | 'CHANGE_TRACKING' | 'CONTACT' ID | 'COMMENT')
    ;
alterTableEventTable_clusteringAction
    : ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')' | ('SUSPEND' | 'RESUME') 'RECLUSTER' | 'DROP' 'CLUSTERING' 'KEY')
    ;
alterTableEventTable_dataGovnPolicyTagAction
    : ('SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)* | 'UNSET' 'TAG' QNAME (',' QNAME)*)
    | ('ADD' 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')' | 'DROP' 'ROW' 'ACCESS' 'POLICY' QNAME | 'DROP' 'ROW' 'ACCESS' 'POLICY' QNAME ',' 'ADD' 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')' | 'DROP' 'ALL' 'ROW' 'ACCESS' 'POLICIES')
    ;
alterTableEventTable_searchOptimizationAction
    : ('ADD' 'SEARCH' 'OPTIMIZATION' ('ON' ('FULL_TEXT' | 'EQUALITY' | 'SUBSTRING' | 'GEO') '(' ID (',' ID)* (',' 'ANALYZER' '=>' STRING)? ')' (',' ('FULL_TEXT' | 'EQUALITY' | 'SUBSTRING' | 'GEO') '(' ID (',' ID)* (',' 'ANALYZER' '=>' STRING)? ')')*)? | 'DROP' 'SEARCH' 'OPTIMIZATION' ('ON' (('FULL_TEXT' | 'EQUALITY' | 'SUBSTRING' | 'GEO') '(' ID (',' ID)* (',' 'ANALYZER' '=>' STRING)? ')' | ID | EXPR) (',' (('FULL_TEXT' | 'EQUALITY' | 'SUBSTRING' | 'GEO') '(' ID (',' ID)* (',' 'ANALYZER' '=>' STRING)? ')' | ID | EXPR))*)?)
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-tag
alterTag
    : 'ALTER' 'TAG' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'TAG' ('IF' 'EXISTS')? QNAME ('ADD' | 'DROP') 'ALLOWED_VALUES' STRING (',' STRING)*
    | 'ALTER' 'TAG' ('IF' 'EXISTS')? QNAME 'SET' ('ALLOWED_VALUES' STRING (',' STRING)* ('MULTI_VALUE' '=' 'TRUE')? ('PROPAGATE' '=' ('ON_DEPENDENCY_AND_DATA_MOVEMENT' | 'ON_DEPENDENCY' | 'ON_DATA_MOVEMENT') ('ON_CONFLICT' '=' (STRING | 'ALLOWED_VALUES_SEQUENCE' | 'MERGE'))?)? ('COMMENT' '=' STRING)? | 'MULTI_VALUE' '=' 'TRUE' ('PROPAGATE' '=' ('ON_DEPENDENCY_AND_DATA_MOVEMENT' | 'ON_DEPENDENCY' | 'ON_DATA_MOVEMENT') ('ON_CONFLICT' '=' (STRING | 'ALLOWED_VALUES_SEQUENCE' | 'MERGE'))?)? ('COMMENT' '=' STRING)? | 'PROPAGATE' '=' ('ON_DEPENDENCY_AND_DATA_MOVEMENT' | 'ON_DEPENDENCY' | 'ON_DATA_MOVEMENT') ('ON_CONFLICT' '=' (STRING | 'ALLOWED_VALUES_SEQUENCE' | 'MERGE'))? ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'TAG' ('IF' 'EXISTS')? QNAME 'UNSET' ('ALLOWED_VALUES' | 'PROPAGATE' | 'ON_CONFLICT' | 'COMMENT')
    | 'ALTER' 'TAG' ('IF' 'EXISTS')? QNAME 'SET' 'DATA' 'MOVEMENT' 'POLICY' QNAME 'FORCE'?
    | 'ALTER' 'TAG' ('IF' 'EXISTS')? QNAME 'UNSET' 'DATA' 'MOVEMENT' 'POLICY'
    | 'ALTER' 'TAG' ('IF' 'EXISTS')? QNAME 'SET' 'MASKING' 'POLICY' QNAME (',' 'MASKING' 'POLICY' QNAME)* 'FORCE'?
    | 'ALTER' 'TAG' ('IF' 'EXISTS')? QNAME 'UNSET' 'MASKING' 'POLICY' QNAME (',' 'MASKING' 'POLICY' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-task
// fix: SERVERLESS_TASK_MIN/MAX_STATEMENT_SIZE: the choice of sizes is written inside the quotes; read as a string
// fix: AFTER names predecessor tasks, as the parameters and examples write them
// fix: the warehouse is a name, as the parameters say
// fix: a session parameter takes a constant value, as the parameters say
alterTask
    : 'ALTER' 'TASK' ('IF' 'EXISTS')? QNAME 'RESUME'
    | 'ALTER' 'TASK' ('IF' 'EXISTS')? QNAME 'SUSPEND'
    | 'ALTER' 'TASK' ('IF' 'EXISTS')? QNAME 'REMOVE' 'AFTER' QNAME (',' QNAME)*
    | 'ALTER' 'TASK' ('IF' 'EXISTS')? QNAME 'ADD' 'AFTER' QNAME (',' QNAME)*
    | 'ALTER' 'TASK' ('IF' 'EXISTS')? QNAME 'SET' (('WAREHOUSE' '=' ID | 'USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE' '=' STRING) ('SCHEDULE' '=' STRING)? ('CONFIG' '=' STRING)? ('OVERLAP_POLICY' '=' ('NO_OVERLAP' | 'ALLOW_CHILD_OVERLAP' | 'ALLOW_ALL_OVERLAP'))? ('USER_TASK_TIMEOUT_MS' '=' NUMBER)? ('SUSPEND_TASK_AFTER_NUM_FAILURES' '=' NUMBER)? ('ERROR_INTEGRATION' '=' ID)? ('SUCCESS_INTEGRATION' '=' ID)? ('LOG_LEVEL' '=' STRING)? ('COMMENT' '=' STRING)? (ID '=' EXPR (',' ID '=' EXPR)*)? ('TASK_AUTO_RETRY_ATTEMPTS' '=' NUMBER)? ('USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS' '=' NUMBER)? ('TARGET_COMPLETION_INTERVAL' '=' STRING)? ('SERVERLESS_TASK_MIN_STATEMENT_SIZE' '=' STRING)? ('SERVERLESS_TASK_MAX_STATEMENT_SIZE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('EXECUTE' 'AS' 'USER' ID)? | 'SCHEDULE' '=' STRING ('CONFIG' '=' STRING)? ('OVERLAP_POLICY' '=' ('NO_OVERLAP' | 'ALLOW_CHILD_OVERLAP' | 'ALLOW_ALL_OVERLAP'))? ('USER_TASK_TIMEOUT_MS' '=' NUMBER)? ('SUSPEND_TASK_AFTER_NUM_FAILURES' '=' NUMBER)? ('ERROR_INTEGRATION' '=' ID)? ('SUCCESS_INTEGRATION' '=' ID)? ('LOG_LEVEL' '=' STRING)? ('COMMENT' '=' STRING)? (ID '=' EXPR (',' ID '=' EXPR)*)? ('TASK_AUTO_RETRY_ATTEMPTS' '=' NUMBER)? ('USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS' '=' NUMBER)? ('TARGET_COMPLETION_INTERVAL' '=' STRING)? ('SERVERLESS_TASK_MIN_STATEMENT_SIZE' '=' STRING)? ('SERVERLESS_TASK_MAX_STATEMENT_SIZE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('EXECUTE' 'AS' 'USER' ID)? | 'CONFIG' '=' STRING ('OVERLAP_POLICY' '=' ('NO_OVERLAP' | 'ALLOW_CHILD_OVERLAP' | 'ALLOW_ALL_OVERLAP'))? ('USER_TASK_TIMEOUT_MS' '=' NUMBER)? ('SUSPEND_TASK_AFTER_NUM_FAILURES' '=' NUMBER)? ('ERROR_INTEGRATION' '=' ID)? ('SUCCESS_INTEGRATION' '=' ID)? ('LOG_LEVEL' '=' STRING)? ('COMMENT' '=' STRING)? (ID '=' EXPR (',' ID '=' EXPR)*)? ('TASK_AUTO_RETRY_ATTEMPTS' '=' NUMBER)? ('USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS' '=' NUMBER)? ('TARGET_COMPLETION_INTERVAL' '=' STRING)? ('SERVERLESS_TASK_MIN_STATEMENT_SIZE' '=' STRING)? ('SERVERLESS_TASK_MAX_STATEMENT_SIZE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('EXECUTE' 'AS' 'USER' ID)? | 'OVERLAP_POLICY' '=' ('NO_OVERLAP' | 'ALLOW_CHILD_OVERLAP' | 'ALLOW_ALL_OVERLAP') ('USER_TASK_TIMEOUT_MS' '=' NUMBER)? ('SUSPEND_TASK_AFTER_NUM_FAILURES' '=' NUMBER)? ('ERROR_INTEGRATION' '=' ID)? ('SUCCESS_INTEGRATION' '=' ID)? ('LOG_LEVEL' '=' STRING)? ('COMMENT' '=' STRING)? (ID '=' EXPR (',' ID '=' EXPR)*)? ('TASK_AUTO_RETRY_ATTEMPTS' '=' NUMBER)? ('USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS' '=' NUMBER)? ('TARGET_COMPLETION_INTERVAL' '=' STRING)? ('SERVERLESS_TASK_MIN_STATEMENT_SIZE' '=' STRING)? ('SERVERLESS_TASK_MAX_STATEMENT_SIZE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('EXECUTE' 'AS' 'USER' ID)? | 'USER_TASK_TIMEOUT_MS' '=' NUMBER ('SUSPEND_TASK_AFTER_NUM_FAILURES' '=' NUMBER)? ('ERROR_INTEGRATION' '=' ID)? ('SUCCESS_INTEGRATION' '=' ID)? ('LOG_LEVEL' '=' STRING)? ('COMMENT' '=' STRING)? (ID '=' EXPR (',' ID '=' EXPR)*)? ('TASK_AUTO_RETRY_ATTEMPTS' '=' NUMBER)? ('USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS' '=' NUMBER)? ('TARGET_COMPLETION_INTERVAL' '=' STRING)? ('SERVERLESS_TASK_MIN_STATEMENT_SIZE' '=' STRING)? ('SERVERLESS_TASK_MAX_STATEMENT_SIZE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('EXECUTE' 'AS' 'USER' ID)? | 'SUSPEND_TASK_AFTER_NUM_FAILURES' '=' NUMBER ('ERROR_INTEGRATION' '=' ID)? ('SUCCESS_INTEGRATION' '=' ID)? ('LOG_LEVEL' '=' STRING)? ('COMMENT' '=' STRING)? (ID '=' EXPR (',' ID '=' EXPR)*)? ('TASK_AUTO_RETRY_ATTEMPTS' '=' NUMBER)? ('USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS' '=' NUMBER)? ('TARGET_COMPLETION_INTERVAL' '=' STRING)? ('SERVERLESS_TASK_MIN_STATEMENT_SIZE' '=' STRING)? ('SERVERLESS_TASK_MAX_STATEMENT_SIZE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('EXECUTE' 'AS' 'USER' ID)? | 'ERROR_INTEGRATION' '=' ID ('SUCCESS_INTEGRATION' '=' ID)? ('LOG_LEVEL' '=' STRING)? ('COMMENT' '=' STRING)? (ID '=' EXPR (',' ID '=' EXPR)*)? ('TASK_AUTO_RETRY_ATTEMPTS' '=' NUMBER)? ('USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS' '=' NUMBER)? ('TARGET_COMPLETION_INTERVAL' '=' STRING)? ('SERVERLESS_TASK_MIN_STATEMENT_SIZE' '=' STRING)? ('SERVERLESS_TASK_MAX_STATEMENT_SIZE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('EXECUTE' 'AS' 'USER' ID)? | 'SUCCESS_INTEGRATION' '=' ID ('LOG_LEVEL' '=' STRING)? ('COMMENT' '=' STRING)? (ID '=' EXPR (',' ID '=' EXPR)*)? ('TASK_AUTO_RETRY_ATTEMPTS' '=' NUMBER)? ('USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS' '=' NUMBER)? ('TARGET_COMPLETION_INTERVAL' '=' STRING)? ('SERVERLESS_TASK_MIN_STATEMENT_SIZE' '=' STRING)? ('SERVERLESS_TASK_MAX_STATEMENT_SIZE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('EXECUTE' 'AS' 'USER' ID)? | 'LOG_LEVEL' '=' STRING ('COMMENT' '=' STRING)? (ID '=' EXPR (',' ID '=' EXPR)*)? ('TASK_AUTO_RETRY_ATTEMPTS' '=' NUMBER)? ('USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS' '=' NUMBER)? ('TARGET_COMPLETION_INTERVAL' '=' STRING)? ('SERVERLESS_TASK_MIN_STATEMENT_SIZE' '=' STRING)? ('SERVERLESS_TASK_MAX_STATEMENT_SIZE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('EXECUTE' 'AS' 'USER' ID)? | 'COMMENT' '=' STRING (ID '=' EXPR (',' ID '=' EXPR)*)? ('TASK_AUTO_RETRY_ATTEMPTS' '=' NUMBER)? ('USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS' '=' NUMBER)? ('TARGET_COMPLETION_INTERVAL' '=' STRING)? ('SERVERLESS_TASK_MIN_STATEMENT_SIZE' '=' STRING)? ('SERVERLESS_TASK_MAX_STATEMENT_SIZE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('EXECUTE' 'AS' 'USER' ID)? | ID '=' EXPR (',' ID '=' EXPR)* ('TASK_AUTO_RETRY_ATTEMPTS' '=' NUMBER)? ('USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS' '=' NUMBER)? ('TARGET_COMPLETION_INTERVAL' '=' STRING)? ('SERVERLESS_TASK_MIN_STATEMENT_SIZE' '=' STRING)? ('SERVERLESS_TASK_MAX_STATEMENT_SIZE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('EXECUTE' 'AS' 'USER' ID)? | 'TASK_AUTO_RETRY_ATTEMPTS' '=' NUMBER ('USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS' '=' NUMBER)? ('TARGET_COMPLETION_INTERVAL' '=' STRING)? ('SERVERLESS_TASK_MIN_STATEMENT_SIZE' '=' STRING)? ('SERVERLESS_TASK_MAX_STATEMENT_SIZE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('EXECUTE' 'AS' 'USER' ID)? | 'USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS' '=' NUMBER ('TARGET_COMPLETION_INTERVAL' '=' STRING)? ('SERVERLESS_TASK_MIN_STATEMENT_SIZE' '=' STRING)? ('SERVERLESS_TASK_MAX_STATEMENT_SIZE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('EXECUTE' 'AS' 'USER' ID)? | 'TARGET_COMPLETION_INTERVAL' '=' STRING ('SERVERLESS_TASK_MIN_STATEMENT_SIZE' '=' STRING)? ('SERVERLESS_TASK_MAX_STATEMENT_SIZE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('EXECUTE' 'AS' 'USER' ID)? | 'SERVERLESS_TASK_MIN_STATEMENT_SIZE' '=' STRING ('SERVERLESS_TASK_MAX_STATEMENT_SIZE' '=' STRING)? ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('EXECUTE' 'AS' 'USER' ID)? | 'SERVERLESS_TASK_MAX_STATEMENT_SIZE' '=' STRING ('CONTACT' ID '=' ID (',' ID '=' ID)*)? ('EXECUTE' 'AS' 'USER' ID)? | 'CONTACT' ID '=' ID (',' ID '=' ID)* ('EXECUTE' 'AS' 'USER' ID)? | 'EXECUTE' 'AS' 'USER' ID)
    | 'ALTER' 'TASK' ('IF' 'EXISTS')? QNAME 'UNSET' ('WAREHOUSE' | 'SCHEDULE' | 'CONFIG' | 'OVERLAP_POLICY' | 'USER_TASK_TIMEOUT_MS' | 'SUSPEND_TASK_AFTER_NUM_FAILURES' | 'LOG_LEVEL' | 'COMMENT' | ID (',' ID)* | 'TARGET_COMPLETION_INTERVAL' | 'SERVERLESS_TASK_MIN_STATEMENT_SIZE' | 'SERVERLESS_TASK_MAX_STATEMENT_SIZE' | 'CONTACT' ID (',' ID)* | 'EXECUTE' 'AS' 'USER') (',' ('WAREHOUSE' | 'SCHEDULE' | 'CONFIG' | 'OVERLAP_POLICY' | 'USER_TASK_TIMEOUT_MS' | 'SUSPEND_TASK_AFTER_NUM_FAILURES' | 'LOG_LEVEL' | 'COMMENT' | ID (',' ID)* | 'TARGET_COMPLETION_INTERVAL' | 'SERVERLESS_TASK_MIN_STATEMENT_SIZE' | 'SERVERLESS_TASK_MAX_STATEMENT_SIZE' | 'CONTACT' ID (',' ID)* | 'EXECUTE' 'AS' 'USER'))*
    | 'ALTER' 'TASK' ('IF' 'EXISTS')? QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'TASK' ('IF' 'EXISTS')? QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'TASK' ('IF' 'EXISTS')? QNAME 'SET' 'FINALIZE' '=' STRING
    | 'ALTER' 'TASK' ('IF' 'EXISTS')? QNAME 'UNSET' 'FINALIZE'
    | 'ALTER' 'TASK' ('IF' 'EXISTS')? QNAME 'MODIFY' 'AS' QUERY
    | 'ALTER' 'TASK' ('IF' 'EXISTS')? QNAME 'MODIFY' 'WHEN' EXPR
    | 'ALTER' 'TASK' ('IF' 'EXISTS')? QNAME 'REMOVE' 'WHEN'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-task-graph
alterTaskGraph
    : 'ALTER' 'TASK' 'GRAPH' ('IF' 'EXISTS')? QNAME 'SET' ID '=' ID (',' ID '=' ID)*
    | 'ALTER' 'TASK' 'GRAPH' ('IF' 'EXISTS')? QNAME 'UNSET' ID (',' ID)*
    | 'ALTER' 'TASK' 'GRAPH' ('IF' 'EXISTS')? QNAME 'SET' 'TASK_DEFAULTS' '(' ID '=' ID (',' ID '=' ID)* ')'
    | 'ALTER' 'TASK' 'GRAPH' ('IF' 'EXISTS')? QNAME 'UNSET' 'TASK_DEFAULTS' '(' ID (',' ID)* ')'
    | 'ALTER' 'TASK' 'GRAPH' ('IF' 'EXISTS')? QNAME 'MODIFY' 'WHEN' EXPR
    | 'ALTER' 'TASK' 'GRAPH' ('IF' 'EXISTS')? QNAME 'REMOVE' 'WHEN'
    | 'ALTER' 'TASK' 'GRAPH' ('IF' 'EXISTS')? QNAME 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'TASK' 'GRAPH' ('IF' 'EXISTS')? QNAME 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'TASK' 'GRAPH' ('IF' 'EXISTS')? QNAME 'SET' 'CONTACT' ID '=' ID (',' ID '=' ID)*
    | 'ALTER' 'TASK' 'GRAPH' ('IF' 'EXISTS')? QNAME 'UNSET' 'CONTACT' ID (',' ID)*
    | 'ALTER' 'TASK' 'GRAPH' ('IF' 'EXISTS')? QNAME 'COMMIT'
    | 'ALTER' 'TASK' 'GRAPH' ('IF' 'EXISTS')? QNAME 'REVERT'
    | 'ALTER' 'TASK' 'GRAPH' ('IF' 'EXISTS')? QNAME 'RESUME'
    | 'ALTER' 'TASK' 'GRAPH' ('IF' 'EXISTS')? QNAME 'SUSPEND'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-type
alterType
    : 'ALTER' 'TYPE' ('IF' 'EXISTS')? QNAME 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'TYPE' ('IF' 'EXISTS')? QNAME 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-user
// fix: mfaActions: the five actions in braces have no | between them; each is one action
// fix: list_of_properties is the workload identity properties the parameters list
// fix: the interfaces are ALL or a list of the ones the parameters name
alterUser
    : 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'RENAME' 'TO' ID
    | 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'RESET' 'PASSWORD'
    | 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'ABORT' 'ALL' 'QUERIES'
    | 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'ADD' 'DELEGATED' 'AUTHORIZATION' 'OF' 'ROLE' ID 'TO' 'SECURITY' 'INTEGRATION' ID
    | 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'REMOVE' 'DELEGATED' ('AUTHORIZATION' 'OF' 'ROLE' ID | 'AUTHORIZATIONS') 'FROM' 'SECURITY' 'INTEGRATION' ID
    | 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'REMOVE' 'DELEGATED' 'AUTHORIZATIONS' 'OF' 'ANY' 'ROLE' 'FROM' 'SECURITY' 'INTEGRATION' ID
    | 'ALTER' 'USER' ('IF' 'EXISTS')? ID? alterUser_mfaActions
    | 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'SET' ('AUTHENTICATION' | 'PASSWORD' | 'SESSION') 'POLICY' QNAME 'FORCE'?
    | 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'UNSET' ('AUTHENTICATION' | 'PASSWORD' | 'SESSION') 'POLICY'
    | 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'SET' (alterUser_objectProperties alterUser_objectParams? alterUser_sessionParams? | alterUser_objectParams alterUser_sessionParams? | alterUser_sessionParams)
    | 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'UNSET' ID (',' ID)*
    ;
alterUser_mfaActions
    : ('ENROLL' 'MFA' | 'SET' 'DEFAULT_MFA_METHOD' '=' ('PASSKEY' | 'TOTP' | 'DUO') | 'REMOVE' 'MFA' 'METHOD' ID | 'MODIFY' 'MFA' 'METHOD' ID 'SET' 'COMMENT' '=' STRING | 'ADD' 'MFA' 'METHOD' 'OTP' ('COUNT' '=' NUMBER)?)
    ;
alterUser_objectProperties
    : ('PASSWORD' '=' STRING | 'LOGIN_NAME' '=' STRING | 'DISPLAY_NAME' '=' STRING | 'FIRST_NAME' '=' STRING | 'MIDDLE_NAME' '=' STRING | 'LAST_NAME' '=' STRING | 'EMAIL' '=' STRING | 'MUST_CHANGE_PASSWORD' '=' ('TRUE' | 'FALSE') | 'DISABLED' '=' ('TRUE' | 'FALSE') | 'ALLOWED_INTERFACES' '=' '(' ('\'ALL\'' | ('\'SNOWFLAKE_INTELLIGENCE\'' | '\'STREAMLIT\'') (',' ('\'SNOWFLAKE_INTELLIGENCE\'' | '\'STREAMLIT\''))*) ')' | 'DAYS_TO_EXPIRY' '=' NUMBER | 'MINS_TO_UNLOCK' '=' NUMBER | 'DEFAULT_WAREHOUSE' '=' STRING | 'DEFAULT_NAMESPACE' '=' STRING | 'DEFAULT_ROLE' '=' STRING | 'DEFAULT_SECONDARY_ROLES' '=' '(' '\'ALL\'' ')' | 'MINS_TO_BYPASS_MFA' '=' NUMBER | 'DISABLE_MFA' '=' ('TRUE' | 'FALSE') | 'RSA_PUBLIC_KEY' '=' STRING | 'RSA_PUBLIC_KEY_FP' '=' STRING | 'RSA_PUBLIC_KEY_2' '=' STRING | 'RSA_PUBLIC_KEY_2_FP' '=' STRING | 'TYPE' '=' ('PERSON' | 'SERVICE' | 'SERVICE_AGENT' | 'LEGACY_SERVICE') | 'WORKLOAD_IDENTITY' '=' '(' 'TYPE' '=' ('AWS' | 'AZURE' | 'GCP' | 'OIDC') ('ARN' '=' STRING)? ('ISSUER' '=' STRING)? ('SUBJECT' '=' STRING)? ('OIDC_AUDIENCE_LIST' '=' '(' STRING (',' STRING)* ')')? ')' | 'COMMENT' '=' STRING)+
    ;
alterUser_objectParams
    : ('ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR' '=' ('TRUE' | 'FALSE') | 'ENABLE_UNREDACTED_SECURE_OBJECT_ERROR' '=' ('TRUE' | 'FALSE') | 'NETWORK_POLICY' '=' STRING | 'PREVENT_UNLOAD_TO_INLINE_URL' '=' ('TRUE' | 'FALSE') | 'PREVENT_UNLOAD_TO_INTERNAL_STAGES' '=' ('TRUE' | 'FALSE'))+
    ;
alterUser_sessionParams
    : ('ABORT_DETACHED_QUERY' '=' ('TRUE' | 'FALSE') | 'AUTOCOMMIT' '=' ('TRUE' | 'FALSE') | 'BINARY_INPUT_FORMAT' '=' STRING | 'BINARY_OUTPUT_FORMAT' '=' STRING | 'DATE_INPUT_FORMAT' '=' STRING | 'DATE_OUTPUT_FORMAT' '=' STRING | 'DEFAULT_NULL_ORDERING' '=' STRING | 'ENABLE_GET_DDL_USE_DATA_TYPE_ALIAS' '=' ('TRUE' | 'FALSE') | 'ENABLE_NOTEBOOK_CREATION_IN_PERSONAL_DB' '=' ('TRUE' | 'FALSE') | 'ERROR_ON_NONDETERMINISTIC_MERGE' '=' ('TRUE' | 'FALSE') | 'ERROR_ON_NONDETERMINISTIC_UPDATE' '=' ('TRUE' | 'FALSE') | 'JSON_INDENT' '=' NUMBER | 'LOCK_TIMEOUT' '=' NUMBER | 'OPT_OUT_ERROR_LOGGING' '=' ('TRUE' | 'FALSE') | 'QUERY_TAG' '=' STRING | 'ROWS_PER_RESULTSET' '=' NUMBER | 'S3_STAGE_VPCE_DNS_NAME' '=' STRING | 'SEARCH_PATH' '=' STRING | 'SIMULATED_DATA_SHARING_CONSUMER' '=' STRING | 'STATEMENT_TIMEOUT_IN_SECONDS' '=' NUMBER | 'STRICT_JSON_OUTPUT' '=' ('TRUE' | 'FALSE') | 'TIMESTAMP_DAY_IS_ALWAYS_24H' '=' ('TRUE' | 'FALSE') | 'TIMESTAMP_INPUT_FORMAT' '=' STRING | 'TIMESTAMP_LTZ_OUTPUT_FORMAT' '=' STRING | 'TIMESTAMP_NTZ_OUTPUT_FORMAT' '=' STRING | 'TIMESTAMP_OUTPUT_FORMAT' '=' STRING | 'TIMESTAMP_TYPE_MAPPING' '=' STRING | 'TIMESTAMP_TZ_OUTPUT_FORMAT' '=' STRING | 'TIMEZONE' '=' STRING | 'TIME_INPUT_FORMAT' '=' STRING | 'TIME_OUTPUT_FORMAT' '=' STRING | 'TRANSACTION_DEFAULT_ISOLATION_LEVEL' '=' STRING | 'TWO_DIGIT_CENTURY_START' '=' NUMBER | 'UNSUPPORTED_DDL_ACTION' '=' STRING | 'USE_CACHED_RESULT' '=' ('TRUE' | 'FALSE') | 'WEEK_OF_YEAR_POLICY' '=' NUMBER | 'WEEK_START' '=' NUMBER)+
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-user-add-key-pair
alterUserAddKeyPair
    : 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'ADD' 'KEY' 'PAIR' ID 'PUBLIC_KEY' '=' STRING ('ROLE_RESTRICTION' '=' STRING)? ('DAYS_TO_EXPIRY' '=' NUMBER)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-user-add-programmatic-access-token
alterUserAddProgrammaticAccessToken
    : 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'ADD' ('PROGRAMMATIC' 'ACCESS' 'TOKEN' | 'PAT') ('IF' 'NOT' 'EXISTS')? ID ('ROLE_RESTRICTION' '=' STRING)? ('DAYS_TO_EXPIRY' '=' NUMBER)? ('MINS_TO_BYPASS_NETWORK_POLICY_REQUIREMENT' '=' NUMBER)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-user-add-workload-identity
alterUserAddWorkloadIdentity
    : 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'ADD' 'WORKLOAD' 'IDENTITY' ID 'TYPE' '=' ('AWS' | 'AZURE' | 'GCP' | 'OIDC') ('ARN' '=' STRING)? ('ISSUER' '=' STRING)? ('SUBJECT' '=' STRING)? ('OIDC_AUDIENCE_LIST' '=' '(' STRING (',' STRING)* ')')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-user-modify-key-pair
alterUserModifyKeyPair
    : 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'MODIFY' 'KEY' 'PAIR' ID 'RENAME' 'TO' ID
    | 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'MODIFY' 'KEY' 'PAIR' ID 'SET' ('DISABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'MODIFY' 'KEY' 'PAIR' ID 'UNSET' ('DISABLED' | 'COMMENT') (',' ('DISABLED' | 'COMMENT'))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-user-modify-programmatic-access-token
alterUserModifyProgrammaticAccessToken
    : 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'MODIFY' ('PROGRAMMATIC' 'ACCESS' 'TOKEN' | 'PAT') ('IF' 'EXISTS')? ID 'RENAME' 'TO' ID
    | 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'MODIFY' ('PROGRAMMATIC' 'ACCESS' 'TOKEN' | 'PAT') ('IF' 'EXISTS')? ID 'SET' ('DISABLED' '=' ('TRUE' | 'FALSE') ('MINS_TO_BYPASS_NETWORK_POLICY_REQUIREMENT' '=' NUMBER)? ('COMMENT' '=' STRING)? | 'MINS_TO_BYPASS_NETWORK_POLICY_REQUIREMENT' '=' NUMBER ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'MODIFY' ('PROGRAMMATIC' 'ACCESS' 'TOKEN' | 'PAT') ('IF' 'EXISTS')? ID 'UNSET' ('DISABLED' | 'MINS_TO_BYPASS_NETWORK_POLICY_REQUIREMENT' | 'COMMENT') (',' ('DISABLED' | 'MINS_TO_BYPASS_NETWORK_POLICY_REQUIREMENT' | 'COMMENT'))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-user-modify-workload-identity
alterUserModifyWorkloadIdentity
    : 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'MODIFY' 'WORKLOAD' 'IDENTITY' ID 'RENAME' 'TO' ID
    | 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'MODIFY' 'WORKLOAD' 'IDENTITY' ID 'SET' ('DISABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)? | 'COMMENT' '=' STRING)
    | 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'MODIFY' 'WORKLOAD' 'IDENTITY' ID 'UNSET' ('DISABLED' | 'COMMENT') (',' ('DISABLED' | 'COMMENT'))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-user-remove-key-pair
alterUserRemoveKeyPair
    : 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'REMOVE' 'KEY' 'PAIR' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-user-remove-programmatic-access-token
alterUserRemoveProgrammaticAccessToken
    : 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'REMOVE' ('PROGRAMMATIC' 'ACCESS' 'TOKEN' | 'PAT') ('IF' 'EXISTS')? ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-user-remove-workload-identity
alterUserRemoveWorkloadIdentity
    : 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'REMOVE' 'WORKLOAD' 'IDENTITY' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-user-rotate-key-pair
alterUserRotateKeyPair
    : 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'ROTATE' 'KEY' 'PAIR' ID 'PUBLIC_KEY' '=' STRING ('EXPIRE_ROTATED_KEY_PAIR_AFTER_HOURS' '=' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-user-rotate-programmatic-access-token
alterUserRotateProgrammaticAccessToken
    : 'ALTER' 'USER' ('IF' 'EXISTS')? ID? 'ROTATE' ('PROGRAMMATIC' 'ACCESS' 'TOKEN' | 'PAT') ('IF' 'EXISTS')? ID ('EXPIRE_ROTATED_TOKEN_AFTER_HOURS' '=' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-user-set-ingress-policy
alterUserSetIngressPolicy
    : 'ALTER' 'USER' QNAME 'SET' 'INGRESS' 'POLICY' QNAME 'FORCE'?
    | 'ALTER' 'USER' QNAME 'UNSET' 'INGRESS' 'POLICY'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-variable
alterVariable
    : 'ALTER' 'VARIABLE' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'VARIABLE' ('IF' 'EXISTS')? QNAME 'SET' 'DEFAULT' EXPR
    | 'ALTER' 'VARIABLE' ('IF' 'EXISTS')? QNAME 'DROP' 'DEFAULT'
    | 'ALTER' 'VARIABLE' ('IF' 'EXISTS')? QNAME 'SET' 'COMMENT' '=' STRING
    | 'ALTER' 'VARIABLE' ('IF' 'EXISTS')? QNAME 'UNSET' 'COMMENT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-view
// fix: list_of_properties is a space-delimited list of the properties the parameters name
// fix: a view's properties are set one at a time, as the SET parameter says
// fix: a view's properties are unset one at a time, as the UNSET parameter says
alterView
    : 'ALTER' 'VIEW' ('IF' 'EXISTS')? QNAME 'RENAME' 'TO' QNAME
    | 'ALTER' 'VIEW' ('IF' 'EXISTS')? QNAME 'SET' ('SECURE' | 'CHANGE_TRACKING' '=' ('TRUE' | 'FALSE') | 'CONTACT' ID '=' ID (',' ID '=' ID)* | 'COMMENT' '=' STRING)
    | 'ALTER' 'VIEW' ('IF' 'EXISTS')? QNAME 'UNSET' ('SECURE' | 'CONTACT' ID | 'COMMENT')
    | 'ALTER' 'VIEW' QNAME alterView_dataMetricFunctionAction
    | 'ALTER' 'VIEW' ('IF' 'EXISTS')? QNAME alterView_dataGovnPolicyTagAction
    ;
alterView_dataMetricFunctionAction
    : 'SET' 'DATA_METRIC_SCHEDULE' '=' (STRING | '\'TRIGGER_ON_CHANGES\'')
    | 'UNSET' 'DATA_METRIC_SCHEDULE'
    | 'ADD' 'DATA' 'METRIC' 'FUNCTION' ID 'ON' '(' ID (',' ID)* (',' 'TABLE' QNAME '(' ID (',' ID)* ')')? ')' ('WITHIN' 'GROUP' '(' ID (',' ID)* ')' ('GROUP' 'LIMIT' NUMBER)?)? ('FILTER' '(' EXPR ')')? ('EXPECTATION' ID '(' EXPR ')' (',' ID '(' EXPR ')')*)? ('EXECUTE' 'AS' 'ROLE' ID)? ('ANOMALY_DETECTION' '=' ('TRUE' | 'FALSE'))? ('SENSITIVITY' '=' ('\'LOW\'' | '\'MEDIUM\'' | '\'HIGH\''))? (',' ID 'ON' '(' ID (',' ID)* (',' 'TABLE' QNAME '(' ID (',' ID)* ')')? ')')? ('WITHIN' 'GROUP' '(' ID (',' ID)* ')' ('GROUP' 'LIMIT' NUMBER)?)? ('FILTER' '(' EXPR ')')? ('EXPECTATION' ID '(' EXPR ')' (',' ID '(' EXPR ')')*)? ('EXECUTE' 'AS' 'ROLE' ID)? ('ANOMALY_DETECTION' '=' ('TRUE' | 'FALSE'))? ('SENSITIVITY' '=' ('\'LOW\'' | '\'MEDIUM\'' | '\'HIGH\''))?
    | 'DROP' 'DATA' 'METRIC' 'FUNCTION' ID 'ON' '(' ID (',' ID)* (',' 'TABLE' QNAME '(' ID (',' ID)* ')')? ')' (',' ID 'ON' '(' ID (',' ID)* (',' 'TABLE' QNAME '(' ID (',' ID)* ')')? ')')?
    | 'MODIFY' 'DATA' 'METRIC' 'FUNCTION' ID 'ON' '(' ID (',' ID)* (',' 'TABLE' QNAME '(' ID (',' ID)* ')')? ')' ('SUSPEND' | 'RESUME') (',' ID 'ON' '(' ID (',' ID)* (',' 'TABLE' QNAME '(' ID (',' ID)* ')')? ')' ('SUSPEND' | 'RESUME'))?
    | 'MODIFY' 'DATA' 'METRIC' 'FUNCTION' ID 'ON' '(' ID (',' ID)* (',' 'TABLE' QNAME '(' ID (',' ID)* ')')? ')' ('ADD' | 'MODIFY') 'EXPECTATION' ID '(' EXPR ')' (',' ID '(' EXPR ')')*
    | 'MODIFY' 'DATA' 'METRIC' 'FUNCTION' ID 'ON' '(' ID (',' ID)* (',' 'TABLE' QNAME '(' ID (',' ID)* ')')? ')' 'DROP' 'EXPECTATION' ID (',' ID)*
    | 'MODIFY' 'DATA' 'METRIC' 'FUNCTION' ID 'ON' '(' ID (',' ID)* (',' 'TABLE' QNAME '(' ID (',' ID)* ')')? ')' 'FILTER' '(' EXPR? ')' (',' ID 'ON' '(' ID (',' ID)* (',' 'TABLE' QNAME '(' ID (',' ID)* ')')? ')' 'FILTER' '(' EXPR? ')')?
    | 'MODIFY' 'DATA' 'METRIC' 'FUNCTION' ID 'SET' ('ANOMALY_DETECTION' '=' ('TRUE' | 'FALSE') | 'SENSITIVITY' '=' ('\'LOW\'' | '\'MEDIUM\'' | '\'HIGH\'') | 'DATA_QUALITY_NOTIFICATION' '=' ('TRUE' | 'FALSE'))+
    ;
alterView_dataGovnPolicyTagAction
    : ('SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)* | 'UNSET' 'TAG' QNAME (',' QNAME)*)
    | ('ADD' 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')' | 'DROP' 'ROW' 'ACCESS' 'POLICY' QNAME | 'DROP' 'ROW' 'ACCESS' 'POLICY' QNAME ',' 'ADD' 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')' | 'DROP' 'ALL' 'ROW' 'ACCESS' 'POLICIES')
    | ('SET' 'AGGREGATION' 'POLICY' QNAME ('ENTITY' 'KEY' '(' ID (',' ID)* ')')? 'FORCE'? | 'UNSET' 'AGGREGATION' 'POLICY')
    | ('SET' 'JOIN' 'POLICY' QNAME 'FORCE'? | 'UNSET' 'JOIN' 'POLICY')
    | 'ADD' 'COLUMN'? ('IF' 'NOT' 'EXISTS')? ID DATA_TYPE ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | (('ALTER' | 'MODIFY') 'COLUMN'? ID 'SET' 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')? 'FORCE'? | ('ALTER' | 'MODIFY') 'COLUMN'? ID 'UNSET' 'MASKING' 'POLICY')
    | (('ALTER' | 'MODIFY') 'COLUMN'? ID 'SET' 'PROJECTION' 'POLICY' QNAME 'FORCE'? | ('ALTER' | 'MODIFY') 'COLUMN'? ID 'UNSET' 'PROJECTION' 'POLICY')
    | ('ALTER' | 'MODIFY') 'COLUMN'? ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)* ',' 'COLUMN'? ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | ('ALTER' | 'MODIFY') 'COLUMN'? ID 'UNSET' 'TAG' QNAME (',' QNAME)* ',' 'COLUMN'? ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/alter-warehouse
// fix: the sizes include the synonyms the parameters list
alterWarehouse
    : 'ALTER' 'WAREHOUSE' ('IF' 'EXISTS')? ID? ('SUSPEND' | 'RESUME' ('IF' 'SUSPENDED')?)
    | 'ALTER' 'WAREHOUSE' ('IF' 'EXISTS')? ID ('ENABLE' | 'DISABLE')
    | 'ALTER' 'WAREHOUSE' ('IF' 'EXISTS')? ID? 'ABORT' 'ALL' 'QUERIES'
    | 'ALTER' 'WAREHOUSE' ('IF' 'EXISTS')? ID 'RENAME' 'TO' ID
    | 'ALTER' 'WAREHOUSE' ('IF' 'EXISTS')? ID 'SET' (alterWarehouse_objectProperties alterWarehouse_objectParams? | alterWarehouse_objectParams)
    | 'ALTER' 'WAREHOUSE' ('IF' 'EXISTS')? ID 'SET' 'TAG' QNAME '=' STRING (',' QNAME '=' STRING)*
    | 'ALTER' 'WAREHOUSE' ('IF' 'EXISTS')? ID 'UNSET' 'TAG' QNAME (',' QNAME)*
    | 'ALTER' 'WAREHOUSE' ('IF' 'EXISTS')? ID 'UNSET' ID (',' ID)*
    | 'ALTER' 'WAREHOUSE' ('IF' 'EXISTS')? ID 'ADD' 'TABLES' '(' QNAME (',' QNAME)* ')'
    | 'ALTER' 'WAREHOUSE' ('IF' 'EXISTS')? ID 'DROP' 'TABLES' '(' QNAME (',' QNAME)* ')'
    ;
alterWarehouse_objectProperties
    : ('WAREHOUSE_TYPE' '=' ('STANDARD' | '\'SNOWPARK-OPTIMIZED\'' | 'ADAPTIVE') | 'WAREHOUSE_SIZE' '=' ('XSMALL' | '\'X-SMALL\'' | 'SMALL' | 'MEDIUM' | 'LARGE' | 'XLARGE' | '\'X-LARGE\'' | 'XXLARGE' | 'X2LARGE' | '\'2X-LARGE\'' | 'XXXLARGE' | 'X3LARGE' | '\'3X-LARGE\'' | 'X4LARGE' | '\'4X-LARGE\'' | 'X5LARGE' | '\'5X-LARGE\'' | 'X6LARGE' | '\'6X-LARGE\'') | 'GENERATION' '=' ('\'1\'' | '\'2\'') | 'RESOURCE_CONSTRAINT' '=' ('STANDARD_GEN_1' | 'STANDARD_GEN_2' | 'MEMORY_1X' | 'MEMORY_1X_x86' | 'MEMORY_16X' | 'MEMORY_16X_x86' | 'MEMORY_64X' | 'MEMORY_64X_x86') | 'WAIT_FOR_COMPLETION' '=' ('TRUE' | 'FALSE') | 'MAX_CLUSTER_COUNT' '=' NUMBER | 'MIN_CLUSTER_COUNT' '=' NUMBER | 'SCALING_POLICY' '=' ('STANDARD' | 'ECONOMY') | 'AUTO_SUSPEND' '=' (NUMBER | 'NULL') | 'AUTO_RESUME' '=' ('TRUE' | 'FALSE') | 'RESOURCE_MONITOR' '=' ID | 'COMMENT' '=' STRING | 'ENABLE_QUERY_ACCELERATION' '=' ('TRUE' | 'FALSE') | 'QUERY_ACCELERATION_MAX_SCALE_FACTOR' '=' NUMBER)+
    ;
alterWarehouse_objectParams
    : ('MAX_CONCURRENCY_LEVEL' '=' NUMBER | 'STATEMENT_QUEUED_TIMEOUT_IN_SECONDS' '=' NUMBER | 'STATEMENT_TIMEOUT_IN_SECONDS' '=' NUMBER)+
    ;

// https://docs.snowflake.com/en/sql-reference/sql/begin
begin
    : 'BEGIN' (('WORK' | 'TRANSACTION'))? ('NAME' QNAME)?
    | 'START' 'TRANSACTION' ('NAME' QNAME)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/call
call
    : 'CALL' QNAME '(' ((ID '=>')? EXPR (',' (ID '=>')? EXPR)*)? ')' ('INTO' ':' ID)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/call-with
// not converted (form): ... stands for text the docs leave out
//   WITH <name> AS PROCEDURE ([ <arg_name> <arg_data_type> ]) [ , ... ] )
//     RETURNS { <result_data_type> [ [ NOT ] NULL ] | TABLE ( [ <col_name> <col_data_type> [ , ... ] ] ) }
//     LANGUAGE { JAVA }
//     RUNTIME_VERSION = '<scala_or_java_runtime_version>'
//     PACKAGES = ( 'com.snowflake:snowpark:<version>' [, '<package_name_and_version>' ...] )
//     [ IMPORTS = ( '<stage_path_and_directory_or_file_name_to_read>' [, '<stage_path_and_directory_or_file_name_to_read>' ...] ) ]
//     HANDLER = '<fully_qualified_method_name>'
//     [ { CALLED ON NULL INPUT | { RETURNS NULL ON NULL INPUT | STRICT } } ]
//     [ AS '<procedure_definition>' ]
//     [ , <cte_nameN> [ ( <cte_column_list> ) ] AS ( SELECT ...  ) ]
// not converted (form): ... stands for text the docs leave out
//   WITH <name> AS PROCEDURE ([ <arg_name> <arg_data_type> ]) [ , ... ] )
//     RETURNS { <result_data_type> [ [ NOT ] NULL ] | TABLE ( [ <col_name> <col_data_type> [ , ... ] ] ) }
//     LANGUAGE { JAVA }
//     RUNTIME_VERSION = '<scala_or_java_runtime_version>'
//     PACKAGES = ( 'com.snowflake:snowpark:<version>' [, '<package_name_and_version>' ...] )
//     [ IMPORTS = ( '<stage_path_and_directory_or_file_name_to_read>' [, '<stage_path_and_directory_or_file_name_to_read>' ...] ) ]
//     HANDLER = '<fully_qualified_method_name>'
//     [ { CALLED ON NULL INPUT | { RETURNS NULL ON NULL INPUT | STRICT } } ]
//     [ , <cte_nameN> [ ( <cte_column_list> ) ] AS ( SELECT ...  ) ]
// not converted (form): ... stands for text the docs leave out
//   WITH <name> AS PROCEDURE ([ <arg_name> <arg_data_type> ]) [ , ... ] )
//     RETURNS <result_data_type> [ [ NOT ] NULL ]
//     LANGUAGE JAVASCRIPT
//     [ { CALLED ON NULL INPUT | { RETURNS NULL ON NULL INPUT | STRICT } } ]
//     AS '<procedure_definition>'
//     [ , <cte_nameN> [ ( <cte_column_list> ) ] AS ( SELECT ...  ) ]
// not converted (form): ... stands for text the docs leave out
//   WITH <name> AS PROCEDURE ( [ <arg_name> <arg_data_type> ] [ , ... ] )
//     RETURNS { <result_data_type> [ [ NOT ] NULL ] | TABLE ( [ <col_name> <col_data_type> [ , ... ] ] ) }
//     LANGUAGE PYTHON
//     RUNTIME_VERSION = '<python_version>'
//     PACKAGES = ( 'snowflake-snowpark-python[==<version>]'[, '<package_name>[==<version>]' ... ])
//     [ IMPORTS = ( '<stage_path_and_directory_or_file_name_to_read>' [, '<stage_path_and_directory_or_file_name_to_read>' ...] ) ]
//     HANDLER = '<function_name>'
//     [ { CALLED ON NULL INPUT | { RETURNS NULL ON NULL INPUT | STRICT } } ]
//     [ , <cte_nameN> [ ( <cte_column_list> ) ] AS ( SELECT ...  ) ]
//     AS '<procedure_definition>'
// not converted (form): ... stands for text the docs leave out
//   WITH <name> AS PROCEDURE ( [ <arg_name> <arg_data_type> ] [ , ... ] )
//     RETURNS { <result_data_type> [ [ NOT ] NULL ] | TABLE ( [ <col_name> <col_data_type> [ , ... ] ] ) }
//     LANGUAGE PYTHON
//     RUNTIME_VERSION = '<python_version>'
//     PACKAGES = ( 'snowflake-snowpark-python[==<version>]'[, '<package_name>[==<version>]' ... ])
//     [ IMPORTS = ( '<stage_path_and_directory_or_file_name_to_read>' [, '<stage_path_and_directory_or_file_name_to_read>' ...] ) ]
//     HANDLER = '<module_file_name>.<function_name>'
//     [ { CALLED ON NULL INPUT | { RETURNS NULL ON NULL INPUT | STRICT } } ]
//     [ , <cte_nameN> [ ( <cte_column_list> ) ] AS ( SELECT ...  ) ]
// not converted (form): ... stands for text the docs leave out
//   WITH <name> AS PROCEDURE ([ <arg_name> <arg_data_type> ]) [ , ... ] )
//     RETURNS { <result_data_type> [ [ NOT ] NULL ] | TABLE ( [ <col_name> <col_data_type> [ , ... ] ] ) }
//     LANGUAGE { SCALA }
//     RUNTIME_VERSION = '<scala_runtime_version>'
//     PACKAGES = ( 'com.snowflake:snowpark_<scala_version>:<version>' [, '<package_name_and_version>' ...] )
//     [ IMPORTS = ( '<stage_path_and_directory_or_file_name_to_read>' [, '<stage_path_and_directory_or_file_name_to_read>' ...] ) ]
//     HANDLER = '<fully_qualified_method_name>'
//     [ { CALLED ON NULL INPUT | { RETURNS NULL ON NULL INPUT | STRICT } } ]
//     [ AS '<procedure_definition>' ]
//     [ , <cte_nameN> [ ( <cte_column_list> ) ] AS ( SELECT ...  ) ]
// not converted (form): ... stands for text the docs leave out
//   WITH <name> AS PROCEDURE ([ <arg_name> <arg_data_type> ]) [ , ... ] )
//     RETURNS { <result_data_type> [ [ NOT ] NULL ] | TABLE ( [ <col_name> <col_data_type> [ , ... ] ] ) }
//     LANGUAGE { SCALA }
//     RUNTIME_VERSION = '<scala_or_java_runtime_version>'
//     PACKAGES = ( 'com.snowflake:snowpark_<scala_version>:<version>' [, '<package_name_and_version>' ...] )
//     [ IMPORTS = ( '<stage_path_and_directory_or_file_name_to_read>' [, '<stage_path_and_directory_or_file_name_to_read>' ...] ) ]
//     HANDLER = '<fully_qualified_method_name>'
//     [ { CALLED ON NULL INPUT | { RETURNS NULL ON NULL INPUT | STRICT } } ]
//     [ , <cte_nameN> [ ( <cte_column_list> ) ] AS ( SELECT ...  ) ]
// not converted (form): ... stands for text the docs leave out
//   WITH <name> AS PROCEDURE ([ <arg_name> <arg_data_type> ]) [ , ... ] )
//     RETURNS { <result_data_type> | TABLE ( [ <col_name> <col_data_type> [ , ... ] ] ) }
//     LANGUAGE SQL
//     [ { CALLED ON NULL INPUT | { RETURNS NULL ON NULL INPUT | STRICT } } ]
//     AS '<procedure_definition>'
//     [ , <cte_nameN> [ ( <cte_column_list> ) ] AS ( SELECT ...  ) ]
callWith
    : 'CALL' QNAME '(' ((ID '=>')? EXPR (',' (ID '=>')? EXPR)*)? ')' ('INTO' ':' ID)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/comment
// fix: the object type is TABLE, SCHEMA, VIEW and so on, as the parameters say
comment
    : 'COMMENT' ('IF' 'EXISTS')? 'ON' ('TABLE' | 'SCHEMA' | 'VIEW' | OBJECT_TYPE) QNAME 'IS' STRING
    | 'COMMENT' ('IF' 'EXISTS')? 'ON' 'COLUMN' QNAME '.' ID 'IS' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/commit
commit
    : 'COMMIT' 'WORK'?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/copy-files
// not converted (form): ... stands for text the docs leave out
//   COPY FILES INTO @[<namespace>.]<stage_name>[/<path>/]
//     FROM ( SELECT <existing_url> [ , <new_filename> ] FROM ... )
//     [ DETAILED_OUTPUT = { TRUE | FALSE } ]
copyFiles
    : 'COPY' 'FILES' 'INTO' STAGE 'FROM' STAGE ('FILES' '=' '(' STRING (',' STRING)* ')')? ('PATTERN' '=' STRING)? ('DETAILED_OUTPUT' '=' ('TRUE' | 'FALSE'))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/copy-into-location
// fix: HEADER takes TRUE or FALSE, as the page's parameters and examples write it
// fix: read: each -- If TYPE = X section of formatTypeOptions is the options of TYPE = X alone, and TYPE = { ... } [ formatTypeOptions ] chooses a TYPE with its own options
copyIntoLocation
    : 'COPY' 'INTO' (copyIntoLocation_internalStage | copyIntoLocation_externalStage | copyIntoLocation_externalLocation) 'FROM' ((ID '.' (ID '.')?)? ID | '(' QUERY ')') ('PARTITION' 'BY' EXPR)? ('FILE_FORMAT' '=' '(' ('FORMAT_NAME' '=' STRING | ('TYPE' '=' 'CSV' copyIntoLocation_formatTypeOptions_CSV? | 'TYPE' '=' 'JSON' copyIntoLocation_formatTypeOptions_JSON? | 'TYPE' '=' 'PARQUET' copyIntoLocation_formatTypeOptions_PARQUET?)) ')')? copyIntoLocation_copyOptions? ('VALIDATION_MODE' '=' 'RETURN_ROWS')? ('HEADER' '=' ('TRUE' | 'FALSE'))?
    ;
copyIntoLocation_internalStage
    : STAGE
    ;
copyIntoLocation_externalStage
    : STAGE
    ;
copyIntoLocation_externalLocation
    : STRING ('STORAGE_INTEGRATION' '=' ID | 'CREDENTIALS' '=' '(' 'AWS_KEY_ID' '=' STRING 'AWS_SECRET_KEY' '=' STRING ('AWS_TOKEN' '=' STRING)? ')')? ('ENCRYPTION' '=' '(' (('TYPE' '=' '\'AWS_CSE\'')? ('MASTER_KEY' '=' STRING)? | ('TYPE' '=' '\'AWS_SSE_S3\'')? | ('TYPE' '=' '\'AWS_SSE_KMS\'' ('KMS_KEY_ID' '=' STRING)?)? | ('TYPE' '=' '\'NONE\'')?) ')')?
    | STRING ('STORAGE_INTEGRATION' '=' ID)? ('ENCRYPTION' '=' '(' (('TYPE' '=' '\'GCS_SSE_KMS\'')? ('KMS_KEY_ID' '=' STRING)? | ('TYPE' '=' '\'NONE\'')?) ')')?
    | STRING ('STORAGE_INTEGRATION' '=' ID | 'CREDENTIALS' '=' '(' ('AZURE_SAS_TOKEN' '=' STRING)? ')')? ('ENCRYPTION' '=' '(' ('TYPE' '=' ('\'AZURE_CSE\'' | '\'NONE\''))? ('MASTER_KEY' '=' STRING)? ')')?
    ;
copyIntoLocation_formatTypeOptions_CSV
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'RECORD_DELIMITER' '=' (STRING | 'NONE') | 'FIELD_DELIMITER' '=' (STRING | 'NONE') | 'FILE_EXTENSION' '=' STRING | 'ESCAPE' '=' (STRING | 'NONE') | 'ESCAPE_UNENCLOSED_FIELD' '=' (STRING | 'NONE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'FIELD_OPTIONALLY_ENCLOSED_BY' '=' (STRING | 'NONE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')' | 'EMPTY_FIELD_AS_NULL' '=' ('TRUE' | 'FALSE')) (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'RECORD_DELIMITER' '=' (STRING | 'NONE') | 'FIELD_DELIMITER' '=' (STRING | 'NONE') | 'FILE_EXTENSION' '=' STRING | 'ESCAPE' '=' (STRING | 'NONE') | 'ESCAPE_UNENCLOSED_FIELD' '=' (STRING | 'NONE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'FIELD_OPTIONALLY_ENCLOSED_BY' '=' (STRING | 'NONE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')' | 'EMPTY_FIELD_AS_NULL' '=' ('TRUE' | 'FALSE')))*
    ;
copyIntoLocation_formatTypeOptions_JSON
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'FILE_EXTENSION' '=' STRING) (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'FILE_EXTENSION' '=' STRING))*
    ;
copyIntoLocation_formatTypeOptions_PARQUET
    : ('COMPRESSION' '=' ('AUTO' | 'LZO' | 'SNAPPY' | 'NONE') | 'SNAPPY_COMPRESSION' '=' ('TRUE' | 'FALSE')) (','? ('COMPRESSION' '=' ('AUTO' | 'LZO' | 'SNAPPY' | 'NONE') | 'SNAPPY_COMPRESSION' '=' ('TRUE' | 'FALSE')))*
    ;
copyIntoLocation_copyOptions
    : ('OVERWRITE' '=' ('TRUE' | 'FALSE') | 'SINGLE' '=' ('TRUE' | 'FALSE') | 'MAX_FILE_SIZE' '=' NUMBER | 'INCLUDE_QUERY_ID' '=' ('TRUE' | 'FALSE') | 'DETAILED_OUTPUT' '=' ('TRUE' | 'FALSE')) (','? ('OVERWRITE' '=' ('TRUE' | 'FALSE') | 'SINGLE' '=' ('TRUE' | 'FALSE') | 'MAX_FILE_SIZE' '=' NUMBER | 'INCLUDE_QUERY_ID' '=' ('TRUE' | 'FALSE') | 'DETAILED_OUTPUT' '=' ('TRUE' | 'FALSE')))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/copy-into-table
// fix: VALIDATION_MODE: RETURN_<n>_ROWS holds its number inside the keyword; written for n = 10
// fix: the parameters list the valid METADATA$ columns
// fix: read: each -- If TYPE = X section of formatTypeOptions is the options of TYPE = X alone, and TYPE = { ... } [ formatTypeOptions ] chooses a TYPE with its own options
copyIntoTable
    : 'COPY' 'INTO' (ID '.' (ID '.')?)? ID 'FROM' (copyIntoTable_internalStage | copyIntoTable_externalStage | copyIntoTable_externalLocation) ('FILES' '=' '(' STRING (',' STRING)* ')')? ('PATTERN' '=' STRING)? ('FILE_FORMAT' '=' '(' ('FORMAT_NAME' '=' STRING | ('TYPE' '=' 'CSV' copyIntoTable_formatTypeOptions_CSV? | 'TYPE' '=' 'JSON' copyIntoTable_formatTypeOptions_JSON? | 'TYPE' '=' 'AVRO' copyIntoTable_formatTypeOptions_AVRO? | 'TYPE' '=' 'ORC' copyIntoTable_formatTypeOptions_ORC? | 'TYPE' '=' 'PARQUET' copyIntoTable_formatTypeOptions_PARQUET? | 'TYPE' '=' 'XML' copyIntoTable_formatTypeOptions_XML?)) ')')? copyIntoTable_copyOptions? ('VALIDATION_MODE' '=' 'RETURN_10_ROWS' | 'VALIDATION_MODE' '=' 'RETURN_ERRORS' | 'VALIDATION_MODE' '=' 'RETURN_ALL_ERRORS')?
    | 'COPY' 'INTO' (ID '.' (ID '.')?)? ID ('(' ID (',' ID)* ')')? 'FROM' '(' 'SELECT' (ID '.')? POSITION ('.' ID)? (',' (ID '.')? POSITION ('.' ID)?)* 'FROM' (copyIntoTable_internalStage | copyIntoTable_externalStage) ')' ('FILES' '=' '(' STRING (',' STRING)* ')')? ('PATTERN' '=' STRING)? ('FILE_FORMAT' '=' '(' ('FORMAT_NAME' '=' STRING | ('TYPE' '=' 'CSV' copyIntoTable_formatTypeOptions_CSV? | 'TYPE' '=' 'JSON' copyIntoTable_formatTypeOptions_JSON? | 'TYPE' '=' 'AVRO' copyIntoTable_formatTypeOptions_AVRO? | 'TYPE' '=' 'ORC' copyIntoTable_formatTypeOptions_ORC? | 'TYPE' '=' 'PARQUET' copyIntoTable_formatTypeOptions_PARQUET? | 'TYPE' '=' 'XML' copyIntoTable_formatTypeOptions_XML?)) ')')? copyIntoTable_copyOptions?
    ;
copyIntoTable_internalStage
    : STAGE
    ;
copyIntoTable_externalStage
    : STAGE
    ;
copyIntoTable_externalLocation
    : STRING ('STORAGE_INTEGRATION' '=' ID | 'CREDENTIALS' '=' '(' 'AWS_KEY_ID' '=' STRING 'AWS_SECRET_KEY' '=' STRING ('AWS_TOKEN' '=' STRING)? ')')? ('ENCRYPTION' '=' '(' (('TYPE' '=' '\'AWS_CSE\'')? ('MASTER_KEY' '=' STRING)? | ('TYPE' '=' '\'AWS_SSE_S3\'')? | ('TYPE' '=' '\'AWS_SSE_KMS\'' ('KMS_KEY_ID' '=' STRING)?)? | ('TYPE' '=' '\'NONE\'')?) ')')?
    | STRING ('STORAGE_INTEGRATION' '=' ID)? ('ENCRYPTION' '=' '(' (('TYPE' '=' '\'GCS_SSE_KMS\'')? ('KMS_KEY_ID' '=' STRING)? | ('TYPE' '=' '\'NONE\'')?) ')')?
    | STRING ('STORAGE_INTEGRATION' '=' ID | 'CREDENTIALS' '=' '(' ('AZURE_SAS_TOKEN' '=' STRING)? ')')? ('ENCRYPTION' '=' '(' ('TYPE' '=' ('\'AZURE_CSE\'' | '\'NONE\''))? ('MASTER_KEY' '=' STRING)? ')')?
    ;
copyIntoTable_formatTypeOptions_CSV
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'RECORD_DELIMITER' '=' (STRING | 'NONE') | 'FIELD_DELIMITER' '=' (STRING | 'NONE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'PARSE_HEADER' '=' ('TRUE' | 'FALSE') | 'SKIP_HEADER' '=' NUMBER | 'SKIP_BLANK_LINES' '=' ('TRUE' | 'FALSE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'ESCAPE' '=' (STRING | 'NONE') | 'ESCAPE_UNENCLOSED_FIELD' '=' (STRING | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'FIELD_OPTIONALLY_ENCLOSED_BY' '=' (STRING | 'NONE') | 'NULL_IF' '=' '(' (STRING (',' STRING)*)? ')' | 'ERROR_ON_COLUMN_COUNT_MISMATCH' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'EMPTY_FIELD_AS_NULL' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE') | 'ENCODING' '=' (STRING | 'UTF8')) (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'RECORD_DELIMITER' '=' (STRING | 'NONE') | 'FIELD_DELIMITER' '=' (STRING | 'NONE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'PARSE_HEADER' '=' ('TRUE' | 'FALSE') | 'SKIP_HEADER' '=' NUMBER | 'SKIP_BLANK_LINES' '=' ('TRUE' | 'FALSE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'ESCAPE' '=' (STRING | 'NONE') | 'ESCAPE_UNENCLOSED_FIELD' '=' (STRING | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'FIELD_OPTIONALLY_ENCLOSED_BY' '=' (STRING | 'NONE') | 'NULL_IF' '=' '(' (STRING (',' STRING)*)? ')' | 'ERROR_ON_COLUMN_COUNT_MISMATCH' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'EMPTY_FIELD_AS_NULL' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE') | 'ENCODING' '=' (STRING | 'UTF8')))*
    ;
copyIntoTable_formatTypeOptions_JSON
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' (STRING (',' STRING)*)? ')' | 'ENABLE_OCTAL' '=' ('TRUE' | 'FALSE') | 'ALLOW_DUPLICATE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ARRAY' '=' ('TRUE' | 'FALSE') | 'STRIP_NULL_VALUES' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'IGNORE_UTF8_ERRORS' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE')) (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' (STRING (',' STRING)*)? ')' | 'ENABLE_OCTAL' '=' ('TRUE' | 'FALSE') | 'ALLOW_DUPLICATE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ARRAY' '=' ('TRUE' | 'FALSE') | 'STRIP_NULL_VALUES' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'IGNORE_UTF8_ERRORS' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE')))*
    ;
copyIntoTable_formatTypeOptions_AVRO
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' (STRING (',' STRING)*)? ')') (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' (STRING (',' STRING)*)? ')'))*
    ;
copyIntoTable_formatTypeOptions_ORC
    : ('TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' (STRING (',' STRING)*)? ')') (','? ('TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' (STRING (',' STRING)*)? ')'))*
    ;
copyIntoTable_formatTypeOptions_PARQUET
    : ('COMPRESSION' '=' ('AUTO' | 'SNAPPY' | 'NONE') | 'BINARY_AS_TEXT' '=' ('TRUE' | 'FALSE') | 'USE_LOGICAL_TYPE' '=' ('TRUE' | 'FALSE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'USE_VECTORIZED_SCANNER' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' (STRING (',' STRING)*)? ')') (','? ('COMPRESSION' '=' ('AUTO' | 'SNAPPY' | 'NONE') | 'BINARY_AS_TEXT' '=' ('TRUE' | 'FALSE') | 'USE_LOGICAL_TYPE' '=' ('TRUE' | 'FALSE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'USE_VECTORIZED_SCANNER' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' (STRING (',' STRING)*)? ')'))*
    ;
copyIntoTable_formatTypeOptions_XML
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'IGNORE_UTF8_ERRORS' '=' ('TRUE' | 'FALSE') | 'PRESERVE_SPACE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ELEMENT' '=' ('TRUE' | 'FALSE') | 'DISABLE_AUTO_CONVERT' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE')) (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'IGNORE_UTF8_ERRORS' '=' ('TRUE' | 'FALSE') | 'PRESERVE_SPACE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ELEMENT' '=' ('TRUE' | 'FALSE') | 'DISABLE_AUTO_CONVERT' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE')))*
    ;
copyIntoTable_copyOptions
    : ('CLUSTER_AT_INGEST_TIME' '=' ('TRUE' | 'FALSE') | 'ENFORCE_LENGTH' '=' ('TRUE' | 'FALSE') | 'FORCE' '=' ('TRUE' | 'FALSE') | 'INCLUDE_METADATA' '=' '(' ID '=' ('METADATA$FILENAME' | 'METADATA$FILE_ROW_NUMBER' | 'METADATA$FILE_CONTENT_KEY' | 'METADATA$FILE_LAST_MODIFIED' | 'METADATA$START_SCAN_TIME') (',' ID '=' ('METADATA$FILENAME' | 'METADATA$FILE_ROW_NUMBER' | 'METADATA$FILE_CONTENT_KEY' | 'METADATA$FILE_LAST_MODIFIED' | 'METADATA$START_SCAN_TIME'))* ')' | 'LOAD_MODE' '=' ('FULL_INGEST' | 'ADD_FILES_COPY') | 'LOAD_UNCERTAIN_FILES' '=' ('TRUE' | 'FALSE') | 'MATCH_BY_COLUMN_NAME' '=' ('CASE_SENSITIVE' | 'CASE_INSENSITIVE' | 'NONE') | 'ON_ERROR' '=' ('CONTINUE' | 'SKIP_FILE' | 'SKIP_FILE_1' | STRING | 'ABORT_STATEMENT') | 'PURGE' '=' ('TRUE' | 'FALSE') | 'RETURN_FAILED_ONLY' '=' ('TRUE' | 'FALSE') | 'SIZE_LIMIT' '=' NUMBER | 'TRUNCATECOLUMNS' '=' ('TRUE' | 'FALSE')) (','? ('CLUSTER_AT_INGEST_TIME' '=' ('TRUE' | 'FALSE') | 'ENFORCE_LENGTH' '=' ('TRUE' | 'FALSE') | 'FORCE' '=' ('TRUE' | 'FALSE') | 'INCLUDE_METADATA' '=' '(' ID '=' ('METADATA$FILENAME' | 'METADATA$FILE_ROW_NUMBER' | 'METADATA$FILE_CONTENT_KEY' | 'METADATA$FILE_LAST_MODIFIED' | 'METADATA$START_SCAN_TIME') (',' ID '=' ('METADATA$FILENAME' | 'METADATA$FILE_ROW_NUMBER' | 'METADATA$FILE_CONTENT_KEY' | 'METADATA$FILE_LAST_MODIFIED' | 'METADATA$START_SCAN_TIME'))* ')' | 'LOAD_MODE' '=' ('FULL_INGEST' | 'ADD_FILES_COPY') | 'LOAD_UNCERTAIN_FILES' '=' ('TRUE' | 'FALSE') | 'MATCH_BY_COLUMN_NAME' '=' ('CASE_SENSITIVE' | 'CASE_INSENSITIVE' | 'NONE') | 'ON_ERROR' '=' ('CONTINUE' | 'SKIP_FILE' | 'SKIP_FILE_1' | STRING | 'ABORT_STATEMENT') | 'PURGE' '=' ('TRUE' | 'FALSE') | 'RETURN_FAILED_ONLY' '=' ('TRUE' | 'FALSE') | 'SIZE_LIMIT' '=' NUMBER | 'TRUNCATECOLUMNS' '=' ('TRUE' | 'FALSE')))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create
// no syntax

// https://docs.snowflake.com/en/sql-reference/sql/create-account
createAccount
    : 'CREATE' 'ACCOUNT' ID 'ADMIN_NAME' '=' STRING ('ADMIN_PASSWORD' '=' STRING | 'ADMIN_RSA_PUBLIC_KEY' '=' STRING) ('ADMIN_USER_TYPE' '=' ('PERSON' | 'SERVICE' | 'LEGACY_SERVICE' | 'NULL'))? ('FIRST_NAME' '=' STRING)? ('LAST_NAME' '=' STRING)? 'EMAIL' '=' STRING ('MUST_CHANGE_PASSWORD' '=' ('TRUE' | 'FALSE'))? 'EDITION' '=' ('STANDARD' | 'ENTERPRISE' | 'BUSINESS_CRITICAL') ('REGION_GROUP' '=' ID)? ('REGION' '=' ID)? ('COMMENT' '=' STRING)? ('POLARIS' '=' ('TRUE' | 'FALSE'))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-agent
createAgent
    : 'CREATE' ('OR' 'REPLACE')? (('TEMP' | 'TEMPORARY'))? 'SECURE'? 'AGENT' QNAME ('COPY' 'GRANTS')? ('SECURE' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('PROFILE' '=' STRING)? 'FROM' 'SPECIFICATION' STRING
    | 'CREATE' (('TEMP' | 'TEMPORARY'))? 'SECURE'? 'AGENT' ('IF' 'NOT' 'EXISTS')? QNAME ('COPY' 'GRANTS')? ('SECURE' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('PROFILE' '=' STRING)? 'FROM' 'SPECIFICATION' STRING
    | 'CREATE' ('OR' 'REPLACE')? 'SECURE'? 'AGENT' QNAME ('COPY' 'GRANTS')? ('COMMENT' '=' STRING)? 'FROM' STAGE
    | 'CREATE' 'SECURE'? 'AGENT' ('IF' 'NOT' 'EXISTS')? QNAME ('COPY' 'GRANTS')? ('COMMENT' '=' STRING)? 'FROM' STAGE
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-aggregation-policy
createAggregationPolicy
    : 'CREATE' ('OR' 'REPLACE')? 'AGGREGATION' 'POLICY' QNAME 'AS' '(' ')' 'RETURNS' 'AGGREGATION_CONSTRAINT' '->' EXPR ('COMMENT' '=' STRING)?
    | 'CREATE' 'AGGREGATION' 'POLICY' ('IF' 'NOT' 'EXISTS')? QNAME 'AS' '(' ')' 'RETURNS' 'AGGREGATION_CONSTRAINT' '->' EXPR ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-ai-function
createAiFunction
    : 'CREATE' ('OR' 'REPLACE')? 'AI' 'FUNCTION' QNAME '(' (ID DATA_TYPE (',' ID DATA_TYPE)*)? ')' 'RETURNS' DATA_TYPE ('COMMENT' '=' STRING)? 'AS' STRING
    | 'CREATE' 'AI' 'FUNCTION' ('IF' 'NOT' 'EXISTS')? QNAME '(' (ID DATA_TYPE (',' ID DATA_TYPE)*)? ')' 'RETURNS' DATA_TYPE ('COMMENT' '=' STRING)? 'AS' STRING
    | 'CREATE' ('OR' 'REPLACE')? 'AI' 'FUNCTION' QNAME '(' (ID DATA_TYPE (',' ID DATA_TYPE)*)? ')' 'RETURNS' DATA_TYPE 'FROM' 'EXPERIMENT' QNAME 'RUN' ID
    | 'CREATE' 'AI' 'FUNCTION' ('IF' 'NOT' 'EXISTS')? QNAME '(' (ID DATA_TYPE (',' ID DATA_TYPE)*)? ')' 'RETURNS' DATA_TYPE 'FROM' 'EXPERIMENT' QNAME 'RUN' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-alert
// fix: <condition> and <action> are SQL statements, as the page's parameters say: a query, and any statement
// fix: <condition> and <action> are SQL statements, as the page's parameters say: a query, and any statement
// not converted (form): DEFINE runs only in DCM project files, not as a statement
//   DEFINE ALERT MY_DB.MY_SCHEMA.ALRT_CHECK_ORDER_VOLUME
//       WAREHOUSE = 'MY_WH'
//       SCHEDULE = '60 MINUTE'
//       STARTED
//       IF (EXISTS (
//         SELECT 1 FROM MY_DB.MY_SCHEMA.ORDERS WHERE ORDER_COUNT > 1000
//       ))
//       THEN
//         CALL SYSTEM$SEND_EMAIL(
//           'my_notification_integration', 'ops@example.com', 'High order volume', 'Order volume exceeded 1000 orders.'
//         )
//   ;
createAlert
    : 'CREATE' ('OR' 'REPLACE')? 'ALERT' QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('SCHEDULE' '=' STRING)? ('WAREHOUSE' '=' ID)? ('COMMENT' '=' STRING)? ('CONFIG' '=' STRING)? ('RUNBOOK' '=' STRING)? ('SUSPEND_ALERT_AFTER_NUM_FAILURES' '=' NUMBER)? 'IF' '(' 'EXISTS' '(' QUERY ')' ')' 'THEN' QUERY
    | 'CREATE' 'ALERT' ('IF' 'NOT' 'EXISTS')? QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('SCHEDULE' '=' STRING)? ('WAREHOUSE' '=' ID)? ('COMMENT' '=' STRING)? ('CONFIG' '=' STRING)? ('RUNBOOK' '=' STRING)? ('SUSPEND_ALERT_AFTER_NUM_FAILURES' '=' NUMBER)? 'IF' '(' 'EXISTS' '(' QUERY ')' ')' 'THEN' QUERY
    | 'CREATE' 'OR' 'ALTER' 'ALERT' QNAME ('SCHEDULE' '=' STRING)? ('WAREHOUSE' '=' ID)? ('COMMENT' '=' STRING)? ('CONFIG' '=' STRING)? ('RUNBOOK' '=' STRING)? ('SUSPEND_ALERT_AFTER_NUM_FAILURES' '=' NUMBER)? 'IF' '(' 'EXISTS' '(' QUERY ')' ')' 'THEN' QUERY
    | 'CREATE' ('OR' 'REPLACE')? 'ALERT' QNAME 'CLONE' QNAME
    | 'CREATE' ('OR' 'REPLACE')? 'ALERT' QNAME 'FROM' 'TEMPLATE' ID ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WAREHOUSE' '=' ID)? ('SCHEDULE' '=' STRING)? ('COMMENT' '=' STRING)? ('RUNBOOK' '=' STRING)? ('SUSPEND_ALERT_AFTER_NUM_FAILURES' '=' NUMBER)? ('TEMPLATE_PARAMS' '=' STRING)?
    | 'CREATE' 'ALERT' ('IF' 'NOT' 'EXISTS')? QNAME 'FROM' 'TEMPLATE' ID ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WAREHOUSE' '=' ID)? ('SCHEDULE' '=' STRING)? ('COMMENT' '=' STRING)? ('RUNBOOK' '=' STRING)? ('SUSPEND_ALERT_AFTER_NUM_FAILURES' '=' NUMBER)? ('TEMPLATE_PARAMS' '=' STRING)?
    | ID ':' ID ':' ID ',' ID ':' ID ':' ID? ',' ID ':' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-api-integration
// fix: read: ('<...>') after _PREFIXES = is a list of URLs, as the parameters say
// not converted (form): undefined fragment "oauth_parameters": the page does not define it
//   CREATE [ OR REPLACE ] API INTEGRATION [ IF NOT EXISTS ] <integration_name>
//     API_PROVIDER = git_https_api
//     API_ALLOWED_PREFIXES = ( '<prefix>' [ , '<prefix>' ... ] )
//     [ API_BLOCKED_PREFIXES = ( '<prefix>' [ , '<prefix>' ... ] ) ]
//     API_USER_AUTHENTICATION = (
//       TYPE = OAUTH2
//       {oauth_parameters}
//     )
//     ENABLED = { TRUE | FALSE }
//     [ COMMENT = '<string_literal>' ]
//     ;
createApiIntegration
    : 'CREATE' ('OR' 'REPLACE')? 'API' 'INTEGRATION' ID 'API_PROVIDER' '=' ('aws_api_gateway' | 'aws_private_api_gateway' | 'aws_gov_api_gateway' | 'aws_gov_private_api_gateway') 'API_AWS_ROLE_ARN' '=' STRING ('API_KEY' '=' STRING)? 'API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')' 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    | 'CREATE' 'API' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'API_PROVIDER' '=' ('aws_api_gateway' | 'aws_private_api_gateway' | 'aws_gov_api_gateway' | 'aws_gov_private_api_gateway') 'API_AWS_ROLE_ARN' '=' STRING ('API_KEY' '=' STRING)? 'API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')' 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    | 'CREATE' ('OR' 'REPLACE')? 'API' 'INTEGRATION' ID 'API_PROVIDER' '=' 'azure_api_management' 'AZURE_TENANT_ID' '=' STRING 'AZURE_AD_APPLICATION_ID' '=' STRING ('API_KEY' '=' STRING)? 'API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')' ('API_BLOCKED_PREFIXES' '=' '(' STRING (',' STRING)* ')')? 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    | 'CREATE' 'API' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'API_PROVIDER' '=' 'azure_api_management' 'AZURE_TENANT_ID' '=' STRING 'AZURE_AD_APPLICATION_ID' '=' STRING ('API_KEY' '=' STRING)? 'API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')' ('API_BLOCKED_PREFIXES' '=' '(' STRING (',' STRING)* ')')? 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    | 'CREATE' ('OR' 'REPLACE')? 'API' 'INTEGRATION' ID 'API_PROVIDER' '=' 'google_api_gateway' 'GOOGLE_AUDIENCE' '=' STRING 'API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')' ('API_BLOCKED_PREFIXES' '=' '(' STRING (',' STRING)* ')')? 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    | 'CREATE' 'API' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'API_PROVIDER' '=' 'google_api_gateway' 'GOOGLE_AUDIENCE' '=' STRING 'API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')' ('API_BLOCKED_PREFIXES' '=' '(' STRING (',' STRING)* ')')? 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    | 'CREATE' ('OR' 'REPLACE')? 'API' 'INTEGRATION' ID 'API_PROVIDER' '=' 'git_https_api' 'API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')' ('API_BLOCKED_PREFIXES' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_AUTHENTICATION_SECRETS' '=' '(' QNAME (',' QNAME)* ')' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'all' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'none')? 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    | 'CREATE' 'API' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'API_PROVIDER' '=' 'git_https_api' 'API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')' ('API_BLOCKED_PREFIXES' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_AUTHENTICATION_SECRETS' '=' '(' QNAME (',' QNAME)* ')' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'all' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'none')? 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    | 'CREATE' ('OR' 'REPLACE')? 'API' 'INTEGRATION' ID 'API_PROVIDER' '=' 'git_https_api' 'API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')' ('API_BLOCKED_PREFIXES' '=' '(' STRING (',' STRING)* ')')? 'API_USER_AUTHENTICATION' '=' '(' 'TYPE' '=' 'SNOWFLAKE_GITHUB_APP' ')' 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    | 'CREATE' 'API' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'API_PROVIDER' '=' 'git_https_api' 'API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')' ('API_BLOCKED_PREFIXES' '=' '(' STRING (',' STRING)* ')')? 'API_USER_AUTHENTICATION' '=' '(' 'TYPE' '=' 'SNOWFLAKE_GITHUB_APP' ')' 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    | 'CREATE' ('OR' 'REPLACE')? 'API' 'INTEGRATION' ID 'API_PROVIDER' '=' 'git_https_api' 'API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')' ('API_BLOCKED_PREFIXES' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_AUTHENTICATION_SECRETS' '=' '(' QNAME (',' QNAME)* ')' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'all' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'none')? 'USE_PRIVATELINK_ENDPOINT' '=' ('TRUE' | 'FALSE') ('TLS_TRUSTED_CERTIFICATES' '=' '(' QNAME (',' QNAME)* ')')? 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    | 'CREATE' 'API' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'API_PROVIDER' '=' 'git_https_api' 'API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')' ('API_BLOCKED_PREFIXES' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_AUTHENTICATION_SECRETS' '=' '(' QNAME (',' QNAME)* ')' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'all' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'none')? 'USE_PRIVATELINK_ENDPOINT' '=' ('TRUE' | 'FALSE') ('TLS_TRUSTED_CERTIFICATES' '=' '(' QNAME (',' QNAME)* ')')? 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    | 'CREATE' ('OR' 'REPLACE')? 'API' 'INTEGRATION' ID 'API_PROVIDER' '=' 'artifact_repository_api' 'API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')' ('API_BLOCKED_PREFIXES' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_AUTHENTICATION_SECRETS' '=' '(' QNAME (',' QNAME)* ')' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'all' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'none')? ('USE_PRIVATELINK_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('TLS_TRUSTED_CERTIFICATES' '=' '(' QNAME (',' QNAME)* ')')? 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    | 'CREATE' 'API' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'API_PROVIDER' '=' 'artifact_repository_api' 'API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')' ('API_BLOCKED_PREFIXES' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_AUTHENTICATION_SECRETS' '=' '(' QNAME (',' QNAME)* ')' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'all' | 'ALLOWED_AUTHENTICATION_SECRETS' '=' 'none')? ('USE_PRIVATELINK_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('TLS_TRUSTED_CERTIFICATES' '=' '(' QNAME (',' QNAME)* ')')? 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    | 'CREATE' ('OR' 'REPLACE')? 'API' 'INTEGRATION' ID 'API_PROVIDER' '=' 'external_mcp' 'API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')' 'API_USER_AUTHENTICATION' '=' '(' 'TYPE' '=' 'OAUTH2' 'OAUTH_CLIENT_ID' '=' STRING 'OAUTH_CLIENT_SECRET' '=' STRING 'OAUTH_TOKEN_ENDPOINT' '=' STRING 'OAUTH_AUTHORIZATION_ENDPOINT' '=' STRING ('OAUTH_CLIENT_AUTH_METHOD' '=' ('CLIENT_SECRET_BASIC' | 'CLIENT_SECRET_POST' | 'NONE'))? ('OAUTH_DISCOVERY_URL' '=' STRING)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_RESOURCE' '=' STRING)? ')' 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    | 'CREATE' 'API' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'API_PROVIDER' '=' 'external_mcp' 'API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')' 'API_USER_AUTHENTICATION' '=' '(' 'TYPE' '=' 'OAUTH2' 'OAUTH_CLIENT_ID' '=' STRING 'OAUTH_CLIENT_SECRET' '=' STRING 'OAUTH_TOKEN_ENDPOINT' '=' STRING 'OAUTH_AUTHORIZATION_ENDPOINT' '=' STRING ('OAUTH_CLIENT_AUTH_METHOD' '=' ('CLIENT_SECRET_BASIC' | 'CLIENT_SECRET_POST' | 'NONE'))? ('OAUTH_DISCOVERY_URL' '=' STRING)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_RESOURCE' '=' STRING)? ')' 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    | 'CREATE' ('OR' 'REPLACE')? 'API' 'INTEGRATION' ID 'API_PROVIDER' '=' 'external_mcp' 'API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')' 'API_USER_AUTHENTICATION' '=' '(' 'TYPE' '=' 'OAUTH_DYNAMIC_CLIENT' 'OAUTH_RESOURCE_URL' '=' STRING ')' 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    | 'CREATE' 'API' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'API_PROVIDER' '=' 'external_mcp' 'API_ALLOWED_PREFIXES' '=' '(' STRING (',' STRING)* ')' 'API_USER_AUTHENTICATION' '=' '(' 'TYPE' '=' 'OAUTH_DYNAMIC_CLIENT' 'OAUTH_RESOURCE_URL' '=' STRING ')' 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-application
createApplication
    : 'CREATE' 'APPLICATION' ID 'FROM' 'APPLICATION' 'PACKAGE' ID ('USING' 'RELEASE' 'CHANNEL' ('QA' | 'ALPHA' | 'DEFAULT'))? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('AUTHORIZE_TELEMETRY_EVENT_SHARING' '=' ('TRUE' | 'FALSE'))? ('AUTHORIZE_RESTRICTED_PROVIDER_REMOTE_OPERATIONS_UNTIL' '=' (STRING | '\'INDEFINITE\'' | '\'NEVER\''))? ('WITH' 'FEATURE' 'POLICY' '=' QNAME)?
    | 'CREATE' 'APPLICATION' ID 'FROM' 'APPLICATION' 'PACKAGE' ID 'USING' STRING ('DEBUG_MODE' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('AUTHORIZE_TELEMETRY_EVENT_SHARING' '=' ('TRUE' | 'FALSE'))? ('AUTHORIZE_RESTRICTED_PROVIDER_REMOTE_OPERATIONS_UNTIL' '=' (STRING | '\'INDEFINITE\'' | '\'NEVER\''))? ('WITH' 'FEATURE' 'POLICY' '=' QNAME)?
    | 'CREATE' 'APPLICATION' ID 'FROM' 'APPLICATION' 'PACKAGE' ID 'USING' 'VERSION' ID ('PATCH' NUMBER)? ('DEBUG_MODE' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('AUTHORIZE_TELEMETRY_EVENT_SHARING' '=' ('TRUE' | 'FALSE'))? ('AUTHORIZE_RESTRICTED_PROVIDER_REMOTE_OPERATIONS_UNTIL' '=' (STRING | '\'INDEFINITE\'' | '\'NEVER\''))? ('WITH' 'FEATURE' 'POLICY' '=' QNAME)?
    | 'CREATE' 'APPLICATION' ID 'FROM' 'LISTING' ID ('USING' 'RELEASE' 'CHANNEL' ('QA' | 'ALPHA' | 'DEFAULT'))? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('BACKGROUND_INSTALL' '=' ('TRUE' | 'FALSE'))? ('AUTHORIZE_TELEMETRY_EVENT_SHARING' '=' ('TRUE' | 'FALSE'))? ('AUTHORIZE_RESTRICTED_PROVIDER_REMOTE_OPERATIONS_UNTIL' '=' (STRING | '\'INDEFINITE\'' | '\'NEVER\''))? ('WITH' 'FEATURE' 'POLICY' '=' QNAME)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-application-package
createApplicationPackage
    : 'CREATE' 'APPLICATION' 'PACKAGE' ('IF' 'NOT' 'EXISTS')? ID ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('COMMENT' '=' STRING)? ('DISTRIBUTION' '=' ('INTERNAL' | 'EXTERNAL'))? ('LISTING_AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('MULTIPLE_INSTANCES' '=' 'TRUE')? ('ENABLE_RELEASE_CHANNELS' '=' ('TRUE' | 'FALSE'))? ('AUTOMATIC_APPLICATION_MAINTENANCE' '=' ('TRUE' | 'FALSE'))? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-application-role
createApplicationRole
    : 'CREATE' ('OR' 'REPLACE')? 'APPLICATION' 'ROLE' QNAME ('COMMENT' '=' STRING)?
    | 'CREATE' 'APPLICATION' 'ROLE' ('IF' 'NOT' 'EXISTS')? QNAME ('COMMENT' '=' STRING)?
    | 'CREATE' 'OR' 'ALTER' 'APPLICATION' 'ROLE' QNAME ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-authentication-policy
// fix: MFA_POLICY: the parameters list the properties, space-delimited
// fix: PAT_POLICY: the parameters list the properties, space-delimited
// fix: WORKLOAD_IDENTITY_POLICY: the parameters list the properties, space-delimited
createAuthenticationPolicy
    : 'CREATE' ('OR' 'REPLACE')? 'AUTHENTICATION' 'POLICY' QNAME ('AUTHENTICATION_METHODS' '=' '(' STRING (',' STRING)* ')')? ('CLIENT_TYPES' '=' '(' STRING (',' STRING)* ')')? ('CLIENT_POLICY' '=' '(' ID '=' '(' 'MINIMUM_VERSION' '=' STRING ')' (',' ID '=' '(' 'MINIMUM_VERSION' '=' STRING ')')* ')')? ('SECURITY_INTEGRATIONS' '=' '(' STRING (',' STRING)* ')')? ('MFA_ENROLLMENT' '=' ('\'REQUIRED\'' | '\'REQUIRED_PASSWORD_ONLY\'' | '\'OPTIONAL\''))? ('MFA_POLICY' '=' '(' ('ALLOWED_METHODS' '=' '(' ('\'ALL\'' | '\'PASSKEY\'' | '\'TOTP\'' | '\'OTP\'' | '\'DUO\'') (',' ('\'ALL\'' | '\'PASSKEY\'' | '\'TOTP\'' | '\'OTP\'' | '\'DUO\''))* ')' | 'ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION' '=' ('\'ALL\'' | '\'NONE\''))+ ')')? ('PAT_POLICY' '=' '(' ('DEFAULT_EXPIRY_IN_DAYS' '=' NUMBER | 'MAX_EXPIRY_IN_DAYS' '=' NUMBER | 'NETWORK_POLICY_EVALUATION' '=' ('ENFORCED_REQUIRED' | 'ENFORCED_NOT_REQUIRED' | 'NOT_ENFORCED') | 'REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS' '=' ('TRUE' | 'FALSE'))+ ')')? ('WORKLOAD_IDENTITY_POLICY' '=' '(' ('ALLOWED_PROVIDERS' '=' '(' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC') (',' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC'))* ')' | 'ALLOWED_AWS_ACCOUNTS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_AZURE_ISSUERS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_OIDC_ISSUERS' '=' '(' STRING (',' STRING)* ')')+ ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' 'AUTHENTICATION' 'POLICY' ('IF' 'NOT' 'EXISTS')? QNAME ('AUTHENTICATION_METHODS' '=' '(' STRING (',' STRING)* ')')? ('CLIENT_TYPES' '=' '(' STRING (',' STRING)* ')')? ('CLIENT_POLICY' '=' '(' ID '=' '(' 'MINIMUM_VERSION' '=' STRING ')' (',' ID '=' '(' 'MINIMUM_VERSION' '=' STRING ')')* ')')? ('SECURITY_INTEGRATIONS' '=' '(' STRING (',' STRING)* ')')? ('MFA_ENROLLMENT' '=' ('\'REQUIRED\'' | '\'REQUIRED_PASSWORD_ONLY\'' | '\'OPTIONAL\''))? ('MFA_POLICY' '=' '(' ('ALLOWED_METHODS' '=' '(' ('\'ALL\'' | '\'PASSKEY\'' | '\'TOTP\'' | '\'OTP\'' | '\'DUO\'') (',' ('\'ALL\'' | '\'PASSKEY\'' | '\'TOTP\'' | '\'OTP\'' | '\'DUO\''))* ')' | 'ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION' '=' ('\'ALL\'' | '\'NONE\''))+ ')')? ('PAT_POLICY' '=' '(' ('DEFAULT_EXPIRY_IN_DAYS' '=' NUMBER | 'MAX_EXPIRY_IN_DAYS' '=' NUMBER | 'NETWORK_POLICY_EVALUATION' '=' ('ENFORCED_REQUIRED' | 'ENFORCED_NOT_REQUIRED' | 'NOT_ENFORCED') | 'REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS' '=' ('TRUE' | 'FALSE'))+ ')')? ('WORKLOAD_IDENTITY_POLICY' '=' '(' ('ALLOWED_PROVIDERS' '=' '(' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC') (',' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC'))* ')' | 'ALLOWED_AWS_ACCOUNTS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_AZURE_ISSUERS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_OIDC_ISSUERS' '=' '(' STRING (',' STRING)* ')')+ ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' 'OR' 'ALTER' 'AUTHENTICATION' 'POLICY' QNAME ('AUTHENTICATION_METHODS' '=' '(' STRING (',' STRING)* ')')? ('CLIENT_TYPES' '=' '(' STRING (',' STRING)* ')')? ('CLIENT_POLICY' '=' '(' ID '=' '(' 'MINIMUM_VERSION' '=' STRING ')' (',' ID '=' '(' 'MINIMUM_VERSION' '=' STRING ')')* ')')? ('SECURITY_INTEGRATIONS' '=' '(' STRING (',' STRING)* ')')? ('MFA_ENROLLMENT' '=' ('\'REQUIRED\'' | '\'REQUIRED_PASSWORD_ONLY\'' | '\'OPTIONAL\''))? ('MFA_POLICY' '=' '(' ('ALLOWED_METHODS' '=' '(' ('\'ALL\'' | '\'PASSKEY\'' | '\'TOTP\'' | '\'OTP\'' | '\'DUO\'') (',' ('\'ALL\'' | '\'PASSKEY\'' | '\'TOTP\'' | '\'OTP\'' | '\'DUO\''))* ')' | 'ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION' '=' ('\'ALL\'' | '\'NONE\''))+ ')')? ('PAT_POLICY' '=' '(' ('DEFAULT_EXPIRY_IN_DAYS' '=' NUMBER | 'MAX_EXPIRY_IN_DAYS' '=' NUMBER | 'NETWORK_POLICY_EVALUATION' '=' ('ENFORCED_REQUIRED' | 'ENFORCED_NOT_REQUIRED' | 'NOT_ENFORCED') | 'REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS' '=' ('TRUE' | 'FALSE'))+ ')')? ('WORKLOAD_IDENTITY_POLICY' '=' '(' ('ALLOWED_PROVIDERS' '=' '(' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC') (',' ('ALL' | 'AWS' | 'AZURE' | 'GCP' | 'OIDC'))* ')' | 'ALLOWED_AWS_ACCOUNTS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_AZURE_ISSUERS' '=' '(' STRING (',' STRING)* ')' | 'ALLOWED_OIDC_ISSUERS' '=' '(' STRING (',' STRING)* ')')+ ')')? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-backup-policy
createBackupPolicy
    : 'CREATE' ('OR' 'REPLACE')? 'BACKUP' 'POLICY' QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'RETENTION' 'LOCK')? ('SCHEDULE' '=' STRING)? ('EXPIRE_AFTER_DAYS' '=' NUMBER)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'BACKUP' 'POLICY' ('IF' 'NOT' 'EXISTS')? QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'RETENTION' 'LOCK')? ('SCHEDULE' '=' STRING)? ('EXPIRE_AFTER_DAYS' '=' NUMBER)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'OR' 'ALTER' 'BACKUP' 'POLICY' QNAME ('WITH' 'RETENTION' 'LOCK')? ('SCHEDULE' '=' STRING)? ('EXPIRE_AFTER_DAYS' '=' NUMBER)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-backup-set
createBackupSet
    : 'CREATE' ('OR' 'REPLACE')? 'BACKUP' 'SET' QNAME 'FOR' 'DYNAMIC'? 'ICEBERG'? 'TABLE' QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'BACKUP' 'POLICY' QNAME)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'BACKUP' 'SET' ('IF' 'NOT' 'EXISTS')? QNAME 'FOR' 'DYNAMIC'? 'ICEBERG'? 'TABLE' QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'BACKUP' 'POLICY' QNAME)? ('COMMENT' '=' STRING)?
    | 'CREATE' ('OR' 'REPLACE')? 'BACKUP' 'SET' QNAME 'FOR' 'SCHEMA' QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'BACKUP' 'POLICY' QNAME)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'BACKUP' 'SET' ('IF' 'NOT' 'EXISTS')? QNAME 'FOR' 'SCHEMA' QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'BACKUP' 'POLICY' QNAME)? ('COMMENT' '=' STRING)?
    | 'CREATE' ('OR' 'REPLACE')? 'BACKUP' 'SET' QNAME 'FOR' 'DATABASE' ID ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'BACKUP' 'POLICY' QNAME)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'BACKUP' 'SET' ('IF' 'NOT' 'EXISTS')? QNAME 'FOR' 'DATABASE' ID ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'BACKUP' 'POLICY' QNAME)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'OR' 'ALTER' 'BACKUP' 'SET' QNAME 'FOR' 'DYNAMIC'? 'ICEBERG'? 'TABLE' QNAME ('WITH' 'BACKUP' 'POLICY' QNAME)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'OR' 'ALTER' 'BACKUP' 'SET' QNAME 'FOR' 'SCHEMA' QNAME ('WITH' 'BACKUP' 'POLICY' QNAME)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'OR' 'ALTER' 'BACKUP' 'SET' QNAME 'FOR' 'DATABASE' ID ('WITH' 'BACKUP' 'POLICY' QNAME)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration
// no syntax

// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-delta-sharing
createCatalogIntegrationDeltaSharing
    : 'CREATE' ('OR' 'REPLACE')? 'CATALOG' 'INTEGRATION' ID 'CATALOG_SOURCE' '=' 'DELTA_SHARING' 'TABLE_FORMAT' '=' 'DELTA' 'REST_CONFIG' '=' '(' 'CATALOG_URI' '=' STRING 'CATALOG_NAME' '=' STRING 'ACCESS_DELEGATION_MODE' '=' 'VENDED_CREDENTIALS' ')' 'REST_AUTHENTICATION' '=' '(' createCatalogIntegrationDeltaSharing_restAuthenticationParams ')' 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    | 'CREATE' 'CATALOG' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'CATALOG_SOURCE' '=' 'DELTA_SHARING' 'TABLE_FORMAT' '=' 'DELTA' 'REST_CONFIG' '=' '(' 'CATALOG_URI' '=' STRING 'CATALOG_NAME' '=' STRING 'ACCESS_DELEGATION_MODE' '=' 'VENDED_CREDENTIALS' ')' 'REST_AUTHENTICATION' '=' '(' createCatalogIntegrationDeltaSharing_restAuthenticationParams ')' 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    ;
createCatalogIntegrationDeltaSharing_restAuthenticationParams
    : 'TYPE' '=' 'BEARER' 'BEARER_TOKEN' '=' STRING
    | 'TYPE' '=' 'OIDC' 'OIDC_AUDIENCE' '=' STRING
    | 'TYPE' '=' 'OAUTH' ('OAUTH_API_TYPE' '=' ('PUBLIC' | 'PRIVATE'))? 'OAUTH_CLIENT_ID' '=' STRING 'OAUTH_CLIENT_SECRET' '=' STRING 'OAUTH_TOKEN_URI' '=' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-delta
createIcebergTableDelta
    : 'CREATE' ('OR' 'REPLACE')? 'ICEBERG' 'TABLE' QNAME ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' STRING)? 'BASE_LOCATION' '=' STRING ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    | 'CREATE' 'ICEBERG' 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' STRING)? 'BASE_LOCATION' '=' STRING ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-glue
createCatalogIntegrationGlue
    : 'CREATE' ('OR' 'REPLACE')? 'CATALOG' 'INTEGRATION' ID 'CATALOG_SOURCE' '=' 'GLUE' 'TABLE_FORMAT' '=' 'ICEBERG' ('CATALOG_NAMESPACE' '=' STRING)? 'REST_CONFIG' '=' '(' 'CATALOG_API_TYPE' '=' ('AWS_GLUE' | 'AWS_PRIVATE_GLUE') 'CATALOG_NAME' '=' STRING ('ACCESS_DELEGATION_MODE' '=' ('VENDED_CREDENTIALS' | 'EXTERNAL_VOLUME_CREDENTIALS'))? ')' 'REST_AUTHENTICATION' '=' '(' 'TYPE' '=' 'SIGV4' 'SIGV4_IAM_ROLE' '=' STRING ('SIGV4_SIGNING_REGION' '=' STRING)? ')' 'ENABLED' '=' ('TRUE' | 'FALSE') ('REFRESH_INTERVAL_SECONDS' '=' NUMBER)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'CATALOG' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'CATALOG_SOURCE' '=' 'GLUE' 'TABLE_FORMAT' '=' 'ICEBERG' ('CATALOG_NAMESPACE' '=' STRING)? 'REST_CONFIG' '=' '(' 'CATALOG_API_TYPE' '=' ('AWS_GLUE' | 'AWS_PRIVATE_GLUE') 'CATALOG_NAME' '=' STRING ('ACCESS_DELEGATION_MODE' '=' ('VENDED_CREDENTIALS' | 'EXTERNAL_VOLUME_CREDENTIALS'))? ')' 'REST_AUTHENTICATION' '=' '(' 'TYPE' '=' 'SIGV4' 'SIGV4_IAM_ROLE' '=' STRING ('SIGV4_SIGNING_REGION' '=' STRING)? ')' 'ENABLED' '=' ('TRUE' | 'FALSE') ('REFRESH_INTERVAL_SECONDS' '=' NUMBER)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-object-storage
createCatalogIntegrationObjectStorage
    : 'CREATE' ('OR' 'REPLACE')? 'CATALOG' 'INTEGRATION' ID 'CATALOG_SOURCE' '=' 'OBJECT_STORE' 'TABLE_FORMAT' '=' ('ICEBERG' | 'DELTA' | 'NONE' | 'HIVE') 'ENABLED' '=' ('TRUE' | 'FALSE') ('REFRESH_INTERVAL_SECONDS' '=' NUMBER)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'CATALOG' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'CATALOG_SOURCE' '=' 'OBJECT_STORE' 'TABLE_FORMAT' '=' ('ICEBERG' | 'DELTA' | 'NONE' | 'HIVE') 'ENABLED' '=' ('TRUE' | 'FALSE') ('REFRESH_INTERVAL_SECONDS' '=' NUMBER)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-horizon
createCatalogIntegrationHorizon
    : 'CREATE' ('OR' 'REPLACE')? 'CATALOG' 'INTEGRATION' ID 'CATALOG_SOURCE' '=' 'POLARIS' 'TABLE_FORMAT' '=' 'ICEBERG' 'REST_CONFIG' '=' '(' 'CATALOG_URI' '=' STRING ('CATALOG_API_TYPE' '=' ('PUBLIC' | 'PRIVATE'))? 'CATALOG_NAME' '=' STRING ('ACCESS_DELEGATION_MODE' '=' ('VENDED_CREDENTIALS' | 'EXTERNAL_VOLUME_CREDENTIALS'))? ')' 'REST_AUTHENTICATION' '=' '(' 'TYPE' '=' 'OAUTH' 'OAUTH_GRANT_TYPE' '=' 'TOKEN_EXCHANGE' 'OAUTH_ALLOWED_SCOPES' '=' '(' STRING ')' ('OAUTH_AUDIENCE' '=' STRING)? ')' ('DEFAULT_STORAGE_CONFIG' '=' '(' 'USE_PRIVATELINK_ENDPOINT' '=' ('TRUE' | 'FALSE') ')')? 'ENABLED' '=' ('TRUE' | 'FALSE') ('REFRESH_INTERVAL_SECONDS' '=' NUMBER)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'CATALOG' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'CATALOG_SOURCE' '=' 'POLARIS' 'TABLE_FORMAT' '=' 'ICEBERG' 'REST_CONFIG' '=' '(' 'CATALOG_URI' '=' STRING ('CATALOG_API_TYPE' '=' ('PUBLIC' | 'PRIVATE'))? 'CATALOG_NAME' '=' STRING ('ACCESS_DELEGATION_MODE' '=' ('VENDED_CREDENTIALS' | 'EXTERNAL_VOLUME_CREDENTIALS'))? ')' 'REST_AUTHENTICATION' '=' '(' 'TYPE' '=' 'OAUTH' 'OAUTH_GRANT_TYPE' '=' 'TOKEN_EXCHANGE' 'OAUTH_ALLOWED_SCOPES' '=' '(' STRING ')' ('OAUTH_AUDIENCE' '=' STRING)? ')' ('DEFAULT_STORAGE_CONFIG' '=' '(' 'USE_PRIVATELINK_ENDPOINT' '=' ('TRUE' | 'FALSE') ')')? 'ENABLED' '=' ('TRUE' | 'FALSE') ('REFRESH_INTERVAL_SECONDS' '=' NUMBER)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-open-catalog
// fix: one or more scopes, as the parameters and examples say
createCatalogIntegrationOpenCatalog
    : 'CREATE' ('OR' 'REPLACE')? 'CATALOG' 'INTEGRATION' ID 'CATALOG_SOURCE' '=' 'POLARIS' 'TABLE_FORMAT' '=' 'ICEBERG' ('CATALOG_NAMESPACE' '=' STRING)? 'REST_CONFIG' '=' '(' 'CATALOG_URI' '=' STRING ('CATALOG_API_TYPE' '=' 'PUBLIC')? 'CATALOG_NAME' '=' STRING ('ACCESS_DELEGATION_MODE' '=' ('VENDED_CREDENTIALS' | 'EXTERNAL_VOLUME_CREDENTIALS'))? ')' 'REST_AUTHENTICATION' '=' '(' 'TYPE' '=' 'OAUTH' ('OAUTH_API_TYPE' '=' ('PUBLIC' | 'PRIVATE'))? ('OAUTH_TOKEN_URI' '=' STRING)? 'OAUTH_CLIENT_ID' '=' STRING 'OAUTH_CLIENT_SECRET' '=' STRING 'OAUTH_ALLOWED_SCOPES' '=' '(' STRING (',' STRING)* ')' ')' 'ENABLED' '=' ('TRUE' | 'FALSE') ('REFRESH_INTERVAL_SECONDS' '=' NUMBER)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'CATALOG' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'CATALOG_SOURCE' '=' 'POLARIS' 'TABLE_FORMAT' '=' 'ICEBERG' ('CATALOG_NAMESPACE' '=' STRING)? 'REST_CONFIG' '=' '(' 'CATALOG_URI' '=' STRING ('CATALOG_API_TYPE' '=' 'PUBLIC')? 'CATALOG_NAME' '=' STRING ('ACCESS_DELEGATION_MODE' '=' ('VENDED_CREDENTIALS' | 'EXTERNAL_VOLUME_CREDENTIALS'))? ')' 'REST_AUTHENTICATION' '=' '(' 'TYPE' '=' 'OAUTH' ('OAUTH_API_TYPE' '=' ('PUBLIC' | 'PRIVATE'))? ('OAUTH_TOKEN_URI' '=' STRING)? 'OAUTH_CLIENT_ID' '=' STRING 'OAUTH_CLIENT_SECRET' '=' STRING 'OAUTH_ALLOWED_SCOPES' '=' '(' STRING (',' STRING)* ')' ')' 'ENABLED' '=' ('TRUE' | 'FALSE') ('REFRESH_INTERVAL_SECONDS' '=' NUMBER)? ('COMMENT' '=' STRING)?
    | 'CREATE' ('OR' 'REPLACE')? 'CATALOG' 'INTEGRATION' ID 'CATALOG_SOURCE' '=' 'POLARIS' 'TABLE_FORMAT' '=' 'ICEBERG' ('CATALOG_NAMESPACE' '=' STRING)? 'REST_CONFIG' '=' '(' 'CATALOG_URI' '=' STRING 'CATALOG_API_TYPE' '=' 'PRIVATE' 'CATALOG_NAME' '=' STRING ('ACCESS_DELEGATION_MODE' '=' ('VENDED_CREDENTIALS' | 'EXTERNAL_VOLUME_CREDENTIALS'))? ')' 'REST_AUTHENTICATION' '=' '(' 'TYPE' '=' 'OAUTH' ('OAUTH_API_TYPE' '=' ('PUBLIC' | 'PRIVATE'))? 'OAUTH_CLIENT_ID' '=' STRING 'OAUTH_CLIENT_SECRET' '=' STRING 'OAUTH_ALLOWED_SCOPES' '=' '(' STRING (',' STRING)* ')' ')' 'ENABLED' '=' ('TRUE' | 'FALSE') ('REFRESH_INTERVAL_SECONDS' '=' NUMBER)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'CATALOG' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'CATALOG_SOURCE' '=' 'POLARIS' 'TABLE_FORMAT' '=' 'ICEBERG' ('CATALOG_NAMESPACE' '=' STRING)? 'REST_CONFIG' '=' '(' 'CATALOG_URI' '=' STRING 'CATALOG_API_TYPE' '=' 'PRIVATE' 'CATALOG_NAME' '=' STRING ('ACCESS_DELEGATION_MODE' '=' ('VENDED_CREDENTIALS' | 'EXTERNAL_VOLUME_CREDENTIALS'))? ')' 'REST_AUTHENTICATION' '=' '(' 'TYPE' '=' 'OAUTH' ('OAUTH_API_TYPE' '=' ('PUBLIC' | 'PRIVATE'))? 'OAUTH_CLIENT_ID' '=' STRING 'OAUTH_CLIENT_SECRET' '=' STRING 'OAUTH_ALLOWED_SCOPES' '=' '(' STRING (',' STRING)* ')' ')' 'ENABLED' '=' ('TRUE' | 'FALSE') ('REFRESH_INTERVAL_SECONDS' '=' NUMBER)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-rest
// fix: one or more scopes, as the example (catalog) shows
createCatalogIntegrationRest
    : 'CREATE' ('OR' 'REPLACE')? 'CATALOG' 'INTEGRATION' ID 'CATALOG_SOURCE' '=' 'ICEBERG_REST' 'TABLE_FORMAT' '=' 'ICEBERG' ('CATALOG_NAMESPACE' '=' STRING)? 'REST_CONFIG' '=' '(' createCatalogIntegrationRest_restConfigParams ')' 'REST_AUTHENTICATION' '=' '(' createCatalogIntegrationRest_restAuthenticationParams ')' 'ENABLED' '=' ('TRUE' | 'FALSE') ('REFRESH_INTERVAL_SECONDS' '=' NUMBER)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'CATALOG' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'CATALOG_SOURCE' '=' 'ICEBERG_REST' 'TABLE_FORMAT' '=' 'ICEBERG' ('CATALOG_NAMESPACE' '=' STRING)? 'REST_CONFIG' '=' '(' createCatalogIntegrationRest_restConfigParams ')' 'REST_AUTHENTICATION' '=' '(' createCatalogIntegrationRest_restAuthenticationParams ')' 'ENABLED' '=' ('TRUE' | 'FALSE') ('REFRESH_INTERVAL_SECONDS' '=' NUMBER)? ('COMMENT' '=' STRING)?
    ;
createCatalogIntegrationRest_restConfigParams
    : 'CATALOG_URI' '=' STRING ('PREFIX' '=' STRING)? ('CATALOG_NAME' '=' STRING)? ('CATALOG_API_TYPE' '=' ('PUBLIC' | 'PRIVATE' | 'AWS_API_GATEWAY' | 'AWS_PRIVATE_API_GATEWAY' | 'AWS_GLUE' | 'AWS_PRIVATE_GLUE' | 'AWS_S3TABLES' | 'AWS_PRIVATE_S3TABLES'))? ('ACCESS_DELEGATION_MODE' '=' ('VENDED_CREDENTIALS' | 'EXTERNAL_VOLUME_CREDENTIALS'))?
    ;
createCatalogIntegrationRest_restAuthenticationParams
    : 'TYPE' '=' 'OAUTH' ('OAUTH_API_TYPE' '=' ('PUBLIC' | 'PRIVATE'))? ('OAUTH_TOKEN_URI' '=' STRING)? 'OAUTH_CLIENT_ID' '=' STRING 'OAUTH_CLIENT_SECRET' '=' STRING 'OAUTH_ALLOWED_SCOPES' '=' '(' STRING (',' STRING)* ')'
    | 'TYPE' '=' 'BEARER' 'BEARER_TOKEN' '=' STRING
    | 'TYPE' '=' 'SIGV4' 'SIGV4_IAM_ROLE' '=' STRING ('SIGV4_SIGNING_REGION' '=' STRING)? ('SIGV4_EXTERNAL_ID' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-snowflake-postgres
createCatalogIntegrationSnowflakePostgres
    : 'CREATE' ('OR' 'REPLACE')? 'CATALOG' 'INTEGRATION' ID 'CATALOG_SOURCE' '=' 'SNOWFLAKE_POSTGRES' 'TABLE_FORMAT' '=' 'ICEBERG' ('CATALOG_NAMESPACE' '=' STRING)? 'REST_CONFIG' '=' '(' createCatalogIntegrationSnowflakePostgres_restConfigParams ')' 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    | 'CREATE' 'CATALOG' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'CATALOG_SOURCE' '=' 'SNOWFLAKE_POSTGRES' 'TABLE_FORMAT' '=' 'ICEBERG' ('CATALOG_NAMESPACE' '=' STRING)? 'REST_CONFIG' '=' '(' createCatalogIntegrationSnowflakePostgres_restConfigParams ')' 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    ;
createCatalogIntegrationSnowflakePostgres_restConfigParams
    : 'POSTGRES_INSTANCE' '=' STRING 'ACCESS_DELEGATION_MODE' '=' 'VENDED_CREDENTIALS' ('CATALOG_NAME' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-clone
// fix: read: Time Travel's OFFSET is a negative number of seconds, -N, and its STATEMENT a quoted query ID, as the Time Travel parameters say
// not converted (form): ... stands for text the docs leave out
//   CREATE [ OR REPLACE ] { DATABASE | SCHEMA } [ IF NOT EXISTS ] <object_name>
//     CLONE <source_object_name>
//       [ { AT | BEFORE } ( { TIMESTAMP => <timestamp> | OFFSET => -<seconds> | STATEMENT => '<query_id>' } ) ]
//       [ IGNORE TABLES WITH INSUFFICIENT DATA RETENTION ]
//       [ IGNORE HYBRID TABLES ]
//       [ INCLUDE INTERNAL STAGES ]
//     ...
// not converted (form): ... stands for text the docs leave out
//   CREATE [ OR REPLACE ] TABLE [ IF NOT EXISTS ] <object_name>
//     CLONE <source_object_name>
//       [ { AT | BEFORE } ( { TIMESTAMP => <timestamp> | OFFSET => -<seconds> | STATEMENT => '<query_id>' } ) ]
//     ...
// not converted (form): ... stands for text the docs leave out
//   CREATE [ OR REPLACE ] ICEBERG TABLE [ IF NOT EXISTS ] <name>
//     CLONE <source_iceberg_table>
//       [ { AT | BEFORE } ( { TIMESTAMP => <timestamp> | OFFSET => -<seconds> | STATEMENT => '<query_id>' } ) ]
//       [ COPY GRANTS ]
//       [ COPY = ( { COMMENTS | TAGS | GRANTS } [ , ... ] ) ]
//       ...
// not converted (form): ... stands for text the docs leave out
//   CREATE [ OR REPLACE ] { ALERT | FILE FORMAT | SEQUENCE | STAGE | STREAM | TASK }
//     [ IF NOT EXISTS ] <object_name>
//     CLONE <source_object_name>
//     ...
createClone
    : 'CREATE' ('OR' 'REPLACE')? 'DYNAMIC' 'TABLE' QNAME 'CLONE' QNAME (('AT' | 'BEFORE') '(' ('TIMESTAMP' '=>' STRING | 'OFFSET' '=>' '-' NUMBER | 'STATEMENT' '=>' STRING) ')')? ('TARGET_LAG' '=' (STRING | 'DOWNSTREAM') 'WAREHOUSE' '=' ID)?
    | 'CREATE' ('OR' 'REPLACE')? 'EVENT' 'TABLE' QNAME 'CLONE' QNAME (('AT' | 'BEFORE') '(' ('TIMESTAMP' '=>' STRING | 'OFFSET' '=>' '-' NUMBER | 'STATEMENT' '=>' STRING) ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'DATABASE' 'ROLE' ID 'CLONE' ID
    | 'CREATE' 'DATABASE' 'ROLE' ('IF' 'NOT' 'EXISTS')? ID 'CLONE' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-compute-pool
createComputePool
    : 'CREATE' 'COMPUTE' 'POOL' ('IF' 'NOT' 'EXISTS')? ID ('FOR' 'APPLICATION' ID)? 'MIN_NODES' '=' NUMBER 'MAX_NODES' '=' NUMBER 'INSTANCE_FAMILY' '=' ID ('AUTO_RESUME' '=' ('TRUE' | 'FALSE'))? ('INITIALLY_SUSPENDED' '=' ('TRUE' | 'FALSE'))? ('AUTO_SUSPEND_SECS' '=' NUMBER)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' '=' STRING)? ('PLACEMENT_GROUP' '=' STRING)? ('BACKUP_INSTANCE_FAMILIES' '=' '(' STRING (',' STRING)* ')')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-connection
createConnection
    : 'CREATE' 'CONNECTION' ('IF' 'NOT' 'EXISTS')? ID ('COMMENT' '=' STRING)?
    | 'CREATE' 'CONNECTION' ('IF' 'NOT' 'EXISTS')? ID 'AS' 'REPLICA' 'OF' ID '.' ID '.' ID ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-contact
createContact
    : 'CREATE' ('OR' 'REPLACE')? 'CONTACT' QNAME (('USERS' '=' '(' STRING (',' STRING)* ')' | 'EMAIL_DISTRIBUTION_LIST' '=' STRING | 'URL' '=' STRING))? ('COMMENT' '=' STRING)?
    | 'CREATE' 'CONTACT' ('IF' 'NOT' 'EXISTS')? QNAME (('USERS' '=' '(' STRING (',' STRING)* ')' | 'EMAIL_DISTRIBUTION_LIST' '=' STRING | 'URL' '=' STRING))? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-cortex-search
// fix: the model is a quoted name, as the example writes it
// fix: the parameters list the column specifications: a text column with a model, or a vector column with a query_model
createCortexSearch
    : 'CREATE' ('OR' 'REPLACE')? 'CORTEX' 'SEARCH' 'SERVICE' QNAME 'ON' ID ('PRIMARY' 'KEY' '(' ID (',' ID)* ')')? 'ATTRIBUTES' ID (',' ID)* 'WAREHOUSE' '=' ID 'TARGET_LAG' '=' STRING ('EMBEDDING_MODEL' '=' STRING)? ('REFRESH_MODE' '=' ('FULL' | 'INCREMENTAL'))? ('INITIALIZE' '=' ('ON_CREATE' | 'ON_SCHEDULE'))? ('FULL_INDEX_BUILD_INTERVAL_DAYS' '=' NUMBER)? ('REQUEST_LOGGING' '=' ('TRUE' | 'FALSE'))? ('AUTO_SUSPEND' '=' NUMBER)? ('COMMENT' '=' STRING)? 'AS' QUERY
    | 'CREATE' 'CORTEX' 'SEARCH' 'SERVICE' ('IF' 'NOT' 'EXISTS')? QNAME 'ON' ID ('PRIMARY' 'KEY' '(' ID (',' ID)* ')')? 'ATTRIBUTES' ID (',' ID)* 'WAREHOUSE' '=' ID 'TARGET_LAG' '=' STRING ('EMBEDDING_MODEL' '=' STRING)? ('REFRESH_MODE' '=' ('FULL' | 'INCREMENTAL'))? ('INITIALIZE' '=' ('ON_CREATE' | 'ON_SCHEDULE'))? ('FULL_INDEX_BUILD_INTERVAL_DAYS' '=' NUMBER)? ('REQUEST_LOGGING' '=' ('TRUE' | 'FALSE'))? ('AUTO_SUSPEND' '=' NUMBER)? ('COMMENT' '=' STRING)? 'AS' QUERY
    | 'CREATE' ('OR' 'REPLACE')? 'CORTEX' 'SEARCH' 'SERVICE' QNAME 'TEXT' 'INDEXES' ID (',' ID)* 'VECTOR' 'INDEXES' (ID ('(' 'model' '=' STRING ')')? | ID ('(' 'query_model' '=' STRING ')')?) (',' (ID ('(' 'model' '=' STRING ')')? | ID ('(' 'query_model' '=' STRING ')')?))* ('PRIMARY' 'KEY' '(' ID (',' ID)* ')')? 'ATTRIBUTES' ID (',' ID)* 'WAREHOUSE' '=' ID 'TARGET_LAG' '=' STRING ('REFRESH_MODE' '=' ('FULL' | 'INCREMENTAL'))? ('INITIALIZE' '=' ('ON_CREATE' | 'ON_SCHEDULE'))? ('FULL_INDEX_BUILD_INTERVAL_DAYS' '=' NUMBER)? ('REQUEST_LOGGING' '=' ('TRUE' | 'FALSE'))? ('AUTO_SUSPEND' '=' NUMBER)? ('COMMENT' '=' STRING)? 'AS' QUERY
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-data-metric-function
// not converted (form): ... stands for text the docs leave out
//   CREATE [ OR ALTER ] DATA METRIC FUNCTION ...
createDataMetricFunction
    : 'CREATE' ('OR' 'REPLACE')? 'SECURE'? 'DATA' 'METRIC' 'FUNCTION' QNAME '(' ID 'TABLE' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')' (',' ID 'TABLE' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')')? ')' 'RETURNS' 'NUMBER' ('NOT'? 'NULL')? ('LANGUAGE' 'SQL')? ('COMMENT' '=' STRING)? 'AS' STRING
    | 'CREATE' 'SECURE'? 'DATA' 'METRIC' 'FUNCTION' ('IF' 'NOT' 'EXISTS')? QNAME '(' ID 'TABLE' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')' (',' ID 'TABLE' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')')? ')' 'RETURNS' 'NUMBER' ('NOT'? 'NULL')? ('LANGUAGE' 'SQL')? ('COMMENT' '=' STRING)? 'AS' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-data-movement-policy
createDataMovementPolicy
    : 'CREATE' ('OR' 'REPLACE')? 'DATA' 'MOVEMENT' 'POLICY' QNAME ('ENFORCE_RULES' '=' '(' ID (',' ID)* ')')? ('ALERT_RULES' '=' '(' ID (',' ID)* ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' 'DATA' 'MOVEMENT' 'POLICY' ('IF' 'NOT' 'EXISTS')? QNAME ('ENFORCE_RULES' '=' '(' ID (',' ID)* ')')? ('ALERT_RULES' '=' '(' ID (',' ID)* ')')? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-data-movement-rule
createDataMovementRule
    : 'CREATE' ('OR' 'REPLACE')? 'DATA' 'MOVEMENT' 'RULE' QNAME 'TYPE' '=' STRING 'MAX_ROWS' 'AS' '(' ')' 'RETURNS' 'INTEGER' '->' '(' EXPR ')' ('COMMENT' '=' STRING)?
    | 'CREATE' 'DATA' 'MOVEMENT' 'RULE' ('IF' 'NOT' 'EXISTS')? QNAME 'TYPE' '=' STRING 'MAX_ROWS' 'AS' '(' ')' 'RETURNS' 'INTEGER' '->' '(' EXPR ')' ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-database
// fix: read: Time Travel's OFFSET is a negative number of seconds, -N, and its STATEMENT a quoted query ID, as the Time Travel parameters say
createDatabase
    : 'CREATE' ('OR' 'REPLACE')? 'TRANSIENT'? 'DATABASE' ID ('CLONE' QNAME (('AT' | 'BEFORE') '(' ('TIMESTAMP' '=>' STRING | 'OFFSET' '=>' '-' NUMBER | 'STATEMENT' '=>' STRING) ')')? ('IGNORE' 'TABLES' 'WITH' 'INSUFFICIENT' 'DATA' 'RETENTION')? ('IGNORE' 'HYBRID' 'TABLES')?)? ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('STAGE_UNDROP_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('EXTERNAL_VOLUME' '=' ID)? ('CATALOG' '=' ID)? ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_VERSION_DEFAULT' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('CATALOG_SYNC_NAMESPACE_MODE' '=' ('NEST' | 'FLATTEN'))? ('CATALOG_SYNC_NAMESPACE_FLATTEN_DELIMITER' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('OBJECT_VISIBILITY' '=' (ID | 'PRIVILEGED'))? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)?
    | 'CREATE' 'TRANSIENT'? 'DATABASE' ('IF' 'NOT' 'EXISTS')? ID ('CLONE' QNAME (('AT' | 'BEFORE') '(' ('TIMESTAMP' '=>' STRING | 'OFFSET' '=>' '-' NUMBER | 'STATEMENT' '=>' STRING) ')')? ('IGNORE' 'TABLES' 'WITH' 'INSUFFICIENT' 'DATA' 'RETENTION')? ('IGNORE' 'HYBRID' 'TABLES')?)? ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('STAGE_UNDROP_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('EXTERNAL_VOLUME' '=' ID)? ('CATALOG' '=' ID)? ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_VERSION_DEFAULT' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('CATALOG_SYNC_NAMESPACE_MODE' '=' ('NEST' | 'FLATTEN'))? ('CATALOG_SYNC_NAMESPACE_FLATTEN_DELIMITER' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('OBJECT_VISIBILITY' '=' (ID | 'PRIVILEGED'))? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)?
    | 'CREATE' 'DATABASE' ID 'FROM' 'BACKUP' 'SET' ID 'IDENTIFIER' STRING
    | 'CREATE' 'DATABASE' ID 'FROM' 'LISTING' STRING
    | 'CREATE' 'DATABASE' ID 'FROM' 'SHARE' ID '.' ID
    | 'CREATE' 'DATABASE' ID 'AS' 'REPLICA' 'OF' ID '.' ID ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)?
    | 'CREATE' 'OR' 'ALTER' 'TRANSIENT'? 'DATABASE' ID ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('EXTERNAL_VOLUME' '=' ID)? ('CATALOG' '=' ID)? ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_VERSION_DEFAULT' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('LOG_LEVEL' '=' STRING)? ('METRIC_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('COMMENT' '=' STRING)? ('OBJECT_VISIBILITY' '=' (ID | 'PRIVILEGED'))?
    | 'CREATE' ('OR' 'REPLACE')? 'DATABASE' ID 'CLONE' ID
    | 'CREATE' 'DATABASE' ('IF' 'NOT' 'EXISTS')? ID 'CLONE' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-database-catalog-linked
// fix: no comma follows LINKED_CATALOG = ( ... ), as the page's examples show
// fix: catalogParams: a comma stands between parameters, before each optional one, as the page's examples write them
createDatabaseCatalogLinked
    : 'CREATE' 'DATABASE' QNAME 'LINKED_CATALOG' '=' '(' createDatabaseCatalogLinked_catalogParams ')' ('EXTERNAL_VOLUME' '=' STRING)? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('CATALOG_CASE_SENSITIVITY' '=' ('CASE_SENSITIVE' | 'CASE_INSENSITIVE'))? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    ;
createDatabaseCatalogLinked_catalogParams
    : 'CATALOG' '=' STRING (',' 'ALLOWED_NAMESPACES' '=' '(' STRING (',' STRING)* ')')? (',' 'BLOCKED_NAMESPACES' '=' '(' STRING (',' STRING)* ')')? (',' 'ALLOWED_WRITE_OPERATIONS' '=' ('NONE' | 'ALL'))? (',' 'NAMESPACE_MODE' '=' ('IGNORE_NESTED_NAMESPACE' | 'FLATTEN_NESTED_NAMESPACE'))? (',' 'NAMESPACE_FLATTEN_DELIMITER' '=' STRING)? (',' 'SYNC_INTERVAL_SECONDS' '=' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-database-role
createDatabaseRole
    : 'CREATE' ('OR' 'REPLACE')? 'DATABASE' 'ROLE' QNAME ('COMMENT' '=' STRING)?
    | 'CREATE' 'DATABASE' 'ROLE' ('IF' 'NOT' 'EXISTS')? QNAME ('COMMENT' '=' STRING)?
    | 'CREATE' 'OR' 'ALTER' 'DATABASE' 'ROLE' QNAME ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-dataset
createDataset
    : 'CREATE' ('OR' 'REPLACE')? 'DATASET' QNAME
    | 'CREATE' ('IF' 'NOT' 'EXISTS')? 'DATASET' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-dbt-project
// fix: the version is a quoted string, as the example and the ALTER and EXECUTE pages write it
createDbtProject
    : 'CREATE' ('OR' 'REPLACE')? 'DBT' 'PROJECT' QNAME 'FROM' STRING ('DBT_VERSION' '=' STRING)? ('DEFAULT_TARGET' '=' ID)? ('DEFAULT_ENVIRONMENT' '=' STRING)? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('AUTO_COMPILE' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_WRITEBACK' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)?
    | 'CREATE' 'DBT' 'PROJECT' ('IF' 'NOT' 'EXISTS')? QNAME 'FROM' STRING ('DBT_VERSION' '=' STRING)? ('DEFAULT_TARGET' '=' ID)? ('DEFAULT_ENVIRONMENT' '=' STRING)? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('AUTO_COMPILE' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_WRITEBACK' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-dcm-project
createDcmProject
    : 'CREATE' ('OR' 'REPLACE')? 'DCM' 'PROJECT' QNAME ('LOG_LEVEL' '=' ('DEBUG' | 'INFO' | 'WARN' | 'ERROR'))? ('COMMENT' '=' STRING)?
    | 'CREATE' 'DCM' 'PROJECT' ('IF' 'NOT' 'EXISTS')? QNAME ('LOG_LEVEL' '=' ('DEBUG' | 'INFO' | 'WARN' | 'ERROR'))? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-dynamic-table
// fix: WITH CONTACT is the table's, after every clause but AS, as the parameters say, not a column's
// fix: WITH CONTACT goes last before AS, as the parameters say
// fix: read: Time Travel's OFFSET is a negative number of seconds, -N, and its STATEMENT a quoted query ID, as the Time Travel parameters say
// fix: read: the additional columns take the whole column definition, which their line abbreviates
createDynamicTable
    : 'CREATE' ('OR' 'REPLACE')? 'TRANSIENT'? 'DYNAMIC' 'TABLE' QNAME '(' ID DATA_TYPE ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)? (',' ID DATA_TYPE ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)?)* ')' 'TARGET_LAG' '=' (STRING | 'DOWNSTREAM') ('SCHEDULER' '=' 'DISABLE' | 'SCHEDULER' '=' 'ENABLE')? 'WAREHOUSE' '=' ID ('INITIALIZATION_WAREHOUSE' '=' ID)? ('REFRESH_MODE' '=' ('AUTO' | 'FULL' | 'INCREMENTAL' | 'ADAPTIVE' | 'CUSTOM_INCREMENTAL'))? ('INITIALIZE' '=' ('ON_CREATE' | 'ON_SCHEDULE'))? ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')')? ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('COMMENT' '=' STRING)? ('COPY' 'GRANTS')? ('COPY' 'TAGS')? ('COPY' '=' '(' ('COMMENTS' | 'TAGS' | 'GRANTS') (',' ('COMMENTS' | 'TAGS' | 'GRANTS'))* ')')? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'AGGREGATION' 'POLICY' QNAME ('ENTITY' 'KEY' '(' ID (',' ID)* ')')?)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('REQUIRE' 'USER')? ('FROZEN' 'WHERE' '(' EXPR ')')? ('WITH'? 'STORAGE' 'LIFECYCLE' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('BACKFILL' 'FROM' QNAME)? ('START' 'AT' '(' ('STREAM' '=>' STRING | 'TIMESTAMP' '=>' STRING | 'STATEMENT' '=>' STRING | 'OFFSET' '=>' '-' NUMBER) ')')? ('EXECUTE' 'AS' 'USER' ID ('USE' 'SECONDARY' 'ROLES' ('ALL' | 'NONE' | ID (',' ID)*))?)? ('ROW_TIMESTAMP' '=' ('TRUE' | 'FALSE'))? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('AS' QUERY | 'REFRESH' 'USING' '(' QUERY ')')
    | 'CREATE' 'TRANSIENT'? 'DYNAMIC' 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME '(' ID DATA_TYPE ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)? (',' ID DATA_TYPE ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)?)* ')' 'TARGET_LAG' '=' (STRING | 'DOWNSTREAM') ('SCHEDULER' '=' 'DISABLE' | 'SCHEDULER' '=' 'ENABLE')? 'WAREHOUSE' '=' ID ('INITIALIZATION_WAREHOUSE' '=' ID)? ('REFRESH_MODE' '=' ('AUTO' | 'FULL' | 'INCREMENTAL' | 'ADAPTIVE' | 'CUSTOM_INCREMENTAL'))? ('INITIALIZE' '=' ('ON_CREATE' | 'ON_SCHEDULE'))? ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')')? ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('COMMENT' '=' STRING)? ('COPY' 'GRANTS')? ('COPY' 'TAGS')? ('COPY' '=' '(' ('COMMENTS' | 'TAGS' | 'GRANTS') (',' ('COMMENTS' | 'TAGS' | 'GRANTS'))* ')')? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'AGGREGATION' 'POLICY' QNAME ('ENTITY' 'KEY' '(' ID (',' ID)* ')')?)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('REQUIRE' 'USER')? ('FROZEN' 'WHERE' '(' EXPR ')')? ('WITH'? 'STORAGE' 'LIFECYCLE' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('BACKFILL' 'FROM' QNAME)? ('START' 'AT' '(' ('STREAM' '=>' STRING | 'TIMESTAMP' '=>' STRING | 'STATEMENT' '=>' STRING | 'OFFSET' '=>' '-' NUMBER) ')')? ('EXECUTE' 'AS' 'USER' ID ('USE' 'SECONDARY' 'ROLES' ('ALL' | 'NONE' | ID (',' ID)*))?)? ('ROW_TIMESTAMP' '=' ('TRUE' | 'FALSE'))? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('AS' QUERY | 'REFRESH' 'USING' '(' QUERY ')')
    | 'CREATE' 'OR' 'ALTER' 'TRANSIENT'? 'DYNAMIC' 'TABLE' QNAME '(' ID DATA_TYPE ('COLLATE' STRING)? ('COMMENT' STRING)? (',' ID DATA_TYPE ('COLLATE' STRING)? ('COMMENT' STRING)?)* ')' 'TARGET_LAG' '=' (STRING | 'DOWNSTREAM') ('SCHEDULER' '=' 'DISABLE' | 'SCHEDULER' '=' 'ENABLE')? 'WAREHOUSE' '=' ID ('INITIALIZATION_WAREHOUSE' '=' ID)? ('REFRESH_MODE' '=' ('AUTO' | 'FULL' | 'INCREMENTAL' | 'ADAPTIVE' | 'CUSTOM_INCREMENTAL'))? ('INITIALIZE' '=' ('ON_CREATE' | 'ON_SCHEDULE'))? ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')')? ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('COMMENT' '=' STRING)? ('REQUIRE' 'USER')? ('FROZEN' 'WHERE' '(' EXPR ')')? ('BACKFILL' 'FROM' QNAME)? ('START' 'AT' '(' ('STREAM' '=>' STRING | 'TIMESTAMP' '=>' STRING | 'STATEMENT' '=>' STRING | 'OFFSET' '=>' '-' NUMBER) ')')? ('EXECUTE' 'AS' 'USER' ID ('USE' 'SECONDARY' 'ROLES' ('ALL' | 'NONE' | ID (',' ID)*))?)? ('ROW_TIMESTAMP' '=' ('TRUE' | 'FALSE'))? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('AS' QUERY | 'REFRESH' 'USING' '(' QUERY ')')
    | 'CREATE' 'DYNAMIC' 'TABLE' QNAME 'FROM' 'BACKUP' 'SET' ID 'IDENTIFIER' STRING
    | 'CREATE' 'DYNAMIC' 'ICEBERG' 'TABLE' QNAME 'FROM' 'BACKUP' 'SET' ID 'IDENTIFIER' STRING
    | 'CREATE' ('OR' 'REPLACE')? 'TRANSIENT'? 'DYNAMIC' 'TABLE' QNAME 'CLONE' QNAME (('AT' | 'BEFORE') '(' ('TIMESTAMP' '=>' STRING | 'OFFSET' '=>' '-' NUMBER | 'STATEMENT' '=>' STRING) ')')? ('COPY' 'GRANTS' 'COPY' 'TAGS' 'COPY' '=' '(' ('COMMENTS' | 'TAGS' | 'GRANTS') (',' ('COMMENTS' | 'TAGS' | 'GRANTS'))* ')' 'TARGET_LAG' '=' (STRING | 'DOWNSTREAM') 'WAREHOUSE' '=' ID 'EXECUTE' 'AS' 'USER' ID 'USE' 'SECONDARY' 'ROLES' ('ALL' | 'NONE' | ID (',' ID)*))?
    | 'CREATE' ('OR' 'REPLACE')? 'DYNAMIC' 'ICEBERG' 'TABLE' QNAME '(' ID DATA_TYPE ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)? (',' ID DATA_TYPE ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)?)* ')' 'TARGET_LAG' '=' (STRING | 'DOWNSTREAM') 'WAREHOUSE' '=' ID ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' '\'SNOWFLAKE\'')? ('BASE_LOCATION' '=' STRING)? ('TARGET_FILE_SIZE' '=' ('\'AUTO\'' | '\'16MB\'' | '\'32MB\'' | '\'64MB\'' | '\'128MB\''))? ('PARTITION' 'BY' '(' createDynamicTable_partitionExpression (',' createDynamicTable_partitionExpression)* ')')? ('PATH_LAYOUT' '=' ('FLAT' | 'HIERARCHICAL'))? ('ICEBERG_VERSION' '=' NUMBER)? ('REFRESH_MODE' '=' ('AUTO' | 'FULL' | 'INCREMENTAL' | 'ADAPTIVE' | 'CUSTOM_INCREMENTAL'))? ('FROZEN' 'WHERE' '(' EXPR ')')? ('INITIALIZE' '=' ('ON_CREATE' | 'ON_SCHEDULE'))? ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')')? ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('COMMENT' '=' STRING)? ('COPY' 'GRANTS')? ('COPY' 'TAGS')? ('COPY' '=' '(' ('COMMENTS' | 'TAGS' | 'GRANTS') (',' ('COMMENTS' | 'TAGS' | 'GRANTS'))* ')')? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'STORAGE' 'LIFECYCLE' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('REQUIRE' 'USER')? ('EXECUTE' 'AS' 'USER' ID ('USE' 'SECONDARY' 'ROLES' ('ALL' | 'NONE' | ID (',' ID)*))?)? ('AS' QUERY | 'REFRESH' 'USING' '(' QUERY ')')
    | 'CREATE' ('OR' 'REPLACE')? 'DYNAMIC' 'ICEBERG' 'TABLE' QNAME 'CLONE' QNAME (('AT' | 'BEFORE') '(' ('TIMESTAMP' '=>' STRING | 'OFFSET' '=>' '-' NUMBER | 'STATEMENT' '=>' STRING) ')')? ('COPY' 'GRANTS' 'COPY' 'TAGS' 'COPY' '=' '(' ('COMMENTS' | 'TAGS' | 'GRANTS') (',' ('COMMENTS' | 'TAGS' | 'GRANTS'))* ')' 'TARGET_LAG' '=' (STRING | 'DOWNSTREAM') 'WAREHOUSE' '=' ID 'EXECUTE' 'AS' 'USER' ID 'USE' 'SECONDARY' 'ROLES' ('ALL' | 'NONE' | ID (',' ID)*))?
    ;
createDynamicTable_partitionExpression
    : ID
    | 'BUCKET' '(' NUMBER ',' ID ')'
    | 'TRUNCATE' '(' NUMBER ',' ID ')'
    | 'YEAR' '(' ID ')'
    | 'MONTH' '(' ID ')'
    | 'DAY' '(' ID ')'
    | 'HOUR' '(' ID ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-event-routing-table
// not converted (form): ... stands for text the docs leave out
//   CREATE EVENT ROUTING TABLE <table_name>
//     WITH RULES
//       {rule name} = (REGION_GROUP={region group}, REGIONS=('{region1}', '{region2}', ...), DESTINATION_ACCOUNT = {organization}.{account_name})
//       ...

// https://docs.snowflake.com/en/sql-reference/sql/create-event-table
// fix: read: Time Travel's OFFSET is a negative number of seconds, -N, and its STATEMENT a quoted query ID, as the Time Travel parameters say
createEventTable
    : 'CREATE' ('OR' 'REPLACE')? 'EVENT' 'TABLE' QNAME ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')')? ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('CHANGE_TRACKING' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('COPY' 'GRANTS')? ('COPY' '=' '(' 'GRANTS' ')')? ('WITH'? 'COMMENT' '=' STRING)? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    | 'CREATE' 'EVENT' 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')')? ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('CHANGE_TRACKING' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('COPY' 'GRANTS')? ('COPY' '=' '(' 'GRANTS' ')')? ('WITH'? 'COMMENT' '=' STRING)? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'EVENT' 'TABLE' QNAME 'CLONE' QNAME (('AT' | 'BEFORE') '(' ('TIMESTAMP' '=>' STRING | 'OFFSET' '=>' '-' NUMBER | 'STATEMENT' '=>' STRING) ')')? ('COPY' 'GRANTS')?
    | 'CREATE' 'EVENT' 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME 'CLONE' QNAME (('AT' | 'BEFORE') '(' ('TIMESTAMP' '=>' STRING | 'OFFSET' '=>' '-' NUMBER | 'STATEMENT' '=>' STRING) ')')? ('COPY' 'GRANTS')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-experiment
createExperiment
    : 'CREATE' ('OR' 'REPLACE')? 'EXPERIMENT' QNAME
    | 'CREATE' 'EXPERIMENT' ('IF' 'NOT' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-external-agent
createExternalAgent
    : 'CREATE' ('OR' 'REPLACE')? 'EXTERNAL' 'AGENT' QNAME ('WITH' 'VERSION' ID)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'EXTERNAL' 'AGENT' ('IF' 'NOT' 'EXISTS')? QNAME ('WITH' 'VERSION' ID)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-external-consumer
createExternalConsumer
    : 'CREATE' ('OR' 'REPLACE')? 'EXTERNAL' 'CONSUMER' QNAME ('REGION' '=' STRING)? ('COMMENT' '=' STRING)? ('EMAIL' '=' STRING)?
    | 'CREATE' 'EXTERNAL' 'CONSUMER' ('IF' 'NOT' 'EXISTS')? QNAME ('REGION' '=' STRING)? ('COMMENT' '=' STRING)? ('EMAIL' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-external-access-integration
createExternalAccessIntegration
    : 'CREATE' ('OR' 'REPLACE')? 'EXTERNAL' 'ACCESS' 'INTEGRATION' ID 'ALLOWED_NETWORK_RULES' '=' '(' ID (',' ID)* ')' ('ALLOWED_API_AUTHENTICATION_INTEGRATIONS' '=' ('(' ID (',' ID)* ')' | 'none'))? ('ALLOWED_AUTHENTICATION_SECRETS' '=' ('(' QNAME (',' QNAME)* ')' | 'all' | 'none'))? 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-external-function
// not converted (form): ... stands for text the docs leave out
//   CREATE [ OR ALTER ] EXTERNAL FUNCTION ...
createExternalFunction
    : 'CREATE' ('OR' 'REPLACE')? 'SECURE'? 'EXTERNAL' 'FUNCTION' QNAME '(' (ID DATA_TYPE (',' ID DATA_TYPE)*)? ')' 'RETURNS' DATA_TYPE ('NOT'? 'NULL')? (('CALLED' 'ON' 'NULL' 'INPUT' | ('RETURNS' 'NULL' 'ON' 'NULL' 'INPUT' | 'STRICT')))? (('VOLATILE' | 'IMMUTABLE'))? ('COMMENT' '=' STRING)? 'API_INTEGRATION' '=' ID ('HEADERS' '=' '(' STRING '=' STRING (',' STRING '=' STRING)* ')')? ('CONTEXT_HEADERS' '=' '(' ID (',' ID)* ')')? ('MAX_BATCH_ROWS' '=' NUMBER)? ('COMPRESSION' '=' ID)? ('REQUEST_TRANSLATOR' '=' ID)? ('RESPONSE_TRANSLATOR' '=' ID)? 'AS' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-external-table
// fix: the ORC options' NULL_IF list is not closed; closed, as the CSV options write it
// fix: each column, with its constraint, follows a comma, as the page's examples write them; the whole list is optional, since a table may have only its VALUE column
// fix: the template is a parenthesized subquery, as the parameters and examples write it
// fix: read: <constraint_properties> is the properties CREATE TABLE's constraint page lists
// fix: read: each -- If TYPE = X section of formatTypeOptions is the options of TYPE = X alone, and TYPE = { ... } [ formatTypeOptions ] chooses a TYPE with its own options
// fix: read: a foreign key's REFERENCES takes MATCH and ON UPDATE/DELETE, the properties the constraint page gives foreign keys
createExternalTable
    : 'CREATE' ('OR' 'REPLACE')? 'EXTERNAL' 'TABLE' QNAME ('(' (ID DATA_TYPE 'AS' EXPR createExternalTable_inlineConstraint? (',' ID DATA_TYPE 'AS' EXPR createExternalTable_inlineConstraint?)*)? ')')? createExternalTable_cloudProviderParams ('PARTITION' 'BY' '(' ID (',' ID)* ')')? 'WITH'? 'LOCATION' '=' createExternalTable_externalStage ('REFRESH_ON_CREATE' '=' ('TRUE' | 'FALSE'))? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('PATTERN' '=' STRING)? 'FILE_FORMAT' '=' '(' ('FORMAT_NAME' '=' STRING | ('TYPE' '=' 'CSV' createExternalTable_formatTypeOptions_CSV? | 'TYPE' '=' 'JSON' createExternalTable_formatTypeOptions_JSON? | 'TYPE' '=' 'AVRO' createExternalTable_formatTypeOptions_AVRO? | 'TYPE' '=' 'ORC' createExternalTable_formatTypeOptions_ORC? | 'TYPE' '=' 'PARQUET' createExternalTable_formatTypeOptions_PARQUET?)) ')' ('AWS_SNS_TOPIC' '=' STRING)? ('COPY' 'GRANTS')? ('COPY' '=' '(' 'GRANTS' ')')? ('COMMENT' '=' STRING)? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' 'VALUE' ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    | 'CREATE' 'EXTERNAL' 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME ('(' (ID DATA_TYPE 'AS' EXPR createExternalTable_inlineConstraint? (',' ID DATA_TYPE 'AS' EXPR createExternalTable_inlineConstraint?)*)? ')')? createExternalTable_cloudProviderParams ('PARTITION' 'BY' '(' ID (',' ID)* ')')? 'WITH'? 'LOCATION' '=' createExternalTable_externalStage ('REFRESH_ON_CREATE' '=' ('TRUE' | 'FALSE'))? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('PATTERN' '=' STRING)? 'FILE_FORMAT' '=' '(' ('FORMAT_NAME' '=' STRING | ('TYPE' '=' 'CSV' createExternalTable_formatTypeOptions_CSV? | 'TYPE' '=' 'JSON' createExternalTable_formatTypeOptions_JSON? | 'TYPE' '=' 'AVRO' createExternalTable_formatTypeOptions_AVRO? | 'TYPE' '=' 'ORC' createExternalTable_formatTypeOptions_ORC? | 'TYPE' '=' 'PARQUET' createExternalTable_formatTypeOptions_PARQUET?)) ')' ('AWS_SNS_TOPIC' '=' STRING)? ('COPY' 'GRANTS')? ('COPY' '=' '(' 'GRANTS' ')')? ('COMMENT' '=' STRING)? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' 'VALUE' ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'EXTERNAL' 'TABLE' QNAME ('(' (ID DATA_TYPE 'AS' EXPR createExternalTable_inlineConstraint? (',' ID DATA_TYPE 'AS' EXPR createExternalTable_inlineConstraint?)*)? ')')? createExternalTable_cloudProviderParams ('PARTITION' 'BY' '(' ID (',' ID)* ')')? 'WITH'? 'LOCATION' '=' createExternalTable_externalStage 'PARTITION_TYPE' '=' 'USER_SPECIFIED' 'FILE_FORMAT' '=' '(' ('FORMAT_NAME' '=' STRING | ('TYPE' '=' 'CSV' createExternalTable_formatTypeOptions_CSV? | 'TYPE' '=' 'JSON' createExternalTable_formatTypeOptions_JSON? | 'TYPE' '=' 'AVRO' createExternalTable_formatTypeOptions_AVRO? | 'TYPE' '=' 'ORC' createExternalTable_formatTypeOptions_ORC? | 'TYPE' '=' 'PARQUET' createExternalTable_formatTypeOptions_PARQUET?)) ')' ('COPY' 'GRANTS')? ('COPY' '=' '(' 'GRANTS' ')')? ('COMMENT' '=' STRING)? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' 'VALUE' ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    | 'CREATE' 'EXTERNAL' 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME ('(' (ID DATA_TYPE 'AS' EXPR createExternalTable_inlineConstraint? (',' ID DATA_TYPE 'AS' EXPR createExternalTable_inlineConstraint?)*)? ')')? createExternalTable_cloudProviderParams ('PARTITION' 'BY' '(' ID (',' ID)* ')')? 'WITH'? 'LOCATION' '=' createExternalTable_externalStage 'PARTITION_TYPE' '=' 'USER_SPECIFIED' 'FILE_FORMAT' '=' '(' ('FORMAT_NAME' '=' STRING | ('TYPE' '=' 'CSV' createExternalTable_formatTypeOptions_CSV? | 'TYPE' '=' 'JSON' createExternalTable_formatTypeOptions_JSON? | 'TYPE' '=' 'AVRO' createExternalTable_formatTypeOptions_AVRO? | 'TYPE' '=' 'ORC' createExternalTable_formatTypeOptions_ORC? | 'TYPE' '=' 'PARQUET' createExternalTable_formatTypeOptions_PARQUET?)) ')' ('COPY' 'GRANTS')? ('COPY' '=' '(' 'GRANTS' ')')? ('COMMENT' '=' STRING)? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' 'VALUE' ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'EXTERNAL' 'TABLE' QNAME ('(' (ID DATA_TYPE 'AS' EXPR createExternalTable_inlineConstraint? (',' ID DATA_TYPE 'AS' EXPR createExternalTable_inlineConstraint?)*)? ')')? createExternalTable_cloudProviderParams ('PARTITION' 'BY' '(' ID (',' ID)* ')')? 'WITH'? 'LOCATION' '=' createExternalTable_externalStage 'PARTITION_TYPE' '=' 'USER_SPECIFIED' 'FILE_FORMAT' '=' '(' ('FORMAT_NAME' '=' STRING | ('TYPE' '=' 'CSV' createExternalTable_formatTypeOptions_CSV? | 'TYPE' '=' 'JSON' createExternalTable_formatTypeOptions_JSON? | 'TYPE' '=' 'AVRO' createExternalTable_formatTypeOptions_AVRO? | 'TYPE' '=' 'ORC' createExternalTable_formatTypeOptions_ORC? | 'TYPE' '=' 'PARQUET' createExternalTable_formatTypeOptions_PARQUET?)) ')' ('TABLE_FORMAT' '=' 'DELTA')? ('COPY' 'GRANTS')? ('COPY' '=' '(' 'GRANTS' ')')? ('COMMENT' '=' STRING)? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' 'VALUE' ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    | 'CREATE' 'EXTERNAL' 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME ('(' (ID DATA_TYPE 'AS' EXPR createExternalTable_inlineConstraint? (',' ID DATA_TYPE 'AS' EXPR createExternalTable_inlineConstraint?)*)? ')')? createExternalTable_cloudProviderParams ('PARTITION' 'BY' '(' ID (',' ID)* ')')? 'WITH'? 'LOCATION' '=' createExternalTable_externalStage 'PARTITION_TYPE' '=' 'USER_SPECIFIED' 'FILE_FORMAT' '=' '(' ('FORMAT_NAME' '=' STRING | ('TYPE' '=' 'CSV' createExternalTable_formatTypeOptions_CSV? | 'TYPE' '=' 'JSON' createExternalTable_formatTypeOptions_JSON? | 'TYPE' '=' 'AVRO' createExternalTable_formatTypeOptions_AVRO? | 'TYPE' '=' 'ORC' createExternalTable_formatTypeOptions_ORC? | 'TYPE' '=' 'PARQUET' createExternalTable_formatTypeOptions_PARQUET?)) ')' ('TABLE_FORMAT' '=' 'DELTA')? ('COPY' 'GRANTS')? ('COPY' '=' '(' 'GRANTS' ')')? ('COMMENT' '=' STRING)? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' 'VALUE' ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'EXTERNAL' 'TABLE' QNAME 'USING' 'TEMPLATE' '(' QUERY (')')+ ('COPY' 'GRANTS')? ('COPY' '=' '(' 'GRANTS' ')')?
    ;
createExternalTable_inlineConstraint
    : ('NOT' 'NULL')? ('CONSTRAINT' ID)? ('UNIQUE' | 'PRIMARY' 'KEY' | ('FOREIGN' 'KEY')? 'REFERENCES' QNAME ('(' ID (',' ID)? ')')? ('MATCH' ('FULL' | 'SIMPLE' | 'PARTIAL'))? ('ON' ('UPDATE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))? ('DELETE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))?)?) ('NOT'? 'ENFORCED')? ('NOT'? 'DEFERRABLE')? ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? (('RELY' | 'NORELY'))?
    ;
createExternalTable_cloudProviderParams
    : ('INTEGRATION' '=' STRING)?
    ;
createExternalTable_externalStage
    : STAGE
    ;
createExternalTable_formatTypeOptions_CSV
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'RECORD_DELIMITER' '=' (STRING | 'NONE') | 'FIELD_DELIMITER' '=' (STRING | 'NONE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'SKIP_HEADER' '=' NUMBER | 'SKIP_BLANK_LINES' '=' ('TRUE' | 'FALSE') | 'ESCAPE_UNENCLOSED_FIELD' '=' (STRING | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'FIELD_OPTIONALLY_ENCLOSED_BY' '=' (STRING | 'NONE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')' | 'EMPTY_FIELD_AS_NULL' '=' ('TRUE' | 'FALSE') | 'ENCODING' '=' STRING) (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'RECORD_DELIMITER' '=' (STRING | 'NONE') | 'FIELD_DELIMITER' '=' (STRING | 'NONE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'SKIP_HEADER' '=' NUMBER | 'SKIP_BLANK_LINES' '=' ('TRUE' | 'FALSE') | 'ESCAPE_UNENCLOSED_FIELD' '=' (STRING | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'FIELD_OPTIONALLY_ENCLOSED_BY' '=' (STRING | 'NONE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')' | 'EMPTY_FIELD_AS_NULL' '=' ('TRUE' | 'FALSE') | 'ENCODING' '=' STRING))*
    ;
createExternalTable_formatTypeOptions_JSON
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'ALLOW_DUPLICATE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ARRAY' '=' ('TRUE' | 'FALSE') | 'STRIP_NULL_VALUES' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE')) (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'ALLOW_DUPLICATE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ARRAY' '=' ('TRUE' | 'FALSE') | 'STRIP_NULL_VALUES' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE')))*
    ;
createExternalTable_formatTypeOptions_AVRO
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE')) (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE')))*
    ;
createExternalTable_formatTypeOptions_ORC
    : ('TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')') (','? ('TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')'))*
    ;
createExternalTable_formatTypeOptions_PARQUET
    : ('COMPRESSION' '=' ('AUTO' | 'SNAPPY' | 'NONE') | 'BINARY_AS_TEXT' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE')) (','? ('COMPRESSION' '=' ('AUTO' | 'SNAPPY' | 'NONE') | 'BINARY_AS_TEXT' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE')))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-external-volume
// not converted (form): ... stands for text the docs leave out
//   CREATE [ OR REPLACE ] EXTERNAL VOLUME [IF NOT EXISTS]
//     <name>
//     STORAGE_LOCATIONS =
//       (
//         (
//           NAME = '<storage_location_name>'
//           { cloudProviderParams | s3CompatibleStorageParams }
//         )
//         [, (...), ...]
//       )
//     [ ALLOW_WRITES = { TRUE | FALSE }]
//     [ COMMENT = '<string_literal>' ]

// https://docs.snowflake.com/en/sql-reference/sql/create-failover-group
// fix: OBJECT_TYPES takes the replication object types the parameters list, not any object type
// fix: ALLOWED_INTEGRATION_TYPES takes the integration types the parameters list
createFailoverGroup
    : 'CREATE' 'FAILOVER' 'GROUP' ('IF' 'NOT' 'EXISTS')? ID 'OBJECT_TYPES' '=' ('ACCOUNT' 'PARAMETERS' | 'DATABASES' | 'EXTERNAL' 'VOLUMES' | 'INTEGRATIONS' | 'LISTINGS' | 'NETWORK' 'POLICIES' | 'PROFILES' | 'RESOURCE' 'MONITORS' | 'ROLES' | 'SHARES' | 'USERS' | 'WAREHOUSES') (',' ('ACCOUNT' 'PARAMETERS' | 'DATABASES' | 'EXTERNAL' 'VOLUMES' | 'INTEGRATIONS' | 'LISTINGS' | 'NETWORK' 'POLICIES' | 'PROFILES' | 'RESOURCE' 'MONITORS' | 'ROLES' | 'SHARES' | 'USERS' | 'WAREHOUSES'))* ('ALLOWED_DATABASES' '=' ID (',' ID)*)? ('ALLOWED_EXTERNAL_VOLUMES' '=' ID (',' ID)*)? ('ALLOWED_SHARES' '=' ID (',' ID)*)? ('ALLOWED_INTEGRATION_TYPES' '=' ('SECURITY' 'INTEGRATIONS' | 'API' 'INTEGRATIONS' | 'STORAGE' 'INTEGRATIONS' | 'EXTERNAL' 'ACCESS' 'INTEGRATIONS' | 'NOTIFICATION' 'INTEGRATIONS') (',' ('SECURITY' 'INTEGRATIONS' | 'API' 'INTEGRATIONS' | 'STORAGE' 'INTEGRATIONS' | 'EXTERNAL' 'ACCESS' 'INTEGRATIONS' | 'NOTIFICATION' 'INTEGRATIONS'))*)? 'ALLOWED_ACCOUNTS' '=' ID '.' ID (',' ID '.' ID)* ('IGNORE' 'EDITION' 'CHECK')? ('REPLICATION_SCHEDULE' '=' STRING)? ('OPTIMIZED_REFRESH' '=' ('TRUE' | 'FALSE'))? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('ERROR_INTEGRATION' '=' ID)?
    | 'CREATE' 'FAILOVER' 'GROUP' ('IF' 'NOT' 'EXISTS')? ID 'AS' 'REPLICA' 'OF' ID '.' ID '.' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-feature-policy
// fix: <type> is an object type in the plural, DATABASES or TABLES, as the page's parameters list
createFeaturePolicy
    : 'CREATE' ('OR' 'REPLACE')? 'FEATURE' 'POLICY' QNAME ('BLOCKED_OBJECT_TYPES_FOR_CREATION' '=' '(' OBJECT_TYPES (',' OBJECT_TYPES)* ')')? ('COMMENT' '=' STRING)? ('AS' STRING)?
    | 'CREATE' 'FEATURE' 'POLICY' ('IF' 'NOT' 'EXISTS')? QNAME ('BLOCKED_OBJECT_TYPES_FOR_CREATION' '=' '(' OBJECT_TYPES (',' OBJECT_TYPES)* ')')? ('COMMENT' '=' STRING)? ('AS' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-file-format
// fix: read: each -- If TYPE = X section of formatTypeOptions is the options of TYPE = X alone, and TYPE = { ... } [ formatTypeOptions ] chooses a TYPE with its own options
createFileFormat
    : 'CREATE' ('OR' 'REPLACE')? (('TEMP' | 'TEMPORARY' | 'VOLATILE'))? 'FILE' 'FORMAT' QNAME (('TYPE' '=' 'CSV' createFileFormat_formatTypeOptions_CSV? | 'TYPE' '=' 'JSON' createFileFormat_formatTypeOptions_JSON? | 'TYPE' '=' 'AVRO' createFileFormat_formatTypeOptions_AVRO? | 'TYPE' '=' 'ORC' createFileFormat_formatTypeOptions_ORC? | 'TYPE' '=' 'PARQUET' createFileFormat_formatTypeOptions_PARQUET? | 'TYPE' '=' 'XML' createFileFormat_formatTypeOptions_XML?))? ('COMMENT' '=' STRING)?
    | 'CREATE' (('TEMP' | 'TEMPORARY' | 'VOLATILE'))? 'FILE' 'FORMAT' ('IF' 'NOT' 'EXISTS')? QNAME (('TYPE' '=' 'CSV' createFileFormat_formatTypeOptions_CSV? | 'TYPE' '=' 'JSON' createFileFormat_formatTypeOptions_JSON? | 'TYPE' '=' 'AVRO' createFileFormat_formatTypeOptions_AVRO? | 'TYPE' '=' 'ORC' createFileFormat_formatTypeOptions_ORC? | 'TYPE' '=' 'PARQUET' createFileFormat_formatTypeOptions_PARQUET? | 'TYPE' '=' 'XML' createFileFormat_formatTypeOptions_XML?))? ('COMMENT' '=' STRING)?
    | 'CREATE' 'OR' 'ALTER' (('TEMP' | 'TEMPORARY' | 'VOLATILE'))? 'FILE' 'FORMAT' QNAME (('TYPE' '=' 'CSV' createFileFormat_formatTypeOptions_CSV? | 'TYPE' '=' 'JSON' createFileFormat_formatTypeOptions_JSON? | 'TYPE' '=' 'AVRO' createFileFormat_formatTypeOptions_AVRO? | 'TYPE' '=' 'ORC' createFileFormat_formatTypeOptions_ORC? | 'TYPE' '=' 'PARQUET' createFileFormat_formatTypeOptions_PARQUET? | 'TYPE' '=' 'XML' createFileFormat_formatTypeOptions_XML?))? ('COMMENT' '=' STRING)?
    ;
createFileFormat_formatTypeOptions_CSV
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'RECORD_DELIMITER' '=' (STRING | 'NONE') | 'FIELD_DELIMITER' '=' (STRING | 'NONE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'FILE_EXTENSION' '=' STRING | 'PARSE_HEADER' '=' ('TRUE' | 'FALSE') | 'SKIP_HEADER' '=' NUMBER | 'SKIP_BLANK_LINES' '=' ('TRUE' | 'FALSE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'ESCAPE' '=' (STRING | 'NONE') | 'ESCAPE_UNENCLOSED_FIELD' '=' (STRING | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'FIELD_OPTIONALLY_ENCLOSED_BY' '=' (STRING | 'NONE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')' | 'ERROR_ON_COLUMN_COUNT_MISMATCH' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'EMPTY_FIELD_AS_NULL' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE') | 'ENCODING' '=' (STRING | 'UTF8')) (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'RECORD_DELIMITER' '=' (STRING | 'NONE') | 'FIELD_DELIMITER' '=' (STRING | 'NONE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'FILE_EXTENSION' '=' STRING | 'PARSE_HEADER' '=' ('TRUE' | 'FALSE') | 'SKIP_HEADER' '=' NUMBER | 'SKIP_BLANK_LINES' '=' ('TRUE' | 'FALSE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'ESCAPE' '=' (STRING | 'NONE') | 'ESCAPE_UNENCLOSED_FIELD' '=' (STRING | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'FIELD_OPTIONALLY_ENCLOSED_BY' '=' (STRING | 'NONE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')' | 'ERROR_ON_COLUMN_COUNT_MISMATCH' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'EMPTY_FIELD_AS_NULL' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE') | 'ENCODING' '=' (STRING | 'UTF8')))*
    ;
createFileFormat_formatTypeOptions_JSON
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')' | 'FILE_EXTENSION' '=' STRING | 'ENABLE_OCTAL' '=' ('TRUE' | 'FALSE') | 'ALLOW_DUPLICATE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ARRAY' '=' ('TRUE' | 'FALSE') | 'STRIP_NULL_VALUES' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'IGNORE_UTF8_ERRORS' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE')) (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')' | 'FILE_EXTENSION' '=' STRING | 'ENABLE_OCTAL' '=' ('TRUE' | 'FALSE') | 'ALLOW_DUPLICATE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ARRAY' '=' ('TRUE' | 'FALSE') | 'STRIP_NULL_VALUES' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'IGNORE_UTF8_ERRORS' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE')))*
    ;
createFileFormat_formatTypeOptions_AVRO
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')') (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')'))*
    ;
createFileFormat_formatTypeOptions_ORC
    : ('TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')') (','? ('TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')'))*
    ;
createFileFormat_formatTypeOptions_PARQUET
    : ('COMPRESSION' '=' ('AUTO' | 'LZO' | 'SNAPPY' | 'NONE') | 'SNAPPY_COMPRESSION' '=' ('TRUE' | 'FALSE') | 'BINARY_AS_TEXT' '=' ('TRUE' | 'FALSE') | 'USE_LOGICAL_TYPE' '=' ('TRUE' | 'FALSE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'USE_VECTORIZED_SCANNER' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')') (','? ('COMPRESSION' '=' ('AUTO' | 'LZO' | 'SNAPPY' | 'NONE') | 'SNAPPY_COMPRESSION' '=' ('TRUE' | 'FALSE') | 'BINARY_AS_TEXT' '=' ('TRUE' | 'FALSE') | 'USE_LOGICAL_TYPE' '=' ('TRUE' | 'FALSE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'USE_VECTORIZED_SCANNER' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')'))*
    ;
createFileFormat_formatTypeOptions_XML
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'IGNORE_UTF8_ERRORS' '=' ('TRUE' | 'FALSE') | 'PRESERVE_SPACE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ELEMENT' '=' ('TRUE' | 'FALSE') | 'DISABLE_AUTO_CONVERT' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE')) (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'IGNORE_UTF8_ERRORS' '=' ('TRUE' | 'FALSE') | 'PRESERVE_SPACE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ELEMENT' '=' ('TRUE' | 'FALSE') | 'DISABLE_AUTO_CONVERT' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE')))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-function
// not converted (form): ... stands for text the docs leave out
//   CREATE [ OR ALTER ] FUNCTION ...
createFunction
    : 'CREATE' ('OR' 'REPLACE')? (('TEMP' | 'TEMPORARY'))? 'SECURE'? 'FUNCTION' QNAME '(' (ID DATA_TYPE ('DEFAULT' EXPR)? (',' ID DATA_TYPE ('DEFAULT' EXPR)?)*)? ')' ('COPY' 'GRANTS')? 'RETURNS' (DATA_TYPE | 'TABLE' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')') ('NOT'? 'NULL')? 'LANGUAGE' 'JAVA' (('CALLED' 'ON' 'NULL' 'INPUT' | ('RETURNS' 'NULL' 'ON' 'NULL' 'INPUT' | 'STRICT')))? (('VOLATILE' | 'IMMUTABLE'))? ('RUNTIME_VERSION' '=' STRING)? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')? ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? ('PACKAGES' '=' '(' STRING (',' STRING)* ')')? 'HANDLER' '=' STRING ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')')? ('TARGET_PATH' '=' STRING)? 'AS' STRING
    | 'CREATE' (('TEMP' | 'TEMPORARY'))? 'SECURE'? 'FUNCTION' ('IF' 'NOT' 'EXISTS')? QNAME '(' (ID DATA_TYPE ('DEFAULT' EXPR)? (',' ID DATA_TYPE ('DEFAULT' EXPR)?)*)? ')' ('COPY' 'GRANTS')? 'RETURNS' (DATA_TYPE | 'TABLE' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')') ('NOT'? 'NULL')? 'LANGUAGE' 'JAVA' (('CALLED' 'ON' 'NULL' 'INPUT' | ('RETURNS' 'NULL' 'ON' 'NULL' 'INPUT' | 'STRICT')))? (('VOLATILE' | 'IMMUTABLE'))? ('RUNTIME_VERSION' '=' STRING)? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')? ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? ('PACKAGES' '=' '(' STRING (',' STRING)* ')')? 'HANDLER' '=' STRING ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')')? ('TARGET_PATH' '=' STRING)? 'AS' STRING
    | 'CREATE' ('OR' 'REPLACE')? (('TEMP' | 'TEMPORARY'))? 'SECURE'? 'FUNCTION' QNAME '(' (ID DATA_TYPE ('DEFAULT' EXPR)? (',' ID DATA_TYPE ('DEFAULT' EXPR)?)*)? ')' ('COPY' 'GRANTS')? 'RETURNS' (DATA_TYPE | 'TABLE' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')') ('NOT'? 'NULL')? 'LANGUAGE' 'JAVA' (('CALLED' 'ON' 'NULL' 'INPUT' | ('RETURNS' 'NULL' 'ON' 'NULL' 'INPUT' | 'STRICT')))? (('VOLATILE' | 'IMMUTABLE'))? ('RUNTIME_VERSION' '=' STRING)? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')? 'IMPORTS' '=' '(' STRING (',' STRING)* ')' 'HANDLER' '=' STRING ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')')?
    | 'CREATE' (('TEMP' | 'TEMPORARY'))? 'SECURE'? 'FUNCTION' ('IF' 'NOT' 'EXISTS')? QNAME '(' (ID DATA_TYPE ('DEFAULT' EXPR)? (',' ID DATA_TYPE ('DEFAULT' EXPR)?)*)? ')' ('COPY' 'GRANTS')? 'RETURNS' (DATA_TYPE | 'TABLE' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')') ('NOT'? 'NULL')? 'LANGUAGE' 'JAVA' (('CALLED' 'ON' 'NULL' 'INPUT' | ('RETURNS' 'NULL' 'ON' 'NULL' 'INPUT' | 'STRICT')))? (('VOLATILE' | 'IMMUTABLE'))? ('RUNTIME_VERSION' '=' STRING)? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')? 'IMPORTS' '=' '(' STRING (',' STRING)* ')' 'HANDLER' '=' STRING ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')')?
    | 'CREATE' ('OR' 'REPLACE')? (('TEMP' | 'TEMPORARY'))? 'SECURE'? 'FUNCTION' QNAME '(' (ID DATA_TYPE ('DEFAULT' EXPR)? (',' ID DATA_TYPE ('DEFAULT' EXPR)?)*)? ')' ('COPY' 'GRANTS')? 'RETURNS' (DATA_TYPE | 'TABLE' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')') ('NOT'? 'NULL')? 'LANGUAGE' 'JAVASCRIPT' (('CALLED' 'ON' 'NULL' 'INPUT' | ('RETURNS' 'NULL' 'ON' 'NULL' 'INPUT' | 'STRICT')))? (('VOLATILE' | 'IMMUTABLE'))? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')? 'AS' STRING
    | 'CREATE' ('OR' 'REPLACE')? (('TEMP' | 'TEMPORARY'))? 'SECURE'? 'AGGREGATE'? 'FUNCTION' QNAME '(' (ID DATA_TYPE ('DEFAULT' EXPR)? (',' ID DATA_TYPE ('DEFAULT' EXPR)?)*)? ')' ('COPY' 'GRANTS')? 'RETURNS' (DATA_TYPE | 'TABLE' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')') ('NOT'? 'NULL')? 'LANGUAGE' 'PYTHON' (('CALLED' 'ON' 'NULL' 'INPUT' | ('RETURNS' 'NULL' 'ON' 'NULL' 'INPUT' | 'STRICT')))? (('VOLATILE' | 'IMMUTABLE'))? 'RUNTIME_VERSION' '=' STRING ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')? ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? ('PACKAGES' '=' '(' STRING (',' STRING)* ')')? ('ARTIFACT_REPOSITORY' '=' STRING)? 'HANDLER' '=' STRING ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')')? 'AS' STRING
    | 'CREATE' (('TEMP' | 'TEMPORARY'))? 'SECURE'? 'AGGREGATE'? 'FUNCTION' ('IF' 'NOT' 'EXISTS')? QNAME '(' (ID DATA_TYPE ('DEFAULT' EXPR)? (',' ID DATA_TYPE ('DEFAULT' EXPR)?)*)? ')' ('COPY' 'GRANTS')? 'RETURNS' (DATA_TYPE | 'TABLE' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')') ('NOT'? 'NULL')? 'LANGUAGE' 'PYTHON' (('CALLED' 'ON' 'NULL' 'INPUT' | ('RETURNS' 'NULL' 'ON' 'NULL' 'INPUT' | 'STRICT')))? (('VOLATILE' | 'IMMUTABLE'))? 'RUNTIME_VERSION' '=' STRING ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')? ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? ('PACKAGES' '=' '(' STRING (',' STRING)* ')')? ('ARTIFACT_REPOSITORY' '=' STRING)? 'HANDLER' '=' STRING ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')')? 'AS' STRING
    | 'CREATE' ('OR' 'REPLACE')? (('TEMP' | 'TEMPORARY'))? 'SECURE'? 'AGGREGATE'? 'FUNCTION' QNAME '(' (ID DATA_TYPE ('DEFAULT' EXPR)? (',' ID DATA_TYPE ('DEFAULT' EXPR)?)*)? ')' ('COPY' 'GRANTS')? 'RETURNS' (DATA_TYPE | 'TABLE' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')') ('NOT'? 'NULL')? 'LANGUAGE' 'PYTHON' (('CALLED' 'ON' 'NULL' 'INPUT' | ('RETURNS' 'NULL' 'ON' 'NULL' 'INPUT' | 'STRICT')))? (('VOLATILE' | 'IMMUTABLE'))? 'RUNTIME_VERSION' '=' STRING ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')? ('ARTIFACT_REPOSITORY' '=' STRING)? 'IMPORTS' '=' '(' STRING (',' STRING)* ')' ('PACKAGES' '=' '(' STRING (',' STRING)* ')')? 'HANDLER' '=' STRING ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')')?
    | 'CREATE' (('TEMP' | 'TEMPORARY'))? 'SECURE'? 'AGGREGATE'? 'FUNCTION' ('IF' 'NOT' 'EXISTS')? QNAME '(' (ID DATA_TYPE ('DEFAULT' EXPR)? (',' ID DATA_TYPE ('DEFAULT' EXPR)?)*)? ')' ('COPY' 'GRANTS')? 'RETURNS' (DATA_TYPE | 'TABLE' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')') ('NOT'? 'NULL')? 'LANGUAGE' 'PYTHON' (('CALLED' 'ON' 'NULL' 'INPUT' | ('RETURNS' 'NULL' 'ON' 'NULL' 'INPUT' | 'STRICT')))? (('VOLATILE' | 'IMMUTABLE'))? 'RUNTIME_VERSION' '=' STRING ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')? ('ARTIFACT_REPOSITORY' '=' STRING)? 'IMPORTS' '=' '(' STRING (',' STRING)* ')' ('PACKAGES' '=' '(' STRING (',' STRING)* ')')? 'HANDLER' '=' STRING ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')')?
    | 'CREATE' ('OR' 'REPLACE')? (('TEMP' | 'TEMPORARY'))? 'SECURE'? 'FUNCTION' QNAME '(' (ID DATA_TYPE ('DEFAULT' EXPR)? (',' ID DATA_TYPE ('DEFAULT' EXPR)?)*)? ')' ('COPY' 'GRANTS')? 'RETURNS' DATA_TYPE ('NOT'? 'NULL')? 'LANGUAGE' 'SCALA' (('CALLED' 'ON' 'NULL' 'INPUT' | ('RETURNS' 'NULL' 'ON' 'NULL' 'INPUT' | 'STRICT')))? (('VOLATILE' | 'IMMUTABLE'))? ('RUNTIME_VERSION' '=' STRING)? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')? ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? ('PACKAGES' '=' '(' STRING (',' STRING)* ')')? 'HANDLER' '=' STRING ('TARGET_PATH' '=' STRING)? 'AS' STRING
    | 'CREATE' (('TEMP' | 'TEMPORARY'))? 'SECURE'? 'FUNCTION' ('IF' 'NOT' 'EXISTS')? QNAME '(' (ID DATA_TYPE ('DEFAULT' EXPR)? (',' ID DATA_TYPE ('DEFAULT' EXPR)?)*)? ')' ('COPY' 'GRANTS')? 'RETURNS' DATA_TYPE ('NOT'? 'NULL')? 'LANGUAGE' 'SCALA' (('CALLED' 'ON' 'NULL' 'INPUT' | ('RETURNS' 'NULL' 'ON' 'NULL' 'INPUT' | 'STRICT')))? (('VOLATILE' | 'IMMUTABLE'))? ('RUNTIME_VERSION' '=' STRING)? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')? ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? ('PACKAGES' '=' '(' STRING (',' STRING)* ')')? 'HANDLER' '=' STRING ('TARGET_PATH' '=' STRING)? 'AS' STRING
    | 'CREATE' ('OR' 'REPLACE')? (('TEMP' | 'TEMPORARY'))? 'SECURE'? 'FUNCTION' QNAME '(' (ID DATA_TYPE ('DEFAULT' EXPR)? (',' ID DATA_TYPE ('DEFAULT' EXPR)?)*)? ')' ('COPY' 'GRANTS')? 'RETURNS' DATA_TYPE ('NOT'? 'NULL')? 'LANGUAGE' 'SCALA' (('CALLED' 'ON' 'NULL' 'INPUT' | ('RETURNS' 'NULL' 'ON' 'NULL' 'INPUT' | 'STRICT')))? (('VOLATILE' | 'IMMUTABLE'))? ('RUNTIME_VERSION' '=' STRING)? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')? 'IMPORTS' '=' '(' STRING (',' STRING)* ')' 'HANDLER' '=' STRING
    | 'CREATE' (('TEMP' | 'TEMPORARY'))? 'SECURE'? 'FUNCTION' ('IF' 'NOT' 'EXISTS')? QNAME '(' (ID DATA_TYPE ('DEFAULT' EXPR)? (',' ID DATA_TYPE ('DEFAULT' EXPR)?)*)? ')' ('COPY' 'GRANTS')? 'RETURNS' DATA_TYPE ('NOT'? 'NULL')? 'LANGUAGE' 'SCALA' (('CALLED' 'ON' 'NULL' 'INPUT' | ('RETURNS' 'NULL' 'ON' 'NULL' 'INPUT' | 'STRICT')))? (('VOLATILE' | 'IMMUTABLE'))? ('RUNTIME_VERSION' '=' STRING)? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')? 'IMPORTS' '=' '(' STRING (',' STRING)* ')' 'HANDLER' '=' STRING
    | 'CREATE' ('OR' 'REPLACE')? (('TEMP' | 'TEMPORARY'))? 'SECURE'? 'FUNCTION' QNAME '(' (ID DATA_TYPE ('DEFAULT' EXPR)? (',' ID DATA_TYPE ('DEFAULT' EXPR)?)*)? ')' ('COPY' 'GRANTS')? 'RETURNS' (DATA_TYPE | 'TABLE' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')') ('NOT'? 'NULL')? (('VOLATILE' | 'IMMUTABLE'))? 'MEMOIZABLE'? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')? 'AS' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-function-spcs
// not converted (form): ... stands for text the docs leave out
//   CREATE [ OR ALTER ] FUNCTION ...
createFunctionSpcs
    : 'CREATE' ('OR' 'REPLACE')? 'FUNCTION' QNAME '(' (ID DATA_TYPE (',' ID DATA_TYPE)*)? ')' 'RETURNS' DATA_TYPE ('NOT'? 'NULL')? (('CALLED' 'ON' 'NULL' 'INPUT' | ('RETURNS' 'NULL' 'ON' 'NULL' 'INPUT' | 'STRICT')))? (('VOLATILE' | 'IMMUTABLE'))? 'SERVICE' '=' QNAME 'ENDPOINT' '=' ID ('COMMENT' '=' STRING)? ('CONTEXT_HEADERS' '=' '(' ID (',' ID)* ')')? ('MAX_BATCH_ROWS' '=' NUMBER)? ('MAX_BATCH_RETRIES' '=' NUMBER)? ('ON_BATCH_FAILURE' '=' ('ABORT' | 'RETURN_NULL'))? ('BATCH_TIMEOUT_SECS' '=' NUMBER)? 'AS' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-gateway
createGateway
    : 'CREATE' ('OR' 'REPLACE')? 'GATEWAY' QNAME 'FROM' 'SPECIFICATION' STRING
    | 'CREATE' 'GATEWAY' ('IF' 'NOT' 'EXISTS')? QNAME 'FROM' 'SPECIFICATION' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-git-repository
createGitRepository
    : 'CREATE' ('OR' 'REPLACE')? 'GIT' 'REPOSITORY' QNAME 'ORIGIN' '=' STRING 'API_INTEGRATION' '=' ID ('GIT_CREDENTIALS' '=' QNAME)? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' 'GIT' 'REPOSITORY' ('IF' 'NOT' 'EXISTS')? QNAME 'ORIGIN' '=' STRING 'API_INTEGRATION' '=' ID ('GIT_CREDENTIALS' '=' QNAME)? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-hybrid-table
// fix: read: <constraint_properties> is the properties CREATE TABLE's constraint page lists
// fix: read: an out-of-line key constraint names its columns, as the constraint page says
// fix: read: a foreign key's REFERENCES takes MATCH and ON UPDATE/DELETE, the properties the constraint page gives foreign keys
createHybridTable
    : 'CREATE' ('OR' 'REPLACE')? 'HYBRID' 'TABLE' QNAME '(' ID DATA_TYPE (('DEFAULT' EXPR | ('AUTOINCREMENT' | 'IDENTITY') (('(' NUMBER ',' NUMBER ')' | 'START' NUMBER 'INCREMENT' NUMBER))? (('ORDER' | 'NOORDER'))?))? ('NOT' 'NULL')? createHybridTable_inlineConstraint? ('COLLATE' STRING)? ('COMMENT' STRING)? (',' ID DATA_TYPE)* (',' createHybridTable_outoflineConstraint)? (',' createHybridTable_outoflineIndex)* ')' ('COMMENT' '=' STRING)?
    | 'CREATE' 'HYBRID' 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME '(' ID DATA_TYPE (('DEFAULT' EXPR | ('AUTOINCREMENT' | 'IDENTITY') (('(' NUMBER ',' NUMBER ')' | 'START' NUMBER 'INCREMENT' NUMBER))? (('ORDER' | 'NOORDER'))?))? ('NOT' 'NULL')? createHybridTable_inlineConstraint? ('COLLATE' STRING)? ('COMMENT' STRING)? (',' ID DATA_TYPE)* (',' createHybridTable_outoflineConstraint)? (',' createHybridTable_outoflineIndex)* ')' ('COMMENT' '=' STRING)?
    ;
createHybridTable_inlineConstraint
    : ('CONSTRAINT' ID)? ('UNIQUE' | 'PRIMARY' 'KEY' | ('FOREIGN' 'KEY')? 'REFERENCES' QNAME ('(' ID ')')? ('MATCH' ('FULL' | 'SIMPLE' | 'PARTIAL'))? ('ON' ('UPDATE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))? ('DELETE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))?)? | 'CHECK' '(' EXPR ')') ('NOT'? 'ENFORCED')? ('NOT'? 'DEFERRABLE')? ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? (('RELY' | 'NORELY'))?
    ;
createHybridTable_outoflineConstraint
    : ('CONSTRAINT' ID)? ('UNIQUE' '(' ID (',' ID)* ')' | 'PRIMARY' 'KEY' '(' ID (',' ID)* ')' | ('FOREIGN' 'KEY')? '(' ID (',' ID)* ')' 'REFERENCES' QNAME ('(' ID (',' ID)* ')')? ('MATCH' ('FULL' | 'SIMPLE' | 'PARTIAL'))? ('ON' ('UPDATE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))? ('DELETE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))?)? | 'CHECK' '(' EXPR ')') ('NOT'? 'ENFORCED')? ('NOT'? 'DEFERRABLE')? ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? (('RELY' | 'NORELY'))? ('COMMENT' STRING)?
    ;
createHybridTable_outoflineIndex
    : 'INDEX' ID '(' ID (',' ID)* ')' ('INCLUDE' '(' ID (',' ID)* ')')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table
// fix: the column definitions have no parentheses; the page's examples write them, CREATE ICEBERG TABLE t (col1 INT)
// fix: the column definitions have no parentheses; the page's examples write them, CREATE ICEBERG TABLE t (col1 INT)
// fix: the column definitions have no parentheses; the page's examples write them, CREATE ICEBERG TABLE t (col1 INT)
// fix: the column definitions have no parentheses; the page's examples write them, CREATE ICEBERG TABLE t (col1 INT)
// fix: each additional column, data, virtual or partition, follows a comma, as the page's examples write them
// fix: read: <constraint_properties> is the properties CREATE TABLE's constraint page lists
// fix: read: <constraint_properties> is the properties CREATE TABLE's constraint page lists
// fix: read: <constraint_properties> is the properties CREATE TABLE's constraint page lists
// fix: read: <constraint_properties> is the properties CREATE TABLE's constraint page lists
// fix: read: an out-of-line key constraint names its columns, as the constraint page says
// fix: read: a foreign key's REFERENCES takes MATCH and ON UPDATE/DELETE, the properties the constraint page gives foreign keys
// fix: read: the additional columns take the whole column definition, which their line abbreviates
createIcebergTable
    : 'CREATE' ('OR' 'REPLACE')? 'TRANSIENT'? 'ICEBERG' 'TABLE' QNAME '(' ID DATA_TYPE ('DEFAULT' EXPR)? createIcebergTable_inlineConstraint? ('NOT' 'NULL')? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)? (',' ID DATA_TYPE ('DEFAULT' EXPR)? createIcebergTable_inlineConstraint? ('NOT' 'NULL')? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)?)* (',' createIcebergTable_outoflineConstraint)* ')' ('PARTITION' 'BY' '(' createIcebergTable_partitionExpression (',' createIcebergTable_partitionExpression)* ')')? ('PATH_LAYOUT' '=' ('FLAT' | 'HIERARCHICAL'))? ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')')? ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' '\'SNOWFLAKE\'')? ('BASE_LOCATION' '=' STRING)? ('TARGET_FILE_SIZE' '=' ('\'AUTO\'' | '\'16MB\'' | '\'32MB\'' | '\'64MB\'' | '\'128MB\''))? ('CATALOG_SYNC' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('CHANGE_TRACKING' '=' ('TRUE' | 'FALSE'))? ('COPY' 'GRANTS')? ('COPY' 'TAGS')? ('COPY' '=' '(' ('COMMENTS' | 'TAGS' | 'GRANTS') (',' ('COMMENTS' | 'TAGS' | 'GRANTS'))* ')')? ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('ICEBERG_VERSION' '=' NUMBER)? ('ALLOW_V3_EXPRESSION_DEFAULTS' '=' ('TRUE' | 'FALSE'))? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'AGGREGATION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))?
    | 'CREATE' 'TRANSIENT'? 'ICEBERG' 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME '(' ID DATA_TYPE ('DEFAULT' EXPR)? createIcebergTable_inlineConstraint? ('NOT' 'NULL')? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)? (',' ID DATA_TYPE ('DEFAULT' EXPR)? createIcebergTable_inlineConstraint? ('NOT' 'NULL')? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)?)* (',' createIcebergTable_outoflineConstraint)* ')' ('PARTITION' 'BY' '(' createIcebergTable_partitionExpression (',' createIcebergTable_partitionExpression)* ')')? ('PATH_LAYOUT' '=' ('FLAT' | 'HIERARCHICAL'))? ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')')? ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' '\'SNOWFLAKE\'')? ('BASE_LOCATION' '=' STRING)? ('TARGET_FILE_SIZE' '=' ('\'AUTO\'' | '\'16MB\'' | '\'32MB\'' | '\'64MB\'' | '\'128MB\''))? ('CATALOG_SYNC' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('CHANGE_TRACKING' '=' ('TRUE' | 'FALSE'))? ('COPY' 'GRANTS')? ('COPY' 'TAGS')? ('COPY' '=' '(' ('COMMENTS' | 'TAGS' | 'GRANTS') (',' ('COMMENTS' | 'TAGS' | 'GRANTS'))* ')')? ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('ICEBERG_VERSION' '=' NUMBER)? ('ALLOW_V3_EXPRESSION_DEFAULTS' '=' ('TRUE' | 'FALSE'))? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'AGGREGATION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))?
    | 'CREATE' ('OR' 'REPLACE')? 'TRANSIENT'? 'ICEBERG' 'TABLE' QNAME ('(' ID DATA_TYPE? ('DEFAULT' EXPR)? (',' ID DATA_TYPE? ('DEFAULT' EXPR)?)* ')')? ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')')? ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' '\'SNOWFLAKE\'')? ('BASE_LOCATION' '=' STRING)? ('COPY' 'GRANTS')? ('COPY' 'TAGS')? ('COPY' '=' '(' ('COMMENTS' | 'TAGS' | 'GRANTS') (',' ('COMMENTS' | 'TAGS' | 'GRANTS'))* ')')? ('ICEBERG_VERSION' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))* 'AS' QUERY
    | 'CREATE' ('OR' 'REPLACE')? 'TRANSIENT'? 'ICEBERG' 'TABLE' QNAME 'LIKE' QNAME ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')')? ('COPY' 'GRANTS')? ('COPY' 'TAGS')? ('COPY' '=' '(' ('COMMENTS' | 'TAGS' | 'GRANTS') (',' ('COMMENTS' | 'TAGS' | 'GRANTS'))* ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'ICEBERG' 'TABLE' QNAME ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' STRING)? 'CATALOG_TABLE_NAME' '=' STRING ('CATALOG_NAMESPACE' '=' STRING)? ('PATH_LAYOUT' '=' ('FLAT' | 'HIERARCHICAL'))? ('TARGET_FILE_SIZE' '=' ('\'AUTO\'' | '\'16MB\'' | '\'32MB\'' | '\'64MB\'' | '\'128MB\''))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    | 'CREATE' 'ICEBERG' 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' STRING)? 'CATALOG_TABLE_NAME' '=' STRING ('CATALOG_NAMESPACE' '=' STRING)? ('PATH_LAYOUT' '=' ('FLAT' | 'HIERARCHICAL'))? ('TARGET_FILE_SIZE' '=' ('\'AUTO\'' | '\'16MB\'' | '\'32MB\'' | '\'64MB\'' | '\'128MB\''))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'ICEBERG' 'TABLE' QNAME ('(' ID DATA_TYPE ('DEFAULT' EXPR)? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? (',' ID DATA_TYPE ('DEFAULT' EXPR)? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)?)* ')')? ('PARTITION' 'BY' '(' createIcebergTable_partitionExpression (',' createIcebergTable_partitionExpression)* ')')? ('PATH_LAYOUT' '=' ('FLAT' | 'HIERARCHICAL'))? ('TARGET_FILE_SIZE' '=' ('\'AUTO\'' | '\'16MB\'' | '\'32MB\'' | '\'64MB\'' | '\'128MB\''))? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('COPY' 'GRANTS')? 'PURGE'? ('COMMENT' '=' STRING)? ('ICEBERG_VERSION' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('BASE_LOCATION' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))?
    | 'CREATE' 'ICEBERG' 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME ('(' ID DATA_TYPE ('DEFAULT' EXPR)? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? (',' ID DATA_TYPE ('DEFAULT' EXPR)? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)?)* ')')? ('PARTITION' 'BY' '(' createIcebergTable_partitionExpression (',' createIcebergTable_partitionExpression)* ')')? ('PATH_LAYOUT' '=' ('FLAT' | 'HIERARCHICAL'))? ('TARGET_FILE_SIZE' '=' ('\'AUTO\'' | '\'16MB\'' | '\'32MB\'' | '\'64MB\'' | '\'128MB\''))? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('COPY' 'GRANTS')? 'PURGE'? ('COMMENT' '=' STRING)? ('ICEBERG_VERSION' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('BASE_LOCATION' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))?
    | 'CREATE' ('OR' 'REPLACE')? 'ICEBERG' 'TABLE' QNAME ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' STRING)? 'BASE_LOCATION' '=' STRING ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    | 'CREATE' 'ICEBERG' 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' STRING)? 'BASE_LOCATION' '=' STRING ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'ICEBERG' 'TABLE' QNAME ('(' ID DATA_TYPE ('NOT' 'NULL')? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)? createIcebergTable_inlineConstraint? (',' (ID DATA_TYPE | ID DATA_TYPE 'AS' EXPR) createIcebergTable_inlineConstraint?)* (',' createIcebergTable_outoflineConstraint)* ')')? ('PARTITION' 'BY' '(' ID (',' ID)* ')')? ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' STRING)? 'BASE_LOCATION' '=' STRING ('ENABLE_SCHEMA_INFERENCE' '=' ('TRUE' | 'FALSE'))? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('MATCH_BY_COLUMN_NAME' '=' ('CASE_INSENSITIVE' | 'CASE_SENSITIVE'))? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('REFRESH_ON_CREATE' '=' ('TRUE' | 'FALSE'))? ('PATTERN' '=' STRING)? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' 'ICEBERG' 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME ('(' ID DATA_TYPE ('NOT' 'NULL')? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)? createIcebergTable_inlineConstraint? (',' (ID DATA_TYPE | ID DATA_TYPE 'AS' EXPR) createIcebergTable_inlineConstraint?)* (',' createIcebergTable_outoflineConstraint)* ')')? ('PARTITION' 'BY' '(' ID (',' ID)* ')')? ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' STRING)? 'BASE_LOCATION' '=' STRING ('ENABLE_SCHEMA_INFERENCE' '=' ('TRUE' | 'FALSE'))? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('MATCH_BY_COLUMN_NAME' '=' ('CASE_INSENSITIVE' | 'CASE_SENSITIVE'))? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('REFRESH_ON_CREATE' '=' ('TRUE' | 'FALSE'))? ('PATTERN' '=' STRING)? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' ('OR' 'REPLACE')? 'ICEBERG' 'TABLE' QNAME ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' STRING)? 'METADATA_FILE_PATH' '=' STRING ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    | 'CREATE' 'ICEBERG' 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' STRING)? 'METADATA_FILE_PATH' '=' STRING ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    ;
createIcebergTable_inlineConstraint
    : ('CONSTRAINT' ID)? ('UNIQUE' | 'PRIMARY' 'KEY' | ('FOREIGN' 'KEY')? 'REFERENCES' QNAME ('(' ID ')')? ('MATCH' ('FULL' | 'SIMPLE' | 'PARTIAL'))? ('ON' ('UPDATE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))? ('DELETE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))?)? | 'CHECK' '(' EXPR ')') ('NOT'? 'ENFORCED')? ('NOT'? 'DEFERRABLE')? ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? (('RELY' | 'NORELY'))?
    | ('CONSTRAINT' ID)? ('UNIQUE' | 'PRIMARY' 'KEY' | ('FOREIGN' 'KEY')? 'REFERENCES' QNAME ('(' ID ')')? ('MATCH' ('FULL' | 'SIMPLE' | 'PARTIAL'))? ('ON' ('UPDATE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))? ('DELETE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))?)?) ('NOT'? 'ENFORCED')? ('NOT'? 'DEFERRABLE')? ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? (('RELY' | 'NORELY'))?
    ;
createIcebergTable_outoflineConstraint
    : ('CONSTRAINT' ID)? ('UNIQUE' '(' ID (',' ID)* ')' | 'PRIMARY' 'KEY' '(' ID (',' ID)* ')' | ('FOREIGN' 'KEY')? '(' ID (',' ID)* ')' 'REFERENCES' QNAME ('(' ID (',' ID)* ')')? ('MATCH' ('FULL' | 'SIMPLE' | 'PARTIAL'))? ('ON' ('UPDATE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))? ('DELETE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))?)? | 'CHECK' '(' EXPR ')') ('NOT'? 'ENFORCED')? ('NOT'? 'DEFERRABLE')? ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? (('RELY' | 'NORELY'))?
    | ('CONSTRAINT' ID)? ('UNIQUE' '(' ID (',' ID)* ')' | 'PRIMARY' 'KEY' '(' ID (',' ID)* ')' | ('FOREIGN' 'KEY')? '(' ID (',' ID)* ')' 'REFERENCES' QNAME ('(' ID (',' ID)* ')')? ('MATCH' ('FULL' | 'SIMPLE' | 'PARTIAL'))? ('ON' ('UPDATE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))? ('DELETE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))?)?) ('NOT'? 'ENFORCED')? ('NOT'? 'DEFERRABLE')? ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? (('RELY' | 'NORELY'))?
    ;
createIcebergTable_partitionExpression
    : ID
    | 'BUCKET' '(' NUMBER ',' ID ')'
    | 'TRUNCATE' '(' NUMBER ',' ID ')'
    | 'YEAR' '(' ID ')'
    | 'MONTH' '(' ID ')'
    | 'DAY' '(' ID ')'
    | 'HOUR' '(' ID ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-aws-glue
createIcebergTableAwsGlue
    : 'CREATE' ('OR' 'REPLACE')? 'ICEBERG' 'TABLE' QNAME ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' STRING)? 'CATALOG_TABLE_NAME' '=' STRING ('CATALOG_NAMESPACE' '=' STRING)? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    | 'CREATE' 'ICEBERG' 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' STRING)? 'CATALOG_TABLE_NAME' '=' STRING ('CATALOG_NAMESPACE' '=' STRING)? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-iceberg-files
createIcebergTableIcebergFiles
    : 'CREATE' ('OR' 'REPLACE')? 'ICEBERG' 'TABLE' QNAME ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' STRING)? 'METADATA_FILE_PATH' '=' STRING ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    | 'CREATE' 'ICEBERG' 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' STRING)? 'METADATA_FILE_PATH' '=' STRING ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-parquet
// fix: the column definitions have no parentheses; the page's examples write them, CREATE ICEBERG TABLE t (col1 INT)
// fix: the column definitions have no parentheses; the page's examples write them, CREATE ICEBERG TABLE t (col1 INT)
// fix: each additional column, data, virtual or partition, follows a comma, as the page's examples write them
// fix: read: <constraint_properties> is the properties CREATE TABLE's constraint page lists
// fix: read: <constraint_properties> is the properties CREATE TABLE's constraint page lists
// fix: read: an out-of-line key constraint names its columns, as the constraint page says
// fix: read: a foreign key's REFERENCES takes MATCH and ON UPDATE/DELETE, the properties the constraint page gives foreign keys
createIcebergTableParquet
    : 'CREATE' ('OR' 'REPLACE')? 'ICEBERG' 'TABLE' QNAME ('(' ID DATA_TYPE ('NOT' 'NULL')? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)? createIcebergTableParquet_inlineConstraint? (',' (ID DATA_TYPE | ID DATA_TYPE 'AS' EXPR) createIcebergTableParquet_inlineConstraint?)* (',' createIcebergTableParquet_outoflineConstraint)* ')')? ('PARTITION' 'BY' '(' ID (',' ID)* ')')? ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' STRING)? 'BASE_LOCATION' '=' STRING ('ENABLE_SCHEMA_INFERENCE' '=' ('TRUE' | 'FALSE'))? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('MATCH_BY_COLUMN_NAME' '=' ('CASE_INSENSITIVE' | 'CASE_SENSITIVE'))? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('REFRESH_ON_CREATE' '=' ('TRUE' | 'FALSE'))? ('PATTERN' '=' STRING)? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' 'ICEBERG' 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME ('(' ID DATA_TYPE ('NOT' 'NULL')? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)? createIcebergTableParquet_inlineConstraint? (',' (ID DATA_TYPE | ID DATA_TYPE 'AS' EXPR) createIcebergTableParquet_inlineConstraint?)* (',' createIcebergTableParquet_outoflineConstraint)* ')')? ('PARTITION' 'BY' '(' ID (',' ID)* ')')? ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' STRING)? 'BASE_LOCATION' '=' STRING ('ENABLE_SCHEMA_INFERENCE' '=' ('TRUE' | 'FALSE'))? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('MATCH_BY_COLUMN_NAME' '=' ('CASE_INSENSITIVE' | 'CASE_SENSITIVE'))? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('REFRESH_ON_CREATE' '=' ('TRUE' | 'FALSE'))? ('PATTERN' '=' STRING)? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' '=' STRING)?
    ;
createIcebergTableParquet_inlineConstraint
    : ('CONSTRAINT' ID)? ('UNIQUE' | 'PRIMARY' 'KEY' | ('FOREIGN' 'KEY')? 'REFERENCES' QNAME ('(' ID ')')? ('MATCH' ('FULL' | 'SIMPLE' | 'PARTIAL'))? ('ON' ('UPDATE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))? ('DELETE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))?)?) ('NOT'? 'ENFORCED')? ('NOT'? 'DEFERRABLE')? ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? (('RELY' | 'NORELY'))?
    ;
createIcebergTableParquet_outoflineConstraint
    : ('CONSTRAINT' ID)? ('UNIQUE' '(' ID (',' ID)* ')' | 'PRIMARY' 'KEY' '(' ID (',' ID)* ')' | ('FOREIGN' 'KEY')? '(' ID (',' ID)* ')' 'REFERENCES' QNAME ('(' ID (',' ID)* ')')? ('MATCH' ('FULL' | 'SIMPLE' | 'PARTIAL'))? ('ON' ('UPDATE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))? ('DELETE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))?)?) ('NOT'? 'ENFORCED')? ('NOT'? 'DEFERRABLE')? ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? (('RELY' | 'NORELY'))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-rest
// fix: read: the additional columns take the whole column definition, which their line abbreviates
createIcebergTableRest
    : 'CREATE' ('OR' 'REPLACE')? 'ICEBERG' 'TABLE' QNAME ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' STRING)? 'CATALOG_TABLE_NAME' '=' STRING ('CATALOG_NAMESPACE' '=' STRING)? ('PATH_LAYOUT' '=' ('FLAT' | 'HIERARCHICAL'))? ('TARGET_FILE_SIZE' '=' ('\'AUTO\'' | '\'16MB\'' | '\'32MB\'' | '\'64MB\'' | '\'128MB\''))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    | 'CREATE' 'ICEBERG' 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' STRING)? 'CATALOG_TABLE_NAME' '=' STRING ('CATALOG_NAMESPACE' '=' STRING)? ('PATH_LAYOUT' '=' ('FLAT' | 'HIERARCHICAL'))? ('TARGET_FILE_SIZE' '=' ('\'AUTO\'' | '\'16MB\'' | '\'32MB\'' | '\'64MB\'' | '\'128MB\''))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'ICEBERG' 'TABLE' QNAME (ID DATA_TYPE ('DEFAULT' EXPR)? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? (',' ID DATA_TYPE ('DEFAULT' EXPR)? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)?)*)? ('PARTITION' 'BY' '(' createIcebergTableRest_partitionExpression (',' createIcebergTableRest_partitionExpression)* ')')? ('PATH_LAYOUT' '=' ('FLAT' | 'HIERARCHICAL'))? ('TARGET_FILE_SIZE' '=' ('\'AUTO\'' | '\'16MB\'' | '\'32MB\'' | '\'64MB\'' | '\'128MB\''))? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('COPY' 'GRANTS')? 'PURGE'? ('COMMENT' '=' STRING)? ('ICEBERG_VERSION' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('BASE_LOCATION' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))?
    | 'CREATE' 'ICEBERG' 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME (ID DATA_TYPE ('DEFAULT' EXPR)? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? (',' ID DATA_TYPE ('DEFAULT' EXPR)? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)?)*)? ('PARTITION' 'BY' '(' createIcebergTableRest_partitionExpression (',' createIcebergTableRest_partitionExpression)* ')')? ('PATH_LAYOUT' '=' ('FLAT' | 'HIERARCHICAL'))? ('TARGET_FILE_SIZE' '=' ('\'AUTO\'' | '\'16MB\'' | '\'32MB\'' | '\'64MB\'' | '\'128MB\''))? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('COPY' 'GRANTS')? 'PURGE'? ('COMMENT' '=' STRING)? ('ICEBERG_VERSION' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('BASE_LOCATION' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))?
    | 'CREATE' ('OR' 'REPLACE')? 'ICEBERG' 'TABLE' QNAME ('(' ID DATA_TYPE? (',' ID DATA_TYPE?)* ')')* 'PURGE'? 'AS' QUERY
    | 'CREATE' ('OR' 'REPLACE')? 'ICEBERG' 'TABLE' QNAME '(' ID DATA_TYPE 'WITH'? 'MASKING' 'POLICY' QNAME (',' QNAME)* (')')+ 'PURGE'? 'AS' QUERY
    ;
createIcebergTableRest_partitionExpression
    : ID
    | 'BUCKET' '(' NUMBER ',' ID ')'
    | 'TRUNCATE' '(' NUMBER ',' ID ')'
    | 'YEAR' '(' ID ')'
    | 'MONTH' '(' ID ')'
    | 'DAY' '(' ID ')'
    | 'HOUR' '(' ID ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-snowflake
// fix: read: Time Travel's OFFSET is a negative number of seconds, -N, and its STATEMENT a quoted query ID, as the Time Travel parameters say
// fix: read: <constraint_properties> is the properties CREATE TABLE's constraint page lists
// fix: read: <constraint_properties> is the properties CREATE TABLE's constraint page lists
// fix: read: an out-of-line key constraint names its columns, as the constraint page says
// fix: read: a foreign key's REFERENCES takes MATCH and ON UPDATE/DELETE, the properties the constraint page gives foreign keys
// fix: read: the additional columns take the whole column definition, which their line abbreviates
// not converted (form): ... stands for text the docs leave out
//   CREATE [ OR REPLACE ] ICEBERG TABLE [ IF NOT EXISTS ] <name>
//     CLONE <source_iceberg_table>
//       [ { AT | BEFORE } ( { TIMESTAMP => <timestamp> | OFFSET => -<seconds> | STATEMENT => '<query_id>' } ) ]
//       [COPY GRANTS]
//       [COPY TAGS]
//       [COPY = ( { COMMENTS | TAGS | GRANTS } [ , ... ] )]
//       ...
createIcebergTableSnowflake
    : 'CREATE' ('OR' 'REPLACE')? 'TRANSIENT'? 'ICEBERG' 'TABLE' QNAME '(' ID DATA_TYPE ('DEFAULT' EXPR)? createIcebergTableSnowflake_inlineConstraint? ('NOT' 'NULL')? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)? (',' ID DATA_TYPE ('DEFAULT' EXPR)? createIcebergTableSnowflake_inlineConstraint? ('NOT' 'NULL')? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)?)* (',' createIcebergTableSnowflake_outoflineConstraint)* ')' ('PARTITION' 'BY' '(' createIcebergTableSnowflake_partitionExpression (',' createIcebergTableSnowflake_partitionExpression)* ')')? ('PATH_LAYOUT' '=' ('FLAT' | 'HIERARCHICAL'))? ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')')? ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' '\'SNOWFLAKE\'')? ('BASE_LOCATION' '=' STRING)? ('TARGET_FILE_SIZE' '=' ('\'AUTO\'' | '\'16MB\'' | '\'32MB\'' | '\'64MB\'' | '\'128MB\''))? ('CATALOG_SYNC' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('CHANGE_TRACKING' '=' ('TRUE' | 'FALSE'))? ('COPY' 'GRANTS')? ('COPY' 'TAGS')? ('COPY' '=' '(' ('COMMENTS' | 'TAGS' | 'GRANTS') (',' ('COMMENTS' | 'TAGS' | 'GRANTS'))* ')')? ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('ICEBERG_VERSION' '=' NUMBER)? ('ALLOW_V3_EXPRESSION_DEFAULTS' '=' ('TRUE' | 'FALSE'))? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'AGGREGATION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))?
    | 'CREATE' 'TRANSIENT'? 'ICEBERG' 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME '(' ID DATA_TYPE ('DEFAULT' EXPR)? createIcebergTableSnowflake_inlineConstraint? ('NOT' 'NULL')? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)? (',' ID DATA_TYPE ('DEFAULT' EXPR)? createIcebergTableSnowflake_inlineConstraint? ('NOT' 'NULL')? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)?)* (',' createIcebergTableSnowflake_outoflineConstraint)* ')' ('PARTITION' 'BY' '(' createIcebergTableSnowflake_partitionExpression (',' createIcebergTableSnowflake_partitionExpression)* ')')? ('PATH_LAYOUT' '=' ('FLAT' | 'HIERARCHICAL'))? ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')')? ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' '\'SNOWFLAKE\'')? ('BASE_LOCATION' '=' STRING)? ('TARGET_FILE_SIZE' '=' ('\'AUTO\'' | '\'16MB\'' | '\'32MB\'' | '\'64MB\'' | '\'128MB\''))? ('CATALOG_SYNC' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('CHANGE_TRACKING' '=' ('TRUE' | 'FALSE'))? ('COPY' 'GRANTS')? ('COPY' 'TAGS')? ('COPY' '=' '(' ('COMMENTS' | 'TAGS' | 'GRANTS') (',' ('COMMENTS' | 'TAGS' | 'GRANTS'))* ')')? ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('ICEBERG_VERSION' '=' NUMBER)? ('ALLOW_V3_EXPRESSION_DEFAULTS' '=' ('TRUE' | 'FALSE'))? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'AGGREGATION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))?
    | 'CREATE' ('OR' 'REPLACE')? 'TRANSIENT'? 'ICEBERG' 'TABLE' QNAME ('(' ID DATA_TYPE? ('DEFAULT' EXPR)? (',' ID DATA_TYPE? ('DEFAULT' EXPR)?)* ')')? ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')')? ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' '\'SNOWFLAKE\'')? ('BASE_LOCATION' '=' STRING)? ('COPY' 'GRANTS')? ('COPY' 'TAGS')? ('COPY' '=' '(' ('COMMENTS' | 'TAGS' | 'GRANTS') (',' ('COMMENTS' | 'TAGS' | 'GRANTS'))* ')')? ('ICEBERG_VERSION' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))* 'AS' QUERY
    | 'CREATE' 'ICEBERG' 'TABLE' QNAME '(' ID DATA_TYPE ('DEFAULT' EXPR)? 'WITH'? 'MASKING' 'POLICY' QNAME (',' QNAME)* ')' ('EXTERNAL_VOLUME' '=' STRING)? ('CATALOG' '=' '\'SNOWFLAKE\'')? ('BASE_LOCATION' '=' STRING)? ('ICEBERG_VERSION' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? 'WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* (')')+ 'AS' QUERY
    | 'CREATE' ('OR' 'REPLACE')? 'TRANSIENT'? 'ICEBERG' 'TABLE' QNAME 'LIKE' QNAME ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')')? ('COPY' 'GRANTS')? ('COPY' 'TAGS')? ('COPY' '=' '(' ('COMMENTS' | 'TAGS' | 'GRANTS') (',' ('COMMENTS' | 'TAGS' | 'GRANTS'))* ')')?
    | 'CREATE' 'ICEBERG' 'TABLE' QNAME 'FROM' 'BACKUP' 'SET' ID 'IDENTIFIER' STRING
    ;
createIcebergTableSnowflake_inlineConstraint
    : ('CONSTRAINT' ID)? ('UNIQUE' | 'PRIMARY' 'KEY' | ('FOREIGN' 'KEY')? 'REFERENCES' QNAME ('(' ID ')')? ('MATCH' ('FULL' | 'SIMPLE' | 'PARTIAL'))? ('ON' ('UPDATE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))? ('DELETE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))?)? | 'CHECK' '(' EXPR ')') ('NOT'? 'ENFORCED')? ('NOT'? 'DEFERRABLE')? ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? (('RELY' | 'NORELY'))?
    ;
createIcebergTableSnowflake_outoflineConstraint
    : ('CONSTRAINT' ID)? ('UNIQUE' '(' ID (',' ID)* ')' | 'PRIMARY' 'KEY' '(' ID (',' ID)* ')' | ('FOREIGN' 'KEY')? '(' ID (',' ID)* ')' 'REFERENCES' QNAME ('(' ID (',' ID)* ')')? ('MATCH' ('FULL' | 'SIMPLE' | 'PARTIAL'))? ('ON' ('UPDATE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))? ('DELETE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))?)? | 'CHECK' '(' EXPR ')') ('NOT'? 'ENFORCED')? ('NOT'? 'DEFERRABLE')? ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? (('RELY' | 'NORELY'))?
    ;
createIcebergTableSnowflake_partitionExpression
    : ID
    | 'BUCKET' '(' NUMBER ',' ID ')'
    | 'TRUNCATE' '(' NUMBER ',' ID ')'
    | 'YEAR' '(' ID ')'
    | 'MONTH' '(' ID ')'
    | 'DAY' '(' ID ')'
    | 'HOUR' '(' ID ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-image-repository
createImageRepository
    : 'CREATE' ('OR' 'REPLACE')? 'IMAGE' 'REPOSITORY' QNAME ('ENCRYPTION' '=' '(' ('TYPE' '=' '\'SNOWFLAKE_FULL\'' | 'TYPE' '=' '\'SNOWFLAKE_SSE\'') ')')? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' 'IMAGE' 'REPOSITORY' ('IF' 'NOT' 'EXISTS')? QNAME ('ENCRYPTION' '=' '(' ('TYPE' '=' '\'SNOWFLAKE_FULL\'' | 'TYPE' '=' '\'SNOWFLAKE_SSE\'') ')')? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-index
createIndex
    : 'CREATE' ('OR' 'REPLACE')? 'INDEX' ID 'ON' QNAME '(' ID (',' ID)* ')' ('INCLUDE' '(' ID (',' ID)* ')')?
    | 'CREATE' 'INDEX' ('IF' 'NOT' 'EXISTS')? ID 'ON' QNAME '(' ID (',' ID)* ')' ('INCLUDE' '(' ID (',' ID)* ')')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-ingress-policy
createIngressPolicy
    : 'CREATE' ('OR' 'REPLACE')? 'INGRESS' 'POLICY' QNAME ('COMMENT' '=' STRING)? 'AS' STRING
    | 'CREATE' 'INGRESS' 'POLICY' ('IF' 'NOT' 'EXISTS')? QNAME ('COMMENT' '=' STRING)? 'AS' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-integration
// not converted (form): <integration_type> stands for a type's keyword, whose page gives its syntax
//   CREATE [ OR REPLACE ] <integration_type> INTEGRATION [ IF NOT EXISTS ] <object_name>
//     [ <integration_type_params> ]
//     [ COMMENT = '<string_literal>' ]

// https://docs.snowflake.com/en/sql-reference/sql/create-interactive-table
createInteractiveTable
    : 'CREATE' ('OR' 'REPLACE')? 'INTERACTIVE' 'TABLE' QNAME '(' ID DATA_TYPE ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? (',' ID DATA_TYPE)* ')' 'CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')' ('TARGET_LAG' '=' STRING)? ('WAREHOUSE' '=' ID)? ('INITIALIZATION_WAREHOUSE' '=' ID)? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'AGGREGATION' 'POLICY' QNAME ('ENTITY' 'KEY' '(' ID (',' ID)* ')')?)? ('WITH'? 'JOIN' 'POLICY' QNAME ('ALLOWED' 'JOIN' 'KEYS' '(' ID (',' ID)* ')')?)? ('WITH'? 'STORAGE' 'LIFECYCLE' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? 'AS' QUERY
    | 'CREATE' 'INTERACTIVE' 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME '(' ID DATA_TYPE ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? (',' ID DATA_TYPE)* ')' 'CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')' ('TARGET_LAG' '=' STRING)? ('WAREHOUSE' '=' ID)? ('INITIALIZATION_WAREHOUSE' '=' ID)? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'AGGREGATION' 'POLICY' QNAME ('ENTITY' 'KEY' '(' ID (',' ID)* ')')?)? ('WITH'? 'JOIN' 'POLICY' QNAME ('ALLOWED' 'JOIN' 'KEYS' '(' ID (',' ID)* ')')?)? ('WITH'? 'STORAGE' 'LIFECYCLE' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? 'AS' QUERY
    | 'CREATE' ('OR' 'REPLACE')? 'INTERACTIVE' 'TABLE' QNAME 'CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')' ('COMMENT' '=' STRING)? 'AS' QUERY
    | 'CREATE' ('OR' 'REPLACE')? 'INTERACTIVE' 'TABLE' QNAME 'CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')' 'TARGET_LAG' '=' STRING 'WAREHOUSE' '=' ID ('INITIALIZATION_WAREHOUSE' '=' ID)? ('COMMENT' '=' STRING)? 'AS' QUERY
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-interactive-warehouse
createInteractiveWarehouse
    : 'CREATE' ('OR' 'REPLACE')? 'INTERACTIVE' 'WAREHOUSE' QNAME ('TABLES' '(' QNAME (',' QNAME)* ')')? ('WITH'? createInteractiveWarehouse_objectProperties)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? createInteractiveWarehouse_objectParams?
    | 'CREATE' 'INTERACTIVE' 'WAREHOUSE' ('IF' 'NOT' 'EXISTS')? QNAME ('TABLES' '(' QNAME (',' QNAME)* ')')? ('WITH'? createInteractiveWarehouse_objectProperties)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? createInteractiveWarehouse_objectParams?
    ;
createInteractiveWarehouse_objectProperties
    : ('WAREHOUSE_SIZE' '=' ('XSMALL' | 'SMALL' | 'MEDIUM' | 'LARGE' | 'XLARGE' | 'XXLARGE' | 'XXXLARGE' | 'X4LARGE') | 'MAX_CLUSTER_COUNT' '=' NUMBER | 'MIN_CLUSTER_COUNT' '=' NUMBER | 'AUTO_SUSPEND' '=' (NUMBER | 'NULL') | 'AUTO_RESUME' '=' ('TRUE' | 'FALSE') | 'INITIALLY_SUSPENDED' '=' ('TRUE' | 'FALSE') | 'RESOURCE_MONITOR' '=' ID | 'FALLBACK_WAREHOUSE' '=' ID | 'COMMENT' '=' STRING)+
    ;
createInteractiveWarehouse_objectParams
    : ('MAX_CONCURRENCY_LEVEL' '=' NUMBER | 'STATEMENT_QUEUED_TIMEOUT_IN_SECONDS' '=' NUMBER | 'STATEMENT_TIMEOUT_IN_SECONDS' '=' NUMBER)+
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-join-policy
createJoinPolicy
    : 'CREATE' ('OR' 'REPLACE')? 'JOIN' 'POLICY' QNAME 'AS' '(' ')' 'RETURNS' 'JOIN_CONSTRAINT' '->' EXPR ('COMMENT' '=' STRING)?
    | 'CREATE' 'JOIN' 'POLICY' ('IF' 'NOT' 'EXISTS')? QNAME 'AS' '(' ')' 'RETURNS' 'JOIN_CONSTRAINT' '->' EXPR ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-listing
createListing
    : 'CREATE' 'EXTERNAL' 'LISTING' ('IF' 'NOT' 'EXISTS')? ID (('SHARE' ID | 'APPLICATION' 'PACKAGE' ID))? 'AS' STRING ('PUBLISH' '=' ('TRUE' | 'FALSE'))? ('REVIEW' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)?
    | 'CREATE' 'EXTERNAL' 'LISTING' ('IF' 'NOT' 'EXISTS')? ID (('SHARE' ID | 'APPLICATION' 'PACKAGE' ID))? 'FROM' STRING ('PUBLISH' '=' ('TRUE' | 'FALSE'))? ('REVIEW' '=' ('TRUE' | 'FALSE'))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-maintenance-policy
createMaintenancePolicy
    : 'CREATE' ('OR' 'REPLACE')? 'MAINTENANCE' 'POLICY' QNAME 'SCHEDULE' '=' STRING ('COMMENT' '=' STRING)?
    | 'CREATE' 'MAINTENANCE' 'POLICY' ('IF' 'NOT' 'EXISTS')? QNAME 'SCHEDULE' '=' STRING ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-managed-account
// fix: the comma before COMMENT goes with it, as the page's example ends at TYPE = READER
createManagedAccount
    : 'CREATE' 'MANAGED' 'ACCOUNT' QNAME 'ADMIN_NAME' '=' ID ',' 'ADMIN_PASSWORD' '=' STRING ',' 'TYPE' '=' 'READER' (',' 'COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-masking-policy
createMaskingPolicy
    : 'CREATE' ('OR' 'REPLACE')? 'MASKING' 'POLICY' QNAME 'AS' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')' 'RETURNS' DATA_TYPE '->' EXPR ('COMMENT' '=' STRING)? ('EXEMPT_OTHER_POLICIES' '=' ('TRUE' | 'FALSE'))?
    | 'CREATE' 'MASKING' 'POLICY' ('IF' 'NOT' 'EXISTS')? QNAME 'AS' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')' 'RETURNS' DATA_TYPE '->' EXPR ('COMMENT' '=' STRING)? ('EXEMPT_OTHER_POLICIES' '=' ('TRUE' | 'FALSE'))?
    | 'CREATE' 'OR' 'ALTER' 'MASKING' 'POLICY' QNAME 'AS' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')' 'RETURNS' DATA_TYPE '->' EXPR ('COMMENT' '=' STRING)? ('EXEMPT_OTHER_POLICIES' '=' ('TRUE' | 'FALSE'))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-materialized-view
// fix: the column list is optional and holds each column, with its policies, inside its parentheses, as the page's parameters and examples write them
createMaterializedView
    : 'CREATE' ('OR' 'REPLACE')? 'SECURE'? 'INTERACTIVE'? 'MATERIALIZED' 'VIEW' QNAME ('COPY' 'GRANTS')? ('COPY' '=' '(' 'GRANTS' ')')? ('(' ID ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? (',' ID ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?)* ')')? ('COMMENT' '=' STRING)? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'AGGREGATION' 'POLICY' QNAME ('ENTITY' 'KEY' '(' ID (',' ID)* ')')?)? ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? 'AS' QUERY
    | 'CREATE' 'SECURE'? 'INTERACTIVE'? 'MATERIALIZED' 'VIEW' ('IF' 'NOT' 'EXISTS')? QNAME ('COPY' 'GRANTS')? ('COPY' '=' '(' 'GRANTS' ')')? ('(' ID ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? (',' ID ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?)* ')')? ('COMMENT' '=' STRING)? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'AGGREGATION' 'POLICY' QNAME ('ENTITY' 'KEY' '(' ID (',' ID)* ')')?)? ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? 'AS' QUERY
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-mcp-server
createMcpServer
    : 'CREATE' ('OR' 'REPLACE')? 'MCP' 'SERVER' QNAME 'FROM' 'SPECIFICATION' STRING
    | 'CREATE' 'MCP' 'SERVER' ('IF' 'NOT' 'EXISTS')? QNAME 'FROM' 'SPECIFICATION' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-model
createModel
    : 'CREATE' ('OR' 'REPLACE')? 'MODEL' QNAME ('WITH' 'VERSION' ID)? 'FROM' 'MODEL' QNAME ('VERSION' ID)?
    | 'CREATE' 'MODEL' ('IF' 'NOT' 'EXISTS')? QNAME ('WITH' 'VERSION' ID)? 'FROM' 'MODEL' QNAME ('VERSION' ID)?
    | 'CREATE' ('OR' 'REPLACE')? 'MODEL' QNAME ('WITH' 'VERSION' ID)? 'FROM' createModel_internalStage
    | 'CREATE' 'MODEL' ('IF' 'NOT' 'EXISTS')? QNAME ('WITH' 'VERSION' ID)? 'FROM' createModel_internalStage
    ;
createModel_internalStage
    : STAGE
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-model-monitor
// fix: read: a <..._array> of column names is ( '<name>' [ , ... ] ), as the examples write it
// fix: read: a <..._array> of column names is ( '<name>' [ , ... ] ), as the examples write it
createModelMonitor
    : 'CREATE' ('OR' 'REPLACE')? 'MODEL' 'MONITOR' ID 'WITH' 'MODEL' '=' QNAME 'VERSION' '=' STRING 'FUNCTION' '=' STRING 'SOURCE' '=' ID 'WAREHOUSE' '=' ID 'REFRESH_INTERVAL' '=' STRING 'AGGREGATION_WINDOW' '=' STRING 'TIMESTAMP_COLUMN' '=' ID ('BASELINE' '=' ID)? ('ID_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('PREDICTION_CLASS_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('PREDICTION_SCORE_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('ACTUAL_CLASS_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('ACTUAL_SCORE_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('SEGMENT_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('CUSTOM_METRIC_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' 'MODEL' 'MONITOR' ('IF' 'NOT' 'EXISTS')? ID 'WITH' 'MODEL' '=' QNAME 'VERSION' '=' STRING 'FUNCTION' '=' STRING 'SOURCE' '=' ID 'WAREHOUSE' '=' ID 'REFRESH_INTERVAL' '=' STRING 'AGGREGATION_WINDOW' '=' STRING 'TIMESTAMP_COLUMN' '=' ID ('BASELINE' '=' ID)? ('ID_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('PREDICTION_CLASS_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('PREDICTION_SCORE_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('ACTUAL_CLASS_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('ACTUAL_SCORE_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('SEGMENT_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('CUSTOM_METRIC_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' ('OR' 'REPLACE')? 'MODEL' 'MONITOR' ID 'WITH' 'MODEL' '=' QNAME 'GATEWAY' '=' QNAME 'FUNCTION' '=' STRING 'WAREHOUSE' '=' ID 'REFRESH_INTERVAL' '=' STRING 'AGGREGATION_WINDOW' '=' STRING ('GROUND_TRUTH' '=' QNAME)? ('ID_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('PREDICTION_CLASS_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('PREDICTION_SCORE_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('ACTUAL_CLASS_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('ACTUAL_SCORE_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' 'MODEL' 'MONITOR' ('IF' 'NOT' 'EXISTS')? ID 'WITH' 'MODEL' '=' QNAME 'GATEWAY' '=' QNAME 'FUNCTION' '=' STRING 'WAREHOUSE' '=' ID 'REFRESH_INTERVAL' '=' STRING 'AGGREGATION_WINDOW' '=' STRING ('GROUND_TRUTH' '=' QNAME)? ('ID_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('PREDICTION_CLASS_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('PREDICTION_SCORE_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('ACTUAL_CLASS_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('ACTUAL_SCORE_COLUMNS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-multi-party-approval-policy
createMultiPartyApprovalPolicy
    : 'CREATE' ('OR' 'REPLACE')? 'MULTI' 'PARTY' 'APPROVAL' 'POLICY' QNAME 'AS' STRING
    | 'CREATE' 'MULTI' 'PARTY' 'APPROVAL' 'POLICY' ('IF' 'NOT' 'EXISTS')? QNAME 'AS' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-network-policy
createNetworkPolicy
    : 'CREATE' ('OR' 'REPLACE')? 'NETWORK' 'POLICY' ID ('ALLOWED_NETWORK_RULE_LIST' '=' '(' STRING (',' STRING)* ')')? ('BLOCKED_NETWORK_RULE_LIST' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_IP_LIST' '=' '(' (STRING (',' STRING)*)? ')')? ('BLOCKED_IP_LIST' '=' '(' (STRING (',' STRING)*)? ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' 'NETWORK' 'POLICY' ('IF' 'NOT' 'EXISTS')? ID ('ALLOWED_NETWORK_RULE_LIST' '=' '(' STRING (',' STRING)* ')')? ('BLOCKED_NETWORK_RULE_LIST' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_IP_LIST' '=' '(' (STRING (',' STRING)*)? ')')? ('BLOCKED_IP_LIST' '=' '(' (STRING (',' STRING)*)? ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' 'OR' 'ALTER' 'NETWORK' 'POLICY' ID ('ALLOWED_NETWORK_RULE_LIST' '=' '(' STRING (',' STRING)* ')')? ('BLOCKED_NETWORK_RULE_LIST' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_IP_LIST' '=' '(' (STRING (',' STRING)*)? ')')? ('BLOCKED_IP_LIST' '=' '(' (STRING (',' STRING)*)? ')')? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-network-rule
createNetworkRule
    : 'CREATE' ('OR' 'REPLACE')? 'NETWORK' 'RULE' QNAME 'TYPE' '=' ('IPV4' | 'IPV6' | 'AWSVPCEID' | 'AZURELINKID' | 'GCPPSCID' | 'HOST_PORT' | 'PRIVATE_HOST_PORT' | 'COMPUTE_POOL') 'VALUE_LIST' '=' '(' STRING (',' STRING)* ')' ('MODE' '=' ('INGRESS' | 'INTERNAL_STAGE' | 'SNOWFLAKE_MANAGED_STORAGE_VOLUME' | 'EGRESS'))? ('COMMENT' '=' STRING)?
    | 'CREATE' 'OR' 'ALTER' 'NETWORK' 'RULE' QNAME 'TYPE' '=' ('IPV4' | 'IPV6' | 'AWSVPCEID' | 'AZURELINKID' | 'GCPPSCID' | 'HOST_PORT' | 'PRIVATE_HOST_PORT' | 'COMPUTE_POOL') 'VALUE_LIST' '=' '(' STRING (',' STRING)* ')' 'MODE' '=' ('INGRESS' | 'INTERNAL_STAGE' | 'SNOWFLAKE_MANAGED_STORAGE_VOLUME' | 'EGRESS') ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-notebook
createNotebook
    : 'CREATE' ('OR' 'REPLACE')? 'NOTEBOOK' QNAME ('FROM' STRING)? ('MAIN_FILE' '=' STRING)? ('COMMENT' '=' STRING)? ('QUERY_WAREHOUSE' '=' ID)? ('IDLE_AUTO_SHUTDOWN_TIME_SECONDS' '=' NUMBER)? ('RUNTIME_NAME' '=' STRING)? ('COMPUTE_POOL' '=' STRING)? ('WAREHOUSE' '=' ID)?
    | 'CREATE' 'NOTEBOOK' ('IF' 'NOT' 'EXISTS')? QNAME ('FROM' STRING)? ('MAIN_FILE' '=' STRING)? ('COMMENT' '=' STRING)? ('QUERY_WAREHOUSE' '=' ID)? ('IDLE_AUTO_SHUTDOWN_TIME_SECONDS' '=' NUMBER)? ('RUNTIME_NAME' '=' STRING)? ('COMPUTE_POOL' '=' STRING)? ('WAREHOUSE' '=' ID)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-notebook-project
createNotebookProject
    : 'CREATE' 'NOTEBOOK' 'PROJECT' ID '.' ID '.' ID 'FROM' STRING ('COMMENT' '=' STRING)?
    | 'CREATE' 'NOTEBOOK' 'PROJECT' ('IF' 'NOT' 'EXISTS')? ID '.' ID '.' ID 'FROM' STRING ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-code-bundle
// not converted (form): DEFINE runs only in DCM project files, not as a statement
//   DEFINE CODE BUNDLE DEMO{{env_suffix}}.JOBS.MY_JOB
//     FROM 'asset://my_job'
//     COMMENT = 'Python job managed by DCM';
createCodeBundle
    : 'CREATE' ('OR' 'REPLACE')? 'CODE' 'BUNDLE' ID '.' ID '.' ID 'FROM' STRING ('COMMENT' '=' STRING)?
    | 'CREATE' 'CODE' 'BUNDLE' ('IF' 'NOT' 'EXISTS')? ID '.' ID '.' ID 'FROM' STRING ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration
// no syntax

// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-email
// not converted (form): ... stands for text the docs leave out
//   CREATE [ OR REPLACE ] NOTIFICATION INTEGRATION [ IF NOT EXISTS ] <name>
//     TYPE = EMAIL
//     ENABLED = { TRUE | FALSE }
//     [ ALLOWED_RECIPIENTS = ( '<email_address>' [ , ... '<email_address>' ] ) ]
//     [ DEFAULT_RECIPIENTS = ( '<email_address>' [ , ... '<email_address>' ] ) ]
//     [ DEFAULT_SUBJECT = '<subject_line>' ]
//     [ COMMENT = '<string_literal>' ]

// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-multi-queue
// fix: the block after the statement defines <queue>, one variant per cloud provider; named here
createNotificationIntegrationMultiQueue
    : 'CREATE' ('OR' 'REPLACE')? 'NOTIFICATION' 'INTEGRATION' ID 'ENABLED' '=' ('TRUE' | 'FALSE') 'TYPE' '=' 'MULTI_QUEUE' 'DIRECTION' '=' 'INBOUND' 'QUEUES' '=' '(' createNotificationIntegrationMultiQueue_queue (',' createNotificationIntegrationMultiQueue_queue)? ')' 'ACTIVE' '=' STRING ('COMMENT' '=' STRING)?
    | 'CREATE' 'NOTIFICATION' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'ENABLED' '=' ('TRUE' | 'FALSE') 'TYPE' '=' 'MULTI_QUEUE' 'DIRECTION' '=' 'INBOUND' 'QUEUES' '=' '(' createNotificationIntegrationMultiQueue_queue (',' createNotificationIntegrationMultiQueue_queue)? ')' 'ACTIVE' '=' STRING ('COMMENT' '=' STRING)?
    ;
createNotificationIntegrationMultiQueue_queue
    : '(' 'NAME' '=' STRING 'NOTIFICATION_PROVIDER' '=' 'AWS_SNS' 'AWS_SNS_TOPIC_ARN' '=' STRING ')'
    | '(' 'NAME' '=' STRING 'NOTIFICATION_PROVIDER' '=' 'GCP_PUBSUB' 'GCP_PUBSUB_SUBSCRIPTION_NAME' '=' STRING ')'
    | '(' 'NAME' '=' STRING 'NOTIFICATION_PROVIDER' '=' 'AZURE_STORAGE_QUEUE' 'AZURE_STORAGE_QUEUE_PRIMARY_URI' '=' STRING 'AZURE_TENANT_ID' '=' STRING ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-queue-inbound-azure
// fix: a ; ends the statement before the parameters the page lists after it
createNotificationIntegrationQueueInboundAzure
    : 'CREATE' ('OR' 'REPLACE')? 'NOTIFICATION' 'INTEGRATION' ID 'ENABLED' '=' ('TRUE' | 'FALSE') 'TYPE' '=' 'QUEUE' 'NOTIFICATION_PROVIDER' '=' 'AZURE_STORAGE_QUEUE' 'AZURE_STORAGE_QUEUE_PRIMARY_URI' '=' STRING 'AZURE_TENANT_ID' '=' STRING ('USE_PRIVATELINK_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)?
    | 'CREATE' 'NOTIFICATION' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'ENABLED' '=' ('TRUE' | 'FALSE') 'TYPE' '=' 'QUEUE' 'NOTIFICATION_PROVIDER' '=' 'AZURE_STORAGE_QUEUE' 'AZURE_STORAGE_QUEUE_PRIMARY_URI' '=' STRING 'AZURE_TENANT_ID' '=' STRING ('USE_PRIVATELINK_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-queue-inbound-gcp
createNotificationIntegrationQueueInboundGcp
    : 'CREATE' ('OR' 'REPLACE')? 'NOTIFICATION' 'INTEGRATION' ID 'ENABLED' '=' ('TRUE' | 'FALSE') 'TYPE' '=' 'QUEUE' 'NOTIFICATION_PROVIDER' '=' 'GCP_PUBSUB' 'GCP_PUBSUB_SUBSCRIPTION_NAME' '=' STRING ('COMMENT' '=' STRING)?
    | 'CREATE' 'NOTIFICATION' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'ENABLED' '=' ('TRUE' | 'FALSE') 'TYPE' '=' 'QUEUE' 'NOTIFICATION_PROVIDER' '=' 'GCP_PUBSUB' 'GCP_PUBSUB_SUBSCRIPTION_NAME' '=' STRING ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-queue-outbound-aws
createNotificationIntegrationQueueOutboundAws
    : 'CREATE' ('OR' 'REPLACE')? 'NOTIFICATION' 'INTEGRATION' ID 'ENABLED' '=' ('TRUE' | 'FALSE') 'TYPE' '=' 'QUEUE' 'DIRECTION' '=' 'OUTBOUND' 'NOTIFICATION_PROVIDER' '=' 'AWS_SNS' 'AWS_SNS_TOPIC_ARN' '=' STRING 'AWS_SNS_ROLE_ARN' '=' STRING ('COMMENT' '=' STRING)?
    | 'CREATE' 'NOTIFICATION' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'ENABLED' '=' ('TRUE' | 'FALSE') 'TYPE' '=' 'QUEUE' 'DIRECTION' '=' 'OUTBOUND' 'NOTIFICATION_PROVIDER' '=' 'AWS_SNS' 'AWS_SNS_TOPIC_ARN' '=' STRING 'AWS_SNS_ROLE_ARN' '=' STRING ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-queue-outbound-azure
// fix: a ; ends the statement before the parameters the page lists after it
createNotificationIntegrationQueueOutboundAzure
    : 'CREATE' ('OR' 'REPLACE')? 'NOTIFICATION' 'INTEGRATION' ID 'ENABLED' '=' ('TRUE' | 'FALSE') 'TYPE' '=' 'QUEUE' 'DIRECTION' '=' 'OUTBOUND' 'NOTIFICATION_PROVIDER' '=' 'AZURE_EVENT_GRID' 'AZURE_EVENT_GRID_TOPIC_ENDPOINT' '=' STRING 'AZURE_TENANT_ID' '=' STRING ('COMMENT' '=' STRING)?
    | 'CREATE' 'NOTIFICATION' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'ENABLED' '=' ('TRUE' | 'FALSE') 'TYPE' '=' 'QUEUE' 'DIRECTION' '=' 'OUTBOUND' 'NOTIFICATION_PROVIDER' '=' 'AZURE_EVENT_GRID' 'AZURE_EVENT_GRID_TOPIC_ENDPOINT' '=' STRING 'AZURE_TENANT_ID' '=' STRING ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-queue-outbound-gcp
createNotificationIntegrationQueueOutboundGcp
    : 'CREATE' ('OR' 'REPLACE')? 'NOTIFICATION' 'INTEGRATION' ID 'ENABLED' '=' ('TRUE' | 'FALSE') 'TYPE' '=' 'QUEUE' 'DIRECTION' '=' 'OUTBOUND' 'NOTIFICATION_PROVIDER' '=' 'GCP_PUBSUB' 'GCP_PUBSUB_TOPIC_NAME' '=' STRING ('COMMENT' '=' STRING)?
    | 'CREATE' 'NOTIFICATION' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'ENABLED' '=' ('TRUE' | 'FALSE') 'TYPE' '=' 'QUEUE' 'DIRECTION' '=' 'OUTBOUND' 'NOTIFICATION_PROVIDER' '=' 'GCP_PUBSUB' 'GCP_PUBSUB_TOPIC_NAME' '=' STRING ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-webhooks
createNotificationIntegrationWebhooks
    : 'CREATE' ('OR' 'REPLACE')? 'NOTIFICATION' 'INTEGRATION' ID 'TYPE' '=' 'WEBHOOK' 'ENABLED' '=' ('TRUE' | 'FALSE') 'WEBHOOK_URL' '=' STRING ('WEBHOOK_SECRET' '=' QNAME)? ('WEBHOOK_BODY_TEMPLATE' '=' STRING)? ('WEBHOOK_HEADERS' '=' '(' STRING '=' STRING (',' STRING '=' STRING)* ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' 'NOTIFICATION' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'TYPE' '=' 'WEBHOOK' 'ENABLED' '=' ('TRUE' | 'FALSE') 'WEBHOOK_URL' '=' STRING ('WEBHOOK_SECRET' '=' QNAME)? ('WEBHOOK_BODY_TEMPLATE' '=' STRING)? ('WEBHOOK_HEADERS' '=' '(' STRING '=' STRING (',' STRING '=' STRING)* ')')? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-openflow-connector
createOpenflowConnector
    : 'CREATE' 'OPENFLOW' 'CONNECTOR' ('IF' 'NOT' 'EXISTS')? QNAME 'IN' 'RUNTIME' ID 'FROM' 'DEFINITION' ID ('DISPLAY_NAME' '=' STRING)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'OPENFLOW' 'CONNECTOR' ('IF' 'NOT' 'EXISTS')? QNAME 'IN' 'RUNTIME' ID 'FROM' STRING ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-openflow-deployment
createOpenflowDeployment
    : 'CREATE' 'OPENFLOW' 'DEPLOYMENT' ('IF' 'NOT' 'EXISTS')? QNAME ('DEPLOYMENT_TYPE' '=' ('SNOWFLAKE' | 'BYOC'))? ('VPC_TYPE' '=' ('\'MANAGED\'' | '\'PROVIDED\''))? ('USE_PRIVATE_LINK' '=' ('TRUE' | 'FALSE'))? ('USE_USER_AUTH_OVER_PRIVATE_LINK' '=' ('TRUE' | 'FALSE'))? ('CUSTOM_INGRESS_HOSTNAME' '=' STRING)? ('DISPLAY_NAME' '=' STRING)? ('COMMENT' '=' STRING)? ('EVENT_TABLE' '=' (STRING | 'NONE'))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-openflow-runtime
createOpenflowRuntime
    : 'CREATE' 'OPENFLOW' 'RUNTIME' ('IF' 'NOT' 'EXISTS')? (ID '.' ID '.')? ID 'IN' 'DEPLOYMENT' ID 'NODE_TYPE' '=' ('SMALL' | 'MEDIUM' | 'LARGE') ('NODE_TYPE_TIER' '=' ('\'S1\'' | '\'S2\'' | '\'S3\'' | '\'M4\'' | '\'M6\'' | '\'L8\''))? 'MIN_NODES' '=' NUMBER 'MAX_NODES' '=' NUMBER 'EXECUTE_AS_ROLE' '=' ID ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('DISPLAY_NAME' '=' STRING)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-or-alter
// no syntax

// https://docs.snowflake.com/en/sql-reference/sql/create-organization-account
createOrganizationAccount
    : 'CREATE' 'ORGANIZATION' 'ACCOUNT' ID 'ADMIN_NAME' '=' STRING ('ADMIN_PASSWORD' '=' STRING | 'ADMIN_RSA_PUBLIC_KEY' '=' STRING) ('FIRST_NAME' '=' STRING)? ('LAST_NAME' '=' STRING)? 'EMAIL' '=' STRING ('MUST_CHANGE_PASSWORD' '=' ('TRUE' | 'FALSE'))? 'EDITION' '=' ('ENTERPRISE' | 'BUSINESS_CRITICAL') ('REGION_GROUP' '=' ID)? ('REGION' '=' ID)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-organization-listing
createOrganizationListing
    : 'CREATE' 'ORGANIZATION' 'LISTING' ('IF' 'NOT' 'EXISTS')? ID (('SHARE' ID | 'APPLICATION' 'PACKAGE' ID))? 'AS' STRING ('PUBLISH' '=' ('TRUE' | 'FALSE'))?
    | 'CREATE' 'ORGANIZATION' 'LISTING' ('IF' 'NOT' 'EXISTS')? ID (('SHARE' ID | 'APPLICATION' 'PACKAGE' ID))? 'FROM' STRING ('PUBLISH' '=' ('TRUE' | 'FALSE'))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-organization-profile
createOrganizationProfile
    : 'CREATE' 'ORGANIZATION' 'PROFILE' ('IF' 'NOT' 'EXISTS')? ID
    | 'CREATE' 'ORGANIZATION' 'PROFILE' ('IF' 'NOT' 'EXISTS')? ID 'AS' STRING ('VERSION' ID)? ('PUBLISH' '=' ('TRUE' | 'FALSE'))?
    | 'CREATE' 'ORGANIZATION' 'PROFILE' ('IF' 'NOT' 'EXISTS')? ID 'FROM' STAGE ('VERSION' ID)? ('PUBLISH' '=' ('TRUE' | 'FALSE'))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-organization-user
createOrganizationUser
    : 'CREATE' 'ORGANIZATION' 'USER' ('IF' 'NOT' 'EXISTS')? ID createOrganizationUser_objectProperties?
    ;
createOrganizationUser_objectProperties
    : ('EMAIL' '=' STRING | 'LOGIN_NAME' '=' STRING | 'DISPLAY_NAME' '=' STRING | 'FIRST_NAME' '=' STRING | 'MIDDLE_NAME' '=' STRING | 'LAST_NAME' '=' STRING | 'TYPE' '=' ('PERSON' | 'SERVICE') | 'COMMENT' '=' STRING)+
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-organization-user-group
createOrganizationUserGroup
    : 'CREATE' 'ORGANIZATION' 'USER' 'GROUP' ('IF' 'NOT' 'EXISTS')? ID ('IS_GRANTABLE' '=' ('TRUE' | 'FALSE'))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-packages-policy
createPackagesPolicy
    : 'CREATE' ('OR' 'REPLACE')? 'PACKAGES' 'POLICY' QNAME 'LANGUAGE' 'PYTHON' ('ALLOWLIST' '=' '(' (STRING (',' STRING)*)? ')')? ('BLOCKLIST' '=' '(' (STRING (',' STRING)*)? ')')? ('ADDITIONAL_CREATION_BLOCKLIST' '=' '(' (STRING (',' STRING)*)? ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' 'PACKAGES' 'POLICY' ('IF' 'NOT' 'EXISTS')? QNAME 'LANGUAGE' 'PYTHON' ('ALLOWLIST' '=' '(' (STRING (',' STRING)*)? ')')? ('BLOCKLIST' '=' '(' (STRING (',' STRING)*)? ')')? ('ADDITIONAL_CREATION_BLOCKLIST' '=' '(' (STRING (',' STRING)*)? ')')? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-password-policy
createPasswordPolicy
    : 'CREATE' ('OR' 'REPLACE')? 'PASSWORD' 'POLICY' QNAME ('PASSWORD_MIN_LENGTH' '=' NUMBER)? ('PASSWORD_MAX_LENGTH' '=' NUMBER)? ('PASSWORD_MIN_UPPER_CASE_CHARS' '=' NUMBER)? ('PASSWORD_MIN_LOWER_CASE_CHARS' '=' NUMBER)? ('PASSWORD_MIN_NUMERIC_CHARS' '=' NUMBER)? ('PASSWORD_MIN_SPECIAL_CHARS' '=' NUMBER)? ('PASSWORD_MIN_AGE_DAYS' '=' NUMBER)? ('PASSWORD_MAX_AGE_DAYS' '=' NUMBER)? ('PASSWORD_MAX_RETRIES' '=' NUMBER)? ('PASSWORD_LOCKOUT_TIME_MINS' '=' NUMBER)? ('PASSWORD_HISTORY' '=' NUMBER)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'PASSWORD' 'POLICY' ('IF' 'NOT' 'EXISTS')? QNAME ('PASSWORD_MIN_LENGTH' '=' NUMBER)? ('PASSWORD_MAX_LENGTH' '=' NUMBER)? ('PASSWORD_MIN_UPPER_CASE_CHARS' '=' NUMBER)? ('PASSWORD_MIN_LOWER_CASE_CHARS' '=' NUMBER)? ('PASSWORD_MIN_NUMERIC_CHARS' '=' NUMBER)? ('PASSWORD_MIN_SPECIAL_CHARS' '=' NUMBER)? ('PASSWORD_MIN_AGE_DAYS' '=' NUMBER)? ('PASSWORD_MAX_AGE_DAYS' '=' NUMBER)? ('PASSWORD_MAX_RETRIES' '=' NUMBER)? ('PASSWORD_LOCKOUT_TIME_MINS' '=' NUMBER)? ('PASSWORD_HISTORY' '=' NUMBER)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-pipe
// fix: AUTO_INGEST takes TRUE or FALSE, as the page's parameters say
createPipe
    : 'CREATE' ('OR' 'REPLACE')? 'PIPE' QNAME ('AUTO_INGEST' '=' ('TRUE' | 'FALSE'))? ('ERROR_INTEGRATION' '=' ID)? ('AWS_SNS_TOPIC' '=' STRING)? ('INTEGRATION' '=' STRING)? ('COMMENT' '=' STRING)? 'AS' COPY_STATEMENT
    | 'CREATE' 'PIPE' ('IF' 'NOT' 'EXISTS')? QNAME ('AUTO_INGEST' '=' ('TRUE' | 'FALSE'))? ('ERROR_INTEGRATION' '=' ID)? ('AWS_SNS_TOPIC' '=' STRING)? ('INTEGRATION' '=' STRING)? ('COMMENT' '=' STRING)? 'AS' COPY_STATEMENT
    | 'CREATE' 'OR' 'ALTER' 'PIPE' QNAME ('AUTO_INGEST' '=' ('TRUE' | 'FALSE'))? ('ERROR_INTEGRATION' '=' ID)? ('AWS_SNS_TOPIC' '=' STRING)? ('INTEGRATION' '=' STRING)? ('PIPE_EXECUTION_PAUSED' '=' 'TRUE' | 'PIPE_EXECUTION_PAUSED' '=' 'FALSE')? ('COMMENT' '=' STRING)? 'AS' COPY_STATEMENT
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-postgres-instance
// fix: the properties are ISSUER, AUDIENCE and ACCEPTED_TOKENS, each token MAPPING_CLAIM, SCOPES_CLAIM and SCOPES, as the parameters and example write them; SCOPES is a list or one string
// fix: read: Time Travel's OFFSET is a negative number of seconds, -N, and its STATEMENT a quoted query ID, as the Time Travel parameters say
createPostgresInstance
    : 'CREATE' 'POSTGRES' 'INSTANCE' ('IF' 'NOT' 'EXISTS')? QNAME 'COMPUTE_FAMILY' '=' STRING 'STORAGE_SIZE_GB' '=' NUMBER 'AUTHENTICATION_AUTHORITY' '=' ('POSTGRES' | 'POSTGRES_OR_SNOWFLAKE' | 'EXTERNAL_OAUTH') ('AUTHENTICATION_PROPERTIES' '=' '(' 'ISSUER' '=' STRING ',' 'AUDIENCE' '=' STRING ',' 'ACCEPTED_TOKENS' '=' '(' '(' 'MAPPING_CLAIM' '=' STRING ',' 'SCOPES_CLAIM' '=' STRING ',' 'SCOPES' '=' ('(' STRING (',' STRING)* ')' | STRING) ')' (',' '(' 'MAPPING_CLAIM' '=' STRING ',' 'SCOPES_CLAIM' '=' STRING ',' 'SCOPES' '=' ('(' STRING (',' STRING)* ')' | STRING) ')')* ')' ')')? ('POSTGRES_VERSION' '=' ('16' | '17' | '18'))? ('NETWORK_POLICY' '=' STRING)? ('HIGH_AVAILABILITY' '=' ('TRUE' | 'FALSE'))? ('STORAGE_INTEGRATION' '=' STRING)? ('POSTGRES_SETTINGS' '=' STRING)? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' 'POSTGRES' 'INSTANCE' ('IF' 'NOT' 'EXISTS')? QNAME 'FORK' ID (('AT' | 'BEFORE') '(' ('TIMESTAMP' '=>' STRING | 'OFFSET' '=>' '-' NUMBER) ')')? ('COMPUTE_FAMILY' '=' STRING)? ('STORAGE_SIZE_GB' '=' NUMBER)? ('HIGH_AVAILABILITY' '=' ('TRUE' | 'FALSE'))? ('POSTGRES_SETTINGS' '=' STRING)? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-privacy-policy
createPrivacyPolicy
    : 'CREATE' ('OR' 'REPLACE')? 'PRIVACY' 'POLICY' QNAME 'AS' '(' ')' 'RETURNS' 'PRIVACY_BUDGET' '->' EXPR ('COMMENT' '=' STRING)?
    | 'CREATE' 'PRIVACY' 'POLICY' ('IF' 'NOT' 'EXISTS')? QNAME 'AS' '(' ')' 'RETURNS' 'PRIVACY_BUDGET' '->' EXPR ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-procedure
// not converted (form): unexpected '`'
//   CREATE [ OR REPLACE ] [ { TEMP | TEMPORARY } ] [ SECURE ] PROCEDURE <name> (
//       [ <arg_name> <arg_data_type> [ DEFAULT <default_value> ] ] [ , ... ] )
//     [ COPY GRANTS ]
//     RETURNS { <result_data_type> [ [ NOT ] NULL ] | TABLE ( [ <col_name> <col_data_type> [ , ... ] ] ) }
//     LANGUAGE PYTHON
//     RUNTIME_VERSION = '<python_version>'
//     [ ARTIFACT_REPOSITORY = `<repository_name>` ]
//     [ PACKAGES = ( '<package_name>' [ , ... ] ) ]
//     [ IMPORTS = ( '<stage_path_and_directory_or_file_name_to_read>' [, '<stage_path_and_directory_or_file_name_to_read>' ...] ) ]
//     HANDLER = '<function_name>'
//     [ EXTERNAL_ACCESS_INTEGRATIONS = ( <name_of_integration> [ , ... ] ) ]
//     [ SECRETS = ('<secret_variable_name>' = <secret_name> [ , ... ] ) ]
//     [ { CALLED ON NULL INPUT | { RETURNS NULL ON NULL INPUT | STRICT } } ]
//     [ { VOLATILE | IMMUTABLE } ] -- Note: VOLATILE and IMMUTABLE are deprecated.
//     [ COMMENT = '<string_literal>' ]
//     [ COPY = ( GRANTS ) ]
//     [ EXECUTE AS { OWNER | CALLER | RESTRICTED CALLER }]
//     AS '<procedure_definition>'
// not converted (form): unexpected '`'
//   CREATE [ OR REPLACE ] [ { TEMP | TEMPORARY } ] PROCEDURE <name> (
//       [ <arg_name> <arg_data_type> [ DEFAULT <default_value> ] ] [ , ... ] )
//     [ COPY GRANTS ]
//     RETURNS { <result_data_type> [ [ NOT ] NULL ] | TABLE ( [ <col_name> <col_data_type> [ , ... ] ] ) }
//     LANGUAGE PYTHON
//     RUNTIME_VERSION = '<python_version>'
//     [ ARTIFACT_REPOSITORY = `<repository_name>` ]
//     [ PACKAGES = ( '<package_name>' [ , ... ] ) ]
//     [ IMPORTS = ( '<stage_path_and_directory_or_file_name_to_read>' [, '<stage_path_and_directory_or_file_name_to_read>' ...] ) ]
//     HANDLER = '<module_file_name>.<function_name>'
//     [ EXTERNAL_ACCESS_INTEGRATIONS = ( <name_of_integration> [ , ... ] ) ]
//     [ SECRETS = ('<secret_variable_name>' = <secret_name> [ , ... ] ) ]
//     [ { CALLED ON NULL INPUT | { RETURNS NULL ON NULL INPUT | STRICT } } ]
//     [ { VOLATILE | IMMUTABLE } ] -- Note: VOLATILE and IMMUTABLE are deprecated.
//     [ COMMENT = '<string_literal>' ]
//     [ COPY = ( GRANTS ) ]
//     [ EXECUTE AS { OWNER | CALLER | RESTRICTED CALLER } ]
// not converted (form): ... stands for text the docs leave out
//   CREATE [ OR ALTER ] PROCEDURE ...
createProcedure
    : 'CREATE' ('OR' 'REPLACE')? (('TEMP' | 'TEMPORARY'))? 'SECURE'? 'PROCEDURE' QNAME '(' (ID DATA_TYPE ('DEFAULT' EXPR)? (',' ID DATA_TYPE ('DEFAULT' EXPR)?)*)? ')' ('COPY' 'GRANTS')? 'RETURNS' (DATA_TYPE ('NOT'? 'NULL')? | 'TABLE' '(' (ID DATA_TYPE (',' ID DATA_TYPE)*)? ')') 'LANGUAGE' 'JAVA' 'RUNTIME_VERSION' '=' STRING 'PACKAGES' '=' '(' STRING (',' STRING)* ')' ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? 'HANDLER' '=' STRING ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')')? ('TARGET_PATH' '=' STRING)? (('CALLED' 'ON' 'NULL' 'INPUT' | ('RETURNS' 'NULL' 'ON' 'NULL' 'INPUT' | 'STRICT')))? (('VOLATILE' | 'IMMUTABLE'))? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')? ('EXECUTE' 'AS' ('OWNER' | 'CALLER' | 'RESTRICTED' 'CALLER'))? 'AS' STRING
    | 'CREATE' ('OR' 'REPLACE')? (('TEMP' | 'TEMPORARY'))? 'SECURE'? 'PROCEDURE' QNAME '(' (ID DATA_TYPE ('DEFAULT' EXPR)? (',' ID DATA_TYPE ('DEFAULT' EXPR)?)*)? ')' ('COPY' 'GRANTS')? 'RETURNS' (DATA_TYPE ('NOT'? 'NULL')? | 'TABLE' '(' (ID DATA_TYPE (',' ID DATA_TYPE)*)? ')') 'LANGUAGE' 'JAVA' 'RUNTIME_VERSION' '=' STRING 'PACKAGES' '=' '(' STRING (',' STRING)* ')' ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? 'HANDLER' '=' STRING ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' QNAME (',' STRING '=' QNAME)* ')')? (('CALLED' 'ON' 'NULL' 'INPUT' | ('RETURNS' 'NULL' 'ON' 'NULL' 'INPUT' | 'STRICT')))? ('VOLATILE' | 'IMMUTABLE')? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')? ('EXECUTE' 'AS' ('OWNER' | 'CALLER' | 'RESTRICTED' 'CALLER'))?
    | 'CREATE' ('OR' 'REPLACE')? (('TEMP' | 'TEMPORARY'))? 'SECURE'? 'PROCEDURE' QNAME '(' (ID DATA_TYPE ('DEFAULT' EXPR)? (',' ID DATA_TYPE ('DEFAULT' EXPR)?)*)? ')' ('COPY' 'GRANTS')? 'RETURNS' DATA_TYPE ('NOT' 'NULL')? 'LANGUAGE' 'JAVASCRIPT' (('CALLED' 'ON' 'NULL' 'INPUT' | ('RETURNS' 'NULL' 'ON' 'NULL' 'INPUT' | 'STRICT')))? ('VOLATILE' | 'IMMUTABLE')? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')? ('EXECUTE' 'AS' ('OWNER' | 'CALLER' | 'RESTRICTED' 'CALLER'))? 'AS' STRING
    | 'CREATE' ('OR' 'REPLACE')? 'SECURE'? 'PROCEDURE' QNAME '(' (ID DATA_TYPE ('DEFAULT' EXPR)? (',' ID DATA_TYPE ('DEFAULT' EXPR)?)*)? ')' ('COPY' 'GRANTS')? 'RETURNS' (DATA_TYPE ('NOT'? 'NULL')? | 'TABLE' '(' (ID DATA_TYPE (',' ID DATA_TYPE)*)? ')') 'LANGUAGE' 'SCALA' 'RUNTIME_VERSION' '=' STRING 'PACKAGES' '=' '(' STRING (',' STRING)* ')' ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? 'HANDLER' '=' STRING ('TARGET_PATH' '=' STRING)? (('CALLED' 'ON' 'NULL' 'INPUT' | ('RETURNS' 'NULL' 'ON' 'NULL' 'INPUT' | 'STRICT')))? (('VOLATILE' | 'IMMUTABLE'))? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')? ('EXECUTE' 'AS' ('OWNER' | 'CALLER' | 'RESTRICTED' 'CALLER'))? 'AS' STRING
    | 'CREATE' ('OR' 'REPLACE')? 'SECURE'? 'PROCEDURE' QNAME '(' (ID DATA_TYPE ('DEFAULT' EXPR)? (',' ID DATA_TYPE ('DEFAULT' EXPR)?)*)? ')' ('COPY' 'GRANTS')? 'RETURNS' (DATA_TYPE ('NOT'? 'NULL')? | 'TABLE' '(' (ID DATA_TYPE (',' ID DATA_TYPE)*)? ')') 'LANGUAGE' 'SCALA' 'RUNTIME_VERSION' '=' STRING 'PACKAGES' '=' '(' STRING (',' STRING)* ')' ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? 'HANDLER' '=' STRING (('CALLED' 'ON' 'NULL' 'INPUT' | ('RETURNS' 'NULL' 'ON' 'NULL' 'INPUT' | 'STRICT')))? ('VOLATILE' | 'IMMUTABLE')? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')? ('EXECUTE' 'AS' ('OWNER' | 'CALLER' | 'RESTRICTED' 'CALLER'))?
    | 'CREATE' ('OR' 'REPLACE')? 'PROCEDURE' QNAME '(' (ID (('IN' | 'INPUT' | 'OUT' | 'OUTPUT'))? DATA_TYPE ('DEFAULT' EXPR)? (',' ID (('IN' | 'INPUT' | 'OUT' | 'OUTPUT'))? DATA_TYPE ('DEFAULT' EXPR)?)*)? ')' ('COPY' 'GRANTS')? 'RETURNS' (DATA_TYPE | 'TABLE' '(' (ID DATA_TYPE (',' ID DATA_TYPE)*)? ')') ('NOT' 'NULL')? 'LANGUAGE' 'SQL' (('CALLED' 'ON' 'NULL' 'INPUT' | ('RETURNS' 'NULL' 'ON' 'NULL' 'INPUT' | 'STRICT')))? (('VOLATILE' | 'IMMUTABLE'))? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')? ('EXECUTE' 'AS' ('OWNER' | 'CALLER' | 'RESTRICTED' 'CALLER'))? 'AS' STRING
    | 'CREATE' ('OR' 'REPLACE')? 'PROCEDURE' 'SCOPED' ('TEMP' | 'TEMPORARY') 'TABLE' QNAME '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-projection-policy
createProjectionPolicy
    : 'CREATE' ('OR' 'REPLACE')? 'PROJECTION' 'POLICY' QNAME 'AS' '(' ')' 'RETURNS' 'PROJECTION_CONSTRAINT' '->' EXPR ('COMMENT' '=' STRING)?
    | 'CREATE' 'PROJECTION' 'POLICY' ('IF' 'NOT' 'EXISTS')? QNAME 'AS' '(' ')' 'RETURNS' 'PROJECTION_CONSTRAINT' '->' EXPR ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-provisioned-throughput
createProvisionedThroughput
    : 'CREATE' ('OR' 'REPLACE')? 'PROVISIONED' 'THROUGHPUT' QNAME 'CLOUD_PROVIDER' '=' STRING 'MODEL' '=' STRING 'PTUS' '=' NUMBER 'TERM_START' '=' STRING 'TERM_END' '=' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-replication-group
// fix: OBJECT_TYPES takes the replication object types the parameters list, not any object type
// fix: ALLOWED_INTEGRATION_TYPES takes the integration types the parameters list
createReplicationGroup
    : 'CREATE' 'REPLICATION' 'GROUP' ('IF' 'NOT' 'EXISTS')? ID 'OBJECT_TYPES' '=' ('ACCOUNT' 'PARAMETERS' | 'DATABASES' | 'EXTERNAL' 'VOLUMES' | 'INTEGRATIONS' | 'NETWORK' 'POLICIES' | 'RESOURCE' 'MONITORS' | 'ROLES' | 'SHARES' | 'USERS' | 'WAREHOUSES') (',' ('ACCOUNT' 'PARAMETERS' | 'DATABASES' | 'EXTERNAL' 'VOLUMES' | 'INTEGRATIONS' | 'NETWORK' 'POLICIES' | 'RESOURCE' 'MONITORS' | 'ROLES' | 'SHARES' | 'USERS' | 'WAREHOUSES'))* ('ALLOWED_DATABASES' '=' ID (',' ID)*)? ('ALLOWED_EXTERNAL_VOLUMES' '=' ID (',' ID)*)? ('ALLOWED_SHARES' '=' ID (',' ID)*)? ('ALLOWED_INTEGRATION_TYPES' '=' ('SECURITY' 'INTEGRATIONS' | 'API' 'INTEGRATIONS' | 'STORAGE' 'INTEGRATIONS' | 'EXTERNAL' 'ACCESS' 'INTEGRATIONS' | 'NOTIFICATION' 'INTEGRATIONS') (',' ('SECURITY' 'INTEGRATIONS' | 'API' 'INTEGRATIONS' | 'STORAGE' 'INTEGRATIONS' | 'EXTERNAL' 'ACCESS' 'INTEGRATIONS' | 'NOTIFICATION' 'INTEGRATIONS'))*)? 'ALLOWED_ACCOUNTS' '=' ID '.' ID (',' ID '.' ID)* ('IGNORE' 'EDITION' 'CHECK')? ('REPLICATION_SCHEDULE' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('ERROR_INTEGRATION' '=' ID)?
    | 'CREATE' 'REPLICATION' 'GROUP' ('IF' 'NOT' 'EXISTS')? ID 'AS' 'REPLICA' 'OF' ID '.' ID '.' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-resource-monitor
createResourceMonitor
    : 'CREATE' ('OR' 'REPLACE')? 'RESOURCE' 'MONITOR' ID 'WITH' ('CREDIT_QUOTA' '=' NUMBER)? ('FREQUENCY' '=' ('MONTHLY' | 'DAILY' | 'WEEKLY' | 'YEARLY' | 'NEVER'))? ('START_TIMESTAMP' '=' (STRING | 'IMMEDIATELY'))? ('END_TIMESTAMP' '=' STRING)? ('NOTIFY_USERS' '=' '(' ID (',' ID)* ')')? ('TRIGGERS' createResourceMonitor_triggerDefinition createResourceMonitor_triggerDefinition*)?
    | 'CREATE' 'RESOURCE' 'MONITOR' ('IF' 'NOT' 'EXISTS')? ID 'WITH' ('CREDIT_QUOTA' '=' NUMBER)? ('FREQUENCY' '=' ('MONTHLY' | 'DAILY' | 'WEEKLY' | 'YEARLY' | 'NEVER'))? ('START_TIMESTAMP' '=' (STRING | 'IMMEDIATELY'))? ('END_TIMESTAMP' '=' STRING)? ('NOTIFY_USERS' '=' '(' ID (',' ID)* ')')? ('TRIGGERS' createResourceMonitor_triggerDefinition createResourceMonitor_triggerDefinition*)?
    ;
createResourceMonitor_triggerDefinition
    : 'ON' NUMBER 'PERCENT' 'DO' ('SUSPEND' | 'SUSPEND_IMMEDIATE' | 'NOTIFY')
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-restricted-session-scope
createRestrictedSessionScope
    : 'CREATE' ('RESTRICTED' 'SESSION' 'SCOPE' | 'RSS') ('IF' 'NOT' 'EXISTS')? QNAME ('COMMENT' '=' STRING)? 'AS' STRING
    | 'CREATE' 'OR' 'ALTER' ('RESTRICTED' 'SESSION' 'SCOPE' | 'RSS') QNAME ('COMMENT' '=' STRING)? 'AS' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-role
createRole
    : 'CREATE' ('OR' 'REPLACE')? 'ROLE' ID ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' 'ROLE' ('IF' 'NOT' 'EXISTS')? ID ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' 'OR' 'ALTER' 'ROLE' ID ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-row-access-policy
createRowAccessPolicy
    : 'CREATE' ('OR' 'REPLACE')? 'ROW' 'ACCESS' 'POLICY' QNAME 'AS' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')' 'RETURNS' 'BOOLEAN' '->' EXPR ('COMMENT' '=' STRING)?
    | 'CREATE' 'ROW' 'ACCESS' 'POLICY' ('IF' 'NOT' 'EXISTS')? QNAME 'AS' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')' 'RETURNS' 'BOOLEAN' '->' EXPR ('COMMENT' '=' STRING)?
    | 'CREATE' 'OR' 'ALTER' 'ROW' 'ACCESS' 'POLICY' QNAME 'AS' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')' 'RETURNS' 'BOOLEAN' '->' EXPR ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-schema
// fix: read: Time Travel's OFFSET is a negative number of seconds, -N, and its STATEMENT a quoted query ID, as the Time Travel parameters say
createSchema
    : 'CREATE' ('OR' 'REPLACE')? 'TRANSIENT'? 'SCHEMA' QNAME ('CLONE' QNAME (('AT' | 'BEFORE') '(' ('TIMESTAMP' '=>' STRING | 'OFFSET' '=>' '-' NUMBER | 'STATEMENT' '=>' STRING) ')')? ('IGNORE' 'TABLES' 'WITH' 'INSUFFICIENT' 'DATA' 'RETENTION')? ('IGNORE' 'HYBRID' 'TABLES')?)? ('WITH' 'MANAGED' 'ACCESS')? ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('STAGE_UNDROP_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('EXTERNAL_VOLUME' '=' ID)? ('CATALOG' '=' ID)? ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_VERSION_DEFAULT' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CLASSIFICATION_PROFILE' '=' STRING)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    | 'CREATE' 'TRANSIENT'? 'SCHEMA' ('IF' 'NOT' 'EXISTS')? QNAME ('CLONE' QNAME (('AT' | 'BEFORE') '(' ('TIMESTAMP' '=>' STRING | 'OFFSET' '=>' '-' NUMBER | 'STATEMENT' '=>' STRING) ')')? ('IGNORE' 'TABLES' 'WITH' 'INSUFFICIENT' 'DATA' 'RETENTION')? ('IGNORE' 'HYBRID' 'TABLES')?)? ('WITH' 'MANAGED' 'ACCESS')? ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('STAGE_UNDROP_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('EXTERNAL_VOLUME' '=' ID)? ('CATALOG' '=' ID)? ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_VERSION_DEFAULT' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('CLASSIFICATION_PROFILE' '=' STRING)? ('COMMENT' '=' STRING)? ('CATALOG_SYNC' '=' STRING)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')? ('ENABLE_DATA_COMPACTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_AUTHORIZATION_SERVER' '=' ID)? ('OAUTH_SCOPES_SUPPORTED' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')?
    | 'CREATE' 'SCHEMA' QNAME 'FROM' 'BACKUP' 'SET' ID 'IDENTIFIER' STRING
    | 'CREATE' 'OR' 'ALTER' 'TRANSIENT'? 'SCHEMA' QNAME ('WITH' 'MANAGED' 'ACCESS')? ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('EXTERNAL_VOLUME' '=' ID)? ('CATALOG' '=' ID)? ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_VERSION_DEFAULT' '=' NUMBER)? ('ICEBERG_MERGE_ON_READ_BEHAVIOR' '=' ('\'AUTO\'' | '\'ENABLED\'' | '\'DISABLED\''))? ('ENABLE_ICEBERG_MERGE_ON_READ' '=' ('TRUE' | 'FALSE'))? ('REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('LOG_LEVEL' '=' STRING)? ('TRACE_LEVEL' '=' STRING)? ('STORAGE_SERIALIZATION_POLICY' '=' ('COMPATIBLE' | 'OPTIMIZED'))? ('COMMENT' '=' STRING)? ('OBJECT_VISIBILITY' '=' 'PRIVILEGED')?
    | 'CREATE' ('OR' 'REPLACE')? 'SCHEMA' QNAME 'CLONE' QNAME
    | 'CREATE' 'SCHEMA' ('IF' 'NOT' 'EXISTS')? QNAME 'CLONE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-secret
// fix: leaving OAUTH_SCOPES out inherits the integration's scopes, as the parameters say
// not converted (form): DEFINE runs only in DCM project files, not as a statement
//   DEFINE SECRET MY_DB.MY_SCHEMA.THIRD_PARTY_API_KEY
//     TYPE = GENERIC_STRING
//     SECRET_STRING = {{ _snow.env_secret("API_KEY") }}
//     COMMENT = 'API key for the third-party billing service';
createSecret
    : 'CREATE' ('OR' 'REPLACE')? 'SECRET' QNAME 'TYPE' '=' 'OAUTH2' 'API_AUTHENTICATION' '=' ID ('OAUTH_SCOPES' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? ('WITH' 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' 'SECRET' ('IF' 'NOT' 'EXISTS')? QNAME 'TYPE' '=' 'OAUTH2' 'API_AUTHENTICATION' '=' ID ('OAUTH_SCOPES' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? ('WITH' 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'SECRET' QNAME 'TYPE' '=' 'OAUTH2' 'OAUTH_REFRESH_TOKEN' '=' STRING 'OAUTH_REFRESH_TOKEN_EXPIRY_TIME' '=' STRING 'API_AUTHENTICATION' '=' ID ('COMMENT' '=' STRING)? ('WITH' 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' 'SECRET' ('IF' 'NOT' 'EXISTS')? QNAME 'TYPE' '=' 'OAUTH2' 'OAUTH_REFRESH_TOKEN' '=' STRING 'OAUTH_REFRESH_TOKEN_EXPIRY_TIME' '=' STRING 'API_AUTHENTICATION' '=' ID ('COMMENT' '=' STRING)? ('WITH' 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'SECRET' QNAME 'TYPE' '=' 'CLOUD_PROVIDER_TOKEN' 'API_AUTHENTICATION' '=' STRING 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)? ('WITH' 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' 'SECRET' ('IF' 'NOT' 'EXISTS')? QNAME 'TYPE' '=' 'CLOUD_PROVIDER_TOKEN' 'API_AUTHENTICATION' '=' STRING 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)? ('WITH' 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'SECRET' QNAME 'TYPE' '=' 'PASSWORD' 'USERNAME' '=' STRING 'PASSWORD' '=' STRING ('COMMENT' '=' STRING)? ('WITH' 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' 'SECRET' ('IF' 'NOT' 'EXISTS')? QNAME 'TYPE' '=' 'PASSWORD' 'USERNAME' '=' STRING 'PASSWORD' '=' STRING ('COMMENT' '=' STRING)? ('WITH' 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'SECRET' QNAME 'TYPE' '=' 'GENERIC_STRING' 'SECRET_STRING' '=' STRING ('COMMENT' '=' STRING)? ('WITH' 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' 'SECRET' ('IF' 'NOT' 'EXISTS')? QNAME 'TYPE' '=' 'GENERIC_STRING' 'SECRET_STRING' '=' STRING ('COMMENT' '=' STRING)? ('WITH' 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'SECRET' QNAME 'TYPE' '=' 'SYMMETRIC_KEY' 'ALGORITHM' '=' 'GENERIC' ('WITH' 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' 'SECRET' ('IF' 'NOT' 'EXISTS')? QNAME 'TYPE' '=' 'SYMMETRIC_KEY' 'ALGORITHM' '=' 'GENERIC' ('WITH' 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'SECRET' QNAME 'TYPE' '=' 'WORKLOAD_IDENTITY_FEDERATION' ('COMMENT' '=' STRING)? ('WITH' 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' 'SECRET' ('IF' 'NOT' 'EXISTS')? QNAME 'TYPE' '=' 'WORKLOAD_IDENTITY_FEDERATION' ('COMMENT' '=' STRING)? ('WITH' 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-security-integration
// not converted (form): ... stands for text the docs leave out
//   CREATE [ OR REPLACE ] SECURITY INTEGRATION [ IF NOT EXISTS ]
//     <name>
//     TYPE = { API_AUTHENTICATION | EXTERNAL_OAUTH | OAUTH | OIDC | SAML2 | SCIM }
//     ...

// https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-api-auth
createSecurityIntegrationApiAuth
    : 'CREATE' 'SECURITY' 'INTEGRATION' ID 'TYPE' '=' 'API_AUTHENTICATION' 'AUTH_TYPE' '=' 'OAUTH2' 'ENABLED' '=' ('TRUE' | 'FALSE') ('OAUTH_TOKEN_ENDPOINT' '=' STRING)? ('OAUTH_CLIENT_AUTH_METHOD' '=' ('CLIENT_SECRET_BASIC' | 'CLIENT_SECRET_POST'))? ('OAUTH_CLIENT_ID' '=' STRING)? ('OAUTH_CLIENT_SECRET' '=' STRING)? ('OAUTH_GRANT' '=' '\'CLIENT_CREDENTIALS\'')? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_ALLOWED_SCOPES' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' 'SECURITY' 'INTEGRATION' ID 'TYPE' '=' 'API_AUTHENTICATION' 'AUTH_TYPE' '=' 'OAUTH2' 'ENABLED' '=' ('TRUE' | 'FALSE') ('OAUTH_AUTHORIZATION_ENDPOINT' '=' STRING)? ('OAUTH_TOKEN_ENDPOINT' '=' STRING)? ('OAUTH_CLIENT_AUTH_METHOD' '=' ('CLIENT_SECRET_BASIC' | 'CLIENT_SECRET_POST'))? ('OAUTH_CLIENT_ID' '=' STRING)? ('OAUTH_CLIENT_SECRET' '=' STRING)? ('OAUTH_GRANT' '=' '\'AUTHORIZATION_CODE\'')? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'SECURITY' 'INTEGRATION' ID 'TYPE' '=' 'API_AUTHENTICATION' 'AUTH_TYPE' '=' 'OAUTH2' 'ENABLED' '=' ('TRUE' | 'FALSE') ('OAUTH_AUTHORIZATION_ENDPOINT' '=' STRING)? ('OAUTH_TOKEN_ENDPOINT' '=' STRING)? ('OAUTH_CLIENT_AUTH_METHOD' '=' ('CLIENT_SECRET_BASIC' | 'CLIENT_SECRET_POST'))? ('OAUTH_CLIENT_ID' '=' STRING)? ('OAUTH_CLIENT_SECRET' '=' STRING)? 'OAUTH_ASSERTION_ISSUER' '=' STRING ('OAUTH_GRANT' '=' '\'JWT_BEARER\'')? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-api-auth-external-secret-provider
createSecurityIntegrationApiAuthExternalSecretProvider
    : 'CREATE' ('OR' 'REPLACE')? 'SECURITY' 'INTEGRATION' ID 'TYPE' '=' 'API_AUTHENTICATION' 'AUTH_TYPE' '=' 'WORKLOAD_IDENTITY_FEDERATION' 'API_PROVIDER' '=' ('AWS_SECRETS_MANAGER' | 'AZURE_KEY_VAULT' | 'GCP_SECRET_MANAGER') 'ENABLED' '=' ('TRUE' | 'FALSE') ('AWS_ROLE_ARN' '=' STRING)? ('AWS_REGION' '=' STRING)? ('AZURE_TENANT_ID' '=' STRING)? ('AZURE_AD_APPLICATION_ID' '=' STRING)? ('AZURE_KEY_VAULT_URI' '=' STRING)? ('GCP_SECRET_MANAGER_PROJECT' '=' STRING)? ('GCP_WIF_PROVIDER' '=' STRING)? ('GCP_SERVICE_ACCOUNT_EMAIL' '=' STRING)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'SECURITY' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'TYPE' '=' 'API_AUTHENTICATION' 'AUTH_TYPE' '=' 'WORKLOAD_IDENTITY_FEDERATION' 'API_PROVIDER' '=' ('AWS_SECRETS_MANAGER' | 'AZURE_KEY_VAULT' | 'GCP_SECRET_MANAGER') 'ENABLED' '=' ('TRUE' | 'FALSE') ('AWS_ROLE_ARN' '=' STRING)? ('AWS_REGION' '=' STRING)? ('AZURE_TENANT_ID' '=' STRING)? ('AZURE_AD_APPLICATION_ID' '=' STRING)? ('AZURE_KEY_VAULT_URI' '=' STRING)? ('GCP_SECRET_MANAGER_PROJECT' '=' STRING)? ('GCP_WIF_PROVIDER' '=' STRING)? ('GCP_SERVICE_ACCOUNT_EMAIL' '=' STRING)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-aws-iam
createSecurityIntegrationAwsIam
    : 'CREATE' 'SECURITY' 'INTEGRATION' ID 'TYPE' '=' 'API_AUTHENTICATION' 'AUTH_TYPE' '=' 'AWS_IAM' 'AWS_ROLE_ARN' '=' STRING 'ENABLED' '=' ('TRUE' | 'FALSE') ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-oauth-external
createSecurityIntegrationOauthExternal
    : 'CREATE' ('OR' 'REPLACE')? 'SECURITY' 'INTEGRATION' ID 'TYPE' '=' 'EXTERNAL_OAUTH' 'ENABLED' '=' ('TRUE' | 'FALSE') 'EXTERNAL_OAUTH_TYPE' '=' ('OKTA' | 'AZURE' | 'PING_FEDERATE' | 'CUSTOM') 'EXTERNAL_OAUTH_ISSUER' '=' STRING 'EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM' '=' (STRING | '(' STRING (',' STRING)* ')') 'EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE' '=' ('\'LOGIN_NAME\'' | '\'EMAIL_ADDRESS\'') ('EXTERNAL_OAUTH_JWS_KEYS_URL' '=' (STRING | '(' STRING (',' STRING)* ')'))? ('EXTERNAL_OAUTH_BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_RSA_PUBLIC_KEY' '=' STRING)? ('EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2' '=' STRING)? ('EXTERNAL_OAUTH_AUDIENCE_LIST' '=' (STRING | '(' STRING (',' STRING)* ')'))? ('EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' ('DISABLE' | 'ENABLE' | 'ENABLE_FOR_PRIVILEGE'))? ('EXTERNAL_OAUTH_SCOPE_DELIMITER' '=' STRING)? ('EXTERNAL_OAUTH_SCOPE_MAPPING_ATTRIBUTE' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('NETWORK_POLICY' '=' STRING)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'SECURITY' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'TYPE' '=' 'EXTERNAL_OAUTH' 'ENABLED' '=' ('TRUE' | 'FALSE') 'EXTERNAL_OAUTH_TYPE' '=' ('OKTA' | 'AZURE' | 'PING_FEDERATE' | 'CUSTOM') 'EXTERNAL_OAUTH_ISSUER' '=' STRING 'EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM' '=' (STRING | '(' STRING (',' STRING)* ')') 'EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE' '=' ('\'LOGIN_NAME\'' | '\'EMAIL_ADDRESS\'') ('EXTERNAL_OAUTH_JWS_KEYS_URL' '=' (STRING | '(' STRING (',' STRING)* ')'))? ('EXTERNAL_OAUTH_BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_OAUTH_RSA_PUBLIC_KEY' '=' STRING)? ('EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2' '=' STRING)? ('EXTERNAL_OAUTH_AUDIENCE_LIST' '=' (STRING | '(' STRING (',' STRING)* ')'))? ('EXTERNAL_OAUTH_ANY_ROLE_MODE' '=' ('DISABLE' | 'ENABLE' | 'ENABLE_FOR_PRIVILEGE'))? ('EXTERNAL_OAUTH_SCOPE_DELIMITER' '=' STRING)? ('EXTERNAL_OAUTH_SCOPE_MAPPING_ATTRIBUTE' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('NETWORK_POLICY' '=' STRING)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-oauth-snowflake
// fix: the partner applications are TABLEAU_DESKTOP, TABLEAU_SERVER and LOOKER, and only LOOKER needs OAUTH_REDIRECT_URI, as the parameters say
createSecurityIntegrationOauthSnowflake
    : 'CREATE' ('OR' 'REPLACE')? 'SECURITY' 'INTEGRATION' ID 'TYPE' '=' 'OAUTH' ('OAUTH_CLIENT' '=' ('TABLEAU_DESKTOP' | 'TABLEAU_SERVER') | 'OAUTH_CLIENT' '=' 'LOOKER' 'OAUTH_REDIRECT_URI' '=' STRING) ('ENABLED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ISSUE_REFRESH_TOKENS' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('NETWORK_POLICY' '=' STRING)? ('ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)?
    | 'CREATE' 'SECURITY' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'TYPE' '=' 'OAUTH' ('OAUTH_CLIENT' '=' ('TABLEAU_DESKTOP' | 'TABLEAU_SERVER') | 'OAUTH_CLIENT' '=' 'LOOKER' 'OAUTH_REDIRECT_URI' '=' STRING) ('ENABLED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ISSUE_REFRESH_TOKENS' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('NETWORK_POLICY' '=' STRING)? ('ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)?
    | 'CREATE' ('OR' 'REPLACE')? 'SECURITY' 'INTEGRATION' ID 'TYPE' '=' 'OAUTH' 'OAUTH_CLIENT' '=' 'CUSTOM' 'OAUTH_CLIENT_TYPE' '=' ('\'CONFIDENTIAL\'' | '\'PUBLIC\'') 'OAUTH_REDIRECT_URI' '=' STRING ('ENABLED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ALLOW_NON_TLS_REDIRECT_URI' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ALTERNATE_REDIRECT_URIS' '=' '(' STRING (',' STRING)* ')')? ('OAUTH_ENFORCE_PKCE' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ENABLE_ROLE_SELECTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('PRE_AUTHORIZED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('OAUTH_ANY_ROLE_MODE' '=' ('DISABLE' | 'ENABLE' | 'ENABLE_FOR_PRIVILEGE'))? ('OAUTH_ISSUE_REFRESH_TOKENS' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('NETWORK_POLICY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY_2' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)?
    | 'CREATE' 'SECURITY' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'TYPE' '=' 'OAUTH' 'OAUTH_CLIENT' '=' 'CUSTOM' 'OAUTH_CLIENT_TYPE' '=' ('\'CONFIDENTIAL\'' | '\'PUBLIC\'') 'OAUTH_REDIRECT_URI' '=' STRING ('ENABLED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ALLOW_NON_TLS_REDIRECT_URI' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ALTERNATE_REDIRECT_URIS' '=' '(' STRING (',' STRING)* ')')? ('OAUTH_ENFORCE_PKCE' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ENABLE_ROLE_SELECTION' '=' ('TRUE' | 'FALSE'))? ('OAUTH_SINGLE_USE_REFRESH_TOKENS_REQUIRED' '=' ('TRUE' | 'FALSE'))? ('OAUTH_USE_SECONDARY_ROLES' '=' ('IMPLICIT' | 'NONE'))? ('PRE_AUTHORIZED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('BLOCKED_ROLES_LIST' '=' '(' STRING (',' STRING)* ')')? ('OAUTH_ANY_ROLE_MODE' '=' ('DISABLE' | 'ENABLE' | 'ENABLE_FOR_PRIVILEGE'))? ('OAUTH_ISSUE_REFRESH_TOKENS' '=' ('TRUE' | 'FALSE'))? ('OAUTH_ACCESS_TOKEN_VALIDITY' '=' NUMBER)? ('OAUTH_REFRESH_TOKEN_VALIDITY' '=' NUMBER)? ('NETWORK_POLICY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY' '=' STRING)? ('OAUTH_CLIENT_RSA_PUBLIC_KEY_2' '=' STRING)? ('IS_AGENTIC' '=' ('TRUE' | 'FALSE'))? ('USE_PRIVATELINK_FOR_AUTHORIZATION_ENDPOINT' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-oidc
createSecurityIntegrationOidc
    : 'CREATE' ('OR' 'REPLACE')? 'SECURITY' 'INTEGRATION' ID 'TYPE' '=' 'OIDC' ('ENABLED' '=' ('TRUE' | 'FALSE'))? 'OIDC_PROVIDER' '=' '\'CUSTOM\'' 'OIDC_ISSUER' '=' STRING 'OIDC_CLIENT_ID' '=' STRING 'OIDC_CLIENT_SECRET' '=' STRING 'OIDC_TOKEN_USER_MAPPING_CLAIM' '=' STRING 'OIDC_SNOWFLAKE_USER_MAPPING_ATTRIBUTE' '=' ('\'LOGIN_NAME\'' | '\'EMAIL_ADDRESS\'') ('OIDC_LOGIN_PAGE_LABEL' '=' STRING)? ('OIDC_ENABLE_SSO_LOGIN_PAGE' '=' ('TRUE' | 'FALSE'))? ('OIDC_SCOPES' '=' '(' STRING (',' STRING)* ')')? ('OIDC_AUTHORIZATION_ENDPOINT' '=' STRING)? ('OIDC_TOKEN_ENDPOINT' '=' STRING)? ('OIDC_JWKS_URI' '=' STRING)? ('OIDC_USERINFO_ENDPOINT' '=' STRING)? ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' 'SECURITY' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'TYPE' '=' 'OIDC' ('ENABLED' '=' ('TRUE' | 'FALSE'))? 'OIDC_PROVIDER' '=' '\'CUSTOM\'' 'OIDC_ISSUER' '=' STRING 'OIDC_CLIENT_ID' '=' STRING 'OIDC_CLIENT_SECRET' '=' STRING 'OIDC_TOKEN_USER_MAPPING_CLAIM' '=' STRING 'OIDC_SNOWFLAKE_USER_MAPPING_ATTRIBUTE' '=' ('\'LOGIN_NAME\'' | '\'EMAIL_ADDRESS\'') ('OIDC_LOGIN_PAGE_LABEL' '=' STRING)? ('OIDC_ENABLE_SSO_LOGIN_PAGE' '=' ('TRUE' | 'FALSE'))? ('OIDC_SCOPES' '=' '(' STRING (',' STRING)* ')')? ('OIDC_AUTHORIZATION_ENDPOINT' '=' STRING)? ('OIDC_TOKEN_ENDPOINT' '=' STRING)? ('OIDC_JWKS_URI' '=' STRING)? ('OIDC_USERINFO_ENDPOINT' '=' STRING)? ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' ('OR' 'REPLACE')? 'SECURITY' 'INTEGRATION' ID 'TYPE' '=' 'OIDC' ('ENABLED' '=' ('TRUE' | 'FALSE'))? 'OIDC_PROVIDER' '=' '\'GOOGLE\'' 'OIDC_ALLOWED_GOOGLE_HOSTED_DOMAINS' '=' '(' STRING (',' STRING)* ')' ('OIDC_ENABLE_SSO_LOGIN_PAGE' '=' ('TRUE' | 'FALSE'))? ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' 'SECURITY' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'TYPE' '=' 'OIDC' ('ENABLED' '=' ('TRUE' | 'FALSE'))? 'OIDC_PROVIDER' '=' '\'GOOGLE\'' 'OIDC_ALLOWED_GOOGLE_HOSTED_DOMAINS' '=' '(' STRING (',' STRING)* ')' ('OIDC_ENABLE_SSO_LOGIN_PAGE' '=' ('TRUE' | 'FALSE'))? ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' ('OR' 'REPLACE')? 'SECURITY' 'INTEGRATION' ID 'TYPE' '=' 'OIDC' ('ENABLED' '=' ('TRUE' | 'FALSE'))? 'OIDC_PROVIDER' '=' '\'MICROSOFT\'' 'OIDC_ALLOWED_MSFT_TENANTS' '=' '(' STRING (',' STRING)* ')' ('OIDC_ENABLE_SSO_LOGIN_PAGE' '=' ('TRUE' | 'FALSE'))? ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' 'SECURITY' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'TYPE' '=' 'OIDC' ('ENABLED' '=' ('TRUE' | 'FALSE'))? 'OIDC_PROVIDER' '=' '\'MICROSOFT\'' 'OIDC_ALLOWED_MSFT_TENANTS' '=' '(' STRING (',' STRING)* ')' ('OIDC_ENABLE_SSO_LOGIN_PAGE' '=' ('TRUE' | 'FALSE'))? ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-saml2
// fix: <idp_parameters> is SAML2_ISSUER, SAML2_SSO_URL, SAML2_PROVIDER and SAML2_X509_CERT, as the required parameters say
createSecurityIntegrationSaml2
    : 'CREATE' ('OR' 'REPLACE')? 'SECURITY' 'INTEGRATION' QNAME 'TYPE' '=' 'SAML2' 'ENABLED' '=' ('TRUE' | 'FALSE') ('METADATA_URL' '=' STRING | 'SAML2_ISSUER' '=' STRING 'SAML2_SSO_URL' '=' STRING 'SAML2_PROVIDER' '=' STRING 'SAML2_X509_CERT' '=' STRING) ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('SAML2_SP_INITIATED_LOGIN_PAGE_LABEL' '=' STRING)? ('SAML2_ENABLE_SP_INITIATED' '=' 'TRUE' | 'SAML2_ENABLE_SP_INITIATED' '=' 'FALSE')? ('SAML2_SNOWFLAKE_X509_CERT' '=' STRING)? ('SAML2_SIGN_REQUEST' '=' 'TRUE' | 'SAML2_SIGN_REQUEST' '=' 'FALSE')? ('SAML2_REQUESTED_NAMEID_FORMAT' '=' STRING)? ('SAML2_POST_LOGOUT_REDIRECT_URL' '=' STRING)? ('SAML2_FORCE_AUTHN' '=' 'TRUE' | 'SAML2_FORCE_AUTHN' '=' 'FALSE')? ('SAML2_SNOWFLAKE_ISSUER_URL' '=' STRING)? ('SAML2_SNOWFLAKE_ACS_URL' '=' STRING)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'SECURITY' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? QNAME 'TYPE' '=' 'SAML2' 'ENABLED' '=' ('TRUE' | 'FALSE') ('METADATA_URL' '=' STRING | 'SAML2_ISSUER' '=' STRING 'SAML2_SSO_URL' '=' STRING 'SAML2_PROVIDER' '=' STRING 'SAML2_X509_CERT' '=' STRING) ('ALLOWED_USER_DOMAINS' '=' '(' STRING (',' STRING)* ')')? ('ALLOWED_EMAIL_PATTERNS' '=' '(' STRING (',' STRING)* ')')? ('SAML2_SP_INITIATED_LOGIN_PAGE_LABEL' '=' STRING)? ('SAML2_ENABLE_SP_INITIATED' '=' 'TRUE' | 'SAML2_ENABLE_SP_INITIATED' '=' 'FALSE')? ('SAML2_SNOWFLAKE_X509_CERT' '=' STRING)? ('SAML2_SIGN_REQUEST' '=' 'TRUE' | 'SAML2_SIGN_REQUEST' '=' 'FALSE')? ('SAML2_REQUESTED_NAMEID_FORMAT' '=' STRING)? ('SAML2_POST_LOGOUT_REDIRECT_URL' '=' STRING)? ('SAML2_FORCE_AUTHN' '=' 'TRUE' | 'SAML2_FORCE_AUTHN' '=' 'FALSE')? ('SAML2_SNOWFLAKE_ISSUER_URL' '=' STRING)? ('SAML2_SNOWFLAKE_ACS_URL' '=' STRING)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-scim
createSecurityIntegrationScim
    : 'CREATE' ('OR' 'REPLACE')? 'SECURITY' 'INTEGRATION' ID 'TYPE' '=' 'SCIM' 'ENABLED' '=' ('TRUE' | 'FALSE') 'SCIM_CLIENT' '=' ('\'OKTA\'' | '\'AZURE\'' | '\'GENERIC\'') 'RUN_AS_ROLE' '=' ('\'OKTA_PROVISIONER\'' | '\'AAD_PROVISIONER\'' | '\'GENERIC_SCIM_PROVISIONER\'' | STRING) ('NETWORK_POLICY' '=' STRING)? ('SYNC_PASSWORD' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)?
    | 'CREATE' 'SECURITY' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'TYPE' '=' 'SCIM' 'ENABLED' '=' ('TRUE' | 'FALSE') 'SCIM_CLIENT' '=' ('\'OKTA\'' | '\'AZURE\'' | '\'GENERIC\'') 'RUN_AS_ROLE' '=' ('\'OKTA_PROVISIONER\'' | '\'AAD_PROVISIONER\'' | '\'GENERIC_SCIM_PROVISIONER\'' | STRING) ('NETWORK_POLICY' '=' STRING)? ('SYNC_PASSWORD' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-semantic-view
// fix: the [ ... ] repeats the whole UNIQUE ( ... ) clause
// fix: the frame clause is window-function syntax, which the command pages do not define
// fix: VERIFIED_AT is seconds since the Unix epoch, as the parameters say
// fix: the flag is a boolean, as the parameters say
createSemanticView
    : 'CREATE' ('OR' 'REPLACE')? 'SEMANTIC' 'VIEW' QNAME 'TABLES' '(' createSemanticView_logicalTable (',' createSemanticView_logicalTable)* ')' ('RELATIONSHIPS' '(' createSemanticView_relationshipDef (',' createSemanticView_relationshipDef)* ')')? ('FACTS' '(' createSemanticView_factExpression (',' createSemanticView_factExpression)* ')')? ('DIMENSIONS' '(' createSemanticView_dimensionExpression (',' createSemanticView_dimensionExpression)* ')')? ('METRICS' '(' (createSemanticView_metricExpression | createSemanticView_windowFunctionMetricExpression) (',' (createSemanticView_metricExpression | createSemanticView_windowFunctionMetricExpression))* ')')? ('COMMENT' '=' STRING)? ('MAX_STALENESS' '=' NUMBER)? ('AI_SQL_GENERATION' STRING)? ('AI_QUESTION_CATEGORIZATION' STRING)? ('AI_VERIFIED_QUERIES' '(' createSemanticView_verifiedQuery (',' createSemanticView_verifiedQuery)* ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COPY' 'GRANTS')?
    | 'CREATE' 'SEMANTIC' 'VIEW' ('IF' 'NOT' 'EXISTS')? QNAME 'TABLES' '(' createSemanticView_logicalTable (',' createSemanticView_logicalTable)* ')' ('RELATIONSHIPS' '(' createSemanticView_relationshipDef (',' createSemanticView_relationshipDef)* ')')? ('FACTS' '(' createSemanticView_factExpression (',' createSemanticView_factExpression)* ')')? ('DIMENSIONS' '(' createSemanticView_dimensionExpression (',' createSemanticView_dimensionExpression)* ')')? ('METRICS' '(' (createSemanticView_metricExpression | createSemanticView_windowFunctionMetricExpression) (',' (createSemanticView_metricExpression | createSemanticView_windowFunctionMetricExpression))* ')')? ('COMMENT' '=' STRING)? ('MAX_STALENESS' '=' NUMBER)? ('AI_SQL_GENERATION' STRING)? ('AI_QUESTION_CATEGORIZATION' STRING)? ('AI_VERIFIED_QUERIES' '(' createSemanticView_verifiedQuery (',' createSemanticView_verifiedQuery)* ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COPY' 'GRANTS')?
    | 'CREATE' 'OR' 'ALTER' 'SEMANTIC' 'VIEW' QNAME 'TABLES' '(' createSemanticView_logicalTable (',' createSemanticView_logicalTable)* ')' ('RELATIONSHIPS' '(' createSemanticView_relationshipDef (',' createSemanticView_relationshipDef)* ')')? ('FACTS' '(' createSemanticView_factExpression (',' createSemanticView_factExpression)* ')')? ('DIMENSIONS' '(' createSemanticView_dimensionExpression (',' createSemanticView_dimensionExpression)* ')')? ('METRICS' '(' (createSemanticView_metricExpression | createSemanticView_windowFunctionMetricExpression) (',' (createSemanticView_metricExpression | createSemanticView_windowFunctionMetricExpression))* ')')? ('COMMENT' '=' STRING)? ('MAX_STALENESS' '=' NUMBER)? ('AI_SQL_GENERATION' STRING)? ('AI_QUESTION_CATEGORIZATION' STRING)? ('AI_VERIFIED_QUERIES' '(' createSemanticView_verifiedQuery (',' createSemanticView_verifiedQuery)* ')')?
    ;
createSemanticView_logicalTable
    : (ID 'AS')? QNAME ('PRIMARY' 'KEY' '(' ID (',' ID)* ')')? (('UNIQUE' '(' ID (',' ID)* ')')+)? ('CONSTRAINT' ID? 'DISTINCT' 'RANGE' 'BETWEEN' ID 'AND' ID 'EXCLUSIVE')? ('WITH' 'SYNONYMS' '='? '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    ;
createSemanticView_relationshipDef
    : (ID 'AS')? ID '(' ID (',' ID)* ')' 'REFERENCES' ID ('(' ('ASOF'? ID (',' ID)* | 'BETWEEN' ID 'AND' ID 'EXCLUSIVE') ')')?
    ;
createSemanticView_factExpression
    : (('PRIVATE' | 'PUBLIC'))? ID '.' ID ('LABELS' '=' '(' 'FILTER' ')')? 'AS' EXPR ('WITH' 'SYNONYMS' '='? '(' STRING (',' STRING)* ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' '=' STRING)?
    ;
createSemanticView_dimensionExpression
    : 'PUBLIC'? ID '.' ID ('LABELS' '=' '(' 'FILTER' ')')? 'AS' EXPR ('WITH' 'SYNONYMS' '='? '(' STRING (',' STRING)* ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' '=' STRING)? ('WITH' 'CORTEX' 'SEARCH' 'SERVICE' QNAME ('USING' ID)?)?
    ;
createSemanticView_metricExpression
    : (('PRIVATE' | 'PUBLIC'))? ID '.' ID ('USING' '(' ID (',' ID)* ')')? ('NON' 'ADDITIVE' 'BY' '(' ID (('ASC' | 'DESC'))? ('NULLS' ('FIRST' | 'LAST'))? (',' ID (('ASC' | 'DESC'))? ('NULLS' ('FIRST' | 'LAST'))?)* ')')? 'AS' EXPR ('WITH' 'SYNONYMS' '='? '(' STRING (',' STRING)* ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' '=' STRING)?
    ;
createSemanticView_windowFunctionMetricExpression
    : (('PRIVATE' | 'PUBLIC'))? ID '.' ID 'AS' ID '(' ID ')' 'OVER' '(' ('PARTITION' 'BY' (EXPR | 'EXCLUDING' ID))? ('ORDER' 'BY' EXPR ('ASC' | 'DESC')? ('NULLS' ('FIRST' | 'LAST'))? (',' EXPR ('ASC' | 'DESC')? ('NULLS' ('FIRST' | 'LAST'))?)*)? ')'
    ;
createSemanticView_verifiedQuery
    : ID 'AS' '(' 'QUESTION' STRING ('VERIFIED_AT' NUMBER)? ('ONBOARDING_QUESTION' ('TRUE' | 'FALSE'))? ('VERIFIED_BY' STRING)? 'SQL' STRING ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-sequence
// fix: the start may be negative, as the parameters say
// fix: the increment may be negative, as the parameters say
createSequence
    : 'CREATE' ('OR' 'REPLACE')? 'SEQUENCE' QNAME 'WITH'? ('START' 'WITH'? '='? '-'? NUMBER)? ('INCREMENT' 'BY'? '='? '-'? NUMBER)? (('ORDER' | 'NOORDER'))? ('COMMENT' '=' STRING)?
    | 'CREATE' 'SEQUENCE' ('IF' 'NOT' 'EXISTS')? QNAME 'WITH'? ('START' 'WITH'? '='? '-'? NUMBER)? ('INCREMENT' 'BY'? '='? '-'? NUMBER)? (('ORDER' | 'NOORDER'))? ('COMMENT' '=' STRING)?
    | 'CREATE' 'OR' 'ALTER' 'SEQUENCE' QNAME 'WITH'? ('START' 'WITH'? '='? '-'? NUMBER)? ('INCREMENT' 'BY'? '='? '-'? NUMBER)? (('ORDER' | 'NOORDER'))? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-service
// fix: the comment is text to fill in, written in braces where the other pages write <string_literal>
// fix: read: a template argument's <value> is a constant, a string or JSON value, as the USING parameters say
createService
    : 'CREATE' 'SERVICE' ('IF' 'NOT' 'EXISTS')? QNAME 'IN' 'COMPUTE' 'POOL' ID (createService_fromSpecification | createService_fromSpecificationTemplate) ('AUTO_SUSPEND_SECS' '=' NUMBER)? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('AUTO_RESUME' '=' ('TRUE' | 'FALSE'))? ('MIN_INSTANCES' '=' NUMBER)? ('MIN_READY_INSTANCES' '=' NUMBER)? ('MAX_INSTANCES' '=' NUMBER)? ('LOG_LEVEL' '=' STRING)? ('QUERY_WAREHOUSE' '=' ID)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' '=' STRING)?
    ;
createService_fromSpecification
    : ('FROM' 'SPECIFICATION_FILE' '=' STRING | 'FROM' STAGE 'SPECIFICATION_FILE' '=' STRING | 'FROM' 'SPECIFICATION' STRING)
    ;
createService_fromSpecificationTemplate
    : ('FROM' 'SPECIFICATION_TEMPLATE_FILE' '=' STRING | 'FROM' STAGE 'SPECIFICATION_TEMPLATE_FILE' '=' STRING | 'FROM' 'SPECIFICATION_TEMPLATE' STRING) 'USING' '(' ID '=>' EXPR (',' ID '=>' EXPR)* ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-session-policy
createSessionPolicy
    : 'CREATE' ('OR' 'REPLACE')? 'SESSION' 'POLICY' QNAME ('SESSION_IDLE_TIMEOUT_MINS' '=' NUMBER)? ('SESSION_UI_IDLE_TIMEOUT_MINS' '=' NUMBER)? ('SESSION_MAX_LIFESPAN_MINS' '=' NUMBER)? ('SESSION_UI_MAX_LIFESPAN_MINS' '=' NUMBER)? ('ALLOWED_SECONDARY_ROLES' '=' '(' (('\'ALL\'' | ID (',' ID)*))? ')')? ('BLOCKED_SECONDARY_ROLES' '=' '(' (('\'ALL\'' | ID (',' ID)*))? ')')? ('AGENT_RESTRICTED_SESSION_SCOPE' '=' STRING)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'SESSION' 'POLICY' ('IF' 'NOT' 'EXISTS')? QNAME ('SESSION_IDLE_TIMEOUT_MINS' '=' NUMBER)? ('SESSION_UI_IDLE_TIMEOUT_MINS' '=' NUMBER)? ('SESSION_MAX_LIFESPAN_MINS' '=' NUMBER)? ('SESSION_UI_MAX_LIFESPAN_MINS' '=' NUMBER)? ('ALLOWED_SECONDARY_ROLES' '=' '(' (('\'ALL\'' | ID (',' ID)*))? ')')? ('BLOCKED_SECONDARY_ROLES' '=' '(' (('\'ALL\'' | ID (',' ID)*))? ')')? ('AGENT_RESTRICTED_SESSION_SCOPE' '=' STRING)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-share
createShare
    : 'CREATE' ('OR' 'REPLACE')? 'SHARE' ID ('SECURE_OBJECTS_ONLY' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)?
    | 'CREATE' 'SHARE' ('IF' 'NOT' 'EXISTS')? ID ('SECURE_OBJECTS_ONLY' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)?
    | 'CREATE' 'OR' 'ALTER' 'SHARE' ID ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-snapshot
createSnapshot
    : 'CREATE' ('OR' 'REPLACE')? 'SNAPSHOT' QNAME 'FROM' 'SERVICE' QNAME 'VOLUME' ID 'INSTANCE' ID ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' 'SNAPSHOT' ('IF' 'NOT' 'EXISTS')? QNAME 'FROM' 'SERVICE' QNAME 'VOLUME' ID 'INSTANCE' ID ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-snapshot-policy
createSnapshotPolicy
    : 'CREATE' ('OR' 'REPLACE')? 'SNAPSHOT' 'POLICY' QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'RETENTION' 'LOCK')? ('SCHEDULE' '=' STRING)? ('EXPIRE_AFTER_DAYS' '=' NUMBER)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'SNAPSHOT' 'POLICY' ('IF' 'NOT' 'EXISTS')? QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'RETENTION' 'LOCK')? ('SCHEDULE' '=' STRING)? ('EXPIRE_AFTER_DAYS' '=' NUMBER)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-snapshot-set
createSnapshotSet
    : 'CREATE' ('OR' 'REPLACE')? 'SNAPSHOT' 'SET' QNAME 'FOR' 'DYNAMIC'? 'TABLE' QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'SNAPSHOT' 'POLICY' QNAME)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'SNAPSHOT' 'SET' ('IF' 'NOT' 'EXISTS')? QNAME 'FOR' 'DYNAMIC'? 'TABLE' QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'SNAPSHOT' 'POLICY' QNAME)? ('COMMENT' '=' STRING)?
    | 'CREATE' ('OR' 'REPLACE')? 'SNAPSHOT' 'SET' QNAME 'FOR' 'SCHEMA' QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'SNAPSHOT' 'POLICY' QNAME)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'SNAPSHOT' 'SET' ('IF' 'NOT' 'EXISTS')? QNAME 'FOR' 'SCHEMA' QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'SNAPSHOT' 'POLICY' QNAME)? ('COMMENT' '=' STRING)?
    | 'CREATE' ('OR' 'REPLACE')? 'SNAPSHOT' 'SET' QNAME 'FOR' 'DATABASE' ID ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'SNAPSHOT' 'POLICY' QNAME)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'SNAPSHOT' 'SET' ('IF' 'NOT' 'EXISTS')? QNAME 'FOR' 'DATABASE' ID ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'SNAPSHOT' 'POLICY' QNAME)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-stage
// fix: the bucket and path are text to fill in, as the parameters say
// fix: read: each -- If TYPE = X section of formatTypeOptions is the options of TYPE = X alone, and TYPE = { ... } [ formatTypeOptions ] chooses a TYPE with its own options
createStage
    : 'CREATE' ('OR' 'REPLACE')? (('TEMP' | 'TEMPORARY'))? 'STAGE' QNAME createStage_internalStageParams createStage_directoryTableParams ('FILE_FORMAT' '=' '(' ('FORMAT_NAME' '=' STRING | ('TYPE' '=' 'CSV' createStage_formatTypeOptions_CSV? | 'TYPE' '=' 'JSON' createStage_formatTypeOptions_JSON? | 'TYPE' '=' 'AVRO' createStage_formatTypeOptions_AVRO? | 'TYPE' '=' 'ORC' createStage_formatTypeOptions_ORC? | 'TYPE' '=' 'PARQUET' createStage_formatTypeOptions_PARQUET? | 'TYPE' '=' 'XML' createStage_formatTypeOptions_XML? | 'TYPE' '=' 'CUSTOM')) ')')? ('STAGE_UNDROP_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' (('TEMP' | 'TEMPORARY'))? 'STAGE' ('IF' 'NOT' 'EXISTS')? QNAME createStage_internalStageParams createStage_directoryTableParams ('FILE_FORMAT' '=' '(' ('FORMAT_NAME' '=' STRING | ('TYPE' '=' 'CSV' createStage_formatTypeOptions_CSV? | 'TYPE' '=' 'JSON' createStage_formatTypeOptions_JSON? | 'TYPE' '=' 'AVRO' createStage_formatTypeOptions_AVRO? | 'TYPE' '=' 'ORC' createStage_formatTypeOptions_ORC? | 'TYPE' '=' 'PARQUET' createStage_formatTypeOptions_PARQUET? | 'TYPE' '=' 'XML' createStage_formatTypeOptions_XML? | 'TYPE' '=' 'CUSTOM')) ')')? ('STAGE_UNDROP_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' ('OR' 'REPLACE')? (('TEMP' | 'TEMPORARY'))? 'STAGE' QNAME createStage_externalStageParams createStage_directoryTableParams ('FILE_FORMAT' '=' '(' ('FORMAT_NAME' '=' STRING | ('TYPE' '=' 'CSV' createStage_formatTypeOptions_CSV? | 'TYPE' '=' 'JSON' createStage_formatTypeOptions_JSON? | 'TYPE' '=' 'AVRO' createStage_formatTypeOptions_AVRO? | 'TYPE' '=' 'ORC' createStage_formatTypeOptions_ORC? | 'TYPE' '=' 'PARQUET' createStage_formatTypeOptions_PARQUET? | 'TYPE' '=' 'XML' createStage_formatTypeOptions_XML? | 'TYPE' '=' 'CUSTOM')) ')')? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' (('TEMP' | 'TEMPORARY'))? 'STAGE' ('IF' 'NOT' 'EXISTS')? QNAME createStage_externalStageParams createStage_directoryTableParams ('FILE_FORMAT' '=' '(' ('FORMAT_NAME' '=' STRING | ('TYPE' '=' 'CSV' createStage_formatTypeOptions_CSV? | 'TYPE' '=' 'JSON' createStage_formatTypeOptions_JSON? | 'TYPE' '=' 'AVRO' createStage_formatTypeOptions_AVRO? | 'TYPE' '=' 'ORC' createStage_formatTypeOptions_ORC? | 'TYPE' '=' 'PARQUET' createStage_formatTypeOptions_PARQUET? | 'TYPE' '=' 'XML' createStage_formatTypeOptions_XML? | 'TYPE' '=' 'CUSTOM')) ')')? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' 'OR' 'ALTER' (('TEMP' | 'TEMPORARY'))? 'STAGE' QNAME createStage_internalStageParams createStage_directoryTableParams ('FILE_FORMAT' '=' '(' ('FORMAT_NAME' '=' STRING | ('TYPE' '=' 'CSV' createStage_formatTypeOptions_CSV? | 'TYPE' '=' 'JSON' createStage_formatTypeOptions_JSON? | 'TYPE' '=' 'AVRO' createStage_formatTypeOptions_AVRO? | 'TYPE' '=' 'ORC' createStage_formatTypeOptions_ORC? | 'TYPE' '=' 'PARQUET' createStage_formatTypeOptions_PARQUET? | 'TYPE' '=' 'XML' createStage_formatTypeOptions_XML? | 'TYPE' '=' 'CUSTOM')) ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' 'OR' 'ALTER' (('TEMP' | 'TEMPORARY'))? 'STAGE' QNAME createStage_externalStageParams createStage_directoryTableParams ('FILE_FORMAT' '=' '(' ('FORMAT_NAME' '=' STRING | ('TYPE' '=' 'CSV' createStage_formatTypeOptions_CSV? | 'TYPE' '=' 'JSON' createStage_formatTypeOptions_JSON? | 'TYPE' '=' 'AVRO' createStage_formatTypeOptions_AVRO? | 'TYPE' '=' 'ORC' createStage_formatTypeOptions_ORC? | 'TYPE' '=' 'PARQUET' createStage_formatTypeOptions_PARQUET? | 'TYPE' '=' 'XML' createStage_formatTypeOptions_XML? | 'TYPE' '=' 'CUSTOM')) ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' ('OR' 'REPLACE')? 'STAGE' QNAME 'CLONE' QNAME
    | 'CREATE' 'STAGE' ('IF' 'NOT' 'EXISTS')? QNAME 'CLONE' QNAME
    ;
createStage_internalStageParams
    : ('ENCRYPTION' '=' '(' ('TYPE' '=' '\'SNOWFLAKE_FULL\'' | 'TYPE' '=' '\'SNOWFLAKE_SSE\'') ')')?
    ;
createStage_externalStageParams
    : 'URL' '=' STRING ('AWS_ACCESS_POINT_ARN' '=' STRING)? ('STORAGE_INTEGRATION' '=' ID | 'CREDENTIALS' '=' '(' ('AWS_KEY_ID' '=' STRING 'AWS_SECRET_KEY' '=' STRING ('AWS_TOKEN' '=' STRING)? | 'AWS_ROLE' '=' STRING) ')')? ('ENCRYPTION' '=' '(' (('TYPE' '=' '\'AWS_CSE\'')? 'MASTER_KEY' '=' STRING | 'TYPE' '=' '\'AWS_SSE_S3\'' | 'TYPE' '=' '\'AWS_SSE_KMS\'' ('KMS_KEY_ID' '=' STRING)? | 'TYPE' '=' '\'NONE\'') ')')? ('USE_PRIVATELINK_ENDPOINT' '=' ('TRUE' | 'FALSE'))?
    | 'URL' '=' STRING ('STORAGE_INTEGRATION' '=' ID)? ('ENCRYPTION' '=' '(' ('TYPE' '=' '\'GCS_SSE_KMS\'' ('KMS_KEY_ID' '=' STRING)? | 'TYPE' '=' '\'NONE\'') ')')? ('USE_PRIVATELINK_ENDPOINT' '=' ('TRUE' | 'FALSE'))?
    | 'URL' '=' STRING ('STORAGE_INTEGRATION' '=' ID | 'CREDENTIALS' '=' '(' ('AZURE_SAS_TOKEN' '=' STRING)? ')')? ('ENCRYPTION' '=' '(' ('TYPE' '=' '\'AZURE_CSE\'' 'MASTER_KEY' '=' STRING | 'TYPE' '=' '\'NONE\'') ')')? ('USE_PRIVATELINK_ENDPOINT' '=' ('TRUE' | 'FALSE'))?
    | 'URL' '=' STRING ('STORAGE_INTEGRATION' '=' ID | 'CREDENTIALS' '=' '(' ('AZURE_SAS_TOKEN' '=' STRING)? ')')? ('ENCRYPTION' '=' '(' ('TYPE' '=' '\'AZURE_CSE\'' 'MASTER_KEY' '=' STRING | 'TYPE' '=' '\'NONE\'') ')')?
    | 'URL' '=' STRING 'ENDPOINT' '=' STRING ('CREDENTIALS' '=' '(' 'AWS_KEY_ID' '=' STRING 'AWS_SECRET_KEY' '=' STRING ')')?
    ;
createStage_directoryTableParams
    : ('DIRECTORY' '=' '(' 'ENABLE' '=' ('TRUE' | 'FALSE') ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ')')?
    | ('DIRECTORY' '=' '(' 'ENABLE' '=' ('TRUE' | 'FALSE') ('REFRESH_ON_CREATE' '=' ('TRUE' | 'FALSE'))? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('AWS_SNS_TOPIC' '=' STRING)? ')')?
    | ('DIRECTORY' '=' '(' 'ENABLE' '=' ('TRUE' | 'FALSE') ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('REFRESH_ON_CREATE' '=' ('TRUE' | 'FALSE'))? ('NOTIFICATION_INTEGRATION' '=' STRING)? ')')?
    | ('DIRECTORY' '=' '(' 'ENABLE' '=' ('TRUE' | 'FALSE') ('REFRESH_ON_CREATE' '=' ('TRUE' | 'FALSE'))? ('AUTO_REFRESH' '=' ('TRUE' | 'FALSE'))? ('NOTIFICATION_INTEGRATION' '=' STRING)? ')')?
    ;
createStage_formatTypeOptions_CSV
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'RECORD_DELIMITER' '=' (STRING | 'NONE') | 'FIELD_DELIMITER' '=' (STRING | 'NONE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'FILE_EXTENSION' '=' STRING | 'PARSE_HEADER' '=' ('TRUE' | 'FALSE') | 'SKIP_HEADER' '=' NUMBER | 'SKIP_BLANK_LINES' '=' ('TRUE' | 'FALSE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'ESCAPE' '=' (STRING | 'NONE') | 'ESCAPE_UNENCLOSED_FIELD' '=' (STRING | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'FIELD_OPTIONALLY_ENCLOSED_BY' '=' (STRING | 'NONE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')' | 'ERROR_ON_COLUMN_COUNT_MISMATCH' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'EMPTY_FIELD_AS_NULL' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE') | 'ENCODING' '=' (STRING | 'UTF8')) (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'RECORD_DELIMITER' '=' (STRING | 'NONE') | 'FIELD_DELIMITER' '=' (STRING | 'NONE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'FILE_EXTENSION' '=' STRING | 'PARSE_HEADER' '=' ('TRUE' | 'FALSE') | 'SKIP_HEADER' '=' NUMBER | 'SKIP_BLANK_LINES' '=' ('TRUE' | 'FALSE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'ESCAPE' '=' (STRING | 'NONE') | 'ESCAPE_UNENCLOSED_FIELD' '=' (STRING | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'FIELD_OPTIONALLY_ENCLOSED_BY' '=' (STRING | 'NONE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')' | 'ERROR_ON_COLUMN_COUNT_MISMATCH' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'EMPTY_FIELD_AS_NULL' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE') | 'ENCODING' '=' (STRING | 'UTF8')))*
    ;
createStage_formatTypeOptions_JSON
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')' | 'FILE_EXTENSION' '=' STRING | 'ENABLE_OCTAL' '=' ('TRUE' | 'FALSE') | 'ALLOW_DUPLICATE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ARRAY' '=' ('TRUE' | 'FALSE') | 'STRIP_NULL_VALUES' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'IGNORE_UTF8_ERRORS' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE')) (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'DATE_FORMAT' '=' (STRING | 'AUTO') | 'TIME_FORMAT' '=' (STRING | 'AUTO') | 'TIMESTAMP_FORMAT' '=' (STRING | 'AUTO') | 'BINARY_FORMAT' '=' ('HEX' | 'BASE64' | 'UTF8') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'MULTI_LINE' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')' | 'FILE_EXTENSION' '=' STRING | 'ENABLE_OCTAL' '=' ('TRUE' | 'FALSE') | 'ALLOW_DUPLICATE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ARRAY' '=' ('TRUE' | 'FALSE') | 'STRIP_NULL_VALUES' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'IGNORE_UTF8_ERRORS' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE')))*
    ;
createStage_formatTypeOptions_AVRO
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')') (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')'))*
    ;
createStage_formatTypeOptions_ORC
    : ('TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')') (','? ('TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')'))*
    ;
createStage_formatTypeOptions_PARQUET
    : ('COMPRESSION' '=' ('AUTO' | 'LZO' | 'SNAPPY' | 'NONE') | 'SNAPPY_COMPRESSION' '=' ('TRUE' | 'FALSE') | 'BINARY_AS_TEXT' '=' ('TRUE' | 'FALSE') | 'USE_LOGICAL_TYPE' '=' ('TRUE' | 'FALSE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'USE_VECTORIZED_SCANNER' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')') (','? ('COMPRESSION' '=' ('AUTO' | 'LZO' | 'SNAPPY' | 'NONE') | 'SNAPPY_COMPRESSION' '=' ('TRUE' | 'FALSE') | 'BINARY_AS_TEXT' '=' ('TRUE' | 'FALSE') | 'USE_LOGICAL_TYPE' '=' ('TRUE' | 'FALSE') | 'TRIM_SPACE' '=' ('TRUE' | 'FALSE') | 'USE_VECTORIZED_SCANNER' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'NULL_IF' '=' '(' STRING (',' STRING)* ')'))*
    ;
createStage_formatTypeOptions_XML
    : ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'IGNORE_UTF8_ERRORS' '=' ('TRUE' | 'FALSE') | 'PRESERVE_SPACE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ELEMENT' '=' ('TRUE' | 'FALSE') | 'DISABLE_AUTO_CONVERT' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE')) (','? ('COMPRESSION' '=' ('AUTO' | 'GZIP' | 'BZ2' | 'BROTLI' | 'ZSTD' | 'DEFLATE' | 'RAW_DEFLATE' | 'NONE') | 'IGNORE_UTF8_ERRORS' '=' ('TRUE' | 'FALSE') | 'PRESERVE_SPACE' '=' ('TRUE' | 'FALSE') | 'STRIP_OUTER_ELEMENT' '=' ('TRUE' | 'FALSE') | 'DISABLE_AUTO_CONVERT' '=' ('TRUE' | 'FALSE') | 'REPLACE_INVALID_CHARACTERS' '=' ('TRUE' | 'FALSE') | 'SKIP_BYTE_ORDER_MARK' '=' ('TRUE' | 'FALSE')))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-storage-integration
// fix: the S3 providers include S3CHINA and S3GOV, as the cloud provider parameters say
createStorageIntegration
    : 'CREATE' ('OR' 'REPLACE')? 'STORAGE' 'INTEGRATION' ID 'TYPE' '=' ('EXTERNAL_STAGE' | 'POSTGRES_EXTERNAL_STORAGE' | 'POSTGRES_INTERNAL_STORAGE') createStorageIntegration_cloudProviderParams 'ENABLED' '=' ('TRUE' | 'FALSE') 'STORAGE_ALLOWED_LOCATIONS' '=' '(' STRING (',' STRING)* ')' ('STORAGE_BLOCKED_LOCATIONS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)?
    | 'CREATE' 'STORAGE' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'TYPE' '=' ('EXTERNAL_STAGE' | 'POSTGRES_EXTERNAL_STORAGE' | 'POSTGRES_INTERNAL_STORAGE') createStorageIntegration_cloudProviderParams 'ENABLED' '=' ('TRUE' | 'FALSE') 'STORAGE_ALLOWED_LOCATIONS' '=' '(' STRING (',' STRING)* ')' ('STORAGE_BLOCKED_LOCATIONS' '=' '(' STRING (',' STRING)* ')')? ('COMMENT' '=' STRING)?
    ;
createStorageIntegration_cloudProviderParams
    : 'STORAGE_PROVIDER' '=' ('\'S3\'' | '\'S3CHINA\'' | '\'S3GOV\'') 'STORAGE_AWS_ROLE_ARN' '=' STRING ('STORAGE_AWS_EXTERNAL_ID' '=' STRING)? ('STORAGE_AWS_OBJECT_ACL' '=' '\'bucket-owner-full-control\'')? ('USE_PRIVATELINK_ENDPOINT' '=' ('TRUE' | 'FALSE'))?
    | 'STORAGE_PROVIDER' '=' '\'GCS\''
    | 'STORAGE_PROVIDER' '=' '\'AZURE\'' 'AZURE_TENANT_ID' '=' STRING ('USE_PRIVATELINK_ENDPOINT' '=' ('TRUE' | 'FALSE'))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-storage-integration-postgres-internal
createStorageIntegrationPostgresInternal
    : 'CREATE' ('OR' 'REPLACE')? 'STORAGE' 'INTEGRATION' ID 'TYPE' '=' 'POSTGRES_INTERNAL_STORAGE' 'POSTGRES_INSTANCE' '=' STRING ('ENABLED' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)?
    | 'CREATE' 'STORAGE' 'INTEGRATION' ('IF' 'NOT' 'EXISTS')? ID 'TYPE' '=' 'POSTGRES_INTERNAL_STORAGE' 'POSTGRES_INSTANCE' '=' STRING ('ENABLED' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-storage-lifecycle-policy
createStorageLifecyclePolicy
    : 'CREATE' ('OR' 'REPLACE')? 'STORAGE' 'LIFECYCLE' 'POLICY' QNAME 'AS' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')' 'RETURNS' 'BOOLEAN' '->' EXPR ('ARCHIVE_TIER' '=' ('COOL' | 'COLD'))? ('ARCHIVE_FOR_DAYS' '=' NUMBER)? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' 'STORAGE' 'LIFECYCLE' 'POLICY' ('IF' 'NOT' 'EXISTS')? QNAME 'AS' '(' ID DATA_TYPE (',' ID DATA_TYPE)* ')' 'RETURNS' 'BOOLEAN' '->' EXPR ('ARCHIVE_TIER' '=' ('COOL' | 'COLD'))? ('ARCHIVE_FOR_DAYS' '=' NUMBER)? ('COMMENT' '=' STRING)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-stream
// fix: read: Time Travel's OFFSET is a negative number of seconds, -N, and its STATEMENT a quoted query ID, as the Time Travel parameters say
createStream
    : 'CREATE' ('OR' 'REPLACE')? 'STREAM' QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COPY' 'GRANTS')? 'ON' 'TABLE' QNAME (('AT' | 'BEFORE') '(' ('TIMESTAMP' '=>' STRING | 'OFFSET' '=>' '-' NUMBER | 'STATEMENT' '=>' STRING | 'STREAM' '=>' STRING) ')')? ('APPEND_ONLY' '=' 'TRUE' | 'APPEND_ONLY' '=' 'FALSE')? ('SHOW_INITIAL_ROWS' '=' 'TRUE' | 'SHOW_INITIAL_ROWS' '=' 'FALSE')? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')?
    | 'CREATE' 'STREAM' ('IF' 'NOT' 'EXISTS')? QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COPY' 'GRANTS')? 'ON' 'TABLE' QNAME (('AT' | 'BEFORE') '(' ('TIMESTAMP' '=>' STRING | 'OFFSET' '=>' '-' NUMBER | 'STATEMENT' '=>' STRING | 'STREAM' '=>' STRING) ')')? ('APPEND_ONLY' '=' 'TRUE' | 'APPEND_ONLY' '=' 'FALSE')? ('SHOW_INITIAL_ROWS' '=' 'TRUE' | 'SHOW_INITIAL_ROWS' '=' 'FALSE')? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'STREAM' QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COPY' 'GRANTS')? 'ON' 'EVENT' 'TABLE' QNAME ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')?
    | 'CREATE' 'STREAM' ('IF' 'NOT' 'EXISTS')? QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COPY' 'GRANTS')? 'ON' 'EVENT' 'TABLE' QNAME ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'STREAM' QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COPY' 'GRANTS')? 'ON' 'EXTERNAL' 'TABLE' QNAME (('AT' | 'BEFORE') '(' ('TIMESTAMP' '=>' STRING | 'OFFSET' '=>' '-' NUMBER | 'STATEMENT' '=>' STRING | 'STREAM' '=>' STRING) ')')? ('INSERT_ONLY' '=' 'TRUE')? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')?
    | 'CREATE' 'STREAM' ('IF' 'NOT' 'EXISTS')? QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COPY' 'GRANTS')? 'ON' 'EXTERNAL' 'TABLE' QNAME (('AT' | 'BEFORE') '(' ('TIMESTAMP' '=>' STRING | 'OFFSET' '=>' '-' NUMBER | 'STATEMENT' '=>' STRING | 'STREAM' '=>' STRING) ')')? ('INSERT_ONLY' '=' 'TRUE')? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'STREAM' QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COPY' 'GRANTS')? 'ON' 'STAGE' QNAME ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')?
    | 'CREATE' 'STREAM' ('IF' 'NOT' 'EXISTS')? QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COPY' 'GRANTS')? 'ON' 'STAGE' QNAME ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'STREAM' QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COPY' 'GRANTS')? 'ON' 'DYNAMIC' 'TABLE' QNAME ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')?
    | 'CREATE' 'STREAM' ('IF' 'NOT' 'EXISTS')? QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COPY' 'GRANTS')? 'ON' 'DYNAMIC' 'TABLE' QNAME ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'STREAM' QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COPY' 'GRANTS')? 'ON' 'VIEW' QNAME (('AT' | 'BEFORE') '(' ('TIMESTAMP' '=>' STRING | 'OFFSET' '=>' '-' NUMBER | 'STATEMENT' '=>' STRING | 'STREAM' '=>' STRING) ')')? ('APPEND_ONLY' '=' 'TRUE' | 'APPEND_ONLY' '=' 'FALSE')? ('SHOW_INITIAL_ROWS' '=' 'TRUE' | 'SHOW_INITIAL_ROWS' '=' 'FALSE')? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')?
    | 'CREATE' 'STREAM' ('IF' 'NOT' 'EXISTS')? QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COPY' 'GRANTS')? 'ON' 'VIEW' QNAME (('AT' | 'BEFORE') '(' ('TIMESTAMP' '=>' STRING | 'OFFSET' '=>' '-' NUMBER | 'STATEMENT' '=>' STRING | 'STREAM' '=>' STRING) ')')? ('APPEND_ONLY' '=' 'TRUE' | 'APPEND_ONLY' '=' 'FALSE')? ('SHOW_INITIAL_ROWS' '=' 'TRUE' | 'SHOW_INITIAL_ROWS' '=' 'FALSE')? ('COMMENT' '=' STRING)? ('COPY' '=' '(' 'GRANTS' ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'STREAM' QNAME 'CLONE' QNAME ('COPY' 'GRANTS')? ('COPY' '=' '(' 'GRANTS' ')')?
    | 'CREATE' 'OR' 'ALTER' 'STREAM' QNAME 'ON' 'TABLE' QNAME ('APPEND_ONLY' '=' 'TRUE' | 'APPEND_ONLY' '=' 'FALSE')? ('SHOW_INITIAL_ROWS' '=' 'TRUE' | 'SHOW_INITIAL_ROWS' '=' 'FALSE')? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-streamlit
// fix: the source is a quoted stage path, as the parameters and examples write it
// not converted (form): DEFINE runs only in DCM project files, not as a statement
//   DEFINE STREAMLIT DEMO{{env_suffix}}.SERVE.MY_DASHBOARD
//       FROM 'asset://my_dashboard/'
//       MAIN_FILE = 'streamlit_app.py'
//       QUERY_WAREHOUSE = DEMO_WH{{env_suffix}}
//   ;
createStreamlit
    : 'CREATE' ('OR' 'REPLACE')? 'STREAMLIT' QNAME ('FROM' STRING)? ('MAIN_FILE' '=' STRING)? ('QUERY_WAREHOUSE' '=' ID)? ('RUNTIME_NAME' '=' STRING)? ('COMPUTE_POOL' '=' ID)? ('COMMENT' '=' STRING)? ('TITLE' '=' STRING)? ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' ID (',' STRING '=' ID)* ')')?
    | 'CREATE' 'STREAMLIT' ('IF' 'NOT' 'EXISTS')? QNAME ('FROM' STRING)? ('MAIN_FILE' '=' STRING)? ('QUERY_WAREHOUSE' '=' ID)? ('RUNTIME_NAME' '=' STRING)? ('COMPUTE_POOL' '=' ID)? ('COMMENT' '=' STRING)? ('TITLE' '=' STRING)? ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' STRING '=' ID (',' STRING '=' ID)* ')')?
    | 'CREATE' ('OR' 'REPLACE')? 'STREAMLIT' QNAME 'ROOT_LOCATION' '=' STRING 'MAIN_FILE' '=' STRING ('QUERY_WAREHOUSE' '=' ID)? ('COMMENT' '=' STRING)? ('TITLE' '=' STRING)? ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')?
    | 'CREATE' 'STREAMLIT' ('IF' 'NOT' 'EXISTS')? QNAME 'ROOT_LOCATION' '=' STRING 'MAIN_FILE' '=' STRING ('QUERY_WAREHOUSE' '=' ID)? ('COMMENT' '=' STRING)? ('TITLE' '=' STRING)? ('IMPORTS' '=' '(' STRING (',' STRING)* ')')? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-table
// fix: the template is a parenthesized subquery, as the parameters and examples write it
// fix: read: Time Travel's OFFSET is a negative number of seconds, -N, and its STATEMENT a quoted query ID, as the Time Travel parameters say
// fix: read: <constraint_properties> is the properties CREATE TABLE's constraint page lists
// fix: read: <constraint_properties> is the properties CREATE TABLE's constraint page lists
// fix: read: an out-of-line key constraint names its columns, as the constraint page says
// fix: read: a foreign key's REFERENCES takes MATCH and ON UPDATE/DELETE, the properties the constraint page gives foreign keys
// fix: read: the additional columns take the whole column definition, which their line abbreviates
// not converted (form): ... stands for text the docs leave out
//   CREATE TABLE <table_name> ( <col1> <data_type> [ WITH ] MASKING POLICY <policy_name> [ , ... ] )
//     ...
//     [ WITH ] ROW ACCESS POLICY <policy_name> ON ( <col1> [ , ... ] )
//     [ ... ]
//     AS <query>
createTable
    : 'CREATE' ('OR' 'REPLACE')? (((('LOCAL' | 'GLOBAL'))? 'TEMP' | 'TEMPORARY' | 'VOLATILE' | 'TRANSIENT'))? 'TABLE' QNAME '(' ID DATA_TYPE (('GENERATED' 'ALWAYS')? 'AS' '(' EXPR ')' 'VIRTUAL'?)? createTable_inlineConstraint? ('NOT' 'NULL')? ('COLLATE' STRING)? (('DEFAULT' EXPR | ('AUTOINCREMENT' | 'IDENTITY') (('(' NUMBER ',' NUMBER ')' | 'START' NUMBER 'INCREMENT' NUMBER))? (('ORDER' | 'NOORDER'))?))? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)? (',' ID DATA_TYPE (('GENERATED' 'ALWAYS')? 'AS' '(' EXPR ')' 'VIRTUAL'?)? createTable_inlineConstraint? ('NOT' 'NULL')? ('COLLATE' STRING)? (('DEFAULT' EXPR | ('AUTOINCREMENT' | 'IDENTITY') (('(' NUMBER ',' NUMBER ')' | 'START' NUMBER 'INCREMENT' NUMBER))? (('ORDER' | 'NOORDER'))?))? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)?)* (',' createTable_outoflineConstraint)* ')' ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')')? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('CHANGE_TRACKING' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('COPY' 'GRANTS')? ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('COPY' 'TAGS')? ('COPY' '=' '(' ('COMMENTS' | 'TAGS' | 'GRANTS') (',' ('COMMENTS' | 'TAGS' | 'GRANTS'))* ')')? ('COMMENT' '=' STRING)? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'AGGREGATION' 'POLICY' QNAME ('ENTITY' 'KEY' '(' ID (',' ID)* ')')?)? ('WITH'? 'JOIN' 'POLICY' QNAME ('ALLOWED' 'JOIN' 'KEYS' '(' ID (',' ID)* ')')?)? ('WITH'? 'STORAGE' 'LIFECYCLE' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('ROW_TIMESTAMP' '=' ('TRUE' | 'FALSE'))?
    | 'CREATE' (((('LOCAL' | 'GLOBAL'))? 'TEMP' | 'TEMPORARY' | 'VOLATILE' | 'TRANSIENT'))? 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME '(' ID DATA_TYPE (('GENERATED' 'ALWAYS')? 'AS' '(' EXPR ')' 'VIRTUAL'?)? createTable_inlineConstraint? ('NOT' 'NULL')? ('COLLATE' STRING)? (('DEFAULT' EXPR | ('AUTOINCREMENT' | 'IDENTITY') (('(' NUMBER ',' NUMBER ')' | 'START' NUMBER 'INCREMENT' NUMBER))? (('ORDER' | 'NOORDER'))?))? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)? (',' ID DATA_TYPE (('GENERATED' 'ALWAYS')? 'AS' '(' EXPR ')' 'VIRTUAL'?)? createTable_inlineConstraint? ('NOT' 'NULL')? ('COLLATE' STRING)? (('DEFAULT' EXPR | ('AUTOINCREMENT' | 'IDENTITY') (('(' NUMBER ',' NUMBER ')' | 'START' NUMBER 'INCREMENT' NUMBER))? (('ORDER' | 'NOORDER'))?))? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('COMMENT' STRING)?)* (',' createTable_outoflineConstraint)* ')' ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')')? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('CHANGE_TRACKING' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('COPY' 'GRANTS')? ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('COPY' 'TAGS')? ('COPY' '=' '(' ('COMMENTS' | 'TAGS' | 'GRANTS') (',' ('COMMENTS' | 'TAGS' | 'GRANTS'))* ')')? ('COMMENT' '=' STRING)? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'AGGREGATION' 'POLICY' QNAME ('ENTITY' 'KEY' '(' ID (',' ID)* ')')?)? ('WITH'? 'JOIN' 'POLICY' QNAME ('ALLOWED' 'JOIN' 'KEYS' '(' ID (',' ID)* ')')?)? ('WITH'? 'STORAGE' 'LIFECYCLE' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('ROW_TIMESTAMP' '=' ('TRUE' | 'FALSE'))?
    | 'CREATE' 'TABLE' QNAME 'FROM' 'BACKUP' 'SET' ID 'IDENTIFIER' STRING
    | 'CREATE' 'OR' 'ALTER' (((('LOCAL' | 'GLOBAL'))? 'TEMP' | 'TEMPORARY' | 'TRANSIENT'))? 'TABLE' QNAME '(' ID DATA_TYPE createTable_inlineConstraint? ('NOT' 'NULL')? ('COLLATE' STRING)? (('DEFAULT' EXPR | ('AUTOINCREMENT' | 'IDENTITY') (('(' NUMBER ',' NUMBER ')' | 'START' NUMBER 'INCREMENT' NUMBER))? (('ORDER' | 'NOORDER'))?))? ('COMMENT' STRING)? (',' ID DATA_TYPE createTable_inlineConstraint? ('NOT' 'NULL')? ('COLLATE' STRING)? (('DEFAULT' EXPR | ('AUTOINCREMENT' | 'IDENTITY') (('(' NUMBER ',' NUMBER ')' | 'START' NUMBER 'INCREMENT' NUMBER))? (('ORDER' | 'NOORDER'))?))? ('COMMENT' STRING)?)* (',' createTable_outoflineConstraint)* ')' ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')')? ('ENABLE_SCHEMA_EVOLUTION' '=' ('TRUE' | 'FALSE'))? ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('CHANGE_TRACKING' '=' ('TRUE' | 'FALSE'))? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('ICEBERG_DEFAULT_DDL_COLLATION' '=' STRING)? ('ERROR_LOGGING' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? ('ROW_TIMESTAMP' '=' ('TRUE' | 'FALSE'))?
    | 'CREATE' ('OR' 'REPLACE')? 'TABLE' QNAME ('(' ID DATA_TYPE? (',' ID DATA_TYPE?)* ')')? ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')')? ('COPY' 'GRANTS')? ('COPY' 'TAGS')? ('COPY' '=' '(' ('COMMENTS' | 'TAGS' | 'GRANTS') (',' ('COMMENTS' | 'TAGS' | 'GRANTS'))* ')')* 'AS' QUERY
    | 'CREATE' ('OR' 'REPLACE')? 'TABLE' QNAME ('COPY' 'GRANTS')? 'USING' 'TEMPLATE' '(' QUERY ')'
    | 'CREATE' ('OR' 'REPLACE')? 'TABLE' QNAME 'LIKE' QNAME ('CLUSTER' 'BY' '(' EXPR (',' EXPR)* ')')? ('COPY' 'GRANTS')? ('COPY' 'TAGS')? ('COPY' '=' '(' ('COMMENTS' | 'TAGS' | 'GRANTS') (',' ('COMMENTS' | 'TAGS' | 'GRANTS'))* ')')?
    | 'CREATE' ('OR' 'REPLACE')? (((('LOCAL' | 'GLOBAL'))? 'TEMP' ('READ' 'ONLY')? | 'TEMPORARY' ('READ' 'ONLY')? | 'VOLATILE' | 'TRANSIENT'))? 'TABLE' QNAME 'CLONE' QNAME (('AT' | 'BEFORE') '(' ('TIMESTAMP' '=>' STRING | 'OFFSET' '=>' '-' NUMBER | 'STATEMENT' '=>' STRING) ')')? ('COPY' 'GRANTS')? ('COPY' 'TAGS')? ('COPY' '=' '(' ('COMMENTS' | 'TAGS' | 'GRANTS') (',' ('COMMENTS' | 'TAGS' | 'GRANTS'))* ')')?
    | 'CREATE' 'TRANSIENT'? 'TABLE' ('IF' 'NOT' 'EXISTS')? QNAME 'FROM' 'ARCHIVE' 'OF' QNAME ('AS'? ID)? 'WHERE' EXPR
    ;
createTable_inlineConstraint
    : ('CONSTRAINT' ID)? ('UNIQUE' | 'PRIMARY' 'KEY' | ('FOREIGN' 'KEY')? 'REFERENCES' QNAME ('(' ID ')')? ('MATCH' ('FULL' | 'SIMPLE' | 'PARTIAL'))? ('ON' ('UPDATE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))? ('DELETE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))?)? | 'CHECK' '(' EXPR ')') ('NOT'? 'ENFORCED')? ('NOT'? 'DEFERRABLE')? ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? (('RELY' | 'NORELY'))?
    ;
createTable_outoflineConstraint
    : ('CONSTRAINT' ID)? ('UNIQUE' '(' ID (',' ID)* ')' | 'PRIMARY' 'KEY' '(' ID (',' ID)* ')' | ('FOREIGN' 'KEY')? '(' ID (',' ID)* ')' 'REFERENCES' QNAME ('(' ID (',' ID)* ')')? ('MATCH' ('FULL' | 'SIMPLE' | 'PARTIAL'))? ('ON' ('UPDATE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))? ('DELETE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))?)? | 'CHECK' '(' EXPR ')') ('NOT'? 'ENFORCED')? ('NOT'? 'DEFERRABLE')? ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? (('RELY' | 'NORELY'))? ('COMMENT' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-table-constraint
// fix: ... stands for the table's other clauses; the form is read without them
// fix: ... stands for the table's other clauses; the form is read without them
createTableConstraint
    : 'CREATE' 'TABLE' QNAME '(' ID DATA_TYPE ('NOT' 'NULL')? (createTableConstraint_inlineUniquePK | createTableConstraint_inlineFK | createTableConstraint_inlineCH) (',' ID DATA_TYPE ('NOT' 'NULL')? (createTableConstraint_inlineUniquePK | createTableConstraint_inlineFK | createTableConstraint_inlineCH))* ')'
    | 'ALTER' 'TABLE' QNAME 'ADD' 'COLUMN' ID DATA_TYPE ('NOT' 'NULL')? (createTableConstraint_inlineUniquePK | createTableConstraint_inlineFK | createTableConstraint_inlineCH)
    | 'CREATE' 'TABLE' QNAME '(' ID DATA_TYPE (',' ID DATA_TYPE)* (',' (createTableConstraint_outoflineUniquePK | createTableConstraint_outoflineFK | createTableConstraint_outoflineCH))? (',' (createTableConstraint_outoflineUniquePK | createTableConstraint_outoflineFK | createTableConstraint_outoflineCH))* ')'
    | 'ALTER' 'TABLE' QNAME 'ADD' (createTableConstraint_outoflineUniquePK | createTableConstraint_outoflineFK | createTableConstraint_outoflineCH)
    ;
createTableConstraint_inlineUniquePK
    : ('CONSTRAINT' ID)? ('UNIQUE' | 'PRIMARY' 'KEY') ('NOT'? 'ENFORCED')? ('NOT'? 'DEFERRABLE')? ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? (('RELY' | 'NORELY'))?
    ;
createTableConstraint_inlineFK
    : ('CONSTRAINT' ID)? ('FOREIGN' 'KEY')? 'REFERENCES' QNAME ('(' ID ')')? ('MATCH' ('FULL' | 'SIMPLE' | 'PARTIAL'))? ('ON' ('UPDATE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))? ('DELETE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))?)? ('NOT'? 'ENFORCED')? ('NOT'? 'DEFERRABLE')? ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? (('RELY' | 'NORELY'))?
    ;
createTableConstraint_inlineCH
    : ('CONSTRAINT' ID)? 'CHECK' '(' EXPR ')' ('ENABLE' ('VALIDATE' | 'NOVALIDATE'))?
    ;
createTableConstraint_outoflineUniquePK
    : ('CONSTRAINT' ID)? ('UNIQUE' | 'PRIMARY' 'KEY') '(' ID (',' ID)* ')' ('NOT'? 'ENFORCED')? ('NOT'? 'DEFERRABLE')? ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? (('RELY' | 'NORELY'))? ('COMMENT' STRING)?
    ;
createTableConstraint_outoflineFK
    : ('CONSTRAINT' ID)? 'FOREIGN' 'KEY' '(' ID (',' ID)* ')' 'REFERENCES' QNAME ('(' ID (',' ID)* ')')? ('MATCH' ('FULL' | 'SIMPLE' | 'PARTIAL'))? ('ON' ('UPDATE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))? ('DELETE' ('CASCADE' | 'SET' 'NULL' | 'SET' 'DEFAULT' | 'RESTRICT' | 'NO' 'ACTION'))?)? ('NOT'? 'ENFORCED')? ('NOT'? 'DEFERRABLE')? ('INITIALLY' ('DEFERRED' | 'IMMEDIATE'))? (('ENABLE' | 'DISABLE'))? (('VALIDATE' | 'NOVALIDATE'))? (('RELY' | 'NORELY'))? ('COMMENT' STRING)?
    ;
createTableConstraint_outoflineCH
    : ('CONSTRAINT' ID)? 'CHECK' '(' EXPR ')' ('ENABLE' ('VALIDATE' | 'NOVALIDATE'))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-tag
createTag
    : 'CREATE' ('OR' 'REPLACE')? 'TAG' QNAME ('ALLOWED_VALUES' STRING (',' STRING)*)? ('MULTI_VALUE' '=' 'TRUE')? ('PROPAGATE' '=' ('ON_DEPENDENCY_AND_DATA_MOVEMENT' | 'ON_DEPENDENCY' | 'ON_DATA_MOVEMENT') ('ON_CONFLICT' '=' (STRING | 'ALLOWED_VALUES_SEQUENCE' | 'MERGE'))?)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'TAG' ('IF' 'NOT' 'EXISTS')? QNAME ('ALLOWED_VALUES' STRING (',' STRING)*)? ('MULTI_VALUE' '=' 'TRUE')? ('PROPAGATE' '=' ('ON_DEPENDENCY_AND_DATA_MOVEMENT' | 'ON_DEPENDENCY' | 'ON_DATA_MOVEMENT') ('ON_CONFLICT' '=' (STRING | 'ALLOWED_VALUES_SEQUENCE' | 'MERGE'))?)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'OR' 'ALTER' 'TAG' QNAME ('ALLOWED_VALUES' STRING (',' STRING)*)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-task
// fix: SERVERLESS_TASK_MIN/MAX_STATEMENT_SIZE: the choice of sizes is written inside the quotes; read as a string
// fix: AFTER names predecessor tasks, as the parameters and examples write them
// fix: the warehouse is a name, as the parameters say
// fix: a session parameter takes a constant value, as the parameters say
// not converted (form): DEFINE runs only in DCM project files, not as a statement
//   DEFINE TASK MY_DB.MY_SCHEMA.TSK_INGEST_DAILY_ORDERS
//       WAREHOUSE = 'MY_WH'
//       SCHEDULE = 'USING CRON 0 5 * * * UTC'
//       STARTED
//   AS
//       SELECT 1
//   ;
createTask
    : 'CREATE' ('OR' 'REPLACE')? 'TASK' QNAME ('WITH' 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('WAREHOUSE' '=' ID | 'USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE' '=' STRING)? ('SCHEDULE' '=' STRING)? ('CONFIG' '=' STRING)? ('OVERLAP_POLICY' '=' ('NO_OVERLAP' | 'ALLOW_CHILD_OVERLAP' | 'ALLOW_ALL_OVERLAP'))? (ID '=' EXPR (',' ID '=' EXPR)*)? ('USER_TASK_TIMEOUT_MS' '=' NUMBER)? ('SUSPEND_TASK_AFTER_NUM_FAILURES' '=' NUMBER)? ('ERROR_INTEGRATION' '=' ID)? ('SUCCESS_INTEGRATION' '=' ID)? ('LOG_LEVEL' '=' STRING)? ('COMMENT' '=' STRING)? ('FINALIZE' '=' STRING)? ('TASK_AUTO_RETRY_ATTEMPTS' '=' NUMBER)? ('USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS' '=' NUMBER)? ('TARGET_COMPLETION_INTERVAL' '=' STRING)? ('SERVERLESS_TASK_MIN_STATEMENT_SIZE' '=' STRING)? ('SERVERLESS_TASK_MAX_STATEMENT_SIZE' '=' STRING)? ('AFTER' QNAME (',' QNAME)*)? ('EXECUTE' 'AS' 'USER' ID)? ('WHEN' EXPR)? 'AS' QUERY
    | 'CREATE' 'TASK' ('IF' 'NOT' 'EXISTS')? QNAME ('WITH' 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('WAREHOUSE' '=' ID | 'USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE' '=' STRING)? ('SCHEDULE' '=' STRING)? ('CONFIG' '=' STRING)? ('OVERLAP_POLICY' '=' ('NO_OVERLAP' | 'ALLOW_CHILD_OVERLAP' | 'ALLOW_ALL_OVERLAP'))? (ID '=' EXPR (',' ID '=' EXPR)*)? ('USER_TASK_TIMEOUT_MS' '=' NUMBER)? ('SUSPEND_TASK_AFTER_NUM_FAILURES' '=' NUMBER)? ('ERROR_INTEGRATION' '=' ID)? ('SUCCESS_INTEGRATION' '=' ID)? ('LOG_LEVEL' '=' STRING)? ('COMMENT' '=' STRING)? ('FINALIZE' '=' STRING)? ('TASK_AUTO_RETRY_ATTEMPTS' '=' NUMBER)? ('USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS' '=' NUMBER)? ('TARGET_COMPLETION_INTERVAL' '=' STRING)? ('SERVERLESS_TASK_MIN_STATEMENT_SIZE' '=' STRING)? ('SERVERLESS_TASK_MAX_STATEMENT_SIZE' '=' STRING)? ('AFTER' QNAME (',' QNAME)*)? ('EXECUTE' 'AS' 'USER' ID)? ('WHEN' EXPR)? 'AS' QUERY
    | 'CREATE' 'OR' 'ALTER' 'TASK' QNAME ('WAREHOUSE' '=' ID | 'USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE' '=' STRING)? ('SCHEDULE' '=' STRING)? ('CONFIG' '=' STRING)? ('OVERLAP_POLICY' '=' ('NO_OVERLAP' | 'ALLOW_CHILD_OVERLAP' | 'ALLOW_ALL_OVERLAP'))? ('USER_TASK_TIMEOUT_MS' '=' NUMBER)? (ID '=' EXPR (',' ID '=' EXPR)*)? ('SUSPEND_TASK_AFTER_NUM_FAILURES' '=' NUMBER)? ('ERROR_INTEGRATION' '=' ID)? ('SUCCESS_INTEGRATION' '=' ID)? ('COMMENT' '=' STRING)? ('FINALIZE' '=' STRING)? ('TASK_AUTO_RETRY_ATTEMPTS' '=' NUMBER)? ('AFTER' QNAME (',' QNAME)*)? ('EXECUTE' 'AS' 'USER' ID)? ('WHEN' EXPR)? 'AS' QUERY
    | 'CREATE' ('OR' 'REPLACE')? 'TASK' QNAME 'CLONE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-task-graph
createTaskGraph
    : 'CREATE' ('OR' 'REPLACE')? 'TASK' 'GRAPH' QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('SCHEDULE' '=' STRING)? ('ROOTS' '=' '(' QNAME ')')? ('FINALIZER' '=' QNAME)? ('OVERLAP_POLICY' '=' ('NO_OVERLAP' | 'ALLOW_CHILD_OVERLAP' | 'ALLOW_ALL_OVERLAP'))? ('ERROR_INTEGRATION' '=' ID)? ('SUCCESS_INTEGRATION' '=' ID)? ('CONFIG' '=' STRING)? ('TASK_AUTO_RETRY_ATTEMPTS' '=' NUMBER)? ('SUSPEND_TASK_AFTER_NUM_FAILURES' '=' NUMBER)? ('COMMENT' '=' STRING)? ('TASK_DEFAULTS' '(' ID '=' ID (',' ID '=' ID)* ')')? ('WHEN' EXPR)?
    | 'CREATE' 'TASK' 'GRAPH' ('IF' 'NOT' 'EXISTS')? QNAME ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? ('SCHEDULE' '=' STRING)? ('ROOTS' '=' '(' QNAME ')')? ('FINALIZER' '=' QNAME)? ('OVERLAP_POLICY' '=' ('NO_OVERLAP' | 'ALLOW_CHILD_OVERLAP' | 'ALLOW_ALL_OVERLAP'))? ('ERROR_INTEGRATION' '=' ID)? ('SUCCESS_INTEGRATION' '=' ID)? ('CONFIG' '=' STRING)? ('TASK_AUTO_RETRY_ATTEMPTS' '=' NUMBER)? ('SUSPEND_TASK_AFTER_NUM_FAILURES' '=' NUMBER)? ('COMMENT' '=' STRING)? ('TASK_DEFAULTS' '(' ID '=' ID (',' ID '=' ID)* ')')? ('WHEN' EXPR)?
    | 'CREATE' 'OR' 'ALTER' 'TASK' 'GRAPH' QNAME ('SCHEDULE' '=' STRING)? ('ROOTS' '=' '(' QNAME ')')? ('FINALIZER' '=' QNAME)? ('OVERLAP_POLICY' '=' ('NO_OVERLAP' | 'ALLOW_CHILD_OVERLAP' | 'ALLOW_ALL_OVERLAP'))? ('ERROR_INTEGRATION' '=' ID)? ('SUCCESS_INTEGRATION' '=' ID)? ('CONFIG' '=' STRING)? ('TASK_AUTO_RETRY_ATTEMPTS' '=' NUMBER)? ('SUSPEND_TASK_AFTER_NUM_FAILURES' '=' NUMBER)? ('COMMENT' '=' STRING)? ('TASK_DEFAULTS' '(' ID '=' ID (',' ID '=' ID)* ')')? ('WHEN' EXPR)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-type
createType
    : 'CREATE' ('OR' 'REPLACE')? 'TYPE' QNAME 'AS' DATA_TYPE ('COMMENT' '=' STRING)?
    | 'CREATE' 'TYPE' ('IF' 'NOT' 'EXISTS')? QNAME 'AS' DATA_TYPE ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-user
// fix: list_of_properties is the workload identity properties the parameters list
// fix: the interfaces are ALL or a list of the ones the parameters name
createUser
    : 'CREATE' ('OR' 'REPLACE')? 'USER' ID createUser_objectProperties? createUser_objectParams? createUser_sessionParams? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    | 'CREATE' 'USER' ('IF' 'NOT' 'EXISTS')? ID createUser_objectProperties? createUser_objectParams? createUser_sessionParams? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?
    ;
createUser_objectProperties
    : ('PASSWORD' '=' STRING | 'LOGIN_NAME' '=' STRING | 'DISPLAY_NAME' '=' STRING | 'FIRST_NAME' '=' STRING | 'MIDDLE_NAME' '=' STRING | 'LAST_NAME' '=' STRING | 'EMAIL' '=' STRING | 'MUST_CHANGE_PASSWORD' '=' ('TRUE' | 'FALSE') | 'DISABLED' '=' ('TRUE' | 'FALSE') | 'ALLOWED_INTERFACES' '=' '(' ('\'ALL\'' | ('\'SNOWFLAKE_INTELLIGENCE\'' | '\'STREAMLIT\'') (',' ('\'SNOWFLAKE_INTELLIGENCE\'' | '\'STREAMLIT\''))*) ')' | 'DAYS_TO_EXPIRY' '=' NUMBER | 'MINS_TO_UNLOCK' '=' NUMBER | 'DEFAULT_WAREHOUSE' '=' STRING | 'DEFAULT_NAMESPACE' '=' STRING | 'DEFAULT_ROLE' '=' STRING | 'DEFAULT_SECONDARY_ROLES' '=' ('(' '\'ALL\'' ')' | '(' ')') | 'MINS_TO_BYPASS_MFA' '=' NUMBER | 'RSA_PUBLIC_KEY' '=' STRING | 'RSA_PUBLIC_KEY_FP' '=' STRING | 'RSA_PUBLIC_KEY_2' '=' STRING | 'RSA_PUBLIC_KEY_2_FP' '=' STRING | 'TYPE' '=' ('PERSON' | 'SERVICE' | 'SERVICE_AGENT' | 'LEGACY_SERVICE') | 'WORKLOAD_IDENTITY' '=' '(' 'TYPE' '=' ('AWS' | 'AZURE' | 'GCP' | 'OIDC') ('ARN' '=' STRING)? ('ISSUER' '=' STRING)? ('SUBJECT' '=' STRING)? ('OIDC_AUDIENCE_LIST' '=' '(' STRING (',' STRING)* ')')? ')' | 'COMMENT' '=' STRING)+
    ;
createUser_objectParams
    : ('ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR' '=' ('TRUE' | 'FALSE') | 'ENABLE_UNREDACTED_SECURE_OBJECT_ERROR' '=' ('TRUE' | 'FALSE') | 'NETWORK_POLICY' '=' STRING)+
    ;
createUser_sessionParams
    : ('ABORT_DETACHED_QUERY' '=' ('TRUE' | 'FALSE') | 'AUTOCOMMIT' '=' ('TRUE' | 'FALSE') | 'BINARY_INPUT_FORMAT' '=' STRING | 'BINARY_OUTPUT_FORMAT' '=' STRING | 'DATE_INPUT_FORMAT' '=' STRING | 'DATE_OUTPUT_FORMAT' '=' STRING | 'DEFAULT_NULL_ORDERING' '=' STRING | 'ENABLE_GET_DDL_USE_DATA_TYPE_ALIAS' '=' ('TRUE' | 'FALSE') | 'ERROR_ON_NONDETERMINISTIC_MERGE' '=' ('TRUE' | 'FALSE') | 'ERROR_ON_NONDETERMINISTIC_UPDATE' '=' ('TRUE' | 'FALSE') | 'JSON_INDENT' '=' NUMBER | 'LOCK_TIMEOUT' '=' NUMBER | 'OPT_OUT_ERROR_LOGGING' '=' ('TRUE' | 'FALSE') | 'QUERY_TAG' '=' STRING | 'ROWS_PER_RESULTSET' '=' NUMBER | 'SIMULATED_DATA_SHARING_CONSUMER' '=' STRING | 'STATEMENT_TIMEOUT_IN_SECONDS' '=' NUMBER | 'STRICT_JSON_OUTPUT' '=' ('TRUE' | 'FALSE') | 'TIMESTAMP_DAY_IS_ALWAYS_24H' '=' ('TRUE' | 'FALSE') | 'TIMESTAMP_INPUT_FORMAT' '=' STRING | 'TIMESTAMP_LTZ_OUTPUT_FORMAT' '=' STRING | 'TIMESTAMP_NTZ_OUTPUT_FORMAT' '=' STRING | 'TIMESTAMP_OUTPUT_FORMAT' '=' STRING | 'TIMESTAMP_TYPE_MAPPING' '=' STRING | 'TIMESTAMP_TZ_OUTPUT_FORMAT' '=' STRING | 'TIMEZONE' '=' STRING | 'TIME_INPUT_FORMAT' '=' STRING | 'TIME_OUTPUT_FORMAT' '=' STRING | 'TRANSACTION_DEFAULT_ISOLATION_LEVEL' '=' STRING | 'TWO_DIGIT_CENTURY_START' '=' NUMBER | 'UNSUPPORTED_DDL_ACTION' '=' STRING | 'USE_CACHED_RESULT' '=' ('TRUE' | 'FALSE') | 'WEEK_OF_YEAR_POLICY' '=' NUMBER | 'WEEK_START' '=' NUMBER)+
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-variable
createVariable
    : 'CREATE' ('OR' 'REPLACE')? 'VARIABLE' QNAME DATA_TYPE ('DEFAULT' EXPR)? ('COMMENT' '=' STRING)?
    | 'CREATE' 'VARIABLE' ('IF' 'NOT' 'EXISTS')? QNAME DATA_TYPE ('DEFAULT' EXPR)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-versioned-schema
// fix: the retention properties have no value; <integer>, as CREATE SCHEMA writes them
// fix: the retention properties have no value; <integer>, as CREATE SCHEMA writes them
createVersionedSchema
    : 'CREATE' 'OR' 'ALTER' 'VERSIONED' 'SCHEMA' QNAME ('WITH' 'MANAGED' 'ACCESS')? ('DATA_RETENTION_TIME_IN_DAYS' '=' NUMBER)? ('MAX_DATA_EXTENSION_TIME_IN_DAYS' '=' NUMBER)? ('DEFAULT_DDL_COLLATION' '=' STRING)? ('COMMENT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-view
// fix: the column list holds each column, with an optional comment and its policies, inside its parentheses, as the page's parameters and examples write them
// fix: CREATE OR ALTER VIEW's column list is CREATE VIEW's: names, each with an optional comment
createView
    : 'CREATE' ('OR' 'REPLACE')? 'SECURE'? (((('LOCAL' | 'GLOBAL'))? 'TEMP' | 'TEMPORARY' | 'VOLATILE'))? 'RECURSIVE'? 'VIEW' QNAME ('(' ID ('COMMENT' STRING)? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? (',' ID ('COMMENT' STRING)? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?)* ')')? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('CHANGE_TRACKING' '=' ('TRUE' | 'FALSE'))? ('COPY' 'GRANTS')? ('COPY' 'TAGS')? ('COPY' '=' '(' ('COMMENTS' | 'TAGS' | 'GRANTS') (',' ('COMMENTS' | 'TAGS' | 'GRANTS'))* ')')? ('COMMENT' '=' STRING)? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'AGGREGATION' 'POLICY' QNAME ('ENTITY' 'KEY' '(' ID (',' ID)* ')')?)? ('WITH'? 'JOIN' 'POLICY' QNAME ('ALLOWED' 'JOIN' 'KEYS' '(' ID (',' ID)* ')')?)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? 'AS' QUERY
    | 'CREATE' 'SECURE'? (((('LOCAL' | 'GLOBAL'))? 'TEMP' | 'TEMPORARY' | 'VOLATILE'))? 'RECURSIVE'? 'VIEW' ('IF' 'NOT' 'EXISTS')? QNAME ('(' ID ('COMMENT' STRING)? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? (',' ID ('COMMENT' STRING)? ('WITH'? 'MASKING' 'POLICY' QNAME ('USING' '(' ID ',' ID (',' ID)* ')')?)? ('WITH'? 'PROJECTION' 'POLICY' QNAME)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')?)* ')')? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('CHANGE_TRACKING' '=' ('TRUE' | 'FALSE'))? ('COPY' 'GRANTS')? ('COPY' 'TAGS')? ('COPY' '=' '(' ('COMMENTS' | 'TAGS' | 'GRANTS') (',' ('COMMENTS' | 'TAGS' | 'GRANTS'))* ')')? ('COMMENT' '=' STRING)? ('WITH'? 'ROW' 'ACCESS' 'POLICY' QNAME 'ON' '(' ID (',' ID)* ')')? ('WITH'? 'AGGREGATION' 'POLICY' QNAME ('ENTITY' 'KEY' '(' ID (',' ID)* ')')?)? ('WITH'? 'JOIN' 'POLICY' QNAME ('ALLOWED' 'JOIN' 'KEYS' '(' ID (',' ID)* ')')?)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? ('WITH' 'CONTACT' '(' ID '=' ID (',' ID '=' ID)* ')')? 'AS' QUERY
    | 'CREATE' 'OR' 'ALTER' 'SECURE'? (((('LOCAL' | 'GLOBAL'))? 'TEMP' | 'TEMPORARY' | 'VOLATILE'))? 'RECURSIVE'? 'VIEW' QNAME ('(' ID ('COMMENT' STRING)? (',' ID ('COMMENT' STRING)?)* ')')? ('CHANGE_TRACKING' '=' ('TRUE' | 'FALSE'))? ('COMMENT' '=' STRING)? 'AS' QUERY
    ;

// https://docs.snowflake.com/en/sql-reference/sql/create-warehouse
// fix: the sizes include the synonyms the parameters list
createWarehouse
    : 'CREATE' ('OR' 'REPLACE')? 'WAREHOUSE' ID ('WITH'? createWarehouse_objectProperties)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? createWarehouse_objectParams?
    | 'CREATE' 'WAREHOUSE' ('IF' 'NOT' 'EXISTS')? ID ('WITH'? createWarehouse_objectProperties)? ('WITH'? 'TAG' '(' QNAME '=' STRING (',' QNAME '=' STRING)* ')')? createWarehouse_objectParams?
    | 'CREATE' 'OR' 'ALTER' 'WAREHOUSE' ID ('WITH'? createWarehouse_objectProperties)? createWarehouse_objectParams?
    ;
createWarehouse_objectProperties
    : ('WAREHOUSE_TYPE' '=' ('STANDARD' | '\'SNOWPARK-OPTIMIZED\'' | 'ADAPTIVE') | 'WAREHOUSE_SIZE' '=' ('XSMALL' | '\'X-SMALL\'' | 'SMALL' | 'MEDIUM' | 'LARGE' | 'XLARGE' | '\'X-LARGE\'' | 'XXLARGE' | 'X2LARGE' | '\'2X-LARGE\'' | 'XXXLARGE' | 'X3LARGE' | '\'3X-LARGE\'' | 'X4LARGE' | '\'4X-LARGE\'' | 'X5LARGE' | '\'5X-LARGE\'' | 'X6LARGE' | '\'6X-LARGE\'') | 'GENERATION' '=' ('\'1\'' | '\'2\'') | 'RESOURCE_CONSTRAINT' '=' ('STANDARD_GEN_1' | 'STANDARD_GEN_2' | 'MEMORY_1X' | 'MEMORY_1X_x86' | 'MEMORY_16X' | 'MEMORY_16X_x86' | 'MEMORY_64X' | 'MEMORY_64X_x86') | 'MAX_CLUSTER_COUNT' '=' NUMBER | 'MIN_CLUSTER_COUNT' '=' NUMBER | 'SCALING_POLICY' '=' ('STANDARD' | 'ECONOMY') | 'AUTO_SUSPEND' '=' (NUMBER | 'NULL') | 'AUTO_RESUME' '=' ('TRUE' | 'FALSE') | 'INITIALLY_SUSPENDED' '=' ('TRUE' | 'FALSE') | 'RESOURCE_MONITOR' '=' ID | 'COMMENT' '=' STRING | 'ENABLE_QUERY_ACCELERATION' '=' ('TRUE' | 'FALSE') | 'QUERY_ACCELERATION_MAX_SCALE_FACTOR' '=' NUMBER)+
    ;
createWarehouse_objectParams
    : ('MAX_CONCURRENCY_LEVEL' '=' NUMBER | 'STATEMENT_QUEUED_TIMEOUT_IN_SECONDS' '=' NUMBER | 'STATEMENT_TIMEOUT_IN_SECONDS' '=' NUMBER)+
    ;

// https://docs.snowflake.com/en/sql-reference/sql/delete
// fix: USING takes any number of tables or parenthesized subqueries, as the parameters say
delete
    : 'DELETE' 'FROM' QNAME ('USING' (QNAME | '(' QUERY ')') (',' (QNAME | '(' QUERY ')'))*)? ('WHERE' EXPR)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc
// no syntax

// https://docs.snowflake.com/en/sql-reference/sql/desc-agent
descAgent
    : ('DESC' | 'DESCRIBE') ('AS' 'RESOURCE')? 'AGENT' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-aggregation-policy
descAggregationPolicy
    : ('DESC' | 'DESCRIBE') 'AGGREGATION' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-alert
descAlert
    : ('DESC' | 'DESCRIBE') 'ALERT' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-application
descApplication
    : ('DESC' | 'DESCRIBE') 'APPLICATION' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-application-package
descApplicationPackage
    : ('DESC' | 'DESCRIBE') 'APPLICATION' 'PACKAGE' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-application-service
descApplicationService
    : ('DESC' | 'DESCRIBE') 'APPLICATION' 'SERVICE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-artifact-repository
descArtifactRepository
    : ('DESC' | 'DESCRIBE') 'ARTIFACT' 'REPOSITORY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-authentication-policy
descAuthenticationPolicy
    : ('DESC' | 'DESCRIBE') 'AUTHENTICATION' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-available-listing
descAvailableListing
    : ('DESC' | 'DESCRIBE') 'AVAILABLE' 'LISTING' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-available-organization-profile
descAvailableOrganizationProfile
    : ('DESC' | 'DESCRIBE') 'AVAILABLE' 'ORGANIZATION' 'PROFILE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-backup-policy
descBackupPolicy
    : ('DESC' | 'DESCRIBE') 'BACKUP' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-backup-set
descBackupSet
    : ('DESC' | 'DESCRIBE') 'BACKUP' 'SET' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-catalog-integration
descCatalogIntegration
    : ('DESC' | 'DESCRIBE') 'CATALOG' 'INTEGRATION' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-committed-task
descCommittedTask
    : ('DESC' | 'DESCRIBE') 'COMMITTED' 'TASK' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-compute-pool
descComputePool
    : ('DESC' | 'DESCRIBE') 'COMPUTE' 'POOL' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-configuration
descConfiguration
    : ('DESC' | 'DESCRIBE') 'CONFIGURATION' ID ('IN' 'APPLICATION' ID)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-cortex-search
descCortexSearch
    : ('DESC' | 'DESCRIBE') 'CORTEX' 'SEARCH' 'SERVICE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-data-movement-policy
descDataMovementPolicy
    : ('DESC' | 'DESCRIBE') 'DATA' 'MOVEMENT' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-data-movement-rule
descDataMovementRule
    : ('DESC' | 'DESCRIBE') 'DATA' 'MOVEMENT' 'RULE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-database
descDatabase
    : ('DESC' | 'DESCRIBE') 'DATABASE' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-dbt-project
descDbtProject
    : ('DESC' | 'DESCRIBE') 'DBT' 'PROJECT' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-dcm-project
descDcmProject
    : ('DESCRIBE' | 'DESC') 'DCM' 'PROJECT' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-dynamic-table
descDynamicTable
    : ('DESC' | 'DESCRIBE') 'DYNAMIC' 'TABLE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-event-table
descEventTable
    : ('DESC' | 'DESCRIBE') 'EVENT' 'TABLE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-external-agent
descExternalAgent
    : ('DESC' | 'DESCRIBE') 'EXTERNAL' 'AGENT' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-external-consumer
descExternalConsumer
    : ('DESC' | 'DESCRIBE') 'EXTERNAL' 'CONSUMER' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-external-table
descExternalTable
    : ('DESC' | 'DESCRIBE') 'EXTERNAL'? 'TABLE' QNAME ('TYPE' '=' ('COLUMNS' | 'STAGE'))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-external-volume
descExternalVolume
    : ('DESC' | 'DESCRIBE') 'EXTERNAL' 'VOLUME' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-feature-policy
descFeaturePolicy
    : ('DESC' | 'DESCRIBE') 'FEATURE' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-file-format
descFileFormat
    : ('DESC' | 'DESCRIBE') 'FILE' 'FORMAT' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-function
descFunction
    : ('DESC' | 'DESCRIBE') 'FUNCTION' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-function-dmf
descFunctionDmf
    : ('DESC' | 'DESCRIBE') 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' 'TABLE' '(' DATA_TYPE (',' DATA_TYPE)* ')' (',' 'TABLE' '(' DATA_TYPE (',' DATA_TYPE)* ')')? ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-function-spcs
descFunctionSpcs
    : ('DESC' | 'DESCRIBE') 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-gateway
descGateway
    : ('DESC' | 'DESCRIBE') 'GATEWAY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-git-repository
descGitRepository
    : ('DESC' | 'DESCRIBE') 'GIT' 'REPOSITORY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-iceberg-table
descIcebergTable
    : ('DESC' | 'DESCRIBE') 'ICEBERG'? 'TABLE' QNAME ('TYPE' '=' ('COLUMNS' | 'STAGE'))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-ingress-policy
descIngressPolicy
    : ('DESC' | 'DESCRIBE') 'INGRESS' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-integration
descIntegration
    : ('DESC' | 'DESCRIBE') (('API' | 'CATALOG' | 'EXTERNAL' 'ACCESS' | 'NOTIFICATION' | 'SECURITY' | 'STORAGE'))? 'INTEGRATION' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-join-policy
descJoinPolicy
    : ('DESCRIBE' | 'DESC') 'JOIN' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-listing
descListing
    : ('DESC' | 'DESCRIBE') 'LISTING' ID ('REVISION' '=' ('DRAFT' | 'PUBLISHED'))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-maintenance-policy
descMaintenancePolicy
    : 'DESCRIBE' 'MAINTENANCE' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-masking-policy
descMaskingPolicy
    : ('DESC' | 'DESCRIBE') 'MASKING' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-materialized-view
descMaterializedView
    : ('DESC' | 'DESCRIBE') 'MATERIALIZED' 'VIEW' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-mcp-server
descMcpServer
    : ('DESC' | 'DESCRIBE') 'MCP' 'SERVER' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-model-monitor
descModelMonitor
    : ('DESCRIBE' | 'DESC') 'MODEL' 'MONITOR' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-multi-party-approval-policy
descMultiPartyApprovalPolicy
    : ('DESCRIBE' | 'DESC') 'MULTI' 'PARTY' 'APPROVAL' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-network-policy
descNetworkPolicy
    : ('DESC' | 'DESCRIBE') 'NETWORK' 'POLICY' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-network-rule
descNetworkRule
    : ('DESC' | 'DESCRIBE') 'NETWORK' 'RULE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-notebook
descNotebook
    : ('DESC' | 'DESCRIBE') 'NOTEBOOK' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-notification-integration
descNotificationIntegration
    : ('DESC' | 'DESCRIBE') 'NOTIFICATION' 'INTEGRATION' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-oflow-data-plane-integration
descOflowDataPlaneIntegration
    : ('DESC' | 'DESCRIBE') 'OPENFLOW' 'DATA' 'PLANE' 'INTEGRATION' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-online-feature-table
descOnlineFeatureTable
    : ('DESC' | 'DESCRIBE') 'ONLINE' 'FEATURE' 'TABLE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-openflow-connector
descOpenflowConnector
    : 'DESCRIBE' 'OPENFLOW' 'CONNECTOR' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-openflow-deployment
descOpenflowDeployment
    : 'DESCRIBE' 'OPENFLOW' 'DEPLOYMENT' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-openflow-runtime
descOpenflowRuntime
    : 'DESCRIBE' 'OPENFLOW' 'RUNTIME' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-organization-profile
descOrganizationProfile
    : ('DESC' | 'DESCRIBE') 'ORGANIZATION' 'PROFILE' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-packages-policy
descPackagesPolicy
    : ('DESC' | 'DESCRIBE') 'PACKAGES' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-password-policy
descPasswordPolicy
    : ('DESC' | 'DESCRIBE') 'PASSWORD' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-pipe
descPipe
    : ('DESC' | 'DESCRIBE') 'PIPE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-postgres-instance
descPostgresInstance
    : ('DESC' | 'DESCRIBE') 'POSTGRES' 'INSTANCE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-privacy-policy
descPrivacyPolicy
    : ('DESC' | 'DESCRIBE') 'PRIVACY' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-procedure
descProcedure
    : ('DESC' | 'DESCRIBE') 'PROCEDURE' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-projection-policy
descProjectionPolicy
    : ('DESC' | 'DESCRIBE') 'PROJECTION' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-result
descResult
    : ('DESC' | 'DESCRIBE') 'RESULT' (STRING | 'LAST_QUERY_ID' '(' ')')
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-restricted-session-scope
descRestrictedSessionScope
    : ('DESCRIBE' | 'DESC') ('RESTRICTED' 'SESSION' 'SCOPE' | 'RSS') QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-row-access-policy
descRowAccessPolicy
    : ('DESC' | 'DESCRIBE') 'ROW' 'ACCESS' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-schema
descSchema
    : ('DESC' | 'DESCRIBE') 'SCHEMA' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-search-optimization
descSearchOptimization
    : ('DESC' | 'DESCRIBE') 'SEARCH' 'OPTIMIZATION' 'ON' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-secret
descSecret
    : ('DESC' | 'DESCRIBE') 'SECRET' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-semantic-view
descSemanticView
    : ('DESCRIBE' | 'DESC') 'SEMANTIC' 'VIEW' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-sequence
descSequence
    : ('DESC' | 'DESCRIBE') 'SEQUENCE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-service
descService
    : ('DESC' | 'DESCRIBE') 'SERVICE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-session-policy
descSessionPolicy
    : ('DESCRIBE' | 'DESC') 'SESSION' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-share
descShare
    : ('DESC' | 'DESCRIBE') 'SHARE' ID
    | ('DESC' | 'DESCRIBE') 'SHARE' ID '.' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-snapshot
descSnapshot
    : ('DESC' | 'DESCRIBE') 'SNAPSHOT' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-snapshot-policy
descSnapshotPolicy
    : ('DESC' | 'DESCRIBE') 'SNAPSHOT' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-snapshot-set
descSnapshotSet
    : ('DESC' | 'DESCRIBE') 'SNAPSHOT' 'SET' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-specification
descSpecification
    : ('DESCRIBE' | 'DESC') 'SPECIFICATION' QNAME ('IN' 'APPLICATION' ID)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-stage
descStage
    : ('DESC' | 'DESCRIBE') 'STAGE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-storage-lifecycle-policy
descStorageLifecyclePolicy
    : ('DESC' | 'DESCRIBE') 'STORAGE' 'LIFECYCLE' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-stream
descStream
    : ('DESC' | 'DESCRIBE') 'STREAM' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-streamlit
descStreamlit
    : ('DESC' | 'DESCRIBE') 'STREAMLIT' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-table
descTable
    : ('DESCRIBE' | 'DESC') 'TABLE' QNAME ('TYPE' '=' ('COLUMNS' | 'STAGE'))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-task
descTask
    : ('DESC' | 'DESCRIBE') 'TASK' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-task-graph
descTaskGraph
    : ('DESC' | 'DESCRIBE') 'TASK' 'GRAPH' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-transaction
descTransaction
    : ('DESC' | 'DESCRIBE') 'TRANSACTION' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-type
descType
    : ('DESC' | 'DESCRIBE') 'TYPE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-user
descUser
    : ('DESC' | 'DESCRIBE') 'USER' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-variable
descVariable
    : ('DESC' | 'DESCRIBE') 'VARIABLE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-view
descView
    : ('DESC' | 'DESCRIBE') 'VIEW' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/desc-warehouse
descWarehouse
    : ('DESC' | 'DESCRIBE') 'WAREHOUSE' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop
// no syntax

// https://docs.snowflake.com/en/sql-reference/sql/drop-account
dropAccount
    : 'DROP' 'ACCOUNT' ('IF' 'EXISTS')? ID 'GRACE_PERIOD_IN_DAYS' '=' NUMBER
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-agent
dropAgent
    : 'DROP' 'AGENT' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-aggregation-policy
dropAggregationPolicy
    : 'DROP' 'AGGREGATION' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-alert
dropAlert
    : 'DROP' 'ALERT' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-application
dropApplication
    : 'DROP' 'APPLICATION' ('IF' 'EXISTS')? ID 'CASCADE'?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-application-package
dropApplicationPackage
    : 'DROP' 'APPLICATION' 'PACKAGE' ('IF' 'EXISTS')? ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-application-role
dropApplicationRole
    : 'DROP' 'APPLICATION' 'ROLE' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-application-service
dropApplicationService
    : 'DROP' 'APPLICATION' 'SERVICE' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-artifact-repository
dropArtifactRepository
    : 'DROP' 'ARTIFACT' 'REPOSITORY' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-authentication-policy
dropAuthenticationPolicy
    : 'DROP' 'AUTHENTICATION' 'POLICY' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-backup-policy
dropBackupPolicy
    : 'DROP' 'BACKUP' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-backup-set
dropBackupSet
    : 'DROP' 'BACKUP' 'SET' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-catalog-integration
dropCatalogIntegration
    : 'DROP' 'CATALOG' 'INTEGRATION' ('IF' 'EXISTS')? ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-compute-pool
dropComputePool
    : 'DROP' 'COMPUTE' 'POOL' ('IF' 'EXISTS')? ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-connection
dropConnection
    : 'DROP' 'CONNECTION' ('IF' 'EXISTS')? ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-contact
dropContact
    : 'DROP' 'CONTACT' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-cortex-search
dropCortexSearch
    : 'DROP' 'CORTEX' 'SEARCH' 'SERVICE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-data-movement-policy
dropDataMovementPolicy
    : 'DROP' 'DATA' 'MOVEMENT' 'POLICY' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-data-movement-rule
dropDataMovementRule
    : 'DROP' 'DATA' 'MOVEMENT' 'RULE' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-database
dropDatabase
    : 'DROP' 'DATABASE' ('IF' 'EXISTS')? ID ('CASCADE' | 'RESTRICT')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-database-role
dropDatabaseRole
    : 'DROP' 'DATABASE' 'ROLE' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-dbt-project
dropDbtProject
    : 'DROP' 'DBT' 'PROJECT' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-dcm-project
dropDcmProject
    : 'DROP' 'DCM' 'PROJECT' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-dynamic-table
dropDynamicTable
    : 'DROP' 'DYNAMIC' 'TABLE' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-event-routing-table
dropEventRoutingTable
    : 'DROP' 'EVENT' 'ROUTING' 'TABLE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-experiment
dropExperiment
    : 'DROP' 'EXPERIMENT' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-external-agent
dropExternalAgent
    : 'DROP' 'EXTERNAL' 'AGENT' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-external-consumer
dropExternalConsumer
    : 'DROP' 'EXTERNAL' 'CONSUMER' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-external-table
dropExternalTable
    : 'DROP' 'EXTERNAL' 'TABLE' ('IF' 'EXISTS')? QNAME ('CASCADE' | 'RESTRICT')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-external-volume
dropExternalVolume
    : 'DROP' 'EXTERNAL' 'VOLUME' ('IF' 'EXISTS')? ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/undrop-external-volume
undropExternalVolume
    : 'UNDROP' 'EXTERNAL' 'VOLUME' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-failover-group
dropFailoverGroup
    : 'DROP' 'FAILOVER' 'GROUP' ('IF' 'EXISTS')? ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-feature-policy
dropFeaturePolicy
    : 'DROP' 'FEATURE' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-file-format
dropFileFormat
    : 'DROP' 'FILE' 'FORMAT' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-function
dropFunction
    : 'DROP' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-function-dmf
dropFunctionDmf
    : 'DROP' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' 'TABLE' '(' DATA_TYPE (',' DATA_TYPE)* ')' (',' 'TABLE' '(' DATA_TYPE (',' DATA_TYPE)* ')')? ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-function-spcs
dropFunctionSpcs
    : 'DROP' 'FUNCTION' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-gateway
dropGateway
    : 'DROP' 'GATEWAY' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-git-repository
dropGitRepository
    : 'DROP' 'GIT' 'REPOSITORY' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-iceberg-table
dropIcebergTable
    : 'DROP' 'ICEBERG'? 'TABLE' ('IF' 'EXISTS')? QNAME ('CASCADE' | 'RESTRICT')? ('PURGE' | 'NOPURGE')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/undrop-iceberg-table
undropIcebergTable
    : 'UNDROP' 'ICEBERG' 'TABLE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-image-repository
dropImageRepository
    : 'DROP' 'IMAGE' 'REPOSITORY' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-index
dropIndex
    : 'DROP' 'INDEX' ('IF' 'EXISTS')? QNAME '.' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-ingress-policy
dropIngressPolicy
    : 'DROP' 'INGRESS' 'POLICY' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-integration
dropIntegration
    : 'DROP' (('API' | 'CATALOG' | 'EXTERNAL' 'ACCESS' | 'NOTIFICATION' | 'SECURITY' | 'STORAGE'))? 'INTEGRATION' ('IF' 'EXISTS')? ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-join-policy
dropJoinPolicy
    : 'DROP' 'JOIN' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-listing
dropListing
    : 'DROP' 'LISTING' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-maintenance-policy
dropMaintenancePolicy
    : 'DROP' 'MAINTENANCE' 'POLICY' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-managed-account
dropManagedAccount
    : 'DROP' 'MANAGED' 'ACCOUNT' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-masking-policy
dropMaskingPolicy
    : 'DROP' 'MASKING' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-materialized-view
dropMaterializedView
    : 'DROP' 'MATERIALIZED' 'VIEW' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-mcp-server
dropMcpServer
    : 'DROP' 'MCP' 'SERVER' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-model
dropModel
    : 'DROP' 'MODEL' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-model-monitor
dropModelMonitor
    : 'DROP' 'MODEL' 'MONITOR' ('IF' 'EXISTS')? ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-multi-party-approval-policy
dropMultiPartyApprovalPolicy
    : 'DROP' 'MULTI' 'PARTY' 'APPROVAL' 'POLICY' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-network-policy
dropNetworkPolicy
    : 'DROP' 'NETWORK' 'POLICY' ('IF' 'EXISTS')? ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-network-rule
dropNetworkRule
    : 'DROP' 'NETWORK' 'RULE' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-notebook
dropNotebook
    : 'DROP' 'NOTEBOOK' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/undrop-notebook
undropNotebook
    : 'UNDROP' 'NOTEBOOK' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-online-feature-table
dropOnlineFeatureTable
    : 'DROP' 'ONLINE' 'FEATURE' 'TABLE' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-openflow-connector
dropOpenflowConnector
    : 'DROP' 'OPENFLOW' 'CONNECTOR' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-openflow-deployment
dropOpenflowDeployment
    : 'DROP' 'OPENFLOW' 'DEPLOYMENT' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-openflow-runtime
dropOpenflowRuntime
    : 'DROP' 'OPENFLOW' 'RUNTIME' ('IF' 'EXISTS')? QNAME 'CASCADE'?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-organization-profile
dropOrganizationProfile
    : 'DROP' 'ORGANIZATION' 'PROFILE' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-organization-user
dropOrganizationUser
    : 'DROP' 'ORGANIZATION' 'USER' ('IF' 'EXISTS')? ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-organization-user-group
dropOrganizationUserGroup
    : 'DROP' 'ORGANIZATION' 'USER' 'GROUP' ('IF' 'EXISTS')? ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-packages-policy
dropPackagesPolicy
    : 'DROP' 'PACKAGES' 'POLICY' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-password-policy
dropPasswordPolicy
    : 'DROP' 'PASSWORD' 'POLICY' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-pipe
dropPipe
    : 'DROP' 'PIPE' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-postgres-instance
dropPostgresInstance
    : 'DROP' 'POSTGRES' 'INSTANCE' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-privacy-policy
dropPrivacyPolicy
    : 'DROP' 'PRIVACY' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-procedure
dropProcedure
    : 'DROP' 'PROCEDURE' ('IF' 'EXISTS')? QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-projection-policy
dropProjectionPolicy
    : 'DROP' 'PROJECTION' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-replication-group
dropReplicationGroup
    : 'DROP' 'REPLICATION' 'GROUP' ('IF' 'EXISTS')? ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-resource-monitor
dropResourceMonitor
    : 'DROP' 'RESOURCE' 'MONITOR' ('IF' 'EXISTS')? ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-restricted-session-scope
dropRestrictedSessionScope
    : 'DROP' ('RESTRICTED' 'SESSION' 'SCOPE' | 'RSS') ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-role
dropRole
    : 'DROP' 'ROLE' ('IF' 'EXISTS')? ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-row-access-policy
dropRowAccessPolicy
    : 'DROP' 'ROW' 'ACCESS' 'POLICY' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-schema
dropSchema
    : 'DROP' 'SCHEMA' ('IF' 'EXISTS')? QNAME ('CASCADE' | 'RESTRICT')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-secret
dropSecret
    : 'DROP' 'SECRET' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-semantic-view
dropSemanticView
    : 'DROP' 'SEMANTIC' 'VIEW' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-sequence
dropSequence
    : 'DROP' 'SEQUENCE' ('IF' 'EXISTS')? QNAME ('CASCADE' | 'RESTRICT')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-service
dropService
    : 'DROP' 'SERVICE' ('IF' 'EXISTS')? QNAME 'FORCE'?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-session-policy
dropSessionPolicy
    : 'DROP' 'SESSION' 'POLICY' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-share
dropShare
    : 'DROP' 'SHARE' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-snapshot
dropSnapshot
    : 'DROP' 'SNAPSHOT' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-snapshot-policy
dropSnapshotPolicy
    : 'DROP' 'SNAPSHOT' 'POLICY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-snapshot-set
dropSnapshotSet
    : 'DROP' 'SNAPSHOT' 'SET' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-stage
dropStage
    : 'DROP' 'STAGE' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-storage-lifecycle-policy
dropStorageLifecyclePolicy
    : 'DROP' 'STORAGE' 'LIFECYCLE' 'POLICY' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-stream
dropStream
    : 'DROP' 'STREAM' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-streamlit
dropStreamlit
    : 'DROP' 'STREAMLIT' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-table
dropTable
    : 'DROP' 'TABLE' ('IF' 'EXISTS')? QNAME ('CASCADE' | 'RESTRICT')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/undrop-table
undropTable
    : 'UNDROP' 'TABLE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-tag
dropTag
    : 'DROP' 'TAG' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-task
dropTask
    : 'DROP' 'TASK' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-task-graph
dropTaskGraph
    : 'DROP' 'TASK' 'GRAPH' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-type
dropType
    : 'DROP' 'TYPE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-user
dropUser
    : 'DROP' 'USER' ('IF' 'EXISTS')? ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-variable
dropVariable
    : 'DROP' 'VARIABLE' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-view
dropView
    : 'DROP' 'VIEW' ('IF' 'EXISTS')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/drop-warehouse
dropWarehouse
    : 'DROP' 'WAREHOUSE' ('IF' 'EXISTS')? ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/execute-alert
executeAlert
    : 'EXECUTE' 'ALERT' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/execute-dbt-project
// fix: the version is text to fill in, as the parameters say (1.9.4 by default)
executeDbtProject
    : 'EXECUTE' 'DBT' 'PROJECT' ('IF' 'EXISTS')? QNAME ('ARGS' '=' STRING)? ('DBT_VERSION' '=' STRING)? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('ENVIRONMENT' '=' STRING)? ('ENV_VARS' '=' '(' STRING '=' STRING (',' STRING '=' STRING)* ')')? ('IMPORTS' '=' '(' (STRING | ID) ('AS' STRING)? (',' (STRING | ID) ('AS' STRING)?)* ')')? ('WRITEBACK' '=' ('TRUE' | 'FALSE'))?
    | 'EXECUTE' 'DBT' 'PROJECT' ('IF' 'EXISTS')? ('FROM' 'WORKSPACE' QNAME)? ('ARGS' '=' STRING)? ('DBT_VERSION' '=' STRING)? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('ENVIRONMENT' '=' STRING)? ('ENV_VARS' '=' '(' STRING '=' STRING (',' STRING '=' STRING)* ')')? ('PROJECT_ROOT' '=' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/execute-dcm-project
executeDcmProject
    : 'EXECUTE' 'DCM' 'PROJECT' QNAME 'PLAN' 'DELTA'? (executeDcmProject_fromSource | executeDcmProject_fromDeployment) ('OUTPUT_PATH' STRING)?
    | 'EXECUTE' 'DCM' 'PROJECT' QNAME 'DEPLOY' ('AS' ID)? (executeDcmProject_fromSource | executeDcmProject_fromDeployment)
    | 'EXECUTE' 'DCM' 'PROJECT' QNAME 'TEST' 'ALL'
    | 'EXECUTE' 'DCM' 'PROJECT' QNAME 'PREVIEW' QNAME 'USING' 'CONFIGURATION' ID ('ENVIRONMENT' '(' ID ')')? 'FROM' STRING ('LIMIT' NUMBER)?
    | 'EXECUTE' 'DCM' 'PROJECT' QNAME 'PURGE' ('AS' ID)?
    ;
executeDcmProject_fromSource
    : ('USING' ('CONFIGURATION' ID)? ('(' EXPR (',' EXPR)* ')')?)? ('ENVIRONMENT' '(' ID ')')? 'FROM' STRING
    ;
executeDcmProject_fromDeployment
    : 'FROM' 'DEPLOYMENT' (ID | 'LAST' | 'LATEST')
    ;

// https://docs.snowflake.com/en/sql-reference/sql/execute-immediate
executeImmediate
    : 'EXECUTE' 'IMMEDIATE' STRING ('USING' '(' ID (',' ID)* ')')?
    | 'EXECUTE' 'IMMEDIATE' ID ('USING' '(' ID (',' ID)* ')')?
    | 'EXECUTE' 'IMMEDIATE' VARIABLE ('USING' '(' ID (',' ID)* ')')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/execute-immediate-from
// fix: a template value is a constant or variable, as the templating notes say, not a name
executeImmediateFrom
    : 'EXECUTE' 'IMMEDIATE' 'FROM' (executeImmediateFrom_absoluteFilePath | executeImmediateFrom_relativeFilePath) ('USING' '(' ID '=>' EXPR (',' ID '=>' EXPR)* ')')? ('DRY_RUN' '=' ('TRUE' | 'FALSE'))?
    ;
executeImmediateFrom_absoluteFilePath
    : STAGE
    ;
executeImmediateFrom_relativeFilePath
    : STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/execute-inference-job-service
executeInferenceJobService
    : 'EXECUTE' 'INFERENCE' 'JOB' 'SERVICE' 'IN' 'COMPUTE' 'POOL' ID 'WITH' 'SPECIFICATION' STRING 'FROM' ('(' QUERY ')' | STAGE) 'MODEL' '=' (ID '.' ID '.')? ID ('VERSION' '=' ID)? ('FUNCTION' '=' STRING)? ('NAME' '=' (ID '.' ID '.')? ID)? ('ASYNC' '=' ('TRUE' | 'FALSE'))? ('REPLICAS' '=' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/execute-job-service
// fix: REPLICAS has its = twice
// fix: read: a template argument's <value> is a constant, a string or JSON value, as the USING parameters say
executeJobService
    : 'EXECUTE' 'JOB' 'SERVICE' 'IN' 'COMPUTE' 'POOL' ID (executeJobService_fromSpecification | executeJobService_fromSpecificationTemplate) ('NAME' '=' (ID '.' ID '.')? ID)? ('ASYNC' '=' ('TRUE' | 'FALSE'))? ('REPLICAS' '=' NUMBER)? ('QUERY_WAREHOUSE' '=' ID)? ('COMMENT' '=' STRING)? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')?
    ;
executeJobService_fromSpecification
    : ('FROM' STAGE 'SPECIFICATION_FILE' '=' STRING | 'FROM' 'SPECIFICATION' STRING)
    ;
executeJobService_fromSpecificationTemplate
    : ('FROM' STAGE 'SPECIFICATION_TEMPLATE_FILE' '=' STRING | 'FROM' 'SPECIFICATION_TEMPLATE' STRING) 'USING' '(' ID '=>' EXPR (',' ID '=>' EXPR)* ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/execute-notebook
executeNotebook
    : 'EXECUTE' 'NOTEBOOK' QNAME '(' (STRING (',' STRING)*)? ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/execute-notebook-project
// fix: notebook.ipynb is an example; the main file is a path to fill in, as the parameters say
executeNotebookProject
    : 'EXECUTE' 'NOTEBOOK' 'PROJECT' ID '.' ID '.' ID 'MAIN_FILE' '=' STRING 'COMPUTE_POOL' '=' STRING 'QUERY_WAREHOUSE' '=' STRING 'RUNTIME' '=' STRING ('ARGUMENTS' '=' STRING)? ('REQUIREMENTS_FILE' '=' STRING)? ('ARTIFACT_REPOSITORIES' '=' '(' QNAME (',' QNAME)* ')')? ('EXTERNAL_ACCESS_INTEGRATIONS' '=' '(' ID (',' ID)* ')')? ('SECRETS' '=' '(' ID '.' ID '.' ID (',' ID '.' ID '.' ID)* ')')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/execute-code-bundle
executeCodeBundle
    : 'EXECUTE' 'CODE' 'BUNDLE' QNAME 'ENTRYPOINT' '=' STRING ('ARGUMENTS' '=' '(' STRING (',' STRING)* ')')? ('WITH' 'SPECIFICATION' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/execute-task
executeTask
    : 'EXECUTE' 'TASK' QNAME ('USING' 'CONFIG' '=' STRING)? 'AWAIT'?
    | 'EXECUTE' 'TASK' QNAME 'RETRY' 'LAST'
    | 'EXECUTE' 'TASK' QNAME 'RETRY' 'GRAPH' 'RUN' 'GROUP' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/execute-task-graph
executeTaskGraph
    : 'EXECUTE' 'TASK' 'GRAPH' QNAME ('USING' 'CONFIG' '=' STRING)?
    | 'EXECUTE' 'TASK' 'GRAPH' QNAME 'RETRY' 'LAST'
    | 'EXECUTE' 'TASK' 'GRAPH' QNAME 'RETRY' 'GRAPH' 'RUN' 'GROUP' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/experiment
// no syntax

// https://docs.snowflake.com/en/sql-reference/sql/explain
explain
    : 'EXPLAIN' ('USING' ('TABULAR' | 'JSON' | 'TEXT'))? QUERY
    ;

// https://docs.snowflake.com/en/sql-reference/sql/get
// fix: PATTERN's string has a quote too many; dropped
get
    : 'GET' get_internalStage FILE ('PARALLEL' '=' NUMBER)? ('PATTERN' '=' STRING)?
    ;
get_internalStage
    : STAGE
    ;

// https://docs.snowflake.com/en/sql-reference/sql/grant-application-role
grantApplicationRole
    : 'GRANT' 'APPLICATION' 'ROLE' QNAME 'TO' 'ROLE' ID ('WITH' 'GRANT' 'OPTION')?
    | 'GRANT' 'APPLICATION' 'ROLE' QNAME 'TO' ('APPLICATION' 'ROLE' ID | 'APPLICATION' ID | 'USER' ID)
    ;

// https://docs.snowflake.com/en/sql-reference/sql/grant-caller
// fix: the privileges are each object type's, listed on Access control privileges rather than here; {object_privileges} leaves the forms unconverted
// not converted (form): undefined fragment "object_privileges": the page does not define it
//   GRANT CALLER {object_privileges} [ , {object_privileges} ... ]
//     ON <object_type> <object_name>
//     TO { ROLE | DATABASE ROLE | APPLICATION } <grantee_name>
// not converted (form): undefined fragment "object_privileges": the page does not define it
//   GRANT INHERITED CALLER {object_privileges} [ , {object_privileges} ... ]
//     ON ALL <object_type_plural>
//     IN { ACCOUNT | DATABASE <db_name> | SCHEMA <schema_name> | APPLICATION <app_name> | APPLICATION PACKAGE <app_pkg_name> }
//     TO { ROLE | DATABASE ROLE | APPLICATION } <grantee_name>
grantCaller
    : 'GRANT' 'ALL' 'CALLER' 'PRIVILEGES' 'ON' OBJECT_TYPE QNAME 'TO' ('ROLE' | 'DATABASE' 'ROLE' | 'APPLICATION') ID
    | 'GRANT' 'ALL' 'INHERITED' 'CALLER' 'PRIVILEGES' 'ON' 'ALL' OBJECT_TYPES 'IN' ('ACCOUNT' | 'DATABASE' ID | 'SCHEMA' QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID) 'TO' ('ROLE' | 'DATABASE' 'ROLE' | 'APPLICATION') ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/grant-database-role
grantDatabaseRole
    : 'GRANT' 'DATABASE' 'ROLE' QNAME 'TO' ('DATABASE' 'ROLE' ID | 'ROLE' ID) ('WITH' 'GRANT' 'OPTION')?
    | 'GRANT' 'DATABASE' 'ROLE' QNAME 'TO' 'USER' ID
    | 'GRANT' 'DATABASE' 'ROLE' QNAME 'TO' 'APPLICATION' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/grant-database-role-share
grantDatabaseRoleShare
    : 'GRANT' 'DATABASE' 'ROLE' QNAME 'TO' 'SHARE' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/grant-ownership
grantOwnership
    : 'GRANT' 'OWNERSHIP' ('ON' (OBJECT_TYPE QNAME | 'ALL' OBJECT_TYPES 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | 'ON' 'FUTURE' OBJECT_TYPES 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) 'TO' ('ROLE' ID | 'DATABASE' 'ROLE' ID) (('REVOKE' | 'COPY') 'CURRENT' 'GRANTS')?
    | 'GRANT' 'OWNERSHIP' 'ON' ID ID 'TO' ('ROLE' ID | 'DATABASE' 'ROLE' ID) (('REVOKE' | 'COPY') 'CURRENT' 'GRANTS')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/grant-privilege
// fix: TO DATABASE ROLE form: the } after the database privileges ends the choice of grants one line early; dropped
// fix: the trailing [ , ... ] repeats the whole choice of schema privileges, as the usage notes allow several privileges on one object type
// fix: closes the { opened before ADD SEARCH OPTIMIZATION
// fix: read: each -- For X section of accountObjectPrivileges and schemaObjectPrivileges is the privileges on X alone, and the ON clause chooses an object type with its own privileges
grantPrivilege
    : 'GRANT' ((grantPrivilege_globalPrivileges | 'ALL' 'PRIVILEGES'?) 'ON' 'ACCOUNT' | ((grantPrivilege_accountObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' 'USER' ID | (grantPrivilege_accountObjectPrivileges_11 | 'ALL' 'PRIVILEGES'?) 'ON' 'RESOURCE' 'MONITOR' ID | (grantPrivilege_accountObjectPrivileges_13 | 'ALL' 'PRIVILEGES'?) 'ON' 'WAREHOUSE' ID | (grantPrivilege_accountObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' 'COMPUTE' 'POOL' ID | (grantPrivilege_accountObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' 'DATABASE' ID | (grantPrivilege_accountObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' 'DATA' 'EXCHANGE' 'LISTING' ID | (grantPrivilege_accountObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' 'INTEGRATION' ID | (grantPrivilege_accountObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' 'CONNECTION' ID | (grantPrivilege_accountObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' 'FAILOVER' 'GROUP' ID | (grantPrivilege_accountObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' 'REPLICATION' 'GROUP' ID | (grantPrivilege_accountObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' 'EXTERNAL' 'VOLUME' ID | (grantPrivilege_accountObjectPrivileges_1 | 'ALL' 'PRIVILEGES'?) 'ON' 'APPLICATION' 'PACKAGE' ID | 'ALL' 'PRIVILEGES'? 'ON' 'APPLICATION' ID) | (grantPrivilege_schemaPrivileges | 'ALL' 'PRIVILEGES'?) 'ON' ('SCHEMA' QNAME | 'ALL' 'SCHEMAS' 'IN' 'DATABASE' ID) | (grantPrivilege_schemaPrivileges | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SCHEMAS' 'IN' 'DATABASE' ID | ((grantPrivilege_schemaObjectPrivileges_1 | 'ALL' 'PRIVILEGES'?) 'ON' ('AGENT' QNAME | 'ALL' 'AGENTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('AGGREGATION' 'POLICY' QNAME | 'ALL' 'AGGREGATION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' ('ALERT' QNAME | 'ALL' 'ALERTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' ('APPLICATION' 'SERVICE' QNAME | 'ALL' 'APPLICATION' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('ARTIFACT' 'REPOSITORY' QNAME | 'ALL' 'ARTIFACT' 'REPOSITORIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('AUTHENTICATION' 'POLICY' QNAME | 'ALL' 'AUTHENTICATION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('CONTACT' QNAME | 'ALL' 'CONTACTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' ('CORTEX' 'EXTENSION' QNAME | 'ALL' 'CORTEX' 'EXTENSIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' ('CORTEX' 'SEARCH' 'SERVICE' QNAME | 'ALL' 'CORTEX' 'SEARCH' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('DATA' 'METRIC' 'FUNCTION' QNAME | 'ALL' 'DATA' 'METRIC' 'FUNCTIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('DATASET' QNAME | 'ALL' 'DATASETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_11 | 'ALL' 'PRIVILEGES'?) 'ON' ('DBT' 'PROJECT' QNAME | 'ALL' 'DBT' 'PROJECTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('DCM' 'PROJECT' QNAME | 'ALL' 'DCM' 'PROJECTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_13 | 'ALL' 'PRIVILEGES'?) 'ON' ('DYNAMIC' 'TABLE' QNAME | 'ALL' 'DYNAMIC' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_15 | 'ALL' 'PRIVILEGES'?) 'ON' ('EVENT' 'TABLE' QNAME | 'ALL' 'EVENT' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_14 | 'ALL' 'PRIVILEGES'?) 'ON' ('EXPERIMENT' QNAME | 'ALL' 'EXPERIMENTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('FILE' 'FORMAT' QNAME | 'ALL' 'FILE' 'FORMATS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('FUNCTION' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' | 'ALL' 'FUNCTIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_16 | 'ALL' 'PRIVILEGES'?) 'ON' ('GATEWAY' QNAME | 'ALL' 'GATEWAYS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_17 | 'ALL' 'PRIVILEGES'?) 'ON' ('GIT' 'REPOSITORY' QNAME | 'ALL' 'GIT' 'REPOSITORIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_18 | 'ALL' 'PRIVILEGES'?) 'ON' ('HYBRID' 'TABLE' QNAME | 'ALL' 'HYBRID' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_20 | 'ALL' 'PRIVILEGES'?) 'ON' ('ICEBERG' 'TABLE' QNAME | 'ALL' 'ICEBERG' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_19 | 'ALL' 'PRIVILEGES'?) 'ON' ('IMAGE' 'REPOSITORY' QNAME | 'ALL' 'IMAGE' 'REPOSITORIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_21 | 'ALL' 'PRIVILEGES'?) 'ON' ('INTERACTIVE' 'TABLE' QNAME | 'ALL' 'INTERACTIVE' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('JOIN' 'POLICY' QNAME | 'ALL' 'JOIN' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' 'MAINTENANCE' 'POLICY' QNAME | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' 'MASKING' 'POLICY' QNAME | (grantPrivilege_schemaObjectPrivileges_22 | 'ALL' 'PRIVILEGES'?) 'ON' ('MATERIALIZED' 'VIEW' QNAME | 'ALL' 'MATERIALIZED' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_23 | 'ALL' 'PRIVILEGES'?) 'ON' ('MCP' 'SERVER' QNAME | 'ALL' 'MCP' 'SERVERS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('MODEL' QNAME | 'ALL' 'MODELS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_24 | 'ALL' 'PRIVILEGES'?) 'ON' ('ONLINE' 'FEATURE' 'TABLE' QNAME | 'ALL' 'ONLINE' 'FEATURE' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('PACKAGES' 'POLICY' QNAME | 'ALL' 'PACKAGES' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('PASSWORD' 'POLICY' QNAME | 'ALL' 'PASSWORD' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_25 | 'ALL' 'PRIVILEGES'?) 'ON' 'PIPE' QNAME | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('PRIVACY' 'POLICY' QNAME | 'ALL' 'PRIVACY' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('PROCEDURE' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' | 'ALL' 'PROCEDURES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('PROJECTION' 'POLICY' QNAME | 'ALL' 'PROJECTION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('ROW' 'ACCESS' 'POLICY' QNAME | 'ALL' 'ROW' 'ACCESS' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_9 | grantPrivilege_schemaObjectPrivileges_27 | 'ALL' 'PRIVILEGES'?) 'ON' ('SECRET' QNAME | 'ALL' 'SECRETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_28 | 'ALL' 'PRIVILEGES'?) 'ON' ('SEMANTIC' 'VIEW' QNAME | 'ALL' 'SEMANTIC' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('SEQUENCE' QNAME | 'ALL' 'SEQUENCES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_29 | 'ALL' 'PRIVILEGES'?) 'ON' ('SERVICE' QNAME | 'ALL' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('SESSION' 'POLICY' QNAME | 'ALL' 'SESSION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('SNAPSHOT' QNAME | 'ALL' 'SNAPSHOTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' ('SNAPSHOT' 'POLICY' QNAME | 'ALL' 'SNAPSHOT' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' ('SNAPSHOT' 'SET' QNAME | 'ALL' 'SNAPSHOT' 'SETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_30 | grantPrivilege_schemaObjectPrivileges_31 | 'ALL' 'PRIVILEGES'?) 'ON' ('STAGE' QNAME | 'ALL' 'STAGES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('STORAGE' 'LIFECYCLE' 'POLICY' QNAME | 'ALL' 'STORAGE' 'LIFECYCLE' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_32 | 'ALL' 'PRIVILEGES'?) 'ON' ('STREAM' QNAME | 'ALL' 'STREAMS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_33 | 'ALL' 'PRIVILEGES'?) 'ON' ('STREAMLIT' QNAME | 'ALL' 'STREAMLITS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_34 | 'ALL' 'PRIVILEGES'?) 'ON' ('TABLE' QNAME | 'ALL' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | grantPrivilege_schemaObjectPrivileges_35 | 'ALL' 'PRIVILEGES'?) 'ON' 'TAG' QNAME | (grantPrivilege_schemaObjectPrivileges_36 | 'ALL' 'PRIVILEGES'?) 'ON' ('TASK' QNAME | 'ALL' 'TASKS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('TYPE' QNAME | 'ALL' 'TYPES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_37 | 'ALL' 'PRIVILEGES'?) 'ON' ('VIEW' QNAME | 'ALL' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_38 | 'ALL' 'PRIVILEGES'?) 'ON' ('WORKSPACE' QNAME | 'ALL' 'WORKSPACES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME))) | ((grantPrivilege_schemaObjectPrivileges_1 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'AGENTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'ALERTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'APPLICATION' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'ARTIFACT' 'REPOSITORIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'AUTHENTICATION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'CONTACTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'CORTEX' 'EXTENSIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'CORTEX' 'SEARCH' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DATA' 'METRIC' 'FUNCTIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DATASETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_11 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DBT' 'PROJECTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DCM' 'PROJECTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_13 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DYNAMIC' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_15 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'EVENT' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_14 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'EXPERIMENTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'FILE' 'FORMATS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'FUNCTIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_16 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'GATEWAYS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_17 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'GIT' 'REPOSITORIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_18 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'HYBRID' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_20 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'ICEBERG' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_21 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'INTERACTIVE' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_22 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'MATERIALIZED' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_23 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'MCP' 'SERVERS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'MODELS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_24 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'ONLINE' 'FEATURE' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PASSWORD' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_25 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PIPES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PRIVACY' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PROCEDURES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_9 | grantPrivilege_schemaObjectPrivileges_27 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SECRETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_28 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SEMANTIC' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SEQUENCES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_29 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SNAPSHOT' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SNAPSHOT' 'SETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_30 | grantPrivilege_schemaObjectPrivileges_31 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'STAGES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_32 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'STREAMS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_33 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'STREAMLITS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_34 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_36 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'TASKS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'TYPES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_37 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_38 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'WORKSPACES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME))) 'TO' 'ROLE'? ID ('WITH' 'GRANT' 'OPTION')?
    | 'GRANT' (('CREATE' 'SCHEMA' | 'MODIFY' | 'MONITOR' | 'USAGE') (',' ('CREATE' 'SCHEMA' | 'MODIFY' | 'MONITOR' | 'USAGE'))* 'ON' 'DATABASE' QNAME | (grantPrivilege_schemaPrivileges | 'ALL' 'PRIVILEGES'?) 'ON' ('SCHEMA' QNAME | 'ALL' 'SCHEMAS' 'IN' 'DATABASE' ID) | (grantPrivilege_schemaPrivileges | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SCHEMAS' 'IN' 'DATABASE' ID | ((grantPrivilege_schemaObjectPrivileges_1 | 'ALL' 'PRIVILEGES'?) 'ON' ('AGENT' QNAME | 'ALL' 'AGENTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('AGGREGATION' 'POLICY' QNAME | 'ALL' 'AGGREGATION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' ('ALERT' QNAME | 'ALL' 'ALERTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' ('APPLICATION' 'SERVICE' QNAME | 'ALL' 'APPLICATION' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('ARTIFACT' 'REPOSITORY' QNAME | 'ALL' 'ARTIFACT' 'REPOSITORIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('AUTHENTICATION' 'POLICY' QNAME | 'ALL' 'AUTHENTICATION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('CONTACT' QNAME | 'ALL' 'CONTACTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' ('CORTEX' 'EXTENSION' QNAME | 'ALL' 'CORTEX' 'EXTENSIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' ('CORTEX' 'SEARCH' 'SERVICE' QNAME | 'ALL' 'CORTEX' 'SEARCH' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('DATA' 'METRIC' 'FUNCTION' QNAME | 'ALL' 'DATA' 'METRIC' 'FUNCTIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('DATASET' QNAME | 'ALL' 'DATASETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_11 | 'ALL' 'PRIVILEGES'?) 'ON' ('DBT' 'PROJECT' QNAME | 'ALL' 'DBT' 'PROJECTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('DCM' 'PROJECT' QNAME | 'ALL' 'DCM' 'PROJECTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_13 | 'ALL' 'PRIVILEGES'?) 'ON' ('DYNAMIC' 'TABLE' QNAME | 'ALL' 'DYNAMIC' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_15 | 'ALL' 'PRIVILEGES'?) 'ON' ('EVENT' 'TABLE' QNAME | 'ALL' 'EVENT' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_14 | 'ALL' 'PRIVILEGES'?) 'ON' ('EXPERIMENT' QNAME | 'ALL' 'EXPERIMENTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('FILE' 'FORMAT' QNAME | 'ALL' 'FILE' 'FORMATS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('FUNCTION' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' | 'ALL' 'FUNCTIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_16 | 'ALL' 'PRIVILEGES'?) 'ON' ('GATEWAY' QNAME | 'ALL' 'GATEWAYS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_17 | 'ALL' 'PRIVILEGES'?) 'ON' ('GIT' 'REPOSITORY' QNAME | 'ALL' 'GIT' 'REPOSITORIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_18 | 'ALL' 'PRIVILEGES'?) 'ON' ('HYBRID' 'TABLE' QNAME | 'ALL' 'HYBRID' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_20 | 'ALL' 'PRIVILEGES'?) 'ON' ('ICEBERG' 'TABLE' QNAME | 'ALL' 'ICEBERG' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_19 | 'ALL' 'PRIVILEGES'?) 'ON' ('IMAGE' 'REPOSITORY' QNAME | 'ALL' 'IMAGE' 'REPOSITORIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_21 | 'ALL' 'PRIVILEGES'?) 'ON' ('INTERACTIVE' 'TABLE' QNAME | 'ALL' 'INTERACTIVE' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('JOIN' 'POLICY' QNAME | 'ALL' 'JOIN' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' 'MAINTENANCE' 'POLICY' QNAME | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' 'MASKING' 'POLICY' QNAME | (grantPrivilege_schemaObjectPrivileges_22 | 'ALL' 'PRIVILEGES'?) 'ON' ('MATERIALIZED' 'VIEW' QNAME | 'ALL' 'MATERIALIZED' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_23 | 'ALL' 'PRIVILEGES'?) 'ON' ('MCP' 'SERVER' QNAME | 'ALL' 'MCP' 'SERVERS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('MODEL' QNAME | 'ALL' 'MODELS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_24 | 'ALL' 'PRIVILEGES'?) 'ON' ('ONLINE' 'FEATURE' 'TABLE' QNAME | 'ALL' 'ONLINE' 'FEATURE' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('PACKAGES' 'POLICY' QNAME | 'ALL' 'PACKAGES' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('PASSWORD' 'POLICY' QNAME | 'ALL' 'PASSWORD' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_25 | 'ALL' 'PRIVILEGES'?) 'ON' 'PIPE' QNAME | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('PRIVACY' 'POLICY' QNAME | 'ALL' 'PRIVACY' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('PROCEDURE' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' | 'ALL' 'PROCEDURES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('PROJECTION' 'POLICY' QNAME | 'ALL' 'PROJECTION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('ROW' 'ACCESS' 'POLICY' QNAME | 'ALL' 'ROW' 'ACCESS' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_9 | grantPrivilege_schemaObjectPrivileges_27 | 'ALL' 'PRIVILEGES'?) 'ON' ('SECRET' QNAME | 'ALL' 'SECRETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_28 | 'ALL' 'PRIVILEGES'?) 'ON' ('SEMANTIC' 'VIEW' QNAME | 'ALL' 'SEMANTIC' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('SEQUENCE' QNAME | 'ALL' 'SEQUENCES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_29 | 'ALL' 'PRIVILEGES'?) 'ON' ('SERVICE' QNAME | 'ALL' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('SESSION' 'POLICY' QNAME | 'ALL' 'SESSION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('SNAPSHOT' QNAME | 'ALL' 'SNAPSHOTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' ('SNAPSHOT' 'POLICY' QNAME | 'ALL' 'SNAPSHOT' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' ('SNAPSHOT' 'SET' QNAME | 'ALL' 'SNAPSHOT' 'SETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_30 | grantPrivilege_schemaObjectPrivileges_31 | 'ALL' 'PRIVILEGES'?) 'ON' ('STAGE' QNAME | 'ALL' 'STAGES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('STORAGE' 'LIFECYCLE' 'POLICY' QNAME | 'ALL' 'STORAGE' 'LIFECYCLE' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_32 | 'ALL' 'PRIVILEGES'?) 'ON' ('STREAM' QNAME | 'ALL' 'STREAMS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_33 | 'ALL' 'PRIVILEGES'?) 'ON' ('STREAMLIT' QNAME | 'ALL' 'STREAMLITS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_34 | 'ALL' 'PRIVILEGES'?) 'ON' ('TABLE' QNAME | 'ALL' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_26 | grantPrivilege_schemaObjectPrivileges_35 | 'ALL' 'PRIVILEGES'?) 'ON' 'TAG' QNAME | (grantPrivilege_schemaObjectPrivileges_36 | 'ALL' 'PRIVILEGES'?) 'ON' ('TASK' QNAME | 'ALL' 'TASKS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('TYPE' QNAME | 'ALL' 'TYPES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_37 | 'ALL' 'PRIVILEGES'?) 'ON' ('VIEW' QNAME | 'ALL' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilege_schemaObjectPrivileges_38 | 'ALL' 'PRIVILEGES'?) 'ON' ('WORKSPACE' QNAME | 'ALL' 'WORKSPACES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME))) | ((grantPrivilege_schemaObjectPrivileges_1 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'AGENTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'ALERTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'APPLICATION' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'ARTIFACT' 'REPOSITORIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'AUTHENTICATION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'CONTACTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'CORTEX' 'EXTENSIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'CORTEX' 'SEARCH' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DATA' 'METRIC' 'FUNCTIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DATASETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_11 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DBT' 'PROJECTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DCM' 'PROJECTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_13 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DYNAMIC' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_15 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'EVENT' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_14 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'EXPERIMENTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'FILE' 'FORMATS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'FUNCTIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_16 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'GATEWAYS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_17 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'GIT' 'REPOSITORIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_18 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'HYBRID' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_20 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'ICEBERG' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_21 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'INTERACTIVE' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_22 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'MATERIALIZED' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_23 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'MCP' 'SERVERS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'MODELS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_24 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'ONLINE' 'FEATURE' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PASSWORD' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_25 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PIPES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PRIVACY' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PROCEDURES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_9 | grantPrivilege_schemaObjectPrivileges_27 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SECRETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_28 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SEMANTIC' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SEQUENCES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_29 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SNAPSHOT' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SNAPSHOT' 'SETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_30 | grantPrivilege_schemaObjectPrivileges_31 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'STAGES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_32 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'STREAMS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_33 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'STREAMLITS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_34 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_36 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'TASKS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'TYPES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_37 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (grantPrivilege_schemaObjectPrivileges_38 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'WORKSPACES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME))) 'TO' 'DATABASE' 'ROLE' ID ('WITH' 'GRANT' 'OPTION')?
    ;
grantPrivilege_globalPrivileges
    : ('CREATE' ('ACCOUNT' | 'APPLICATION' | 'APPLICATION' 'PACKAGE' | 'COMPUTE' 'POOL' | 'LISTING' | 'DATABASE' | 'EXTERNAL' 'VOLUME' | 'FAILOVER' 'GROUP' | 'INTEGRATION' | 'NETWORK' 'POLICY' | 'ORGANIZATION' 'LISTING' | 'ORGANIZATION' 'PROFILE' | 'PREVIEW' 'APPLICATION' | 'REPLICATION' 'GROUP' | 'ROLE' | 'SHARE' | 'USER' | 'WAREHOUSE') | 'ATTACH' 'POLICY' | 'AUDIT' | 'BIND' 'SERVICE' 'ENDPOINT' | 'APPLY' (('AGGREGATION' | 'AUTHENTICATION' | 'JOIN' | 'MAINTENANCE' | 'MASKING' | 'PACKAGES' | 'PASSWORD' | 'PROJECTION' | 'ROW' 'ACCESS' | 'SESSION' | 'STORAGE' 'LIFECYCLE') 'POLICY' | 'CONTACT' | 'TAG') | 'EXECUTE' ('ALERT' | 'DATA' 'METRIC' 'FUNCTION' | 'MANAGED' 'ALERT' | 'MANAGED' 'TASK' | 'TASK') | 'IMPORT' ('SHARE' | 'ORGANIZATION' 'LISTING') | 'MANAGE' ('ACCOUNT' 'SUPPORT' 'CASES' | 'APPLICATION' 'SPECIFICATIONS' | 'CORTEX' 'EXTENSION' 'SECURITY' 'SCAN' | 'EVENT' 'SHARING' | 'GRANTS' | 'LISTING' 'AUTO' 'FULFILLMENT' | 'ORGANIZATION' 'SUPPORT' 'CASES' | 'SHARE' 'TARGET' | 'USER' 'SUPPORT' 'CASES' | 'VISIBILITY' | 'WAREHOUSES') | 'MODIFY' ('LOG' 'LEVEL' | 'TRACE' 'LEVEL' | 'SESSION' 'LOG' 'LEVEL' | 'SESSION' 'TRACE' 'LEVEL') | 'MONITOR' ('CORTEX' 'EXTENSIONS' | 'EXECUTION' | 'SECURITY' | 'USAGE') | 'OVERRIDE' 'SHARE' 'RESTRICTIONS' | 'PURCHASE' 'DATA' 'EXCHANGE' 'LISTING' | 'RESOLVE' 'ALL' | 'READ' 'SESSION' | 'READ' 'UNREDACTED' 'ERROR' 'TABLE' | 'USE' 'AI' 'FUNCTION' QNAME | 'USE' 'AI' 'FUNCTIONS') (',' ('CREATE' ('ACCOUNT' | 'APPLICATION' | 'APPLICATION' 'PACKAGE' | 'COMPUTE' 'POOL' | 'LISTING' | 'DATABASE' | 'EXTERNAL' 'VOLUME' | 'FAILOVER' 'GROUP' | 'INTEGRATION' | 'NETWORK' 'POLICY' | 'ORGANIZATION' 'LISTING' | 'ORGANIZATION' 'PROFILE' | 'PREVIEW' 'APPLICATION' | 'REPLICATION' 'GROUP' | 'ROLE' | 'SHARE' | 'USER' | 'WAREHOUSE') | 'ATTACH' 'POLICY' | 'AUDIT' | 'BIND' 'SERVICE' 'ENDPOINT' | 'APPLY' (('AGGREGATION' | 'AUTHENTICATION' | 'JOIN' | 'MAINTENANCE' | 'MASKING' | 'PACKAGES' | 'PASSWORD' | 'PROJECTION' | 'ROW' 'ACCESS' | 'SESSION' | 'STORAGE' 'LIFECYCLE') 'POLICY' | 'CONTACT' | 'TAG') | 'EXECUTE' ('ALERT' | 'DATA' 'METRIC' 'FUNCTION' | 'MANAGED' 'ALERT' | 'MANAGED' 'TASK' | 'TASK') | 'IMPORT' ('SHARE' | 'ORGANIZATION' 'LISTING') | 'MANAGE' ('ACCOUNT' 'SUPPORT' 'CASES' | 'APPLICATION' 'SPECIFICATIONS' | 'CORTEX' 'EXTENSION' 'SECURITY' 'SCAN' | 'EVENT' 'SHARING' | 'GRANTS' | 'LISTING' 'AUTO' 'FULFILLMENT' | 'ORGANIZATION' 'SUPPORT' 'CASES' | 'SHARE' 'TARGET' | 'USER' 'SUPPORT' 'CASES' | 'VISIBILITY' | 'WAREHOUSES') | 'MODIFY' ('LOG' 'LEVEL' | 'TRACE' 'LEVEL' | 'SESSION' 'LOG' 'LEVEL' | 'SESSION' 'TRACE' 'LEVEL') | 'MONITOR' ('CORTEX' 'EXTENSIONS' | 'EXECUTION' | 'SECURITY' | 'USAGE') | 'OVERRIDE' 'SHARE' 'RESTRICTIONS' | 'PURCHASE' 'DATA' 'EXCHANGE' 'LISTING' | 'RESOLVE' 'ALL' | 'READ' 'SESSION' | 'READ' 'UNREDACTED' 'ERROR' 'TABLE' | 'USE' 'AI' 'FUNCTION' QNAME | 'USE' 'AI' 'FUNCTIONS'))*
    ;
grantPrivilege_accountObjectPrivileges_1
    : ('ATTACH' 'LISTING' | 'DEVELOP' | 'INSTALL' | 'MANAGE' 'VERSIONS' | 'MANAGE' 'RELEASES') (',' ('ATTACH' 'LISTING' | 'DEVELOP' | 'INSTALL' | 'MANAGE' 'VERSIONS' | 'MANAGE' 'RELEASES'))*
    ;
grantPrivilege_accountObjectPrivileges_2
    : ('MODIFY' | 'MONITOR' | 'OPERATE' | 'USAGE') (',' ('MODIFY' | 'MONITOR' | 'OPERATE' | 'USAGE'))*
    ;
grantPrivilege_accountObjectPrivileges_3
    : 'FAILOVER' (',' 'FAILOVER')*
    ;
grantPrivilege_accountObjectPrivileges_4
    : ('MODIFY' | 'USAGE') (',' ('MODIFY' | 'USAGE'))*
    ;
grantPrivilege_accountObjectPrivileges_5
    : ('APPLYBUDGET' | 'CREATE' ('DATABASE' 'ROLE' | 'SCHEMA') | 'IMPORTED' 'PRIVILEGES' | 'MODIFY' | 'MONITOR' | 'USAGE') (',' ('APPLYBUDGET' | 'CREATE' ('DATABASE' 'ROLE' | 'SCHEMA') | 'IMPORTED' 'PRIVILEGES' | 'MODIFY' | 'MONITOR' | 'USAGE'))*
    ;
grantPrivilege_accountObjectPrivileges_6
    : 'USAGE' (',' 'USAGE')*
    ;
grantPrivilege_accountObjectPrivileges_7
    : ('FAILOVER' | 'MODIFY' | 'MONITOR' | 'REPLICATE') (',' ('FAILOVER' | 'MODIFY' | 'MONITOR' | 'REPLICATE'))*
    ;
grantPrivilege_accountObjectPrivileges_8
    : ('USAGE' | 'USE_ANY_ROLE') (',' ('USAGE' | 'USE_ANY_ROLE'))*
    ;
grantPrivilege_accountObjectPrivileges_10
    : ('MODIFY' | 'MONITOR' | 'REPLICATE') (',' ('MODIFY' | 'MONITOR' | 'REPLICATE'))*
    ;
grantPrivilege_accountObjectPrivileges_11
    : ('MODIFY' | 'MONITOR') (',' ('MODIFY' | 'MONITOR'))*
    ;
grantPrivilege_accountObjectPrivileges_12
    : ('IMPERSONATE' | 'MODIFY' 'PROGRAMMATIC' 'AUTHENTICATION' 'METHODS' | 'MONITOR') (',' ('IMPERSONATE' | 'MODIFY' 'PROGRAMMATIC' 'AUTHENTICATION' 'METHODS' | 'MONITOR'))*
    ;
grantPrivilege_accountObjectPrivileges_13
    : ('APPLYBUDGET' | 'MODIFY' | 'MONITOR' | 'USAGE' | 'OPERATE') (',' ('APPLYBUDGET' | 'MODIFY' | 'MONITOR' | 'USAGE' | 'OPERATE'))*
    ;
grantPrivilege_schemaPrivileges
    : ('ADD' 'SEARCH' 'OPTIMIZATION' | 'APPLYBUDGET' | 'CREATE' ('AGENT' | 'ALERT' | 'APPLICATION' 'SERVICE' | 'ARTIFACT' 'REPOSITORY' | 'CONTACT' | 'CORTEX' 'SEARCH' 'SERVICE' | 'DATA' 'METRIC' 'FUNCTION' | 'DATASET' | 'DBT' 'PROJECT' | 'DCM' 'PROJECT' | 'EVENT' 'TABLE' | 'EXPERIMENT' | 'FILE' 'FORMAT' | 'FUNCTION' | 'GATEWAY' | ('GIT' | 'IMAGE') 'REPOSITORY' | 'MCP' 'SERVER' | 'MODEL' | 'NETWORK' 'RULE' | 'NOTEBOOK' | 'PIPE' | 'PROCEDURE' | ('AGGREGATION' | 'AUTHENTICATION' | 'MAINTENANCE' | 'MASKING' | 'PACKAGES' | 'PASSWORD' | 'PRIVACY' | 'PROJECTION' | 'ROW' 'ACCESS' | 'SESSION' | 'STORAGE' 'LIFECYCLE') 'POLICY' | 'SECRET' | 'SEQUENCE' | 'SERVICE' | 'SNAPSHOT' | 'SNAPSHOT' 'POLICY' | 'SNAPSHOT' 'SET' | 'STAGE' | 'STREAM' | 'STREAMLIT' | 'SNOWFLAKE' '.' 'CORE' '.' 'BUDGET' | 'SNOWFLAKE' '.' 'DATA_PRIVACY' '.' 'CLASSIFICATION_PROFILE' | 'SNOWFLAKE' '.' 'DATA_PRIVACY' '.' 'CUSTOM_CLASSIFIER' | 'SNOWFLAKE' '.' 'ML' '.' 'ANOMALY_DETECTION' | 'SNOWFLAKE' '.' 'ML' '.' 'CLASSIFICATION' | 'SNOWFLAKE' '.' 'ML' '.' 'FORECAST' | 'SNOWFLAKE' '.' 'ML' '.' 'TOP_INSIGHTS' | 'SNOWFLAKE' '.' 'ML' '.' 'DOCUMENT_INTELLIGENCE' | (('DYNAMIC' | 'EXTERNAL' | 'HYBRID' | 'ICEBERG' | 'INTERACTIVE' | 'ONLINE' 'FEATURE'))? 'TABLE' | 'TAG' | 'TASK' | 'TYPE' | 'WORKSPACE' | (('MATERIALIZED' | 'SEMANTIC'))? 'VIEW') | 'MODIFY' | 'MONITOR' | 'USAGE') (',' ('ADD' 'SEARCH' 'OPTIMIZATION' | 'APPLYBUDGET' | 'CREATE' ('AGENT' | 'ALERT' | 'APPLICATION' 'SERVICE' | 'ARTIFACT' 'REPOSITORY' | 'CONTACT' | 'CORTEX' 'SEARCH' 'SERVICE' | 'DATA' 'METRIC' 'FUNCTION' | 'DATASET' | 'DBT' 'PROJECT' | 'DCM' 'PROJECT' | 'EVENT' 'TABLE' | 'EXPERIMENT' | 'FILE' 'FORMAT' | 'FUNCTION' | 'GATEWAY' | ('GIT' | 'IMAGE') 'REPOSITORY' | 'MCP' 'SERVER' | 'MODEL' | 'NETWORK' 'RULE' | 'NOTEBOOK' | 'PIPE' | 'PROCEDURE' | ('AGGREGATION' | 'AUTHENTICATION' | 'MAINTENANCE' | 'MASKING' | 'PACKAGES' | 'PASSWORD' | 'PRIVACY' | 'PROJECTION' | 'ROW' 'ACCESS' | 'SESSION' | 'STORAGE' 'LIFECYCLE') 'POLICY' | 'SECRET' | 'SEQUENCE' | 'SERVICE' | 'SNAPSHOT' | 'SNAPSHOT' 'POLICY' | 'SNAPSHOT' 'SET' | 'STAGE' | 'STREAM' | 'STREAMLIT' | 'SNOWFLAKE' '.' 'CORE' '.' 'BUDGET' | 'SNOWFLAKE' '.' 'DATA_PRIVACY' '.' 'CLASSIFICATION_PROFILE' | 'SNOWFLAKE' '.' 'DATA_PRIVACY' '.' 'CUSTOM_CLASSIFIER' | 'SNOWFLAKE' '.' 'ML' '.' 'ANOMALY_DETECTION' | 'SNOWFLAKE' '.' 'ML' '.' 'CLASSIFICATION' | 'SNOWFLAKE' '.' 'ML' '.' 'FORECAST' | 'SNOWFLAKE' '.' 'ML' '.' 'TOP_INSIGHTS' | 'SNOWFLAKE' '.' 'ML' '.' 'DOCUMENT_INTELLIGENCE' | (('DYNAMIC' | 'EXTERNAL' | 'HYBRID' | 'ICEBERG' | 'INTERACTIVE' | 'ONLINE' 'FEATURE'))? 'TABLE' | 'TAG' | 'TASK' | 'TYPE' | 'WORKSPACE' | (('MATERIALIZED' | 'SEMANTIC'))? 'VIEW') | 'MODIFY' | 'MONITOR' | 'USAGE'))*
    ;
grantPrivilege_schemaObjectPrivileges_1
    : ('MODIFY' | 'MONITOR' | 'USAGE') (',' ('MODIFY' | 'MONITOR' | 'USAGE'))*
    ;
grantPrivilege_schemaObjectPrivileges_2
    : ('MONITOR' | 'OPERATE') (',' ('MONITOR' | 'OPERATE'))*
    ;
grantPrivilege_schemaObjectPrivileges_3
    : ('MONITOR' | 'OPERATE' | 'USAGE') (',' ('MONITOR' | 'OPERATE' | 'USAGE'))*
    ;
grantPrivilege_schemaObjectPrivileges_4
    : 'READ' (',' 'READ')*
    ;
grantPrivilege_schemaObjectPrivileges_5
    : ('APPLY' | 'MODIFY') (',' ('APPLY' | 'MODIFY'))*
    ;
grantPrivilege_schemaObjectPrivileges_6
    : ('READ' | 'WRITE') (',' ('READ' | 'WRITE'))*
    ;
grantPrivilege_schemaObjectPrivileges_7
    : ('MONITOR' | 'OPERATE' | 'USAGE') (',' ('MONITOR' | 'OPERATE' | 'USAGE'))*
    ;
grantPrivilege_schemaObjectPrivileges_8
    : 'USAGE' (',' 'USAGE')*
    ;
grantPrivilege_schemaObjectPrivileges_9
    : 'USAGE' (',' 'USAGE')*
    ;
grantPrivilege_schemaObjectPrivileges_10
    : 'USAGE' (',' 'USAGE')*
    ;
grantPrivilege_schemaObjectPrivileges_11
    : 'USAGE' ',' 'MONITOR' (',' 'MONITOR')*
    ;
grantPrivilege_schemaObjectPrivileges_12
    : 'READ' ',' 'MONITOR' (',' 'MONITOR')*
    ;
grantPrivilege_schemaObjectPrivileges_13
    : 'MONITOR' ',' 'OPERATE' ',' 'SELECT' (',' 'SELECT')*
    ;
grantPrivilege_schemaObjectPrivileges_14
    : ('CREATE' | 'MODIFY' | 'USAGE') (',' ('CREATE' | 'MODIFY' | 'USAGE'))*
    ;
grantPrivilege_schemaObjectPrivileges_15
    : ('APPLYBUDGET' | 'DELETE' | 'OWNERSHIP' | 'REFERENCES' | 'SELECT' | 'TRUNCATE') (',' ('APPLYBUDGET' | 'DELETE' | 'OWNERSHIP' | 'REFERENCES' | 'SELECT' | 'TRUNCATE'))*
    ;
grantPrivilege_schemaObjectPrivileges_16
    : ('CREATE' | 'MODIFY' | 'USAGE') (',' ('CREATE' | 'MODIFY' | 'USAGE'))*
    ;
grantPrivilege_schemaObjectPrivileges_17
    : 'READ' ',' 'WRITE' (',' 'READ' ',' 'WRITE')*
    ;
grantPrivilege_schemaObjectPrivileges_18
    : ('APPLYBUDGET' | 'DELETE' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE') (',' ('APPLYBUDGET' | 'DELETE' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE'))*
    ;
grantPrivilege_schemaObjectPrivileges_19
    : 'READ' ',' 'WRITE' (',' 'READ' ',' 'WRITE')*
    ;
grantPrivilege_schemaObjectPrivileges_20
    : ('APPLYBUDGET' | 'DELETE' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE') (',' ('APPLYBUDGET' | 'DELETE' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE'))*
    ;
grantPrivilege_schemaObjectPrivileges_21
    : ('REFERENCES' | 'SELECT') (',' ('REFERENCES' | 'SELECT'))*
    ;
grantPrivilege_schemaObjectPrivileges_22
    : ('APPLYBUDGET' | 'REFERENCES' | 'SELECT') (',' ('APPLYBUDGET' | 'REFERENCES' | 'SELECT'))*
    ;
grantPrivilege_schemaObjectPrivileges_23
    : ('MODIFY' | 'USAGE') (',' ('MODIFY' | 'USAGE'))*
    ;
grantPrivilege_schemaObjectPrivileges_24
    : ('MONITOR' | 'SELECT') (',' ('MONITOR' | 'SELECT'))*
    ;
grantPrivilege_schemaObjectPrivileges_25
    : ('APPLYBUDGET' | 'MONITOR' | 'OPERATE') (',' ('APPLYBUDGET' | 'MONITOR' | 'OPERATE'))*
    ;
grantPrivilege_schemaObjectPrivileges_26
    : 'APPLY' (',' 'APPLY')*
    ;
grantPrivilege_schemaObjectPrivileges_27
    : ('READ' | 'USAGE') (',' ('READ' | 'USAGE'))*
    ;
grantPrivilege_schemaObjectPrivileges_28
    : ('SELECT' | 'REFERENCES' | 'MONITOR') (',' ('SELECT' | 'REFERENCES' | 'MONITOR'))*
    ;
grantPrivilege_schemaObjectPrivileges_29
    : ('MONITOR' | 'OPERATE') (',' ('MONITOR' | 'OPERATE'))*
    ;
grantPrivilege_schemaObjectPrivileges_30
    : 'USAGE' (',' 'USAGE')*
    ;
grantPrivilege_schemaObjectPrivileges_31
    : 'READ' (',' 'WRITE')? (',' 'READ' (',' 'WRITE')?)*
    ;
grantPrivilege_schemaObjectPrivileges_32
    : 'SELECT' (',' 'SELECT')*
    ;
grantPrivilege_schemaObjectPrivileges_33
    : 'USAGE' (',' 'USAGE')*
    ;
grantPrivilege_schemaObjectPrivileges_34
    : ('APPLYBUDGET' | 'DELETE' | 'EVOLVE' 'SCHEMA' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'SELECT' 'ERROR' 'TABLE' | 'TRUNCATE' | 'UPDATE') (',' ('APPLYBUDGET' | 'DELETE' | 'EVOLVE' 'SCHEMA' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'SELECT' 'ERROR' 'TABLE' | 'TRUNCATE' | 'UPDATE'))*
    ;
grantPrivilege_schemaObjectPrivileges_35
    : 'READ'
    ;
grantPrivilege_schemaObjectPrivileges_36
    : ('APPLYBUDGET' | 'MONITOR' | 'OPERATE') (',' ('APPLYBUDGET' | 'MONITOR' | 'OPERATE'))*
    ;
grantPrivilege_schemaObjectPrivileges_37
    : ('REFERENCES' | 'SELECT') (',' ('REFERENCES' | 'SELECT'))*
    ;
grantPrivilege_schemaObjectPrivileges_38
    : ('READ' | 'WRITE') (',' ('READ' | 'WRITE'))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/grant-privilege-application
// fix: ON { <object_type> ... is not closed; closed
// fix: read: each -- For X section of accountObjectPrivileges and schemaObjectPrivileges is the privileges on X alone, and the ON clause chooses an object type with its own privileges
grantPrivilegeApplication
    : 'GRANT' (grantPrivilegeApplication_globalPrivileges 'ON' 'ACCOUNT' | ((grantPrivilegeApplication_accountObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'USER' ID | (grantPrivilegeApplication_accountObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' 'RESOURCE' 'MONITOR' ID | (grantPrivilegeApplication_accountObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' 'WAREHOUSE' ID | (grantPrivilegeApplication_accountObjectPrivileges_1 | 'ALL' 'PRIVILEGES'?) 'ON' 'COMPUTE' 'POOL' ID | (grantPrivilegeApplication_accountObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' 'DATABASE' ID | (grantPrivilegeApplication_accountObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' 'INTEGRATION' ID | (grantPrivilegeApplication_accountObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' 'CONNECTION' ID | (grantPrivilegeApplication_accountObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' 'FAILOVER' 'GROUP' ID | (grantPrivilegeApplication_accountObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' 'REPLICATION' 'GROUP' ID | (grantPrivilegeApplication_accountObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' 'EXTERNAL' 'VOLUME' ID | 'ALL' 'PRIVILEGES'? 'ON' 'APPLICATION' 'PACKAGE' ID | 'ALL' 'PRIVILEGES'? 'ON' 'APPLICATION' ID) | (grantPrivilegeApplication_schemaPrivileges | 'ALL' 'PRIVILEGES'?) 'ON' ('SCHEMA' QNAME | 'ALL' 'SCHEMAS' 'IN' 'DATABASE' ID) | ((grantPrivilegeApplication_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('AGGREGATION' 'POLICY' QNAME | 'ALL' 'AGGREGATION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_1 | 'ALL' 'PRIVILEGES'?) 'ON' ('ALERT' QNAME | 'ALL' 'ALERTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' ('DYNAMIC' 'TABLE' QNAME | 'ALL' 'DYNAMIC' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' ('EVENT' 'TABLE' QNAME | 'ALL' 'EVENT' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('FILE' 'FORMAT' QNAME | 'ALL' 'FILE' 'FORMATS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('FUNCTION' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' | 'ALL' 'FUNCTIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('IMAGE' 'REPOSITORY' QNAME | 'ALL' 'IMAGE' 'REPOSITORIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' 'MASKING' 'POLICY' QNAME | (grantPrivilegeApplication_schemaObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' ('MATERIALIZED' 'VIEW' QNAME | 'ALL' 'MATERIALIZED' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('PACKAGES' 'POLICY' QNAME | 'ALL' 'PACKAGES' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('PASSWORD' 'POLICY' QNAME | 'ALL' 'PASSWORD' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' 'PIPE' QNAME | (grantPrivilegeApplication_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('PROCEDURE' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' | 'ALL' 'PROCEDURES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('PROJECTION' 'POLICY' QNAME | 'ALL' 'PROJECTION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('ROW' 'ACCESS' 'POLICY' QNAME | 'ALL' 'ROW' 'ACCESS' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_4 | grantPrivilegeApplication_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('SECRET' QNAME | 'ALL' 'SECRETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' ('SEMANTIC' 'VIEW' QNAME | 'ALL' 'SEMANTIC' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('SEQUENCE' QNAME | 'ALL' 'SEQUENCES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_11 | 'ALL' 'PRIVILEGES'?) 'ON' ('SERVICE' QNAME | 'ALL' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('SESSION' 'POLICY' QNAME | 'ALL' 'SESSION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('SNAPSHOT' QNAME | 'ALL' 'SNAPSHOTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_12 | grantPrivilegeApplication_schemaObjectPrivileges_13 | 'ALL' 'PRIVILEGES'?) 'ON' ('STAGE' QNAME | 'ALL' 'STAGES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_14 | 'ALL' 'PRIVILEGES'?) 'ON' ('STREAM' QNAME | 'ALL' 'STREAMS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_15 | 'ALL' 'PRIVILEGES'?) 'ON' ('TABLE' QNAME | 'ALL' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_8 | grantPrivilegeApplication_schemaObjectPrivileges_16 | 'ALL' 'PRIVILEGES'?) 'ON' 'TAG' QNAME | (grantPrivilegeApplication_schemaObjectPrivileges_17 | 'ALL' 'PRIVILEGES'?) 'ON' ('TASK' QNAME | 'ALL' 'TASKS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('TYPE' QNAME | 'ALL' 'TYPES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplication_schemaObjectPrivileges_18 | 'ALL' 'PRIVILEGES'?) 'ON' ('VIEW' QNAME | 'ALL' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)))) 'TO' 'APPLICATION' QNAME
    ;
grantPrivilegeApplication_globalPrivileges
    : ('CREATE' ('COMPUTE' 'POOL' | 'DATABASE' | 'WAREHOUSE') | 'BIND' 'SERVICE' 'ENDPOINT' | 'EXECUTE' 'MANAGED' 'TASK' | 'MANAGE' 'WAREHOUSES' | 'READ' 'SESSION') (',' ('CREATE' ('COMPUTE' 'POOL' | 'DATABASE' | 'WAREHOUSE') | 'BIND' 'SERVICE' 'ENDPOINT' | 'EXECUTE' 'MANAGED' 'TASK' | 'MANAGE' 'WAREHOUSES' | 'READ' 'SESSION'))*
    ;
grantPrivilegeApplication_accountObjectPrivileges_1
    : ('MODIFY' | 'MONITOR' | 'OPERATE' | 'USAGE') (',' ('MODIFY' | 'MONITOR' | 'OPERATE' | 'USAGE'))*
    ;
grantPrivilegeApplication_accountObjectPrivileges_2
    : 'FAILOVER' (',' 'FAILOVER')*
    ;
grantPrivilegeApplication_accountObjectPrivileges_3
    : ('APPLYBUDGET' | 'CREATE' ('DATABASE' 'ROLE' | 'SCHEMA') | 'IMPORTED' 'PRIVILEGES' | 'MODIFY' | 'MONITOR' | 'USAGE') (',' ('APPLYBUDGET' | 'CREATE' ('DATABASE' 'ROLE' | 'SCHEMA') | 'IMPORTED' 'PRIVILEGES' | 'MODIFY' | 'MONITOR' | 'USAGE'))*
    ;
grantPrivilegeApplication_accountObjectPrivileges_4
    : 'USAGE' (',' 'USAGE')*
    ;
grantPrivilegeApplication_accountObjectPrivileges_5
    : ('FAILOVER' | 'MODIFY' | 'MONITOR' | 'REPLICATE') (',' ('FAILOVER' | 'MODIFY' | 'MONITOR' | 'REPLICATE'))*
    ;
grantPrivilegeApplication_accountObjectPrivileges_6
    : ('USAGE' | 'USE_ANY_ROLE') (',' ('USAGE' | 'USE_ANY_ROLE'))*
    ;
grantPrivilegeApplication_accountObjectPrivileges_7
    : ('MODIFY' | 'MONITOR' | 'REPLICATE') (',' ('MODIFY' | 'MONITOR' | 'REPLICATE'))*
    ;
grantPrivilegeApplication_accountObjectPrivileges_8
    : ('MODIFY' | 'MONITOR') (',' ('MODIFY' | 'MONITOR'))*
    ;
grantPrivilegeApplication_accountObjectPrivileges_9
    : 'MONITOR' (',' 'MONITOR')*
    ;
grantPrivilegeApplication_accountObjectPrivileges_10
    : ('APPLYBUDGET' | 'MODIFY' | 'MONITOR' | 'USAGE' | 'OPERATE') (',' ('APPLYBUDGET' | 'MODIFY' | 'MONITOR' | 'USAGE' | 'OPERATE'))*
    ;
grantPrivilegeApplication_schemaPrivileges
    : ('ADD' 'SEARCH' 'OPTIMIZATION' | 'CREATE' ('ALERT' | 'EXTERNAL' 'TABLE' | 'FILE' 'FORMAT' | 'FUNCTION' | 'IMAGE' 'REPOSITORY' | 'MATERIALIZED' 'VIEW' | 'PIPE' | 'PROCEDURE' | ('AGGREGATION' | 'MASKING' | 'PASSWORD' | 'PROJECTION' | 'ROW' 'ACCESS' | 'SESSION') 'POLICY' | 'SECRET' | 'SEMANTIC' 'VIEW' | 'SEQUENCE' | 'SERVICE' | 'SNAPSHOT' | 'STAGE' | 'STREAM' | 'TAG' | 'TABLE' | 'TASK' | 'VIEW') | 'MODIFY' | 'MONITOR' | 'USAGE') (',' ('ADD' 'SEARCH' 'OPTIMIZATION' | 'CREATE' ('ALERT' | 'EXTERNAL' 'TABLE' | 'FILE' 'FORMAT' | 'FUNCTION' | 'IMAGE' 'REPOSITORY' | 'MATERIALIZED' 'VIEW' | 'PIPE' | 'PROCEDURE' | ('AGGREGATION' | 'MASKING' | 'PASSWORD' | 'PROJECTION' | 'ROW' 'ACCESS' | 'SESSION') 'POLICY' | 'SECRET' | 'SEMANTIC' 'VIEW' | 'SEQUENCE' | 'SERVICE' | 'SNAPSHOT' | 'STAGE' | 'STREAM' | 'TAG' | 'TABLE' | 'TASK' | 'VIEW') | 'MODIFY' | 'MONITOR' | 'USAGE'))*
    ;
grantPrivilegeApplication_schemaObjectPrivileges_1
    : ('MONITOR' | 'OPERATE') (',' ('MONITOR' | 'OPERATE'))*
    ;
grantPrivilegeApplication_schemaObjectPrivileges_2
    : 'OPERATE' ',' 'SELECT' (',' 'SELECT')*
    ;
grantPrivilegeApplication_schemaObjectPrivileges_3
    : ('INSERT' | 'SELECT') (',' ('INSERT' | 'SELECT'))*
    ;
grantPrivilegeApplication_schemaObjectPrivileges_4
    : 'USAGE' (',' 'USAGE')*
    ;
grantPrivilegeApplication_schemaObjectPrivileges_5
    : 'READ' ',' 'WRITE' (',' 'READ' ',' 'WRITE')*
    ;
grantPrivilegeApplication_schemaObjectPrivileges_6
    : ('APPLYBUDGET' | 'REFERENCES' | 'SELECT') (',' ('APPLYBUDGET' | 'REFERENCES' | 'SELECT'))*
    ;
grantPrivilegeApplication_schemaObjectPrivileges_7
    : ('APPLYBUDGET' | 'MONITOR' | 'OPERATE') (',' ('APPLYBUDGET' | 'MONITOR' | 'OPERATE'))*
    ;
grantPrivilegeApplication_schemaObjectPrivileges_8
    : 'APPLY' (',' 'APPLY')*
    ;
grantPrivilegeApplication_schemaObjectPrivileges_9
    : 'READ' ',' 'USAGE' (',' 'USAGE')*
    ;
grantPrivilegeApplication_schemaObjectPrivileges_10
    : 'REFERENCES' (',' 'REFERENCES')*
    ;
grantPrivilegeApplication_schemaObjectPrivileges_11
    : ('MONITOR' | 'OPERATE') (',' ('MONITOR' | 'OPERATE'))*
    ;
grantPrivilegeApplication_schemaObjectPrivileges_12
    : 'USAGE' (',' 'USAGE')*
    ;
grantPrivilegeApplication_schemaObjectPrivileges_13
    : 'READ' (',' 'WRITE')? (',' 'READ' (',' 'WRITE')?)*
    ;
grantPrivilegeApplication_schemaObjectPrivileges_14
    : 'SELECT' (',' 'SELECT')*
    ;
grantPrivilegeApplication_schemaObjectPrivileges_15
    : ('APPLYBUDGET' | 'DELETE' | 'EVOLVE' 'SCHEMA' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE') (',' ('APPLYBUDGET' | 'DELETE' | 'EVOLVE' 'SCHEMA' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE'))*
    ;
grantPrivilegeApplication_schemaObjectPrivileges_16
    : 'READ'
    ;
grantPrivilegeApplication_schemaObjectPrivileges_17
    : ('APPLYBUDGET' | 'MONITOR' | 'OPERATE') (',' ('APPLYBUDGET' | 'MONITOR' | 'OPERATE'))*
    ;
grantPrivilegeApplication_schemaObjectPrivileges_18
    : ('REFERENCES' | 'SELECT') (',' ('REFERENCES' | 'SELECT'))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/grant-privilege-application-role
// fix: ON { <object_type> ... is not closed; closed
// fix: read: each -- For X section of accountObjectPrivileges and schemaObjectPrivileges is the privileges on X alone, and the ON clause chooses an object type with its own privileges
grantPrivilegeApplicationRole
    : 'GRANT' ((grantPrivilegeApplicationRole_schemaPrivileges | 'ALL' 'PRIVILEGES'?) 'ON' 'SCHEMA' QNAME | ((grantPrivilegeApplicationRole_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('AGGREGATION' 'POLICY' QNAME | 'ALL' 'AGGREGATION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_1 | 'ALL' 'PRIVILEGES'?) 'ON' ('ALERT' QNAME | 'ALL' 'ALERTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' ('DYNAMIC' 'TABLE' QNAME | 'ALL' 'DYNAMIC' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' ('EVENT' 'TABLE' QNAME | 'ALL' 'EVENT' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('FILE' 'FORMAT' QNAME | 'ALL' 'FILE' 'FORMATS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('FUNCTION' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' | 'ALL' 'FUNCTIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('IMAGE' 'REPOSITORY' QNAME | 'ALL' 'IMAGE' 'REPOSITORIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' 'MASKING' 'POLICY' QNAME | (grantPrivilegeApplicationRole_schemaObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' ('MATERIALIZED' 'VIEW' QNAME | 'ALL' 'MATERIALIZED' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('PACKAGES' 'POLICY' QNAME | 'ALL' 'PACKAGES' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('PASSWORD' 'POLICY' QNAME | 'ALL' 'PASSWORD' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' 'PIPE' QNAME | (grantPrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('PROCEDURE' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' | 'ALL' 'PROCEDURES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('PROJECTION' 'POLICY' QNAME | 'ALL' 'PROJECTION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('ROW' 'ACCESS' 'POLICY' QNAME | 'ALL' 'ROW' 'ACCESS' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_4 | grantPrivilegeApplicationRole_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('SECRET' QNAME | 'ALL' 'SECRETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' ('SEMANTIC' 'VIEW' QNAME | 'ALL' 'SEMANTIC' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('SEQUENCE' QNAME | 'ALL' 'SEQUENCES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_11 | 'ALL' 'PRIVILEGES'?) 'ON' ('SERVICE' QNAME | 'ALL' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('SESSION' 'POLICY' QNAME | 'ALL' 'SESSION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('SNAPSHOT' QNAME | 'ALL' 'SNAPSHOTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_12 | grantPrivilegeApplicationRole_schemaObjectPrivileges_13 | 'ALL' 'PRIVILEGES'?) 'ON' ('STAGE' QNAME | 'ALL' 'STAGES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_14 | 'ALL' 'PRIVILEGES'?) 'ON' ('STREAM' QNAME | 'ALL' 'STREAMS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_15 | 'ALL' 'PRIVILEGES'?) 'ON' ('TABLE' QNAME | 'ALL' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_8 | grantPrivilegeApplicationRole_schemaObjectPrivileges_16 | 'ALL' 'PRIVILEGES'?) 'ON' 'TAG' QNAME | (grantPrivilegeApplicationRole_schemaObjectPrivileges_17 | 'ALL' 'PRIVILEGES'?) 'ON' ('TASK' QNAME | 'ALL' 'TASKS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('TYPE' QNAME | 'ALL' 'TYPES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeApplicationRole_schemaObjectPrivileges_18 | 'ALL' 'PRIVILEGES'?) 'ON' ('VIEW' QNAME | 'ALL' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME))) | ((grantPrivilegeApplicationRole_schemaObjectPrivileges_1 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'ALERTS' 'IN' 'SCHEMA' QNAME | (grantPrivilegeApplicationRole_schemaObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DYNAMIC' 'TABLES' 'IN' 'SCHEMA' QNAME | (grantPrivilegeApplicationRole_schemaObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'EVENT' 'TABLES' 'IN' 'SCHEMA' QNAME | (grantPrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'FILE' 'FORMATS' 'IN' 'SCHEMA' QNAME | (grantPrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'FUNCTIONS' 'IN' 'SCHEMA' QNAME | (grantPrivilegeApplicationRole_schemaObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'MATERIALIZED' 'VIEWS' 'IN' 'SCHEMA' QNAME | (grantPrivilegeApplicationRole_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PASSWORD' 'POLICIES' 'IN' 'SCHEMA' QNAME | (grantPrivilegeApplicationRole_schemaObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PIPES' 'IN' 'SCHEMA' QNAME | (grantPrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PROCEDURES' 'IN' 'SCHEMA' QNAME | (grantPrivilegeApplicationRole_schemaObjectPrivileges_4 | grantPrivilegeApplicationRole_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SECRETS' 'IN' 'SCHEMA' QNAME | (grantPrivilegeApplicationRole_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SEMANTIC' 'VIEWS' 'IN' 'SCHEMA' QNAME | (grantPrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SEQUENCES' 'IN' 'SCHEMA' QNAME | (grantPrivilegeApplicationRole_schemaObjectPrivileges_11 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SERVICES' 'IN' 'SCHEMA' QNAME | (grantPrivilegeApplicationRole_schemaObjectPrivileges_12 | grantPrivilegeApplicationRole_schemaObjectPrivileges_13 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'STAGES' 'IN' 'SCHEMA' QNAME | (grantPrivilegeApplicationRole_schemaObjectPrivileges_14 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'STREAMS' 'IN' 'SCHEMA' QNAME | (grantPrivilegeApplicationRole_schemaObjectPrivileges_15 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'TABLES' 'IN' 'SCHEMA' QNAME | (grantPrivilegeApplicationRole_schemaObjectPrivileges_17 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'TASKS' 'IN' 'SCHEMA' QNAME | (grantPrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'TYPES' 'IN' 'SCHEMA' QNAME | (grantPrivilegeApplicationRole_schemaObjectPrivileges_18 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'VIEWS' 'IN' 'SCHEMA' QNAME)) 'TO' 'APPLICATION' 'ROLE' QNAME ('WITH' 'GRANT' 'OPTION')?
    ;
grantPrivilegeApplicationRole_schemaPrivileges
    : ('ADD' 'SEARCH' 'OPTIMIZATION' | 'CREATE' ('ALERT' | 'EXTERNAL' 'TABLE' | 'FILE' 'FORMAT' | 'FUNCTION' | 'IMAGE' 'REPOSITORY' | 'MATERIALIZED' 'VIEW' | 'PIPE' | 'PROCEDURE' | ('AGGREGATION' | 'MASKING' | 'PASSWORD' | 'PROJECTION' | 'ROW' 'ACCESS' | 'SESSION') 'POLICY' | 'SECRET' | 'SEMANTIC' 'VIEW' | 'SEQUENCE' | 'SERVICE' | 'SNAPSHOT' | 'STAGE' | 'STREAM' | 'TAG' | 'TABLE' | 'TASK' | 'VIEW') | 'MODIFY' | 'MONITOR' | 'USAGE') (',' ('ADD' 'SEARCH' 'OPTIMIZATION' | 'CREATE' ('ALERT' | 'EXTERNAL' 'TABLE' | 'FILE' 'FORMAT' | 'FUNCTION' | 'IMAGE' 'REPOSITORY' | 'MATERIALIZED' 'VIEW' | 'PIPE' | 'PROCEDURE' | ('AGGREGATION' | 'MASKING' | 'PASSWORD' | 'PROJECTION' | 'ROW' 'ACCESS' | 'SESSION') 'POLICY' | 'SECRET' | 'SEMANTIC' 'VIEW' | 'SEQUENCE' | 'SERVICE' | 'SNAPSHOT' | 'STAGE' | 'STREAM' | 'TAG' | 'TABLE' | 'TASK' | 'VIEW') | 'MODIFY' | 'MONITOR' | 'USAGE'))*
    ;
grantPrivilegeApplicationRole_schemaObjectPrivileges_1
    : ('MONITOR' | 'OPERATE') (',' ('MONITOR' | 'OPERATE'))*
    ;
grantPrivilegeApplicationRole_schemaObjectPrivileges_2
    : 'OPERATE' ',' 'SELECT' (',' 'SELECT')*
    ;
grantPrivilegeApplicationRole_schemaObjectPrivileges_3
    : ('INSERT' | 'SELECT') (',' ('INSERT' | 'SELECT'))*
    ;
grantPrivilegeApplicationRole_schemaObjectPrivileges_4
    : 'USAGE' (',' 'USAGE')*
    ;
grantPrivilegeApplicationRole_schemaObjectPrivileges_5
    : 'READ' ',' 'WRITE' (',' 'READ' ',' 'WRITE')*
    ;
grantPrivilegeApplicationRole_schemaObjectPrivileges_6
    : ('APPLYBUDGET' | 'REFERENCES' | 'SELECT') (',' ('APPLYBUDGET' | 'REFERENCES' | 'SELECT'))*
    ;
grantPrivilegeApplicationRole_schemaObjectPrivileges_7
    : ('APPLYBUDGET' | 'MONITOR' | 'OPERATE') (',' ('APPLYBUDGET' | 'MONITOR' | 'OPERATE'))*
    ;
grantPrivilegeApplicationRole_schemaObjectPrivileges_8
    : 'APPLY' (',' 'APPLY')*
    ;
grantPrivilegeApplicationRole_schemaObjectPrivileges_9
    : 'READ' ',' 'USAGE' (',' 'USAGE')*
    ;
grantPrivilegeApplicationRole_schemaObjectPrivileges_10
    : 'REFERENCES' (',' 'REFERENCES')*
    ;
grantPrivilegeApplicationRole_schemaObjectPrivileges_11
    : ('MONITOR' | 'OPERATE') (',' ('MONITOR' | 'OPERATE'))*
    ;
grantPrivilegeApplicationRole_schemaObjectPrivileges_12
    : 'USAGE' (',' 'USAGE')*
    ;
grantPrivilegeApplicationRole_schemaObjectPrivileges_13
    : 'READ' (',' 'WRITE')? (',' 'READ' (',' 'WRITE')?)*
    ;
grantPrivilegeApplicationRole_schemaObjectPrivileges_14
    : 'SELECT' (',' 'SELECT')*
    ;
grantPrivilegeApplicationRole_schemaObjectPrivileges_15
    : ('APPLYBUDGET' | 'DELETE' | 'EVOLVE' 'SCHEMA' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE') (',' ('APPLYBUDGET' | 'DELETE' | 'EVOLVE' 'SCHEMA' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE'))*
    ;
grantPrivilegeApplicationRole_schemaObjectPrivileges_16
    : 'READ'
    ;
grantPrivilegeApplicationRole_schemaObjectPrivileges_17
    : ('APPLYBUDGET' | 'MONITOR' | 'OPERATE') (',' ('APPLYBUDGET' | 'MONITOR' | 'OPERATE'))*
    ;
grantPrivilegeApplicationRole_schemaObjectPrivileges_18
    : ('REFERENCES' | 'SELECT') (',' ('REFERENCES' | 'SELECT'))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/grant-privilege-share
grantPrivilegeShare
    : 'GRANT' grantPrivilegeShare_objectPrivilege 'ON' ('DATABASE' QNAME | 'SCHEMA' QNAME | 'FUNCTION' QNAME | 'MODEL' QNAME | 'SEMANTIC' 'VIEW' QNAME | ('TABLE' QNAME | 'ALL' 'TABLES' 'IN' 'SCHEMA' QNAME) | ('EXTERNAL' 'TABLE' QNAME | 'ALL' 'EXTERNAL' 'TABLES' 'IN' 'SCHEMA' QNAME) | ('ICEBERG' 'TABLE' QNAME | 'ALL' 'ICEBERG' 'TABLES' 'IN' 'SCHEMA' QNAME) | ('DYNAMIC' 'TABLE' QNAME | 'ALL' 'DYNAMIC' 'TABLES' 'IN' 'SCHEMA' QNAME) | 'TAG' QNAME | 'VIEW' QNAME) 'TO' 'SHARE' ID
    ;
grantPrivilegeShare_objectPrivilege
    : 'REFERENCE_USAGE' (',' 'REFERENCE_USAGE')*
    | 'USAGE' (',' 'USAGE')*
    | 'READ'
    | ('REFERENCES' | 'SELECT') (',' ('REFERENCES' | 'SELECT'))*
    | 'EVOLVE' 'SCHEMA' (',' 'EVOLVE' 'SCHEMA')*
    | 'SELECT' (',' 'SELECT')*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/grant-privilege-user
// fix: read: each -- For X section of accountObjectPrivileges and schemaObjectPrivileges is the privileges on X alone, and the ON clause chooses an object type with its own privileges
grantPrivilegeUser
    : 'GRANT' ((grantPrivilegeUser_globalPrivileges | 'ALL' 'PRIVILEGES'?) 'ON' 'ACCOUNT' | ((grantPrivilegeUser_accountObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'USER' ID | (grantPrivilegeUser_accountObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' 'RESOURCE' 'MONITOR' ID | (grantPrivilegeUser_accountObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' 'WAREHOUSE' ID | (grantPrivilegeUser_accountObjectPrivileges_1 | 'ALL' 'PRIVILEGES'?) 'ON' 'COMPUTE' 'POOL' ID | (grantPrivilegeUser_accountObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' 'DATABASE' ID | (grantPrivilegeUser_accountObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' 'INTEGRATION' ID | (grantPrivilegeUser_accountObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' 'CONNECTION' ID | (grantPrivilegeUser_accountObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' 'FAILOVER' 'GROUP' ID | (grantPrivilegeUser_accountObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' 'REPLICATION' 'GROUP' ID | (grantPrivilegeUser_accountObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' 'EXTERNAL' 'VOLUME' ID) | (grantPrivilegeUser_schemaPrivileges | 'ALL' 'PRIVILEGES'?) 'ON' ('SCHEMA' QNAME | 'ALL' 'SCHEMAS' 'IN' 'DATABASE' ID) | ((grantPrivilegeUser_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('AGGREGATION' 'POLICY' QNAME | 'ALL' 'AGGREGATION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_1 | 'ALL' 'PRIVILEGES'?) 'ON' ('ALERT' QNAME | 'ALL' 'ALERTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('AUTHENTICATION' 'POLICY' QNAME | 'ALL' 'AUTHENTICATION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' ('DATA' 'METRIC' 'FUNCTION' QNAME | 'ALL' 'DATA' 'METRIC' 'FUNCTIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' ('DYNAMIC' 'TABLE' QNAME | 'ALL' 'DYNAMIC' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('EVENT' 'TABLE' QNAME | 'ALL' 'EVENT' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('FILE' 'FORMAT' QNAME | 'ALL' 'FILE' 'FORMATS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('FUNCTION' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' | 'ALL' 'FUNCTIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' ('GIT' 'REPOSITORY' QNAME | 'ALL' 'GIT' 'REPOSITORIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' ('HYBRID' 'TABLE' QNAME | 'ALL' 'HYBRID' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('ICEBERG' 'TABLE' QNAME | 'ALL' 'ICEBERG' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('IMAGE' 'REPOSITORY' QNAME | 'ALL' 'IMAGE' 'REPOSITORIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('JOIN' 'POLICY' QNAME | 'ALL' 'JOIN' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' 'MASKING' 'POLICY' QNAME | (grantPrivilegeUser_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' ('MATERIALIZED' 'VIEW' QNAME | 'ALL' 'MATERIALIZED' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('MODEL' QNAME | 'ALL' 'MODELS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('PACKAGES' 'POLICY' QNAME | 'ALL' 'PACKAGES' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('PASSWORD' 'POLICY' QNAME | 'ALL' 'PASSWORD' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_11 | 'ALL' 'PRIVILEGES'?) 'ON' 'PIPE' QNAME | (grantPrivilegeUser_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('PRIVACY' 'POLICY' QNAME | 'ALL' 'PRIVACY' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('PROCEDURE' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' | 'ALL' 'PROCEDURES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('PROJECTION' 'POLICY' QNAME | 'ALL' 'PROJECTION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('ROW' 'ACCESS' 'POLICY' QNAME | 'ALL' 'ROW' 'ACCESS' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_5 | grantPrivilegeUser_schemaObjectPrivileges_13 | 'ALL' 'PRIVILEGES'?) 'ON' ('SECRET' QNAME | 'ALL' 'SECRETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_14 | 'ALL' 'PRIVILEGES'?) 'ON' ('SEMANTIC' 'VIEW' QNAME | 'ALL' 'SEMANTIC' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('SEQUENCE' QNAME | 'ALL' 'SEQUENCES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_15 | 'ALL' 'PRIVILEGES'?) 'ON' ('SERVICE' QNAME | 'ALL' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('SESSION' 'POLICY' QNAME | 'ALL' 'SESSION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('SNAPSHOT' QNAME | 'ALL' 'SNAPSHOTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_16 | grantPrivilegeUser_schemaObjectPrivileges_17 | 'ALL' 'PRIVILEGES'?) 'ON' ('STAGE' QNAME | 'ALL' 'STAGES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_18 | 'ALL' 'PRIVILEGES'?) 'ON' ('STREAM' QNAME | 'ALL' 'STREAMS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_19 | 'ALL' 'PRIVILEGES'?) 'ON' ('STREAMLIT' QNAME | 'ALL' 'STREAMLITS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_20 | 'ALL' 'PRIVILEGES'?) 'ON' ('TABLE' QNAME | 'ALL' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_12 | grantPrivilegeUser_schemaObjectPrivileges_21 | 'ALL' 'PRIVILEGES'?) 'ON' 'TAG' QNAME | (grantPrivilegeUser_schemaObjectPrivileges_22 | 'ALL' 'PRIVILEGES'?) 'ON' ('TASK' QNAME | 'ALL' 'TASKS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('TYPE' QNAME | 'ALL' 'TYPES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (grantPrivilegeUser_schemaObjectPrivileges_23 | 'ALL' 'PRIVILEGES'?) 'ON' ('VIEW' QNAME | 'ALL' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)))) 'TO' 'USER'? ID ('WITH' 'GRANT' 'OPTION')?
    ;
grantPrivilegeUser_globalPrivileges
    : ('ATTACH' 'POLICY' | 'AUDIT' | 'BIND' 'SERVICE' 'ENDPOINT' | 'APPLY' (('AGGREGATION' | 'AUTHENTICATION' | 'JOIN' | 'MASKING' | 'PACKAGES' | 'PASSWORD' | 'PROJECTION' | 'ROW' 'ACCESS' | 'SESSION') 'POLICY' | 'TAG') | 'EXECUTE' ('ALERT' | 'DATA' 'METRIC' 'FUNCTION' | 'MANAGED' 'ALERT' | 'MANAGED' 'TASK' | 'TASK') | 'IMPORT' 'SHARE' | 'MANAGE' ('ACCOUNT' 'SUPPORT' 'CASES' | 'EVENT' 'SHARING' | 'GRANTS' | 'LISTING' 'AUTO' 'FULFILLMENT' | 'ORGANIZATION' 'SUPPORT' 'CASES' | 'USER' 'SUPPORT' 'CASES' | 'WAREHOUSES') | 'MODIFY' ('LOG' 'LEVEL' | 'TRACE' 'LEVEL' | 'SESSION' 'LOG' 'LEVEL' | 'SESSION' 'TRACE' 'LEVEL') | 'MONITOR' ('EXECUTION' | 'SECURITY' | 'USAGE') | 'OVERRIDE' 'SHARE' 'RESTRICTIONS' | 'PURCHASE' 'DATA' 'EXCHANGE' 'LISTING' | 'RESOLVE' 'ALL' | 'READ' 'SESSION') (',' ('ATTACH' 'POLICY' | 'AUDIT' | 'BIND' 'SERVICE' 'ENDPOINT' | 'APPLY' (('AGGREGATION' | 'AUTHENTICATION' | 'JOIN' | 'MASKING' | 'PACKAGES' | 'PASSWORD' | 'PROJECTION' | 'ROW' 'ACCESS' | 'SESSION') 'POLICY' | 'TAG') | 'EXECUTE' ('ALERT' | 'DATA' 'METRIC' 'FUNCTION' | 'MANAGED' 'ALERT' | 'MANAGED' 'TASK' | 'TASK') | 'IMPORT' 'SHARE' | 'MANAGE' ('ACCOUNT' 'SUPPORT' 'CASES' | 'EVENT' 'SHARING' | 'GRANTS' | 'LISTING' 'AUTO' 'FULFILLMENT' | 'ORGANIZATION' 'SUPPORT' 'CASES' | 'USER' 'SUPPORT' 'CASES' | 'WAREHOUSES') | 'MODIFY' ('LOG' 'LEVEL' | 'TRACE' 'LEVEL' | 'SESSION' 'LOG' 'LEVEL' | 'SESSION' 'TRACE' 'LEVEL') | 'MONITOR' ('EXECUTION' | 'SECURITY' | 'USAGE') | 'OVERRIDE' 'SHARE' 'RESTRICTIONS' | 'PURCHASE' 'DATA' 'EXCHANGE' 'LISTING' | 'RESOLVE' 'ALL' | 'READ' 'SESSION'))*
    ;
grantPrivilegeUser_accountObjectPrivileges_1
    : ('MODIFY' | 'MONITOR' | 'OPERATE' | 'USAGE') (',' ('MODIFY' | 'MONITOR' | 'OPERATE' | 'USAGE'))*
    ;
grantPrivilegeUser_accountObjectPrivileges_2
    : 'FAILOVER' (',' 'FAILOVER')*
    ;
grantPrivilegeUser_accountObjectPrivileges_3
    : ('APPLYBUDGET' | 'IMPORTED' 'PRIVILEGES' | 'MODIFY' | 'MONITOR' | 'USAGE') (',' ('APPLYBUDGET' | 'IMPORTED' 'PRIVILEGES' | 'MODIFY' | 'MONITOR' | 'USAGE'))*
    ;
grantPrivilegeUser_accountObjectPrivileges_4
    : 'USAGE' (',' 'USAGE')*
    ;
grantPrivilegeUser_accountObjectPrivileges_5
    : ('FAILOVER' | 'MODIFY' | 'MONITOR' | 'REPLICATE') (',' ('FAILOVER' | 'MODIFY' | 'MONITOR' | 'REPLICATE'))*
    ;
grantPrivilegeUser_accountObjectPrivileges_6
    : ('USAGE' | 'USE_ANY_ROLE') (',' ('USAGE' | 'USE_ANY_ROLE'))*
    ;
grantPrivilegeUser_accountObjectPrivileges_7
    : ('MODIFY' | 'MONITOR' | 'REPLICATE') (',' ('MODIFY' | 'MONITOR' | 'REPLICATE'))*
    ;
grantPrivilegeUser_accountObjectPrivileges_8
    : ('MODIFY' | 'MONITOR') (',' ('MODIFY' | 'MONITOR'))*
    ;
grantPrivilegeUser_accountObjectPrivileges_9
    : 'MONITOR' (',' 'MONITOR')*
    ;
grantPrivilegeUser_accountObjectPrivileges_10
    : ('APPLYBUDGET' | 'MODIFY' | 'MONITOR' | 'USAGE' | 'OPERATE') (',' ('APPLYBUDGET' | 'MODIFY' | 'MONITOR' | 'USAGE' | 'OPERATE'))*
    ;
grantPrivilegeUser_schemaPrivileges
    : ('ADD' 'SEARCH' 'OPTIMIZATION' | 'APPLYBUDGET' | 'MODIFY' | 'MONITOR' | 'USAGE') (',' ('ADD' 'SEARCH' 'OPTIMIZATION' | 'APPLYBUDGET' | 'MODIFY' | 'MONITOR' | 'USAGE'))*
    ;
grantPrivilegeUser_schemaObjectPrivileges_1
    : ('MONITOR' | 'OPERATE') (',' ('MONITOR' | 'OPERATE'))*
    ;
grantPrivilegeUser_schemaObjectPrivileges_2
    : 'USAGE' (',' 'USAGE')*
    ;
grantPrivilegeUser_schemaObjectPrivileges_3
    : 'MONITOR' ',' 'OPERATE' ',' 'SELECT' (',' 'SELECT')*
    ;
grantPrivilegeUser_schemaObjectPrivileges_4
    : ('APPLYBUDGET' | 'DELETE' | 'REFERENCES' | 'SELECT' | 'TRUNCATE') (',' ('APPLYBUDGET' | 'DELETE' | 'REFERENCES' | 'SELECT' | 'TRUNCATE'))*
    ;
grantPrivilegeUser_schemaObjectPrivileges_5
    : 'USAGE' (',' 'USAGE')*
    ;
grantPrivilegeUser_schemaObjectPrivileges_6
    : 'READ' ',' 'WRITE' (',' 'READ' ',' 'WRITE')*
    ;
grantPrivilegeUser_schemaObjectPrivileges_7
    : ('APPLYBUDGET' | 'DELETE' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE') (',' ('APPLYBUDGET' | 'DELETE' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE'))*
    ;
grantPrivilegeUser_schemaObjectPrivileges_8
    : 'READ' ',' 'WRITE' (',' 'READ' ',' 'WRITE')*
    ;
grantPrivilegeUser_schemaObjectPrivileges_9
    : ('APPLYBUDGET' | 'DELETE' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE') (',' ('APPLYBUDGET' | 'DELETE' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE'))*
    ;
grantPrivilegeUser_schemaObjectPrivileges_10
    : ('APPLYBUDGET' | 'REFERENCES' | 'SELECT') (',' ('APPLYBUDGET' | 'REFERENCES' | 'SELECT'))*
    ;
grantPrivilegeUser_schemaObjectPrivileges_11
    : ('APPLYBUDGET' | 'MONITOR' | 'OPERATE') (',' ('APPLYBUDGET' | 'MONITOR' | 'OPERATE'))*
    ;
grantPrivilegeUser_schemaObjectPrivileges_12
    : 'APPLY' (',' 'APPLY')*
    ;
grantPrivilegeUser_schemaObjectPrivileges_13
    : ('READ' | 'USAGE') (',' ('READ' | 'USAGE'))*
    ;
grantPrivilegeUser_schemaObjectPrivileges_14
    : 'REFERENCES' (',' 'REFERENCES')*
    ;
grantPrivilegeUser_schemaObjectPrivileges_15
    : ('MONITOR' | 'OPERATE') (',' ('MONITOR' | 'OPERATE'))*
    ;
grantPrivilegeUser_schemaObjectPrivileges_16
    : 'USAGE' (',' 'USAGE')*
    ;
grantPrivilegeUser_schemaObjectPrivileges_17
    : 'READ' (',' 'WRITE')? (',' 'READ' (',' 'WRITE')?)*
    ;
grantPrivilegeUser_schemaObjectPrivileges_18
    : 'SELECT' (',' 'SELECT')*
    ;
grantPrivilegeUser_schemaObjectPrivileges_19
    : 'USAGE' (',' 'USAGE')*
    ;
grantPrivilegeUser_schemaObjectPrivileges_20
    : ('APPLYBUDGET' | 'DELETE' | 'EVOLVE' 'SCHEMA' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'SELECT' 'ERROR' 'TABLE' | 'TRUNCATE' | 'UPDATE') (',' ('APPLYBUDGET' | 'DELETE' | 'EVOLVE' 'SCHEMA' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'SELECT' 'ERROR' 'TABLE' | 'TRUNCATE' | 'UPDATE'))*
    ;
grantPrivilegeUser_schemaObjectPrivileges_21
    : 'READ'
    ;
grantPrivilegeUser_schemaObjectPrivileges_22
    : ('APPLYBUDGET' | 'MONITOR' | 'OPERATE') (',' ('APPLYBUDGET' | 'MONITOR' | 'OPERATE'))*
    ;
grantPrivilegeUser_schemaObjectPrivileges_23
    : ('REFERENCES' | 'SELECT') (',' ('REFERENCES' | 'SELECT'))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/grant-role
grantRole
    : 'GRANT' 'ROLE' ID 'TO' 'ROLE' ID ('WITH' 'GRANT' 'OPTION')?
    | 'GRANT' 'ROLE' ID 'TO' 'USER' ID (('FOR' 'INTERVAL' NUMBER | 'UNTIL' STRING) ('RESET' 'EXPIRATION')?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/grant-service-role
grantServiceRole
    : 'GRANT' 'SERVICE' 'ROLE' QNAME 'TO' ('ROLE' ID | 'APPLICATION' 'ROLE' ID | 'DATABASE' 'ROLE' ID)
    ;

// https://docs.snowflake.com/en/sql-reference/sql/insert
// not converted (form): ... stands for text the docs leave out
//   INSERT [ OVERWRITE ] INTO <target_table> [ ( <target_col_name> [ , ... ] ) ]
//          {
//              VALUES ( { <value> | DEFAULT | NULL } [ , ... ] ) [ , ( ... ) ]
//            | <query>
//          }

// https://docs.snowflake.com/en/sql-reference/sql/insert-multi-table
insertMultiTable
    : 'INSERT' 'OVERWRITE'? 'ALL' insertMultiTable_intoClause+ QUERY
    | 'INSERT' 'OVERWRITE'? ('FIRST' | 'ALL') ('WHEN' EXPR 'THEN' insertMultiTable_intoClause+)+ ('ELSE' insertMultiTable_intoClause)? QUERY
    ;
insertMultiTable_intoClause
    : 'INTO' QNAME ('(' ID (',' ID)* ')')? ('VALUES' '(' (ID | 'DEFAULT' | 'NULL') (',' (ID | 'DEFAULT' | 'NULL'))* ')')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/list
list
    : 'LIST' (list_internalStage | list_externalStage) ('PATTERN' '=' STRING)?
    | 'LIST' list_repositoryClone ('PATTERN' '=' STRING)?
    ;
list_internalStage
    : STAGE
    ;
list_externalStage
    : STAGE
    ;
list_repositoryClone
    : STAGE
    ;

// https://docs.snowflake.com/en/sql-reference/sql/merge
// fix: matchedClause: one action a clause; the clauses repeat in MERGE's own form, { matchedClause | notMatchedClause } [ ... ]
// fix: the source is a table or a subquery, as the parameters say
merge
    : 'MERGE' 'INTO' QNAME 'USING' (QNAME | '(' QUERY ')') 'ON' EXPR (merge_matchedClause | merge_notMatchedClause)+
    ;
merge_matchedClause
    : 'WHEN' 'MATCHED' ('AND' EXPR)? 'THEN' ('UPDATE' ('ALL' 'BY' 'NAME' | 'SET' ID '=' EXPR (',' ID '=' EXPR)*) | 'DELETE')
    ;
merge_notMatchedClause
    : 'WHEN' 'NOT' 'MATCHED' ('AND' EXPR)? 'THEN' 'INSERT' ('ALL' 'BY' 'NAME' | ('(' ID (',' ID)* ')')? 'VALUES' '(' EXPR (',' EXPR)* ')')
    ;

// https://docs.snowflake.com/en/sql-reference/sql/put
put
    : 'PUT' FILE put_internalStage ('PARALLEL' '=' NUMBER)? ('AUTO_COMPRESS' '=' 'TRUE' | 'AUTO_COMPRESS' '=' 'FALSE')? ('SOURCE_COMPRESSION' '=' 'AUTO_DETECT' | 'SOURCE_COMPRESSION' '=' 'GZIP' | 'SOURCE_COMPRESSION' '=' 'BZ2' | 'SOURCE_COMPRESSION' '=' 'BROTLI' | 'SOURCE_COMPRESSION' '=' 'ZSTD' | 'SOURCE_COMPRESSION' '=' 'DEFLATE' | 'SOURCE_COMPRESSION' '=' 'RAW_DEFLATE' | 'SOURCE_COMPRESSION' '=' 'NONE')? ('OVERWRITE' '=' 'TRUE' | 'OVERWRITE' '=' 'FALSE')?
    ;
put_internalStage
    : STAGE
    ;

// https://docs.snowflake.com/en/sql-reference/sql/remove
remove
    : 'REMOVE' (remove_internalStage | remove_externalStage) ('PATTERN' '=' STRING)?
    ;
remove_internalStage
    : STAGE
    ;
remove_externalStage
    : STAGE
    ;

// https://docs.snowflake.com/en/sql-reference/sql/revoke-application-role
revokeApplicationRole
    : 'REVOKE' ('GRANT' 'OPTION' 'FOR')? 'APPLICATION' 'ROLE' QNAME 'FROM' 'ROLE' ID ('RESTRICT' | 'CASCADE')?
    | 'REVOKE' 'APPLICATION' 'ROLE' QNAME 'FROM' ('APPLICATION' 'ROLE' ID | 'APPLICATION' ID | 'USER' ID)
    ;

// https://docs.snowflake.com/en/sql-reference/sql/revoke-caller
// fix: the privileges are each object type's, listed on Access control privileges rather than here; {object_privileges} leaves the forms unconverted
// not converted (form): undefined fragment "object_privileges": the page does not define it
//   REVOKE CALLER {object_privileges} [ , {object_privileges} ... ]
//     ON <object_type> <object_name>
//     FROM { ROLE | DATABASE ROLE } <grantee_name>
// not converted (form): undefined fragment "object_privileges": the page does not define it
//   REVOKE INHERITED CALLER {object_privileges} [ , {object_privileges} ... ]
//     ON ALL <object_type_plural>
//     IN { ACCOUNT | DATABASE <db_name> | SCHEMA <schema_name> | APPLICATION <app_name> | APPLICATION PACKAGE <app_pkg_name> }
//     FROM { ROLE | DATABASE ROLE } <grantee_name>
revokeCaller
    : 'REVOKE' 'ALL' 'CALLER' 'PRIVILEGES' 'ON' OBJECT_TYPE QNAME 'FROM' ('ROLE' | 'DATABASE' 'ROLE') ID
    | 'REVOKE' 'ALL' 'INHERITED' 'CALLER' 'PRIVILEGES' 'ON' 'ALL' OBJECT_TYPES 'IN' ('ACCOUNT' | 'DATABASE' ID | 'SCHEMA' QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID) 'FROM' ('ROLE' | 'DATABASE' 'ROLE') ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/revoke-database-role
revokeDatabaseRole
    : 'REVOKE' ('GRANT' 'OPTION' 'FOR')? 'DATABASE' 'ROLE' QNAME 'FROM' ('ROLE' | 'DATABASE' 'ROLE') ID ('RESTRICT' | 'CASCADE')?
    | 'REVOKE' 'DATABASE' 'ROLE' QNAME 'FROM' 'APPLICATION' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/revoke-database-role-share
revokeDatabaseRoleShare
    : 'REVOKE' 'DATABASE' 'ROLE' QNAME 'FROM' 'SHARE' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/revoke-privilege
// fix: FROM DATABASE ROLE form: the } after the database privileges ends the choice of grants one line early; dropped
// fix: FROM DATABASE ROLE form: the database privileges are one choice among the others; a | is missing before the next
// fix: the trailing [ , ... ] repeats the whole choice of schema privileges, as the usage notes allow several privileges on one object type
// fix: closes the { opened before ADD SEARCH OPTIMIZATION
// fix: read: each -- For X section of accountObjectPrivileges and schemaObjectPrivileges is the privileges on X alone, and the ON clause chooses an object type with its own privileges
revokePrivilege
    : 'REVOKE' ('GRANT' 'OPTION' 'FOR')? ((revokePrivilege_globalPrivileges | 'ALL' 'PRIVILEGES'?) 'ON' 'ACCOUNT' | ((revokePrivilege_accountObjectPrivileges_11 | 'ALL' 'PRIVILEGES'?) 'ON' 'RESOURCE' 'MONITOR' ID | (revokePrivilege_accountObjectPrivileges_13 | 'ALL' 'PRIVILEGES'?) 'ON' 'WAREHOUSE' ID | (revokePrivilege_accountObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' 'COMPUTE' 'POOL' ID | (revokePrivilege_accountObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' 'DATABASE' ID | (revokePrivilege_accountObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' 'DATA' 'EXCHANGE' 'LISTING' ID | (revokePrivilege_accountObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' 'INTEGRATION' ID | (revokePrivilege_accountObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' 'CONNECTION' ID | (revokePrivilege_accountObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' 'FAILOVER' 'GROUP' ID | (revokePrivilege_accountObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' 'REPLICATION' 'GROUP' ID | (revokePrivilege_accountObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' 'EXTERNAL' 'VOLUME' ID) | (revokePrivilege_schemaPrivileges | 'ALL' 'PRIVILEGES'?) 'ON' ('SCHEMA' QNAME | 'ALL' 'SCHEMAS' 'IN' 'DATABASE' ID) | (revokePrivilege_schemaPrivileges | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SCHEMAS' 'IN' 'DATABASE' ID | ((revokePrivilege_schemaObjectPrivileges_1 | 'ALL' 'PRIVILEGES'?) 'ON' ('AGENT' QNAME | 'ALL' 'AGENTS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('AGGREGATION' 'POLICY' QNAME | 'ALL' 'AGGREGATION' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' ('ALERT' QNAME | 'ALL' 'ALERTS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' ('APPLICATION' 'SERVICE' QNAME | 'ALL' 'APPLICATION' 'SERVICES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('ARTIFACT' 'REPOSITORY' QNAME | 'ALL' 'ARTIFACT' 'REPOSITORIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('AUTHENTICATION' 'POLICY' QNAME | 'ALL' 'AUTHENTICATION' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('CONTACT' QNAME | 'ALL' 'CONTACTS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' ('CORTEX' 'EXTENSION' QNAME | 'ALL' 'CORTEX' 'EXTENSIONS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' ('CORTEX' 'SEARCH' 'SERVICE' QNAME | 'ALL' 'CORTEX' 'SEARCH' 'SERVICES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('DATA' 'METRIC' 'FUNCTION' QNAME | 'ALL' 'DATA' 'METRIC' 'FUNCTIONS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('DATASET' QNAME | 'ALL' 'DATASETS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_11 | 'ALL' 'PRIVILEGES'?) 'ON' ('DBT' 'PROJECT' QNAME | 'ALL' 'DBT' 'PROJECTS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('DCM' 'PROJECT' QNAME | 'ALL' 'DCM' 'PROJECTS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_13 | 'ALL' 'PRIVILEGES'?) 'ON' ('DYNAMIC' 'TABLE' QNAME | 'ALL' 'DYNAMIC' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_15 | 'ALL' 'PRIVILEGES'?) 'ON' ('EVENT' 'TABLE' QNAME | 'ALL' 'EVENT' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_14 | 'ALL' 'PRIVILEGES'?) 'ON' ('EXPERIMENT' QNAME | 'ALL' 'EXPERIMENTS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('FILE' 'FORMAT' QNAME | 'ALL' 'FILE' 'FORMATS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('FUNCTION' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' | 'ALL' 'FUNCTIONS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_16 | 'ALL' 'PRIVILEGES'?) 'ON' ('GATEWAY' QNAME | 'ALL' 'GATEWAYS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_17 | 'ALL' 'PRIVILEGES'?) 'ON' ('GIT' 'REPOSITORY' QNAME | 'ALL' 'GIT' 'REPOSITORIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_18 | 'ALL' 'PRIVILEGES'?) 'ON' ('HYBRID' 'TABLE' QNAME | 'ALL' 'HYBRID' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_20 | 'ALL' 'PRIVILEGES'?) 'ON' ('ICEBERG' 'TABLE' QNAME | 'ALL' 'ICEBERG' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_19 | 'ALL' 'PRIVILEGES'?) 'ON' ('IMAGE' 'REPOSITORY' QNAME | 'ALL' 'IMAGE' 'REPOSITORIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_21 | 'ALL' 'PRIVILEGES'?) 'ON' ('INTERACTIVE' 'TABLE' QNAME | 'ALL' 'INTERACTIVE' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('JOIN' 'POLICY' QNAME | 'ALL' 'JOIN' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' 'MAINTENANCE' 'POLICY' QNAME | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' 'MASKING' 'POLICY' QNAME | (revokePrivilege_schemaObjectPrivileges_22 | 'ALL' 'PRIVILEGES'?) 'ON' ('MATERIALIZED' 'VIEW' QNAME | 'ALL' 'MATERIALIZED' 'VIEWS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_23 | 'ALL' 'PRIVILEGES'?) 'ON' ('MCP' 'SERVER' QNAME | 'ALL' 'MCP' 'SERVERS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('MODEL' QNAME | 'ALL' 'MODELS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_24 | 'ALL' 'PRIVILEGES'?) 'ON' ('ONLINE' 'FEATURE' 'TABLE' QNAME | 'ALL' 'ONLINE' 'FEATURE' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('PACKAGES' 'POLICY' QNAME | 'ALL' 'PACKAGES' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('PASSWORD' 'POLICY' QNAME | 'ALL' 'PASSWORD' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_25 | 'ALL' 'PRIVILEGES'?) 'ON' 'PIPE' QNAME | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('PRIVACY' 'POLICY' QNAME | 'ALL' 'PRIVACY' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('PROCEDURE' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' | 'ALL' 'PROCEDURES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('PROJECTION' 'POLICY' QNAME | 'ALL' 'PROJECTION' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('ROW' 'ACCESS' 'POLICY' QNAME | 'ALL' 'ROW' 'ACCESS' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | revokePrivilege_schemaObjectPrivileges_27 | 'ALL' 'PRIVILEGES'?) 'ON' ('SECRET' QNAME | 'ALL' 'SECRETS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_28 | 'ALL' 'PRIVILEGES'?) 'ON' ('SEMANTIC' 'VIEW' QNAME | 'ALL' 'SEMANTIC' 'VIEWS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('SEQUENCE' QNAME | 'ALL' 'SEQUENCES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_29 | 'ALL' 'PRIVILEGES'?) 'ON' ('SERVICE' QNAME | 'ALL' 'SERVICES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('SESSION' 'POLICY' QNAME | 'ALL' 'SESSION' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('SNAPSHOT' QNAME | 'ALL' 'SNAPSHOTS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' ('SNAPSHOT' 'POLICY' QNAME | 'ALL' 'SNAPSHOT' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' ('SNAPSHOT' 'SET' QNAME | 'ALL' 'SNAPSHOT' 'SETS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_30 | revokePrivilege_schemaObjectPrivileges_31 | 'ALL' 'PRIVILEGES'?) 'ON' ('STAGE' QNAME | 'ALL' 'STAGES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('STORAGE' 'LIFECYCLE' 'POLICY' QNAME | 'ALL' 'STORAGE' 'LIFECYCLE' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_32 | 'ALL' 'PRIVILEGES'?) 'ON' ('STREAM' QNAME | 'ALL' 'STREAMS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_33 | 'ALL' 'PRIVILEGES'?) 'ON' ('STREAMLIT' QNAME | 'ALL' 'STREAMLITS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_34 | 'ALL' 'PRIVILEGES'?) 'ON' ('TABLE' QNAME | 'ALL' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | revokePrivilege_schemaObjectPrivileges_35 | 'ALL' 'PRIVILEGES'?) 'ON' 'TAG' QNAME | (revokePrivilege_schemaObjectPrivileges_36 | 'ALL' 'PRIVILEGES'?) 'ON' ('TASK' QNAME | 'ALL' 'TASKS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('TYPE' QNAME | 'ALL' 'TYPES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_37 | 'ALL' 'PRIVILEGES'?) 'ON' ('VIEW' QNAME | 'ALL' 'VIEWS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_38 | 'ALL' 'PRIVILEGES'?) 'ON' ('WORKSPACE' QNAME | 'ALL' 'WORKSPACES' 'IN' 'SCHEMA' QNAME)) | ((revokePrivilege_schemaObjectPrivileges_1 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'AGENTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'ALERTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'APPLICATION' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'ARTIFACT' 'REPOSITORIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'AUTHENTICATION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'CONTACTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'CORTEX' 'EXTENSIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'CORTEX' 'SEARCH' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DATA' 'METRIC' 'FUNCTIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DATASETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_11 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DBT' 'PROJECTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DCM' 'PROJECTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_13 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DYNAMIC' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_15 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'EVENT' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_14 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'EXPERIMENTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'FILE' 'FORMATS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'FUNCTIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_16 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'GATEWAYS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_17 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'GIT' 'REPOSITORIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_18 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'HYBRID' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_20 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'ICEBERG' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_21 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'INTERACTIVE' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_22 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'MATERIALIZED' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_23 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'MCP' 'SERVERS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'MODELS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_24 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'ONLINE' 'FEATURE' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PASSWORD' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_25 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PIPES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PRIVACY' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PROCEDURES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | revokePrivilege_schemaObjectPrivileges_27 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SECRETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_28 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SEMANTIC' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SEQUENCES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_29 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SNAPSHOT' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SNAPSHOT' 'SETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_30 | revokePrivilege_schemaObjectPrivileges_31 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'STAGES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_32 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'STREAMS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_33 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'STREAMLITS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_34 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_36 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'TASKS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'TYPES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_37 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_38 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'WORKSPACES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME))) 'FROM' 'ROLE'? ID ('RESTRICT' | 'CASCADE')?
    | 'REVOKE' ('GRANT' 'OPTION' 'FOR')? (('CREATE' 'SCHEMA' | 'MODIFY' | 'MONITOR' | 'USAGE') (',' ('CREATE' 'SCHEMA' | 'MODIFY' | 'MONITOR' | 'USAGE'))* 'ON' 'DATABASE' QNAME | (revokePrivilege_globalPrivileges | 'ALL' 'PRIVILEGES'?) 'ON' 'ACCOUNT' | ((revokePrivilege_accountObjectPrivileges_11 | 'ALL' 'PRIVILEGES'?) 'ON' 'RESOURCE' 'MONITOR' ID | (revokePrivilege_accountObjectPrivileges_13 | 'ALL' 'PRIVILEGES'?) 'ON' 'WAREHOUSE' ID | (revokePrivilege_accountObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' 'COMPUTE' 'POOL' ID | (revokePrivilege_accountObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' 'DATABASE' ID | (revokePrivilege_accountObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' 'INTEGRATION' ID | (revokePrivilege_accountObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' 'EXTERNAL' 'VOLUME' ID) | (revokePrivilege_schemaPrivileges | 'ALL' 'PRIVILEGES'?) 'ON' ('SCHEMA' QNAME | 'ALL' 'SCHEMAS' 'IN' 'DATABASE' ID) | (revokePrivilege_schemaPrivileges | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SCHEMAS' 'IN' 'DATABASE' ID | ((revokePrivilege_schemaObjectPrivileges_1 | 'ALL' 'PRIVILEGES'?) 'ON' ('AGENT' QNAME | 'ALL' 'AGENTS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('AGGREGATION' 'POLICY' QNAME | 'ALL' 'AGGREGATION' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' ('ALERT' QNAME | 'ALL' 'ALERTS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' ('APPLICATION' 'SERVICE' QNAME | 'ALL' 'APPLICATION' 'SERVICES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('ARTIFACT' 'REPOSITORY' QNAME | 'ALL' 'ARTIFACT' 'REPOSITORIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('AUTHENTICATION' 'POLICY' QNAME | 'ALL' 'AUTHENTICATION' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('CONTACT' QNAME | 'ALL' 'CONTACTS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' ('CORTEX' 'EXTENSION' QNAME | 'ALL' 'CORTEX' 'EXTENSIONS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' ('CORTEX' 'SEARCH' 'SERVICE' QNAME | 'ALL' 'CORTEX' 'SEARCH' 'SERVICES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('DATA' 'METRIC' 'FUNCTION' QNAME | 'ALL' 'DATA' 'METRIC' 'FUNCTIONS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('DATASET' QNAME | 'ALL' 'DATASETS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_11 | 'ALL' 'PRIVILEGES'?) 'ON' ('DBT' 'PROJECT' QNAME | 'ALL' 'DBT' 'PROJECTS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('DCM' 'PROJECT' QNAME | 'ALL' 'DCM' 'PROJECTS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_13 | 'ALL' 'PRIVILEGES'?) 'ON' ('DYNAMIC' 'TABLE' QNAME | 'ALL' 'DYNAMIC' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_15 | 'ALL' 'PRIVILEGES'?) 'ON' ('EVENT' 'TABLE' QNAME | 'ALL' 'EVENT' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_14 | 'ALL' 'PRIVILEGES'?) 'ON' ('EXPERIMENT' QNAME | 'ALL' 'EXPERIMENTS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('FILE' 'FORMAT' QNAME | 'ALL' 'FILE' 'FORMATS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('FUNCTION' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' | 'ALL' 'FUNCTIONS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_16 | 'ALL' 'PRIVILEGES'?) 'ON' ('GATEWAY' QNAME | 'ALL' 'GATEWAYS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_17 | 'ALL' 'PRIVILEGES'?) 'ON' ('GIT' 'REPOSITORY' QNAME | 'ALL' 'GIT' 'REPOSITORIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_18 | 'ALL' 'PRIVILEGES'?) 'ON' ('HYBRID' 'TABLE' QNAME | 'ALL' 'HYBRID' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_20 | 'ALL' 'PRIVILEGES'?) 'ON' ('ICEBERG' 'TABLE' QNAME | 'ALL' 'ICEBERG' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_19 | 'ALL' 'PRIVILEGES'?) 'ON' ('IMAGE' 'REPOSITORY' QNAME | 'ALL' 'IMAGE' 'REPOSITORIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_21 | 'ALL' 'PRIVILEGES'?) 'ON' ('INTERACTIVE' 'TABLE' QNAME | 'ALL' 'INTERACTIVE' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('JOIN' 'POLICY' QNAME | 'ALL' 'JOIN' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' 'MAINTENANCE' 'POLICY' QNAME | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' 'MASKING' 'POLICY' QNAME | (revokePrivilege_schemaObjectPrivileges_22 | 'ALL' 'PRIVILEGES'?) 'ON' ('MATERIALIZED' 'VIEW' QNAME | 'ALL' 'MATERIALIZED' 'VIEWS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_23 | 'ALL' 'PRIVILEGES'?) 'ON' ('MCP' 'SERVER' QNAME | 'ALL' 'MCP' 'SERVERS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('MODEL' QNAME | 'ALL' 'MODELS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_24 | 'ALL' 'PRIVILEGES'?) 'ON' ('ONLINE' 'FEATURE' 'TABLE' QNAME | 'ALL' 'ONLINE' 'FEATURE' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('PACKAGES' 'POLICY' QNAME | 'ALL' 'PACKAGES' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('PASSWORD' 'POLICY' QNAME | 'ALL' 'PASSWORD' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_25 | 'ALL' 'PRIVILEGES'?) 'ON' 'PIPE' QNAME | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('PRIVACY' 'POLICY' QNAME | 'ALL' 'PRIVACY' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('PROCEDURE' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' | 'ALL' 'PROCEDURES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('PROJECTION' 'POLICY' QNAME | 'ALL' 'PROJECTION' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('ROW' 'ACCESS' 'POLICY' QNAME | 'ALL' 'ROW' 'ACCESS' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | revokePrivilege_schemaObjectPrivileges_27 | 'ALL' 'PRIVILEGES'?) 'ON' ('SECRET' QNAME | 'ALL' 'SECRETS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_28 | 'ALL' 'PRIVILEGES'?) 'ON' ('SEMANTIC' 'VIEW' QNAME | 'ALL' 'SEMANTIC' 'VIEWS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('SEQUENCE' QNAME | 'ALL' 'SEQUENCES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_29 | 'ALL' 'PRIVILEGES'?) 'ON' ('SERVICE' QNAME | 'ALL' 'SERVICES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('SESSION' 'POLICY' QNAME | 'ALL' 'SESSION' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('SNAPSHOT' QNAME | 'ALL' 'SNAPSHOTS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' ('SNAPSHOT' 'POLICY' QNAME | 'ALL' 'SNAPSHOT' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' ('SNAPSHOT' 'SET' QNAME | 'ALL' 'SNAPSHOT' 'SETS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_30 | revokePrivilege_schemaObjectPrivileges_31 | 'ALL' 'PRIVILEGES'?) 'ON' ('STAGE' QNAME | 'ALL' 'STAGES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' ('STORAGE' 'LIFECYCLE' 'POLICY' QNAME | 'ALL' 'STORAGE' 'LIFECYCLE' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_32 | 'ALL' 'PRIVILEGES'?) 'ON' ('STREAM' QNAME | 'ALL' 'STREAMS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_33 | 'ALL' 'PRIVILEGES'?) 'ON' ('STREAMLIT' QNAME | 'ALL' 'STREAMLITS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_34 | 'ALL' 'PRIVILEGES'?) 'ON' ('TABLE' QNAME | 'ALL' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | revokePrivilege_schemaObjectPrivileges_35 | 'ALL' 'PRIVILEGES'?) 'ON' 'TAG' QNAME | (revokePrivilege_schemaObjectPrivileges_36 | 'ALL' 'PRIVILEGES'?) 'ON' ('TASK' QNAME | 'ALL' 'TASKS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('TYPE' QNAME | 'ALL' 'TYPES' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_37 | 'ALL' 'PRIVILEGES'?) 'ON' ('VIEW' QNAME | 'ALL' 'VIEWS' 'IN' 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_38 | 'ALL' 'PRIVILEGES'?) 'ON' ('WORKSPACE' QNAME | 'ALL' 'WORKSPACES' 'IN' 'SCHEMA' QNAME)) | ((revokePrivilege_schemaObjectPrivileges_1 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'AGENTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'ALERTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'APPLICATION' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'ARTIFACT' 'REPOSITORIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'AUTHENTICATION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'CONTACTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'CORTEX' 'EXTENSIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'CORTEX' 'SEARCH' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DATA' 'METRIC' 'FUNCTIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DATASETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_11 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DBT' 'PROJECTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DCM' 'PROJECTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_13 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DYNAMIC' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_15 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'EVENT' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_14 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'EXPERIMENTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'FILE' 'FORMATS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'FUNCTIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_16 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'GATEWAYS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_17 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'GIT' 'REPOSITORIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_18 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'HYBRID' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_20 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'ICEBERG' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_21 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'INTERACTIVE' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_22 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'MATERIALIZED' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_23 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'MCP' 'SERVERS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'MODELS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_24 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'ONLINE' 'FEATURE' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PASSWORD' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_25 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PIPES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_26 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PRIVACY' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PROCEDURES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | revokePrivilege_schemaObjectPrivileges_27 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SECRETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_28 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SEMANTIC' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SEQUENCES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_29 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SNAPSHOT' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SNAPSHOT' 'SETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_30 | revokePrivilege_schemaObjectPrivileges_31 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'STAGES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_32 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'STREAMS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_33 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'STREAMLITS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_34 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_36 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'TASKS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'TYPES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_37 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilege_schemaObjectPrivileges_38 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'WORKSPACES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME))) 'FROM' 'DATABASE' 'ROLE' ID ('RESTRICT' | 'CASCADE')?
    ;
revokePrivilege_globalPrivileges
    : ('CREATE' ('ACCOUNT' | 'APPLICATION' | 'APPLICATION' 'PACKAGE' | 'COMPUTE' 'POOL' | 'LISTING' | 'DATABASE' | 'EXTERNAL' 'VOLUME' | 'FAILOVER' 'GROUP' | 'INTEGRATION' | 'NETWORK' 'POLICY' | 'ORGANIZATION' 'LISTING' | 'ORGANIZATION' 'PROFILE' | 'PREVIEW' 'APPLICATION' | 'REPLICATION' 'GROUP' | 'ROLE' | 'SHARE' | 'USER' | 'WAREHOUSE') | 'ATTACH' 'POLICY' | 'AUDIT' | 'BIND' 'SERVICE' 'ENDPOINT' | 'APPLY' (('AGGREGATION' | 'AUTHENTICATION' | 'JOIN' | 'MAINTENANCE' | 'MASKING' | 'PACKAGES' | 'PASSWORD' | 'PROJECTION' | 'ROW' 'ACCESS' | 'SESSION' | 'STORAGE' 'LIFECYCLE') 'POLICY' | 'CONTACT' | 'TAG') | 'EXECUTE' ('ALERT' | 'DATA' 'METRIC' 'FUNCTION' | 'MANAGED' 'ALERT' | 'MANAGED' 'TASK' | 'TASK') | 'IMPORT' ('SHARE' | 'ORGANIZATION' 'LISTING') | 'MANAGE' ('ACCOUNT' 'SUPPORT' 'CASES' | 'APPLICATION' 'SPECIFICATIONS' | 'CORTEX' 'EXTENSION' 'SECURITY' 'SCAN' | 'EVENT' 'SHARING' | 'GRANTS' | 'LISTING' 'AUTO' 'FULFILLMENT' | 'ORGANIZATION' 'SUPPORT' 'CASES' | 'SHARE' 'TARGET' | 'USER' 'SUPPORT' 'CASES' | 'VISIBILITY' | 'WAREHOUSES') | 'MODIFY' ('LOG' 'LEVEL' | 'TRACE' 'LEVEL' | 'SESSION' 'LOG' 'LEVEL' | 'SESSION' 'TRACE' 'LEVEL') | 'MONITOR' ('CORTEX' 'EXTENSIONS' | 'EXECUTION' | 'SECURITY' | 'USAGE') | 'OVERRIDE' 'SHARE' 'RESTRICTIONS' | 'PURCHASE' 'DATA' 'EXCHANGE' 'LISTING' | 'RESOLVE' 'ALL' | 'READ' 'SESSION' | 'READ' 'UNREDACTED' 'ERROR' 'TABLE' | 'USE' 'AI' 'FUNCTION' QNAME | 'USE' 'AI' 'FUNCTIONS') (',' ('CREATE' ('ACCOUNT' | 'APPLICATION' | 'APPLICATION' 'PACKAGE' | 'COMPUTE' 'POOL' | 'LISTING' | 'DATABASE' | 'EXTERNAL' 'VOLUME' | 'FAILOVER' 'GROUP' | 'INTEGRATION' | 'NETWORK' 'POLICY' | 'ORGANIZATION' 'LISTING' | 'ORGANIZATION' 'PROFILE' | 'PREVIEW' 'APPLICATION' | 'REPLICATION' 'GROUP' | 'ROLE' | 'SHARE' | 'USER' | 'WAREHOUSE') | 'ATTACH' 'POLICY' | 'AUDIT' | 'BIND' 'SERVICE' 'ENDPOINT' | 'APPLY' (('AGGREGATION' | 'AUTHENTICATION' | 'JOIN' | 'MAINTENANCE' | 'MASKING' | 'PACKAGES' | 'PASSWORD' | 'PROJECTION' | 'ROW' 'ACCESS' | 'SESSION' | 'STORAGE' 'LIFECYCLE') 'POLICY' | 'CONTACT' | 'TAG') | 'EXECUTE' ('ALERT' | 'DATA' 'METRIC' 'FUNCTION' | 'MANAGED' 'ALERT' | 'MANAGED' 'TASK' | 'TASK') | 'IMPORT' ('SHARE' | 'ORGANIZATION' 'LISTING') | 'MANAGE' ('ACCOUNT' 'SUPPORT' 'CASES' | 'APPLICATION' 'SPECIFICATIONS' | 'CORTEX' 'EXTENSION' 'SECURITY' 'SCAN' | 'EVENT' 'SHARING' | 'GRANTS' | 'LISTING' 'AUTO' 'FULFILLMENT' | 'ORGANIZATION' 'SUPPORT' 'CASES' | 'SHARE' 'TARGET' | 'USER' 'SUPPORT' 'CASES' | 'VISIBILITY' | 'WAREHOUSES') | 'MODIFY' ('LOG' 'LEVEL' | 'TRACE' 'LEVEL' | 'SESSION' 'LOG' 'LEVEL' | 'SESSION' 'TRACE' 'LEVEL') | 'MONITOR' ('CORTEX' 'EXTENSIONS' | 'EXECUTION' | 'SECURITY' | 'USAGE') | 'OVERRIDE' 'SHARE' 'RESTRICTIONS' | 'PURCHASE' 'DATA' 'EXCHANGE' 'LISTING' | 'RESOLVE' 'ALL' | 'READ' 'SESSION' | 'READ' 'UNREDACTED' 'ERROR' 'TABLE' | 'USE' 'AI' 'FUNCTION' QNAME | 'USE' 'AI' 'FUNCTIONS'))*
    ;
revokePrivilege_accountObjectPrivileges_2
    : ('MODIFY' | 'MONITOR' | 'OPERATE' | 'USAGE') (',' ('MODIFY' | 'MONITOR' | 'OPERATE' | 'USAGE'))*
    ;
revokePrivilege_accountObjectPrivileges_3
    : 'FAILOVER' (',' 'FAILOVER')*
    ;
revokePrivilege_accountObjectPrivileges_4
    : ('MODIFY' | 'USAGE') (',' ('MODIFY' | 'USAGE'))*
    ;
revokePrivilege_accountObjectPrivileges_5
    : ('APPLYBUDGET' | 'CREATE' ('DATABASE' 'ROLE' | 'SCHEMA') | 'IMPORTED' 'PRIVILEGES' | 'MODIFY' | 'MONITOR' | 'USAGE') (',' ('APPLYBUDGET' | 'CREATE' ('DATABASE' 'ROLE' | 'SCHEMA') | 'IMPORTED' 'PRIVILEGES' | 'MODIFY' | 'MONITOR' | 'USAGE'))*
    ;
revokePrivilege_accountObjectPrivileges_6
    : 'USAGE' (',' 'USAGE')*
    ;
revokePrivilege_accountObjectPrivileges_7
    : ('FAILOVER' | 'MODIFY' | 'MONITOR' | 'REPLICATE') (',' ('FAILOVER' | 'MODIFY' | 'MONITOR' | 'REPLICATE'))*
    ;
revokePrivilege_accountObjectPrivileges_8
    : ('USAGE' | 'USE_ANY_ROLE') (',' ('USAGE' | 'USE_ANY_ROLE'))*
    ;
revokePrivilege_accountObjectPrivileges_10
    : ('MODIFY' | 'MONITOR' | 'REPLICATE') (',' ('MODIFY' | 'MONITOR' | 'REPLICATE'))*
    ;
revokePrivilege_accountObjectPrivileges_11
    : ('MODIFY' | 'MONITOR') (',' ('MODIFY' | 'MONITOR'))*
    ;
revokePrivilege_accountObjectPrivileges_13
    : ('APPLYBUDGET' | 'MODIFY' | 'MONITOR' | 'USAGE' | 'OPERATE') (',' ('APPLYBUDGET' | 'MODIFY' | 'MONITOR' | 'USAGE' | 'OPERATE'))*
    ;
revokePrivilege_schemaPrivileges
    : ('ADD' 'SEARCH' 'OPTIMIZATION' | 'APPLYBUDGET' | 'CREATE' ('AGENT' | 'ALERT' | 'APPLICATION' 'SERVICE' | 'ARTIFACT' 'REPOSITORY' | 'CONTACT' | 'CORTEX' 'SEARCH' 'SERVICE' | 'DATA' 'METRIC' 'FUNCTION' | 'DATASET' | 'DBT' 'PROJECT' | 'DCM' 'PROJECT' | 'EVENT' 'TABLE' | 'EXPERIMENT' | 'FILE' 'FORMAT' | 'FUNCTION' | 'GATEWAY' | ('GIT' | 'IMAGE') 'REPOSITORY' | 'MCP' 'SERVER' | 'MODEL' | 'NETWORK' 'RULE' | 'NOTEBOOK' | 'PIPE' | 'PROCEDURE' | ('AGGREGATION' | 'AUTHENTICATION' | 'MAINTENANCE' | 'MASKING' | 'PACKAGES' | 'PASSWORD' | 'PRIVACY' | 'PROJECTION' | 'ROW' 'ACCESS' | 'SESSION' | 'STORAGE' 'LIFECYCLE') 'POLICY' | 'SECRET' | 'SEQUENCE' | 'SERVICE' | 'SNAPSHOT' | 'SNAPSHOT' 'POLICY' | 'SNAPSHOT' 'SET' | 'STAGE' | 'STREAM' | 'STREAMLIT' | 'SNOWFLAKE' '.' 'CORE' '.' 'BUDGET' | 'SNOWFLAKE' '.' 'DATA_PRIVACY' '.' 'CLASSIFICATION_PROFILE' | 'SNOWFLAKE' '.' 'DATA_PRIVACY' '.' 'CUSTOM_CLASSIFIER' | 'SNOWFLAKE' '.' 'ML' '.' 'ANOMALY_DETECTION' | 'SNOWFLAKE' '.' 'ML' '.' 'CLASSIFICATION' | 'SNOWFLAKE' '.' 'ML' '.' 'FORECAST' | 'SNOWFLAKE' '.' 'ML' '.' 'TOP_INSIGHTS' | 'SNOWFLAKE' '.' 'ML' '.' 'DOCUMENT_INTELLIGENCE' | (('DYNAMIC' | 'EXTERNAL' | 'HYBRID' | 'ICEBERG' | 'INTERACTIVE' | 'ONLINE' 'FEATURE'))? 'TABLE' | 'TAG' | 'TASK' | 'TYPE' | 'WORKSPACE' | (('MATERIALIZED' | 'SEMANTIC'))? 'VIEW') | 'MODIFY' | 'MONITOR' | 'USAGE') (',' ('ADD' 'SEARCH' 'OPTIMIZATION' | 'APPLYBUDGET' | 'CREATE' ('AGENT' | 'ALERT' | 'APPLICATION' 'SERVICE' | 'ARTIFACT' 'REPOSITORY' | 'CONTACT' | 'CORTEX' 'SEARCH' 'SERVICE' | 'DATA' 'METRIC' 'FUNCTION' | 'DATASET' | 'DBT' 'PROJECT' | 'DCM' 'PROJECT' | 'EVENT' 'TABLE' | 'EXPERIMENT' | 'FILE' 'FORMAT' | 'FUNCTION' | 'GATEWAY' | ('GIT' | 'IMAGE') 'REPOSITORY' | 'MCP' 'SERVER' | 'MODEL' | 'NETWORK' 'RULE' | 'NOTEBOOK' | 'PIPE' | 'PROCEDURE' | ('AGGREGATION' | 'AUTHENTICATION' | 'MAINTENANCE' | 'MASKING' | 'PACKAGES' | 'PASSWORD' | 'PRIVACY' | 'PROJECTION' | 'ROW' 'ACCESS' | 'SESSION' | 'STORAGE' 'LIFECYCLE') 'POLICY' | 'SECRET' | 'SEQUENCE' | 'SERVICE' | 'SNAPSHOT' | 'SNAPSHOT' 'POLICY' | 'SNAPSHOT' 'SET' | 'STAGE' | 'STREAM' | 'STREAMLIT' | 'SNOWFLAKE' '.' 'CORE' '.' 'BUDGET' | 'SNOWFLAKE' '.' 'DATA_PRIVACY' '.' 'CLASSIFICATION_PROFILE' | 'SNOWFLAKE' '.' 'DATA_PRIVACY' '.' 'CUSTOM_CLASSIFIER' | 'SNOWFLAKE' '.' 'ML' '.' 'ANOMALY_DETECTION' | 'SNOWFLAKE' '.' 'ML' '.' 'CLASSIFICATION' | 'SNOWFLAKE' '.' 'ML' '.' 'FORECAST' | 'SNOWFLAKE' '.' 'ML' '.' 'TOP_INSIGHTS' | 'SNOWFLAKE' '.' 'ML' '.' 'DOCUMENT_INTELLIGENCE' | (('DYNAMIC' | 'EXTERNAL' | 'HYBRID' | 'ICEBERG' | 'INTERACTIVE' | 'ONLINE' 'FEATURE'))? 'TABLE' | 'TAG' | 'TASK' | 'TYPE' | 'WORKSPACE' | (('MATERIALIZED' | 'SEMANTIC'))? 'VIEW') | 'MODIFY' | 'MONITOR' | 'USAGE'))*
    ;
revokePrivilege_schemaObjectPrivileges_1
    : ('MODIFY' | 'MONITOR' | 'USAGE') (',' ('MODIFY' | 'MONITOR' | 'USAGE'))*
    ;
revokePrivilege_schemaObjectPrivileges_2
    : ('MONITOR' | 'OPERATE') (',' ('MONITOR' | 'OPERATE'))*
    ;
revokePrivilege_schemaObjectPrivileges_3
    : ('MONITOR' | 'OPERATE' | 'USAGE') (',' ('MONITOR' | 'OPERATE' | 'USAGE'))*
    ;
revokePrivilege_schemaObjectPrivileges_4
    : 'READ' (',' 'READ')*
    ;
revokePrivilege_schemaObjectPrivileges_5
    : ('APPLY' | 'MODIFY') (',' ('APPLY' | 'MODIFY'))*
    ;
revokePrivilege_schemaObjectPrivileges_6
    : ('READ' | 'WRITE') (',' ('READ' | 'WRITE'))*
    ;
revokePrivilege_schemaObjectPrivileges_7
    : ('MONITOR' | 'OPERATE' | 'USAGE') (',' ('MONITOR' | 'OPERATE' | 'USAGE'))*
    ;
revokePrivilege_schemaObjectPrivileges_8
    : 'USAGE' (',' 'USAGE')*
    ;
revokePrivilege_schemaObjectPrivileges_9
    : 'USAGE' (',' 'USAGE')*
    ;
revokePrivilege_schemaObjectPrivileges_10
    : 'USAGE' (',' 'USAGE')*
    ;
revokePrivilege_schemaObjectPrivileges_11
    : 'USAGE' ',' 'MONITOR' (',' 'MONITOR')*
    ;
revokePrivilege_schemaObjectPrivileges_12
    : 'READ' ',' 'MONITOR' (',' 'MONITOR')*
    ;
revokePrivilege_schemaObjectPrivileges_13
    : 'MONITOR' ',' 'OPERATE' ',' 'SELECT' (',' 'SELECT')*
    ;
revokePrivilege_schemaObjectPrivileges_14
    : ('CREATE' | 'MODIFY' | 'USAGE') (',' ('CREATE' | 'MODIFY' | 'USAGE'))*
    ;
revokePrivilege_schemaObjectPrivileges_15
    : ('APPLYBUDGET' | 'DELETE' | 'OWNERSHIP' | 'REFERENCES' | 'SELECT' | 'TRUNCATE') (',' ('APPLYBUDGET' | 'DELETE' | 'OWNERSHIP' | 'REFERENCES' | 'SELECT' | 'TRUNCATE'))*
    ;
revokePrivilege_schemaObjectPrivileges_16
    : ('CREATE' | 'MODIFY' | 'USAGE') (',' ('CREATE' | 'MODIFY' | 'USAGE'))*
    ;
revokePrivilege_schemaObjectPrivileges_17
    : 'READ' ',' 'WRITE' (',' 'READ' ',' 'WRITE')*
    ;
revokePrivilege_schemaObjectPrivileges_18
    : ('APPLYBUDGET' | 'DELETE' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE') (',' ('APPLYBUDGET' | 'DELETE' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE'))*
    ;
revokePrivilege_schemaObjectPrivileges_19
    : 'READ' ',' 'WRITE' (',' 'READ' ',' 'WRITE')*
    ;
revokePrivilege_schemaObjectPrivileges_20
    : ('APPLYBUDGET' | 'DELETE' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE') (',' ('APPLYBUDGET' | 'DELETE' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE'))*
    ;
revokePrivilege_schemaObjectPrivileges_21
    : ('REFERENCES' | 'SELECT') (',' ('REFERENCES' | 'SELECT'))*
    ;
revokePrivilege_schemaObjectPrivileges_22
    : ('APPLYBUDGET' | 'REFERENCES' | 'SELECT') (',' ('APPLYBUDGET' | 'REFERENCES' | 'SELECT'))*
    ;
revokePrivilege_schemaObjectPrivileges_23
    : ('MODIFY' | 'USAGE') (',' ('MODIFY' | 'USAGE'))*
    ;
revokePrivilege_schemaObjectPrivileges_24
    : ('MONITOR' | 'SELECT') (',' ('MONITOR' | 'SELECT'))*
    ;
revokePrivilege_schemaObjectPrivileges_25
    : ('APPLYBUDGET' | 'MONITOR' | 'OPERATE') (',' ('APPLYBUDGET' | 'MONITOR' | 'OPERATE'))*
    ;
revokePrivilege_schemaObjectPrivileges_26
    : 'APPLY' (',' 'APPLY')*
    ;
revokePrivilege_schemaObjectPrivileges_27
    : ('READ' | 'USAGE') (',' ('READ' | 'USAGE'))*
    ;
revokePrivilege_schemaObjectPrivileges_28
    : ('SELECT' | 'REFERENCES' | 'MONITOR') (',' ('SELECT' | 'REFERENCES' | 'MONITOR'))*
    ;
revokePrivilege_schemaObjectPrivileges_29
    : ('MONITOR' | 'OPERATE') (',' ('MONITOR' | 'OPERATE'))*
    ;
revokePrivilege_schemaObjectPrivileges_30
    : 'USAGE' (',' 'USAGE')*
    ;
revokePrivilege_schemaObjectPrivileges_31
    : 'READ' (',' 'WRITE')? (',' 'READ' (',' 'WRITE')?)*
    ;
revokePrivilege_schemaObjectPrivileges_32
    : 'SELECT' (',' 'SELECT')*
    ;
revokePrivilege_schemaObjectPrivileges_33
    : 'USAGE' (',' 'USAGE')*
    ;
revokePrivilege_schemaObjectPrivileges_34
    : ('APPLYBUDGET' | 'DELETE' | 'EVOLVE' 'SCHEMA' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'SELECT' 'ERROR' 'TABLE' | 'TRUNCATE' | 'UPDATE') (',' ('APPLYBUDGET' | 'DELETE' | 'EVOLVE' 'SCHEMA' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'SELECT' 'ERROR' 'TABLE' | 'TRUNCATE' | 'UPDATE'))*
    ;
revokePrivilege_schemaObjectPrivileges_35
    : 'READ'
    ;
revokePrivilege_schemaObjectPrivileges_36
    : ('APPLYBUDGET' | 'MONITOR' | 'OPERATE') (',' ('APPLYBUDGET' | 'MONITOR' | 'OPERATE'))*
    ;
revokePrivilege_schemaObjectPrivileges_37
    : ('REFERENCES' | 'SELECT') (',' ('REFERENCES' | 'SELECT'))*
    ;
revokePrivilege_schemaObjectPrivileges_38
    : ('READ' | 'WRITE') (',' ('READ' | 'WRITE'))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/revoke-privilege-application
// fix: ON { <object_type> ... is not closed; closed
// fix: read: each -- For X section of accountObjectPrivileges and schemaObjectPrivileges is the privileges on X alone, and the ON clause chooses an object type with its own privileges
revokePrivilegeApplication
    : 'REVOKE' (revokePrivilegeApplication_globalPrivileges 'ON' 'ACCOUNT' | ((revokePrivilegeApplication_accountObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'USER' ID | (revokePrivilegeApplication_accountObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' 'RESOURCE' 'MONITOR' ID | (revokePrivilegeApplication_accountObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' 'WAREHOUSE' ID | (revokePrivilegeApplication_accountObjectPrivileges_1 | 'ALL' 'PRIVILEGES'?) 'ON' 'COMPUTE' 'POOL' ID | (revokePrivilegeApplication_accountObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' 'DATABASE' ID | (revokePrivilegeApplication_accountObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' 'INTEGRATION' ID | (revokePrivilegeApplication_accountObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' 'CONNECTION' ID | (revokePrivilegeApplication_accountObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' 'FAILOVER' 'GROUP' ID | (revokePrivilegeApplication_accountObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' 'REPLICATION' 'GROUP' ID | (revokePrivilegeApplication_accountObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' 'EXTERNAL' 'VOLUME' ID) | (revokePrivilegeApplication_schemaPrivileges | 'ALL' 'PRIVILEGES'?) 'ON' ('SCHEMA' QNAME | 'ALL' 'SCHEMAS' 'IN' 'DATABASE' ID) | ((revokePrivilegeApplication_schemaObjectPrivileges_1 | 'ALL' 'PRIVILEGES'?) 'ON' ('ALERT' QNAME | 'ALL' 'ALERTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' ('DYNAMIC' 'TABLE' QNAME | 'ALL' 'DYNAMIC' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' ('EVENT' 'TABLE' QNAME | 'ALL' 'EVENT' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('FILE' 'FORMAT' QNAME | 'ALL' 'FILE' 'FORMATS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('FUNCTION' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' | 'ALL' 'FUNCTIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('IMAGE' 'REPOSITORY' QNAME | 'ALL' 'IMAGE' 'REPOSITORIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' 'MASKING' 'POLICY' QNAME | (revokePrivilegeApplication_schemaObjectPrivileges_18 | 'ALL' 'PRIVILEGES'?) 'ON' ('MATERIALIZED' 'VIEW' QNAME | 'ALL' 'MATERIALIZED' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' ('PACKAGES' 'POLICY' QNAME | 'ALL' 'PACKAGES' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' ('PASSWORD' 'POLICY' QNAME | 'ALL' 'PASSWORD' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' 'PIPE' QNAME | (revokePrivilegeApplication_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('PROCEDURE' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' | 'ALL' 'PROCEDURES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' ('ROW' 'ACCESS' 'POLICY' QNAME | 'ALL' 'ROW' 'ACCESS' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_4 | revokePrivilegeApplication_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('SECRET' QNAME | 'ALL' 'SECRETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('SEMANTIC' 'VIEW' QNAME | 'ALL' 'SEMANTIC' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('SEQUENCE' QNAME | 'ALL' 'SEQUENCES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' ('SERVICE' QNAME | 'ALL' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' ('SESSION' 'POLICY' QNAME | 'ALL' 'SESSION' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('SNAPSHOT' QNAME | 'ALL' 'SNAPSHOTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_11 | revokePrivilegeApplication_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('STAGE' QNAME | 'ALL' 'STAGES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_13 | 'ALL' 'PRIVILEGES'?) 'ON' ('STREAM' QNAME | 'ALL' 'STREAMS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_14 | 'ALL' 'PRIVILEGES'?) 'ON' ('TABLE' QNAME | 'ALL' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_7 | revokePrivilegeApplication_schemaObjectPrivileges_15 | 'ALL' 'PRIVILEGES'?) 'ON' 'TAG' QNAME | (revokePrivilegeApplication_schemaObjectPrivileges_16 | 'ALL' 'PRIVILEGES'?) 'ON' ('TASK' QNAME | 'ALL' 'TASKS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('TYPE' QNAME | 'ALL' 'TYPES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)) | (revokePrivilegeApplication_schemaObjectPrivileges_17 | 'ALL' 'PRIVILEGES'?) 'ON' ('VIEW' QNAME | 'ALL' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME)))) 'FROM' 'APPLICATION' QNAME
    ;
revokePrivilegeApplication_globalPrivileges
    : ('CREATE' ('COMPUTE' 'POOL' | 'DATABASE' | 'WAREHOUSE') | 'BIND' 'SERVICE' 'ENDPOINT' | 'EXECUTE' 'MANAGED' 'TASK' | 'MANAGE' 'WAREHOUSES' | 'READ' 'SESSION') (',' ('CREATE' ('COMPUTE' 'POOL' | 'DATABASE' | 'WAREHOUSE') | 'BIND' 'SERVICE' 'ENDPOINT' | 'EXECUTE' 'MANAGED' 'TASK' | 'MANAGE' 'WAREHOUSES' | 'READ' 'SESSION'))*
    ;
revokePrivilegeApplication_accountObjectPrivileges_1
    : ('MODIFY' | 'MONITOR' | 'OPERATE' | 'USAGE') (',' ('MODIFY' | 'MONITOR' | 'OPERATE' | 'USAGE'))*
    ;
revokePrivilegeApplication_accountObjectPrivileges_2
    : 'FAILOVER' (',' 'FAILOVER')*
    ;
revokePrivilegeApplication_accountObjectPrivileges_3
    : ('APPLYBUDGET' | 'CREATE' ('DATABASE' 'ROLE' | 'SCHEMA') | 'IMPORTED' 'PRIVILEGES' | 'MODIFY' | 'MONITOR' | 'USAGE') (',' ('APPLYBUDGET' | 'CREATE' ('DATABASE' 'ROLE' | 'SCHEMA') | 'IMPORTED' 'PRIVILEGES' | 'MODIFY' | 'MONITOR' | 'USAGE'))*
    ;
revokePrivilegeApplication_accountObjectPrivileges_4
    : 'USAGE' (',' 'USAGE')*
    ;
revokePrivilegeApplication_accountObjectPrivileges_5
    : ('FAILOVER' | 'MODIFY' | 'MONITOR' | 'REPLICATE') (',' ('FAILOVER' | 'MODIFY' | 'MONITOR' | 'REPLICATE'))*
    ;
revokePrivilegeApplication_accountObjectPrivileges_6
    : ('USAGE' | 'USE_ANY_ROLE') (',' ('USAGE' | 'USE_ANY_ROLE'))*
    ;
revokePrivilegeApplication_accountObjectPrivileges_7
    : ('MODIFY' | 'MONITOR' | 'REPLICATE') (',' ('MODIFY' | 'MONITOR' | 'REPLICATE'))*
    ;
revokePrivilegeApplication_accountObjectPrivileges_8
    : ('MODIFY' | 'MONITOR') (',' ('MODIFY' | 'MONITOR'))*
    ;
revokePrivilegeApplication_accountObjectPrivileges_9
    : 'MONITOR' (',' 'MONITOR')*
    ;
revokePrivilegeApplication_accountObjectPrivileges_10
    : ('APPLYBUDGET' | 'MODIFY' | 'MONITOR' | 'USAGE' | 'OPERATE') (',' ('APPLYBUDGET' | 'MODIFY' | 'MONITOR' | 'USAGE' | 'OPERATE'))*
    ;
revokePrivilegeApplication_schemaPrivileges
    : ('ADD' 'SEARCH' 'OPTIMIZATION' | 'CREATE' ('ALERT' | 'EXTERNAL' 'TABLE' | 'FILE' 'FORMAT' | 'FUNCTION' | 'IMAGE' 'REPOSITORY' | 'MATERIALIZED' 'VIEW' | 'PIPE' | 'PROCEDURE' | ('AGGREGATION' | 'MASKING' | 'PASSWORD' | 'PROJECTION' | 'ROW' 'ACCESS' | 'SESSION') 'POLICY' | 'SECRET' | 'SEMANTIC' 'VIEW' | 'SEQUENCE' | 'SERVICE' | 'SNAPSHOT' | 'STAGE' | 'STREAM' | 'TAG' | 'TABLE' | 'TASK' | 'VIEW') | 'MODIFY' | 'MONITOR' | 'USAGE') (',' ('ADD' 'SEARCH' 'OPTIMIZATION' | 'CREATE' ('ALERT' | 'EXTERNAL' 'TABLE' | 'FILE' 'FORMAT' | 'FUNCTION' | 'IMAGE' 'REPOSITORY' | 'MATERIALIZED' 'VIEW' | 'PIPE' | 'PROCEDURE' | ('AGGREGATION' | 'MASKING' | 'PASSWORD' | 'PROJECTION' | 'ROW' 'ACCESS' | 'SESSION') 'POLICY' | 'SECRET' | 'SEMANTIC' 'VIEW' | 'SEQUENCE' | 'SERVICE' | 'SNAPSHOT' | 'STAGE' | 'STREAM' | 'TAG' | 'TABLE' | 'TASK' | 'VIEW') | 'MODIFY' | 'MONITOR' | 'USAGE'))*
    ;
revokePrivilegeApplication_schemaObjectPrivileges_1
    : ('MONITOR' | 'OPERATE') (',' ('MONITOR' | 'OPERATE'))*
    ;
revokePrivilegeApplication_schemaObjectPrivileges_2
    : 'OPERATE' ',' 'SELECT' (',' 'SELECT')*
    ;
revokePrivilegeApplication_schemaObjectPrivileges_3
    : ('INSERT' | 'SELECT') (',' ('INSERT' | 'SELECT'))*
    ;
revokePrivilegeApplication_schemaObjectPrivileges_4
    : 'USAGE' (',' 'USAGE')*
    ;
revokePrivilegeApplication_schemaObjectPrivileges_5
    : 'READ' ',' 'WRITE' (',' 'READ' ',' 'WRITE')*
    ;
revokePrivilegeApplication_schemaObjectPrivileges_6
    : ('APPLYBUDGET' | 'MONITOR' | 'OPERATE') (',' ('APPLYBUDGET' | 'MONITOR' | 'OPERATE'))*
    ;
revokePrivilegeApplication_schemaObjectPrivileges_7
    : 'APPLY' (',' 'APPLY')*
    ;
revokePrivilegeApplication_schemaObjectPrivileges_8
    : 'READ' ',' 'USAGE' (',' 'USAGE')*
    ;
revokePrivilegeApplication_schemaObjectPrivileges_9
    : 'REFERENCES' (',' 'REFERENCES')*
    ;
revokePrivilegeApplication_schemaObjectPrivileges_10
    : ('MONITOR' | 'OPERATE') (',' ('MONITOR' | 'OPERATE'))*
    ;
revokePrivilegeApplication_schemaObjectPrivileges_11
    : 'USAGE' (',' 'USAGE')*
    ;
revokePrivilegeApplication_schemaObjectPrivileges_12
    : 'READ' (',' 'WRITE')? (',' 'READ' (',' 'WRITE')?)*
    ;
revokePrivilegeApplication_schemaObjectPrivileges_13
    : 'SELECT' (',' 'SELECT')*
    ;
revokePrivilegeApplication_schemaObjectPrivileges_14
    : ('APPLYBUDGET' | 'DELETE' | 'EVOLVE' 'SCHEMA' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE') (',' ('APPLYBUDGET' | 'DELETE' | 'EVOLVE' 'SCHEMA' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE'))*
    ;
revokePrivilegeApplication_schemaObjectPrivileges_15
    : 'READ'
    ;
revokePrivilegeApplication_schemaObjectPrivileges_16
    : ('APPLYBUDGET' | 'MONITOR' | 'OPERATE') (',' ('APPLYBUDGET' | 'MONITOR' | 'OPERATE'))*
    ;
revokePrivilegeApplication_schemaObjectPrivileges_17
    : ('REFERENCES' | 'SELECT') (',' ('REFERENCES' | 'SELECT'))*
    ;
revokePrivilegeApplication_schemaObjectPrivileges_18
    : ('APPLYBUDGET' | 'REFERENCES' | 'SELECT') (',' ('APPLYBUDGET' | 'REFERENCES' | 'SELECT'))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/revoke-privilege-application-role
// fix: read: each -- For X section of accountObjectPrivileges and schemaObjectPrivileges is the privileges on X alone, and the ON clause chooses an object type with its own privileges
revokePrivilegeApplicationRole
    : 'REVOKE' ('GRANT' 'OPTION' 'FOR')? ((revokePrivilegeApplicationRole_schemaPrivileges | 'ALL' 'PRIVILEGES'?) 'ON' ('SCHEMA' QNAME | 'ALL' 'SCHEMAS' 'IN' 'DATABASE' ID) | (revokePrivilegeApplicationRole_schemaPrivileges | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SCHEMAS' 'IN' 'DATABASE' ID | ((revokePrivilegeApplicationRole_schemaObjectPrivileges_1 | 'ALL' 'PRIVILEGES'?) 'ON' ('ALERT' QNAME | 'ALL' 'ALERTS' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' ('DYNAMIC' 'TABLE' QNAME | 'ALL' 'DYNAMIC' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' ('EVENT' 'TABLE' QNAME | 'ALL' 'EVENT' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('FILE' 'FORMAT' QNAME | 'ALL' 'FILE' 'FORMATS' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('FUNCTION' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' | 'ALL' 'FUNCTIONS' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('IMAGE' 'REPOSITORY' QNAME | 'ALL' 'IMAGE' 'REPOSITORIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' 'MASKING' 'POLICY' QNAME | (revokePrivilegeApplicationRole_schemaObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' ('MATERIALIZED' 'VIEW' QNAME | 'ALL' 'MATERIALIZED' 'VIEWS' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('PACKAGES' 'POLICY' QNAME | 'ALL' 'PACKAGES' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('PASSWORD' 'POLICY' QNAME | 'ALL' 'PASSWORD' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' 'PIPE' QNAME | (revokePrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('PROCEDURE' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' | 'ALL' 'PROCEDURES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('ROW' 'ACCESS' 'POLICY' QNAME | 'ALL' 'ROW' 'ACCESS' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_4 | revokePrivilegeApplicationRole_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('SECRET' QNAME | 'ALL' 'SECRETS' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' ('SEMANTIC' 'VIEW' QNAME | 'ALL' 'SEMANTIC' 'VIEWS' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('SEQUENCE' QNAME | 'ALL' 'SEQUENCES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_11 | 'ALL' 'PRIVILEGES'?) 'ON' ('SERVICE' QNAME | 'ALL' 'SERVICES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('SESSION' 'POLICY' QNAME | 'ALL' 'SESSION' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('SNAPSHOT' QNAME | 'ALL' 'SNAPSHOTS' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_12 | revokePrivilegeApplicationRole_schemaObjectPrivileges_13 | 'ALL' 'PRIVILEGES'?) 'ON' ('STAGE' QNAME | 'ALL' 'STAGES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_14 | 'ALL' 'PRIVILEGES'?) 'ON' ('STREAM' QNAME | 'ALL' 'STREAMS' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_15 | 'ALL' 'PRIVILEGES'?) 'ON' ('TABLE' QNAME | 'ALL' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_8 | revokePrivilegeApplicationRole_schemaObjectPrivileges_16 | 'ALL' 'PRIVILEGES'?) 'ON' 'TAG' QNAME | (revokePrivilegeApplicationRole_schemaObjectPrivileges_17 | 'ALL' 'PRIVILEGES'?) 'ON' ('TASK' QNAME | 'ALL' 'TASKS' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('TYPE' QNAME | 'ALL' 'TYPES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_18 | 'ALL' 'PRIVILEGES'?) 'ON' ('VIEW' QNAME | 'ALL' 'VIEWS' 'IN' 'SCHEMA' QNAME)) | ((revokePrivilegeApplicationRole_schemaObjectPrivileges_1 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'ALERTS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'DYNAMIC' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'EVENT' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'FILE' 'FORMATS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'FUNCTIONS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'MATERIALIZED' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PASSWORD' 'POLICIES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PIPES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'PROCEDURES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_4 | revokePrivilegeApplicationRole_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SECRETS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SEMANTIC' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SEQUENCES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_11 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'SERVICES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_12 | revokePrivilegeApplicationRole_schemaObjectPrivileges_13 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'STAGES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_14 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'STREAMS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_15 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'TABLES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_17 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'TASKS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'TYPES' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME) | (revokePrivilegeApplicationRole_schemaObjectPrivileges_18 | 'ALL' 'PRIVILEGES'?) 'ON' 'FUTURE' 'VIEWS' 'IN' ('DATABASE' ID | 'SCHEMA' QNAME))) 'FROM' 'APPLICATION' 'ROLE' QNAME ('RESTRICT' | 'CASCADE')?
    ;
revokePrivilegeApplicationRole_schemaPrivileges
    : ('ADD' 'SEARCH' 'OPTIMIZATION' | 'CREATE' ('ALERT' | 'EXTERNAL' 'TABLE' | 'FILE' 'FORMAT' | 'FUNCTION' | 'IMAGE' 'REPOSITORY' | 'MATERIALIZED' 'VIEW' | 'PIPE' | 'PROCEDURE' | ('AGGREGATION' | 'MASKING' | 'PASSWORD' | 'PROJECTION' | 'ROW' 'ACCESS' | 'SESSION') 'POLICY' | 'SECRET' | 'SEMANTIC' 'VIEW' | 'SEQUENCE' | 'SERVICE' | 'SNAPSHOT' | 'STAGE' | 'STREAM' | 'TAG' | 'TABLE' | 'TASK' | 'VIEW') | 'MODIFY' | 'MONITOR' | 'USAGE') (',' ('ADD' 'SEARCH' 'OPTIMIZATION' | 'CREATE' ('ALERT' | 'EXTERNAL' 'TABLE' | 'FILE' 'FORMAT' | 'FUNCTION' | 'IMAGE' 'REPOSITORY' | 'MATERIALIZED' 'VIEW' | 'PIPE' | 'PROCEDURE' | ('AGGREGATION' | 'MASKING' | 'PASSWORD' | 'PROJECTION' | 'ROW' 'ACCESS' | 'SESSION') 'POLICY' | 'SECRET' | 'SEMANTIC' 'VIEW' | 'SEQUENCE' | 'SERVICE' | 'SNAPSHOT' | 'STAGE' | 'STREAM' | 'TAG' | 'TABLE' | 'TASK' | 'VIEW') | 'MODIFY' | 'MONITOR' | 'USAGE'))*
    ;
revokePrivilegeApplicationRole_schemaObjectPrivileges_1
    : ('MONITOR' | 'OPERATE') (',' ('MONITOR' | 'OPERATE'))*
    ;
revokePrivilegeApplicationRole_schemaObjectPrivileges_2
    : 'OPERATE' ',' 'SELECT' (',' 'SELECT')*
    ;
revokePrivilegeApplicationRole_schemaObjectPrivileges_3
    : ('INSERT' | 'SELECT') (',' ('INSERT' | 'SELECT'))*
    ;
revokePrivilegeApplicationRole_schemaObjectPrivileges_4
    : 'USAGE' (',' 'USAGE')*
    ;
revokePrivilegeApplicationRole_schemaObjectPrivileges_5
    : 'READ' ',' 'WRITE' (',' 'READ' ',' 'WRITE')*
    ;
revokePrivilegeApplicationRole_schemaObjectPrivileges_6
    : ('APPLYBUDGET' | 'REFERENCES' | 'SELECT') (',' ('APPLYBUDGET' | 'REFERENCES' | 'SELECT'))*
    ;
revokePrivilegeApplicationRole_schemaObjectPrivileges_7
    : ('APPLYBUDGET' | 'MONITOR' | 'OPERATE') (',' ('APPLYBUDGET' | 'MONITOR' | 'OPERATE'))*
    ;
revokePrivilegeApplicationRole_schemaObjectPrivileges_8
    : 'APPLY' (',' 'APPLY')*
    ;
revokePrivilegeApplicationRole_schemaObjectPrivileges_9
    : 'READ' ',' 'USAGE' (',' 'USAGE')*
    ;
revokePrivilegeApplicationRole_schemaObjectPrivileges_10
    : 'REFERENCES' (',' 'REFERENCES')*
    ;
revokePrivilegeApplicationRole_schemaObjectPrivileges_11
    : ('MONITOR' | 'OPERATE') (',' ('MONITOR' | 'OPERATE'))*
    ;
revokePrivilegeApplicationRole_schemaObjectPrivileges_12
    : 'USAGE' (',' 'USAGE')*
    ;
revokePrivilegeApplicationRole_schemaObjectPrivileges_13
    : 'READ' (',' 'WRITE')? (',' 'READ' (',' 'WRITE')?)*
    ;
revokePrivilegeApplicationRole_schemaObjectPrivileges_14
    : 'SELECT' (',' 'SELECT')*
    ;
revokePrivilegeApplicationRole_schemaObjectPrivileges_15
    : ('APPLYBUDGET' | 'DELETE' | 'EVOLVE' 'SCHEMA' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE') (',' ('APPLYBUDGET' | 'DELETE' | 'EVOLVE' 'SCHEMA' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE'))*
    ;
revokePrivilegeApplicationRole_schemaObjectPrivileges_16
    : 'READ'
    ;
revokePrivilegeApplicationRole_schemaObjectPrivileges_17
    : ('APPLYBUDGET' | 'MONITOR' | 'OPERATE') (',' ('APPLYBUDGET' | 'MONITOR' | 'OPERATE'))*
    ;
revokePrivilegeApplicationRole_schemaObjectPrivileges_18
    : ('REFERENCES' | 'SELECT') (',' ('REFERENCES' | 'SELECT'))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/revoke-privilege-share
revokePrivilegeShare
    : 'REVOKE' revokePrivilegeShare_objectPrivilege 'ON' ('DATABASE' QNAME | 'SCHEMA' QNAME | 'MODEL' QNAME | 'SEMANTIC' 'VIEW' QNAME | ('TABLE' QNAME | 'ALL' 'TABLES' 'IN' 'SCHEMA' QNAME) | ('EXTERNAL' 'TABLE' QNAME | 'ALL' 'EXTERNAL' 'TABLES' 'IN' 'SCHEMA' QNAME) | ('ICEBERG' 'TABLE' QNAME | 'ALL' 'ICEBERG' 'TABLES' 'IN' 'SCHEMA' QNAME) | ('DYNAMIC' 'TABLE' QNAME | 'ALL' 'DYNAMIC' 'TABLES' 'IN' 'SCHEMA' QNAME) | ('VIEW' QNAME | 'ALL' 'VIEWS' 'IN' 'SCHEMA' QNAME)) 'FROM' 'SHARE' ID
    ;
revokePrivilegeShare_objectPrivilege
    : 'REFERENCE_USAGE' (',' 'REFERENCE_USAGE')*
    | 'USAGE' (',' 'USAGE')*
    | 'READ'
    | ('REFERENCES' | 'SELECT') (',' ('REFERENCES' | 'SELECT'))*
    | 'EVOLVE' 'SCHEMA' (',' 'EVOLVE' 'SCHEMA')*
    | 'SELECT' (',' 'SELECT')*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/revoke-privilege-user
// fix: read: each -- For X section of accountObjectPrivileges and schemaObjectPrivileges is the privileges on X alone, and the ON clause chooses an object type with its own privileges
revokePrivilegeUser
    : 'REVOKE' ('GRANT' 'OPTION' 'FOR')? ((revokePrivilegeUser_globalPrivileges | 'ALL' 'PRIVILEGES'?) 'ON' 'ACCOUNT' | ((revokePrivilegeUser_accountObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' 'RESOURCE' 'MONITOR' ID | (revokePrivilegeUser_accountObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' 'WAREHOUSE' ID | (revokePrivilegeUser_accountObjectPrivileges_1 | 'ALL' 'PRIVILEGES'?) 'ON' 'COMPUTE' 'POOL' ID | (revokePrivilegeUser_accountObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' 'DATABASE' ID | (revokePrivilegeUser_accountObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' 'INTEGRATION' ID | (revokePrivilegeUser_accountObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' 'CONNECTION' ID | (revokePrivilegeUser_accountObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' 'FAILOVER' 'GROUP' ID | (revokePrivilegeUser_accountObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' 'REPLICATION' 'GROUP' ID | (revokePrivilegeUser_accountObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' 'EXTERNAL' 'VOLUME' ID) | (revokePrivilegeUser_schemaPrivileges | 'ALL' 'PRIVILEGES'?) 'ON' ('SCHEMA' QNAME | 'ALL' 'SCHEMAS' 'IN' 'DATABASE' ID) | ((revokePrivilegeUser_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('AGGREGATION' 'POLICY' QNAME | 'ALL' 'AGGREGATION' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_1 | 'ALL' 'PRIVILEGES'?) 'ON' ('ALERT' QNAME | 'ALL' 'ALERTS' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('AUTHENTICATION' 'POLICY' QNAME | 'ALL' 'AUTHENTICATION' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_2 | 'ALL' 'PRIVILEGES'?) 'ON' ('DATA' 'METRIC' 'FUNCTION' QNAME | 'ALL' 'DATA' 'METRIC' 'FUNCTIONS' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_3 | 'ALL' 'PRIVILEGES'?) 'ON' ('DYNAMIC' 'TABLE' QNAME | 'ALL' 'DYNAMIC' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_4 | 'ALL' 'PRIVILEGES'?) 'ON' ('EVENT' 'TABLE' QNAME | 'ALL' 'EVENT' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('FILE' 'FORMAT' QNAME | 'ALL' 'FILE' 'FORMATS' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('FUNCTION' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' | 'ALL' 'FUNCTIONS' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_6 | 'ALL' 'PRIVILEGES'?) 'ON' ('GIT' 'REPOSITORY' QNAME | 'ALL' 'GIT' 'REPOSITORIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_7 | 'ALL' 'PRIVILEGES'?) 'ON' ('HYBRID' 'TABLE' QNAME | 'ALL' 'HYBRID' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_9 | 'ALL' 'PRIVILEGES'?) 'ON' ('ICEBERG' 'TABLE' QNAME | 'ALL' 'ICEBERG' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_8 | 'ALL' 'PRIVILEGES'?) 'ON' ('IMAGE' 'REPOSITORY' QNAME | 'ALL' 'IMAGE' 'REPOSITORIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('JOIN' 'POLICY' QNAME | 'ALL' 'JOIN' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' 'MASKING' 'POLICY' QNAME | (revokePrivilegeUser_schemaObjectPrivileges_10 | 'ALL' 'PRIVILEGES'?) 'ON' ('MATERIALIZED' 'VIEW' QNAME | 'ALL' 'MATERIALIZED' 'VIEWS' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('MODEL' QNAME | 'ALL' 'MODELS' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('PACKAGES' 'POLICY' QNAME | 'ALL' 'PACKAGES' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('PASSWORD' 'POLICY' QNAME | 'ALL' 'PASSWORD' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_11 | 'ALL' 'PRIVILEGES'?) 'ON' 'PIPE' QNAME | (revokePrivilegeUser_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('PRIVACY' 'POLICY' QNAME | 'ALL' 'PRIVACY' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('PROCEDURE' QNAME '(' (DATA_TYPE (',' DATA_TYPE)*)? ')' | 'ALL' 'PROCEDURES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('PROJECTION' 'POLICY' QNAME | 'ALL' 'PROJECTION' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('ROW' 'ACCESS' 'POLICY' QNAME | 'ALL' 'ROW' 'ACCESS' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_5 | revokePrivilegeUser_schemaObjectPrivileges_13 | 'ALL' 'PRIVILEGES'?) 'ON' ('SECRET' QNAME | 'ALL' 'SECRETS' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_14 | 'ALL' 'PRIVILEGES'?) 'ON' ('SEMANTIC' 'VIEW' QNAME | 'ALL' 'SEMANTIC' 'VIEWS' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('SEQUENCE' QNAME | 'ALL' 'SEQUENCES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_15 | 'ALL' 'PRIVILEGES'?) 'ON' ('SERVICE' QNAME | 'ALL' 'SERVICES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_12 | 'ALL' 'PRIVILEGES'?) 'ON' ('SESSION' 'POLICY' QNAME | 'ALL' 'SESSION' 'POLICIES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('SNAPSHOT' QNAME | 'ALL' 'SNAPSHOTS' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_16 | revokePrivilegeUser_schemaObjectPrivileges_17 | 'ALL' 'PRIVILEGES'?) 'ON' ('STAGE' QNAME | 'ALL' 'STAGES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_18 | 'ALL' 'PRIVILEGES'?) 'ON' ('STREAM' QNAME | 'ALL' 'STREAMS' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_19 | 'ALL' 'PRIVILEGES'?) 'ON' ('STREAMLIT' QNAME | 'ALL' 'STREAMLITS' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_20 | 'ALL' 'PRIVILEGES'?) 'ON' ('TABLE' QNAME | 'ALL' 'TABLES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_12 | revokePrivilegeUser_schemaObjectPrivileges_21 | 'ALL' 'PRIVILEGES'?) 'ON' 'TAG' QNAME | (revokePrivilegeUser_schemaObjectPrivileges_22 | 'ALL' 'PRIVILEGES'?) 'ON' ('TASK' QNAME | 'ALL' 'TASKS' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_5 | 'ALL' 'PRIVILEGES'?) 'ON' ('TYPE' QNAME | 'ALL' 'TYPES' 'IN' 'SCHEMA' QNAME) | (revokePrivilegeUser_schemaObjectPrivileges_23 | 'ALL' 'PRIVILEGES'?) 'ON' ('VIEW' QNAME | 'ALL' 'VIEWS' 'IN' 'SCHEMA' QNAME))) 'FROM' 'USER'? ID ('RESTRICT' | 'CASCADE')?
    ;
revokePrivilegeUser_globalPrivileges
    : ('ATTACH' 'POLICY' | 'AUDIT' | 'BIND' 'SERVICE' 'ENDPOINT' | 'APPLY' (('AGGREGATION' | 'AUTHENTICATION' | 'JOIN' | 'MASKING' | 'PACKAGES' | 'PASSWORD' | 'PROJECTION' | 'ROW' 'ACCESS' | 'SESSION') 'POLICY' | 'TAG') | 'EXECUTE' ('ALERT' | 'DATA' 'METRIC' 'FUNCTION' | 'MANAGED' 'ALERT' | 'MANAGED' 'TASK' | 'TASK') | 'IMPORT' 'SHARE' | 'MANAGE' ('ACCOUNT' 'SUPPORT' 'CASES' | 'EVENT' 'SHARING' | 'GRANTS' | 'LISTING' 'AUTO' 'FULFILLMENT' | 'ORGANIZATION' 'SUPPORT' 'CASES' | 'USER' 'SUPPORT' 'CASES' | 'WAREHOUSES') | 'MODIFY' ('LOG' 'LEVEL' | 'TRACE' 'LEVEL' | 'SESSION' 'LOG' 'LEVEL' | 'SESSION' 'TRACE' 'LEVEL') | 'MONITOR' ('EXECUTION' | 'SECURITY' | 'USAGE') | 'OVERRIDE' 'SHARE' 'RESTRICTIONS' | 'PURCHASE' 'DATA' 'EXCHANGE' 'LISTING' | 'RESOLVE' 'ALL' | 'READ' 'SESSION') (',' ('ATTACH' 'POLICY' | 'AUDIT' | 'BIND' 'SERVICE' 'ENDPOINT' | 'APPLY' (('AGGREGATION' | 'AUTHENTICATION' | 'JOIN' | 'MASKING' | 'PACKAGES' | 'PASSWORD' | 'PROJECTION' | 'ROW' 'ACCESS' | 'SESSION') 'POLICY' | 'TAG') | 'EXECUTE' ('ALERT' | 'DATA' 'METRIC' 'FUNCTION' | 'MANAGED' 'ALERT' | 'MANAGED' 'TASK' | 'TASK') | 'IMPORT' 'SHARE' | 'MANAGE' ('ACCOUNT' 'SUPPORT' 'CASES' | 'EVENT' 'SHARING' | 'GRANTS' | 'LISTING' 'AUTO' 'FULFILLMENT' | 'ORGANIZATION' 'SUPPORT' 'CASES' | 'USER' 'SUPPORT' 'CASES' | 'WAREHOUSES') | 'MODIFY' ('LOG' 'LEVEL' | 'TRACE' 'LEVEL' | 'SESSION' 'LOG' 'LEVEL' | 'SESSION' 'TRACE' 'LEVEL') | 'MONITOR' ('EXECUTION' | 'SECURITY' | 'USAGE') | 'OVERRIDE' 'SHARE' 'RESTRICTIONS' | 'PURCHASE' 'DATA' 'EXCHANGE' 'LISTING' | 'RESOLVE' 'ALL' | 'READ' 'SESSION'))*
    ;
revokePrivilegeUser_accountObjectPrivileges_1
    : ('MODIFY' | 'MONITOR' | 'OPERATE' | 'USAGE') (',' ('MODIFY' | 'MONITOR' | 'OPERATE' | 'USAGE'))*
    ;
revokePrivilegeUser_accountObjectPrivileges_2
    : 'FAILOVER' (',' 'FAILOVER')*
    ;
revokePrivilegeUser_accountObjectPrivileges_3
    : ('APPLYBUDGET' | 'IMPORTED' 'PRIVILEGES' | 'MODIFY' | 'MONITOR' | 'USAGE') (',' ('APPLYBUDGET' | 'IMPORTED' 'PRIVILEGES' | 'MODIFY' | 'MONITOR' | 'USAGE'))*
    ;
revokePrivilegeUser_accountObjectPrivileges_4
    : 'USAGE' (',' 'USAGE')*
    ;
revokePrivilegeUser_accountObjectPrivileges_5
    : ('FAILOVER' | 'MODIFY' | 'MONITOR' | 'REPLICATE') (',' ('FAILOVER' | 'MODIFY' | 'MONITOR' | 'REPLICATE'))*
    ;
revokePrivilegeUser_accountObjectPrivileges_6
    : ('USAGE' | 'USE_ANY_ROLE') (',' ('USAGE' | 'USE_ANY_ROLE'))*
    ;
revokePrivilegeUser_accountObjectPrivileges_7
    : ('MODIFY' | 'MONITOR' | 'REPLICATE') (',' ('MODIFY' | 'MONITOR' | 'REPLICATE'))*
    ;
revokePrivilegeUser_accountObjectPrivileges_8
    : ('MODIFY' | 'MONITOR') (',' ('MODIFY' | 'MONITOR'))*
    ;
revokePrivilegeUser_accountObjectPrivileges_10
    : ('APPLYBUDGET' | 'MODIFY' | 'MONITOR' | 'USAGE' | 'OPERATE') (',' ('APPLYBUDGET' | 'MODIFY' | 'MONITOR' | 'USAGE' | 'OPERATE'))*
    ;
revokePrivilegeUser_schemaPrivileges
    : ('ADD' 'SEARCH' 'OPTIMIZATION' | 'APPLYBUDGET' | 'MODIFY' | 'MONITOR' | 'USAGE') (',' ('ADD' 'SEARCH' 'OPTIMIZATION' | 'APPLYBUDGET' | 'MODIFY' | 'MONITOR' | 'USAGE'))*
    ;
revokePrivilegeUser_schemaObjectPrivileges_1
    : ('MONITOR' | 'OPERATE') (',' ('MONITOR' | 'OPERATE'))*
    ;
revokePrivilegeUser_schemaObjectPrivileges_2
    : 'USAGE' (',' 'USAGE')*
    ;
revokePrivilegeUser_schemaObjectPrivileges_3
    : 'MONITOR' ',' 'OPERATE' ',' 'SELECT' (',' 'SELECT')*
    ;
revokePrivilegeUser_schemaObjectPrivileges_4
    : ('APPLYBUDGET' | 'DELETE' | 'REFERENCES' | 'SELECT' | 'TRUNCATE') (',' ('APPLYBUDGET' | 'DELETE' | 'REFERENCES' | 'SELECT' | 'TRUNCATE'))*
    ;
revokePrivilegeUser_schemaObjectPrivileges_5
    : 'USAGE' (',' 'USAGE')*
    ;
revokePrivilegeUser_schemaObjectPrivileges_6
    : 'READ' ',' 'WRITE' (',' 'READ' ',' 'WRITE')*
    ;
revokePrivilegeUser_schemaObjectPrivileges_7
    : ('APPLYBUDGET' | 'DELETE' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE') (',' ('APPLYBUDGET' | 'DELETE' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE'))*
    ;
revokePrivilegeUser_schemaObjectPrivileges_8
    : 'READ' ',' 'WRITE' (',' 'READ' ',' 'WRITE')*
    ;
revokePrivilegeUser_schemaObjectPrivileges_9
    : ('APPLYBUDGET' | 'DELETE' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE') (',' ('APPLYBUDGET' | 'DELETE' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'TRUNCATE' | 'UPDATE'))*
    ;
revokePrivilegeUser_schemaObjectPrivileges_10
    : ('APPLYBUDGET' | 'REFERENCES' | 'SELECT') (',' ('APPLYBUDGET' | 'REFERENCES' | 'SELECT'))*
    ;
revokePrivilegeUser_schemaObjectPrivileges_11
    : ('APPLYBUDGET' | 'MONITOR' | 'OPERATE') (',' ('APPLYBUDGET' | 'MONITOR' | 'OPERATE'))*
    ;
revokePrivilegeUser_schemaObjectPrivileges_12
    : 'APPLY' (',' 'APPLY')*
    ;
revokePrivilegeUser_schemaObjectPrivileges_13
    : ('READ' | 'USAGE') (',' ('READ' | 'USAGE'))*
    ;
revokePrivilegeUser_schemaObjectPrivileges_14
    : 'REFERENCES' (',' 'REFERENCES')*
    ;
revokePrivilegeUser_schemaObjectPrivileges_15
    : ('MONITOR' | 'OPERATE') (',' ('MONITOR' | 'OPERATE'))*
    ;
revokePrivilegeUser_schemaObjectPrivileges_16
    : 'USAGE' (',' 'USAGE')*
    ;
revokePrivilegeUser_schemaObjectPrivileges_17
    : 'READ' (',' 'WRITE')? (',' 'READ' (',' 'WRITE')?)*
    ;
revokePrivilegeUser_schemaObjectPrivileges_18
    : 'SELECT' (',' 'SELECT')*
    ;
revokePrivilegeUser_schemaObjectPrivileges_19
    : 'USAGE' (',' 'USAGE')*
    ;
revokePrivilegeUser_schemaObjectPrivileges_20
    : ('APPLYBUDGET' | 'DELETE' | 'EVOLVE' 'SCHEMA' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'SELECT' 'ERROR' 'TABLE' | 'TRUNCATE' | 'UPDATE') (',' ('APPLYBUDGET' | 'DELETE' | 'EVOLVE' 'SCHEMA' | 'INSERT' | 'REFERENCES' | 'SELECT' | 'SELECT' 'ERROR' 'TABLE' | 'TRUNCATE' | 'UPDATE'))*
    ;
revokePrivilegeUser_schemaObjectPrivileges_21
    : 'READ'
    ;
revokePrivilegeUser_schemaObjectPrivileges_22
    : ('APPLYBUDGET' | 'MONITOR' | 'OPERATE') (',' ('APPLYBUDGET' | 'MONITOR' | 'OPERATE'))*
    ;
revokePrivilegeUser_schemaObjectPrivileges_23
    : ('REFERENCES' | 'SELECT') (',' ('REFERENCES' | 'SELECT'))*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/revoke-role
revokeRole
    : 'REVOKE' ('GRANT' 'OPTION' 'FOR')? 'ROLE' ID 'FROM' 'ROLE' ID ('RESTRICT' | 'CASCADE')?
    | 'REVOKE' 'ROLE' ID 'FROM' 'USER' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/revoke-service-role
revokeServiceRole
    : 'REVOKE' 'SERVICE' 'ROLE' QNAME 'FROM' ('ROLE' ID | 'APPLICATION' 'ROLE' ID | 'DATABASE' 'ROLE' ID)
    ;

// https://docs.snowflake.com/en/sql-reference/sql/rollback
rollback
    : 'ROLLBACK' 'WORK'?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/select
// not converted (form): ... stands for text the docs leave out
//   [ ... ]
// not converted (form): ... stands for text the docs leave out
//   SELECT * ILIKE ... REPLACE ...
// not converted (form): ... stands for text the docs leave out
//   SELECT * ILIKE ... RENAME ...
// not converted (form): ... stands for text the docs leave out
//   SELECT * ILIKE ... REPLACE ... RENAME ...
// not converted (form): ... stands for text the docs leave out
//   SELECT * EXCLUDE ... REPLACE ...
// not converted (form): ... stands for text the docs leave out
//   SELECT * EXCLUDE ... RENAME ...
// not converted (form): ... stands for text the docs leave out
//   SELECT * EXCLUDE ... REPLACE ... RENAME ...
// not converted (form): ... stands for text the docs leave out
//   SELECT * REPLACE ... RENAME ...
// not converted (form): ... stands for text the docs leave out
//   [ ... ]
// not converted (form): undefined fragment "emp_id": the page does not define it
//   SELECT emp_id,
//          name,
//          dept,
//     FROM employees;
select
    : 'SELECT' (('ALL' | 'DISTINCT'))? ('TOP' NUMBER)? ((QNAME | ID) '.')? '*' ('ILIKE' STRING)? ('EXCLUDE' (ID | '(' ID (',' ID)* ')'))? ('REPLACE' '(' EXPR 'AS' ID (',' EXPR 'AS' ID)* ')')? ('RENAME' (ID 'AS' ID | '(' ID 'AS' ID (',' ID 'AS' ID)* ')'))?
    | 'SELECT' (('ALL' | 'DISTINCT'))? ('TOP' NUMBER)? (((QNAME | ID) '.')? ID | ((QNAME | ID) '.')? POSITION | EXPR) ('AS'? ID)? (',' (((QNAME | ID) '.')? ID | ((QNAME | ID) '.')? POSITION | EXPR) ('AS'? ID)?)*
    ;

// https://docs.snowflake.com/en/sql-reference/sql/set
// fix: the variables and the values must be as many, as the usage notes say; written out for up to three
set
    : 'SET' ID '=' EXPR
    | 'SET' '(' ID ')' '=' '(' EXPR ')'
    | 'SET' '(' ID ',' ID ')' '=' '(' EXPR ',' EXPR ')'
    | 'SET' '(' ID ',' ID ',' ID ')' '=' '(' EXPR ',' EXPR ',' EXPR ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show
// no syntax

// https://docs.snowflake.com/en/sql-reference/sql/show-accounts
showAccounts
    : 'SHOW' 'ACCOUNTS' 'HISTORY'? ('LIKE' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-agents
showAgents
    : 'SHOW' 'AGENTS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID | 'SCHEMA' (ID '.')? ID))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-aggregation-policies
showAggregationPolicies
    : 'SHOW' 'AGGREGATION' 'POLICIES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'SCHEMA' QNAME?))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-alerts
// fix: read: the IN list's <schema_name> line has no | before APPLICATION; one is read there
showAlerts
    : 'SHOW' 'TERSE'? 'ALERTS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-application-packages
showApplicationPackages
    : 'SHOW' 'APPLICATION' 'PACKAGES' ('LIKE' STRING)? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-application-roles
// fix: LIKE takes a quoted pattern, as the page's parameters and examples write it
showApplicationRoles
    : 'SHOW' 'APPLICATION' 'ROLES' ('LIKE' STRING)? 'IN' 'APPLICATION' QNAME ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-application-services
showApplicationServices
    : 'SHOW' 'APPLICATION' 'SERVICES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'SCHEMA' QNAME?))? ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-applications
showApplications
    : 'SHOW' 'APPLICATIONS' ('LIKE' STRING)? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-artifact-repositories
showArtifactRepositories
    : 'SHOW' 'ARTIFACT' 'REPOSITORIES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'SCHEMA' QNAME?))? ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-authentication-policies
showAuthenticationPolicies
    : 'SHOW' 'AUTHENTICATION' 'POLICIES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID) | 'ON' ('ACCOUNT' | 'USER' ID))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-available-listings
showAvailableListings
    : 'SHOW' 'AVAILABLE' 'LISTINGS'
    | 'SHOW' 'TERSE'? 'AVAILABLE' 'LISTINGS' ('LIMIT' NUMBER)? ('IS_IMPORTED' '=' 'TRUE')? ('IS_ORGANIZATION' '=' 'TRUE')? ('IS_SHARED_WITH_ME' '=' 'TRUE')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-available-offers
showAvailableOffers
    : 'SHOW' 'AVAILABLE' 'OFFERS' ('LIKE' STRING)? 'IN' 'LISTING' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-available-organization-profiles
showAvailableOrganizationProfiles
    : 'SHOW' 'AVAILABLE' 'ORGANIZATION' 'PROFILES'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-backup-policies
showBackupPolicies
    : 'SHOW' 'BACKUP' 'POLICIES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME) ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-backup-sets
showBackupSets
    : 'SHOW' 'BACKUP' 'SETS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME) ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-backups-in-backup-set
showBackupsInBackupSet
    : 'SHOW' 'BACKUPS' 'IN' 'BACKUP' 'SET' QNAME ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-caller-grants
showCallerGrants
    : 'SHOW' 'CALLER' 'GRANTS' (('ON' OBJECT_TYPE QNAME | 'ON' 'ACCOUNT') | 'TO' ('ROLE' | 'DATABASE' 'ROLE') ID)
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-catalog-integrations
showCatalogIntegrations
    : 'SHOW' 'CATALOG' 'INTEGRATIONS' ('LIKE' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-channels
showChannels
    : 'SHOW' 'CHANNELS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME | 'TABLE' | 'TABLE' QNAME | 'PIPE' | 'PIPE' QNAME))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-classes
showClasses
    : 'SHOW' 'CLASSES' ('LIKE' STRING)? ('IN' 'DATABASE' ID?)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-columns
showColumns
    : 'SHOW' 'COLUMNS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'SCHEMA' QNAME? | 'TABLE' | 'TABLE'? QNAME | 'VIEW' | 'VIEW'? QNAME) | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-committed-tasks
showCommittedTasks
    : 'SHOW' 'COMMITTED' 'TASKS' ('LIKE' STRING)? ('IN' ('GRAPH' ID | 'ACCOUNT' | 'DATABASE' ID? | 'SCHEMA'? QNAME?))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-compute-pool-instance-families
showComputePoolInstanceFamilies
    : 'SHOW' 'COMPUTE' 'POOL' 'INSTANCE' 'FAMILIES'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-compute-pools
showComputePools
    : 'SHOW' 'COMPUTE' 'POOLS' ('LIKE' STRING)? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-configurations
showConfigurations
    : 'SHOW' 'CONFIGURATIONS' ('IN' 'APPLICATION' ID)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-connections
showConnections
    : 'SHOW' 'CONNECTIONS' ('LIKE' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-contacts
showContacts
    : 'SHOW' 'CONTACTS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-cortex-base-models
showCortexBaseModels
    : 'SHOW' 'CORTEX' 'BASE' 'MODELS' ('LIKE' STRING)? ('IN' 'SCHEMA'? 'SNOWFLAKE' '.' 'MODELS')? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-cortex-search
// fix: LIKE takes the pattern alone, as the page's parameters and examples write it
showCortexSearch
    : 'SHOW' 'CORTEX' 'SEARCH' 'SERVICES' ('LIKE' STRING)? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-data-metric-functions
showDataMetricFunctions
    : 'SHOW' 'DATA' 'METRIC' 'FUNCTIONS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME))? ('STARTS' 'WITH' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-data-movement-policies
showDataMovementPolicies
    : 'SHOW' 'DATA' 'MOVEMENT' 'POLICIES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-data-movement-rules
showDataMovementRules
    : 'SHOW' 'DATA' 'MOVEMENT' 'RULES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-database-roles
showDatabaseRoles
    : 'SHOW' 'DATABASE' 'ROLES' 'IN' 'DATABASE' QNAME ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-databases
// fix: the privileges to filter by are each object type's, listed on Access control privileges rather than here; the filter is left out
showDatabases
    : 'SHOW' 'TERSE'? 'DATABASES' 'HISTORY'? ('LIKE' STRING)? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-databases-in-failover-group
showDatabasesInFailoverGroup
    : 'SHOW' 'DATABASES' 'IN' 'FAILOVER' 'GROUP' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-databases-in-replication-group
showDatabasesInReplicationGroup
    : 'SHOW' 'DATABASES' 'IN' 'REPLICATION' 'GROUP' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-datasets
showDatasets
    : 'SHOW' 'DATASETS' ('LIKE' STRING)? ('IN' ('SCHEMA' QNAME | 'DATABASE' ID | 'ACCOUNT'))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-dbt-projects
showDbtProjects
    : 'SHOW' 'DBT' 'PROJECTS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-dcm-projects
showDcmProjects
    : 'SHOW' 'TERSE'? 'DCM' 'PROJECTS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME))? ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-delegated-authorizations
showDelegatedAuthorizations
    : 'SHOW' 'DELEGATED' 'AUTHORIZATIONS'
    | 'SHOW' 'DELEGATED' 'AUTHORIZATIONS' 'BY' 'USER' ID
    | 'SHOW' 'DELEGATED' 'AUTHORIZATIONS' 'TO' 'SECURITY' 'INTEGRATION' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-deployments-in-dcm-project
showDeploymentsInDcmProject
    : 'SHOW' 'DEPLOYMENTS' 'IN' 'DCM' 'PROJECT' QNAME ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-dynamic-tables
showDynamicTables
    : 'SHOW' 'DYNAMIC' 'TABLES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-endpoints
showEndpoints
    : 'SHOW' 'ENDPOINTS' 'IN' 'SERVICE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-entities-in-dcm-project
// fix: LIKE takes a quoted string, as the page's parameters and example write it
// fix: STARTS WITH takes a quoted string, as the page's parameters and example write it
// fix: FROM takes a quoted string, as the page's parameters and example write it
showEntitiesInDcmProject
    : 'SHOW' 'ENTITIES' 'IN' 'DCM' 'PROJECT' QNAME ('LIMIT' NUMBER)?
    | 'SHOW' 'ENTITIES' 'LIKE' STRING 'IN' 'DCM' 'PROJECT' QNAME
    | 'SHOW' 'ENTITIES' 'IN' 'DCM' 'PROJECT' QNAME 'STARTS' 'WITH' STRING
    | 'SHOW' 'ENTITIES' 'IN' 'DCM' 'PROJECT' QNAME 'LIMIT' NUMBER 'FROM' STRING
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-event-routing-table-on-organization
showEventRoutingTableOnOrganization
    : 'SHOW' 'EVENT' 'ROUTING' 'TABLE' 'ON' 'ORGANIZATION' 'FOR' 'ALL' 'APPLICATION' 'LISTINGS'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-event-routing-tables
showEventRoutingTables
    : 'SHOW' 'EVENT' 'ROUTING' 'TABLES'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-event-tables
showEventTables
    : 'SHOW' 'TERSE'? 'EVENT' 'TABLES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'SCHEMA' QNAME?))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-experiments
showExperiments
    : 'SHOW' 'EXPERIMENTS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'SCHEMA' QNAME?))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-external-agents
showExternalAgents
    : 'SHOW' 'EXTERNAL' 'AGENTS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'SCHEMA' QNAME?))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-external-consumer-programmatic-access-tokens
showExternalConsumerProgrammaticAccessTokens
    : 'SHOW' 'USER' 'PROGRAMMATIC' 'ACCESS' 'TOKENS' 'FOR' 'EXTERNAL' 'CONSUMER' QNAME
    | 'SHOW' 'USER' 'PATS' 'FOR' 'EXTERNAL' 'CONSUMER' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-external-consumers
showExternalConsumers
    : 'SHOW' 'EXTERNAL' 'CONSUMERS' ('LIKE' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-external-functions
showExternalFunctions
    : 'SHOW' 'EXTERNAL' 'FUNCTIONS' ('LIKE' STRING)? ('IN' ('APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-external-tables
// fix: read: the IN list's <schema_name> line has no | before APPLICATION; one is read there
showExternalTables
    : 'SHOW' 'TERSE'? 'EXTERNAL' 'TABLES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-external-volumes
showExternalVolumes
    : 'SHOW' 'EXTERNAL' 'VOLUMES' ('LIKE' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-failover-groups
showFailoverGroups
    : 'SHOW' 'FAILOVER' 'GROUPS' ('IN' 'ACCOUNT' ID)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-feature-policies
showFeaturePolicies
    : 'SHOW' 'FEATURE' 'POLICIES' ('IN' ('ACCOUNT' | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID | 'DATABASE' ID | 'SCHEMA' QNAME))?
    | 'SHOW' 'FEATURE' 'POLICIES' 'ON' 'ACCOUNT'
    | 'SHOW' 'FEATURE' 'POLICIES' 'ON' 'APPLICATION' ID
    | 'SHOW' 'FEATURE' 'POLICIES' 'ON' 'DATABASE' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-file-formats
// fix: read: the IN list's <schema_name> line has no | before APPLICATION; one is read there
showFileFormats
    : 'SHOW' 'FILE' 'FORMATS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-files-in-artifact-repository
showFilesInArtifactRepository
    : 'SHOW' 'FILES' ('LIKE' STRING)? 'IN' 'ARTIFACT' 'REPOSITORY' QNAME ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-functions
showFunctions
    : 'SHOW' 'FUNCTIONS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'CLASS' ID | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-functions-in-model
showFunctionsInModel
    : 'SHOW' 'FUNCTIONS' ('LIKE' STRING)? 'IN' 'MODEL' QNAME ('VERSION' ID)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-gateways
showGateways
    : 'SHOW' 'GATEWAYS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-git-branches
showGitBranches
    : 'SHOW' 'GIT' 'BRANCHES' ('LIKE' STRING)? 'IN' ('GIT' 'REPOSITORY')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-git-repositories
showGitRepositories
    : 'SHOW' 'GIT' 'REPOSITORIES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME))? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-git-tags
showGitTags
    : 'SHOW' 'GIT' 'TAGS' ('LIKE' STRING)? 'IN' ('GIT' 'REPOSITORY')? QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-global-accounts
showGlobalAccounts
    : 'SHOW' 'GLOBAL' 'ACCOUNTS' ('LIKE' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-grants
showGrants
    : 'SHOW' 'GRANTS' ('LIMIT' NUMBER)?
    | 'SHOW' 'GRANTS' 'ON' 'ACCOUNT' ('LIMIT' NUMBER)?
    | 'SHOW' 'GRANTS' 'ON' OBJECT_TYPE QNAME ('LIMIT' NUMBER)?
    | 'SHOW' 'GRANTS' 'TO' ('APPLICATION' ID | 'APPLICATION' 'ROLE' (ID '.')? ID | 'SERVICE' 'ROLE' QNAME '!' QNAME | ID 'ROLE' ID '!' ID | 'ROLE' ID | 'SHARE' ID ('IN' 'APPLICATION' 'PACKAGE' ID)? | 'USER' ID) ('LIMIT' NUMBER)?
    | 'SHOW' 'GRANTS' 'OF' ('APPLICATION' 'ROLE' ID | 'SERVICE' 'ROLE' QNAME '!' QNAME | 'ROLE' ID) ('LIMIT' NUMBER)?
    | 'SHOW' 'GRANTS' 'OF' 'SHARE' ID ('LIMIT' NUMBER)?
    | 'SHOW' 'FUTURE' 'GRANTS' 'IN' 'SCHEMA' QNAME ('LIMIT' NUMBER)?
    | 'SHOW' 'FUTURE' 'GRANTS' 'IN' 'DATABASE' ID ('LIMIT' NUMBER)?
    | 'SHOW' 'FUTURE' 'GRANTS' 'TO' 'ROLE' ID ('LIMIT' NUMBER)?
    | 'SHOW' 'FUTURE' 'GRANTS' 'TO' 'DATABASE' 'ROLE' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-grants-in-dcm-project
showGrantsInDcmProject
    : 'SHOW' 'GRANTS' 'IN' 'DCM' 'PROJECT' QNAME ('LIMIT' NUMBER)?
    | 'SHOW' 'FUTURE' 'GRANTS' 'IN' 'DCM' 'PROJECT' QNAME ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-hybrid-tables
showHybridTables
    : 'SHOW' 'TERSE'? 'HYBRID'? 'TABLES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'SCHEMA' QNAME?))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-iceberg-tables
showIcebergTables
    : 'SHOW' 'TERSE'? 'ICEBERG'? 'TABLES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-image-repositories
showImageRepositories
    : 'SHOW' 'IMAGE' 'REPOSITORIES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-images-in-image-repository
showImagesInImageRepository
    : 'SHOW' 'IMAGES' 'IN' 'IMAGE' 'REPOSITORY' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-indexes
showIndexes
    : 'SHOW' 'TERSE'? 'INDEXES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'SCHEMA' QNAME? | 'TABLE' | 'TABLE' QNAME))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-indexes-standard-tables
showIndexesStandardTables
    : 'SHOW' 'INDEXES' 'IN' 'TABLE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-ingress-policies
showIngressPolicies
    : 'SHOW' 'INGRESS' 'POLICIES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME) | 'ON' ('ACCOUNT' | 'USER' ID))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-integrations
showIntegrations
    : 'SHOW' (('API' | 'CATALOG' | 'EXTERNAL' 'ACCESS' | 'NOTIFICATION' | 'SECURITY' | 'STORAGE'))? 'INTEGRATIONS' ('LIKE' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-interactive-tables
showInteractiveTables
    : 'SHOW' 'INTERACTIVE' 'TABLES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-join-policies
showJoinPolicies
    : 'SHOW' 'JOIN' 'POLICIES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'SCHEMA' QNAME?))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-listings
showListings
    : 'SHOW' 'LISTINGS' ('LIKE' STRING)? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-listings-in-failover-group
showListingsInFailoverGroup
    : 'SHOW' 'LISTINGS' 'IN' 'FAILOVER' 'GROUP' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-locks
showLocks
    : 'SHOW' 'LOCKS' ('IN' 'ACCOUNT')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-maintenance-policies
// fix: the parameters list the entity types: DATABASE, APPLICATION PACKAGE and SCHEMA
showMaintenancePolicies
    : 'SHOW' 'MAINTENANCE' 'POLICIES' ('ON' | 'IN') ('ACCOUNT' | 'APPLICATION' ID | ('DATABASE' | 'APPLICATION' 'PACKAGE' | 'SCHEMA') ID)
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-managed-accounts
showManagedAccounts
    : 'SHOW' 'MANAGED' 'ACCOUNTS' ('LIKE' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-masking-policies
// fix: read: the IN list's <schema_name> line has no | before APPLICATION; one is read there
showMaskingPolicies
    : 'SHOW' 'MASKING' 'POLICIES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-materialized-views
// fix: read: the IN list's <schema_name> line has no | before APPLICATION; one is read there
showMaterializedViews
    : 'SHOW' 'MATERIALIZED' 'VIEWS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-mcp-servers
showMcpServers
    : 'SHOW' 'MCP' 'SERVERS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-mfa-methods
showMfaMethods
    : 'SHOW' 'MFA' 'METHODS' ('FOR' 'USER' ID)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-model-monitors
// fix: LIKE takes a quoted pattern, as the page's parameters and examples write it
showModelMonitors
    : 'SHOW' 'MODEL' 'MONITORS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME | 'GATEWAY' QNAME))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-models
showModels
    : 'SHOW' 'MODELS' ('LIKE' STRING)? ('IN' ('DATABASE' ID? | 'SCHEMA' QNAME?))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-multi-party-approval-policies
showMultiPartyApprovalPolicies
    : 'SHOW' 'MULTI' 'PARTY' 'APPROVAL' 'POLICIES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER)?
    | 'SHOW' 'MULTI' 'PARTY' 'APPROVAL' 'POLICIES' 'ON' 'ACCOUNT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-network-policies
showNetworkPolicies
    : 'SHOW' 'NETWORK' 'POLICIES'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-network-rules
showNetworkRules
    : 'SHOW' 'NETWORK' 'RULES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'SCHEMA'? QNAME?))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-nodes-compute-pool
showNodesComputePool
    : 'SHOW' 'NODES' 'IN' 'COMPUTE' 'POOL' QNAME ('LIKE' STRING)? ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-notebook-projects
showNotebookProjects
    : 'SHOW' 'NOTEBOOK' 'PROJECTS'
    | 'SHOW' 'NOTEBOOK' 'PROJECTS' 'IN' 'SCHEMA' ID '.' ID
    | 'SHOW' 'NOTEBOOK' 'PROJECTS' 'IN' 'DATABASE' ID
    | 'SHOW' 'NOTEBOOK' 'PROJECTS' 'IN' 'ACCOUNT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-code-bundles
showCodeBundles
    : 'SHOW' 'CODE' 'BUNDLES'
    | 'SHOW' 'CODE' 'BUNDLES' 'IN' 'SCHEMA' ID '.' ID
    | 'SHOW' 'CODE' 'BUNDLES' 'IN' 'DATABASE' ID
    | 'SHOW' 'CODE' 'BUNDLES' 'IN' 'ACCOUNT'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-notebooks
showNotebooks
    : 'SHOW' 'NOTEBOOKS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-notification-integrations
showNotificationIntegrations
    : 'SHOW' 'NOTIFICATION' 'INTEGRATIONS' ('LIKE' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-objects
// fix: read: the IN list's <schema_name> line has no | before APPLICATION; one is read there
showObjects
    : 'SHOW' 'TERSE'? 'OBJECTS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-objects-owned-by-application
showObjectsOwnedByApplication
    : 'SHOW' 'OBJECTS' 'OWNED' 'BY' 'APPLICATION' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-offers
showOffers
    : 'SHOW' 'OFFERS' ('LIKE' STRING)? 'IN' 'LISTING' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-oflow-data-plane-integration
showOflowDataPlaneIntegration
    : 'SHOW' 'OPENFLOW' 'DATA' 'PLANE' 'INTEGRATIONS' ('LIKE' STRING)? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-online-feature-tables
showOnlineFeatureTables
    : 'SHOW' 'ONLINE' 'FEATURE' 'TABLES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-openflow-connector-definitions
showOpenflowConnectorDefinitions
    : 'SHOW' 'OPENFLOW' 'CONNECTOR' 'DEFINITIONS' ('LIKE' STRING)? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-openflow-connectors
showOpenflowConnectors
    : 'SHOW' 'OPENFLOW' 'CONNECTORS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'SCHEMA'? QNAME?))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-openflow-deployments
showOpenflowDeployments
    : 'SHOW' 'OPENFLOW' 'DEPLOYMENTS' ('LIKE' STRING)? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-openflow-runtimes
showOpenflowRuntimes
    : 'SHOW' 'OPENFLOW' 'RUNTIMES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'SCHEMA'? QNAME?))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-organization-accounts
showOrganizationAccounts
    : 'SHOW' 'ORGANIZATION' 'ACCOUNTS' ('LIKE' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-organization-contracts
showOrganizationContracts
    : 'SHOW' 'ORGANIZATION' 'CONTRACTS'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-organization-profiles
showOrganizationProfiles
    : 'SHOW' 'ORGANIZATION' 'PROFILES'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-organization-user-groups
showOrganizationUserGroups
    : 'SHOW' 'ORGANIZATION' 'USER' 'GROUPS'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-organization-users
showOrganizationUsers
    : 'SHOW' 'ORGANIZATION' 'USERS' ('IN' 'ORGANIZATION' 'USER' 'GROUP' ID)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-packages-in-artifact-repository
showPackagesInArtifactRepository
    : 'SHOW' 'PACKAGES' ('LIKE' STRING)? 'IN' 'ARTIFACT' 'REPOSITORY' QNAME ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-packages-policies
showPackagesPolicies
    : 'SHOW' 'PACKAGES' 'POLICIES' ('IN' ('SCHEMA' | 'SCHEMA' QNAME | QNAME))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-parameters
showParameters
    : 'SHOW' 'PARAMETERS' ('LIKE' STRING)? (('IN' | 'FOR') (('SESSION' | 'ACCOUNT') | ('USER' | 'WAREHOUSE' | 'DATABASE' | 'SCHEMA' | 'TASK' | 'PIPE') QNAME? | 'TABLE' QNAME?))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-password-policies
showPasswordPolicies
    : 'SHOW' 'PASSWORD' 'POLICIES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID) | 'ON' ('ACCOUNT' | 'USER' ID))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-pipes
// fix: read: the IN list's <schema_name> line has no | before APPLICATION; one is read there
showPipes
    : 'SHOW' 'PIPES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-postgres-instances
showPostgresInstances
    : 'SHOW' 'POSTGRES' 'INSTANCES' ('LIKE' STRING)? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-pricing-plans
showPricingPlans
    : 'SHOW' 'PRICING' 'PLANS' ('LIKE' STRING)? 'IN' 'LISTING' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-primary-keys
showPrimaryKeys
    : 'SHOW' 'TERSE'? 'PRIMARY' 'KEYS' ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'SCHEMA' QNAME? | 'TABLE' | 'TABLE'? QNAME))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-privacy-policies
showPrivacyPolicies
    : 'SHOW' 'PRIVACY' 'POLICIES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'SCHEMA' QNAME?))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-privileges
showPrivileges
    : 'SHOW' 'PRIVILEGES' 'IN' 'APPLICATION' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-procedures
// fix: read: the IN list's <schema_name> line has no | before APPLICATION; one is read there
showProcedures
    : 'SHOW' 'PROCEDURES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'CLASS' ID | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-projection-policies
showProjectionPolicies
    : 'SHOW' 'PROJECTION' 'POLICIES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'SCHEMA' QNAME?))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-references
showReferences
    : 'SHOW' 'REFERENCES' 'IN' 'APPLICATION' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-regions
showRegions
    : 'SHOW' 'REGIONS' ('LIKE' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-release-channels
showReleaseChannels
    : 'SHOW' 'RELEASE' 'CHANNELS' 'IN' 'APPLICATION' 'PACKAGE' ID
    | 'SHOW' 'RELEASE' 'CHANNELS' 'IN' 'LISTING' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-release-directives
showReleaseDirectives
    : 'SHOW' 'RELEASE' 'DIRECTIVES' ('LIKE' STRING)? 'IN' 'APPLICATION' 'PACKAGE' QNAME ('FOR' 'RELEASE' 'CHANNEL' ID)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-replication-accounts
showReplicationAccounts
    : 'SHOW' 'REPLICATION' 'ACCOUNTS' ('LIKE' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-replication-databases
showReplicationDatabases
    : 'SHOW' 'REPLICATION' 'DATABASES' ('LIKE' STRING)? ('WITH' 'PRIMARY' ID '.' ID)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-replication-groups
showReplicationGroups
    : 'SHOW' 'REPLICATION' 'GROUPS' ('IN' 'ACCOUNT' ID)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-resource-monitors
showResourceMonitors
    : 'SHOW' 'RESOURCE' 'MONITORS' ('LIKE' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-restricted-session-scopes
showRestrictedSessionScopes
    : 'SHOW' ('RESTRICTED' 'SESSION' 'SCOPES' | 'RSS') ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'SCHEMA' QNAME?))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-roles
showRoles
    : 'SHOW' 'TERSE'? 'ROLES' ('LIKE' STRING)? ('IN' 'CLASS' ID)? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-roles-in-service
showRolesInService
    : 'SHOW' 'ROLES' 'IN' 'SERVICE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-row-access-policies
// fix: read: the IN list's <schema_name> line has no | before APPLICATION; one is read there
showRowAccessPolicies
    : 'SHOW' 'ROW' 'ACCESS' 'POLICIES' ('LIKE' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-rules-in-event-routing-table
showRulesInEventRoutingTable
    : 'SHOW' 'RULES' 'IN' 'EVENT' 'ROUTING' 'TABLE' '(' QNAME ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-run-in-experiment
showRunInExperiment
    : 'SHOW' 'RUN' 'METRICS' ('LIKE' STRING)? 'IN' 'EXPERIMENT' QNAME ('RUN' ID)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    | 'SHOW' 'RUN' 'PARAMETERS' ('LIKE' STRING)? 'IN' 'EXPERIMENT' QNAME ('RUN' ID)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-runs-in-experiment
showRunsInExperiment
    : 'SHOW' 'RUNS' ('LIKE' STRING)? 'IN' 'EXPERIMENT' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-schema-variables
showSchemaVariables
    : 'SHOW' 'SCHEMA' 'VARIABLES' ('LIKE' STRING)? 'IN' 'SCHEMA' (QNAME | ID '.' ID)
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-schemas
// fix: the privileges to filter by are each object type's, listed on Access control privileges rather than here; the filter is left out
showSchemas
    : 'SHOW' 'TERSE'? 'SCHEMAS' 'HISTORY'? ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-secrets
showSecrets
    : 'SHOW' 'SECRETS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE'? ID | 'SCHEMA'? QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-semantic-dimensions
showSemanticDimensions
    : 'SHOW' 'SEMANTIC' 'DIMENSIONS' ('LIKE' STRING)? ('IN' (QNAME | 'ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' ID '.' ID))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-semantic-dimensions-for-metric
showSemanticDimensionsForMetric
    : 'SHOW' 'SEMANTIC' 'DIMENSIONS' ('LIKE' STRING)? 'IN' QNAME 'FOR' 'METRIC' ID ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-semantic-facts
showSemanticFacts
    : 'SHOW' 'SEMANTIC' 'FACTS' ('LIKE' STRING)? ('IN' (QNAME | 'ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' ID '.' ID))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-semantic-metrics
showSemanticMetrics
    : 'SHOW' 'SEMANTIC' 'METRICS' ('LIKE' STRING)? ('IN' (QNAME | 'ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' ID '.' ID))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-semantic-views
showSemanticViews
    : 'SHOW' 'TERSE'? 'SEMANTIC' 'VIEWS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-sequences
// fix: read: the IN list's <schema_name> line has no | before APPLICATION; one is read there
showSequences
    : 'SHOW' 'SEQUENCES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-service-containers-in-service
showServiceContainersInService
    : 'SHOW' 'SERVICE' 'CONTAINERS' 'IN' 'SERVICE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-service-instances-in-service
showServiceInstancesInService
    : 'SHOW' 'SERVICE' 'INSTANCES' 'IN' 'SERVICE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-service-volumes-in-service
showServiceVolumesInService
    : 'SHOW' 'SERVICE' 'VOLUMES' 'IN' 'SERVICE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-services
showServices
    : 'SHOW' 'JOB'? 'SERVICES' ('EXCLUDE' 'JOBS')? ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME | 'COMPUTE' 'POOL' ID))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)? ('OF' 'TYPE' ID (',' ID)*)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-session-policies
showSessionPolicies
    : 'SHOW' 'SESSION' 'POLICIES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID) | 'ON' ('ACCOUNT' | 'USER' ID))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-shared-content
showSharedContent
    : 'SHOW' 'SHARED' 'CONTENT' 'IN' 'APPLICATION' 'PACKAGE' ID 'FOR' 'VERSION' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-shares
showShares
    : 'SHOW' 'SHARES' ('LIKE' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-shares-in-failover-group
showSharesInFailoverGroup
    : 'SHOW' 'SHARES' 'IN' 'FAILOVER' 'GROUP' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-shares-in-replication-group
showSharesInReplicationGroup
    : 'SHOW' 'SHARES' 'IN' 'REPLICATION' 'GROUP' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-snapshot-policies
// fix: the ] closing the IN clause stands after LIMIT instead of after IN's }; STARTS WITH and LIMIT are optional on their own, as the parameters say
showSnapshotPolicies
    : 'SHOW' 'SNAPSHOT' 'POLICIES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-snapshot-sets
// fix: the ] closing the IN clause stands after LIMIT instead of after IN's }; STARTS WITH and LIMIT are optional on their own, as the parameters say
showSnapshotSets
    : 'SHOW' 'SNAPSHOT' 'SETS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-snapshots
showSnapshots
    : 'SHOW' 'SNAPSHOTS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-snapshots-in-snapshot-set
showSnapshotsInSnapshotSet
    : 'SHOW' 'SNAPSHOTS' 'IN' 'SNAPSHOT' 'SET' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-specifications
showSpecifications
    : 'SHOW' (('APPROVED' | 'DECLINED' | 'PENDING'))? 'SPECIFICATIONS' ('IN' 'APPLICATION' ID)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-stages
// fix: read: the IN list's <schema_name> line has no | before APPLICATION; one is read there
showStages
    : 'SHOW' 'STAGES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-storage-lifecycle-policies
showStorageLifecyclePolicies
    : 'SHOW' 'STORAGE' 'LIFECYCLE' 'POLICIES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-streamlits
// fix: [ LIMIT ... is not closed; closed
// fix: IN's SCHEMA and SCHEMA <schema_name> are two choices; the | between them is missing
showStreamlits
    : 'SHOW' 'TERSE'? 'STREAMLITS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME))? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-streams
showStreams
    : 'SHOW' 'TERSE'? 'STREAMS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'SCHEMA'? QNAME? | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-tables
// fix: read: the IN list's <schema_name> line has no | before APPLICATION; one is read there
showTables
    : 'SHOW' 'TERSE'? 'TABLES' 'HISTORY'? ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-tags
// fix: read: the IN list's <schema_name> line has no | before APPLICATION; one is read there
showTags
    : 'SHOW' 'TAGS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-task-graphs
showTaskGraphs
    : 'SHOW' 'TASK' 'GRAPHS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'SCHEMA'? QNAME?))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-tasks
showTasks
    : 'SHOW' 'TERSE'? 'TASKS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'SCHEMA'? QNAME? | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))? ('STARTS' 'WITH' STRING)? ('ROOT' 'ONLY')? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-telemetry-event-definitions
showTelemetryEventDefinitions
    : 'SHOW' 'TELEMETRY' 'EVENT' 'DEFINITIONS' 'IN' 'APPLICATION' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-transactions
showTransactions
    : 'SHOW' 'TRANSACTIONS' ('IN' 'ACCOUNT')?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-types
// fix: read: the IN list's <schema_name> line has no | before APPLICATION; one is read there
showTypes
    : 'SHOW' 'TYPES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))? ('STARTS' 'WITH' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-user-functions
// fix: read: the IN list's <schema_name> line has no | before APPLICATION; one is read there
showUserFunctions
    : 'SHOW' 'USER' 'FUNCTIONS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-user-key-pairs
showUserKeyPairs
    : 'SHOW' 'USER' 'KEY' 'PAIRS' ('FOR' 'USER' ID)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-user-procedures
// fix: read: the IN list's <schema_name> line has no | before APPLICATION; one is read there
showUserProcedures
    : 'SHOW' 'USER' 'PROCEDURES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-user-programmatic-access-tokens
showUserProgrammaticAccessTokens
    : 'SHOW' 'USER' ('PROGRAMMATIC' 'ACCESS' 'TOKENS' | 'PATS') ('FOR' 'USER' ID)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-user-workload-identity-authentication-methods
showUserWorkloadIdentityAuthenticationMethods
    : 'SHOW' 'USER' 'WORKLOAD' 'IDENTITY' 'AUTHENTICATION' 'METHODS' ('FOR' 'USER' ID)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-users
showUsers
    : 'SHOW' 'TERSE'? 'USERS' ('LIKE' STRING)? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-variables
showVariables
    : 'SHOW' 'VARIABLES' ('LIKE' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-versions
// fix: LIKE takes a quoted pattern, as the page's parameters write it
showVersions
    : 'SHOW' 'VERSIONS' ('LIKE' STRING)? 'IN' 'APPLICATION' 'PACKAGE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-artifact-repository
showVersionsInArtifactRepository
    : 'SHOW' 'VERSIONS' ('LIKE' STRING)? 'IN' 'ARTIFACT' 'REPOSITORY' QNAME 'FOR' 'PACKAGE' ID ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-dataset
showVersionsInDataset
    : 'SHOW' 'VERSIONS' ('LIKE' STRING)? 'IN' 'DATASET' QNAME ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-dbt-project
showVersionsInDbtProject
    : 'SHOW' 'VERSIONS' 'IN' 'DBT' 'PROJECT' QNAME ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-listing
showVersionsInListing
    : 'SHOW' 'VERSIONS' 'IN' 'LISTING' QNAME ('LIMIT' NUMBER)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-model
showVersionsInModel
    : 'SHOW' 'VERSIONS' ('LIKE' STRING)? 'IN' 'MODEL' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-openflow-connector
showVersionsInOpenflowConnector
    : 'SHOW' 'VERSIONS' 'IN' 'OPENFLOW' 'CONNECTOR' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-organization-profile
showVersionsInOrganizationProfile
    : 'SHOW' 'VERSIONS' 'IN' 'ORGANIZATION' 'PROFILE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-views
showViews
    : 'SHOW' 'TERSE'? 'VIEWS' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' ID? | 'SCHEMA'? QNAME? | 'APPLICATION' ID | 'APPLICATION' 'PACKAGE' ID))? ('STARTS' 'WITH' STRING)? ('LIMIT' NUMBER ('FROM' STRING)?)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-warehouses
// fix: the privileges to filter by are each object type's, listed on Access control privileges rather than here; the filter is left out
showWarehouses
    : 'SHOW' 'WAREHOUSES' ('LIKE' STRING)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/show-workspaces
// fix: [ IN { ... } is not closed; closed
showWorkspaces
    : 'SHOW' 'WORKSPACES' ('LIKE' STRING)? ('IN' ('ACCOUNT' | 'DATABASE' | 'DATABASE' ID | 'SCHEMA' | 'SCHEMA' QNAME | QNAME))?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/truncate-materialized-view
truncateMaterializedView
    : 'TRUNCATE' 'MATERIALIZED' 'VIEW' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/truncate-table
truncateTable
    : 'TRUNCATE' 'TABLE'? ('IF' 'EXISTS')? QNAME
    | 'TRUNCATE' 'TABLE'? ('IF' 'EXISTS')? 'ERROR_TABLE' '(' QNAME ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/undrop
// no syntax

// https://docs.snowflake.com/en/sql-reference/sql/undrop-account
undropAccount
    : 'UNDROP' 'ACCOUNT' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/undrop-agent
undropAgent
    : 'UNDROP' 'AGENT' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/undrop-database
undropDatabase
    : 'UNDROP' 'DATABASE' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/undrop-dynamic-table
undropDynamicTable
    : 'UNDROP' 'DYNAMIC' 'TABLE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/undrop-schema
undropSchema
    : 'UNDROP' 'SCHEMA' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/undrop-snapshot
undropSnapshot
    : 'UNDROP' 'SNAPSHOT' (QNAME | 'IDENTIFIER' '(' ID ')') ('RENAME' 'TO' ID)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/undrop-stage
undropStage
    : 'UNDROP' 'STAGE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/undrop-streamlit
undropStreamlit
    : 'UNDROP' 'STREAMLIT' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/undrop-tag
undropTag
    : 'UNDROP' 'TAG' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/undrop-type
undropType
    : 'UNDROP' 'TYPE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/undrop-variable
undropVariable
    : 'UNDROP' 'VARIABLE' QNAME
    ;

// https://docs.snowflake.com/en/sql-reference/sql/unset
unset
    : 'UNSET' ID
    | 'UNSET' '(' ID (',' ID)* ')'
    ;

// https://docs.snowflake.com/en/sql-reference/sql/update
// fix: the value assigned is an expression, as the examples show
update
    : 'UPDATE' QNAME 'SET' ID '=' EXPR (',' ID '=' EXPR)* ('FROM' QNAME)? ('WHERE' EXPR)?
    ;

// https://docs.snowflake.com/en/sql-reference/sql/use
// no syntax

// https://docs.snowflake.com/en/sql-reference/sql/use-database
useDatabase
    : 'USE' 'DATABASE'? ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/use-role
useRole
    : 'USE' 'ROLE' ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/use-schema
useSchema
    : 'USE' 'SCHEMA'? (ID '.')? ID
    ;

// https://docs.snowflake.com/en/sql-reference/sql/use-secondary-roles
useSecondaryRoles
    : 'USE' 'SECONDARY' 'ROLES' ('ALL' | 'NONE' | ID (',' ID)*)
    ;

// https://docs.snowflake.com/en/sql-reference/sql/use-warehouse
useWarehouse
    : 'USE' 'WAREHOUSE' ID
    ;
