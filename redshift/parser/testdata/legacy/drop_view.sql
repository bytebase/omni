-- DROP VIEW command test cases
-- Based on AWS Redshift documentation: https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_VIEW.html

-- Basic DROP VIEW
DROP VIEW event;

-- DROP VIEW with schema qualification
DROP VIEW public.eventview;
DROP VIEW myschema.sales_view;

-- DROP VIEW with IF EXISTS clause
DROP VIEW IF EXISTS eventview;
DROP VIEW IF EXISTS public.eventview;
DROP VIEW IF EXISTS nonexistent_view;

-- DROP VIEW with CASCADE option
DROP VIEW eventview CASCADE;
DROP VIEW public.myeventview CASCADE;
DROP VIEW IF EXISTS eventview CASCADE;

-- DROP VIEW with RESTRICT option (default behavior)
DROP VIEW eventview RESTRICT;
DROP VIEW public.eventview RESTRICT;
DROP VIEW IF EXISTS eventview RESTRICT;

-- DROP multiple views in single statement
DROP VIEW view1, view2, view3;
DROP VIEW public.view1, schema2.view2, view3;

-- DROP multiple views with IF EXISTS
DROP VIEW IF EXISTS view1, view2, view3;
DROP VIEW IF EXISTS public.view1, schema2.view2, view3;

-- DROP multiple views with CASCADE
DROP VIEW view1, view2, view3 CASCADE;
DROP VIEW public.view1, schema2.view2, view3 CASCADE;

-- DROP multiple views with RESTRICT
DROP VIEW view1, view2, view3 RESTRICT;
DROP VIEW public.view1, schema2.view2, view3 RESTRICT;

-- Complex combinations
DROP VIEW IF EXISTS public.eventview, myschema.sales_view, customer_view CASCADE;
DROP VIEW IF EXISTS view1, public.view2, schema3.view3 RESTRICT;

-- Views with reserved keywords as names (quoted identifiers)
DROP VIEW "order";
DROP VIEW "user";
DROP VIEW "table";
DROP VIEW IF EXISTS "select" CASCADE;

-- Views with special characters in names
DROP VIEW view_with_underscores;
DROP VIEW view123;
DROP VIEW my_schema.view_name_2023;

-- Case sensitivity tests
DROP VIEW MyView;
DROP VIEW EVENTVIEW;
DROP VIEW eventview;

-- Additional AWS documentation examples
-- Example from documentation: drop view eventview cascade;
DROP VIEW eventview CASCADE;

-- Example scenarios from documentation
-- After creating: create view eventview as select dateid, eventname, catid from event where catid = 1;
-- After creating dependent: create view myeventview as select eventname, catid from eventview where eventname <> ' ';
-- This would fail: DROP VIEW eventview;
-- This succeeds: DROP VIEW eventview CASCADE;

-- Real-world examples
DROP VIEW sales_summary;
DROP VIEW quarterly_report;
DROP VIEW customer_analytics CASCADE;
DROP VIEW IF EXISTS temp_analysis_view;
DROP VIEW regional_sales, monthly_totals, yearly_summary RESTRICT;

-- Edge cases with whitespace and formatting
DROP VIEW
    eventview;

DROP    VIEW    IF    EXISTS    eventview    CASCADE;

DROP VIEW view1,
          view2,
          view3;

-- Mixed case keywords
Drop View eventview;
drop view if exists eventview cascade;
DROP view eventview RESTRICT;

-- Schema names with various formats
DROP VIEW prod.sales_view;
DROP VIEW staging.temp_view;
DROP VIEW analytics.customer_segment_view;