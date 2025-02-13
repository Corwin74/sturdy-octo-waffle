package service

import (
	"context"
	"shop/api/shop/v1"
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

// // Info -- Получить информацию о монетах, инвентаре и истории транзакций
// func (s *ShopService) Info(ctx context.Context, in *v1.InfoRequest) (*v1.InfoResponse, error) {

// 	return &v1.InfoResponse{Coins: 50, Inventory: []*v1.InventoryItem{}, 
// 	CoinHistory: nil}, nil
// }

// // SendCoin -- Отправить монеты другому пользователю
// func (s *ShopService) SendCoin(ctx context.Context, in *v1.SentTransaction) (*v1.SuccessResponse, error) {
// 	if in.Amount == 0 {
// 		return &v1.SuccessResponse{Message: "Нулевые входные данные, караул!!!"}, errors.New("fhfhfhfh")
// 	}
// 	outMessage := fmt.Sprintf("%v, %v", in.Amount, in.ToUser)
// 	return &v1.SuccessResponse{Message: outMessage}, nil
// }


// // BuyItem -- Купить предмет за монеты
// func (s *ShopService) BuyItem(ctx context.Context, in *v1.Item) (*v1.SuccessResponse, error) {
// 	if in.Name == "" {
// 		return &v1.SuccessResponse{Message: "Нулевые входные данные, караул!!!"}, errors.New("fhfhfhfh")
// 	}
// 	outMessage := fmt.Sprintf("%v", in.Name)
// 	return &v1.SuccessResponse{Success: true, Message: outMessage}, nil
// }

// Auth -- Получение токена
func (s *ShopService) Auth(ctx context.Context, in *v1.AuthRequest) (*v1.AuthResponse, error) {
	token, err := s.userUsecase.Auth(ctx, in.Username, in.Password)
	if err != nil {
		return &v1.AuthResponse{Data: &v1.AuthResponse_Error{Error: err.Error()}}, nil
	}

	return &v1.AuthResponse{Data: &v1.AuthResponse_Token{Token: token}}, nil
}
