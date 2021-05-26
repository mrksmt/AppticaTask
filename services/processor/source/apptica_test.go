package source

import (
	"log"
	"testing"
	"time"
)

func TestGet(t *testing.T) {

	applicationId := "1421444"
	countryId := "1"
	today := time.Now().Format("2006-01-02")
	monthAgo := time.Now().AddDate(0, -1, 0).Format("2006-01-02")

	log.Println(applicationId, countryId, monthAgo, today)

	src := New()
	data, err := src.Get(applicationId, countryId, monthAgo, today)
	log.Println(data, err)

	t.Error("MOCK")
}
