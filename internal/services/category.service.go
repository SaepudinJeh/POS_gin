package services

import (
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

func (s *CategoryService) GetAll() ([]*models.Category, error) {
	return s.repo.FindAll()
}

func (s *CategoryService) GetByID(id uint) (*models.Category, error) {
	return s.repo.FindByID(id)
}

func (s *CategoryService) Create(name, description string) (*models.Category, error) {
	// Cek nama duplikat
	if _, err := s.repo.FindByName(name); err == nil {
		return nil, errors.New("nama kategori sudah ada")
	}

	category := models.Category{
		Name:        name,
		Description: description,
	}

	if err := s.repo.Create(&category); err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *CategoryService) Update(id uint, name, description string) (*models.Category, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("kategori tidak ditemukan")
	}

	// Cek duplikat nama, kecuali dirinya sendiri
	if existing, err := s.repo.FindByName(name); err == nil && existing.ID != id {
		return nil, errors.New("nama kategori sudah ada")
	}

	category.Name = name
	category.Description = description

	if err := s.repo.Update(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) Delete(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return errors.New("kategori tidak ditemukan")
	}
	return s.repo.Delete(id)
}
