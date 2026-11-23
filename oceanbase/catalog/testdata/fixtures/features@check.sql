CREATE TABLE bookings (id INT PRIMARY KEY, starts DATE, ends DATE, guests INT CHECK (guests > 0), nights INT, note VARCHAR(100), CONSTRAINT chk_dates CHECK (starts <= ends), CONSTRAINT chk_nights CHECK (nights > 0));
CREATE TABLE rooms (id INT PRIMARY KEY, beds INT, CONSTRAINT chk_always CHECK (1 = 1))
