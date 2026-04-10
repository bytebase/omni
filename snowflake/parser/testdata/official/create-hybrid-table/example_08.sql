CREATE OR REPLACE HYBRID TABLE team
  (team_id INT PRIMARY KEY,
  team_name VARCHAR(40),
  stadium VARCHAR(40));

CREATE OR REPLACE HYBRID TABLE player
  (player_id INT PRIMARY KEY,
  first_name VARCHAR(40),
  last_name VARCHAR(40),
  team_id INT,
  FOREIGN KEY (team_id) REFERENCES team(team_id));
