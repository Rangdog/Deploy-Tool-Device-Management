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

### **Auth**

| Method | Endpoint           | Description   |
| ------ | ------------------ | ------------- |
| POST   | /api/auth/login    | Đăng nhập     |
| POST   | /api/auth/register | Đăng ký user  |
| POST   | /api/auth/logout   | Đăng xuất     |
| POST   | /api/auth/refresh  | Refresh token |

### **Assets**

| Method | Endpoint                    | Description             |
| ------ | --------------------------- | ----------------------- |
| GET    | /api/assets                 | Lấy danh sách thiết bị |
| POST   | /api/assets                 | Tạo thiết bị mới       |
| GET    | /api/assets/{id}            | Lấy thiết bị theo id   |
| PUT    | /api/assets/{id}            | Cập nhật thiết bị      |
| DELETE | /api/assets/{id}            | Xoá thiết bị           |
| PATCH  | /api/assets-retired/{id}    | Nghỉ hưu thiết bị      |
| GET    | /api/assets/filter          | Lọc thiết bị           |
| GET    | /api/assets/filter-dashboard| Dashboard thiết bị     |
| GET    | /api/assets-log/{asset_id}  | Lịch sử thiết bị       |
| GET    | /api/assets/maintenance-schedules | Thiết bị chưa có lịch bảo trì |

### **Assignments**

| Method | Endpoint                | Description             |
| ------ | ----------------------- | ----------------------- |
| GET    | /api/assignments/filter | Lọc danh sách bàn giao |
| GET    | /api/assignments/{id}   | Lấy bàn giao theo id   |
| PUT    | /api/assignments/{id}   | Cập nhật bàn giao      |

### **Bills**

| Method | Endpoint                 | Description              |
| ------ | ------------------------ | ------------------------ |
| POST   | /api/bills               | Tạo bill                |
| GET    | /api/bills-un-paid/      | Lấy bill chưa thanh toán |
| GET    | /api/bills/filter        | Lọc bills               |
| GET    | /api/bills/{billNumber}  | Lấy bill theo số        |
| PATCH  | /api/bills/{billNumber}  | Cập nhật trạng thái bill |

### **Categories**

| Method | Endpoint              | Description            |
| ------ | --------------------- | ---------------------- |
| GET    | /api/categories      | Lấy danh sách category |
| POST   | /api/categories      | Tạo category           |
| DELETE | /api/categories/{id} | Xoá category           |

### **Departments & Companies**

| Method | Endpoint               | Description           |
| ------ | ---------------------- | --------------------- |
| GET    | /api/departments      | Lấy danh sách phòng ban |
| POST   | /api/departments      | Tạo phòng ban         |
| DELETE | /api/departments/{id} | Xoá phòng ban         |
| POST   | /api/company         | Tạo công ty           |
| GET    | /api/company/{id}    | Lấy công ty theo id   |

### **Locations**

| Method | Endpoint             | Description            |
| ------ | -------------------- | ---------------------- |
| GET    | /api/locations      | Lấy danh sách location |
| POST   | /api/locations      | Tạo location           |
| DELETE | /api/locations/{id} | Xoá location           |

### **Maintenance Schedules**

| Method | Endpoint                        | Description                   |
| ------ | ------------------------------- | ----------------------------- |
| GET    | /api/maintenance-schedules     | Lấy danh sách lịch bảo trì   |
| POST   | /api/maintenance-schedules     | Tạo lịch bảo trì             |
| GET    | /api/maintenance-schedules/{id} | Lấy lịch bảo trì theo id    |
| PATCH  | /api/maintenance-schedules/{id} | Cập nhật lịch bảo trì       |
| DELETE | /api/maintenance-schedules/{id} | Xoá lịch bảo trì            |

### **Notifications**

| Method | Endpoint                 | Description          |
| ------ | ------------------------ | -------------------- |
| GET    | /api/notifications      | Lấy thông báo        |
| PUT    | /api/notifications/{id} | Cập nhật thông báo   |

### **Request Transfer**

| Method | Endpoint                        | Description              |
| ------ | ------------------------------- | ------------------------ |
| POST   | /api/request-transfer          | Gửi request transfer    |
| PATCH  | /api/request-transfer/confirm/{id} | Duyệt transfer          |
| PATCH  | /api/request-transfer/deny/{id} | Từ chối transfer        |
| GET    | /api/request-transfer/filter  | Lọc request transfer    |
| GET    | /api/request-transfer/{id}    | Lấy request transfer theo id |

### **Users**

| Method | Endpoint                      | Description                  |
| ------ | ----------------------------- | ---------------------------- |
| GET    | /api/users                   | Lấy danh sách user         |
| GET    | /api/users/not-dep           | User chưa có phòng ban     |
| PATCH  | /api/users/role              | Cập nhật role user         |
| PATCH  | /api/user/password-reset     | Reset password             |
| PATCH  | /api/user/information        | Cập nhật thông tin user    |
| PATCH  | /api/user/department         | Cập nhật phòng ban user    |
| PATCH  | /api/user/manager-department/{user_id} | Cập nhật manager department |
| PATCH  | /api/user/can-export/{user_id} | Cập nhật can-export        |
| GET    | /api/user/session            | Lấy session user           |
| POST   | /api/user/forget-password    | Gửi email reset password   |
| DELETE | /api/user/{email}            | Xoá user theo email        |

### **Roles**

| Method | Endpoint    | Description         |
| ------ | ----------- | ------------------- |
| GET    | /api/roles | Lấy danh sách roles |

### **Cron Jobs**

| Method | Endpoint                             | Description                    |
| ------ | ------------------------------------- | ------------------------------ |
| GET    | /api/CheckAndSenMaintenanceNotification | Gửi notification bảo trì      |
| GET    | /api/SendEmailsForWarrantyExpiry      | Gửi email hết hạn bảo hành    |
| GET    | /api/UpdateStatusWhenFinishMaintenance | Cập nhật trạng thái sau bảo trì |

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


