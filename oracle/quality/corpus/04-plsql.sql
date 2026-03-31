-- @name: anonymous block
-- @valid: true
-- @source: Oracle SQL Reference
BEGIN
  NULL;
END;

-- @name: block with variable
-- @valid: true
-- @source: Oracle SQL Reference
DECLARE
  v_count NUMBER;
BEGIN
  SELECT COUNT(*) INTO v_count FROM employees;
END;

-- @name: block with if
-- @valid: true
-- @source: Oracle SQL Reference
DECLARE
  v_sal NUMBER;
BEGIN
  SELECT salary INTO v_sal FROM employees WHERE employee_id = 100;
  IF v_sal > 10000 THEN
    NULL;
  ELSIF v_sal > 5000 THEN
    NULL;
  ELSE
    NULL;
  END IF;
END;

-- @name: block with loop
-- @valid: true
-- @source: Oracle SQL Reference
BEGIN
  FOR i IN 1..10 LOOP
    NULL;
  END LOOP;
END;

-- @name: block with cursor loop
-- @valid: true
-- @source: Oracle SQL Reference
DECLARE
  CURSOR c_emp IS SELECT employee_id, salary FROM employees;
BEGIN
  FOR rec IN c_emp LOOP
    NULL;
  END LOOP;
END;

-- @name: block with exception
-- @valid: true
-- @source: Oracle SQL Reference
BEGIN
  SELECT 1 INTO :x FROM dual;
EXCEPTION
  WHEN NO_DATA_FOUND THEN
    NULL;
  WHEN OTHERS THEN
    NULL;
END;

-- @name: create procedure
-- @valid: true
-- @source: Oracle SQL Reference
CREATE PROCEDURE update_salary(
  p_emp_id IN NUMBER,
  p_pct IN NUMBER DEFAULT 10
) IS
BEGIN
  UPDATE employees SET salary = salary * (1 + p_pct/100)
  WHERE employee_id = p_emp_id;
END;

-- @name: create function
-- @valid: true
-- @source: Oracle SQL Reference
CREATE FUNCTION get_dept_count(p_dept_id IN NUMBER) RETURN NUMBER IS
  v_count NUMBER;
BEGIN
  SELECT COUNT(*) INTO v_count FROM employees WHERE department_id = p_dept_id;
  RETURN v_count;
END;

-- @name: create package spec
-- @valid: true
-- @source: Oracle SQL Reference
CREATE PACKAGE emp_pkg IS
  PROCEDURE hire_employee(p_name VARCHAR2, p_dept NUMBER);
  FUNCTION get_salary(p_emp_id NUMBER) RETURN NUMBER;
END emp_pkg;

-- @name: create package body
-- @valid: true
-- @source: Oracle SQL Reference
CREATE PACKAGE BODY emp_pkg IS
  PROCEDURE hire_employee(p_name VARCHAR2, p_dept NUMBER) IS
  BEGIN
    NULL;
  END;
  FUNCTION get_salary(p_emp_id NUMBER) RETURN NUMBER IS
  BEGIN
    RETURN 0;
  END;
END emp_pkg;

-- @name: create trigger
-- @valid: true
-- @source: Oracle SQL Reference
CREATE TRIGGER trg_emp_audit
BEFORE INSERT OR UPDATE ON employees
FOR EACH ROW
BEGIN
  :NEW.created_date := SYSDATE;
END;
