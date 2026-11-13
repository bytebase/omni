package v1review

import "testing"

func TestNamingConfigCompile(t *testing.T) {
	re, err := NamingConfig{Pattern: "^[a-z]+$"}.Compile()
	if err != nil || re == nil {
		t.Fatalf("Compile = %v, %v", re, err)
	}
	if !re.MatchString("ok") || re.MatchString("NO") {
		t.Errorf("pattern %q compiled wrong", re)
	}
	if re, err := (NamingConfig{MaxLength: 10}).Compile(); re != nil || err != nil {
		t.Errorf("empty Pattern: got %v, %v", re, err)
	}
	if re, err := (NamingConfig{Pattern: "("}).Compile(); re != nil || err == nil {
		t.Errorf("invalid pattern: got %v, %v", re, err)
	}
}

func TestTemplateConfigExpand(t *testing.T) {
	cfg := TemplateConfig{Format: "^fk_{{referencing_table}}_{{referencing_column}}_{{referenced_table}}$"}
	re, err := cfg.Expand(map[string]string{
		ReferencingTableToken:  "book",
		ReferencingColumnToken: "author_id",
		ReferencedTableToken:   "author",
	})
	if err != nil {
		t.Fatalf("Expand: %v", err)
	}
	if got, want := re.String(), "^fk_book_author_id_author$"; got != want {
		t.Errorf("pattern = %q, want %q", got, want)
	}

	// A token with no value stays as written; a repeated token is replaced everywhere.
	re, err = (TemplateConfig{Format: "{{table}}_{{table}}_{{column_list}}"}).Expand(map[string]string{TableToken: "t"})
	if err != nil {
		t.Fatalf("Expand: %v", err)
	}
	if got, want := re.String(), "t_t_{{column_list}}"; got != want {
		t.Errorf("pattern = %q, want %q", got, want)
	}

	// A value is not quoted: a metacharacter in it changes the pattern.
	re, err = (TemplateConfig{Format: "^idx_{{table}}$"}).Expand(map[string]string{TableToken: "a.b"})
	if err != nil || !re.MatchString("idx_axb") {
		t.Errorf("unquoted value: %v, %v", re, err)
	}
	if _, err := (TemplateConfig{Format: "^idx_{{table}}$"}).Expand(map[string]string{TableToken: "a(b"}); err == nil {
		t.Error("a value that breaks the pattern: no error")
	}

	// An empty format matches every name.
	re, err = (TemplateConfig{}).Expand(nil)
	if err != nil || !re.MatchString("anything") {
		t.Errorf("empty format: %v, %v", re, err)
	}
}
