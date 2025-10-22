package storage

import (
	"errors"
	"fmt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	UserRole string `json:"user_role"`
}

var (
	ErrUserExists        = errors.New("user exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrIncorrectPassowrd = errors.New("incorrect password")
)

func (s *Storage) SaveUser(username, password, userRole string) error {
	const op = "storage.postgres.SaveUser"

	_, err := s.db.Exec(`
	INSERT INTO users (username, password, user_role)
	VALUES ($1, $2, $3)
`, username, password, userRole)

	if err != nil {
		if isDuplicateError(err) {
			return ErrUserExists
		}
		return ErrUserNotFound
	}

	return nil
}

func (s *Storage) GetRecordByUsername(username string) (User, error) {
	var record User
	const op = "storage.postgres.GetRecordByUsername"

	row := s.db.QueryRow(`
	SELECT * FROM users
	WHERE username = $1
`, username)
	if err := row.Scan(&record.ID, &record.Username, &record.Password, &record.UserRole); err != nil {
		return User{}, fmt.Errorf("%s: failed to parse string: %w", op, err)
	}

	return record, nil
}

func (s *Storage) GetAllUsers() (*[]User, error) {
	const op = "storage.postgres.GetAllUsers"

	rows, err := s.db.Query(`
	SELECT * FROM users
`)
	if err != nil {
		return &[]User{}, fmt.Errorf("%s: failed to parse users: %w", op, err)
	}
	defer rows.Close()
	users := []User{}
	for rows.Next() {
		userGotten := User{}
		err := rows.Scan(&userGotten.ID, &userGotten.Username, &userGotten.Password, &userGotten.UserRole)
		if err != nil {
			fmt.Println(fmt.Errorf("%s: failed to parse user: %w", op, err))
			continue
		}
		users = append(users, userGotten)
	}

	return &users, nil
}

func (s *Storage) DeleteUser(username string) error {
	const op = "storage.postgres.DeleteUser"

	result, err := s.db.Exec(`
	DELETE FROM users
	WHERE username = $1
`, username)
	if err != nil {
		return fmt.Errorf("%s: failed to delete user: %w", op, err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (s *Storage) AuthenticateUser(username, password string) (User, error) {
	const op = "storage.postgres.AuthenticateUser"

	user, err := s.GetRecordByUsername(username)
	if err != nil {
		return User{}, ErrUserNotFound
	}
	if user.Username == username && user.Password == password {
		return user, nil
	} else {
		return User{}, ErrIncorrectPassowrd
	}
}
