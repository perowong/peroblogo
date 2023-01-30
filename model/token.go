package model

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/perowong/peroblogo/utils"
)

type UserToken struct {
	UserID     int64     `db:"user_id"`
	Token      string    `db:"token"`
	Ct         time.Time `db:"ct"`
	Ut         time.Time `db:"ut"`
	ExpireTime time.Time `db:"expire_time"`
}

func (m *Model) AddUserToken(userID int64) (token string, err error) {
	now := time.Now()
	token = utils.GetMD5Hash(strconv.FormatInt(userID+now.Unix(), 10))

	_, err = m.DB.Exec(`
		INSERT INTO user_token (
			user_id,
			token,
			expire_time
		) VALUES (?, ?, ?)
	`, userID,
		token,
		now.UTC().AddDate(0, 0, 7),
	)

	return
}

func (m *Model) GetUserToken(token string) (userToken *UserToken, err error) {
	userToken = &UserToken{}
	err = m.DB.QueryRow(`
		SELECT user_id,token,expire_time FROM user_token WHERE token=?
	`, token).Scan(
		&userToken.UserID,
		&userToken.Token,
		&userToken.ExpireTime,
	)

	return
}

func (m *Model) CheckUserTokenExistBy(userID int64) (bool, error) {
	userToken := &UserToken{}
	err := m.DB.Get(userToken, `SELECT token FROM user_token WHERE user_id=?`, userID)

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func (m *Model) UpdateUserTokenExpireTime(token string) (err error) {
	_, err = m.DB.Exec(
		`UPDATE user_token set expire_time=? WHERE token=?`,
		time.Now().UTC().AddDate(0, 0, 7),
		token,
	)

	return
}

func (m *Model) UpdateUserTokenByUserID(userID int64) (token string, err error) {
	now := time.Now()
	token = utils.GetMD5Hash(strconv.FormatInt(userID+now.Unix(), 10))

	_, err = m.DB.Exec(
		`UPDATE user_token set expire_time=?,token=? WHERE user_id=?`,
		now.UTC().AddDate(0, 0, 7),
		token,
		userID,
	)

	return
}
