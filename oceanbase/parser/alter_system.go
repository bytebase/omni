package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// ALTER SYSTEM and the cluster statements that share its rule.
//
// alter_system_stmt has some 125 alternatives, nearly all of them a run of
// fixed words with a few `WORD [=] value` clauses after it, and no alternative
// builds a tree a consumer reads: what a caller needs is which statement it
// is. So the rule is written here as a table, one entry per alternative in the
// grammar's order, of small matchers (sysPat) that follow the grammar's own
// shapes, and the node records the words that name the statement.
//
// Ref: alter_system_stmt, switchover_tenant_stmt, recover_tenant_stmt,
// transfer_partition_stmt and service_name_stmt in sql_parser_mysql_mode.y
// (oceanbase v4.3.5_CE_BP6, line 20299 on). The rule's rows in
// testdata/compat.tsv hold every alternative, and its near-misses, to the
// pinned server.

// sysPat is one piece of a statement's shape. It matches at the current token
// and calls k for the rest; where k refuses, it undoes its match and reports
// false. A pattern that can match in more than one way tries each, so a
// statement parses exactly when some reading of the table reaches its end.
type sysPat func(p *Parser, k func() bool) bool

// sysMatch is what a match keeps beside the scan position: the furthest token
// a matcher refused, for the error, and the statement's naming words.
type sysMatch struct {
	far     Token
	hasFar  bool
	action  string
	nameEnd int    // where the naming words of the form end
	schema  string // the object the statement names, where a form reads one
	name    string
}

// sysMiss records the current token as one no matcher took.
func (p *Parser) sysMiss() {
	if !p.sys.hasFar || p.cur.Loc > p.sys.far.Loc {
		p.sys.far, p.sys.hasFar = p.cur, true
	}
}

// sysLeaf makes a pattern of read, which consumes what it matches and reports
// whether it matched.
func sysLeaf(read func(p *Parser) bool) sysPat {
	return func(p *Parser, k func() bool) bool {
		s := p.saveScan()
		if !read(p) {
			p.restoreScan(s)
			p.sysMiss()
			return false
		}
		if k() {
			return true
		}
		p.restoreScan(s)
		return false
	}
}

// sysWord matches one of OceanBase's keywords, written bare.
func sysWord(word string) sysPat {
	return sysLeaf(func(p *Parser) bool {
		if !p.isBareWord(word) {
			return false
		}
		p.advance()
		return true
	})
}

// sysTok matches a token by its type.
func sysTok(typ int) sysPat {
	return sysLeaf(func(p *Parser) bool {
		if p.cur.Type != typ {
			return false
		}
		p.advance()
		return true
	})
}

// sysPats turns a pattern's parts into patterns: a string is its words in
// order, each a keyword.
func sysPats(parts []any) []sysPat {
	var ps []sysPat
	for _, part := range parts {
		switch v := part.(type) {
		case string:
			for _, w := range strings.Fields(v) {
				ps = append(ps, sysWord(w))
			}
		case sysPat:
			ps = append(ps, v)
		}
	}
	return ps
}

func sysRun(p *Parser, ps []sysPat, k func() bool) bool {
	if len(ps) == 0 {
		return k()
	}
	return ps[0](p, func() bool { return sysRun(p, ps[1:], k) })
}

// sysSeq matches its parts in order.
func sysSeq(parts ...any) sysPat {
	ps := sysPats(parts)
	return func(p *Parser, k func() bool) bool { return sysRun(p, ps, k) }
}

// sysAlt matches any one of its alternatives, each a string of words or a
// pattern.
func sysAlt(alts ...any) sysPat {
	ps := make([]sysPat, len(alts))
	for i, a := range alts {
		ps[i] = sysSeq(a)
	}
	return func(p *Parser, k func() bool) bool {
		for _, alt := range ps {
			if alt(p, k) {
				return true
			}
		}
		return false
	}
}

// sysOpt matches its parts in order, or nothing.
func sysOpt(parts ...any) sysPat {
	seq := sysSeq(parts...)
	return func(p *Parser, k func() bool) bool { return seq(p, k) || k() }
}

