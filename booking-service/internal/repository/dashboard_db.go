package repository

import (
	"context"
	"fmt"

	"gitlab.com/final_project1240930/booking_service/internal/logs"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Gorm implementation of DashboardRepository
type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

// GetDailySales returns the daily sales for today
func (r *dashboardRepository) GetDailySales(ctx context.Context) (float64, error) {
	var dailySales float64
	err := r.db.Raw(`
		SELECT COALESCE(SUM(total_price), 0) AS daily_sales
		FROM bookings
		WHERE DATE(booking_date_time) = CURRENT_DATE
	`).Scan(&dailySales).Error
	if err != nil {
		logs.Error("Failed to get daily sales", zap.Error(err))
		return 0, fmt.Errorf("failed to get daily sales: %w", err)
	}
	return dailySales, nil
}

// GetDailyBookings returns the daily bookings count for today
func (r *dashboardRepository) GetDailyBookings(ctx context.Context) (int32, error) {
	var dailyBookings int32
	err := r.db.Raw(`
		SELECT COALESCE(COUNT(*), 0) AS daily_bookings
		FROM bookings
		WHERE DATE(booking_date_time) = CURRENT_DATE
	`).Scan(&dailyBookings).Error
	if err != nil {
		logs.Error("Failed to get daily bookings", zap.Error(err))
		return 0, fmt.Errorf("failed to get daily bookings: %w", err)
	}
	return dailyBookings, nil
}

// GetDailyCustomers returns the total number of daily customers (children + adults) for today
func (r *dashboardRepository) GetDailyCustomers(ctx context.Context) (int32, error) {
	var dailyCustomers int32
	err := r.db.Raw(`
		SELECT COALESCE(SUM(num_children + num_adults), 0) AS daily_customers
		FROM bookings
		WHERE DATE(booking_date_time) = CURRENT_DATE
	`).Scan(&dailyCustomers).Error
	if err != nil {
		logs.Error("Failed to get daily customers", zap.Error(err))
		return 0, fmt.Errorf("failed to get daily customers: %w", err)
	}
	return dailyCustomers, nil
}

// GetTotalUsers returns the total user count
func (r *dashboardRepository) GetTotalUsers(ctx context.Context) (int32, error) {
	var totalUsers int32
	err := r.db.Raw(`
		SELECT COUNT(*) AS total_users
		FROM users
	`).Scan(&totalUsers).Error
	if err != nil {
		logs.Error("Failed to get total users", zap.Error(err))
		return 0, fmt.Errorf("failed to get total users: %w", err)
	}
	return totalUsers, nil
}

// GetMonthlySales returns monthly sales for the current year
func (r *dashboardRepository) GetMonthlySales(ctx context.Context) (GetMonthlySalesResponse, error) {
	var response GetMonthlySalesResponse
	err := r.db.Raw(`
		WITH months AS (
			SELECT generate_series(1, 12) AS month
		)
		SELECT 
			months.month,
			COALESCE(SUM(b.total_price), 0) AS total_sales
		FROM months
		LEFT JOIN bookings b
			ON EXTRACT(MONTH FROM b.booking_date_time) = months.month
			AND EXTRACT(YEAR FROM b.booking_date_time) = EXTRACT(YEAR FROM CURRENT_DATE)
		GROUP BY months.month
		ORDER BY months.month
	`).Scan(&response.Sales).Error
	if err != nil {
		logs.Error("Failed to get monthly sales", zap.Error(err))
		return response, fmt.Errorf("failed to get monthly sales: %w", err)
	}
	logs.Info("Successfully fetched monthly sales")
	return response, nil
}

// GetMonthlyBookingAndCustomers returns monthly bookings and customer counts for the current year, showing all 12 months
func (r *dashboardRepository) GetMonthlyBookingAndCustomers(ctx context.Context) (GetMonthlyBookingAndCustomersResponse, error) {
	var response GetMonthlyBookingAndCustomersResponse
	err := r.db.Raw(`
		WITH months AS (
			SELECT generate_series(1, 12) AS month
		)
		SELECT 
			months.month,
			COALESCE(COUNT(b.booking_date_time), 0) AS total_bookings,
			COALESCE(SUM(b.num_children + b.num_adults), 0) AS total_customers
		FROM months
		LEFT JOIN bookings b
			ON EXTRACT(MONTH FROM b.booking_date_time) = months.month
			AND EXTRACT(YEAR FROM b.booking_date_time) = EXTRACT(YEAR FROM CURRENT_DATE)
		GROUP BY months.month
		ORDER BY months.month
	`).Scan(&response.Data).Error
	if err != nil {
		logs.Error("Failed to get monthly bookings and customers", zap.Error(err))
		return response, fmt.Errorf("failed to get monthly bookings and customers: %w", err)
	}
	logs.Info("Successfully fetched monthly bookings and customers")
	return response, nil
}

// GetBestSellers returns the best selling menu sets and a la carte items
func (r *dashboardRepository) GetBestSellers(ctx context.Context) (GetBestSellersResponse, error) {
	var response GetBestSellersResponse

	// Menu Sets Best Sellers
	err := r.db.Raw(`
		SELECT 
			ms.name AS menu_set_name,
			SUM(bms.quantity) AS total_quantity_sold
		FROM booking_menu_sets bms
		JOIN menu_sets ms ON bms.menu_set_id = ms.uuid
		GROUP BY ms.name
		ORDER BY total_quantity_sold DESC
		LIMIT 5
	`).Scan(&response.TopMenuSets).Error
	if err != nil {
		logs.Error("Failed to get best selling menu sets", zap.Error(err))
		return response, fmt.Errorf("failed to get best selling menu sets: %w", err)
	}

	// A La Carte Best Sellers
	err = r.db.Raw(`
		SELECT 
			mi.name_th,
			mi.name_en,
			mi.image_url,
			SUM(bmi.quantity) AS total_quantity_sold
		FROM booking_menu_items bmi
		JOIN menu_items mi ON bmi.menu_item_id = mi.uuid
		GROUP BY mi.name_th, mi.name_en, mi.image_url
		ORDER BY total_quantity_sold DESC
		LIMIT 5
	`).Scan(&response.TopAlaCarte).Error
	if err != nil {
		logs.Error("Failed to get best selling a la carte", zap.Error(err))
		return response, fmt.Errorf("failed to get best selling a la carte: %w", err)
	}

	logs.Info("Successfully fetched best sellers")
	return response, nil
}
