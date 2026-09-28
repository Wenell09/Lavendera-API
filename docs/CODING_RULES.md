# Coding Rules & Conventions — Lavendera API

> **Status**: Diidentifikasi langsung dari kode implementasi aktual saat ini.  
> **Source of Truth**: Seluruh file pada repositori `lavendera-api`.

Dokumen ini mencatat konvensi dan pola penulisan kode (*coding conventions*) yang **benar-benar digunakan** dalam project ini. Setiap AI Agent atau pengembang yang memodifikasi atau menambahkan kode baru diharapkan mengikuti aturan-aturan berikut agar tetap konsisten dengan basis kode yang sudah ada.

---

## 1. Naming Conventions

### 1.1. Package Naming
- Menggunakan huruf kecil murni (*lowercase*), berupa kata tunggal atau pemisah *underscore* (`snake_case`) jika terdiri dari dua kata.
- Contoh aktual:
  - `package service_category`
  - `package appcontext`
  - `package applogger`
  - `package appvalidator`
  - `package apperror`
  - `package models`

### 1.2. File Naming
- Menggunakan `snake_case` di seluruh project.
- Pola file yang memisahkan interface dan implementasi:
  - `<domain>_<layer>.go` untuk deklarasi interface (contoh: `service_repository.go`, `outlet_service.go`, `discount_controller.go`).
  - `<domain>_<layer>_impl.go` untuk implementasi konkret struct (contoh: `service_repository_impl.go`, `outlet_service_impl.go`, `discount_controller_impl.go`).
- Pola DTO:
  - `<domain>_request.go`
  - `<domain>_response.go`
  - `<domain>_filter.go`
- Pola Routes:
  - `<domain>_routes.go`
- Pola Wire:
  - `wire.go` (pada masing-masing modul)

### 1.3. Interface Naming
- Interface dinamai dengan PascalCase dan mencerminkan perannya:
  - `AuthRepository`, `AuthService`, `AuthController`
  - `ServiceCategoryRepository`, `ServiceCategoryService`, `ServiceCategoryController`
  - `ServiceRepository`, `Service`, `ServiceController`
  - `DiscountRepository`, `DiscountService`, `DiscountController`
  - `OutletRepository`, `OutletService`, `OutletController`

### 1.4. Implementation Struct Naming
- Struct implementasi menggunakan akhiran `Impl`:
  - `AuthRepositoryImpl`, `AuthServiceImpl`, `AuthControllerImpl`
  - `ServiceCategoryRepositoryImpl`, `ServiceCategoryServiceImpl`, `ServiceCategoryControllerImpl`
  - `ServiceRepositoryImpl`, `ServiceImpl`, `ServiceControllerImpl`
  - `DiscountRepositoryImpl`, `DiscountServiceImpl`, `DiscountControllerImpl`
  - `OutletRepositoryImpl`, `OutletServiceImpl`, `OutletControllerImpl`

---

## 2. Constructor & Dependency Injection Pattern

### 2.1. Constructor Pattern
- Setiap komponen diekspos melalui fungsi konstruktor publik berawalan `New...` yang mengembalikan tipe interface (bukan struct pointer konkret).
- Contoh:
  ```go
  func NewOutletRepository(db *gorm.DB) OutletRepository {
      return &OutletRepositoryImpl{DB: db}
  }

  func NewOutletService(repository repository.OutletRepository, logger *logrus.Logger) OutletService {
      return &OutletServiceImpl{Repository: repository, Logger: logger}
  }

  func NewOutletController(service service.OutletService, validator *validator.Validate) OutletController {
      return &OutletControllerImpl{Service: service, Validator: validator}
  }
  ```

### 2.2. Wire Sets
- Setiap modul menyediakan file `wire.go` berisi variabel `WireSet`:
  ```go
  var WireSet = wire.NewSet(
      repository.NewOutletRepository,
      service.NewOutletService,
      controller.NewOutletController,
  )
  ```
