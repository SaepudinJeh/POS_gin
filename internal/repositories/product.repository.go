package repositories

import (
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
	FindAll(filter ProductFilter) ([]*models.Product, int64, error)
	FindByID(id uint) (*models.Product, error)
	FindBySKU(sku string) (*models.Product, error)
	Create(product *models.Product) error
	Update(product *models.Product) error
	Delete(id uint) error
	CountByCategoryID(categoryID uint) (int64, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll(filter ProductFilter) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	query := r.db.Model(&models.Product{})

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

func (r *productRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	if err := r.db.Preload("Category").First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindBySKU(sku string) (*models.Product, error) {
	var product models.Product
	if err := r.db.Where("sku = ?", sku).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *productRepository) CountByCategoryID(categoryID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Product{}).
		Where("category_id = ?", categoryID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
