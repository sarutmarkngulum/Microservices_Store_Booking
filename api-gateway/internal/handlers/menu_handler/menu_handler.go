package handlers

import (
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gitlab.com/final_project1240930/api_gateway/internal/logs"
	services "gitlab.com/final_project1240930/api_gateway/internal/services/menu"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type menuHandler struct {
	menuSrv services.MenuService
}

func NewMenuHandler(menuSrv services.MenuService) *menuHandler {
	return &menuHandler{menuSrv: menuSrv}
}

func createErrorResponse(err error) map[string]string {
	return map[string]string{"error": err.Error()}
}

func (h *menuHandler) CreateMenuItem(c echo.Context) error {
	var req services.CreateMenuItemRequest

	req.NameTh = c.FormValue("name_th")
	req.NameEn = c.FormValue("name_en")
	price := c.FormValue("price")
	req.Description = c.FormValue("description")
	categoryStr := c.FormValue("category")

	if price == "" {
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("price is required")))
	}
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid price format")))
	}
	req.Price = priceFloat

	if categoryStr == "" {
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("category is required")))
	}
	categoryInt, err := strconv.Atoi(categoryStr)
	if err != nil || categoryInt < 0 || categoryInt > int(services.MenuCategory_DESSERT) { // ตรวจสอบค่า category
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid category format")))
	}
	req.Category = services.MenuCategory(categoryInt)

	file, err := c.FormFile("image_data")
	if err == nil {
		fileContent, err := readFileAsBytes(file)
		if err != nil {
			return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("failed to read image file")))
		}
		req.ImageData = fileContent
	}

	resp, err := h.menuSrv.CreateMenuItem(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to create menu item", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func readFileAsBytes(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	fileBytes, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}

	// ตรวจสอบขนาดไฟล์ (ไฟล์ไม่เกิน 5MB)
	if len(fileBytes) > 5*1024*1024 {
		return nil, errors.New("file is too large")
	}
	return fileBytes, nil
}

