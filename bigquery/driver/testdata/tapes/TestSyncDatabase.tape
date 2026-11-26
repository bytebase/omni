-- pin: driver gcp:bigquery:bytebase-dev:US sha256:92a85df45321648c0feb59ee0c222e8156cb9332d322ddb9f15eaf37c8c12533 2026-11-26

== TestSyncDatabase driver#1
> raw
| {"op":"ping"}
< raw
| {}
> raw
| {"op":"datasetMetadata","dataset":"omni_record"}
< raw
| {"dataset":{"fullID":"bytebase-dev:omni_record","location":"US"}}
> raw
| {"op":"read","dataset":"omni_record","statement":"\n\t\tSELECT\n\t\t\ttable_name,\n\t\t\tcolumn_name,\n\t\t\tordinal_position,\n\t\t\tis_nullable,\n\t\t\tdata_type,\n\t\t\tcollation_name,\n\t\t\tcolumn_default\n\t\tFROM `omni_record`.INFORMATION_SCHEMA.COLUMNS ORDER BY table_name, ordinal_position;"}
< raw
| {"jobID":"job-1","schema":[{"name":"table_name","type":"STRING"},{"name":"column_name","type":"STRING"},{"name":"ordinal_position","type":"INTEGER"},{"name":"is_nullable","type":"STRING"},{"name":"data_type","type":"STRING"},{"name":"collation_name","type":"STRING"},{"name":"column_default","type":"STRING"}],"rows":[[{"k":"string","s":"customers"},{"k":"string","s":"id"},{"k":"int64","i":1},{"k":"string","s":"NO"},{"k":"string","s":"INT64"},{"k":"string","s":"NULL"},{"k":"string","s":"NULL"}],[{"k":"string","s":"customers"},{"k":"string","s":"city"},{"k":"int64","i":2},{"k":"string","s":"YES"},{"k":"string","s":"STRING"},{"k":"string","s":"NULL"},{"k":"string","s":"NULL"}],[{"k":"string","s":"customers"},{"k":"string","s":"signed_up"},{"k":"int64","i":3},{"k":"string","s":"YES"},{"k":"string","s":"TIMESTAMP"},{"k":"string","s":"NULL"},{"k":"string","s":"NULL"}],[{"k":"string","s":"customers"},{"k":"string","s":"score"},{"k":"int64","i":4},{"k":"string","s":"YES"},{"k":"string","s":"NUMERIC"},{"k":"string","s":"NULL"},{"k":"string","s":"NULL"}],[{"k":"string","s":"customers"},{"k":"string","s":"tags"},{"k":"int64","i":5},{"k":"string","s":"NO"},{"k":"string","s":"ARRAY\u003cSTRING\u003e"},{"k":"string","s":"NULL"},{"k":"string","s":"NULL"}],[{"k":"string","s":"partitioned"},{"k":"string","s":"_PARTITIONTIME"},{"k":"null"},{"k":"string","s":"NO"},{"k":"string","s":"TIMESTAMP"},{"k":"null"},{"k":"null"}],[{"k":"string","s":"partitioned"},{"k":"string","s":"id"},{"k":"int64","i":1},{"k":"string","s":"YES"},{"k":"string","s":"INT64"},{"k":"string","s":"NULL"},{"k":"string","s":"NULL"}],[{"k":"string","s":"recent"},{"k":"string","s":"id"},{"k":"int64","i":1},{"k":"string","s":"YES"},{"k":"string","s":"INT64"},{"k":"string","s":"NULL"},{"k":"string","s":"NULL"}],[{"k":"string","s":"recent"},{"k":"string","s":"city"},{"k":"int64","i":2},{"k":"string","s":"YES"},{"k":"string","s":"STRING"},{"k":"string","s":"NULL"},{"k":"string","s":"NULL"}]]}
> raw
| {"op":"read","dataset":"omni_record","statement":"\n\t\tSELECT\n\t\t\ttable_name,\n\t\t\tview_definition\n\t\tFROM `omni_record`.INFORMATION_SCHEMA.VIEWS ORDER BY table_name;"}
< raw
| {"jobID":"job-2","schema":[{"name":"table_name","type":"STRING"},{"name":"view_definition","type":"STRING"}],"rows":[[{"k":"string","s":"recent"},{"k":"string","s":"SELECT id, city FROM `omni_record.customers` WHERE id \u003e 1"}]]}
> raw
| {"op":"listTables","dataset":"omni_record"}
< raw
| {"tables":["customers","partitioned","recent"]}
> raw
| {"op":"tableMetadata","dataset":"omni_record","table":"customers"}
< raw
| {"table":{"description":"a synced comment","numRows":3,"numBytes":113}}
> raw
| {"op":"tableMetadata","dataset":"omni_record","table":"partitioned"}
< raw
| {"table":{}}
> raw
| {"op":"tableMetadata","dataset":"omni_record","table":"recent"}
< raw
| {"table":{}}