// sysList matches item, then sep and item for as long as both follow. It is a
// loop and not a recursion, so a long list costs no stack: a separator that no
// item follows is left for what comes next.
func sysList(item, sep sysPat) sysPat {
	matched := func() bool { return true }
	more := sysSeq(sep, item)
	return func(p *Parser, k func() bool) bool {
		s := p.saveScan()
		if !item(p, matched) {
			return false
		}
		for more(p, matched) {
		}
		if k() {
			return true
		}
		p.restoreScan(s)
		return false
	}
}

// sysLooseList matches a list whose rule may be empty before its first comma:
//
//	list: %empty | item | list ',' item
//
// so it takes nothing, `item, item` and `, item` alike, as CALIBRATION_INFO = (, 's')
// and ALTER SYSTEM ADD ZONE z, REGION r are taken.
func sysLooseList(item sysPat) sysPat {
	matched := func() bool { return true }
	more := sysSeq(sysTok(','), item)
	return func(p *Parser, k func() bool) bool {
		s := p.saveScan()
		item(p, matched)
		for more(p, matched) {
		}
		if k() {
			return true
		}
		p.restoreScan(s)
		return false
	}
}

// sysForm is one alternative of a statement rule: the words that name it, then
// its clauses.
func sysForm(name string, clauses ...any) sysPat {
	words, rest := sysSeq(name), sysSeq(clauses...)
	return func(p *Parser, k func() bool) bool {
		return words(p, func() bool {
			p.sys.action, p.sys.nameEnd = name, p.prev.End
			return rest(p, k)
		})
	}
}

