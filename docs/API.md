# API Documentation — Lavendera API

> **Status**: Verifikasi dari implementasi aktual kode.  
> **Source of Truth**: [internal/app/routes.go](file:///home/wenell/projects/lavendera-API/internal/app/routes.go) dan controller masing-masing modul.

---

## Daftar Lengkap Endpoint

Semua endpoint API (kecuali root `/`) berada di bawah base path `/api/v1`.

| No | Method | Path | Auth | Role | Deskripsi Singkat |
|---|---|---|---|---|---|
| 1 | `GET` | `/` | Public | None | API Gateway Health Check |
| 2 | `POST` | `/api/v1/auth/register` | Public | None | Pendaftaran Tenant & Admin Default |
| 3 | `POST` | `/api/v1/auth/login` | Public | None | Login User & Penerbitan JWT Token |
| 4 | `POST` | `/api/v1/outlets` | Bearer JWT | `ADMIN` | Membuat Outlet Baru |
| 5 | `GET` | `/api/v1/outlets` | Bearer JWT | `ADMIN` | Mengambil Daftar Outlet (Paginated + Search) |
| 6 | `GET` | `/api/v1/outlets/:id` | Bearer JWT | `ADMIN` | Mengambil Detail Outlet Berdasarkan ID |
| 7 | `PATCH` | `/api/v1/outlets/:id` | Bearer JWT | `ADMIN` | Mengubah Data Outlet Secara Parsial |
| 8 | `DELETE` | `/api/v1/outlets/:id` | Bearer JWT | `ADMIN` | Menghapus Outlet Berdasarkan ID |
| 9 | `GET` | `/api/v1/service-categories` | Bearer JWT | `ADMIN` | Mengambil Seluruh Kategori Layanan |
| 10 | `GET` | `/api/v1/service-categories/:id` | Bearer JWT | `ADMIN` | Mengambil Kategori Layanan Berdasarkan ID |
| 11 | `POST` | `/api/v1/service-categories` | Bearer JWT | `ADMIN` | Membuat Kategori Layanan Baru |
| 12 | `PATCH` | `/api/v1/service-categories/:id` | Bearer JWT | `ADMIN` | Mengubah Nama Kategori Layanan |
| 13 | `DELETE` | `/api/v1/service-categories/:id` | Bearer JWT | `ADMIN` | Menghapus Kategori Layanan Berdasarkan ID |
| 14 | `POST` | `/api/v1/services` | Bearer JWT | `ADMIN` | Membuat Layanan Laundry Baru |
| 15 | `GET` | `/api/v1/services` | Bearer JWT | `ADMIN` | Mengambil Daftar Layanan (Paginated + Filter) |
| 16 | `GET` | `/api/v1/services/:id` | Bearer JWT | `ADMIN` | Mengambil Detail Layanan Berdasarkan ID |
| 17 | `PATCH` | `/api/v1/services/:id` | Bearer JWT | `ADMIN` | Mengubah Data Layanan Secara Parsial |
| 18 | `DELETE` | `/api/v1/services/:id` | Bearer JWT | `ADMIN` | Menghapus Layanan Berdasarkan ID |
| 19 | `POST` | `/api/v1/discounts` | Bearer JWT | `ADMIN` | Membuat Diskon Promosi Baru |
| 20 | `GET` | `/api/v1/discounts` | Bearer JWT | `ADMIN` | Mengambil Daftar Diskon (Paginated + Search) |
| 21 | `GET` | `/api/v1/discounts/:id` | Bearer JWT | `ADMIN` | Mengambil Detail Diskon Berdasarkan ID |
| 22 | `PATCH` | `/api/v1/discounts/:id` | Bearer JWT | `ADMIN` | Mengubah Data Diskon Secara Parsial |
| 23 | `DELETE` | `/api/v1/discounts/:id` | Bearer JWT | `ADMIN` | Menghapus Diskon Berdasarkan ID |

---

## 1. System / Gateway

### GET /
- **Auth Required**: No (Public)
- **Role Requirement**: None
- **Parameters**: None
- **Request Body**: None
- **Response** (`200 OK`):
  ```json
  {
    "status": "success",
    "message": "Lavendera API Gateway is running 🚀"
  }
  ```

---

## 2. Authentication Module

### POST /api/v1/auth/register
- **Auth Required**: No (Public)
- **Role Requirement**: None
- **Request Body**:
  ```json
  {
    "tenant_name": "Laundry Sejahtera",
    "name": "Budi Santoso",
    "email": "budi@laundrysejahtera.com",
    "password": "password123"
  }
  ```
- **Validation Rules**:
  - `tenant_name`: required, min 3, max 100
  - `name`: required, min 2, max 100
  - `email`: required, valid email, max 255
  - `password`: required, min 8, max 72
- **Response** (`201 Created`):
  ```json
  {
    "status": 201,
    "message": "Registration successful",
    "success": true,
    "data": {
      "tenant": {
        "id": "1e5db5d2-0056-4c4f-9e79-5ee45b98df99",
        "name": "Laundry Sejahtera",
        "slug": "laundry-sejahtera",
        "email": "budi@laundrysejahtera.com"
      },
      "user": {
        "id": "76d54cf8-a1fb-45ba-b072-467f53f93ce0",
        "tenant_id": "1e5db5d2-0056-4c4f-9e79-5ee45b98df99",
        "name": "Budi Santoso",
        "email": "budi@laundrysejahtera.com",
        "role": "ADMIN"
      }
    },
    "meta": {}
  }
  ```
- **Errors**:
  - `400 Bad Request`: Body JSON tidak valid atau validasi field gagal.
  - `409 Conflict`: `"email already registered"` atau `"tenant slug already exists"`.
  - `500 Internal Server Error`: Kegagalan transaksi database atau hashing password.

---

### POST /api/v1/auth/login
- **Auth Required**: No (Public)
- **Role Requirement**: None
- **Request Body**:
  ```json
  {
    "email": "budi@laundrysejahtera.com",
    "password": "password123"
  }
  ```
- **Validation Rules**:
  - `email`: required, email, max 255
  - `password`: required, min 8, max 72
- **Response** (`200 OK`):
  ```json
  {
    "status": 200,
    "message": "Login successful",
    "success": true,
    "data": {
      "user": {
        "id": "76d54cf8-a1fb-45ba-b072-467f53f93ce0",
        "tenant_id": "1e5db5d2-0056-4c4f-9e79-5ee45b98df99",
        "name": "Budi Santoso",
        "email": "budi@laundrysejahtera.com",
        "role": "ADMIN"
      },
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    },
    "meta": {}
  }
  ```
- **Errors**:
  - `400 Bad Request`: Body JSON tidak valid atau format salah.
  - `401 Unauthorized`: `"invalid email or password"` atau `"user account is inactive"`.
  - `500 Internal Server Error`: Kesalahan server saat generate JWT token.

---

## 3. Outlet Module

### POST /api/v1/outlets
- **Auth Required**: Yes (`Bearer <token>`)
- **Role Requirement**: `ADMIN`
- **Request Body**:
  ```json
  {
    "name": "Outlet Cabang Melati",
    "phone": "081234567890",
    "address": "Jl. Melati No. 12, Jakarta",
    "is_public_order_enabled": true,
    "is_active": true
  }
  ```
- **Validation Rules**:
  - `name`: required, min 5, max 255
  - `phone`: required, format nomor HP Indonesia (`id_phone`)
  - `address`: required
  - `is_public_order_enabled`: boolean
  - `is_active`: boolean
- **Special Behavior**: Slug di-generate otomatis dari `name` menggunakan `utils.GenerateSlug`.
- **Response** (`201 Created`):
  ```json
  {
    "status": 201,
    "message": "outlet created successfully",
    "success": true,
    "data": {
      "id": "e674b9e2-3cf7-4fbc-b4e8-4680dc1c1f7b",
      "name": "Outlet Cabang Melati",
      "slug": "outlet-cabang-melati",
      "phone": "081234567890",
      "address": "Jl. Melati No. 12, Jakarta",
      "is_public_order_enabled": true,
      "is_active": true,
      "created_at": "2026-09-28T10:00:00Z",
      "updated_at": "2026-09-28T10:00:00Z"
    },
    "meta": {}
  }
  ```
- **Errors**: `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `409 Conflict` (`"outlet name already exists"`).

---

### GET /api/v1/outlets
- **Auth Required**: Yes (`Bearer <token>`)
- **Role Requirement**: `ADMIN`
- **Query Parameters**:
  - `search` (string, opsional): filter nama awalan (`name ILIKE search + "%"`).
  - `page` (int, opsional, default: `1`).
  - `limit` (int, opsional, default: `10`, max: `100`).
- **Pagination**: Yes (`meta.pagination`).
- **Response** (`200 OK`):
  ```json
  {
    "status": 200,
    "message": "outlets retrieved successfully",
    "success": true,
    "data": [
      {
        "id": "e674b9e2-3cf7-4fbc-b4e8-4680dc1c1f7b",
        "name": "Outlet Cabang Melati",
        "slug": "outlet-cabang-melati",
        "phone": "081234567890",
        "address": "Jl. Melati No. 12, Jakarta",
        "is_public_order_enabled": true,
        "is_active": true,
        "created_at": "2026-09-28T10:00:00Z",
        "updated_at": "2026-09-28T10:00:00Z"
      }
    ],
    "meta": {
      "pagination": {
        "page": 1,
        "limit": 10,
        "total": 1,
        "total_pages": 1
      }
    }
  }
  ```

---

### GET /api/v1/outlets/:id
- **Auth Required**: Yes (`Bearer <token>`)
- **Role Requirement**: `ADMIN`
- **Path Parameters**: `id` (UUID format)
- **Response** (`200 OK`): Single `OutletResponse` object.
- **Errors**: `400 Bad Request` (`"invalid outlet id"`), `404 Not Found` (`"outlet not found"`).

---

### PATCH /api/v1/outlets/:id
- **Auth Required**: Yes (`Bearer <token>`)
- **Role Requirement**: `ADMIN`
- **Path Parameters**: `id` (UUID format)
- **Request Body** (Semua field opsional / partial update):
  ```json
  {
    "name": "Outlet Melati Utama",
    "slug": "outlet-melati-utama",
    "phone": "081298765432",
    "address": "Jl. Melati Raya No. 15",
    "is_public_order_enabled": false,
    "is_active": true
  }
  ```
- **Response** (`200 OK`): Single `OutletResponse` object yang sudah diperbarui.
- **Errors**: `400 Bad Request`, `404 Not Found`, `409 Conflict`.

---

### DELETE /api/v1/outlets/:id
- **Auth Required**: Yes (`Bearer <token>`)
- **Role Requirement**: `ADMIN`
- **Path Parameters**: `id` (UUID format)
- **Response** (`200 OK`):
  ```json
  {
    "status": 200,
    "message": "outlet deleted successfully",
    "success": true,
    "data": null,
    "meta": {}
  }
  ```
- **Errors**: `400 Bad Request`, `404 Not Found`.

---

## 4. Service Category Module

### GET /api/v1/service-categories
- **Auth Required**: Yes (`Bearer <token>`)
- **Role Requirement**: `ADMIN`
- **Pagination / Search**: No (Mengambil semua kategori milik tenant urut `created_at DESC`).
- **Response** (`200 OK`):
  ```json
  {
    "status": 200,
    "message": "service categories retrieved successfully",
    "success": true,
    "data": [
      {
        "id": "9059f13e-63f6-4927-a068-07e4d8fb8756",
        "name": "Kiloan",
        "created_at": "2026-09-28T09:00:00Z",
        "updated_at": "2026-09-28T09:00:00Z"
      }
    ],
    "meta": {}
  }
  ```

---

### GET /api/v1/service-categories/:id
- **Auth Required**: Yes (`Bearer <token>`)
- **Role Requirement**: `ADMIN`
- **Path Parameters**: `id` (UUID)
- **Response** (`200 OK`): Single `ServiceCategoryResponse`.
- **Errors**: `400 Bad Request` (`"invalid service id"`), `404 Not Found` (`"service category not found"`).

---

### POST /api/v1/service-categories
- **Auth Required**: Yes (`Bearer <token>`)
- **Role Requirement**: `ADMIN`
- **Request Body**:
  ```json
  {
    "name": "Dry Clean"
  }
  ```
- **Validation Rules**: `name`: required, min 2, max 100.
- **Response** (`201 Created`): Single `ServiceCategoryResponse`.
- **Errors**: `400 Bad Request`, `409 Conflict` (`"service category name already exists"`).

---

### PATCH /api/v1/service-categories/:id
- **Auth Required**: Yes (`Bearer <token>`)
- **Role Requirement**: `ADMIN`
- **Path Parameters**: `id` (UUID)
- **Request Body**:
  ```json
  {
    "name": "Premium Dry Clean"
  }
  ```
- **Response** (`200 OK`): Single `ServiceCategoryResponse`.
- **Errors**: `400 Bad Request`, `404 Not Found`, `409 Conflict`.

---

### DELETE /api/v1/service-categories/:id
- **Auth Required**: Yes (`Bearer <token>`)
- **Role Requirement**: `ADMIN`
- **Path Parameters**: `id` (UUID)
- **Response** (`200 OK`):
  ```json
  {
    "status": 200,
    "message": "service category deleted successfully",
    "success": true,
    "data": null,
    "meta": {}
  }
  ```
- **Errors**: `400 Bad Request`, `404 Not Found`.

---

## 5. Service Module

### POST /api/v1/services
- **Auth Required**: Yes (`Bearer <token>`)
- **Role Requirement**: `ADMIN`
- **Request Body**:
  ```json
  {
    "category_id": "9059f13e-63f6-4927-a068-07e4d8fb8756",
    "name": "Cuci Setrika Reguler",
    "price": 8000,
    "min_quantity": 3,
    "unit": "gram",
    "duration_days": 2,
    "is_active": true
  }
  ```
- **Validation Rules**:
  - `category_id`: required (UUID)
  - `name`: required, min 3, max 255
  - `price`: gte=0
  - `min_quantity`: gt=0
  - `unit`: required (database CHECK: `'gram'`, `'pcs'`, `'pair'`)
  - `duration_days`: gt=0
  - `is_active`: boolean
- **Response** (`201 Created`):
  ```json
  {
    "status": 201,
    "message": "service created successfully",
    "success": true,
    "data": {
      "id": "e44d348a-6b87-43cf-be72-c5180f96894c",
      "Category": {
        "id": "9059f13e-63f6-4927-a068-07e4d8fb8756",
        "name": "Kiloan"
      },
      "name": "Cuci Setrika Reguler",
      "price": 8000,
      "min_quantity": 3,
      "unit": "gram",
      "duration_days": 2,
      "is_active": true,
      "created_at": "2026-09-28T10:15:00Z",
      "updated_at": "2026-09-28T10:15:00Z"
    },
    "meta": {}
  }
  ```
- **Errors**: `400 Bad Request`, `409 Conflict` (`"service name already exists"`).

---

### GET /api/v1/services
- **Auth Required**: Yes (`Bearer <token>`)
- **Role Requirement**: `ADMIN`
- **Query Parameters**:
  - `search` (string, opsional): prefix search nama layanan (`name ILIKE search + "%"`).
  - `category_id` (UUID, opsional): filter berdasarkan ID kategori layanan.
  - `page` (int, opsional, default: `1`).
  - `limit` (int, opsional, default: `10`, max: `100`).
- **Pagination**: Yes (`meta.pagination`).
- **Response** (`200 OK`):
  ```json
  {
    "status": 200,
    "message": "services retrieved successfully",
    "success": true,
    "data": [
      {
        "id": "e44d348a-6b87-43cf-be72-c5180f96894c",
        "Category": {
          "id": "9059f13e-63f6-4927-a068-07e4d8fb8756",
          "name": "Kiloan"
        },
        "name": "Cuci Setrika Reguler",
        "price": 8000,
        "min_quantity": 3,
        "unit": "gram",
        "duration_days": 2,
        "is_active": true,
        "created_at": "2026-09-28T10:15:00Z",
        "updated_at": "2026-09-28T10:15:00Z"
      }
    ],
    "meta": {
      "pagination": {
        "page": 1,
        "limit": 10,
        "total": 1,
        "total_pages": 1
      }
    }
  }
  ```

---

### GET /api/v1/services/:id
- **Auth Required**: Yes (`Bearer <token>`)
- **Role Requirement**: `ADMIN`
- **Path Parameters**: `id` (UUID)
- **Response** (`200 OK`): Single `ServiceResponse` object.
- **Errors**: `400 Bad Request` (`"invalid service id"`), `404 Not Found` (`"service not found"`).

---

### PATCH /api/v1/services/:id
- **Auth Required**: Yes (`Bearer <token>`)
- **Role Requirement**: `ADMIN`
- **Path Parameters**: `id` (UUID)
- **Request Body** (Semua field opsional / partial update):
  ```json
  {
    "category_id": "9059f13e-63f6-4927-a068-07e4d8fb8756",
    "name": "Cuci Setrika Express",
    "price": 12000,
    "min_quantity": 2,
    "unit": "gram",
    "duration_days": 1,
    "is_active": true
  }
  ```
- **Validation**:
  - Jika `category_id` diberikan, service memeriksa apakah kategori tersebut ada pada tenant via `CategoryRepository.FindByID`. Jika tidak ada, mengembalikan `404 Not Found` (`"service category not found"`).
- **Response** (`200 OK`): Single `ServiceResponse` object.
- **Errors**: `400 Bad Request`, `404 Not Found`, `409 Conflict`.

---

### DELETE /api/v1/services/:id
- **Auth Required**: Yes (`Bearer <token>`)
- **Role Requirement**: `ADMIN`
- **Path Parameters**: `id` (UUID)
- **Response** (`200 OK`):
  ```json
  {
    "status": 200,
    "message": "service deleted successfully",
    "success": true,
    "data": null,
    "meta": {}
  }
  ```
- **Errors**: `400 Bad Request`, `404 Not Found`.

---

## 6. Discount Module

### POST /api/v1/discounts
- **Auth Required**: Yes (`Bearer <token>`)
- **Role Requirement**: `ADMIN`
- **Request Body**:
  ```json
  {
    "name": "Promo Kemerdekaan",
    "type": "percentage",
    "value": 15,
    "is_active": true
  }
  ```
- **Validation Rules**:
  - `name`: required, min 3, max 255
  - `type`: required, one of `percentage` atau `fixed`
  - `value`: required, gt 0
  - `is_active`: boolean
- **Special Business Rule**: Jika `type == "percentage"` dan `value > 100`, value otomatis dibatasi menjadi `100`.
- **Response** (`201 Created`):
  ```json
  {
    "status": 201,
    "message": "discount created successfully",
    "success": true,
    "data": {
      "id": "e44d348a-6b87-43cf-be72-c5180f96894c",
      "name": "Promo Kemerdekaan",
      "type": "percentage",
      "value": 15,
      "is_active": true,
      "created_at": "2026-09-28T10:30:00Z",
      "updated_at": "2026-09-28T10:30:00Z"
    },
    "meta": {}
  }
  ```
- **Errors**: `400 Bad Request`, `409 Conflict` (`"discount name already exists"`).

---

### GET /api/v1/discounts
- **Auth Required**: Yes (`Bearer <token>`)
- **Role Requirement**: `ADMIN`
- **Query Parameters**:
  - `search` (string, opsional): prefix search nama diskon.
  - `page` (int, opsional, default: `1`).
  - `limit` (int, opsional, default: `10`, max: `100`).
- **Pagination**: Yes (`meta.pagination`).
- **Response** (`200 OK`):
  ```json
  {
    "status": 200,
    "message": "discounts retrieved successfully",
    "success": true,
    "data": [
      {
        "id": "e44d348a-6b87-43cf-be72-c5180f96894c",
        "name": "Promo Kemerdekaan",
        "type": "percentage",
        "value": 15,
        "is_active": true,
        "created_at": "2026-09-28T10:30:00Z",
        "updated_at": "2026-09-28T10:30:00Z"
      }
    ],
    "meta": {
      "pagination": {
        "page": 1,
        "limit": 10,
        "total": 1,
        "total_pages": 1
      }
    }
  }
  ```

---

### GET /api/v1/discounts/:id
- **Auth Required**: Yes (`Bearer <token>`)
- **Role Requirement**: `ADMIN`
- **Path Parameters**: `id` (UUID)
- **Response** (`200 OK`): Single `DiscountResponse` object.
- **Errors**: `400 Bad Request` (`"invalid discount id"`), `404 Not Found` (`"discount not found"`).

---

### PATCH /api/v1/discounts/:id
- **Auth Required**: Yes (`Bearer <token>`)
- **Role Requirement**: `ADMIN`
- **Path Parameters**: `id` (UUID)
- **Request Body** (Semua field opsional):
  ```json
  {
    "name": "Promo Kemerdekaan 20%",
    "type": "percentage",
    "value": 20,
    "is_active": true
  }
  ```
- **Validation**: Jika tipe atau value berubah dan tipe efektif adalah `percentage` dengan `value > 100`, value dibatasi ke `100`.
- **Response** (`200 OK`): Single `DiscountResponse` object.
- **Errors**: `400 Bad Request`, `404 Not Found`, `409 Conflict`.

---

### DELETE /api/v1/discounts/:id
- **Auth Required**: Yes (`Bearer <token>`)
- **Role Requirement**: `ADMIN`
- **Path Parameters**: `id` (UUID)
- **Response** (`200 OK`):
  ```json
  {
    "status": 200,
    "message": "discount deleted successfully",
    "success": true,
    "data": null,
    "meta": {}
  }
  ```
- **Errors**: `400 Bad Request`, `404 Not Found`.
