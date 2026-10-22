package catalog

// Facts about built-in functions that pg_proc.dat, the source of
// pg_proc_generated.go, does not hold. TestBuiltinProcExtrasMatchPostgres
// checks these tables against a PostgreSQL 18 server.

// builtinProcArgDefaults is pg_proc.pronargdefaults for each built-in
// function that has parameter defaults. PostgreSQL declares those defaults in
// system_functions.sql, which redefines the functions after bootstrap, so
// pg_proc.dat carries none of them.
//
// pg: src/backend/catalog/system_functions.sql
var builtinProcArgDefaults = map[uint32]int16{
	1177: 2, // jsonb_path_exists_tz(jsonb,jsonpath,jsonb,boolean)
	1179: 2, // jsonb_path_query_tz(jsonb,jsonpath,jsonb,boolean)
	1180: 2, // jsonb_path_query_array_tz(jsonb,jsonpath,jsonb,boolean)
	1268: 1, // parse_ident(text,boolean)
	2023: 2, // jsonb_path_query_first_tz(jsonb,jsonpath,jsonb,boolean)
	2030: 2, // jsonb_path_match_tz(jsonb,jsonpath,jsonb,boolean)
	2096: 1, // pg_terminate_backend(integer,bigint)
	2172: 1, // pg_backup_start(text,boolean)
	2307: 1, // pg_stat_reset_slru(text)
	2739: 1, // pg_backup_stop(boolean)
	3261: 1, // json_strip_nulls(json,boolean)
	3262: 1, // jsonb_strip_nulls(jsonb,boolean)
	3305: 1, // jsonb_set(jsonb,text[],jsonb,boolean)
	3436: 2, // pg_promote(boolean,integer)
	3464: 7, // make_interval(integer,integer,integer,integer,integer,integer,double precision)
	3577: 1, // pg_logical_emit_message(boolean,text,text,boolean)
	3578: 1, // pg_logical_emit_message(boolean,text,bytea,boolean)
	3579: 1, // jsonb_insert(jsonb,text[],jsonb,boolean)
	3775: 1, // pg_stat_reset_shared(text)
	3779: 2, // pg_create_physical_replication_slot(name,boolean,boolean)
	3782: 1, // pg_logical_slot_get_changes(name,pg_lsn,integer,text[])
	3783: 1, // pg_logical_slot_get_binary_changes(name,pg_lsn,integer,text[])
	3784: 1, // pg_logical_slot_peek_changes(name,pg_lsn,integer,text[])
	3785: 1, // pg_logical_slot_peek_binary_changes(name,pg_lsn,integer,text[])
	3786: 3, // pg_create_logical_replication_slot(name,name,boolean,boolean,boolean)
	3960: 1, // json_populate_record(anyelement,json,boolean)
	3961: 1, // json_populate_recordset(anyelement,json,boolean)
	4005: 2, // jsonb_path_exists(jsonb,jsonpath,jsonb,boolean)
	4006: 2, // jsonb_path_query(jsonb,jsonpath,jsonb,boolean)
	4007: 2, // jsonb_path_query_array(jsonb,jsonpath,jsonb,boolean)
	4008: 2, // jsonb_path_query_first(jsonb,jsonpath,jsonb,boolean)
	4009: 2, // jsonb_path_match(jsonb,jsonpath,jsonb,boolean)
	4350: 1, // normalize(text,text)
	4351: 1, // is_normalized(text,text)
	5054: 2, // jsonb_set_lax(jsonb,text[],jsonb,boolean,text)
	6212: 2, // random_normal(double precision,double precision)
}

