package usecase

import (
	"finance/adaptor"
	"finance/entity"
)

type UserUsecaseInterface interface {
	LoginByPassword(username string, hashPassword string) (*entity.User, error)
	LoginBySession(sessionId string) (*entity.User, error)
	Logout(userID int)
}

type UserUsecase struct {
	userAdaptor adaptor.UserAdaptorInterface
}

func NewUserUsecase(adaptor adaptor.UserAdaptorInterface) *UserUsecase {
	return &UserUsecase{userAdaptor: adaptor}
}

func (u *UserUsecase) LoginBySession(sessionID string) (*entity.User, error) {
	user, err := u.userAdaptor.GetUserBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUsecase) LoginByPassword(username string, hashPassword string) (*entity.Session, error) {
	user, err := u.userAdaptor.GetUserByCredentails(username, hashPassword)
	if err != nil {
		return nil, err
	}

	var session *entity.Session
	session, err = u.userAdaptor.CreateNewUserSession(user)
	if err != nil {
		return nil, err
	}

	return session, nil
}
