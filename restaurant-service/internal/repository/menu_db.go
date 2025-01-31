package repository

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofrs/uuid"
	"gitlab.com/final_project1240930/booking_service/internal/logs"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type menuRepositoryDB struct {
	db         *gorm.DB
	cloudinary *cloudinary.Cloudinary
}

func NewMenuRepository(db *gorm.DB, cloudinaryURL string) MenuRepository {
	// ตรวจสอบว่า URL เป็นค่าว่างหรือไม่
	if cloudinaryURL == "" {
		logs.Fatal("Cloudinary URL is missing")
	}
	logs.Info("Using Cloudinary URL", zap.String("cloudinaryURL", cloudinaryURL))

	// สร้าง Cloudinary instance จาก URL
	cloudinaryInstance, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		logs.Fatal("Failed to initialize Cloudinary", zap.Error(err))
	}

	// ทดสอบการอัปโหลดไฟล์โดยใช้ไฟล์ทดสอบ
	_, err = cloudinaryInstance.Upload.Upload(context.Background(), bytes.NewReader([]byte("test")), uploader.UploadParams{})
	if err != nil {
		logs.Error("Cloudinary upload test failed", zap.Error(err))
	} else {
		logs.Info("Cloudinary upload test successful")
	}

	// สร้างและคืนค่า menuRepositoryDB
	return &menuRepositoryDB{
		db:         db,
		cloudinary: cloudinaryInstance,
	}
}

// ---------------- Images Upload & Delete ------------------------

// UploadImage method
func (r *menuRepositoryDB) UploadImage(ctx context.Context, imageData []byte, fileName string) (string, error) {
	var resp *uploader.UploadResult
	var err error

	// ล็อกข้อมูลก่อนอัปโหลด
	logs.Info("Uploading image to Cloudinary", zap.String("file_name", fileName))

	resp, err = r.cloudinary.Upload.Upload(ctx, bytes.NewReader(imageData), uploader.UploadParams{
		Folder:   "menu_images",
		PublicID: fileName,
	})

	if err != nil {
		logs.Error("Failed to upload image to Cloudinary", zap.Error(err))
		return "", fmt.Errorf("failed to upload image: %v", err)
	}

	// ล็อกการตอบกลับจาก Cloudinary
	logs.Info("Cloudinary upload response", zap.Any("response", resp))

	if resp == nil || resp.SecureURL == "" {
		logs.Error("Cloudinary upload response is empty", zap.String("file_name", fileName))
		return "", fmt.Errorf("received empty response from Cloudinary")
	}

	logs.Info("Image uploaded successfully", zap.String("image_url", resp.SecureURL))
	return resp.SecureURL, nil
}

// DeleteImage method
func (r *menuRepositoryDB) DeleteImage(ctx context.Context, image Image) error {
	publicID := extractPublicID(image.ImageURL)
	if publicID == "" {
		logs.Error("Invalid Cloudinary image URL", zap.String("image_url", image.ImageURL))
		return fmt.Errorf("invalid image URL")
	}

	// ลบภาพจาก Cloudinary
	_, err := r.cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		logs.Error("Failed to delete image from Cloudinary", zap.Error(err))
		return fmt.Errorf("failed to delete image: %v", err)
	}

	logs.Info("Image deleted successfully", zap.String("image_url", image.ImageURL))
	return nil
}

func extractPublicID(imageURL string) string {
	parts := strings.Split(imageURL, "/")
	if len(parts) < 8 {
		return ""
	}
	return parts[len(parts)-1]
}

// ---------------- Menu ------------------------

// CreateMenuItem
func (r *menuRepositoryDB) CreateMenuItem(ctx context.Context, item MenuItem) (uuid.UUID, error) {
	// สร้าง UUID ใหม่สำหรับเมนู
	newID, err := uuid.NewV4()
	if err != nil {
		logs.Error("Failed to generate UUID", zap.Error(err))
		return uuid.Nil, err
	}

	// กำหนดค่า ID
	item.UUID = newID

	// บันทึกข้อมูลลงฐานข้อมูล
	if err := r.db.WithContext(ctx).Create(&item).Error; err != nil {
		logs.Error("Failed to create menu item", zap.Error(err), zap.Any("MenuItem", item))
		return uuid.Nil, err
	}

	// Log รายละเอียดเมนูที่บันทึก
	logs.Info("Menu item created successfully",
		zap.String("ID", item.UUID.String()),
		zap.String("NameTH", item.NameTH),
		zap.String("NameEN", item.NameEN),
		zap.Float64("Price", item.Price),
		zap.String("Category", string(item.Category)),
		zap.String("Images", string(item.ImageURL)),
	)

	return item.UUID, nil
}

// UpdateMenuItem
func (r *menuRepositoryDB) UpdateMenuItem(ctx context.Context, item MenuItem) error {
	if err := r.db.WithContext(ctx).Save(&item).Error; err != nil {
		logs.Error("Failed to update menu item", zap.Error(err), zap.Any("MenuItem", item))
		return err
	}
	return nil
}

