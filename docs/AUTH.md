# Authentication & Authorization Documentation — Lavendera API

> **Status**: Verifikasi dari implementasi aktual kode.  
> **Source of Truth**: [internal/auth/](file:///home/wenell/projects/lavendera-API/internal/auth/), [internal/shared/middleware/](file:///home/wenell/projects/lavendera-API/internal/shared/middleware/), [internal/shared/appcontext/](file:///home/wenell/projects/lavendera-API/internal/shared/appcontext/), [internal/shared/utils/jwt.go](file:///home/wenell/projects/lavendera-API/internal/shared/utils/jwt.go).

---

## 1. Authentication Architecture

Lavendera API menggunakan arsitektur autentikasi berbasis **JSON Web Token (JWT)** stateless:
- **Signing Algorithm**: HMAC-SHA256 (`HS256`).
- **Secret Key**: Dibaca dari file `.env` melalui environment variable `JWT_SECRET`.
- **Password Hashing**: Menggunakan `golang.org/x/crypto/bcrypt` dengan `bcrypt.DefaultCost` (cost 10).

---

## 2. Registration Flow

**Endpoint**: `POST /api/v1/auth/register` (Public)

1. Client mengirim payload JSON:
   - `tenant_name`: Nama laundry tenant (min 3, max 100).
   - `name`: Nama lengkap admin owner (min 2, max 100).
   - `email`: Email admin (valid email, max 255).
   - `password`: Password (min 8, max 72 karakter).
2. Controller [AuthControllerImpl.Register](file:///home/wenell/projects/lavendera-API/internal/auth/controller/auth_controller_impl.go#L49-L71) memvalidasi struct request via `validator.Validate`.
3. Service [AuthServiceImpl.Register](file:///home/wenell/projects/lavendera-API/internal/auth/service/auth_service_impl.go#L67-L110) mengeksekusi logika:
   - Normalisasi email: `strings.ToLower(strings.TrimSpace(req.Email))`.
   - Generate slug tenant otomatis dari nama tenant: `slug := utils.GenerateSlug(req.TenantName)`.
   - Pengecekan ketersediaan email: `Repository.ExistsUserByEmail(ctx, email)`. Jika sudah ada, return `apperror.ConflictError{Msg: "email already registered"}` (HTTP 409).
   - Pengecekan ketersediaan slug: `Repository.ExistsTenantBySlug(ctx, slug)`. Jika sudah ada, return `apperror.ConflictError{Msg: "tenant slug already exists"}` (HTTP 409).
   - Hashing password: `bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)`.
   - Setup entity:
     - Tenant: `&models.Tenant{Name: req.TenantName, Slug: slug, Email: email}`
     - User: `&models.User{Name: req.Name, Email: email, Password: string(hashedPassword), Role: "ADMIN", IsActive: true}`
     - Service Categories default: `[]models.ServiceCategory{{Name: "Kiloan"}, {Name: "Satuan"}}`
4. Repository [AuthRepositoryImpl.CreateDefaultAdmin](file:///home/wenell/projects/lavendera-API/internal/auth/repository/auth_repository_impl.go#L19-L36) mengeksekusi pembuatan Tenant, Admin User, dan Kategori Layanan default dalam satu database transaction `tx.Transaction`.
5. Mengembalikan HTTP `201 Created` dengan data Tenant dan User (tanpa mengembalikan token JWT secara otomatis; user diharapkan melakukan login terlebih dahulu).

---

## 3. Login Flow

**Endpoint**: `POST /api/v1/auth/login` (Public)

1. Client mengirim JSON berisi `email` dan `password`.
2. Controller [AuthControllerImpl.Login](file:///home/wenell/projects/lavendera-API/internal/auth/controller/auth_controller_impl.go#L25-L47) memvalidasi format input.
3. Service [AuthServiceImpl.Login](file:///home/wenell/projects/lavendera-API/internal/auth/service/auth_service_impl.go#L29-L65):
   - Mencari user di database: `Repository.FindUserByEmail(ctx, email)`.
   - Jika record tidak ditemukan (`gorm.ErrRecordNotFound`): mengembalikan `apperror.UnauthorizedError{Msg: "invalid email or password"}` (HTTP 401).
   - Memeriksa keaktifan akun `if !user.IsActive`: mengembalikan `apperror.UnauthorizedError{Msg: "user account is inactive"}` (HTTP 401).
   - Memeriksa kecocokan password via `bcrypt.CompareHashAndPassword`: jika gagal, mengembalikan `apperror.UnauthorizedError{Msg: "invalid email or password"}` (HTTP 401).
   - Menerbitkan token JWT:
     ```go
     token, err := utils.GenerateToken(user.ID.String(), user.TenantID.String(), user.Role, a.JWTConfig)
     ```
4. Mengembalikan HTTP `200 OK` dengan user profile dan string token JWT.

---

## 4. JWT Token Specification & Claims

Didefinisikan pada [internal/shared/utils/jwt.go](file:///home/wenell/projects/lavendera-API/internal/shared/utils/jwt.go):

```go
type JWTClaims struct {
    UserID   string `json:"user_id"`
    TenantID string `json:"tenant_id"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}
```

### Karakteristik Token Saat Ini:
- **Header**: `{"alg": "HS256", "typ": "JWT"}`
- **Payload Claims**:
  - `user_id`: string representasi UUID v4 dari entitas User.
  - `tenant_id`: string representasi UUID v4 dari entitas Tenant pemilik User.
  - `role`: string role user (`"ADMIN"` atau `"STAFF"`).
- **Expiration Behavior**:
  > [!NOTE]
  > Pada fungsi `utils.GenerateToken`, field `jwt.RegisteredClaims` (seperti `ExpiresAt`) **belum diisi** secara eksplisit. Dengan demikian, token yang dihasilkan saat ini tidak memiliki batas waktu kedaluwarsa (tanpa TTL), kecuali jika diubah di masa mendatang.

---

## 5. Middleware Pipeline

### 5.1. Authentication Middleware (`JWTMiddleware`)
File: [internal/shared/middleware/jwt_middleware.go](file:///home/wenell/projects/lavendera-API/internal/shared/middleware/jwt_middleware.go)

Urutan verifikasi:
1. Membaca header `Authorization`. Jika kosong -> abort `401 Unauthorized` (`"Header otorisasi tidak ditemukan"`).
2. Memeriksa format Bearer: `len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer")`. Jika salah -> abort `401 Unauthorized` (`"Format token salah (gunakan Bearer <token>)"`).
3. Mem-parse token dengan claims `utils.JWTClaims`:
   - Validasi algoritma wajib `HS256`. Jika metode tanda tangan berbeda -> ditolak (`jwt.ErrSignatureInvalid`).
   - Verifikasi tanda tangan dengan `jwtConfig.SecretKey`.
4. Memvalidasi kelengkapan claims payload:
   - `UserID` tidak boleh kosong.
   - `TenantID` tidak boleh kosong.
   - `Role` tidak boleh kosong.
5. Mem-parse `UserID` dan `TenantID` ke tipe `uuid.UUID`. Jika format UUID tidak valid -> abort `401 Unauthorized`.
6. **Context Propagation**: Menyimpan identitas ke context HTTP request:
   ```go
   ctx := c.Request.Context()
   ctx = appcontext.WithUserID(ctx, userID)
   ctx = appcontext.WithTenantID(ctx, tenantID)
   ctx = appcontext.WithRole(ctx, claims.Role)
   c.Request = c.Request.WithContext(ctx)
   c.Next()
   ```

### 5.2. Role-Based Access Control (`RequireRole`)
File: [internal/shared/middleware/rbac_middleware.go](file:///home/wenell/projects/lavendera-API/internal/shared/middleware/rbac_middleware.go)

1. Mengambil role dari context via `appcontext.RoleFromContext(c.Request.Context())`.
2. Jika role tidak ada atau kosong -> abort `403 Forbidden` (`"User role not found"`).
3. Memeriksa apakah role user ada dalam parameter variadik `roles ...string`:
   - Jika cocok -> memanggil `c.Next()`.
   - Jika tidak ada satupun yang cocok -> abort `403 Forbidden` (`"You do not have access to this feature"`).

---

## 6. Request Context Storage (`appcontext`)

File: [internal/shared/appcontext/context.go](file:///home/wenell/projects/lavendera-API/internal/shared/appcontext/context.go)

Kunci context dikelola menggunakan unexported/typed key `contextKey string`:
```go
const (
    UserIDKey   contextKey = "user_id"
    TenantIDKey contextKey = "tenant_id"
    roleKey     contextKey = "role"
)
```

Fungsi helper yang tersedia:
- `WithUserID(ctx context.Context, userID uuid.UUID) context.Context`
- `WithTenantID(ctx context.Context, tenantID uuid.UUID) context.Context`
- `WithRole(ctx context.Context, role string) context.Context`
- `UserIDFromContext(ctx context.Context) (uuid.UUID, bool)`
- `TenantIDFromContext(ctx context.Context) (uuid.UUID, bool)`
- `RoleFromContext(ctx context.Context) (string, bool)`

---

## 7. Multi-Tenant Isolation Enforcement

Isolasi data diberlakukan ketat pada tingkat akses database melalui GORM Scope `database.TenantScope(ctx)`:
1. Setiap pemanggilan query repository memanggil:
   ```go
   s.DB.WithContext(ctx).Scopes(database.TenantScope(ctx))
   ```
2. Helper `TenantScope`:
   - Membaca `TenantID` dari context menggunakan `appcontext.TenantIDFromContext(ctx)`.
   - Jika tenant ID tidak ditemukan atau bernilai `uuid.Nil`: menyuntikkan query `db.Where("1 = 0")`. Hal ini menjamin tidak ada row yang bisa dibaca/ditulis jika user tidak memiliki konteks tenant yang sah.
   - Jika valid: menyuntikkan query `db.Where("tenant_id = ?", tenantID)`.
3. Pada operasi Create, Service layer secara eksplisit menetapkan `entity.TenantID = tenantID` yang diambil dari context, bukan mempercayai input dari payload request pengguna.

---

## 8. Public vs Protected Routes Map

- **Public Routes**:
  - `GET /`
  - `POST /api/v1/auth/register`
  - `POST /api/v1/auth/login`
- **Protected Routes (Harus Berisi Bearer Token + Role ADMIN)**:
  - `POST /api/v1/outlets`
  - `GET /api/v1/outlets`
  - `GET /api/v1/outlets/:id`
  - `PATCH /api/v1/outlets/:id`
  - `DELETE /api/v1/outlets/:id`
  - `GET /api/v1/service-categories`
  - `GET /api/v1/service-categories/:id`
  - `POST /api/v1/service-categories`
  - `PATCH /api/v1/service-categories/:id`
  - `DELETE /api/v1/service-categories/:id`
  - `POST /api/v1/services`
  - `GET /api/v1/services`
  - `GET /api/v1/services/:id`
  - `PATCH /api/v1/services/:id`
  - `DELETE /api/v1/services/:id`
  - `POST /api/v1/discounts`
  - `GET /api/v1/discounts`
  - `GET /api/v1/discounts/:id`
  - `PATCH /api/v1/discounts/:id`
  - `DELETE /api/v1/discounts/:id`
- **Staff Routes**:
  - Disediakan placeholder komentar pada [routes.go](file:///home/wenell/projects/lavendera-API/internal/app/routes.go#L55): `// route khusus staff`. Belum ada rute yang didaftarkan khusus STAFF.
