MERGE INTO target_table USING (
    SELECT id, description FROM source_table
) AS filtered_source_table
ON target_table.id = filtered_source_table.id
WHEN MATCHED THEN UPDATE SET a = filtered_source_table.b;