// DeleteMenuItem
func (r *menuRepositoryDB) DeleteMenuItem(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&MenuItem{}, "uuid = ?", id).Error; err != nil {
		logs.Error("Failed to delete menu item", zap.Error(err))
		return err
	}
	return nil
}

// GetMenuItems
func (r *menuRepositoryDB) GetMenuItems(ctx context.Context) ([]MenuItem, error) {
	var menuItems []MenuItem

	if err := r.db.WithContext(ctx).Find(&menuItems).Error; err != nil {
		logs.Error("Failed to get menu items", zap.Error(err))
		return nil, fmt.Errorf("failed to get menu items: %v", err)
	}

	return menuItems, nil
}

// GetMenuItemByID
func (r *menuRepositoryDB) GetMenuItemByID(ctx context.Context, id uuid.UUID) (MenuItem, error) {
	var menuItem MenuItem
	if err := r.db.WithContext(ctx).First(&menuItem, "uuid = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return MenuItem{}, fmt.Errorf("menu item not found with ID %s", id)
		}
		return MenuItem{}, fmt.Errorf("failed to get menu item: %v", err)
	}

	return menuItem, nil
}

// ---------------- Menu Set ------------------------

// CreateMenuSet สร้างเมนูเซตใหม่
func (r *menuRepositoryDB) CreateMenuSet(ctx context.Context, menuSet MenuSet) (uuid.UUID, error) {
	// สร้าง UUID ใหม่สำหรับ MenuSet
	newID, err := uuid.NewV4()
	if err != nil {
		logs.Error("Failed to generate UUID", zap.Error(err))
		return uuid.Nil, err
	}

	// กำหนด ID ของ MenuSet
	menuSet.UUID = newID

	// บันทึกข้อมูล MenuSet ลงฐานข้อมูล
	if err := r.db.WithContext(ctx).Create(&menuSet).Error; err != nil {
		logs.Error("Failed to create menu set", zap.Error(err), zap.Any("MenuSet", menuSet))
		return uuid.Nil, err
	}

	// Log รายละเอียดของ MenuSet ที่บันทึก
	logs.Info("Menu set created successfully", zap.String("ID", menuSet.UUID.String()), zap.String("SetName", menuSet.Name))

	return menuSet.UUID, nil
}

// UpdateMenuSet อัปเดต MenuSet
func (r *menuRepositoryDB) UpdateMenuSet(ctx context.Context, menuSet MenuSet) error {
	if err := r.db.WithContext(ctx).Save(&menuSet).Error; err != nil {
		logs.Error("Failed to update menu set", zap.Error(err), zap.Any("MenuSet", menuSet))
		return err
	}
	return nil
}

// DeleteMenuSet ลบ MenuSet
func (r *menuRepositoryDB) DeleteMenuSet(ctx context.Context, setID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("uuid = ?", setID).Delete(&MenuSet{}).Error; err != nil {
		logs.Error("Failed to delete menu set", zap.Error(err))
		return err
	}
	return nil
}

// GetMenuSets ดึงข้อมูลทั้งหมดของ MenuSets
func (r *menuRepositoryDB) GetMenuSets(ctx context.Context) ([]MenuSet, error) {
	var menuSet []MenuSet
	if err := r.db.WithContext(ctx).Find(&menuSet).Error; err != nil {
		logs.Error("Failed to get menu sets", zap.Error(err))
		return nil, fmt.Errorf("failed to get menu sets: %v", err)
	}
	return menuSet, nil
}

// GetMenuSetByID ดึงข้อมูล MenuSet ตาม ID
func (r *menuRepositoryDB) GetMenuSetByID(ctx context.Context, setID uuid.UUID) (MenuSet, error) {
	var menuSet MenuSet
	if err := r.db.WithContext(ctx).Where("uuid = ?", setID).First(&menuSet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return MenuSet{}, fmt.Errorf("menu set not found")
		}
		return MenuSet{}, fmt.Errorf("failed to get menu set: %v", err)
	}
	return menuSet, nil
}

// ---------------- Menu Set Item ------------------------

