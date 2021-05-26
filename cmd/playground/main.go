package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"task/api"

	"github.com/gorilla/mux"
)

func main() {

	log.SetFlags(log.Lshortfile)

	resp, err := http.Get("https://api.apptica.com/package/top_history/1421444/1?date_from=2021-04-22&date_to=2021-05-23&B4NKGg=fVN5Q9KVOlOHDx9mOsKPAQsFBlEhBOwguLkNEDTZvKzJzT3l")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	r := mux.NewRouter()
	r.Host("{subdomain}.apptica.com").
		Path("/package/{applicationId:[0-9]+}/{countryId:[0-9]+}").
		Queries("date_from", "{date_from:[0-9]{4}-[0-9]{2}-[0-9]{2}}").
		Queries("date_to", "{date_to:[0-9]{4}-[0-9]{2}-[0-9]{2}}").
		Queries("B4NKGg", "{B4NKGg}").
		Name("rawdata")

	// url.String() will be "http://news.example.com/articles/technology/42?filter=gorilla"
	url, err := r.Get("rawdata").URL(
		"subdomain", "api",
		"subdomain", "api",
		"applicationId", "1421444",
		"countryId", "1",
		"date_from", "2021-04-22",
		"date_to", "2021-05-23",
		"B4NKGg", "fVN5Q9KVOlOHDx9mOsKPAQsFBlEhBOwguLkNEDTZvKzJzT3l",
	)

	log.Println(url.String())
	return

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		// bodyString := string(bodyBytes)
		// log.Println(bodyString)

		aa := new(api.RawData)
		json.Unmarshal(bodyBytes, aa)

		log.Printf("%#v", aa)
		bb, _ := json.MarshalIndent(aa, "", "   ")
		fmt.Println(string(bb))
	}
}

/*
https://api.apptica.com/package/top_history/{{applicationId}}/{{countryId}}?date_from={{dateFrom}}&date_to={{dateTo}}&B4NKGg=fVN5Q9KVOlOHDx9mOsKPAQsFBlEhBOwguLkNEDTZvKzJzT3l
https://api.apptica.com/package/top_history/1421444/1?date_from=2021-05-1&date_to=2021-05-23&B4NKGg=fVN5Q9KVOlOHDx9mOsKPAQsFBlEhBOwguLkNEDTZvKzJzT3l
*/
