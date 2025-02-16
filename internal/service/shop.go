package service

import (
	"context"
	"errors"
	v1 "shop/api/shop/v1"
	"shop/internal/common"
)

// ShopService is a shop service.
type ShopService struct {
	v1.UnimplementedShopServer
	userUsecase UserUsecase
}

// NewShopService new a shop service.
func NewShopService(uu UserUsecase) *ShopService {
	return &ShopService{
		userUsecase: uu,
	}
}

// Info -- Получить информацию о монетах, инвентаре и истории транзакций
func (s *ShopService) Info(ctx context.Context, in *v1.InfoRequest) (*v1.InfoResponse, error) {
	userInfo, err := s.userUsecase.Info(ctx)
	if err != nil {
		if errors.Is(err, common.ErrUnauthorized) {
			return nil, ErrUnauthorized
		}
		return nil, ErrBadRequest
	}

    protoInventory := make([]*v1.InventoryItem, len(userInfo.Inventory))
    for i, item := range userInfo.Inventory {
        protoInventory[i] = &v1.InventoryItem{
            Type:     item.Type,
            Quantity: int32(item.Quantity),
        }
    }

	protoReceived := make([]*v1.ReceivedTransaction, len(userInfo.CoinHistory.Received))
	for i, received := range userInfo.CoinHistory.Received {
		protoReceived[i] = &v1.ReceivedTransaction{
			FromUser: received.FromUser,
			Amount:   int32(received.Amount),
		}
	}

	protoSent := make([]*v1.SentTransaction, len(userInfo.CoinHistory.Sent))
	for i, sent := range userInfo.CoinHistory.Sent {
		protoSent[i] = &v1.SentTransaction{
			ToUser: sent.ToUser,
			Amount: int32(sent.Amount),
		}
	}

	response := &v1.InfoResponse{
		Coins:     int32(userInfo.Coins),
		Inventory: protoInventory,
		CoinHistory: &v1.CoinHistoryDetails{
			Received: protoReceived,
			Sent:     protoSent,
		},
	}
	return response, nil
}

// SendCoin -- Отправить монеты другому пользователю
func (s *ShopService) SendCoin(ctx context.Context, in *v1.SentTransaction) (*v1.BaseResponse, error) {
	if in.Amount == 0 {
		return &v1.BaseResponse{Error: "400"}, ErrBadRequest
	}
	if in.ToUser == "" {
		return &v1.BaseResponse{Error: "400"}, ErrBadRequest
	}
	err := s.userUsecase.TransferCoins(ctx, in.ToUser, uint(in.Amount))
	if err != nil {
		if errors.Is(err, common.ErrUnauthorized) {
			return &v1.BaseResponse{Error: "401"}, ErrUnauthorized
		}
		return &v1.BaseResponse{Error: "400"}, ErrBadRequest
	}

	return &v1.BaseResponse{Error: ""}, nil
}

// BuyItem -- Купить предмет за монеты
func (s *ShopService) BuyItem(ctx context.Context, in *v1.Item) (*v1.BaseResponse, error) {
	if in.Name == "" {
		return &v1.BaseResponse{Error: "400"}, ErrBadRequest
	}
	err := s.userUsecase.Buy(ctx, in.Name)
	if err != nil {
		if errors.Is(err, common.ErrUnauthorized) {
			return &v1.BaseResponse{Error: "401"}, ErrUnauthorized
		}
		return &v1.BaseResponse{Error: "400"}, ErrBadRequest
	}
	return &v1.BaseResponse{Error: ""}, nil
}

// Auth -- Получение токена
func (s *ShopService) Auth(ctx context.Context, in *v1.AuthRequest) (*v1.AuthResponse, error) {
	if in.Password == "" || in.Username == "" {
		return &v1.AuthResponse{Data: &v1.AuthResponse_Error{Error: "400"}}, ErrBadRequest
	}
	
	token, err := s.userUsecase.Auth(ctx, in.Username, in.Password)
	if err != nil {
		return &v1.AuthResponse{Data: &v1.AuthResponse_Error{Error: "401"}}, ErrUnauthorized
	}

	return &v1.AuthResponse{Data: &v1.AuthResponse_Token{Token: token}}, nil
}
