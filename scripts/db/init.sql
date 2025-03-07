CREATE SCHEMA IF NOT EXISTS synapse;

SET
search_path TO synapse;

-- Kubernetes resources
DROP TABLE IF EXISTS kubernetes_resources;
CREATE TABLE kubernetes_resources
(
    id                       BIGSERIAL PRIMARY KEY,
    external_id              VARCHAR(255) NOT NULL UNIQUE,
    namespace                VARCHAR(255) NOT NULL,
    deployment_name          VARCHAR(255),
    deployment_replicas      INT,
    service_name             VARCHAR(255),
    ingress_jupyter_name     VARCHAR(255),
    ingress_jupyter_protocol VARCHAR(10),
    ingress_jupyter_cname    VARCHAR(255),
    ingress_jupyter_url      VARCHAR(255),
    ingress_ray_name         VARCHAR(255),
    ingress_ray_protocol     VARCHAR(10),
    ingress_ray_cname        VARCHAR(255),
    ingress_ray_url          VARCHAR(255),
    status                   INT,
    created_at               TIMESTAMP(3) NOT NULL,
    updated_at               TIMESTAMP(3),
    UNIQUE (external_id)
);

COMMENT
ON TABLE kubernetes_resources IS 'Kubernetes Resources';
COMMENT
ON COLUMN kubernetes_resources.external_id IS 'External Business Primary Key';
COMMENT
ON COLUMN kubernetes_resources.namespace IS 'Kubernetes Namespace';
COMMENT
ON COLUMN kubernetes_resources.deployment_name IS 'Deployment Name';
COMMENT
ON COLUMN kubernetes_resources.deployment_replicas IS 'Deployment Replica';
COMMENT
ON COLUMN kubernetes_resources.service_name IS 'Service Name';
COMMENT
ON COLUMN kubernetes_resources.ingress_jupyter_name IS 'Ingress jupyter';
COMMENT
ON COLUMN kubernetes_resources.ingress_jupyter_protocol IS 'Ingress jupyter protocol';
COMMENT
ON COLUMN kubernetes_resources.ingress_jupyter_cname IS 'Ingress jupyter cname';
COMMENT
ON COLUMN kubernetes_resources.ingress_jupyter_url IS 'Ingress jupyter url';
COMMENT
ON COLUMN kubernetes_resources.ingress_ray_name IS 'Ingress ray';
COMMENT
ON COLUMN kubernetes_resources.ingress_ray_protocol IS 'Ingress ray protocol';
COMMENT
ON COLUMN kubernetes_resources.ingress_ray_cname IS 'Ingress ray cname';
COMMENT
ON COLUMN kubernetes_resources.ingress_ray_url IS 'Ingress ray url';
COMMENT
ON COLUMN kubernetes_resources.status IS 'Status: -1-stopped, 0-starting, 1-running';
COMMENT
ON COLUMN kubernetes_resources.created_at IS 'Create time';
COMMENT
ON COLUMN kubernetes_resources.updated_at IS 'Update time';

-- KeyPair
DROP TABLE IF EXISTS key_pair_resources;
CREATE TABLE key_pair_resources
(
    id          BIGSERIAL PRIMARY KEY,
    external_id VARCHAR(255) NOT NULL UNIQUE,
    private_key TEXT,
    public_key  TEXT,
    status      INT,
    created_at  TIMESTAMP(3) NOT NULL,
    updated_at  TIMESTAMP(3)
);

COMMENT
ON TABLE key_pair_resources IS 'KeyPair Resources';
COMMENT
ON COLUMN key_pair_resources.external_id IS 'External Business Primary Key';
COMMENT
ON COLUMN key_pair_resources.private_key IS 'Private Key';
COMMENT
ON COLUMN key_pair_resources.public_key IS 'Public Key';
COMMENT
ON COLUMN key_pair_resources.status IS 'Status: 0-inactive, 1-active';
COMMENT
ON COLUMN key_pair_resources.created_at IS 'Create time';
COMMENT
ON COLUMN key_pair_resources.updated_at IS 'Update time';

