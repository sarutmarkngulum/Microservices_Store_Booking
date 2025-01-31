package services

import (
	"context"
	"time"

	"gitlab.com/final_project1240930/api_gateway/internal/logs"
	"go.uber.org/zap"
)

type BookingService interface {
	CreateBooking(ctx context.Context, req *CreateBookingRequest) (*CreateBookingResponse, error)
	UpdateBooking(ctx context.Context, req *CreateBookingRequest) (*UpdateBookingResponse, error)
	DeleteBooking(ctx context.Context, req *DeleteBookingRequest) (*DeleteBookingResponse, error)
	GetBookingDetails(ctx context.Context, req *GetBookingDetailsRequest) (*GetBookingDetailsResponse, error)
	GetBookingDetailsByID(ctx context.Context, req *GetBookingDetailsByIDRequest) (*GetBookingDetailsByIDResponse, error)
}
type bookingService struct {
	bookingClient BookingServiceClient
}

func NewBookingService(bookingClient BookingServiceClient) BookingService {
	return &bookingService{bookingClient: bookingClient}
}

// Handle timeout logic
func (s *bookingService) createWithTimeout(ctx context.Context, call func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	// ตั้งเวลา timeout 5 วินาที
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	res, err := call(ctx)
	if err != nil {
		logs.Error("Error: %v", zap.Error(err))
		return nil, err
	}
	return res, nil
}

// Business logic for creating a booking
func (s *bookingService) CreateBooking(ctx context.Context, req *CreateBookingRequest) (*CreateBookingResponse, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.bookingClient.CreateBooking(ctx, req)
	})
	if res != nil {
		return res.(*CreateBookingResponse), nil
	}
	return nil, err
}

// Business logic for updating a booking
func (s *bookingService) UpdateBooking(ctx context.Context, req *CreateBookingRequest) (*UpdateBookingResponse, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.bookingClient.UpdateBooking(ctx, req)
	})
	if res != nil {
		return res.(*UpdateBookingResponse), nil
	}
	return nil, err
}

// Business logic for deleting a booking
func (s *bookingService) DeleteBooking(ctx context.Context, req *DeleteBookingRequest) (*DeleteBookingResponse, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.bookingClient.DeleteBooking(ctx, req)
	})
	if res != nil {
		return res.(*DeleteBookingResponse), nil
	}
	return nil, err
}

// Business logic for retrieving bookings
func (s *bookingService) GetBookingDetails(ctx context.Context, req *GetBookingDetailsRequest) (*GetBookingDetailsResponse, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.bookingClient.GetBookingDetails(ctx, req)
	})
	if res != nil {
		return res.(*GetBookingDetailsResponse), nil
	}
	return nil, err
}

func (s *bookingService) GetBookingDetailsByID(ctx context.Context, req *GetBookingDetailsByIDRequest) (*GetBookingDetailsByIDResponse, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.bookingClient.GetBookingDetailsByID(ctx, req)
	})
	if res != nil {
		return res.(*GetBookingDetailsByIDResponse), nil
	}
	return nil, err
}
