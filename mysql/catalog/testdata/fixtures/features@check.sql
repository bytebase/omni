CREATE TABLE bookings (id INT PRIMARY KEY, starts DATE, ends DATE, guests INT CHECK (guests > 0), note VARCHAR(100), CONSTRAINT chk_dates CHECK (starts <= ends))
