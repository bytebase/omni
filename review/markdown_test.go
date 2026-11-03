package review

import (
	"strings"
	"testing"
)

func TestCode(t *testing.T) {
	tests := []struct{ in, want string }{
		{"orders", "`orders`"},
		{"_tmp_", "`_tmp_`"},
		{"retail.shop", "`retail.shop`"},
		{"`orders`", "`` `orders` ``"},
		{"a``b", "```a``b```"},
		{"a\nb\tc\r", "`a\\nb\\tc\\r`"},
		{" a ", "`  a  `"},
		{"  ", "`  `"},
		{"", "`  `"},
		{strings.Repeat("a", 130), "`" + strings.Repeat("a", 119) + "…`"},
	}
	for _, tt := range tests {
		if got := Code(tt.in); got != tt.want {
			t.Errorf("Code(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestCodeBlock(t *testing.T) {
	tests := []struct{ in, want string }{
		{"SELECT 1;\n       ^", "```\nSELECT 1;\n       ^\n```"},
		{"x\n", "```\nx\n```"},
		{"a ``` b", "````\na ``` b\n````"},
	}
	for _, tt := range tests {
		if got := CodeBlock(tt.in); got != tt.want {
			t.Errorf("CodeBlock(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestEscape(t *testing.T) {
	if got, want := Escape("a*b_c [d] <e> `f` ~g~ h|i \\"), `a\*b\_c \[d\] \<e\> \`+"`"+`f\`+"`"+` \~g\~ h\|i \\`; got != want {
		t.Errorf("Escape = %q, want %q", got, want)
	}
	if got := Escape("zero-length name (1.5)."); got != "zero-length name (1.5)." {
		t.Errorf("Escape changed prose with no inline Markdown: %q", got)
	}
}

func TestNamesAndJoin(t *testing.T) {
	tests := []struct {
		in   []string
		want string
	}{
		{nil, ""},
		{[]string{"t"}, "`t`"},
		{[]string{"t", "p"}, "`t` and `p`"},
		{[]string{"a", "b", "c"}, "`a`, `b`, and `c`"},
		{[]string{"f(integer, text)", "_g_"}, "`f(integer, text)` and `_g_`"},
	}
	for _, tt := range tests {
		if got := Names(tt.in...); got != tt.want {
			t.Errorf("Names(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
	if got, want := Join("`ix` of `dbo.t`", "`ix` of `dbo.p`"), "`ix` of `dbo.t` and `ix` of `dbo.p`"; got != want {
		t.Errorf("Join = %q, want %q", got, want)
	}
}

func TestPlural(t *testing.T) {
	tests := []struct {
		kind string
		n    int
		want string
	}{
		{"table", 1, "table"},
		{"table", 2, "tables"},
		{"index", 2, "indexes"},
		{"policy", 3, "policies"},
		{"operator class", 2, "operator classes"},
		{"operator family", 2, "operator families"},
		{"text search dictionary", 2, "text search dictionaries"},
		{"statistics", 2, "statistics"},
		{"materialized view", 2, "materialized views"},
		{"user mapping", 2, "user mappings"},
		{"key", 2, "keys"},
		{"temporary table", 0, "temporary tables"},
	}
	for _, tt := range tests {
		if got := Plural(tt.kind, tt.n); got != tt.want {
			t.Errorf("Plural(%q, %d) = %q, want %q", tt.kind, tt.n, got, tt.want)
		}
	}
}
