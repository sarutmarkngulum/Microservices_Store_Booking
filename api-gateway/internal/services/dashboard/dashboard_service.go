package services

import (
	context "context"
	"time"

	"gitlab.com/final_project1240930/api_gateway/internal/logs"
	"go.uber.org/zap"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type DashBoardService interface {
	GetDailySummary(ctx context.Context, req *emptypb.Empty) (*GetDailySummaryResponse, error)
	GetMonthlySales(ctx context.Context, req *emptypb.Empty) (*GetMonthlySalesResponse, error)
	GetMonthlyBookingAndCustomers(ctx context.Context, req *emptypb.Empty) (*GetMonthlyBookingAndCustomersResponse, error)
	GetBestSellers(ctx context.Context, req *emptypb.Empty) (*GetBestSellersResponse, error)
}

type dashboardService struct {
	dashboardClient DashboardServiceClient
}

func NewDashboardService(dashboardClient DashboardServiceClient) DashBoardService {
	return &dashboardService{dashboardClient: dashboardClient}
}

func (s *dashboardService) callWithTimeout(ctx context.Context, call func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5) // Set timeout
	defer cancel()

	res, err := call(ctx)
	if err != nil {
		logs.Error("Error calling service", zap.Error(err))
		return nil, err
	}
	return res, nil
}

func (s *dashboardService) GetDailySummary(ctx context.Context, req *emptypb.Empty) (*GetDailySummaryResponse, error) {
	res, err := s.callWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.dashboardClient.GetDailySummary(ctx, req)
	})
	if res != nil {
		return res.(*GetDailySummaryResponse), nil
	}
	return nil, err
}

func (s *dashboardService) GetMonthlySales(ctx context.Context, req *emptypb.Empty) (*GetMonthlySalesResponse, error) {
	res, err := s.callWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.dashboardClient.GetMonthlySales(ctx, req)
	})
	if res != nil {
		return res.(*GetMonthlySalesResponse), nil
	}
	return nil, err
}
func (s *dashboardService) GetMonthlyBookingAndCustomers(ctx context.Context, req *emptypb.Empty) (*GetMonthlyBookingAndCustomersResponse, error) {
	res, err := s.callWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.dashboardClient.GetMonthlyBookingAndCustomers(ctx, req)
	})
	if res != nil {
		return res.(*GetMonthlyBookingAndCustomersResponse), nil
	}
	return nil, err
}

func (s *dashboardService) GetBestSellers(ctx context.Context, req *emptypb.Empty) (*GetBestSellersResponse, error) {
	res, err := s.callWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.dashboardClient.GetBestSellers(ctx, req)
	})
	if res != nil {
		return res.(*GetBestSellersResponse), nil
	}
	return nil, err
}