- Di tingkat aplikasi, seluruh set diimpor ke dalam `wire.Build(...)` di [internal/app/wire.go](file:///home/wenell/projects/lavendera-API/internal/app/wire.go).
- `wire_gen.go` di-generate menggunakan perintah `wire ./internal/app` dan tidak boleh diedit manual.

---

## 3. Context Usage & Tenant Scope

1. **Context First Parameter**:
   - Seluruh method pada Repository dan Service wajib menerima `ctx context.Context` sebagai parameter pertama.
   - Controller mengambil context dari Gin menggunakan `c.Request.Context()`.
2. **GORM Context Propagation**:
   - Setiap query GORM wajib memanggil `.WithContext(ctx)`.
3. **Tenant Scoping**:
   - Setiap query data per-tenant wajib memanggil `.Scopes(database.TenantScope(ctx))`.
   - Pola query standar:
     ```go
     s.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Entity{})...
     ```
4. **Context Identity Accessors**:
   - Pengambilan identitas dari context wajib menggunakan helper [appcontext](file:///home/wenell/projects/lavendera-API/internal/shared/appcontext/context.go):
     - `appcontext.TenantIDFromContext(ctx)`
     - `appcontext.UserIDFromContext(ctx)`
     - `appcontext.RoleFromContext(ctx)`

---

## 4. DTO & Partial Update Pattern

### 4.1. Request DTOs
- Struct request untuk Create menggunakan value biasa dengan tag `validate:"required,..."`:
  ```go
  type CreateDiscountRequest struct {
      Name     string `json:"name" validate:"required,min=3,max=255"`
      Type     string `json:"type" validate:"required,oneof=percentage fixed"`
      Value    int64  `json:"value" validate:"required,gt=0"`
      IsActive bool   `json:"is_active"`
  }
  ```
- Struct request untuk Update menggunakan tipe pointer (`*string`, `*int64`, `*bool`) dengan tag `validate:"omitempty,..."`:
  ```go
  type UpdateDiscountRequest struct {
      Name     *string `json:"name" validate:"omitempty,min=3,max=255"`
      Type     *string `json:"type" validate:"omitempty,oneof=percentage fixed"`
      Value    *int64  `json:"value" validate:"omitempty,min=0"`
      IsActive *bool   `json:"is_active"`
  }
  ```

### 4.2. Partial Updates via Map
- Pada Service layer, update parsial dirangkai ke dalam `map[string]interface{}`:
  ```go
  updateData := make(map[string]interface{})
  if req.Name != nil {
      updateData["name"] = *req.Name
  }
  if req.IsActive != nil {
      updateData["is_active"] = *req.IsActive
  }
  if len(updateData) == 0 {
      return currentEntity, nil
  }
  if err := s.Repository.Update(ctx, id, updateData); err != nil {
      return nil, err
  }
  ```
- Repository mengeksekusi:
  ```go
  if len(data) == 0 {
      return nil
  }
  return r.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Entity{}).Where("id = ?", id).Updates(data).Error
  ```

---

## 5. Validation Pattern

Project saat ini menggunakan kombinasi **Go Playground Validator (v10)** dan validasi manual di Service:
1. **Controller Layer**:
   ```go
   var req dto.CreateServiceRequest
   if err := c.ShouldBindJSON(&req); err != nil {
       apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid request body"})
       return
   }
   if err := s.Validator.Struct(req); err != nil {
       if ve, ok := err.(validator.ValidationErrors); ok {
           apperror.NewHandleError(c, apperror.NewFieldError(ve))
           return
       }
       apperror.NewHandleError(c, err)
       return
   }
   ```
2. **Custom Validator**:
   - `id_phone`: Memvalidasi nomor telepon region Indonesia (`ID`) menggunakan pustaka `nyaruka/phonenumbers`. Didaftarkan di [validator.go](file:///home/wenell/projects/lavendera-API/internal/shared/appvalidator/validator.go).
3. **URL Param Validation**:
   - Parameter UUID di URL divalidasi via `uuid.Parse(c.Param("id"))`. Jika gagal -> `apperror.ValidationError{Msg: "invalid ... id"}`.
4. **Service-level Validation**:
   - Validasi keunikan (`ExistsByName`).
   - Validasi aturan bisnis spesifik (misal membatasi diskon persentase maks 100).

---

## 6. Response & Error Handling Patterns

### 6.1. Standard Response Envelope
Semua endpoint HTTP mengembalikan response ber-envelope seragam:
```go
// Success
c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "message", data, response.ResponseMeta{}))

// Success with Pagination
c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "message", result.Data, response.ResponseMeta{Pagination: result.Pagination}))

// Error
apperror.NewHandleError(c, err)
```

### 6.2. Error Types in `apperror`
Semua custom error mengimplementasikan interface:
```go
type AppError interface {
    error
    StatusCode() int
    ResponseMessage() string
    ErrorData() interface{}
}
```
Error yang telah diimplementasikan:
- `apperror.ValidationError`: HTTP 400 (`"Validation Error"`, message string).
- `apperror.FieldError`: HTTP 400 (`"Validation Error"`, array `FieldDetail{Field, Error}`).
- `apperror.NotFoundError`: HTTP 404 (`"Not Found"`, message string).
- `apperror.UnauthorizedError`: HTTP 401 (`"Unauthorized"`, message string).
- `apperror.ConflictError`: HTTP 409 (`"Conflict"`, message string).
- Error standar yang tidak dikenali otomatis diterjemahkan menjadi HTTP 500 (`"Internal Server Error"`).

---

## 7. Pagination & Search Patterns

Pola pagination dan pencarian seragam pada modul `service`, `outlet`, dan `discount`:
1. **Filter Struct**:
   ```go
   type FeatureFilter struct {
       Search string `json:"search"`
       Page   int    `json:"page"`
       Limit  int    `json:"limit"`
   }

   func (f *FeatureFilter) SetDefault() {
       if f.Page < 1 { f.Page = 1 }
       if f.Limit < 1 { f.Limit = 10 }
       if f.Limit > 100 { f.Limit = 100 }
   }
   ```
2. **Controller Query Parsing**:
   ```go
   search := c.Query("search")
   page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
   limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
   ```
3. **Repository Execution**:
   - Hitung total: `query.Count(&total)`.
   - Hitung offset: `offset := (filter.Page - 1) * filter.Limit`.
   - Eksekusi query: `query.Order("created_at DESC").Offset(offset).Limit(filter.Limit).Find(&data)`.
4. **Service Total Pages**:
   - `totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))`.

---

## 8. Logging Pattern

Project menggunakan **Sirupsen Logrus** dengan format JSON ke `stdout`:
1. **Middleware Logging**:
   [LoggerMiddleware](file:///home/wenell/projects/lavendera-API/internal/shared/middleware/logger_middleware.go) mencatat setiap request dengan field:
   `status`, `latency_ms`, `ip`, `method`, `path`.
2. **Service Logging**:
   Service layer menyertakan context tenant melalui `utils.LogWithContext(s.Logger, ctx)`:
   ```go
   logger := utils.LogWithContext(s.Logger, ctx).WithField("name", req.Name)
   logger.Info("creating new entity")
   ```
   Jika `tenant_id` ada di context, field `"tenant_id"` otomatis disisipkan ke log entry.

---

## 9. Transaction Pattern

Jika suatu operasi memodifikasi lebih dari satu tabel, gunakan `DB.WithContext(ctx).Transaction`:
```go
return a.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(parent).Error; err != nil {
        return err
    }
    child.ParentID = parent.ID
    if err := tx.Create(child).Error; err != nil {
        return err
    }
    return nil
})
```
