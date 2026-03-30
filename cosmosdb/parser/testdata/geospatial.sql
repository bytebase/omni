SELECT c.name, ST_DISTANCE(c.location, {"type": "Point", "coordinates": [55.2708, 25.2048]}) AS distFromDubaiInMeters FROM c
