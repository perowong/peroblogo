package dao

import (
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Setup(dbCtx *sqlx.DB) {
	db = dbCtx
}

type Dao struct {
	DB *sqlx.DB
}

func NewDao() (d *Dao) {
	return &Dao{
		DB: db,
	}
}
