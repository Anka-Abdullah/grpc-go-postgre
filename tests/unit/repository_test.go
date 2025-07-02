package unit

import (
	"database/sql/driver"
	"testing"
	"time"

	"grpc-exmpl/internal/model"
	"grpc-exmpl/internal/repository"
	"grpc-exmpl/tests/testdb"
)

func TestProductRepositoryCreate(t *testing.T) {
	db, stub, err := testdb.New()
	if err != nil {
		t.Fatalf("failed to create stub db: %v", err)
	}
	repo := repository.NewProductRepository(db)

	stub.QueryFunc = func(query string, args []driver.NamedValue) ([][]driver.Value, error) {
		return [][]driver.Value{{int64(1)}}, nil
	}

	p := &model.Product{Name: "Book", Description: "nice", Price: 10, Stock: 2, UserID: 5}
	if err := repo.Create(p); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if p.ID != 1 {
		t.Fatalf("expected ID 1 got %d", p.ID)
	}
	if p.CreatedAt.IsZero() || p.UpdatedAt.IsZero() {
		t.Fatalf("timestamps not set")
	}
}

func TestProductRepositoryGetByID(t *testing.T) {
	db, stub, err := testdb.New()
	if err != nil {
		t.Fatalf("failed to create stub db: %v", err)
	}
	repo := repository.NewProductRepository(db)

	now := time.Now()
	stub.QueryFunc = func(query string, args []driver.NamedValue) ([][]driver.Value, error) {
		return [][]driver.Value{{int64(2), "Item", "desc", float64(9.9), int64(3), int64(1), now, now}}, nil
	}

	p, err := repo.GetByID(2)
	if err != nil {
		t.Fatalf("GetByID error: %v", err)
	}
	if p.Name != "Item" || p.Price != 9.9 {
		t.Fatalf("unexpected product: %v", p)
	}
}
