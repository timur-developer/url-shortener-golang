package storage

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

var (
	ErrURLExists         = errors.New("url exists")
	ErrURLNotFound       = errors.New("url not found")
	ErrAliasNotFound     = errors.New("alias not found")
	ErrIdNotFound        = errors.New("id not found")
	ErrWrongPatchRequest = errors.New("to make complete update use put request")
)

type URL struct {
	ID     int    `db:"id"`
	Alias  string `db:"alias"`
	URL    string `db:"url"`
	UserID int    `db:"user_id"`
}

func (s *Storage) SaveURL(urlToSave string, alias string, user_id int) error {
	const op = "storage.postgres.SaveURL"

	_, err := s.db.Exec(`
	INSERT INTO url(url, alias, user_id)
	VALUES ($1, $2, $3)`, urlToSave, alias, user_id)

	if err != nil {
		if isDuplicateError(err) {
			return ErrURLExists
		}
		return err
	}

	return nil
}

func (s *Storage) GetURL(alias string, user_id int) (string, error) {
	var urlGotten string
	const op = "storage.postgres.GetURL"

	row := s.db.QueryRow(`
	SELECT url FROM url
	WHERE alias = $1 AND user_id = $2
	`, alias, user_id)
	if err := row.Scan(&urlGotten); err != nil {
		if urlGotten == "" {
			return "", ErrAliasNotFound
		} else {
			return "", fmt.Errorf("%s: failed to parse string: %w", op, err)
		}
	}

	return urlGotten, nil
}

func (s *Storage) GetRecordByURL(url string, user_id int) (URL, error) {
	var recordGotten URL
	const op = "storage.postgres.GetRecordByURL"

	row := s.db.QueryRow(`
	SELECT * FROM url
	WHERE url = $1 AND user_id = $2
	`, url, user_id)
	if err := row.Scan(&recordGotten.ID, &recordGotten.Alias, &recordGotten.URL, &recordGotten.UserID); err != nil {
		return URL{}, fmt.Errorf("%s: failed to parse string: %w", op, err)
	}

	return recordGotten, nil
}

func (s *Storage) GetRecordByAlias(alias string, user_id int) (URL, error) {
	var recordGotten URL
	const op = "storage.postgres.GetRecordByAlias"

	row := s.db.QueryRow(`
	SELECT * FROM url
	WHERE alias = $1 AND user_id = $2
	`, alias, user_id)
	if err := row.Scan(&recordGotten.ID, &recordGotten.Alias, &recordGotten.URL, &recordGotten.UserID); err != nil {
		return URL{}, fmt.Errorf("%s: failed to parse string: %w", op, err)
	}

	return recordGotten, nil
}

func (s *Storage) DeleteURL(alias string, user_id int) error {
	const op = "storage.postgres.DeleteURL"

	result, err := s.db.Exec(`
	DELETE FROM url 
	WHERE alias = $1 AND user_id = $2
`, alias, user_id)
	if err != nil {
		return fmt.Errorf("%s: failed to delete url: %w", op, err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrAliasNotFound
	}
	return nil
}

func (s *Storage) GetAllURLs(user_id int) (*[]URL, error) {
	const op = "storage.postgres.GetAllURLs"

	rows, err := s.db.Query(`
	SELECT * FROM url
	WHERE user_id = $1
`, user_id)
	if err != nil {
		return &[]URL{}, fmt.Errorf("%s: failed to parse urls: %w", op, err)
	}
	defer rows.Close()
	urls := []URL{}
	for rows.Next() {
		urlGotten := URL{}
		err := rows.Scan(&urlGotten.ID, &urlGotten.Alias, &urlGotten.URL, &urlGotten.UserID)
		if err != nil {
			fmt.Println(fmt.Errorf("%s: failed to parse url: %w", op, err))
			continue
		}
		urls = append(urls, urlGotten)
	}
	return &urls, nil
}

func (s *Storage) GetAllURLsAdmin() (*[]URL, error) {
	const op = "storage.postgres.GetAllURLs"

	rows, err := s.db.Query(`
	SELECT * FROM url
`)
	if err != nil {
		return &[]URL{}, fmt.Errorf("%s: failed to parse urls: %w", op, err)
	}
	defer rows.Close()
	urls := []URL{}
	for rows.Next() {
		urlGotten := URL{}
		err := rows.Scan(&urlGotten.ID, &urlGotten.Alias, &urlGotten.URL, &urlGotten.UserID)
		if err != nil {
			fmt.Println(fmt.Errorf("%s: failed to parse url: %w", op, err))
			continue
		}
		urls = append(urls, urlGotten)
	}
	return &urls, nil
}

func (s *Storage) UpdateRecordPartly(alias, newAlias, newUrl string, user_id int) (URL, error) {
	const op = "storage.postgres.UpdateRecordPartly"

	var result sql.Result
	var err error

	if newAlias == "" && newUrl != "" {
		result, err = s.db.Exec(`
		UPDATE url set url = $1
		WHERE alias = $2 AND user_id = $3
	`, newUrl, alias, user_id)
		if err != nil {
			return URL{}, fmt.Errorf("%s: failed to update record: %w", op, err)
		}
	} else if newAlias != "" && newUrl == "" {
		result, err = s.db.Exec(`
		UPDATE url set alias = $1
		WHERE alias = $2 AND user_id = $3
	`, newAlias, alias, user_id)
		if err != nil {
			return URL{}, fmt.Errorf("%s: failed to update record: %w", op, err)
		}
		alias = newAlias
	} else {
		return URL{}, ErrWrongPatchRequest
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return URL{}, ErrAliasNotFound
	}
	var updatedRecord URL
	row := s.db.QueryRow(`
	SELECT * FROM url
	WHERE alias = $1 AND user_id = $2
`, alias, user_id)
	if err := row.Scan(&updatedRecord.ID, &updatedRecord.Alias, &updatedRecord.URL, &updatedRecord.UserID); err != nil {
		return URL{}, fmt.Errorf("%s: failed to get updated record: %w", op, err)
	}

	return updatedRecord, nil
}

func (s *Storage) UpdateRecordCompletely(alias, newAlias, newUrl string, user_id int) (URL, error) {
	const op = "storage.postgres.UpdateRecordCompletely"

	var result sql.Result
	var err error

	result, err = s.db.Exec(`
	UPDATE url set alias = $1,
	               url = $2
	WHERE alias = $3 AND user_id = $4
`, newAlias, newUrl, alias, user_id)
	if err != nil {
		return URL{}, fmt.Errorf("%s: failed to update record: %w", op, err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return URL{}, ErrAliasNotFound
	}

	var updatedRecord URL
	row := s.db.QueryRow(`
	SELECT * FROM url
	WHERE alias = $1 AND user_id = $2
`, newAlias, user_id)
	if err := row.Scan(&updatedRecord.ID, &updatedRecord.Alias, &updatedRecord.URL, &updatedRecord.UserID); err != nil {
		return URL{}, fmt.Errorf("%s: failed to get updated record: %w", op, err)
	}

	return updatedRecord, nil
}
