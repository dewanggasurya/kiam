package acm

import (
	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/dbflex/orm"
)

type Policy struct {
	orm.DataModelBase `bson:"-" json:"-"`
	ID                string `bson:"_id,omitempty" json:"_id,omitempty"`
	Name              string
}

func (g *Policy) TableName() string {
	return "ACMPolicy"
}

func (g *Policy) GetID(_ dbflex.IConnection) ([]string, []interface{}) {
	return []string{"_id"}, []interface{}{g.ID}
}

func (g *Policy) SetID(keys ...interface{}) {
	if len(keys) > 0 {
		g.ID = keys[0].(string)
	}
}
