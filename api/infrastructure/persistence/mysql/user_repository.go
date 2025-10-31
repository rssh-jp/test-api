package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/rssh-jp/test-api/api/domain"
)

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	query := `SELECT id, username, email, created_at, updated_at FROM users ORDER BY created_at DESC`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		// Note: age is not in the database schema, so we don't set it
		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?`
	
	var user domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	// Note: age is not in the database schema, so we don't set it
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (username, email, password_hash, status, email_verified, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Use default values for required fields that don't exist in domain.User
	passwordHash := "$2a$10$defaultpasswordhash" // Default placeholder
	status := "active"
	emailVerified := 0

	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, passwordHash, status, emailVerified, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET username = ?, email = ?, updated_at = ? WHERE id = ?`
	user.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.UpdatedAt, user.ID)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
