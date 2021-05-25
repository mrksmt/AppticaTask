package apptopcat

import (
	"encoding/json"
	"log"
	"net/http"
	"task/api"
	"time"

	"github.com/pkg/errors"
)

// Handler
type Handler struct{}

// проверка реализации типом требуемых интерфейсов
var _ api.HttpRouter = (*Handler)(nil)
var stor api.Storage

// New возвращает хэндлер
func New(storage api.Storage) *Handler {
	stor = storage
	s := &Handler{}
	return s
}

// Routes реализация интерфейса HttpRouter, возвращает набор путей и хэндлеров для них
func (s *Handler) Routes() []api.Route {
	routes := []api.Route{
		{
			Name:       "GetAppTopCategory",
			Method:     http.MethodGet,
			Pattern:    "/appTopCategory",
			Handler:    s.appTopCategory,
			QueryPairs: []string{"date", "{date}"}, // https://stackoverflow.com/a/45378656/10487940
		},
	}
	return routes
}

// appTopCategory хэндлер запроса http://<ваш домен>/appTopCategory?date=2021-04-01
func (s *Handler) appTopCategory(w http.ResponseWriter, r *http.Request) {

	date := r.FormValue("date")
	log.Println(date)

	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		err = errors.Wrap(err, "Bad request err")
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		log.Println(err)
		return
	}

	if stor == nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "No data storage")
		return
	}

	data, exist, err := stor.Get([]byte(date))

	if err != nil {
		err = errors.Wrap(err, "Storage Get err")
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !exist {
		WriteErrorResponse(w, http.StatusNotFound, "Not found")
		return
	}

	WriteJSONResponse(w, http.StatusOK, []byte(data))
}

// Запись http ответа в виде json с  набором заголовков
func WriteJSONResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(status)
	w.Write(data)
}

func WriteErrorResponse(w http.ResponseWriter, code int, message string) {
	dataMsg := &api.ResultData{
		StatusCode: code,
		Message:    message,
	}
	data, err := json.Marshal(dataMsg)
	if err != nil {
		log.Println(errors.Wrap(err, "ResultData marshaling err"))
		return
	}
	WriteJSONResponse(w, http.StatusInternalServerError, data)
}
