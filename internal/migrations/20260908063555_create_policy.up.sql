CREATE POLICY rls_tenants ON tenants FOR ALL
    USING (id = current_tenant_id())
    WITH CHECK (id = current_tenant_id());

CREATE POLICY rls_customers ON customers FOR ALL
    USING (tenant_id = current_tenant_id())
    WITH CHECK (tenant_id = current_tenant_id());

CREATE POLICY rls_users ON users FOR ALL
    USING (tenant_id = current_tenant_id())
    WITH CHECK (tenant_id = current_tenant_id());

CREATE POLICY rls_outlets ON outlets FOR ALL
    USING (tenant_id = current_tenant_id())
    WITH CHECK (tenant_id = current_tenant_id());

CREATE POLICY rls_service_categories ON service_categories FOR ALL
    USING (tenant_id = current_tenant_id())
    WITH CHECK (tenant_id = current_tenant_id());

CREATE POLICY rls_discounts ON discounts FOR ALL
    USING (tenant_id = current_tenant_id())
    WITH CHECK (tenant_id = current_tenant_id());

CREATE POLICY rls_services ON services FOR ALL
    USING (tenant_id = current_tenant_id())
    WITH CHECK (tenant_id = current_tenant_id());

CREATE POLICY rls_orders ON orders FOR ALL
    USING (tenant_id = current_tenant_id())
    WITH CHECK (tenant_id = current_tenant_id());

CREATE POLICY rls_audit_logs ON audit_logs FOR ALL
    USING (tenant_id = current_tenant_id())
    WITH CHECK (tenant_id = current_tenant_id());

CREATE POLICY rls_payments ON payments FOR ALL
    USING (
        EXISTS (
            SELECT 1 FROM orders
            WHERE orders.id = payments.order_id
              AND orders.tenant_id = current_tenant_id()
        )
    )
    WITH CHECK (
        EXISTS (
            SELECT 1 FROM orders
            WHERE orders.id = payments.order_id
              AND orders.tenant_id = current_tenant_id()
        )
    );

CREATE POLICY rls_notifications ON notifications FOR ALL
    USING (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = notifications.user_id
              AND users.tenant_id = current_tenant_id()
        )
    )
    WITH CHECK (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = notifications.user_id
              AND users.tenant_id = current_tenant_id()
        )
    );

CREATE POLICY rls_order_items ON order_items FOR ALL
    USING (
        EXISTS (
            SELECT 1 FROM orders
            WHERE orders.id = order_items.order_id
              AND orders.tenant_id = current_tenant_id()
        )
    )
    WITH CHECK (
        EXISTS (
            SELECT 1 FROM orders
            WHERE orders.id = order_items.order_id
              AND orders.tenant_id = current_tenant_id()
        )
    );

CREATE POLICY rls_order_status_histories ON order_status_histories FOR ALL
    USING (
        EXISTS (
            SELECT 1 FROM orders
            WHERE orders.id = order_status_histories.order_id
              AND orders.tenant_id = current_tenant_id()
        )
    )
    WITH CHECK (
        EXISTS (
            SELECT 1 FROM orders
            WHERE orders.id = order_status_histories.order_id
              AND orders.tenant_id = current_tenant_id()
        )
    );

CREATE POLICY rls_outlet_payment_methods ON outlet_payment_methods FOR ALL
    USING (
        EXISTS (
            SELECT 1 FROM outlets
            WHERE outlets.id = outlet_payment_methods.outlet_id
              AND outlets.tenant_id = current_tenant_id()
        )
    )
    WITH CHECK (
        EXISTS (
            SELECT 1 FROM outlets
            WHERE outlets.id = outlet_payment_methods.outlet_id
              AND outlets.tenant_id = current_tenant_id()
        )
    );

CREATE POLICY rls_outlet_users ON outlet_users FOR ALL
    USING (
        EXISTS (
            SELECT 1 FROM outlets
            WHERE outlets.id = outlet_users.outlet_id
              AND outlets.tenant_id = current_tenant_id()
        )
    )
    WITH CHECK (
        EXISTS (
            SELECT 1 FROM outlets
            WHERE outlets.id = outlet_users.outlet_id
              AND outlets.tenant_id = current_tenant_id()
        )
    );