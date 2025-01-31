package repository

import (
	"context"
	"time"
)

type Booking struct {
	BookingID        string            `gorm:"type:uuid;primary_key" json:"booking_id"`
	CustomerName     string            `json:"customer_name"`
	CompanyName      string            `json:"company_name"`
	BookingDateTime  time.Time         `json:"booking_date_time"`
	PhoneNumber      string            `json:"phone_number"`
	NumChildren      int32             `json:"num_children"`
	NumAdults        int32             `json:"num_adults"`
	NumTables        int32             `json:"num_tables"`
	TotalPrice       float64           `json:"total_price"`
	Tables           []BookingTable    `json:"tables"`
	BookingMenuSets  []BookingMenuSet  `json:"menu_sets"`
	BookingMenuItems []BookingMenuItem `json:"menu_items"`
	Status           string            `json:"status"`
}

type BookingTable struct {
	TableID     string `json:"table_id"`
	TableNumber string `json:"table_number"`
	Type        string `json:"type"`
	SeatCount   int32  `json:"seat_count"`
}

type BookingMenuSet struct {
	MenuSetID    string            `json:"menu_set_id"`
	MenuSetName  string            `json:"menu_set_name"`
	MenuSetPrice float32           `json:"menu_set_price"`
	Quantity     int32             `json:"quantity"`
	MenuItems    []BookingMenuItem `json:"menu_items"`
}

