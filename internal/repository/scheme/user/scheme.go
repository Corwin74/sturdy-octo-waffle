package scheme_user

import (
	"fmt"
	"shop/internal/models"

	"github.com/google/uuid"
)

const (
	Table = "users"

	ID = "id"
	Name = "name"
	Password = "password"
	Balance = "balance"
)


type User struct {
	ID string
	Name string
	Password string
	Balance int
}


func ConvertToDBModel(user models.User) User {
	return User{
		ID: user.ID.String(),
		Name: user.Name,
		Password:  user.Password,
		Balance: int(user.Balance),
	}
}


func (u User) ConvertToDomainModel() (models.User, error) {
	id, err := uuid.Parse(u.ID)
	if err != nil {
		return models.User{}, fmt.Errorf("parsing id: %w", err)
	}

	return models.User{
		ID: id,
		Name: u.Name,
		Password: u.Password,
		Balance: uint(u.Balance),
	}, nil
}
