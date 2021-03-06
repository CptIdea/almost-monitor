package name_cache

import (
	"almost-monitor/pkg"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"strings"
	"sync"
)

type NameCache struct {
	vk    *api.VK
	cache map[int64]string
	mu    sync.Mutex
}

func NewNameCache(vk *api.VK) *NameCache {
	return &NameCache{vk: vk, cache: make(map[int64]string)}
}

func (n *NameCache) GetUserName(ID int64) string {
	n.mu.Lock()
	defer n.mu.Unlock()
	if ans, ok := n.cache[ID]; ok {
		return ans
	}

	users, err := n.vk.UsersGet(params.NewUsersGetBuilder().UserIDs([]string{fmt.Sprint(ID)}).Params)
	if err != nil {
		return "UNKNOWN"
	}
	if len(users) > 0 {
		n.cache[ID] = fmt.Sprintf("%s %s", users[0].FirstName, users[0].LastName)
		return fmt.Sprintf("%s %s", users[0].FirstName, users[0].LastName)
	}

	return ""
}

func (n *NameCache) FillNames(status *pkg.AlmostStatus) {
	for _, id := range status.Users {
		status.UsersName = append(status.UsersName, strings.Split(n.GetUserName(id), " ")[0])
	}
}
