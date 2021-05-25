package manager

import (
	"encoding/json"
	"log"
	dataProcessor "task/services/processor/processor"
	"task/services/processor/source"
	"testing"
	"time"
)

func TestAAAA(t *testing.T) {

	src := source.New()
	_ = src

	today := time.Now().Format("2006-01-02")
	monthAgo := time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	log.Println(monthAgo, today)
	raw, err := src.Get(monthAgo, today)
	if err != nil {
		t.Error(err)
	}

	pretty, _ := json.MarshalIndent(raw, "", "   ")
	log.Println(string(pretty))

	dp := dataProcessor.New()

	processedData, _ := dp.ProcessRawData(raw)

	for date, dayData := range processedData {
		_ = date
		_ = dayData
	}

}
