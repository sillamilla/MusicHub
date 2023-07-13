package model

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Session   string    `json:"session"`
	Bio       string    `json:"bio"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"create_at"`
}

type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Icon     string `json:"icon"`
}

type UpdateUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
	Icon     string `json:"icon"`
}

type Input struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ChangePassword struct {
	Old        string `json:"old"`
	New        string `json:"new"`
	ConfirmNew string `json:"confirm_new"`
}

func UserFromInput(ID string, user Input, session string, createdAt time.Time) User {
	return User{
		ID:        ID,
		Username:  user.Username,
		Password:  user.Password,
		Session:   session,
		CreatedAt: createdAt,
		Email:     "",
		Bio:       "",
		Icon:      "",
	}
}
