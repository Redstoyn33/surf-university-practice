package domain

import "context"

type Booking struct {
	ID                 int64   `db:"id" json:"id"`
	ClientID           int64   `db:"client_id" json:"clientId"`
	SlotID             int64   `db:"slot_id" json:"-"`
	Status             string  `db:"status" json:"status"`
	RentalSelected     bool    `db:"rental_selected" json:"rentalSelected"`
	CreatedAt          string  `db:"created_at" json:"createdAt"`
	CancellationReason *string `db:"cancellation_reason" json:"cancellationReason,omitempty"`
	Slot               *Slot   `json:"slot,omitempty"`
}

type BookingRepository interface {
	InsertBooking(ctx context.Context, clientID, slotID int64, rentalSelected bool) (Booking, error)
	QueryBookingsByClient(ctx context.Context, clientID int64, statusFilter string) ([]Booking, error)
	GetBookingByID(ctx context.Context, id int64) (Booking, error)
	UpdateBookingStatus(ctx context.Context, id int64, status string, reason *string) error
	CancelBookingTx(ctx context.Context, bookingID, slotID int64) error
}
