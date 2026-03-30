SELECT c.name, c.value = Infinity AS isInfinite, IS_NUMBER(NaN) AS nanCheck FROM c
