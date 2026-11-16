-- pin: driver docker.elastic.co/elasticsearch/elasticsearch:8.19.25 sha256:5bfdf4313b9bc05c0d46fe219d6e67710eb0c7654ed6c275167bdf9beefc56f0 8.19.25

== TestPing driver#1
> raw
| {"method":"GET","route":"/"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"version\":{\"number\":\"8.19.25\"}}"}
