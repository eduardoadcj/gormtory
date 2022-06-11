package gormtory

import (
	"errors"
	"gormtory/data"

	"gorm.io/gorm"
)

type GormRepository struct {
	connector        data.GormConnector
	connectionLayers []*gorm.DB
}

func NewGormRepository(connector data.GormConnector) *GormRepository {
	return &GormRepository{connector: connector}
}

func (r *GormRepository) inferConnection() (*gorm.DB, error) {
	if r.connectionLayers == nil {
		db, err := r.connector.GetConnection()
		if err != nil {
			return nil, err
		}

		r.connectionLayers = []*gorm.DB{db}
		return r.connectionLayers[0], nil
	}

	return r.connectionLayers[len(r.connectionLayers)-1], nil
}

func (r *GormRepository) Close() {
	r.connectionLayers = nil
}

// region Transaction Control implementation

func (r *GormRepository) Begin() error {
	db, err := r.inferConnection()
	if err != nil {
		return err
	}

	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	r.connectionLayers = append(r.connectionLayers, tx)
	return nil
}

func (r *GormRepository) Commit() error {
	if len(r.connectionLayers) < 2 {
		return errors.New("no transaction available")
	}

	tx, err := r.inferConnection()
	if err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	r.connectionLayers = r.connectionLayers[:len(r.connectionLayers)-1]
	return nil
}

func (r *GormRepository) Rollback() error {
	if len(r.connectionLayers) < 2 {
		return errors.New("no transaction available")
	}

	tx, err := r.inferConnection()
	if err != nil {
		return err
	}

	if err := tx.Rollback().Error; err != nil {
		return err
	}

	r.connectionLayers = r.connectionLayers[:len(r.connectionLayers)-1]
	return err
}

// endregion

func (r *GormRepository) Create(value interface{}) error {
	db, err := r.inferConnection()
	if err != nil {
		return err
	}

	if err := db.Create(value).Error; err != nil {
		return err
	}

	return nil
}
