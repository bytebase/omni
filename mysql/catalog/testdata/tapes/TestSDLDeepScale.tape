-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestSDLDeepScale default#1
> begin
< ok
> exec
| CREATE TABLE `t000` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t000_code` (code),
| 	KEY `idx_t000_name` (name),
| 	KEY `idx_t000_amount` (amount)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t001` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t001_code` (code),
| 	KEY `idx_t001_name` (name),
| 	KEY `idx_t001_prev` (prev_id),
| 	CONSTRAINT `fk_t001_prev` FOREIGN KEY (prev_id) REFERENCES `t000` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t002` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t002_code` (code),
| 	KEY `idx_t002_name` (name),
| 	KEY `idx_t002_prev` (prev_id),
| 	CONSTRAINT `fk_t002_prev` FOREIGN KEY (prev_id) REFERENCES `t001` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t003` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t003_code` (code),
| 	KEY `idx_t003_name` (name),
| 	KEY `idx_t003_prev` (prev_id),
| 	CONSTRAINT `fk_t003_prev` FOREIGN KEY (prev_id) REFERENCES `t002` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t004` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t004_code` (code),
| 	KEY `idx_t004_name` (name),
| 	KEY `idx_t004_prev` (prev_id),
| 	CONSTRAINT `fk_t004_prev` FOREIGN KEY (prev_id) REFERENCES `t003` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t005` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t005_code` (code),
| 	KEY `idx_t005_name` (name),
| 	KEY `idx_t005_prev` (prev_id),
| 	CONSTRAINT `fk_t005_prev` FOREIGN KEY (prev_id) REFERENCES `t004` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t006` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t006_code` (code),
| 	KEY `idx_t006_name` (name),
| 	KEY `idx_t006_prev` (prev_id),
| 	CONSTRAINT `fk_t006_prev` FOREIGN KEY (prev_id) REFERENCES `t005` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t007` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t007_code` (code),
| 	KEY `idx_t007_name` (name),
| 	KEY `idx_t007_prev` (prev_id),
| 	CONSTRAINT `fk_t007_prev` FOREIGN KEY (prev_id) REFERENCES `t006` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t008` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t008_code` (code),
| 	KEY `idx_t008_name` (name),
| 	KEY `idx_t008_prev` (prev_id),
| 	CONSTRAINT `fk_t008_prev` FOREIGN KEY (prev_id) REFERENCES `t007` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t009` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t009_code` (code),
| 	KEY `idx_t009_name` (name),
| 	KEY `idx_t009_prev` (prev_id),
| 	CONSTRAINT `fk_t009_prev` FOREIGN KEY (prev_id) REFERENCES `t008` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t010` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t010_code` (code),
| 	KEY `idx_t010_name` (name),
| 	KEY `idx_t010_prev` (prev_id),
| 	CONSTRAINT `fk_t010_prev` FOREIGN KEY (prev_id) REFERENCES `t009` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t011` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t011_code` (code),
| 	KEY `idx_t011_name` (name),
| 	KEY `idx_t011_prev` (prev_id),
| 	CONSTRAINT `fk_t011_prev` FOREIGN KEY (prev_id) REFERENCES `t010` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t012` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t012_code` (code),
| 	KEY `idx_t012_name` (name),
| 	KEY `idx_t012_prev` (prev_id),
| 	CONSTRAINT `fk_t012_prev` FOREIGN KEY (prev_id) REFERENCES `t011` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t013` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t013_code` (code),
| 	KEY `idx_t013_name` (name),
| 	KEY `idx_t013_prev` (prev_id),
| 	CONSTRAINT `fk_t013_prev` FOREIGN KEY (prev_id) REFERENCES `t012` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t014` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t014_code` (code),
| 	KEY `idx_t014_name` (name),
| 	KEY `idx_t014_prev` (prev_id),
| 	CONSTRAINT `fk_t014_prev` FOREIGN KEY (prev_id) REFERENCES `t013` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t015` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t015_code` (code),
| 	KEY `idx_t015_name` (name),
| 	KEY `idx_t015_prev` (prev_id),
| 	CONSTRAINT `fk_t015_prev` FOREIGN KEY (prev_id) REFERENCES `t014` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t016` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t016_code` (code),
| 	KEY `idx_t016_name` (name),
| 	KEY `idx_t016_prev` (prev_id),
| 	CONSTRAINT `fk_t016_prev` FOREIGN KEY (prev_id) REFERENCES `t015` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t017` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t017_code` (code),
| 	KEY `idx_t017_name` (name),
| 	KEY `idx_t017_prev` (prev_id),
| 	CONSTRAINT `fk_t017_prev` FOREIGN KEY (prev_id) REFERENCES `t016` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t018` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t018_code` (code),
| 	KEY `idx_t018_name` (name),
| 	KEY `idx_t018_prev` (prev_id),
| 	CONSTRAINT `fk_t018_prev` FOREIGN KEY (prev_id) REFERENCES `t017` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t019` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t019_code` (code),
| 	KEY `idx_t019_name` (name),
| 	KEY `idx_t019_prev` (prev_id),
| 	CONSTRAINT `fk_t019_prev` FOREIGN KEY (prev_id) REFERENCES `t018` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t020` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t020_code` (code),
| 	KEY `idx_t020_name` (name),
| 	KEY `idx_t020_prev` (prev_id),
| 	CONSTRAINT `fk_t020_prev` FOREIGN KEY (prev_id) REFERENCES `t019` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t021` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t021_code` (code),
| 	KEY `idx_t021_name` (name),
| 	KEY `idx_t021_prev` (prev_id),
| 	CONSTRAINT `fk_t021_prev` FOREIGN KEY (prev_id) REFERENCES `t020` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t022` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t022_code` (code),
| 	KEY `idx_t022_name` (name),
| 	KEY `idx_t022_prev` (prev_id),
| 	CONSTRAINT `fk_t022_prev` FOREIGN KEY (prev_id) REFERENCES `t021` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t023` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t023_code` (code),
| 	KEY `idx_t023_name` (name),
| 	KEY `idx_t023_prev` (prev_id),
| 	CONSTRAINT `fk_t023_prev` FOREIGN KEY (prev_id) REFERENCES `t022` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t024` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t024_code` (code),
| 	KEY `idx_t024_name` (name),
| 	KEY `idx_t024_prev` (prev_id),
| 	CONSTRAINT `fk_t024_prev` FOREIGN KEY (prev_id) REFERENCES `t023` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t025` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t025_code` (code),
| 	KEY `idx_t025_name` (name),
| 	KEY `idx_t025_prev` (prev_id),
| 	CONSTRAINT `fk_t025_prev` FOREIGN KEY (prev_id) REFERENCES `t024` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t026` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t026_code` (code),
| 	KEY `idx_t026_name` (name),
| 	KEY `idx_t026_prev` (prev_id),
| 	CONSTRAINT `fk_t026_prev` FOREIGN KEY (prev_id) REFERENCES `t025` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t027` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t027_code` (code),
| 	KEY `idx_t027_name` (name),
| 	KEY `idx_t027_prev` (prev_id),
| 	CONSTRAINT `fk_t027_prev` FOREIGN KEY (prev_id) REFERENCES `t026` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t028` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t028_code` (code),
| 	KEY `idx_t028_name` (name),
| 	KEY `idx_t028_prev` (prev_id),
| 	CONSTRAINT `fk_t028_prev` FOREIGN KEY (prev_id) REFERENCES `t027` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t029` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t029_code` (code),
| 	KEY `idx_t029_name` (name),
| 	KEY `idx_t029_prev` (prev_id),
| 	CONSTRAINT `fk_t029_prev` FOREIGN KEY (prev_id) REFERENCES `t028` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t030` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t030_code` (code),
| 	KEY `idx_t030_name` (name),
| 	KEY `idx_t030_prev` (prev_id),
| 	CONSTRAINT `fk_t030_prev` FOREIGN KEY (prev_id) REFERENCES `t029` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t031` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t031_code` (code),
| 	KEY `idx_t031_name` (name),
| 	KEY `idx_t031_prev` (prev_id),
| 	CONSTRAINT `fk_t031_prev` FOREIGN KEY (prev_id) REFERENCES `t030` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t032` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t032_code` (code),
| 	KEY `idx_t032_name` (name),
| 	KEY `idx_t032_prev` (prev_id),
| 	CONSTRAINT `fk_t032_prev` FOREIGN KEY (prev_id) REFERENCES `t031` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t033` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t033_code` (code),
| 	KEY `idx_t033_name` (name),
| 	KEY `idx_t033_prev` (prev_id),
| 	CONSTRAINT `fk_t033_prev` FOREIGN KEY (prev_id) REFERENCES `t032` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t034` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t034_code` (code),
| 	KEY `idx_t034_name` (name),
| 	KEY `idx_t034_prev` (prev_id),
| 	CONSTRAINT `fk_t034_prev` FOREIGN KEY (prev_id) REFERENCES `t033` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t035` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t035_code` (code),
| 	KEY `idx_t035_name` (name),
| 	KEY `idx_t035_prev` (prev_id),
| 	CONSTRAINT `fk_t035_prev` FOREIGN KEY (prev_id) REFERENCES `t034` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t036` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t036_code` (code),
| 	KEY `idx_t036_name` (name),
| 	KEY `idx_t036_prev` (prev_id),
| 	CONSTRAINT `fk_t036_prev` FOREIGN KEY (prev_id) REFERENCES `t035` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t037` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t037_code` (code),
| 	KEY `idx_t037_name` (name),
| 	KEY `idx_t037_prev` (prev_id),
| 	CONSTRAINT `fk_t037_prev` FOREIGN KEY (prev_id) REFERENCES `t036` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t038` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t038_code` (code),
| 	KEY `idx_t038_name` (name),
| 	KEY `idx_t038_prev` (prev_id),
| 	CONSTRAINT `fk_t038_prev` FOREIGN KEY (prev_id) REFERENCES `t037` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t039` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t039_code` (code),
| 	KEY `idx_t039_name` (name),
| 	KEY `idx_t039_prev` (prev_id),
| 	CONSTRAINT `fk_t039_prev` FOREIGN KEY (prev_id) REFERENCES `t038` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t040` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t040_code` (code),
| 	KEY `idx_t040_name` (name),
| 	KEY `idx_t040_prev` (prev_id),
| 	CONSTRAINT `fk_t040_prev` FOREIGN KEY (prev_id) REFERENCES `t039` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t041` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t041_code` (code),
| 	KEY `idx_t041_name` (name),
| 	KEY `idx_t041_prev` (prev_id),
| 	CONSTRAINT `fk_t041_prev` FOREIGN KEY (prev_id) REFERENCES `t040` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t042` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t042_code` (code),
| 	KEY `idx_t042_name` (name),
| 	KEY `idx_t042_prev` (prev_id),
| 	CONSTRAINT `fk_t042_prev` FOREIGN KEY (prev_id) REFERENCES `t041` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t043` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t043_code` (code),
| 	KEY `idx_t043_name` (name),
| 	KEY `idx_t043_prev` (prev_id),
| 	CONSTRAINT `fk_t043_prev` FOREIGN KEY (prev_id) REFERENCES `t042` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t044` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t044_code` (code),
| 	KEY `idx_t044_name` (name),
| 	KEY `idx_t044_prev` (prev_id),
| 	CONSTRAINT `fk_t044_prev` FOREIGN KEY (prev_id) REFERENCES `t043` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t045` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t045_code` (code),
| 	KEY `idx_t045_name` (name),
| 	KEY `idx_t045_prev` (prev_id),
| 	CONSTRAINT `fk_t045_prev` FOREIGN KEY (prev_id) REFERENCES `t044` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t046` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t046_code` (code),
| 	KEY `idx_t046_name` (name),
| 	KEY `idx_t046_prev` (prev_id),
| 	CONSTRAINT `fk_t046_prev` FOREIGN KEY (prev_id) REFERENCES `t045` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t047` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t047_code` (code),
| 	KEY `idx_t047_name` (name),
| 	KEY `idx_t047_prev` (prev_id),
| 	CONSTRAINT `fk_t047_prev` FOREIGN KEY (prev_id) REFERENCES `t046` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t048` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t048_code` (code),
| 	KEY `idx_t048_name` (name),
| 	KEY `idx_t048_prev` (prev_id),
| 	CONSTRAINT `fk_t048_prev` FOREIGN KEY (prev_id) REFERENCES `t047` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t049` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t049_code` (code),
| 	KEY `idx_t049_name` (name),
| 	KEY `idx_t049_prev` (prev_id),
| 	CONSTRAINT `fk_t049_prev` FOREIGN KEY (prev_id) REFERENCES `t048` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t050` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t050_code` (code),
| 	KEY `idx_t050_name` (name),
| 	KEY `idx_t050_prev` (prev_id),
| 	CONSTRAINT `fk_t050_prev` FOREIGN KEY (prev_id) REFERENCES `t049` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t051` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t051_code` (code),
| 	KEY `idx_t051_name` (name),
| 	KEY `idx_t051_prev` (prev_id),
| 	CONSTRAINT `fk_t051_prev` FOREIGN KEY (prev_id) REFERENCES `t050` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t052` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t052_code` (code),
| 	KEY `idx_t052_name` (name),
| 	KEY `idx_t052_prev` (prev_id),
| 	CONSTRAINT `fk_t052_prev` FOREIGN KEY (prev_id) REFERENCES `t051` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t053` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t053_code` (code),
| 	KEY `idx_t053_name` (name),
| 	KEY `idx_t053_prev` (prev_id),
| 	CONSTRAINT `fk_t053_prev` FOREIGN KEY (prev_id) REFERENCES `t052` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t054` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t054_code` (code),
| 	KEY `idx_t054_name` (name),
| 	KEY `idx_t054_prev` (prev_id),
| 	CONSTRAINT `fk_t054_prev` FOREIGN KEY (prev_id) REFERENCES `t053` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t055` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t055_code` (code),
| 	KEY `idx_t055_name` (name),
| 	KEY `idx_t055_prev` (prev_id),
| 	CONSTRAINT `fk_t055_prev` FOREIGN KEY (prev_id) REFERENCES `t054` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t056` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t056_code` (code),
| 	KEY `idx_t056_name` (name),
| 	KEY `idx_t056_prev` (prev_id),
| 	CONSTRAINT `fk_t056_prev` FOREIGN KEY (prev_id) REFERENCES `t055` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t057` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t057_code` (code),
| 	KEY `idx_t057_name` (name),
| 	KEY `idx_t057_prev` (prev_id),
| 	CONSTRAINT `fk_t057_prev` FOREIGN KEY (prev_id) REFERENCES `t056` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t058` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t058_code` (code),
| 	KEY `idx_t058_name` (name),
| 	KEY `idx_t058_prev` (prev_id),
| 	CONSTRAINT `fk_t058_prev` FOREIGN KEY (prev_id) REFERENCES `t057` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t059` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t059_code` (code),
| 	KEY `idx_t059_name` (name),
| 	KEY `idx_t059_prev` (prev_id),
| 	CONSTRAINT `fk_t059_prev` FOREIGN KEY (prev_id) REFERENCES `t058` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t060` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t060_code` (code),
| 	KEY `idx_t060_name` (name),
| 	KEY `idx_t060_prev` (prev_id),
| 	CONSTRAINT `fk_t060_prev` FOREIGN KEY (prev_id) REFERENCES `t059` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t061` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t061_code` (code),
| 	KEY `idx_t061_name` (name),
| 	KEY `idx_t061_prev` (prev_id),
| 	CONSTRAINT `fk_t061_prev` FOREIGN KEY (prev_id) REFERENCES `t060` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t062` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t062_code` (code),
| 	KEY `idx_t062_name` (name),
| 	KEY `idx_t062_prev` (prev_id),
| 	CONSTRAINT `fk_t062_prev` FOREIGN KEY (prev_id) REFERENCES `t061` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t063` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t063_code` (code),
| 	KEY `idx_t063_name` (name),
| 	KEY `idx_t063_prev` (prev_id),
| 	CONSTRAINT `fk_t063_prev` FOREIGN KEY (prev_id) REFERENCES `t062` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t064` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t064_code` (code),
| 	KEY `idx_t064_name` (name),
| 	KEY `idx_t064_prev` (prev_id),
| 	CONSTRAINT `fk_t064_prev` FOREIGN KEY (prev_id) REFERENCES `t063` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t065` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t065_code` (code),
| 	KEY `idx_t065_name` (name),
| 	KEY `idx_t065_prev` (prev_id),
| 	CONSTRAINT `fk_t065_prev` FOREIGN KEY (prev_id) REFERENCES `t064` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t066` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t066_code` (code),
| 	KEY `idx_t066_name` (name),
| 	KEY `idx_t066_prev` (prev_id),
| 	CONSTRAINT `fk_t066_prev` FOREIGN KEY (prev_id) REFERENCES `t065` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t067` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t067_code` (code),
| 	KEY `idx_t067_name` (name),
| 	KEY `idx_t067_prev` (prev_id),
| 	CONSTRAINT `fk_t067_prev` FOREIGN KEY (prev_id) REFERENCES `t066` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t068` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t068_code` (code),
| 	KEY `idx_t068_name` (name),
| 	KEY `idx_t068_prev` (prev_id),
| 	CONSTRAINT `fk_t068_prev` FOREIGN KEY (prev_id) REFERENCES `t067` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t069` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t069_code` (code),
| 	KEY `idx_t069_name` (name),
| 	KEY `idx_t069_prev` (prev_id),
| 	CONSTRAINT `fk_t069_prev` FOREIGN KEY (prev_id) REFERENCES `t068` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t070` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t070_code` (code),
| 	KEY `idx_t070_name` (name),
| 	KEY `idx_t070_prev` (prev_id),
| 	CONSTRAINT `fk_t070_prev` FOREIGN KEY (prev_id) REFERENCES `t069` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t071` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t071_code` (code),
| 	KEY `idx_t071_name` (name),
| 	KEY `idx_t071_prev` (prev_id),
| 	CONSTRAINT `fk_t071_prev` FOREIGN KEY (prev_id) REFERENCES `t070` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t072` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t072_code` (code),
| 	KEY `idx_t072_name` (name),
| 	KEY `idx_t072_prev` (prev_id),
| 	CONSTRAINT `fk_t072_prev` FOREIGN KEY (prev_id) REFERENCES `t071` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t073` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t073_code` (code),
| 	KEY `idx_t073_name` (name),
| 	KEY `idx_t073_prev` (prev_id),
| 	CONSTRAINT `fk_t073_prev` FOREIGN KEY (prev_id) REFERENCES `t072` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t074` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t074_code` (code),
| 	KEY `idx_t074_name` (name),
| 	KEY `idx_t074_prev` (prev_id),
| 	CONSTRAINT `fk_t074_prev` FOREIGN KEY (prev_id) REFERENCES `t073` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t075` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t075_code` (code),
| 	KEY `idx_t075_name` (name),
| 	KEY `idx_t075_prev` (prev_id),
| 	CONSTRAINT `fk_t075_prev` FOREIGN KEY (prev_id) REFERENCES `t074` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t076` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t076_code` (code),
| 	KEY `idx_t076_name` (name),
| 	KEY `idx_t076_prev` (prev_id),
| 	CONSTRAINT `fk_t076_prev` FOREIGN KEY (prev_id) REFERENCES `t075` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t077` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t077_code` (code),
| 	KEY `idx_t077_name` (name),
| 	KEY `idx_t077_prev` (prev_id),
| 	CONSTRAINT `fk_t077_prev` FOREIGN KEY (prev_id) REFERENCES `t076` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t078` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t078_code` (code),
| 	KEY `idx_t078_name` (name),
| 	KEY `idx_t078_prev` (prev_id),
| 	CONSTRAINT `fk_t078_prev` FOREIGN KEY (prev_id) REFERENCES `t077` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t079` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t079_code` (code),
| 	KEY `idx_t079_name` (name),
| 	KEY `idx_t079_prev` (prev_id),
| 	CONSTRAINT `fk_t079_prev` FOREIGN KEY (prev_id) REFERENCES `t078` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t080` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t080_code` (code),
| 	KEY `idx_t080_name` (name),
| 	KEY `idx_t080_prev` (prev_id),
| 	CONSTRAINT `fk_t080_prev` FOREIGN KEY (prev_id) REFERENCES `t079` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t081` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t081_code` (code),
| 	KEY `idx_t081_name` (name),
| 	KEY `idx_t081_prev` (prev_id),
| 	CONSTRAINT `fk_t081_prev` FOREIGN KEY (prev_id) REFERENCES `t080` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t082` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t082_code` (code),
| 	KEY `idx_t082_name` (name),
| 	KEY `idx_t082_prev` (prev_id),
| 	CONSTRAINT `fk_t082_prev` FOREIGN KEY (prev_id) REFERENCES `t081` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t083` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t083_code` (code),
| 	KEY `idx_t083_name` (name),
| 	KEY `idx_t083_prev` (prev_id),
| 	CONSTRAINT `fk_t083_prev` FOREIGN KEY (prev_id) REFERENCES `t082` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t084` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t084_code` (code),
| 	KEY `idx_t084_name` (name),
| 	KEY `idx_t084_prev` (prev_id),
| 	CONSTRAINT `fk_t084_prev` FOREIGN KEY (prev_id) REFERENCES `t083` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t085` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t085_code` (code),
| 	KEY `idx_t085_name` (name),
| 	KEY `idx_t085_prev` (prev_id),
| 	CONSTRAINT `fk_t085_prev` FOREIGN KEY (prev_id) REFERENCES `t084` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t086` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t086_code` (code),
| 	KEY `idx_t086_name` (name),
| 	KEY `idx_t086_prev` (prev_id),
| 	CONSTRAINT `fk_t086_prev` FOREIGN KEY (prev_id) REFERENCES `t085` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t087` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t087_code` (code),
| 	KEY `idx_t087_name` (name),
| 	KEY `idx_t087_prev` (prev_id),
| 	CONSTRAINT `fk_t087_prev` FOREIGN KEY (prev_id) REFERENCES `t086` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t088` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t088_code` (code),
| 	KEY `idx_t088_name` (name),
| 	KEY `idx_t088_prev` (prev_id),
| 	CONSTRAINT `fk_t088_prev` FOREIGN KEY (prev_id) REFERENCES `t087` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE `t089` (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	code VARCHAR(40) NOT NULL,
| 	name VARCHAR(120) NOT NULL,
| 	description TEXT,
| 	amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	qty INT NOT NULL DEFAULT 0,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	status ENUM('new','open','closed') NOT NULL DEFAULT 'new',
| 	prev_id INT NULL,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY `uk_t089_code` (code),
| 	KEY `idx_t089_name` (name),
| 	KEY `idx_t089_prev` (prev_id),
| 	CONSTRAINT `fk_t089_prev` FOREIGN KEY (prev_id) REFERENCES `t088` (id) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> commit
< ok
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR" "string"
< row
: bytes
| 8.0.46
> query
|
| 			SELECT
| 				TABLE_NAME,
| 				INDEX_NAME,
| 				COLUMN_NAME,
| 				COLLATION,
| 				IFNULL(SUB_PART, -1),
| 				EXPRESSION,
| 				SEQ_IN_INDEX,
| 				INDEX_TYPE,
| 				NON_UNIQUE,
| 				CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END,
| 				INDEX_COMMENT
| 			FROM information_schema.STATISTICS
| 			WHERE TABLE_SCHEMA = ?
| 			ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| sdl_deep_scale_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT" "int64"
< column "EXPRESSION" "TEXT"
< column "SEQ_IN_INDEX" "UNSIGNED INT" "uint32"
< column "INDEX_TYPE" "VARCHAR" "string"
< column "NON_UNIQUE" "INT" "int32"
< column "CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END" "INT" "int32"
< column "INDEX_COMMENT" "VARCHAR" "string"
< row
: bytes
| t000
: bytes
| idx_t000_amount
: bytes
| amount
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t000
: bytes
| idx_t000_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t000
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t000
: bytes
| uk_t000_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t001
: bytes
| idx_t001_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t001
: bytes
| idx_t001_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t001
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t001
: bytes
| uk_t001_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t002
: bytes
| idx_t002_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t002
: bytes
| idx_t002_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t002
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t002
: bytes
| uk_t002_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t003
: bytes
| idx_t003_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t003
: bytes
| idx_t003_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t003
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t003
: bytes
| uk_t003_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t004
: bytes
| idx_t004_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t004
: bytes
| idx_t004_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t004
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t004
: bytes
| uk_t004_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t005
: bytes
| idx_t005_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t005
: bytes
| idx_t005_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t005
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t005
: bytes
| uk_t005_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t006
: bytes
| idx_t006_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t006
: bytes
| idx_t006_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t006
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t006
: bytes
| uk_t006_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t007
: bytes
| idx_t007_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t007
: bytes
| idx_t007_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t007
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t007
: bytes
| uk_t007_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t008
: bytes
| idx_t008_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t008
: bytes
| idx_t008_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t008
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t008
: bytes
| uk_t008_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t009
: bytes
| idx_t009_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t009
: bytes
| idx_t009_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t009
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t009
: bytes
| uk_t009_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t010
: bytes
| idx_t010_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t010
: bytes
| idx_t010_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t010
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t010
: bytes
| uk_t010_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t011
: bytes
| idx_t011_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t011
: bytes
| idx_t011_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t011
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t011
: bytes
| uk_t011_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t012
: bytes
| idx_t012_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t012
: bytes
| idx_t012_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t012
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t012
: bytes
| uk_t012_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t013
: bytes
| idx_t013_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t013
: bytes
| idx_t013_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t013
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t013
: bytes
| uk_t013_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t014
: bytes
| idx_t014_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t014
: bytes
| idx_t014_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t014
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t014
: bytes
| uk_t014_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t015
: bytes
| idx_t015_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t015
: bytes
| idx_t015_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t015
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t015
: bytes
| uk_t015_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t016
: bytes
| idx_t016_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t016
: bytes
| idx_t016_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t016
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t016
: bytes
| uk_t016_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t017
: bytes
| idx_t017_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t017
: bytes
| idx_t017_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t017
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t017
: bytes
| uk_t017_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t018
: bytes
| idx_t018_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t018
: bytes
| idx_t018_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t018
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t018
: bytes
| uk_t018_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t019
: bytes
| idx_t019_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t019
: bytes
| idx_t019_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t019
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t019
: bytes
| uk_t019_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t020
: bytes
| idx_t020_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t020
: bytes
| idx_t020_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t020
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t020
: bytes
| uk_t020_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t021
: bytes
| idx_t021_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t021
: bytes
| idx_t021_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t021
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t021
: bytes
| uk_t021_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t022
: bytes
| idx_t022_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t022
: bytes
| idx_t022_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t022
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t022
: bytes
| uk_t022_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t023
: bytes
| idx_t023_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t023
: bytes
| idx_t023_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t023
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t023
: bytes
| uk_t023_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t024
: bytes
| idx_t024_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t024
: bytes
| idx_t024_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t024
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t024
: bytes
| uk_t024_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t025
: bytes
| idx_t025_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t025
: bytes
| idx_t025_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t025
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t025
: bytes
| uk_t025_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t026
: bytes
| idx_t026_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t026
: bytes
| idx_t026_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t026
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t026
: bytes
| uk_t026_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t027
: bytes
| idx_t027_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t027
: bytes
| idx_t027_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t027
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t027
: bytes
| uk_t027_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t028
: bytes
| idx_t028_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t028
: bytes
| idx_t028_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t028
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t028
: bytes
| uk_t028_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t029
: bytes
| idx_t029_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t029
: bytes
| idx_t029_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t029
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t029
: bytes
| uk_t029_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t030
: bytes
| idx_t030_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t030
: bytes
| idx_t030_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t030
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t030
: bytes
| uk_t030_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t031
: bytes
| idx_t031_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t031
: bytes
| idx_t031_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t031
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t031
: bytes
| uk_t031_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t032
: bytes
| idx_t032_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t032
: bytes
| idx_t032_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t032
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t032
: bytes
| uk_t032_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t033
: bytes
| idx_t033_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t033
: bytes
| idx_t033_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t033
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t033
: bytes
| uk_t033_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t034
: bytes
| idx_t034_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t034
: bytes
| idx_t034_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t034
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t034
: bytes
| uk_t034_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t035
: bytes
| idx_t035_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t035
: bytes
| idx_t035_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t035
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t035
: bytes
| uk_t035_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t036
: bytes
| idx_t036_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t036
: bytes
| idx_t036_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t036
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t036
: bytes
| uk_t036_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t037
: bytes
| idx_t037_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t037
: bytes
| idx_t037_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t037
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t037
: bytes
| uk_t037_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t038
: bytes
| idx_t038_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t038
: bytes
| idx_t038_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t038
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t038
: bytes
| uk_t038_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t039
: bytes
| idx_t039_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t039
: bytes
| idx_t039_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t039
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t039
: bytes
| uk_t039_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t040
: bytes
| idx_t040_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t040
: bytes
| idx_t040_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t040
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t040
: bytes
| uk_t040_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t041
: bytes
| idx_t041_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t041
: bytes
| idx_t041_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t041
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t041
: bytes
| uk_t041_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t042
: bytes
| idx_t042_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t042
: bytes
| idx_t042_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t042
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t042
: bytes
| uk_t042_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t043
: bytes
| idx_t043_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t043
: bytes
| idx_t043_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t043
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t043
: bytes
| uk_t043_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t044
: bytes
| idx_t044_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t044
: bytes
| idx_t044_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t044
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t044
: bytes
| uk_t044_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t045
: bytes
| idx_t045_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t045
: bytes
| idx_t045_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t045
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t045
: bytes
| uk_t045_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t046
: bytes
| idx_t046_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t046
: bytes
| idx_t046_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t046
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t046
: bytes
| uk_t046_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t047
: bytes
| idx_t047_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t047
: bytes
| idx_t047_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t047
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t047
: bytes
| uk_t047_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t048
: bytes
| idx_t048_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t048
: bytes
| idx_t048_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t048
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t048
: bytes
| uk_t048_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t049
: bytes
| idx_t049_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t049
: bytes
| idx_t049_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t049
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t049
: bytes
| uk_t049_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t050
: bytes
| idx_t050_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t050
: bytes
| idx_t050_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t050
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t050
: bytes
| uk_t050_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t051
: bytes
| idx_t051_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t051
: bytes
| idx_t051_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t051
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t051
: bytes
| uk_t051_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t052
: bytes
| idx_t052_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t052
: bytes
| idx_t052_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t052
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t052
: bytes
| uk_t052_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t053
: bytes
| idx_t053_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t053
: bytes
| idx_t053_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t053
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t053
: bytes
| uk_t053_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t054
: bytes
| idx_t054_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t054
: bytes
| idx_t054_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t054
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t054
: bytes
| uk_t054_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t055
: bytes
| idx_t055_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t055
: bytes
| idx_t055_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t055
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t055
: bytes
| uk_t055_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t056
: bytes
| idx_t056_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t056
: bytes
| idx_t056_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t056
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t056
: bytes
| uk_t056_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t057
: bytes
| idx_t057_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t057
: bytes
| idx_t057_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t057
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t057
: bytes
| uk_t057_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t058
: bytes
| idx_t058_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t058
: bytes
| idx_t058_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t058
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t058
: bytes
| uk_t058_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t059
: bytes
| idx_t059_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t059
: bytes
| idx_t059_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t059
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t059
: bytes
| uk_t059_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t060
: bytes
| idx_t060_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t060
: bytes
| idx_t060_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t060
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t060
: bytes
| uk_t060_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t061
: bytes
| idx_t061_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t061
: bytes
| idx_t061_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t061
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t061
: bytes
| uk_t061_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t062
: bytes
| idx_t062_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t062
: bytes
| idx_t062_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t062
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t062
: bytes
| uk_t062_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t063
: bytes
| idx_t063_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t063
: bytes
| idx_t063_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t063
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t063
: bytes
| uk_t063_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t064
: bytes
| idx_t064_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t064
: bytes
| idx_t064_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t064
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t064
: bytes
| uk_t064_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t065
: bytes
| idx_t065_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t065
: bytes
| idx_t065_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t065
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t065
: bytes
| uk_t065_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t066
: bytes
| idx_t066_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t066
: bytes
| idx_t066_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t066
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t066
: bytes
| uk_t066_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t067
: bytes
| idx_t067_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t067
: bytes
| idx_t067_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t067
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t067
: bytes
| uk_t067_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t068
: bytes
| idx_t068_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t068
: bytes
| idx_t068_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t068
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t068
: bytes
| uk_t068_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t069
: bytes
| idx_t069_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t069
: bytes
| idx_t069_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t069
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t069
: bytes
| uk_t069_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t070
: bytes
| idx_t070_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t070
: bytes
| idx_t070_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t070
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t070
: bytes
| uk_t070_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t071
: bytes
| idx_t071_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t071
: bytes
| idx_t071_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t071
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t071
: bytes
| uk_t071_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t072
: bytes
| idx_t072_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t072
: bytes
| idx_t072_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t072
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t072
: bytes
| uk_t072_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t073
: bytes
| idx_t073_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t073
: bytes
| idx_t073_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t073
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t073
: bytes
| uk_t073_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t074
: bytes
| idx_t074_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t074
: bytes
| idx_t074_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t074
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t074
: bytes
| uk_t074_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t075
: bytes
| idx_t075_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t075
: bytes
| idx_t075_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t075
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t075
: bytes
| uk_t075_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t076
: bytes
| idx_t076_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t076
: bytes
| idx_t076_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t076
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t076
: bytes
| uk_t076_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t077
: bytes
| idx_t077_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t077
: bytes
| idx_t077_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t077
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t077
: bytes
| uk_t077_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t078
: bytes
| idx_t078_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t078
: bytes
| idx_t078_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t078
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t078
: bytes
| uk_t078_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t079
: bytes
| idx_t079_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t079
: bytes
| idx_t079_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t079
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t079
: bytes
| uk_t079_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t080
: bytes
| idx_t080_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t080
: bytes
| idx_t080_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t080
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t080
: bytes
| uk_t080_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t081
: bytes
| idx_t081_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t081
: bytes
| idx_t081_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t081
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t081
: bytes
| uk_t081_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t082
: bytes
| idx_t082_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t082
: bytes
| idx_t082_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t082
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t082
: bytes
| uk_t082_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t083
: bytes
| idx_t083_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t083
: bytes
| idx_t083_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t083
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t083
: bytes
| uk_t083_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t084
: bytes
| idx_t084_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t084
: bytes
| idx_t084_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t084
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t084
: bytes
| uk_t084_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t085
: bytes
| idx_t085_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t085
: bytes
| idx_t085_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t085
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t085
: bytes
| uk_t085_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t086
: bytes
| idx_t086_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t086
: bytes
| idx_t086_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t086
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t086
: bytes
| uk_t086_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t087
: bytes
| idx_t087_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t087
: bytes
| idx_t087_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t087
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t087
: bytes
| uk_t087_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t088
: bytes
| idx_t088_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t088
: bytes
| idx_t088_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t088
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t088
: bytes
| uk_t088_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t089
: bytes
| idx_t089_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t089
: bytes
| idx_t089_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t089
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t089
: bytes
| uk_t089_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			IFNULL(COLUMN_NAME, ''),
| 			ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE,
| 			COLUMN_TYPE,
| 			IFNULL(CHARACTER_SET_NAME, ''),
| 			IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT),
| 			convert(GENERATION_EXPRESSION using BINARY),
| 			EXTRA,
| 			SRS_ID
| 		FROM information_schema.COLUMNS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| sdl_deep_scale_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR" "string"
< column "ORDINAL_POSITION" "UNSIGNED INT" "uint32"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "TEXT"
< column "IS_NULLABLE" "VARCHAR" "string"
< column "COLUMN_TYPE" "TEXT" "string"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR" "string"
< column "QUOTE(COLUMN_COMMENT)" "TEXT"
< column "convert(GENERATION_EXPRESSION using BINARY)" "BLOB" "[]uint8"
< column "EXTRA" "VARCHAR"
< column "SRS_ID" "UNSIGNED INT"
< row
: bytes
| t000
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t000
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t000
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t000
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t000
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t000
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t000
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t000
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t000
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t000
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t001
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t001
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t001
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t001
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t001
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t001
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t001
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t001
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t001
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t001
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t001
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t002
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t002
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t002
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t002
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t002
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t002
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t002
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t002
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t002
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t002
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t002
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t003
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t003
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t003
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t003
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t003
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t003
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t003
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t003
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t003
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t003
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t003
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t004
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t004
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t004
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t004
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t004
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t004
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t004
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t004
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t004
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t004
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t004
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t005
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t005
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t005
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t005
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t005
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t005
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t005
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t005
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t005
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t005
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t005
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t006
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t006
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t006
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t006
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t006
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t006
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t006
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t006
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t006
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t006
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t006
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t007
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t007
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t007
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t007
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t007
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t007
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t007
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t007
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t007
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t007
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t007
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t008
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t008
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t008
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t008
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t008
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t008
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t008
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t008
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t008
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t008
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t008
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t009
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t009
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t009
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t009
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t009
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t009
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t009
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t009
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t009
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t009
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t009
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t010
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t010
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t010
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t010
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t010
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t010
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t010
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t010
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t010
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t010
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t010
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t011
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t011
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t011
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t011
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t011
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t011
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t011
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t011
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t011
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t011
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t011
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t012
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t012
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t012
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t012
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t012
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t012
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t012
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t012
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t012
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t012
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t012
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t013
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t013
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t013
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t013
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t013
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t013
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t013
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t013
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t013
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t013
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t013
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t014
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t014
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t014
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t014
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t014
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t014
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t014
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t014
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t014
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t014
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t014
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t015
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t015
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t015
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t015
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t015
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t015
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t015
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t015
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t015
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t015
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t015
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t016
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t016
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t016
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t016
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t016
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t016
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t016
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t016
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t016
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t016
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t016
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t017
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t017
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t017
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t017
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t017
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t017
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t017
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t017
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t017
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t017
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t017
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t018
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t018
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t018
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t018
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t018
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t018
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t018
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t018
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t018
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t018
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t018
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t019
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t019
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t019
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t019
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t019
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t019
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t019
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t019
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t019
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t019
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t019
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t020
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t020
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t020
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t020
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t020
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t020
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t020
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t020
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t020
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t020
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t020
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t021
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t021
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t021
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t021
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t021
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t021
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t021
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t021
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t021
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t021
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t021
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t022
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t022
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t022
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t022
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t022
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t022
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t022
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t022
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t022
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t022
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t022
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t023
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t023
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t023
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t023
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t023
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t023
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t023
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t023
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t023
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t023
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t023
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t024
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t024
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t024
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t024
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t024
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t024
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t024
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t024
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t024
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t024
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t024
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t025
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t025
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t025
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t025
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t025
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t025
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t025
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t025
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t025
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t025
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t025
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t026
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t026
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t026
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t026
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t026
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t026
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t026
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t026
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t026
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t026
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t026
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t027
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t027
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t027
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t027
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t027
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t027
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t027
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t027
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t027
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t027
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t027
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t028
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t028
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t028
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t028
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t028
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t028
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t028
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t028
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t028
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t028
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t028
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t029
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t029
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t029
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t029
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t029
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t029
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t029
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t029
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t029
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t029
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t029
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t030
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t030
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t030
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t030
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t030
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t030
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t030
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t030
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t030
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t030
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t030
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t031
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t031
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t031
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t031
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t031
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t031
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t031
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t031
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t031
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t031
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t031
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t032
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t032
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t032
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t032
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t032
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t032
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t032
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t032
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t032
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t032
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t032
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t033
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t033
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t033
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t033
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t033
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t033
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t033
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t033
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t033
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t033
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t033
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t034
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t034
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t034
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t034
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t034
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t034
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t034
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t034
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t034
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t034
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t034
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t035
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t035
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t035
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t035
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t035
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t035
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t035
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t035
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t035
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t035
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t035
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t036
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t036
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t036
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t036
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t036
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t036
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t036
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t036
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t036
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t036
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t036
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t037
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t037
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t037
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t037
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t037
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t037
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t037
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t037
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t037
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t037
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t037
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t038
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t038
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t038
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t038
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t038
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t038
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t038
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t038
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t038
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t038
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t038
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t039
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t039
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t039
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t039
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t039
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t039
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t039
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t039
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t039
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t039
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t039
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t040
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t040
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t040
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t040
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t040
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t040
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t040
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t040
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t040
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t040
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t040
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t041
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t041
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t041
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t041
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t041
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t041
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t041
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t041
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t041
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t041
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t041
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t042
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t042
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t042
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t042
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t042
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t042
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t042
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t042
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t042
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t042
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t042
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t043
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t043
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t043
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t043
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t043
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t043
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t043
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t043
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t043
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t043
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t043
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t044
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t044
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t044
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t044
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t044
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t044
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t044
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t044
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t044
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t044
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t044
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t045
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t045
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t045
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t045
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t045
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t045
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t045
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t045
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t045
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t045
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t045
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t046
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t046
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t046
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t046
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t046
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t046
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t046
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t046
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t046
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t046
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t046
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t047
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t047
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t047
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t047
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t047
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t047
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t047
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t047
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t047
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t047
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t047
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t048
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t048
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t048
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t048
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t048
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t048
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t048
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t048
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t048
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t048
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t048
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t049
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t049
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t049
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t049
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t049
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t049
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t049
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t049
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t049
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t049
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t049
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t050
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t050
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t050
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t050
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t050
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t050
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t050
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t050
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t050
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t050
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t050
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t051
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t051
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t051
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t051
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t051
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t051
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t051
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t051
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t051
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t051
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t051
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t052
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t052
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t052
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t052
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t052
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t052
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t052
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t052
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t052
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t052
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t052
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t053
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t053
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t053
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t053
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t053
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t053
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t053
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t053
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t053
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t053
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t053
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t054
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t054
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t054
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t054
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t054
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t054
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t054
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t054
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t054
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t054
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t054
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t055
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t055
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t055
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t055
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t055
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t055
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t055
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t055
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t055
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t055
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t055
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t056
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t056
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t056
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t056
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t056
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t056
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t056
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t056
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t056
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t056
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t056
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t057
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t057
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t057
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t057
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t057
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t057
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t057
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t057
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t057
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t057
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t057
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t058
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t058
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t058
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t058
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t058
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t058
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t058
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t058
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t058
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t058
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t058
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t059
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t059
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t059
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t059
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t059
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t059
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t059
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t059
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t059
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t059
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t059
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t060
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t060
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t060
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t060
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t060
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t060
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t060
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t060
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t060
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t060
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t060
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t061
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t061
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t061
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t061
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t061
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t061
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t061
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t061
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t061
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t061
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t061
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t062
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t062
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t062
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t062
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t062
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t062
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t062
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t062
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t062
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t062
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t062
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t063
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t063
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t063
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t063
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t063
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t063
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t063
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t063
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t063
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t063
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t063
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t064
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t064
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t064
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t064
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t064
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t064
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t064
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t064
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t064
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t064
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t064
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t065
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t065
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t065
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t065
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t065
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t065
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t065
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t065
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t065
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t065
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t065
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t066
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t066
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t066
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t066
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t066
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t066
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t066
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t066
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t066
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t066
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t066
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t067
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t067
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t067
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t067
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t067
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t067
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t067
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t067
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t067
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t067
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t067
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t068
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t068
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t068
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t068
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t068
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t068
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t068
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t068
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t068
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t068
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t068
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t069
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t069
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t069
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t069
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t069
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t069
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t069
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t069
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t069
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t069
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t069
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t070
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t070
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t070
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t070
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t070
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t070
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t070
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t070
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t070
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t070
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t070
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t071
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t071
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t071
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t071
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t071
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t071
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t071
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t071
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t071
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t071
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t071
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t072
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t072
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t072
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t072
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t072
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t072
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t072
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t072
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t072
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t072
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t072
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t073
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t073
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t073
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t073
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t073
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t073
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t073
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t073
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t073
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t073
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t073
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t074
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t074
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t074
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t074
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t074
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t074
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t074
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t074
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t074
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t074
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t074
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t075
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t075
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t075
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t075
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t075
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t075
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t075
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t075
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t075
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t075
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t075
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t076
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t076
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t076
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t076
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t076
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t076
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t076
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t076
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t076
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t076
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t076
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t077
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t077
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t077
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t077
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t077
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t077
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t077
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t077
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t077
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t077
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t077
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t078
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t078
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t078
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t078
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t078
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t078
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t078
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t078
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t078
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t078
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t078
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t079
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t079
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t079
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t079
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t079
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t079
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t079
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t079
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t079
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t079
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t079
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t080
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t080
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t080
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t080
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t080
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t080
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t080
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t080
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t080
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t080
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t080
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t081
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t081
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t081
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t081
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t081
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t081
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t081
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t081
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t081
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t081
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t081
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t082
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t082
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t082
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t082
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t082
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t082
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t082
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t082
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t082
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t082
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t082
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t083
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t083
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t083
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t083
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t083
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t083
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t083
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t083
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t083
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t083
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t083
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t084
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t084
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t084
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t084
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t084
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t084
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t084
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t084
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t084
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t084
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t084
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t085
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t085
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t085
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t085
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t085
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t085
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t085
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t085
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t085
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t085
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t085
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t086
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t086
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t086
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t086
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t086
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t086
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t086
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t086
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t086
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t086
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t086
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t087
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t087
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t087
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t087
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t087
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t087
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t087
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t087
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t087
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t087
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t087
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t088
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t088
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t088
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t088
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t088
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t088
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t088
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t088
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t088
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t088
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t088
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t089
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t089
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t089
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t089
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t089
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t089
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t089
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t089
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t089
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t089
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t089
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
> query
|
| 			SELECT
| 				tc.TABLE_NAME,
| 				cc.CONSTRAINT_NAME,
| 				cc.CHECK_CLAUSE,
| 				tc.ENFORCED
| 			FROM information_schema.CHECK_CONSTRAINTS cc
| 				JOIN information_schema.TABLE_CONSTRAINTS tc
| 					ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 					AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 			WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ?
| 			ORDER BY tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		
: string
| sdl_deep_scale_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< column "ENFORCED" "VARCHAR" "string"
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| sdl_deep_scale_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT"
> query
|
| 	SELECT
| 		TRIGGER_NAME,
| 		EVENT_OBJECT_TABLE,
| 		EVENT_MANIPULATION,
| 		ACTION_TIMING,
| 		ACTION_STATEMENT,
| 		SQL_MODE,
| 		CHARACTER_SET_CLIENT,
| 		COLLATION_CONNECTION
| 	FROM INFORMATION_SCHEMA.TRIGGERS
| 	WHERE EVENT_OBJECT_SCHEMA = ?
| 	ORDER BY EVENT_OBJECT_TABLE ASC, EVENT_MANIPULATION ASC, ACTION_TIMING ASC, ACTION_ORDER ASC;
| 	
: string
| sdl_deep_scale_1
< result
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "EVENT_OBJECT_TABLE" "VARCHAR" "string"
< column "EVENT_MANIPULATION" "ENUM" "string"
< column "ACTION_TIMING" "ENUM" "string"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
> query
|
| 	SELECT
| 		EVENT_NAME,
| 		TIME_ZONE,
| 		SQL_MODE,
| 		CHARACTER_SET_CLIENT,
| 		COLLATION_CONNECTION
| 	FROM INFORMATION_SCHEMA.EVENTS
| 	WHERE EVENT_SCHEMA = ?
| 	ORDER BY EVENT_NAME ASC;
| 	
: string
| sdl_deep_scale_1
< result
< column "EVENT_NAME" "VARCHAR" "string"
< column "TIME_ZONE" "VARCHAR" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			CONSTRAINT_NAME,
| 			REFERENCED_TABLE_NAME,
| 			DELETE_RULE,
| 			UPDATE_RULE,
| 			MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ?;
| 	
: string
| sdl_deep_scale_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
< row
: bytes
| t001
: bytes
| fk_t001_prev
: bytes
| t000
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t002
: bytes
| fk_t002_prev
: bytes
| t001
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t003
: bytes
| fk_t003_prev
: bytes
| t002
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t004
: bytes
| fk_t004_prev
: bytes
| t003
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t005
: bytes
| fk_t005_prev
: bytes
| t004
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t006
: bytes
| fk_t006_prev
: bytes
| t005
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t007
: bytes
| fk_t007_prev
: bytes
| t006
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t008
: bytes
| fk_t008_prev
: bytes
| t007
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t009
: bytes
| fk_t009_prev
: bytes
| t008
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t010
: bytes
| fk_t010_prev
: bytes
| t009
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t011
: bytes
| fk_t011_prev
: bytes
| t010
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t012
: bytes
| fk_t012_prev
: bytes
| t011
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t013
: bytes
| fk_t013_prev
: bytes
| t012
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t014
: bytes
| fk_t014_prev
: bytes
| t013
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t015
: bytes
| fk_t015_prev
: bytes
| t014
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t016
: bytes
| fk_t016_prev
: bytes
| t015
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t017
: bytes
| fk_t017_prev
: bytes
| t016
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t018
: bytes
| fk_t018_prev
: bytes
| t017
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t019
: bytes
| fk_t019_prev
: bytes
| t018
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t020
: bytes
| fk_t020_prev
: bytes
| t019
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t021
: bytes
| fk_t021_prev
: bytes
| t020
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t022
: bytes
| fk_t022_prev
: bytes
| t021
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t023
: bytes
| fk_t023_prev
: bytes
| t022
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t024
: bytes
| fk_t024_prev
: bytes
| t023
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t025
: bytes
| fk_t025_prev
: bytes
| t024
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t026
: bytes
| fk_t026_prev
: bytes
| t025
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t027
: bytes
| fk_t027_prev
: bytes
| t026
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t028
: bytes
| fk_t028_prev
: bytes
| t027
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t029
: bytes
| fk_t029_prev
: bytes
| t028
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t030
: bytes
| fk_t030_prev
: bytes
| t029
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t031
: bytes
| fk_t031_prev
: bytes
| t030
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t032
: bytes
| fk_t032_prev
: bytes
| t031
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t033
: bytes
| fk_t033_prev
: bytes
| t032
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t034
: bytes
| fk_t034_prev
: bytes
| t033
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t035
: bytes
| fk_t035_prev
: bytes
| t034
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t036
: bytes
| fk_t036_prev
: bytes
| t035
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t037
: bytes
| fk_t037_prev
: bytes
| t036
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t038
: bytes
| fk_t038_prev
: bytes
| t037
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t039
: bytes
| fk_t039_prev
: bytes
| t038
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t040
: bytes
| fk_t040_prev
: bytes
| t039
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t041
: bytes
| fk_t041_prev
: bytes
| t040
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t042
: bytes
| fk_t042_prev
: bytes
| t041
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t043
: bytes
| fk_t043_prev
: bytes
| t042
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t044
: bytes
| fk_t044_prev
: bytes
| t043
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t045
: bytes
| fk_t045_prev
: bytes
| t044
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t046
: bytes
| fk_t046_prev
: bytes
| t045
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t047
: bytes
| fk_t047_prev
: bytes
| t046
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t048
: bytes
| fk_t048_prev
: bytes
| t047
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t049
: bytes
| fk_t049_prev
: bytes
| t048
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t050
: bytes
| fk_t050_prev
: bytes
| t049
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t051
: bytes
| fk_t051_prev
: bytes
| t050
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t052
: bytes
| fk_t052_prev
: bytes
| t051
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t053
: bytes
| fk_t053_prev
: bytes
| t052
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t054
: bytes
| fk_t054_prev
: bytes
| t053
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t055
: bytes
| fk_t055_prev
: bytes
| t054
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t056
: bytes
| fk_t056_prev
: bytes
| t055
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t057
: bytes
| fk_t057_prev
: bytes
| t056
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t058
: bytes
| fk_t058_prev
: bytes
| t057
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t059
: bytes
| fk_t059_prev
: bytes
| t058
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t060
: bytes
| fk_t060_prev
: bytes
| t059
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t061
: bytes
| fk_t061_prev
: bytes
| t060
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t062
: bytes
| fk_t062_prev
: bytes
| t061
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t063
: bytes
| fk_t063_prev
: bytes
| t062
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t064
: bytes
| fk_t064_prev
: bytes
| t063
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t065
: bytes
| fk_t065_prev
: bytes
| t064
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t066
: bytes
| fk_t066_prev
: bytes
| t065
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t067
: bytes
| fk_t067_prev
: bytes
| t066
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t068
: bytes
| fk_t068_prev
: bytes
| t067
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t069
: bytes
| fk_t069_prev
: bytes
| t068
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t070
: bytes
| fk_t070_prev
: bytes
| t069
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t071
: bytes
| fk_t071_prev
: bytes
| t070
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t072
: bytes
| fk_t072_prev
: bytes
| t071
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t073
: bytes
| fk_t073_prev
: bytes
| t072
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t074
: bytes
| fk_t074_prev
: bytes
| t073
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t075
: bytes
| fk_t075_prev
: bytes
| t074
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t076
: bytes
| fk_t076_prev
: bytes
| t075
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t077
: bytes
| fk_t077_prev
: bytes
| t076
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t078
: bytes
| fk_t078_prev
: bytes
| t077
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t079
: bytes
| fk_t079_prev
: bytes
| t078
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t080
: bytes
| fk_t080_prev
: bytes
| t079
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t081
: bytes
| fk_t081_prev
: bytes
| t080
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t082
: bytes
| fk_t082_prev
: bytes
| t081
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t083
: bytes
| fk_t083_prev
: bytes
| t082
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t084
: bytes
| fk_t084_prev
: bytes
| t083
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t085
: bytes
| fk_t085_prev
: bytes
| t084
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t086
: bytes
| fk_t086_prev
: bytes
| t085
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t087
: bytes
| fk_t087_prev
: bytes
| t086
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t088
: bytes
| fk_t088_prev
: bytes
| t087
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t089
: bytes
| fk_t089_prev
: bytes
| t088
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
> query
|
| 		SELECT
| 			TABLE_SCHEMA,
| 			TABLE_NAME,
| 			CONSTRAINT_NAME,
| 			COLUMN_NAME,
| 			REFERENCED_TABLE_SCHEMA,
| 			REFERENCED_COLUMN_NAME
| 		FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
| 		WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION;
| 	
: string
| sdl_deep_scale_1
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| sdl_deep_scale_1
: bytes
| t001
: bytes
| fk_t001_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t002
: bytes
| fk_t002_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t003
: bytes
| fk_t003_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t004
: bytes
| fk_t004_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t005
: bytes
| fk_t005_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t006
: bytes
| fk_t006_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t007
: bytes
| fk_t007_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t008
: bytes
| fk_t008_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t009
: bytes
| fk_t009_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t010
: bytes
| fk_t010_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t011
: bytes
| fk_t011_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t012
: bytes
| fk_t012_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t013
: bytes
| fk_t013_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t014
: bytes
| fk_t014_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t015
: bytes
| fk_t015_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t016
: bytes
| fk_t016_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t017
: bytes
| fk_t017_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t018
: bytes
| fk_t018_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t019
: bytes
| fk_t019_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t020
: bytes
| fk_t020_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t021
: bytes
| fk_t021_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t022
: bytes
| fk_t022_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t023
: bytes
| fk_t023_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t024
: bytes
| fk_t024_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t025
: bytes
| fk_t025_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t026
: bytes
| fk_t026_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t027
: bytes
| fk_t027_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t028
: bytes
| fk_t028_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t029
: bytes
| fk_t029_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t030
: bytes
| fk_t030_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t031
: bytes
| fk_t031_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t032
: bytes
| fk_t032_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t033
: bytes
| fk_t033_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t034
: bytes
| fk_t034_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t035
: bytes
| fk_t035_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t036
: bytes
| fk_t036_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t037
: bytes
| fk_t037_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t038
: bytes
| fk_t038_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t039
: bytes
| fk_t039_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t040
: bytes
| fk_t040_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t041
: bytes
| fk_t041_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t042
: bytes
| fk_t042_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t043
: bytes
| fk_t043_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t044
: bytes
| fk_t044_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t045
: bytes
| fk_t045_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t046
: bytes
| fk_t046_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t047
: bytes
| fk_t047_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t048
: bytes
| fk_t048_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t049
: bytes
| fk_t049_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t050
: bytes
| fk_t050_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t051
: bytes
| fk_t051_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t052
: bytes
| fk_t052_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t053
: bytes
| fk_t053_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t054
: bytes
| fk_t054_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t055
: bytes
| fk_t055_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t056
: bytes
| fk_t056_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t057
: bytes
| fk_t057_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t058
: bytes
| fk_t058_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t059
: bytes
| fk_t059_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t060
: bytes
| fk_t060_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t061
: bytes
| fk_t061_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t062
: bytes
| fk_t062_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t063
: bytes
| fk_t063_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t064
: bytes
| fk_t064_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t065
: bytes
| fk_t065_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t066
: bytes
| fk_t066_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t067
: bytes
| fk_t067_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t068
: bytes
| fk_t068_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t069
: bytes
| fk_t069_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t070
: bytes
| fk_t070_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t071
: bytes
| fk_t071_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t072
: bytes
| fk_t072_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t073
: bytes
| fk_t073_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t074
: bytes
| fk_t074_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t075
: bytes
| fk_t075_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t076
: bytes
| fk_t076_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t077
: bytes
| fk_t077_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t078
: bytes
| fk_t078_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t079
: bytes
| fk_t079_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t080
: bytes
| fk_t080_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t081
: bytes
| fk_t081_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t082
: bytes
| fk_t082_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t083
: bytes
| fk_t083_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t084
: bytes
| fk_t084_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t085
: bytes
| fk_t085_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t086
: bytes
| fk_t086_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t087
: bytes
| fk_t087_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t088
: bytes
| fk_t088_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t089
: bytes
| fk_t089_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			PARTITION_NAME,
| 			SUBPARTITION_NAME,
| 			PARTITION_METHOD,
| 			SUBPARTITION_METHOD,
| 			PARTITION_EXPRESSION,
| 			SUBPARTITION_EXPRESSION,
| 			PARTITION_DESCRIPTION,
| 			PARTITION_COMMENT
| 		FROM INFORMATION_SCHEMA.PARTITIONS
| 		WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME ASC, PARTITION_ORDINAL_POSITION ASC, SUBPARTITION_ORDINAL_POSITION ASC;
| 	
: string
| sdl_deep_scale_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "SUBPARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "VARCHAR"
< column "SUBPARTITION_EXPRESSION" "VARCHAR"
< column "PARTITION_DESCRIPTION" "TEXT"
< column "PARTITION_COMMENT" "TEXT" "string"
> query
|
| 		SELECT
| 			ROUTINE_NAME,
| 			ROUTINE_TYPE,
| 			SQL_MODE,
| 			CHARACTER_SET_CLIENT,
| 			COLLATION_CONNECTION,
| 			DATABASE_COLLATION,
| 			IFNULL(ROUTINE_COMMENT, '')
| 		FROM
| 			INFORMATION_SCHEMA.ROUTINES
| 		WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE IN ('FUNCTION', 'PROCEDURE')
| 		ORDER BY ROUTINE_TYPE, ROUTINE_NAME;
| 	
: string
| sdl_deep_scale_1
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "ENUM" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT" "string"
> query
|
| 		SELECT
| 			TABLES.TABLE_NAME,
| 			TABLES.TABLE_TYPE,
| 			IFNULL(TABLES.ENGINE, ''),
| 			IFNULL(TABLES.TABLE_COLLATION, ''),
| 			IFNULL(TABLES.TABLE_ROWS, 0),
| 			IFNULL(TABLES.DATA_LENGTH, 0),
| 			IFNULL(TABLES.INDEX_LENGTH, 0),
| 			IFNULL(TABLES.DATA_FREE, 0),
| 			IFNULL(TABLES.CREATE_OPTIONS, ''),
| 			QUOTE(IFNULL(TABLES.TABLE_COMMENT, '')),
| 			IFNULL(CCSA.CHARACTER_SET_NAME, '')
| 		FROM information_schema.TABLES TABLES
| 		LEFT JOIN information_schema.COLLATION_CHARACTER_SET_APPLICABILITY CCSA
| 		ON TABLES.TABLE_COLLATION = CCSA.COLLATION_NAME
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| sdl_deep_scale_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "TABLE_TYPE" "ENUM" "string"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLES.TABLE_COMMENT, ''))" "TEXT"
< column "IFNULL(CCSA.CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< row
: bytes
| t000
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t001
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t002
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t003
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t004
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t005
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t006
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t007
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t008
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t009
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t010
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t011
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t012
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t013
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t014
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t015
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t016
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t017
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t018
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t019
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t020
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t021
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t022
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t023
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t024
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t025
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t026
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t027
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t028
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t029
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t030
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t031
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t032
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t033
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t034
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t035
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t036
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t037
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t038
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t039
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t040
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t041
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t042
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t043
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t044
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t045
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t046
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t047
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t048
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t049
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t050
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t051
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t052
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t053
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t054
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t055
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t056
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t057
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t058
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t059
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t060
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t061
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t062
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t063
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t064
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t065
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t066
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t067
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t068
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t069
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t070
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t071
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t072
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t073
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t074
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t075
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t076
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t077
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t078
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t079
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t080
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t081
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t082
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t083
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t084
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t085
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t086
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t087
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t088
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t089
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
> query
|
| 		SELECT t.NAME
| 		FROM information_schema.INNODB_TABLES t
| 		LEFT JOIN information_schema.INNODB_TABLESPACES ts ON t.SPACE = ts.SPACE
| 		WHERE t.NAME LIKE CONCAT(?, '/%') AND (ts.SPACE_TYPE IS NULL OR ts.SPACE_TYPE <> 'Single')
: string
| sdl_deep_scale_1
< result
< column "NAME" "VARCHAR" "string"
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| sdl_deep_scale_1
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci

== TestSDLDeepScale default#2
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR" "string"
< row
: bytes
| 8.0.46
> query
|
| 			SELECT
| 				TABLE_NAME,
| 				INDEX_NAME,
| 				COLUMN_NAME,
| 				COLLATION,
| 				IFNULL(SUB_PART, -1),
| 				EXPRESSION,
| 				SEQ_IN_INDEX,
| 				INDEX_TYPE,
| 				NON_UNIQUE,
| 				CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END,
| 				INDEX_COMMENT
| 			FROM information_schema.STATISTICS
| 			WHERE TABLE_SCHEMA = ?
| 			ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| sdl_deep_scale_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT" "int64"
< column "EXPRESSION" "TEXT"
< column "SEQ_IN_INDEX" "UNSIGNED INT" "uint32"
< column "INDEX_TYPE" "VARCHAR" "string"
< column "NON_UNIQUE" "INT" "int32"
< column "CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END" "INT" "int32"
< column "INDEX_COMMENT" "VARCHAR" "string"
< row
: bytes
| t000
: bytes
| idx_t000_amount
: bytes
| amount
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t000
: bytes
| idx_t000_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t000
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t000
: bytes
| uk_t000_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t001
: bytes
| idx_t001_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t001
: bytes
| idx_t001_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t001
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t001
: bytes
| uk_t001_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t002
: bytes
| idx_t002_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t002
: bytes
| idx_t002_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t002
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t002
: bytes
| uk_t002_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t003
: bytes
| idx_t003_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t003
: bytes
| idx_t003_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t003
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t003
: bytes
| uk_t003_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t004
: bytes
| idx_t004_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t004
: bytes
| idx_t004_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t004
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t004
: bytes
| uk_t004_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t005
: bytes
| idx_t005_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t005
: bytes
| idx_t005_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t005
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t005
: bytes
| uk_t005_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t006
: bytes
| idx_t006_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t006
: bytes
| idx_t006_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t006
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t006
: bytes
| uk_t006_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t007
: bytes
| idx_t007_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t007
: bytes
| idx_t007_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t007
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t007
: bytes
| uk_t007_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t008
: bytes
| idx_t008_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t008
: bytes
| idx_t008_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t008
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t008
: bytes
| uk_t008_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t009
: bytes
| idx_t009_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t009
: bytes
| idx_t009_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t009
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t009
: bytes
| uk_t009_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t010
: bytes
| idx_t010_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t010
: bytes
| idx_t010_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t010
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t010
: bytes
| uk_t010_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t011
: bytes
| idx_t011_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t011
: bytes
| idx_t011_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t011
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t011
: bytes
| uk_t011_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t012
: bytes
| idx_t012_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t012
: bytes
| idx_t012_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t012
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t012
: bytes
| uk_t012_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t013
: bytes
| idx_t013_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t013
: bytes
| idx_t013_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t013
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t013
: bytes
| uk_t013_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t014
: bytes
| idx_t014_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t014
: bytes
| idx_t014_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t014
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t014
: bytes
| uk_t014_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t015
: bytes
| idx_t015_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t015
: bytes
| idx_t015_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t015
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t015
: bytes
| uk_t015_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t016
: bytes
| idx_t016_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t016
: bytes
| idx_t016_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t016
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t016
: bytes
| uk_t016_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t017
: bytes
| idx_t017_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t017
: bytes
| idx_t017_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t017
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t017
: bytes
| uk_t017_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t018
: bytes
| idx_t018_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t018
: bytes
| idx_t018_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t018
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t018
: bytes
| uk_t018_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t019
: bytes
| idx_t019_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t019
: bytes
| idx_t019_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t019
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t019
: bytes
| uk_t019_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t020
: bytes
| idx_t020_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t020
: bytes
| idx_t020_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t020
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t020
: bytes
| uk_t020_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t021
: bytes
| idx_t021_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t021
: bytes
| idx_t021_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t021
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t021
: bytes
| uk_t021_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t022
: bytes
| idx_t022_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t022
: bytes
| idx_t022_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t022
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t022
: bytes
| uk_t022_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t023
: bytes
| idx_t023_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t023
: bytes
| idx_t023_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t023
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t023
: bytes
| uk_t023_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t024
: bytes
| idx_t024_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t024
: bytes
| idx_t024_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t024
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t024
: bytes
| uk_t024_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t025
: bytes
| idx_t025_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t025
: bytes
| idx_t025_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t025
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t025
: bytes
| uk_t025_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t026
: bytes
| idx_t026_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t026
: bytes
| idx_t026_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t026
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t026
: bytes
| uk_t026_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t027
: bytes
| idx_t027_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t027
: bytes
| idx_t027_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t027
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t027
: bytes
| uk_t027_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t028
: bytes
| idx_t028_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t028
: bytes
| idx_t028_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t028
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t028
: bytes
| uk_t028_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t029
: bytes
| idx_t029_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t029
: bytes
| idx_t029_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t029
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t029
: bytes
| uk_t029_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t030
: bytes
| idx_t030_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t030
: bytes
| idx_t030_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t030
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t030
: bytes
| uk_t030_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t031
: bytes
| idx_t031_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t031
: bytes
| idx_t031_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t031
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t031
: bytes
| uk_t031_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t032
: bytes
| idx_t032_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t032
: bytes
| idx_t032_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t032
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t032
: bytes
| uk_t032_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t033
: bytes
| idx_t033_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t033
: bytes
| idx_t033_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t033
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t033
: bytes
| uk_t033_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t034
: bytes
| idx_t034_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t034
: bytes
| idx_t034_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t034
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t034
: bytes
| uk_t034_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t035
: bytes
| idx_t035_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t035
: bytes
| idx_t035_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t035
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t035
: bytes
| uk_t035_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t036
: bytes
| idx_t036_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t036
: bytes
| idx_t036_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t036
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t036
: bytes
| uk_t036_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t037
: bytes
| idx_t037_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t037
: bytes
| idx_t037_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t037
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t037
: bytes
| uk_t037_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t038
: bytes
| idx_t038_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t038
: bytes
| idx_t038_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t038
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t038
: bytes
| uk_t038_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t039
: bytes
| idx_t039_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t039
: bytes
| idx_t039_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t039
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t039
: bytes
| uk_t039_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t040
: bytes
| idx_t040_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t040
: bytes
| idx_t040_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t040
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t040
: bytes
| uk_t040_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t041
: bytes
| idx_t041_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t041
: bytes
| idx_t041_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t041
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t041
: bytes
| uk_t041_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t042
: bytes
| idx_t042_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t042
: bytes
| idx_t042_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t042
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t042
: bytes
| uk_t042_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t043
: bytes
| idx_t043_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t043
: bytes
| idx_t043_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t043
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t043
: bytes
| uk_t043_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t044
: bytes
| idx_t044_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t044
: bytes
| idx_t044_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t044
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t044
: bytes
| uk_t044_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t045
: bytes
| idx_t045_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t045
: bytes
| idx_t045_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t045
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t045
: bytes
| uk_t045_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t046
: bytes
| idx_t046_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t046
: bytes
| idx_t046_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t046
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t046
: bytes
| uk_t046_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t047
: bytes
| idx_t047_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t047
: bytes
| idx_t047_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t047
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t047
: bytes
| uk_t047_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t048
: bytes
| idx_t048_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t048
: bytes
| idx_t048_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t048
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t048
: bytes
| uk_t048_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t049
: bytes
| idx_t049_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t049
: bytes
| idx_t049_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t049
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t049
: bytes
| uk_t049_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t050
: bytes
| idx_t050_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t050
: bytes
| idx_t050_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t050
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t050
: bytes
| uk_t050_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t051
: bytes
| idx_t051_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t051
: bytes
| idx_t051_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t051
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t051
: bytes
| uk_t051_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t052
: bytes
| idx_t052_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t052
: bytes
| idx_t052_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t052
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t052
: bytes
| uk_t052_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t053
: bytes
| idx_t053_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t053
: bytes
| idx_t053_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t053
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t053
: bytes
| uk_t053_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t054
: bytes
| idx_t054_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t054
: bytes
| idx_t054_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t054
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t054
: bytes
| uk_t054_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t055
: bytes
| idx_t055_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t055
: bytes
| idx_t055_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t055
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t055
: bytes
| uk_t055_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t056
: bytes
| idx_t056_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t056
: bytes
| idx_t056_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t056
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t056
: bytes
| uk_t056_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t057
: bytes
| idx_t057_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t057
: bytes
| idx_t057_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t057
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t057
: bytes
| uk_t057_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t058
: bytes
| idx_t058_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t058
: bytes
| idx_t058_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t058
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t058
: bytes
| uk_t058_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t059
: bytes
| idx_t059_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t059
: bytes
| idx_t059_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t059
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t059
: bytes
| uk_t059_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t060
: bytes
| idx_t060_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t060
: bytes
| idx_t060_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t060
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t060
: bytes
| uk_t060_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t061
: bytes
| idx_t061_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t061
: bytes
| idx_t061_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t061
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t061
: bytes
| uk_t061_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t062
: bytes
| idx_t062_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t062
: bytes
| idx_t062_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t062
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t062
: bytes
| uk_t062_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t063
: bytes
| idx_t063_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t063
: bytes
| idx_t063_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t063
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t063
: bytes
| uk_t063_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t064
: bytes
| idx_t064_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t064
: bytes
| idx_t064_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t064
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t064
: bytes
| uk_t064_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t065
: bytes
| idx_t065_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t065
: bytes
| idx_t065_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t065
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t065
: bytes
| uk_t065_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t066
: bytes
| idx_t066_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t066
: bytes
| idx_t066_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t066
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t066
: bytes
| uk_t066_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t067
: bytes
| idx_t067_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t067
: bytes
| idx_t067_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t067
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t067
: bytes
| uk_t067_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t068
: bytes
| idx_t068_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t068
: bytes
| idx_t068_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t068
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t068
: bytes
| uk_t068_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t069
: bytes
| idx_t069_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t069
: bytes
| idx_t069_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t069
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t069
: bytes
| uk_t069_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t070
: bytes
| idx_t070_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t070
: bytes
| idx_t070_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t070
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t070
: bytes
| uk_t070_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t071
: bytes
| idx_t071_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t071
: bytes
| idx_t071_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t071
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t071
: bytes
| uk_t071_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t072
: bytes
| idx_t072_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t072
: bytes
| idx_t072_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t072
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t072
: bytes
| uk_t072_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t073
: bytes
| idx_t073_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t073
: bytes
| idx_t073_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t073
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t073
: bytes
| uk_t073_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t074
: bytes
| idx_t074_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t074
: bytes
| idx_t074_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t074
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t074
: bytes
| uk_t074_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t075
: bytes
| idx_t075_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t075
: bytes
| idx_t075_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t075
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t075
: bytes
| uk_t075_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t076
: bytes
| idx_t076_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t076
: bytes
| idx_t076_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t076
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t076
: bytes
| uk_t076_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t077
: bytes
| idx_t077_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t077
: bytes
| idx_t077_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t077
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t077
: bytes
| uk_t077_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t078
: bytes
| idx_t078_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t078
: bytes
| idx_t078_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t078
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t078
: bytes
| uk_t078_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t079
: bytes
| idx_t079_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t079
: bytes
| idx_t079_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t079
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t079
: bytes
| uk_t079_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t080
: bytes
| idx_t080_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t080
: bytes
| idx_t080_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t080
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t080
: bytes
| uk_t080_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t081
: bytes
| idx_t081_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t081
: bytes
| idx_t081_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t081
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t081
: bytes
| uk_t081_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t082
: bytes
| idx_t082_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t082
: bytes
| idx_t082_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t082
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t082
: bytes
| uk_t082_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t083
: bytes
| idx_t083_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t083
: bytes
| idx_t083_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t083
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t083
: bytes
| uk_t083_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t084
: bytes
| idx_t084_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t084
: bytes
| idx_t084_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t084
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t084
: bytes
| uk_t084_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t085
: bytes
| idx_t085_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t085
: bytes
| idx_t085_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t085
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t085
: bytes
| uk_t085_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t086
: bytes
| idx_t086_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t086
: bytes
| idx_t086_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t086
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t086
: bytes
| uk_t086_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t087
: bytes
| idx_t087_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t087
: bytes
| idx_t087_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t087
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t087
: bytes
| uk_t087_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t088
: bytes
| idx_t088_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t088
: bytes
| idx_t088_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t088
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t088
: bytes
| uk_t088_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t089
: bytes
| idx_t089_name
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t089
: bytes
| idx_t089_prev
: bytes
| prev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| t089
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| t089
: bytes
| uk_t089_code
: bytes
| code
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			IFNULL(COLUMN_NAME, ''),
| 			ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE,
| 			COLUMN_TYPE,
| 			IFNULL(CHARACTER_SET_NAME, ''),
| 			IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT),
| 			convert(GENERATION_EXPRESSION using BINARY),
| 			EXTRA,
| 			SRS_ID
| 		FROM information_schema.COLUMNS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| sdl_deep_scale_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR" "string"
< column "ORDINAL_POSITION" "UNSIGNED INT" "uint32"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "TEXT"
< column "IS_NULLABLE" "VARCHAR" "string"
< column "COLUMN_TYPE" "TEXT" "string"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR" "string"
< column "QUOTE(COLUMN_COMMENT)" "TEXT"
< column "convert(GENERATION_EXPRESSION using BINARY)" "BLOB" "[]uint8"
< column "EXTRA" "VARCHAR"
< column "SRS_ID" "UNSIGNED INT"
< row
: bytes
| t000
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t000
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t000
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t000
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t000
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t000
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t000
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t000
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t000
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t000
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t001
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t001
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t001
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t001
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t001
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t001
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t001
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t001
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t001
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t001
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t001
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t002
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t002
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t002
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t002
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t002
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t002
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t002
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t002
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t002
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t002
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t002
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t003
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t003
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t003
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t003
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t003
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t003
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t003
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t003
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t003
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t003
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t003
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t004
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t004
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t004
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t004
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t004
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t004
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t004
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t004
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t004
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t004
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t004
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t005
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t005
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t005
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t005
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t005
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t005
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t005
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t005
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t005
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t005
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t005
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t006
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t006
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t006
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t006
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t006
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t006
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t006
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t006
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t006
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t006
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t006
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t007
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t007
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t007
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t007
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t007
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t007
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t007
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t007
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t007
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t007
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t007
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t008
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t008
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t008
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t008
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t008
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t008
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t008
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t008
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t008
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t008
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t008
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t009
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t009
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t009
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t009
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t009
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t009
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t009
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t009
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t009
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t009
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t009
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t010
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t010
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t010
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t010
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t010
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t010
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t010
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t010
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t010
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t010
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t010
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t011
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t011
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t011
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t011
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t011
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t011
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t011
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t011
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t011
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t011
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t011
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t012
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t012
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t012
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t012
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t012
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t012
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t012
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t012
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t012
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t012
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t012
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t013
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t013
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t013
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t013
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t013
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t013
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t013
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t013
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t013
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t013
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t013
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t014
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t014
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t014
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t014
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t014
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t014
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t014
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t014
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t014
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t014
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t014
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t015
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t015
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t015
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t015
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t015
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t015
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t015
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t015
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t015
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t015
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t015
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t016
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t016
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t016
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t016
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t016
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t016
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t016
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t016
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t016
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t016
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t016
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t017
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t017
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t017
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t017
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t017
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t017
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t017
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t017
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t017
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t017
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t017
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t018
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t018
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t018
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t018
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t018
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t018
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t018
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t018
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t018
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t018
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t018
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t019
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t019
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t019
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t019
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t019
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t019
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t019
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t019
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t019
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t019
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t019
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t020
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t020
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t020
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t020
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t020
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t020
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t020
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t020
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t020
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t020
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t020
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t021
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t021
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t021
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t021
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t021
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t021
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t021
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t021
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t021
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t021
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t021
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t022
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t022
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t022
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t022
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t022
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t022
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t022
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t022
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t022
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t022
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t022
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t023
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t023
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t023
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t023
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t023
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t023
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t023
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t023
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t023
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t023
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t023
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t024
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t024
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t024
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t024
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t024
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t024
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t024
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t024
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t024
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t024
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t024
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t025
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t025
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t025
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t025
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t025
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t025
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t025
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t025
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t025
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t025
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t025
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t026
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t026
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t026
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t026
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t026
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t026
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t026
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t026
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t026
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t026
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t026
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t027
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t027
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t027
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t027
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t027
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t027
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t027
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t027
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t027
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t027
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t027
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t028
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t028
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t028
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t028
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t028
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t028
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t028
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t028
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t028
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t028
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t028
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t029
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t029
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t029
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t029
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t029
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t029
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t029
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t029
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t029
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t029
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t029
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t030
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t030
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t030
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t030
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t030
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t030
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t030
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t030
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t030
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t030
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t030
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t031
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t031
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t031
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t031
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t031
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t031
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t031
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t031
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t031
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t031
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t031
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t032
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t032
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t032
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t032
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t032
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t032
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t032
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t032
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t032
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t032
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t032
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t033
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t033
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t033
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t033
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t033
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t033
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t033
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t033
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t033
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t033
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t033
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t034
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t034
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t034
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t034
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t034
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t034
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t034
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t034
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t034
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t034
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t034
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t035
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t035
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t035
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t035
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t035
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t035
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t035
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t035
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t035
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t035
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t035
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t036
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t036
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t036
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t036
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t036
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t036
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t036
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t036
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t036
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t036
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t036
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t037
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t037
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t037
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t037
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t037
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t037
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t037
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t037
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t037
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t037
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t037
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t038
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t038
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t038
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t038
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t038
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t038
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t038
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t038
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t038
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t038
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t038
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t039
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t039
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t039
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t039
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t039
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t039
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t039
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t039
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t039
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t039
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t039
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t040
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t040
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t040
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t040
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t040
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t040
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t040
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t040
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t040
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t040
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t040
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t041
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t041
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t041
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t041
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t041
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t041
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t041
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t041
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t041
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t041
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t041
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t042
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t042
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t042
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t042
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t042
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t042
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t042
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t042
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t042
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t042
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t042
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t043
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t043
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t043
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t043
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t043
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t043
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t043
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t043
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t043
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t043
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t043
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t044
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t044
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t044
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t044
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t044
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t044
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t044
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t044
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t044
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t044
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t044
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t045
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t045
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t045
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t045
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t045
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t045
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t045
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t045
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t045
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t045
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t045
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t046
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t046
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t046
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t046
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t046
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t046
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t046
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t046
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t046
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t046
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t046
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t047
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t047
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t047
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t047
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t047
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t047
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t047
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t047
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t047
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t047
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t047
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t048
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t048
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t048
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t048
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t048
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t048
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t048
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t048
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t048
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t048
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t048
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t049
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t049
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t049
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t049
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t049
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t049
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t049
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t049
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t049
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t049
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t049
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t050
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t050
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t050
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t050
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t050
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t050
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t050
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t050
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t050
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t050
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t050
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t051
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t051
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t051
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t051
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t051
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t051
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t051
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t051
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t051
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t051
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t051
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t052
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t052
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t052
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t052
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t052
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t052
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t052
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t052
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t052
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t052
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t052
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t053
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t053
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t053
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t053
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t053
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t053
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t053
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t053
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t053
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t053
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t053
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t054
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t054
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t054
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t054
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t054
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t054
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t054
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t054
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t054
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t054
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t054
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t055
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t055
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t055
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t055
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t055
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t055
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t055
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t055
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t055
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t055
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t055
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t056
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t056
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t056
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t056
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t056
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t056
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t056
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t056
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t056
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t056
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t056
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t057
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t057
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t057
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t057
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t057
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t057
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t057
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t057
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t057
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t057
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t057
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t058
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t058
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t058
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t058
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t058
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t058
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t058
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t058
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t058
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t058
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t058
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t059
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t059
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t059
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t059
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t059
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t059
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t059
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t059
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t059
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t059
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t059
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t060
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t060
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t060
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t060
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t060
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t060
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t060
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t060
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t060
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t060
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t060
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t061
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t061
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t061
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t061
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t061
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t061
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t061
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t061
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t061
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t061
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t061
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t062
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t062
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t062
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t062
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t062
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t062
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t062
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t062
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t062
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t062
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t062
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t063
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t063
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t063
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t063
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t063
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t063
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t063
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t063
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t063
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t063
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t063
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t064
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t064
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t064
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t064
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t064
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t064
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t064
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t064
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t064
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t064
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t064
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t065
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t065
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t065
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t065
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t065
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t065
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t065
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t065
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t065
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t065
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t065
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t066
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t066
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t066
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t066
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t066
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t066
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t066
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t066
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t066
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t066
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t066
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t067
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t067
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t067
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t067
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t067
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t067
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t067
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t067
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t067
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t067
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t067
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t068
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t068
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t068
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t068
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t068
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t068
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t068
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t068
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t068
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t068
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t068
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t069
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t069
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t069
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t069
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t069
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t069
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t069
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t069
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t069
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t069
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t069
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t070
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t070
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t070
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t070
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t070
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t070
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t070
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t070
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t070
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t070
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t070
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t071
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t071
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t071
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t071
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t071
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t071
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t071
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t071
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t071
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t071
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t071
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t072
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t072
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t072
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t072
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t072
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t072
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t072
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t072
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t072
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t072
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t072
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t073
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t073
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t073
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t073
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t073
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t073
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t073
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t073
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t073
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t073
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t073
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t074
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t074
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t074
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t074
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t074
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t074
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t074
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t074
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t074
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t074
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t074
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t075
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t075
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t075
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t075
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t075
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t075
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t075
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t075
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t075
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t075
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t075
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t076
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t076
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t076
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t076
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t076
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t076
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t076
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t076
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t076
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t076
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t076
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t077
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t077
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t077
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t077
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t077
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t077
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t077
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t077
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t077
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t077
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t077
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t078
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t078
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t078
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t078
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t078
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t078
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t078
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t078
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t078
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t078
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t078
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t079
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t079
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t079
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t079
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t079
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t079
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t079
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t079
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t079
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t079
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t079
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t080
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t080
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t080
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t080
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t080
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t080
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t080
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t080
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t080
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t080
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t080
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t081
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t081
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t081
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t081
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t081
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t081
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t081
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t081
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t081
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t081
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t081
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t082
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t082
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t082
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t082
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t082
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t082
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t082
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t082
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t082
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t082
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t082
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t083
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t083
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t083
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t083
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t083
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t083
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t083
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t083
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t083
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t083
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t083
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t084
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t084
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t084
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t084
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t084
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t084
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t084
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t084
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t084
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t084
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t084
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t085
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t085
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t085
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t085
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t085
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t085
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t085
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t085
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t085
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t085
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t085
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t086
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t086
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t086
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t086
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t086
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t086
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t086
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t086
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t086
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t086
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t086
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t087
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t087
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t087
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t087
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t087
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t087
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t087
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t087
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t087
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t087
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t087
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t088
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t088
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t088
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t088
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t088
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t088
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t088
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t088
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t088
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t088
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t088
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t089
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| t089
: bytes
| code
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t089
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(120)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t089
: bytes
| description
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t089
: bytes
| amount
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t089
: bytes
| qty
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t089
: bytes
| active
: int64 7
: bytes
| '1'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t089
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| t089
: bytes
| updated_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| t089
: bytes
| status
: int64 10
: bytes
| 'new'
: bytes
| NO
: bytes
| enum('new','open','closed')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| t089
: bytes
| prev_id
: int64 11
: null
: bytes
| YES
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
> query
|
| 			SELECT
| 				tc.TABLE_NAME,
| 				cc.CONSTRAINT_NAME,
| 				cc.CHECK_CLAUSE,
| 				tc.ENFORCED
| 			FROM information_schema.CHECK_CONSTRAINTS cc
| 				JOIN information_schema.TABLE_CONSTRAINTS tc
| 					ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 					AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 			WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ?
| 			ORDER BY tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		
: string
| sdl_deep_scale_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< column "ENFORCED" "VARCHAR" "string"
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| sdl_deep_scale_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT"
> query
|
| 	SELECT
| 		TRIGGER_NAME,
| 		EVENT_OBJECT_TABLE,
| 		EVENT_MANIPULATION,
| 		ACTION_TIMING,
| 		ACTION_STATEMENT,
| 		SQL_MODE,
| 		CHARACTER_SET_CLIENT,
| 		COLLATION_CONNECTION
| 	FROM INFORMATION_SCHEMA.TRIGGERS
| 	WHERE EVENT_OBJECT_SCHEMA = ?
| 	ORDER BY EVENT_OBJECT_TABLE ASC, EVENT_MANIPULATION ASC, ACTION_TIMING ASC, ACTION_ORDER ASC;
| 	
: string
| sdl_deep_scale_1
< result
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "EVENT_OBJECT_TABLE" "VARCHAR" "string"
< column "EVENT_MANIPULATION" "ENUM" "string"
< column "ACTION_TIMING" "ENUM" "string"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
> query
|
| 	SELECT
| 		EVENT_NAME,
| 		TIME_ZONE,
| 		SQL_MODE,
| 		CHARACTER_SET_CLIENT,
| 		COLLATION_CONNECTION
| 	FROM INFORMATION_SCHEMA.EVENTS
| 	WHERE EVENT_SCHEMA = ?
| 	ORDER BY EVENT_NAME ASC;
| 	
: string
| sdl_deep_scale_1
< result
< column "EVENT_NAME" "VARCHAR" "string"
< column "TIME_ZONE" "VARCHAR" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			CONSTRAINT_NAME,
| 			REFERENCED_TABLE_NAME,
| 			DELETE_RULE,
| 			UPDATE_RULE,
| 			MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ?;
| 	
: string
| sdl_deep_scale_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
< row
: bytes
| t001
: bytes
| fk_t001_prev
: bytes
| t000
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t002
: bytes
| fk_t002_prev
: bytes
| t001
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t003
: bytes
| fk_t003_prev
: bytes
| t002
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t004
: bytes
| fk_t004_prev
: bytes
| t003
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t005
: bytes
| fk_t005_prev
: bytes
| t004
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t006
: bytes
| fk_t006_prev
: bytes
| t005
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t007
: bytes
| fk_t007_prev
: bytes
| t006
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t008
: bytes
| fk_t008_prev
: bytes
| t007
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t009
: bytes
| fk_t009_prev
: bytes
| t008
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t010
: bytes
| fk_t010_prev
: bytes
| t009
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t011
: bytes
| fk_t011_prev
: bytes
| t010
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t012
: bytes
| fk_t012_prev
: bytes
| t011
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t013
: bytes
| fk_t013_prev
: bytes
| t012
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t014
: bytes
| fk_t014_prev
: bytes
| t013
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t015
: bytes
| fk_t015_prev
: bytes
| t014
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t016
: bytes
| fk_t016_prev
: bytes
| t015
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t017
: bytes
| fk_t017_prev
: bytes
| t016
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t018
: bytes
| fk_t018_prev
: bytes
| t017
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t019
: bytes
| fk_t019_prev
: bytes
| t018
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t020
: bytes
| fk_t020_prev
: bytes
| t019
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t021
: bytes
| fk_t021_prev
: bytes
| t020
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t022
: bytes
| fk_t022_prev
: bytes
| t021
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t023
: bytes
| fk_t023_prev
: bytes
| t022
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t024
: bytes
| fk_t024_prev
: bytes
| t023
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t025
: bytes
| fk_t025_prev
: bytes
| t024
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t026
: bytes
| fk_t026_prev
: bytes
| t025
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t027
: bytes
| fk_t027_prev
: bytes
| t026
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t028
: bytes
| fk_t028_prev
: bytes
| t027
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t029
: bytes
| fk_t029_prev
: bytes
| t028
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t030
: bytes
| fk_t030_prev
: bytes
| t029
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t031
: bytes
| fk_t031_prev
: bytes
| t030
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t032
: bytes
| fk_t032_prev
: bytes
| t031
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t033
: bytes
| fk_t033_prev
: bytes
| t032
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t034
: bytes
| fk_t034_prev
: bytes
| t033
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t035
: bytes
| fk_t035_prev
: bytes
| t034
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t036
: bytes
| fk_t036_prev
: bytes
| t035
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t037
: bytes
| fk_t037_prev
: bytes
| t036
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t038
: bytes
| fk_t038_prev
: bytes
| t037
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t039
: bytes
| fk_t039_prev
: bytes
| t038
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t040
: bytes
| fk_t040_prev
: bytes
| t039
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t041
: bytes
| fk_t041_prev
: bytes
| t040
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t042
: bytes
| fk_t042_prev
: bytes
| t041
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t043
: bytes
| fk_t043_prev
: bytes
| t042
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t044
: bytes
| fk_t044_prev
: bytes
| t043
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t045
: bytes
| fk_t045_prev
: bytes
| t044
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t046
: bytes
| fk_t046_prev
: bytes
| t045
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t047
: bytes
| fk_t047_prev
: bytes
| t046
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t048
: bytes
| fk_t048_prev
: bytes
| t047
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t049
: bytes
| fk_t049_prev
: bytes
| t048
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t050
: bytes
| fk_t050_prev
: bytes
| t049
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t051
: bytes
| fk_t051_prev
: bytes
| t050
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t052
: bytes
| fk_t052_prev
: bytes
| t051
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t053
: bytes
| fk_t053_prev
: bytes
| t052
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t054
: bytes
| fk_t054_prev
: bytes
| t053
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t055
: bytes
| fk_t055_prev
: bytes
| t054
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t056
: bytes
| fk_t056_prev
: bytes
| t055
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t057
: bytes
| fk_t057_prev
: bytes
| t056
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t058
: bytes
| fk_t058_prev
: bytes
| t057
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t059
: bytes
| fk_t059_prev
: bytes
| t058
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t060
: bytes
| fk_t060_prev
: bytes
| t059
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t061
: bytes
| fk_t061_prev
: bytes
| t060
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t062
: bytes
| fk_t062_prev
: bytes
| t061
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t063
: bytes
| fk_t063_prev
: bytes
| t062
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t064
: bytes
| fk_t064_prev
: bytes
| t063
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t065
: bytes
| fk_t065_prev
: bytes
| t064
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t066
: bytes
| fk_t066_prev
: bytes
| t065
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t067
: bytes
| fk_t067_prev
: bytes
| t066
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t068
: bytes
| fk_t068_prev
: bytes
| t067
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t069
: bytes
| fk_t069_prev
: bytes
| t068
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t070
: bytes
| fk_t070_prev
: bytes
| t069
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t071
: bytes
| fk_t071_prev
: bytes
| t070
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t072
: bytes
| fk_t072_prev
: bytes
| t071
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t073
: bytes
| fk_t073_prev
: bytes
| t072
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t074
: bytes
| fk_t074_prev
: bytes
| t073
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t075
: bytes
| fk_t075_prev
: bytes
| t074
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t076
: bytes
| fk_t076_prev
: bytes
| t075
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t077
: bytes
| fk_t077_prev
: bytes
| t076
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t078
: bytes
| fk_t078_prev
: bytes
| t077
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t079
: bytes
| fk_t079_prev
: bytes
| t078
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t080
: bytes
| fk_t080_prev
: bytes
| t079
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t081
: bytes
| fk_t081_prev
: bytes
| t080
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t082
: bytes
| fk_t082_prev
: bytes
| t081
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t083
: bytes
| fk_t083_prev
: bytes
| t082
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t084
: bytes
| fk_t084_prev
: bytes
| t083
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t085
: bytes
| fk_t085_prev
: bytes
| t084
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t086
: bytes
| fk_t086_prev
: bytes
| t085
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t087
: bytes
| fk_t087_prev
: bytes
| t086
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t088
: bytes
| fk_t088_prev
: bytes
| t087
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| t089
: bytes
| fk_t089_prev
: bytes
| t088
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
> query
|
| 		SELECT
| 			TABLE_SCHEMA,
| 			TABLE_NAME,
| 			CONSTRAINT_NAME,
| 			COLUMN_NAME,
| 			REFERENCED_TABLE_SCHEMA,
| 			REFERENCED_COLUMN_NAME
| 		FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
| 		WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION;
| 	
: string
| sdl_deep_scale_1
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| sdl_deep_scale_1
: bytes
| t001
: bytes
| fk_t001_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t002
: bytes
| fk_t002_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t003
: bytes
| fk_t003_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t004
: bytes
| fk_t004_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t005
: bytes
| fk_t005_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t006
: bytes
| fk_t006_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t007
: bytes
| fk_t007_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t008
: bytes
| fk_t008_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t009
: bytes
| fk_t009_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t010
: bytes
| fk_t010_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t011
: bytes
| fk_t011_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t012
: bytes
| fk_t012_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t013
: bytes
| fk_t013_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t014
: bytes
| fk_t014_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t015
: bytes
| fk_t015_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t016
: bytes
| fk_t016_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t017
: bytes
| fk_t017_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t018
: bytes
| fk_t018_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t019
: bytes
| fk_t019_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t020
: bytes
| fk_t020_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t021
: bytes
| fk_t021_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t022
: bytes
| fk_t022_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t023
: bytes
| fk_t023_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t024
: bytes
| fk_t024_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t025
: bytes
| fk_t025_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t026
: bytes
| fk_t026_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t027
: bytes
| fk_t027_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t028
: bytes
| fk_t028_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t029
: bytes
| fk_t029_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t030
: bytes
| fk_t030_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t031
: bytes
| fk_t031_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t032
: bytes
| fk_t032_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t033
: bytes
| fk_t033_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t034
: bytes
| fk_t034_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t035
: bytes
| fk_t035_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t036
: bytes
| fk_t036_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t037
: bytes
| fk_t037_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t038
: bytes
| fk_t038_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t039
: bytes
| fk_t039_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t040
: bytes
| fk_t040_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t041
: bytes
| fk_t041_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t042
: bytes
| fk_t042_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t043
: bytes
| fk_t043_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t044
: bytes
| fk_t044_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t045
: bytes
| fk_t045_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t046
: bytes
| fk_t046_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t047
: bytes
| fk_t047_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t048
: bytes
| fk_t048_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t049
: bytes
| fk_t049_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t050
: bytes
| fk_t050_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t051
: bytes
| fk_t051_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t052
: bytes
| fk_t052_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t053
: bytes
| fk_t053_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t054
: bytes
| fk_t054_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t055
: bytes
| fk_t055_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t056
: bytes
| fk_t056_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t057
: bytes
| fk_t057_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t058
: bytes
| fk_t058_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t059
: bytes
| fk_t059_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t060
: bytes
| fk_t060_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t061
: bytes
| fk_t061_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t062
: bytes
| fk_t062_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t063
: bytes
| fk_t063_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t064
: bytes
| fk_t064_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t065
: bytes
| fk_t065_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t066
: bytes
| fk_t066_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t067
: bytes
| fk_t067_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t068
: bytes
| fk_t068_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t069
: bytes
| fk_t069_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t070
: bytes
| fk_t070_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t071
: bytes
| fk_t071_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t072
: bytes
| fk_t072_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t073
: bytes
| fk_t073_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t074
: bytes
| fk_t074_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t075
: bytes
| fk_t075_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t076
: bytes
| fk_t076_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t077
: bytes
| fk_t077_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t078
: bytes
| fk_t078_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t079
: bytes
| fk_t079_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t080
: bytes
| fk_t080_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t081
: bytes
| fk_t081_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t082
: bytes
| fk_t082_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t083
: bytes
| fk_t083_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t084
: bytes
| fk_t084_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t085
: bytes
| fk_t085_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t086
: bytes
| fk_t086_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t087
: bytes
| fk_t087_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t088
: bytes
| fk_t088_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
< row
: bytes
| sdl_deep_scale_1
: bytes
| t089
: bytes
| fk_t089_prev
: bytes
| prev_id
: bytes
| sdl_deep_scale_1
: bytes
| id
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			PARTITION_NAME,
| 			SUBPARTITION_NAME,
| 			PARTITION_METHOD,
| 			SUBPARTITION_METHOD,
| 			PARTITION_EXPRESSION,
| 			SUBPARTITION_EXPRESSION,
| 			PARTITION_DESCRIPTION,
| 			PARTITION_COMMENT
| 		FROM INFORMATION_SCHEMA.PARTITIONS
| 		WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME ASC, PARTITION_ORDINAL_POSITION ASC, SUBPARTITION_ORDINAL_POSITION ASC;
| 	
: string
| sdl_deep_scale_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "SUBPARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "VARCHAR"
< column "SUBPARTITION_EXPRESSION" "VARCHAR"
< column "PARTITION_DESCRIPTION" "TEXT"
< column "PARTITION_COMMENT" "TEXT" "string"
> query
|
| 		SELECT
| 			ROUTINE_NAME,
| 			ROUTINE_TYPE,
| 			SQL_MODE,
| 			CHARACTER_SET_CLIENT,
| 			COLLATION_CONNECTION,
| 			DATABASE_COLLATION,
| 			IFNULL(ROUTINE_COMMENT, '')
| 		FROM
| 			INFORMATION_SCHEMA.ROUTINES
| 		WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE IN ('FUNCTION', 'PROCEDURE')
| 		ORDER BY ROUTINE_TYPE, ROUTINE_NAME;
| 	
: string
| sdl_deep_scale_1
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "ENUM" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT" "string"
> query
|
| 		SELECT
| 			TABLES.TABLE_NAME,
| 			TABLES.TABLE_TYPE,
| 			IFNULL(TABLES.ENGINE, ''),
| 			IFNULL(TABLES.TABLE_COLLATION, ''),
| 			IFNULL(TABLES.TABLE_ROWS, 0),
| 			IFNULL(TABLES.DATA_LENGTH, 0),
| 			IFNULL(TABLES.INDEX_LENGTH, 0),
| 			IFNULL(TABLES.DATA_FREE, 0),
| 			IFNULL(TABLES.CREATE_OPTIONS, ''),
| 			QUOTE(IFNULL(TABLES.TABLE_COMMENT, '')),
| 			IFNULL(CCSA.CHARACTER_SET_NAME, '')
| 		FROM information_schema.TABLES TABLES
| 		LEFT JOIN information_schema.COLLATION_CHARACTER_SET_APPLICABILITY CCSA
| 		ON TABLES.TABLE_COLLATION = CCSA.COLLATION_NAME
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| sdl_deep_scale_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "TABLE_TYPE" "ENUM" "string"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLES.TABLE_COMMENT, ''))" "TEXT"
< column "IFNULL(CCSA.CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< row
: bytes
| t000
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t001
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t002
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t003
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t004
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t005
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t006
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t007
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t008
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t009
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t010
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t011
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t012
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t013
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t014
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t015
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t016
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t017
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t018
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t019
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t020
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t021
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t022
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t023
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t024
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t025
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t026
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t027
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t028
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t029
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t030
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t031
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t032
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t033
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t034
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t035
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t036
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t037
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t038
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t039
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t040
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t041
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t042
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t043
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t044
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t045
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t046
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t047
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t048
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t049
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t050
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t051
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t052
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t053
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t054
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t055
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t056
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t057
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t058
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t059
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t060
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t061
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t062
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t063
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t064
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t065
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t066
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t067
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t068
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t069
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t070
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t071
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t072
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t073
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t074
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t075
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t076
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t077
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t078
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t079
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t080
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t081
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t082
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t083
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t084
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t085
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t086
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t087
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t088
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| t089
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 49152
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
> query
|
| 		SELECT t.NAME
| 		FROM information_schema.INNODB_TABLES t
| 		LEFT JOIN information_schema.INNODB_TABLESPACES ts ON t.SPACE = ts.SPACE
| 		WHERE t.NAME LIKE CONCAT(?, '/%') AND (ts.SPACE_TYPE IS NULL OR ts.SPACE_TYPE <> 'Single')
: string
| sdl_deep_scale_1
< result
< column "NAME" "VARCHAR" "string"
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| sdl_deep_scale_1
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
