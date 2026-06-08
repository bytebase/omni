-- FETCH command test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_FETCH.html

-- Basic FETCH statements (NEXT is default)
FETCH FROM lollapalooza;
fetch from lollapalooza;
FETCH NEXT FROM lollapalooza;
fetch next from lollapalooza;

-- FETCH with cursor names
FETCH FROM my_cursor;
FETCH FROM cursor_name;
FETCH FROM test_cursor;
FETCH FROM sales_cursor;
FETCH FROM product_cursor;

-- FETCH NEXT variations
FETCH NEXT FROM cursor1;
FETCH NEXT FROM my_cursor_name;
FETCH NEXT FROM test_cursor_123;

-- FETCH ALL variations
FETCH ALL FROM cursor1;
FETCH ALL FROM my_cursor;
FETCH ALL FROM sales_report_cursor;
fetch all from customer_cursor;

-- FETCH FORWARD with count
FETCH FORWARD 1 FROM lollapalooza;
FETCH FORWARD 5 FROM lollapalooza;
FETCH FORWARD 10 FROM my_cursor;
FETCH FORWARD 100 FROM sales_cursor;
FETCH FORWARD 1000 FROM product_cursor;

-- FETCH FORWARD 0 (current row)
FETCH FORWARD 0 FROM cursor1;
FETCH FORWARD 0 FROM my_cursor;

-- FETCH FORWARD ALL
FETCH FORWARD ALL FROM cursor1;
FETCH FORWARD ALL FROM my_cursor;
fetch forward all from sales_cursor;

-- Mixed case variations
Fetch From cursor1;
FETCH from cursor2;
fetch FROM cursor3;
Fetch Next From cursor4;
FETCH ALL from cursor5;
fetch forward 5 from cursor6;

-- FETCH with quoted cursor names
FETCH FROM "my cursor";
FETCH FROM "Cursor Name";
FETCH FROM "cursor-with-dashes";
FETCH FROM "cursor.with.dots";
FETCH FROM "cursor with spaces";
FETCH NEXT FROM "quoted cursor";
FETCH ALL FROM "Special Cursor Name";
FETCH FORWARD 10 FROM "cursor@special";

-- FETCH with complex cursor names
FETCH FROM cursor_with_long_name;
FETCH FROM complex_cursor_name_with_underscores;
FETCH FROM cursor_123_test_abc;
FETCH NEXT FROM monthly_summary_cursor;
FETCH ALL FROM customer_detail_cursor;
FETCH FORWARD 50 FROM transaction_history_cursor;

-- Examples from AWS documentation
-- Fetch the first 5 rows in the cursor lollapalooza
FETCH FORWARD 5 FROM lollapalooza;

-- Fetch the next row
FETCH NEXT FROM lollapalooza;

-- Common cursor usage patterns
FETCH FROM emp_cursor;
FETCH NEXT FROM dept_cursor;
FETCH ALL FROM sales_cursor;
FETCH FORWARD 20 FROM product_cursor;
FETCH FORWARD ALL FROM customer_cursor;

-- Procedural cursor patterns
FETCH FROM proc_cursor;
FETCH NEXT FROM func_cursor;
FETCH ALL FROM temp_cursor;
FETCH FORWARD 1 FROM result_cursor;
FETCH FORWARD ALL FROM output_cursor;

-- Database object cursor patterns
FETCH FROM users_cur;
FETCH NEXT FROM orders_cur;
FETCH ALL FROM items_cur;
FETCH FORWARD 100 FROM transactions_cur;
FETCH FORWARD ALL FROM reports_cur;

-- Single character cursor names
FETCH FROM a;
FETCH NEXT FROM b;
FETCH ALL FROM c;
FETCH FORWARD 1 FROM x;
FETCH FORWARD 5 FROM y;
FETCH FORWARD ALL FROM z;

-- Numeric suffixes in cursor names
FETCH FROM cursor1;
FETCH NEXT FROM cursor2;
FETCH ALL FROM cursor123;
FETCH FORWARD 10 FROM my_cursor_1;
FETCH FORWARD ALL FROM test_cursor_99;

-- Reserved words as quoted cursor names
FETCH FROM "fetch";
FETCH FROM "next";
FETCH FROM "all";
FETCH FROM "forward";
FETCH FROM "from";
FETCH FROM "cursor";
FETCH FROM "table";
FETCH FROM "select";

-- Mixed case cursor names
FETCH FROM MyCursor;
FETCH NEXT FROM CURSOR_NAME;
FETCH ALL FROM TestCursor;
FETCH FORWARD 5 FROM myCursor;
FETCH FORWARD ALL FROM cursorName;

-- Special characters in quoted cursor names
FETCH FROM "cursor@special";
FETCH NEXT FROM "cursor#123";
FETCH ALL FROM "cursor$name";
FETCH FORWARD 10 FROM "cursor%test";
FETCH FORWARD ALL FROM "cursor&data";

-- Typical business cursor patterns
FETCH FROM customer_orders_cursor;
FETCH NEXT FROM inventory_items_cursor;
FETCH ALL FROM financial_reports_cursor;
FETCH FORWARD 25 FROM employee_details_cursor;
FETCH FORWARD ALL FROM product_sales_cursor;

-- Time-based cursor patterns
FETCH FROM daily_report_cursor;
FETCH NEXT FROM weekly_summary_cursor;
FETCH ALL FROM monthly_stats_cursor;
FETCH FORWARD 7 FROM quarterly_data_cursor;
FETCH FORWARD ALL FROM yearly_totals_cursor;

-- Large count values (up to 1000 for single-node clusters)
FETCH FORWARD 999 FROM large_dataset_cursor;
FETCH FORWARD 1000 FROM max_count_cursor;

-- Edge cases with whitespace and formatting
FETCH    FROM    cursor1;
FETCH NEXT    FROM    cursor2;
FETCH    ALL FROM cursor3;
FETCH FORWARD    5    FROM cursor4;
FETCH    FORWARD ALL FROM cursor5;