package repositories

import (
	"context"

	"POS/internal/models"

	"gorm.io/gorm"
)

type StockFilter struct {
	ProductID *uint
	Type      string
	Page      int
	Limit     int
}

type StockRepository interface {
	Create(ctx context.Context, mutation *models.StockMutation) error
	FindAll(ctx context.Context, filter StockFilter) ([]*models.StockMutation, int64, error)
	FindByID(ctx context.Context, id uint) (*models.StockMutation, error)
	FindByProductID(ctx context.Context, productID uint) ([]*models.StockMutation, error)
}

type stockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) StockRepository {
	return &stockRepository{db: db}
}

func (r *stockRepository) Create(ctx context.Context, mutation *models.StockMutation) error {
	return r.db.WithContext(ctx).Create(mutation).Error
}

func (r *stockRepository) FindAll(ctx context.Context, filter StockFilter) ([]*models.StockMutation, int64, error) {
	var mutations []*models.StockMutation
	var total int64

	query := r.db.WithContext(ctx).Model(&models.StockMutation{})

	if filter.ProductID != nil {
		query = query.Where("product_id = ?", *filter.ProductID)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit

	if err := query.
		Preload("Product").
		Order("id DESC").
		Offset(offset).
		Limit(filter.Limit).
		Find(&mutations).Error; err != nil {
		return nil, 0, err
	}

	return mutations, total, nil
}

func (r *stockRepository) FindByID(ctx context.Context, id uint) (*models.StockMutation, error) {
	var mutation models.StockMutation
	if err := r.db.WithContext(ctx).
		Preload("Product").
		First(&mutation, id).Error; err != nil {
		return nil, err
	}
	return &mutation, nil
}

func (r *stockRepository) FindByProductID(ctx context.Context, productID uint) ([]*models.StockMutation, error) {
	var mutations []*models.StockMutation
	if err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("id DESC").
		Find(&mutations).Error; err != nil {
		return nil, err
	}
	return mutations, nil
}
