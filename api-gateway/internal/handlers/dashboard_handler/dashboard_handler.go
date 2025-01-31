package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/final_project1240930/api_gateway/internal/logs"
	services "gitlab.com/final_project1240930/api_gateway/internal/services/dashboard"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type dashboardHandler struct {
	dashboardSrv services.DashBoardService
}

func NewDashboardHandler(dashboardSrv services.DashBoardService) *dashboardHandler {
	return &dashboardHandler{dashboardSrv: dashboardSrv}
}

func createErrorResponse(err error) map[string]string {
	return map[string]string{"error": err.Error()}
}

func (h *dashboardHandler) GetDailySummary(c echo.Context) error {
	resp, err := h.dashboardSrv.GetDailySummary(c.Request().Context(), &emptypb.Empty{})
	if err != nil {
		logs.Error("Failed to get tables", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *dashboardHandler) GetMonthlySales(c echo.Context) error {
	resp, err := h.dashboardSrv.GetMonthlySales(c.Request().Context(), &emptypb.Empty{})
	if err != nil {
		logs.Error("Failed to get monthly sales", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	for i := range resp.Sales {
		if resp.Sales[i].TotalSales == 0 {
			resp.Sales[i].TotalSales = 0
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *dashboardHandler) GetMonthlyBookingAndCustomers(c echo.Context) error {
	resp, err := h.dashboardSrv.GetMonthlyBookingAndCustomers(c.Request().Context(), &emptypb.Empty{})
	if err != nil {
		logs.Error("Failed to get monthly bookings and customers", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	for i := range resp.Data {
		if resp.Data[i].TotalBookings == 0 {
			resp.Data[i].TotalBookings = 0
		}
		if resp.Data[i].TotalCustomers == 0 {
			resp.Data[i].TotalCustomers = 0
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *dashboardHandler) GetBestSellers(c echo.Context) error {
	resp, err := h.dashboardSrv.GetBestSellers(c.Request().Context(), &emptypb.Empty{})
	if err != nil {
		logs.Error("Failed to get tables", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}
