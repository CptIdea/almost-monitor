package pkg

import "gorm.io/gorm"

type AlmostStatus struct {
	gorm.Model
	Users []int
}
