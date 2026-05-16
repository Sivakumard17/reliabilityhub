CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE clusters (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(255) NOT NULL UNIQUE,
    environment VARCHAR(50)  NOT NULL DEFAULT 'production',
    endpoint    TEXT,
    status      VARCHAR(50)  NOT NULL DEFAULT 'unknown',
    metadata    JSONB        NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE incidents (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title        VARCHAR(500) NOT NULL,
    description  TEXT,
    severity     VARCHAR(20)  NOT NULL DEFAULT 'medium',
    status       VARCHAR(50)  NOT NULL DEFAULT 'open',
    cluster_id   UUID REFERENCES clusters(id) ON DELETE SET NULL,
    service      VARCHAR(255),
    alert_name   VARCHAR(255),
    labels       JSONB        NOT NULL DEFAULT '{}',
    annotations  JSONB        NOT NULL DEFAULT '{}',
    started_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    resolved_at  TIMESTAMPTZ,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_severity CHECK (
        severity IN ('critical', 'high', 'medium', 'low', 'info')
    ),
    CONSTRAINT valid_status CHECK (
        status IN ('open', 'investigating', 'identified', 'monitoring', 'resolved')
    )
);

CREATE TABLE slos (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    cluster_id   UUID REFERENCES clusters(id) ON DELETE CASCADE,
    service      VARCHAR(255) NOT NULL,
    slo_type     VARCHAR(50)  NOT NULL DEFAULT 'availability',
    target       NUMERIC(5,4) NOT NULL,
    window_days  INTEGER      NOT NULL DEFAULT 30,
    promql_good  TEXT,
    promql_total TEXT,
    is_active    BOOLEAN      NOT NULL DEFAULT true,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_target CHECK (target > 0 AND target <= 1),
    CONSTRAINT valid_slo_type CHECK (
        slo_type IN ('availability', 'latency', 'error_rate', 'throughput')
    )
);

CREATE TABLE slo_snapshots (
    id                 UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    slo_id             UUID NOT NULL REFERENCES slos(id) ON DELETE CASCADE,
    compliance         NUMERIC(5,4) NOT NULL,
    error_budget_total NUMERIC(10,6),
    error_budget_used  NUMERIC(10,6),
    error_budget_pct   NUMERIC(5,4),
    burn_rate_1h       NUMERIC(8,4),
    burn_rate_6h       NUMERIC(8,4),
    burn_rate_24h      NUMERIC(8,4),
    snapshot_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE remediation_policies (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    cluster_id   UUID REFERENCES clusters(id) ON DELETE CASCADE,
    trigger_type VARCHAR(50)  NOT NULL,
    conditions   JSONB        NOT NULL DEFAULT '{}',
    actions      JSONB        NOT NULL DEFAULT '[]',
    is_active    BOOLEAN      NOT NULL DEFAULT true,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE audit_logs (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    actor       VARCHAR(255) NOT NULL DEFAULT 'system',
    action      VARCHAR(255) NOT NULL,
    resource    VARCHAR(255) NOT NULL,
    resource_id UUID,
    details     JSONB        NOT NULL DEFAULT '{}',
    ip_address  INET,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_incidents_status     ON incidents(status);
CREATE INDEX idx_incidents_severity   ON incidents(severity);
CREATE INDEX idx_incidents_started_at ON incidents(started_at DESC);
CREATE INDEX idx_incidents_service    ON incidents(service);
CREATE INDEX idx_incidents_cluster_id ON incidents(cluster_id);
CREATE INDEX idx_slos_cluster_id      ON slos(cluster_id);
CREATE INDEX idx_slos_is_active       ON slos(is_active);
CREATE INDEX idx_slo_snapshots_slo_id ON slo_snapshots(slo_id);
CREATE INDEX idx_slo_snapshots_at     ON slo_snapshots(snapshot_at DESC);
CREATE INDEX idx_audit_logs_resource  ON audit_logs(resource, resource_id);
CREATE INDEX idx_audit_logs_at        ON audit_logs(created_at DESC);

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_clusters
    BEFORE UPDATE ON clusters
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER set_updated_at_incidents
    BEFORE UPDATE ON incidents
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER set_updated_at_slos
    BEFORE UPDATE ON slos
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
