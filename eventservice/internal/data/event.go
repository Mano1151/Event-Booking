package data

import (
	"context"
	"fmt"
	"strings"
	"time"

	eventv1 "eventservice/api/eventservice/v1"
	"eventservice/internal/biz"
	"gorm.io/gorm"
)

// parseDate with normalization
func parseDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Now()
	}

	dateStr = strings.TrimSpace(dateStr)

	// Normalize single-digit month/day
	parts := strings.Split(dateStr, "-")
	if len(parts) == 3 {
		for i := 1; i <= 2; i++ {
			if len(parts[i]) == 1 {
				parts[i] = "0" + parts[i]
			}
		}
		dateStr = strings.Join(parts, "-")
	}

	// Use local timezone
	loc := time.Local

	layouts := []string{
		time.RFC3339,       // 2025-12-07T18:30:00Z
		"2006-01-02 15:04", // 2025-12-07 18:30
		"2006-01-02T15:04", // 2025-12-07T18:30
		"2006-01-02",       // 2025-12-07
	}

	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, dateStr, loc); err == nil {
			return t
		}
	}

	// fallback
	return time.Now()
}


type ShowEvent struct {
	ID             uint64    `gorm:"primaryKey"`
	Title          string
	Description    string
	Date           time.Time
	TotalSeats     int32
	AvailableSeats int32 `gorm:"column:available_seats"`
	PricePerSeat   float32
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type showEventRepo struct {
	db *gorm.DB
}

func NewShowEventRepo(db *gorm.DB) biz.ShowEventRepo {
	db.AutoMigrate(&ShowEvent{})
	return &showEventRepo{db: db}
}

func (r *showEventRepo) Create(ctx context.Context, req *eventv1.CreateShowEventRequest) (*eventv1.ShowEvent, error) {
	ev := &ShowEvent{
		Title:          req.Title,
		Description:    req.Description,
		Date:           parseDate(req.Date),
		TotalSeats:     req.TotalSeats,
		AvailableSeats: req.TotalSeats,
		PricePerSeat:   req.PricePerSeat,
	}

	if err := r.db.WithContext(ctx).Create(ev).Error; err != nil {
		return nil, err
	}

	return &eventv1.ShowEvent{
		Id:             ev.ID,
		Title:          ev.Title,
		Description:    ev.Description,
		Date:           ev.Date.Format(time.RFC3339),
		TotalSeats:     ev.TotalSeats,
		AvailableSeats: ev.AvailableSeats,
		PricePerSeat:   ev.PricePerSeat,
	}, nil
}

func (r *showEventRepo) Get(ctx context.Context, id uint64) (*biz.ShowEvent, error) {
	var model ShowEvent
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return &biz.ShowEvent{
		ID:             model.ID,
		Title:          model.Title,
		Description:    model.Description,
		Date:           model.Date,
		TotalSeats:     model.TotalSeats,
		AvailableSeats: model.AvailableSeats,
		PricePerSeat:   model.PricePerSeat,
	}, nil
}

func (r *showEventRepo) List(ctx context.Context) ([]*biz.ShowEvent, error) {
	var models []ShowEvent
	fmt.Println("----->> Fetching ShowEvents")
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	res := make([]*biz.ShowEvent, 0, len(models))
	for _, m := range models {
		res = append(res, &biz.ShowEvent{
			ID:             m.ID,
			Title:          m.Title,
			Description:    m.Description,
			Date:           m.Date,
			TotalSeats:     m.TotalSeats,
			AvailableSeats: m.AvailableSeats,
			PricePerSeat:   m.PricePerSeat,
		})
	}
	return res, nil
}

func (r *showEventRepo) Update(ctx context.Context, ev *biz.ShowEvent) (*biz.ShowEvent, error) {
	var dbEv ShowEvent
	if err := r.db.First(&dbEv, ev.ID).Error; err != nil {
		return nil, err
	}

	dbEv.AvailableSeats = ev.AvailableSeats
	fmt.Printf("Updating DB: EventID=%d NewAvailable=%d\n", ev.ID, ev.AvailableSeats)

	if err := r.db.Save(&dbEv).Error; err != nil {
		return nil, err
	}

	return &biz.ShowEvent{
		ID:             dbEv.ID,
		Title:          dbEv.Title,
		Description:    dbEv.Description,
		Date:           dbEv.Date,
		TotalSeats:     dbEv.TotalSeats,
		AvailableSeats: dbEv.AvailableSeats,
		PricePerSeat:   dbEv.PricePerSeat,
	}, nil
}

func (r *showEventRepo) Delete(ctx context.Context, id uint64) (bool, error) {
	if err := r.db.WithContext(ctx).Delete(&ShowEvent{}, id).Error; err != nil {
		return false, err
	}
	return true, nil
}
