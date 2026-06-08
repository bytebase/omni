package parser

import "testing"

func TestRedshiftAlterExternalTableOptionsParse(t *testing.T) {
	tests := []string{
		"ALTER TABLE spectrum.sales SET TABLE PROPERTIES ('numRows'='1000000', 'compressionType'='gzip');",
		"ALTER TABLE spectrum.events SET LOCATION 's3://mybucket/data/events/2024/';",
		"ALTER TABLE spectrum.logs SET FILE FORMAT PARQUET;",
		"ALTER TABLE spectrum.sales ADD PARTITION (year=2024, month=1) LOCATION 's3://mybucket/sales/2024/01/';",
		"ALTER TABLE spectrum.events ADD IF NOT EXISTS PARTITION (date='2024-01-01') LOCATION 's3://mybucket/events/2024-01-01/';",
		"ALTER TABLE spectrum.sales DROP PARTITION (year=2023, month=12);",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}
		})
	}
}