// The terminals and the small rules the forms share.
var (
	sysEq = sysOpt(sysTok('=')) // opt_equal_mark
	// STRING_VALUE. The lexer joins strings that stand side by side into one, so
	// SERVER 'a' 'b' holds a single value.
	sysStr = sysLeaf(func(p *Parser) bool {
		if p.cur.Type != tokSCONST {
			return false
		}
		for p.cur.Type == tokSCONST {
			p.advance()
		}
		return true
	})
	// INTNUM: an integer of no more than 63 bits. A longer run of digits is a
	// DECIMAL_VAL.
	sysInt = sysLeaf(func(p *Parser) bool {
		if p.cur.Type != tokICONST || !p.cur.isInt64() {
			return false
		}
		p.advance()
		return true
	})

	// relation_name: NAME_OB or an unreserved keyword.
	sysName = sysLeaf(func(p *Parser) bool {
		_, _, err := p.parseIdent()
		return err == nil
	})
	// NAME_OB: an identifier that is no keyword.
	sysNameOB = sysLeaf(func(p *Parser) bool {
		_, _, err := p.parseNameOB()
		return err == nil
	})
	// relation_name_or_string: relation_name | STRING_VALUE | ALL.
	sysNS = sysAlt(sysName, sysStr, "ALL")
	// int_or_decimal: INTNUM | DECIMAL_VAL. A run of digits too long for an INTNUM
	// is a DECIMAL_VAL, so an integer of any width is one or the other:
	// MINVALUE -9223372036854775808 is '-' DECIMAL_VAL.
	sysNumber = sysAlt(sysTok(tokICONST), sysTok(tokFCONST))
	// text_string: STRING_VALUE | HEX_STRING_VALUE | BIN_STRING_VALUE.
	sysText = sysAlt(sysStr, sysTok(tokXCONST), sysTok(tokBCONST))
	// conf_const: a string, a number with or without a minus, a temporal
	// literal, TRUE, FALSE, NULL or a system variable.
	sysConf = sysLeaf(func(p *Parser) bool {
		switch p.cur.Type {
		case '-':
			p.advance()
			if p.cur.Type != tokICONST && p.cur.Type != tokFCONST {
				return false
			}
			p.advance()
			return true
		case tokAt2:
			_, err := p.parseVariableRef()
			return err == nil
		case tokXCONST, tokBCONST:
			return false
		}
		_, err := p.parseLiteralOnly()
		return err == nil
	})

	sysComma = sysTok(',')

	sysTenantName     = sysSeq("TENANT", sysEq, sysNS)                                  // tenant_name
	sysOptTenantName  = sysOpt(sysTenantName)                                           // opt_tenant_name
	sysTenantNames    = sysList(sysNS, sysComma)                                        // tenant_name_list
	sysTenantTuple    = sysSeq("TENANT", sysEq, sysTenantNames)                         // tenant_list_tuple
	sysOptTenantTuple = sysOpt(sysTenantTuple)                                          // opt_tenant_list_v2, opt_backup_tenant_list
	sysOptTenantList  = sysOpt("TENANT", sysTok('='), sysTenantNames)                   // opt_tenant_list
	sysIPPort         = sysAlt(sysSeq("SERVER", sysEq, sysStr), sysSeq("HOST", sysStr)) // ip_port
	sysZoneDesc       = sysSeq("ZONE", sysEq, sysNS)                                    // zone_desc
	sysOptZoneDesc    = sysOpt(sysZoneDesc)
	sysServerOrZone   = sysAlt(sysIPPort, sysZoneDesc)
	sysOptServerZone  = sysOpt(sysServerOrZone)
	sysLS             = sysSeq("LS", sysEq, sysInt)
	sysTabletID       = sysSeq("TABLET_ID", sysTok('='), sysInt)
	// opt_tenant_list_or_ls_or_tablet_id.
	sysFreezeScope    = sysOpt(sysAlt(sysSeq(sysTenantTuple, sysOpt(sysLS), sysOpt(sysTabletID)), sysTabletID))
	sysOptDescription = sysOpt("DESCRIPTION", sysEq, sysStr)
	sysOptBackupTo    = sysOpt("TO", sysEq, sysStr)
	sysOptBackupDest  = sysOpt("FROM", sysStr)
	sysOptUntil       = sysOpt(sysAlt(sysSeq("UNTIL TIME", sysTok('='), sysStr), sysSeq("UNTIL SCN", sysTok('='), sysInt)))
	sysOptEncryptKey  = sysOpt("ENCRYPTED BY", sysStr)
	sysOptCopyID      = sysOpt("COPY", sysInt)
	sysOptTenantInfo  = sysOpt(sysAlt(sysSeq("TENANT_ID", sysEq, sysInt), sysSeq("TENANT", sysEq, sysNS)))
	sysOptBackupCopy  = sysOpt("BACKUP_BACKUP_DEST", sysEq, sysStr)
	sysOptActive      = sysOpt("WITH ACTIVE")
	sysOptPaxos       = sysOpt("PAXOS_REPLICA_NUM", sysEq, sysInt)
	sysOptDataSource  = sysOpt("DATA_SOURCE", sysEq, sysStr)
	sysPath           = sysSeq("PATH", sysEq, sysStr) // opt_path_info, which is not optional
	sysOptSet         = sysOpt("SET")
	sysOptAccessInfo  = sysOpt("ACCESS_INFO", sysEq, sysStr)
	sysOptAttribute   = sysOpt("ATTRIBUTE", sysEq, sysStr)
	sysStorageFor     = sysSeq("FOR", sysNS)                                                   // opt_storage_use_for
	sysScopeType      = sysAlt(sysSeq("ZONE", sysEq, sysStr), sysSeq("REGION", sysEq, sysStr)) // opt_scope_type
	sysPolicy         = sysSeq("POLICY", sysEq, sysStr)
	sysServerList     = sysList(sysStr, sysComma)
	// cache_type. CLOG, ILOG, COLUMN_STAT, BLOCK_INDEX and BLOOM_FILTER are in
	// the rule and are no keywords of the lexer's, so no statement reaches them.
	sysCacheType = sysAlt("ALL", "LOCATION", "BLOCK", "ROW", "SCHEMA", "PLAN", "AUDIT", "PL", "PS", "LIB")
	// add_or_alter_zone_options.
	sysZoneOptions = sysLooseList(sysSeq(sysAlt("REGION", "IDC", "ZONE_TYPE"), sysEq, sysNS))
	// The clauses every parameter of alter_system_set_parameter_action ends with.
	sysConfigScope = sysOpt("SCOPE", sysTok('='), sysAlt("MEMORY", "SPFILE", "BOTH"))
	sysParamTail   = sysSeq(sysOpt("COMMENT", sysStr), sysConfigScope, sysOptServerZone, sysOptTenantName)
	sysSetParam    = sysAlt(
		sysSeq(sysAlt(sysNameOB, "TABLET_SIZE", "CLUSTER_ID"), sysTok('='), sysConf, sysParamTail),
		sysSeq(sysAlt("ROOTSERVICE_LIST", "BACKUP_BACKUP_DEST", "OBCONFIG_URL", "LOG_DISK_SIZE", "LOG_RESTORE_SOURCE"), sysTok('='), sysStr, sysParamTail),
	)
	sysSetTPOption = sysAlt(
		sysSeq(sysAlt("TP_NO", "OCCUR", "FREQUENCY", "ERROR_CODE", "MATCH"), sysEq, sysInt),
		sysSeq("TP_NAME", sysEq, sysNS),
	)
	sysThrottleMetric = sysAlt(
		sysSeq(sysAlt("RT", "CPU", "NETWORK", "QUEUE_TIME"), sysTok('='), sysNumber),
		sysSeq(sysAlt("IO", "LOGICAL_READS"), sysTok('='), sysInt),
	)
	sysThrottleMetrics = sysList(sysThrottleMetric, sysSeq())

	// remap_table_val, without the alternatives of REMAP_TABLE_NAME and
	// REMAP_DATABASE_TABLE_NAME, tokens the lexer never makes.
	sysDotted   = sysSeq(sysName, sysTok('.'), sysName)
	sysDotStar  = sysSeq(sysName, sysTok('.'), sysTok('*'))
	sysRemapVal = sysAlt(
		sysSeq(sysDotStar, sysTok(':'), sysDotStar),
		sysSeq(sysDotted, sysTok(':'), sysName, sysTok(':'), sysDotted),
		sysSeq(sysDotted, sysTok(':'), sysName, sysTok(':'), sysName),
		sysSeq(sysDotted, sysTok(':'), sysDotted),
		sysSeq(sysDotted, sysTok(':'), sysName),
		sysSeq(sysName, sysTok(':'), sysName),
	)
	// recover_table_relation_name. A table is named by a relation_name or by a
	// word of mysql_reserved_keyword (table_relation_name).
	sysRecoverTable = sysAlt(
		sysSeq(sysTok('*'), sysTok('.'), sysTok('*')),
		sysRemapVal,
		sysDotStar,
		sysSeq(sysName, sysTok('.'), sysLeaf(func(p *Parser) bool {
			_, err := p.parseNameAfterDot()
			return err == nil
		})),
	)
	// restore_with_config_list and restore_remap_item_list: one or more, with
	// nothing between them.
	sysRestoreWith = sysList(sysSeq("WITH", sysAlt(sysSeq("KEY FROM", sysStr, sysOptEncryptKey), sysSeq("TENANT_STS_CREDENTIAL", sysStr))), sysSeq())
	sysRemapItems  = sysList(sysSeq("REMAP", sysAlt("TABLE", "TABLEGROUP", "TABLESPACE"), sysList(sysRemapVal, sysComma)), sysSeq())

	sysReplica  = sysSeq(sysLS, "SERVER", sysEq, sysStr)
	sysPartInfo = sysSeq("TABLE_ID", sysEq, sysInt, sysComma, "OBJECT_ID", sysEq, sysInt) // part_info
)

