package jsonstore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ariefdarmawan/kiam"
)

type store struct {
	folderPath string
}

func NewStorage(loc string) *store {
	s := new(store)
	s.folderPath = loc
	return s
}

func (s *store) Load(pool *kiam.SessionPool) error {
	fis, e := ioutil.ReadDir(s.folderPath)
	if e != nil {
		return e
	}

	for _, fi := range fis {
		if strings.HasPrefix(fi.Name(), ".") {
			continue
		}

		bs, e := ioutil.ReadFile(filepath.Join(s.folderPath, fi.Name()))
		if e != nil {
			return fmt.Errorf("fail read file %s. %s", fi.Name(), e.Error())
		}

		sess := new(kiam.Session)
		if e = json.Unmarshal(bs, sess); e != nil {
			return fmt.Errorf("fail serializing file %s. %s", fi.Name(), e.Error())
		}

		if _, e = pool.Create(sess.ReferenceID, sess.Data, int(sess.Duration)/int(time.Second)); e != nil {
			//return fmt.Errorf("fail to create session. %s", e.Error())
			os.Remove(filepath.Join(s.folderPath, fi.Name()))
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
	locPath := filepath.Join(s.folderPath, id+".json")
	bs, e := ioutil.ReadFile(locPath)
	if e != nil {
		return nil, fmt.Errorf("fail read file %s. %s", id, e.Error())
	}

	sess := new(kiam.Session)
	if e = json.Unmarshal(bs, sess); e != nil {
		return nil, fmt.Errorf("fail serializing file %s. %s", id, e.Error())
	}

	return sess, nil
}

func (s *store) Write(sess *kiam.Session) error {
	bs, e := json.Marshal(sess)
	if e != nil {
		return fmt.Errorf("fail serializing session %s. %s", sess.SessionID, e.Error())
	}

	e = ioutil.WriteFile(filepath.Join(s.folderPath, sess.SessionID+".json"), bs, 0644)
	if e != nil {
		return fmt.Errorf("fail write file %s. %s", sess.SessionID, e.Error())
	}
	return nil
}

func (s *store) Remove(id string) {
	locPath := filepath.Join(s.folderPath, id+".json")
	os.Remove(locPath)
}

func (s *store) Close() {
}
