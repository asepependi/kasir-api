package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
	"strconv"
)

type ProductService struct {
	productRepo *repositories.ProductRepository
}

func NewProductService(productRepo *repositories.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) GetAll(name string) ([]models.Product, error) {
	return s.productRepo.GetAll(name)
}

func (s *ProductService) Create(data *models.Product) (*models.Product, error) {
	if err := s.productRepo.Create(data); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	product, err := s.productRepo.GetByID(strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) Update(id int, data *models.Product) (*models.Product, error) {
	if err := s.productRepo.Update(strconv.Itoa(id), data); err != nil {
		return nil, err
	}
	return s.GetByID(id)
}

func (s *ProductService) Delete(id int) error {
	return s.productRepo.Delete(strconv.Itoa(id))
}