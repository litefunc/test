package model

import (
	"cloud/lib/null"
)

type Account struct {
	Bio
	Pwd string `json:"pwd"`
}

type Bio struct {
	ID      uint64      `json:"id"`
	Email   string      `json:"email"`
	Name    string      `json:"name"`
	Company uint64      `json:"com_id" gorm:"column:com"`
	Admin   bool        `json:"admin"`
	Type    string      `json:"type"`
	Phone   string      `json:"phone"`
	Note    null.String `json:"note"`
	Photo   null.String `json:"photo"`
}

type Accounts []Account

func (self Account) Json() string {
	return Json(self)
}

func (self Accounts) Json() string {
	return Json(self)
}
