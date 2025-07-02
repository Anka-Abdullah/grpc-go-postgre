package repository

import (
	"database/sql"
	"fmt"
	"time"

	"grpc-exmpl/internal/model"
)

// ProductRepository defines contract for product operations
type ProductRepository interface {
	Create(product *model.Product) error
	GetByID(id int64) (*model.Product, error)
	ListByUserID(userID int64) ([]*model.Product, error)
	Update(product *model.Product) error
	Delete(id int64) error
}

type productRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new instance of ProductRepository
func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *model.Product) error {
	query := `
		INSERT INTO products (name, description, price, stock, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	err := r.db.QueryRow(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.UserID,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&product.ID)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

func (r *productRepository) GetByID(id int64) (*model.Product, error) {
	query := `
		SELECT id, name, description, price, stock, user_id, created_at, updated_at
		FROM products WHERE id = $1
	`

	p := &model.Product{}
	err := r.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Stock,
		&p.UserID,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return p, nil
}

func (r *productRepository) ListByUserID(userID int64) ([]*model.Product, error) {
	query := `
		SELECT id, name, description, price, stock, user_id, created_at, updated_at
		FROM products
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		p := &model.Product{}
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.UserID,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return products, nil
}

func (r *productRepository) Update(product *model.Product) error {
	query := `
		UPDATE products
		SET name = $2, description = $3, price = $4, stock = $5, updated_at = $6
		WHERE id = $1
	`

	product.UpdatedAt = time.Now()

	result, err := r.db.Exec(
		query,
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (r *productRepository) Delete(id int64) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
