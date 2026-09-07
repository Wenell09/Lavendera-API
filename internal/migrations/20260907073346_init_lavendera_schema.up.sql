-- =============================================================================
-- LAVENDERA DATABASE SCHEMA 
-- =============================================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Trigger Function untuk updated_at Otomatis
CREATE OR REPLACE FUNCTION update_timestamp_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- -----------------------------------------------------------------------------
-- LEVEL 1 & 2: TENANTS & TENANT ENTITIES
-- -----------------------------------------------------------------------------

CREATE TABLE "TENANTS" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "name" VARCHAR(255) NOT NULL,
  "logo_url" VARCHAR(500),
  "domain" VARCHAR(255) UNIQUE NOT NULL,
  "email" VARCHAR(255) NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TRIGGER update_tenants_modtime 
  BEFORE UPDATE ON "TENANTS" 
  FOR EACH ROW EXECUTE PROCEDURE update_timestamp_column();


CREATE TABLE "CUSTOMERS" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "tenant_id" UUID NOT NULL REFERENCES "TENANTS"("id") ON DELETE CASCADE,
  "name" VARCHAR(255) NOT NULL,
  "phone" VARCHAR(50) NOT NULL,
  "address" TEXT,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TRIGGER update_customers_modtime 
  BEFORE UPDATE ON "CUSTOMERS" 
  FOR EACH ROW EXECUTE PROCEDURE update_timestamp_column();


CREATE TABLE "USERS" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "tenant_id" UUID NOT NULL REFERENCES "TENANTS"("id") ON DELETE CASCADE,
  "name" VARCHAR(255) NOT NULL,
  "email" VARCHAR(255) NOT NULL,
  "password_hash" VARCHAR(255) NOT NULL,
  "role" VARCHAR(50) NOT NULL DEFAULT 'STAFF',
  "is_active" BOOLEAN NOT NULL DEFAULT TRUE,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  CONSTRAINT "unique_user_email_per_tenant" UNIQUE ("tenant_id", "email")
);

CREATE TRIGGER update_users_modtime 
  BEFORE UPDATE ON "USERS" 
  FOR EACH ROW EXECUTE PROCEDURE update_timestamp_column();


CREATE TABLE "AUDIT_LOGS" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "tenant_id" UUID NOT NULL REFERENCES "TENANTS"("id") ON DELETE CASCADE,
  "user_id" UUID REFERENCES "USERS"("id") ON DELETE SET NULL,
  "action" VARCHAR(100) NOT NULL,
  "entity_type" VARCHAR(100) NOT NULL,
  "old_data" JSONB,
  "new_data" JSONB,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);


CREATE TABLE "OUTLETS" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "tenant_id" UUID NOT NULL REFERENCES "TENANTS"("id") ON DELETE CASCADE,
  "name" VARCHAR(255) NOT NULL,
  "slug" VARCHAR(100) NOT NULL,
  "phone" VARCHAR(50),
  "address" TEXT,
  "is_public_order_enabled" BOOLEAN NOT NULL DEFAULT TRUE,
  "is_active" BOOLEAN NOT NULL DEFAULT TRUE,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  CONSTRAINT "unique_outlet_slug_per_tenant" UNIQUE ("tenant_id", "slug")
);

CREATE TRIGGER update_outlets_modtime 
  BEFORE UPDATE ON "OUTLETS" 
  FOR EACH ROW EXECUTE PROCEDURE update_timestamp_column();


CREATE TABLE "SERVICE_CATEGORIES" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "tenant_id" UUID NOT NULL REFERENCES "TENANTS"("id") ON DELETE CASCADE,
  "name" VARCHAR(255) NOT NULL,
  "description" TEXT,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TRIGGER update_service_categories_modtime 
  BEFORE UPDATE ON "SERVICE_CATEGORIES" 
  FOR EACH ROW EXECUTE PROCEDURE update_timestamp_column();


CREATE TABLE "DISCOUNTS" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "tenant_id" UUID NOT NULL REFERENCES "TENANTS"("id") ON DELETE CASCADE,
  "code" VARCHAR(50) NOT NULL,
  "name" VARCHAR(255) NOT NULL,
  "type" VARCHAR(20) NOT NULL CHECK ("type" IN ('percentage', 'fixed')),
  "value" NUMERIC(12, 2) NOT NULL DEFAULT 0,
  "minimum_transaction" NUMERIC(12, 2) DEFAULT 0,
  "start_at" TIMESTAMP WITH TIME ZONE,
  "end_at" TIMESTAMP WITH TIME ZONE,
  "is_active" BOOLEAN NOT NULL DEFAULT TRUE,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  CONSTRAINT "unique_discount_code_per_tenant" UNIQUE ("tenant_id", "code")
);

