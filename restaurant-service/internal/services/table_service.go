package services

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/final_project1240930/booking_service/internal/logs"
	"gitlab.com/final_project1240930/booking_service/internal/repository"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type tableServer struct {
	tableRepo repository.TableRepository
}

func NewTableServer(tableRepo repository.TableRepository) TableServiceServer {
	return &tableServer{tableRepo: tableRepo}
}

func (s *tableServer) mustEmbedUnimplementedTableServiceServer() {}

// ---------------- Table ------------------------

// CreateTable
func (s *tableServer) CreateTable(ctx context.Context, req *CreateTableRequest) (*CreateTableResponse, error) {
	logs.Info("Received CreateTableRequest", zap.Int32("NumTable", req.GetNumTable()))

	if req.GetNumTable() <= 0 {
		logs.Error("Validation error: NumTable is invalid", zap.Int32("NumTable", req.GetNumTable()))
		return nil, status.Errorf(codes.InvalidArgument, "NumTable must be greater than 0")
	}

	// แปลงค่า TableType จาก Protobuf ไปเป็น Repository
	dbTableType, err := mapProtoTableTypeToRepository(req.GetType())
	if err != nil {
		logs.Error("Invalid table type", zap.String("Type", req.GetType().String()))
		return nil, status.Errorf(codes.InvalidArgument, "Invalid table type: %v", req.GetType().String())
	}

	// ตรวจสอบว่าเลขโต๊ะซ้ำหรือไม่
	if err := s.checkTableExistsByNum(ctx, req.GetNumTable()); err != nil {
		return nil, err
	}

	table := repository.Table{
		NumTable: req.GetNumTable(),
		Type:     dbTableType,
	}

	tableID, err := s.tableRepo.CreateTable(ctx, table)
	if err != nil {
		logs.Error("Failed to create table", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to create table: %v", err)
	}

	logs.Info("Successfully created table", zap.String("TableID", tableID.String()))

	return &CreateTableResponse{
		Id:     tableID.String(),
		Status: "success",
	}, nil
}

// Helper function to check if table exists by NumTable
func (s *tableServer) checkTableExistsByNum(ctx context.Context, numTable int32) error {
	_, err := s.tableRepo.GetTableByNumTable(ctx, numTable)
	if err == nil {
		logs.Error("Validation error: NumTable already exists", zap.Int32("NumTable", numTable))
		return status.Errorf(codes.AlreadyExists, "Table with NumTable %d already exists", numTable)
	}
	if err != gorm.ErrRecordNotFound {
		logs.Error("Error fetching table by NumTable", zap.Int32("NumTable", numTable), zap.Error(err))
		return status.Errorf(codes.Internal, "Error checking table: %v", err)
	}
	return nil
}

func (s *tableServer) UpdateTable(ctx context.Context, req *UpdateTableRequest) (*UpdateTableResponse, error) {
	logs.Info("Received UpdateTableRequest", zap.String("TableID", req.GetId()))

	tableID, err := uuid.FromString(req.GetId())
	if err != nil {
		logs.Error("Validation error: Invalid table ID", zap.String("TableID", req.GetId()))
		return nil, status.Errorf(codes.InvalidArgument, "Invalid table ID")
	}

	table, err := s.tableRepo.GetTableByID(ctx, tableID)
	if err != nil {
		logs.Error("Failed to fetch table", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to fetch table: %v", err)
	}

	// ตรวจสอบว่า NumTable ที่อัพเดตไม่ซ้ำกับที่มีอยู่แล้ว
	if req.GetNumTable() != table.NumTable {
		if err := s.checkTableExistsByNum(ctx, req.GetNumTable()); err != nil {
			logs.Error("Validation error: Table number already exists", zap.Int32("NumTable", req.GetNumTable()))
			return nil, status.Errorf(codes.InvalidArgument, "Table number already exists: %v", req.GetNumTable())
		}
	}

	tableType, err := mapProtoTableTypeToRepository(req.GetType())
	if err != nil {
		return nil, fmt.Errorf("invalid type: %v", req.GetType())
	}

	table.NumTable = req.GetNumTable()
	table.Type = tableType

	if err := s.tableRepo.UpdateTable(ctx, table); err != nil {
		logs.Error("Failed to update table", zap.Error(err))
		return nil, err
	}

	logs.Info("Successfully updated table", zap.String("TableID", table.UUID.String()))

	return &UpdateTableResponse{
		Status: "success",
	}, nil
}

func (s *tableServer) DeleteTable(ctx context.Context, req *DeleteTableRequest) (*DeleteTableResponse, error) {
	logs.Info("Received DeleteTableRequest", zap.String("TableID", req.GetId()))

	tableID, err := uuid.FromString(req.GetId())
	if err != nil {
		logs.Error("Validation error: Invalid table ID", zap.String("TableID", req.GetId()))
		return nil, status.Errorf(codes.InvalidArgument, "Invalid table ID")
	}

	if err := s.tableRepo.DeleteTable(ctx, tableID); err != nil {
		logs.Error("Failed to delete table", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to delete table: %v", err)
	}

	logs.Info("Successfully deleted table", zap.String("TableID", tableID.String()))

	return &DeleteTableResponse{
		Status: "success",
	}, nil
}

func (s *tableServer) GetTables(ctx context.Context, _ *emptypb.Empty) (*TableList, error) {
	logs.Info("Received GetTablesRequest")

	tables, err := s.tableRepo.GetTables(ctx)
	if err != nil {
		logs.Error("Failed to fetch tables", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to fetch tables: %v", err)
	}

	var tableList []*Table
	for _, table := range tables {
		pbType, err := mapRepositoryTableTypeToProto(table.Type)
		if err != nil {
			logs.Error("Unknown table type", zap.String("Type", string(table.Type)))
			continue
		}

		tableList = append(tableList, &Table{
			Id:       table.UUID.String(),
			NumTable: table.NumTable,
			Type:     pbType,
		})
	}

	logs.Info("Successfully fetched tables", zap.Int("TotalTables", len(tableList)))

	return &TableList{Tables: tableList}, nil
}

func (s *tableServer) GetTableByNumTable(ctx context.Context, req *GetTableByNumTableRequest) (*GetTableByNumTableResponse, error) {

	table, err := s.tableRepo.GetTableByNumTable(ctx, req.NumTable)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get table by num_table: %v", err)
	}

	pbType, err := mapRepositoryTableTypeToProto(table.Type)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map table type: %v", err)
	}

	return &GetTableByNumTableResponse{
		Table: &Table{
			Id:       table.UUID.String(),
			NumTable: table.NumTable,
			Type:     pbType,
		},
	}, nil
}

func (s *tableServer) GetAvailableTables(ctx context.Context, req *GetAvailableTablesRequest) (*GetAvailableTablesResponse, error) {

	tableAvailability, err := s.tableRepo.GetAvailableTables(ctx, req.Date)
	if err != nil {
		return nil, fmt.Errorf("could not get available tables: %v", err)
	}

	timeSlotMap := make(map[string][]*Table)

	for _, availability := range tableAvailability {

		pbType, err := mapRepositoryTableTypeToProto(availability.Type)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to map table type: %v", err)
		}
		table := &Table{
			Id:       availability.TableID,
			NumTable: int32(availability.NumTable),
			Type:     TableType(pbType),
		}

		timeSlotMap[availability.TimeSlot] = append(timeSlotMap[availability.TimeSlot], table)
	}

	var availabilityList []*TableAvailability
	for timeSlot, tables := range timeSlotMap {
		availabilityList = append(availabilityList, &TableAvailability{
			TimeSlot: timeSlot,
			Tables:   tables,
		})
	}

	sort.Slice(availabilityList, func(i, j int) bool {
		time1, _ := time.Parse("15:04", availabilityList[i].TimeSlot)
		time2, _ := time.Parse("15:04", availabilityList[j].TimeSlot)
		return time1.Before(time2)
	})

	return &GetAvailableTablesResponse{
		AvailableTables: availabilityList,
	}, nil
}

// ---------------- Table Type ------------------------

func (s *tableServer) UpdateTableType(ctx context.Context, req *UpdateTableTypeRequest) (*UpdateTableTypeResponse, error) {
	logs.Info("Received UpdateTableTypeRequest", zap.String("TableType", req.GetType().String()))

	var tableType repository.Type
	switch req.GetType() {
	case TableType_STANDARD:
		tableType = repository.STANDARD
	case TableType_LARGE:
		tableType = repository.LARGE
	default:
		logs.Error("Invalid table type received")
		return nil, status.Errorf(codes.InvalidArgument, "Invalid table type: %v", req.GetType().String())
	}

	err := s.tableRepo.UpdateTableType(ctx, tableType, req.GetSeatCount())
	if err != nil {
		logs.Error("Failed to update table type", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to update table type: %v", err)
	}

	logs.Info("Successfully updated table type")

	return &UpdateTableTypeResponse{
		Status: "success",
	}, nil
}

// ListTableTypes
func (s *tableServer) ListTableTypes(ctx context.Context, _ *emptypb.Empty) (*TableTypeList, error) {
	logs.Info("Received ListTableTypesRequest")

	tableTypeCounts, err := s.tableRepo.ListTableTypes(ctx)
	if err != nil {
		logs.Error("Failed to list table types", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to list table types: %v", err)
	}

	var pbTableTypeCounts []*TableTypeCount
	for _, tableTypeCount := range tableTypeCounts {
		logs.Info("Processing table type", zap.String("Type", string(tableTypeCount.Type)), zap.Int("SeatCount", tableTypeCount.SeatCount))

		protoType, err := mapRepositoryTableTypeToProto(tableTypeCount.Type) // ใช้ฟังก์ชันที่ถูกต้อง
		if err != nil {
			logs.Error("Skipping unknown table type", zap.String("Type", string(tableTypeCount.Type)))
			continue
		}

		pbTableTypeCounts = append(pbTableTypeCounts, &TableTypeCount{
			Type:  protoType,
			Count: int32(tableTypeCount.SeatCount),
		})
	}

	return &TableTypeList{
		TableTypes: pbTableTypeCounts,
	}, nil
}

func mapRepositoryTableTypeToProto(tableType repository.Type) (TableType, error) {
	switch tableType {
	case repository.STANDARD:
		return TableType_STANDARD, nil
	case repository.LARGE:
		return TableType_LARGE, nil
	default:
		return TableType_TABLE_TYPE_UNKNOWN, fmt.Errorf("unknown table type: %v", tableType)
	}
}

func mapProtoTableTypeToRepository(tableType TableType) (repository.Type, error) {
	switch tableType {
	case TableType_STANDARD:
		return repository.STANDARD, nil
	case TableType_LARGE:
		return repository.LARGE, nil
	default:
		return "", fmt.Errorf("unknown table type: %v", tableType)
	}
}
