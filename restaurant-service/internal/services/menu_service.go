package services

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"gitlab.com/final_project1240930/booking_service/internal/logs"
	"gitlab.com/final_project1240930/booking_service/internal/repository"
	"go.uber.org/zap"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type menuServer struct {
	menuRepo repository.MenuRepository
}

func NewMenuServer(menuRepo repository.MenuRepository) MenuServiceServer {
	return &menuServer{menuRepo: menuRepo}
}

func (s menuServer) mustEmbedUnimplementedMenuServiceServer() {
}

func mapProtoCategoryToRepoCategory(protoCategory MenuCategory) (repository.MenuCategory, error) {
	switch protoCategory {
	case MenuCategory_MAIN_COURSE:
		return repository.MainCourse, nil
	case MenuCategory_BEVERAGE:
		return repository.Beverage, nil
	case MenuCategory_DESSERT:
		return repository.Dessert, nil
	default:
		return "", fmt.Errorf("invalid category: %v", protoCategory)
	}
}

// ---------------- Image Upload & Delete ------------------------

func (s *menuServer) UploadImage(ctx context.Context, req *UploadImageRequest) (*UploadImageResponse, error) {
	logs.Info("Received UploadImageRequest", zap.String("FileName", req.FileName))

	// ตรวจสอบว่ามี image_data หรือไม่
	if len(req.ImageData) == 0 {
		logs.Error("Validation error: No image data provided")
		return nil, status.Errorf(codes.InvalidArgument, "Image data must be provided")
	}

	// อัปโหลดไฟล์ไปยัง Cloudinary
	imageUrl, err := s.UploadImageToCloudinary(req.ImageData, req.FileName)
	if err != nil {
		logs.Error("Failed to upload image", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to upload image: %v", err)
	}

	// ส่งคืน URL ของภาพที่อัปโหลด
	return &UploadImageResponse{
		ImageUrl: imageUrl,
		Status:   Status_SUCCESS,
	}, nil
}

// ฟังก์ชันช่วยอัปโหลดรูปภาพไปยัง Cloudinary
func (s *menuServer) UploadImageToCloudinary(imageData []byte, fileName string) (string, error) {
	if len(imageData) == 0 {
		logs.Error("No image data received")
		return "", fmt.Errorf("no image data received")
	}

	resp, err := s.menuRepo.UploadImage(context.Background(), imageData, fileName)
	if err != nil {
		logs.Error("Failed to upload image to Cloudinary", zap.Error(err))
		return "", fmt.Errorf("failed to upload image to cloudinary: %v", err)
	}

	logs.Info("Image uploaded to Cloudinary", zap.String("image_url", resp))
	return resp, nil
}

// ฟังก์ชันการลบรูปภาพ
func (s *menuServer) DeleteImage(ctx context.Context, req *DeleteImageRequest) (*DeleteImageResponse, error) {
	// ตรวจสอบว่ามี URL หรือไม่
	if req.ImageUrl == "" {
		return nil, fmt.Errorf("image URL must be provided for deletion")
	}

	// สร้าง repository.Image จาก req.ImageUrl
	image := repository.Image{
		ImageURL: req.ImageUrl, // ใช้ req.ImageUrl ในการตั้งค่า
	}

	// ลบรูปภาพจาก Cloud Storage หรือที่เก็บข้อมูล
	err := s.menuRepo.DeleteImage(ctx, image)
	if err != nil {
		logs.Error("Failed to delete image", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to delete image: %v", err)
	}

	// ส่งคืนสถานะการลบสำเร็จ
	return &DeleteImageResponse{
		Status: Status_SUCCESS,
	}, nil
}

// ---------------- Menu ------------------------

func (s *menuServer) CreateMenuItem(ctx context.Context, req *CreateMenuItemRequest) (*CreateMenuItemResponse, error) {

	if req.NameTh == "" || req.NameEn == "" {
		logs.Error("Validation error: NameTh or NameEn is empty")
		return nil, status.Errorf(codes.InvalidArgument, "NameTh and NameEn must be provided")
	}
	if req.Price <= 0 {
		logs.Error("Validation error: Invalid Price", zap.Float64("Price", req.Price))
		return nil, status.Errorf(codes.InvalidArgument, "Price must be greater than 0")
	}
	if req.Description == "" {
		logs.Error("Validation error: Description is empty")
		return nil, status.Errorf(codes.InvalidArgument, "Description must be provided")
	}
	println("Images", req.ImageData)

	// แปลงประเภท Category
	category, err := mapProtoCategoryToRepoCategory(req.Category)
	if err != nil {
		logs.Error("Validation error: Invalid Category", zap.String("Category", req.Category.String()))
		return nil, status.Errorf(codes.InvalidArgument, "Invalid category")
	}

	// อัปโหลดรูปภาพถ้ามี
	imageUrl := "no image"      // ถ้าไม่มีรูปภาพจะตั้งค่าเป็น "no image"
	if len(req.ImageData) > 0 { // ถ้ามี ImageData
		logs.Info("Image data received", zap.Int("ImageDataLength", len(req.ImageData)))

		// สร้าง UUID สำหรับชื่อไฟล์
		newUUID, err := uuid.NewV4()
		if err != nil {
			logs.Error("Error generating UUID", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "Error generating UUID: %v", err)
		}
		fileName := fmt.Sprintf("menu_image_%s.jpg", newUUID.String())

		// อัปโหลดไฟล์ไปยัง Cloudinary
		imageUrl, err = s.UploadImageToCloudinary(req.ImageData, fileName)
		if err != nil {
			logs.Error("Failed to upload image", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "Failed to upload image: %v", err)
		}
	}

	// สร้าง MenuItem
	item := repository.MenuItem{
		NameTH:      req.NameTh,
		NameEN:      req.NameEn,
		Description: req.Description,
		Price:       req.Price,
		Category:    category,
		ImageURL:    imageUrl, // ถ้ามีการอัปโหลดจะได้ URL กลับมา, ถ้าไม่มีจะเป็น "no image"
	}

	newID, err := s.menuRepo.CreateMenuItem(ctx, item)
	if err != nil {
		logs.Error("Failed to create menu item", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to create menu item: %v", err)
	}

	logs.Info("Successfully created menu item", zap.String("MenuID", newID.String()))

	return &CreateMenuItemResponse{
		Id:     newID.String(),
		Status: Status_SUCCESS,
	}, nil
}

func (s *menuServer) UpdateMenuItem(ctx context.Context, req *UpdateMenuItemRequest) (*UpdateMenuItemResponse, error) {

	if req.Id == "" {
		logs.Error("Validation error: ID is empty")
		return nil, status.Errorf(codes.InvalidArgument, "ID must be provided")
	}
	if req.NameTh == "" || req.NameEn == "" {
		logs.Error("Validation error: NameTh or NameEn is empty")
		return nil, status.Errorf(codes.InvalidArgument, "NameTh and NameEn must be provided")
	}
	if req.Price <= 0 {
		logs.Error("Validation error: Invalid Price", zap.Float64("Price", req.Price))
		return nil, status.Errorf(codes.InvalidArgument, "Price must be greater than 0")
	}
	if req.Description == "" {
		logs.Error("Validation error: Description is empty")
		return nil, status.Errorf(codes.InvalidArgument, "Description must be provided")
	}

	// แปลงประเภท Category
	category, err := mapProtoCategoryToRepoCategory(req.Category)
	if err != nil {
		logs.Error("Validation error: Invalid Category", zap.String("Category", req.Category.String()))
		return nil, status.Errorf(codes.InvalidArgument, "Invalid category")
	}

	id, err := uuid.FromString(req.Id)
	if err != nil {
		logs.Error("Invalid UUID format", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "Invalid menu item ID format")
	}

	menuItem, err := s.menuRepo.GetMenuItemByID(ctx, id)
	if err != nil {
		logs.Error("Menu item not found", zap.String("ID", req.Id))
		return nil, status.Errorf(codes.NotFound, "Menu item not found")
	}

	// อัปเดตค่าต่าง ๆ ของเมนู
	if req.NameTh != "" {
		menuItem.NameTH = req.NameTh
	}
	if req.NameEn != "" {
		menuItem.NameEN = req.NameEn
	}
	if req.Description != "" {
		menuItem.Description = req.Description
	}
	if req.Price > 0 {
		menuItem.Price = req.Price
	}
	menuItem.Category = category

	/// อัปโหลดรูปภาพใหม่ถ้ามี
	if len(req.ImageData) > 0 {
		// ลบรูปภาพเก่าก่อนหากมี
		if menuItem.ImageURL != "no image" {
			err := s.menuRepo.DeleteImage(ctx, repository.Image{ImageURL: menuItem.ImageURL})
			if err != nil {
				logs.Error("Failed to delete old image from Cloudinary", zap.Error(err))
				return nil, status.Errorf(codes.Internal, "Failed to delete old image: %v", err)
			}
		}

		// สร้าง UUID สำหรับชื่อไฟล์
		newUUID, err := uuid.NewV4()
		if err != nil {
			logs.Error("Error generating UUID", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "Error generating UUID: %v", err)
		}
		fileName := fmt.Sprintf("menu_image_%s.jpg", newUUID.String())

		// อัปโหลดไฟล์ไปยัง Cloudinary
		imageUrl, err := s.UploadImageToCloudinary(req.ImageData, fileName)
		if err != nil {
			logs.Error("Failed to upload image", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "Failed to upload image: %v", err)
		}
		menuItem.ImageURL = imageUrl
	}

	if err := s.menuRepo.UpdateMenuItem(ctx, menuItem); err != nil {
		logs.Error("Failed to update menu item", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to update menu item: %v", err)
	}

	logs.Info("Successfully updated menu item", zap.String("ID", menuItem.UUID.String()))

	return &UpdateMenuItemResponse{
		Status: Status_SUCCESS,
	}, nil
}

func (s *menuServer) DeleteMenuItem(ctx context.Context, req *DeleteMenuItemRequest) (*DeleteMenuItemResponse, error) {
	logs.Info("Received DeleteMenuItemRequest", zap.String("ID", req.Id))

	id, err := uuid.FromString(req.Id)
	if err != nil {
		logs.Error("Invalid UUID format", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "Invalid menu item ID format")
	}

	menuItem, err := s.menuRepo.GetMenuItemByID(ctx, id)
	if err != nil {
		logs.Error("Failed to find menu item", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "Menu item not found")
	}

	// ลบรูปภาพจาก Cloudinary หากมี
	if menuItem.ImageURL != "no image" {
		err := s.menuRepo.DeleteImage(ctx, repository.Image{ImageURL: menuItem.ImageURL})
		if err != nil {
			logs.Error("Failed to delete image from Cloudinary", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "Failed to delete image: %v", err)
		}
	}

	// ลบข้อมูลเมนูจากฐานข้อมูล
	if err := s.menuRepo.DeleteMenuItem(ctx, menuItem.UUID); err != nil {
		logs.Error("Failed to delete menu item", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to delete menu item: %v", err)
	}

	logs.Info("Menu item deleted successfully", zap.String("ID", req.Id))

	return &DeleteMenuItemResponse{
		Status: Status_SUCCESS,
	}, nil
}

func (s *menuServer) GetMenuItems(ctx context.Context, _ *emptypb.Empty) (*MenuItemList, error) {

	menuItems, err := s.menuRepo.GetMenuItems(ctx)
	if err != nil {
		logs.Error("Failed to get menu items", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to get menu items: %v", err)
	}

	var menuItemList []*MenuItem
	for _, item := range menuItems {
		var protoCategory MenuCategory
		switch item.Category {
		case repository.MainCourse:
			protoCategory = MenuCategory_MAIN_COURSE
		case repository.Beverage:
			protoCategory = MenuCategory_BEVERAGE
		case repository.Dessert:
			protoCategory = MenuCategory_DESSERT
		default:
			return nil, fmt.Errorf("invalid category: %v", item.Category)
		}

		menuItemList = append(menuItemList, &MenuItem{
			Id:          item.UUID.String(),
			NameTh:      item.NameTH,
			NameEn:      item.NameEN,
			Description: item.Description,
			Price:       item.Price,
			Category:    protoCategory,
			ImageUrl:    item.ImageURL,
		})
	}

	return &MenuItemList{
		MenuItems: menuItemList,
	}, nil
}

func (s *menuServer) GetMenuItemById(ctx context.Context, req *GetMenuItemByIdRequest) (*MenuItem, error) {
	logs.Info("Received GetMenuItemByIdRequest", zap.String("ItemID", req.Id))

	id, err := uuid.FromString(req.Id)
	if err != nil {
		logs.Error("Invalid UUID format", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "Invalid menu item ID format")
	}

	menuItem, err := s.menuRepo.GetMenuItemByID(ctx, id)
	if err != nil {
		logs.Error("Failed to find menu item", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "Menu item not found")
	}

	var protoCategory MenuCategory
	switch menuItem.Category {
	case repository.MainCourse:
		protoCategory = MenuCategory_MAIN_COURSE
	case repository.Beverage:
		protoCategory = MenuCategory_BEVERAGE
	case repository.Dessert:
		protoCategory = MenuCategory_DESSERT
	default:
		logs.Error("Invalid category in menu item", zap.String("Category", string(menuItem.Category)))
		return nil, status.Errorf(codes.Internal, "Invalid category in menu item")
	}

	response := &MenuItem{
		Id:          menuItem.UUID.String(),
		NameTh:      menuItem.NameTH,
		NameEn:      menuItem.NameEN,
		Description: menuItem.Description,
		Price:       menuItem.Price,
		Category:    protoCategory,
		ImageUrl:    menuItem.ImageURL,
	}

	logs.Info("Successfully retrieved menu item", zap.String("ItemID", menuItem.UUID.String()))

	return response, nil
}

// ---------------- Menu Set ------------------------

func (s *menuServer) CreateMenuSet(ctx context.Context, req *CreateMenuSetRequest) (*CreateMenuSetResponse, error) {
	logs.Info("Received CreateMenuSetRequest", zap.String("Name", req.Name))

	if req.Name == "" {
		logs.Error("Validation error: Name is empty")
		return nil, status.Errorf(codes.InvalidArgument, "Name must be provided")
	}

	if req.Price <= 0 {
		logs.Error("Validation error: Invalid Price", zap.Float64("Price", req.Price))
		return nil, status.Errorf(codes.InvalidArgument, "Price must be greater than 0")
	}

	menuSet := repository.MenuSet{
		Name:  req.Name,
		Price: req.Price,
	}

	newID, err := s.menuRepo.CreateMenuSet(ctx, menuSet)
	if err != nil {
		logs.Error("Failed to create menu set", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to create menu set: %v", err)
	}

	logs.Info("Successfully created menu set", zap.String("MenuSetID", newID.String()))

	return &CreateMenuSetResponse{
		Id:     newID.String(),
		Status: Status_SUCCESS,
	}, nil
}

func (s *menuServer) UpdateMenuSet(ctx context.Context, req *UpdateMenuSetRequest) (*UpdateMenuSetResponse, error) {
	logs.Info("Received UpdateMenuSetRequest", zap.String("SetID", req.Id))

	id, err := uuid.FromString(req.Id)
	if err != nil {
		logs.Error("Invalid UUID format", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "Invalid menu set ID format")
	}

	menuSet, err := s.menuRepo.GetMenuSetByID(ctx, id)
	if err != nil {
		logs.Error("Failed to find menu set", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "Menu set not found")
	}

	if req.Name != "" {
		menuSet.Name = req.Name
		menuSet.Price = req.Price
	}

	if err := s.menuRepo.UpdateMenuSet(ctx, menuSet); err != nil {
		logs.Error("Failed to update menu set", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to update menu set: %v", err)
	}

	logs.Info("Menu set updated successfully", zap.String("ID", menuSet.UUID.String()))

	return &UpdateMenuSetResponse{
		Status: Status_SUCCESS,
	}, nil
}

func (s *menuServer) DeleteMenuSet(ctx context.Context, req *DeleteMenuSetRequest) (*DeleteMenuSetResponse, error) {
	logs.Info("Received DeleteMenuSetRequest", zap.String("SetID", req.Id))

	id, err := uuid.FromString(req.Id)
	if err != nil {
		logs.Error("Invalid UUID format", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "Invalid menu set ID format")
	}

	if err := s.menuRepo.DeleteMenuSet(ctx, id); err != nil {
		logs.Error("Failed to delete menu set", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to delete menu set: %v", err)
	}

	logs.Info("Menu set deleted successfully", zap.String("ID", req.Id))

	return &DeleteMenuSetResponse{
		Status: Status_SUCCESS,
	}, nil
}

func (s *menuServer) GetMenuSets(ctx context.Context, _ *emptypb.Empty) (*MenuSetList, error) {

	menuSets, err := s.menuRepo.GetMenuSets(ctx)
	if err != nil {
		logs.Error("Failed to get menu sets", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to get menu sets: %v", err)
	}

	var menuSetList []*MenuSet
	for _, set := range menuSets {
		menuSetList = append(menuSetList, &MenuSet{
			Id:   set.UUID.String(),
			Name: set.Name,
		})
	}

	return &MenuSetList{
		MenuSets: menuSetList,
	}, nil
}

func (s *menuServer) GetMenuSetById(ctx context.Context, req *GetMenuSetByIdRequest) (*MenuSet, error) {
	logs.Info("Received GetMenuSetByIdRequest", zap.String("SetID", req.Id))

	id, err := uuid.FromString(req.Id)
	if err != nil {
		logs.Error("Invalid UUID format", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "Invalid menu set ID format")
	}

	menuSet, err := s.menuRepo.GetMenuSetByID(ctx, id)
	if err != nil {
		logs.Error("Failed to find menu set", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "Menu set not found")
	}

	response := &MenuSet{
		Id:   menuSet.UUID.String(),
		Name: menuSet.Name,
	}

	logs.Info("Successfully retrieved menu set", zap.String("SetID", menuSet.UUID.String()))

	return response, nil
}

// ---------------- Menu Set Item ------------------------

func (s *menuServer) CreateMenuSetItem(ctx context.Context, req *CreateMenuSetItemRequest) (*CreateMenuSetItemResponse, error) {

	if req.MenuSetId == "" || len(req.MenuItemId) == 0 {
		logs.Error("Validation error: MenuSetId or MenuItemId is empty",
			zap.String("MenuSetId", req.MenuSetId),
			zap.Int("MenuItemIds", len(req.MenuItemId)),
		)
		return nil, status.Errorf(codes.InvalidArgument, "MenuSetId and at least one MenuItemId must be provided")
	}

	menuSetID, err := uuid.FromString(req.MenuSetId)
	if err != nil {
		logs.Error("Invalid MenuSetId format", zap.String("MenuSetId", req.MenuSetId), zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "Invalid MenuSetId format")
	}

	var menuItemIDs []uuid.UUID
	for _, itemID := range req.MenuItemId {
		menuItemID, err := uuid.FromString(itemID)
		if err != nil {
			logs.Error("Invalid MenuItemId format", zap.String("MenuItemId", itemID), zap.Error(err))
			return nil, status.Errorf(codes.InvalidArgument, "Invalid MenuItemId format")
		}
		menuItemIDs = append(menuItemIDs, menuItemID)
	}

	err = s.menuRepo.CreateMenuSetItems(ctx, menuSetID, menuItemIDs)
	if err != nil {
		logs.Error("Failed to create menu set item", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to create menu set item: %v", err)
	}

	return &CreateMenuSetItemResponse{
		Status: Status_SUCCESS,
	}, nil
}

func (s *menuServer) UpdateMenuSetItem(ctx context.Context, req *UpdateMenuSetItemRequest) (*UpdateMenuSetItemResponse, error) {
	if req.MenuSetId == "" || len(req.MenuItemId) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "MenuSetId and at least one MenuItemId must be provided")
	}

	menuSetID, err := uuid.FromString(req.MenuSetId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid MenuSetId format")
	}

	var menuItemIDs []uuid.UUID
	for _, itemID := range req.MenuItemId {
		menuItemID, err := uuid.FromString(itemID)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Invalid MenuItemId format")
		}
		menuItemIDs = append(menuItemIDs, menuItemID)
	}

	err = s.menuRepo.UpdateMenuSetItems(ctx, menuSetID, menuItemIDs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update menu set items: %v", err)
	}

	return &UpdateMenuSetItemResponse{
		Status: Status_SUCCESS,
	}, nil
}

func (s *menuServer) DeleteMenuSetItem(ctx context.Context, req *DeleteMenuSetItemRequest) (*DeleteMenuSetItemResponse, error) {
	menuSetID, err := uuid.FromString(req.MenuSetId)
	if err != nil {
		logs.Error("Invalid MenuSetId format", zap.String("MenuSetId", req.MenuSetId), zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "Invalid MenuSetId format")
	}

	err = s.menuRepo.DeleteMenuSetItem(ctx, menuSetID)
	if err != nil {
		logs.Error("Failed to delete menu set items", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to delete menu set items: %v", err)
	}

	return &DeleteMenuSetItemResponse{
		Status: Status_SUCCESS,
	}, nil
}

func (s *menuServer) GetMenuSetItems(ctx context.Context, _ *emptypb.Empty) (*MenuSetItemList, error) {

	items, err := s.menuRepo.GetMenuSetItems(ctx)
	if err != nil {
		logs.Error("Failed to get menu set items", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to get menu set items: %v", err)
	}

	var menuSetItems []*MenuSetItem
	for _, item := range items {
		menuSetItems = append(menuSetItems, &MenuSetItem{
			MenuSetId:    item.MenuSetId,
			MenuItemId:   item.MenuItemId,
			MenuSetName:  item.MenuSetName,
			MenuSetPrice: float32(item.MenuSetPrice),
			MenuNameTh:   item.MenuNameTh,
			MenuNameEn:   item.MenuNameEn,
			MenuPrice:    float32(item.MenuPrice),
			MenuCategory: item.MenuCategory,
			ImageUrl:     item.ImageUrl,
		})
	}

	return &MenuSetItemList{
		MenuSetItems: menuSetItems,
	}, nil
}

func (s *menuServer) GetMenuSetItemByMenuSetID(ctx context.Context, req *GetMenuSetItemByIdRequest) (*MenuSetItemList, error) {

	menuSetID, err := uuid.FromString(req.MenuSetId)
	if err != nil {
		logs.Error("Invalid MenuSetId format", zap.String("MenuSetId", req.MenuSetId), zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "Invalid MenuSetId format")
	}

	items, err := s.menuRepo.GetMenuSetItemByMenuSetID(ctx, menuSetID)
	if err != nil {
		logs.Error("Failed to get menu set items by MenuSetId", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to get menu set items by MenuSetId: %v", err)
	}

	var menuSetItems []*MenuSetItem
	for _, item := range items {
		menuSetItems = append(menuSetItems, &MenuSetItem{
			MenuSetId:    item.MenuSetId,
			MenuItemId:   item.MenuItemId,
			MenuSetName:  item.MenuSetName,
			MenuSetPrice: float32(item.MenuSetPrice),
			MenuNameTh:   item.MenuNameTh,
			MenuNameEn:   item.MenuNameEn,
			MenuPrice:    float32(item.MenuPrice),
			MenuCategory: item.MenuCategory,
			ImageUrl:     item.ImageUrl,
		})
	}

	return &MenuSetItemList{
		MenuSetItems: menuSetItems,
	}, nil
}