CREATE TRIGGER update_discounts_modtime 
  BEFORE UPDATE ON "DISCOUNTS" 
  FOR EACH ROW EXECUTE PROCEDURE update_timestamp_column();

-- -----------------------------------------------------------------------------
-- LEVEL 3 & 4: OPERATIONAL, SERVICES & PAYMENT METHODS
-- -----------------------------------------------------------------------------

CREATE TABLE "OUTLET_USERS" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "outlet_id" UUID NOT NULL REFERENCES "OUTLETS"("id") ON DELETE CASCADE,
  "user_id" UUID NOT NULL REFERENCES "USERS"("id") ON DELETE CASCADE,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  CONSTRAINT "unique_user_per_outlet" UNIQUE ("outlet_id", "user_id")
);


CREATE TABLE "SERVICES" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "tenant_id" UUID NOT NULL REFERENCES "TENANTS"("id") ON DELETE CASCADE,
  "category_id" UUID REFERENCES "SERVICE_CATEGORIES"("id") ON DELETE SET NULL,
  "name" VARCHAR(255) NOT NULL,
  "price" NUMERIC(12, 2) NOT NULL DEFAULT 0,
  "min_quantity" NUMERIC(8, 2) NOT NULL DEFAULT 1.0,
  "unit" VARCHAR(20) NOT NULL DEFAULT 'kg',
  "description" TEXT,
  "is_active" BOOLEAN NOT NULL DEFAULT TRUE,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TRIGGER update_services_modtime 
  BEFORE UPDATE ON "SERVICES" 
  FOR EACH ROW EXECUTE PROCEDURE update_timestamp_column();


CREATE TABLE "OUTLET_PAYMENT_METHODS" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "outlet_id" UUID NOT NULL REFERENCES "OUTLETS"("id") ON DELETE CASCADE,
  "type" VARCHAR(50) NOT NULL CHECK ("type" IN ('qris', 'bank_transfer', 'ewallet')),
  "provider_name" VARCHAR(100) NOT NULL,
  "account_number" VARCHAR(100),
  "account_name" VARCHAR(255) NOT NULL,
  "qr_image_url" VARCHAR(500),
  "is_active" BOOLEAN NOT NULL DEFAULT TRUE,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TRIGGER update_outlet_payment_methods_modtime 
  BEFORE UPDATE ON "OUTLET_PAYMENT_METHODS" 
  FOR EACH ROW EXECUTE PROCEDURE update_timestamp_column();

-- -----------------------------------------------------------------------------
-- LEVEL 5 TO 8: ORDERS, PAYMENTS, NOTIFICATIONS
-- -----------------------------------------------------------------------------

CREATE TABLE "ORDERS" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "tenant_id" UUID NOT NULL REFERENCES "TENANTS"("id") ON DELETE CASCADE,
  "outlet_id" UUID NOT NULL REFERENCES "OUTLETS"("id") ON DELETE RESTRICT,
  "customer_id" UUID NOT NULL REFERENCES "CUSTOMERS"("id") ON DELETE RESTRICT,
  "discount_id" UUID REFERENCES "DISCOUNTS"("id") ON DELETE SET NULL,
  "created_by" UUID REFERENCES "USERS"("id") ON DELETE SET NULL,
  "order_number" VARCHAR(100) NOT NULL,
  "tracking_code" VARCHAR(50) NOT NULL,
  "tracking_token" VARCHAR(255) NOT NULL,
  "status" VARCHAR(50) NOT NULL DEFAULT 'MENUNGGU_KONFIRMASI',
  "subtotal" NUMERIC(12, 2) NOT NULL DEFAULT 0,
  "discount_amount" NUMERIC(12, 2) NOT NULL DEFAULT 0,
  "total" NUMERIC(12, 2) NOT NULL DEFAULT 0,
  "paid_amount" NUMERIC(12, 2) NOT NULL DEFAULT 0,
  "remaining_amount" NUMERIC(12, 2) NOT NULL DEFAULT 0,
  "notes" TEXT,
  "estimated_done" TIMESTAMP WITH TIME ZONE,
  "completed_at" TIMESTAMP WITH TIME ZONE,
  "canceled_at" TIMESTAMP WITH TIME ZONE,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  CONSTRAINT "unique_order_number_per_tenant" UNIQUE ("tenant_id", "order_number")
);

