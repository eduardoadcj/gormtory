package data

import "gorm.io/gorm"

type GormConnector interface {
	GetConnection() (*gorm.DB, error)
}
