package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/perowong/peroblogo/scripts"
)

var createCommentTableSql string = `CREATE TABLE Comment (
	ID int NOT NULL AUTO_INCREMENT,
	BlogID char(16) NOT NULL, -- index
	ParentID int NOT NULL,
	ReplyID int NOT NULL,
	FromUid char(36) NOT NULL,
	FromNickname varchar(100) NOT NULL,
	FromEmail varchar(255),
	ToUid char(36),
	ToNickname varchar(100),
	ToEmail varchar(255),
	Likes int NOT NULL DEFAULT 0,
	Content varchar(600) NOT NULL DEFAULT '',
	SubCount int NOT NULL DEFAULT 0,
	IsTop tinyint NOT NULL DEFAULT 0, -- 0: no, 1: yes
	Ct timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (ID),
	INDEX idx_blogid (BlogID),
	INDEX idx_parentid (ParentID)
) ENGINE=InnoDB
	AUTO_INCREMENT=10000
	DEFAULT
	CHARSET=utf8mb4
	COLLATE=utf8mb4_0900_ai_ci;`

var createLikeTableSql string = `CREATE TABLE UserLike (
	ID int NOT NULL AUTO_INCREMENT,
	SubjectID int NOT NULL, -- index, refer to comment id
	SubjectType tinyint NOT NULL, -- 1: blog, 2: comment
	FromUid char(36) NOT NULL,
	FromEmail varchar(255),
	Liked tinyint NOT NULL DEFAULT 0, -- 0: unlike, 1: like
	PRIMARY KEY (ID),
	INDEX idx_subjectid (SubjectID)
) ENGINE=InnoDB
	AUTO_INCREMENT=10000
	DEFAULT
	CHARSET=utf8mb4
	COLLATE=utf8mb4_0900_ai_ci;`

func main() {
	sqlCtx := scripts.GetSqlExecContext()

	sqlCtx(createCommentTableSql)
	sqlCtx(createLikeTableSql)
}
