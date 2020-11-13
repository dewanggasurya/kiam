package kiam

import (
	"errors"

	"git.kanosolution.net/kano/kaos"
	"github.com/eaciit/toolkit"
)

type iam struct {
	pool *SessionPool
	st   IAMStorage
}

func NewIAM(logger *toolkit.LogEngine) *iam {
	ae := new(iam)
	ae.pool = NewSessionPool(logger)
	return ae
}

func NewIAMWithStorage(logger *toolkit.LogEngine, st IAMStorage) *iam {
	ae := new(iam)
	ae.pool = NewSessionPool(logger)
	ae.SetStorage(st)
	return ae
}

func (a *iam) SetStorage(st IAMStorage) *iam {
	a.st = st
	return a
}

func (a *iam) Get(ctx *kaos.Context, parm toolkit.M) (*Session, error) {
	id := parm.GetString("id")
	if id == "" {
		return nil, errors.New("ID is mandatory")
	}
	session, ok := a.pool.GetBySessionID(id)
	if !ok {
		if a.st == nil {
			return nil, errors.New("Session not found")
		}
		return a.st.Get(a.pool, id)
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
	if e == nil && a.st != nil {
		go a.st.Write(a.pool, s)
	}
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
	if a.st != nil {
		go a.st.Write(a.pool, se)
	}
	return se, nil
}

func (a *iam) Store() error {
	if a.st != nil {
		return a.st.Store(a.pool)
	}
	return nil
}

func (a *iam) Load() error {
	if a.st != nil {
		return a.st.Load(a.pool)
	}
	return nil
}

func (a *iam) Close() {
	// do nothing for now
}