type BookingMenuItem struct {
	MenuItemID  string  `json:"menu_item_id"`
	NameTh      string  `json:"name_th"`
	NameEn      string  `json:"name_en"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
	Quantity    int32   `json:"quantity"`
}

type BookingEntity struct {
	BookingID                   string    `gorm:"column:booking_id"`
	CustomerName                string    `gorm:"column:customer_name"`
	CompanyName                 string    `gorm:"column:company_name"`
	BookingDateTime             time.Time `gorm:"column:booking_date_time"`
	PhoneNumber                 string    `gorm:"column:phone_number"`
	NumChildren                 int32     `gorm:"column:num_children"`
	NumAdults                   int32     `gorm:"column:num_adults"`
	NumTables                   int32     `gorm:"column:num_tables"`
	TableID                     string    `gorm:"column:table_id"`
	TableNumber                 string    `gorm:"column:num_table"`
	TotalPrice                  float64   `gorm:"column:total_price"`
	TableType                   string    `gorm:"column:type"`
	SeatCount                   int32     `gorm:"column:seat_count"`
	MenuSetID                   string    `gorm:"column:menu_set_id"`
	MenuSetName                 string    `gorm:"column:menu_set_name"`
	MenuSetQuantity             int32     `gorm:"column:menu_set_quantity"`
	MenuSetPrice                float64   `gorm:"column:menu_set_price"`
	MenuItemID                  string    `gorm:"column:menu_item_id"`
	MenuItemNameTh              string    `gorm:"column:menu_item_name_th"`
	MenuItemNameEn              string    `gorm:"column:menu_item_name_en"`
	MenuItemDescription         string    `gorm:"column:menu_item_description"`
	MenuItemPrice               float64   `gorm:"column:menu_item_price"`
	MenuItemCategory            string    `gorm:"column:menu_item_category"`
	MenuItemImageURL            string    `gorm:"column:menu_item_image_url"`
	MenuItemQuantity            int32     `gorm:"column:menu_item_quantity"`
	SeparateMenuItemID          string    `gorm:"column:separate_menu_item_id"`
	SeparateMenuItemNameTh      string    `gorm:"column:separate_menu_item_name_th"`
	SeparateMenuItemNameEn      string    `gorm:"column:separate_menu_item_name_en"`
	SeparateMenuItemDescription string    `gorm:"column:separate_menu_item_description"`
	SeparateMenuItemPrice       float64   `gorm:"column:separate_menu_item_price"`
	SeparateMenuItemCategory    string    `gorm:"column:separate_menu_item_category"`
	SeparateMenuItemImageURL    string    `gorm:"column:separate_menu_item_image_url"`
	SeparateMenuItemQuantity    int32     `gorm:"column:separate_menu_item_quantity"`
	Status                      string    `json:"status"`
}

// Request
type CreateBookingRequest struct {
	BookingID       string                  `json:"booking_id" binding:"required"`
	CustomerName    string                  `json:"customer_name" binding:"required"`     // ชื่อลูกค้า
	CompanyName     string                  `json:"company_name"`                         // ชื่อบริษัท (ถ้ามี)
	BookingDateTime time.Time               `json:"booking_date_time" binding:"required"` // วันและเวลาที่จอง
	PhoneNumber     string                  `json:"phone_number" binding:"required"`      // เบอร์โทรศัพท์
	NumChildren     int32                   `json:"num_children"`                         // จำนวนเด็ก
	NumAdults       int32                   `json:"num_adults" binding:"required"`        // จำนวนผู้ใหญ่
	NumTables       int32                   `json:"num_tables" binding:"required"`        // จำนวนโต๊ะที่จอง
	Tables          []CreateBookingTable    `json:"tables"`                               // รายการโต๊ะ
	MenuSets        []CreateBookingMenuSet  `json:"menu_sets"`                            // รายการเมนูเซ็ต
	MenuItems       []CreateBookingMenuItem `json:"menu_items"`                           // รายการเมนูอาหาร
	Status          string                  `gorm:"column:status" json:"status"`
	TotalPrice      float64                 `gorm:"column:total_price" json:"total_price"`
}

type CreateBookingTable struct {
	TableID string `json:"table_id" binding:"required"`
}

type CreateBookingMenuSet struct {
	MenuSetID string `json:"menu_set_id" binding:"required"`
	Quantity  int32  `json:"quantity" binding:"required"`
}

type CreateBookingMenuItem struct {
	MenuItemID string `json:"menu_item_id" binding:"required"`
	Quantity   int32  `json:"quantity" binding:"required"`
}

type CreateBooking struct {
	BookingID       string    `gorm:"column:uuid" json:"booking_id"`
	CustomerName    string    `gorm:"column:customer_name" json:"customer_name"`
	CompanyName     string    `gorm:"column:company_name" json:"company_name"`
	BookingDateTime time.Time `gorm:"column:booking_date_time" json:"booking_date_time"`
	PhoneNumber     string    `gorm:"column:phone_number" json:"phone_number"`
	NumChildren     int32     `gorm:"column:num_children" json:"num_children"`
	NumAdults       int32     `gorm:"column:num_adults" json:"num_adults"`
	NumTables       int32     `gorm:"column:num_tables" json:"num_tables"`
	Status          string    `gorm:"column:status" json:"status"`
	TotalPrice      float64   `gorm:"column:total_price" json:"total_price"`
}

func (CreateBooking) TableName() string {
	return "bookings"
}

type BookingTableEntity struct {
	BookingID string `gorm:"column:booking_id;primaryKey" json:"booking_id"`
	TableID   string `gorm:"column:table_id;primaryKey" json:"table_id"`
}

func (BookingTableEntity) TableName() string {
	return "booking_tables"
}

type BookingMenuSetEntity struct {
	BookingID string `gorm:"column:booking_id;primaryKey"`
	MenuSetID string `gorm:"column:menu_set_id"`
	Quantity  int32  `gorm:"column:quantity"`
}

func (BookingMenuSetEntity) TableName() string {
	return "booking_menu_sets"
}

type BookingMenuItemEntity struct {
	BookingID  string `gorm:"column:booking_id;primaryKey"`
	MenuItemID string `gorm:"column:menu_item_id"`
	Quantity   int32  `gorm:"column:quantity"`
}

func (BookingMenuItemEntity) TableName() string {
	return "booking_menu_items"
}

type BookingRepository interface {
	GetBookingDetails(ctx context.Context) ([]Booking, error)
	GetBookingDetailsByID(ctx context.Context, bookingID string) (*Booking, error)

	CreateBooking(ctx context.Context, booking *CreateBookingRequest) error
	UpdateBooking(ctx context.Context, bookingID string, req *CreateBookingRequest) error
	DeleteBooking(ctx context.Context, bookingID string) error
}
