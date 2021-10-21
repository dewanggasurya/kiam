package kiam

import (
	"time"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/dbflex/orm"
	"github.com/eaciit/toolkit"
)

type Session struct {
	orm.DataModelBase `json:"-"`
	SessionID         string
	ReferenceID       string
	Data              toolkit.M
	LastUpdate        time.Time
	Duration          int
}

// TableName get model table name
func (o *Session) TableName() string {
	return "sessions"
}

// GetID get model id
func (o *Session) GetID(_ dbflex.IConnection) ([]string, []interface{}) {
	return []string{"SessionID"}, []interface{}{o.SessionID}
}

// SetID set model  id
func (o *Session) SetID(keys ...interface{}) {
	if len(keys) > 0 {
		o.SessionID = keys[0].(string)
	}
}
