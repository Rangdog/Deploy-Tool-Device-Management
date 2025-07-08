# Device Manager

## 📝 Mô tả

Device Manager là dự án quản lý thiết bị, sử dụng:

- Go + Gin
- PostgreSQL (GORM migrate)
- Redis
- JWT Authentication
- Robfig cron
- Goroutine + channel

Deploy trên **Railway**.

---

## 🚀 Tech stack

- Go v1.xx
- Gin – HTTP framework
- GORM – ORM cho PostgreSQL
- Redis – Cache, pubsub
- Robfig cron – Scheduler jobs
- JWT – Authentication
- Docker + Railway

---

## 📁 Cấu trúc thư mục

```
.
├── .github/workflows/   # CI/CD workflows
├── .vscode/             # VSCode config
├── api/                 # Router, handler APIs
├── cmd/                 # Main application entry point
├── config/              # Env, configs
├── constant/            # Constants, enums, status codes
├── internal/
│   └── domain/
│       ├── dto/         # Data Transfer Objects
│       ├── entity/      # Database models
│       ├── filter/      # Filter struct, query params
│       ├── mocks/       # Mock data for testing
│       ├── repository/  # Repository interfaces & implementations
│       └── service/     # Business logic
├── pkg/                 # Common packages (JWT, hash, utils)
├── .env.template        # Env file template
├── .gitignore
├── Dockerfile
├── go.mod
└── go.sum
```

---

## ⚙️ Cài đặt

```bash
# Clone repo
git clone https://github.com/Khalac/Tool-Device-Management
cd be

# Cài dependencies
go mod tidy

# Chạy ứng dụng
go run cmd/server/main.go
```

Hoặc chạy bằng Docker:

```bash
docker build -t device-manager .
docker run -p 8080:8080 device-manager
```

---

## 🔧 API Overview

*(Danh sách API rút gọn để tránh dài, có thể copy từ phần trên khi cần)*

---

## ✅ Database Migration

Sử dụng **GORM AutoMigrate** trong `config/db.go` hoặc `main.go`:

```go
db.AutoMigrate(&entity.Device{}, &entity.User{}, ...)
```

---

## 🔄 Cron Jobs

Được định nghĩa trong package `cron/`, sử dụng:

```go
import "github.com/robfig/cron/v3"
```

Ví dụ: cron chạy mỗi 5 phút cập nhật trạng thái thiết bị.

---

## 🧵 Concurrency

Dự án sử dụng **goroutine + channel** cho worker pool xử lý đồng thời. Các implement concurrency nằm trong `internal/service/` hoặc `pkg/`.

---

## 🔐 Environment Variables

| Key          | Mô tả              |
| ------------ | ------------------ |
| PORT         | Cổng server        |
| DB\_USER     | Database username  |
| DB\_PASSWORD | Database password  |
| DB\_NAME     | Tên database       |
| DB\_HOST     | Địa chỉ DB         |
| REDIS\_URL   | Redis URL          |
| JWT\_SECRET  | Secret key cho JWT |

Tạo file `.env` dựa trên `.env.template`.

---

## 🛠️ Testing

```bash
go test ./...
```

---

## ✨ Deploy

Deploy bằng **Railway**:

1. Kết nối GitHub repo
2. Thiết lập environment variables theo `.env.template`
3. Railway auto build & deploy container

---

## 📄 License

MIT

---

## 🤝 Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---


