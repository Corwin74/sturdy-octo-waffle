package service

import (
	"context"
	"fmt"
	v1 "shop/api/shop/v1"

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
	info, err := s.userUsecase.Info(ctx)
	fmt.Println(info, err)
	return &v1.InfoResponse{Coins: 50, Inventory: []*v1.InventoryItem{},
	CoinHistory: nil}, nil
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
	fmt.Println(err)
	if err != nil {
		return &v1.BaseResponse{Error: "500"}, ErrInternal
	}
	outMessage := fmt.Sprintf("%v, %v", in.Amount, in.ToUser)
	return &v1.BaseResponse{Error: outMessage}, nil
}

// BuyItem -- Купить предмет за монеты
func (s *ShopService) BuyItem(ctx context.Context, in *v1.Item) (*v1.BaseResponse, error) {
	if in.Name == "" {
		return &v1.BaseResponse{Error: "400"}, ErrBadRequest
	}
	err := s.userUsecase.Buy(ctx, in.Name)
	if err != nil {
		return &v1.BaseResponse{Error: "400"}, ErrBadRequest
	}
	return &v1.BaseResponse{Error: "Cool!"}, nil
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
