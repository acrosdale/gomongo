package pgdb

import (
	"time"

	"github.com/uptrace/bun"
)

type Poll struct {
	bun.BaseModel `bun:"table:polls,alias:u"`

	ID        int64     `bun:"id,pk,autoincrement"`
	Question  string    `bun:"question,notnull"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

type Comment struct {
	bun.BaseModel `bun:"table:comments,alias:u"`

	ID        int64     `bun:"id,pk,autoincrement"`
	Body      string    `bun:"body,notnull"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`

	PollID int64 `bun:"poll_id"`
	Poll   Poll  `bun:"rel:belongs-to,join:poll_id=id"`
}
