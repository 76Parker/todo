package table

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Task struct {
	id          int64                     `db:"id"`
	title       string                    `db:"title"`
	status      string                    `db:"status"`
	description pgtype.Text               `db:"description"`
	category    pgtype.Text               `db:"category"`
	tags        pgtype.Array[pgtype.Text] `db:"tags"`
	createdAt   time.Time                 `db:"created_at"`
}
