# Go WebSocket Template with Prisma Go Client

A modern Go template featuring WebSocket support and Prisma Go client for efficient database operations.

## 🚀 Features

- **WebSocket Integration**: Real-time bidirectional communication
- **Prisma Go Client**: Type-safe database operations
- **RESTful API**: Simple and clean HTTP endpoints
- **Database Integration**: Efficient database operations using Prisma Go client
- **Modern Architecture**: Follows Go best practices

## 📋 Prerequisites

- Go 1.16 or higher
- Prisma Go client
- Database (PostgreSQL recommended)

## 🛠️ Installation

1. Clone the repository
```bash
git clone https://github.com/yourusername/project-name
cd project-name
```

2. Install dependencies
```bash
go mod download
```

3. Set up your environment variables
```bash
cp .env.example .env
# Edit .env with your database credentials
```

4. Generate Prisma Go client
```bash
go run github.com/steebchen/prisma-client-go generate
```

5. Start the server
```bash
go run main.go
```

## 🏗️ Project Structure

```
.
├── main.go
├── client.go
├── manager.go
├── events.go
├── config/
│   └── database.go
├── controller/
│   └── user_controller.go
├── data/
│   └── request/
│       ├── user_create_req.go
│       └── user_update_req.go
│   └── response/
│       ├── user_response.go
│       └── web_response.go
├── helper/
│   └── json.go
├── repository/
│   ├── post_repo.go
│   └── user_repo_impl.go
├── service/
│   ├── user_service_impl.go
│   └── user_service.go
├── model/
│   ├── file.go
│   ├── user.go
│   └── room.go
├── prisma/
│   ├── schema.prisma
│   └── migrations/
├── pkg/
│   ├── database/
│   └── websocket/
├── .env.example
├── go.mod
└── README.md
```

## 🔌 WebSocket Usage

Connect to the WebSocket endpoint:
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');
```

## 🗄️ Database Operations

Example of using Prisma Go client:
```go
// Create a new user
user, err := client.User.CreateOne(
    db.User.Email.Set("example@email.com"),
    db.User.Name.Set("John Doe"),
).Exec(ctx)

// Find a user
user, err := client.User.FindFirst(
    db.User.Email.Equals("example@email.com"),
).Exec(ctx)

// Update a user
user, err := client.User.FindUnique(
    db.User.ID.Equals("user-id"),
).Update(
    db.User.Name.Set("Jane Doe"),
).Exec(ctx)
```

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | /api/users | Get all users |
| POST   | /api/users | Create new user |
| GET    | /api/users/:id | Get user by ID |
| PUT    | /api/users/:id | Update user |
| DELETE | /api/users/:id | Delete user |
| WS     | /ws | WebSocket connection |

## 🛡️ Environment Variables

Required environment variables:
```env
DATABASE_URL="postgresql://username:password@localhost:5432/dbname"
PORT=8080
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
