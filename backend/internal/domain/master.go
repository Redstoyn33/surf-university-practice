package domain

import "context"

type Master struct {
	ID         int64   `db:"id" json:"id"`
	Name       string  `db:"name" json:"name"`
	PhotoURL   string  `db:"photo_url" json:"photo"`
	Rating     float64 `db:"rating" json:"rating"`
	Level      string  `db:"level" json:"level"`
	ProgramIDs []int64 `json:"programIds"`
}

type MasterRepository interface {
	QueryMasters(ctx context.Context) ([]Master, error)
	GetMasterByID(ctx context.Context, id int64) (Master, error)
}