func (h *menuHandler) UpdateMenuItem(c echo.Context) error {

	menuID := c.Param("id")
	if menuID == "" {
		logs.Error("Menu ID is missing from URL", zap.String("menu_id", menuID))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("menu ID is required")))
	}

	var req services.UpdateMenuItemRequest
	req.Id = menuID

	req.NameTh = c.FormValue("name_th")
	req.NameEn = c.FormValue("name_en")
	price := c.FormValue("price")
	req.Description = c.FormValue("description")
	categoryStr := c.FormValue("category")

	if price == "" {
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("price is required")))
	}
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid price format")))
	}
	req.Price = priceFloat

	if categoryStr == "" {
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("category is required")))
	}
	categoryInt, err := strconv.Atoi(categoryStr)
	if err != nil || categoryInt < 0 || categoryInt > int(services.MenuCategory_DESSERT) { // ตรวจสอบค่า category
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid category format")))
	}
	req.Category = services.MenuCategory(categoryInt)

	file, err := c.FormFile("image_data")
	if err == nil {
		fileContent, err := readFileAsBytes(file)
		if err != nil {
			return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("failed to read image file")))
		}
		req.ImageData = fileContent
	}

	resp, err := h.menuSrv.UpdateMenuItem(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to update menu item", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *menuHandler) DeleteMenuItem(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		logs.Error("ID is missing in the URL")
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("ID must be provided")))
	}

	req := services.DeleteMenuItemRequest{Id: id}

	resp, err := h.menuSrv.DeleteMenuItem(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to delete menu item", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *menuHandler) GetMenuItems(c echo.Context) error {
	resp, err := h.menuSrv.GetMenuItems(c.Request().Context(), &emptypb.Empty{})
	if err != nil {
		logs.Error("Failed to get menu items", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *menuHandler) GetMenuItemById(c echo.Context) error {
	id := c.Param("id")
	req := services.GetMenuItemByIdRequest{Id: id}

	resp, err := h.menuSrv.GetMenuItemById(c.Request().Context(), &req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			logs.Error("Menu item not found", zap.String("menuItemId", id))
			return c.JSON(http.StatusNotFound, createErrorResponse(errors.New("menu item not found")))
		}
		logs.Error("Failed to get menu item", zap.String("menuItemId", id), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *menuHandler) CreateMenuSet(c echo.Context) error {
	var req services.CreateMenuSetRequest
	if err := c.Bind(&req); err != nil {
		logs.Error("Invalid request format for CreateMenuSet", zap.Error(err))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid request format")))
	}

	resp, err := h.menuSrv.CreateMenuSet(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to create menu set", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *menuHandler) UpdateMenuSet(c echo.Context) error {
	menuSetID := c.Param("id")
	if menuSetID == "" {
		logs.Error("Menu Set ID is missing from URL", zap.String("menu_set_id", menuSetID))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("menu ID is required")))
	}

	var req services.UpdateMenuSetRequest
	req.Id = menuSetID

	if err := c.Bind(&req); err != nil {
		logs.Error("Invalid request format for UpdateMenuSet", zap.Error(err))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid request format")))
	}

	resp, err := h.menuSrv.UpdateMenuSet(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to update menu set", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *menuHandler) DeleteMenuSet(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		logs.Error("ID is missing in the URL")
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("ID must be provided")))
	}

	req := services.DeleteMenuSetRequest{
		Id: id,
	}

	resp, err := h.menuSrv.DeleteMenuSet(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to delete menu set", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)

}

func (h *menuHandler) GetMenuSets(c echo.Context) error {
	resp, err := h.menuSrv.GetMenuSets(c.Request().Context(), &emptypb.Empty{})
	if err != nil {
		logs.Error("Failed to get menu sets", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *menuHandler) GetMenuSetById(c echo.Context) error {
	id := c.Param("id")
	req := services.GetMenuSetByIdRequest{Id: id}

	resp, err := h.menuSrv.GetMenuSetById(c.Request().Context(), &req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			logs.Error("Menu set not found", zap.String("menuSetId", id))
			return c.JSON(http.StatusNotFound, createErrorResponse(errors.New("menu set not found")))
		}
		logs.Error("Failed to get menu set", zap.String("menuSetId", id), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *menuHandler) CreateMenuSetItem(c echo.Context) error {
	var req services.CreateMenuSetItemRequest
	if err := c.Bind(&req); err != nil {
		logs.Error("Invalid request format for CreateMenuSetItem", zap.Error(err))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid request format")))
	}

	resp, err := h.menuSrv.CreateMenuSetItem(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to create menu set item", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *menuHandler) GetMenuSetItems(c echo.Context) error {
	resp, err := h.menuSrv.GetMenuSetItems(c.Request().Context(), &emptypb.Empty{})
	if err != nil {
		logs.Error("Failed to get menu set items", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *menuHandler) GetMenuSetItemById(c echo.Context) error {
	id := c.Param("id")
	req := services.GetMenuSetItemByIdRequest{MenuSetId: id}

	resp, err := h.menuSrv.GetMenuSetItemByMenuSetID(c.Request().Context(), &req)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			logs.Error("Menu set item not found", zap.String("menuSetItemId", id))
			return c.JSON(http.StatusNotFound, createErrorResponse(errors.New("menu set item not found")))
		}
		logs.Error("Failed to get menu set item", zap.String("menuSetItemId", id), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *menuHandler) UpdateMenuSetItem(c echo.Context) error {
	var req services.UpdateMenuSetItemRequest
	if err := c.Bind(&req); err != nil {
		logs.Error("Invalid request format for UpdateMenuSetItem", zap.Error(err))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid request format")))
	}

	resp, err := h.menuSrv.UpdateMenuSetItem(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to update menu set item", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *menuHandler) DeleteMenuSetItem(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		logs.Error("ID is missing in the URL")
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("ID must be provided")))
	}

	req := services.DeleteMenuSetItemRequest{
		MenuSetId: id,
	}

	resp, err := h.menuSrv.DeleteMenuSetItem(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to delete menu set item", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}