// sysBackup builds the four BACKUP forms that differ by INCREMENTAL and PLUS
// ARCHIVELOG, of a database or of tenants.
func sysBackup(name string, target sysPat) []sysPat {
	return []sysPat{
		sysForm(name, target, sysOptBackupTo, "PLUS ARCHIVELOG", sysOptDescription),
		sysForm(name, target, sysOptBackupTo, sysOptDescription),
	}
}

// alterSystemForms is alter_system_stmt after ALTER SYSTEM, with the four
// statements that begin with the same two words.
var alterSystemForms = func() sysPat {
	forms := []sysPat{
		sysForm("BOOTSTRAP",
			sysList(sysSeq(sysOpt("REGION", sysEq, sysNS), "ZONE", sysEq, sysNS, "SERVER", sysEq, sysStr), sysComma),
			sysOpt(sysComma, "SHARED_STORAGE_INFO", sysEq, sysStr)),
		sysForm("FLUSH", sysCacheType, "CACHE",
			sysOpt("NAMESPACE", sysEq, sysStr),
			sysOpt(sysAlt(sysSeq("SQL_ID", sysEq, sysStr), sysSeq("SCHEMA_ID", sysEq, sysInt))),
			sysOpt("DATABASES", sysEq, sysStr), sysOptTenantList, sysOpt("GLOBAL")),
		sysForm("FLUSH SQL", sysCacheType, sysOptTenantList, sysOpt("GLOBAL")),
		sysForm("FLUSH KVCACHE", sysOptTenantName, sysOpt("CACHE", sysEq, sysNS)),
		sysForm("FLUSH SS_MICRO_CACHE", sysOptTenantName),
		sysForm("FLUSH DAG WARNINGS"),
		sysForm("FLUSH ILOGCACHE", sysOpt("FILE_ID", sysEq, sysInt)),
		sysForm("SWITCH REPLICA", sysAlt("LEADER", "FOLLOWER", "DEFAULT"),
			sysAlt(sysSeq(sysLS, sysIPPort, sysTenantName), sysSeq(sysIPPort, sysOptTenantName), sysSeq(sysZoneDesc, sysOptTenantName))),
		sysForm("SWITCH ROOTSERVICE", sysAlt("LEADER", "FOLLOWER"), sysServerOrZone),
		sysForm("REPORT REPLICA", sysOptServerZone),
		sysForm("RECYCLE REPLICA", sysOptServerZone),
		sysForm("START MERGE", sysZoneDesc),
		sysForm("SUSPEND MERGE", sysOptTenantTuple),
		sysForm("RESUME MERGE", sysOptTenantTuple),
		sysForm("SUSPEND RECOVERY", sysOptZoneDesc),
		sysForm("RESUME RECOVERY", sysOptZoneDesc),
		sysForm("CLEAR MERGE ERROR", sysOptTenantTuple),
		sysForm("ADD ARBITRATION SERVICE", sysStr),
		sysForm("REMOVE ARBITRATION SERVICE", sysStr),
		sysForm("REPLACE ARBITRATION SERVICE", sysStr, "WITH", sysStr),
		sysForm("CANCEL", sysOpt("PARTITION MIGRATION"), "TASK", sysStr),
		sysForm("MAJOR FREEZE", sysFreezeScope, sysOpt("REBUILD COLUMN GROUP")),
		sysForm("CHECKPOINT"),
		sysForm("MINOR FREEZE", sysFreezeScope, sysOpt("SERVER", sysEq, sysTok('('), sysServerList, sysTok(')')), sysOptZoneDesc),
		sysForm("CHECKPOINT SLOG", sysOptTenantInfo, sysIPPort),
		sysForm("CLEAR ROOTTABLE", sysOptTenantName),
		sysForm("ADD ZONE", sysNS, sysZoneOptions),
		sysForm("ADD SHARED_STORAGE_DEST", sysPath, sysOptAccessInfo, sysOptAttribute, sysStorageFor, sysScopeType),
		sysForm("DROP SHARED_STORAGE_DEST", sysPath, sysStorageFor, sysScopeType),
		sysForm("REFRESH SCHEMA", sysOptServerZone),
		sysForm("REFRESH MEMORY STAT", sysOptServerZone),
		sysForm("WASH MEMORY FRAGMENTATION", sysOptServerZone),
		sysForm("REFRESH IO CALIBRATION",
			sysOpt("STORAGE", sysEq, sysStr),
			sysOpt("CALIBRATION_INFO", sysEq, sysTok('('), sysLooseList(sysStr), sysTok(')')),
			sysOptServerZone),
		sysForm("SET_TP", sysLooseList(sysSetTPOption), sysOptServerZone),
		sysForm("CLEAR LOCATION CACHE", sysOptServerZone),
		sysForm("REMOVE BALANCE TASK", sysOptTenantList,
			sysOpt("ZONE", sysTok('='), sysList(sysStr, sysOpt(sysComma))),
			sysOpt("TYPE", sysEq, sysAlt("AUTO", "MANUAL", "ALL"))),
		sysForm("RELOAD GTS"),
		sysForm("RELOAD UNIT"),
		sysForm("RELOAD SERVER"),
		sysForm("RELOAD ZONE"),
		sysForm("MIGRATE UNIT", sysEq, sysInt, "DESTINATION", sysEq, sysStr),
		sysForm("CANCEL MIGRATE UNIT", sysInt),
		sysForm("ADD REPLICA", sysReplica, "REPLICA_TYPE", sysEq, sysStr, sysOptDataSource, sysOptPaxos, sysOptTenantName),
		sysForm("REMOVE REPLICA", sysReplica, sysOptPaxos, sysOptTenantName),
		sysForm("MIGRATE REPLICA", sysLS, "SOURCE", sysEq, sysStr, "DESTINATION", sysEq, sysStr, sysOptDataSource, sysOptTenantName),
		sysForm("MODIFY REPLICA", sysReplica, "REPLICA_TYPE", sysEq, sysStr, sysOptPaxos, sysOptTenantName),
		sysForm("MODIFY", sysLS, "PAXOS_REPLICA_NUM", sysEq, sysInt, sysOptTenantName),
		sysForm("CANCEL REPLICA TASK TASK_ID", sysEq, sysStr, sysOptTenantName),
		sysForm("UPGRADE VIRTUAL SCHEMA"),
		sysForm("RUN JOB", sysStr, sysOptServerZone),
		sysForm("BEGIN UPGRADE"),
		sysForm("END UPGRADE"),
		sysForm("BEGIN ROLLING UPGRADE"),
		sysForm("END ROLLING UPGRADE"),
		sysForm("RUN UPGRADE JOB", sysStr, sysOptTenantTuple),
		sysForm("STOP UPGRADE JOB"),
		sysForm("ENABLE SQL THROTTLE", sysOpt("FOR PRIORITY", sysTok(tokLessEq), sysInt), "USING", sysThrottleMetrics),
		sysForm("DISABLE SQL THROTTLE"),
		sysForm("SET DISK VALID", sysIPPort),
		sysForm("SET NETWORK BANDWIDTH REGION", sysNS, "TO", sysNS, sysConf),
		sysForm("ADD RESTORE SOURCE", sysStr),
		sysForm("CLEAR RESTORE SOURCE"),
		sysForm("RECOVER TABLE", sysList(sysRecoverTable, sysComma),
			"TO", sysOpt("TENANT", sysEq), sysNS, sysOptBackupDest, sysOptUntil,
			"WITH", sysStr, sysOptEncryptKey, sysOpt(sysRestoreWith), sysOpt(sysRemapItems), sysOptDescription),
		sysForm("RESTORE FROM", sysStr, sysOptUntil, "PREVIEW"),
		sysForm("RESTORE", sysName, sysOptBackupDest, sysOptUntil, "WITH", sysStr, sysOptEncryptKey, sysOpt(sysRestoreWith), sysOptDescription),
		sysForm("CHANGE TENANT", sysAlt(sysSeq("TENANT_ID", sysEq, sysInt), sysNS)),
		sysForm("DROP TABLES IN SESSION", sysInt),
		sysForm("REFRESH TABLES IN SESSION", sysInt),
		sysForm("ARCHIVELOG", sysOptTenantTuple, sysOptDescription),
		sysForm("NOARCHIVELOG", sysOptTenantTuple, sysOptDescription),
		sysForm("BACKUP KEY", sysOptTenantTuple, sysOptBackupTo, sysOptEncryptKey),
		sysForm("CANCEL BACKUP", sysOptTenantTuple),
		sysForm("BACKUP CLUSTER PARAMETERS TO", sysStr),
		sysForm("CANCEL RESTORE", sysName),
		sysForm("CANCEL RECOVER TABLE", sysName),
		sysForm("SUSPEND BACKUP"),
		sysForm("RESUME BACKUP"),
		sysForm("TRIGGER TTL", sysOptTenantTuple),
		sysForm("SUSPEND TTL", sysOptTenantTuple),
		sysForm("RESUME TTL", sysOptTenantTuple),
		sysForm("CANCEL TTL", sysOptTenantTuple),
		sysForm("VALIDATE DATABASE", sysOptCopyID),
		sysForm("VALIDATE BACKUPSET", sysInt, sysOptCopyID),
		sysForm("CANCEL VALIDATE", sysInt, sysOptCopyID),
		sysForm("CANCEL BACKUP BACKUPSET"),
		sysForm("CANCEL BACKUP BACKUPPIECE"),
		sysForm("CANCEL ALL BACKUP FORCE"),
		sysForm("DELETE BACKUPSET", sysInt, sysOptCopyID, sysOptTenantTuple, sysOptDescription),
		sysForm("DELETE BACKUPPIECE", sysInt, sysOptCopyID, sysOptTenantTuple, sysOptDescription),
		sysForm("DELETE OBSOLETE BACKUP", sysOptTenantTuple, sysOptDescription),
		sysForm("CANCEL DELETE BACKUP", sysOptTenantTuple, sysOptDescription),
		sysForm("ADD DELETE BACKUP", sysPolicy,
			sysOpt("RECOVERY_WINDOW", sysEq, sysStr), sysOpt("REDUNDANCY", sysEq, sysInt), sysOpt("BACKUP_COPIES", sysEq, sysInt), sysOptTenantTuple),
		sysForm("DROP DELETE BACKUP", sysPolicy, sysOptTenantTuple),
		sysForm("BACKUP BACKUPSET ALL NOT BACKED UP", sysInt, "TIMES", sysOptTenantInfo, sysOptBackupCopy),
		sysForm("BACKUP BACKUPSET ALL", sysOptTenantInfo, sysOptBackupCopy),
		sysForm("BACKUP BACKUPSET", sysEq, sysInt, sysOptTenantInfo, sysOptBackupCopy),
		sysForm("START BACKUP ARCHIVELOG"),
		sysForm("STOP BACKUP ARCHIVELOG"),
		sysForm("BACKUP BACKUPPIECE ALL NOT BACKED UP", sysInt, "TIMES", sysOptActive, sysOptTenantInfo, sysOptBackupCopy),
		sysForm("BACKUP BACKUPPIECE ALL", sysOptActive, sysOptTenantInfo, sysOptBackupCopy),
		sysForm("BACKUP BACKUPPIECE", sysEq, sysInt, sysOptActive, sysOptTenantInfo, sysOptBackupCopy),
		sysForm("LOAD MODULE DATA MODULE", sysEq, sysNS, sysTenantName, sysOpt("INFILE", sysEq, sysStr)),
		sysForm("CHECK MODULE DATA MODULE", sysEq, sysNS, sysTenantName, sysOpt("INFILE", sysEq, sysStr)),
		sysForm("TRIGGER STORAGE_CACHE_POLICY_EXECUTOR", sysOptTenantName),
		sysForm("REBUILD TABLET", sysLS, "TABLETS", sysEq, sysStr, "DESTINATION", sysEq, sysStr, "SOURCE", sysEq, sysStr, sysOptTenantName),
		sysForm("LOAD LICENSE", sysStr),
		sysForm("RESET", sysList(sysSeq(sysNameOB, sysConfigScope, sysOptTenantName), sysComma)),
		sysForm("CANCEL CLONE", sysName),

		// switchover_tenant_stmt, recover_tenant_stmt, transfer_partition_stmt
		// and service_name_stmt.
		sysForm("ACTIVATE STANDBY", sysOptTenantName, sysOpt("VERIFY")),
		sysForm("SWITCHOVER TO PRIMARY", sysOptTenantName, sysOpt("VERIFY")),
		sysForm("SWITCHOVER TO STANDBY", sysOptTenantName, sysOpt("VERIFY")),
		sysForm("RECOVER STANDBY", sysOptTenantName, sysAlt("UNTIL UNLIMITED", "CANCEL", sysOptUntil)),
		sysForm("TRANSFER PARTITION", sysPartInfo, "TO LS", sysInt, sysOptTenantName),
		sysForm("CANCEL TRANSFER PARTITION", sysAlt("ALL", sysPartInfo), sysOptTenantName),
		sysForm("CANCEL BALANCE JOB", sysOptTenantName),
	}
	// server_action SERVER and zone_action ZONE.
	for _, action := range []string{"ADD", "DELETE", "CANCEL DELETE", "START", "STOP", "FORCE STOP", "ISOLATE"} {
		forms = append(forms, sysForm(action+" SERVER", sysServerList, sysOptZoneDesc))
	}
	for _, action := range []string{"DELETE", "START", "STOP", "FORCE STOP", "ISOLATE"} {
		forms = append(forms, sysForm(action+" ZONE", sysNS))
	}
	// alter_or_change_or_modify.
	for _, verb := range []string{"ALTER", "CHANGE", "MODIFY"} {
		forms = append(forms,
			sysForm(verb+" EXTERNAL_STORAGE_DEST", sysPath, sysOptSet, sysOptAccessInfo, sysOptAttribute),
			sysForm(verb+" ZONE", sysNS, sysOptSet, sysZoneOptions),
			sysForm(verb+" SHARED_STORAGE_DEST", sysPath, sysOptSet, sysOptAccessInfo, sysOptAttribute))
	}
	// service_op SERVICE.
	for _, op := range []string{"CREATE", "DELETE", "START", "STOP"} {
		forms = append(forms, sysForm(op+" SERVICE", sysName, sysOptTenantName))
	}
	forms = append(forms, sysBackup("BACKUP INCREMENTAL DATABASE", sysSeq())...)
	forms = append(forms, sysBackup("BACKUP DATABASE", sysSeq())...)
	forms = append(forms, sysBackup("BACKUP INCREMENTAL", sysOptTenantTuple)...)
	forms = append(forms, sysBackup("BACKUP", sysOptTenantTuple)...)
	// opt_set alter_system_set_parameter_actions comes last: SET alone names
	// nothing, and every other form that opens with it is tried first.
	forms = append(forms,
		sysForm("SET", sysList(sysSetParam, sysComma)),
		sysForm("", sysList(sysSetParam, sysComma)))
	alts := make([]any, len(forms))
	for i, f := range forms {
		alts[i] = f
	}
	return sysAlt(alts...)
}()

