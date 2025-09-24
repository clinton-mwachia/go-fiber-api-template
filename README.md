# Fiber + MongoDB API Template 🚀

A production-ready **Go Fiber API template** with **MongoDB** and **JWT Authentication**, designed for scalability, security, and maintainability.  
This template can serve as a starting point for building APIs of any size.

---

## ✨ Features

- ⚡ Fast API server using [Fiber](https://github.com/gofiber/fiber)
- 🗄️ MongoDB integration (with official driver)
- 🔐 JWT authentication
- 👤 User & Todo example models
- 🛡️ Role-based access (Admin vs Normal User)
- 📂 Modular project structure (`models`, `routes`, `middlewares`, `utils`, `config`)
- 📝 Environment-based configuration
- 🚀 Ready for containerization and cloud deployment

---

## 📂 Project Structure

```

go-fiber-api-template/
│── main.go
│── go.mod
│── go.sum
│── config/
│ └── config.go
│── models/
│ ├── user.go
│ └── todo.go
│── routes/
│ └── router.go
│── controllers/
│ ├── todo.go
│ └── user.go
│── middlewares/
│ ├── auth.go
│ └── ownership.go
│── utils/
│ ├── folderCreate.go
│ └── password.go
│── .env
│── README.md

```

---

## ⚙️ Setup Instructions

### 1. Clone Repo

```bash
git https://github.com/clinton-mwachia/go-fiber-api-template.git
cd go-fiber-api-template
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Setup Environment

Create a `.env` file in the root:

```env
PORT=8080
MONGO_URI=mongodb://localhost:27017
DB_NAME=fiber_api_db
JWT_SECRET=supersecretkey
```

### 4. Run Server

```bash
go run main.go
```

---

## 🔑 Authentication

### Register User

`POST /api/user/register`

```json
{
  "username": "john",
  "email": "john@example.com",
  "role": "user",
  "password": "password123"
}
```

### Login

`POST /api/login`

```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

Response:

```json
{
  "token": "jwt-token-string",
  "expires_at": 1758965623
}
```

Include JWT token in **Authorization Header**:

```
Authorization: Bearer <token>
```

---

## 📌 API Endpoints

### Auth

- `POST /api/register` – Register new user
- `POST /api/login` – Login and receive JWT

### Users

- `GET /api/users` – Get all users
- `GET /api/user/:id` – Get user by id

### Todos

- `POST /api/todos` – Create a todo (user only sees their todos)
- `GET /api/todos` – Get logged-in user’s todos
- `GET /api/todos/all` – Admin only → Get all todos
- `PUT /api/todo/:id` – Update own todo
- `DELETE /api/todo/:id` – Delete own todo

---

## 🛡️ Roles

- **Normal User** → Can access only their own todos
- **Admin** → Can access all todos & users

---

## 🧪 Testing

You can use **Postman** or **cURL**:

```bash
curl -X GET http://localhost:8080/api/todos \
  -H "Authorization: Bearer <your-jwt>"
```

---

## 🐳 Docker Support (Optional)

COMING SOON

---

## 🛠️ Tech Stack

- [Go](https://golang.org/)
- [Fiber](https://gofiber.io/)
- [MongoDB](https://www.mongodb.com/)
- JWT for authentication
- bcrypt for password hashing

---

## 🚀 Deployment

You can deploy to:

- Docker + Kubernetes
- Render / Railway / Fly.io
- AWS, GCP, Azure

---

## 📜 License

MIT License © 2025

---

## 🙌 Contributing

PRs are welcome! Please open an issue for discussion first.
