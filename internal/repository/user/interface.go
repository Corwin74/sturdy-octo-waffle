package user

import "github.com/google/uuid"

type Filter struct {
	ID       *uuid.UUID
	Username *string
}

type Update struct {
	Balance *uint
}

type GetOptions struct {
	ForUpdate bool
}
