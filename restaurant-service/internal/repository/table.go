package repository

import (
	"context"

	"github.com/gofrs/uuid"
)

type Type string

const (
	STANDARD Type = "STANDARD"
	LARGE    Type = "LARGE"
)

type Table struct {
	UUID     uuid.UUID `gorm:"column:uuid;type:uuid;default:gen_random_uuid();primaryKey"`
	NumTable int32     `gorm:"type:int;not null;unique"`
	Type     Type      `gorm:"type:table_type;not null"`
}

type TableType struct {
	Type      Type `gorm:"column:type;primaryKey"`
	SeatCount int  `gorm:"column:seat_count"`
}

type TableAvailability struct {
	TimeSlot string `json:"time_slot"`
	TableID  string `json:"table_id"`
	NumTable int    `json:"num_table"`
	Type     Type   `json:"type"`
}

type TableRepository interface {
	// CRUD for Tables
	CreateTable(ctx context.Context, table Table) (uuid.UUID, error)
	UpdateTable(ctx context.Context, table Table) error
	DeleteTable(ctx context.Context, tableID uuid.UUID) error
	GetTables(ctx context.Context) ([]Table, error)
	GetTableByID(ctx context.Context, tableID uuid.UUID) (Table, error)
	GetTableByNumTable(ctx context.Context, numTable int32) (Table, error)

	GetAvailableTables(ctx context.Context, date string) ([]TableAvailability, error)

	// CRUD for Table Types
	UpdateTableType(ctx context.Context, Type Type, seatCount int32) error
	ListTableTypes(ctx context.Context) ([]TableType, error)
}
