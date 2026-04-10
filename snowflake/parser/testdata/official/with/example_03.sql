WITH RECURSIVE current_layer (indent, layer_ID, parent_component_ID, component_id, description, sort_key) AS (
  SELECT 
      '...', 
      1, 
      parent_component_ID, 
      component_id, 
      description, 
      '0001'
    FROM components WHERE component_id = 1
  UNION ALL
  SELECT indent || '...',
      layer_ID + 1,
      components.parent_component_ID,
      components.component_id, 
      components.description,
      sort_key || SUBSTRING('000' || components.component_ID, -4)
    FROM current_layer JOIN components 
      ON (components.parent_component_id = current_layer.component_id)
  )
SELECT
  indent || description AS description, 
  component_id,
  parent_component_ID
FROM current_layer
ORDER BY sort_key;
