CREATE TABLE IF NOT EXISTS roles (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- user_id is an opaque string; no FK to users table to keep modules decoupled.
CREATE TABLE IF NOT EXISTS user_roles (
    user_id VARCHAR(36) NOT NULL,
    role_id VARCHAR(36) NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);
