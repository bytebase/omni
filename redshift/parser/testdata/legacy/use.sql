-- USE command test cases
-- Changes the current database
-- https://docs.aws.amazon.com/redshift/latest/dg/r_USE.html

-- Basic USE command
USE mydb;
USE database1;
USE dev;
USE prod;
USE test;

-- Database names with underscores
USE my_database;
USE test_db;
USE dev_environment;
USE prod_environment;

-- Database names with numbers
USE db1;
USE database123;
USE test_2023;
USE dev_v2;

-- Mixed case database names (Redshift is case-insensitive by default)
USE MyDatabase;
USE TestDB;
USE DevEnvironment;

-- Longer database names
USE production_analytics_database;
USE development_testing_environment;
USE staging_application_database;

-- Common database naming patterns
USE dwh;
USE data_warehouse;
USE analytics;
USE reporting;
USE staging;
USE sandbox;

-- Database names that follow common conventions
USE raw_data;
USE processed_data;
USE curated_data;
USE mart_sales;
USE mart_finance;
USE mart_marketing;

-- Edge cases with valid identifiers
USE db;
USE d;
USE database_with_very_long_name_but_still_valid;

-- Real-world examples
USE ecommerce;
USE customer_360;
USE revenue_analytics;
USE user_behavior;
USE product_catalog;
USE inventory_management;