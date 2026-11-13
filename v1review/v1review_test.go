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
