-- SHOW VIEW command tests
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_VIEW.html

-- Basic usage: show view without schema
SHOW VIEW LA_Venues_v;

-- Show view with schema name
SHOW VIEW public.Sports_v;

-- Show view with quoted identifiers
SHOW VIEW "my_view";
SHOW VIEW "my_schema"."my_view";

-- Show views with mixed case names
SHOW VIEW MyView;
SHOW VIEW MySchema.MyView;

-- Show view with underscores in names
SHOW VIEW sales_summary_view;
SHOW VIEW reporting.monthly_sales_view;

-- Show view with numbers in names
SHOW VIEW view1;
SHOW VIEW schema1.view2;

-- Show view for materialized views (mentioned in documentation)
SHOW VIEW materialized_sales_view;
SHOW VIEW analytics.materialized_customer_view;

-- Show view for late-binding views (mentioned in documentation)
SHOW VIEW late_binding_external_view;
SHOW VIEW external_schema.late_binding_view;

-- Edge cases with special naming
SHOW VIEW "_view";
SHOW VIEW "123view";
SHOW VIEW "view-with-dash";
SHOW VIEW "schema-name"."view-name";

-- Real-world example based on documentation
-- First create a view (commented out as this is just testing SHOW VIEW)
-- CREATE VIEW LA_Venues_v AS SELECT * FROM venue WHERE venuecity='Los Angeles';
SHOW VIEW LA_Venues_v;

-- Another real-world example from documentation
-- CREATE VIEW public.Sports_v AS SELECT * FROM category WHERE catgroup='Sports';
SHOW VIEW public.Sports_v;

-- Common view patterns
SHOW VIEW customer_summary;
SHOW VIEW sales.revenue_by_region;
SHOW VIEW analytics.user_activity_daily;
SHOW VIEW reporting.quarterly_metrics;

-- Views with longer names
SHOW VIEW very_long_view_name_that_describes_complex_business_logic;
SHOW VIEW long_schema_name.another_long_view_name_for_testing;

-- Common abbreviations in view names
SHOW VIEW vw_customers;
SHOW VIEW v_sales;
SHOW VIEW dw.dim_product_v;
SHOW VIEW etl.stg_orders_vw;