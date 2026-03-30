SELECT c.name, IS_STRING(c.name) AS isStr, IS_NUMBER(c.population) AS isNum, IS_DEFINED(c.extra) AS isDef FROM c
