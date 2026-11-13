package parser

import (
	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// OceanBase's own objects: tenants, resource units and pools, tablegroups and
// sequences. Their statements are a name and a list of `WORD [=] value`
// options, so they are written as alter_system.go writes ALTER SYSTEM: a table
// of forms, one per alternative, over the same matchers. Each parses into an
// ObjectDDLStmt that names the statement and its object.
//
// Ref: create_tenant_stmt, clone_tenant_stmt, create_standby_tenant_stmt,
// alter_tenant_stmt, drop_tenant_stmt, create_resource_stmt,
// alter_resource_stmt and drop_resource_stmt in sql_parser_mysql_mode.y
// (oceanbase v4.3.5_CE_BP6). Their rows in testdata/compat.tsv hold every
// alternative, and its near-misses, to the pinned server.
var (
	// sysObjName is the relation_name that names the statement's object.
	sysObjName = sysLeaf(func(p *Parser) bool {
		name, _, err := p.parseIdent()
		if err != nil {
			return false
		}
		p.sys.name = name
		return true
	})
	sysIfNotExists = sysOpt("IF NOT EXISTS")
	sysIfExists    = sysOpt("IF EXISTS")
	sysZoneList    = sysList(sysStr, sysOptComma) // zone_list
	sysStrList     = sysList(sysStr, sysComma)    // resource_pool_list
	sysIDList      = sysList(sysInt, sysComma)    // id_list
	// USER_VARIABLE: @name, or the @ and the quoted name after it, which reach
	// this parser as two tokens and parseVariableRef reads whole.
	sysUserVar = sysLeaf(func(p *Parser) bool {
		if !p.isUserVariable() {
			return false
		}
		_, err := p.parseVariableRef()
		return err == nil
	})
	// var_name: NAME_OB or unreserved_keyword_normal.
	sysVarName = sysLeaf(func(p *Parser) bool {
		if p.curIsNotNormalKeyword() {
			return false
		}
		_, _, err := p.parseIdent()
		return err == nil
	})
	// sys_var_and_val_list: var_name {TO | = | :=} {expr | ON | BINARY | DEFAULT}.
	// The expression is tried first: BINARY also opens one, as in x = BINARY 'a',
	// and a list keeps the first reading of an item that matches.
	sysVarAssigns = sysList(sysSeq(sysVarName, sysAlt("TO", sysTok('='), sysTok(tokAssign)), sysAlt(sysExpr, "ON", "BINARY", "DEFAULT")), sysComma)

	// tenant_option, and the list that may be empty before its first comma.
	sysTenantOption = sysAlt(
		sysSeq(sysAlt("LOGONLY_REPLICA_NUM", "REPLICA_NUM", "PROGRESSIVE_MERGE_NUM"), sysEq, sysInt),
		sysSeq("LOCALITY", sysEq, sysStr, sysOpt("FORCE")),
		sysSeq("PRIMARY_ZONE", sysEq, sysAlt("DEFAULT", "RANDOM", sysUserVar, sysNS)),
		sysSeq("RESOURCE_POOL_LIST", sysEq, sysTok('('), sysStrList, sysTok(')')),
		sysSeq(sysAlt("ENABLE_ARBITRATION_SERVICE", "ENABLE_EXTENDED_ROWID"), sysEq, sysBool),
		sysSeq("ZONE_LIST", sysEq, sysTok('('), sysZoneList, sysTok(')')),
		sysSeq(sysAlt("CHARSET", "CHARACTER SET"), sysEq, sysAlt(sysNameOB, sysStr, "BINARY", "ASCII")),
		sysSeq("COLLATE", sysEq, sysAlt(sysNameOB, sysStr)),
		"READ ONLY", "READ WRITE",
		sysSeq("COMMENT", sysEq, sysStr),
		sysSeq(sysOpt("DEFAULT"), "TABLEGROUP", sysEq, sysAlt("NULL", sysName)),
	)
	sysTenantOptions = sysLooseList(sysTenantOption)
	// clone_tenant_option. A clone takes exactly two, with a comma between:
	// clone_tenant_option_list is no list.
	sysCloneOption = sysSeq(sysAlt("RESOURCE_POOL", "UNIT"), sysEq, sysNS)
	// The same list with an option in it.
	sysSomeTenantOptions = sysSeq(sysOptComma, sysList(sysTenantOption, sysComma))

	sysUnitOptions = sysLooseList(sysSeq(sysAlt("MIN_CPU", "MIN_IOPS", "MAX_CPU", "MEMORY_SIZE", "MAX_IOPS", "IOPS_WEIGHT",
		"LOG_DISK_SIZE", "DATA_DISK_SIZE", "MAX_NET_BANDWIDTH", "NET_BANDWIDTH_WEIGHT"), sysEq, sysConf))
	sysPoolZoneList = sysSeq("ZONE_LIST", sysEq, sysTok('('), sysZoneList, sysTok(')'))
	sysParenStrs    = sysSeq(sysTok('('), sysStrList, sysTok(')'))
)

// A tablegroup's options and partitioning, and a sequence's options.
//
// Ref: create_tablegroup_stmt, alter_tablegroup_stmt, drop_tablegroup_stmt,
// create_sequence_stmt, alter_sequence_stmt and drop_sequence_stmt in
// sql_parser_mysql_mode.y (oceanbase v4.3.5_CE_BP6).
var (
	// sysObjFactor is the relation_factor that names the statement's object.
	sysObjFactor = sysLeaf(func(p *Parser) bool {
		ref, err := p.parseRelationFactor()
		if err != nil {
			return false
		}
		p.sys.schema, p.sys.name = ref.Schema, ref.Name
		return true
	})
	sysRelFactor = sysLeaf(func(p *Parser) bool {
		_, err := p.parseRelationFactor()
		return err == nil
	})
	// name_list: NAME_OB | unreserved_keyword | name_list [','] NAME_OB. The first
	// name may be a keyword, and each one after it is a NAME_OB, with a comma
	// before it or without.
	sysNameList = sysSeq(sysName, sysOpt(sysList(sysSeq(sysOptComma, sysNameOB), sysSeq())))

	sysTablegroupOption = sysAlt(
		sysSeq(sysAlt("TABLEGROUP_ID", "MAX_USED_PART_ID"), sysEq, sysInt),
		sysSeq("BINDING", sysEq, sysBool),
		sysSeq("SHARDING", sysEq, sysStr),
	)
	// tablegroup_option_list is right-recursive:
	//
	//	list: space_separated | option ',' list
	//
	// so options that end in a comma come first and a run with no comma ends
	// the list: `A, B C` parses and `A B, C` does not.
	sysTablegroupOptions = sysSeq(
		sysOpt(sysList(sysSeq(sysTablegroupOption, sysComma), sysSeq())),
		sysList(sysTablegroupOption, sysSeq()))

	sysEngineOption = sysOpt("ENGINE", sysTok('='), "INNODB")                                          // opt_engine_option
	sysRangeExpr    = sysAlt("MAXVALUE", sysSeq(sysTok('('), sysList(sysExpr, sysComma), sysTok(')'))) // range_partition_expr
	sysListExpr     = sysSeq(sysTok('('), sysAlt("DEFAULT", sysList(sysExpr, sysComma)), sysTok(')'))  // list_partition_expr
	sysRangeSubpart = sysSeq("SUBPARTITION", sysRelFactor, "VALUES LESS THAN", sysRangeExpr, sysEngineOption)
	sysListSubpart  = sysSeq("SUBPARTITION", sysRelFactor, "VALUES IN", sysListExpr, sysEngineOption)
	sysHashSubpart  = sysSeq("SUBPARTITION", sysRelFactor, sysEngineOption)
	// opt_subpartition_list.
	sysSubpartList = sysOpt(sysAlt(sysParens(sysHashSubpart), sysParens(sysRangeSubpart), sysParens(sysListSubpart)))
	sysRangeParts  = sysParens(sysSeq("PARTITION", sysRelFactor, "VALUES LESS THAN", sysRangeExpr, sysOpt("ID", sysInt), sysEngineOption, sysSubpartList))
	sysListParts   = sysParens(sysSeq("PARTITION", sysRelFactor, "VALUES IN", sysListExpr, sysOpt("ID", sysInt), sysEngineOption, sysSubpartList))
	// tg_subpartition_option. A tablegroup names no column: its keys are counted.
	sysTablegroupSubpart = sysOpt(sysAlt(
		sysSeq("SUBPARTITION BY RANGE", sysOpt("COLUMNS", sysInt), "SUBPARTITION TEMPLATE", sysParens(sysRangeSubpart)),
		sysSeq("SUBPARTITION BY LIST", sysOpt("COLUMNS", sysInt), "SUBPARTITION TEMPLATE", sysParens(sysListSubpart)),
		sysSeq("SUBPARTITION BY HASH", sysOpt("SUBPARTITIONS", sysInt)),
		sysSeq("SUBPARTITION BY KEY", sysInt, sysOpt("SUBPARTITIONS", sysInt)),
	))
	sysOptPartitions = sysOpt("PARTITIONS", sysInt)
	// opt_tg_partition_option, without its empty alternative.
	sysTablegroupPartition = sysAlt(
		sysSeq("PARTITION BY HASH", sysTablegroupSubpart, sysOptPartitions),
		sysSeq("PARTITION BY KEY", sysInt, sysTablegroupSubpart, sysOptPartitions),
		sysSeq("PARTITION BY RANGE", sysOpt("COLUMNS", sysInt), sysTablegroupSubpart, sysOptPartitions, sysRangeParts),
		sysSeq("PARTITION BY LIST", sysOpt("COLUMNS", sysInt), sysTablegroupSubpart, sysOptPartitions, sysListParts),
	)
	sysPartsOfEither = sysAlt(sysRangeParts, sysListParts) // opt_partition_range_or_list

	// simple_num and sequence_option.
	sysSimpleNum      = sysSeq(sysOpt(sysAlt(sysTok('+'), sysTok('-'))), sysNumber)
	sysSequenceOption = sysAlt(
		sysSeq(sysAlt("INCREMENT BY", "START WITH", "MAXVALUE", "MINVALUE", "CACHE"), sysSimpleNum),
		"NOMAXVALUE", "NOMINVALUE", "CYCLE", "NOCYCLE", "NOCACHE", "ORDER", "NOORDER", "RESTART",
	)
	sysSequenceOptions = sysOpt(sysList(sysSequenceOption, sysSeq()))
)

// sysParens is a comma-separated list of item in parentheses.
func sysParens(item sysPat) sysPat {
	return sysSeq(sysTok('('), sysList(item, sysComma), sysTok(')'))
}

// clusterForms holds, for each word an object statement opens with, the forms
// that follow it.
var clusterForms = map[string]sysPat{
	"CREATE": sysAlt(
		sysForm("TENANT", sysIfNotExists, sysObjName, sysTenantOptions,
			sysOpt(sysAlt(sysSeq("SET VARIABLES", sysVarAssigns), sysSeq("SET", sysVarAssigns), sysSeq("VARIABLES", sysVarAssigns)))),
		// clone_tenant_stmt.
		sysForm("TENANT", sysIfNotExists, sysObjName, "FROM", sysName, sysOpt("USING SNAPSHOT", sysName),
			"WITH", sysCloneOption, sysComma, sysCloneOption),
		sysForm("STANDBY TENANT", sysIfNotExists, sysObjName, sysOpt("LOG_RESTORE_SOURCE", sysEq, sysConf), sysTenantOptions),
		sysForm("RESOURCE UNIT", sysIfNotExists, sysObjName, sysUnitOptions),
		sysForm("RESOURCE POOL", sysIfNotExists, sysObjName, sysLooseList(sysAlt(
			sysSeq("UNIT", sysEq, sysNS), sysSeq("UNIT_NUM", sysEq, sysInt), sysPoolZoneList, sysSeq("REPLICA_TYPE", sysEq, sysStr)))),
		sysForm("TABLEGROUP", sysIfNotExists, sysObjName, sysOpt(sysTablegroupOptions), sysOpt(sysTablegroupPartition)),
		sysForm("SEQUENCE", sysIfNotExists, sysObjFactor, sysSequenceOptions),
	),
	"ALTER": sysAlt(
		sysForm("TENANT", sysObjName, "RENAME GLOBAL_NAME TO", sysName),
		sysForm("TENANT", sysObjName, sysAlt("LOCK", "UNLOCK")),
		// The rule lets both the options and the variables be empty, and its
		// action then refuses the statement: ALTER TENANT t and ALTER TENANT t SET
		// are 1064s.
		sysForm("TENANT", sysAlt("ALL", sysObjName), sysOptSet, sysAlt(
			sysSeq(sysSomeTenantOptions, sysOpt("VARIABLES", sysVarAssigns)),
			sysSeq("VARIABLES", sysVarAssigns))),
		sysForm("RESOURCE UNIT", sysObjName, sysUnitOptions),
		sysForm("RESOURCE POOL", "MERGE", sysParenStrs, "INTO", sysParenStrs),
		sysForm("RESOURCE POOL", sysObjName, "SPLIT INTO", sysParenStrs, "ON", sysTok('('), sysZoneList, sysTok(')')),
		sysForm("RESOURCE POOL", sysObjName, sysList(sysAlt(
			sysSeq("UNIT", sysEq, sysNS),
			sysSeq("UNIT_NUM", sysEq, sysInt, sysOpt("DELETE UNIT", sysEq, sysTok('('), sysIDList, sysTok(')'))),
			sysPoolZoneList), sysComma)),
		sysForm("RESOURCE TENANT", sysObjName, "UNIT_NUM", sysEq, sysInt, sysOpt("DELETE UNIT_GROUP", sysEq, sysTok('('), sysIDList, sysTok(')'))),
		sysForm("TABLEGROUP", sysObjName, "ADD", sysOpt("TABLE"), sysList(sysRelFactor, sysComma)),
		// alter_tablegroup_actions: [SET] options, with a comma between two runs.
		sysForm("TABLEGROUP", sysObjName, sysList(sysSeq(sysOptSet, sysList(sysTablegroupOption, sysSeq())), sysComma)),
		// alter_tg_partition_option.
		sysForm("TABLEGROUP", sysObjName, sysAlt(
			sysSeq("DROP PARTITION", sysAlt(sysNameList, sysSeq(sysTok('('), sysNameList, sysTok(')')))),
			sysSeq("ADD PARTITION", sysPartsOfEither),
			sysTablegroupPartition,
			sysSeq("REORGANIZE PARTITION", sysNameList, "INTO", sysPartsOfEither),
			sysSeq("TRUNCATE PARTITION", sysNameList),
		)),
		sysForm("SEQUENCE", sysObjFactor, sysSequenceOptions),
	),
	"DROP": sysAlt(
		sysForm("TENANT", sysIfExists, sysObjName, sysOpt(sysAlt("FORCE", "PURGE"))),
		sysForm("RESOURCE UNIT", sysIfExists, sysObjName),
		sysForm("RESOURCE POOL", sysIfExists, sysObjName),
		sysForm("TABLEGROUP", sysIfExists, sysObjName),
		sysForm("SEQUENCE", sysIfExists, sysObjFactor),
	),
}

// curIsNotNormalKeyword reports whether the current token is an unreserved keyword
// that unreserved_keyword_normal leaves out, written bare: the words a role may not
// be called (unreserved_keyword_ambiguous_roles), PASSWORD, ACCESS, and the two of
// unreserved_keyword_adb_compatible, MAXVALUE and RANGE. Each is a name where the
// grammar writes relation_name and none where it writes var_name: SET range = 1
// is a 1064, and so is CREATE TENANT t SET range = 1.
//
// Ref: unreserved_keyword and unreserved_keyword_for_role_name in
// sql_parser_mysql_mode.y (oceanbase v4.3.5_CE_BP6, line 25806).
func (p *Parser) curIsNotNormalKeyword() bool {
	if isAmbiguousRoleWord(p.cur.Type) {
		return true
	}
	for _, w := range []string{"PASSWORD", "ACCESS", "MAXVALUE", "RANGE"} {
		if p.isBareWord(w) {
			return true
		}
	}
	return false
}

// clusterObjectAhead reports whether the token after CREATE, ALTER or DROP
// opens one of these statements: TENANT, TABLEGROUP, SEQUENCE, STANDBY TENANT,
// or RESOURCE before UNIT, POOL or TENANT.
func (p *Parser) clusterObjectAhead() bool {
	switch {
	case p.isBareWord("TENANT"), p.isBareWord("TABLEGROUP"), p.isBareWord("SEQUENCE"):
		return true
	case p.isBareWord("STANDBY"):
		return p.nextIsBareWord("TENANT")
	case p.isBareWord("RESOURCE"):
		return p.nextIsBareWord("UNIT", "POOL", "TENANT")
	}
	return false
}

// parseClusterDDL parses a cluster statement after verb, the word it opened
// with, when one stands there.
func (p *Parser) parseClusterDDL(start int, verb string) (nodes.Node, bool, error) {
	if !p.clusterObjectAhead() {
		return nil, false, nil
	}
	object, err := p.matchSystemForms(clusterForms[verb])
	if err != nil {
		return nil, true, err
	}
	return &nodes.ObjectDDLStmt{
		Loc:    nodes.Loc{Start: start, End: p.prev.End},
		Action: verb + " " + object,
		Schema: p.sys.schema,
		Name:   p.sys.name,
	}, true, nil
}
