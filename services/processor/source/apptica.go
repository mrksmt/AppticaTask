package source

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"task/api"

	"github.com/pkg/errors"
)

type AppticaDataSource struct{}

var _ api.DataSource = (*AppticaDataSource)(nil)

func New() *AppticaDataSource {
	return &AppticaDataSource{}
}

func (s *AppticaDataSource) Get(from, to string) (*api.RawData, error) {

	getString := fmt.Sprintf(`https://api.apptica.com/package/top_history/1421444/1?date_from=%s&date_to=%s&B4NKGg=fVN5Q9KVOlOHDx9mOsKPAQsFBlEhBOwguLkNEDTZvKzJzT3l`, from, to)
	// fmt.Println(getString)
	resp, err := http.Get(getString)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Resp status ont OK: %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Read resp body err")
	}

	// log.Println()

	rawData := new(api.RawData)
	err = json.Unmarshal(bodyBytes, rawData)
	if err != nil {
		return nil, errors.Wrap(err, "Unmarshal resp body err")
	}

	return rawData, nil
}
