-- pin: default gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestQueryPostgresDollarQuote default#1
> raw
| {"op":"read","statement":"SELECT $$a;b$$ AS x;"}
< raw
| {"columns":[{"name":"x","type":{"code":"STRING"}}],"rows":[[{"type":{"code":"STRING"},"value":"a;b"}]]}
> raw
| {"op":"read","statement":"SELECT 2 AS y"}
< raw
| {"columns":[{"name":"y","type":{"code":"INT64"}}],"rows":[[{"type":{"code":"INT64"},"value":"2"}]]}
> raw
| {"op":"close"}
< raw
| {}
