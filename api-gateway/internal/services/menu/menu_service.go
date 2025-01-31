package services

import (
	"context"
	"time"

	"gitlab.com/final_project1240930/api_gateway/internal/logs"
	"go.uber.org/zap"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type MenuService interface {
	// Handle Menu
	CreateMenuItem(ctx context.Context, req *CreateMenuItemRequest) (*CreateMenuItemResponse, error)
	UpdateMenuItem(ctx context.Context, req *UpdateMenuItemRequest) (*UpdateMenuItemResponse, error)
	DeleteMenuItem(ctx context.Context, req *DeleteMenuItemRequest) (*DeleteMenuItemResponse, error)
	GetMenuItems(ctx context.Context, req *emptypb.Empty) (*MenuItemList, error)
	GetMenuItemById(ctx context.Context, req *GetMenuItemByIdRequest) (*MenuItem, error)

	// Handle Menu Set
	CreateMenuSet(ctx context.Context, req *CreateMenuSetRequest) (*CreateMenuSetResponse, error)
	UpdateMenuSet(ctx context.Context, req *UpdateMenuSetRequest) (*UpdateMenuSetResponse, error)
	DeleteMenuSet(ctx context.Context, req *DeleteMenuSetRequest) (*DeleteMenuSetResponse, error)
	GetMenuSets(ctx context.Context, req *emptypb.Empty) (*MenuSetList, error)
	GetMenuSetById(ctx context.Context, req *GetMenuSetByIdRequest) (*MenuSet, error)

	// Handle Menu Set Item
	CreateMenuSetItem(ctx context.Context, req *CreateMenuSetItemRequest) (*CreateMenuSetItemResponse, error)
	GetMenuSetItems(ctx context.Context, req *emptypb.Empty) (*MenuSetItemList, error)
	GetMenuSetItemByMenuSetID(ctx context.Context, req *GetMenuSetItemByIdRequest) (*MenuSetItemList, error)
	UpdateMenuSetItem(ctx context.Context, req *UpdateMenuSetItemRequest) (*UpdateMenuSetItemResponse, error)
	DeleteMenuSetItem(ctx context.Context, req *DeleteMenuSetItemRequest) (*DeleteMenuSetItemResponse, error)

	// Handle Upload Img On Cloud
	UploadImage(ctx context.Context, req *UploadImageRequest) (*UploadImageResponse, error)
	DeleteImage(ctx context.Context, req *DeleteImageRequest) (*DeleteImageResponse, error)
}
type menuService struct {
	menuClient MenuServiceClient
}

func NewMenuService(menuClient MenuServiceClient) MenuService {
	return &menuService{menuClient: menuClient}
}

func (s *menuService) createWithTimeout(ctx context.Context, call func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	res, err := call(ctx)
	if err != nil {
		logs.Error("Error: %v", zap.Error(err))
		return nil, err
	}
	return res, nil
}

// Handle Menu
func (s *menuService) CreateMenuItem(ctx context.Context, req *CreateMenuItemRequest) (*CreateMenuItemResponse, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.menuClient.CreateMenuItem(ctx, req)
	})
	if res != nil {
		return res.(*CreateMenuItemResponse), nil
	}
	return nil, err
}

func (s *menuService) UpdateMenuItem(ctx context.Context, req *UpdateMenuItemRequest) (*UpdateMenuItemResponse, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.menuClient.UpdateMenuItem(ctx, req)
	})
	if res != nil {
		return res.(*UpdateMenuItemResponse), nil
	}
	return nil, err
}

func (s *menuService) DeleteMenuItem(ctx context.Context, req *DeleteMenuItemRequest) (*DeleteMenuItemResponse, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.menuClient.DeleteMenuItem(ctx, req)
	})
	if res != nil {
		return res.(*DeleteMenuItemResponse), nil
	}
	return nil, err
}

func (s *menuService) GetMenuItems(ctx context.Context, req *emptypb.Empty) (*MenuItemList, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.menuClient.GetMenuItems(ctx, req)
	})
	if res != nil {
		return res.(*MenuItemList), nil
	}
	return nil, err
}

