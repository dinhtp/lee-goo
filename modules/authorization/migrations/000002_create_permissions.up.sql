CREATE TABLE IF NOT EXISTS permissions (
    id VARCHAR(36) PRIMARY KEY,
    action VARCHAR(100) NOT NULL,
    resource VARCHAR(100) NOT NULL,
    UNIQUE (action, resource)
);
