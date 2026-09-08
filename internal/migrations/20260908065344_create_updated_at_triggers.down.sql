DROP TRIGGER IF EXISTS update_tenants_modtime ON tenants;
DROP TRIGGER IF EXISTS update_customers_modtime ON customers;
DROP TRIGGER IF EXISTS update_users_modtime ON users;
DROP TRIGGER IF EXISTS update_outlets_modtime ON outlets;
DROP TRIGGER IF EXISTS update_service_categories_modtime ON service_categories;
DROP TRIGGER IF EXISTS update_discounts_modtime ON discounts;
DROP TRIGGER IF EXISTS update_services_modtime ON services;
DROP TRIGGER IF EXISTS update_orders_modtime ON orders;
DROP TRIGGER IF EXISTS update_order_items_modtime ON order_items;
DROP TRIGGER IF EXISTS update_payments_modtime ON payments;
DROP TRIGGER IF EXISTS update_outlet_payment_methods_modtime ON outlet_payment_methods;

DROP FUNCTION IF EXISTS update_timestamp_column();