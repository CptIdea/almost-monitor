package main

import (
	"almost-monitor/internal/api"
	nameCache2 "almost-monitor/internal/name_cache"
	"almost-monitor/internal/repo/almost_status_repo"
	"almost-monitor/internal/status_notificator"
	"almost-monitor/internal/status_scanner"
	vkApi "github.com/SevereCloud/vksdk/v2/api"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
)

func main() {
	token := os.Getenv("TOKEN")
	groupID, err := strconv.Atoi(os.Getenv("GROUP_ID"))
	if err != nil {
		log.Fatalf("ошибка конвертации id группы в число: %s", err)
	}
	chatID, err := strconv.Atoi(os.Getenv("CHAT_ID"))
	if err != nil {
		log.Fatalf("ошибка конвертации id чата в число: %s", err)
	}
	port := os.Getenv("PORT")

	vk := vkApi.NewVK(token)

	dsn := "host=localhost user=mvp password=mvp dbname=almost_status port=5432"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("ошибка подключения базы данных: %s", err)
	}

	repo, err := almost_status_repo.NewAlmostStatusRepo(db)
	if err != nil {
		log.Fatalf("ошибка создания репозитория: %s", err)
	}

	scanner := status_scanner.NewScanner(vk, repo)

	nameCache := nameCache2.NewNameCache(vk)
	notificator := status_notificator.NewStatusNotificator(vk, repo, nameCache)

	httpServer := api.NewHttpServer(repo, nameCache)

	select {
	case err := <-api.ListenAndServe(port, httpServer):
		log.Fatalf("выход http сервера с ошибкой: %s", err)

	case err := <-scanner.Start(groupID, chatID):
		log.Fatalf("выход сканера с ошибкой: %s", err)

	case err := <-notificator.Start(groupID, chatID):
		log.Fatalf("выход нотификатора с ошибкой: %s", err)
	}
}
