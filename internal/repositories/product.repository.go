package repositories

import (
	"context"

	"POS/internal/models"

	"gorm.io/gorm"
)

type ProductFilter struct {
	CategoryID *uint
	Search     string
	IsActive   *bool
	Page       int
	Limit      int
}

type ProductRepository interface {
	FindAll(ctx context.Context, filter ProductFilter) ([]*models.Product, int64, error)
	FindByID(ctx context.Context, id uint) (*models.Product, error)
	FindBySKU(ctx context.Context, sku string) (*models.Product, error)
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id uint) error
	CountByCategoryID(ctx context.Context, categoryID uint) (int64, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll(ctx context.Context, filter ProductFilter) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Product{})

	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR sku ILIKE ?", search, search)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
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
		Preload("Category").
		Order("id DESC").
		Offset(offset).
		Limit(filter.Limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) FindByID(ctx context.Context, id uint) (*models.Product, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).Preload("Category").First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindBySKU(ctx context.Context, sku string) (*models.Product, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}

func (r *productRepository) CountByCategoryID(ctx context.Context, categoryID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Product{}).
		Where("category_id = ?", categoryID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
