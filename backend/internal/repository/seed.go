package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

func SeedData(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO masters (id, name, photo_url, rating, level) VALUES
		(1, 'Анна Кузнецова', 'https://example.com/photos/anna.jpg', 4.8, 'опытный'),
		(2, 'Иван Петров', 'https://example.com/photos/ivan.jpg', 4.5, 'опытный'),
		(3, 'Елена Соколова', 'https://example.com/photos/elena.jpg', 4.2, 'новичок')`)
	if err != nil {
		return fmt.Errorf("insert masters: %w", err)
	}

	_, err = tx.Exec(`INSERT INTO programs (id, name, description, max_capacity) VALUES
		(1, 'Лепка для новичков', 'Основы ручной лепки из глины. Создайте своё первое изделие.', 6),
		(2, 'Гончарный круг', 'Освойте работу на гончарном круге под руководством мастера.', 4),
		(3, 'Роспись керамики', 'Декоративная роспись готовых керамических изделий.', 8)`)
	if err != nil {
		return fmt.Errorf("insert programs: %w", err)
	}

	_, err = tx.Exec(`INSERT INTO masters_programs (master_id, program_id) VALUES
		(1, 1), (1, 2), (1, 3),
		(2, 1), (2, 2),
		(3, 3)`)
	if err != nil {
		return fmt.Errorf("insert masters_programs: %w", err)
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())
	if now.Hour() >= 19 {
		startOfDay = startOfDay.Add(24 * time.Hour)
	}

	for i := 0; i < 7; i++ {
		day := startOfDay.Add(time.Duration(i) * 24 * time.Hour)
		if day.Before(now) {
			continue
		}
		for _, slot := range []struct {
			masterID   int64
			programID  int64
			hour       int
			spots      int
			rental     bool
			rentalPrice float64
		}{
			{1, 1, 10, 6, false, 0},
			{1, 2, 13, 4, true, 500},
			{2, 1, 10, 6, false, 0},
			{2, 2, 14, 4, true, 500},
			{3, 3, 11, 8, false, 0},
		} {
			dt := day.Add(time.Duration(slot.hour) * time.Hour)
			endTime := dt.Add(2 * time.Hour)
			rentalAvail := 0
			if slot.rental {
				rentalAvail = 1
			}
			_, err = tx.Exec(
				`INSERT INTO slots (date_time, end_time, program_id, master_id, total_spots, available_spots, rental_available, rental_price)
				 VALUES ($1, $2, $3, $4, $5, $5, $6, $7)`,
				dt.Format(time.RFC3339), endTime.Format(time.RFC3339),
				slot.programID, slot.masterID, slot.spots, rentalAvail, slot.rentalPrice,
			)
			if err != nil {
				return fmt.Errorf("insert slot day=%d hour=%d: %w", i, slot.hour, err)
			}
		}
	}

	return tx.Commit()
}
