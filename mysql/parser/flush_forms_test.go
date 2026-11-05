package parser

import "testing"

// FLUSH took any word as an option, so `FLUSH foo` and TiDB's `FLUSH tidb` parsed, and
// it skipped what it could not place in WITH READ LOCK, FOR EXPORT and FOR CHANNEL. It
// also took a string as a table name and a name as the channel. flushFormCases holds the
// flush_option list, the table forms and their edges.
//
// Every case was measured against MySQL 8.0 (mysql:8.0, 8.0.46) through PREPARE on
// 2026-11-05. TestFlushFormsOracle keeps it a recording.
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

	// MariaDB's options, which MySQL does not have.
	{"FLUSH QUERY CACHE", false},
	{"FLUSH DES_KEY_FILE", false},
	{"FLUSH SSL", false},
	{"FLUSH THREADS", false},
	{"FLUSH SLAVE", false},
	{"FLUSH MASTER", false},
	{"FLUSH SESSION STATUS", false},
	{"FLUSH GLOBAL STATUS", false},
	{"FLUSH user_statistics", false},
	{"FLUSH BINARY LOGS DELETE_DOMAIN_ID = (1)", false},
	{"FLUSH RELAY LOGS 'c'", false},
	{"FLUSH TABLES WITH READ LOCK AND DISABLE CHECKPOINT", false},
	{"FLUSH RELAY LOG foo", false},

	// The options.
	{"FLUSH HOSTS", true},
	{"FLUSH LOGS", true},
	{"FLUSH PRIVILEGES", true},
	{"FLUSH PRIVILEGES;", true},
	{"FLUSH OPTIMIZER_COSTS", true},
	{"FLUSH STATUS", true},
	{"FLUSH USER_RESOURCES", true},
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
	{"FLUSH BINARY LOGS, ENGINE LOGS, ERROR LOGS, GENERAL LOGS, SLOW LOGS, RELAY LOGS, HOSTS, PRIVILEGES, LOGS, STATUS, USER_RESOURCES, OPTIMIZER_COSTS", true},
	{"FLUSH HOSTS,", false},
	{"FLUSH , HOSTS", false},
	{"FLUSH HOSTS HOSTS", false},
	{"FLUSH status, tables", false},

	// RELAY LOGS takes a channel, a plain string with no line feed in it.
	{"FLUSH RELAY LOGS FOR CHANNEL 'c'", true},
	{"FLUSH RELAY LOGS FOR CHANNEL \"c\"", true},
	{"FLUSH RELAY LOGS FOR CHANNEL ''", true},
	{"FLUSH RELAY LOGS FOR CHANNEL 'c', RELAY LOGS", true},
	{"FLUSH RELAY LOGS FOR CHANNEL c", false},
	{"FLUSH RELAY LOGS FOR CHANNEL _utf8mb4'c'", false},
	{"FLUSH RELAY LOGS FOR CHANNEL N'c'", false},
	{"FLUSH RELAY LOGS FOR CHANNEL 'a' 'b'", false},
	{"FLUSH RELAY LOGS FOR CHANNEL 'a\nb'", false}, // 1525, while parsing

	// The table forms.
	{"FLUSH TABLES", true},
	{"FLUSH TABLE", true},
	{"FLUSH TABLES t", true},
	{"FLUSH TABLES t, u", true},
	{"FLUSH TABLES test.t", true},
	{"FLUSH TABLES `t`", true},
	{"FLUSH TABLES WITH READ LOCK", true},
	{"FLUSH TABLES t WITH READ LOCK", true},
	{"FLUSH TABLES t FOR EXPORT", true},
	{"FLUSH TABLES t, u FOR EXPORT", true},
	{"FLUSH LOCAL TABLES", true},
	{"FLUSH LOCAL TABLES t WITH READ LOCK", true},
	{"FLUSH TABLES t,", false},
	{"FLUSH TABLES 't'", false},
	{"FLUSH TABLES t.*", false},
	{"FLUSH TABLES t AS x", false},
	{"FLUSH TABLES FOR EXPORT", false},
	{"FLUSH TABLES WITH READ", false},
	{"FLUSH TABLES WITH", false},
	{"FLUSH TABLES t FOR", false},
	{"FLUSH TABLES t HOSTS", false},
	{"FLUSH TABLES, HOSTS", false},
	{"FLUSH TABLES t FOR EXPORT, u", false},
	{"FLUSH TABLES t WITH READ LOCK FOR EXPORT", false},
	{"FLUSH TABLES @a", false},
	{"FLUSH TABLES _utf8", false},
}

// TestFlushForms holds omni to flushFormCases.
func TestFlushForms(t *testing.T) {
	for _, tc := range flushFormCases {
		_, err := Parse(tc.sql)
		if got := err == nil; got != tc.accept {
			t.Errorf("Parse(%q) accepts=%v, MySQL 8.0 accepts=%v: %v", tc.sql, got, tc.accept, err)
		}
	}
}
