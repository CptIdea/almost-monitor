package api

import (
	"log"
	"net/http"
	"strings"
)

func ListenAndServe(port string) chan error {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	log.Printf("запуск http сервера на порту %q", port)

	cErr := make(chan error)
	go func() {
		err := http.ListenAndServe(port, nil)
		cErr <- err
	}()

	return nil

}
