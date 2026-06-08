-- ALTER FUNCTION examples for Amazon Redshift
-- Based on: https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_FUNCTION.html

-- Example 1: Basic function rename
-- Rename a function with three parameters
ALTER FUNCTION first_quarter_revenue(bigint, numeric, int)
RENAME TO quarterly_revenue;

-- Example 2: Function rename with schema qualification
-- Rename a function in a specific schema
ALTER FUNCTION sales.calculate_commission(numeric, numeric)
RENAME TO compute_commission;

-- Example 3: Change function owner to a specific user
ALTER FUNCTION quarterly_revenue(bigint, numeric)
OWNER TO etl_user;

-- Example 4: Change function owner to CURRENT_USER
ALTER FUNCTION process_data(varchar, timestamp)
OWNER TO CURRENT_USER;

-- Example 5: Change function owner to SESSION_USER
ALTER FUNCTION analytics.generate_report(date, date, varchar)
OWNER TO SESSION_USER;

-- Example 6: Python UDF with named parameters - rename
-- Note: Python UDFs will be deprecated after November 1, 2025
ALTER FUNCTION calculate_tax(amount numeric, rate numeric)
RENAME TO compute_tax;

-- Example 7: Function with complex types
ALTER FUNCTION transform_json(json_data varchar(max), options varchar)
OWNER TO data_team;

-- Example 8: Function with no parameters
ALTER FUNCTION get_current_timestamp()
RENAME TO fetch_current_time;

-- Example 9: Multiple functions with same name but different signatures
-- First function with two parameters
ALTER FUNCTION process_order(order_id bigint, status varchar)
RENAME TO handle_order;

-- Second function with three parameters (different signature)
ALTER FUNCTION process_order(order_id bigint, status varchar, priority int)
RENAME TO handle_priority_order;

-- Example 10: Schema-qualified function with owner change
ALTER FUNCTION warehouse.calculate_inventory(product_id int, warehouse_id int, as_of_date date)
OWNER TO inventory_manager;

-- Example 11: Function with decimal precision types
ALTER FUNCTION finance.calculate_interest(principal decimal(15,2), rate decimal(5,4), years int)
RENAME TO compute_interest;

-- Example 12: Function with character varying types
ALTER FUNCTION marketing.parse_campaign(campaign_data varchar(500), delimiter char(1))
OWNER TO marketing_team;

-- Example 13: Function with timestamp types
ALTER FUNCTION log_analytics.parse_log_entry(log_line text, log_timestamp timestamp without time zone)
RENAME TO analyze_log_entry;

-- Example 14: Function with boolean return type
ALTER FUNCTION validate_email(email varchar(255))
OWNER TO validation_service;

-- Example 15: Complex multi-schema scenario
-- Function in one schema being renamed
ALTER FUNCTION staging.clean_customer_data(raw_data text, validation_rules json)
RENAME TO sanitize_customer_data;

-- Then changing owner
ALTER FUNCTION staging.sanitize_customer_data(raw_data text, validation_rules json)
OWNER TO data_quality_team;