# Features Matrix & Analysis — Lavendera API

> **Status**: Verifikasi dari implementasi aktual kode per 28 September 2026.  
> **Source of Truth**: Seluruh modul di `internal/`, model di `internal/models/`, dan script migrasi di `internal/db/migrations/`.

---

## 1. Feature Matrix

Status diklasifikasikan menjadi:
- **DONE**: Fitur diimplementasikan lengkap (Model, Migration, Repository, Service, Controller, Routes).
- **PARTIAL**: Sebagian layer telah diimplementasikan (misal hanya Model/Auth atau belum ada CRUD).
- **TODO**: Telah didefinisikan skema database dan Model struct, namun belum ada layer API/Service/Repository.
- **UNKNOWN**: Belum dapat dikonfirmasi dari implementasi saat ini.

| Feature / Domain | Status | API | Repository | Service | Controller | Notes |
|---|---|---|---|---|---|---|
| **Authentication & Registration** | **DONE** | Ya | Ya | Ya | Ya | Register tenant + default admin + login JWT HS256. |
| **Outlet Management** | **DONE** | Ya | Ya | Ya | Ya | CRUD lengkap, pagination, search prefix, slug otomatis. |
| **Service Category** | **DONE** | Ya | Ya | Ya | Ya | CRUD lengkap untuk kategori layanan (Kiloan, Satuan). |
| **Service (Layanan Laundry)** | **DONE** | Ya | Ya | Ya | Ya | CRUD lengkap, pagination, search, filter `category_id`. |
| **Discounts** | **DONE** | Ya | Ya | Ya | Ya | CRUD lengkap, pagination, search, auto cap 100%. |
| **User / Staff Management** | **DONE** | Ya | Ya | Ya | Ya | CRUD lengkap, pagination, search, filter role, proteksi hapus akun sendiri. |
| **Outlet User Assignment** | **TODO** | Tidak | Tidak | Tidak | Tidak | Model `OutletUser` dan tabel migrasi `outlet_users` ada, namun belum ada layer aplikasi. |
| **Outlet Payment Methods** | **TODO** | Tidak | Tidak | Tidak | Tidak | Model `OutletPaymentMethod` dan tabel migrasi ada, belum ada layer aplikasi. |
| **Customer Management** | **TODO** | Tidak | Tidak | Tidak | Tidak | Model `Customer` dan tabel migrasi ada, belum ada layer aplikasi. |
| **Order Management** | **TODO** | Tidak | Tidak | Tidak | Tidak | Model `Order` dan `OrderItem` serta tabel migrasi ada, belum ada layer aplikasi. |
| **Payment Management** | **TODO** | Tidak | Tidak | Tidak | Tidak | Model `Payment` dan tabel migrasi ada, belum ada layer aplikasi. |
| **Notifications** | **TODO** | Tidak | Tidak | Tidak | Tidak | Model `Notification` dan tabel migrasi ada, belum ada layer aplikasi. |
| **Audit Logs** | **TODO** | Tidak | Tidak | Tidak | Tidak | Model `AuditLog` dan tabel migrasi ada, belum ada layer aplikasi. |
| **Staff Specific Routes** | **PARTIAL** | Tidak | N/A | N/A | N/A | Komentar placeholder ada di [routes.go](file:///home/wenell/projects/lavendera-API/internal/app/routes.go#L55), belum ada endpoint staff. |

---

## 2. Detailed Feature Analysis (Implemented Modules)

### 2.1. Authentication & Onboarding
- **Purpose**: Pendaftaran tenant baru beserta akun admin pertama dan login untuk mendapatkan token JWT.
- **Status**: **DONE**
- **Model**: [models.Tenant](file:///home/wenell/projects/lavendera-API/internal/models/tenant.go), [models.User](file:///home/wenell/projects/lavendera-API/internal/models/user.go), [models.ServiceCategory](file:///home/wenell/projects/lavendera-API/internal/models/service_category.go)
- **Repository**: [AuthRepository](file:///home/wenell/projects/lavendera-API/internal/auth/repository/auth_repository.go) / [AuthRepositoryImpl](file:///home/wenell/projects/lavendera-API/internal/auth/repository/auth_repository_impl.go)
- **Service**: [AuthService](file:///home/wenell/projects/lavendera-API/internal/auth/service/auth_service.go) / [AuthServiceImpl](file:///home/wenell/projects/lavendera-API/internal/auth/service/auth_service_impl.go)
- **Controller**: [AuthController](file:///home/wenell/projects/lavendera-API/internal/auth/controller/auth_controller.go) / [AuthControllerImpl](file:///home/wenell/projects/lavendera-API/internal/auth/controller/auth_controller_impl.go)
- **Routes**: [auth_routes.go](file:///home/wenell/projects/lavendera-API/internal/auth/routes/auth_routes.go) (`POST /register`, `POST /login`)
- **Middleware**: Public (tidak ada middleware auth).
- **Request DTO**: `RegisterRequest`, `LoginRequest` di [auth_request.go](file:///home/wenell/projects/lavendera-API/internal/auth/dto/auth_request.go)
- **Response DTO**: `RegisterResponse`, `LoginResponse`, `TenantResponse`, `UserResponse` di [auth_response.go](file:///home/wenell/projects/lavendera-API/internal/auth/dto/auth_response.go)
- **Validation**:
  - `RegisterRequest`: `tenant_name` (3-100), `name` (2-100), `email` (email, max 255), `password` (8-72).
  - `LoginRequest`: `email` (email, max 255), `password` (8-72).
- **Business Rules**:
  - Email dan tenant slug dicek keunikannya sebelum insert.
  - Slug dibuat otomatis menggunakan fungsi regex `utils.GenerateSlug(req.TenantName)`.
  - Saat pendaftaran tenant, otomatis dibuatkan user admin (`Role: "ADMIN"`, `IsActive: true`) dan 2 kategori layanan awal: `"Kiloan"` dan `"Satuan"`.
  - Operasi insert dieksekusi dalam satu transaksi database GORM.
  - Login memvalidasi status `is_active` user. Jika non-aktif, login ditolak dengan HTTP 401.

---

### 2.2. Outlet Management
- **Purpose**: Pengelolaan data cabang outlet laundry per tenant.
- **Status**: **DONE**
- **Model**: [models.Outlet](file:///home/wenell/projects/lavendera-API/internal/models/outlet.go)
- **Repository**: [OutletRepository](file:///home/wenell/projects/lavendera-API/internal/outlet/repository/outlet_repository.go) / [OutletRepositoryImpl](file:///home/wenell/projects/lavendera-API/internal/outlet/repository/outlet_repository_impl.go)
- **Service**: [OutletService](file:///home/wenell/projects/lavendera-API/internal/outlet/service/outlet_service.go) / [OutletServiceImpl](file:///home/wenell/projects/lavendera-API/internal/outlet/service/outlet_service_impl.go)
- **Controller**: [OutletController](file:///home/wenell/projects/lavendera-API/internal/outlet/controller/outlet_controller.go) / [OutletControllerImpl](file:///home/wenell/projects/lavendera-API/internal/outlet/controller/outlet_controller_impl.go)
- **Routes**: [outlet_routes.go](file:///home/wenell/projects/lavendera-API/internal/outlet/routes/outlet_routes.go) (`POST`, `GET`, `GET /:id`, `PATCH /:id`, `DELETE /:id`)
- **Middleware**: `JWTMiddleware` + `RequireRole("ADMIN")`
- **Request DTO**: `CreateOutletRequest`, `UpdateOutletRequest` di [outlet_request.go](file:///home/wenell/projects/lavendera-API/internal/outlet/dto/outlet_request.go)
- **Response DTO**: `OutletResponse`, `OutletListResponse` di [outlet_response.go](file:///home/wenell/projects/lavendera-API/internal/outlet/dto/outlet_response.go)
- **Validation**:
  - `name`: min 5, max 255.
  - `phone`: wajib nomor telepon valid Indonesia (`id_phone`).
- **Pagination & Search**:
  - `search`: Filter prefix nama (`name ILIKE search + "%"`).
  - `page`: default 1, `limit`: default 10, max 100.
- **Special Behavior**:
  - `Create`: Slug dibuat otomatis dari nama outlet via `utils.GenerateSlug(req.Name)`.
  - `Update`: Slug dapat di-override melalui body `UpdateOutletRequest.Slug`.

---

### 2.3. Service Category Management
- **Purpose**: Pengelompokan jenis layanan (misal Kiloan, Satuan, Karpet, Sepatu).
- **Status**: **DONE**
- **Model**: [models.ServiceCategory](file:///home/wenell/projects/lavendera-API/internal/models/service_category.go)
- **Repository**: [ServiceCategoryRepository](file:///home/wenell/projects/lavendera-API/internal/service_category/repository/service_category_repository.go) / [ServiceCategoryRepositoryImpl](file:///home/wenell/projects/lavendera-API/internal/service_category/repository/service_category_repository_impl.go)
- **Service**: [ServiceCategoryService](file:///home/wenell/projects/lavendera-API/internal/service_category/service/service_category_service.go) / [ServiceCategoryServiceImpl](file:///home/wenell/projects/lavendera-API/internal/service_category/service/service_category_service_impl.go)
- **Controller**: [ServiceCategoryController](file:///home/wenell/projects/lavendera-API/internal/service_category/controller/service_category_controller.go) / [ServiceCategoryControllerImpl](file:///home/wenell/projects/lavendera-API/internal/service_category/controller/service_category_controller_impl.go)
- **Routes**: [service_category_routes.go](file:///home/wenell/projects/lavendera-API/internal/service_category/routes/service_category_routes.go) (`POST`, `GET`, `GET /:id`, `PATCH /:id`, `DELETE /:id`)
- **Middleware**: `JWTMiddleware` + `RequireRole("ADMIN")`
- **Request DTO**: `CreateServiceCategoryRequest`, `UpdateServiceCategoryRequest`
- **Response DTO**: `ServiceCategoryResponse`
- **Validation**: `name`: min 2, max 100.
- **Pagination & Search**: Tidak ada pagination; mengambil seluruh data diurutkan berdasarkan `created_at DESC`.

---

### 2.4. Service (Layanan Laundry) Management
- **Purpose**: Pengelolaan katalog harga dan durasi layanan laundry.
- **Status**: **DONE**
- **Model**: [models.Service](file:///home/wenell/projects/lavendera-API/internal/models/service.go)
- **Repository**: [ServiceRepository](file:///home/wenell/projects/lavendera-API/internal/service/repository/service_repository.go) / [ServiceRepositoryImpl](file:///home/wenell/projects/lavendera-API/internal/service/repository/service_repository_impl.go)
- **Service**: [Service](file:///home/wenell/projects/lavendera-API/internal/service/service/service.go) / [ServiceImpl](file:///home/wenell/projects/lavendera-API/internal/service/service/service_impl.go)
- **Controller**: [ServiceController](file:///home/wenell/projects/lavendera-API/internal/service/controller/service_controller.go) / [ServiceControllerImpl](file:///home/wenell/projects/lavendera-API/internal/service/controller/service_controller_impl.go)
- **Routes**: [service_routes.go](file:///home/wenell/projects/lavendera-API/internal/service/routes/service_routes.go) (`POST`, `GET`, `GET /:id`, `PATCH /:id`, `DELETE /:id`)
- **Middleware**: `JWTMiddleware` + `RequireRole("ADMIN")`
- **Request DTO**: `CreateServiceRequest`, `UpdateServiceRequest` di [service_request.go](file:///home/wenell/projects/lavendera-API/internal/service/dto/service_request.go)
- **Response DTO**: `ServiceResponse`, `CategoryResponse`, `ServiceListResponse` di [service_response.go](file:///home/wenell/projects/lavendera-API/internal/service/dto/service_response.go)
- **Validation**:
  - `name`: min 3, max 255.
  - `price`: gte 0.
  - `min_quantity`: gt 0.
  - `unit`: required (`gram`, `pcs`, `pair` pada constraint SQL).
  - `duration_days`: gt 0.
- **Pagination & Search**:
  - `search`: Filter prefix nama (`name ILIKE search + "%"`).
  - `category_id`: Filter spesifik kategori layanan.
  - `page`: default 1, `limit`: default 10, max 100.
- **Relasi Khusus**: Service menginjeksi `CategoryRepository` untuk memvalidasi keberadaan `category_id` per tenant saat update layanan.

---

### 2.5. Discount Management
- **Purpose**: Pengelolaan kupon/potongan harga transaksi (persentase atau nominal tetap).
- **Status**: **DONE**
- **Model**: [models.Discount](file:///home/wenell/projects/lavendera-API/internal/models/discount.go)
- **Repository**: [DiscountRepository](file:///home/wenell/projects/lavendera-API/internal/discount/repository/discount_repository.go) / [DiscountRepositoryImpl](file:///home/wenell/projects/lavendera-API/internal/discount/repository/discount_repository_impl.go)
- **Service**: [DiscountService](file:///home/wenell/projects/lavendera-API/internal/discount/service/discount_service.go) / [DiscountServiceImpl](file:///home/wenell/projects/lavendera-API/internal/discount/service/discount_service_impl.go)
- **Controller**: [DiscountController](file:///home/wenell/projects/lavendera-API/internal/discount/controller/discount_controller.go) / [DiscountControllerImpl](file:///home/wenell/projects/lavendera-API/internal/discount/controller/discount_controller_impl.go)
- **Routes**: [discount_routes.go](file:///home/wenell/projects/lavendera-API/internal/discount/routes/discount_routes.go) (`POST`, `GET`, `GET /:id`, `PATCH /:id`, `DELETE /:id`)
- **Middleware**: `JWTMiddleware` + `RequireRole("ADMIN")`
- **Request DTO**: `CreateDiscountRequest`, `UpdateDiscountRequest`
- **Response DTO**: `DiscountResponse`, `DiscountListResponse`
- **Validation**: `type`: `oneof=percentage fixed`, `value`: `gt=0`.
- **Special Business Rule**:
  - Pada `Create` dan `Update`, jika `type == "percentage"` dan `value > 100`, value dipaksa bernilai maksimal `100`.

---

### 2.6. User / Staff Management
- **Purpose**: Pengelolaan staf dan admin tambahan per tenant.
- **Status**: **DONE**
- **Model**: [models.User](file:///home/wenell/projects/lavendera-API/internal/models/user.go)
- **Repository**: [UserRepository](file:///home/wenell/projects/lavendera-API/internal/user/repository/user_repository.go) / [UserRepositoryImpl](file:///home/wenell/projects/lavendera-API/internal/user/repository/user_repository_impl.go)
- **Service**: [UserService](file:///home/wenell/projects/lavendera-API/internal/user/service/user_service.go) / [UserServiceImpl](file:///home/wenell/projects/lavendera-API/internal/user/service/user_service_impl.go)
- **Controller**: [UserController](file:///home/wenell/projects/lavendera-API/internal/user/controller/user_controller.go) / [UserControllerImpl](file:///home/wenell/projects/lavendera-API/internal/user/controller/user_controller_impl.go)
- **Routes**: [user_routes.go](file:///home/wenell/projects/lavendera-API/internal/user/routes/user_routes.go) (`POST`, `GET`, `GET /:id`, `PATCH /:id`, `DELETE /:id`)
- **Middleware**: `JWTMiddleware` + `RequireRole("ADMIN")`
- **Request DTO**: `CreateUserRequest`, `UpdateUserRequest`
- **Response DTO**: `UserResponse`, `UserListResponse`
- **Validation**:
  - `name`: min 2, max 150
  - `email`: valid email
  - `password`: min 8, max 72
  - `role`: oneof `ADMIN` `STAFF`
- **Special Business Rule**:
  - Pengecekan unik `email` per tenant saat create dan update.
  - Saat `DELETE`, mencegah akun yang sedang login (ID sama dengan di token JWT) menghapus dirinya sendiri (`apperror.ConflictError`).

---

## 3. Unimplemented / Planned Modules (Database Models Only)

Tabel dan model berikut telah tersedia di skema database ([internal/models/](file:///home/wenell/projects/lavendera-API/internal/models/) dan migrasi SQL), tetapi **belum diimplementasikan** layer Repository, Service, Controller, maupun Route-nya:

1. **`Customer`** ([internal/models/customer.go](file:///home/wenell/projects/lavendera-API/internal/models/customer.go)): Pengelolaan pelanggan laundry per tenant.
2. **`Order` & `OrderItem`** ([internal/models/order.go](file:///home/wenell/projects/lavendera-API/internal/models/order.go), [internal/models/order_item.go](file:///home/wenell/projects/lavendera-API/internal/models/order_item.go)): Pembuatan dan tracking pesanan laundry, status order (`MENUNGGU_KONFIRMASI`, `DIPROSES`, `SELESAI`, `DIAMBIL`, `DIBATALKAN`), serta perhitungan total & diskon.
3. **`Payment` & `OutletPaymentMethod`** ([internal/models/payment.go](file:///home/wenell/projects/lavendera-API/internal/models/payment.go), [internal/models/outlet_payment_method.go](file:///home/wenell/projects/lavendera-API/internal/models/outlet_payment_method.go)): Transaksi pembayaran (cash, QRIS, bank transfer), verifikasi, dan konfigurasi pembayaran outlet.
4. **`Notification`** ([internal/models/notification.go](file:///home/wenell/projects/lavendera-API/internal/models/notification.go)): Notifikasi pesanan atau sistem kepada user.
5. **`AuditLog`** ([internal/models/audit_logs.go](file:///home/wenell/projects/lavendera-API/internal/models/audit_logs.go)): Pencatatan riwayat perubahan data (*old_data* dan *new_data* JSONB).
6. **`OutletUser`** ([internal/models/outlet_user.go](file:///home/wenell/projects/lavendera-API/internal/models/outlet_user.go)): Pemetaan penugasan staff/user ke outlet tertentu.
