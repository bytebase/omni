package parser

import (
	"strings"
	"testing"
	"time"
)

// TestParseEventNotificationQueueName: a queue name may leave out any part
// but the last, as SQL Server 2022 takes, and has at most three parts.
func TestParseEventNotificationQueueName(t *testing.T) {
	for _, tt := range []struct {
		sql string
		ok  bool
	}{
		{"CREATE EVENT NOTIFICATION a ON QUEUE .q FOR QUEUE_ACTIVATION TO SERVICE 'a', 'a'", true},
		{"CREATE EVENT NOTIFICATION a ON QUEUE ..q FOR QUEUE_ACTIVATION TO SERVICE 'a', 'a'", true},
		{"CREATE EVENT NOTIFICATION a ON QUEUE .s.q FOR QUEUE_ACTIVATION TO SERVICE 'a', 'a'", true},
		{"CREATE EVENT NOTIFICATION a ON QUEUE d..q FOR QUEUE_ACTIVATION TO SERVICE 'a', 'a'", true},
		{"DROP EVENT NOTIFICATION a ON QUEUE ..q", true},
		{"CREATE EVENT NOTIFICATION a ON QUEUE d.s. FOR QUEUE_ACTIVATION TO SERVICE 'a', 'a'", false},
		{"CREATE EVENT NOTIFICATION a ON QUEUE ...q FOR QUEUE_ACTIVATION TO SERVICE 'a', 'a'", false},
		{"DROP EVENT NOTIFICATION a ON QUEUE a.b.c.d", false},
	} {
		_, err := Parse(tt.sql)
		if (err == nil) != tt.ok {
			t.Errorf("Parse(%q): err = %v, want ok = %v", tt.sql, err, tt.ok)
		}
	}
}

// TestParseEventSessionLongNotChain: a predicate's chain of NOT is read in
// linear time; building its text a NOT at a time took gigabytes.
func TestParseEventSessionLongNotChain(t *testing.T) {
	sql := "ALTER EVENT SESSION a ON DATABASE ADD EVENT a.a(WHERE " + strings.Repeat("NOT ", 200000) + "a = 1)"
	err, ok := parseWithin(sql, 10*time.Second)
	if !ok {
		t.Fatal("no result within 10s")
	}
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
}
