# Workflow Documentation — Lavendera API

> **Status**: Disimpulkan langsung dari implementasi aktual kode.  
> **Source of Truth**: Pola implementasi pada [internal/service/](file:///home/wenell/projects/lavendera-API/internal/service/), [internal/outlet/](file:///home/wenell/projects/lavendera-API/internal/outlet/), dan [internal/discount/](file:///home/wenell/projects/lavendera-API/internal/discount/).

---

## 1. Feature Development Workflow

Berdasarkan pola yang konsisten pada 5 modul yang ada (`auth`, `service_category`, `service`, `discount`, `outlet`), penambahan fitur baru harus mengikuti urutan tahapan berikut:

```
Step 1: Database Migration & Model
        ↓
Step 2: DTO (Request, Response, Filter)
        ↓
Step 3: Repository (Interface & Impl)
        ↓
Step 4: Service (Interface & Impl)
        ↓
Step 5: Controller (Interface & Impl)
        ↓
Step 6: Routes Registration
        ↓
Step 7: WireSet Declaration
        ↓
Step 8: App Wiring (wire.go & routes.go)
        ↓
Step 9: Regenerate Google Wire (wire_gen.go)
```

---

## 2. Detailed CRUD Implementation Workflow

### Step 1: Migration & Model
1. Baca file migrasi SQL di [internal/db/migrations/](file:///home/wenell/projects/lavendera-API/internal/db/migrations/) dengan timestamp:
   - `<timestamp>_create_<table_name>.up.sql`
   - `<timestamp>_create_<table_name>.down.sql`


### Step 2: DTO Definition
Buat file di `internal/<feature>/dto/`:
- `<feature>_request.go`:
  - `Create<Feature>Request`: semua field wajib diberi tag `validate:"required,..."`.
  - `Update<Feature>Request`: semua field opsional berupa pointer (`*string`, `*int64`, `*bool`) dengan tag `validate:"omitempty,..."`.
- `<feature>_response.go`:
  - `<Feature>Response`: memetakan field model ke representasi JSON.
  - `<Feature>ListResponse`: membungkus slice data dan pagination meta.
- `<feature>_filter.go` (jika mendukung pagination/filter):
  - Struct `Filter` dengan method `SetDefault()` (misal `Page` min 1, `Limit` min 10 max 100).

### Step 3: Repository Layer
Buat interface di `internal/<feature>/repository/<feature>_repository.go` dan implementasi di `_impl.go`:
```go
type FeatureRepository interface {
    FindAll(ctx context.Context, filter dto.FeatureFilter) ([]models.Feature, int64, error)
    FindByID(ctx context.Context, id uuid.UUID) (*models.Feature, error)
    ExistsByName(ctx context.Context, name string) (bool, error)
    Create(ctx context.Context, item *models.Feature) error
    Update(ctx context.Context, id uuid.UUID, data map[string]interface{}) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```
**Aturan Wajib Repository**:
- Gunakan `.Scopes(database.TenantScope(ctx))` pada semua operasi `FindAll`, `FindByID`, `ExistsByName`, `Update`, dan `Delete`.
- Method `Update` menerima `map[string]interface{}` untuk partial update dan melakukan pengecekan `if len(data) == 0 { return nil }`.

### Step 4: Service Layer
Buat interface di `internal/<feature>/service/<feature>_service.go` dan implementasi di `_impl.go`:
```go
type FeatureService interface {
    FindAll(ctx context.Context, filter dto.FeatureFilter) (*dto.FeatureListResponse, error)
    FindByID(ctx context.Context, id uuid.UUID) (*dto.FeatureResponse, error)
    Create(ctx context.Context, req dto.CreateFeatureRequest) (*dto.FeatureResponse, error)
    Update(ctx context.Context, id uuid.UUID, req dto.UpdateFeatureRequest) (*dto.FeatureResponse, error)
    Delete(ctx context.Context, id uuid.UUID) error
}
```
**Aturan Wajib Service**:
1. Gunakan `utils.LogWithContext(s.Logger, ctx)` untuk inisialisasi logger berstruktur.
2. Ambil `tenantID, exists := appcontext.TenantIDFromContext(ctx)`. Jika tidak ada, kembalikan `apperror.UnauthorizedError{Msg: "tenant_id not found"}`.
3. Sebelum `Create`, periksa keunikan nama via `Repository.ExistsByName(ctx, req.Name)`. Jika ada, kembalikan `apperror.ConflictError{Msg: "... already exists"}`.
4. Tangkap error PostgreSQL duplikasi `errors.Is(err, gorm.ErrDuplicatedKey)` dan kembalikan `apperror.ConflictError`.
5. Tangkap record tidak ditemukan `errors.Is(err, gorm.ErrRecordNotFound)` dan kembalikan `apperror.NotFoundError`.
6. Pada operasi `Update`:
   - Ambil data existing via `FindByID` (validasi keberadaan).
   - Jika nama diubah, periksa apakah nama baru sudah digunakan via `ExistsByName`.
   - Bangun `map[string]interface{}` hanya untuk field yang tidak nil.
   - Panggil `Repository.Update`, lalu fetch ulang data terbaru via `FindByID`.

### Step 5: Controller Layer
Buat interface di `internal/<feature>/controller/<feature>_controller.go` dan implementasi di `_impl.go`:
```go
type FeatureController interface {
    Create(c *gin.Context)
    FindAll(c *gin.Context)
    FindByID(c *gin.Context)
    Update(c *gin.Context)
    Delete(c *gin.Context)
}
```
**Aturan Wajib Controller**:
1. Parsing body via `c.ShouldBindJSON(&req)`. Jika gagal: `apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid request body"})`.
2. Validasi via `s.Validator.Struct(req)`. Jika error:
   ```go
   if ve, ok := err.(validator.ValidationErrors); ok {
       apperror.NewHandleError(c, apperror.NewFieldError(ve))
       return
   }
   apperror.NewHandleError(c, err)
   ```
3. Parsing parameter ID dari URL:
   ```go
   id, err := uuid.Parse(c.Param("id"))
   if err != nil {
       apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid <feature> id"})
       return
   }
   ```
4. Mengembalikan response seragam:
   - Create: `c.JSON(http.StatusCreated, response.NewResponseSuccess(http.StatusCreated, "...", result, response.ResponseMeta{}))`
   - FindAll / FindByID / Update: `c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "...", result, ...))`
   - Delete: `c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "...", nil, response.ResponseMeta{}))`

### Step 6: Routes Registration
Buat di `internal/<feature>/routes/<feature>_routes.go`:
```go
func RegisterRoutes(router *gin.RouterGroup, controller controller.FeatureController) {
    group := router.Group("/<features>")
    group.POST("", controller.Create)
    group.GET("", controller.FindAll)
    group.GET("/:id", controller.FindByID)
    group.PATCH("/:id", controller.Update)
    group.DELETE("/:id", controller.Delete)
}
```

### Step 7: WireSet Declaration
Buat di `internal/<feature>/wire.go`:
```go
var WireSet = wire.NewSet(
    repository.NewFeatureRepository,
    service.NewFeatureService,
    controller.NewFeatureController,
)
```

### Step 8 & 9: Injeksi ke App & Build
1. Tambahkan `feature.WireSet` pada `wire.Build(...)` di [internal/app/wire.go](file:///home/wenell/projects/lavendera-API/internal/app/wire.go).
2. Tambahkan parameter `featureController` di `NewRouter` pada [internal/app/routes.go](file:///home/wenell/projects/lavendera-API/internal/app/routes.go) dan daftarkan rutenya di grup admin/staff.
3. Jalankan Google Wire generator di shell:
   ```bash
   wire ./internal/app
   ```
   *(File [internal/app/wire_gen.go](file:///home/wenell/projects/lavendera-API/internal/app/wire_gen.go) akan otomatis terbarui).*

---

## 3. Validation Workflow

Validasi dijalankan secara deklaratif pada Controller layer:
1. Struct request mendefinisikan tag validasi:
   - `validate:"required"`
   - `validate:"min=3,max=255"`
   - `validate:"email"`
   - `validate:"gt=0"`, `validate:"gte=0"`
   - `validate:"oneof=percentage fixed"`
   - `validate:"id_phone"` (custom validator nomor telepon Indonesia via libphonenumber)
2. Controller mengeksekusi `Validator.Struct(req)`.
3. Jika terdapat error validasi:
   - Dikonversi ke `apperror.FieldError` via helper `apperror.NewFieldError(validationErr)`.
   - Mengembalikan array object `[{"field": "...", "error": "..."}]` dengan status HTTP `400 Bad Request`.

---

## 4. Authentication Workflow

1. Request publik dikirim ke `POST /api/v1/auth/login` atau `POST /api/v1/auth/register`.
2. Pada register:
   - Cek email unik (`ExistsUserByEmail`).
   - Cek slug tenant unik (`ExistsTenantBySlug`).
   - Hash password dengan bcrypt.
   - Buat Tenant, User (Role="ADMIN"), dan kategori default ("Kiloan", "Satuan") dalam 1 database transaction.
3. Pada login:
   - Cari user berdasarkan email.
   - Verifikasi status `is_active`.
   - Bandingkan hash password (`bcrypt.CompareHashAndPassword`).
   - Generate token JWT HS256 dengan claims `user_id`, `tenant_id`, `role`.
4. Pada request berikutnya:
   - Client mengirim header `Authorization: Bearer <token>`.
   - `JWTMiddleware` memvalidasi signature token, memverifikasi claims, dan menaruh `user_id`, `tenant_id`, `role` ke dalam `c.Request.Context()`.

---

## 5. Authorization & RBAC Workflow

1. Route dibungkus ke dalam group dengan middleware `RequireRole("ADMIN")` atau role tertentu.
2. Saat request masuk:
   - Middleware memanggil `appcontext.RoleFromContext(c.Request.Context())`.
   - Jika role tidak cocok dengan daftar role yang diizinkan, request langsung di-abort dengan HTTP `403 Forbidden` dan pesan `"You do not have access to this feature"`.

---

## 6. Database Operation Workflow & Scoping

1. Seluruh query database harus melewati GORM instance yang telah disuntik context: `s.DB.WithContext(ctx)`.
2. Scope multi-tenant dipasang pada setiap query tenant data:
   ```go
   query := s.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Feature{})
   ```
3. Jika `TenantID` hilang dari context, `TenantScope` menghasilkan clause `WHERE 1 = 0` sehingga data dari tenant lain tidak akan pernah bocor atau termodifikasi.

---

## 7. Testing Status

> **Periksa Kode**: Tidak ditemukan file test (`*_test.go`) pada seluruh codebase project saat ini.
> Status pengujian otomatis: **NOT IMPLEMENTED**.
