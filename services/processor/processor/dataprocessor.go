package processor

import (
	"net/http"
	"task/api"
)

type DataProcessorType1 struct{}

var _ api.DataProcessor = (*DataProcessorType1)(nil)

func New() *DataProcessorType1 { return &DataProcessorType1{} }

func (dp *DataProcessorType1) ProcessRawData(rd *api.RawData) (map[string]*api.ResultData, error) {

	result := make(map[string]*api.ResultData)

	for category, categoryData := range rd.Data {
		for _, subCategoryData := range categoryData {
			for date, position := range subCategoryData {

				currDateResultData, exist := result[date]
				if !exist {
					currDateResultData = &api.ResultData{
						StatusCode: http.StatusOK,
						Message:    "OK",
						Data:       make(map[int]api.Position),
					}
					result[date] = currDateResultData
				}

				categoryPosition, exist := currDateResultData.Data[category]
				if !exist || position < categoryPosition {
					currDateResultData.Data[category] = position
				}

			}
		}
	}

	return result, nil
}
