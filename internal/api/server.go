package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"almost-monitor/internal/repo/almost_status_repo"
)

type HttpServer struct {
	repo almost_status_repo.AlmostStatusRepo
}

func NewHttpServer(repo almost_status_repo.AlmostStatusRepo) *HttpServer {
	return &HttpServer{repo: repo}
}

func (s *HttpServer) GetFromTime(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("ошибка парсинга формы: %s", err)
		s.SendError(w, err)
		return
	}
	t, err := strconv.Atoi(r.Form.Get("time"))

	list, err := s.repo.GetListFrom(time.Unix(int64(t), 0))
	if err != nil {
		log.Printf("ошибка получения списка: %s", err)
		s.SendError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		log.Printf("ошибка отправки ответа:%s", err)
	}
}

func (s *HttpServer) SendError(w http.ResponseWriter, err error) {
	_, err = fmt.Fprint(w, err)
	if err != nil {
		log.Printf("ошибка отправки ошибки: %s", err)
	}
}
