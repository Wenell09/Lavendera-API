# Lavendera API

Backend API **multi-tenant** untuk sistem manajemen laundry (SaaS) berbasis Go dan PostgreSQL. Dilengkapi dengan isolasi tenant (*tenant isolation*), manajemen role berbasis RBAC, sistem order lengkap dengan tracking publik, dan manajemen pembayaran.

---

## 🛠 Tech Stack

| Kategori | Library / Tools |
|---|---|
| **Language** | Go 1.27+ |
| **Framework** | Gin Gonic `v1.12.0` |
| **ORM & Database** | GORM `v1.31.2` + PostgreSQL `v1.6.2` |
| **Dependency Injection** | Google Wire `v0.7.0` |
| **Auth & Security** | JWT HS256, bcrypt |
| **Configuration** | Viper |
| **Validation** | Go Playground Validator `v10` |
| **Migration** | golang-migrate `v4.19.1` |
| **Logging** | Logrus |
| **Containerization** | Docker + Docker Compose |

---

## 🏗 Arsitektur & Struktur Proyek

Menggunakan **Layered Clean Architecture / Modular Monolith**. Setiap domain/fitur berdiri mandiri di direktori `internal/<feature>/`:

```
lavendera-API/
├── main.go
├── docker-compose.yaml
├── dockerfile
├── .env.example
├── docs/
│   └── API.md                  # Dokumentasi API lengkap
└── internal/
    ├── app/
    │   └── routes.go            # Registrasi seluruh route
    ├── auth/                    # Registrasi & Login
    ├── outlet/                  # Manajemen Outlet
    ├── outlet_payment_method/   # Payment Method per Outlet
    ├── outlet_user/             # Penugasan User ke Outlet
    ├── service_category/        # Kategori Layanan Laundry
    ├── service/                 # Layanan Laundry
    ├── discount/                # Diskon & Promosi
    ├── user/                    # Manajemen User / Staff
    ├── customer/                # Manajemen Pelanggan
    ├── tenant/                  # Profil Tenant
    ├── order/                   # Sistem Order & Pembayaran
    ├── models/                  # Entity Model terpusat (GORM)
    └── shared/
        ├── appcontext/          # Context helper (tenant_id, user_id)
        ├── apperror/            # Error types & handler
        ├── applogger/           # Logrus logger setup
        ├── appvalidator/        # Custom validator (id_phone, dll)
        ├── config/              # Viper config & JWT config
        ├── database/            # GORM setup & TenantScope
        ├── middleware/          # JWT, RBAC, Logger middleware
        ├── response/            # Standar response format
        └── utils/               # JWT utils, slug, phone, dll
```

Setiap modul memiliki struktur internal:
```
<module>/
├── controller/   # HTTP handler & validasi input
├── service/      # Business logic & aturan bisnis
├── repository/   # Data Access Layer via GORM
├── dto/          # Request, Response, & Filter struct
├── routes/       # Registrasi route
└── wire.go       # Dependency injection (Wire)
```

---

## 🔐 Keamanan & Tenant Isolation

- **JWT HS256**: Token menyimpan `user_id`, `tenant_id`, dan `role`. Diverifikasi oleh `JWTMiddleware`.
- **RBAC**: `RequireRole` middleware memastikan hanya role yang sesuai yang dapat mengakses endpoint tertentu (`ADMIN` / `STAFF`).
- **Tenant Isolation**: Setiap query ke database secara otomatis disaring berdasarkan `tenant_id` yang diambil dari JWT menggunakan GORM `TenantScope`. Data antar tenant tidak bisa bocor satu sama lain.

---

## 🚀 Fitur Utama

| Fitur | Deskripsi |
|---|---|
| **Auth & Onboarding** | Registrasi tenant + admin pertama, login JWT |
| **Outlet Management** | CRUD outlet per tenant dengan auto-slug, pagination & search |
| **Service Category** | CRUD kategori layanan laundry |
| **Service Management** | CRUD layanan laundry, filter kategori, unit: `kg`/`pcs`/`pair` |
| **Discount Management** | CRUD diskon tipe `percentage` & `fixed` dengan batas auto-cap 100% |
| **User / Staff Management** | CRUD akun user/staf per tenant, filter role, proteksi hapus diri sendiri |
| **Customer Management** | CRUD pelanggan per tenant, paginasi & search |
| **Tenant Management** | Lihat & update profil tenant |
| **Outlet Payment Methods** | Manajemen metode pembayaran spesifik outlet (`qris` / `bank_transfer`) |
| **Outlet User Assignment** | Penugasan staf ke outlet tertentu |
| **Order Management** | Order kasir & order publik mandiri, kalkulasi harga & diskon otomatis |
| **Payment Management** | Konfirmasi pembayaran, input bayar tunai & hitung kembalian |
| **Order Tracking** | Tracking order publik via token unik tanpa login |
| **Dashboard & Reports** | Ringkasan statistik operasional, rekap status, metode pembayaran, grafik penjualan & laporan keuangan terpaginasi |

