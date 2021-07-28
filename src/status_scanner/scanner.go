package status_scanner

import (
	"almost-monitor/pkg"
	"almost-monitor/src/repo/almost_status_repo"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"log"
	"time"
)

type Scanner struct {
	vk   *api.VK
	repo almost_status_repo.AlmostStatusRepo
}

func NewScanner(vk *api.VK, repo almost_status_repo.AlmostStatusRepo) *Scanner {
	return &Scanner{vk: vk, repo: repo}
}

func (s *Scanner) Start(groupID int, chatID int) chan error {
	cErr := make(chan error)
	go func() {
		cErr <- s.asyncStart(groupID, chatID)
	}()
	return cErr
}

func (s *Scanner) asyncStart(groupID int, chatID int) error {
	log.Printf("запуск сканера активности группой %d в чате %d", groupID, chatID)
	errCounter := 0
	for {
		if errCounter > 10 {
			return fmt.Errorf("слишком много ошибок. выход")
		}

		req := params.NewMessagesGetConversationMembersBuilder().GroupID(groupID).PeerID(chatID).Fields([]string{"online"})
		response, err := s.vk.MessagesGetConversationMembers(req.Params)
		if err != nil {
			errCounter++
			log.Printf("ошибка получения пользователей беседы: %s", err)
			continue
		}

		status := new(pkg.AlmostStatus)
		for _, profile := range response.Profiles {
			if profile.Online {
				status.Users = append(status.Users, int64(profile.ID))
			}
		}
		_, err = s.repo.Create(status)
		if err != nil {
			errCounter++
			log.Printf("ошибка создания статуса пользователей: %s", err)
			continue
		}
		time.Sleep(5 * time.Minute)
	}
}
