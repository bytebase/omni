DROP DATABASE IF EXISTS omni_simulation_other;
CREATE DATABASE omni_simulation_other;
CREATE TABLE omni_simulation_other.archived (id INT PRIMARY KEY);
CREATE TABLE accounts (id INT PRIMARY KEY, name VARCHAR(20));
CREATE TABLE ledger (id INT PRIMARY KEY, account_id INT, CONSTRAINT fk_account FOREIGN KEY (account_id) REFERENCES accounts (id));
CREATE TABLE entries (id INT PRIMARY KEY, archived_id INT, CONSTRAINT fk_archived FOREIGN KEY (archived_id) REFERENCES omni_simulation_other.archived (id));
CREATE VIEW v_ledger AS SELECT id, account_id FROM ledger
