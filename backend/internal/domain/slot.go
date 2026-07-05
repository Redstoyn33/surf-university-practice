package domain

import "context"

type Slot struct {
	ID              int64    `db:"id" json:"id"`
	DateTime        string   `db:"date_time" json:"dateTime"`
	EndTime         string   `db:"end_time" json:"endTime"`
	ProgramID       int64    `db:"program_id" json:"-"`
	MasterID        int64    `db:"master_id" json:"-"`
	TotalSpots      int      `db:"total_spots" json:"totalSpots"`
	AvailableSpots  int      `db:"available_spots" json:"availableSpots"`
	RentalAvailable bool     `db:"rental_available" json:"rentalAvailable"`
	RentalPrice     float64  `db:"rental_price" json:"rentalPrice"`
	Program         *Program `json:"program,omitempty"`
	Master          *Master  `json:"master,omitempty"`
}

type SlotFilter struct {
	DateFrom  string
	DateTo    string
	MasterID  *int64
	ProgramID *int64
}

type SlotRepository interface {
	QuerySlots(ctx context.Context, filter SlotFilter) ([]Slot, error)
	GetSlotByID(ctx context.Context, id int64) (Slot, error)
	DecrementSpots(ctx context.Context, slotID int64) error
	IncrementSpots(ctx context.Context, slotID int64) error
}
