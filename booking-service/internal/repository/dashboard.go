package repository

import "context"

type GetDailySummaryResponse struct {
	DailySales     float64 `json:"daily_sales"`
	DailyBookings  int32   `json:"daily_bookings"`
	DailyCustomers int32   `json:"daily_customers"`
	TotalUsers     int32   `json:"total_users"`
}

type GetMonthlySalesResponse struct {
	Sales []MonthlySales `json:"sales"`
}

type MonthlySales struct {
	Month      int32   `json:"month"`
	TotalSales float64 `json:"total_sales"`
}

type GetMonthlyBookingAndCustomersResponse struct {
	Data []MonthlyBookingAndCustomers `json:"data"`
}

type MonthlyBookingAndCustomers struct {
	Month          int32 `json:"month"`
	TotalBookings  int32 `json:"total_bookings"`
	TotalCustomers int32 `json:"total_customers"`
}

type GetBestSellersResponse struct {
	TopMenuSets []MenuSets `json:"top_menu_sets"`
	TopAlaCarte []Menu     `json:"top_a_la_carte"`
}

type MenuSets struct {
	MenuSetName       string `json:"menu_set_name"`
	TotalQuantitySold int32  `json:"total_quantity_sold"`
}
type Menu struct {
	NameTh            string `json:"name_th"`
	NameEn            string `json:"name_en"`
	ImageUrl          string `json:"image_url"`
	TotalQuantitySold int32  `json:"total_quantity_sold"`
}

type DashboardRepository interface {
	GetDailySales(ctx context.Context) (float64, error)
	GetDailyBookings(ctx context.Context) (int32, error)
	GetDailyCustomers(ctx context.Context) (int32, error)
	GetTotalUsers(ctx context.Context) (int32, error)

	GetMonthlySales(ctx context.Context) (GetMonthlySalesResponse, error)
	GetMonthlyBookingAndCustomers(ctx context.Context) (GetMonthlyBookingAndCustomersResponse, error)
	GetBestSellers(ctx context.Context) (GetBestSellersResponse, error)
}
