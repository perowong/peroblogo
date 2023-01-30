package model

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID           int64     `db:"id"`
	BlogID       string    `db:"blog_id"`
	ParentID     int64     `db:"parent_id"`
	ReplyID      int64     `db:"reply_id"`
	FromUid      int64     `db:"from_uid"`
	FromNickname string    `db:"from_nickname"`
	FromAvatar   string    `db:"from_avatar"`
	ToUid        int64     `db:"to_uid"`
	ToNickname   string    `db:"to_nickname"`
	ToAvatar     string    `db:"to_avatar"`
	Likes        int64     `db:"likes"`
	Content      string    `db:"content"`
	SubCount     int64     `db:"sub_count"`
	IsTop        int       `db:"is_top"`
	Ct           time.Time `db:"ct"`
	Children     []*Comment
}

func (m *Model) AddComment(comment *Comment) (id int64, err error) {
	result, err := m.DB.Exec(`
		INSERT INTO comment (
			blog_id,
			parent_id,
			reply_id,
			from_uid,
			from_nickname,
			from_avatar,
			to_uid,
			to_nickname,
			to_avatar,
			content
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, comment.BlogID,
		comment.ParentID,
		comment.ReplyID,
		comment.FromUid,
		comment.FromNickname,
		comment.FromAvatar,
		comment.ToUid,
		comment.ToNickname,
		comment.ToAvatar,
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

func (m *Model) CheckCommentExistBy(id int64) (bool, error) {
	comment := &Comment{}
	err := m.DB.Get(comment, `SELECT id FROM comment WHERE id=?`, id)

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func (m *Model) ReadComment(id int64) (comment *Comment, err error) {
	comment = &Comment{}
	err = m.DB.QueryRow(`
		SELECT * FROM comment WHERE id=?
	`, id).Scan(
		&comment.ID,
		&comment.BlogID,
		&comment.ParentID,
		&comment.ReplyID,
		&comment.FromUid,
		&comment.FromNickname,
		&comment.FromAvatar,
		&comment.ToUid,
		&comment.ToNickname,
		&comment.ToAvatar,
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

func (m *Model) UpdateSubCount(id int64, count int64) (err error) {
	_, err = m.DB.Exec(
		`UPDATE comment set sub_count=? WHERE id=?`,
		count,
		id,
	)
	if err != nil {
		return
	}

	return
}

func (m *Model) getCommentList(rows *sql.Rows) (list []*Comment, err error) {
	for rows.Next() {
		comment := &Comment{}
		err = rows.Scan(
			&comment.ID,
			&comment.BlogID,
			&comment.ParentID,
			&comment.ReplyID,
			&comment.FromUid,
			&comment.FromNickname,
			&comment.FromAvatar,
			&comment.ToUid,
			&comment.ToNickname,
			&comment.ToAvatar,
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

func (m *Model) ListCommentByBlogID(blogId string) (list []*Comment, err error) {
	rows, err := m.DB.Query(`
		SELECT * FROM comment WHERE blog_id=? AND parent_id=0
		ORDER BY likes, ct DESC
	`, blogId)
	if err != nil {
		return
	}
	defer rows.Close()

	return m.getCommentList(rows)
}

func (m *Model) ListCommentByParentID(parentId int64) (list []*Comment, err error) {
	rows, err := m.DB.Query(`
		SELECT * FROM comment WHERE parent_id=?
		ORDER BY likes, ct DESC
	`, parentId)
	if err != nil {
		return
	}
	defer rows.Close()

	return m.getCommentList(rows)
}
