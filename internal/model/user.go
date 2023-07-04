package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64     `db:"id"`
	OpenID    string    `db:"openid"`
	AuthType  string    `db:"auth_type"`
	Nickname  string    `db:"nickname"`
	AvatarUrl string    `db:"avatar_url"`
	Email     string    `db:"email"`
	Ct        time.Time `json:"Ct,omitempty" db:"ct"`
}

func (m *Model) AddUser(user *User) (id int64, err error) {
	result, err := m.DB.Exec(`
		INSERT INTO user (
			openid,
			auth_type,
			nickname,
			avatar_url,
			email
		) VALUES (?, ?, ?, ?, ?)
	`, user.OpenID,
		user.AuthType,
		user.Nickname,
		user.AvatarUrl,
		user.Email,
	)
	if err != nil {
		return
	}
	id, err = result.LastInsertId()

	return
}

func (m *Model) GetUserIDBy(openID string) (int64, error) {
	user := &User{}
	err := m.DB.Get(user, "SELECT id,openid FROM user WHERE openid=?", openID)

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	if err == sql.ErrNoRows {
		return 0, nil
	}

	return user.ID, nil
}