---

## 📡 Ringkasan Endpoint API (`/api/v1`)

> Dokumentasi lengkap dengan contoh request/response ada di [`docs/API.md`](docs/API.md).

### Public (Tanpa Auth)
| Method | Path | Deskripsi |
|---|---|---|
| `GET` | `/` | Health Check |
| `POST` | `/api/v1/auth/register` | Registrasi Tenant & Admin |
| `POST` | `/api/v1/auth/login` | Login & Dapatkan JWT |
| `POST` | `/api/v1/public/orders` | Buat Order Mandiri (Customer) |
| `GET` | `/api/v1/public/orders/track/:token` | Lacak Order via Token |

### Outlet
| Method | Path | Role | Deskripsi |
|---|---|---|---|
| `POST` | `/api/v1/outlets` | `ADMIN` | Buat Outlet |
| `GET` | `/api/v1/outlets` | `ADMIN` | List Outlet |
| `GET` | `/api/v1/outlets/:id` | `ADMIN` | Detail Outlet |
| `PATCH` | `/api/v1/outlets/:id` | `ADMIN` | Update Outlet |
| `DELETE` | `/api/v1/outlets/:id` | `ADMIN` | Hapus Outlet |

### Service Category
| Method | Path | Role | Deskripsi |
|---|---|---|---|
| `GET` | `/api/v1/service-categories` | `ADMIN, STAFF` | List Kategori |
| `GET` | `/api/v1/service-categories/:id` | `ADMIN, STAFF` | Detail Kategori |
| `POST` | `/api/v1/service-categories` | `ADMIN` | Buat Kategori |
| `PATCH` | `/api/v1/service-categories/:id` | `ADMIN` | Update Kategori |
| `DELETE` | `/api/v1/service-categories/:id` | `ADMIN` | Hapus Kategori |

### Service
| Method | Path | Role | Deskripsi |
|---|---|---|---|
| `GET` | `/api/v1/services` | `ADMIN, STAFF` | List Layanan |
| `GET` | `/api/v1/services/:id` | `ADMIN, STAFF` | Detail Layanan |
| `POST` | `/api/v1/services` | `ADMIN` | Buat Layanan |
| `PATCH` | `/api/v1/services/:id` | `ADMIN` | Update Layanan |
| `DELETE` | `/api/v1/services/:id` | `ADMIN` | Hapus Layanan |

### Discount
| Method | Path | Role | Deskripsi |
|---|---|---|---|
| `GET` | `/api/v1/discounts` | `ADMIN, STAFF` | List Diskon |
| `GET` | `/api/v1/discounts/:id` | `ADMIN, STAFF` | Detail Diskon |
| `POST` | `/api/v1/discounts` | `ADMIN` | Buat Diskon |
| `PATCH` | `/api/v1/discounts/:id` | `ADMIN` | Update Diskon |
| `DELETE` | `/api/v1/discounts/:id` | `ADMIN` | Hapus Diskon |

### User / Staff
| Method | Path | Role | Deskripsi |
|---|---|---|---|
| `GET` | `/api/v1/users` | `ADMIN` | List User |
| `GET` | `/api/v1/users/:id` | `ADMIN` | Detail User |
| `POST` | `/api/v1/users` | `ADMIN` | Buat User/Staff |
| `PATCH` | `/api/v1/users/:id` | `ADMIN` | Update User |
| `DELETE` | `/api/v1/users/:id` | `ADMIN` | Hapus User |

### Customer
| Method | Path | Role | Deskripsi |
|---|---|---|---|
| `GET` | `/api/v1/customers` | `ADMIN, STAFF` | List Customer |
| `GET` | `/api/v1/customers/:id` | `ADMIN, STAFF` | Detail Customer |
| `POST` | `/api/v1/customers` | `ADMIN, STAFF` | Buat Customer |
| `PATCH` | `/api/v1/customers/:id` | `ADMIN, STAFF` | Update Customer |
| `DELETE` | `/api/v1/customers/:id` | `ADMIN, STAFF` | Hapus Customer |

