# Database Documentation — Lavendera API

> **Status**: Verifikasi dari implementasi aktual kode.  
> **Source of Truth**: [internal/models/](file:///home/wenell/projects/lavendera-API/internal/models/) dan file migrasi SQL di [internal/db/migrations/](file:///home/wenell/projects/lavendera-API/internal/db/migrations/).

---

## 1. Database Engine & Driver

- **Engine**: PostgreSQL (Didukung ekstensi `uuid-ossp`)
- **Go ORM**: GORM v1.31.2 (`gorm.io/gorm`)
- **Driver**: `gorm.io/driver/postgres` v1.6.2 (menggunakan `pgx/v5`)
- **Connection Pool Configuration** ([database.go](file:///home/wenell/projects/lavendera-API/internal/shared/database/database.go)):
  - `SetMaxIdleConns(10)`
  - `SetMaxOpenConns(100)`
  - `SetConnMaxLifetime(1 * time.Hour)`
  - Logging GORM: `logger.Info`

---

## 2. Proven Relationship Map

Berdasarkan deklarasi struct di [internal/models/](file:///home/wenell/projects/lavendera-API/internal/models/) dan foreign keys pada SQL migration:

```
Tenant
 ├── User (1:N)
 │    ├── OutletUser (1:N via user_id)
 │    ├── Order (1:N via created_by)
 │    ├── Payment (1:N via received_by)
 │    ├── Notification (1:N via user_id)
 │    └── AuditLog (1:N via user_id)
 │
 ├── Outlet (1:N)
 │    ├── OutletUser (1:N via outlet_id)
 │    ├── OutletPaymentMethod (1:N via outlet_id)
 │    ├── Order (1:N via outlet_id)
 │    └── Payment (1:N via outlet_id)
 │
 ├── ServiceCategory (1:N)
 │    └── Service (1:N via category_id)
 │         └── OrderItem (1:N via service_id)
 │
 ├── Discount (1:N)
 │    └── Order (1:N via discount_id)
 │
 ├── Customer (1:N)
 │    └── Order (1:N via customer_id)
 │         ├── OrderItem (1:N via order_id)
 │         ├── Payment (1:N via order_id)
 │         └── Notification (1:N via order_id)
 │
 └── AuditLog (1:N via tenant_id)
```

---

## 3. Schema & Model Details

### 3.1. `tenants`
- **Model**: [models.Tenant](file:///home/wenell/projects/lavendera-API/internal/models/tenant.go)
- **Migration**: [20260908053004_create_tenants.up.sql](file:///home/wenell/projects/lavendera-API/internal/db/migrations/20260908053004_create_tenants.up.sql)
- **Columns**:
  - `id` (UUID, PK, Default: `uuid_generate_v4()`)
  - `name` (VARCHAR(150), NOT NULL)
  - `slug` (VARCHAR(100), NOT NULL, UNIQUE)
  - `email` (VARCHAR(255), NOT NULL)
  - `created_at` (TIMESTAMPTZ, DEFAULT NOW())
  - `updated_at` (TIMESTAMPTZ, DEFAULT NOW())
- **Constraints**: `CONSTRAINT unique_tenant_slug UNIQUE (slug)`

---

### 3.2. `users`
- **Model**: [models.User](file:///home/wenell/projects/lavendera-API/internal/models/user.go)
- **Migration**: [20260908061134_create_users.up.sql](file:///home/wenell/projects/lavendera-API/internal/db/migrations/20260908061134_create_users.up.sql)
- **Columns**:
  - `id` (UUID, PK, Default: `uuid_generate_v4()`)
  - `tenant_id` (UUID, NOT NULL, FK `tenants(id)` ON DELETE CASCADE)
  - `name` (VARCHAR(150), NOT NULL)
  - `email` (VARCHAR(255), NOT NULL)
  - `password` (VARCHAR(255), NOT NULL — bcrypt hash)
  - `role` (VARCHAR(20), NOT NULL, DEFAULT `'STAFF'`)
  - `is_active` (BOOLEAN, NOT NULL, DEFAULT `TRUE`)
  - `created_at` (TIMESTAMPTZ, DEFAULT NOW())
  - `updated_at` (TIMESTAMPTZ, DEFAULT NOW())
- **Constraints**:
  - `CHECK (role IN ('ADMIN', 'STAFF'))`
  - `CONSTRAINT unique_user_email_per_tenant UNIQUE (tenant_id, email)`
- **Indexes**: `CREATE INDEX idx_users_tenant ON users (tenant_id)`

---

### 3.3. `outlets`
- **Model**: [models.Outlet](file:///home/wenell/projects/lavendera-API/internal/models/outlet.go)
- **Migration**: [20260908061245_create_outlets.up.sql](file:///home/wenell/projects/lavendera-API/internal/db/migrations/20260908061245_create_outlets.up.sql)
- **Columns**:
  - `id` (UUID, PK, Default: `uuid_generate_v4()`)
  - `tenant_id` (UUID, NOT NULL, FK `tenants(id)` ON DELETE CASCADE)
  - `name` (VARCHAR(150), NOT NULL)
  - `slug` (VARCHAR(100), NOT NULL)
  - `phone` (VARCHAR(20), NULLABLE)
  - `address` (TEXT, NULLABLE)
  - `is_public_order_enabled` (BOOLEAN, NOT NULL, DEFAULT `TRUE`)
  - `is_active` (BOOLEAN, NOT NULL, DEFAULT `TRUE`)
  - `created_at` (TIMESTAMPTZ, DEFAULT NOW())
  - `updated_at` (TIMESTAMPTZ, DEFAULT NOW())
- **Constraints**: `CONSTRAINT unique_outlet_slug_per_tenant UNIQUE (tenant_id, slug)`
- **Indexes**: `CREATE INDEX idx_outlets_tenant ON outlets (tenant_id)`

---

### 3.4. `outlet_users`
- **Model**: [models.OutletUser](file:///home/wenell/projects/lavendera-API/internal/models/outlet_user.go)
- **Migration**: [20260908061613_create_outlet_users.up.sql](file:///home/wenell/projects/lavendera-API/internal/db/migrations/20260908061613_create_outlet_users.up.sql)
- **Columns**:
  - `outlet_id` (UUID, FK `outlets(id)` ON DELETE CASCADE)
  - `user_id` (UUID, FK `users(id)` ON DELETE CASCADE)
  - `created_at` (TIMESTAMPTZ, DEFAULT NOW())
- **Constraints**: Composite Primary Key `PRIMARY KEY (outlet_id, user_id)`

---

### 3.5. `service_categories`
- **Model**: [models.ServiceCategory](file:///home/wenell/projects/lavendera-API/internal/models/service_category.go)
- **Migration**: [20260908061502_create_service_categories.up.sql](file:///home/wenell/projects/lavendera-API/internal/db/migrations/20260908061502_create_service_categories.up.sql)
- **Columns**:
  - `id` (UUID, PK, Default: `uuid_generate_v4()`)
  - `tenant_id` (UUID, NOT NULL, FK `tenants(id)` ON DELETE CASCADE)
  - `name` (VARCHAR(100), NOT NULL)
  - `created_at` (TIMESTAMPTZ, DEFAULT NOW())
  - `updated_at` (TIMESTAMPTZ, DEFAULT NOW())

---

### 3.6. `services`
- **Model**: [models.Service](file:///home/wenell/projects/lavendera-API/internal/models/service.go)
- **Migration**: [20260908061644_create_services.up.sql](file:///home/wenell/projects/lavendera-API/internal/db/migrations/20260908061644_create_services.up.sql)
- **Columns**:
  - `id` (UUID, PK, Default: `uuid_generate_v4()`)
  - `tenant_id` (UUID, NOT NULL, FK `tenants(id)` ON DELETE CASCADE)
  - `category_id` (UUID, FK `service_categories(id)` ON DELETE SET NULL)
  - `name` (VARCHAR(255), NOT NULL)
  - `price` (BIGINT, NOT NULL, DEFAULT `0`)
  - `min_quantity` (BIGINT, NOT NULL, DEFAULT `1`)
  - `unit` (VARCHAR(20), NOT NULL, DEFAULT `'gram'`)
  - `duration_days` (INTEGER, NOT NULL, DEFAULT `1`)
  - `is_active` (BOOLEAN, NOT NULL, DEFAULT `TRUE`)
  - `created_at` (TIMESTAMPTZ, DEFAULT NOW())
  - `updated_at` (TIMESTAMPTZ, DEFAULT NOW())
- **Constraints**: `CHECK (unit IN ('gram', 'pcs', 'pair'))`
- **Indexes**: `CREATE INDEX idx_services_tenant ON services (tenant_id)`

---

### 3.7. `discounts`
- **Model**: [models.Discount](file:///home/wenell/projects/lavendera-API/internal/models/discount.go)
- **Migration**: [20260908061536_create_discounts.up.sql](file:///home/wenell/projects/lavendera-API/internal/db/migrations/20260908061536_create_discounts.up.sql)
- **Columns**:
  - `id` (UUID, PK, Default: `uuid_generate_v4()`)
  - `tenant_id` (UUID, NOT NULL, FK `tenants(id)` ON DELETE CASCADE)
  - `name` (VARCHAR(150), NOT NULL)
  - `type` (VARCHAR(20), NOT NULL)
  - `value` (BIGINT, NOT NULL, DEFAULT `0`)
  - `is_active` (BOOLEAN, NOT NULL, DEFAULT `TRUE`)
  - `created_at` (TIMESTAMPTZ, DEFAULT NOW())
  - `updated_at` (TIMESTAMPTZ, DEFAULT NOW())
- **Constraints**: `CHECK (type IN ('percentage', 'fixed'))`

---

### 3.8. `customers`
- **Model**: [models.Customer](file:///home/wenell/projects/lavendera-API/internal/models/customer.go)
- **Migration**: [20260908053325_create_customers.up.sql](file:///home/wenell/projects/lavendera-API/internal/db/migrations/20260908053325_create_customers.up.sql)
- **Columns**:
  - `id` (UUID, PK, Default: `uuid_generate_v4()`)
  - `tenant_id` (UUID, NOT NULL, FK `tenants(id)` ON DELETE CASCADE)
  - `name` (VARCHAR(150), NOT NULL)
  - `phone` (VARCHAR(20), NOT NULL)
  - `address` (TEXT, NULLABLE)
  - `created_at` (TIMESTAMPTZ, DEFAULT NOW())
  - `updated_at` (TIMESTAMPTZ, DEFAULT NOW())
- **Indexes**: `CREATE INDEX idx_customers_tenant ON customers (tenant_id)`
- **Module Status**: Belum diimplementasikan layer controller/service/repo (Model only).

---

### 3.9. `outlet_payment_methods`
- **Model**: [models.OutletPaymentMethod](file:///home/wenell/projects/lavendera-API/internal/models/outlet_payment_method.go)
- **Migration**: [20260908061713_create_outlet_payment_methods.up.sql](file:///home/wenell/projects/lavendera-API/internal/db/migrations/20260908061713_create_outlet_payment_methods.up.sql)
- **Columns**:
  - `id` (UUID, PK, Default: `uuid_generate_v4()`)
  - `outlet_id` (UUID, NOT NULL, FK `outlets(id)` ON DELETE CASCADE)
  - `type` (VARCHAR(50), NOT NULL)
  - `provider_name` (VARCHAR(100), NOT NULL)
  - `account_number` (VARCHAR(50), NULLABLE)
  - `account_name` (VARCHAR(150), NOT NULL)
  - `qr_image_url` (VARCHAR(500), NULLABLE)
  - `is_active` (BOOLEAN, NOT NULL, DEFAULT `TRUE`)
  - `created_at` (TIMESTAMPTZ, DEFAULT NOW())
  - `updated_at` (TIMESTAMPTZ, DEFAULT NOW())
- **Constraints**: `CHECK (type IN ('qris', 'bank_transfer'))`
- **Module Status**: Model only.

---

### 3.10. `orders`
- **Model**: [models.Order](file:///home/wenell/projects/lavendera-API/internal/models/order.go)
- **Migration**: [20260908062453_create_orders.up.sql](file:///home/wenell/projects/lavendera-API/internal/db/migrations/20260908062453_create_orders.up.sql)
- **Columns**:
  - `id` (UUID, PK, Default: `uuid_generate_v4()`)
  - `tenant_id` (UUID, NOT NULL, FK `tenants(id)` ON DELETE CASCADE)
  - `outlet_id` (UUID, NOT NULL, FK `outlets(id)` ON DELETE RESTRICT)
  - `customer_id` (UUID, NOT NULL, FK `customers(id)` ON DELETE RESTRICT)
  - `discount_id` (UUID, NULLABLE, FK `discounts(id)` ON DELETE SET NULL)
  - `created_by` (UUID, NULLABLE, FK `users(id)` ON DELETE SET NULL)
  - `order_number` (VARCHAR(100), NOT NULL)
  - `tracking_token` (VARCHAR(100), NOT NULL, UNIQUE)
  - `status` (VARCHAR(50), NOT NULL, DEFAULT `'MENUNGGU_KONFIRMASI'`)
  - `subtotal` (BIGINT, NOT NULL, DEFAULT `0`)
  - `discount_amount` (BIGINT, NOT NULL, DEFAULT `0`)
  - `total` (BIGINT, NOT NULL, DEFAULT `0`)
  - `paid_amount` (BIGINT, NOT NULL, DEFAULT `0`)
  - `notes` (TEXT, NULLABLE)
  - `estimated_done` (TIMESTAMPTZ, NULLABLE)
  - `completed_at` (TIMESTAMPTZ, NULLABLE)
  - `canceled_at` (TIMESTAMPTZ, NULLABLE)
  - `created_at` (TIMESTAMPTZ, DEFAULT NOW())
  - `updated_at` (TIMESTAMPTZ, DEFAULT NOW())
- **Constraints**:
  - `CHECK (status IN ('MENUNGGU_KONFIRMASI', 'DIPROSES', 'SELESAI', 'DIAMBIL', 'DIBATALKAN'))`
  - `CONSTRAINT unique_order_number_per_tenant UNIQUE (tenant_id, order_number)`
- **Indexes**: `idx_orders_tenant`, `idx_orders_outlet`, `idx_orders_customer`, `idx_orders_number`
- **Module Status**: Model only.

---

### 3.11. `order_items`
- **Model**: [models.OrderItem](file:///home/wenell/projects/lavendera-API/internal/models/order_item.go)
- **Migration**: [20260908062517_create_order_items.up.sql](file:///home/wenell/projects/lavendera-API/internal/db/migrations/20260908062517_create_order_items.up.sql)
- **Columns**:
  - `id` (UUID, PK, Default: `uuid_generate_v4()`)
  - `order_id` (UUID, NOT NULL, FK `orders(id)` ON DELETE CASCADE)
  - `service_id` (UUID, FK `services(id)` ON DELETE RESTRICT)
  - `price` (BIGINT, NOT NULL, DEFAULT `0`)
  - `quantity` (BIGINT, NOT NULL, DEFAULT `0`)
  - `subtotal` (BIGINT, NOT NULL, DEFAULT `0`)
  - `created_at` (TIMESTAMPTZ, DEFAULT NOW())
  - `updated_at` (TIMESTAMPTZ, DEFAULT NOW())
- **Indexes**: `CREATE INDEX idx_order_items_order ON order_items (order_id)`
- **Module Status**: Model only.

---

### 3.12. `payments`
- **Model**: [models.Payment](file:///home/wenell/projects/lavendera-API/internal/models/payment.go)
- **Migration**: [20260908062851_create_payments.up.sql](file:///home/wenell/projects/lavendera-API/internal/db/migrations/20260908062851_create_payments.up.sql)
- **Columns**:
  - `id` (UUID, PK, Default: `uuid_generate_v4()`)
  - `outlet_id` (UUID, NOT NULL, FK `outlets(id)` ON DELETE RESTRICT)
  - `order_id` (UUID, NOT NULL, FK `orders(id)` ON DELETE CASCADE)
  - `received_by` (UUID, NULLABLE, FK `users(id)` ON DELETE SET NULL)
  - `amount` (BIGINT, NOT NULL, DEFAULT `0`)
  - `payment_method` (VARCHAR(50), NOT NULL)
  - `status` (VARCHAR(20), NOT NULL, DEFAULT `'pending'`)
  - `reference_number` (VARCHAR(100), NULLABLE)
  - `notes` (TEXT, NULLABLE)
  - `paid_at` (TIMESTAMPTZ, DEFAULT NOW())
  - `verified_at` (TIMESTAMPTZ, NULLABLE)
  - `created_at` (TIMESTAMPTZ, DEFAULT NOW())
  - `updated_at` (TIMESTAMPTZ, DEFAULT NOW())
- **Constraints**:
  - `CHECK (payment_method IN ('cash', 'qris', 'bank_transfer'))`
  - `CHECK (status IN ('pending', 'verified', 'rejected'))`
- **Indexes**: `CREATE INDEX idx_payments_order ON payments (order_id)`
- **Module Status**: Model only.

---

### 3.13. `notifications`
- **Model**: [models.Notification](file:///home/wenell/projects/lavendera-API/internal/models/notification.go)
- **Migration**: [20260908062933_create_notifications.up.sql](file:///home/wenell/projects/lavendera-API/internal/db/migrations/20260908062933_create_notifications.up.sql)
- **Columns**:
  - `id` (UUID, PK, Default: `uuid_generate_v4()`)
  - `order_id` (UUID, NULLABLE, FK `orders(id)` ON DELETE CASCADE)
  - `user_id` (UUID, NULLABLE, FK `users(id)` ON DELETE CASCADE)
  - `type` (VARCHAR(50), NOT NULL)
  - `message` (TEXT, NOT NULL)
  - `is_read` (BOOLEAN, NOT NULL, DEFAULT `false`)
  - `sent_at` (TIMESTAMPTZ, NULLABLE)
  - `created_at` (TIMESTAMPTZ, DEFAULT NOW())
- **Module Status**: Model only.

---

### 3.14. `audit_logs`
- **Model**: [models.AuditLog](file:///home/wenell/projects/lavendera-API/internal/models/audit_logs.go)
- **Migration**: [20260908061215_create_audit_logs.up.sql](file:///home/wenell/projects/lavendera-API/internal/db/migrations/20260908061215_create_audit_logs.up.sql)
- **Columns**:
  - `id` (UUID, PK, Default: `uuid_generate_v4()`)
  - `tenant_id` (UUID, NOT NULL, FK `tenants(id)` ON DELETE CASCADE)
  - `user_id` (UUID, NULLABLE, FK `users(id)` ON DELETE SET NULL)
  - `action` (VARCHAR(100), NOT NULL)
  - `entity_type` (VARCHAR(100), NOT NULL)
  - `old_data` (JSONB, NULLABLE)
  - `new_data` (JSONB, NULLABLE)
  - `created_at` (TIMESTAMPTZ, DEFAULT NOW())
- **Module Status**: Model only.

---

## 4. Database Triggers

Trigger PostgreSQL `update_timestamp_column()` didefinisikan pada migrasi [20260908065344_create_updated_at_triggers.up.sql](file:///home/wenell/projects/lavendera-API/internal/db/migrations/20260908065344_create_updated_at_triggers.up.sql).
Trigger ini otomatis memperbarui nilai kolom `updated_at` menjadi `NOW()` sesaat sebelum baris data diperbarui (`BEFORE UPDATE`) pada 11 tabel:
- `tenants`, `customers`, `users`, `outlets`, `service_categories`, `discounts`, `services`, `orders`, `order_items`, `payments`, `outlet_payment_methods`.

---

## 5. Transaction Usage in Codebase

Penggunaan transaksi database aktual ditemukan pada [AuthRepositoryImpl.CreateDefaultAdmin](file:///home/wenell/projects/lavendera-API/internal/auth/repository/auth_repository_impl.go#L19-L36):
```go
func (a *AuthRepositoryImpl) CreateDefaultAdmin(ctx context.Context, tenant *models.Tenant, user *models.User, categories []models.ServiceCategory) error {
    return a.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(tenant).Error; err != nil {
            return err
        }
        user.TenantID = tenant.ID
        if err := tx.Create(user).Error; err != nil {
            return err
        }
        for i := range categories {
            categories[i].TenantID = tenant.ID
        }
        if err := tx.Create(categories).Error; err != nil {
            return err
        }
        return nil
    })
}
```

---

## 6. Multi-Tenant Isolation Mechanism

Isolasi tenant tidak mengandalkan PostgreSQL Row-Level Security (RLS) di database melainkan diatur pada application-level melalui GORM Scope [TenantScope](file:///home/wenell/projects/lavendera-API/internal/shared/database/scope.go):
```go
func TenantScope(ctx context.Context) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        tenantID, ok := appcontext.TenantIDFromContext(ctx)
        if !ok || tenantID == uuid.Nil {
            return db.Where("1 = 0")
        }
        return db.Where("tenant_id = ?", tenantID)
    }
}
```
Jika `tenant_id` tidak dapat diekstrak dari JWT context, klausa `1 = 0` menjamin tidak ada satupun row yang dapat dibaca atau dimodifikasi oleh query tersebut.
