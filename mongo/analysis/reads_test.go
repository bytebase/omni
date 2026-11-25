package analysis_test

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/mongo/analysis"
)

func TestReads(t *testing.T) {
	ns := func(db, coll string) analysis.Namespace { return analysis.Namespace{Database: db, Collection: coll} }
	tests := []struct {
		input string
		want  []analysis.Namespace
	}{
		{`db.users.find()`, []analysis.Namespace{ns("", "users")}},
		{`db.users.find().explain()`, []analysis.Namespace{ns("", "users")}},
		{`db.users.insertOne({a: 1})`, nil},
		{`db.users.aggregate([{$out: "report"}])`, []analysis.Namespace{ns("", "users")}},
		{`db.users.aggregate([{$lookup: {from: "orders", pipeline: [{$lookup: {from: "items", localField: "a", foreignField: "b", as: "i"}}], as: "o"}}])`,
			[]analysis.Namespace{ns("", "users"), ns("", "orders"), ns("", "items")}},
		{`db.users.aggregate([{$lookup: {from: {db: "hr", coll: "emp"}, pipeline: [{$unionWith: "dept"}], as: "e"}}])`,
			[]analysis.Namespace{ns("", "users"), ns("hr", "emp"), ns("hr", "dept")}},
		{`db.users.aggregate([{$unionWith: {coll: "a", pipeline: [{$facet: {x: [{$graphLookup: {from: "b", startWith: "$x", connectFromField: "x", connectToField: "y", as: "g"}}]}}]}}])`,
			[]analysis.Namespace{ns("", "users"), ns("", "a"), ns("", "b")}},
		{`db.getSiblingDB("hr").emp.find().limit(1)`, []analysis.Namespace{ns("hr", "emp")}},
		{`db.getSiblingDB("hr").getCollection("emp").aggregate([{$unionWith: "dept"}])`, []analysis.Namespace{ns("hr", "emp"), ns("hr", "dept")}},
		{`db.getSiblingDB("hr").getMongo().find()`, nil},
		{`db.getSiblingDB("hr").emp.insertOne({a: 1})`, nil},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := mustAnalyze(t, tt.input).Reads
			if !slices.Equal(got, tt.want) {
				t.Errorf("Reads = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExpandViews(t *testing.T) {
	views := map[analysis.Namespace]string{
		{Collection: "v"}:                 `{"viewOn": "w", "pipeline": [{"$lookup": {"from": "orders", "localField": "a", "foreignField": "b", "as": "o"}}]}`,
		{Collection: "w"}:                 `{"viewOn": "users", "pipeline": [{"$unionWith": {"coll": "v"}}]}`,
		{Database: "hr", Collection: "s"}: `{"viewOn": "emp", "pipeline": []}`,
	}
	definition := func(ns analysis.Namespace) (string, bool) {
		d, ok := views[ns]
		return d, ok
	}
	got, err := analysis.ExpandViews([]analysis.Namespace{{Collection: "v"}, {Database: "hr", Collection: "s"}}, definition)
	if err != nil {
		t.Fatal(err)
	}
	want := []analysis.Namespace{{Collection: "v"}, {Database: "hr", Collection: "s"}, {Collection: "w"}, {Collection: "orders"}, {Database: "hr", Collection: "emp"}, {Collection: "users"}}
	if !slices.Equal(got, want) {
		t.Errorf("ExpandViews = %v, want %v", got, want)
	}

	for _, bad := range []string{`not json`, `{"pipeline": []}`} {
		_, err := analysis.ExpandViews([]analysis.Namespace{{Collection: "x"}}, func(analysis.Namespace) (string, bool) { return bad, true })
		if err == nil {
			t.Errorf("ExpandViews(%s): no error", bad)
		}
	}
}
