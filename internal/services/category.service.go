package services

import (
	"context"
	"errors"

	"POS/internal/models"
	"POS/internal/repositories"
)

type CategoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll(ctx context.Context) ([]*models.Category, error) {
	return s.repo.FindAll(ctx)
}

func (s *CategoryService) GetByID(ctx context.Context, id uint) (*models.Category, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *CategoryService) Create(ctx context.Context, name, description string) (*models.Category, error) {
	if _, err := s.repo.FindByName(ctx, name); err == nil {
		return nil, errors.New("nama kategori sudah ada")
	}

	category := models.Category{
		Name:        name,
		Description: description,
	}

	if err := s.repo.Create(ctx, &category); err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *CategoryService) Update(ctx context.Context, id uint, name, description string) (*models.Category, error) {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("kategori tidak ditemukan")
	}

	if existing, err := s.repo.FindByName(ctx, name); err == nil && existing.ID != id {
		return nil, errors.New("nama kategori sudah ada")
	}

	category.Name = name
	category.Description = description

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) Delete(ctx context.Context, id uint) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return errors.New("kategori tidak ditemukan")
	}
	return s.repo.Delete(ctx, id)
}
