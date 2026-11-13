package docs

// What a converter's own tests read: one page or one form, converted.

// A Page is one page converted.
type Page struct {
	Forms  []string   // alternatives of the page's rule, as ANTLR text
	Frags  []Fragment // the fragments its forms reach
	Failed []string   // forms and definitions not converted, with why
	Fixes  []string   // the reasons of the fixes applied to it
}

// A Fragment is a rule a page defines with ::= or where ... is.
type Fragment struct {
	Name string
	Alts []string
}

// ConvertPage converts one page's blocks as Run does.
func ConvertPage(d *Dialect, slug string, blocks []string, fixes []*Fix) Page {
	pg := convertPage(d, slug, d.Site+"/"+slug, blocks, fixes)
	out := Page{Forms: pg.forms, Fixes: pg.fixes}
	for _, f := range pg.frags {
		out.Frags = append(out.Frags, Fragment{Name: f.name, Alts: f.alts})
	}
	for _, f := range pg.failed {
		out.Failed = append(out.Failed, f.what+": "+f.reason)
	}
	return out
}

// Form converts notation text to ANTLR alternatives, on a page named p
// that defines no fragment.
func Form(d *Dialect, notation string) ([]string, error) {
	return (&converter{d: d, slug: "p", page: "p", defs: map[string]bool{}}).form(notation)
}

// Statement converts a statement form, which scopes a | outside brackets
// to what follows the object's name.
func Statement(d *Dialect, text string) ([]string, error) {
	return (&converter{d: d, slug: "p", page: "p", defs: map[string]bool{}}).statement(text)
}

// The errors Form and Statement return for a form they do not convert.
var (
	ErrScope  = errScope
	ErrElided = errElided
)
