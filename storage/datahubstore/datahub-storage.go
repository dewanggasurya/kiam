package datahubstore

import (
	"fmt"
	"time"

	"git.kanosolution.net/kano/dbflex"
	"github.com/ariefdarmawan/datahub"
	"github.com/ariefdarmawan/kiam"
)

type store struct {
	hub *datahub.Hub
}

func NewStorage(hub *datahub.Hub) (*store, error) {
	s := new(store)
	s.hub = hub

	session := new(kiam.Session)
	keys, _ := session.GetID(nil)
	e := hub.EnsureTable(session.TableName(), keys, session)
	if e != nil {
		return nil, fmt.Errorf("kiam storage is failed ensuring table %s : %s", session.TableName(), e.Error())
	}

	return s, nil
}

func (s *store) Load(pool *kiam.SessionPool) error {
	sessions := []kiam.Session{}

	e := s.hub.Gets(new(kiam.Session), dbflex.NewQueryParam().SetWhere(
		dbflex.Gte(
			"LastUpdate",
			time.Now().Add(-1*24*time.Hour), // TODO this only my prediction, since load is not called yet
		),
	), &sessions)
	if e != nil {
		return e
	}

	for _, session := range sessions {
		if _, e = pool.Create(session.ReferenceID, session.Data, int(session.Duration)/int(time.Second)); e != nil {
			return fmt.Errorf("fail to create session. %s", e.Error())
		}
	}

	return nil
}

func (s *store) Store(pool *kiam.SessionPool) error {
	ids := pool.GetIDs()
	for _, id := range ids {
		if sess, ok := pool.GetBySessionID(id); ok {
			e := s.Write(sess)
			if e != nil {
				return fmt.Errorf("fail to store session %s. %s", id, e.Error())
			}
		}
	}
	return nil
}

func (s *store) Get(id string) (*kiam.Session, error) {
	sess := new(kiam.Session)

	e := s.hub.GetByID(sess, id)
	if e != nil {
		return nil, fmt.Errorf("fail read data %s. %s", id, e.Error())
	}

	return sess, nil
}

func (s *store) Write(sess *kiam.Session) error {
	e := s.hub.Save(sess)
	if e != nil {
		return fmt.Errorf("fail write session %s. %s", sess.SessionID, e.Error())
	}
	return nil
}

func (s *store) Remove(id string) {
	s.hub.DeleteQuery(new(kiam.Session), dbflex.Eq("SessionID", id))
}

func (s *store) Close() {
}
