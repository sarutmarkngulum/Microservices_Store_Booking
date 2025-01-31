package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/final_project1240930/api_gateway/internal/logs"
	services "gitlab.com/final_project1240930/api_gateway/internal/services/booking"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

type bookingHandler struct {
	bookingSrv services.BookingService
}

func NewBookingHandler(bookingSrv services.BookingService) *bookingHandler {
	return &bookingHandler{bookingSrv: bookingSrv}
}

func createErrorResponse(err error) map[string]string {
	return map[string]string{"error": err.Error()}
}

func (h *bookingHandler) CreateBooking(c echo.Context) error {
	var req services.CreateBookingRequest

	if err := c.Bind(&req); err != nil {
		logs.Error("Invalid request format for CreateBooking", zap.Error(err))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid request format")))
	}

	resp, err := h.bookingSrv.CreateBooking(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to create booking", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *bookingHandler) UpdateBooking(c echo.Context) error {

	bookingID := c.Param("booking_id")
	if bookingID == "" {
		logs.Error("Booking ID is missing from URL", zap.String("booking_id", bookingID))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("booking ID is required")))
	}

	logs.Info("Received booking_id", zap.String("booking_id", bookingID))

	var req services.CreateBookingRequest

	req.BookingId = bookingID

	body := c.Request().Body

	data, err := io.ReadAll(body)
	if err != nil {
		logs.Error("Error reading request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("could not read request body")))
	}

	if err := protojson.Unmarshal(data, &req); err != nil {
		logs.Error("Error unmarshaling request data", zap.Error(err))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid request format")))
	}

	if req.BookingId == "" {
		logs.Error("Booking ID is empty after unmarshaling", zap.String("booking_id", req.BookingId))
		req.BookingId = bookingID
	}

	resp, err := h.bookingSrv.UpdateBooking(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to update booking", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *bookingHandler) DeleteBooking(c echo.Context) error {
	var req services.DeleteBookingRequest

	req.BookingId = c.Param("booking_id")

	if req.BookingId == "" {
		logs.Error("booking ID cannot be empty")
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("booking ID cannot be empty")))
	}

	resp, err := h.bookingSrv.DeleteBooking(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to delete booking", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *bookingHandler) GetBookingById(c echo.Context) error {

	id := c.Param("booking_id")
	req := services.GetBookingDetailsByIDRequest{BookingId: id}

	resp, err := h.bookingSrv.GetBookingDetailsByID(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to get booking", zap.String("bookingId", id), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *bookingHandler) GetBookings(c echo.Context) error {
	resp, err := h.bookingSrv.GetBookingDetails(c.Request().Context(), &services.GetBookingDetailsRequest{})
	if err != nil {
		logs.Error("Failed to get bookings", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}
