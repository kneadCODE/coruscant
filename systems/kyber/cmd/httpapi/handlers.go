package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"github.com/kneadCODE/coruscant/shared/golib/pg"
	"github.com/kneadCODE/coruscant/shared/golib/telemetry"
)

// User represents a user in the system
type User struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Created time.Time `json:"created"`
}

// UserService handles user-related database operations with automatic observability
type UserService struct {
	db *pg.Client
}

// NewUserService creates a new user service
func NewUserService(db *pg.Client) *UserService {
	return &UserService{db: db}
}

// UserService database operations with automatic observability

// CreateUser creates a new user in the database
func (s *UserService) CreateUser(ctx context.Context, name, email string) (*User, error) {
	var user User
	
	// Example of using WithTx for transactional operations
	err := s.db.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		// Insert user and get the generated ID
		query := `INSERT INTO users (name, email, created) VALUES ($1, $2, NOW()) RETURNING id, name, email, created`
		err := tx.QueryRow(ctx, query, name, email).Scan(&user.ID, &user.Name, &user.Email, &user.Created)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		
		// Log user creation for observability
		telemetry.RecordInfoEvent(ctx, "User created successfully", 
			"user_id", user.ID,
			"user_name", user.Name,
		)
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, userID int) (*User, error) {
	var user User
	
	query := `SELECT id, name, email, created FROM users WHERE id = $1`
	err := s.db.QueryRow(ctx, query, userID).Scan(&user.ID, &user.Name, &user.Email, &user.Created)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %d: %w", userID, err)
	}
	
	return &user, nil
}

// ListUsers retrieves all users with pagination
func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]User, error) {
	query := `SELECT id, name, email, created FROM users ORDER BY created DESC LIMIT $1 OFFSET $2`
	rows, err := s.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()
	
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Created)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}
	
	return users, nil
}

// BatchCreateUsers demonstrates batch operations
func (s *UserService) BatchCreateUsers(ctx context.Context, users []User) error {
	return s.db.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		batch := &pgx.Batch{}
		
		for _, user := range users {
			batch.Queue("INSERT INTO users (name, email, created) VALUES ($1, $2, NOW())", 
				user.Name, user.Email)
		}
		
		results := tx.SendBatch(ctx, batch)
		defer results.Close()
		
		// Process results
		for i := 0; i < len(users); i++ {
			_, err := results.Exec()
			if err != nil {
				return fmt.Errorf("failed to create user %d in batch: %w", i, err)
			}
		}
		
		telemetry.RecordInfoEvent(ctx, "Batch user creation completed", 
			"users_created", len(users),
		)
		
		return nil
	})
}

// HTTP Handlers

// setupRoutes configures all HTTP routes
func setupRoutes(router chi.Router, userService *UserService) {
	// Original testing endpoints
	router.Route("/testing", func(r chi.Router) {
		r.Get("/", testingHandler)
		r.Get("/abc", testingHandler)
	})
	router.Get("/testing2", testingHandler)

	// Database-powered user endpoints
	router.Route("/users", func(r chi.Router) {
		r.Post("/", createUserHandler(userService))
		r.Get("/", listUsersHandler(userService))
		r.Get("/{id}", getUserHandler(userService))
		r.Post("/batch", batchCreateUsersHandler(userService))
	})

	// Health check with database connectivity
	router.Get("/health", healthCheckHandler)
}

func testingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	telemetry.RecordInfoEvent(ctx, "testing endpoint called")

	someFunc(ctx)

	telemetry.RecordInfoEvent(ctx, "testing endpoint response sent")
}

func someFunc(ctx context.Context) {
	ctx, end := telemetry.Measure(ctx, "response-preparation")
	defer end(nil)

	// Simulate response preparation work
	time.Sleep(10 * time.Millisecond)

	telemetry.RecordInfoEvent(ctx, "response prepared")
}

// createUserHandler handles POST /users
func createUserHandler(userService *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		
		var req struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		
		if req.Name == "" || req.Email == "" {
			http.Error(w, "Name and email are required", http.StatusBadRequest)
			return
		}
		
		user, err := userService.CreateUser(ctx, req.Name, req.Email)
		if err != nil {
			telemetry.RecordErrorEvent(ctx, fmt.Errorf("failed to create user: %w", err))
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

// getUserHandler handles GET /users/{id}
func getUserHandler(userService *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		idStr := chi.URLParam(r, "id")
		
		userID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		
		user, err := userService.GetUser(ctx, userID)
		if err != nil {
			telemetry.RecordErrorEvent(ctx, fmt.Errorf("failed to get user: %w", err))
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// listUsersHandler handles GET /users
func listUsersHandler(userService *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		
		// Parse query parameters
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")
		
		limit := 10 // default
		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
				limit = l
			}
		}
		
		offset := 0 // default
		if offsetStr != "" {
			if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
				offset = o
			}
		}
		
		users, err := userService.ListUsers(ctx, limit, offset)
		if err != nil {
			telemetry.RecordErrorEvent(ctx, fmt.Errorf("failed to list users: %w", err))
			http.Error(w, "Failed to list users", http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

// batchCreateUsersHandler handles POST /users/batch
func batchCreateUsersHandler(userService *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		
		var users []User
		if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		
		if len(users) == 0 {
			http.Error(w, "No users provided", http.StatusBadRequest)
			return
		}
		
		// Validate users
		for i, user := range users {
			if user.Name == "" || user.Email == "" {
				http.Error(w, fmt.Sprintf("User %d missing name or email", i), http.StatusBadRequest)
				return
			}
		}
		
		if err := userService.BatchCreateUsers(ctx, users); err != nil {
			telemetry.RecordErrorEvent(ctx, fmt.Errorf("failed to batch create users: %w", err))
			http.Error(w, "Failed to create users", http.StatusInternalServerError)
			return
		}
		
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Users created successfully"))
	}
}

// healthCheckHandler handles GET /health
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	if err := dbClient.Ping(ctx); err != nil {
		telemetry.RecordErrorEvent(ctx, fmt.Errorf("health check failed: %w", err))
		http.Error(w, "Database unhealthy", http.StatusServiceUnavailable)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// createSchema creates the database schema (in production, use proper migrations)
func createSchema(ctx context.Context, db *pg.Client) error {
	schema := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			created TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
		
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
		CREATE INDEX IF NOT EXISTS idx_users_created ON users(created);
	`
	
	_, err := db.Exec(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}
	
	telemetry.RecordInfoEvent(ctx, "Database schema created successfully")
	return nil
}