package main

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/grammar/docs"
)

// TestNotationReadsAWSBlocks pins how a Syntax block of AWS's, as its HTML
// writes it, is read into the notation docs converts.
func TestNotationReadsAWSBlocks(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		// Italics are values, words in lower case are too, keywords written
		// in lower case are not.
		{"CREATE GROUP <em>group_name</em> [ WITH USER username ]", "CREATE GROUP <group_name> [ WITH USER <username> ]"},
		{"IAM_ROLE { default | 'SESSION' }", "IAM_ROLE { DEFAULT | 'SESSION' }"},
		{"<em>column_name data_type</em>", "<column_name> <data_type>"},
		{"<em>schema_name</em>.<em>table</em>", "<schema_name>.<table>"},
		// A quoted string in lower case is a value; one in upper case a literal.
		{"URI '<em>hostname</em>' REGION 'aws-region'", "URI '<hostname>' REGION '<aws-region>'"},
		{"'arn:aws:iam::&lt;account-id&gt;:role/&lt;role-name&gt;'", "'arn:aws:iam::<account-id>:role/<role-name>'"},
		{"TO <em>&lt;prefix&gt;:&lt;username&gt;</em>", "TO <prefix>:<username>"},
		// A definition's head, and a blank line between its variants.
		{"ALTER USER <em>name</em> <em>option</em>\n\nwhere <em>option</em> is\n\nCREATEDB | NOCREATEDB",
			"ALTER USER <name> <option>\n\noption ::=\n\nCREATEDB | NOCREATEDB"},
		{"X <em>c</em>\n\nwhere <em>c</em> is one of:\n\nGRANT A\n  ON B\n\nREVOKE A\n  ON B",
			"X <c>\n\nc ::=\n\nGRANT A\n  ON B\n\n--\nREVOKE A\n  ON B"},
		{"x ‘a’ … [, …]", "<x> '<a>' ... [, ...]"},
	} {
		if got := notation(c.in); got != c.want {
			t.Errorf("%q\n got %q\nwant %q", c.in, got, c.want)
		}
	}
}

// TestSyntaxBlocks pins which code blocks are a page's syntax: those under a
// heading that says Syntax, to the next heading of its level or above.
func TestSyntaxBlocks(t *testing.T) {
	page := `<h2 id="x-privileges">Required privileges</h2><pre class="programlisting">NOT THIS</pre>
<h2 id="x-synopsis">Syntax
</h2><pre class="programlisting"><code>DROP X <em>name</em></code></pre>
<pre class="screen">DROP Y</pre>
<h2 id="x-parameters">Parameters</h2><pre class="programlisting">NOR THIS</pre>
<h2 id="x-full">Full X</h2><h3 id="x-full-syntax">Full X syntax</h3><pre class="programlisting">CREATE X</pre><h3>Notes</h3><pre class="programlisting">NOR THAT</pre>`
	if got := strings.Join(syntaxBlocks(page), " | "); got != "DROP X <name> | DROP Y | CREATE X" {
		t.Errorf("got %q", got)
	}
}

func TestCommandPages(t *testing.T) {
	index := `<a href="./c_other.html">x</a><div class="highlights"><ul><li><a href="./r_ABORT.html">ABORT</a></li><li><a href="./close.html">CLOSE</a></li><li><a href="./r_ABORT.html">ABORT</a></li></ul></div><a href="./r_in_condition.html">x</a>`
	if got := strings.Join(commandPages(index), " "); got != "/redshift/latest/dg/r_ABORT.html /redshift/latest/dg/close.html" {
		t.Errorf("got %q", got)
	}
	if got := ruleName("r_CREATE_TABLE_NEW"); got != "createTableNew" {
		t.Errorf("rule name %q", got)
	}
}

func TestPlaceholderClasses(t *testing.T) {
	for name, want := range map[string]string{
		"table_name": docs.ClassName, "username": docs.ClassName, "encode_type": docs.ClassName,
		"query": docs.ClassQuery, "select_statement": docs.ClassQuery, "statement": docs.ClassQuery,
		"data_type": docs.ClassDataType, "argtype": docs.ClassDataType, "column_type": docs.ClassDataType,
		"expression": docs.ClassExpr, "using_predicate_exp": docs.ClassExpr, "condition": docs.ClassExpr,
		"seed": docs.ClassNumber, "row_limit": docs.ClassNumber, "max-size": docs.ClassNumber, "process_id": docs.ClassNumber, "start": docs.ClassNumber,
		"with_subquery": classCTE, "common_table_expression": classCTE, "argmode": classArgMode, "schema_element": classElement,
	} {
		if got := placeholderClass(name); got != want {
			t.Errorf("<%s> is %s, want %s", name, got, want)
		}
	}
}
