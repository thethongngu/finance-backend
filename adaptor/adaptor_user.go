package adaptor

import (
	"database/sql"
	"fmt"
)

type User struct {
	UserID         int    `json:"user_id"`
	Username       string `json:"username"`
	HashedPassword string `json:"hash_password"`
	AvatarURL      string `json:"avatar_url"`
}

type Session struct {
	SessionID string `json:"session_id"`
	UserID    int    `json:"user_id"`
}

type UserAdaptorInterface interface {
	GetUserByID(userID int) (*User, error)
	GetUserByCredentails(username string, password string) (*User, error)
	GetUserBySessionID(sessionID string) (*User, error)
	CreateNewUserSession(user *User) (*Session, error)
}

type UserMySQLAdaptor struct {
	conn *sql.DB
}

func NewUserMySQLAdaptor() UserMySQLAdaptor {
	return UserMySQLAdaptor{conn: GetMySQLConnection()}
}

func (mysql UserMySQLAdaptor) GetUserByID(userID int) (*User, error) {
	var user User
	err := mysql.conn.
		QueryRow(`SELECT user_id, username, password, avatar_url FROM user WHERE user_id = ?`, userID).
		Scan(&user.UserID, &user.Username, &user.HashedPassword, &user.AvatarURL)
	if err != nil {
		err = fmt.Errorf("[Error] GetUserByID: %v", err)
		return nil, err
	}
	return &user, nil
}

func (mysql UserMySQLAdaptor) GetUserByCredentails(username string, hashPassword string) (*User, error) {
	var user User
	err := mysql.conn.
		QueryRow(`SELECT user_id, username, avatar_url FROM user WHERE username = ? AND password = ?`, username, hashPassword).
		Scan(&user.UserID, &user.Username, &user.AvatarURL)
	if err != nil {
		fmt.Printf("[Error] GetUserByCredentails: %v\n", err)
		return nil, err
	}
	return &user, nil
}

func (mysql UserMySQLAdaptor) GetUserBySessionID(sessionID string) (*User, error) {
	var userID int
	err := mysql.conn.
		QueryRow(`SELECT user_id FROM session WHERE session_id = ?`, sessionID).Scan(&userID)
	if err != nil {
		err = fmt.Errorf("[Error] GetUserBySessionID: %v", err)
		return nil, err
	}

	return mysql.GetUserByID(userID)
}

func (mysql UserMySQLAdaptor) CreateNewUserSession(user *User) (*Session, error) {
	_, err := mysql.conn.Exec(`INSERT INTO session (user_id) VALUES (?)`, user.UserID)
	if err != nil {
		err = fmt.Errorf("[Error] CreateNewUserSession: %v", err)
		return nil, err
	}

	var session Session
	err = mysql.conn.QueryRow(`SELECT session_id, user_id FROM session WHERE user_id = (?)`, user.UserID).
		Scan(&session.SessionID, &session.UserID)
	if err != nil {
		err = fmt.Errorf("[Error] CreateNewUserSession: %v", err)
		return nil, err
	}

	return &session, nil
}
