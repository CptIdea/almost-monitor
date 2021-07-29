package almost_status_repo

import (
	"almost-monitor/pkg"
	"time"
)

type AlmostStatusRepo interface {
	Create(status *pkg.AlmostStatus) (*pkg.AlmostStatus, error)
	Update(status *pkg.AlmostStatus) (*pkg.AlmostStatus, error)
	Delete(id uint) error
	Get(id uint) (*pkg.AlmostStatus, error)
	GetListFrom(time time.Time) ([]*pkg.AlmostStatus, error)
	GetLast() (*pkg.AlmostStatus, error)
}
