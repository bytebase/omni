package testcontainer

import (
	"context"
	"errors"
	"testing"
	"time"

	gomysql "github.com/go-sql-driver/mysql"
)

// TestRetryOceanBase pins what the retrying connector promises a verdict helper:
// a transient answer is asked again, a real answer comes back as the server's
// error, and a transient answer that never clears is no answer at all.
func TestRetryOceanBase(t *testing.T) {
	window, pause := retryWindow, retryPause
	t.Cleanup(func() { retryWindow, retryPause = window, pause })
	retryWindow, retryPause = 50*time.Millisecond, time.Millisecond
	ctx := context.Background()
	shutdown := &gomysql.MySQLError{Number: 4550, Message: "rootservice is shutdown"}

	calls := 0
	if err := retry(ctx, func() error {
		calls++
		if calls < 3 {
			return shutdown
		}
		return nil
	}); err != nil || calls != 3 {
		t.Errorf("a transient answer that clears: err=%v after %d calls, want nil after 3", err, calls)
	}

	calls = 0
	err := retry(ctx, func() error {
		calls++
		return &gomysql.MySQLError{Number: 1064, Message: "syntax"}
	})
	var me *gomysql.MySQLError
	if !errors.As(err, &me) || me.Number != 1064 || calls != 1 {
		t.Errorf("a parse error: err=%v after %d calls, want 1064 after 1", err, calls)
	}

	err = retry(ctx, func() error { return shutdown })
	if err == nil || errors.As(err, &me) {
		t.Errorf("a transient answer that never clears came back as %#v, want a plain error", err)
	}
}
