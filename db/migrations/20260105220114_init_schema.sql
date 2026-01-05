-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE services (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    identifier TEXT NOT NULL UNIQUE,
    image_url TEXT,
    registered_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE service_access_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    access_key TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    rotated_from UUID REFERENCES service_access_settings(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE service_subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    service_id UUID NOT NULL UNIQUE REFERENCES services(id) ON DELETE CASCADE,
    plan_name TEXT NOT NULL,
    monthly_email_quota INTEGER NOT NULL,
    monthly_storage_quota BIGINT NOT NULL,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE service_usages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    service_id UUID NOT NULL UNIQUE REFERENCES services(id) ON DELETE CASCADE,
    emails_sent INTEGER DEFAULT 0,
    storage_used BIGINT DEFAULT 0,
    period_start TIMESTAMP NOT NULL DEFAULT NOW(),
    period_end TIMESTAMP NOT NULL DEFAULT NOW() + INTERVAL '30 days',
    last_reset TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE email_tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "to" TEXT NOT NULL,
    "from" TEXT NOT NULL,
    subject TEXT NOT NULL,
    body TEXT NOT NULL,
    cc TEXT[],
    bcc TEXT[],
    reply_to VARCHAR(255),
    status TEXT NOT NULL DEFAULT 'pending',
    retry_count INTEGER DEFAULT 0,
    last_attempt TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE attachments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email_task_id UUID NOT NULL REFERENCES email_tasks(id) ON DELETE CASCADE,
    filename TEXT NOT NULL,
    content_type TEXT NOT NULL,
    data TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS attachments CASCADE;
DROP TABLE IF EXISTS email_tasks CASCADE;
DROP TABLE IF EXISTS service_usages CASCADE;
DROP TABLE IF EXISTS service_subscriptions CASCADE;
DROP TABLE IF EXISTS service_access_settings CASCADE;
DROP TABLE IF EXISTS services CASCADE;
