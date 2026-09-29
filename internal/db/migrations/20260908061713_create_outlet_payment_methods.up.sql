CREATE TABLE outlet_payment_methods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    outlet_id UUID NOT NULL REFERENCES outlets(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('qris', 'bank_transfer')),
    provider_name VARCHAR(100) NOT NULL,
    account_number VARCHAR(50),
    account_name VARCHAR(150) NOT NULL,
    qr_image_url VARCHAR(500),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

