// Package table contains table struct for DB representation
package table

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Task table-struct for DB mapping
type Task struct {
	ID          int64                     `db:"id"`
	Title       string                    `db:"title"`
	Status      string                    `db:"status"`
	Description pgtype.Text               `db:"description"`
	Category    pgtype.Text               `db:"category"`
	Tags        pgtype.Array[pgtype.Text] `db:"tags"`
	CreatedAt   time.Time                 `db:"created_at"`
}
