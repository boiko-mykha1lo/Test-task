package repository

import (
	"database/sql"
	"fmt"

	"github.com/yourname/products/internal/models"
)

type ProductRepository interface {
	Create(product *models.CreateProductRequest) (*models.Product, error)
	Delete(id int) error
	List(page, pageSize int) ([]models.Product, int, error)
	GetByID(id int) (*models.Product, error)
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(req *models.CreateProductRequest) (*models.Product, error) {
	query := `
        INSERT INTO products (name, description, price)
        VALUES ($1, $2, $3)
        RETURNING id, name, description, price, created_at, updated_at
    `

	product := &models.Product{}
	err := r.db.QueryRow(query, req.Name, req.Description, req.Price).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (r *productRepository) Delete(id int) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (r *productRepository) List(page, pageSize int) ([]models.Product, int, error) {
	offset := (page - 1) * pageSize

	var total int
	countQuery := `SELECT COUNT(*) FROM products`
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	query := `
        SELECT id, name, description, price, created_at, updated_at
        FROM products
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}

	return products, total, nil
}

func (r *productRepository) GetByID(id int) (*models.Product, error) {
	query := `
        SELECT id, name, description, price, created_at, updated_at
        FROM products
        WHERE id = $1
    `

	product := &models.Product{}
	err := r.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}
