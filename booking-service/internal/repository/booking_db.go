package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type bookingRepository struct {
	DB *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{DB: db}
}

// NewDatabase creates a new database connection
func NewDatabase(host string, port int, user, password, dbname string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Successfully connected to the database.")
	return db, nil
}

// CloseDatabase closes the database connection
func CloseDatabase(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get DB instance: %w", err)
	}
	return sqlDB.Close()
}

func (r *bookingRepository) mapBookingDetails(entities []BookingEntity) []Booking {
	bookingMap := make(map[string]Booking) // สร้างแผนที่เพื่อเก็บ Booking

	for _, entity := range entities {
		// ตรวจสอบว่า BookingID นี้มีในแผนที่แล้วหรือยัง ถ้ายังให้สร้าง Booking ใหม่
		if _, exists := bookingMap[entity.BookingID]; !exists {
			bookingMap[entity.BookingID] = Booking{
				BookingID:        entity.BookingID,
				CustomerName:     entity.CustomerName,
				CompanyName:      entity.CompanyName,
				BookingDateTime:  entity.BookingDateTime,
				PhoneNumber:      entity.PhoneNumber,
				NumChildren:      entity.NumChildren,
				NumAdults:        entity.NumAdults,
				NumTables:        entity.NumTables,
				TotalPrice:       entity.TotalPrice,
				Tables:           []BookingTable{},    // เริ่มต้น slice ของตาราง
				BookingMenuSets:  []BookingMenuSet{},  // เริ่มต้น slice ของ Menu Sets
				BookingMenuItems: []BookingMenuItem{}, // เริ่มต้น slice ของ Menu
				Status:           entity.Status,
			}
		}

		// ดึงข้อมูลจากแผนที่มาใช้งาน
		booking := bookingMap[entity.BookingID]

		// จัดการข้อมูลตาราง
		if entity.TableID != "" {
			// เช็คว่ามีการเพิ่มตารางที่ซ้ำกันหรือไม่
			tableExists := false
			for _, table := range booking.Tables {
				if table.TableID == entity.TableID {
					tableExists = true
					break
				}
			}
			if !tableExists {
				// ถ้าไม่มีตารางนี้ในรายการให้เพิ่มเข้าไป
				booking.Tables = append(booking.Tables, BookingTable{
					TableID:     entity.TableID,
					TableNumber: entity.TableNumber,
					Type:        entity.TableType,
					SeatCount:   entity.SeatCount,
				})
			}
		}

		// จัดการข้อมูล Menu Items (จาก Separate Menu Item)
		if entity.SeparateMenuItemID != "" {
			// เช็คว่ามีการเพิ่ม Menu Item นี้ไปแล้วหรือไม่
			menuItemExists := false
			for _, item := range booking.BookingMenuItems {
				if item.MenuItemID == entity.SeparateMenuItemID {
					// ถ้า Menu Item ซ้ำกัน ไม่ต้องบวก Quantity, ใช้ค่าเดิม
					menuItemExists = true
					break
				}
			}
			if !menuItemExists {
				// ถ้ายังไม่มี Menu Item นี้ในรายการให้เพิ่มเข้าไป
				booking.BookingMenuItems = append(booking.BookingMenuItems, BookingMenuItem{
					MenuItemID:  entity.SeparateMenuItemID,
					NameTh:      entity.SeparateMenuItemNameTh,
					NameEn:      entity.SeparateMenuItemNameEn,
					Description: entity.SeparateMenuItemDescription,
					Price:       float32(entity.SeparateMenuItemPrice),
					Category:    entity.SeparateMenuItemCategory,
					ImageURL:    entity.SeparateMenuItemImageURL,
					Quantity:    entity.SeparateMenuItemQuantity,
				})
			}
		}

		// จัดการข้อมูล Menu Sets
		if entity.MenuSetID != "" {
			// ตรวจสอบว่ามีการเพิ่ม Menu Set นี้แล้วหรือยัง
			menuSetExists := false
			var updatedMenuSet *BookingMenuSet
			for i, menuSet := range booking.BookingMenuSets {
				if menuSet.MenuSetID == entity.MenuSetID {
					// ถ้ามี Menu Set นี้แล้ว จะเพิ่ม Menu Items เข้าไป
					menuSetExists = true
					updatedMenuSet = &booking.BookingMenuSets[i]
					break
				}
			}

			// ถ้า Menu Set นี้ยังไม่เคยถูกเพิ่มเข้าไป ให้เพิ่มเข้าไปใหม่
			if !menuSetExists {
				// เพิ่ม Menu Set ใหม่
				booking.BookingMenuSets = append(booking.BookingMenuSets, BookingMenuSet{
					MenuSetID:    entity.MenuSetID,
					MenuSetName:  entity.MenuSetName,
					MenuSetPrice: float32(entity.MenuSetPrice),
					Quantity:     entity.MenuSetQuantity,
					MenuItems: []BookingMenuItem{
						{
							MenuItemID:  entity.MenuItemID,
							NameTh:      entity.MenuItemNameTh,
							NameEn:      entity.MenuItemNameEn,
							Description: entity.MenuItemDescription,
							Price:       float32(entity.MenuItemPrice),
							Category:    entity.MenuItemCategory,
							ImageURL:    entity.MenuItemImageURL,
						},
					},
				})
			} else {
				// ถ้า Menu Set นี้มีอยู่แล้ว ให้เพิ่ม Menu Items ลงไป
				// ตรวจสอบว่ามี Menu Item ซ้ำกันในรายการหรือไม่
				menuItemExists := false
				for _, item := range updatedMenuSet.MenuItems {
					if item.MenuItemID == entity.MenuItemID {
						menuItemExists = true
						break
					}
				}
				if !menuItemExists {
					updatedMenuSet.MenuItems = append(updatedMenuSet.MenuItems, BookingMenuItem{
						MenuItemID:  entity.MenuItemID,
						NameTh:      entity.MenuItemNameTh,
						NameEn:      entity.MenuItemNameEn,
						Description: entity.MenuItemDescription,
						Price:       float32(entity.MenuItemPrice),
						Category:    entity.MenuItemCategory,
						ImageURL:    entity.MenuItemImageURL,
					})
				}
			}
		}

		// อัพเดทแผนที่ด้วยข้อมูลที่ปรับปรุง
		bookingMap[entity.BookingID] = booking
	}

	// แปลงแผนที่ BookingMap เป็น slice ของ Booking
	var results []Booking
	for _, booking := range bookingMap {
		results = append(results, booking)
	}

	return results
}

