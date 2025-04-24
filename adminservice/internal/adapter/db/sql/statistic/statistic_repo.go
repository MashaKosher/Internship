package statistic

import "gorm.io/gorm"

type StatisticRepo struct {
	*gorm.DB
}

func New(db *gorm.DB) *StatisticRepo {
	return &StatisticRepo{db}
}
