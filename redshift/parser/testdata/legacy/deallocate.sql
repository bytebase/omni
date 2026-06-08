-- DEALLOCATE command test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_DEALLOCATE.html

-- Basic DEALLOCATE statements
DEALLOCATE statement1;
deallocate statement1;

-- DEALLOCATE with PREPARE keyword (optional)
DEALLOCATE PREPARE statement2;
deallocate prepare statement2;

-- DEALLOCATE ALL to deallocate all prepared statements
DEALLOCATE ALL;
deallocate all;

-- DEALLOCATE with simple statement names
DEALLOCATE my_statement;
DEALLOCATE stmt_name;
DEALLOCATE test_statement;

-- DEALLOCATE with underscore in statement names
DEALLOCATE my_statement_name;
DEALLOCATE statement_with_long_name;
DEALLOCATE test_stmt_123;

-- DEALLOCATE with numeric suffixes
DEALLOCATE stmt1;
DEALLOCATE stmt2;
DEALLOCATE stmt123;
DEALLOCATE my_stmt_1;
DEALLOCATE test_stmt_99;

-- DEALLOCATE with quoted identifiers
DEALLOCATE "my statement";
DEALLOCATE "Statement Name";
DEALLOCATE "statement-with-dashes";
DEALLOCATE "statement.with.dots";
DEALLOCATE "statement with spaces";

-- DEALLOCATE with mixed case statement names
DEALLOCATE MyStatement;
DEALLOCATE STATEMENT_NAME;
DEALLOCATE TestStatement;
DEALLOCATE myStatement;
DEALLOCATE statementName;

-- DEALLOCATE with special characters in quoted names
DEALLOCATE "statement@special";
DEALLOCATE "statement#123";
DEALLOCATE "statement$name";
DEALLOCATE "statement%test";

-- DEALLOCATE with reserved words as statement names (when quoted)
DEALLOCATE "deallocate";
DEALLOCATE "prepare";
DEALLOCATE "transaction";
DEALLOCATE "work";
DEALLOCATE "table";

-- DEALLOCATE with single character statement names
DEALLOCATE a;
DEALLOCATE b;
DEALLOCATE c;
DEALLOCATE x;
DEALLOCATE y;
DEALLOCATE z;

-- DEALLOCATE with typical prepared statement naming patterns
DEALLOCATE select_query;
DEALLOCATE insert_stmt;
DEALLOCATE update_stmt;
DEALLOCATE delete_stmt;
DEALLOCATE report_query;

-- DEALLOCATE with business logic statement names
DEALLOCATE sales_report;
DEALLOCATE customer_lookup;
DEALLOCATE inventory_check;
DEALLOCATE order_summary;
DEALLOCATE revenue_analysis;

-- DEALLOCATE with procedural statement names
DEALLOCATE proc_stmt;
DEALLOCATE func_stmt;
DEALLOCATE temp_stmt;
DEALLOCATE result_stmt;
DEALLOCATE output_stmt;

-- DEALLOCATE with database operation naming patterns
DEALLOCATE users_select;
DEALLOCATE orders_insert;
DEALLOCATE items_update;
DEALLOCATE transactions_delete;
DEALLOCATE reports_query;

-- DEALLOCATE with PREPARE keyword variations
DEALLOCATE PREPARE select_customers;
DEALLOCATE prepare insert_orders;
DEALLOCATE PREPARE update_products;
DEALLOCATE prepare delete_temp_data;

-- Mixed case variations
Deallocate statement1;
DEALLOCATE statement2;
DeAlLoCaTe statement3;
deallocate STATEMENT4;
DEALLOCATE StatementName;

-- Complex statement names
DEALLOCATE complex_statement_name_with_underscores;
DEALLOCATE "Complex Statement Name With Spaces";
DEALLOCATE statement_123_test_abc;
DEALLOCATE "statement-name-with-dashes-123";

-- Typical prepared statement lifecycle examples
-- First prepare statements, then deallocate them
DEALLOCATE sales_summary_query;
DEALLOCATE monthly_report_stmt;
DEALLOCATE customer_analysis_prep;
DEALLOCATE inventory_lookup_stmt;

-- DEALLOCATE ALL variations
DEALLOCATE ALL;
deallocate all;
Deallocate All;
DEALLOCATE all;
deallocate ALL;

-- DEALLOCATE with PREPARE and ALL
DEALLOCATE PREPARE ALL;
deallocate prepare all;
DEALLOCATE prepare ALL;
deallocate PREPARE all;

-- Real-world examples based on typical prepared statement usage
DEALLOCATE get_user_by_id;
DEALLOCATE insert_new_order;
DEALLOCATE update_customer_info;
DEALLOCATE delete_expired_sessions;
DEALLOCATE calculate_monthly_revenue;
DEALLOCATE find_top_selling_products;
DEALLOCATE check_inventory_levels;
DEALLOCATE generate_sales_report;
DEALLOCATE validate_user_permissions;
DEALLOCATE process_payment_batch;

-- Examples with schema-qualified-like names (still just identifiers)
DEALLOCATE analytics_sales_query;
DEALLOCATE reporting_monthly_stmt;
DEALLOCATE etl_data_import_prep;
DEALLOCATE staging_validation_check;

-- Examples following typical SQL naming conventions
DEALLOCATE sp_get_customer_orders;
DEALLOCATE fn_calculate_discount;
DEALLOCATE vw_active_users_query;
DEALLOCATE tmp_processing_stmt;

-- Edge cases with valid identifier patterns
DEALLOCATE a1;
DEALLOCATE a1b2c3;
DEALLOCATE _private_stmt;
DEALLOCATE stmt_;
DEALLOCATE "special_name";
DEALLOCATE statement123abc;

-- Quoted edge cases
DEALLOCATE "123numeric_start";
DEALLOCATE "statement with \ttab";
DEALLOCATE "statement ""with"" quotes";
DEALLOCATE "CasE SeNsItIvE StAtEmEnT";
DEALLOCATE "statement@#$%^&*()";