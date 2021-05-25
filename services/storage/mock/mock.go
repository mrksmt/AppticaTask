package mock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"task/api"
)

type EmptyStorage struct{}

var _ api.Storage = (*EmptyStorage)(nil)

func (s *EmptyStorage) Put(key, data []byte) error           { return nil }
func (s *EmptyStorage) Get(key []byte) ([]byte, bool, error) { return nil, false, nil }

type ErrorStorage struct{ *EmptyStorage }

var _ api.Storage = (*ErrorStorage)(nil)

func (s *ErrorStorage) Get(key []byte) ([]byte, bool, error) { return nil, false, fmt.Errorf("MOCK") }

type AnyDateStorage struct{ *EmptyStorage }

var _ api.Storage = (*AnyDateStorage)(nil)

func (s *AnyDateStorage) Get(key []byte) ([]byte, bool, error) { return getMockData(), true, nil }

func getMockData() []byte {
	d := api.ResultData{
		StatusCode: http.StatusOK,
		Message:    "OK",
		Data:       map[int]api.Position{1: 2, 3: 4, 5: 6},
	}
	res, _ := json.Marshal(d)
	return res
}
