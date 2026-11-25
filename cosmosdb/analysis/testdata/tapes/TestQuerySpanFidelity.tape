-- pin: emulator mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:2.14.28 sha256:9a4f69121f20c1d6502db84625c71ebb441ab61b01880a7b2d3d8031c6dac429 2.14.28

== TestQuerySpanFidelity emulator#1
> raw
| {"op":"load","fixture":"document","documents":"{\"id\":\"«$.id»\",\"pk\":\"«$.pk»\",\"name\":\"«$.name»\",\"nick\":\"«$.nick»\",\"n\":1101,\"address\":{\"city\":\"«$.address.city»\",\"zip\":\"«$.address.zip»\",\"geo\":{\"lat\":\"«$.address.geo.lat»\"}},\"tags\":[{\"name\":\"«$.tags[0].name»\",\"skills\":[\"«$.tags[0].skills[0]»\",\"«$.tags[0].skills[1]»\"]},{\"name\":\"«$.tags[1].name»\",\"skills\":[\"«$.tags[1].skills[0]»\"]}],\"phones\":[\"«$.phones[0]»\",\"«$.phones[1]»\"]}"}
< raw
| {}
> raw
| {"op":"query","fixture":"document","statement":"SELECT * FROM c"}
< raw
| {"items":[{"id":"«$.id»","pk":"«$.pk»","name":"«$.name»","nick":"«$.nick»","n":1101,"address":{"city":"«$.address.city»","zip":"«$.address.zip»","geo":{"lat":"«$.address.geo.lat»"}},"tags":[{"name":"«$.tags[0].name»","skills":["«$.tags[0].skills[0]»","«$.tags[0].skills[1]»"]},{"name":"«$.tags[1].name»","skills":["«$.tags[1].skills[0]»"]}],"phones":["«$.phones[0]»","«$.phones[1]»"],"_rid":"(system)","_self":"(system)","_etag":"(system)","_attachments":"(system)","_ts":"(system)"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.name FROM c"}
< raw
| {"items":[{"name":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.name, c.address.city FROM c"}
< raw
| {"items":[{"name":"«$.name»","city":"«$.address.city»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.address FROM c"}
< raw
| {"items":[{"address":{"city":"«$.address.city»","zip":"«$.address.zip»","geo":{"lat":"«$.address.geo.lat»"}}}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.address.geo.lat FROM c"}
< raw
| {"items":[{"lat":"«$.address.geo.lat»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.tags FROM c"}
< raw
| {"items":[{"tags":[{"name":"«$.tags[0].name»","skills":["«$.tags[0].skills[0]»","«$.tags[0].skills[1]»"]},{"name":"«$.tags[1].name»","skills":["«$.tags[1].skills[0]»"]}]}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.id, c.pk FROM c"}
< raw
| {"items":[{"id":"«$.id»","pk":"«$.pk»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c._ts FROM c"}
< raw
| {"items":[{"_ts":"(system)"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c[\"name\"] FROM c"}
< raw
| {"items":[{"name":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.address[\"city\"] FROM c"}
< raw
| {"items":[{"city":"«$.address.city»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c[\"address\"][\"geo\"] FROM c"}
< raw
| {"items":[{"geo":{"lat":"«$.address.geo.lat»"}}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c[\"tags\"][1][\"name\"] FROM c"}
< raw
| {"items":[{"name":"«$.tags[1].name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.tags[0] FROM c"}
< raw
| {"items":[{"$1":{"name":"«$.tags[0].name»","skills":["«$.tags[0].skills[0]»","«$.tags[0].skills[1]»"]}}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.tags[0].name FROM c"}
< raw
| {"items":[{"name":"«$.tags[0].name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.tags[0][\"name\"] FROM c"}
< raw
| {"items":[{"name":"«$.tags[0].name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.tags[0].skills FROM c"}
< raw
| {"items":[{"skills":["«$.tags[0].skills[0]»","«$.tags[0].skills[1]»"]}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT (c.name) FROM c"}
< raw
| {"items":[{"name":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT * FROM c"}
< raw
| {"items":[{"id":"«$.id»","pk":"«$.pk»","name":"«$.name»","nick":"«$.nick»","n":1101,"address":{"city":"«$.address.city»","zip":"«$.address.zip»","geo":{"lat":"«$.address.geo.lat»"}},"tags":[{"name":"«$.tags[0].name»","skills":["«$.tags[0].skills[0]»","«$.tags[0].skills[1]»"]},{"name":"«$.tags[1].name»","skills":["«$.tags[1].skills[0]»"]}],"phones":["«$.phones[0]»","«$.phones[1]»"],"_rid":"(system)","_self":"(system)","_etag":"(system)","_attachments":"(system)","_ts":"(system)"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c FROM c"}
< raw
| {"items":[{"c":{"id":"«$.id»","pk":"«$.pk»","name":"«$.name»","nick":"«$.nick»","n":1101,"address":{"city":"«$.address.city»","zip":"«$.address.zip»","geo":{"lat":"«$.address.geo.lat»"}},"tags":[{"name":"«$.tags[0].name»","skills":["«$.tags[0].skills[0]»","«$.tags[0].skills[1]»"]},{"name":"«$.tags[1].name»","skills":["«$.tags[1].skills[0]»"]}],"phones":["«$.phones[0]»","«$.phones[1]»"],"_rid":"(system)","_self":"(system)","_etag":"(system)","_attachments":"(system)","_ts":"(system)"}}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT VALUE c FROM c"}
< raw
| {"selectValue":true,"items":[{"id":"«$.id»","pk":"«$.pk»","name":"«$.name»","nick":"«$.nick»","n":1101,"address":{"city":"«$.address.city»","zip":"«$.address.zip»","geo":{"lat":"«$.address.geo.lat»"}},"tags":[{"name":"«$.tags[0].name»","skills":["«$.tags[0].skills[0]»","«$.tags[0].skills[1]»"]},{"name":"«$.tags[1].name»","skills":["«$.tags[1].skills[0]»"]}],"phones":["«$.phones[0]»","«$.phones[1]»"],"_rid":"(system)","_self":"(system)","_etag":"(system)","_attachments":"(system)","_ts":"(system)"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT root FROM root"}
< raw
| {"items":[{"root":{"id":"«$.id»","pk":"«$.pk»","name":"«$.name»","nick":"«$.nick»","n":1101,"address":{"city":"«$.address.city»","zip":"«$.address.zip»","geo":{"lat":"«$.address.geo.lat»"}},"tags":[{"name":"«$.tags[0].name»","skills":["«$.tags[0].skills[0]»","«$.tags[0].skills[1]»"]},{"name":"«$.tags[1].name»","skills":["«$.tags[1].skills[0]»"]}],"phones":["«$.phones[0]»","«$.phones[1]»"],"_rid":"(system)","_self":"(system)","_etag":"(system)","_attachments":"(system)","_ts":"(system)"}}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT TOP 1 * FROM c"}
< raw
| {"items":[{"id":"«$.id»","pk":"«$.pk»","name":"«$.name»","nick":"«$.nick»","n":1101,"address":{"city":"«$.address.city»","zip":"«$.address.zip»","geo":{"lat":"«$.address.geo.lat»"}},"tags":[{"name":"«$.tags[0].name»","skills":["«$.tags[0].skills[0]»","«$.tags[0].skills[1]»"]},{"name":"«$.tags[1].name»","skills":["«$.tags[1].skills[0]»"]}],"phones":["«$.phones[0]»","«$.phones[1]»"],"_rid":"(system)","_self":"(system)","_etag":"(system)","_attachments":"(system)","_ts":"(system)"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.name AS n FROM c"}
< raw
| {"items":[{"n":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.name n FROM c"}
< raw
| {"items":[{"n":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.address.city AS name FROM c"}
< raw
| {"items":[{"name":"«$.address.city»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.name AS address, c.address AS name FROM c"}
< raw
| {"items":[{"address":"«$.name»","name":{"city":"«$.address.city»","zip":"«$.address.zip»","geo":{"lat":"«$.address.geo.lat»"}}}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.address AS a FROM c"}
< raw
| {"items":[{"a":{"city":"«$.address.city»","zip":"«$.address.zip»","geo":{"lat":"«$.address.geo.lat»"}}}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT r.name FROM root r"}
< raw
| {"items":[{"name":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT root.name FROM root"}
< raw
| {"items":[{"name":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT x.name FROM mycontainer x"}
< raw
| {"items":[{"name":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c FROM container c"}
< raw
| {"items":[{"c":{"id":"«$.id»","pk":"«$.pk»","name":"«$.name»","nick":"«$.nick»","n":1101,"address":{"city":"«$.address.city»","zip":"«$.address.zip»","geo":{"lat":"«$.address.geo.lat»"}},"tags":[{"name":"«$.tags[0].name»","skills":["«$.tags[0].skills[0]»","«$.tags[0].skills[1]»"]},{"name":"«$.tags[1].name»","skills":["«$.tags[1].skills[0]»"]}],"phones":["«$.phones[0]»","«$.phones[1]»"],"_rid":"(system)","_self":"(system)","_etag":"(system)","_attachments":"(system)","_ts":"(system)"}}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.name FROM container c"}
< raw
| {"items":[{"name":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT UPPER(c.name) FROM c"}
< raw
| {"items":[{"$1":"«$.NAME»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.name, UPPER(c.name) FROM c"}
< raw
| {"items":[{"name":"«$.name»","$1":"«$.NAME»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT UPPER(c.name), c.name, LOWER(c.nick) FROM c"}
< raw
| {"items":[{"$1":"«$.NAME»","name":"«$.name»","$2":"«$.nick»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.name, c.nick, c.n + 1, UPPER(c.name) AS u, LOWER(c.name) FROM c"}
< raw
| {"items":[{"name":"«$.name»","nick":"«$.nick»","$1":1102,"u":"«$.NAME»","$2":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.name AS n, c.tags[0], [c.name] FROM c"}
< raw
| {"items":[{"n":"«$.name»","$1":{"name":"«$.tags[0].name»","skills":["«$.tags[0].skills[0]»","«$.tags[0].skills[1]»"]},"$2":["«$.name»"]}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT -c.n FROM c"}
< raw
| {"items":[{"$1":-1101}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT {\"a\": c.name} FROM c"}
< raw
| {"items":[{"$1":{"a":"«$.name»"}}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT IS_DEFINED(c.nick) FROM c"}
< raw
| {"items":[{"$1":true}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT {\"k\": c.name}.k FROM c"}
< raw
| {"items":[{"k":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT CONCAT(c.name, \"-\", c.address.city) AS label FROM c"}
< raw
| {"items":[{"label":"«$.name»-«$.address.city»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.name ?? c.nick AS who FROM c"}
< raw
| {"items":[{"who":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT IIF(c.n \u003e 0, c.name, c.address.city) AS pick FROM c"}
< raw
| {"items":[{"pick":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.n \u003e 0 ? c.name : c.address.zip AS pick FROM c"}
< raw
| {"items":[{"pick":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT {\"a\": c.name, \"b\": c.address.city} AS o FROM c"}
< raw
| {"items":[{"o":{"a":"«$.name»","b":"«$.address.city»"}}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT [c.name, c.address.zip] AS arr FROM c"}
< raw
| {"items":[{"arr":["«$.name»","«$.address.zip»"]}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT [c.name, c.address.city][0] AS first FROM c"}
< raw
| {"items":[{"first":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT {\"k\": c.name}.k AS kk FROM c"}
< raw
| {"items":[{"kk":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT (c.address ?? c.nick).city AS city FROM c"}
< raw
| {"items":[{"city":"«$.address.city»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT (c.n \u003e 0 ? c.address : c.tags).zip AS zip FROM c"}
< raw
| {"items":[{"zip":"«$.address.zip»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT IIF(c.n \u003e 0, c.address, c.tags).geo.lat FROM c"}
< raw
| {"items":[{"lat":"«$.address.geo.lat»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT LOWER(c.name) AS l, UPPER(c.nick) AS u FROM c"}
< raw
| {"items":[{"l":"«$.name»","u":"«$.NICK»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT REPLACE(c.name, \"x\", \"y\") AS r FROM c"}
< raw
| {"items":[{"r":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT ToString(c.n) AS s, c.n * 2 AS d FROM c"}
< raw
| {"items":[{"s":"1101","d":2202}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT ARRAY_SLICE(c.phones, 0, 1) AS ph FROM c"}
< raw
| {"items":[{"ph":["«$.phones[0]»"]}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT ARRAY_CONCAT(c.phones, [c.name]) AS all FROM c"}
< raw
| {"items":[{"all":["«$.phones[0]»","«$.phones[1]»","«$.name»"]}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT ARRAY_LENGTH(c.tags) AS n FROM c"}
< raw
| {"items":[{"n":2}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT 1 AS one, \"k\" AS k FROM c"}
< raw
| {"items":[{"one":1,"k":"k"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT VALUE c.name FROM c"}
< raw
| {"selectValue":true,"items":["«$.name»"]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT VALUE c.address FROM c"}
< raw
| {"selectValue":true,"items":[{"city":"«$.address.city»","zip":"«$.address.zip»","geo":{"lat":"«$.address.geo.lat»"}}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT VALUE {\"n\": c.name, \"z\": c.address.zip} FROM c"}
< raw
| {"selectValue":true,"items":[{"n":"«$.name»","z":"«$.address.zip»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT VALUE UPPER(c.name) FROM c"}
< raw
| {"selectValue":true,"items":["«$.NAME»"]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT VALUE [c.name, c.n] FROM c"}
< raw
| {"selectValue":true,"items":[["«$.name»",1101]]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT DISTINCT VALUE c.address.city FROM c"}
< raw
| {"selectValue":true,"items":["«$.address.city»"]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT t.name FROM c JOIN t IN c.tags"}
< raw
| {"items":[{"name":"«$.tags[0].name»"},{"name":"«$.tags[1].name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT t FROM c JOIN t IN c.tags"}
< raw
| {"items":[{"t":{"name":"«$.tags[0].name»","skills":["«$.tags[0].skills[0]»","«$.tags[0].skills[1]»"]}},{"t":{"name":"«$.tags[1].name»","skills":["«$.tags[1].skills[0]»"]}}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.name, t.name AS tag FROM c JOIN t IN c.tags"}
< raw
| {"items":[{"name":"«$.name»","tag":"«$.tags[0].name»"},{"name":"«$.name»","tag":"«$.tags[1].name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT VALUE t FROM c JOIN t IN c.tags"}
< raw
| {"selectValue":true,"items":[{"name":"«$.tags[0].name»","skills":["«$.tags[0].skills[0]»","«$.tags[0].skills[1]»"]},{"name":"«$.tags[1].name»","skills":["«$.tags[1].skills[0]»"]}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT t.name AS tag, s AS skill FROM c JOIN t IN c.tags JOIN s IN t.skills"}
< raw
| {"items":[{"tag":"«$.tags[0].name»","skill":"«$.tags[0].skills[0]»"},{"tag":"«$.tags[0].name»","skill":"«$.tags[0].skills[1]»"},{"tag":"«$.tags[1].name»","skill":"«$.tags[1].skills[0]»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT t.skills[0] FROM c JOIN t IN c.tags"}
< raw
| {"items":[{"$1":"«$.tags[0].skills[0]»"},{"$1":"«$.tags[1].skills[0]»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT p FROM c JOIN p IN c.phones"}
< raw
| {"items":[{"p":"«$.phones[0]»"},{"p":"«$.phones[1]»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.name, t FROM c JOIN t IN c.phones"}
< raw
| {"items":[{"name":"«$.name»","t":"«$.phones[0]»"},{"name":"«$.name»","t":"«$.phones[1]»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT VALUE t.name FROM t IN c.tags"}
< raw
| {"selectValue":true,"items":["«$.tags[0].name»","«$.tags[1].name»"]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT t.name FROM t IN c.tags WHERE t.name != \"\""}
< raw
| {"items":[{"name":"«$.tags[0].name»"},{"name":"«$.tags[1].name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT * FROM t IN c.tags"}
< raw
| {"items":[{"name":"«$.tags[0].name»","skills":["«$.tags[0].skills[0]»","«$.tags[0].skills[1]»"]},{"name":"«$.tags[1].name»","skills":["«$.tags[1].skills[0]»"]}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT a.city FROM c.address a"}
< raw
| {"items":[{"city":"«$.address.city»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT a FROM c.address a"}
< raw
| {"items":[{"a":{"city":"«$.address.city»","zip":"«$.address.zip»","geo":{"lat":"«$.address.geo.lat»"}}}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT * FROM c.address"}
< raw
| {"items":[{"city":"«$.address.city»","zip":"«$.address.zip»","geo":{"lat":"«$.address.geo.lat»"}}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT address.city FROM c.address"}
< raw
| {"items":[{"city":"«$.address.city»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT t.name FROM c JOIN t IN c.tags WHERE ARRAY_CONTAINS(t.skills, \"«$.tags[0].skills[0]»\")"}
< raw
| {"items":[{"name":"«$.tags[0].name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT ARRAY(SELECT VALUE t.name FROM t IN c.tags) AS names FROM c"}
< raw
| {"items":[{"names":["«$.tags[0].name»","«$.tags[1].name»"]}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT ARRAY(SELECT t.name FROM t IN c.tags) AS names FROM c"}
< raw
| {"items":[{"names":[{"name":"«$.tags[0].name»"},{"name":"«$.tags[1].name»"}]}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT ARRAY(SELECT VALUE p FROM p IN c.phones) FROM c"}
< raw
| {"items":[{"$1":["«$.phones[0]»","«$.phones[1]»"]}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT ARRAY(SELECT VALUE s FROM s IN t.skills) AS sk FROM c JOIN t IN c.tags"}
< raw
| {"items":[{"sk":["«$.tags[0].skills[0]»","«$.tags[0].skills[1]»"]},{"sk":["«$.tags[1].skills[0]»"]}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT VALUE ARRAY(SELECT VALUE t.name FROM t IN c.tags) FROM c"}
< raw
| {"selectValue":true,"items":[["«$.tags[0].name»","«$.tags[1].name»"]]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT (SELECT VALUE COUNT(1) FROM t IN c.tags) AS cnt FROM c"}
< raw
| {"items":[{"cnt":2}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT (SELECT VALUE c.name) AS n FROM c"}
< raw
| {"items":[{"n":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT FIRST(SELECT VALUE t.name FROM t IN c.tags) AS f FROM c"}
< raw
| {"items":[{"f":"«$.tags[0].name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT EXISTS(SELECT VALUE t FROM t IN c.tags) AS has FROM c"}
< raw
| {"items":[{"has":true}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.name FROM c WHERE EXISTS(SELECT VALUE t FROM t IN c.tags WHERE t.name = \"«$.tags[1].name»\")"}
< raw
| {"items":[{"name":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT t.name FROM c JOIN (SELECT VALUE x FROM x IN c.tags) t"}
< raw
| {"items":[{"name":"«$.tags[0].name»"},{"name":"«$.tags[1].name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT s.name FROM (SELECT c.name FROM c) s"}
< raw
| {"items":[{"name":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT s.city FROM (SELECT c.name, c.address.city FROM c) s"}
< raw
| {"items":[{"city":"«$.address.city»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT s FROM (SELECT c.name, c.address.city FROM c) s"}
< raw
| {"items":[{"s":{"name":"«$.name»","city":"«$.address.city»"}}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT * FROM (SELECT c.name, c.address.city FROM c) s"}
< raw
| {"items":[{"name":"«$.name»","city":"«$.address.city»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT * FROM (SELECT VALUE c.address FROM c) s"}
< raw
| {"items":[{"city":"«$.address.city»","zip":"«$.address.zip»","geo":{"lat":"«$.address.geo.lat»"}}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT s.zip FROM (SELECT VALUE c.address FROM c) s"}
< raw
| {"items":[{"zip":"«$.address.zip»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT VALUE s FROM (SELECT VALUE c.address FROM c) s"}
< raw
| {"selectValue":true,"items":[{"city":"«$.address.city»","zip":"«$.address.zip»","geo":{"lat":"«$.address.geo.lat»"}}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT s.n FROM (SELECT c.name AS n FROM c) s"}
< raw
| {"items":[{"n":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT VALUE s.name FROM (SELECT c.name FROM c) s"}
< raw
| {"selectValue":true,"items":["«$.name»"]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT COUNT(1) AS n FROM c"}
< raw
| {"items":[{"n":1}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT COUNT(1) FROM c"}
< raw
| {"items":[{"$1":1}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT VALUE COUNT(1) FROM c"}
< raw
| {"selectValue":true,"items":[1]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT MAX(c.name) AS m FROM c"}
< raw
| {"items":[{"m":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT MIN(c.n) FROM c"}
< raw
| {"items":[{"$1":1101}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT SUM(c.n) AS total, AVG(c.n) AS mean FROM c"}
< raw
| {"items":[{"total":1101,"mean":1101}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.address.city, COUNT(1) AS n FROM c GROUP BY c.address.city"}
< raw
| {"items":[{"city":"«$.address.city»","n":1}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.address.city, COUNT(1) FROM c GROUP BY c.address.city"}
< raw
| {"items":[{"city":"«$.address.city»","$1":1}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT t.name, COUNT(1) AS n FROM c JOIN t IN c.tags GROUP BY t.name"}
< raw
| {"items":[{"name":"«$.tags[0].name»","n":1},{"name":"«$.tags[1].name»","n":1}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT DISTINCT c.address.city FROM c"}
< raw
| {"items":[{"city":"«$.address.city»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT DISTINCT VALUE t.name FROM t IN c.tags"}
< raw
| {"selectValue":true,"items":["«$.tags[0].name»","«$.tags[1].name»"]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT TOP 1 c.name FROM c ORDER BY c.name"}
< raw
| {"items":[{"name":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.name FROM c ORDER BY c.address.city DESC"}
< raw
| {"items":[{"name":"«$.name»"}]}
> raw
| {"op":"query","fixture":"document","statement":"SELECT c.address.city FROM c WHERE c.n \u003e 0 OFFSET 0 LIMIT 1"}
< raw
| {"items":[{"city":"«$.address.city»"}]}
