CREATE SCHEMA IF NOT EXISTS synapse;

SET
search_path TO synapse;

-- Kubernetes 资源表
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
ON TABLE kubernetes_resources IS 'Kubernetes 资源表';
COMMENT
ON COLUMN kubernetes_resources.external_id IS '外部业务主键';
COMMENT
ON COLUMN kubernetes_resources.namespace IS 'Kubernetes 命名空间';
COMMENT
ON COLUMN kubernetes_resources.deployment_name IS 'Deployment 名称';
COMMENT
ON COLUMN kubernetes_resources.deployment_replicas IS 'Deployment 副本数';
COMMENT
ON COLUMN kubernetes_resources.service_name IS 'Service 名称';
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
ON COLUMN kubernetes_resources.status IS 'Status: -1-终止, 0-启动中, 1-正常';
COMMENT
ON COLUMN kubernetes_resources.created_at IS '创建时间';
COMMENT
ON COLUMN kubernetes_resources.updated_at IS '更新时间';

-- KeyPair表
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
ON TABLE key_pair_resources IS 'KeyPair 资源表';
COMMENT
ON COLUMN key_pair_resources.external_id IS '外部业务主键';
COMMENT
ON COLUMN key_pair_resources.private_key IS '私钥';
COMMENT
ON COLUMN key_pair_resources.public_key IS '公钥';
COMMENT
ON COLUMN key_pair_resources.status IS 'Status: 0-停用, 1-正常';
COMMENT
ON COLUMN key_pair_resources.created_at IS '创建时间';
COMMENT
ON COLUMN key_pair_resources.updated_at IS '更新时间';

-- Instance表
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
ON TABLE instance_resources IS 'Instance 资源表';
COMMENT
ON COLUMN instance_resources.external_id IS '外部业务主键';
COMMENT
ON COLUMN instance_resources.instance_type IS '实例类型';
COMMENT
ON COLUMN instance_resources.image_id IS '镜像ID';
COMMENT
ON COLUMN instance_resources.device_name IS '设备名称';
COMMENT
ON COLUMN instance_resources.volume_size IS '磁盘大小';
COMMENT
ON COLUMN instance_resources.key_pair_id IS '密钥对id';
COMMENT
ON COLUMN instance_resources.security_group_id IS '安全组ID';
COMMENT
ON COLUMN instance_resources.private_key IS '私钥';
COMMENT
ON COLUMN instance_resources.public_key IS '公钥';
COMMENT
ON COLUMN instance_resources.instance_id IS '实例ID';
COMMENT
ON COLUMN instance_resources.allocation_id IS '分配ID';
COMMENT
ON COLUMN instance_resources.ip_address IS 'IP地址';
COMMENT
ON COLUMN instance_resources.domain_name IS '域名';
COMMENT
ON COLUMN instance_resources.status IS 'Status: -1:失败,0:启动中,1:正在运行,2:已经停止,3:已终止';
COMMENT
ON COLUMN instance_resources.created_at IS '创建时间';
COMMENT
ON COLUMN instance_resources.updated_at IS '更新时间';

-- 实例模型类型表(todo: 暂不使用)
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
ON COLUMN instance_model_type.type IS '类型: 1-CPU, 2-GPU';
COMMENT
ON COLUMN instance_model_type.instance_type IS '实例类型: g4dn.12xlarge';
COMMENT
ON COLUMN instance_model_type.cpu_count IS 'CPU 核心数';
COMMENT
ON COLUMN instance_model_type.memory IS '内存大小';
COMMENT
ON COLUMN instance_model_type.gpu_sku IS 'GPU的SKU';
COMMENT
ON COLUMN instance_model_type.gpu_count IS 'GPU 数量';
COMMENT
ON COLUMN instance_model_type.gpu_memory IS 'GPU 内存';
COMMENT
ON COLUMN instance_model_type.storage_type IS '存储类型';
COMMENT
ON COLUMN instance_model_type.storage_capacity IS '存储容量';
COMMENT
ON COLUMN instance_model_type.display_name IS '显示名称';