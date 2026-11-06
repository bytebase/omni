package parser

import "testing"

// FLUSH took any word as an option, so `FLUSH foo` and TiDB's `FLUSH tidb` parsed, and
// it skipped what it could not place in WITH READ LOCK and FOR EXPORT. It also read
// MySQL's FOR CHANNEL, which OceanBase does not have, and none of RELAY LOG, QUERY CACHE
// or DES_KEY_FILE, which it does. flushFormCases holds OceanBase's flush_option list, the
// table forms and their edges.
//
// Every case was measured against OceanBase CE 4.3.5.6 (oceanbase/oceanbase-ce:4.3.5-lts)
// through PREPARE on 2026-11-05. TestNameTokenFormsOceanBaseOracle keeps it a recording.
var flushFormCases = []parserFormCase{
	// Words that are not options.
	{"FLUSH tidb", false},
	{"FLUSH foo", false},
	{"FLUSH", false},
	{"FLUSH LOCAL", false},
	{"FLUSH `hosts`", false},
	{"FLUSH 'hosts'", false},
	{"FLUSH COLUMNS", false},
	{"FLUSH RESOURCES", false},
	{"FLUSH QUERY_CACHE", false},
	{"FLUSH ERROR", false},
	{"FLUSH ENGINE", false},
	{"FLUSH BINARY", false},
	{"FLUSH RELAY", false},
	{"FLUSH LOGS RELAY", false},
	{"FLUSH QUERY", false},
	{"FLUSH QUERY HOSTS", false},
	{"FLUSH RELAY foo bar", false},
	{"FLUSH RELAY HOSTS foo", false},
	{"FLUSH RELAY LOGS foo", false},

	// MariaDB's options, which OceanBase does not have.
	{"FLUSH SSL", false},
	{"FLUSH THREADS", false},
	{"FLUSH SLAVE", false},
	{"FLUSH REPLICA", false},
	{"FLUSH MASTER", false},
	{"FLUSH SESSION STATUS", false},
	{"FLUSH GLOBAL STATUS", false},
	{"FLUSH user_statistics", false},
	{"FLUSH BINARY LOGS DELETE_DOMAIN_ID = (1)", false},
	{"FLUSH RELAY LOGS 'c'", false},
	{"FLUSH TABLES WITH READ LOCK AND DISABLE CHECKPOINT", false},

	// The options.
	{"FLUSH HOSTS", true},
	{"FLUSH LOGS", true},
	{"FLUSH PRIVILEGES", true},
	{"FLUSH PRIVILEGES;", true},
	{"FLUSH OPTIMIZER_COSTS", true},
	{"FLUSH STATUS", true},
	{"FLUSH USER_RESOURCES", true},
	{"FLUSH DES_KEY_FILE", true},
	{"FLUSH QUERY CACHE", true},
	{"FLUSH LOCAL QUERY CACHE", true},
	{"FLUSH ERROR LOGS", true},
	{"FLUSH ENGINE LOGS", true},
	{"FLUSH GENERAL LOGS", true},
	{"FLUSH SLOW LOGS", true},
	{"FLUSH BINARY LOGS", true},
	{"FLUSH RELAY LOGS", true},
	{"FLUSH LOCAL HOSTS", true},
	{"FLUSH LOCAL STATUS", true},
	{"FLUSH NO_WRITE_TO_BINLOG PRIVILEGES", true},
	{"FLUSH hosts, status", true},
	{"FLUSH LOGS, LOGS", true},
	{"FLUSH STATUS, STATUS", true},
	{"FLUSH LOGS, RELAY LOGS", true},
	{"FLUSH RELAY LOGS, RELAY LOGS", true},
	{"FLUSH BINARY LOGS, ENGINE LOGS, ERROR LOGS, GENERAL LOGS, SLOW LOGS, RELAY LOGS, HOSTS, PRIVILEGES, LOGS, STATUS, USER_RESOURCES, OPTIMIZER_COSTS", true},
	{"FLUSH HOSTS,", false},
	{"FLUSH , HOSTS", false},
	{"FLUSH HOSTS HOSTS", false},
	{"FLUSH status, tables", false},

	// No FOR CHANNEL: RELAY LOG takes a name instead, one that is no keyword.
	{"FLUSH RELAY LOGS FOR CHANNEL 'c'", false},
	{"FLUSH RELAY LOGS FOR CHANNEL c", false},
	{"FLUSH RELAY LOG foo", true},
	{"FLUSH RELAY LOG `foo`", true},
	{"FLUSH RELAY LOG 1a", true},
	{"FLUSH RELAY LOG _cp1250", true},
	{"FLUSH RELAY LOG foo, HOSTS", true},
	{"FLUSH RELAY LOG", false},
	{"FLUSH RELAY LOG 'foo'", false},
	{"FLUSH RELAY LOG a.b", false},
	{"FLUSH RELAY LOG @a", false},
	{"FLUSH RELAY LOG _utf8", false},
	{"FLUSH RELAY LOG 1", false},
	{"FLUSH RELAY LOG foo bar", false},
	{"FLUSH RELAY LOG select", false},
	{"FLUSH RELAY LOG status", false},
	{"FLUSH RELAY LOG hosts", false},
	{"FLUSH RELAY LOG relay", false},
	{"FLUSH RELAY LOG logs", false},
	{"FLUSH RELAY LOG tables", false},
	{"FLUSH RELAY LOG account", false},
	{"FLUSH RELAY LOG names", false},

	// The table forms.
	{"FLUSH TABLES", true},
	{"FLUSH TABLE", true},
	{"FLUSH TABLES t", true},
	{"FLUSH TABLE t", true},
	{"FLUSH TABLES t, u", true},
	{"FLUSH TABLES test.t", true},
	{"FLUSH TABLES `t`", true},
	{"FLUSH TABLES WITH READ LOCK", true},
	{"FLUSH TABLES t WITH READ LOCK", true},
	{"FLUSH TABLES t, u WITH READ LOCK", true},
	{"FLUSH TABLES t FOR EXPORT", true},
	{"FLUSH TABLE t FOR EXPORT", true},
	{"FLUSH TABLES test.t FOR EXPORT", true},
	{"FLUSH TABLES t, u FOR EXPORT", true},
	{"FLUSH LOCAL TABLES", true},
	{"FLUSH LOCAL TABLES t WITH READ LOCK", true},
	{"FLUSH NO_WRITE_TO_BINLOG TABLES t WITH READ LOCK", true},
	{"FLUSH TABLES t,", false},
	{"FLUSH TABLES 't'", false},
	{"FLUSH TABLES t.*", false},
	{"FLUSH TABLES t AS x", false},
	{"FLUSH TABLES t PARTITION (p0)", false},
	{"FLUSH TABLES FOR EXPORT", false},
	{"FLUSH TABLES WITH READ", false},
	{"FLUSH TABLES WITH", false},
	{"FLUSH TABLES t FOR", false},
	{"FLUSH TABLES t HOSTS", false},
	{"FLUSH TABLES, HOSTS", false},
	{"FLUSH TABLES t FOR EXPORT, u", false},
	{"FLUSH TABLES t WITH READ LOCK FOR EXPORT", false},
}

// TestFlushForms holds omni to flushFormCases.
func TestFlushForms(t *testing.T) {
	for _, tc := range flushFormCases {
		_, err := Parse(tc.sql)
		if got := err == nil; got != tc.accept {
			t.Errorf("Parse(%q) accepts=%v, OceanBase 4.3.5.6 accepts=%v: %v", tc.sql, got, tc.accept, err)
		}
	}
}
