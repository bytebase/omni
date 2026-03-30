SELECT CONCAT(c.name, ' - ', c.country) AS label, StringToNumber(c.population) / 1000 AS populationInThousands FROM c
