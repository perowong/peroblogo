package model

import (
	"github.com/jmoiron/sqlx"
	"github.com/perowong/peroblogo/dao"
)

type Model struct {
	DB *sqlx.DB
}

func NewModel() (m *Model) {
	return &Model{
		DB: dao.DB,
	}
}
