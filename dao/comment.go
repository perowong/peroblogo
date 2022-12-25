package dao

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID           int64
	BlogID       string
	ParentID     int64
	ReplyID      int64
	FromUid      string
	FromNickname string
	FromEmail    string
	ToUid        string
	ToNickname   string
	ToEmail      string
	Likes        int64
	Content      string
	SubCount     int64
	IsTop        int
	Ct           time.Time
	Children     []*Comment
}

func (d *Dao) AddComment(comment *Comment) (id int64, err error) {
	result, err := d.DB.Exec(`
		INSERT INTO Comment (
			BlogID,
			ParentID,
			ReplyID,
			FromUid,
			FromNickname,
			FromEmail,
			ToUid,
			ToNickname,
			ToEmail,
			Content
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, comment.BlogID,
		comment.ParentID,
		comment.ReplyID,
		comment.FromUid,
		comment.FromNickname,
		comment.FromEmail,
		comment.ToUid,
		comment.ToNickname,
		comment.ToEmail,
		comment.Content,
	)
	if err != nil {
		return
	}

	// get the last insert id
	id, err = result.LastInsertId()
	if err != nil {
		return
	}

	return
}

func (d *Dao) CheckExistByID(id int64) (comment *Comment, err error) {
	err = d.DB.QueryRow(`
		SELECT ID FROM Comment WHERE ID=?
	`, id).Scan(&comment.ID)

	if err != nil {
		return
	}

	return
}

func (d *Dao) ReadComment(id int64) (comment *Comment, err error) {
	comment = &Comment{}
	err = d.DB.QueryRow(`
		SELECT * FROM Comment WHERE ID=?
	`, id).Scan(
		&comment.ID,
		&comment.BlogID,
		&comment.ParentID,
		&comment.ReplyID,
		&comment.FromUid,
		&comment.FromNickname,
		&comment.FromEmail,
		&comment.ToUid,
		&comment.ToNickname,
		&comment.ToEmail,
		&comment.Likes,
		&comment.Content,
		&comment.SubCount,
		&comment.IsTop,
		&comment.Ct,
	)

	if err != nil {
		return
	}

	return
}

func (d *Dao) UpdateSubCount(id int64, count int64) (err error) {
	_, err = d.DB.Exec(
		`UPDATE Comment set SubCount=? WHERE ID=?`,
		count,
		id,
	)
	if err != nil {
		return
	}

	return
}

func (d *Dao) getCommentList(rows *sql.Rows) (list []*Comment, err error) {
	for rows.Next() {
		comment := &Comment{}
		err = rows.Scan(
			&comment.ID,
			&comment.BlogID,
			&comment.ParentID,
			&comment.ReplyID,
			&comment.FromUid,
			&comment.FromNickname,
			&comment.FromEmail,
			&comment.ToUid,
			&comment.ToNickname,
			&comment.ToEmail,
			&comment.Likes,
			&comment.Content,
			&comment.SubCount,
			&comment.IsTop,
			&comment.Ct,
		)
		if err != nil {
			return
		}
		list = append(list, comment)
	}
	return
}

func (d *Dao) ListCommentByBlogID(blogId string) (list []*Comment, err error) {
	rows, err := d.DB.Query(`
		SELECT * FROM Comment WHERE BlogID=? AND ParentID=0
		ORDER BY Likes, Ct DESC
	`, blogId)
	if err != nil {
		return
	}
	defer rows.Close()

	return d.getCommentList(rows)
}

func (d *Dao) ListCommentByParentID(parentId int64) (list []*Comment, err error) {
	rows, err := d.DB.Query(`
		SELECT * FROM Comment WHERE ParentID=?
		ORDER BY Likes, Ct DESC
	`, parentId)
	if err != nil {
		return
	}
	defer rows.Close()

	return d.getCommentList(rows)
}
