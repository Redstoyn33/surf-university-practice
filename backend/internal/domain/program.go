package domain

import "context"

type Program struct {
	ID           int64  `db:"id" json:"id"`
	Name         string `db:"name" json:"name"`
	Description  string `db:"description" json:"description"`
	MaxCapacity  int    `db:"max_capacity" json:"maxCapacity"`
	MasterIDs    []int64 `json:"masterIds"`
}

type ProgramRepository interface {
	QueryPrograms(ctx context.Context) ([]Program, error)
	GetProgramByID(ctx context.Context, id int64) (Program, error)
}
