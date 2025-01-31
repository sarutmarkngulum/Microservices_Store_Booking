package services

import (
	"context"
	"time"

	"gitlab.com/final_project1240930/api_gateway/internal/logs"
	"go.uber.org/zap"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type TableService interface {
	// Handle Table
	CreateTable(ctx context.Context, req *CreateTableRequest) (*CreateTableResponse, error)
	UpdateTable(ctx context.Context, req *UpdateTableRequest) (*UpdateTableResponse, error)
	DeleteTable(ctx context.Context, req *DeleteTableRequest) (*DeleteTableResponse, error)
	GetTables(ctx context.Context, req *emptypb.Empty) (*TableList, error)
	GetTableByNumTable(ctx context.Context, req *GetTableByNumTableRequest) (*GetTableByNumTableResponse, error)
	GetAvailableTables(ctx context.Context, req *GetAvailableTablesRequest) (*GetAvailableTablesResponse, error)

	// Handle Table Type Count
	UpdateTableType(ctx context.Context, req *UpdateTableTypeRequest) (*UpdateTableTypeResponse, error)
	ListTableTypes(ctx context.Context, req *emptypb.Empty) (*TableTypeList, error)
}

type tableService struct {
	tableClient TableServiceClient
}

func NewTableService(tableClient TableServiceClient) TableService {
	return &tableService{tableClient: tableClient}
}

// wrapper function with timeout, logging, and error handling
func (s *tableService) callWithTimeout(ctx context.Context, call func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5) // Set timeout
	defer cancel()

	res, err := call(ctx)
	if err != nil {
		logs.Error("Error calling service", zap.Error(err))
		return nil, err
	}
	return res, nil
}

// Implement the methods with the wrapper

func (s *tableService) CreateTable(ctx context.Context, req *CreateTableRequest) (*CreateTableResponse, error) {
	res, err := s.callWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.tableClient.CreateTable(ctx, req)
	})
	if res != nil {
		return res.(*CreateTableResponse), nil
	}
	return nil, err
}

func (s *tableService) UpdateTable(ctx context.Context, req *UpdateTableRequest) (*UpdateTableResponse, error) {
	res, err := s.callWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.tableClient.UpdateTable(ctx, req)
	})
	if res != nil {
		return res.(*UpdateTableResponse), nil
	}
	return nil, err
}

func (s *tableService) DeleteTable(ctx context.Context, req *DeleteTableRequest) (*DeleteTableResponse, error) {
	res, err := s.callWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.tableClient.DeleteTable(ctx, req)
	})
	if res != nil {
		return res.(*DeleteTableResponse), nil
	}
	return nil, err
}

func (s *tableService) GetTables(ctx context.Context, req *emptypb.Empty) (*TableList, error) {
	res, err := s.callWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.tableClient.GetTables(ctx, req)
	})
	if res != nil {
		return res.(*TableList), nil
	}
	return nil, err
}

func (s *tableService) GetTableByNumTable(ctx context.Context, req *GetTableByNumTableRequest) (*GetTableByNumTableResponse, error) {
	res, err := s.callWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.tableClient.GetTableByNumTable(ctx, req)
	})
	if res != nil {
		return res.(*GetTableByNumTableResponse), nil
	}
	return nil, err
}

func (s *tableService) GetAvailableTables(ctx context.Context, req *GetAvailableTablesRequest) (*GetAvailableTablesResponse, error) {
	res, err := s.callWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.tableClient.GetAvailableTables(ctx, req)
	})
	if res != nil {
		return res.(*GetAvailableTablesResponse), nil
	}
	return nil, err
}

func (s *tableService) UpdateTableType(ctx context.Context, req *UpdateTableTypeRequest) (*UpdateTableTypeResponse, error) {
	res, err := s.callWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.tableClient.UpdateTableType(ctx, req)
	})
	if res != nil {
		return res.(*UpdateTableTypeResponse), nil
	}
	return nil, err
}

func (s *tableService) ListTableTypes(ctx context.Context, req *emptypb.Empty) (*TableTypeList, error) {
	res, err := s.callWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.tableClient.ListTableTypes(ctx, req)
	})
	if res != nil {
		return res.(*TableTypeList), nil
	}
	return nil, err
}
