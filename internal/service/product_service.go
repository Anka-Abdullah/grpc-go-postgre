package service

import (
	"fmt"
	"strings"

	"grpc-exmpl/internal/model"
	"grpc-exmpl/internal/repository"
)

// ProductService defines business logic for product operations
type ProductService interface {
	CreateProduct(req *model.CreateProductRequest) (*model.Product, error)
	GetProductByID(id int64) (*model.Product, error)
	ListProductsByUser(userID int64) ([]*model.Product, error)
	UpdateProduct(req *model.UpdateProductRequest) (*model.Product, error)
	DeleteProduct(id int64) error
}

type productService struct {
	repo repository.ProductRepository
}

// NewProductService creates a new instance of ProductService
func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

// CreateProduct handles product creation logic
func (s *productService) CreateProduct(req *model.CreateProductRequest) (*model.Product, error) {
	if err := s.validateCreate(req); err != nil {
		return nil, err
	}

	p := &model.Product{
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Price:       req.Price,
		Stock:       req.Stock,
		UserID:      req.UserID,
	}

	if err := s.repo.Create(p); err != nil {
		return nil, err
	}

	return p, nil
}

// GetProductByID retrieves a product by ID
func (s *productService) GetProductByID(id int64) (*model.Product, error) {
	return s.repo.GetByID(id)
}

// ListProductsByUser retrieves all products by a user ID
func (s *productService) ListProductsByUser(userID int64) ([]*model.Product, error) {
	return s.repo.ListByUserID(userID)
}

// UpdateProduct updates existing product data
func (s *productService) UpdateProduct(req *model.UpdateProductRequest) (*model.Product, error) {
	if err := s.validateUpdate(req); err != nil {
		return nil, err
	}

	existing, err := s.repo.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	existing.Name = strings.TrimSpace(req.Name)
	existing.Description = strings.TrimSpace(req.Description)
	existing.Price = req.Price
	existing.Stock = req.Stock

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// DeleteProduct deletes a product by ID
func (s *productService) DeleteProduct(id int64) error {
	return s.repo.Delete(id)
}

// validateCreate validates product creation request
func (s *productService) validateCreate(req *model.CreateProductRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if req.Price <= 0 {
		return fmt.Errorf("price must be greater than zero")
	}
	if req.Stock < 0 {
		return fmt.Errorf("stock cannot be negative")
	}
	if req.UserID <= 0 {
		return fmt.Errorf("user_id is required")
	}
	return nil
}

// validateUpdate validates product update request
func (s *productService) validateUpdate(req *model.UpdateProductRequest) error {
	if req.ID <= 0 {
		return fmt.Errorf("id is required")
	}
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if req.Price <= 0 {
		return fmt.Errorf("price must be greater than zero")
	}
	if req.Stock < 0 {
		return fmt.Errorf("stock cannot be negative")
	}
	return nil
}
