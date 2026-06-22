-- Vector similarity search using ANN (Approximate Nearest Neighbor)

-- Basic vector search with ORDER BY ANN
SELECT * FROM products 
WHERE category = 'electronics' 
ORDER BY embedding ANN OF [0.1, 0.2, 0.3, 0.4, 0.5] 
LIMIT 10;

-- Vector search without WHERE clause
SELECT id, name, price FROM products 
ORDER BY embedding ANN OF [0.25, 0.5, 0.75] 
LIMIT 5;

-- Combined with other ORDER BY (if supported)
SELECT * FROM products 
WHERE price < 1000 
ORDER BY embedding ANN OF [0.1, 0.2, 0.3] LIMIT 20;