package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
)

func TestUnicodeStrings(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "CJK characters in string value",
			input: `db.users.insertOne({ name: "张三" })`,
		},
		{
			name:  "emoji in string value",
			input: `db.posts.insertOne({ content: "Hello 🌍🎉" })`,
		},
		{
			name:  "Japanese in string value",
			input: `db.docs.find({ title: "東京タワー" })`,
		},
		{
			name:  "Korean in string value",
			input: `db.docs.find({ city: "서울" })`,
		},
		{
			name:  "unicode escape in string",
			input: `db.users.insertOne({ name: "\u0041\u0042\u0043" })`,
		},
		{
			name:  "mixed ASCII and CJK",
			input: `db.users.find({ name: "Alice-张三-Bob" })`,
		},
		{
			name:  "CJK in collection name via bracket access",
			input: `db["用户"].find()`,
		},
		{
			name:  "emoji in collection name via bracket access",
			input: `db["📊data"].find()`,
		},
		{
			name:  "unicode key in document",
			input: `db.data.insertOne({ "名前": "太郎", "年齢": 25 })`,
		},
		{
			name:  "unicode in array values",
			input: `db.tags.insertOne({ tags: ["数据库", "MongoDB", "测试"] })`,
		},
		{
			name:  "unicode in nested document",
			input: `db.people.insertOne({ address: { city: "北京", country: "中国" } })`,
		},
		{
			name:  "unicode identifier as collection name",
			input: `db.données.find()`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := mustParse(t, tc.input)
			loc := node.GetLoc()
			if loc.Start < 0 {
				t.Errorf("Start = %d, want >= 0", loc.Start)
			}
			if loc.End <= loc.Start {
				t.Errorf("End = %d <= Start = %d", loc.End, loc.Start)
			}
		})
	}
}

func TestUnicodeStringContent(t *testing.T) {
	// Verify that unicode string content is preserved correctly.
	node := mustParse(t, `db.users.insertOne({ name: "张三" })`)
	cs, ok := node.(*ast.CollectionStatement)
	if !ok {
		t.Fatalf("expected CollectionStatement, got %T", node)
	}
	if len(cs.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(cs.Args))
	}
	doc, ok := cs.Args[0].(*ast.Document)
	if !ok {
		t.Fatalf("expected Document arg, got %T", cs.Args[0])
	}
	if len(doc.Pairs) != 1 {
		t.Fatalf("expected 1 pair, got %d", len(doc.Pairs))
	}
	val, ok := doc.Pairs[0].Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expected StringLiteral, got %T", doc.Pairs[0].Value)
	}
	if val.Value != "张三" {
		t.Errorf("string value = %q, want %q", val.Value, "张三")
	}
}

func TestUnicodeEscapeContent(t *testing.T) {
	// \u0041 = A, \u0042 = B, \u0043 = C
	node := mustParse(t, `db.test.insertOne({ v: "\u0041\u0042\u0043" })`)
	cs := node.(*ast.CollectionStatement)
	doc := cs.Args[0].(*ast.Document)
	val := doc.Pairs[0].Value.(*ast.StringLiteral)
	if val.Value != "ABC" {
		t.Errorf("unicode escape value = %q, want %q", val.Value, "ABC")
	}
}
