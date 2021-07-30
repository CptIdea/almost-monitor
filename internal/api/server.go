package api

import (
	"almost-monitor/internal/name_cache"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"almost-monitor/internal/repo/almost_status_repo"
)

type HttpServer struct {
	repo      almost_status_repo.AlmostStatusRepo
	nameCache *name_cache.NameCache
}

func NewHttpServer(repo almost_status_repo.AlmostStatusRepo, nameCache *name_cache.NameCache) *HttpServer {
	return &HttpServer{repo: repo, nameCache: nameCache}
}

func (s *HttpServer) GetFromTime(w http.ResponseWriter, r *http.Request) {
	//CORS
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	if r.Method == "OPTIONS" {
		return
	}

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

	for i := range list {
		s.nameCache.FillNames(list[i])
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
