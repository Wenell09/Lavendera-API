# Lavendera API

Backend API multi-tenant untuk sistem manajemen laundry (SaaS) berbasis Go dan PostgreSQL. Sistem dilengkapi dengan isolasi tenant (*tenant isolation*), manajemen role (*RBAC*), dan RESTful API.

---

## 🛠 Tech Stack

- **Language**: Go 1.27+
- **Framework**: Gin Gonic (`v1.12.0`)
- **ORM & DB**: GORM (`v1.31.2`) + PostgreSQL (`v1.6.2`)
- **Dependency Injection**: Google Wire (`v0.7.0`)
- **Auth & Security**: JWT HS256, bcrypt
- **Configuration & Validation**: Viper, Go Playground Validator
- **Migration**: golang-migrate (`v4.19.1`)

---

## 🏗 Arsitektur & Struktur Proyek

Menggunakan **Layered Clean Architecture / Modular Monolith**:
- Setiap domain/fitur berada di `internal/<feature>/` dengan struktur:
  - `dto/`: Request, response, dan filter struct
  - `controller/`: HTTP Request Handler & validation
  - `service/`: Business logic & tenant validation
  - `repository/`: Data Access Layer via GORM + `TenantScope`
- Entity models terpusat di `internal/models/`
- Shared utilities & middleware di `internal/shared/`

---

## 🚀 Fitur Utama

- **Authentication & Onboarding**: Registrasi tenant + admin default, login JWT.
- **Outlet Management**: CRUD outlet per tenant, auto-slug, pagination, & search.
- **Service Category**: CRUD kategori layanan laundry.
- **Service Management**: CRUD layanan laundry, filter kategori, pagination, & search.
- **Discount Management**: CRUD diskon promosi (persentase / nominal).
- **User / Staff Management**: CRUD akun user/staf per tenant, pagination, filter role, & proteksi hapus diri sendiri.

---

## 📡 Endpoint API Utama (`/api/v1`)

| Method | Path | Auth | Role | Deskripsi |
|---|---|---|---|---|
| `GET` | `/` | Public | - | Health Check |
| `POST` | `/api/v1/auth/register` | Public | - | Registrasi Tenant & Admin |
| `POST` | `/api/v1/auth/login` | Public | - | Login & Obtain JWT |
| `POST` | `/api/v1/outlets` | Bearer | `ADMIN` | Buat Outlet Baru |
| `GET` | `/api/v1/outlets` | Bearer | `ADMIN` | List Outlet (Paginated) |
| `GET` | `/api/v1/outlets/:id` | Bearer | `ADMIN` | Detail Outlet |
| `PATCH` | `/api/v1/outlets/:id` | Bearer | `ADMIN` | Update Outlet |
| `DELETE` | `/api/v1/outlets/:id` | Bearer | `ADMIN` | Hapus Outlet |
| `GET` | `/api/v1/service-categories` | Bearer | `ADMIN` | List Kategori Layanan |
| `POST` | `/api/v1/service-categories` | Bearer | `ADMIN` | Buat Kategori Layanan |
| `PATCH` | `/api/v1/service-categories/:id` | Bearer | `ADMIN` | Update Kategori Layanan |
| `DELETE` | `/api/v1/service-categories/:id` | Bearer | `ADMIN` | Hapus Kategori Layanan |
| `GET` | `/api/v1/services` | Bearer | `ADMIN` | List Layanan Laundry |
| `POST` | `/api/v1/services` | Bearer | `ADMIN` | Buat Layanan Laundry |
| `PATCH` | `/api/v1/services/:id` | Bearer | `ADMIN` | Update Layanan Laundry |
| `DELETE` | `/api/v1/services/:id` | Bearer | `ADMIN` | Hapus Layanan Laundry |
| `GET` | `/api/v1/discounts` | Bearer | `ADMIN` | List Diskon Promosi |
| `POST` | `/api/v1/discounts` | Bearer | `ADMIN` | Buat Diskon Promosi |
| `PATCH` | `/api/v1/discounts/:id` | Bearer | `ADMIN` | Update Diskon Promosi |
| `DELETE` | `/api/v1/discounts/:id` | Bearer | `ADMIN` | Hapus Diskon Promosi |
| `GET` | `/api/v1/users` | Bearer | `ADMIN` | List User / Staff (Paginated + Filter) |
| `POST` | `/api/v1/users` | Bearer | `ADMIN` | Buat User / Staff Baru |
| `GET` | `/api/v1/users/:id` | Bearer | `ADMIN` | Detail User / Staff |
| `PATCH` | `/api/v1/users/:id` | Bearer | `ADMIN` | Update User / Staff |
| `DELETE` | `/api/v1/users/:id` | Bearer | `ADMIN` | Hapus User / Staff |

---

## ⚙️ Cara Menjalankan

### 1. Requirements
- Go 1.27+
- PostgreSQL
- Wire (`go install github.com/google/wire/cmd/wire@latest`)

### 2. Environment Setup
Salin `.env.example` ke `.env` dan konfigurasikan variabel environment:
```bash
cp .env.example .env
```

### 3. Generate Wire Dependency
```bash
wire ./internal/app
```

### 4. Run Application
```bash
go run main.go
```
