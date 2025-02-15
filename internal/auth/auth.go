// package auth

// import (
// 	"context"
// 	"errors"
// 	"shop/internal/models"
// 	user_repo "shop/internal/repository/user"
// 	"strings"

// 	"github.com/go-kratos/kratos/v2/metadata"
// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/google/uuid"
// )

// var (
//     ErrInvalidToken = errors.New("invalid token")
//     ErrNoToken      = errors.New("no token provided")
//     ErrExpiredToken = errors.New("token expired")
// 	ErrNotFound = errors.New("user not auth")
// )

// type UserClaims struct {
//     jwt.RegisteredClaims
//     UserID uuid.UUID `json:"user_id"`
// }

// type Auther struct {
//     userRepo UserRepo
//     secret   string
// }

// // New создает новый экземпляр Auther
// func NewAuther(userRepo UserRepo, secret string) *Auther {
//     return &Auther{
//         userRepo: userRepo,
//         secret:   secret,
//     }
// }

// // extractToken извлекает токен из заголовка Authorization
// func extractToken(authorization string) (string, error) {
//     if authorization == "" {
//         return "", ErrNoToken
//     }

//     parts := strings.Split(authorization, " ")
//     if len(parts) != 2 || parts[0] != "Bearer" {
//         return "", ErrInvalidToken
//     }

//     return parts[1], nil
// }

// // IsAuth проверяет валидность токена и возвращает пользователя
// func (a *Auther) IsAuth(ctx context.Context) (models.User, error) {
//     md, ok := metadata.FromServerContext(ctx)
//     if !ok {
//         return models.User{}, ErrNoToken
//     }

//     authorization := md.Get("authorization")
//     tokenString, err := extractToken(authorization)
//     if err != nil {
//         return models.User{}, err
//     }

//     token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//         if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//             return nil, ErrInvalidToken
//         }
//         return []byte(a.secret), nil
//     })

//     if err != nil {
//         if errors.Is(err, jwt.ErrTokenExpired) {
//             return models.User{}, ErrExpiredToken
//         }
//         return models.User{}, ErrInvalidToken
//     }

//     claims, ok := token.Claims.(jwt.MapClaims)
//     if !ok || !token.Valid {
//         return models.User{}, ErrInvalidToken
//     }

//     // Проверяем наличие необходимых полей
//     userIDStr, ok := claims["id"].(string)
//     if !ok {
//         return models.User{}, ErrInvalidToken
//     }

//     userID, err := uuid.Parse(userIDStr)
//     if err != nil {
//         return models.User{}, ErrInvalidToken
//     }

//     // Получаем пользователя из БД
//     user, err := a.userRepo.Get(ctx, user_repo.Filter{ID: &userID})
//     if err != nil {
//         return models.User{}, err
//     }

//     return user, nil
// }
