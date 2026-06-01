CREATE TABLE IF NOT EXISTS modules (
    name                    VARCHAR(100) PRIMARY KEY,
    version                 VARCHAR(50)  NOT NULL,
    status                  VARCHAR(30)  NOT NULL,
    path                    VARCHAR(255) NOT NULL,
    checksum                VARCHAR(255),
    installed_at            TIMESTAMP,
    enabled_at              TIMESTAMP,
    disabled_at             TIMESTAMP,
    upgraded_at             TIMESTAMP,
    uninstalled_at          TIMESTAMP,
    removed_from_codebase_at TIMESTAMP,
    created_at              TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMP    NOT NULL DEFAULT NOW()
);
