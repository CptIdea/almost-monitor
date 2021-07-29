package status_notificator

import (
	"almost-monitor/src/nameCache"
	"almost-monitor/src/repo/almost_status_repo"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"log"
	"time"
)

type StatusNotificator struct {
	vk        *api.VK
	repo      almost_status_repo.AlmostStatusRepo
	nameCache *nameCache.NameCache
}

func NewStatusNotificator(vk *api.VK, repo almost_status_repo.AlmostStatusRepo, nameCache *nameCache.NameCache) *StatusNotificator {
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
	ticker := time.NewTicker(24 * time.Hour)
	tickerChan := ticker.C
	errCounter := 0
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
			message += fmt.Sprintf("\n%s: %d%%", s.nameCache.GetUserName(userID), (count*100)/288)
		}

		req := params.NewMessagesSendBuilder().PeerID(chatID).Message(message)
		_, err = s.vk.MessagesSend(req.Params)
		if err != nil {
			errCounter++
			log.Printf("ошибка получения статусов: %s", err)
			continue
		}

		<-tickerChan
	}
}
