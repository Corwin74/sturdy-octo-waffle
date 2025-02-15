package service

import "github.com/go-kratos/kratos/v2/errors"

var (
	ErrUnauthorized = errors.New(401, "Unauthorized", "Неавторизован")
	ErrBadRequest = errors.New(400, "Bad Request", "Неверный запрос")
	ErrInternal = errors.New(500, "Internal Error", "Ошибка сервера")
)