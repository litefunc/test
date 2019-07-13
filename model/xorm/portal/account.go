package model

import (
	"cloud/lib/null"
)

type Account struct {
	Bio `xorm:"extends"`
	Pwd string `json:"pwd"`
}

type Bio struct {
	ID      uint64      `json:"id" xorm:"'id' pk autoincr"`
	Email   string      `json:"email"`
	Name    string      `json:"name"`
	Company uint64      `json:"com_id" xorm:"'com'"`
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

func (md Account) TableName() string {
	return "cloud.account"
}
