CREATE SCHEMA IF NOT EXISTS woodpecker;

-- Agents
CREATE TABLE woodpecker.agents (
    id BIGSERIAL PRIMARY KEY,
    created BIGINT,
    updated BIGINT,
    name VARCHAR(255),
    owner_id BIGINT,
    token VARCHAR(255),
    last_contact BIGINT,
    last_work BIGINT,
    platform VARCHAR(100),
    backend VARCHAR(100),
    capacity INTEGER,
    version VARCHAR(255),
    no_schedule BOOLEAN,
    custom_labels JSON,
    org_id BIGINT
);

-- Configs
CREATE TABLE woodpecker.configs (
    id BIGSERIAL PRIMARY KEY,
    repo_id BIGINT,
    hash VARCHAR(255),
    name VARCHAR(255),
    data BYTEA
);

-- Crons
CREATE TABLE woodpecker.crons (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255),
    repo_id BIGINT,
    creator_id BIGINT,
    next_exec BIGINT,
    schedule VARCHAR(255) NOT NULL,
    created BIGINT DEFAULT 0 NOT NULL,
    branch VARCHAR(255)
);

-- Forges
CREATE TABLE woodpecker.forges (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(250),
    url VARCHAR(500),
    oauth_client_id VARCHAR(250),
    oauth_client_secret VARCHAR(250),
    skip_verify BOOLEAN,
    oauth_host VARCHAR(250),
    additional_options JSON
);

-- Log entries
CREATE TABLE woodpecker.log_entries (
    id BIGSERIAL PRIMARY KEY,
    step_id BIGINT,
    "time" BIGINT,
    line INTEGER,
    data BYTEA,
    created BIGINT,
    type INTEGER
);

-- Migration
CREATE TABLE woodpecker.migration (
    id VARCHAR(255) PRIMARY KEY,
    description VARCHAR(255)
);

-- Orgs
CREATE TABLE woodpecker.orgs (
    id BIGSERIAL PRIMARY KEY,
    forge_id BIGINT,
    name VARCHAR(255),
    is_user BOOLEAN,
    private BOOLEAN
);

-- Perms
CREATE TABLE woodpecker.perms (
    user_id BIGINT NOT NULL,
    repo_id BIGINT NOT NULL,
    pull BOOLEAN,
    push BOOLEAN,
    admin BOOLEAN,
    synced BIGINT,
    created BIGINT,
    updated BIGINT,
    PRIMARY KEY(user_id, repo_id)
);

-- Pipeline configs
CREATE TABLE woodpecker.pipeline_configs (
    config_id BIGINT NOT NULL,
    pipeline_id BIGINT NOT NULL,
    PRIMARY KEY(config_id, pipeline_id)
);

-- Pipelines
CREATE TABLE woodpecker.pipelines (
    id BIGSERIAL PRIMARY KEY,
    repo_id BIGINT,
    number BIGINT,
    author VARCHAR(255),
    parent BIGINT,
    event VARCHAR(255),
    event_reason JSON,
    status VARCHAR(255),
    errors JSON,
    created BIGINT DEFAULT 0 NOT NULL,
    updated BIGINT DEFAULT 0 NOT NULL,
    started BIGINT,
    finished BIGINT,
    deploy VARCHAR(255),
    deploy_task VARCHAR(255),
    commit VARCHAR(255),
    branch VARCHAR(255),
    ref VARCHAR(255),
    refspec VARCHAR(255),
    title VARCHAR(255),
    message TEXT,
    "timestamp" BIGINT,
    sender VARCHAR(255),
    avatar VARCHAR(500),
    email VARCHAR(500),
    forge_url VARCHAR(255),
    reviewer VARCHAR(255),
    reviewed BIGINT,
    changed_files TEXT,
    additional_variables JSON,
    pr_labels JSON,
    pr_milestone VARCHAR(255),
    is_prerelease BOOLEAN,
    from_fork BOOLEAN
);

-- Redirections
CREATE TABLE woodpecker.redirections (
    id BIGSERIAL PRIMARY KEY,
    repo_id BIGINT,
    repo_full_name VARCHAR(255)
);

-- Registries
CREATE TABLE woodpecker.registries (
    id BIGSERIAL PRIMARY KEY,
    org_id BIGINT DEFAULT 0 NOT NULL,
    repo_id BIGINT DEFAULT 0 NOT NULL,
    address VARCHAR(255) NOT NULL,
    username VARCHAR(2000),
    password TEXT
);

-- Repos
CREATE TABLE woodpecker.repos (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    forge_id BIGINT,
    forge_remote_id VARCHAR(255),
    org_id BIGINT,
    owner VARCHAR(255),
    name VARCHAR(255),
    full_name VARCHAR(255),
    avatar VARCHAR(500),
    forge_url VARCHAR(1000),
    clone VARCHAR(1000),
    clone_ssh VARCHAR(1000),
    branch VARCHAR(500),
    pr_enabled BOOLEAN DEFAULT true,
    timeout BIGINT,
    visibility VARCHAR(10),
    private BOOLEAN,
    trusted JSON,
    require_approval VARCHAR(50),
    approval_allowed_users JSON,
    active BOOLEAN,
    allow_pr BOOLEAN,
    allow_deploy BOOLEAN,
    config_path VARCHAR(500),
    hash VARCHAR(500),
    cancel_previous_pipeline_events JSON,
    netrc_trusted JSON,
    config_extension_endpoint VARCHAR(500)
);