func (s *menuService) GetMenuItemById(ctx context.Context, req *GetMenuItemByIdRequest) (*MenuItem, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.menuClient.GetMenuItemById(ctx, req)
	})
	if res != nil {
		return res.(*MenuItem), nil
	}
	return nil, err
}

// Handle Menu Set
func (s *menuService) CreateMenuSet(ctx context.Context, req *CreateMenuSetRequest) (*CreateMenuSetResponse, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.menuClient.CreateMenuSet(ctx, req)
	})
	if res != nil {
		return res.(*CreateMenuSetResponse), nil
	}
	return nil, err
}

func (s *menuService) UpdateMenuSet(ctx context.Context, req *UpdateMenuSetRequest) (*UpdateMenuSetResponse, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.menuClient.UpdateMenuSet(ctx, req)
	})
	if res != nil {
		return res.(*UpdateMenuSetResponse), nil
	}
	return nil, err
}

func (s *menuService) DeleteMenuSet(ctx context.Context, req *DeleteMenuSetRequest) (*DeleteMenuSetResponse, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.menuClient.DeleteMenuSet(ctx, req)
	})
	if res != nil {
		return res.(*DeleteMenuSetResponse), nil
	}
	return nil, err
}

func (s *menuService) GetMenuSets(ctx context.Context, req *emptypb.Empty) (*MenuSetList, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.menuClient.GetMenuSets(ctx, req)
	})
	if res != nil {
		return res.(*MenuSetList), nil
	}
	return nil, err
}

func (s *menuService) GetMenuSetById(ctx context.Context, req *GetMenuSetByIdRequest) (*MenuSet, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.menuClient.GetMenuSetById(ctx, req)
	})
	if res != nil {
		return res.(*MenuSet), nil
	}
	return nil, err
}

// Handle Menu Set Item
func (s *menuService) CreateMenuSetItem(ctx context.Context, req *CreateMenuSetItemRequest) (*CreateMenuSetItemResponse, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.menuClient.CreateMenuSetItem(ctx, req)
	})
	if res != nil {
		return res.(*CreateMenuSetItemResponse), nil
	}
	return nil, err
}

func (s *menuService) GetMenuSetItems(ctx context.Context, req *emptypb.Empty) (*MenuSetItemList, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.menuClient.GetMenuSetItems(ctx, req)
	})
	if res != nil {
		return res.(*MenuSetItemList), nil
	}
	return nil, err
}

func (s *menuService) GetMenuSetItemByMenuSetID(ctx context.Context, req *GetMenuSetItemByIdRequest) (*MenuSetItemList, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.menuClient.GetMenuSetItemByMenuSetID(ctx, req)
	})
	if res != nil {
		return res.(*MenuSetItemList), nil
	}
	return nil, err
}

func (s *menuService) UpdateMenuSetItem(ctx context.Context, req *UpdateMenuSetItemRequest) (*UpdateMenuSetItemResponse, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.menuClient.UpdateMenuSetItem(ctx, req)
	})
	if res != nil {
		return res.(*UpdateMenuSetItemResponse), nil
	}
	return nil, err
}

func (s *menuService) DeleteMenuSetItem(ctx context.Context, req *DeleteMenuSetItemRequest) (*DeleteMenuSetItemResponse, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.menuClient.DeleteMenuSetItem(ctx, req)
	})
	if res != nil {
		return res.(*DeleteMenuSetItemResponse), nil
	}
	return nil, err
}

// Handle Upload Img On Cloud
func (s *menuService) UploadImage(ctx context.Context, req *UploadImageRequest) (*UploadImageResponse, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.menuClient.UploadImage(ctx, req)
	})
	if res != nil {
		return res.(*UploadImageResponse), nil
	}
	return nil, err
}

func (s *menuService) DeleteImage(ctx context.Context, req *DeleteImageRequest) (*DeleteImageResponse, error) {
	res, err := s.createWithTimeout(ctx, func(ctx context.Context) (interface{}, error) {
		return s.menuClient.DeleteImage(ctx, req)
	})
	if res != nil {
		return res.(*DeleteImageResponse), nil
	}
	return nil, err
}
