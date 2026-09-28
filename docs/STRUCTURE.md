# Directory Structure Documentation — Lavendera API

> **Status**: Verifikasi dari implementasi aktual kode.  
> **Source of Truth**: Seluruh susunan direktori dan file pada repositori saat ini.

---

## 1. Complete Directory Tree

```
/home/wenell/projects/lavendera-API/
├── .env                              # File konfigurasi lokal (diabaikan git)
├── .env.example                      # Contoh template environment variable
├── .gitignore                        # Konfigurasi pengabaian git
├── go.mod                            # Deklarasi module Go dan daftar dependency
├── go.sum                            # Checksum dependency Go
├── main.go                           # Entry point server HTTP
├── internal/
│   ├── app/                          # Wiring dependency, routing induk, lifecycle aplikasi
│   │   ├── app.go                    # Struct App (Router, DB, Logger)
│   │   ├── providers.go              # Database provider untuk Wire
│   │   ├── routes.go                 # Registrasi router Gin, grouping middleware & controller
│   │   ├── wire.go                   # Definisi Wire Build
│   │   └── wire_gen.go               # Hasil generate Google Wire
│   │
│   ├── auth/                         # Modul Autentikasi & Registrasi Tenant/Admin
│   │   ├── controller/               # HTTP Handlers
│   │   │   ├── auth_controller.go
│   │   │   └── auth_controller_impl.go
│   │   ├── dto/                      # Data Transfer Objects
│   │   │   ├── auth_request.go
│   │   │   └── auth_response.go
│   │   ├── repository/               # Akses Data database
│   │   │   ├── auth_repository.go
│   │   │   └── auth_repository_impl.go
│   │   ├── routes/                   # Routing auth
│   │   │   └── auth_routes.go
│   │   ├── service/                  # Business Logic
│   │   │   ├── auth_service.go
│   │   │   └── auth_service_impl.go
│   │   └── wire.go                   # WireSet untuk auth
│   │
│   ├── outlet/                       # Modul Manajemen Outlet Laundry
│   │   ├── controller/
│   │   │   ├── outlet_controller.go
│   │   │   └── outlet_controller_impl.go
│   │   ├── dto/
│   │   │   ├── outlet_filter.go
│   │   │   ├── outlet_request.go
│   │   │   └── outlet_response.go
│   │   ├── repository/
│   │   │   ├── outlet_repository.go
│   │   │   └── outlet_repository_impl.go
│   │   ├── routes/
│   │   │   └── outlet_routes.go
│   │   ├── service/
│   │   │   ├── outlet_service.go
│   │   │   └── outlet_service_impl.go
│   │   └── wire.go
│   │
│   ├── service_category/             # Modul Kategori Layanan (misal: Kiloan, Satuan)
│   │   ├── controller/
│   │   │   ├── service_category_controller.go
│   │   │   └── service_category_controller_impl.go
│   │   ├── dto/
│   │   │   ├── service_category_request.go
│   │   │   └── service_category_response.go
│   │   ├── repository/
│   │   │   ├── service_category_repository.go
│   │   │   └── service_category_repository_impl.go
│   │   ├── routes/
│   │   │   └── service_category_routes.go
│   │   ├── service/
│   │   │   ├── service_category_service.go
│   │   │   └── service_category_service_impl.go
│   │   └── wire.go
│   │
│   ├── service/                      # Modul Layanan Laundry
│   │   ├── controller/
│   │   │   ├── service_controller.go
│   │   │   └── service_controller_impl.go
│   │   ├── dto/
│   │   │   ├── service_filter.go
│   │   │   ├── service_request.go
│   │   │   └── service_response.go
│   │   ├── repository/
│   │   │   ├── service_repository.go
│   │   │   └── service_repository_impl.go
│   │   ├── routes/
│   │   │   └── service_routes.go
│   │   ├── service/
│   │   │   ├── service_service.go
│   │   │   └── service_service_impl.go
│   │   └── wire.go
│   │
│   ├── discount/                     # Modul Diskon Promosi
│   │   ├── controller/
│   │   │   ├── discount_controller.go
│   │   │   └── discount_controller_impl.go
│   │   ├── dto/
│   │   │   ├── discount_filter.go
│   │   │   ├── discount_request.go
│   │   │   └── discount_response.go
│   │   ├── repository/
│   │   │   ├── discount_repository.go
│   │   │   └── discount_repository_impl.go
│   │   ├── routes/
│   │   │   └── discount_routes.go
│   │   ├── service/
│   │   │   ├── discount_service.go
│   │   │   └── discount_service_impl.go
│   │   └── wire.go
│   │
│   ├── models/                       # GORM Database Models (Entity terpusat)
│   │   ├── audit_logs.go
│   │   ├── customer.go
│   │   ├── discount.go
│   │   ├── notification.go
│   │   ├── order.go
│   │   ├── order_item.go
│   │   ├── outlet.go
│   │   ├── outlet_payment_method.go
│   │   ├── outlet_user.go
│   │   ├── payment.go
│   │   ├── service.go
│   │   ├── service_category.go
│   │   ├── tenant.go
│   │   └── user.go
│   │
│   ├── db/
│   │   └── migrations/               # Script Migrasi SQL PostgreSQL (Up & Down)
│   │       ├── 20260908052809_create_extension.{up,down}.sql
│   │       ├── 20260908053004_create_tenants.{up,down}.sql
│   │       ├── 20260908053325_create_customers.{up,down}.sql
│   │       ├── 20260908061134_create_users.{up,down}.sql
│   │       ├── 20260908061215_create_audit_logs.{up,down}.sql
│   │       ├── 20260908061245_create_outlets.{up,down}.sql
│   │       ├── 20260908061502_create_service_categories.{up,down}.sql
│   │       ├── 20260908061536_create_discounts.{up,down}.sql
│   │       ├── 20260908061613_create_outlet_users.{up,down}.sql
│   │       ├── 20260908061644_create_services.{up,down}.sql
│   │       ├── 20260908061713_create_outlet_payment_methods.{up,down}.sql
│   │       ├── 20260908062453_create_orders.{up,down}.sql
│   │       ├── 20260908062517_create_order_items.{up,down}.sql
│   │       ├── 20260908062851_create_payments.{up,down}.sql
│   │       ├── 20260908062933_create_notifications.{up,down}.sql
│   │       ├── 20260908063056_create_index.{up,down}.sql
│   │       └── 20260908065344_create_updated_at_triggers.{up,down}.sql
│   │
│   └── shared/                       # Komponen Shared (Utilitas, Middleware, Konfigurasi)
│       ├── appcontext/
│       │   └── context.go            # UserIDKey, TenantIDKey, roleKey getter/setter
│       ├── apperror/
│       │   ├── app_error.go          # Interface AppError
│       │   ├── auth_error.go         # UnauthorizedError, ConflictError
│       │   ├── error_handler.go      # NewHandleError
│       │   ├── not_found_error.go    # NotFoundError
│       │   └── validation_error.go   # ValidationError, FieldError, NewFieldError
│       ├── applogger/
│       │   └── logger.go             # NewLogger (Logrus instance)
│       ├── appvalidator/
│       │   └── validator.go          # NewValidator + custom "id_phone" registration
│       ├── config/
│       │   ├── config.go             # Struct Config & LoadConfig() via Viper
│       │   └── jwt.go                # Struct JWTConfig & NewJWTConfig()
│       ├── database/
│       │   ├── database.go           # DBConnection() (GORM + connection pool)
│       │   └── scope.go              # TenantScope(ctx)
│       ├── middleware/
│       │   ├── jwt_middleware.go     # JWTMiddleware(jwtConfig)
│       │   ├── logger_middleware.go  # LoggerMiddleware(logger)
│       │   └── rbac_middleware.go    # RequireRole(roles...)
│       ├── response/
│       │   └── response.go           # ResponseSuccess, ResponseError, ResponseMeta
│       └── utils/
│           ├── jwt.go                # GenerateToken()
│           ├── log_with_ctx.go       # LogWithContext()
│           ├── phone_utils.go        # IsValidIndonesianPhone()
│           └── slug.go               # GenerateSlug()
```

