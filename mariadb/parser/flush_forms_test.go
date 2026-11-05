package parser

import "testing"

// FLUSH took any word as an option, so `FLUSH foo` and TiDB's `FLUSH tidb` parsed, and
// it skipped what it could not place in WITH READ LOCK, FOR EXPORT and FOR CHANNEL. It
// also took a string as a table name, a name as the channel, and none of MariaDB's own
// options past the first word. flushFormCases holds MariaDB's flush_option list, the table
// forms and their edges.
//
// Every case was measured against MariaDB 11.8.8 (mariadb:11.8.8) through PREPARE on
// 2026-11-05. TestNameTokenFormsOracle keeps it a recording. The three
// QUERY_RESPONSE_TIME tables answer only once the plugin is installed, so
// flushPluginFormCases holds them apart, measured after INSTALL SONAME
// 'query_response_time' on the same image.
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
	{"FLUSH OPTIMIZER_COSTS", false},
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
	{"FLUSH RELAY LOG foo", false},
	{"FLUSH THREAD_POOL_GROUPS", false},
	{"FLUSH QUERY_RESPONSE_TIME_AUDIT", false},

	// The options.
	{"FLUSH HOSTS", true},
	{"FLUSH LOGS", true},
	{"FLUSH PRIVILEGES", true},
	{"FLUSH PRIVILEGES;", true},
	{"FLUSH STATUS", true},
	{"FLUSH SESSION STATUS", true},
	{"FLUSH GLOBAL STATUS", true},
	{"FLUSH LOCAL STATUS", true},
	{"FLUSH USER_RESOURCES", true},
	{"FLUSH QUERY CACHE", true},
	{"FLUSH DES_KEY_FILE", true},
	{"FLUSH SSL", true},
	{"FLUSH THREADS", true},
	{"FLUSH MASTER", true},
	{"FLUSH SLAVE", true},
	{"FLUSH REPLICA", true},
	{"FLUSH SLAVE 'c'", true},
	{"FLUSH SLAVE c", false},
	{"FLUSH SLAVE FOR CHANNEL 'c'", false},
	{"FLUSH ERROR LOGS", true},
	{"FLUSH ENGINE LOGS", true},
	{"FLUSH GENERAL LOGS", true},
	{"FLUSH SLOW LOGS", true},
	{"FLUSH BINARY LOGS", true},
	{"FLUSH LOCAL HOSTS", true},
	{"FLUSH NO_WRITE_TO_BINLOG PRIVILEGES", true},
	{"FLUSH hosts, status", true},
	{"FLUSH LOGS, LOGS", true},
	{"FLUSH STATUS, STATUS", true},
	{"FLUSH MASTER, MASTER", true},
	{"FLUSH HOSTS,", false},
	{"FLUSH , HOSTS", false},
	{"FLUSH HOSTS HOSTS", false},
	{"FLUSH status, tables", false},

	// An information_schema table with a reset function.
	{"FLUSH user_statistics", true},
	{"FLUSH CLIENT_STATISTICS", true},
	{"FLUSH INDEX_STATISTICS", true},
	{"FLUSH TABLE_STATISTICS", true},
	{"FLUSH USER_VARIABLES", true},
	{"FLUSH THREAD_POOL_STATS", true},
	{"FLUSH THREAD_POOL_WAITS", true},
	{"FLUSH `user_statistics`", true},
	{"FLUSH user_statistics, hosts", true},
	{"FLUSH HOSTS, user_statistics, LOGS", true},

	// BINARY LOGS takes a list of domain ids, each at most 2^32-1.
	{"FLUSH BINARY LOGS DELETE_DOMAIN_ID = (1)", true},
	{"FLUSH BINARY LOGS DELETE_DOMAIN_ID = ()", true},
	{"FLUSH BINARY LOGS DELETE_DOMAIN_ID = (1, 2)", true},
	{"FLUSH BINARY LOGS DELETE_DOMAIN_ID = (+1)", true},
	{"FLUSH BINARY LOGS DELETE_DOMAIN_ID = (1.5)", true},
	{"FLUSH BINARY LOGS DELETE_DOMAIN_ID = (4294967295)", true},
	{"FLUSH BINARY LOGS DELETE_DOMAIN_ID = (4294967296)", false},           // 1076, while parsing
	{"FLUSH BINARY LOGS DELETE_DOMAIN_ID = (18446744073709551616)", false}, // 1076, while parsing
	{"FLUSH BINARY LOGS DELETE_DOMAIN_ID = (1,)", false},
	{"FLUSH BINARY LOGS DELETE_DOMAIN_ID = (-1)", false},
	{"FLUSH BINARY LOGS DELETE_DOMAIN_ID = (0x10)", false},
	{"FLUSH BINARY LOGS DELETE_DOMAIN_ID = ('1')", false},
	{"FLUSH BINARY LOGS DELETE_DOMAIN_ID := (1)", false},
	{"FLUSH BINARY LOGS DELETE_DOMAIN_ID = 1", false},
	{"FLUSH BINARY LOGS DELETE_DOMAIN_ID (1)", false},

	// RELAY LOGS takes a connection name or a channel, plain strings, not both.
	{"FLUSH RELAY LOGS", true},
	{"FLUSH RELAY LOGS 'c'", true},
	{"FLUSH RELAY LOGS \"c\"", true},
	{"FLUSH RELAY LOGS FOR CHANNEL 'c'", true},
	{"FLUSH RELAY LOGS FOR CHANNEL ''", true},
	{"FLUSH RELAY LOGS FOR CHANNEL 'a\nb'", true},
	{"FLUSH RELAY LOGS 'c' FOR CHANNEL 'd'", false}, // 1210, while parsing
	{"FLUSH RELAY LOGS FOR CHANNEL c", false},
	{"FLUSH RELAY LOGS FOR CHANNEL _utf8mb4'c'", false},
	{"FLUSH RELAY LOGS FOR CHANNEL N'c'", false},
	{"FLUSH RELAY LOGS FOR CHANNEL 'a' 'b'", false},

	// RELAY LOGS and SLAVE once each in a list.
	{"FLUSH RELAY LOGS, RELAY LOGS", false}, // 1221, while parsing
	{"FLUSH SLAVE, SLAVE", false},           // 1221, while parsing
	{"FLUSH RELAY LOGS, SLAVE", true},
	{"FLUSH LOGS, RELAY LOGS", true},
	{"FLUSH SLAVE 'c', RELAY LOGS", true},

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
	{"FLUSH TABLES WITH READ LOCK AND DISABLE CHECKPOINT", true},
	{"FLUSH TABLES t WITH READ LOCK AND DISABLE CHECKPOINT", true},
	{"FLUSH TABLES t FOR EXPORT", true},
	{"FLUSH TABLES t, u FOR EXPORT", true},
	{"FLUSH LOCAL TABLES", true},
	{"FLUSH LOCAL TABLES t WITH READ LOCK", true},
	{"FLUSH NO_WRITE_TO_BINLOG TABLES t WITH READ LOCK", true},
	{"FLUSH TABLES t,", false},
	{"FLUSH TABLES 't'", false},
	{"FLUSH TABLES t.*", false},
	{"FLUSH TABLES t AS x", false},
	{"FLUSH TABLES FOR EXPORT", false},
	{"FLUSH TABLES WITH READ", false},
	{"FLUSH TABLES WITH", false},
	{"FLUSH TABLES t WITH READ LOCK AND DISABLE", false},
	{"FLUSH TABLES t FOR", false},
	{"FLUSH TABLES t HOSTS", false},
	{"FLUSH TABLES, HOSTS", false},
	{"FLUSH TABLES t FOR EXPORT, u", false},
	{"FLUSH TABLES t WITH READ LOCK FOR EXPORT", false},
	{"FLUSH TABLES @a", false},
	{"FLUSH TABLES _utf8", false},
}

// flushPluginFormCases are the information_schema tables a plugin MariaDB ships adds:
// query_response_time's three, measured on mariadb:11.8.8 after INSTALL SONAME
// 'query_response_time' on 2026-11-05. The default image has no such plugin and answers
// 1064, which TestNameTokenFormsOracle would read as a refusal, so they are held to omni
// only.
var flushPluginFormCases = []parserFormCase{
	{"FLUSH QUERY_RESPONSE_TIME", true},
	{"FLUSH QUERY_RESPONSE_TIME_READ", true},
	{"FLUSH QUERY_RESPONSE_TIME_WRITE", true},
}

// TestFlushForms holds omni to flushFormCases and flushPluginFormCases.
func TestFlushForms(t *testing.T) {
	for _, cases := range [][]parserFormCase{flushFormCases, flushPluginFormCases} {
		for _, tc := range cases {
			_, err := Parse(tc.sql)
			if got := err == nil; got != tc.accept {
				t.Errorf("Parse(%q) accepts=%v, MariaDB 11.8.8 accepts=%v: %v", tc.sql, got, tc.accept, err)
			}
		}
	}
}
