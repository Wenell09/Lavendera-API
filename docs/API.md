# API Documentation — Lavendera API

> **Base URL**: `http://localhost:8000`  
> **API Version**: `/api/v1`  
> **Content-Type**: `application/json`

---

## Daftar Isi

1. [Konvensi Response](#konvensi-response)
2. [Autentikasi](#autentikasi)
3. [Daftar Lengkap Endpoint](#daftar-lengkap-endpoint)
4. [System / Gateway](#1-system--gateway)
5. [Authentication Module](#2-authentication-module)
6. [Outlet Module](#3-outlet-module)
7. [Service Category Module](#4-service-category-module)
8. [Service Module](#5-service-module)
9. [Discount Module](#6-discount-module)
10. [User Module](#7-user-module)
11. [Customer Module](#8-customer-module)
12. [Tenant Module](#9-tenant-module)
13. [Outlet Payment Methods Module](#10-outlet-payment-methods-module)
14. [Outlet Users Module](#11-outlet-users-module)
15. [Order Module](#12-order-module)
16. [Report Module](#13-report-module)
17. [Profile Module](#14-profile-module)

---

## Konvensi Response

Semua respons API menggunakan format JSON yang konsisten.

### Success Response
```json
{
  "status": 200,
  "message": "deskripsi pesan sukses",
  "success": true,
  "data": { ... },
  "meta": {
    "pagination": { ... }
  }
}
```

### Error Response
```json
{
  "status": 400,
  "message": "deskripsi pesan error",
  "success": false,
  "error": "detail error atau object field errors"
}
```

### Pagination Object
Digunakan pada endpoint yang mendukung paginasi (`meta.pagination`):
```json
{
  "page": 1,
  "limit": 10,
  "total": 50,
  "total_pages": 5
}
```

### Validation Error Response
```json
{
  "status": 400,
  "message": "validation error",
  "success": false,
  "error": {
    "field_name": "error message"
  }
}
```

---

## Autentikasi

Semua endpoint yang memerlukan autentikasi menggunakan **JWT Bearer Token**.

```
Authorization: Bearer <token>
```

Token diperoleh dari endpoint login dan bersifat **HS256 JWT**. Token menyimpan:
- `user_id` — UUID user
- `tenant_id` — UUID tenant
- `role` — Role user (`ADMIN` atau `STAFF`)

### Role yang Tersedia
| Role | Deskripsi |
|---|---|
| `ADMIN` | Akses penuh ke semua fitur manajemen |
| `STAFF` | Akses terbatas (baca, transaksi kasir) |

---

## Daftar Lengkap Endpoint

| No | Method | Path | Auth | Role | Deskripsi |
|---|---|---|---|---|---|
| 1 | `GET` | `/` | Public | — | Health Check |
| 2 | `POST` | `/api/v1/auth/register` | Public | — | Registrasi Tenant & Admin Default |
| 3 | `POST` | `/api/v1/auth/login` | Public | — | Login & Penerbitan JWT Token |
| 4 | `POST` | `/api/v1/outlets` | Bearer | `ADMIN` | Membuat Outlet Baru |
| 5 | `GET` | `/api/v1/outlets` | Bearer | `ADMIN` | List Outlet (Paginated + Search) |
| 6 | `GET` | `/api/v1/outlets/:id` | Bearer | `ADMIN` | Detail Outlet |
| 7 | `PATCH` | `/api/v1/outlets/:id` | Bearer | `ADMIN` | Update Outlet (Partial) |
| 8 | `DELETE` | `/api/v1/outlets/:id` | Bearer | `ADMIN` | Hapus Outlet |
| 9 | `GET` | `/api/v1/service-categories` | Bearer | `ADMIN, STAFF` | List Kategori Layanan |
| 10 | `GET` | `/api/v1/service-categories/:id` | Bearer | `ADMIN, STAFF` | Detail Kategori Layanan |
| 11 | `POST` | `/api/v1/service-categories` | Bearer | `ADMIN` | Buat Kategori Layanan |
| 12 | `PATCH` | `/api/v1/service-categories/:id` | Bearer | `ADMIN` | Update Kategori Layanan |
| 13 | `DELETE` | `/api/v1/service-categories/:id` | Bearer | `ADMIN` | Hapus Kategori Layanan |
| 14 | `POST` | `/api/v1/services` | Bearer | `ADMIN` | Buat Layanan Laundry |
| 15 | `GET` | `/api/v1/services` | Bearer | `ADMIN, STAFF` | List Layanan (Paginated + Filter) |
| 16 | `GET` | `/api/v1/services/:id` | Bearer | `ADMIN, STAFF` | Detail Layanan |
| 17 | `PATCH` | `/api/v1/services/:id` | Bearer | `ADMIN` | Update Layanan (Partial) |
| 18 | `DELETE` | `/api/v1/services/:id` | Bearer | `ADMIN` | Hapus Layanan |
| 19 | `POST` | `/api/v1/discounts` | Bearer | `ADMIN` | Buat Diskon Promosi |
| 20 | `GET` | `/api/v1/discounts` | Bearer | `ADMIN, STAFF` | List Diskon (Paginated + Search) |
| 21 | `GET` | `/api/v1/discounts/:id` | Bearer | `ADMIN, STAFF` | Detail Diskon |
| 22 | `PATCH` | `/api/v1/discounts/:id` | Bearer | `ADMIN` | Update Diskon (Partial) |
| 23 | `DELETE` | `/api/v1/discounts/:id` | Bearer | `ADMIN` | Hapus Diskon |
| 24 | `POST` | `/api/v1/users` | Bearer | `ADMIN` | Buat User / Staff Baru |
| 25 | `GET` | `/api/v1/users` | Bearer | `ADMIN` | List User (Paginated + Search + Filter Role) |
| 26 | `GET` | `/api/v1/users/:id` | Bearer | `ADMIN` | Detail User |
| 27 | `PATCH` | `/api/v1/users/:id` | Bearer | `ADMIN` | Update User (Partial) |
| 28 | `DELETE` | `/api/v1/users/:id` | Bearer | `ADMIN` | Hapus User |
| 29 | `POST` | `/api/v1/customers` | Bearer | `ADMIN, STAFF` | Buat Customer |
| 30 | `GET` | `/api/v1/customers` | Bearer | `ADMIN, STAFF` | List Customer (Paginated + Search) |
| 31 | `GET` | `/api/v1/customers/:id` | Bearer | `ADMIN, STAFF` | Detail Customer |
| 32 | `PATCH` | `/api/v1/customers/:id` | Bearer | `ADMIN, STAFF` | Update Customer (Partial) |
| 33 | `DELETE` | `/api/v1/customers/:id` | Bearer | `ADMIN, STAFF` | Hapus Customer |
| 34 | `GET` | `/api/v1/tenants/:id` | Bearer | `ADMIN` | Detail Tenant |
| 35 | `PATCH` | `/api/v1/tenants/:id` | Bearer | `ADMIN` | Update Tenant (Partial) |
| 36 | `POST` | `/api/v1/outlets/:id/payment-methods` | Bearer | `ADMIN` | Tambah Payment Method ke Outlet |
| 37 | `GET` | `/api/v1/outlets/:id/payment-methods` | Bearer | `ADMIN, STAFF` | List Payment Methods Outlet |
| 38 | `GET` | `/api/v1/outlets/:id/payment-methods/:payment_method_id` | Bearer | `ADMIN, STAFF` | Detail Payment Method |
| 39 | `PATCH` | `/api/v1/outlets/:id/payment-methods/:payment_method_id` | Bearer | `ADMIN` | Update Payment Method |
| 40 | `DELETE` | `/api/v1/outlets/:id/payment-methods/:payment_method_id` | Bearer | `ADMIN` | Hapus Payment Method |
| 41 | `GET` | `/api/v1/outlets/:id/users` | Bearer | `ADMIN` | List User di Outlet |
| 42 | `GET` | `/api/v1/users/:id/outlets` | Bearer | `ADMIN, STAFF` | Daftar Outlet yang Di-assign ke User |
| 43 | `POST` | `/api/v1/outlets/:id/users` | Bearer | `ADMIN` | Assign User ke Outlet |
| 44 | `DELETE` | `/api/v1/outlets/:id/users/:user_id` | Bearer | `ADMIN` | Unassign User dari Outlet |
| 45 | `POST` | `/api/v1/orders` | Bearer | `ADMIN, STAFF` | Buat Order (Kasir) |
| 46 | `GET` | `/api/v1/orders` | Bearer | `ADMIN, STAFF` | List Order (Paginated + Filter) |
| 47 | `GET` | `/api/v1/orders/:id` | Bearer | `ADMIN, STAFF` | Detail Order |
| 48 | `PATCH` | `/api/v1/orders/:id/status` | Bearer | `ADMIN, STAFF` | Update Status Order |
| 49 | `PATCH` | `/api/v1/orders/:id/cancel` | Bearer | `ADMIN, STAFF` | Batalkan Order |
| 50 | `PATCH` | `/api/v1/orders/:id/payments/:payment_id/confirm` | Bearer | `ADMIN, STAFF` | Konfirmasi Pembayaran |
| 51 | `PATCH` | `/api/v1/orders/:id/cash-payment` | Bearer | `ADMIN, STAFF` | Input Jumlah Bayar Tunai |
| 52 | `POST` | `/api/v1/public/orders` | Public | — | Buat Order Mandiri (Publik) |
| 53 | `GET` | `/api/v1/public/orders/track/:token` | Public | — | Lacak Order dengan Token |
| 54 | `GET` | `/api/v1/profile` | Bearer | `ADMIN, STAFF` | Lihat Data Profile |
| 55 | `GET` | `/api/v1/dashboard` | Bearer | `ADMIN, STAFF` | Ringkasan Statistik Dashboard |
| 56 | `GET` | `/api/v1/reports/financial` | Bearer | `ADMIN, STAFF` | Laporan Keuangan dan Riwayat Transaksi |


---

## 1. System / Gateway

### GET /
- **Auth**: Public
- **Deskripsi**: Health check endpoint untuk memverifikasi server berjalan.
- **Response** `200 OK`:
  ```json
  {
    "status": "success",
    "message": "Lavendera API Gateway is running 🚀"
  }
  ```

---

## 2. Authentication Module

### POST /api/v1/auth/register
- **Auth**: Public
- **Deskripsi**: Mendaftarkan tenant baru beserta akun Admin pertama secara atomik (dalam satu transaksi database).
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
  | Field | Aturan |
  |---|---|
  | `tenant_name` | required, min=3, max=100 |
  | `name` | required, min=2, max=100 |
  | `email` | required, email, max=255 |
  | `password` | required, min=8, max=72 |

- **Response** `201 Created`:
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
        "email": "budi@laundrysejahtera.com",
        "created_at": "2026-10-16T10:00:00Z",
        "updated_at": "2026-10-16T10:00:00Z"
      },
      "user": {
        "id": "76d54cf8-a1fb-45ba-b072-467f53f93ce0",
        "tenant_id": "1e5db5d2-0056-4c4f-9e79-5ee45b98df99",
        "name": "Budi Santoso",
        "email": "budi@laundrysejahtera.com",
        "role": "ADMIN",
        "is_active": true,
        "created_at": "2026-10-16T10:00:00Z",
        "updated_at": "2026-10-16T10:00:00Z"
      }
    },
    "meta": {}
  }
  ```
- **Error Responses**:
  | Status | Kondisi |
  |---|---|
  | `400` | Body JSON tidak valid atau validasi field gagal |
  | `409` | Email sudah terdaftar atau slug tenant sudah ada |
  | `500` | Kegagalan transaksi database atau hashing password |

---

### POST /api/v1/auth/login
- **Auth**: Public
- **Deskripsi**: Login user dan mendapatkan JWT token. Hanya user dengan status `is_active = true` yang dapat login.
- **Request Body**:
  ```json
  {
    "email": "budi@laundrysejahtera.com",
    "password": "password123"
  }
  ```
- **Validation Rules**:
  | Field | Aturan |
  |---|---|
  | `email` | required, email, max=255 |
  | `password` | required, min=8, max=72 |

- **Response** `200 OK`:
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
        "role": "ADMIN",
        "is_active": true,
        "created_at": "2026-10-16T10:00:00Z",
        "updated_at": "2026-10-16T10:00:00Z"
      },
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNzZkNTRjZjgtYTFmYi00NWJhLWIwNzItNDY3ZjUzZjkzY2UwIiwidGVuYW50X2lkIjoiMWU1ZGI1ZDItMDA1Ni00YzRmLTllNzktNWVlNDViOThkZjk5Iiwicm9sZSI6IkFETUlOIn0..."
    },
    "meta": {}
  }
  ```
- **Error Responses**:
  | Status | Kondisi |
  |---|---|
  | `400` | Body JSON tidak valid atau format field salah |
  | `401` | Email atau password salah, atau akun non-aktif |
  | `500` | Kesalahan server saat generate JWT |

---

## 3. Outlet Module

### POST /api/v1/outlets
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Deskripsi**: Membuat outlet baru. Slug di-generate otomatis dari `name`.
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
  | Field | Aturan |
  |---|---|
  | `name` | required, min=5, max=255 |
  | `phone` | required, format nomor HP Indonesia |
  | `address` | required |
  | `is_public_order_enabled` | boolean, opsional |
  | `is_active` | boolean, opsional |

- **Response** `201 Created`:
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
      "created_at": "2026-10-16T10:00:00Z",
      "updated_at": "2026-10-16T10:00:00Z"
    },
    "meta": {}
  }
  ```
- **Error Responses**:
  | Status | Kondisi |
  |---|---|
  | `400` | Validasi field gagal |
  | `401` | Token tidak valid atau tidak ada |
  | `403` | Bukan ADMIN |
  | `409` | Nama outlet sudah digunakan |

---

### GET /api/v1/outlets
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Deskripsi**: Mengambil daftar outlet milik tenant dengan paginasi dan pencarian.
- **Query Parameters**:
  | Parameter | Tipe | Default | Deskripsi |
  |---|---|---|---|
  | `search` | string | — | Prefix search pada nama outlet |
  | `page` | integer | `1` | Halaman saat ini |
  | `limit` | integer | `10` | Jumlah per halaman (max: 100) |

- **Response** `200 OK`:
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
        "created_at": "2026-10-16T10:00:00Z",
        "updated_at": "2026-10-16T10:00:00Z"
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
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Path Params**: `id` (UUID)
- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "outlet retrieved successfully",
    "success": true,
    "data": {
      "id": "e674b9e2-3cf7-4fbc-b4e8-4680dc1c1f7b",
      "name": "Outlet Cabang Melati",
      "slug": "outlet-cabang-melati",
      "phone": "081234567890",
      "address": "Jl. Melati No. 12, Jakarta",
      "is_public_order_enabled": true,
      "is_active": true,
      "created_at": "2026-10-16T10:00:00Z",
      "updated_at": "2026-10-16T10:00:00Z"
    },
    "meta": {}
  }
  ```
- **Error Responses**: `400` (ID tidak valid), `404` (outlet tidak ditemukan)

---

### PATCH /api/v1/outlets/:id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Path Params**: `id` (UUID)
- **Deskripsi**: Partial update — hanya field yang dikirim yang akan diperbarui.
- **Request Body** (semua field opsional):
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
- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "outlet updated successfully",
    "success": true,
    "data": {
      "id": "e674b9e2-3cf7-4fbc-b4e8-4680dc1c1f7b",
      "name": "Outlet Melati Utama",
      "slug": "outlet-melati-utama",
      "phone": "081298765432",
      "address": "Jl. Melati Raya No. 15",
      "is_public_order_enabled": false,
      "is_active": true,
      "created_at": "2026-10-16T10:00:00Z",
      "updated_at": "2026-10-16T11:00:00Z"
    },
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `404`, `409` (nama/slug sudah ada)

---

### DELETE /api/v1/outlets/:id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Path Params**: `id` (UUID)
- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "outlet deleted successfully",
    "success": true,
    "data": null,
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `404`

---

## 4. Service Category Module

### GET /api/v1/service-categories
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Deskripsi**: Mengambil semua kategori layanan milik tenant, diurutkan berdasarkan `created_at DESC`. Tidak ada paginasi.
- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "service categories retrieved successfully",
    "success": true,
    "data": [
      {
        "id": "9059f13e-63f6-4927-a068-07e4d8fb8756",
        "name": "Kiloan",
        "created_at": "2026-10-16T09:00:00Z",
        "updated_at": "2026-10-16T09:00:00Z"
      },
      {
        "id": "b34cd9f0-0123-4abc-def0-123456789abc",
        "name": "Dry Clean",
        "created_at": "2026-10-16T09:05:00Z",
        "updated_at": "2026-10-16T09:05:00Z"
      }
    ],
    "meta": {}
  }
  ```

---

### GET /api/v1/service-categories/:id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Path Params**: `id` (UUID)
- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "service category retrieved successfully",
    "success": true,
    "data": {
      "id": "9059f13e-63f6-4927-a068-07e4d8fb8756",
      "name": "Kiloan",
      "created_at": "2026-10-16T09:00:00Z",
      "updated_at": "2026-10-16T09:00:00Z"
    },
    "meta": {}
  }
  ```
- **Error Responses**: `400` (ID tidak valid), `404` (kategori tidak ditemukan)

---

### POST /api/v1/service-categories
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Request Body**:
  ```json
  {
    "name": "Dry Clean"
  }
  ```
- **Validation Rules**:
  | Field | Aturan |
  |---|---|
  | `name` | required, min=2, max=100 |

- **Response** `201 Created`:
  ```json
  {
    "status": 201,
    "message": "service category created successfully",
    "success": true,
    "data": {
      "id": "b34cd9f0-0123-4abc-def0-123456789abc",
      "name": "Dry Clean",
      "created_at": "2026-10-16T10:00:00Z",
      "updated_at": "2026-10-16T10:00:00Z"
    },
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `409` (nama kategori sudah ada)

---

### PATCH /api/v1/service-categories/:id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Path Params**: `id` (UUID)
- **Request Body**:
  ```json
  {
    "name": "Premium Dry Clean"
  }
  ```
- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "service category updated successfully",
    "success": true,
    "data": {
      "id": "b34cd9f0-0123-4abc-def0-123456789abc",
      "name": "Premium Dry Clean",
      "created_at": "2026-10-16T10:00:00Z",
      "updated_at": "2026-10-16T11:00:00Z"
    },
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `404`, `409`

---

### DELETE /api/v1/service-categories/:id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Path Params**: `id` (UUID)
- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "service category deleted successfully",
    "success": true,
    "data": null,
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `404`

---

## 5. Service Module

### POST /api/v1/services
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Deskripsi**: Membuat layanan laundry baru yang terikat pada kategori layanan.
- **Request Body**:
  ```json
  {
    "category_id": "9059f13e-63f6-4927-a068-07e4d8fb8756",
    "name": "Cuci Setrika Reguler",
    "price": 8000,
    "min_quantity": 3,
    "unit": "kg",
    "duration_days": 2,
    "is_active": true
  }
  ```
- **Validation Rules**:
  | Field | Aturan |
  |---|---|
  | `category_id` | required, UUID |
  | `name` | required, min=3, max=255 |
  | `price` | gte=0 (dalam rupiah) |
  | `min_quantity` | gt=0 |
  | `unit` | required, one of: `kg`, `pcs`, `pair` |
  | `duration_days` | gt=0 |
  | `is_active` | boolean |

- **Response** `201 Created`:
  ```json
  {
    "status": 201,
    "message": "service created successfully",
    "success": true,
    "data": {
      "id": "e44d348a-6b87-43cf-be72-c5180f96894c",
      "category": {
        "id": "9059f13e-63f6-4927-a068-07e4d8fb8756",
        "name": "Kiloan"
      },
      "name": "Cuci Setrika Reguler",
      "price": 8000,
      "min_quantity": 3,
      "unit": "kg",
      "duration_days": 2,
      "is_active": true,
      "created_at": "2026-10-16T10:15:00Z",
      "updated_at": "2026-10-16T10:15:00Z"
    },
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `404` (kategori tidak ditemukan), `409` (nama layanan sudah ada)

---

### GET /api/v1/services
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Query Parameters**:
  | Parameter | Tipe | Default | Deskripsi |
  |---|---|---|---|
  | `search` | string | — | Prefix search pada nama layanan |
  | `category_id` | UUID | — | Filter berdasarkan ID kategori |
  | `page` | integer | `1` | Halaman saat ini |
  | `limit` | integer | `10` | Jumlah per halaman (max: 100) |

- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "services retrieved successfully",
    "success": true,
    "data": [
      {
        "id": "e44d348a-6b87-43cf-be72-c5180f96894c",
        "category": {
          "id": "9059f13e-63f6-4927-a068-07e4d8fb8756",
          "name": "Kiloan"
        },
        "name": "Cuci Setrika Reguler",
        "price": 8000,
        "min_quantity": 3,
        "unit": "kg",
        "duration_days": 2,
        "is_active": true,
        "created_at": "2026-10-16T10:15:00Z",
        "updated_at": "2026-10-16T10:15:00Z"
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
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Path Params**: `id` (UUID)
- **Response** `200 OK`: Single `ServiceResponse` object (sama seperti item di list).
- **Error Responses**: `400` (ID tidak valid), `404` (layanan tidak ditemukan)

---

### PATCH /api/v1/services/:id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Path Params**: `id` (UUID)
- **Deskripsi**: Partial update. Jika `category_id` diberikan, sistem memverifikasi kategori tersebut milik tenant.
- **Request Body** (semua field opsional):
  ```json
  {
    "category_id": "9059f13e-63f6-4927-a068-07e4d8fb8756",
    "name": "Cuci Setrika Express",
    "price": 12000,
    "min_quantity": 2,
    "unit": "kg",
    "duration_days": 1,
    "is_active": true
  }
  ```
- **Response** `200 OK`: Single `ServiceResponse` object yang sudah diperbarui.
- **Error Responses**: `400`, `404` (layanan atau kategori tidak ditemukan), `409`

---

### DELETE /api/v1/services/:id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Path Params**: `id` (UUID)
- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "service deleted successfully",
    "success": true,
    "data": null,
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `404`

---

## 6. Discount Module

### POST /api/v1/discounts
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Deskripsi**: Membuat diskon promosi. Jika `type = "percentage"` dan `value > 100`, nilai otomatis dibatasi ke `100`.
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
  | Field | Aturan |
  |---|---|
  | `name` | required, min=3, max=255 |
  | `type` | required, one of: `percentage`, `fixed` |
  | `value` | required, gt=0 |
  | `is_active` | boolean |

- **Response** `201 Created`:
  ```json
  {
    "status": 201,
    "message": "discount created successfully",
    "success": true,
    "data": {
      "id": "c12d3456-7890-abcd-ef01-234567890abc",
      "name": "Promo Kemerdekaan",
      "type": "percentage",
      "value": 15,
      "is_active": true,
      "created_at": "2026-10-16T10:30:00Z",
      "updated_at": "2026-10-16T10:30:00Z"
    },
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `409` (nama diskon sudah ada)

---

### GET /api/v1/discounts
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Query Parameters**:
  | Parameter | Tipe | Default | Deskripsi |
  |---|---|---|---|
  | `search` | string | — | Prefix search pada nama diskon |
  | `page` | integer | `1` | Halaman saat ini |
  | `limit` | integer | `10` | Jumlah per halaman (max: 100) |

- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "discounts retrieved successfully",
    "success": true,
    "data": [
      {
        "id": "c12d3456-7890-abcd-ef01-234567890abc",
        "name": "Promo Kemerdekaan",
        "type": "percentage",
        "value": 15,
        "is_active": true,
        "created_at": "2026-10-16T10:30:00Z",
        "updated_at": "2026-10-16T10:30:00Z"
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
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Path Params**: `id` (UUID)
- **Response** `200 OK`: Single `DiscountResponse` object.
- **Error Responses**: `400` (ID tidak valid), `404` (diskon tidak ditemukan)

---

### PATCH /api/v1/discounts/:id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Path Params**: `id` (UUID)
- **Request Body** (semua field opsional):
  ```json
  {
    "name": "Promo Kemerdekaan 20%",
    "type": "percentage",
    "value": 20,
    "is_active": true
  }
  ```
- **Response** `200 OK`: Single `DiscountResponse` object yang sudah diperbarui.
- **Error Responses**: `400`, `404`, `409`

---

### DELETE /api/v1/discounts/:id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Path Params**: `id` (UUID)
- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "discount deleted successfully",
    "success": true,
    "data": null,
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `404`

---

## 7. User Module

### POST /api/v1/users
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Deskripsi**: Membuat akun user/staff baru dalam tenant yang sama dengan admin yang membuat.
- **Request Body**:
  ```json
  {
    "name": "Sari Dewi",
    "email": "sari@laundrysejahtera.com",
    "password": "password123",
    "role": "STAFF",
    "is_active": true
  }
  ```
- **Validation Rules**:
  | Field | Aturan |
  |---|---|
  | `name` | required, min=2, max=150 |
  | `email` | required, email, max=255 |
  | `password` | required, min=8, max=72 |
  | `role` | required, one of: `ADMIN`, `STAFF` |
  | `is_active` | boolean, opsional (default: `true`) |

- **Response** `201 Created`:
  ```json
  {
    "status": 201,
    "message": "user created successfully",
    "success": true,
    "data": {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "tenant_id": "1e5db5d2-0056-4c4f-9e79-5ee45b98df99",
      "name": "Sari Dewi",
      "email": "sari@laundrysejahtera.com",
      "role": "STAFF",
      "is_active": true,
      "created_at": "2026-10-16T10:00:00Z",
      "updated_at": "2026-10-16T10:00:00Z"
    },
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `401`, `403`, `409` (email sudah terdaftar)

---

### GET /api/v1/users
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Query Parameters**:
  | Parameter | Tipe | Default | Deskripsi |
  |---|---|---|---|
  | `search` | string | — | Prefix search pada nama atau email user |
  | `page` | integer | `1` | Halaman saat ini |
  | `limit` | integer | `10` | Jumlah per halaman (max: 100) |

- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "users retrieved successfully",
    "success": true,
    "data": [
      {
        "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "tenant_id": "1e5db5d2-0056-4c4f-9e79-5ee45b98df99",
        "name": "Sari Dewi",
        "email": "sari@laundrysejahtera.com",
        "role": "STAFF",
        "is_active": true,
        "created_at": "2026-10-16T10:00:00Z",
        "updated_at": "2026-10-16T10:00:00Z"
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

### GET /api/v1/users/:id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Path Params**: `id` (UUID)
- **Response** `200 OK`: Single `UserResponse` object.
- **Error Responses**: `400` (ID tidak valid), `404` (user tidak ditemukan)

---

### PATCH /api/v1/users/:id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Path Params**: `id` (UUID)
- **Deskripsi**: Partial update. Jika `email` diubah, keunikannya per tenant akan diperiksa. Jika `password` diberikan, akan di-hash ulang.
- **Request Body** (semua field opsional):
  ```json
  {
    "name": "Sari Dewi Utami",
    "email": "sari.utami@laundrysejahtera.com",
    "password": "newpassword123",
    "role": "ADMIN",
    "is_active": false
  }
  ```
- **Response** `200 OK`: Single `UserResponse` object yang sudah diperbarui.
- **Error Responses**: `400`, `404`, `409` (email sudah dipakai user lain)

---

### DELETE /api/v1/users/:id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Path Params**: `id` (UUID)
- **Business Rule**: Admin tidak dapat menghapus akun miliknya sendiri.
- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "user deleted successfully",
    "success": true,
    "data": null,
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `404`, `409` (mencoba hapus akun sendiri)

---

## 8. Customer Module

### POST /api/v1/customers
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Request Body**:
  ```json
  {
    "name": "Siti Rahma",
    "phone": "081234567891",
    "address": "Jl. Mawar No. 5, Jakarta"
  }
  ```
- **Validation Rules**:
  | Field | Aturan |
  |---|---|
  | `name` | required, min=3, max=150 |
  | `phone` | required, format nomor HP Indonesia |
  | `address` | opsional, max=500 |

- **Response** `201 Created`:
  ```json
  {
    "status": 201,
    "message": "customer created successfully",
    "success": true,
    "data": {
      "id": "d4e5f6a7-b8c9-0123-def4-567890abcdef",
      "name": "Siti Rahma",
      "phone": "081234567891",
      "address": "Jl. Mawar No. 5, Jakarta",
      "created_at": "2026-10-16T10:00:00Z",
      "updated_at": "2026-10-16T10:00:00Z"
    },
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `409` (nomor telepon sudah terdaftar)

---

### GET /api/v1/customers
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Query Parameters**:
  | Parameter | Tipe | Default | Deskripsi |
  |---|---|---|---|
  | `search` | string | — | Prefix search pada nama atau nomor handphone customer |
  | `page` | integer | `1` | Halaman saat ini |
  | `limit` | integer | `10` | Jumlah per halaman (max: 100) |

- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "customers retrieved successfully",
    "success": true,
    "data": [
      {
        "id": "d4e5f6a7-b8c9-0123-def4-567890abcdef",
        "name": "Siti Rahma",
        "phone": "081234567891",
        "address": "Jl. Mawar No. 5, Jakarta",
        "created_at": "2026-10-16T10:00:00Z",
        "updated_at": "2026-10-16T10:00:00Z"
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

### GET /api/v1/customers/:id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Path Params**: `id` (UUID)
- **Response** `200 OK`: Single `CustomerResponse` object.
- **Error Responses**: `400` (ID tidak valid), `404` (customer tidak ditemukan)

---

### PATCH /api/v1/customers/:id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Path Params**: `id` (UUID)
- **Request Body** (semua field opsional):
  ```json
  {
    "name": "Siti Rahma Wati",
    "phone": "081298765432",
    "address": "Jl. Mawar Baru No. 10"
  }
  ```
- **Response** `200 OK`: Single `CustomerResponse` object yang sudah diperbarui.
- **Error Responses**: `400`, `404`

---

### DELETE /api/v1/customers/:id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Path Params**: `id` (UUID)
- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "customer deleted successfully",
    "success": true,
    "data": null,
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `404`

---

## 9. Tenant Module

### GET /api/v1/tenants/:id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Path Params**: `id` (UUID tenant)
- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "tenant retrieved successfully",
    "success": true,
    "data": {
      "id": "1e5db5d2-0056-4c4f-9e79-5ee45b98df99",
      "name": "Laundry Sejahtera",
      "slug": "laundry-sejahtera",
      "email": "budi@laundrysejahtera.com",
      "created_at": "2026-10-16T10:00:00Z",
      "updated_at": "2026-10-16T10:00:00Z"
    },
    "meta": {}
  }
  ```
- **Error Responses**: `400` (ID tidak valid), `404` (tenant tidak ditemukan)

---

### PATCH /api/v1/tenants/:id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Path Params**: `id` (UUID tenant)
- **Request Body** (semua field opsional):
  ```json
  {
    "name": "Laundry Sejahtera Baru",
    "slug": "laundry-sejahtera-baru"
  }
  ```
- **Validation Rules**:
  | Field | Aturan |
  |---|---|
  | `name` | opsional, min=3, max=100 |
  | `slug` | opsional, min=3, max=100 |

- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "tenant updated successfully",
    "success": true,
    "data": {
      "id": "1e5db5d2-0056-4c4f-9e79-5ee45b98df99",
      "name": "Laundry Sejahtera Baru",
      "slug": "laundry-sejahtera-baru",
      "email": "budi@laundrysejahtera.com",
      "created_at": "2026-10-16T10:00:00Z",
      "updated_at": "2026-10-16T11:00:00Z"
    },
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `404`, `409` (slug sudah digunakan)

---

## 10. Outlet Payment Methods Module

### POST /api/v1/outlets/:id/payment-methods
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Path Params**: `id` (UUID outlet)
- **Request Body**:
  ```json
  {
    "type": "bank_transfer",
    "provider_name": "BCA",
    "account_number": "1234567890",
    "account_name": "Laundry Sejahtera",
    "qr_image_url": null,
    "is_active": true
  }
  ```
- **Validation Rules**:
  | Field | Aturan |
  |---|---|
  | `type` | required, one of: `qris`, `bank_transfer` |
  | `provider_name` | required |
  | `account_number` | opsional (wajib untuk `bank_transfer`) |
  | `account_name` | required |
  | `qr_image_url` | opsional (URL gambar QR untuk `qris`) |
  | `is_active` | boolean, opsional |

- **Response** `201 Created`:
  ```json
  {
    "status": 201,
    "message": "outlet payment method created successfully",
    "success": true,
    "data": {
      "id": "f1a2b3c4-d5e6-7890-abcd-ef0123456789",
      "type": "bank_transfer",
      "provider_name": "BCA",
      "account_number": "1234567890",
      "account_name": "Laundry Sejahtera",
      "qr_image_url": null,
      "is_active": true,
      "created_at": "2026-10-16T10:00:00Z",
      "updated_at": "2026-10-16T10:00:00Z"
    },
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `404` (outlet tidak ditemukan)

---

### GET /api/v1/outlets/:id/payment-methods
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Path Params**: `id` (UUID outlet)
- **Query Parameters**:
  | Parameter | Tipe | Default | Deskripsi |
  |---|---|---|---|
  | `search` | string | — | Prefix search pada provider name |
  | `page` | integer | `1` | Halaman saat ini |
  | `limit` | integer | `10` | Jumlah per halaman (max: 100) |

- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "outlet payment methods retrieved successfully",
    "success": true,
    "data": [
      {
        "id": "f1a2b3c4-d5e6-7890-abcd-ef0123456789",
        "type": "bank_transfer",
        "provider_name": "BCA",
        "account_number": "1234567890",
        "account_name": "Laundry Sejahtera",
        "qr_image_url": null,
        "is_active": true,
        "created_at": "2026-10-16T10:00:00Z",
        "updated_at": "2026-10-16T10:00:00Z"
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

### GET /api/v1/outlets/:id/payment-methods/:payment_method_id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Path Params**: `id` (UUID outlet), `payment_method_id` (UUID)
- **Response** `200 OK`: Single `OutletPaymentMethodResponse` object.
- **Error Responses**: `400`, `404`

---

### PATCH /api/v1/outlets/:id/payment-methods/:payment_method_id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Path Params**: `id` (UUID outlet), `payment_method_id` (UUID)
- **Request Body** (semua field opsional):
  ```json
  {
    "type": "qris",
    "provider_name": "GoPay",
    "account_number": null,
    "account_name": "Laundry Sejahtera",
    "qr_image_url": "https://example.com/qr.png",
    "is_active": true
  }
  ```
- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "outlet payment method updated successfully",
    "success": true,
    "data": {
      "id": "f1a2b3c4-d5e6-7890-abcd-ef0123456789",
      "type": "qris",
      "provider_name": "GoPay",
      "account_number": null,
      "account_name": "Laundry Sejahtera",
      "qr_image_url": "https://example.com/qr.png",
      "is_active": true,
      "created_at": "2026-10-16T10:00:00Z",
      "updated_at": "2026-10-16T11:00:00Z"
    },
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `404`

---

### DELETE /api/v1/outlets/:id/payment-methods/:payment_method_id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Path Params**: `id` (UUID outlet), `payment_method_id` (UUID)
- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "outlet payment method deleted successfully",
    "success": true,
    "data": null,
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `404`

---

## 11. Outlet Users Module

### POST /api/v1/outlets/:id/users
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Path Params**: `id` (UUID outlet)
- **Deskripsi**: Menetapkan (assign) user/staff ke outlet tertentu.
- **Request Body**:
  ```json
  {
    "outlet_id": "e674b9e2-3cf7-4fbc-b4e8-4680dc1c1f7b",
    "user_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  }
  ```
  > **Catatan**: `outlet_id` dalam body harus sama dengan `:id` di path.

- **Response** `201 Created`:
  ```json
  {
    "status": 201,
    "message": "user assigned to outlet successfully",
    "success": true,
    "data": {
      "outlet": {
        "id": "e674b9e2-3cf7-4fbc-b4e8-4680dc1c1f7b",
        "name": "Outlet Cabang Melati",
        "slug": "outlet-cabang-melati"
      },
      "users": [
        {
        "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "name": "Sari Dewi",
        "email": "sari@laundrysejahtera.com",
        "role": "STAFF",
        "is_active": true
      }
      ]
    },
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `404` (outlet atau user tidak ditemukan), `409` (user sudah di-assign ke outlet ini)

---

### GET /api/v1/outlets/:id/users
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Path Params**: `id` (UUID outlet)
- **Query Parameters**:
  | Parameter | Tipe | Default | Deskripsi |
  |---|---|---|---|
  | `search` | string | — | Prefix search pada nama user |
  | `page` | integer | `1` | Halaman saat ini |
  | `limit` | integer | `10` | Jumlah per halaman (max: 100) |

- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "outlet users retrieved successfully",
    "success": true,
    "data": {
      "outlet": {
        "id": "e674b9e2-3cf7-4fbc-b4e8-4680dc1c1f7b",
        "name": "Outlet Cabang Melati",
        "slug": "outlet-cabang-melati"
      },
      "users": [
        {
        "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "name": "Sari Dewi",
        "email": "sari@laundrysejahtera.com",
        "role": "STAFF",
        "is_active": true
      },
      {
        "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567891",
        "name": "Budi",
        "email": "budi@laundrysejahtera.com",
        "role": "STAFF",
        "is_active": true
      },
      ]
    },
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

### DELETE /api/v1/outlets/:id/users/:user_id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN`
- **Path Params**: `id` (UUID outlet), `user_id` (UUID user)
- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "user unassigned from outlet successfully",
    "success": true,
    "data": null,
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `404`

---

### GET /api/v1/users/:id/outlets
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Path Params**: `id` (UUID user)
- **Deskripsi**: Mengambil daftar outlet yang sudah di-assign ke user tertentu. Berguna bagi staff untuk mengetahui outlet mana saja yang menjadi tanggung jawabnya.
- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "user outlet retrieved successfully",
    "success": true,
    "data": {
      "user": {
          "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
          "name": "Sari Dewi",
          "email": "sari@laundrysejahtera.com",
          "role": "STAFF",
          "is_active": true
        },
      "outlets": [
        {
          "id": "e674b9e2-3cf7-4fbc-b4e8-4680dc1c1f7b",
          "name": "Outlet Cabang Melati",
          "slug": "outlet-cabang-melati"
        },
        {
          "id": "e674b9e2-3cf7-4fbc-b4e8-4680dc1c1f7c",
          "name": "Outlet Cabang Matahari",
          "slug": "outlet-cabang-matahari"
        },
      ]
    },
    "meta": {}
  }
  ```
- **Error Responses**: `400` (ID tidak valid), `404` (user tidak ditemukan)

---

## 12. Order Module

### Order Status Flow

```
MENUNGGU_KONFIRMASI → DIPROSES → SELESAI → DIAMBIL
         ↓
      DIBATALKAN
```

| Status | Deskripsi |
|---|---|
| `MENUNGGU_KONFIRMASI` | Order baru dibuat, menunggu konfirmasi staf |
| `DIPROSES` | Order sedang dikerjakan |
| `SELESAI` | Order telah selesai dikerjakan |
| `DIAMBIL` | Order sudah diambil oleh customer |
| `DIBATALKAN` | Order dibatalkan |

### Payment Status

| Status | Deskripsi |
|---|---|
| `UNPAID` | Belum ada pembayaran |
| `PAID` | Pembayaran sudah dikonfirmasi |
| `PENDING` | Pembayaran dalam proses verifikasi |

---

### POST /api/v1/orders _(Cashier Order)_
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Deskripsi**: Membuat order melalui kasir. Customer harus sudah terdaftar di sistem. Total dihitung otomatis berdasarkan layanan, kuantitas, dan diskon.
- **Request Body**:
  ```json
  {
    "outlet_id": "e674b9e2-3cf7-4fbc-b4e8-4680dc1c1f7b",
    "customer_id": "d4e5f6a7-b8c9-0123-def4-567890abcdef",
    "discount_id": "c12d3456-7890-abcd-ef01-234567890abc",
    "payment_method": "cash",
    "paid_amount": 50000,
    "notes": "Tolong pakai pewangi lavender",
    "status": "DIPROSES",
    "items": [
      {
        "service_id": "e44d348a-6b87-43cf-be72-c5180f96894c",
        "quantity": 5000
      }
    ]
  }
  ```
- **Validation Rules**:
  | Field | Aturan |
  |---|---|
  | `outlet_id` | required, UUID |
  | `customer_id` | required, UUID |
  | `discount_id` | opsional, UUID |
  | `payment_method` | required, one of: `cash`, `bank_transfer`, `qris` |
  | `paid_amount` | opsional, gte=0 |
  | `notes` | opsional, max=1000 |
  | `status` | opsional, one of: `MENUNGGU_KONFIRMASI`, `DIPROSES` |
  | `items` | required, min=1 array item |
  | `items[].service_id` | required, UUID |
  | `items[].quantity` | required, gt=0 (dalam satuan sesuai unit layanan) |

- **Response** `201 Created`:
  ```json
  {
    "status": 201,
    "message": "order created successfully",
    "success": true,
    "data": {
      "id": "ord-uuid-0001-abcd-ef0123456789",
      "tenant_id": "1e5db5d2-0056-4c4f-9e79-5ee45b98df99",
      "order_number": "ORD-20261016-0001",
      "tracking_token": "TRK-AbCdEfGhIjKl",
      "status": "DIPROSES",
      "payment_status": "PAID",
      "subtotal": 40000,
      "discount_amount": 6000,
      "total_amount": 34000,
      "paid_amount": 0,
      "change_amount": 0,
      "outlet": {
        "id": "e674b9e2-3cf7-4fbc-b4e8-4680dc1c1f7b",
        "name": "Outlet Cabang Melati",
        "phone": "081234567890",
        "address": "Jl. Melati No. 12, Jakarta"
      },
      "customer": {
        "id": "d4e5f6a7-b8c9-0123-def4-567890abcdef",
        "name": "Siti Rahma",
        "phone": "081234567891",
        "address": "Jl. Mawar No. 5, Jakarta"
      },
      "discount": {
        "id": "c12d3456-7890-abcd-ef01-234567890abc",
        "name": "Promo Kemerdekaan",
        "type": "percentage",
        "value": 15
      },
      "creator": {
        "id": "76d54cf8-a1fb-45ba-b072-467f53f93ce0",
        "name": "Budi Santoso"
      },
      "items": [
        {
          "id": "item-uuid-0001-abcd-ef0123456789",
          "price": 8000,
          "quantity": 5000,
          "unit": "kg",
          "subtotal": 40000,
          "service": {
            "id": "e44d348a-6b87-43cf-be72-c5180f96894c",
            "name": "Cuci Setrika Reguler",
            "price": 8000,
            "duration_days": 2,
            "category": {
              "id": "9059f13e-63f6-4927-a068-07e4d8fb8756",
              "name": "Kiloan"
            }
          }
        }
      ],
      "payment": {
        "id": "pay-uuid-0001-abcd-ef0123456789",
        "amount": 34000,
        "payment_method": "cash",
        "status": "verified",
        "reference_number": null,
        "notes": null,
        "paid_at": "2026-10-16T10:00:00Z",
        "verified_at": "2026-10-16T10:00:00Z",
        "created_at": "2026-10-16T10:00:00Z",
        "updated_at": "2026-10-16T10:00:00Z"
      },
      "notes": "Tolong pakai pewangi lavender",
      "estimated_done": "2026-10-18T10:00:00Z",
      "completed_at": null,
      "canceled_at": null,
      "created_at": "2026-10-16T10:00:00Z",
      "updated_at": "2026-10-16T10:00:00Z"
    },
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `404` (outlet, customer, atau service tidak ditemukan), `422` (outlet tidak aktif)

---

### GET /api/v1/orders
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Deskripsi**: Mengambil daftar order milik tenant dengan berbagai filter.
- **Query Parameters**:
  | Parameter | Tipe | Default | Deskripsi |
  |---|---|---|---|
  | `search` | string | — | Search pada nomor order atau nama customer |
  | `status` | string | — | Filter status: `MENUNGGU_KONFIRMASI`, `DIPROSES`, `SELESAI`, `DIAMBIL`, `DIBATALKAN` |
  | `outlet_id` | UUID | — | Filter berdasarkan outlet |
  | `customer_id` | UUID | — | Filter berdasarkan customer |
  | `start_date` | string | — | Filter tanggal mulai (format: `YYYY-MM-DD`) |
  | `end_date` | string | — | Filter tanggal akhir (format: `YYYY-MM-DD`) |
  | `page` | integer | `1` | Halaman saat ini |
  | `limit` | integer | `10` | Jumlah per halaman (max: 100) |

- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "orders retrieved successfully",
    "success": true,
    "data": [
      {
        "id": "ord-uuid-0001-abcd-ef0123456789",
        "tenant_id": "1e5db5d2-0056-4c4f-9e79-5ee45b98df99",
        "order_number": "ORD-20261016-0001",
        "tracking_token": "TRK-AbCdEfGhIjKl",
        "status": "DIPROSES",
        "payment_status": "UNPAID",
        "subtotal": 40000,
        "discount_amount": 6000,
        "total_amount": 34000,
        "paid_amount": 0,
        "change_amount": 0,
        "outlet": {
          "id": "e674b9e2-3cf7-4fbc-b4e8-4680dc1c1f7b",
          "name": "Outlet Cabang Melati",
          "phone": "081234567890",
          "address": "Jl. Melati No. 12, Jakarta"
        },
        "customer": {
          "id": "d4e5f6a7-b8c9-0123-def4-567890abcdef",
          "name": "Siti Rahma",
          "phone": "081234567891",
          "address": "Jl. Mawar No. 5, Jakarta"
        },
        "discount": null,
        "creator": {
          "id": "76d54cf8-a1fb-45ba-b072-467f53f93ce0",
          "name": "Budi Santoso"
        },
        "items": [ ... ],
        "payment": { ... },
        "notes": null,
        "estimated_done": "2026-10-18T10:00:00Z",
        "completed_at": null,
        "canceled_at": null,
        "created_at": "2026-10-16T10:00:00Z",
        "updated_at": "2026-10-16T10:00:00Z"
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

### GET /api/v1/orders/:id
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Path Params**: `id` (UUID order)
- **Response** `200 OK`: Single `OrderResponse` object (format lengkap seperti di atas).
- **Error Responses**: `400` (ID tidak valid), `404` (order tidak ditemukan)

---

### PATCH /api/v1/orders/:id/status
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Path Params**: `id` (UUID order)
- **Deskripsi**: Memperbarui status pemrosesan order.
- **Request Body**:
  ```json
  {
    "status": "SELESAI"
  }
  ```
- **Validation Rules**:
  | Field | Aturan |
  |---|---|
  | `status` | required, one of: `MENUNGGU_KONFIRMASI`, `DIPROSES`, `SELESAI`, `DIAMBIL` |

- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "order status updated successfully",
    "success": true,
    "data": { ...OrderResponse... },
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `404`, `422` (transisi status tidak valid)

---

### PATCH /api/v1/orders/:id/cancel
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Path Params**: `id` (UUID order)
- **Deskripsi**: Membatalkan order. Hanya order dengan status `MENUNGGU_KONFIRMASI` yang dapat dibatalkan.
- **Request Body**: Tidak diperlukan.
- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "order canceled successfully",
    "success": true,
    "data": { ...OrderResponse dengan status DIBATALKAN... },
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `404`, `422` (status order tidak memungkinkan dibatalkan)

---

### PATCH /api/v1/orders/:id/payments/:payment_id/confirm
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Path Params**: `id` (UUID order), `payment_id` (UUID payment)
- **Deskripsi**: Mengkonfirmasi pembayaran order. Payment status berubah menjadi `PAID`.
- **Request Body**: Tidak diperlukan.
- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "payment confirmed successfully",
    "success": true,
    "data": { ...OrderResponse dengan payment_status PAID... },
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `404`

---

### PATCH /api/v1/orders/:id/cash-payment
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Path Params**: `id` (UUID order)
- **Deskripsi**: Menginput jumlah uang tunai yang dibayarkan customer. Sistem akan menghitung kembalian (`change_amount`).
- **Request Body**:
  ```json
  {
    "paid_amount": 50000
  }
  ```
- **Validation Rules**:
  | Field | Aturan |
  |---|---|
  | `paid_amount` | required, gte=0 |

- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "cash payment input successfully",
    "success": true,
    "data": {
      "id": "ord-uuid-0001-abcd-ef0123456789",
      "order_number": "ORD-20261016-0001",
      "total_amount": 34000,
      "paid_amount": 50000,
      "change_amount": 16000,
      "payment_status": "PAID",
      ...
    },
    "meta": {}
  }
  ```
- **Error Responses**: `400`, `404`, `422` (metode pembayaran bukan `cash`)

---

## Public Order Endpoints

Endpoint berikut **tidak memerlukan autentikasi** dan ditujukan untuk customer yang melakukan order mandiri melalui link publik outlet.

---

### POST /api/v1/public/orders
- **Auth**: Public (tidak perlu token)
- **Deskripsi**: Customer membuat order secara mandiri melalui halaman publik outlet. Customer tidak harus terdaftar — jika nomor telepon belum ada, customer baru akan dibuat otomatis. `is_public_order_enabled` pada outlet harus bernilai `true`.
- **Request Body**:
  ```json
  {
    "outlet_slug": "outlet-cabang-melati",
    "name": "Dewi Sartika",
    "phone": "081234567899",
    "address": "Jl. Anggrek No. 3",
    "discount_id": null,
    "payment_method": "bank_transfer",
    "paid_amount": null,
    "notes": "Pisahkan baju putih dan berwarna",
    "items": [
      {
        "service_id": "e44d348a-6b87-43cf-be72-c5180f96894c",
        "quantity": 3000
      }
    ]
  }
  ```
- **Validation Rules**:
  | Field | Aturan |
  |---|---|
  | `outlet_slug` | required, min=1, max=100 |
  | `name` | required, min=3, max=150 |
  | `phone` | required, format nomor HP Indonesia |
  | `address` | opsional, max=500 |
  | `discount_id` | opsional, UUID |
  | `payment_method` | required, one of: `cash`, `bank_transfer`, `qris` |
  | `paid_amount` | opsional, gte=0 |
  | `notes` | opsional, max=1000 |
  | `items` | required, min=1 array item |
  | `items[].service_id` | required, UUID |
  | `items[].quantity` | required, gt=0 |

- **Response** `201 Created`:
  ```json
  {
    "status": 201,
    "message": "public order created successfully",
    "success": true,
    "data": {
      "id": "ord-uuid-0002-abcd-ef0123456789",
      "order_number": "ORD-20261016-0002",
      "tracking_token": "TRK-XyZaBcDeFgHi",
      "status": "MENUNGGU_KONFIRMASI",
      "payment_status": "UNPAID",
      "subtotal": 24000,
      "discount_amount": 0,
      "total_amount": 24000,
      "paid_amount": 0,
      "change_amount": 0,
      "outlet": {
        "id": "e674b9e2-3cf7-4fbc-b4e8-4680dc1c1f7b",
        "name": "Outlet Cabang Melati",
        "phone": "081234567890",
        "address": "Jl. Melati No. 12, Jakarta"
      },
      "customer": {
        "id": "new-customer-uuid-abcd",
        "name": "Dewi Sartika",
        "phone": "081234567899",
        "address": "Jl. Anggrek No. 3"
      },
      "discount": null,
      "creator": null,
      "items": [ ... ],
      "payment": { ... },
      "notes": "Pisahkan baju putih dan berwarna",
      "estimated_done": "2026-10-18T10:00:00Z",
      "completed_at": null,
      "canceled_at": null,
      "created_at": "2026-10-16T10:00:00Z",
      "updated_at": "2026-10-16T10:00:00Z"
    },
    "meta": {}
  }
  ```
- **Error Responses**:
  | Status | Kondisi |
  |---|---|
  | `400` | Validasi request gagal |
  | `404` | Outlet dengan slug tersebut tidak ditemukan |
  | `422` | Outlet tidak mengaktifkan public order (`is_public_order_enabled = false`) atau outlet tidak aktif |

---

### GET /api/v1/public/orders/track/:token
- **Auth**: Public (tidak perlu token)
- **Path Params**: `token` — tracking token yang didapat saat order dibuat
- **Deskripsi**: Customer dapat melacak status ordernya menggunakan tracking token.
- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "order tracked successfully",
    "success": true,
    "data": {
      "id": "ord-uuid-0002-abcd-ef0123456789",
      "order_number": "ORD-20261016-0002",
      "tracking_token": "TRK-XyZaBcDeFgHi",
      "status": "DIPROSES",
      "payment_status": "UNPAID",
      "subtotal": 24000,
      "discount_amount": 0,
      "total_amount": 24000,
      "paid_amount": 0,
      "change_amount": 0,
      "outlet": {
        "id": "e674b9e2-3cf7-4fbc-b4e8-4680dc1c1f7b",
        "name": "Outlet Cabang Melati",
        "phone": "081234567890",
        "address": "Jl. Melati No. 12, Jakarta"
      },
      "customer": {
        "id": "new-customer-uuid-abcd",
        "name": "Dewi Sartika",
        "phone": "081234567899",
        "address": "Jl. Anggrek No. 3"
      },
      "discount": null,
      "creator": null,
      "items": [
        {
          "id": "item-uuid-0002-abcd",
          "price": 8000,
          "quantity": 3000,
          "unit": "kg",
          "subtotal": 24000,
          "service": {
            "id": "e44d348a-6b87-43cf-be72-c5180f96894c",
            "name": "Cuci Setrika Reguler",
            "price": 8000,
            "duration_days": 2,
            "category": {
              "id": "9059f13e-63f6-4927-a068-07e4d8fb8756",
              "name": "Kiloan"
            }
          }
        }
      ],
      "payment": {
        "id": "pay-uuid-0002-abcd",
        "amount": 24000,
        "payment_method": "bank_transfer",
        "status": "unpaid",
        "reference_number": null,
        "notes": null,
        "paid_at": null,
        "verified_at": null,
        "created_at": "2026-10-16T10:00:00Z",
        "updated_at": "2026-10-16T10:00:00Z"
      },
      "notes": "Pisahkan baju putih dan berwarna",
      "estimated_done": "2026-10-18T10:00:00Z",
      "completed_at": null,
      "canceled_at": null,
      "created_at": "2026-10-16T10:00:00Z",
      "updated_at": "2026-10-16T10:00:00Z"
    },
    "meta": {}
  }
  ```
- **Error Responses**: `400` (token kosong), `404` (order dengan token tidak ditemukan)

---

## 13. Report Module

### GET /api/v1/dashboard
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Deskripsi**: Mengambil ringkasan statistik operasional dan keuangan (summary, rekap status order, rekap metode pembayaran, dan grafik penjualan).
- **Query Parameters**:
  | Parameter | Tipe | Default | Deskripsi |
  |---|---|---|---|
  | `start_date` | string | Bulan berjalan | Tanggal mulai filter (format: `YYYY-MM-DD`) |
  | `end_date` | string | Bulan berjalan | Tanggal akhir filter (format: `YYYY-MM-DD`) |
  | `outlet_id` | UUID | — | Filter berdasarkan outlet tertentu |

- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "dashboard retrieved successfully",
    "success": true,
    "data": {
      "summary": {
        "total_sales": 1500000,
        "total_payments_received": 1200000,
        "total_orders": 10,
        "active_orders": 3
      },
      "order_status": [
        {"status": "MENUNGGU_KONFIRMASI", "count": 2},
        {"status": "DIPROSES", "count": 1},
        {"status": "SELESAI", "count": 5},
        {"status": "DIAMBIL", "count": 2},
        {"status": "DIBATALKAN", "count": 0}
      ],
      "payment_methods": [
        {"method": "cash", "total": 1200000, "count": 8}
      ],
      "sales_chart": [
        {"date": "2026-10-01", "total_sales": 500000, "total_payments_received": 400000}
      ]
    },
    "meta": {}
  }
  ```

---

### GET /api/v1/reports/financial
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Deskripsi**: Mengambil laporan keuangan dan daftar transaksi pembayaran terverifikasi dengan paginasi.
- **Query Parameters**:
  | Parameter | Tipe | Default | Deskripsi |
  |---|---|---|---|
  | `start_date` | string | Bulan berjalan | Tanggal mulai filter (format: `YYYY-MM-DD`) |
  | `end_date` | string | Bulan berjalan | Tanggal akhir filter (format: `YYYY-MM-DD`) |
  | `outlet_id` | UUID | — | Filter berdasarkan outlet tertentu |
  | `page` | integer | `1` | Halaman saat ini |
  | `limit` | integer | `10` | Jumlah per halaman (max: 100) |

- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "financial report retrieved successfully",
    "success": true,
    "data": {
      "summary": {
        "total_sales": 1500000,
        "total_discount": 100000,
        "total_payments_received": 1200000,
        "total_orders": 10
      },
      "payments": [
        {
          "payment_id": "pay-uuid-0001",
          "order_id": "ord-uuid-0001",
          "order_number": "ORD-20261016-0001",
          "customer_name": "Siti Rahma",
          "outlet_name": "Outlet Cabang Melati",
          "payment_method": "cash",
          "amount": 1200000,
          "status": "verified",
          "paid_at": "2026-10-16T10:00:00Z"
        }
      ]
    },
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

## 14. Profile Module

### GET /api/v1/profile
- **Auth**: `Bearer <token>`
- **Role**: `ADMIN, STAFF`
- **Deskripsi**: Mendapatkan data profile.
- **Response** `200 OK`:
  ```json
  {
    "status": 200,
    "message": "profile retrieved successfully",
    "success": true,
    "data": {
      "id": "1e5db5d2-0056-4c4f-9e79-5ee45b98df99",
      "name": "profile",
      "role": "STAFF",
      "email": "profile@gmail.com",
      "created_at": "2026-10-16T10:00:00Z",
    },
    "meta": {}
  }
  ```
- **Error Responses**: `404` (profile tidak ditemukan)

---

## Kode Error Umum

| HTTP Status | Kode | Kondisi Umum |
|---|---|---|
| `400` | Bad Request | Request body tidak valid, format JSON salah, atau validasi field gagal |
| `401` | Unauthorized | Token JWT tidak ada, tidak valid, atau sudah expired |
| `403` | Forbidden | Token valid tetapi role tidak memiliki akses |
| `404` | Not Found | Resource yang diminta tidak ditemukan |
| `409` | Conflict | Data duplikat (email, nama, slug yang sudah ada) |
| `422` | Unprocessable | Business rule dilanggar (e.g. outlet tidak aktif) |
| `500` | Internal Server Error | Kesalahan server tak terduga |
