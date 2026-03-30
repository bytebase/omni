SELECT p.name, t.tag FROM products p JOIN t IN p.tags WHERE t.active = true
