package domain

import "context"

type Rating struct {
	ID        int64  `db:"id" json:"id"`
	ClientID  int64  `db:"client_id" json:"clientId"`
	MasterID  int64  `db:"master_id" json:"masterId"`
	SlotID    int64  `db:"slot_id" json:"slotId"`
	Score     int    `db:"score" json:"score"`
	CreatedAt string `db:"created_at" json:"createdAt"`
}

type RatingRepository interface {
	InsertRating(ctx context.Context, clientID, masterID, slotID int64, score int) (Rating, error)
	GetRatingByClientAndSlot(ctx context.Context, clientID, slotID int64) (Rating, error)
}
