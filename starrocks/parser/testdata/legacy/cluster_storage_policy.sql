CREATE STORAGE POLICY testPolicy
PROPERTIES(
    "storage_resource" = "s3",
    "cooldown_datetime" = "2022-06-08 00:00:00"
);
CREATE STORAGE POLICY testPolicy
PROPERTIES(
    "storage_resource" = "s3",
    "cooldown_ttl" = "1d"
);
DROP STORAGE POLICY policy1
