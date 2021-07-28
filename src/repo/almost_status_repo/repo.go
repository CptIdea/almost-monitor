package almost_status_repo

import (
	"almost-monitor/pkg"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type almostStatusRepo struct {
	db *gorm.DB
}

func NewAlmostStatusRepo(db *gorm.DB) (AlmostStatusRepo, error) {
	err := db.AutoMigrate(&pkg.AlmostStatus{})
	if err != nil {
		return nil, fmt.Errorf("ошибка автомиграции: %w", err)
	}
	return &almostStatusRepo{db: db}, nil
}

func (a *almostStatusRepo) Create(status *pkg.AlmostStatus) (*pkg.AlmostStatus, error) {
	err := a.db.Create(status).Error
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (a *almostStatusRepo) Update(status *pkg.AlmostStatus) (*pkg.AlmostStatus, error) {
	err := a.db.Save(status).Error
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (a *almostStatusRepo) Delete(id uint) error {
	err := a.db.Where("id = ?", id).Delete(&pkg.AlmostStatus{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *almostStatusRepo) Get(id uint) (*pkg.AlmostStatus, error) {
	var status = new(pkg.AlmostStatus)
	err := a.db.Where("id = ?", id).First(status).Error
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (a *almostStatusRepo) GetListFrom(time time.Time) ([]*pkg.AlmostStatus, error) {
	var status = make([]*pkg.AlmostStatus, 0)
	err := a.db.Where("created_at > ?", time).Find(status).Error
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (a *almostStatusRepo) GetLast() (*pkg.AlmostStatus, error) {
	var status = new(pkg.AlmostStatus)
	err := a.db.Last(status).Error
	if err != nil {
		return nil, err
	}
	return status, nil
}
