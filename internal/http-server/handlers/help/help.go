package help

import (
	"fmt"
	"github.com/go-chi/render"
	"net/http"
)

type EndpointInfo struct {
	Method          string      `json:"method"`
	Path            string      `json:"path"`
	Description     string      `json:"description"`
	AuthRequired    bool        `json:"auth_required"`
	AdminRequired   bool        `json:"admin_required"`
	RequestExample  interface{} `json:"request_example,omitempty"`
	ResponseExample interface{} `json:"response_example,omitempty"`
}

type HelpResponse struct {
	Service   string         `json:"serivce"`
	Version   string         `json:"version"`
	Endpoints []EndpointInfo `json:"endpoints"`
}

func Help() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("HELP HANDLER CALLED!")
		endpoints := []EndpointInfo{
			{
				Method:        "POST",
				Path:          "/register",
				Description:   "Регистрация нового пользователя",
				AuthRequired:  false,
				AdminRequired: false,
				RequestExample: map[string]interface{}{
					"username":   "john_doe",
					"password":   "secret123",
					"admin_pass": "optional_admin_password",
				},
				ResponseExample: map[string]interface{}{
					"status": "user registered successfully",
					"tokens": map[string]string{
						"access_token":  "eyJhbGci...",
						"refresh_token": "eyJhbGci...",
					},
				},
			},
			{
				Method:        "POST",
				Path:          "/login",
				Description:   "Вход в систему",
				AuthRequired:  false,
				AdminRequired: false,
				RequestExample: map[string]interface{}{
					"username": "john_doe",
					"password": "secret123",
				},
				ResponseExample: map[string]interface{}{
					"status": "login successful",
					"tokens": map[string]string{
						"access_token":  "eyJhbGci...",
						"refresh_token": "eyJhbGci...",
					},
				},
			},
			{
				Method:        "POST",
				Path:          "/refresh",
				Description:   "Обновление access токена",
				AuthRequired:  false,
				AdminRequired: false,
				RequestExample: map[string]interface{}{
					"refresh_token": "eyJhbGci...",
				},
				ResponseExample: map[string]interface{}{
					"status": "tokens refreshed",
					"tokens": map[string]string{
						"access_token":  "new_eyJhbGci...",
						"refresh_token": "new_eyJhbGci...",
					},
				},
			},

			{
				Method:          "GET",
				Path:            "/{alias}",
				Description:     "Редирект по короткой ссылке",
				AuthRequired:    false,
				AdminRequired:   false,
				RequestExample:  "GET /abc123 → 301 Redirect to original URL",
				ResponseExample: "HTTP 301 Redirect",
			},

			{
				Method:        "POST",
				Path:          "/url",
				Description:   "Создание короткой ссылки",
				AuthRequired:  true,
				AdminRequired: false,
				RequestExample: map[string]interface{}{
					"url":   "https://example.com/very-long-url",
					"alias": "my-short-link",
				},
				ResponseExample: map[string]interface{}{
					"id":      1,
					"alias":   "my-short-link",
					"url":     "https://example.com/very-long-url",
					"user_id": 123,
				},
			},
			{
				Method:          "GET",
				Path:            "/url",
				Description:     "Получить все мои ссылки",
				AuthRequired:    true,
				AdminRequired:   false,
				ResponseExample: []interface{}{},
			},
			{
				Method:          "GET",
				Path:            "/url/{alias}",
				Description:     "Информация о конкретной ссылке",
				AuthRequired:    true,
				AdminRequired:   false,
				ResponseExample: map[string]interface{}{},
			},
			{
				Method:        "DELETE",
				Path:          "/{alias}",
				Description:   "Удалить ссылку",
				AuthRequired:  true,
				AdminRequired: false,
				ResponseExample: map[string]interface{}{
					"status": "URL deleted successfully",
				},
			},
			{
				Method:        "PATCH",
				Path:          "/{alias}",
				Description:   "Частичное обновление ссылки (только alias ИЛИ только URL)",
				AuthRequired:  true,
				AdminRequired: false,
				RequestExample: map[string]interface{}{
					"alias": "new-alias",
					"url":   "https://new-url.com",
				},
				ResponseExample: map[string]interface{}{},
			},
			{
				Method:        "PUT",
				Path:          "/{alias}",
				Description:   "Полное обновление ссылки (alias И URL)",
				AuthRequired:  true,
				AdminRequired: false,
				RequestExample: map[string]interface{}{
					"alias": "new-alias",
					"url":   "https://new-url.com",
				},
				ResponseExample: map[string]interface{}{},
			},

			{
				Method:          "GET",
				Path:            "/users",
				Description:     "Получить список всех пользователей (только для админов)",
				AuthRequired:    true,
				AdminRequired:   true,
				ResponseExample: []interface{}{},
			},
			{
				Method:          "DELETE",
				Path:            "/users/{username}",
				Description:     "Удалить пользователя (только для админов)",
				AuthRequired:    true,
				AdminRequired:   true,
				ResponseExample: map[string]interface{}{},
			},
			{
				Method:          "GET",
				Path:            "/admin/url",
				Description:     "Получить все ссылки всех пользователей (только для админов)",
				AuthRequired:    true,
				AdminRequired:   true,
				ResponseExample: []interface{}{},
			},
			{
				Method:          "GET",
				Path:            "/{alias}/{user_id}",
				Description:     "Редирект для админа по ссылке любого пользователя",
				AuthRequired:    true,
				AdminRequired:   true,
				ResponseExample: "HTTP 301 Redirect",
			},
			{
				Method:          "DELETE",
				Path:            "/{alias}/{user_id}",
				Description:     "Удалить ссылку любого пользователя (только для админов)",
				AuthRequired:    true,
				AdminRequired:   true,
				ResponseExample: map[string]interface{}{},
			},
		}

		response := HelpResponse{
			Service:   "URL Shortener API",
			Version:   "1.0.0",
			Endpoints: endpoints,
		}

		render.JSON(w, r, response)
	}
}
