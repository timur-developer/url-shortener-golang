package users

import (
	"go-to-do-checklist/internal/storage"
)

type UserStorage interface {
	SaveUser(username, password, userRole string) error
	GetRecordByUsername(username string) (storage.User, error)
	GetAllUsers() (*[]storage.User, error)
	DeleteUser(username string) error
	AuthenticateUser(username, password string) (storage.User, error)
}
