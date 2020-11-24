package acm

import (
	"errors"

	"git.kanosolution.net/kano/dbflex"
	dbf "git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/dbflex/orm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	orm.DataModelBase `bson:"-" json:"-"`
	ID                string `bson:"_id" json:"_id"`
	Name              string
}

func (g *Group) TableName() string {
	return "ACMGroups"
}

func (g *Group) GetID(_ dbflex.IConnection) ([]string, []interface{}) {
	return []string{"_id"}, []interface{}{g.ID}
}

func (g *Group) SetID(keys ...interface{}) {
	if len(keys) > 0 {
		g.ID = keys[0].(string)
	}
}

func (m *manager) CreateGroup(id, name string) (*Group, error) {
	g := new(Group)
	g.ID = id
	g.Name = name
	if e := m.h.Save(g); e != nil {
		return nil, errors.New("fail to create group: " + e.Error())
	}
	return g, nil
}

func (m *manager) AddUserToGroup(uid, gid string) error {
	u := new(User)
	g := new(Group)

	m.h.GetByID(u, uid)
	m.h.GetByID(g, gid)

	if u.ID == "" {
		return errors.New("fail to add user to group: invalid user")
	}

	if g.ID == "" {
		return errors.New("fail to add user to group: invalid group")
	}

	gm := new(GroupMember)
	w := dbflex.And(dbflex.Eq("GroupID", gid), dbflex.Eq("UserID", uid))
	m.h.GetByParm(gm, dbf.NewQueryParam().SetWhere(w))
	if gm.ID != "" {
		return errors.New("user has been added to group before")
	}

	gm.ID = primitive.NewObjectID().Hex()
	gm.UserID = uid
	gm.GroupID = gid
	if e := m.h.Save(gm); e != nil {
		return errors.New("fail to add user to group: " + e.Error())
	}

	return nil
}

func (m *manager) RemoveUserFromGroup(uid, gid string) error {
	w := dbflex.And(dbflex.Eq("GroupID", gid), dbflex.Eq("UserID", uid))
	cmd := dbflex.From(new(GroupMember).TableName()).Where(w).Delete()
	if _, e := m.h.Execute(cmd, nil); e != nil {
		return e
	}
	return nil
}

func (m *manager) GroupMembers(gid string, f *dbflex.Filter, sortBy string, lastIndex string, take int) ([]User, error) {
	res := []User{}
	if sortBy == "" {
		sortBy = "_id"
	}

	var w *dbflex.Filter
	if lastIndex != "" {
		w = dbflex.Lt(sortBy, lastIndex)
	} else {
		w = dbflex.And(dbflex.Lt(sortBy, lastIndex), f)
	}
	if f != nil {
		w = dbf.And(f, w)
	}

	if take == 0 {
		take = 20
	}

	cmd := dbflex.From(new(User).TableName()).OrderBy(sortBy).Where(w)
	if take > 0 {
		cmd.Take(take)
	}

	if _, e := m.h.Populate(cmd, &res); e != nil {
		return res, e
	}

	return res, nil
}
