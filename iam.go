package kiam

import (
	"errors"

	"git.kanosolution.net/kano/kaos"
	"github.com/eaciit/toolkit"
)

type iam struct {
	pool *SessionPool
}

func NewIAM(logger *toolkit.LogEngine) *iam {
	ae := new(iam)
	ae.pool = NewSessionPool(logger)
	return ae
}

func (a *iam) Get(ctx *kaos.Context, parm toolkit.M) (*Session, error) {
	id := parm.GetString("id")
	if id == "" {
		return nil, errors.New("ID is mandatory")
	}
	session, ok := a.pool.GetBySessionID(id)
	if !ok {
		return nil, errors.New("Session not found")
	}
	a.pool.Update(session.SessionID, 0)
	return session, nil
}

func (a *iam) Create(ctx *kaos.Context, parm toolkit.M) (*Session, error) {
	id := parm.GetString("id")
	duration := parm.GetInt("second")
	if id == "" {
		return nil, errors.New("ID is mandatory")
	}

	s, e := a.pool.Create(id, nil, duration)
	return s, e
}

func (a *iam) Renew(ctx *kaos.Context, parm toolkit.M) (*Session, error) {
	id := parm.GetString("id")
	duration := parm.GetInt("second")
	if id == "" {
		return nil, errors.New("ID is mandatory")
	}

	seid, e := a.pool.Renew(id, duration)
	if e != nil {
		return nil, e
	}
	se, _ := a.pool.GetBySessionID(seid)
	return se, nil
}
