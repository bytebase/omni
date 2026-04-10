CREATE TABLE departments (department_id INTEGER, name VARCHAR);
CREATE TABLE employees (employee_ID INTEGER, last_name VARCHAR,
  department_ID INTEGER, project_names ARRAY);

INSERT INTO departments (department_ID, name) VALUES
  (1, 'Engineering'),
  (2, 'Support');
INSERT INTO employees (employee_ID, last_name, department_ID) VALUES
  (101, 'Richards', 1),
  (102, 'Paulson',  1),
  (103, 'Johnson',  2);
