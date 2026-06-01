package role

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dinhtp/lee-goo/modules/authorization/contracts"
	domainRole "github.com/dinhtp/lee-goo/modules/authorization/internal/domain/role"
	"github.com/dinhtp/lee-goo/system/database"
)

type roleRepository struct {
	db *sql.DB
}

// compile-time interface check
var _ domainRole.RolePort = (*roleRepository)(nil)

// NewRepository constructs a RolePort backed by the platform database connection.
func NewRepository(conn database.Connection) domainRole.RolePort {
	return &roleRepository{db: conn.DB()}
}

func (r *roleRepository) Create(ctx context.Context, role domainRole.Role) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO roles (id, name, description, created_at) VALUES ($1, $2, $3, $4)`,
		role.ID, role.Name, role.Description, role.CreatedAt,
	)
	return err
}

func (r *roleRepository) FindByID(ctx context.Context, id string) (*domainRole.Role, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, description, created_at FROM roles WHERE id = $1`, id,
	)
	role, err := scanRole(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, contracts.ErrRoleNotFound
		}
		return nil, err
	}
	return role, nil
}

func (r *roleRepository) List(ctx context.Context) ([]domainRole.Role, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, description, created_at FROM roles ORDER BY name`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []domainRole.Role
	for rows.Next() {
		var role domainRole.Role
		var createdAt time.Time
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &createdAt); err != nil {
			return nil, err
		}
		role.CreatedAt = createdAt
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (r *roleRepository) AssignToUser(ctx context.Context, userID, roleID string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO user_roles (user_id, role_id, assigned_at) VALUES ($1, $2, NOW()) ON CONFLICT DO NOTHING`,
		userID, roleID,
	)
	return err
}

func (r *roleRepository) RemoveFromUser(ctx context.Context, userID, roleID string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`,
		userID, roleID,
	)
	return err
}

func (r *roleRepository) FindByUserID(ctx context.Context, userID string) ([]domainRole.Role, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT r.id, r.name, r.description, r.created_at
		 FROM roles r
		 JOIN user_roles ur ON r.id = ur.role_id
		 WHERE ur.user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []domainRole.Role
	for rows.Next() {
		var role domainRole.Role
		var createdAt time.Time
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &createdAt); err != nil {
			return nil, err
		}
		role.CreatedAt = createdAt
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

// scanRole reads a single role row from a QueryRow result.
func scanRole(row *sql.Row) (*domainRole.Role, error) {
	var role domainRole.Role
	var createdAt time.Time
	if err := row.Scan(&role.ID, &role.Name, &role.Description, &createdAt); err != nil {
		return nil, err
	}
	role.CreatedAt = createdAt
	return &role, nil
}