---

## 2. Directory Responsibilities & Boundaries

### 2.1. Root Files
- **[main.go](file:///home/wenell/projects/lavendera-API/main.go)**: Titik masuk utama aplikasi. Membaca konfigurasi, menjalankan dependency injection via Wire, menyiapkan penutupan pool database via `defer sqlDB.Close()`, dan menjalankan listener HTTP server.
- **[go.mod](file:///home/wenell/projects/lavendera-API/go.mod)**: Mendefinisikan nama modul dan pin dependency.

### 2.2. `internal/app/`
- **Tanggung Jawab**: Integrasi dan wiring tingkat aplikasi. Menyatukan semua controller, database, logger, middleware, dan validator menjadi instance `*App` yang siap dijalankan.
- **Yang Boleh Berada di Sini**: Definisi Wire (`wire.go`), provider database/router tingkat aplikasi, registrasi grup endpoint induk (`routes.go`), struct `App`.
- **Yang Tidak Boleh**: Implementasi logika bisnis atau query SQL langsung.

### 2.3. `internal/<feature>/` (Modul Domain)
Setiap folder fitur mengisolasi domainnya masing-masing.
- **`controller/`**:
  - Tanggung Jawab: Menerima request HTTP Gin, parsing JSON / URL query / URL param, memanggil struct validator, memanggil service layer, mengembalikan response JSON seragam.
  - Aturan: Tidak boleh memanggil database atau repository secara langsung.
- **`service/`**:
  - Tanggung Jawab: Menjalankan business logic, validasi keunikan nama, mengekstrak tenant ID dari context, memetakan DTO ke entity model.
  - Aturan: Tidak boleh mengimpor package Gin (`*gin.Context`). Hanya menggunakan `context.Context`.
- **`repository/`**:
  - Tanggung Jawab: Mengeksekusi query database GORM dengan `TenantScope(ctx)`.
  - Aturan: Tidak boleh mengembalikan DTO response; selalu mengembalikan Model entity atau error.
- **`dto/`**:
  - Tanggung Jawab: Mendefinisikan kontrak data request, response, dan filter khusus fitur terkait.
  - Aturan: Struct request harus memiliki tag JSON dan validasi (`validate:"..."`).
- **`routes/`**:
  - Tanggung Jawab: Mendaftarkan rute fitur ke dalam `*gin.RouterGroup`.
- **`wire.go`**:
  - Tanggung Jawab: Mendefinisikan `WireSet` berisikan konstruktor repository, service, dan controller untuk diekspor ke `internal/app/wire.go`.

### 2.4. `internal/models/`
- **Tanggung Jawab**: Menyimpan seluruh representasi tabel database GORM dalam satu package terpusat.
- **Yang Boleh Berada di Sini**: Definisi struct entitas, tag GORM (`gorm:"..."`), relasi foreign key GORM.
- **Yang Tidak Boleh**: Logika bisnis, HTTP DTO, atau fungsi query database.

### 2.5. `internal/db/migrations/`
- **Tanggung Jawab**: Script SQL versi skema PostgreSQL untuk membuat ekstensi, tabel, indeks, foreign key, dan trigger `updated_at`. Menggunakan format penamaan timestamp `YYYYMMDDHHMMSS_<description>.{up,down}.sql`.

### 2.6. `internal/shared/`
- **Tanggung Jawab**: Komponen cross-cutting yang dibutuhkan oleh banyak modul fitur:
  - `appcontext`: Standar penyimpanan dan pembacaan user/tenant/role di context Go.
  - `apperror`: Definisi error aplikasi dan mapping HTTP status code.
  - `applogger`: Setup logger logrus standar.
  - `appvalidator`: Validator instance dengan custom rules.
  - `config`: Konfigurasi terpusat berbasis Viper.
  - `database`: GORM connection pool & scoping multi-tenant.
  - `middleware`: Gin middlewares.
  - `response`: Format response envelope standar.
  - `utils`: Helper independen (slug, phone, jwt, context log).

---

## 3. Inter-Directory Dependencies

```
[internal/app]
      │
      ├──> [internal/auth]
      ├──> [internal/outlet]
      ├──> [internal/service]
      ├──> [internal/service_category]
      ├──> [internal/discount]
      └──> [internal/shared]

[internal/<feature>]
      │
      ├──> [internal/models]
      └──> [internal/shared]
             ├── appcontext
             ├── apperror
             ├── appvalidator
             ├── database
             ├── response
             └── utils

[internal/service] (Cross-Feature Dependency)
      │
      └──> [internal/service_category/repository] (Untuk validasi CategoryID)
```
*(Catatan: Modul `service` secara eksplisit bergantung pada `service_category/repository` untuk memvalidasi keberadaan `category_id` per tenant saat update atau create).*
