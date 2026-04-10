CREATE OR REPLACE PROCEDURE min_max_invoices_sp(
    minimum_price NUMBER(12,2),
    maximum_price NUMBER(12,2))
  RETURNS TABLE (id INTEGER, price NUMBER(12, 2))
  LANGUAGE SQL
AS
DECLARE
  rs RESULTSET;
  query VARCHAR DEFAULT 'SELECT * FROM invoices WHERE price > ? AND price < ?';
BEGIN
  rs := (EXECUTE IMMEDIATE :query USING (minimum_price, maximum_price));
  RETURN TABLE(rs);
END;
