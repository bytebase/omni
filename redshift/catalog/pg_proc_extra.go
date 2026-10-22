package catalog

// Facts about built-in functions that pg_proc.dat, the source of
// pg_proc_generated.go, does not hold. TestBuiltinProcExtrasMatchPostgres
// checks both tables against a PostgreSQL 18 server.

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

// builtinOrderedSetAggKinds is pg_aggregate.aggkind for each built-in
// ordered-set ('o') and hypothetical-set ('h') aggregate. Every other
// built-in aggregate is a normal one ('n').
//
// pg: src/include/catalog/pg_aggregate.dat
var builtinOrderedSetAggKinds = map[uint32]byte{
	3972: 'o', // percentile_disc(double precision,anyelement)
	3974: 'o', // percentile_cont(double precision,double precision)
	3976: 'o', // percentile_cont(double precision,interval)
	3978: 'o', // percentile_disc(double precision[],anyelement)
	3980: 'o', // percentile_cont(double precision[],double precision)
	3982: 'o', // percentile_cont(double precision[],interval)
	3984: 'o', // mode(anyelement)
	3986: 'h', // rank("any")
	3988: 'h', // percent_rank("any")
	3990: 'h', // cume_dist("any")
	3992: 'h', // dense_rank("any")
}

func init() {
	for i := range BuiltinProcs {
		if n, ok := builtinProcArgDefaults[BuiltinProcs[i].OID]; ok {
			BuiltinProcs[i].NArgDefaults = n
		}
	}
}

// builtinAggKind returns pg_aggregate.aggkind for a built-in aggregate, or 0
// for one created by CREATE AGGREGATE, whose kind the catalog does not
// record.
func builtinAggKind(proc *BuiltinProc) byte {
	if proc.OID >= FirstNormalObjectId {
		return 0
	}
	if kind, ok := builtinOrderedSetAggKinds[proc.OID]; ok {
		return kind
	}
	return 'n'
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
