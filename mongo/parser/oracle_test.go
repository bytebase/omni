package parser

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/mongo/internal/testcontainer"
)

var mongoOracle testcontainer.Container

func TestMain(m *testing.M) {
	code := m.Run()
	_ = testcontainer.TerminateContainer(mongoOracle)
	os.Exit(code)
}

func TestMongoShellOracle(t *testing.T) {
	var err error
	mongoOracle, err = testcontainer.Start(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		script string
		accept bool
	}{
		{"db.items.insertOne({id:1, value:'one'})", true},
		{"db.items.find({})", true},
		{"db.items.find({id:1}).limit(1)", true},
		{"db.items.findOne({id:1})", true},
		{"db.items.updateOne({id:1}, {$set:{value:'two'}})", true},
		{"db.items.aggregate([{$match:{id:1}}])", true},
		{"db.items.deleteOne({id:1})", true},
		{"db.items.find({)", false},
		{"db.items.find({id:})", false},
		{"db.items.find(", false},
	} {
		t.Run(tc.script, func(t *testing.T) {
			_, parseErr := Parse(tc.script)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			code, output, err := mongoOracle.Exec(ctx, []string{"mongosh", "--quiet", "--norc", "omni", "--eval", tc.script})
			if err != nil {
				t.Fatal(err)
			}
			data, err := io.ReadAll(output)
			if err != nil {
				t.Fatal(err)
			}
			if code != 0 && !strings.Contains(string(data), "SyntaxError") {
				t.Fatalf("no syntax verdict: %s", data)
			}
			if (parseErr == nil) != tc.accept || (code == 0) != tc.accept {
				t.Errorf("parser error %v, mongosh exit %d: %s", parseErr, code, data)
			}
		})
	}
}
