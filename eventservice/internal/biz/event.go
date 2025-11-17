package biz

import (
	"context"
	"fmt"
	"strings"
	"time"

	userv1 "userservice/api/userservice/v1"
	eventv1 "eventservice/api/eventservice/v1"
	"github.com/go-kratos/kratos/v2/log"
)

// ---------------- Helpers ----------------
// parseDate safely converts various date formats (with or without time) into time.Time
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


func toProto(ev *ShowEvent) *eventv1.ShowEvent {
	return &eventv1.ShowEvent{
		Id:             ev.ID,
		Title:          ev.Title,
		Description:    ev.Description,
		Date:           ev.Date.Format(time.RFC3339),
		TotalSeats:     ev.TotalSeats,
		AvailableSeats: ev.AvailableSeats,
		PricePerSeat:   ev.PricePerSeat,
	}
}

// ---------------- Entity ----------------
type ShowEvent struct {
	ID             uint64
	Title          string
	Description    string
	Date           time.Time
	TotalSeats     int32
	AvailableSeats int32
	PricePerSeat   float32
}

// ---------------- Repo Interface ----------------
type ShowEventRepo interface {
	Create(ctx context.Context, req *eventv1.CreateShowEventRequest) (*eventv1.ShowEvent, error)
	Get(ctx context.Context, id uint64) (*ShowEvent, error)
	List(ctx context.Context) ([]*ShowEvent, error)
	Update(ctx context.Context, ev *ShowEvent) (*ShowEvent, error)
	Delete(ctx context.Context, id uint64) (bool, error)
}

// ---------------- Usecase ----------------
type ShowEventUsecase struct {
	repo       ShowEventRepo
	userClient userv1.UserServiceClient
	log        *log.Helper
}

func NewShowEventUsecase(repo ShowEventRepo, userClient userv1.UserServiceClient, logger log.Logger) *ShowEventUsecase {
	return &ShowEventUsecase{
		repo:       repo,
		userClient: userClient,
		log:        log.NewHelper(logger),
	}
}

// Create event
func (uc *ShowEventUsecase) Create(ctx context.Context, req *eventv1.CreateShowEventRequest) (*ShowEvent, error) {
	protoEv, err := uc.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	ev := &ShowEvent{
		ID:             protoEv.Id,
		Title:          protoEv.Title,
		Description:    protoEv.Description,
		Date:           parseDate(protoEv.Date),
		TotalSeats:     protoEv.TotalSeats,
		AvailableSeats: protoEv.AvailableSeats,
		PricePerSeat:   protoEv.PricePerSeat,
	}

	return ev, nil
}

// Get event
func (uc *ShowEventUsecase) Get(ctx context.Context, id uint64) (*ShowEvent, error) {
	return uc.repo.Get(ctx, id)
}

// List events
func (uc *ShowEventUsecase) List(ctx context.Context) ([]*ShowEvent, error) {
	return uc.repo.List(ctx)
}

// Update event
func (uc *ShowEventUsecase) Update(ctx context.Context, ev *ShowEvent) (*ShowEvent, error) {
	return uc.repo.Update(ctx, ev)
}

// Delete event
func (uc *ShowEventUsecase) Delete(ctx context.Context, id uint64) (bool, error) {
	return uc.repo.Delete(ctx, id)
}

// ---------------- Seat Management -------------------

func (uc *ShowEventUsecase) DecrementSeats(ctx context.Context, eventID uint64, seatIDs []string) error {
	ev, err := uc.repo.Get(ctx, eventID)
	if err != nil {
		return fmt.Errorf("event not found: %w", err)
	}

	if ev.AvailableSeats < int32(len(seatIDs)) {
		return fmt.Errorf("not enough seats available")
	}

	ev.AvailableSeats -= int32(len(seatIDs))
	_, err = uc.repo.Update(ctx, ev)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	uc.log.Infof("EventID=%d: seats decremented, AvailableSeats=%d", eventID, ev.AvailableSeats)
	return nil
}

func (uc *ShowEventUsecase) IncrementSeats(ctx context.Context, eventID uint64, seatIDs []string) error {
	ev, err := uc.repo.Get(ctx, eventID)
	if err != nil {
		return fmt.Errorf("event not found: %w", err)
	}

	ev.AvailableSeats += int32(len(seatIDs))
	if ev.AvailableSeats > ev.TotalSeats {
		ev.AvailableSeats = ev.TotalSeats
	}

	_, err = uc.repo.Update(ctx, ev)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	uc.log.Infof("EventID=%d: seats incremented, AvailableSeats=%d", eventID, ev.AvailableSeats)
	return nil
}
