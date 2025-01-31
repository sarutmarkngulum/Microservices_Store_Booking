package services

import (
	"context"

	"gitlab.com/final_project1240930/booking_service/internal/logs"
	"gitlab.com/final_project1240930/booking_service/internal/repository"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type dashboardServer struct {
	repo repository.DashboardRepository
}

func NewDashboardServer(repo repository.DashboardRepository) DashboardServiceServer {
	return &dashboardServer{repo: repo}
}

func (s *dashboardServer) mustEmbedUnimplementedDashboardServiceServer() {}

func (s *dashboardServer) GetDailySummary(ctx context.Context, req *emptypb.Empty) (*GetDailySummaryResponse, error) {
	logs.Info("Received GetDailySummaryRequest")

	dailySales, err := s.repo.GetDailySales(ctx)
	if err != nil {
		logs.Error("Failed to get daily sales", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to get daily sales: %v", err)
	}

	dailyBookings, err := s.repo.GetDailyBookings(ctx)
	if err != nil {
		logs.Error("Failed to get daily bookings", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to get daily bookings: %v", err)
	}

	dailyCustomers, err := s.repo.GetDailyCustomers(ctx)
	if err != nil {
		logs.Error("Failed to get daily customers", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to get daily customers: %v", err)
	}

	totalUsers, err := s.repo.GetTotalUsers(ctx)
	if err != nil {
		logs.Error("Failed to get total users", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to get total users: %v", err)
	}

	// ส่งค่าผลลัพธ์
	return &GetDailySummaryResponse{
		DailySales:     dailySales,
		DailyBookings:  dailyBookings,
		DailyCustomers: dailyCustomers,
		TotalUsers:     totalUsers,
	}, nil
}

func (s *dashboardServer) GetMonthlySales(ctx context.Context, req *emptypb.Empty) (*GetMonthlySalesResponse, error) {
	logs.Info("Received GetMonthlySalesRequest")

	data, err := s.repo.GetMonthlySales(ctx)
	if err != nil {
		logs.Error("Failed to fetch monthly sales", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to fetch monthly sales: %v", err)
	}

	// Map the data to protobuf response
	var salesList []*MonthlySales
	for _, sale := range data.Sales {
		salesList = append(salesList, &MonthlySales{
			Month:      sale.Month,
			TotalSales: sale.TotalSales,
		})
	}

	return &GetMonthlySalesResponse{
		Sales: salesList,
	}, nil
}

func (s *dashboardServer) GetMonthlyBookingAndCustomers(ctx context.Context, req *emptypb.Empty) (*GetMonthlyBookingAndCustomersResponse, error) {
	logs.Info("Received GetMonthlyBookingAndCustomersRequest")

	data, err := s.repo.GetMonthlyBookingAndCustomers(ctx)
	if err != nil {
		logs.Error("Failed to fetch monthly booking and customers", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to fetch monthly booking and customers: %v", err)
	}

	// Map the data to protobuf response
	var bookingAndCustomersList []*MonthlyBookingAndCustomers
	for _, item := range data.Data {
		bookingAndCustomersList = append(bookingAndCustomersList, &MonthlyBookingAndCustomers{
			Month:          item.Month,
			TotalBookings:  item.TotalBookings,
			TotalCustomers: item.TotalCustomers,
		})
	}

	return &GetMonthlyBookingAndCustomersResponse{
		Data: bookingAndCustomersList,
	}, nil
}

func (s *dashboardServer) GetBestSellers(ctx context.Context, req *emptypb.Empty) (*GetBestSellersResponse, error) {
	logs.Info("Received GetBestSellersRequest")

	data, err := s.repo.GetBestSellers(ctx)
	if err != nil {
		logs.Error("Failed to fetch best sellers", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to fetch best sellers: %v", err)
	}

	var topMenuSetsList []*MenuSets
	for _, menu := range data.TopMenuSets {
		topMenuSetsList = append(topMenuSetsList, &MenuSets{
			NameSetName:       menu.MenuSetName,
			TotalQuantitySold: menu.TotalQuantitySold,
		})
	}

	var topAlaCarteList []*Menu
	for _, menu := range data.TopAlaCarte {
		topAlaCarteList = append(topAlaCarteList, &Menu{
			NameTh:            menu.NameTh,
			NameEn:            menu.NameEn,
			ImageUrl:          menu.ImageUrl,
			TotalQuantitySold: menu.TotalQuantitySold,
		})
	}

	return &GetBestSellersResponse{
		TopMenuSets: topMenuSetsList,
		TopALaCarte: topAlaCarteList,
	}, nil
}
