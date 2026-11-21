-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestSDLRealWorldIdempotence default#1
> begin
< ok
> exec
| SET FOREIGN_KEY_CHECKS=0;
< ok 0 0
> exec
|
|
| -- Sakila Sample Database Schema
| -- Version 0.8
|
| -- Copyright (c) 2006, MySQL AB
| -- All rights reserved.
|
| -- Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
|
| --  * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
| --  * Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
| --  * Neither the name of MySQL AB nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
|
| -- THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
|
|
|
|
| --
| -- Table structure for table `actor`
| --
|
| CREATE TABLE actor (
|   actor_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
|   first_name VARCHAR(45) NOT NULL,
|   last_name VARCHAR(45) NOT NULL,
|   last_update TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY  (actor_id),
|   KEY idx_actor_last_name (last_name)
| )ENGINE=InnoDB DEFAULT CHARSET=utf8;
< ok 0 0
> exec
|
|
| --
| -- Table structure for table `address`
| --
|
| CREATE TABLE address (
|   address_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
|   address VARCHAR(50) NOT NULL,
|   address2 VARCHAR(50) DEFAULT NULL,
|   district VARCHAR(20) NOT NULL,
|   city_id INT UNSIGNED NOT NULL,
|   postal_code VARCHAR(10) DEFAULT NULL,
|   phone VARCHAR(20) NOT NULL,
|   last_update TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY  (address_id),
|   KEY idx_fk_city_id (city_id),
|   CONSTRAINT `fk_address_city` FOREIGN KEY (city_id) REFERENCES city (city_id) ON DELETE RESTRICT ON UPDATE CASCADE
| )ENGINE=InnoDB DEFAULT CHARSET=utf8;
< ok 0 0
> exec
|
|
| --
| -- Table structure for table `category`
| --
|
| CREATE TABLE category (
|   category_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
|   name VARCHAR(25) NOT NULL,
|   last_update TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY  (category_id)
| )ENGINE=InnoDB DEFAULT CHARSET=utf8;
< ok 0 0
> exec
|
|
| --
| -- Table structure for table `city`
| --
|
| CREATE TABLE city (
|   city_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
|   city VARCHAR(50) NOT NULL,
|   country_id INT UNSIGNED NOT NULL,
|   last_update TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY  (city_id),
|   KEY idx_fk_country_id (country_id),
|   CONSTRAINT `fk_city_country` FOREIGN KEY (country_id) REFERENCES country (country_id) ON DELETE RESTRICT ON UPDATE CASCADE
| )ENGINE=InnoDB DEFAULT CHARSET=utf8;
< ok 0 0
> exec
|
|
| --
| -- Table structure for table `country`
| --
|
| CREATE TABLE country (
|   country_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
|   country VARCHAR(50) NOT NULL,
|   last_update TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY  (country_id)
| )ENGINE=InnoDB DEFAULT CHARSET=utf8;
< ok 0 0
> exec
|
|
| --
| -- Table structure for table `customer`
| --
|
| CREATE TABLE customer (
|   customer_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
|   store_id INT UNSIGNED NOT NULL,
|   first_name VARCHAR(45) NOT NULL,
|   last_name VARCHAR(45) NOT NULL,
|   email VARCHAR(50) DEFAULT NULL,
|   address_id INT UNSIGNED NOT NULL,
|   active BOOLEAN NOT NULL DEFAULT TRUE,
|   create_date DATETIME NOT NULL,
|   last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY  (customer_id),
|   KEY idx_fk_store_id (store_id),
|   KEY idx_fk_address_id (address_id),
|   KEY idx_last_name (last_name),
|   CONSTRAINT fk_customer_address FOREIGN KEY (address_id) REFERENCES address (address_id) ON DELETE RESTRICT ON UPDATE CASCADE,
|   CONSTRAINT fk_customer_store FOREIGN KEY (store_id) REFERENCES store (store_id) ON DELETE RESTRICT ON UPDATE CASCADE
| )ENGINE=InnoDB DEFAULT CHARSET=utf8;
< ok 0 0
> exec
|
|
| --
| -- Table structure for table `film`
| --
|
| CREATE TABLE film (
|   film_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
|   title VARCHAR(255) NOT NULL,
|   description TEXT DEFAULT NULL,
|   release_year YEAR DEFAULT NULL,
|   language_id INT UNSIGNED NOT NULL,
|   original_language_id INT UNSIGNED DEFAULT NULL,
|   rental_duration TINYINT UNSIGNED NOT NULL DEFAULT 3,
|   rental_rate DECIMAL(4,2) NOT NULL DEFAULT 4.99,
|   length SMALLINT UNSIGNED DEFAULT NULL,
|   replacement_cost DECIMAL(5,2) NOT NULL DEFAULT 19.99,
|   rating ENUM('G','PG','PG-13','R','NC-17') DEFAULT 'G',
|   special_features SET('Trailers','Commentaries','Deleted Scenes','Behind the Scenes') DEFAULT NULL,
|   last_update TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY  (film_id),
|   KEY idx_title (title),
|   KEY idx_fk_language_id (language_id),
|   KEY idx_fk_original_language_id (original_language_id),
|   CONSTRAINT fk_film_language FOREIGN KEY (language_id) REFERENCES language (language_id) ON DELETE RESTRICT ON UPDATE CASCADE,
|   CONSTRAINT fk_film_language_original FOREIGN KEY (original_language_id) REFERENCES language (language_id) ON DELETE RESTRICT ON UPDATE CASCADE
| )ENGINE=InnoDB DEFAULT CHARSET=utf8;
< ok 0 0
> exec
|
|
| --
| -- Table structure for table `film_actor`
| --
|
| CREATE TABLE film_actor (
|   actor_id INT UNSIGNED NOT NULL,
|   film_id INT UNSIGNED NOT NULL,
|   last_update TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY  (actor_id,film_id),
|   KEY idx_fk_film_id (`film_id`),
|   CONSTRAINT fk_film_actor_actor FOREIGN KEY (actor_id) REFERENCES actor (actor_id) ON DELETE RESTRICT ON UPDATE CASCADE,
|   CONSTRAINT fk_film_actor_film FOREIGN KEY (film_id) REFERENCES film (film_id) ON DELETE RESTRICT ON UPDATE CASCADE
| )ENGINE=InnoDB DEFAULT CHARSET=utf8;
< ok 0 0
> exec
|
|
| --
| -- Table structure for table `film_category`
| --
|
| CREATE TABLE film_category (
|   film_id INT UNSIGNED NOT NULL,
|   category_id INT UNSIGNED NOT NULL,
|   last_update TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY (film_id, category_id),
|   CONSTRAINT fk_film_category_film FOREIGN KEY (film_id) REFERENCES film (film_id) ON DELETE RESTRICT ON UPDATE CASCADE,
|   CONSTRAINT fk_film_category_category FOREIGN KEY (category_id) REFERENCES category (category_id) ON DELETE RESTRICT ON UPDATE CASCADE
| )ENGINE=InnoDB DEFAULT CHARSET=utf8;
< ok 0 0
> exec
|
|
| --
| -- Table structure for table `film_text`
| --
|
| CREATE TABLE film_text (
|   film_id INT NOT NULL,
|   title VARCHAR(255) NOT NULL,
|   description TEXT,
|   PRIMARY KEY  (film_id),
|   FULLTEXT KEY idx_title_description (title,description)
| )ENGINE=InnoDB DEFAULT CHARSET=utf8;
< ok 0 0
> exec
|
|
| --
| -- Triggers for loading film_text from film
| --
|
| CREATE TRIGGER `ins_film` AFTER INSERT ON `film` FOR EACH ROW BEGIN
|     INSERT INTO film_text (film_id, title, description)
|         VALUES (new.film_id, new.title, new.description);
|   END;
< ok 0 0
> exec
|
|
|
| CREATE TRIGGER `upd_film` AFTER UPDATE ON `film` FOR EACH ROW BEGIN
|     IF (old.title != new.title) or (old.description != new.description)
|     THEN
|         UPDATE film_text
|             SET title=new.title,
|                 description=new.description,
|                 film_id=new.film_id
|         WHERE film_id=old.film_id;
|     END IF;
|   END;
< ok 0 0
> exec
|
|
|
| CREATE TRIGGER `del_film` AFTER DELETE ON `film` FOR EACH ROW BEGIN
|     DELETE FROM film_text WHERE film_id = old.film_id;
|   END;
< ok 0 0
> exec
|
|
|
| --
| -- Table structure for table `inventory`
| --
|
| CREATE TABLE inventory (
|   inventory_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
|   film_id INT UNSIGNED NOT NULL,
|   store_id INT UNSIGNED NOT NULL,
|   last_update TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY  (inventory_id),
|   KEY idx_fk_film_id (film_id),
|   KEY idx_store_id_film_id (store_id,film_id),
|   CONSTRAINT fk_inventory_store FOREIGN KEY (store_id) REFERENCES store (store_id) ON DELETE RESTRICT ON UPDATE CASCADE,
|   CONSTRAINT fk_inventory_film FOREIGN KEY (film_id) REFERENCES film (film_id) ON DELETE RESTRICT ON UPDATE CASCADE
| )ENGINE=InnoDB DEFAULT CHARSET=utf8;
< ok 0 0
> exec
|
|
| --
| -- Table structure for table `language`
| --
|
| CREATE TABLE language (
|   language_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
|   name CHAR(20) NOT NULL,
|   last_update TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY (language_id)
| )ENGINE=InnoDB DEFAULT CHARSET=utf8;
< ok 0 0
> exec
|
|
| --
| -- Table structure for table `payment`
| --
|
| CREATE TABLE payment (
|   payment_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
|   customer_id INT UNSIGNED NOT NULL,
|   staff_id INT UNSIGNED NOT NULL,
|   rental_id INT DEFAULT NULL,
|   amount DECIMAL(5,2) NOT NULL,
|   payment_date DATETIME NOT NULL,
|   last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY  (payment_id),
|   KEY idx_fk_staff_id (staff_id),
|   KEY idx_fk_customer_id (customer_id),
|   CONSTRAINT fk_payment_rental FOREIGN KEY (rental_id) REFERENCES rental (rental_id) ON DELETE SET NULL ON UPDATE CASCADE,
|   CONSTRAINT fk_payment_customer FOREIGN KEY (customer_id) REFERENCES customer (customer_id) ON DELETE RESTRICT ON UPDATE CASCADE,
|   CONSTRAINT fk_payment_staff FOREIGN KEY (staff_id) REFERENCES staff (staff_id) ON DELETE RESTRICT ON UPDATE CASCADE
| )ENGINE=InnoDB DEFAULT CHARSET=utf8;
< ok 0 0
> exec
|
|
|
| --
| -- Table structure for table `rental`
| --
|
| CREATE TABLE rental (
|   rental_id INT NOT NULL AUTO_INCREMENT,
|   rental_date DATETIME NOT NULL,
|   inventory_id INT UNSIGNED NOT NULL,
|   customer_id INT UNSIGNED NOT NULL,
|   return_date DATETIME DEFAULT NULL,
|   staff_id INT UNSIGNED NOT NULL,
|   last_update TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY (rental_id),
|   UNIQUE KEY  (rental_date,inventory_id,customer_id),
|   KEY idx_fk_inventory_id (inventory_id),
|   KEY idx_fk_customer_id (customer_id),
|   KEY idx_fk_staff_id (staff_id),
|   CONSTRAINT fk_rental_staff FOREIGN KEY (staff_id) REFERENCES staff (staff_id) ON DELETE RESTRICT ON UPDATE CASCADE,
|   CONSTRAINT fk_rental_inventory FOREIGN KEY (inventory_id) REFERENCES inventory (inventory_id) ON DELETE RESTRICT ON UPDATE CASCADE,
|   CONSTRAINT fk_rental_customer FOREIGN KEY (customer_id) REFERENCES customer (customer_id) ON DELETE RESTRICT ON UPDATE CASCADE
| )ENGINE=InnoDB DEFAULT CHARSET=utf8;
< ok 0 0
> exec
|
|
| --
| -- Table structure for table `staff`
| --
|
| CREATE TABLE staff (
|   staff_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
|   first_name VARCHAR(45) NOT NULL,
|   last_name VARCHAR(45) NOT NULL,
|   address_id INT UNSIGNED NOT NULL,
|   picture MEDIUMBLOB DEFAULT NULL,
|   email VARCHAR(50) DEFAULT NULL,
|   store_id INT UNSIGNED NOT NULL,
|   active BOOLEAN NOT NULL DEFAULT TRUE,
|   username VARCHAR(16) NOT NULL,
|   password VARCHAR(40) BINARY DEFAULT NULL,
|   last_update TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY  (staff_id),
|   KEY idx_fk_store_id (store_id),
|   KEY idx_fk_address_id (address_id),
|   CONSTRAINT fk_staff_store FOREIGN KEY (store_id) REFERENCES store (store_id) ON DELETE RESTRICT ON UPDATE CASCADE,
|   CONSTRAINT fk_staff_address FOREIGN KEY (address_id) REFERENCES address (address_id) ON DELETE RESTRICT ON UPDATE CASCADE
| )ENGINE=InnoDB DEFAULT CHARSET=utf8;
< ok 0 0
> exec
|
|
| --
| -- Table structure for table `store`
| --
|
| CREATE TABLE store (
|   store_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
|   manager_staff_id INT UNSIGNED NOT NULL,
|   address_id INT UNSIGNED NOT NULL,
|   last_update TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   PRIMARY KEY  (store_id),
|   UNIQUE KEY idx_unique_manager (manager_staff_id),
|   KEY idx_fk_address_id (address_id),
|   CONSTRAINT fk_store_staff FOREIGN KEY (manager_staff_id) REFERENCES staff (staff_id) ON DELETE RESTRICT ON UPDATE CASCADE,
|   CONSTRAINT fk_store_address FOREIGN KEY (address_id) REFERENCES address (address_id) ON DELETE RESTRICT ON UPDATE CASCADE
| )ENGINE=InnoDB DEFAULT CHARSET=utf8;
< ok 0 0
> exec
|
|
| --
| -- View structure for view `customer_list`
| --
|
| CREATE VIEW customer_list
| AS
| SELECT cu.customer_id AS ID, CONCAT(cu.first_name, _utf8' ', cu.last_name) AS name, a.address AS address, a.postal_code AS `zip code`,
| 	a.phone AS phone, city.city AS city, country.country AS country, IF(cu.active, _utf8'active',_utf8'') AS notes, cu.store_id AS SID
| FROM customer AS cu JOIN address AS a ON cu.address_id = a.address_id JOIN city ON a.city_id = city.city_id
| 	JOIN country ON city.country_id = country.country_id;
< ok 0 0
> exec
|
|
| --
| -- View structure for view `film_list`
| --
|
| CREATE VIEW film_list
| AS
| SELECT film.film_id AS FID, film.title AS title, film.description AS description, category.name AS category, film.rental_rate AS price,
| 	film.length AS length, film.rating AS rating, GROUP_CONCAT(CONCAT(actor.first_name, _utf8' ', actor.last_name) SEPARATOR ', ') AS actors
| FROM category LEFT JOIN film_category ON category.category_id = film_category.category_id LEFT JOIN film ON film_category.film_id = film.film_id
|         JOIN film_actor ON film.film_id = film_actor.film_id
| 	JOIN actor ON film_actor.actor_id = actor.actor_id
| GROUP BY film.film_id, film.title, film.description, film.rental_rate, film.length, film.rating, category.name;
< ok 0 0
> exec
|
|
| --
| -- View structure for view `nicer_but_slower_film_list`
| --
|
| CREATE VIEW nicer_but_slower_film_list
| AS
| SELECT film.film_id AS FID, film.title AS title, film.description AS description, category.name AS category, film.rental_rate AS price,
| 	film.length AS length, film.rating AS rating, GROUP_CONCAT(CONCAT(CONCAT(UCASE(SUBSTR(actor.first_name,1,1)),
| 	LCASE(SUBSTR(actor.first_name,2,LENGTH(actor.first_name))),_utf8' ',CONCAT(UCASE(SUBSTR(actor.last_name,1,1)),
| 	LCASE(SUBSTR(actor.last_name,2,LENGTH(actor.last_name)))))) SEPARATOR ', ') AS actors
| FROM category LEFT JOIN film_category ON category.category_id = film_category.category_id LEFT JOIN film ON film_category.film_id = film.film_id
|         JOIN film_actor ON film.film_id = film_actor.film_id
| 	JOIN actor ON film_actor.actor_id = actor.actor_id
| GROUP BY film.film_id, film.title, film.description, film.rental_rate, film.length, film.rating, category.name;
< ok 0 0
> exec
|
|
| --
| -- View structure for view `staff_list`
| --
|
| CREATE VIEW staff_list
| AS
| SELECT s.staff_id AS ID, CONCAT(s.first_name, _utf8' ', s.last_name) AS name, a.address AS address, a.postal_code AS `zip code`, a.phone AS phone,
| 	city.city AS city, country.country AS country, s.store_id AS SID
| FROM staff AS s JOIN address AS a ON s.address_id = a.address_id JOIN city ON a.city_id = city.city_id
| 	JOIN country ON city.country_id = country.country_id;
< ok 0 0
> exec
|
|
| --
| -- View structure for view `sales_by_store`
| --
|
| CREATE VIEW sales_by_store
| AS
| SELECT
| CONCAT(c.city, _utf8',', cy.country) AS store
| , CONCAT(m.first_name, _utf8' ', m.last_name) AS manager
| , SUM(p.amount) AS total_sales
| FROM payment AS p
| INNER JOIN rental AS r ON p.rental_id = r.rental_id
| INNER JOIN inventory AS i ON r.inventory_id = i.inventory_id
| INNER JOIN store AS s ON i.store_id = s.store_id
| INNER JOIN address AS a ON s.address_id = a.address_id
| INNER JOIN city AS c ON a.city_id = c.city_id
| INNER JOIN country AS cy ON c.country_id = cy.country_id
| INNER JOIN staff AS m ON s.manager_staff_id = m.staff_id
| GROUP BY s.store_id
| ORDER BY cy.country, c.city;
< ok 0 0
> exec
|
|
| --
| -- View structure for view `sales_by_film_category`
| --
| -- Note that total sales will add up to >100% because
| -- some titles belong to more than 1 category
| --
|
| CREATE VIEW sales_by_film_category
| AS
| SELECT
| c.name AS category
| , SUM(p.amount) AS total_sales
| FROM payment AS p
| INNER JOIN rental AS r ON p.rental_id = r.rental_id
| INNER JOIN inventory AS i ON r.inventory_id = i.inventory_id
| INNER JOIN film AS f ON i.film_id = f.film_id
| INNER JOIN film_category AS fc ON f.film_id = fc.film_id
| INNER JOIN category AS c ON fc.category_id = c.category_id
| GROUP BY c.name
| ORDER BY total_sales DESC;
< ok 0 0
> exec
|
|
| --
| -- View structure for view `actor_info`
| --
|
| CREATE DEFINER=CURRENT_USER SQL SECURITY INVOKER VIEW actor_info
| AS
| SELECT
| a.actor_id,
| a.first_name,
| a.last_name,
| GROUP_CONCAT(DISTINCT CONCAT(c.name, ': ',
| 		(SELECT GROUP_CONCAT(f.title ORDER BY f.title SEPARATOR ', ')
|                     FROM sakila.film f
|                     INNER JOIN sakila.film_category fc
|                       ON f.film_id = fc.film_id
|                     INNER JOIN sakila.film_actor fa
|                       ON f.film_id = fa.film_id
|                     WHERE fc.category_id = c.category_id
|                     AND fa.actor_id = a.actor_id
|                  )
|              )
|              ORDER BY c.name SEPARATOR '; ')
| AS film_info
| FROM sakila.actor a
| LEFT JOIN sakila.film_actor fa
|   ON a.actor_id = fa.actor_id
| LEFT JOIN sakila.film_category fc
|   ON fa.film_id = fc.film_id
| LEFT JOIN sakila.category c
|   ON fc.category_id = c.category_id
| GROUP BY a.actor_id, a.first_name, a.last_name;
< ok 0 0
> exec
|
|
| --
| -- Procedure structure for procedure `rewards_report`
| --
|
|
| CREATE PROCEDURE rewards_report (
|     IN min_monthly_purchases TINYINT UNSIGNED
|     , IN min_dollar_amount_purchased DECIMAL(10,2) UNSIGNED
|     , OUT count_rewardees INT
| )
| LANGUAGE SQL
| NOT DETERMINISTIC
| READS SQL DATA
| SQL SECURITY DEFINER
| COMMENT 'Provides a customizable report on best customers'
| proc: BEGIN
|
|     DECLARE last_month_start DATE;
|     DECLARE last_month_end DATE;
|
|     /* Some sanity checks... */
|     IF min_monthly_purchases = 0 THEN
|         SELECT 'Minimum monthly purchases parameter must be > 0';
|         LEAVE proc;
|     END IF;
|     IF min_dollar_amount_purchased = 0.00 THEN
|         SELECT 'Minimum monthly dollar amount purchased parameter must be > $0.00';
|         LEAVE proc;
|     END IF;
|
|     /* Determine start and end time periods */
|     SET last_month_start = DATE_SUB(CURRENT_DATE(), INTERVAL 1 MONTH);
|     SET last_month_start = STR_TO_DATE(CONCAT(YEAR(last_month_start),'-',MONTH(last_month_start),'-01'),'%Y-%m-%d');
|     SET last_month_end = LAST_DAY(last_month_start);
|
|     /*
|         Create a temporary storage area for
|         Customer IDs.
|     */
|     CREATE TEMPORARY TABLE tmpCustomer (customer_id INT UNSIGNED NOT NULL PRIMARY KEY);
|
|     /*
|         Find all customers meeting the
|         monthly purchase requirements
|     */
|     INSERT INTO tmpCustomer (customer_id)
|     SELECT p.customer_id
|     FROM payment AS p
|     WHERE DATE(p.payment_date) BETWEEN last_month_start AND last_month_end
|     GROUP BY customer_id
|     HAVING SUM(p.amount) > min_dollar_amount_purchased
|     AND COUNT(customer_id) > min_monthly_purchases;
|
|     /* Populate OUT parameter with count of found customers */
|     SELECT COUNT(*) FROM tmpCustomer INTO count_rewardees;
|
|     /*
|         Output ALL customer information of matching rewardees.
|         Customize output as needed.
|     */
|     SELECT c.*
|     FROM tmpCustomer AS t
|     INNER JOIN customer AS c ON t.customer_id = c.customer_id;
|
|     /* Clean up */
|     DROP TABLE tmpCustomer;
| END ;
< ok 0 0
> exec
|
|
|
|
| CREATE FUNCTION get_customer_balance(p_customer_id INT, p_effective_date DATETIME) RETURNS DECIMAL(5,2)
|     DETERMINISTIC
|     READS SQL DATA
| BEGIN
|
|        #OK, WE NEED TO CALCULATE THE CURRENT BALANCE GIVEN A CUSTOMER_ID AND A DATE
|        #THAT WE WANT THE BALANCE TO BE EFFECTIVE FOR. THE BALANCE IS:
|        #   1) RENTAL FEES FOR ALL PREVIOUS RENTALS
|        #   2) ONE DOLLAR FOR EVERY DAY THE PREVIOUS RENTALS ARE OVERDUE
|        #   3) IF A FILM IS MORE THAN RENTAL_DURATION * 2 OVERDUE, CHARGE THE REPLACEMENT_COST
|        #   4) SUBTRACT ALL PAYMENTS MADE BEFORE THE DATE SPECIFIED
|
|   DECLARE v_rentfees DECIMAL(5,2); #FEES PAID TO RENT THE VIDEOS INITIALLY
|   DECLARE v_overfees INTEGER;      #LATE FEES FOR PRIOR RENTALS
|   DECLARE v_payments DECIMAL(5,2); #SUM OF PAYMENTS MADE PREVIOUSLY
|
|   SELECT IFNULL(SUM(film.rental_rate),0) INTO v_rentfees
|     FROM film, inventory, rental
|     WHERE film.film_id = inventory.film_id
|       AND inventory.inventory_id = rental.inventory_id
|       AND rental.rental_date <= p_effective_date
|       AND rental.customer_id = p_customer_id;
|
|   SELECT IFNULL(SUM(IF((TO_DAYS(rental.return_date) - TO_DAYS(rental.rental_date)) > film.rental_duration,
|         ((TO_DAYS(rental.return_date) - TO_DAYS(rental.rental_date)) - film.rental_duration),0)),0) INTO v_overfees
|     FROM rental, inventory, film
|     WHERE film.film_id = inventory.film_id
|       AND inventory.inventory_id = rental.inventory_id
|       AND rental.rental_date <= p_effective_date
|       AND rental.customer_id = p_customer_id;
|
|
|   SELECT IFNULL(SUM(payment.amount),0) INTO v_payments
|     FROM payment
|
|     WHERE payment.payment_date <= p_effective_date
|     AND payment.customer_id = p_customer_id;
|
|   RETURN v_rentfees + v_overfees - v_payments;
| END ;
< ok 0 0
> exec
|
|
|
|
| CREATE PROCEDURE film_in_stock(IN p_film_id INT, IN p_store_id INT, OUT p_film_count INT)
| READS SQL DATA
| BEGIN
|      SELECT inventory_id
|      FROM inventory
|      WHERE film_id = p_film_id
|      AND store_id = p_store_id
|      AND inventory_in_stock(inventory_id);
|
|      SELECT FOUND_ROWS() INTO p_film_count;
| END ;
< ok 0 0
> exec
|
|
|
|
| CREATE PROCEDURE film_not_in_stock(IN p_film_id INT, IN p_store_id INT, OUT p_film_count INT)
| READS SQL DATA
| BEGIN
|      SELECT inventory_id
|      FROM inventory
|      WHERE film_id = p_film_id
|      AND store_id = p_store_id
|      AND NOT inventory_in_stock(inventory_id);
|
|      SELECT FOUND_ROWS() INTO p_film_count;
| END ;
< ok 0 0
> exec
|
|
|
|
| CREATE FUNCTION inventory_held_by_customer(p_inventory_id INT) RETURNS INT
| READS SQL DATA
| BEGIN
|   DECLARE v_customer_id INT;
|   DECLARE EXIT HANDLER FOR NOT FOUND RETURN NULL;
|
|   SELECT customer_id INTO v_customer_id
|   FROM rental
|   WHERE return_date IS NULL
|   AND inventory_id = p_inventory_id;
|
|   RETURN v_customer_id;
| END ;
< ok 0 0
> exec
|
|
|
|
| CREATE FUNCTION inventory_in_stock(p_inventory_id INT) RETURNS BOOLEAN
| READS SQL DATA
| BEGIN
|     DECLARE v_rentals INT;
|     DECLARE v_out     INT;
|
|     #AN ITEM IS IN-STOCK IF THERE ARE EITHER NO ROWS IN THE rental TABLE
|     #FOR THE ITEM OR ALL ROWS HAVE return_date POPULATED
|
|     SELECT COUNT(*) INTO v_rentals
|     FROM rental
|     WHERE inventory_id = p_inventory_id;
|
|     IF v_rentals = 0 THEN
|       RETURN TRUE;
|     END IF;
|
|     SELECT COUNT(rental_id) INTO v_out
|     FROM inventory LEFT JOIN rental USING(inventory_id)
|     WHERE inventory.inventory_id = p_inventory_id
|     AND rental.return_date IS NULL;
|
|     IF v_out > 0 THEN
|       RETURN FALSE;
|     ELSE
|       RETURN TRUE;
|     END IF;
| END ;
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
| sakila
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
| actor
: bytes
| idx_actor_last_name
: bytes
| last_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| actor
: bytes
| PRIMARY
: bytes
| actor_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| address
: bytes
| idx_fk_city_id
: bytes
| city_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| address
: bytes
| PRIMARY
: bytes
| address_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| category
: bytes
| PRIMARY
: bytes
| category_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| city
: bytes
| idx_fk_country_id
: bytes
| country_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| city
: bytes
| PRIMARY
: bytes
| city_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| country
: bytes
| PRIMARY
: bytes
| country_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| customer
: bytes
| idx_fk_address_id
: bytes
| address_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| customer
: bytes
| idx_fk_store_id
: bytes
| store_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| customer
: bytes
| idx_last_name
: bytes
| last_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| customer
: bytes
| PRIMARY
: bytes
| customer_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| film
: bytes
| idx_fk_language_id
: bytes
| language_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| film
: bytes
| idx_fk_original_language_id
: bytes
| original_language_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| film
: bytes
| idx_title
: bytes
| title
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| film
: bytes
| PRIMARY
: bytes
| film_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| film_actor
: bytes
| idx_fk_film_id
: bytes
| film_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| film_actor
: bytes
| PRIMARY
: bytes
| actor_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| film_actor
: bytes
| PRIMARY
: bytes
| film_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| film_category
: bytes
| fk_film_category_category
: bytes
| category_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| film_category
: bytes
| PRIMARY
: bytes
| film_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| film_category
: bytes
| PRIMARY
: bytes
| category_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| film_text
: bytes
| idx_title_description
: bytes
| title
: null
: int64 -1
: null
: int64 1
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
< row
: bytes
| film_text
: bytes
| idx_title_description
: bytes
| description
: null
: int64 -1
: null
: int64 2
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
< row
: bytes
| film_text
: bytes
| PRIMARY
: bytes
| film_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| inventory
: bytes
| idx_fk_film_id
: bytes
| film_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| inventory
: bytes
| idx_store_id_film_id
: bytes
| store_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| inventory
: bytes
| idx_store_id_film_id
: bytes
| film_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| inventory
: bytes
| PRIMARY
: bytes
| inventory_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| language
: bytes
| PRIMARY
: bytes
| language_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| payment
: bytes
| fk_payment_rental
: bytes
| rental_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| payment
: bytes
| idx_fk_customer_id
: bytes
| customer_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| payment
: bytes
| idx_fk_staff_id
: bytes
| staff_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| payment
: bytes
| PRIMARY
: bytes
| payment_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| rental
: bytes
| idx_fk_customer_id
: bytes
| customer_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| rental
: bytes
| idx_fk_inventory_id
: bytes
| inventory_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| rental
: bytes
| idx_fk_staff_id
: bytes
| staff_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| rental
: bytes
| PRIMARY
: bytes
| rental_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| rental
: bytes
| rental_date
: bytes
| rental_date
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| rental
: bytes
| rental_date
: bytes
| inventory_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| rental
: bytes
| rental_date
: bytes
| customer_id
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| staff
: bytes
| idx_fk_address_id
: bytes
| address_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| staff
: bytes
| idx_fk_store_id
: bytes
| store_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| staff
: bytes
| PRIMARY
: bytes
| staff_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| store
: bytes
| idx_fk_address_id
: bytes
| address_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| store
: bytes
| idx_unique_manager
: bytes
| manager_staff_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| store
: bytes
| PRIMARY
: bytes
| store_id
: bytes
| A
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
| sakila
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
| actor
: bytes
| actor_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| actor
: bytes
| first_name
: int64 2
: null
: bytes
| NO
: bytes
| varchar(45)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| actor
: bytes
| last_name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(45)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| actor
: bytes
| last_update
: int64 4
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| actor_info
: bytes
| actor_id
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| actor_info
: bytes
| first_name
: int64 2
: null
: bytes
| NO
: bytes
| varchar(45)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| actor_info
: bytes
| last_name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(45)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| actor_info
: bytes
| film_info
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| address
: bytes
| address_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| address
: bytes
| address
: int64 2
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| address
: bytes
| address2
: int64 3
: null
: bytes
| YES
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| address
: bytes
| district
: int64 4
: null
: bytes
| NO
: bytes
| varchar(20)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| address
: bytes
| city_id
: int64 5
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| address
: bytes
| postal_code
: int64 6
: null
: bytes
| YES
: bytes
| varchar(10)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| address
: bytes
| phone
: int64 7
: null
: bytes
| NO
: bytes
| varchar(20)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| address
: bytes
| last_update
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
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| category
: bytes
| category_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| category
: bytes
| name
: int64 2
: null
: bytes
| NO
: bytes
| varchar(25)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| category
: bytes
| last_update
: int64 3
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| city
: bytes
| city_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| city
: bytes
| city
: int64 2
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| city
: bytes
| country_id
: int64 3
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| city
: bytes
| last_update
: int64 4
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| country
: bytes
| country_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| country
: bytes
| country
: int64 2
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| country
: bytes
| last_update
: int64 3
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| customer
: bytes
| customer_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| customer
: bytes
| store_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer
: bytes
| first_name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(45)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer
: bytes
| last_name
: int64 4
: null
: bytes
| NO
: bytes
| varchar(45)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer
: bytes
| email
: int64 5
: null
: bytes
| YES
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer
: bytes
| address_id
: int64 6
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer
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
| customer
: bytes
| create_date
: int64 8
: null
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer
: bytes
| last_update
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| YES
: bytes
| timestamp
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
| customer_list
: bytes
| ID
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer_list
: bytes
| name
: int64 2
: null
: bytes
| YES
: bytes
| varchar(91)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer_list
: bytes
| address
: int64 3
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer_list
: bytes
| zip code
: int64 4
: null
: bytes
| YES
: bytes
| varchar(10)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer_list
: bytes
| phone
: int64 5
: null
: bytes
| NO
: bytes
| varchar(20)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer_list
: bytes
| city
: int64 6
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer_list
: bytes
| country
: int64 7
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer_list
: bytes
| notes
: int64 8
: bytes
| ''
: bytes
| NO
: bytes
| varchar(6)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer_list
: bytes
| SID
: int64 9
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| film_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| film
: bytes
| title
: int64 2
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| description
: int64 3
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| release_year
: int64 4
: null
: bytes
| YES
: bytes
| year
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| language_id
: int64 5
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| original_language_id
: int64 6
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| rental_duration
: int64 7
: bytes
| '3'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| rental_rate
: int64 8
: bytes
| '4.99'
: bytes
| NO
: bytes
| decimal(4,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| length
: int64 9
: null
: bytes
| YES
: bytes
| smallint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| replacement_cost
: int64 10
: bytes
| '19.99'
: bytes
| NO
: bytes
| decimal(5,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| rating
: int64 11
: bytes
| 'G'
: bytes
| YES
: bytes
| enum('G','PG','PG-13','R','NC-17')
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| special_features
: int64 12
: null
: bytes
| YES
: bytes
| set('Trailers','Commentaries','Deleted Scenes','Behind the Scenes')
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| last_update
: int64 13
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| film_actor
: bytes
| actor_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_actor
: bytes
| film_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_actor
: bytes
| last_update
: int64 3
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| film_category
: bytes
| film_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_category
: bytes
| category_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_category
: bytes
| last_update
: int64 3
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| film_list
: bytes
| FID
: int64 1
: bytes
| '0'
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_list
: bytes
| title
: int64 2
: null
: bytes
| YES
: bytes
| varchar(255)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_list
: bytes
| description
: int64 3
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_list
: bytes
| category
: int64 4
: null
: bytes
| NO
: bytes
| varchar(25)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_list
: bytes
| price
: int64 5
: bytes
| '4.99'
: bytes
| YES
: bytes
| decimal(4,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_list
: bytes
| length
: int64 6
: null
: bytes
| YES
: bytes
| smallint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_list
: bytes
| rating
: int64 7
: bytes
| 'G'
: bytes
| YES
: bytes
| enum('G','PG','PG-13','R','NC-17')
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_list
: bytes
| actors
: int64 8
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_text
: bytes
| film_id
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
|
: null
< row
: bytes
| film_text
: bytes
| title
: int64 2
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_text
: bytes
| description
: int64 3
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| inventory
: bytes
| inventory_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| inventory
: bytes
| film_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| inventory
: bytes
| store_id
: int64 3
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| inventory
: bytes
| last_update
: int64 4
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| language
: bytes
| language_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| language
: bytes
| name
: int64 2
: null
: bytes
| NO
: bytes
| char(20)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| language
: bytes
| last_update
: int64 3
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| nicer_but_slower_film_list
: bytes
| FID
: int64 1
: bytes
| '0'
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| nicer_but_slower_film_list
: bytes
| title
: int64 2
: null
: bytes
| YES
: bytes
| varchar(255)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| nicer_but_slower_film_list
: bytes
| description
: int64 3
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| nicer_but_slower_film_list
: bytes
| category
: int64 4
: null
: bytes
| NO
: bytes
| varchar(25)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| nicer_but_slower_film_list
: bytes
| price
: int64 5
: bytes
| '4.99'
: bytes
| YES
: bytes
| decimal(4,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| nicer_but_slower_film_list
: bytes
| length
: int64 6
: null
: bytes
| YES
: bytes
| smallint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| nicer_but_slower_film_list
: bytes
| rating
: int64 7
: bytes
| 'G'
: bytes
| YES
: bytes
| enum('G','PG','PG-13','R','NC-17')
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| nicer_but_slower_film_list
: bytes
| actors
: int64 8
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| payment
: bytes
| payment_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| payment
: bytes
| customer_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| payment
: bytes
| staff_id
: int64 3
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| payment
: bytes
| rental_id
: int64 4
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
| payment
: bytes
| amount
: int64 5
: null
: bytes
| NO
: bytes
| decimal(5,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| payment
: bytes
| payment_date
: int64 6
: null
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| payment
: bytes
| last_update
: int64 7
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| YES
: bytes
| timestamp
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
| rental
: bytes
| rental_id
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
| rental
: bytes
| rental_date
: int64 2
: null
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| rental
: bytes
| inventory_id
: int64 3
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| rental
: bytes
| customer_id
: int64 4
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| rental
: bytes
| return_date
: int64 5
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| rental
: bytes
| staff_id
: int64 6
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| rental
: bytes
| last_update
: int64 7
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| sales_by_film_category
: bytes
| category
: int64 1
: null
: bytes
| NO
: bytes
| varchar(25)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sales_by_film_category
: bytes
| total_sales
: int64 2
: null
: bytes
| YES
: bytes
| decimal(27,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sales_by_store
: bytes
| store
: int64 1
: null
: bytes
| YES
: bytes
| varchar(101)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sales_by_store
: bytes
| manager
: int64 2
: null
: bytes
| YES
: bytes
| varchar(91)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sales_by_store
: bytes
| total_sales
: int64 3
: null
: bytes
| YES
: bytes
| decimal(27,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff
: bytes
| staff_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| staff
: bytes
| first_name
: int64 2
: null
: bytes
| NO
: bytes
| varchar(45)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff
: bytes
| last_name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(45)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff
: bytes
| address_id
: int64 4
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff
: bytes
| picture
: int64 5
: null
: bytes
| YES
: bytes
| mediumblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff
: bytes
| email
: int64 6
: null
: bytes
| YES
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff
: bytes
| store_id
: int64 7
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff
: bytes
| active
: int64 8
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
| staff
: bytes
| username
: int64 9
: null
: bytes
| NO
: bytes
| varchar(16)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff
: bytes
| password
: int64 10
: null
: bytes
| YES
: bytes
| varchar(40)
: bytes
| utf8mb3
: bytes
| utf8mb3_bin
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff
: bytes
| last_update
: int64 11
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| staff_list
: bytes
| ID
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff_list
: bytes
| name
: int64 2
: null
: bytes
| YES
: bytes
| varchar(91)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff_list
: bytes
| address
: int64 3
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff_list
: bytes
| zip code
: int64 4
: null
: bytes
| YES
: bytes
| varchar(10)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff_list
: bytes
| phone
: int64 5
: null
: bytes
| NO
: bytes
| varchar(20)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff_list
: bytes
| city
: int64 6
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff_list
: bytes
| country
: int64 7
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff_list
: bytes
| SID
: int64 8
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| store
: bytes
| store_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| store
: bytes
| manager_staff_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| store
: bytes
| address_id
: int64 3
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| store
: bytes
| last_update
: int64 4
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| sakila
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
| sakila
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| actor_info
: bytes
| select `a`.`actor_id` AS `actor_id`,`a`.`first_name` AS `first_name`,`a`.`last_name` AS `last_name`,group_concat(distinct concat(`c`.`name`,': ',(select group_concat(`f`.`title` order by `f`.`title` ASC separator ', ') from ((`sakila`.`film` `f` join `sakila`.`film_category` `fc` on((`f`.`film_id` = `fc`.`film_id`))) join `sakila`.`film_actor` `fa` on((`f`.`film_id` = `fa`.`film_id`))) where ((`fc`.`category_id` = `c`.`category_id`) and (`fa`.`actor_id` = `a`.`actor_id`)))) order by `c`.`name` ASC separator '; ') AS `film_info` from (((`sakila`.`actor` `a` left join `sakila`.`film_actor` `fa` on((`a`.`actor_id` = `fa`.`actor_id`))) left join `sakila`.`film_category` `fc` on((`fa`.`film_id` = `fc`.`film_id`))) left join `sakila`.`category` `c` on((`fc`.`category_id` = `c`.`category_id`))) group by `a`.`actor_id`,`a`.`first_name`,`a`.`last_name`
< row
: bytes
| customer_list
: bytes
| select `cu`.`customer_id` AS `ID`,concat(`cu`.`first_name`,' ',`cu`.`last_name`) AS `name`,`a`.`address` AS `address`,`a`.`postal_code` AS `zip code`,`a`.`phone` AS `phone`,`sakila`.`city`.`city` AS `city`,`sakila`.`country`.`country` AS `country`,if(`cu`.`active`,'active','') AS `notes`,`cu`.`store_id` AS `SID` from (((`sakila`.`customer` `cu` join `sakila`.`address` `a` on((`cu`.`address_id` = `a`.`address_id`))) join `sakila`.`city` on((`a`.`city_id` = `sakila`.`city`.`city_id`))) join `sakila`.`country` on((`sakila`.`city`.`country_id` = `sakila`.`country`.`country_id`)))
< row
: bytes
| film_list
: bytes
| select `sakila`.`film`.`film_id` AS `FID`,`sakila`.`film`.`title` AS `title`,`sakila`.`film`.`description` AS `description`,`sakila`.`category`.`name` AS `category`,`sakila`.`film`.`rental_rate` AS `price`,`sakila`.`film`.`length` AS `length`,`sakila`.`film`.`rating` AS `rating`,group_concat(concat(`sakila`.`actor`.`first_name`,' ',`sakila`.`actor`.`last_name`) separator ', ') AS `actors` from ((((`sakila`.`category` left join `sakila`.`film_category` on((`sakila`.`category`.`category_id` = `sakila`.`film_category`.`category_id`))) left join `sakila`.`film` on((`sakila`.`film_category`.`film_id` = `sakila`.`film`.`film_id`))) join `sakila`.`film_actor` on((`sakila`.`film`.`film_id` = `sakila`.`film_actor`.`film_id`))) join `sakila`.`actor` on((`sakila`.`film_actor`.`actor_id` = `sakila`.`actor`.`actor_id`))) group by `sakila`.`film`.`film_id`,`sakila`.`film`.`title`,`sakila`.`film`.`description`,`sakila`.`film`.`rental_rate`,`sakila`.`film`.`length`,`sakila`.`film`.`rating`,`sakila`.`category`.`name`
< row
: bytes
| nicer_but_slower_film_list
: bytes
| select `sakila`.`film`.`film_id` AS `FID`,`sakila`.`film`.`title` AS `title`,`sakila`.`film`.`description` AS `description`,`sakila`.`category`.`name` AS `category`,`sakila`.`film`.`rental_rate` AS `price`,`sakila`.`film`.`length` AS `length`,`sakila`.`film`.`rating` AS `rating`,group_concat(concat(concat(upper(substr(`sakila`.`actor`.`first_name`,1,1)),lower(substr(`sakila`.`actor`.`first_name`,2,length(`sakila`.`actor`.`first_name`))),' ',concat(upper(substr(`sakila`.`actor`.`last_name`,1,1)),lower(substr(`sakila`.`actor`.`last_name`,2,length(`sakila`.`actor`.`last_name`)))))) separator ', ') AS `actors` from ((((`sakila`.`category` left join `sakila`.`film_category` on((`sakila`.`category`.`category_id` = `sakila`.`film_category`.`category_id`))) left join `sakila`.`film` on((`sakila`.`film_category`.`film_id` = `sakila`.`film`.`film_id`))) join `sakila`.`film_actor` on((`sakila`.`film`.`film_id` = `sakila`.`film_actor`.`film_id`))) join `sakila`.`actor` on((`sakila`.`film_actor`.`actor_id` = `sakila`.`actor`.`actor_id`))) group by `sakila`.`film`.`film_id`,`sakila`.`film`.`title`,`sakila`.`film`.`description`,`sakila`.`film`.`rental_rate`,`sakila`.`film`.`length`,`sakila`.`film`.`rating`,`sakila`.`category`.`name`
< row
: bytes
| sales_by_film_category
: bytes
| select `c`.`name` AS `category`,sum(`p`.`amount`) AS `total_sales` from (((((`sakila`.`payment` `p` join `sakila`.`rental` `r` on((`p`.`rental_id` = `r`.`rental_id`))) join `sakila`.`inventory` `i` on((`r`.`inventory_id` = `i`.`inventory_id`))) join `sakila`.`film` `f` on((`i`.`film_id` = `f`.`film_id`))) join `sakila`.`film_category` `fc` on((`f`.`film_id` = `fc`.`film_id`))) join `sakila`.`category` `c` on((`fc`.`category_id` = `c`.`category_id`))) group by `c`.`name` order by `total_sales` desc
< row
: bytes
| sales_by_store
: bytes
| select concat(`c`.`city`,',',`cy`.`country`) AS `store`,concat(`m`.`first_name`,' ',`m`.`last_name`) AS `manager`,sum(`p`.`amount`) AS `total_sales` from (((((((`sakila`.`payment` `p` join `sakila`.`rental` `r` on((`p`.`rental_id` = `r`.`rental_id`))) join `sakila`.`inventory` `i` on((`r`.`inventory_id` = `i`.`inventory_id`))) join `sakila`.`store` `s` on((`i`.`store_id` = `s`.`store_id`))) join `sakila`.`address` `a` on((`s`.`address_id` = `a`.`address_id`))) join `sakila`.`city` `c` on((`a`.`city_id` = `c`.`city_id`))) join `sakila`.`country` `cy` on((`c`.`country_id` = `cy`.`country_id`))) join `sakila`.`staff` `m` on((`s`.`manager_staff_id` = `m`.`staff_id`))) group by `s`.`store_id` order by `cy`.`country`,`c`.`city`
< row
: bytes
| staff_list
: bytes
| select `s`.`staff_id` AS `ID`,concat(`s`.`first_name`,' ',`s`.`last_name`) AS `name`,`a`.`address` AS `address`,`a`.`postal_code` AS `zip code`,`a`.`phone` AS `phone`,`sakila`.`city`.`city` AS `city`,`sakila`.`country`.`country` AS `country`,`s`.`store_id` AS `SID` from (((`sakila`.`staff` `s` join `sakila`.`address` `a` on((`s`.`address_id` = `a`.`address_id`))) join `sakila`.`city` on((`a`.`city_id` = `sakila`.`city`.`city_id`))) join `sakila`.`country` on((`sakila`.`city`.`country_id` = `sakila`.`country`.`country_id`)))
> query
| SHOW CREATE VIEW `sakila`.`actor_info`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| actor_info
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY INVOKER VIEW `actor_info` AS select `a`.`actor_id` AS `actor_id`,`a`.`first_name` AS `first_name`,`a`.`last_name` AS `last_name`,group_concat(distinct concat(`c`.`name`,': ',(select group_concat(`f`.`title` order by `f`.`title` ASC separator ', ') from ((`film` `f` join `film_category` `fc` on((`f`.`film_id` = `fc`.`film_id`))) join `film_actor` `fa` on((`f`.`film_id` = `fa`.`film_id`))) where ((`fc`.`category_id` = `c`.`category_id`) and (`fa`.`actor_id` = `a`.`actor_id`)))) order by `c`.`name` ASC separator '; ') AS `film_info` from (((`actor` `a` left join `film_actor` `fa` on((`a`.`actor_id` = `fa`.`actor_id`))) left join `film_category` `fc` on((`fa`.`film_id` = `fc`.`film_id`))) left join `category` `c` on((`fc`.`category_id` = `c`.`category_id`))) group by `a`.`actor_id`,`a`.`first_name`,`a`.`last_name`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `sakila`.`customer_list`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| customer_list
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `customer_list` AS select `cu`.`customer_id` AS `ID`,concat(`cu`.`first_name`,_utf8mb3' ',`cu`.`last_name`) AS `name`,`a`.`address` AS `address`,`a`.`postal_code` AS `zip code`,`a`.`phone` AS `phone`,`city`.`city` AS `city`,`country`.`country` AS `country`,if(`cu`.`active`,_utf8mb3'active',_utf8mb3'') AS `notes`,`cu`.`store_id` AS `SID` from (((`customer` `cu` join `address` `a` on((`cu`.`address_id` = `a`.`address_id`))) join `city` on((`a`.`city_id` = `city`.`city_id`))) join `country` on((`city`.`country_id` = `country`.`country_id`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `sakila`.`film_list`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| film_list
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `film_list` AS select `film`.`film_id` AS `FID`,`film`.`title` AS `title`,`film`.`description` AS `description`,`category`.`name` AS `category`,`film`.`rental_rate` AS `price`,`film`.`length` AS `length`,`film`.`rating` AS `rating`,group_concat(concat(`actor`.`first_name`,_utf8mb3' ',`actor`.`last_name`) separator ', ') AS `actors` from ((((`category` left join `film_category` on((`category`.`category_id` = `film_category`.`category_id`))) left join `film` on((`film_category`.`film_id` = `film`.`film_id`))) join `film_actor` on((`film`.`film_id` = `film_actor`.`film_id`))) join `actor` on((`film_actor`.`actor_id` = `actor`.`actor_id`))) group by `film`.`film_id`,`film`.`title`,`film`.`description`,`film`.`rental_rate`,`film`.`length`,`film`.`rating`,`category`.`name`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `sakila`.`nicer_but_slower_film_list`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| nicer_but_slower_film_list
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `nicer_but_slower_film_list` AS select `film`.`film_id` AS `FID`,`film`.`title` AS `title`,`film`.`description` AS `description`,`category`.`name` AS `category`,`film`.`rental_rate` AS `price`,`film`.`length` AS `length`,`film`.`rating` AS `rating`,group_concat(concat(concat(upper(substr(`actor`.`first_name`,1,1)),lower(substr(`actor`.`first_name`,2,length(`actor`.`first_name`))),_utf8mb3' ',concat(upper(substr(`actor`.`last_name`,1,1)),lower(substr(`actor`.`last_name`,2,length(`actor`.`last_name`)))))) separator ', ') AS `actors` from ((((`category` left join `film_category` on((`category`.`category_id` = `film_category`.`category_id`))) left join `film` on((`film_category`.`film_id` = `film`.`film_id`))) join `film_actor` on((`film`.`film_id` = `film_actor`.`film_id`))) join `actor` on((`film_actor`.`actor_id` = `actor`.`actor_id`))) group by `film`.`film_id`,`film`.`title`,`film`.`description`,`film`.`rental_rate`,`film`.`length`,`film`.`rating`,`category`.`name`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `sakila`.`sales_by_film_category`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| sales_by_film_category
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `sales_by_film_category` AS select `c`.`name` AS `category`,sum(`p`.`amount`) AS `total_sales` from (((((`payment` `p` join `rental` `r` on((`p`.`rental_id` = `r`.`rental_id`))) join `inventory` `i` on((`r`.`inventory_id` = `i`.`inventory_id`))) join `film` `f` on((`i`.`film_id` = `f`.`film_id`))) join `film_category` `fc` on((`f`.`film_id` = `fc`.`film_id`))) join `category` `c` on((`fc`.`category_id` = `c`.`category_id`))) group by `c`.`name` order by `total_sales` desc
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `sakila`.`sales_by_store`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| sales_by_store
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `sales_by_store` AS select concat(`c`.`city`,_utf8mb3',',`cy`.`country`) AS `store`,concat(`m`.`first_name`,_utf8mb3' ',`m`.`last_name`) AS `manager`,sum(`p`.`amount`) AS `total_sales` from (((((((`payment` `p` join `rental` `r` on((`p`.`rental_id` = `r`.`rental_id`))) join `inventory` `i` on((`r`.`inventory_id` = `i`.`inventory_id`))) join `store` `s` on((`i`.`store_id` = `s`.`store_id`))) join `address` `a` on((`s`.`address_id` = `a`.`address_id`))) join `city` `c` on((`a`.`city_id` = `c`.`city_id`))) join `country` `cy` on((`c`.`country_id` = `cy`.`country_id`))) join `staff` `m` on((`s`.`manager_staff_id` = `m`.`staff_id`))) group by `s`.`store_id` order by `cy`.`country`,`c`.`city`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `sakila`.`staff_list`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| staff_list
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `staff_list` AS select `s`.`staff_id` AS `ID`,concat(`s`.`first_name`,_utf8mb3' ',`s`.`last_name`) AS `name`,`a`.`address` AS `address`,`a`.`postal_code` AS `zip code`,`a`.`phone` AS `phone`,`city`.`city` AS `city`,`country`.`country` AS `country`,`s`.`store_id` AS `SID` from (((`staff` `s` join `address` `a` on((`s`.`address_id` = `a`.`address_id`))) join `city` on((`a`.`city_id` = `city`.`city_id`))) join `country` on((`city`.`country_id` = `country`.`country_id`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| sakila
< result
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "EVENT_OBJECT_TABLE" "VARCHAR" "string"
< column "EVENT_MANIPULATION" "ENUM" "string"
< column "ACTION_TIMING" "ENUM" "string"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< row
: bytes
| ins_film
: bytes
| film
: bytes
| INSERT
: bytes
| AFTER
: bytes
| BEGIN
|     INSERT INTO film_text (film_id, title, description)
|         VALUES (new.film_id, new.title, new.description);
|   END
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
< row
: bytes
| upd_film
: bytes
| film
: bytes
| UPDATE
: bytes
| AFTER
: bytes
| BEGIN
|     IF (old.title != new.title) or (old.description != new.description)
|     THEN
|         UPDATE film_text
|             SET title=new.title,
|                 description=new.description,
|                 film_id=new.film_id
|         WHERE film_id=old.film_id;
|     END IF;
|   END
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
< row
: bytes
| del_film
: bytes
| film
: bytes
| DELETE
: bytes
| AFTER
: bytes
| BEGIN
|     DELETE FROM film_text WHERE film_id = old.film_id;
|   END
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| sakila
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
| sakila
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
< row
: bytes
| address
: bytes
| fk_address_city
: bytes
| city
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| city
: bytes
| fk_city_country
: bytes
| country
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| customer
: bytes
| fk_customer_address
: bytes
| address
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| customer
: bytes
| fk_customer_store
: bytes
| store
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| film_actor
: bytes
| fk_film_actor_actor
: bytes
| actor
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| film_actor
: bytes
| fk_film_actor_film
: bytes
| film
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| film_category
: bytes
| fk_film_category_category
: bytes
| category
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| film_category
: bytes
| fk_film_category_film
: bytes
| film
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| film
: bytes
| fk_film_language
: bytes
| language
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| film
: bytes
| fk_film_language_original
: bytes
| language
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| inventory
: bytes
| fk_inventory_film
: bytes
| film
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| inventory
: bytes
| fk_inventory_store
: bytes
| store
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| payment
: bytes
| fk_payment_customer
: bytes
| customer
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| payment
: bytes
| fk_payment_rental
: bytes
| rental
: bytes
| SET NULL
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| payment
: bytes
| fk_payment_staff
: bytes
| staff
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| rental
: bytes
| fk_rental_customer
: bytes
| customer
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| rental
: bytes
| fk_rental_inventory
: bytes
| inventory
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| rental
: bytes
| fk_rental_staff
: bytes
| staff
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| staff
: bytes
| fk_staff_address
: bytes
| address
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| staff
: bytes
| fk_staff_store
: bytes
| store
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| store
: bytes
| fk_store_address
: bytes
| address
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| store
: bytes
| fk_store_staff
: bytes
| staff
: bytes
| RESTRICT
: bytes
| CASCADE
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
| sakila
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| sakila
: bytes
| address
: bytes
| fk_address_city
: bytes
| city_id
: bytes
| sakila
: bytes
| city_id
< row
: bytes
| sakila
: bytes
| city
: bytes
| fk_city_country
: bytes
| country_id
: bytes
| sakila
: bytes
| country_id
< row
: bytes
| sakila
: bytes
| customer
: bytes
| fk_customer_address
: bytes
| address_id
: bytes
| sakila
: bytes
| address_id
< row
: bytes
| sakila
: bytes
| customer
: bytes
| fk_customer_store
: bytes
| store_id
: bytes
| sakila
: bytes
| store_id
< row
: bytes
| sakila
: bytes
| film
: bytes
| fk_film_language
: bytes
| language_id
: bytes
| sakila
: bytes
| language_id
< row
: bytes
| sakila
: bytes
| film
: bytes
| fk_film_language_original
: bytes
| original_language_id
: bytes
| sakila
: bytes
| language_id
< row
: bytes
| sakila
: bytes
| film_actor
: bytes
| fk_film_actor_actor
: bytes
| actor_id
: bytes
| sakila
: bytes
| actor_id
< row
: bytes
| sakila
: bytes
| film_actor
: bytes
| fk_film_actor_film
: bytes
| film_id
: bytes
| sakila
: bytes
| film_id
< row
: bytes
| sakila
: bytes
| film_category
: bytes
| fk_film_category_category
: bytes
| category_id
: bytes
| sakila
: bytes
| category_id
< row
: bytes
| sakila
: bytes
| film_category
: bytes
| fk_film_category_film
: bytes
| film_id
: bytes
| sakila
: bytes
| film_id
< row
: bytes
| sakila
: bytes
| inventory
: bytes
| fk_inventory_film
: bytes
| film_id
: bytes
| sakila
: bytes
| film_id
< row
: bytes
| sakila
: bytes
| inventory
: bytes
| fk_inventory_store
: bytes
| store_id
: bytes
| sakila
: bytes
| store_id
< row
: bytes
| sakila
: bytes
| payment
: bytes
| fk_payment_customer
: bytes
| customer_id
: bytes
| sakila
: bytes
| customer_id
< row
: bytes
| sakila
: bytes
| payment
: bytes
| fk_payment_rental
: bytes
| rental_id
: bytes
| sakila
: bytes
| rental_id
< row
: bytes
| sakila
: bytes
| payment
: bytes
| fk_payment_staff
: bytes
| staff_id
: bytes
| sakila
: bytes
| staff_id
< row
: bytes
| sakila
: bytes
| rental
: bytes
| fk_rental_customer
: bytes
| customer_id
: bytes
| sakila
: bytes
| customer_id
< row
: bytes
| sakila
: bytes
| rental
: bytes
| fk_rental_inventory
: bytes
| inventory_id
: bytes
| sakila
: bytes
| inventory_id
< row
: bytes
| sakila
: bytes
| rental
: bytes
| fk_rental_staff
: bytes
| staff_id
: bytes
| sakila
: bytes
| staff_id
< row
: bytes
| sakila
: bytes
| staff
: bytes
| fk_staff_address
: bytes
| address_id
: bytes
| sakila
: bytes
| address_id
< row
: bytes
| sakila
: bytes
| staff
: bytes
| fk_staff_store
: bytes
| store_id
: bytes
| sakila
: bytes
| store_id
< row
: bytes
| sakila
: bytes
| store
: bytes
| fk_store_address
: bytes
| address_id
: bytes
| sakila
: bytes
| address_id
< row
: bytes
| sakila
: bytes
| store
: bytes
| fk_store_staff
: bytes
| manager_staff_id
: bytes
| sakila
: bytes
| staff_id
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
| sakila
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
| sakila
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "ENUM" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT" "string"
< row
: bytes
| get_customer_balance
: bytes
| FUNCTION
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
|
< row
: bytes
| inventory_held_by_customer
: bytes
| FUNCTION
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
|
< row
: bytes
| inventory_in_stock
: bytes
| FUNCTION
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
|
< row
: bytes
| film_in_stock
: bytes
| PROCEDURE
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
|
< row
: bytes
| film_not_in_stock
: bytes
| PROCEDURE
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
|
< row
: bytes
| rewards_report
: bytes
| PROCEDURE
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| Provides a customizable report on best customers
> query
| SHOW CREATE FUNCTION `sakila`.`get_customer_balance`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| get_customer_balance
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `get_customer_balance`(p_customer_id INT, p_effective_date DATETIME) RETURNS decimal(5,2)
|     READS SQL DATA
|     DETERMINISTIC
| BEGIN
|
|        #OK, WE NEED TO CALCULATE THE CURRENT BALANCE GIVEN A CUSTOMER_ID AND A DATE
|        #THAT WE WANT THE BALANCE TO BE EFFECTIVE FOR. THE BALANCE IS:
|        #   1) RENTAL FEES FOR ALL PREVIOUS RENTALS
|        #   2) ONE DOLLAR FOR EVERY DAY THE PREVIOUS RENTALS ARE OVERDUE
|        #   3) IF A FILM IS MORE THAN RENTAL_DURATION * 2 OVERDUE, CHARGE THE REPLACEMENT_COST
|        #   4) SUBTRACT ALL PAYMENTS MADE BEFORE THE DATE SPECIFIED
|
|   DECLARE v_rentfees DECIMAL(5,2); #FEES PAID TO RENT THE VIDEOS INITIALLY
|   DECLARE v_overfees INTEGER;      #LATE FEES FOR PRIOR RENTALS
|   DECLARE v_payments DECIMAL(5,2); #SUM OF PAYMENTS MADE PREVIOUSLY
|
|   SELECT IFNULL(SUM(film.rental_rate),0) INTO v_rentfees
|     FROM film, inventory, rental
|     WHERE film.film_id = inventory.film_id
|       AND inventory.inventory_id = rental.inventory_id
|       AND rental.rental_date <= p_effective_date
|       AND rental.customer_id = p_customer_id;
|
|   SELECT IFNULL(SUM(IF((TO_DAYS(rental.return_date) - TO_DAYS(rental.rental_date)) > film.rental_duration,
|         ((TO_DAYS(rental.return_date) - TO_DAYS(rental.rental_date)) - film.rental_duration),0)),0) INTO v_overfees
|     FROM rental, inventory, film
|     WHERE film.film_id = inventory.film_id
|       AND inventory.inventory_id = rental.inventory_id
|       AND rental.rental_date <= p_effective_date
|       AND rental.customer_id = p_customer_id;
|
|
|   SELECT IFNULL(SUM(payment.amount),0) INTO v_payments
|     FROM payment
|
|     WHERE payment.payment_date <= p_effective_date
|     AND payment.customer_id = p_customer_id;
|
|   RETURN v_rentfees + v_overfees - v_payments;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE FUNCTION `sakila`.`inventory_held_by_customer`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| inventory_held_by_customer
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `inventory_held_by_customer`(p_inventory_id INT) RETURNS int
|     READS SQL DATA
| BEGIN
|   DECLARE v_customer_id INT;
|   DECLARE EXIT HANDLER FOR NOT FOUND RETURN NULL;
|
|   SELECT customer_id INTO v_customer_id
|   FROM rental
|   WHERE return_date IS NULL
|   AND inventory_id = p_inventory_id;
|
|   RETURN v_customer_id;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE FUNCTION `sakila`.`inventory_in_stock`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| inventory_in_stock
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `inventory_in_stock`(p_inventory_id INT) RETURNS tinyint(1)
|     READS SQL DATA
| BEGIN
|     DECLARE v_rentals INT;
|     DECLARE v_out     INT;
|
|     #AN ITEM IS IN-STOCK IF THERE ARE EITHER NO ROWS IN THE rental TABLE
|     #FOR THE ITEM OR ALL ROWS HAVE return_date POPULATED
|
|     SELECT COUNT(*) INTO v_rentals
|     FROM rental
|     WHERE inventory_id = p_inventory_id;
|
|     IF v_rentals = 0 THEN
|       RETURN TRUE;
|     END IF;
|
|     SELECT COUNT(rental_id) INTO v_out
|     FROM inventory LEFT JOIN rental USING(inventory_id)
|     WHERE inventory.inventory_id = p_inventory_id
|     AND rental.return_date IS NULL;
|
|     IF v_out > 0 THEN
|       RETURN FALSE;
|     ELSE
|       RETURN TRUE;
|     END IF;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `sakila`.`film_in_stock`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| film_in_stock
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `film_in_stock`(IN p_film_id INT, IN p_store_id INT, OUT p_film_count INT)
|     READS SQL DATA
| BEGIN
|      SELECT inventory_id
|      FROM inventory
|      WHERE film_id = p_film_id
|      AND store_id = p_store_id
|      AND inventory_in_stock(inventory_id);
|
|      SELECT FOUND_ROWS() INTO p_film_count;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `sakila`.`film_not_in_stock`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| film_not_in_stock
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `film_not_in_stock`(IN p_film_id INT, IN p_store_id INT, OUT p_film_count INT)
|     READS SQL DATA
| BEGIN
|      SELECT inventory_id
|      FROM inventory
|      WHERE film_id = p_film_id
|      AND store_id = p_store_id
|      AND NOT inventory_in_stock(inventory_id);
|
|      SELECT FOUND_ROWS() INTO p_film_count;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `sakila`.`rewards_report`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| rewards_report
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `rewards_report`(
|     IN min_monthly_purchases TINYINT UNSIGNED
|     , IN min_dollar_amount_purchased DECIMAL(10,2) UNSIGNED
|     , OUT count_rewardees INT
| )
|     READS SQL DATA
|     COMMENT 'Provides a customizable report on best customers'
| proc: BEGIN
|
|     DECLARE last_month_start DATE;
|     DECLARE last_month_end DATE;
|
|     /* Some sanity checks... */
|     IF min_monthly_purchases = 0 THEN
|         SELECT 'Minimum monthly purchases parameter must be > 0';
|         LEAVE proc;
|     END IF;
|     IF min_dollar_amount_purchased = 0.00 THEN
|         SELECT 'Minimum monthly dollar amount purchased parameter must be > $0.00';
|         LEAVE proc;
|     END IF;
|
|     /* Determine start and end time periods */
|     SET last_month_start = DATE_SUB(CURRENT_DATE(), INTERVAL 1 MONTH);
|     SET last_month_start = STR_TO_DATE(CONCAT(YEAR(last_month_start),'-',MONTH(last_month_start),'-01'),'%Y-%m-%d');
|     SET last_month_end = LAST_DAY(last_month_start);
|
|     /*
|         Create a temporary storage area for
|         Customer IDs.
|     */
|     CREATE TEMPORARY TABLE tmpCustomer (customer_id INT UNSIGNED NOT NULL PRIMARY KEY);
|
|     /*
|         Find all customers meeting the
|         monthly purchase requirements
|     */
|     INSERT INTO tmpCustomer (customer_id)
|     SELECT p.customer_id
|     FROM payment AS p
|     WHERE DATE(p.payment_date) BETWEEN last_month_start AND last_month_end
|     GROUP BY customer_id
|     HAVING SUM(p.amount) > min_dollar_amount_purchased
|     AND COUNT(customer_id) > min_monthly_purchases;
|
|     /* Populate OUT parameter with count of found customers */
|     SELECT COUNT(*) FROM tmpCustomer INTO count_rewardees;
|
|     /*
|         Output ALL customer information of matching rewardees.
|         Customize output as needed.
|     */
|     SELECT c.*
|     FROM tmpCustomer AS t
|     INNER JOIN customer AS c ON t.customer_id = c.customer_id;
|
|     /* Clean up */
|     DROP TABLE tmpCustomer;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
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
| sakila
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
| actor
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| actor_info
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| address
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| category
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| city
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| country
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| customer
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
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
| utf8mb3
< row
: bytes
| customer_list
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| film
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
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
| utf8mb3
< row
: bytes
| film_actor
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| film_category
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| film_list
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| film_text
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 32768
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| inventory
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 32768
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| language
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| nicer_but_slower_film_list
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| payment
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
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
| utf8mb3
< row
: bytes
| rental
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 65536
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| sales_by_film_category
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| sales_by_store
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| staff
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 32768
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| staff_list
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| store
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 32768
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
> query
|
| 		SELECT t.NAME
| 		FROM information_schema.INNODB_TABLES t
| 		LEFT JOIN information_schema.INNODB_TABLESPACES ts ON t.SPACE = ts.SPACE
| 		WHERE t.NAME LIKE CONCAT(?, '/%') AND (ts.SPACE_TYPE IS NULL OR ts.SPACE_TYPE <> 'Single')
: string
| sakila
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
| sakila
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci

== TestSDLRealWorldIdempotence default#10
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
| sdl_rw_mediawiki_4
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
| actor
: bytes
| actor_name
: bytes
| actor_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| actor
: bytes
| actor_user
: bytes
| actor_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| actor
: bytes
| PRIMARY
: bytes
| actor_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| archive
: bytes
| ar_actor_timestamp
: bytes
| ar_actor
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| archive
: bytes
| ar_actor_timestamp
: bytes
| ar_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| archive
: bytes
| ar_revid
: bytes
| ar_rev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| archive
: bytes
| ar_usertext_timestamp
: bytes
| ar_user_text
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| archive
: bytes
| ar_usertext_timestamp
: bytes
| ar_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| archive
: bytes
| name_title_timestamp
: bytes
| ar_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| archive
: bytes
| name_title_timestamp
: bytes
| ar_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| archive
: bytes
| name_title_timestamp
: bytes
| ar_timestamp
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| archive
: bytes
| PRIMARY
: bytes
| ar_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| bot_passwords
: bytes
| PRIMARY
: bytes
| bp_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| bot_passwords
: bytes
| PRIMARY
: bytes
| bp_app_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| category
: bytes
| cat_pages
: bytes
| cat_pages
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| category
: bytes
| cat_title
: bytes
| cat_title
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| category
: bytes
| PRIMARY
: bytes
| cat_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| cl_collation_ext
: bytes
| cl_collation
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| cl_collation_ext
: bytes
| cl_to
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| cl_collation_ext
: bytes
| cl_type
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| cl_collation_ext
: bytes
| cl_from
: bytes
| A
: int64 -1
: null
: int64 4
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| cl_sortkey
: bytes
| cl_to
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| cl_sortkey
: bytes
| cl_type
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| cl_sortkey
: bytes
| cl_sortkey
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| cl_sortkey
: bytes
| cl_from
: bytes
| A
: int64 -1
: null
: int64 4
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| cl_timestamp
: bytes
| cl_to
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| cl_timestamp
: bytes
| cl_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| PRIMARY
: bytes
| cl_from
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| PRIMARY
: bytes
| cl_to
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| change_tag_log_tag
: bytes
| ct_log_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| change_tag_log_tag
: bytes
| ct_tag
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| change_tag_rc_tag
: bytes
| ct_rc_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| change_tag_rc_tag
: bytes
| ct_tag
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| change_tag_rev_tag
: bytes
| ct_rev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| change_tag_rev_tag
: bytes
| ct_tag
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| change_tag_tag_id
: bytes
| ct_tag
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| change_tag_tag_id
: bytes
| ct_rc_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| change_tag_tag_id
: bytes
| ct_rev_id
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| change_tag_tag_id
: bytes
| ct_log_id
: bytes
| A
: int64 -1
: null
: int64 4
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| PRIMARY
: bytes
| ct_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| comment
: bytes
| comment_hash
: bytes
| comment_hash
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| comment
: bytes
| PRIMARY
: bytes
| comment_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| content
: bytes
| PRIMARY
: bytes
| content_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| content_models
: bytes
| model_name
: bytes
| model_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| content_models
: bytes
| PRIMARY
: bytes
| model_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| el_from
: bytes
| el_from
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| el_from
: bytes
| el_to
: bytes
| A
: int64 40
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| el_from_index_60
: bytes
| el_from
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| el_from_index_60
: bytes
| el_index_60
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| el_from_index_60
: bytes
| el_id
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| el_index
: bytes
| el_index
: bytes
| A
: int64 60
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| el_index_60
: bytes
| el_index_60
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| el_index_60
: bytes
| el_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| el_to
: bytes
| el_to
: bytes
| A
: int64 60
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| el_to
: bytes
| el_from
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| PRIMARY
: bytes
| el_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| fa_actor_timestamp
: bytes
| fa_actor
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| fa_actor_timestamp
: bytes
| fa_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| fa_deleted_timestamp
: bytes
| fa_deleted_timestamp
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| fa_name
: bytes
| fa_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| fa_name
: bytes
| fa_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| fa_sha1
: bytes
| fa_sha1
: bytes
| A
: int64 10
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| fa_storage_group
: bytes
| fa_storage_group
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| fa_storage_group
: bytes
| fa_storage_key
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| fa_user_timestamp
: bytes
| fa_user_text
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| fa_user_timestamp
: bytes
| fa_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| PRIMARY
: bytes
| fa_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_actor_timestamp
: bytes
| img_actor
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_actor_timestamp
: bytes
| img_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_media_mime
: bytes
| img_media_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_media_mime
: bytes
| img_major_mime
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_media_mime
: bytes
| img_minor_mime
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_sha1
: bytes
| img_sha1
: bytes
| A
: int64 10
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_size
: bytes
| img_size
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_timestamp
: bytes
| img_timestamp
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_user_timestamp
: bytes
| img_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_user_timestamp
: bytes
| img_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_usertext_timestamp
: bytes
| img_user_text
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_usertext_timestamp
: bytes
| img_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| PRIMARY
: bytes
| img_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| image_comment_temp
: bytes
| imgcomment_name
: bytes
| imgcomment_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| image_comment_temp
: bytes
| PRIMARY
: bytes
| imgcomment_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| image_comment_temp
: bytes
| PRIMARY
: bytes
| imgcomment_description_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| imagelinks
: bytes
| il_backlinks_namespace
: bytes
| il_from_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| imagelinks
: bytes
| il_backlinks_namespace
: bytes
| il_to
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| imagelinks
: bytes
| il_backlinks_namespace
: bytes
| il_from
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| imagelinks
: bytes
| il_to
: bytes
| il_to
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| imagelinks
: bytes
| il_to
: bytes
| il_from
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| imagelinks
: bytes
| PRIMARY
: bytes
| il_from
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| imagelinks
: bytes
| PRIMARY
: bytes
| il_to
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| interwiki
: bytes
| iw_prefix
: bytes
| iw_prefix
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| ip_changes
: bytes
| ipc_hex_time
: bytes
| ipc_hex
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| ip_changes
: bytes
| ipc_hex_time
: bytes
| ipc_rev_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| ip_changes
: bytes
| ipc_rev_timestamp
: bytes
| ipc_rev_timestamp
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| ip_changes
: bytes
| PRIMARY
: bytes
| ipc_rev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| ipb_address
: bytes
| ipb_address
: bytes
| A
: int64 255
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| ipb_address
: bytes
| ipb_user
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| ipb_address
: bytes
| ipb_auto
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| ipb_address
: bytes
| ipb_anon_only
: bytes
| A
: int64 -1
: null
: int64 4
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| ipb_expiry
: bytes
| ipb_expiry
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| ipb_parent_block_id
: bytes
| ipb_parent_block_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| ipb_range
: bytes
| ipb_range_start
: bytes
| A
: int64 8
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| ipb_range
: bytes
| ipb_range_end
: bytes
| A
: int64 8
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| ipb_timestamp
: bytes
| ipb_timestamp
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| ipb_user
: bytes
| ipb_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| PRIMARY
: bytes
| ipb_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| iwlinks
: bytes
| iwl_prefix_from_title
: bytes
| iwl_prefix
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| iwlinks
: bytes
| iwl_prefix_from_title
: bytes
| iwl_from
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| iwlinks
: bytes
| iwl_prefix_from_title
: bytes
| iwl_title
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| iwlinks
: bytes
| iwl_prefix_title_from
: bytes
| iwl_prefix
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| iwlinks
: bytes
| iwl_prefix_title_from
: bytes
| iwl_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| iwlinks
: bytes
| iwl_prefix_title_from
: bytes
| iwl_from
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| iwlinks
: bytes
| PRIMARY
: bytes
| iwl_from
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| iwlinks
: bytes
| PRIMARY
: bytes
| iwl_prefix
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| iwlinks
: bytes
| PRIMARY
: bytes
| iwl_title
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_cmd
: bytes
| job_cmd
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_cmd
: bytes
| job_namespace
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_cmd
: bytes
| job_title
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_cmd
: bytes
| job_params
: bytes
| A
: int64 128
: null
: int64 4
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_cmd_token
: bytes
| job_cmd
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_cmd_token
: bytes
| job_token
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_cmd_token
: bytes
| job_random
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_cmd_token_id
: bytes
| job_cmd
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_cmd_token_id
: bytes
| job_token
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_cmd_token_id
: bytes
| job_id
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_sha1
: bytes
| job_sha1
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_timestamp
: bytes
| job_timestamp
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| PRIMARY
: bytes
| job_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| l10n_cache
: bytes
| PRIMARY
: bytes
| lc_lang
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| l10n_cache
: bytes
| PRIMARY
: bytes
| lc_key
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| langlinks
: bytes
| ll_lang
: bytes
| ll_lang
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| langlinks
: bytes
| ll_lang
: bytes
| ll_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| langlinks
: bytes
| PRIMARY
: bytes
| ll_from
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| langlinks
: bytes
| PRIMARY
: bytes
| ll_lang
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| log_search
: bytes
| ls_log_id
: bytes
| ls_log_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| log_search
: bytes
| PRIMARY
: bytes
| ls_field
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| log_search
: bytes
| PRIMARY
: bytes
| ls_value
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| log_search
: bytes
| PRIMARY
: bytes
| ls_log_id
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| actor_time
: bytes
| log_actor
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| actor_time
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_actor_type_time
: bytes
| log_actor
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_actor_type_time
: bytes
| log_type
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_actor_type_time
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_page_id_time
: bytes
| log_page
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_page_id_time
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_user_text_time
: bytes
| log_user_text
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_user_text_time
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_user_text_type_time
: bytes
| log_user_text
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_user_text_type_time
: bytes
| log_type
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_user_text_type_time
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_user_type_time
: bytes
| log_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_user_type_time
: bytes
| log_type
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_user_type_time
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| page_time
: bytes
| log_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| page_time
: bytes
| log_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| page_time
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| PRIMARY
: bytes
| log_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| times
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| type_action
: bytes
| log_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| type_action
: bytes
| log_action
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| type_action
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| type_time
: bytes
| log_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| type_time
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| user_time
: bytes
| log_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| user_time
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| module_deps
: bytes
| PRIMARY
: bytes
| md_module
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| module_deps
: bytes
| PRIMARY
: bytes
| md_skin
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| objectcache
: bytes
| exptime
: bytes
| exptime
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| objectcache
: bytes
| PRIMARY
: bytes
| keyname
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| oldimage
: bytes
| oi_actor_timestamp
: bytes
| oi_actor
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| oldimage
: bytes
| oi_actor_timestamp
: bytes
| oi_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| oldimage
: bytes
| oi_name_archive_name
: bytes
| oi_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| oldimage
: bytes
| oi_name_archive_name
: bytes
| oi_archive_name
: bytes
| A
: int64 14
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| oldimage
: bytes
| oi_name_timestamp
: bytes
| oi_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| oldimage
: bytes
| oi_name_timestamp
: bytes
| oi_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| oldimage
: bytes
| oi_sha1
: bytes
| oi_sha1
: bytes
| A
: int64 10
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| oldimage
: bytes
| oi_usertext_timestamp
: bytes
| oi_user_text
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| oldimage
: bytes
| oi_usertext_timestamp
: bytes
| oi_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| page
: bytes
| name_title
: bytes
| page_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page
: bytes
| name_title
: bytes
| page_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page
: bytes
| page_len
: bytes
| page_len
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| page
: bytes
| page_random
: bytes
| page_random
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| page
: bytes
| page_redirect_namespace_len
: bytes
| page_is_redirect
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| page
: bytes
| page_redirect_namespace_len
: bytes
| page_namespace
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| page
: bytes
| page_redirect_namespace_len
: bytes
| page_len
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| page
: bytes
| PRIMARY
: bytes
| page_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page_props
: bytes
| pp_page_propname
: bytes
| pp_page
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page_props
: bytes
| pp_page_propname
: bytes
| pp_propname
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page_props
: bytes
| pp_propname_page
: bytes
| pp_propname
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page_props
: bytes
| pp_propname_page
: bytes
| pp_page
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page_props
: bytes
| pp_propname_sortkey_page
: bytes
| pp_propname
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page_props
: bytes
| pp_propname_sortkey_page
: bytes
| pp_sortkey
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page_props
: bytes
| pp_propname_sortkey_page
: bytes
| pp_page
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page_restrictions
: bytes
| pr_cascade
: bytes
| pr_cascade
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| page_restrictions
: bytes
| pr_level
: bytes
| pr_level
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| page_restrictions
: bytes
| pr_pagetype
: bytes
| pr_page
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page_restrictions
: bytes
| pr_pagetype
: bytes
| pr_type
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page_restrictions
: bytes
| pr_typelevel
: bytes
| pr_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| page_restrictions
: bytes
| pr_typelevel
: bytes
| pr_level
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| page_restrictions
: bytes
| PRIMARY
: bytes
| pr_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| pagelinks
: bytes
| pl_backlinks_namespace
: bytes
| pl_from_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| pagelinks
: bytes
| pl_backlinks_namespace
: bytes
| pl_namespace
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| pagelinks
: bytes
| pl_backlinks_namespace
: bytes
| pl_title
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| pagelinks
: bytes
| pl_backlinks_namespace
: bytes
| pl_from
: bytes
| A
: int64 -1
: null
: int64 4
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| pagelinks
: bytes
| pl_namespace
: bytes
| pl_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| pagelinks
: bytes
| pl_namespace
: bytes
| pl_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| pagelinks
: bytes
| pl_namespace
: bytes
| pl_from
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| pagelinks
: bytes
| PRIMARY
: bytes
| pl_from
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| pagelinks
: bytes
| PRIMARY
: bytes
| pl_namespace
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| pagelinks
: bytes
| PRIMARY
: bytes
| pl_title
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| protected_titles
: bytes
| pt_namespace_title
: bytes
| pt_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| protected_titles
: bytes
| pt_namespace_title
: bytes
| pt_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| protected_titles
: bytes
| pt_timestamp
: bytes
| pt_timestamp
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| querycache
: bytes
| qc_type
: bytes
| qc_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| querycache
: bytes
| qc_type
: bytes
| qc_value
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| querycache_info
: bytes
| PRIMARY
: bytes
| qci_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| querycachetwo
: bytes
| qcc_title
: bytes
| qcc_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| querycachetwo
: bytes
| qcc_title
: bytes
| qcc_namespace
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| querycachetwo
: bytes
| qcc_title
: bytes
| qcc_title
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| querycachetwo
: bytes
| qcc_titletwo
: bytes
| qcc_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| querycachetwo
: bytes
| qcc_titletwo
: bytes
| qcc_namespacetwo
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| querycachetwo
: bytes
| qcc_titletwo
: bytes
| qcc_titletwo
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| querycachetwo
: bytes
| qcc_type
: bytes
| qcc_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| querycachetwo
: bytes
| qcc_type
: bytes
| qcc_value
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| new_name_timestamp
: bytes
| rc_new
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| new_name_timestamp
: bytes
| rc_namespace
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| new_name_timestamp
: bytes
| rc_timestamp
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| PRIMARY
: bytes
| rc_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_actor
: bytes
| rc_actor
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_actor
: bytes
| rc_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_cur_id
: bytes
| rc_cur_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_ip
: bytes
| rc_ip
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_name_type_patrolled_timestamp
: bytes
| rc_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_name_type_patrolled_timestamp
: bytes
| rc_type
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_name_type_patrolled_timestamp
: bytes
| rc_patrolled
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_name_type_patrolled_timestamp
: bytes
| rc_timestamp
: bytes
| A
: int64 -1
: null
: int64 4
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_namespace_title_timestamp
: bytes
| rc_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_namespace_title_timestamp
: bytes
| rc_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_namespace_title_timestamp
: bytes
| rc_timestamp
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_ns_actor
: bytes
| rc_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_ns_actor
: bytes
| rc_actor
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_ns_usertext
: bytes
| rc_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_ns_usertext
: bytes
| rc_user_text
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_timestamp
: bytes
| rc_timestamp
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_user_text
: bytes
| rc_user_text
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_user_text
: bytes
| rc_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| redirect
: bytes
| PRIMARY
: bytes
| rd_from
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| redirect
: bytes
| rd_ns_title
: bytes
| rd_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| redirect
: bytes
| rd_ns_title
: bytes
| rd_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| redirect
: bytes
| rd_ns_title
: bytes
| rd_from
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| page_timestamp
: bytes
| rev_page
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| page_timestamp
: bytes
| rev_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| page_user_timestamp
: bytes
| rev_page
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| page_user_timestamp
: bytes
| rev_user
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| page_user_timestamp
: bytes
| rev_timestamp
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| PRIMARY
: bytes
| rev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| rev_page_id
: bytes
| rev_page
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| rev_page_id
: bytes
| rev_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| rev_timestamp
: bytes
| rev_timestamp
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| user_timestamp
: bytes
| rev_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| user_timestamp
: bytes
| rev_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| usertext_timestamp
: bytes
| rev_user_text
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| usertext_timestamp
: bytes
| rev_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision_actor_temp
: bytes
| actor_timestamp
: bytes
| revactor_actor
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision_actor_temp
: bytes
| actor_timestamp
: bytes
| revactor_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision_actor_temp
: bytes
| page_actor_timestamp
: bytes
| revactor_page
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision_actor_temp
: bytes
| page_actor_timestamp
: bytes
| revactor_actor
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision_actor_temp
: bytes
| page_actor_timestamp
: bytes
| revactor_timestamp
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision_actor_temp
: bytes
| PRIMARY
: bytes
| revactor_rev
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| revision_actor_temp
: bytes
| PRIMARY
: bytes
| revactor_actor
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| revision_actor_temp
: bytes
| revactor_rev
: bytes
| revactor_rev
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| revision_comment_temp
: bytes
| PRIMARY
: bytes
| revcomment_rev
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| revision_comment_temp
: bytes
| PRIMARY
: bytes
| revcomment_comment_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| revision_comment_temp
: bytes
| revcomment_rev
: bytes
| revcomment_rev
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| searchindex
: bytes
| si_page
: bytes
| si_page
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| searchindex
: bytes
| si_text
: bytes
| si_text
: null
: int64 -1
: null
: int64 1
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
< row
: bytes
| searchindex
: bytes
| si_title
: bytes
| si_title
: null
: int64 -1
: null
: int64 1
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
< row
: bytes
| site_identifiers
: bytes
| site_ids_key
: bytes
| si_key
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| site_identifiers
: bytes
| site_ids_site
: bytes
| si_site
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| site_identifiers
: bytes
| site_ids_type
: bytes
| si_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| site_identifiers
: bytes
| site_ids_type
: bytes
| si_key
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| site_stats
: bytes
| PRIMARY
: bytes
| ss_row_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| sites
: bytes
| PRIMARY
: bytes
| site_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| sites
: bytes
| sites_domain
: bytes
| site_domain
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| sites
: bytes
| sites_forward
: bytes
| site_forward
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| sites
: bytes
| sites_global_key
: bytes
| site_global_key
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| sites
: bytes
| sites_group
: bytes
| site_group
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| sites
: bytes
| sites_language
: bytes
| site_language
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| sites
: bytes
| sites_protocol
: bytes
| site_protocol
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| sites
: bytes
| sites_source
: bytes
| site_source
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| sites
: bytes
| sites_type
: bytes
| site_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| slot_roles
: bytes
| PRIMARY
: bytes
| role_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| slot_roles
: bytes
| role_name
: bytes
| role_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| slots
: bytes
| PRIMARY
: bytes
| slot_revision_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| slots
: bytes
| PRIMARY
: bytes
| slot_role_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| slots
: bytes
| slot_revision_origin_role
: bytes
| slot_revision_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| slots
: bytes
| slot_revision_origin_role
: bytes
| slot_origin
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| slots
: bytes
| slot_revision_origin_role
: bytes
| slot_role_id
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| tag_summary
: bytes
| PRIMARY
: bytes
| ts_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| tag_summary
: bytes
| tag_summary_log_id
: bytes
| ts_log_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| tag_summary
: bytes
| tag_summary_rc_id
: bytes
| ts_rc_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| tag_summary
: bytes
| tag_summary_rev_id
: bytes
| ts_rev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| templatelinks
: bytes
| PRIMARY
: bytes
| tl_from
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| templatelinks
: bytes
| PRIMARY
: bytes
| tl_namespace
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| templatelinks
: bytes
| PRIMARY
: bytes
| tl_title
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| templatelinks
: bytes
| tl_backlinks_namespace
: bytes
| tl_from_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| templatelinks
: bytes
| tl_backlinks_namespace
: bytes
| tl_namespace
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| templatelinks
: bytes
| tl_backlinks_namespace
: bytes
| tl_title
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| templatelinks
: bytes
| tl_backlinks_namespace
: bytes
| tl_from
: bytes
| A
: int64 -1
: null
: int64 4
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| templatelinks
: bytes
| tl_namespace
: bytes
| tl_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| templatelinks
: bytes
| tl_namespace
: bytes
| tl_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| templatelinks
: bytes
| tl_namespace
: bytes
| tl_from
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| text
: bytes
| PRIMARY
: bytes
| old_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| transcache
: bytes
| PRIMARY
: bytes
| tc_url
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| updatelog
: bytes
| PRIMARY
: bytes
| ul_key
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| uploadstash
: bytes
| PRIMARY
: bytes
| us_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| uploadstash
: bytes
| us_key
: bytes
| us_key
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| uploadstash
: bytes
| us_timestamp
: bytes
| us_timestamp
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| uploadstash
: bytes
| us_user
: bytes
| us_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| user
: bytes
| PRIMARY
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| user
: bytes
| user_email
: bytes
| user_email
: bytes
| A
: int64 50
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| user
: bytes
| user_email_token
: bytes
| user_email_token
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| user
: bytes
| user_name
: bytes
| user_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| user_former_groups
: bytes
| PRIMARY
: bytes
| ufg_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| user_former_groups
: bytes
| PRIMARY
: bytes
| ufg_group
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| user_groups
: bytes
| PRIMARY
: bytes
| ug_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| user_groups
: bytes
| PRIMARY
: bytes
| ug_group
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| user_groups
: bytes
| ug_expiry
: bytes
| ug_expiry
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| user_groups
: bytes
| ug_group
: bytes
| ug_group
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| user_newtalk
: bytes
| un_user_id
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| user_newtalk
: bytes
| un_user_ip
: bytes
| user_ip
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| user_properties
: bytes
| PRIMARY
: bytes
| up_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| user_properties
: bytes
| PRIMARY
: bytes
| up_property
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| user_properties
: bytes
| user_properties_property
: bytes
| up_property
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| valid_tag
: bytes
| PRIMARY
: bytes
| vt_tag
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| watchlist
: bytes
| namespace_title
: bytes
| wl_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| watchlist
: bytes
| namespace_title
: bytes
| wl_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| watchlist
: bytes
| PRIMARY
: bytes
| wl_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| watchlist
: bytes
| wl_user
: bytes
| wl_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| watchlist
: bytes
| wl_user
: bytes
| wl_namespace
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| watchlist
: bytes
| wl_user
: bytes
| wl_title
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| watchlist
: bytes
| wl_user_notificationtimestamp
: bytes
| wl_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| watchlist
: bytes
| wl_user_notificationtimestamp
: bytes
| wl_notificationtimestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
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
| sdl_rw_mediawiki_4
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
| actor
: bytes
| actor_id
: int64 1
: null
: bytes
| NO
: bytes
| bigint unsigned
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
| actor
: bytes
| actor_user
: int64 2
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| actor
: bytes
| actor_name
: int64 3
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| archive
: bytes
| ar_namespace
: int64 2
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
| archive
: bytes
| ar_title
: int64 3
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_comment
: int64 4
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(767)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_comment_id
: int64 5
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_user
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_user_text
: int64 7
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_actor
: int64 8
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_timestamp
: int64 9
: bytes
| '0x'
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_minor_edit
: int64 10
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_rev_id
: int64 11
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_text_id
: int64 12
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_deleted
: int64 13
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_len
: int64 14
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_page_id
: int64 15
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_parent_id
: int64 16
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_sha1
: int64 17
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_content_model
: int64 18
: null
: bytes
| YES
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_content_format
: int64 19
: null
: bytes
| YES
: bytes
| varbinary(64)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| bot_passwords
: bytes
| bp_user
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| bot_passwords
: bytes
| bp_app_id
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| bot_passwords
: bytes
| bp_password
: int64 3
: null
: bytes
| NO
: bytes
| tinyblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| bot_passwords
: bytes
| bp_token
: int64 4
: bytes
| '0x'
: bytes
| NO
: bytes
| binary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| bot_passwords
: bytes
| bp_restrictions
: int64 5
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| bot_passwords
: bytes
| bp_grants
: int64 6
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| category
: bytes
| cat_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| category
: bytes
| cat_title
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| category
: bytes
| cat_pages
: int64 3
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
| category
: bytes
| cat_subcats
: int64 4
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
| category
: bytes
| cat_files
: int64 5
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
| categorylinks
: bytes
| cl_from
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| categorylinks
: bytes
| cl_to
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| categorylinks
: bytes
| cl_sortkey
: int64 3
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(230)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| categorylinks
: bytes
| cl_sortkey_prefix
: int64 4
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| categorylinks
: bytes
| cl_timestamp
: int64 5
: null
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| categorylinks
: bytes
| cl_collation
: int64 6
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| categorylinks
: bytes
| cl_type
: int64 7
: bytes
| '0x70616765'
: bytes
| NO
: bytes
| enum('page','subcat','file')
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| change_tag
: bytes
| ct_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| change_tag
: bytes
| ct_rc_id
: int64 2
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
| change_tag
: bytes
| ct_log_id
: int64 3
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| change_tag
: bytes
| ct_rev_id
: int64 4
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| change_tag
: bytes
| ct_tag
: int64 5
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| change_tag
: bytes
| ct_params
: int64 6
: null
: bytes
| YES
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| comment
: bytes
| comment_id
: int64 1
: null
: bytes
| NO
: bytes
| bigint unsigned
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
| comment
: bytes
| comment_hash
: int64 2
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
|
: null
< row
: bytes
| comment
: bytes
| comment_text
: int64 3
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| comment
: bytes
| comment_data
: int64 4
: null
: bytes
| YES
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| content
: bytes
| content_id
: int64 1
: null
: bytes
| NO
: bytes
| bigint unsigned
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
| content
: bytes
| content_size
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| content
: bytes
| content_sha1
: int64 3
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| content
: bytes
| content_model
: int64 4
: null
: bytes
| NO
: bytes
| smallint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| content
: bytes
| content_address
: int64 5
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| content_models
: bytes
| model_id
: int64 1
: null
: bytes
| NO
: bytes
| smallint
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
| content_models
: bytes
| model_name
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(64)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| externallinks
: bytes
| el_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| externallinks
: bytes
| el_from
: int64 2
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| externallinks
: bytes
| el_to
: int64 3
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| externallinks
: bytes
| el_index
: int64 4
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| externallinks
: bytes
| el_index_60
: int64 5
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(60)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_id
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
| filearchive
: bytes
| fa_name
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_archive_name
: int64 3
: bytes
| ''
: bytes
| YES
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_storage_group
: int64 4
: null
: bytes
| YES
: bytes
| varbinary(16)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_storage_key
: int64 5
: bytes
| ''
: bytes
| YES
: bytes
| varbinary(64)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_deleted_user
: int64 6
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
| filearchive
: bytes
| fa_deleted_timestamp
: int64 7
: bytes
| '0x'
: bytes
| YES
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_deleted_reason
: int64 8
: bytes
| ''
: bytes
| YES
: bytes
| varbinary(767)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_deleted_reason_id
: int64 9
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_size
: int64 10
: bytes
| '0'
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_width
: int64 11
: bytes
| '0'
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
| filearchive
: bytes
| fa_height
: int64 12
: bytes
| '0'
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
| filearchive
: bytes
| fa_metadata
: int64 13
: null
: bytes
| YES
: bytes
| mediumblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_bits
: int64 14
: bytes
| '0'
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
| filearchive
: bytes
| fa_media_type
: int64 15
: null
: bytes
| YES
: bytes
| enum('UNKNOWN','BITMAP','DRAWING','AUDIO','VIDEO','MULTIMEDIA','OFFICE','TEXT','EXECUTABLE','ARCHIVE','3D')
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_major_mime
: int64 16
: bytes
| '0x756E6B6E6F776E'
: bytes
| YES
: bytes
| enum('unknown','application','audio','image','text','video','message','model','multipart','chemical')
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_minor_mime
: int64 17
: bytes
| '0x756E6B6E6F776E'
: bytes
| YES
: bytes
| varbinary(100)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_description
: int64 18
: bytes
| ''
: bytes
| YES
: bytes
| varbinary(767)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_description_id
: int64 19
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_user
: int64 20
: bytes
| '0'
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_user_text
: int64 21
: bytes
| ''
: bytes
| YES
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_actor
: int64 22
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_timestamp
: int64 23
: bytes
| '0x'
: bytes
| YES
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_deleted
: int64 24
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_sha1
: int64 25
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_name
: int64 1
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_size
: int64 2
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_width
: int64 3
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
| image
: bytes
| img_height
: int64 4
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
| image
: bytes
| img_metadata
: int64 5
: null
: bytes
| NO
: bytes
| mediumblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_bits
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
| image
: bytes
| img_media_type
: int64 7
: null
: bytes
| YES
: bytes
| enum('UNKNOWN','BITMAP','DRAWING','AUDIO','VIDEO','MULTIMEDIA','OFFICE','TEXT','EXECUTABLE','ARCHIVE','3D')
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_major_mime
: int64 8
: bytes
| '0x756E6B6E6F776E'
: bytes
| NO
: bytes
| enum('unknown','application','audio','image','text','video','message','model','multipart','chemical')
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_minor_mime
: int64 9
: bytes
| '0x756E6B6E6F776E'
: bytes
| NO
: bytes
| varbinary(100)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_description
: int64 10
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(767)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_description_id
: int64 11
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_user
: int64 12
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_user_text
: int64 13
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_actor
: int64 14
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_timestamp
: int64 15
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_sha1
: int64 16
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image_comment_temp
: bytes
| imgcomment_name
: int64 1
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image_comment_temp
: bytes
| imgcomment_description_id
: int64 2
: null
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| imagelinks
: bytes
| il_from
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| imagelinks
: bytes
| il_from_namespace
: int64 2
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
| imagelinks
: bytes
| il_to
: int64 3
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| interwiki
: bytes
| iw_prefix
: int64 1
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| interwiki
: bytes
| iw_url
: int64 2
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| interwiki
: bytes
| iw_api
: int64 3
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| interwiki
: bytes
| iw_wikiid
: int64 4
: null
: bytes
| NO
: bytes
| varbinary(64)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| interwiki
: bytes
| iw_local
: int64 5
: null
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
| interwiki
: bytes
| iw_trans
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ip_changes
: bytes
| ipc_rev_id
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ip_changes
: bytes
| ipc_rev_timestamp
: int64 2
: bytes
| '0x'
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ip_changes
: bytes
| ipc_hex
: int64 3
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(35)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_id
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
| ipblocks
: bytes
| ipb_address
: int64 2
: null
: bytes
| NO
: bytes
| tinyblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_user
: int64 3
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_by
: int64 4
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_by_text
: int64 5
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_by_actor
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_reason
: int64 7
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(767)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_reason_id
: int64 8
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_timestamp
: int64 9
: bytes
| '0x'
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_auto
: int64 10
: bytes
| '0'
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
| ipblocks
: bytes
| ipb_anon_only
: int64 11
: bytes
| '0'
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
| ipblocks
: bytes
| ipb_create_account
: int64 12
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
| ipblocks
: bytes
| ipb_enable_autoblock
: int64 13
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
| ipblocks
: bytes
| ipb_expiry
: int64 14
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_range_start
: int64 15
: null
: bytes
| NO
: bytes
| tinyblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_range_end
: int64 16
: null
: bytes
| NO
: bytes
| tinyblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_deleted
: int64 17
: bytes
| '0'
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
| ipblocks
: bytes
| ipb_block_email
: int64 18
: bytes
| '0'
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
| ipblocks
: bytes
| ipb_allow_usertalk
: int64 19
: bytes
| '0'
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
| ipblocks
: bytes
| ipb_parent_block_id
: int64 20
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
| iwlinks
: bytes
| iwl_from
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| iwlinks
: bytes
| iwl_prefix
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(20)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| iwlinks
: bytes
| iwl_title
: int64 3
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| job
: bytes
| job_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| job
: bytes
| job_cmd
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(60)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| job
: bytes
| job_namespace
: int64 3
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
|
: null
< row
: bytes
| job
: bytes
| job_title
: int64 4
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| job
: bytes
| job_timestamp
: int64 5
: null
: bytes
| YES
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| job
: bytes
| job_params
: int64 6
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| job
: bytes
| job_random
: int64 7
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| job
: bytes
| job_attempts
: int64 8
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| job
: bytes
| job_token
: int64 9
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| job
: bytes
| job_token_timestamp
: int64 10
: null
: bytes
| YES
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| job
: bytes
| job_sha1
: int64 11
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| l10n_cache
: bytes
| lc_lang
: int64 1
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| l10n_cache
: bytes
| lc_key
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| l10n_cache
: bytes
| lc_value
: int64 3
: null
: bytes
| NO
: bytes
| mediumblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| langlinks
: bytes
| ll_from
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| langlinks
: bytes
| ll_lang
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(20)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| langlinks
: bytes
| ll_title
: int64 3
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| log_search
: bytes
| ls_field
: int64 1
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| log_search
: bytes
| ls_value
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| log_search
: bytes
| ls_log_id
: int64 3
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| logging
: bytes
| log_type
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_action
: int64 3
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_timestamp
: int64 4
: bytes
| '0x3139373030313031303030303030'
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_user
: int64 5
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_user_text
: int64 6
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_actor
: int64 7
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_namespace
: int64 8
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
| logging
: bytes
| log_title
: int64 9
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_page
: int64 10
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_comment
: int64 11
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(767)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_comment_id
: int64 12
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_params
: int64 13
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_deleted
: int64 14
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| module_deps
: bytes
| md_module
: int64 1
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| module_deps
: bytes
| md_skin
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| module_deps
: bytes
| md_deps
: int64 3
: null
: bytes
| NO
: bytes
| mediumblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| objectcache
: bytes
| keyname
: int64 1
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| objectcache
: bytes
| value
: int64 2
: null
: bytes
| YES
: bytes
| mediumblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| objectcache
: bytes
| exptime
: int64 3
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_name
: int64 1
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_archive_name
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_size
: int64 3
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_width
: int64 4
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
| oldimage
: bytes
| oi_height
: int64 5
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
| oldimage
: bytes
| oi_bits
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
| oldimage
: bytes
| oi_description
: int64 7
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(767)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_description_id
: int64 8
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_user
: int64 9
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_user_text
: int64 10
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_actor
: int64 11
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_timestamp
: int64 12
: bytes
| '0x'
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_metadata
: int64 13
: null
: bytes
| NO
: bytes
| mediumblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_media_type
: int64 14
: null
: bytes
| YES
: bytes
| enum('UNKNOWN','BITMAP','DRAWING','AUDIO','VIDEO','MULTIMEDIA','OFFICE','TEXT','EXECUTABLE','ARCHIVE','3D')
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_major_mime
: int64 15
: bytes
| '0x756E6B6E6F776E'
: bytes
| NO
: bytes
| enum('unknown','application','audio','image','text','video','message','model','multipart','chemical')
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_minor_mime
: int64 16
: bytes
| '0x756E6B6E6F776E'
: bytes
| NO
: bytes
| varbinary(100)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_deleted
: int64 17
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_sha1
: int64 18
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| page
: bytes
| page_namespace
: int64 2
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
|
: null
< row
: bytes
| page
: bytes
| page_title
: int64 3
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_restrictions
: int64 4
: null
: bytes
| NO
: bytes
| tinyblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_is_redirect
: int64 5
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_is_new
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_random
: int64 7
: null
: bytes
| NO
: bytes
| double unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_touched
: int64 8
: bytes
| '0x'
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_links_updated
: int64 9
: null
: bytes
| YES
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_latest
: int64 10
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_len
: int64 11
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_content_model
: int64 12
: null
: bytes
| YES
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_lang
: int64 13
: null
: bytes
| YES
: bytes
| varbinary(35)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page_props
: bytes
| pp_page
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
|
: null
< row
: bytes
| page_props
: bytes
| pp_propname
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(60)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page_props
: bytes
| pp_value
: int64 3
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page_props
: bytes
| pp_sortkey
: int64 4
: null
: bytes
| YES
: bytes
| float
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page_restrictions
: bytes
| pr_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| page_restrictions
: bytes
| pr_page
: int64 2
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
|
: null
< row
: bytes
| page_restrictions
: bytes
| pr_type
: int64 3
: null
: bytes
| NO
: bytes
| varbinary(60)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page_restrictions
: bytes
| pr_level
: int64 4
: null
: bytes
| NO
: bytes
| varbinary(60)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page_restrictions
: bytes
| pr_cascade
: int64 5
: null
: bytes
| NO
: bytes
| tinyint
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page_restrictions
: bytes
| pr_user
: int64 6
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page_restrictions
: bytes
| pr_expiry
: int64 7
: null
: bytes
| YES
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| pagelinks
: bytes
| pl_from
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| pagelinks
: bytes
| pl_from_namespace
: int64 2
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
| pagelinks
: bytes
| pl_namespace
: int64 3
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
| pagelinks
: bytes
| pl_title
: int64 4
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| protected_titles
: bytes
| pt_namespace
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
|
: null
< row
: bytes
| protected_titles
: bytes
| pt_title
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| protected_titles
: bytes
| pt_user
: int64 3
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| protected_titles
: bytes
| pt_reason
: int64 4
: bytes
| ''
: bytes
| YES
: bytes
| varbinary(767)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| protected_titles
: bytes
| pt_reason_id
: int64 5
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| protected_titles
: bytes
| pt_timestamp
: int64 6
: null
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| protected_titles
: bytes
| pt_expiry
: int64 7
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| protected_titles
: bytes
| pt_create_perm
: int64 8
: null
: bytes
| NO
: bytes
| varbinary(60)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| querycache
: bytes
| qc_type
: int64 1
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| querycache
: bytes
| qc_value
: int64 2
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| querycache
: bytes
| qc_namespace
: int64 3
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
| querycache
: bytes
| qc_title
: int64 4
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| querycache_info
: bytes
| qci_type
: int64 1
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| querycache_info
: bytes
| qci_timestamp
: int64 2
: bytes
| '0x3139373030313031303030303030'
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| querycachetwo
: bytes
| qcc_type
: int64 1
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| querycachetwo
: bytes
| qcc_value
: int64 2
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| querycachetwo
: bytes
| qcc_namespace
: int64 3
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
| querycachetwo
: bytes
| qcc_title
: int64 4
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| querycachetwo
: bytes
| qcc_namespacetwo
: int64 5
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
| querycachetwo
: bytes
| qcc_titletwo
: int64 6
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_id
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
| recentchanges
: bytes
| rc_timestamp
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_user
: int64 3
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_user_text
: int64 4
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_actor
: int64 5
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_namespace
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
| recentchanges
: bytes
| rc_title
: int64 7
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_comment
: int64 8
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(767)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_comment_id
: int64 9
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_minor
: int64 10
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_bot
: int64 11
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_new
: int64 12
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_cur_id
: int64 13
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_this_oldid
: int64 14
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_last_oldid
: int64 15
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_type
: int64 16
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_source
: int64 17
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(16)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_patrolled
: int64 18
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_ip
: int64 19
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(40)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_old_len
: int64 20
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
| recentchanges
: bytes
| rc_new_len
: int64 21
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
| recentchanges
: bytes
| rc_deleted
: int64 22
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_logid
: int64 23
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_log_type
: int64 24
: null
: bytes
| YES
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_log_action
: int64 25
: null
: bytes
| YES
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_params
: int64 26
: null
: bytes
| YES
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| redirect
: bytes
| rd_from
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| redirect
: bytes
| rd_namespace
: int64 2
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
| redirect
: bytes
| rd_title
: int64 3
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| redirect
: bytes
| rd_interwiki
: int64 4
: null
: bytes
| YES
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| redirect
: bytes
| rd_fragment
: int64 5
: null
: bytes
| YES
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| revision
: bytes
| rev_page
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_text_id
: int64 3
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_comment
: int64 4
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(767)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_user
: int64 5
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_user_text
: int64 6
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_timestamp
: int64 7
: bytes
| '0x'
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_minor_edit
: int64 8
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_deleted
: int64 9
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_len
: int64 10
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_parent_id
: int64 11
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_sha1
: int64 12
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_content_model
: int64 13
: null
: bytes
| YES
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_content_format
: int64 14
: null
: bytes
| YES
: bytes
| varbinary(64)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision_actor_temp
: bytes
| revactor_rev
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision_actor_temp
: bytes
| revactor_actor
: int64 2
: null
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision_actor_temp
: bytes
| revactor_timestamp
: int64 3
: bytes
| '0x'
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision_actor_temp
: bytes
| revactor_page
: int64 4
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision_comment_temp
: bytes
| revcomment_rev
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision_comment_temp
: bytes
| revcomment_comment_id
: int64 2
: null
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| searchindex
: bytes
| si_page
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| searchindex
: bytes
| si_title
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| searchindex
: bytes
| si_text
: int64 3
: null
: bytes
| NO
: bytes
| mediumtext
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| site_identifiers
: bytes
| si_site
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| site_identifiers
: bytes
| si_type
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| site_identifiers
: bytes
| si_key
: int64 3
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| site_stats
: bytes
| ss_row_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| site_stats
: bytes
| ss_total_edits
: int64 2
: null
: bytes
| YES
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| site_stats
: bytes
| ss_good_articles
: int64 3
: null
: bytes
| YES
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| site_stats
: bytes
| ss_total_pages
: int64 4
: null
: bytes
| YES
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| site_stats
: bytes
| ss_users
: int64 5
: null
: bytes
| YES
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| site_stats
: bytes
| ss_active_users
: int64 6
: null
: bytes
| YES
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| site_stats
: bytes
| ss_images
: int64 7
: null
: bytes
| YES
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sites
: bytes
| site_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| sites
: bytes
| site_global_key
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sites
: bytes
| site_type
: int64 3
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sites
: bytes
| site_group
: int64 4
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sites
: bytes
| site_source
: int64 5
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sites
: bytes
| site_language
: int64 6
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sites
: bytes
| site_protocol
: int64 7
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sites
: bytes
| site_domain
: int64 8
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sites
: bytes
| site_data
: int64 9
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sites
: bytes
| site_forward
: int64 10
: null
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
| sites
: bytes
| site_config
: int64 11
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| slot_roles
: bytes
| role_id
: int64 1
: null
: bytes
| NO
: bytes
| smallint
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
| slot_roles
: bytes
| role_name
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(64)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| slots
: bytes
| slot_revision_id
: int64 1
: null
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| slots
: bytes
| slot_role_id
: int64 2
: null
: bytes
| NO
: bytes
| smallint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| slots
: bytes
| slot_content_id
: int64 3
: null
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| slots
: bytes
| slot_origin
: int64 4
: null
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| tag_summary
: bytes
| ts_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| tag_summary
: bytes
| ts_rc_id
: int64 2
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
| tag_summary
: bytes
| ts_log_id
: int64 3
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| tag_summary
: bytes
| ts_rev_id
: int64 4
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| tag_summary
: bytes
| ts_tags
: int64 5
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| templatelinks
: bytes
| tl_from
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| templatelinks
: bytes
| tl_from_namespace
: int64 2
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
| templatelinks
: bytes
| tl_namespace
: int64 3
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
| templatelinks
: bytes
| tl_title
: int64 4
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| text
: bytes
| old_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| text
: bytes
| old_text
: int64 2
: null
: bytes
| NO
: bytes
| mediumblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| text
: bytes
| old_flags
: int64 3
: null
: bytes
| NO
: bytes
| tinyblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| transcache
: bytes
| tc_url
: int64 1
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| transcache
: bytes
| tc_contents
: int64 2
: null
: bytes
| YES
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| transcache
: bytes
| tc_time
: int64 3
: null
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| updatelog
: bytes
| ul_key
: int64 1
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| updatelog
: bytes
| ul_value
: int64 2
: null
: bytes
| YES
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| uploadstash
: bytes
| us_user
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_key
: int64 3
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_orig_path
: int64 4
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_path
: int64 5
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_source_type
: int64 6
: null
: bytes
| YES
: bytes
| varbinary(50)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_timestamp
: int64 7
: null
: bytes
| NO
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_status
: int64 8
: null
: bytes
| NO
: bytes
| varbinary(50)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_chunk_inx
: int64 9
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_props
: int64 10
: null
: bytes
| YES
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_size
: int64 11
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_sha1
: int64 12
: null
: bytes
| NO
: bytes
| varbinary(31)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_mime
: int64 13
: null
: bytes
| YES
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_media_type
: int64 14
: null
: bytes
| YES
: bytes
| enum('UNKNOWN','BITMAP','DRAWING','AUDIO','VIDEO','MULTIMEDIA','OFFICE','TEXT','EXECUTABLE','ARCHIVE','3D')
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_image_width
: int64 15
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_image_height
: int64 16
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_image_bits
: int64 17
: null
: bytes
| YES
: bytes
| smallint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| user
: bytes
| user_name
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_real_name
: int64 3
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_password
: int64 4
: null
: bytes
| NO
: bytes
| tinyblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_newpassword
: int64 5
: null
: bytes
| NO
: bytes
| tinyblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_newpass_time
: int64 6
: null
: bytes
| YES
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_email
: int64 7
: null
: bytes
| NO
: bytes
| tinyblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_touched
: int64 8
: bytes
| '0x'
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_token
: int64 9
: bytes
| '0x'
: bytes
| NO
: bytes
| binary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_email_authenticated
: int64 10
: null
: bytes
| YES
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_email_token
: int64 11
: null
: bytes
| YES
: bytes
| binary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_email_token_expires
: int64 12
: null
: bytes
| YES
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_registration
: int64 13
: null
: bytes
| YES
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_editcount
: int64 14
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
| user
: bytes
| user_password_expires
: int64 15
: null
: bytes
| YES
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_former_groups
: bytes
| ufg_user
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_former_groups
: bytes
| ufg_group
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_groups
: bytes
| ug_user
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_groups
: bytes
| ug_group
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_groups
: bytes
| ug_expiry
: int64 3
: null
: bytes
| YES
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_newtalk
: bytes
| user_id
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_newtalk
: bytes
| user_ip
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(40)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_newtalk
: bytes
| user_last_timestamp
: int64 3
: null
: bytes
| YES
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_properties
: bytes
| up_user
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_properties
: bytes
| up_property
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_properties
: bytes
| up_value
: int64 3
: null
: bytes
| YES
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| valid_tag
: bytes
| vt_tag
: int64 1
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| watchlist
: bytes
| wl_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| watchlist
: bytes
| wl_user
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| watchlist
: bytes
| wl_namespace
: int64 3
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
| watchlist
: bytes
| wl_title
: int64 4
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| watchlist
: bytes
| wl_notificationtimestamp
: int64 5
: null
: bytes
| YES
: bytes
| varbinary(14)
: bytes
|
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
| sdl_rw_mediawiki_4
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
| sdl_rw_mediawiki_4
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
| sdl_rw_mediawiki_4
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
| sdl_rw_mediawiki_4
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
| sdl_rw_mediawiki_4
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
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
| sdl_rw_mediawiki_4
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
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
| sdl_rw_mediawiki_4
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
| sdl_rw_mediawiki_4
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
| sdl_rw_mediawiki_4
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
| actor
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| archive
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| bot_passwords
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| category
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| categorylinks
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| change_tag
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| comment
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| content
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| content_models
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| externallinks
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| filearchive
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| image
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| image_comment_temp
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| imagelinks
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| interwiki
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| ip_changes
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| ipblocks
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| iwlinks
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| job
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| l10n_cache
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| langlinks
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| log_search
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| logging
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| module_deps
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| objectcache
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| oldimage
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| page
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| page_props
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| page_restrictions
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| pagelinks
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| protected_titles
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| querycache
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| querycache_info
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| querycachetwo
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| recentchanges
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| redirect
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| revision
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
| max_rows=10000000 avg_row_length=1024
: bytes
| ''
: bytes
| binary
< row
: bytes
| revision_actor_temp
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| revision_comment_temp
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| searchindex
: bytes
| BASE TABLE
: bytes
| MyISAM
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 0
: bytes
| 1024
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| site_identifiers
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| site_stats
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| sites
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| slot_roles
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| slots
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| tag_summary
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| templatelinks
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| text
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
| max_rows=10000000 avg_row_length=10240
: bytes
| ''
: bytes
| binary
< row
: bytes
| transcache
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| updatelog
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| uploadstash
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| user
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| user_former_groups
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| user_groups
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| user_newtalk
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| user_properties
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| valid_tag
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| watchlist
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
> query
|
| 		SELECT t.NAME
| 		FROM information_schema.INNODB_TABLES t
| 		LEFT JOIN information_schema.INNODB_TABLESPACES ts ON t.SPACE = ts.SPACE
| 		WHERE t.NAME LIKE CONCAT(?, '/%') AND (ts.SPACE_TYPE IS NULL OR ts.SPACE_TYPE <> 'Single')
: string
| sdl_rw_mediawiki_4
< result
< column "NAME" "VARCHAR" "string"
> query
| SHOW CREATE TABLE `sdl_rw_mediawiki_4`.`searchindex`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| searchindex
: bytes
| CREATE TABLE `searchindex` (
|   `si_page` int unsigned NOT NULL,
|   `si_title` varchar(255) NOT NULL DEFAULT '',
|   `si_text` mediumtext NOT NULL,
|   UNIQUE KEY `si_page` (`si_page`),
|   FULLTEXT KEY `si_title` (`si_title`),
|   FULLTEXT KEY `si_text` (`si_text`)
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb3
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| sdl_rw_mediawiki_4
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci

== TestSDLRealWorldIdempotence default#2
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
| sakila
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
| actor
: bytes
| idx_actor_last_name
: bytes
| last_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| actor
: bytes
| PRIMARY
: bytes
| actor_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| address
: bytes
| idx_fk_city_id
: bytes
| city_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| address
: bytes
| PRIMARY
: bytes
| address_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| category
: bytes
| PRIMARY
: bytes
| category_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| city
: bytes
| idx_fk_country_id
: bytes
| country_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| city
: bytes
| PRIMARY
: bytes
| city_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| country
: bytes
| PRIMARY
: bytes
| country_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| customer
: bytes
| idx_fk_address_id
: bytes
| address_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| customer
: bytes
| idx_fk_store_id
: bytes
| store_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| customer
: bytes
| idx_last_name
: bytes
| last_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| customer
: bytes
| PRIMARY
: bytes
| customer_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| film
: bytes
| idx_fk_language_id
: bytes
| language_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| film
: bytes
| idx_fk_original_language_id
: bytes
| original_language_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| film
: bytes
| idx_title
: bytes
| title
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| film
: bytes
| PRIMARY
: bytes
| film_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| film_actor
: bytes
| idx_fk_film_id
: bytes
| film_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| film_actor
: bytes
| PRIMARY
: bytes
| actor_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| film_actor
: bytes
| PRIMARY
: bytes
| film_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| film_category
: bytes
| fk_film_category_category
: bytes
| category_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| film_category
: bytes
| PRIMARY
: bytes
| film_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| film_category
: bytes
| PRIMARY
: bytes
| category_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| film_text
: bytes
| idx_title_description
: bytes
| title
: null
: int64 -1
: null
: int64 1
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
< row
: bytes
| film_text
: bytes
| idx_title_description
: bytes
| description
: null
: int64 -1
: null
: int64 2
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
< row
: bytes
| film_text
: bytes
| PRIMARY
: bytes
| film_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| inventory
: bytes
| idx_fk_film_id
: bytes
| film_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| inventory
: bytes
| idx_store_id_film_id
: bytes
| store_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| inventory
: bytes
| idx_store_id_film_id
: bytes
| film_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| inventory
: bytes
| PRIMARY
: bytes
| inventory_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| language
: bytes
| PRIMARY
: bytes
| language_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| payment
: bytes
| fk_payment_rental
: bytes
| rental_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| payment
: bytes
| idx_fk_customer_id
: bytes
| customer_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| payment
: bytes
| idx_fk_staff_id
: bytes
| staff_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| payment
: bytes
| PRIMARY
: bytes
| payment_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| rental
: bytes
| idx_fk_customer_id
: bytes
| customer_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| rental
: bytes
| idx_fk_inventory_id
: bytes
| inventory_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| rental
: bytes
| idx_fk_staff_id
: bytes
| staff_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| rental
: bytes
| PRIMARY
: bytes
| rental_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| rental
: bytes
| rental_date
: bytes
| rental_date
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| rental
: bytes
| rental_date
: bytes
| inventory_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| rental
: bytes
| rental_date
: bytes
| customer_id
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| staff
: bytes
| idx_fk_address_id
: bytes
| address_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| staff
: bytes
| idx_fk_store_id
: bytes
| store_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| staff
: bytes
| PRIMARY
: bytes
| staff_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| store
: bytes
| idx_fk_address_id
: bytes
| address_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| store
: bytes
| idx_unique_manager
: bytes
| manager_staff_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| store
: bytes
| PRIMARY
: bytes
| store_id
: bytes
| A
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
| sakila
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
| actor
: bytes
| actor_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| actor
: bytes
| first_name
: int64 2
: null
: bytes
| NO
: bytes
| varchar(45)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| actor
: bytes
| last_name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(45)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| actor
: bytes
| last_update
: int64 4
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| actor_info
: bytes
| actor_id
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| actor_info
: bytes
| first_name
: int64 2
: null
: bytes
| NO
: bytes
| varchar(45)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| actor_info
: bytes
| last_name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(45)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| actor_info
: bytes
| film_info
: int64 4
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| address
: bytes
| address_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| address
: bytes
| address
: int64 2
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| address
: bytes
| address2
: int64 3
: null
: bytes
| YES
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| address
: bytes
| district
: int64 4
: null
: bytes
| NO
: bytes
| varchar(20)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| address
: bytes
| city_id
: int64 5
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| address
: bytes
| postal_code
: int64 6
: null
: bytes
| YES
: bytes
| varchar(10)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| address
: bytes
| phone
: int64 7
: null
: bytes
| NO
: bytes
| varchar(20)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| address
: bytes
| last_update
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
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| category
: bytes
| category_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| category
: bytes
| name
: int64 2
: null
: bytes
| NO
: bytes
| varchar(25)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| category
: bytes
| last_update
: int64 3
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| city
: bytes
| city_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| city
: bytes
| city
: int64 2
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| city
: bytes
| country_id
: int64 3
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| city
: bytes
| last_update
: int64 4
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| country
: bytes
| country_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| country
: bytes
| country
: int64 2
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| country
: bytes
| last_update
: int64 3
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| customer
: bytes
| customer_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| customer
: bytes
| store_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer
: bytes
| first_name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(45)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer
: bytes
| last_name
: int64 4
: null
: bytes
| NO
: bytes
| varchar(45)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer
: bytes
| email
: int64 5
: null
: bytes
| YES
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer
: bytes
| address_id
: int64 6
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer
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
| customer
: bytes
| create_date
: int64 8
: null
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer
: bytes
| last_update
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| YES
: bytes
| timestamp
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
| customer_list
: bytes
| ID
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer_list
: bytes
| name
: int64 2
: null
: bytes
| YES
: bytes
| varchar(91)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer_list
: bytes
| address
: int64 3
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer_list
: bytes
| zip code
: int64 4
: null
: bytes
| YES
: bytes
| varchar(10)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer_list
: bytes
| phone
: int64 5
: null
: bytes
| NO
: bytes
| varchar(20)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer_list
: bytes
| city
: int64 6
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer_list
: bytes
| country
: int64 7
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer_list
: bytes
| notes
: int64 8
: bytes
| ''
: bytes
| NO
: bytes
| varchar(6)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| customer_list
: bytes
| SID
: int64 9
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| film_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| film
: bytes
| title
: int64 2
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| description
: int64 3
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| release_year
: int64 4
: null
: bytes
| YES
: bytes
| year
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| language_id
: int64 5
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| original_language_id
: int64 6
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| rental_duration
: int64 7
: bytes
| '3'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| rental_rate
: int64 8
: bytes
| '4.99'
: bytes
| NO
: bytes
| decimal(4,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| length
: int64 9
: null
: bytes
| YES
: bytes
| smallint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| replacement_cost
: int64 10
: bytes
| '19.99'
: bytes
| NO
: bytes
| decimal(5,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| rating
: int64 11
: bytes
| 'G'
: bytes
| YES
: bytes
| enum('G','PG','PG-13','R','NC-17')
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| special_features
: int64 12
: null
: bytes
| YES
: bytes
| set('Trailers','Commentaries','Deleted Scenes','Behind the Scenes')
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film
: bytes
| last_update
: int64 13
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| film_actor
: bytes
| actor_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_actor
: bytes
| film_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_actor
: bytes
| last_update
: int64 3
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| film_category
: bytes
| film_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_category
: bytes
| category_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_category
: bytes
| last_update
: int64 3
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| film_list
: bytes
| FID
: int64 1
: bytes
| '0'
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_list
: bytes
| title
: int64 2
: null
: bytes
| YES
: bytes
| varchar(255)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_list
: bytes
| description
: int64 3
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_list
: bytes
| category
: int64 4
: null
: bytes
| NO
: bytes
| varchar(25)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_list
: bytes
| price
: int64 5
: bytes
| '4.99'
: bytes
| YES
: bytes
| decimal(4,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_list
: bytes
| length
: int64 6
: null
: bytes
| YES
: bytes
| smallint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_list
: bytes
| rating
: int64 7
: bytes
| 'G'
: bytes
| YES
: bytes
| enum('G','PG','PG-13','R','NC-17')
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_list
: bytes
| actors
: int64 8
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_text
: bytes
| film_id
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
|
: null
< row
: bytes
| film_text
: bytes
| title
: int64 2
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| film_text
: bytes
| description
: int64 3
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| inventory
: bytes
| inventory_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| inventory
: bytes
| film_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| inventory
: bytes
| store_id
: int64 3
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| inventory
: bytes
| last_update
: int64 4
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| language
: bytes
| language_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| language
: bytes
| name
: int64 2
: null
: bytes
| NO
: bytes
| char(20)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| language
: bytes
| last_update
: int64 3
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| nicer_but_slower_film_list
: bytes
| FID
: int64 1
: bytes
| '0'
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| nicer_but_slower_film_list
: bytes
| title
: int64 2
: null
: bytes
| YES
: bytes
| varchar(255)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| nicer_but_slower_film_list
: bytes
| description
: int64 3
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| nicer_but_slower_film_list
: bytes
| category
: int64 4
: null
: bytes
| NO
: bytes
| varchar(25)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| nicer_but_slower_film_list
: bytes
| price
: int64 5
: bytes
| '4.99'
: bytes
| YES
: bytes
| decimal(4,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| nicer_but_slower_film_list
: bytes
| length
: int64 6
: null
: bytes
| YES
: bytes
| smallint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| nicer_but_slower_film_list
: bytes
| rating
: int64 7
: bytes
| 'G'
: bytes
| YES
: bytes
| enum('G','PG','PG-13','R','NC-17')
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| nicer_but_slower_film_list
: bytes
| actors
: int64 8
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| payment
: bytes
| payment_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| payment
: bytes
| customer_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| payment
: bytes
| staff_id
: int64 3
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| payment
: bytes
| rental_id
: int64 4
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
| payment
: bytes
| amount
: int64 5
: null
: bytes
| NO
: bytes
| decimal(5,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| payment
: bytes
| payment_date
: int64 6
: null
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| payment
: bytes
| last_update
: int64 7
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| YES
: bytes
| timestamp
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
| rental
: bytes
| rental_id
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
| rental
: bytes
| rental_date
: int64 2
: null
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| rental
: bytes
| inventory_id
: int64 3
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| rental
: bytes
| customer_id
: int64 4
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| rental
: bytes
| return_date
: int64 5
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| rental
: bytes
| staff_id
: int64 6
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| rental
: bytes
| last_update
: int64 7
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| sales_by_film_category
: bytes
| category
: int64 1
: null
: bytes
| NO
: bytes
| varchar(25)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sales_by_film_category
: bytes
| total_sales
: int64 2
: null
: bytes
| YES
: bytes
| decimal(27,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sales_by_store
: bytes
| store
: int64 1
: null
: bytes
| YES
: bytes
| varchar(101)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sales_by_store
: bytes
| manager
: int64 2
: null
: bytes
| YES
: bytes
| varchar(91)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sales_by_store
: bytes
| total_sales
: int64 3
: null
: bytes
| YES
: bytes
| decimal(27,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff
: bytes
| staff_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| staff
: bytes
| first_name
: int64 2
: null
: bytes
| NO
: bytes
| varchar(45)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff
: bytes
| last_name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(45)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff
: bytes
| address_id
: int64 4
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff
: bytes
| picture
: int64 5
: null
: bytes
| YES
: bytes
| mediumblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff
: bytes
| email
: int64 6
: null
: bytes
| YES
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff
: bytes
| store_id
: int64 7
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff
: bytes
| active
: int64 8
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
| staff
: bytes
| username
: int64 9
: null
: bytes
| NO
: bytes
| varchar(16)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff
: bytes
| password
: int64 10
: null
: bytes
| YES
: bytes
| varchar(40)
: bytes
| utf8mb3
: bytes
| utf8mb3_bin
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff
: bytes
| last_update
: int64 11
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| staff_list
: bytes
| ID
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff_list
: bytes
| name
: int64 2
: null
: bytes
| YES
: bytes
| varchar(91)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff_list
: bytes
| address
: int64 3
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff_list
: bytes
| zip code
: int64 4
: null
: bytes
| YES
: bytes
| varchar(10)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff_list
: bytes
| phone
: int64 5
: null
: bytes
| NO
: bytes
| varchar(20)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff_list
: bytes
| city
: int64 6
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff_list
: bytes
| country
: int64 7
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| staff_list
: bytes
| SID
: int64 8
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| store
: bytes
| store_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| store
: bytes
| manager_staff_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| store
: bytes
| address_id
: int64 3
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| store
: bytes
| last_update
: int64 4
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
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
| sakila
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
| sakila
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| actor_info
: bytes
| select `a`.`actor_id` AS `actor_id`,`a`.`first_name` AS `first_name`,`a`.`last_name` AS `last_name`,group_concat(distinct concat(`c`.`name`,': ',(select group_concat(`f`.`title` order by `f`.`title` ASC separator ', ') from ((`sakila`.`film` `f` join `sakila`.`film_category` `fc` on((`f`.`film_id` = `fc`.`film_id`))) join `sakila`.`film_actor` `fa` on((`f`.`film_id` = `fa`.`film_id`))) where ((`fc`.`category_id` = `c`.`category_id`) and (`fa`.`actor_id` = `a`.`actor_id`)))) order by `c`.`name` ASC separator '; ') AS `film_info` from (((`sakila`.`actor` `a` left join `sakila`.`film_actor` `fa` on((`a`.`actor_id` = `fa`.`actor_id`))) left join `sakila`.`film_category` `fc` on((`fa`.`film_id` = `fc`.`film_id`))) left join `sakila`.`category` `c` on((`fc`.`category_id` = `c`.`category_id`))) group by `a`.`actor_id`,`a`.`first_name`,`a`.`last_name`
< row
: bytes
| customer_list
: bytes
| select `cu`.`customer_id` AS `ID`,concat(`cu`.`first_name`,' ',`cu`.`last_name`) AS `name`,`a`.`address` AS `address`,`a`.`postal_code` AS `zip code`,`a`.`phone` AS `phone`,`sakila`.`city`.`city` AS `city`,`sakila`.`country`.`country` AS `country`,if(`cu`.`active`,'active','') AS `notes`,`cu`.`store_id` AS `SID` from (((`sakila`.`customer` `cu` join `sakila`.`address` `a` on((`cu`.`address_id` = `a`.`address_id`))) join `sakila`.`city` on((`a`.`city_id` = `sakila`.`city`.`city_id`))) join `sakila`.`country` on((`sakila`.`city`.`country_id` = `sakila`.`country`.`country_id`)))
< row
: bytes
| film_list
: bytes
| select `sakila`.`film`.`film_id` AS `FID`,`sakila`.`film`.`title` AS `title`,`sakila`.`film`.`description` AS `description`,`sakila`.`category`.`name` AS `category`,`sakila`.`film`.`rental_rate` AS `price`,`sakila`.`film`.`length` AS `length`,`sakila`.`film`.`rating` AS `rating`,group_concat(concat(`sakila`.`actor`.`first_name`,' ',`sakila`.`actor`.`last_name`) separator ', ') AS `actors` from ((((`sakila`.`category` left join `sakila`.`film_category` on((`sakila`.`category`.`category_id` = `sakila`.`film_category`.`category_id`))) left join `sakila`.`film` on((`sakila`.`film_category`.`film_id` = `sakila`.`film`.`film_id`))) join `sakila`.`film_actor` on((`sakila`.`film`.`film_id` = `sakila`.`film_actor`.`film_id`))) join `sakila`.`actor` on((`sakila`.`film_actor`.`actor_id` = `sakila`.`actor`.`actor_id`))) group by `sakila`.`film`.`film_id`,`sakila`.`film`.`title`,`sakila`.`film`.`description`,`sakila`.`film`.`rental_rate`,`sakila`.`film`.`length`,`sakila`.`film`.`rating`,`sakila`.`category`.`name`
< row
: bytes
| nicer_but_slower_film_list
: bytes
| select `sakila`.`film`.`film_id` AS `FID`,`sakila`.`film`.`title` AS `title`,`sakila`.`film`.`description` AS `description`,`sakila`.`category`.`name` AS `category`,`sakila`.`film`.`rental_rate` AS `price`,`sakila`.`film`.`length` AS `length`,`sakila`.`film`.`rating` AS `rating`,group_concat(concat(concat(upper(substr(`sakila`.`actor`.`first_name`,1,1)),lower(substr(`sakila`.`actor`.`first_name`,2,length(`sakila`.`actor`.`first_name`))),' ',concat(upper(substr(`sakila`.`actor`.`last_name`,1,1)),lower(substr(`sakila`.`actor`.`last_name`,2,length(`sakila`.`actor`.`last_name`)))))) separator ', ') AS `actors` from ((((`sakila`.`category` left join `sakila`.`film_category` on((`sakila`.`category`.`category_id` = `sakila`.`film_category`.`category_id`))) left join `sakila`.`film` on((`sakila`.`film_category`.`film_id` = `sakila`.`film`.`film_id`))) join `sakila`.`film_actor` on((`sakila`.`film`.`film_id` = `sakila`.`film_actor`.`film_id`))) join `sakila`.`actor` on((`sakila`.`film_actor`.`actor_id` = `sakila`.`actor`.`actor_id`))) group by `sakila`.`film`.`film_id`,`sakila`.`film`.`title`,`sakila`.`film`.`description`,`sakila`.`film`.`rental_rate`,`sakila`.`film`.`length`,`sakila`.`film`.`rating`,`sakila`.`category`.`name`
< row
: bytes
| sales_by_film_category
: bytes
| select `c`.`name` AS `category`,sum(`p`.`amount`) AS `total_sales` from (((((`sakila`.`payment` `p` join `sakila`.`rental` `r` on((`p`.`rental_id` = `r`.`rental_id`))) join `sakila`.`inventory` `i` on((`r`.`inventory_id` = `i`.`inventory_id`))) join `sakila`.`film` `f` on((`i`.`film_id` = `f`.`film_id`))) join `sakila`.`film_category` `fc` on((`f`.`film_id` = `fc`.`film_id`))) join `sakila`.`category` `c` on((`fc`.`category_id` = `c`.`category_id`))) group by `c`.`name` order by `total_sales` desc
< row
: bytes
| sales_by_store
: bytes
| select concat(`c`.`city`,',',`cy`.`country`) AS `store`,concat(`m`.`first_name`,' ',`m`.`last_name`) AS `manager`,sum(`p`.`amount`) AS `total_sales` from (((((((`sakila`.`payment` `p` join `sakila`.`rental` `r` on((`p`.`rental_id` = `r`.`rental_id`))) join `sakila`.`inventory` `i` on((`r`.`inventory_id` = `i`.`inventory_id`))) join `sakila`.`store` `s` on((`i`.`store_id` = `s`.`store_id`))) join `sakila`.`address` `a` on((`s`.`address_id` = `a`.`address_id`))) join `sakila`.`city` `c` on((`a`.`city_id` = `c`.`city_id`))) join `sakila`.`country` `cy` on((`c`.`country_id` = `cy`.`country_id`))) join `sakila`.`staff` `m` on((`s`.`manager_staff_id` = `m`.`staff_id`))) group by `s`.`store_id` order by `cy`.`country`,`c`.`city`
< row
: bytes
| staff_list
: bytes
| select `s`.`staff_id` AS `ID`,concat(`s`.`first_name`,' ',`s`.`last_name`) AS `name`,`a`.`address` AS `address`,`a`.`postal_code` AS `zip code`,`a`.`phone` AS `phone`,`sakila`.`city`.`city` AS `city`,`sakila`.`country`.`country` AS `country`,`s`.`store_id` AS `SID` from (((`sakila`.`staff` `s` join `sakila`.`address` `a` on((`s`.`address_id` = `a`.`address_id`))) join `sakila`.`city` on((`a`.`city_id` = `sakila`.`city`.`city_id`))) join `sakila`.`country` on((`sakila`.`city`.`country_id` = `sakila`.`country`.`country_id`)))
> query
| SHOW CREATE VIEW `sakila`.`actor_info`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| actor_info
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY INVOKER VIEW `actor_info` AS select `a`.`actor_id` AS `actor_id`,`a`.`first_name` AS `first_name`,`a`.`last_name` AS `last_name`,group_concat(distinct concat(`c`.`name`,': ',(select group_concat(`f`.`title` order by `f`.`title` ASC separator ', ') from ((`film` `f` join `film_category` `fc` on((`f`.`film_id` = `fc`.`film_id`))) join `film_actor` `fa` on((`f`.`film_id` = `fa`.`film_id`))) where ((`fc`.`category_id` = `c`.`category_id`) and (`fa`.`actor_id` = `a`.`actor_id`)))) order by `c`.`name` ASC separator '; ') AS `film_info` from (((`actor` `a` left join `film_actor` `fa` on((`a`.`actor_id` = `fa`.`actor_id`))) left join `film_category` `fc` on((`fa`.`film_id` = `fc`.`film_id`))) left join `category` `c` on((`fc`.`category_id` = `c`.`category_id`))) group by `a`.`actor_id`,`a`.`first_name`,`a`.`last_name`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `sakila`.`customer_list`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| customer_list
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `customer_list` AS select `cu`.`customer_id` AS `ID`,concat(`cu`.`first_name`,_utf8mb3' ',`cu`.`last_name`) AS `name`,`a`.`address` AS `address`,`a`.`postal_code` AS `zip code`,`a`.`phone` AS `phone`,`city`.`city` AS `city`,`country`.`country` AS `country`,if(`cu`.`active`,_utf8mb3'active',_utf8mb3'') AS `notes`,`cu`.`store_id` AS `SID` from (((`customer` `cu` join `address` `a` on((`cu`.`address_id` = `a`.`address_id`))) join `city` on((`a`.`city_id` = `city`.`city_id`))) join `country` on((`city`.`country_id` = `country`.`country_id`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `sakila`.`film_list`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| film_list
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `film_list` AS select `film`.`film_id` AS `FID`,`film`.`title` AS `title`,`film`.`description` AS `description`,`category`.`name` AS `category`,`film`.`rental_rate` AS `price`,`film`.`length` AS `length`,`film`.`rating` AS `rating`,group_concat(concat(`actor`.`first_name`,_utf8mb3' ',`actor`.`last_name`) separator ', ') AS `actors` from ((((`category` left join `film_category` on((`category`.`category_id` = `film_category`.`category_id`))) left join `film` on((`film_category`.`film_id` = `film`.`film_id`))) join `film_actor` on((`film`.`film_id` = `film_actor`.`film_id`))) join `actor` on((`film_actor`.`actor_id` = `actor`.`actor_id`))) group by `film`.`film_id`,`film`.`title`,`film`.`description`,`film`.`rental_rate`,`film`.`length`,`film`.`rating`,`category`.`name`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `sakila`.`nicer_but_slower_film_list`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| nicer_but_slower_film_list
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `nicer_but_slower_film_list` AS select `film`.`film_id` AS `FID`,`film`.`title` AS `title`,`film`.`description` AS `description`,`category`.`name` AS `category`,`film`.`rental_rate` AS `price`,`film`.`length` AS `length`,`film`.`rating` AS `rating`,group_concat(concat(concat(upper(substr(`actor`.`first_name`,1,1)),lower(substr(`actor`.`first_name`,2,length(`actor`.`first_name`))),_utf8mb3' ',concat(upper(substr(`actor`.`last_name`,1,1)),lower(substr(`actor`.`last_name`,2,length(`actor`.`last_name`)))))) separator ', ') AS `actors` from ((((`category` left join `film_category` on((`category`.`category_id` = `film_category`.`category_id`))) left join `film` on((`film_category`.`film_id` = `film`.`film_id`))) join `film_actor` on((`film`.`film_id` = `film_actor`.`film_id`))) join `actor` on((`film_actor`.`actor_id` = `actor`.`actor_id`))) group by `film`.`film_id`,`film`.`title`,`film`.`description`,`film`.`rental_rate`,`film`.`length`,`film`.`rating`,`category`.`name`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `sakila`.`sales_by_film_category`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| sales_by_film_category
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `sales_by_film_category` AS select `c`.`name` AS `category`,sum(`p`.`amount`) AS `total_sales` from (((((`payment` `p` join `rental` `r` on((`p`.`rental_id` = `r`.`rental_id`))) join `inventory` `i` on((`r`.`inventory_id` = `i`.`inventory_id`))) join `film` `f` on((`i`.`film_id` = `f`.`film_id`))) join `film_category` `fc` on((`f`.`film_id` = `fc`.`film_id`))) join `category` `c` on((`fc`.`category_id` = `c`.`category_id`))) group by `c`.`name` order by `total_sales` desc
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `sakila`.`sales_by_store`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| sales_by_store
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `sales_by_store` AS select concat(`c`.`city`,_utf8mb3',',`cy`.`country`) AS `store`,concat(`m`.`first_name`,_utf8mb3' ',`m`.`last_name`) AS `manager`,sum(`p`.`amount`) AS `total_sales` from (((((((`payment` `p` join `rental` `r` on((`p`.`rental_id` = `r`.`rental_id`))) join `inventory` `i` on((`r`.`inventory_id` = `i`.`inventory_id`))) join `store` `s` on((`i`.`store_id` = `s`.`store_id`))) join `address` `a` on((`s`.`address_id` = `a`.`address_id`))) join `city` `c` on((`a`.`city_id` = `c`.`city_id`))) join `country` `cy` on((`c`.`country_id` = `cy`.`country_id`))) join `staff` `m` on((`s`.`manager_staff_id` = `m`.`staff_id`))) group by `s`.`store_id` order by `cy`.`country`,`c`.`city`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `sakila`.`staff_list`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| staff_list
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `staff_list` AS select `s`.`staff_id` AS `ID`,concat(`s`.`first_name`,_utf8mb3' ',`s`.`last_name`) AS `name`,`a`.`address` AS `address`,`a`.`postal_code` AS `zip code`,`a`.`phone` AS `phone`,`city`.`city` AS `city`,`country`.`country` AS `country`,`s`.`store_id` AS `SID` from (((`staff` `s` join `address` `a` on((`s`.`address_id` = `a`.`address_id`))) join `city` on((`a`.`city_id` = `city`.`city_id`))) join `country` on((`city`.`country_id` = `country`.`country_id`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| sakila
< result
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "EVENT_OBJECT_TABLE" "VARCHAR" "string"
< column "EVENT_MANIPULATION" "ENUM" "string"
< column "ACTION_TIMING" "ENUM" "string"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< row
: bytes
| ins_film
: bytes
| film
: bytes
| INSERT
: bytes
| AFTER
: bytes
| BEGIN
|     INSERT INTO film_text (film_id, title, description)
|         VALUES (new.film_id, new.title, new.description);
|   END
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
< row
: bytes
| upd_film
: bytes
| film
: bytes
| UPDATE
: bytes
| AFTER
: bytes
| BEGIN
|     IF (old.title != new.title) or (old.description != new.description)
|     THEN
|         UPDATE film_text
|             SET title=new.title,
|                 description=new.description,
|                 film_id=new.film_id
|         WHERE film_id=old.film_id;
|     END IF;
|   END
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
< row
: bytes
| del_film
: bytes
| film
: bytes
| DELETE
: bytes
| AFTER
: bytes
| BEGIN
|     DELETE FROM film_text WHERE film_id = old.film_id;
|   END
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| sakila
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
| sakila
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
< row
: bytes
| address
: bytes
| fk_address_city
: bytes
| city
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| city
: bytes
| fk_city_country
: bytes
| country
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| customer
: bytes
| fk_customer_address
: bytes
| address
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| customer
: bytes
| fk_customer_store
: bytes
| store
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| film_actor
: bytes
| fk_film_actor_actor
: bytes
| actor
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| film_actor
: bytes
| fk_film_actor_film
: bytes
| film
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| film_category
: bytes
| fk_film_category_category
: bytes
| category
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| film_category
: bytes
| fk_film_category_film
: bytes
| film
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| film
: bytes
| fk_film_language
: bytes
| language
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| film
: bytes
| fk_film_language_original
: bytes
| language
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| inventory
: bytes
| fk_inventory_film
: bytes
| film
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| inventory
: bytes
| fk_inventory_store
: bytes
| store
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| payment
: bytes
| fk_payment_customer
: bytes
| customer
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| payment
: bytes
| fk_payment_rental
: bytes
| rental
: bytes
| SET NULL
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| payment
: bytes
| fk_payment_staff
: bytes
| staff
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| rental
: bytes
| fk_rental_customer
: bytes
| customer
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| rental
: bytes
| fk_rental_inventory
: bytes
| inventory
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| rental
: bytes
| fk_rental_staff
: bytes
| staff
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| staff
: bytes
| fk_staff_address
: bytes
| address
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| staff
: bytes
| fk_staff_store
: bytes
| store
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| store
: bytes
| fk_store_address
: bytes
| address
: bytes
| RESTRICT
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| store
: bytes
| fk_store_staff
: bytes
| staff
: bytes
| RESTRICT
: bytes
| CASCADE
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
| sakila
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| sakila
: bytes
| address
: bytes
| fk_address_city
: bytes
| city_id
: bytes
| sakila
: bytes
| city_id
< row
: bytes
| sakila
: bytes
| city
: bytes
| fk_city_country
: bytes
| country_id
: bytes
| sakila
: bytes
| country_id
< row
: bytes
| sakila
: bytes
| customer
: bytes
| fk_customer_address
: bytes
| address_id
: bytes
| sakila
: bytes
| address_id
< row
: bytes
| sakila
: bytes
| customer
: bytes
| fk_customer_store
: bytes
| store_id
: bytes
| sakila
: bytes
| store_id
< row
: bytes
| sakila
: bytes
| film
: bytes
| fk_film_language
: bytes
| language_id
: bytes
| sakila
: bytes
| language_id
< row
: bytes
| sakila
: bytes
| film
: bytes
| fk_film_language_original
: bytes
| original_language_id
: bytes
| sakila
: bytes
| language_id
< row
: bytes
| sakila
: bytes
| film_actor
: bytes
| fk_film_actor_actor
: bytes
| actor_id
: bytes
| sakila
: bytes
| actor_id
< row
: bytes
| sakila
: bytes
| film_actor
: bytes
| fk_film_actor_film
: bytes
| film_id
: bytes
| sakila
: bytes
| film_id
< row
: bytes
| sakila
: bytes
| film_category
: bytes
| fk_film_category_category
: bytes
| category_id
: bytes
| sakila
: bytes
| category_id
< row
: bytes
| sakila
: bytes
| film_category
: bytes
| fk_film_category_film
: bytes
| film_id
: bytes
| sakila
: bytes
| film_id
< row
: bytes
| sakila
: bytes
| inventory
: bytes
| fk_inventory_film
: bytes
| film_id
: bytes
| sakila
: bytes
| film_id
< row
: bytes
| sakila
: bytes
| inventory
: bytes
| fk_inventory_store
: bytes
| store_id
: bytes
| sakila
: bytes
| store_id
< row
: bytes
| sakila
: bytes
| payment
: bytes
| fk_payment_customer
: bytes
| customer_id
: bytes
| sakila
: bytes
| customer_id
< row
: bytes
| sakila
: bytes
| payment
: bytes
| fk_payment_rental
: bytes
| rental_id
: bytes
| sakila
: bytes
| rental_id
< row
: bytes
| sakila
: bytes
| payment
: bytes
| fk_payment_staff
: bytes
| staff_id
: bytes
| sakila
: bytes
| staff_id
< row
: bytes
| sakila
: bytes
| rental
: bytes
| fk_rental_customer
: bytes
| customer_id
: bytes
| sakila
: bytes
| customer_id
< row
: bytes
| sakila
: bytes
| rental
: bytes
| fk_rental_inventory
: bytes
| inventory_id
: bytes
| sakila
: bytes
| inventory_id
< row
: bytes
| sakila
: bytes
| rental
: bytes
| fk_rental_staff
: bytes
| staff_id
: bytes
| sakila
: bytes
| staff_id
< row
: bytes
| sakila
: bytes
| staff
: bytes
| fk_staff_address
: bytes
| address_id
: bytes
| sakila
: bytes
| address_id
< row
: bytes
| sakila
: bytes
| staff
: bytes
| fk_staff_store
: bytes
| store_id
: bytes
| sakila
: bytes
| store_id
< row
: bytes
| sakila
: bytes
| store
: bytes
| fk_store_address
: bytes
| address_id
: bytes
| sakila
: bytes
| address_id
< row
: bytes
| sakila
: bytes
| store
: bytes
| fk_store_staff
: bytes
| manager_staff_id
: bytes
| sakila
: bytes
| staff_id
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
| sakila
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
| sakila
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "ENUM" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT" "string"
< row
: bytes
| get_customer_balance
: bytes
| FUNCTION
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
|
< row
: bytes
| inventory_held_by_customer
: bytes
| FUNCTION
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
|
< row
: bytes
| inventory_in_stock
: bytes
| FUNCTION
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
|
< row
: bytes
| film_in_stock
: bytes
| PROCEDURE
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
|
< row
: bytes
| film_not_in_stock
: bytes
| PROCEDURE
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
|
< row
: bytes
| rewards_report
: bytes
| PROCEDURE
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
| Provides a customizable report on best customers
> query
| SHOW CREATE FUNCTION `sakila`.`get_customer_balance`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| get_customer_balance
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `get_customer_balance`(p_customer_id INT, p_effective_date DATETIME) RETURNS decimal(5,2)
|     READS SQL DATA
|     DETERMINISTIC
| BEGIN
|
|        #OK, WE NEED TO CALCULATE THE CURRENT BALANCE GIVEN A CUSTOMER_ID AND A DATE
|        #THAT WE WANT THE BALANCE TO BE EFFECTIVE FOR. THE BALANCE IS:
|        #   1) RENTAL FEES FOR ALL PREVIOUS RENTALS
|        #   2) ONE DOLLAR FOR EVERY DAY THE PREVIOUS RENTALS ARE OVERDUE
|        #   3) IF A FILM IS MORE THAN RENTAL_DURATION * 2 OVERDUE, CHARGE THE REPLACEMENT_COST
|        #   4) SUBTRACT ALL PAYMENTS MADE BEFORE THE DATE SPECIFIED
|
|   DECLARE v_rentfees DECIMAL(5,2); #FEES PAID TO RENT THE VIDEOS INITIALLY
|   DECLARE v_overfees INTEGER;      #LATE FEES FOR PRIOR RENTALS
|   DECLARE v_payments DECIMAL(5,2); #SUM OF PAYMENTS MADE PREVIOUSLY
|
|   SELECT IFNULL(SUM(film.rental_rate),0) INTO v_rentfees
|     FROM film, inventory, rental
|     WHERE film.film_id = inventory.film_id
|       AND inventory.inventory_id = rental.inventory_id
|       AND rental.rental_date <= p_effective_date
|       AND rental.customer_id = p_customer_id;
|
|   SELECT IFNULL(SUM(IF((TO_DAYS(rental.return_date) - TO_DAYS(rental.rental_date)) > film.rental_duration,
|         ((TO_DAYS(rental.return_date) - TO_DAYS(rental.rental_date)) - film.rental_duration),0)),0) INTO v_overfees
|     FROM rental, inventory, film
|     WHERE film.film_id = inventory.film_id
|       AND inventory.inventory_id = rental.inventory_id
|       AND rental.rental_date <= p_effective_date
|       AND rental.customer_id = p_customer_id;
|
|
|   SELECT IFNULL(SUM(payment.amount),0) INTO v_payments
|     FROM payment
|
|     WHERE payment.payment_date <= p_effective_date
|     AND payment.customer_id = p_customer_id;
|
|   RETURN v_rentfees + v_overfees - v_payments;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE FUNCTION `sakila`.`inventory_held_by_customer`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| inventory_held_by_customer
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `inventory_held_by_customer`(p_inventory_id INT) RETURNS int
|     READS SQL DATA
| BEGIN
|   DECLARE v_customer_id INT;
|   DECLARE EXIT HANDLER FOR NOT FOUND RETURN NULL;
|
|   SELECT customer_id INTO v_customer_id
|   FROM rental
|   WHERE return_date IS NULL
|   AND inventory_id = p_inventory_id;
|
|   RETURN v_customer_id;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE FUNCTION `sakila`.`inventory_in_stock`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| inventory_in_stock
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `inventory_in_stock`(p_inventory_id INT) RETURNS tinyint(1)
|     READS SQL DATA
| BEGIN
|     DECLARE v_rentals INT;
|     DECLARE v_out     INT;
|
|     #AN ITEM IS IN-STOCK IF THERE ARE EITHER NO ROWS IN THE rental TABLE
|     #FOR THE ITEM OR ALL ROWS HAVE return_date POPULATED
|
|     SELECT COUNT(*) INTO v_rentals
|     FROM rental
|     WHERE inventory_id = p_inventory_id;
|
|     IF v_rentals = 0 THEN
|       RETURN TRUE;
|     END IF;
|
|     SELECT COUNT(rental_id) INTO v_out
|     FROM inventory LEFT JOIN rental USING(inventory_id)
|     WHERE inventory.inventory_id = p_inventory_id
|     AND rental.return_date IS NULL;
|
|     IF v_out > 0 THEN
|       RETURN FALSE;
|     ELSE
|       RETURN TRUE;
|     END IF;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `sakila`.`film_in_stock`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| film_in_stock
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `film_in_stock`(IN p_film_id INT, IN p_store_id INT, OUT p_film_count INT)
|     READS SQL DATA
| BEGIN
|      SELECT inventory_id
|      FROM inventory
|      WHERE film_id = p_film_id
|      AND store_id = p_store_id
|      AND inventory_in_stock(inventory_id);
|
|      SELECT FOUND_ROWS() INTO p_film_count;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `sakila`.`film_not_in_stock`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| film_not_in_stock
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `film_not_in_stock`(IN p_film_id INT, IN p_store_id INT, OUT p_film_count INT)
|     READS SQL DATA
| BEGIN
|      SELECT inventory_id
|      FROM inventory
|      WHERE film_id = p_film_id
|      AND store_id = p_store_id
|      AND NOT inventory_in_stock(inventory_id);
|
|      SELECT FOUND_ROWS() INTO p_film_count;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `sakila`.`rewards_report`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| rewards_report
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `rewards_report`(
|     IN min_monthly_purchases TINYINT UNSIGNED
|     , IN min_dollar_amount_purchased DECIMAL(10,2) UNSIGNED
|     , OUT count_rewardees INT
| )
|     READS SQL DATA
|     COMMENT 'Provides a customizable report on best customers'
| proc: BEGIN
|
|     DECLARE last_month_start DATE;
|     DECLARE last_month_end DATE;
|
|     /* Some sanity checks... */
|     IF min_monthly_purchases = 0 THEN
|         SELECT 'Minimum monthly purchases parameter must be > 0';
|         LEAVE proc;
|     END IF;
|     IF min_dollar_amount_purchased = 0.00 THEN
|         SELECT 'Minimum monthly dollar amount purchased parameter must be > $0.00';
|         LEAVE proc;
|     END IF;
|
|     /* Determine start and end time periods */
|     SET last_month_start = DATE_SUB(CURRENT_DATE(), INTERVAL 1 MONTH);
|     SET last_month_start = STR_TO_DATE(CONCAT(YEAR(last_month_start),'-',MONTH(last_month_start),'-01'),'%Y-%m-%d');
|     SET last_month_end = LAST_DAY(last_month_start);
|
|     /*
|         Create a temporary storage area for
|         Customer IDs.
|     */
|     CREATE TEMPORARY TABLE tmpCustomer (customer_id INT UNSIGNED NOT NULL PRIMARY KEY);
|
|     /*
|         Find all customers meeting the
|         monthly purchase requirements
|     */
|     INSERT INTO tmpCustomer (customer_id)
|     SELECT p.customer_id
|     FROM payment AS p
|     WHERE DATE(p.payment_date) BETWEEN last_month_start AND last_month_end
|     GROUP BY customer_id
|     HAVING SUM(p.amount) > min_dollar_amount_purchased
|     AND COUNT(customer_id) > min_monthly_purchases;
|
|     /* Populate OUT parameter with count of found customers */
|     SELECT COUNT(*) FROM tmpCustomer INTO count_rewardees;
|
|     /*
|         Output ALL customer information of matching rewardees.
|         Customize output as needed.
|     */
|     SELECT c.*
|     FROM tmpCustomer AS t
|     INNER JOIN customer AS c ON t.customer_id = c.customer_id;
|
|     /* Clean up */
|     DROP TABLE tmpCustomer;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
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
| sakila
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
| actor
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| actor_info
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| address
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| category
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| city
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| country
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| customer
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
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
| utf8mb3
< row
: bytes
| customer_list
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| film
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
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
| utf8mb3
< row
: bytes
| film_actor
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| film_category
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| film_list
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| film_text
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 32768
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| inventory
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 32768
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| language
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| nicer_but_slower_film_list
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| payment
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
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
| utf8mb3
< row
: bytes
| rental
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 65536
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| sales_by_film_category
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| sales_by_store
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| staff
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 32768
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| staff_list
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| store
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 32768
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
> query
|
| 		SELECT t.NAME
| 		FROM information_schema.INNODB_TABLES t
| 		LEFT JOIN information_schema.INNODB_TABLESPACES ts ON t.SPACE = ts.SPACE
| 		WHERE t.NAME LIKE CONCAT(?, '/%') AND (ts.SPACE_TYPE IS NULL OR ts.SPACE_TYPE <> 'Single')
: string
| sakila
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
| sakila
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci

== TestSDLRealWorldIdempotence default#3
> begin
< ok
> exec
| --  Sample employee database 
| --  See changelog table for details
| --  Copyright (C) 2007,2008, MySQL AB
| --  
| --  Original data created by Fusheng Wang and Carlo Zaniolo
| --  http://www.cs.aau.dk/TimeCenter/software.htm
| --  http://www.cs.aau.dk/TimeCenter/Data/employeeTemporalDataSet.zip
| -- 
| --  Current schema by Giuseppe Maxia 
| --  Data conversion from XML to relational by Patrick Crews
| -- 
| -- This work is licensed under the 
| -- Creative Commons Attribution-Share Alike 3.0 Unported License. 
| -- To view a copy of this license, visit 
| -- http://creativecommons.org/licenses/by-sa/3.0/ or send a letter to 
| -- Creative Commons, 171 Second Street, Suite 300, San Francisco, 
| -- California, 94105, USA.
| -- 
| --  DISCLAIMER
| --  To the best of our knowledge, this data is fabricated, and
| --  it does not correspond to real people. 
| --  Any similarity to existing people is purely coincidental.
| -- 
|
|
|
|
|
| /*!50503 set default_storage_engine = InnoDB */;
< ok 0 0
> exec
|
| /*!50503 select CONCAT('storage engine: ', @@default_storage_engine) as INFO */;
< ok 0 0
> exec
|
|
| CREATE TABLE employees (
|     emp_no      INT             NOT NULL,
|     birth_date  DATE            NOT NULL,
|     first_name  VARCHAR(14)     NOT NULL,
|     last_name   VARCHAR(16)     NOT NULL,
|     gender      ENUM ('M','F')  NOT NULL,    
|     hire_date   DATE            NOT NULL,
|     PRIMARY KEY (emp_no)
| );
< ok 0 0
> exec
|
|
| CREATE TABLE departments (
|     dept_no     CHAR(4)         NOT NULL,
|     dept_name   VARCHAR(40)     NOT NULL,
|     PRIMARY KEY (dept_no),
|     UNIQUE  KEY (dept_name)
| );
< ok 0 0
> exec
|
|
| CREATE TABLE dept_manager (
|    emp_no       INT             NOT NULL,
|    dept_no      CHAR(4)         NOT NULL,
|    from_date    DATE            NOT NULL,
|    to_date      DATE            NOT NULL,
|    FOREIGN KEY (emp_no)  REFERENCES employees (emp_no)    ON DELETE CASCADE,
|    FOREIGN KEY (dept_no) REFERENCES departments (dept_no) ON DELETE CASCADE,
|    PRIMARY KEY (emp_no,dept_no)
| );
< ok 0 0
> exec
|  
|
| CREATE TABLE dept_emp (
|     emp_no      INT             NOT NULL,
|     dept_no     CHAR(4)         NOT NULL,
|     from_date   DATE            NOT NULL,
|     to_date     DATE            NOT NULL,
|     FOREIGN KEY (emp_no)  REFERENCES employees   (emp_no)  ON DELETE CASCADE,
|     FOREIGN KEY (dept_no) REFERENCES departments (dept_no) ON DELETE CASCADE,
|     PRIMARY KEY (emp_no,dept_no)
| );
< ok 0 0
> exec
|
|
| CREATE TABLE titles (
|     emp_no      INT             NOT NULL,
|     title       VARCHAR(50)     NOT NULL,
|     from_date   DATE            NOT NULL,
|     to_date     DATE,
|     FOREIGN KEY (emp_no) REFERENCES employees (emp_no) ON DELETE CASCADE,
|     PRIMARY KEY (emp_no,title, from_date)
| ) 
| ;
< ok 0 0
> exec
|  
|
| CREATE TABLE salaries (
|     emp_no      INT             NOT NULL,
|     salary      INT             NOT NULL,
|     from_date   DATE            NOT NULL,
|     to_date     DATE            NOT NULL,
|     FOREIGN KEY (emp_no) REFERENCES employees (emp_no) ON DELETE CASCADE,
|     PRIMARY KEY (emp_no, from_date)
| ) 
| ;
< ok 0 0
> exec
|  
|
| CREATE OR REPLACE VIEW dept_emp_latest_date AS
|     SELECT emp_no, MAX(from_date) AS from_date, MAX(to_date) AS to_date
|     FROM dept_emp
|     GROUP BY emp_no;
< ok 0 0
> exec
|
|
| # shows only the current department for each employee
| CREATE OR REPLACE VIEW current_dept_emp AS
|     SELECT l.emp_no, dept_no, l.from_date, l.to_date
|     FROM dept_emp d
|         INNER JOIN dept_emp_latest_date l
|         ON d.emp_no=l.emp_no AND d.from_date=l.from_date AND l.to_date = d.to_date;
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
| sdl_rw_employees_1
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
| departments
: bytes
| dept_name
: bytes
| dept_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| departments
: bytes
| PRIMARY
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dept_emp
: bytes
| dept_no
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| dept_emp
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dept_emp
: bytes
| PRIMARY
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dept_manager
: bytes
| dept_no
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| dept_manager
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dept_manager
: bytes
| PRIMARY
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| employees
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| salaries
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| salaries
: bytes
| PRIMARY
: bytes
| from_date
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| titles
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| titles
: bytes
| PRIMARY
: bytes
| title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| titles
: bytes
| PRIMARY
: bytes
| from_date
: bytes
| A
: int64 -1
: null
: int64 3
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
| sdl_rw_employees_1
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
| current_dept_emp
: bytes
| emp_no
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
|
: null
< row
: bytes
| current_dept_emp
: bytes
| dept_no
: int64 2
: null
: bytes
| NO
: bytes
| char(4)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| current_dept_emp
: bytes
| from_date
: int64 3
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| current_dept_emp
: bytes
| to_date
: int64 4
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| departments
: bytes
| dept_no
: int64 1
: null
: bytes
| NO
: bytes
| char(4)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| departments
: bytes
| dept_name
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
| dept_emp
: bytes
| emp_no
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
|
: null
< row
: bytes
| dept_emp
: bytes
| dept_no
: int64 2
: null
: bytes
| NO
: bytes
| char(4)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp
: bytes
| from_date
: int64 3
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp
: bytes
| to_date
: int64 4
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp_latest_date
: bytes
| emp_no
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
|
: null
< row
: bytes
| dept_emp_latest_date
: bytes
| from_date
: int64 2
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp_latest_date
: bytes
| to_date
: int64 3
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_manager
: bytes
| emp_no
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
|
: null
< row
: bytes
| dept_manager
: bytes
| dept_no
: int64 2
: null
: bytes
| NO
: bytes
| char(4)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_manager
: bytes
| from_date
: int64 3
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_manager
: bytes
| to_date
: int64 4
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| emp_no
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
|
: null
< row
: bytes
| employees
: bytes
| birth_date
: int64 2
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| first_name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(14)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| last_name
: int64 4
: null
: bytes
| NO
: bytes
| varchar(16)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| gender
: int64 5
: null
: bytes
| NO
: bytes
| enum('M','F')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| hire_date
: int64 6
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| salaries
: bytes
| emp_no
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
|
: null
< row
: bytes
| salaries
: bytes
| salary
: int64 2
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
|
: null
< row
: bytes
| salaries
: bytes
| from_date
: int64 3
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| salaries
: bytes
| to_date
: int64 4
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| titles
: bytes
| emp_no
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
|
: null
< row
: bytes
| titles
: bytes
| title
: int64 2
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| titles
: bytes
| from_date
: int64 3
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| titles
: bytes
| to_date
: int64 4
: null
: bytes
| YES
: bytes
| date
: bytes
|
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
| sdl_rw_employees_1
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
| sdl_rw_employees_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| current_dept_emp
: bytes
| select `sdl_rw_employees_1`.`l`.`emp_no` AS `emp_no`,`d`.`dept_no` AS `dept_no`,`sdl_rw_employees_1`.`l`.`from_date` AS `from_date`,`sdl_rw_employees_1`.`l`.`to_date` AS `to_date` from (`sdl_rw_employees_1`.`dept_emp` `d` join `sdl_rw_employees_1`.`dept_emp_latest_date` `l` on(((`d`.`emp_no` = `sdl_rw_employees_1`.`l`.`emp_no`) and (`d`.`from_date` = `sdl_rw_employees_1`.`l`.`from_date`) and (`sdl_rw_employees_1`.`l`.`to_date` = `d`.`to_date`))))
< row
: bytes
| dept_emp_latest_date
: bytes
| select `sdl_rw_employees_1`.`dept_emp`.`emp_no` AS `emp_no`,max(`sdl_rw_employees_1`.`dept_emp`.`from_date`) AS `from_date`,max(`sdl_rw_employees_1`.`dept_emp`.`to_date`) AS `to_date` from `sdl_rw_employees_1`.`dept_emp` group by `sdl_rw_employees_1`.`dept_emp`.`emp_no`
> query
| SHOW CREATE VIEW `sdl_rw_employees_1`.`current_dept_emp`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| current_dept_emp
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `current_dept_emp` AS select `l`.`emp_no` AS `emp_no`,`d`.`dept_no` AS `dept_no`,`l`.`from_date` AS `from_date`,`l`.`to_date` AS `to_date` from (`dept_emp` `d` join `dept_emp_latest_date` `l` on(((`d`.`emp_no` = `l`.`emp_no`) and (`d`.`from_date` = `l`.`from_date`) and (`l`.`to_date` = `d`.`to_date`))))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `sdl_rw_employees_1`.`dept_emp_latest_date`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| dept_emp_latest_date
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `dept_emp_latest_date` AS select `dept_emp`.`emp_no` AS `emp_no`,max(`dept_emp`.`from_date`) AS `from_date`,max(`dept_emp`.`to_date`) AS `to_date` from `dept_emp` group by `dept_emp`.`emp_no`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| sdl_rw_employees_1
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
| sdl_rw_employees_1
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
| sdl_rw_employees_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
< row
: bytes
| dept_emp
: bytes
| dept_emp_ibfk_1
: bytes
| employees
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| dept_emp
: bytes
| dept_emp_ibfk_2
: bytes
| departments
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| dept_manager
: bytes
| dept_manager_ibfk_1
: bytes
| employees
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| dept_manager
: bytes
| dept_manager_ibfk_2
: bytes
| departments
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| salaries
: bytes
| salaries_ibfk_1
: bytes
| employees
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| titles
: bytes
| titles_ibfk_1
: bytes
| employees
: bytes
| CASCADE
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
| sdl_rw_employees_1
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| sdl_rw_employees_1
: bytes
| dept_emp
: bytes
| dept_emp_ibfk_1
: bytes
| emp_no
: bytes
| sdl_rw_employees_1
: bytes
| emp_no
< row
: bytes
| sdl_rw_employees_1
: bytes
| dept_emp
: bytes
| dept_emp_ibfk_2
: bytes
| dept_no
: bytes
| sdl_rw_employees_1
: bytes
| dept_no
< row
: bytes
| sdl_rw_employees_1
: bytes
| dept_manager
: bytes
| dept_manager_ibfk_1
: bytes
| emp_no
: bytes
| sdl_rw_employees_1
: bytes
| emp_no
< row
: bytes
| sdl_rw_employees_1
: bytes
| dept_manager
: bytes
| dept_manager_ibfk_2
: bytes
| dept_no
: bytes
| sdl_rw_employees_1
: bytes
| dept_no
< row
: bytes
| sdl_rw_employees_1
: bytes
| salaries
: bytes
| salaries_ibfk_1
: bytes
| emp_no
: bytes
| sdl_rw_employees_1
: bytes
| emp_no
< row
: bytes
| sdl_rw_employees_1
: bytes
| titles
: bytes
| titles_ibfk_1
: bytes
| emp_no
: bytes
| sdl_rw_employees_1
: bytes
| emp_no
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
| sdl_rw_employees_1
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
| sdl_rw_employees_1
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
| sdl_rw_employees_1
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
| current_dept_emp
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| departments
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
| 16384
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
| dept_emp
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
| 16384
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
| dept_emp_latest_date
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| dept_manager
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
| 16384
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
| employees
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
| 0
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
| salaries
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
| 0
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
| titles
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
| 0
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
| sdl_rw_employees_1
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
| sdl_rw_employees_1
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci

== TestSDLRealWorldIdempotence default#4
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
| sdl_rw_employees_1
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
| departments
: bytes
| dept_name
: bytes
| dept_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| departments
: bytes
| PRIMARY
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dept_emp
: bytes
| dept_no
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| dept_emp
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dept_emp
: bytes
| PRIMARY
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dept_manager
: bytes
| dept_no
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| dept_manager
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dept_manager
: bytes
| PRIMARY
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| employees
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| salaries
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| salaries
: bytes
| PRIMARY
: bytes
| from_date
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| titles
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| titles
: bytes
| PRIMARY
: bytes
| title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| titles
: bytes
| PRIMARY
: bytes
| from_date
: bytes
| A
: int64 -1
: null
: int64 3
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
| sdl_rw_employees_1
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
| current_dept_emp
: bytes
| emp_no
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
|
: null
< row
: bytes
| current_dept_emp
: bytes
| dept_no
: int64 2
: null
: bytes
| NO
: bytes
| char(4)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| current_dept_emp
: bytes
| from_date
: int64 3
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| current_dept_emp
: bytes
| to_date
: int64 4
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| departments
: bytes
| dept_no
: int64 1
: null
: bytes
| NO
: bytes
| char(4)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| departments
: bytes
| dept_name
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
| dept_emp
: bytes
| emp_no
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
|
: null
< row
: bytes
| dept_emp
: bytes
| dept_no
: int64 2
: null
: bytes
| NO
: bytes
| char(4)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp
: bytes
| from_date
: int64 3
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp
: bytes
| to_date
: int64 4
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp_latest_date
: bytes
| emp_no
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
|
: null
< row
: bytes
| dept_emp_latest_date
: bytes
| from_date
: int64 2
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp_latest_date
: bytes
| to_date
: int64 3
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_manager
: bytes
| emp_no
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
|
: null
< row
: bytes
| dept_manager
: bytes
| dept_no
: int64 2
: null
: bytes
| NO
: bytes
| char(4)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_manager
: bytes
| from_date
: int64 3
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_manager
: bytes
| to_date
: int64 4
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| emp_no
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
|
: null
< row
: bytes
| employees
: bytes
| birth_date
: int64 2
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| first_name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(14)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| last_name
: int64 4
: null
: bytes
| NO
: bytes
| varchar(16)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| gender
: int64 5
: null
: bytes
| NO
: bytes
| enum('M','F')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| hire_date
: int64 6
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| salaries
: bytes
| emp_no
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
|
: null
< row
: bytes
| salaries
: bytes
| salary
: int64 2
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
|
: null
< row
: bytes
| salaries
: bytes
| from_date
: int64 3
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| salaries
: bytes
| to_date
: int64 4
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| titles
: bytes
| emp_no
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
|
: null
< row
: bytes
| titles
: bytes
| title
: int64 2
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| titles
: bytes
| from_date
: int64 3
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| titles
: bytes
| to_date
: int64 4
: null
: bytes
| YES
: bytes
| date
: bytes
|
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
| sdl_rw_employees_1
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
| sdl_rw_employees_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| current_dept_emp
: bytes
| select `sdl_rw_employees_1`.`l`.`emp_no` AS `emp_no`,`d`.`dept_no` AS `dept_no`,`sdl_rw_employees_1`.`l`.`from_date` AS `from_date`,`sdl_rw_employees_1`.`l`.`to_date` AS `to_date` from (`sdl_rw_employees_1`.`dept_emp` `d` join `sdl_rw_employees_1`.`dept_emp_latest_date` `l` on(((`d`.`emp_no` = `sdl_rw_employees_1`.`l`.`emp_no`) and (`d`.`from_date` = `sdl_rw_employees_1`.`l`.`from_date`) and (`sdl_rw_employees_1`.`l`.`to_date` = `d`.`to_date`))))
< row
: bytes
| dept_emp_latest_date
: bytes
| select `sdl_rw_employees_1`.`dept_emp`.`emp_no` AS `emp_no`,max(`sdl_rw_employees_1`.`dept_emp`.`from_date`) AS `from_date`,max(`sdl_rw_employees_1`.`dept_emp`.`to_date`) AS `to_date` from `sdl_rw_employees_1`.`dept_emp` group by `sdl_rw_employees_1`.`dept_emp`.`emp_no`
> query
| SHOW CREATE VIEW `sdl_rw_employees_1`.`current_dept_emp`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| current_dept_emp
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `current_dept_emp` AS select `l`.`emp_no` AS `emp_no`,`d`.`dept_no` AS `dept_no`,`l`.`from_date` AS `from_date`,`l`.`to_date` AS `to_date` from (`dept_emp` `d` join `dept_emp_latest_date` `l` on(((`d`.`emp_no` = `l`.`emp_no`) and (`d`.`from_date` = `l`.`from_date`) and (`l`.`to_date` = `d`.`to_date`))))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `sdl_rw_employees_1`.`dept_emp_latest_date`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| dept_emp_latest_date
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `dept_emp_latest_date` AS select `dept_emp`.`emp_no` AS `emp_no`,max(`dept_emp`.`from_date`) AS `from_date`,max(`dept_emp`.`to_date`) AS `to_date` from `dept_emp` group by `dept_emp`.`emp_no`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| sdl_rw_employees_1
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
| sdl_rw_employees_1
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
| sdl_rw_employees_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
< row
: bytes
| dept_emp
: bytes
| dept_emp_ibfk_1
: bytes
| employees
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| dept_emp
: bytes
| dept_emp_ibfk_2
: bytes
| departments
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| dept_manager
: bytes
| dept_manager_ibfk_1
: bytes
| employees
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| dept_manager
: bytes
| dept_manager_ibfk_2
: bytes
| departments
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| salaries
: bytes
| salaries_ibfk_1
: bytes
| employees
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| titles
: bytes
| titles_ibfk_1
: bytes
| employees
: bytes
| CASCADE
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
| sdl_rw_employees_1
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| sdl_rw_employees_1
: bytes
| dept_emp
: bytes
| dept_emp_ibfk_1
: bytes
| emp_no
: bytes
| sdl_rw_employees_1
: bytes
| emp_no
< row
: bytes
| sdl_rw_employees_1
: bytes
| dept_emp
: bytes
| dept_emp_ibfk_2
: bytes
| dept_no
: bytes
| sdl_rw_employees_1
: bytes
| dept_no
< row
: bytes
| sdl_rw_employees_1
: bytes
| dept_manager
: bytes
| dept_manager_ibfk_1
: bytes
| emp_no
: bytes
| sdl_rw_employees_1
: bytes
| emp_no
< row
: bytes
| sdl_rw_employees_1
: bytes
| dept_manager
: bytes
| dept_manager_ibfk_2
: bytes
| dept_no
: bytes
| sdl_rw_employees_1
: bytes
| dept_no
< row
: bytes
| sdl_rw_employees_1
: bytes
| salaries
: bytes
| salaries_ibfk_1
: bytes
| emp_no
: bytes
| sdl_rw_employees_1
: bytes
| emp_no
< row
: bytes
| sdl_rw_employees_1
: bytes
| titles
: bytes
| titles_ibfk_1
: bytes
| emp_no
: bytes
| sdl_rw_employees_1
: bytes
| emp_no
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
| sdl_rw_employees_1
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
| sdl_rw_employees_1
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
| sdl_rw_employees_1
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
| current_dept_emp
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| departments
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
| 16384
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
| dept_emp
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
| 16384
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
| dept_emp_latest_date
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| dept_manager
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
| 16384
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
| employees
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
| 0
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
| salaries
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
| 0
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
| titles
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
| 0
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
| sdl_rw_employees_1
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
| sdl_rw_employees_1
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci

== TestSDLRealWorldIdempotence default#5
> begin
< ok
> exec
| --  Sample employee database 
| --  See changelog table for details
| --  Copyright (C) 2007,2008, MySQL AB
| --  
| --  Original data created by Fusheng Wang and Carlo Zaniolo
| --  http://www.cs.aau.dk/TimeCenter/software.htm
| --  http://www.cs.aau.dk/TimeCenter/Data/employeeTemporalDataSet.zip
| -- 
| --  Current schema by Giuseppe Maxia 
| --  Data conversion from XML to relational by Patrick Crews
| -- 
| -- This work is licensed under the 
| -- Creative Commons Attribution-Share Alike 3.0 Unported License. 
| -- To view a copy of this license, visit 
| -- http://creativecommons.org/licenses/by-sa/3.0/ or send a letter to 
| -- Creative Commons, 171 Second Street, Suite 300, San Francisco, 
| -- California, 94105, USA.
| -- 
| --  DISCLAIMER
| --  To the best of our knowledge, this data is fabricated, and
| --  it does not correspond to real people. 
| --  Any similarity to existing people is purely coincidental.
| -- 
|
|
|
|
|
| /*!50503 set default_storage_engine = InnoDB */;
< ok 0 0
> exec
|
| /*!50503 select CONCAT('storage engine: ', @@default_storage_engine) as INFO */;
< ok 0 0
> exec
|
|
| CREATE TABLE employees (
|     emp_no      INT             NOT NULL,
|     birth_date  DATE            NOT NULL,
|     first_name  VARCHAR(14)     NOT NULL,
|     last_name   VARCHAR(16)     NOT NULL,
|     gender      ENUM ('M','F')  NOT NULL,    
|     hire_date   DATE            NOT NULL,
|     PRIMARY KEY (emp_no)
| );
< ok 0 0
> exec
|
|
| CREATE TABLE departments (
|     dept_no     CHAR(4)         NOT NULL,
|     dept_name   VARCHAR(40)     NOT NULL,
|     PRIMARY KEY (dept_no),
|     UNIQUE  KEY (dept_name)
| );
< ok 0 0
> exec
|
|
| CREATE TABLE dept_manager (
|    emp_no       INT             NOT NULL,
|    dept_no      CHAR(4)         NOT NULL,
|    from_date    DATE            NOT NULL,
|    to_date      DATE            NOT NULL,
|    FOREIGN KEY (emp_no)  REFERENCES employees (emp_no)    ON DELETE CASCADE,
|    FOREIGN KEY (dept_no) REFERENCES departments (dept_no) ON DELETE CASCADE,
|    PRIMARY KEY (emp_no,dept_no)
| );
< ok 0 0
> exec
|  
|
| CREATE TABLE dept_emp (
|     emp_no      INT             NOT NULL,
|     dept_no     CHAR(4)         NOT NULL,
|     from_date   DATE            NOT NULL,
|     to_date     DATE            NOT NULL,
|     FOREIGN KEY (emp_no)  REFERENCES employees   (emp_no)  ON DELETE CASCADE,
|     FOREIGN KEY (dept_no) REFERENCES departments (dept_no) ON DELETE CASCADE,
|     PRIMARY KEY (emp_no,dept_no)
| );
< ok 0 0
> exec
|
|
| CREATE TABLE titles (
|     emp_no      INT             NOT NULL,
|     title       VARCHAR(50)     NOT NULL,
|     from_date   DATE            NOT NULL,
|     to_date     DATE,
|     # FOREIGN KEY (emp_no) REFERENCES employees (emp_no) ON DELETE CASCADE,
|     PRIMARY KEY (emp_no,title, from_date)
| );
< ok 0 0
> exec
|  
|
| /*!50510
| ALTER TABLE titles 
| partition by range COLUMNS (from_date)
| (
|     partition p01 values less than ('1985-12-31'),
|     partition p02 values less than ('1986-12-31'),
|     partition p03 values less than ('1987-12-31'),
|     partition p04 values less than ('1988-12-31'),
|     partition p05 values less than ('1989-12-31'),
|     partition p06 values less than ('1990-12-31'),
|     partition p07 values less than ('1991-12-31'),
|     partition p08 values less than ('1992-12-31'),
|     partition p09 values less than ('1993-12-31'),
|     partition p10 values less than ('1994-12-31'),
|     partition p11 values less than ('1995-12-31'),
|     partition p12 values less than ('1996-12-31'),
|     partition p13 values less than ('1997-12-31'),
|     partition p14 values less than ('1998-12-31'),
|     partition p15 values less than ('1999-12-31'),
|     partition p16 values less than ('2000-12-31'),
|     partition p17 values less than ('2001-12-31'),
|     partition p18 values less than ('2002-12-31'),
|     partition p19 values less than (MAXVALUE)
| )
| */;
< ok 0 0
> exec
|
|
| CREATE TABLE salaries (
|     emp_no      INT             NOT NULL,
|     salary      INT             NOT NULL,
|     from_date   DATE            NOT NULL,
|     to_date     DATE            NOT NULL,
|     # FOREIGN KEY (emp_no) REFERENCES employees (emp_no) ON DELETE CASCADE,
|     PRIMARY KEY (emp_no, from_date)
| );
< ok 0 0
> exec
|  
|
| /*!50510
| ALTER TABLE salaries
| partition by range COLUMNS (from_date)
| (
|     partition p01 values less than ('1985-12-31'),
|     partition p02 values less than ('1986-12-31'),
|     partition p03 values less than ('1987-12-31'),
|     partition p04 values less than ('1988-12-31'),
|     partition p05 values less than ('1989-12-31'),
|     partition p06 values less than ('1990-12-31'),
|     partition p07 values less than ('1991-12-31'),
|     partition p08 values less than ('1992-12-31'),
|     partition p09 values less than ('1993-12-31'),
|     partition p10 values less than ('1994-12-31'),
|     partition p11 values less than ('1995-12-31'),
|     partition p12 values less than ('1996-12-31'),
|     partition p13 values less than ('1997-12-31'),
|     partition p14 values less than ('1998-12-31'),
|     partition p15 values less than ('1999-12-31'),
|     partition p16 values less than ('2000-12-31'),
|     partition p17 values less than ('2001-12-31'),
|     partition p18 values less than ('2002-12-31'),
|     partition p19 values less than (MAXVALUE)
| )
| */;
< ok 0 0
> exec
|
|
| CREATE OR REPLACE VIEW dept_emp_latest_date AS
|     SELECT emp_no, MAX(from_date) AS from_date, MAX(to_date) AS to_date
|     FROM dept_emp
|     GROUP BY emp_no;
< ok 0 0
> exec
|
|
| # shows only the current department for each employee
| CREATE OR REPLACE VIEW current_dept_emp AS
|     SELECT l.emp_no, dept_no, l.from_date, l.to_date
|     FROM dept_emp d
|         INNER JOIN dept_emp_latest_date l
|         ON d.emp_no=l.emp_no AND d.from_date=l.from_date AND l.to_date = d.to_date;
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
| sdl_rw_employees_part_2
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
| departments
: bytes
| dept_name
: bytes
| dept_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| departments
: bytes
| PRIMARY
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dept_emp
: bytes
| dept_no
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| dept_emp
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dept_emp
: bytes
| PRIMARY
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dept_manager
: bytes
| dept_no
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| dept_manager
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dept_manager
: bytes
| PRIMARY
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| employees
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| salaries
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| salaries
: bytes
| PRIMARY
: bytes
| from_date
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| titles
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| titles
: bytes
| PRIMARY
: bytes
| title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| titles
: bytes
| PRIMARY
: bytes
| from_date
: bytes
| A
: int64 -1
: null
: int64 3
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
| sdl_rw_employees_part_2
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
| current_dept_emp
: bytes
| emp_no
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
|
: null
< row
: bytes
| current_dept_emp
: bytes
| dept_no
: int64 2
: null
: bytes
| NO
: bytes
| char(4)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| current_dept_emp
: bytes
| from_date
: int64 3
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| current_dept_emp
: bytes
| to_date
: int64 4
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| departments
: bytes
| dept_no
: int64 1
: null
: bytes
| NO
: bytes
| char(4)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| departments
: bytes
| dept_name
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
| dept_emp
: bytes
| emp_no
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
|
: null
< row
: bytes
| dept_emp
: bytes
| dept_no
: int64 2
: null
: bytes
| NO
: bytes
| char(4)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp
: bytes
| from_date
: int64 3
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp
: bytes
| to_date
: int64 4
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp_latest_date
: bytes
| emp_no
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
|
: null
< row
: bytes
| dept_emp_latest_date
: bytes
| from_date
: int64 2
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp_latest_date
: bytes
| to_date
: int64 3
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_manager
: bytes
| emp_no
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
|
: null
< row
: bytes
| dept_manager
: bytes
| dept_no
: int64 2
: null
: bytes
| NO
: bytes
| char(4)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_manager
: bytes
| from_date
: int64 3
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_manager
: bytes
| to_date
: int64 4
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| emp_no
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
|
: null
< row
: bytes
| employees
: bytes
| birth_date
: int64 2
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| first_name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(14)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| last_name
: int64 4
: null
: bytes
| NO
: bytes
| varchar(16)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| gender
: int64 5
: null
: bytes
| NO
: bytes
| enum('M','F')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| hire_date
: int64 6
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| salaries
: bytes
| emp_no
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
|
: null
< row
: bytes
| salaries
: bytes
| salary
: int64 2
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
|
: null
< row
: bytes
| salaries
: bytes
| from_date
: int64 3
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| salaries
: bytes
| to_date
: int64 4
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| titles
: bytes
| emp_no
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
|
: null
< row
: bytes
| titles
: bytes
| title
: int64 2
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| titles
: bytes
| from_date
: int64 3
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| titles
: bytes
| to_date
: int64 4
: null
: bytes
| YES
: bytes
| date
: bytes
|
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
| sdl_rw_employees_part_2
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
| sdl_rw_employees_part_2
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| current_dept_emp
: bytes
| select `sdl_rw_employees_part_2`.`l`.`emp_no` AS `emp_no`,`d`.`dept_no` AS `dept_no`,`sdl_rw_employees_part_2`.`l`.`from_date` AS `from_date`,`sdl_rw_employees_part_2`.`l`.`to_date` AS `to_date` from (`sdl_rw_employees_part_2`.`dept_emp` `d` join `sdl_rw_employees_part_2`.`dept_emp_latest_date` `l` on(((`d`.`emp_no` = `sdl_rw_employees_part_2`.`l`.`emp_no`) and (`d`.`from_date` = `sdl_rw_employees_part_2`.`l`.`from_date`) and (`sdl_rw_employees_part_2`.`l`.`to_date` = `d`.`to_date`))))
< row
: bytes
| dept_emp_latest_date
: bytes
| select `sdl_rw_employees_part_2`.`dept_emp`.`emp_no` AS `emp_no`,max(`sdl_rw_employees_part_2`.`dept_emp`.`from_date`) AS `from_date`,max(`sdl_rw_employees_part_2`.`dept_emp`.`to_date`) AS `to_date` from `sdl_rw_employees_part_2`.`dept_emp` group by `sdl_rw_employees_part_2`.`dept_emp`.`emp_no`
> query
| SHOW CREATE VIEW `sdl_rw_employees_part_2`.`current_dept_emp`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| current_dept_emp
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `current_dept_emp` AS select `l`.`emp_no` AS `emp_no`,`d`.`dept_no` AS `dept_no`,`l`.`from_date` AS `from_date`,`l`.`to_date` AS `to_date` from (`dept_emp` `d` join `dept_emp_latest_date` `l` on(((`d`.`emp_no` = `l`.`emp_no`) and (`d`.`from_date` = `l`.`from_date`) and (`l`.`to_date` = `d`.`to_date`))))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `sdl_rw_employees_part_2`.`dept_emp_latest_date`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| dept_emp_latest_date
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `dept_emp_latest_date` AS select `dept_emp`.`emp_no` AS `emp_no`,max(`dept_emp`.`from_date`) AS `from_date`,max(`dept_emp`.`to_date`) AS `to_date` from `dept_emp` group by `dept_emp`.`emp_no`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| sdl_rw_employees_part_2
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
| sdl_rw_employees_part_2
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
| sdl_rw_employees_part_2
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
< row
: bytes
| dept_emp
: bytes
| dept_emp_ibfk_1
: bytes
| employees
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| dept_emp
: bytes
| dept_emp_ibfk_2
: bytes
| departments
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| dept_manager
: bytes
| dept_manager_ibfk_1
: bytes
| employees
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| dept_manager
: bytes
| dept_manager_ibfk_2
: bytes
| departments
: bytes
| CASCADE
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
| sdl_rw_employees_part_2
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| sdl_rw_employees_part_2
: bytes
| dept_emp
: bytes
| dept_emp_ibfk_1
: bytes
| emp_no
: bytes
| sdl_rw_employees_part_2
: bytes
| emp_no
< row
: bytes
| sdl_rw_employees_part_2
: bytes
| dept_emp
: bytes
| dept_emp_ibfk_2
: bytes
| dept_no
: bytes
| sdl_rw_employees_part_2
: bytes
| dept_no
< row
: bytes
| sdl_rw_employees_part_2
: bytes
| dept_manager
: bytes
| dept_manager_ibfk_1
: bytes
| emp_no
: bytes
| sdl_rw_employees_part_2
: bytes
| emp_no
< row
: bytes
| sdl_rw_employees_part_2
: bytes
| dept_manager
: bytes
| dept_manager_ibfk_2
: bytes
| dept_no
: bytes
| sdl_rw_employees_part_2
: bytes
| dept_no
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
| sdl_rw_employees_part_2
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
< row
: bytes
| salaries
: bytes
| p01
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1985-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p02
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1986-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p03
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1987-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p04
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1988-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p05
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1989-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p06
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1990-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p07
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1991-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p08
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1992-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p09
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1993-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p10
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1994-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p11
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1995-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p12
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1996-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p13
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1997-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p14
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1998-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p15
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1999-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p16
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '2000-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p17
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '2001-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p18
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '2002-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p19
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| MAXVALUE
: bytes
|
< row
: bytes
| titles
: bytes
| p01
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1985-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p02
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1986-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p03
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1987-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p04
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1988-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p05
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1989-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p06
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1990-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p07
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1991-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p08
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1992-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p09
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1993-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p10
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1994-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p11
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1995-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p12
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1996-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p13
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1997-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p14
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1998-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p15
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1999-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p16
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '2000-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p17
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '2001-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p18
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '2002-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p19
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| MAXVALUE
: bytes
|
> query
| SHOW CREATE TABLE `sdl_rw_employees_part_2`.`salaries`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| salaries
: bytes
| CREATE TABLE `salaries` (
|   `emp_no` int NOT NULL,
|   `salary` int NOT NULL,
|   `from_date` date NOT NULL,
|   `to_date` date NOT NULL,
|   PRIMARY KEY (`emp_no`,`from_date`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50500 PARTITION BY RANGE  COLUMNS(from_date)
| (PARTITION p01 VALUES LESS THAN ('1985-12-31') ENGINE = InnoDB,
|  PARTITION p02 VALUES LESS THAN ('1986-12-31') ENGINE = InnoDB,
|  PARTITION p03 VALUES LESS THAN ('1987-12-31') ENGINE = InnoDB,
|  PARTITION p04 VALUES LESS THAN ('1988-12-31') ENGINE = InnoDB,
|  PARTITION p05 VALUES LESS THAN ('1989-12-31') ENGINE = InnoDB,
|  PARTITION p06 VALUES LESS THAN ('1990-12-31') ENGINE = InnoDB,
|  PARTITION p07 VALUES LESS THAN ('1991-12-31') ENGINE = InnoDB,
|  PARTITION p08 VALUES LESS THAN ('1992-12-31') ENGINE = InnoDB,
|  PARTITION p09 VALUES LESS THAN ('1993-12-31') ENGINE = InnoDB,
|  PARTITION p10 VALUES LESS THAN ('1994-12-31') ENGINE = InnoDB,
|  PARTITION p11 VALUES LESS THAN ('1995-12-31') ENGINE = InnoDB,
|  PARTITION p12 VALUES LESS THAN ('1996-12-31') ENGINE = InnoDB,
|  PARTITION p13 VALUES LESS THAN ('1997-12-31') ENGINE = InnoDB,
|  PARTITION p14 VALUES LESS THAN ('1998-12-31') ENGINE = InnoDB,
|  PARTITION p15 VALUES LESS THAN ('1999-12-31') ENGINE = InnoDB,
|  PARTITION p16 VALUES LESS THAN ('2000-12-31') ENGINE = InnoDB,
|  PARTITION p17 VALUES LESS THAN ('2001-12-31') ENGINE = InnoDB,
|  PARTITION p18 VALUES LESS THAN ('2002-12-31') ENGINE = InnoDB,
|  PARTITION p19 VALUES LESS THAN (MAXVALUE) ENGINE = InnoDB) */
> query
| SHOW CREATE TABLE `sdl_rw_employees_part_2`.`titles`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| titles
: bytes
| CREATE TABLE `titles` (
|   `emp_no` int NOT NULL,
|   `title` varchar(50) NOT NULL,
|   `from_date` date NOT NULL,
|   `to_date` date DEFAULT NULL,
|   PRIMARY KEY (`emp_no`,`title`,`from_date`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50500 PARTITION BY RANGE  COLUMNS(from_date)
| (PARTITION p01 VALUES LESS THAN ('1985-12-31') ENGINE = InnoDB,
|  PARTITION p02 VALUES LESS THAN ('1986-12-31') ENGINE = InnoDB,
|  PARTITION p03 VALUES LESS THAN ('1987-12-31') ENGINE = InnoDB,
|  PARTITION p04 VALUES LESS THAN ('1988-12-31') ENGINE = InnoDB,
|  PARTITION p05 VALUES LESS THAN ('1989-12-31') ENGINE = InnoDB,
|  PARTITION p06 VALUES LESS THAN ('1990-12-31') ENGINE = InnoDB,
|  PARTITION p07 VALUES LESS THAN ('1991-12-31') ENGINE = InnoDB,
|  PARTITION p08 VALUES LESS THAN ('1992-12-31') ENGINE = InnoDB,
|  PARTITION p09 VALUES LESS THAN ('1993-12-31') ENGINE = InnoDB,
|  PARTITION p10 VALUES LESS THAN ('1994-12-31') ENGINE = InnoDB,
|  PARTITION p11 VALUES LESS THAN ('1995-12-31') ENGINE = InnoDB,
|  PARTITION p12 VALUES LESS THAN ('1996-12-31') ENGINE = InnoDB,
|  PARTITION p13 VALUES LESS THAN ('1997-12-31') ENGINE = InnoDB,
|  PARTITION p14 VALUES LESS THAN ('1998-12-31') ENGINE = InnoDB,
|  PARTITION p15 VALUES LESS THAN ('1999-12-31') ENGINE = InnoDB,
|  PARTITION p16 VALUES LESS THAN ('2000-12-31') ENGINE = InnoDB,
|  PARTITION p17 VALUES LESS THAN ('2001-12-31') ENGINE = InnoDB,
|  PARTITION p18 VALUES LESS THAN ('2002-12-31') ENGINE = InnoDB,
|  PARTITION p19 VALUES LESS THAN (MAXVALUE) ENGINE = InnoDB) */
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
| sdl_rw_employees_part_2
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
| sdl_rw_employees_part_2
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
| current_dept_emp
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| departments
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
| 16384
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
| dept_emp
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
| 16384
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
| dept_emp_latest_date
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| dept_manager
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
| 16384
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
| employees
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
| 0
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
| salaries
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 311296
: bytes
| 0
: bytes
| 0
: bytes
| partitioned
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| titles
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 311296
: bytes
| 0
: bytes
| 0
: bytes
| partitioned
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
| sdl_rw_employees_part_2
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
| sdl_rw_employees_part_2
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci

== TestSDLRealWorldIdempotence default#6
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
| sdl_rw_employees_part_2
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
| departments
: bytes
| dept_name
: bytes
| dept_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| departments
: bytes
| PRIMARY
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dept_emp
: bytes
| dept_no
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| dept_emp
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dept_emp
: bytes
| PRIMARY
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dept_manager
: bytes
| dept_no
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| dept_manager
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dept_manager
: bytes
| PRIMARY
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| employees
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| salaries
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| salaries
: bytes
| PRIMARY
: bytes
| from_date
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| titles
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| titles
: bytes
| PRIMARY
: bytes
| title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| titles
: bytes
| PRIMARY
: bytes
| from_date
: bytes
| A
: int64 -1
: null
: int64 3
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
| sdl_rw_employees_part_2
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
| current_dept_emp
: bytes
| emp_no
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
|
: null
< row
: bytes
| current_dept_emp
: bytes
| dept_no
: int64 2
: null
: bytes
| NO
: bytes
| char(4)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| current_dept_emp
: bytes
| from_date
: int64 3
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| current_dept_emp
: bytes
| to_date
: int64 4
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| departments
: bytes
| dept_no
: int64 1
: null
: bytes
| NO
: bytes
| char(4)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| departments
: bytes
| dept_name
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
| dept_emp
: bytes
| emp_no
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
|
: null
< row
: bytes
| dept_emp
: bytes
| dept_no
: int64 2
: null
: bytes
| NO
: bytes
| char(4)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp
: bytes
| from_date
: int64 3
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp
: bytes
| to_date
: int64 4
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp_latest_date
: bytes
| emp_no
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
|
: null
< row
: bytes
| dept_emp_latest_date
: bytes
| from_date
: int64 2
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp_latest_date
: bytes
| to_date
: int64 3
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_manager
: bytes
| emp_no
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
|
: null
< row
: bytes
| dept_manager
: bytes
| dept_no
: int64 2
: null
: bytes
| NO
: bytes
| char(4)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_manager
: bytes
| from_date
: int64 3
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_manager
: bytes
| to_date
: int64 4
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| emp_no
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
|
: null
< row
: bytes
| employees
: bytes
| birth_date
: int64 2
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| first_name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(14)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| last_name
: int64 4
: null
: bytes
| NO
: bytes
| varchar(16)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| gender
: int64 5
: null
: bytes
| NO
: bytes
| enum('M','F')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| hire_date
: int64 6
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| salaries
: bytes
| emp_no
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
|
: null
< row
: bytes
| salaries
: bytes
| salary
: int64 2
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
|
: null
< row
: bytes
| salaries
: bytes
| from_date
: int64 3
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| salaries
: bytes
| to_date
: int64 4
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| titles
: bytes
| emp_no
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
|
: null
< row
: bytes
| titles
: bytes
| title
: int64 2
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| titles
: bytes
| from_date
: int64 3
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| titles
: bytes
| to_date
: int64 4
: null
: bytes
| YES
: bytes
| date
: bytes
|
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
| sdl_rw_employees_part_2
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
| sdl_rw_employees_part_2
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| current_dept_emp
: bytes
| select `sdl_rw_employees_part_2`.`l`.`emp_no` AS `emp_no`,`d`.`dept_no` AS `dept_no`,`sdl_rw_employees_part_2`.`l`.`from_date` AS `from_date`,`sdl_rw_employees_part_2`.`l`.`to_date` AS `to_date` from (`sdl_rw_employees_part_2`.`dept_emp` `d` join `sdl_rw_employees_part_2`.`dept_emp_latest_date` `l` on(((`d`.`emp_no` = `sdl_rw_employees_part_2`.`l`.`emp_no`) and (`d`.`from_date` = `sdl_rw_employees_part_2`.`l`.`from_date`) and (`sdl_rw_employees_part_2`.`l`.`to_date` = `d`.`to_date`))))
< row
: bytes
| dept_emp_latest_date
: bytes
| select `sdl_rw_employees_part_2`.`dept_emp`.`emp_no` AS `emp_no`,max(`sdl_rw_employees_part_2`.`dept_emp`.`from_date`) AS `from_date`,max(`sdl_rw_employees_part_2`.`dept_emp`.`to_date`) AS `to_date` from `sdl_rw_employees_part_2`.`dept_emp` group by `sdl_rw_employees_part_2`.`dept_emp`.`emp_no`
> query
| SHOW CREATE VIEW `sdl_rw_employees_part_2`.`current_dept_emp`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| current_dept_emp
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `current_dept_emp` AS select `l`.`emp_no` AS `emp_no`,`d`.`dept_no` AS `dept_no`,`l`.`from_date` AS `from_date`,`l`.`to_date` AS `to_date` from (`dept_emp` `d` join `dept_emp_latest_date` `l` on(((`d`.`emp_no` = `l`.`emp_no`) and (`d`.`from_date` = `l`.`from_date`) and (`l`.`to_date` = `d`.`to_date`))))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `sdl_rw_employees_part_2`.`dept_emp_latest_date`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| dept_emp_latest_date
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `dept_emp_latest_date` AS select `dept_emp`.`emp_no` AS `emp_no`,max(`dept_emp`.`from_date`) AS `from_date`,max(`dept_emp`.`to_date`) AS `to_date` from `dept_emp` group by `dept_emp`.`emp_no`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| sdl_rw_employees_part_2
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
| sdl_rw_employees_part_2
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
| sdl_rw_employees_part_2
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
< row
: bytes
| dept_emp
: bytes
| dept_emp_ibfk_1
: bytes
| employees
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| dept_emp
: bytes
| dept_emp_ibfk_2
: bytes
| departments
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| dept_manager
: bytes
| dept_manager_ibfk_1
: bytes
| employees
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| dept_manager
: bytes
| dept_manager_ibfk_2
: bytes
| departments
: bytes
| CASCADE
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
| sdl_rw_employees_part_2
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| sdl_rw_employees_part_2
: bytes
| dept_emp
: bytes
| dept_emp_ibfk_1
: bytes
| emp_no
: bytes
| sdl_rw_employees_part_2
: bytes
| emp_no
< row
: bytes
| sdl_rw_employees_part_2
: bytes
| dept_emp
: bytes
| dept_emp_ibfk_2
: bytes
| dept_no
: bytes
| sdl_rw_employees_part_2
: bytes
| dept_no
< row
: bytes
| sdl_rw_employees_part_2
: bytes
| dept_manager
: bytes
| dept_manager_ibfk_1
: bytes
| emp_no
: bytes
| sdl_rw_employees_part_2
: bytes
| emp_no
< row
: bytes
| sdl_rw_employees_part_2
: bytes
| dept_manager
: bytes
| dept_manager_ibfk_2
: bytes
| dept_no
: bytes
| sdl_rw_employees_part_2
: bytes
| dept_no
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
| sdl_rw_employees_part_2
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
< row
: bytes
| salaries
: bytes
| p01
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1985-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p02
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1986-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p03
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1987-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p04
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1988-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p05
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1989-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p06
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1990-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p07
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1991-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p08
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1992-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p09
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1993-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p10
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1994-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p11
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1995-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p12
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1996-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p13
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1997-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p14
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1998-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p15
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1999-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p16
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '2000-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p17
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '2001-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p18
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '2002-12-31'
: bytes
|
< row
: bytes
| salaries
: bytes
| p19
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| MAXVALUE
: bytes
|
< row
: bytes
| titles
: bytes
| p01
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1985-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p02
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1986-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p03
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1987-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p04
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1988-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p05
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1989-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p06
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1990-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p07
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1991-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p08
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1992-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p09
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1993-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p10
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1994-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p11
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1995-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p12
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1996-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p13
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1997-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p14
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1998-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p15
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '1999-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p16
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '2000-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p17
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '2001-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p18
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| '2002-12-31'
: bytes
|
< row
: bytes
| titles
: bytes
| p19
: null
: bytes
| RANGE COLUMNS
: null
: bytes
| `from_date`
: null
: bytes
| MAXVALUE
: bytes
|
> query
| SHOW CREATE TABLE `sdl_rw_employees_part_2`.`salaries`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| salaries
: bytes
| CREATE TABLE `salaries` (
|   `emp_no` int NOT NULL,
|   `salary` int NOT NULL,
|   `from_date` date NOT NULL,
|   `to_date` date NOT NULL,
|   PRIMARY KEY (`emp_no`,`from_date`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50500 PARTITION BY RANGE  COLUMNS(from_date)
| (PARTITION p01 VALUES LESS THAN ('1985-12-31') ENGINE = InnoDB,
|  PARTITION p02 VALUES LESS THAN ('1986-12-31') ENGINE = InnoDB,
|  PARTITION p03 VALUES LESS THAN ('1987-12-31') ENGINE = InnoDB,
|  PARTITION p04 VALUES LESS THAN ('1988-12-31') ENGINE = InnoDB,
|  PARTITION p05 VALUES LESS THAN ('1989-12-31') ENGINE = InnoDB,
|  PARTITION p06 VALUES LESS THAN ('1990-12-31') ENGINE = InnoDB,
|  PARTITION p07 VALUES LESS THAN ('1991-12-31') ENGINE = InnoDB,
|  PARTITION p08 VALUES LESS THAN ('1992-12-31') ENGINE = InnoDB,
|  PARTITION p09 VALUES LESS THAN ('1993-12-31') ENGINE = InnoDB,
|  PARTITION p10 VALUES LESS THAN ('1994-12-31') ENGINE = InnoDB,
|  PARTITION p11 VALUES LESS THAN ('1995-12-31') ENGINE = InnoDB,
|  PARTITION p12 VALUES LESS THAN ('1996-12-31') ENGINE = InnoDB,
|  PARTITION p13 VALUES LESS THAN ('1997-12-31') ENGINE = InnoDB,
|  PARTITION p14 VALUES LESS THAN ('1998-12-31') ENGINE = InnoDB,
|  PARTITION p15 VALUES LESS THAN ('1999-12-31') ENGINE = InnoDB,
|  PARTITION p16 VALUES LESS THAN ('2000-12-31') ENGINE = InnoDB,
|  PARTITION p17 VALUES LESS THAN ('2001-12-31') ENGINE = InnoDB,
|  PARTITION p18 VALUES LESS THAN ('2002-12-31') ENGINE = InnoDB,
|  PARTITION p19 VALUES LESS THAN (MAXVALUE) ENGINE = InnoDB) */
> query
| SHOW CREATE TABLE `sdl_rw_employees_part_2`.`titles`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| titles
: bytes
| CREATE TABLE `titles` (
|   `emp_no` int NOT NULL,
|   `title` varchar(50) NOT NULL,
|   `from_date` date NOT NULL,
|   `to_date` date DEFAULT NULL,
|   PRIMARY KEY (`emp_no`,`title`,`from_date`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50500 PARTITION BY RANGE  COLUMNS(from_date)
| (PARTITION p01 VALUES LESS THAN ('1985-12-31') ENGINE = InnoDB,
|  PARTITION p02 VALUES LESS THAN ('1986-12-31') ENGINE = InnoDB,
|  PARTITION p03 VALUES LESS THAN ('1987-12-31') ENGINE = InnoDB,
|  PARTITION p04 VALUES LESS THAN ('1988-12-31') ENGINE = InnoDB,
|  PARTITION p05 VALUES LESS THAN ('1989-12-31') ENGINE = InnoDB,
|  PARTITION p06 VALUES LESS THAN ('1990-12-31') ENGINE = InnoDB,
|  PARTITION p07 VALUES LESS THAN ('1991-12-31') ENGINE = InnoDB,
|  PARTITION p08 VALUES LESS THAN ('1992-12-31') ENGINE = InnoDB,
|  PARTITION p09 VALUES LESS THAN ('1993-12-31') ENGINE = InnoDB,
|  PARTITION p10 VALUES LESS THAN ('1994-12-31') ENGINE = InnoDB,
|  PARTITION p11 VALUES LESS THAN ('1995-12-31') ENGINE = InnoDB,
|  PARTITION p12 VALUES LESS THAN ('1996-12-31') ENGINE = InnoDB,
|  PARTITION p13 VALUES LESS THAN ('1997-12-31') ENGINE = InnoDB,
|  PARTITION p14 VALUES LESS THAN ('1998-12-31') ENGINE = InnoDB,
|  PARTITION p15 VALUES LESS THAN ('1999-12-31') ENGINE = InnoDB,
|  PARTITION p16 VALUES LESS THAN ('2000-12-31') ENGINE = InnoDB,
|  PARTITION p17 VALUES LESS THAN ('2001-12-31') ENGINE = InnoDB,
|  PARTITION p18 VALUES LESS THAN ('2002-12-31') ENGINE = InnoDB,
|  PARTITION p19 VALUES LESS THAN (MAXVALUE) ENGINE = InnoDB) */
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
| sdl_rw_employees_part_2
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
| sdl_rw_employees_part_2
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
| current_dept_emp
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| departments
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
| 16384
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
| dept_emp
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
| 16384
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
| dept_emp_latest_date
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| dept_manager
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
| 16384
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
| employees
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
| 0
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
| salaries
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 311296
: bytes
| 0
: bytes
| 0
: bytes
| partitioned
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| titles
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 311296
: bytes
| 0
: bytes
| 0
: bytes
| partitioned
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
| sdl_rw_employees_part_2
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
| sdl_rw_employees_part_2
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci

== TestSDLRealWorldIdempotence default#7
> begin
< ok
> exec
| -- Roundcube Webmail initial database structure
|
|
| SET FOREIGN_KEY_CHECKS=0;
< ok 0 0
> exec
|
|
| -- Table structure for table `session`
|
| CREATE TABLE `session` (
|  `sess_id` varchar(128) NOT NULL,
|  `expires_at` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
|  `ip` varchar(40) NOT NULL,
|  `vars` mediumtext NOT NULL,
|  PRIMARY KEY(`sess_id`),
|  INDEX `expires_at_index` (`expires_at`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `users`
|
| CREATE TABLE `users` (
|  `user_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
|  `username` varchar(128) BINARY NOT NULL,
|  `mail_host` varchar(128) NOT NULL,
|  `created` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
|  `last_login` datetime,
|  `failed_login` datetime,
|  `failed_login_counter` int(10) UNSIGNED,
|  `language` varchar(16),
|  `preferences` longtext,
|  PRIMARY KEY(`user_id`),
|  UNIQUE `username` (`username`, `mail_host`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `cache`
|
| CREATE TABLE `cache` (
|  `user_id` int(10) UNSIGNED NOT NULL,
|  `cache_key` varchar(128) BINARY NOT NULL,
|  `expires` datetime,
|  `data` longtext NOT NULL,
|  PRIMARY KEY (`user_id`, `cache_key`),
|  CONSTRAINT `user_id_fk_cache` FOREIGN KEY (`user_id`)
|    REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|  INDEX `expires_index` (`expires`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `cache_shared`
|
| CREATE TABLE `cache_shared` (
|  `cache_key` varchar(255) BINARY NOT NULL,
|  `expires` datetime,
|  `data` longtext NOT NULL,
|  PRIMARY KEY (`cache_key`),
|  INDEX `expires_index` (`expires`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `cache_index`
|
| CREATE TABLE `cache_index` (
|  `user_id` int(10) UNSIGNED NOT NULL,
|  `mailbox` varchar(255) BINARY NOT NULL,
|  `expires` datetime,
|  `valid` tinyint(1) NOT NULL DEFAULT '0',
|  `data` longtext NOT NULL,
|  CONSTRAINT `user_id_fk_cache_index` FOREIGN KEY (`user_id`)
|    REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|  INDEX `expires_index` (`expires`),
|  PRIMARY KEY (`user_id`, `mailbox`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `cache_thread`
|
| CREATE TABLE `cache_thread` (
|  `user_id` int(10) UNSIGNED NOT NULL,
|  `mailbox` varchar(255) BINARY NOT NULL,
|  `expires` datetime,
|  `data` longtext NOT NULL,
|  CONSTRAINT `user_id_fk_cache_thread` FOREIGN KEY (`user_id`)
|    REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|  INDEX `expires_index` (`expires`),
|  PRIMARY KEY (`user_id`, `mailbox`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `cache_messages`
|
| CREATE TABLE `cache_messages` (
|  `user_id` int(10) UNSIGNED NOT NULL,
|  `mailbox` varchar(255) BINARY NOT NULL,
|  `uid` int(11) UNSIGNED NOT NULL DEFAULT '0',
|  `expires` datetime,
|  `data` longtext NOT NULL,
|  `flags` int(11) NOT NULL DEFAULT '0',
|  CONSTRAINT `user_id_fk_cache_messages` FOREIGN KEY (`user_id`)
|    REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|  INDEX `expires_index` (`expires`),
|  PRIMARY KEY (`user_id`, `mailbox`, `uid`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `collected_addresses`
|
| CREATE TABLE `collected_addresses` (
|  `address_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
|  `changed` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
|  `name` varchar(255) NOT NULL DEFAULT '',
|  `email` varchar(255) NOT NULL,
|  `user_id` int(10) UNSIGNED NOT NULL,
|  `type` int(10) UNSIGNED NOT NULL,
|  PRIMARY KEY(`address_id`),
|  CONSTRAINT `user_id_fk_collected_addresses` FOREIGN KEY (`user_id`)
|    REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|  UNIQUE INDEX `user_email_collected_addresses_index` (`user_id`, `type`, `email`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `contacts`
|
| CREATE TABLE `contacts` (
|  `contact_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
|  `changed` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
|  `del` tinyint(1) NOT NULL DEFAULT '0',
|  `name` varchar(128) NOT NULL DEFAULT '',
|  `email` text NOT NULL,
|  `firstname` varchar(128) NOT NULL DEFAULT '',
|  `surname` varchar(128) NOT NULL DEFAULT '',
|  `vcard` longtext,
|  `words` text,
|  `user_id` int(10) UNSIGNED NOT NULL,
|  PRIMARY KEY(`contact_id`),
|  CONSTRAINT `user_id_fk_contacts` FOREIGN KEY (`user_id`)
|    REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|  INDEX `user_contacts_index` (`user_id`,`del`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `contactgroups`
|
| CREATE TABLE `contactgroups` (
|   `contactgroup_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
|   `user_id` int(10) UNSIGNED NOT NULL,
|   `changed` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
|   `del` tinyint(1) NOT NULL DEFAULT '0',
|   `name` varchar(128) NOT NULL DEFAULT '',
|   PRIMARY KEY(`contactgroup_id`),
|   CONSTRAINT `user_id_fk_contactgroups` FOREIGN KEY (`user_id`)
|     REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|   INDEX `contactgroups_user_index` (`user_id`,`del`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `contactgroupmembers`
|
| CREATE TABLE `contactgroupmembers` (
|   `contactgroup_id` int(10) UNSIGNED NOT NULL,
|   `contact_id` int(10) UNSIGNED NOT NULL,
|   `created` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
|   PRIMARY KEY (`contactgroup_id`, `contact_id`),
|   CONSTRAINT `contactgroup_id_fk_contactgroups` FOREIGN KEY (`contactgroup_id`)
|     REFERENCES `contactgroups`(`contactgroup_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|   CONSTRAINT `contact_id_fk_contacts` FOREIGN KEY (`contact_id`)
|     REFERENCES `contacts`(`contact_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|   INDEX `contactgroupmembers_contact_index` (`contact_id`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `identities`
|
| CREATE TABLE `identities` (
|  `identity_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
|  `user_id` int(10) UNSIGNED NOT NULL,
|  `changed` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
|  `del` tinyint(1) NOT NULL DEFAULT '0',
|  `standard` tinyint(1) NOT NULL DEFAULT '0',
|  `name` varchar(128) NOT NULL,
|  `organization` varchar(128) NOT NULL DEFAULT '',
|  `email` varchar(128) NOT NULL,
|  `reply-to` varchar(128) NOT NULL DEFAULT '',
|  `bcc` varchar(128) NOT NULL DEFAULT '',
|  `signature` longtext,
|  `html_signature` tinyint(1) NOT NULL DEFAULT '0',
|  PRIMARY KEY(`identity_id`),
|  CONSTRAINT `user_id_fk_identities` FOREIGN KEY (`user_id`)
|    REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|  INDEX `user_identities_index` (`user_id`, `del`),
|  INDEX `email_identities_index` (`email`, `del`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `responses`
|
| CREATE TABLE `responses` (
|  `response_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
|  `user_id` int(10) UNSIGNED NOT NULL,
|  `name` varchar(255) NOT NULL,
|  `data` longtext NOT NULL,
|  `is_html` tinyint(1) NOT NULL DEFAULT '0',
|  `changed` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
|  `del` tinyint(1) NOT NULL DEFAULT '0',
|  PRIMARY KEY (`response_id`),
|  CONSTRAINT `user_id_fk_responses` FOREIGN KEY (`user_id`)
|    REFERENCES `users` (`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|  INDEX `user_responses_index` (`user_id`, `del`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `dictionary`
|
| CREATE TABLE `dictionary` (
|   `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY, -- redundant, for compat. with Galera Cluster
|   `user_id` int(10) UNSIGNED, -- NULL here is for "shared dictionaries"
|   `language` varchar(16) NOT NULL,
|   `data` longtext NOT NULL,
|   CONSTRAINT `user_id_fk_dictionary` FOREIGN KEY (`user_id`)
|     REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|   UNIQUE `uniqueness` (`user_id`, `language`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `searches`
|
| CREATE TABLE `searches` (
|  `search_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
|  `user_id` int(10) UNSIGNED NOT NULL,
|  `type` int(3) NOT NULL DEFAULT '0',
|  `name` varchar(128) NOT NULL,
|  `data` text,
|  PRIMARY KEY(`search_id`),
|  CONSTRAINT `user_id_fk_searches` FOREIGN KEY (`user_id`)
|    REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|  UNIQUE `uniqueness` (`user_id`, `type`, `name`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
| -- Table structure for table `filestore`
|
| CREATE TABLE `filestore` (
|  `file_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
|  `user_id` int(10) UNSIGNED NOT NULL,
|  `context` varchar(32) NOT NULL,
|  `filename` varchar(128) NOT NULL,
|  `mtime` int(10) NOT NULL,
|  `data` longtext NOT NULL,
|  PRIMARY KEY (`file_id`),
|  CONSTRAINT `user_id_fk_filestore` FOREIGN KEY (`user_id`)
|    REFERENCES `users` (`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|  UNIQUE `uniqueness` (`user_id`, `context`, `filename`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
| -- Table structure for table `uploads`
|
| CREATE TABLE `uploads` (
|  `upload_id` varchar(64) NOT NULL,
|  `session_id` varchar(128) NOT NULL,
|  `group` varchar(128) NOT NULL,
|  `metadata` mediumtext NOT NULL,
|  `created` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
|  PRIMARY KEY (`upload_id`),
|  INDEX `uploads_session_group_index` (`session_id`, `group`, `created`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
| -- Table structure for table `system`
|
| CREATE TABLE `system` (
|  `name` varchar(64) NOT NULL,
|  `value` mediumtext,
|  PRIMARY KEY(`name`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
| SET FOREIGN_KEY_CHECKS=1;
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
| sdl_rw_roundcube_3
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
| cache
: bytes
| expires_index
: bytes
| expires
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| cache
: bytes
| PRIMARY
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache
: bytes
| PRIMARY
: bytes
| cache_key
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_index
: bytes
| expires_index
: bytes
| expires
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| cache_index
: bytes
| PRIMARY
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_index
: bytes
| PRIMARY
: bytes
| mailbox
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_messages
: bytes
| expires_index
: bytes
| expires
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| cache_messages
: bytes
| PRIMARY
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_messages
: bytes
| PRIMARY
: bytes
| mailbox
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_messages
: bytes
| PRIMARY
: bytes
| uid
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_shared
: bytes
| expires_index
: bytes
| expires
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| cache_shared
: bytes
| PRIMARY
: bytes
| cache_key
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_thread
: bytes
| expires_index
: bytes
| expires
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| cache_thread
: bytes
| PRIMARY
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_thread
: bytes
| PRIMARY
: bytes
| mailbox
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| collected_addresses
: bytes
| PRIMARY
: bytes
| address_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| collected_addresses
: bytes
| user_email_collected_addresses_index
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| collected_addresses
: bytes
| user_email_collected_addresses_index
: bytes
| type
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| collected_addresses
: bytes
| user_email_collected_addresses_index
: bytes
| email
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| contactgroupmembers
: bytes
| contactgroupmembers_contact_index
: bytes
| contact_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| contactgroupmembers
: bytes
| PRIMARY
: bytes
| contactgroup_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| contactgroupmembers
: bytes
| PRIMARY
: bytes
| contact_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| contactgroups
: bytes
| contactgroups_user_index
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| contactgroups
: bytes
| contactgroups_user_index
: bytes
| del
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| contactgroups
: bytes
| PRIMARY
: bytes
| contactgroup_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| contacts
: bytes
| PRIMARY
: bytes
| contact_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| contacts
: bytes
| user_contacts_index
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| contacts
: bytes
| user_contacts_index
: bytes
| del
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| dictionary
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
| dictionary
: bytes
| uniqueness
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dictionary
: bytes
| uniqueness
: bytes
| language
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| filestore
: bytes
| PRIMARY
: bytes
| file_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| filestore
: bytes
| uniqueness
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| filestore
: bytes
| uniqueness
: bytes
| context
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| filestore
: bytes
| uniqueness
: bytes
| filename
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| identities
: bytes
| email_identities_index
: bytes
| email
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| identities
: bytes
| email_identities_index
: bytes
| del
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| identities
: bytes
| PRIMARY
: bytes
| identity_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| identities
: bytes
| user_identities_index
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| identities
: bytes
| user_identities_index
: bytes
| del
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| responses
: bytes
| PRIMARY
: bytes
| response_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| responses
: bytes
| user_responses_index
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| responses
: bytes
| user_responses_index
: bytes
| del
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| searches
: bytes
| PRIMARY
: bytes
| search_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| searches
: bytes
| uniqueness
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| searches
: bytes
| uniqueness
: bytes
| type
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| searches
: bytes
| uniqueness
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| session
: bytes
| expires_at_index
: bytes
| expires_at
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| session
: bytes
| PRIMARY
: bytes
| sess_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| system
: bytes
| PRIMARY
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| uploads
: bytes
| PRIMARY
: bytes
| upload_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| uploads
: bytes
| uploads_session_group_index
: bytes
| session_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| uploads
: bytes
| uploads_session_group_index
: bytes
| group
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| uploads
: bytes
| uploads_session_group_index
: bytes
| created
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| users
: bytes
| PRIMARY
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| users
: bytes
| username
: bytes
| username
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| users
: bytes
| username
: bytes
| mail_host
: bytes
| A
: int64 -1
: null
: int64 2
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
| sdl_rw_roundcube_3
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
| cache
: bytes
| user_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache
: bytes
| cache_key
: int64 2
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache
: bytes
| expires
: int64 3
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache
: bytes
| data
: int64 4
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_index
: bytes
| user_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_index
: bytes
| mailbox
: int64 2
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_index
: bytes
| expires
: int64 3
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_index
: bytes
| valid
: int64 4
: bytes
| '0'
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
| cache_index
: bytes
| data
: int64 5
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_messages
: bytes
| user_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_messages
: bytes
| mailbox
: int64 2
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_messages
: bytes
| uid
: int64 3
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_messages
: bytes
| expires
: int64 4
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_messages
: bytes
| data
: int64 5
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_messages
: bytes
| flags
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
| cache_shared
: bytes
| cache_key
: int64 1
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_shared
: bytes
| expires
: int64 2
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_shared
: bytes
| data
: int64 3
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_thread
: bytes
| user_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_thread
: bytes
| mailbox
: int64 2
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_thread
: bytes
| expires
: int64 3
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_thread
: bytes
| data
: int64 4
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| collected_addresses
: bytes
| address_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| collected_addresses
: bytes
| changed
: int64 2
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| collected_addresses
: bytes
| name
: int64 3
: bytes
| ''
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| collected_addresses
: bytes
| email
: int64 4
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| collected_addresses
: bytes
| user_id
: int64 5
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| collected_addresses
: bytes
| type
: int64 6
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contactgroupmembers
: bytes
| contactgroup_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contactgroupmembers
: bytes
| contact_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contactgroupmembers
: bytes
| created
: int64 3
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contactgroups
: bytes
| contactgroup_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| contactgroups
: bytes
| user_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contactgroups
: bytes
| changed
: int64 3
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contactgroups
: bytes
| del
: int64 4
: bytes
| '0'
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
| contactgroups
: bytes
| name
: int64 5
: bytes
| ''
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| contact_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| contacts
: bytes
| changed
: int64 2
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| del
: int64 3
: bytes
| '0'
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
| contacts
: bytes
| name
: int64 4
: bytes
| ''
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| email
: int64 5
: null
: bytes
| NO
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| firstname
: int64 6
: bytes
| ''
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| surname
: int64 7
: bytes
| ''
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| vcard
: int64 8
: null
: bytes
| YES
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| words
: int64 9
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| user_id
: int64 10
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dictionary
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| dictionary
: bytes
| user_id
: int64 2
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dictionary
: bytes
| language
: int64 3
: null
: bytes
| NO
: bytes
| varchar(16)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dictionary
: bytes
| data
: int64 4
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filestore
: bytes
| file_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| filestore
: bytes
| user_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filestore
: bytes
| context
: int64 3
: null
: bytes
| NO
: bytes
| varchar(32)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filestore
: bytes
| filename
: int64 4
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filestore
: bytes
| mtime
: int64 5
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
|
: null
< row
: bytes
| filestore
: bytes
| data
: int64 6
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| identity_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| identities
: bytes
| user_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| changed
: int64 3
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| del
: int64 4
: bytes
| '0'
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
| identities
: bytes
| standard
: int64 5
: bytes
| '0'
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
| identities
: bytes
| name
: int64 6
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| organization
: int64 7
: bytes
| ''
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| email
: int64 8
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| reply-to
: int64 9
: bytes
| ''
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| bcc
: int64 10
: bytes
| ''
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| signature
: int64 11
: null
: bytes
| YES
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| html_signature
: int64 12
: bytes
| '0'
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
| responses
: bytes
| response_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| responses
: bytes
| user_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| responses
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| responses
: bytes
| data
: int64 4
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| responses
: bytes
| is_html
: int64 5
: bytes
| '0'
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
| responses
: bytes
| changed
: int64 6
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| responses
: bytes
| del
: int64 7
: bytes
| '0'
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
| searches
: bytes
| search_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| searches
: bytes
| user_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| searches
: bytes
| type
: int64 3
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
| searches
: bytes
| name
: int64 4
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| searches
: bytes
| data
: int64 5
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| session
: bytes
| sess_id
: int64 1
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| session
: bytes
| expires_at
: int64 2
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| session
: bytes
| ip
: int64 3
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| session
: bytes
| vars
: int64 4
: null
: bytes
| NO
: bytes
| mediumtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| system
: bytes
| name
: int64 1
: null
: bytes
| NO
: bytes
| varchar(64)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| system
: bytes
| value
: int64 2
: null
: bytes
| YES
: bytes
| mediumtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploads
: bytes
| upload_id
: int64 1
: null
: bytes
| NO
: bytes
| varchar(64)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploads
: bytes
| session_id
: int64 2
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploads
: bytes
| group
: int64 3
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploads
: bytes
| metadata
: int64 4
: null
: bytes
| NO
: bytes
| mediumtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploads
: bytes
| created
: int64 5
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| user_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| users
: bytes
| username
: int64 2
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| mail_host
: int64 3
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| created
: int64 4
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| last_login
: int64 5
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| failed_login
: int64 6
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| failed_login_counter
: int64 7
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| language
: int64 8
: null
: bytes
| YES
: bytes
| varchar(16)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| preferences
: int64 9
: null
: bytes
| YES
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
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
| sdl_rw_roundcube_3
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
| sdl_rw_roundcube_3
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
| sdl_rw_roundcube_3
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
| sdl_rw_roundcube_3
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
| sdl_rw_roundcube_3
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
< row
: bytes
| contactgroupmembers
: bytes
| contactgroup_id_fk_contactgroups
: bytes
| contactgroups
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| contactgroupmembers
: bytes
| contact_id_fk_contacts
: bytes
| contacts
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| cache
: bytes
| user_id_fk_cache
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| cache_index
: bytes
| user_id_fk_cache_index
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| cache_messages
: bytes
| user_id_fk_cache_messages
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| cache_thread
: bytes
| user_id_fk_cache_thread
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| collected_addresses
: bytes
| user_id_fk_collected_addresses
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| contactgroups
: bytes
| user_id_fk_contactgroups
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| contacts
: bytes
| user_id_fk_contacts
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| dictionary
: bytes
| user_id_fk_dictionary
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| filestore
: bytes
| user_id_fk_filestore
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| identities
: bytes
| user_id_fk_identities
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| responses
: bytes
| user_id_fk_responses
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| searches
: bytes
| user_id_fk_searches
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
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
| sdl_rw_roundcube_3
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| cache
: bytes
| user_id_fk_cache
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| cache_index
: bytes
| user_id_fk_cache_index
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| cache_messages
: bytes
| user_id_fk_cache_messages
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| cache_thread
: bytes
| user_id_fk_cache_thread
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| collected_addresses
: bytes
| user_id_fk_collected_addresses
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| contactgroupmembers
: bytes
| contact_id_fk_contacts
: bytes
| contact_id
: bytes
| sdl_rw_roundcube_3
: bytes
| contact_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| contactgroupmembers
: bytes
| contactgroup_id_fk_contactgroups
: bytes
| contactgroup_id
: bytes
| sdl_rw_roundcube_3
: bytes
| contactgroup_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| contactgroups
: bytes
| user_id_fk_contactgroups
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| contacts
: bytes
| user_id_fk_contacts
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| dictionary
: bytes
| user_id_fk_dictionary
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| filestore
: bytes
| user_id_fk_filestore
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| identities
: bytes
| user_id_fk_identities
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| responses
: bytes
| user_id_fk_responses
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| searches
: bytes
| user_id_fk_searches
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
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
| sdl_rw_roundcube_3
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
| sdl_rw_roundcube_3
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
| sdl_rw_roundcube_3
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
| cache
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| cache_index
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| cache_messages
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| cache_shared
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| cache_thread
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| collected_addresses
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| contactgroupmembers
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
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| contactgroups
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| contacts
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| dictionary
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| filestore
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| identities
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 32768
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| responses
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| searches
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| session
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| system
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| uploads
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| users
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
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
| sdl_rw_roundcube_3
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
| sdl_rw_roundcube_3
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci

== TestSDLRealWorldIdempotence default#8
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
| sdl_rw_roundcube_3
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
| cache
: bytes
| expires_index
: bytes
| expires
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| cache
: bytes
| PRIMARY
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache
: bytes
| PRIMARY
: bytes
| cache_key
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_index
: bytes
| expires_index
: bytes
| expires
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| cache_index
: bytes
| PRIMARY
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_index
: bytes
| PRIMARY
: bytes
| mailbox
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_messages
: bytes
| expires_index
: bytes
| expires
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| cache_messages
: bytes
| PRIMARY
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_messages
: bytes
| PRIMARY
: bytes
| mailbox
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_messages
: bytes
| PRIMARY
: bytes
| uid
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_shared
: bytes
| expires_index
: bytes
| expires
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| cache_shared
: bytes
| PRIMARY
: bytes
| cache_key
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_thread
: bytes
| expires_index
: bytes
| expires
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| cache_thread
: bytes
| PRIMARY
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_thread
: bytes
| PRIMARY
: bytes
| mailbox
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| collected_addresses
: bytes
| PRIMARY
: bytes
| address_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| collected_addresses
: bytes
| user_email_collected_addresses_index
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| collected_addresses
: bytes
| user_email_collected_addresses_index
: bytes
| type
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| collected_addresses
: bytes
| user_email_collected_addresses_index
: bytes
| email
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| contactgroupmembers
: bytes
| contactgroupmembers_contact_index
: bytes
| contact_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| contactgroupmembers
: bytes
| PRIMARY
: bytes
| contactgroup_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| contactgroupmembers
: bytes
| PRIMARY
: bytes
| contact_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| contactgroups
: bytes
| contactgroups_user_index
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| contactgroups
: bytes
| contactgroups_user_index
: bytes
| del
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| contactgroups
: bytes
| PRIMARY
: bytes
| contactgroup_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| contacts
: bytes
| PRIMARY
: bytes
| contact_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| contacts
: bytes
| user_contacts_index
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| contacts
: bytes
| user_contacts_index
: bytes
| del
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| dictionary
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
| dictionary
: bytes
| uniqueness
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dictionary
: bytes
| uniqueness
: bytes
| language
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| filestore
: bytes
| PRIMARY
: bytes
| file_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| filestore
: bytes
| uniqueness
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| filestore
: bytes
| uniqueness
: bytes
| context
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| filestore
: bytes
| uniqueness
: bytes
| filename
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| identities
: bytes
| email_identities_index
: bytes
| email
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| identities
: bytes
| email_identities_index
: bytes
| del
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| identities
: bytes
| PRIMARY
: bytes
| identity_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| identities
: bytes
| user_identities_index
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| identities
: bytes
| user_identities_index
: bytes
| del
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| responses
: bytes
| PRIMARY
: bytes
| response_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| responses
: bytes
| user_responses_index
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| responses
: bytes
| user_responses_index
: bytes
| del
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| searches
: bytes
| PRIMARY
: bytes
| search_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| searches
: bytes
| uniqueness
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| searches
: bytes
| uniqueness
: bytes
| type
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| searches
: bytes
| uniqueness
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| session
: bytes
| expires_at_index
: bytes
| expires_at
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| session
: bytes
| PRIMARY
: bytes
| sess_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| system
: bytes
| PRIMARY
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| uploads
: bytes
| PRIMARY
: bytes
| upload_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| uploads
: bytes
| uploads_session_group_index
: bytes
| session_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| uploads
: bytes
| uploads_session_group_index
: bytes
| group
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| uploads
: bytes
| uploads_session_group_index
: bytes
| created
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| users
: bytes
| PRIMARY
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| users
: bytes
| username
: bytes
| username
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| users
: bytes
| username
: bytes
| mail_host
: bytes
| A
: int64 -1
: null
: int64 2
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
| sdl_rw_roundcube_3
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
| cache
: bytes
| user_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache
: bytes
| cache_key
: int64 2
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache
: bytes
| expires
: int64 3
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache
: bytes
| data
: int64 4
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_index
: bytes
| user_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_index
: bytes
| mailbox
: int64 2
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_index
: bytes
| expires
: int64 3
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_index
: bytes
| valid
: int64 4
: bytes
| '0'
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
| cache_index
: bytes
| data
: int64 5
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_messages
: bytes
| user_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_messages
: bytes
| mailbox
: int64 2
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_messages
: bytes
| uid
: int64 3
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_messages
: bytes
| expires
: int64 4
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_messages
: bytes
| data
: int64 5
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_messages
: bytes
| flags
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
| cache_shared
: bytes
| cache_key
: int64 1
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_shared
: bytes
| expires
: int64 2
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_shared
: bytes
| data
: int64 3
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_thread
: bytes
| user_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_thread
: bytes
| mailbox
: int64 2
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_thread
: bytes
| expires
: int64 3
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_thread
: bytes
| data
: int64 4
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| collected_addresses
: bytes
| address_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| collected_addresses
: bytes
| changed
: int64 2
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| collected_addresses
: bytes
| name
: int64 3
: bytes
| ''
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| collected_addresses
: bytes
| email
: int64 4
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| collected_addresses
: bytes
| user_id
: int64 5
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| collected_addresses
: bytes
| type
: int64 6
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contactgroupmembers
: bytes
| contactgroup_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contactgroupmembers
: bytes
| contact_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contactgroupmembers
: bytes
| created
: int64 3
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contactgroups
: bytes
| contactgroup_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| contactgroups
: bytes
| user_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contactgroups
: bytes
| changed
: int64 3
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contactgroups
: bytes
| del
: int64 4
: bytes
| '0'
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
| contactgroups
: bytes
| name
: int64 5
: bytes
| ''
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| contact_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| contacts
: bytes
| changed
: int64 2
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| del
: int64 3
: bytes
| '0'
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
| contacts
: bytes
| name
: int64 4
: bytes
| ''
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| email
: int64 5
: null
: bytes
| NO
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| firstname
: int64 6
: bytes
| ''
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| surname
: int64 7
: bytes
| ''
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| vcard
: int64 8
: null
: bytes
| YES
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| words
: int64 9
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| user_id
: int64 10
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dictionary
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| dictionary
: bytes
| user_id
: int64 2
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dictionary
: bytes
| language
: int64 3
: null
: bytes
| NO
: bytes
| varchar(16)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dictionary
: bytes
| data
: int64 4
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filestore
: bytes
| file_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| filestore
: bytes
| user_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filestore
: bytes
| context
: int64 3
: null
: bytes
| NO
: bytes
| varchar(32)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filestore
: bytes
| filename
: int64 4
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filestore
: bytes
| mtime
: int64 5
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
|
: null
< row
: bytes
| filestore
: bytes
| data
: int64 6
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| identity_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| identities
: bytes
| user_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| changed
: int64 3
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| del
: int64 4
: bytes
| '0'
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
| identities
: bytes
| standard
: int64 5
: bytes
| '0'
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
| identities
: bytes
| name
: int64 6
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| organization
: int64 7
: bytes
| ''
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| email
: int64 8
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| reply-to
: int64 9
: bytes
| ''
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| bcc
: int64 10
: bytes
| ''
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| signature
: int64 11
: null
: bytes
| YES
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| html_signature
: int64 12
: bytes
| '0'
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
| responses
: bytes
| response_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| responses
: bytes
| user_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| responses
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| responses
: bytes
| data
: int64 4
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| responses
: bytes
| is_html
: int64 5
: bytes
| '0'
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
| responses
: bytes
| changed
: int64 6
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| responses
: bytes
| del
: int64 7
: bytes
| '0'
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
| searches
: bytes
| search_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| searches
: bytes
| user_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| searches
: bytes
| type
: int64 3
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
| searches
: bytes
| name
: int64 4
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| searches
: bytes
| data
: int64 5
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| session
: bytes
| sess_id
: int64 1
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| session
: bytes
| expires_at
: int64 2
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| session
: bytes
| ip
: int64 3
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| session
: bytes
| vars
: int64 4
: null
: bytes
| NO
: bytes
| mediumtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| system
: bytes
| name
: int64 1
: null
: bytes
| NO
: bytes
| varchar(64)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| system
: bytes
| value
: int64 2
: null
: bytes
| YES
: bytes
| mediumtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploads
: bytes
| upload_id
: int64 1
: null
: bytes
| NO
: bytes
| varchar(64)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploads
: bytes
| session_id
: int64 2
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploads
: bytes
| group
: int64 3
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploads
: bytes
| metadata
: int64 4
: null
: bytes
| NO
: bytes
| mediumtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploads
: bytes
| created
: int64 5
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| user_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| users
: bytes
| username
: int64 2
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| mail_host
: int64 3
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| created
: int64 4
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| last_login
: int64 5
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| failed_login
: int64 6
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| failed_login_counter
: int64 7
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| language
: int64 8
: null
: bytes
| YES
: bytes
| varchar(16)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| preferences
: int64 9
: null
: bytes
| YES
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
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
| sdl_rw_roundcube_3
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
| sdl_rw_roundcube_3
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
| sdl_rw_roundcube_3
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
| sdl_rw_roundcube_3
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
| sdl_rw_roundcube_3
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
< row
: bytes
| contactgroupmembers
: bytes
| contactgroup_id_fk_contactgroups
: bytes
| contactgroups
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| contactgroupmembers
: bytes
| contact_id_fk_contacts
: bytes
| contacts
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| cache
: bytes
| user_id_fk_cache
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| cache_index
: bytes
| user_id_fk_cache_index
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| cache_messages
: bytes
| user_id_fk_cache_messages
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| cache_thread
: bytes
| user_id_fk_cache_thread
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| collected_addresses
: bytes
| user_id_fk_collected_addresses
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| contactgroups
: bytes
| user_id_fk_contactgroups
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| contacts
: bytes
| user_id_fk_contacts
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| dictionary
: bytes
| user_id_fk_dictionary
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| filestore
: bytes
| user_id_fk_filestore
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| identities
: bytes
| user_id_fk_identities
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| responses
: bytes
| user_id_fk_responses
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| searches
: bytes
| user_id_fk_searches
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
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
| sdl_rw_roundcube_3
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| cache
: bytes
| user_id_fk_cache
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| cache_index
: bytes
| user_id_fk_cache_index
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| cache_messages
: bytes
| user_id_fk_cache_messages
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| cache_thread
: bytes
| user_id_fk_cache_thread
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| collected_addresses
: bytes
| user_id_fk_collected_addresses
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| contactgroupmembers
: bytes
| contact_id_fk_contacts
: bytes
| contact_id
: bytes
| sdl_rw_roundcube_3
: bytes
| contact_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| contactgroupmembers
: bytes
| contactgroup_id_fk_contactgroups
: bytes
| contactgroup_id
: bytes
| sdl_rw_roundcube_3
: bytes
| contactgroup_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| contactgroups
: bytes
| user_id_fk_contactgroups
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| contacts
: bytes
| user_id_fk_contacts
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| dictionary
: bytes
| user_id_fk_dictionary
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| filestore
: bytes
| user_id_fk_filestore
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| identities
: bytes
| user_id_fk_identities
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| responses
: bytes
| user_id_fk_responses
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_3
: bytes
| searches
: bytes
| user_id_fk_searches
: bytes
| user_id
: bytes
| sdl_rw_roundcube_3
: bytes
| user_id
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
| sdl_rw_roundcube_3
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
| sdl_rw_roundcube_3
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
| sdl_rw_roundcube_3
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
| cache
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| cache_index
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| cache_messages
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| cache_shared
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| cache_thread
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| collected_addresses
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| contactgroupmembers
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
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| contactgroups
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| contacts
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| dictionary
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| filestore
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| identities
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 32768
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| responses
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| searches
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| session
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| system
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| uploads
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| users
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
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
| sdl_rw_roundcube_3
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
| sdl_rw_roundcube_3
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci

== TestSDLRealWorldIdempotence default#9
> begin
< ok
> exec
|# sha256:6b5949d9dbb9b133b2fedc1cb3f6af6ad86cbef74576a097b409fc8c9e7e4d08 75716
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
| sdl_rw_mediawiki_4
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
| actor
: bytes
| actor_name
: bytes
| actor_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| actor
: bytes
| actor_user
: bytes
| actor_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| actor
: bytes
| PRIMARY
: bytes
| actor_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| archive
: bytes
| ar_actor_timestamp
: bytes
| ar_actor
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| archive
: bytes
| ar_actor_timestamp
: bytes
| ar_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| archive
: bytes
| ar_revid
: bytes
| ar_rev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| archive
: bytes
| ar_usertext_timestamp
: bytes
| ar_user_text
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| archive
: bytes
| ar_usertext_timestamp
: bytes
| ar_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| archive
: bytes
| name_title_timestamp
: bytes
| ar_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| archive
: bytes
| name_title_timestamp
: bytes
| ar_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| archive
: bytes
| name_title_timestamp
: bytes
| ar_timestamp
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| archive
: bytes
| PRIMARY
: bytes
| ar_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| bot_passwords
: bytes
| PRIMARY
: bytes
| bp_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| bot_passwords
: bytes
| PRIMARY
: bytes
| bp_app_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| category
: bytes
| cat_pages
: bytes
| cat_pages
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| category
: bytes
| cat_title
: bytes
| cat_title
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| category
: bytes
| PRIMARY
: bytes
| cat_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| cl_collation_ext
: bytes
| cl_collation
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| cl_collation_ext
: bytes
| cl_to
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| cl_collation_ext
: bytes
| cl_type
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| cl_collation_ext
: bytes
| cl_from
: bytes
| A
: int64 -1
: null
: int64 4
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| cl_sortkey
: bytes
| cl_to
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| cl_sortkey
: bytes
| cl_type
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| cl_sortkey
: bytes
| cl_sortkey
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| cl_sortkey
: bytes
| cl_from
: bytes
| A
: int64 -1
: null
: int64 4
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| cl_timestamp
: bytes
| cl_to
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| cl_timestamp
: bytes
| cl_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| PRIMARY
: bytes
| cl_from
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| categorylinks
: bytes
| PRIMARY
: bytes
| cl_to
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| change_tag_log_tag
: bytes
| ct_log_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| change_tag_log_tag
: bytes
| ct_tag
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| change_tag_rc_tag
: bytes
| ct_rc_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| change_tag_rc_tag
: bytes
| ct_tag
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| change_tag_rev_tag
: bytes
| ct_rev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| change_tag_rev_tag
: bytes
| ct_tag
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| change_tag_tag_id
: bytes
| ct_tag
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| change_tag_tag_id
: bytes
| ct_rc_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| change_tag_tag_id
: bytes
| ct_rev_id
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| change_tag_tag_id
: bytes
| ct_log_id
: bytes
| A
: int64 -1
: null
: int64 4
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| change_tag
: bytes
| PRIMARY
: bytes
| ct_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| comment
: bytes
| comment_hash
: bytes
| comment_hash
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| comment
: bytes
| PRIMARY
: bytes
| comment_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| content
: bytes
| PRIMARY
: bytes
| content_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| content_models
: bytes
| model_name
: bytes
| model_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| content_models
: bytes
| PRIMARY
: bytes
| model_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| el_from
: bytes
| el_from
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| el_from
: bytes
| el_to
: bytes
| A
: int64 40
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| el_from_index_60
: bytes
| el_from
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| el_from_index_60
: bytes
| el_index_60
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| el_from_index_60
: bytes
| el_id
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| el_index
: bytes
| el_index
: bytes
| A
: int64 60
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| el_index_60
: bytes
| el_index_60
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| el_index_60
: bytes
| el_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| el_to
: bytes
| el_to
: bytes
| A
: int64 60
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| el_to
: bytes
| el_from
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| externallinks
: bytes
| PRIMARY
: bytes
| el_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| fa_actor_timestamp
: bytes
| fa_actor
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| fa_actor_timestamp
: bytes
| fa_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| fa_deleted_timestamp
: bytes
| fa_deleted_timestamp
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| fa_name
: bytes
| fa_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| fa_name
: bytes
| fa_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| fa_sha1
: bytes
| fa_sha1
: bytes
| A
: int64 10
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| fa_storage_group
: bytes
| fa_storage_group
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| fa_storage_group
: bytes
| fa_storage_key
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| fa_user_timestamp
: bytes
| fa_user_text
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| fa_user_timestamp
: bytes
| fa_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| filearchive
: bytes
| PRIMARY
: bytes
| fa_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_actor_timestamp
: bytes
| img_actor
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_actor_timestamp
: bytes
| img_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_media_mime
: bytes
| img_media_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_media_mime
: bytes
| img_major_mime
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_media_mime
: bytes
| img_minor_mime
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_sha1
: bytes
| img_sha1
: bytes
| A
: int64 10
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_size
: bytes
| img_size
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_timestamp
: bytes
| img_timestamp
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_user_timestamp
: bytes
| img_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_user_timestamp
: bytes
| img_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_usertext_timestamp
: bytes
| img_user_text
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| img_usertext_timestamp
: bytes
| img_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| image
: bytes
| PRIMARY
: bytes
| img_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| image_comment_temp
: bytes
| imgcomment_name
: bytes
| imgcomment_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| image_comment_temp
: bytes
| PRIMARY
: bytes
| imgcomment_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| image_comment_temp
: bytes
| PRIMARY
: bytes
| imgcomment_description_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| imagelinks
: bytes
| il_backlinks_namespace
: bytes
| il_from_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| imagelinks
: bytes
| il_backlinks_namespace
: bytes
| il_to
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| imagelinks
: bytes
| il_backlinks_namespace
: bytes
| il_from
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| imagelinks
: bytes
| il_to
: bytes
| il_to
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| imagelinks
: bytes
| il_to
: bytes
| il_from
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| imagelinks
: bytes
| PRIMARY
: bytes
| il_from
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| imagelinks
: bytes
| PRIMARY
: bytes
| il_to
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| interwiki
: bytes
| iw_prefix
: bytes
| iw_prefix
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| ip_changes
: bytes
| ipc_hex_time
: bytes
| ipc_hex
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| ip_changes
: bytes
| ipc_hex_time
: bytes
| ipc_rev_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| ip_changes
: bytes
| ipc_rev_timestamp
: bytes
| ipc_rev_timestamp
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| ip_changes
: bytes
| PRIMARY
: bytes
| ipc_rev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| ipb_address
: bytes
| ipb_address
: bytes
| A
: int64 255
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| ipb_address
: bytes
| ipb_user
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| ipb_address
: bytes
| ipb_auto
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| ipb_address
: bytes
| ipb_anon_only
: bytes
| A
: int64 -1
: null
: int64 4
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| ipb_expiry
: bytes
| ipb_expiry
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| ipb_parent_block_id
: bytes
| ipb_parent_block_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| ipb_range
: bytes
| ipb_range_start
: bytes
| A
: int64 8
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| ipb_range
: bytes
| ipb_range_end
: bytes
| A
: int64 8
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| ipb_timestamp
: bytes
| ipb_timestamp
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| ipb_user
: bytes
| ipb_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| ipblocks
: bytes
| PRIMARY
: bytes
| ipb_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| iwlinks
: bytes
| iwl_prefix_from_title
: bytes
| iwl_prefix
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| iwlinks
: bytes
| iwl_prefix_from_title
: bytes
| iwl_from
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| iwlinks
: bytes
| iwl_prefix_from_title
: bytes
| iwl_title
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| iwlinks
: bytes
| iwl_prefix_title_from
: bytes
| iwl_prefix
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| iwlinks
: bytes
| iwl_prefix_title_from
: bytes
| iwl_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| iwlinks
: bytes
| iwl_prefix_title_from
: bytes
| iwl_from
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| iwlinks
: bytes
| PRIMARY
: bytes
| iwl_from
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| iwlinks
: bytes
| PRIMARY
: bytes
| iwl_prefix
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| iwlinks
: bytes
| PRIMARY
: bytes
| iwl_title
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_cmd
: bytes
| job_cmd
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_cmd
: bytes
| job_namespace
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_cmd
: bytes
| job_title
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_cmd
: bytes
| job_params
: bytes
| A
: int64 128
: null
: int64 4
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_cmd_token
: bytes
| job_cmd
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_cmd_token
: bytes
| job_token
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_cmd_token
: bytes
| job_random
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_cmd_token_id
: bytes
| job_cmd
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_cmd_token_id
: bytes
| job_token
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_cmd_token_id
: bytes
| job_id
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_sha1
: bytes
| job_sha1
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| job_timestamp
: bytes
| job_timestamp
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| job
: bytes
| PRIMARY
: bytes
| job_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| l10n_cache
: bytes
| PRIMARY
: bytes
| lc_lang
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| l10n_cache
: bytes
| PRIMARY
: bytes
| lc_key
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| langlinks
: bytes
| ll_lang
: bytes
| ll_lang
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| langlinks
: bytes
| ll_lang
: bytes
| ll_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| langlinks
: bytes
| PRIMARY
: bytes
| ll_from
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| langlinks
: bytes
| PRIMARY
: bytes
| ll_lang
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| log_search
: bytes
| ls_log_id
: bytes
| ls_log_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| log_search
: bytes
| PRIMARY
: bytes
| ls_field
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| log_search
: bytes
| PRIMARY
: bytes
| ls_value
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| log_search
: bytes
| PRIMARY
: bytes
| ls_log_id
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| actor_time
: bytes
| log_actor
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| actor_time
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_actor_type_time
: bytes
| log_actor
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_actor_type_time
: bytes
| log_type
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_actor_type_time
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_page_id_time
: bytes
| log_page
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_page_id_time
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_user_text_time
: bytes
| log_user_text
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_user_text_time
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_user_text_type_time
: bytes
| log_user_text
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_user_text_type_time
: bytes
| log_type
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_user_text_type_time
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_user_type_time
: bytes
| log_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_user_type_time
: bytes
| log_type
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| log_user_type_time
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| page_time
: bytes
| log_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| page_time
: bytes
| log_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| page_time
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| PRIMARY
: bytes
| log_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| times
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| type_action
: bytes
| log_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| type_action
: bytes
| log_action
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| type_action
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| type_time
: bytes
| log_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| type_time
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| user_time
: bytes
| log_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| logging
: bytes
| user_time
: bytes
| log_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| module_deps
: bytes
| PRIMARY
: bytes
| md_module
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| module_deps
: bytes
| PRIMARY
: bytes
| md_skin
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| objectcache
: bytes
| exptime
: bytes
| exptime
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| objectcache
: bytes
| PRIMARY
: bytes
| keyname
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| oldimage
: bytes
| oi_actor_timestamp
: bytes
| oi_actor
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| oldimage
: bytes
| oi_actor_timestamp
: bytes
| oi_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| oldimage
: bytes
| oi_name_archive_name
: bytes
| oi_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| oldimage
: bytes
| oi_name_archive_name
: bytes
| oi_archive_name
: bytes
| A
: int64 14
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| oldimage
: bytes
| oi_name_timestamp
: bytes
| oi_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| oldimage
: bytes
| oi_name_timestamp
: bytes
| oi_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| oldimage
: bytes
| oi_sha1
: bytes
| oi_sha1
: bytes
| A
: int64 10
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| oldimage
: bytes
| oi_usertext_timestamp
: bytes
| oi_user_text
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| oldimage
: bytes
| oi_usertext_timestamp
: bytes
| oi_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| page
: bytes
| name_title
: bytes
| page_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page
: bytes
| name_title
: bytes
| page_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page
: bytes
| page_len
: bytes
| page_len
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| page
: bytes
| page_random
: bytes
| page_random
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| page
: bytes
| page_redirect_namespace_len
: bytes
| page_is_redirect
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| page
: bytes
| page_redirect_namespace_len
: bytes
| page_namespace
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| page
: bytes
| page_redirect_namespace_len
: bytes
| page_len
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| page
: bytes
| PRIMARY
: bytes
| page_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page_props
: bytes
| pp_page_propname
: bytes
| pp_page
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page_props
: bytes
| pp_page_propname
: bytes
| pp_propname
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page_props
: bytes
| pp_propname_page
: bytes
| pp_propname
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page_props
: bytes
| pp_propname_page
: bytes
| pp_page
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page_props
: bytes
| pp_propname_sortkey_page
: bytes
| pp_propname
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page_props
: bytes
| pp_propname_sortkey_page
: bytes
| pp_sortkey
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page_props
: bytes
| pp_propname_sortkey_page
: bytes
| pp_page
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page_restrictions
: bytes
| pr_cascade
: bytes
| pr_cascade
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| page_restrictions
: bytes
| pr_level
: bytes
| pr_level
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| page_restrictions
: bytes
| pr_pagetype
: bytes
| pr_page
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page_restrictions
: bytes
| pr_pagetype
: bytes
| pr_type
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| page_restrictions
: bytes
| pr_typelevel
: bytes
| pr_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| page_restrictions
: bytes
| pr_typelevel
: bytes
| pr_level
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| page_restrictions
: bytes
| PRIMARY
: bytes
| pr_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| pagelinks
: bytes
| pl_backlinks_namespace
: bytes
| pl_from_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| pagelinks
: bytes
| pl_backlinks_namespace
: bytes
| pl_namespace
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| pagelinks
: bytes
| pl_backlinks_namespace
: bytes
| pl_title
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| pagelinks
: bytes
| pl_backlinks_namespace
: bytes
| pl_from
: bytes
| A
: int64 -1
: null
: int64 4
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| pagelinks
: bytes
| pl_namespace
: bytes
| pl_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| pagelinks
: bytes
| pl_namespace
: bytes
| pl_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| pagelinks
: bytes
| pl_namespace
: bytes
| pl_from
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| pagelinks
: bytes
| PRIMARY
: bytes
| pl_from
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| pagelinks
: bytes
| PRIMARY
: bytes
| pl_namespace
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| pagelinks
: bytes
| PRIMARY
: bytes
| pl_title
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| protected_titles
: bytes
| pt_namespace_title
: bytes
| pt_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| protected_titles
: bytes
| pt_namespace_title
: bytes
| pt_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| protected_titles
: bytes
| pt_timestamp
: bytes
| pt_timestamp
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| querycache
: bytes
| qc_type
: bytes
| qc_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| querycache
: bytes
| qc_type
: bytes
| qc_value
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| querycache_info
: bytes
| PRIMARY
: bytes
| qci_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| querycachetwo
: bytes
| qcc_title
: bytes
| qcc_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| querycachetwo
: bytes
| qcc_title
: bytes
| qcc_namespace
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| querycachetwo
: bytes
| qcc_title
: bytes
| qcc_title
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| querycachetwo
: bytes
| qcc_titletwo
: bytes
| qcc_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| querycachetwo
: bytes
| qcc_titletwo
: bytes
| qcc_namespacetwo
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| querycachetwo
: bytes
| qcc_titletwo
: bytes
| qcc_titletwo
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| querycachetwo
: bytes
| qcc_type
: bytes
| qcc_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| querycachetwo
: bytes
| qcc_type
: bytes
| qcc_value
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| new_name_timestamp
: bytes
| rc_new
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| new_name_timestamp
: bytes
| rc_namespace
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| new_name_timestamp
: bytes
| rc_timestamp
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| PRIMARY
: bytes
| rc_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_actor
: bytes
| rc_actor
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_actor
: bytes
| rc_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_cur_id
: bytes
| rc_cur_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_ip
: bytes
| rc_ip
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_name_type_patrolled_timestamp
: bytes
| rc_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_name_type_patrolled_timestamp
: bytes
| rc_type
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_name_type_patrolled_timestamp
: bytes
| rc_patrolled
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_name_type_patrolled_timestamp
: bytes
| rc_timestamp
: bytes
| A
: int64 -1
: null
: int64 4
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_namespace_title_timestamp
: bytes
| rc_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_namespace_title_timestamp
: bytes
| rc_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_namespace_title_timestamp
: bytes
| rc_timestamp
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_ns_actor
: bytes
| rc_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_ns_actor
: bytes
| rc_actor
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_ns_usertext
: bytes
| rc_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_ns_usertext
: bytes
| rc_user_text
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_timestamp
: bytes
| rc_timestamp
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_user_text
: bytes
| rc_user_text
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| recentchanges
: bytes
| rc_user_text
: bytes
| rc_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| redirect
: bytes
| PRIMARY
: bytes
| rd_from
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| redirect
: bytes
| rd_ns_title
: bytes
| rd_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| redirect
: bytes
| rd_ns_title
: bytes
| rd_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| redirect
: bytes
| rd_ns_title
: bytes
| rd_from
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| page_timestamp
: bytes
| rev_page
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| page_timestamp
: bytes
| rev_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| page_user_timestamp
: bytes
| rev_page
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| page_user_timestamp
: bytes
| rev_user
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| page_user_timestamp
: bytes
| rev_timestamp
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| PRIMARY
: bytes
| rev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| rev_page_id
: bytes
| rev_page
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| rev_page_id
: bytes
| rev_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| rev_timestamp
: bytes
| rev_timestamp
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| user_timestamp
: bytes
| rev_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| user_timestamp
: bytes
| rev_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| usertext_timestamp
: bytes
| rev_user_text
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision
: bytes
| usertext_timestamp
: bytes
| rev_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision_actor_temp
: bytes
| actor_timestamp
: bytes
| revactor_actor
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision_actor_temp
: bytes
| actor_timestamp
: bytes
| revactor_timestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision_actor_temp
: bytes
| page_actor_timestamp
: bytes
| revactor_page
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision_actor_temp
: bytes
| page_actor_timestamp
: bytes
| revactor_actor
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision_actor_temp
: bytes
| page_actor_timestamp
: bytes
| revactor_timestamp
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| revision_actor_temp
: bytes
| PRIMARY
: bytes
| revactor_rev
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| revision_actor_temp
: bytes
| PRIMARY
: bytes
| revactor_actor
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| revision_actor_temp
: bytes
| revactor_rev
: bytes
| revactor_rev
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| revision_comment_temp
: bytes
| PRIMARY
: bytes
| revcomment_rev
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| revision_comment_temp
: bytes
| PRIMARY
: bytes
| revcomment_comment_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| revision_comment_temp
: bytes
| revcomment_rev
: bytes
| revcomment_rev
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| searchindex
: bytes
| si_page
: bytes
| si_page
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| searchindex
: bytes
| si_text
: bytes
| si_text
: null
: int64 -1
: null
: int64 1
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
< row
: bytes
| searchindex
: bytes
| si_title
: bytes
| si_title
: null
: int64 -1
: null
: int64 1
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
< row
: bytes
| site_identifiers
: bytes
| site_ids_key
: bytes
| si_key
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| site_identifiers
: bytes
| site_ids_site
: bytes
| si_site
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| site_identifiers
: bytes
| site_ids_type
: bytes
| si_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| site_identifiers
: bytes
| site_ids_type
: bytes
| si_key
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| site_stats
: bytes
| PRIMARY
: bytes
| ss_row_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| sites
: bytes
| PRIMARY
: bytes
| site_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| sites
: bytes
| sites_domain
: bytes
| site_domain
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| sites
: bytes
| sites_forward
: bytes
| site_forward
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| sites
: bytes
| sites_global_key
: bytes
| site_global_key
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| sites
: bytes
| sites_group
: bytes
| site_group
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| sites
: bytes
| sites_language
: bytes
| site_language
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| sites
: bytes
| sites_protocol
: bytes
| site_protocol
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| sites
: bytes
| sites_source
: bytes
| site_source
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| sites
: bytes
| sites_type
: bytes
| site_type
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| slot_roles
: bytes
| PRIMARY
: bytes
| role_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| slot_roles
: bytes
| role_name
: bytes
| role_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| slots
: bytes
| PRIMARY
: bytes
| slot_revision_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| slots
: bytes
| PRIMARY
: bytes
| slot_role_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| slots
: bytes
| slot_revision_origin_role
: bytes
| slot_revision_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| slots
: bytes
| slot_revision_origin_role
: bytes
| slot_origin
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| slots
: bytes
| slot_revision_origin_role
: bytes
| slot_role_id
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| tag_summary
: bytes
| PRIMARY
: bytes
| ts_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| tag_summary
: bytes
| tag_summary_log_id
: bytes
| ts_log_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| tag_summary
: bytes
| tag_summary_rc_id
: bytes
| ts_rc_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| tag_summary
: bytes
| tag_summary_rev_id
: bytes
| ts_rev_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| templatelinks
: bytes
| PRIMARY
: bytes
| tl_from
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| templatelinks
: bytes
| PRIMARY
: bytes
| tl_namespace
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| templatelinks
: bytes
| PRIMARY
: bytes
| tl_title
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| templatelinks
: bytes
| tl_backlinks_namespace
: bytes
| tl_from_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| templatelinks
: bytes
| tl_backlinks_namespace
: bytes
| tl_namespace
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| templatelinks
: bytes
| tl_backlinks_namespace
: bytes
| tl_title
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| templatelinks
: bytes
| tl_backlinks_namespace
: bytes
| tl_from
: bytes
| A
: int64 -1
: null
: int64 4
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| templatelinks
: bytes
| tl_namespace
: bytes
| tl_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| templatelinks
: bytes
| tl_namespace
: bytes
| tl_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| templatelinks
: bytes
| tl_namespace
: bytes
| tl_from
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| text
: bytes
| PRIMARY
: bytes
| old_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| transcache
: bytes
| PRIMARY
: bytes
| tc_url
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| updatelog
: bytes
| PRIMARY
: bytes
| ul_key
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| uploadstash
: bytes
| PRIMARY
: bytes
| us_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| uploadstash
: bytes
| us_key
: bytes
| us_key
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| uploadstash
: bytes
| us_timestamp
: bytes
| us_timestamp
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| uploadstash
: bytes
| us_user
: bytes
| us_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| user
: bytes
| PRIMARY
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| user
: bytes
| user_email
: bytes
| user_email
: bytes
| A
: int64 50
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| user
: bytes
| user_email_token
: bytes
| user_email_token
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| user
: bytes
| user_name
: bytes
| user_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| user_former_groups
: bytes
| PRIMARY
: bytes
| ufg_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| user_former_groups
: bytes
| PRIMARY
: bytes
| ufg_group
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| user_groups
: bytes
| PRIMARY
: bytes
| ug_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| user_groups
: bytes
| PRIMARY
: bytes
| ug_group
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| user_groups
: bytes
| ug_expiry
: bytes
| ug_expiry
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| user_groups
: bytes
| ug_group
: bytes
| ug_group
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| user_newtalk
: bytes
| un_user_id
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| user_newtalk
: bytes
| un_user_ip
: bytes
| user_ip
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| user_properties
: bytes
| PRIMARY
: bytes
| up_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| user_properties
: bytes
| PRIMARY
: bytes
| up_property
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| user_properties
: bytes
| user_properties_property
: bytes
| up_property
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| valid_tag
: bytes
| PRIMARY
: bytes
| vt_tag
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| watchlist
: bytes
| namespace_title
: bytes
| wl_namespace
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| watchlist
: bytes
| namespace_title
: bytes
| wl_title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| watchlist
: bytes
| PRIMARY
: bytes
| wl_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| watchlist
: bytes
| wl_user
: bytes
| wl_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| watchlist
: bytes
| wl_user
: bytes
| wl_namespace
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| watchlist
: bytes
| wl_user
: bytes
| wl_title
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| watchlist
: bytes
| wl_user_notificationtimestamp
: bytes
| wl_user
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| watchlist
: bytes
| wl_user_notificationtimestamp
: bytes
| wl_notificationtimestamp
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
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
| sdl_rw_mediawiki_4
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
| actor
: bytes
| actor_id
: int64 1
: null
: bytes
| NO
: bytes
| bigint unsigned
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
| actor
: bytes
| actor_user
: int64 2
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| actor
: bytes
| actor_name
: int64 3
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| archive
: bytes
| ar_namespace
: int64 2
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
| archive
: bytes
| ar_title
: int64 3
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_comment
: int64 4
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(767)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_comment_id
: int64 5
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_user
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_user_text
: int64 7
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_actor
: int64 8
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_timestamp
: int64 9
: bytes
| '0x'
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_minor_edit
: int64 10
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_rev_id
: int64 11
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_text_id
: int64 12
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_deleted
: int64 13
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_len
: int64 14
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_page_id
: int64 15
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_parent_id
: int64 16
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_sha1
: int64 17
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_content_model
: int64 18
: null
: bytes
| YES
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| archive
: bytes
| ar_content_format
: int64 19
: null
: bytes
| YES
: bytes
| varbinary(64)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| bot_passwords
: bytes
| bp_user
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| bot_passwords
: bytes
| bp_app_id
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| bot_passwords
: bytes
| bp_password
: int64 3
: null
: bytes
| NO
: bytes
| tinyblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| bot_passwords
: bytes
| bp_token
: int64 4
: bytes
| '0x'
: bytes
| NO
: bytes
| binary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| bot_passwords
: bytes
| bp_restrictions
: int64 5
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| bot_passwords
: bytes
| bp_grants
: int64 6
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| category
: bytes
| cat_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| category
: bytes
| cat_title
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| category
: bytes
| cat_pages
: int64 3
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
| category
: bytes
| cat_subcats
: int64 4
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
| category
: bytes
| cat_files
: int64 5
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
| categorylinks
: bytes
| cl_from
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| categorylinks
: bytes
| cl_to
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| categorylinks
: bytes
| cl_sortkey
: int64 3
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(230)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| categorylinks
: bytes
| cl_sortkey_prefix
: int64 4
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| categorylinks
: bytes
| cl_timestamp
: int64 5
: null
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| categorylinks
: bytes
| cl_collation
: int64 6
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| categorylinks
: bytes
| cl_type
: int64 7
: bytes
| '0x70616765'
: bytes
| NO
: bytes
| enum('page','subcat','file')
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| change_tag
: bytes
| ct_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| change_tag
: bytes
| ct_rc_id
: int64 2
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
| change_tag
: bytes
| ct_log_id
: int64 3
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| change_tag
: bytes
| ct_rev_id
: int64 4
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| change_tag
: bytes
| ct_tag
: int64 5
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| change_tag
: bytes
| ct_params
: int64 6
: null
: bytes
| YES
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| comment
: bytes
| comment_id
: int64 1
: null
: bytes
| NO
: bytes
| bigint unsigned
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
| comment
: bytes
| comment_hash
: int64 2
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
|
: null
< row
: bytes
| comment
: bytes
| comment_text
: int64 3
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| comment
: bytes
| comment_data
: int64 4
: null
: bytes
| YES
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| content
: bytes
| content_id
: int64 1
: null
: bytes
| NO
: bytes
| bigint unsigned
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
| content
: bytes
| content_size
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| content
: bytes
| content_sha1
: int64 3
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| content
: bytes
| content_model
: int64 4
: null
: bytes
| NO
: bytes
| smallint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| content
: bytes
| content_address
: int64 5
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| content_models
: bytes
| model_id
: int64 1
: null
: bytes
| NO
: bytes
| smallint
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
| content_models
: bytes
| model_name
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(64)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| externallinks
: bytes
| el_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| externallinks
: bytes
| el_from
: int64 2
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| externallinks
: bytes
| el_to
: int64 3
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| externallinks
: bytes
| el_index
: int64 4
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| externallinks
: bytes
| el_index_60
: int64 5
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(60)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_id
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
| filearchive
: bytes
| fa_name
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_archive_name
: int64 3
: bytes
| ''
: bytes
| YES
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_storage_group
: int64 4
: null
: bytes
| YES
: bytes
| varbinary(16)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_storage_key
: int64 5
: bytes
| ''
: bytes
| YES
: bytes
| varbinary(64)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_deleted_user
: int64 6
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
| filearchive
: bytes
| fa_deleted_timestamp
: int64 7
: bytes
| '0x'
: bytes
| YES
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_deleted_reason
: int64 8
: bytes
| ''
: bytes
| YES
: bytes
| varbinary(767)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_deleted_reason_id
: int64 9
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_size
: int64 10
: bytes
| '0'
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_width
: int64 11
: bytes
| '0'
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
| filearchive
: bytes
| fa_height
: int64 12
: bytes
| '0'
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
| filearchive
: bytes
| fa_metadata
: int64 13
: null
: bytes
| YES
: bytes
| mediumblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_bits
: int64 14
: bytes
| '0'
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
| filearchive
: bytes
| fa_media_type
: int64 15
: null
: bytes
| YES
: bytes
| enum('UNKNOWN','BITMAP','DRAWING','AUDIO','VIDEO','MULTIMEDIA','OFFICE','TEXT','EXECUTABLE','ARCHIVE','3D')
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_major_mime
: int64 16
: bytes
| '0x756E6B6E6F776E'
: bytes
| YES
: bytes
| enum('unknown','application','audio','image','text','video','message','model','multipart','chemical')
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_minor_mime
: int64 17
: bytes
| '0x756E6B6E6F776E'
: bytes
| YES
: bytes
| varbinary(100)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_description
: int64 18
: bytes
| ''
: bytes
| YES
: bytes
| varbinary(767)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_description_id
: int64 19
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_user
: int64 20
: bytes
| '0'
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_user_text
: int64 21
: bytes
| ''
: bytes
| YES
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_actor
: int64 22
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_timestamp
: int64 23
: bytes
| '0x'
: bytes
| YES
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_deleted
: int64 24
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filearchive
: bytes
| fa_sha1
: int64 25
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_name
: int64 1
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_size
: int64 2
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_width
: int64 3
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
| image
: bytes
| img_height
: int64 4
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
| image
: bytes
| img_metadata
: int64 5
: null
: bytes
| NO
: bytes
| mediumblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_bits
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
| image
: bytes
| img_media_type
: int64 7
: null
: bytes
| YES
: bytes
| enum('UNKNOWN','BITMAP','DRAWING','AUDIO','VIDEO','MULTIMEDIA','OFFICE','TEXT','EXECUTABLE','ARCHIVE','3D')
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_major_mime
: int64 8
: bytes
| '0x756E6B6E6F776E'
: bytes
| NO
: bytes
| enum('unknown','application','audio','image','text','video','message','model','multipart','chemical')
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_minor_mime
: int64 9
: bytes
| '0x756E6B6E6F776E'
: bytes
| NO
: bytes
| varbinary(100)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_description
: int64 10
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(767)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_description_id
: int64 11
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_user
: int64 12
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_user_text
: int64 13
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_actor
: int64 14
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_timestamp
: int64 15
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image
: bytes
| img_sha1
: int64 16
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image_comment_temp
: bytes
| imgcomment_name
: int64 1
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| image_comment_temp
: bytes
| imgcomment_description_id
: int64 2
: null
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| imagelinks
: bytes
| il_from
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| imagelinks
: bytes
| il_from_namespace
: int64 2
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
| imagelinks
: bytes
| il_to
: int64 3
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| interwiki
: bytes
| iw_prefix
: int64 1
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| interwiki
: bytes
| iw_url
: int64 2
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| interwiki
: bytes
| iw_api
: int64 3
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| interwiki
: bytes
| iw_wikiid
: int64 4
: null
: bytes
| NO
: bytes
| varbinary(64)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| interwiki
: bytes
| iw_local
: int64 5
: null
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
| interwiki
: bytes
| iw_trans
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ip_changes
: bytes
| ipc_rev_id
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ip_changes
: bytes
| ipc_rev_timestamp
: int64 2
: bytes
| '0x'
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ip_changes
: bytes
| ipc_hex
: int64 3
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(35)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_id
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
| ipblocks
: bytes
| ipb_address
: int64 2
: null
: bytes
| NO
: bytes
| tinyblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_user
: int64 3
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_by
: int64 4
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_by_text
: int64 5
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_by_actor
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_reason
: int64 7
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(767)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_reason_id
: int64 8
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_timestamp
: int64 9
: bytes
| '0x'
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_auto
: int64 10
: bytes
| '0'
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
| ipblocks
: bytes
| ipb_anon_only
: int64 11
: bytes
| '0'
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
| ipblocks
: bytes
| ipb_create_account
: int64 12
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
| ipblocks
: bytes
| ipb_enable_autoblock
: int64 13
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
| ipblocks
: bytes
| ipb_expiry
: int64 14
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_range_start
: int64 15
: null
: bytes
| NO
: bytes
| tinyblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_range_end
: int64 16
: null
: bytes
| NO
: bytes
| tinyblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| ipblocks
: bytes
| ipb_deleted
: int64 17
: bytes
| '0'
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
| ipblocks
: bytes
| ipb_block_email
: int64 18
: bytes
| '0'
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
| ipblocks
: bytes
| ipb_allow_usertalk
: int64 19
: bytes
| '0'
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
| ipblocks
: bytes
| ipb_parent_block_id
: int64 20
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
| iwlinks
: bytes
| iwl_from
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| iwlinks
: bytes
| iwl_prefix
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(20)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| iwlinks
: bytes
| iwl_title
: int64 3
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| job
: bytes
| job_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| job
: bytes
| job_cmd
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(60)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| job
: bytes
| job_namespace
: int64 3
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
|
: null
< row
: bytes
| job
: bytes
| job_title
: int64 4
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| job
: bytes
| job_timestamp
: int64 5
: null
: bytes
| YES
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| job
: bytes
| job_params
: int64 6
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| job
: bytes
| job_random
: int64 7
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| job
: bytes
| job_attempts
: int64 8
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| job
: bytes
| job_token
: int64 9
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| job
: bytes
| job_token_timestamp
: int64 10
: null
: bytes
| YES
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| job
: bytes
| job_sha1
: int64 11
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| l10n_cache
: bytes
| lc_lang
: int64 1
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| l10n_cache
: bytes
| lc_key
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| l10n_cache
: bytes
| lc_value
: int64 3
: null
: bytes
| NO
: bytes
| mediumblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| langlinks
: bytes
| ll_from
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| langlinks
: bytes
| ll_lang
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(20)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| langlinks
: bytes
| ll_title
: int64 3
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| log_search
: bytes
| ls_field
: int64 1
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| log_search
: bytes
| ls_value
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| log_search
: bytes
| ls_log_id
: int64 3
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| logging
: bytes
| log_type
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_action
: int64 3
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_timestamp
: int64 4
: bytes
| '0x3139373030313031303030303030'
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_user
: int64 5
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_user_text
: int64 6
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_actor
: int64 7
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_namespace
: int64 8
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
| logging
: bytes
| log_title
: int64 9
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_page
: int64 10
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_comment
: int64 11
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(767)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_comment_id
: int64 12
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_params
: int64 13
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| logging
: bytes
| log_deleted
: int64 14
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| module_deps
: bytes
| md_module
: int64 1
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| module_deps
: bytes
| md_skin
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| module_deps
: bytes
| md_deps
: int64 3
: null
: bytes
| NO
: bytes
| mediumblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| objectcache
: bytes
| keyname
: int64 1
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| objectcache
: bytes
| value
: int64 2
: null
: bytes
| YES
: bytes
| mediumblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| objectcache
: bytes
| exptime
: int64 3
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_name
: int64 1
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_archive_name
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_size
: int64 3
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_width
: int64 4
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
| oldimage
: bytes
| oi_height
: int64 5
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
| oldimage
: bytes
| oi_bits
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
| oldimage
: bytes
| oi_description
: int64 7
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(767)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_description_id
: int64 8
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_user
: int64 9
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_user_text
: int64 10
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_actor
: int64 11
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_timestamp
: int64 12
: bytes
| '0x'
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_metadata
: int64 13
: null
: bytes
| NO
: bytes
| mediumblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_media_type
: int64 14
: null
: bytes
| YES
: bytes
| enum('UNKNOWN','BITMAP','DRAWING','AUDIO','VIDEO','MULTIMEDIA','OFFICE','TEXT','EXECUTABLE','ARCHIVE','3D')
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_major_mime
: int64 15
: bytes
| '0x756E6B6E6F776E'
: bytes
| NO
: bytes
| enum('unknown','application','audio','image','text','video','message','model','multipart','chemical')
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_minor_mime
: int64 16
: bytes
| '0x756E6B6E6F776E'
: bytes
| NO
: bytes
| varbinary(100)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_deleted
: int64 17
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| oldimage
: bytes
| oi_sha1
: int64 18
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| page
: bytes
| page_namespace
: int64 2
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
|
: null
< row
: bytes
| page
: bytes
| page_title
: int64 3
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_restrictions
: int64 4
: null
: bytes
| NO
: bytes
| tinyblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_is_redirect
: int64 5
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_is_new
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_random
: int64 7
: null
: bytes
| NO
: bytes
| double unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_touched
: int64 8
: bytes
| '0x'
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_links_updated
: int64 9
: null
: bytes
| YES
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_latest
: int64 10
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_len
: int64 11
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_content_model
: int64 12
: null
: bytes
| YES
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page
: bytes
| page_lang
: int64 13
: null
: bytes
| YES
: bytes
| varbinary(35)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page_props
: bytes
| pp_page
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
|
: null
< row
: bytes
| page_props
: bytes
| pp_propname
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(60)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page_props
: bytes
| pp_value
: int64 3
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page_props
: bytes
| pp_sortkey
: int64 4
: null
: bytes
| YES
: bytes
| float
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page_restrictions
: bytes
| pr_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| page_restrictions
: bytes
| pr_page
: int64 2
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
|
: null
< row
: bytes
| page_restrictions
: bytes
| pr_type
: int64 3
: null
: bytes
| NO
: bytes
| varbinary(60)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page_restrictions
: bytes
| pr_level
: int64 4
: null
: bytes
| NO
: bytes
| varbinary(60)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page_restrictions
: bytes
| pr_cascade
: int64 5
: null
: bytes
| NO
: bytes
| tinyint
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page_restrictions
: bytes
| pr_user
: int64 6
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| page_restrictions
: bytes
| pr_expiry
: int64 7
: null
: bytes
| YES
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| pagelinks
: bytes
| pl_from
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| pagelinks
: bytes
| pl_from_namespace
: int64 2
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
| pagelinks
: bytes
| pl_namespace
: int64 3
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
| pagelinks
: bytes
| pl_title
: int64 4
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| protected_titles
: bytes
| pt_namespace
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
|
: null
< row
: bytes
| protected_titles
: bytes
| pt_title
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| protected_titles
: bytes
| pt_user
: int64 3
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| protected_titles
: bytes
| pt_reason
: int64 4
: bytes
| ''
: bytes
| YES
: bytes
| varbinary(767)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| protected_titles
: bytes
| pt_reason_id
: int64 5
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| protected_titles
: bytes
| pt_timestamp
: int64 6
: null
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| protected_titles
: bytes
| pt_expiry
: int64 7
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| protected_titles
: bytes
| pt_create_perm
: int64 8
: null
: bytes
| NO
: bytes
| varbinary(60)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| querycache
: bytes
| qc_type
: int64 1
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| querycache
: bytes
| qc_value
: int64 2
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| querycache
: bytes
| qc_namespace
: int64 3
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
| querycache
: bytes
| qc_title
: int64 4
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| querycache_info
: bytes
| qci_type
: int64 1
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| querycache_info
: bytes
| qci_timestamp
: int64 2
: bytes
| '0x3139373030313031303030303030'
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| querycachetwo
: bytes
| qcc_type
: int64 1
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| querycachetwo
: bytes
| qcc_value
: int64 2
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| querycachetwo
: bytes
| qcc_namespace
: int64 3
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
| querycachetwo
: bytes
| qcc_title
: int64 4
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| querycachetwo
: bytes
| qcc_namespacetwo
: int64 5
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
| querycachetwo
: bytes
| qcc_titletwo
: int64 6
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_id
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
| recentchanges
: bytes
| rc_timestamp
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_user
: int64 3
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_user_text
: int64 4
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_actor
: int64 5
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_namespace
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
| recentchanges
: bytes
| rc_title
: int64 7
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_comment
: int64 8
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(767)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_comment_id
: int64 9
: bytes
| '0'
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_minor
: int64 10
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_bot
: int64 11
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_new
: int64 12
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_cur_id
: int64 13
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_this_oldid
: int64 14
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_last_oldid
: int64 15
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_type
: int64 16
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_source
: int64 17
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(16)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_patrolled
: int64 18
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_ip
: int64 19
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(40)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_old_len
: int64 20
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
| recentchanges
: bytes
| rc_new_len
: int64 21
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
| recentchanges
: bytes
| rc_deleted
: int64 22
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_logid
: int64 23
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_log_type
: int64 24
: null
: bytes
| YES
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_log_action
: int64 25
: null
: bytes
| YES
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| recentchanges
: bytes
| rc_params
: int64 26
: null
: bytes
| YES
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| redirect
: bytes
| rd_from
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| redirect
: bytes
| rd_namespace
: int64 2
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
| redirect
: bytes
| rd_title
: int64 3
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| redirect
: bytes
| rd_interwiki
: int64 4
: null
: bytes
| YES
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| redirect
: bytes
| rd_fragment
: int64 5
: null
: bytes
| YES
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| revision
: bytes
| rev_page
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_text_id
: int64 3
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_comment
: int64 4
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(767)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_user
: int64 5
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_user_text
: int64 6
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_timestamp
: int64 7
: bytes
| '0x'
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_minor_edit
: int64 8
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_deleted
: int64 9
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_len
: int64 10
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_parent_id
: int64 11
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_sha1
: int64 12
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_content_model
: int64 13
: null
: bytes
| YES
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision
: bytes
| rev_content_format
: int64 14
: null
: bytes
| YES
: bytes
| varbinary(64)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision_actor_temp
: bytes
| revactor_rev
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision_actor_temp
: bytes
| revactor_actor
: int64 2
: null
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision_actor_temp
: bytes
| revactor_timestamp
: int64 3
: bytes
| '0x'
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision_actor_temp
: bytes
| revactor_page
: int64 4
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision_comment_temp
: bytes
| revcomment_rev
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| revision_comment_temp
: bytes
| revcomment_comment_id
: int64 2
: null
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| searchindex
: bytes
| si_page
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| searchindex
: bytes
| si_title
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| searchindex
: bytes
| si_text
: int64 3
: null
: bytes
| NO
: bytes
| mediumtext
: bytes
| utf8mb3
: bytes
| utf8mb3_general_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| site_identifiers
: bytes
| si_site
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| site_identifiers
: bytes
| si_type
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| site_identifiers
: bytes
| si_key
: int64 3
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| site_stats
: bytes
| ss_row_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| site_stats
: bytes
| ss_total_edits
: int64 2
: null
: bytes
| YES
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| site_stats
: bytes
| ss_good_articles
: int64 3
: null
: bytes
| YES
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| site_stats
: bytes
| ss_total_pages
: int64 4
: null
: bytes
| YES
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| site_stats
: bytes
| ss_users
: int64 5
: null
: bytes
| YES
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| site_stats
: bytes
| ss_active_users
: int64 6
: null
: bytes
| YES
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| site_stats
: bytes
| ss_images
: int64 7
: null
: bytes
| YES
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sites
: bytes
| site_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| sites
: bytes
| site_global_key
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sites
: bytes
| site_type
: int64 3
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sites
: bytes
| site_group
: int64 4
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sites
: bytes
| site_source
: int64 5
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sites
: bytes
| site_language
: int64 6
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sites
: bytes
| site_protocol
: int64 7
: null
: bytes
| NO
: bytes
| varbinary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sites
: bytes
| site_domain
: int64 8
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sites
: bytes
| site_data
: int64 9
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| sites
: bytes
| site_forward
: int64 10
: null
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
| sites
: bytes
| site_config
: int64 11
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| slot_roles
: bytes
| role_id
: int64 1
: null
: bytes
| NO
: bytes
| smallint
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
| slot_roles
: bytes
| role_name
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(64)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| slots
: bytes
| slot_revision_id
: int64 1
: null
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| slots
: bytes
| slot_role_id
: int64 2
: null
: bytes
| NO
: bytes
| smallint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| slots
: bytes
| slot_content_id
: int64 3
: null
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| slots
: bytes
| slot_origin
: int64 4
: null
: bytes
| NO
: bytes
| bigint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| tag_summary
: bytes
| ts_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| tag_summary
: bytes
| ts_rc_id
: int64 2
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
| tag_summary
: bytes
| ts_log_id
: int64 3
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| tag_summary
: bytes
| ts_rev_id
: int64 4
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| tag_summary
: bytes
| ts_tags
: int64 5
: null
: bytes
| NO
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| templatelinks
: bytes
| tl_from
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| templatelinks
: bytes
| tl_from_namespace
: int64 2
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
| templatelinks
: bytes
| tl_namespace
: int64 3
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
| templatelinks
: bytes
| tl_title
: int64 4
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| text
: bytes
| old_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| text
: bytes
| old_text
: int64 2
: null
: bytes
| NO
: bytes
| mediumblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| text
: bytes
| old_flags
: int64 3
: null
: bytes
| NO
: bytes
| tinyblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| transcache
: bytes
| tc_url
: int64 1
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| transcache
: bytes
| tc_contents
: int64 2
: null
: bytes
| YES
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| transcache
: bytes
| tc_time
: int64 3
: null
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| updatelog
: bytes
| ul_key
: int64 1
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| updatelog
: bytes
| ul_value
: int64 2
: null
: bytes
| YES
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| uploadstash
: bytes
| us_user
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_key
: int64 3
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_orig_path
: int64 4
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_path
: int64 5
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_source_type
: int64 6
: null
: bytes
| YES
: bytes
| varbinary(50)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_timestamp
: int64 7
: null
: bytes
| NO
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_status
: int64 8
: null
: bytes
| NO
: bytes
| varbinary(50)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_chunk_inx
: int64 9
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_props
: int64 10
: null
: bytes
| YES
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_size
: int64 11
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_sha1
: int64 12
: null
: bytes
| NO
: bytes
| varbinary(31)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_mime
: int64 13
: null
: bytes
| YES
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_media_type
: int64 14
: null
: bytes
| YES
: bytes
| enum('UNKNOWN','BITMAP','DRAWING','AUDIO','VIDEO','MULTIMEDIA','OFFICE','TEXT','EXECUTABLE','ARCHIVE','3D')
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_image_width
: int64 15
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_image_height
: int64 16
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploadstash
: bytes
| us_image_bits
: int64 17
: null
: bytes
| YES
: bytes
| smallint unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| user
: bytes
| user_name
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_real_name
: int64 3
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_password
: int64 4
: null
: bytes
| NO
: bytes
| tinyblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_newpassword
: int64 5
: null
: bytes
| NO
: bytes
| tinyblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_newpass_time
: int64 6
: null
: bytes
| YES
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_email
: int64 7
: null
: bytes
| NO
: bytes
| tinyblob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_touched
: int64 8
: bytes
| '0x'
: bytes
| NO
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_token
: int64 9
: bytes
| '0x'
: bytes
| NO
: bytes
| binary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_email_authenticated
: int64 10
: null
: bytes
| YES
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_email_token
: int64 11
: null
: bytes
| YES
: bytes
| binary(32)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_email_token_expires
: int64 12
: null
: bytes
| YES
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_registration
: int64 13
: null
: bytes
| YES
: bytes
| binary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user
: bytes
| user_editcount
: int64 14
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
| user
: bytes
| user_password_expires
: int64 15
: null
: bytes
| YES
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_former_groups
: bytes
| ufg_user
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_former_groups
: bytes
| ufg_group
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_groups
: bytes
| ug_user
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_groups
: bytes
| ug_group
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_groups
: bytes
| ug_expiry
: int64 3
: null
: bytes
| YES
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_newtalk
: bytes
| user_id
: int64 1
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_newtalk
: bytes
| user_ip
: int64 2
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(40)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_newtalk
: bytes
| user_last_timestamp
: int64 3
: null
: bytes
| YES
: bytes
| varbinary(14)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_properties
: bytes
| up_user
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_properties
: bytes
| up_property
: int64 2
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| user_properties
: bytes
| up_value
: int64 3
: null
: bytes
| YES
: bytes
| blob
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| valid_tag
: bytes
| vt_tag
: int64 1
: null
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| watchlist
: bytes
| wl_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
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
| watchlist
: bytes
| wl_user
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| watchlist
: bytes
| wl_namespace
: int64 3
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
| watchlist
: bytes
| wl_title
: int64 4
: bytes
| ''
: bytes
| NO
: bytes
| varbinary(255)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| watchlist
: bytes
| wl_notificationtimestamp
: int64 5
: null
: bytes
| YES
: bytes
| varbinary(14)
: bytes
|
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
| sdl_rw_mediawiki_4
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
| sdl_rw_mediawiki_4
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
| sdl_rw_mediawiki_4
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
| sdl_rw_mediawiki_4
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
| sdl_rw_mediawiki_4
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
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
| sdl_rw_mediawiki_4
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
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
| sdl_rw_mediawiki_4
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
| sdl_rw_mediawiki_4
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
| sdl_rw_mediawiki_4
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
| actor
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| archive
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| bot_passwords
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| category
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| categorylinks
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| change_tag
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| comment
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| content
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| content_models
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| externallinks
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| filearchive
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| image
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| image_comment_temp
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| imagelinks
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| interwiki
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| ip_changes
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| ipblocks
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| iwlinks
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| job
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| l10n_cache
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| langlinks
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| log_search
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| logging
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| module_deps
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| objectcache
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| oldimage
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| page
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| page_props
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| page_restrictions
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| pagelinks
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| protected_titles
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| querycache
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| querycache_info
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| querycachetwo
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| recentchanges
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| redirect
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| revision
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
| max_rows=10000000 avg_row_length=1024
: bytes
| ''
: bytes
| binary
< row
: bytes
| revision_actor_temp
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| revision_comment_temp
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| searchindex
: bytes
| BASE TABLE
: bytes
| MyISAM
: bytes
| utf8mb3_general_ci
: bytes
| 0
: bytes
| 0
: bytes
| 1024
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb3
< row
: bytes
| site_identifiers
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| site_stats
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| sites
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| slot_roles
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| slots
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| tag_summary
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| templatelinks
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| text
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
| max_rows=10000000 avg_row_length=10240
: bytes
| ''
: bytes
| binary
< row
: bytes
| transcache
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| updatelog
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| uploadstash
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| user
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| user_former_groups
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| user_groups
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| user_newtalk
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| user_properties
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| valid_tag
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
< row
: bytes
| watchlist
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| binary
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| binary
> query
|
| 		SELECT t.NAME
| 		FROM information_schema.INNODB_TABLES t
| 		LEFT JOIN information_schema.INNODB_TABLESPACES ts ON t.SPACE = ts.SPACE
| 		WHERE t.NAME LIKE CONCAT(?, '/%') AND (ts.SPACE_TYPE IS NULL OR ts.SPACE_TYPE <> 'Single')
: string
| sdl_rw_mediawiki_4
< result
< column "NAME" "VARCHAR" "string"
> query
| SHOW CREATE TABLE `sdl_rw_mediawiki_4`.`searchindex`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| searchindex
: bytes
| CREATE TABLE `searchindex` (
|   `si_page` int unsigned NOT NULL,
|   `si_title` varchar(255) NOT NULL DEFAULT '',
|   `si_text` mediumtext NOT NULL,
|   UNIQUE KEY `si_page` (`si_page`),
|   FULLTEXT KEY `si_title` (`si_title`),
|   FULLTEXT KEY `si_text` (`si_text`)
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb3
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| sdl_rw_mediawiki_4
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
