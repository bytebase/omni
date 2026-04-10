SHOW TABLES IN ACCOUNT ->> SELECT "database_name" || '.' || "schema_name" || '.' || "name" AS fully_qualified_name FROM $1 ORDER BY fully_qualified_name;
