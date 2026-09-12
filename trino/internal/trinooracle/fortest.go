package trinooracle

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var forTest struct {
	once sync.Once
	o    *Oracle
	err  error
}

// ForTest returns the Oracle the differential tests run against: the server at
// $TRINO_ORACLE_URL or the conventional localhost:18080 when one answers, else
// a Trino testcontainer started once per test binary. Docker is assumed, so a
// server that can be neither reached nor started fails the test.
func ForTest(t interface {
	Helper()
	Fatalf(format string, args ...any)
}) *Oracle {
	t.Helper()
	forTest.once.Do(func() {
		o := Connect("")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := o.Ping(ctx)
		cancel()
		if err == nil {
			forTest.o = o
			return
		}
		if URLFromEnv() != "" {
			forTest.err = fmt.Errorf("trino at %s unreachable: %w", o.baseURL, err)
			return
		}
		o, _, err = StartContainer(context.Background())
		forTest.o, forTest.err = o, err
	})
	if forTest.err != nil {
		t.Fatalf("trino oracle: %v", forTest.err)
	}
	return forTest.o
}