-- Instance
DROP TABLE IF EXISTS instance_resources;
CREATE TABLE instance_resources
(
    id                BIGSERIAL PRIMARY KEY,
    external_id       VARCHAR(255) NOT NULL UNIQUE,
    instance_type     VARCHAR(50)  NOT NULL,
    image_id          VARCHAR(100) NOT NULL,
    device_name       VARCHAR(50)  NOT NULL,
    volume_size       INT          NOT NULL,
    key_pair_id       VARCHAR(50)  NOT NULL,
    security_group_id VARCHAR(50)  NOT NULL,
    private_key       TEXT,
    public_key        TEXT,
    instance_id       VARCHAR(255) NOT NULL UNIQUE,
    association_id    VARCHAR(255),
    allocation_id     VARCHAR(255),
    ip_address        VARCHAR(255),
    domain_name       VARCHAR(255),
    status            INT,
    created_at        TIMESTAMP(3) NOT NULL,
    updated_at        TIMESTAMP(3)
);

COMMENT
ON TABLE instance_resources IS 'Instance Resources';
COMMENT
ON COLUMN instance_resources.external_id IS 'External Business Primary Key';
COMMENT
ON COLUMN instance_resources.instance_type IS 'Instance Type';
COMMENT
ON COLUMN instance_resources.image_id IS 'Image ID(AMI)';
COMMENT
ON COLUMN instance_resources.device_name IS 'Device Name';
COMMENT
ON COLUMN instance_resources.volume_size IS 'Volume Size';
COMMENT
ON COLUMN instance_resources.key_pair_id IS 'Key Pair Id';
COMMENT
ON COLUMN instance_resources.security_group_id IS 'Security Group Id';
COMMENT
ON COLUMN instance_resources.private_key IS 'Private Key';
COMMENT
ON COLUMN instance_resources.public_key IS 'Public Key';
COMMENT
ON COLUMN instance_resources.instance_id IS 'Instance Id';
COMMENT
ON COLUMN instance_resources.allocation_id IS 'Allocation Id';
COMMENT
ON COLUMN instance_resources.ip_address IS 'IP Address';
COMMENT
ON COLUMN instance_resources.domain_name IS 'Domain Name';
COMMENT
ON COLUMN instance_resources.status IS 'Status: -1:failed,0:starting,1:running,2:stopped,3:terminated';
COMMENT
ON COLUMN instance_resources.created_at IS 'Create time';
COMMENT
ON COLUMN instance_resources.updated_at IS 'Update time';

-- Model type
DROP TABLE IF EXISTS instance_model_type;
CREATE TABLE instance_model_type
(
    id               BIGSERIAL PRIMARY KEY,
    type             INT          NOT NULL,
    instance_type    varchar(50)  NOT NULL,
    cpu_count        INT          NOT NULL,
    memory           INT          NOT NULL,
    gpu_sku          varchar(50)  NOT NULL,
    gpu_count        INT          NOT NULL,
    gpu_memory       INT          NOT NULL,
    storage_type     INT          NOT NULL,
    storage_capacity BIGINT       NOT NULL,
    display_name     VARCHAR(255) NOT NULL
);

COMMENT
ON COLUMN instance_model_type.type IS 'instance machine type: 1-CPU, 2-GPU';
COMMENT
ON COLUMN instance_model_type.instance_type IS 'instance aws type: g4dn.12xlarge';
COMMENT
ON COLUMN instance_model_type.cpu_count IS 'CPU count';
COMMENT
ON COLUMN instance_model_type.memory IS 'Memory';
COMMENT
ON COLUMN instance_model_type.gpu_sku IS 'GPU SKU';
COMMENT
ON COLUMN instance_model_type.gpu_count IS 'GPU count';
COMMENT
ON COLUMN instance_model_type.gpu_memory IS 'GPU Memory';
COMMENT
ON COLUMN instance_model_type.storage_type IS 'Storage Type';
COMMENT
ON COLUMN instance_model_type.storage_capacity IS 'Storage Capacity';
COMMENT
ON COLUMN instance_model_type.display_name IS 'Display Name';