// GetBookingDetails ดึงข้อมูลการจองจากฐานข้อมูลและแปลงเป็น BookingDetails
func (r *bookingRepository) GetBookingDetails(ctx context.Context) ([]Booking, error) {
	var entities []BookingEntity

	query := `
		SELECT 
			b.uuid AS booking_id,
			b.customer_name,
			b.company_name,
			b.booking_date_time,
			b.phone_number,
			b.num_children,
			b.num_adults,
			b.num_tables,
			b.total_price,
			bt.table_id,
			t.num_table,
			t.type,
			tt.seat_count,
			ms.menu_set_id,
			ms.quantity as menu_set_quantity,
			set_menu.name AS menu_set_name,
			set_menu.price AS menu_set_price,
			mi.uuid AS menu_item_id,
			mi.name_th AS menu_item_name_th,
			mi.name_en AS menu_item_name_en,
			mi.description AS menu_item_description,
			mi.price AS menu_item_price,
			mi.category AS menu_item_category,
			mi.image_url AS menu_item_image_url,
			bmi.quantity as separate_menu_item_quantity,
			mi2.uuid AS separate_menu_item_id,
			mi2.name_th AS separate_menu_item_name_th,
			mi2.name_en AS separate_menu_item_name_en,
			mi2.description AS separate_menu_item_description,
			mi2.price AS separate_menu_item_price,
			mi2.category AS separate_menu_item_category,
			mi2.image_url AS separate_menu_item_image_url,
			b.status 
		FROM bookings b
		LEFT JOIN booking_tables bt ON bt.booking_id = b.uuid
		LEFT JOIN tables t ON bt.table_id = t.uuid
		LEFT JOIN table_types tt ON t.type = tt.type
		LEFT JOIN booking_menu_sets ms ON ms.booking_id = b.uuid
		LEFT JOIN menu_sets set_menu ON ms.menu_set_id = set_menu.uuid
		LEFT JOIN menu_set_items msi ON msi.menu_set_id = ms.menu_set_id
		LEFT JOIN menu_items mi ON msi.menu_item_id = mi.uuid
		LEFT JOIN booking_menu_items bmi ON bmi.booking_id = b.uuid
		LEFT JOIN menu_items mi2 ON bmi.menu_item_id = mi2.uuid 
		ORDER BY b.booking_date_time;
	`

	err := r.DB.Raw(query).Scan(&entities).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query bookings: %w", err)
	}

	return r.mapBookingDetails(entities), nil
}

