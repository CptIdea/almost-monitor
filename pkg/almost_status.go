package pkg

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type AlmostStatus struct {
	gorm.Model
	Users     pq.Int64Array `gorm:"type:int[]"`
	UsersName []string      `gorm:"-"`
}
