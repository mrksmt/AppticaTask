package ledisstorage

import (
	"task/api"

	lediscfg "github.com/ledisdb/ledisdb/config"
	"github.com/ledisdb/ledisdb/ledis"
	"github.com/pkg/errors"
)

type LedisStorage struct {
	db *ledis.DB
}

var _ api.Storage = (*LedisStorage)(nil)
var _ api.Runnable = (*LedisStorage)(nil)

func New() *LedisStorage {
	s := &LedisStorage{}
	return s
}

func (s *LedisStorage) Run(mainParams *api.MainParams) error {

	cfg := lediscfg.NewConfigDefault()
	cfg.DataDir = "../var"
	l, err := ledis.Open(cfg)
	if err != nil {
		return errors.Wrap(err, "LedisDB open cfg")
	}

	s.db, err = l.Select(0)
	if err != nil {
		return errors.Wrap(err, "LedisDB select err")
	}

	return nil
}

func (s *LedisStorage) Put(key, value []byte) error {
	return s.db.Set(key, value)
}

func (s *LedisStorage) Get(key []byte) ([]byte, bool, error) {

	exist, err := s.db.Exists(key)
	if err != nil {
		return []byte{}, false, errors.Wrap(err, "LedisDB Key exist err")
	}
	if exist == 0 {
		return []byte{}, false, nil
	}

	val, err := s.db.Get(key)
	if err != nil {
		return []byte{}, false, errors.Wrap(err, "LedisDB Get err")
	}
	return val, true, nil
}
