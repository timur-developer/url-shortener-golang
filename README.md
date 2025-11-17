# URL Shortener

A secure and scalable URL shortening service with authentication and role-based access control.

## Features

- URL shortening with custom aliases
- User authentication (Basic Auth)
- Role-based permissions (Anonymous/User/Admin)
- RESTful API
- PostgreSQL storage
- Admin dashboard

## Tech Stack

- Go 1.21+
- PostgreSQL
- Chi Router
- Basic Authentication

---

## Quick Start

### 1. Clone the repository

```bash
git clone <your-repo>
cd go-url-shortener
```

### 2. Configure the database

Set up PostgreSQL connection in `config/local.yaml`

### 3. Run the application

```bash
go run cmd/url-shortener/main.go
```

---

## API Examples

### Create short URL

```http
POST /url
Content-Type: application/json

{
  "url": "https://example.com",
  "alias": "my-link"
}
```

### Register user

```http
POST /register
Content-Type: application/json

{
  "username": "user",
  "password": "pass"
}
```

### Redirect to original URL

```http
GET /my-link
```

---

## Role-Based Access

| Role      | Permissions                               |
|-----------|-------------------------------------------|
| Anonymous | Create links, Redirect                    |
| Users     | Manage own links, View analytics          |
| Admins    | Manage all links, User management         |

---

## Git Setup

If you're setting up this repository for the first time:

```bash
git init
git add .
git commit -m "Initial commit: URL shortener with auth"
git remote add origin https://github.com/YOUR_USERNAME/YOUR_REPO_NAME.git
git branch -M main
git push -u origin main
```
