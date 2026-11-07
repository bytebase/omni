package parser

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/cassandra/internal/testcontainer"
)

var cassandraOracle testcontainer.Container

func TestMain(m *testing.M) {
	code := m.Run()
	_ = testcontainer.TerminateContainer(cassandraOracle)
	os.Exit(code)
}

func executeCQL(t *testing.T, sql string) (int, string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	code, output, err := cassandraOracle.Exec(ctx, []string{"cqlsh", "-e", sql})
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(output)
	if err != nil {
		t.Fatal(err)
	}
	return code, string(data)
}

func TestCassandraQueryOracle(t *testing.T) {
	var err error
	cassandraOracle, err = testcontainer.Start(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	code, output := executeCQL(t, "CREATE KEYSPACE IF NOT EXISTS omni WITH replication = {'class':'SimpleStrategy','replication_factor':1}; CREATE TABLE IF NOT EXISTS omni.items (id int PRIMARY KEY, value text); TRUNCATE omni.items;")
	if code != 0 {
		t.Fatalf("fixture: %s", output)
	}
	for _, tc := range []struct {
		sql    string
		accept bool
	}{
		{"SELECT * FROM omni.items;", true},
		{"INSERT INTO omni.items (id, value) VALUES (1, 'one');", true},
		{"SELECT value FROM omni.items WHERE id = 1;", true},
		{"UPDATE omni.items SET value = 'two' WHERE id = 1;", true},
		{"DELETE FROM omni.items WHERE id = 1;", true},
		{"SELECT FROM;", false},
		{"SELECT * FROM omni.items WHERE;", false},
		{"INSERT INTO omni.items (id) VALUES ();", false},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			_, parseErr := Parse(tc.sql)
			code, output := executeCQL(t, tc.sql)
			if code != 0 && !strings.Contains(output, "SyntaxException") {
				t.Fatalf("no syntax verdict: %s", output)
			}
			if (parseErr == nil) != tc.accept || (code == 0) != tc.accept {
				t.Errorf("parser error %v, cqlsh exit %d: %s", parseErr, code, output)
			}
		})
	}
}
