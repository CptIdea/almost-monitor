package status_notificator

import (
	"almost-monitor/internal/name_cache"
	"almost-monitor/internal/repo/almost_status_repo"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"log"
	"math/rand"
	"time"
)

type StatusNotificator struct {
	vk        *api.VK
	repo      almost_status_repo.AlmostStatusRepo
	nameCache *name_cache.NameCache
}

func NewStatusNotificator(vk *api.VK, repo almost_status_repo.AlmostStatusRepo, nameCache *name_cache.NameCache) *StatusNotificator {
	return &StatusNotificator{vk: vk, repo: repo, nameCache: nameCache}
}

func (s *StatusNotificator) Start(groupID int, chatID int) chan error {
	cErr := make(chan error)
	go func() {
		cErr <- s.asyncStart(groupID, chatID)
	}()
	return cErr
}

func (s *StatusNotificator) asyncStart(groupID int, chatID int) error {
	rand.Seed(time.Now().Unix())
	ticker := time.NewTicker(24 * time.Hour)
	tickerChan := ticker.C
	errCounter := 0
	<-tickerChan
	for {
		if errCounter > 10 {
			return fmt.Errorf("слишком много ошибок. выход")
		}

		list, err := s.repo.GetListFrom(time.Now().AddDate(0, 0, -1))
		if err != nil {
			errCounter++
			log.Printf("ошибка получения статусов: %s", err)
			continue
		}
		counters := make(map[int64]int)
		for _, status := range list {
			for _, user := range status.Users {
				if _, ok := counters[user]; ok {
					counters[user]++
				} else {
					counters[user] = 1
				}
			}
		}
		message := "Аптайм пользователей за последние 24 часа"
		for userID, count := range counters {
			message += fmt.Sprintf("\n👉%s:  %s", s.nameCache.GetUserName(userID), time.Unix(int64(count*5*60), 0).Format("15:04"))
		}

		req := params.NewMessagesSendBuilder().PeerID(chatID).Message(message)
		req.RandomID(rand.Int())
		_, err = s.vk.MessagesSend(req.Params)
		if err != nil {
			errCounter++
			log.Printf("ошибка получения статусов: %s", err)
			continue
		}
		<-tickerChan
	}
}
