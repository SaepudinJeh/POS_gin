package services

import (
	"context"
	"errors"

	"POS/internal/models"
	"POS/internal/repositories"
)

type StockInput struct {
	ProductID uint
	Type      string
	Quantity  int
	Note      string
	UserID    uint
}

type StockService struct {
	stockRepo   repositories.StockRepository
	productRepo repositories.ProductRepository
}

func NewStockService(
	stockRepo repositories.StockRepository,
	productRepo repositories.ProductRepository,
) *StockService {
	return &StockService{
		stockRepo:   stockRepo,
		productRepo: productRepo,
	}
}

func (s *StockService) GetAll(ctx context.Context, filter repositories.StockFilter) ([]*models.StockMutation, int64, error) {
	return s.stockRepo.FindAll(ctx, filter)
}

func (s *StockService) GetByID(ctx context.Context, id uint) (*models.StockMutation, error) {
	return s.stockRepo.FindByID(ctx, id)
}

func (s *StockService) GetByProductID(ctx context.Context, productID uint) ([]*models.StockMutation, error) {
	if _, err := s.productRepo.FindByID(ctx, productID); err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}
	return s.stockRepo.FindByProductID(ctx, productID)
}

func (s *StockService) CreateMutation(ctx context.Context, input StockInput) (*models.StockMutation, error) {
	product, err := s.productRepo.FindByID(ctx, input.ProductID)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	var newStock int
	switch input.Type {
	case models.StockTypeIn:
		newStock = product.Stock + input.Quantity
	case models.StockTypeOut:
		if product.Stock < input.Quantity {
			return nil, errors.New("stok tidak mencukupi")
		}
		newStock = product.Stock - input.Quantity
	case models.StockTypeAdjustment:
		newStock = input.Quantity
	default:
		return nil, errors.New("tipe mutasi tidak valid")
	}

	mutation := &models.StockMutation{
		ProductID:   product.ID,
		Type:        input.Type,
		Quantity:    input.Quantity,
		StockBefore: product.Stock,
		StockAfter:  newStock,
		Note:        input.Note,
		UserID:      input.UserID,
	}

	if err := s.stockRepo.Create(ctx, mutation); err != nil {
		return nil, err
	}

	product.Stock = newStock
	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return s.stockRepo.FindByID(ctx, mutation.ID)
}