// alterDiskgroupForms is alter_system_stmt's two ALTER DISKGROUP alternatives,
// after those two words.
var alterDiskgroupForms = sysAlt(
	sysForm("ADD DISK", sysStr, sysOpt("NAME", sysEq, sysNS), sysIPPort, sysOptZoneDesc),
	sysForm("DROP DISK", sysStr, sysIPPort, sysOptZoneDesc),
)

// setEncryptionForms is alter_system_stmt's two alternatives that open with
// SET, after that word.
var setEncryptionForms = sysAlt(
	sysForm("ENCRYPTION ON IDENTIFIED BY", sysStr, "ONLY"),
	sysForm("DECRYPTION IDENTIFIED BY", sysList(sysText, sysComma)),
)

// matchSystemForms matches forms at the current token through to the end of
// the statement, and reports the words that name the form it was. A statement
// that matches none is a syntax error at the furthest token no form took.
func (p *Parser) matchSystemForms(forms sysPat) (string, error) {
	p.sys = sysMatch{}
	atEnd := func() bool {
		if p.cur.Type == tokEOF || p.cur.Type == ';' {
			return true
		}
		p.sysMiss()
		return false
	}
	if !forms(p, atEnd) {
		return "", p.syntaxErrorAtTok(p.sys.far)
	}
	return p.sys.action, nil
}

// parseAlterSystemStmt parses ALTER SYSTEM at SYSTEM.
func (p *Parser) parseAlterSystemStmt(start int) (*nodes.AlterSystemStmt, error) {
	p.advance() // consume SYSTEM
	action, err := p.matchSystemForms(alterSystemForms)
	if err != nil {
		return nil, err
	}
	if action == "" {
		// ALTER SYSTEM name = value, where SET is optional.
		action = "SET"
	}
	return &nodes.AlterSystemStmt{Loc: nodes.Loc{Start: start, End: p.prev.End}, Action: action}, nil
}

// parseAlterDiskgroupStmt parses ALTER DISKGROUP at DISKGROUP.
func (p *Parser) parseAlterDiskgroupStmt(start int) (*nodes.AlterSystemStmt, error) {
	p.advance() // consume DISKGROUP
	if _, _, err := p.parseIdent(); err != nil {
		return nil, err
	}
	action, err := p.matchSystemForms(alterDiskgroupForms)
	if err != nil {
		return nil, err
	}
	return &nodes.AlterSystemStmt{Loc: nodes.Loc{Start: start, End: p.prev.End}, Action: "DISKGROUP " + action}, nil
}