-- Secrets
CREATE TABLE woodpecker.secrets (
    id BIGSERIAL PRIMARY KEY,
    org_id BIGINT DEFAULT 0 NOT NULL,
    repo_id BIGINT DEFAULT 0 NOT NULL,
    name VARCHAR(255) NOT NULL,
    value TEXT,
    images JSON,
    events JSON
);

-- Server configs
CREATE TABLE woodpecker.server_configs (
    key VARCHAR(255) PRIMARY KEY,
    value VARCHAR(255)
);

-- Steps
CREATE TABLE woodpecker.steps (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR(255),
    pipeline_id BIGINT,
    pid INTEGER,
    ppid INTEGER,
    name VARCHAR(255),
    state VARCHAR(255),
    error TEXT,
    failure VARCHAR(255),
    exit_code INTEGER,
    started BIGINT,
    finished BIGINT,
    type VARCHAR(255)
);

-- Tasks
CREATE TABLE woodpecker.tasks (
    id VARCHAR(255) PRIMARY KEY,
    pid INTEGER,
    name VARCHAR(255),
    data BYTEA,
    labels JSON,
    dependencies JSON,
    run_on JSON,
    dependencies_status JSON,
    agent_id BIGINT,
    pipeline_id BIGINT,
    repo_id BIGINT
);

-- Workflows
CREATE TABLE woodpecker.workflows (
    id BIGSERIAL PRIMARY KEY,
    pipeline_id BIGINT,
    pid INTEGER,
    name VARCHAR(255),
    state VARCHAR(255),
    error TEXT,
    started BIGINT,
    finished BIGINT,
    agent_id BIGINT,
    platform VARCHAR(255),
    environ JSON,
    axis_id INTEGER
);

CREATE INDEX "IDX_agents_org_id" ON public.agents USING btree (org_id);

CREATE INDEX "IDX_crons_creator_id" ON public.crons USING btree (creator_id);

CREATE INDEX "IDX_crons_name" ON public.crons USING btree (name);

CREATE INDEX "IDX_crons_repo_id" ON public.crons USING btree (repo_id);

CREATE INDEX "IDX_log_entries_step_id" ON public.log_entries USING btree (step_id);

CREATE INDEX "IDX_perms_repo_id" ON public.perms USING btree (repo_id);

CREATE INDEX "IDX_perms_user_id" ON public.perms USING btree (user_id);

CREATE INDEX "IDX_pipelines_author" ON public.pipelines USING btree (author);

CREATE INDEX "IDX_pipelines_repo_id" ON public.pipelines USING btree (repo_id);

CREATE INDEX "IDX_pipelines_status" ON public.pipelines USING btree (status);

CREATE INDEX "IDX_registries_address" ON public.registries USING btree (address);

CREATE INDEX "IDX_registries_org_id" ON public.registries USING btree (org_id);

CREATE INDEX "IDX_registries_repo_id" ON public.registries USING btree (repo_id);

CREATE INDEX "IDX_repos_org_id" ON public.repos USING btree (org_id);

CREATE INDEX "IDX_repos_user_id" ON public.repos USING btree (user_id);

CREATE INDEX "IDX_secrets_name" ON public.secrets USING btree (name);

CREATE INDEX "IDX_secrets_org_id" ON public.secrets USING btree (org_id);

CREATE INDEX "IDX_secrets_repo_id" ON public.secrets USING btree (repo_id);

CREATE INDEX "IDX_steps_pipeline_id" ON public.steps USING btree (pipeline_id);

CREATE INDEX "IDX_steps_uuid" ON public.steps USING btree (uuid);

CREATE INDEX "IDX_workflows_pipeline_id" ON public.workflows USING btree (pipeline_id);

CREATE UNIQUE INDEX "UQE_configs_s" ON public.configs USING btree (repo_id, hash, name);

CREATE UNIQUE INDEX "UQE_crons_s" ON public.crons USING btree (name, repo_id);

CREATE UNIQUE INDEX "UQE_orgs_s" ON public.orgs USING btree (forge_id, name);

CREATE UNIQUE INDEX "UQE_perms_s" ON public.perms USING btree (user_id, repo_id);

CREATE UNIQUE INDEX "UQE_pipeline_configs_s" ON public.pipeline_configs USING btree (config_id, pipeline_id);

CREATE UNIQUE INDEX "UQE_pipelines_s" ON public.pipelines USING btree (repo_id, number);

CREATE UNIQUE INDEX "UQE_redirections_repo_full_name" ON public.redirections USING btree (repo_full_name);

CREATE UNIQUE INDEX "UQE_registries_s" ON public.registries USING btree (org_id, repo_id, address);

CREATE UNIQUE INDEX "UQE_repos_full_name" ON public.repos USING btree (full_name);

CREATE UNIQUE INDEX "UQE_repos_name" ON public.repos USING btree (owner, name);

CREATE UNIQUE INDEX "UQE_secrets_s" ON public.secrets USING btree (org_id, repo_id, name);

CREATE UNIQUE INDEX "UQE_steps_s" ON public.steps USING btree (pipeline_id, pid);

CREATE UNIQUE INDEX "UQE_tasks_id" ON public.tasks USING btree (id);

CREATE UNIQUE INDEX "UQE_users_hash" ON public.users USING btree (hash);

CREATE UNIQUE INDEX "UQE_users_login" ON public.users USING btree (login);

CREATE UNIQUE INDEX "UQE_workflows_s" ON public.workflows USING btree (pipeline_id, pid);
