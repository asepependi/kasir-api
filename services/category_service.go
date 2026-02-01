package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
	"strconv"
)

type CategoryService struct {
	repoCategory *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repoCategory: repo}
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.repoCategory.GetAll()
}

func (s *CategoryService) Create(data *models.Category) (*models.Category, error) {
	if err := s.repoCategory.Create(data); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	category, err := s.repoCategory.GetByID(strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) Update(id int, data *models.Category) (*models.Category, error) {
	if err := s.repoCategory.Update(strconv.Itoa(id), data); err != nil {
		return nil, err
	}
	return s.GetByID(id)
}

func (s *CategoryService) Delete(id int) error {
	return s.repoCategory.Delete(strconv.Itoa(id))
}
