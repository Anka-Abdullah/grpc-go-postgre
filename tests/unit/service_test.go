package unit

import (
	"testing"

	"grpc-exmpl/internal/model"
	"grpc-exmpl/internal/service"
)

type fakeProductRepo struct {
	created *model.Product
	stored  *model.Product
}

func (f *fakeProductRepo) Create(p *model.Product) error                       { f.created = p; p.ID = 1; return nil }
func (f *fakeProductRepo) GetByID(id int64) (*model.Product, error)            { return f.stored, nil }
func (f *fakeProductRepo) ListByUserID(userID int64) ([]*model.Product, error) { return nil, nil }
func (f *fakeProductRepo) Update(p *model.Product) error                       { f.stored = p; return nil }
func (f *fakeProductRepo) Delete(id int64) error                               { return nil }

func TestProductServiceCreateValidation(t *testing.T) {
	repo := &fakeProductRepo{}
	svc := service.NewProductService(repo)

	_, err := svc.CreateProduct(&model.CreateProductRequest{Name: "", Price: 1, Stock: 1, UserID: 1})
	if err == nil {
		t.Fatal("expected validation error")
	}

	p, err := svc.CreateProduct(&model.CreateProductRequest{Name: "Book", Description: "good", Price: 2, Stock: 1, UserID: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.created == nil || p.ID != 1 {
		t.Fatalf("repo not called or id not set")
	}
}

func TestProductServiceUpdate(t *testing.T) {
	repo := &fakeProductRepo{stored: &model.Product{ID: 1, Name: "old", Description: "d", Price: 1, Stock: 1, UserID: 2}}
	svc := service.NewProductService(repo)

	upd := &model.UpdateProductRequest{ID: 1, Name: "new", Description: "d2", Price: 2, Stock: 5}
	res, err := svc.UpdateProduct(upd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Name != "new" || repo.stored.Price != 2 {
		t.Fatalf("update did not apply")
	}
}
