CREATE TABLE IF NOT EXISTS tag(
		id   INT(11)      NOT NULL AUTO_INCREMENT,
		name VARCHAR(250) NOT NULL,
		PRIMARY KEY (id),
		UNIQUE KEY tag_name_UNIQUE (name))
		CHARACTER SET utf8mb4;
