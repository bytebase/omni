CREATE TYPE mood AS ENUM ('happy', 'sad');
CREATE TABLE notes (t text, v varchar(10), n int, m mood);
