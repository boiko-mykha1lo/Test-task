package repository

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourname/products/internal/models"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/products_test?sslmode=disable")
	require.NoError(t, err)

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS products (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            description TEXT,
            price DECIMAL(10, 2) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
	require.NoError(t, err)

	return db
}

func teardownTestDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS products")
	require.NoError(t, err)
	db.Close()
}

func TestProductRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	repo := NewProductRepository(db)

	req := &models.CreateProductRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
	}

	product, err := repo.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Greater(t, product.ID, 0)
	assert.Equal(t, req.Name, product.Name)
	assert.Equal(t, req.Description, product.Description)
	assert.Equal(t, req.Price, product.Price)
}

func TestProductRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	repo := NewProductRepository(db)

	req := &models.CreateProductRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
	}
	product, err := repo.Create(req)
	require.NoError(t, err)

	err = repo.Delete(product.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(product.ID)
	assert.Error(t, err)
}

func TestProductRepository_List(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	repo := NewProductRepository(db)

	for i := 0; i < 5; i++ {
		req := &models.CreateProductRequest{
			Name:        "Product " + string(rune(i)),
			Description: "Description",
			Price:       float64(i * 10),
		}
		_, err := repo.Create(req)
		require.NoError(t, err)
	}

	products, total, err := repo.List(1, 3)

	assert.NoError(t, err)
	assert.Equal(t, 5, total)
	assert.Len(t, products, 3)
}
