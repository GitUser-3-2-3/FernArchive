package data

import (
	"context"
	"database/sql"
	"log/slog"
	"slices"
	"time"

	"github.com/lib/pq"
)

type Permissions []string

func (psn Permissions) Include(code string) bool {
	return slices.Contains(psn, code)
}

type PermissionModel struct {
	Db *sql.DB
}

func (mdl PermissionModel) GetAllForUser(userId int64) (Permissions, error) {
	query := `SELECT permissions.code FROM permissions 
                INNER JOIN user_permissions ON user_permissions.permission_id = permissions.id
                INNER JOIN users ON user_permissions.user_id = users.id
                WHERE users.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := mdl.Db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			slog.Error("Failed to close rows: ", err)
		}
	}(rows)
	var permissions Permissions
	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return permissions, nil
}

func (mdl PermissionModel) AddForUser(userId int64, codes ...string) error {
	query := `INSERT INTO user_permissions
                SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := mdl.Db.ExecContext(ctx, query, userId, pq.Array(codes))
	return err
}
