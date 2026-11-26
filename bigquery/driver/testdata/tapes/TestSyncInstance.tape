-- pin: driver gcp:bigquery:bytebase-dev:US sha256:92a85df45321648c0feb59ee0c222e8156cb9332d322ddb9f15eaf37c8c12533 2026-11-26

== TestSyncInstance driver#1
> raw
| {"op":"ping"}
< raw
| {}
> raw
| {"op":"listDatasets"}
< raw
| {"datasets":["ds1","ds2","dwh","hello","omni_record","verizon_demo","zp"]}
