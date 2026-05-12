package catalog

import "testing"

func TestPGAnalyzerScopeRegressions(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "lateral_setop_can_see_prior_lateral_alias",
			sql: `
				CREATE TABLE clips (body text);
				CREATE VIEW v_lateral_setop AS
					SELECT exp.token
					FROM clips c
					CROSS JOIN LATERAL (
						SELECT lower(c.body) AS norm
					) sp
					LEFT JOIN LATERAL (
						SELECT x.token
						FROM (
							SELECT regexp_split_to_table(COALESCE(sp.norm, ''), ',') AS token
							UNION ALL
							SELECT regexp_split_to_table(COALESCE(sp.norm, ''), ';') AS token
						) x
					) exp ON true;
			`,
		},
		{
			name: "derived_table_output_alias_visible_without_inner_join_leak",
			sql: `
				CREATE TABLE assortment_source (assortment_sku integer);
				CREATE TABLE group_source (a_group integer);
				CREATE VIEW v_derived_alias AS
					SELECT assortment_sku, a_group, row_num
					FROM (
						SELECT
							a.assortment_sku,
							group_join.a_group AS a_group,
							row_number() OVER () AS row_num
						FROM assortment_source a
						JOIN group_source group_join ON true
					) abc;
			`,
		},
		{
			name: "derived_table_does_not_expose_group_by_resjunk_alias",
			sql: `
				CREATE TABLE item_master_scope (assortment_sku integer, item_number numeric);
				CREATE TABLE group_source_scope (representative_item_number bigint, a_group text);
				CREATE VIEW v_derived_alias_group_by AS
					SELECT assortment_sku, a_group, row_num
					FROM (
						SELECT
							im.assortment_sku,
							group_join.a_group::character varying AS a_group,
							row_number() OVER (
								PARTITION BY im.assortment_sku
								ORDER BY group_join.a_group::character varying
							) AS row_num
						FROM item_master_scope im
						JOIN (
							SELECT representative_item_number, a_group
							FROM group_source_scope
						) group_join ON im.item_number = group_join.representative_item_number::numeric
						GROUP BY im.assortment_sku, group_join.a_group::character varying
					) abc
					WHERE row_num = 1;
			`,
		},
		{
			name: "unqualified_column_from_join_not_duplicated_by_join_rte",
			sql: `
				CREATE TABLE left_source (id integer);
				CREATE TABLE right_source (a_group integer);
				CREATE VIEW v_join_unqualified AS
					SELECT a_group
					FROM left_source l
					JOIN right_source r ON true;
			`,
		},
		{
			name: "top_level_with_visible_to_setop_branches",
			sql: `
				CREATE TABLE cte_source (id integer);
				CREATE VIEW v_cte_setop AS
					WITH core_hi_items AS (
						SELECT id FROM cte_source
					)
					SELECT id FROM core_hi_items
					UNION ALL
					SELECT id FROM core_hi_items;
			`,
		},
		{
			name: "setop_column_count_ignores_branch_resjunk_targets",
			sql: `
				CREATE TABLE union_source (id integer, sort_key integer);
				CREATE VIEW v_union_branch_order AS
					(SELECT id FROM union_source ORDER BY sort_key)
					UNION ALL
					SELECT id FROM union_source;
			`,
		},
		{
			name: "sublink_column_count_ignores_order_by_resjunk_targets",
			sql: `
				CREATE TABLE sublink_source (id integer, sort_key integer, new_value text);
				CREATE VIEW v_scalar_sublink_order AS
					SELECT (
						SELECT new_value
						FROM sublink_source inner_source
						ORDER BY inner_source.sort_key DESC
						LIMIT 1
					) AS latest_value
					FROM sublink_source outer_source;
				CREATE VIEW v_any_sublink_order AS
					SELECT id
					FROM sublink_source outer_source
					WHERE id IN (
						SELECT id
						FROM sublink_source inner_source
						GROUP BY id
						ORDER BY count(*) DESC
					);
			`,
		},
		{
			name: "materialized_view_with_cte_visible_to_setop_branches",
			sql: `
				CREATE TABLE mat_cte_source (id integer);
				CREATE MATERIALIZED VIEW mv_cte_setop AS
					WITH lagged_data AS (
						SELECT id FROM mat_cte_source
					)
					SELECT id FROM lagged_data
					UNION ALL
					SELECT id FROM lagged_data;
			`,
		},
		{
			name: "schema_record_function_with_column_definition_list",
			sql: `
				CREATE SCHEMA redshift;
				CREATE FUNCTION redshift.week_item_data()
				RETURNS SETOF record
				LANGUAGE sql
				AS $$
					SELECT 1::integer, 'sku-1'::text;
				$$;
				CREATE VIEW v_week_item_data AS
					SELECT wid.week_id, wid.sku
					FROM redshift.week_item_data() wid(week_id integer, sku text);
			`,
		},
		{
			name: "schema_returns_table_function_available_to_earlier_view",
			sql: `
				CREATE SCHEMA redshift;
				CREATE VIEW v_week_item_data_before_function AS
					SELECT wid.ad_end_date, wid.lob
					FROM redshift.week_item_data() wid(ad_end_date, lob);
				CREATE OR REPLACE FUNCTION redshift.week_item_data()
				RETURNS TABLE(ad_end_date date, lob text)
				LANGUAGE plpgsql
				AS $$ BEGIN RETURN; END; $$;
			`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := LoadSQL(tt.sql); err != nil {
				t.Fatalf("LoadSQL failed: %v", err)
			}
		})
	}
}

func TestPGAnalyzerScopeRejectsNonLateralSetopSiblingReference(t *testing.T) {
	_, err := LoadSQL(`
		CREATE TABLE non_lateral_base (id integer);
		CREATE VIEW v_non_lateral_setop AS
			SELECT s.id
			FROM non_lateral_base b,
			     (
				     SELECT b.id AS id
				     UNION ALL
				     SELECT b.id AS id
			     ) s;
	`)
	if err == nil {
		t.Fatal("LoadSQL unexpectedly accepted non-LATERAL set-op reference to sibling alias")
	}
}
