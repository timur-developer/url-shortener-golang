package url

import "go-to-do-checklist/internal/storage"

type URLStorage interface {
	SaveURL(urlToSave string, alias string, user_id int) error
	GetURL(alias string, user_id int) (string, error)
	GetAllURLs(user_id int) (*[]storage.URL, error)
	DeleteURL(alias string, user_id int) error
	UpdateRecordPartly(alias, newAlias, newUrl string, user_id int) (storage.URL, error)
	GetRecordByURL(url string, user_id int) (storage.URL, error)
	GetRecordByAlias(alias string, user_id int) (storage.URL, error)
	UpdateRecordCompletely(alias, newAlias, newUrl string, user_id int) (storage.URL, error)
	GetAllURLsAdmin() (*[]storage.URL, error)
}
