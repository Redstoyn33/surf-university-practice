package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/glini/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type RatingRepo struct {
	db *sqlx.DB
}

func NewRatingRepo(db *sqlx.DB) *RatingRepo {
	return &RatingRepo{db: db}
}

func (r *RatingRepo) InsertRating(ctx context.Context, clientID, masterID, slotID int64, score int) (domain.Rating, error) {
	var rating domain.Rating
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO ratings (client_id, master_id, slot_id, score, created_at)
		 VALUES ($1, $2, $3, $4, datetime('now'))
		 RETURNING id, client_id, master_id, slot_id, score, created_at`,
		clientID, masterID, slotID, score,
	).StructScan(&rating)
	if err != nil {
		if isUniqueConstraintErr(err) {
			return domain.Rating{}, fmt.Errorf("duplicate rating: %w", domain.ErrConflict)
		}
		return domain.Rating{}, fmt.Errorf("insert rating: %w", err)
	}

	if err := r.updateMasterRating(ctx, masterID); err != nil {
		return domain.Rating{}, fmt.Errorf("update master rating: %w", err)
	}

	return rating, nil
}

func (r *RatingRepo) GetRatingByClientAndSlot(ctx context.Context, clientID, slotID int64) (domain.Rating, error) {
	var rating domain.Rating
	err := r.db.GetContext(ctx, &rating,
		`SELECT id, client_id, master_id, slot_id, score, created_at
		 FROM ratings WHERE client_id = $1 AND slot_id = $2`, clientID, slotID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Rating{}, fmt.Errorf("rating not found: %w", domain.ErrNotFound)
		}
		return domain.Rating{}, fmt.Errorf("get rating: %w", err)
	}
	return rating, nil
}

func (r *RatingRepo) updateMasterRating(ctx context.Context, masterID int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE masters SET rating = (
			SELECT ROUND(AVG(CAST(score AS REAL)), 1) FROM ratings WHERE master_id = $1
		) WHERE id = $1`, masterID)
	if err != nil {
		return fmt.Errorf("update master rating: %w", err)
	}
	return nil
}
