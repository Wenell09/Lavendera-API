# AGENTS.md — Developer & AI Agent Guide for Lavendera API

> **Notice**: File ini dibuat khusus untuk memandu AI coding agent dan software engineer dalam memahami, menavigasi, dan mengembangkan fitur pada project **Lavendera API** secara konsisten dengan kode yang benar-benar ada saat ini.
>
> Seluruh informasi di dokumen ini diverifikasi langsung dari source code per **28 September 2026**.

---

## 1. Project Overview

**Lavendera API** adalah backend API multi-tenant untuk sistem manajemen laundry (SaaS). Sistem dirancang dengan isolasi tenant (*tenant isolation*), manajemen role (*Role-Based Access Control / RBAC*), dan endpoint RESTful berbasis JSON.

- **Bahasa**: Go (versi `1.27.0` pada [go.mod](file:///home/wenell/projects/lavendera-API/go.mod))
- **Entry Point**: [main.go](file:///home/wenell/projects/lavendera-API/main.go)
- **Module Name**: `github.com/Wenell09/lavendera-api`

---

## 2. Tech Stack

Berdasarkan implementasi aktual pada [go.mod](file:///home/wenell/projects/lavendera-API/go.mod):

| Kategori | Library / Tool | Versi / Keterangan |
|---|---|---|
| **Web Framework** | `github.com/gin-gonic/gin` | `v1.12.0` |
| **ORM & Database Driver** | `gorm.io/gorm`<br>`gorm.io/driver/postgres` | `v1.31.2`<br>`v1.6.2` (PostgreSQL) |
| **Dependency Injection** | `github.com/google/wire` | `v0.7.0` (Compile-time code generator) |
| **Configuration** | `github.com/spf13/viper` | `v1.21.0` (Membaca file `.env` & OS environment) |
| **Structured Logging** | `github.com/sirupsen/logrus` | `v1.10.2` (JSON Formatter ke `stdout`) |
| **Validation** | `github.com/go-playground/validator/v10` | `v10.30.4` + Custom validator `id_phone` |
| **JWT** | `github.com/golang-jwt/jwt/v5` | `v5.3.1` (HMAC SHA-256 / HS256) |
| **Password Hashing** | `golang.org/x/crypto/bcrypt` | `v0.55.0` (`bcrypt.DefaultCost`) |
| **UUID Generator** | `github.com/google/uuid` | `v1.6.0` (UUID v4) |
| **Phone Validation** | `github.com/nyaruka/phonenumbers/v2` | `v2.0.12` (Khusus format nomor Indonesia `ID`) |
| **Database Migration** | `github.com/golang-migrate/migrate/v4` | `v4.19.1` (SQL migration scripts di [internal/db/migrations](file:///home/wenell/projects/lavendera-API/internal/db/migrations)) |

---

## 3. Architecture Summary

Project mengadopsi variasi **Layered Clean Architecture / Modular Monolith**:
1. Setiap fitur/domain dibungkus ke dalam direktori independen di dalam `internal/<feature>/` (contoh: [internal/auth](file:///home/wenell/projects/lavendera-API/internal/auth), [internal/outlet](file:///home/wenell/projects/lavendera-API/internal/outlet), [internal/service](file:///home/wenell/projects/lavendera-API/internal/service), [internal/service_category](file:///home/wenell/projects/lavendera-API/internal/service_category), [internal/discount](file:///home/wenell/projects/lavendera-API/internal/discount)).
2. Setiap fitur memiliki 4 layer utama:
   - **`dto/`**: Data Transfer Object untuk request, filter, dan response.
   - **`controller/`**: Interface & implementasi Gin HTTP handler, input binding, validasi DTO, dan response mapping.
   - **`service/`**: Interface & implementasi business logic, tenant check, name conflict check, DTO-to-entity mapping.
   - **`repository/`**: Interface & implementasi database access via GORM dengan context dan `TenantScope`.
3. Komponen bersama disimpan di [internal/shared/](file:///home/wenell/projects/lavendera-API/internal/shared):
   - `appcontext`: Helper context untuk menyimpan dan mengekstrak `user_id`, `tenant_id`, `role`.
   - `apperror`: Struktur error terpusat (`NotFoundError`, `ValidationError`, `ConflictError`, `UnauthorizedError`) dan handler `NewHandleError`.
   - `applogger`: Logger terpusat berbasis Logrus JSON.
   - `appvalidator`: Inisialisasi Go Playground Validator dengan custom tag `id_phone`.
   - `config`: Viper configuration loader (`PORT`, `DATABASE_URL`, `JWT_SECRET`).
   - `database`: Inisialisasi koneksi GORM PostgreSQL dan GORM Scope `TenantScope`.
   - `middleware`: `JWTMiddleware`, `RequireRole`, `LoggerMiddleware`.
   - `response`: Format response konsisten `ResponseSuccess` dan `ResponseError`.
   - `utils`: Helper slug, validasi nomor HP Indonesia, JWT generator, context logger.
4. Model entitas database terpusat di [internal/models/](file:///home/wenell/projects/lavendera-API/internal/models).
5. Dependency Injection ditangani via **Google Wire** di [internal/app/wire.go](file:///home/wenell/projects/lavendera-API/internal/app/wire.go) dan di-generate ke [internal/app/wire_gen.go](file:///home/wenell/projects/lavendera-API/internal/app/wire_gen.go).

---

## 4. Important Directories

```
.
├── main.go                       # Entry point aplikasi
├── go.mod / go.sum               # Dependency management
├── .env.example / .env           # Environment variable template
├── internal/
│   ├── app/                      # Wiring aplikasi, router setup, providers
│   ├── auth/                     # Fitur Autentikasi (Register, Login)
│   ├── outlet/                   # Fitur Outlet CRUD
│   ├── service/                  # Fitur Layanan Laundry CRUD + Search + Filter + Pagination
│   ├── service_category/         # Fitur Kategori Layanan CRUD
│   ├── discount/                 # Fitur Diskon Promosi CRUD + Pagination
│   ├── models/                   # GORM Entity models
│   ├── db/migrations/            # File migrasi SQL up & down
│   └── shared/                   # Shared utilities, middleware, context, error, config
```

---

## 5. Coding Rules Observed in Codebase

1. **Explicit Interface & Implementation Separation**:
   Setiap layer (`repository`, `service`, `controller`) mendefinisikan interface di `<layer>.go` dan implementasi di `<layer>_impl.go`.
   - Contoh: [service_category_repository.go](file:///home/wenell/projects/lavendera-API/internal/service_category/repository/service_category_repository.go) dan [service_category_repository_impl.go](file:///home/wenell/projects/lavendera-API/internal/service_category/repository/service_category_repository_impl.go).
2. **Constructor Injection**:
   Setiap komponen memiliki fungsi konstruktor berawalan `New...` yang mengembalikan interface.
   - Contoh: `func NewServiceCategoryRepository(db *gorm.DB) ServiceCategoryRepository`
   - Konstruktor di-export dan didaftarkan pada `WireSet` masing-masing fitur.
3. **Context Propagation**:
   - Seluruh method repository dan service menerima `ctx context.Context` sebagai argumen pertama.
   - Gin handler mengoper `c.Request.Context()` ke service layer.
   - GORM queries selalu memanggil `.WithContext(ctx)`.
4. **Pointer untuk Update DTO (Partial Update)**:
   - Pada `Update...Request`, semua field opsional didefinisikan sebagai pointer (`*string`, `*int64`, `*bool`) dengan tag `validate:"omitempty,..."`.
   - Service menyusun `map[string]interface{}` hanya untuk field yang tidak nil sebelum mengupdate ke repository.
5. **No Direct Model Exposure**:
   Controller tidak pernah menerima model GORM secara langsung dari client ataupun mengembalikannya langsung. Selalu ada konversi DTO Request -> Entity Model -> DTO Response.

---

## 6. Feature Development Workflow

Jika ingin menambahkan fitur baru (misal `customer`), ikuti alur konsisten berikut:
```
1. Model: internal/models/<entity>.go
   ↓
2. DTO: internal/<feature>/dto/<feature>_{request,response,filter}.go
   ↓
3. Repository: internal/<feature>/repository/<feature>_repository{.go,_impl.go}
   ↓
4. Service: internal/<feature>/service/<feature>_service{.go,_impl.go}
   ↓
5. Controller: internal/<feature>/controller/<feature>_controller{.go,_impl.go}
   ↓
6. Routes: internal/<feature>/routes/<feature>_routes.go
   ↓
7. WireSet: internal/<feature>/wire.go
   ↓
8. Register ke internal/app/wire.go dan internal/app/routes.go
   ↓
9. Jalankan 'wire ./internal/app' untuk regenerasi wire_gen.go
```

---

## 7. Authentication Rules

1. **Header Format**: Request wajib menyertakan header `Authorization: Bearer <token>`.
2. **Token Verification**:
   - Algoritma wajib: `HS256`.
   - Secret key diambil dari `config.ENV.JWTSecret`.
   - Claims wajib: `user_id`, `tenant_id`, `role`.
   - Token diproses oleh [JWTMiddleware](file:///home/wenell/projects/lavendera-API/internal/shared/middleware/jwt_middleware.go).
3. **Context Injection**:
   Setelah token valid, middleware mengekstrak dan memasukkan data identity ke context:
   - `appcontext.WithUserID(ctx, userID)` (bertipe `uuid.UUID`)
   - `appcontext.WithTenantID(ctx, tenantID)` (bertipe `uuid.UUID`)
   - `appcontext.WithRole(ctx, role)` (bertipe `string`)
4. **Public Routes vs Protected Routes**:
   - Public: `/`, `/api/v1/auth/register`, `/api/v1/auth/login`.
   - Protected: Semua route di bawah `/api/v1` selain `/auth/*`.

---

## 8. Authorization & RBAC Rules

1. **Role Values**: Hanya 2 role yang didukung pada database constraint: `"ADMIN"` dan `"STAFF"` ([20260908061134_create_users.up.sql](file:///home/wenell/projects/lavendera-API/internal/db/migrations/20260908061134_create_users.up.sql)).
2. **Enforcement Middleware**: [RequireRole](file:///home/wenell/projects/lavendera-API/internal/shared/middleware/rbac_middleware.go) mengambil role dari context. Jika role tidak cocok, abort dengan HTTP `403 Forbidden`.
3. **Current Role Assignment**:
   - Seluruh route outlet, service category, service, dan discount saat ini dilindungi oleh `RequireRole("ADMIN")`.
   - Route staff dikomentari namun belum ada endpoint yang didaftarkan.

---

## 9. Tenant Isolation Rules

1. **Wajib Menggunakan `TenantScope`**:
   Semua operasi query, find, update, dan delete di repository yang berkaitan dengan data tenant **WAJIB** menambahkan:
   ```go
   s.DB.WithContext(ctx).Scopes(database.TenantScope(ctx))
   ```
2. **Perilaku TenantScope**:
   - Jika `tenant_id` tidak ada di context atau bernilai `uuid.Nil`, scope mengeksekusi `db.Where("1 = 0")`, memastikan query tidak mengembalikan atau memodifikasi data apapun.
   - Jika valid, menambahkan `db.Where("tenant_id = ?", tenantID)`.
3. **Create Entity**:
   Service layer wajib mengekstrak `tenant_id` via `appcontext.TenantIDFromContext(ctx)` dan menetapkannya ke struct model sebelum memanggil repository:
   ```go
   tenantID, exists := appcontext.TenantIDFromContext(ctx)
   if !exists {
       return nil, apperror.UnauthorizedError{Msg: "tenant_id not found"}
   }
   entity.TenantID = tenantID
   ```

---

## 10. Database Rules

1. **Primary Keys**: Menggunakan UUID v4 (`uuid_generate_v4()`).
2. **Foreign Keys**: Selalu memiliki relasi `REFERENCES <table_name>(id)` dengan klausul `ON DELETE CASCADE` atau `ON DELETE RESTRICT` / `SET NULL` sesuai kebutuhan relasi.
3. **Timestamps**:
   - Kolom `created_at` dan `updated_at` bertipe `TIMESTAMP WITH TIME ZONE DEFAULT NOW()`.
   - Update `updated_at` otomatis di-handle oleh database trigger `update_timestamp_column()` PostgreSQL.
4. **Transactions**:
   - Operasi multi-tabel wajib dibungkus dalam `tx := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error { ... })`.
   - Contoh implementasi aktual: pendaftaran tenant + admin + default categories di [AuthRepositoryImpl.CreateDefaultAdmin](file:///home/wenell/projects/lavendera-API/internal/auth/repository/auth_repository_impl.go#L19-L36).

---

## 11. API Conventions

1. **Response Envelope (Success)**:
   ```json
   {
     "status": 200,
     "message": "resource retrieved successfully",
     "success": true,
     "data": { ... },
     "meta": {
       "pagination": { ... } // opsional, hanya jika paginated list
     }
   }
   ```
2. **Response Envelope (Error)**:
   ```json
   {
     "status": 400,
     "message": "Validation Error",
     "success": false,
     "error": [ ... ] // atau string message
   }
   ```
3. **Pagination & Filter Query Params**:
   - `page`: default `1`
   - `limit`: default `10`, max `100`
   - `search`: string prefix matching (`name ILIKE search + "%"`)
4. **HTTP Status Codes**:
   - `200 OK`: Sukses GET, PATCH, DELETE
   - `201 Created`: Sukses POST (Create)
   - `400 Bad Request`: Body JSON parsing gagal atau validasi input gagal
   - `401 Unauthorized`: Token tidak ada/tidak valid, user non-aktif, atau password salah
   - `403 Forbidden`: Role tidak diizinkan mengakses resource
   - `404 Not Found`: Record tidak ditemukan di database (scoped per tenant)
   - `409 Conflict`: Data duplikat (nama layanan duplikat, email/slug sudah terdaftar)
   - `500 Internal Server Error`: Kesalahan server yang tidak tertangani

---

## 12. Error Handling Pattern

1. **Standard Error Helper**:
   Selalu gunakan [apperror.NewHandleError(c, err)](file:///home/wenell/projects/lavendera-API/internal/shared/apperror/error_handler.go) di controller untuk memformat response error HTTP.
2. **AppError Interface**:
   Jika error mengimplementasikan interface [AppError](file:///home/wenell/projects/lavendera-API/internal/shared/apperror/app_error.go), helper akan mengekstrak `StatusCode()`, `ResponseMessage()`, dan `ErrorData()`.
3. **Validator Error**:
   Jika `err` bertipe `validator.ValidationErrors`, bungkus terlebih dahulu dengan `apperror.NewFieldError(ve)` sebelum diteruskan ke `apperror.NewHandleError(c, ...)`.

---

## 13. Important Constraints

- **Single Token Role Enforcement**: Role diperiksa langsung dari JWT claims, bukan query ulang ke database setiap request.
- **Tenant Scope Everywhere**: Jangan pernah query tabel multi-tenant tanpa `TenantScope(ctx)`.
- **Database Unique Constraints**:
  - `tenants.slug` unik secara global.
  - `users.email` unik per tenant: `UNIQUE (tenant_id, email)`.
  - `outlets.slug` unik per tenant: `UNIQUE (tenant_id, slug)`.
  - `orders.order_number` unik per tenant: `UNIQUE (tenant_id, order_number)`.

---

## 14. Things AI Agents MUST NOT Do

1. **JANGAN** memanggil database langsung dari controller tanpa melalui service layer.
2. **JANGAN** membuat query ke tabel tenant tanpa menyertakan `Scopes(database.TenantScope(ctx))`.
3. **JANGAN** mempercayai `tenant_id` dari payload JSON request client jika endpoint berada di dalam protected group; selalu ambil `tenant_id` dari context JWT (`appcontext.TenantIDFromContext(ctx)`).
4. **JANGAN** mengubah file [internal/app/wire_gen.go](file:///home/wenell/projects/lavendera-API/internal/app/wire_gen.go) secara manual tanpa menjalankan `wire ./internal/app`.
5. **JANGAN** mengembalikan response raw GORM model ke client; selalu mapping ke DTO Response.
6. **JANGAN** menggunakan raw SQL query jika GORM scopes & methods sudah mencukupi pattern yang ada.
7. **JANGAN** mengabaikan validasi struct DTO dengan `validator.Validate`.
8. **JANGAN** mengasumsikan role lain selain `"ADMIN"` dan `"STAFF"` tersedia pada database schema.

---

## 15. Before Modifying Code

Sebelum membuat atau mengubah kode pada project ini, AI Agent **WAJIB** melakukan 10 langkah berikut:

1. **Cari implementasi fitur serupa**:
   - Jika membuat resource CRUD baru dengan pagination/search, buka [internal/service/](file:///home/wenell/projects/lavendera-API/internal/service/) atau [internal/discount/](file:///home/wenell/projects/lavendera-API/internal/discount/) sebagai referensi template utama.
   - Jika membuat master data sederhana tanpa pagination, buka [internal/service_category/](file:///home/wenell/projects/lavendera-API/internal/service_category/).
2. **Baca model database**: Periksa model di [internal/models/](file:///home/wenell/projects/lavendera-API/internal/models/) dan periksa file SQL migrasi di [internal/db/migrations/](file:///home/wenell/projects/lavendera-API/internal/db/migrations/) untuk memastikan tipe kolom, nullability, dan constraints.
3. **Baca DTO**: Periksa request/response/filter DTO fitur serupa di `internal/<feature>/dto/`.
4. **Baca repository**: Periksa interface repository dan metode `Scopes(database.TenantScope(ctx))` pada implementasinya.
5. **Baca service**: Periksa business validation (check exists by name, duplicate error handling, context tenant extraction).
6. **Baca controller**: Periksa alur request binding `ShouldBindJSON`, struct validation, UUID param parsing, dan pemanggilan `apperror.NewHandleError`.
7. **Baca routes**: Periksa registrasi route dan middleware group di `routes.go` masing-masing fitur.
8. **Baca middleware terkait**: Periksa apakah route membutuhkan `JWTMiddleware` dan `RequireRole("ADMIN")`.
9. **Periksa Dependency Injection**:
   - Tambahkan provider ke `WireSet` fitur (`internal/<feature>/wire.go`).
   - Daftarkan `WireSet` tersebut ke `internal/app/wire.go`.
   - Update `internal/app/routes.go` untuk menerima controller baru.
10. **Periksa dampak terhadap fitur lain**: Pastikan relasi foreign key dan cascade delete tidak merusak integritas data modul lain.
