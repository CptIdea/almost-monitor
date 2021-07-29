package nameCache

import (
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
)

type NameCache struct {
	vk    *api.VK
	cache map[int64]string
}

func NewNameCache(vk *api.VK) *NameCache {
	return &NameCache{vk: vk, cache: make(map[int64]string)}
}

func (n *NameCache) GetUserName(ID int64) string {
	if ans, ok := n.cache[ID]; ok {
		return ans
	}

	users, err := n.vk.UsersGet(params.NewUsersGetBuilder().UserIDs([]string{fmt.Sprint(ID)}).Params)
	if err != nil {
		return "UNKNOWN"
	}
	if len(users) > 0 {
		n.cache[ID] = fmt.Sprint(users[0].FirstName, users[0].LastName)
		return fmt.Sprintf("%s %s", users[0].FirstName, users[0].LastName)
	}

	return ""
}
