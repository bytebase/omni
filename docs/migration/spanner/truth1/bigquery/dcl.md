# BigQuery GoogleSQL truth1 — DCL Syntax

All forms extracted from:
https://cloud.google.com/bigquery/docs/reference/standard-sql/data-control-language

---

## DCL-001: GRANT statement

```
GRANT role_list
  ON resource_type resource_name
  TO user_list

role_list:
  role [, role ...]

resource_type:
  SCHEMA | TABLE | VIEW | EXTERNAL TABLE | PROJECT

user_list:
  user_identifier [, user_identifier ...]

-- User identifier formats:
-- "user:email@domain.com"
-- "group:group@domain.com"
-- "serviceAccount:name@project.iam.gserviceaccount.com"
-- "domain:example.com"
-- "specialGroup:allAuthenticatedUsers"
-- "specialGroup:allUsers"
-- "connection:[project_id.]location.connection_id"
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-control-language#grant_statement
version: bigquery-current

```sql
GRANT `roles/bigquery.dataViewer` ON SCHEMA `myProject`.myDataset
TO "user:raha@example-pet-store.com", "user:sasha@example-pet-store.com";

GRANT `roles/aiplatform.user`, `roles/run.invoker`
ON PROJECT `my-vertex-project`
TO "connection:my-bq-project.us.my-connection", "connection:another-bq-project.eu.other-connection";
```

---

## DCL-002: REVOKE statement

```
REVOKE role_list
  ON resource_type resource_name
  FROM user_list

-- resource_type: SCHEMA | TABLE | VIEW | EXTERNAL TABLE | PROJECT
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-control-language#revoke_statement
version: bigquery-current

```sql
REVOKE `roles/bigquery.admin` ON SCHEMA `myProject`.myDataset
FROM "group:example-team@example-pet-store.com", "serviceAccount:user@test-project.iam.gserviceaccount.com";

REVOKE `roles/run.invoker`
ON PROJECT `my-vertex-project`
FROM "connection:my-bq-project.us.my-connection";
```
