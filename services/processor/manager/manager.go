package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"task/api"
	"task/services/processor/processor"
	"time"

	"github.com/pkg/errors"
)

type Manager struct {
	cancel context.CancelFunc
}

var _ api.Runnable = (*Manager)(nil)
var cfg *Config
var dataSource api.DataSource
var dp api.DataProcessor
var stor api.Storage

func New(src api.DataSource, dataProc api.DataProcessor, storage api.Storage, config *Config) *Manager {
	dataSource = src
	dp = dataProc
	stor = storage
	checkConfig(config)
	return &Manager{}
}

func (s *Manager) Run(mainParams *api.MainParams) error {

	var localCtx context.Context
	localCtx, s.cancel = context.WithCancel(mainParams.Ctx)

	mainParams.Wg.Add(1)
	go s.mainLoop(localCtx, mainParams.Wg)

	dp = processor.New()

	return nil
}

func (s *Manager) mainLoop(ctx context.Context, wg *sync.WaitGroup) {

	defer wg.Done()
	ctxDone := ctx.Done()
	ticker := time.NewTicker(time.Second * time.Duration(cfg.Rate))

	for {

		s.updateData()

		select {
		case <-ctxDone:
			return
		case <-ticker.C:
		}
	}
}

func (s *Manager) updateData() {

	// get raw
	today := time.Now().Format("2006-01-02")
	monthAgo := time.Now().AddDate(0, -1, 0).Format("2006-01-02")

	log.Printf("Get raw data from %s to %s", monthAgo, today)

	raw, err := dataSource.Get(monthAgo, today)
	if err != nil {
		log.Println(errors.Wrap(err, "Get raw data err"))
		return
	}

	// process
	processedData, err := dp.ProcessRawData(raw)
	if err != nil {
		log.Println(errors.Wrap(err, "Process raw data err"))
	}

	// save result
	for date, dateData := range processedData {

		// convert to JSON
		preparedJSON, err := json.Marshal(dateData)
		if err != nil {
			log.Println(errors.Wrap(err, "Marshaling result data err"))
		}

		// save to storage
		err = stor.Put([]byte(date), preparedJSON)
		if err != nil {
			log.Println(errors.Wrap(err, "Saving result data err"))
		}
		log.Printf("Update data for %s", date)

	}
	fmt.Println()
}