// GetBookingDetailsByID
func (r *bookingRepository) GetBookingDetailsByID(ctx context.Context, bookingID string) (*Booking, error) {
	var entities []BookingEntity

	query := `
		SELECT 
			b.uuid AS booking_id,
			b.customer_name,
			b.company_name,
			b.booking_date_time,
			b.phone_number,
			b.num_children,
			b.num_adults,
			b.num_tables,
			b.total_price,
			bt.table_id,
			t.num_table,
			t.type,
			tt.seat_count,
			ms.menu_set_id,
			ms.quantity as menu_set_quantity,
			set_menu.name AS menu_set_name,
			set_menu.price AS menu_set_price,
			mi.uuid AS menu_item_id,
			mi.name_th AS menu_item_name_th,
			mi.name_en AS menu_item_name_en,
			mi.description AS menu_item_description,
			mi.price AS menu_item_price,
			mi.category AS menu_item_category,
			mi.image_url AS menu_item_image_url,
			bmi.quantity as separate_menu_item_quantity,
			mi2.uuid AS separate_menu_item_id,
			mi2.name_th AS separate_menu_item_name_th,
			mi2.name_en AS separate_menu_item_name_en,
			mi2.description AS separate_menu_item_description,
			mi2.price AS separate_menu_item_price,
			mi2.category AS separate_menu_item_category,
			mi2.image_url AS separate_menu_item_image_url,
			b.status 
		FROM bookings b
		LEFT JOIN booking_tables bt ON bt.booking_id = b.uuid
		LEFT JOIN tables t ON bt.table_id = t.uuid
		LEFT JOIN table_types tt ON t.type = tt.type
		LEFT JOIN booking_menu_sets ms ON ms.booking_id = b.uuid
		LEFT JOIN menu_sets set_menu ON ms.menu_set_id = set_menu.uuid
		LEFT JOIN menu_set_items msi ON msi.menu_set_id = ms.menu_set_id
		LEFT JOIN menu_items mi ON msi.menu_item_id = mi.uuid
		LEFT JOIN booking_menu_items bmi ON bmi.booking_id = b.uuid
		LEFT JOIN menu_items mi2 ON bmi.menu_item_id = mi2.uuid 
		WHERE b.uuid = ? 
		ORDER BY b.booking_date_time;
	`

	// ใช้การ query ข้อมูลจากฐานข้อมูล
	err := r.DB.Raw(query, bookingID).Scan(&entities).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query booking by ID: %w", err)
	}

	bookings := r.mapBookingDetails(entities)

	// ถ้าหากไม่มีการจองที่ตรงกับ bookingID ที่ระบุ ให้คืนค่า nil และ error
	if len(bookings) == 0 {
		return nil, fmt.Errorf("booking not found for ID: %s", bookingID)
	}

	return &bookings[0], nil
}

