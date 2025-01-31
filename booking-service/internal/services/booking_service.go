package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gitlab.com/final_project1240930/booking_service/internal/repository"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type bookingServer struct {
	bookingRepo repository.BookingRepository
}

func NewBookingServer(bookingRepo repository.BookingRepository) BookingServiceServer {
	return &bookingServer{bookingRepo: bookingRepo}
}

func (s *bookingServer) mustEmbedUnimplementedBookingServiceServer() {

}

func convertToProto(bookings []repository.Booking) []*BookingDetail {
	var protoBookings []*BookingDetail

	for _, booking := range bookings {

		protoBooking := &BookingDetail{
			BookingId:       booking.BookingID,
			CustomerName:    booking.CustomerName,
			CompanyName:     booking.CompanyName,
			BookingDateTime: booking.BookingDateTime.Format(time.RFC3339),
			PhoneNumber:     booking.PhoneNumber,
			NumChildren:     booking.NumChildren,
			NumAdults:       booking.NumAdults,
			NumTables:       booking.NumTables,
			Status:          booking.Status,
			TotalPrice:      booking.TotalPrice,
			// แปลงข้อมูล
			Tables:    convertTablesToProto(booking.Tables),
			MenuSets:  convertMenuSetsToProto(booking.BookingMenuSets),
			MenuItems: convertMenuItemsToProto(booking.BookingMenuItems),
		}

		// เพิ่ม BookingDetail ลงใน slice
		protoBookings = append(protoBookings, protoBooking)
	}

	return protoBookings
}
func convertTablesToProto(tables []repository.BookingTable) []*BookingTable {
	var protoTables []*BookingTable
	for _, table := range tables {
		protoTable := &BookingTable{
			TableId:     table.TableID,
			TableNumber: table.TableNumber,
			Type:        table.Type,
			SeatCount:   table.SeatCount,
		}
		protoTables = append(protoTables, protoTable)
	}
	return protoTables
}

func convertMenuSetsToProto(menuSets []repository.BookingMenuSet) []*BookingMenuSet {
	var protoMenuSets []*BookingMenuSet
	for _, menuSet := range menuSets {
		protoMenuSet := &BookingMenuSet{
			MenuSetName:  menuSet.MenuSetName,
			MenuSetPrice: menuSet.MenuSetPrice,
			Quantity:     menuSet.Quantity,
			MenuItems:    convertMenuItemsToProto(menuSet.MenuItems),
		}
		protoMenuSets = append(protoMenuSets, protoMenuSet)
	}
	return protoMenuSets
}

func convertMenuItemsToProto(menuItems []repository.BookingMenuItem) []*BookingMenuItem {
	var protoMenuItems []*BookingMenuItem
	for _, menuItem := range menuItems {
		protoMenuItem := &BookingMenuItem{
			MenuItemId:  menuItem.MenuItemID,
			NameTh:      menuItem.NameTh,
			NameEn:      menuItem.NameEn,
			Description: menuItem.Description,
			Price:       menuItem.Price,
			Category:    menuItem.Category,
			ImageUrl:    menuItem.ImageURL,
			Quantity:    menuItem.Quantity,
		}
		protoMenuItems = append(protoMenuItems, protoMenuItem)
	}
	return protoMenuItems
}

func (s *bookingServer) GetBookingDetails(ctx context.Context, req *GetBookingDetailsRequest) (*GetBookingDetailsResponse, error) {
	// ดึงข้อมูลจาก repository
	bookings, err := s.bookingRepo.GetBookingDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to fetch booking details")
	}

	// ตรวจสอบว่ามีข้อมูลหรือไม่
	if len(bookings) == 0 {
		return nil, status.Error(codes.NotFound, "No bookings found")
	}

	// แปลงข้อมูลจาก Booking เป็น Protobuf
	protoBookings := convertToProto(bookings)

	return &GetBookingDetailsResponse{
		BookingDetails: protoBookings,
	}, nil
}

