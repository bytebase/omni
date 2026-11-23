CREATE TABLE customers (id INT PRIMARY KEY, email VARCHAR(100) NOT NULL, UNIQUE KEY uk_email (email));
CREATE TABLE orders (id INT PRIMARY KEY, customer_id INT, email VARCHAR(100), parent_id INT, KEY idx_customer (customer_id), CONSTRAINT fk_customer FOREIGN KEY (customer_id) REFERENCES customers (id) ON DELETE CASCADE ON UPDATE SET NULL, CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES orders (id));
CREATE TABLE items (order_id INT NOT NULL, line INT NOT NULL, PRIMARY KEY (order_id, line), CONSTRAINT fk_items_order FOREIGN KEY (order_id) REFERENCES orders (id))
