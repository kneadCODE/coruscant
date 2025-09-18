package pgtest_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kneadCODE/coruscant/shared/golib/testing/pgtest"
)

// TestWithDB_Isolation demonstrates test isolation with automatic rollback
func TestWithDB_Isolation(t *testing.T) {
	// Test 1: Insert data
	pgtest.TestWithDB(t, func(ctx context.Context, tx pgx.Tx) error {
		// Create a temporary table for testing
		_, err := tx.Exec(ctx, `
			CREATE TEMP TABLE test_users (
				id SERIAL PRIMARY KEY,
				name TEXT NOT NULL,
				email TEXT UNIQUE NOT NULL
			)
		`)
		require.NoError(t, err)

		// Insert test data
		_, err = tx.Exec(ctx,
			"INSERT INTO test_users (name, email) VALUES ($1, $2)",
			"John Doe", "john@example.com")
		require.NoError(t, err)

		// Verify data was inserted
		var count int
		err = tx.QueryRow(ctx, "SELECT COUNT(*) FROM test_users").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count, "Expected 1 user to be inserted")

		return nil
	})

	// Test 2: Verify isolation - previous data should not exist
	pgtest.TestWithDB(t, func(ctx context.Context, tx pgx.Tx) error {
		// Create the same temp table (previous one was rolled back)
		_, err := tx.Exec(ctx, `
			CREATE TEMP TABLE test_users (
				id SERIAL PRIMARY KEY,
				name TEXT NOT NULL,
				email TEXT UNIQUE NOT NULL
			)
		`)
		require.NoError(t, err)

		// Verify no data exists from previous test
		var count int
		err = tx.QueryRow(ctx, "SELECT COUNT(*) FROM test_users").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 0, count, "Expected no users (previous test data rolled back)")

		// Insert different data
		_, err = tx.Exec(ctx,
			"INSERT INTO test_users (name, email) VALUES ($1, $2)",
			"Jane Smith", "jane@example.com")
		require.NoError(t, err)

		// Verify our data is there
		var name string
		err = tx.QueryRow(ctx, "SELECT name FROM test_users WHERE email = $1", "jane@example.com").Scan(&name)
		require.NoError(t, err)
		assert.Equal(t, "Jane Smith", name)

		return nil
	})
}

// TestWithDB_ErrorHandling demonstrates error handling with automatic rollback
func TestWithDB_ErrorHandling(t *testing.T) {
	pgtest.TestWithDB(t, func(ctx context.Context, tx pgx.Tx) error {
		// Create temp table
		_, err := tx.Exec(ctx, `
			CREATE TEMP TABLE test_products (
				id SERIAL PRIMARY KEY,
				name TEXT NOT NULL UNIQUE,
				price DECIMAL(10,2) NOT NULL CHECK (price > 0)
			)
		`)
		require.NoError(t, err)

		// Insert valid data
		_, err = tx.Exec(ctx,
			"INSERT INTO test_products (name, price) VALUES ($1, $2)",
			"Laptop", 999.99)
		require.NoError(t, err)

		// Verify insertion worked
		var count int
		err = tx.QueryRow(ctx, "SELECT COUNT(*) FROM test_products").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count)

		// Test constraint violation (duplicate name)
		_, err = tx.Exec(ctx,
			"INSERT INTO test_products (name, price) VALUES ($1, $2)",
			"Laptop", 1299.99) // Same name, should fail
		assert.Error(t, err, "Expected unique constraint violation")

		return nil // Test completes successfully despite constraint error
	})
}

// TestWithDB_ComplexQuery demonstrates complex queries with joins
func TestWithDB_ComplexQuery(t *testing.T) {
	pgtest.TestWithDB(t, func(ctx context.Context, tx pgx.Tx) error {
		// Create related tables
		_, err := tx.Exec(ctx, `
			CREATE TEMP TABLE categories (
				id SERIAL PRIMARY KEY,
				name TEXT NOT NULL
			);
			
			CREATE TEMP TABLE products (
				id SERIAL PRIMARY KEY,
				name TEXT NOT NULL,
				category_id INTEGER REFERENCES categories(id),
				price DECIMAL(10,2)
			);
		`)
		require.NoError(t, err)

		// Insert test data
		var categoryID int
		err = tx.QueryRow(ctx,
			"INSERT INTO categories (name) VALUES ($1) RETURNING id",
			"Electronics").Scan(&categoryID)
		require.NoError(t, err)

		_, err = tx.Exec(ctx,
			"INSERT INTO products (name, category_id, price) VALUES ($1, $2, $3)",
			"Smartphone", categoryID, 699.99)
		require.NoError(t, err)

		// Test complex query with join
		type ProductInfo struct {
			Name         string
			CategoryName string
			Price        float64
		}

		var product ProductInfo
		err = tx.QueryRow(ctx, `
			SELECT p.name, c.name, p.price
			FROM products p
			JOIN categories c ON p.category_id = c.id
			WHERE p.name = $1
		`, "Smartphone").Scan(&product.Name, &product.CategoryName, &product.Price)

		require.NoError(t, err)
		assert.Equal(t, "Smartphone", product.Name)
		assert.Equal(t, "Electronics", product.CategoryName)
		assert.Equal(t, 699.99, product.Price)

		return nil
	})
}
