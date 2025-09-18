package pg_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/kneadCODE/coruscant/shared/golib/pg"
)

// ExampleNewClient demonstrates basic client usage with options pattern
func ExampleNewClient() {
	// Create client with explicit options
	client, err := pg.NewClient(context.Background(),
		pg.WithHost("localhost"),
		pg.WithPort(5432),
		pg.WithDatabase("myapp"),
		pg.WithCredentials("user", "password"),
		pg.WithSSLMode("prefer"),
		pg.WithMaxConnections(25),
		pg.WithMinConnections(5),
	)
	if err != nil {
		log.Fatal("Failed to create database client:", err)
	}
	defer client.Close()

	// Test connection
	if err := client.Ping(context.Background()); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Successfully connected to PostgreSQL")
}

// ExampleNewClient_environment demonstrates environment-based configuration
func ExampleNewClient_environment() {
	// Configure from environment variables
	options := []pg.Option{
		pg.WithHost(getEnv("DB_HOST", "localhost")),
		pg.WithPort(5432),
		pg.WithDatabase(getEnv("DB_NAME", "myapp")),
		pg.WithCredentials(getEnv("DB_USER", "user"), getEnv("DB_PASS", "password")),
		pg.WithSSLMode("prefer"),
	}

	client, err := pg.NewClient(context.Background(), options...)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}
	defer client.Close()

	fmt.Println("Client created from environment configuration")
}

// ExampleClient_Query demonstrates querying data
func ExampleClient_Query() {
	client := setupExampleClient()
	defer client.Close()

	ctx := context.Background()
	rows, err := client.Query(ctx, "SELECT id, name, email FROM users WHERE active = $1", true)
	if err != nil {
		log.Fatal("Query failed:", err)
	}
	defer rows.Close()

	type User struct {
		ID    int
		Name  string
		Email string
	}

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			log.Fatal("Scan failed:", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Row iteration failed:", err)
	}

	fmt.Printf("Found %d active users\n", len(users))
}

// ExampleClient_WithTx demonstrates transaction usage
func ExampleClient_WithTx() {
	client := setupExampleClient()
	defer client.Close()

	ctx := context.Background()

	err := client.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		// Insert user
		var userID int
		err := tx.QueryRow(ctx,
			"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
			"John Doe", "john@example.com").Scan(&userID)
		if err != nil {
			return err
		}

		// Insert user profile
		_, err = tx.Exec(ctx,
			"INSERT INTO user_profiles (user_id, bio) VALUES ($1, $2)",
			userID, "Software engineer")
		return err
	})
	if err != nil {
		log.Fatal("Transaction failed:", err)
	}

	fmt.Println("User and profile created successfully")
}

// ExampleClient_Ping demonstrates health checking with client.Ping
func ExampleClient_Ping() {
	client := setupExampleClient()
	defer client.Close()

	ctx := context.Background()
	if err := client.Ping(ctx); err != nil {
		log.Printf("Database is unhealthy: %v\n", err)
	} else {
		fmt.Println("Database is healthy")
	}

	// Ping with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx); err != nil {
		log.Printf("Database ping failed: %v\n", err)
	} else {
		fmt.Println("Database ping successful")
	}
}

// Helper functions
func setupExampleClient() *pg.Client {
	client, err := pg.NewClient(context.Background(),
		pg.WithHost("localhost"),
		pg.WithPort(5432),
		pg.WithDatabase("postgres"),
		pg.WithCredentials("coruscant", "trust"),
		pg.WithSSLMode("disable"),
	)
	if err != nil {
		log.Fatal("Failed to create example client:", err)
	}
	return client
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
