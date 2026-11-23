CREATE TABLE customers (id NUMBER(10) CONSTRAINT customers_pk PRIMARY KEY, name VARCHAR2(100) NOT NULL, region VARCHAR2(10));
CREATE TABLE orders (id NUMBER(10) CONSTRAINT orders_pk PRIMARY KEY, customer_id NUMBER(10) CONSTRAINT orders_customer_fk REFERENCES customers (id), amount NUMBER(10,2), status VARCHAR2(20));
CREATE TABLE audit_log (id NUMBER(10), what VARCHAR2(100), at DATE);
CREATE VIEW v_customers AS SELECT id, name FROM customers;
CREATE SEQUENCE order_seq;
CREATE SYNONYM cust FOR customers;
CREATE FUNCTION f_double (p_x IN NUMBER) RETURN NUMBER IS
BEGIN
  RETURN p_x * 2;
END;
/
CREATE PROCEDURE p_existing (p_id IN NUMBER) IS
  v_name customers.name%TYPE;
BEGIN
  SELECT name INTO v_name FROM customers WHERE id = p_id;
END;
/
CREATE PACKAGE pk AS
  g_count NUMBER := 0;
  PROCEDURE log_it (p_what IN VARCHAR2);
  FUNCTION total (p_customer IN NUMBER) RETURN NUMBER;
END pk;
/
CREATE PACKAGE BODY pk AS
  PROCEDURE log_it (p_what IN VARCHAR2) IS
  BEGIN
    INSERT INTO audit_log (id, what, at) VALUES (1, p_what, SYSDATE);
    g_count := g_count + 1;
  END;
  FUNCTION total (p_customer IN NUMBER) RETURN NUMBER IS
    v_total NUMBER;
  BEGIN
    SELECT SUM(amount) INTO v_total FROM orders WHERE customer_id = p_customer;
    RETURN v_total;
  END;
END pk;
/
CREATE PROCEDURE p_calls IS
  v NUMBER;
BEGIN
  v := f_double(1);
  p_existing(v);
  pk.log_it('called');
END;
/
CREATE TYPE t_point AS OBJECT (x NUMBER, y NUMBER, MEMBER FUNCTION norm RETURN NUMBER);
/
CREATE TYPE BODY t_point AS
  MEMBER FUNCTION norm RETURN NUMBER IS
  BEGIN
    RETURN SQRT(x * x + y * y);
  END;
END;
/
CREATE TRIGGER trg_existing BEFORE INSERT ON orders FOR EACH ROW
BEGIN
  :NEW.status := NVL(:NEW.status, 'new');
END;
/
