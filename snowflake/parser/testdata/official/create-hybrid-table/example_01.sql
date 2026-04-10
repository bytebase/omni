CREATE HYBRID TABLE mytable (
  customer_id INT AUTOINCREMENT PRIMARY KEY,
  full_name VARCHAR(255),
  email VARCHAR(255) UNIQUE,
  extended_customer_info VARIANT,
  INDEX index_full_name (full_name)
);
