-- CLOSE command test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_CLOSE.html

-- Basic CLOSE statements
CLOSE cursor1;
close cursor1;

-- CLOSE with simple cursor names
CLOSE my_cursor;
CLOSE cursor_name;
CLOSE test_cursor;

-- CLOSE with underscore in cursor names
CLOSE my_cursor_name;
CLOSE cursor_with_long_name;
CLOSE test_cursor_123;

-- CLOSE with numeric suffixes
CLOSE cursor1;
CLOSE cursor2;
CLOSE cursor123;
CLOSE my_cursor_1;
CLOSE test_cursor_99;

-- CLOSE with quoted identifiers
CLOSE "my cursor";
CLOSE "Cursor Name";
CLOSE "cursor-with-dashes";
CLOSE "cursor.with.dots";
CLOSE "cursor with spaces";

-- CLOSE with mixed case cursor names
CLOSE MyCursor;
CLOSE CURSOR_NAME;
CLOSE TestCursor;
CLOSE myCursor;
CLOSE cursorName;

-- CLOSE with special characters in quoted names
CLOSE "cursor@special";
CLOSE "cursor#123";
CLOSE "cursor$name";
CLOSE "cursor%test";

-- CLOSE with reserved words as cursor names (when quoted)
CLOSE "close";
CLOSE "begin";
CLOSE "transaction";
CLOSE "work";
CLOSE "table";

-- CLOSE with single character cursor names
CLOSE a;
CLOSE b;
CLOSE c;
CLOSE x;
CLOSE y;
CLOSE z;

-- CLOSE with typical cursor naming patterns
CLOSE emp_cursor;
CLOSE dept_cursor;
CLOSE sales_cursor;
CLOSE product_cursor;
CLOSE customer_cursor;

-- CLOSE with procedural cursor names
CLOSE proc_cursor;
CLOSE func_cursor;
CLOSE temp_cursor;
CLOSE result_cursor;
CLOSE output_cursor;

-- CLOSE with database object naming patterns
CLOSE users_cur;
CLOSE orders_cur;
CLOSE items_cur;
CLOSE transactions_cur;
CLOSE reports_cur;

-- Example from typical cursor usage patterns
close sales_report_cursor;
close monthly_summary_cursor;
close customer_detail_cursor;

-- Mixed case variations
Close cursor1;
CLOSE cursor2;
ClOsE cursor3;
close CURSOR4;
CLOSE CursorName;

-- Complex cursor names
CLOSE complex_cursor_name_with_underscores;
CLOSE "Complex Cursor Name With Spaces";
CLOSE cursor_123_test_abc;
CLOSE "cursor-name-with-dashes-123";