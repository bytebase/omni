SELECT c.name, ARRAY(SELECT VALUE t FROM t IN c.tags WHERE t.active = true) AS activeTags FROM c