// builtinProcArgNames is pg_proc.proargnames for each built-in function
// whose parameter names, like its defaults, only system_functions.sql
// declares.
//
// pg: src/backend/catalog/system_functions.sql
var builtinProcArgNames = map[uint32][]string{
	1177: {"target", "path", "vars", "silent"},                                             // jsonb_path_exists_tz(jsonb,jsonpath,jsonb,boolean)
	1179: {"target", "path", "vars", "silent"},                                             // jsonb_path_query_tz(jsonb,jsonpath,jsonb,boolean)
	1180: {"target", "path", "vars", "silent"},                                             // jsonb_path_query_array_tz(jsonb,jsonpath,jsonb,boolean)
	2023: {"target", "path", "vars", "silent"},                                             // jsonb_path_query_first_tz(jsonb,jsonpath,jsonb,boolean)
	2030: {"target", "path", "vars", "silent"},                                             // jsonb_path_match_tz(jsonb,jsonpath,jsonb,boolean)
	2172: {"label", "fast"},                                                                // pg_backup_start(text,boolean)
	3261: {"target", "strip_in_arrays"},                                                    // json_strip_nulls(json,boolean)
	3262: {"target", "strip_in_arrays"},                                                    // jsonb_strip_nulls(jsonb,boolean)
	3305: {"jsonb_in", "path", "replacement", "create_if_missing"},                         // jsonb_set(jsonb,text[],jsonb,boolean)
	3577: {"transactional", "prefix", "message", "flush"},                                  // pg_logical_emit_message(boolean,text,text,boolean)
	3578: {"transactional", "prefix", "message", "flush"},                                  // pg_logical_emit_message(boolean,text,bytea,boolean)
	3579: {"jsonb_in", "path", "replacement", "insert_after"},                              // jsonb_insert(jsonb,text[],jsonb,boolean)
	3775: {"target"},                                                                       // pg_stat_reset_shared(text)
	3960: {"base", "from_json", "use_json_as_text"},                                        // json_populate_record(anyelement,json,boolean)
	3961: {"base", "from_json", "use_json_as_text"},                                        // json_populate_recordset(anyelement,json,boolean)
	4005: {"target", "path", "vars", "silent"},                                             // jsonb_path_exists(jsonb,jsonpath,jsonb,boolean)
	4006: {"target", "path", "vars", "silent"},                                             // jsonb_path_query(jsonb,jsonpath,jsonb,boolean)
	4007: {"target", "path", "vars", "silent"},                                             // jsonb_path_query_array(jsonb,jsonpath,jsonb,boolean)
	4008: {"target", "path", "vars", "silent"},                                             // jsonb_path_query_first(jsonb,jsonpath,jsonb,boolean)
	4009: {"target", "path", "vars", "silent"},                                             // jsonb_path_match(jsonb,jsonpath,jsonb,boolean)
	5054: {"jsonb_in", "path", "replacement", "create_if_missing", "null_value_treatment"}, // jsonb_set_lax(jsonb,text[],jsonb,boolean,text)
	6212: {"mean", "stddev"},                                                               // random_normal(double precision,double precision)
}

// builtinOrderedSetAggs is pg_aggregate.aggkind and aggnumdirectargs for
// each built-in ordered-set ('o') and hypothetical-set ('h') aggregate.
// Every other built-in aggregate is a normal one ('n').
//
// pg: src/include/catalog/pg_aggregate.dat
var builtinOrderedSetAggs = map[uint32]struct {
	kind      byte
	numDirect int
}{
	3972: {'o', 1}, // percentile_disc(double precision,anyelement)
	3974: {'o', 1}, // percentile_cont(double precision,double precision)
	3976: {'o', 1}, // percentile_cont(double precision,interval)
	3978: {'o', 1}, // percentile_disc(double precision[],anyelement)
	3980: {'o', 1}, // percentile_cont(double precision[],double precision)
	3982: {'o', 1}, // percentile_cont(double precision[],interval)
	3984: {'o', 0}, // mode(anyelement)
	3986: {'h', 1}, // rank("any")
	3988: {'h', 1}, // percent_rank("any")
	3990: {'h', 1}, // cume_dist("any")
	3992: {'h', 1}, // dense_rank("any")
}

func init() {
	for i := range BuiltinProcs {
		if n, ok := builtinProcArgDefaults[BuiltinProcs[i].OID]; ok {
			BuiltinProcs[i].NArgDefaults = n
		}
		if names, ok := builtinProcArgNames[BuiltinProcs[i].OID]; ok {
			BuiltinProcs[i].ArgNames = names
		}
	}
}

// aggKind returns pg_aggregate.aggkind for an aggregate: 'n' normal, 'o'
// ordered-set, 'h' hypothetical-set.
func (c *Catalog) aggKind(proc *BuiltinProc) byte {
	if proc.OID >= FirstNormalObjectId {
		if up := c.userProcs[proc.OID]; up != nil && up.AggKind != 0 {
			return up.AggKind
		}
		return 'n'
	}
	return builtinAggKind(proc)
}

// builtinAggKind returns pg_aggregate.aggkind for a built-in aggregate.
func builtinAggKind(proc *BuiltinProc) byte {
	if agg, ok := builtinOrderedSetAggs[proc.OID]; ok {
		return agg.kind
	}
	return 'n'
}

// aggNumDirectArgs returns pg_aggregate.aggnumdirectargs for an ordered-set
// or hypothetical-set aggregate, or -1 for any other.
func (c *Catalog) aggNumDirectArgs(proc *BuiltinProc) int {
	if proc.OID >= FirstNormalObjectId {
		if up := c.userProcs[proc.OID]; up != nil && up.AggKind != 0 && up.AggKind != 'n' {
			return int(up.AggNumDirectArgs)
		}
		return -1
	}
	if agg, ok := builtinOrderedSetAggs[proc.OID]; ok {
		return agg.numDirect
	}
	return -1
}

// countStarOID returns the OID of count(*), the zero-argument count.
func (c *Catalog) countStarOID() uint32 {
	for _, p := range c.procByName["count"] {
		if p.NArgs == 0 && p.Kind == 'a' {
			return p.OID
		}
	}
	return 0
}
