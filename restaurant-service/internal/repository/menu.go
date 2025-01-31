package repository

import (
	"context"

	"github.com/gofrs/uuid"
)

type MenuCategory string

const (
	MainCourse MenuCategory = "MAIN_COURSE"
	Beverage   MenuCategory = "BEVERAGE"
	Dessert    MenuCategory = "DESSERT"
)

type MenuItem struct {
	UUID        uuid.UUID    `gorm:"column:uuid;type:uuid;default:gen_random_uuid();primaryKey" json:"menu_item_id"`
	NameTH      string       `gorm:"type:varchar(255);not null;unique" json:"name_th"`
	NameEN      string       `gorm:"type:varchar(255);not null;unique" json:"name_en"`
	Description string       `gorm:"type:text" json:"description"`
	Price       float64      `gorm:"type:decimal(10,2);not null" json:"price"`
	Category    MenuCategory `gorm:"type:menu_category;not null" json:"category"`
	ImageURL    string       `gorm:"type:varchar(255)" json:"image_url"`
}

type MenuSet struct {
	UUID  uuid.UUID `gorm:"column:uuid;type:uuid;default:gen_random_uuid();primaryKey"`
	Name  string    `gorm:"type:varchar(255);not null;unique"`
	Price float64   `gorm:"type:double precision;not null"`
}

type MenuSetItem struct {
	UUID       uuid.UUID `gorm:"column:uuid;type:uuid;default:gen_random_uuid();primaryKey"`
	MenuSetID  uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	MenuItemID uuid.UUID `gorm:"type:uuid;not null;primaryKey"`

	// Foreign Key
	MenuSet  MenuSet  `gorm:"foreignKey:MenuSetID;constraint:onDelete:CASCADE"`
	MenuItem MenuItem `gorm:"foreignKey:MenuItemID;constraint:onDelete:CASCADE"`
}
type Image struct {
	UUID     uuid.UUID `gorm:"column:uuid;type:uuid;default:gen_random_uuid();primaryKey"`
	ImageURL string    `gorm:"type:varchar(255);not null"`
}

type MenuSetItemDetails struct {
	MenuSetId    string  `json:"menu_set_id"`
	MenuItemId   string  `json:"menu_item_id"`
	MenuSetName  string  `json:"menu_set_name"`
	MenuSetPrice float64 `json:"menu_set_price"`
	MenuNameTh   string  `json:"menu_name_th"`
	MenuNameEn   string  `json:"menu_name_en"`
	MenuPrice    float64 `json:"menu_price"`
	MenuCategory string  `json:"menu_category"`
	ImageUrl     string  `json:"image_url"`
}

type MenuRepository interface {
	// Image Methods
	UploadImage(ctx context.Context, imageData []byte, fileName string) (string, error)
	DeleteImage(ctx context.Context, image Image) error

	// Menu Item Methods
	CreateMenuItem(ctx context.Context, item MenuItem) (uuid.UUID, error)
	UpdateMenuItem(ctx context.Context, item MenuItem) error
	DeleteMenuItem(ctx context.Context, id uuid.UUID) error
	GetMenuItems(ctx context.Context) ([]MenuItem, error)
	GetMenuItemByID(ctx context.Context, id uuid.UUID) (MenuItem, error)

	// Menu Set Methods
	CreateMenuSet(ctx context.Context, menuSet MenuSet) (uuid.UUID, error)
	UpdateMenuSet(ctx context.Context, set MenuSet) error
	DeleteMenuSet(ctx context.Context, id uuid.UUID) error
	GetMenuSets(ctx context.Context) ([]MenuSet, error)
	GetMenuSetByID(ctx context.Context, id uuid.UUID) (MenuSet, error)

	// // Menu Set Item Methods
	CreateMenuSetItems(ctx context.Context, menuSetID uuid.UUID, menuItemIDs []uuid.UUID) error
	UpdateMenuSetItems(ctx context.Context, menuSetID uuid.UUID, menuItemIDs []uuid.UUID) error
	DeleteMenuSetItem(ctx context.Context, uuid uuid.UUID) error
	GetMenuSetItems(ctx context.Context) ([]MenuSetItemDetails, error)
	GetMenuSetItemByMenuSetID(ctx context.Context, uuid uuid.UUID) ([]MenuSetItemDetails, error)
}
