package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/final_project1240930/booking_service/internal/logs"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Gorm implementation of TableRepository
type tableRepository struct {
	db *gorm.DB
}

func NewTableRepository(db *gorm.DB) TableRepository {
	return &tableRepository{db: db}
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

// ----------------  CRUD Methods for Table ------------------------

// CreateTable
func (r *tableRepository) CreateTable(ctx context.Context, table Table) (uuid.UUID, error) {
	// ตรวจสอบว่าเลขโต๊ะซ้ำหรือไม่
	var existingTable Table
	err := r.db.First(&existingTable, "num_table = ?", table.NumTable).Error
	if err == nil {
		// หากพบข้อมูลซ้ำ
		logs.Error("Table with NumTable already exists", zap.Int32("NumTable", table.NumTable))
		return uuid.Nil, fmt.Errorf("table with num_table %d already exists", table.NumTable)
	}

	// หากไม่พบข้อมูลซ้ำ ก็ทำการสร้างโต๊ะใหม่
	if err := r.db.Create(&table).Error; err != nil {
		logs.Error("Failed to create table", zap.Error(err))
		return uuid.Nil, fmt.Errorf("failed to create table: %w", err)
	}
	logs.Info("Successfully created table", zap.String("TableID", table.UUID.String()))
	return table.UUID, nil
}

// UpdateTable updates an existing table in the database
func (r *tableRepository) UpdateTable(ctx context.Context, table Table) error {
	if err := r.db.Save(&table).Error; err != nil {
		logs.Error("Failed to update table", zap.Error(err))
		return fmt.Errorf("failed to update table: %w", err)
	}
	logs.Info("Successfully updated table", zap.String("TableID", table.UUID.String()))
	return nil
}

// DeleteTable deletes a table by its ID
func (r *tableRepository) DeleteTable(ctx context.Context, tableID uuid.UUID) error {
	if err := r.db.Delete(&Table{}, "uuid = ?", tableID).Error; err != nil {
		logs.Error("Failed to delete table", zap.String("TableID", tableID.String()), zap.Error(err))
		return fmt.Errorf("failed to delete table: %w", err)
	}
	logs.Info("Successfully deleted table", zap.String("TableID", tableID.String()))
	return nil
}

// GetTables retrieves all tables from the database
func (r *tableRepository) GetTables(ctx context.Context) ([]Table, error) {
	var tables []Table
	if err := r.db.Find(&tables).Error; err != nil {
		logs.Error("Failed to fetch tables", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch tables: %w", err)
	}
	logs.Info("Successfully fetched tables", zap.Int("TotalTables", len(tables)))
	return tables, nil
}

// GetTableByID retrieves a single table by its ID
func (r *tableRepository) GetTableByID(ctx context.Context, tableID uuid.UUID) (Table, error) {
	var table Table
	if err := r.db.First(&table, "uuid = ?", tableID).Error; err != nil {
		logs.Error("Failed to fetch table by ID", zap.String("TableID", tableID.String()), zap.Error(err))
		return Table{}, fmt.Errorf("failed to fetch table by ID: %w", err)
	}
	logs.Info("Successfully fetched table by ID", zap.String("TableID", table.UUID.String()))
	return table, nil
}

func (r *tableRepository) GetTableByNumTable(ctx context.Context, numTable int32) (Table, error) {
	var table Table
	err := r.db.First(&table, "num_table = ?", numTable).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return Table{}, err
		}

		logs.Error("Failed to fetch table by num_table", zap.Int32("NumTable", numTable), zap.Error(err))
		return Table{}, fmt.Errorf("failed to fetch table by num_table: %w", err)
	}
	logs.Info("Successfully fetched table by num_table", zap.Int32("NumTable", numTable))
	return table, nil
}

func (r *tableRepository) GetAvailableTables(ctx context.Context, date string) ([]TableAvailability, error) {
	startTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v", err)
	}
	startTime = startTime.Add(10 * time.Hour) // 10:00 AM
	endTime := startTime.Add(12 * time.Hour)  // 10:00 PM

	var result []TableAvailability

	query := `
		WITH all_hours AS (
			SELECT generate_series(
				?::timestamp, 
				?::timestamp, 
				interval '1 hour'
			) AS hour_start
		),
		booked_tables AS (
			SELECT
				bt.table_id AS table_id,
				a.hour_start
			FROM
				all_hours a
			JOIN bookings b
			ON a.hour_start <= b.booking_date_time
				AND a.hour_start + interval '1 hour' > b.booking_date_time
			JOIN booking_tables bt ON bt.booking_id = b.uuid
			WHERE b.status NOT IN ('CANCELLED', 'COMPLETED')
		),
		available_tables AS (
			SELECT
				a.hour_start,
				t.uuid AS table_id,
				t.num_table,
				t."type"
			FROM
				all_hours a
			CROSS JOIN tables t
			WHERE NOT EXISTS (
				SELECT 1
				FROM booked_tables b
				WHERE b.hour_start = a.hour_start
				AND b.table_id = t.uuid
			)
		)
		SELECT
			to_char(a.hour_start, 'HH24:MI') AS time_slot,
			t.uuid AS table_id,
			t.num_table,
			t."type"
		FROM
			all_hours a
		CROSS JOIN tables t
		WHERE NOT EXISTS (
			SELECT 1
			FROM booked_tables b
			WHERE b.hour_start = a.hour_start
			AND b.table_id = t.uuid
		)
		ORDER BY
			a.hour_start, t.num_table;
	`

	// Execute query using GORM
	err = r.db.Raw(query, startTime, endTime).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ----------------  CRUD Methods for TableType ------------------------
func (r *tableRepository) UpdateTableType(ctx context.Context, tableType Type, seatCount int32) error {
	// ตรวจสอบว่า seatCount เป็นค่าที่ถูกต้องหรือไม่
	if seatCount <= 0 {
		logs.Error("Invalid seat count", zap.Int32("SeatCount", seatCount))
		return fmt.Errorf("seat count must be greater than 0")
	}

	// ใช้ GORM อัพเดทข้อมูลในตาราง table_types
	var tableTypeEnum string
	switch tableType {
	case STANDARD:
		tableTypeEnum = "STANDARD"
	case LARGE:
		tableTypeEnum = "LARGE"
	default:
		return fmt.Errorf("invalid table type: %v", tableType)
	}

	// อัพเดทข้อมูลในตาราง table_types
	err := r.db.Model(&TableType{}).
		Where("type = ?", tableTypeEnum).
		Update("seat_count", seatCount).Error

	if err != nil {
		logs.Error("Failed to update table type in database", zap.Error(err))
		return fmt.Errorf("failed to update table type in database: %v", err)
	}

	logs.Info("Table type updated successfully", zap.String("TableType", tableTypeEnum), zap.Int32("SeatCount", seatCount))
	return nil
}

// ListTableTypes all available table types (Standard, Large)
func (r *tableRepository) ListTableTypes(ctx context.Context) ([]TableType, error) {
	var tableTypes []TableType

	// ดึงข้อมูลประเภทโต๊ะทั้งหมดจากฐานข้อมูล
	err := r.db.Model(&TableType{}).Find(&tableTypes).Error
	if err != nil {
		logs.Error("Failed to fetch table types", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch table types: %v", err)
	}

	logs.Info("Successfully fetched table types", zap.Int("TotalTableTypes", len(tableTypes)))
	return tableTypes, nil
}
