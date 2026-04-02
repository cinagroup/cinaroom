-- CinaSeek Initial Schema Migration
-- Generated: 2026-04-02

BEGIN;

-- Enable UUID extension (optional, for future use)
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================
-- Users
-- ============================================================
CREATE TABLE IF NOT EXISTS users (
    id                  BIGSERIAL PRIMARY KEY,
    cinatoken_id        BIGINT        NOT NULL,
    username            VARCHAR(20)   NOT NULL,
    email               VARCHAR(100)  NOT NULL,
    password            VARCHAR(255)  DEFAULT '',
    nickname            VARCHAR(50)   DEFAULT '',
    phone               VARCHAR(20)   DEFAULT '',
    avatar              VARCHAR(255)  DEFAULT '',
    provider            VARCHAR(50)   DEFAULT '',
    active              BOOLEAN       NOT NULL DEFAULT TRUE,
    two_factor_enabled  BOOLEAN       NOT NULL DEFAULT FALSE,
    last_login_at       TIMESTAMPTZ,
    created_at          TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ   NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_users_cinatoken_id UNIQUE (cinatoken_id),
    CONSTRAINT uq_users_username     UNIQUE (username),
    CONSTRAINT uq_users_email        UNIQUE (email)
);

CREATE INDEX idx_users_cinatoken_id ON users (cinatoken_id);
CREATE INDEX idx_users_email        ON users (email);
CREATE INDEX idx_users_username     ON users (username);

