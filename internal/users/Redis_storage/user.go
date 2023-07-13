package Redis_storage

import (
	"github.com/sillamilla/user_microservice/internal/users/model"
	"time"
)

type UserStorage struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Bio       string    `json:"bio"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"create_at"`
}

func UserFromInput(ID string, user model.Input, createdAt time.Time) UserStorage {
	return UserStorage{
		ID:        ID,
		Username:  user.Username,
		Password:  user.Password,
		CreatedAt: createdAt,
		Email:     "",
		Bio:       "",
		Icon:      "",
	}
}