func (r *menuRepositoryDB) CreateMenuSetItems(ctx context.Context, menuSetID uuid.UUID, menuItemIDs []uuid.UUID) error {
	if menuSetID == uuid.Nil || len(menuItemIDs) == 0 {
		logs.Error("Invalid input: MenuSetID or MenuItemIDs")
		return fmt.Errorf("invalid input: MenuSetID or MenuItemIDs")
	}

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		logs.Error("Failed to begin transaction", zap.Error(err))
		return fmt.Errorf("transaction error: %v", err)
	}

	for _, menuItemID := range menuItemIDs {
		menuSetItem := MenuSetItem{
			MenuSetID:  menuSetID,
			MenuItemID: menuItemID,
		}
		if err := tx.Create(&menuSetItem).Error; err != nil {
			tx.Rollback()
			logs.Error("Failed to create menu set item", zap.Error(err), zap.Any("MenuSetItem", menuSetItem))
			return fmt.Errorf("failed to create menu set item: %v", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		logs.Error("Failed to commit transaction", zap.Error(err))
		return fmt.Errorf("transaction commit error: %v", err)
	}

	logs.Info("Menu set items created successfully", zap.String("MenuSetID", menuSetID.String()), zap.Any("MenuItemIDs", menuItemIDs))
	return nil
}

func (r *menuRepositoryDB) UpdateMenuSetItems(ctx context.Context, menuSetID uuid.UUID, menuItemIDs []uuid.UUID) error {
	if menuSetID == uuid.Nil || len(menuItemIDs) == 0 {
		logs.Error("Invalid input: MenuSetID or MenuItemIDs")
		return fmt.Errorf("invalid input: MenuSetID or MenuItemIDs")
	}

	// เริ่มต้น transaction
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		logs.Error("Failed to begin transaction", zap.Error(err))
		return fmt.Errorf("transaction error: %v", err)
	}

	// ลบรายการที่มีอยู่ก่อน (ถ้ามี) เพื่ออัปเดตเป็นรายการใหม่
	if err := tx.Where("menu_set_id = ?", menuSetID).Delete(&MenuSetItem{}).Error; err != nil {
		tx.Rollback()
		logs.Error("Failed to delete existing menu set items", zap.Error(err))
		return fmt.Errorf("failed to delete existing menu set items: %v", err)
	}

	// สร้างรายการใหม่
	for _, menuItemID := range menuItemIDs {
		menuSetItem := MenuSetItem{
			MenuSetID:  menuSetID,
			MenuItemID: menuItemID,
		}
		if err := tx.Create(&menuSetItem).Error; err != nil {
			tx.Rollback()
			logs.Error("Failed to create menu set item", zap.Error(err), zap.Any("MenuSetItem", menuSetItem))
			return fmt.Errorf("failed to create menu set item: %v", err)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		logs.Error("Failed to commit transaction", zap.Error(err))
		return fmt.Errorf("transaction commit error: %v", err)
	}

	logs.Info("Menu set items updated successfully", zap.String("MenuSetID", menuSetID.String()), zap.Any("MenuItemIDs", menuItemIDs))
	return nil
}

// DeleteMenuSetItem ลบ MenuSetItem ตาม menuSetID
func (r *menuRepositoryDB) DeleteMenuSetItem(ctx context.Context, menuSetID uuid.UUID) error {
	// ลบทุกแถวที่มี menu_set_id เดียวกันในตาราง menu_set_items
	if err := r.db.WithContext(ctx).Delete(&MenuSetItem{}, "menu_set_id = ?", menuSetID).Error; err != nil {
		logs.Error("Failed to delete menu set items", zap.Error(err), zap.String("MenuSetID", menuSetID.String()))
		return fmt.Errorf("failed to delete menu set items: %v", err)
	}

	logs.Info("Menu set items deleted successfully", zap.String("MenuSetID", menuSetID.String()))
	return nil
}

// GetMenuSetItems ดึงข้อมูลทั้งหมด
func (r *menuRepositoryDB) GetMenuSetItems(ctx context.Context) ([]MenuSetItemDetails, error) {
	var result []MenuSetItemDetails

	if err := r.db.WithContext(ctx).
		Table("menu_sets ms").
		Select("ms.uuid AS menu_set_id",
			"ms.name AS menu_set_name",
			"mi.uuid AS menu_item_id",
			"ms.price AS menu_set_price",
			"mi.name_th AS menu_name_th",
			"mi.name_en AS menu_name_en",
			"mi.price AS menu_price",
			"mi.category AS menu_category",
			"mi.image_url AS image_url").
		Joins("JOIN menu_set_items msi ON ms.uuid = msi.menu_set_id").
		Joins("JOIN menu_items mi ON msi.menu_item_id = mi.uuid").
		Scan(&result).Error; err != nil {
		logs.Error("Failed to get menu set items", zap.Error(err))
		return nil, fmt.Errorf("failed to get menu set items: %v", err)
	}

	return result, nil
}

// GetMenuSetItem By MenuSetID
func (r *menuRepositoryDB) GetMenuSetItemByMenuSetID(ctx context.Context, menuSetID uuid.UUID) ([]MenuSetItemDetails, error) {
	var result []MenuSetItemDetails
	if err := r.db.WithContext(ctx).
		Table("menu_sets ms").
		Select("ms.uuid AS menu_set_id",
			"mi.uuid AS menu_item_id",
			"ms.name AS menu_set_name",
			"ms.price AS menu_set_price",
			"mi.name_th AS menu_name_th",
			"mi.name_en AS menu_name_en",
			"mi.price AS menu_price",
			"mi.category AS menu_category",
			"mi.image_url AS image_url").
		Joins("JOIN menu_set_items msi ON ms.uuid = msi.menu_set_id").
		Joins("JOIN menu_items mi ON msi.menu_item_id = mi.uuid").
		Where("ms.uuid = ?", menuSetID).
		Scan(&result).Error; err != nil {
		logs.Error("Failed to get menu set items by menu set ID", zap.Error(err))
		return nil, fmt.Errorf("failed to get menu set items by menu set ID: %v", err)
	}

	return result, nil
}
