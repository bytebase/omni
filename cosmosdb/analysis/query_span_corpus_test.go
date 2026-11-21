package analysis_test

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/cosmosdb/analysis"
)

func TestAnalyzeQuerySpanCorpus(t *testing.T) {
	type projectionWant struct {
		name  string
		paths []string
	}
	tests := []struct {
		name        string
		sql         string
		selectStar  bool
		projections []projectionWant
		predicates  []string
	}{
		{
			name: "basic field projection with alias and bracket string access",
			sql:  `select c.category, c.salary["value"] FROM container c`,
			projections: []projectionWant{
				{name: "category", paths: []string{"container.category"}},
				{name: "value", paths: []string{"container.salary.value"}},
			},
		},
		{
			name:       "SELECT TOP with star",
			sql:        `SELECT TOP 10 * FROM c`,
			selectStar: true,
		},
		{
			name: "SELECT TOP with field projection",
			sql:  `SELECT TOP 5 c.name, c.email FROM container c`,
			projections: []projectionWant{
				{name: "name", paths: []string{"container.name"}},
				{name: "email", paths: []string{"container.email"}},
			},
		},
		{
			name: "SELECT with ORDER BY",
			sql:  `SELECT c.name, c.population FROM container c ORDER BY c.population DESC`,
			projections: []projectionWant{
				{name: "name", paths: []string{"container.name"}},
				{name: "population", paths: []string{"container.population"}},
			},
		},
		{
			name: "SELECT with OFFSET LIMIT",
			sql:  `SELECT c.name FROM container c ORDER BY c.name OFFSET 0 LIMIT 10`,
			projections: []projectionWant{
				{name: "name", paths: []string{"container.name"}},
			},
		},
		{
			name: "SELECT with function in projection",
			sql:  `SELECT UPPER(c.name) AS upperName, c.email FROM container c`,
			projections: []projectionWant{
				{name: "upperName", paths: []string{"container.name"}},
				{name: "email", paths: []string{"container.email"}},
			},
		},
		{
			name: "SELECT with aggregation",
			sql:  `SELECT COUNT(1) AS totalRecords FROM c`,
			projections: []projectionWant{
				{name: "totalRecords"},
			},
		},
		{
			name: "SELECT with GROUP BY",
			sql:  `SELECT c.country, COUNT(1) AS cnt FROM container c GROUP BY c.country`,
			projections: []projectionWant{
				{name: "country", paths: []string{"container.country"}},
				{name: "cnt"},
			},
		},
		{
			name: "SELECT DISTINCT VALUE",
			sql:  `SELECT DISTINCT VALUE c.country FROM container c`,
			projections: []projectionWant{
				{name: "country", paths: []string{"container.country"}},
			},
		},
		{
			name: "SELECT VALUE with container alias",
			sql:  `SELECT VALUE c.name FROM container c`,
			projections: []projectionWant{
				{name: "name", paths: []string{"container.name"}},
			},
		},
		{
			name: "SELECT with underscore system fields",
			sql:  `SELECT c._ts, c._etag, c.name FROM container c`,
			projections: []projectionWant{
				{name: "_ts", paths: []string{"container._ts"}},
				{name: "_etag", paths: []string{"container._etag"}},
				{name: "name", paths: []string{"container.name"}},
			},
		},
		{
			name:       "WHERE with not-equal operator",
			sql:        `SELECT * FROM container c WHERE c.status != "inactive"`,
			selectStar: true,
			predicates: []string{"container.status"},
		},
		{
			name:       "WHERE with IN operator",
			sql:        `SELECT * FROM container c WHERE c.country IN ("US", "UK", "CA")`,
			selectStar: true,
			predicates: []string{"container.country"},
		},
		{
			name:       "WHERE with BETWEEN operator",
			sql:        `SELECT * FROM container c WHERE c.population BETWEEN 100000 AND 5000000`,
			selectStar: true,
			predicates: []string{"container.population"},
		},
		{
			name: "WHERE with compound condition and projections",
			sql:  `SELECT c.name, c.email FROM container c WHERE c.country != "US" AND c._ts > 1000`,
			projections: []projectionWant{
				{name: "name", paths: []string{"container.name"}},
				{name: "email", paths: []string{"container.email"}},
			},
			predicates: []string{"container._ts", "container.country"},
		},
		{
			name:       "WHERE with CONTAINS function call",
			sql:        `SELECT * FROM container c WHERE CONTAINS(c.name, "test")`,
			selectStar: true,
			predicates: []string{"container.name"},
		},
		{
			name:       "WHERE without alias",
			sql:        `SELECT * FROM container WHERE container.name = 'test'`,
			selectStar: true,
			predicates: []string{"container.name"},
		},
		{
			name:       "WHERE with array index path",
			sql:        `SELECT * FROM container WHERE container.addresses[1].country = 'Canada'`,
			selectStar: true,
			predicates: []string{"container.addresses[1].country"},
		},
		{
			name:       "WHERE with UDF call and array index path",
			sql:        `SELECT * FROM container c WHERE udf.foo(c.name, c.salary, c.addresses[1].country) = c.age`,
			selectStar: true,
			predicates: []string{"container.addresses[1].country", "container.age", "container.name", "container.salary"},
		},
		{
			name:       "WHERE with LIKE operator",
			sql:        `SELECT * FROM container c WHERE c.name LIKE "Alice%"`,
			selectStar: true,
			predicates: []string{"container.name"},
		},
		{
			name:       "WHERE with unary NOT",
			sql:        `SELECT * FROM container c WHERE NOT c.active`,
			selectStar: true,
			predicates: []string{"container.active"},
		},
		{
			name:       "WHERE with OR operator",
			sql:        `SELECT * FROM container c WHERE c.email = "x" OR c.name = "y"`,
			selectStar: true,
			predicates: []string{"container.email", "container.name"},
		},
		{
			name:       "WHERE with string concatenation",
			sql:        `SELECT * FROM container c WHERE c.firstName || " " || c.lastName = "Alice Smith"`,
			selectStar: true,
			predicates: []string{"container.firstName", "container.lastName"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			qa := analysis.Analyze(parseSelect(t, tc.sql))
			if qa.SelectStar != tc.selectStar {
				t.Errorf("SelectStar = %v, want %v", qa.SelectStar, tc.selectStar)
			}
			if len(qa.Projections) != len(tc.projections) {
				t.Fatalf("len(Projections) = %d, want %d", len(qa.Projections), len(tc.projections))
			}
			for i, want := range tc.projections {
				got := qa.Projections[i]
				if got.Name != want.name {
					t.Errorf("Projections[%d].Name = %q, want %q", i, got.Name, want.name)
				}
				var gotPaths []string
				for _, fp := range got.SourcePaths {
					gotPaths = append(gotPaths, fp.String())
				}
				if !slices.Equal(gotPaths, want.paths) {
					t.Errorf("Projections[%d].SourcePaths = %v, want %v", i, gotPaths, want.paths)
				}
			}

			var gotPreds []string
			for _, fp := range qa.Predicates {
				gotPreds = append(gotPreds, fp.String())
			}
			slices.Sort(gotPreds)
			wantPreds := slices.Clone(tc.predicates)
			slices.Sort(wantPreds)
			if !slices.Equal(gotPreds, wantPreds) {
				t.Errorf("Predicates = %v, want %v", gotPreds, wantPreds)
			}
		})
	}
}
