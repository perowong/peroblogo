-- create the databases
CREATE DATABASE IF NOT EXISTS peroblog;

USE peroblog;

CREATE TABLE IF NOT EXISTS comment (
	id int NOT NULL AUTO_INCREMENT,
	blog_id char(16) NOT NULL, -- index
	parent_id int NOT NULL,
	reply_id int NOT NULL,
	from_uid int,
	from_nickname varchar(100),
	from_avatar varchar(255),
	to_uid int,
	to_nickname varchar(100),
	to_avatar varchar(255),
	likes int NOT NULL DEFAULT 0,
	content varchar(600) NOT NULL DEFAULT '',
	sub_count int NOT NULL DEFAULT 0,
	is_top tinyint NOT NULL DEFAULT 0, -- 0: no, 1: yes
	ct timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	INDEX idx_blog_id (blog_id),
	INDEX idx_parent_id (parent_id)
) ENGINE=InnoDB
	AUTO_INCREMENT=10000
	DEFAULT
	CHARSET=utf8mb4
	COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS user_like (
	id int NOT NULL AUTO_INCREMENT,
	subject_id int NOT NULL, -- index, refer to comment id
	subject_type tinyint NOT NULL, -- 1: blog, 2: comment
	from_uid int,
	from_nickname varchar(100),
	from_avatar varchar(255),
	liked tinyint NOT NULL DEFAULT 0, -- 0: unlike, 1: like
	PRIMARY KEY (id),
	INDEX idx_subject_id (subject_id)
) ENGINE=InnoDB
	AUTO_INCREMENT=10000
	DEFAULT
	CHARSET=utf8mb4
	COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS user (
	id int NOT NULL AUTO_INCREMENT,
	openid varchar(255),
	auth_type char(36),
	nickname varchar(255) NOT NULL,
	avatar_url varchar(255) NOT NULL,
	email varchar(255),
	ct timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	INDEX idx_openid (openid)
) ENGINE=InnoDB
	AUTO_INCREMENT=10000
	DEFAULT
	CHARSET=utf8mb4
	COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS user_token (
	user_id int,
	token varchar(255),
	ct timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	ut timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	expire_time timestamp,
	PRIMARY KEY (user_id),
	INDEX idx_token (token)
) ENGINE=InnoDB
	DEFAULT
	CHARSET=utf8mb4
	COLLATE=utf8mb4_0900_ai_ci;
