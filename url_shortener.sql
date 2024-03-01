drop database if exists ulrshortener;
create database ulrshortener;
use ulrshortener;

CREATE TABLE Shortener (
    id VARCHAR(45) NOT NULL,
    url VARCHAR(255) NOT NULL,
    shortUrl VARCHAR(45) NOT NULL UNIQUE,
    expiredAt datetime NOT NULL,
    count INT Default 0,
    PRIMARY KEY (id)
);

DELIMITER $$
CREATE EVENT delete_url_every_day
	ON SCHEDULE EVERY 1 DAY
	STARTS CURRENT_TIMESTAMP + INTERVAL 1 DAY
	DO
	BEGIN
		DELETE from shortener where expiredAt < now();
	END $$

