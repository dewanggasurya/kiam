package acm

import (
	"fmt"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/dbflex/orm"
)

type AccessForScope string

const (
	AccessForUserKind string = "UserKind"
	AccessForUser            = "User"
	AccessForGroup           = "Group"
)

type access struct {
	orm.DataModelBase `bson:"-" json:"-" ecname:"-"`
	ID                string `bson:"_id" json:"_id" key:"1"`
	Kind              string
	ObjectID          string
	PolicyID          string
	Value             int
	DimensionID       string
}

func (a *access) TableName() string {
	return "ACMAccess"
}

func (g *access) GetID(_ dbflex.IConnection) ([]string, []interface{}) {
	return []string{"_id"}, []interface{}{g.ID}
}

func (g *access) SetID(keys ...interface{}) {
	if len(keys) > 0 {
		g.ID = keys[0].(string)
	}
}

func (m *manager) Grant(kind, objectid, policyid, dimid string, value int) (int, error) {
	w := dbflex.And(dbflex.Eq("Kind", kind),
		dbflex.Eq("ObjectID", objectid),
		dbflex.Eq("PolicyID", policyid),
		dbflex.Eq("DimensionID", dimid))
	a := new(access)
	m.h.GetByParm(a, dbflex.NewQueryParam().SetWhere(w).SetTake(1))
	if a.ID == "" {
		a.Kind = kind
		a.ObjectID = objectid
		a.PolicyID = policyid
	}
	a.Value = value
	if e := m.h.Save(a); e != nil {
		return 0, fmt.Errorf("fail grant access: %s", e.Error())
	}
	return a.Value, nil
}

func (m *manager) getAccess(kind, objectid, pid, dimid string) int {
	w := dbflex.And(dbflex.Eq("Kind", kind),
		dbflex.Eq("ObjectID", objectid),
		dbflex.Eq("PolicyID", pid),
		dbflex.Eq("DimensionID", dimid))
	a := new(access)
	m.h.GetByParm(a, dbflex.NewQueryParam().SetWhere(w).SetTake(1))
	return a.Value
}

func (m *manager) HasAccess(uid, pid, dimid string, value int) bool {
	u := new(User)
	m.h.GetByID(u, uid)
	if u.ID == "" || !u.Enable {
		return false
	}

	if v := m.getAccess(AccessForUser, uid, pid, dimid); v >= value {
		return true
	}

	if u.Kind != "" {
		if v := m.getAccess(AccessForUserKind, u.Kind, pid, dimid); v >= value {
			return true
		}
	}

	groups, e := m.UserGroups(uid, nil, "", "", -1)
	if e != nil {
		return false
	}
	for _, g := range groups {
		if v := m.getAccess(AccessForGroup, g.ID, pid, dimid); v >= value {
			return true
		}
	}

	return false
}
