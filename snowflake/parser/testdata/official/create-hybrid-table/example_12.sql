INSERT INTO team VALUES (0, 'Unknown', 'Unknown');
INSERT INTO player VALUES (200, 'Tommy', 'Atkins', 0);

SELECT * FROM team t, player p WHERE t.team_id=p.team_id;
