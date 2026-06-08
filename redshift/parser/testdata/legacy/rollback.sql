-- rollback.sql
-- Test cases for ROLLBACK command in Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_ROLLBACK.html

-- Basic ROLLBACK
ROLLBACK;

-- ROLLBACK with WORK keyword
ROLLBACK WORK;

-- ROLLBACK with TRANSACTION keyword
ROLLBACK TRANSACTION;

-- Case variations
rollback;
Rollback;
ROLLBACK;

rollback work;
Rollback Work;
ROLLBACK WORK;

rollback transaction;
Rollback Transaction;
ROLLBACK TRANSACTION;

-- Full example from AWS documentation
-- Create a table
create table movie_gross( name varchar(30), gross bigint );

-- Start a transaction
begin;

-- Insert data
insert into movie_gross values ( 'Raiders of the Lost Ark', 23400000);
insert into movie_gross values ( 'Star Wars', 10000000 );

-- Roll back the transaction (discards all updates)
rollback;

-- Another example with WORK keyword
begin;
insert into movie_gross values ( 'The Empire Strikes Back', 20000000);
insert into movie_gross values ( 'Return of the Jedi', 15000000 );
rollback work;

-- Example with TRANSACTION keyword
begin;
insert into movie_gross values ( 'Indiana Jones and the Temple of Doom', 25000000);
insert into movie_gross values ( 'Indiana Jones and the Last Crusade', 19700000 );
rollback transaction;

-- ROLLBACK in a more complex transaction
begin;

-- Multiple operations
create table temp_sales (id int, amount decimal(10,2));
insert into temp_sales values (1, 100.50);
insert into temp_sales values (2, 250.75);
update temp_sales set amount = amount * 1.1 where id = 1;
delete from temp_sales where amount < 200;

-- Roll back all changes
rollback;

-- ROLLBACK after DDL operations
begin;
create table test_rollback (col1 int, col2 varchar(50));
alter table test_rollback add column col3 date;
rollback;

-- Note: ROLLBACK is functionally equivalent to ABORT
-- Both commands stop the current transaction and discard all updates

-- Test with explicit transaction blocks
BEGIN TRANSACTION;
-- Some operations would go here
ROLLBACK TRANSACTION;

BEGIN WORK;
-- Some operations would go here
ROLLBACK WORK;