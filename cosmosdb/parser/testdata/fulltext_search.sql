SELECT TOP 10 c.name, c.description FROM c WHERE FullTextContains(c.description, 'bicycle') ORDER BY RANK FullTextScore(c.description, ['bicycle'])
