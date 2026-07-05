package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/glini/backend/internal/domain"
)

type BookingService struct {
	bookingRepo domain.BookingRepository
	slotRepo    domain.SlotRepository
}

func NewBookingService(bookingRepo domain.BookingRepository, slotRepo domain.SlotRepository) *BookingService {
	return &BookingService{bookingRepo: bookingRepo, slotRepo: slotRepo}
}

func (s *BookingService) CreateBooking(ctx context.Context, clientID, slotID int64, rentalSelected bool) (domain.Booking, error) {
	slot, err := s.slotRepo.GetSlotByID(ctx, slotID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Booking{}, fmt.Errorf("slot not found: %w", domain.ErrNotFound)
		}
		return domain.Booking{}, fmt.Errorf("get slot: %w", err)
	}

	if slot.AvailableSpots <= 0 {
		return domain.Booking{}, fmt.Errorf("no available spots: %w", domain.ErrConflict)
	}

	if rentalSelected && !slot.RentalAvailable {
		return domain.Booking{}, fmt.Errorf("rental not available: %w", domain.ErrValidation)
	}

	booking, err := s.bookingRepo.InsertBooking(ctx, clientID, slotID, rentalSelected)
	if err != nil {
		if errors.Is(err, domain.ErrConflict) {
			return domain.Booking{}, fmt.Errorf("booking conflict: %w", domain.ErrConflict)
		}
		return domain.Booking{}, fmt.Errorf("create booking: %w", err)
	}

	return booking, nil
}

func (s *BookingService) CancelBooking(ctx context.Context, bookingID, clientID int64) (domain.Booking, error) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Booking{}, fmt.Errorf("booking not found: %w", domain.ErrNotFound)
		}
		return domain.Booking{}, fmt.Errorf("get booking: %w", err)
	}

	if booking.ClientID != clientID {
		return domain.Booking{}, fmt.Errorf("not your booking: %w", domain.ErrForbidden)
	}

	if booking.Status != "активна" {
		return domain.Booking{}, fmt.Errorf("booking is not active: %w", domain.ErrNotActive)
	}

	if booking.Slot == nil || booking.Slot.DateTime == "" {
		slot, err := s.slotRepo.GetSlotByID(ctx, booking.SlotID)
		if err != nil {
			return domain.Booking{}, fmt.Errorf("get slot: %w", err)
		}
		booking.Slot = &slot
	}

	slotTime, err := time.Parse(time.RFC3339, booking.Slot.DateTime)
	if err != nil {
		return domain.Booking{}, fmt.Errorf("parse slot time: %w", err)
	}

	if time.Until(slotTime) < 4*time.Hour {
		return domain.Booking{}, fmt.Errorf("less than 4 hours before slot: %w", domain.ErrValidation)
	}

	if err := s.bookingRepo.CancelBookingTx(ctx, bookingID, booking.SlotID); err != nil {
		return domain.Booking{}, fmt.Errorf("cancel booking: %w", err)
	}

	updated, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return domain.Booking{}, fmt.Errorf("get updated booking: %w", err)
	}

	return updated, nil
}

func (s *BookingService) GetBookingByID(ctx context.Context, bookingID, clientID int64) (domain.Booking, error) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Booking{}, fmt.Errorf("booking not found: %w", domain.ErrNotFound)
		}
		return domain.Booking{}, fmt.Errorf("get booking: %w", err)
	}
	if booking.ClientID != clientID {
		return domain.Booking{}, fmt.Errorf("booking not found: %w", domain.ErrNotFound)
	}
	return booking, nil
}

func (s *BookingService) GetMyBookings(ctx context.Context, clientID int64, statusFilter string) ([]domain.Booking, error) {
	return s.bookingRepo.QueryBookingsByClient(ctx, clientID, statusFilter)
}
