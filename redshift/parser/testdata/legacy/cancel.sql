-- CANCEL statement examples for Redshift
-- https://docs.aws.amazon.com/redshift/latest/dg/r_CANCEL.html

-- Basic CANCEL statement
CANCEL 18764;

-- CANCEL with custom message
CANCEL 18764 'Long-running query terminated by administrator';

-- Multiple CANCEL statements
CANCEL 12345;
CANCEL 67890 'Query canceled due to resource constraints';
CANCEL 11111 'Maintenance window - canceling all active queries';