package parser

import "testing"

// TestLoadDataSourceForms: one FROM, before LOCAL and the file source, or
// Aurora's FROM S3; LOAD takes DATA or XML and no other word.
func TestLoadDataSourceForms(t *testing.T) {
	for _, tc := range []struct {
		sql    string
		accept bool
	}{
		{"LOAD DATA FROM LOCAL INFILE 'f' INTO TABLE t", true},
		{"LOAD DATA FROM S3 's3://b/x' INTO TABLE t", true},
		{"LOAD DATA FROM FROM S3 's3://b/x' INTO TABLE t", false},
		{"LOAD XML FROM FROM S3 's3://b/x' INTO TABLE t", false},
		{"LOAD DATA FROM FROM INFILE 'f' INTO TABLE t", false},
		{"LOAD DATABASE INFILE 'f' INTO TABLE t", false},
	} {
		if _, err := Parse(tc.sql); (err == nil) != tc.accept {
			t.Errorf("Parse(%q) error = %v, want accept %v", tc.sql, err, tc.accept)
		}
	}
}
