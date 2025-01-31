package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gitlab.com/final_project1240930/api_gateway/internal/logs"
	services "gitlab.com/final_project1240930/api_gateway/internal/services/table"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type tableHandler struct {
	tableSrv services.TableService
}

func NewTableHandler(tableSrv services.TableService) *tableHandler {
	return &tableHandler{tableSrv: tableSrv}
}

func createErrorResponse(err error) map[string]string {
	return map[string]string{"error": err.Error()}
}

func marshalProtoMessage(resp proto.Message) (map[string]interface{}, error) {
	marshaler := protojson.MarshalOptions{
		EmitUnpopulated: true,
		UseEnumNumbers:  false, // ใช้ชื่อ enum แทนตัวเลข
	}

	jsonBytes, err := marshaler.Marshal(resp)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (h *tableHandler) CreateTable(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		logs.Error("Failed to read request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("failed to read request body")))
	}

	var req services.CreateTableRequest
	unmarshaler := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}

	if err := unmarshaler.Unmarshal(body, &req); err != nil {
		logs.Error("Invalid request format for CreateTable", zap.Error(err))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid request format")))
	}

	if req.Type == services.TableType_TABLE_TYPE_UNKNOWN {
		logs.Error("Invalid table type", zap.Any("Type", req.Type))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid table type")))
	}
	resp, err := h.tableSrv.CreateTable(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to create table", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	result, err := marshalProtoMessage(resp)
	if err != nil {
		logs.Error("Failed to marshal CreateTable response", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, result)
}

func (h *tableHandler) UpdateTable(c echo.Context) error {

	id := c.Param("id")
	if id == "" {
		logs.Error("Table ID is required")
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("table ID is required")))
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		logs.Error("Failed to read request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("failed to read request body")))
	}

	var req services.UpdateTableRequest
	unmarshaler := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}

	if err := unmarshaler.Unmarshal(body, &req); err != nil {
		logs.Error("Invalid request format for UpdateTable", zap.Error(err))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid request format")))
	}

	req.Id = id

	if req.Type == services.TableType_TABLE_TYPE_UNKNOWN {
		logs.Error("Invalid table type", zap.Any("Type", req.Type))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid table type")))
	}

	if req.NumTable <= 0 {
		logs.Error("Invalid table number", zap.Int32("NumTable", req.NumTable))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid table number")))
	}

	resp, err := h.tableSrv.UpdateTable(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to update table", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	result, err := marshalProtoMessage(resp)
	if err != nil {
		logs.Error("Failed to marshal UpdateTable response", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, result)
}

func (h *tableHandler) DeleteTable(c echo.Context) error {
	id := c.Param("id")
	req := services.DeleteTableRequest{Id: id}

	resp, err := h.tableSrv.DeleteTable(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to delete table", zap.String("tableId", id), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *tableHandler) GetTables(c echo.Context) error {
	resp, err := h.tableSrv.GetTables(c.Request().Context(), &emptypb.Empty{})
	if err != nil {
		logs.Error("Failed to get tables", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	result, err := marshalProtoMessage(resp)
	if err != nil {
		logs.Error("Failed to marshal tables response", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, result)
}

func (h *tableHandler) GetTableByNumTable(c echo.Context) error {
	number := c.Param("number")

	numTable, err := strconv.Atoi(number)
	if err != nil {
		logs.Error("Invalid table number", zap.String("tableNumber", number), zap.Error(err))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid table number")))
	}

	req := services.GetTableByNumTableRequest{NumTable: int32(numTable)}

	resp, err := h.tableSrv.GetTableByNumTable(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to get table by number", zap.String("tableNumber", number), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	result, err := marshalProtoMessage(resp)
	if err != nil {
		logs.Error("Failed to marshal table response", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, result)
}

func (h *tableHandler) GetAvailableTables(c echo.Context) error {
	date := c.Param("date")

	req := services.GetAvailableTablesRequest{Date: date}

	resp, err := h.tableSrv.GetAvailableTables(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to get table by number", zap.String("Date", date), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	result, err := marshalProtoMessage(resp)
	if err != nil {
		logs.Error("Failed to marshal table response", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, result)
}

func (h *tableHandler) UpdateTableType(c echo.Context) error {

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		logs.Error("Failed to read request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("failed to read request body")))
	}

	var req services.UpdateTableTypeRequest
	unmarshaler := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}

	if err := unmarshaler.Unmarshal(body, &req); err != nil {
		logs.Error("Invalid request format for UpdateTableType", zap.Error(err))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid request format")))
	}

	if req.Type == services.TableType_TABLE_TYPE_UNKNOWN {
		logs.Error("Invalid table type", zap.Any("Type", req.Type))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid table type")))
	}

	resp, err := h.tableSrv.UpdateTableType(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to update table type", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	result, err := marshalProtoMessage(resp)
	if err != nil {
		logs.Error("Failed to marshal UpdateTableType response", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, result)
}

func (h *tableHandler) ListTableTypes(c echo.Context) error {
	resp, err := h.tableSrv.ListTableTypes(c.Request().Context(), &emptypb.Empty{})
	if err != nil {
		logs.Error("Failed to list table types", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	result, err := marshalProtoMessage(resp)
	if err != nil {
		logs.Error("Failed to marshal table types response", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, result)
}
