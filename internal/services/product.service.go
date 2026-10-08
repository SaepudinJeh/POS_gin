package services

import (
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

func (s *ProductService) GetAll(filter repositories.ProductFilter) ([]*models.Product, int64, error) {
	return s.productRepo.FindAll(filter)
}

func (s *ProductService) GetByID(id uint) (*models.Product, error) {
	return s.productRepo.FindByID(id)
}

func (s *ProductService) Create(input ProductInput) (*models.Product, error) {
	// SKU harus unik
	if _, err := s.productRepo.FindBySKU(input.SKU); err == nil {
		return nil, errors.New("SKU sudah dipakai")
	}

	// Kategori harus ada
	if _, err := s.categoryRepo.FindByID(input.CategoryID); err != nil {
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

	if err := s.productRepo.Create(&product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *ProductService) Update(id uint, input ProductInput) (*models.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	// SKU unik, kecuali dirinya sendiri
	if existing, err := s.productRepo.FindBySKU(input.SKU); err == nil && existing.ID != id {
		return nil, errors.New("SKU sudah dipakai")
	}

	// Kategori harus ada
	if _, err := s.categoryRepo.FindByID(input.CategoryID); err != nil {
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

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) Delete(id uint) error {
	if _, err := s.productRepo.FindByID(id); err != nil {
		return errors.New("produk tidak ditemukan")
	}
	return s.productRepo.Delete(id)
}
