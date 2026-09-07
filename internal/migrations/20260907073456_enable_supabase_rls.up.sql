-- ==============================================================================
-- 1. HELPER FUNCTION (DIOPTIMALKAN DENGAN PURE SQL & INLINING)
-- ==============================================================================
CREATE OR REPLACE FUNCTION current_tenant_id()
RETURNS UUID AS $$
  SELECT NULLIF(current_setting('app.current_tenant_id', true), '')::UUID;
$$ LANGUAGE sql STABLE PARALLEL SAFE;


-- ==============================================================================
-- 2. ENABLE ROW LEVEL SECURITY
-- ==============================================================================
ALTER TABLE "TENANTS" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "CUSTOMERS" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "USERS" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "OUTLETS" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "SERVICE_CATEGORIES" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "DISCOUNTS" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "SERVICES" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "ORDERS" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "PAYMENTS" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "NOTIFICATIONS" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "AUDIT_LOGS" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "ORDER_ITEMS" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "ORDER_STATUS_HISTORIES" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "OUTLET_PAYMENT_METHODS" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "OUTLET_USERS" ENABLE ROW LEVEL SECURITY;


-- ==============================================================================
-- 3. POLICY UNTUK TABEL UTAMA (DIRECT TENANT FILTER + WITH CHECK)
-- ==============================================================================
CREATE POLICY rls_tenants ON "TENANTS" FOR ALL 
  USING (id = current_tenant_id()) 
  WITH CHECK (id = current_tenant_id());

CREATE POLICY rls_customers ON "CUSTOMERS" FOR ALL 
  USING (tenant_id = current_tenant_id()) 
  WITH CHECK (tenant_id = current_tenant_id());

CREATE POLICY rls_users ON "USERS" FOR ALL 
  USING (tenant_id = current_tenant_id()) 
  WITH CHECK (tenant_id = current_tenant_id());

CREATE POLICY rls_outlets ON "OUTLETS" FOR ALL 
  USING (tenant_id = current_tenant_id()) 
  WITH CHECK (tenant_id = current_tenant_id());

CREATE POLICY rls_service_categories ON "SERVICE_CATEGORIES" FOR ALL 
  USING (tenant_id = current_tenant_id()) 
  WITH CHECK (tenant_id = current_tenant_id());

CREATE POLICY rls_discounts ON "DISCOUNTS" FOR ALL 
  USING (tenant_id = current_tenant_id()) 
  WITH CHECK (tenant_id = current_tenant_id());

CREATE POLICY rls_services ON "SERVICES" FOR ALL 
  USING (tenant_id = current_tenant_id()) 
  WITH CHECK (tenant_id = current_tenant_id());

CREATE POLICY rls_orders ON "ORDERS" FOR ALL 
  USING (tenant_id = current_tenant_id()) 
  WITH CHECK (tenant_id = current_tenant_id());

CREATE POLICY rls_payments ON "PAYMENTS" FOR ALL 
  USING (tenant_id = current_tenant_id()) 
  WITH CHECK (tenant_id = current_tenant_id());

CREATE POLICY rls_notifications ON "NOTIFICATIONS" FOR ALL 
  USING (tenant_id = current_tenant_id()) 
  WITH CHECK (tenant_id = current_tenant_id());

CREATE POLICY rls_audit_logs ON "AUDIT_LOGS" FOR ALL 
  USING (tenant_id = current_tenant_id()) 
  WITH CHECK (tenant_id = current_tenant_id());


-- ==============================================================================
-- 4. POLICY UNTUK CHILD TABLES (SESUAI DESAIN QUERY EXISTS ANDA)
-- ==============================================================================
CREATE POLICY rls_order_items ON "ORDER_ITEMS" FOR ALL 
  USING (
    EXISTS (
      SELECT 1 FROM "ORDERS" 
      WHERE "ORDERS".id = "ORDER_ITEMS".order_id 
        AND "ORDERS".tenant_id = current_tenant_id()
    )
  )
  WITH CHECK (
    EXISTS (
      SELECT 1 FROM "ORDERS" 
      WHERE "ORDERS".id = "ORDER_ITEMS".order_id 
        AND "ORDERS".tenant_id = current_tenant_id()
    )
  );

CREATE POLICY rls_order_status_histories ON "ORDER_STATUS_HISTORIES" FOR ALL 
  USING (
    EXISTS (
      SELECT 1 FROM "ORDERS" 
      WHERE "ORDERS".id = "ORDER_STATUS_HISTORIES".order_id 
        AND "ORDERS".tenant_id = current_tenant_id()
    )
  )
  WITH CHECK (
    EXISTS (
      SELECT 1 FROM "ORDERS" 
      WHERE "ORDERS".id = "ORDER_STATUS_HISTORIES".order_id 
        AND "ORDERS".tenant_id = current_tenant_id()
    )
  );

CREATE POLICY rls_outlet_payment_methods ON "OUTLET_PAYMENT_METHODS" FOR ALL 
  USING (
    EXISTS (
      SELECT 1 FROM "OUTLETS" 
      WHERE "OUTLETS".id = "OUTLET_PAYMENT_METHODS".outlet_id 
        AND "OUTLETS".tenant_id = current_tenant_id()
    )
  )
  WITH CHECK (
    EXISTS (
      SELECT 1 FROM "OUTLETS" 
      WHERE "OUTLETS".id = "OUTLET_PAYMENT_METHODS".outlet_id 
        AND "OUTLETS".tenant_id = current_tenant_id()
    )
  );

CREATE POLICY rls_outlet_users ON "OUTLET_USERS" FOR ALL 
  USING (
    EXISTS (
      SELECT 1 FROM "OUTLETS" 
      WHERE "OUTLETS".id = "OUTLET_USERS".outlet_id 
        AND "OUTLETS".tenant_id = current_tenant_id()
    )
  )
  WITH CHECK (
    EXISTS (
      SELECT 1 FROM "OUTLETS" 
      WHERE "OUTLETS".id = "OUTLET_USERS".outlet_id 
        AND "OUTLETS".tenant_id = current_tenant_id()
    )
  );