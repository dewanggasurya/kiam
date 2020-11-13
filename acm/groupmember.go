package acm

import (
	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/dbflex/orm"
)

type GroupMember struct {
	orm.DataModelBase `bson:"-" json:"-"`
	ID                string `bson:"_id,omitempty" json:"_id,omitempty"`
	GroupID           string
	UserID            string
}

func (gm *GroupMember) TableName() string {
	return "ACMGroupMembers"
}

func (gm *GroupMember) GetID(_ dbflex.IConnection) ([]string, []interface{}) {
	return []string{"_id"}, []interface{}{gm.ID}
}

func (gm *GroupMember) SetID(keys ...interface{}) {
	if len(keys) > 0 {
		gm.ID = keys[0].(string)
	}
}
