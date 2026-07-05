package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/glini/backend/internal/domain"
)

type RatingService struct {
	ratingRepo  domain.RatingRepository
	bookingRepo domain.BookingRepository
	slotRepo    domain.SlotRepository
}

func NewRatingService(ratingRepo domain.RatingRepository, bookingRepo domain.BookingRepository, slotRepo domain.SlotRepository) *RatingService {
	return &RatingService{
		ratingRepo:  ratingRepo,
		bookingRepo: bookingRepo,
		slotRepo:    slotRepo,
	}
}

func (s *RatingService) CreateRating(ctx context.Context, clientID, masterID, slotID int64, score int) (domain.Rating, error) {
	if score < 1 || score > 5 {
		return domain.Rating{}, fmt.Errorf("score must be between 1 and 5: %w", domain.ErrValidation)
	}

	bookings, err := s.bookingRepo.QueryBookingsByClient(ctx, clientID, "активна")
	if err != nil {
		return domain.Rating{}, fmt.Errorf("query bookings: %w", err)
	}

	var hasBooking bool
	for _, b := range bookings {
		if b.SlotID == slotID {
			hasBooking = true
			break
		}
	}
	if !hasBooking {
		bookingsAll, err := s.bookingRepo.QueryBookingsByClient(ctx, clientID, "")
		if err != nil {
			return domain.Rating{}, fmt.Errorf("query all bookings: %w", err)
		}
		for _, b := range bookingsAll {
			if b.SlotID == slotID && b.Status == "отменена клиентом" {
				return domain.Rating{}, fmt.Errorf("booking was cancelled: %w", domain.ErrValidation)
			}
		}
		return domain.Rating{}, fmt.Errorf("no active booking for this slot: %w", domain.ErrValidation)
	}

	slot, err := s.slotRepo.GetSlotByID(ctx, slotID)
	if err != nil {
		return domain.Rating{}, fmt.Errorf("get slot: %w", err)
	}

	endTime, err := time.Parse(time.RFC3339, slot.EndTime)
	if err != nil {
		return domain.Rating{}, fmt.Errorf("parse end time: %w", err)
	}

	now := time.Now()
	if now.Before(endTime.Add(1 * time.Hour)) {
		return domain.Rating{}, fmt.Errorf("too early to rate: %w", domain.ErrValidation)
	}
	if now.After(endTime.Add(48 * time.Hour)) {
		return domain.Rating{}, fmt.Errorf("too late to rate: %w", domain.ErrValidation)
	}

	rating, err := s.ratingRepo.InsertRating(ctx, clientID, masterID, slotID, score)
	if err != nil {
		if errors.Is(err, domain.ErrConflict) {
			return domain.Rating{}, fmt.Errorf("already rated: %w", domain.ErrConflict)
		}
		return domain.Rating{}, fmt.Errorf("insert rating: %w", err)
	}

	return rating, nil
}
