package ledisstorage

import (
	"context"
	"log"
	"os"
	"sync"
	"task/api"
	"testing"

	lediscfg "github.com/ledisdb/ledisdb/config"
)

func TestRun(t *testing.T) {

	cfg := lediscfg.NewConfigDefault()
	os.RemoveAll(cfg.DataDir)
	defer os.RemoveAll(cfg.DataDir)

	stor := New()
	mainParams := &api.MainParams{
		Ctx: context.TODO(),
		Wg:  new(sync.WaitGroup),
	}

	err := stor.Run(mainParams)
	if err != nil {
		log.Println(err)
	}

	err = stor.Put([]byte("key1"), []byte("value1"))
	err = stor.Put([]byte("key2"), []byte("value2"))

	val1, exist1, err1 := stor.Get([]byte("key1"))
	val2, exist2, err2 := stor.Get([]byte("key2"))

	log.Println(string(val1), exist1, err1)
	log.Println(string(val2), exist2, err2)

	t.Error("MOCK")
}
