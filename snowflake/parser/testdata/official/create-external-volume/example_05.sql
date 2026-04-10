SHOW ICEBERG TABLES
  ->> SELECT *
        FROM $1
        WHERE "external_volume_name" = 'my_external_volume_1';
