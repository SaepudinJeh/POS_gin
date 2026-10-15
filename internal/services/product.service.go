package services

import (
	"context"
	"errors"

	"POS/internal/models"
	"POS/internal/repositories"
)

type ProductInput struct {
	SKU         string
	Name        string
	Description string
	Price       float64
	Cost        float64
	Stock       int
	CategoryID  uint
	IsActive    *bool
}

type ProductService struct {
	productRepo  repositories.ProductRepository
	categoryRepo repositories.CategoryRepository
}

func NewProductService(
	productRepo repositories.ProductRepository,
	categoryRepo repositories.CategoryRepository,
) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *ProductService) GetAll(ctx context.Context, filter repositories.ProductFilter) ([]*models.Product, int64, error) {
	return s.productRepo.FindAll(ctx, filter)
}

func (s *ProductService) GetByID(ctx context.Context, id uint) (*models.Product, error) {
	return s.productRepo.FindByID(ctx, id)
}

func (s *ProductService) Create(ctx context.Context, input ProductInput) (*models.Product, error) {
	if _, err := s.productRepo.FindBySKU(ctx, input.SKU); err == nil {
		return nil, errors.New("SKU sudah dipakai")
	}

	if _, err := s.categoryRepo.FindByID(ctx, input.CategoryID); err != nil {
		return nil, errors.New("kategori tidak ditemukan")
	}

	product := models.Product{
		SKU:         input.SKU,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Cost:        input.Cost,
		Stock:       input.Stock,
		CategoryID:  input.CategoryID,
		IsActive:    true,
	}

	if err := s.productRepo.Create(ctx, &product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *ProductService) Update(ctx context.Context, id uint, input ProductInput) (*models.Product, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	if existing, err := s.productRepo.FindBySKU(ctx, input.SKU); err == nil && existing.ID != id {
		return nil, errors.New("SKU sudah dipakai")
	}

	if _, err := s.categoryRepo.FindByID(ctx, input.CategoryID); err != nil {
		return nil, errors.New("kategori tidak ditemukan")
	}

	product.SKU = input.SKU
	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price
	product.Cost = input.Cost
	product.Stock = input.Stock
	product.CategoryID = input.CategoryID
	if input.IsActive != nil {
		product.IsActive = *input.IsActive
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) Delete(ctx context.Context, id uint) error {
	if _, err := s.productRepo.FindByID(ctx, id); err != nil {
		return errors.New("produk tidak ditemukan")
	}
	return s.productRepo.Delete(ctx, id)
}