-- ============================================================
-- Virtual Machines
-- ============================================================
CREATE TABLE IF NOT EXISTS vms (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name          VARCHAR(100)  NOT NULL,
    status        VARCHAR(20)   NOT NULL DEFAULT 'stopped',
    ip            VARCHAR(50)   DEFAULT '',
    image         VARCHAR(50)   NOT NULL,
    cpu           INT           NOT NULL DEFAULT 1,
    memory        INT           NOT NULL DEFAULT 1,
    disk          INT           NOT NULL DEFAULT 10,
    network_type  VARCHAR(20)   NOT NULL DEFAULT 'nat',
    ssh_key       TEXT          DEFAULT '',
    init_script   TEXT          DEFAULT '',
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vms_user_id ON vms (user_id);
CREATE INDEX idx_vms_status  ON vms (status);
CREATE INDEX idx_vms_name    ON vms (name);

-- ============================================================
-- VM Snapshots
-- ============================================================
CREATE TABLE IF NOT EXISTS vm_snapshots (
    id         BIGSERIAL PRIMARY KEY,
    vm_id      BIGINT        NOT NULL REFERENCES vms(id) ON DELETE CASCADE,
    name       VARCHAR(100)  NOT NULL,
    size       BIGINT        NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vm_snapshots_vm_id ON vm_snapshots (vm_id);

-- ============================================================
-- VM Operation Logs
-- ============================================================
CREATE TABLE IF NOT EXISTS vm_logs (
    id         BIGSERIAL PRIMARY KEY,
    vm_id      BIGINT        NOT NULL REFERENCES vms(id) ON DELETE CASCADE,
    operation  VARCHAR(50)   NOT NULL,
    result     VARCHAR(20)   NOT NULL,
    message    TEXT          DEFAULT '',
    created_at TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vm_logs_vm_id      ON vm_logs (vm_id);
CREATE INDEX idx_vm_logs_created_at ON vm_logs (created_at);

-- ============================================================
-- Directory Mounts
-- ============================================================
CREATE TABLE IF NOT EXISTS mounts (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    vm_id       BIGINT        NOT NULL REFERENCES vms(id) ON DELETE CASCADE,
    name        VARCHAR(100)  NOT NULL,
    host_path   VARCHAR(500)  NOT NULL,
    vm_path     VARCHAR(500)  NOT NULL,
    status      VARCHAR(20)   NOT NULL DEFAULT 'unmounted',
    permission  VARCHAR(10)   NOT NULL DEFAULT 'rw',
    auto_mount  BOOLEAN       NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_mounts_user_id ON mounts (user_id);
CREATE INDEX idx_mounts_vm_id   ON mounts (vm_id);

-- ============================================================
-- OpenClaw Configurations
-- ============================================================
CREATE TABLE IF NOT EXISTS openclaw_configs (
    id                 BIGSERIAL PRIMARY KEY,
    vm_id              BIGINT        NOT NULL REFERENCES vms(id) ON DELETE CASCADE,
    status             VARCHAR(20)   NOT NULL DEFAULT 'not_installed',
    version            VARCHAR(20)   DEFAULT '',
    running_time       BIGINT        NOT NULL DEFAULT 0,
    default_model      VARCHAR(100)  DEFAULT '',
    api_key            VARCHAR(255)  DEFAULT '',
    enabled_tools      TEXT          DEFAULT '',
    enabled_skills     TEXT          DEFAULT '',
    workspace_path     VARCHAR(500)  DEFAULT '',
    skills_path        VARCHAR(500)  DEFAULT '',
    sync_openclaw_json BOOLEAN       NOT NULL DEFAULT TRUE,
    sync_tool_configs  BOOLEAN       NOT NULL DEFAULT TRUE,
    last_deployed_at   TIMESTAMPTZ,
    created_at         TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_openclaw_configs_vm_id ON openclaw_configs (vm_id);

-- ============================================================
-- Remote Access
-- ============================================================
CREATE TABLE IF NOT EXISTS remote_access (
    id             BIGSERIAL PRIMARY KEY,
    vm_id          BIGINT        NOT NULL REFERENCES vms(id) ON DELETE CASCADE,
    enabled        BOOLEAN       NOT NULL DEFAULT FALSE,
    access_address VARCHAR(255)  DEFAULT '',
    qr_code        VARCHAR(500)  DEFAULT '',
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_remote_access_vm_id UNIQUE (vm_id)
);

-- ============================================================
-- IP Whitelists
-- ============================================================
CREATE TABLE IF NOT EXISTS ip_whitelists (
    id         BIGSERIAL PRIMARY KEY,
    vm_id      BIGINT        NOT NULL REFERENCES vms(id) ON DELETE CASCADE,
    ip         VARCHAR(50)   NOT NULL,
    is_cidr    BOOLEAN       NOT NULL DEFAULT FALSE,
    note       VARCHAR(200)  DEFAULT '',
    created_at TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ip_whitelists_vm_id ON ip_whitelists (vm_id);

-- ============================================================
-- Remote Access Logs
-- ============================================================
CREATE TABLE IF NOT EXISTS remote_logs (
    id             BIGSERIAL PRIMARY KEY,
    vm_id          BIGINT        NOT NULL REFERENCES vms(id) ON DELETE CASCADE,
    access_time    TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    access_ip      VARCHAR(50)   NOT NULL,
    access_path    VARCHAR(500)  DEFAULT '',
    user_agent     VARCHAR(500)  DEFAULT '',
    response_code  INT           NOT NULL DEFAULT 0
);

CREATE INDEX idx_remote_logs_vm_id       ON remote_logs (vm_id);
CREATE INDEX idx_remote_logs_access_time ON remote_logs (access_time);

-- ============================================================
-- Login Logs
-- ============================================================
CREATE TABLE IF NOT EXISTS login_logs (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    login_time TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    ip         VARCHAR(50)   NOT NULL,
    location   VARCHAR(200)  DEFAULT '',
    device     VARCHAR(200)  DEFAULT ''
);

CREATE INDEX idx_login_logs_user_id    ON login_logs (user_id);
CREATE INDEX idx_login_logs_login_time ON login_logs (login_time);

-- ============================================================
-- Sessions
-- ============================================================
CREATE TABLE IF NOT EXISTS sessions (
    id             BIGSERIAL PRIMARY KEY,
    user_id        BIGINT        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token          VARCHAR(255)  NOT NULL,
    device         VARCHAR(200)  DEFAULT '',
    location       VARCHAR(200)  DEFAULT '',
    login_time     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    last_active_at TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    expired_at     TIMESTAMPTZ   NOT NULL
);

CREATE INDEX idx_sessions_user_id ON sessions (user_id);
CREATE UNIQUE INDEX uq_sessions_token ON sessions (token);

-- ============================================================
-- System Settings
-- ============================================================
CREATE TABLE IF NOT EXISTS system_settings (
    id        BIGSERIAL PRIMARY KEY,
    key       VARCHAR(100)  NOT NULL,
    value     TEXT          DEFAULT '',
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_system_settings_key UNIQUE (key)
);

-- ============================================================
-- VM Metrics
-- ============================================================
CREATE TABLE IF NOT EXISTS vm_metrics (
    id           BIGSERIAL PRIMARY KEY,
    vm_id        BIGINT        NOT NULL REFERENCES vms(id) ON DELETE CASCADE,
    cpu_usage    DOUBLE PRECISION NOT NULL DEFAULT 0,
    memory_usage DOUBLE PRECISION NOT NULL DEFAULT 0,
    disk_io      DOUBLE PRECISION NOT NULL DEFAULT 0,
    network_rx   DOUBLE PRECISION NOT NULL DEFAULT 0,
    network_tx   DOUBLE PRECISION NOT NULL DEFAULT 0,
    timestamp    TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vm_metrics_vm_id     ON vm_metrics (vm_id);
CREATE INDEX idx_vm_metrics_timestamp ON vm_metrics (timestamp);

-- ============================================================
-- Seed: default system settings
-- ============================================================
INSERT INTO system_settings (key, value) VALUES
    ('system.name',       'CinaSeek'),
    ('system.version',    '1.0.0'),
    ('maintenance.mode',  'false'),
    ('default_vm_image',  'ubuntu:22.04'),
    ('max_vms_per_user',  '10'),
    ('openclaw.version',  'latest')
ON CONFLICT (key) DO NOTHING;

COMMIT;
