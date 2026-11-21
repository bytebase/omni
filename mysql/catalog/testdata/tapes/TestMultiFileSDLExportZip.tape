-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestMultiFileSDLExportZip default#1
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
