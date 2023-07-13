package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/sillamilla/user_microservice/internal/users/Mongo_storage"
	"github.com/sillamilla/user_microservice/internal/users/Redis_storage"
	"github.com/sillamilla/user_microservice/internal/users/model"
	"github.com/sillamilla/user_microservice/internal/users/service/helper"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Service interface {
	SignUp(ctx context.Context, input model.Input) (model.User, error)
	SignIn(ctx context.Context, input model.Input) (model.User, error)

	Logout(ctx context.Context, session string) error

	GetBySession(ctx context.Context, session string) (model.User, error)
	UpsertSessions(ctx context.Context, id string) (string, error)
	GetSession(ctx context.Context, id string) (string, error)

	EditProfile(ctx context.Context, id string, input model.UpdateUser) error
	EditPassword(ctx context.Context, id string, input model.ChangePassword) error

	GetByID(ctx context.Context, id string) (model.User, error)
	GetByUsername(ctx context.Context, username string) (model.User, error)
	SearchByUsername(ctx context.Context, username string) (model.UserInfo, error)
}

type service struct {
	re Redis_storage.Storage
	mo Mongo_storage.Storage
}

func New(re Redis_storage.Storage, mo Mongo_storage.Storage) Service {
	return &service{
		re: re,
		mo: mo,
	}
}

func (s *service) SignUp(ctx context.Context, input model.Input) (model.User, error) {
	user, err := s.GetByUsername(ctx, input.Username)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return model.User{}, errors.Wrap(err, "service.SignUp.GetByUsername")
	}
	if user.Username != "" {
		return model.User{}, errors.New("Username already taken")
	}

	password, err := helper.HashPassword(input.Password)
	if err != nil {
		return model.User{}, errors.Wrap(err, "service.SignUp.HashPassword")
	}
	input.Password = password

	id := uuid.NewString()
	session, err := s.UpsertSessions(ctx, id)
	if err != nil {
		return model.User{}, errors.Wrap(err, "service.SignUp.UpsertSessions")
	}
	//todo dont work upsert session

	newUser := model.UserFromInput(id, input, session, time.Now())
	err = s.mo.SignUp(ctx, newUser)
	if err != nil {
		return model.User{}, errors.Wrap(err, "service.SignUp")
	}

	context.WithValue(ctx, "user", user)

	return user, nil
}

func (s *service) SignIn(ctx context.Context, input model.Input) (model.User, error) {
	user, err := s.GetByUsername(ctx, input.Username)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.User{}, errors.New("User not found")
	} else if err != nil {
		return model.User{}, errors.Wrap(err, "service.SignIn.GetByUsername")
	}

	err = helper.ComparePassword(user.Password, input.Password)
	if err != nil {
		return model.User{}, errors.New("Invalid password")
	}

	input.Password = user.Password

	signUser, err := s.mo.SignIn(ctx, input)
	if err != nil {
		return model.User{}, errors.Wrap(err, "service.SignIn")
	}

	err = s.re.UpsertSession(ctx, signUser.ID, signUser.Session)
	if err != nil {
		return model.User{}, errors.Wrap(err, "service.SignIn.UpsertSession")
	}

	context.WithValue(ctx, "user", signUser)

	return signUser, nil
}

func (s *service) EditProfile(ctx context.Context, id string, input model.UpdateUser) error {
	user, err := s.GetByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "service.EditProfile.GetByID")
	}
	if user == (model.User{}) {
		return errors.New("User not found")
	}

	searchUser, err := s.GetByUsername(ctx, input.Username)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return errors.Wrap(err, "service.EditProfile.GetByUsername")
	}
	if searchUser.Username != "" {
		return errors.New("Username is already taken")
	}

	err = s.mo.EditProfile(ctx, id, input)
	if err != nil {
		return errors.Wrap(err, "service.EditProfile")
	}

	return nil
}

func (s *service) EditPassword(ctx context.Context, id string, input model.ChangePassword) error {
	user, err := s.GetByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "service.EditPassword.GetByID")
	}

	err = helper.ComparePassword(user.Password, input.Old)
	if err != nil {
		return errors.New("Invalid password")
	}

	if input.New != input.ConfirmNew {
		return errors.New("New password and confirm new password do not match")
	}

	password, err := helper.HashPassword(input.New)
	if err != nil {
		return errors.Wrap(err, "service.EditPassword.HashPassword")
	}

	err = s.mo.EditPassword(ctx, id, password)
	if err != nil {
		return errors.Wrap(err, "service.EditPassword")
	}

	return nil
}

func (s *service) Logout(ctx context.Context, id string) error {
	err := s.re.Logout(ctx, id)
	if err != nil {
		return errors.Wrap(err, "service.Logout")
	}

	return nil
}

func (s *service) GetByID(ctx context.Context, id string) (model.User, error) {
	byID, err := s.mo.GetByID(ctx, id)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.User{}, errors.New("User not found")
	} else if err != nil {
		return model.User{}, errors.Wrap(err, "service.GetByID")
	}

	return byID, nil
}

func (s *service) GetByUsername(ctx context.Context, username string) (model.User, error) {
	byUsername, err := s.mo.GetByUsername(ctx, username)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.User{}, errors.Wrap(err, "User not found")
	} else if err != nil {
		return model.User{}, errors.Wrap(err, "service.searchByUsername")
	}

	return byUsername, nil
}

func (s *service) SearchByUsername(ctx context.Context, username string) (model.UserInfo, error) {
	byUsername, err := s.mo.SearchByUsername(ctx, username)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.UserInfo{}, errors.New("User not found")
	} else if err != nil {
		return model.UserInfo{}, errors.Wrap(err, "service.searchByUsername")
	}

	return byUsername, nil
}

func (s *service) UpsertSessions(ctx context.Context, id string) (string, error) {
	sessionID, err := helper.HashPassword(uuid.New().String())
	if err != nil {
		return "", errors.Wrap(err, "service.SetSession.GenerateSessionID")
	}

	err = s.re.UpsertSession(ctx, id, sessionID)
	if err != nil {
		return "", errors.Wrap(err, "service.SetSession")
	}
	//todo перенести
	err = s.mo.UpsertSession(ctx, id, sessionID)
	if err != nil {
		return "", errors.Wrap(err, "service.SetSession")
	}

	context.WithValue(ctx, "session", sessionID)

	return sessionID, nil
}

func (s *service) GetBySession(ctx context.Context, session string) (model.User, error) {
	user, err := s.mo.GetBySession(ctx, session)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.User{}, errors.New("User not found")
	} else if err != nil {
		return model.User{}, errors.Wrap(err, "service.GetBySession")
	}

	return user, nil
}

func (s *service) GetSession(ctx context.Context, id string) (string, error) {
	session, err := s.re.GetSession(ctx, id)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			user, err := s.GetByID(ctx, id)
			if err != nil {
				return "", errors.Wrap(err, "service.GetSession.GetByID")
			}

			err = s.re.UpsertSession(ctx, id, user.Session)
			if err != nil {
				return "", errors.Wrap(err, "service.GetSession.UpsertSession")
			}

			return user.Session, nil
		} else {
			return "", errors.Wrap(err, "service.GetSession")
		}
	}

	return session, nil
}
