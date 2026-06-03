package permission

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	domainRole "github.com/dinhtp/lee-goo/modules/authorization/internal/domain/role"
)

type permissionRepository struct {
	db *sqlx.DB
}

// compile-time interface check
var _ domainRole.PermissionPort = (*permissionRepository)(nil)

// NewRepository constructs a PermissionPort backed by a *sqlx.DB.
func NewRepository(db *sqlx.DB) domainRole.PermissionPort {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) Create(ctx context.Context, p domainRole.Permission) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO permissions (id, action, resource) VALUES ($1, $2, $3)
		 ON CONFLICT (action, resource) DO NOTHING`,
		p.ID, p.Action, p.Resource,
	)
	return err
}

func (r *permissionRepository) AssignToRole(ctx context.Context, roleID, permissionID string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2)
		 ON CONFLICT DO NOTHING`,
		roleID, permissionID,
	)
	return err
}

// FindByRoleIDs returns distinct permissions associated with any of the given role IDs.
// Uses individual placeholders to stay compatible with the pgx stdlib driver.
func (r *permissionRepository) FindByRoleIDs(ctx context.Context, roleIDs []string) ([]domainRole.Permission, error) {
	if len(roleIDs) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(roleIDs))
	args := make([]interface{}, len(roleIDs))
	for i, id := range roleIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(
		`SELECT DISTINCT p.id, p.action, p.resource
		 FROM permissions p
		 JOIN role_permissions rp ON p.id = rp.permission_id
		 WHERE rp.role_id IN (%s)`,
		strings.Join(placeholders, ","),
	)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []domainRole.Permission
	for rows.Next() {
		var p domainRole.Permission
		if err := rows.Scan(&p.ID, &p.Action, &p.Resource); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, rows.Err()
}
