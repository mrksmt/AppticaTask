package source

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"task/api"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type AppticaDataSource struct{}

var _ api.DataSource = (*AppticaDataSource)(nil)

func New() *AppticaDataSource {
	return &AppticaDataSource{}
}

func (s *AppticaDataSource) Get(from, to string) (*api.RawData, error) {

	r := mux.NewRouter()
	r.Host("{subdomain}.apptica.com").
		Path("/package/{applicationId}/{countryId:[0-9]+}").
		Queries("date_from", "{date_from}").
		Queries("date_to", "{date_to}").
		Queries("B4NKGg", "{B4NKGg}").
		Name("rawdata")

	url, err := r.Get("rawdata").URL(
		"subdomain", "api",
		"applicationId", "1421444",
		"countryId", "1",
		"date_from", from,
		"date_to", to,
		"B4NKGg", "fVN5Q9KVOlOHDx9mOsKPAQsFBlEhBOwguLkNEDTZvKzJzT3l",
	)

	resp, err := http.Get(url.String())
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
