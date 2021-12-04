package adaptor

import (
	"database/sql"
	"finance/entity"
	"fmt"
)

type UserAdaptorInterface interface {
	GetUserByID(userID int) (*entity.User, error)
	GetUserByCredentails(username string, password string) (*entity.User, error)
	GetUserBySessionID(sessionID string) (*entity.User, error)
	CreateNewUserSession(user *entity.User) (*entity.Session, error)
}

type UserMySQLAdaptor struct {
	conn *sql.DB
}

func NewUserMySQLAdaptor() UserMySQLAdaptor {
	return UserMySQLAdaptor{conn: GetMySQLConnection()}
}

func (mysql UserMySQLAdaptor) GetUserByID(userID int) (*entity.User, error) {
	var user entity.User
	err := mysql.conn.
		QueryRow(`SELECT user_id, name, hashed_password, avatar_url FROM user WHERE user_id = ?`, userID).
		Scan(&user.UserID, &user.Username, &user.HashedPassword, &user.AvatarURL)
	if err != nil {
		err = fmt.Errorf("[Error] Cannot query user: %v", err)
		return nil, err
	}
	return &user, nil
}

func (mysql UserMySQLAdaptor) GetUserByCredentails(username string, hashPassword string) (*entity.User, error) {
	var user *entity.User
	err := mysql.conn.
		QueryRow(`SELECT user_id, name, avatar_url FROM user WHERE username = ? AND password = ?`, username, hashPassword).
		Scan(user.UserID, user.Username, user.AvatarURL)
	if err != nil {
		err = fmt.Errorf("[Error] sql: %v", err)
		return nil, err
	}
	return user, nil
}

func (mysql UserMySQLAdaptor) GetUserBySessionID(sessionID string) (*entity.User, error) {
	var userID int
	err := mysql.conn.
		QueryRow(`SELECT user_id FROM session WHERE session_id = ?`, sessionID).Scan(&userID)
	if err != nil {
		err = fmt.Errorf("[Error] sql: %v", err)
		return nil, err
	}

	return mysql.GetUserByID(userID)
}

func (mysql UserMySQLAdaptor) CreateNewUserSession(user *entity.User) (*entity.Session, error) {
	_, err := mysql.conn.Exec(`INSERT INTO session (user_id) VALUES (?)`, user.UserID)
	if err != nil {
		return nil, err
	}

	var session *entity.Session
	err = mysql.conn.QueryRow(`SELECT session_id, user_id FROM session WHERE userID = (?)`, user.UserID).
		Scan(&session.SessionID, &session.UserID)
	if err != nil {
		return nil, err
	}

	return session, nil
}
