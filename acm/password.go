package acm

import (
	"git.kanosolution.net/kano/dbflex/orm"
)

type passwrd struct {
	orm.DataModelBase `bson:"-" json:"-" ecname:"-" sql:"-"`
	ID                string `bson:"_id" json:"_id" sql:"id" ecname:"_id" key:"1"`
	Password          string
}

func (p *passwrd) TableName() string {
	return "ACMPasswds"
}

func (p *passwrd) SetID(keys ...interface{}) {
	p.ID = keys[0].(string)
}

func (m *manager) RecoverPassword(email string) error {
	return nil
}
