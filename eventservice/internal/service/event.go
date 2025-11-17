package service

import (
	"context"
	"fmt"
	"time"

	v1 "eventservice/api/eventservice/v1"
	"eventservice/internal/biz"
	userv1 "userservice/api/userservice/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ---------------- ShowEventService ----------------
type ShowEventService struct {
	v1.UnimplementedEventServiceServer
	uc         *biz.ShowEventUsecase
	userClient userv1.UserServiceClient
}

func NewShowEventService(uc *biz.ShowEventUsecase, userClient userv1.UserServiceClient) *ShowEventService {
	return &ShowEventService{uc: uc, userClient: userClient}
}

// Create ShowEvent
func (s *ShowEventService) CreateShowEvent(ctx context.Context, req *v1.CreateShowEventRequest) (*v1.ShowEventReply, error) {
	createdEv, err := s.uc.Create(ctx, req) // pass proto directly
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create event: %v", err)
	}

	return &v1.ShowEventReply{
		ShowEvent: &v1.ShowEvent{
			Id:             createdEv.ID,
			Title:          createdEv.Title,
			Description:    createdEv.Description,
			Date:           createdEv.Date.Format(time.RFC3339),
			TotalSeats:     createdEv.TotalSeats,
			AvailableSeats: createdEv.AvailableSeats,
			PricePerSeat:   createdEv.PricePerSeat,
		},
	}, nil
}

// Get ShowEvent
func (s *ShowEventService) GetShowEvent(ctx context.Context, req *v1.GetShowEventRequest) (*v1.ShowEventReply, error) {
	ev, err := s.uc.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "event not found: %v", err)
	}
	return &v1.ShowEventReply{ShowEvent: toProto(ev)}, nil
}

// List ShowEvents
func (s *ShowEventService) ListShowEvents(ctx context.Context, req *v1.ListShowEventsRequest) (*v1.ListShowEventsReply, error) {
	fmt.Println("----->>showwvwnts-biz")
	evs, err := s.uc.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list events: %v", err)
	}

	result := make([]*v1.ShowEvent, 0, len(evs))
	for _, e := range evs {
		result = append(result, toProto(e))
	}

	return &v1.ListShowEventsReply{ShowEvents: result}, nil
}

// Update ShowEvent
func (s *ShowEventService) UpdateShowEvent(ctx context.Context, req *v1.UpdateShowEventRequest) (*v1.ShowEventReply, error) {
	ev := &biz.ShowEvent{
		ID:             req.Id,
		Title:          req.Title,
		Description:    req.Description,
		Date:           parseDate(req.Date),
		TotalSeats:     req.TotalSeats,
		AvailableSeats: req.AvailableSeats,
		PricePerSeat:   req.PricePerSeat,
	}

	updatedEv, err := s.uc.Update(ctx, ev)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update event: %v", err)
	}

	return &v1.ShowEventReply{ShowEvent: toProto(updatedEv)}, nil
}

// Delete ShowEvent
func (s *ShowEventService) DeleteShowEvent(ctx context.Context, req *v1.DeleteShowEventRequest) (*v1.DeleteShowEventReply, error) {
	ok, err := s.uc.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete event: %v", err)
	}
	return &v1.DeleteShowEventReply{Success: ok}, nil
}

// DecrementSeats - called from BookingService when booking is CONFIRMED
func (s *ShowEventService) DecrementSeats(ctx context.Context, req *v1.DecrementSeatsRequest) (*v1.DecrementSeatsReply, error) {
	err := s.uc.DecrementSeats(ctx, req.EventId, req.SeatIds) // SeatIds is now []string
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to decrement seats: %v", err)
	}

	ev, err := s.uc.Get(ctx, req.EventId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch event: %v", err)
	}

	return &v1.DecrementSeatsReply{
		Success:        true,
		AvailableSeats: ev.AvailableSeats,
	}, nil
}

func (s *ShowEventService) IncrementSeats(ctx context.Context, req *v1.IncrementSeatsRequest) (*v1.IncrementSeatsReply, error) {
	err := s.uc.IncrementSeats(ctx, req.EventId, req.SeatIds) // SeatIds is now []string
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to increment seats: %v", err)
	}

	ev, err := s.uc.Get(ctx, req.EventId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch event: %v", err)
	}

	return &v1.IncrementSeatsReply{
		Success:        true,
		AvailableSeats: ev.AvailableSeats,
	}, nil
}

// Helpers
func parseDate(dateStr string) time.Time {
    if dateStr == "" {
        return time.Now()
    }

    // Try ISO full format first (e.g., "2025-12-07T18:30:00Z" or with offset)
    if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
        return t
    }

    // Try "YYYY-MM-DD HH:MM" (common frontend format)
    if t, err := time.Parse("2006-01-02 15:04", dateStr); err == nil {
        return t
    }

    // Try "YYYY-MM-DD" (date only)
    if t, err := time.Parse("2006-01-02", dateStr); err == nil {
        return t
    }

    // Fallback: return current time if parsing fails
    return time.Now()
}

func toProto(ev *biz.ShowEvent) *v1.ShowEvent {
	return &v1.ShowEvent{
		Id:             ev.ID,
		Title:          ev.Title,
		Description:    ev.Description,
		Date:           ev.Date.Format(time.RFC3339),
		TotalSeats:     ev.TotalSeats,
		AvailableSeats: ev.AvailableSeats,
		PricePerSeat:   ev.PricePerSeat,
	}
}

func (s *ShowEventService) ValidateUser(ctx context.Context, req *v1.ValidateUserRequest) (*v1.ValidateUserReply, error) {
	// Call UserService
	resp, err := s.userClient.GetUser(ctx, &userv1.GetUserRequest{Id: req.Id})
	if err != nil {
		return &v1.ValidateUserReply{Found: false}, nil
	}
	if resp.Id == 0 {
		return &v1.ValidateUserReply{Found: false}, nil
	}
	return &v1.ValidateUserReply{Found: true}, nil
}
