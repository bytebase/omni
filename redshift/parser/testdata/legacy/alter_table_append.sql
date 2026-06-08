-- ALTER TABLE APPEND command tests
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_TABLE_APPEND.html

-- Basic ALTER TABLE APPEND
ALTER TABLE sales APPEND FROM sales_monthly;
ALTER TABLE events APPEND FROM staging_events;
ALTER TABLE customers APPEND FROM temp_customers;

-- With schema qualification
ALTER TABLE public.sales APPEND FROM staging.sales_monthly;
ALTER TABLE analytics.fact_sales APPEND FROM staging.sales_data;
ALTER TABLE warehouse.products APPEND FROM temp.product_import;

-- With IGNOREEXTRA option
-- Discards extra columns in source table that are not in target table
ALTER TABLE sales APPEND FROM sales_listing IGNOREEXTRA;
ALTER TABLE events APPEND FROM staging_events_extended IGNOREEXTRA;
ALTER TABLE customers APPEND FROM temp_customers_full IGNOREEXTRA;

-- With schema qualification and IGNOREEXTRA
ALTER TABLE public.sales APPEND FROM staging.sales_extended IGNOREEXTRA;
ALTER TABLE analytics.fact_sales APPEND FROM staging.sales_data_full IGNOREEXTRA;

-- With FILLTARGET option
-- Fills extra columns in target table with default values or NULL
ALTER TABLE sales_report APPEND FROM sales_month FILLTARGET;
ALTER TABLE event_summary APPEND FROM event_data FILLTARGET;
ALTER TABLE customer_analytics APPEND FROM customer_base FILLTARGET;

-- With schema qualification and FILLTARGET
ALTER TABLE public.sales_report APPEND FROM staging.sales_month FILLTARGET;
ALTER TABLE analytics.fact_sales_extended APPEND FROM staging.sales_basic FILLTARGET;

-- From materialized views
ALTER TABLE target_tbl APPEND FROM my_streaming_materialized_view;
ALTER TABLE sales_fact APPEND FROM sales_stream_mv;
ALTER TABLE events_archive APPEND FROM events_stream_mv;

-- With schema qualification for materialized views
ALTER TABLE public.target_table APPEND FROM streaming.my_materialized_view;
ALTER TABLE analytics.fact_table APPEND FROM staging.stream_mv;

-- Materialized views with options
ALTER TABLE sales_archive APPEND FROM sales_stream_mv IGNOREEXTRA;
ALTER TABLE events_extended APPEND FROM events_stream_mv FILLTARGET;

-- Complex examples with various naming conventions
ALTER TABLE "Sales_Table" APPEND FROM "Monthly_Sales";
ALTER TABLE table_with_underscores APPEND FROM source_table_123;
ALTER TABLE "MixedCase_Table" APPEND FROM "source_MixedCase" IGNOREEXTRA;
ALTER TABLE target123 APPEND FROM source456 FILLTARGET;

-- Long identifiers
ALTER TABLE very_long_table_name_with_many_words APPEND FROM another_very_long_source_table_name;
ALTER TABLE schema_with_long_name.table_with_long_name APPEND FROM source_schema.source_table_with_long_name IGNOREEXTRA;

-- Edge cases and various patterns
ALTER TABLE t1 APPEND FROM t2;
ALTER TABLE t APPEND FROM s IGNOREEXTRA;
ALTER TABLE target APPEND FROM source FILLTARGET;
ALTER TABLE "TARGET" APPEND FROM "SOURCE";
ALTER TABLE schema1.table1 APPEND FROM schema2.table2;
ALTER TABLE s1.t1 APPEND FROM s2.t2 IGNOREEXTRA;
ALTER TABLE s.t APPEND FROM s.t2 FILLTARGET;