CREATE TRIGGER update_orders_modtime 
  BEFORE UPDATE ON "ORDERS" 
  FOR EACH ROW EXECUTE PROCEDURE update_timestamp_column();


CREATE TABLE "ORDER_ITEMS" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "order_id" UUID NOT NULL REFERENCES "ORDERS"("id") ON DELETE CASCADE,
  "service_id" UUID REFERENCES "SERVICES"("id") ON DELETE RESTRICT,
  "price" NUMERIC(12, 2) NOT NULL DEFAULT 0,
  "quantity" NUMERIC(8, 2) NOT NULL DEFAULT 1.0,
  "subtotal" NUMERIC(12, 2) NOT NULL DEFAULT 0,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TRIGGER update_order_items_modtime 
  BEFORE UPDATE ON "ORDER_ITEMS" 
  FOR EACH ROW EXECUTE PROCEDURE update_timestamp_column();


CREATE TABLE "ORDER_STATUS_HISTORIES" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "order_id" UUID NOT NULL REFERENCES "ORDERS"("id") ON DELETE CASCADE,
  "status" VARCHAR(50) NOT NULL,
  "changed_by" UUID REFERENCES "USERS"("id") ON DELETE SET NULL,
  "notes" TEXT,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);


CREATE TABLE "PAYMENTS" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "tenant_id" UUID NOT NULL REFERENCES "TENANTS"("id") ON DELETE CASCADE,
  "outlet_id" UUID NOT NULL REFERENCES "OUTLETS"("id") ON DELETE RESTRICT,
  "order_id" UUID NOT NULL REFERENCES "ORDERS"("id") ON DELETE CASCADE,
  "received_by" UUID REFERENCES "USERS"("id") ON DELETE SET NULL,
  "amount" NUMERIC(12, 2) NOT NULL DEFAULT 0,
  "payment_method" VARCHAR(50) NOT NULL CHECK ("payment_method" IN ('cash', 'qris', 'bank_transfer')),
  "status" VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK ("status" IN ('pending', 'verified', 'rejected')),
  "reference_number" VARCHAR(100),
  "proof_image_url" VARCHAR(500),
  "notes" TEXT,
  "paid_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  "verified_at" TIMESTAMP WITH TIME ZONE,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TRIGGER update_payments_modtime 
  BEFORE UPDATE ON "PAYMENTS" 
  FOR EACH ROW EXECUTE PROCEDURE update_timestamp_column();


CREATE TABLE "NOTIFICATIONS" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "tenant_id" UUID NOT NULL REFERENCES "TENANTS"("id") ON DELETE CASCADE,
  "order_id" UUID REFERENCES "ORDERS"("id") ON DELETE CASCADE,
  "user_id" UUID REFERENCES "USERS"("id") ON DELETE CASCADE,
  "type" VARCHAR(50) NOT NULL,
  "message" TEXT NOT NULL,
  "sent_at" TIMESTAMP WITH TIME ZONE,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexing untuk Performa Multi-Tenant
CREATE INDEX "idx_customers_tenant" ON "CUSTOMERS" ("tenant_id");
CREATE INDEX "idx_users_tenant" ON "USERS" ("tenant_id");
CREATE INDEX "idx_outlets_tenant" ON "OUTLETS" ("tenant_id");
CREATE INDEX "idx_services_tenant" ON "SERVICES" ("tenant_id");
CREATE INDEX "idx_orders_tenant" ON "ORDERS" ("tenant_id");
CREATE INDEX "idx_payments_tenant" ON "PAYMENTS" ("tenant_id");
CREATE INDEX "idx_orders_outlet" ON "ORDERS" ("outlet_id");
CREATE INDEX "idx_orders_customer" ON "ORDERS" ("customer_id");
CREATE INDEX "idx_orders_tracking" ON "ORDERS" ("tracking_code");
CREATE INDEX "idx_order_items_order" ON "ORDER_ITEMS" ("order_id");
CREATE INDEX "idx_payments_order" ON "PAYMENTS" ("order_id");
CREATE INDEX "idx_order_status_histories_order" ON "ORDER_STATUS_HISTORIES" ("order_id");