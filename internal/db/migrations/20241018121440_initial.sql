-- +goose Up
-- +goose StatementBegin

-- AUTH
CREATE TABLE auth_users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255),
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    phone_number VARCHAR(20) UNIQUE,
    telegram_id BIGINT UNIQUE, 
    oauth_id VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE auth_permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    codename VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE auth_user_permissions (
    user_id INT NOT NULL,
    permission_id INT NOT NULL,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, permission_id),
    FOREIGN KEY (user_id) REFERENCES auth_users (id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES auth_permissions (id) ON DELETE CASCADE
);


CREATE TABLE organizations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE,
    notes TEXT DEFAULT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE security_posts (
    id SERIAL PRIMARY KEY,
    organization_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    location TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organizations (id) ON DELETE CASCADE
);

CREATE TABLE access_passes (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    organization_id INT NOT NULL,
    security_post_id INT,
    qr_code TEXT UNIQUE,
    vehicle_number VARCHAR(50),
    purpose TEXT,
    valid_from TIMESTAMP NOT NULL,
    valid_until TIMESTAMP NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES auth_users (id) ON DELETE CASCADE,
    FOREIGN KEY (organization_id) REFERENCES organizations (id) ON DELETE CASCADE,
    FOREIGN KEY (security_post_id) REFERENCES security_posts (id) ON DELETE SET NULL
);

CREATE TABLE access_logs (
    id SERIAL PRIMARY KEY,
    user_id INT,
    security_post_id INT,
    pass_id INT,
    action_type VARCHAR(50) NOT NULL, -- e.g., "entry", "exit"
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES auth_users (id) ON DELETE CASCADE,
    FOREIGN KEY (security_post_id) REFERENCES security_posts (id) ON DELETE CASCADE,
    FOREIGN KEY (pass_id) REFERENCES access_passes (id) ON DELETE CASCADE
);

-- Triggers and Trigger Functions

-- Update `updated_at` timestamp on row update
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers for updating `updated_at` automatically
CREATE TRIGGER set_users_updated_at
BEFORE UPDATE ON auth_users
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER set_permissions_updated_at
BEFORE UPDATE ON auth_permissions
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER set_organizations_updated_at
BEFORE UPDATE ON organizations
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER set_security_posts_updated_at
BEFORE UPDATE ON security_posts
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER set_access_passes_updated_at
BEFORE UPDATE ON access_passes
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE OR REPLACE FUNCTION log_access_event()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO access_logs (user_id, security_post_id, pass_id, action_type, timestamp)
    VALUES (NEW.user_id, NEW.security_post_id, NEW.id, 'created', CURRENT_TIMESTAMP);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER log_access_pass_creation
AFTER INSERT ON access_passes
FOR EACH ROW
EXECUTE FUNCTION log_access_event();

CREATE OR REPLACE FUNCTION deactivate_expired_passes()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.valid_until < CURRENT_TIMESTAMP THEN
        NEW.is_active = FALSE;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_pass_expiry
BEFORE INSERT OR UPDATE ON access_passes
FOR EACH ROW
EXECUTE FUNCTION deactivate_expired_passes();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger set_users_updated_at on auth_users;
drop trigger set_permissions_updated_at on auth_permissions;
drop trigger set_organizations_updated_at on organizations;
drop trigger set_security_posts_updated_at on security_posts;
drop trigger set_access_passes_updated_at on access_passes;
drop trigger log_access_pass_creation on access_passes;
drop trigger check_pass_expiry on access_passes;

drop function update_timestamp;
drop function log_access_event;
drop function deactivate_expired_passes;

drop table auth_users cascade;
drop table auth_permissions cascade;
drop table auth_user_permissions cascade;
drop table organizations cascade;
drop table security_posts cascade;
drop table access_passes cascade;
drop table access_logs cascade;
-- +goose StatementEnd