### Tenant
| Method | Path | Role | Deskripsi |
|---|---|---|---|
| `GET` | `/api/v1/tenants/:id` | `ADMIN` | Detail Tenant |
| `PATCH` | `/api/v1/tenants/:id` | `ADMIN` | Update Tenant |

### Outlet Payment Methods
| Method | Path | Role | Deskripsi |
|---|---|---|---|
| `POST` | `/api/v1/outlets/:id/payment-methods` | `ADMIN` | Tambah Payment Method |
| `GET` | `/api/v1/outlets/:id/payment-methods` | `ADMIN, STAFF` | List Payment Methods |
| `GET` | `/api/v1/outlets/:id/payment-methods/:pm_id` | `ADMIN, STAFF` | Detail Payment Method |
| `PATCH` | `/api/v1/outlets/:id/payment-methods/:pm_id` | `ADMIN` | Update Payment Method |
| `DELETE` | `/api/v1/outlets/:id/payment-methods/:pm_id` | `ADMIN` | Hapus Payment Method |

### Outlet Users
| Method | Path | Role | Deskripsi |
|---|---|---|---|
| `GET` | `/api/v1/outlets/:id/users` | `ADMIN` | List User di Outlet |
| `GET` | `/api/v1/users/:id/outlets` | `ADMIN, STAFF` | List Outlet yang Di-assign ke User |
| `POST` | `/api/v1/outlets/:id/users` | `ADMIN` | Assign User ke Outlet |
| `DELETE` | `/api/v1/outlets/:id/users/:user_id` | `ADMIN` | Unassign User dari Outlet |

### Order
| Method | Path | Role | Deskripsi |
|---|---|---|---|
| `POST` | `/api/v1/orders` | `ADMIN, STAFF` | Buat Order (Kasir) |
| `GET` | `/api/v1/orders` | `ADMIN, STAFF` | List Order |
| `GET` | `/api/v1/orders/:id` | `ADMIN, STAFF` | Detail Order |
| `PATCH` | `/api/v1/orders/:id/status` | `ADMIN, STAFF` | Update Status Order |
| `PATCH` | `/api/v1/orders/:id/cancel` | `ADMIN, STAFF` | Batalkan Order |
| `PATCH` | `/api/v1/orders/:id/payments/:payment_id/confirm` | `ADMIN, STAFF` | Konfirmasi Pembayaran |
| `PATCH` | `/api/v1/orders/:id/cash-payment` | `ADMIN, STAFF` | Input Bayar Tunai |

### Dashboard & Reports
| Method | Path | Role | Deskripsi |
|---|---|---|---|
| `GET` | `/api/v1/dashboard` | `ADMIN, STAFF` | Ringkasan Statistik Dashboard |
| `GET` | `/api/v1/reports/financial` | `ADMIN, STAFF` | Laporan Keuangan & Riwayat Pembayaran |

---

## ⚙️ Cara Menjalankan

### Prasyarat
- Go 1.27+
- PostgreSQL (atau via Docker)
- Wire CLI: `go install github.com/google/wire/cmd/wire@latest`

### 1. Clone & Setup Environment
```bash
git clone <repo-url>
cd lavendera-API
cp .env.example .env
```

Edit `.env` dan isi variabel berikut:
```env
APP_ENV=development
PORT=8000
DATABASE_URL=postgresql://user:password@localhost:5432/lavendera?sslmode=disable
JWT_SECRET=your-secret-key-min-32-chars
```

### 2. Jalankan dengan Docker Compose (Direkomendasikan)
```bash
docker compose up --build
```
Ini akan menjalankan aplikasi dan PostgreSQL secara bersamaan.

### 3. Jalankan Secara Manual

**Generate Wire dependency:**
```bash
wire ./internal/app
```

**Jalankan database migration:**
```bash
migrate -path ./migrations -database "$DATABASE_URL" up
```

**Jalankan server:**
```bash
go run main.go
```

Server akan berjalan di `http://localhost:8000`.

---

## 📖 Dokumentasi API Lengkap

Lihat [`docs/API.md`](docs/API.md) untuk dokumentasi lengkap semua endpoint beserta contoh request & response JSON.
