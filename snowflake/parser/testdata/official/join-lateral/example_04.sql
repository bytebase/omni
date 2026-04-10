CREATE OR REPLACE TABLE persons AS
  SELECT column1 AS id, PARSE_JSON(column2) AS c
    FROM VALUES
      (12712555,
       '{ "name": { "first": "John", "last": "Smith" },
          "contact": [{ "business": [
            { "type": "phone", "content": "555-1234" },
            { "type": "email", "content": "j.smith@example.com" }
          ]}]}'),
      (98127771,
       '{ "name": { "first": "Jane", "last": "Doe" },
          "contact": [{ "business": [
            { "type": "phone", "content": "555-1236" },
            { "type": "email", "content": "j.doe@example.com" }
          ]}]}');
