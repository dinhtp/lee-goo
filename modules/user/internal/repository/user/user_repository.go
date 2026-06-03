package user

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/dinhtp/lee-goo/modules/user/contracts"
	domainUser "github.com/dinhtp/lee-goo/modules/user/internal/domain/user"
)

type userRepository struct {
	db *sqlx.DB
}

// compile-time interface check
var _ domainUser.UserPort = (*userRepository)(nil)

// NewRepository constructs a SQL-backed UserPort.
func NewRepository(db *sqlx.DB) domainUser.UserPort {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*domainUser.User, error) {
	const q = `SELECT id, email, name, password_hash, created_at, updated_at
	           FROM users WHERE id = $1`

	u := &domainUser.User{}
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, contracts.ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domainUser.User, error) {
	const q = `SELECT id, email, name, password_hash, created_at, updated_at
	           FROM users WHERE email = $1`

	u := &domainUser.User{}
	err := r.db.QueryRowContext(ctx, q, email).Scan(
		&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, contracts.ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *userRepository) List(ctx context.Context) ([]domainUser.User, error) {
	const q = `SELECT id, email, name, password_hash, created_at, updated_at
	           FROM users ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domainUser.User
	for rows.Next() {
		var u domainUser.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *userRepository) Create(ctx context.Context, u domainUser.User) error {
	const q = `INSERT INTO users (id, email, name, password_hash, created_at, updated_at)
	           VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, q, u.ID, u.Email, u.Name, u.PasswordHash, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return contracts.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (r *userRepository) Update(ctx context.Context, u domainUser.User) error {
	const q = `UPDATE users SET name=$2, updated_at=$3 WHERE id=$1`

	_, err := r.db.ExecContext(ctx, q, u.ID, u.Name, u.UpdatedAt)
	return err
}

// isUniqueViolation detects PostgreSQL unique-constraint errors without
// importing a pgx-specific error package, keeping the repository driver-agnostic
// at the cost of a string check (acceptable until pgerrcode is added in Phase 6).
func isUniqueViolation(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "unique constraint") ||
		strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "UNIQUE constraint failed")
}

// Ensure time package is used (suppress unused import if compiler complains).
var _ = time.Now
