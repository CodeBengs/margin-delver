package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(up, down)
}

func up(ctx context.Context, tx *sql.Tx) error {
	query := `
CREATE TABLE IF NOT EXISTS ms_users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    username VARCHAR(64) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(128) NOT NULL,
    flag_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    deleted_at DATETIME(3) NULL,
    PRIMARY KEY (id),
    UNIQUE KEY idx_users_username (username),
    KEY idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	insert := `
INSERT INTO ms_users (
    username,
    password_hash,
    name,
    flag_active,
    created_at,
    updated_at
) VALUES (
    'DELVERADMIN1',
    '$2a$10$hA8jOylTEJbBLlfMCDCUc.eUugFOFo/d.ctWfj2dLF2J2L8IMW9iC',
    'Delver Administrator',
    TRUE,
    CURRENT_TIMESTAMP(3),
    CURRENT_TIMESTAMP(3)
) ON DUPLICATE KEY UPDATE
    username = username;
`

	if _, err := tx.ExecContext(ctx, insert); err != nil {
		return err
	}

	return nil
}

func down(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS ms_users")
	return err
}
