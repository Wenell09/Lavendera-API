# Architecture Documentation — Lavendera API

> **Status**: Verifikasi dari implementasi aktual kode.  
> **Source of Truth**: [main.go](file:///home/wenell/projects/lavendera-API/main.go), [internal/app/](file:///home/wenell/projects/lavendera-API/internal/app/), [internal/shared/](file:///home/wenell/projects/lavendera-API/internal/shared/).

---

## 1. Architectural Style

Lavendera API mengadopsi arsitektur **Modular Monolith** dengan pendekatan **Layered Architecture (Controller-Service-Repository)** dan prinsip Dependency Inversion (Interface-driven design).

Karakteristik arsitektur yang teridentifikasi:
- **Feature-based Packaging**: Kode aplikasi diorganisasikan per fitur/domain (`auth`, `outlet`, `service`, `service_category`, `discount`) di dalam direktori `internal/`.
- **Inverted Dependencies**: Layer yang lebih tinggi (misal Controller) bergantung pada interface yang didefinisikan untuk layer berikutnya (Service), bukan pada implementasi konkretnya.
- **Compile-time Dependency Injection**: Menggunakan Google Wire (`github.com/google/wire`) untuk merangkai objek saat build-time tanpa menggunakan runtime reflection.
- **Centralized Multi-Tenancy**: Isolasi data antar-tenant diberlakukan pada data-access layer menggunakan custom GORM Scope `TenantScope(ctx)`.

---

## 2. Layers & Responsibilities

Setiap modul fitur terbagi menjadi 4 layer standar:

```
┌────────────────────────────────────────────────────────┐
│                   HTTP Request                         │
└──────────────────────────┬─────────────────────────────┘
                           │
                           ▼
┌────────────────────────────────────────────────────────┐
│                     Routing Layer                      │
│        (gin.Engine / RouterGroup / Middleware)         │
└──────────────────────────┬─────────────────────────────┘
                           │
                           ▼
┌────────────────────────────────────────────────────────┐
│                   Controller Layer                     │
│  - Bind JSON / Query Params / URL Params               │
│  - Validate input via validator.Validate               │
│  - Call Service layer via interface                    │
│  - Handle errors via apperror.NewHandleError           │
│  - Format response via response.ResponseSuccess/Error   │
└──────────────────────────┬─────────────────────────────┘
                           │
                           ▼
┌────────────────────────────────────────────────────────┐
│                    Service Layer                       │
│  - Business logic enforcement                          │
│  - Check name conflicts / business rules               │
│  - Extract tenant_id & identity from Context           │
│  - Map DTO Request <-> GORM Entity Models              │
│  - Call Repository layer via interface                 │
│  - Structured logging with tenant context              │
└──────────────────────────┬─────────────────────────────┘
                           │
                           ▼
┌────────────────────────────────────────────────────────┐
│                   Repository Layer                     │
│  - Direct Database queries via *gorm.DB                │
│  - Apply database.TenantScope(ctx)                     │
│  - Execute SQL queries, transactions, preloads         │
│  - Return GORM models or errors                        │
└──────────────────────────┬─────────────────────────────┘
                           │
                           ▼
┌────────────────────────────────────────────────────────┐
│                   PostgreSQL Database                  │
└────────────────────────────────────────────────────────┘
```

### Rincian Tanggung Jawab Tiap Layer:

1. **Routing & Middleware Layer**:
   - Menangani routing HTTP melalui Gin.
   - Menjalankan pipeline middleware: Logging ([LoggerMiddleware](file:///home/wenell/projects/lavendera-API/internal/shared/middleware/logger_middleware.go)), Recovery, JWT authentication ([JWTMiddleware](file:///home/wenell/projects/lavendera-API/internal/shared/middleware/jwt_middleware.go)), dan role authorization ([RequireRole](file:///home/wenell/projects/lavendera-API/internal/shared/middleware/rbac_middleware.go)).

2. **Controller Layer**:
   - Membaca input dari HTTP request (`c.ShouldBindJSON`, `c.Param`, `c.Query`, `c.DefaultQuery`).
   - Melakukan validasi deklaratif menggunakan `validator.Validate`.
   - Meneruskan data request yang telah tervalidasi ke Service interface.
   - Mengembalikan output HTTP menggunakan format response standar (`response.ResponseSuccess` atau `response.ResponseError`).

3. **Service Layer**:
   - Menjalankan seluruh logika bisnis (validasi unik bisnis seperti `ExistsByName`, hashing password, logic pembatasan diskon).
   - Membaca konteks request (`appcontext.TenantIDFromContext(ctx)`).
   - Melakukan mapping antara DTO Request dan Model Database.
   - Mencatat log operasional berstruktur dengan menyertakan context tenant ID (`utils.LogWithContext`).

4. **Repository Layer**:
   - Bertanggung jawab penuh terhadap akses data persistence (PostgreSQL via GORM).
   - Menjamin isolasi data tenant menggunakan `Scopes(database.TenantScope(ctx))`.
   - Mengoperasikan pagination (limit, offset, count) dan pencarian filter (`ILIKE`).

5. **Entity Models Layer** ([internal/models](file:///home/wenell/projects/lavendera-API/internal/models)):
   - Representasi struktur tabel GORM yang dibagikan ke seluruh modul aplikasi.

---

## 3. Dependency Direction

Dependency mengalir secara searah dari luar ke dalam:

```
Routes ──> Controller ──> Service ──> Repository ──> *gorm.DB / Models
  │            │             │              │
  └────────────┴─────────────┴──────────────┴──────> internal/shared
```

Aturan ketergantungan yang teramati:
- Controller bergantung pada **Service Interface**.
- Service bergantung pada **Repository Interface**.
- Repository bergantung pada `*gorm.DB` dan `internal/models`.
- Repository dan Service **TIDAK PERNAH** bergantung pada Gin context (`*gin.Context`), melainkan standar `context.Context`.
- Semua layer dapat bergantung pada modul pembantu di `internal/shared/` (`apperror`, `appcontext`, `response`, `utils`, `applogger`, `database`).
- Tidak ditemukan circular dependency antar fitur (misal modul `service` hanya memanggil `categoryRepository.ServiceCategoryRepository` yang diinjeksi melalui konstruktor via Wire).

---

## 4. Application Startup Lifecycle

Urutan eksekusi aplikasi saat dimulai ([main.go](file:///home/wenell/projects/lavendera-API/main.go)):

```
 1. main()
    │
    ├──> 2. config.LoadConfig()
    │       - Inisialisasi Viper (".env", type "env", default PORT="8080")
    │       - Membaca OS Environment Variables
    │       - Unmarshal ke config.ENV
    │       - Validasi FATAL: DATABASE_URL & JWT_SECRET tidak boleh kosong
    │
    ├──> 3. app.InitializeApp() (via Google Wire)
    │       ├── NewDB() -> database.DBConnection()
    │       │     - Open gorm.Open(postgres.Open(dsn))
    │       │     - Set Connection Pool (MaxIdleConns: 10, MaxOpenConns: 100, MaxLifetime: 1h)
    │       ├── applogger.NewLogger()
    │       ├── config.NewJWTConfig()
    │       ├── appvalidator.NewValidator() (Register custom rule: id_phone)
    │       ├── Instansiasi seluruh Repository (Auth, ServiceCategory, Service, Discount, Outlet)
    │       ├── Instansiasi seluruh Service
    │       ├── Instansiasi seluruh Controller
    │       ├── NewRouter(...) -> Buat gin.Engine, pasang middleware, daftarkan route
    │       └── NewApp(engine, db, logger)
    │
    ├──> 4. Ambil sqlDB := application.DB.DB(), daftarkan defer sqlDB.Close()
    │
    └──> 5. application.Router.Run(":" + config.ENV.Port)
```

---

## 5. Dependency Injection (Google Wire)

Injeksi dependensi dirangkai secara deklaratif menggunakan **Google Wire**.

### Komponen Wire:
1. **Set Fitur Individu**:
   - `auth.WireSet`: Repository, Service, Controller, JWTConfig.
   - `service_category.WireSet`: Repository, Service, Controller.
   - `service.WireSet`: Repository, Service, Controller.
   - `discount.WireSet`: Repository, Service, Controller.
   - `outlet.WireSet`: Repository, Service, Controller.
2. **Injector Utama** ([internal/app/wire.go](file:///home/wenell/projects/lavendera-API/internal/app/wire.go)):
   ```go
   func InitializeApp() (*App, error) {
       wire.Build(
           auth.WireSet,
           service_category.WireSet,
           service.WireSet,
           discount.WireSet,
           outlet.WireSet,
           appvalidator.NewValidator,
           applogger.NewLogger,
           NewDB,
           NewRouter,
           NewApp,
       )
       return nil, nil
   }
   ```
3. **Hasil Generate** ([internal/app/wire_gen.go](file:///home/wenell/projects/lavendera-API/internal/app/wire_gen.go)):
   Wire mengurai pohon dependensi secara otomatis saat `go run github.com/google/wire/cmd/wire` dijalankan.

---

## 6. Request Lifecycle

Berikut adalah alur lengkap pemrosesan HTTP request:

```
Client
  │  HTTP Request (Header: Authorization: Bearer <jwt>)
  ▼
[LoggerMiddleware] (Catat start time)
  │
  ▼
[Recovery] (Tangkal panic)
  │
  ▼
[JWTMiddleware] (Jika protected route)
  │  ├── Parse header "Authorization" -> format Bearer
  │  ├── Verifikasi token HS256 dengan JWTSecret
  │  ├── Ekstrak claims: user_id, tenant_id, role
  │  └── Injeksi identity ke c.Request.Context()
  ▼
[RequireRole("ADMIN")] (Jika admin route)
  │  └── Periksa role dari context; abort 403 jika tidak cocok
  ▼
[Controller]
  │  ├── Bind JSON / Param
  │  ├── Validasi Struct via appvalidator
  │  └── Panggil Service(ctx, req)
  ▼
[Service]
  │  ├── Ambil tenant_id dari context
  │  ├── Periksa aturan bisnis (cek eksistensi nama / duplikasi)
  │  ├── Siapkan Entity Model
  │  └── Panggil Repository(ctx, ...)
  ▼
[Repository]
  │  ├── Terapkan Scopes(database.TenantScope(ctx)) -> WHERE tenant_id = ?
  │  └── Eksekusi GORM ke PostgreSQL
  ▼
[Service]
  │  └── Mapping Model Entity ke DTO Response
  ▼
[Controller]
  │  └── Kirim JSON Response (response.NewResponseSuccess)
  ▼
[LoggerMiddleware]
  │  └── Hitung latency, cetak log JSON terstruktur (status, latency_ms, ip, method, path)
  ▼
Client (Menerima HTTP Response)
```

---

## 7. Shared Components Directory Map

| Komponen | Path | Tanggung Jawab |
|---|---|---|
| `appcontext` | [internal/shared/appcontext](file:///home/wenell/projects/lavendera-API/internal/shared/appcontext) | Kunci context dan accessor/mutator typed UUID untuk `user_id`, `tenant_id`, `role`. |
| `apperror` | [internal/shared/apperror](file:///home/wenell/projects/lavendera-API/internal/shared/apperror) | Interface `AppError`, tipe error spesifik (`NotFoundError`, `ValidationError`, `ConflictError`, `UnauthorizedError`), dan fungsi penangan terpusat `NewHandleError`. |
| `applogger` | [internal/shared/applogger](file:///home/wenell/projects/lavendera-API/internal/shared/applogger) | Inisialisasi Logrus logger dengan format JSON dan level Info ke `stdout`. |
| `appvalidator` | [internal/shared/appvalidator](file:///home/wenell/projects/lavendera-API/internal/shared/appvalidator) | Inisialisasi validator Go Playground dengan custom rule `id_phone`. |
| `config` | [internal/shared/config](file:///home/wenell/projects/lavendera-API/internal/shared/config) | Konfigurasi aplikasi via Viper (`PORT`, `DATABASE_URL`, `JWT_SECRET`) dan `JWTConfig`. |
| `database` | [internal/shared/database](file:///home/wenell/projects/lavendera-API/internal/shared/database) | Manajemen koneksi GORM PostgreSQL, pool koneksi, dan GORM scope `TenantScope`. |
| `middleware` | [internal/shared/middleware](file:///home/wenell/projects/lavendera-API/internal/shared/middleware) | Gin middleware: `JWTMiddleware`, `RequireRole`, `LoggerMiddleware`. |
| `response` | [internal/shared/response](file:///home/wenell/projects/lavendera-API/internal/shared/response) | Envelope standar format JSON untuk response sukses (`ResponseSuccess`) dan error (`ResponseError`). |
| `utils` | [internal/shared/utils](file:///home/wenell/projects/lavendera-API/internal/shared/utils) | Generator slug regex, helper JWT HS256, validasi nomor HP Indonesia via libphonenumber, dan context logger. |
