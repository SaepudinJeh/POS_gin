package repositories

import (
	"context"

	"POS/internal/models"

	"gorm.io/gorm"
)

type SettingRepository interface {
	FindAll(ctx context.Context) ([]*models.Setting, error)
	FindByKey(ctx context.Context, key string) (*models.Setting, error)
	Upsert(ctx context.Context, key, value string) error
	BulkUpsert(ctx context.Context, settings map[string]string) error
}

type settingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) SettingRepository {
	return &settingRepository{db: db}
}

func (r *settingRepository) FindAll(ctx context.Context) ([]*models.Setting, error) {
	var settings []*models.Setting
	if err := r.db.WithContext(ctx).Order("key ASC").Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

func (r *settingRepository) FindByKey(ctx context.Context, key string) (*models.Setting, error) {
	var s models.Setting
	if err := r.db.WithContext(ctx).Where("key = ?", key).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *settingRepository) Upsert(ctx context.Context, key, value string) error {
	var s models.Setting
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&s).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.WithContext(ctx).Create(&models.Setting{Key: key, Value: value}).Error
	}
	if err != nil {
		return err
	}
	s.Value = value
	return r.db.WithContext(ctx).Save(&s).Error
}

func (r *settingRepository) BulkUpsert(ctx context.Context, settings map[string]string) error {
	return r.db.WithContext(ctx).Transaction(func(dbtx *gorm.DB) error {
		for key, value := range settings {
			var s models.Setting
			err := dbtx.Where("key = ?", key).First(&s).Error
			if err == gorm.ErrRecordNotFound {
				if err := dbtx.Create(&models.Setting{Key: key, Value: value}).Error; err != nil {
					return err
				}
				continue
			}
			if err != nil {
				return err
			}
			s.Value = value
			if err := dbtx.Save(&s).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
