# URL Shortener
A secure and scalable URL shortening service with authentication and role-based access control.

## Features
URL shortening with custom aliases

User authentication (Basic Auth)

Role-based permissions (Anonymous/User/Admin)

RESTful API

PostgreSQL storage

Admin dashboard

### Quick Start
Clone and configure:
git clone <your-repo>
cd go-url-shortener
Set up PostgreSQL connection in config/local.yaml

Run:
go run cmd/url-shortener/main.go

### API Examples
Create short URL:


POST /url
{"url": "https://example.com", "alias": "my-link"}

Register user: POST /register  
{"username": "user", "password": "pass"}
Redirect:

GET /my-link
Role Access
Anonymous: Create links, Redirect
Users: Manage own links, View analytics
Admins: Manage all links, User management

### Tech Stack
Go 1.21+
PostgreSQL
Chi Router
Basic Authentication

## Git Setup Commands
bash
git init
git add .
git commit -m "Initial commit: URL shortener with auth"
git remote add origin https://github.com/YOUR_USERNAME/YOUR_REPO_NAME.git
git branch -M main
git push -u origin main
