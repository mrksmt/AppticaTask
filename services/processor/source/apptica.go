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

func (s *AppticaDataSource) Get(applicationId, countryId, from, to string) (*api.RawData, error) {

	r := mux.NewRouter()
	r.Host("{subdomain}.apptica.com").
		Path("/package/top_history/{applicationId:[0-9]+}/{countryId:[0-9]+}").
		Queries("date_from", "{date_from:[0-9]{4}-[0-9]{2}-[0-9]{2}}").
		Queries("date_to", "{date_to:[0-9]{4}-[0-9]{2}-[0-9]{2}}").
		Queries("B4NKGg", "{B4NKGg}").
		Name("rawdata")

	url, err := r.Get("rawdata").URL(
		"subdomain", "api",
		"applicationId", applicationId,
		"countryId", countryId,
		"date_from", from,
		"date_to", to,
		"B4NKGg", "fVN5Q9KVOlOHDx9mOsKPAQsFBlEhBOwguLkNEDTZvKzJzT3l",
	)

	// http://api.apptica.com/package/1421444/1?date_from=2021-04-26&date_to=2021-05-26&B4NKGg=fVN5Q9KVOlOHDx9mOsKPAQsFBlEhBOwguLkNEDTZvKzJzT3l
	// https://api.apptica.com/package/top_history/1421444/1?date_from=2021-05-1&date_to=2021-05-23&B4NKGg=fVN5Q9KVOlOHDx9mOsKPAQsFBlEhBOwguLkNEDTZvKzJzT3l

	log.Println(url.String())
	if err != nil {
		return nil, errors.Wrap(err, "Get request URL err")
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return nil, errors.Wrap(err, "http.Get err")
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