func (s *bookingServer) GetBookingDetailsByID(ctx context.Context, req *GetBookingDetailsByIDRequest) (*GetBookingDetailsByIDResponse, error) {
	// ดึงข้อมูลจาก repository
	booking, err := s.bookingRepo.GetBookingDetailsByID(ctx, req.BookingId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to fetch booking details")
	}

	// แปลงข้อมูลจาก Booking เป็น Protobuf
	protoBookings := convertToProto([]repository.Booking{*booking})

	return &GetBookingDetailsByIDResponse{
		BookingDetail: protoBookings[0],
	}, nil
}

// แปลง gRPC request เป็น repository request
func ConvertCreateBookingRequestToRepositoryRequest(req *CreateBookingRequest, bookingDateTime time.Time) *repository.CreateBookingRequest {
	// Process table IDs
	tables := make([]repository.CreateBookingTable, len(req.TableIds))
	for i, table := range req.TableIds {
		tables[i] = repository.CreateBookingTable{
			TableID: table,
		}
	}

	menuSets := make([]repository.CreateBookingMenuSet, len(req.MenuSets))
	for i, menuSet := range req.MenuSets {
		menuSets[i] = repository.CreateBookingMenuSet{
			MenuSetID: menuSet.MenuSetId,
			Quantity:  menuSet.Quantity,
		}
	}

	menuItems := make([]repository.CreateBookingMenuItem, len(req.MenuItems))
	for i, menuItem := range req.MenuItems {
		menuItems[i] = repository.CreateBookingMenuItem{
			MenuItemID: menuItem.MenuItemId,
			Quantity:   menuItem.Quantity,
		}
	}

	return &repository.CreateBookingRequest{
		BookingID:       req.BookingId,
		CustomerName:    req.CustomerName,
		CompanyName:     req.CompanyName,
		BookingDateTime: bookingDateTime,
		PhoneNumber:     req.PhoneNumber,
		NumChildren:     req.NumChildren,
		NumAdults:       req.NumAdults,
		NumTables:       req.NumTables,
		Tables:          tables,
		MenuSets:        menuSets,
		MenuItems:       menuItems,
		Status:          req.Status,
		TotalPrice:      req.TotalPrice,
	}
}

func (s *bookingServer) CreateBooking(ctx context.Context, req *CreateBookingRequest) (*CreateBookingResponse, error) {
	// Load Bangkok timezone
	bangkok, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("could not load Bangkok timezone: %v", err))
	}

	bookingDateTime, err := time.Parse(time.RFC3339, req.BookingDateTime)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid booking_date_time format: %v", err))
	}

	// Convert to Bangkok timezone
	bookingDateTimeInBangkok := bookingDateTime.In(bangkok)

	if req.BookingId == "" {
		req.BookingId = uuid.New().String()
	}

	repositoryReq := ConvertCreateBookingRequestToRepositoryRequest(req, bookingDateTimeInBangkok)

	err = s.bookingRepo.CreateBooking(ctx, repositoryReq)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("could not create booking: %v", err))
	}

	return &CreateBookingResponse{BookingId: req.BookingId}, nil
}

func (s *bookingServer) UpdateBooking(ctx context.Context, req *CreateBookingRequest) (*UpdateBookingResponse, error) {
	if req.BookingId == "" {
		return nil, status.Error(codes.InvalidArgument, "booking_id is required for update")
	}

	// Load Bangkok timezone
	bangkok, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("could not load Bangkok timezone: %v", err))
	}

	bookingDateTime, err := time.Parse(time.RFC3339, req.BookingDateTime)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid booking_date_time format: %v", err))
	}

	bookingDateTimeInBangkok := bookingDateTime.In(bangkok)

	repositoryReq := ConvertCreateBookingRequestToRepositoryRequest(req, bookingDateTimeInBangkok)

	err = s.bookingRepo.UpdateBooking(ctx, req.BookingId, repositoryReq)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("could not update booking: %v", err))
	}

	return &UpdateBookingResponse{Success: true}, nil
}

func (s *bookingServer) DeleteBooking(ctx context.Context, req *DeleteBookingRequest) (*DeleteBookingResponse, error) {

	err := s.bookingRepo.DeleteBooking(ctx, req.BookingId)
	if err != nil {
		return nil, fmt.Errorf("could not delete booking: %v", err)
	}

	return &DeleteBookingResponse{Success: true}, nil
}
