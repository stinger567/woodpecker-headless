CREATE SCHEMA IF NOT EXISTS woodpecker;

-- Agents
CREATE TABLE woodpecker.agents (
    id bigserial PRIMARY KEY,
    created timestamptz DEFAULT NOW() NOT NULL,
    updated timestamptz,
    name varchar(255),
    owner_id bigint,
    token varchar(255),
    last_contact bigint,
    last_work bigint,
    platform varchar(100),
    backend varchar(100),
    capacity integer,
    version varchar(255),
    no_schedule boolean,
    custom_labels json,
    org_id bigint
);

-- Configs
CREATE TABLE woodpecker.configs (
    id bigserial PRIMARY KEY,
    repo_id uuid,
    hash VARCHAR(255),
    name varchar(255),
    data bytea
);

-- Crons
CREATE TABLE woodpecker.crons (
    id bigserial PRIMARY KEY,
    name varchar(255),
    repo_id uuid,
    creator_id uuid,
    next_exec bigint,
    schedule varchar(255) NOT NULL,
    created timestamptz DEFAULT NOW() NOT NULL,
    branch varchar(255),
    internal boolean
);

-- Log entries
CREATE TABLE woodpecker.log_entries (
    id bigserial PRIMARY KEY,
    step_id bigint,
    "time" bigint,
    line integer,
    data bytea,
    created timestamptz DEFAULT NOW() NOT NULL,
    type INTEGER
);

-- Perms
CREATE TABLE woodpecker.perms (
    user_id uuid NOT NULL,
    repo_id uuid NOT NULL,
    pull boolean,
    push boolean,
    admin boolean,
    synced bigint,
    created bigint,
    updated bigint,
    PRIMARY KEY (user_id, repo_id)
);

-- Pipeline configs
CREATE TABLE woodpecker.pipeline_configs (
    config_id bigint NOT NULL,
    pipeline_id bigint NOT NULL,
    PRIMARY KEY (config_id, pipeline_id)
);

-- Pipelines
CREATE TABLE woodpecker.pipelines (
    id bigserial PRIMARY KEY,
    repo_id uuid,
    number bigint,
    author varchar(255),
    parent bigint,
    event varchar(255),
    event_reason json,
    status varchar(255),
    errors json,
    created timestamptz DEFAULT NOW() NOT NULL,
    updated timestamptz DEFAULT NOW() NOT NULL,
    started timestamptz,
    finished timestamptz,
    deploy varchar(255),
    deploy_task varchar(255),
    commit VARCHAR(255),
    branch varchar(255),
    ref varchar(255),
    refspec varchar(255),
    title varchar(255),
    message text,
    "timestamp" timestamptz,
    sender varchar(255),
    avatar varchar(500),
    email varchar(500),
    forge_url varchar(255),
    reviewer varchar(255),
    reviewed bigint,
    changed_files text,
    additional_variables json,
    pr_labels json,
    pr_milestone varchar(255),
    is_prerelease boolean,
    from_fork boolean,
    internal boolean
);

-- Redirections
CREATE TABLE woodpecker.redirections (
    id bigserial PRIMARY KEY,
    repo_id uuid,
    repo_full_name varchar(255)
);

-- Registries
CREATE TABLE woodpecker.registries (
    id bigserial PRIMARY KEY,
    org_id bigint DEFAULT 0 NOT NULL,
    repo_id uuid DEFAULT "" NOT NULL,
    address varchar(255) NOT NULL,
    username varchar(2000),
    password TEXT
);

-- Repos
CREATE TABLE woodpecker.repos (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    forge_account_id uuid NOT NULL REFERENCES oauth.account (id) ON DELETE CASCADE,
    forge_id uuid NOT NULL REFERENCES oauth.providers (id) ON DELETE CASCADE,
    forge_remote_id varchar UNIQUE NOT NULL,
    org_id bigint,
    owner varchar,
    name varchar,
    full_name varchar,
    avatar varchar,
    forge_url varchar,
    clone varchar,
    clone_ssh varchar,
    branch varchar,
    pr_enabled boolean DEFAULT TRUE,
    timeout bigint,
    visibility varchar,
    private boolean,
    trusted json,
    require_approval varchar,
    approval_allowed_users json,
    active boolean,
    allow_pr boolean,
    allow_deploy boolean,
    config_path varchar,
    hash varchar,
    cancel_previous_pipeline_events json,
    netrc_trusted json,
    config_extension_endpoint varchar
);

-- Secrets
CREATE TABLE woodpecker.secrets (
    id bigserial PRIMARY KEY,
    org_id uuid NOT NULL,
    repo_id uuid NOT NULL,
    name varchar(255) NOT NULL,
    value text,
    images json,
    events json
);

-- Server configs
CREATE TABLE woodpecker.server_configs (
    key VARCHAR(255) PRIMARY KEY,
    value varchar(255)
);

-- Steps
CREATE TABLE woodpecker.steps (
    id bigserial PRIMARY KEY,
    uuid varchar(255),
    pipeline_id bigint,
    pid integer,
    ppid integer,
    name varchar(255),
    state varchar(255),
    error text,
    failure varchar(255),
    exit_code integer,
    started timestamptz,
    finished timestamptz,
    type VARCHAR(255)
);

-- Tasks
CREATE TABLE woodpecker.tasks (
    id varchar(255) PRIMARY KEY,
    pid integer,
    name varchar(255),
    data bytea,
    labels json,
    dependencies json,
    run_on json,
    dependencies_status json,
    agent_id bigint,
    pipeline_id bigint,
    repo_id bigint
);

-- Workflows
CREATE TABLE woodpecker.workflows (
    id bigserial PRIMARY KEY,
    pipeline_id bigint,
    pid integer,
    name varchar(255),
    state varchar(255),
    error text,
    started timestamptz,
    finished timestamptz,
    agent_id bigint,
    platform varchar(255),
    environ json,
    axis_id integer
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

CREATE UNIQUE INDEX "UQE_repos_name" ON public.repos USING btree (OWNER, name);

CREATE UNIQUE INDEX "UQE_secrets_s" ON public.secrets USING btree (org_id, repo_id, name);

CREATE UNIQUE INDEX "UQE_steps_s" ON public.steps USING btree (pipeline_id, pid);

CREATE UNIQUE INDEX "UQE_tasks_id" ON public.tasks USING btree (id);

CREATE UNIQUE INDEX "UQE_users_hash" ON public.users USING btree (hash);

CREATE UNIQUE INDEX "UQE_users_login" ON public.users USING btree (login);

CREATE UNIQUE INDEX "UQE_workflows_s" ON public.workflows USING btree (pipeline_id, pid);

