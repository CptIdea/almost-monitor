package api

import (
	"log"
	"net/http"
	"strings"
)

func ListenAndServe(port string, server *HttpServer) chan error {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	http.HandleFunc("/getFromTime", server.GetFromTime)

	log.Printf("запуск http сервера на порту %q", port)

	cErr := make(chan error)
	go func() {
		err := http.ListenAndServe(port, nil)
		cErr <- err
	}()

	return cErr

}