func (r *bookingRepository) CreateBooking(ctx context.Context, req *CreateBookingRequest) error {
	// Start a transaction
	tx := r.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// ตรวจสอบการจองในเวลาเดียวกัน (ไม่ตรวจสอบเบอร์โทรซ้ำในเวลาต่างกัน)
	var existingBooking CreateBooking
	err := tx.Where("booking_date_time = ?", req.BookingDateTime).First(&existingBooking).Error
	if err == nil {
		// ถ้ามีข้อมูลที่ซ้ำกัน (ไม่พบ error) ให้คืนค่าผลลัพธ์ว่า Booking นี้มีอยู่แล้ว
		tx.Rollback() // Rollback the transaction
		return fmt.Errorf("Booking at the same time already exists")
	} else if err != gorm.ErrRecordNotFound {
		// ถ้ามีข้อผิดพลาดอื่นๆ ในการค้นหา
		tx.Rollback()
		return fmt.Errorf("Error checking existing booking: %s", err)
	}

	// สร้าง UUID สำหรับ Booking
	bookingID := uuid.New().String()

	// เพิ่มข้อมูล Booking หลัก
	booking := CreateBooking{
		BookingID:       bookingID,
		CustomerName:    req.CustomerName,
		CompanyName:     req.CompanyName,
		BookingDateTime: req.BookingDateTime,
		PhoneNumber:     req.PhoneNumber,
		NumChildren:     req.NumChildren,
		NumAdults:       req.NumAdults,
		NumTables:       req.NumTables,
		Status:          "CONFIRMED",
		TotalPrice:      req.TotalPrice,
	}

	if err := tx.Create(&booking).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, tableID := range req.Tables {
		tableEntity := BookingTableEntity{
			BookingID: bookingID,
			TableID:   tableID.TableID,
		}
		if err := tx.Create(&tableEntity).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, menuSet := range req.MenuSets {
		menuSetEntity := BookingMenuSetEntity{
			BookingID: bookingID,
			MenuSetID: menuSet.MenuSetID,
			Quantity:  menuSet.Quantity,
		}
		if err := tx.Create(&menuSetEntity).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("error creating menu set: %v", err)
		}
	}

	for _, menuItem := range req.MenuItems {
		menuItemEntity := BookingMenuItemEntity{
			BookingID:  bookingID,
			MenuItemID: menuItem.MenuItemID,
			Quantity:   menuItem.Quantity,
		}
		if err := tx.Create(&menuItemEntity).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction if all operations succeed
	return tx.Commit().Error
}

func (r *bookingRepository) UpdateBooking(ctx context.Context, bookingID string, req *CreateBookingRequest) error {
	// เริ่มต้นการทำธุรกรรมเพื่อความสอดคล้อง
	tx := r.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// ดึงข้อมูลการจองที่มีอยู่
	var booking CreateBooking
	if err := tx.Where("uuid = ?", bookingID).First(&booking).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to query booking by ID: %w", err)
	}

	// อัปเดตข้อมูลการจองหลัก
	if err := tx.Model(&booking).Where("uuid = ?", bookingID).Updates(map[string]interface{}{
		"customer_name":     req.CustomerName,
		"company_name":      req.CompanyName,
		"booking_date_time": req.BookingDateTime,
		"phone_number":      req.PhoneNumber,
		"num_children":      req.NumChildren,
		"num_adults":        req.NumAdults,
		"num_tables":        req.NumTables,
		"total_price":       req.TotalPrice,
		"status":            req.Status,
	}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to : %w", err)
	}

	// ลบข้อมูลที่เกี่ยวข้องเก่าก่อนที่จะเพิ่มข้อมูลใหม่
	tx.Where("booking_id = ?", bookingID).Delete(&BookingTableEntity{})
	tx.Where("booking_id = ?", bookingID).Delete(&BookingMenuSetEntity{})
	tx.Where("booking_id = ?", bookingID).Delete(&BookingMenuItemEntity{})

	// เพิ่มข้อมูลที่เกี่ยวข้องใหม่
	for _, table := range req.Tables {
		if err := tx.Create(&BookingTableEntity{
			BookingID: bookingID,
			TableID:   table.TableID,
		}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, menuSet := range req.MenuSets {
		if err := tx.Create(&BookingMenuSetEntity{
			BookingID: bookingID,
			MenuSetID: menuSet.MenuSetID,
			Quantity:  menuSet.Quantity,
		}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, menuItem := range req.MenuItems {
		if err := tx.Create(&BookingMenuItemEntity{
			BookingID:  bookingID,
			MenuItemID: menuItem.MenuItemID,
			Quantity:   menuItem.Quantity,
		}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// คอมมิทการทำธุรกรรม
	return tx.Commit().Error
}

func (r *bookingRepository) DeleteBooking(ctx context.Context, bookingID string) error {
	// ตรวจสอบว่า bookingID เป็นค่าว่าง
	if bookingID == "" {
		return fmt.Errorf("booking ID cannot be empty")
	}

	// Start a transaction
	tx := r.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Delete related Tables, MenuSets, and MenuItems
	if err := tx.Where("booking_id = ?", bookingID).Delete(&BookingTableEntity{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("booking_id = ?", bookingID).Delete(&BookingMenuSetEntity{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("booking_id = ?", bookingID).Delete(&BookingMenuItemEntity{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete the Booking
	if err := tx.Where("uuid = ?", bookingID).Delete(&CreateBooking{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting booking: %v", err)
	}

	// Commit the transaction if all operations succeed
	return tx.Commit().Error
}
