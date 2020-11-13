package kiam

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type IAMStorage interface {
	Load(pool *SessionPool) error
	Store(pool *SessionPool) error
	Get(pool *SessionPool, id string) (*Session, error)
	Write(pool *SessionPool, sess *Session) error
	Close()
}

type jsonStorage struct {
	fp string
}

func NewJsonStorage(fp string) *jsonStorage {
	j := new(jsonStorage)
	j.fp = fp
	return j
}

func (j *jsonStorage) Load(pool *SessionPool) error {
	if _, err := os.Stat(j.fp); os.IsNotExist(err) {
		return nil
	}
	bs, err := ioutil.ReadFile(j.fp)
	if err != nil {
		return err
	}

	sess := map[string]*Session{}
	if err = json.Unmarshal(bs, &sess); err != nil {
		return err
	}
	return nil
}

func (j *jsonStorage) Store(pool *SessionPool) error {
	bs, err := json.Marshal(pool.sessions)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(j.fp, bs, 0644); err != nil {
		return err
	}
	return nil
}

func (j *jsonStorage) Get(pool *SessionPool, id string) (*Session, error) {
	return nil, nil
}

func (j *jsonStorage) Write(pool *SessionPool, s *Session) error {
	return nil
}

func (j *jsonStorage) Close() {
	//panic("not implemented") // TODO: Implement
}
