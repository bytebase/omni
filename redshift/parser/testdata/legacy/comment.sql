-- COMMENT command test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_COMMENT.html

-- Basic TABLE comments (from AWS documentation examples)
COMMENT ON TABLE sales IS 'This table stores tickets sales data';
COMMENT ON TABLE event IS 'Contains listings of individual events.';

-- Remove table comment using NULL
COMMENT ON TABLE sales IS NULL;

-- COLUMN comments (from AWS documentation examples)
COMMENT ON COLUMN sales.eventid IS 'Foreign-key reference to the EVENT table.';

-- More COLUMN comment examples
COMMENT ON COLUMN event.eventname IS 'Name of the event';
COMMENT ON COLUMN sales.qtysold IS 'Quantity of tickets sold';
COMMENT ON COLUMN sales.pricepaid IS 'Total price paid for tickets';

-- Remove column comment
COMMENT ON COLUMN sales.eventid IS NULL;

-- DATABASE comments
COMMENT ON DATABASE dev IS 'Development database';
COMMENT ON DATABASE prod IS 'Production database';
COMMENT ON DATABASE test IS 'Testing environment database';

-- Remove database comment
COMMENT ON DATABASE dev IS NULL;

-- VIEW comments
COMMENT ON VIEW sales_summary IS 'Summary view of sales data';
COMMENT ON VIEW monthly_revenue IS 'Monthly revenue aggregation view';

-- Remove view comment
COMMENT ON VIEW sales_summary IS NULL;

-- SCHEMA comments
COMMENT ON SCHEMA public IS 'Default public schema';
COMMENT ON SCHEMA analytics IS 'Schema for analytics tables';
COMMENT ON SCHEMA staging IS 'Staging area for data loading';

-- Remove schema comment
COMMENT ON SCHEMA analytics IS NULL;

-- CONSTRAINT comments
COMMENT ON CONSTRAINT sales_pkey ON sales IS 'Primary key constraint';
COMMENT ON CONSTRAINT fk_sales_event ON sales IS 'Foreign key to event table';
COMMENT ON CONSTRAINT chk_qty_positive ON sales IS 'Check constraint for positive quantity';

-- Remove constraint comment
COMMENT ON CONSTRAINT sales_pkey ON sales IS NULL;

-- Schema-qualified object names
COMMENT ON TABLE public.sales IS 'Sales table in public schema';
COMMENT ON TABLE analytics.monthly_sales IS 'Monthly sales aggregation';
COMMENT ON COLUMN public.event.eventname IS 'Event name in public schema';
COMMENT ON VIEW analytics.revenue_view IS 'Revenue view in analytics schema';

-- Quoted identifiers
COMMENT ON TABLE "Sales Data" IS 'Table with quoted name';
COMMENT ON COLUMN "Sales Data"."Event ID" IS 'Column with quoted names';
COMMENT ON SCHEMA "Analytics Schema" IS 'Schema with quoted name';

-- Mixed case variations
comment on table Sales IS 'Mixed case comment';
Comment On Table EVENT is 'Mixed case keywords';
COMMENT on column sales.eventid is 'Mixed case column comment';

-- Comments with special characters and escapes
COMMENT ON TABLE sales IS 'This table stores ''quoted'' text';
COMMENT ON TABLE event IS 'Table with "double quotes" in comment';
COMMENT ON TABLE venue IS 'Table with newline
and continuation';

-- Long comments
COMMENT ON TABLE sales IS 'This is a very long comment that describes the sales table in great detail, including information about its purpose, structure, data sources, update frequency, and business rules that apply to the data stored within it';

-- Empty comments (effectively removes comment)
COMMENT ON TABLE sales IS '';

-- Complex schema and table names
COMMENT ON TABLE my_schema.my_table IS 'Table in custom schema';
COMMENT ON COLUMN my_schema.my_table.my_column IS 'Column in custom schema';

-- Numeric and special characters in names (quoted)
COMMENT ON TABLE "Table123" IS 'Table with numbers';
COMMENT ON COLUMN "Table_with_underscores".col_name IS 'Column with underscores';

-- Real-world example scenarios
COMMENT ON TABLE customer IS 'Customer master data table containing demographic and contact information';
COMMENT ON TABLE orders IS 'Order transaction table with purchase details and timestamps';
COMMENT ON TABLE order_items IS 'Line items for each order with product details and quantities';

COMMENT ON COLUMN customer.customer_id IS 'Unique identifier for customer (primary key)';
COMMENT ON COLUMN customer.email IS 'Customer email address (unique)';
COMMENT ON COLUMN customer.created_date IS 'Date when customer record was created';

COMMENT ON COLUMN orders.order_id IS 'Unique order identifier (primary key)';
COMMENT ON COLUMN orders.customer_id IS 'Foreign key reference to customer table';
COMMENT ON COLUMN orders.order_date IS 'Date when order was placed';
COMMENT ON COLUMN orders.total_amount IS 'Total order amount in USD';

-- Remove all the example comments
COMMENT ON TABLE customer IS NULL;
COMMENT ON TABLE orders IS NULL;
COMMENT ON TABLE order_items IS NULL;
COMMENT ON COLUMN customer.customer_id IS NULL;
COMMENT ON COLUMN orders.order_id IS NULL;