package service

import (
	"github.com/yourname/products/internal/broker"
	"github.com/yourname/products/internal/metrics"
	"github.com/yourname/products/internal/models"
	"github.com/yourname/products/internal/repository"
)

type ProductService interface {
	Create(req *models.CreateProductRequest) (*models.Product, error)
	Delete(id int) error
	List(page, pageSize int) (*models.ProductListResponse, error)
}

type productService struct {
	repo   repository.ProductRepository
	broker broker.MessageBroker
}

func NewProductService(repo repository.ProductRepository, broker broker.MessageBroker) ProductService {
	return &productService{
		repo:   repo,
		broker: broker,
	}
}

func (s *productService) Create(req *models.CreateProductRequest) (*models.Product, error) {
	product, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	if err := s.broker.PublishProductCreated(product.ID, product.Name); err != nil {
	}

	metrics.ProductsCreated.Inc()

	return product, nil
}

func (s *productService) Delete(id int) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	if err := s.broker.PublishProductDeleted(id); err != nil {
	}

	metrics.ProductsDeleted.Inc()

	return nil
}

func (s *productService) List(page, pageSize int) (*models.ProductListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	products, total, err := s.repo.List(page, pageSize)
	if err != nil {
		return nil, err
	}

	return &models.ProductListResponse{
		Products: products,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
